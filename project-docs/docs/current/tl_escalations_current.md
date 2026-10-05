# TL Escalations — Current

## Open escalations

### ESC-24-001 - Task 6 Developer role unavailable after bounded model-capacity retry

- **Raised:** 2026-10-05 by TL-Lite in TRG-24-003.
- **Protected category:** unrepairable external infrastructure/capability failure after bounded retry.
- **Originating task:** task-6.
- **Trigger signature:** task-6:session-liveness:developer-role-capacity-after-bounded-retry.
- **Attempt count:** 2 interrupted spawn attempts; neither returned a typed result, so neither consumes the task's no-progress budget.
- **Condition:** The required Developer role failed to spawn twice because the selected model reported capacity exhaustion. The first interruption routed ACTION REPAIR; the bounded retry routed ACTION TRIAGE. No implementation, report, or task outcome came from either spawn.
- **Status:** OPEN; Phase Runner must stop Task 6 dispatch and phase advance.
- **Resolution condition:** Restore the required Developer role capability and have the Phase Runner verify that the role can spawn. Then resume the same Task 6 recovery from its existing checkpoint under TRG-24-003/ESC-24-001, preserving the trigger and no-progress count. Do not substitute an unapproved role or model override.
- **Human/external action:** Arrange restoration of the selected Developer model's capacity or the required-role spawn capability, then provide the Phase Runner with the capability evidence for RUN-RESUME.
- **Links:** [TRG-24-003](../phases/phase-24/triage-log.md); [Task 6 checkpoint](../phases/phase-24/tasks/task-6-done.md); `project-docs/specs/09_AI_WORKFLOW.md` §2.3 and §9.

## Recently resolved

All resolved and archived in
[Phase 23 escalations](../phases/phase-23/tl-escalations.md):
- ESC-23-012: phase-end final review blockers and Human judgment; resolved
  2026-10-03 by DEC-267..271 (tasks 19-21, final review 2 PASS, Human
  checklist 2 7/7 PASS).
- ESC-23-011 (model capacity), ESC-23-010 (DEC-265/266), ESC-23-009
  (DEC-264), ESC-23-008 (DEC-263), ESC-23-007 (DEC-262), ESC-23-006
  (DEC-260), ESC-23-005 (DEC-258), ESC-23-004 (DEC-249..256), ESC-23-003
  (DEC-239), ESC-23-002 (DEC-238), ESC-23-001 (DEC-235).

## Standing handling

Register an escalation when allocated, not after resolution. Ordinary causal
repair and incomplete in-scope work do not themselves require a Human stop.
Read accepted ownership before calling it ambiguous (DEC-194). A real contract/
schema contradiction needs the missing authority, not a weakened assertion
(DEC-195). External interruptions do not consume repair no-progress budget.

Prior archives: [Phase 22](../phases/phase-22/tl-escalations.md),
[Phase 21](../phases/phase-21/tl-escalations.md),
[Phase 20](../phases/phase-20/tl-escalations.md).
