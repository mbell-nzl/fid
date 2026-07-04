---
name: task-pack
description: Use to convert a phase objective or feature slice into ordered Developer task specs with validation contract coverage and quality gate commands.
---

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
- Every task must state explicit non-goals under `### Out Of Scope`.

## Required Task Headings

Use `project-docs/templates/task-spec-template.md`. Each task must include the
template's headings exactly as written there:

- `## Objective`
- `## Validation Contract Coverage` (`VAL-NNN` IDs or explicit no-VAL
  rationale; contract rows carry the source requirement references)
- `## Binding Constraints` (interpretation decisions, ordering rules, and
  contract semantics that bind the implementation; omit only for trivial
  tasks)
- `## Scope`, containing:
  - `### Files To Create Or Modify`
  - `### Files To Read For Context`
  - `### Out Of Scope`
- `## Acceptance Criteria`
- `## First Failing Test`
- `## Documentation Updates`
- `## Exact Commands To Run`
- `## Notes For Developer` (constraints, ordering, dependencies on prior
  tasks, escalation tripwires, and the completion report path when it deviates
  from the `task-M-done.md` convention)

Scenario, parity, and audit obligations are declared in the phase brief and
QA plan; add task-spec sections for them only when the phase brief requires it
(parity or heavy-mode phases).

If a heading here stops matching the template, fix the template or this list
in the same change; `workflow-lint.sh` checks this list against the template.

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
