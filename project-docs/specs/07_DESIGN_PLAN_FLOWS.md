# END-TO-END FLOWS - fid

## 1. Standalone Issue Inventory

1. Human initializes fid and establishes project identity.
2. Human creates Epics/Tasks/Subtasks, acceptance criteria, dependencies,
   labels, and workflow assignments.
3. fid attributes each mutation, validates hierarchy/workflow invariants,
   appends one Event, and updates the projection.
4. `fid ready` derives actionable Issues and typed reasons.

This flow claims no protected Human Outcome or trusted execution evidence. It
keeps fid useful as an Issue tracker without requiring the V2 kit.

## 2. Kit Activation

1. Kit discovers backend capabilities and validates required ranges.
2. Kit prepares one activation naming old/new kit, overlay, provider, gate
   authoring, and backend identities.
3. Each participant stages and validates its component under the activation ID.
4. Before the point of no return, failure aborts and restores the old active
   set. After irreversible backend commit, recovery completes forward.
5. fid permits normal kit work only after every required component agrees on
   the authoritative active version.

## 3. Work Envelope Activation

1. Orchestrator selects one primary existing Issue and any explicitly bounded
   related Issues.
2. Kit evaluates risk, roles/routing, authority, protected Decisions, gates,
   validation/review/Human Obligations, budgets, Stops, and repository scope.
3. Human approval is obtained where kit policy requires it.
4. fid creates immutable Gate Definition Registrations from the activated Gate
   Definition values that the Envelope will reference.
5. Kit submits the complete versioned Envelope through the backend port.
6. fid validates references, capability support, provenance, non-weakening, and
   revision invariants, then records the approved/active revision.

fid does not persist or recalculate the external kit's risk lane and does not
create Issue copies. A caller cannot substitute a mutable Gate Definition for
the activated Registration reference.

## 4. Builder Execution

1. Harness requests Next Actions for Actor capabilities.
2. fid derives actionable state from Issue, Envelope, authority, dependencies,
   Obligations, Claim/Run, evidence, Stops, and Attention.
3. Harness acquires a mutation Claim for one selected Issue, exact Envelope
   revision, worktree, and declared scope.
4. fid rejects overlapping or unknown parallel mutation with
   `mutation-conflict`.
5. Harness starts a Run and records observed telemetry incrementally.
6. Builder changes the repository within approved scope and records Concerns or
   Deviations only when they govern future action.
7. Transition evaluation returns all unsatisfied preconditions; apply either
   commits the complete transition Event or nothing.

## 5. Registered Gate Execution

1. Project-owned Gate Definition is validated by the kit and activated into an
   immutable fid Registration.
2. Harness requests execution by Registration reference only.
3. fid resolves the exact active definition/binding and appends the Attempt
   before process spawn.
4. fid executes structured argv under registered restrictions and streams live
   output.
5. On observed termination, fid writes one immutable Receipt with exact target,
   toolchain, status, sanitized output identity, and retention state.
6. Evidence links the Receipt to applicable Obligations.
7. A crash, timeout, cancellation, harness error, or missing Receipt cannot
   satisfy the Obligation.

## 6. Interruption Recovery

1. A Run loses its agent session or executor before normal terminalization.
2. Harness confirms prior-executor termination or obtains force-recovery
   authority; lease expiry alone is insufficient.
3. Repository port reports worktree identity and changes against the last
   durable checkpoint.
4. Unknown liveness returns `executor-liveness-unknown`; unsafe or ambiguous
   state returns `recovery-state-conflict`.
5. On success fid atomically terminalizes the prior Run as interrupted,
   disposes its Claim, records reconciliation/adoption, creates the linked new
   Claim/Run, and preserves Issues, Envelope, Obligations, evidence policy, and
   consumed budgets.

Recovered code is not trusted evidence; gates/review run normally.

## 7. Review And Rework

1. A fresh Reviewer Run binds the exact Envelope revision and repository
   identity.
2. Reviewer records coverage, verdict, evidence, and zero or more Findings.
3. `pass` may satisfy the review Obligation when every other condition holds.
4. `changes-required` leaves it unsatisfied and derives Builder rework on the
   original Issue/Envelope.
5. Builder acquires a new Claim/Run and may mark Findings addressed.
6. Only a later qualifying independent Review Record resolves a Finding.
7. Target changes stale bound proof; rework consumes existing budgets.
8. Exhaustion creates `review-rework-exhausted` and Attention.

No QA Work Item, QA Plan, handoff, or completion report is created.

## 8. Protected Decision And Human Outcome

1. Execution encounters a protected choice or Human-only check.
2. fid records a typed Stop and deduplicated Attention Item.
3. Orchestrator presents a Human Validation Brief derived from current facts.
4. Human records project Decision prose in the canonical project file where
   applicable.
5. Kit submits the Decision Reference and provenanced Human Outcome.
6. fid validates active authority, resolves matching Attention through the
   Outcome, and permits only the newly authorized action.

Routine status never raises Attention.

## 9. Envelope Revision And Multi-Issue Progress

1. Progress is read from selected Issue states; it is never written as partial
   completion.
2. Material selection/contract change proposes a new immutable revision.
3. Active mutation against the prior revision blocks activation.
4. Removing unfinished work requires a visible authorized disposition; adding
   work selects an existing Issue.
5. Completed Issue state remains complete, while target changes may stale
   bound evidence and Review Records.
6. Completion succeeds only when every active-revision Issue is satisfying
   terminal and every applicable Obligation is current and satisfied.

## 10. Rebuild And Portable Recovery

1. Operator creates a versioned Portable Recovery Export or deletes the SQLite
   projection.
2. Validation checks schema, hashes, referential integrity, Event order,
   immutable record uniqueness, and capability compatibility without mutation.
3. Rebuild/import replays ordered Events through the same appliers used live.
4. Current Issues, Envelopes, Claims/Runs, Stops/Attention, registrations,
   Receipts/Evidence, Reviews/Findings, Outcomes, activation, and telemetry
   match the source semantics.
5. Projections are regenerated; project prose is not copied.

Phase 21 implements deterministic export creation and non-mutating validation
only. Import, identity/history preservation, and full live/rebuild/import
equivalence enter with Phase 22 recovery; `export.portable.v1` remains
unavailable until that parity passes.

## 11. Pre-Dogfood Convergence

1. A D3 slice specifies its JSON contracts, DDL, Event fixtures, conformance
   scenarios, and displaced V1 surface.
2. Replacement behavior is implemented under red/green TDD.
3. Displaced commands, tables, packages, formatters, and tests are removed in
   the same or an explicitly dependent buildable task.
4. No compatibility shim is added.
5. Final convergence replaces the development migration chain with the first
   supported schema/Event baseline and rehearses a disposable project.

*Last updated: 2026-08-09 (ADR-0011)*
