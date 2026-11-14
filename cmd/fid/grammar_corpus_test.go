package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/mbell-nzl/fid/internal/cli"
)

func TestGrammarUnknownActivationCommand(t *testing.T) {
	t.Parallel()
	result := executeCommandForTest(t.TempDir(), "--json", "activation", "bogus")
	if cli.ExitCode(result.err) != 2 || result.stdout != "" {
		t.Fatalf("activation bogus = %#v, want stderr failure, exit 2", result)
	}
	var response struct {
		Operation string      `json:"operation"`
		Error     cli.Failure `json:"error"`
	}
	if err := json.Unmarshal([]byte(result.stderr), &response); err != nil {
		t.Fatal(err)
	}
	if response.Operation != "activation" || response.Error.Code != cli.Code("operation-unavailable") || response.Error.Message != `unknown command "bogus" for "fid activation"` {
		t.Fatalf("response = %#v", response)
	}
}

type grammarRow struct {
	name, operation, code, message, argument string
	args                                     []string
	value                                    any
	allowed                                  []string
}

func grammarCommands(root *cobra.Command) (leaves, groups []*cobra.Command) {
	var walk func(*cobra.Command)
	walk = func(cmd *cobra.Command) {
		if cmd.Name() == "help" || cmd.Name() == "completion" {
			return
		}
		if cmd.HasSubCommands() {
			groups = append(groups, cmd)
			for _, child := range cmd.Commands() {
				walk(child)
			}
		} else if cmd.CommandPath() != "fid journal inspect" && cmd.CommandPath() != "fid sql" {
			leaves = append(leaves, cmd)
		}
	}
	walk(root)
	return
}

func grammarOperation(cmd *cobra.Command) string {
	switch cmd.CommandPath() {
	case "fid init":
		return "project.init"
	case "fid rebuild":
		return "project.rebuild"
	case "fid backup":
		return "project.backup"
	case "fid restore":
		return "project.restore"
	case "fid status":
		return "project.status"
	case "fid show":
		return "issue.show"
	case "fid capabilities":
		return "capabilities.get"
	case "fid ready":
		return "ready.get"
	case "fid check":
		return "check.execute"
	case "fid handoff":
		return "task.handoff"
	case "fid start":
		return "task.start"
	case "fid heartbeat":
		return "hold.heartbeat"
	case "fid park":
		return "hold.park"
	case "fid answer":
		return "task.answer"
	case "fid takeover":
		return "task.takeover"
	case "fid next":
		return "next.get"
	case "fid resume":
		return "projection.resume"
	}
	return strings.ReplaceAll(strings.TrimPrefix(cmd.CommandPath(), "fid "), " ", ".")
}

func grammarUnknownFlagRow(cmd *cobra.Command) grammarRow {
	allowed := []string{}
	if !cmd.HasSubCommands() {
		cmd.LocalNonPersistentFlags().VisitAll(func(flag *pflag.Flag) {
			if flag.Name != "help" && flag.Name != "json" && flag.Name != "project" {
				allowed = append(allowed, "--"+flag.Name)
			}
		})
	}
	sort.Strings(allowed)
	return grammarRow{name: cmd.CommandPath() + " flag", operation: grammarOperation(cmd), code: "invalid-argument", message: "unknown flag: --bogus", argument: "--bogus", value: "--bogus", allowed: allowed, args: append(strings.Fields(strings.TrimPrefix(cmd.CommandPath(), "fid")), "--bogus")}
}

func grammarClassRows(t *testing.T) []grammarRow {
	t.Helper()
	root := newRootCommandAt(t.TempDir())
	leaves, groups := grammarCommands(root)
	if len(leaves) != 105 || len(groups) != 29 {
		t.Fatalf("live registry = %d leaves / %d groups", len(leaves), len(groups))
	}
	rows := []grammarRow{}
	for _, cmd := range append(leaves, groups...) {
		row := grammarUnknownFlagRow(cmd)
		rows = append(rows, row)
		short := row
		short.name = cmd.CommandPath() + " shorthand"
		short.argument = "-x"
		short.value = "-x"
		short.message = "unknown shorthand flag: 'x' in -x"
		short.args = append(strings.Fields(strings.TrimPrefix(cmd.CommandPath(), "fid")), "-x")
		rows = append(rows, short)
	}
	for _, cmd := range leaves {
		parts := strings.Fields(strings.TrimPrefix(cmd.CommandPath(), "fid"))
		use := strings.Fields(cmd.Use)[1:]
		args := append([]string{}, parts...)
		allowed := append([]string{}, use...)
		var value any
		message := fmt.Sprintf("accepts %d arg(s), received 0", len(use))
		if len(use) == 0 {
			args = append(args, "extra")
			value = "extra"
			allowed = []string{"no arguments"}
			message = strings.TrimPrefix(cmd.CommandPath(), "fid ") + " accepts no arguments"
		}
		if cmd.CommandPath() == "fid check" {
			continue
		}
		rows = append(rows, grammarRow{name: cmd.CommandPath() + " positional", operation: grammarOperation(cmd), code: "invalid-argument", message: message, argument: "arguments", value: value, allowed: allowed, args: args})
		if len(use) > 0 {
			extra := append([]string{}, parts...)
			for range use {
				extra = append(extra, "ID-1")
			}
			extra = append(extra, "extra")
			rows = append(rows, grammarRow{name: cmd.CommandPath() + " extra positional", operation: grammarOperation(cmd), code: "invalid-argument", message: fmt.Sprintf("accepts %d arg(s), received %d", len(use), len(use)+1), argument: "arguments", value: "extra", allowed: allowed, args: extra})
		}
		required := []string{}
		// Existing checks are exhaustive for input leaves; bootstrap Actor is optional.
		if cmd.Flags().Lookup("input") != nil {
			required = append(required, "input")
			if cmd.Flags().Lookup("actor") != nil && !strings.HasPrefix(cmd.CommandPath(), "fid activation ") {
				required = append(required, "actor")
			}
			if cmd.Flags().Lookup("idempotency-key") != nil {
				required = append(required, "idempotency-key")
			}
		} else {
			switch strings.Join(parts, " ") {
			case "issue create":
				required = []string{"actor", "idempotency-key", "kind", "title"}
			case "issue transition":
				required = []string{"actor", "idempotency-key", "to"}
			case "issue update", "label add", "issue label add", "issue label remove":
				required = []string{"actor", "idempotency-key"}
			case "issue ac add":
				required = []string{"text", "ordinal", "actor", "idempotency-key"}
			case "link add", "link remove":
				required = []string{"actor", "idempotency-key", "from-type", "to-type"}
			case "link list":
				required = []string{"type"}
			case "gate resolve":
				required = []string{"version"}
			case "export import":
				required = []string{"repository", "idempotency-key"}
			case "workflow define":
				required = []string{"name", "applies-to", "states", "categories", "transitions", "actor", "idempotency-key"}
			case "workflow assign":
				required = []string{"name", "actor", "idempotency-key"}
			}
		}
		for i, name := range required {
			args = append([]string{}, parts...)
			for range use {
				args = append(args, "ID-1")
			}
			for _, earlier := range required[:i] {
				args = append(args, "--"+earlier, grammarValidFlag(earlier))
			}
			rows = append(rows, grammarRow{name: cmd.CommandPath() + " missing " + name, operation: grammarOperation(cmd), code: "invalid-argument", message: "--" + name + " is required", argument: "--" + name, value: nil, allowed: []string{"--" + name}, args: args})
		}
	}
	for _, cmd := range groups {
		allowed := []string{}
		for _, child := range cmd.Commands() {
			if child.IsAvailableCommand() {
				allowed = append(allowed, child.Name())
			}
		}
		if cmd == root {
			allowed = append(allowed, "completion", "help")
		}
		sort.Strings(allowed)
		rows = append(rows, grammarRow{name: cmd.CommandPath() + " unknown", operation: grammarOperation(cmd), code: "operation-unavailable", message: fmt.Sprintf("unknown command %q for %q", "bogus", cmd.CommandPath()), argument: "command", value: "bogus", allowed: allowed, args: append(strings.Fields(strings.TrimPrefix(cmd.CommandPath(), "fid")), "bogus")})
	}
	return rows
}

func grammarValidFlag(name string) string {
	switch name {
	case "ordinal", "version":
		return "1"
	case "actor":
		return "agent:test"
	case "from-type", "to-type", "type":
		return "issue"
	case "name":
		return "example"
	case "applies-to", "kind":
		return "task"
	case "categories":
		return "unstarted,completed"
	case "states":
		return "todo,done"
	case "transitions":
		return "todo->done"
	}
	return "value"
}

func grammarExecute(dir string, args ...string) commandExecutionResult {
	root := newRootCommandAt(dir)
	var out, errout bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&errout)
	root.SetArgs(args)
	err := root.Execute()
	return commandExecutionResult{stdout: out.String(), stderr: errout.String(), err: err}
}

func grammarCheck(t *testing.T, row grammarRow, jsonMode bool, result commandExecutionResult, exit int) {
	t.Helper()
	if exit != 2 || result.stdout != "" {
		t.Fatalf("%s: exit %d result %#v", row.name, exit, result)
	}
	if !jsonMode {
		if result.stderr != p24GrammarHumanText(t, row) {
			t.Fatalf("%s Human = %q want %q", row.name, result.stderr, row.message)
		}
		return
	}
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal([]byte(result.stderr), &envelope); err != nil {
		t.Fatalf("%s: %v: %q", row.name, err, result.stderr)
	}
	if len(envelope) != 4 || string(envelope["contract_version"]) != `"fid-cli.v1"` || string(envelope["ok"]) != "false" || string(envelope["operation"]) != fmt.Sprintf("%q", row.operation) {
		t.Fatalf("%s envelope %s", row.name, result.stderr)
	}
	var failure cli.Failure
	if err := json.Unmarshal(envelope["error"], &failure); err != nil {
		t.Fatal(err)
	}
	wanted, _ := json.Marshal(map[string]any{"argument": row.argument, "value": row.value, "allowed": row.allowed})
	got, _ := json.Marshal(failure.Details)
	if string(got) != string(wanted) || string(failure.Code) != row.code || failure.Category != cli.CategoryInvalidInput || failure.Retryable || failure.Subject != nil || failure.Message != row.message {
		t.Fatalf("%s failure %s; expected details %s message %q", row.name, result.stderr, wanted, row.message)
	}
}

func grammarTree(t *testing.T, dir string) map[string]string {
	t.Helper()
	tree := map[string]string{}
	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		if entry.IsDir() {
			tree[rel] = "directory"
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		tree[rel] = fmt.Sprintf("%x", sha256.Sum256(data))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return tree
}

func TestGrammarWholeTreeClasses(t *testing.T) {
	t.Parallel()
	rows := grammarClassRows(t)
	workspace := clonePristineWorkspace(t)
	outside := t.TempDir()
	for _, dir := range []string{workspace, outside} {
		before := grammarTree(t, dir)
		for _, row := range rows {
			t.Run(filepath.Base(dir)+"/"+row.name, func(t *testing.T) {
				for _, jsonMode := range []bool{true, false} {
					args := append([]string{}, row.args...)
					if jsonMode {
						args = append([]string{"--json"}, args...)
					}
					result := grammarExecute(dir, args...)
					grammarCheck(t, row, jsonMode, result, cli.ExitCode(result.err))
				}
			})
		}
		if after := grammarTree(t, dir); !reflect.DeepEqual(before, after) {
			t.Fatal("grammar corpus changed workspace")
		}
	}
}

func grammarMessageRows() []grammarRow {
	rows := []grammarRow{}
	for _, kind := range []string{"flag", "word"} {
		prefix, suffix := "unknown flag: --", ""
		if kind == "word" {
			prefix = `unknown command "`
			suffix = `" for "fid"`
		}
		add := func(name, value, expected string) {
			row := grammarRow{name: kind + "/" + name, operation: "fid", code: "operation-unavailable", message: expected, argument: "command", value: value, allowed: []string{}, args: []string{value}}
			if kind == "flag" {
				row.operation = "stop.raise"
				row.code = "invalid-argument"
				row.argument = "--" + value
				row.value = "--" + value
				row.args = []string{"stop", "raise", "--" + value}
				row.allowed = []string{"--actor", "--example", "--idempotency-key", "--input"}
			}
			rows = append(rows, row)
		}
		if kind == "word" { // root help includes all listed groups/leaves, including Cobra commands.
			root := newRootCommandAt(".")
			allowed := []string{"completion", "help"}
			for _, cmd := range root.Commands() {
				if cmd.IsAvailableCommand() {
					allowed = append(allowed, cmd.Name())
				}
			}
			sort.Strings(allowed)
			defer func() {
				for i := range rows {
					if rows[i].argument == "command" {
						rows[i].allowed = allowed
					}
				}
			}()
		}
		add("short", "bogus", prefix+"bogus"+suffix)
		length := 600
		if kind == "flag" {
			length -= 2
		} // Include the two dashes in the 600-character spelling.
		add("long", strings.Repeat("a", length), prefix+strings.Repeat("a", 511-len(prefix))+"…")
		add("unicode", strings.Repeat("界", 600), prefix+strings.Repeat("界", 511-len(prefix))+"…")
		escape := "\\u001b"
		if kind == "word" {
			escape = "\\x1b"
		}
		for _, length := range []int{512, 513} {
			n := length - len(prefix) - len(suffix) - len(escape)
			value := strings.Repeat("a", n) + "\x1b"
			expected := prefix + strings.Repeat("a", n) + escape + suffix
			if length == 513 {
				// Put the escape near the beginning so only a literal scalar is cut.
				value = "\x1b" + strings.Repeat("a", n)
				expected = prefix + escape + strings.Repeat("a", n) + suffix
				expected = expected[:511] + "…"
			}
			add(fmt.Sprintf("escaped-%d", length), value, expected)
		}
		value := strings.Repeat("a", 509-len(prefix)) + "\x1b" + "tail"
		add("escape-crosses-cut", value, prefix+strings.Repeat("a", 509-len(prefix))+"…")
	}
	return rows
}

func TestGrammarMessageBoundaries(t *testing.T) {
	t.Parallel()
	for _, row := range grammarMessageRows() {
		t.Run(row.name, func(t *testing.T) {
			for _, jsonMode := range []bool{true, false} {
				args := append([]string{}, row.args...)
				if jsonMode {
					args = append([]string{"--json"}, args...)
				}
				result := grammarExecute(t.TempDir(), args...)
				grammarCheck(t, row, jsonMode, result, cli.ExitCode(result.err))
			}
		})
	}
}

func TestGrammarFlagDomainBeforeWorkspace(t *testing.T) {
	t.Parallel()
	rows := []grammarRow{
		{name: "limit", operation: "claim.list", code: "invalid-argument", args: []string{"claim", "list", "--limit", "0"}, argument: "--limit", value: "0", allowed: []string{"1..10000"}, message: `invalid argument "0" for "--limit" flag`},
		{name: "status", operation: "run.list", code: "invalid-argument", args: []string{"run", "list", "--status", "bogus"}, argument: "--status", value: "bogus", allowed: []string{"abandoned", "active", "failed", "interrupted", "succeeded"}, message: `invalid argument "bogus" for "--status" flag`},
		{name: "outcome", operation: "receipt.list", code: "invalid-argument", args: []string{"receipt", "list", "--outcome", "bogus"}, argument: "--outcome", value: "bogus", allowed: []string{"cancelled", "failed", "signaled", "spawn-failed", "succeeded", "timed-out"}, message: `invalid argument "bogus" for "--outcome" flag`},
		{name: "bool", operation: "issue.list", code: "invalid-argument", args: []string{"issue", "list", "--has-concerns=maybe"}, argument: "--has-concerns", value: "maybe", allowed: []string{"bool"}, message: `invalid argument "maybe" for "--has-concerns" flag: strconv.ParseBool: parsing "maybe": invalid syntax`},
		{name: "missing limit", operation: "claim.list", code: "invalid-argument", args: []string{"claim", "list", "--limit"}, argument: "--limit", value: nil, allowed: []string{"1..10000"}, message: `flag needs an argument: --limit`},
	}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			for _, dir := range []string{clonePristineWorkspace(t), t.TempDir()} {
				before := grammarTree(t, dir)
				for _, jsonMode := range []bool{true, false} {
					args := append([]string{}, row.args...)
					if jsonMode {
						args = append([]string{"--json"}, args...)
					}
					result := grammarExecute(dir, args...)
					grammarCheck(t, row, jsonMode, result, cli.ExitCode(result.err))
				}
				if !reflect.DeepEqual(before, grammarTree(t, dir)) {
					t.Fatal("domain failure wrote files")
				}
			}
		})
	}
}

func grammarDomainRows(t *testing.T) []grammarRow {
	t.Helper()
	leaves, _ := grammarCommands(newRootCommandAt(t.TempDir()))
	rows := []grammarRow{}
	typed := 0
	for _, cmd := range leaves {
		parts := strings.Fields(strings.TrimPrefix(cmd.CommandPath(), "fid"))
		positional := len(strings.Fields(cmd.Use)) - 1
		args := append([]string{}, parts...)
		for range positional {
			args = append(args, "ID-1")
		}
		cmd.LocalNonPersistentFlags().VisitAll(func(flag *pflag.Flag) {
			if flag.Name == "help" {
				return
			}
			domain := flagDomainFor(cmd, flag)
			if flag.Value.Type() != "string" && flag.Value.Type() != "stringSlice" {
				typed++
			}
			appendRow := func(name, raw, message string, value any) {
				rows = append(rows, grammarRow{name: cmd.CommandPath() + " " + name, operation: grammarOperation(cmd), code: "invalid-argument", argument: "--" + flag.Name, value: value, allowed: domain.Allowed, message: message, args: append(append([]string{}, args...), "--"+flag.Name+raw)})
			}
			if flag.NoOptDefVal == "" {
				appendRow("missing value "+flag.Name, "", "flag needs an argument: --"+flag.Name, nil)
			}
			switch flag.Value.Type() {
			case "bool":
				_, err := strconv.ParseBool("maybe")
				appendRow("malformed "+flag.Name, "=maybe", fmt.Sprintf("invalid argument %q for %q flag: %v", "maybe", "--"+flag.Name, err), "maybe")
			case "int":
				_, err := strconv.ParseInt("bogus", 0, 64)
				appendRow("malformed "+flag.Name, "=bogus", fmt.Sprintf("invalid argument %q for %q flag: %v", "bogus", "--"+flag.Name, err), "bogus")
				raw := "0"
				if domain.Accept(raw) {
					raw = "6"
				}
				if domain.Accept(raw) {
					t.Fatalf("no invalid bound for %s --%s", cmd.CommandPath(), flag.Name)
				}
				appendRow("outside domain "+flag.Name, "="+raw, fmt.Sprintf("invalid argument %q for %q flag", raw, "--"+flag.Name), raw)
			case "stringSlice":
				if !domain.Accept("\"bad\nlabel\"") {
					raw := "\"bad\nlabel\""
					appendRow("outside domain "+flag.Name, "="+raw, fmt.Sprintf("invalid argument %q for %q flag", raw, "--"+flag.Name), raw)
				}
			case "string":
				for _, raw := range []string{"!", "", "a\nb"} {
					if !domain.Accept(raw) {
						appendRow("outside domain "+flag.Name, "="+raw, fmt.Sprintf("invalid argument %q for %q flag", raw, "--"+flag.Name), raw)
						break
					}
				}
			default:
				t.Fatalf("uncovered live flag type %s", flag.Value.Type())
			}
		})
	}
	if typed != 65 {
		t.Fatalf("R18 corpus = %d want 65", typed)
	}
	return rows
}

func grammarRunRows(t *testing.T, rows []grammarRow, process bool) {
	t.Helper()
	workspace := clonePristineWorkspace(t)
	outside := t.TempDir()
	binary := ""
	if process {
		binary = buildAdoptedFID(t)
	}
	for _, dir := range []string{workspace, outside} {
		before := grammarTree(t, dir)
		for _, row := range rows {
			t.Run(filepath.Base(dir)+"/"+row.name, func(t *testing.T) {
				for _, jsonMode := range []bool{true, false} {
					args := append([]string{}, row.args...)
					if jsonMode {
						args = append([]string{"--json"}, args...)
					}
					var result commandExecutionResult
					var exit int
					if process {
						result.stdout, result.stderr, exit = runAdoptedFID(t, binary, dir, args...)
					} else {
						result = grammarExecute(dir, args...)
						exit = cli.ExitCode(result.err)
					}
					grammarCheck(t, row, jsonMode, result, exit)
				}
			})
		}
		if after := grammarTree(t, dir); !reflect.DeepEqual(before, after) {
			t.Fatal("grammar corpus changed workspace")
		}
	}
}

func TestGrammarWholeTreeFlagDomains(t *testing.T) {
	t.Parallel()
	grammarRunRows(t, grammarDomainRows(t), false)
}

func TestGrammarBuiltBinaryCorpus(t *testing.T) {
	t.Parallel()
	rows := append(grammarClassRows(t), grammarDomainRows(t)...)
	rows = append(rows, grammarMessageRows()...)
	root := newRootCommandAt(t.TempDir())
	near := grammarUnknownFlagRow(root)
	near.name = "near miss"
	near.code = "operation-unavailable"
	near.argument = "command"
	near.value = "isue"
	near.args = []string{"isue"}
	near.message = `unknown command "isue" for "fid"`
	for _, row := range rows {
		if row.name == "fid unknown" {
			near.allowed = row.allowed
		}
	}
	rows = append(rows, near)
	grammarRunRows(t, rows, true)
}

func TestGrammarGroupHelpWithoutWord(t *testing.T) {
	t.Parallel()
	root := newRootCommandAt(t.TempDir())
	_, groups := grammarCommands(root)
	binary := buildAdoptedFID(t)
	for _, cmd := range groups {
		args := strings.Fields(strings.TrimPrefix(cmd.CommandPath(), "fid"))
		result := grammarExecute(t.TempDir(), args...)
		if result.err != nil || result.stderr != "" || !strings.Contains(result.stdout, "Available Commands:") {
			t.Fatalf("%s help: %#v", cmd.CommandPath(), result)
		}
		out, errout, exit := runAdoptedFID(t, binary, t.TempDir(), args...)
		if exit != 0 || errout != "" || out != result.stdout {
			t.Fatalf("%s process help exit=%d stdout=%q stderr=%q", cmd.CommandPath(), exit, out, errout)
		}
	}
}

func grammarUndeclaredTyped(root *cobra.Command) []string {
	leaves, _ := grammarCommands(root)
	missing := []string{}
	for _, cmd := range leaves {
		cmd.LocalNonPersistentFlags().VisitAll(func(flag *pflag.Flag) {
			kind := flag.Value.Type()
			if flag.Name != "help" && kind != "string" && kind != "stringSlice" && flagDomainFor(cmd, flag).Citation == "pflag type; P24 §8 for --example" {
				missing = append(missing, cmd.CommandPath()+" --"+flag.Name)
			}
		})
	}
	sort.Strings(missing)
	return missing
}

func TestGrammarFlagDomainRegistryCompleteness(t *testing.T) {
	t.Parallel()
	root := newRootCommandAt(t.TempDir())
	if missing := grammarUndeclaredTyped(root); len(missing) > 0 {
		t.Fatalf("R18 flags without declarations: %v", missing)
	}
	cmd, _, err := root.Find([]string{"claim", "list"})
	if err != nil {
		t.Fatal(err)
	}
	cmd.Flags().Int("seeded-unregistered", 1, "seeded mismatch")
	if missing := grammarUndeclaredTyped(root); !reflect.DeepEqual(missing, []string{"fid claim list --seeded-unregistered"}) {
		t.Fatalf("seeded declaration gap = %v", missing)
	}
	leaves, _ := grammarCommands(root)
	for _, cmd := range leaves {
		for spelling, domain := range allDeclaredStringDomains(t, cmd) {
			if cmd.Flags().Lookup(spelling) == nil {
				t.Fatalf("domain absent from live flag set: %s --%s", cmd.CommandPath(), spelling)
			}
			if domain.Citation == "" {
				t.Fatalf("no accepted row citation: %s --%s", cmd.CommandPath(), spelling)
			}
		}
	}
}

func TestGrammarCSVAndKeyDomains(t *testing.T) {
	t.Parallel()
	rows := []grammarRow{
		{name: "blank state", operation: "workflow.define", code: "invalid-argument", args: []string{"workflow", "define", "--states="}, argument: "--states", value: "", allowed: []string{"string"}, message: `invalid argument "" for "--states" flag`},
		{name: "duplicate states", operation: "workflow.define", code: "invalid-argument", args: []string{"workflow", "define", "--states=todo,todo"}, argument: "--states", value: "todo,todo", allowed: []string{"string"}, message: `invalid argument "todo,todo" for "--states" flag`},
		{name: "malformed edge", operation: "workflow.define", code: "invalid-argument", args: []string{"workflow", "define", "--transitions=todo"}, argument: "--transitions", value: "todo", allowed: []string{"string"}, message: `invalid argument "todo" for "--transitions" flag`},
		{name: "key control", operation: "issue.create", code: "invalid-argument", args: []string{"issue", "create", "--idempotency-key=a\n"}, argument: "--idempotency-key", value: "a\n", allowed: []string{"string"}, message: `invalid argument "a\n" for "--idempotency-key" flag`},
	}
	grammarRunRows(t, rows, false)
}

func TestGrammarIntegerDomainsPreservePflagValues(t *testing.T) {
	t.Parallel()
	root := newRootCommandAt(t.TempDir())
	cmd, _, err := root.Find([]string{"claim", "list"})
	if err != nil {
		t.Fatal(err)
	}
	for _, raw := range []string{"1", "0x1", "01", "+1", "10000"} {
		if err := cmd.Flags().Set("limit", raw); err != nil {
			t.Errorf("in-range pflag integer %q rejected: %v", raw, err)
		}
	}
}

func TestGrammarAmendedFlagDomains(t *testing.T) {
	t.Parallel()
	root := newRootCommandAt(t.TempDir())
	for _, row := range []struct {
		path        []string
		flag, value string
	}{
		{[]string{"claim", "list"}, "status", "superseded"},
		{[]string{"receipt", "list"}, "outcome", "cancelled"},
	} {
		cmd, _, err := root.Find(row.path)
		if err != nil {
			t.Fatal(err)
		}
		if err := cmd.Flags().Set(row.flag, row.value); err != nil {
			t.Errorf("accepted amended domain %s --%s=%s: %v", cmd.CommandPath(), row.flag, row.value, err)
		}
	}
}

func TestGrammarInheritedStringDomains(t *testing.T) {
	t.Parallel()
	rows := []grammarRow{}
	for _, entry := range []struct {
		path        []string
		flag, value string
	}{
		{[]string{"issue", "create"}, "title", strings.Repeat("a", 513)},
		{[]string{"issue", "ac", "add", "ID-1"}, "text", "bad\ntext"},
		{[]string{"label", "add", "name"}, "colour", ""},
		{[]string{"issue", "list"}, "label", "bad\nlabel"},
		{[]string{"issue", "list"}, "status", "bad\nstate"},
		{[]string{"issue", "transition", "ID-1"}, "to", "bad\nstate"},
		{[]string{"claim", "list"}, "issue", "bad\nid"},
		{[]string{"decision", "list"}, "class", ""},
	} {
		path := entry.path
		operation := strings.Join(path, ".")
		if path[len(path)-1] == "ID-1" || path[len(path)-1] == "name" {
			operation = strings.Join(path[:len(path)-1], ".")
		}
		message := fmt.Sprintf("invalid argument %q for %q flag", entry.value, "--"+entry.flag)
		if entry.flag == "title" {
			prefix := `invalid argument "`
			message = prefix + strings.Repeat("a", 511-len(prefix)) + "…"
		}
		rows = append(rows, grammarRow{name: operation + " " + entry.flag, operation: operation, code: "invalid-argument", argument: "--" + entry.flag, value: entry.value, allowed: []string{"string"}, message: message, args: append(append([]string{}, path...), "--"+entry.flag+"="+entry.value)})
	}
	grammarRunRows(t, rows, false)
}

func allDeclaredStringDomains(t *testing.T, cmd *cobra.Command) map[string]flagDomain {
	result := map[string]flagDomain{}
	for _, domains := range []map[string]flagDomain{declaredStringDomains[rootCommandKey(cmd)], declaredTextDomains[rootCommandKey(cmd)]} {
		for spelling, domain := range domains {
			if _, exists := result[spelling]; exists {
				t.Fatalf("duplicate string domain %s --%s", cmd.CommandPath(), spelling)
			}
			result[spelling] = domain
		}
	}
	return result
}

func TestGrammarOpaqueStateUTF8Domain(t *testing.T) {
	t.Parallel()
	if declaredStringDomains["workflow define"]["states"].Accept(string([]byte{0xff})) {
		t.Fatal("invalid UTF-8 state accepted")
	}
}
