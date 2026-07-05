# Phase 16 - Task 3: CF-011 backfill part 2 (concerns / unclaim / link add / dev gate) + timeline-consistency e2e

## Objective

Finish the CF-011 reconciliation (DC-A, approved at kickoff) and prove the
headline: the activity timeline is consistent across **every** write command.

1. Backfill the remaining unlogged writers, each gaining one `activity`
   effect inside its existing single event (PI-2), mirroring the existing
   emitters in shape:

   | Writer | Event type (existing) | New activity event name | ref |
   |---|---|---|---|
   | `concern raise` | `issue.concern_raised` | `concern_raised` | `issue` |
   | `concern address` | `issue.concern_addressed` | `concern_addressed` | `issue` |
   | `unclaim` | `issue.claim.released` | `claim_released` | `issue` |
   | `link add` | `link.added` | `link_added` | `issue` |
   | `dev gate` | `dev.gate.recorded` | `gate_recorded` | `issue` |

   `link add` mirrors the existing `link_removed` projection
   (`internal/links/remove.go`, DEC-031) — one row recording the declared link;
   an issue→requirement link logs an `issue`-scoped row.
   `dev gate` is in this task by DEC-043 because it appends an issue-scoped
   `gate_runs` fact and is required by the VAL-007 lifecycle proof.

2. Add a phase-end consistency e2e proving `fid activity list` reflects every
   write across a full lifecycle with no silent gaps, and that delete-and-rebuild
   reproduces an identical `activity_log` (VAL-007).

## Requirements

- PRD refs: PRD-OUT-017; PRD-DEV-004/010 (`dev gate`), PRD-DEV-005..006
  (concerns), PRD-CLAIM-004 (unclaim), PRD-LINK-001..006a (link add)
- TRD refs: `03_TRD_SCHEMA.md` §9, §14.2, §14.6; ADR-0001 (journal-canonical,
  deterministic rebuild); `02_TRD.md` §4.7 (`gate_runs`)
- Design refs: `06_DESIGN_PLAN_CLI.md` §8
- Quality refs: `08_QUALITY_PLAN.md` §5.2 (e2e + rebuild determinism)

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-16/qa/validation-contract.md`.

- VAL-006: `concern raise`, `concern address`, `unclaim`, `link add`, and
  `dev gate` each append one named `activity_log` row inside their single event;
  `link add` mirrors `link_removed`.
- VAL-007: full-lifecycle `fid activity list` covers every write with no gaps;
  delete-and-rebuild reproduces an identical `activity_log`.

## Scope

### Files To Create Or Modify

- `internal/issues/concern_raise.go` - append one `activity` effect
  (`concern_raised`, `ref_type=issue`, `details` = concern_id/summary) to the
  `issue.concern_raised` event.
- `internal/issues/concern_address.go` - append one `activity` effect
  (`concern_addressed`, `ref_type=issue`, `details` = concern_id) to the
  `issue.concern_addressed` event.
- `internal/claims/release.go` - append one `activity` effect (`claim_released`,
  `ref_type=issue`, `details` = released_by) to the `issue.claim.released`
  event; reuse `claims` `activityIDExists` (`internal/claims/shared.go`).
- `internal/links/add.go` - append one `activity` effect (`link_added`,
  `ref_type=issue`, `ref_id`=from issue, `details` = relation/to_type/to_id)
  mirroring `internal/links/remove.go`'s `link_removed` shape (DEC-031). Add the
  `activityIDExists` helper to `internal/links` if absent (mirror claims/issues).
- `internal/dev/gate.go` - append one `activity` effect (`gate_recorded`,
  `ref_type=issue`, `ref_id`=issue, `details` = gate_run_id/attempt_number/
  command/exit_code/agent_run_id when present) to the `dev.gate.recorded` event;
  reuse `internal/dev` `activityIDExists`.
- `internal/issues/*_test.go`, `internal/claims/release_test.go`,
  `internal/links/add_test.go`, `internal/dev/gate_test.go` - one event + one
  named activity row each; reject paths add zero rows; issue↔issue and
  issue→requirement link add both covered.
- `cmd/fid/e2e_v01e_activity_test.go` (new) - integration/e2e: drive a full
  lifecycle (create → `dev start` (claim) → ac add/verify → `dev gate` →
  `dev finish` → `link add`, plus an `unclaim` path and a `concern raise`/
  `address` path), assert `fid activity list` returns one entry per write with no
  gap, then delete `fid.sqlite`, rebuild, and assert `activity_log` is identical.

### Files To Read For Context

- `internal/links/remove.go` (canonical `link_removed` activity shape to mirror,
  DEC-031) and `internal/links/add.go` (event construction)
- `internal/claims/renew.go`, `internal/claims/acquire.go`,
  `internal/claims/shared.go` (claim activity emitter + `activityIDExists`)
- `internal/issues/concern_raise.go`, `internal/issues/concern_address.go`
- `internal/issues/transition.go` (activity `Values` shape)
- `internal/dev/gate.go`, `internal/dev/gate_test.go`,
  `internal/dev/dev.go` (`activityIDExists`)
- `cmd/fid/e2e_v01e_dev_test.go` (existing e2e + delete-and-rebuild pattern)
- `internal/journal/rebuild` helpers / `cmd/fid/rebuild.go` (rebuild invocation)
- Task 1 `internal/activity` (the read used to assert timeline coverage)

### Out Of Scope

- Task 2 writers (already backfilled).
- Excluded writers: `fid label add`, `fid run record`, `fid requirement sync`
  (DC-A).
- `fid changes` digest (v0.2c); maintenance audit + v0.1 closure (Phase 17).

## Acceptance Criteria

- [ ] `concern raise`, `concern address`, `unclaim`, `link add` each append
  exactly one event with one named `activity_log` row (`ref_type=issue`); a
  rejected attempt of each adds zero rows. (VAL-006)
- [ ] `dev gate` appends exactly one `dev.gate.recorded` event with one
  `gate_runs` row and one `gate_recorded` activity row (`ref_type=issue`);
  rejected gate attempts add zero rows. (VAL-006, DEC-043)
- [ ] `link add` activity mirrors the `link_removed` shape (DEC-031); an
  issue→requirement link add logs an `issue`-scoped `link_added` row. (VAL-006)
- [ ] A full-lifecycle e2e shows `fid activity list` with one entry per write and
  no silent gaps. (VAL-007)
- [ ] After deleting `fid.sqlite` and rebuilding, `activity_log` is identical to
  the pre-delete projection. (VAL-007)
- [ ] `03_TRD_SCHEMA.md` §9 `activity_log.event` vocabulary lists the new names.
  (doc)

## TDD Plan

- First failing test from the resumed point: `dev gate` now writes one
  `gate_recorded` activity row in the same event as the `gate_runs` row — fails
  red against current behavior (today zero rows).
- Additional tests: concern raise/address and link add named rows; reject paths
  add zero rows; issue→requirement link add; gate rejection adds zero rows; the
  e2e timeline-coverage + rebuild-equality assertions.
- Regression risks: one event per writer (PI-2); validate-before-append; mirror
  the `link_removed`/claim/transition shapes byte-for-byte so rebuild stays
  deterministic; the rebuild-equality check must compare a stable projection of
  `activity_log` (e.g. ordered rows minus volatile columns if any); functions
  ≤60 lines, files ≤400.

## Documentation Updates

- `project-docs/specs/03_TRD_SCHEMA.md` - extend the `activity_log.event`
  vocabulary comment with `concern_raised`/`concern_addressed`/`claim_released`/
  `link_added`/`gate_recorded` (combine with Task 2's edit if sequenced after
  it). Bump "Last updated".
- Completion report: confirm the full DC-A/DEC-043 writer set is now consistent
  and the CF-011 carry-forward can close.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: `link add` vs `link remove` activity shape parity
  (DEC-031).
- Audit or release-readiness impact: the rebuild-determinism proof (VAL-007)
  feeds the Phase 17 v0.1 closure smoke.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Record each command in
the completion report as `PASS`, or `BLOCKED` with reason. Record the
`internal/issues`, `internal/claims`, `internal/links`, and `cmd/fid` coverage
percentages.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` — **required**: this task adds the
  integration/e2e timeline-consistency + rebuild test and changes
  `cmd/fid`-exercised write paths; record PASS/BLOCKED.
- `go build -o /tmp/fid ./cmd/fid` — required: CLI build used by the e2e.

## Notes For Developer

- PI-2: one event per writer; the activity effect joins the existing effects.
- `unclaim` (`internal/claims/release.go`) currently emits only the release
  effect — add the `claim_released` activity effect alongside it in the same
  event.
- `link add` must mirror `internal/links/remove.go` exactly (DEC-031) so add/
  remove are symmetric in the timeline.
- `dev gate` is included by DEC-043. Add `gate_recorded` without changing the
  gate persistence contract: the event still remains one `dev.gate.recorded`
  journal line, now with one `gate_runs` insert and one `activity_log` effect.
- The VAL-007 rebuild-equality check is the determinism gate: capture
  `activity_log` before deleting `fid.sqlite`, rebuild from the journal, and
  assert equality.
- DC-A is human-confirmed for Phase 16; do not reopen it unless implementation
  reveals a new spec contradiction.
