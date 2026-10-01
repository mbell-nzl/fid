# Progress — Current

Updated 2026-10-01. **All 18 tasks are complete. The phase-end final
review (`74182db`) returned BLOCK under TRG-23-016 / ESC-23-012. DEC-267
disposed it, and task 18 has since repaired MF-1 (`b84b17c`, `b37c169`).
The Human checklist is recorded as FAIL (55/2/6). DEC-268 repairs A1-A3 in
Phase 23.** Branch: `phase-23`.

- **MF-2 / BC-A (VAL-005):** DEC-267 keeps the two-column DDL and corrects
  the criterion text. No implementation is needed.
- **MF-1:** task 18 repaired the full-checkptr golden child, and its report
  passes report-lint. The final review must confirm it.
- **Human checklist (VAL-070..073):** HUMAN-03b-resolution and
  HUMAN-03b-internals FAIL, plus 6 PARTIAL. A1-A8 are analysed in
  [the checklist](../phases/phase-23/qa/human-checklist.md).
  - **A1:** the `review-rework-exhausted` Stop cannot be resolved in
    Phase 23.
  - **A2:** Human-mode `claim acquire` prints nothing on that denial.
  - **A3:** `RequirementRef` does not show the stable `requirement_id`.

Next: DEC-269 plans the A1-A3 repairs as supplemental tasks 19 (A2), 20
(A3), and 21 (A1), adding VAL-085..089 (89 rows). The A1 resolution design
and the A3 interface amendment **await Human approval**. Once DEC-269 reads
ACCEPTED, the Runner resumes using
[the handoff](../phases/phase-23/qa/resume-after-dec-269.md). The order is:
tasks 19-21, then `review-15`, then a fresh QA Lead `final-review-1.md`
(recorded as review-16, which regenerates the fixtures and authors
`human-checklist-2.md`), then the Human re-check, then RUN-END.

**Wrap-up agenda (DEC-268):** the Human and Tech Lead discuss the remaining
Human operating-experience scope and decide when to plan it. That scope is
A4-A8: Stop details rendered as raw JSON, denial reasons, hint lines, the
`finding address` reason, and `--input` authoring and caller-supplied
repository/digest values. It also includes the leftover PARTIAL rows. No
carry-forward is recorded before that discussion.

Earlier, task 18 converged the `authority-*` details under DEC-266,
completed after ESC-23-011's model-capacity stop, and closed review-13 MF-2.
Earlier, DEC-265 clarified 19 declared Failure codes (18 added and 1
retained).
Earlier, task 15 closed under DEC-264 (0018 down restores schema 17), and
tasks 16 and 17 closed without escalation.
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
- DEC-265: section 5 declares 19 Failure codes (18 added, 1 retained).
  Corpora cover all 19; claims about the added set count 18.
- DEC-266: all four `authority-*` codes carry the Phase 21 details
  `{operation,actor_id,required_actions,required_capabilities}`, null subject.
- DEC-243: VAL-025 closes with real durable Human Outcome links and all nine
  compatibility cells; the checkpoint must independently review the evidence.
- DEC-267: VAL-005's criterion names the two added Evidence columns. MF-1
  is a task-18 test-harness repair; the 5000 ms deadline, golden bytes, and
  commands are preserved.
- DEC-268: the Human directs A1-A3 to be repaired in Phase 23. DEC-269
  allocates tasks 19-21; its A1 and A3 designs need Human approval. The rest of the Human operating-experience
  scope is on the wrap-up agenda.

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
