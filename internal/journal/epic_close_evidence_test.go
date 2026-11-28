package journal

import (
	"testing"

	"github.com/mbell-nzl/fid/internal/epicplan"
)

func TestEpicCloseEvidenceAdmission(t *testing.T) {
	t.Parallel()
	for _, probe := range []string{"valid", "approval", "extra result", "extra merge", "missing binding"} {
		t.Run(probe, func(t *testing.T) {
			t.Parallel()
			for _, mode := range admissionModes(epicApprovalEvent()) {
				t.Run(mode, func(t *testing.T) {
					t.Parallel()
					prefix := append(stateMovePrefix(), epicApprovalEvent())
					b := epicplan.CloseBinding{ApprovalID: "PDN-epic", Results: []string{}, Merges: []string{}}
					switch probe {
					case "approval":
						b.ApprovalID = "PDN-other"
					case "extra result":
						b.Results = []string{"EVR-other"}
					case "extra merge":
						b.Merges = []string{"unrecorded"}
					}
					decision := Effect{Op: "insert", Table: "person_decisions", PK: map[string]any{"id": "PDN-close"}, Values: map[string]any{"kind": "close-epic", "actor": "human:mark", "subject_type": "epic", "subject_id": "FID-e", "reason": "All evidence current", "binding_id": "FID-e", "binding_json": epicChangeJSON(b), "decided_at": admissionStamp(60)}}
					if probe == "missing binding" {
						delete(decision.Values, "binding_json")
					}
					move := Effect{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-e"}, Values: map[string]any{"status": "done", "revision": 2, "updated_at": admissionStamp(60)}, Preconditions: map[string]any{"status": "todo", "revision": 1}}
					e := admissionEvent("EVT-20261110-000007", "epic.validated-close", "epic", "FID-e", 2, decision, move)
					e.CommandName, e.ActorV2 = "task.done", statePerson
					w := newAdmissionWorkspace(t, prefix, mode)
					err := w.admit(t, mode, e)
					if (err == nil) != (probe == "valid") {
						t.Fatalf("%s/%s: %v", probe, mode, err)
					}
				})
			}
		})
	}
}
