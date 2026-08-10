# Phase 20 - Task 3j: Bounded command-test parallelism

## Objective

Reduce `cmd/fid` wall time by running proven-isolated command tests in
parallel, after Task 3i removes process-wide directory mutation. Keep tests
that use environment, package globals, shared fixtures, or other process-wide
state serial.

## Validation Contract Coverage

- No new VAL-NNN applies. Existing assertions remain unchanged; timing and
  race-stability evidence govern this harness optimization.

## Binding Constraints

- Add `t.Parallel` only to top-level tests with unique cloned workspaces and no
  process-global mutation.
- Use the Go test scheduler's bounded parallelism; do not add a custom worker
  pool or unbounded goroutines.
- Tests using `t.Setenv`, package globals, shared mutable fixtures, init paths,
  or explicit serialization remain serial.
- Race stability and correctness are hard gates; timing improvement is a
  measured target, never a reason to retry for a green result.

## Scope

### Files To Create Or Modify

- Selected isolated `cmd/fid/*_test.go` files
- `cmd/fid/workspace_fixture_test.go` if concurrency assertions require it
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/tasks/task-3j-done.md`

### Files To Read For Context

- Task 3i report and its remaining global-state inventory
- All `cmd/fid` uses of `t.Setenv`, package-level mutable state, and fixture
  helpers

### Out Of Scope

- Parallelizing integration subprocess journeys.
- Changing product concurrency or persistence behavior.

## Acceptance Criteria

- [ ] A committed isolation inventory classifies selected parallel and retained
      serial tests with reasons. (no VAL)
- [ ] Selected tests use unique workspaces and cannot mutate the pristine
      template or a sibling clone.
- [ ] `go test -race -count=5 ./cmd/fid` passes in one observed run.
- [ ] Before/after ordinary-package wall-time samples are recorded without
      mining for a favorable result.
- [ ] No test assertion, timeout, or race/coverage command is weakened.

## First Failing Test

Add a concurrency isolation regression that launches two cloned-workspace
command executions concurrently and proves their journal/projection/output
state cannot cross; confirm failure before the harness supports concurrent
execution safely.

## Documentation Updates

- `infra-improvement-pack.md` - isolation inventory, repeated-race evidence,
  and timing.

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

- Record the default `-parallel` value and host CPU count with the timing
  evidence; do not bake either value into product or tests.
- Completion report: `tasks/task-3j-done.md`.
