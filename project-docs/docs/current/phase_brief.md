# Phase Brief - Phase 15: `fid dev *` Composites + Cleanup (v0.1e slice 2/4)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Replaced at each phase kickoff by the incoming Tech Lead.
>
> **Size target:** under 150 lines.

## Where We Are

**Phase 14 (Telemetry + pricing) is closed — PASS.** v0.1e (Telemetry,
Composites, Maintenance) is the final v0.1 version, sliced into four phases
(14-17) — see the v0.1e phase-slicing note in `10_IMPLEMENTATION_ROADMAP.md`.

**Phase 15 is active: `fid dev *` composites + orchestration cleanup.** It ships
the developer composite commands on top of the existing claim/transition/
completion/`agent_runs` surfaces, plus the three scheduled carry-forward cleanups
on those same paths (CF-013/CF-014/CF-015). Exit: an agent can run `fid dev
start → gate → finish` to drive a task from claim through `dev_complete` in
composite steps, `fid dev block` to record a failed gate with reason, and each
composite emits a single auditable event with the correct claim/transition/
completion/gate effects.

## Phase 15 Scope

Source: roadmap v0.1e Phase 15; PRD-DEV-004, 008-010; PRD-CLAIM-005;
PRD-WF-007..010; TRD-ARCH-010; `02_TRD.md` §4.7; CF-013/CF-014/CF-015.

- New `internal/dev` package: each composite = **one journal event** carrying all
  effects, validate-before-append (CF-007), gated by the Phase 10 matrix.
- `fid dev start` — claim acquire + transition `todo → in_progress` + open
  `agent_runs` row.
- `fid dev gate` — append `gate_runs` on the derived active attempt (started
  state + live caller claim required).
- `fid dev finish` — close open run + completion(`dev_complete`, gate snapshot) +
  transition `in_progress → dev_complete` + release claim.
- `fid dev block` — completion(`failed_gate`, gate snapshot, reason) + transition
  `in_progress → failed_gate` + release claim.
- Matrix rows `dev.start/gate/finish/block` (section 13: developer + human only).
- CF-014 file split; CF-015 + CF-013 function extraction (behavior-preserving).

Already in place (no rebuild risk): `claims`/`issues`/`agent_runs`/`gate_runs`/
`task_completions` tables + triggers (`migrations/0002`), the generic effect
projector, the claim/transition/run/completion effect builders, and the Phase 10
permissions matrix. No schema migration.

## Ordered Tasks

1. **Task 1 - Cleanup (CF-014 + CF-015).** Split `cmd/fid/issue.go` (984) and
   `internal/journal/replay.go` (485) ≤400 lines; extract helpers in
   `internal/labels/add.go` and `internal/issues/label.go` ≤60-line functions.
   Behavior-preserving. (VAL-001, VAL-002)
2. **Task 2 - `internal/dev` + `dev start` + `dev gate`.** New package + the four
   matrix rows; `dev start` (claim+transition+open run) and `dev gate` (gate_runs
   on derived attempt). Spec updates: `02_TRD` boundary, CLI flags
   (`--model`, `--exit-code`, `--output-excerpt`). (VAL-003..006)
3. **Task 3 - CF-013 + `dev finish` + `dev block`.** Extract the reusable
   completion-effect builder from `complete.go` (CF-013), then the two closing
   composites consuming it. (VAL-006..010)
4. **Task 4 - E2E composite smoke + dogfood prep.** Integration test over both
   composite paths with permission denial + rebuild determinism. (VAL-011;
   VAL-012 phase-end)

## Validation Contract Index

`project-docs/docs/phases/phase-15/qa/validation-contract.md` (VAL-001..012):

- VAL-001 CF-014: `issue.go`/`replay.go` ≤400 lines, no behavior change (task-1)
- VAL-002 CF-015: label add/remove functions ≤60 lines, no behavior change (task-1)
- VAL-003 `dev start` one event = claim + `in_progress` + open run (task-2)
- VAL-004 `dev start` validate-before-append guards + claim-in-same-event (task-2)
- VAL-005 `dev gate` one `gate_runs` on derived attempt; records, not executes (task-2)
- VAL-006 `dev.*` matrix rows match section 13; show=enforce; deny rejected (tasks-2-3)
- VAL-007 `dev finish` one event = close run + completion + transition + release (task-3)
- VAL-008 `dev finish` enforces `all_acs_verified[developer]`; report in-event (task-3)
- VAL-009 `dev block` one event = completion(`failed_gate`) + transition + release (task-3)
- VAL-010 CF-013: shared completion-effect builder; `issue complete` unchanged (task-3)
- VAL-011 e2e composite smoke + rebuild determinism (task-4)
- VAL-012 human dogfood: start/gate/finish/block observable from CLI (phase-end)

## Decisions (confirmed by human 2026-06-29)

DN-A, DN-B, and PI-1..PI-6 are **all confirmed**. Developers may proceed without
re-escalating; record them as phase decisions during execution.

- **DN-A** `fid dev gate` **records** a reported result (does not execute the
  command); add `--exit-code` (required) + `--output-excerpt` (optional) to the
  CLI spec. Rationale: GLOSSARY "developers run gates"; ADR-0001 deterministic
  recorder; avoids a command-execution security surface. Public-interface
  addition. (VAL-005)
- **DN-B** `fid dev start --model <m>` (required) opens the base run, since
  `agent_runs.model` is `NOT NULL` and the spec gives start no model flag;
  run tokens/prices/cost stay NULL (Phase 14 `fid run record` remains the
  cost-capture path). Public-interface addition. (VAL-003)

## Planning Interpretations (confirmed 2026-06-29; recorded as decisions on execution)

- **PI-1** new `internal/dev` domain boundary; record in `02_TRD.md` §§2-3.
- **PI-2** one event per composite; effect-driven projector needs no upgrader;
  validate-before-append on every reject path (zero new journal lines).
- **PI-3** `dev start` claim-in-same-event satisfies PRD-CLAIM-005; targets
  `todo → in_progress`; rejects a live claim held by another role.
- **PI-4** `dev finish` rule ordering: `has_completion_report` satisfied by the
  in-event completion; reject before append only on unmet `all_acs_verified[
  developer]`.
- **PI-5** `dev finish`/`dev block` release the caller's claim (matches CLI
  output) though those target states stay started-category; plain
  `fid issue transition` semantics unchanged.
- **PI-6** `dev finish` closes the open run when one exists; `dev block` leaves
  runs alone.

## Standing Constraints (carry into v0.1e)

- New write commands refresh the projection with the held-lock freshness helper
  inside `journal.WithLock` (DEC-028); one logical write = one event — even for
  composites.
- Validate-before-append on every rejection path (CF-007): non-zero exit,
  actionable error, zero new journal lines.
- Append matrix rows per-implemented-action (DEC-026); copy `03_TRD_SCHEMA.md`
  section 13 booleans verbatim.
- Maintain `triage-log.md` per review window even on CONTINUE; pin recurring
  improvement notes to a carry-forward ID on first appearance.
- Required integration gate (`go test -tags integration -race ./...`) must be
  recorded PASS/BLOCKED at every task touching `cmd/fid` or adding
  integration-tagged tests (Phase 14 ESC-001 lesson).

## Carry-Forward Register (open, see `decisions_active.md`)

| ID | Item | Target |
|---|---|---|
| CF-008 | Resolution-reason field / cancel `--reason` flag | v0.2 |
| CF-011 | Broad `activity_log` reconciliation (with `fid activity list`) | v0.1e Phase 16 |
| CF-013 | `internal/issues/complete.go` `Complete` orchestration cleanup | **Phase 15 (Task 3)** |
| CF-014 | Split `cmd/fid/issue.go` (984) and `internal/journal/replay.go` (485) | **Phase 15 (Task 1)** |
| CF-015 | Function cleanup of `internal/labels/add.go` / `internal/issues/label.go` | **Phase 15 (Task 1)** |

CF-013/014/015 all land this phase. CF-011 stays Phase 16: the composites only
project the per-command `activity_log` rows their transition/claim writes already
emit; they do not widen or backfill CF-011.

## Open Questions

- [x] DN-A, DN-B, and PI-1..PI-6 confirmed by human 2026-06-29. None block any
  task; Task 2 implements the DN-A/DN-B flag additions in the CLI spec.

## Next Planning Checkpoint

Phase 16 — `fid activity list` + CF-011 reconciliation + `fid sql` /
`fid journal inspect`.
