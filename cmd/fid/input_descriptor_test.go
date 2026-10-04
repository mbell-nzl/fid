package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/spf13/cobra"
)

func TestAttentionResolveUnknownFieldNamesInputPath(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "request.json")
	if err := os.WriteFile(path, []byte(`{"attention_id":"ATT-1","expected_version":1,"fact_event_ids":[],"bogus":true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	result := executeCommandForTest(t.TempDir(), "--json", "attention", "resolve", "--input", path, "--actor", "ACT-1", "--idempotency-key", "key-1")
	if result.err == nil {
		t.Fatalf("expected invalid input: %#v", result)
	}
	var envelope struct {
		Error struct {
			Details struct {
				Argument string `json:"argument"`
			} `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(result.stderr), &envelope); err != nil {
		t.Fatalf("decode stderr %q: %v", result.stderr, err)
	}
	if envelope.Error.Details.Argument != "input.bogus" {
		t.Fatalf("argument = %q, want input.bogus; stderr=%s", envelope.Error.Details.Argument, result.stderr)
	}
}

func TestAttentionResolveIntegerBoundWritesJSONFailure(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	path := filepath.Join(workspace, "request.json")
	if err := os.WriteFile(path, []byte(`{"attention_id":"ATT-1","expected_version":0,"fact_event_ids":["EVT-1"]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	result := executeCommandForTest(workspace, "--json", "attention", "resolve", "--input", path, "--actor", "ACT-1", "--idempotency-key", "bound-1")
	if cli.ExitCode(result.err) != 2 || result.stdout != "" {
		t.Fatalf("exit = %d, stdout = %q, stderr = %q", cli.ExitCode(result.err), result.stdout, result.stderr)
	}
	var envelope struct {
		ContractVersion string           `json:"contract_version"`
		OK              bool             `json:"ok"`
		Operation       string           `json:"operation"`
		Event           *json.RawMessage `json:"event"`
		Error           struct {
			Category string `json:"category"`
			Code     string `json:"code"`
			Details  struct {
				Argument string   `json:"argument"`
				Value    *string  `json:"value"`
				Allowed  []string `json:"allowed"`
			} `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(result.stderr), &envelope); err != nil {
		t.Fatalf("decode stderr %q: %v", result.stderr, err)
	}
	if envelope.ContractVersion != "fid-cli.v1" || envelope.OK || envelope.Operation != "attention.resolve" || envelope.Event != nil ||
		envelope.Error.Category != "invalid-input" || envelope.Error.Code != "invalid-argument" ||
		envelope.Error.Details.Argument != "input.expected_version" || envelope.Error.Details.Value == nil || *envelope.Error.Details.Value != "0" ||
		!reflect.DeepEqual(envelope.Error.Details.Allowed, []string{"integer range 1..unbounded"}) {
		t.Fatalf("unexpected structural failure: %s", result.stderr)
	}
	if _, err := os.Stat(filepath.Join(workspace, ".fid")); !os.IsNotExist(err) {
		t.Fatalf("structural failure touched workspace: %v", err)
	}
}

func TestInputDescriptorRegistryMatchesLiveTreeAndTypes(t *testing.T) {
	t.Parallel()
	root := newRootCommandAt(t.TempDir())
	live := make(map[string]bool)
	var walk func(*cobra.Command)
	walk = func(command *cobra.Command) {
		if len(command.Commands()) == 0 {
			if command.Flags().Lookup("input") != nil {
				live[rootCommandKey(command)] = true
			}
			return
		}
		for _, child := range command.Commands() {
			walk(child)
		}
	}
	walk(root)
	for key, descriptor := range inputDescriptorRegistry {
		if pendingInputDescriptors[key] || !live[key] {
			t.Errorf("registered %s overlaps fence or is not live", key)
		}
		if err := cli.CheckInputType(descriptor); err != nil {
			t.Errorf("%s: %v", key, err)
		}
	}
	for key := range pendingInputDescriptors {
		if !live[key] {
			t.Errorf("pending %s is not live", key)
		}
	}
	for key := range live {
		if _, registered := inputDescriptorRegistry[key]; !registered && !pendingInputDescriptors[key] {
			t.Errorf("live --input leaf %s is unaccounted", key)
		}
	}
	mismatch := stopResolveDescriptor()
	mismatch.Fields[0].Name = "wrong_stop_id"
	if err := cli.CheckInputType(mismatch); err == nil {
		t.Error("seeded descriptor/type mismatch passed")
	}
}

func TestEveryLiveLeafHasOneLineShort(t *testing.T) {
	t.Parallel()
	root := newRootCommandAt(t.TempDir())
	var walk func(*cobra.Command)
	walk = func(command *cobra.Command) {
		if len(command.Commands()) == 0 && (strings.TrimSpace(command.Short) == "" || strings.Contains(command.Short, "\n")) {
			t.Errorf("%s has no one-line Short", command.CommandPath())
		}
		for _, child := range command.Commands() {
			walk(child)
		}
	}
	walk(root)
}

func TestAdoptedInputHelpAndExamples(t *testing.T) {
	t.Parallel()
	templates := map[string]string{
		"stop raise": `{
  "subject": {
    "type": "run",
    "id": "<subject-id>"
  },
  "code": "executor-liveness-unknown",
  "blocked_actions": ["claim.acquire"],
  "details": {
    "run_id": "<run-id>",
    "executor_id": "<executor-id>",
    "observation_id": null
  },
  "resolution_condition": {
    "type": "executor-terminated",
    "subject": {
      "type": "run",
      "id": "<subject-id>"
    },
    "budget_name": null,
    "capability_id": null
  },
  "resolver_capability": "<capability-id>"
}` + "\n",
		"stop resolve": `{
  "stop_id": "<stop-id>",
  "fact_event_id": "<fact-event-id>"
}` + "\n",
		"attention resolve": `{
  "attention_id": "<attention-id>",
  "expected_version": 1,
  "fact_event_ids": []
}` + "\n",
	}
	for _, key := range []string{"stop raise", "stop resolve", "attention resolve"} {
		descriptor := inputDescriptorRegistry[key]
		t.Run(key, func(t *testing.T) {
			t.Parallel()
			workspace := t.TempDir()
			parts := strings.Split(key, " ")
			help := executeCommandForTest(workspace, append(parts, "--help")...)
			if help.err != nil || !strings.Contains(help.stdout, descriptor.Purpose) || !strings.Contains(help.stdout, "fid "+key+" --example") {
				t.Fatalf("help for %s = %#v", key, help)
			}
			if key == "stop raise" {
				const fullHelpSHA256 = "b3a71078d99b3047a9ffdba35e599bd8298da2e6f6e9b06f709ba3283af5bf75"
				if got := fmt.Sprintf("%x", sha256.Sum256([]byte(help.stdout))); got != fullHelpSHA256 {
					t.Fatalf("full stop raise help changed: SHA-256 %s, want %s; help=%s", got, fullHelpSHA256, help.stdout)
				}
			}
			for _, field := range descriptor.Fields {
				if !strings.Contains(help.stdout, "input."+field.Name+":") {
					t.Errorf("help missing input.%s", field.Name)
				}
			}
			example := executeCommandForTest(workspace, append(parts, "--example")...)
			if example.err != nil || example.stderr != "" {
				t.Fatalf("example for %s = %#v", key, example)
			}
			if example.stdout != templates[key] {
				t.Fatalf("template for %s = %q, want %q", key, example.stdout, templates[key])
			}
			decoder := json.NewDecoder(strings.NewReader(example.stdout))
			decoder.DisallowUnknownFields()
			request := reflect.New(descriptor.Type).Interface()
			if err := decoder.Decode(request); err != nil {
				t.Fatalf("template does not decode: %v; %s", err, example.stdout)
			}
			if len(example.stdout) == 0 || example.stdout[0] != '{' || !strings.Contains(example.stdout, "\n  \"") {
				t.Fatalf("template is not indented JSON: %q", example.stdout)
			}
			entries, err := os.ReadDir(workspace)
			if err != nil || len(entries) != 0 {
				t.Fatalf("example changed non-project directory: entries=%v, err=%v", entries, err)
			}
		})
	}
}

func TestExampleRejectsOtherInputsInFixedOrder(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		args     []string
		argument string
	}{
		{[]string{"--json", "--project", "/missing", "--input", "/missing"}, "--json"},
		{[]string{"--project", "/missing"}, "--project"},
		{[]string{"--input", "/missing"}, "--input"},
		{[]string{"--actor", "ACT-1"}, "--actor"},
		{[]string{"--idempotency-key", "key-1"}, "--idempotency-key"},
		{[]string{"extra"}, "arguments"},
	} {
		args := append([]string{"attention", "resolve", "--example"}, test.args...)
		result := executeCommandForTest(t.TempDir(), args...)
		if cli.ExitCode(result.err) != 2 {
			t.Errorf("%v: exit = %d, stderr=%s", args, cli.ExitCode(result.err), result.stderr)
			continue
		}
		if strings.Contains(argsToText(args), "--json") {
			var envelope struct {
				Error struct {
					Details struct {
						Argument string `json:"argument"`
					} `json:"details"`
				} `json:"error"`
			}
			if err := json.Unmarshal([]byte(result.stderr), &envelope); err != nil || envelope.Error.Details.Argument != test.argument {
				t.Errorf("%v: JSON conflict = %s (%v)", args, result.stderr, err)
			}
		} else if !strings.Contains(result.stderr, "--example cannot be combined with") {
			t.Errorf("%v: Human conflict = %q", args, result.stderr)
		}
	}
}

func argsToText(args []string) string { return strings.Join(args, " ") }

func TestAdoptedStructuralValidationPaths(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, key, payload, path string
	}{
		{"top unknown", "attention resolve", `{"attention_id":"ATT-1","expected_version":1,"fact_event_ids":[],"bogus":1}`, "input.bogus"},
		{"top missing", "attention resolve", `{"attention_id":"ATT-1","fact_event_ids":[]}`, "input.expected_version"},
		{"top null", "attention resolve", `{"attention_id":null,"expected_version":1,"fact_event_ids":[]}`, "input.attention_id"},
		{"top type", "attention resolve", `{"attention_id":"ATT-1","expected_version":"one","fact_event_ids":[]}`, "input.expected_version"},
		{"nested unknown", "stop raise", `{"subject":{"type":"run","id":"RUN-1","bogus":1},"code":"operator-stop","blocked_actions":["run.start"],"details":{"reason":"x","reference":"x"},"resolution_condition":{"type":"work-ended","subject":{"type":"run","id":"RUN-1"},"budget_name":null,"capability_id":null},"resolver_capability":"stop.resolve"}`, "input.subject.bogus"},
		{"nested enum", "stop raise", `{"subject":{"type":"bogus","id":"RUN-1"},"code":"operator-stop","blocked_actions":["run.start"],"details":{"reason":"x","reference":"x"},"resolution_condition":{"type":"work-ended","subject":{"type":"run","id":"RUN-1"},"budget_name":null,"capability_id":null},"resolver_capability":"stop.resolve"}`, "input.subject.type"},
		{"nested missing", "stop raise", `{"subject":{"type":"run"},"code":"operator-stop","blocked_actions":["run.start"],"details":{"reason":"x","reference":"x"},"resolution_condition":{"type":"work-ended","subject":{"type":"run","id":"RUN-1"},"budget_name":null,"capability_id":null},"resolver_capability":"stop.resolve"}`, "input.subject.id"},
		{"nested deep type", "stop raise", `{"subject":{"type":"run","id":"RUN-1"},"code":"operator-stop","blocked_actions":["run.start"],"details":{"reason":"x","reference":"x"},"resolution_condition":{"type":"work-ended","subject":{"type":"run","id":1},"budget_name":null,"capability_id":null},"resolver_capability":"stop.resolve"}`, "input.resolution_condition.subject.id"},
		{"duplicate", "attention resolve", `{"attention_id":"ATT-1","attention_id":"ATT-2","expected_version":1,"fact_event_ids":[]}`, "input.attention_id"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var field *cli.InputError
			if err := cli.ValidateInput([]byte(test.payload), inputDescriptorRegistry[test.key]); err == nil || !strings.Contains(err.Error(), test.path) {
				t.Fatalf("validation = %v, want %s", err, test.path)
			} else if typed, ok := err.(*cli.InputError); ok {
				field = typed
			}
			if field == nil || field.Path != test.path || len(field.Allowed) == 0 {
				t.Fatalf("field error = %#v", field)
			}
		})
	}
}

func TestStructuralFailuresCoverEachAdoptedLeaf(t *testing.T) {
	t.Parallel()
	baseline := map[string]string{
		"stop resolve":      `{"stop_id":"STP-1","fact_event_id":"EVT-1"}`,
		"attention resolve": `{"attention_id":"ATT-1","expected_version":1,"fact_event_ids":["EVT-1"]}`,
		"stop raise":        `{"subject":{"type":"run","id":"RUN-1"},"code":"operator-stop","blocked_actions":["run.start"],"details":{"reason":"reason","reference":"REF"},"resolution_condition":{"type":"work-ended","subject":{"type":"run","id":"RUN-1"},"budget_name":null,"capability_id":null},"resolver_capability":"stop.resolve"}`,
	}
	for _, test := range []struct {
		name, leaf, old, replacement, path string
	}{
		{"stop resolve unknown", "stop resolve", `"fact_event_id":"EVT-1"`, `"fact_event_id":"EVT-1","bogus":1`, "input.bogus"},
		{"stop resolve duplicate", "stop resolve", `"stop_id":"STP-1"`, `"stop_id":"STP-1","stop_id":"STP-2"`, "input.stop_id"},
		{"stop resolve missing", "stop resolve", `"stop_id":"STP-1",`, ``, "input.stop_id"},
		{"stop resolve null", "stop resolve", `"stop_id":"STP-1"`, `"stop_id":null`, "input.stop_id"},
		{"stop resolve type", "stop resolve", `"stop_id":"STP-1"`, `"stop_id":1`, "input.stop_id"},
		{"stop resolve ID", "stop resolve", `"stop_id":"STP-1"`, `"stop_id":""`, "input.stop_id"},
		{"attention unknown", "attention resolve", `"expected_version":1`, `"expected_version":1,"bogus":1`, "input.bogus"},
		{"attention duplicate", "attention resolve", `"attention_id":"ATT-1"`, `"attention_id":"ATT-1","attention_id":"ATT-2"`, "input.attention_id"},
		{"attention missing", "attention resolve", `"attention_id":"ATT-1",`, ``, "input.attention_id"},
		{"attention null", "attention resolve", `"attention_id":"ATT-1"`, `"attention_id":null`, "input.attention_id"},
		{"attention type", "attention resolve", `"expected_version":1`, `"expected_version":"1"`, "input.expected_version"},
		{"attention bound", "attention resolve", `"expected_version":1`, `"expected_version":0`, "input.expected_version"},
		{"attention ID", "attention resolve", `"attention_id":"ATT-1"`, `"attention_id":""`, "input.attention_id"},
		{"stop raise unknown", "stop raise", `"code":"operator-stop"`, `"code":"operator-stop","bogus":1`, "input.bogus"},
		{"stop raise duplicate", "stop raise", `"code":"operator-stop"`, `"code":"operator-stop","code":"operator-stop"`, "input.code"},
		{"stop raise missing", "stop raise", `"code":"operator-stop",`, ``, "input.code"},
		{"stop raise null", "stop raise", `"code":"operator-stop"`, `"code":null`, "input.code"},
		{"stop raise type", "stop raise", `"code":"operator-stop"`, `"code":1`, "input.code"},
		{"stop raise enum", "stop raise", `"code":"operator-stop"`, `"code":"bogus"`, "input.code"},
		{"stop raise format", "stop raise", `"resolver_capability":"stop.resolve"`, `"resolver_capability":""`, "input.resolver_capability"},
		{"nested unknown", "stop raise", `"id":"RUN-1"},"code"`, `"id":"RUN-1","bogus":1},"code"`, "input.subject.bogus"},
		{"nested duplicate", "stop raise", `"id":"RUN-1"},"code"`, `"id":"RUN-1","id":"RUN-2"},"code"`, "input.subject.id"},
		{"nested missing", "stop raise", `"type":"run","id":"RUN-1"},"code"`, `"type":"run"},"code"`, "input.subject.id"},
		{"nested null", "stop raise", `"id":"RUN-1"},"code"`, `"id":null},"code"`, "input.subject.id"},
		{"nested type", "stop raise", `"id":"RUN-1"},"code"`, `"id":1},"code"`, "input.subject.id"},
		{"nested enum", "stop raise", `"type":"run","id":"RUN-1"},"code"`, `"type":"bogus","id":"RUN-1"},"code"`, "input.subject.type"},
		{"nested format", "stop raise", `"id":"RUN-1"},"code"`, `"id":""},"code"`, "input.subject.id"},
		{"deep unknown", "stop raise", `"id":"RUN-1"},"budget_name"`, `"id":"RUN-1","bogus":1},"budget_name"`, "input.resolution_condition.subject.bogus"},
		{"array index", "stop raise", `"blocked_actions":["run.start"]`, `"blocked_actions":[1]`, "input.blocked_actions[0]"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			payload := strings.Replace(baseline[test.leaf], test.old, test.replacement, 1)
			if payload == baseline[test.leaf] {
				t.Fatal("test mutation did not match baseline")
			}
			invalid, ok := cli.ValidateInput([]byte(payload), inputDescriptorRegistry[test.leaf]).(*cli.InputError)
			if !ok || invalid.Path != test.path || len(invalid.Allowed) == 0 {
				t.Fatalf("%s: result = %#v, want %s", test.name, invalid, test.path)
			}
		})
	}
}
