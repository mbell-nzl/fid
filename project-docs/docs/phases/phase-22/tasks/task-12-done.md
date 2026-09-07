# Phase 22 - Task 12: Recovery freshness and actionable resume - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Recovery attempts: 4
- Human attention required: no

## Summary of what was implemented

- Projected recovery freshness from the currently observed repository while retaining immutable Receipt and Evidence identities.
- Allowed a recovered successor Run to execute within its original project and worktree while rejecting a different project or worktree.
- Exercised real CLI recovery from target A to adopted target B, then created successor Gate, Receipt and Evidence facts, completed the terminal Issue transition, closed the successor Run, and completed the Envelope.
- Converged terminal retained fixtures to observed Git repositories and retained Budget replay snapshots for all public Run mutations.

## Files created or modified

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
- VAL-050: PASS - `TestArtifactMaterialPolicy` in `cmd/fid/artifact_material_policy_test.go` proves metadata-only loss remains admissible while required Artifact loss changes only availability and policy outcome without replacing Evidence or Receipt facts.
- VAL-051: PASS - `TestP22ResumeHumanMatchesEligibleAndStoppedRecoveryState`, `TestP22ResumeHumanReportsExpiredClaimEffectiveStatus`, and `TestP22RunReplayRetainsHistoricalBudgetSnapshot` exercise Resume reconstruction, Human/JSON projection parity, effective lease status, derived Budget state, and historical snapshots.
- VAL-052: PASS - `TestP22ResumeReadyCandidatesExecuteTheirActualOperations`, `TestP22ResumeObserveRejectsUntrustedProducerDespiteCapabilities`, `TestP22RecoveryAdoptedTargetPermitsTerminalTransitionWithFreshSuccessorProof`, and `TestExecutionBoundTransitionPredicateMatrix` execute qualifying mutations and retain exact rejection predicates for untrusted producers and foreign repository identity.

## Tests run

- `go test -tags integration ./cmd/fid/integration -run '^TestP22RecoveryAdoptsChangedTargetForSuccessorEvidence$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): the first recovery fixture used a stale Run version after Gate A.
- `go test -tags integration ./cmd/fid/integration -run '^TestP22RecoveryAdoptsChangedTargetForSuccessorEvidence$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 1.381s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP22RecoveryAdoptedTargetPermitsTerminalTransitionWithFreshSuccessorProof$' -count=1` - FAIL - `FAIL  github.com/mbell-nzl/fid/cmd/fid/integration 1.708s` (terminated `exit_code: 1`): the actual target-B Issue transition reported `run-inactive` because the active Run was compared with the original Envelope target.
- `go test ./cmd/fid -run '^TestExecutionBoundTransitionPredicateMatrix/run_repository_project_mismatched$' -count=1` - FAIL - `FAIL  github.com/mbell-nzl/fid/cmd/fid 0.407s` (terminated `exit_code: 1`): worktree-only matching admitted a different Project ID.
- `go test ./cmd/fid -run '^TestExecutionBoundTransitionPredicateMatrix/(run_repository_mismatched|run_repository_project_mismatched)$' -count=1 && go test -tags integration ./cmd/fid/integration -run '^TestP22RecoveryAdoptedTargetPermitsTerminalTransitionWithFreshSuccessorProof$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 1.971s` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - FAIL - `FAIL  github.com/mbell-nzl/fid/cmd/fid/integration 49.665s` (terminated `exit_code: 1`): retained Phase 21 terminal journey had a synthetic repository target and the new observer correctly rejected it.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 51.320s` (terminated `exit_code: 0`)
- Root-attributed DEC-201 replay regression `TestP22RunReplayRetainsHistoricalBudgetSnapshot` - FAIL then PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 3.943s` (terminated `exit_code: 0`); raw receipts and four-case history are in `task-12-root-evidence.md`.
- The changed behavior was exercised through the built `fid` CLI in the real integration journeys named above, including `run recover`, `gate run`, `evidence link`, `issue transition`, `run close`, and `envelope complete` at adopted target B.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - FAIL - `* unused: 1` (terminated `exit_code: 1`): Root's observed-fixture convergence left an unused wrapper; the obsolete helper was removed without changing assertions.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 85.521s coverage: 52.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:                                  (statements)                  62.4%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 51.743s` (terminated `exit_code: 0`)

## Deviations from spec

- The specified First Failing Test was not captured before the first Task 12 domain edits. The later actual A→B CLI terminal-transition red and project-identity red are retained above; this is a test-order deviation, not retrospective TDD evidence.
- Earlier full-gate attempts were invalidated by later accepted in-scope replay and retained-fixture repairs. Their observed receipts remain in the recovery history; the five results above are the final source-bound gate.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `c2d8098bfe271ce87e14b60a2778eb4532a60506`

## Commit subject

- `feat(workflow): reconcile recovered execution targets`
