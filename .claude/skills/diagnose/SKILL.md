---
name: diagnose
description: Use when debugging a failing test, quality gate, bug report, or unexplained behavior before changing code.
---

<!-- Generated from project-docs/skills/diagnose/SKILL.md by project-docs/tools/sync-skills.sh. -->
<!-- Do not edit directly; update the canonical skill and rerun the sync command. -->

# Diagnose

Provenance: External-Adapted Skill. Based on Matt Pocock's diagnose skill.

Use this skill to understand a failure before choosing a fix.

## Procedure

1. Reproduce the failure with the smallest command, input, or scenario that
   still shows the problem. Capture the exact error or observed behavior.
2. State the expected behavior from the task spec, validation contract, spec,
   test, or existing user-facing contract.
3. Minimize the search area: identify the changed files, call path, data flow,
   configuration, and external boundary involved.
4. Form a short ranked hypothesis list. For each hypothesis, name what would
   prove or disprove it.
5. Inspect or instrument narrowly. Prefer logs, assertions, focused tests,
   traces, or file-state checks that can be removed after diagnosis.
6. Change one causal factor at a time. When the failure disappears, add or
   update a regression test that would have caught it.
7. Remove temporary instrumentation and record the diagnostic evidence in the
   task report, QA artifact, or review note that is already in scope.

## References

- `project-docs/templates/completion-report-template.md`
- `project-docs/templates/mid-phase-review-template.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
