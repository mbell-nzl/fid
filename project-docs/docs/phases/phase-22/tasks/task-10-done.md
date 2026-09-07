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
- `cmd/fid/activation.go`, `cmd/fid/activation_test.go`, `cmd/fid/activation_binding_test.go`, `cmd/fid/activation_export_test.go`, `cmd/fid/activation_resource_matrix_test.go`, and `cmd/fid/integration/p21_conformance_test.go`
- `internal/evidence/registration.go`, `internal/cli/codes.go`, and their matching tests
- `internal/journal/project_effects.go`, `internal/journal/checkpoint_projection_test.go`, and `internal/journal/recovery_schema_test.go`
- Retained surface convergence tests in `cmd/fid/contract_surface_test.go`, `cmd/fid/envelope_test.go`, and `cmd/fid/gate_test.go`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-10-done.md`

## Validation Contract Coverage

- VAL-040: PASS - `TestActivationFreezesFirstAndUpgradeBindingsAcrossStagingAndRecovery` proves first and upgrade predecessor/resource bindings remain immutable through staging, observation, restore, and completion; `TestActivationCommitRejectsMissingMismatchedAndMixedComponentDigests`, `TestStagedActivationContentCannotBeReplacedUnderANewKey`, and `TestActivationComponentResponseDigestsRemainSortedAfterSmallerSecondStage` prove exact sorted staged component selection and immutable content; `TestUpgradeForwardRecoverySelectsOnlyNewFinalTrust` proves forward recovery supersedes the predecessor and selects only the new trust.
- VAL-041: PASS - `TestCommittedActivationCanOnlyRecoverForward` proves the public committed-abort direction rejection, while `TestActivationRecoverySelectionRaisesAndCompletionDisposesIncompleteStop` proves pre-commit restore and rejects a fresh terminal selection.
- VAL-042: PASS - `TestActivationRequiresEveryResourceAtCommitAndFinalize` covers each missing, null, and mismatched backend, kit, overlay, project, and provider state at both commit and finalization; `TestActivationCompleteResourceAgreementFinalizesTrust` supplies the all-agree control.
- VAL-043: PASS - `TestIncompleteUpgradeBlocksNormalWorkWhileInspectionAndRestoreRemainAvailable` rejects trusted actor resolution and a real `gate register` mutation while allowing status and restore; `TestIncompleteActivationAllowsExportWithoutChangingActivation` proves real export create/validate leave the nonterminal activation and Event count unchanged.
- VAL-044: PASS - `TestActivationRecoveryPreservesDirectionAndIdentityAcrossReplay` and `TestActivationRecoverySelectionRaisesAndCompletionDisposesIncompleteStop` prove one selected Event identity, immutable direction, completion, and response-loss replay after completion.

## Tests run

- `go test ./cmd/fid -run '^TestCommittedActivationCanOnlyRecoverForward$' -count=1` - RED then PASS - `FAIL` (terminated `exit_code: 1`); later `0.188s` (terminated `exit_code: 0`)
- `go test ./internal/workflow -run '^TestActivationRecoverySelectionRaisesAndCompletionDisposesIncompleteStop$' -count=1` - PASS - `0.160s` (terminated `exit_code: 0`)
- `go test ./internal/workflow -run '^TestCompleteActivationRecoveryRejectsMissingPriorResourceSet$' -count=1` - PASS - `0.087s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestActivation(RequiresEveryResourceAtCommitAndFinalize|CompleteResourceAgreementFinalizesTrust)$' -count=1` - PASS - `1.397s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^Test(ActivationFreezesFirstAndUpgradeBindingsAcrossStagingAndRecovery|ActivationCommitRejectsMissingMismatchedAndMixedComponentDigests|StagedActivationContentCannotBeReplacedUnderANewKey)$' -count=1` - RED then PASS - `FAIL` (terminated `exit_code: 1`); later `0.672s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestUpgradeForwardRecoverySelectsOnlyNewFinalTrust$' -count=1` - PASS - `0.503s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestActivationComponentResponseDigestsRemainSortedAfterSmallerSecondStage$' -count=1` - RED then PASS - `FAIL` (terminated `exit_code: 1`); later `0.103s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/...` - PASS - `11.407s` (terminated `exit_code: 0`)
- Real `fid` CLI: `/tmp/fid-task10 --project /tmp/fid-task10-runtime.wYHpfL --json activation prepare --input /tmp/fid-task10-runtime.wYHpfL/activation.json --idempotency-key runtime-activation-prepare` - PASS - output contained `"ok":true,"operation":"activation.prepare"` and five pending resources (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `coverage: 55.3% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `63.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `24.661s` (terminated `exit_code: 0`)

## Deviations from spec

- The first required `golangci-lint run ./...` attempt ended `exit_code: 1` with `* unused: 1`: an unused activation helper, unchecked rows close, and redundant test type were introduced by this task's recovery convergence. Removed the unused helper, checked the close result, simplified the declaration, and reran the byte-identical command successfully with `0 issues.` (terminated `exit_code: 0`).
- Root assistance captured the initial resource-matrix failure and a test-construction correction in `project-docs/docs/phases/phase-22/delegate-handoff.md` at commit `523de68`; this task reran the exact matrix command and records its passing receipt above.
- The same handoff preserves the VAL-040 binding-test red receipt: its exact focused command ended `exit_code: 1` with `FAIL` because a successor was finalized before its predecessor was superseded, violating the unique finalized activation constraint. This task ordered predecessor supersession before successor finalization in the same Event and reran the byte-identical binding command successfully.
- `TestActivationComponentResponseDigestsRemainSortedAfterSmallerSecondStage` first ended `exit_code: 1` with `FAIL`: staging appended a lexicographically smaller second digest without sorting the response snapshot. The staged snapshot now sorts component digests before canonical result serialization; its focused rerun passed.
- Initial implementation checkpoint: `33f85bbd6abffd0425a6a446780e8b0a37ce457a`; the final binding, effect-order, and response-order repair is recorded in the task implementation commit below.
- Converged retained callers/assertions displaced by this task's accepted input and activation-incomplete behavior: `cmd/fid/integration/p21_conformance_test.go` was introduced by `a99b4270f70e98111b5bae43c54ed92aed7078c6` (phase-21/task-22), `cmd/fid/envelope_test.go` by `56dad8fa0dcb9b79f7df2a31b356f1b2eb8a694e` (phase-21/task-6), `cmd/fid/gate_test.go` by `0c4c04ea0c7d0786e2e6cf88a7573d8a3105aef5` (phase-21/task-5), and `cmd/fid/contract_surface_test.go` by `ad693228a9e1df57cab1d8f37159b47ecdb08e29` (phase-20/task-15). The final required race and tagged integration gates execute each converged subject.
- Converged the retained DEC-192 task-1 assertion introduced by `72c5974f6698fa08a4fe1a24099b3a93b3882f1d` (phase-22/task-1): task 10's accepted fence removal displaced its activation-family rejection subject. `TestRecoverySchemaFreshInitAndReplay` now proves live/replay rejection for a genuinely unknown Event with an invalid effect and retains the corrupt Stop-code branch.

## Concerns, technical debt, or notes for the Tech Lead

- DEC-199 assigns typed activation-incomplete agreement detail enrichment and any new activation recovery CLI leaves to task 11; task 10 retains the accepted domain predicates and existing public leaves.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `03cc6c764412d44d53c526f6244af89f22dd9eca`

## Commit subject

- `fix(workflow): preserve activation recovery bindings`
