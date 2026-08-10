# Phase 20 - TL-Lite Triage Log

> **Purpose.** Append-only log of mid-phase triage events handled by TL-Lite.
> Phase Runner reads TL-Lite's `VERDICT:` line to decide whether to resume or
> stop. Tech Lead audits this log at phase wrap-up.
>
> **Numbering.** Identifiers are phase-scoped: `TRG-20-NNN` and `ESC-20-NNN`.
>
> **This file is the allocation source for both.** Read this file, take the
> highest `NNN` already recorded in that series, plus one; start at `001` when
> empty. Never recall a number from memory, never renumber a published entry,
> and never emit a bare placeholder.
>
> **Write the entry before reporting the verdict.** Newest entry at top.

| TRG-20-NNN | Date | Trigger | Source | Verdict | ESC-20-NNN (if BLOCK) | Summary |
|---|---|---|---|---|---|---|
| TRG-20-004 | 2026-08-10 | Infrastructure Review 1 found a retained integration-tagged in-process test outside both accepted gates | `qa/infra-review-1.md` | BLOCK | ESC-20-004 | Task 3f's direct-parent inventory omitted `internal/issues/transition_rules_test.go`; repository-wide classification and ordinary ownership are required before Task 4. |
| TRG-20-003 | 2026-08-09 | Task 3d blocker report: required full-suite coverage gate has no captured terminal output and report-lint fails on three go-test prose false matches | `tasks/task-3d-done.md` | BLOCK | ESC-20-003 | The required gate is UNVERIFIED, not PASS; the new report-lint control also has an unresolved three-failure false-match defect. |
| TRG-20-002 | 2026-08-09 | Bounded QA re-review review-1b returned BLOCK with two unresolved findings before Task 4 | `qa/review-1b.md` | BLOCK | ESC-20-002 | VAL-043 has an exhaustive evidence gap and the required integration gate fails in three init scenarios. |
| TRG-20-001 | 2026-08-09 | QA Review 1 has two unresolved Blocking evidence gaps before Task 4 can start | `qa/review-1.md` | BLOCK | ESC-20-001 | VAL-002 and VAL-043 lack required exhaustive boundary evidence across Tasks 1-3. |

---

## TRG-20-004 - Infrastructure Review 1 found an orphaned tagged test

- **Trigger:** The fresh infrastructure review found
  `internal/issues/transition_rules_test.go` still tagged while Task 3f's
  canonical dedicated command runs only `cmd/fid/integration`; the ordinary
  command excludes the file and the classification regression scanned only
  the direct `cmd/fid` parent.
- **Source:** `project-docs/docs/phases/phase-20/qa/infra-review-1.md`
- **Verdict:** BLOCK
- **Reasoning:** This is a cross-task gate-definition regression against
  DEC-133's complete tagged-test accounting. The test's real persistence
  assertion is valid, but no accepted command executes it. The Human-approved
  pack explicitly requires blocking findings to be corrected and returned to
  the same reviewer before Task 4.
- **Local decisions made:** none; correction follows DEC-133's existing rule
  that in-process tests belong to the ordinary suite and that every tagged
  test must be inventoried.
- **References:** `qa/infra-review-1.md`; `decisions.md` DEC-133;
  `infra-improvement-pack.md` Invariants and Final Evidence Contract.
- **If BLOCK:** ESC-20-004 raised in `tl_escalations_current.md`; correction
  commit `75c856a` awaits complete post-fix gates and same-reviewer re-review.

## TRG-20-003 - Task 3d blocked by uncaptured coverage gate and report-lint false matches

- **Trigger:** Task 3d Developer produced a blocker report because the required
  `go test -race -coverprofile=coverage.out ./...` terminated in the session
  without captured terminal output, making its gate UNVERIFIED rather than
  PASS; independent report-lint then failed three checks because the new
  go-test token control matched explanatory `go test` prose in the blocker
  report. Task 3d has no commit.
- **Source:** `project-docs/docs/phases/phase-20/tasks/task-3d-done.md`
- **Verdict:** BLOCK
- **Reasoning:** This is a genuine unresolved required-gate and report-lint
  hard stop. AGENTS.md and DEC-131 require a real terminal token for a
  full-suite `go test` PASS and classify uncaptured output as UNVERIFIED, not
  success. The report-lint failure is a second genuine hard stop; no attached
  verified commit resolves either condition, and DEC-131 requires the new
  control to land rather than remain a promise. The resulting partial-work
  recovery changes delivery tooling and evidence across Task 3d and the five
  corrected completion reports, exceeding local-decision authority.
- **Local decisions made:** none
- **References:** `project-docs/docs/phases/phase-20/tasks/task-3d-done.md`
  (Exact quality-gate results and blocker); `AGENTS.md` (Commands, evidence
  token and UNVERIFIED rules); `project-docs/docs/phases/phase-20/decisions.md`
  (DEC-131, decisions 4-6 and 9); `project-docs/specs/09_AI_WORKFLOW.md` §2.3
  and §9.
- **If BLOCK:** ESC-20-003 raised in `tl_escalations_current.md`; recommended
  next action: invoke full Tech Lead to resolve ESC-20-003 before any Task 3d
  recovery or review-1c work resumes.

## TRG-20-002 - Review 1b unresolved VAL-043 evidence and integration-gate failure

- **Trigger:** bounded QA re-review review-1b returned BLOCK with two new
  unresolved findings before Task 4: VAL-043 lacks explicit rejected
  lower-bound/type/overflow cases for `executor_limits.max_excerpt_bytes`, and
  the exact required integration command fails in three `cmd/fid init`
  scenarios with invalid project initialization argument, contradicting
  task-3b's recorded PASS gate.
- **Source:** `project-docs/docs/phases/phase-20/qa/review-1b.md`
- **Verdict:** BLOCK
- **Reasoning:** These are QA Reviewer Blocking findings and a required-gate
  failure. They are not misclassified, no supplied verified commit resolves
  them, and no accepted decision covers them: DEC-130 applies only to the
  earlier ESC-20-001 recovery and does not weaken or close VAL-043. The repair
  affects strict-config evidence and the integration journey across subsystems,
  so it exceeds local-decision authority.
- **Local decisions made:** none
- **References:** `project-docs/docs/phases/phase-20/qa/review-1b.md`
  (VAL-043 and Task 3b Blocking findings);
  `project-docs/docs/current/tl_escalations_current.md` (ESC-20-001 archive
  note); `project-docs/docs/phases/phase-20/decisions.md` (DEC-130);
  `project-docs/specs/09_AI_WORKFLOW.md` §2.3.
- **If BLOCK:** ESC-20-002 raised in `tl_escalations_current.md`; recommended
  next action: invoke full Tech Lead to resolve ESC-20-002 before Task 4
  starts.

## TRG-20-001 - Review 1 unresolved cross-task evidence gaps

- **Trigger:** QA Review 1 has two unresolved Blocking evidence gaps before
  Task 4 can start: Task 1 does not exhaustively prove exact detail
  value/null/nested schemas and redaction behavior for every registered reason
  code; Tasks 2-3 do not prove the required strict-config and init boundary
  matrix, including one-character/128-byte/MaxInt64 and overflow
  just-inside/outside cases.
- **Source:** `project-docs/docs/phases/phase-20/qa/review-1.md`
- **Verdict:** BLOCK
- **Reasoning:** These are QA Reviewer Blocking findings that are neither
  misclassified, resolved by the supplied verified review commit, nor covered
  by an accepted decision. They leave VAL-002 and VAL-043 without the required
  qualifying/non-qualifying evidence and span completed Tasks 1-3.
- **Local decisions made:** none
- **References:** `project-docs/docs/phases/phase-20/qa/review-1.md` (Task 1
  and Task 2 Blocking findings); `project-docs/docs/phases/phase-20/qa/validation-contract.md`
  (VAL-002, VAL-043); `project-docs/specs/09_AI_WORKFLOW.md` \u00a72.3.
- **If BLOCK:** ESC-20-001 raised in `tl_escalations_current.md`; recommended
  next action: invoke full Tech Lead to resolve the cross-task evidence
  recovery before Task 4 starts.

<!--
Phase wrap-up checklist:

- [ ] Walked every TRG-20-NNN above.
- [ ] Confirmed each RESUME verdict was correct in hindsight.
- [ ] Flagged any TRG that should have escalated.
- [ ] Confirmed every entry carries a real, unique, non-placeholder ID, and
      that every BLOCK names the ESC it raised.
- [ ] Confirmed autonomous TL-Lite resume in effect (DEC-002; no warm-up tally).
-->
