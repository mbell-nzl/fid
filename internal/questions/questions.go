// Package questions reads the durable questions a task has asked a person.
package questions

import (
	"fmt"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

type Answer struct {
	QuestionID string `json:"question_id"`
	Question   string `json:"question"`
	Text       string `json:"text"`
	DecisionID string `json:"decision_id"`
}

type Question struct {
	ID         string  `json:"id"`
	TaskID     string  `json:"task_id"`
	HoldID     string  `json:"hold_id"`
	Text       string  `json:"question"`
	Asker      string  `json:"asker"`
	AskedAt    string  `json:"asked_at"`
	ResolvedAt *string `json:"resolved_at"`
	Resolution *string `json:"resolution"`
	Answer     *Answer `json:"answer"`
}

// Read preserves journal insertion order, including questions asked in the
// same second. Filtering in SQL keeps show bounded to its task's history.
func Read(conn fiddb.Queryer, task string) ([]Question, error) {
	rows, err := conn.Query(`SELECT q.id,q.task_id,q.hold_id,q.text,q.asker,q.asked_at,q.resolved_at,q.resolution,d.id,d.reason
 FROM questions q LEFT JOIN person_decisions d ON d.id=q.decision_id
 WHERE (?='' OR q.task_id=?) ORDER BY q.rowid`, task, task)
	if err != nil {
		return nil, fmt.Errorf("read questions: %w", err)
	}
	defer func() { _ = rows.Close() }()
	result := []Question{}
	for rows.Next() {
		var q Question
		var decision, answer *string
		if err := rows.Scan(&q.ID, &q.TaskID, &q.HoldID, &q.Text, &q.Asker, &q.AskedAt, &q.ResolvedAt, &q.Resolution, &decision, &answer); err != nil {
			return nil, fmt.Errorf("scan question: %w", err)
		}
		if decision != nil {
			q.Answer = &Answer{q.ID, q.Text, *answer, *decision}
		}
		result = append(result, q)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate questions: %w", err)
	}
	return result, nil
}

func (q Question) Resolve(at, resolution string, decision any) journal.Effect {
	return journal.Effect{Op: "update", Table: "questions", PK: map[string]any{"id": q.ID}, Values: map[string]any{"resolved_at": at, "resolution": resolution, "decision_id": decision}, Preconditions: map[string]any{"resolved_at": nil}}
}
