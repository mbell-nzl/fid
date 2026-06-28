# Phase Brief - v0.1e Kickoff (Planning)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Replaced at each phase kickoff by the incoming Tech Lead.
>
> **Size target:** under 150 lines.

## Where We Are

**v0.1d is closed.** Phase 13 (Labels + v0.1d closure) wrapped 2026-06-28 on
branch `phase-13` with final review PASS WITH IMPROVEMENTS and human checklist
PASS. The v0.1 headline lifecycle exit criterion is proven end-to-end and
rebuild-stable: create → claim → transition → complete → unclaim with permission
enforcement, a typed link, and a label.

**No phase is active.** The next target is **v0.1e — Telemetry, Composites,
Maintenance**, the final v0.1 version. The incoming Tech Lead slices v0.1e into
phases (this brief is a planning placeholder, not a task plan).

## v0.1e Scope (roadmap, to be sliced)

Source: `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` v0.1e.

- `fid run record` + `.fid/pricing.toml` defaults: `agent_runs` row with role,
  model, tokens in/out, snapshot pricing, computed `cost_usd`, optional provider
  reconciliation; warn on unknown model. (PRD-TELEM-001..003a, PRD-NFR-010)
- `fid stats cost` / `fid stats tokens`: aggregate by role/model/epic.
  (PRD-TELEM-004)
- `fid dev start` / `gate` / `finish` / `block` composites: claim + transition +
  agent_run + `gate_runs` + completion orchestration. (PRD-DEV-004, 008..010,
  PRD-CLAIM-005)
- `fid sql` read-only escape hatch (SELECT/EXPLAIN/PRAGMA/WITH only).
  (PRD-OUT-008, TRD-SEC-002)
- `fid journal inspect` by `event_id` (read-only). (PRD-OPS-013)
- `fid maintenance check` (consistency audit, `--warn-as-error`) / `vacuum`.
  (PRD-OPS-005..008)
- `fid activity list`: projected timeline read with entity/event/date filters —
  the home for the CF-011 broad `activity_log` reconciliation. (PRD-OUT-017)
- Stub-table CLI guard: `decision create`, `escalation create`, etc. return
  "commands for <entity> ship in <version>" with non-zero exit. (PRD-ERR-012)

## Slicing Guidance For The Incoming TL

- Slice as vertical, agent-runnable increments (roadmap Guiding Principle #3).
  Each phase must exit on a command an agent can run against a real `.fid/`
  workspace.
- Likely ordering: telemetry (`run record` + pricing + `stats`) is self-contained
  and a good first slice; the `fid dev *` composites build on existing
  claim/transition/complete and are the natural home for CF-013/CF-014/CF-015
  cleanup; `fid activity list` carries the CF-011 reconciliation; maintenance and
  the read-only `fid sql` / `fid journal inspect` surfaces can batch.
- Confirm the v0.1 exit criteria (roadmap "v0.1 exit criteria") are all met by
  the end of v0.1e before declaring v0.1 done.

## Standing Constraints (carry into v0.1e)

- Every new write command refreshes the projection with the held-lock freshness
  helper inside `journal.WithLock` (DEC-028); one logical write = one event.
- Validate-before-append on every rejection path (CF-007): non-zero exit,
  actionable error, zero new journal lines.
- Append matrix rows per-implemented-action (DEC-026); copy `03_TRD_SCHEMA.md`
  §13 booleans verbatim, edit the spec only if implementation proves it wrong.
  Any new v0.1e write command needs its `<action>` matrix row.
- Maintain `triage-log.md` per review window even on CONTINUE (Phase 13 gap), and
  pin recurring improvement notes to a carry-forward ID on first appearance.

## Carry-Forward Register (open, see `decisions_active.md`)

| ID | Item | Target |
|---|---|---|
| CF-008 | Resolution-reason field / cancel `--reason` flag | v0.2 |
| CF-011 | Broad `activity_log` reconciliation across all writes (with `fid activity list`) | v0.1e |
| CF-013 | `internal/issues/complete.go` `Complete` orchestration cleanup | v0.1e (with `fid dev finish`/`block`) |
| CF-014 | Split `cmd/fid/issue.go` (984) and `internal/journal/replay.go` (485) under 400 lines | v0.1e |
| CF-015 | Function-level cleanup of `internal/labels/add.go` `Add` and `internal/issues/label.go` add/remove flows | v0.1e (with CF-013/CF-014) |

## Open Questions

- [ ] None blocking. v0.1e slicing decisions belong to the incoming Tech Lead at
  phase kickoff.

## Next Planning Checkpoint

After v0.1e, v0.1 is complete (full issue layer dogfoodable). Next is **v0.2 —
Governance, Workflow Engine, Reports**, starting with v0.2a (workflow engine
migration to data-backed workflows + `role_permissions` table).
