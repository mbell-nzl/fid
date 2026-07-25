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

## Checkpoint Decisions (2026-07-18, pending Human ratification)

**DEC-059..DEC-066** — trusted-evidence checkpoint
(`../reviews/2026-07-18-roadmap-update-decisions.md`): registered gate
execution with receipts and trust sources; typed evidence + freshness;
optional completion prose behind structured `has_completion`; work envelopes
+ runner primitives; extended economics; roadmap re-slice (evidence
foundation v0.2b before governance/reports; envelope v0.3a). Binding for all
v0.2+ planning once ratified. DEC-066 relabels: old v0.2b→v0.2d,
old v0.2c→v0.2e (CF-017's target is now v0.2e).

## Carry-Forward Register

Tech Lead reviews this register at each phase start and schedules, re-defers,
closes, or cancels each item.

| ID | Source | Item | Target | Priority |
|---|---|---|---|---|
| CF-008 | Phase 5 DEC-015 (human HUMAN-013c) | Resolution type / cancellation reason: `cancelled` is a single undifferentiated bucket and cannot distinguish "redundant", "duplicate", "won't fix", etc. The enhancement is a resolution-reason field + cancel `--reason` flag, touching the data model and `activity_log` audit. **Dispositioned at Phase 19 planning (DEC-082): re-deferred to v0.2d**, where the governance entities and the `activity_log` ref-type extension (DEC-081) land together; it is orthogonal to workflow authoring. Human confirmation pending. | v0.2d (DEC-082) | Low (product) |
| CF-017 | Phase 15 Review 1 + final review | Behavior-preserving function extraction for `internal/dev` `Finish` (~109 lines) and `Block` (~97 lines), both above the 60-line non-test function target. Trigger fired in Phase 18 (Task 5 touched `internal/dev`) and was correctly not taken during a behavior-preserving migration; now **scheduled**, not trigger-based (DEC-073). | v0.2e (`fid dev resume`/`abandon`) planning item | Low |
| CF-018 | Phase 16 Review 2 + final review | Behavior-preserving split of `internal/sqlquery/sqlquery.go` (595 lines) under the 400-line non-test file target. Keep the public package API stable. | Next SQL package reopen | Low |
| CF-020 | Phase 17 DEC-056 | Burst issue writes stamp `occurred_at` ahead of the wall clock via `nextIssueUpdateTime`'s monotonicity bump, and `recent_activity_7d` excludes future rows. Bounded and self-correcting; do not fix inside a v0.1 closure audit. | v0.2 | Low (product) |
| CF-021 | Phase 17 Review 3 / STRUCT-003 | Add Go doc comments to exported `maintenance.VacuumResult` and `Vacuum`. **Scheduled at Phase 19 planning (DEC-082): Phase 19 Task 6**, which already touches `internal/maintenance` and runs the full Go gate. | **Phase 19 Task 6 (scheduled)** | Low |
| CF-022 | Phase 17 Review 1 | Add broader exact-output markdown characterization tests to lock rendering behavior beyond the current suite. Trigger fires in Phase 19 Task 5 (new `internal/format` markdown renderer). **Partially taken (DEC-082):** Task 5 ships a byte-exact golden for the new workflow-list renderer; the backfill for pre-existing renderers stays deferred. | Remainder: next `internal/format` markdown reopen | Low |
| CF-023 | Phase 17 DEC-058; expanded by DEC-074; **narrowed by DEC-088** | Workflow-tooling hardening. **Landed 2026-07-26** by the retrospective pass: (d) ID uniqueness → DEC-084 phase-scoped IDs; (e) mint ID + triage row on `BLOCK` → DEC-085; (f) run-ledger chronological append + correction entries; (g) wrap-time reopen audit; (h) evidence honesty → DEC-086. **Already satisfied** before the pass: (a) worked commit-metadata example, (b) report-is-its-own-`docs(task)`-commit rule — the Phase 17 recurrence did not repeat in Phase 18. **Reclassified low-value:** (c) report-lint dirty-vs-parse messaging. **Remaining scope: DEC-087 only** — `drill-09` for the non-terminating-command defect class plus the `report-lint.sh` evidence-token check it needs, both requiring a shell session to demonstrate failing-then-passing. | **Workflow overhaul**, or any session with a shell (DEC-087) | Low |
| CF-024 | Phase 18 DEC-073 | TRD-PERF-003 headroom is thin: after the DEC-071 rebuild-connection fix the 100k-event rebuild measures `26.99–27.66 s` against the `<30 s` control (~8–10% margin). v0.2b receipts/evidence raise event volume on this path. Next lever is batching commits per N events in the shared `ProjectEvent` / `WriteProjectionMeta` path (Tech Lead-owned). **Do not weaken the target.** Phase 19 Task 6 carries a measurement obligation (re-record rebuild + read against the controls); the fix itself stays v0.2b (DEC-082). | v0.2b (first projection/rebuild work); measured in Phase 19 Task 6 | Medium |

## Phase 19 Planning Decisions (2026-07-26, pending Human approval)

**DEC-075..DEC-082** — recorded in `../phases/phase-19/decisions.md`: heavy
operating mode; `workflow_assignments` as a journal-written projected table
with kind→project→built-in resolution precedence; derived cancel edges and
initial-state resolution; define restricted to the four implemented rules;
custom state names not transitionable until Phase 20; `--applies-to`
vocabulary incl. `escalation`; no `activity_log` rows for workflow authoring;
carry-forward disposition. Binding for Phase 19 once approved.

## Workflow Retrospective Decisions (2026-07-26)

**DEC-084..DEC-088** — recorded in
`../reviews/2026-07-26-workflow-retrospective.md`. Binding now:

- **DEC-084** — triage/escalation IDs are **phase-scoped**
  (`TRG-<phase>-<NNN>`, `ESC-<phase>-<NNN>`), allocated from the phase triage
  log as the single source. Historical IDs are not renumbered.
- **DEC-085** — TL-Lite writes the TRG (and on `BLOCK` the ESC) entry **before**
  reporting its verdict; a placeholder ID is a defective triage. Phase Runner
  records the concrete IDs in the ledger's `Triage log:` field.
- **DEC-086** — every required command's result carries a verbatim terminal
  token; **a command that did not terminate is never `PASS`** (mode-invariant).
  `VAL-NNN` coverage cites a test, `file:line`, or output — not a restatement.
- **DEC-087** — `drill-09` + the `report-lint` termination check are specified
  but **not landed** (no shell in the authoring session); sole remaining CF-023
  scope.
- **DEC-088** — CF-023 narrowed to DEC-087; see the register row.

**Unverified:** the retrospective pass could not execute `sync-agents.sh`,
`sync-skills.sh`, `workflow-lint.sh`, `report-lint.sh --self-test`, or
`run-drills.sh`. Run them before trusting the changed surface.

## Recently Closed

- **Phase 18 closed** at wrap (DEC-073) — v0.2a slice 1 of 3 delivered:
  workflow structure is data (migrations `0007`/`0008`), every consumer resolves
  through `workflow.Load`, the frozen Go `definitions` map is test-only, and the
  DEC-066 gate-obligation seam is present with zero behavior. VAL-001..024 all
  PASS; final review PASS; human checklist PASS.
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
- **CF-023 is narrowed to DEC-087** (drill-09 + report-lint termination check).
  Its two former stand-in rules are now standing rules, not phase-local
  mitigations: DEC-086 (a non-terminating required command is never `PASS`) and
  DEC-084/085 (phase-scoped IDs, written before the verdict).
- **No production (non-test) Go path may read the frozen v0.1 workflow
  `definitions` map** — it is a test-parity fixture only; the projected workflow
  tables are the runtime source of truth (Phase 18 VAL-024).
- **The gate-obligation seam stays behavior-free until v0.2b**: no seeded use,
  no receipt/gate evaluation or query for `required_gates_current` /
  `has_completion` (DEC-066, Phase 18 VAL-016/017).

## TL-Lite Resume Policy

- Autonomous TL-Lite resume: **ENABLED (permanent)**. The warm-up safeguard is
  retired — see DEC-002 in `../phases/phase-0/decisions.md`.
- Phase Runner does not stop on `VERDICT: RESUME`; it records the
  `TRG-<phase>-<NNN>` and continues. `VERDICT: BLOCK` always stops the run.
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
