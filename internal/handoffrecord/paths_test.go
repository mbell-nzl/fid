package handoffrecord

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

func pathRecord() Record {
	return Record{
		View:   View{ID: "HOF-1", PR: "https://github.com/owner/repo/pull/98", Head: strings.Repeat("a", 40), Base: strings.Repeat("b", 40), At: "2026-11-22T00:00:00Z", RequiredChecks: []string{}, ResultIDs: []string{}},
		TaskID: "FID-1", HoldID: "HLD-1", Actor: "agent:builder", Tree: strings.Repeat("c", 40), Lane: "assured", PolicySnapshot: 1, DefinitionDigests: map[string]string{},
	}
}

func TestHandoffPathsRoundTripAndReuse(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		paths []string
		bytes bool
	}{
		{"empty", []string{}, false},
		{"ordinary", []string{"docs/a.md", "main.go"}, false},
		{"replacement rune", []string{"notes\ufffd.md"}, false},
		{"invalid UTF-8", []string{"main.go", "notes\xff.md"}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			record := pathRecord()
			// Reusing a record must discard the previous comparison's bytes.
			record.SetPaths([]string{"previous\xff.md"})
			record.SetPaths(tc.paths)
			if (record.ChangedPathBytes != nil) != tc.bytes {
				t.Fatalf("exact-byte list=%q, want present=%t", record.ChangedPathBytes, tc.bytes)
			}
			raw, err := json.Marshal(record)
			if err != nil {
				t.Fatal(err)
			}
			var reloaded Record
			if err := json.Unmarshal(raw, &reloaded); err != nil {
				t.Fatal(err)
			}
			if err := reloaded.Validate(); err != nil {
				t.Fatalf("valid comparison refused: %v", err)
			}
			if !slices.Equal(reloaded.Paths(), tc.paths) {
				t.Fatalf("comparison after JSON=%q, want %q", reloaded.Paths(), tc.paths)
			}
		})
	}
}

func TestHandoffPathsRejectInconsistentByteLists(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		bytes [][]byte
	}{
		{"empty bytes", [][]byte{}},
		{"missing path", [][]byte{[]byte("main.go")}},
		{"extra path", [][]byte{[]byte("main.go"), []byte("notes\xff.md"), []byte("extra.go")}},
		{"different path", [][]byte{[]byte("main.go"), []byte("other.md")}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			record := pathRecord()
			record.ChangedPaths = []string{"main.go", "notes\ufffd.md"}
			record.ChangedPathBytes = tc.bytes
			if err := record.Validate(); err == nil || !strings.Contains(err.Error(), "path bytes differ from paths") {
				t.Fatalf("inconsistent byte list accepted: %v", err)
			}
		})
	}
}
