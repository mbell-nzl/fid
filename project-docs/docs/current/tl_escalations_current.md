# TL Escalations — Current

## Open escalations

None. The Runner resumes task 18 from
[the handoff](../phases/phase-23/qa/resume-after-esc-23-010.md).

## Recently resolved

### ESC-23-010 - Task 18 Failure-code count contradicts the named public set

- **Status:** RESOLVED 2026-09-30 under DEC-265 (Tech Lead, Human-delegated).
  - **Cause:** there was no member conflict. Section 5, section 5.1,
    `p23FailureContracts`, and `failures.json` all name the same 19 codes.
    DEC-230's "18" counts the codes new to `codes.go`, and
    `requirement-not-found` is a retained Phase 20 code whose Phase 23
    producers section 5 re-pins. The contradiction came from section 5's
    lead sentence, which called all 19 "added".
  - **Fix:** a wording-only clarification that 19 are declared, 18 added
    and 1 retained. Codes, VAL text, ownership, and commands are unchanged.
  - Task 18 continues same-task at attempt 2 (`in-scope-implementation`,
    no-progress 0/2) to finish VAL-068 semantic parity and VAL-067
    provenance. See the [full record](../phases/phase-23/tl-escalations.md).

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

Older, all resolved and archived in
[Phase 23 escalations](../phases/phase-23/tl-escalations.md):
- ESC-23-006 under DEC-260 (requirement writes need a trusted Actor only;
  closed set stays thirty).
- ESC-23-005 under DEC-258 (parallel tests; gate definition unchanged; see
  [repair evidence](../phases/phase-23/qa/esc-23-005-repair.md)).
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
