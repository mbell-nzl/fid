---
status: accepted
date: 2026-08-08
---

# Keep review rework on the original Issue

A completed Review Run produces an immutable Review Record bound to the exact
Work Envelope revision and repository identity reviewed. A `changes-required`
record proves that review occurred but does not satisfy the governing review
Obligation. Its Findings remain durable and the original Issue and Envelope
remain active with derived Next action `builder-rework`; fid creates neither a
Work Item nor a follow-up Issue automatically.

Ordinary correction of an in-scope Finding uses a new Builder Claim and Run on
the same Issue and Envelope revision. A Finding does not by itself expand
authority. The Envelope must be revised when correction changes approved
scope, a protected decision, a dependency, risk, acceptance, an Obligation, or
authority. A separate nonblocking improvement may instead become an explicitly
linked Issue, but blocking work cannot be deferred merely to satisfy the
current review Obligation.

Only a new independent Review Run and immutable Review Record can resolve a
Finding or satisfy the review Obligation. A Builder may record that a Finding
was addressed and request re-review, but cannot mark it resolved. When the
repository target changes, target-bound receipts and Review Records become
stale; the initial contract conservatively re-runs the whole applicable gate
and review profile rather than attempting path-aware reuse.

Review-rework cycles consume the Envelope's existing attempt, time, token,
spend, and review-cycle budgets. Exhaustion makes the Next action
non-actionable with `review-rework-exhausted` and creates Attention. A blocking
Finding may cease to block only through a resolving Review Record, an
authoritative inapplicability fact, an approved Envelope revision, or explicit
Human risk acceptance where project policy permits it.

## Consequences

- Review evidence stays immutable without duplicating the durable work graph.
- The kit/harness routes rework and selects policy; fid enforces actor,
  identity, freshness, budget, and transition preconditions.
- Shared conformance must cover changes-required review, same-Issue rework,
  stale target proof, re-review resolution, forbidden Builder self-resolution,
  blocking deferral, legitimate linked follow-up work, Envelope-revision
  triggers, budget exhaustion, idempotent retry, and atomic failure.
- This resolves the review-rework scenario required by ADR-0005. ADR-0009
  resolves multi-Issue Envelope revision and partial progress.
