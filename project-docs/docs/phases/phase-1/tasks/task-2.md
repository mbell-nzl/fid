# Phase 1 - Task 2: Journal Event v1 Model, ID Generation, Canonical JSON, Upgrader Chain

## Objective

Add the pure-data core of the journal: the v1 event/effect Go types, canonical
JSON serialization (sorted keys, stable form) and event hashing, generated ID
helpers (`<PREFIX>-<hash>` and `EVT-<UTC YYYYMMDD>-<sequence>`), and the
event-schema upgrader-chain mechanism. No I/O, no DB, no lock — only data.

## Requirements

- PRD refs: PRD-OPS-001
- TRD refs: TRD §4.2 (ID generation), TRD §4.3 (event v1 contract, effect
  shape, canonical JSON, schema versioning, upgrader chain), TRD-ARCH-002,
  TRD-ARCH-003
- Design refs: N/A
- Quality refs: Quality Plan §5.1

## Validation Contract Coverage

- VAL-007: The v1 envelope carries `event_id`, `event_schema`, `event_type`,
  `occurred_at`, `actor`, `command`, and ordered `effects` (data-model side).
- VAL-009: The version-upgrader chain dispatches by `event_schema` and is a
  verified no-op at v1.
- VAL-015: Generated ID shapes (`<PREFIX>-<hash>` and event IDs) and
  collision-retry behavior at the data layer.

## Scope

### Files To Create Or Modify

- `internal/journal/event.go` (event + effect types, v1 schema constant)
- `internal/journal/canonical.go` (canonical JSON + event hash)
- `internal/journal/upgrade.go` (upgrader chain registry + v1 identity)
- `internal/journal/id.go` (or `internal/db/id.go`) ID generation helpers
- `internal/journal/*_test.go`
- `testdata/journal/` golden v1 event fixtures
- `project-docs/docs/phases/phase-1/tasks/task-2-completion.md`

### Files To Read For Context

- `project-docs/specs/02_TRD.md` §4.2, §4.3 (event v1 contract, effect shape,
  event-ID allocation, canonical JSON, upgrader chain), §4.3 example event
- `project-docs/specs/03_TRD_SCHEMA.md` §1 (ID shapes, prefixes)
- `project-docs/docs/phases/phase-1/qa/validation-contract.md`

### Out Of Scope

- Writing to disk, acquiring the lock, allocating sequences against real files
  (Task 3), and projecting effects into SQLite (Task 4).
- Any actual N→N+1 upgrader (none exists at v1); build only the chain mechanism
  and the v1 identity entry.

## Acceptance Criteria

- [ ] Event/effect types serialize to the v1 envelope with all required fields
  and an ordered `effects` array; a golden fixture round-trips. (VAL-007)
- [ ] Canonical JSON uses UTF-8, sorted object keys, no insignificant
  whitespace, and preserves array order; the same event hashes deterministically
  across runs. (VAL-007)
- [ ] The upgrader chain selects upgraders by `event_schema` and returns a v1
  event unchanged (identity) at the current version; a unit test asserts
  dispatch. (VAL-009)
- [ ] Issue ID helper produces `<PREFIX>-<hash>` with 6 hex chars by default
  and supports 8; event ID helper produces `EVT-<UTC YYYYMMDD>-<sequence>` from
  a given date + sequence. (VAL-015)
- [ ] A forced-collision scenario (caller-supplied existence predicate)
  retries with a fresh hash and returns a non-colliding ID. (VAL-015)

## TDD Plan

- First failing test: marshal a known event and assert canonical bytes + a
  stable hash.
- Additional tests: ID shape/length; event-ID formatting from date+sequence;
  collision retry; upgrader-chain dispatch is identity at v1; golden v1 event
  fixture round-trips through marshal/unmarshal.
- Regression risks: non-deterministic map ordering in JSON; float/precision in
  timestamps; mixing event IDs with row IDs.

## Documentation Updates

- `project-docs/docs/phases/phase-1/tasks/task-2-completion.md` - completion
  report citing VAL-007/009/015 evidence and gate results.

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

- None expected; this task is pure-data unit-tested.

## Notes For Developer

- Keep this package pure: no `os`, no DB, no clock reads inside the data types
  (callers pass timestamps/sequences). This is what makes replay pure data
  (TRD-ARCH-002) and the upgraders golden-testable (TRD-ARCH-003).
- Represent every primary key as a JSON object even for single-column keys
  (`{"id":"FID-a1b2c3"}`), per TRD §4.3.
- Capture at least one golden v1 event under `testdata/` so Task 4 can reuse it
  for replay tests (VAL-009).
