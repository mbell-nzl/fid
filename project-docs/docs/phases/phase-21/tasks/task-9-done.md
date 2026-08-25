# Phase 21 - Task 9: Retained actor-policy removal - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: The retained internal Task 9 packages pass in the required race gate, and stats/root-capability convergence advanced, but the full command gate still contains CLI V1 fixtures and replacement-bound dev tests that expect `--by`, `policy show`, or Phase-23 judgment writers.
- Recovery attempts: 5
- Next action: Migrate the remaining retained CLI fixture subjects to exact actor/key/Event-v2 and policy-absence/fail-closed contracts; converge replacement-bound dev expectations without restoring judgment writers, then rerun all five required commands byte-for-byte.
- Human attention required: no

## Summary of what was implemented

- Removed public registration of the V1 `fid policy` command while retaining its unregistered helper for enumerated legacy dev/completion callers.
- Added standalone project-local actor attribution and Event-v2 idempotency for workflow definition, workflow assignment, and acceptance-criterion addition.
- Made AC verification, concern judgment mutations, and maintenance vacuum return typed capability-unavailable failures; renamed the public stats grouping to `actor` and labelled it as historical attribution.
- Added a focused CLI conformance matrix for the converted public surfaces.
- Converted label add, issue-label add/remove, link add/remove, and requirement sync to standalone actor/key Event-v2 paths; started canonical activity event normalization and transport migration.
- Registered `concern list` as an accepted project/JSON root-flag leaf and strengthened stats CLI assertions to the public `actor_attribution` envelope with `--group-by actor|epic`.

## Files created or modified

- `cmd/fid/concern.go`
- `cmd/fid/capabilities_test.go`
- `cmd/fid/contract_surface_test.go`
- `cmd/fid/issue_ac.go`
- `cmd/fid/maintenance.go`
- `cmd/fid/policy.go`
- `cmd/fid/root.go`
- `cmd/fid/stats.go`
- `cmd/fid/stats_test.go`
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
- `internal/activity/activity.go`
- `internal/labels/attribution.go`
- `internal/labels/add.go`
- `internal/labels/types.go`
- `internal/links/add.go`
- `internal/links/attribution.go`
- `internal/links/idempotency_results.go`
- `internal/links/remove.go`
- `internal/links/types.go`
- `internal/requirements/attribution.go`
- `internal/requirements/sync.go`
- `internal/requirements/types.go`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-9-done.md`

## Validation Contract Coverage

- VAL-028: BLOCKED - `TestTask9RetainedInventorySurface` in `cmd/fid/task9_contract_test.go` now has red rows for label add, issue-label add, link add, and requirement sync; the current required package gate ended `FAIL	github.com/mbell-nzl/fid/cmd/fid	147.763s` while retained assertion migration and remaining families remain unconverged.
- VAL-029: BLOCKED - `cmd/fid/root.go` no longer registers `newPolicyCommand`, but the required package gate has not passed and the full generic-caller isolation proof remains unfinished.
- VAL-030: BLOCKED - `TestTask9RetainedInventorySurface` proves typed capability-unavailable results for `issue ac verify`, `concern raise`, and `maintenance vacuum`; continuation must add persistence-negative evidence for every fail-closed judgment mutation.
- VAL-031: BLOCKED - `TestTask9RetainedInventorySurface` proves `fid stats token --group-by actor` returns `items`; continuation must finish the accepted grouping matrix and its historical-attribution-only evidence.

- Attempt 2 update: VAL-030 now has direct no-writer evidence in `internal/issues/judgment_unavailable_test.go` and `internal/maintenance/vacuum_unavailable_test.go`; the remaining CLI persistence-negative matrix is in scope for recovery. VAL-028, VAL-029, and VAL-031 remain PARTIAL because the required gate still fails on retained V1 fixture expectations.

- Attempt 4 update: VAL-028 internal workflow, links, requirements, and activity assertions now pass as Event-v2/actor/key and normalized-read coverage. VAL-030 continues to retain typed no-writer assertions; the incomplete CLI migration still prevents COMPLETE for VAL-028 through VAL-031.

- Attempt 5 update: VAL-028 now includes the exact retained root inventory in `TestLiveLeafCommandsMatchContractInventories` (`cmd/fid/contract_surface_test.go`) and an accepted `concern list` root capability. VAL-029 and VAL-030 retain the policy-absence and capability-unavailable coverage in `TestDeferredCapabilitySurfacesStayUnavailableWithoutSideEffects` (`cmd/fid/capabilities_test.go`). VAL-031 CLI assertions now use `TestStatsTokensDefaultsToJSONByActorAttribution`, `TestStatsTokensUsesDirectLegacyJSONAndEpicFilter`, and `TestStatsCostUsesDirectLegacyJSONAndExcludesNullCostsFromSum` (`cmd/fid/stats_test.go`) to check public envelopes and `--group-by`; the required command gate remains BLOCKED by other unmigrated CLI fixtures.

## Tests run

- `go test ./cmd/fid` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 11.645s` (terminated `exit_code: 1`; remaining V1 CLI fixtures and policy/judgment expectations are enumerated in the output)
- `go run ./cmd/fid --json concern list` - PASS - `[]` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestTask9RetainedInventorySurface$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 2.072s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestTask9RetainedInventorySurface/AC_add_accepts_standalone_actor_and_idempotency$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 1.913s` (terminated `exit_code: 0`)
- `go test ./internal/labels ./internal/issues` - FAIL - `FAIL	github.com/mbell-nzl/fid/internal/issues	6.617s` (terminated `exit_code: 1`; retained AC and issue-label Event-v1 assertions remain)
- `go test ./internal/requirements ./internal/links` - FAIL - `FAIL	github.com/mbell-nzl/fid/internal/links	1.355s` (terminated `exit_code: 1`; retained role and activity-name assertions remain)

## Exact quality-gate results

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL	github.com/mbell-nzl/fid/cmd/fid	146.305s` (terminated `exit_code: 1`; retained V1 contract tests require in-scope convergence)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; the first attempt also observed the temporarily removed unregistered policy helper before it was restored, and retained V1 tests still require convergence)
- `go tool cover -func=coverage.out` - PASS - `68.3%` (terminated `exit_code: 0`)
- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL	github.com/mbell-nzl/fid/cmd/fid	147.763s` (terminated `exit_code: 1`; retained test migration remains in scope)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; retained test migration remains in scope)
- `go tool cover -func=coverage.out` - PASS - `68.3%` (terminated `exit_code: 0`)

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL github.com/mbell-nzl/fid/cmd/fid 130.786s` (terminated `exit_code: 1`; workflow/link/activity and retained CLI V1 assertion migration remains)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - BLOCKED - `* unused: 2` (terminated `exit_code: 1`; obsolete helper symbols remained after writer isolation)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; the full suite still contains retained V1 test expectations)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 64.8%` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL github.com/mbell-nzl/fid/cmd/fid 135.581s` (terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 66.5%` (terminated `exit_code: 0`)

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL github.com/mbell-nzl/fid/cmd/fid 131.849s` (terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 66.6%` (terminated `exit_code: 0`)

## Deviations from spec

- No accepted-contract deviations.
- Required-gate first attempts are preserved above: the focused package test failed with `FAIL	github.com/mbell-nzl/fid/cmd/fid	146.305s` (terminated `exit_code: 1`) because retained tests still expect replaced V1 behavior; `golangci-lint run ./...` first failed with `* unused: 6` (terminated `exit_code: 1`) after intermediate migration symbols became unused, those symbols were removed, and its byte-identical rerun later passed with `0 issues.` (terminated `exit_code: 0`); coverage failed with `FAIL` (terminated `exit_code: 1`) while the temporary removal of the unregistered policy helper and retained V1 test expectations were being repaired. No gate, assertion, scope, or accepted interface was weakened. A continuation must rerun each affected command byte-identically after completing the in-scope migration.
- Recovery attempt 1 added label, issue-label, link, requirement, and activity conversion. Its required package gate failed with `FAIL	github.com/mbell-nzl/fid/cmd/fid	147.763s` (terminated `exit_code: 1`) and coverage failed with `FAIL` (terminated `exit_code: 1`) because retained V1 `--by`, policy-show, direct-JSON, role, snake-case activity, and Event-v1 assertions have not yet been migrated. `gofmt -l cmd internal` and `golangci-lint run ./...` both passed with observed `exit_code: 0`. No assertion, threshold, command, policy, or public contract was weakened; the next continuation must preserve these failures while migrating their accepted replacements.

- Recovery attempt 2 made direct judgment and vacuum entries unavailable before policy, projection, lock, Event, or row work. The writer-oriented tests for those exact subjects were adapted to assert typed unavailability before the same writer inputs can be evaluated; this replaces obsolete mutation-success coverage with the accepted fail-closed no-writer contract, not a narrowed assertion. Its required package gate failed with `FAIL github.com/mbell-nzl/fid/cmd/fid 130.786s` (terminated `exit_code: 1`) and coverage failed with `FAIL` (terminated `exit_code: 1`) because broad retained V1 CLI and internal test fixtures remain. The first lint attempt failed with `* unused: 2` (terminated `exit_code: 1`) due to obsolete helper symbols after writer isolation; removing those symbols and separating inactive legacy implementation produced the byte-identical later PASS `0 issues.` (terminated `exit_code: 0`). No gate, assertion semantics, threshold, policy, or accepted public contract changed.
- Recovery attempt 4 migrated the retained internal workflow, links, requirements, activity, and CLI-code tests to Event-v2/idempotency, standalone attribution, normalized activity, and capability-code contracts. Its focused internal package set was green before the required gate. The required package and coverage gates remain red on unmigrated CLI V1 fixtures and replacement-bound dev tests that invoke removed Phase-23 writers; no test subject was removed or weakened, and no V1 policy or judgment writer was restored.
- Recovery attempt 5 added the missing accepted `concern list` root capability and updated stats fixtures from internal `TokensRow`/`CostRow` serialization to exact public success envelopes with `actor_attribution`, numeric aggregate fields, and accepted `--group-by` values. The removed internal-row expectation lines were obsolete because they asserted an implementation type (`Role`) rather than the accepted public contract; their replacements assert every returned public field and retain all grouping/filter aggregation subjects. The first required package gate failed with `FAIL github.com/mbell-nzl/fid/cmd/fid 131.849s` (terminated `exit_code: 1`) and the coverage gate failed with `FAIL` (terminated `exit_code: 1`) due to remaining retained V1 fixtures plus replacement-bound dev tests that call fail-closed judgment writers. Formatting and lint passed. No assertion, threshold, command, policy, or accepted public contract changed.

## Concerns, technical debt, or notes for the Tech Lead

- The remaining package failures include tests under `cmd/fid` and `internal/dev` whose subjects span retained V1 fixture migration and replacement-bound dev completion behavior. Continue only the accepted fixture/expectation convergence; do not restore `policy show` or Phase-23 judgment writers.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `e4f486bbadb4fa396e4859a0000dff70116bbf7d`

## Commit subject

- `test(convergence): checkpoint task-9 CLI contract migration`
