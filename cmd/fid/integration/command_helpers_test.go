//go:build integration

package integration

// Shared process/JSON/filesystem helpers have no kernel fixture dependency.
import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func p21Input(t *testing.T, workspace, name string, value map[string]any) string {
	t.Helper()
	contents, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal %s input: %v", name, err)
	}
	path := filepath.Join(workspace, name+".json")
	if err := os.WriteFile(path, contents, 0o600); err != nil {
		t.Fatalf("write %s input: %v", name, err)
	}
	return path
}

func p21Run(t *testing.T, binary, workspace string, args ...string) map[string]any {
	t.Helper()
	stdout, stderr, err := closureRunCommand(workspace, binary, args...)
	if err != nil || stderr != "" {
		t.Fatalf("fid %s = err:%v stdout:%q stderr:%q", strings.Join(args, " "), err, stdout, stderr)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(stdout), &envelope); err != nil {
		t.Fatalf("decode fid %s response: %v; stdout=%q", strings.Join(args, " "), err, stdout)
	}
	return envelope
}

func p21PathString(t *testing.T, value map[string]any, path ...string) string {
	t.Helper()
	item := p21Path(t, value, path...)
	text, ok := item.(string)
	if !ok || text == "" {
		t.Fatalf("path %s = %#v, want non-empty string", strings.Join(path, "."), item)
	}
	return text
}

func p21Path(t *testing.T, value map[string]any, path ...string) any {
	t.Helper()
	var item any = value
	for _, key := range path {
		object, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("path %s traverses %#v, want object", strings.Join(path, "."), item)
		}
		item = object[key]
	}
	return item
}

func p21JournalEventCount(t *testing.T, workspace string) int {
	t.Helper()
	return len(p21JournalEventTypes(t, workspace))
}

func p21AssertJournalTypes(t *testing.T, workspace string, want []string) {
	t.Helper()
	got := p21JournalEventTypes(t, workspace)
	if len(got) != len(want) {
		t.Fatalf("journal Event count = %d, want %d; types=%v", len(got), len(want), got)
	}
	for index, eventType := range want {
		if got[index] != eventType {
			t.Fatalf("journal Event %d = %q, want %q; types=%v", index+1, got[index], eventType, got)
		}
	}
}

func p21JournalEventTypes(t *testing.T, workspace string) []string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(config.LedgerDir(workspace), "journal"))
	if err != nil {
		t.Fatalf("read journal: %v", err)
	}
	var types []string
	for _, entry := range entries {
		contents := readP21File(t, filepath.Join(config.LedgerDir(workspace), "journal", entry.Name()))
		for _, line := range bytes.FieldsFunc(contents, func(r rune) bool { return r == '\n' }) {
			var event struct {
				Type string `json:"event_type"`
			}
			if err := json.Unmarshal(line, &event); err != nil {
				t.Fatalf("decode journal Event: %v", err)
			}
			if event.Type == "" {
				t.Fatal("journal Event has empty type")
			}
			types = append(types, event.Type)
		}
	}
	return types
}

func p21JournalBytes(t *testing.T, workspace string) []byte {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(config.LedgerDir(workspace), "journal"))
	if err != nil {
		t.Fatalf("read journal: %v", err)
	}
	var journal []byte
	for _, entry := range entries {
		journal = append(journal, readP21File(t, filepath.Join(config.LedgerDir(workspace), "journal", entry.Name()))...)
	}
	return journal
}

func p21JSON(t *testing.T, value any) []byte {
	t.Helper()
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	return append(encoded, '\n')
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

func readP21File(t *testing.T, path string) []byte {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return contents
}

func p23Envelope(t *testing.T, stdout string) map[string]any {
	t.Helper()
	var result map[string]any
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("decode Review response: %v; stdout=%q", err, stdout)
	}
	return result
}

func p23Number(t *testing.T, value any) int {
	t.Helper()
	number, ok := value.(float64)
	if !ok {
		t.Fatalf("number = %#v, want float64 JSON number", value)
	}
	return int(number)
}

func p23NumberValue(value any) int {
	switch n := value.(type) {
	case int:
		return n
	case float64:
		return int(n)
	}
	return 0
}

func p23Clone(value map[string]any) map[string]any {
	payload, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	var clone map[string]any
	if err := json.Unmarshal(payload, &clone); err != nil {
		panic(err)
	}
	return clone
}

func p22Write(t *testing.T, workspace, name, contents string) {
	t.Helper()
	path := filepath.Join(workspace, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create the directory for %s: %v", name, err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func p22Git(t *testing.T, workspace string, arguments ...string) {
	t.Helper()
	closureRun(t, workspace, "git", arguments...)
}

func p22GitOutput(t *testing.T, root string, args ...string) string {
	t.Helper()
	output, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v %s", args, err, output)
	}
	return strings.TrimSpace(string(output))
}

func p22SameMembers(t *testing.T, name string, got, want []string) {
	t.Helper()
	got, want = append([]string{}, got...), append([]string{}, want...)
	sort.Strings(got)
	sort.Strings(want)
	if len(got) == 0 || !reflect.DeepEqual(got, want) {
		t.Fatalf("%s = %v, want exact %v", name, got, want)
	}
}
