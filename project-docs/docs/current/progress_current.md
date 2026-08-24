# Progress - Current

> In-flight status only. Historical progress remains in phase archives.

## Project State

**Phase 21 is active at task 8 with the ledger valid at `REPAIR_PENDING`;
DEC-164 authorizes the ordinary-TRIAGE verifier correction and exact task-8
recovery route.** Tasks 1-7 and reviews 1/1b and 2/2b are mechanically
complete. Phase 20
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
the matching `workflow-sequence-lint.sh` correction with new drills. DEC-164
keeps ordinary schema-v2 repair authority byte-exact and locates its governing
TRG reference in TASK-RECOVERY `Source`; its bound same-task recovery replaces
the pending dispatch whose blocker candidate the TRIAGE classified.

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
- Dispatch state: **TASKS 1-7 COMPLETE; TASK 8 REPAIR_PENDING AFTER
  TRG-21-007**. The DEC-164 verifier correction lands first; then the Phase
  Runner appends the exact ordinary task-8 recovery governed by TRG-21-007 and
  dispatches only its same-task continuation. No RUN-RESUME or new TRIAGE is
  required because the ledger never accepted TRG-21-008 or entered STOPPED.
- Planning verification: **PASS** - deterministic coverage/inventory audit,
  `project-docs/tools/workflow-lint.sh`, and `git diff --check` all terminated
  `exit_code: 0` after the DEC-159 approval-state update.

## Carry-Forward

DEC-155 schedules CF-020/027 in task 8, CF-032 in task 12, and CF-024 in task
21. CF-008/018/023/026a/026b/028/029/030 are explicitly re-deferred; closed or
cancelled rows remain closed/cancelled.

## Next Work

1. Human-directed workflow maintenance lands DEC-164's protected verifier,
   drill, and canonical-doc correction with full deterministic evidence.
2. Phase Runner validates `REPAIR_PENDING`, appends the pending TRG-21-007
   task-8 recovery with exact external authority, and dispatches only task 8.
3. Developer preserves failed attempts and retained assertion strength, fixes
   the remaining in-scope report/gate defects, reruns all five exact commands,
   and returns task 8 for two-argument report lint and Window 3 review.
