---
status: accepted
date: 2026-08-09
---

# Deliver convergence as vertical scenario phases

> **Superseded on 2026-10-31 by ADR-0013.** No phases and no convergence
> with the V2 kit. `RESET.md` tasks ship in small PRs. The text below is
> historical reasoning.

The first convergence proposal grouped work mainly by technical dependency:
machine contracts, execution state, evidence, review, and portability. That
order made architectural dependencies visible, but it delayed proof that the
parts formed a usable workflow and risked producing several individually
complete horizontal layers before any representative journey worked.

fid will instead organize convergence primarily as phase-sized vertical
scenario slices. A phase may contain several Developer tasks and may include
horizontal schema, event, adapter, migration, and deletion work, but the phase
acceptance boundary is an end-to-end user or agent journey.

The approved order is:

1. a clean standalone Issue journey;
2. a thin one-Issue Routine walking skeleton spanning activation, authority,
   execution, one registered gate, trusted evidence, completion, resume, and
   portable export;
3. failure and recovery;
4. independent review and Human judgment;
5. parallel and evolving work; and
6. full qualification and clean-baseline convergence.

The first slice becomes a newly planned Phase 20. The earlier Phase 20 role-
policy proposal remains superseded; no `phase-20/` artifact directory was
created for it, so reusing the number does not rewrite a historical phase
record.

Each phase owns explicit legacy removals. A surface with no retained purpose
may leave in the first phase that reaches it. A safety-relevant V1 surface
remains fail-closed until its generic replacement passes the owning phase's
tests. The final convergence phase is not a holding area for known legacy
code; it removes only cross-slice scaffolding, superseded migrations and
fixtures, and residue discovered by whole-product qualification.

No partial slice is a live-migration point. A capability is advertised only
when its mandatory contract and conformance cases pass, even when an earlier
walking skeleton has introduced some of its internal model.

fid will continue using the repository's V1 phase, role, task, validation, and
evidence process to build these slices until the V2 kit is ready. That process
is delivery scaffolding, not fid product scope. The Tech Lead may now plan the
new Phase 20 from the accepted roadmap. The exact one-Issue Routine walking-
skeleton boundary is the next product discussion before Phase 21 is planned.

## Gate D3.2 amendment - 2026-08-22

The Human accepted the one-Issue walking-skeleton boundary as one **Heavy
Phase 21**, not formal 21A/21B phases. Phase 21 uses two internal task groups
with a fresh independent review and executor-feasibility checkpoint after the
authority/activation/Gate Registration/Work Envelope/Claim/Run group and before
trusted execution/evidence. This is internal sequencing, not a new acceptance
boundary; Phase 21 still closes only through the complete vertical journey.

The accepted thinness decisions are:

- Attention, typed Stops, budget consumption/exhaustion, interruption
  recovery, automatic respawn, activation recovery, and export import belong
  to Phase 22;
- Phase 21 stores immutable Envelope/Run budget contracts with consumption and
  remaining unavailable (`null`);
- Phase 21 supports one complete generic Obligation value shape and one type,
  deterministic gate evidence from fid-executed Receipts;
- executor policy is limited to standard-library-enforceable controls and
  rejects unsupported restrictive policy rather than claiming OS containment;
- export create and non-mutating validate land, but `export.portable.v1`
  remains unavailable until import identity/history parity passes; and
- every incomplete semantic capability remains unavailable in discovery.

`Routine` is an external kit risk-lane term. fid neither persists nor
evaluates it; fid receives only the activated, approved operational values
compiled from kit policy.

## Phase 21 terminal-gate amendment - 2026-08-22

The one gate-evidence Obligation guards the selected Issue's terminal-success
transition before it also guards Work Envelope completion. The representative
journey therefore executes the registered Gate, records and links its Receipt,
and obtains a satisfied evaluation before terminal Issue mutation. Completion
then independently re-evaluates the same current facts and requires the
terminal-success Issue. One shared evaluator prevents the two guards from
becoming separate authorities.

The post-rebuild half is a parity proof, not a second journey mutation: rebuild
replays the one set of durable facts, then repeats reads, evaluation,
completion state, resume, export, and negative conformance. It does not create
a second one-revision Envelope in the project.

## Consequences

- architecture is still developed in dependency order inside a phase, but
  progress is accepted only through a coherent scenario;
- cross-package phases are expected and require several bounded Developer
  tasks rather than one oversized task;
- shared fake/fid conformance, rebuild, and removal checks enter as soon as
  the owning scenario needs them; and
- the roadmap records both the forward outcome and the displaced V1 surface
  for every phase.
