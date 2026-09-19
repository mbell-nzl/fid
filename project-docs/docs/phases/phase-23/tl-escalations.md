# Phase 23 Escalations

## ESC-23-001 - Widening the journal `Condition` value breaks journal replay and Attention deduplication

- **Phase:** 23 (planning; raised before any task dispatch)
- **Status:** **RESOLVED 2026-09-20 - Human chose Option A.** Recorded as
  DEC-235; `public-interface.md` section 2's `Condition` row amended, inventory
  row 8 restated and row 10 created, tasks 9, 10, and 11 rewritten. Nothing is
  blocked; all eighteen tasks are dispatchable. The record below is kept as
  raised, with the resolution at the end.
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
is not monotonic, for three independent reasons, none addressed anywhere in the
accepted contract set or the eighteen task specs: a search of every Phase 23
document for `semantic_key`, `EventSchema`, and `omitempty` returns nothing.

### Evidence

Reproduced under `go test -overlay`, working tree untouched, with only
`currentValueShapes["Condition"]` widened.

1. **Journal replay of every pre-widening Event fails.**
   `validateCurrentShape` (`current_event_value_types.go:112`) requires an
   **exact** field-count match: `len(object) != len(fields)` rejects. A Condition
   in today's four-field form returns `Condition: closed object has wrong shape`
   against the widened shape. That value is reached on the strict path through
   `validateCurrentJSON`'s `case "stops.resolution_condition_json"`
   (`current_event_json.go:27`), and `stop.raised`, `stop.resolved`, and
   `attention.resolved` are all in `currentEventRules`
   (`current_event_catalog.go:46-48`) and behind no fence. `ReplayJournal`
   (`replay.go:153`) calls `ProjectEvent` -> `applyEventEffects`
   (`project_effects.go:11`) -> `validateCurrentEvent` for **every historical
   Event**, so rebuild-from-journal revalidates old Events against current
   shapes. Phase 23 does not bump `EventSchema`, so `UpgradeToCurrent`
   (`upgrade.go:50`) passes a Phase 22 V2 Event through unchanged and no hook
   normalizes it. Any workspace holding a Stop raised before the widening can no
   longer rebuild - which also puts VAL-004 and VAL-064 at risk at tasks 15/16.
   **Four reach points, not one** (QA confirmation, `qa-plan.md` 17.2). Besides
   the `stops.resolution_condition_json` column, `Condition` is revalidated as a
   nested field of the `Stop` value in the typed results of `stop.raised`,
   `stop.resolved`, and `run.budget-observed` (`validateCurrentResultJSON`,
   `current_event_json.go:134-141`): a Phase 22 `Stop` carrying a four-field
   Condition is rejected with `Stop: resolution_condition: Condition: closed
   object has wrong shape`. Confirmed by execution rather than by reading: the
   same recorded journal bytes rebuild `REBUILD OK: events=1 stops=1
   attention=1` on the baseline and fail under the widening as a
   `CorruptJournalError` - the operator is told the journal is **malformed**,
   not that the binary is too new. Those four sites are what Option B's upgrader
   must rewrite.
2. **Attention deduplication silently breaks across the upgrade.**
   `attentionSemanticKey` (`stop.go:422`) hashes the whole `Condition` struct
   via `MarshalCanonicalJSON`, which preserves explicit nulls. Measured: the same
   `budget-observed` condition hashes
   `sha256:8aa57e2a0374b532…` before and `sha256:3407ff3b2d279107…` after the
   three nulls are added. `attention_items.semantic_key` is a stored column, so
   after the widening `loadOpenAttentionByKey` (`stop.go:557`) no longer
   matches an Attention that was open at upgrade time, and an identical Stop
   raises a **second** open Attention. Confirmed end to end across two binaries
   sharing one workspace: an identical second raise joins the open Attention on
   the baseline (`openAttention=1`) and opens a second one under the widening
   (`openAttention=2`), both keys stored side by side. Phase 22's
   one-open-Attention-per-semantic-key invariant is lost with nothing recording
   that it was given up.
3. **Every live `stop.raise` idempotency record is invalidated.** Measured by
   the QA Lead and not in the original record (`qa-plan.md` 17.2).
   `StopRaiseRequest` (`stop.go:72-79`) embeds `ResolutionCondition Condition`,
   and `stopDigest` (`stop.go:431`) canonicalises the whole request into
   `idempotency_results.request_digest`, so widening the struct changes the
   digest of every `stop.raise`. On a workspace whose Stop was raised by the
   baseline binary, replaying the **same** idempotency key under the widened
   struct returns `idempotency conflict for stop.raise:project-… key
   "retained-test-1"` instead of the stored result: a retry that was idempotent
   before the upgrade becomes a hard refusal after it. That is a third stored
   identity the widening breaks, and a third column Option B's migration must
   recompute.

None of the three is caught by any current assertion, and none is a defect in
the plan's ownership: tasks 9, 10, and 11 own exactly the right files. The
accepted encoding itself is what does not fit.

### Why this needs the Human

Every remedy either changes an accepted public value's JSON encoding or adds
unaccepted machinery to the phase's highest-risk area. That is the same class as
DEC-223 to DEC-227, all of which the Human approved. Tech Lead planning
authority covers ownership, file scope, and recorded cause; it does not cover
changing what a Phase 23 value looks like on the wire.

### Who this reaches, and what is actually being decided

Recorded so the choice is made on the real population rather than an assumed
one (QA confirmation, `qa-plan.md` 17.4).

`schema.md` section 1 accepts task 1's edits to the shipped `0013` and `0015`
files as "a deliberate pre-adoption baseline replacement requiring fresh init
and rebuild; it is not an upgrade of old fid data", and `golang-migrate`
(`internal/db/migrate.go:8`) records only version and dirty state - it never
re-runs or checksums an applied file. Under that accepted policy there is no
released population of Phase 22 workspaces expected to run Phase 23 features,
so Option B's migration machinery is work for a population the policy says does
not exist. That does **not** dissolve this escalation:

1. **The incompatibility is in-phase, not only post-release.** The widening
   lands at tasks 9-11, and task 6 raises `review-rework-exhausted` Stops
   carrying ordinary `budget-observed` Conditions. Any workspace - developer,
   gate, or fixture - whose journal was written by a tasks 1-8 binary and is
   rebuilt by a task 9+ binary hits the measured `CorruptJournalError`. That is
   the VAL-004 and VAL-064 risk at tasks 15/16, and it is real rather than
   theoretical. One mitigating fact, checked: **no** pinned journal fixture in
   `testdata/` carries a `resolution_condition` today, so the exposure is live
   workspaces plus any fixture task 6 chooses to pin.
2. **The Attention and idempotency consequences need no rebuild at all.** Both
   are stored-column comparisons, so any workspace carrying an open Attention or
   a live `stop.raise` idempotency record across the upgrade is affected whether
   or not it ever replays.

The decision is therefore better framed as **"what encoding do we want for V1,
given that either choice costs in-phase churn"** than as "how do we protect
existing user data". Framed that way, Option A costs 31 lines and buys the
inventory row 8 invariant.

### Options

- **Option A (recommended) - the three new Condition fields are absent, not
  null, on the types they do not apply to.** `omitempty` on the **three new
  `*string` fields only** - `outcome_class`, `decision_class`, `finding_id` -
  plus per-type Condition shapes in the journal keyed on `type`, as
  `currentStopDetailShapes` is already keyed on `code` and already owned by
  tasks 6, 9, 10, and 11.

  **Read the encoding exactly; one reading of it destroys the remedy.**
  `budget_name` and `capability_id` keep today's present-and-null treatment, so
  the six Phase 22 types marshal byte-for-byte as they do now and a new type
  carries **five** fields - `type`, `subject`, `budget_name:null`,
  `capability_id:null`, and its own field - not three. The alternative reading,
  "a new type carries literally only its own field", would put `omitempty` on
  `budget_name` and `capability_id` as well; that reading was measured and it
  changes **every** existing semantic key (`38cda009…` -> `cf865c69…`,
  `da2f2a11…` -> `ef2f1a6a…`, `28607d77…` -> `9fb954bc…`), because `work-ended`
  drops from `{"budget_name":null,"capability_id":null,…}` to
  `{"subject":…,"type":"work-ended"}`. That is precisely the harm Option A
  exists to prevent. The correct reading was implemented under overlay and
  measured: the recorded pre-widening journal replays `REBUILD OK: events=1
  stops=1 attention=1`, and all three Phase 22 semantic keys stay
  **byte-identical to the baseline** (`da2f2a11…`, `38cda009…`, `28607d77…`).

  **Benefit.** Historical Events replay unchanged, and every stored
  `semantic_key` and `stop.raise` `request_digest` stays valid, so all three
  measured consequences above are eliminated. The exclusivity rule also becomes
  **structural**: measured, the wrong field is rejected with `Condition: missing
  required field "outcome_class"`, the own field null with `Condition:
  outcome_class: invalid ID`, two fields with `Condition: closed object has
  wrong shape`, and the own field alone is accepted - all four cases inventory
  row 8 leaves open, closed by the shape itself, which **removes** the relation
  work row 8 assigns to tasks 9, 10, and 11. Against the shipped corpus
  `internal/journal`, `internal/workflow`, and `internal/export` are all `ok`
  with **zero failures and zero fixture edits**.

  **Cost.** The public JSON for a Condition omits the new fields that do not
  apply to its type instead of showing them as null. That is a real amendment to
  the `public-interface.md` section 2 row, which presupposes present-and-null.
  Size measured at 16 lines of catalog map plus a 15-line dispatch branch - and
  it is **one file larger than this record first stated**. Unlike
  `currentStopDetailShapes`, which is selected at the call site from the row's
  `code` column (`current_event_json.go:35-40`), a Condition's discriminant
  lives **inside** the value, and the value is also reached generically as a
  field of `Stop` through `if shape, ok := currentValueShapes[kind]; ok`
  (`current_event_value_types.go:39`). Per-type selection therefore needs a
  `Condition` branch in `validateCurrentValue`, in
  `internal/journal/current_event_value_types.go` - a file the inventory names
  only as a *reader* of `currentStopDetailShapes` and that **no task's file
  scope lists**. Under the phase's own ownership discipline it must join the
  scope of whichever task implements the remedy.

  **One implementation constraint to carry into the decision: the lookup must
  use a checked type assertion.** The dispatch reads `object["type"]` **before**
  `validateCurrentShape` has established it is a string. Measured: the baseline
  degrades cleanly to `Condition: type: invalid closed enum
  =executor-terminated,…`, while the obvious `object["type"].(string)` panics
  with `interface conversion: interface {} is int, not string`. This code
  validates untrusted journal and archive bytes during replay and `import`, so
  the branch must return the ordinary rejection instead. One line of
  implementation; worth one sentence in the decision, because a panic here is a
  crash on a crafted archive.
- **Option B - keep the accepted encoding literally, and pay for it.** Every
  Condition always carries all seven fields. Cost, corrected upward by
  measurement (`qa-plan.md` 17.2-17.3): an `EventSchema` bump to V3, plus an
  upgrader that rewrites historical Condition JSON at **four** payload sites
  rather than one - the `stops.resolution_condition_json` column and the nested
  `Stop` inside the typed results of `stop.raised`, `stop.resolved`, and
  `run.budget-observed`; **and** a migration-0017 recompute of
  `attention_items.semantic_key` for every open Attention **and** of
  `idempotency_results.request_digest` for every live `stop.raise` record;
  **and** an accepted rule for what happens to an Attention whose key changes
  mid-life; **and** repair of the **11 failing tests and subtests** across 7
  top-level tests that the present-and-null encoding breaks in the shipped
  corpus before any Option B machinery is written
  (`TestExecutionBlockedEventProjectsIdenticallyLiveAndReplay`,
  `TestActivationRecoveryEventsProjectIdenticallyLiveAndReplay/activation.recovery-selected`,
  `TestCurrentEventRejectsConditionalBindingCorruption` - three cases x
  live/replay, `TestTerminalSuccessEvaluationIncludesEnvelopeStop`,
  `TestRecoverRunRejectsApplicableStopWithoutEffects`,
  `TestStopAndAttentionListsFilterLimitAndRemainReadOnly`,
  `TestEquivalentStopsShareAttentionButKeepDistinctFacts`). It also **keeps**
  inventory row 8's relation work, which Option A removes. In its favour: it
  preserves the accepted value exactly, it keeps one shape per value kind as
  every other `currentValueShapes` entry does, and it is genuinely
  implementable - `ReplayJournal` hashes the stored Event before upgrading it
  (`replay.go:141`, then `:147`), so a payload-rewriting upgrader does not break
  hash chains. It simply costs four to five times Option A, in tasks 1 and 15,
  the two tasks already carrying the most migration risk.
- **Option C - defer the three fields to Phase 24.** Ship the four new Condition
  **types** (a safe enum widening) with no new fields, and carry the class and
  Finding discriminators in Stop details, where `decision_class` and
  `outcome_class` already exist for `protected-decision-required` and
  `human-outcome-required`. Rejected on inspection: the Attention semantic key is
  derived from `(subject, code, condition)` only, so with no class field in the
  Condition, two human-outcome Attentions of different classes on the same
  Envelope produce the same key, `loadOpenAttentionByKey` joins them into one
  Attention, and resolving one class closes the Attention covering both.
  Verified by the QA Lead, who adds that the same collision applies a second
  time to `decision-satisfied` and `decision_class`. Recorded here only so it is
  not re-derived.

### Resolution as decided

**The Human chose Option A on 2026-09-20, in the corrected reading**: `omitempty`
on the three new fields only, the six Phase 22 types keeping `budget_name` and
`capability_id` present-and-null, with per-type journal shapes keyed on `type`.
Recorded as DEC-235, which carries the full disposition. Applied in this batch:

- `public-interface.md` section 2's `Condition` row now states the encoding, the
  five-field judgment Condition, the untouched four-field Phase 22 Condition,
  and that exclusivity is structural. (*Corrected by DEC-236 / R-26: only the
  three judgment types that have a discriminator carry five fields;
  `review-recorded` carries four. `public-interface.md` itself is correct - this
  summary line is not, and is historical.*)
- `execution-scope.md` section 6: inventory **row 8 restated, not deleted** -
  the relation branch keeps the Phase 22 `budget_name` / `capability_id`
  binding and **no Phase 23 task extends it**. New **row 10** is the per-type
  shape map and its dispatch, failing closed, with task 9 creating it and owning
  `internal/journal/current_event_value_types.go`, and tasks 10 and 11 adding
  rows. *(Corrected by DEC-236: the reason first recorded here - "which the
  shapes cannot express" - was wrong, and the branch must still be **reached**
  by task 9's new dispatch or the binding is lost silently.)*
- Tasks 9, 10, and 11: the `omitempty`-on-new-fields-only rule in all three
  (each edits the struct), the checked type assertion in task 9, byte-identity
  of the six Phase 22 rows as an acceptance bar, and the relation branch marked
  do-not-touch.
- VAL-078 keeps its assertion, method, pass and fail criteria, validator, and
  owning task; only its **evidence class** is annotated, from relation rejection
  to shape rejection.
- The three held revisions are applied with it: R-22 (proving Event named in
  tasks 9 and 10, as a direct strict-path assertion rather than a raise), R-23
  (`task-11.md` now has a VAL-078 criterion), I-37 (`schema.md` section 1).

### Original required resolution, as raised

Human chooses A, B, or C (or directs another route). Whichever is chosen, record
it as a DEC, amend `public-interface.md` section 2's `Condition` row to match,
and update tasks 9, 10, and 11 plus inventory row 8 in `execution-scope.md`
section 6 before those tasks are dispatched. Three specifics for the Option A
branch, so the disposition is not rederived: the implementing task also gains
`internal/journal/current_event_value_types.go` in its file scope and the
checked-assertion requirement in its acceptance criteria; inventory row 8 is
**restated**, not deleted, because the invariant survives and only its enforcer
moves from the relation branch to the shape; and VAL-078's evidence class
becomes shape rejection rather than relation rejection, with the assertion
itself unchanged.

Three bounded QA revisions are deliberately held for the same disposition
batch, because each is phrased against whichever encoding is chosen and each
touches only tasks that are already blocked: R-22 (name the proving Event for
the `currentStopDetailShapes` and `Condition` entries tasks 9 and 10 add, since
three of DEC-224's four raising families are retained and do not exercise them),
R-23 (`task-11.md` has no acceptance criterion citing VAL-078), and I-37 (one
recorded sentence in `schema.md` section 1 on the silent in-place-migration
divergence). None of them is a dispatch blocker.

### Scope of the block

Tasks 1 through 8 are unaffected and may dispatch on schedule: none of them
touches the `Condition` shape or struct. Only tasks 9, 10, and 11 are blocked,
and task 12 onward only through the strict task ordering. The R-19 disposition
and everything else in DEC-233 stand independently of this escalation.

Independently confirmed by the QA Lead across all eighteen specs, not only the
three this record names (`qa-plan.md` 17.5): no task before 9 touches the shape
or struct, and task 6 - the closest, because it owns the same Stop registries
and co-edits the same catalog file on unrelated lines - explicitly forbids
adding a Condition type. One consequence worth the Runner's attention: **task 12
is substantively independent of the `Condition` shape**, so if the decision
takes time, task 12 is available out of order at the Runner's discretion.

### References

- `project-docs/docs/phases/phase-23/public-interface.md` section 2, `Condition`
  row; `contracts.md` section 11.
- `project-docs/docs/phases/phase-23/execution-scope.md` section 6, inventory
  row 8.
- `project-docs/docs/phases/phase-23/decisions.md` DEC-229 (Condition types are
  registered by the task that first raises a Stop carrying them), DEC-233, and
  DEC-234 (this record's R-20 amendments and what is held for the disposition
  batch).
- `project-docs/docs/phases/phase-23/qa/qa-plan.md` sections 17.2, 17.3, 17.4 -
  the QA Lead's independent reproduction of all three consequences, the measured
  Option A implementation, and the harm-population framing.
- `internal/journal/current_event_value_types.go:39,112`;
  `internal/journal/current_event_json.go:27,35-40,134-141`;
  `internal/journal/current_event_catalog.go:46-48`;
  `internal/journal/replay.go:141,147,153`;
  `internal/journal/project_effects.go:11`;
  `internal/journal/upgrade.go:50`; `internal/db/migrate.go:8`;
  `internal/workflow/stop.go:72-79,422,431,557`.

## Blocker candidates raised without an identifier

Recorded here so the phase register is complete; neither became an escalation.

- **QA pass 7, `qa-plan.md` 18.12 - conditional on branch B.** Withdrawn by the
  QA Lead once DEC-236 chose branch A, and independently re-checked on three
  legs in `qa-plan.md` 19.11: the discriminator pairing was accepted before the
  escalation existed, the Human approved a mechanism rather than a per-type
  field assignment, and the contradicting sentence was new in the `932fbea`
  batch. Restoring the accepted reading is Tech Lead planning authority, so no
  Human decision was required.
- **QA pass 8, `qa-plan.md` 20.5 - R-28, VAL-078's granularity.** Correctly
  raised as a candidate with no identifier allocated: the protected category is
  validation-contract authority, the repair is a Tech Lead contract edit, and no
  product, encoding, or Human judgment is requested. Dispositioned under DEC-237
  by splitting the row into VAL-078 and VAL-080. A Developer could not have
  repaired it under task 11's product-file scope, which is why it was a
  candidate and not a Developer deviation.
