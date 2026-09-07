# Phase 22 - Task 12: Recovery freshness and actionable resume - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added current-observation Resume projection, typed next actions, durable Stops/Attention/Recovery context, Claim effective status, and shared JSON/Human rendering.
- Added derived Budget output state, strict output-only input rejection, immutable-policy stripping, and idempotency replay of stored historical result snapshots.
- Changed terminal Issue transition and Envelope completion evaluation to observe the current workspace target rather than trust a stale stored or caller-provided target.
- Converged retained Resume fixtures to fresh real Git repositories, avoiding copied-worktree identity drift.

- Projected recovery freshness from the currently observed repository while retaining immutable Receipt and Evidence identities.
- Allowed a recovered successor Run to execute within its original project and worktree while rejecting a different project or worktree.
- Exercised real CLI recovery from target A to adopted target B, then created successor Gate, Receipt and Evidence facts, completed the terminal Issue transition, closed the successor Run, and completed the Envelope.
- Converged terminal retained fixtures to observed Git repositories and retained Budget replay snapshots for all public Run mutations.

## Files created or modified

- `cmd/fid/envelope.go`, `cmd/fid/resume.go`, `cmd/fid/gate.go`, and matching command tests.
- `cmd/fid/integration/p22_recovery_journey_test.go` and Phase 22 Resume/Budget/Artifact contract tests.
- `internal/workflow/budget.go`, `envelope.go`, `evaluation.go`, `resume.go`, `resume_actions.go`, `run.go`, and matching tests.
- `internal/cli/codes.go` and `internal/cli/codes_test.go`.
- `project-docs/docs/phases/phase-22/tasks/task-12-done.md`.

- `internal/workflow/evaluation.go` and `cmd/fid/task12_execution_transition_test.go`.
- `cmd/fid/integration/p22_recovery_adopted_target_test.go`, `p22_resume_actions_test.go`, and `p21_conformance_test.go`.
- `cmd/fid/evidence_test.go` and Root-owned DEC-201 replay and retained-fixture convergence files committed during Task 12.
- `project-docs/docs/phases/phase-22/tasks/task-12-done.md`.

## Documentation updated

- `project-docs/docs/phases/phase-22/decisions.md` records DEC-201 output/replay closure.
- `project-docs/docs/phases/phase-22/tasks/task-12-root-evidence.md` retains Root-attributed red/green evidence and convergence provenance.

## Validation Contract Coverage

- VAL-048: PASS - `TestP22EarlyRecoveryJourney` and `TestP22RecoveryAdoptsChangedTargetForSuccessorEvidence` in `cmd/fid/integration` retain old Evidence identity while a one-byte dirty target change stales its proof and produces new successor facts.
- VAL-049: PASS - `TestP22EarlyRecoveryJourney` covers unchanged clean and dirty targets across recovery and branch rename; `TestP22RecoveryAdoptsChangedTargetForSuccessorEvidence` preserves the original Registration and Obligation binding at the adopted in-scope target.
- VAL-050: PASS - `TestArtifactLossHonorsRequiredMaterialAndPreservesProofIdentity` in `cmd/fid/artifact_material_policy_test.go` proves metadata-only loss remains admissible while required Artifact loss changes only availability and policy outcome without replacing Evidence or Receipt facts.
- VAL-051: PASS - `TestP22EarlyRecoveryJourney` also compares Resume before and after projection rebuild; `TestP22ResumeHumanMatchesEligibleAndStoppedRecoveryState`, `TestP22ResumeHumanDistinguishesElapsedLeaseFromRecordedStatus`, and `TestP22RunReplayRetainsHistoricalBudgetSnapshot` exercise Resume reconstruction, Human/JSON projection parity, effective lease status, derived Budget state, and historical snapshots.
- VAL-052: PASS - `TestP22ResumeReadyCandidatesExecuteTheirActualOperations`, `TestP22ResumeObserveRejectsUntrustedProducerDespiteCapabilities`, `TestP22RecoveryAdoptedTargetPermitsTerminalTransitionWithFreshSuccessorProof`, and `TestExecutionBoundTransitionPredicateMatrix` execute qualifying mutations and retain exact rejection predicates for untrusted producers and foreign repository identity.

## Tests run

- `go test -tags integration ./cmd/fid/integration -run '^TestP22RecoveryAdoptsChangedTargetForSuccessorEvidence$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): the first recovery fixture used a stale Run version after Gate A.
- `go test -tags integration ./cmd/fid/integration -run '^TestP22RecoveryAdoptsChangedTargetForSuccessorEvidence$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 1.381s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP22RecoveryAdoptedTargetPermitsTerminalTransitionWithFreshSuccessorProof$' -count=1` - FAIL - `FAIL  github.com/mbell-nzl/fid/cmd/fid/integration 1.708s` (terminated `exit_code: 1`): the actual target-B Issue transition reported `run-inactive` because the active Run was compared with the original Envelope target.
- `go test ./cmd/fid -run '^TestExecutionBoundTransitionPredicateMatrix/run_repository_project_mismatched$' -count=1` - FAIL - `FAIL  github.com/mbell-nzl/fid/cmd/fid 0.407s` (terminated `exit_code: 1`): worktree-only matching admitted a different Project ID.
- `go test ./cmd/fid -run '^TestExecutionBoundTransitionPredicateMatrix/(run_repository_mismatched|run_repository_project_mismatched)$' -count=1 && go test -tags integration ./cmd/fid/integration -run '^TestP22RecoveryAdoptedTargetPermitsTerminalTransitionWithFreshSuccessorProof$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 1.971s` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - FAIL - `FAIL  github.com/mbell-nzl/fid/cmd/fid/integration 49.665s` (terminated `exit_code: 1`): retained Phase 21 terminal journey had a synthetic repository target and the new observer correctly rejected it.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 51.320s` (terminated `exit_code: 0`)
- Root-attributed DEC-201 replay regression `TestP22RunReplayRetainsHistoricalBudgetSnapshot`: exact raw receipts and four-case red/green history are retained in `task-12-root-evidence.md`; the Developer final integration gate also executes this table.
- The changed behavior was exercised through the built `fid` CLI in the real integration journeys named above, including `run recover`, `gate run`, `evidence link`, `issue transition`, `run close`, and `envelope complete` at adopted target B.

## Exact quality-gate results

The Developer recorded these final receipts in `task-12-gate-addendum.md`
(commit 833b4fc); they retain that session attribution. That addendum is part
of this report and preserves all continuation attempts separately.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 52.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					62.4%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration	51.743s` (terminated `exit_code: 0`)

## Deviations from spec

- Mandatory race-gate deviation: the first continuation full race attempt printed `FAIL` (terminated `exit_code: 1`) in session 90704. Current-target observation exposed retained synthetic Stop/terminal Evidence fixtures. Root and Developer converted those fixtures to fresh observed Git targets, preserving their predicates and recording provenance in the root evidence. Subsequent accepted replay changes invalidated an intermediate green; the unchanged command's final closing receipt is the cached 52.9% output in the Developer addendum above.
- Mandatory integration-gate deviation: the later exact integration gate printed `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration	49.665s` (terminated `exit_code: 1`). The retained Phase 21 task-22 journey used a synthetic repository target. Its fixture now initializes Git and ignores generated inputs/artifacts before binding. The byte-identical causal rerun printed the 51.320s green recorded above, and the final ordered gate printed the 51.743s green in the Developer addendum.
- Mandatory lint-gate deviation: `* unused: 1` (terminated `exit_code: 1`) identified an obsolete test wrapper after fixture convergence. Removing that unused wrapper retained all assertions; the byte-identical lint command then printed `0 issues.` (terminated `exit_code: 0`).
- The specified First Failing Test was not captured before the first Task 12 domain edits. The later actual A→B CLI terminal-transition red and project-identity red are retained above; this is a test-order deviation, not retrospective TDD evidence.
- Earlier full-gate attempts were invalidated by later accepted in-scope replay and retained-fixture repairs. Their observed receipts remain in the recovery history; the five results above are the final source-bound gate.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `c2d80980b9bfa9c1c20723cc7751526eb8baf2bf`

## Commit subject

- `feat(workflow): reconcile recovered execution targets`

## Root report preservation and provenance

The retained real-process Phase 21 fixture was introduced by
`a99b4270f70e98111b5bae43c54ed92aed7078c6` (phase-21/task-22) before the Phase 22 planning base.
Task 12 changed its repository setup to actual Git observation while retaining
its activated-Envelope journey and terminal assertions. The exact full tagged
gate executed both its failed and repaired forms, as recorded above.

Root restored the earlier checkpoint verbatim below after the completion-report
rewrite omitted its failed and superseded receipts. The quoted checkpoint is
historical, including its then-pending conditions, local repair counter and
BLOCKED rows; the current Runner ledger has one same-task continuation. The
original Developer session remains the source of its receipts. Original
root-labeled control receipts retain root attribution. Final-source gate
corrections and continuation receipts are separately retained in the Developer's
`task-12-gate-addendum.md`; root does not attest those as root executions.

Root's observed report-control correction at 2278953 and successful recovery
report lint are retained in `task-12-root-evidence.md`. Its later root-observed
coverage cross-check confirmed 62.4% with an actual terminated exit_code: 0.

### Retained checkpoint at 2278953

> # Phase 22 - Task 12: Recovery freshness and actionable resume - Completion Report
>
> ## Workflow outcome
>
> - Task outcome: RECOVERY_REQUIRED
> - Category: in-scope-implementation
> - Condition: The A→B adopted-target successor Gate/Receipt/Evidence case and the paired executable action matrix remain in progress under the same Task 12 scope. The final five-command gate must be rerun from the beginning after those tests and their causal repairs are complete.
> - Recovery attempts: 4
> - Next action: Complete the retained A→B recovery test and action-mutation matrix, then run every exact task gate in order and finalize this report.
> - Human attention required: no
>
> ## Summary of what was implemented
>
> - Added current-observation Resume projection, typed next actions, durable Stops/Attention/Recovery context, Claim effective status, and shared JSON/Human rendering.
> - Added derived Budget output state, strict output-only input rejection, immutable-policy stripping, and idempotency replay of stored historical result snapshots.
> - Changed terminal Issue transition and Envelope completion evaluation to observe the current workspace target rather than trust a stale stored or caller-provided target.
> - Converged retained Resume fixtures to fresh real Git repositories, avoiding copied-worktree identity drift.
>
> ## Files created or modified
>
> - `cmd/fid/envelope.go`, `cmd/fid/resume.go`, `cmd/fid/gate.go`, and matching command tests.
> - `cmd/fid/integration/p22_recovery_journey_test.go` and Phase 22 Resume/Budget/Artifact contract tests.
> - `internal/workflow/budget.go`, `envelope.go`, `evaluation.go`, `resume.go`, `resume_actions.go`, `run.go`, and matching tests.
> - `internal/cli/codes.go` and `internal/cli/codes_test.go`.
> - `project-docs/docs/phases/phase-22/tasks/task-12-done.md`.
>
> ## Documentation updated
>
> - `project-docs/docs/phases/phase-22/tasks/task-12-root-evidence.md` was added in Root-owned commits `43ae4c7`, `8670e66`, `8cb8771`, `decc95f`, and `2a7abad`.
> - This recovery checkpoint preserves the Developer-observed receipts below.
>
> ## Validation Contract Coverage
>
> - VAL-048: PARTIAL - `TestP22EarlyRecoveryJourney` exercises same-HEAD one-byte staleness and Evidence identity; Root's clean/dirty table extension was observed as a focused green only through Root's committed evidence, so it is not re-attested here.
> - VAL-049: PARTIAL - current journey covers unchanged clean/dirty recovery and branch rename; the distinct A→B adopted-target successor mutation remains outstanding.
> - VAL-050: PARTIAL - Root commit `2a7abad` adds the metadata-only versus required-material Artifact case. The Developer observed the prior Artifact-loss focused command pass: `ok  github.com/mbell-nzl/fid/cmd/fid 0.533s` (terminated `exit_code: 0`).
> - VAL-051: PARTIAL - focused retained Resume/Task21 rebuild reads passed: `ok  github.com/mbell-nzl/fid/cmd/fid 2.833s` (terminated `exit_code: 0`); final source changed afterward.
> - VAL-052: PARTIAL - terminal target rejection and producer predicate were repaired, but the paired actual-mutation action matrix is outstanding.
>
> ## Tests run
>
> - `go test ./internal/workflow -run '^TestDeriveBudgetProjectsEveryAcceptedOutputState$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): initial red before `Budget.State` existed.
> - `go test ./internal/workflow -run '^TestDeriveBudgetProjectsEveryAcceptedOutputState$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 0.004s` (terminated `exit_code: 0`)
> - `go test ./cmd/fid -run '^TestDecodeClosedEnvelopeInputRejectsAuthoredBudgetState$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): authored empty/null state was initially accepted.
> - `go test ./cmd/fid -run '^TestDecodeClosedEnvelopeInputRejectsAuthoredBudgetState$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 0.072s` (terminated `exit_code: 0`)
> - `go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): dirty recovery initially omitted the adopted path and reconciliation digest did not bind.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 1.599s` (terminated `exit_code: 0`)
> - `go test ./internal/workflow -run '^TestResumeKeepsOpenAttentionAfterAllLinkedStopsAreResolved$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): all resolved Stops hid a still-open Attention item.
> - `go test ./internal/workflow -run '^TestResumeKeepsOpenAttentionAfterAllLinkedStopsAreResolved$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 0.113s` (terminated `exit_code: 0`)
> - `go test ./cmd/fid -run '^TestBudgetProjectionAndHistoricalReplayAcrossCLI$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): Envelope and Run replay returned later counters/state.
> - `go test ./cmd/fid -run '^TestBudgetProjectionAndHistoricalReplayAcrossCLI$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 0.570s` (terminated `exit_code: 0`)
> - `go test ./cmd/fid -run '^TestResumeProjectsExactActiveAndCompletedDurableStateWithoutRecovery/completed$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): actual stale current-target terminal transition and Envelope completion were initially admitted.
> - `go test ./cmd/fid -run '^TestResumeProjectsExactActiveAndCompletedDurableStateWithoutRecovery/completed$' - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 1.139s` (terminated `exit_code: 0`)
> - `go test ./cmd/fid ./internal/workflow -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 7.112s` (terminated `exit_code: 0`)
>
> ## Exact quality-gate results
>
> - `gofmt -l cmd internal` - BLOCKED - `exit 0 (no output)` (terminated `exit_code: 0`); source changed after this initial pass and must be rerun.
> - `golangci-lint run ./...` - FAIL - `* unused: 1` (terminated `exit_code: 1`); the unused Artifact test helper was removed.
> - `golangci-lint run ./...` - BLOCKED - `0 issues.` (terminated `exit_code: 0`); later source changes require a fresh exact-gate run.
> - `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `ok  github.com/mbell-nzl/fid/internal/workflow 92.425s coverage: 52.8% of statements` (terminated `exit_code: 0`); later source changes require a fresh exact-gate run.
> - `go tool cover -func=coverage.out` - BLOCKED - not run; source changed after the coverage-producing command and the exact sequence is incomplete.
> - `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - not run; final source freeze has not occurred.
>
> ## Deviations from spec
>
> - The specified First Failing Test was not captured before the first Task 12 domain edits. Later focused red/green tests were preserved, including dirty recovery, open Attention, Budget replay, terminal transition, and Envelope completion. This is a test-order deviation, not retrospective TDD evidence.
> - The all-Stops-resolved Attention reader repair preceded its decisive focused red case; the subsequent red and green receipts are preserved above.
>
> ## Concerns, technical debt, or notes for the Tech Lead
>
> - Root control repair: `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-22/tasks/task-12-done.md project-docs/docs/phases/phase-22/tasks/task-12.md` printed `FAIL report-lint found 4 failure(s)` (terminated `exit_code: 4`). Four focused PASS rows ended with punctuation after termination; root removed only that punctuation. The byte-identical rerun remains required. Earlier Developer report-lint findings and report-commit sandbox failures are retained in the root contribution evidence; no gate definition or result changed.
>
> - Root owns the still-uncommitted `cmd/fid/integration/p22_resume_actions_test.go` and the A→B changed-target recovery continuation. Do not finalize this checkpoint or claim final gate success until those are integrated and the exact gate sequence is observed anew.
>
> ## Dependency-budget note
>
> - No dependency changes.
>
> ## Commit hash
>
> - `2339dc9ebd8d0e6b5502ca32e1562892cf704ce6`
>
> ## Commit subject
>
> - `feat(workflow): project current recovery state`
