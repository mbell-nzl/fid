package epicplan

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/mbell-nzl/fid/internal/db"
)

// Views evaluates immutable results against today's validation revisions and
// completed task merges. A failure replaces an earlier pass for that revision.
func Views(q Queryer, epic string, roots ...string) ([]View, error) {
	list, _, _, err := Read(q, epic)
	if err != nil {
		return nil, err
	}
	merges, mergeErr := Merges(q, epic)
	closed, err := readCloseBinding(q, epic)
	if err != nil {
		return nil, err
	}
	views := []View{}
	for _, v := range list {
		view := View{Validation: v, State: "missing"}
		var raw string
		err := q.QueryRow(`SELECT record_json FROM epic_validation_results WHERE epic_id=? AND ordinal=? ORDER BY rowid DESC LIMIT 1`, epic, v.Ordinal).Scan(&raw)
		if err == nil {
			var p Pass
			if err := json.Unmarshal([]byte(raw), &p); err != nil {
				return nil, err
			}
			view.Pass = &p
			view.State = "failed"
			if p.Outcome == "passed" {
				view.State = "current"
			}
			if p.Generation != v.Generation || mergeErr != nil || !coversMerges(p, merges, closed, roots) {
				view.State = "stale"
			}
		} else if err != sql.ErrNoRows {
			return nil, err
		}
		if profile, ok := strings.CutPrefix(v.Verifier, "review:"); ok {
			rows, err := q.Query(`SELECT r.id,r.record_json FROM review_rounds r JOIN issues t ON t.id=r.task_id WHERE t.parent_id=? ORDER BY r.rowid`, epic)
			if err != nil {
				return nil, err
			}
			for rows.Next() {
				var id, data string
				if err := rows.Scan(&id, &data); err != nil {
					_ = rows.Close()
					return nil, err
				}
				var r struct {
					Profiles []string `json:"profiles"`
				}
				if err := json.Unmarshal([]byte(data), &r); err != nil {
					_ = rows.Close()
					return nil, err
				}
				if slices.Contains(r.Profiles, profile) {
					view.Rounds = append(view.Rounds, id)
				}
			}
			if err := rows.Err(); err != nil {
				_ = rows.Close()
				return nil, err
			}
			_ = rows.Close()
		}
		views = append(views, view)
	}
	return views, nil
}

// BlockerExceptions affects acceptance only; subtasks never enter scheduling.
func BlockerExceptions(q db.Queryer) (map[string]bool, error) {
	rows, err := q.Query(`SELECT DISTINCT b.id,CASE WHEN b.kind='subtask' THEN p.parent_id ELSE b.parent_id END
 FROM issues b LEFT JOIN issues p ON p.id=b.parent_id
 JOIN links l ON ((l.from_id=b.id AND l.relation='blocks') OR (l.to_id=b.id AND l.relation='blocked_by')) AND l.from_type='issue' AND l.to_type='issue'
 JOIN issues t ON t.id=CASE WHEN l.relation='blocks' THEN l.to_id ELSE l.from_id END AND t.kind='task' AND t.status NOT IN ('done','cancelled')
 WHERE b.kind IN ('task','subtask') AND t.parent_id=CASE WHEN b.kind='subtask' THEN p.parent_id ELSE b.parent_id END`)
	if err != nil {
		return nil, err
	}
	type candidate struct{ id, epic string }
	candidates := []candidate{}
	for rows.Next() {
		var c candidate
		if err := rows.Scan(&c.id, &c.epic); err != nil {
			_ = rows.Close()
			return nil, err
		}
		candidates = append(candidates, c)
	}
	err = rows.Err()
	_ = rows.Close()
	if err != nil {
		return nil, err
	}
	result := map[string]bool{}
	for _, c := range candidates {
		a, err := Current(q, c.epic)
		if err != nil {
			return nil, err
		}
		result[c.id] = a.State == "current"
	}
	return result, nil
}

func Missing(q Queryer, epic string, roots ...string) ([]string, error) {
	a, err := Current(q, epic)
	if err != nil {
		return nil, err
	}
	missing := []string{}
	if a.State != "current" {
		missing = append(missing, "approval: "+a.State)
	}
	views, err := Views(q, epic, roots...)
	if err != nil {
		return nil, err
	}
	for _, v := range views {
		if v.State != "current" {
			missing = append(missing, fmt.Sprintf("%s %s: %s", v.ID, v.Verifier, v.State))
		}
	}
	return missing, nil
}
