# Progress - Current

> In-flight status only. Historical implementation progress remains in the
> phase archives.

## Project state

**No implementation phase is active.** Phase 19 completed and wrapped on
2026-08-08. fid has no adopted project workspace or compatibility constituency.

Gates D1/D2 supersede the old Phase 20 role-policy proposal and the historical
process-shaped roadmap. Gate D3.1 establishes a new Phase 20 as the first of
six vertical scenario slices. The target product remains the standalone Issue
tracker plus the fid/V2 durable operational kernel described by ADR-0011.

Branch: `main`.

## Gate D2 documentation

Completed and committed:

- product decision ADR, decision-session ledger, and reconciliation report;
- charter, PRD, and glossary rewrite;
- technical, logical-schema, and runtime-configuration rewrite;
- CLI, Human-interface, and end-to-end flow rewrite; and
- quality, V2 workflow-integration, and convergence-roadmap rewrite.

The rewrite removes forward promises for global role policy, role composites,
QA Plans, Validation Contracts, templates, skills, handoffs, routine
retrospectives, generic TOON/YAML/Markdown output, and pre-dogfood
compatibility. JSON is the sole canonical machine representation.

## Implementation baseline

Phase 19's Issue workflow authoring implementation remains useful foundation,
not a frozen public contract. Retain useful Issue, journal, SQLite, rebuild,
locking, identity, and repository primitives only where they conform to the
new specs.

Product code has not yet been removed or reshaped under ADR-0011. ADR-0012 now
assigns replacement and removal to phase-sized vertical scenarios beginning
with the standalone Issue journey.

## Carry-forward

The active disposition is in `decisions_active.md`. In summary:

- cancel `internal/dev` refactoring and generic Markdown formatter work;
- keep Issue-path and rebuild-performance debt only when their retained paths
  reopen;
- treat workflow inspection/lifecycle as candidates, not commitments; and
- keep V1 process-tooling concerns separate from product scope.

## Current work

1. Invoke the V1 Tech Lead to plan the new Phase 20 standalone Issue slice.
2. Approve its exact JSON contracts, removal inventory, validation contract,
   Developer tasks, and quality commands before implementation.
3. Resolve Gate D3.2's exact Routine walking-skeleton boundary before Phase 21
   planning.

No implementation task should cite the superseded historical version roadmap
as forward authority.
