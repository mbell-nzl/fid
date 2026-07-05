# Phase 16 Validation Contract - Activity Timeline + Read-Only Escape Hatches

> Phase: 16 - `fid activity list` + CF-011 `activity_log` reconciliation + read-only
> escape hatches `fid sql` / `fid journal inspect` (roadmap v0.1e, slice 3 of 4)
> Author: Tech Lead
> Last refined: 2026-07-05 by Tech Lead (kickoff confirmation)
>
> Location: `project-docs/docs/phases/phase-16/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 16 is complete. The contract specifies **what must be proven**; the QA
plan specifies **how**.

If a required assertion cannot be satisfied, stop and escalate rather than
weakening it.

## Phase Scoping Notes

Phase 16 is the third v0.1e slice: the human-facing **activity timeline** read
(`fid activity list`), the long-scheduled **CF-011** broad `activity_log`
reconciliation, and the two read-only **escape hatches** (`fid sql`,
`fid journal inspect`). It adds **no new persistence contract**: `activity_log`,
`gate_runs`, the journal/projection machinery, and the effect projector all
exist (`migrations/0002_v01_tables.up.sql`; `03_TRD_SCHEMA.md` §§7, 9, 14.2,
14.6). Everything here is read-heavy; CF-011 is the only write-path change and is
**additive and rebuild-safe** (the journal already records every write — old
events replay unchanged, new events carry one extra `activity` effect).

Scope:

- **`fid activity list`** (PRD-OUT-017): projected timeline read over
  `activity_log` with `--on <type>/<id>`, `--since <date>`, `--event <name>`
  filters and all four read formats (`06_DESIGN_PLAN_CLI.md` §8).
- **CF-011** reconciliation: backfill the unlogged entity-state writers so the
  timeline is consistent across every write (Phase 7 DEC-020). Per-command
  activity is already projected for transition/update/issue-label/link-remove/
  claim writes (Phases 5/11/12/13/15).
- **`fid sql`** (PRD-OUT-008, TRD-SEC-002): read-only escape hatch; reject
  anything that is not SELECT/EXPLAIN/PRAGMA/WITH (`06_DESIGN_PLAN_CLI.md` §15).
- **`fid journal inspect <event-id>`** (PRD-OPS-013): read-only event inspection
  — actor, command metadata, effects, journal location, canonical hash, and the
  event-type-keyed revert verdict (`06_DESIGN_PLAN_CLI.md` §16).
- **CF-016** (Task 6, confirmed at kickoff): `fid issue get --section gates`
  read over the append-only `gate_runs` history.

## Scope Decisions

The kickoff items were human-confirmed on 2026-07-05 (workspace timezone).
DEC-043 was added by the Tech Lead during the Task 3 unblock. They are recorded
in `../decisions.md`; Developers proceed without re-escalating these points.

- **DC-A — CF-011 timeline scope approved.** Every write command that mutates or appends
  a fact about an **issue** (or its issue→requirement link) emits exactly one
  `activity_log` row inside its single journal event. **Backfilled writers:**
  `issue create`, `issue ac add`, `issue ac verify`, `issue complete`,
  `concern raise`, `concern address`, `unclaim`, `link add`, and `dev gate`
  (confirmed by DEC-043 after the Task 3 blocker). **Deliberately excluded:**
  `fid label add` (label *creation* — `label` is not a valid
  `activity_log.ref_type`; the §14.2 trigger allows only
  `issue`/`requirement`/`decision`/`escalation`, and issue-label add/remove
  already log with `ref_type=issue`); `fid run record` (telemetry — owned by
  `fid stats`, not an entity-timeline event in v0.1); `fid requirement sync`
  (bulk pointer upsert — noisy, and requirement reads/drift are v0.2). Event-name
  vocabulary (mirroring `status_changed`/`issue_updated`/`issue_label_added`/
  `link_removed`/claim events): `issue_created`, `ac_added`, `ac_verified`,
  `task_completed`, `concern_raised`, `concern_addressed`, `claim_released`,
  `link_added`, `gate_recorded`. (VAL-005/006/007)
- **DC-B — CF-016 scheduled in Phase 16.** Add
  `fid issue get <id> --section gates` (read-only `gate_runs` history) as Task 6.
  Public-interface addition; no schema change (`gate_runs` exists). (VAL-012)
- **Mode — Heavy Mode approved.** Phase 16 adds public CLI surfaces and includes
  the TRD-SEC-002 read-only SQL security requirement, so it follows
  `09_AI_WORKFLOW.md` §4.3.

## Planning Interpretations (recorded as phase decisions during execution)

- **PI-1 (reads are ungated).** `fid activity list`, `fid sql`,
  `fid journal inspect`, and `fid issue get --section gates` are read-only: no
  `--by` flag and no permissions-matrix rows (`03_TRD_SCHEMA.md` §13 governs
  write commands only). They follow the read freshness pattern
  (`journal.EnsureProjectionFresh` then open + query), matching `fid issue
  list`/`fid stats`.
- **PI-2 (CF-011 is additive, rebuild-safe, no schema bump).** Adding an
  `activity` effect to *new* events of the backfilled writers does not change how
  *prior* events replay (they simply lack the row — the existing behavior).
  `event_schema` stays v1; the effect projector already supports the `activity`
  op; no upgrader is required. Each backfilled writer still emits exactly **one**
  event (the activity effect joins the existing effects), and
  validate-before-append is preserved (a rejected write appends zero rows).
  Mirror the existing per-command emitters (`internal/issues/transition.go`,
  `internal/links/remove.go`, `internal/claims/renew.go`) in shape so rebuild
  stays deterministic.
- **PI-3 (`fid sql` is read-only by construction).** A leading-keyword allowlist
  (SELECT/EXPLAIN/PRAGMA/WITH, case-insensitive, comment/whitespace-stripped)
  produces the friendly rejection error; the underlying SQLite connection is
  opened read-only (e.g. `query_only`/`mode=ro`) as the hard guarantee so a
  mutation that slips past the classifier (`WITH ... DELETE`, a `;`-separated
  second statement) still cannot alter the DB. Reject multi-statement input.
  (TRD-SEC-002)
- **PI-4 (new read packages + journal placement).** `internal/activity` (timeline
  read) and `internal/sqlquery` (read-only escape hatch) are new domain read
  packages recorded in `02_TRD.md` §§2-3. The journal-inspect reader (find event
  by ID with file/line + canonical hash + revert verdict) lives in
  `internal/journal`; the event-type-keyed revert classifier will be reused and
  extended by v0.2 `fid event revert`.
- **PI-5 (`fid journal inspect` revert verdict is report-only).** v0.1 only
  *reports* the verdict (mutable-state → `inverse` with restore + precondition;
  append-only fact → `retraction`; creation/completion/claim/telemetry/composite
  → `eligible:false` with a forward `correction`). The classifier is exhaustive
  over current event types with a safe default (unknown → not eligible + generic
  correction). No compensating event is appended (v0.2).
- **PI-6 (composites inherit the completion activity row).** Once `fid issue
  complete` emits a `task_completed` activity row (DC-A), the shared
  `issues.BuildCompletionEffects` consumed by `internal/dev` finish/block emits
  it too, so the composites stay consistent (still one event each, now carrying
  `status_changed` + `task_completed` activity). Verify each composite still
  emits exactly one event.
- **PI-7 (`dev gate` emits a gate activity row).** Per DEC-043, `fid dev gate`
  records an issue-scoped `gate_recorded` activity row in the same
  `dev.gate.recorded` event as the `gate_runs` insert. This closes the VAL-007
  lifecycle gap without changing the gate persistence contract or event schema.

## Numbering

- IDs use `VAL-NNN`, scoped to Phase 16, sequential from `VAL-001`.
- Carry-forward assertions are re-issued with a new `VAL-NNN` and a source note.
- Retired/superseded assertions stay in the table; IDs are never reused.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | PRD-OUT-017; `06_DESIGN_PLAN_CLI.md` §8 | `fid activity list` returns projected `activity_log` entries with the documented shape (`id`, `journal_event_id`, `ref_type`, `ref_id`, `event`, `actor`, `occurred_at`, `details`) in chronological order, rendered in all four formats (json/toon/markdown/yaml). | CLI + domain: seed writes, run `fid activity list` in each format; assert the field set, ordering, and that `details` deserialises. | QA Lead | task-1 | Entries carry the documented fields in chronological order and all four formats render. | A documented field is missing, ordering is undefined, or a format fails to render. | Ready | `fid changes` (entity-grouped digest) is v0.2c, not this command. |
| VAL-002 | PRD-OUT-017; `06_DESIGN_PLAN_CLI.md` §8 | `fid activity list --on <type>/<id>` returns exactly the entries whose `(ref_type, ref_id)` match and excludes entries for other entities. | CLI: with activity for two issues, filter `--on issue/<id>`; assert only that issue's entries appear. | QA Lead | task-1 | Matching entries included; non-matching entries excluded. | A non-matching entry appears or a matching one is dropped. | Ready | Malformed `--on` handling is an error-path detail, not this assertion. |
| VAL-003 | PRD-OUT-017; `06_DESIGN_PLAN_CLI.md` §8 | `fid activity list --since <date>` includes entries with `occurred_at` on or after the boundary and excludes earlier ones. | CLI: with entries before and after a boundary, filter `--since`; assert the boundary inclusion/exclusion. | QA Lead | task-1 | On/after-boundary entries included; earlier entries excluded. | An earlier entry appears or an on/after entry is dropped. | Ready | Datetime vs date-only parsing nuance noted in the task; the boundary semantics are what is proven. |
| VAL-004 | PRD-OUT-017; `06_DESIGN_PLAN_CLI.md` §8 | `fid activity list --event <name>` includes only entries whose `event` equals the name and excludes others. | CLI: with mixed event names, filter `--event status_changed`; assert only those appear. | QA Lead | task-1 | Only matching-event entries appear. | An entry with a different event name appears or a matching one is dropped. | Ready | Combining multiple filters is exercised incidentally; each filter is proven independently. |
| VAL-005 | CF-011 (Phase 7 DEC-020); PRD-OUT-017; `03_TRD_SCHEMA.md` §§9, 14.2/14.6; DC-A; PI-2 | `issue create`, `issue ac add`, `issue ac verify`, and `issue complete` each append exactly one `activity_log` row (events `issue_created`/`ac_added`/`ac_verified`/`task_completed`, `ref_type=issue`) **inside their single journal event**, so `fid activity list --on issue/<id>` reflects each; a rejected attempt of each writer appends zero activity rows (validate-before-append). | Domain + CLI: run each writer, assert one new event and one new `activity_log` row of the named event; assert a rejection path adds zero rows; `fid activity list` shows the four entries. | QA Lead | task-2 | Each of the four writers adds exactly one event + one named activity row; rejections add zero. | A writer emits zero/extra activity rows, splits into >1 event, or a rejection appends a row. | Ready | Excluded writers (label create, run record, requirement sync) intentionally add no row (DC-A). |
| VAL-006 | CF-011 (Phase 7 DEC-020); PRD-OUT-017; PRD-DEV-004/010; DEC-031/DEC-033/DEC-043; DC-A; PI-2/PI-7 | `concern raise`, `concern address`, `unclaim`, `link add`, and `dev gate` each append exactly one `activity_log` row (events `concern_raised`/`concern_addressed`/`claim_released`/`link_added`/`gate_recorded`) inside their single journal event; `link add` mirrors the existing `link_removed` shape (`internal/links/remove.go`, DEC-031) and an issue→requirement link logs an `issue`-scoped row; `fid activity list` reflects each. | Domain + CLI: run each writer; assert one event + one named activity row each, including issue↔issue and issue→requirement link add; assert `dev gate` records `gate_recorded`; `fid activity list` shows them. | QA Lead | task-3 | Each of the five writers adds exactly one event + one named activity row; link add mirrors link remove; rejected attempts add zero rows. | A writer emits zero/extra rows, splits into >1 event, a rejection appends a row, or link add diverges from the link-remove shape. | Ready | `fid label add` (label creation) adds no activity row (DC-A: invalid ref_type). |
| VAL-007 | CF-011; PRD-OUT-017; `08_QUALITY_PLAN.md` §5.2; ADR-0001; DEC-043 | Over a full lifecycle (create → claim/`dev start` → ac add/verify → `dev gate` → `dev finish`/`dev block` and/or `unclaim` → `link add`), `fid activity list` returns one chronological entry per write command with **no silent gaps**, and delete-and-rebuild reproduces an identical `activity_log` (deterministic projection from the journal). | Integration/e2e: run the lifecycle on a real `.fid` workspace; assert the timeline covers every write including `gate_recorded`; delete `fid.sqlite`, rebuild, and assert `activity_log` is byte-for-byte equal. | QA Lead | task-3 | Timeline reflects every write with no gap; rebuild reproduces an identical `activity_log`. | A write is missing from the timeline or rebuild diverges. | Ready | Maintenance audit + v0.1 closure smoke are Phase 17. |
| VAL-008 | PRD-OUT-008; `06_DESIGN_PLAN_CLI.md` §15 | `fid sql <query>` accepts SELECT, EXPLAIN, a read-only PRAGMA, and `WITH ... SELECT`, returning result rows in all four formats. | CLI: run a SELECT (all formats), an EXPLAIN, a read PRAGMA, and a `WITH ... SELECT`; assert rows render. | QA Lead | task-4 | All four read statement kinds return rendered rows in every format. | A permitted read statement is rejected or a format fails to render. | Ready | A query DSL (`fid query`) is conditional v0.3; this is raw read-only SQL only. |
| VAL-009 | PRD-OUT-008; TRD-SEC-002; `06_DESIGN_PLAN_CLI.md` §15; PI-3 | `fid sql` rejects INSERT/UPDATE/DELETE/DROP and any other non-read statement before execution (non-zero exit, the documented `accepts only SELECT, EXPLAIN, PRAGMA query, and WITH` error), and the connection is read-only so a mutation that bypasses classification still cannot alter the DB. | Integration (TRD-SEC-002): attempt INSERT/UPDATE/DELETE/DROP via `fid sql`; assert each is rejected non-zero and the projected DB is unchanged; assert a multi-statement / `WITH ... DELETE` cannot mutate. | QA Lead | task-4 | Non-read statements rejected non-zero with the documented error and zero DB mutation. | A non-read statement runs, mutates the DB, or exits zero. | Ready | Defense in depth: classifier error + read-only connection are both required. |
| VAL-010 | PRD-OPS-013; `06_DESIGN_PLAN_CLI.md` §16 | `fid journal inspect <event-id>` for an existing event reports `event_id`, `event_schema`, `event_type`, `journal_file` + `line`, canonical `hash` (equal to `journal.HashEvent`), `actor`, `command`, and ordered `effects` in all four formats; an unknown event id exits non-zero with an actionable error. | CLI + domain: inspect a known event in each format and assert the fields + that `hash` equals `HashEvent` over the stored event; inspect a missing id and assert non-zero + error. | QA Lead | task-5 | Known event reports all documented fields with a hash matching `HashEvent`; unknown id exits non-zero with an actionable error. | A documented field is wrong/missing, the hash diverges, or a missing id is reported as success. | Ready | `fid event revert` (appending the compensating event) is v0.2. |
| VAL-011 | PRD-OPS-013; `06_DESIGN_PLAN_CLI.md` §16; PI-5 | The `revert` verdict is event-type-keyed: a mutable-state event (e.g. `issue.transitioned`) reports `eligible:true, mode:inverse` with `restores` + `precondition`; a non-reversible event (e.g. `issue.created`, `issue.claim.acquired`, a `dev.*` composite) reports `eligible:false` with a forward `correction`. | CLI: inspect a transition event (assert inverse verdict) and a creation/claim/composite event (assert `eligible:false` + `correction`). | QA Lead | task-5 | Mutable-state event reports inverse with restore+precondition; non-reversible event reports not-eligible with a correction. | A verdict is inverted, missing the restore/precondition, or missing the correction note. | Ready | Exact per-type mapping is refined by v0.2 `fid event revert`; v0.1 only reports the verdict. |
| VAL-012 | CF-016 (Phase 15 human wrap note); PRD-DEV-004/010; `03_TRD_SCHEMA.md` §7 | `fid issue get <id> --section gates` returns the issue's append-only `gate_runs` history (`command`, `exit_code`, `output_excerpt`, `ran_at`, `attempt_number`, `recorded_by_role`) for both in-progress and finished tasks in all four formats; an issue with no recorded gates returns an empty list; the section is distinct from the completion snapshot (`task_completions.gate_results`). | CLI: record gates on an in-progress task and read `--section gates`; read it again after `dev finish`; read an issue with no gates (assert empty); all four formats. | QA Lead | task-6 | Gate history returned for in-progress and finished tasks (empty when none), distinct from the completion snapshot, all formats. | History missing while gates exist, conflated with the completion snapshot, or a non-empty/empty state wrong. | Ready | CF-016 is scheduled in Phase 16 by DC-B. |
| VAL-013 | `08_QUALITY_PLAN.md` §§5.2, 5.6; phase goal | A human dogfooding a real `.fid` workspace can, from CLI output alone: (a) read an activity timeline that reflects each of their writes, (b) run a read-only `fid sql` SELECT and see a mutation rejected, and (c) inspect a journal event by ID. | Human checklist with one row per observable; QA Lead generates from this VAL after deterministic QA passes. | Human | phase-end | Human confirms the three read surfaces from CLI output without reading implementation. | Human cannot confirm one of the observables, or the checklist is not run. | Ready | QA Lead fans this into per-observable human-checklist rows (a/b/c). |

## Coverage Rules

Confirmed for Phase 16 execution:

- Phase-scope requirements each map to at least one `VAL-NNN`: PRD-OUT-017 to
  VAL-001/002/003/004 (+ consistency VAL-007, dogfood VAL-013); CF-011 to
  VAL-005/006/007; PRD-OUT-008 + TRD-SEC-002 to VAL-008/009 (+ VAL-013);
  PRD-OPS-013 to VAL-010/011 (+ VAL-013); CF-016 to VAL-012.
- Every `VAL-NNN` is referenced by a task spec or explicitly marked `phase-end`
  (VAL-013).
- Binary/conditional assertions name both states: VAL-002 (included/excluded),
  VAL-003 (on-after/before), VAL-004 (match/non-match), VAL-005/006 (row vs
  zero-on-reject), VAL-009 (read accepted vs write rejected + no mutation),
  VAL-010 (found/not-found), VAL-011 (inverse/not-eligible), VAL-012
  (history/empty).
- Determinism/parity proofs: VAL-007 carries a delete-and-rebuild `activity_log`
  equality proof; VAL-010 carries a canonical-hash parity proof (`hash` ==
  `journal.HashEvent`).
- Coverage is enforced at the gate, not as a `VAL-NNN`: `internal/activity`,
  `internal/sqlquery`, `internal/journal`, `internal/issues`, `internal/claims`,
  `internal/links`, and `cmd/fid` coverage must stay above thresholds; each task
  runs the required coverage gate and records package percentages.
- The phase touches `cmd/fid` and adds CLI paths exercised by integration-tagged
  tests, so the required integration gate (`go test -tags integration -race
  ./...`) must be recorded PASS (or BLOCKED with reason) at every task that adds
  or moves integration-tagged tests — explicitly required for `fid sql`
  (TRD-SEC-002 integration test, VAL-009) and the consistency e2e (VAL-007).
- QA Lead maps each `VAL-NNN` to deterministic `(qa-lead)`, mid-phase review, or
  `(human)` verification, and fans VAL-013 into per-observable human rows.
- Standing QA-plan checks carry forward: validate-before-append on every
  backfilled rejection path (VAL-005/006), one logical write = one event on each
  backfilled writer (VAL-005/006/007), deterministic rebuild of the widened
  projection (VAL-007), and read-only enforcement defense in depth (VAL-009).

## Outcomes

Phase-end final review records every `VAL-NNN` outcome: `PASS`, `FAIL`,
`BLOCKED`, or `CARRY-FORWARD`. `PASS` without evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-30 | VAL-001..013 | created | Phase 16 planning draft | roadmap v0.1e Phase 16; PRD-OUT-017/008; TRD-SEC-002; PRD-OPS-013; CF-011; CF-016 |
| 2026-07-05 | DC-A/DC-B/Mode; VAL-001..013 | marked ready | Human confirmed CF-011 scope, CF-016 scheduling, and Heavy Mode at kickoff | DEC-040; DEC-041; DEC-042 |
| 2026-07-05 | VAL-006/007 | added `dev gate` / `gate_recorded` to Task 3 CF-011 reconciliation | Task 3 blocker showed VAL-007 lifecycle could not prove no silent gaps while `dev.gate.recorded` emitted no activity row | DEC-043 |
