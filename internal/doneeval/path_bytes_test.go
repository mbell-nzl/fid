package doneeval

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
)

func persistPathHandoff(t *testing.T, s *Snapshot, path string) {
	t.Helper()
	s.Handoff.SetPaths([]string{path})
	raw, err := json.Marshal(s.Handoff)
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "handoff.json")
	if err := os.WriteFile(file, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	raw, err = os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	var reloaded handoffrecord.Record
	if err := json.Unmarshal(raw, &reloaded); err != nil {
		t.Fatal(err)
	}
	if reloaded.ChangedPaths[0] != "notes\ufffd.md" || reloaded.Paths()[0] != path {
		t.Fatalf("persisted paths=%q, exact paths=%q", reloaded.ChangedPaths, reloaded.Paths())
	}
	s.Handoff = &reloaded
}

func TestDoneConditionalChecksSurviveInvalidUTF8Handoff(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, path string
		required   bool
	}{
		{"invalid UTF-8", "notes\xff.md", true},
		{"readable replacement", "notes\ufffd.md", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := completeSnapshot("assured", true)
			check := s.Policy.Checks["test"]
			check.WhenChanged = []string{"**/*.go"}
			s.Policy.Checks["test"] = check
			digest := checkrecord.Digest(checkrecord.DefinitionFor(check))
			s.Handoff.DefinitionDigests["test"] = digest
			s.Checks[0].Start.DefinitionDigest = digest
			if !tc.required {
				s.Handoff.RequiredChecks, s.Handoff.ResultIDs = []string{}, []string{}
				s.Handoff.DefinitionDigests = map[string]string{}
				s.Checks = nil
			}
			persistPathHandoff(t, &s, tc.path)
			if got := Evaluate(s); got.Status != Complete {
				t.Fatalf("serialized hand-off refused: %+v", got)
			}
			if tc.required {
				s.Checks = nil
				if got := Evaluate(s); got.Status != Refused || got.Gates[ChecksGate].Reason != ChecksUnsatisfied {
					t.Fatalf("missing conditional result accepted: %+v", got)
				}
			}
		})
	}
}

func TestDoneFloorSurvivesInvalidUTF8Handoff(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, path string
		want       Status
	}{
		{"invalid UTF-8", "notes\xff.md", Refused},
		{"readable replacement", "notes\ufffd.md", Complete},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := completeSnapshot("standard", true)
			persistPathHandoff(t, &s, tc.path)
			got := Evaluate(s)
			if got.Status != tc.want || tc.want == Refused && got.Gates[PolicyGate].Reason != PolicyMismatch {
				t.Fatalf("floor after persistence: %+v, want %s", got, tc.want)
			}
		})
	}
}
