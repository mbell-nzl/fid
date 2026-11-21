package export

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

// The same corpus reaches shared admission through import's private projection.
func TestImportRowIdentity(t *testing.T) {
	t.Parallel()
	raw, err := os.ReadFile("../../testdata/journal/row-identity.json")
	if err != nil {
		t.Fatal(err)
	}
	type probe struct {
		Name   string
		Prefix []journal.Event
		Event  journal.Event
	}
	var corpus struct {
		Prefix              []journal.Event
		Controls, Malformed []probe
	}
	if err := json.Unmarshal(raw, &corpus); err != nil {
		t.Fatal(err)
	}
	for _, group := range []struct {
		name    string
		probes  []probe
		refused bool
	}{{"control", corpus.Controls, false}, {"malformed", corpus.Malformed, true}} {
		for _, p := range group.probes {
			if !p.Event.IsV2() {
				continue
			} // Import already refuses all Event-v1 history.
			t.Run(group.name+"/"+p.Name, func(t *testing.T) {
				t.Parallel()
				v, err := newSemanticValidator(testMigrationsDir(t), 35)
				if err != nil {
					t.Fatal(err)
				}
				defer func() {
					if err := v.close(); err != nil {
						t.Error(err)
					}
				}()
				for _, e := range append(append([]journal.Event{}, corpus.Prefix...), p.Prefix...) {
					if err := applyLeanEmptyEvent(v, "import", e); err != nil {
						t.Fatalf("prefix: %v", err)
					}
				}
				before := handoffProjectionSnapshot(t, v)
				err = applyLeanEmptyEvent(v, "import", p.Event)
				if !group.refused {
					if err != nil {
						t.Fatalf("control: %v", err)
					}
					return
				}
				var semantic *semanticValidationError
				if !errors.As(err, &semantic) {
					t.Fatalf("malformed import = %v, want typed semantic failure", err)
				}
				if before != handoffProjectionSnapshot(t, v) {
					t.Fatal("refusal changed projection")
				}
				var valid journal.Event
				encoded, marshalErr := journal.MarshalCanonicalJSON(corpus.Controls[0].Event)
				if marshalErr != nil {
					t.Fatal(marshalErr)
				}
				if err := valid.UnmarshalJSON(encoded); err != nil {
					t.Fatal(err)
				}
				valid.EventID = p.Event.EventID
				valid.Effects[2].Values["journal_event_id"] = valid.EventID
				if len(p.Prefix) > 0 {
					valid.Effects[1].Values["revision"] = 3
					valid.Effects[1].Preconditions["revision"] = 2
					valid.Effects[1].Prior["revision"] = 2
				}
				if err := applyLeanEmptyEvent(v, "import", valid); err != nil {
					t.Fatalf("valid write after refusal: %v", err)
				}
			})
		}
	}
}
