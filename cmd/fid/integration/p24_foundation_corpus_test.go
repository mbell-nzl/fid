//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

var p24FoundationFields = map[string]string{
	"issue.create": "--title", "activation.prepare": "kit_identity",
	"activation.resource": "evidence_ref", "activation.component": "content",
	"activation.commit": "expected_version", "activation.finalize": "expected_version",
	"gate.register": "definition_version", "envelope.create": "outcome",
	"envelope.activate": "expected_revision", "claim.acquire": "ttl_seconds",
	"run.start": "claim_version",
}

func TestP24FoundationKeyedLeafConflictAndCanonicalRetryCorpus(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	p24StartKeyedCapture(workspace)
	_ = p22NewFIDRecoveryAdapterAt(t, binary, workspace, "agent", 3600, []string{}, nil)
	captured := p24EndKeyedCapture(workspace)
	for operation, name := range p24FoundationFields {
		entry, ok := captured[operation]
		if !ok {
			t.Fatalf("fixture did not capture keyed %s; captured = %v", operation, p24CapturedNames(captured))
		}
		p24ExerciseCapturedLeaf(t, workspace, binary, operation, name, entry)
	}
	for operation, second := range map[string]string{"activation.prepare": "resources", "envelope.create": "non_goals", "run.start": "executor", "claim.acquire": "mode", "gate.register": "definition_id"} {
		entry := captured[operation]
		p24AssertFamilyConflict(t, workspace, binary, operation, entry, []string{p24FoundationFields[operation], second})
		p24AssertFamilyConflict(t, workspace, binary, operation, entry, []string{"--actor"})
	}
	p24AssertFamilyConflict(t, workspace, binary, "envelope.create", captured["envelope.create"], []string{"protected_decisions"})
	p24AssertRebuildMapParity(t, workspace, binary)
}

func p24ExerciseCapturedLeaf(t *testing.T, workspace, binary, operation, name string, entry p24CapturedKeyedCommand) {
	t.Helper()
	p24AssertStoredFieldNames(t, workspace, operation, entry)
	before := p21JournalEventCount(t, workspace)
	canonical := append([]string(nil), entry.args...)
	if len(entry.input) != 0 {
		var request map[string]any
		if err := json.Unmarshal(entry.input, &request); err != nil {
			t.Fatal(err)
		}
		if _, present := request[name]; !present {
			t.Fatalf("%s planned field %s absent from caller input %v", operation, name, request)
		}
		canonical = p24ReplaceInput(t, workspace, canonical, "canonical-"+strings.ReplaceAll(operation, ".", "-"), p24ReverseJSON(request))
	} else {
		canonical = p24PermuteFlagPairs(canonical)
	}
	stdout, stderr, err := closureRunCommand(workspace, binary, canonical...)
	if err != nil || stderr != "" {
		t.Fatalf("%s canonical retry = %v stderr=%q", operation, err, stderr)
	}
	exactResponse := operation == "hold.park" || operation == "task.takeover" || operation == "task.answer"
	byteWant := strings.Replace(entry.stdout, `"replayed":false`, `"replayed":true`, 1)
	if exactResponse {
		byteWant = entry.stdout
	}
	if (!exactResponse && byteWant == entry.stdout) || stdout != byteWant {
		at := 0
		for at < len(stdout) && at < len(byteWant) && stdout[at] == byteWant[at] {
			at++
		}
		start, end := at-80, at+160
		if start < 0 {
			start = 0
		}
		if end > len(stdout) {
			end = len(stdout)
		}
		wantEnd := at + 160
		if wantEnd > len(byteWant) {
			wantEnd = len(byteWant)
		}
		t.Fatalf("%s canonical retry bytes first differ at %d: got %q want %q", operation, at, stdout[start:end], byteWant[start:wantEnd])
	}
	want := p23Envelope(t, entry.stdout)
	if !exactResponse {
		want["result"].(map[string]any)["replayed"] = true
	}
	if got := p23Envelope(t, stdout); !reflect.DeepEqual(got, want) {
		t.Fatalf("%s canonical retry = %#v, want %#v", operation, got, want)
	}
	changed := append([]string(nil), entry.args...)
	if len(entry.input) != 0 {
		var request map[string]any
		if err := json.Unmarshal(entry.input, &request); err != nil {
			t.Fatal(err)
		}
		p24ChangeFoundationValue(t, request, name)
		body, err := json.Marshal(request)
		if err != nil {
			t.Fatal(err)
		}
		changed = p24ReplaceInput(t, workspace, changed, "changed-"+strings.ReplaceAll(operation, ".", "-"), body)
	} else {
		if strings.HasPrefix(name, "--") {
			for i, arg := range changed {
				if arg == name {
					changed[i+1] += "-changed"
				}
			}
		} else {
			for i, arg := range changed {
				if arg == "link" && i+3 < len(changed) && (changed[i+1] == "add" || changed[i+1] == "remove") && name == "relation" {
					changed[i+3] = "relates_to"
				}
				if arg == "label" && i+3 < len(changed) && changed[i+1] == "add" && (changed[i+2] == "--actor" || changed[i+2] == "--idempotency-key") {
					changed[i+2] = "p24-other-label"
				}
			}
			if operation == "issue.label.add" || operation == "issue.label.remove" {
				for i, arg := range changed {
					if arg == "issue" && i+4 < len(changed) && changed[i+1] == "label" {
						changed[i+4] = "p24-other-label"
					}
				}
			}
		}
	}
	conflict, conflictErr, runErr := closureRunCommand(workspace, binary, changed...)
	if runErr == nil || processExitCode(t, runErr) != 4 || conflict != "" {
		t.Fatalf("%s changed-input conflict = %v stdout=%q stderr=%q", operation, runErr, conflict, conflictErr)
	}
	if got := p21Path(t, p23Envelope(t, conflictErr), "error", "details", "differing_fields"); !reflect.DeepEqual(got, []any{name}) {
		t.Fatalf("%s differing_fields = %#v, want [%s]", operation, got, name)
	}
	if after := p21JournalEventCount(t, workspace); after != before {
		t.Fatalf("%s retry/conflict appended Events: %d -> %d", operation, before, after)
	}
}

func p24AssertStoredFieldNames(t *testing.T, workspace, operation string, entry p24CapturedKeyedCommand) {
	t.Helper()
	key := ""
	for i, arg := range entry.args {
		if arg == "--idempotency-key" && i+1 < len(entry.args) {
			key = entry.args[i+1]
		}
	}
	if key == "" {
		t.Fatalf("%s has no key", operation)
	}
	want := map[string]bool{"--actor": true}
	if len(entry.input) != 0 {
		var members map[string]json.RawMessage
		if err := json.Unmarshal(entry.input, &members); err != nil {
			t.Fatal(err)
		}
		for name, raw := range members {
			if operation == "review.record" && name == "inapplicable_findings" && string(raw) == "[]" {
				continue
			}
			want[name] = true
		}
	} else {
		for i, arg := range entry.args {
			if strings.HasPrefix(arg, "--") && arg != "--project" && arg != "--json" && arg != "--idempotency-key" && i+1 < len(entry.args) {
				want[arg] = true
			}
		}
		for _, name := range p24PositionalNames[operation] {
			want[name] = true
		}
	}
	conn := p22FaultDB(t, workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var encoded string
	if err := conn.QueryRow(`SELECT request_field_digests_json FROM idempotency_results WHERE operation_scope LIKE ? AND idempotency_key=?`, operation+":%", key).Scan(&encoded); err != nil {
		t.Fatal(err)
	}
	var fields map[string]string
	if err := json.Unmarshal([]byte(encoded), &fields); err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for name := range fields {
		got[name] = true
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s input-as-written fields = %v, want %v; input=%s", operation, got, want, entry.input)
	}
}

var p24PositionalNames = map[string][]string{
	"task.answer":   {"task", "text"},
	"task.takeover": {"task"},
	"issue.update":  {"id"}, "issue.transition": {"id"}, "issue.ac.add": {"issue-id"},
	"issue.label.add": {"issue-id", "label"}, "issue.label.remove": {"issue-id", "label"},
	"label.add": {"name"}, "link.add": {"from-id", "relation", "to-id"}, "link.remove": {"from-id", "relation", "to-id"},
}

func p24CapturedNames(entries map[string]p24CapturedKeyedCommand) []string {
	names := make([]string, 0, len(entries))
	for name := range entries {
		names = append(names, name)
	}
	return names
}

func p24ReplaceInput(t *testing.T, workspace string, args []string, name string, body []byte) []string {
	t.Helper()
	path := filepath.Join(workspace, "p24-"+name+".json")
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
	for i, arg := range args {
		if arg == "--input" {
			args[i+1] = path
			return args
		}
	}
	t.Fatal("missing --input")
	return nil
}

func p24PermuteFlagPairs(args []string) []string {
	start := -1
	for i := 0; i < len(args); i++ {
		if !strings.HasPrefix(args[i], "--") || args[i] == "--json" || args[i] == "--project" {
			continue
		}
		start = i
		break
	}
	if start < 0 || (len(args)-start)%2 != 0 {
		return args
	}
	for left, right := start, len(args)-2; left < right; left, right = left+2, right-2 {
		args[left], args[left+1], args[right], args[right+1] = args[right], args[right+1], args[left], args[left+1]
	}
	return args
}

func p24ChangeFoundationValue(t *testing.T, request map[string]any, name string) {
	t.Helper()
	if _, found := request[name]; !found && name == "protected_decisions" {
		request[name] = []any{}
		return
	}
	switch value := request[name].(type) {
	case string:
		if name == "mode" {
			request[name] = "review"
		} else if name == "reconciliation_digest" {
			request[name] = "sha256:" + p22Zeros(64)
		} else {
			request[name] = value + "-p24-changed"
		}
	case float64:
		request[name] = value + 1
	case map[string]any:
		if reason, ok := value["reason"].(string); ok {
			value["reason"] = reason + "-p24-changed"
			return
		}
		if id, ok := value["id"].(string); ok {
			value["id"] = id + "-p24-changed"
			return
		}
		t.Fatalf("%s has no reason field: %#v", name, value)
	case []any:
		if len(value) == 0 {
			request[name] = []any{"phase 24 non-goal"}
			return
		}
		first, ok := value[0].(map[string]any)
		if !ok {
			t.Fatalf("%s first member = %#v", name, value[0])
		}
		for _, member := range []string{"provenance_ref", "definition_id", "actor_id", "new_identity"} {
			if text, ok := first[member].(string); ok {
				first[member] = text + "-p24-changed"
				return
			}
		}
		t.Fatalf("%s has no safe text member: %#v", name, first)
	default:
		t.Fatalf("%s cannot be changed: %#v", name, request[name])
	}
}

func p24AssertFamilyConflict(t *testing.T, workspace, binary, operation string, entry p24CapturedKeyedCommand, changedNames []string) {
	t.Helper()
	before := p21JournalEventCount(t, workspace)
	args := append([]string(nil), entry.args...)
	if len(entry.input) != 0 {
		var request map[string]any
		if err := json.Unmarshal(entry.input, &request); err != nil {
			t.Fatal(err)
		}
		for _, name := range changedNames {
			if name != "--actor" {
				p24ChangeFoundationValue(t, request, name)
			}
		}
		body, err := json.Marshal(request)
		if err != nil {
			t.Fatal(err)
		}
		args = p24ReplaceInput(t, workspace, args, "family-"+strings.ReplaceAll(operation, ".", "-"), body)
	} else {
		for _, name := range changedNames {
			if name != "--actor" {
				args = p24ChangeFlagOrPositional(t, operation, args, name)
			}
		}
	}
	for _, name := range changedNames {
		if name == "--actor" {
			args = p24ChangeActor(args)
		}
	}
	stdout, stderr, err := closureRunCommand(workspace, binary, args...)
	if err == nil || processExitCode(t, err) != 4 || stdout != "" {
		t.Fatalf("%s %v conflict args=%v = %v stdout=%q stderr=%q", operation, changedNames, args, err, stdout, stderr)
	}
	want := append([]string(nil), changedNames...)
	sort.Strings(want)
	got := p21Path(t, p23Envelope(t, stderr), "error", "details", "differing_fields")
	values := make([]any, len(want))
	for i, name := range want {
		values[i] = name
	}
	if !reflect.DeepEqual(got, values) {
		t.Fatalf("%s %v differing_fields = %#v, want %#v", operation, changedNames, got, values)
	}
	if after := p21JournalEventCount(t, workspace); after != before {
		t.Fatalf("%s %v conflict appended Events: %d -> %d", operation, changedNames, before, after)
	}
}

func p24ChangeFlagOrPositional(t *testing.T, operation string, args []string, name string) []string {
	t.Helper()
	for i, arg := range args {
		if operation == "task.answer" && arg == "answer" && name == "text" {
			args[i+2] += "-changed"
			return args
		}
		if arg == name && strings.HasPrefix(name, "--") && i+1 < len(args) {
			if name == "--expected-revision" {
				value, err := strconv.Atoi(args[i+1])
				if err != nil {
					t.Fatal(err)
				}
				args[i+1] = strconv.Itoa(value + 1)
			} else {
				args[i+1] += "-p24-changed"
			}
			return args
		}
		if operation == "link.add" && arg == "link" && i+4 < len(args) && name == "to-id" {
			args[i+4] = "P22-missing"
			return args
		}
		if operation == "link.add" && arg == "link" && i+3 < len(args) && name == "relation" {
			args[i+3] = "relates_to"
			return args
		}
	}
	if name == "--description" {
		return append(args, name, "phase 24 description")
	}
	t.Fatalf("%s has no public field %s in %v", operation, name, args)
	return nil
}

func p24ChangeActor(args []string) []string {
	for i, arg := range args {
		if arg == "--actor" && i+1 < len(args) {
			if strings.HasPrefix(args[i+1], "agent:") {
				args[i+1] = "agent:other"
			} else if strings.HasPrefix(args[i+1], "human:") {
				args[i+1] = "human:other"
			} else if args[i+1] == "actor-denied" {
				args[i+1] = "actor-reviewer-two"
			} else {
				args[i+1] = "actor-denied"
			}
			return args
		}
	}
	return append(args, "--actor", "agent:other")
}
