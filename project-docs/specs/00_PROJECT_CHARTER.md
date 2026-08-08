# PROJECT CHARTER - fid

> This document records the durable intent of the project: why it exists, who
> it serves, what is in and out of scope, and which decisions should survive
> beyond individual implementation phases.

## 1. Vision

`fid` is a local-first, agent-first issue tracker and durable operational
work-state and evidence backend for AI-assisted software development. It
provides a structured, queryable graph over issues and their lifecycle, work
relationships, claims and runs, events and projections, receipts and evidence,
reviews and outcomes, attention, and telemetry. Its append-only journal is
canonical and its SQLite projection is driven entirely through a CLI.

The name comes from a tapered hand tool used to open up tightly-laid rope so
the strands can be worked separately. The metaphor fits: fid opens up
tightly-coupled project work into separable, trackable parts that agents and
humans can manipulate independently.

The long-term vision is a single binary that replaces scattered operational
Markdown and manual bookkeeping with shared, schema-enforced, queryable work
state without requiring a server, auth layer, or external service. fid does
not own the engineering method around that state: an external kit owns risk,
authority, role routing, deterministic orchestration, and provider projections.
Project files and Git remain canonical for product intent, architecture,
decisions, code, and repository identity; Humans retain product judgment and
protected-decision authority. (ADR-0004)

## 2. Problem Statement

AI-assisted development workflows produce a high volume of structured
artifacts: task specs, completion reports, QA reviews, decisions, escalations,
triage records, and retrospectives. Today these live as markdown files in a
directory tree (e.g. the `whet` project's `project-docs/` layout), queried
via `grep`, `git log`, and agent reading.

This is insufficient because:

- **Cross-cutting queries are expensive.** "Every decision that touched the
  data model", "all tasks where actual tokens exceeded 1.5x estimate", or
  "blocking findings still open across all epics" require multi-file grep
  plus agent interpretation, burning context window and tokens.
- **Workflow invariants are unenforced.** Required fields on completion
  reports, AC verification before transitions, and concern resolution before
  dev-complete are conventions, not constraints. Violations are caught late
  (QA review) or not at all.
- **Non-task entities are second-class.** Decisions, escalations, triage
  events, and carry-forward items lack identity, status, and relationships.
  They exist as sections within files, not as queryable rows.
- **Token and cost telemetry is ad-hoc.** Agent run costs are estimated from
  conversation length, not measured and recorded per task and per role.

Meaningful improvement looks like: single CLI commands that answer cross-cutting
questions, schema-enforced workflow transitions, first-class identity for every
project artifact type, and computed retrospective signals.

## 3. Target Users / Consumers

| Segment | Description | Primary Needs |
|---|---|---|
| Primary | Solo developer or small team running an agentic coding workflow directly or through an external workflow kit | Queryable operational work state, enforced integrity/evidence invariants, telemetry, cross-issue traceability |
| Secondary | Orchestrator, Builder, Reviewer, and deterministic-harness consumers | Versioned machine contracts, issue/work projections, claims, receipts, evidence, attention, and resume/recovery state |
| Secondary | Human product owners reviewing and approving work | Markdown export for human reading, human verification checklists, dashboard (v0.5) |
| Not target | Teams needing real-time multi-tenant collaboration | fid is local-first, single-project-per-DB, last-write-wins; no auth, no remote sync, no CRDTs |
| Not target | Users wanting a custom-fields platform or general-purpose project manager | Schema is opinionated; no user-defined fields in v1 |

## 4. Design Principles

These principles guide tradeoffs when requirements conflict.

1. **Rows hold structure; markdown holds prose.** The schema captures identity,
   status, relationships, and metadata. Prose stays in markdown files
   referenced by `*_md_path` columns. Don't shred prose into columns; don't
   leave structured data in unindexed files.
2. **Journal is canonical; DB is a derived cache.** The append-only JSONL
   journal committed to git is the source of truth. The SQLite DB is
   gitignored and rebuilt from the journal. This gives full SQL query power
   without binary merge conflicts. (ADR-0001)
3. **Opinionated defaults; configurable surface.** Ship working out of the box
   with a default role profile, default workflows, and a three-tier QA model.
   Make roles, workflows, and permissions per-project configurable so fid
   supports diverse workflows without dictating one. (ADR-0002, ADR-0003)
4. **Schema-enforced invariants over convention.** Required fields, workflow
   transition rules, kind constraints, and permission boundaries are enforced
   at write time — application-layer primary, SQLite triggers as backstop.
   Don't rely on agents following instructions.
5. **Append-only history for audit-grade entities.** Task completions, AC
   verifications, QA reviews, triage events, and assertion outcomes are
   append-only. Re-attempts create new rows; they don't edit old ones. The
   history is the value.
6. **Agent-friendly by default.** JSON output for universal pipeability, TOON
   for token efficiency, role-tuned context packs to minimise context window
   pressure, claim-based exclusive access to prevent concurrent agent
   conflicts.
7. **One authoritative owner per kind of truth.** fid owns operational work
   state and evidence. The workflow kit owns engineering-workflow policy and
   orchestration. Project files and Git own product/code truth; Humans own
   judgment and protected decisions. Integrations exchange versioned records
   and references rather than creating competing copies. (ADR-0004)

## 5. Scope

### 5.1 Core Concepts

Define domain concepts without committing to implementation details. Canonical
term definitions live in `project-docs/specs/GLOSSARY.md`; keep this section as
a brief concept inventory, not a duplicate glossary.

| Concept | Description |
|---|---|
| Issue | Unified entity for epics, tasks, and subtasks; Issues and their relationships are the only durable work graph (ADR-0005) |
| Operational work state | Durable issues plus execution facts, evidence, attention, and projections that agents query and resume from |
| Epic | Top-level container of related tasks; sequential or prioritised scheduling |
| Workflow | Named state machine with states, transitions, and per-transition validation rules |
| Claim | Short-lived exclusive-access lease on an entity (issue, QA review) |
| Journal | Append-only JSONL event log; canonical source of truth, committed to git |
| Concern / Finding / Escalation | Quality signals raised at different workflow stages by different roles |
| Validation contract | Assertion-level specification for an epic's proof obligations |
| QA plan | Verification design citing contract assertions; declares evidence methods |
| Gate | Deterministic command producing pass/fail (tests, linters, etc.); trusted gates are registered definitions executed by fid |
| Gate definition / Execution receipt | Versioned registered gate entry, and the immutable attested record of one execution (trust source, repository identity, terminal status) |
| Evidence obligation | Declared requirement that a transition/verification/outcome be backed by specific, valid (current) evidence |
| Work envelope | Versioned execution contract over a primary Issue and an explicitly bounded related Issue set: outcome, authority, obligations, budgets, and stop conditions. It does not duplicate the Issue graph. |
| Role | Workflow position (developer, qa_reviewer, qa_lead, tl_lite, tl, workflow_runner, human) with permission boundaries |

### 5.2 In Scope

- Unified issue tracker (epic / task / subtask) with configurable workflows
- Acceptance criteria with multi-role verification
- Registered gate execution with trusted, immutable receipts; evidence
  records with obligations and freshness (v0.2)
- Structured completion path with generated completion views; completion
  prose optional (v0.2)
- Work envelopes and deterministic orchestration primitives for an external
  runner kit (v0.3); fid hosts no LLM loop
- Append-only task completions and structured developer concerns
- QA reviews with findings, QA plans with steps (v0.3)
- Validation contracts with assertions, required validators, and outcomes (v0.4)
- Decisions, escalations, triage events, carry-forward items, retrospectives
- Typed cross-entity links with auto-inverse
- Claim-based exclusive access with TTL
- Actor/capability and claim precondition enforcement over approved authority;
  the shipped role matrix remains a compatibility surface pending roadmap
  disposition
- JSONL journal as source of truth; SQLite as derived query cache
- Multi-format output: JSON (default), TOON (opt-in), Markdown, YAML
- Token/cost telemetry per agent run with snapshot pricing
- Named reports and stats aggregations
- Markdown export with YAML frontmatter contract (v0.4)
- TUI kanban board (v0.3)
- CLI/skill-first agent integration; MCP adapter deferred until dogfooding
  proves a concrete MCP-only client need
- Templates with monotonic versioning (v0.3)
- Single static Go binary, project-scoped `.fid/` directory

Items in the pre-ADR-0004 roadmap that encode a particular engineering method
(including role choreography, QA-plan/validation hierarchies, templates, and
role-specific reports) require explicit forward-roadmap disposition. Their
presence here does not assign workflow-policy ownership back to fid.

### 5.3 Out of Scope

- Multi-tenant server, remote sync, auth, or hosting
- Real-time collaboration or CRDT-based conflict resolution
- Hosting models/agents, general job scheduling, shell automation, or git
  mutation — fid executes only registered gate definitions (DEC-059)
- Owning risk classification, authority/role routing, provider projections, or
  the engineering-workflow policy and orchestration loop
- Replacement for Git, PRDs/TRDs, canonical project decisions, or Human
  product judgment
- User-defined custom fields
- Multi-project queries across separate databases
- SLA tracking, burn-down charts, customer requests, OKRs
- Multiple assignees, subscribers, or watchers per issue
- Initiatives or hierarchy levels above epics
- Time tracking in hours (token telemetry replaces this)
- Bidirectional markdown sync (deferred to v0.5+ pending evidence of use case)

## 6. Interfaces / Surfaces

List the user-facing or integration surfaces the project will expose.

| Surface | Purpose | Notes |
|---|---|---|
| `fid` CLI | Primary interface for all reads and writes | Single Go binary; nested command tree via Cobra. Ships v0.1. |
| JSONL journal | Canonical event store committed to git | `.fid/journal/YYYY-MM-DD.jsonl`; one event per CLI write. Ships v0.1. |
| SQLite database | Derived query cache for fast reads | `.fid/fid.sqlite`; gitignored; rebuilt from journal. Ships v0.1. |
| TUI board | Kanban visualisation in the terminal | `fid board`; Bubbletea + Lipgloss. Ships v0.3. |
| Agent skills | Agent workflow guidance layered over the CLI | Repo skills explain fid-aware workflows for Codex/Claude-style agents. CLI remains the primary integration surface. |
| Markdown export | Human-readable project state as a file tree | `fid export`; YAML frontmatter contract. Ships v0.4. |
| Web dashboard | Browser-based project visualisation | Separate `fid-dashboard` package; reads SQLite. Ships v0.5. |

## 7. Success Criteria

Define objective or observable outcomes.

| Horizon | Criteria |
|---|---|
| Early dogfood | v0.1 core issue layer replaces ~70% of daily interaction with whet's task artifacts. `fid init`, issue CRUD, AC verification, claims, ready queue, link management, and journal/rebuild all work reliably. The author uses fid for real project work. |
| Public release | v0.3 ships QA reviews, templates, and TUI board. The full agentic task loop (TL plan → Developer execute → QA review → TL wrap-up) runs through fid's CLI with enforced workflow invariants. Cross-epic queries answer in single commands. |
| Long-term | fid is a stable issue-tracking and operational work-state/evidence layer for AI-assisted development workflows. Projects and evolving workflow kits adopt it through versioned capability contracts without forking or depending on storage internals. Telemetry and trusted evidence make workflow effectiveness measurable. |

## 8. Constraints

- **Time / staffing:** Solo developer building with AI agent assistance. No team; no external contributors initially.
- **Budget:** Anthropic API costs for agentic development are the primary expense. Snapshot pricing in `agent_runs` tracks this.
- **Security / privacy / compliance:** Local-first; no data leaves the machine. No auth surface. Journal committed to git inherits the repo's access control.
- **Operational:** Single static binary; no daemon, no background process, no external service dependencies. SQLite only.
- **Legal / licensing:** Open source under the MIT licence. No third-party service integrations that impose licensing constraints. Pure-Go SQLite driver (`modernc.org/sqlite`) avoids CGo. No monetisation planned.

## 9. Landscape

fid occupies a narrow niche: a local-first, schema-enforced issue tracker
designed for AI-assisted development workflows. Nothing else combines
journal-canonical persistence, role-based agent permissions, and a
three-tier QA model in a single CLI binary.

| Tool | Overlap | Key Differences |
|---|---|---|
| **Linear / Jira / GitHub Issues** | Issue tracking, workflows, labels, priorities | Cloud-hosted, multi-tenant, UI-first. No journal-canonical persistence, no agent role model, no AC verification or QA review workflow. fid borrows lightweight features (labels, inbox, status categories, estimates) but not the collaboration or integration surface. |
| **Plain markdown + grep** | The status quo fid replaces. Markdown files in a git repo queried by agents reading files. | No schema enforcement, no cross-cutting queries without agent interpretation, no workflow invariants, no structured telemetry. fid keeps prose in markdown but indexes structure in SQLite. |
| **Beads** | Go CLI, SQLite storage, project-scoped. Similar architectural choices. | General-purpose task tracking without agent-specific concerns. No journal-canonical model, no role-based permissions, no QA layer, no validation contracts. fid's architecture was informed by Beads' Go/SQLite patterns. |
| **Beans** | Go CLI, local-first, developer-oriented project tracking. | Focused on time tracking and invoicing. Different domain; no workflow enforcement or agent integration. |
| **Claude Code TaskCreate / built-in tasks** | In-conversation task tracking for AI agents. | Ephemeral, conversation-scoped, no persistence beyond the session. fid provides the durable, queryable, cross-session layer underneath. |
| **Task Master (various)** | MCP-based task management for AI agents. | Typically exposes a broad tool surface to agents, leading to context bloat. fid stays CLI/skill-first until dogfooding proves an MCP-only client need. Task Master's lessons informed fid's lean tool design. |

## 10. Durable Decisions Log

Keep only project-shaping decisions here. Phase-local decisions belong in
`project-docs/docs/phases/phase-N/decisions.md`.

| Decision | Rationale | Date / Reference |
|---|---|---|
| Journal-canonical, DB-derived, no daemon | Gives full SQL query power without binary merge conflicts or daemon failure modes. Journal committed to git; DB is gitignored and rebuilt. | ADR-0001 |
| Roles configurable via seeded default profile | Supports diverse workflows without dictating one. Roles table in v0.1; permissions table in v0.2. | ADR-0002 |
| Three-tier QA model shipped in full | Primary user's workflow requires assertion-level verification. Plans + reviews in v0.3; contracts + assertions + outcomes in v0.4. All layers schema-optional per project. | ADR-0003 |
| Implementation language: Go | Agentic build pipeline proven on whet transfers directly. Single static binary distribution. Agent error-recovery loop more efficient than Rust during design iteration. | 2026-05-09, spec §6.16 |
| Application-layer validation primary, triggers as backstop | Rich error messages and parameterised rules need app-layer context. Triggers guard structural invariants for non-CLI write paths (future dashboard, direct SQL). | spec §6.14 |
| CLI/skill-first agent integration | Coding agents already use shell commands effectively. `fid --help`, structured output, examples, and repo skills are the primary agent interface. MCP is deferred until a concrete MCP-only client need appears and a product/security decision records the transport and threat model. | 2026-05-29 decision |
| Hash-based IDs, not sequential | Reduces coordination needs for concurrent agent/branch issue creation. Collision risk managed by progressive length scaling (6→8 hex chars). | spec §6.0 |
| fid executes registered gates; reported evidence is never trusted proof | Agent-reported exit codes cannot satisfy deterministic obligations; trusted evidence is an executed/attested receipt from a versioned registry. fid remains local-first: no agents, models, git mutation, or scheduling. | 2026-07-18, DEC-059/060 |
| Evidence has identity, trust source, and freshness | Receipts and typed evidence records are first-class rows; validity (current/stale) is separate from existence; conservative whole-tree staleness first. | 2026-07-18, DEC-061 |
| Completion prose is optional; views are generated | Structured `has_completion` obligations replace the mandatory completion Markdown; human-readable completion views and the validation brief are generated projections of structured state. | 2026-07-18, DEC-062/064 |
| fid exposes runner primitives; the runner is a separate deterministic kit | Work envelopes + machine-readable orchestration primitives in fid; no hosted LLM loop. Humans keep protected-decision authority. | 2026-07-18, DEC-063 |
| fid owns operational work state; the kit owns engineering workflow policy | Keeps fid independently useful as an issue tracker and stable backend while risk, roles, orchestration, and provider capabilities evolve in the kit. Project/Git/Human truth stays with its authoritative owner. | ADR-0004 |

## 11. Document Suite

This charter is part of the core specification suite:

- `00_PROJECT_CHARTER.md` - vision, scope, constraints, durable decisions.
- `GLOSSARY.md` - canonical project language and flagged ambiguities.
- `01_PRD.md` - product requirements.
- `02_TRD.md` - technical requirements.
- `03_TRD_SCHEMA.md` - full SQL schema, triggers, and validation.
- `04_TRD_CONFIG.md` - `.fid/` directory structure and config formats.
- `05_DESIGN_PLAN.md` - user experience and interface design.
- `06_DESIGN_PLAN_CLI.md` - CLI flag specs and output examples.
- `07_DESIGN_PLAN_FLOWS.md` - end-to-end flows, whet mapping, export layout.
- `08_QUALITY_PLAN.md` - testing strategy and quality gates.
- `09_AI_WORKFLOW.md` - agent roles and orchestration.
- `10_IMPLEMENTATION_ROADMAP.md` - implementation phases.

*Last updated: 2026-08-08 (ADR-0004 product-boundary ratification)*
