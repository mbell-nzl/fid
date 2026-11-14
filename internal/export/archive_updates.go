package export

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
)

// These are the mutable columns in the accepted operational Event families.
// Physical presence in a table does not authorize rewriting immutable facts.
var semanticUpdateColumns = map[string]string{
	"check_requests":       "response_json",
	"referenced_artifacts": "availability",
	"idempotency_results":  "result_entity_id result_event_id aggregate_sequence result_json",
	"questions":            "resolved_at resolution decision_id",
	"holds":                "ended_at end_reason last_heartbeat_at expires_at kept park_reason",
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
