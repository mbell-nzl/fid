# Progress — Current

Updated 2026-10-06. **Phase 24 "Operating experience and input
convergence" is executing** (Heavy mode, DEC-277; 15 tasks with 9a). Tasks
1-5 are complete, and QA window A passed (review-2). Task 6 is
`RECOVERY_REQUIRED` at `1baaecc`/`7298964`. Its focused slice passes, but
both full gates fail at 21 retained fixtures it displaced.

DEC-283 closes ESC-24-001 (capacity restored) and ESC-24-002. The accepted
20:48 RUN-RESUME bound an unsatisfiable recovery tuple and wedged the
ledger. DEC-283 reconciles the tuple without rewriting history and
authorizes a bounded verifier repair outside the run sequence. It also
amends Task 6's scope and rules on stored Claim authority. Planning history:
DEC-278..282. Phase 23 closed under DEC-271 (PR #25, `a46ec77`).

- Contract: [public interface](../phases/phase-24/public-interface.md);
  [inventory](../phases/phase-24/input-inventory.md);
  [validation contract](../phases/phase-24/qa/validation-contract.md);
  [decisions](../phases/phase-24/decisions.md).
- Phase 23 evidence: [final acceptance](../phases/phase-23/qa/final-acceptance.md),
  [wrap audit](../phases/phase-23/wrap-up-audit.md).

## Next

1. Phase Runner: DEC-283 clause 4 verifier maintenance in a clean temporary
   worktree (`--ff-only` integration). Then append clause 6 (a)-(c) and
   dispatch the Task 6 attempt-4 continuation.
2. QA windows: B after 6, C after 8, D1 after 9a, D2 after 11, E after 14.
   The Human walkthrough follows task 11.
3. Phase 25 "Parallel and evolving work" and Phase 26 "Qualification and
   baseline convergence" follow.

A deep workflow retrospective is recommended under CF-036/039 (see Phase 23
workflow notes). It is separate from Phase 24.

## Earlier phases

Phase 22 closed under DEC-220/221: [wrap audit](../phases/phase-22/wrap-up-audit.md).
Claims qualification is Phase 25; external adapter qualification Phase 26.
Live carry-forwards are in [active decisions](decisions_active.md).
