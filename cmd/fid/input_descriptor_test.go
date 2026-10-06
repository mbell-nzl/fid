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
	"slices"
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
		"requirement record-reference", "requirement link", "human-outcome record", "concern record", "concern dispose", "deviation record", "review record", "finding address", "gate register", "envelope create", "envelope activate", "envelope complete", "claim acquire", "claim renew", "claim release", "claim expire", "run start", "run close", "run abandon", "run observe-budget", "run observe", "run checkpoint", "run recover", "gate run", "evidence link", "evidence evaluate", "artifact verify"}
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
	if got := len(leaves) * len(conflicts); got != 186 {
		t.Fatalf("exclusivity corpus has %d rows, want 186", got)
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

func TestP24GateEnvelopeClaimDescriptors(t *testing.T) {
	t.Parallel()
	for _, leaf := range []string{"gate register", "envelope create", "envelope activate", "envelope complete", "claim acquire", "claim renew", "claim release", "claim expire"} {
		descriptor, ok := inputDescriptorRegistry[leaf]
		if !ok || pendingInputDescriptors[leaf] {
			t.Errorf("%s has no adopted descriptor", leaf)
			continue
		}
		if err := cli.CheckInputType(descriptor); err != nil {
			t.Errorf("%s descriptor/type: %v", leaf, err)
		}
	}
}

func TestP24RemovedBindingMembers(t *testing.T) {
	t.Parallel()
	for leaf, members := range map[string][]string{"gate register": {"repository", "worktree_id"}, "envelope create": {"repository_scope"}, "envelope complete": {"repository"}, "claim acquire": {"worktree_id"}} {
		for _, member := range members {
			t.Run(leaf+"/"+member, func(t *testing.T) {
				t.Parallel()
				descriptor := inputDescriptorRegistry[leaf]
				data, err := cli.InputExample(descriptor)
				if err != nil {
					t.Fatal(err)
				}
				var input map[string]any
				if err := json.Unmarshal(data, &input); err != nil {
					t.Fatal(err)
				}
				p24ExpandDescriptorObjects(t, input, descriptor.Fields)
				input[member] = nil
				result := p24BindingMutation(t, t.TempDir(), leaf, input, "removed")
				if cli.ExitCode(result.err) != 2 || result.stdout != "" || !strings.Contains(result.stderr, `"argument":"input.`+member+`"`) {
					t.Fatalf("removed member: %#v", result)
				}
			})
		}
	}
}

func TestP24BindingDescriptorProcessCorpus(t *testing.T) {
	t.Parallel()
	leaves := bindingDescriptors()
	if len(leaves) != 8 {
		t.Fatalf("binding leaf corpus = %d, want 8", len(leaves))
	}
	p24DescriptorProcessCorpus(t, leaves)
}

func TestP24RunDescriptorProcessCorpus(t *testing.T) {
	t.Parallel()
	leaves := runDescriptors()
	if len(leaves) != 7 {
		t.Fatalf("Run corpus = %d, want 7", len(leaves))
	}
	p24DescriptorProcessCorpus(t, leaves)
}

func p24DescriptorProcessCorpus(t *testing.T, leaves map[string]cli.InputDescriptor) {
	t.Helper()
	binary := buildAdoptedFID(t)
	for leaf, descriptor := range leaves {
		t.Run(leaf, func(t *testing.T) {
			t.Parallel()
			if pendingInputDescriptors[leaf] || cli.CheckInputType(descriptor) != nil {
				t.Fatal("descriptor missing or mismatched")
			}
			workspace := t.TempDir()
			parts := strings.Split(leaf, " ")
			help := executeCommandForTest(workspace, append(parts, "--help")...)
			if help.err != nil || !strings.Contains(help.stdout, "fid "+leaf+" --example") {
				t.Fatalf("help: %#v", help)
			}
			var helpFields func([]cli.InputField, string)
			helpFields = func(fields []cli.InputField, prefix string) {
				for _, field := range fields {
					path := prefix + "." + field.Name
					typeName := field.Type
					if field.Element != nil {
						typeName += "<" + field.Element.Type + ">"
					}
					if !strings.Contains(help.stdout, path+": "+typeName+", "+field.Requirement) || !strings.Contains(help.stdout, field.Meaning) {
						t.Errorf("help lacks %s type, requirement or meaning", path)
					}
					line := ""
					for _, candidate := range strings.Split(help.stdout, "\n") {
						if strings.HasPrefix(strings.TrimSpace(candidate), path+": ") {
							line = candidate
						}
					}
					constraints := []string{}
					if field.Format != "" {
						constraints = append(constraints, field.Format)
					}
					if len(field.Members) > 0 {
						constraints = append(constraints, "allowed: "+strings.Join(field.Members, "|"))
					}
					if field.Element != nil {
						if len(field.Element.Members) > 0 {
							constraints = append(constraints, "elements allowed: "+strings.Join(field.Element.Members, "|"))
						} else if field.Element.Format != "" {
							constraints = append(constraints, "element format: "+field.Element.Format)
						}
					}
					minimum, maximum := "unbounded", "unbounded"
					if field.MinInt != nil {
						minimum = fmt.Sprint(*field.MinInt)
					}
					if field.MaxInt != nil {
						maximum = fmt.Sprint(*field.MaxInt)
					}
					if field.MinInt != nil || field.MaxInt != nil {
						constraints = append(constraints, "integer range "+minimum+".."+maximum)
					}
					minimum, maximum = "0", "unbounded"
					if field.MinLen != nil {
						minimum = fmt.Sprint(*field.MinLen)
					}
					if field.MaxLen != nil {
						maximum = fmt.Sprint(*field.MaxLen)
					}
					if field.MinLen != nil || field.MaxLen != nil {
						constraints = append(constraints, "length "+minimum+".."+maximum)
					}
					for _, constraint := range constraints {
						if !strings.Contains(line, constraint) {
							t.Errorf("help lacks %s constraint %s", path, constraint)
						}
					}
					helpFields(field.Children, path)
					if field.Element != nil {
						helpFields(field.Element.Children, path+"[i]")
					}
				}
			}
			helpFields(descriptor.Fields, "input")
			template, err := cli.InputExample(descriptor)
			if err != nil {
				t.Fatal(err)
			}
			stdout, stderr, exit := runAdoptedFID(t, binary, workspace, append(parts, "--example")...)
			if exit != 0 || stderr != "" || stdout != string(template) {
				t.Fatalf("template: exit=%d stderr=%q stdout=%q", exit, stderr, stdout)
			}
			request := reflect.New(descriptor.Type).Interface()
			decoder := json.NewDecoder(strings.NewReader(stdout))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(request); err != nil {
				t.Fatal(err)
			}

			entries, err := os.ReadDir(workspace)
			if err != nil || len(entries) != 0 {
				t.Fatalf("example side effects: %v / %v", entries, err)
			}
			var baseline map[string]any
			if err := json.Unmarshal(template, &baseline); err != nil {
				t.Fatal(err)
			}
			p24ExpandDescriptorObjects(t, baseline, descriptor.Fields)
			payload, err := json.Marshal(baseline)
			if err != nil {
				t.Fatal(err)
			}
			if err := cli.ValidateInput(payload, descriptor); err != nil {
				t.Fatalf("expanded qualifying structure: %v", err)
			}
			check := func(name, path string, object map[string]any, change func(), restore func()) {
				t.Helper()
				t.Run(name, func(t *testing.T) {
					change()
					data, err := json.Marshal(baseline)
					restore()
					if err != nil {
						t.Fatal(err)
					}
					if strings.HasSuffix(name, "/duplicate") {
						field := path[strings.LastIndex(path, ".")+1:]
						value, err := json.Marshal(object[field])
						if err != nil {
							t.Fatal(err)
						}
						object["__p24_duplicate_target"] = "marker"
						marked, err := json.Marshal(baseline)
						delete(object, "__p24_duplicate_target")
						if err != nil {
							t.Fatal(err)
						}
						data = []byte(strings.Replace(string(marked), `"__p24_duplicate_target":"marker",`, `"`+field+`":`+string(value)+`,`, 1))
					}
					invalid, ok := cli.ValidateInput(data, descriptor).(*cli.InputError)
					if !ok || invalid.Path != path || len(invalid.Allowed) == 0 {
						t.Fatalf("structural path: %#v, want %s", invalid, path)
					}
					input := filepath.Join(t.TempDir(), "input.json")
					if err := os.WriteFile(input, data, 0o600); err != nil {
						t.Fatal(err)
					}
					args := append([]string{"--json"}, parts...)
					args = append(args, "--input", input)
					if leaf != "evidence evaluate" {
						args = append(args, "--actor", "ACT-1", "--idempotency-key", "structure")
					}
					out, failure, exit := runAdoptedFID(t, binary, workspace, args...)
					if exit != 2 || out != "" {
						t.Fatalf("structural refusal: exit=%d stdout=%q stderr=%q", exit, out, failure)
					}
					value := invalid.Value
					if value != nil {
						if _, ok := value.(string); !ok {
							encoded, err := json.Marshal(value)
							if err != nil {
								t.Fatal(err)
							}
							value = string(encoded)
						}
					}
					assertAdoptedJSONFailure(t, failure, descriptor.Operation, path, value, invalid.Allowed)
					entries, err := os.ReadDir(workspace)
					if err != nil || len(entries) != 0 {
						t.Fatalf("structural refusal side effects: %v / %v", entries, err)
					}
				})
			}
			var walk func(map[string]any, []cli.InputField, string)
			walk = func(object map[string]any, fields []cli.InputField, prefix string) {
				check(prefix+"/unknown", prefix+".bogus", object, func() { object["bogus"] = true }, func() { delete(object, "bogus") })
				for _, field := range fields {
					path, original := prefix+"."+field.Name, object[field.Name]
					noop := func() {}
					check(prefix+"/"+field.Name+"/duplicate", path, object, noop, noop)
					if field.Requirement != "optional" {
						check(path+"/missing", path, object, func() { delete(object, field.Name) }, func() { object[field.Name] = original })
					}
					if field.Requirement == "required" {
						check(path+"/null", path, object, func() { object[field.Name] = nil }, func() { object[field.Name] = original })
					}
					wrong := any(true)
					if field.Type == "boolean" {
						wrong = "true"
					}
					check(path+"/type", path, object, func() { object[field.Name] = wrong }, func() { object[field.Name] = original })
					if len(field.Members) > 0 {
						check(path+"/enum", path, object, func() { object[field.Name] = "not-a-member" }, func() { object[field.Name] = original })
					}
					if field.Type == "string" && field.Format != "" && field.Format != "text" {
						check(path+"/format", path, object, func() { object[field.Name] = "" }, func() { object[field.Name] = original })
					}
					if field.MinInt != nil {
						check(path+"/bound", path, object, func() { object[field.Name] = *field.MinInt - 1 }, func() { object[field.Name] = original })
					}
					if field.Type == "string" && field.MinLen != nil && *field.MinLen > 0 {
						check(path+"/short", path, object, func() { object[field.Name] = strings.Repeat("x", *field.MinLen-1) }, func() { object[field.Name] = original })
					}
					if field.Type == "string" && field.MaxLen != nil {
						check(path+"/long", path, object, func() { object[field.Name] = strings.Repeat("x", *field.MaxLen+1) }, func() { object[field.Name] = original })
					}
					if field.Type == "array" {
						if field.MinLen != nil && *field.MinLen > 0 {
							check(path+"/short", path, object, func() { object[field.Name] = []any{} }, func() { object[field.Name] = original })
						}
						if field.MaxLen != nil {
							check(path+"/long", path, object, func() { object[field.Name] = make([]any, *field.MaxLen+1) }, func() { object[field.Name] = original })
						}
						if field.Element != nil {
							check(path+"/element-type", path+"[0]", object, func() { object[field.Name] = []any{true} }, func() { object[field.Name] = original })
							if field.Element.Type == "string" && field.Element.Format != "" && field.Element.Format != "text" {
								check(path+"/element-format", path+"[0]", object, func() { object[field.Name] = []any{""} }, func() { object[field.Name] = original })
							}
							if field.Format == "sorted-unique" {
								check(path+"/order", path, object, func() { object[field.Name] = []any{"duplicate", "duplicate"} }, func() { object[field.Name] = original })
							}
						}
					}
					if len(field.Children) > 0 {
						walk(original.(map[string]any), field.Children, path)
					}
					if field.Element != nil && len(field.Element.Children) > 0 {
						walk(original.([]any)[0].(map[string]any), field.Element.Children, path+"[0]")
					}
				}
			}
			walk(baseline, descriptor.Fields, "input")
		})
	}
}

func p24ExpandDescriptorObjects(t *testing.T, object map[string]any, fields []cli.InputField) {
	t.Helper()
	for _, field := range fields {
		if field.Requirement == "optional" {
			valueField := field
			valueField.Requirement = "required"
			data, err := cli.InputExample(cli.InputDescriptor{Fields: []cli.InputField{valueField}})
			if err != nil {
				t.Fatal(err)
			}
			var values map[string]any
			if err := json.Unmarshal(data, &values); err != nil {
				t.Fatal(err)
			}
			object[field.Name] = values[field.Name]
		}
		if field.Format == "digest" && object[field.Name] != nil {
			object[field.Name] = "sha256:" + strings.Repeat("1", 64)
		}
		if (field.Format == "oid" || field.Format == "OID") && object[field.Name] != nil {
			object[field.Name] = strings.Repeat("1", 40)
		}
		if field.Format == "timestamp" && object[field.Name] != nil {
			object[field.Name] = "2026-10-06T00:00:00Z"
		}
		children := field.Children
		if field.Element != nil {
			children = field.Element.Children
		}
		if len(children) == 0 {
			continue
		}
		data, err := cli.InputExample(cli.InputDescriptor{Fields: children})
		if err != nil {
			t.Fatal(err)
		}
		var child map[string]any
		if err := json.Unmarshal(data, &child); err != nil {
			t.Fatal(err)
		}
		p24ExpandDescriptorObjects(t, child, children)
		if field.Type == "array" {
			object[field.Name] = []any{child}
		} else {
			object[field.Name] = child
		}
	}
}

func p24ReviewHelpMatches(help, leaf string) bool {
	return p24JudgmentHelpMatches(help, inputDescriptorRegistry[leaf].Fields, "input") && strings.Contains(help, "fid "+leaf+" --example") && (leaf != "review record" ||
		strings.Contains(help, "git rev-parse HEAD") && strings.Contains(help, "git status --porcelain"))
}

func TestP24ReviewHelpEvidenceRejectsIncompleteMetadata(t *testing.T) {
	t.Parallel()
	for leaf, descriptor := range reviewDescriptors() {
		t.Run(leaf, func(t *testing.T) {
			t.Parallel()
			help := executeCommandForTest(t.TempDir(), append(strings.Split(leaf, " "), "--help")...).stdout
			if !p24ReviewHelpMatches(help, leaf) {
				t.Fatal("valid help rejected")
			}
			var corrupt func([]cli.InputField, string)
			corrupt = func(fields []cli.InputField, prefix string) {
				for _, field := range fields {
					path := prefix + "." + field.Name
					for _, metadata := range append(p24JudgmentHelpMetadata(field), field.Meaning) {
						lines := strings.Split(help, "\n")
						for i, line := range lines {
							if strings.HasPrefix(strings.TrimSpace(line), path+": ") {
								lines[i] = strings.Replace(line, metadata, "WRONG", 1)
								if lines[i] == line {
									t.Fatalf("seed did not change %s metadata %s", path, metadata)
								}
							}
						}
						if p24ReviewHelpMatches(strings.Join(lines, "\n"), leaf) {
							t.Errorf("help checker accepted wrong %s metadata %s", path, metadata)
						}
					}
					corrupt(field.Children, path)
					if field.Element != nil {
						corrupt(field.Element.Children, path+"[i]")
					}
				}
			}
			corrupt(descriptor.Fields, "input")
		})
	}
}

func TestP24UnknownInputFormatRejects(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"unknown-format", "Oid", "oID"} {
		t.Run(format, func(t *testing.T) {
			t.Parallel()
			descriptor := cli.InputDescriptor{Fields: []cli.InputField{{Name: "head_oid", Type: "string", Requirement: "required", Format: format}}}
			payload := []byte(`{"head_oid":"a0ca854133ae98b77e36cf51ad596ff0e8c6e096"}`)
			invalid, ok := cli.ValidateInput(payload, descriptor).(*cli.InputError)
			want := &cli.InputError{Path: "input.head_oid", Value: "a0ca854133ae98b77e36cf51ad596ff0e8c6e096", Allowed: []string{format}}
			if !ok || !p24JudgmentDetailsMatch(invalid, want) || !reflect.DeepEqual(invalid.Details(), want.Details()) {
				t.Fatalf("unknown format %q: got %#v, want %#v", format, invalid, want)
			}
		})
	}
}

func TestP24OIDFormatDispatch(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"OID", "oid"} {
		t.Run(format, func(t *testing.T) {
			t.Parallel()
			descriptor := cli.InputDescriptor{Fields: []cli.InputField{{Name: "head_oid", Type: "string", Requirement: "required", Format: format}}}
			for _, row := range []struct {
				name, value string
				valid       bool
			}{
				{"valid", "a0ca854133ae98b77e36cf51ad596ff0e8c6e096", true},
				{"malformed", "not-an-oid", false},
				{"empty", "", false},
				{"short", "a0ca854", false},
				{"nonhex", strings.Repeat("g", 40), false},
				{"uppercase hex", strings.Repeat("A", 40), false},
				{"alternate length", strings.Repeat("a", 64), false},
				{"placeholder", "<head-oid>", false},
			} {
				t.Run(row.name, func(t *testing.T) {
					t.Parallel()
					payload, err := json.Marshal(map[string]string{"head_oid": row.value})
					if err != nil {
						t.Fatal(err)
					}
					validation := cli.ValidateInput(payload, descriptor)
					if row.valid {
						if validation != nil {
							t.Fatalf("valid %s rejected: %v", format, validation)
						}
						return
					}
					invalid, ok := validation.(*cli.InputError)
					want := &cli.InputError{Path: "input.head_oid", Value: row.value, Allowed: []string{format}}
					if !ok || !p24JudgmentDetailsMatch(invalid, want) || !reflect.DeepEqual(invalid.Details(), want.Details()) {
						t.Fatalf("invalid %s: got %#v, want %#v", format, invalid, want)
					}
				})
			}
		})
	}
}

func TestP24ReviewMalformedOIDProcess(t *testing.T) {
	t.Parallel()
	descriptor := inputDescriptorRegistry["review record"]
	example, err := cli.InputExample(descriptor)
	if err != nil {
		t.Fatal(err)
	}
	binary := buildAdoptedFID(t)
	for _, value := range []string{"not-an-oid", "", "a0ca854", strings.Repeat("g", 40)} {
		t.Run(value, func(t *testing.T) {
			t.Parallel()
			workspace := t.TempDir()
			input := filepath.Join(t.TempDir(), "input.json")
			payload := bytes.ReplaceAll(example, []byte("<head-oid>"), []byte(value))
			if err := os.WriteFile(input, payload, 0o600); err != nil {
				t.Fatal(err)
			}
			stdout, stderr, exit := runAdoptedFID(t, binary, workspace, "--json", "review", "record", "--input", input, "--actor", "ACT-1", "--idempotency-key", "oid-process")
			if exit != 2 || stdout != "" {
				t.Fatalf("process exit=%d stdout=%q stderr=%q", exit, stdout, stderr)
			}
			assertAdoptedJSONFailure(t, stderr, "review.record", "input.expected_repository.head_oid", value, []string{"oid"})
			var envelope map[string]any
			if err := json.Unmarshal([]byte(stderr), &envelope); err != nil {
				t.Fatal(err)
			}
			failure := envelope["error"].(map[string]any)
			if len(envelope) != 4 || len(failure) != 6 || failure["retryable"] != false || failure["subject"] != nil || failure["message"] != "invalid input.expected_repository.head_oid" {
				t.Fatalf("closed structural envelope: %s", stderr)
			}
			entries, err := os.ReadDir(workspace)
			if err != nil || len(entries) != 0 {
				t.Fatalf("process changed workspace: %v / %v", entries, err)
			}
		})
	}
}

func TestP24ReviewAndFindingDescriptorHelpAndStructure(t *testing.T) {
	t.Parallel()
	binary := buildAdoptedFID(t)
	for _, leaf := range []string{"review record", "finding address"} {
		t.Run(leaf, func(t *testing.T) {
			descriptor := inputDescriptorRegistry[leaf]
			if pendingInputDescriptors[leaf] {
				t.Fatalf("%s remains behind the descriptor fence", leaf)
			}
			example, err := cli.InputExample(descriptor)
			if err != nil {
				t.Fatalf("template generation: %v", err)
			}
			// VAL-014 checks typed placeholders, not executable OID values.
			var template map[string]any
			if err := json.Unmarshal(example, &template); err != nil {
				t.Fatal(err)
			}
			p24ReviewTemplateObject(t, template, descriptor.Fields, "input")
			typed := reflect.New(descriptor.Type).Interface()
			decoder := json.NewDecoder(bytes.NewReader(example))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(typed); err != nil {
				t.Fatalf("template does not decode: %v", err)
			}
			parts := strings.Split(leaf, " ")
			workspace := t.TempDir()
			stdout, stderr, exit := runAdoptedFID(t, binary, workspace, append(parts, "--example")...)
			if exit != 0 || stderr != "" || stdout != string(example) {
				t.Fatalf("ordered example: exit=%d stdout=%q stderr=%q", exit, stdout, stderr)
			}
			entries, err := os.ReadDir(workspace)
			if err != nil || len(entries) != 0 {
				t.Fatalf("example changed non-project directory: %v / %v", entries, err)
			}
			if leaf == "review record" {
				if template["expected_repository"].(map[string]any)["head_oid"] != "<head-oid>" {
					t.Fatal("OID placeholder changed")
				}
				// Runtime controls use a real Git commit OID.
				example = bytes.ReplaceAll(example, []byte("<head-oid>"), []byte("a0ca854133ae98b77e36cf51ad596ff0e8c6e096"))
			}
			if err := cli.ValidateInput(example, descriptor); err != nil {
				t.Fatalf("qualifying runtime structure: %v", err)
			}
			help := executeCommandForTest(t.TempDir(), append(parts, "--help")...)
			if help.err != nil || !p24ReviewHelpMatches(help.stdout, leaf) {
				t.Fatalf("descriptor help = %#v", help)
			}
			if leaf == "review record" && (!strings.Contains(help.stdout, "git rev-parse HEAD") || !strings.Contains(help.stdout, "git status --porcelain")) {
				t.Fatalf("review help omits git commands: %s", help.stdout)
			}
			p24ReviewRecursiveStructuralDetails(t, leaf, descriptor)
			first := descriptor.Fields[0].Name
			cases := []struct {
				name, path string
				change     func(map[string]any)
			}{
				{"unknown", "input.bogus", func(v map[string]any) { v["bogus"] = true }},
				{"missing", "input." + first, func(v map[string]any) { delete(v, first) }},
				{"null", "input." + first, func(v map[string]any) { v[first] = nil }},
				{"type", "input." + first, func(v map[string]any) { v[first] = 1 }},
				{"format", "input." + first, func(v map[string]any) { v[first] = "" }},
			}
			if leaf == "review record" {
				cases = append(cases,
					struct {
						name, path string
						change     func(map[string]any)
					}{"nested unknown", "input.expected_repository.bogus", func(v map[string]any) { v["expected_repository"].(map[string]any)["bogus"] = true }},
					struct {
						name, path string
						change     func(map[string]any)
					}{"deep enum", "input.coverage.items[0].source.source_type", func(v map[string]any) {
						v["coverage"].(map[string]any)["items"].([]any)[0].(map[string]any)["source"].(map[string]any)["source_type"] = "bogus"
					}},
				)
			}
			for _, test := range cases {
				t.Run(test.name, func(t *testing.T) {
					var value map[string]any
					if err := json.Unmarshal(example, &value); err != nil {
						t.Fatal(err)
					}
					test.change(value)
					payload, err := json.Marshal(value)
					if err != nil {
						t.Fatal(err)
					}
					invalid, ok := cli.ValidateInput(payload, descriptor).(*cli.InputError)
					if !ok || invalid.Path != test.path {
						t.Fatalf("validation = %#v, want %s", invalid, test.path)
					}
					workspace := t.TempDir()
					input := filepath.Join(workspace, "input.json")
					if err := os.WriteFile(input, payload, 0o600); err != nil {
						t.Fatal(err)
					}
					args := append([]string{"--json"}, parts...)
					args = append(args, "--input", input, "--actor", "ACT-1", "--idempotency-key", "structure")
					result := executeCommandForTest(workspace, args...)
					if cli.ExitCode(result.err) != 2 || result.stdout != "" {
						t.Fatalf("structural failure = %#v", result)
					}
					assertAdoptedJSONFailure(t, result.stderr, descriptor.Operation, test.path, invalid.Value, invalid.Allowed)
					if _, err := os.Stat(filepath.Join(workspace, ".fid")); !os.IsNotExist(err) {
						t.Fatalf("structural failure changed workspace: %v", err)
					}
				})
			}
			duplicates := []struct{ name, key, path string }{{"top duplicate", first, "input." + first}}
			if leaf == "review record" {
				duplicates = append(duplicates, struct{ name, key, path string }{"deep duplicate", "head_oid", "input.expected_repository.head_oid"})
			}
			for _, test := range duplicates {
				t.Run(test.name, func(t *testing.T) {
					needle := `"` + test.key + `":`
					payload := []byte(strings.Replace(string(example), needle, needle+`"duplicate",`+needle, 1))
					invalid, ok := cli.ValidateInput(payload, descriptor).(*cli.InputError)
					if !ok || invalid.Path != test.path {
						t.Fatalf("duplicate validation = %#v, want %s", invalid, test.path)
					}
					workspace := t.TempDir()
					input := filepath.Join(workspace, "input.json")
					if err := os.WriteFile(input, payload, 0o600); err != nil {
						t.Fatal(err)
					}
					args := append([]string{"--json"}, parts...)
					args = append(args, "--input", input, "--actor", "ACT-1", "--idempotency-key", "duplicate")
					result := executeCommandForTest(workspace, args...)
					if cli.ExitCode(result.err) != 2 || result.stdout != "" {
						t.Fatalf("duplicate failure = %#v", result)
					}
					assertAdoptedJSONFailure(t, result.stderr, descriptor.Operation, test.path, invalid.Value, invalid.Allowed)
					if _, err := os.Stat(filepath.Join(workspace, ".fid")); !os.IsNotExist(err) {
						t.Fatalf("duplicate failure changed workspace: %v", err)
					}
				})
			}
		})
	}
}

func p24ReviewTemplateObject(t *testing.T, object map[string]any, fields []cli.InputField, prefix string) {
	t.Helper()
	want := 0
	for _, field := range fields {
		value, present := object[field.Name]
		if field.Requirement == "optional" {
			if present {
				t.Fatalf("optional template member %s.%s", prefix, field.Name)
			}
			continue
		}
		want++
		if !present {
			t.Fatalf("missing template member %s.%s", prefix, field.Name)
		}
		p24ReviewTemplateValue(t, value, field, prefix+"."+field.Name)
	}
	if len(object) != want {
		t.Fatalf("template object %s has %d keys, want %d", prefix, len(object), want)
	}
}

func p24ReviewTemplateValue(t *testing.T, value any, field cli.InputField, path string) {
	t.Helper()
	if value == nil && field.Requirement == "nullable" {
		return
	}
	valid := false
	switch field.Type {
	case "string":
		_, valid = value.(string)
	case "boolean":
		_, valid = value.(bool)
	case "integer":
		number, ok := value.(float64)
		valid = ok && number == float64(int64(number))
	case "object":
		object, ok := value.(map[string]any)
		valid = ok
		if ok {
			p24ReviewTemplateObject(t, object, field.Children, path)
		}
	case "array":
		array, ok := value.([]any)
		valid = ok
		if ok {
			if field.MinLen != nil && len(array) < *field.MinLen {
				t.Fatalf("template %s lacks required elements", path)
			}
			for i, entry := range array {
				p24ReviewTemplateValue(t, entry, *field.Element, fmt.Sprintf("%s[%d]", path, i))
			}
		}
	}
	if !valid {
		t.Fatalf("template %s type = %T, want %s", path, value, field.Type)
	}
}

func p24JudgmentHelpMatches(help string, fields []cli.InputField, prefix string) bool {
	for _, field := range fields {
		path := prefix + "." + field.Name
		want := "  " + path + ": " + strings.Join(p24JudgmentHelpMetadata(field), ", ") + " — " + field.Meaning
		found := 0
		for _, line := range strings.Split(help, "\n") {
			if line == want {
				found++
			}
		}
		if found != 1 {
			return false
		}
		if !p24JudgmentHelpMatches(help, field.Children, path) {
			return false
		}
		if field.Element != nil && !p24JudgmentHelpMatches(help, field.Element.Children, path+"[i]") {
			return false
		}
	}
	return true
}

func p24JudgmentHelpMetadata(field cli.InputField) []string {
	typeName := field.Type
	if field.Element != nil {
		typeName += "<" + field.Element.Type + ">"
	}
	metadata := []string{typeName, field.Requirement}
	if field.Format != "" {
		metadata = append(metadata, field.Format)
	}
	if len(field.Members) > 0 {
		metadata = append(metadata, "allowed: "+strings.Join(field.Members, "|"))
	}
	if field.Element != nil {
		if len(field.Element.Members) > 0 {
			metadata = append(metadata, "elements allowed: "+strings.Join(field.Element.Members, "|"))
		} else if field.Element.Format != "" {
			metadata = append(metadata, "element format: "+field.Element.Format)
		}
	}
	if field.MinInt != nil || field.MaxInt != nil {
		minimum, maximum := "unbounded", "unbounded"
		if field.MinInt != nil {
			minimum = fmt.Sprint(*field.MinInt)
		}
		if field.MaxInt != nil {
			maximum = fmt.Sprint(*field.MaxInt)
		}
		metadata = append(metadata, "integer range "+minimum+".."+maximum)
	}
	if field.MinLen != nil || field.MaxLen != nil {
		minimum, maximum := "0", "unbounded"
		if field.MinLen != nil {
			minimum = fmt.Sprint(*field.MinLen)
		}
		if field.MaxLen != nil {
			maximum = fmt.Sprint(*field.MaxLen)
		}
		metadata = append(metadata, "length "+minimum+".."+maximum)
	}
	return metadata
}

func p24JudgmentDetailsMatch(got, want *cli.InputError) bool {
	return got != nil && got.Path == want.Path && reflect.DeepEqual(got.Value, want.Value) && reflect.DeepEqual(got.Allowed, want.Allowed)
}

func TestP24JudgmentEvidenceRejectsIncompleteMetadata(t *testing.T) {
	t.Parallel()
	for leaf, descriptor := range judgmentDescriptors() {
		t.Run(leaf, func(t *testing.T) {
			t.Parallel()
			help := executeCommandForTest(t.TempDir(), append(strings.Split(leaf, " "), "--help")...).stdout
			if !p24JudgmentHelpMatches(help, descriptor.Fields, "input") {
				t.Fatal("valid help rejected")
			}
			var corrupt func([]cli.InputField, string)
			corrupt = func(fields []cli.InputField, prefix string) {
				for _, field := range fields {
					path := prefix + "." + field.Name
					for _, metadata := range append(p24JudgmentHelpMetadata(field), field.Meaning) {
						if metadata == "" {
							continue
						}
						lines := strings.Split(help, "\n")
						for i, line := range lines {
							if strings.HasPrefix(strings.TrimSpace(line), path+": ") {
								lines[i] = strings.Replace(line, metadata, "WRONG", 1)
								if lines[i] == line {
									t.Fatalf("seed did not change %s metadata %s", path, metadata)
								}
							}
						}
						if p24JudgmentHelpMatches(strings.Join(lines, "\n"), descriptor.Fields, "input") {
							t.Errorf("help checker accepted wrong %s metadata %s", path, metadata)
						}
					}
					corrupt(field.Children, path)
					if field.Element != nil {
						corrupt(field.Element.Children, path+"[i]")
					}
				}
			}
			corrupt(descriptor.Fields, "input")
			want := &cli.InputError{Path: "input." + descriptor.Fields[0].Name, Value: "integer", Allowed: []string{"string"}}
			if !p24JudgmentDetailsMatch(want, want) {
				t.Fatal("valid details rejected")
			}
			for _, got := range []*cli.InputError{
				{Path: want.Path, Value: want.Value, Allowed: []string{"wrong-nonempty"}},
				{Path: want.Path, Value: nil, Allowed: want.Allowed},
			} {
				if p24JudgmentDetailsMatch(got, want) {
					t.Errorf("details checker accepted wrong details %#v", got)
				}
			}
		})
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
		"run observe-budget":           "input.measurement_ref",
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
			if !strings.Contains(help.stdout, descriptor.Purpose) || !p24JudgmentHelpMatches(help.stdout, descriptor.Fields, "input") {
				t.Fatalf("help for %s lacks exact recursive field metadata: %s", leaf, help.stdout)
			}
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
			typed := reflect.New(descriptor.Type).Interface()
			decoder := json.NewDecoder(strings.NewReader(example.stdout))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(typed); err != nil {
				t.Fatal(err)
			}
			var expanded map[string]any
			if err := json.Unmarshal([]byte(example.stdout), &expanded); err != nil {
				t.Fatal(err)
			}
			p24ExpandDescriptorObjects(t, expanded, descriptor.Fields)
			payload, err := json.Marshal(expanded)
			if err != nil {
				t.Fatal(err)
			}
			if err := cli.ValidateInput(payload, descriptor); err != nil {
				t.Fatalf("template for %s: %v", leaf, err)
			}
			entries, err := os.ReadDir(workspace)
			if err != nil || len(entries) != 0 {
				t.Fatalf("example changed non-project directory: %v, %v", entries, err)
			}
		})
	}
}

func p24ReviewRecursiveStructuralDetails(t *testing.T, leaf string, descriptor cli.InputDescriptor) {
	t.Helper()
	template, err := cli.InputExample(descriptor)
	if err != nil {
		t.Fatal(err)
	}
	var baseline map[string]any
	if err := json.Unmarshal(template, &baseline); err != nil {
		t.Fatal(err)
	}
	p24ExpandDescriptorObjects(t, baseline, descriptor.Fields)
	if leaf == "review record" {
		baseline["expected_repository"].(map[string]any)["head_oid"] = "a0ca854133ae98b77e36cf51ad596ff0e8c6e096"
	}
	// Valid controls go through the same CLI, from a non-project directory.
	// Reaching the family fallback proves structural decoding was accepted.
	exercise := func(name string, payload []byte, want *cli.InputError) {
		t.Helper()
		t.Run("recursive/"+name, func(t *testing.T) {
			t.Parallel()
			validation := cli.ValidateInput(payload, descriptor)
			if want == nil {
				if validation != nil {
					t.Fatalf("qualifying structure: %v", validation)
				}
			} else if invalid, ok := validation.(*cli.InputError); !ok || !p24JudgmentDetailsMatch(invalid, want) {
				t.Fatalf("details = %#v, want %#v", validation, want)
			}
			workspace := t.TempDir()
			input := filepath.Join(t.TempDir(), "request.json")
			if err := os.WriteFile(input, payload, 0o600); err != nil {
				t.Fatal(err)
			}
			args := append([]string{"--json"}, strings.Split(leaf, " ")...)
			args = append(args, "--input", input, "--actor", "ACT-1", "--idempotency-key", "recursive")
			result := executeCommandForTest(workspace, args...)
			if cli.ExitCode(result.err) != 2 || result.stdout != "" {
				t.Fatalf("CLI: %#v", result)
			}
			if want == nil {
				want = &cli.InputError{Path: "input", Allowed: []string{"accepted " + descriptor.Operation + " input"}}
			}
			wireValue := want.Value
			if wireValue != nil {
				if _, ok := wireValue.(string); !ok {
					value, err := json.Marshal(wireValue)
					if err != nil {
						t.Fatal(err)
					}
					wireValue = string(value)
				}
			}
			assertAdoptedJSONFailure(t, result.stderr, descriptor.Operation, want.Path, wireValue, want.Allowed)
			var envelope map[string]any
			if err := json.Unmarshal([]byte(result.stderr), &envelope); err != nil {
				t.Fatal(err)
			}
			failure := envelope["error"].(map[string]any)
			if failure["retryable"] != false || failure["subject"] != nil {
				t.Fatalf("wrong failure metadata: %s", result.stderr)
			}
			entries, err := os.ReadDir(workspace)
			if err != nil || len(entries) != 0 {
				t.Fatalf("CLI changed non-project directory: %v / %v", entries, err)
			}
		})
	}
	encode := func() []byte {
		t.Helper()
		payload, err := json.Marshal(baseline)
		if err != nil {
			t.Fatal(err)
		}
		return payload
	}
	exercise("valid", encode(), nil)
	var walk func(map[string]any, []cli.InputField, string)
	walk = func(object map[string]any, fields []cli.InputField, prefix string) {
		keys := make([]string, len(fields))
		for i, field := range fields {
			keys[i] = field.Name
		}
		slices.Sort(keys)
		object["bogus"] = true
		payload := encode()
		delete(object, "bogus")
		exercise(prefix+"/unknown", payload, &cli.InputError{Path: prefix + ".bogus", Value: "bogus", Allowed: keys})
		for _, field := range fields {
			path := prefix + "." + field.Name
			original := object[field.Name]
			mutate := func(class string, value any, allowed []string, valid bool) {
				object[field.Name] = value
				payload := encode()
				object[field.Name] = original
				var want *cli.InputError
				if !valid {
					want = &cli.InputError{Path: path, Value: value, Allowed: allowed}
					switch class {
					case "null":
						want.Value = "null"
					case "type":
						if field.Type == "integer" {
							want.Value = "boolean"
						} else {
							want.Value = "integer"
						}
					}
				}
				exercise(path+"/"+class, payload, want)
			}
			delete(object, field.Name)
			payload := encode()
			object[field.Name] = original
			var missing *cli.InputError
			if field.Requirement != "optional" {
				missing = &cli.InputError{Path: path, Allowed: []string{field.Name}}
			}
			exercise(path+"/missing", payload, missing)
			mutate("null", nil, []string{field.Type}, field.Requirement == "nullable")
			if field.Type == "integer" {
				mutate("type", false, []string{field.Type}, false)
			} else {
				mutate("type", 1, []string{field.Type}, false)
			}
			if len(field.Members) > 0 {
				mutate("enum", "not-a-member", field.Members, false)
			}
			switch field.Format {
			case "oid":
				mutate("format", "not-an-oid", []string{"oid"}, false)
			case "ID", "digest", "relative-path":
				mutate("format", "", []string{field.Format}, false)
			case "text":
				allowed := []string{"text"}
				if len(field.Members) > 0 {
					// Closed membership rejects before the text-length check.
					allowed = field.Members
				}
				mutate("format", strings.Repeat("x", 4097), allowed, false)
			}
			if field.MinInt != nil {
				mutate("integer-minimum", float64(*field.MinInt-1), []string{fmt.Sprintf("integer range %d..unbounded", *field.MinInt)}, false)
			}
			if field.MinLen != nil && field.Type == "string" {
				mutate("minimum", "", []string{fmt.Sprintf("length %d..%d", *field.MinLen, *field.MaxLen)}, false)
			}
			if field.MaxLen != nil && field.Type == "string" {
				mutate("maximum", strings.Repeat("x", *field.MaxLen+1), []string{fmt.Sprintf("length %d..%d", *field.MinLen, *field.MaxLen)}, false)
			}
			if field.Element != nil {
				for _, bad := range []struct {
					value any
					kind  string
				}{{nil, "null"}, {1, "integer"}} {
					object[field.Name] = []any{bad.value}
					payload := encode()
					object[field.Name] = original
					exercise(path+"/element-"+bad.kind, payload, &cli.InputError{Path: path + "[0]", Value: bad.kind, Allowed: []string{field.Element.Type}})
				}
				if field.Element.Format == "ID" {
					object[field.Name] = []any{""}
					payload := encode()
					object[field.Name] = original
					exercise(path+"/element-format", payload, &cli.InputError{Path: path + "[0]", Value: "", Allowed: []string{"ID"}})
				}
				if field.MinLen != nil || field.MaxLen != nil || field.Format == "sorted-unique" {
					minimum, maximum := "0", "unbounded"
					if field.MinLen != nil {
						minimum = fmt.Sprint(*field.MinLen)
					}
					if field.MaxLen != nil {
						maximum = fmt.Sprint(*field.MaxLen)
					}
					arrays := map[string][]any{}
					if field.MinLen != nil && *field.MinLen > 0 {
						arrays["array-minimum"] = []any{}
					}
					if field.MaxLen != nil {
						oversize := make([]any, *field.MaxLen+1)
						var element any = "scope/path"
						if entries := original.([]any); len(entries) > 0 {
							element = entries[0]
						}
						for i := range oversize {
							oversize[i] = element
						}
						arrays["array-maximum"] = oversize
					}
					if field.Format == "sorted-unique" {
						arrays["sorted-unique"] = []any{"z", "a"}
					}
					for class, value := range arrays {
						object[field.Name] = value
						payload := encode()
						object[field.Name] = original
						encoded, err := json.Marshal(value)
						if err != nil {
							t.Fatal(err)
						}
						var raw []json.RawMessage
						if err := json.Unmarshal(encoded, &raw); err != nil {
							t.Fatal(err)
						}
						allowed := "length " + minimum + ".." + maximum
						if class == "sorted-unique" {
							allowed = "sorted-unique"
						}
						exercise(path+"/"+class, payload, &cli.InputError{Path: path, Value: raw, Allowed: []string{allowed}})
					}
				}
			}
			// Insert a duplicate into this exact object, even when sibling
			// objects contain the same field spelling.
			object["__duplicate_target"] = true
			marked := encode()
			delete(object, "__duplicate_target")
			value, err := json.Marshal(original)
			if err != nil {
				t.Fatal(err)
			}
			duplicate := strings.Replace(string(marked), `"__duplicate_target":true,`, `"`+field.Name+`":`+string(value)+`,`, 1)
			exercise(path+"/duplicate", []byte(duplicate), &cli.InputError{Path: path, Value: field.Name, Allowed: keys})
			if field.Type == "object" {
				walk(original.(map[string]any), field.Children, path)
			} else if field.Element != nil && field.Element.Type == "object" {
				walk(original.([]any)[0].(map[string]any), field.Element.Children, path+"[0]")
			}
		}
	}
	walk(baseline, descriptor.Fields, "input")
}

func p24JudgmentRecursiveStructuralDetails(t *testing.T, leaf string, descriptor cli.InputDescriptor) {
	t.Helper()
	template, err := cli.InputExample(descriptor)
	if err != nil {
		t.Fatal(err)
	}
	var baseline map[string]any
	if err := json.Unmarshal(template, &baseline); err != nil {
		t.Fatal(err)
	}
	p24ExpandDescriptorObjects(t, baseline, descriptor.Fields)
	// Valid controls go through the same CLI, from a non-project directory.
	// Reaching the family fallback proves structural decoding was accepted.
	exercise := func(name string, payload []byte, want *cli.InputError) {
		t.Helper()
		t.Run("recursive/"+name, func(t *testing.T) {
			t.Parallel()
			validation := cli.ValidateInput(payload, descriptor)
			if want == nil {
				if validation != nil {
					t.Fatalf("qualifying structure: %v", validation)
				}
			} else if invalid, ok := validation.(*cli.InputError); !ok || !p24JudgmentDetailsMatch(invalid, want) {
				t.Fatalf("details = %#v, want %#v", validation, want)
			}
			workspace := t.TempDir()
			input := filepath.Join(t.TempDir(), "request.json")
			if err := os.WriteFile(input, payload, 0o600); err != nil {
				t.Fatal(err)
			}
			args := append([]string{"--json"}, strings.Split(leaf, " ")...)
			args = append(args, "--input", input, "--actor", "ACT-1", "--idempotency-key", "recursive")
			result := executeCommandForTest(workspace, args...)
			if cli.ExitCode(result.err) != 2 || result.stdout != "" {
				t.Fatalf("CLI: %#v", result)
			}
			if want == nil {
				want = &cli.InputError{Path: "input", Allowed: []string{"accepted " + descriptor.Operation + " input"}}
			}
			wireValue := want.Value
			if wireValue != nil {
				if _, ok := wireValue.(string); !ok {
					value, err := json.Marshal(wireValue)
					if err != nil {
						t.Fatal(err)
					}
					wireValue = string(value)
				}
			}
			assertAdoptedJSONFailure(t, result.stderr, descriptor.Operation, want.Path, wireValue, want.Allowed)
			var envelope map[string]any
			if err := json.Unmarshal([]byte(result.stderr), &envelope); err != nil {
				t.Fatal(err)
			}
			failure := envelope["error"].(map[string]any)
			if failure["retryable"] != false || failure["subject"] != nil {
				t.Fatalf("wrong failure metadata: %s", result.stderr)
			}
			entries, err := os.ReadDir(workspace)
			if err != nil || len(entries) != 0 {
				t.Fatalf("CLI changed non-project directory: %v / %v", entries, err)
			}
		})
	}
	encode := func() []byte {
		t.Helper()
		payload, err := json.Marshal(baseline)
		if err != nil {
			t.Fatal(err)
		}
		return payload
	}
	exercise("valid", encode(), nil)
	var walk func(map[string]any, []cli.InputField, string)
	walk = func(object map[string]any, fields []cli.InputField, prefix string) {
		keys := make([]string, len(fields))
		for i, field := range fields {
			keys[i] = field.Name
		}
		slices.Sort(keys)
		object["bogus"] = true
		payload := encode()
		delete(object, "bogus")
		exercise(prefix+"/unknown", payload, &cli.InputError{Path: prefix + ".bogus", Value: "bogus", Allowed: keys})
		for _, field := range fields {
			path := prefix + "." + field.Name
			original := object[field.Name]
			mutate := func(class string, value any, allowed []string, valid bool) {
				object[field.Name] = value
				payload := encode()
				object[field.Name] = original
				var want *cli.InputError
				if !valid {
					want = &cli.InputError{Path: path, Value: value, Allowed: allowed}
					switch class {
					case "null":
						want.Value = "null"
					case "type":
						want.Value = "integer"
					}
				}
				exercise(path+"/"+class, payload, want)
			}
			delete(object, field.Name)
			payload := encode()
			object[field.Name] = original
			var missing *cli.InputError
			if field.Requirement != "optional" {
				missing = &cli.InputError{Path: path, Allowed: []string{field.Name}}
			}
			exercise(path+"/missing", payload, missing)
			mutate("null", nil, []string{field.Type}, field.Requirement == "nullable")
			mutate("type", 1, []string{field.Type}, false)
			if len(field.Members) > 0 {
				mutate("enum", "not-a-member", field.Members, false)
			}
			switch field.Format {
			case "ID", "digest":
				mutate("format", "", []string{field.Format}, false)
			case "text":
				allowed := []string{"text"}
				if len(field.Members) > 0 {
					// Closed membership rejects before the text-length check.
					allowed = field.Members
				}
				mutate("format", strings.Repeat("x", 4097), allowed, false)
			}
			if field.MinLen != nil && field.Type == "string" {
				mutate("minimum", "", []string{fmt.Sprintf("length %d..%d", *field.MinLen, *field.MaxLen)}, false)
			}
			if field.MaxLen != nil && field.Type == "string" {
				mutate("maximum", strings.Repeat("x", *field.MaxLen+1), []string{fmt.Sprintf("length %d..%d", *field.MinLen, *field.MaxLen)}, false)
			}
			if field.Element != nil {
				for _, bad := range []struct {
					value any
					kind  string
				}{{nil, "null"}, {1, "integer"}} {
					object[field.Name] = []any{bad.value}
					payload := encode()
					object[field.Name] = original
					exercise(path+"/element-"+bad.kind, payload, &cli.InputError{Path: path + "[0]", Value: bad.kind, Allowed: []string{field.Element.Type}})
				}
				if field.Element.Format == "ID" {
					object[field.Name] = []any{""}
					payload := encode()
					object[field.Name] = original
					exercise(path+"/element-format", payload, &cli.InputError{Path: path + "[0]", Value: "", Allowed: []string{"ID"}})
				}
				if field.MinLen != nil || field.MaxLen != nil || field.Format == "sorted-unique" {
					minimum, maximum := "0", "unbounded"
					if field.MinLen != nil {
						minimum = fmt.Sprint(*field.MinLen)
					}
					if field.MaxLen != nil {
						maximum = fmt.Sprint(*field.MaxLen)
					}
					arrays := map[string][]any{}
					if field.MinLen != nil && *field.MinLen > 0 {
						arrays["array-minimum"] = []any{}
					}
					if field.MaxLen != nil {
						oversize := make([]any, *field.MaxLen+1)
						var element any = "scope/path"
						if entries := original.([]any); len(entries) > 0 {
							element = entries[0]
						}
						for i := range oversize {
							oversize[i] = element
						}
						arrays["array-maximum"] = oversize
					}
					if field.Format == "sorted-unique" {
						arrays["sorted-unique"] = []any{"z", "a"}
					}
					for class, value := range arrays {
						object[field.Name] = value
						payload := encode()
						object[field.Name] = original
						encoded, err := json.Marshal(value)
						if err != nil {
							t.Fatal(err)
						}
						var raw []json.RawMessage
						if err := json.Unmarshal(encoded, &raw); err != nil {
							t.Fatal(err)
						}
						allowed := "length " + minimum + ".." + maximum
						if class == "sorted-unique" {
							allowed = "sorted-unique"
						}
						exercise(path+"/"+class, payload, &cli.InputError{Path: path, Value: raw, Allowed: []string{allowed}})
					}
				}
			}
			// Insert a duplicate into this exact object, even when sibling
			// objects contain the same field spelling.
			object["__duplicate_target"] = true
			marked := encode()
			delete(object, "__duplicate_target")
			value, err := json.Marshal(original)
			if err != nil {
				t.Fatal(err)
			}
			duplicate := strings.Replace(string(marked), `"__duplicate_target":true,`, `"`+field.Name+`":`+string(value)+`,`, 1)
			exercise(path+"/duplicate", []byte(duplicate), &cli.InputError{Path: path, Value: field.Name, Allowed: keys})
			if field.Type == "object" {
				walk(original.(map[string]any), field.Children, path)
			} else if field.Element != nil && field.Element.Type == "object" {
				walk(original.([]any)[0].(map[string]any), field.Element.Children, path+"[0]")
			}
		}
	}
	walk(baseline, descriptor.Fields, "input")
}

func p24JudgmentExpectedError(descriptor cli.InputDescriptor, name, path string) *cli.InputError {
	fields := descriptor.Fields
	parts := strings.Split(strings.TrimPrefix(path, "input."), ".")
	var field cli.InputField
	for i, part := range parts {
		key := strings.Split(part, "[")[0]
		for _, candidate := range fields {
			if candidate.Name == key {
				field = candidate
			}
		}
		if i < len(parts)-1 {
			fields = field.Children
			if field.Element != nil {
				fields = field.Element.Children
			}
		}
	}
	class := name[strings.LastIndex(name, " ")+1:]
	result := &cli.InputError{Path: path}
	switch class {
	case "unknown", "duplicate":
		result.Value = parts[len(parts)-1]
		for _, candidate := range fields {
			result.Allowed = append(result.Allowed, candidate.Name)
		}
		slices.Sort(result.Allowed)
	case "missing":
		result.Allowed = []string{field.Name}
	case "null", "type":
		result.Value = map[string]string{"null": "null", "type": "integer"}[class]
		result.Allowed = []string{field.Type}
	case "format":
		result.Value = ""
		result.Allowed = []string{field.Format}
		if field.Name == "display_label" {
			result.Allowed = []string{"length 1..128"}
		}
	case "enum":
		result.Value = "not-a-member"
		if len(parts) > 1 {
			result.Value = "unknown"
		}
		result.Allowed = field.Members
	}
	return result
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
			p24JudgmentRecursiveStructuralDetails(t, leaf, descriptor)
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
				wantDetails := p24JudgmentExpectedError(descriptor, name, want)
				t.Helper()
				t.Run(name, func(t *testing.T) {
					invalid, ok := cli.ValidateInput(payload, descriptor).(*cli.InputError)
					if !ok || !p24JudgmentDetailsMatch(invalid, wantDetails) {
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
					assertAdoptedJSONFailure(t, result.stderr, descriptor.Operation, want, wantDetails.Value, wantDetails.Allowed)
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

func TestP24RunDescriptorsAdoptSevenLeaves(t *testing.T) {
	t.Parallel()
	for _, leaf := range []string{"run start", "run close", "run abandon", "run observe-budget", "run observe", "run checkpoint", "run recover"} {
		descriptor, found := inputDescriptorRegistry[leaf]
		if !found || pendingInputDescriptors[leaf] {
			t.Errorf("%s lacks adopted descriptor", leaf)
			continue
		}
		if err := cli.CheckInputType(descriptor); err != nil {
			t.Errorf("%s: %v", leaf, err)
		}
	}
}

func TestP24RemovedRunRepositoryMembers(t *testing.T) {
	t.Parallel()
	for _, leaf := range []string{"run start", "run close", "run abandon", "run observe-budget"} {
		t.Run(leaf, func(t *testing.T) {
			t.Parallel()
			descriptor := inputDescriptorRegistry[leaf]
			data, err := cli.InputExample(descriptor)
			if err != nil {
				t.Fatal(err)
			}
			var input map[string]any
			if err := json.Unmarshal(data, &input); err != nil {
				t.Fatal(err)
			}
			p24ExpandDescriptorObjects(t, input, descriptor.Fields)
			path := "input.repository"
			if leaf == "run observe-budget" {
				input["measurement_ref"].(map[string]any)["repository"] = map[string]any{}
				path = "input.measurement_ref.repository"
			} else {
				input["repository"] = map[string]any{}
			}
			payload, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}
			root := t.TempDir()
			filename := filepath.Join(t.TempDir(), "input.json")
			if err := os.WriteFile(filename, payload, 0o600); err != nil {
				t.Fatal(err)
			}
			args := append([]string{"--json"}, strings.Split(leaf, " ")...)
			args = append(args, "--input", filename, "--actor", "ACT-1", "--idempotency-key", "removed")
			response := executeCommandForTest(root, args...)
			if cli.ExitCode(response.err) != 2 || response.stdout != "" || !strings.Contains(response.stderr, `"argument":"`+path+`"`) {
				t.Fatalf("removed member: %#v", response)
			}
			entries, err := os.ReadDir(root)
			if err != nil || len(entries) != 0 {
				t.Fatalf("refusal changed root: %v %v", entries, err)
			}
		})
	}
}
