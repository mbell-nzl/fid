# Phase Runner handoff after ESC-23-002

Read DEC-238, Task 1's amended scope and completion report,
`qa/esc-23-002-repair.md`, and the escalation resolution. Reconcile actual
branch HEAD/worktree and verify the implementation hash and subject in the
report. Run report-lint and sequence-lint yourself:

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-1-done.md project-docs/docs/phases/phase-23/tasks/task-1.md
project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-23/run-ledger.md
```

The retained ledger ends in the historical STOP, not a newly open escalation.
Preserve earlier failures, TRG records, and STOP; allocate no replacement IDs.
Append the schema-v2 RUN-RESUME with the actual verified checkpoint and your
own verification evidence. Its STOP fields are:

- Governing TRG: TRG-23-002
- Governing ESC: ESC-23-002
- Trigger signature: task-1:gate-repair:final-required-gates
- Attempt count: 3
- Resolution authority: DEC-238
- Continuing task: task-1
- Next recovery: task-1 | gate-repair | task-scope | task-1:gate-repair:final-required-gates | 3

Include `Sequence contract: v2` and `Recovery directive schema: v2`.
The prior-STOP reference names TRG-23-002 and ESC-23-002.

The immediately following TASK-RECOVERY uses the same task, trigger, attempt,
category `gate-repair`, and mode `task-scope`. Both `Source` and `Repair
authority` reference **DEC-238 and TRG-23-002**. Restart no-progress at `0/2`;
control-conformance fields are `N/A`. Use the full task-scope field set from
the ledger template.

DEC-238 permits reconciliation of the already committed delegated repair:
verify its report/evidence, record REPAIR `RESOLVED`, then TASK task-1 COMPLETE
after report-lint succeeds. Cite the repair agent's committed test evidence;
record your own mechanical checks, without claiming fresh execution of those
tests in the Runner session. No repeat implementation dispatch is needed solely
because the Human delegated work across fixed role boundaries.

Validate each appended event using the normal sequence control. Continue with
Task 2 and the ordered task/review sequence through Task 8. **Tasks 9-11 remain
held pending explicit Human dispatch clearance.** Later review windows and
all task/phase gates still apply.

The handoff tuple was mechanically checked on a scratch copy of the retained
ledger plus RUN-RESUME, TASK-RECOVERY, REPAIR RESOLVED and TASK COMPLETE:
`workflow-sequence-lint.sh` returned `VALID state=ADVANCEABLE` (terminated
`exit_code: 0`). All hypothetical receipts were explicitly labelled
STRUCTURE SIMULATION ONLY. This proves the event binding is accepted, not that
the Runner has already resumed or performed its own report verification.
The canonical ledger separately returned `VALID state=STOPPED` (terminated
`exit_code: 0`). Replace simulation values with actual verified checkpoints
and receipts; do not copy simulated receipt text as evidence.
