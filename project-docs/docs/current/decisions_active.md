# Decisions - Active

> **Purpose.** Live carry-forward constraints and still-binding commitments.
> Loaded by default in every agent session. Keep this small.
>
> **What belongs here:**
> - Carry-forward items that are still open.
> - Scheduled commitments that affect upcoming phases.
> - Short-lived workflow exceptions that running agents need.
>
> **What does not belong here:**
> - Standing rules: use `decisions_principles.md`.
> - Closed phase decision history: use `../phases/phase-N/decisions.md`.
>
> **Size target:** under 150 lines.

## Pre-Scaffold Design Decisions

Three ADRs predate the phase system and live in `docs/adr/`:

- `0001-journal-canonical-db-derived.md` — journal is source of truth,
  SQLite is a derived projection. No daemon.
- `0002-configurable-roles-default-profile.md` — roles and permissions are
  configurable, seeded from a default profile.
- `0003-three-tier-qa-model.md` — validation contracts, QA plans, and QA
  reviews ship together as the supported quality model.

These remain binding unless a phase decision explicitly supersedes one.

## Carry-Forward Register

Tech Lead reviews this register at each phase start and schedules, re-defers,
closes, or cancels each item.

| ID | Source | Item | Target | Priority |
|---|---|---|---|---|
| CF-001 | Scaffold | Fill remaining `<TBD>` stack choices and quality commands in `AGENTS.md`. | First implementation phase | High |
| CF-002 | Scaffold | Split `scaffold/initial-spec.md` progressively into the seven spec files under `project-docs/specs/`. | Pre-Phase 0 | High |

## Scheduled Commitments

- The first implementation phase must validate that every command in
  `AGENTS.md` is concrete and executable for the chosen stack.
- Every planned phase must create
  `project-docs/docs/phases/phase-N/qa/validation-contract.md` with `VAL-NNN`
  assertions before task execution begins.
- Before release-readiness work, decide where held-out scenarios live and how
  implementation agents are prevented from loading them.

## TL-Lite Warm-Up State

- Autonomous TL-Lite resume enabled: No.
- Warm-up phases completed: 0 of 2.
- During warm-up, Phase Runner stops for human approval on every
  `VERDICT: RESUME`.

## Size Tripwire

| File | Target | Tripwire |
|---|---:|---:|
| `phase_brief.md` | 150 | 200 |
| `decisions_active.md` | 150 | 200 |
| `progress_current.md` | 100 | 133 |
| `tl_escalations_current.md` | 80 | 107 |
