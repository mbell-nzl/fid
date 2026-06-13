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
| CF-008 | Phase 5 DEC-015 (human HUMAN-013c) | Resolution type / cancellation reason: `cancelled` is a single undifferentiated bucket and cannot distinguish "redundant", "duplicate", "won't fix", etc. Closing a redundant backlog task already works via `todo → cancelled` (PRD-WF-013); the enhancement is a new resolution-reason field + cancel `--reason` flag (touches data model + `activity_log` audit), **not** a new workflow edge. Product-owner decision; do not widen v0.1 default graphs. | v0.2 (with configurable workflows PRD-WF-011/012) | Low (product) |
| CF-009 | Phase 5 DEC-016 (human HUMAN-013a) | `fid issue update` confirmation renders all changed fields as one long comma-separated line; render one updated field per line for readability. Non-blocking CLI polish; touches `06_DESIGN_PLAN_CLI.md` + `cmd/fid` update presenter. Fold into the next `fid issue update`-adjacent change. | Next update-adjacent change | Low |
| CF-010 | Phase 5 QA final review (review-4) | Non-test file size above `08_QUALITY_PLAN.md` §4.4 400-line target: `internal/issues/transition.go` (426), `internal/issues/update.go` (551), `internal/journal/project.go` (524). Split event assembly from append/project orchestration; reduce size rather than adding behavior. Behaviour/contract correct — improvement only. | Next transition/update/projection-adjacent change | Low |

CF-002 closed as satisfied at Phase 0 kickoff; see `../phases/phase-0/decisions.md`.
CF-004 (sequence-vs-scheduling backstop) closed at Phase 4 wrap — delivered by
Phase 4 Task 1 with CLI evidence in Task 3 (VAL-001 PASS); see
`../phases/phase-4/decisions.md`.
CF-005 (`activity` effect op + `activity_log` projection discipline) closed at
Phase 5 wrap — delivered by Task 2 (op + sole-writer discipline, SQL session
guard still deferred per DEC-006) and the Task 3/5 emitters (VAL-004 PASS); see
`../phases/phase-5/decisions.md`.
CF-006 (enriched read-model centralisation) closed at Phase 4 wrap — delivered by
Phase 4 Task 4 (VAL-006 PASS, `cmd/fid` adapter-boundary scan clean); see
`../phases/phase-4/decisions.md`.
CF-007 (rejected-write journal poisoning) closed at Phase 5 wrap — delivered by
Task 1/1b validate-before-append seam (VAL-005 PASS, reinforced by VAL-007); see
`../phases/phase-5/decisions.md`.

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
