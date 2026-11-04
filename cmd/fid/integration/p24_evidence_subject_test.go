//go:build integration

package integration

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func TestP24ReviewAndHumanEvidenceSubjectRepositories(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"review", "human-outcome"} {
		t.Run(kind, func(t *testing.T) {
			t.Parallel()
			actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "evidence.link", "gate.register", "human-outcome.record", "review.record", "run.close", "run.observe", "run.start"}
			a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, "mutation", p22ExecutionFixture{ReviewObligation: true, HumanObligations: true})
			input := map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": kind, "receipt_id": nil, "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": []string{}}
			var expected map[string]any
			if kind == "review" {
				p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "subject-builder")
				claim := p23AcquireClaim(t, a, "actor-reviewer-two", "review", "subject-review")
				run := p23StartRun(t, a, "actor-reviewer-two", claim, "subject-review")
				recorded := p23ReviewRecord(t, a, "actor-reviewer-two", "subject-review", p23ReviewInput(a, p21PathString(t, run, "id"), p23Number(t, run["version"]), "pass", []any{}, nil))
				p23AssertReviewSuccess(t, "subject Review", recorded)
				review := p23ReviewFromResponse(t, recorded)
				input["review_id"], expected = review["id"], review["repository"].(map[string]any)
			} else {
				recorded := p23RecordHumanOutcome(t, a, "subject-human", "actor-human-outcome", "release", "approved")
				assertP23HumanOutcomeSuccess(t, recorded)
				input["obligation_id"] = "OBL-human-release"
				input["human_outcome_id"] = p21Path(t, p23Envelope(t, recorded.Stdout), "result", "outcome", "id")
				expected = p24SubjectEventRepository(t, a.workspace, p21PathString(t, p23Envelope(t, recorded.Stdout), "event", "id"))
			}
			// Current state differs, then becomes unobservable. Linking must still
			// copy the immutable subject, without Git observation.
			p22Write(t, a.workspace, "base.txt", "edited after subject\n")
			current := p22CorpusRepository(t, a.workspace, a.projectID)
			if reflect.DeepEqual(current, expected) {
				t.Fatal("subject and current repository must differ")
			}
			restoreHead := unbornGitHead(t, a.workspace)
			linked := p23RunInput(t, a, "evidence", "link", a.actor, "subject-link", input)
			p23AssertReviewSuccess(t, "subject link without Git", linked)
			envelope := p23Envelope(t, linked.Stdout)
			if target := p21Path(t, envelope, "result", "evidence", "target_repository"); !reflect.DeepEqual(target, expected) {
				t.Fatalf("%s target = %#v, want subject %#v", kind, target, expected)
			}
			if event := p24SubjectEventRepository(t, a.workspace, p21PathString(t, envelope, "event", "id")); !reflect.DeepEqual(event, expected) {
				t.Fatalf("%s linked Event = %#v, want %#v", kind, event, expected)
			}
			restoreHead()
			evaluated := a.p22Run(t, "evidence", "evaluate", "--input", p21Input(t, a.workspace, "subject-evaluate", map[string]any{"obligation_id": input["obligation_id"]}))
			p23AssertReviewSuccess(t, "subject evaluation", evaluated)
			want := "stale"
			if kind == "human-outcome" {
				want = "fresh"
			}
			if got := p21Path(t, p23Envelope(t, evaluated.Stdout), "result", "evaluation", "freshness"); got != want {
				t.Fatalf("%s freshness = %v, want %s", kind, got, want)
			}
		})
	}
}

func p24SubjectEventRepository(t *testing.T, root, eventID string) map[string]any {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(config.LedgerDir(root), "journal", "*.jsonl"))
	p22FaultOK(t, err)
	for _, path := range paths {
		data, err := os.ReadFile(path)
		p22FaultOK(t, err)
		for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
			var event map[string]any
			p22FaultOK(t, json.Unmarshal([]byte(line), &event))
			if event["event_id"] == eventID {
				return event["repository"].(map[string]any)
			}
		}
	}
	t.Fatalf("subject Event %s missing", eventID)
	return nil
}

func TestP24GateChangedFakeFIDInputParity(t *testing.T) {
	t.Parallel()
	changed := map[string]string{"gate.run": "repository", "evidence.link": "target_repository", "evidence.evaluate": "target_repository"}
	leaves, corpus := map[string]bool{}, map[string]bool{}
	file, err := parser.ParseFile(token.NewFileSet(), "p22_semantic_fake_test.go", nil, 0)
	p22FaultOK(t, err)
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Name.Name != "execute" {
			continue
		}
		ast.Inspect(function, func(node ast.Node) bool {
			literal, ok := node.(*ast.BasicLit)
			if ok && literal.Kind == token.STRING {
				operation, err := strconv.Unquote(literal.Value)
				p22FaultOK(t, err)
				if _, changed := changed[operation]; changed {
					leaves[operation] = true
				}
			}
			return true
		})
	}
	for _, fixture := range p22RecoveryFixtures(t) {
		if _, ok := changed[fixture.Operation]; ok {
			corpus[fixture.Operation] = true
		}
	}
	if !reflect.DeepEqual(leaves, corpus) {
		t.Fatalf("fake dispatch and shared corpus disagree: dispatch=%v corpus=%v", leaves, corpus)
	}
	if !reflect.DeepEqual(leaves, map[string]bool{"gate.run": true}) {
		t.Fatalf("Task 8 changed shared dispatch = %v", leaves)
	}
	for operation := range leaves {
		a := p24RunAdapter(t)
		p22Write(t, a.workspace, "base.txt", "changed before budget denial\n")
		a.repository = p22CorpusRepository(t, a.workspace, a.projectID)
		fake := p22NewFakeRecoveryAdapter()
		_, input := a.p22BuildArgs(t, p22ConformanceCase{Operation: operation})
		removed := p23Clone(input)
		removed[changed[operation]] = a.repository
		before := p21JournalEventCount(t, a.workspace)
		fid := p23RunInput(t, a, "gate", "run", a.actor, "parity-removed", removed)
		model := fake.model.execute(t, operation, "parity-removed", a.actor, removed, false)
		assertP22ProcessFailure(t, operation, "invalid-argument", 2, fid)
		assertP22ProcessFailure(t, operation, "invalid-argument", 2, model)
		if !reflect.DeepEqual(p21Path(t, p23Envelope(t, fid.Stderr), "error", "details"), p21Path(t, p23Envelope(t, model.Stderr), "error", "details")) {
			t.Fatalf("removed member details differ: fid=%s fake=%s", fid.Stderr, model.Stderr)
		}
		if p21JournalEventCount(t, a.workspace) != before || fake.model.sequence != 0 {
			t.Fatal("removed input changed state")
		}
		fid = p23RunInput(t, a, "gate", "run", a.actor, "parity-accepted", input)
		model = fake.model.execute(t, operation, "parity-accepted", a.actor, input, false)
		assertP22LiveDenial(t, fid, false)
		assertP22LiveDenial(t, model, false)
		p24AssertRunEventRepository(t, a, p23Clone(p22CorpusRepository(t, a.workspace, a.projectID)))
	}
}
