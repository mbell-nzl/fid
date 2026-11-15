// Package epicrecord reads the shared eligibility rule and person's epic close.
package epicrecord

import (
	"database/sql"
	"fmt"

	"github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issuestate"
)

type Eligibility struct {
	ID, Kind, State string
	Revision        int
	Children        []issuestate.ChildTaskStatus
}

// CanClose includes empty epics; status additionally requires at least one task.
func (e Eligibility) CanClose() bool {
	return e.Kind == "epic" && (e.State == "todo" || e.State == "in_progress") && issuestate.AllChildTasksDoneOrCancelled(e.Children).Satisfied
}
func (e Eligibility) OpenTasks() []string {
	result := issuestate.AllChildTasksDoneOrCancelled(e.Children)
	if result.Satisfied {
		return []string{}
	}
	return result.Unmet[0].EntityIDs
}
func ReadEligibility(q db.Queryer, id string) (Eligibility, error) {
	e := Eligibility{ID: id}
	rows, err := q.Query(`SELECT e.kind,e.status,e.revision,t.id,t.status FROM issues e LEFT JOIN issues t ON t.parent_id=e.id AND t.kind='task' WHERE e.id=? ORDER BY t.id`, id)
	if err != nil {
		return e, fmt.Errorf("read epic tasks: %w", err)
	}
	defer func() { _ = rows.Close() }()
	found := false
	for rows.Next() {
		var childID, state sql.NullString
		if err := rows.Scan(&e.Kind, &e.State, &e.Revision, &childID, &state); err != nil {
			return e, err
		}
		found = true
		if childID.Valid {
			e.Children = append(e.Children, issuestate.ChildTaskStatus{IssueID: childID.String, State: state.String})
		}
	}
	if err := rows.Err(); err != nil {
		return e, err
	}
	if !found {
		return e, sql.ErrNoRows
	}
	return e, nil
}

type Decision struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Actor       string `json:"actor"`
	SubjectType string `json:"subject_type"`
	SubjectID   string `json:"subject_id"`
	Reason      string `json:"reason"`
	BindingID   string `json:"binding_id"`
	At          string `json:"decided_at"`
}
type Result struct {
	Replayed bool     `json:"replayed"`
	Status   string   `json:"status"`
	Decision Decision `json:"decision"`
	Revision int      `json:"revision"`
	EventID  string   `json:"-"`
}

func ReadDecision(q interface{ QueryRow(string, ...any) *sql.Row }, id string) (Decision, error) {
	var d Decision
	err := q.QueryRow(`SELECT id,kind,actor,subject_type,subject_id,reason,binding_id,decided_at FROM person_decisions WHERE id=? AND kind='close-epic'`, id).Scan(&d.ID, &d.Kind, &d.Actor, &d.SubjectType, &d.SubjectID, &d.Reason, &d.BindingID, &d.At)
	return d, err
}
