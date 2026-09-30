# TL Escalations — Current

## Open escalations

### ESC-23-010 - Task 18 Failure-code count contradicts the named public set

- **Status:** OPEN 2026-09-30; governing triage TRG-23-011.
- **Category:** protected public-interface and specification contradiction.
- **Source:** `project-docs/docs/phases/phase-23/tasks/task-18-done.md`; task-18 closure and VAL-068 remain blocked.
- **Conflict:** `public-interface.md` section 5 names nineteen distinct added Failure codes while its prose and accepted DEC-230 require exactly eighteen. The existing nineteen-row task-18 corpus reflects the named table. No local choice of a member or new count is authorized.
- **Decision needed:** The contract owner must state the exact accepted member set and count, then reconcile the protected public-interface text and DEC-230 under the proper authority. Keep the original task-18 audit and all exact gate results as evidence.
- **Independent ordinary work:** The shared fake/fid audit compares fields and replay but not complete semantic result values. Once the protected contract is resolved, task 18 continues within its existing scope to prove VAL-068 using independent fake expectations, reruns the unchanged required verification, updates its report, and obtains affected independent review. No contract choice or VAL-068 weakening is authorized by this escalation.
- **Next:** invoke full Tech Lead with prompt: "Resolve ESC-23-010 in tl_escalations_current.md."
- **References:** [TRG-23-011](../phases/phase-23/triage-log.md), [DEC-230](../phases/phase-23/decisions.md), [public interface section 5](../phases/phase-23/public-interface.md), [task-18 report](../phases/phase-23/tasks/task-18-done.md).

## Recently resolved

### ESC-23-009 - Migration 0018 down conflicts with the accepted 0016 rollback

- **Status:** RESOLVED 2026-09-30 under DEC-264 (Tech Lead proposal,
  Human-approved).
  - **Cause:** `schema.md` section 5 restored objects as they were first
    created, not as they stood at schema 17. As a result:
    - V1 `concerns` kept `escalated_to_esc_id`, so 0016 down failed with a
      duplicate column;
    - `issues_enriched` cited 0003 instead of 0011;
    - the two 0004 kind triggers were missing.
  - **Fix:** 0018 down restores the fresh-schema-17 definition of each object.
    Section 5 now lists nine exact steps, and the Tech Lead verified them
    through the full down chain on a scratch database. Migrations 0001-0017,
    0018 up, VAL text, and all commands are unchanged.
  - Task 15 continues same-task at attempt 2 (`in-scope-implementation`,
    no-progress 0/2). See the
    [full record](../phases/phase-23/tl-escalations.md).

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

Older, all resolved and archived in
[Phase 23 escalations](../phases/phase-23/tl-escalations.md):
- ESC-23-004 (checkpoint BC-1..3) under DEC-249..256.
- ESC-23-003 under DEC-239 (VAL-015 evidence route; see the
  [resume handoff](../phases/phase-23/qa/resume-after-esc-23-003.md)).
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
