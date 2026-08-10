# Phase 20 - Task 3g: Reuse one integration binary

## Objective

Build one immutable `fid` binary for the dedicated integration package and
reuse it across all real-process journeys. Preserve isolated workspaces and
subprocess assertions while removing four redundant builds.

## Validation Contract Coverage

- No VAL-NNN applies. Existing binary journeys retain their product evidence;
  this task optimizes their test harness.

## Binding Constraints

- Package-owned `TestMain` or equivalent builds exactly once per package run.
- Tests may not mutate the binary and must still run real OS processes.
- A build failure fails the package with captured output; no fallback build.
- No dependency or product source change.

## Scope

### Files To Create Or Modify

- `cmd/fid/integration/main_test.go`
- `cmd/fid/integration/*_test.go` helper/call sites
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/tasks/task-3g-done.md`

### Files To Read For Context

- `cmd/fid/integration/` after Task 3f
- Task 3f report and DEC-133

### Out Of Scope

- Scenario classification, workspace cloning, or product build changes.

## Acceptance Criteria

- [ ] A harness regression proves all four journeys receive one shared binary
      path and no per-test build helper remains. (no VAL)
- [ ] The package builds once, reports failures, and cleans its temporary root.
- [ ] Each journey retains a distinct workspace and real subprocess execution.
- [ ] Dedicated integration race gate passes with all four journeys.

## First Failing Test

Add a package-harness inventory assertion requiring one TestMain-owned binary
and no per-test build function; confirm failure while Task 3f's helper remains.

## Documentation Updates

- `infra-improvement-pack.md` - record one-build evidence and timing.

## Exact Commands To Run

Run exactly:

```text
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

- Do not use per-test `sync.Once`; package lifecycle owns the build.
- Completion report: `tasks/task-3g-done.md`.
