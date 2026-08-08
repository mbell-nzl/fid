# Decisions - Active

> **Purpose.** Live carry-forward constraints and still-binding commitments.
> Loaded by default in every agent session. Keep this small.
>
> **Size target:** under 150 lines.

## Durable Architecture Decisions

The first three ADRs predate the phase system. ADR-0004 records the approved
fid/V2 product boundary, ADR-0005 resolves durable work identity, ADR-0006
defines parallel mutation isolation, and ADR-0007 defines interruption
recovery. All remain binding unless explicitly superseded:

- `docs/adr/0001-journal-canonical-db-derived.md` — journal is source of truth,
  SQLite is a derived projection. No daemon.
- `docs/adr/0002-configurable-roles-default-profile.md` — roles/permissions are
  configurable, seeded from a default profile.
- `docs/adr/0003-three-tier-qa-model.md` — validation contracts, QA plans, and
  QA reviews ship together.
- `docs/adr/0004-operational-work-state-boundary.md` — fid owns issue tracking
  plus durable operational work state/evidence; the external kit owns
  engineering-workflow policy, risk/authority/role routing, orchestration, and
  provider projections. ADR-0002/0003 remain compatibility assumptions pending
  explicit forward-roadmap disposition.
- `docs/adr/0005-issues-are-the-durable-work-graph.md` — fid Issues are the
  only durable work graph. Work Envelopes reference selected Issues; the
  harness derives noncanonical Next actions instead of persisting a second
  Work Item graph.
- `docs/adr/0006-isolate-parallel-builder-mutations.md` — one active mutating
  Builder Claim per Issue; parallel mutation requires distinct Issues,
  separate worktrees, and provably disjoint scopes.
- `docs/adr/0007-recover-interruptions-with-linked-runs.md` — a recovery
  creates a linked new Run only after executor termination and non-destructive
  worktree reconciliation; ambiguity stops.

**Gate C delivery decision (2026-08-08):** qualify all twelve V2 kit
backend-port capabilities before migrating a live project. Implement them as
bounded vertical slices with contract/conformance evidence, controlled
fixtures, and a disposable migration rehearsal; this is not a big-bang phase
or approval of the whole historical fid roadmap. Exact slicing follows the
open semantics decision and roadmap disposition.

**Gate B1 policy/enforcement decision (2026-08-08):** the kit evaluates
engineering-workflow policy and compiles the approved Work Envelope; fid does
not recalculate risk, routing, or obligation-selection policy. fid enforces
generic integrity plus the approved envelope's recorded preconditions,
authorization, non-weakening, trusted-evidence, freshness, and lifecycle
rules. Shared conformance fixtures prevent contract drift.

**Gate B2 trusted-executor decision (2026-08-08):** fid is the primary local
executor for registered gates and creates receipts from observed execution;
callers supply only immutable registration references. External attested
runners are a secondary CI/Assured/release extension. Reported results never
satisfy trusted obligations. Live output is required; framing remains an
interface-design detail.

**Gate B3 Gate Definition authority decision (2026-08-08):** the project owns
approved machine-readable Gate Definition values; the kit owns their versioned
schema, validation, and activation; fid owns immutable registrations,
executions, and receipts and never executes directly from mutable authoring
configuration. `QUALITY_PLAN.md` is explanatory. Exact authoring-file location
is deferred; standalone clients use the registration contract.

**Gate B4 fid-native work identity decision (2026-08-08):** fid Issues and
their relationships are the only durable work graph. A Work Envelope selects a
primary Issue and, where necessary, a bounded related Issue set; it does not
own persistent Work Items. Durable decomposition uses tasks/subtasks; smaller
steps remain transient plans or Run checkpoints. The harness computes Next
actions from durable state. Claims/Runs bind Issue, Envelope revision, actor
role, and optional Obligation; reviews use Review Runs/Records. Parallel
claims, interrupted-Run recovery, review rework, and multi-Issue Envelope
revision are resolved by Gate B4a/B4b/B4c/B4d. The Gate B4 lifecycle design is
complete; exact transport and projection schemas remain implementation work.

**Gate B4a parallel-Builder decision (2026-08-08):** one Issue may have at
most one active mutating Builder Claim. Claims bind Issue, Envelope revision,
Builder, and worktree; the Envelope itself is not claimed. Parallel mutation
requires distinct selected Issues, separate worktrees, and declared provably
disjoint scopes. Otherwise the second Next action is non-actionable with
`mutation-conflict`. Read-only Review Runs may coexist but are repository-
identity-bound and become stale when their target changes. Gate B4b later
resolves interrupted-Run recovery, Gate B4c resolves review rework, and Gate
B4d resolves multi-Issue Envelope revision/partial progress.

**Gate B4b interrupted-Run recovery decision (2026-08-08):** continuation
always creates a new Run linked to a terminally `interrupted` prior Run. The
old executor must be confirmed terminated and its Claim released, expired, or
authorizedly superseded; lease expiry alone is insufficient. The kit/harness
reconciles the worktree non-destructively. In-scope changes may be explicitly
adopted but are not trusted evidence; unsafe or ambiguous state returns
`recovery-state-conflict`, and unknown liveness returns
`executor-liveness-unknown`. Issue/Envelope state, obligations, evidence rules,
and consumed budgets carry forward. Gate B4c resolves review rework and Gate
B4d resolves multi-Issue Envelope revision/partial progress.

**Gate B4c review-rework decision (2026-08-08):** a completed Review Run
creates an immutable, exact-target Review Record. `changes-required` leaves the
review Obligation unsatisfied and derives `builder-rework` on the original
Issue/Envelope; no Work Item or follow-up Issue is created automatically. A
new Builder Claim/Run addresses in-scope Findings, but only a new independent
Review Record may resolve them. Target changes stale prior proof; initial reuse
is whole-profile conservative. Scope, protected-decision, dependency, risk,
acceptance, Obligation, or authority changes require an Envelope revision.
Blocking work cannot be deferred merely to pass, and exhausted cycles return
`review-rework-exhausted` plus Attention. Gate B4d resolves the remaining
multi-Issue scenario.

**Gate B4d multi-Issue Envelope decision (2026-08-08):** partial progress is a
projection of selected fid Issue states, not an Envelope state or membership
lifecycle. Selection or material contract changes create immutable revisions;
a new revision cannot activate while mutation remains live against the old
one. Removing unfinished work requires an authorized visible disposition, and
adding work selects an existing Issue. Completed Issues remain complete, but
later target changes may stale bound proof. Envelope completion requires every
Issue selected by the active revision to be satisfying-terminal and every
applicable Obligation current and satisfied; cancellation is not satisfaction.
Gate B4 lifecycle design is complete.

## Checkpoint Decisions (2026-07-18, pending Human ratification)

**DEC-059..DEC-066** — trusted-evidence checkpoint
(`../reviews/2026-07-18-roadmap-update-decisions.md`): registered gate execution
with receipts and trust sources; typed evidence + freshness; optional completion
prose behind structured `has_completion`; work envelopes + runner primitives;
extended economics; roadmap re-slice. Binding for all v0.2+ planning once
ratified. DEC-066 relabels: old v0.2b→v0.2d, old v0.2c→v0.2e.

## Carry-Forward Register

Reviewed at each phase start (schedule, re-defer, close, or cancel) and audited
at each wrap against what the phase actually changed.

| ID | Source | Item | Target | Priority |
|---|---|---|---|---|
| CF-008 | Phase 5 DEC-015 (human HUMAN-013c) | Resolution type / cancellation reason: `cancelled` is one undifferentiated bucket (no "redundant", "duplicate", "won't fix"). Needs a resolution-reason field + cancel `--reason`, touching the data model and `activity_log` audit. Re-deferred (DEC-082/093) to land with the governance entities and the DEC-081 ref-type extension. **Now paired with CF-026b** — both are governance-shaped data-model questions about vocabulary and lifecycle. | v0.2d | Low (product) |
| CF-017 | Phase 15 Review 1 + final review | Behavior-preserving function extraction for `internal/dev` `Finish` (~109 lines) and `Block` (~97 lines), both above the 60-line non-test function target. Scheduled, not trigger-based (DEC-073). Same class as CF-018 and CF-027. **Trigger did not fire in Phase 19** (`internal/dev` untouched). | v0.2e (`fid dev resume`/`abandon`) planning item | Low |
| CF-018 | Phase 16 Review 2 + final review | Behavior-preserving split of `internal/sqlquery/sqlquery.go` (595 lines) under the 400-line non-test file target. Keep the public package API stable. **Trigger did not fire in Phase 19** (`internal/sqlquery` untouched). | Next SQL package reopen | Low |
| CF-020 | Phase 17 DEC-056 | Burst issue writes stamp `occurred_at` ahead of the wall clock via `nextIssueUpdateTime`'s monotonicity bump, and `recent_activity_7d` excludes future rows. Bounded and self-correcting. | v0.2 | Low (product) |
| CF-022 | Phase 17 Review 1 | Broader exact-output markdown characterization tests to lock rendering behavior. **Trigger fired in Phase 19 Task 5 and was partially taken:** the new workflow-list renderer ships a byte-exact golden. The backfill for pre-existing renderers stays deferred — and should be taken **with CF-027's `writeMarkdown` split**, since both reopen `internal/format`. | Next `internal/format` reopen (pair with CF-027) | Low |
| CF-023 | Phase 17 DEC-058; expanded by DEC-074/DEC-095; raised to Medium by DEC-095 | **Workflow-tooling hardening.** Landed history and closed sub-items are in `../phases/phase-19/decisions.md`. **Open scope:** (a) `drill-09` plus the `report-lint.sh` evidence-token check (DEC-087), both needing a shell session to demonstrate failing-then-passing; (b) `report-lint.sh` should name `UNVERIFIED`/`FAIL` explicitly instead of the generic "no PASS/BLOCKED result" (DEC-103 — behaviour is correct, message is not); (c) **mechanical `Refs: phase-N/task-M` enforcement** via a `commit-msg` hook or a `report-lint.sh` check, covering Conventional-Commit subject format too — in `50533e2..HEAD` every Developer/QA-Lead task commit carries the trailer **except `32614bb` and `ce4cf5b`, both from the same irregular Task 4 session**, and the six other gaps are five `docs(tl)` governance commits plus one raw human commit (`ebbcd16`, which also has no Conventional-Commit type and names the wrong phase). Decide whether the hook covers `docs(tl)` or `AGENTS.md` exempts it explicitly; (d) **spec-sync obligations must appear in the owning task's `## Documentation Updates`**, not only in the phase sync table (DEC-106 item 5). **Five accumulated drill cases:** (1) a non-terminating or wrong-exit command recorded `PASS` — flag `rg`/`grep` results written as `exit 0 (no output)`; (2) a provider-shared prompt body asserting a capability limit only one provider's frontmatter enforces (DEC-097); (3) a required command failing on a *task-authored harness*, repaired and rerun byte-identically — correct handling is `PASS` **plus** a deviation note (DEC-100); (4) a report declining an out-of-scope repair as "pre-existing" without citing an introducing commit older than phase-start (DEC-102); (5) **a count presented as a rate without its denominator** (DEC-106 item 2) — the cheapest of the five to detect. **Record-keeping remains unfixed across two phases:** three of eight Phase 19 triage entries required Tech Lead backfill. | **Workflow overhaul**, or any session with a shell. **Decide at Phase 20 planning whether this is scheduled work or a wish list.** | **Medium** |
| CF-024 | Phase 18 DEC-073 | TRD-PERF-003 headroom. **Measured in Phase 19 Task 6: 100k-event rebuild `24.749 s` against the `<30 s` control (~17.5% headroom, up from ~8-10%)**, and 10k registry read `44.29 ms` against `<500 ms`. Improved but **not closed** — v0.2b receipts/evidence raise event volume on this path. Next lever is batching commits per N events in the shared `ProjectEvent` / `WriteProjectionMeta` path (Tech Lead-owned). **Do not weaken the target.** | v0.2b (first projection/rebuild work) | Medium |
| CF-026a | Phase 19 Human VAL-026 dogfood (DEC-105) | **Workflow inspection — the read half.** `fid workflow list` exposes only counts; no named-workflow `show`/`get` renders states, categories, ordered transitions, requirements, and assignments together. The shipped surface matches the documented surface exactly (`06_DESIGN_PLAN_CLI.md`), so this is a **new gap found by dogfooding**, not unshipped scope. Bounded read surface over already-projected tables; no data-model decision needed. | **v0.2a closure — Phase 20 planning candidate** (small, additive, no schema change) | Medium (product) |
| CF-026b | Phase 19 Human VAL-026 dogfood (DEC-105) | **Workflow lifecycle management — the design half.** Custom workflows cannot be revised, deactivated, retired, or deleted; reusing a slug is rejected. Needs a journal-safe model *decided before it is scheduled*: whether revision creates a successor/version, how kind/project assignment changes apply, what existing issues retain on revision, and which history-preserving retirement/deletion rules apply. No direct SQLite edits. | **Design decision first**, then v0.2d governance work alongside CF-008 | Medium (product) |
| CF-027 | Phase 19 QA-BLK-FINAL-05 (DEC-105); extended by DEC-106 (MINOR-01) | Behavior-preserving decomposition of **`internal/issues/create.go`** (412 lines, over STRUCT-002's 400; `Create` ~228 lines, over STRUCT-001's 60) and **`internal/format/markdown.go` `writeMarkdown`** (68 lines, over 60). **Both are inherited debt, not phase regressions** — `create.go` was 419 lines at phase start and Phase 19 made it *smaller*; `writeMarkdown` was 66 and grew by one unavoidable type-switch dispatch arm. Same class as CF-017 and CF-018. **Sequencing was the safety property and it is now satisfied:** CF-025's Go parity test landed in Phase 19 Task 6 and is the exact regression guard for the `create.go` path, with VAL-058/059 as the net. | **Now unblocked** — next `internal/issues` reopen; pair the `writeMarkdown` half with CF-022 | Low |

## Standing Rules From Phase 19

Eleven rule-bearing decisions; **five items retired at wrap** (DEC-106 item 6).
**Every surviving rule is now a sentence in a canonical spec**, so load the
spec, not the decision. Reasoning: `../phases/phase-19/decisions.md`.

- **Uncaptured gate results are `UNVERIFIED`, never `PASS`** (DEC-102 item 1;
  `09_AI_WORKFLOW.md` §8, `AGENTS.md`). Attestation from memory, transcription
  from another agent, and a promised future control are all `UNVERIFIED`.
- **"Pre-existing" is a provenance claim** (DEC-102 item 5; same sections).
  Cite the introducing commit and show it predates phase start; `git show HEAD`
  proves a line is committed, not inherited.
- **A required command is executed by its author before publication**, and
  working-directory sensitivity is declared (DEC-103 items 4-5; same sections).
- **Evidence is an automated test** unless the validation contract records why
  not (DEC-104 item 4; `09_AI_WORKFLOW.md` §6).
- **Cleanup is never evidence**; scratch state is unique per attempt; task specs
  declare required environment capabilities (DEC-101 items 1-3).
- **A task-authored harness that fails while being authored is a deviation, not
  a failed gate** (DEC-100), under its four conditions.
- **No-match inventories terminate `exit 1 (no output)`** (DEC-094).
- **Provider-shared prompt bodies must not assert capability limits only one
  provider's frontmatter enforces** (DEC-097; `09_AI_WORKFLOW.md` §7).
- **Second-escalation tripwire** (DEC-104 item 3): a second escalation on the
  same assertion, procedure, or artifact forces reconsideration of the
  **approach** — the disposition must say whether the mechanism should exist.
- **The Tech Lead has a shell** (DEC-104 item 1) and executes any command it
  authors before publishing it. It still never writes product code or tests, and
  never runs quality gates *as task evidence*.
- **Architecture:** `internal/workflow` owns authoring, assignment, and
  resolution and may use the journal write API (DEC-089/105); no in-package
  (`package db` / `package journal`) test file may import it (DEC-096).
  `02_TRD.md` §3/§3.1 is canonical.

**Standing lesson:** eight Phase 19 escalations, twelve decisions, **zero
product defects** — every stop was in verification apparatus. When an evidence
apparatus grows larger than the observable it protects, fix the apparatus, not
the brick that fell, and prefer expressing the observable as a test.

## Recently Closed

- **Phase 19 closed** (DEC-106, 2026-08-08) — v0.2a slice 2 of 3: workflow
  authoring is data (`fid workflow define/assign/list`, `workflow_assignments`
  with kind→project→built-in precedence, `escalation-default` seeded).
  VAL-001..070 all PASS; final review PASS WITH IMPROVEMENTS; checklist 9/9.
  **CF-021 and CF-025 closed** in Task 6; CF-026a/b and CF-027 opened.
- **Phase 18 closed** (DEC-073) — v0.2a slice 1 of 3: workflow structure is data
  (migrations `0007`/`0008`), every consumer resolves through `workflow.Load`,
  the frozen Go `definitions` map is test-only, the DEC-066 seam has zero
  behavior.
- **v0.1 closed** at Phase 17 wrap (DEC-057) — every v0.1-targeted PRD row maps
  to passing evidence, binary 12.27 MiB (<30 MB), accepted for early dogfood.
  **CF-019** closed there too.
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
- **Every spec-sync obligation appears in the owning task's
  `## Documentation Updates`**, not only in the phase sync table (DEC-106).
- Before release-readiness work, decide where held-out scenarios live and how
  agents are prevented from loading them.
- **No production Go path may read the frozen v0.1 workflow `definitions`
  map** — test-parity fixture only (Phase 18 VAL-024, Phase 19 VAL-070).
- **The gate-obligation seam stays behavior-free until v0.2b**: no seeded use,
  no receipt/gate evaluation or query for `required_gates_current` /
  `has_completion` (DEC-066; Phase 18 VAL-016/017, Phase 19 VAL-068/069).
- **A deep retrospective is recommended before Phase 20 planning** (DEC-104,
  DEC-106); brief in `../phases/phase-19/decisions.md`.

## TL-Lite Resume Policy

- Autonomous TL-Lite resume: **ENABLED (permanent)**. The warm-up safeguard is
  retired — see DEC-002 in `../phases/phase-0/decisions.md`.
- Phase Runner does not stop on `VERDICT: RESUME`; it records the
  `TRG-<phase>-<NNN>` and continues. `VERDICT: BLOCK` always stops the run.
- Per DEC-047, blocker reports, required gate failures, report-lint failures,
  missing/skipped/narrower command evidence, and ad hoc gate substitutions go to
  TL-Lite before final STOP/resume.

## Size Tripwire

| File | Target | Tripwire |
|---|---:|---:|
| `phase_brief.md` | 150 | 200 |
| `decisions_active.md` | 150 | 200 |
| `progress_current.md` | 100 | 133 |
| `tl_escalations_current.md` | 80 | 107 |
