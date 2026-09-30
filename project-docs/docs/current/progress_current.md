# Progress — Current

Updated 2026-09-30. **Tasks 1-17 are complete, and review-12 passed window
13-15. Task 18 stopped at attempt 6 on an authority-details divergence,
which DEC-266 resolves. The Runner is ready to resume.** Branch: `phase-23`.

Task 18, attempts 2-6, pinned DEC-265's 19/18/1 partition. It also
completed the full semantic fake/fid comparison for all 23 operations, the
rejection matrix, and the VAL-067 provenance. One comparison stayed red:
`claim.acquire` returned `authority-denied` as `{operation,actor}`.
DEC-266 finds that the registry never converged three `authority-*` codes
to the Phase 21 four-field contract:
- eight producers followed the registry;
- `gate.go`'s four-field paths are rejected at construction, a latent
  defect;
- the fake's other authority rows shared the wrong shape.

The contract governs. Task 18 converges the registry and producers under a
planning amendment.

Next: the Runner resumes using
[the handoff](../phases/phase-23/qa/resume-after-dec-266.md) (task 18
attempt 7), then QA window 16-18 (review-13), then phase-end QA.

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
