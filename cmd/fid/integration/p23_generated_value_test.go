//go:build integration

package integration

import "testing"

func TestP23GeneratedValueSyntax(t *testing.T) {
	t.Parallel()
	cases := []struct {
		kind, value string
		valid       bool
	}{
		{"RREF", "RREF-abcdef", true},
		{"RREF", "RREF-abcdef01", true},
		{"RREF", "RREF-nothex", false},
		{"RREF", "RREF-abc", false},
		{"EventRef", "EVT-20260930-000001", true},
		{"EventRef", "EVT-no-date-1", false},
		{"digest", "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", true},
		{"digest", "sha256:zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz", false},
		{"timestamp", "2026-09-30T12:00:00Z", true},
		{"timestamp", "tomorrow", false},
	}
	for _, tc := range cases {
		if got := p23ValidGeneratedToken(tc.kind, tc.value); got != tc.valid {
			t.Errorf("%s %q valid = %v, want %v", tc.kind, tc.value, got, tc.valid)
		}
	}
}
