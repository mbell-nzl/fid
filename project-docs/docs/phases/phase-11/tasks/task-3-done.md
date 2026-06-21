## Summary of what was implemented

- Added the root `fid claim` command with the read-only `fid claim list` subcommand.
- Added `internal/claims.List` with `--by` filtering over stored `claimed_by`, stale/live classification, and deterministic `ORDER BY rowid` default ordering.
- Added markdown rendering support for claim-list rows and CLI/domain tests covering format parity, stale filtering, rowid order, rebuild stability, human display, unknown-role empty filtering, and read-only behavior.

## Files created or modified

- `cmd/fid/claim.go`
- `cmd/fid/claim_test.go`
- `cmd/fid/root.go`
- `internal/claims/list.go`
- `internal/claims/list_test.go`
- `internal/claims/shared.go`
- `internal/format/markdown.go`
- `project-docs/docs/phases/phase-11/tasks/task-3-done.md`

## Documentation updated

- Added this completion report at `project-docs/docs/phases/phase-11/tasks/task-3-done.md`.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: verified-no-change for `fid claim list`.
- `project-docs/docs/phases/phase-11/decisions.md`: no update required because implementation matched DEC-027 and the existing CLI contract.

## Validation Contract Coverage

- `VAL-008` PASS for the Task 3 slice: `TestClaimListDefaultsToJSONAndListsLiveClaims` and `TestClaimListDispatchesFormatsAndFiltersByStoredRole` prove required fields, default JSON output, `--by` filtering over stored `claimed_by`, and json/toon/markdown/yaml parity.
- `VAL-009` PASS for the Task 3 slice: `TestListFiltersClaimsByRoleAndStalenessAndOrdersByRowID`, `TestClaimListStaleFilterIsReadOnlyAndUsesRowIDOrder`, and `TestClaimListOrderIsStableAcrossRebuild` prove stale-only membership, live exclusion from default output, `ORDER BY rowid`, and rebuild-stable ordering for the same journal.
- `VAL-010` PASS for the Task 3 read-surface slice: `TestListFiltersClaimsByRoleAndStalenessAndOrdersByRowID` and `TestClaimListDispatchesFormatsAndFiltersByStoredRole` prove a human-held claim lists as `claimed_by: human`. The transition-blocking half of `VAL-010` remains assigned to Task 4 per the validation contract.

## Tests run

- `go test ./internal/claims -run TestListFiltersClaimsByRoleAndStalenessAndOrdersByRowID -count=1`
- `go test ./cmd/fid -run TestClaimListDefaultsToJSONAndListsLiveClaims -count=1`
- `go test ./cmd/fid -run 'TestClaimListDispatchesFormatsAndFiltersByStoredRole|TestClaimListStaleFilterIsReadOnlyAndUsesRowIDOrder|TestClaimListOrderIsStableAcrossRebuild' -count=1`
- `go test ./cmd/fid`
- `go test ./internal/claims`
- `go test ./internal/format`

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS
- `golangci-lint run ./...` — PASS (`0 issues.`)
- `go test -race -coverprofile=coverage.out ./...` — PASS
- `go tool cover -func=coverage.out` — PASS
- Coverage snapshot: `cmd/fid` 77.8%, `internal/claims` 83.0%, `internal/issues` 81.9%, `internal/journal` 80.3%, `internal/policy` 81.2%, `internal/workflow` 95.5%
- `go test -tags integration -race ./...` — PASS
- `go test -bench . -benchmem ./...` — PASS
- Benchmark command was required because Task 3 added a new list query/read path.
- `go build -o /tmp/fid ./cmd/fid` — PASS

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending reconciliation after task commit.

## Commit subject

- Pending reconciliation after task commit.
