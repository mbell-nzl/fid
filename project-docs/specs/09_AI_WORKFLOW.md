# AI WORKFLOW INTEGRATION - fid

**Status:** Target V2 integration contract
**Last updated:** 2026-08-09
**Authority:** [ADR-0011](../../docs/adr/0011-realign-fid-around-the-durable-operational-kernel.md)

## 1. Scope

This document explains how an external agentic-engineering kit uses fid. It is
not a seven-role process manual and it does not prescribe how fid itself is
currently developed.

The integration has three distinct owners:

| Owner | Owns |
|---|---|
| Project | Requirements, architecture, decisions, source code, validation strategy, and approved Gate Definition values. |
| Kit | Risk policy, work selection, role routing, prompts, skills, provider projections, Work Envelope compilation, and orchestration. |
| fid | Issues, immutable Envelope Revisions, authority, Claims, Runs, Stops, registered gates, Receipts, Evidence Records, reviews, Human Outcomes, Attention, recovery, and projections. |

The Human retains product judgment and protected-decision authority. fid
records the operational effect of that authority but does not invent it.

## 2. External collaboration roles

The V2 kit uses three agent roles:

- **Orchestrator:** selects Issues, applies project policy, compiles and proposes
  a Work Envelope, activates actors, routes Stops, and requests Human input.
- **Builder:** claims one selected Issue and declared mutation scope, changes
  the project, runs applicable gates through fid, and records concise outcomes.
- **Reviewer:** performs independent review against the active Envelope
  Revision and records a Review Record and Findings.

These names belong to the kit. fid stores Actors, capabilities, provenance, and
Authority Grants. It must not encode an assumption that every client uses these
role names or that role membership alone grants authority.

A Human and a trusted execution harness are also Actors. They are not extra kit
agent roles: the Human supplies protected judgment, while the harness supplies
process observations that agents cannot self-report as trusted Receipts.

## 3. Unit of work

An Issue is the only durable unit in fid's work graph. There is no second
`WorkItem` entity.

A Work Envelope authorizes execution over one primary Issue and an explicitly
selected related Issue set. Each Builder claims one selected Issue, not the
whole Envelope. A procedural step that has no independent lifecycle remains
transient orchestration context; it becomes a child Issue only when it needs
its own claim, dependency, acceptance, or durable state.

This preserves a useful issue tracker while allowing the kit to package enough
policy for safe execution without duplicating the backlog.

## 4. Policy compilation

The project stores policy inputs in machine-readable project configuration and
supporting prose. The kit validates those inputs and compiles the applicable
result into an immutable Envelope Revision:

- selected Issues and intended outcome;
- exact project-source references and digests;
- repository and mutation scope;
- risk result and routing consequences;
- Actor identities, capabilities, and Authority Grants;
- protected Decision classes;
- registered gate references;
- typed Obligations and their trust, freshness, coverage, and independence
  requirements;
- budgets and stop conditions.

An Obligation is an enforceable condition, not another markdown artifact or
generic dictionary synonym. It may require deterministic gate evidence,
acceptance evidence, independent review, Human judgment, or a declared
combination. Only applicable policy is compiled; routine work does not acquire
ceremonial records merely because the data model can represent them.

Risk lanes or capability profiles are kit policy labels. fid may record their
result or provenance when needed for audit, but it does not define a universal
`Routine`, `Standard`, or `Escalated` workflow. The operational contract is the
actual authority and Obligations on the Envelope Revision.

## 5. Gate registration and execution

The project owns approved Gate Definition values. The kit owns their schema,
validation, overlays, and activation. fid stores an immutable Gate Definition
Registration binding an approved value and digest to its project, repository,
scope, provenance, and registration version.

When an authorized Actor requests a gate run, fid:

1. resolves the immutable registration;
2. creates an Execution Attempt before spawning the process;
3. runs the registered executable and argument vector under its declared
   working-directory, environment, timeout, output, and path restrictions;
4. observes termination and captures sanitized output and artifact identities;
5. appends an immutable Execution Receipt; and
6. evaluates only the Obligations to which that Receipt can apply.

The terminal may stream command output for ergonomics. The agent still cannot
satisfy deterministic evidence by reporting command text, an exit code, or a
log it ran independently. Explicit attestation can be accepted only where the
Obligation permits that trust class.

## 6. Normal execution flow

1. The Orchestrator asks fid for ready Issues and relevant bounded context.
2. The kit applies project policy and proposes an Envelope Revision.
3. fid validates references, capabilities, authority, gate registrations, and
   invariants, then activates the revision atomically.
4. A Builder acquires a Mutation Claim for one selected Issue and declared
   scope, then starts a Run bound to the Claim and repository identity.
5. The Builder changes the project and invokes applicable registered gates
   through fid.
6. fid records Receipts and Evidence Records and exposes the remaining
   Obligations and Next Action as projections.
7. A Reviewer claims an independent Review Run when policy requires it and
   records a Review Record and Findings.
8. A Human records any required protected Decision or Human Outcome.
9. fid permits completion only when the Issue workflow, Claims, Runs, Stops,
   Obligations, review coverage, Findings, and Human Outcomes allow it.
10. The Orchestrator selects the next permitted action from current state.

The kit orchestrates; fid enforces and remembers. Neither writes a narrative
completion report merely to hand state to the other.

## 7. Parallel Builder rule

Multiple Builders may work within one active Envelope Revision only when each
holds a Mutation Claim for a different selected Issue and the declared mutation
scopes do not overlap.

fid rejects:

- two live Mutation Claims for the same Issue;
- overlapping or unknown mutation scopes when exclusivity is required;
- mutations made against a stale Claim or Envelope Revision;
- completion while relevant Claims or Runs remain live; and
- a Builder resolving its own Review Finding without later qualifying review.

The kit may serialize work more conservatively. It may not weaken these
invariants.

## 8. Interruption and recovery

An abnormal executor exit terminates the current Run as interrupted. Recovery
does not reopen it. The next Orchestrator or Builder:

1. reads the recovery projection;
2. verifies repository and worktree identity;
3. reconciles Claim liveness and any partial effects;
4. records a typed Stop or authorized disposition when state is unsafe or
   ambiguous; and
5. starts a linked new Run only after continuation is permitted.

Routine status is available through projections and does not interrupt the
Human. Attention Items are reserved for decisions, blockers, circuit breakers,
risk acceptance, or readiness for Human validation.

## 9. Revision and change

Material changes to outcome, selected Issues, repository scope, risk result,
authority, protected Decisions, Obligations, budgets, or stop conditions create
a new Envelope Revision. Existing Claims and Runs follow explicit revision and
staleness rules; prior Receipts, reviews, and Human Outcomes retain their
original lineage.

Changes to project prose are detected through stable references and digests.
The kit decides what policy should be recompiled. fid decides whether existing
operational records remain applicable under the accepted freshness rules. It
does not parse prose to recreate policy.

## 10. Context and Human interfaces

Agents request bounded JSON projections for a named question, such as ready
work, current obligations, recovery, evidence trace, or changed-since state.
JSON is the only canonical machine representation.

Humans may receive concise terminal text or a named Human Validation Brief in
purpose-built Markdown or HTML. Such output is derived, decision-focused, and
noncanonical. Generic Markdown, YAML, and TOON format switches are outside the
workflow contract.

## 11. Stops and protected decisions

Work stops when continuation would violate authority, integrity, an active
Envelope Revision, or an applicable Obligation. Typical reasons include:

- missing or ambiguous protected-decision authority;
- invalid or stale project references;
- conflicting Claims or repository identity;
- failed deterministic gates;
- unavailable required artifacts;
- unresolved blocking Review Findings;
- budget or retry circuit breakers; and
- inability to establish trusted evidence.

The Orchestrator may resolve operational routing within its Authority Grant.
Only the Human may decide a protected class. fid records the Decision Reference
and Outcome; canonical rationale remains in the project.

## 12. Explicit exclusions

fid does not own or generate:

- role prompts, provider agents, or skills;
- phase briefs, task specifications, completion reports, or handoffs;
- QA Plans, Validation Contracts, assertion hierarchies, or command prose;
- architecture and product decision documents;
- operating-mode artifact matrices;
- triage logs, escalation reports, or routine retrospectives; or
- a second durable work graph beside Issues.

Useful operational facts formerly buried in those artifacts are represented as
typed Claims, Runs, Stops, Evidence Records, Reviews, Findings, Decisions,
Outcomes, or Attention. The rest remains project-owned prose or disappears as
reconstructible narration.

## 13. Transition note

fid's repository still uses the V1 seven-role development process while the V2
kit and target backend port are incomplete. That process remains binding for
current implementation work under `AGENTS.md`; this document neither repeals
it nor makes it part of the product. The first live project adoption occurs
only after the clean fid/V2 baseline passes the quality plan.
