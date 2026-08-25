# Phase 21 - Task 9: Retained actor-policy removal - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: The retained V1 test corpus still asserts superseded `--by`, public policy, V1 event, and direct-JSON behavior; the Task 9 convergence is not yet complete.
- Recovery attempts: 0
- Next action: Converge the remaining retained inventory writers and their tests to the accepted Task 9 `--actor`, Event-v2, fail-closed, and historical-attribution contracts, then rerun all five required commands byte-for-byte.
- Human attention required: no

## Summary of what was implemented

- Removed public registration of the V1 `fid policy` command while retaining its unregistered helper for enumerated legacy dev/completion callers.
- Added standalone project-local actor attribution and Event-v2 idempotency for workflow definition, workflow assignment, and acceptance-criterion addition.
- Made AC verification, concern judgment mutations, and maintenance vacuum return typed capability-unavailable failures; renamed the public stats grouping to `actor` and labelled it as historical attribution.
- Added a focused CLI conformance matrix for the converted public surfaces.

## Files created or modified

- `cmd/fid/concern.go`
- `cmd/fid/issue_ac.go`
- `cmd/fid/maintenance.go`
- `cmd/fid/policy.go`
- `cmd/fid/root.go`
- `cmd/fid/stats.go`
- `cmd/fid/task9_contract_test.go`
- `cmd/fid/unavailable.go`
- `cmd/fid/workflow.go`
- `internal/cli/codes.go`
- `internal/issues/ac.go`
- `internal/issues/idempotency_results.go`
- `internal/stats/stats.go`
- `internal/workflow/assign.go`
- `internal/workflow/attribution.go`
- `internal/workflow/define.go`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-9-done.md`

## Validation Contract Coverage

- VAL-028: BLOCKED - `TestTask9RetainedInventorySurface` in `cmd/fid/task9_contract_test.go` proves the converted workflow and AC entrypoints, but the first required focused package gate ended `FAIL\tgithub.com/mbell-nzl/fid/cmd/fid\t146.305s` while retained labels, links, requirements, activity, and other inventory coverage remains unconverged.
- VAL-029: BLOCKED - `cmd/fid/root.go` no longer registers `newPolicyCommand`, but the required package gate has not passed and the full generic-caller isolation proof remains unfinished.
- VAL-030: BLOCKED - `TestTask9RetainedInventorySurface` proves typed capability-unavailable results for `issue ac verify`, `concern raise`, and `maintenance vacuum`; continuation must add persistence-negative evidence for every fail-closed judgment mutation.
- VAL-031: BLOCKED - `TestTask9RetainedInventorySurface` proves `fid stats token --group-by actor` returns `items`; continuation must finish the accepted grouping matrix and its historical-attribution-only evidence.

## Tests run

- `go test -race ./cmd/fid -run '^TestTask9RetainedInventorySurface$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 2.072s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestTask9RetainedInventorySurface/AC_add_accepts_standalone_actor_and_idempotency$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 1.913s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL\tgithub.com/mbell-nzl/fid/cmd/fid\t146.305s` (terminated `exit_code: 1`; retained V1 contract tests require in-scope convergence)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; the first attempt also observed the temporarily removed unregistered policy helper before it was restored, and retained V1 tests still require convergence)
- `go tool cover -func=coverage.out` - PASS - `68.3%` (terminated `exit_code: 0`)

## Deviations from spec

- No accepted-contract deviations.
- Required-gate first attempts are preserved above: the focused package test failed with `FAIL\tgithub.com/mbell-nzl/fid/cmd/fid\t146.305s` (terminated `exit_code: 1`) because retained tests still expect replaced V1 behavior; `golangci-lint run ./...` first failed with `* unused: 6` (terminated `exit_code: 1`) after intermediate migration symbols became unused, those symbols were removed, and its byte-identical rerun later passed with `0 issues.` (terminated `exit_code: 0`); coverage failed with `FAIL` (terminated `exit_code: 1`) while the temporary removal of the unregistered policy helper and retained V1 test expectations were being repaired. No gate, assertion, scope, or accepted interface was weakened. A continuation must rerun each affected command byte-identically after completing the in-scope migration.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `e51afe3c3e0ea157ee960fb643996453250a37c6`

## Commit subject

- `feat(convergence): checkpoint retained inventory actor migration`
