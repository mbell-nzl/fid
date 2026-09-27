# Progress — Current

Updated 2026-09-27. **STOPPED on TRG-23-006, resolved by DEC-260; ready to
resume task 12.** Branch: `phase-23`. Tasks 1-11 are complete, and
`checkpoint-after-task-11-3.md` passed (`8cd44ae`).

Task 12's real-process journey could not stage an Actor profile that held the
requirement capabilities. `public-interface.md` section 4 required those
capabilities, but section 4.1 and R-19 exclude them from the closed set. Its
expected-nonzero declaration also failed report-lint grammar (ESC-23-006).

DEC-260 made two corrections:
- Requirement writes need a trusted Actor but no capability, so the set stays
  at thirty.
- The task-12 declaration is corrected to one line.

Next: the Runner resumes task 12 at attempt 2 using
[the handoff](../phases/phase-23/qa/resume-after-esc-23-006.md). Review window
10-12 follows task 12.

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
