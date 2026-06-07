## Summary of what was implemented

- Updated `internal/db.Open` to build a `modernc.org/sqlite` file DSN with `_pragma=foreign_keys(1)` and `_pragma=recursive_triggers(1)` so both required pragmas are applied on every physical SQLite connection the pool opens.
- Kept the `db.Open` public signature unchanged and retained the post-open `Ping`, with no schema, journal, replay, init, or issue behavior changes.
- Added a multi-connection regression test that forces multiple concurrently held `*sql.Conn` instances, proves more than one distinct physical SQLite connection was opened, and asserts both pragmas are `ON` on each connection path.
- This implements the DEC-006 remediation for ESC-001 and closes the VAL-014 gap identified in QA Review 1.

## Files created or modified

- `internal/db/db.go`
- `internal/db/db_test.go`
- `project-docs/docs/phases/phase-1/tasks/task-7-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-1/tasks/task-7-completion.md`

## Validation Contract Coverage

- `VAL-014`:
  `TestOpenEnablesRequiredPragmasOnEveryPhysicalConnection` now forces three concurrently held pooled connections, verifies at least two distinct physical SQLite driver connections were opened, and asserts both `PRAGMA foreign_keys` and `PRAGMA recursive_triggers` report `1` on every held connection.
- Red evidence for `VAL-014` before the fix:
  `go test -run TestOpenEnablesRequiredPragmasOnEveryPhysicalConnection ./internal/db`
  failed with `pragma "foreign_keys" = 0, want 1`, proving the prior single-`Exec` implementation did not satisfy every connection path.
- Green evidence for `VAL-014` after the fix:
  the same focused test passed, and the full repo gate passed afterward.
- Decision and escalation trace:
  `DEC-006` selected DSN `_pragma` enforcement over pool pinning or connection hooks, and `ESC-001` required proof across multiple physical connections before Phase 1 sign-off.

## Tests run

- `go test -run TestOpenEnablesRequiredPragmasOnEveryPhysicalConnection ./internal/db`
  Red: failed first with `pragma "foreign_keys" = 0, want 1`.
  Green: passed after moving pragma enforcement into the SQLite DSN.
- `go test ./internal/db` - PASS
- `gofmt -w internal/db/db.go internal/db/db_test.go` - PASS
- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS (`total: 52.7%`)

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
- `PASS` `go tool cover -func=coverage.out`
- `Not triggered` `go test -tags integration -race ./...` because the multi-connection regression stayed in the default test suite and was not placed behind the `integration` build tag.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- The new DSN helper normalizes paths through an absolute `file:` URI so `_pragma` parameters can be attached safely and consistently. If future scope needs `:memory:` or named in-memory SQLite DSNs through `db.Open`, that case should be specified explicitly rather than inferred from file-path behavior.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Implementation commit: `232a7f3`

## Commit subject

- Implementation subject: `fix(db): enforce sqlite pragmas per connection`
