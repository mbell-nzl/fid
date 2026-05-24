---
name: task-pack
description: Use to convert a phase objective or feature slice into ordered Developer task specs with validation contract coverage and quality gate commands.
---

<!-- Generated from project-docs/skills/task-pack/SKILL.md by project-docs/tools/sync-skills.sh. -->
<!-- Do not edit directly; update the canonical skill and rerun the sync command. -->

# Task Pack

Provenance: Workflow-Extracted Skill. Based on yoyu task-pack extraction
and adapted for stack-neutral starter kit use.

Use this skill when writing task specs for a planned phase. The Tech Lead
is the primary user during phase planning.

## Task Slicing Rules

- One Developer session implements exactly one task.
- Prefer focused diffs over broad rewrites.
- Logic/data tasks come before UI/wiring tasks.
- Schema, dependency, platform config, and public interface changes stay
  isolated unless the phase spec explicitly plans the exception.
- Every behavior change needs paired test work.
- Every task must list files not to modify.

## Required Task Headings

Use `project-docs/templates/task-spec-template.md`. Each task must include:

- Objective
- Requirements (PRD/TRD references)
- Validation Contract Coverage (`VAL-NNN` IDs or explicit no-VAL rationale)
- Files to create or modify
- Documentation to update
- Files not to modify
- Dependencies on prior tasks
- Implementation notes
- Exact quality gate commands
- Acceptance criteria
- Output requirements (completion report path)

## Quality Gate Guidance

Use the project's configured commands from `AGENTS.md`. If the project has
task-specific commands, prefer the narrowest scope that covers the task.
Stop and report if any required command is still a placeholder.

## Output

- Ordered task list
- Task spec paths created
- Notes on intentionally deferred work
- Risks and open questions

## References

- `project-docs/templates/task-spec-template.md`
- `project-docs/templates/completion-report-template.md`
- `project-docs/specs/05_AI_WORKFLOW.md`
- `project-docs/specs/04_QUALITY_PLAN.md`
