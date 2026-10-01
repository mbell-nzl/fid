# TL Escalations — Current

## Open escalations

### ESC-23-012 - Phase-end accepted Evidence column count conflicts and Human judgments remain pending

- **Status:** PARTIALLY RESOLVED 2026-10-01 under DEC-267 (Tech Lead,
  Human-delegated). It stays open for BC-B only. Governing triage:
  TRG-23-016.
  - **BC-A / MF-2, resolved:** the accepted DDL governs. 0017 adds two
    Evidence columns, `review_id` and `human_outcome_id`. "Three" was a
    planning miscount from `f94e354`. The VAL-005 criterion, the task-1
    criterion, and DEC-223 now name the two columns. No implementation
    change is needed.
  - **MF-1, routed:** task 18 owns it by planning amendment, under the same
    repair boundary QA set. The Runner resumes with
    [the handoff](../phases/phase-23/qa/resume-after-dec-267.md)
    (task-18 ledger attempt 2).
  - **BC-B, open:** the Human records VAL-070..073 in the
    [Human checklist](../phases/phase-23/qa/human-checklist.md).
- **Resolution condition (remaining):** the Human records all VAL-070..073
  outcomes. A fresh independent QA Lead then writes `final-review-1.md`,
  which covers the MF-1 repair, VAL-005 under DEC-267, and the Human
  results. Phase advance needs that review.
- **References:** [TRG-23-016](../phases/phase-23/triage-log.md);
  [final review](../phases/phase-23/qa/final-review.md);
  [DEC-267](../phases/phase-23/decisions.md); VAL-005, VAL-070..073;
  [full record](../phases/phase-23/tl-escalations.md).

## Recently resolved

All resolved and archived in
[Phase 23 escalations](../phases/phase-23/tl-escalations.md):
- ESC-23-011: Developer model capacity; resolved 2026-10-01 by restored
  capability (07:19 RUN-RESUME), and task 18 then completed.
- ESC-23-010 under DEC-265 (19 Failure codes: 18 added, 1 retained) and
  DEC-266 (`authority-*` details converged to the Phase 21 four fields).
- ESC-23-009 under DEC-264 (0018 down restores the schema-17 state).
- ESC-23-008 under DEC-263 (Stop-blocked denial fields in failure details).
- ESC-23-007 under DEC-262 (ActivityEvent rows split at 0018).
- ESC-23-006 under DEC-260; ESC-23-005 under DEC-258; ESC-23-004 under
  DEC-249..256; ESC-23-003 under DEC-239; ESC-23-002 under DEC-238;
  ESC-23-001 under DEC-235. R-28 was dispositioned under DEC-237 without an ID.
- ESC-22-005 under DEC-221; see
  [Phase 22 escalations](../phases/phase-22/tl-escalations.md).

## Standing handling

Register an escalation when allocated, not after resolution. Ordinary causal
repair and incomplete in-scope work do not themselves require a Human stop.
Read accepted ownership before calling it ambiguous (DEC-194). A real contract/
schema contradiction needs the missing authority, not a weakened assertion
(DEC-195). External interruptions do not consume repair no-progress budget.

Prior archives: [Phase 21](../phases/phase-21/tl-escalations.md),
[Phase 20](../phases/phase-20/tl-escalations.md),
[Phase 19](../phases/phase-19/tl-escalations.md).
