# Phase Runner handoff after ESC-23-005

Read DEC-258, task 11's completion report, `qa/esc-23-005-repair.md`, and the
ESC-23-005 resolution. Reconcile the actual branch HEAD and worktree, then
verify the report's commit hash and subject. Run report-lint and
sequence-lint yourself:

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-11-done.md project-docs/docs/phases/phase-23/tasks/task-11.md
project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-23/run-ledger.md
```

The retained ledger ends in the TRG-23-005 STOP. Preserve every earlier
failure, TRG record, and the STOP, and allocate no replacement IDs. Append the
schema-v2 RUN-RESUME with the actual verified checkpoint and your own
verification evidence. Its STOP fields are:

- Governing TRG: TRG-23-005
- Governing ESC: ESC-23-005
- Trigger signature: review-7:qa-rework:task-11-BC-3
- Attempt count: 5 (increments TRG-23-005's attempt 4)
- Resolution authority: DEC-258
- Continuing task: task-11
- Next recovery: task-11 | gate-repair | task-scope | review-7:qa-rework:task-11-BC-3 | 5

Include `Sequence contract: v2` and `Recovery directive schema: v2`. The
prior-STOP reference names TRG-23-005 and ESC-23-005. Omit `Repair owners`,
because the open review-7 owner list from DEC-256 still stands.

The next event is a TASK-RECOVERY with the same task, trigger, and attempt,
category `gate-repair`, and mode `task-scope`. Both `Source` and `Repair
authority` reference **DEC-258 and TRG-23-005**. Reset no-progress to `0/2`
and set the control-conformance fields to `N/A`.

DEC-258 lets the Runner reconcile the committed delegated repair without
redispatching implementation:

1. Verify the report and evidence.
2. Record REPAIR with `Outcome: RESOLVED` and `Next checkpoint: VERIFY/REPORT`.
3. After report-lint succeeds, record TASK task-11 COMPLETE, naming commit
   `3dff4574a71438900100449e36d7044969ba9162`.

Cite the repair agent's committed gate evidence and record your own
mechanical checks. Do not claim fresh execution of the gates in the Runner
session.

After TASK task-11 COMPLETE, the sequence is `REREVIEW_PENDING`. Dispatch an
independent QA Lead for `qa/checkpoint-after-task-11-3.md`, the fresh re-review
of review-7 required by DEC-255. **Task 12 stays held until it returns PASS
with Must-fix 0.**

Tell the reviewer that DEC-258 records a headroom warning: `cmd/fid` is at
531.533 s, which is 88.6% of its budget. This passes VAL-079, but the Tech
Lead must schedule a speed item before task 13 adds `cmd/fid` tests.

## Simulation check

The handoff tuple was checked on a scratch copy of the whole `phase-23`
directory. The copy held the retained ledger plus RUN-RESUME, TASK-RECOVERY,
REPAIR RESOLVED, and TASK COMPLETE, each labelled STRUCTURE SIMULATION ONLY.
`workflow-sequence-lint.sh` returned `VALID state=REREVIEW_PENDING`
(terminated `exit_code: 0`). A first simulation that used attempt 4 was
rejected with `RUN-RESUME Attempt count must increment its governing TRIAGE`.
That is why attempt 5 is used above.

The canonical ledger separately returned `VALID state=STOPPED` (terminated
`exit_code: 0`). This proves the event binding is accepted. It does not show
that the Runner has resumed. Replace the simulation values with real
checkpoints and receipts.
