# Phase 20 - Task 3e: Production initialization benchmark

## Objective

Add a reproducible benchmark for the complete production
`InitializeWorkspace` durability path and record stage-level timing without
changing that path. This creates the measurement baseline needed to distinguish
test-fixture improvements from production durability work.

## Validation Contract Coverage

- No VAL-NNN applies. This is test/performance instrumentation governed by
  DEC-133 and does not change product behavior.

## Binding Constraints

- Exercise real initialization with migrations, journal append, rebuild
  verification, filesystem sync, and activation.
- Use existing `initHooks` only for observation; do not add production timing
  branches or weaken any durability stage.
- No product optimization and no dependency change.

## Scope

### Files To Create Or Modify

- `internal/bootstrap/init_benchmark_test.go`
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/tasks/task-3e-done.md`

### Files To Read For Context

- `internal/bootstrap/init.go`
- `internal/bootstrap/init_test.go`
- `project-docs/reviews/2026-08-10-gate-execution-and-test-runtime-diagnosis.md`
- `project-docs/docs/phases/phase-20/decisions.md` DEC-133

### Out Of Scope

- Optimizing initialization, migrations, rebuild, SQLite, or filesystem sync.
- Reusable command-test fixtures; Task 3h owns them.

## Acceptance Criteria

- [ ] `BenchmarkInitializeWorkspace` repeatedly initializes a fresh root through
      the real production path and reports allocations. (no VAL)
- [ ] Successful activated output is verified and cleanup is outside the timed
      region where appropriate. (no VAL)
- [ ] Stage observations cannot alter control flow or production semantics.
- [ ] Three samples and baseline metrics are recorded in the report and anchor.

## First Failing Test

Explicit measurement-task TDD exception: first prove the benchmark is absent
with `rg -n '^func BenchmarkInitializeWorkspace' internal/bootstrap` returning
`exit 1 (no output)`, then add and execute it. No product behavior changes.

## Documentation Updates

- `infra-improvement-pack.md` - record Task 3e state and benchmark baseline.

## Exact Commands To Run

Run exactly:

```text
go test -run '^$' -bench '^BenchmarkInitializeWorkspace$' -benchmem -count=3 ./internal/bootstrap
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

Required-command rules live in `AGENTS.md` (Commands). Every command is `PASS`
or `BLOCKED`; benchmark/focused commands do not substitute for the full gate.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Preserve and poll long-command sessions until observed termination.
- Completion report: `tasks/task-3e-done.md`.
