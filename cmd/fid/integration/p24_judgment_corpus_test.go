//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var p24JudgmentFields = map[string][]string{
	"concern.record":               {"id", "summary"},
	"concern.dispose":              {"reason", "concern_id"},
	"deviation.record":             {"id", "reasoning"},
	"human-outcome.record":         {"id", "reasoning"},
	"decision.record-reference":    {"id", "class"},
	"requirement.record-reference": {"id", "source"},
	"requirement.link":             {"relation", "requirement_id"},
	"finding.address":              {"reason", "finding_id"},
	"review.record":                {"verdict", "run_id"},
}

func TestP24JudgmentKeyedLeafConflictAndCanonicalRetryCorpus(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	root := t.TempDir()
	workspace := filepath.Join(root, "leaf-contract")
	p24StartKeyedCapture(workspace)
	pairs := p23GenerateLeafPairs(t, binary, root)
	initial := p24EndKeyedCapture(workspace)
	seen := map[string]bool{}
	for _, pair := range pairs {
		names, ok := p24JudgmentFields[pair.Operation]
		if !ok {
			continue
		}
		seen[pair.Operation] = true
		for i, arg := range pair.JSONArgs {
			if arg == "--input" && i+1 < len(pair.JSONArgs) {
				input, err := os.ReadFile(pair.JSONArgs[i+1])
				if err != nil {
					t.Fatal(err)
				}
				p24AssertStoredFieldNames(t, pair.CWD, pair.Operation, p24CapturedKeyedCommand{args: pair.JSONArgs, input: input})
				break
			}
		}
		before := p21JournalEventCount(t, pair.CWD)
		original, originalErr, err := closureRunCommand(pair.CWD, binary, pair.JSONArgs...)
		if err != nil || originalErr != "" {
			t.Fatalf("%s baseline replay = %v stderr=%q", pair.Operation, err, originalErr)
		}
		entry, found := initial[pair.Operation]
		if !found {
			t.Fatalf("%s initial response not captured", pair.Operation)
		}
		if want := strings.Replace(entry.stdout, `"replayed":false`, `"replayed":true`, 1); original != want || want == entry.stdout {
			t.Fatalf("%s initial/replay JSON bytes differ beyond replayed", pair.Operation)
		}
		request := p23LeafRequest(t, pair)
		canonical := p24JudgmentArgs(t, pair, request, "canonical")
		for index, value := range canonical {
			if value == "--input" {
				path := canonical[index+1]
				if err := os.WriteFile(path, p24ReverseJSON(request), 0o600); err != nil {
					t.Fatal(err)
				}
			} else if value == "--idempotency-key" {
				canonical[index+1] = pair.JSONArgs[index+1]
			}
		}
		retry, retryErr, retryRunErr := closureRunCommand(pair.CWD, binary, canonical...)
		if retryRunErr != nil || retryErr != "" || retry != original {
			t.Fatalf("%s reordered retry = %v stderr=%q\n got=%q\nwant=%q", pair.Operation, retryRunErr, retryErr, retry, original)
		}
		for _, changedNames := range [][]string{{names[0]}, names, {"--actor"}} {
			changed := p23Clone(request)
			actor := false
			for _, name := range changedNames {
				if name == "--actor" {
					actor = true
					continue
				}
				switch value := changed[name].(type) {
				case string:
					if pair.Operation == "requirement.link" && name == "relation" {
						changed[name] = "verifies"
					} else if pair.Operation == "review.record" && name == "verdict" {
						if value != "changes-required" {
							t.Fatalf("review verdict fixture = %q, want changes-required", value)
						}
						changed[name] = "blocked"
					} else {
						changed[name] = value + "-p24-different"
					}
				case map[string]any:
					value["display_label"] = value["display_label"].(string) + "-p24-different"
				default:
					t.Fatalf("%s field %s cannot be changed: %#v", pair.Operation, name, changed[name])
				}
			}
			args := p24JudgmentArgs(t, pair, changed, strings.Join(changedNames, "-"))
			for index, value := range args {
				if value == "--idempotency-key" {
					args[index+1] = pair.JSONArgs[index+1]
				}
				if value == "--actor" && actor {
					args[index+1] = "actor-denied"
					if pair.JSONArgs[index+1] == "actor-denied" {
						args[index+1] = "actor-reviewer-two"
					}
				}
			}
			stdout, stderr, runErr := closureRunCommand(pair.CWD, binary, args...)
			if runErr == nil || processExitCode(t, runErr) != 4 || stdout != "" {
				t.Fatalf("%s %v conflict = %v stdout=%q stderr=%q", pair.Operation, changedNames, runErr, stdout, stderr)
			}
			got := p21Path(t, p23Envelope(t, stderr), "error", "details", "differing_fields")
			want := append([]string(nil), changedNames...)
			sort.Strings(want)
			values := make([]any, len(want))
			for i, name := range want {
				values[i] = name
			}
			if !reflect.DeepEqual(got, values) {
				t.Fatalf("%s %v differing_fields = %#v, want %#v", pair.Operation, changedNames, got, values)
			}
		}
		if after := p21JournalEventCount(t, pair.CWD); after != before {
			t.Fatalf("%s retry/conflicts appended Events: %d -> %d", pair.Operation, before, after)
		}
	}
	if len(seen) != len(p24JudgmentFields) {
		t.Fatalf("judgment corpus covered %v, want %v", seen, p24JudgmentFields)
	}
	p24AssertRebuildMapParity(t, pairs[0].CWD, binary)
}
