package export

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
)

// These are the mutable columns in the accepted operational Event families.
// Physical presence in a table does not authorize rewriting immutable facts.
var semanticUpdateColumns = map[string]string{
	"kit_activations":       "status point_of_no_return updated_at version",
	"activation_components": "status",
	"activation_resources":  "version observed_identity state evidence_ref error_code",
	"activation_recoveries": "completed_event_id",
	"work_envelopes":        "status",
	"envelope_revisions":    "status approved_actor_id activated_at completed_at",
	"claims":                "version expires_at status released_at release_reason",
	"runs":                  "status terminal_repository_json ended_at",
	"run_runtime":           "version",
	"attention_items":       "status version resolution_event_id",
	"referenced_artifacts":  "availability",
	"idempotency_results":   "result_entity_id result_event_id aggregate_sequence result_json",
	"holds":                 "ended_at end_reason",
}

func validateSemanticUpdate(effect journal.Effect) error {
	if effect.Op != "update" {
		return nil
	}
	columns, limited := semanticUpdateColumns[effect.Table]
	if !limited {
		return nil
	}
	for column := range effect.Values {
		if !strings.Contains(" "+columns+" ", " "+column+" ") {
			return fmt.Errorf("%s.%s is immutable in this Event family", effect.Table, column)
		}
	}
	return nil
}
