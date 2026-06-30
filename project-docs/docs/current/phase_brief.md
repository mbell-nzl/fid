# Phase Brief - Phase 16 (kickoff pending): Activity Timeline + Read-Only Escape Hatches (v0.1e slice 3/4)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Replaced at each phase kickoff by the incoming Tech Lead.
>
> **Size target:** under 150 lines.
>
> **State:** Phase 15 is CLOSED (PASS). This is a transitional kickoff brief for
> Phase 16. The Phase 16 validation contract (`VAL-NNN`) and ordered task specs
> are drafted at Phase 16 **Start**, not yet present.

## Where We Are

**Phase 15 (`fid dev *` composites + CF-013/014/015 cleanup) is closed — PASS**
(wrapped 2026-06-29). The four developer composites ship one-event each
(`dev start/gate/finish/block`), the three scheduled cleanups closed, and the
human dogfood (VAL-012) passed. See `../phases/phase-15/` for the archive.

v0.1e (Telemetry, Composites, Maintenance) is the final v0.1 version, sliced into
four phases (14-17) — see the v0.1e phase-slicing note in
`10_IMPLEMENTATION_ROADMAP.md`. Phases 14-15 are done; **Phase 16 is next**.

## Phase 16 Scope (from roadmap v0.1e)

Source: roadmap v0.1e Phase 16; PRD-OUT-017, PRD-OUT-008, TRD-SEC-002,
PRD-OPS-013; CF-011; **CF-016 pending human scheduling confirmation**.

- **`fid activity list`** — projected timeline read over `activity_log` with
  entity / event / date filters, all four read formats (PRD-OUT-017).
- **CF-011 broad `activity_log` reconciliation** — decide which write commands
  populate `activity_log` and backfill the unlogged writers (issue create, AC
  writes, completions, and any other write that does not yet emit a row) so the
  timeline is consistent across every write. Derived-projection only, **no
  rebuild risk** (the journal already records every write). Per-command
  projections already delivered for claim renewal (Phase 11), link removal
  (Phase 12), and label add/remove (Phase 13).
- **`fid sql`** — read-only escape hatch; reject anything that is not
  SELECT/EXPLAIN/PRAGMA/WITH (PRD-OUT-008, TRD-SEC-002).
- **`fid journal inspect <event_id>`** — read-only journal event inspection:
  actor, command metadata, effects, journal location, event hash, revert
  eligibility (PRD-OPS-013).

Already in place (no rebuild risk): the `activity_log` table + triggers, the
journal/projection machinery, the effect projector, and the permissions matrix.
No schema migration expected.

## Decision To Confirm At Phase 16 Kickoff

- **CF-016 (gate read surface).** Surfaced by the human at Phase 15 wrap:
  `fid dev gate` writes append-only `gate_runs`, but there is no first-class
  issue-scoped read surface for recorded gates. **Tech Lead recommendation:** add
  `fid issue get <id> --section gates` (fits the existing issue-section read model;
  keeps raw append-only history distinct from the completion snapshot
  `task_completions.gate_results`) rather than a new `fid dev gate list` verb, and
  schedule it **into Phase 16** since `fid sql` and `fid journal inspect` already
  land here and the `--section gates` addition is a small additive read over the
  existing `gate_runs` table (no schema change). Fallback if Phase 16 has no room:
  v0.2c alongside `fid run show` / `fid stats gates`. **This changes a public
  interface and touches product intent — confirm scope with the human at kickoff.**
  Note the interim path lands in Phase 16 regardless (`fid sql` SELECT over
  `gate_runs`; `fid journal inspect` of `dev.gate.recorded`).

## Standing Constraints (carry into v0.1e)

- New write commands refresh the projection with the held-lock freshness helper
  inside `journal.WithLock` (DEC-028); one logical write = one event.
- Validate-before-append on every rejection path (CF-007): non-zero exit,
  actionable error, zero new journal lines. (Mostly relevant to any new writes;
  Phase 16 is read-heavy, but CF-011 backfill touches projection discipline.)
- `fid sql` must reject non-read statements before executing (TRD-SEC-002).
- Append matrix rows per-implemented-action (DEC-026); copy `03_TRD_SCHEMA.md`
  section 13 booleans verbatim for any new actions (`sql`, `journal.inspect`,
  `activity.list` if gated).
- Required integration gate (`go test -tags integration -race ./...`) must be
  recorded PASS/BLOCKED at every task touching `cmd/fid` or adding
  integration-tagged tests (Phase 14 ESC-001 lesson; held in Phase 15).
- Maintain `triage-log.md` per review window even on CONTINUE; **pin recurring
  improvement notes to a carry-forward ID on first appearance** (Phase 15 lesson:
  the function-size debt was pinned a phase late).
- Keep the template H1 in developer completion reports (Phase 15 Review-1 note).

## Carry-Forward Register (open, see `decisions_active.md`)

| ID | Item | Target |
|---|---|---|
| CF-008 | Resolution-reason field / cancel `--reason` flag | v0.2 |
| CF-011 | Broad `activity_log` reconciliation (with `fid activity list`) | **Phase 16** |
| CF-016 | First-class issue-scoped gate read surface (`--section gates` rec.) | **Phase 16 rec. / v0.2c fallback (human to confirm)** |
| CF-017 | `internal/dev` `Finish`/`Block` function-size cleanup | v0.2c |

CF-013/CF-014/CF-015 all closed at Phase 15 wrap.

## Next Planning Checkpoint

**Phase 16 Start (this brief's phase):** draft
`../phases/phase-16/qa/validation-contract.md` (`VAL-NNN`), break into ordered
one-session tasks, confirm CF-016 scheduling with the human, and refresh this
brief with the validation contract index.

Then **Phase 17** — `fid maintenance check` / `vacuum` + stub-table CLI guard +
the v0.1 closure exit-criteria smoke (closes v0.1).
