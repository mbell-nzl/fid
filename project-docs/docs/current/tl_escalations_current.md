# TL Escalations — Current

## Open escalations

### ESC-23-012 - Phase-end accepted Evidence column count conflicts and Human judgments remain pending (Human checklist FAIL; A1-A3 repairs)

- **Status:** PARTIALLY RESOLVED 2026-10-01 under DEC-267 and DEC-268.
  It stays open for the A1-A3 repairs. Governing triage: TRG-23-016.
  - **BC-A / MF-2, resolved:** the accepted DDL governs, and VAL-005 names
    the two added Evidence columns (DEC-267).
  - **MF-1, repaired:** task 18 made the fix (`b84b17c`, report `b37c169`,
    report-lint PASS). The final review must confirm it.
  - **BC-B, recorded:** the Human checklist reads 55 pass, 2 fail, 6
    partial, overall FAIL. DEC-268 (Human-directed) repairs A1 (the
    rework Stop cannot be resolved), A2 (Human-mode `claim acquire` is
    silent), and A3 (the `requirement_id` is not shown) in Phase 23. The
    plan is DEC-269 (tasks 19-21, VAL-085..089). The A1 and A3 designs are
    pending Human approval. Resume via `qa/resume-after-dec-269.md`.
- **Resolution condition (remaining):** the A1-A3 repairs land, and the
  Human re-checks the affected rows. Then one fresh independent QA Lead
  writes `final-review-1.md` covering MF-1, VAL-005, A1-A3, and the Human
  results. Do not dispatch the final review before the repairs land.
- **References:** [TRG-23-016](../phases/phase-23/triage-log.md);
  [final review](../phases/phase-23/qa/final-review.md);
  [DEC-267, DEC-268](../phases/phase-23/decisions.md); VAL-005, VAL-070..073;
  [Human checklist](../phases/phase-23/qa/human-checklist.md);
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
