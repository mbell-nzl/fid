# Phase 24 - Task 1: Idempotency field-digest storage (migration 0019)

## Objective

Add durable storage for per-field request digests (DEC-275). This covers
migration 0019, the `idempotency_results` column, the value carried by the
idempotency Effect so rebuild reproduces it, fail-closed replay validation,
and projection schema 19 in the export manifest. Call sites do not compute
real maps yet: task 2 adopts them. Until then, every commit writes the empty
map `{}`.

## Validation Contract Coverage

- VAL-001: 0019 adds exactly the one column; down restores schema 18 exactly.
- VAL-002 (storage half): rebuild reproduces `request_field_digests_json`
  byte-for-byte. Task 2 extends the evidence to every family.
- VAL-003: a keyed Event whose idempotency Effect lacks the map fails closed.
- VAL-004: projection schema and manifest `projection_schema` 19; import
  accepts only 19.
- VAL-005: the 100k rebuild stays under 30000 ms.
- VAL-052: gate wall-time budgets.

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Contract section 7 and DEC-275. The accepted DDL is exactly
  `request_field_digests_json TEXT NOT NULL CHECK (json_valid(request_field_digests_json))`
  on `idempotency_results`. Shipped SQL must equal it (DEC-178). The value
  is a canonical JSON object from public input name to `digest`. In this task
  it is always `{}`.
- The down migration restores the schema-18 `idempotency_results` exactly
  (DEC-264). Use the established table-rebuild pattern if SQLite `DROP
  COLUMN` cannot restore byte-equal DDL.
- The map travels as a value of the idempotency Effect built centrally in
  `internal/journal/idempotency.go`, on both the insert path and the
  gate-receipt `update` path. Every other path that constructs an
  idempotency Effect or `EventIdempotency` directly (inventory R3) carries
  it too. Current-event validation requires the key and an object of
  `digest` values. When it is absent, replay stops with
  `journal-replay-stopped` (VAL-003).
- `IdempotencyRequest` gains the field-digest map. `IdempotencyConflictError`
  gains the stored and supplied maps for task 2's use. The public failure
  details do not change in this task.
- `fiddb.ActiveMigrationVersion` becomes 19, and so does the export
  manifest. `ProjectionSchemaVersion` (projection meta format, currently 5)
  is a different constant and does not change.
- Journals from earlier binaries are unsupported (DEC-270). Do not add a
  compatibility path or criterion.

## Scope

### Files To Create Or Modify

- `migrations/0019_p24_idempotency_fields.up.sql`, `.down.sql` (new)
- `internal/db/projection_meta.go`; `internal/db/migration0019_test.go` (new)
- `internal/journal/idempotency.go`, `current_event_idempotency.go`,
  `current_event_validation.go` and the projection applier, if it needs the
  new column
- Every direct idempotency-Effect or `EventIdempotency` constructor found by
  `rg -n "EventIdempotency\{" internal -g '!*_test.go'`
- `internal/journal/rebuild_bench_test.go` and its journal generator, if the
  generator builds keyed Events
- Retained tests pinning schema 18 or the manifest value 18 (convergence,
  with provenance)

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` section 7
- `project-docs/docs/phases/phase-24/input-inventory.md` R3, R4
- `project-docs/docs/phases/phase-24/decisions.md` DEC-275, DEC-277
- `migrations/0013_p21_operational_kernel.up.sql` (table DDL);
  `migrations/0018_p23_remove_v1_judgment.*` (style); `internal/db/migration0018_test.go`
- `project-docs/docs/phases/phase-23/decisions.md` DEC-264, DEC-270

### Out Of Scope

- Computing real field digests at call sites, and `differing_fields` in
  failure details (task 2).
- Any input shape, observation or Human output change.

## Acceptance Criteria

- [ ] The migration test shows the schema-18 dump, up adding exactly the one
  column, and down byte-equal to 18. (VAL-001)
- [ ] A journal with keyed mutations from at least the issue, claim, run and
  review families rebuilds with `request_field_digests_json` byte-equal to
  the live projection. (VAL-002)
- [ ] The paired journal fixture replays with the map, and stops with
  `journal-replay-stopped` at the map-less Event. (VAL-003)
- [ ] Export, validate and import of a schema-19 archive succeed. An archive
  whose manifest declares 18 is rejected with nothing published. (VAL-004)
- [ ] The 100k rebuild benchmark reports `< 30000 ms/rebuild` with its parity
  self-assert. (VAL-005)
- [ ] Wall times recorded and inside budget. (VAL-052)

## First Failing Test

`internal/db/migration0019_test.go`: after migrating to 19, `idempotency_results` has `request_field_digests_json`. Today it does not.

## Documentation Updates

- `project-docs/specs/03_TRD_SCHEMA.md`: the `idempotency_results` column
  and the schema-19 baseline.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root the Runner records. No command is a placeholder.
Apply gofmt and goimports to edited Go files before the final gate. Every
command below is required before task close. The tagged integration gate is
mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/
rg -n "idempotency_results" migrations internal -g '!*_test.go'
```

The `-timeout 30m` bounds follow DEC-227/259. Record the wall time `go test`
reports for `cmd/fid/integration`, `cmd/fid` and `internal/workflow`. The
budgets are tagged integration under 900 s and no ordinary package over 600 s
(VAL-052). A breach caused by this task's own new journeys is a bounded speed
repair owned by this task; otherwise stop and escalate. Never record a
timed-out gate as `PASS`.

Required-command rules live in `AGENTS.md` (Commands). Every required command
is reported `PASS` or `BLOCKED` with reason, with its own output token and
observed termination. Focused red/green selections never replace these
commands. Report lint is a Runner control, not a Developer gate.

List any conditional command required by this task and why:

- The benchmark line is required because the projection and every keyed
  Effect change (VAL-005; 9-11% headroom at Phase 23 end).
- The `rg` line regenerates registry R3 for the report.

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-24/tasks/task-1-done.md`.
  Commit with `Refs: phase-24/task-1`.
- Red first: record the first failing test before changing product code.
- Retained-test convergence follows DEC-162/165/166. Record each displaced
  assertion with its introducing commit, phase and task. Never infer
  ownership from a filename or task label. Never delete or weaken an
  assertion to pass a gate.
- Any missing registry member, contract ambiguity, or need to change an
  accepted contract shape or assertion is a `BLOCKED` escalation candidate,
  not a local decision (DEC-195).
- If down-migration exactness needs a decision beyond this spec's DDL, stop:
  DDL is accepted text.
