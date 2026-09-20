# Phase 23 - Task 6: Review evidence, review cycles, and rework exhaustion - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added the typed Evidence link wire for Receipt, Review, and Human Outcome subjects, with exactly one nullable subject binding, optional all-Obligation evaluation, and the accepted production type/trust admission classifier. Passing Reviews now link durably; `changes-required` and `blocked` Reviews reject with `review-changes-required` and sorted blocking Finding IDs.
- Reused one workflow-owned independence query at Review recording and Review Evidence link time. Claim history has precedence over Run history, query failures propagate, and a Review that linked while independent does not become retroactively invalid if its Actor later mutates.
- Added review Evidence evaluation. An unresolved blocking Finding on the same Issue, Envelope, and revision yields `finding-unresolved` even after a later passing Review links, while duplicate reason tuples are removed without disturbing stable precedence.
- Added the `review-cycles` fid-owned Budget unit and exact `review-rework-exhausted` authoring pair. A `changes-required` Review charges one delta keyed by its Review ID. Pass and blocked verdicts do not charge, replay and rebuild do not charge again, and external `run.observe-budget` cannot charge the unit.
- Made the last allowed Review, charge, Envelope Stop, and Attention one `review.recorded` Event. The durable Stop blocks exactly `claim.acquire` and `run.start`; both adapters return the accepted non-committing denial with the original Stop and Attention IDs, exhausted Budget, null Event, and `replayed:false`.
- Converged the strict current-Event interpretation for the accepted unit: known-zero Run projection, fid sample ownership, Review-to-Run sample binding, and actual Review result Budget/Stop/Attention projection. The `review.recorded` family remains outside the journal budget-fact switch, so the charge and exhaustion derivation are producer-side guarantees as Task 6 specifies.
- Added all four Task 6 CLI failure codes, the four new evaluation-reason precedence entries, the Stop/detail registries, all five authoring cases, compact domain matrices, retained receipt fixture convergence, and built-binary Review Evidence and cycle journeys.

## Files created or modified

- `internal/evidence/evidence.go`, `internal/evidence/evidence_admission.go`, and `internal/evidence/evidence_task6_test.go`
- `internal/review/review.go` and `internal/review/review_test.go`
- `internal/workflow/budget.go`, `internal/workflow/envelope.go`, `internal/workflow/evaluation.go`, `internal/workflow/run.go`, `internal/workflow/stop.go`, and matching tests
- `internal/claims/generic.go` and `internal/claims/generic_test.go`
- `internal/journal/current_event_budget_values.go`, `current_event_references.go`, `current_event_result_budgets.go`, `current_event_result_projection.go`, `current_event_value_catalog.go`, and `current_event_task6_values_test.go`
- `internal/cli/codes.go` and `internal/cli/codes_test.go`
- `cmd/fid/evidence.go`, `claim.go`, `review.go`, `stop.go`, their matching tests, and retained Evidence request fixtures
- `cmd/fid/integration/p23_review_evidence_budget_test.go` and the retained integration Evidence fixtures
- `project-docs/docs/phases/phase-23/tasks/task-6-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-23/tasks/task-6-done.md`

## Validation Contract Coverage

- VAL-008: PASS - `TestEvidenceEvaluateNullObligationReturnsActiveRevisionInIDOrder` exercises the production leaf with a null Obligation ID and returns every active-revision evaluation in Obligation ID order.
- VAL-024: PASS - `TestReviewAdmissionReturnsBlockingFindingIDs`, `TestLoadReviewEvidenceBindingEnforcesVerdictTypeAndLinkTimeIndependence`, and `TestP23ReviewEvidenceRealProcessJourney` prove a passing Review links durably, while `changes-required` and `blocked` reject with `review-changes-required`, exact sorted blocking IDs, and no Event.
- VAL-025: TASK 6 EVIDENCE PASS; FULL ASSERTION REMAINS PARTIAL UNTIL TASK 10 - the production classifier matrix and durable Receipt/Review paths cover the Task 6 allocation below. The two Human mismatch cells are pure production-classifier evidence over otherwise qualifying seeded values. The matching Human cell proves type compatibility only; no Human Outcome writer or positive durable Human link exists in this task.

| Evidence type | Obligation | Expected | Observed | Evidence route and producing test |
|---|---|---|---|---|
| receipt | gate-evidence | durable admission | durable admission | real Receipt link through `TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent` and retained built-command fixtures |
| receipt | review | `evidence-type-not-accepted` | `evidence-type-not-accepted` | production classifier, `TestEvidenceAdmissionTask6PureProductionClassifierMatrix` |
| receipt | human-outcome | `evidence-type-not-accepted` | `evidence-type-not-accepted` | production classifier, `TestEvidenceAdmissionTask6PureProductionClassifierMatrix` |
| review | gate-evidence | `evidence-type-not-accepted` | `evidence-type-not-accepted` | production classifier, `TestEvidenceAdmissionTask6PureProductionClassifierMatrix` |
| review | review | durable admission | durable admission | actual built binary, `TestP23ReviewEvidenceRealProcessJourney` |
| review | human-outcome | `evidence-type-not-accepted` | `evidence-type-not-accepted` | production classifier, `TestEvidenceAdmissionTask6PureProductionClassifierMatrix` |
| human-outcome | gate-evidence | `outcome-not-permitted` | `outcome-not-permitted` | pure intermediate production-classifier evidence, `TestEvidenceAdmissionTask6PureProductionClassifierMatrix` |
| human-outcome | review | `outcome-not-permitted` | `outcome-not-permitted` | pure intermediate production-classifier evidence, `TestEvidenceAdmissionTask6PureProductionClassifierMatrix` |
| human-outcome | human-outcome | type-compatible; durable link deferred | type-compatible only | pure intermediate production-classifier evidence, `TestEvidenceAdmissionTask6PureProductionClassifierMatrix`; Task 10 owns the real row and positive link |

- VAL-026: PASS - `TestP23ReviewEvidenceRealProcessJourney` records an unresolved blocking Finding, links a later passing Review, and observes `finding-unresolved` with the exact Finding ID. This closes Task 5 VAL-021's non-vacuous blocking proof.
- VAL-027: PASS for Task 6 producer-side enforcement - `TestChangesRequiredReviewChargesExactlyOneReviewCycle` proves one Review-keyed delta and no pass/blocked charge. `TestP23ReviewCyclesRealProcessJourney` proves the durable charge, rebuild, idempotent replay, and no second sample. The import leg remains assigned to Task 16 by the validation contract.
- VAL-028: PASS - `TestP23ReviewCyclesRealProcessJourney` invokes the built binary's `run.observe-budget` against `review-cycles`, receives a rejection, and proves the Event and sample counts remain unchanged.
- VAL-029: PASS - `TestP23ReviewCyclesRealProcessJourney` proves the exhausting Review, charge, Stop, and Attention are one Event, rebuild cleanly, replay without duplication, and leave later Claim and Run requests denied with no Event or entity.
- VAL-030: PASS - `TestEnvelopeAuthoringReviewCycleBudgetMatrix` exercises the matched pair, both unit/action mismatches, review cycles without a Review Obligation, and an unbudgeted Review Obligation. `TestEnvelopeAuthoringReviewCycleAlongsideRetainedBudgets` preserves ordinary Budget acceptance alongside the new pair.
- VAL-077: PASS - `TestP23ReviewEvidenceRealProcessJourney` links a passing Review once, gives its Actor a later mutating Claim, then proves the next link rejects `review-not-independent`, identifies the Claim, and appends no Event.
- VAL-079: PASS - the exact ordinary race gate reported `cmd/fid` at 535.185s and `internal/workflow` at 211.439s, both below 600s. The exact tagged integration rerun completed in 643.828s, below 900s.

## Tests run

- Required first RED: `go test ./internal/review -run '^TestChangesRequiredReviewChargesExactlyOneReviewCycle$' -count=1` - FAIL - `undefined: reviewCycleEffects` / `FAIL` (terminated `exit_code: 1`).
- Required first GREEN: the byte-identical command above - PASS - `ok  github.com/mbell-nzl/fid/internal/review` (terminated `exit_code: 0`).
- Focused domain and adapter convergence: `go test ./cmd/fid ./internal/review ./internal/workflow ./internal/claims ./internal/journal ./internal/cli ./internal/evidence -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/evidence 0.214s` (terminated `exit_code: 0`); preceding package tokens included `cmd/fid 45.265s`, `internal/workflow 22.930s`, and `internal/journal 10.063s`.
- Focused cycle real-process GREEN after strict projection and adapter repairs: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/... -run '^TestP23ReviewCyclesRealProcessJourney$' -count=1 -v` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 4.728s` (terminated `exit_code: 0`).
- Integration helper focused evidence, attributed to `/root/task4_process_journey`: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/... -run '^(TestP23ReviewEvidenceRealProcessJourney|TestP23ReviewCyclesRealProcessJourney)$'` - UNVERIFIED - reported `ok  github.com/mbell-nzl/fid/cmd/fid/integration 8.681s` with `exit_code: 0`; the committed helper change is `d05c288`.
- Reconstructed prechange conditional baseline under DEC-246: clean detached worktree `/tmp/fid-task6-baseline.N7uwYo` at Task 5 completion `8e6a76b5a4d083cf6b9dde76acf9ac3a5d169f0f`; `git status --short` was empty.
- `go test -race -timeout 30m -cover ./internal/evidence/` at that baseline - PASS - `ok  github.com/mbell-nzl/fid/internal/evidence 2.281s coverage: 23.6% of statements` (terminated `exit_code: 0`).
- `go test -race -timeout 30m -cover ./internal/workflow/` at that baseline - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 192.336s coverage: 64.1% of statements` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` - PASS after causal repair - `0 issues.` (terminated `exit_code: 0`). The first byte-identical attempt failed with `cmd/fid/evidence.go:74:52: ST1005: error strings should not be capitalized` and `exit_code: 1`.
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 211.439s coverage: 63.4% of statements` (terminated `exit_code: 0`).
- Ordinary package budgets: `cmd/fid` was 535.185s, `internal/workflow` was 211.439s, and every ordinary package remained below 600s.
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 64.2%` (terminated `exit_code: 0`).
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS after retained portable-Effect fixture repair - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 643.828s` (terminated `exit_code: 0`). The first byte-identical attempt failed with `unknown Effect column evidence_records.review_id` and `unknown Effect column evidence_records.human_outcome_id`; token `FAIL github.com/mbell-nzl/fid/cmd/fid/integration 351.564s` (terminated `exit_code: 1`).
- Tagged integration budget: `cmd/fid/integration` was 643.828s, below 900s.
- `go test -race -timeout 30m -cover ./internal/evidence/` - PASS - `ok  github.com/mbell-nzl/fid/internal/evidence 2.293s coverage: 25.7% of statements` (terminated `exit_code: 0`).
- `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 183.638s coverage: 63.4% of statements` (terminated `exit_code: 0`).

## Deviations from spec

- The required conditional before-edit coverage commands were missed chronologically. DEC-246 authorized reconstruction at the immutable Task 5 completion commit. Both byte-identical commands were personally run in a clean detached worktree, with the source commit and working directory recorded above. The after-edit commands remain ordinary exact gates.
- The first exact lint attempt found one capitalized internal error string in the new Evidence adapter. Commit `e861e1a7888ec870983246c8822644cc9aa83325` lowercased it; the byte-identical lint command then passed.
- The first exact tagged integration attempt found the retained portable-Effect interpreter had not admitted the accepted nullable Review and Human Outcome Evidence columns. The retained fixture was converged without changing product behavior or the gate, then the byte-identical full command was rerun.
- During focused integration, the initial Task 6 authoring fixture supplied output-only Budget `state`, then used a Review profile that differed from its Obligation. Those fixtures were corrected to the accepted input. The product tracer then exposed four real strict-path gaps: known-zero projection for `review-cycles`, fid source ownership and Review-to-Run sample reference, actual Review Budget/Stop/Attention result projection, and numeric Stop-detail validation. DEC-247 authorized only those four journal convergence paths; no Event catalog, result grammar, or budget-fact switch changed.
- DEC-244 added locked generic Claim admission to the task scope, and DEC-245 added the existing Review CLI result adapter. These corrections implemented the already accepted VAL-029 behavior without changing a validation assertion or public contract.
- No product, schema, protected gate, threshold, validation-contract, or dependency deviation was made.

## Concerns, technical debt, or notes for the Tech Lead

- Full VAL-025 remains partial until Task 10 records a real Human Outcome and proves its positive durable link. Task 6 supplies all six mismatch classifications and type compatibility without adding a temporary Human writer.
- VAL-027's import no-double-charge leg remains Task 16, when the archive profile covers the Phase 23 tables.
- Producer code enforces the review-cycle charge and exhaustion derivation. As accepted in Task 6, `review.recorded` remains outside `validateCurrentBudgetFacts`; the report does not claim a journal-layer automatic-charge proof.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `0cb705c`

## Commit subject

- `feat(review): enforce evidence rework cycles`
