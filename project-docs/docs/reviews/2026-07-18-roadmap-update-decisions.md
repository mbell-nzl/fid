# Decision Record - Trusted Evidence And Lean Workflow Checkpoint

Date: 2026-07-18
Author: Tech Lead (planning/specification checkpoint; no implementation)
Source brief: `2026-07-18-fid-roadmap-update-tl-brief.md`
Status: recorded, **pending Human ratification**; no product code changed.

This record continues the project DEC numbering after Phase 17 (DEC-058).
It adopts the brief's recommended defaults D1-D7. One recorded position is
deliberately superseded: PRD-QA-009's "fid is a verifier, not an executor"
note, which predates the Phase 17 / CF-023 evidence that agent-reported
deterministic claims are the workflow's weakest link. No other recorded
product or architecture decision contradicts the defaults; ADR-0001
(journal-canonical, no daemon), ADR-0002 (configurable roles), and ADR-0003
(three-tier QA) are all preserved.

---

## DEC-059: fid executes registered gates; reported exit codes are not trusted evidence

- **Date:** 2026-07-18
- **Context:** `fid dev gate` records a caller-supplied command, exit code,
  and excerpt (PRD-DEV-004/010). Nothing attests that the command ran. The
  brief's checkpoint evidence (Phase 17 wrap, CF-023, DEC-058) shows the
  workflow paying prose/protocol costs to compensate for untrustworthy
  deterministic claims.
- **Decision:** fid executes **registered gate definitions** — versioned
  entries in trusted project configuration with structured argv, explicit
  cwd, timeout, and environment allowlist — and records the outcome itself.
  Alternatively, fid accepts receipts from a separately trusted runner
  through a versioned attestation contract. An agent-reported exit code
  remains importable historical data (`trust_source=reported`) but can never
  satisfy a trusted-gate obligation. Boundary: fid does not run agents, host
  models, mutate git, evaluate shell strings, or become a general job
  scheduler; if execution lives in a companion process, fid still owns the
  receipt schema, validation, journal events, linkage, and state-transition
  enforcement.
- **Consequences:** Supersedes the "verifier, not executor" note on
  PRD-QA-009 for registered gates only. New PRD area `GATE`; TRD gains
  executor-boundary requirements (TRD-ARCH-016) and security requirements
  (TRD-SEC-005..009). `fid dev gate --cmd --exit-code` survives temporarily
  as an explicitly deprecated reported-evidence path.
- **Refs:** brief D1; PRD-GATE-001..008; TRD §4.7, §6; 06_DESIGN_PLAN_CLI
  §5a; CF-023.

## DEC-060: Execution attempts and terminal receipts supersede gate_runs as evidence

- **Date:** 2026-07-18
- **Context:** `gate_runs` rows are single append-only facts with no
  execution identity, no repository binding, and no freshness anchor.
- **Decision:** Model trusted gate evidence as **immutable execution
  attempts plus terminal receipts**. A successful receipt captures at least:
  receipt ID, schema version, producer, trust source, issue/attempt/agent-run
  linkage, role, timestamps; gate definition ID+version, structured argv,
  cwd, timeout, exit code or signal, duration; git branch/worktree plus
  base, head, and tree identity, dirty state, capture source; toolchain
  fingerprint; bounded sanitized stdout/stderr excerpts, full-log reference
  where retained, content hashes, byte counts, truncation/redaction flags;
  and terminal status `pass | fail | timed_out | cancelled | harness_error`.
  Because external execution and journal append cannot share one filesystem
  transaction, the model is honest: append an execution-start fact, append
  the terminal receipt only after the process terminates, and let
  `fid maintenance check` surface abandoned starts. **A crash may lose a
  pass; it must never fabricate one.**
- **Consequences:** New `gate_executions` / `gate_receipts` entities
  (03_TRD_SCHEMA §2.4 draft DDL); legacy `gate_runs` rows are never
  rewritten and project as `trust_source=reported`.
- **Refs:** brief D2; PRD-GATE-002..005; TRD §4.3/§4.7; DEC-059.

## DEC-061: Typed evidence records, obligations, and conservative freshness

- **Date:** 2026-07-18
- **Context:** Planned `fid evidence check` (PRD-QA-009) proves only that
  references resolve; command and manual evidence stay `unverifiable` and
  nothing establishes freshness.
- **Decision:** Add typed **evidence records** that link AC verifications,
  QA review coverage, validation outcomes, and workflow gates to receipts,
  commits, source state, screenshots, artifacts, or human outcomes. Rules:
  (a) deterministic obligations require trusted terminal receipts or another
  checkable evidence type — an unverifiable command string cannot satisfy
  them; (b) evidence **validity is separate from existence**: `current`,
  `stale`, `missing`, `unverifiable`, or `inapplicable`; (c) freshness starts
  conservative — any relevant target-tree change makes a gate receipt stale;
  path-aware invalidation waits for dogfood evidence plus a sound
  affected-path model and a new decision; (d) required gates are identified
  by gate ID/version, never display strings; (e) manual Human judgment
  remains valid first-class evidence when the obligation requires or allows
  a Human — record actor, checklist item, outcome, timestamp, optional note,
  and never mislabel it deterministic proof.
- **Consequences:** New PRD area `EVID`; `evidence_records`/`evidence_links`
  schema draft; PRD-QA-009 amended to add receipt freshness reporting once
  capture-time hashes exist.
- **Refs:** brief D3; PRD-EVID-001..006; 03_TRD_SCHEMA §2.4.

## DEC-062: Structured completion replaces mandatory completion prose

- **Date:** 2026-07-18
- **Context:** `has_completion_report` (PRD-WF-008) forces a Markdown
  report per attempt. Phase 17 spent five of eight triages and three
  escalations on completion-report commit-metadata friction (DEC-058,
  CF-023) — cost paid to prose that structured state should replace.
- **Decision:** Introduce a structured **`has_completion`** transition rule:
  satisfied by an append-only completion row whose applicable AC, concern,
  receipt, and evidence-obligation rules are met. Completion Markdown becomes
  optional through an additive workflow/profile change; an optional short
  authored summary/notes path is retained for explanation. Completion
  **views** are generated (JSON, TOON, Markdown, HTML) from changed
  commits/files and repository identity, gate receipts and evidence links,
  AC/validation coverage, telemetry/interventions/deviations/concerns/
  decisions, and the optional authored summary. Views link to append-only
  receipt rows; fid does not snapshot independent copies of mutable
  gate-result JSON.
- **Consequences:** `has_completion_report` retained for compatibility until
  a recorded flag-removal decision; PRD-WF-014, PRD-OUT-025 added; the
  `task_completions.gate_results` JSON snapshot is superseded by receipt
  links for trusted-mode completions.
- **Refs:** brief D4; DEC-058/CF-023; PRD-WF-008/014; 06_DESIGN_PLAN_CLI
  `fid dev finish`.

## DEC-063: Work envelopes and deterministic-runner primitives; runner is not an LLM obligation

- **Date:** 2026-07-18
- **Context:** PRD §7 carried the open question "should `workflow_runner`
  become a deterministic harness?" gated on a real failure. Phase 17's
  report-lint stop/resume churn (DEC-058) is that failure: orchestration
  turns burned on mechanical enforcement an LLM prompt was carrying.
- **Decision:** Promote the conditional task-pack concept (PRD-OPS-020) into
  a versioned, structured **work envelope**: epic/phase outcome, risk lane,
  approved authority, protected decisions, task graph, validation and gate
  obligations, budgets, and stop conditions. fid exposes **orchestration
  primitives**, not an LLM loop: next actionable work with typed skip
  reasons; claim/lease and resumable run state; retry/intervention budgets
  and circuit-breaker state; machine-readable transition rejection and stop
  reasons; attention items requiring Human or product/architecture judgment.
  A deterministic runner may remain a separate kit component driving these
  primitives. Lifecycle bookkeeping, evidence completeness, and stop-policy
  enforcement move out of runner prompts into fid/runner behavior. Humans
  retain authority over protected decisions.
- **Consequences:** Resolves the PRD §7 open question. PRD-OPS-020 is
  superseded by PRD-OPS-021/022 (envelope + primitives, v0.3); the
  build-pipeline workflow doc gains only a future-marked migration note
  (09_AI_WORKFLOW §11) and is otherwise unchanged until fid adoption.
- **Refs:** brief D5; PRD-OPS-020..022; GLOSSARY "Work envelope".

## DEC-064: Generated human validation brief and structured review coverage

- **Date:** 2026-07-18
- **Decision:** Extend attention/wrap-up reporting into a generated **human
  validation brief** containing only: outcomes delivered and user-visible
  changes; product/architecture/data/dependency/security decisions;
  deviations, unresolved risks, failed or carried-forward obligations, open
  escalations; receipt and validation roll-up plus spend, retries, and
  interventions; and the human-only checks still requiring judgment. QA
  reviews record structured coverage: scope, review profile, risk
  dimensions, files/contracts examined, receipts reproduced, verdict, and
  findings. A complete review with zero findings is valid. Security, code,
  validation, UX, data, and performance reviews are **profiles** over one
  common review contract, not separate persisted entity types.
- **Consequences:** PRD-OUT-024 (brief, v0.4 after validation outcomes
  exist; ahead of export/dashboard/sync), PRD-QA-010 (structured coverage,
  shipping with the first QA implementation so no string-reference
  migration is needed later).
- **Refs:** brief D6; PRD-OUT-020/024; PRD-QA-010.

## DEC-065: Trajectory economics beyond tokens

- **Date:** 2026-07-18
- **Decision:** Extend run telemetry beyond tokens/nominal cost to include
  wall duration, retry count, resumptions, context compactions where
  observable, tool or harness failures, normalized failure class, Human
  interventions and minutes, review findings, and escaped defects when
  later recorded. Capture automatically when the harness knows; keep
  unknown values **null** rather than asking an agent to estimate.
- **Consequences:** PRD-TELEM-007 added; additive `agent_runs` columns land
  with the v0.2 run-ops slice; stats surfaces extend accordingly.
- **Refs:** brief D7; PRD-TELEM-007.

## DEC-066: Roadmap re-slice — evidence before prose; compatibility for v0.1 data

- **Date:** 2026-07-18
- **Decision:** Re-slice the roadmap dependency-first (see
  `10_IMPLEMENTATION_ROADMAP.md`): keep **v0.2a Workflow Engine Migration**
  (plus versioned gate-obligation seams); insert **v0.2b Evidence
  Foundation** (gate definitions, executions/receipts, trust source, repo
  binding via `internal/gitctx` pulled forward, freshness, maintenance
  checks, capability discovery); then **v0.2c Lean Completion**, **v0.2d
  Governance Entities**, **v0.2e Reports/Context Packs/Run Ops/Economics**;
  **v0.3a Work Envelope & Runner Primitives** before the QA layer; QA and
  validation entities link first-class evidence from their first
  implementation; the human validation brief lands ahead of one-way export,
  dashboard, MCP evaluation, and bidirectional sync. Compatibility: never
  rewrite historical events; legacy gate rows project as
  `trust_source=reported`; new journal events and projection migrations are
  versioned; old and trusted gate modes are visibly distinct; completion
  Markdown becomes optional additively before any flag removal; receipts
  and capability discovery have a JSON/TOON contract so a workflow kit
  integrates without touching SQLite or `.fid/` internals; integrations
  detect missing capabilities and stop — never silently falling back from
  trusted receipts to reported evidence.
- **Consequences:** First proposed implementation phase: **Phase 18 =
  v0.2a Workflow Engine Migration (heavy mode)**, preceded by the CF-023
  workflow-tooling pass. Dashboard/MCP/sync remain v0.5 conditionals.
- **Refs:** brief "Roadmap Re-Slicing Guidance"; 10_IMPLEMENTATION_ROADMAP;
  phase_brief.

---

## Spec-sync ledger for this checkpoint

| Document | Change |
|---|---|
| `00_PROJECT_CHARTER.md` | Scope + concepts + durable decisions rows for registered gate execution, receipts, work envelope. |
| `GLOSSARY.md` | New: Gate definition, Execution attempt, Execution receipt, Trust source, Evidence obligation, Evidence freshness, Work envelope, Human validation brief. Amended: Gate, Evidence check, Workflow Runner. |
| `01_PRD.md` | New areas GATE/EVID; PRD-GATE-001..008, PRD-EVID-001..006, PRD-WF-014, PRD-QA-010, PRD-TELEM-007, PRD-OUT-024/025, PRD-OPS-021..023; amendments to PRD-DEV-004/010, PRD-WF-008, PRD-QA-009, PRD-OUT-020/021, PRD-OPS-020; traceability + open questions updated. |
| `02_TRD.md` | Executor/receipt architecture (§2/§3/§4.1/§4.3/§4.7), security TRD-SEC-005..009, interfaces (receipt + capability contracts), TRD-ARCH-016..018, open questions. |
| `03_TRD_SCHEMA.md` | Draft v0.2 DDL: `gate_executions`, `gate_receipts`, `evidence_records`, `evidence_links`; prefixes; permissions additions. |
| `04_TRD_CONFIG.md` | `[gates]` registry file, artifacts/retention layout. |
| `06_DESIGN_PLAN_CLI.md` | `fid gate` family, `fid capabilities`, deprecation of reported `fid dev gate`, lean `fid dev finish`, completion view, validation brief, envelope family sketch, evidence check amendment. |
| `07_DESIGN_PLAN_FLOWS.md` | Trusted-gate developer flow. |
| `08_QUALITY_PLAN.md` | §5.10 evidence-foundation validation scenarios (10 scenarios). |
| `09_AI_WORKFLOW.md` | §11 future-marked enforcement-migration note only (per recorded feedback the build pipeline stays fid-free). |
| `10_IMPLEMENTATION_ROADMAP.md` | Full re-slice per DEC-066. |
| `docs/current/*` | Phase brief, progress, decisions_active refreshed. |
