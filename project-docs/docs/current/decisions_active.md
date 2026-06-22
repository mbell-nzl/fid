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
| CF-008 | Phase 5 DEC-015 (human HUMAN-013c) | Resolution type / cancellation reason: `cancelled` is a single undifferentiated bucket and cannot distinguish "redundant", "duplicate", "won't fix", etc. Closing a redundant backlog task already works via `todo → cancelled` (PRD-WF-013); the enhancement is a new resolution-reason field + cancel `--reason` flag (touches data model + `activity_log` audit), **not** a new workflow edge. Product-owner decision; do not widen v0.1 default graphs. | v0.2 (with configurable workflows PRD-WF-011/012) | Low (product) |
| CF-011 | Phase 7 DEC-020 | Reconcile `activity_log` (derived timeline projection) population **consistently** across all write commands. Today `fid issue update`/`transition` emit an `activity` effect but `fid issue create` does not; Phase 7 AC writes and Phase 8 completions (DEC-D) intentionally match `create` (no `activity_log` row) and defer. The journal already records every write — this is about the human-facing timeline table only (no rebuild risk). When `fid activity list` is designed (PRD-OUT-017), decide which write commands populate `activity_log` and backfill the gap (create + AC writes + completions + any other unlogged creations). | v0.1e (with `fid activity list` / activity-log work) | Medium |
| CF-013 | Phase 8 final review / Review 1 | `internal/issues/complete.go` `Complete` is a large orchestration function (~194 lines; file is 391, under the 400-line tripwire). Not a blocker — behavior and gates passed. Candidate for a function-level structure cleanup (extract narrow helpers for completion-event construction and append/project result assembly). | v0.1e, alongside the `fid dev finish`/`block` composites that also touch this path | Low |
| CF-014 | Phase 11 final review (Improvement) | `cmd/fid/issue.go` (984 lines) and `internal/journal/replay.go` (485 lines) exceed the 400-line non-test file target (AGENTS.md / `08_QUALITY_PLAN.md`). Not a blocker — behavior gates passed and the QA plan defines no file-size gate. Candidate for a behavior-preserving split; batch with the `fid dev *` composites that also touch `cmd/fid/issue.go`. | v0.1e, alongside CF-013 | Low |

CF-003 closed at Phase 10 wrap — role-discovery corrupt-tail tolerance reconciled
by Phase 10 Task 1 (DEC-D, VAL-001 PASS): `KnownRoles`/`collectRoles` now report a
malformed or un-upgradeable trailing journal line as corruption with journal
file/line context, matching the replayer (TRD-ARCH-002) and ADR-0001; the TRG-002
test rewrite is decision-backed, not weakening; see `../phases/phase-10/decisions.md`.
CF-012 closed at Phase 8 wrap — behavior-preserving split of
`internal/issues/issue.go` delivered by Phase 8 Task 4 (VAL-009 PASS):
`create.go` 387, `get.go` 120, residual `issue.go` 192, no exported-surface or
behavior change; see `../phases/phase-8/decisions.md`.
CF-009 closed at Phase 6 wrap — `fid issue update` one-field-per-line output
delivered by Task 1 (VAL-009 PASS), `06_DESIGN_PLAN_CLI.md` §`fid issue update`
updated; see `../phases/phase-6/decisions.md`.
CF-010 closed at Phase 6 wrap — `transition.go` (262), `update.go` (260),
`project.go` (137) all ≤400 lines after behavior-preserving split (Task 2,
VAL-010 PASS); see `../phases/phase-6/decisions.md`.
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
