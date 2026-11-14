package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/leaktest"
)

// updateGoldens rewrites golden files from the output under test:
// go test ./cmd/fid -run 'TestStatusGolden|TestShowGolden' -update
var updateGoldens = flag.Bool("update", false, "rewrite golden files")

// assertGolden compares output with testdata/<name>, or rewrites the file
// with -update. JSON output is indented first, so its golden reads well.
// Text, written to a buffer, must have no colour codes and name no kernel
// words.
func assertGolden(t *testing.T, name, output string) {
	t.Helper()
	if filepath.Ext(name) == ".txt" {
		leaktest.AssertNoKernelWords(t, output)
		if strings.Contains(output, "\033") {
			t.Fatalf("golden %s to a buffer has colour codes:\n%q", name, output)
		}
	}
	if filepath.Ext(name) == ".json" {
		var indented bytes.Buffer
		if err := json.Indent(&indented, []byte(output), "", "  "); err != nil {
			t.Fatalf("golden %s: output isn't JSON: %v\n%s", name, err, output)
		}
		output = indented.String()
	}
	path := filepath.Join("testdata", name)
	if *updateGoldens {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(output), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("golden %s: %v (rewrite goldens with -update)", name, err)
	}
	if string(want) != output {
		t.Fatalf("output differs from golden %s (rewrite goldens with -update):\n--- got\n%s\n--- want\n%s", name, output, want)
	}
}
