## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: gate-repair
- Condition: The required real-process integration gate fails in retained pre-Phase-21 fixtures that still send removed `--by` flags to `workflow define`, `workflow assign`, `label add`, and related commands. Those fixtures are outside task 12 scope and predate the Phase 21 start.
- Recovery attempts: 0
- Next action: Phase Runner must route retained integration-fixture convergence to the owning accepted scope, then dispatch the same task 12 continuation to rerun the six exact commands.
- Human attention required: no

## Summary of what was implemented

- Added atomic invalid-input validation to Issue revision Effect construction.
- Added the execution-bound nonterminal transition predicate for Envelope, trusted Actor authority, generic Claim/Run, repository, Issue revision, and workflow edge qualification.
- Added the shared fail-closed Obligation evaluator. Until Group 2 Evidence exists, terminal-success reports `obligation-evidence-absent` and commits no Event.

## Files created or modified

- `cmd/fid/issue_mutation_response.go`
- `cmd/fid/task12_execution_transition_test.go`
- `internal/issues/revision.go`
- `internal/issues/revision_test.go`
- `internal/issues/transition.go`
- `internal/workflow/evaluation.go`
- `internal/workflow/evaluation_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-12-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-12-done.md`

## Validation Contract Coverage

- VAL-038: PASS - `TestExecutionBoundNonterminalTransitionRequiresActiveEnvelopeClaimRunAndRevision` at `cmd/fid/task12_execution_transition_test.go:9` advances one all-qualified nonterminal transition through the CLI root command.
- VAL-039: PASS - `TestExecutionBoundTransitionReturnsEveryOrderedReasonWithoutMutation` at `cmd/fid/task12_execution_transition_test.go:56` proves ordered all-reasons output and an unchanged Event count.
- VAL-040: PASS - `TestEvaluateObligationsReturnsIdenticalUnsatisfiedReasonsForBothGuards` at `internal/workflow/evaluation_test.go:8` compares the terminal-success and `envelope.complete` evaluator outputs byte-for-byte; `cmd/fid/task12_execution_transition_test.go:40` proves terminal transition remains atomic before Evidence exists.
- VAL-041: PASS - `TestBuildIssueRevisionEffectsRejectsInvalidInputsWithoutEffects` at `internal/issues/revision_test.go:42` covers whitespace IDs, normalized duplicates, and overflow with nil Effects.

## Tests run

- `go test ./internal/issues --run TestBuildIssueRevisionEffectsRejectsInvalidInputsWithoutEffects -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/issues\t0.004s` (terminated `exit_code: 0`)
- `go test ./cmd/fid --run TestExecutionBoundNonterminalTransitionRequiresActiveEnvelopeClaimRunAndRevision -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.313s` (terminated `exit_code: 0`)
- `go test ./cmd/fid --run TestExecutionBoundTransitionReturnsEveryOrderedReasonWithoutMutation -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.208s` (terminated `exit_code: 0`)
- `go test ./internal/workflow --run TestEvaluateObligationsReturnsIdenticalUnsatisfiedReasonsForBothGuards -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.003s` (terminated `exit_code: 0`)
- `go test ./internal/issues ./internal/workflow ./cmd/fid --run 'Test(BuildIssueRevisionEffectsRejectsInvalidInputsWithoutEffects|ExecutionBound)' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.452s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race ./internal/workflow ./internal/issues ./internal/policy ./internal/validation ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t(cached)` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`): retained fixtures use removed `--by` flags; provenance is `e712f85d08a0b3259b23badd34d9306bf0a10f9b`, `7119d5b38d79663f544f093d55f676c375746419`, and `09aa0a081d2ecfee025c3b4d8158bc6203d02d98`, all before Phase 21.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t24.081s\tcoverage: 33.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `65.3%` (terminated `exit_code: 0`)

## Deviations from spec

- Required integration gate failed before task completion. The failure is preserved; no command, assertion, scope, or acceptance contract was changed.

## Concerns, technical debt, or notes for the Tech Lead

- Retained integration fixtures in `cmd/fid/integration/` are inconsistent with the accepted Phase 21 `--actor` command contract. They cannot be changed within task 12 scope.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `93da8a150d76b10d63bb8330517c508728bbe45f`

## Commit subject

- `feat(workflow): add execution-bound transition predicate`
