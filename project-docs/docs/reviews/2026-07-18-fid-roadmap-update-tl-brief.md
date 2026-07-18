# fid Roadmap Update Brief - Trusted Evidence And Lean Workflow Support

Date: 2026-07-18

Audience: `fid` Tech Lead

Checkpoint: v0.1 is closed; perform this review before planning v0.2a

Mode: planning and specification update only; do not implement product code

## Assignment

Revisit the `fid` product direction at the v0.1 wrap point and update the
canonical specifications and roadmap so `fid` can serve as the structured
state and trusted-evidence control plane for a slimmer agentic engineering
workflow.

Preserve `fid`'s local-first, journal-canonical, schema-enforced design. The
goal is not to turn it into a model host, daemon, CI service, or general shell
automation framework. The goal is to stop relying on agents to author prose
that claims deterministic commands ran, and instead capture trustworthy facts
at the point they occur.

Unless the Human overrides them when invoking this brief, use the recommended
defaults below as the planning direction. Stop after the specification,
decision, and roadmap update. Do not begin a v0.2 implementation phase.

## Required Context

Read these sources before editing:

- `AGENTS.md` and `project-docs/docs/current/phase_brief.md`.
- `project-docs/specs/00_PROJECT_CHARTER.md`, especially the design principles,
  scope, interfaces, and success criteria.
- `project-docs/specs/GLOSSARY.md`.
- `project-docs/specs/01_PRD.md`, especially `PRD-DEV-004`, `PRD-DEV-010`,
  `PRD-QA-009`, `PRD-TELEM-*`, `PRD-OUT-020` through `PRD-OUT-023`,
  `PRD-OPS-017` through `PRD-OPS-020`, and the open product questions.
- `project-docs/specs/02_TRD.md` section 4.7 and the observability, security,
  persistence, and compatibility sections.
- `project-docs/specs/03_TRD_SCHEMA.md` definitions for `gate_runs`,
  `task_completions`, `agent_runs`, QA entities, and validation outcomes.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` definitions for `fid dev gate`,
  `fid dev finish`, run telemetry, handoff, evidence, QA, validation, reports,
  and stats.
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` for the intended end-to-end
  agent workflow.
- `project-docs/specs/08_QUALITY_PLAN.md` and
  `project-docs/specs/09_AI_WORKFLOW.md`.
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md`, including all v0.2-v0.5
  sequencing and conditional features.
- Phase 17 wrap artifacts and `CF-023`, which record completion-report and
  report-lint protocol friction.

## Why This Checkpoint Exists

The v0.1 core provides the correct foundation: append-only journal events,
structured task attempts and completions, gate rows, claims, AC verification,
telemetry, and rebuildable projections. The current contract nevertheless has
three material limitations for a receipt-oriented workflow:

1. `fid dev gate` records a caller-supplied command, exit code, and excerpt; it
   does not execute or independently attest the command.
2. The planned `fid evidence check` proves only that certain references still
   resolve. Command and manual evidence remain `unverifiable`, and no
   capture-time content or repository-state hash establishes freshness.
3. `fid dev finish` and workflow transition rules still require a completion
   Markdown report, preserving the prose generation, linting, self-referential
   commit, and synchronization costs that structured state should remove.

The PRD says to revisit whether `workflow_runner` should become a deterministic
harness when runner cost or prompt drift causes a real failure. Phase 17 and
`CF-023` provide sufficient evidence to revisit that question now.

## Recommended Product Decisions

Record accepted decisions in the canonical decision mechanism and update all
affected specs. If the Human does not approve the defaults, present a compact
decision slate before editing rather than guessing.

### D1. Registered Gate Execution

Recommended default: `fid` executes registered deterministic gates, or accepts
receipts from a separately trusted runner through an equivalent versioned
attestation contract. An agent-reported exit code remains importable historical
data but does not satisfy a trusted-gate obligation.

Use a narrow boundary:

- Execute only versioned gate definitions approved in project configuration.
- Prefer structured argv arrays, explicit working directories, timeouts, and
  environment allowlists over shell command strings.
- Do not run agents, host models, mutate git, or become a general job scheduler.
- If execution lives in a companion process, `fid` still owns receipt schema,
  validation, journal events, linkage, and state-transition enforcement.

### D2. Receipts Are First-Class Evidence

Recommended default: evolve `gate_runs` into, or supersede it with, immutable
execution attempts and terminal receipts. A successful receipt should capture
at least:

- Receipt ID, schema version, producer, trust source, issue, attempt, agent run,
  role, and timestamps.
- Gate definition ID and version, structured argv, cwd, timeout, exit code or
  signal, and duration.
- Git branch/worktree plus base, head, and tree identity; dirty state and the
  capture source.
- Relevant toolchain fingerprint.
- Bounded sanitized stdout/stderr excerpts, full-log reference where retained,
  content hashes, byte counts, and truncation/redaction flags.
- Terminal status such as pass, fail, timed_out, cancelled, or harness_error.

External command execution and journal append cannot be one filesystem
transaction. Model this honestly: append an execution-start fact, append a
terminal receipt only after the process terminates, and let maintenance checks
surface abandoned starts. A crash may lose a pass, but must never fabricate one.

### D3. Evidence Registry, Obligations, And Freshness

Recommended default: add typed evidence records or references that can link AC
verifications, QA review coverage, validation outcomes, and workflow gates to
receipts, commits, source state, screenshots, artifacts, or human outcomes.

Define these rules:

- Deterministic obligations require trusted, terminal receipts or another
  checkable evidence type; an unverifiable command string cannot satisfy them.
- Evidence validity is separate from existence. It can be current, stale,
  missing, unverifiable, or inapplicable.
- Begin conservatively: any relevant target-tree change makes a gate receipt
  stale. Add path-aware invalidation only after dogfooding proves the coarse
  rule too expensive and a sound affected-path model exists.
- Required gates are identified by gate ID/version, not fragile display strings.
- Manual product judgment remains valid first-class evidence when the obligation
  explicitly requires a Human. Record actor, checklist item, outcome, timestamp,
  and optional note; do not mislabel it deterministic proof.

### D4. Completion Prose Becomes Optional

Recommended default: replace `has_completion_report` with structured
`has_completion` plus applicable AC, concern, receipt, and evidence-obligation
rules. Retain an optional short authored summary or notes path for explanation.

Generate completion views in JSON, TOON, Markdown, or HTML from:

- Changed commits/files and repository identity.
- Gate receipts and evidence links.
- AC/validation coverage and outcomes.
- Agent telemetry, interventions, deviations, concerns, and decisions.
- The optional authored summary.

Do not snapshot independent copies of mutable gate-result JSON when stable links
to append-only receipt rows can provide the same projection without drift.

### D5. Work Envelope And Deterministic Runner State

Recommended default: promote the conditional task-pack concept into a versioned,
structured work-envelope capability. It should capture an epic or phase outcome,
risk lane, approved authority, protected decisions, task graph, validation and
gate obligations, budgets, and stop conditions.

`fid` should expose orchestration primitives, not host an LLM loop:

- Next actionable work with typed skip reasons.
- Claim/lease and resumable run state.
- Retry/intervention budget and circuit-breaker state.
- Machine-readable transition rejection and stop reason.
- Attention items that require Human or architecture/product judgment.

A deterministic runner may remain a separate kit component that drives these
primitives. Move lifecycle bookkeeping, evidence completeness, and stop-policy
enforcement out of the runner prompt.

### D6. Generated Human Brief And Structured Review Coverage

Recommended default: extend attention and wrap-up reporting into a generated
human validation brief containing only:

- Outcomes delivered and user-visible changes.
- Product, architecture, data, dependency, and security decisions.
- Deviations, unresolved risks, failed/carried-forward obligations, and open
  escalations.
- Receipt and validation roll-up, spend, retries, and interventions.
- Human-only checks still requiring judgment.

QA reviews should record structured scope, review profile, risk dimensions,
files/contracts examined, receipts reproduced, verdict, and findings. A complete
review with zero findings is valid. Security, code, validation, UX, data, and
performance are review profiles or skills over a common review contract, not
necessarily separate persisted entity types.

### D7. Broader Trajectory Economics

Recommended default: extend telemetry beyond tokens and nominal cost to include
wall duration, retry count, resumptions, context compactions where observable,
tool or harness failures, normalized failure class, Human interventions and
minutes, review findings, and escaped defects when later recorded.

Capture automatically when the harness knows. Keep unknown values null rather
than asking an agent to estimate them.

## Security And Retention Constraints

Update the threat model and security requirements for registered execution:

- The gate registry is trusted project configuration and changes to it are
  auditable, reviewable changes.
- No implicit shell expansion. An explicit shell-mode escape, if retained, must
  be visibly higher risk and policy-controlled.
- Scrub sensitive flags and environment data before journal persistence.
- Bound output capture and record redaction/truncation metadata.
- Keep only safe metadata, excerpts, and hashes in the Git-committed journal.
  Store optional full logs as content-addressed local or CI artifacts under an
  explicit retention policy.
- A receipt imported from an untrusted caller must declare its trust source and
  must not silently satisfy a trusted-execution rule.
- Preserve the no-network, no-daemon, local-first default.

## Compatibility And Migration

Preserve v0.1 history and the journal's append-only contract:

- Never rewrite historical events. Project legacy gate rows as
  `source=reported` or an equivalent explicit trust classification.
- Version new journal events and projection migrations.
- Keep the existing `fid dev gate --cmd --exit-code` surface temporarily as a
  deprecated reported-evidence/import path, or introduce an explicitly named
  replacement. Do not let old and trusted modes look identical.
- Make completion Markdown optional through an additive workflow/profile change
  before considering flag removal.
- Define JSON/TOON compatibility for receipts and capability discovery so a
  workflow kit can integrate without querying SQLite or depending on `.fid/`
  internals.
- Add a capability/version command or extend an existing machine-readable
  surface so integrations can fail clearly when required receipt features are
  absent. Never silently fall back from trusted receipts to reported evidence.

## Roadmap Re-Slicing Guidance

Do not simply append all work to v0.5. Propose a dependency-aware re-slice that
puts trustworthy evidence before prose-oriented context, export, dashboard, or
sync features.

Evaluate this starting sequence:

1. Keep v0.2a Workflow Engine Migration, adding only the seams needed for
   versioned gate obligations and future transition rules.
2. Insert a v0.2 evidence foundation: gate definitions, execution attempts and
   receipts, trust source, repository binding, freshness, maintenance checks,
   and compatibility migration.
3. Continue governance entities, but reassess whether reports/context packs and
   work-envelope/runner primitives should be grouped or separated.
4. Pull generated human attention/validation views and expanded economics ahead
   of one-way Markdown export, dashboard, MCP evaluation, and bidirectional sync.
5. Re-sequence QA plans, reviews, validation contracts, and outcomes so their
   first implementation links directly to first-class evidence rather than
   embedding path/string references that will immediately require migration.

The Tech Lead may choose different version labels. Preserve these dependency
rules:

- Workflow obligations can reference stable gate definitions.
- Receipt/evidence identity exists before QA and validation entities depend on it.
- A lean completion path exists before a workflow kit adopts `fid` as its state
  backend.
- Human/report projections follow the structured state they render.
- Dashboard, MCP, and Markdown sync do not displace the evidence foundation.

## Required Validation Scenarios

Add or update quality-plan scenarios for at least:

1. Registered gate passes on a known clean target tree and satisfies its exact
   obligation.
2. Nonzero exit, timeout, signal, harness error, and abandoned execution-start
   paths never produce a passing receipt.
3. A tree change after a pass makes the evidence stale and blocks a transition
   that requires current evidence.
4. A legacy caller-reported exit code remains queryable but cannot satisfy a
   trusted gate.
5. Secret-like arguments/environment values and oversized output do not leak
   into the journal; hashes and redaction/truncation flags remain correct.
6. Journal rebuild reproduces receipt, evidence, and freshness state exactly;
   merge/conflict behavior remains deterministic.
7. Completion without a Markdown report succeeds only when all structured
   completion obligations are met.
8. QA review with complete structured coverage and zero findings is accepted.
9. Manual Human evidence can satisfy only obligations that require or allow
   Human judgment.
10. A consumer using the machine-readable compatibility contract detects an
    unsupported `fid` version/capability and stops without weakening evidence.

## Non-Goals For This Update

- Do not implement the features.
- Do not design an LLM/model-hosting API.
- Do not add a daemon, network listener, remote service, or git mutation.
- Do not prioritize the web dashboard, MCP adapter, or bidirectional Markdown
  sync as a substitute for receipts.
- Do not introduce generic custom fields or a generic workflow scripting
  language.
- Do not move qualitative product/architecture/security judgment into opaque
  automated scoring.

## Required Outputs

Update the smallest coherent set of canonical documents, but cover every
affected contract. Expected outputs are:

1. A decision record for the executor/attester boundary, receipt trust model,
   completion-prose policy, and deterministic-runner direction.
2. Glossary entries for gate definition, execution attempt, execution receipt,
   trust source, evidence obligation, evidence freshness, work envelope, and
   human validation brief, or explicit reasons to use different terms.
3. New or amended PRD requirements with priorities and target versions.
4. TRD/schema changes covering entities, invariants, persistence, security,
   compatibility, retention, and module boundaries.
5. CLI/flow changes sufficient to make trusted and reported evidence visibly
   distinct and usable by an external workflow kit.
6. Quality-plan requirements and the scenarios listed above.
7. AI-workflow changes that relocate mechanical enforcement from prompts into
   `fid`/runner behavior while retaining Human authority over protected
   decisions.
8. A re-sliced implementation roadmap with dependency rationale, explicit
   deferrals, v0.1 compatibility treatment, and phase-sized increments.
9. A refreshed current phase brief stating that no implementation phase is
   active and naming the first proposed post-review phase, pending Human
   approval.

Run all documentation/specification lint and traceability checks applicable to
planning changes. Report any command that remains unavailable or placeholder.

## Completion Response To The Human

Return a concise summary containing:

- Decisions recorded and any decisions still requiring Human approval.
- Requirements added or materially changed.
- Old versus new roadmap sequence.
- Compatibility and migration consequences for v0.1 dogfood data.
- The proposed first implementation phase, its risk/operating mode, and why it
  is the correct first dependency.
- Exact checks run and their outcomes.
- A clear statement that no product implementation began.

Suggested invocation:

```text
Use the Tech Lead role. Read
project-docs/docs/reviews/2026-07-18-fid-roadmap-update-tl-brief.md and perform
the planning/specification checkpoint it defines. Use its recommended defaults
unless you identify a concrete contradiction with fid's recorded product or
architecture decisions. Update the canonical specs and roadmap, run applicable
documentation checks, and stop before implementation for my review.
```
