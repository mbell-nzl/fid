# Phase 22 - Task 3: Lease lifecycle and executor observations - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: session-continuation
- Condition: The accepted task contract requires the exact Run `duplicate-entity` and `run-version-conflict` CLI envelopes, but `internal/cli/codes.go` is outside task 3's named source scope. Its current `duplicate-entity` definition excludes `SubjectRun`, and it has no `run-version-conflict` code definition.
- Recovery attempts: 1
- Next action: Obtain a bounded task-3 scope amendment that authorizes `internal/cli/codes.go` and its direct tests, then resume from this red-test checkpoint without changing accepted error semantics.
- Human attention required: no

## Summary of what was implemented

- Added the required first failing regression `TestExpiredClaimWithUnresolvedExecutorCannotBeReacquired`; no product implementation was made because the accepted CLI error contract cannot be rendered within the task's allowed files.

## Files created or modified

- `cmd/fid/generic_claim_test.go`
- `project-docs/docs/phases/phase-22/tasks/task-3-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-3-done.md`

## Validation Contract Coverage

- VAL-008: BLOCKED - `TestExpiredClaimWithUnresolvedExecutorCannotBeReacquired` is red at `cmd/fid/generic_claim_test.go:250`; `claim.expire` is not yet registered.
- VAL-009: BLOCKED - implementation cannot reach the required real-process race while exact Task 3 error envelopes are outside source scope.
- VAL-010: BLOCKED - implementation cannot reach trusted observer admission while exact Task 3 error envelopes are outside source scope.
- VAL-011: BLOCKED - implementation cannot reach current-version runtime fencing while exact Task 3 error envelopes are outside source scope.
- VAL-012: BLOCKED - the accepted reused-executor case requires `duplicate-entity` with a Run subject, but `internal/cli/codes.go` excludes `SubjectRun` and is out of scope.

## Tests run

- `go test ./cmd/fid -run '^TestExpiredClaimWithUnresolvedExecutorCannotBeReacquired$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  - Attempt 1: `generic_claim_test.go:250: expire Claim = main.commandExecutionResult{stdout:"", stderr:"", err:(*pflag.NotExistError)(0x814ac82f8c0)}, want explicit expiry success`

## Exact quality-gate results

- `gofmt -l cmd internal` - BLOCKED - not run; no conforming in-scope implementation path exists while required CLI code definitions are excluded from task scope.
- `golangci-lint run ./...` - BLOCKED - not run; no conforming in-scope implementation path exists while required CLI code definitions are excluded from task scope.
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - not run; no conforming in-scope implementation path exists while required CLI code definitions are excluded from task scope.
- `go tool cover -func=coverage.out` - BLOCKED - not run; its required coverage profile cannot exist before the blocked test suite gate.
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - not run; no conforming in-scope implementation path exists while required CLI code definitions are excluded from task scope.

## Deviations from spec

- The task is unfinished. The accepted public interface requires closed error definitions that are not writable under task 3's declared file scope; no contract, gate, tolerance, or test was changed.

## Concerns, technical debt, or notes for the Tech Lead

- Provenance: `internal/cli/codes.go` was introduced by `ab4ac8f288b3f180b698bab8924b92f375d09eb0` and its current generic Run lifecycle integration by `5ca732ff163ee1d31fc155d3d9ebb3dac1501838`, both predate Phase 22. The task's accepted `public-interface.md` rows require additions to that closed definition surface, but the source file is omitted from task 3 scope.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
