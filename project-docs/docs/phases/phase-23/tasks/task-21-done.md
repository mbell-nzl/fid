# Phase 23 - Task 21: Resolvable review-rework exhaustion - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: Scope gap: the required shared-fake mirror needs `cmd/fid/integration/p22_semantic_fake_test.go` and `cmd/fid/integration/p22_semantic_denial_test.go`, neither of which is in task 21's allowed file list. The scoped implementation and required gates are green, but the fake does not model review-cycle exhaustion.
- Recovery attempts: 1
- Next action: Phase Runner obtains a bounded same-task scope amendment for the two named fake files, then dispatches a task-21 continuation to add fake condition/admission conformance, rerun the exact gate, and close this report. Do not advance to another task.
- Human attention required: no

## Summary of what was implemented

- Review-cycle exhaustion now raises a `human-outcome-recorded` Stop conditioned on an approved `review-rework-extension` Outcome for the Envelope. The existing DEC-250 predicate handles atomic and manual resolution.
- One shared helper excludes rework budgets from generic gate admission and resume reasons. The Phase 22 admission function remains available for other budgets.
- Human denial adds the exact guidance line after task 19's error line; JSON behavior is unchanged.
- Built-binary journeys cover wrong class, rejected result, another Issue, non-human authority denial, qualifying atomic and manual resolution, admitted retry, and two later single-charge Reviews.

## Files created or modified

- `internal/workflow/budget.go`; `internal/workflow/resume_actions.go`; `internal/workflow/budget_test.go`
- `internal/evidence/execution.go`; `internal/evidence/execution_test.go`
- `internal/claims/generic_test.go`
- `cmd/fid/stop.go`; `cmd/fid/integration/p23_review_evidence_budget_test.go`; `cmd/fid/integration/p23_human_fixtures_test.go`; `cmd/fid/integration/p23_human_failure_test.go`
- `project-docs/docs/phases/phase-23/tasks/task-21-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-23/qa/human-fixtures.md` rework section adds the resolution and admitted retry instructions. Its prior verbatim snapshot is identified as historical; Phase-end QA regenerates it.

## Validation Contract Coverage

- VAL-086: PASS for the real-process behavior in `TestP23ReviewCyclesRealProcessJourney` (`cmd/fid/integration/p23_review_evidence_budget_test.go:77`) and the separate manual path in `TestP23ReviewReworkManualResolutionWithPriorOutcome` (`:204`). `TestReviewReworkStopRequiresApprovedOutcome` (`internal/workflow/budget_test.go:146`) checks the DEC-235 nullable fields. The CLI check recorded the Outcome Event `EVT-20261002-000027` resolving Stop `STP-d08849` and the next Claim Event `EVT-20261002-000028`.
- VAL-087: PASS for the scoped code paths: `TestReviewCyclesDoNotBlockGateAdmission` (`internal/evidence/execution_test.go:86`) covers the exhausted rework chain; `TestResumeGateReasonsExcludeReworkBudget` (`internal/workflow/budget_test.go:181`) covers active-Run resume reasons and retains the attempts reason; `TestGenericBudgetAdmissionExcludesOnlyReworkAction` (`:168`) checks unknown/exhausted non-rework budgets. The isolated-fixture route is used because the accepted Claim/Run rules do not admit the active Builder Run and exhausting Review together through the CLI. The shared-fake mirror remains the scope gap above.
- VAL-088: PASS for the two later `changes-required` Reviews and their consumed totals 2 and 3, remaining 0, and zero new Stops in `TestP23ReviewCyclesRealProcessJourney` (`cmd/fid/integration/p23_review_evidence_budget_test.go:77`). `TestP23StopFailureHumanAndJSONParity` (`cmd/fid/integration/p23_human_failure_test.go:14`) checks the two exact Human lines; `TestP23StopFailureCodeMatrix` (`:41`) checks unchanged JSON and status across both rework denials.
- VAL-072: DETERMINISTIC FIXTURE READY; Human judgment remains phase-end. `TestP23HumanFixtureStepsRemainRunnable` (`cmd/fid/integration/p23_human_fixtures_test.go:165`) runs denial, qualifying Outcome, and admitted retry. The fresh manifest at `/tmp/fid-task21-l8zUT4/checkpoints-3467353164/manifest.json` records the activated Human Actor and Outcome input. Direct CLI output showed the new Stop condition and successful retry.

## Tests run

- `go test -run '^TestReviewReworkStopRequiresApprovedOutcome$' -count=1 ./internal/workflow` - FAIL (red) - `FAIL` (terminated `exit_code: 1`); old Condition was `budget-observed`. Later PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 0.164s` (terminated `exit_code: 0`).
- `go test -run '^TestReviewCyclesDoNotBlockGateAdmission$' -count=1 ./internal/evidence` - FAIL (red) - `FAIL` (terminated `exit_code: 1`); old admission returned `budget-exhausted`. Later PASS - `ok  github.com/mbell-nzl/fid/internal/evidence 0.005s` (terminated `exit_code: 0`).
- `go test -tags integration -run '^TestP23StopFailureHumanAndJSONParity$' -count=1 ./cmd/fid/integration` - FAIL (red) - `FAIL` (terminated `exit_code: 1`); the old Human denial had only the error line.
- `go test -tags integration -run '^TestP23Review(CyclesRealProcessJourney|ReworkManualResolutionWithPriorOutcome)$' -count=1 ./cmd/fid/integration` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 6.626s` (terminated `exit_code: 0`).
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/workflow ./internal/evidence ./internal/claims ./cmd/fid` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 230.728s` (terminated `exit_code: 0`).
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP23(ReviewCyclesRealProcessJourney|ReviewReworkManualResolutionWithPriorOutcome|HumanFixtureStepsRemainRunnable|HumanFixtures|StopFailureHumanAndJSONParity|StopFailureCodeMatrix)$' ./cmd/fid/integration` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 19.016s` (terminated `exit_code: 0`).
- `FID_P23_HUMAN_FIXTURE_ROOT=/tmp/fid-task21-l8zUT4 go test -tags integration -run '^TestP23HumanFixtures$' -count=1 -v ./cmd/fid/integration` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 11.653s` (terminated `exit_code: 0`).
- Direct `fid claim acquire` in `/tmp/fid-task21-l8zUT4/checkpoints-3467353164/review-exhaustion` - `next: an activated Human records an approved Human Outcome of class review-rework-extension for this Envelope's Issue; see fid stop list --subject envelope:ENV-282fc5` (terminated `exit_code: 4`); direct `fid stop list --subject envelope:ENV-282fc5` named `human-outcome-recorded` (terminated `exit_code: 0`); direct `fid human-outcome record` returned resolved Stop `STP-d08849` (terminated `exit_code: 0`); the same direct `fid claim acquire` returned Claim `CLM-81c95d`, revision 1 (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`). Repeated after the lint repair with the same result.
- `golangci-lint run ./...` - FAIL - `* staticcheck: 1` (terminated `exit_code: 1`). Same exact command after the in-scope repair: PASS - `0 issues.` (terminated `exit_code: 0`).
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 113.517s coverage: 62.7% of statements` (terminated `exit_code: 0`). Package times: `cmd/fid` 274.158s, `internal/workflow` 113.517s; no ordinary package exceeded 600s.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)				60.7%` (terminated `exit_code: 0`).
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 268.910s` (terminated `exit_code: 0`); under 900s.

## Deviations from spec

- Shared-fake mirror remains incomplete because the fake dispatcher and gate denial live in two out-of-scope files named under Workflow outcome. This is the reason for `RECOVERY_REQUIRED`.
- First lint attempt failed with `* staticcheck: 1` (terminated `exit_code: 1`) because the new Condition branch used an if/else chain flagged by QF1003. Changed only that branch to a tagged switch; the byte-identical retry returned `0 issues.` (terminated `exit_code: 0`). The failed attempt remains part of this report.
- The first guidance implementation used an invalid pointer dereference on `cli.Subject.ID`; its focused tagged command ended `FAIL` (terminated `exit_code: 1`). Repaired the string rendering and retained the original red guidance evidence.

## Concerns, technical debt, or notes for the Tech Lead

- The shared fake's `gate.run` dispatch is `cmd/fid/integration/p22_semantic_fake_test.go:71`, and its `denyGate` implementation is `cmd/fid/integration/p22_semantic_denial_test.go:12`. Both fall outside task 21 scope, so no fake change was made. The latter was introduced by Phase 22 commit `25caa7d`; the former is present in `5ec69e6`.
- Retained `internal/claims/generic_test.go` condition fixture came from `0cb705c` and was converged to DEC-269 while preserving the claim-denial assertion. The prior `budget-observed` expectation in the review journey came from `59d2653`, and the single-line Human denial expectation came from task 19's `180f516`; both were converged to this task's accepted change.
- DEC-270 applies if a pre-task-21 journal with the old rework Stop fails to rebuild: that is expected fail-closed behavior. No compatibility path was added, and the new-shape Review journal rebuilt in `TestP23ReviewCyclesRealProcessJourney`.
- No package reached 80% of its declared time budget.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
