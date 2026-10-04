# Phase 24 - Task 1: Idempotency field-digest storage (migration 0019)

## Objective

Add durable storage for per-field request digests (DEC-275). This covers
migration 0019, the `idempotency_results` column, the value carried by the
idempotency Effect so rebuild reproduces it, fail-closed replay validation,
and projection schema 19 in the export manifest. Call sites do not compute
real maps yet: task 2 adopts them. Until then, every commit writes the empty
map `{}`.

## Validation Contract Coverage

- VAL-001: 0019 up adds exactly the one column.
- VAL-058: 0019 down restores schema 18 exactly.
- VAL-002 (storage half): rebuild reproduces `request_field_digests_json`
  byte-for-byte. Task 2 extends the evidence to every family.
- VAL-003: a keyed Event whose idempotency Effect lacks a valid map fails
  closed.
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
- **Where the map travels (DEC-278).** It is an insert value of every
  `idempotency_results` insert Effect: the central builder in
  `internal/journal/idempotency.go` and the direct `gate.execution-started`
  insert in `internal/evidence/execution.go`. On the `gate.receipt-recorded`
  update it is a **guard** beside `request_digest`, never a value, so the
  column stays immutable after insert. Current-event validation requires the
  map as an object of `digest` values on every insert, and as the matching
  guard on the receipt update. When it is absent, not an object, or holds a
  non-digest value, replay stops with `journal-replay-stopped` (VAL-003).
- **Closed registries the column enters** (inventory R3, re-measured at
  `eb34889`):
  - `internal/journal/current_event_columns.go`: the `idempotency_results`
    column and type list;
  - `internal/journal/current_event_json.go`: register
    `idempotency_results.request_field_digests_json` as a JSON kind (today
    an unlisted `_json` column fails `unregistered current JSON column`);
  - `internal/journal/current_event_updates.go`: the receipt-update guards
    become `request_digest request_field_digests_json`. The update values
    are unchanged.

  `internal/export/archive_updates.go` (`semanticUpdateColumns`) stays
  unchanged, because the map is never an update value. Confirm this in the
  report.
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
- `internal/journal/current_event_columns.go`, `current_event_json.go`,
  `current_event_updates.go` (closed registries above)
- `internal/evidence/execution.go`: the direct `gate.execution-started`
  insert Effect and the `gate.receipt-recorded` update guard. Also every other
  direct constructor found by
  `rg -n "EventIdempotency\{|Table: \"idempotency_results\"" internal -g '!*_test.go'`
- Retained tests displaced by the new NOT NULL column or pinning schema 18 or
  manifest 18 (convergence, with provenance), including:
  `internal/maintenance/check_test.go` (`seedTask21CorruptProjection`),
  `internal/db/operational_kernel_migration_test.go` (direct
  `idempotency_results` inserts), and the closed portable column list in
  `cmd/fid/integration/p22_portable_schema_test.go`. The tagged gate below
  executes the last one.

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` section 7
- `project-docs/docs/phases/phase-24/input-inventory.md` R3, R4
- `project-docs/docs/phases/phase-24/decisions.md` DEC-275, DEC-277
- `internal/export/archive_updates.go` (stays unchanged; see above)
- `migrations/0013_p21_operational_kernel.up.sql` (table DDL);
  `migrations/0018_p23_remove_v1_judgment.*` (style); `internal/db/migration0018_test.go`
- `project-docs/docs/phases/phase-23/decisions.md` DEC-264, DEC-270

### Out Of Scope

- Computing real field digests at call sites, and `differing_fields` in
  failure details (task 2).
- Any input shape, observation or Human output change.

## Acceptance Criteria

- [ ] The migration test shows the schema-18 dump and up adding exactly the
  one column with the accepted DDL. (VAL-001)
- [ ] The same test shows down byte-equal to the schema-18 dump. (VAL-058)
- [ ] A journal with keyed mutations from at least the issue, claim, run and
  review families rebuilds with `request_field_digests_json` byte-equal to
  the live projection. (VAL-002)
- [ ] The journal fixture with the map replays. Each fixture with the map
  removed, as a non-object, or holding a non-digest value stops with
  `journal-replay-stopped` (exit 6) at that line. A receipt update whose
  guard map differs from the stored map also fails closed. (VAL-003)
- [ ] Export, validate and import of a schema-19 archive succeed. A copy whose
  manifest declares 18 fails with `archive-invalid`, check
  `schema-unsupported` on `manifest.json`, and the target tree digest is
  unchanged. (VAL-004)
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
reports for `cmd/fid/integration`, `cmd/fid`, `internal/workflow` and the
slowest other ordinary package. The budgets are tagged integration under
900 s and no ordinary package over 600 s (VAL-052); report any ordinary
package over the 480 s warning. A breach caused by this task's own new journeys is a bounded speed
repair owned by this task; otherwise stop and escalate. Never record a
timed-out gate as `PASS`.

Required-command rules live in `AGENTS.md` (Commands). Every required command
is reported `PASS` or `BLOCKED` with reason, with its own output token and
observed termination. Focused red/green selections never replace these
commands. Report lint is a Runner control, not a Developer gate.

List any conditional command required by this task and why:

- The benchmark line is required because the projection and every keyed
  Effect change (VAL-005). The QA plan measured about 19.6% headroom (worst
  fresh sample 24128 ms). The benchmark journal contains no keyed Events, and
  its generator is a protected definition (DEC-158): do not change it.
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
- `ADD COLUMN` with the accepted NOT NULL text fails on a populated table
  (`Cannot add a NOT NULL column with default value NULL`, QA plan 6.4).
  Init, rebuild and import migrate empty tables, but catch-up migrates the
  live projection in place. Record in the report what an existing Phase 23
  workspace with keyed rows does under the new binary: the catch-up failure,
  and whether `fid rebuild` then fails closed with `journal-replay-stopped`.
  This is a note, not a criterion. No compatibility path is added (DEC-270).
