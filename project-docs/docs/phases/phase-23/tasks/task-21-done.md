# Phase 23 - Task 21: Resolvable review-rework exhaustion - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: All task-21 work and required gates complete.
- Recovery attempts: 2
- Next action: None.
- Same-task continuation attempt 2 completed the bounded shared-fake conformance repair. Attempt 1 ended `RECOVERY_REQUIRED` because the required fake mirror lived in `p22_semantic_fake_test.go` and `p22_semantic_denial_test.go`, outside the older task file list. The Runner's attempt-2 recovery authorized exactly those two paths and this report.
- Human attention required: no

## Summary of what was implemented

- Review-cycle exhaustion now raises a `human-outcome-recorded` Stop conditioned on an approved `review-rework-extension` Outcome for the Envelope. The existing DEC-250 predicate handles atomic and manual resolution.
- One shared helper excludes rework budgets from generic gate admission and resume reasons. The Phase 22 admission function remains available for other budgets.
- Human denial adds the exact guidance line after task 19's error line; JSON behavior is unchanged.
- Built-binary journeys cover wrong class, rejected result, another Issue, non-human authority denial, qualifying atomic and manual resolution, admitted retry, and two later single-charge Reviews.
- The shared semantic fake now charges a review cycle into the DEC-269 `human-outcome-recorded` Stop shape and skips rework budgets during generic gate admission. Its existing token-budget denial remains in the dual-adapter corpus.

## Files created or modified

- `internal/workflow/budget.go`; `internal/workflow/resume_actions.go`; `internal/workflow/budget_test.go`
- `internal/evidence/execution.go`; `internal/evidence/execution_test.go`
- `internal/claims/generic_test.go`
- `cmd/fid/stop.go`; `cmd/fid/integration/p23_review_evidence_budget_test.go`; `cmd/fid/integration/p23_human_fixtures_test.go`; `cmd/fid/integration/p23_human_failure_test.go`
- `cmd/fid/integration/p22_semantic_fake_test.go`; `cmd/fid/integration/p22_semantic_denial_test.go` (bounded continuation)
- `project-docs/docs/phases/phase-23/tasks/task-21-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-23/qa/human-fixtures.md` rework section adds the resolution and admitted retry instructions. Its prior verbatim snapshot is identified as historical; Phase-end QA regenerates it.

## Validation Contract Coverage

- VAL-086: PASS for the real-process behavior in `TestP23ReviewCyclesRealProcessJourney` (`cmd/fid/integration/p23_review_evidence_budget_test.go:77`) and the separate manual path in `TestP23ReviewReworkManualResolutionWithPriorOutcome` (`:204`). `TestReviewReworkStopRequiresApprovedOutcome` (`internal/workflow/budget_test.go:146`) checks the DEC-235 nullable fields. The CLI check recorded the Outcome Event `EVT-20261002-000027` resolving Stop `STP-d08849` and the next Claim Event `EVT-20261002-000028`.
- VAL-087: PASS for the scoped code paths: `TestReviewCyclesDoNotBlockGateAdmission` (`internal/evidence/execution_test.go:86`) covers the exhausted rework chain; `TestResumeGateReasonsExcludeReworkBudget` (`internal/workflow/budget_test.go:181`) covers active-Run resume reasons and retains the attempts reason; `TestGenericBudgetAdmissionExcludesOnlyReworkAction` (`:168`) checks unknown/exhausted non-rework budgets. The isolated-fixture route is used because the accepted Claim/Run rules do not admit the active Builder Run and exhausting Review together through the CLI. `TestP22FakeGenericAdmissionExcludesReviewCycles` checks unknown and exhausted rework exclusion while both other-budget denial classes remain selected.
- VAL-088: PASS for the two later `changes-required` Reviews and their consumed totals 2 and 3, remaining 0, and zero new Stops in `TestP23ReviewCyclesRealProcessJourney` (`cmd/fid/integration/p23_review_evidence_budget_test.go:77`). `TestP23StopFailureHumanAndJSONParity` (`cmd/fid/integration/p23_human_failure_test.go:14`) checks the two exact Human lines; `TestP23StopFailureCodeMatrix` (`:41`) checks unchanged JSON and status across both rework denials.
- VAL-072: DETERMINISTIC FIXTURE READY; Human judgment remains phase-end. `TestP23HumanFixtureStepsRemainRunnable` (`cmd/fid/integration/p23_human_fixtures_test.go:165`) runs denial, qualifying Outcome, and admitted retry. The fresh manifest at `/tmp/fid-task21-l8zUT4/checkpoints-3467353164/manifest.json` records the activated Human Actor and Outcome input. Direct CLI output showed the new Stop condition and successful retry.
- VAL-068: PASS for this continuation's shared-fake obligation. `TestP22FakeReviewReworkStopCondition` pins the Envelope subject, two blocked actions, resolver, DEC-235 nullable fields, Outcome class, budget details, and no re-trip on a later charge. `TestP22FakeGenericAdmissionExcludesReviewCycles` pins both excluded rework states, both retained generic denial classes, and the shared gate denial for tokens. `TestP22DualAdapterRecoveryConformance` ran in the full tagged gate; expectations are authored in the fake tests, not read from fid output.

## Tests run

- `go test -run '^TestReviewReworkStopRequiresApprovedOutcome$' -count=1 ./internal/workflow` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -run '^TestReviewReworkStopRequiresApprovedOutcome$' -count=1 ./internal/workflow` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 0.164s` (terminated `exit_code: 0`)
- `go test -run '^TestReviewCyclesDoNotBlockGateAdmission$' -count=1 ./internal/evidence` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -run '^TestReviewCyclesDoNotBlockGateAdmission$' -count=1 ./internal/evidence` - PASS - `ok  github.com/mbell-nzl/fid/internal/evidence 0.005s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP23StopFailureHumanAndJSONParity$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -run '^TestP23Review(CyclesRealProcessJourney|ReworkManualResolutionWithPriorOutcome)$' -count=1 ./cmd/fid/integration` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 6.626s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/workflow ./internal/evidence ./internal/claims ./cmd/fid` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 230.728s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP23(ReviewCyclesRealProcessJourney|ReviewReworkManualResolutionWithPriorOutcome|HumanFixtureStepsRemainRunnable|HumanFixtures|StopFailureHumanAndJSONParity|StopFailureCodeMatrix)$' ./cmd/fid/integration` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 19.016s` (terminated `exit_code: 0`)
- `FID_P23_HUMAN_FIXTURE_ROOT=/tmp/fid-task21-l8zUT4 go test -tags integration -run '^TestP23HumanFixtures$' -count=1 -v ./cmd/fid/integration` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 11.653s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP22Fake(ReviewReworkStopCondition|GenericAdmissionExcludesReviewCycles)$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`; missing fake budget type and admission method)
- `go test -tags integration -run '^TestP22Fake(ReviewReworkStopCondition|GenericAdmissionExcludesReviewCycles)$' -count=1 ./cmd/fid/integration` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 1.086s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP22Fake(ReviewReworkStopCondition|GenericAdmissionExcludesReviewCycles)$' -count=1 ./cmd/fid/integration` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 0.570s` (terminated `exit_code: 0`) — final focused assertions.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP22Fake(ReviewReworkStopCondition|GenericAdmissionExcludesReviewCycles)$' ./cmd/fid/integration` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 1.605s` (terminated `exit_code: 0`)

The first two red tests observed `budget-observed` and `budget-exhausted`,
respectively. The guidance red test observed only task 19's error line.
- Direct `fid claim acquire` in `/tmp/fid-task21-l8zUT4/checkpoints-3467353164/review-exhaustion` - `next: an activated Human records an approved Human Outcome of class review-rework-extension for this Envelope's Issue; see fid stop list --subject envelope:ENV-282fc5` (terminated `exit_code: 4`); direct `fid stop list --subject envelope:ENV-282fc5` named `human-outcome-recorded` (terminated `exit_code: 0`); direct `fid human-outcome record` returned resolved Stop `STP-d08849` (terminated `exit_code: 0`); the same direct `fid claim acquire` returned Claim `CLM-81c95d`, revision 1 (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - FAIL - `* staticcheck: 1` (terminated `exit_code: 1`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 113.517s coverage: 62.7% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)				60.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 268.910s` (terminated `exit_code: 0`)

Package times were `cmd/fid` 274.158s, `internal/workflow` 113.517s, and
`cmd/fid/integration` 268.910s. All stayed within their declared budgets.

Attempt 2, same exact commands in the task-spec order from repository root:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`) — wall 0.036s.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`) — wall about 1.004s.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.7% of statements` (terminated `exit_code: 0`) — reported `cmd/fid` wall 215.093s.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)				60.7%` (terminated `exit_code: 0`) — wall 0.114s.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	291.867s` (terminated `exit_code: 0`) — reported package wall 291.867s.

Attempt 2 reported `cmd/fid` 215.093s, `internal/workflow` as cached (attempt 1 measured 113.517s), and `cmd/fid/integration` 291.867s. These are below 600s/900s; no headroom warning applies.

## Deviations from spec

- Attempt 1's `RECOVERY_REQUIRED` scope gap was that the shared-fake dispatcher and denial lived in two files outside the older task list. The Runner's bounded attempt-2 continuation authorized those exact paths; both are now converged and the five byte-identical gates passed.
- The first committed report-lint attempt found 12 report-format failures (`FAIL report-lint found 12 failure(s)`, terminated `exit_code: 12`): missing COMPLETE outcome fields and termination annotations with extra semicolon text. This report-only correction restores the template fields and puts wall-time notes after each termination annotation.
- First lint attempt failed with `* staticcheck: 1` (terminated `exit_code: 1`) because the new Condition branch used an if/else chain flagged by QF1003. Changed only that branch to a tagged switch; the byte-identical retry returned `0 issues.` (terminated `exit_code: 0`). The failed attempt remains part of this report.
- The first guidance implementation used an invalid pointer dereference on `cli.Subject.ID`; its focused tagged command ended `FAIL` (terminated `exit_code: 1`). Repaired the string rendering and retained the original red guidance evidence.

## Concerns, technical debt, or notes for the Tech Lead

- Attempt 1 established that the shared fake's `gate.run` dispatch was `cmd/fid/integration/p22_semantic_fake_test.go:71`, and its `denyGate` implementation was `cmd/fid/integration/p22_semantic_denial_test.go:12`. Both were outside the original task 21 list and unchanged in attempt 1. The latter was introduced by Phase 22 commit `25caa7d`; the former is present in `5ec69e6`. Attempt 2 changed only these fake paths under the bounded recovery.
- Attempt 1's product implementation commit was `c6e070b1ea0b3aff339564c47bb1cd39cfccbcd0` (`feat(review): make rework exhaustion resolvable`); this report's final hash identifies the attempt-2 fake conformance commit.
- Retained `internal/claims/generic_test.go` condition fixture came from `0cb705c` and was converged to DEC-269 while preserving the claim-denial assertion. The prior `budget-observed` expectation in the review journey came from `59d2653`, and the single-line Human denial expectation came from task 19's `180f516`; both were converged to this task's accepted change.
- DEC-270 applies if a pre-task-21 journal with the old rework Stop fails to rebuild: that is expected fail-closed behavior. No compatibility path was added, and the new-shape Review journal rebuilt in `TestP23ReviewCyclesRealProcessJourney`.
- No package reached 80% of its declared time budget.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `343b458fc4c1b7a726ae9218315edf319c5e40ba`

## Commit subject

- `test(review): align semantic fake with rework exhaustion`
