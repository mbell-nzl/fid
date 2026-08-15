# Phase 20 - Task 15: Format shim removal and whole-tree caller convergence - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: The red-first repository guard still reports remaining non-rejection `--format` callers that must be converted to direct `--json` invocations or their replacement command-family surface.
- Recovery attempts: 1
- Next action: Continue the caller-by-caller conversion, preserving each semantic assertion; run the focused guard to green, then the unchanged six-command gate.
- Human attention required: no

## Summary of what was implemented

- Added `TestLiveTestsContainNoUniversalFormatInvocation`, an AST-based source guard that exempts only named retained rejection tests and lists every remaining live caller.
- Removed `migratedReadTestArgs`, `containsFormatJSON`, `isLegacyIssueSection`, `legacyMigratedReadOutput`, and all argument/output rewriting from `cmd/fid/issue_test.go`.
- Began converting Issue-get fixture reads to direct canonical JSON in `mustGetIssueJSON`; caller convergence is incomplete and deliberately remains in this checkpoint.

## Files created or modified

- `cmd/fid/contract_surface_test.go`
- `cmd/fid/issue_test.go`
- `project-docs/docs/phases/phase-20/tasks/task-15-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-20/tasks/task-15-done.md` - durable same-task recovery checkpoint.

## Validation Contract Coverage

- VAL-029: BLOCKED - Legacy-leaf regression tests are unchanged, but the Task 15 caller convergence is incomplete; the required full gate has not run.
- VAL-031: BLOCKED - `TestLiveTestsContainNoUniversalFormatInvocation` is present at `cmd/fid/contract_surface_test.go:13`, but its focused execution still lists live callers and exits non-zero.

## Tests run

- `go test ./cmd/fid -run '^TestLiveTestsContainNoUniversalFormatInvocation$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestLiveTestsContainNoUniversalFormatInvocation$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`) - After shim deletion, remaining direct caller inventory is still red.

## Exact quality-gate results

- `gofmt -l cmd internal` - BLOCKED - Not run: caller-convergence implementation remains incomplete; no termination observed.
- `golangci-lint run ./...` - BLOCKED - Not run: caller-convergence implementation remains incomplete; no termination observed.
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - Not run: caller-convergence implementation remains incomplete; no termination observed.
- `go tool cover -func=coverage.out` - BLOCKED - Not run: prerequisite coverage command has not run; no termination observed.
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - Not run: caller-convergence implementation remains incomplete; no termination observed.
- `go build -o /tmp/fid ./cmd/fid` - BLOCKED - Not run: caller-convergence implementation remains incomplete; no termination observed.

## Deviations from spec

- This is a same-task recovery checkpoint, not a task close. The immutable red-first evidence is the focused source-guard result `FAIL` with observed `exit_code: 1`; it lists all remaining live caller locations. No validation assertion, quality gate, threshold, dependency, or protected task boundary changed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
