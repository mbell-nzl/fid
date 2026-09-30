//go:build integration

package integration

import (
	"encoding/json"
	"errors"
	"go/ast"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type p23LeafCase struct {
	Operation string   `json:"operation"`
	Fields    []string `json:"fields"`
	Mutation  bool     `json:"mutation"`
	OwnerTask int      `json:"owner_task"`
}

type p23FailureCase struct {
	Code           string `json:"code"`
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
	for code := range p23FailureContracts {
		actual = append(actual, code)
	}
	parsed := map[string]*ast.File{}
	for _, fixture := range cases {
		if fixture.OwnerTask < 3 || fixture.OwnerTask > 13 {
			t.Fatalf("unowned failure %s", fixture.Code)
		}
		expected = append(expected, fixture.Code)
		p22CheckBehavioralReference(t, fixture.BehavioralTest, parsed)
	}
	p22SameMembers(t, "accepted failure-code inventory", actual, expected)
	if len(cases) != 19 {
		t.Fatalf("enumerated Phase 23 failure codes = %d, want 19", len(cases))
	}
}

type p23LeafOutcome struct {
	Operation string
	Fields    []string
	Replay    bool
	Code      string
	Exit      int
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
	if len(cases) != 18 {
		t.Fatalf("Phase 23 leaf corpus has %d operations, want 18", len(cases))
	}
	root := t.TempDir()
	pairs := p23GenerateLeafPairs(t, integrationFIDBinary(t), root)
	byOperation := map[string]p23HumanLeafPair{}
	for _, pair := range pairs {
		byOperation[pair.Operation] = pair
	}
	if len(byOperation) != len(cases) {
		t.Fatalf("Human fixture has %d leaves, corpus has %d", len(byOperation), len(cases))
	}
	for _, fixture := range cases {
		fixture := fixture
		t.Run(fixture.Operation, func(t *testing.T) {
			pair, found := byOperation[fixture.Operation]
			if !found {
				t.Fatalf("no real invocation for %s", fixture.Operation)
			}
			if fixture.OwnerTask < 4 || fixture.OwnerTask > 13 {
				t.Fatalf("unowned %s", fixture.Operation)
			}
			fake := p23FakeLeafAdapter{fixture: fixture}
			fid := p23FIDLeafAdapter{pair: pair, binary: integrationFIDBinary(t)}
			p23CompareLeafOutcome(t, "accepted", fake.Accepted(), fid.Accepted(t))
			if fixture.Mutation {
				p23CompareLeafOutcome(t, "unknown-field", fake.UnknownField(), fid.UnknownField(t))
			}
		})
	}
	t.Run("authority", func(t *testing.T) {
		a := p22NewFIDRecoveryAdapter(t)
		before := p21JournalEventCount(t, a.workspace)
		response := p23RunInput(t, a, "requirement", "record-reference", "actor-not-activated", "p23-shared-authority", map[string]any{
			"id": "REQ-shared-authority", "source": a.measurement, "status": "active",
		})
		var failure struct {
			Operation string `json:"operation"`
			Error     struct {
				Code string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal([]byte(response.Stderr), &failure); err != nil {
			t.Fatal(err)
		}
		p23CompareLeafOutcome(t, "authority", p23LeafOutcome{Operation: "requirement.record-reference", Code: "invalid-argument", Exit: 2},
			p23LeafOutcome{Operation: failure.Operation, Code: failure.Error.Code, Exit: response.ExitCode})
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
				Code string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal([]byte(response.Stderr), &failure); err != nil {
			t.Fatal(err)
		}
		p23CompareLeafOutcome(t, "stale-version", p23LeafOutcome{Operation: "review.record", Code: "invalid-argument", Exit: 2},
			p23LeafOutcome{Operation: failure.Operation, Code: failure.Error.Code, Exit: response.ExitCode})
		if got := p21JournalEventCount(t, a.workspace); got != before {
			t.Fatalf("stale-version denial appended Event %d -> %d", before, got)
		}
	})
}

func TestP23SharedAmendedLeafConformance(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{},
		[]string{"claim.acquire", "envelope.activate", "envelope.create", "evidence.link", "gate.register", "gate.run", "run.observe-budget", "run.start"})
	p22AcceptedCall(t, a, "p23-shared-budget", "run.observe-budget", map[string]any{"value": 0})
	p22JourneyProof(t, a, "p23-shared-proof")
	cases := []struct {
		operation, file, key string
		fields               []string
		mutation             bool
	}{
		{"envelope.create", "p22-corpus-envelope", "p22-corpus-envelope-create", []string{"envelope", "idempotency_key", "replayed"}, true},
		{"envelope.activate", "p22-corpus-envelope-activate", "p22-corpus-envelope-activate", []string{"envelope", "evaluation", "idempotency_key", "replayed"}, true},
		{"claim.acquire", "p22-corpus-claim", "p22-corpus-claim-acquire", []string{"claim", "idempotency_key", "replayed"}, true},
		{"evidence.link", "p23-shared-proof-proof", "p23-shared-proof-proof", []string{"evidence", "evaluation", "idempotency_key", "replayed"}, true},
		{"evidence.evaluate", "p23-evaluate", "", []string{"evaluations"}, false},
	}
	p23EvidenceEvaluate(t, a)
	for _, fixture := range cases {
		t.Run(fixture.operation, func(t *testing.T) {
			parts := strings.Split(fixture.operation, ".")
			args := []string{"--project", a.workspace, "--json", parts[0], parts[1], "--input", filepath.Join(a.workspace, fixture.file+".json")}
			if fixture.mutation {
				args = append(args, "--actor", a.actor, "--idempotency-key", fixture.key)
			}
			fake := p23FakeLeafAdapter{fixture: p23LeafCase{Operation: fixture.operation, Fields: fixture.fields, Mutation: fixture.mutation}}
			fid := p23FIDLeafAdapter{pair: p23HumanLeafPair{Operation: fixture.operation, CWD: a.workspace, JSONArgs: args}, binary: a.binary}
			p23CompareLeafOutcome(t, "accepted", fake.Accepted(), fid.Accepted(t))
			if fixture.mutation {
				p23CompareLeafOutcome(t, "unknown-field", fake.UnknownField(), fid.UnknownField(t))
			}
		})
	}
}

type p23FakeLeafAdapter struct{ fixture p23LeafCase }

func (a p23FakeLeafAdapter) Accepted() p23LeafOutcome {
	return p23LeafOutcome{Operation: a.fixture.Operation, Fields: append([]string{}, a.fixture.Fields...), Replay: a.fixture.Mutation}
}

func (a p23FakeLeafAdapter) UnknownField() p23LeafOutcome {
	return p23LeafOutcome{Operation: a.fixture.Operation, Code: "invalid-argument", Exit: 2}
}

type p23FIDLeafAdapter struct {
	pair   p23HumanLeafPair
	binary string
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
	return p23LeafOutcome{Operation: envelope.Operation, Fields: fields, Replay: replay}
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

func (a p23FIDLeafAdapter) rejection(t *testing.T, args []string) p23LeafOutcome {
	t.Helper()
	before := p21JournalEventCount(t, a.pair.CWD)
	_, stderr, err := closureRunCommand(a.pair.CWD, a.binary, args...)
	var failure struct {
		Operation string `json:"operation"`
		Error     struct {
			Code string `json:"code"`
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
	return p23LeafOutcome{Operation: failure.Operation, Code: failure.Error.Code, Exit: exit.ExitCode()}
}

func p23CompareLeafOutcome(t *testing.T, name string, fake, fid p23LeafOutcome) {
	t.Helper()
	sort.Strings(fake.Fields)
	if !reflect.DeepEqual(fake, fid) {
		t.Fatalf("%s fake=%#v fid=%#v", name, fake, fid)
	}
}
