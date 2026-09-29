# Progress — Current

Updated 2026-09-30. **Tasks 1-14 are complete, and review-10 passed window
10-12. Task 15 stopped on ESC-23-009, which DEC-264 resolves. The Runner is
ready to resume.** Branch: `phase-23`.

Task 13 closed under DEC-263. Task 14 closed after TRG-23-009's bounded
generic Concern-list repair. Task 15's first continuation left a partial
checkpoint (`fa6355a`) that includes 0018 up. It blocked because `schema.md`
section 5 restored 0002 `concerns` with `escalated_to_esc_id`, which 0016
down adds again. DEC-264 changes section 5 to restore the schema-17 state,
fixing the view source and the missing kind triggers as well. The Tech Lead
verified the section 5 steps through the full down chain on a scratch
database.

Next: the Runner resumes using
[the handoff](../phases/phase-23/qa/resume-after-esc-23-009.md) (task 15
attempt 2), then QA window 13-15, then tasks 16-18.

Earlier, under ESC-23-007, task 13 attempt 1 added `review.list`, the
failure codes, and request closure. Its ActivityEvent effects could not commit, for two reasons:
- `activity_log` is missing from the strict journal table registry, which
  was a planning gap.
- The 0016 trigger cannot resolve `requirement` or `decision` refs until
  task 15's 0018.

DEC-262 splits the work at 0018:
- Task 13 registers the table and emits the `issue`-subject names.
- Task 15 emits the `requirement`/`decision` rows and closes VAL-055.

Earlier, review-9's two Must-fix repairs closed under DEC-261, and
review-10 recorded PASS.
Earlier, ESC-23-006 was resolved by DEC-260. Under DEC-260, task 12 closed
at `b763127`. Requirement writes need a trusted Actor but no capability.

Earlier: ESC-23-005 was resolved by DEC-258/259, which closed task 11.
Earlier history: DEC-249..257 dispositioned review-7 (ESC-23-004). The DEC-256
owners (tasks 6, 2, and 10) are repaired and complete.

## Current phase contract

Phase 23 is Heavy, approved under DEC-222, with 18 tasks and 84 validation
assertions (80 deterministic, 4 Human). The accepted contracts and task pack
are indexed in the [phase brief](phase_brief.md). DEC-258 satisfies DEC-254's
test-speed step.

Retain these constraints:

- DEC-249: only required Obligations guarding the requested action block it.
- DEC-250/253: Stop resolution uses shared qualifying-fact predicates keyed on
  the recorded Condition; task 11 owns the ratified scope and required evidence.
- DEC-251: Human Outcomes need no Envelope grant; Activation supplies authority.
- DEC-227: exact task commands and 600-second ordinary-package / 900-second
  tagged-integration budgets remain binding. A prior implementation commit
  supplies no fresh gate evidence.
- DEC-262: task 13 writes no `requirement`/`decision` activity row;
  task 15 does, alongside 0018's trigger.
- DEC-263: Stop-blocked denial fields live in failure `details` only.
  StopDetails are unchanged. Wrap-up note: DEC-263 does not decide how a
  protected Decision that becomes unsatisfied again after its Stop resolved
  (Outcome supersession) is reported.
- DEC-264: 0018 down follows `schema.md` section 5 exactly. It restores the
  schema-17 state and edits no migration from 0001 to 0017.
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
