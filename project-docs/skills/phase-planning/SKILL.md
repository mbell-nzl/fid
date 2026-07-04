---
name: phase-planning
description: Use for Tech Lead phase kickoff including operating-mode selection, carry-forward disposition, validation contract drafting, task slicing, and phase brief refresh.
---

# Phase Planning

Provenance: Workflow-Extracted Skill.

Use this skill when planning a new approved phase. The Tech Lead is the
primary user; other roles may reference it for understanding the expected
phase shape.

## Procedure

1. Read the roadmap entry for the target phase and relevant spec sections
   (PRD, TRD, design, quality).
2. Review carry-forward items in `project-docs/docs/current/decisions_active.md`.
   For each item: schedule in this phase, re-defer with rationale, or close as
   satisfied/obsolete.
3. Skim prior phase workflow notes when available. Treat as calibration only
   unless promoted to a current decision.
4. Propose the phase operating mode (light, normal, or heavy) per
   `project-docs/specs/09_AI_WORKFLOW.md` section 4. Record it in the phase
   brief; light mode skips the validation contract and QA plan per the
   artifact matrix.
5. Create phase directories:
   - `project-docs/docs/phases/phase-N/tasks/`
   - `project-docs/docs/phases/phase-N/qa/`
6. For normal and heavy modes, draft
   `project-docs/docs/phases/phase-N/qa/validation-contract.md` from
   `project-docs/templates/validation-contract-template.md`. Express each
   in-scope requirement as one observable `VAL-NNN` assertion. Conditional or
   binary assertions must name both states; parity or "matches reference"
   assertions need deterministic conformance checks where possible.
7. Slice work using the one-logic rule: logic/data tasks before UI/wiring,
   schema and dependency changes isolated unless the phase explicitly plans
   the exception.
8. Write task specs from `project-docs/templates/task-spec-template.md`.
   Each task must cite covered `VAL-NNN` IDs or state why none apply, and its
   `## Exact Commands To Run` section must list the full task gate before task
   close unless the spec records an explicit exception.
9. Refresh `project-docs/docs/current/phase_brief.md` from the template.
   Keep it as an execution index under 150 lines, not a duplicate of full
   task specs. Include the operating mode and validation contract index when
   required.
10. Present the operating mode, validation contract when required, and ordered
    task list to the human.

## Stop Conditions

Stop if required roadmap, spec, or current-state artifacts are missing.

Stop if the requested phase boundary mixes unrelated schema, logic, UI,
dependency, and workflow changes without an explicit exception.

## References

- `project-docs/specs/09_AI_WORKFLOW.md`
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md`
- `project-docs/templates/validation-contract-template.md`
- `project-docs/templates/task-spec-template.md`
- `project-docs/templates/phase-brief-template.md`
