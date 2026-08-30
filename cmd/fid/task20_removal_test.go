package main

import "testing"

func TestTask20V1CommandsAreAbsent(t *testing.T) {
	root := newRootCommandAt(t.TempDir())
	for _, path := range [][]string{
		{"dev"},
		{"policy"},
		{"issue", "complete"},
		{"issue", "claim"},
		{"issue", "unclaim"},
		{"run", "record"},
	} {
		command, _, err := root.Find(path)
		if err == nil || command != root {
			t.Fatalf("fid %v resolved command %q with err %v, want unknown command", path, command.CommandPath(), err)
		}
	}
}
