package journal

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"
)

type rowProbe struct {
	Name   string  `json:"name"`
	Prefix []Event `json:"prefix"`
	Event  Event   `json:"event"`
}

func TestRowIdentityOnEveryJournalPath(t *testing.T) {
	t.Parallel()
	raw, err := os.ReadFile("../../testdata/journal/row-identity.json")
	if err != nil {
		t.Fatal(err)
	}
	var corpus struct {
		Prefix              []Event
		Controls, Malformed []rowProbe
	}
	if err := json.Unmarshal(raw, &corpus); err != nil {
		t.Fatal(err)
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
				for _, mode := range admissionModes(probe.Event) {
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
						var refusal interface{ RowAdmissionRefusal() }
						if !errors.As(err, &refusal) {
							t.Fatalf("malformed Event = %v (%T), want typed row refusal", err, err)
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
