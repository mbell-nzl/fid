// Package persondecision records decisions made by a person (ADR-0014).
package persondecision

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

type DeniedError struct{ Actor string }

func (err *DeniedError) Error() string { return "only a person can make this decision" }

type ReasonRequiredError struct{}

func (*ReasonRequiredError) Error() string { return "a reason is required" }

// Effect creates an immutable decision initially bound to its subject. An
// answer command sets its question binding before the atomic write.
func Effect(conn *sql.DB, kind, actor, subjectType, subjectID, reason string, at time.Time) (string, journal.Effect, error) {
	if !policy.Allows(actor, policy.PersonDecisions) {
		return "", journal.Effect{}, &DeniedError{Actor: actor}
	}
	if strings.TrimSpace(reason) == "" {
		return "", journal.Effect{}, &ReasonRequiredError{}
	}
	id, err := journal.GenerateID("PDN", func(id string) (bool, error) {
		var count int
		err := conn.QueryRow(`SELECT count(*) FROM person_decisions WHERE id=?`, id).Scan(&count)
		return count > 0, err
	})
	if err != nil {
		return "", journal.Effect{}, fmt.Errorf("new person's decision: %w", err)
	}
	return id, journal.Effect{Op: "insert", Table: "person_decisions", PK: map[string]any{"id": id}, Values: map[string]any{
		"kind": kind, "actor": actor, "subject_type": subjectType, "subject_id": subjectID, "reason": reason, "binding_id": subjectID, "decided_at": holds.Timestamp(at),
	}}, nil
}
