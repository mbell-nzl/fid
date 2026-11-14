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

func TestAdoptedExampleProcessExclusivity(t *testing.T) {
	t.Parallel()
	leaves := p24LiveInputLeaves(t)
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
	if got := len(leaves) * len(conflicts); got != 234 {
		t.Fatalf("exclusivity corpus has %d rows, want 234", got)
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
					} else if stderr != "error: --example cannot be combined with "+conflict.argument+"\nreason: invalid-argument; argument="+conflict.argument+", allowed=--example\nnext: see fid "+leaf+" --help; template: fid "+leaf+" --example\n" {
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

func p24LiveInputLeaves(t *testing.T) []string {
	t.Helper()
	var leaves []string
	var walk func(*cobra.Command)
	walk = func(command *cobra.Command) {
		if len(command.Commands()) == 0 && command.Flags().Lookup("input") != nil {
			leaves = append(leaves, rootCommandKey(command))
		}
		for _, child := range command.Commands() {
			walk(child)
		}
	}
	walk(newRootCommandAt(t.TempDir()))
	slices.Sort(leaves)
	if len(leaves) != 1 {
		t.Fatalf("live input leaves = %d, want 1", len(leaves))
	}
	return leaves
}

func TestP24WholeTreeDescriptorProcessCorpus(t *testing.T) {
	t.Parallel()
	leaves := make(map[string]cli.InputDescriptor)
	for _, leaf := range p24LiveInputLeaves(t) {
		descriptor, exists := inputDescriptorRegistry[leaf]
		if !exists {
			t.Fatalf("missing descriptor %s", leaf)
		}
		leaves[leaf] = descriptor
	}
	p24DescriptorProcessCorpus(t, leaves)
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
		if !live[key] {
			t.Errorf("registered %s is not live", key)
		}
		if err := cli.CheckInputType(descriptor); err != nil {
			t.Errorf("%s: %v", key, err)
		}
	}
	for key := range live {
		if _, registered := inputDescriptorRegistry[key]; !registered {
			t.Errorf("live --input leaf %s is unaccounted", key)
		}
	}
	mismatch := inputDescriptorRegistry["artifact verify"]
	mismatch.Fields[0].Name = "wrong_artifact_id"
	if err := cli.CheckInputType(mismatch); err == nil {
		t.Error("seeded descriptor/type mismatch passed")
	}
}

func TestP24InputDescriptorsHaveNoTopLevelSensitiveFields(t *testing.T) {
	t.Parallel()
	for leaf, descriptor := range inputDescriptorRegistry {
		for _, field := range descriptor.Fields {
			if field.Sensitive {
				t.Errorf("%s: top-level input.%s is sensitive", leaf, field.Name)
			}
		}
	}
}

func p24DescriptorProcessCorpus(t *testing.T, leaves map[string]cli.InputDescriptor) {
	t.Helper()
	binary := buildAdoptedFID(t)
	for leaf, descriptor := range leaves {
		t.Run(leaf, func(t *testing.T) {
			t.Parallel()
			if cli.CheckInputType(descriptor) != nil {
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
						for _, constraint := range p24JudgmentHelpMetadata(*field.Element) {
							if strings.HasPrefix(constraint, "integer range ") || strings.HasPrefix(constraint, "length ") {
								constraints = append(constraints, "element "+constraint)
							}
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
					for _, children := range field.Variants {
						variantPath := path
						if field.Type == "array" {
							variantPath += "[i]"
						}
						helpFields(children, variantPath)
					}
					if field.Element != nil {
						helpFields(field.Element.Children, path+"[i]")
					}
				}
			}
			helpFields(inputDescriptorRegistry[leaf].Fields, "input")
			template, err := cli.InputExample(inputDescriptorRegistry[leaf])
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
			baselineTemplate, err := cli.InputExample(descriptor)
			if err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(baselineTemplate, &baseline); err != nil {
				t.Fatal(err)
			}
			p24ReviewTemplateObject(t, baseline, p24SelectedFields(descriptor.Fields, baseline), "input")
			p24TemplateOrder(t, baselineTemplate, p24SelectedFields(descriptor.Fields, baseline), baseline)
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
					invalid, ok := cli.ValidateInput(data, inputDescriptorRegistry[leaf]).(*cli.InputError)
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
				fields = p24SelectedFields(fields, baseline)
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
						mutateElement := func(value any) {
							entries := append([]any(nil), original.([]any)...)
							if len(entries) == 0 {
								entries = []any{value}
							} else {
								entries[0] = value
							}
							object[field.Name] = entries
						}
						if field.MinLen != nil && *field.MinLen > 0 {
							check(path+"/short", path, object, func() { object[field.Name] = []any{} }, func() { object[field.Name] = original })
						}
						if field.MaxLen != nil {
							check(path+"/long", path, object, func() { object[field.Name] = make([]any, *field.MaxLen+1) }, func() { object[field.Name] = original })
						}
						if field.Element != nil {
							check(path+"/element-type", path+"[0]", object, func() { mutateElement(true) }, func() { object[field.Name] = original })
							if field.Element.Type == "string" && field.Element.Format != "" && field.Element.Format != "text" {
								check(path+"/element-format", path+"[0]", object, func() { mutateElement("") }, func() { object[field.Name] = original })
							}
							if len(field.Element.Members) > 0 {
								check(path+"/element-enum", path+"[0]", object, func() { mutateElement("not-a-member") }, func() { object[field.Name] = original })
							}
							if field.Element.MinInt != nil {
								check(path+"/element-min", path+"[0]", object, func() { mutateElement(*field.Element.MinInt - 1) }, func() { object[field.Name] = original })
							}
							if field.Element.MaxInt != nil {
								check(path+"/element-max", path+"[0]", object, func() { mutateElement(*field.Element.MaxInt + 1) }, func() { object[field.Name] = original })
							}
							if field.Element.MinLen != nil && *field.Element.MinLen > 0 {
								check(path+"/element-short", path+"[0]", object, func() { mutateElement("") }, func() { object[field.Name] = original })
							}
							if field.Element.MaxLen != nil {
								check(path+"/element-long", path+"[0]", object, func() { mutateElement(strings.Repeat("x", *field.Element.MaxLen+1)) }, func() { object[field.Name] = original })
							}
							if field.Format == "sorted-unique" {
								check(path+"/order", path, object, func() {
									entries := append([]any(nil), original.([]any)...)
									if len(entries) == 0 {
										entries = []any{"duplicate", "duplicate"}
									} else if len(entries) == 1 {
										entries = append(entries, entries[0])
									} else {
										entries[1] = entries[0]
									}
									object[field.Name] = entries
								}, func() { object[field.Name] = original })
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

func p24TemplateOrder(t *testing.T, data []byte, fields []cli.InputField, root map[string]any) {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(data))
	if token, err := decoder.Token(); err != nil || token != json.Delim('{') {
		t.Fatalf("template object: %v / %v", token, err)
	}
	for _, field := range p24SelectedFields(fields, root) {
		if field.Requirement == "optional" {
			continue
		}
		if token, err := decoder.Token(); err != nil || token != field.Name {
			t.Fatalf("template key order: %v / %v, want %s", token, err, field.Name)
		}
		var value json.RawMessage
		if err := decoder.Decode(&value); err != nil {
			t.Fatal(err)
		}
		if field.Requirement == "nullable" {
			if string(value) != "null" {
				t.Fatal("nullable template is not null")
			}
			continue
		}
		if field.Type == "object" {
			p24TemplateOrder(t, value, field.Children, root)
		}
		if field.Type == "array" {
			var entries []json.RawMessage
			if err := json.Unmarshal(value, &entries); err != nil {
				t.Fatal(err)
			}
			want := 0
			if field.MinLen != nil && *field.MinLen > 0 {
				want = 1
			}
			if len(entries) != want {
				t.Fatalf("template array %s: length=%d want %d", field.Name, len(entries), want)
			}
			for _, entry := range entries {
				if field.Element.Type == "object" {
					p24TemplateOrder(t, entry, field.Element.Children, root)
				}
			}
		}
	}
	if decoder.More() {
		t.Fatal("extra template member")
	}
	if _, err := decoder.Token(); err != nil {
		t.Fatal(err)
	}
}

func p24ExpandDescriptorObjects(t *testing.T, object map[string]any, fields []cli.InputField) {
	t.Helper()
	fields = p24SelectedFields(fields, object)
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
		if field.Format == "environment name" {
			object[field.Name] = "SECRET"
		}
		if field.Element != nil && field.Element.Format == "digest" {
			count := 1
			if field.MinLen != nil && *field.MinLen > count {
				count = *field.MinLen
			}
			values := make([]any, count)
			for i := range values {
				values[i] = fmt.Sprintf("sha256:%064x", i+1)
			}
			object[field.Name] = values
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
			count := 1
			if field.MinLen != nil && *field.MinLen > count {
				count = *field.MinLen
			}
			values := make([]any, count)
			for i := range values {
				values[i] = child
			}
			object[field.Name] = values
		} else {
			object[field.Name] = child
		}
	}
}

func p24SelectedFields(fields []cli.InputField, root map[string]any) []cli.InputField {
	selected := append([]cli.InputField(nil), fields...)
	for i, field := range selected {
		if len(field.Variants) == 0 {
			continue
		}
		selector := field.Discriminator
		if selector == "" {
			selector = "code"
		}
		children := field.Variants[fmt.Sprint(root[selector])]
		if field.Type == "array" {
			element := *field.Element
			element.Children = children
			selected[i].Element = &element
		} else {
			selected[i].Children = children
		}
	}
	return selected
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
			if field.MinLen != nil && *field.MinLen > 0 && len(array) != 1 {
				t.Fatalf("template %s must contain one editable element", path)
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
