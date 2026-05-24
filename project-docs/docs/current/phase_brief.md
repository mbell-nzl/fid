# Phase Brief - No Active Phase

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Replaced at each phase kickoff by the incoming Tech Lead.
>
> **Size target:** under 150 lines.

## Phase Goal

No implementation phase is active yet.

The human is progressively filling the seven spec files from
`scaffold/initial-spec.md`. Once specs and quality commands are concrete,
invoke Tech Lead to plan Phase 0 from
`project-docs/specs/06_IMPLEMENTATION_ROADMAP.md`.

## Status

Pre-phase scaffolding. Specs are being filled progressively.

## Pre-Scaffold Material

- `scaffold/initial-spec.md` — 2600-line design spec; the input for filling
  the seven specs.
- `docs/adr/0001-0003` — three binding design decisions predating the phase
  system.
- `project-docs/specs/GLOSSARY.md` — domain language already populated from
  earlier design work.

## Active Constraints

- Do not implement stack-dependent code while `<TBD>` stack choices remain
  in `AGENTS.md`.
- Do not run autonomous implementation while quality commands are placeholders.
- Preserve the artifact workflow unless a recorded decision changes it.

## Open Questions

- [ ] Fill remaining `<TBD>` stack items (framework, test tools, build runner).
- [ ] Fill quality commands in `AGENTS.md`.
- [ ] Split `scaffold/initial-spec.md` into the seven spec files.
- [ ] Plan Phase 0 with Tech Lead.

## Next Planning Checkpoint

Phase 0 planning begins after specs and quality commands are concrete.
