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
| CF-005 | Phase 2 planning (TL) | TRD SCHEMA §14.6 activity_log "inserts only during projection" guard trigger and the §4.3 `activity` effect op were deferred: no activity_log population path exists until mutable issue lifecycle writes ship, and a SQL session-scoped projection guard is unreliable under the pooled multi-connection model (DEC-006). Phase 2 ships the §14.2 polymorphic-ref backstop on `activity_log` and exercises it via `insert` effects. Re-deferred to the v0.1b mutation/workflow slice. | Phase 5 | Medium |
| CF-007 | Phase 4 VAL-010 / DEC-014 | Rejected `fid issue create` poisons the journal tail: `internal/issues.Create` appends the event before projecting it, so a create that the Task 1 sequence backstop trigger rejects still commits to the journal and breaks all later reads (`fid issue get` replays the invalid tail and fails) until the line is removed. Implementation conformance gap against TRD §4.3 (validate before append; append is the commit point). Fix: validate scheduling/sequence before append or add a pre-append projection-dry-run; add a regression test that a rejected create leaves the journal line count unchanged plus a poisoned-tail read-recovery test, and a VAL assertion for rejected-write journal hygiene. | Phase 5 | High |

CF-002 closed as satisfied at Phase 0 kickoff; see `../phases/phase-0/decisions.md`.
CF-004 (sequence-vs-scheduling backstop) closed at Phase 4 wrap — delivered by
Phase 4 Task 1 with CLI evidence in Task 3 (VAL-001 PASS); see
`../phases/phase-4/decisions.md`.
CF-006 (enriched read-model centralisation) closed at Phase 4 wrap — delivered by
Phase 4 Task 4 (VAL-006 PASS, `cmd/fid` adapter-boundary scan clean); see
`../phases/phase-4/decisions.md`.

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
