# Phase 22 - Task 10: Activation recovery domain - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Persisted immutable predecessor and five-resource activation bindings, versioned resource observations, exact staged authority/Gate digest selection, and exact resource agreement checks before commit, finalization, and recovery completion.
- Added same-activation recovery selection and completion with forward-only post-commit direction, durable idempotency snapshots, atomic activation-incomplete Stop/Attention facts, and terminal Stop disposition.
- Blocked trusted actor and Registration resolution while any activation is nonterminal, while retaining inspection, recovery, and export paths.
- Removed the final DEC-192 activation recovery projection fence and converged the task-1 live/replay schema assertion to an unknown Event with an invalid effect, retaining corrupt Stop-code parity.

## Files created or modified

- `internal/workflow/activation.go`, `internal/workflow/activation_recovery.go`, and `internal/workflow/activation_recovery_test.go`
- `cmd/fid/activation.go`, `cmd/fid/activation_test.go`, `cmd/fid/activation_export_test.go`, `cmd/fid/activation_resource_matrix_test.go`, and `cmd/fid/integration/p21_conformance_test.go`
- `internal/evidence/registration.go`, `internal/cli/codes.go`, and their matching tests
- `internal/journal/project_effects.go`, `internal/journal/checkpoint_projection_test.go`, and `internal/journal/recovery_schema_test.go`
- Retained surface convergence tests in `cmd/fid/contract_surface_test.go`, `cmd/fid/envelope_test.go`, and `cmd/fid/gate_test.go`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-10-done.md`

## Validation Contract Coverage

- VAL-040: PASS - `TestActivationPrepareCreatesPreparedActivation` and `TestActivationRequiresEveryResourceAtCommitAndFinalize` in `cmd/fid/activation_test.go` and `cmd/fid/activation_resource_matrix_test.go` prove first-activation bindings, separately staged components, exact digest selection, and all five resource agreements; `TestIncompleteUpgradeBlocksNormalWorkWhileInspectionAndRestoreRemainAvailable` proves an upgrade retains exact predecessor resources.
- VAL-041: PASS - `TestCommittedActivationCanOnlyRecoverForward` proves the public committed-abort direction rejection, while `TestActivationRecoverySelectionRaisesAndCompletionDisposesIncompleteStop` proves pre-commit restore and rejects a fresh terminal selection.
- VAL-042: PASS - `TestActivationRequiresEveryResourceAtCommitAndFinalize` covers each missing, null, and mismatched backend, kit, overlay, project, and provider state at both commit and finalization; `TestCompleteResourceAgreementFinalizesTrust` supplies the all-agree control.
- VAL-043: PASS - `TestIncompleteUpgradeBlocksNormalWorkWhileInspectionAndRestoreRemainAvailable` rejects trusted actor resolution and a real `gate register` mutation while allowing status and restore; `TestIncompleteActivationAllowsExportWithoutChangingActivation` proves real export create/validate leave the nonterminal activation and Event count unchanged.
- VAL-044: PASS - `TestActivationRecoveryPreservesDirectionAndIdentityAcrossReplay` and `TestActivationRecoverySelectionRaisesAndCompletionDisposesIncompleteStop` prove one selected Event identity, immutable direction, completion, and response-loss replay after completion.

## Tests run

- `go test ./cmd/fid -run '^TestCommittedActivationCanOnlyRecoverForward$' -count=1` - RED then PASS - `FAIL` (terminated `exit_code: 1`); later `ok github.com/mbell-nzl/fid/cmd/fid 0.188s` (terminated `exit_code: 0`)
- `go test ./internal/workflow -run '^TestActivationRecoverySelectionRaisesAndCompletionDisposesIncompleteStop$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow 0.160s` (terminated `exit_code: 0`)
- `go test ./internal/workflow -run '^TestCompleteActivationRecoveryRejectsMissingPriorResourceSet$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow 0.087s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestActivation(RequiresEveryResourceAtCommitAndFinalize|CompleteResourceAgreementFinalizesTrust)$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 1.397s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/...` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid/integration 11.407s` (terminated `exit_code: 0`)
- Real `fid` CLI: `/tmp/fid-task10 --project /tmp/fid-task10-runtime.wYHpfL --json activation prepare --input /tmp/fid-task10-runtime.wYHpfL/activation.json --idempotency-key runtime-activation-prepare` - PASS - output contained `"ok":true,"operation":"activation.prepare"` and five pending resources (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow 91.017s coverage: 55.3% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 63.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid/integration 25.182s` (terminated `exit_code: 0`)

## Deviations from spec

- The first required `golangci-lint run ./...` attempt ended `exit_code: 1` with `* unused: 1`: an unused activation helper, unchecked rows close, and redundant test type were introduced by this task's recovery convergence. Removed the unused helper, checked the close result, simplified the declaration, and reran the byte-identical command successfully with `0 issues.` (terminated `exit_code: 0`).
- Root assistance captured the initial resource-matrix failure and a test-construction correction in `project-docs/docs/phases/phase-22/tasks/task-10-delegate-handoff.md` at commit `523de68`; this task reran the exact matrix command and records its passing receipt above.
- Converged the retained DEC-192 task-1 assertion introduced by `72c5974f6698fa08a4fe1a24099b3a93b3882f1b` (phase-22/task-1): task 10's accepted fence removal displaced its activation-family rejection subject. `TestRecoverySchemaFreshInitAndReplay` now proves live/replay rejection for a genuinely unknown Event with an invalid effect and retains the corrupt Stop-code branch.

## Concerns, technical debt, or notes for the Tech Lead

- DEC-199 assigns typed activation-incomplete agreement detail enrichment and any new activation recovery CLI leaves to task 11; task 10 retains the accepted domain predicates and existing public leaves.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
