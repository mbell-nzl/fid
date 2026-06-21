# Phase 11 Triage Log

## TRG-001 - Claim-write projection refresh lock ordering

- **Date:** 2026-06-21
- **Trigger:** Task 4 completion report notes claim-write lock ordering during
  claim writes. Source: task report.
- **Disposition:** BLOCK.
- **Escalation:** ESC-008 in
  `project-docs/docs/current/tl_escalations_current.md`.
- **Resolution:** Tech Lead resolved through DEC-028. Phase 11 adds Task 5 and
  VAL-015 before final QA; ESC-008 archived in
  `project-docs/docs/phases/phase-11/tl-escalations.md`.
- **Reason:** The note concerns `internal/claims.Acquire` and related claim
  writes calling `journal.EnsureProjectionFresh` inside the claim write-lock
  path. That crosses `internal/claims`, journal/projection freshness, and
  concurrency behavior, and Task 4 itself is documented as a cross-subsystem
  task touching issues, claims, workflow categories, journal effects, and CLI
  behavior.
- **Local decisions made:** None; disposition only.
- **References:** `project-docs/docs/phases/phase-11/tasks/task-4-done.md`;
  `project-docs/docs/phases/phase-11/tasks/task-4.md`;
  `project-docs/docs/phases/phase-11/qa/validation-contract.md`.
