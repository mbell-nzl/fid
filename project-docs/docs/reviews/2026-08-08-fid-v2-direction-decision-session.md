# fid and Agentic Engineering Kit V2 - Direction Decision Session

Date: 2026-08-08
Status: In progress; nine decisions approved, remaining decisions open
Method: one decision at a time, challenged against canonical fid and kit docs

## Decision 1 - Product Boundary

Status: **APPROVED by Human on 2026-08-08**

fid remains a local-first, agent-first issue tracker and owns durable
operational work state and evidence. The kit owns engineering-workflow policy,
risk classification, authority and role routing, deterministic orchestration,
and provider projections. Project files and Git retain product, architecture,
decision, code, and repository truth; Humans retain product judgment and
protected-decision authority.

Rationale: this keeps fid independently useful and stable while model/provider
capabilities and the kit's workflow policy evolve. It avoids both an
all-encompassing fid workflow framework and split durable truth.

Consequent option disposition: the reconciliation report's Option 4
(kit-native durable backend replacing fid) is rejected by this boundary
decision, as is hybrid persistence. Options 1-3 remain available pending the
delivery-breadth decision.

Recorded in: `docs/adr/0004-operational-work-state-boundary.md` and
`project-docs/specs/00_PROJECT_CHARTER.md`.

Explicitly not decided at this point in the session:

- How Work Envelope execution identity binds to fid Issues. This is resolved
  later by Decision 6 / Gate B4 and ADR-0005.
- The final disposition of individual roadmap features, ADR-0002, or ADR-0003.

## Terminology Clarification - No Decision Required

`Routine capability profile` and `Routine production path` were analyst
shorthand in the reconciliation report, not canonical kit terms. The kit has a
capital-R **Routine risk lane** and independently versioned **backend
capabilities**. The precise phrase is: **the minimum backend capability set
needed to execute and recover one real Routine-lane task in production**.

## Decision 2 / Gate C - Delivery And Qualification Breadth

Status: **APPROVED by Human on 2026-08-08**

The active-project portfolio does not currently provide a cheap,
representative place to trial a partial fid/V2 capability set. Active projects
use the V1 workflow, so an early production trial would require temporary
workflow adaptation, dual operation, or partial migration. That setup could
consume much of the expected benefit of the narrow option and introduce risk
into the project used as the trial.

A lightweight artificial project can provide deterministic evidence for
capability correctness, interruption recovery, stale-evidence handling,
rebuild, export, and adapter conformance. It cannot by itself establish the
workflow economics, nuisance cost, Human attention burden, or maintainability
of V2 on sustained real work. The accepted kit plan already reflects this
distinction: V2-3 proves the backend contract and end-to-end recovery path,
while V2-6 compares matched V1 and V2 tasks and evaluates efficiency and
quality outcomes.

This weakens the case for Option 2 in the current portfolio. Its early
production feedback is not cheap or clean enough to justify making partial
live-project adoption the delivery constraint.

### Approved Resolution

Reframe Option 3 as **full-capability qualification with staged
implementation**:

- Implement all twelve specified backend-port capabilities before migrating a
  live project to V2.
- Deliver the capabilities as bounded vertical slices with contract and
  conformance evidence at each slice; do not treat the work as one undivided
  implementation phase.
- Use controlled fixtures and artificial workloads for deterministic backend
  and recovery evidence during implementation.
- Rehearse migration against a disposable copy or worktree of a V1 project
  before changing a live project.
- Start live adoption at a project or phase boundary, then perform the matched
  V1/V2 comparison required by V2-6 on representative work.

Here, `all twelve` means the twelve V2 kit backend-port capabilities. It does
not mean completing fid's entire historical roadmap or retaining role-, QA-,
template-, report-, or prose-heavy workflow features that do not fit the
approved product boundary.

Agent execution reduces Human keyboard time, but does not remove specification
uncertainty, integration risk, token and review cost, or the need for Human
judgment at protected decisions. Staged contract evidence is therefore still
required even if live migration waits for the full capability set.

This is a qualification and migration-threshold decision, not a requirement
that every workflow invoke all twelve capabilities. Independent capability
negotiation remains part of the port contract. The exact slice plan remains a
Gate D roadmap decision after Gate B settles the shared semantics.

## Decision 3 / Gate B1 - Policy And Enforcement Split

Status: **APPROVED by Human on 2026-08-08**

The kit evaluates engineering-workflow policy and compiles the result into an
approved, versioned Work Envelope. This includes risk classification,
authority, routing, budgets, protected decisions, and the selection of
obligations. fid does not independently recalculate those policy choices.

fid enforces generic operational integrity and the preconditions recorded in
the approved envelope. These include identity and expected-revision checks,
idempotency, claim/lease ownership, immutable history, authorization to approve
or revise an envelope, non-weakening of approved obligations, trusted-evidence
binding and freshness, and legal lifecycle changes. Shared conformance
fixtures prove the public kit/fid contract without creating two engineering
policy engines.

Consequently, fid will not infer that a particular file or change implies a
higher risk lane. If policy inputs change, the kit must promote or revise the
envelope; fid then prevents completion against the superseded revision and
enforces the strengthened recorded obligations. The exact authorization and
policy-provenance fields remain a joint-contract design detail.

## Decision 4 / Gate B2 - Trusted Gate Executor

Status: **APPROVED by Human on 2026-08-08**

fid is the primary trusted executor for local registered gates. Invocation
supplies an immutable Gate Definition Registration reference, never replacement
command fields. fid resolves the registration, records the execution-start
fact before process spawn, executes under its registered restrictions, observes
the terminal result, and records the immutable receipt. The kit receives and
references the receipt identity; an agent does not report or pass a trusted
result into fid.

Separately trusted external runners remain an extension for CI, Assured work,
and release gates through a versioned attestation contract. Reported agent or
Human command results remain queryable compatibility/diagnostic facts but
cannot satisfy deterministic trusted-evidence obligations. Live gate output is
an approved usability requirement; its exact interactive and machine framing
is a later interface-design detail.

## Decision 5 / Gate B3 - Gate Definition Authority Model

Status: **APPROVED by Human on 2026-08-08**

Clarification: `.fid/gates.toml` is itself committed, machine-readable,
project-owned configuration. The open question is not simply `fid-owned`
versus `project-owned`. It combines three concerns that must remain distinct:

1. Which namespace contains the authored definition.
2. Which product owns the versioned schema and validates it.
3. Whether the trusted executor reads mutable authoring configuration directly
   or resolves an explicitly authorized immutable registration.

`.fid/gates.toml` is not an original pre-V2 fid artifact. It first appeared in
fid commit `49fae77` at 17:23 on 2026-07-18 as part of the trusted-evidence
checkpoint, after the kit repository committed its V2 review and plan in
`ba706fe` at 15:03 that day. It belongs to the same convergence period, but
predates the precise joint ownership and immutable-registration model.

The current fid draft treats `.fid/gates.toml` as the live registry parsed by
fid at execution and snapshots its fields into receipts; Gate Definitions are
not stored as first-class registration entities. The V2 kit contract instead
requires an authoring source to be validated into an immutable, identity-bound
Gate Definition Registration that the executor resolves by reference. The
approved Gate B2 decision selects the latter runtime model regardless of the
eventual authoring-file location.

The kit owns the versioned Gate Definition schema, validation, and activation
procedure; the consuming project owns and approves the configured values in
machine-readable project configuration; fid owns the immutable registrations,
execution attempts, and receipts. `QUALITY_PLAN.md` may explain and reference
gate IDs, but is not parsed as executable configuration. fid never treats a
mutable authoring file as runtime authority and executes only an active
immutable registration.

The exact authoring filename and format are deliberately deferred to joint
interface design. Credible locations include a `gates` section in
`.agentic/project.yaml`, a separate project-overlay gate file, or a respecified
`.fid/gates.toml` used only as authoring input. A separate project-overlay gate
file is the current implementation preference because it keeps diffs isolated
without making policy appear backend-owned. Standalone clients use the generic
registration contract rather than a second live registry.

## Decision 6 / Gate B4 - fid-Native Work Identity

Status: **APPROVED by Human on 2026-08-08**

fid Issues and their parent/dependency relationships are the only durable work
graph. A Work Envelope is a versioned execution contract over one primary
Issue and, when explicitly selected, a bounded related Issue set. It does not
contain or create a second persistent Work Item identity, lifecycle, table, or
editable source of truth.

Durable decomposition becomes a task or subtask when it is independently
claimable, parallelizable, dependency-bearing, separately acceptable, or
expected to survive the current Run. Smaller procedural steps stay in the
agent's transient plan or a concise Run checkpoint. The harness derives a
noncanonical Next action from current Issue, Envelope, Obligation, Claim, Run,
Review, evidence, stop, and attention state.

Claims and Runs bind the Issue, Envelope revision, actor role, and optional
governing Obligation. Independent review is scheduled as a Review Run and
recorded as a Review Record satisfying a review Obligation, not as a Review
Work Item. The current kit V2-2 `WorkItem`/`WorkGraph` fake is superseded work
in progress and needs a correction slice; its historical phase evidence should
not be rewritten.

This approves the domain direction, not the exact transport or storage schema.
Before implementation, joint conformance scenarios must settle:

- parallel Builder claims over independently claimable Issues (resolved by
  Decision 7 / Gate B4a and ADR-0006);
- review findings that return work for rework or create follow-up work
  (resolved by Decision 9 / Gate B4c and ADR-0008);
- crash recovery during an active Run (resolved by Decision 8 / Gate B4b and
  ADR-0007); and
- Envelope revision and partial completion across multiple selected Issues.

Recorded in: `docs/adr/0005-issues-are-the-durable-work-graph.md`.

## Decision 7 / Gate B4a - Parallel Builder Isolation

Status: **APPROVED by Human on 2026-08-08**

One fid Issue may have at most one active mutating Builder Claim. A Work
Envelope is not claimed; each Claim binds one selected Issue, the Envelope
revision, the Builder, and the repository worktree identity. Parallel Builder
mutation is permitted only for distinct selected Issues in separate worktrees
with declared, provably disjoint mutation scopes. If isolation cannot be
established, the second Next action is non-actionable with typed reason
`mutation-conflict`.

Parallel work inside one Issue requires decomposition into independently
meaningful fid subtasks; procedural steps are not promoted merely to unlock
parallelism. Read-only Reviewer Runs may coexist, but each is bound to the
exact repository identity reviewed and becomes stale when that target changes.
The kit/harness establishes approved isolation; fid atomically enforces the
recorded Claim-conflict preconditions.

Recorded in: `docs/adr/0006-isolate-parallel-builder-mutations.md`.

## Decision 8 / Gate B4b - Interrupted-Run Recovery

Status: **APPROVED by Human on 2026-08-08**

Continuation after an agent-session or executor interruption always creates a
new Run linked to the prior Run. Run identity is never reused. The prior Run
becomes terminally `interrupted`, distinct from successful closure and
intentional abandonment; the Issue and Work Envelope remain active.

Takeover requires the prior executor to be confirmed terminated and its Claim
released, expired, or explicitly superseded through authorized force recovery.
Lease expiry alone is not proof that filesystem mutation stopped. Unknown
liveness returns `executor-liveness-unknown` and blocks mutation.

The kit/harness captures and reconciles the current worktree against the Run
start, last durable checkpoint, approved scope, and repository history before
the new Claim/Run may mutate. A new Builder may explicitly adopt in-scope
recovered work, but that work is not trusted evidence. Protected, unexpected,
unattributable, or unsafe state returns `recovery-state-conflict`. Recovery
never automatically resets, cleans, stashes, checks out, commits, or discards
Git state.

Obligations, concerns, evidence rules, and consumed budgets carry forward.
Receipts retain exact-target freshness semantics, and an execution start
without a terminal receipt satisfies nothing. Human authority is required only
for force takeover or ambiguous/protected repository state.

Recorded in: `docs/adr/0007-recover-interruptions-with-linked-runs.md`.

## Decision 9 / Gate B4c - Review Rework

Status: **APPROVED by Human on 2026-08-08**

A completed Review Run creates an immutable Review Record bound to the exact
Envelope revision and repository identity reviewed. A `changes-required`
record proves review occurred but leaves the review Obligation unsatisfied.
The original Issue and Envelope remain active; the harness derives
`builder-rework` rather than creating a Work Item or automatically opening a
second Issue.

An in-scope correction uses a new Builder Claim and Run on the same Issue and
Envelope revision. The Envelope is revised only when the correction changes
approved scope, a protected decision, dependency, risk, acceptance,
Obligation, or authority. A separate nonblocking improvement may become an
explicitly linked Issue, but blocking work cannot be deferred merely to pass
the current review.

The Builder may record that a Finding was addressed and request re-review, but
only a new independent Review Run and immutable Review Record can resolve it
or satisfy the review Obligation. Repository target changes stale prior
target-bound receipts and Review Records; the initial contract conservatively
re-runs the applicable whole gate/review profile. Rework consumes existing
budgets, and exhaustion returns `review-rework-exhausted` plus Attention.

A blocking Finding may otherwise cease to block only through an authoritative
inapplicability fact, an approved Envelope revision, or explicit Human risk
acceptance where project policy permits it.

Recorded in: `docs/adr/0008-keep-review-rework-on-the-original-issue.md`.

## Open Decision Order

1. Remaining Gate B4 lifecycle scenario: multi-Issue Envelope revision and
   partial completion.
2. Roadmap disposition: reslice/order and treatment of role-, QA-, template-,
   and report-oriented features.

The order may change if discussion shows that a downstream choice depends on
an unresolved semantic decision.
