# Phase 22 - Task 11: Activation CLI and config boundary - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added the public `activation resource`, `activation recover`, and `activation recover-complete` CLI leaves with closed JSON results and error envelopes.
- Preserved distinct activation version, competing-prepare, and recovery-direction failures. Agreement failures now carry the activation ID, disagreeing components, and accepted `activation-incomplete` Reasons.
- Made the activation status Human view name `restore-prior` before commit and `complete-forward; irreversible commit` after it, matching the JSON recovery direction.
- Validated project configuration at activation mutation entry, so malformed configuration rejects consuming activation operations before effects while help and export validation remain usable.
- Exercised prepare, first-component stage, commit, and finalize durable boundaries through the built `fid` binary. Each mutation loses its response through a closed pipe, then a fresh process confirms status/recovery and byte-identical same-key replay.

## Files created or modified

- `cmd/fid/activation.go`, `cmd/fid/activation_test.go`, `cmd/fid/activation_resource_matrix_test.go`, `cmd/fid/activation_export_test.go`, and `cmd/fid/contract_surface_test.go`
- `internal/workflow/activation.go`, `internal/workflow/activation_recovery.go`, and `internal/workflow/activation_recovery_test.go`
- `internal/cli/codes.go` and `internal/cli/codes_test.go`
- `cmd/fid/integration/p21_conformance_test.go` and `cmd/fid/integration/p22_activation_recovery_test.go`
- `project-docs/docs/phases/phase-22/tasks/task-11-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-11-done.md`
- No contract or schema change was made. CF-035's accepted consuming-path proposal is recorded in `project-docs/docs/phases/phase-22/decisions.md:77` and implemented here.

## Validation Contract Coverage

- VAL-045: PASS - `cmd/fid/activation_resource_matrix_test.go:15` and `internal/workflow/activation_recovery_test.go:96` assert complete and incomplete resource-agreement results, including activation ID, components, accepted Reasons, no effects, and typed prior-agreement failure. `cmd/fid/activation_test.go:478` asserts paired Human/JSON pre-commit restore-prior and post-commit complete-forward/irreversible direction. `cmd/fid/activation_test.go:505` asserts closed not-found and invalid leaf variants. Root's independent paired corpus is recorded with attribution in `project-docs/docs/phases/phase-22/tasks/task-11-root-evidence.md`; `cmd/fid/integration/p22_activation_semantics_test.go:16` executes the closed version, direction, agreement, replay, and no-effect cases against fake and fid.
- VAL-046: PASS - `cmd/fid/integration/p22_activation_recovery_test.go:49` runs separate built-binary prepare, first-stage, committed, and finalized fixtures; its closed-pipe subprocess helper loses each selected mutation response and a fresh binary checks durable identity, permitted recovery direction, agreement/completion, same-key replay, and unchanged provider sentinel. The required tagged integration gate passed with `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	25.372s`.
- VAL-047: PASS - `cmd/fid/activation_export_test.go:13` corrupts an otherwise valid project config, proves consuming `activation prepare` rejects before journal effects, and proves `activation --help` plus `export validate` remain available.
- Convergence matrix: `TestLiveLeafCommandsMatchContractInventories` was introduced by `37b59d629b296b7b5383454170f9f017fdd94253` (phase-20/task-18); Task 11's accepted three public activation leaves invalidated its 75-leaf expectation, so it now asserts all 78 and runs in `go test -race -coverprofile=coverage.out ./...`. `TestCompleteActivationRecoveryRejectsMissingPriorResourceSet` and the tagged `p21ObserveNewActivationResources` helper originated with `33f85bbd6abffd0425a6a446780e8b0a37ce457a` (phase-22/task-10); DEC-199's Task 11 typed agreement result and public-resource-leaf replacement invalidated their generic-error/direct-domain expectations. The former runs in the full race gate and the latter runs in `go test -tags integration -race ./cmd/fid/integration/...`. `TestActivationLifecycleRejectsSkippedAndExtraTransitionsWithoutEvents` was introduced by `0444ccc458f3ac5a720724825f7f81589cb41e09` (phase-21/task-4); Task 11's accepted closed activation-invalid variant displaced its generic invalid-argument expectation, and it runs in the full race gate.

## Tests run

- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestActivationRecoveryAtEveryDurableBoundary$' -count=1` - FAIL - `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration	0.805s` (terminated `exit_code: 1`). First behavioral red: the public `activation resource` leaf was absent.
- `go test ./cmd/fid ./internal/workflow ./internal/cli -run 'TestActivation|TestReasonCodeDetailSchemaIsExactForEveryCode' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	0.026s` (terminated `exit_code: 0`).
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestActivationRecoveryAtEveryDurableBoundary$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.555s` (terminated `exit_code: 0`).
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22DualAdapterActivationConformance$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.834s` (terminated `exit_code: 0`). This root-owned independent fake/fid corpus receipt is cited from `task-11-root-evidence.md`; it is not presented as Developer-run evidence.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	82.053s	coverage: 55.0% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					63.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	25.372s` (terminated `exit_code: 0`)

## Deviations from spec

- No specification deviations. Required gate attempt history: the first `golangci-lint run ./...` attempt printed `1 issues.` (terminated `exit_code: 1`) because the typed prior-agreement refactor left an obsolete unused helper; the helper was removed and the byte-identical rerun printed `0 issues.` (terminated `exit_code: 0`). The first `go test -race -coverprofile=coverage.out ./...` attempt printed `FAIL	github.com/mbell-nzl/fid/internal/workflow	89.158s` (terminated `exit_code: 1`) because the retained live-leaf inventory omitted the three accepted activation leaves and a retained prior-resource expectation still required generic activation-invalid. The inventory was converged to 78 leaves and the assertion updated to the accepted typed activation-incomplete error; the byte-identical rerun passed as recorded above.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
