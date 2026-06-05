---
name: tdd
description: Use for implementation tasks that require red/green TDD from a task spec, especially when tests, validation assertions, and completion-report evidence must line up.
---

<!-- Generated from project-docs/skills/tdd/SKILL.md by project-docs/tools/sync-skills.sh. -->
<!-- Do not edit directly; update the canonical skill and rerun the sync command. -->

# TDD

Provenance: External-Adapted Skill. Based on Matt Pocock's TDD skill and
adapted for task specs, validation contracts, and completion reports.

Use this skill for behavior-changing implementation work. Documentation-only
or exploratory tasks may use the task spec's explicit exception instead.

## Procedure

1. Read the task spec, named `VAL-NNN` assertions, relevant requirements,
   existing tests, and the public interface where behavior should be observed.
2. Choose the thinnest behavior slice that proves progress toward the task.
   Prefer a vertical tracer through the real boundary over a broad batch of
   isolated internals.
3. Write one failing test first. Assert externally visible behavior or a public
   contract, not incidental implementation details.
4. Run the narrowest useful test command and preserve the failing result as
   evidence for the completion report.
5. Implement the smallest change that makes the test pass. Keep business rules
   in the existing core/domain boundary and adapters thin.
6. Refactor only with tests green, then repeat red/green for remaining
   acceptance criteria.
7. Run every command in the task spec's `## Exact Commands To Run` section.
   The full task gate is required before close unless the task spec explicitly
   records an exception; focused red/green commands are additional evidence,
   not substitutes.
8. In the completion report, record each required command as `PASS`, or
   `BLOCKED` with a reason in a blocker report. Missing, skipped, or narrower
   substituted commands are failed gates.
9. Copy commit hashes and subjects from git output when reporting them; never
   reconstruct hashes from memory.
10. In the completion report, connect tests and gate evidence to the covered
   `VAL-NNN` IDs or explain why no validation assertion applies.

## References

- `project-docs/templates/task-spec-template.md`
- `project-docs/templates/completion-report-template.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/09_AI_WORKFLOW.md`
