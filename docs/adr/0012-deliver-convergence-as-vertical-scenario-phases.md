---
status: accepted
date: 2026-08-09
---

# Deliver convergence as vertical scenario phases

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

## Consequences

- architecture is still developed in dependency order inside a phase, but
  progress is accepted only through a coherent scenario;
- cross-package phases are expected and require several bounded Developer
  tasks rather than one oversized task;
- shared fake/fid conformance, rebuild, and removal checks enter as soon as
  the owning scenario needs them; and
- the roadmap records both the forward outcome and the displaced V1 surface
  for every phase.
