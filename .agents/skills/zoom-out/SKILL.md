---
name: zoom-out
description: Use before modifying unfamiliar or cross-cutting code to map requirements, terminology, callers, module boundaries, and likely blast radius.
---

<!-- Generated from project-docs/skills/zoom-out/SKILL.md by project-docs/tools/sync-skills.sh. -->
<!-- Do not edit directly; update the canonical skill and rerun the sync command. -->

# Zoom Out

Provenance: External-Adapted Skill. Based on Matt Pocock's zoom-out skill.

Use this skill before making a change whose scope, terminology, or ownership
boundary is not yet clear.

## Procedure

1. Read the task or question, then identify the relevant spec sections,
   glossary terms, active decisions, and validation assertions.
2. Map the behavior path at a high level: user-facing adapter, core/domain
   logic, persistence or external integrations, shared utilities, tests, and
   docs.
3. Search for callers, neighboring tests, templates, prompts, and generated
   artifacts that may need to move with the change.
4. Separate the required change from adjacent improvements. Name explicit
   non-goals so the implementation stays scoped.
5. Note public-interface, data-model, dependency, workflow, or terminology
   questions that require a recorded decision or escalation.
6. Pick the smallest implementation route that preserves existing boundaries
   and gives clear verification evidence.

## References

- `project-docs/specs/GLOSSARY.md`
- `project-docs/specs/02_TRD.md`
- `project-docs/specs/05_AI_WORKFLOW.md`
- `project-docs/templates/task-spec-template.md`
