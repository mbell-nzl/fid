//go:build integration

package integration

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestV01ParityGoldenOutputs(t *testing.T) {
	binary := buildClosureFIDBinary(t)
	workspace := t.TempDir()
	copyParityFixture(t, filepath.Join("..", "testdata", "v01_parity", "fixture"), workspace)

	for _, fixture := range []struct {
		name   string
		args   []string
		golden string
	}{
		{"issue get epic", []string{"issue", "get", "FID-a10001", "--format", "json"}, "issue-get-epic.json"},
		{"issue get task", []string{"issue", "get", "FID-b20002", "--format", "json"}, "issue-get-task.json"},
		{"issue get subtask", []string{"issue", "get", "FID-c30003", "--format", "json"}, "issue-get-subtask.json"},
		{"issue list epic", []string{"issue", "list", "--kind", "epic", "--format", "json"}, "issue-list-epic.json"},
		{"issue list task", []string{"issue", "list", "--kind", "task", "--format", "json"}, "issue-list-task.json"},
		{"issue list subtask", []string{"issue", "list", "--kind", "subtask", "--format", "json"}, "issue-list-subtask.json"},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			got, _ := closureRunFID(t, binary, workspace, fixture.args...)
			want, err := os.ReadFile(filepath.Join("..", "testdata", "v01_parity", fixture.golden))
			if err != nil {
				t.Fatalf("read golden %s: %v", fixture.golden, err)
			}
			if got != string(want) {
				t.Fatalf("stdout differs from %s:\n%s", fixture.golden, parityDiff(want, []byte(got)))
			}
		})
	}
}

func copyParityFixture(t *testing.T, source, destination string) {
	t.Helper()

	entries, err := os.ReadDir(source)
	if err != nil {
		t.Fatalf("read fixture directory %s: %v", source, err)
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		destinationPath := filepath.Join(destination, entry.Name())
		if entry.IsDir() {
			if err := os.MkdirAll(destinationPath, 0o755); err != nil {
				t.Fatalf("create fixture directory %s: %v", destinationPath, err)
			}
			copyParityFixture(t, sourcePath, destinationPath)
			continue
		}
		copyParityFixtureFile(t, sourcePath, destinationPath)
	}
}

func copyParityFixtureFile(t *testing.T, source, destination string) {
	t.Helper()

	input, err := os.Open(source)
	if err != nil {
		t.Fatalf("open fixture file %s: %v", source, err)
	}
	defer func() { _ = input.Close() }()
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		t.Fatalf("create fixture parent %s: %v", destination, err)
	}
	output, err := os.Create(destination)
	if err != nil {
		t.Fatalf("create fixture file %s: %v", destination, err)
	}
	if _, err := io.Copy(output, input); err != nil {
		_ = output.Close()
		t.Fatalf("copy fixture file %s: %v", source, err)
	}
	if err := output.Close(); err != nil {
		t.Fatalf("close fixture file %s: %v", destination, err)
	}
}

func parityDiff(want, got []byte) string {
	limit := len(want)
	if len(got) < limit {
		limit = len(got)
	}
	for index := 0; index < limit; index++ {
		if want[index] != got[index] {
			return fmt.Sprintf("first difference at byte %d: want %q, got %q", index, want[index], got[index])
		}
	}
	return fmt.Sprintf("length differs: want %d bytes, got %d bytes", len(want), len(got))
}
