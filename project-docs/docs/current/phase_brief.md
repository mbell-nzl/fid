# Phase Brief - Phase 14: Telemetry + Pricing (v0.1e slice 1/4)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Replaced at each phase kickoff by the incoming Tech Lead.
>
> **Size target:** under 150 lines.

## Where We Are

**v0.1d is closed.** v0.1e (Telemetry, Composites, Maintenance) is the final
v0.1 version, sliced into **four phases (14-17)** — see the v0.1e phase-slicing
note in `10_IMPLEMENTATION_ROADMAP.md`.

**Phase 14 is active: Telemetry + pricing.** The first, self-contained v0.1e
slice. It leads the version because the Phase 15 `fid dev *` composites open and
close the `agent_runs` rows this phase introduces. Exit: an agent can record an
agent run with a snapshot price and computed USD cost, see an unknown-model
warning, and aggregate cost and tokens by role/model/epic.

## Phase 14 Scope

Source: roadmap v0.1e Phase 14; PRD-TELEM-001..004; PRD-NFR-010; TRD-ARCH-010.

- Expand `.fid/pricing.toml` shipped defaults to Anthropic + OpenAI + Google,
  with sub-cent price precision preserved (PRD-NFR-010).
- `fid run record`: one `agent_runs` row (role, model, tokens in/out, snapshot
  pricing, computed `cost_usd`, optional provider reconciliation fields), warn
  on unknown model (PRD-TELEM-001..003a).
- `fid stats tokens` / `fid stats cost`: aggregate by role/model/epic, all four
  read formats (PRD-TELEM-004, PRD-OUT-001..004).

Already in place (no rebuild risk): `agent_runs` table + all columns
(`migrations/0002`), generic insert-effect projection, `run.record` in
`03_TRD_SCHEMA.md` section 13, and `internal/config` pricing loader.

## Ordered Tasks

1. **Task 1 - Pricing defaults + precision.** Expand `DefaultPricingConfig` to
   three providers; fix emit precision for round-trip; update
   `04_TRD_CONFIG.md`. (VAL-001)
2. **Task 2 - `fid run record`.** New `internal/runs` package; `run.record`
   matrix row (copy section 13); cost computation + unknown-model warning +
   provider fields; `cmd/fid/run.go`; one `run.recorded` event, no
   `activity_log`. (VAL-002..006)
3. **Task 3 - `fid stats tokens`/`cost`.** New `internal/stats` package;
   `cmd/fid/stats.go`; aggregate by role/model/epic; four-format parity; NULL-cost
   handling. (VAL-007, VAL-008)

## Validation Contract Index

`project-docs/docs/phases/phase-14/qa/validation-contract.md` (VAL-001..009):

- VAL-001 pricing defaults (3 providers) round-trip (task-1)
- VAL-002 `fid run record` one event + `agent_runs` row, no activity_log (task-2)
- VAL-003 `cost_usd` = formula over snapshot prices (task-2)
- VAL-004 unknown model: NULL cost/prices + warning + exit 0 (task-2)
- VAL-005 optional provider fields persist/NULL without altering normalized (task-2)
- VAL-006 `run.record` matrix row matches section 13; show=enforce; unknown role
  rejected before append (task-2)
- VAL-007 `fid stats tokens` by role/model/epic, four formats (task-3)
- VAL-008 `fid stats cost` by role/model/epic, NULL-cost excluded from sum,
  four formats (task-3)
- VAL-009 human dogfood: record+cost, unknown-model warning, aggregates (phase-end)

## Planning Interpretations (confirm at kickoff)

- **PI-1** new `internal/runs` domain boundary for the `agent_runs` write;
  record in `02_TRD.md` (table + matrix key already specced).
- **PI-2** `fid run record` writes no `activity_log` row (telemetry is not a
  domain mutation; CF-011 stays Phase 16).
- **PI-3** unknown model = warning, not rejection: NULL cost/prices, exit 0.
- **PI-4** `--by epic` groups by the run's owning epic (null bucket for
  unattached); reconcile `06_DESIGN_PLAN_CLI.md` `stats tokens` `--by` to add
  `epic`; NULL-cost runs counted but excluded from cost sums.
- **PI-5** bump `pricing.toml` emit precision so sub-cent defaults round-trip.

## Standing Constraints (carry into v0.1e)

- New write commands refresh the projection with the held-lock freshness helper
  inside `journal.WithLock` (DEC-028); one logical write = one event.
- Validate-before-append on every rejection path (CF-007): non-zero exit,
  actionable error, zero new journal lines. For `fid run record` the only reject
  path is an unknown/denied role; unknown model is a warning, not a rejection.
- Append matrix rows per-implemented-action (DEC-026); copy `03_TRD_SCHEMA.md`
  section 13 booleans verbatim.
- Maintain `triage-log.md` per review window even on CONTINUE; pin recurring
  improvement notes to a carry-forward ID on first appearance.

## Carry-Forward Register (open, see `decisions_active.md`)

| ID | Item | Target |
|---|---|---|
| CF-008 | Resolution-reason field / cancel `--reason` flag | v0.2 |
| CF-011 | Broad `activity_log` reconciliation (with `fid activity list`) | v0.1e Phase 16 |
| CF-013 | `internal/issues/complete.go` `Complete` orchestration cleanup | v0.1e Phase 15 |
| CF-014 | Split `cmd/fid/issue.go` (984) and `internal/journal/replay.go` (485) | v0.1e Phase 15 |
| CF-015 | Function cleanup of `internal/labels/add.go` / `internal/issues/label.go` | v0.1e Phase 15 |

None of the carry-forward items land in Phase 14 (telemetry adds no
`activity_log` rows and does not touch the dev-composite/cleanup files).

## Open Questions

- [ ] Confirm PI-1..PI-5 at kickoff. None block planning.

## Next Planning Checkpoint

Phase 15 — `fid dev *` composites + CF-013/CF-014/CF-015 cleanup.
