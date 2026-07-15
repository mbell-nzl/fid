# Phase 17 - Task 5: Wire `fid maintenance check` and All Formats - Completion Report

## Summary of what was implemented

- Added the read-only `fid maintenance check` command family with JSON default plus TOON, Markdown, and YAML output.
- Refreshed the derived projection before auditing, loaded the configured maintenance grace period, and rendered the complete structured result before returning a failing audit outcome.
- Added focused maintenance Markdown rendering and concise operator examples without changing pre-existing Markdown surfaces.

## Files created or modified

- `cmd/fid/maintenance.go`
- `cmd/fid/maintenance_test.go`
- `cmd/fid/maintenance_integration_test.go`
- `cmd/fid/root.go`
- `internal/format/markdown.go`
- `internal/format/markdown_maintenance.go`
- `internal/format/markdown_test.go`
- `project-docs/docs/playbook.md`
- `project-docs/docs/phases/phase-17/tasks/task-5-done.md`

## Documentation updated

- `project-docs/docs/playbook.md` — added maintenance audit usage, read-only freshness behavior, formats, and status/exit semantics.

## Validation Contract Coverage

- VAL-012: PASS — command tests prove `pass`/zero for healthy output and `fail`/non-zero with a listed structural finding in JSON, TOON, Markdown, and YAML.
- VAL-013: PASS — the identical warning fixture returns `warn`/zero by default and `fail`/non-zero with `--warn-as-error`; warning payloads remain byte-for-byte equivalent after JSON decoding.
- VAL-014: PASS — integration tests cover both missing-projection rebuild and stale-projection catch-up, asserting healthy fresh results and byte-identical journals before and after each audit.

## Tests run

- `go test ./cmd/fid -run '^TestMaintenanceCheckCommandRendersHealthyAndFailingResults$' -count=1` - expected RED: `unknown command "maintenance" for "fid"`.
- `go test ./cmd/fid -run '^TestMaintenanceCheckCommandRendersHealthyAndFailingResults$' -count=1` - PASS after command implementation.
- `go test -race ./cmd/fid -run '^TestMaintenanceCheck' -count=1` - PASS.
- `go test -race ./internal/format -run '^(TestWriteMarkdownRendersMaintenanceCheck|TestWriteDispatchesMarkdown)$' -count=1` - PASS.
- Runtime entrypoint: `workspace=$(mktemp -d); cd "$workspace"; /tmp/fid init >/dev/null; /tmp/fid maintenance check --format markdown` - PASS; disposable workspace rendered `# Maintenance Check` with `**Status:** pass`, no findings, and zero counts.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS (`0 issues.`)
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS (total statement coverage: 76.4%)
- `go test -tags integration -race ./...` - PASS
- `go build -o /tmp/fid ./cmd/fid` - PASS
- Conditional commands - N/A (the task specification lists none).

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
