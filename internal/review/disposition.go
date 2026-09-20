package review

import (
	"database/sql"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/mbell-nzl/fid/internal/journal"
)

func resolveFindingEffects(conn *sql.DB, request RecordRequest, binding recordBinding, actorID, reviewID string) ([]journal.Effect, error) {
	type disposition struct{ findingID, action, reason string }
	dispositions := make([]disposition, 0, len(request.ResolvesFindingIDs)+len(request.InapplicableFindings))
	for _, findingID := range request.ResolvesFindingIDs {
		dispositions = append(dispositions, disposition{findingID: findingID, action: "resolved", reason: "resolved by qualifying Review"})
	}
	for _, input := range request.InapplicableFindings {
		dispositions = append(dispositions, disposition{findingID: input.FindingID, action: "inapplicable", reason: *input.Reason})
	}
	created := make(map[string]bool, len(request.Findings))
	for _, finding := range request.Findings {
		created[finding.ID] = true
	}
	effects := make([]journal.Effect, 0, len(dispositions))
	for _, input := range dispositions {
		effect, err := terminalFindingEffect(conn, input.findingID, input.action, input.reason, binding, actorID, reviewID, created[input.findingID])
		if err != nil {
			return nil, err
		}
		effects = append(effects, effect)
	}
	return effects, nil
}

func terminalFindingEffect(conn *sql.DB, findingID, action, reason string, binding recordBinding, actorID, reviewID string, created bool) (journal.Effect, error) {
	if created {
		return journal.Effect{}, &FindingTargetInvalidError{ID: findingID}
	}
	var issueID, envelopeID, creatorReviewID string
	var revision int
	err := conn.QueryRow(`SELECT r.issue_id, r.envelope_id, r.envelope_revision, f.review_id FROM review_findings f JOIN review_records r ON r.id=f.review_id WHERE f.id=?`, findingID).Scan(&issueID, &envelopeID, &revision, &creatorReviewID)
	if err == sql.ErrNoRows {
		return journal.Effect{}, &FindingNotFoundError{ID: findingID}
	}
	if err != nil {
		return journal.Effect{}, fmt.Errorf("load terminal Finding %q: %w", findingID, err)
	}
	if issueID != binding.IssueID || envelopeID != binding.EnvelopeID || revision != binding.EnvelopeRevision || creatorReviewID == reviewID {
		return journal.Effect{}, &FindingTargetInvalidError{ID: findingID}
	}
	var state string
	if err := conn.QueryRow(`SELECT action FROM finding_dispositions WHERE finding_id=? AND action IN ('resolved','inapplicable')`, findingID).Scan(&state); err != sql.ErrNoRows {
		if err == nil {
			return journal.Effect{}, &FindingAlreadyDisposedError{ID: findingID, State: state}
		}
		return journal.Effect{}, fmt.Errorf("review resolution Finding %q: %w", findingID, err)
	}
	dispositionID, err := nextDispositionID(conn)
	if err != nil {
		return journal.Effect{}, err
	}
	return journal.Effect{Op: "insert", Table: "finding_dispositions", PK: map[string]any{"id": dispositionID}, Values: map[string]any{"finding_id": findingID, "action": action, "actor_id": actorID, "claim_id": nil, "qualifying_review_id": reviewID, "reason": reason, "disposed_at": "<occurred-at>", "created_event_id": "<event-id>"}}, nil
}

func validateTerminalFindingInputs(resolved []string, inapplicable []InapplicableFindingInput) error {
	if len(resolved)+len(inapplicable) > 64 {
		return fmt.Errorf("review terminal Finding count exceeds 64")
	}
	seen := make(map[string]bool, len(resolved)+len(inapplicable))
	for _, id := range resolved {
		if strings.TrimSpace(id) == "" || seen[id] {
			return fmt.Errorf("review terminal Finding ID is invalid or duplicated")
		}
		seen[id] = true
	}
	for _, input := range inapplicable {
		if strings.TrimSpace(input.FindingID) == "" || seen[input.FindingID] || input.Reason == nil || strings.TrimSpace(*input.Reason) == "" || len(*input.Reason) > 4096 || !utf8.ValidString(*input.Reason) {
			return fmt.Errorf("review inapplicable Finding input is invalid")
		}
		seen[input.FindingID] = true
	}
	return nil
}

func validateTerminalFindingVerdict(request RecordRequest) error {
	if len(request.ResolvesFindingIDs)+len(request.InapplicableFindings) > 0 && request.Verdict != "pass" && request.Verdict != "changes-required" {
		return &VerdictInvalidError{Verdict: request.Verdict, Reason: "resolving Findings requires pass or changes-required"}
	}
	return nil
}
