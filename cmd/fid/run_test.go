package main

import "testing"

func TestLegacyRunRecordIsNotRegistered(t *testing.T) {
	t.Parallel()
	for _, leaf := range liveExecutableLeafCommands(newRootCommandAt(t.TempDir())) {
		if leaf == "run record" {
			t.Fatal("legacy run record remains registered")
		}
	}
}
