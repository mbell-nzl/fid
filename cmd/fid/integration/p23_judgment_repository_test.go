//go:build integration

package integration

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

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
)

type p24JudgmentRow struct {
	operation, key string
	input          map[string]any
}

func p24JudgmentExpectedRepository(t *testing.T, a *p22FIDRecoveryAdapter) map[string]any {
	t.Helper()
	git := func(args ...string) string {
		t.Helper()
		output, err := exec.Command("git", append([]string{"-C", a.workspace}, args...)...).Output()
		if err != nil {
			t.Fatalf("independent Git observation %v: %v", args, err)
		}
		return strings.TrimSpace(string(output))
	}
	head, tree := git("rev-parse", "HEAD"), git("rev-parse", "HEAD^{tree}")
	// These normal-state fixtures are clean, so the full dirty contract is
	// independently known: false with a null digest. Refuse to assume that
	// when Git reports a change; unavailable-state fixtures never call here.
	if status := git("status", "--porcelain", "--untracked-files=all", "--", "."); status != "" {
		t.Fatalf("repository parity fixture must be clean: %s", status)
	}
	root, err := filepath.Abs(a.workspace)
	if err != nil {
		t.Fatal(err)
	}
	projectConfig, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(a.workspace), "config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	return map[string]any{"base_oid": head, "capture_source": "git", "dirty": false, "dirty_digest": nil,
		"head_oid": head, "project_id": projectConfig.ProjectID, "tree_oid": tree,
		"worktree_id": fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(root)))}
}

func TestP24JudgmentRepositoryOracleUsesConfig(t *testing.T) {
	t.Parallel()
	a := p24JudgmentAdapter(t)
	configProjectID := a.projectID
	a.projectID = "decoy-cached-project-id"
	if got := p24JudgmentExpectedRepository(t, a)["project_id"]; got != configProjectID {
		t.Fatalf("oracle project_id = %v, want config %s", got, configProjectID)
	}
}

func p24RecordedSource(t *testing.T, operation string, envelope map[string]any) map[string]any {
	t.Helper()
	result := envelope["result"].(map[string]any)
	if operation == "decision.record-reference" {
		return result["decision"].(map[string]any)["source"].(map[string]any)
	}
	if operation == "requirement.record-reference" {
		return result["requirement"].(map[string]any)["source"].(map[string]any)
	}
	name := map[string]string{"human-outcome.record": "outcome", "concern.record": "concern", "deviation.record": "deviation"}[operation]
	return result[name].(map[string]any)["references"].([]any)[0].(map[string]any)
}

func p24JSONMap(t *testing.T, value any) map[string]any {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	return decoded
}

func p24EventByID(t *testing.T, root, id string) journal.Event {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(config.LedgerDir(root), "journal"))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		contents, err := os.ReadFile(filepath.Join(config.LedgerDir(root), "journal", entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range bytes.Split(bytes.TrimSpace(contents), []byte{'\n'}) {
			var event journal.Event
			if err := event.UnmarshalJSON(line); err != nil {
				t.Fatal(err)
			}
			if event.EventID == id {
				return event
			}
		}
	}
	t.Fatalf("Event %q was not found", id)
	return journal.Event{}
}

func TestP24JudgmentSourceDigestsAndRepositoryParity(t *testing.T) {
	t.Parallel()
	a := p24JudgmentAdapter(t)
	matching := p23DecisionSource(a, true)["digest"].(string)
	wrong := "sha256:" + p22Zeros(64)
	for index, row := range p24JudgmentRows(a, "matching", "CON-unused", "implements")[:5] {
		p24RowSource(row)["digest"] = matching
		wantRepository := p24JudgmentExpectedRepository(t, a)
		response := p24RunJudgment(t, a, row)
		if response.ExitCode != 0 || response.Stderr != "" {
			t.Fatalf("matching Source %d %s: %#v", index, row.operation, response)
		}
		envelope := p23Envelope(t, response.Stdout)
		stored := p24RecordedSource(t, row.operation, envelope)
		if stored["digest"] != matching || !reflect.DeepEqual(stored["repository"], wantRepository) {
			t.Fatalf("%s stored Source = %#v, observation = %#v", row.operation, stored, wantRepository)
		}
		event := p24EventByID(t, a.workspace, envelope["event"].(map[string]any)["id"].(string))
		if !reflect.DeepEqual(p24JSONMap(t, event.Repository), wantRepository) {
			t.Fatalf("%s Event repository differs from observation", row.operation)
		}
	}
	for _, row := range p24JudgmentRows(a, "omitted", "CON-unused", "implements")[:5] {
		response := p24RunJudgment(t, a, row)
		if response.ExitCode != 0 || p24RecordedSource(t, row.operation, p23Envelope(t, response.Stdout))["digest"] != matching {
			t.Fatalf("omitted Source %s: %#v", row.operation, response)
		}
	}
	for _, row := range p24JudgmentRows(a, "mismatch", "CON-unused", "implements")[:5] {
		p24RowSource(row)["digest"] = wrong
		before := p21JournalEventCount(t, a.workspace)
		response := p24RunJudgment(t, a, row)
		if response.ExitCode != 4 || response.Stdout != "" || p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("mismatching Source %s: %#v", row.operation, response)
		}
		var envelope map[string]any
		if err := json.Unmarshal([]byte(response.Stderr), &envelope); err != nil {
			t.Fatal(err)
		}
		failure := envelope["error"].(map[string]any)
		details := failure["details"].(map[string]any)
		if failure["code"] != "source-digest-mismatch" || details["expected_digest"] != wrong || details["observed_digest"] != matching {
			t.Fatalf("mismatch details %s: %#v", row.operation, failure)
		}
	}
}

func TestP24JudgmentSourceInputValidation(t *testing.T) {
	t.Parallel()
	a := p24JudgmentAdapter(t)
	for _, row := range p24JudgmentRows(a, "supplied-repository", "CON-unused", "implements")[:5] {
		p24RowSource(row)["repository"] = a.repository
		before := p21JournalEventCount(t, a.workspace)
		response := p24RunJudgment(t, a, row)
		want := "input.references[0].repository"
		if row.operation == "decision.record-reference" || row.operation == "requirement.record-reference" {
			want = "input.source.repository"
		}
		if response.ExitCode != 2 || !p23HasJSON(response.Stderr, `"argument":"`+want+`"`) || p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("supplied repository %s: %#v", row.operation, response)
		}
	}
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "outside.md"), []byte("outside"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "outside.md"), filepath.Join(a.workspace, "escape.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("base.txt", filepath.Join(a.workspace, "inside.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(a.workspace, "directory"), 0o700); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, path, code string }{
		{"absolute", filepath.Join(a.workspace, "base.txt"), "invalid-argument"},
		{"parent", "../base.txt", "invalid-argument"},
		{"escape", "escape.md", "source-digest-mismatch"},
		{"directory", "directory", "source-digest-mismatch"},
		{"inside", "inside.md", ""},
	} {
		row := p24JudgmentRows(a, "path-"+tc.name, "CON-unused", "implements")[0]
		p24RowSource(row)["path"] = tc.path
		before := p21JournalEventCount(t, a.workspace)
		response := p24RunJudgment(t, a, row)
		if tc.code == "" {
			if response.ExitCode != 0 || p24RecordedSource(t, row.operation, p23Envelope(t, response.Stdout))["digest"] != p23DecisionSource(a, true)["digest"] {
				t.Fatalf("in-root symlink: %#v", response)
			}
			continue
		}
		if response.ExitCode != map[string]int{"invalid-argument": 2, "source-digest-mismatch": 4}[tc.code] ||
			!p23HasJSON(response.Stderr, `"code":"`+tc.code+`"`) || p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("path %s: %#v", tc.name, response)
		}
		if tc.code == "invalid-argument" && !p23HasJSON(response.Stderr, `"argument":"input.source.path"`) {
			t.Fatalf("path %s did not name field: %s", tc.name, response.Stderr)
		}
		if tc.code == "source-digest-mismatch" && !p23HasJSON(response.Stderr, `"observed_digest":null`) {
			t.Fatalf("path %s reported observed digest: %s", tc.name, response.Stderr)
		}
	}
	row := p24JudgmentRows(a, "external-missing", "CON-unused", "implements")[3]
	source := p24RowSource(row)
	source["source_type"], source["path"], source["external_id"] = "external-reference", nil, "external-id-with-no-local-file"
	response := p24RunJudgment(t, a, row)
	if response.ExitCode != 2 || !p23HasJSON(response.Stderr, `"argument":"input.references[0].digest"`) {
		t.Fatalf("external digest omitted: %#v", response)
	}
	row = p24JudgmentRows(a, "external-supplied", "CON-unused", "implements")[3]
	source = p24RowSource(row)
	source["source_type"], source["path"], source["external_id"] = "external-reference", nil, "external-id-with-no-local-file"
	digest := "sha256:" + p22Ones(64)
	source["digest"] = digest
	response = p24RunJudgment(t, a, row)
	if response.ExitCode != 0 || p24RecordedSource(t, row.operation, p23Envelope(t, response.Stdout))["digest"] != digest {
		t.Fatalf("external digest supplied: %#v", response)
	}
}

func p24JudgmentAdapter(t *testing.T) *p22FIDRecoveryAdapter {
	t.Helper()
	actions := []string{"claim.acquire", "claim.release", "concern.record", "concern.dispose", "decision.record-reference", "deviation.record",
		"envelope.activate", "envelope.create", "finding.address", "gate.register", "human-outcome.record", "review.record",
		"run.close", "run.observe", "run.start"}
	return p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, actions,
		p22ExecutionFixture{})
}

func p24JudgmentRows(a *p22FIDRecoveryAdapter, suffix, concernID, relation string) []p24JudgmentRow {
	source := p23DecisionSource(a, true)
	delete(source, "digest")
	return []p24JudgmentRow{
		{"decision.record-reference", "decision-" + suffix, map[string]any{"id": "DREF-" + suffix, "class": "release", "source": source,
			"authority": map[string]any{"actor_id": a.actor, "required_capabilities": []string{"decision.record-reference"}}, "outcome": "approved", "scope": map[string]any{"issue_id": a.issueID}}},
		{"requirement.record-reference", "requirement-" + suffix, map[string]any{"id": "REQ-" + suffix, "source": source, "status": "active"}},
		{"human-outcome.record", "outcome-" + suffix, map[string]any{"id": "OUT-" + suffix, "class": "release", "subject": map[string]any{"type": "issue", "id": a.issueID},
			"result": "approved", "scope": map[string]any{"issue_id": a.issueID}, "reasoning": "approved with reference", "references": []any{source}}},
		{"concern.record", "concern-" + suffix, map[string]any{"id": "CON-" + suffix, "subject": map[string]any{"type": "issue", "id": a.issueID},
			"scope": map[string]any{"paths": []string{}}, "severity": "major", "summary": "source concern", "reasoning": "source evidence", "references": []any{source}}},
		{"deviation.record", "deviation-" + suffix, map[string]any{"id": "DEV-" + suffix, "subject": map[string]any{"type": "issue", "id": a.issueID},
			"expected": "expected", "observed": "observed", "authority": map[string]any{"actor_id": a.actor, "reference": nil},
			"disposition": "pending-decision", "reasoning": "source evidence", "references": []any{source}, "supersedes_deviation_id": nil}},
		{"concern.dispose", "dispose-" + suffix, map[string]any{"concern_id": concernID, "action": "carried-forward", "issue_id": a.issueID,
			"supersedes_concern_id": nil, "supersedes_finding_id": nil, "reason": "accepted disposition"}},
		{"requirement.link", "link-" + suffix, map[string]any{"requirement_id": "REQ-pre", "issue_id": a.issueID, "relation": relation}},
	}
}

func p24RunJudgment(t *testing.T, a *p22FIDRecoveryAdapter, row p24JudgmentRow) p22ConformanceResponse {
	t.Helper()
	parts := strings.Split(row.operation, ".")
	return p23RunInput(t, a, parts[0], parts[1], a.actor, "p24-"+row.key, row.input)
}

func p24PrepareJudgmentTargets(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	for _, id := range []string{"CON-pre-one", "CON-pre-two"} {
		response := p23RunInput(t, a, "concern", "record", a.actor, "p24-"+id, map[string]any{
			"id": id, "subject": map[string]any{"type": "issue", "id": a.issueID}, "scope": map[string]any{"paths": []string{}},
			"severity": "minor", "summary": "disposable concern", "reasoning": "disposition setup", "references": []any{},
		})
		if response.ExitCode != 0 {
			t.Fatalf("setup concern %s: %#v", id, response)
		}
	}
	source := p23DecisionSource(a, true)
	delete(source, "digest")
	response := p23RunInput(t, a, "requirement", "record-reference", a.actor, "p24-pre-requirement", map[string]any{
		"id": "REQ-pre", "source": source, "status": "active",
	})
	if response.ExitCode != 0 {
		t.Fatalf("setup requirement: %#v", response)
	}
}

func TestP24JudgmentLeavesRepositoryUnavailable(t *testing.T) {
	t.Parallel()
	for _, state := range []string{"not-a-repository", "unmerged-paths"} {
		t.Run(state, func(t *testing.T) {
			a := p24JudgmentAdapter(t)
			p24PrepareJudgmentTargets(t, a)
			normalRows := p24JudgmentRows(a, "normal", "CON-pre-one", "implements")
			originals := make(map[string]map[string]any, len(normalRows))
			for _, row := range normalRows {
				observation := p24JudgmentExpectedRepository(t, a)
				response := p24RunJudgment(t, a, row)
				if response.ExitCode != 0 || response.Stderr != "" {
					t.Fatalf("normal %s: %#v", row.operation, response)
				}
				envelope := p23Envelope(t, response.Stdout)
				event := p24EventByID(t, a.workspace, envelope["event"].(map[string]any)["id"].(string))
				if !reflect.DeepEqual(p24JSONMap(t, event.Repository), observation) {
					t.Fatalf("%s Event repository differs from independent observation", row.operation)
				}
				// The seven record shapes have no standalone Repository. Five
				// carry Sources; disposal's setup Concern has no references, and
				// RequirementLink has no Source. Walk every applicable member,
				// and pin the count so a missing Source cannot pass vacuously.
				recordName := map[string]string{"decision.record-reference": "decision", "requirement.record-reference": "requirement",
					"human-outcome.record": "outcome", "concern.record": "concern", "deviation.record": "deviation",
					"concern.dispose": "concern", "requirement.link": "link"}[row.operation]
				count := 0
				var compare func(any)
				compare = func(value any) {
					switch value := value.(type) {
					case map[string]any:
						for key, member := range value {
							if key == "repository" {
								count++
								if !reflect.DeepEqual(member, observation) {
									t.Fatalf("%s stored repository = %#v, want %#v", row.operation, member, observation)
								}
							} else {
								compare(member)
							}
						}
					case []any:
						for _, member := range value {
							compare(member)
						}
					}
				}
				compare(envelope["result"].(map[string]any)[recordName].(map[string]any))
				wantCount := 1
				if row.operation == "concern.dispose" || row.operation == "requirement.link" {
					wantCount = 0
				}
				if count != wantCount {
					t.Fatalf("%s stored repository members = %d, want %d", row.operation, count, wantCount)
				}
				originals[row.operation] = envelope
			}
			if state == "not-a-repository" {
				gittest.UnbornHead(t, a.workspace)
			} else {
				p24MakeUnmerged(t, a.workspace)
			}
			before := p21JournalEventCount(t, a.workspace)
			for _, row := range normalRows {
				response := p24RunJudgment(t, a, row)
				if response.ExitCode != 0 || response.Stderr != "" {
					t.Fatalf("%s unobservable replay: %#v", row.operation, response)
				}
				original := originals[row.operation]
				original["result"].(map[string]any)["replayed"] = true
				if !reflect.DeepEqual(original, p23Envelope(t, response.Stdout)) {
					t.Fatalf("%s replay differs from original", row.operation)
				}
				if after := p21JournalEventCount(t, a.workspace); after != before {
					t.Fatalf("%s replay appended %d Events", row.operation, after-before)
				}
			}
			for _, row := range p24JudgmentRows(a, "unsafe", "CON-pre-two", "verifies") {
				response := p24RunJudgment(t, a, row)
				if response.ExitCode != 4 || response.Stdout != "" {
					t.Fatalf("%s unavailable: %#v", row.operation, response)
				}
				var envelope map[string]any
				if err := json.Unmarshal([]byte(response.Stderr), &envelope); err != nil {
					t.Fatal(err)
				}
				failure := envelope["error"].(map[string]any)
				details := failure["details"].(map[string]any)
				if failure["code"] != "repository-unavailable" || failure["category"] != "precondition" || failure["retryable"] != true ||
					failure["subject"] != nil || !p23SameKeys(details, "detail reason") || details["reason"] != state || details["detail"] == "" {
					t.Fatalf("%s failure = %#v", row.operation, failure)
				}
				if after := p21JournalEventCount(t, a.workspace); after != before {
					t.Fatalf("%s appended %d Events", row.operation, after-before)
				}
			}
		})
	}
}

func p24MakeUnmerged(t *testing.T, root string) {
	t.Helper()
	command := exec.Command("git", "-C", root, "branch", "--show-current")
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	base := strings.TrimSpace(string(output))
	p22Git(t, root, "checkout", "-qb", "p24-conflict-side")
	p22Write(t, root, "base.txt", "side\n")
	p22Git(t, root, "commit", "-qam", "side")
	p22Git(t, root, "checkout", "-q", base)
	p22Write(t, root, "base.txt", "main\n")
	p22Git(t, root, "commit", "-qam", "main")
	merge := exec.Command("git", "-C", root, "merge", "p24-conflict-side")
	if output, err := merge.CombinedOutput(); err == nil {
		t.Fatalf("expected conflict, merge output: %s", output)
	}
}

func p24RowSource(row p24JudgmentRow) map[string]any {
	if source, ok := row.input["source"].(map[string]any); ok {
		return source
	}
	return row.input["references"].([]any)[0].(map[string]any)
}

func TestP24JudgmentLeavesReplayAfterHeadAndFileChange(t *testing.T) {
	t.Parallel()
	a := p24JudgmentAdapter(t)
	p24PrepareJudgmentTargets(t, a)
	rows := p24JudgmentRows(a, "replay", "CON-pre-one", "implements")
	originals := make(map[string]map[string]any, len(rows))
	for _, row := range rows {
		response := p24RunJudgment(t, a, row)
		if response.ExitCode != 0 {
			t.Fatalf("initial %s: %#v", row.operation, response)
		}
		originals[row.operation] = p23Envelope(t, response.Stdout)
	}
	p22Write(t, a.workspace, "base.txt", "advanced source\n")
	p22Git(t, a.workspace, "add", "base.txt")
	p22Git(t, a.workspace, "commit", "-qm", "advance HEAD and source")
	before := p21JournalEventCount(t, a.workspace)
	for index, row := range rows {
		response := p24RunJudgment(t, a, row)
		if response.ExitCode != 0 || response.Stderr != "" {
			t.Fatalf("changed HEAD replay %s: %#v", row.operation, response)
		}
		original := originals[row.operation]
		original["result"].(map[string]any)["replayed"] = true
		if !reflect.DeepEqual(original, p23Envelope(t, response.Stdout)) || p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("changed HEAD replay %s changed result or Event count", row.operation)
		}
		if index >= 5 {
			continue
		}
		conflict := p24JudgmentRows(a, "replay", "CON-pre-one", "implements")[index]
		p24RowSource(conflict)["digest"] = p23DecisionSource(a, true)["digest"]
		refused := p24RunJudgment(t, a, conflict)
		if refused.ExitCode != 4 || !p23HasJSON(refused.Stderr, `"code":"idempotency-conflict"`) ||
			!p23HasJSON(refused.Stderr, `"differing_fields":["`+map[bool]string{true: "source", false: "references"}[index < 2]+`"]`) ||
			p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("supplied/omitted digest conflict %s: %#v", row.operation, refused)
		}
	}
	if err := os.Remove(filepath.Join(a.workspace, "base.txt")); err != nil {
		t.Fatal(err)
	}
	for index, row := range rows[:5] {
		response := p24RunJudgment(t, a, row)
		if response.ExitCode != 0 || !reflect.DeepEqual(originals[row.operation], p23Envelope(t, response.Stdout)) ||
			p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("deleted-source replay %s: %#v", row.operation, response)
		}
		control := p24JudgmentRows(a, "deleted", "CON-pre-two", "verifies")[index]
		refused := p24RunJudgment(t, a, control)
		if refused.ExitCode != 4 || refused.Stdout != "" || p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("deleted-source control %s: %#v", row.operation, refused)
		}
		var envelope map[string]any
		if err := json.Unmarshal([]byte(refused.Stderr), &envelope); err != nil {
			t.Fatal(err)
		}
		failure := envelope["error"].(map[string]any)
		details := failure["details"].(map[string]any)
		if failure["code"] != "source-digest-mismatch" || details["expected_digest"] != nil || details["observed_digest"] != nil {
			t.Fatalf("deleted-source control %s: %#v", row.operation, failure)
		}
	}
	p22Write(t, a.workspace, "base.txt", "advanced source\n")
	for index, row := range p24JudgmentRows(a, "deleted", "CON-pre-two", "verifies")[:5] {
		response := p24RunJudgment(t, a, row)
		if response.ExitCode != 0 {
			t.Fatalf("restored-source control %d %s: %#v", index, row.operation, response)
		}
	}
}
