# Phase 23 - Task 6: Review evidence, review cycles, and rework exhaustion - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
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
- QA repair 2026-09-26 (review-7 MF-4, DEC-252/DEC-256): added the VAL-027 linked-recovery evidence (test-only) and repaired `validateGateRunReferences` so `gate.run` is authorized against every Envelope Grant. See the two dated repair sections below.

## Files created or modified

- `internal/evidence/evidence.go`, `internal/evidence/evidence_admission.go`, and `internal/evidence/evidence_task6_test.go`
- `internal/review/review.go` and `internal/review/review_test.go`
- `internal/workflow/budget.go`, `internal/workflow/envelope.go`, `internal/workflow/evaluation.go`, `internal/workflow/run.go`, `internal/workflow/stop.go`, and matching tests
- `internal/claims/generic.go` and `internal/claims/generic_test.go`
- `internal/journal/current_event_budget_values.go`, `current_event_references.go`, `current_event_result_budgets.go`, `current_event_result_projection.go`, `current_event_value_catalog.go`, and `current_event_task6_values_test.go`
- `internal/cli/codes.go` and `internal/cli/codes_test.go`
- `cmd/fid/evidence.go`, `claim.go`, `review.go`, `stop.go`, their matching tests, and retained Evidence request fixtures
- `cmd/fid/integration/p23_review_evidence_budget_test.go` and the retained integration Evidence fixtures
- QA repair 2026-09-26: `cmd/fid/integration/p23_review_evidence_budget_test.go` (MF-4), `cmd/fid/integration/p23_gate_grants_test.go` (new, I-5), `internal/evidence/execution.go` (I-5), `internal/evidence/execution_gate_authority_test.go` (new, I-5)
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
- VAL-027: PASS for Task 6 producer-side enforcement - `TestChangesRequiredReviewChargesExactlyOneReviewCycle` proves one Review-keyed delta and no pass/blocked charge. `TestP23ReviewCyclesRealProcessJourney` proves the durable charge, rebuild, idempotent replay, and no second sample. The import leg remains assigned to Task 16 by the validation contract. Repair note 2026-09-26 (review-7 MF-4): the linked-recovery leg was unevidenced above and is now proven by `TestP23ReviewCycleChargeSurvivesLinkedRecoveryAndRebuild` (`cmd/fid/integration/p23_review_evidence_budget_test.go`); see "QA repair MF-4" below. The earlier claim that "recording, rebuild, idempotent replay" covered the assertion did not include recovery and stands corrected.
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
- `go test -race -timeout 30m -cover ./internal/evidence/` - PASS - `ok  github.com/mbell-nzl/fid/internal/evidence 2.281s coverage: 23.6% of statements` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 192.336s coverage: 64.1% of statements` (terminated `exit_code: 0`)

- QA repair 2026-09-26 focused runs are recorded in the dated sections "QA repair MF-4" and "QA repair I-5" below, and the repair's exact gates in "QA repair 2026-09-26 gate results".

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- First lint attempt - FAIL - `cmd/fid/evidence.go:74:52: ST1005: error strings should not be capitalized` (terminated `exit_code: 1`); commit `e861e1a` made the causal one-line repair before the byte-identical passing attempt.
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 211.439s coverage: 63.4% of statements` (terminated `exit_code: 0`)
- Ordinary package budgets: `cmd/fid` was 535.185s, `internal/workflow` was 211.439s, and every ordinary package remained below 600s.
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 64.2%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 643.828s` (terminated `exit_code: 0`)
- First tagged integration attempt - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid/integration 351.564s` (terminated `exit_code: 1`); the retained interpreter rejected `evidence_records.review_id` and `human_outcome_id`, and commit `0e0f2aa` converged that accepted fixture before the byte-identical passing attempt.
- Tagged integration budget: `cmd/fid/integration` was 643.828s, below 900s.
- `go test -race -timeout 30m -cover ./internal/evidence/` - PASS - `ok  github.com/mbell-nzl/fid/internal/evidence 2.293s coverage: 25.7% of statements` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 183.638s coverage: 63.4% of statements` (terminated `exit_code: 0`)

### QA repair 2026-09-26 gate results (run serially in this recovery session against implementation commit `8531522e9602b0f72717d6c9d889684228cfa2ee`, byte-identical spec commands, no `-count=1`)

These rows are this recovery's runs; the rows above are the original session's history and are preserved.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`); the same run printed `ok  	github.com/mbell-nzl/fid/cmd/fid	478.295s	coverage: 69.3% of statements`, `ok  	github.com/mbell-nzl/fid/internal/evidence	2.350s	coverage: 29.2% of statements`, and `ok  	github.com/mbell-nzl/fid/internal/export	232.335s`; no package line failed.
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					61.6%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	681.447s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -cover ./internal/evidence/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/evidence	2.338s	coverage: 29.2% of statements` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	185.338s	coverage: 62.1% of statements` (terminated `exit_code: 0`)

DEC-227 package wall times (VAL-079): `cmd/fid/integration` 681.447 s (tagged, budget 900 s, PASS); `cmd/fid` 478.295 s (ordinary, budget 600 s, PASS); `internal/workflow` 185.338 s (ordinary, budget 600 s, PASS; taken from the dedicated `./internal/workflow/` command because the repository-wide run printed `(cached)` for that unchanged package). No ordinary package exceeded 600 s; the repository-wide run's slowest packages were `cmd/fid` 478.295 s and `internal/export` 232.335 s. `goimports -l internal/evidence cmd/fid/integration` printed nothing (checked before the gates; exit status not separately recorded). Each gate's `exit_code` was captured by the launching shell as the command terminated, not asserted.

## QA repair MF-4 (review-7, VAL-027 linked-recovery leg)

Dated 2026-09-26. Trigger `review-7:qa-rework:task-6-MF-4-I-5`; authority: `project-docs/docs/phases/phase-23/qa/checkpoint-after-task-11-2.md` Repair Plan row MF-4, and DEC-256 (one recovery covers MF-4 and I-5). Test-only: no contract, schema, Event, gate, or product file changed.

Finding restated: the validation contract's VAL-027 method includes "after linked recovery at task 6", and no test recorded a `changes-required` Review and then ran `run.recover`. The product behaved correctly (review-7's scratch probe); the committed evidence was missing.

Evidence added: `TestP23ReviewCycleChargeSurvivesLinkedRecoveryAndRebuild` in `cmd/fid/integration/p23_review_evidence_budget_test.go` (built binary, tagged package; not in `cmd/fid`, whose ordinary time budget is nearly used). It reuses the unchanged `p22Observe`, `p22Expire`, `p22Recover` helpers and the review-Run recovery pattern of `p23_claim_mode_test.go`. The journey:

1. Builds a review-mode Claim and Run (`review-cycles` maximum 3, so no exhaustion Stop blocks recovery) and records a `changes-required` Review in it.
2. Asserts exactly one `review-cycles` `budget_samples` row: `source = fid`, `source_sample_id` equal to the Review ID, `value = 1` (consumed 1), bound to the review Run.
3. Observes the Run terminated, expires its Claim, and runs `run.recover` into a linked successor (`result.recovery.prior_run_id` is the review Run and the successor Run ID differs).
4. Re-reads the rows and compares them with the pre-recovery snapshot through `p23AssertReviewCyclesUnchanged`: first row count and consumed sum, then every row field-by-field, so a second charge, a removed charge, or a re-keyed charge all fail. Then runs `fid rebuild` and compares again.

Observed result (rows on the passing run: one sample, count 1, consumed 1, unchanged after recovery and after rebuild; no second `review-cycles` row): `--- PASS: TestP23ReviewCycleChargeSurvivesLinkedRecoveryAndRebuild (1.71s)` / `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.207s` (terminated `exit_code: 0`). The product was already correct, so the new test could not be red before an edit; the mutation checks below show the assertions bite. No double charge was found, so no `internal/workflow` file was touched.

Mutation checks (each reverted immediately; `git status --short` afterwards showed no product-file change):

- Double-charge attempt: in `internal/workflow/recovery_budget.go`, the `force` guard was disabled and the unit filter changed to `review-cycles`, making recovery charge a review-cycles sample. Result: `run.recover` was refused by the strict journal (`invalid-argument`, exit 2), so the test failed at `run.recover of the review Run ... want success` (terminated `exit_code: 1`). The double-charge path is fail-closed and cannot append a second sample; this proves the journey depends on recovery succeeding, not that the count assertion fired.
- Count/consumed assertion: a test-side mutation inserted a second `review-cycles` row between recovery and the comparison. Result: `review-cycles after linked recovery = count 2 consumed 2, want count 1 consumed 1` (terminated `exit_code: 1`). The reverted test file was restored byte-identical from a saved copy.
- Key assertion: in `internal/review/review.go` the charge's `SourceSampleID` was changed from the Review ID to the Run ID. Result: `review-cycles samples after recording = [...SourceSampleID:"RUN-190e93"...], want exactly one fid sample keyed on Review "REV-3e9657" with consumed 1 on Run "RUN-190e93"` (terminated `exit_code: 1`).
- Not mutated: a rebuild-only double charge. A rebuild replays the journal, so a database-only insert would be erased before the comparison and a journal-side mutation would need a second Event; the rebuild leg therefore relies on the same comparison helper as the recovery leg, and this report does not claim a separate rebuild mutation.

VAL-027 count/consumed comparison for this repair: 1 sample / consumed 1 after recording, 1 / 1 after `run.recover`, 1 / 1 after `fid rebuild`. The import leg remains Task 16.

## QA repair I-5 (review-7 Required Judgment 1, DEC-252)

Dated 2026-09-26. Authority: DEC-252, as ordered by DEC-256.

Defect: `validateGateRunReferences` (`internal/evidence/execution.go`) read a single unordered `envelope_authority_grants` row and required it to name the running Actor. In a multi-grant Envelope `gate.run` therefore worked only for the Actor whose Grant sorted first. It failed closed. It contradicted Phase 21 contracts section 6 ("exact action, capability, Issue, Envelope, and revision Grant match; missing is deny").

Provenance (verified in this session, not transcribed):

- `git log --format='%h %ad %s%n%b' --date=short -S'authority grant does not permit gate.run' -- internal/evidence/execution.go` printed `b694991 2026-08-30 feat(evidence): persist gate attempts and receipts` with body `Refs: phase-21/task-15`.
- `git blame -w -L '/var actionsJSON, selectorJSON, scopeJSON string/,+3' b694991 -- internal/evidence/execution.go` shows the single-row `conn.QueryRow(... envelope_authority_grants ...).Scan(...)` at lines 261-263 authored in `b6949911` on 2026-08-30.
- `git merge-base --is-ancestor b694991 HEAD` printed `ancestor`; Phase 23 planning began on 2026-09-14 (`34bd90e`, first commit under `project-docs/docs/phases/phase-23`, and `f94e354 2026-09-14 docs(phase-23): plan review and human judgment`). The defect predates phase start by fifteen days and is not a Phase 23 regression.
- Existing workaround, not edited: `cmd/fid/integration/p23_human_outcome_test.go` `p23RunQualifyingGateSession` runs the gate as the first-sorted Actor `actor-denied`; DEC-256 confirms it still passes, and it did in the full tagged run above.

Repair: `validateGateRunReferences` now delegates to `authorizeGateRunGrants`, which reads every Grant of the Envelope revision (`ORDER BY ordinal`) and accepts when any one names the running Actor with the exact `gate.run` action, `gate.run` required capability, and the Run's Issue, all for the Envelope and revision already bound to the Run and Obligation. This is the same action, capability, Issue, Envelope, and revision match applied by `internal/policy` `EvaluateEnvelopeGrant`, restricted to the checks the function already made. The failure detail is unchanged (`run gate: authority grant does not permit gate.run`); a revision with no Grants returns the previous `run gate: load authority grant: sql: no rows in result set` wrapping; query and scan errors keep the previous `load authority grant` wrapping; no other denial path was altered.

Red/green:

- In-package regression `TestValidateGateRunReferencesAuthorizesAgainstEveryEnvelopeGrant` (`internal/evidence/execution_gate_authority_test.go`): five Grants on one Envelope revision, with `actor-a` and `actor-b` granted, `actor-c` lacking the action, `actor-d` granted for a different Issue, and `actor-e` lacking the required capability. Red before the fix: `--- FAIL: .../actor-b ... granted Actor actor-b = run gate: authority grant does not permit gate.run` / `FAIL	github.com/mbell-nzl/fid/internal/evidence	0.007s` (terminated `exit_code: 1`). Green after: `--- PASS: TestValidateGateRunReferencesAuthorizesAgainstEveryEnvelopeGrant` with all five subtests passing (terminated `exit_code: 0`).
- Built-binary journey `TestP23GateRunAuthorizesEveryEnvelopeGrantThroughBuiltBinary` (`cmd/fid/integration/p23_gate_grants_test.go`), one active Envelope for a fresh Issue with three Grants: `actor-denied` (first sorted) and `actor-reviewer-two` (last sorted) granted `gate.run`; `actor-recovery` (between them) holds the `gate.run` capability in its activated profile but its Grant omits `gate.run`. Each Actor acquires its own mutation Claim and starts its own Run through the CLI, then runs `fid gate run`. Observed: `actor-denied` and `actor-reviewer-two` each receive an `fid-executed` `succeeded` Receipt naming their own Actor and Run; `actor-recovery` is denied with the adapter's `invalid-argument` (exit 2, empty stdout) and the Event count and a projection summary (row counts of `execution_attempts`, `execution_receipts`, `evidence_records`, `budget_samples`, `idempotency_results`, `run_budget_bindings`, plus every Run's id and status) are identical before and after. An in-process `evidence.RunGate` call as `actor-recovery` returns exactly `run gate: authority grant does not permit gate.run`, pinning the cause to the Grant match, again with unchanged Event count and projection. Red before the fix (the same journey, fix reverted): `actor-reviewer-two gate.run with its own gate.run Grant = ... invalid-argument ... ExitCode:2, want stdout-only success` / `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration	2.303s` (terminated `exit_code: 1`). Green after: `--- PASS: TestP23GateRunAuthorizesEveryEnvelopeGrantThroughBuiltBinary (3.05s)` / `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.535s` (terminated `exit_code: 0`).

Scope note on the ungranted leg: an Actor with no `gate.run` Grant cannot own an active Run through the CLI in the standard fixture (its `claim.acquire` and `run.start` would be denied), and a Run owned by another Actor is rejected earlier by the unchanged "active Run binding is invalid" check. The journey therefore builds one Envelope in which the ungranted Actor holds `claim.acquire` and `run.start` Grants but not `gate.run`, which is the only route that reaches the Grant check through the built binary. The adapter maps the unchanged denial to `invalid-argument`; DEC-252 requires no code or detail change, and none was made.

## Deviations from spec

- The first post-commit report-lint attempt terminated `exit_code: 11` with `FAIL report-lint found 11 failure(s)`. Required-command statuses and termination receipts were reformatted without changing their evidence; the byte-identical control then passed.
- The required conditional before-edit coverage commands were missed chronologically. DEC-246 authorized reconstruction at the immutable Task 5 completion commit. Both byte-identical commands were personally run in a clean detached worktree, with the source commit and working directory recorded above. The after-edit commands remain ordinary exact gates.
- The first exact lint attempt found one capitalized internal error string in the new Evidence adapter. Commit `e861e1a7888ec870983246c8822644cc9aa83325` lowercased it; the byte-identical lint command then passed.
- The first exact tagged integration attempt found the retained portable-Effect interpreter had not admitted the accepted nullable Review and Human Outcome Evidence columns. The retained fixture was converged without changing product behavior or the gate, then the byte-identical full command was rerun.
- During focused integration, the initial Task 6 authoring fixture supplied output-only Budget `state`, then used a Review profile that differed from its Obligation. Those fixtures were corrected to the accepted input. The product tracer then exposed four real strict-path gaps: known-zero projection for `review-cycles`, fid source ownership and Review-to-Run sample reference, actual Review Budget/Stop/Attention result projection, and numeric Stop-detail validation. DEC-247 authorized only those four journal convergence paths; no Event catalog, result grammar, or budget-fact switch changed.
- DEC-244 added locked generic Claim admission to the task scope, and DEC-245 added the existing Review CLI result adapter. These corrections implemented the already accepted VAL-029 behavior without changing a validation assertion or public contract.
- No product, schema, protected gate, threshold, validation-contract, or dependency deviation was made.
- QA repair 2026-09-26: DEC-256 folds I-5 into the MF-4 recovery, superseding the `task-6.md` amendment wording that calls them separate dispatches; both ran in this one recovery session. All seven required commands passed on the first attempt in this recovery, so there is no failed-attempt deviation for it.
- The repository-wide `go test -race -timeout 30m -coverprofile=coverage.out ./...` printed `(cached)` for the unchanged `internal/workflow` package. The DEC-227 `internal/workflow` time was taken from the dedicated `go test -race -timeout 30m -cover ./internal/workflow/` command in the same session, which executed (185.338s). `cmd/fid`, `internal/evidence`, and `internal/export` executed uncached in the repository-wide run.

## Concerns, technical debt, or notes for the Tech Lead

- Full VAL-025 remains partial until Task 10 records a real Human Outcome and proves its positive durable link. Task 6 supplies all six mismatch classifications and type compatibility without adding a temporary Human writer.
- VAL-027's import no-double-charge leg remains Task 16, when the archive profile covers the Phase 23 tables.
- Producer code enforces the review-cycle charge and exhaustion derivation. As accepted in Task 6, `review.recorded` remains outside `validateCurrentBudgetFacts`; the report does not claim a journal-layer automatic-charge proof.

- I-5: the task-10 Human Outcome journey still runs its gate as the first-sorted Actor. It was left alone as instructed; MF-5 may drop that workaround. Task 14 and later fixtures need not route around the defect.
- I-5: the ungranted-Actor denial through the built binary surfaces as `invalid-argument` (exit 2) because the existing gate adapter maps every non-authority reference error that way; the message is pinned in-process. If the Tech Lead wants `authority-denied` for this path, that is a code change and outside DEC-252.
- MF-4 shows the double-charge path is fail-closed at the strict journal, so nothing in `internal/workflow/{budget,run_budget,recovery_budget}.go` needed a repair.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `8531522`

## Commit subject

- `fix(evidence): authorize gate.run against every Envelope grant`
