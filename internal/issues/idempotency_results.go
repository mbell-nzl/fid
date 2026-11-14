package issues

import (
	"encoding/json"
	"fmt"

	"github.com/mbell-nzl/fid/internal/journal"
)

func createResultJSON(result CreateResult) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{
		"event_occurred_at": result.OccurredAt,
		"idempotency_key":   result.IdempotencyKey,
		"issue":             map[string]any{"lane": result.Lane, "id": result.ID, "kind": result.Kind, "revision": result.Revision, "status": result.Status, "title": result.Title},
		"replayed":          false,
	})
	return string(payload), err
}

func createResultFromJSON(raw string, eventID string, replayed bool) (CreateResult, error) {
	var value struct {
		EventOccurredAt string `json:"event_occurred_at"`
		IdempotencyKey  string `json:"idempotency_key"`
		Issue           struct {
			Lane     *string `json:"lane"`
			ID       string  `json:"id"`
			Kind     string  `json:"kind"`
			Revision int     `json:"revision"`
			Status   string  `json:"status"`
			Title    string  `json:"title"`
		} `json:"issue"`
	}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return CreateResult{}, fmt.Errorf("decode stored Issue create result: %w", err)
	}
	return CreateResult{Lane: value.Issue.Lane, ID: value.Issue.ID, Kind: value.Issue.Kind, Title: value.Issue.Title, Status: value.Issue.Status, Revision: value.Issue.Revision, EventID: eventID, OccurredAt: value.EventOccurredAt, IdempotencyKey: value.IdempotencyKey, Replayed: replayed}, nil
}

func updateResultJSON(result UpdateResult) (string, error) {
	return resultJSON(map[string]any{
		"event_occurred_at": result.OccurredAt,
		"idempotency_key":   result.IdempotencyKey,
		"issue":             map[string]any{"id": result.ID, "revision": result.Revision},
		"changes":           result.Changes,
		"replayed":          false,
	})
}

func updateResultFromJSON(raw, eventID string, replayed bool) (UpdateResult, error) {
	if err := journal.ValidateIssueUpdateResult(raw); err != nil {
		return UpdateResult{}, err
	}
	var value struct {
		EventOccurredAt string        `json:"event_occurred_at"`
		IdempotencyKey  string        `json:"idempotency_key"`
		Changes         []FieldChange `json:"changes"`
		Issue           struct {
			Lane     *string `json:"lane"`
			ID       string  `json:"id"`
			Revision int     `json:"revision"`
		} `json:"issue"`
	}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return UpdateResult{}, fmt.Errorf("decode stored Issue update result: %w", err)
	}
	return UpdateResult{ID: value.Issue.ID, Revision: value.Issue.Revision, Changes: value.Changes, EventID: eventID, OccurredAt: value.EventOccurredAt, IdempotencyKey: value.IdempotencyKey, Replayed: replayed}, nil
}

func transitionResultJSON(result TransitionResult) (string, error) {
	return resultJSON(map[string]any{
		"event_occurred_at": result.OccurredAt,
		"idempotency_key":   result.IdempotencyKey,
		"issue":             map[string]any{"id": result.ID, "revision": result.Revision},
		"from_state":        result.FromState,
		"decision_id":       result.DecisionID,
		"to_state":          result.ToState,
		"replayed":          false,
	})
}

func transitionResultFromJSON(raw, eventID string, replayed bool) (TransitionResult, error) {
	var value struct {
		EventOccurredAt string  `json:"event_occurred_at"`
		IdempotencyKey  string  `json:"idempotency_key"`
		DecisionID      *string `json:"decision_id"`
		FromState       string  `json:"from_state"`
		ToState         string  `json:"to_state"`
		Issue           struct {
			Lane     *string `json:"lane"`
			ID       string  `json:"id"`
			Revision int     `json:"revision"`
		} `json:"issue"`
	}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return TransitionResult{}, fmt.Errorf("decode stored Issue transition result: %w", err)
	}
	return TransitionResult{DecisionID: value.DecisionID, ID: value.Issue.ID, Revision: value.Issue.Revision, FromState: value.FromState, ToState: value.ToState, EventID: eventID, OccurredAt: value.EventOccurredAt, IdempotencyKey: value.IdempotencyKey, Replayed: replayed}, nil
}

func acAddResultJSON(result AddAcceptanceCriterionResult) (string, error) {
	return resultJSON(map[string]any{
		"acceptance_criterion": map[string]any{"id": result.ID, "issue_id": result.IssueID, "ordinal": result.Ordinal, "text": result.Text},
		"event_occurred_at":    result.OccurredAt,
		"idempotency_key":      result.IdempotencyKey,
		"issue_revision":       result.IssueRevision,
		"replayed":             false,
	})
}

func acAddResultFromJSON(raw, eventID string, replayed bool) (AddAcceptanceCriterionResult, error) {
	var value struct {
		Criterion struct {
			ID      string `json:"id"`
			IssueID string `json:"issue_id"`
			Ordinal int    `json:"ordinal"`
			Text    string `json:"text"`
		} `json:"acceptance_criterion"`
		IdempotencyKey string `json:"idempotency_key"`
		OccurredAt     string `json:"event_occurred_at"`
		IssueRevision  int    `json:"issue_revision"`
	}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return AddAcceptanceCriterionResult{}, fmt.Errorf("decode stored AC add result: %w", err)
	}
	return AddAcceptanceCriterionResult{ID: value.Criterion.ID, IssueID: value.Criterion.IssueID, Ordinal: value.Criterion.Ordinal, Text: value.Criterion.Text, IssueRevision: value.IssueRevision, EventID: eventID, OccurredAt: value.OccurredAt, IdempotencyKey: value.IdempotencyKey, Replayed: replayed}, nil
}

func resultJSON(value map[string]any) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(value)
	return string(payload), err
}
