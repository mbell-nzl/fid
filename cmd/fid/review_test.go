package main

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

// TestReviewRecordRejectsNonIndependentActor is the Task 4 vertical tracer:
// recording under the mutating actor must be rejected before an Event is added.
func TestReviewRecordRejectsNonIndependentActor(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := runEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "review-independent-builder")
	repository := runRepositoryForTest(t, workspace)
	runPath := filepath.Join(workspace, "review-independent-run.json")
	writeActivationInput(t, runPath, runStartInput(issueID, envelopeID, claim["id"], repository))
	started := executeCommandForTest(workspace,
		"--project", workspace, "--json", "run", "start", "--actor", "actor-builder",
		"--input", runPath, "--idempotency-key", "review-independent-run",
	)
	if started.err != nil || started.stderr != "" {
		t.Fatalf("start mutating Run = %#v, want stdout-only success", started)
	}
	run := runFromResponse(t, started.stdout)
	inputPath := filepath.Join(workspace, "review-record.json")
	writeActivationInput(t, inputPath, reviewRecordInput(issueID, run, repository,
		map[string]any{"source_type": "external-reference", "display_label": "Independence probe", "path": nil, "section": nil,
			"external_id": "independence-probe", "digest": "sha256:" + strings.Repeat("0", 64)},
		"pass", []any{}, []any{}, []any{}))

	result := executeCommandForTest(workspace,
		"--project", workspace, "--json", "review", "record", "--actor", "actor-builder",
		"--input", inputPath, "--idempotency-key", "review-not-independent",
	)
	if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"review-not-independent"`) {
		t.Fatalf("review record by mutating actor = %#v, want typed review-not-independent rejection", result)
	}
}

func TestReviewRecordCreatesImmutableReview(t *testing.T) {
	t.Parallel()
	workspace, issueID, registration := finalizedObservedEnvelopeWorkspaceWithDefinitions(t, gateDefinitionsForTest(nil))
	input := envelopeCreateInput(t, workspace, issueID, registration)
	grant := input["authority_grants"].([]any)[0].(map[string]any)
	grant["actions"] = []string{"claim.acquire", "envelope.activate", "review.record", "run.start"}
	grant["actor_selector"] = map[string]any{"actor_id": "actor-builder", "required_capabilities": []string{"claim.acquire", "envelope.activate", "review.record", "run.start"}}
	path := filepath.Join(workspace, "review-envelope.json")
	writeActivationInput(t, path, input)
	created := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", path, "--idempotency-key", "review-envelope-create")
	if created.err != nil {
		t.Fatalf("create envelope = %#v", created)
	}
	envelopeID := assertEnvelopeShape(t, created.stdout, "draft")
	activatePath := filepath.Join(workspace, "review-envelope-activate.json")
	writeActivationInput(t, activatePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	activated := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "activate", "--actor", "actor-builder", "--input", activatePath, "--idempotency-key", "review-envelope-activate")
	if activated.err != nil {
		t.Fatalf("activate envelope = %#v", activated)
	}
	claim, _ := acquireGenericClaimForTestWithMode(t, workspace, issueID, envelopeID, "review-record-claim", 60, "review")
	repository := runRepositoryForTest(t, workspace)
	runPath := filepath.Join(workspace, "review-record-run.json")
	writeActivationInput(t, runPath, runStartInput(issueID, envelopeID, claim["id"], repository))
	started := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", runPath, "--idempotency-key", "review-record-run")
	if started.err != nil {
		t.Fatalf("start review Run = %#v", started)
	}
	run := runFromResponse(t, started.stdout)
	source := input["sources"].([]any)[0]
	recordPath := filepath.Join(workspace, "review-record-input.json")
	findings := []any{
		map[string]any{"id": "FND-review-one", "subject": map[string]any{"type": "issue", "id": issueID}, "scope": map[string]any{"paths": []any{}}, "severity": "minor", "summary": "first non-blocking finding", "blocking": false},
		map[string]any{"id": "FND-review-two", "subject": map[string]any{"type": "issue", "id": issueID}, "scope": map[string]any{"paths": []any{}}, "severity": "nit", "summary": "second non-blocking finding", "blocking": false},
	}
	recordInput := reviewRecordInput(issueID, run, repository, source, "pass", findings, []any{}, []any{})
	invalidVerdicts := []struct {
		name     string
		verdict  string
		findings []any
	}{
		{name: "pass-blocking", verdict: "pass", findings: []any{reviewFinding(issueID, "FND-pass-blocking", "major", true)}},
		{name: "changes-required-without-blocker", verdict: "changes-required", findings: []any{}},
		{name: "blocked-without-finding-or-exclusion", verdict: "blocked", findings: []any{}},
	}
	for _, test := range invalidVerdicts {
		invalid := cloneJSONMap(t, recordInput)
		invalid["verdict"] = test.verdict
		invalid["findings"] = test.findings
		writeActivationInput(t, recordPath, invalid)
		before := snapshotWorkspacePersistence(t, workspace)
		result := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "record", "--actor", "actor-builder", "--input", recordPath, "--idempotency-key", "review-invalid-"+test.name)
		if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"review-verdict-invalid"`) {
			t.Fatalf("invalid %s Review = %#v, want typed review-verdict-invalid rejection", test.name, result)
		}
		if after := snapshotWorkspacePersistence(t, workspace); !after.equal(before) {
			t.Fatalf("invalid %s Review appended an Event or changed the projection", test.name)
		}
	}
	writeActivationInput(t, recordPath, recordInput)
	staleInput := cloneJSONMap(t, recordInput)
	staleRepository := cloneJSONMap(t, recordInput["expected_repository"].(map[string]any))
	staleRepository["head_oid"] = strings.Repeat("0", 40)
	staleInput["expected_repository"] = staleRepository
	writeActivationInput(t, recordPath, staleInput)
	beforeMismatch := snapshotWorkspacePersistence(t, workspace)
	mismatch := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "record", "--actor", "actor-builder", "--input", recordPath, "--idempotency-key", "review-record-stale")
	if mismatch.err == nil || mismatch.stdout != "" || !strings.Contains(mismatch.stderr, `"code":"repository-mismatch"`) || !strings.Contains(mismatch.stderr, `"expected"`) || !strings.Contains(mismatch.stderr, `"observed"`) {
		t.Fatalf("record Review with stale repository = %#v, err=%v, want typed repository mismatch", mismatch, mismatch.err)
	}
	if afterMismatch := snapshotWorkspacePersistence(t, workspace); !afterMismatch.equal(beforeMismatch) {
		t.Fatal("repository mismatch appended an Event or changed the projection")
	}
	writeActivationInput(t, recordPath, recordInput)
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "record", "--actor", "actor-builder", "--input", recordPath, "--idempotency-key", "review-record")
	if result.err != nil || result.stderr != "" || !strings.Contains(result.stdout, `"verdict":"pass"`) {
		t.Fatalf("record Review = %#v, want immutable pass Review", result)
	}
	recorded := decodeCommandJSONStructure(t, result.stdout).(map[string]any)["result"].(map[string]any)["review"].(map[string]any)
	got := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "get", recorded["id"].(string))
	if got.err != nil || got.stderr != "" || !strings.Contains(got.stdout, `"verdict":"pass"`) {
		t.Fatalf("get recorded Review = %#v, want immutable pass Review", got)
	}
	if !strings.Contains(got.stdout, `"FND-review-one"`) || !strings.Contains(got.stdout, `"FND-review-two"`) {
		t.Fatalf("get recorded Review = %#v, want both immutable Findings", got)
	}
	if rebuilt := executeCommandForTest(workspace, "--project", workspace, "--json", "rebuild"); rebuilt.err != nil || rebuilt.stderr != "" {
		t.Fatalf("rebuild Review projection = %#v", rebuilt)
	}
	replayedProjection := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "get", recorded["id"].(string))
	if replayedProjection.err != nil || replayedProjection.stderr != "" || replayedProjection.stdout != got.stdout {
		t.Fatalf("review.get after rebuild = %#v, want byte-identical live projection %#v", replayedProjection, got)
	}
	recordInput["inapplicable_findings"] = []any{}
	writeActivationInput(t, recordPath, recordInput)
	replayed := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "record", "--actor", "actor-builder", "--input", recordPath, "--idempotency-key", "review-record")
	if replayed.err != nil || replayed.stderr != "" || !strings.Contains(replayed.stdout, `"replayed":true`) {
		t.Fatalf("replay Review = %#v, want original typed replay", replayed)
	}
	beforeDuplicate := snapshotWorkspacePersistence(t, workspace)
	duplicate := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "record", "--actor", "actor-builder", "--input", recordPath, "--idempotency-key", "review-record-duplicate")
	if duplicate.err == nil || duplicate.stdout != "" || !strings.Contains(duplicate.stderr, `"code":"duplicate-entity"`) {
		t.Fatalf("second Review for one Run = %#v, want typed duplicate entity conflict", duplicate)
	}
	if afterDuplicate := snapshotWorkspacePersistence(t, workspace); !afterDuplicate.equal(beforeDuplicate) {
		t.Fatal("second Review appended an Event or changed the projection")
	}
	conflictingInput := cloneJSONMap(t, recordInput)
	conflictingInput["profile"] = "review-other"
	writeActivationInput(t, recordPath, conflictingInput)
	beforeConflict := snapshotWorkspacePersistence(t, workspace)
	conflicting := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "record", "--actor", "actor-builder", "--input", recordPath, "--idempotency-key", "review-record")
	if conflicting.err == nil || conflicting.stdout != "" || !strings.Contains(conflicting.stderr, `"code":"idempotency-conflict"`) {
		t.Fatalf("conflicting Review replay = %#v, want typed idempotency conflict", conflicting)
	}
	if afterConflict := snapshotWorkspacePersistence(t, workspace); !afterConflict.equal(beforeConflict) {
		t.Fatal("idempotency conflict appended an Event or changed the projection")
	}
}

func TestReviewRecordIndependentReworkAndReceiptBinding(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID, registration, source := reviewWorkspaceWithTwoActivatedReviewers(t)
	repository := runRepositoryForTest(t, workspace)

	mutationClaim := acquireReviewClaim(t, workspace, issueID, envelopeID, "actor-builder", "mutation", "review-mutation")
	mutationRun := startReviewRun(t, workspace, issueID, envelopeID, mutationClaim, "actor-builder", "mutation")
	gatePath := filepath.Join(workspace, "review-mutation-gate.json")
	writeActivationInput(t, gatePath, map[string]any{"registration_id": registration.ID, "registration_version": registration.Version, "obligation_id": testEvidenceObligationID, "run_id": mutationRun["id"], "repository": repository})
	gate := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "run", "--actor", "actor-builder", "--input", gatePath, "--idempotency-key", "review-mutation-gate")
	if gate.err != nil || gate.stderr != "" {
		t.Fatalf("record Receipt on mutation Run = %#v", gate)
	}
	var terminal gateTerminalFrame
	lines := strings.Split(strings.TrimSpace(gate.stdout), "\n")
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &terminal); err != nil || terminal.Data.ReceiptID == "" {
		t.Fatalf("decode mutation Receipt: %v, output=%q", err, gate.stdout)
	}
	closeReviewRun(t, workspace, mutationRun, repository, "actor-builder", "review-mutation")
	releaseReviewClaim(t, workspace, mutationClaim, "actor-builder", "review-mutation")

	reviewClaim := acquireReviewClaim(t, workspace, issueID, envelopeID, "actor-reviewer-one", "review", "review-first")
	reviewRun := startReviewRun(t, workspace, issueID, envelopeID, reviewClaim, "actor-reviewer-one", "review-first")
	wrongReceipt := reviewRecordInput(issueID, reviewRun, repository, source, "changes-required", []any{reviewFinding(issueID, "FND-prior", "major", true)}, []any{}, []any{terminal.Data.ReceiptID})
	wrongReceiptPath := filepath.Join(workspace, "review-wrong-receipt.json")
	writeActivationInput(t, wrongReceiptPath, wrongReceipt)
	beforeReceipt := snapshotWorkspacePersistence(t, workspace)
	rejectedReceipt := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "record", "--actor", "actor-reviewer-one", "--input", wrongReceiptPath, "--idempotency-key", "review-wrong-receipt")
	if rejectedReceipt.err == nil || rejectedReceipt.stdout != "" || !strings.Contains(rejectedReceipt.stderr, `"code":"invalid-argument"`) {
		t.Fatalf("Review binding another Run Receipt = %#v, want typed no-Event rejection", rejectedReceipt)
	}
	if afterReceipt := snapshotWorkspacePersistence(t, workspace); !afterReceipt.equal(beforeReceipt) {
		t.Fatal("different-Run Receipt rejection appended an Event or changed the projection")
	}
	firstPath := filepath.Join(workspace, "review-first.json")
	firstInput := reviewRecordInput(issueID, reviewRun, repository, source, "changes-required", []any{reviewFinding(issueID, "FND-prior", "major", true), reviewFinding(issueID, "FND-prior-two", "major", true)}, []any{}, []any{})
	writeActivationInput(t, firstPath, firstInput)
	first := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "record", "--actor", "actor-reviewer-one", "--input", firstPath, "--idempotency-key", "review-first")
	if first.err != nil || first.stderr != "" {
		t.Fatalf("independent first Review = %#v", first)
	}
	firstReview := decodeCommandJSONStructure(t, first.stdout).(map[string]any)["result"].(map[string]any)["review"].(map[string]any)
	closeReviewRun(t, workspace, reviewRun, repository, "actor-reviewer-one", "review-first")
	releaseReviewClaim(t, workspace, reviewClaim, "actor-reviewer-one", "review-first")

	reworkClaim := acquireReviewClaim(t, workspace, issueID, envelopeID, "actor-builder", "mutation", "review-rework")
	reworkRun := startReviewRun(t, workspace, issueID, envelopeID, reworkClaim, "actor-builder", "review-rework")
	closeReviewRun(t, workspace, reworkRun, repository, "actor-builder", "review-rework")
	releaseReviewClaim(t, workspace, reworkClaim, "actor-builder", "review-rework")

	secondClaim := acquireReviewClaim(t, workspace, issueID, envelopeID, "actor-reviewer-one", "review", "review-second")
	secondRun := startReviewRun(t, workspace, issueID, envelopeID, secondClaim, "actor-reviewer-one", "review-second")
	secondPath := filepath.Join(workspace, "review-second.json")
	writeActivationInput(t, secondPath, reviewRecordInput(issueID, secondRun, repository, source, "pass", []any{}, []any{"FND-prior", "FND-prior-two"}, []any{}))
	second := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "record", "--actor", "actor-reviewer-one", "--input", secondPath, "--idempotency-key", "review-second")
	if second.err != nil || second.stderr != "" {
		t.Fatalf("same independent Reviewer after rework resolves prior Findings = %#v", second)
	}
	secondReview := decodeCommandJSONStructure(t, second.stdout).(map[string]any)["result"].(map[string]any)["review"].(map[string]any)
	wantResolved := []any{"FND-prior", "FND-prior-two"}
	if !reflect.DeepEqual(secondReview["resolved_finding_ids"], wantResolved) {
		t.Fatalf("resolving Review result IDs = %#v, want %#v", secondReview["resolved_finding_ids"], wantResolved)
	}
	assertAtomicReviewResolutionEvent(t, workspace, secondReview["id"].(string), secondReview["event_id"].(string), wantResolved)
	gotFirst := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "get", firstReview["id"].(string))
	gotSecond := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "get", secondReview["id"].(string))
	if gotFirst.err != nil || gotSecond.err != nil {
		t.Fatalf("review.get prior/resolving Reviews = first %#v second %#v, want both Findings atomically resolved through CLI", gotFirst, gotSecond)
	}
	assertResolvedReviewReadback(t, gotFirst.stdout, gotSecond.stdout, secondReview["id"].(string), wantResolved)
	if rebuilt := executeCommandForTest(workspace, "--project", workspace, "--json", "rebuild"); rebuilt.err != nil || rebuilt.stderr != "" {
		t.Fatalf("rebuild resolved Review projection = %#v", rebuilt)
	}
	rebuiltFirst := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "get", firstReview["id"].(string))
	rebuiltSecond := executeCommandForTest(workspace, "--project", workspace, "--json", "review", "get", secondReview["id"].(string))
	if rebuiltFirst.stdout != gotFirst.stdout || rebuiltSecond.stdout != gotSecond.stdout || rebuiltFirst.err != nil || rebuiltSecond.err != nil {
		t.Fatalf("resolved Review readback after rebuild = first %#v second %#v, want byte-identical live projections", rebuiltFirst, rebuiltSecond)
	}
}

func assertAtomicReviewResolutionEvent(t *testing.T, workspace, reviewID, eventID string, findingIDs []any) {
	t.Helper()
	event := mustReadLastJournalEvent(t, mustJournalPath(t, workspace))
	if event.EventType != "review.recorded" || event.EventID != eventID {
		t.Fatalf("resolving Event identity = %q/%q, want review.recorded/%q", event.EventType, event.EventID, eventID)
	}
	var dispositionEffects int
	for _, effect := range event.Effects {
		if effect.Op == "insert" && effect.Table == "finding_dispositions" {
			dispositionEffects++
			if effect.Values["qualifying_review_id"] != reviewID {
				t.Fatalf("resolution Effect qualifying Review = %#v, want %q", effect.Values["qualifying_review_id"], reviewID)
			}
		}
	}
	if dispositionEffects != len(findingIDs) {
		t.Fatalf("resolution disposition Effects = %d, want %d in one Event", dispositionEffects, len(findingIDs))
	}
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open Review projection: %v", err)
	}
	defer func() { _ = conn.Close() }()
	var projected int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM finding_dispositions WHERE qualifying_review_id=? AND action='resolved'`, reviewID).Scan(&projected); err != nil {
		t.Fatalf("count projected Finding dispositions: %v", err)
	}
	if projected != len(findingIDs) {
		t.Fatalf("projected Finding dispositions = %d, want %d", projected, len(findingIDs))
	}
}

func assertResolvedReviewReadback(t *testing.T, firstJSON, secondJSON, resolvingReviewID string, wantResolved []any) {
	t.Helper()
	first := decodeCommandJSONStructure(t, firstJSON).(map[string]any)["result"].(map[string]any)["review"].(map[string]any)
	second := decodeCommandJSONStructure(t, secondJSON).(map[string]any)["result"].(map[string]any)["review"].(map[string]any)
	if !reflect.DeepEqual(second["resolved_finding_ids"], wantResolved) {
		t.Fatalf("resolving Review readback IDs = %#v, want %#v", second["resolved_finding_ids"], wantResolved)
	}
	findings := first["findings"].([]any)
	if len(findings) != len(wantResolved) {
		t.Fatalf("prior Review Finding count = %d, want %d", len(findings), len(wantResolved))
	}
	for index, value := range findings {
		finding := value.(map[string]any)
		dispositions := finding["dispositions"].([]any)
		if finding["id"] != wantResolved[index] || finding["state"] != "resolved" || len(dispositions) != 1 || dispositions[0].(map[string]any)["qualifying_review_id"] != resolvingReviewID {
			t.Fatalf("resolved Finding %d = %#v, want one disposition from %q", index, finding, resolvingReviewID)
		}
	}
}

func reviewWorkspaceWithTwoActivatedReviewers(t *testing.T) (string, string, string, gateRegistrationView, any) {
	t.Helper()
	workspace := clonePristineWorkspace(t)
	created := executeCommandForTest(workspace, "--project", workspace, "--json", "issue", "create", "--kind", "task", "--title", "Review fixture", "--actor", "agent:developer")
	if created.err != nil {
		t.Fatalf("create fixture Issue: %#v", created)
	}
	issueID := decodeCommandJSONStructure(t, created.stdout).(map[string]any)["result"].(map[string]any)["issue"].(map[string]any)["id"].(string)
	activationID := prepareActivationForTest(t, workspace)
	profiles := []any{
		map[string]any{"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/builder", "session_ref": nil, "capabilities": authorityCapabilities()},
		map[string]any{"actor_id": "actor-harness", "actor_type": "service", "producer": "activated-harness", "provenance_ref": "activation:profile/harness", "session_ref": nil, "capabilities": []string{"run.observe"}},
		map[string]any{"actor_id": "actor-reviewer-one", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/reviewer-one", "session_ref": nil, "capabilities": authorityCapabilities()},
		map[string]any{"actor_id": "actor-reviewer-two", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/reviewer-two", "session_ref": nil, "capabilities": authorityCapabilities()},
	}
	profileDigest := stageActivationComponentForTest(t, workspace, activationID, "authority-profile", "reviewers", profiles)
	gateDigest := stageActivationComponentForTest(t, workspace, activationID, "gate-definitions", "review-gates", gateDefinitionsForTest(func(definition map[string]any) { definition["argv"] = []string{"/bin/true"} }))
	observeActivationResourcesForTest(t, workspace, activationID, "new")
	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	commitPath := filepath.Join(workspace, "review-activation-commit.json")
	writeActivationInput(t, commitPath, map[string]any{"activation_id": activationID, "expected_status": "staged", "expected_version": activationVersionForTest(t, workspace, activationID), "component_digests": digests})
	if committed := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "review-activation-commit"); committed.err != nil {
		t.Fatalf("commit review activation: %#v", committed)
	}
	finalPath := filepath.Join(workspace, "review-activation-finalize.json")
	writeActivationInput(t, finalPath, map[string]any{"activation_id": activationID, "expected_status": "committed", "expected_version": activationVersionForTest(t, workspace, activationID)})
	if finalized := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "finalize", "--input", finalPath, "--idempotency-key", "review-activation-finalize"); finalized.err != nil {
		t.Fatalf("finalize review activation: %#v", finalized)
	}
	initializeObservedRepositoryForTest(t, workspace)
	registrationPath := filepath.Join(workspace, "review-gate-register.json")
	registrationInput := gateRegisterInput(t, workspace)
	registrationInput["scope"] = map[string]any{"issue_id": issueID}
	writeActivationInput(t, registrationPath, registrationInput)
	registration := registerGateForTest(t, workspace, registrationPath, "review-gate-register")
	input := envelopeCreateInput(t, workspace, issueID, registration)
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "finding.address", "gate.run", "review.record", "run.close", "run.start"}
	input["authority_grants"] = []any{
		map[string]any{"actions": actions, "actor_selector": map[string]any{"actor_id": "actor-builder", "required_capabilities": actions}, "scope": map[string]any{"issue_id": issueID}},
		map[string]any{"actions": []string{"run.observe"}, "actor_selector": map[string]any{"actor_id": "actor-harness", "required_capabilities": []string{"run.observe"}}, "scope": map[string]any{"issue_id": issueID}},
		map[string]any{"actions": actions, "actor_selector": map[string]any{"actor_id": "actor-reviewer-one", "required_capabilities": actions}, "scope": map[string]any{"issue_id": issueID}},
		map[string]any{"actions": actions, "actor_selector": map[string]any{"actor_id": "actor-reviewer-two", "required_capabilities": actions}, "scope": map[string]any{"issue_id": issueID}},
	}
	envelopePath := filepath.Join(workspace, "review-two-actor-envelope.json")
	writeActivationInput(t, envelopePath, input)
	createdEnvelope := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", envelopePath, "--idempotency-key", "review-two-actor-envelope")
	if createdEnvelope.err != nil {
		t.Fatalf("create review Envelope: %#v", createdEnvelope)
	}
	envelopeID := assertEnvelopeShape(t, createdEnvelope.stdout, "draft")
	activatePath := filepath.Join(workspace, "review-two-actor-activate.json")
	writeActivationInput(t, activatePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	if activated := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "activate", "--actor", "actor-builder", "--input", activatePath, "--idempotency-key", "review-two-actor-activate"); activated.err != nil {
		t.Fatalf("activate review Envelope: %#v", activated)
	}
	return workspace, issueID, envelopeID, registration, input["sources"].([]any)[0]
}

func acquireReviewClaim(t *testing.T, workspace, issueID, envelopeID, actor, mode, key string) map[string]any {
	t.Helper()
	path := filepath.Join(workspace, key+"-claim.json")
	writeActivationInput(t, path, map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "worktree_id": runRepositoryForTest(t, workspace)["worktree_id"], "scope": map[string]any{"paths": []string{"cmd/fid/review.go"}}, "mode": mode, "ttl_seconds": 60})
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "acquire", "--actor", actor, "--input", path, "--idempotency-key", key+"-claim")
	if result.err != nil || result.stderr != "" {
		t.Fatalf("acquire %s Claim as %s = %#v", mode, actor, result)
	}
	return claimFromResponse(t, result.stdout)
}

func startReviewRun(t *testing.T, workspace, issueID, envelopeID string, claim map[string]any, actor, key string) map[string]any {
	t.Helper()
	path := filepath.Join(workspace, key+"-run.json")
	input := runStartInput(issueID, envelopeID, claim["id"], runRepositoryForTest(t, workspace))
	input["executor"] = map[string]any{"id": "executor-" + key, "source": "activated-harness"}
	writeActivationInput(t, path, input)
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", actor, "--input", path, "--idempotency-key", key+"-run")
	if result.err != nil || result.stderr != "" {
		t.Fatalf("start Run as %s = %#v", actor, result)
	}
	return runFromResponse(t, result.stdout)
}

func closeReviewRun(t *testing.T, workspace string, run map[string]any, repository map[string]any, actor, key string) {
	t.Helper()
	runID := run["id"].(string)
	current := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "get", runID)
	if current.err != nil {
		t.Fatalf("get Run before close = %#v", current)
	}
	run = runFromResponse(t, current.stdout)
	observePath := filepath.Join(workspace, key+"-observe.json")
	writeActivationInput(t, observePath, map[string]any{"run_id": runID, "expected_version": run["version"], "liveness": map[string]any{"executor_id": run["executor"].(map[string]any)["id"], "state": "terminated", "source_ref": "process:test", "observed_at": "2026-09-20T00:00:00Z"}})
	observed := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "observe", "--actor", "actor-harness", "--input", observePath, "--idempotency-key", key+"-observe")
	if observed.err != nil {
		t.Fatalf("observe Run = %#v", observed)
	}
	response := decodeCommandJSONStructure(t, observed.stdout).(map[string]any)["result"].(map[string]any)
	closePath := filepath.Join(workspace, key+"-close.json")
	writeActivationInput(t, closePath, map[string]any{"run_id": runID, "expected_version": response["run"].(map[string]any)["version"], "observation_id": response["observation"].(map[string]any)["id"], "repository": repository})
	closed := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "close", "--actor", actor, "--input", closePath, "--idempotency-key", key+"-close")
	if closed.err != nil {
		t.Fatalf("close Run as %s = %#v", actor, closed)
	}
}

func releaseReviewClaim(t *testing.T, workspace string, claim map[string]any, actor, key string) {
	t.Helper()
	path := filepath.Join(workspace, key+"-release.json")
	writeActivationInput(t, path, map[string]any{"claim_id": claim["id"], "expected_version": claim["version"], "reason": "completed"})
	released := executeCommandForTest(workspace, "--project", workspace, "--json", "claim", "release", "--actor", actor, "--input", path, "--idempotency-key", key+"-release")
	if released.err != nil || released.stderr != "" {
		t.Fatalf("release Claim as %s = %#v", actor, released)
	}
}

func reviewRecordInput(issueID string, run map[string]any, repository map[string]any, source any, verdict string, findings, resolves, receipts []any) map[string]any {
	sourceInput := map[string]any{}
	for key, value := range source.(map[string]any) {
		if key != "repository" {
			sourceInput[key] = value
		}
	}
	// The inherited envelope fixture uses a placeholder project-file Source.
	// Review SourceInput now checks files, so this retained coverage declaration
	// names the same digest as an external reference.
	sourceInput["source_type"] = "external-reference"
	sourceInput["path"] = nil
	sourceInput["external_id"] = "review-corpus-source"
	return map[string]any{"run_id": run["id"], "expected_run_version": run["version"], "profile": "review-default", "expected_repository": map[string]any{"head_oid": repository["head_oid"], "dirty": repository["dirty"]}, "coverage": map[string]any{"items": []any{map[string]any{"source": sourceInput, "assessed": true}}, "exclusions": []any{}}, "verdict": verdict, "findings": findings, "resolves_finding_ids": resolves, "evidence_refs": []any{}, "receipt_ids": receipts}
}

func reviewFinding(issueID, id, severity string, blocking bool) map[string]any {
	return map[string]any{"id": id, "subject": map[string]any{"type": "issue", "id": issueID}, "scope": map[string]any{"paths": []any{}}, "severity": severity, "summary": "review finding", "blocking": blocking}
}
