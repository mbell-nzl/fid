package waiverrecord

import (
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/fingerprint"
)

func TestWaiverRecordValidity(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		edit  func(*Record)
		valid bool
	}{
		{"complete", func(*Record) {}, true},
		{"agent", func(r *Record) { r.Actor = "agent:mark" }, false},
		{"reason", func(r *Record) { r.Reason = "\t" }, false},
		{"binding", func(r *Record) { r.Binding.Fingerprint = "" }, false},
		{"algorithm", func(r *Record) { r.Binding.Algorithm = "other" }, false},
		{"lane", func(r *Record) { r.Binding.Lane = "assured" }, false},
		{"time", func(r *Record) { r.At = "yesterday" }, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := Record{ID: "PDN-a", TaskID: "FID-a", Actor: "human:mark", Reason: "Reviewed risk", At: "2026-11-23T00:00:00Z", Binding: Binding{Lane: "standard", Algorithm: fingerprint.Algorithm, Fingerprint: strings.Repeat("f", 64), HandoffID: "HOF-a", Base: strings.Repeat("a", 40), Head: strings.Repeat("b", 40)}}
			tc.edit(&r)
			if err := r.Validate(); (err == nil) != tc.valid {
				t.Fatalf("valid=%t: %v", tc.valid, err)
			}
		})
	}
}
