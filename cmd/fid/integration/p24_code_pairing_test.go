//go:build integration

package integration

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/gittest"
)

type p24CodePair struct {
	root, operation string
	args            []string
}
type p24CodeFactory func(*testing.T, string) p24CodePair

// The corpus is enumerated against RegisteredCodes below. A row invokes the
// real fid binary, with no injected mapper errors or simulated responses.
var p24CodeRows = map[string]p24CodeFactory{
	"lane-not-set":           p24LifecycleCodePair,
	"not-waiting":            p24SimplePair("task.answer", "answer", "FID-missing", "Yes", "--actor", "human:mark"),
	"hold-not-current":       p24LifecycleCodePair,
	"worktree-busy":          p24LifecycleCodePair,
	"already-held":           p24LifecycleCodePair,
	"unsupported-kind":       p24LifecycleCodePair,
	"internal-error":         p24InternalPair,
	"journal-replay-stopped": p24JournalPair,
	"lock-unavailable":       p24LockPair,
	"projection-divergence":  p24ProjectionPair,
	"import-target-conflict": p24ImportPair,

	"repository-unavailable": p24RetainedCodePair,

	"result-limit-exceeded": p24RetainedCodePair,

	"invalid-argument":      p24SimplePair("issue.create", "issue", "create", "--bogus"),
	"operation-unavailable": p24SimplePair("fid", "isue"),
	"project-not-found": func(t *testing.T, _ string) p24CodePair {
		return p24CodePair{gittest.NewWorkspace(t), "activity.list", []string{"activity", "list"}}
	},
	"project-already-initialized": p24SimplePair("project.init", "init", "--prefix", "P24"),
	"issue-not-found":             p24SimplePair("issue.get", "issue", "get", "FID-missing"),
	"label-not-found":             p24IssueCodePair,
	"state-conflict":              p24IssueCodePair,
	"hierarchy-invalid":           p24IssueCodePair,
	"transition-not-permitted":    p24IssueCodePair,
	"link-not-found":              p24IssueCodePair,
	"unsupported-relation":        p24IssueCodePair,
	"dependency-cycle":            p24IssueCodePair,

	"capability-unavailable": p24SimplePair("maintenance.vacuum", "maintenance", "vacuum", "--actor", "agent:a"),

	"archive-invalid": p24ArchivePair,

	"authority-denied": p24RetainedCodePair,

	"duplicate-entity": p24RetainedCodePair,
}

// Only exclusions proven by a current producer/caller audit belong here.
// Missing executable rows are failures, never silently classified unreachable.
var p24CodeExclusions = map[string]string{}

var p24SpecialCodePairs = map[string]func(*testing.T){
	"checks-unsatisfied":        handoffCodePair("checks-unsatisfied"),
	"test-paths-not-configured": handoffCodePair("test-paths-not-configured"),
	"tree-dirty":                handoffCodePair("tree-dirty"),
	"comparison-unavailable":    p24LeanPairName("comparison-unavailable"),
	"checks-not-configured":     p24LeanPairName("checks-not-configured"), "check-unsuccessful": p24LeanPairName("check-unsuccessful"), "check-not-qualified": p24LeanPairName("check-not-qualified"), "check-still-running": p24LeanPairName("check-still-running"), "check-incomplete": p24LeanPairName("check-incomplete"),
	"policy-unavailable": TestPolicyUnavailableHumanAndJSON,
}

// idempotency-conflict is a defined, reachable code that RegisteredCodes does
// not enumerate today. Keep its next-table process coverage explicitly, without
// changing that protected inventory or pretending it is a registry member.
var p24SupplementalPairs = map[string]p24CodeFactory{"idempotency-conflict": p24RetainedCodePair}

func TestP24RegisteredCodePairingCompleteness(t *testing.T) {
	t.Parallel()
	registered := map[string]bool{}
	for _, code := range cli.RegisteredCodes() {
		registered[string(code)] = true
	}
	for code := range p24CodeRows {
		if !registered[code] {
			t.Errorf("stale pairing row %s", code)
		}
	}
	for code := range p24SpecialCodePairs {
		if !registered[code] || p24CodeRows[code] != nil {
			t.Errorf("invalid special pairing row %s", code)
		}
	}
	for code, reason := range p24CodeExclusions {
		if !registered[code] || (p24CodeRows[code] != nil || p24SpecialCodePairs[code] != nil) || strings.TrimSpace(reason) == "" {
			t.Errorf("invalid exclusion %s: %s", code, reason)
		}
	}
	for code := range registered {
		if p24CodeRows[code] == nil && p24SpecialCodePairs[code] == nil && p24CodeExclusions[code] == "" {
			t.Errorf("registered code has no process pair or pinned exclusion: %s", code)
		}
	}
}

func TestP24RegisteredCodePairs(t *testing.T) {
	t.Parallel()
	for code, test := range p24SpecialCodePairs {
		t.Run(code, test)
	}
	rows := map[string]p24CodeFactory{}
	for code, prepare := range p24CodeRows {
		rows[code] = prepare
	}
	for code, prepare := range p24SupplementalPairs {
		rows[code] = prepare
	}
	for code, prepare := range rows {
		t.Run(code, func(t *testing.T) {
			t.Parallel()
			pair := prepare(t, code)
			binary := integrationFIDBinary(t)
			args := append([]string{"--project", pair.root}, pair.args...)
			if pair.operation == "export.import" {
				args = pair.args
			}
			out, stderr, err := closureRunCommand(pair.root, binary, append([]string{"--json"}, args...)...)
			if err == nil || out != "" {
				t.Fatalf("JSON %s: %v %q %q", code, err, out, stderr)
			}
			payload := p23Envelope(t, stderr)
			f := payload["error"].(map[string]any)
			if payload["operation"] != pair.operation || f["code"] != code {
				t.Fatalf("wrong producer for %s: %#v", code, payload)
			}
			status := processExitCode(t, err)
			for _, r := range f["message"].(string) {
				if unicode.IsControl(r) {
					t.Fatal("control character in JSON message")
				}
			}
			out, stderr, err = closureRunCommand(pair.root, binary, args...)
			want := p24ExpectedFailureText(t, payload, pair.operation, strings.Contains(strings.Join(args, " "), "--input") || pair.operation == "stop.raise")
			if err == nil || processExitCode(t, err) != status || out != "" || stderr != want {
				t.Fatalf("Human %s: %v stdout=%q stderr=%q want=%q", code, err, out, stderr, want)
			}
		})
	}
}

func p24SimplePair(operation string, args ...string) p24CodeFactory {
	return func(t *testing.T, _ string) p24CodePair {
		t.Helper()
		root := gittest.NewWorkspace(t)
		p21Run(t, integrationFIDBinary(t), root, "init", "--prefix", "P24", "--json")
		return p24CodePair{root, operation, args}
	}
}

func p24IssueCodePair(t *testing.T, code string) p24CodePair {
	t.Helper()
	root, binary := gittest.NewWorkspace(t), integrationFIDBinary(t)
	p21Run(t, binary, root, "init", "--prefix", "P24", "--json")
	create := func(title string) string {
		t.Helper()
		v := p21Run(t, binary, root, "--json", "issue", "create", "--kind", "task", "--title", title, "--actor", "agent:a", "--idempotency-key", title)
		return p21PathString(t, v, "result", "issue", "id")
	}
	id := create("first")
	args, operation := []string{}, ""
	actor := []string{"--actor", "agent:a", "--idempotency-key", "pair"}
	switch code {
	case "label-not-found":
		operation = "issue.label.add"
		args = []string{"issue", "label", "add", id, "missing"}
	case "state-conflict":
		operation = "issue.update"
		args = []string{"issue", "update", id, "--title", "changed", "--expected-revision", "999"}
	case "hierarchy-invalid":
		operation = "issue.create"
		args = []string{"issue", "create", "--kind", "subtask", "--title", "child"}
	case "transition-not-permitted":
		operation = "issue.transition"
		args = []string{"issue", "transition", id, "--to", "qa_pass", "--expected-revision", "1"}
	case "link-not-found":
		operation = "link.remove"
		args = []string{"link", "remove", id, "blocks", create("second"), "--from-type", "issue", "--to-type", "issue"}
	case "unsupported-relation":
		operation = "link.add"
		args = []string{"link", "add", id, "implements", create("second"), "--from-type", "issue", "--to-type", "issue"}
	case "dependency-cycle":
		peer := create("second")
		p21Run(t, binary, root, "--json", "link", "add", id, "blocks", peer, "--from-type", "issue", "--to-type", "issue", "--actor", "agent:a", "--idempotency-key", "edge")
		operation = "link.add"
		args = []string{"link", "add", peer, "blocks", id, "--from-type", "issue", "--to-type", "issue"}
	}
	return p24CodePair{root, operation, append(args, actor...)}
}

func p24ArchivePair(t *testing.T, _ string) p24CodePair {
	t.Helper()
	root := t.TempDir()
	path := filepath.Join(root, "invalid.zip")
	if err := os.WriteFile(path, []byte("invalid ZIP"), 0600); err != nil {
		t.Fatal(err)
	}
	return p24CodePair{root, "export.validate", []string{"export", "validate", path}}
}

// Independent section-9 oracle. It reads JSON values and never calls the
// production renderer or BoundedMessage. Unit fixtures separately pin the cuts.
func p24ExpectedFailureText(t *testing.T, payload map[string]any, operation string, input bool) string {
	t.Helper()
	f := payload["error"].(map[string]any)
	d := f["details"].(map[string]any)
	code := f["code"].(string)
	scalar := func(value any) string {
		if number, ok := value.(float64); ok && math.Trunc(number) == number {
			return strconv.FormatFloat(number, 'f', 0, 64)
		}
		return fmt.Sprint(value)
	}
	join := func(key, separator string) string {
		values, _ := d[key].([]any)
		parts := []string{}
		for _, v := range values {
			parts = append(parts, fmt.Sprint(v))
		}
		return strings.Join(parts, separator)
	}
	detail := ""
	switch code {
	case "invalid-argument":
		detail = "argument=" + fmt.Sprint(d["argument"]) + ", allowed=" + join("allowed", "|")
	case "idempotency-conflict":
		detail = "differing fields: " + join("differing_fields", ",")
	case "transition-not-permitted":
		detail = join("unsatisfied", ",")
	default:
		keys := []string{}
		for key := range d {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		parts := []string{}
		for _, key := range keys {
			value := d[key]
			text := ""
			switch value := value.(type) {
			case map[string]any:
				continue
			case []any:
				// Maintenance findings are object arrays even when empty.
				if code == "projection-divergence" && (key == "warnings" || key == "errors") {
					continue
				}
				scalars := []string{}
				nested := false
				for _, v := range value {
					if _, ok := v.(map[string]any); ok {
						nested = true
						break
					}
					scalars = append(scalars, fmt.Sprint(v))
				}
				if nested {
					continue
				}
				text = strings.Join(scalars, ",")
			case nil:
				text = "null"
			default:
				text = scalar(value)
			}
			parts = append(parts, key+"="+text)
		}
		detail = strings.Join(parts, ", ")
	}
	if code == "authority-denied" && operation == "next.claim" {
		detail = "actor_id=" + fmt.Sprint(d["actor_id"])
	}
	reason := "reason: " + code
	if detail != "" {
		reason += "; " + detail
	}
	reason = p24PairReasonCut(reason)
	text := "error: " + f["message"].(string) + "\n" + reason + "\n"
	next := map[string]string{
		"test-paths-not-configured": "add [review] test_adequacy, then merge .fid.toml and git fetch",
		"checks-not-configured":     "add checks to [lanes." + fmt.Sprint(d["lane"]) + "], then merge .fid.toml and git fetch",
		"check-still-running":       "wait for it to finish, or stop pid " + scalar(d["pid"]) + ", then check again",
		"lane-not-set":              "fid issue update " + fmt.Sprint(d["task_id"]) + " --lane <routine|standard|assured> --actor <actor>",
		"hold-not-current":          "fid next",

		"idempotency-conflict": "resend the original input, or use a new --idempotency-key",
	}[code]
	if code == "invalid-argument" && input {
		leaf := "fid " + strings.ReplaceAll(operation, ".", " ")
		next = "see " + leaf + " --help; template: " + leaf + " --example"
	}
	if next != "" {
		text += "next: " + next + "\n"
	}
	return text
}

func p24PairReasonCut(text string) string {
	var escaped strings.Builder
	for _, r := range text {
		switch r {
		case '\n':
			escaped.WriteString(`\n`)
		case '\r':
			escaped.WriteString(`\r`)
		case '\t':
			escaped.WriteString(`\t`)
		default:
			if unicode.IsControl(r) {
				fmt.Fprintf(&escaped, `\u%04x`, r)
			} else {
				escaped.WriteRune(r)
			}
		}
	}
	value := escaped.String()
	if utf8.RuneCountInString(value) <= 200 {
		return value
	}
	tokens := regexp.MustCompile(`\\(?:[abfnrtv\\"']|x[0-9a-fA-F]{2}|u[0-9a-fA-F]{4}|U[0-9a-fA-F]{8})|.`).FindAllString(value, -1)
	var result strings.Builder
	count := 0
	for _, token := range tokens {
		size := utf8.RuneCountInString(token)
		if count+size > 200 {
			break
		}
		result.WriteString(token)
		count += size
	}
	return result.String() + "…"
}
