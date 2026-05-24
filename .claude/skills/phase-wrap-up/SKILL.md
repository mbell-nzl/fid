---
name: phase-wrap-up
description: Use at phase end when the Tech Lead audits validation outcomes, syncs specs, records decisions, and compresses current context.
---

<!-- Generated from project-docs/skills/phase-wrap-up/SKILL.md by project-docs/tools/sync-skills.sh. -->
<!-- Do not edit directly; update the canonical skill and rerun the sync command. -->

# Phase Wrap-Up

Provenance: Workflow-Extracted Skill.

Use this skill for the phase-end Tech Lead procedure after QA review and human
verification have produced their artifacts.

## Procedure

1. Load the final review, human checklist outcome, validation contract, task
   completion reports, phase decisions, triage log, escalations, and current
   context files for the phase.
2. Audit every `VAL-NNN` final outcome. Disposition failed, blocked, partial,
   or carry-forward assertions through the canonical phase decisions and
   roadmap artifacts.
3. Walk every in-phase decision, Developer deviation, validation-contract
   deviation, and QA recommendation marked for immediate action. For each spec,
   record whether it was touched, verified unchanged, or not applicable.
4. Audit all `TRG-NNN` entries and update warm-up or autonomous-resume notes
   according to the workflow spec.
5. Append workflow notes and a deep-retrospective recommendation to the phase
   decisions file.
6. Check whether phase work changed project language. Update
   `project-docs/specs/GLOSSARY.md` or flag ambiguity.
7. Move resolved escalations into the phase archive, refresh active decisions
   and progress, and trim default-loaded current docs to their budgets.

## References

- `project-docs/specs/05_AI_WORKFLOW.md`
- `project-docs/specs/04_QUALITY_PLAN.md`
- `project-docs/templates/phase-decisions-template.md`
- `project-docs/templates/phase-brief-template.md`
- `project-docs/templates/final-review-template.md`
