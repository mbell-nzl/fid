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
| CF-003 | Phase 1 TRG-002 | Revisit role-validation corrupt-tail tolerance: the bootstrap role-discovery scan in `internal/policy` tolerates a malformed trailing journal line after seeded roles are found. Acceptable Phase 1 debt; reconcile before mutable role events or richer policy semantics ship. | Before v0.1d role permission matrix / mutable role events | Medium |
| CF-004 | Phase 2 planning (TL) | Sequence-vs-scheduling backstop trigger (TRD SCHEMA §14.1 `issues.sequence`) was deferred: the cross-table rule (tasks under sequential epics require a sequence; prioritised/standalone must not) is coupled to v0.1b scheduling semantics and has an ambiguous missing-`epic_metadata` case. Static positive-`sequence` CHECK + `UNIQUE(parent_id, sequence)` already ship. Add the cross-table backstop when `epic_metadata` population, tasks-under-epics, and `fid ready` land. | v0.1b Issue Lifecycle | Medium |
| CF-005 | Phase 2 planning (TL) | TRD SCHEMA §14.6 activity_log "inserts only during projection" guard trigger and the §4.3 `activity` effect op were deferred: no activity_log population path exists until v0.1b, and a SQL session-scoped projection guard is unreliable under the pooled multi-connection model (DEC-006). Phase 2 ships the §14.2 polymorphic-ref backstop on `activity_log` and exercises it via `insert` effects. Add the projection-only guard + `activity` op when activity_log population ships. | v0.1b (activity_log population) | Medium |
| CF-006 | Phase 3 TRG-001 / review-3 | Enriched issue read-model centralisation: `fid issue get --format markdown` assembles parent title, labels, acceptance criteria, latest verifications, and links inside `cmd/fid/issue.go` because `internal/issues.Get` does not yet expose that full enriched view. Acceptable for one read command, but it brushes TRD-ARCH-005 (`cmd/fid` stays thin) if a second read command reuses the shape. Centralise the enriched read model in `internal/issues` before reusing it across read commands. | v0.1b (`fid issue get --section` + `fid issue list`) | Medium |

CF-002 closed as satisfied at Phase 0 kickoff; see `../phases/phase-0/decisions.md`.

## Environment Constraints

- Go, module, and golangci-lint caches are redirected into gitignored in-repo
  dirs (`.gocache/`, `.gomodcache/`, `.gopath/`, `.golangci-cache/`, `.cache/`)
  via `[shell_environment_policy.set]` in `.codex/config.toml` so quality gates
  work under the workspace-write sandbox. See DEC-001 in
  `../phases/phase-0/decisions.md`. The env policy applies at session start, so
  a fresh agent session is required for it to take effect.

## Scheduled Commitments

- Every planned phase must create
  `project-docs/docs/phases/phase-N/qa/validation-contract.md` with `VAL-NNN`
  assertions before task execution begins.
- Before release-readiness work, decide where held-out scenarios live and how
  implementation agents are prevented from loading them.

## TL-Lite Resume Policy

- Autonomous TL-Lite resume: **ENABLED (permanent)**. The warm-up safeguard is
  retired — see DEC-002 in `../phases/phase-0/decisions.md`.
- Phase Runner does not stop on `VERDICT: RESUME`; it records the TRG-NNN and
  continues. `VERDICT: BLOCK` always stops the run.
- No warm-up tally to maintain at phase wrap.

## Size Tripwire

| File | Target | Tripwire |
|---|---:|---:|
| `phase_brief.md` | 150 | 200 |
| `decisions_active.md` | 150 | 200 |
| `progress_current.md` | 100 | 133 |
| `tl_escalations_current.md` | 80 | 107 |
