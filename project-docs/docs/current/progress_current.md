# Progress — Current

Updated 2026-09-27. **REPAIR_PENDING: task 11's recovery is active; no active
STOP.** Branch: `phase-23`. Tasks 1-11 initially completed. The post-task-11
re-review `qa/checkpoint-after-task-11-2.md` (`9cb9420`) returned BLOCK with
BC-1..BC-3 and MF-4/MF-5; the resulting STOP was consumed by the accepted
DEC-256 RUN-RESUME on 2026-09-26 at 17:38.

The Tech Lead recorded DEC-249..255 (2026-09-26). The Human approved DEC-254
Part 2 (task 13 opens with a `cmd/fid` speed step) and chose Option B for BC-2
(DEC-251).

DEC-257 implemented the route for DEC-256's four ordered owners. The ledger
records tasks 6, 2, and 10 repaired and complete, then TASK-RECOVERY task-11
at 19:19. Commit `466745f` implements the task-11 repair; the existing completion
report still describes the earlier `c5bf0f8` implementation. Finish the active
recovery's exact verification and report before recording its result.

A replacement Runner uses a normal session RUN-RESUME: STOP-only fields `N/A`,
`Repair owners` omitted, `Continuing task: task-11`. Preserve the existing
recovery and completed-owner progress. The rejected repeated STOP resume did
not change the ledger. See the
[resume diagnosis](../reviews/2026-09-27-phase-runner-resume.md).
After task 11 closes, run `checkpoint-after-task-11-3.md`; task 12 stays held
until that re-review returns PASS.

## Current phase contract

Phase 23 is Heavy, approved under DEC-222, with 18 tasks and 84 validation
assertions (80 deterministic, 4 Human). The accepted contracts and task pack
are indexed in the [phase brief](phase_brief.md). Task 13 opens with the
Human-approved test-speed step (DEC-254); task 12 awaits the checkpoint PASS.

Retain these constraints while finishing the repairs:

- DEC-249: only required Obligations guarding the requested action block it.
- DEC-250/253: Stop resolution uses shared qualifying-fact predicates keyed on
  the recorded Condition; task 11 owns the ratified scope and required evidence.
- DEC-251: Human Outcomes need no Envelope grant; Activation supplies authority.
- DEC-227: exact task commands and 600-second ordinary-package / 900-second
  tagged-integration budgets remain binding. A prior implementation commit
  supplies no fresh gate evidence.
- DEC-243: VAL-025 closes with real durable Human Outcome links and all nine
  compatibility cells; the checkpoint must independently review the evidence.

## Earlier work and carry-forwards

Task 1's ESC-23-002 recovery and task 4's ESC-23-003 recovery completed before
the current checkpoint. Earlier STOP/hold instructions are historical. Their
full progress snapshot is retained in [phase history](../phases/phase-23/progress.md),
with [task 1 evidence](../phases/phase-23/qa/esc-23-002-repair.md) and
[task 4 recovery authority](../phases/phase-23/qa/resume-after-esc-23-003.md).

The eight planning/QA passes and all accepted amendments are retained in
[phase decisions](../phases/phase-23/decisions.md) and the
[QA plan](../phases/phase-23/qa/qa-plan.md). DEC-237/240's final contract has
84 assertions; the earlier 79/82 counts are historical.

Phase 22 is closed under DEC-220; DEC-221/task 20's CI performance repair is
also complete. See its [wrap audit](../phases/phase-22/wrap-up-audit.md) and
[hosted CI evidence](../phases/phase-22/qa/task-20-hosted-ci.md).
Live carry-forward owners/triggers remain in [active decisions](decisions_active.md).
Claims qualification remains Phase 24, external adapter qualification Phase 25.
