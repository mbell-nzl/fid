package journal

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/epicplan"
)

func epicValidationEvent(action string) Event {
	e := admissionEvent(stateMoveEventID, "epic.validation."+action, "epic", "FID-e", 1,
		Effect{Op: "insert", Table: "epic_validation_changes", PK: map[string]any{"epic_id": "FID-e", "generation": 1}, Values: map[string]any{"ordinal": 1, "claim": "Inspect fixture", "verifier": "human", "recipe": "docs/fixture.md", "removed": false}})
	e.CommandName = e.EventType
	if action != "add" {
		e.EventID = "EVT-20261110-000007"
		e.Aggregate.Sequence = 2
		e.Effects[0].PK["generation"] = 2
	}
	if action == "remove" {
		e.Effects[0].Values["removed"] = true
	}
	return e
}

func TestEpicValidationEvidenceAdmission(t *testing.T) {
	t.Parallel()
	for _, probe := range []string{"person valid", "person invalid commit", "check valid", "check nonzero pass", "check unrelated artifact"} {
		t.Run(probe, func(t *testing.T) {
			t.Parallel()
			change := epicValidationEvent("add")
			person := strings.HasPrefix(probe, "person")
			if !person {
				change.Effects[0].Values["verifier"] = "check:proof"
				change.Effects[0].Values["recipe"] = ""
			}
			prefix := append(stateMovePrefix(), change)
			p := epicplan.Pass{ID: "EVR-a1b2c3", EpicID: "FID-e", ValidationID: "V1", Generation: 1, Commit: strings.Repeat("a", 40), Merges: []string{}, Actor: "human:mark", At: admissionStamp(60), Outcome: "passed"}
			family := "epic.validation.verify"
			effects := []Effect{}
			if person {
				family = "epic.validation.pass"
				p.DecisionID = "PDN-pass"
				if probe == "person invalid commit" {
					p.Commit = strings.Repeat("z", 40)
				}
				binding := p
				binding.ID = ""
				effects = append(effects, Effect{Op: "insert", Table: "person_decisions", PK: map[string]any{"id": p.DecisionID}, Values: map[string]any{"kind": "pass-validation", "actor": p.Actor, "subject_type": "epic", "subject_id": p.EpicID, "reason": "Inspected fixture", "binding_id": p.ValidationID, "binding_json": epicChangeJSON(binding), "decided_at": p.At}})
			} else {
				code := 0
				if probe == "check nonzero pass" {
					code = 1
				}
				p.Definition = &checkrecord.Definition{Run: []string{"true"}, Timeout: time.Second, Environment: checkrecord.EnvironmentNames(), Artifacts: []string{}}
				p.Check = &checkrecord.Result{Outcome: "passed", ExitCode: &code, Counts: true, Tree: checkrecord.Tree{Head: p.Commit, Tree: strings.Repeat("b", 40)}, FinishedAt: p.At, Reasons: []string{}, Artifacts: []checkrecord.Artifact{}, RefusedArtifacts: []checkrecord.RefusedArtifact{}, OutputDigest: checkrecord.BytesDigest(nil)}
				if probe == "check unrelated artifact" {
					effects = append(effects, Effect{Op: "insert", Table: "referenced_artifacts", PK: map[string]any{"id": "ART-a1b2c3"}, Values: map[string]any{"digest": "sha256:" + strings.Repeat("c", 64), "media_type": "text/plain", "size_bytes": 1, "retention_class": "evidence", "availability": "available", "location_policy": "content-addressed", "created_at": p.At, "expires_at": nil}})
				}
			}
			effects = append(effects, Effect{Op: "insert", Table: "epic_validation_results", PK: map[string]any{"id": p.ID}, Values: map[string]any{"epic_id": p.EpicID, "ordinal": 1, "generation": 1, "record_json": epicChangeJSON(p)}})
			e := admissionEvent("EVT-20261110-000007", family, "epic", p.EpicID, 2, effects...)
			e.CommandName = family
			if !person {
				request := IdempotencyRequest{OperationScope: family + ":project-1", IdempotencyKey: "verify", RequestDigest: "sha256:" + strings.Repeat("d", 64)}
				e.Idempotency = &EventIdempotency{Key: request.IdempotencyKey, Scope: request.OperationScope, RequestDigest: request.RequestDigest}
				e.Effects = append(e.Effects, idempotencyResultEffect(request, e, p.EpicID, epicChangeJSON(epicplan.Result{EpicID: p.EpicID, At: p.At, Checks: []epicplan.Pass{p}})))
			}
			w := newAdmissionWorkspace(t, prefix, "commit")
			err := w.admit(t, "commit", e)
			want := strings.HasSuffix(probe, " valid")
			if (err == nil) != want {
				t.Fatalf("%s admission=%v", probe, err)
			}
		})
	}
}

func epicApprovalEvent() Event {
	binding := epicplan.ListBinding([]epicplan.Validation{}, 0)
	raw, _ := json.Marshal(binding)
	e := admissionEvent(stateMoveEventID, "epic.approve", "epic", "FID-e", 1, Effect{Op: "insert", Table: "person_decisions", PK: map[string]any{"id": "PDN-epic"}, Values: map[string]any{"kind": "approve-epic", "actor": "human:mark", "subject_type": "epic", "subject_id": "FID-e", "reason": "Approve the plan", "binding_id": binding.Digest, "binding_json": string(raw), "decided_at": admissionStamp(60)}})
	e.CommandName = e.EventType
	return e
}

func TestEpicPlansAcrossJournalPaths(t *testing.T) {
	t.Parallel()
	for _, action := range []string{"add", "update", "remove", "approve"} {
		t.Run(action, func(t *testing.T) {
			t.Parallel()
			for _, probe := range []string{"valid", "actor", "subject", "identity", "extra"} {
				t.Run(probe, func(t *testing.T) {
					t.Parallel()
					for _, mode := range admissionModes(epicApprovalEvent()) {
						t.Run(mode, func(t *testing.T) {
							prefix := stateMovePrefix()
							e := epicValidationEvent(action)
							if action == "update" || action == "remove" {
								prefix = append(prefix, epicValidationEvent("add"))
							}
							if action == "approve" {
								e = epicApprovalEvent()
							}
							switch probe {
							case "actor":
								e.ActorV2.ActorType = "service"
							case "subject":
								e.Aggregate.ID = "FID-a"
							case "identity":
								if action == "approve" {
									e.Effects[0].Values["binding_id"] = "wrong"
								} else {
									e.Effects[0].PK["generation"] = 9
								}
							case "extra":
								e.Effects = append(e.Effects, e.Effects[0])
							}
							w := newAdmissionWorkspace(t, prefix, mode)
							err := w.admit(t, mode, e)
							if (err == nil) != (probe == "valid") {
								t.Fatalf("%s/%s/%s: %v", action, probe, mode, err)
							}
						})
					}
				})
			}
		})
	}
}

func TestValidationIDsNeverReusedInAdmission(t *testing.T) {
	t.Parallel()
	prefix := append(stateMovePrefix(), epicValidationEvent("add"), epicValidationEvent("remove"))
	e := epicValidationEvent("add")
	e.EventID = "EVT-20261110-000008"
	e.Aggregate.Sequence = 3
	e.Effects[0].PK["generation"] = 3
	for _, ordinal := range []int{1, 2} {
		t.Run(fmt.Sprint(ordinal), func(t *testing.T) {
			t.Parallel()
			w := newAdmissionWorkspace(t, prefix, "commit")
			event := e
			event.Effects = append([]Effect{}, e.Effects...)
			event.Effects[0].Values = map[string]any{"ordinal": ordinal, "claim": "New claim", "verifier": "review:code", "recipe": "", "removed": false}
			err := w.admit(t, "commit", event)
			if (err == nil) != (ordinal == 2) {
				t.Fatalf("ordinal %d: %v", ordinal, err)
			}
		})
	}
}
