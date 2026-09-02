//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

const p21CorpusDir = "testdata/p21_conformance"

var (
	p21CodeSpan = regexp.MustCompile("`([^`]+)`")
	p21Quoted   = regexp.MustCompile(`"([^"]+)"`)
	p21Kebab    = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)
)

func TestPhase21ConformanceCorpusIsCompleteAndByteExact(t *testing.T) {
	repoRoot := integrationRepositoryRoot()
	publicInterface := string(readP21File(t, filepath.Join(repoRoot, "project-docs/docs/phases/phase-21/public-interface.md")))
	expected := p21ExpectedCorpus(t, publicInterface)
	if os.Getenv("P21_WRITE_CONFORMANCE_FIXTURES") == "1" {
		p21WriteCorpus(t, expected)
	}

	for path, want := range expected {
		got, err := os.ReadFile(filepath.Join(p21CorpusDir, path))
		if err != nil {
			t.Fatalf("read committed conformance fixture %s: %v", path, err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("conformance fixture %s differs from regenerated corpus\n%s", path, p21CorpusDiff(want, got))
		}
	}

	entries, err := p21CorpusEntries(p21CorpusDir)
	if err != nil {
		t.Fatalf("read committed corpus: %v", err)
	}
	if len(entries) != len(expected) {
		t.Fatalf("committed conformance fixture count = %d, want %d", len(entries), len(expected))
	}
	for _, entry := range entries {
		if _, ok := expected[entry]; !ok {
			t.Fatalf("unexpected committed conformance fixture %s", entry)
		}
	}

	manifest := p21DecodeFixture(t, filepath.Join(p21CorpusDir, "manifest.json"))
	observed := p21StringSlice(t, manifest["closed_values"], "manifest closed_values")
	closedMembers := p21ClosedMembers(publicInterface)
	missing := p21Missing(closedMembers, observed)
	if len(missing) != 0 {
		t.Fatalf("closed-value positive control missing %d member(s): %s", len(missing), strings.Join(missing, ", "))
	}
	for _, value := range observed {
		if !p21Kebab.MatchString(value) {
			t.Fatalf("closed public enum value %q is not lowercase kebab-case", value)
		}
		if strings.Contains(value, "_") {
			t.Fatalf("closed-value negative underscore control matched %q", value)
		}
	}
}

func p21ExpectedCorpus(t *testing.T, publicInterface string) map[string][]byte {
	t.Helper()
	operations := p21OperationRows(t, publicInterface)
	for _, operation := range []p21OperationFixture{
		{Operation: "issue.ac.verify", Command: "fid issue ac verify", Kind: "fail-closed"},
		{Operation: "concern.raise", Command: "fid concern raise", Kind: "fail-closed"},
		{Operation: "concern.address", Command: "fid concern address", Kind: "fail-closed"},
		{Operation: "maintenance.vacuum", Command: "fid maintenance vacuum", Kind: "fail-closed"},
	} {
		operations = append(operations, operation)
	}
	if len(operations) == 0 {
		t.Fatal("operation extraction produced no rows")
	}
	failures := p21FailureCodes(t, publicInterface)
	if len(failures) == 0 {
		t.Fatal("failure extraction produced no codes")
	}
	live := p21LiveVariants()

	expected := map[string][]byte{}
	expected["manifest.json"] = p21JSON(t, map[string]any{
		"closed_values": p21ClosedMembers(publicInterface),
		"failures":      failures,
		"live":          live,
		"operations":    operations,
	})
	for _, operation := range operations {
		expected["operations/"+operation.Operation+".json"] = p21JSON(t, operation)
	}
	for _, code := range failures {
		expected["failures/"+code+".json"] = p21JSON(t, map[string]string{"failure_code": code})
	}
	for _, variant := range live {
		encoded, err := json.Marshal(p21LiveFixture(variant))
		if err != nil {
			t.Fatalf("marshal live fixture %s: %v", variant, err)
		}
		expected["live/"+variant+".jsonl"] = append(encoded, '\n')
	}
	return expected
}

type p21OperationFixture struct {
	Operation string `json:"operation"`
	Command   string `json:"command"`
	Kind      string `json:"kind"`
}

func p21OperationRows(t *testing.T, source string) []p21OperationFixture {
	t.Helper()
	section := p21Section(t, source, "### 4.1 Retained standalone and inventory operations", "### 4.2 Operational-kernel operations")
	section += p21Section(t, source, "### 4.2 Operational-kernel operations", "## 5. Fail-Closed And Absent Leaves")
	seen := map[string]bool{}
	var operations []p21OperationFixture
	for _, line := range strings.Split(section, "\n") {
		if !strings.HasPrefix(line, "| `") {
			continue
		}
		cells := strings.Split(line, "|")
		if len(cells) < 3 {
			t.Fatalf("malformed public operation row: %q", line)
		}
		values := p21CodeSpan.FindAllStringSubmatch(cells[1], -1)
		if len(values) != 2 {
			t.Fatalf("operation row must contain operation and command: %q", line)
		}
		operation, command := values[0][1], values[1][1]
		if seen[operation] {
			t.Fatalf("duplicate public operation %q", operation)
		}
		seen[operation] = true
		operations = append(operations, p21OperationFixture{Operation: operation, Command: command, Kind: "accepted"})
	}
	return operations
}

func p21FailureCodes(t *testing.T, source string) []string {
	t.Helper()
	section := p21Section(t, source, "## 7. Complete Failure And Exit Contract", "## 8. Exact `fid.live.v1` Frames")
	seen := map[string]bool{}
	var codes []string
	for _, line := range strings.Split(section, "\n") {
		if !strings.HasPrefix(line, "| `") {
			continue
		}
		cell := strings.Split(line, "|")[1]
		for _, match := range p21CodeSpan.FindAllStringSubmatch(cell, -1) {
			code := match[1]
			if seen[code] {
				continue
			}
			seen[code] = true
			codes = append(codes, code)
		}
	}
	return codes
}

func p21LiveVariants() []string {
	return []string{
		"stdout", "stderr", "control-output-truncated", "control-redaction-applied",
		"control-artifact-rejected", "terminal-failed", "terminal-fid-error",
		"terminal-signaled", "terminal-spawn-failed", "terminal-succeeded",
		"terminal-timed-out", "replay-succeeded", "replay-gate-unsuccessful",
	}
}

func p21LiveFixture(variant string) map[string]any {
	stream := "terminal"
	data := map[string]any{"outcome": "succeeded", "replayed": false}
	switch {
	case variant == "stdout" || variant == "stderr":
		stream = variant
		data = map[string]any{"encoding": "utf-8", "text": "conformance"}
	case strings.HasPrefix(variant, "control-"):
		stream = "control"
		data = map[string]any{"kind": strings.TrimPrefix(variant, "control-")}
	case strings.HasPrefix(variant, "terminal-"):
		data["outcome"] = strings.TrimPrefix(variant, "terminal-")
	case variant == "replay-succeeded":
		data["replayed"] = true
	case variant == "replay-gate-unsuccessful":
		data["outcome"] = "failed"
		data["replayed"] = true
	}
	return map[string]any{"frame_schema": "fid.live.v1", "sequence": 1, "stream": stream, "data": data}
}

func p21ClosedMembers(source string) []string {
	sections := []string{
		p21Section(nil, source, "## 3. Closed Result Types", "## 4. Exact Accepted Leaves"),
		p21Section(nil, source, "## 7. Complete Failure And Exit Contract", "## 8. Exact `fid.live.v1` Frames"),
		p21Section(nil, source, "## 8. Exact `fid.live.v1` Frames", ""),
	}
	seen := map[string]bool{}
	var members []string
	for _, section := range sections {
		for _, match := range p21Quoted.FindAllStringSubmatch(section, -1) {
			value := match[1]
			if !p21Kebab.MatchString(value) || seen[value] {
				continue
			}
			seen[value] = true
			members = append(members, value)
		}
	}
	sort.Strings(members)
	return members
}

func p21Section(t *testing.T, source, start, end string) string {
	if t != nil {
		t.Helper()
	}
	startAt := strings.Index(source, start)
	if startAt < 0 {
		if t != nil {
			t.Fatalf("missing public-interface section %q", start)
		}
		return ""
	}
	section := source[startAt:]
	if end == "" {
		return section
	}
	endAt := strings.Index(section, end)
	if endAt < 0 {
		if t != nil {
			t.Fatalf("missing public-interface section boundary %q", end)
		}
		return ""
	}
	return section[:endAt]
}

func p21JSON(t *testing.T, value any) []byte {
	t.Helper()
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	return append(encoded, '\n')
}

func p21CorpusEntries(root string) ([]string, error) {
	var entries []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		entries = append(entries, relative)
		return nil
	})
	sort.Strings(entries)
	return entries, err
}

func p21WriteCorpus(t *testing.T, expected map[string][]byte) {
	t.Helper()
	for path, contents := range expected {
		fixturePath := filepath.Join(p21CorpusDir, path)
		if err := os.MkdirAll(filepath.Dir(fixturePath), 0o755); err != nil {
			t.Fatalf("create conformance fixture directory for %s: %v", path, err)
		}
		if err := os.WriteFile(fixturePath, contents, 0o644); err != nil {
			t.Fatalf("write conformance fixture %s: %v", path, err)
		}
	}
}

func p21DecodeFixture(t *testing.T, path string) map[string]any {
	t.Helper()
	contents := readP21File(t, path)
	var fixture map[string]any
	if err := json.Unmarshal(contents, &fixture); err != nil {
		t.Fatalf("decode fixture %s: %v", path, err)
	}
	return fixture
}

func p21StringSlice(t *testing.T, value any, label string) []string {
	t.Helper()
	items, ok := value.([]any)
	if !ok {
		t.Fatalf("%s = %#v, want array", label, value)
	}
	stringsOnly := make([]string, 0, len(items))
	for _, item := range items {
		stringValue, ok := item.(string)
		if !ok {
			t.Fatalf("%s contains %#v, want string", label, item)
		}
		stringsOnly = append(stringsOnly, stringValue)
	}
	return stringsOnly
}

func p21Missing(expected, observed []string) []string {
	present := map[string]bool{}
	for _, value := range observed {
		present[value] = true
	}
	var missing []string
	for _, value := range expected {
		if !present[value] {
			missing = append(missing, value)
		}
	}
	return missing
}

func p21CorpusDiff(want, got []byte) string {
	return fmt.Sprintf("want:\n%s\ngot:\n%s", want, got)
}

func readP21File(t *testing.T, path string) []byte {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return contents
}
