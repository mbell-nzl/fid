## Summary of what was implemented

- Added thin `fid issue create` and `fid issue get` Cobra adapters in `cmd/fid/issue.go`, wired into the root command, with stderr confirmations and clear error rendering.
- Added `internal/issues` to build and append a single `issue.created` event, project it into `issues`, and read the projected issue back in the Design Plan JSON shape.
- Added the `internal/format` seam with JSON output and explicit "not yet supported" errors for `toon`, `markdown`, and `yaml`.
- Added minimal `internal/policy` role-existence validation backed by the journal so `fid issue create` and `fid rebuild` both enforce `--by <role>` even when `fid.sqlite` is absent.

## Files created or modified

- `cmd/fid/issue.go`
- `cmd/fid/issue_test.go`
- `cmd/fid/rebuild.go`
- `cmd/fid/root.go`
- `internal/format/format.go`
- `internal/issues/issue.go`
- `internal/policy/policy.go`
- `project-docs/docs/phases/phase-1/tasks/task-6-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-1/tasks/task-6-completion.md`

## Validation Contract Coverage

- `VAL-006`:
  `TestIssueCreateThenGetReturnsJSONByDefaultAndWithFormatFlag` proves a minimal issue is created through the write path and read back as JSON both by default and via `--format json`.
- `VAL-007`:
  `TestIssueCreateAppendsOneEventAndProjectsOneIssue` proves one logical `fid issue create` write appends exactly one `issue.created` v1 event and projects one `issues` row.
- `VAL-008`:
  `TestIssueCreateAndRebuildRoundTripReturnsIdenticalJSON` proves delete-DB plus `fid rebuild --by developer` preserves identical `fid issue get` JSON output.
- `VAL-016`:
  `TestIssueCreateAndRebuildValidateRoleExistence` proves missing `--by` and unknown roles are rejected for both `fid issue create` and `fid rebuild`, while the positive create/rebuild paths succeed in the other Task 6 tests.

## Tests run

- `go test ./cmd/fid -run TestIssueCreateThenGetReturnsJSONByDefaultAndWithFormatFlag` - FAIL first, then PASS after wiring `issue` into the root command and implementing the minimal service path.
- `go test ./cmd/fid -run 'TestIssueCreateAppendsOneEventAndProjectsOneIssue|TestIssueCreateAndRebuildRoundTripReturnsIdenticalJSON|TestIssueCreateAndRebuildValidateRoleExistence|TestIssueGetRejectsUnsupportedFormats|TestRebuildCommand'` - PASS
- `go test ./cmd/fid ./internal/issues ./internal/policy ./internal/format` - PASS
- `go test ./...` - PASS
- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS (`total: 52.5%`)
- `go build -o /tmp/fid ./cmd/fid` - PASS

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
- `PASS` `go tool cover -func=coverage.out`
- `PASS` `go build -o /tmp/fid ./cmd/fid`
- `Not triggered` `go test -tags integration -race ./...` because the create/get/rebuild round-trip tests were kept in the default test suite and not placed behind the `integration` build tag.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- Role validation now scans journal events directly and deliberately tolerates a corrupt trailing line after seeded roles have already been discovered so `fid rebuild` can still report the intended corruption error. Revisit this fallback when mutable role events ship beyond Phase 1 bootstrap.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit creation.

## Commit subject

- Pending commit creation.
