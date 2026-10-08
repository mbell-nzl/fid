# Progress — Current

Updated 2026-10-09. **Phase 24 "Operating experience and input
convergence" is executing** (Heavy mode, DEC-277; 15 tasks with 9a). Tasks
1-9a are complete. QA windows A, B, C (late, DEC-284) and D1 have passed.
Task 10 is in recovery: implementation `6607143`, report `e0495b5`. The
ledger is `STOPPED` at TRG-24-011/ESC-24-005.

DEC-286 closes ESC-24-005. `export import --project` printed raw text and
exited 7 in both modes. It is a Phase 21 section 2 forbidden flag, so it
now returns `invalid-argument`/2 with the flag spelling and task 9a's
allowed list. Task 10 gains `root.go` and a pure `canonical_request.go`
extraction; `sql` and `journal inspect` keep their raw refusal. Recovery
continues at attempt 2 under TRG-24-011's trigger. Earlier: DEC-285
(`activity.list`, ESC-24-004), DEC-284 (Task 9, ESC-24-003) and DEC-283
(Task 6, ESC-24-001/002). Planning history: DEC-278..282. Phase 23 closed
under DEC-271 (PR #25, `a46ec77`).

- Contract: [public interface](../phases/phase-24/public-interface.md);
  [inventory](../phases/phase-24/input-inventory.md);
  [validation contract](../phases/phase-24/qa/validation-contract.md);
  [decisions](../phases/phase-24/decisions.md).
- Phase 23 evidence: [final acceptance](../phases/phase-23/qa/final-acceptance.md),
  [wrap audit](../phases/phase-23/wrap-up-audit.md).

## Next

1. Phase Runner: DEC-286 clause 8. Append the STOP-resume RUN-RESUME and
   TASK-RECOVERY task-10 (attempt 2), then dispatch the Task 10
   continuation.
2. Task 11, then QA window D2 (tasks 10-11) and the Human walkthrough. Then
   tasks 12-14 and window E.
3. Phase 25 "Parallel and evolving work" and Phase 26 "Qualification and
   baseline convergence" follow.

A deep workflow retrospective is recommended under CF-036/039 (see Phase 23
workflow notes). Add the DEC-283/284/285/286 workflow notes to it. It is separate from
Phase 24.

## Earlier phases

Phase 22 closed under DEC-220/221: [wrap audit](../phases/phase-22/wrap-up-audit.md).
Claims qualification is Phase 25; external adapter qualification Phase 26.
Live carry-forwards are in [active decisions](decisions_active.md).
