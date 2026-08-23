# Progress - Current

> In-flight status only. Historical progress remains in phase archives.

## Project State

**Phase 21 is active and mechanically STOPPED after QA Review 1; DEC-163
authorizes the review-1 repair route and the ordered-owner verifier
correction.** Tasks 1-3 are mechanically complete. Phase 20
completed/wrapped on 2026-08-22. The Phase 20 retrospective is committed in
`../reviews/2026-08-22-phase-20-retrospective.md`; Gate D3.2 is closed by
DEC-153; the initial slice/spec contract is DEC-154; carry-forward disposition
is DEC-155; the terminal gate amendment is DEC-156; and the complete interface,
bootstrap, inventory, sequencing, and validation correction is DEC-157.
QA-plan refinements and the unchanged deterministic rebuild target are DEC-158.
QA readiness and Human execution approval are DEC-159. DEC-160 preserves the
role-free bootstrap while bounding temporary V1 policy compatibility; DEC-161
selects the valid ledger recovery record without a verifier change; DEC-162
converges retained V1 tests inside the task whose accepted change breaks them;
DEC-163 makes a QA `REPAIR` review carry an ordered owner list and authorizes
the matching `workflow-sequence-lint.sh` correction with new drills.

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
- Dispatch state: **TASKS 1-3 COMPLETE; REVIEW 1 REPAIR / STOPPED AT
  TRG-21-006**. The DEC-163 verifier correction lands first; then the Phase
  Runner appends DEC-163's exact RUN-RESUME and TASK-RECOVERY (governing
  TRG-21-006 / ESC-21-005, attempt 2) for review-1 MF-1 on task-1, then MF-2 on
  task-3, then fresh re-review `review-1b`. Task 4 remains gated.
- Planning verification: **PASS** - deterministic coverage/inventory audit,
  `project-docs/tools/workflow-lint.sh`, and `git diff --check` all terminated
  `exit_code: 0` after the DEC-159 approval-state update.

## Carry-Forward

DEC-155 schedules CF-020/027 in task 8, CF-032 in task 12, and CF-024 in task
21. CF-008/018/023/026a/026b/028/029/030 are explicitly re-deferred; closed or
cancelled rows remain closed/cancelled.

## Next Work

1. Human applies DEC-163 items 4-5 to `workflow-sequence-lint.sh` and the drills,
   confirms `run-drills.sh` PASS plus a standalone `VALID state=STOPPED` on the
   phase-21 ledger, and commits on `phase-21`.
2. Phase Runner appends DEC-163's RUN-RESUME attempt 2 and its bound task-1
   `qa-rework` TASK-RECOVERY, then dispatches only review-1 MF-1.
3. After task 1 completes, Phase Runner runs review-1 MF-2 as a task-3
   `qa-rework` recovery, then spawns a fresh independent `review-1b` of tasks
   1-3. Task 4 is released only on a re-review PASS.
