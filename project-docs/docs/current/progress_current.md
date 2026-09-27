# Progress — Current

Updated 2026-09-28. **Tasks 1-12 are complete, and review window 10-12 is
in QA REPAIR. The Runner is ready to resume.** Branch: `phase-23`.

Review-9 (`3dbf5c6`) found two bounded Must-fix items and no blocker
candidates:
- MF-1: task 10, VAL-046 artifact-producer evidence.
- MF-2: task 12, VAL-053 portable-fake `requirement.synced`.

The previous Runner session ended before it recorded REVIEW. The ledger is
still valid and advanceable.

DEC-261 makes two changes:
- Reopened task 10 runs the DEC-259 race gates.
- Both repairs are bounded.

Next: the Runner resumes using
[the handoff](../phases/phase-23/qa/resume-after-review-9.md). It repairs
task 10, then task 12. A fresh re-review follows in `review-10.md`, which
must PASS before task 13.

Earlier, ESC-23-006 was resolved by DEC-260. Under DEC-260, task 12 closed
at `b763127`. Requirement writes need a trusted Actor but no capability.

Earlier: ESC-23-005 was resolved by DEC-258/259, which closed task 11.
Earlier history: DEC-249..257 dispositioned review-7 (ESC-23-004). The DEC-256
owners (tasks 6, 2, and 10) are repaired and complete.

## Current phase contract

Phase 23 is Heavy, approved under DEC-222, with 18 tasks and 84 validation
assertions (80 deterministic, 4 Human). The accepted contracts and task pack
are indexed in the [phase brief](phase_brief.md). DEC-258 satisfies DEC-254's
test-speed step; task 12 awaits the checkpoint PASS.

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
