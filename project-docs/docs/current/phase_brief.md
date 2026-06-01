# Phase Brief - No Active Phase

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Replaced at each phase kickoff by the incoming Tech Lead.
>
> **Size target:** under 150 lines.

## Phase Goal

No implementation phase is active yet.

The human is progressively filling the canonical spec suite from
`scaffold/initial-spec.md`. Once remaining final spec audit blockers are
resolved, invoke Tech Lead to plan Phase 0 from
`project-docs/specs/10_IMPLEMENTATION_ROADMAP.md`.

## Status

Pre-phase scaffolding. Specs are being filled progressively.

## Pre-Scaffold Material

- `scaffold/initial-spec.md` — 2600-line design spec; the input for filling
  the canonical spec suite.
- `docs/adr/0001-0003` — three binding design decisions predating the phase
  system.
- `project-docs/specs/GLOSSARY.md` — domain language already populated from
  earlier design work.

## Active Constraints

- Do not start implementation while the final spec audit has open BLOCKER
  findings.
- Preserve the artifact workflow unless a recorded decision changes it.

## Open Questions

- [ ] Resolve remaining final spec audit blockers.
- [ ] Split `scaffold/initial-spec.md` into the canonical spec suite.
- [ ] Plan Phase 0 with Tech Lead.

## Next Planning Checkpoint

Phase 0 planning begins after remaining final spec audit blockers are resolved.
