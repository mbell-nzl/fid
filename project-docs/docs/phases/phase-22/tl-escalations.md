# Phase 22 Escalations

## Retained original escalation

### ESC-22-001 - Protected workflow-sequence control cannot represent DEC-189 supersession

- **Phase:** 22
- **Status:** OPEN
- **Raised by:** TRG-22-006
- **Category:** protected workflow-control definition/policy
- **Summary:** `project-docs/tools/workflow-sequence-lint.sh` rejects the user-authorized TRIAGE TRG-22-005 and bound task-2 TASK-RECOVERY despite the accepted contract permitting a later same-task schema-v2 TRIAGE REPAIR to replace its pending dispatch and retain the review-form trigger identity.
- **Evidence:** From `/home/mark/projects/fid`, the supplied exact append reproducer terminated `exit_code: 1` with `INVALID TRIAGE REPAIR must be followed by its bound TASK-RECOVERY before TRIAGE; INVALID schema v2 TRIAGE Originating task does not match Trigger signature`; the ledger remained unchanged and standalone validation remained `VALID state=REPAIR_PENDING`.
- **Why blocked:** The required corrective edit is to `project-docs/tools/workflow-sequence-lint.sh`. `09_AI_WORKFLOW.md` §§2.3 and 9 prohibit normal control-conformance self-repair of that verifier; changing it requires separate protected workflow-control authority.
- **Required resolution:** Full Tech Lead or authorized protected workflow-control decision must reconcile the verifier implementation with `project-docs/templates/run-ledger-template.md` and DEC-189's same-task supersession and review-form trigger contract, then provide a concrete authorized recovery route. Preserve the valid Phase 22 ledger; do not bypass the validator or authorize product/task recovery first.
- **References:** TRG-22-006; `project-docs/specs/09_AI_WORKFLOW.md:427-458,586-614,703-713`; `project-docs/templates/run-ledger-template.md`; `project-docs/docs/phases/phase-22/decisions.md:222-270` (DEC-189); `project-docs/docs/phases/phase-22/run-ledger.md:158-209`.

## Resolution ESC-22-001

- **Status:** RESOLVED, 2026-09-06.
- **Authority:** DEC-190 and DEC-191; Human-authorized workflow maintenance.
- **Repair commit:** `3bb42b4efd391b216b752db2c31fe793f3538bac`.
- **Evidence:** All maintenance gates passed. The real ledger accepted normal
  RUN-RESUME, TRIAGE TRG-22-005 and the exact task-2 attempt-2 recovery, each
  with `APPEND_CONFIRMED state=REPAIR_PENDING` (terminated `exit_code: 0`).
  Standalone sequence: `VALID state=REPAIR_PENDING` (terminated `exit_code: 0`).
- **History:** The rejected BLOCK and STOP appends never persisted. Their
  original failure evidence is retained above and in TRG-22-006. No fabricated
  STOP/resume pair or control-conformance amendment was introduced.
- **Remaining work:** Task-2 product recovery under DEC-189, task-3 MF-3,
  and fresh independent review-1 re-review. No Human decision is pending.
- **Full receipts:** `../../reviews/2026-09-06-phase-22-workflow-recovery.md`.
