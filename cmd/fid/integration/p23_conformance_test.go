//go:build integration

package integration

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"go/ast"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
)

type p23LeafCase struct {
	Operation     string            `json:"operation"`
	Fields        []string          `json:"fields"`
	Mutation      bool              `json:"mutation"`
	OwnerTask     int               `json:"owner_task"`
	Semantics     map[string]any    `json:"semantics,omitempty"`
	Amended       bool              `json:"amended,omitempty"`
	NotFoundID    string            `json:"not_found_id,omitempty"`
	NotFoundField string            `json:"not_found_field,omitempty"`
	NotFoundCases []p23NotFoundCase `json:"not_found_cases,omitempty"`
}

type p23NotFoundCase struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Field       string `json:"field"`
	Code        string `json:"code"`
	SubjectType string `json:"subject_type"`
}

type p23FailureCase struct {
	Code           string `json:"code"`
	Partition      string `json:"partition"`
	OwnerTask      int    `json:"owner_task"`
	BehavioralTest string `json:"behavioral_test"`
}

func TestP23FailureCorpusHasEveryAcceptedCode(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "p23_conformance", "failures.json"))
	if err != nil {
		t.Fatal(err)
	}
	var cases []p23FailureCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	var actual, expected []string
	added, retained := 0, 0
	for code := range p23FailureContracts {
		actual = append(actual, code)
	}
	parsed := map[string]*ast.File{}
	for _, fixture := range cases {
		if fixture.OwnerTask < 3 || fixture.OwnerTask > 13 {
			t.Fatalf("unowned failure %s", fixture.Code)
		}
		expected = append(expected, fixture.Code)
		switch fixture.Partition {
		case "added":
			added++
			if fixture.Code == "requirement-not-found" {
				t.Fatal("requirement-not-found must be retained")
			}
		case "retained":
			retained++
			if fixture.Code != "requirement-not-found" {
				t.Fatalf("unexpected retained failure %s", fixture.Code)
			}
		default:
			t.Fatalf("unpartitioned failure %s", fixture.Code)
		}
		p22CheckBehavioralReference(t, fixture.BehavioralTest, parsed)
	}
	p22SameMembers(t, "accepted failure-code inventory", actual, expected)
	if len(cases) != 19 || added != 18 || retained != 1 {
		t.Fatalf("declared/added/retained = %d/%d/%d, want 19/18/1", len(cases), added, retained)
	}
}

type p23LeafOutcome struct {
	Operation string
	Fields    []string
	Replay    bool
	Code      string
	Category  string
	Subject   any
	Details   map[string]any
	Exit      int
	Semantics map[string]any
}

func TestP23CallerProvidedReviewIDsRemainDistinct(t *testing.T) {
	t.Parallel()
	const first, second = "REV-abcdef", "REV-123456"
	a := p24RunAdapter(t)
	p22AcceptedCall(t, a, "review-id-proof-budget", "run.observe-budget", map[string]any{"value": 0})
	_, gateInput := a.p22BuildArgs(t, p22ConformanceCase{Operation: "gate.run"})
	gate := p23RunInput(t, a, "gate", "run", a.actor, "review-id-proof", gateInput)
	if gate.ExitCode != 0 || gate.Stderr != "" {
		t.Fatalf("Review ID fixture gate = %d %q", gate.ExitCode, gate.Stderr)
	}
	receiptID := p21PathString(t, p24GateTerminal(t, gate.Stdout), "response", "result", "receipt", "id")
	var models []map[string]any
	for _, reviewID := range []string{first, second} {
		request := map[string]any{
			"obligation_id": "OBL-input", "evidence_type": "receipt", "receipt_id": receiptID,
			"review_id": reviewID, "human_outcome_id": nil,
			"required_artifact_ids": []any{},
		}
		workspace := a.workspace
		pair := p23HumanLeafPair{Operation: "evidence.link", CWD: workspace, JSONArgs: []string{
			"--input", p21Input(t, workspace, "evidence", request), "--idempotency-key", "key",
		}}
		model := p23EvidenceLinkModel(t, pair)["evidence.link"]
		if got := model["evidence"].(map[string]any)["review_id"]; got != reviewID {
			t.Fatalf("Evidence review_id = %v, input = %s", got, reviewID)
		}
		models = append(models, model)
		got := p23NormalizeLeafValues(t, model, model)["evidence"].(map[string]any)["review_id"]
		if got != reviewID {
			t.Fatalf("normalized Evidence review_id = %v, input = %s", got, reviewID)
		}
	}
	if reflect.DeepEqual(p23NormalizeLeafValues(t, models[0], models[0]), p23NormalizeLeafValues(t, models[0], models[1])) {
		t.Fatal("different caller-provided Review IDs compare equal")
	}
	for _, field := range []string{"id", "envelope_id"} {
		firstValue := map[string]any{"input": map[string]any{field: "ENV-abcdef"}}
		secondValue := map[string]any{"input": map[string]any{field: "ENV-123456"}}
		if reflect.DeepEqual(p23NormalizeLeafValues(t, firstValue, firstValue), p23NormalizeLeafValues(t, firstValue, secondValue)) {
			t.Fatalf("caller-provided %s values compare equal", field)
		}
	}
}

func TestP23SharedLeafConformance(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "p23_conformance", "operations.json"))
	if err != nil {
		t.Fatal(err)
	}
	var cases []p23LeafCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) != 23 {
		t.Fatalf("Phase 23 leaf corpus has %d operations, want 23", len(cases))
	}
	root := t.TempDir()
	pairs := p23GenerateLeafPairs(t, integrationFIDBinary(t), root)
	byOperation := map[string]p23HumanLeafPair{}
	for _, pair := range pairs {
		byOperation[pair.Operation] = pair
	}
	model := p23RequirementModel(t, byOperation)
	concerns := p23ConcernModel(t, byOperation)
	for operation, result := range p23JudgmentModel(t, byOperation) {
		concerns[operation] = result
	}
	for operation, result := range p23ReviewFindingModel(t, byOperation) {
		concerns[operation] = result
	}
	accepted := map[string]p23LeafOutcome{}
	if len(byOperation) != 18 {
		t.Fatalf("Human fixture has %d new leaves, want 18", len(byOperation))
	}
	for _, fixture := range cases {
		if fixture.Amended {
			continue
		}
		fixture := fixture
		t.Run(fixture.Operation, func(t *testing.T) {
			if len(fixture.Semantics) == 0 {
				t.Fatalf("%s lacks contract-authored semantic expectations", fixture.Operation)
			}
			pair, found := byOperation[fixture.Operation]
			if !found {
				t.Fatalf("no real invocation for %s", fixture.Operation)
			}
			if fixture.OwnerTask < 4 || fixture.OwnerTask > 13 {
				t.Fatalf("unowned %s", fixture.Operation)
			}
			fake := p23FakeLeafAdapter{fixture: fixture, requirement: model, other: concerns}
			fid := p23FIDLeafAdapter{pair: pair, binary: integrationFIDBinary(t), fixture: fixture}
			actual := fid.Accepted(t)
			accepted[fixture.Operation] = actual
			p23CompareLeafOutcome(t, "accepted", fake.Accepted(), actual)
			if fixture.Mutation {
				p23CompareLeafOutcome(t, "unknown-field", fake.UnknownField(), fid.UnknownField(t))
				p23CompareLeafOutcome(t, "different-payload", fake.Conflict(pair), fid.Conflict(t))
				if fixture.Operation != "review.record" && fixture.Operation != "finding.address" {
					p23CompareLeafOutcome(t, "untrusted-actor", fake.UntrustedActor(), fid.UntrustedActor(t))
				}
			}
			if fixture.NotFoundID != "" && fixture.Operation != "review.record" {
				p23CompareLeafOutcome(t, "not-found", fake.NotFound(), fid.NotFound(t))
			}
			for _, missing := range fixture.NotFoundCases {
				p23CompareLeafOutcome(t, "not-found/"+missing.Name, p23ExpectedNotFound(fixture.Operation, missing), fid.NotFoundVariant(t, missing))
			}
			if fixture.Operation == "finding.address" {
				request := p23LeafRequest(t, pair)
				issue := p23LeafRequest(t, byOperation["concern.record"])["subject"].(map[string]any)["id"]
				args := append([]string{}, pair.JSONArgs...)
				for i, arg := range args {
					if arg == "--idempotency-key" {
						args[i+1] += "-stale-claim"
					}
				}
				p23CompareLeafOutcome(t, "stale-claim", p23LeafOutcome{Operation: fixture.Operation,
					Code: "claim-stale", Category: "conflict", Exit: 4,
					Subject: map[string]any{"type": "claim", "id": request["claim_id"]},
					Details: map[string]any{"issue_id": issue, "claim_id": request["claim_id"],
						"expected_version": request["expected_claim_version"], "actual_version": float64(2),
						"expires_at": "<generated:timestamp>"}}, fid.rejection(t, args))
			}
		})
	}
	p23CheckRequirementIdentity(t, accepted)
	p23CheckConcernIdentity(t, accepted)
	p23CheckJudgmentIdentity(t, accepted)
	p23CheckReviewFindingIdentity(t, accepted)
	t.Run("authority", func(t *testing.T) {
		a := p22NewFIDRecoveryAdapter(t)
		before := p21JournalEventCount(t, a.workspace)
		response := p23RunInput(t, a, "requirement", "record-reference", "actor-not-activated", "p23-shared-authority", map[string]any{
			"id": "REQ-shared-authority", "source": p23DecisionSource(a, true), "status": "active",
		})
		var failure struct {
			Operation string `json:"operation"`
			Error     struct {
				Code     string         `json:"code"`
				Category string         `json:"category"`
				Subject  any            `json:"subject"`
				Details  map[string]any `json:"details"`
			} `json:"error"`
		}
		if err := json.Unmarshal([]byte(response.Stderr), &failure); err != nil {
			t.Fatal(err)
		}
		p23CompareLeafOutcome(t, "authority", p23LeafOutcome{Operation: "requirement.record-reference", Code: "invalid-argument", Category: "invalid-input", Exit: 2,
			Details: map[string]any{"argument": "input", "value": nil, "allowed": []any{"accepted requirement.record-reference input"}}},
			p23LeafOutcome{Operation: failure.Operation, Code: failure.Error.Code, Category: failure.Error.Category,
				Subject: failure.Error.Subject, Details: failure.Error.Details, Exit: response.ExitCode})
		if got := p21JournalEventCount(t, a.workspace); got != before {
			t.Fatalf("authority denial appended Event %d -> %d", before, got)
		}
	})
	t.Run("stale-version", func(t *testing.T) {
		a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{},
			[]string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "gate.register", "review.record", "run.close", "run.observe", "run.start"})
		p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "p23-shared-before-review")
		claim := p23AcquireClaim(t, a, "actor-denied", "review", "p23-shared-review-claim")
		run := p23StartRun(t, a, "actor-denied", claim, "p23-shared-review-run")
		before := p21JournalEventCount(t, a.workspace)
		response := p23ReviewRecord(t, a, "actor-denied", "p23-shared-stale", p23ReviewInput(a, run["id"].(string), p23Number(t, run["version"])+1, "pass", []any{}, []any{}))
		var failure struct {
			Operation string `json:"operation"`
			Error     struct {
				Code     string         `json:"code"`
				Category string         `json:"category"`
				Subject  any            `json:"subject"`
				Details  map[string]any `json:"details"`
			} `json:"error"`
		}
		if err := json.Unmarshal([]byte(response.Stderr), &failure); err != nil {
			t.Fatal(err)
		}
		p23CompareLeafOutcome(t, "stale-version", p23LeafOutcome{Operation: "review.record", Code: "invalid-argument", Category: "invalid-input", Exit: 2,
			Details: map[string]any{"argument": "input", "value": nil, "allowed": []any{"accepted review.record input"}}},
			p23LeafOutcome{Operation: failure.Operation, Code: failure.Error.Code, Category: failure.Error.Category,
				Subject: failure.Error.Subject, Details: failure.Error.Details, Exit: response.ExitCode})
		if got := p21JournalEventCount(t, a.workspace); got != before {
			t.Fatalf("stale-version denial appended Event %d -> %d", before, got)
		}
	})
	t.Run("review.record/not-found", func(t *testing.T) {
		var fixture p23LeafCase
		for _, candidate := range cases {
			if candidate.Operation == "review.record" {
				fixture = candidate
				break
			}
		}
		a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{},
			[]string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "gate.register", "review.record", "run.close", "run.observe", "run.start"})
		p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "p23-shared-missing-before-review")
		claim := p23AcquireClaim(t, a, "actor-denied", "review", "p23-shared-missing-review-claim")
		run := p23StartRun(t, a, "actor-denied", claim, "p23-shared-missing-review-run")
		input := p23ReviewInput(a, run["id"].(string), p23Number(t, run["version"]), "pass", []any{}, []any{})
		args := []string{"--project", a.workspace, "--json", "review", "record", "--input", p21Input(t, a.workspace, "p23-shared-missing-review", input),
			"--actor", "actor-denied", "--idempotency-key", "p23-shared-missing-review"}
		pair := p23HumanLeafPair{Operation: fixture.Operation, CWD: a.workspace, JSONArgs: args}
		fake := p23FakeLeafAdapter{fixture: fixture}
		fid := p23FIDLeafAdapter{pair: pair, binary: a.binary, fixture: fixture}
		p23CompareLeafOutcome(t, "untrusted-actor", fake.UntrustedActor(), fid.UntrustedActor(t))
		p23CompareLeafOutcome(t, "not-found", fake.NotFound(), fid.NotFound(t))
	})
	t.Run("finding.address/missing-capability", func(t *testing.T) {
		var fixture p23LeafCase
		for _, candidate := range cases {
			if candidate.Operation == "finding.address" {
				fixture = candidate
				break
			}
		}
		a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{},
			[]string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "gate.register", "review.record", "run.close", "run.observe", "run.start"})
		p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "p23-shared-finding-before-review")
		reviewClaim := p23AcquireClaim(t, a, "actor-denied", "review", "p23-shared-finding-review-claim")
		reviewRun := p23StartRun(t, a, "actor-denied", reviewClaim, "p23-shared-finding-review-run")
		p23ReviewRecord(t, a, "actor-denied", "p23-shared-finding-review", p23ReviewInput(a,
			reviewRun["id"].(string), p23Number(t, reviewRun["version"]), "changes-required",
			[]any{p23Finding(a.issueID, "FND-p23-shared-authority", "major", true)}, nil))
		p23CloseAndRelease(t, a, "actor-denied", reviewRun["id"].(string), reviewClaim["id"].(string), "p23-shared-finding-review")
		claim := p23AcquireClaim(t, a, a.actor, "mutation", "p23-shared-finding-mutation-claim")
		input := map[string]any{"finding_id": "FND-p23-shared-authority", "claim_id": claim["id"],
			"expected_claim_version": claim["version"], "reason": "authority probe"}
		args := []string{"--project", a.workspace, "--json", "finding", "address", "--input", p21Input(t, a.workspace, "p23-shared-finding-authority", input),
			"--actor", a.actor, "--idempotency-key", "p23-shared-finding-authority"}
		pair := p23HumanLeafPair{Operation: fixture.Operation, CWD: a.workspace, JSONArgs: args}
		fake := p23FakeLeafAdapter{fixture: fixture}
		fid := p23FIDLeafAdapter{pair: pair, binary: a.binary, fixture: fixture}
		p23CompareLeafOutcome(t, "missing-capability", fake.MissingCapability(a.actor), fid.rejection(t, args))
	})
}

// p24ModelSource derives stored fields independently from the caller's input,
// project config and Git capture for the Phase 24 SourceInput convergence.
func p24ModelSource(t *testing.T, pair p23HumanLeafPair, raw any) map[string]any {
	t.Helper()
	input := raw.(map[string]any)
	stored := make(map[string]any, len(input)+1)
	for key, value := range input {
		stored[key] = value
	}
	if _, supplied := stored["digest"]; !supplied && stored["source_type"] == "project-file" {
		path := stored["path"].(string)
		contents, err := os.ReadFile(filepath.Join(pair.CWD, filepath.FromSlash(path)))
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(contents)
		stored["digest"] = "sha256:" + hex.EncodeToString(sum[:])
	}
	projectConfig, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(pair.CWD), config.ConfigFileName))
	if err != nil {
		t.Fatal(err)
	}
	repository := p22CorpusRepository(t, pair.CWD, projectConfig.ProjectID)
	if digest, ok := repository["dirty_digest"].(*string); ok {
		if digest == nil {
			repository["dirty_digest"] = nil
		} else {
			repository["dirty_digest"] = *digest
		}
	}
	stored["repository"] = repository
	return stored
}

func p24ModelSources(t *testing.T, pair p23HumanLeafPair, raw any) []any {
	t.Helper()
	inputs := raw.([]any)
	sources := make([]any, 0, len(inputs))
	for _, input := range inputs {
		sources = append(sources, p24ModelSource(t, pair, input))
	}
	return sources
}

func TestP23SharedAmendedLeafConformance(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "p23_conformance", "operations.json"))
	if err != nil {
		t.Fatal(err)
	}
	var shared []p23LeafCase
	if err := json.Unmarshal(data, &shared); err != nil {
		t.Fatal(err)
	}
	amended := map[string]p23LeafCase{}
	for _, fixture := range shared {
		if fixture.Amended {
			amended[fixture.Operation] = fixture
		}
	}
	if len(shared) != 23 || len(amended) != 5 {
		t.Fatalf("shared operations = %d, amended = %d; want 23 and 5", len(shared), len(amended))
	}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{},
		[]string{"claim.acquire", "envelope.activate", "envelope.create", "evidence.link", "gate.register", "gate.run", "run.observe-budget", "run.start"})
	p22AcceptedCall(t, a, "p23-shared-budget", "run.observe-budget", map[string]any{"value": 0})
	p22JourneyProof(t, a, "p23-shared-proof")
	cases := []struct {
		operation, file, key string
	}{
		{"envelope.create", "p22-corpus-envelope", "p22-corpus-envelope-create"},
		{"envelope.activate", "p22-corpus-envelope-activate", "p22-corpus-envelope-activate"},
		{"claim.acquire", "p22-corpus-claim", "p22-corpus-claim-acquire"},
		{"evidence.link", "p23-shared-proof-proof", "p23-shared-proof-proof"},
		{"evidence.evaluate", "p23-evaluate", ""},
	}
	p23EvidenceEvaluate(t, a)
	altGate := p22GateRun(t, a.binary, a.workspace, a.actor,
		p21Input(t, a.workspace, "p23-shared-alt-gate", map[string]any{
			"run_id": a.runID, "registration_id": a.registrationID, "registration_version": a.registrationVersion,
			"obligation_id": "OBL-p22-corpus", "expected_repository": p24ExpectedRepository(a.repository),
		}), "p23-shared-alt-gate")
	altReceiptID := p21PathString(t, altGate, "result", "receipt", "id")
	accepted := map[string]p23LeafOutcome{}
	for _, fixture := range cases {
		t.Run(fixture.operation, func(t *testing.T) {
			leaf, found := amended[fixture.operation]
			if !found || len(leaf.Semantics) == 0 {
				t.Fatalf("%s lacks contract-authored semantic expectations", fixture.operation)
			}
			parts := strings.Split(fixture.operation, ".")
			args := []string{"--project", a.workspace, "--json", parts[0], parts[1], "--input", filepath.Join(a.workspace, fixture.file+".json")}
			if leaf.Mutation {
				args = append(args, "--actor", a.actor, "--idempotency-key", fixture.key)
			}
			fake := p23FakeLeafAdapter{fixture: leaf, other: p23AmendedModel(t, fixture.operation, p23HumanLeafPair{Operation: fixture.operation, CWD: a.workspace, JSONArgs: args})}
			fid := p23FIDLeafAdapter{pair: p23HumanLeafPair{Operation: fixture.operation, CWD: a.workspace, JSONArgs: args}, binary: a.binary, fixture: leaf, alternateReceiptID: altReceiptID}
			actual := fid.Accepted(t)
			accepted[fixture.operation] = actual
			if leaf.Mutation {
				first, found := a.firstResults[fixture.operation]
				if !found {
					t.Fatalf("%s has no captured first response", fixture.operation)
				}
				first["replayed"] = true
				if !reflect.DeepEqual(first, actual.Semantics) {
					t.Fatalf("%s first response differs from replay", fixture.operation)
				}
			} else {
				second := fid.Accepted(t)
				if !reflect.DeepEqual(actual, second) {
					t.Fatalf("%s repeated read changed result", fixture.operation)
				}
			}
			p23CompareLeafOutcome(t, "accepted", fake.Accepted(), actual)
			if fixture.operation == "claim.acquire" {
				p23CheckClaimLease(t, actual)
			}
			if leaf.Mutation {
				p23CompareLeafOutcome(t, "unknown-field", fake.UnknownField(), fid.UnknownField(t))
				p23CompareLeafOutcome(t, "different-payload", fake.Conflict(fid.pair), fid.Conflict(t))
			}
			removed := map[string][]string{"envelope.create": {"repository_scope", "risk_result.source.repository"}, "claim.acquire": {"worktree_id"}, "evidence.link": {"target_repository"}, "evidence.evaluate": {"target_repository"}}
			for _, field := range removed[fixture.operation] {
				allowed := p24AdoptedTopLevelFields[fixture.operation]
				if strings.Contains(field, ".") {
					allowed = []any{"digest", "display_label", "external_id", "path", "section", "source_type"}
				}
				value := field[strings.LastIndex(field, ".")+1:]
				p23CompareLeafOutcome(t, "removed-member/"+field, p23LeafOutcome{Operation: fixture.operation,
					Code: "invalid-argument", Category: "invalid-input", Exit: 2,
					Details: map[string]any{"argument": "input." + field, "value": value, "allowed": allowed}}, fid.InvalidBinding(t, field, nil))
			}
			for _, missing := range leaf.NotFoundCases {
				p23CompareLeafOutcome(t, "not-found/"+missing.Name, p23ExpectedNotFound(fixture.operation, missing), fid.NotFoundVariant(t, missing))
			}
			if fixture.operation == "claim.acquire" || fixture.operation == "envelope.activate" {
				field, value := "envelope_revision", float64(2)
				if fixture.operation == "envelope.activate" {
					field, value = "expected_revision", float64(2)
				}
				p23CompareLeafOutcome(t, "stale-binding", p23LeafOutcome{Operation: fixture.operation,
					Code: "invalid-argument", Category: "invalid-input", Exit: 2,
					Details: map[string]any{"argument": "input", "value": nil, "allowed": []any{"accepted " + fixture.operation + " input"}}},
					fid.InvalidBinding(t, field, value))
			}
			if fixture.operation == "evidence.link" {
				p23CompareLeafOutcome(t, "repository-binding", p23LeafOutcome{Operation: fixture.operation,
					Code: "invalid-argument", Category: "invalid-input", Exit: 2,
					Details: map[string]any{"argument": "input.target_repository", "value": "target_repository", "allowed": p24AdoptedTopLevelFields["evidence.link"]}},
					fid.InvalidBinding(t, "target_repository", map[string]any{"worktree_id": "worktree-other"}))
			}
			if leaf.Mutation {
				p23CompareLeafOutcome(t, "untrusted-actor", fake.UntrustedActor(), fid.UntrustedActor(t))
			}
		})
	}
	p23CheckAmendedIdentity(t, accepted)
}

type p23FakeLeafAdapter struct {
	fixture     p23LeafCase
	requirement p23RequirementExpectation
	other       map[string]map[string]any
}

func (a p23FakeLeafAdapter) Accepted() p23LeafOutcome {
	if result, ok := a.requirement.results[a.fixture.Operation]; ok {
		return p23LeafOutcome{Operation: a.fixture.Operation, Fields: append([]string{}, a.fixture.Fields...), Replay: a.fixture.Mutation, Semantics: result}
	}
	if result, ok := a.other[a.fixture.Operation]; ok {
		return p23LeafOutcome{Operation: a.fixture.Operation, Fields: append([]string{}, a.fixture.Fields...), Replay: a.fixture.Mutation, Semantics: result}
	}
	return p23LeafOutcome{Operation: a.fixture.Operation, Fields: append([]string{}, a.fixture.Fields...), Replay: a.fixture.Mutation, Semantics: a.fixture.Semantics}
}

func (a p23FakeLeafAdapter) UnknownField() p23LeafOutcome {
	if allowed, ok := p24AdoptedTopLevelFields[a.fixture.Operation]; ok {
		return p23LeafOutcome{Operation: a.fixture.Operation, Code: "invalid-argument", Category: "invalid-input", Exit: 2,
			Details: map[string]any{"argument": "input.unaccepted_field", "value": "unaccepted_field", "allowed": allowed}}
	}
	return p23LeafOutcome{Operation: a.fixture.Operation, Code: "invalid-argument", Category: "invalid-input", Exit: 2,
		Details: map[string]any{"argument": "input", "value": nil, "allowed": []any{"accepted " + a.fixture.Operation + " input"}}}
}

var p24AdoptedTopLevelFields = map[string][]any{
	"evidence.link":                {"evidence_type", "human_outcome_id", "obligation_id", "receipt_id", "required_artifact_ids", "review_id"},
	"evidence.evaluate":            {"obligation_id"},
	"envelope.create":              {"authority_grants", "budgets", "non_goals", "obligations", "outcome", "primary_issue_id", "protected_decisions", "related_issue_ids", "risk_result", "sources"},
	"envelope.activate":            {"envelope_id", "expected_revision"},
	"envelope.complete":            {"envelope_id", "expected_revision"},
	"claim.acquire":                {"envelope_id", "envelope_revision", "issue_id", "mode", "scope", "ttl_seconds"},
	"review.record":                {"coverage", "evidence_refs", "expected_repository", "expected_run_version", "findings", "inapplicable_findings", "profile", "receipt_ids", "resolves_finding_ids", "run_id", "verdict"},
	"finding.address":              {"claim_id", "expected_claim_version", "finding_id", "reason"},
	"decision.record-reference":    {"authority", "class", "id", "outcome", "scope", "source"},
	"requirement.record-reference": {"id", "source", "status"},
	"requirement.link":             {"issue_id", "relation", "requirement_id"},
	"human-outcome.record":         {"class", "id", "reasoning", "references", "result", "scope", "subject"},
	"concern.record":               {"id", "reasoning", "references", "scope", "severity", "subject", "summary"},
	"concern.dispose":              {"action", "concern_id", "issue_id", "reason", "supersedes_concern_id", "supersedes_finding_id"},
	"deviation.record":             {"authority", "disposition", "expected", "id", "observed", "reasoning", "references", "subject", "supersedes_deviation_id"},
}

func (a p23FakeLeafAdapter) UntrustedActor() p23LeafOutcome {
	switch a.fixture.Operation {
	case "envelope.create", "envelope.activate":
		return p23LeafOutcome{Operation: a.fixture.Operation, Code: "authority-unavailable", Category: "authority", Exit: 5,
			Details: map[string]any{"operation": a.fixture.Operation, "actor_id": "actor-not-activated",
				"required_actions": []any{a.fixture.Operation}, "required_capabilities": []any{a.fixture.Operation}}}
	case "claim.acquire":
		return p23LeafOutcome{Operation: a.fixture.Operation, Code: "authority-denied", Category: "authority", Exit: 5,
			Details: map[string]any{"operation": a.fixture.Operation, "actor_id": "actor-not-activated",
				"required_actions": []any{a.fixture.Operation}, "required_capabilities": []any{a.fixture.Operation}}}
	}
	if a.fixture.Operation == "human-outcome.record" {
		return p23LeafOutcome{Operation: a.fixture.Operation, Code: "authority-missing", Category: "authority", Exit: 5,
			Details: map[string]any{"operation": a.fixture.Operation, "actor_id": "actor-not-activated",
				"required_actions": []any{a.fixture.Operation}, "required_capabilities": []any{a.fixture.Operation}}}
	}
	return p23LeafOutcome{Operation: a.fixture.Operation, Code: "invalid-argument", Category: "invalid-input", Exit: 2,
		Details: map[string]any{"argument": "input", "value": nil, "allowed": []any{"accepted " + a.fixture.Operation + " input"}}}
}

func (a p23FakeLeafAdapter) MissingCapability(actor string) p23LeafOutcome {
	return p23LeafOutcome{Operation: a.fixture.Operation, Code: "authority-denied", Category: "authority", Exit: 5,
		Details: map[string]any{"operation": a.fixture.Operation, "actor_id": actor,
			"required_actions": []any{a.fixture.Operation}, "required_capabilities": []any{a.fixture.Operation}}}
}

func (a p23FakeLeafAdapter) Conflict(pair p23HumanLeafPair) p23LeafOutcome {
	key := ""
	for i, arg := range pair.JSONArgs {
		if arg == "--idempotency-key" && i+1 < len(pair.JSONArgs) {
			key = pair.JSONArgs[i+1]
		}
	}
	changedField := map[string]string{
		"review.record": "profile", "finding.address": "reason", "concern.record": "summary",
		"concern.dispose": "reason", "deviation.record": "reasoning", "decision.record-reference": "outcome",
		"human-outcome.record": "reasoning", "requirement.record-reference": "status", "requirement.link": "relation",
		"envelope.create": "outcome", "envelope.activate": "--actor", "claim.acquire": "ttl_seconds", "evidence.link": "receipt_id",
	}[a.fixture.Operation]
	return p23LeafOutcome{Operation: a.fixture.Operation, Code: "idempotency-conflict", Category: "conflict", Exit: 4,
		Details: map[string]any{"operation_scope": a.fixture.Operation + ":<generated:project>", "idempotency_key": key,
			"stored_request_digest": "<generated:digest>", "supplied_request_digest": "<generated:digest>", "differing_fields": []any{changedField}}}
}

func (a p23FakeLeafAdapter) NotFound() p23LeafOutcome {
	code, subjectType := "", ""
	switch a.fixture.Operation {
	case "review.get":
		code, subjectType = "review-not-found", "review"
	case "requirement.get":
		code, subjectType = "requirement-not-found", "requirement"
	case "review.record", "finding.address":
		code, subjectType = "finding-not-found", "finding"
	case "concern.dispose":
		code, subjectType = "concern-not-found", "concern"
	case "deviation.record":
		code, subjectType = "deviation-not-found", "deviation"
	case "requirement.link":
		code, subjectType = "requirement-not-found", "requirement"
	}
	return p23LeafOutcome{Operation: a.fixture.Operation, Code: code, Category: "not-found", Exit: 3,
		Subject: map[string]any{"type": subjectType, "id": a.fixture.NotFoundID},
		Details: map[string]any{"id": a.fixture.NotFoundID}}
}

func p23ExpectedNotFound(operation string, missing p23NotFoundCase) p23LeafOutcome {
	if missing.Code == "invalid-argument" {
		return p23LeafOutcome{Operation: operation, Code: "invalid-argument", Category: "invalid-input", Exit: 2,
			Details: map[string]any{"argument": "input", "value": nil, "allowed": []any{"accepted " + operation + " input"}}}
	}
	return p23LeafOutcome{Operation: operation, Code: missing.Code, Category: "not-found", Exit: 3,
		Subject: map[string]any{"type": missing.SubjectType, "id": missing.ID},
		Details: map[string]any{"id": missing.ID}}
}

type p23FIDLeafAdapter struct {
	pair               p23HumanLeafPair
	binary             string
	fixture            p23LeafCase
	alternateReceiptID string
}

func (a p23FIDLeafAdapter) NotFoundVariant(t *testing.T, missing p23NotFoundCase) p23LeafOutcome {
	t.Helper()
	args := append([]string{}, a.pair.JSONArgs...)
	request := p23LeafRequest(t, a.pair)
	if a.fixture.Operation == "evidence.link" && missing.Name != "obligation" {
		request["evidence_type"] = missing.Name
		request["receipt_id"] = nil
	}
	fields := strings.Split(missing.Field, ".")
	target := request
	for _, field := range fields[:len(fields)-1] {
		nested, ok := target[field].(map[string]any)
		if !ok {
			t.Fatalf("%s has no object %s", a.fixture.Operation, field)
		}
		target = nested
	}
	target[fields[len(fields)-1]] = missing.ID
	for i, arg := range args {
		if arg == "--input" && i+1 < len(args) {
			args[i+1] = p21Input(t, a.pair.CWD, "p23-missing-"+missing.Name, request)
		}
		if arg == "--idempotency-key" && i+1 < len(args) {
			args[i+1] += "-missing-" + missing.Name
		}
	}
	return a.rejection(t, args)
}

func (a p23FIDLeafAdapter) NotFound(t *testing.T) p23LeafOutcome {
	t.Helper()
	args := append([]string{}, a.pair.JSONArgs...)
	if len(args) == 0 || a.fixture.NotFoundID == "" {
		t.Fatal("missing ID-taking case")
	}
	if a.fixture.NotFoundField == "" {
		args[len(args)-1] = a.fixture.NotFoundID
		return a.rejection(t, args)
	}
	request := p23LeafRequest(t, a.pair)
	if a.fixture.NotFoundField == "resolves_finding_ids" {
		request[a.fixture.NotFoundField] = []string{a.fixture.NotFoundID}
	} else {
		request[a.fixture.NotFoundField] = a.fixture.NotFoundID
	}
	if a.fixture.Operation == "deviation.record" {
		request["id"] = "DEV-missing-reference-probe"
	}
	for i, arg := range args {
		if arg == "--input" && i+1 < len(args) {
			args[i+1] = p21Input(t, a.pair.CWD, "p23-missing-"+strings.ReplaceAll(a.fixture.Operation, ".", "-"), request)
		}
		if arg == "--idempotency-key" && i+1 < len(args) {
			args[i+1] += "-missing"
		}
	}
	return a.rejection(t, args)
}

func (a p23FIDLeafAdapter) InvalidBinding(t *testing.T, field string, value any) p23LeafOutcome {
	t.Helper()
	args := append([]string{}, a.pair.JSONArgs...)
	request := p23LeafRequest(t, a.pair)
	fields := strings.Split(field, ".")
	target := request
	for _, part := range fields[:len(fields)-1] {
		nested, ok := target[part].(map[string]any)
		if !ok {
			t.Fatalf("%s has no object %s", a.fixture.Operation, part)
		}
		target = nested
	}
	target[fields[len(fields)-1]] = value
	for i, arg := range args {
		if arg == "--input" {
			args[i+1] = p21Input(t, a.pair.CWD, "p23-invalid-binding-"+strings.ReplaceAll(a.fixture.Operation, ".", "-"), request)
		}
		if arg == "--idempotency-key" {
			args[i+1] += "-invalid-binding"
		}
	}
	return a.rejection(t, args)
}

func (a p23FIDLeafAdapter) Accepted(t *testing.T) p23LeafOutcome {
	t.Helper()
	stdout, stderr, err := closureRunCommand(a.pair.CWD, a.binary, a.pair.JSONArgs...)
	if err != nil || stderr != "" {
		t.Fatalf("%s accepted: %v %q", a.pair.Operation, err, stderr)
	}
	var envelope struct {
		Operation string         `json:"operation"`
		Result    map[string]any `json:"result"`
	}
	if err := json.Unmarshal([]byte(stdout), &envelope); err != nil {
		t.Fatal(err)
	}
	if a.pair.JSONFacts != nil && !reflect.DeepEqual(envelope.Result, a.pair.JSONFacts) {
		t.Fatalf("%s duplicate delivery changed accepted result", a.pair.Operation)
	}
	fields := make([]string, 0, len(envelope.Result))
	for field := range envelope.Result {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	replay, _ := envelope.Result["replayed"].(bool)
	semantics := envelope.Result
	full := strings.HasPrefix(a.fixture.Operation, "requirement.") || strings.HasPrefix(a.fixture.Operation, "concern.") ||
		strings.HasPrefix(a.fixture.Operation, "deviation.") || strings.HasPrefix(a.fixture.Operation, "decision.") ||
		strings.HasPrefix(a.fixture.Operation, "human-outcome.") || strings.HasPrefix(a.fixture.Operation, "review.") ||
		strings.HasPrefix(a.fixture.Operation, "finding.") || a.fixture.Operation == "claim.acquire" ||
		a.fixture.Operation == "evidence.evaluate" || strings.HasPrefix(a.fixture.Operation, "envelope.") ||
		a.fixture.Operation == "evidence.link"
	if !full {
		semantics = p23SelectedSemantics(t, envelope.Result, a.fixture.Semantics)
	}
	return p23LeafOutcome{Operation: envelope.Operation, Fields: fields, Replay: replay, Semantics: semantics}
}

func p23SelectedSemantics(t *testing.T, result, paths map[string]any) map[string]any {
	t.Helper()
	selected := map[string]any{}
	for path := range paths {
		value := any(result)
		for _, segment := range strings.Split(path, ".") {
			switch current := value.(type) {
			case map[string]any:
				value = current[segment]
			case []any:
				index, err := strconv.Atoi(segment)
				if err != nil || index < 0 || index >= len(current) {
					t.Fatalf("semantic path %s has no index %s", path, segment)
				}
				value = current[index]
			default:
				t.Fatalf("semantic path %s stopped at %s", path, segment)
			}
		}
		selected[path] = value
	}
	return selected
}

func (a p23FIDLeafAdapter) UnknownField(t *testing.T) p23LeafOutcome {
	t.Helper()
	args := append([]string{}, a.pair.JSONArgs...)
	index := -1
	for i, arg := range args {
		if arg == "--input" {
			index = i + 1
			break
		}
	}
	if index < 0 || index >= len(args) {
		t.Fatal("mutation has no --input")
	}
	data, err := os.ReadFile(args[index])
	if err != nil {
		t.Fatal(err)
	}
	var request map[string]any
	if err := json.Unmarshal(data, &request); err != nil {
		t.Fatal(err)
	}
	request["unaccepted_field"] = true
	path := p21Input(t, a.pair.CWD, "p23-unknown-"+strings.ReplaceAll(a.pair.Operation, ".", "-"), request)
	args[index] = path
	for i, arg := range args {
		if arg == "--idempotency-key" && i+1 < len(args) {
			args[i+1] += "-unknown"
		}
	}
	return a.rejection(t, args)
}

func (a p23FIDLeafAdapter) UntrustedActor(t *testing.T) p23LeafOutcome {
	t.Helper()
	args := append([]string{}, a.pair.JSONArgs...)
	for i, arg := range args {
		if arg == "--actor" && i+1 < len(args) {
			args[i+1] = "actor-not-activated"
		}
		if arg == "--idempotency-key" && i+1 < len(args) {
			args[i+1] += "-untrusted"
		}
	}
	return a.rejection(t, args)
}

func (a p23FIDLeafAdapter) Conflict(t *testing.T) p23LeafOutcome {
	t.Helper()
	args := append([]string{}, a.pair.JSONArgs...)
	input := p23LeafRequest(t, a.pair)
	changed := p23ConflictInput(a.fixture.Operation, input)
	if changed == nil {
		changed = p23Clone(input)
		switch a.fixture.Operation {
		case "envelope.create":
			changed["outcome"] = input["outcome"].(string) + " different"
		case "envelope.activate":
			for i, arg := range args {
				if arg == "--actor" && i+1 < len(args) {
					args[i+1] = "actor-denied"
				}
			}
		case "evidence.link":
			if a.alternateReceiptID == "" {
				t.Fatal("missing alternative receipt")
			}
			changed["receipt_id"] = a.alternateReceiptID
		default:
			t.Fatalf("%s has no independent conflicting input", a.fixture.Operation)
		}
	}
	for i, arg := range args {
		if arg == "--input" && i+1 < len(args) {
			args[i+1] = p21Input(t, a.pair.CWD, "p23-conflict-"+strings.ReplaceAll(a.fixture.Operation, ".", "-"), changed)
		}
	}
	result := a.rejection(t, args)
	if result.Code == "idempotency-conflict" && result.Details["stored_request_digest"] == result.Details["supplied_request_digest"] {
		t.Fatalf("%s conflicting payload has equal request digests", a.fixture.Operation)
	}
	return result
}

func (a p23FIDLeafAdapter) rejection(t *testing.T, args []string) p23LeafOutcome {
	t.Helper()
	before := p21JournalEventCount(t, a.pair.CWD)
	_, stderr, err := closureRunCommand(a.pair.CWD, a.binary, args...)
	var failure struct {
		Operation string `json:"operation"`
		Error     struct {
			Code     string         `json:"code"`
			Category string         `json:"category"`
			Subject  any            `json:"subject"`
			Details  map[string]any `json:"details"`
		} `json:"error"`
	}
	if json.Unmarshal([]byte(stderr), &failure) != nil {
		t.Fatalf("%s unknown-field stderr: %q", a.pair.Operation, stderr)
	}
	var exit *exec.ExitError
	if !errors.As(err, &exit) {
		t.Fatalf("%s unknown-field exit: %v", a.pair.Operation, err)
	}
	if after := p21JournalEventCount(t, a.pair.CWD); after != before {
		t.Fatalf("%s rejection appended Event %d -> %d", a.pair.Operation, before, after)
	}
	return p23LeafOutcome{Operation: failure.Operation, Code: failure.Error.Code, Category: failure.Error.Category,
		Subject: failure.Error.Subject, Details: failure.Error.Details, Exit: exit.ExitCode()}
}

func p23CompareLeafOutcome(t *testing.T, name string, fake, fid p23LeafOutcome) {
	t.Helper()
	sort.Strings(fake.Fields)
	fid.Semantics = p23NormalizeLeafValues(t, fake.Semantics, fid.Semantics)
	fake.Semantics = p23NormalizeLeafValues(t, fake.Semantics, fake.Semantics)
	fid.Details = p23NormalizeLeafValues(t, fake.Details, fid.Details)
	fake.Details = p23NormalizeLeafValues(t, fake.Details, fake.Details)
	if !reflect.DeepEqual(fake, fid) {
		t.Fatalf("%s fake=%#v fid=%#v", name, fake, fid)
	}
}

// The independent fake marks fid-generated positions with placeholders. Use
// those positions for both adapters; every unmarked caller value stays exact.
func p23NormalizeLeafValues(t *testing.T, expected, raw map[string]any) map[string]any {
	t.Helper()
	if raw == nil {
		return nil
	}
	return p23NormalizeLeafValue(t, expected, raw).(map[string]any)
}

func p23NormalizeLeafValue(t *testing.T, expected, raw any) any {
	t.Helper()
	switch value := raw.(type) {
	case map[string]any:
		declared, _ := expected.(map[string]any)
		out := make(map[string]any, len(value))
		for child, item := range value {
			out[child] = p23NormalizeLeafValue(t, declared[child], item)
		}
		return out
	case []any:
		declared, _ := expected.([]any)
		out := make([]any, len(value))
		for i, item := range value {
			if i < len(declared) {
				out[i] = p23NormalizeLeafValue(t, declared[i], item)
			} else {
				out[i] = item
			}
		}
		return out
	case string:
		marker, ok := expected.(string)
		if !ok {
			return raw
		}
		if strings.HasSuffix(marker, ":<generated:project>") {
			prefix := strings.TrimSuffix(marker, "<generated:project>")
			if value != marker && (!strings.HasPrefix(value, prefix) || !p23ValidGeneratedToken("project", strings.TrimPrefix(value, prefix))) {
				t.Fatalf("invalid generated project scope %q", value)
			}
			return marker
		}
		if strings.HasPrefix(marker, "<generated:") && strings.HasSuffix(marker, ">") {
			kind := strings.TrimSuffix(strings.TrimPrefix(marker, "<generated:"), ">")
			if value != marker && !p23ValidGeneratedToken(kind, value) {
				t.Fatalf("invalid generated %s %q", kind, value)
			}
			return marker
		}
	}
	return raw
}

var p23GeneratedIDPattern = regexp.MustCompile(`^[0-9a-f]{6}([0-9a-f]{2})?$`)
var p23EventRefPattern = regexp.MustCompile(`^EVT-([0-9]{8})-([0-9]{6,})$`)

func p23ValidGeneratedToken(kind, value string) bool {
	switch kind {
	case "digest":
		if !strings.HasPrefix(value, "sha256:") || len(value) != 71 {
			return false
		}
		_, err := hex.DecodeString(strings.TrimPrefix(value, "sha256:"))
		return err == nil && value == strings.ToLower(value)
	case "timestamp":
		_, err := time.Parse(time.RFC3339Nano, value)
		return err == nil
	case "EventRef":
		parts := p23EventRefPattern.FindStringSubmatch(value)
		if parts == nil {
			return false
		}
		_, err := time.Parse("20060102", parts[1])
		return err == nil && strings.TrimLeft(parts[2], "0") != ""
	case "project":
		if !strings.HasPrefix(value, "project-") || len(value) != len("project-")+32 {
			return false
		}
		_, err := hex.DecodeString(strings.TrimPrefix(value, "project-"))
		return err == nil && value == strings.ToLower(value)
	default:
		return strings.HasPrefix(value, kind+"-") && p23GeneratedIDPattern.MatchString(strings.TrimPrefix(value, kind+"-"))
	}
}
