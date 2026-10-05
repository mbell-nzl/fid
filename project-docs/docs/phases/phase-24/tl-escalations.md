# Phase 24 Escalations

## ESC-24-002 - Accepted Task 6 RUN-RESUME binds an unsatisfiable recovery tuple

- **Raised:** 2026-10-06 by the Phase Runner (reported to the Human).
  Allocated by the Tech Lead in DEC-283. The ledger admits no TRIAGE event
  while resume-pending, so no TRG exists (as with DEC-149 clause 6).
- **Protected category:** protected verifier (`workflow-sequence-lint.sh`)
  defect; orchestration integrity.
- **Originating task:** task-6.
- **Trigger signature:** task-6:session-liveness:developer-role-capacity-after-bounded-retry
- **Condition:** the 2026-10-05 20:48 RUN-RESUME records trigger
  `task-6:session-liveness:...` at attempt 3. Its `Next recovery` names
  `task-6:session-continuation:...` at attempt 1. The verifier binds the
  next TASK-RECOVERY to both, so no candidate validates. Rejected
  candidates: a REPAIR (`RUN-RESUME must be followed by its bound
  TASK-RECOVERY`) and a TASK-RECOVERY (tuple mismatch). The ledger is
  unchanged and still `VALID state=REPAIR_PENDING`.
- **Status:** RESOLVED 2026-10-06 by DEC-283 (decision). The Runner then
  carries out the clause 4 verifier maintenance and the clause 6 ledger
  sequence; no further Tech Lead or Human input is needed unless the
  verification items cannot be produced as stated.
- **Links:** [DEC-283](decisions.md); [run ledger](run-ledger.md);
  [Task 6 report](tasks/task-6-done.md).

## ESC-24-001 - Task 6 Developer role unavailable after bounded model-capacity retry

- **Raised:** 2026-10-05 by TL-Lite in TRG-24-003.
- **Protected category:** unrepairable external infrastructure/capability
  failure after bounded retry.
- **Originating task:** task-6.
- **Trigger signature:** task-6:session-liveness:developer-role-capacity-after-bounded-retry
- **Attempt count:** 2 interrupted spawn attempts. Neither returned a typed
  result, so neither consumes the no-progress budget.
- **Condition:** the selected Developer model reported capacity exhaustion
  on two spawns.
- **Status:** RESOLVED 2026-10-05 20:48. The Runner verified the restored
  capability:
  `restored-capability:developer-role-spawn-01a10b08-efe6-7050-8f88-65c524e62519`.
  The resumed session committed `1baaecc` and `7298964`. DEC-283 records
  the closure and reconciles the resume's recovery tuple.
- **Links:** [TRG-24-003](triage-log.md); [DEC-283](decisions.md).
