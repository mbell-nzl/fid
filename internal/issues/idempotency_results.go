package issues

import (
	"encoding/json"
	"fmt"

	"github.com/mbell-nzl/fid/internal/journal"
)

func createResultJSON(result CreateResult) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{
		"idempotency_key": result.IdempotencyKey,
		"issue":           map[string]any{"id": result.ID, "kind": result.Kind, "revision": result.Revision, "status": result.Status, "title": result.Title},
		"replayed":        false,
	})
	return string(payload), err
}

func createResultFromJSON(raw string, eventID string, replayed bool) (CreateResult, error) {
	var value struct {
		IdempotencyKey string `json:"idempotency_key"`
		Issue          struct {
			ID       string `json:"id"`
			Kind     string `json:"kind"`
			Revision int    `json:"revision"`
			Status   string `json:"status"`
			Title    string `json:"title"`
		} `json:"issue"`
	}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return CreateResult{}, fmt.Errorf("decode stored Issue create result: %w", err)
	}
	return CreateResult{ID: value.Issue.ID, Kind: value.Issue.Kind, Title: value.Issue.Title, Status: value.Issue.Status, Revision: value.Issue.Revision, EventID: eventID, IdempotencyKey: value.IdempotencyKey, Replayed: replayed}, nil
}
