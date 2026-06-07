# Phase 2 - Task 4: Polymorphic Ref, Link Scope, And Claim Scope Backstop Triggers + Phase-End Determinism (Migration 0005)

## Objective

Add the remaining v0.1 backstop triggers — polymorphic ref integrity on `links`,
`claims`, and `activity_log`; link endpoint scope and relation/endpoint
compatibility; and claim `ref_type` scope (migration `0005`). Then prove the
phase-level invariants: `fid rebuild` over a journal that exercises every new
v0.1 table is deterministic and non-destructive with all triggers active, and the
adapter boundary still holds.

## Requirements

- PRD refs: PRD-OPS-003, PRD-OPS-004 (deterministic rebuild)
- TRD refs: TRD §4.5; TRD SCHEMA §14.2 (polymorphic ref), §14.3 + §10.1 (link
  scope/relation compatibility), §14.4 (claim scope); TRD-ARCH-001 (journal as
  source of truth), TRD-ARCH-005 (adapter boundary), TRD-ARCH-006, TRD-SEC-004
- Design refs: N/A
- Quality refs: Quality Plan §5.1; TRD-SEC-004 verification

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-2/qa/validation-contract.md`.

- VAL-013: polymorphic ref integrity on `links`, `claims`, `activity_log`.
- VAL-014: link endpoint scope + relation/endpoint compatibility (§10.1).
- VAL-015: claim `ref_type` scope (v0.1: `issue` only).
- VAL-016: deterministic, non-destructive rebuild over the full v0.1 schema with
  all triggers active (phase-end).
- VAL-017: adapter boundary preserved (phase-end).

## Scope

### Files To Create Or Modify

- `migrations/0005_v01_ref_triggers.up.sql` (polymorphic ref integrity on
  `links` (both endpoints), `claims`, `activity_log`; link scope + relation
  compatibility; claim scope)
- `migrations/0005_v01_ref_triggers.down.sql` (drop all triggers added)
- `internal/db/*_test.go` (raw-SQL backstop tests, both states per rule)
- `internal/journal/*_test.go` (full-schema deterministic-rebuild test)
- `internal/db/projection_meta.go` (bump `ProjectionSchemaVersion` if warranted;
  see Notes)
- `project-docs/docs/phases/phase-2/tasks/task-4-completion.md`

### Files To Read For Context

- `project-docs/specs/03_TRD_SCHEMA.md` §14.2, §14.3, §14.4, §10, §10.1, §11
- `project-docs/specs/02_TRD.md` §4.5, TRD-ARCH-001, TRD-ARCH-005
- `migrations/0002_v01_tables.up.sql` (`links`, `claims`, `activity_log`,
  `requirements`)
- `internal/journal/replay.go`, `internal/journal/project.go`
- `project-docs/docs/phases/phase-2/qa/validation-contract.md`

### Out Of Scope

- `fid link add`/`fid issue claim` commands and auto-inverse link projection
  (v0.1d).
- §14.6 activity_log "inserts only during projection" guard and the `activity`
  effect op (deferred — see contract Scope Deferrals).
- `decision`/`escalation` link endpoints and `qa_review` claim scope (v0.2/v0.3).

## Acceptance Criteria

- [ ] BEFORE INSERT triggers reject `links` whose `from`/`to` endpoint does not
  resolve, `claims` whose `(ref_type, ref_id)` does not resolve, and
  `activity_log` whose `(ref_type, ref_id)` does not resolve; rows whose
  endpoints resolve insert successfully. (VAL-013)
- [ ] `links` trigger accepts only v0.1 endpoint types (`issue`, `requirement`)
  and relation/endpoint combinations per §10.1: rejects `decision`/`escalation`
  endpoints and an incompatible combo (e.g. `implements` issue→issue); accepts
  `implements` issue→requirement, `blocks` issue→issue, and `relates_to`
  issue↔requirement. (VAL-014)
- [ ] `claims` trigger accepts `ref_type='issue'` and rejects any other
  `ref_type` (e.g. `qa_review`). (VAL-015)
- [ ] A journal exercising every new v0.1 table replays cleanly with all triggers
  active; delete-DB + `fid rebuild` yields identical projection state (row/
  structural equality); no trigger falsely rejects a valid journal-applied row.
  (VAL-016)
- [ ] Import-graph check shows no `internal/...` package imports `cmd/fid`; the
  enriched-read change stays in `internal/issues`. (VAL-017)

## TDD Plan

- First failing test: raw-SQL `INSERT INTO claims` with a dangling `ref_id`
  should RAISE — red before the trigger exists.
- Additional tests: dangling-ref reject + resolvable accept for `links` (both
  endpoints), `claims`, `activity_log`; link scope rejects (decision/escalation
  endpoints), relation/type compatibility accept+reject cases; claim ref_type
  reject (`qa_review`) + accept (`issue`); full-schema rebuild determinism with
  delete-DB parity; VAL-017 import-graph + adapter scan.
- Regression risks: composite-PK `links` endpoint resolution; `requirement`
  endpoints resolving against `requirements.id`; trigger ordering vs the existing
  insert effects; `relates_to` allowed across any supported endpoint pair.

## Documentation Updates

- `project-docs/docs/phases/phase-2/tasks/task-4-completion.md` - completion
  report citing VAL-013..017 evidence and gate outcomes.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: None.
- Audit or release-readiness impact: Triggers satisfy TRD-SEC-004 backstop
  verification; capture raw-SQL RAISE evidence for the security review.

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

This task closes the phase, so also run the conditional commands below and the
phase-level extras (build + benchmarks) so the completion report carries the
phase-end evidence the QA Lead needs:

```text
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go build -o /tmp/fid ./cmd/fid` is required for the phase-end full-schema
  rebuild smoke (`/tmp/fid init`, `/tmp/fid rebuild --by developer`).
- `go test -tags integration -race ./...` is required: rebuild-determinism and
  trigger backstop tests use temp DBs and may be tagged.
- `go test -bench . -benchmem ./...` for the rebuild performance spot-check
  (informational; tuning is Phase 3).

## Notes For Developer

- Test triggers via raw SQL (the non-CLI backstop path) per TRD-SEC-004, with a
  control proving valid journal-applied rows still project (no false positives).
- For VAL-016, build the determinism fixture with journal events (not direct
  SQL) so the rebuild proof exercises the proven write/replay path end-to-end.
- `links` is a composite-PK polymorphic edge: resolve both `(from_type,from_id)`
  and `(to_type,to_id)`. `requirement` endpoints resolve against
  `requirements.id`; `issue` endpoints against `issues.id`.
- Encode §10.1 relation/endpoint compatibility in the trigger: `implements`/
  `implemented_by` only issue↔requirement; issue↔issue relations (`blocks`,
  `duplicates`, `supersedes`, `caused_by`, `discovered_from` and inverses) only
  issue↔issue; `relates_to` across any supported endpoint pair.
- Keep `cmd/fid` a thin adapter (VAL-017); no domain/persistence logic leaks in.
- Coordinate the `ProjectionSchemaVersion` value with Tasks 1-3 so it stays
  monotonic and the final value reflects the highest schema generation; document
  the final value in the completion report.
