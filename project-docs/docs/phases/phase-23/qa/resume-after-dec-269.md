# Phase Runner handoff after DEC-268/DEC-269 (Human checklist FAIL, tasks 19-21)

This handoff **replaces** steps 3-5 of "Runner closure" in
`qa/resume-after-dec-267.md`. Do **not** dispatch the final review yet. Keep
every earlier ledger event, TRG, and ESC. ESC-23-012 stays open. Do not
allocate a new TRG or ESC.

## Precondition: Human approval

Before any action, confirm that the DEC-269 `Status:` line in
`project-docs/docs/phases/phase-23/decisions.md` reads **ACCEPTED**. If it still
reads PROPOSED, append nothing, stop, and tell the Human that DEC-269 approval
is pending.

## What changed

- The Human recorded `qa/human-checklist.md`: 55 pass, 2 fail, 6 partial,
  overall FAIL. DEC-268 (Human) directs that findings A1-A3 be repaired in
  Phase 23.
- DEC-269 allocates three supplemental tasks, in this order:
  1. `tasks/task-19.md`: Human-mode failures from the Stop-error family
     (A2; VAL-085).
  2. `tasks/task-20.md`: `RequirementRef.requirement_id` and the requirement
     Human views (A3; VAL-089).
  3. `tasks/task-21.md`: rework exhaustion resolved by a
     `review-rework-extension` Human Outcome; the rework budget excluded from
     generic admission; the denial guidance line (A1; VAL-086..088).
- The validation contract now holds 89 rows (VAL-085..089 added).

## Ledger binding

The ledger currently validates as `VALID state=ADVANCEABLE`, after
`2026-10-01 20:59 TASK task-18 COMPLETE`. No STOP is active; the 16:20 STOP was
consumed by the 20:17 RUN-RESUME. Append one **normal** RUN-RESUME:

- Sequence contract: `v2`
- Recovery directive schema: `v2`
- Prior run/session or STOP: the 2026-10-01 20:17 RUN-RESUME session, ending
  at 20:59 TASK task-18 COMPLETE
- Governing TRG, Governing ESC, Trigger signature, Attempt count, Resolution
  authority, Next recovery: `N/A`. Omit `Repair owners`.
- Continuing task: `task-19`
- Checkpoint and Resolution evidence: what you actually verified, namely:
  - the HEAD commit and its subject;
  - DEC-268 and DEC-269 present, and DEC-269 ACCEPTED;
  - tasks 19-21 present;
  - VAL-085..089 present;
  - this handoff present;
  - a clean product tree.

Then run the normal task loop, one task at a time and in order:

1. Dispatch a Developer for task 19 with that task's paths, its exact commands,
   and its required reads.
2. Run report-lint yourself and verify the commit.
3. Append `TASK task-19 COMPLETE` with `QA-plan commands beyond the task
   spec: none required`.
4. Repeat steps 1-3 for task 20, then task 21.

Ordinary recovery, triage, and STOP rules apply.

## Reviews

1. **review-15 (window 19-21).** After task 21, dispatch a fresh QA Reviewer
   for `qa/review-15.md` covering task-19, task-20, and task-21 and
   VAL-085..089. Append `REVIEW review-15`. A REPAIR goes to the named owning
   task and needs a fresh independent re-review.
2. **Final review (recorded as review-16).** Only after review-15 PASSes,
   dispatch a **fresh, independent** QA Lead for `qa/final-review-1.md`.
   Append it as `REVIEW review-16` with `Path:
   project-docs/docs/phases/phase-23/qa/final-review-1.md` and `Tasks covered:
   task-1 through task-21`. It must:
   - cite and preserve `final-review.md`;
   - verify the MF-1 repair (`b84b17c`, `b37c169`) and rerun
     `go test -race -timeout 30m ./...` and the live golden tests;
   - rerun the seven-command phase gate and the section 22 timing controls;
   - rerun the linked-seed 16 to 17 migration with row and link preservation
     and `foreign_key_check`, and evaluate VAL-005 under DEC-267;
   - record an outcome for all 89 rows. VAL-070..073 stay Human-owned.
   - regenerate the Human fixtures against the reviewed HEAD, using the
     `qa/human-fixtures.md` fresh-generation command, and verify every
     invocation still runs;
   - author `qa/human-checklist-2.md` with only the seven re-check rows:
     HUMAN-03a, HUMAN-03b-resolution, HUMAN-03b-internals, HUMAN-04a-17,
     HUMAN-04b-17, HUMAN-04c-17, and HUMAN-04a-18. The rework section includes
     task 21's resolution step. `qa/human-checklist.md` stays immutable.
3. **If the final review PASSes technically but the Human rows are still
   blank,** stop at ADVANCEABLE and notify the Human. After the Human records
   `human-checklist-2.md`, resume with a normal RUN-RESUME.
   - If the result is PASS, append RUN-END naming `final-review-1.md` and
     `human-checklist-2.md`.
   - If a Human row FAILs, stop and route to the Tech Lead.

## Simulation check

The Tech Lead checked the event binding on a scratch copy of the ledger,
triage log, and decisions at `671f358`, using
`workflow-sequence-lint.sh --append-stdin`:

- normal RUN-RESUME (`Continuing task: task-19`, all recovery fields `N/A`);
- TASK task-19, task-20, and task-21 COMPLETE;
- REVIEW review-15 (tasks 19-21);
- REVIEW review-16 (`final-review-1.md`, tasks 1-21);
- RUN-END naming `final-review-1.md` and `human-checklist-2.md`.

Each event returned `APPEND_CONFIRMED state=ADVANCEABLE` with `exit_code: 0`,
and the standalone check returned `VALID state=ADVANCEABLE`. Two malformed
candidates were rejected with `exit_code: 1`: a non-`review-N` REVIEW ID, and
an invalid heading time. **The sequence lint does not enforce task order or
review cadence. This handoff and the QA cadence do.** The simulation proves
only the event grammar; use real checkpoints and receipts.

This handoff grants authority and scope. It is not a verification receipt.
The Tech Lead ran no implementation gate.
