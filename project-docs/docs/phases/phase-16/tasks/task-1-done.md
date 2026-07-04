# Task 1 Completion Report

## Summary of what was implemented
- Added `internal/activity` with `List(options)` over `activity_log`, including `--on <type>/<id>`, `--since`, and `--event` filters, projection-freshness enforcement, chronological ordering (`occurred_at`, `id`), and JSON `details` deserialisation back to structured values.
- Added `fid activity list` in `cmd/fid/activity.go`, registered it on the root command, and supported `json`, `toon`, `yaml`, and markdown output for the new read surface.
- Updated the task-owned architecture docs to record `internal/activity` as a domain read package.
- `--since` accepts both `YYYY-MM-DD` and RFC3339 input. Date-only values are normalised to UTC midnight before querying, so on/after-boundary semantics match VAL-003.

## Files created or modified
- `internal/activity/activity.go`
- `internal/activity/activity_test.go`
- `cmd/fid/activity.go`
- `cmd/fid/activity_test.go`
- `cmd/fid/root.go`
- `project-docs/specs/02_TRD.md`
- `AGENTS.md`
- `project-docs/docs/phases/phase-16/tasks/task-1-done.md`

## Documentation updated
- `project-docs/specs/02_TRD.md`
  Added `internal/activity` to the architecture diagram and module-boundary table; updated the document timestamp to `2026-07-05`.
- `AGENTS.md`
  Added `internal/activity/` to the core/domain logic package list.

## Validation Contract Coverage
- `VAL-001`
  Covered by `TestListReturnsChronologicalEntriesWithStructuredDetails`, `TestActivityListDispatchesFormats`, and the built-binary smoke run of `fid activity list --on issue/<id> --format json`.
- `VAL-002`
  Covered by `TestListFiltersByOnSinceAndEvent/on filter` and `TestActivityListAppliesFilters/on filter`.
- `VAL-003`
  Covered by `TestListFiltersByOnSinceAndEvent/since boundary includes exact match`, `TestListFiltersByOnSinceAndEvent/since date only`, and `TestActivityListAppliesFilters/since filter`.
- `VAL-004`
  Covered by `TestListFiltersByOnSinceAndEvent/event filter` and `TestActivityListAppliesFilters/event filter`.

## Tests run
- `go test ./internal/activity -run TestListReturnsChronologicalEntriesWithStructuredDetails`
  Red: failed to build with `undefined: List`, `undefined: ListOptions`, and `undefined: Entry`.
- `go test ./internal/activity -run TestListReturnsChronologicalEntriesWithStructuredDetails`
  Green: PASS.
- `go test ./internal/activity`
  Green: PASS.
- `go test ./cmd/fid -run TestActivityList`
  Red: failed with `unknown command "activity" for "fid"`.
- `go test ./cmd/fid -run TestActivityList`
  Green: PASS.
- `/tmp/fid init`
  PASS in `/tmp/fid-activity-YV93Dj`; workspace initialised.
- `/tmp/fid issue create --kind task --title "Activity smoke" --by developer`
  PASS; created `FID-ACTIVITY-YV93DJ-ec98ce`.
- `/tmp/fid issue claim FID-ACTIVITY-YV93DJ-ec98ce --by developer`
  PASS; live claim recorded so an existing activity-producing writer could run.
- `/tmp/fid issue transition FID-ACTIVITY-YV93DJ-ec98ce --to in_progress --by developer`
  PASS; status changed from `todo` to `in_progress`.
- `/tmp/fid activity list --on issue/FID-ACTIVITY-YV93DJ-ec98ce --format json`
  PASS; returned one `status_changed` entry with structured `details` (`from: todo`, `to: in_progress`).

## Exact quality-gate results
- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
  `cmd/fid`: `76.8%`
  `internal/activity`: `81.1%`
- `PASS` `go tool cover -func=coverage.out`
- `NOT REQUIRED` `go test -tags integration -race ./...`
  No integration-tagged tests were added in Task 1.
- `PASS` `go build -o /tmp/fid ./cmd/fid`

## Deviations from spec
- Markdown output for `fid activity list` is rendered in `cmd/fid/activity.go` instead of extending `internal/format`'s markdown type switch. This kept the change inside the task-owned files while still satisfying the four-format acceptance contract.

## Concerns, technical debt, or notes for the Tech Lead
- None.

## Dependency-budget note
- No dependency changes.

## Commit hash
- Pending commit.

## Commit subject
- Pending commit.
