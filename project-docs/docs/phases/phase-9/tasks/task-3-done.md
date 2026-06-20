# Phase 9 Task 3 Completion Report

## Summary of what was implemented

- Added `fid concern list` with `--issue`, `--status`, and `--format` support,
  backed by a new `internal/issues.ListConcerns` read path.
- Added markdown rendering for concern lists so `json`, `toon`, `markdown`, and
  `yaml` all work for the new command.
- Added real integration coverage for concern rebuild parity and the full v0.1c
  closure smoke from ACs through concern raise/address/list.

## Files created or modified

- `cmd/fid/concern.go`
- `cmd/fid/concern_test.go`
- `cmd/fid/e2e_v01c_concern_test.go`
- `internal/issues/concern_list.go`
- `internal/issues/concern_list_test.go`
- `internal/format/markdown.go`
- `internal/issues/transition_rules_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-9/tasks/task-3-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
  - `fid concern list` now explicitly documents composable `--issue` +
    `--status` filters, empty-result exit-0 behavior, and the JSON shape
    including `issue_id` and `updated_at`.
- `project-docs/docs/phases/phase-9/tasks/task-3-done.md`

## Validation Contract Coverage

- `VAL-008`
  - Evidence: `internal/issues/concern_list_test.go` proves domain filtering and
    empty results; `cmd/fid/concern_test.go` proves `fid concern list` issue and
    status filters, combined filters, empty matches, and format parity across
    `json`, `toon`, `markdown`, and `yaml`.
- `VAL-009`
  - Evidence: `cmd/fid/e2e_v01c_concern_test.go` `TestV01CConcernRebuildParity`
    snapshots `concerns` rows plus `fid concern list`, `fid issue get`, and
    `fid issue list` output, deletes `.fid/fid.sqlite`, runs `fid rebuild
    --by developer`, and asserts exact parity afterward.
- `VAL-010`
  - Evidence: `cmd/fid/e2e_v01c_concern_test.go` `TestV01CClosureSmoke`
    exercises `ac add` -> `ac verify` -> `ac list` -> `issue complete` ->
    `concern raise` -> derived flag lights on `issue get` and `issue list` ->
    `concern address` -> derived flag clears -> `concern list`.
- Coverage thresholds
  - `cmd/fid`: `77.8%`
  - `internal/issues`: `81.8%`
  - `internal/journal`: `80.3%`
  - `internal/policy`: `85.0%`
  - `internal/workflow`: `95.5%`

## Tests run

- `RED` `go test ./internal/issues -run TestListConcernsFiltersByIssueAndStatus`
  - Expected compile failure before implementation: undefined
    `ListConcerns`/`ConcernListOptions`.
- `GREEN` `go test ./internal/issues -run TestListConcernsFiltersByIssueAndStatus`
  - PASS.
- `RED` `go test ./cmd/fid -run TestConcernListFiltersByIssueAndStatus`
  - Expected CLI failure before wiring: unknown `concern list` flags.
- `GREEN` `go test ./cmd/fid -run 'TestConcernList(FiltersByIssueAndStatus|DispatchesFormats)'`
  - PASS.
- `RED` `go test ./cmd/fid -run TestConcernListDispatchesFormats`
  - Expected markdown failure before renderer support:
    `unsupported markdown value type []issues.ConcernListItem`.
- `GREEN` `go test -tags integration ./cmd/fid -run 'TestV01C(ConcernRebuildParity|ClosureSmoke)'`
  - PASS.
- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
- `PASS` `go tool cover -func=coverage.out`
- `PASS` `go test -tags integration -race ./...`
- `PASS` `go build -o /tmp/fid ./cmd/fid`
- `N/A` `go test -bench . -benchmem ./...`
  - No new hot concern-read query path was added; Task 3 only added a standard
    filtered read command plus tests.

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
  - No files listed.
- `PASS` `golangci-lint run ./...`
  - Output: `0 issues.`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
  - `cmd/fid` coverage: `77.8%`
  - `internal/issues` coverage: `81.8%`
- `PASS` `go tool cover -func=coverage.out`
  - DEC-007 thresholds satisfied:
    `cmd/fid` `77.8%`, `internal/issues` `81.8%`, `internal/journal` `80.3%`,
    `internal/policy` `85.0%`, `internal/workflow` `95.5%`.
- `PASS` `go test -tags integration -race ./...`
  - All integration-tagged packages passed.
- `PASS` `go build -o /tmp/fid ./cmd/fid`
  - Build completed successfully.
- `N/A` `go test -bench . -benchmem ./...`
  - No new hot concern-read query path was added.

## Deviations from spec

- Expanded the nominal write set to `internal/format/markdown.go` so
  `fid concern list --format markdown` could satisfy `VAL-008`.
- Adjusted `internal/issues/transition_rules_test.go` to remove an existing
  concern-ID ordering assumption surfaced by the required integration gate.
  This was a test-only stability fix; no workflow or concern production logic
  changed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit creation.

## Commit subject

- Pending commit creation.
