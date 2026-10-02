# Phase Runner handoff after DEC-267 (ESC-23-012, phase-end final review)

DEC-267 resolves the Tech Lead part of the 2026-10-01 16:20 STOP. Keep all
of the following:

- TRG-23-016 (BLOCK) and ESC-23-012;
- the 16:20 STOP;
- `qa/final-review.md`, `qa/human-checklist.md`, and `qa/final-evidence/` at
  `74182db`, including the failed full-checkptr run;
- every earlier ledger event.

Do not allocate a replacement TRG or ESC.

## What DEC-267 decided

1. **MF-2 / BC-A is resolved, with no implementation.** The accepted DDL
   governs: 0017 adds two Evidence columns, `review_id` and
   `human_outcome_id`. "Three" was a planning miscount. The VAL-005 pass
   criterion, the task-1 criterion, and DEC-223 now name the two columns.
2. **MF-1 goes to task 18** under "Planning Amendment - DEC-267" in
   `tasks/task-18.md`. The repair boundary is QA's own; only the owner
   differs. DEC-267 Part 2 gives the reasons, one of which is that the
   validator binds this STOP's resume to task 18.
3. **BC-B stays with the Human.** VAL-070..073 are recorded in
   `qa/human-checklist.md`. ESC-23-012 stays open for BC-B only.
4. **Unchanged:** contracts, DDL, product code, every command and flag,
   budgets, validators, and the other 83 assertions.

## Resume binding

The governing TRIAGE is TRG-23-016, with `Attempt count` 1. The resume is
therefore attempt **2** on the same trigger.

- Governing TRG: TRG-23-016
- Governing ESC: ESC-23-012
- Trigger signature: task-18:qa-block:phase-end-final-review
- Attempt count: 2
- Resolution authority: DEC-267
- Continuing task: task-18
- Next recovery: task-18 | qa-rework | task-scope | task-18:qa-block:phase-end-final-review | 2

Include `Sequence contract: v2` and `Recovery directive schema: v2`. In the
prior-STOP reference, name the 16:20 STOP, TRG-23-016, and ESC-23-012.
Record the checkpoint you actually verified and your own resolution
evidence:

- the commit that records DEC-267, and its subject;
- DEC-267 present in phase decisions;
- the task-18 Planning Amendment - DEC-267 present;
- the corrected VAL-005 row;
- ESC-23-012 marked partially resolved, with BC-B pending;
- this handoff present.

The next event is a TASK-RECOVERY for task-18 that repeats the resume tuple
exactly: the same trigger, attempt 2, and a `Checkpoint` byte-identical to
the RUN-RESUME `Checkpoint`. Its fields are:

- Category: `qa-rework`
- Repair mode: `task-scope`
- Source: `DEC-267 and TRG-23-016; project-docs/docs/phases/phase-23/qa/resume-after-dec-267.md`
- Repair authority: `DEC-267 and TRG-23-016`
- Consecutive no-progress attempts: `0/2`
- Every control-conformance field: `N/A`. This is a test-harness repair
  inside the amended task scope, not a control-conformance repair.
- Allowed files: `ordinary task scope`. The DEC-267 amendment names them.

## Bounded task-18 continuation (same-task Developer)

Start from HEAD. Read `task-18-done.md`, the DEC-267 amendment, and the
`final-review.md` MF-1 section first.

1. **Diagnose.** Run the amendment's direct child timing probe in fresh
   scratch and record the result. If the cause is not harness startup,
   return `RECOVERY_REQUIRED`.
2. **Red first.** Write the deterministic regression and record its
   failure.
3. **Repair.** Keep the change inside `workspace_fixture_test.go`, and in
   `live_output_test.go` only as far as needed.
4. **Commands.** Run the focused command, then the five task-18 commands,
   then `go test -race -timeout 30m ./...`, all exactly as the amendment
   lists them. Record each with observed termination.
5. **Report.** Update `task-18-done.md`:
   - keep all earlier evidence;
   - add a DEC-267 deviation note that quotes QA's failed token
     (`terminated exit_code: 1`), and gives the root cause, the repair, and
     the later result;
   - record the provenance (`f0eec23`, `bca799f`).

   Commit the implementation and the report with
   `Refs: phase-23/task-18`.

## Runner closure

1. Run report-lint and sequence-lint yourself, and verify the commit hash
   and subject.
2. Record REPAIR task-18 with `Outcome: RESOLVED`, then TASK task-18
   COMPLETE. In TASK COMPLETE, `QA-plan commands beyond the task spec` is
   `none required`, because the amendment already lists the section 22
   control.
3. **Wait for the Human.** Do not dispatch the final review until
   `qa/human-checklist.md` records an outcome for every VAL-070..073 row.
   If the Human has not finished, stop at `ADVANCEABLE`, notify the Human,
   and resume afterwards with a normal non-STOP RUN-RESUME.
4. Dispatch a **fresh, independent** QA Lead. Its final review goes in
   `qa/final-review-1.md`, which cites and preserves `final-review.md`. It
   must:
   - verify the MF-1 repair and the task-18 report, and rerun
     `go test -race -timeout 30m ./...` and the live golden tests;
   - rerun the seven-command phase gate and the section 22 timing controls,
     because the harness change invalidates their evidence;
   - rerun the linked-seed 16 to 17 migration, row and link preservation,
     and `foreign_key_check`, and evaluate VAL-005 under the DEC-267
     criterion;
   - read the Human's VAL-070..073 outcomes. If the Human's fixtures came
     from an earlier HEAD, confirm that
     `git diff --name-only <fixture HEAD> <reviewed HEAD> -- . ':(exclude)*_test.go' ':(exclude)project-docs'`
     prints nothing. Otherwise the Human must regenerate the fixtures and
     re-judge.
5. On PASS, append RUN-END that names `final-review-1.md` and the Human
   checklist. On REPAIR or BLOCK, triage as usual.

## Simulation check

The binding was checked on a scratch copy of the `phase-23` ledger, triage
log, and decisions, taken at `74182db` plus the Runner's uncommitted
TRG-23-016 / STOP entries. A stub DEC-267 was used, and every event went
through `workflow-sequence-lint.sh --append-stdin`:

- The unmodified ledger returned `VALID state=STOPPED`.
- The RUN-RESUME above returned `APPEND_CONFIRMED state=REPAIR_PENDING`.
- The TASK-RECOVERY above returned `APPEND_CONFIRMED state=REPAIR_PENDING`.
- REPAIR RESOLVED returned `COMPLETION_PENDING`. TASK COMPLETE and then
  RUN-END each returned `APPEND_CONFIRMED state=ADVANCEABLE`. The standalone
  check returned `VALID state=ADVANCEABLE`.
- A task-1 continuation was rejected:
  `INVALID RUN-RESUME Continuing task does not match its Trigger signature`.
  So was recording the final review as REVIEW before the bound recovery:
  `INVALID RUN-RESUME must be followed by its bound TASK-RECOVERY before REVIEW`.

Each passing run terminated `exit_code: 0`, and each rejection
`exit_code: 1`. This proves only the event binding; use real checkpoints
and receipts.

This handoff grants authority and scope. It is not a verification receipt.
The Tech Lead ran no implementation gate.
