---
status: accepted
date: 2026-08-09
---

# Realign fid around the durable operational kernel

> **Amended on 2026-10-31 by ADR-0013.** fid evaluates project lane policy
> from `.fid.toml` on `origin/<default branch>`. Activation survives only as
> "the policy fid trusts, from `origin/main` at commit X". The convergence
> target with the V2 kit is withdrawn. The record criteria and retirement list
> stand.

> **Amended on 2026-11-10 by ADR-0014.** The durable model's Work Envelopes,
> Claims and Runs, Stops, Attention, Kit Activation state, Gate Definition
> Registrations, Execution Attempts, Evidence Records, Obligation evaluation,
> Deviations and Human Outcomes go with the kernel after task 1.5. Issues, the
> journal, check results with Artifacts, review rounds and findings stay,
> joined by holds, the policy snapshot and a person's decisions. The policy
> snapshot replaces the surviving Activation. The ownership and durability
> test stands.

fid has not been adopted by a real project. No user workspace, journal history,
or external integration depends on its pre-dogfood commands, schemas, output
formats, or event contracts. Its first intended production use will follow the
fid/V2-kit convergence. Carrying speculative compatibility through that work
would preserve two competing product models without protecting a real user.

fid therefore commits only to its core issue tracker, the durable operational
kernel required by the V2 backend port, and the semantic projections needed to
operate, recover, evaluate, and understand that state.

## Ownership and durability test

fid creates a first-class record only when the fact:

- affects a permitted future action or completion;
- is required for resume, recovery, rebuild, or audit;
- proves authority, an Obligation, evidence, review, or Human outcome; or
- provides trustworthy operational measurement.

Project requirements, architecture, validation strategy, decisions, and other
durable prose remain project files. The kit owns risk policy, role routing,
workflow method, templates, skills, and provider projections. fid may store
stable references, digests, immutable registrations, compiled Work Envelope
Obligations, and operational outcomes without copying or authoring that policy
or prose.

## Durable model

The retained product model comprises:

- Issues, relationships, acceptance criteria, labels, and Issue workflows;
- Work Envelopes and immutable Envelope Revisions over selected Issues;
- generic actor authority, Claims, Runs, typed Stops, Attention, repository
  identities, Kit Activation state, and ordered events;
- Gate Definition Registrations, Execution Attempts, Execution Receipts,
  Evidence Records, Referenced Artifact state, and Obligation evaluation;
- Review Records, Review Findings, Concerns, Deviations, protected Decision
  references, and Human Outcomes;
- trajectory telemetry, integrity checks, rebuild, and portable structured
  export.

Next Action, partial progress, completion views, context, changed-since state,
Human Validation Briefs, reports, and statistics are noncanonical projections.
Not every Work Envelope invokes every record type; policy and observed events
select only the applicable records.

## Pre-dogfood convergence

fid makes no compatibility promise for its current pre-dogfood CLI, database
schema, migration chain, role policy, generic output-format matrix, or journal
events. Superseded surfaces are removed or replaced before first adoption,
without a legacy mode or deprecated-command period. Git history preserves the
discarded implementation.

The product boundary is sharp, but implementation proceeds in buildable
vertical slices. Each slice specifies and tests its generic replacement before
removing displaced code and tests. Before first dogfood, migrations and event
schemas are consolidated into a clean initial product baseline.

This amends ADR-0010: Phase 20 remains superseded and configurable fid-global
roles remain rejected, but the V1 permissions matrix and `fid policy show` are
no longer compatibility commitments. Generic authority must be designed before
the matrix is removed; absence of authority never becomes implicit permission.

## Output contract

JSON is fid's only canonical machine representation. The universal TOON, YAML,
and Markdown formatters and their `--format` promises are removed. Purpose-built
Human-readable output, including Markdown or HTML where it materially improves
a Human decision interface, is allowed only for a named projection and is
never a machine contract.

The current in-tree `TOON` encoder is a bespoke JSON-like syntax rather than a
conforming implementation of the evolving TOON specification. Token-efficient
projection encoding remains an evidence-gated candidate, not a promised
format.

## Feature disposition ledger

`Core` is committed product behavior. `Replace` keeps the need but supersedes
the old model. `Retire` removes the forward requirement and, because there is
no compatibility constituency, removes displaced implementation. `Candidate`
has no promised version or advance interface design.

| Feature family | Disposition | Forward boundary |
|---|---|---|
| Journal-canonical persistence, SQLite projection, locking, rebuild, maintenance, IDs | Core | Extend the same foundation for every retained aggregate. |
| Issue CRUD, hierarchy, dependencies, labels, ready queue, configurable Issue workflows | Core | fid remains an independently useful issue tracker. |
| Acceptance criteria | Core/Replace | Keep Issue outcome statements; replace role verification with applicable Evidence, Review, or Human Outcome records. |
| Claims, Runs, stops, resume, activity, telemetry | Replace | Bind generic actor, Issue, Envelope revision, authority, repository identity, budgets, and optional Obligation. |
| V1 role matrix and `fid policy show` | Retire | Replace with generic actor capabilities and recorded Envelope authority; no global kit-role taxonomy. |
| `fid dev start/gate/finish/block` and completion-report path | Retire | Generic port operations and structured completion replace role choreography. |
| Reported `gate_runs` and gate snapshots | Retire | Only registered execution Attempts/Receipts and explicitly trusted attestation can satisfy deterministic Obligations. |
| Requirements and protected Decisions | Replace | Project prose stays canonical; fid stores identity, digest/reference, authority, outcome, drift, and affected work. |
| Escalation workflow/entities, triage records/composites, TL/QA composites | Retire | Typed Stops, Decision references, Outcomes, Review Records, and Attention preserve the operational facts. |
| Concerns, Deviations, Review Records, Findings, Human Outcomes, Attention | Core/Replace | Use generic subjects and actors; keep concise behavior-affecting facts, not report prose. |
| Gate registry, execution, receipts, evidence, freshness, artifacts, capabilities | Core | Implement the V2 backend-port trust and conformance contract. |
| Work Envelope, revision, authority, Obligation, completion, activation, export | Core/Replace | Adopt the kit semantics; fid stores/enforces but does not recalculate policy. |
| QA Plan aggregate and command family | Retire | Project strategy is compiled into Envelope Obligations. |
| Validation Contract/assertion/validator hierarchy | Retire/Replace | Keep coverage, independence, trust, freshness, and outcome invariants directly on Obligations. |
| Templates, template validation/provenance, agent skills, handoff composer | Retire | Kit owns scaffold/provider inputs; generated context comes from semantic projections. |
| Carry-forward aggregate and retrospectives | Retire/Replace | Represent real follow-up as an Issue, Decision, Concern, risk, or Attention; occasional retrospective prose stays project-owned. |
| Ready explanations, context/resume/changed-since, completion/evidence trace, Human brief | Core projections | Named, bounded projections over durable state; never separately authored records. |
| Basic trajectory summaries | Core projection | Required to evaluate workflow economics; speculative pattern analytics remain conditional. |
| Portable structured export/import validation | Core | Backend-neutral recovery contract; no whet-style Markdown tree or bidirectional sync. |
| JSON | Core | Sole canonical machine representation. |
| Generic TOON, YAML, and Markdown output | Retire | Purpose-built Human Markdown/HTML is allowed only where justified. |
| `fid sql` | Development diagnostic | Explicitly unstable and never an integration or schema compatibility surface. |
| Inbox promotion, workflow inspection, requirement drift, `--dry-run` | Candidate tracker improvements | Schedule only when dogfood or a selected vertical slice establishes the need. |
| Governance command stubs and frozen V1 parity fixtures | Retire | Remove when their displaced contracts are removed. |

## Concise candidate ledger

Candidate entries record only the unmet-need trigger. They do not reserve a
version, dependency, schema, command, or implementation.

| Candidate | Reconsider when |
|---|---|
| TUI board or local dashboard | Dogfood repeatedly exposes a Human navigation/comparison problem that concise CLI or generated Human projections cannot solve. This is the leading optional UI candidate. |
| Lightweight search, then possibly a query DSL | Real questions repeatedly require unstable SQL or expensive multi-command joins and cannot be served by a stable named projection. This is the leading optional query candidate. |
| Saved views | Users repeatedly reconstruct the same useful query after the query surface exists. |
| Token-efficient projection encoding, including TOON | Representative fid payloads show at least 20% end-to-end context savings versus compact JSON across supported models, including instruction overhead and retries, with no correctness regression. JSON remains canonical. |
| MCP adapter | A supported client cannot use the CLI/machine contract reliably, or measured transport overhead materially harms operation. |
| Cycles or richer scheduling | Real backlog use exposes a need not served by priority, dependencies, and epics. |
| Journal merge or event-correction tooling | A real conflict or correction incident demonstrates semantics that normal append-only compensating events cannot handle ergonomically. Immutable Receipts are never reverted. |
| Retrospective and recurring-pattern analytics | Sufficient trusted trajectory data exists and a named Human decision depends on the result. |
| Multi-project federation | A concrete workflow must coordinate durable state across repositories. |
| Performance tier, checkpoints, or compaction | Representative benchmarks approach an accepted latency or rebuild limit. |

## Consequences

- The forward PRD, TRD, schema, design, workflow, quality, and roadmap documents
  require a coordinated rewrite. Historical phase records remain historical.
- Existing product code is not removed in this decision-only change. Removal is
  planned as TDD implementation work after the replacement contracts are
  coherent.
- fid's own current V1 development procedure remains an engineering-process
  concern until the V2 kit can replace it; it is not a fid product requirement.
