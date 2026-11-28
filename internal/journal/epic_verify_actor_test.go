package journal

import (
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/epicplan"
)

func TestEpicVerificationWorkActorsAcrossJournalPaths(t *testing.T) {
	t.Parallel()
	for _, actor := range []string{"agent", "human", "service"} {
		t.Run(actor, func(t *testing.T) {
			t.Parallel()
			e := admissionEvent(stateMoveEventID, "epic.validation.verify", "epic", "FID-e", 1)
			e.CommandName = e.EventType
			e.ActorV2.ActorType = actor
			r := epicplan.Result{EpicID: "FID-e", At: admissionStamp(60), Checks: []epicplan.Pass{}}
			request := IdempotencyRequest{OperationScope: e.EventType + ":project-1", IdempotencyKey: "verify", RequestDigest: "sha256:" + strings.Repeat("d", 64)}
			e.Idempotency = &EventIdempotency{Key: request.IdempotencyKey, Scope: request.OperationScope, RequestDigest: request.RequestDigest}
			e.Effects = []Effect{idempotencyResultEffect(request, e, "FID-e", epicChangeJSON(r))}
			for _, mode := range admissionModes(e) {
				t.Run(mode, func(t *testing.T) {
					t.Parallel()
					w := newAdmissionWorkspace(t, stateMovePrefix(), mode)
					err := w.admit(t, mode, e)
					if (err == nil) != (actor != "service") {
						t.Fatalf("%s verification admission: %v", actor, err)
					}
				})
			}
		})
	}
}
