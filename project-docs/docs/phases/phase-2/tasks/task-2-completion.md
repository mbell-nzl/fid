## Summary of what was implemented

Implemented Phase 2 Task 2 by adding migration `0003` for the `issues_enriched` view plus the required concern and triage indexes, bumping `internal/db.ProjectionSchemaVersion` from `2` to `3`, and switching `internal/issues.Get` to read `concerns_raised`, `open_concerns_count`, and `triage_status` from `issues_enriched` instead of hard-coded defaults.

The task also added package-level `internal/issues` coverage for enriched concern and triage states, `internal/db` coverage for view correctness and `EXPLAIN QUERY PLAN` index use, and `cmd/fid` CLI JSON coverage proving the enriched fields survive delete-DB plus `fid rebuild`.

## Files created or modified

- `migrations/0003_v01_enriched_view.up.sql`
- `migrations/0003_v01_enriched_view.down.sql`
- `internal/db/projection_meta.go`
- `internal/db/migrate_v01_tables_test.go`
- `internal/db/issues_enriched_test.go`
- `internal/issues/issue.go`
- `internal/issues/issue_test.go`
- `cmd/fid/issue_test.go`
- `project-docs/docs/phases/phase-2/tasks/task-2-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-2/tasks/task-2-completion.md`

## Validation Contract Coverage

- `VAL-005` PASS: `internal/db/issues_enriched_test.go` proves `issues_enriched` returns `open_concerns_count=2` and `concerns_raised=true` for open concerns, plus `0`/`false` for addressed-only and absent-concern states.
- `VAL-006` PASS: `internal/db/issues_enriched_test.go` proves `issues_enriched.triage_status` maps the latest verdict to `'resumed'`, `'blocked'`, and `NULL`; `internal/issues/issue_test.go` proves `Get` surfaces latest `RESUME`, latest `BLOCK`, and no-triage states through the package API.
- `VAL-007` PASS: `internal/issues/issue_test.go` and `cmd/fid/issue_test.go` prove `fid issue get --format json` now surfaces the enriched fields from `issues_enriched`, reports `0`/`false`/`null` when no concern or triage exists, and returns identical JSON after deleting `.fid/fid.sqlite` and running `fid rebuild`.
- `VAL-008` PASS: `internal/db/issues_enriched_test.go` asserts `EXPLAIN QUERY PLAN` references `idx_concerns_issue_status` and `idx_triage_issue_created` and does not fall back to scanning the base `concerns` or `triage_events` tables.

## Tests run

- `go test ./internal/issues -run TestGetReturnsEnrichedConcernFields -count=1`
  Red: failed before implementation with `Get() concerns_raised = false, want true`.
  Green: passed after adding migration `0003` and switching `Get` to `issues_enriched`.
- `go test ./internal/issues -count=1`
- `go test ./internal/db -count=1`
- `go test ./cmd/fid -count=1`

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS
- `golangci-lint run ./...` — PASS (`0 issues.`)
- `go test -race -coverprofile=coverage.out ./...` — PASS
- `go tool cover -func=coverage.out` — PASS (`total: 72.6%`; `internal/issues: 85.1%`; `internal/journal: 80.6%`; `cmd/fid: 70.5%`)
- `go test -tags integration -race ./...` — N/A; Task 2's CLI/rebuild parity coverage lives in untagged `cmd/fid/issue_test.go`, so the QA-plan condition for the integration-tag command did not apply.
- `go build -o /tmp/fid ./cmd/fid` — PASS

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `00d1820`

## Commit subject

- `feat(issue): read enriched issue state from issues_enriched`
