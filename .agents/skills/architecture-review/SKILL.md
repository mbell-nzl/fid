---
name: architecture-review
description: Use to find deepening opportunities in the codebase, informed by the glossary and recorded decisions. Surfaces shallow modules, missing seams, and testability friction.
---

<!-- Generated from project-docs/skills/architecture-review/SKILL.md by project-docs/tools/sync-skills.sh. -->
<!-- Do not edit directly; update the canonical skill and rerun the sync command. -->

# Architecture Review

Provenance: External-Adapted Skill. Based on Matt Pocock's
improve-codebase-architecture skill, adapted for this starter kit's spec
system and review output conventions.

Use this skill to audit architecture quality, find refactoring
opportunities, or assess module depth before a phase that touches
cross-cutting structure.

## Key Vocabulary

Use these terms consistently:

- **Module** — anything with an interface and an implementation.
- **Interface** — everything a caller must know: types, invariants, error
  modes, ordering, config. Not just the type signature.
- **Depth** — leverage at the interface: significant behaviour behind a
  small interface. Deep = high leverage. Shallow = interface nearly as
  complex as the implementation.
- **Seam** — where an interface lives; a place behaviour can be altered
  without editing in place.
- **Adapter** — a concrete thing satisfying an interface at a seam.
- **Locality** — change, bugs, and knowledge concentrated in one place.
- **Deletion test** — imagine deleting the module. If complexity vanishes,
  it was a pass-through. If complexity reappears across callers, it was
  earning its keep.

## Procedure

1. Read `project-docs/specs/GLOSSARY.md`, relevant specs, and any recorded
   decisions (`DEC-NNN`) in the area under review.
2. Explore the codebase organically. Note where understanding one concept
   requires bouncing between many small modules, where modules are shallow,
   where pure-function extraction hides the real complexity, and where
   tightly-coupled modules leak across seams.
3. Apply the deletion test to anything suspected of being shallow.
4. For each candidate, document:
   - **Files** — which files/modules are involved.
   - **Problem** — why the current architecture causes friction.
   - **Solution** — plain English description of what would change.
   - **Benefits** — explained in terms of locality, leverage, and test
     improvement.
   - **Recommendation strength** — `Strong`, `Worth exploring`, or
     `Speculative`.
5. If a candidate contradicts a recorded decision, surface it only when the
   friction is real enough to warrant revisiting. Mark the conflict clearly.
6. Write the review output to
   `project-docs/docs/reviews/<date>-architecture-review.md`.

## Integration With Project Decisions

- Use `GLOSSARY.md` vocabulary for domain concepts.
- Use the vocabulary above for architecture concepts.
- Do not create a parallel decision system. Record architecture decisions
  through the project's `DEC-NNN` system.
- Do not propose interfaces during the review. Present candidates first and
  let the user choose which to explore further.

## References

- `project-docs/specs/GLOSSARY.md`
- `project-docs/specs/02_TRD.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/templates/phase-decisions-template.md`
