# Phase 20 - Task 3i: Explicit command working directory - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Changed the production root constructor to capture `os.Getwd` once, return a
  wrapped error, and build the command with an explicit directory in root
  context.
- Added `newRootCommandAt` for isolated tests and a narrow
  `commandWorkingDirectory` accessor used by every handler.
- Preserved blank and relative `--project` behavior by resolving relative paths
  from the captured/injected directory before absolute normalization.
- Removed every command-test `os.Chdir` and `mustChdir`; command execution no
  longer mutates process working directory.
- Added focused tests for injected-directory execution, unchanged process cwd,
  production cwd capture, and relative-project resolution.

## Files created or modified

- `cmd/fid/root.go`
- 24 command handler files formerly calling `os.Getwd`
- `cmd/fid/root_flags_test.go`
- `cmd/fid/init_test.go`
- `cmd/fid/issue_test.go`
- `cmd/fid/rebuild_test.go`
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/tasks/task-3i-done.md`

## Documentation updated

- `infra-improvement-pack.md` records the final `Getwd`/`Chdir` inventories and
  the observed five-count race-session completion.

## Validation Contract Coverage

- No VAL-NNN applies. Existing command tests retain product assertions.
  `TestRootCommandUsesInjectedWorkingDirectoryWithoutChangingProcess`,
  `TestRootCommandCapturesCurrentWorkingDirectory`, and
  `TestRootCommandResolvesRelativeProjectFromInjectedWorkingDirectory` prove
  the dependency seam and preserved defaults.

## Tests run

- `cd /tmp/fid-task-3i-entrypoint.mCi1vD && /tmp/fid init --json --prefix T3I` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"project.init","result":{"journal_path":".fid/journal/2026-08-16.jsonl","project":{"id":"project-983138bf390bb71653eb2809cad501b1","issue_prefix":"T3I","root":"/tmp/fid-task-3i-entrypoint.mCi1vD"}},"event":{"id":"EVT-20260816-000001","occurred_at":"2026-08-16T10:49:09.532299051Z"}}` (terminated `exit_code: 0`)
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-3i-done.md project-docs/docs/phases/phase-20/tasks/task-3i.md` - FAIL - `FAIL report-lint found 27 failure(s)` (terminated `exit_code: 27`)

## Historical focused test evidence

- `go test -run '^TestRootCommandUsesInjectedWorkingDirectoryWithoutChangingProcess$' ./cmd/fid` - EXPECTED RED - `FAIL`.
- `go test -run '^TestRootCommandUsesInjectedWorkingDirectoryWithoutChangingProcess$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.041s`.
- `go test -run '^TestRootCommand(CapturesCurrentWorkingDirectory|ResolvesRelativeProjectFromInjectedWorkingDirectory|UsesInjectedWorkingDirectoryWithoutChangingProcess)$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.073s`.
- `rg -n 'os\.Chdir|mustChdir' cmd/fid --glob '*.go'` - PASS no-match - `exit 1 (no output)`.

## Exact quality-gate results

- `go test -race -count=5 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	515.474s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					78.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.537s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Historical quality-gate evidence

- `go test -race -count=5 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	431.955s`.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`.
- `golangci-lint run ./...` - PASS - `0 issues.`.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.5% of statements`.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					77.7%`.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.718s`.
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)`.

## Deviations from spec

- The initial canonical report-lint attempt failed with `FAIL report-lint found 27 failure(s)` (terminated `exit_code: 27`) because the completed report lacked the typed workflow outcome and observed termination syntax. This report-only repair adds those canonical fields and fresh gate evidence; the historical focused and quality-gate records remain intact below their original content. No product behavior, command, assertion, threshold, or dependency changed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `3e830b49c8a90472f091e7b852848cf6f86694d8`

## Commit subject

- `refactor(cli): inject command working directory`
