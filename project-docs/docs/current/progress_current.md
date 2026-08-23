# Progress - Current

> In-flight status only. Historical progress remains in phase archives.

## Project State

**Phase 21 is active and mechanically STOPPED after task 2; DEC-160/DEC-161
authorize bounded task-2 recovery.** Phase 20
completed/wrapped on 2026-08-22. The Phase 20 retrospective is committed in
`../reviews/2026-08-22-phase-20-retrospective.md`; Gate D3.2 is closed by
DEC-153; the initial slice/spec contract is DEC-154; carry-forward disposition
is DEC-155; the terminal gate amendment is DEC-156; and the complete interface,
bootstrap, inventory, sequencing, and validation correction is DEC-157.
QA-plan refinements and the unchanged deterministic rebuild target are DEC-158.
QA readiness and Human execution approval are DEC-159. DEC-160 preserves the
role-free bootstrap while bounding temporary V1 policy compatibility; DEC-161
selects the valid ledger recovery record without a verifier change.

## Planned Boundary

One Heavy Phase 21 uses 22 bounded Developer tasks in two internal groups. An
independent checkpoint after task 12 must prove Group 1 contract closure and
fresh Group 2 executor/verifier feasibility before task 13. The validation
contract contains VAL-001..087. Seven three-task review windows, Phase Runner
two-argument post-commit report lint before each next dispatch, final QA, and
one Human terminal-output parent assertion with four Human subchecks apply.

## Approval State

- Human approval of revised Heavy mode, exact contracts, VAL contract, and task
  pack: **APPROVED 2026-08-23 / DEC-159**.
- QA Lead Heavy QA plan with fresh feasibility: **READY** at `e8a636e`; no
  outstanding Must-fix.
- Product implementation: **AUTHORIZED UNDER THE ORDERED TASK PACK**.
- Dispatch state: **TASK 1 COMPLETE; TASK 2 STOPPED / RECOVERY AUTHORIZED**.
  The Phase Runner must append DEC-161's exact RUN-RESUME and TASK-RECOVERY
  before dispatching the same-task continuation. Task 3 remains gated.
- Planning verification: **PASS** - deterministic coverage/inventory audit,
  `project-docs/tools/workflow-lint.sh`, and `git diff --check` all terminated
  `exit_code: 0` after the DEC-159 approval-state update.

## Carry-Forward

DEC-155 schedules CF-020/027 in task 8, CF-032 in task 12, and CF-024 in task
21. CF-008/018/023/026a/026b/028/029/030 are explicitly re-deferred; closed or
cancelled rows remain closed/cancelled.

## Next Work

1. Phase Runner validates the unchanged stopped ledger and appends DEC-161's
   RUN-RESUME attempt 3 with its exact task-2 recovery tuple.
2. Developer resumes only task 2 under DEC-160, adds focused compatibility
   tests first, and reruns the two failed required commands byte-for-byte.
3. Phase Runner records resolved recovery, verifies the committed completion
   report with two-argument report lint, and only then releases task 3.
