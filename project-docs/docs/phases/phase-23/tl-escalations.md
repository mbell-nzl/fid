# Phase 23 Escalations

## ESC-23-001 - Widening the journal `Condition` value breaks journal replay and Attention deduplication

- **Phase:** 23 (planning; raised before any task dispatch)
- **Status:** OPEN - blocks dispatch of tasks 9, 10, and 11 only
- **Raised by:** Tech Lead, while dispositioning R-19 under the widened
  inventory frame the QA Lead asked for in I-33
- **Category:** accepted public-value encoding, with migration and
  durable-replay consequences
- **Decision needed from:** Human

### Summary

`public-interface.md` section 2 amends the `Condition` value to add
`outcome_class:?ID`, `decision_class:?ID`, and `finding_id:?ID`, with "only the
field matching the type is non-null". Tasks 9, 10, and 11 each add one of those
fields to `internal/workflow/stop.go`'s `Condition` struct and to
`currentValueShapes["Condition"]`.

`Condition` is the **only** field-adding widening of a live journal value shape
in this phase - every other amendment in section 2 either widens a closed enum
(`Claim.mode`, `Budget.unit`, `Budget.exhaustion_action`, `StopCode`) or lands
on a family in `retainedCurrentRules`, which the retained branch returns from at
`current_event_validation.go:64` before any value validation. Adding a **field**
is not monotonic, for two independent reasons, and neither is addressed anywhere
in the accepted contract set or the eighteen task specs: a search of every Phase
23 document for `semantic_key`, `EventSchema`, and `omitempty` returns nothing.

### Evidence

Reproduced under `go test -overlay`, working tree untouched, with only
`currentValueShapes["Condition"]` widened.

1. **Journal replay of every pre-widening Event fails.**
   `validateCurrentShape` (`current_event_value_types.go:112`) requires an
   **exact** field-count match: `len(object) != len(fields)` rejects. A Condition
   in today's four-field form returns `Condition: closed object has wrong shape`
   against the widened shape. That value is reached on the strict path through
   `validateCurrentJSON`'s `case "stops.resolution_condition_json"`
   (`current_event_json.go:26`), and `stop.raised`, `stop.resolved`, and
   `attention.resolved` are all in `currentEventRules`
   (`current_event_catalog.go:46-48`) and behind no fence. `ReplayJournal`
   (`replay.go:153`) calls `ProjectEvent` -> `applyEventEffects`
   (`project_effects.go:11`) -> `validateCurrentEvent` for **every historical
   Event**, so rebuild-from-journal revalidates old Events against current
   shapes. Phase 23 does not bump `EventSchema`, so `UpgradeToCurrent`
   (`upgrade.go:50`) passes a Phase 22 V2 Event through unchanged and no hook
   normalizes it. Any workspace holding a Stop raised before the widening can no
   longer rebuild - which also puts VAL-004 and VAL-064 at risk at tasks 15/16.
2. **Attention deduplication silently breaks across the upgrade.**
   `attentionSemanticKey` (`stop.go:422`) hashes the whole `Condition` struct
   via `MarshalCanonicalJSON`, which preserves explicit nulls. Measured: the same
   `budget-observed` condition hashes
   `sha256:8aa57e2a0374b532…` before and `sha256:3407ff3b2d279107…` after the
   three nulls are added. `attention_items.semantic_key` is a stored column, so
   after the widening `openAttentionBySemanticKey` (`stop.go:558`) no longer
   matches an Attention that was open at upgrade time, and an identical Stop
   raises a **second** open Attention. Phase 22's one-open-Attention-per-semantic-key
   invariant is lost with nothing recording that it was given up.

Neither effect is caught by any current assertion, and neither is a defect in
the plan's ownership: tasks 9, 10, and 11 own exactly the right files. The
accepted encoding itself is what does not fit.

### Why this needs the Human

Every remedy either changes an accepted public value's JSON encoding or adds
unaccepted machinery to the phase's highest-risk area. That is the same class as
DEC-223 to DEC-227, all of which the Human approved. Tech Lead planning
authority covers ownership, file scope, and recorded cause; it does not cover
changing what a Phase 23 value looks like on the wire.

### Options

- **Option A (recommended) - non-applicable Condition fields are absent, not
  null.** `omitempty` on the three new `*string` fields, and per-type Condition
  shapes in the journal keyed on `type`, exactly as `currentStopDetailShapes` is
  already keyed on `code` and already owned by tasks 6, 9, 10, and 11. The six
  Phase 22 types keep their current shape byte-for-byte, so historical Events
  replay unchanged and every existing semantic key is stable; each new type
  carries only its own field, which makes the exclusivity rule **structural** and
  closes inventory row 8 for free. Cost: the public JSON for a Condition omits
  the fields that do not apply to its type, instead of showing them as null.
  That is a real amendment to the section 2 row, which currently presupposes
  present-and-null.
- **Option B - keep the accepted encoding literally, and pay for it.** Every
  Condition always carries all seven fields. This requires an `EventSchema` bump
  to V3 plus an upgrader that rewrites historical Condition JSON during replay,
  **and** a migration-0017 recompute of `attention_items.semantic_key` for every
  open Attention, **and** an accepted rule for what happens to an Attention whose
  key changes mid-life. It preserves the accepted value exactly and lands the
  cost in tasks 1 and 15, the two tasks already carrying the most migration risk.
- **Option C - defer the three fields to Phase 24.** Ship the four new Condition
  **types** (a safe enum widening) with no new fields, and carry the class and
  Finding discriminators in Stop details, where `decision_class` and
  `outcome_class` already exist for `protected-decision-required` and
  `human-outcome-required`. Rejected on inspection: the Attention semantic key is
  derived from the Condition, so two human-outcome Attentions of different
  classes would collide on one key. Recorded here only so it is not re-derived.

### Required resolution

Human chooses A, B, or C (or directs another route). Whichever is chosen, record
it as a DEC, amend `public-interface.md` section 2's `Condition` row to match,
and update tasks 9, 10, and 11 plus inventory row 8 in `execution-scope.md`
section 6 before those tasks are dispatched.

### Scope of the block

Tasks 1 through 8 are unaffected and may dispatch on schedule: none of them
touches the `Condition` shape or struct. Only tasks 9, 10, and 11 are blocked,
and task 12 onward only through the strict task ordering. The R-19 disposition
and everything else in DEC-233 stand independently of this escalation.

### References

- `project-docs/docs/phases/phase-23/public-interface.md` section 2, `Condition`
  row; `contracts.md` section 11.
- `project-docs/docs/phases/phase-23/execution-scope.md` section 6, inventory
  row 8.
- `project-docs/docs/phases/phase-23/decisions.md` DEC-229 (Condition types are
  registered by the task that first raises a Stop carrying them), DEC-233.
- `internal/journal/current_event_value_types.go:112`;
  `internal/journal/current_event_json.go:26`;
  `internal/journal/current_event_catalog.go:46-48`;
  `internal/journal/replay.go:153`; `internal/journal/project_effects.go:11`;
  `internal/journal/upgrade.go:50`; `internal/workflow/stop.go:422,558`.
