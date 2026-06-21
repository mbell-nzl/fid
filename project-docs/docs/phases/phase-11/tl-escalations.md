# Phase 11 TL Escalations

## ESC-008: Claim-write projection refresh lock ordering

- **Raised:** 2026-06-21
- **Resolved:** 2026-06-21
- **By:** TL-Lite triage (TRG-001), resolved by Tech Lead
- **Trigger:** Affects more than one task or subsystem.
- **Context:** Phase 11 Task 4 completion report noted that
  `internal/claims.Acquire` and related claim writes still called
  `journal.EnsureProjectionFresh` from inside the claim write-lock path. Task
  fixtures surfaced nested-lock failures when the journal was ahead of
  projection.
- **Decision:** Not acceptable as-is for Phase 11 final QA. DEC-028 requires a
  do-now Task 5 to make claim write paths validate against fresh projection
  state without recursively acquiring `.fid/write.lock`.
- **Disposition:** Resolved to do-now implementation task, not carry-forward.
- **Required follow-up:** Complete phase-11/task-5, then Review 5 must verify
  VAL-015 before final QA.
- **References:** DEC-028; `project-docs/docs/phases/phase-11/tasks/task-5.md`;
  `project-docs/docs/phases/phase-11/qa/validation-contract.md` VAL-015;
  `project-docs/docs/phases/phase-11/triage-log.md` TRG-001.
