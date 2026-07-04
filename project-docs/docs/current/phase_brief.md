# Phase Brief - Phase 16: Activity Timeline + Read-Only Escape Hatches (v0.1e slice 3/4)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Replaced at each phase kickoff by the incoming Tech Lead.
>
> **Size target:** under 150 lines.
>
> **State:** Phase 16 **PLANNED**. Validation contract (`VAL-001..013`) and the
> six ordered task specs are drafted. Two items need **human confirmation at
> kickoff** (DC-A CF-011 scope, DC-B CF-016 scheduling) before Tasks 2/3/6 start.

## Where We Are

Phase 15 (`fid dev *` composites + CF-013/014/015 cleanup) is CLOSED — PASS
(wrapped 2026-06-29). v0.1e (Telemetry, Composites, Maintenance) is the final
v0.1 version, sliced into four phases (14-17). Phases 14-15 done; **Phase 16 is
next**, Phase 17 closes v0.1.

## Phase 16 Scope (roadmap v0.1e Phase 16)

Source: PRD-OUT-017, PRD-OUT-008, TRD-SEC-002, PRD-OPS-013; CF-011; CF-016. No
new persistence contract; read-heavy. The only write-path change (CF-011) is
**additive and rebuild-safe** (no `event_schema` bump, no upgrader).

- **`fid activity list`** (PRD-OUT-017) — projected timeline read over
  `activity_log` with `--on`/`--since`/`--event` filters, all four formats (T1).
- **CF-011 reconciliation** — backfill the unlogged entity-state writers so the
  timeline is consistent across every write (T2-T3). Already projected:
  transition/update/issue-label/link-remove/claim. Backfill: issue create, ac
  add, ac verify, complete (T2); concern raise/address, unclaim, link add (T3).
- **`fid sql`** (PRD-OUT-008, TRD-SEC-002) — read-only escape hatch; reject
  non-SELECT/EXPLAIN/PRAGMA/WITH before execution + read-only connection (T4).
- **`fid journal inspect <event-id>`** (PRD-OPS-013) — read-only event
  inspection: actor, command, effects, location, canonical hash, revert verdict
  (T5).
- **CF-016** (T6, conditional) — `fid issue get --section gates` read over
  `gate_runs`.

## Decisions To Confirm With Human At Kickoff

- **DC-A — CF-011 timeline scope.** Every write that mutates/appends an **issue**
  fact (or its issue→requirement link) emits one `activity_log` row in its single
  event. Backfilled: issue create, ac add, ac verify, complete, concern raise,
  concern address, unclaim, link add. **Excluded:** `fid label add` (label
  creation — `label` is not a valid `activity_log.ref_type`; §14.2 allows only
  issue/requirement/decision/escalation), `fid run record` (telemetry, owned by
  `fid stats`), `fid requirement sync` (bulk upsert, noisy; requirement reads are
  v0.2). Event names: `issue_created`/`ac_added`/`ac_verified`/`task_completed`/
  `concern_raised`/`concern_addressed`/`claim_released`/`link_added`. **Changes
  projected timeline content + journal effect shape → confirm excluded edge
  cases + vocabulary.**
- **DC-B — CF-016 scheduling.** Add `fid issue get --section gates` into Phase 16
  as Task 6 (recommend yes; fallback v0.2c). Public-interface addition; no schema
  change. **Confirm scheduling.**

If the human defers either, Tasks 2/3 proceed only after DC-A confirmation;
Task 6 is cut on DC-B deferral (interim path: `fid sql` over `gate_runs` +
`fid journal inspect` of `dev.gate.recorded`).

## Validation Contract Index (`phase-16/qa/validation-contract.md`)

| VAL | Owning task | Assertion (short) |
|---|---|---|
| VAL-001 | T1 | `fid activity list` documented entry shape, chronological, all formats |
| VAL-002 | T1 | `--on <type>/<id>` includes only that entity |
| VAL-003 | T1 | `--since` boundary include/exclude |
| VAL-004 | T1 | `--event` includes only matching event |
| VAL-005 | T2 | create/ac add/ac verify/complete each emit one activity row; reject→zero |
| VAL-006 | T3 | concern raise/address/unclaim/link add each emit one row; link add mirrors link remove |
| VAL-007 | T3 | full-lifecycle timeline has no gaps; delete-and-rebuild reproduces `activity_log` |
| VAL-008 | T4 | `fid sql` accepts SELECT/EXPLAIN/PRAGMA/WITH, all formats |
| VAL-009 | T4 | `fid sql` rejects non-read before execution + read-only connection (TRD-SEC-002) |
| VAL-010 | T5 | inspect reports fields + hash==`HashEvent`; unknown id non-zero |
| VAL-011 | T5 | revert verdict event-type-keyed: inverse vs not-eligible+correction |
| VAL-012 | T6 | `--section gates` history (empty when none), distinct from completion snapshot |
| VAL-013 | phase-end (Human) | dogfood: read timeline, run `fid sql` + see reject, inspect an event |

## Ordered Task List

1. **T1 — `fid activity list`** (`internal/activity` + cmd). VAL-001..004.
2. **T2 — CF-011 backfill part 1** (create/ac add/ac verify/complete). VAL-005.
   *Needs DC-A.*
3. **T3 — CF-011 backfill part 2 + consistency e2e** (concern/unclaim/link add).
   VAL-006, VAL-007. *Needs DC-A.*
4. **T4 — `fid sql`** (`internal/sqlquery` + cmd + TRD-SEC-002 integration test).
   VAL-008, VAL-009.
5. **T5 — `fid journal inspect`** (`internal/journal` reader + revert classifier
   + cmd). VAL-010, VAL-011.
6. **T6 — CF-016 `--section gates`** (conditional on DC-B). VAL-012.

T1 and T4/T5 are independent of DC-A and can run first. T2-T3 deliver the CF-011
headline. VAL-013 (human dogfood) is generated by QA Lead at phase end.

## Standing Constraints (carry into v0.1e)

- **Reads are ungated** (PI-1): `fid activity list`/`sql`/`journal inspect`/
  `--section gates` take no `--by` and add no §13 matrix rows (§13 = writes only).
  Follow the read pattern: `journal.EnsureProjectionFresh` → open → query →
  `fidformat.Write`.
- **CF-011 is additive** (PI-2): the activity effect joins the existing single
  event; no second journal line; validate-before-append (rejected write appends
  zero rows, CF-007); no `event_schema` bump; mirror existing emitters byte-for-
  byte so rebuild stays deterministic.
- **`fid sql` defense in depth** (PI-3): leading-keyword allowlist (friendly
  error) **and** a read-only connection (hard backstop); reject multi-statement.
- **Required integration gate** (`go test -tags integration -race ./...`) must be
  recorded PASS/BLOCKED at every task touching `cmd/fid` or adding
  integration-tagged tests (Phase 14 ESC-001 lesson) — **explicitly required for
  T3 (e2e) and T4 (TRD-SEC-002)**.
- Maintain `triage-log.md` per review window even on CONTINUE; pin recurring
  improvement notes to a carry-forward ID on first appearance (Phase 15 lesson).
- Keep the template H1 in developer completion reports (Phase 15 lesson).

## Carry-Forward Register (open, see `decisions_active.md`)

| ID | Item | Target |
|---|---|---|
| CF-008 | Resolution-reason field / cancel `--reason` flag | v0.2 |
| CF-011 | Broad `activity_log` reconciliation (this phase, T2-T3) | **Phase 16** |
| CF-016 | First-class gate read surface (`--section gates`, T6) | **Phase 16 (DC-B) / v0.2c fallback** |
| CF-017 | `internal/dev` `Finish`/`Block` function-size cleanup | v0.2c |

CF-013/CF-014/CF-015 closed at Phase 15 wrap.

## Next Planning Checkpoint

**Phase 17** — `fid maintenance check` / `vacuum` + stub-table CLI guard + the
v0.1 closure exit-criteria smoke (closes v0.1). The Phase 16 VAL-007 rebuild
determinism proof feeds that closure smoke.
