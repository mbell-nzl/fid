# Phase 20 - Task 3i: Explicit command working directory

## Objective

Remove process-wide `os.Chdir` from the `cmd/fid` test harness by giving the
root command an explicit working directory. Preserve the production CLI's
current behavior by capturing the caller's real working directory when no test
directory is supplied.

## Validation Contract Coverage

- No new VAL-NNN applies. Existing command behavior and workspace-selection
  assertions remain authoritative; this task changes dependency delivery.

## Binding Constraints

- The production constructor captures `os.Getwd` once and returns its error.
- Command handlers resolve implicit and relative paths from the injected
  directory; explicit absolute paths retain their meaning.
- Test helpers supply an isolated directory without changing process state.
- No public flag, output, schema, journal event, or workspace-selection rule
  changes.
- Do not add parallel tests here; Task 3j owns concurrency after this seam is
  proven race-safe.

## Scope

### Files To Create Or Modify

- `cmd/fid/root.go` and command files that currently call `os.Getwd`
- `cmd/fid/root_test.go`, `cmd/fid/issue_test.go`, `cmd/fid/init_test.go`
- Selected `cmd/fid/*_test.go` harness call sites
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/tasks/task-3i-done.md`

### Files To Read For Context

- All `cmd/fid` production and test references to `os.Getwd` or `os.Chdir`
- Task 3h report and DEC-133

### Out Of Scope

- Public working-directory flags or environment variables.
- Parallel test execution; Task 3j owns it.

## Acceptance Criteria

- [ ] `TestRootCommandUsesInjectedWorkingDirectoryWithoutChangingProcess`
      proves command execution uses the supplied root while the process working
      directory remains unchanged. (no VAL)
- [ ] The production constructor preserves implicit and relative workspace
      selection from the caller's actual working directory.
- [ ] The shared command-test harness contains no `os.Chdir`.
- [ ] An inventory records any remaining `os.Chdir` and explains why it is not
      in a command-test execution path.
- [ ] Repeated `cmd/fid` race tests pass before Task 3j adds parallelism.

## First Failing Test

Add `TestRootCommandUsesInjectedWorkingDirectoryWithoutChangingProcess`; confirm
it fails because the current harness can select a workspace only by calling
`os.Chdir`.

## Documentation Updates

- `infra-improvement-pack.md` - working-directory seam and final inventory.

## Exact Commands To Run

Run exactly:

```text
go test -race -count=5 ./cmd/fid
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go build -o /tmp/fid ./cmd/fid
```

Required-command rules live in `AGENTS.md` (Commands).

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Keep the directory dependency package-private and narrow; do not place test
  state in package globals.
- Completion report: `tasks/task-3i-done.md`.
