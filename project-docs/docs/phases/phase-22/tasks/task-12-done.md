# Phase 22 - Task 12: Recovery freshness and actionable resume - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: The A→B adopted-target successor Gate/Receipt/Evidence case and the paired executable action matrix remain in progress under the same Task 12 scope. The final five-command gate must be rerun from the beginning after those tests and their causal repairs are complete.
- Recovery attempts: 4
- Next action: Complete the retained A→B recovery test and action-mutation matrix, then run every exact task gate in order and finalize this report.
- Human attention required: no

## Summary of what was implemented

- Added current-observation Resume projection, typed next actions, durable Stops/Attention/Recovery context, Claim effective status, and shared JSON/Human rendering.
- Added derived Budget output state, strict output-only input rejection, immutable-policy stripping, and idempotency replay of stored historical result snapshots.
- Changed terminal Issue transition and Envelope completion evaluation to observe the current workspace target rather than trust a stale stored or caller-provided target.
- Converged retained Resume fixtures to fresh real Git repositories, avoiding copied-worktree identity drift.

## Files created or modified

- `cmd/fid/envelope.go`, `cmd/fid/resume.go`, `cmd/fid/gate.go`, and matching command tests.
- `cmd/fid/integration/p22_recovery_journey_test.go` and Phase 22 Resume/Budget/Artifact contract tests.
- `internal/workflow/budget.go`, `envelope.go`, `evaluation.go`, `resume.go`, `resume_actions.go`, `run.go`, and matching tests.
- `internal/cli/codes.go` and `internal/cli/codes_test.go`.
- `project-docs/docs/phases/phase-22/tasks/task-12-done.md`.

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-12-root-evidence.md` was added in Root-owned commits `43ae4c7`, `8670e66`, `8cb8771`, `decc95f`, and `2a7abad`.
- This recovery checkpoint preserves the Developer-observed receipts below.

## Validation Contract Coverage

- VAL-048: PARTIAL - `TestP22EarlyRecoveryJourney` exercises same-HEAD one-byte staleness and Evidence identity; Root's clean/dirty table extension was observed as a focused green only through Root's committed evidence, so it is not re-attested here.
- VAL-049: PARTIAL - current journey covers unchanged clean/dirty recovery and branch rename; the distinct A→B adopted-target successor mutation remains outstanding.
- VAL-050: PARTIAL - Root commit `2a7abad` adds the metadata-only versus required-material Artifact case. The Developer observed the prior Artifact-loss focused command pass: `ok  github.com/mbell-nzl/fid/cmd/fid 0.533s` (terminated `exit_code: 0`).
- VAL-051: PARTIAL - focused retained Resume/Task21 rebuild reads passed: `ok  github.com/mbell-nzl/fid/cmd/fid 2.833s` (terminated `exit_code: 0`); final source changed afterward.
- VAL-052: PARTIAL - terminal target rejection and producer predicate were repaired, but the paired actual-mutation action matrix is outstanding.

## Tests run

- `go test ./internal/workflow -run '^TestDeriveBudgetProjectsEveryAcceptedOutputState$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): initial red before `Budget.State` existed.
- `go test ./internal/workflow -run '^TestDeriveBudgetProjectsEveryAcceptedOutputState$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 0.004s` (terminated `exit_code: 0`).
- `go test ./cmd/fid -run '^TestDecodeClosedEnvelopeInputRejectsAuthoredBudgetState$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): authored empty/null state was initially accepted.
- `go test ./cmd/fid -run '^TestDecodeClosedEnvelopeInputRejectsAuthoredBudgetState$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 0.072s` (terminated `exit_code: 0`).
- `go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): dirty recovery initially omitted the adopted path and reconciliation digest did not bind.
- `go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 1.599s` (terminated `exit_code: 0`)
- `go test ./internal/workflow -run '^TestResumeKeepsOpenAttentionAfterAllLinkedStopsAreResolved$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): all resolved Stops hid a still-open Attention item.
- `go test ./internal/workflow -run '^TestResumeKeepsOpenAttentionAfterAllLinkedStopsAreResolved$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 0.113s` (terminated `exit_code: 0`).
- `go test ./cmd/fid -run '^TestBudgetProjectionAndHistoricalReplayAcrossCLI$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): Envelope and Run replay returned later counters/state.
- `go test ./cmd/fid -run '^TestBudgetProjectionAndHistoricalReplayAcrossCLI$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 0.570s` (terminated `exit_code: 0`).
- `go test ./cmd/fid -run '^TestResumeProjectsExactActiveAndCompletedDurableStateWithoutRecovery/completed$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): actual stale current-target terminal transition and Envelope completion were initially admitted.
- `go test ./cmd/fid -run '^TestResumeProjectsExactActiveAndCompletedDurableStateWithoutRecovery/completed$' - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 1.139s` (terminated `exit_code: 0`)
- `go test ./cmd/fid ./internal/workflow -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 7.112s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - BLOCKED - `exit 0 (no output)` (terminated `exit_code: 0`); source changed after this initial pass and must be rerun.
- `golangci-lint run ./...` - FAIL - `* unused: 1` (terminated `exit_code: 1`); the unused Artifact test helper was removed.
- `golangci-lint run ./...` - BLOCKED - `0 issues.` (terminated `exit_code: 0`); later source changes require a fresh exact-gate run.
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `ok  github.com/mbell-nzl/fid/internal/workflow 92.425s coverage: 52.8% of statements` (terminated `exit_code: 0`); later source changes require a fresh exact-gate run.
- `go tool cover -func=coverage.out` - BLOCKED - not run; source changed after the coverage-producing command and the exact sequence is incomplete.
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - not run; final source freeze has not occurred.

## Deviations from spec

- The specified First Failing Test was not captured before the first Task 12 domain edits. Later focused red/green tests were preserved, including dirty recovery, open Attention, Budget replay, terminal transition, and Envelope completion. This is a test-order deviation, not retrospective TDD evidence.
- The all-Stops-resolved Attention reader repair preceded its decisive focused red case; the subsequent red and green receipts are preserved above.

## Concerns, technical debt, or notes for the Tech Lead

- Root owns the still-uncommitted `cmd/fid/integration/p22_resume_actions_test.go` and the A→B changed-target recovery continuation. Do not finalize this checkpoint or claim final gate success until those are integrated and the exact gate sequence is observed anew.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `2339dc9ebd8d0e6b5502ca32e1562892cf704ce6`

## Commit subject

- `feat(workflow): project current recovery state`
