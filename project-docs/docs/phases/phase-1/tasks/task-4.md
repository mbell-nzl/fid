# Phase 1 - Task 4: Replayer, Projection, projection_meta, Catch-Up, And `fid rebuild`

## Objective

Implement the read/projection half of the loop: replay journal files in order
through the upgrader chain, apply each event's effects to SQLite in one
transaction, maintain the `projection_meta` checkpoint, decide clean-read vs
catch-up vs full rebuild, and surface corruption and precondition conflicts as
structured errors. Ship the `fid rebuild` command on top.

## Requirements

- PRD refs: PRD-OPS-003, PRD-OPS-004, PRD-OPS-012
- TRD refs: TRD §4.3 (write path step 6, rebuild, catch-up, malformed journal,
  preconditions, canonical-hash checkpoint), TRD-ARCH-001, TRD-ARCH-002,
  TRD-ARCH-003
- Design refs: Design Plan CLI §18 (`fid rebuild` output and error text)
- Quality refs: Quality Plan §5.1

## Validation Contract Coverage

- VAL-008: Replay is pure data; delete-DB + rebuild yields identical state.
- VAL-009: Replay routes events through the upgrader chain (golden v1 replay).
- VAL-010: Malformed trailing line reported as corruption, not skipped.
- VAL-011: Precondition mismatch stops replay with a structured conflict error.
- VAL-012: `projection_meta` drives clean-read vs catch-up vs rebuild.

## Scope

### Files To Create Or Modify

- `internal/journal/replay.go` (ordered replay, upgrader routing, corruption +
  precondition handling)
- `internal/journal/project.go` (apply effects in a transaction, update
  `projection_meta`)
- `internal/db/projection_meta.go` (checkpoint read/write + freshness decision)
- `internal/journal/*_test.go`, `internal/db/*_test.go`
- `cmd/fid/rebuild.go` (thin `fid rebuild` adapter)
- `testdata/journal/` corrupt-tail and conflicting-precondition fixtures
- `project-docs/docs/phases/phase-1/tasks/task-4-completion.md`

### Files To Read For Context

- `project-docs/specs/02_TRD.md` §4.3 (rebuild, catch-up, malformed journal,
  preconditions, projection_meta semantics)
- `project-docs/specs/03_TRD_SCHEMA.md` §9.5 (`projection_meta`)
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` §18
- Task 2/3 outputs (event types, canonical hash, writer/lock APIs)
- `project-docs/docs/phases/phase-1/qa/validation-contract.md`

### Out Of Scope

- `fid init` bootstrap and issue create/get (Tasks 5/6). This task provides the
  reusable replay/projection engine and the operator `fid rebuild` command.
- Performance/benchmark tuning of rebuild (Phase 3).
- Tables beyond `projection_meta`/`roles`/`issues` (Phase 2 schema widening).

## Acceptance Criteria

- [ ] Replaying a journal applies effects to SQLite in a single transaction per
  event and writes a `projection_meta` checkpoint (journal_file, offset,
  last_event_id, last_event_hash, schema versions, replay_status). (VAL-008, VAL-012)
- [ ] After creating state, deleting `fid.sqlite`, and rebuilding, queries are
  identical and replay neither re-runs command logic nor reads DB state to
  resolve effect values. (VAL-008)
- [ ] Replay routes every event through the upgrader chain; a golden v1 journal
  replays to the expected DB state. (VAL-009)
- [ ] A malformed trailing JSONL line stops replay with a non-zero error naming
  the journal file and line; a clean journal replays without error. (VAL-010)
- [ ] An effect precondition mismatch stops replay with a structured error
  naming event_id, event_type, table, pk, expected, actual, and file/line; a
  matching precondition applies the effect. (VAL-011)
- [ ] Freshness decision: a `projection_meta` matching the journal serves reads
  without rebuild; missing meta, schema-version mismatch, position mismatch, or
  last-event-hash mismatch triggers rebuild/catch-up before serving. (VAL-012)
- [ ] `fid rebuild --by <role>` replays from scratch into a fresh DB under the
  write lock and reports replayed event/file counts. (VAL-008, VAL-012)

## TDD Plan

- First failing test: project a single-effect event and assert the row plus the
  `projection_meta` checkpoint.
- Additional tests: delete-DB + rebuild identical-state; golden v1 replay;
  corrupt-tail error names file/line; precondition mismatch structured error +
  matching-precondition control; freshness matrix (clean match, missing meta,
  hash mismatch, position mismatch); `fid rebuild` happy path and corrupt-journal
  error.
- Regression risks: file ordering (ascending UTC date then line order);
  transaction atomicity per event; hash recomputation matching the writer's
  canonical form; serving stale reads on mismatch.

## Documentation Updates

- `project-docs/docs/phases/phase-1/tasks/task-4-completion.md` - completion
  report citing VAL-008/009/010/011/012 evidence and gate results.

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

The full task gate above is required before task close. Every required command
must be recorded in the completion report as `PASS`, or as `BLOCKED` with reason
in a blocker report.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go build -o /tmp/fid ./cmd/fid` is required because this task adds the first
  real subcommand (`fid rebuild`); confirm the binary builds.
- `go test -tags integration -race ./...` if replay/rebuild tests use the
  `integration` tag with temp DBs.

## Notes For Developer

- Reuse the canonical-JSON hash from Task 2 for `last_event_hash`; the writer
  and replayer must agree byte-for-byte.
- Any path mutating `fid.sqlite` (rebuild, catch-up) acquires `.fid/write.lock`
  (Task 3). Clean reads matching the checkpoint may serve without the lock.
- Do not silently skip, truncate, or reorder journal entries — corruption and
  conflicts are reported, never repaired automatically (TRD §4.3).
- Keep `cmd/fid/rebuild.go` a thin adapter; the replay/projection logic lives in
  `internal/` (TRD-ARCH-005, VAL-017).
