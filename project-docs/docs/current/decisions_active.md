# Decisions - Active

> **Purpose.** Live carry-forward constraints and still-binding commitments.
> Loaded by default in every agent session. Keep this small.
>
> **Size target:** under 150 lines.

## Pre-Scaffold Design Decisions

Three ADRs predate the phase system and remain binding unless a phase decision
explicitly supersedes them:

- `docs/adr/0001-journal-canonical-db-derived.md` — journal is source of truth,
  SQLite is a derived projection. No daemon.
- `docs/adr/0002-configurable-roles-default-profile.md` — roles and permissions
  are configurable, seeded from a default profile.
- `docs/adr/0003-three-tier-qa-model.md` — validation contracts, QA plans, and
  QA reviews ship together as the supported quality model.

## Carry-Forward Register

Tech Lead reviews this register at each phase start and schedules, re-defers,
closes, or cancels each item.

| ID | Source | Item | Target | Priority |
|---|---|---|---|---|
| CF-008 | Phase 5 DEC-015 (human HUMAN-013c) | Resolution type / cancellation reason: `cancelled` is a single undifferentiated bucket and cannot distinguish "redundant", "duplicate", "won't fix", etc. The enhancement is a resolution-reason field + cancel `--reason` flag, touching the data model and `activity_log` audit. | v0.2 (with configurable workflows PRD-WF-011/012) | Low (product) |
| CF-017 | Phase 15 Review 1 + final review | Behavior-preserving function extraction for `internal/dev` `Finish` (~109 lines) and `Block` (~97 lines), both above the 60-line non-test function target. Reopen with the next `internal/dev` work. | v0.2c (`fid dev resume`/`abandon`) | Low |
| CF-018 | Phase 16 Review 2 + final review | Behavior-preserving split of `internal/sqlquery/sqlquery.go` (595 lines) under the 400-line non-test file target. Keep the public package API stable. | Next SQL package reopen | Low |
| CF-020 | Phase 17 DEC-056 | Burst issue writes stamp `occurred_at` ahead of the wall clock via `nextIssueUpdateTime`'s monotonicity bump, and `recent_activity_7d` excludes future rows. Bounded and self-correcting; do not fix inside a v0.1 closure audit. | v0.2 | Low (product) |
| CF-021 | Phase 17 Review 3 / STRUCT-003 | Add Go doc comments to exported `maintenance.VacuumResult` and `Vacuum`. | Next `internal/maintenance` reopen | Low |
| CF-022 | Phase 17 Review 1 | Add broader exact-output markdown characterization tests to lock rendering behavior beyond the current suite. | Next `internal/format` markdown reopen | Low |
| CF-023 | Phase 17 DEC-058 / retrospective | Workflow-tooling hardening: (a) worked backtick-wrapped commit-metadata example in the completion-report template, (b) Developer-prompt rule that the report is its own committed `docs(task)` step, (c) report-lint messaging distinguishing dirty vs parse failure, and (d) ESC-NNN global-uniqueness enforcement. Flagged for deep retrospective. | Workflow-tooling pass before v0.2a execution | Medium |

## Recently Closed

- **CF-019** closed at Phase 17 wrap — `internal/format/markdown.go` split into
  four non-test files (95/164/321/121 lines, all <=400) with byte-identical
  behavior; VAL-001 PASS.
- **v0.1 closed** at Phase 17 wrap (DEC-057) — every v0.1-targeted PRD row maps
  to passing evidence (VAL-024), static binary 12.27 MiB (<30 MB), and the Human
  accepted v0.1 for early dogfood (VAL-025 / HUMAN-25a..e PASS).
- **CF-011** closed at Phase 16 wrap — broad `activity_log` reconciliation
  delivered for issue create, AC add/verify, complete, concern raise/address,
  unclaim, link add, and `dev gate` / `gate_recorded`; VAL-005/006/007 PASS.
- **CF-016** closed at Phase 16 wrap — `fid issue get --section gates`
  delivered as a raw `gate_runs` read surface; VAL-012 PASS.
- Earlier closed carry-forward items are archived in their phase decisions:
  CF-013/014/015 (Phase 15), CF-003 (Phase 10), CF-012 (Phase 8), CF-009/010
  (Phase 6), CF-004/006 (Phase 4), CF-005/007 (Phase 5), CF-002 (Phase 0).

## Environment Constraints

- Go, module, and golangci-lint caches are redirected into gitignored in-repo
  dirs (`.gocache/`, `.gomodcache/`, `.gopath/`, `.golangci-cache/`, `.cache/`)
  via `[shell_environment_policy.set]` in `.codex/config.toml` so quality gates
  work under the workspace-write sandbox. See DEC-001 in
  `../phases/phase-0/decisions.md`. The env policy applies at session start, so
  a fresh agent session is required for changes to take effect.

## Scheduled Commitments

- Every planned normal/heavy phase must create
  `project-docs/docs/phases/phase-N/qa/validation-contract.md` with `VAL-NNN`
  assertions before task execution begins.
- Before release-readiness work, decide where held-out scenarios live and how
  implementation agents are prevented from loading them.

## TL-Lite Resume Policy

- Autonomous TL-Lite resume: **ENABLED (permanent)**. The warm-up safeguard is
  retired — see DEC-002 in `../phases/phase-0/decisions.md`.
- Phase Runner does not stop on `VERDICT: RESUME`; it records the TRG-NNN and
  continues. `VERDICT: BLOCK` always stops the run.
- Per DEC-047, task-level blocker reports, required gate failures,
  report-lint failures, missing/skipped/narrower command evidence, and ad hoc
  gate substitutions go to TL-Lite before final STOP/resume.

## Size Tripwire

| File | Target | Tripwire |
|---|---:|---:|
| `phase_brief.md` | 150 | 200 |
| `decisions_active.md` | 150 | 200 |
| `progress_current.md` | 100 | 133 |
| `tl_escalations_current.md` | 80 | 107 |
