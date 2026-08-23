# Phase 21 - TL Escalations

## Open Escalations

None. DEC-160 through DEC-163 resolve ESC-21-001..005. The ledger remains in its
stopped state until the DEC-163 verifier change lands and the Phase Runner
appends the authorized resume.

## Resolved Escalations

### ESC-21-005 - Multi-owner QA review recovery cannot bind the single-owner sequence control

- **Raised by:** TRG-21-005, corrected for external identity by TRG-21-006,
  after QA Review 1 of tasks 1-3.
- **Resolution:** DEC-163. A QA `REPAIR` review carries an ordered
  comma-separated owner list in exactly one `Owning task` field, matching the
  review artifact's Repair Plan row order. Every listed owner completes
  TASK-RECOVERY, REPAIR, and TASK complete before a linked re-review may close
  or replace the review. `CORRECTION` remains evidence-only and is rejected as
  an owner-selection mechanism. `workflow-sequence-lint.sh` is corrected under
  protected decision authority to implement this, with new drills, because it
  accepted an append that left the ledger unrecoverable. No historical record is
  rewritten; MF-1 and MF-2 keep their exact boundaries and evidence.
- **Recovery:** after the verifier change lands, RUN-RESUME attempt 2 under
  TRG-21-006 / ESC-21-005 / DEC-163 into `task-1 | qa-rework | task-scope |
  task-1:orchestration-integrity:review-1:multi-owner-ledger-recovery-binding |
  2` for review-1 MF-1, then an ordinary task-3 `qa-rework` recovery for MF-2,
  then fresh re-review `review-1b`. Task 4 stays gated.

### ESC-21-004 - Retained rebuild and role-row tests versus accepted Event-v2 contracts

- **Raised by:** TRG-21-004 after the DEC-160 policy-bridge recovery.
- **Resolution:** DEC-162 keeps the no-leaf-flags `project.rebuild` interface and
  the exact role-free Event-v2 bootstrap. Retained V1 tests are not protected
  contracts: the task whose accepted change invalidates a retained test adapts it
  in that same task, preserving subject and assertion strength. Task 2 gains the
  test-file scope, ordered adaptation rules, and a closed no-weakening list;
  `journal.RebuildWorkspace` drops its vestigial role parameter.
- **Recovery:** same-task, task-scope `in-scope-implementation` continuation of
  task 2 at attempt 6 under the TRG-21-004 trigger signature, with both
  previously failed required commands rerun byte-identically.

### ESC-21-001 - Role-free bootstrap versus retained V1 policy

- **Raised by:** TRG-21-001 after task 2.
- **Resolution:** DEC-160 keeps the exact zero-Effect Event-v2 bootstrap and
  authorizes a closed, nonpersistent seven-name V1 policy vocabulary only for
  inventoried legacy callers. Explicit legacy role Effects remain
  authoritative; new generic paths cannot consult it; task 20 removes it.
- **Recovery:** same-task task-2 continuation with focused policy red/green
  tests and byte-identical reruns of both failed required commands.

### ESC-21-002 - Markdown-delimited external trigger identity

- **Raised by:** TRG-21-002 after rejected TRIAGE candidates.
- **Resolution:** DEC-161 makes schema-v2 identity values plain single-line
  scalars. TRG-21-001/TRG-21-002 remain immutable malformed evidence; neither
  governs recovery.
- **Recovery:** none against the verifier or historical records.

### ESC-21-003 - Protected recovery route

- **Raised by:** TRG-21-003 and the valid STOP ledger.
- **Resolution:** DEC-161 uses already-valid TRG-21-003 as the governing record
  and supplies exact RUN-RESUME attempt 3 authority. No protected verifier edit
  or control-conformance repair is required.
- **Recovery:** Phase Runner appends RUN-RESUME, then task-scope TASK-RECOVERY
  for task 2 under DEC-160/DEC-161. Task 3 stays dependency-gated.

## Archive Rule

Any Phase 21 `BLOCK` must allocate an `ESC-21-NNN` through `triage-log.md` and
record its protected resolution condition here before Human attention.
