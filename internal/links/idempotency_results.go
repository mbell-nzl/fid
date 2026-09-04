package links

import (
	"encoding/json"
	"fmt"

	"github.com/mbell-nzl/fid/internal/journal"
)

func addResultJSON(result AddResult) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"event_occurred_at": result.OccurredAt, "idempotency_key": result.IdempotencyKey, "inverse": result.Inverse, "issue_revisions": result.IssueRevisions, "link": result.Declared, "replayed": false})
	return string(payload), err
}

func addResultFromJSON(raw, eventID string, replayed bool) (AddResult, error) {
	var value struct {
		OccurredAt     string          `json:"event_occurred_at"`
		IdempotencyKey string          `json:"idempotency_key"`
		IssueRevisions []IssueRevision `json:"issue_revisions"`
		Inverse        Link            `json:"inverse"`
		Link           Link            `json:"link"`
	}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return AddResult{}, fmt.Errorf("decode stored link add result: %w", err)
	}
	return AddResult{Declared: value.Link, Inverse: value.Inverse, IssueRevisions: value.IssueRevisions, EventID: eventID, OccurredAt: value.OccurredAt, IdempotencyKey: value.IdempotencyKey, Replayed: replayed}, nil
}

func removeResultJSON(result RemoveResult) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"event_occurred_at": result.OccurredAt, "idempotency_key": result.IdempotencyKey, "inverse": result.Inverse, "issue_revisions": result.IssueRevisions, "link": result.Declared, "replayed": false})
	return string(payload), err
}

func removeResultFromJSON(raw, eventID string, replayed bool) (RemoveResult, error) {
	var value struct {
		OccurredAt     string          `json:"event_occurred_at"`
		IdempotencyKey string          `json:"idempotency_key"`
		IssueRevisions []IssueRevision `json:"issue_revisions"`
		Inverse        Link            `json:"inverse"`
		Link           Link            `json:"link"`
	}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return RemoveResult{}, fmt.Errorf("decode stored link remove result: %w", err)
	}
	return RemoveResult{Declared: value.Link, Inverse: value.Inverse, IssueRevisions: value.IssueRevisions, EventID: eventID, OccurredAt: value.OccurredAt, IdempotencyKey: value.IdempotencyKey, Replayed: replayed}, nil
}
