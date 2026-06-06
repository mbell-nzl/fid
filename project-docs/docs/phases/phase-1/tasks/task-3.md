# Phase 1 - Task 3: Write Lock And Journal Writer (Append + Flush)

## Objective

Implement the checkout-scoped write lock (`.fid/write.lock`) and the journal
writer: under the lock, read the UTC clock once, choose the journal file,
allocate the next per-day event sequence, append one flushed JSONL line as the
commit point, and detect clock skew. This is the write half of the persistence
loop; projection is wired in Task 4.

## Requirements

- PRD refs: PRD-OPS-001, PRD-OPS-002, PRD-NFR-008
- TRD refs: TRD §4.3 (write path steps 1-5, event-ID/sequence allocation,
  clock-skew rule, concurrency/lock metadata), TRD-ARCH-004, TRD-PERF-006
- Design refs: Design Plan CLI §18 (clock-skew error text)
- Quality refs: Quality Plan §5.1

## Validation Contract Coverage

- VAL-007: One logical write appends exactly one v1 JSONL line (append side).
- VAL-013: `.fid/write.lock` serializes same-checkout writes with owner metadata
  and rejects a waiting caller with a holder-named error.
- VAL-015: Event-ID sequence allocation per UTC journal file under the lock.

## Scope

### Files To Create Or Modify

- `internal/journal/lock.go` (lock-dir acquire/release + owner metadata)
- `internal/journal/writer.go` (single clock read, file selection, sequence
  allocation, append + flush, clock-skew detection)
- `internal/journal/*_test.go`
- `project-docs/docs/phases/phase-1/tasks/task-3-completion.md`

### Files To Read For Context

- `project-docs/specs/02_TRD.md` §4.3 (write path, event-ID and journal-file
  allocation, concurrency, lock metadata)
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` §18 (clock-skew error)
- `internal/journal/event.go`, `canonical.go`, `id.go` (Task 2 outputs)
- `project-docs/docs/phases/phase-1/qa/validation-contract.md`

### Out Of Scope

- Projecting effects into SQLite, `projection_meta`, catch-up, rebuild,
  corruption/precondition handling on replay (all Task 4).
- The validate/build domain logic of specific commands (`fid init`, issue
  create) — Tasks 5/6. Provide a reusable append API those tasks call.

## Acceptance Criteria

- [ ] Acquiring `.fid/write.lock` creates a lock directory containing PID,
  hostname, started_at, command, and actor role (when available). (VAL-013)
- [ ] A second acquire while the lock is held waits the configured short timeout,
  then exits non-zero with an error naming the holder and lock path; the lock is
  never silently broken while live. (VAL-013)
- [ ] Two concurrent writes (goroutines) both append distinct events with no
  lost line and monotonically increasing sequences. (VAL-007, VAL-013)
- [ ] Appending one event reads the UTC clock once and produces exactly one
  flushed JSONL line in `.fid/journal/<UTC date>.jsonl`. (VAL-007)
- [ ] Event ID sequence is allocated per UTC journal file: first event of a day
  is `EVT-<date>-000001`, next is `000002`, scanning existing lines under the
  lock. (VAL-015)
- [ ] A chosen UTC date earlier than the greatest existing journal date exits
  non-zero with a clock-skew error and appends nothing. (VAL-007)

## TDD Plan

- First failing test: append one event, assert exactly one line written and the
  parsed event ID sequence is `000001`.
- Additional tests: second append increments sequence; lock contention yields a
  holder-named error; two concurrent writers both persist; clock-skew date is
  rejected with no append; lock metadata fields present.
- Regression risks: partial writes on crash (flush ordering — append is commit
  point); lock-dir cleanup on error paths; sequence races without the lock.

## Documentation Updates

- `project-docs/docs/phases/phase-1/tasks/task-3-completion.md` - completion
  report citing VAL-007/013/015 evidence and gate results.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: None.
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. The `-race` flag is
mandatory because this task adds concurrency. Every required command must be
recorded in the completion report as `PASS`, or as `BLOCKED` with reason in a
blocker report.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` if the concurrency/lock tests are
  placed behind the `integration` tag.

## Notes For Developer

- The journal append is the commit point: if append/flush fails, no DB state is
  touched (TRD §4.3). Keep projection out of this task.
- Acquire the lock by atomically creating a lock directory (mkdir is atomic);
  do not implement timeout-based force-break of a live lock in v0.1.
- Expose a small append API (`AppendEvent` under an already-held lock, plus a
  `WithLock` helper) that Tasks 4-6 reuse, so the full write path composes
  validate → build → append → project under one lock.
