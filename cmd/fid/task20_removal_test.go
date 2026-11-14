package main

import (
	"strings"
	"testing"
)

func TestTask20V1CommandsAreAbsent(t *testing.T) {
	t.Parallel()
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
		word := path[len(path)-1]
		if path[0] == "run" {
			word = "run"
		}
		if err == nil || !strings.Contains(err.Error(), `unknown command "`+word+`"`) {
			t.Fatalf("fid %v Execute() error = %v, want unknown command", path, err)
		}
	}
}
