# TL Escalations — Current

## Open escalations

None. ESC-23-008 is resolved under DEC-263; the Runner resumes task 13 via
`qa/resume-after-esc-23-008.md`.

## Recently resolved

### ESC-23-008 - Task-13 protected-decision denial response contradicts closed failure schema

- **Status:** RESOLVED 2026-09-29 under DEC-263 (Tech Lead, Human-delegated).
  - **Cause:** section 5 copied the StopDetails of `protected-decision-required`
    and `human-outcome-required` as their failure details. It never added the
    denial fields that DEC-226 part 4 and VAL-057 require.
  - **Fix:** both codes' details gain `stop_ids`, `attention_ids`, `event`
    (null), and `replayed` (`false`), in the `review-rework-exhausted` form.
    The envelope, StopDetails, codes, precedence, VAL text, and gates are
    unchanged.
  - Task 13 continues same-task at attempt 2 (`in-scope-implementation`,
    no-progress 0/2). See the
    [full record](../phases/phase-23/tl-escalations.md).

### ESC-23-007 - Task-13 ActivityEvent journal and requirement-reference schema boundary

- **Status:** RESOLVED 2026-09-28 under DEC-262 (Tech Lead, Human-delegated).
  Task 13 owns the strict-path `activity_log` entry and the `issue`-subject
  rows; task 15 writes the `requirement`/`decision` rows with 0018 and closes
  VAL-055. See the [full record](../phases/phase-23/tl-escalations.md).

### ESC-23-006 - Task-12 requirement authority and report-lint contract conflict

- **Status:** RESOLVED 2026-09-27 under DEC-260 (Tech Lead, Human-delegated).
  Requirement writes need a trusted Actor but no Actor-profile capability or
  grant. `public-interface.md` section 4 is corrected, and the closed set
  stays at thirty. Report-lint grammar is unchanged. Task 12's expected-nonzero
  declaration is rewritten on one line. Task 12 continues same-task at
  attempt 2 (`in-scope-implementation`). See the
  [full record](../phases/phase-23/tl-escalations.md).

### ESC-23-005 - Task-11 ordinary gate budget

- **Status:** RESOLVED 2026-09-27 under DEC-258 (Human-delegated repair).
  The cause was serial test structure plus measurement spread, not task 11.
  Test-only `t.Parallel()` across the suites; gate definition unchanged. Final
  exact gates at `3dff457`: `cmd/fid` 531.533 s (later 299.6 s under DEC-259),
  `internal/workflow` 165.920 s, tagged 326.421 s. See
  [repair evidence](../phases/phase-23/qa/esc-23-005-repair.md) and the
  [full record](../phases/phase-23/tl-escalations.md).

### ESC-23-004 - Post-task-11 checkpoint re-review: three protected blocker candidates (BC-1, BC-2, BC-3)

- **Status:** RESOLVED 2026-09-26. BC-1 by DEC-249, BC-2 by DEC-251 (Human,
  Option B), BC-3 by DEC-250; carry items by DEC-252..254; repair order by
  DEC-255/256. See the [full record](../phases/phase-23/tl-escalations.md).

### ESC-23-003 - VAL-015 Run-holder-only CLI evidence is unrepresentable under the accepted lifecycle

- **Status:** RESOLVED 2026-09-20 under DEC-239 (Human-delegated through
  Task 9). Run/Claim Actor binding and durable history are retained; only the
  impossible Run-holder-only CLI evidence route became an isolated query
  fixture. See the [full record](../phases/phase-23/tl-escalations.md) and
  [resume handoff](../phases/phase-23/qa/resume-after-esc-23-003.md).

Older, all resolved and archived in
[Phase 23 escalations](../phases/phase-23/tl-escalations.md):
- ESC-23-002 under DEC-238 (task 1 gate closure;
  [repair evidence](../phases/phase-23/qa/esc-23-002-repair.md)).
- ESC-23-001 under DEC-235 (Human Option A: `omitempty` on the three new
  `Condition` fields only).
- R-28 was a blocker candidate with no identifier, correctly dispositioned
  under DEC-237 without an escalation.
- ESC-22-005 under DEC-221; see
  [Phase 22 escalations](../phases/phase-22/tl-escalations.md).

## Standing handling

Register an escalation when allocated, not after resolution. Ordinary causal
repair and incomplete in-scope work do not themselves require a Human stop.
Read accepted ownership before calling it ambiguous (DEC-194). A real contract/
schema contradiction needs the missing authority, not a weakened assertion
(DEC-195). External interruptions do not consume repair no-progress budget.

Prior archives: [Phase 21](../phases/phase-21/tl-escalations.md),
[Phase 20](../phases/phase-20/tl-escalations.md),
[Phase 19](../phases/phase-19/tl-escalations.md).
