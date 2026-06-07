# Phase 1 - TL Escalations (Archived)

> Resolved escalations moved here from `tl_escalations_current.md` at Phase 1
> wrap-up. ESC-NNN numbering is global and never resets.

## ESC-001: SQLite pragmas may not apply to every pooled connection

- **Raised:** 2026-06-06
- **By:** QA Reviewer
- **Trigger:** QA Reviewer trigger 5 - spans tasks or subsystems beyond the review window.
- **Context:** Phase 1 Task 1 implemented `internal/db.Open` by running
  `PRAGMA foreign_keys = ON` and `PRAGMA recursive_triggers = ON` on the
  returned `*sql.DB`. That proves only the physical connection used for those
  Exec calls, but `database/sql` may open later physical SQLite connections for
  the same pool without those per-connection settings. VAL-014 requires every
  connection path to report both pragmas ON, and later replay/init/issue tasks
  all consume this DB helper.
- **Question / proposed direction:** Decide whether to enforce pragmas through
  per-physical-connection DSN/connector configuration or deliberately constrain
  the pool; add evidence that multiple opened physical connections report both
  pragmas ON before Phase 1 sign-off.
- **Links:** `qa/review-1.md`, `tasks/task-1.md`, `internal/db/db.go`.
- **Disposition (2026-06-07, Tech Lead):** Do-now in Phase 1. Resolved by
  DEC-006 (`decisions.md`): enforce both pragmas via the `modernc.org/sqlite`
  `_pragma` DSN parameters so every physical connection is configured by the
  driver, proven by a multi-connection regression test. Implemented as Task 7
  (`tasks/task-7.md`).
- **Resolution:** CLOSED at Phase 1 wrap-up. Task 7 landed with passing
  multi-connection VAL-014 evidence; QA Review 2 verdict CONTINUE; final review
  records VAL-014 PASS. TRG-003 was the controlling triage entry.
