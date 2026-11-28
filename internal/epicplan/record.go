// Package epicplan reads epic plans and their immutable validation evidence.
// It has no journal dependency, so admission uses the same readers as commands.
package epicplan

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/db"
)

type Queryer interface {
	db.Queryer
	QueryRow(string, ...any) *sql.Row
}

type Validation struct {
	ID         string `json:"id"`
	Ordinal    int    `json:"ordinal"`
	Generation int    `json:"generation"`
	Claim      string `json:"claim"`
	Verifier   string `json:"verifier"`
	Recipe     string `json:"recipe,omitempty"`
}

type Binding struct {
	Digest     string `json:"digest"`
	Generation int    `json:"generation"`
}

type Approval struct {
	State   string  `json:"state"`
	ID      string  `json:"decision_id,omitempty"`
	Actor   string  `json:"actor,omitempty"`
	At      string  `json:"at,omitempty"`
	Reason  string  `json:"reason,omitempty"`
	Binding Binding `json:"binding"`
}

type Task struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Title    string `json:"title"`
	State    string `json:"state"`
	Lane     string `json:"lane,omitempty"`
	Revision int    `json:"revision"`
	Accept   bool   `json:"accept"`
}

type Plan struct {
	EpicID      string       `json:"epic_id"`
	State       string       `json:"state"`
	Revision    int          `json:"revision"`
	Tasks       []Task       `json:"tasks"`
	Validations []Validation `json:"validations"`
	Binding     Binding      `json:"binding"`
}

type Pass struct {
	ID           string                  `json:"id"`
	EpicID       string                  `json:"epic_id"`
	ValidationID string                  `json:"validation_id"`
	Generation   int                     `json:"generation"`
	Commit       string                  `json:"commit"`
	Merges       []string                `json:"merges"`
	Actor        string                  `json:"actor"`
	At           string                  `json:"at"`
	Outcome      string                  `json:"outcome"`
	DecisionID   string                  `json:"decision_id,omitempty"`
	Check        *checkrecord.Result     `json:"check,omitempty"`
	Definition   *checkrecord.Definition `json:"definition,omitempty"`
}

type View struct {
	Validation
	State  string   `json:"state"`
	Pass   *Pass    `json:"pass,omitempty"`
	Rounds []string `json:"rounds,omitempty"`
}

func Read(q db.Queryer, epic string) ([]Validation, int, int, error) {
	rows, err := q.Query(`SELECT c.generation,c.ordinal,c.claim,c.verifier,c.recipe,c.removed FROM epic_validation_changes c WHERE c.epic_id=? AND c.generation=(SELECT max(n.generation) FROM epic_validation_changes n WHERE n.epic_id=c.epic_id AND n.ordinal=c.ordinal) ORDER BY c.ordinal`, epic)
	if err != nil {
		return nil, 0, 0, err
	}
	defer func() { _ = rows.Close() }()
	list := []Validation{}
	generation, ordinal := 0, 0
	for rows.Next() {
		var v Validation
		var removed bool
		if err := rows.Scan(&v.Generation, &v.Ordinal, &v.Claim, &v.Verifier, &v.Recipe, &removed); err != nil {
			return nil, 0, 0, err
		}
		generation = max(generation, v.Generation)
		ordinal = max(ordinal, v.Ordinal)
		v.ID = fmt.Sprintf("V%d", v.Ordinal)
		if !removed {
			list = append(list, v)
		}
	}
	return list, generation, ordinal, rows.Err()
}

func ListBinding(list []Validation, generation int) Binding {
	// Only the validation list is signed; task/spec edits do not invalidate it.
	raw, _ := json.Marshal(list)
	return Binding{Digest: fmt.Sprintf("sha256:%x", sha256.Sum256(raw)), Generation: generation}
}

func Current(q db.Queryer, epic string) (Approval, error) {
	list, generation, _, err := Read(q, epic)
	if err != nil {
		return Approval{}, err
	}
	rows, err := q.Query(`SELECT id,actor,reason,decided_at,binding_json FROM person_decisions WHERE kind='approve-epic' AND subject_id=? ORDER BY rowid DESC LIMIT 1`, epic)
	if err != nil {
		return Approval{}, err
	}
	defer func() { _ = rows.Close() }()
	a := Approval{State: "missing"}
	if rows.Next() {
		var raw string
		if err := rows.Scan(&a.ID, &a.Actor, &a.Reason, &a.At, &raw); err != nil {
			return a, err
		}
		if err := json.Unmarshal([]byte(raw), &a.Binding); err != nil {
			return a, err
		}
		a.State = "stale"
		if a.Binding == ListBinding(list, generation) {
			a.State = "current"
		}
	}
	return a, rows.Err()
}

func Preview(q Queryer, epic string) (Plan, error) {
	p := Plan{EpicID: epic, Tasks: []Task{}}
	if err := q.QueryRow(`SELECT status,revision FROM issues WHERE id=? AND kind='epic'`, epic).Scan(&p.State, &p.Revision); err != nil {
		return p, err
	}
	var err error
	var generation int
	p.Validations, generation, _, err = Read(q, epic)
	if err != nil {
		return p, err
	}
	p.Binding = ListBinding(p.Validations, generation)
	rows, err := q.Query(`SELECT i.id,i.kind,i.title,i.status,coalesce(i.lane,''),i.revision,
 coalesce(c.actor,'') LIKE 'agent:%' AND s.category NOT IN ('completed','cancelled')
 AND NOT EXISTS(SELECT 1 FROM person_decisions d WHERE d.kind='accept' AND d.subject_id=i.id)
 AND NOT EXISTS(SELECT 1 FROM epic_plan_acceptances a WHERE a.issue_id=i.id)
 FROM issues i LEFT JOIN issues parent ON parent.id=i.parent_id
 JOIN workflow_states s ON s.workflow_id=i.workflow_id AND s.state=i.status
 LEFT JOIN issue_creators c ON c.issue_id=i.id
 WHERE (i.kind='task' AND i.parent_id=?) OR (i.kind='subtask' AND parent.parent_id=?) ORDER BY i.kind,i.sequence IS NULL,i.sequence,i.id`, epic, epic)
	if err != nil {
		return p, err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Kind, &t.Title, &t.State, &t.Lane, &t.Revision, &t.Accept); err != nil {
			return p, err
		}
		p.Tasks = append(p.Tasks, t)
	}
	return p, rows.Err()
}

func Normalize(v Validation) (Validation, error) {
	v.Claim = strings.TrimSpace(v.Claim)
	v.Verifier = strings.TrimSpace(v.Verifier)
	v.Recipe = strings.TrimSpace(v.Recipe)
	if v.Claim == "" {
		return v, fmt.Errorf("validation claim is required")
	}
	if v.Recipe != "" {
		if strings.HasPrefix(v.Recipe, "/") || strings.ContainsAny(v.Recipe, "\\\x00") {
			return v, fmt.Errorf("recipe must be project-relative")
		}
		v.Recipe = path.Clean(v.Recipe)
		if v.Recipe == "." || v.Recipe == ".." || strings.HasPrefix(v.Recipe, "../") {
			return v, fmt.Errorf("recipe must be project-relative")
		}
	}
	if v.Verifier == "human" {
		if v.Recipe == "" {
			return v, fmt.Errorf("human validation requires a recipe")
		}
		return v, nil
	}
	kind, name, ok := strings.Cut(v.Verifier, ":")
	if !ok || name == "" || kind != "check" && kind != "review" {
		return v, fmt.Errorf("verifier must be check:<name>, review:<profile> or human")
	}
	return v, nil
}

func Merges(q db.Queryer, epic string) ([]string, error) {
	rows, err := q.Query(`SELECT m.record_json FROM issues i LEFT JOIN completions c ON c.task_id=i.id LEFT JOIN merge_observations m ON m.id=c.merge_id WHERE i.parent_id=? AND i.kind='task' AND i.status='done' ORDER BY i.id`, epic)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	merges := []string{}
	for rows.Next() {
		var raw sql.NullString
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		var r struct {
			MergeCommit string `json:"merge_commit"`
		}
		if !raw.Valid {
			return nil, fmt.Errorf("completed task has no completion receipt")
		}
		if err := json.Unmarshal([]byte(raw.String), &r); err != nil {
			return nil, err
		}
		if r.MergeCommit == "" {
			return nil, fmt.Errorf("completed task has no merge commit")
		}
		merges = append(merges, r.MergeCommit)
	}
	return merges, rows.Err()
}

type Result struct {
	Replayed   bool        `json:"replayed"`
	EpicID     string      `json:"epic_id"`
	Validation *Validation `json:"validation,omitempty"`
	DecisionID string      `json:"decision_id,omitempty"`
	Accepted   []string    `json:"accepted,omitempty"`
	Checks     []Pass      `json:"checks"`
	Warnings   []string    `json:"warnings,omitempty"`
	At         string      `json:"occurred_at"`
	EventID    string      `json:"-"`
}
