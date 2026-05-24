---
name: phase-planning
description: Use for Tech Lead phase kickoff including carry-forward disposition, validation contract drafting, task slicing, and phase brief refresh.
---

<!-- Generated from project-docs/skills/phase-planning/SKILL.md by project-docs/tools/sync-skills.sh. -->
<!-- Do not edit directly; update the canonical skill and rerun the sync command. -->

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
4. Create phase directories:
   - `project-docs/docs/phases/phase-N/tasks/`
   - `project-docs/docs/phases/phase-N/qa/`
5. Draft `project-docs/docs/phases/phase-N/qa/validation-contract.md` from
   `project-docs/templates/validation-contract-template.md`. Express each
   in-scope requirement as one observable `VAL-NNN` assertion.
6. Slice work using the one-logic rule: logic/data tasks before UI/wiring,
   schema and dependency changes isolated unless the phase explicitly plans
   the exception.
7. Write task specs from `project-docs/templates/task-spec-template.md`.
   Each task must cite covered `VAL-NNN` IDs or state why none apply.
8. Refresh `project-docs/docs/current/phase_brief.md` from the template.
   Keep it as an execution index under 150 lines, not a duplicate of full
   task specs.
9. Present the validation contract and ordered task list to the human.

## Stop Conditions

Stop if required roadmap, spec, or current-state artifacts are missing.

Stop if the requested phase boundary mixes unrelated schema, logic, UI,
dependency, and workflow changes without an explicit exception.

## References

- `project-docs/specs/05_AI_WORKFLOW.md`
- `project-docs/specs/06_IMPLEMENTATION_ROADMAP.md`
- `project-docs/templates/validation-contract-template.md`
- `project-docs/templates/task-spec-template.md`
- `project-docs/templates/phase-brief-template.md`
