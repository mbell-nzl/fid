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
- `docs/adr/0002-configurable-roles-default-profile.md` — roles/permissions are
  configurable, seeded from a default profile.
- `docs/adr/0003-three-tier-qa-model.md` — validation contracts, QA plans, and
  QA reviews ship together.

## Checkpoint Decisions (2026-07-18, pending Human ratification)

**DEC-059..DEC-066** — trusted-evidence checkpoint
(`../reviews/2026-07-18-roadmap-update-decisions.md`): registered gate execution
with receipts and trust sources; typed evidence + freshness; optional completion
prose behind structured `has_completion`; work envelopes + runner primitives;
extended economics; roadmap re-slice (evidence foundation v0.2b before
governance/reports; envelope v0.3a). Binding for all v0.2+ planning once
ratified. DEC-066 relabels: old v0.2b→v0.2d, old v0.2c→v0.2e.

## Carry-Forward Register

Reviewed at each phase start: schedule, re-defer, close, or cancel each item.

| ID | Source | Item | Target | Priority |
|---|---|---|---|---|
| CF-008 | Phase 5 DEC-015 (human HUMAN-013c) | Resolution type / cancellation reason: `cancelled` is one undifferentiated bucket (no "redundant", "duplicate", "won't fix"). Needs a resolution-reason field + cancel `--reason`, touching the data model and `activity_log` audit. **Re-deferred to v0.2d (DEC-082), human-confirmed (DEC-093)** — it lands with the governance entities and the DEC-081 ref-type extension; orthogonal to workflow authoring. | v0.2d (DEC-082) | Low (product) |
| CF-017 | Phase 15 Review 1 + final review | Behavior-preserving function extraction for `internal/dev` `Finish` (~109 lines) and `Block` (~97 lines), both above the 60-line non-test function target. Trigger fired in Phase 18 (Task 5 touched `internal/dev`) and was correctly not taken during a behavior-preserving migration; now **scheduled**, not trigger-based (DEC-073). | v0.2e (`fid dev resume`/`abandon`) planning item | Low |
| CF-018 | Phase 16 Review 2 + final review | Behavior-preserving split of `internal/sqlquery/sqlquery.go` (595 lines) under the 400-line non-test file target. Keep the public package API stable. | Next SQL package reopen | Low |
| CF-020 | Phase 17 DEC-056 | Burst issue writes stamp `occurred_at` ahead of the wall clock via `nextIssueUpdateTime`'s monotonicity bump, and `recent_activity_7d` excludes future rows. Bounded and self-correcting; do not fix inside a v0.1 closure audit. | v0.2 | Low (product) |
| CF-021 | Phase 17 Review 3 / STRUCT-003 | Add Go doc comments to exported `maintenance.VacuumResult` and `Vacuum`. **Scheduled at Phase 19 planning (DEC-082): Phase 19 Task 6**, which already touches `internal/maintenance` and runs the full Go gate. | **Phase 19 Task 6 (scheduled)** | Low |
| CF-022 | Phase 17 Review 1 | Add broader exact-output markdown characterization tests to lock rendering behavior beyond the current suite. Trigger fires in Phase 19 Task 5 (new `internal/format` markdown renderer). **Partially taken (DEC-082):** Task 5 ships a byte-exact golden for the new workflow-list renderer; the backfill for pre-existing renderers stays deferred. | Remainder: next `internal/format` markdown reopen | Low |
| CF-023 | Phase 17 DEC-058; expanded by DEC-074; narrowed by DEC-088; raised to Medium by DEC-095 | Workflow-tooling hardening. Landed history and the closed sub-items are in `../phases/phase-19/decisions.md` (DEC-084..088, DEC-095, DEC-097, DEC-100). **Remaining scope: DEC-087** — `drill-09` plus the `report-lint.sh` evidence-token check, both needing a shell session to demonstrate failing-then-passing. **Three drill cases accumulated across Phase 19:** (1) a non-terminating or wrong-exit command recorded `PASS` — flag `rg`/`grep` results written as `exit 0 (no output)`; (2) a provider-shared prompt body asserting a capability limit only one provider's frontmatter enforces (DEC-097) — cheaply detectable mechanically; (3) a required command failing on a *task-authored harness*, repaired and rerun byte-identically (DEC-100) — correct handling is `PASS` **plus** a deviation note, not a blocker report and not a silent rerun. | **Workflow overhaul**, or any session with a shell (DEC-087) | **Medium** |
| CF-024 | Phase 18 DEC-073 | TRD-PERF-003 headroom is thin: after the DEC-071 rebuild-connection fix the 100k-event rebuild measures `26.99–27.66 s` against the `<30 s` control (~8–10% margin). v0.2b receipts/evidence raise event volume on this path. Next lever is batching commits per N events in the shared `ProjectEvent` / `WriteProjectionMeta` path (Tech Lead-owned). **Do not weaken the target.** Phase 19 Task 6 carries a measurement obligation (re-record rebuild + read against the controls); the fix itself stays v0.2b (DEC-082). | v0.2b (first projection/rebuild work); measured in Phase 19 Task 6 | Medium |

## Phase 19 Decisions (2026-07-26, approved DEC-093)

**DEC-075..DEC-083** — recorded in `../phases/phase-19/decisions.md`: heavy
mode; `workflow_assignments` as a journal-written projected table with
kind→project→built-in precedence; derived cancel edges and initial state;
define restricted to the four implemented rules; custom state names not
transitionable until Phase 20; `--applies-to` incl. `escalation`; no
`activity_log` rows for workflow authoring; carry-forward disposition.

**DEC-089..DEC-093** — QA-plan blocker dispositions (QA-BLK-01..07), same
file. Binding beyond Phase 19:

- **DEC-089** — `internal/workflow` owns workflow authoring, assignment, and
  resolution; `02_TRD.md` §3 synchronized (v0.1's "read-only registry / no
  authoring" wording retired).
- **DEC-090** — PRD-WF-013 stays unamended: an issue-applicable workflow
  (`epic`/`task`/`subtask`/`any`) must declare **exactly one**
  `cancelled`-category state, so no issue can become non-cancellable. Zero
  stays legal only for `escalation` / `qa_review`.
- **DEC-091/092/093** (phase-local) — contract refined to 70 `Ready` assertions;
  `05_DESIGN_PLAN.md` §3.1 gains the `workflow` command family; Human approval
  of mode, task list, contract, and CF-008.

**DEC-094..DEC-095** (2026-07-29, resolving ESC-19-001 / TRG-19-001), same
file. Binding beyond Phase 19:

- **DEC-094** — a search/inventory command proving an *absence* terminates
  with **exit 1 and no output**; record `PASS (no match) - exit 1 (no output)`.
  `exit 0 (no output)` is reserved for silent success and is fabricated
  evidence for `rg`/`grep`. `AGENTS.md` and `09_AI_WORKFLOW.md` §11 updated.
  Task specs and QA plans relying on such an inventory state its expected exit
  status.
- **DEC-095** — the Tech Lead **seeds `phases/phase-N/triage-log.md` at phase
  start** so ID allocation never depends on an agent creating a file; a verdict
  without a durable record is defective triage, and the Phase Runner stops
  rather than accept a placeholder ID.

**DEC-096..DEC-097** (2026-07-29, resolving ESC-19-002 / TRG-19-002), same
file. Binding beyond Phase 19:

- **DEC-096** — `internal/workflow` may depend on `internal/journal`, **bounded
  to the `Effect` data contract** (no append, projection, or replay surface);
  the authoring core's `[]journal.Effect` return type stands. Because
  `internal/journal` imports `internal/db`, **no in-package (`package db`) test
  file in `internal/db` may import `internal/workflow`** — such tests belong in
  the external `package db_test`. The same rule applies to `internal/journal`
  in-package tests. `02_TRD.md` §3 records the edge.
- **DEC-097** — TL-Lite's missing triage records are a **provider-prompt
  defect on Codex**, not agent unreliability: the *shared* body asserted a
  "cannot run commands" limit enforced only by Claude's frontmatter, which on
  Codex reads as "cannot write files." `project-docs/agents/tl-lite.md` fixed
  and regenerated (`sync-agents.sh` run and verified). **Standing lesson:
  provider-shared agent bodies must not assert capability limits that only one
  provider's frontmatter enforces.**

**DEC-099..DEC-100** (2026-08-01, resolving ESC-19-004 / TRG-19-004), same
file. Binding beyond Phase 19:

- **DEC-099** — Task 4 resumes from **preserved** partial work. The VAL-059
  guard is **structural, not temporal**: the baseline is built in a detached
  worktree at the pinned commit from a committed fixed fixture, so working-tree
  edits cannot affect it and capture order is not binding. **Provenance of a
  captured artifact is settled by reproducing it from its pinned inputs and
  comparing digests**, not by chronology or narrative.
- **DEC-100** — a required command whose **harness is a deliverable of the task
  being authored** (fixture, capture/measurement script, golden generator) may
  fail, be repaired, and be rerun to a pass, recorded `PASS` **with a mandatory
  deviation note** giving the first failure's token, root cause, and fix. All
  four must hold: the artifact is a task deliverable, **not** the product under
  test, a behavior-asserting test, or a pre-existing gate/tool; the rerun is
  **byte-identical**; the pass is **deterministic and third-party reproducible
  from committed artifacts**; both failure and final result are recorded.
  Otherwise the gate is failed. At triage, **"task incomplete, stop"** and
  **"evidence tainted, restart"** are different verdicts — a harness bug is
  normally the first. `AGENTS.md` and `09_AI_WORKFLOW.md` §11 updated.

## Workflow Retrospective Decisions (2026-07-26)

**DEC-084..DEC-088** — recorded in
`../reviews/2026-07-26-workflow-retrospective.md`. Binding now:

- **DEC-084** — triage/escalation IDs are **phase-scoped**
  (`TRG-<phase>-<NNN>`, `ESC-<phase>-<NNN>`), allocated from the phase triage
  log. Historical IDs are not renumbered.
- **DEC-085** — TL-Lite writes the TRG (and on `BLOCK` the ESC) entry **before**
  reporting its verdict; a placeholder ID is a defective triage. Phase Runner
  records the concrete IDs in the ledger's `Triage log:` field.
- **DEC-086** — every required command's result carries a verbatim terminal
  token; **a command that did not terminate is never `PASS`** (mode-invariant).
  `VAL-NNN` coverage cites a test, `file:line`, or output — not a restatement.
- **DEC-087 / DEC-088** — `drill-09` + the `report-lint` termination check are
  specified but **not landed**; sole remaining CF-023 scope (see the register).

**Partly unverified:** `sync-agents.sh` has since been run and verified (DEC-097
fix confirmed in the generated prompts). `sync-skills.sh`, `workflow-lint.sh`,
`report-lint.sh --self-test`, and `run-drills.sh` remain unexecuted — run them
before trusting that surface.

## Recently Closed

- **Phase 18 closed** (DEC-073) — v0.2a slice 1 of 3: workflow structure is data
  (migrations `0007`/`0008`), every consumer resolves through `workflow.Load`,
  the frozen Go `definitions` map is test-only, the DEC-066 seam has zero
  behavior. VAL-001..024, final review, and human checklist all PASS.
- **v0.1 closed** at Phase 17 wrap (DEC-057) — every v0.1-targeted PRD row maps
  to passing evidence (VAL-024), static binary 12.27 MiB (<30 MB), Human
  accepted for early dogfood. **CF-019** closed there too.
- Earlier closed carry-forward items are archived in their phase decisions:
  CF-011/016 (Phase 16), CF-013/014/015 (Phase 15), CF-003 (Phase 10), CF-012
  (Phase 8), CF-009/010 (Phase 6), CF-004/006 (Phase 4), CF-005/007 (Phase 5),
  CF-002 (Phase 0).

## Environment Constraints

- Go, module, and golangci-lint caches are redirected into gitignored in-repo
  dirs (`.gocache/`, `.gomodcache/`, `.gopath/`, `.golangci-cache/`, `.cache/`)
  via `[shell_environment_policy.set]` in `.codex/config.toml` so quality gates
  work under the workspace-write sandbox (DEC-001,
  `../phases/phase-0/decisions.md`). The policy applies at session start, so a
  fresh agent session is required for changes to take effect.

## Scheduled Commitments

- Every planned normal/heavy phase must create its
  `phases/phase-N/qa/validation-contract.md` with `VAL-NNN` assertions before
  task execution begins.
- **Every phase seeds an empty `phases/phase-N/triage-log.md` at planning
  time** (DEC-095), so TRG/ESC allocation never depends on a mid-phase agent
  creating the file.
- Before release-readiness work, decide where held-out scenarios live and how
  agents are prevented from loading them.
- **No production (non-test) Go path may read the frozen v0.1 workflow
  `definitions` map** — test-parity fixture only; the projected workflow tables
  are the runtime source of truth (Phase 18 VAL-024; Phase 19 VAL-070).
- **The gate-obligation seam stays behavior-free until v0.2b**: no seeded use,
  no receipt/gate evaluation or query for `required_gates_current` /
  `has_completion` (DEC-066; Phase 18 VAL-016/017, Phase 19 VAL-068/069).

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
