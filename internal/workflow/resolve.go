package workflow

import (
	"database/sql"
	"fmt"
)

// GoverningHumanOutcome is the latest durable judgment for one class and
// Issue scope. Event identity, rather than a caller-authored timestamp, fixes
// the order.
type GoverningHumanOutcome struct {
	ID      string
	Result  string
	EventID string
}

// ResolveGoverningHumanOutcome selects by immutable Event identity so caller
// timestamps cannot change which judgment governs.
func ResolveGoverningHumanOutcome(conn *sql.DB, class, scopeIssueID string) (GoverningHumanOutcome, error) {
	var outcome GoverningHumanOutcome
	err := conn.QueryRow(`
		SELECT id,result,created_event_id
		FROM human_outcomes
		WHERE class=? AND scope_issue_id=?
		ORDER BY created_event_id DESC
		LIMIT 1`, class, scopeIssueID).Scan(&outcome.ID, &outcome.Result, &outcome.EventID)
	if err != nil {
		return GoverningHumanOutcome{}, fmt.Errorf("resolve governing Human Outcome: %w", err)
	}
	return outcome, nil
}

// HumanOutcomeSuperseded is the single supersession predicate (contracts.md
// section 7): an Outcome is superseded exactly when the governing Outcome for
// its class and scope is a different Outcome with a different result. Both the
// Obligation evaluator and the Stop resolution predicates call it.
func HumanOutcomeSuperseded(governing GoverningHumanOutcome, outcomeID, result string) bool {
	return governing.ID != outcomeID && governing.Result != result
}
