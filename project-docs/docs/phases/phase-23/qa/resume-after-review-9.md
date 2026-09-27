# Phase Runner handoff after review-9

Review-9 (`qa/review-9.md`, commit `3dbf5c6`) returned `VERDICT: REPAIR`.
It found Must-fix 2, Improvement 0 and no blocker candidates. This is an
ordinary QA REPAIR, so there is **no TRG, no ESC and no STOP**. The previous
Runner session ended after the review commit and before it appended the
REVIEW event. The ledger validates as `VALID state=ADVANCEABLE`, and its last
event is `TASK task-12 COMPLETE` (21:25). Do not spawn another reviewer for
window 10-12. Record the existing review.

Tech Lead decision: DEC-261 in phase decisions.

## Ledger sequence

1. **Normal `RUN-RESUME`.**
   - Set `Governing TRG`, `Governing ESC`, `Trigger signature`,
     `Attempt count`, `Resolution authority` and `Next recovery` to `N/A`.
   - Omit `Repair owners`.
   - Set `Continuing task` to `task-10`.
   - Record the checkpoint you actually verify.
2. **`REVIEW review-9`.**
   - `Path: project-docs/docs/phases/phase-23/qa/review-9.md`
   - `Tasks covered: task-10, task-11, task-12`
   - `Outcome: REPAIR`
   - `Owning task: task-10, task-12` (this is the Repair Plan row order)
   - `Re-review of: none`
   - `Must-fix: 2 / Improvement: 0 / blocker candidates: 0; QA minted no ESC`
3. **`TASK-RECOVERY task-10`.**
   - Category `qa-rework`, mode `task-scope`, authority `QA review`.
   - Trigger `review-9:qa-rework:task-10-MF-1`, attempt 1, no-progress `0/2`.
   - `Allowed files: ordinary task scope`; control fields `N/A`.
   - Then `REPAIR task-10` and `TASK task-10 COMPLETE`.
4. **`TASK-RECOVERY task-12`.**
   - Trigger `review-9:qa-rework:task-12-MF-2`, with the same field pattern
     as step 3.
   - Then `REPAIR task-12` and `TASK task-12 COMPLETE`.
5. **`REVIEW review-10`.** A fresh, independent QA Reviewer writes
   `qa/review-10.md` with `Re-review of: review-9`. Task 13 dispatches only
   after that review records PASS.

Tech Lead dry run: steps 2 and 3 were appended to a scratch copy of the
phase directory through `workflow-sequence-lint.sh --append-stdin`. Both
returned `APPEND_CONFIRMED state=REPAIR_PENDING`.

## Task 10 - MF-1 (VAL-046)

Dispatch a same-task Developer to follow `tasks/task-10.md`
`## Repair Amendment - review-9 MF-1 under DEC-261`. In summary:

- Prove the `artifact.availability-observed` sequence on an Envelope whose
  sole remaining unsatisfied Obligation is the Human one.
- Show the preceding non-qualifying state, and exactly one Stop and one
  Attention after the qualifying flip.
- Show a repeat qualifying Event while the Stop is open, with no duplicate.
- Show resolution by the matching Outcome.
- Correct the VAL-046 citations in `task-10-done.md`.

**Gate change (DEC-261):** task 10's exact block now uses the DEC-259 race
commands (`-gcflags=modernc.org/...=-d=checkptr=0`). Pass the Developer the
five commands exactly as they now appear in `task-10.md`. Report-lint on the
current report against the amended spec fails only on those two new race
commands, which the repair must record. Earlier rows in the old form stay as
history. Budgets and VAL-079 are unchanged.

## Task 12 - MF-2 (VAL-053)

Dispatch a same-task Developer to follow `tasks/task-12.md`
`## Repair Amendment - review-9 MF-2 under DEC-261`. In summary:

- Remove `requirement.synced` from `p22PortableEventTypes`. It was
  introduced in `5ec69e6`, before Phase 23.
- Add a focused test in which the fake reader rejects that type and still
  accepts a retained family.
- Record a source inventory that includes test files, and correct the
  report's absence claim.
- Do **not** add Phase 23 families to the fake; that is task 16.
- Run all seven exact task-12 commands, unchanged.

Tech Lead inventory at `3dbf5c6`: `rg -n "requirement\.synced" cmd internal`
finds exactly one hit, `cmd/fid/integration/p22_portable_fake_test.go:267`.

## Notes for the review-10 re-reviewer

- Judge task 10's five gates against `task-10.md` as amended by DEC-261,
  not against the pre-DEC-259 form cited in review-9's gate table.
- For MF-2, "catalog parity" means parity for the removed family
  (DEC-261 Part 2).
