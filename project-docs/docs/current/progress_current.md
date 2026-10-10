# Progress — Current

Updated 2026-10-11. **Phase 24 "Operating experience and input
convergence" is executing** (Heavy mode, DEC-277; 15 tasks with 9a). Tasks
1-13 are complete (Task 12 at `fc4bdff7`, Task 13 at `fd6783f9` under
DEC-289). QA windows A, B, C (late, DEC-284), D1 and D2 have passed.

Task 14 attempt 1 returned BLOCKED (`4b6a3837` / `98d45617`). The
take-control journey had no permitted source for the terminal-success
`--to` state (`qa_pass`) or for the Issue revision. TRG-24-016 raised
ESC-24-008, and the ledger is `STOPPED`. DEC-290 closes it: `resume` shows
`Revision:` and `Next states:` (JSON `issue.next_states`), and the
`issue transition` help points there. Task 14 owns the change (VAL-068),
and DEC-278's provenance classes are unchanged. Earlier: DEC-289 (Task 13
scope), DEC-288 (ESC-24-007), DEC-287 (ESC-24-006), DEC-286 (ESC-24-005),
DEC-285 (ESC-24-004), DEC-284 (ESC-24-003) and DEC-283 (ESC-24-001/002).
Planning history: DEC-278..282. Phase 23 closed under DEC-271 (PR #25,
`a46ec77`).

- Contract: [public interface](../phases/phase-24/public-interface.md);
  [inventory](../phases/phase-24/input-inventory.md);
  [validation contract](../phases/phase-24/qa/validation-contract.md);
  [decisions](../phases/phase-24/decisions.md).
- Phase 23 evidence: [final acceptance](../phases/phase-23/qa/final-acceptance.md),
  [wrap audit](../phases/phase-23/wrap-up-audit.md).

## Next

1. Phase Runner: DEC-290 clause 8. Append RUN-RESUME and TASK-RECOVERY
   task-14 (attempt 2), dispatch one same-task continuation, then REPAIR
   `RESOLVED` and TASK task-14 COMPLETE.
2. QA window E (tasks 12-14 and VAL-068), then the final review and Human
   checklist.
3. Phase 25 "Parallel and evolving work" and Phase 26 "Qualification and
   baseline convergence" follow.

A deep workflow retrospective is recommended under CF-036/039 (see Phase 23
workflow notes). Add the DEC-283..290 workflow notes to it. It is separate from
Phase 24.

## Earlier phases

Phase 22 closed under DEC-220/221: [wrap audit](../phases/phase-22/wrap-up-audit.md).
Claims qualification is Phase 25; external adapter qualification Phase 26.
Live carry-forwards are in [active decisions](decisions_active.md).
