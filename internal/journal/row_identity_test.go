package journal

import (
	"bytes"
	"errors"
	"os"
	"reflect"
	"testing"
)

type rowProbe struct {
	Name   string  `json:"name"`
	Prefix []Event `json:"prefix"`
	Event  Event   `json:"event"`
	Table  string
	Reason string
}

func TestRowIdentityOnEveryJournalPath(t *testing.T) {
	t.Parallel()
	corpus := readRowCorpus(t)
	representatives := map[string]bool{}
	for _, probe := range rowIdentityRepresentatives(corpus.Malformed) {
		representatives[probe.Name] = true
	}
	for _, group := range []struct {
		name    string
		probes  []rowProbe
		refused bool
	}{{"control", corpus.Controls, false}, {"malformed", corpus.Malformed, true}} {
		for _, probe := range group.probes {
			t.Run(group.name+"/"+probe.Name, func(t *testing.T) {
				t.Parallel()
				prefix := append(append([]Event{}, corpus.Prefix...), probe.Prefix...)
				modes := admissionModes(probe.Event)
				if group.refused && !representatives[probe.Name] {
					// Every malformed probe exercises the complete admission against a
					// real projection. Other paths need one per schema/table/reason;
					// controls keep every shape on every applicable path. Event-v1
					// cannot be committed, so its full boundary is projection.
					modes = []string{"project"}
					if probe.Event.IsV2() {
						modes = []string{"commit"}
					}
				}
				for _, mode := range modes {
					t.Run(mode, func(t *testing.T) {
						w := newAdmissionWorkspace(t, prefix, mode)
						reference := w
						if mode == "replay" && group.refused {
							reference = newAdmissionWorkspace(t, prefix, "project")
						}
						before := stateProjection(t, reference.root)
						err := w.admit(t, mode, probe.Event)
						if !group.refused {
							if err != nil {
								t.Fatalf("control refused: %v", err)
							}
							return
						}
						// Typed admission refusals survive the replay and commit wrappers.
						var refusal *RowIdentityError
						if !errors.As(err, &refusal) {
							t.Fatalf("malformed Event = %v (%T), want typed row refusal", err, err)
						}
						if refusal.Table != probe.Table || refusal.Reason != probe.Reason {
							t.Errorf("refusal = (%q, %q), want (%q, %q)", refusal.Table, refusal.Reason, probe.Table, probe.Reason)
						}
						if !reflect.DeepEqual(before, stateProjection(t, w.root)) {
							t.Fatal("refusal changed projection")
						}
						stored, err := os.ReadFile(w.file)
						if err != nil {
							t.Fatal(err)
						}
						live := mode == "commit" || mode == "dry-run" || mode == "project"
						want := w.journalWith(t, probe.Event)
						if live {
							want = w.prefix
						}
						if !bytes.Equal(stored, want) {
							t.Fatal("refusal changed journal")
						}
						if !live {
							if err := os.WriteFile(w.file, w.prefix, 0o644); err != nil {
								t.Fatal(err)
							}
						}
						valid := admissionLaneRaise()
						valid.EventID = probe.Event.EventID
						valid.Effects[1].Values["revision"] = 2
						if len(probe.Prefix) > 0 {
							valid.Effects[1].Values["revision"] = 3
							valid.Effects[1].Preconditions["revision"] = 2
							valid.Effects[1].Prior["revision"] = 2
							for _, prior := range probe.Prefix {
								if prior.Aggregate.Type == "issue" && prior.Aggregate.ID == "FID-a" {
									valid.Aggregate.Sequence = prior.Aggregate.Sequence + 1
								}
							}
						}
						valid.Effects[2].Values["journal_event_id"] = valid.EventID
						if err := w.admit(t, "commit", valid); err != nil {
							t.Fatalf("valid write after refusal: %v", err)
						}
					})
				}
			})
		}
	}
}

func TestMalformedCancellationHoldEnd(t *testing.T) {
	t.Parallel()
	hold := admissionHoldTaken()
	hold.EventID = "EVT-20261110-000006"
	prefix := append(stateMovePrefix(), hold)
	valid := stateCancel("FID-a", "task")
	valid.EventID = "EVT-20261110-000007"
	valid.Effects[3].Values["journal_event_id"] = valid.EventID
	ended := Effect{Op: "update", Table: "holds", PK: map[string]any{"id": "HLD-a"}, Values: map[string]any{"ended_at": valid.OccurredAt.Format("2006-01-02T15:04:05Z07:00"), "end_reason": "task-cancelled"}}
	valid.Effects = append(valid.Effects[:2:2], append([]Effect{ended}, valid.Effects[2:]...)...)
	for _, variant := range []string{"task-completed", "two holds"} {
		for _, mode := range admissionModes(valid) {
			t.Run(variant+"/"+mode, func(t *testing.T) {
				w := newAdmissionWorkspace(t, prefix, mode)
				malformed := valid
				raw, err := MarshalCanonicalJSON(valid)
				if err != nil {
					t.Fatal(err)
				}
				if err := malformed.UnmarshalJSON(raw); err != nil {
					t.Fatal(err)
				}
				if variant == "task-completed" {
					malformed.Effects[2].Values["end_reason"] = "task-completed"
				} else {
					malformed.Effects = append(malformed.Effects[:3:3], append([]Effect{malformed.Effects[2]}, malformed.Effects[3:]...)...)
				}
				reference := w
				if mode == "replay" && variant != "task-completed" {
					reference = newAdmissionWorkspace(t, prefix, "project")
				}
				before := stateProjection(t, reference.root)
				err = w.admit(t, mode, malformed)
				var refused *RowIdentityError
				reason := "invalid task hold end"
				if variant == "two holds" {
					reason = "stage may write at most 1 rows"
				}
				if !errors.As(err, &refused) || refused.Table != "holds" || refused.Reason != reason {
					t.Fatalf("malformed hold end = %v", err)
				}
				if !reflect.DeepEqual(before, stateProjection(t, w.root)) {
					t.Fatal("refusal changed projection")
				}
				stored, readErr := os.ReadFile(w.file)
				if readErr != nil {
					t.Fatal(readErr)
				}
				want := w.journalWith(t, malformed)
				if mode == "commit" || mode == "dry-run" || mode == "project" {
					want = w.prefix
				}
				if !bytes.Equal(stored, want) {
					t.Fatal("refusal changed journal")
				}
				if err := os.WriteFile(w.file, w.prefix, 0o644); err != nil {
					t.Fatal(err)
				}
				if mode == "replay" && variant == "task-completed" {
					if _, err := EnsureProjectionFresh(w.root, testProjectMigrationsDir(t)); err != nil {
						t.Fatal(err)
					}
				}
				if err := w.admit(t, "commit", valid); err != nil {
					t.Fatalf("valid cancellation after refusal: %v", err)
				}
			})
		}
	}
}
