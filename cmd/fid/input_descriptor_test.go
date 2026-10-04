package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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

func TestAdoptedExampleProcessExclusivity(t *testing.T) {
	t.Parallel()
	leaves := []string{"stop raise", "stop resolve", "attention resolve", "decision record-reference",
		"requirement record-reference", "requirement link", "human-outcome record", "concern record", "concern dispose", "deviation record"}
	conflicts := []struct {
		name, argument string
		value          any
		args           func(string) []string
	}{
		{"json first", "--json", nil, func(dir string) []string {
			return []string{"--json", "--project", filepath.Join(dir, "absent-project"), "--input", filepath.Join(dir, "absent-input"), "--actor", "ACT-1", "--idempotency-key", "key-1", "extra"}
		}},
		{"project first", "--project", "ABSENT", func(dir string) []string {
			return []string{"--project", filepath.Join(dir, "absent-project"), "--input", filepath.Join(dir, "absent-input"), "--actor", "ACT-1", "--idempotency-key", "key-1", "extra"}
		}},
		{"input first", "--input", "ABSENT", func(dir string) []string {
			return []string{"--input", filepath.Join(dir, "absent-input"), "--actor", "ACT-1", "--idempotency-key", "key-1", "extra"}
		}},
		{"actor first", "--actor", "ACT-1", func(string) []string { return []string{"--actor", "ACT-1", "--idempotency-key", "key-1", "extra"} }},
		{"key first", "--idempotency-key", "key-1", func(string) []string { return []string{"--idempotency-key", "key-1", "extra"} }},
		{"positional", "arguments", "extra", func(string) []string { return []string{"extra"} }},
	}
	if got := len(leaves) * len(conflicts); got != 60 {
		t.Fatalf("exclusivity corpus has %d rows, want 60", got)
	}
	binary := buildAdoptedFID(t)
	for _, leaf := range leaves {
		t.Run(leaf, func(t *testing.T) {
			parts := strings.Split(leaf, " ")
			operation := inputDescriptorRegistry[leaf].Operation
			for _, conflict := range conflicts {
				t.Run(conflict.name, func(t *testing.T) {
					workspace := t.TempDir()
					args := append(append([]string{}, parts...), "--example")
					args = append(args, conflict.args(workspace)...)
					stdout, stderr, exit := runAdoptedFID(t, binary, workspace, args...)
					if exit != 2 || stdout != "" {
						t.Fatalf("%s %s: exit=%d stdout=%q stderr=%q", leaf, conflict.name, exit, stdout, stderr)
					}
					value := conflict.value
					if value == "ABSENT" {
						if conflict.argument == "--project" {
							value = filepath.Join(workspace, "absent-project")
						} else {
							value = filepath.Join(workspace, "absent-input")
						}
					}
					if conflict.argument == "--json" {
						assertAdoptedFailure(t, stderr, operation, conflict.argument, value)
					} else if stderr != "error: --example cannot be combined with "+conflict.argument+"\n" {
						t.Fatalf("first conflict = %q, want %s", stderr, conflict.argument)
					}
					// Human output predates typed details; exercise its shared detail writer in JSON mode.
					root := newRootCommandAt(workspace)
					var details bytes.Buffer
					root.SetErr(&details)
					if err := root.PersistentFlags().Set("json", "true"); err != nil {
						t.Fatal(err)
					}
					command, _, err := root.Find(parts)
					if err != nil {
						t.Fatal(err)
					}
					if err := writeExampleConflict(command, operation, conflict.argument, value); cli.ExitCode(err) != 2 {
						t.Fatalf("detail writer exit = %d", cli.ExitCode(err))
					}
					assertAdoptedFailure(t, details.String(), operation, conflict.argument, value)
					entries, err := os.ReadDir(workspace)
					if err != nil || len(entries) != 0 {
						t.Fatalf("conflict touched non-project workspace: entries=%v err=%v", entries, err)
					}
				})
			}
			workspace := t.TempDir()
			stdout, stderr, exit := runAdoptedFID(t, binary, workspace, append(parts, "--example")...)
			if exit != 0 || stderr != "" || !json.Valid([]byte(stdout)) {
				t.Fatalf("standalone %s: exit=%d stdout=%q stderr=%q", leaf, exit, stdout, stderr)
			}
			entries, err := os.ReadDir(workspace)
			if err != nil || len(entries) != 0 {
				t.Fatalf("standalone touched non-project workspace: entries=%v err=%v", entries, err)
			}
		})
	}
}

func buildAdoptedFID(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "fid")
	command := exec.Command("go", "build", "-o", binary, ".")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("build fid: %v: %s", err, output)
	}
	return binary
}

func runAdoptedFID(t *testing.T, binary, workspace string, args ...string) (string, string, int) {
	t.Helper()
	command := exec.Command(binary, args...)
	command.Dir = workspace
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	err := command.Run()
	if err == nil {
		return stdout.String(), stderr.String(), 0
	}
	if exit, ok := err.(*exec.ExitError); ok {
		return stdout.String(), stderr.String(), exit.ExitCode()
	}
	t.Fatalf("run fid: %v", err)
	return "", "", -1
}

func assertAdoptedFailure(t *testing.T, stderr, operation, argument string, value any) {
	assertAdoptedJSONFailure(t, stderr, operation, argument, value, []string{"--example"})
}

func assertAdoptedJSONFailure(t *testing.T, stderr, operation, argument string, value any, allowed []string) {
	t.Helper()
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(stderr), &fields); err != nil {
		t.Fatalf("decode failure %q: %v", stderr, err)
	}
	if _, present := fields["event"]; present {
		t.Fatalf("failure emitted Event: %s", stderr)
	}
	var envelope struct {
		ContractVersion string           `json:"contract_version"`
		OK              bool             `json:"ok"`
		Operation       string           `json:"operation"`
		Event           *json.RawMessage `json:"event"`
		Error           struct {
			Category string         `json:"category"`
			Code     string         `json:"code"`
			Details  map[string]any `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(stderr), &envelope); err != nil {
		t.Fatalf("decode failure %q: %v", stderr, err)
	}
	wantAllowed := make([]any, len(allowed))
	for i, item := range allowed {
		wantAllowed[i] = item
	}
	want := map[string]any{"argument": argument, "value": value, "allowed": wantAllowed}
	if envelope.ContractVersion != "fid-cli.v1" || envelope.OK || envelope.Operation != operation || envelope.Event != nil || envelope.Error.Category != "invalid-input" || envelope.Error.Code != "invalid-argument" || !reflect.DeepEqual(envelope.Error.Details, want) {
		t.Fatalf("failure = %s; want operation %s details %#v", stderr, operation, want)
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

func TestP24JudgmentDescriptorSourcePositionsAndExamples(t *testing.T) {
	t.Parallel()
	wantSources := map[string]string{
		"decision record-reference":    "input.source",
		"requirement record-reference": "input.source",
		"human-outcome record":         "input.references[]",
		"concern record":               "input.references[]",
		"deviation record":             "input.references[]",
		"requirement link":             "",
		"concern dispose":              "",
	}
	var positions func([]cli.InputField, string) []string
	positions = func(fields []cli.InputField, path string) []string {
		found := []string{}
		for _, field := range fields {
			current := path + "." + field.Name
			if field.Name == "source_type" {
				found = append(found, path)
			}
			found = append(found, positions(field.Children, current)...)
			if field.Element != nil {
				found = append(found, positions(field.Element.Children, current+"[]")...)
			}
		}
		return found
	}
	for leaf, want := range wantSources {
		leaf, want := leaf, want
		t.Run(leaf, func(t *testing.T) {
			t.Parallel()
			descriptor, exists := inputDescriptorRegistry[leaf]
			if !exists || pendingInputDescriptors[leaf] {
				t.Fatalf("%s lacks adopted descriptor", leaf)
			}
			found := positions(descriptor.Fields, "input")
			if want == "" && len(found) != 0 || want != "" && (len(found) != 1 || found[0] != want) {
				t.Fatalf("%s SourceInput positions = %v, want %q", leaf, found, want)
			}
			workspace := t.TempDir()
			parts := strings.Split(leaf, " ")
			help := executeCommandForTest(workspace, append(parts, "--help")...)
			if help.err != nil || !strings.Contains(help.stdout, "fid "+leaf+" --example") {
				t.Fatalf("help for %s = %#v", leaf, help)
			}
			var checkHelp func([]cli.InputField, string)
			checkHelp = func(fields []cli.InputField, prefix string) {
				for _, field := range fields {
					path := prefix + "." + field.Name
					if !strings.Contains(help.stdout, path+": ") || !strings.Contains(help.stdout, field.Meaning) {
						t.Errorf("help for %s lacks %s or its meaning", leaf, path)
					}
					checkHelp(field.Children, path)
					if field.Element != nil {
						checkHelp(field.Element.Children, path+"[i]")
					}
				}
			}
			checkHelp(descriptor.Fields, "input")
			example := executeCommandForTest(workspace, append(parts, "--example")...)
			if example.err != nil || example.stderr != "" || !json.Valid([]byte(example.stdout)) {
				t.Fatalf("example for %s = %#v", leaf, example)
			}
			template, err := cli.InputExample(descriptor)
			if err != nil || example.stdout != string(template) {
				t.Fatalf("example for %s differs from ordered template: %v", leaf, err)
			}
			if err := cli.CheckInputType(descriptor); err != nil {
				t.Fatalf("descriptor type for %s: %v", leaf, err)
			}
			var typed any
			if err := json.Unmarshal([]byte(example.stdout), &typed); err != nil {
				t.Fatal(err)
			}
			if err := cli.ValidateInput([]byte(example.stdout), descriptor); err != nil {
				t.Fatalf("template for %s: %v", leaf, err)
			}
			entries, err := os.ReadDir(workspace)
			if err != nil || len(entries) != 0 {
				t.Fatalf("example changed non-project directory: %v, %v", entries, err)
			}
		})
	}
}

func TestP24JudgmentStructuralValidationCorpus(t *testing.T) {
	t.Parallel()
	leaves := []string{"decision record-reference", "requirement record-reference", "requirement link",
		"human-outcome record", "concern record", "concern dispose", "deviation record"}
	for _, leaf := range leaves {
		t.Run(leaf, func(t *testing.T) {
			descriptor := inputDescriptorRegistry[leaf]
			template, err := cli.InputExample(descriptor)
			if err != nil {
				t.Fatal(err)
			}
			first := descriptor.Fields[0].Name
			var baseline map[string]any
			if err := json.Unmarshal(template, &baseline); err != nil {
				t.Fatal(err)
			}
			sourcePath := ""
			if _, ok := baseline["source"]; ok {
				sourcePath = "input.source"
			} else if _, ok := baseline["references"]; ok {
				sourcePath = "input.references[0]"
				baseline["references"] = []any{map[string]any{"source_type": "project-file", "display_label": "Example",
					"path": nil, "section": nil, "external_id": nil}}
			}
			fresh := func() map[string]any {
				t.Helper()
				encoded, err := json.Marshal(baseline)
				if err != nil {
					t.Fatal(err)
				}
				var copy map[string]any
				if err := json.Unmarshal(encoded, &copy); err != nil {
					t.Fatal(err)
				}
				return copy
			}
			check := func(name string, payload []byte, want string) {
				t.Helper()
				t.Run(name, func(t *testing.T) {
					invalid, ok := cli.ValidateInput(payload, descriptor).(*cli.InputError)
					if !ok || invalid.Path != want || len(invalid.Allowed) == 0 {
						t.Fatalf("%s: invalid = %#v, want %s", name, invalid, want)
					}
					workspace := t.TempDir()
					inputPath := filepath.Join(workspace, "request.json")
					if err := os.WriteFile(inputPath, payload, 0o600); err != nil {
						t.Fatal(err)
					}
					parts := append([]string{"--json"}, strings.Split(leaf, " ")...)
					parts = append(parts, "--input", inputPath, "--actor", "ACT-1", "--idempotency-key", "key-1")
					result := executeCommandForTest(workspace, parts...)
					if cli.ExitCode(result.err) != 2 || result.stdout != "" {
						t.Fatalf("%s: result = %#v", name, result)
					}
					var envelope struct {
						Error struct {
							Details struct {
								Argument string   `json:"argument"`
								Allowed  []string `json:"allowed"`
							} `json:"details"`
						} `json:"error"`
					}
					if err := json.Unmarshal([]byte(result.stderr), &envelope); err != nil || envelope.Error.Details.Argument != want || len(envelope.Error.Details.Allowed) == 0 {
						t.Fatalf("%s: stderr = %s, err = %v", name, result.stderr, err)
					}
					entries, err := os.ReadDir(workspace)
					if err != nil || len(entries) != 1 || entries[0].Name() != "request.json" {
						t.Fatalf("%s: unexpected effects: %v, %v", name, entries, err)
					}
				})
			}
			mutate := func(name, want string, change func(map[string]any)) {
				t.Helper()
				value := fresh()
				change(value)
				payload, err := json.Marshal(value)
				if err != nil {
					t.Fatal(err)
				}
				check(name, payload, want)
			}
			mutate("unknown", "input.bogus", func(v map[string]any) { v["bogus"] = true })
			mutate("missing", "input."+first, func(v map[string]any) { delete(v, first) })
			mutate("null", "input."+first, func(v map[string]any) { v[first] = nil })
			mutate("type", "input."+first, func(v map[string]any) { v[first] = 1 })
			mutate("format", "input."+first, func(v map[string]any) { v[first] = "" })
			encoded, err := json.Marshal(fresh())
			if err != nil {
				t.Fatal(err)
			}
			needle := `"` + first + `":`
			check("duplicate", []byte(strings.Replace(string(encoded), needle, needle+`"duplicate",`+needle, 1)), "input."+first)
			for _, field := range descriptor.Fields {
				if len(field.Members) > 0 {
					mutate("enum", "input."+field.Name, func(v map[string]any) { v[field.Name] = "not-a-member" })
					break
				}
			}
			if sourcePath == "" {
				return
			}
			getSource := func(v map[string]any) map[string]any {
				if sourcePath == "input.source" {
					return v["source"].(map[string]any)
				}
				return v["references"].([]any)[0].(map[string]any)
			}
			mutate("nested unknown", sourcePath+".bogus", func(v map[string]any) { getSource(v)["bogus"] = true })
			mutate("nested missing", sourcePath+".source_type", func(v map[string]any) { delete(getSource(v), "source_type") })
			mutate("nested null", sourcePath+".source_type", func(v map[string]any) { getSource(v)["source_type"] = nil })
			mutate("nested type", sourcePath+".source_type", func(v map[string]any) { getSource(v)["source_type"] = 1 })
			mutate("nested enum", sourcePath+".source_type", func(v map[string]any) { getSource(v)["source_type"] = "unknown" })
			mutate("nested format", sourcePath+".display_label", func(v map[string]any) { getSource(v)["display_label"] = "" })
			check("nested duplicate", []byte(strings.Replace(string(encoded), `"source_type":"project-file"`,
				`"source_type":"project-file","source_type":"external-reference"`, 1)), sourcePath+".source_type")
		})
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

func TestAdoptedStructuralProcessDetails(t *testing.T) {
	t.Parallel()
	baseline := map[string]string{
		"stop resolve":      `{"stop_id":"STP-1","fact_event_id":"EVT-1"}`,
		"attention resolve": `{"attention_id":"ATT-1","expected_version":1,"fact_event_ids":["EVT-1"]}`,
		"stop raise":        `{"subject":{"type":"run","id":"RUN-1"},"code":"operator-stop","blocked_actions":["run.start"],"details":{"reason":"reason","reference":"REF"},"resolution_condition":{"type":"work-ended","subject":{"type":"run","id":"RUN-1"},"budget_name":null,"capability_id":null},"resolver_capability":"stop.resolve"}`,
	}
	topResolve := []string{"fact_event_id", "stop_id"}
	topAttention := []string{"attention_id", "expected_version", "fact_event_ids"}
	topRaise := []string{"blocked_actions", "code", "details", "resolution_condition", "resolver_capability", "subject"}
	subject := []string{"id", "type"}
	codeMembers := []string{"executor-liveness-unknown", "recovery-state-conflict", "budget-exhausted", "budget-data-unavailable", "review-rework-exhausted", "protected-decision-required", "human-outcome-required", "activation-incomplete", "capability-unavailable", "operator-stop"}
	conditionMembers := []string{"executor-terminated", "run-recovered", "activation-recovered", "capability-available", "budget-observed", "work-ended", "decision-satisfied", "human-outcome-recorded", "review-recorded", "finding-resolved"}
	cases := []struct {
		name, leaf, old, replacement, path string
		value                              any
		allowed                            []string
	}{
		{"resolve unknown", "stop resolve", `"fact_event_id":"EVT-1"`, `"fact_event_id":"EVT-1","bogus":1`, "input.bogus", "bogus", topResolve},
		{"resolve duplicate", "stop resolve", `"stop_id":"STP-1"`, `"stop_id":"STP-1","stop_id":"STP-2"`, "input.stop_id", "stop_id", topResolve},
		{"resolve missing", "stop resolve", `"stop_id":"STP-1",`, ``, "input.stop_id", nil, []string{"stop_id"}},
		{"resolve null", "stop resolve", `"stop_id":"STP-1"`, `"stop_id":null`, "input.stop_id", "null", []string{"string"}},
		{"resolve type", "stop resolve", `"stop_id":"STP-1"`, `"stop_id":1`, "input.stop_id", "integer", []string{"string"}},
		{"resolve format", "stop resolve", `"stop_id":"STP-1"`, `"stop_id":""`, "input.stop_id", "", []string{"ID"}},
		{"attention unknown", "attention resolve", `"expected_version":1`, `"expected_version":1,"bogus":1`, "input.bogus", "bogus", topAttention},
		{"attention duplicate", "attention resolve", `"attention_id":"ATT-1"`, `"attention_id":"ATT-1","attention_id":"ATT-2"`, "input.attention_id", "attention_id", topAttention},
		{"attention missing", "attention resolve", `"attention_id":"ATT-1",`, ``, "input.attention_id", nil, []string{"attention_id"}},
		{"attention null", "attention resolve", `"attention_id":"ATT-1"`, `"attention_id":null`, "input.attention_id", "null", []string{"string"}},
		{"attention type", "attention resolve", `"expected_version":1`, `"expected_version":"1"`, "input.expected_version", "string", []string{"integer"}},
		{"attention format", "attention resolve", `"attention_id":"ATT-1"`, `"attention_id":""`, "input.attention_id", "", []string{"ID"}},
		{"attention bound", "attention resolve", `"expected_version":1`, `"expected_version":0`, "input.expected_version", float64(0), []string{"integer range 1..unbounded"}},
		{"attention array member", "attention resolve", `"fact_event_ids":["EVT-1"]`, `"fact_event_ids":[1]`, "input.fact_event_ids[0]", "integer", []string{"string"}},
		{"attention array null", "attention resolve", `"fact_event_ids":["EVT-1"]`, `"fact_event_ids":[null]`, "input.fact_event_ids[0]", "null", []string{"string"}},
		{"attention array format", "attention resolve", `"fact_event_ids":["EVT-1"]`, `"fact_event_ids":[""]`, "input.fact_event_ids[0]", "", []string{"ID"}},
		{"attention sorted unique", "attention resolve", `"fact_event_ids":["EVT-1"]`, `"fact_event_ids":["EVT-1","EVT-1"]`, "input.fact_event_ids", []json.RawMessage{json.RawMessage(`"EVT-1"`), json.RawMessage(`"EVT-1"`)}, []string{"sorted-unique"}},
		{"raise unknown", "stop raise", `"code":"operator-stop"`, `"code":"operator-stop","bogus":1`, "input.bogus", "bogus", topRaise},
		{"raise duplicate", "stop raise", `"code":"operator-stop"`, `"code":"operator-stop","code":"operator-stop"`, "input.code", "code", topRaise},
		{"raise missing", "stop raise", `"code":"operator-stop",`, ``, "input.code", nil, []string{"code"}},
		{"raise null", "stop raise", `"code":"operator-stop"`, `"code":null`, "input.code", "null", []string{"string"}},
		{"raise type", "stop raise", `"code":"operator-stop"`, `"code":1`, "input.code", "integer", []string{"string"}},
		{"raise enum", "stop raise", `"code":"operator-stop"`, `"code":"bogus"`, "input.code", "bogus", codeMembers},
		{"raise format", "stop raise", `"resolver_capability":"stop.resolve"`, `"resolver_capability":""`, "input.resolver_capability", "", []string{"ID"}},
		{"subject unknown", "stop raise", `"id":"RUN-1"},"code"`, `"id":"RUN-1","bogus":1},"code"`, "input.subject.bogus", "bogus", subject},
		{"subject duplicate", "stop raise", `"id":"RUN-1"},"code"`, `"id":"RUN-1","id":"RUN-2"},"code"`, "input.subject.id", "id", subject},
		{"subject missing", "stop raise", `"type":"run","id":"RUN-1"},"code"`, `"type":"run"},"code"`, "input.subject.id", nil, []string{"id"}},
		{"subject null", "stop raise", `"id":"RUN-1"},"code"`, `"id":null},"code"`, "input.subject.id", "null", []string{"string"}},
		{"subject type", "stop raise", `"id":"RUN-1"},"code"`, `"id":1},"code"`, "input.subject.id", "integer", []string{"string"}},
		{"subject enum", "stop raise", `"type":"run","id":"RUN-1"},"code"`, `"type":"bogus","id":"RUN-1"},"code"`, "input.subject.type", "bogus", []string{"run", "envelope", "activation"}},
		{"subject format", "stop raise", `"id":"RUN-1"},"code"`, `"id":""},"code"`, "input.subject.id", "", []string{"ID"}},
		{"deep unknown", "stop raise", `"id":"RUN-1"},"budget_name"`, `"id":"RUN-1","bogus":1},"budget_name"`, "input.resolution_condition.subject.bogus", "bogus", subject},
		{"deep duplicate", "stop raise", `"id":"RUN-1"},"budget_name"`, `"id":"RUN-1","id":"RUN-2"},"budget_name"`, "input.resolution_condition.subject.id", "id", subject},
		{"deep missing", "stop raise", `"type":"run","id":"RUN-1"},"budget_name"`, `"type":"run"},"budget_name"`, "input.resolution_condition.subject.id", nil, []string{"id"}},
		{"deep null", "stop raise", `"id":"RUN-1"},"budget_name"`, `"id":null},"budget_name"`, "input.resolution_condition.subject.id", "null", []string{"string"}},
		{"deep type", "stop raise", `"id":"RUN-1"},"budget_name"`, `"id":1},"budget_name"`, "input.resolution_condition.subject.id", "integer", []string{"string"}},
		{"deep enum", "stop raise", `"type":"run","id":"RUN-1"},"budget_name"`, `"type":"bogus","id":"RUN-1"},"budget_name"`, "input.resolution_condition.subject.type", "bogus", []string{"run", "envelope", "activation"}},
		{"deep format", "stop raise", `"id":"RUN-1"},"budget_name"`, `"id":""},"budget_name"`, "input.resolution_condition.subject.id", "", []string{"ID"}},
		{"condition enum", "stop raise", `"type":"work-ended","subject"`, `"type":"bogus","subject"`, "input.resolution_condition.type", "bogus", conditionMembers},
		{"details unknown", "stop raise", `"reference":"REF"},"resolution_condition"`, `"reference":"REF","bogus":1},"resolution_condition"`, "input.details.bogus", "bogus", []string{"reason", "reference"}},
		{"action member type", "stop raise", `"blocked_actions":["run.start"]`, `"blocked_actions":[1]`, "input.blocked_actions[0]", "integer", []string{"string"}},
		{"action member enum", "stop raise", `"blocked_actions":["run.start"]`, `"blocked_actions":["bogus"]`, "input.blocked_actions[0]", "bogus", []string{"claim.acquire", "run.start", "run.recover", "gate.run", "issue.transition-terminal-success", "envelope.complete"}},
	}
	if len(cases) != 42 {
		t.Fatalf("structural corpus has %d rows, want 42", len(cases))
	}
	binary := buildAdoptedFID(t)
	for leaf, payload := range baseline {
		t.Run(leaf+" valid structure", func(t *testing.T) {
			if err := cli.ValidateInput([]byte(payload), inputDescriptorRegistry[leaf]); err != nil {
				t.Fatalf("qualifying structure: %v", err)
			}
			workspace := t.TempDir()
			input := filepath.Join(t.TempDir(), "valid.json")
			if err := os.WriteFile(input, []byte(payload), 0o600); err != nil {
				t.Fatal(err)
			}
			args := append([]string{"--json"}, strings.Split(leaf, " ")...)
			args = append(args, "--input", input, "--actor", "ACT-1", "--idempotency-key", "valid-structure")
			_, stderr, exit := runAdoptedFID(t, binary, workspace, args...)
			if exit != 2 || !strings.Contains(stderr, `"argument":"input"`) || !strings.Contains(stderr, `"allowed":["accepted `+inputDescriptorRegistry[leaf].Operation+` input"]`) {
				t.Fatalf("valid structure did not reach domain validation: exit=%d stderr=%q", exit, stderr)
			}
		})
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			payload := strings.Replace(baseline[test.leaf], test.old, test.replacement, 1)
			if payload == baseline[test.leaf] {
				t.Fatal("mutation did not match baseline")
			}
			invalid, ok := cli.ValidateInput([]byte(payload), inputDescriptorRegistry[test.leaf]).(*cli.InputError)
			if !ok || invalid.Path != test.path || !reflect.DeepEqual(invalid.Value, test.value) || !reflect.DeepEqual(invalid.Allowed, test.allowed) {
				t.Fatalf("typed details = %#v, want path=%s value=%#v allowed=%v", invalid, test.path, test.value, test.allowed)
			}
			workspace := t.TempDir()
			input := filepath.Join(t.TempDir(), "input.json")
			if err := os.WriteFile(input, []byte(payload), 0o600); err != nil {
				t.Fatal(err)
			}
			args := append([]string{"--json"}, strings.Split(test.leaf, " ")...)
			args = append(args, "--input", input, "--actor", "ACT-1", "--idempotency-key", "structure-1")
			stdout, stderr, exit := runAdoptedFID(t, binary, workspace, args...)
			if exit != 2 || stdout != "" {
				t.Fatalf("exit=%d stdout=%q stderr=%q", exit, stdout, stderr)
			}
			wireValue := test.value
			if wireValue != nil {
				if _, ok := wireValue.(string); !ok {
					encoded, err := json.Marshal(wireValue)
					if err != nil {
						t.Fatal(err)
					}
					wireValue = string(encoded)
				}
			}
			assertAdoptedJSONFailure(t, stderr, inputDescriptorRegistry[test.leaf].Operation, test.path, wireValue, test.allowed)
			entries, err := os.ReadDir(workspace)
			if err != nil || len(entries) != 0 {
				t.Fatalf("failure touched workspace: entries=%v err=%v", entries, err)
			}
		})
	}
}
