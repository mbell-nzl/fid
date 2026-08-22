# Phase 20 - Task 3l: Deterministic maintenance time

## Objective

Remove the remaining maintenance/closure wall-clock sleep by using the
existing `maintenance.CheckOptions.Now` seam or deterministic persisted
fixture timestamps. Preserve the real CLI default and maintenance eligibility
rules.

## Validation Contract Coverage

- No new VAL-NNN applies. Existing maintenance and closure scenario assertions
  retain their coverage; this task only makes their time setup deterministic.

## Binding Constraints

- Reuse `maintenance.CheckOptions.Now`; do not create a second maintenance
  clock abstraction.
- The production CLI's zero-value path remains real UTC time.
- If a real-binary journey needs aged state, construct that state through
  supported journal/projection test helpers and preserve rebuild parity.
- No shortened sleeps, polling, public test-only flag, or environment override.
- Do not change grace periods, comparison boundaries, findings, output, event
  data, or repair behavior.

## Scope

### Files To Create Or Modify

- `internal/maintenance/` focused tests if boundary coverage is missing
- Maintenance command construction/tests under `cmd/fid`
- The dedicated closure integration fixture under `cmd/fid/integration`
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/tasks/task-3l-done.md`

### Files To Read For Context

- `internal/maintenance` check options and timestamp comparisons
- The real-binary closure scenario after Task 3f
- Journal/projector helpers for deterministic occurred-at timestamps
- Task 3k report and DEC-133

### Out Of Scope

- New maintenance behavior or a public clock-control interface.
- Unrelated sleeps in lock-contention or cancellation tests.

## Acceptance Criteria

- [ ] Fixed-time maintenance tests cover just-before, exact-boundary, and
      just-after eligibility with `CheckOptions.Now`. (no VAL)
- [ ] The closure scenario uses deterministic aged state and contains no
      wall-clock sleep.
- [ ] Rebuilding the deterministic scenario yields byte-equivalent projection
      and the same maintenance result.
- [ ] A repository sleep inventory classifies every remaining `time.Sleep` as
      an intentional non-expiry timing test or a follow-up finding.
- [ ] Repeated maintenance and dedicated integration race tests pass.

## First Failing Test

Add or strengthen an exact-boundary maintenance eligibility test using
`CheckOptions.Now`; confirm the closure fixture cannot yet express the same
deterministic aged state without sleeping.

## Documentation Updates

- `infra-improvement-pack.md` - deterministic maintenance setup, remaining
  sleep inventory, and final per-task timing.

## Exact Commands To Run

Run exactly:

```text
go test -tags integration -race -count=10 ./internal/maintenance ./cmd/fid/integration/...
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

- Keep `-tags integration` on every dedicated-package command because those
  journey files retain their build constraint.
- Completion report: `tasks/task-3l-done.md`.
