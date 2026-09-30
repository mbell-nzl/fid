//go:build integration

package integration

import (
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
)

type p23LeafCase struct {
	Operation  string         `json:"operation"`
	Fields     []string       `json:"fields"`
	Mutation   bool           `json:"mutation"`
	OwnerTask  int            `json:"owner_task"`
	Semantics  map[string]any `json:"semantics,omitempty"`
	Amended    bool           `json:"amended,omitempty"`
	NotFoundID string         `json:"not_found_id,omitempty"`
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
			}
			if fixture.NotFoundID != "" {
				p23CompareLeafOutcome(t, "not-found", fake.NotFound(), fid.NotFound(t))
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
			"obligation_id": "OBL-p22-corpus", "repository": a.repository,
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
	return p23LeafOutcome{Operation: a.fixture.Operation, Code: "invalid-argument", Category: "invalid-input", Exit: 2,
		Details: map[string]any{"argument": "input", "value": nil, "allowed": []any{"accepted " + a.fixture.Operation + " input"}}}
}

func (a p23FakeLeafAdapter) Conflict(pair p23HumanLeafPair) p23LeafOutcome {
	key := ""
	for i, arg := range pair.JSONArgs {
		if arg == "--idempotency-key" && i+1 < len(pair.JSONArgs) {
			key = pair.JSONArgs[i+1]
		}
	}
	return p23LeafOutcome{Operation: a.fixture.Operation, Code: "idempotency-conflict", Category: "conflict", Exit: 4,
		Details: map[string]any{"operation_scope": a.fixture.Operation + ":<generated:project>", "idempotency_key": key,
			"stored_request_digest": "<generated:digest>", "supplied_request_digest": "<generated:digest>"}}
}

func (a p23FakeLeafAdapter) NotFound() p23LeafOutcome {
	code, subjectType := "", ""
	switch a.fixture.Operation {
	case "review.get":
		code, subjectType = "review-not-found", "review"
	case "requirement.get":
		code, subjectType = "requirement-not-found", "requirement"
	}
	return p23LeafOutcome{Operation: a.fixture.Operation, Code: code, Category: "not-found", Exit: 3,
		Subject: map[string]any{"type": subjectType, "id": a.fixture.NotFoundID},
		Details: map[string]any{"id": a.fixture.NotFoundID}}
}

type p23FIDLeafAdapter struct {
	pair               p23HumanLeafPair
	binary             string
	fixture            p23LeafCase
	alternateReceiptID string
}

func (a p23FIDLeafAdapter) NotFound(t *testing.T) p23LeafOutcome {
	t.Helper()
	args := append([]string{}, a.pair.JSONArgs...)
	if len(args) == 0 || a.fixture.NotFoundID == "" {
		t.Fatal("missing ID-taking read case")
	}
	args[len(args)-1] = a.fixture.NotFoundID
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
	fake.Semantics = p23NormalizeLeafValues(t, fake.Semantics)
	fid.Semantics = p23NormalizeLeafValues(t, fid.Semantics)
	fake.Details = p23NormalizeLeafValues(t, fake.Details)
	fid.Details = p23NormalizeLeafValues(t, fid.Details)
	if !reflect.DeepEqual(fake, fid) {
		t.Fatalf("%s fake=%#v fid=%#v", name, fake, fid)
	}
}

// The only normalized result members here are the generated Requirement row ID,
// EventRef, and creation time. Both adapters pass through this same traversal.
func p23NormalizeLeafValues(t *testing.T, raw map[string]any) map[string]any {
	t.Helper()
	if raw == nil {
		return nil
	}
	return p23NormalizeLeafValue(t, "", raw).(map[string]any)
}

func p23NormalizeLeafValue(t *testing.T, key string, raw any) any {
	t.Helper()
	switch value := raw.(type) {
	case map[string]any:
		out := make(map[string]any, len(value))
		for child, item := range value {
			out[child] = p23NormalizeLeafValue(t, child, item)
		}
		return out
	case []any:
		out := make([]any, len(value))
		for i, item := range value {
			out[i] = p23NormalizeLeafValue(t, key, item)
		}
		return out
	case string:
		switch key {
		case "operation_scope":
			parts := strings.SplitN(value, ":", 2)
			if len(parts) == 2 && (parts[1] == "<generated:project>" || p23ValidGeneratedToken("project", parts[1])) {
				return parts[0] + ":<generated:project>"
			}
		case "stored_request_digest", "supplied_request_digest":
			if value != "<generated:digest>" && !p23ValidGeneratedToken("digest", value) {
				t.Fatalf("invalid generated request digest %q", value)
			}
			return "<generated:digest>"
		case "created_at", "disposed_at", "recorded_at", "acquired_at", "expires_at", "activated_at":
			if value != "<generated:timestamp>" && !p23ValidGeneratedToken("timestamp", value) {
				t.Fatalf("invalid generated timestamp %q", value)
			}
			return "<generated:timestamp>"
		case "event_id":
			if value != "<generated:EventRef>" && !p23ValidGeneratedToken("EventRef", value) {
				t.Fatalf("invalid generated EventRef %q", value)
			}
			return "<generated:EventRef>"
		case "id":
			if value == "<generated:RREF>" || p23ValidGeneratedToken("RREF", value) {
				return "<generated:RREF>"
			}
			if value == "<generated:REV>" || p23ValidGeneratedToken("REV", value) {
				return "<generated:REV>"
			}
			if value == "<generated:FDP>" || p23ValidGeneratedToken("FDP", value) {
				return "<generated:FDP>"
			}
			if value == "<generated:CLM>" || p23ValidGeneratedToken("CLM", value) {
				return "<generated:CLM>"
			}
			if value == "<generated:ENV>" || p23ValidGeneratedToken("ENV", value) {
				return "<generated:ENV>"
			}
			if value == "<generated:EVD>" || p23ValidGeneratedToken("EVD", value) {
				return "<generated:EVD>"
			}
		case "review_id":
			if value == "<generated:REV>" || p23ValidGeneratedToken("REV", value) {
				return "<generated:REV>"
			}
		case "envelope_id":
			if value == "<generated:ENV>" || p23ValidGeneratedToken("ENV", value) {
				return "<generated:ENV>"
			}
		case "activation_id":
			if value == "<generated:ACT>" || p23ValidGeneratedToken("ACT", value) {
				return "<generated:ACT>"
			}
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
