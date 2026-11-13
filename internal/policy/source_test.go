package policy

import (
	"reflect"
	"testing"
)

func TestSourceErrorNextAction(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, ref, commit, want string
	}{
		{"unknown default", "", "", "git remote set-head origin <branch>"},
		{"deleted default", "origin/trunk", "", "git remote set-head origin <branch>"},
		{"unreadable policy", "origin/trunk", "abc123", "fix .fid.toml on the default branch, merge it and git fetch"},
		{"known commit without ref", "", "abc123", "fix .fid.toml on the default branch, merge it and git fetch"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := &SourceError{Problem: "unavailable", Ref: test.ref, Commit: test.commit}
			if got := err.NextAction(); got != test.want {
				t.Fatalf("NextAction() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestSourceErrorDetails(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name   string
		source SourceError
		want   map[string]any
	}{
		{"unknown default", SourceError{Problem: "unknown"}, map[string]any{"problem": "unknown", "ref": nil, "commit": nil}},
		{"deleted default", SourceError{Problem: "deleted", Ref: "origin/trunk"}, map[string]any{"problem": "deleted", "ref": "origin/trunk", "commit": nil}},
		{"invalid file", SourceError{Problem: "invalid", Ref: "origin/trunk", Commit: "abc123"}, map[string]any{"problem": "invalid", "ref": "origin/trunk", "commit": "abc123"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := test.source.Details(); !reflect.DeepEqual(got, test.want) {
				t.Fatalf("Details() = %#v, want %#v", got, test.want)
			}
		})
	}
}
