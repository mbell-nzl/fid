package main

import "testing"

func TestRequirementCommandHasReferenceLeaves(t *testing.T) {
	t.Parallel()
	command := newRequirementCommand()
	got := map[string]bool{}
	for _, child := range command.Commands() {
		got[child.Name()] = true
	}
	for _, leaf := range []string{"record-reference", "link", "get", "list"} {
		if !got[leaf] {
			t.Errorf("requirement %s leaf is missing", leaf)
		}
	}
	if got["sync"] {
		t.Fatal("requirement sync leaf remains available")
	}
}
