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
- Exact Commands To Run
- Acceptance criteria
- Output requirements (completion report path)

## Exact Command Guidance

Use the project's configured commands from `AGENTS.md`, plus any task-specific
commands required by the QA plan. The full task gate is required before task
close unless the task spec explicitly records an exception. Focused commands
used during red/green work do not replace the required gate.

Every required command must be reported as `PASS`, or `BLOCKED` with a reason
in a blocker report. Missing, silently skipped, or narrower substituted
commands are failed gates. Stop and report if any required command is still a
placeholder.

## Output

- Ordered task list
- Task spec paths created
- Notes on intentionally deferred work
- Risks and open questions

## References

- `project-docs/templates/task-spec-template.md`
- `project-docs/templates/completion-report-template.md`
- `project-docs/specs/09_AI_WORKFLOW.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
