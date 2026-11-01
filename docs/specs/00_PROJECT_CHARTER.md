# PROJECT CHARTER - fid

## 1. Purpose

fid is a local-first, agent-first Issue tracker and durable operational
work-state and evidence backend for AI-assisted software development. It gives
Humans, agents, and an external workflow kit one rebuildable record of what
work exists, what execution is authorized, what happened, what evidence proves
it, and what requires attention.

fid is independently useful as an Issue tracker. When paired with the V2
Agentic Engineering Kit, it implements the kit's versioned backend-port
contract without owning the engineering method around that contract.

## 2. Ownership Boundary

| Truth | Authoritative owner |
|---|---|
| Issues, relationships, Issue workflow state | fid |
| Work Envelopes, Claims, Runs, Stops, operational events | fid stores and enforces the approved contract |
| Gate registrations, Attempts, Receipts, Evidence, Reviews, Outcomes, Attention, telemetry | fid |
| Risk classification, role routing, obligation selection, orchestration policy | External kit and consuming project |
| Requirements, architecture, validation strategy, decision prose, code, repository identity | Project files and Git |
| Product judgment, protected decisions, risk acceptance, release authority | Human |

fid never recalculates kit risk or routing policy. The kit never creates a
second durable work graph or evidence store. Integrations exchange versioned
records and references through the backend port, never SQLite or journal
internals.

## 3. Product Principles

1. **Issues are the only durable work graph.** A Work Envelope selects one
   primary Issue and an explicitly bounded related Issue set; it does not
   create persistent Work Items.
2. **Journal canonical, SQLite derived.** Every durable mutation appends an
   ordered event. SQLite is a rebuildable query projection.
3. **One durable fact, many projections.** Next Action, progress, completion,
   context, reports, and Human briefs are regenerated views, not authored state.
4. **Policy outside, enforcement inside.** The kit compiles approved policy;
   fid enforces identity, revision, authority, lease, evidence, freshness,
   non-weakening, idempotency, and lifecycle preconditions.
5. **Evidence is observed, not reported.** Deterministic Obligations are
   satisfied only by immutable executed or qualifying attested Receipts.
6. **Records must earn durability.** A first-class record must affect permitted
   action/completion, enable recovery/rebuild/audit, prove an Obligation or
   authority, or provide trustworthy operational measurement.
7. **JSON is the machine contract.** Human output is concise terminal text or
   a purpose-built named projection. Universal format matrices are out of scope.
8. **Fail closed.** Missing capability, authority, evidence, freshness,
   executor liveness, or repository certainty blocks rather than defaults.

## 4. Users

| User | Primary need |
|---|---|
| Solo developer or small team | Local Issue tracking, clear readiness, durable evidence and recovery without a service |
| Orchestrator, Builder, Reviewer, deterministic harness | Versioned JSON contracts for work, authority, Claims, Runs, evidence, review, Stops, and resume |
| Human product owner | Decision-focused Attention and Human Validation Briefs without reading agent transcripts |
| Maintainer/operator | Rebuild, integrity checks, portable export, diagnostics, and trustworthy trajectory measurement |

Not target users: teams requiring a hosted multi-tenant service, real-time
collaboration, remote synchronization, auth administration, or a general
custom-fields project-management platform.

## 5. Core Product

### 5.1 Issue tracker

- Epic, task, and subtask Issues with parent/dependency relationships.
- Configurable Issue workflows and assignments.
- Acceptance criteria, labels, links, prioritization, and actionable-work
  queries with typed reasons.
- Stable IDs, append-only history, integrity checks, and rebuild.

### 5.2 Durable operational kernel

- Work Envelopes and immutable revisions over selected Issues.
- Generic actor identity, capabilities, approved authority grants, Claims,
  Runs, budgets, Stops, Attention, recovery, and Kit Activation.
- Immutable Gate Definition Registrations, Attempts, Receipts, Evidence
  Records, Referenced Artifact state, and freshness-aware Obligations.
- Review Records and Findings, Concerns, Deviations, protected Decision
  references, Human Outcomes, and structured completion.
- Repository-bound trajectory telemetry and portable structured export.

### 5.3 Required projections

- Actionable work and typed non-actionable reasons.
- Resume and changed-since context.
- Current Envelope, Obligation, evidence, review, and Attention status.
- Completion/evidence trace and Human Validation Brief.
- Basic trajectory economics and integrity diagnostics.

## 6. Explicitly Out Of Scope

- Hosting models, agents, or the orchestration loop.
- Owning risk policy, roles, routing, provider projections, or Human judgment.
- Replacing Git or canonical requirements, architecture, decisions, validation
  strategy, or source code.
- QA Plan, Validation Contract/assertion, template, handoff, triage,
  role-composite, or retrospective product aggregates.
- Caller-reported deterministic evidence or completion-report enforcement.
- Configurable fid-global roles or permissions.
- Universal TOON, YAML, or Markdown output.
- Markdown import, bidirectional sync, or whet-style file-tree export.
- A daemon, listener, remote service, multi-tenant database, or CRDT sync.
- A supported direct-SQL integration contract.

## 7. Pre-Dogfood Baseline

No real project has adopted fid and no journal history requires compatibility.
The pre-dogfood CLI, schema, migrations, fixed role matrix, generic format
engine, and event shapes may be removed or replaced. The first supported
baseline is the aligned product after the Gate D convergence work.

Implementation remains incremental and test-driven. Each replacement slice
keeps the repository buildable; schema and event migrations are consolidated
before first dogfood.

## 8. Success Criteria

| Horizon | Observable result |
|---|---|
| Contract qualification | All required backend capabilities pass shared fake/fid conformance, rebuild, interruption, evidence, and export scenarios |
| First dogfood | A real Routine-lane task is selected from fid Issues, executed through registered gates, reviewed where required, completed on current evidence, rebuilt, and resumed from a fresh session |
| V2 adoption | A project replaces V1 workflow files with kit policy plus fid operational state without dual persistence or hand-authored completion/evidence reports |
| Sustained use | Human attention, quality, retries, token/cost trajectory, and recovery nuisance are measurable and compare favorably with matched V1 work |

## 9. Constraints

- Go, SQLite, local filesystem, single static CLI binary.
- Local-first, single project/workspace, single-user trust root.
- No network or daemon by default.
- Additive or breaking public contracts are versioned and tested once the first
  supported baseline is declared.
- Dependency additions require recorded need, alternatives, and license review.

## 10. Durable Decisions

| Decision | Record |
|---|---|
| Journal is canonical; SQLite is derived | ADR-0001 |
| fid owns operational state; kit owns workflow policy | ADR-0004 |
| Issues are the only durable work graph | ADR-0005 |
| Parallel mutation requires Issue/worktree/scope isolation | ADR-0006 |
| Interrupted work resumes through a linked new Run | ADR-0007 |
| Review rework stays on the original Issue | ADR-0008 |
| Multi-Issue progress is derived; material changes create Envelope revisions | ADR-0009 |
| Phase 20 and fid-global configurable role policy are superseded | ADR-0010 |
| Product scope is the durable kernel plus necessary projections; no pre-dogfood compatibility; JSON only | ADR-0011 |

The evidence-gated candidate ledger is in ADR-0011. TUI/dashboard and search
are the leading optional candidates, but neither has a version promise.

*Last updated: 2026-08-09 (Gate D2 / ADR-0011 realignment)*
