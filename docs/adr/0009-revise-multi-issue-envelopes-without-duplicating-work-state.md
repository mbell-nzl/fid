---
status: accepted
date: 2026-08-08
---

# Revise multi-Issue Envelopes without duplicating work state

> **Superseded on 2026-11-10 by ADR-0014.** Envelopes are removed with the
> kernel, so multi-Issue Envelope revision and partial progress no longer
> apply. Issue state stays the only durable work lifecycle (ADR-0005). The
> text below is historical reasoning.

Progress across a Work Envelope's selected Issues is a derived projection of
their fid Issue states; it is not a `partially-completed` Envelope state,
membership lifecycle, or Work Item graph. Completion of one selected Issue is
durable and is not reversed merely because another selected Issue remains open
or the Envelope is later revised.

A selection, outcome, scope, authority, risk, acceptance, Obligation, budget,
or stop-condition change creates a new immutable Envelope revision. The prior
revision remains auditable and Runs, Claims, Reviews, and receipts retain their
original revision bindings. A new revision may activate only when the prior
revision has no active mutating Claim or Run; live work must finish, be
interrupted and recovered, or be explicitly stopped first.

Removing an unfinished selected Issue cannot silently weaken the approved
outcome. The revision records an authorized disposition and reason: a
successor Envelope reference or explicit deferral that leaves the fid Issue
open and visible, an authoritative inapplicability fact, or an approved outcome
change. The kit re-evaluates policy and obtains the authority required by the
material change; fid enforces the recorded revision and non-weakening
preconditions. Adding work selects an existing fid Issue and never creates an
Envelope-owned work identity.

The Envelope completes only when every Issue selected by its active revision
is in a satisfying terminal state and every applicable Obligation is current
and satisfied. Cancellation does not satisfy completion or a dependency.
Completed Issue state carries across revisions, while subsequent repository
changes may stale exact-target receipts and Review Records under their normal
freshness rules and require final Envelope-level proof to run again.

## Consequences

- Issue state remains the only durable work lifecycle; Envelope progress and
  selection history are projections over Issues plus immutable revisions.
- An unresolved removed Issue cannot disappear from recovery, reporting, or
  Human attention merely because the active selection changed.
- Shared conformance must cover partial progress, completed-Issue carry-forward,
  add/remove revision, unauthorized weakening, missing removal disposition,
  live-Claim/Run revision conflict, cancellation, dependency behavior,
  exact-target staleness, completion guards, idempotent retry, and atomic
  failure.
- This resolves the final lifecycle scenario required by ADR-0005.
