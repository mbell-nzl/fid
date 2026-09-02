package main

import (
	"strings"
	"testing"
)

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
		root.SetArgs(path)
		err := root.Execute()
		if err == nil || !strings.Contains(err.Error(), `unknown command "`+path[len(path)-1]+`"`) {
			t.Fatalf("fid %v Execute() error = %v, want unknown command", path, err)
		}
	}
}
