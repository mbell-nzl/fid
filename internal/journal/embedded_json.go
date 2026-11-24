package journal

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mbell-nzl/fid/internal/approvalrecord"
	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/completionrecord"
	"github.com/mbell-nzl/fid/internal/defaultconflict"
	"github.com/mbell-nzl/fid/internal/epicrecord"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/mergerecord"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/waiverrecord"
)

// JSON text is an Event boundary, independent of the family or SQL operation.
// Scan even unused insert prior/preconditions: they are retained in history and
// archives. Names ending in _json are the schema's JSON-text convention, also
// used by archive validation. Unknown columns still face the column contracts.
func validateEmbeddedJSON(event Event) error {
	for _, effect := range event.Effects {
		for _, placement := range []struct {
			name   string
			values map[string]any
		}{
			{"pk", effect.PK}, {"values", effect.Values}, {"prior", effect.Prior}, {"preconditions", effect.Preconditions},
		} {
			for column, value := range placement.values {
				value = currentEventValue(value)
				if !strings.HasSuffix(column, "_json") || value == nil {
					continue
				}
				raw, ok := value.(string)
				if !ok {
					return rowFailure(event, effect.Table, fmt.Sprintf("%s.%s requires JSON text", placement.name, column))
				}
				if err := validateTypedJSONText(raw, embeddedJSONType(event, effect.Table, column)); err != nil {
					return rowFailure(event, effect.Table, fmt.Sprintf("%s.%s: %v", placement.name, column, err))
				}
			}
		}
	}
	return nil
}

// Destinations describe key matching, not payload validity. Semantic readers
// still check their complete records. Structs fold keys like encoding/json;
// map/interface destinations keep exact keys at every nested map boundary.
var embeddedJSONTypes = map[string]reflect.Type{
	"finding_responses.record_json":                  reflect.TypeFor[reviewrecord.Response](),
	"default_conflicts.record_json":                  reflect.TypeFor[defaultconflict.Record](),
	"review_claims.binding_json":                     reflect.TypeFor[fingerprint.Result](),
	"completions.record_json":                        reflect.TypeFor[completionrecord.Record](),
	"review_rounds.record_json":                      reflect.TypeFor[reviewrecord.Round](),
	"round_findings.record_json":                     reflect.TypeFor[reviewrecord.Finding](),
	"merge_observations.record_json":                 reflect.TypeFor[mergerecord.Merge](),
	"ci_results.record_json":                         reflect.TypeFor[mergerecord.CIResult](),
	"handoffs.record_json":                           reflect.TypeFor[handoffrecord.Record](),
	"check_requests.names_json":                      reflect.TypeFor[[]string](),
	"check_requests.response_json":                   reflect.TypeFor[checkrecord.Response](),
	"check_starts.start_json":                        reflect.TypeFor[checkrecord.Start](),
	"check_processes.process_json":                   reflect.TypeFor[checkrecord.Process](),
	"check_results.result_json":                      reflect.TypeFor[checkrecord.Result](),
	"policy_snapshots.effective_json":                reflect.TypeFor[policy.Effective](),
	"idempotency_results.request_field_digests_json": reflect.TypeFor[map[string]string](),
}

func embeddedJSONType(event Event, table, column string) reflect.Type {
	if table == "person_decisions" && column == "binding_json" && event.EventType == "task.waived" {
		return reflect.TypeFor[waiverrecord.Binding]()
	}
	if table == "person_decisions" && column == "binding_json" && event.EventType == "task.approved" {
		return reflect.TypeFor[approvalrecord.Binding]()
	}
	if table != "idempotency_results" || column != "result_json" {
		return embeddedJSONTypes[table+"."+column]
	}
	switch event.EventType {
	case "task.finished":
		return reflect.TypeFor[completionrecord.Result]()
	case "task.waived":
		return reflect.TypeFor[waiverrecord.Result]()
	case "task.approved":
		return reflect.TypeFor[approvalrecord.Result]()
	case "epic.closed":
		return reflect.TypeFor[epicrecord.Result]()
	case "finding.closed":
		return reflect.TypeFor[findingCloseResult]()
	case "review.round-recorded":
		return reflect.TypeFor[reviewRetryJSON]()
	case "check.completed", "check.stopped", "check.started", "check.result-recorded":
		return reflect.TypeFor[checkrecord.Response]()
	case "issue.created", "issue.transitioned":
		return reflect.TypeFor[issueObjectRetryJSON]()
	case "issue.updated":
		return reflect.TypeFor[issueUpdateRetryJSON]()
	case "issue.ac-added":
		return reflect.TypeFor[struct {
			Criterion struct{} `json:"acceptance_criterion"`
		}]()
	case "link.added", "link.removed":
		return reflect.TypeFor[struct {
			Link    struct{} `json:"link"`
			Inverse struct{} `json:"inverse"`
		}]()
	case "label.added":
		return reflect.TypeFor[struct {
			Label struct{} `json:"label"`
		}]()
	case "artifact.availability-observed":
		return reflect.TypeFor[struct {
			Artifact struct{} `json:"artifact"`
		}]()
	default:
		return reflect.TypeFor[struct {
			Fields map[string]any `json:"fields"`
		}]()
	}
}

type reviewRetryJSON struct {
	Round reviewrecord.Round `json:"round"`
	At    string             `json:"occurred_at"`
}

// These Issue objects contain scalar fields; unknown nested members remain
// untyped. Keep the shapes per family so unknown response maps stay exact.
type issueObjectRetryJSON struct {
	Issue struct{} `json:"issue"`
}
type issueUpdateRetryJSON struct {
	Issue   struct{} `json:"issue"`
	Changes []struct {
		Field    string
		From, To any
	} `json:"changes"`
}
