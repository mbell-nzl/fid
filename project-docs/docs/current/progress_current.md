# Progress — Current

Updated 2026-10-11. **Phase 24 "Operating experience and input
convergence" is at phase end** (Heavy mode, DEC-277; 15 tasks with 9a).
Tasks 1-14 are complete and QA windows A-E have passed (review-12 closed
window E under DEC-291).

The final review (`qa/final-review.md`, `6d633351`) returned BLOCK. All 63
deterministic VAL rows PASS and every gate passes. VAL-053..057 wait for
the Human checklist (B1-H), and QA-36 fails on two missing owned inventory
receipts: Task 1 R4 (F1) and Task 10 R10 (F2). TRG-24-019 raised ESC-24-010
and the ledger is `STOPPED`. DEC-292 disposes it: the Human records the
checklist, then an evidence-only Task 14 recovery, report-only Task 1 and
Task 10 recoveries, and a fresh `final-review-1.md`. ESC-24-010 stays open
until the Human's checklist is committed and DEC-292 clause 9 is appended.
Earlier: DEC-291 (ESC-24-009), DEC-290 (ESC-24-008), DEC-289 (Task 13
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

1. Human: record and commit `../phases/phase-24/qa/human-checklist.md`
   (DEC-292 clause 1). Live workspaces expire about 11:50 NZDT 2026-10-12.
2. Tech Lead: DEC-292 clause 9 closes ESC-24-010 (or a new decision
   dispositions any FAIL/PARTIAL row).
3. Phase Runner: DEC-292 clause 6 route, then review-13 /
   `final-review-1.md` and RUN-END.
4. Tech Lead wrap-up; then Phase 25 "Parallel and evolving work" and
   Phase 26 "Qualification and baseline convergence".

A deep workflow retrospective is recommended under CF-036/039 (see Phase 23
workflow notes). Add the DEC-283..292 workflow notes to it. It is separate from
Phase 24.

## Earlier phases

Phase 22 closed under DEC-220/221: [wrap audit](../phases/phase-22/wrap-up-audit.md).
Claims qualification is Phase 25; external adapter qualification Phase 26.
Live carry-forwards are in [active decisions](decisions_active.md).
