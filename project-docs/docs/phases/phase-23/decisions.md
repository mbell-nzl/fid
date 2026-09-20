# Phase 23 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes,
> and wrap-up decisions. DEC numbering is global and never resets.
>
> Status: accepted, with one open Human decision. Phase 23 is approved and not
> yet started. Five QA passes have returned nineteen Must-fix revisions -
> thirteen first-pass, R-14/R-15 on the readiness refresh, R-16 on the DEC-231
> confirmation, R-17/R-18 on the DEC-232 confirmation, and R-19 on the DEC-233
> confirmation - and all are dispositioned and applied, together with five the
> Tech Lead found while dispositioning them. All five contract amendments,
> DEC-223 through DEC-227, are approved by the Human on 2026-09-19, DEC-226
> together with the part 5 vocabulary fix added during that review. DEC-229
> through DEC-233 are Tech Lead planning dispositions. **One decision is open:
> ESC-23-001**, how the journal `Condition` value may gain fields, raised by the
> Tech Lead under DEC-233; it blocks tasks 9, 10, and 11 only. Next: a fresh QA
> Lead confirmation of the DEC-233 edits, the Human's ESC-23-001 choice, then the
> Runner's readiness check.

## Decisions

### DEC-222: Accept the Phase 23 slice contract and task pack

- **Status:** ACCEPTED by Human, 2026-09-14.
- **Date:** proposed and accepted 2026-09-14
- **Context:** The approved roadmap calls Phase 23 "Review and Human judgment".
  Its exact contracts and task pack were explicitly left to this planning step
  by the Phase 22 wrap. Roadmap section 5 requires the JSON, event, and schema
  contracts to be accepted before any product code in the slice changes.
- **Decision:** Accepted, together and as one unit,
  `contracts.md`, `public-interface.md`, `schema.md`,
  `execution-scope.md`, `v1-convergence-inventory.md`, `contract-coverage.md`,
  `qa/validation-contract.md`, and the ordered 18-task pack; declare the phase
  Heavy; and accept the six canonical spec amendments listed in
  `contracts.md` section 15. The Human planning decisions below are already
  folded into those artifacts.
- **Consequences:** Implementation may begin at task 1 after the QA plan and
  the Runner's QA readiness check. `reviews.coverage.v1` and `human.outcomes.v1`
  become advertisable at task 17. The named V1 judgment surfaces are removed at
  task 15 under replacement-bound removal. `work.claims.v1` stays unavailable
  until Phase 24 and `telemetry.trajectory.v1` until Phase 25.
- **Preconditions still open:** the Heavy-mode QA plan and fresh command and
  feasibility validation do not exist yet. The Human is activating QA Lead
  directly; Tech Lead performed no handoff. No task may be dispatched until the
  QA plan exists and the Phase Runner's readiness check passes.
- **Refs:** `10_IMPLEMENTATION_ROADMAP.md` sections 4, 5, 5.1, 7; PRD sections
  7 and 8; TRD sections 4, 5.4, and `03_TRD_SCHEMA.md` sections 3.7 and 7;
  Phase 22 DEC-187 kit dispositions C-03, C-04, C-10; DEC-220.

### DEC-222 Human planning decisions, 2026-09-14

Both open planning questions were resolved before acceptance. They amend
DEC-222 in place and are folded into the accepted artifacts; they are not
separate DECs.

- **Concern lifecycle.** `fid concern dispose` is accepted as a new accepted CLI
  leaf with a matching `concern.disposed` Event family, closing the gap between
  PRD-REV-007 (which requires status and disposition) and the CLI and Event
  catalogs (which omitted the mechanism). The disposition set is exactly
  `carried-forward`, `rejected`, and `superseded`. The Tech Lead's original
  `accepted` and `converted-to-finding` values are withdrawn: `accepted` was
  vacuous, and `converted-to-finding` was structurally wrong because Findings
  are created only inside a Review Record (PRD-REV-002) and describe rework on
  the current revision rather than carried-forward work. `carried-forward` is
  the only route from a Concern to future work because ADR-0004/0005 makes
  Issues the only durable work graph.
- **Concern disposal authority.** Self-disposal is permitted: the Actor that
  raised a Concern may dispose it. Concerns are deliberately lightweight and off
  the critical path, and a Builder later establishing that its own observation
  was not a problem is the ordinary case. The raising and disposing Actors are
  stored and surfaced separately so a self-dismissal is legible rather than
  prohibited. No independence rule and no Human Outcome requirement attaches to
  disposal; where a rejection is an authorized risk acceptance, the Human
  Outcome is recorded separately and referenced. This also avoids reordering the
  task pack, since a Human Outcome requirement would have pushed Concern
  disposal from task 8 to after task 10.
- **Requirement traceability.** `links` is restricted to Issue endpoints so
  `requirement_issue_links` is the single authoritative model, as required by
  the slice acceptance rule in `10_IMPLEMENTATION_ROADMAP.md` section 7. The new
  `verifies` relation is retained, and `03_TRD_SCHEMA.md` section 3.2 is added
  to the accepted spec-amendment list to record both changes.

### DEC-223: Link-preserving `evidence_records` rebuild in migration 0017

- **Status:** **ACCEPTED by Human, 2026-09-19.** Proposed by Tech Lead
  2026-09-15; amends DDL accepted under DEC-222, now approved.
- **Context:** QA Lead R-01 measured that accepted 0017 fails on any schema-16
  workspace whose receipt Evidence has Artifact links:
  `Runtime error near line 265: FOREIGN KEY constraint failed (19)`, rc=1. The
  runner opens `foreign_keys(1)` (`internal/db/db.go`) and wraps one migration
  file in one transaction, so `DROP TABLE evidence_records` performs an implicit
  `DELETE` that violates the immediate `evidence_artifact_links.evidence_id`
  reference (`0013:293`). `PRAGMA defer_foreign_keys = ON` only moves the failure
  to `COMMIT`, because a table rename never decrements the deferred counter.
- **Decision:** Amend `schema.md` sections 2 and 3 with the measured
  link-preserving order - stage `evidence_artifact_links` in a TEMP table, delete
  the rows, rebuild and rename `evidence_records`, restore the rows, drop the
  staging table - in both directions. The down direction restores only links
  whose Evidence survives the receipt-only filter. `evidence_artifact_links`
  keeps its original DDL and is never itself rebuilt. VAL-005 is strengthened
  rather than narrowed: its seed must carry Artifact links, must migrate in
  place, and must show an empty `foreign_key_check`.
- **Rejected alternative:** narrowing VAL-005 to a link-free seed. That would
  make the assertion pass while shipping a migration that fails in ordinary use.
  `catchUpWorkspaceLocked` runs `MigrateUp` against the existing populated
  projection, so this is a real upgrade path, not a synthetic test condition.
- **Consequences:** task 1 ships the staging SQL as accepted DDL and may not
  omit it as an optimization. DEC-178 still binds: shipped SQL equals the
  amended `schema.md` exactly.
- **Independent Tech Lead validation, 2026-09-15:** the failure and the fix were
  reproduced in a scratch sandbox with `sqlite3` 3.50.6 on a minimal schema
  holding the four relevant tables. The accepted text fails at the `DROP` with
  `Runtime error near line 17: FOREIGN KEY constraint failed (19)`, rc=1. The
  amended text returns rc=0 in both directions, with `foreign_key_check` empty,
  `integrity_check` `ok`, the Evidence row identical by named column with the
  three new columns null, both link rows restored with their exact
  `required_material`, `evidence_artifact_links` DDL byte-identical, and no
  leaked staging table. The down direction drops exactly the link whose `review`
  Evidence the receipt-only filter discards. This is a sandbox proof of the SQL
  mechanism with the CLI, not of the shipped migration under
  `modernc.org/sqlite`; task 1 still proves it on the real chain, and VAL-005 is
  where that proof is recorded.
- **Refs:** `qa-plan.md` R-01 and sections 2.4, 6.4; `schema.md` sections 2, 3.

### DEC-224: Named raising Events for the `human-outcome-required` Stop

- **Status:** **ACCEPTED by Human, 2026-09-19** ("ok as described"), after the
  Human asked for the four raising moments and the exact added permissions;
  both are enumerated below and in the `contracts.md` section 10 producer table.
  Proposed by Tech Lead 2026-09-15; adds effect rows to an accepted Event
  contract.
- **Context:** QA Lead R-02 established that no accepted effect row can raise
  the Stop that `contracts.md` section 10 requires. The clause states a state
  predicate ("the Human Obligation is the only unsatisfied one") while Stops are
  durable rows that only a committing Event can write. DEC-195 governs: add the
  missing rows, never weaken the assertion.
- **Decision:** Enumerate the complete set of Events whose commit can make the
  predicate true, and give each the effect rows it needs:
  `envelope.activated` (rows already added by task 9), `evidence.linked` (rows
  added by task 10), `review.recorded` (rows already optional), and
  `artifact.availability-observed` (rows added by task 10). At most one
  unresolved Stop exists per Obligation; a repeat qualifying Event raises no
  duplicate; the Stop is never withdrawn by later staling and resolves only
  through a qualifying Human Outcome. VAL-046 asserts all four sequences, the
  at-most-once rule, and both non-qualifying states.
- **Why four and not one:** `evidence.linked` alone covers the ordinary case,
  but a later independent Review that resolves the last blocking Finding
  (`review.recorded`) and an observation that flips the last required Artifact
  (`artifact.availability-observed`) both reach the same readiness state, and an
  Envelope whose only Obligation is the Human one reaches it at activation.
  Naming three of the four and silently dropping the fourth would leave
  PRD-HUM-002 unserved on a reachable path.
- **Consequences:** task 10 implements one shared readiness evaluation invoked
  from four producers and owns the two new effect-row additions.
- **Refs:** `qa-plan.md` R-02 and section 2.4; `contracts.md` section 10;
  `schema.md` section 6.2; DEC-195.

### DEC-225: Review independence conflicts with the mutating side only

- **Status:** **CONFIRMED by Human, 2026-09-19** ("your interpretation is
  correct"). Proposed by Tech Lead 2026-09-15 as determined by prior accepted
  material, and recorded for Human confirmation because it fixes product
  semantics.
- **Context:** QA Lead R-04: `contracts.md` section 3 said the recording Actor
  "must not equal any Actor that held a mutating Claim or a Run on the reviewed
  Issue", without saying whether "a Run" includes a Run under a `review` Claim.
  The two readings produce opposite fixtures and opposite product behavior.
- **Decision:** The conflicting set is mutating Claims and Runs under such
  Claims. A `review` Claim and its Runs never conflict, so one independent
  Reviewer may record successive Reviews across rework cycles and may resolve a
  Finding created by its own earlier Review Record; only a Review resolving a
  Finding **it** created is prohibited.
- **Why prior accepted material determines this:** PRD-REV-004 and the Glossary
  define independence against Builder self-resolution, not against repeat
  review; `contracts.md` section 5 scopes the self-resolution ban to the Review
  Record rather than the Actor, which is vacuous under the broad reading; and
  section 6's review-cycle budget, rework, and exhaustion model presumes
  repeated review, which the broad reading would make unusable without a fresh
  Actor per cycle.
- **Consequences:** VAL-015 gains the repeat-Review case and VAL-022 gains
  resolution by the same independent Reviewer's later Review. Independence is
  also re-evaluated at link time (VAL-077), so a Reviewer who later takes a
  mutating Claim cannot have that Review linked as Evidence.
- **Refs:** `qa-plan.md` R-04 and section 2.4; PRD-REV-004; GLOSSARY "Review
  Finding"; `contracts.md` sections 3, 5, 6, 13.

### DEC-226: Failure-code producers, precedence, and non-committing denials

- **Status:** **ACCEPTED by Human, 2026-09-19**, together with the part 5
  vocabulary fix added during that review. Proposed by Tech Lead 2026-09-15;
  fixes accepted public-interface semantics and removes one failure code.
- **Root cause found during Human review, 2026-09-19:** the Human could not
  reconcile the two code lists, and the reason was a real naming collision in
  the project's own language, not a gap in their understanding. `GLOSSARY.md`
  carried a single entry, `Reason code`, defined as "why an operation failed" -
  the failure set - while the slice contracts used `ReasonCode` for the
  evaluation-reason set. The implementation carries the same collision: the
  failure-code registry is named `reasonCodeDefinitions`
  (`internal/cli/codes.go`). One term denoted two closed sets, which is exactly
  how `finding-unresolved` came to be filed as a failure code. Parts 1 and 3
  below are both downstream of that collision.
- **Context:** QA Lead R-03, R-05, and R-06 found three contradictions: one
  input with two accepted codes; three new codes with no named producer and no
  precedence; and a committed-denial replay rule with no denial Event.
- **Decision:**
  1. `evidence.link` returns exactly one code, in the order
     `outcome-not-permitted` (any Human Outcome subject), `review-changes-required`,
     `review-not-independent`, `evidence-type-not-accepted`. The Human Outcome
     case is the most specific and therefore wins over the generic type
     mismatch.
  2. `public-interface.md` section 5.1 names the producing operation for every
     added code, and fixes the guarded-action denial order: retained Phase 21/22
     failures, then `protected-decision-required` (authority outranks
     precondition, lowest ordinal first), then `review-rework-exhausted`
     (disjoint action set), then `human-outcome-required`, then
     `obligation-unsatisfied`, which is returned only when no open Stop blocks
     the action.
  3. `finding-unresolved` becomes a ReasonCode only. No leaf may reject on it:
     linking a `pass` Review that leaves an earlier blocking Finding unresolved
     is accepted and the Obligation stays unsatisfied (VAL-026). Task 6 adds it
     to the evaluation reason precedence instead of `internal/cli/codes.go`.
  4. Phase 23 adds no denial Event. Stops are raised by the successful Events
     that cause them; the requests those Stops block commit nothing and return
     null `EventRef`, `replayed:false`, and the immutable Stop and Attention IDs.
     Determinism comes from the durable Stop, not a stored snapshot. The Phase 22
     `run.execution-blocked` committed-denial rule is unchanged.
  5. **Name the two sets** (added 2026-09-19). `GLOSSARY.md` replaces the single
     `Reason code` entry with **Failure code** (names one refused operation;
     exactly one per failed response; carries category, retryability, and exit
     status) and **Evaluation reason** (describes why one Obligation is
     unsatisfied; zero or more; no category, retryability, or exit status), plus
     the rule that an identifier belongs to both sets only where the same fact
     both refuses an action and describes a state. The Tech Lead applied that
     edit under this decision. Task 13 renames `reasonCodeDefinitions` to
     `failureCodeDefinitions` and its sibling test identifiers: a pure rename of
     unexported names, no behavior, JSON key, or public value change.
- **Why the vocabulary fix belongs here:** without it Phase 24 inherits the same
  collision, and the next misfiled identifier is as invisible as this one was.
  Per `AGENTS.md`, `GLOSSARY.md` is the canonical language source and unclear
  terms are resolved or flagged there rather than worked around, which is the
  same move DEC-195 makes for a missing schema row.
- **Rejected alternative:** adding denial Events for `claim.acquire`,
  `run.start`, `issue.transition-terminal-success`, and `envelope.complete`.
  That would expand the Event surface and write journal entries for refused
  requests, to reproduce a guarantee the existing Stop already provides.
- **Consequences:** VAL-025, VAL-039, VAL-056, and VAL-057 are restated
  accordingly; task 6's code list loses `finding-unresolved`; `GLOSSARY.md` is
  amended now and recorded `touched` in the spec sync table; task 13 carries the
  rename. No new assertion is added for part 5: vocabulary and an unexported
  rename change no observable behavior, and the existing gates compile the
  rename.
- **Refs:** `qa-plan.md` R-03, R-05, R-06; `public-interface.md` section 5;
  `contracts.md` sections 6, 10, 13; `schema.md` section 6.3; DEC-216.

### DEC-227: Explicit test-time bound, wall-time budgets, and speed ownership

- **Status:** **ACCEPTED by Human, 2026-09-19.** Proposed by Tech Lead
  2026-09-15; changes a required command, now approved.
- **Context:** QA Lead R-13 measured the tagged integration suite at 560.2 s
  against Go's implicit 600 s per-binary default, with `cmd/fid` at 472.0 s and
  `internal/workflow` at 458.4 s, and Phase 23 adds real-process journeys to
  exactly those packages at six tasks. A breach produces a killed command, which
  is `BLOCKED` and not a repairable failure inside a task.
- **Decision:** Run the two long gates with an explicit `-timeout 30m` in every
  Phase 23 task spec; record the `cmd/fid/integration`, `cmd/fid`, and
  `internal/workflow` package times in every completion report; hold budgets of
  900 s for the tagged suite and 600 s for any ordinary package, asserted at
  phase end by VAL-079; and assign speed repair to the task whose own new
  journeys caused a breach, with `t.Parallel()` permitted in **new** Phase 23
  journey files only. A pre-existing cause is escalated, not absorbed. Task 1
  additionally prunes stale `fid-integration-binary-*` directories in
  `cmd/fid/integration/main_test.go`.
- **Spec sync if approved:** `AGENTS.md` (Commands) still lists the two gates
  without a bound. Task-spec command lists take precedence, so Phase 23 needs no
  edit to dispatch; the wrap-up spec sync records whether the same explicit bound
  becomes the project default. Noted in the spec sync verdict table below.
- **Why not DEC-221's answer:** DEC-221 chose performance repair over raising a
  hosted CI job timeout, where the timeout was masking a real parallelism
  problem. Here the limit is an undeclared language default, and the measured
  budget replaces it as the real control; the explicit bound weakens no
  assertion, changes no scope, and removes no gate.
- **Refs:** `qa-plan.md` R-13, I-16; `execution-scope.md` section 7; DEC-221.

### DEC-228: Tech Lead dispositions of the remaining QA first-pass revisions

- **Status:** ACCEPTED by Tech Lead 2026-09-15 under planning authority. No
  protected definition, command, threshold, or accepted contract semantics
  change; these are planning-artifact corrections.
- **Decision:** Apply R-07 through R-12 and improvements I-01 through I-16 as
  recorded in the disposition table below, including the three DEC-167/173 file
  scope corrections (R-08, R-09), the R-10 Human fixture owner, the R-11 split of
  VAL-002, and the three new assertion rows from R-12. Also record the
  `contracts.md` section 8 rule that a Concern or Deviation subject need not
  resolve to an existing row (I-14) and the per-type evaluation reason codes
  (I-15).
- **Consequences:** all 79 validation rows move from `Draft` to `Ready`; rows
  depending on DEC-223..227 revert to `Draft` if those are not approved.
- **Refs:** `qa-plan.md` section 2.3; DEC-167/173, DEC-192, DEC-194, DEC-195.

### DEC-229: Tech Lead dispositions of the QA readiness-refresh revisions

- **Status:** ACCEPTED by Tech Lead 2026-09-19 under planning authority. No
  protected definition, command, threshold, budget, or accepted assertion
  semantics changes; these are file-scope and recorded-cause corrections of the
  same class as DEC-228.
- **Context:** The readiness refresh (`qa-plan.md` section 12) re-proved all
  five Human-approved decisions against fresh evidence and returned two
  Must-fix revisions, R-14 and R-15, plus improvements I-17..I-26.
- **Decision - R-14, four closed registries with no owning task.** Accepted in
  substance, with one owner changed and one sub-item rejected:
  - `internal/workflow/stop.go` and
    `internal/journal/current_event_value_catalog.go` join tasks 6, 9, 10, and
    11, each for its own additions. Verified: `stopCodes` (7), `conditionTypes`
    (6), the 4-field `Condition`, and `validStopDetails` are closed maps, and
    `StopCode` validation keys on `currentStopDetailShapes`, so without these
    two files every Phase 23 Stop is rejected at raise time. A Condition type
    is registered by the task that first **raises** a Stop carrying it - tasks
    9 and 10 - not deferred to task 11, which registers the remaining two and
    owns the resolution behavior of all four. `public-interface.md` section 2
    already fixes every added StopCode, StopDetails variant, and Condition
    field, so no contract addition was needed.
  - `internal/journal/current_event_columns.go` joins **task 1**, not the five
    domain tasks QA proposed. Task 1's own binding constraints already register
    all nine Event types with their complete `schema.md` 6.2 effect rules; the
    ten effect tables those rules name belong with them, the 0017 DDL that
    fixes their column domains ships in the same task, and
    `TestCurrentEventColumnDomainsMatchAcceptedDDL` then proves the agreement
    there instead of five tasks editing one shared map. Task 15 renames the
    `judgment_concerns` key with its table.
  - `internal/export/archive_semantics_catalog.go` gains explicit owners under
    a standing rule recorded in `execution-scope.md` section 6: a task that
    adds, amends, or removes a journal Event rule mirrors that same change in
    `semanticEventRules` in the same task. Task 1 mirrors the nine new rules,
    task 9 its two amendments, task 10 its two, task 12 the deletion (R-09),
    task 15 the rename.
  - **Rejected:** adding `internal/journal/current_event_catalog.go` to tasks 4
    and 5. QA inferred the gap from the file lists alone. Task 1 registers all
    nine Event types, their operation mappings, and their complete effect rules
    including `review.recorded` and `finding.addressed`, so tasks 4 and 5 own
    only their fence key in `project_effects.go`. Adding the file would create
    a shared-map conflict for no behavior.
- **Decision - R-15, CF-042's recorded cause is wrong.** Accepted, with option
  (i), the causal fix. The leak is per-run, not per-kill: five tests re-execute
  the test binary, each child rebuilds a ~16 MB binary, and the crash-process
  children never reach `RemoveAll`; one clean run leaks 20 directories and
  306 MB, so a 24-hour age bound prunes nothing intra-phase. Task 1 therefore
  publishes the built binary path and has a re-executing child reuse it,
  building and removing nothing. All five call sites already spawn with
  `append(os.Environ(), ...)`, so the change stays inside
  `cmd/fid/integration/main_test.go`. Option (ii) - pruning everything but the
  current process's directory - is rejected: it would delete a concurrent run's
  live binary. The 24-hour prune is retained only as the residual safety net
  for a parent killed mid-run, and task 1 records the directory count and
  `/tmp` free space either side of its tagged gate. Reuse also removes about
  twenty redundant `go build` invocations per run, which helps the DEC-227
  budgets.
- **Decision - improvements.** I-17 through I-26 are all applied; see the table
  below. I-25 additionally surfaced an unreachable failure code, which is
  DEC-230.
- **Consequences:** eight task specs, `execution-scope.md`, `contracts.md`,
  `public-interface.md`, `schema.md`, `v1-convergence-inventory.md`, the
  validation contract, and both carry-forward registers are edited. The
  contract stays at 79 rows; no row is added, removed, renumbered, or weakened.
  A fresh QA Lead readiness confirmation is required before the Runner's
  readiness check.
- **Refs:** `qa-plan.md` sections 12.4, 12.7, 12.10; DEC-167/173, DEC-224,
  DEC-226, DEC-227, DEC-228; CF-042.

### DEC-230: Remove `decision-not-found` from the added failure code set

- **Status:** ACCEPTED by Tech Lead 2026-09-19 under planning authority,
  applying DEC-226's own acceptance criterion. Reversible in one line if the
  Human prefers the alternative below; flagged to the Human with this batch.
- **Context:** Pinning each `*-not-found` producer under I-25 showed that six
  of the seven codes have a real producing leaf, and `decision-not-found` has
  none. No Phase 23 leaf takes a Decision Reference ID as input: section 4 gives
  Decision References `record-reference` and `list` only, with no `get`, and no
  mutation names one. `linked_to_decision` is a transition rule that reports
  satisfied or unsatisfied, not a not-found error. VAL-056 requires one
  reproduction per added code through its named producer, so the row was
  unsatisfiable exactly for this code.
- **Decision:** Remove `decision-not-found` from `public-interface.md` section 5
  and from task 9's `internal/cli/codes.go` list. The added set becomes 18
  codes, not 19. Reading a Decision Reference stays `decision.list --class`.
- **Alternative considered and rejected:** add a `decision.get` leaf. It would
  make the code reachable, but it adds public surface to fix dead vocabulary,
  and four other Phase 23 entities already ship with `list` and no `get`, their
  not-found codes reachable only because other leaves take their IDs. Removing
  an unreturnable code is the smaller and more consistent change; DEC-226 exists
  precisely to stop codes with no producer.
- **Consequences:** VAL-056 is satisfiable in full at 18 codes. QA's 12.6 count
  of 19 new codes becomes 18. If a `decision.get` leaf is ever added, the code
  returns with it.
- **Refs:** `public-interface.md` sections 4, 5, 5.1; `qa-plan.md` 12.6, I-25;
  DEC-226; VAL-056.

### DEC-231: Typed-result registries are owned by the task that unfences the family

- **Status:** ACCEPTED by Tech Lead 2026-09-19 under planning authority. No
  protected definition, command, threshold, budget, public surface, or accepted
  assertion semantics changes; this is a file-scope and ownership correction of
  the same class as DEC-228 and DEC-229.
- **Context:** The QA Lead's DEC-229/DEC-230 confirmation pass (`qa-plan.md`
  section 13) confirmed every edit it was asked to confirm and surfaced one new
  Must-fix, **R-16**, while doing so. Registering the nine Phase 23 families in
  `currentEventRules` - which task 1 must do - routes them onto the journal's
  strict validation path, and three further closed registries on that path were
  named by no task spec. Independently re-verified in source before
  disposition: `project_effects.go:20` calls `validateCurrentAppliedResults`,
  which gates on `currentEventRules` membership
  (`current_event_result_projection.go:12`) and then requires the family to
  clear (a) the inline `fields` grammar map in `current_event_json.go:135-141`,
  which otherwise returns `unregistered current result family`; (b)
  `currentExpectedResult`'s switch (`current_event_result_projection.go:69-121`),
  whose `default` returns `unknown current result`; and (c) `currentValueShapes`
  in `current_event_value_catalog.go`, from which `currentProjectedValue`
  (`:175-203`) reads the projected field list. `schema.md` 6.3 makes a typed
  snapshot mandatory for every Phase 23 mutation, so all nine families need all
  three. Task 1 stays green because its own fence rejects all nine; the failure
  first lands at task 4 with its cause outside task 4's scope.
- **Decision:** Accept R-16. **A family's typed-result grammar, expected-result
  branch, value shapes, and any composite fill helper are added by the task
  that deletes that family's fence key**, not by task 1. Recorded as a standing
  rule in `execution-scope.md` section 6. The nine keys map one-to-one onto six
  tasks: task 4 `review.recorded`; task 5 `finding.addressed`; task 8
  `concern.recorded`, `concern.disposed`, `deviation.recorded`; task 9
  `decision.referenced`; task 10 `human-outcome.recorded`; task 12
  `requirement.referenced`, `requirement.linked`.
- **Why the unfencing task and not task 1, where DEC-229 put the column
  catalog.** The two cases look identical and are not, on the one test that
  matters - whether the owning task can execute what it registers:
  - *DEC-229's column catalog is provable at task 1.*
    `TestCurrentEventColumnDomainsMatchAcceptedDDL` iterates `currentEventTables`
    and introspects the migrated database with `PRAGMA table_info`. It needs no
    Event to be appended, it runs against the 0017 DDL task 1 itself ships, and
    it fails closed. Column domains are also fully derivable from that DDL.
  - *These three registries are not provable at task 1.*
    `validateCurrentResultJSON`, `currentExpectedResult`, and
    `currentProjectedValue` execute only when an Event of the family is
    appended, and `unavailableJudgmentEvents` makes that impossible for all nine
    until their unfencing task. Putting all nine at task 1 would write roughly
    nine result grammars, nine switch branches, fourteen value shapes, and three
    composite fill helpers behind a fence guaranteeing that **none of them is
    ever executed in the task that wrote them**. Ownership is only assignable
    where it is executable; this would convert a scoped, provable obligation
    into an unverifiable one and hand task 4 a silent inherited defect instead
    of its own failing test.
  - *Derivability differs too.* `currentProjectedValue` maps a projected public
    row plus child-table queries onto a field list. That mapping is the domain
    the unfencing task builds; at task 1 only the DDL exists.
  - *It is the consistent reading of DEC-229, not a departure from it.*
    DEC-229 already ruled that a Condition type is registered by the task that
    first **raises** a Stop carrying it. Same principle: registration lands
    where it first executes. Task-1-for-everything would be the inconsistent
    choice.
  - *The shared-map objection does not apply.* The eighteen tasks are strictly
    sequential, so there is no concurrent edit, and no test iterates
    `currentEventRules` to cross-check the other three registries - verified by
    search across `internal/journal` and `internal/export` test files - so a
    partial registration cannot be observed before the fence opens.
- **Value-shape ownership**, each first registered by the task named and reused
  unchanged afterwards: task 4 `Review`, `Coverage`, `CoverageItem`,
  `CoverageExclusion`, `Finding`, `FindingDisposition`, plus the
  `currentProjectedValue` fill helpers for `Review` (`findings[]`, `coverage`)
  and `Finding` (`dispositions[]`); task 5 reuses `Finding` and
  `FindingDisposition` and registers no new shape; task 8 `Concern`,
  `ConcernDisposition`, `Deviation`, plus the fill helper for
  `Concern.disposition`; task 9 `DecisionRef` and `ProtectedDecision`; task 10
  `HumanOutcome`, reusing `ProtectedDecision`; task 12 `RequirementRef` and
  `RequirementLink`. Task 4 also widens the journal `Subject` value shape - today
  `type:=run,envelope,activation` - to the complete accepted SubjectType set in
  `public-interface.md` section 2, in one monotonic edit, as the first task to
  register a value shape that references `Subject`; no later task re-narrows it.
- **Two consequential sub-items, both accepted.** `currentValueShapes["Claim"]`
  closes `mode:=mutation`, which `public-interface.md` section 2 widens to
  `mutation|review`: task 3, which is where review-mode Claims become
  reachable, and it is genuinely reachable because `run.abandoned` and
  `run.recovered` both project a `Claim`. `currentValueShapes["Budget"]` and
  `["ImmutableBudget"]` close `unit` without `review-cycles` and
  `exhaustion_action` without `review-rework-exhausted`: task 6, which adds
  both.
- **Recorded for the Developer, not a new finding:** `idempotency_results` is
  encoded `currentStage(0, ...)` - optional - for nearly every existing family,
  yet `validateCurrentAppliedResults` requires the result row unconditionally
  once a family is in `currentEventRules`. The Phase 23 rows in `schema.md` 6.2
  mirror that existing encoding exactly. "Optional" there is the stage
  encoding, not permission to omit the typed snapshot, which DEC-216 and
  `schema.md` 6.3 make mandatory for every Phase 23 mutation.
- **Improvements:** I-27, I-28, and I-29 are all applied.
- **Consequences:** eight task specs (3, 4, 5, 6, 8, 9, 10, 11, 12),
  `execution-scope.md` section 6, `schema.md` 6.2, and the validation contract
  header are edited. Task 1 is unchanged and gains nothing. The contract stays
  at **79** rows; no row is added, removed, renumbered, reworded, or weakened,
  and no `VAL-NNN` covers or needs to cover these registries directly - each is
  proved by the first successful Event of its family, which every unfencing
  task already asserts. A fresh QA Lead confirmation pass is required before the
  Runner's readiness check.
- **Refs:** `qa-plan.md` section 13.5 (R-16), 13.7; `execution-scope.md` 6;
  `schema.md` 6.2, 6.3; `public-interface.md` 2, 3, 4; DEC-192, DEC-216,
  DEC-228, DEC-229.

### DEC-232: Close the journal strict-path registry set and keep `Subject` narrow

- **Status:** ACCEPTED by Tech Lead 2026-09-19 under planning authority. No
  protected definition, command, threshold, budget, public surface, or accepted
  assertion semantics changes; file-scope, ownership, and recorded-cause
  corrections of the same class as DEC-228, DEC-229, and DEC-231. Two items
  edit `public-interface.md` and are flagged to the Human with this batch as
  reversible: naming the section 2 row `ErrorSubjectType`, and naming the
  already-accepted judgment subject domain `JudgmentSubject`. Neither adds,
  removes, or changes a single accepted value.
- **Context:** the QA Lead's DEC-231 confirmation pass (`qa-plan.md` section 14)
  confirmed DEC-231 - the unfencing task is the right owner, the executability
  test is the right test, and all four reachability paths are clean - and
  returned two new Must-fix revisions, R-17 and R-18, plus I-30 and I-31. R-17
  asked me to derive the affected task set myself rather than adopt QA's, and
  I-31 asked for the rule to be stated once with an inventory. Doing both
  surfaced two further defects of the same family, recorded below as the Tech
  Lead's own findings.

**R-17, `validateCurrentCardinality` has no owning task. Accepted; owner is
task 4 alone.** `current_event_cardinality.go:8`, called from the strict branch
at `current_event_validation.go:96`, defaults `maximum := 1` per `table/op` key
and widens it only through a closed `switch effect.Table` at `:14` naming
Phase 21/22 tables. Independently reproduced: a `review.recorded` Event with two
`review_findings/insert` Effects returns `current Event review.recorded: too
many review_findings/insert Effects`, and the same for two
`finding_dispositions/insert`, while one Finding, two `stops/insert`, two
`stop_resolutions/insert`, two `attention_items/update`, and
`requirement_refs` insert-plus-update all pass. `public-interface.md` section 4
accepts `findings:FindingInput[0..64]` and `resolves_finding_ids:ID[0..64]`, so
the ordinary two-Finding Review was blocked at task 4 with its cause outside
task 4's scope.

QA's proposed task set is confirmed, by a different route than QA's. Only
`review.recorded` repeats a key: every other new family writes one row per key,
`requirement.referenced`'s insert and update are distinct keys, and the Stop
bundle and `budget_samples` tables are already widened by the existing branches.
The four **amended** families needed separate treatment and also need nothing:
`envelope.created`, `envelope.activated`, `evidence.linked`, and
`artifact.availability-observed` are all in `retainedCurrentRules`
(`current_event_catalog.go:82-83,87-88`), and the retained branch returns at
`current_event_validation.go:64` before cardinality is ever reached. That is
what makes the 16-entry `protected_decisions` input safe without an entry. The
widening is conditioned on `review.recorded` so `finding.addressed` keeps its
implicit maximum of one.

**R-17b, `validateCurrentUpdate` has no owning task. Tech Lead finding, same
class, accepted.** Found while deriving R-17's task set by walking every
validator the strict path calls. `current_event_updates.go:14` is a closed
`switch effect.Table` whose `default` returns `current Event update table not
admitted`. It lists no Phase 23 table, and two Phase 23 families emit an update:
`concern.disposed`'s **required** `judgment_concerns/update` and
`requirement.referenced`'s optional `requirement_refs/update`. Reproduced: both
return that error; `attention_items/update` is admitted, as expected. The file
is named in **zero** phase documents. Owners: task 8 and task 12, each supplying
the exact changed-field and guard sets from the accepted DDL; task 15 renames
the `judgment_concerns` entry with its table.
`validateCurrentTerminalValues` in the same file fails open and is recorded as
optional strengthening, not a requirement.

**R-17c, the `validateCurrentJSON` column switch. Tech Lead finding, described
rather than reassigned.** `current_event_json.go:14` is a second closed registry
in the file DEC-231 already assigned, keyed `table.column`, returning
`unregistered current JSON column`. Reproduced for all **twelve** Phase 23
`*_json` columns on unfenced families - task 4 five, task 8 three, task 9 two,
task 10 one, task 12 one; the count was written as eleven here and in the
confirmation brief and is corrected under I-32 (DEC-233), with the per-task
assignment unchanged and already complete. This is **not** a new ownership
gap - the file is already in the scope of tasks 4, 5, 8, 9, 10, and 12, and the
failure is loud and inside the owning task - but their spec text named only the
`fields` map, so each task's bullet now names its own columns and their
accepted kinds.
`evidence_records.target_repository_json` needs nothing because `evidence.linked`
is retained.

**R-18, the `Subject` widening removes the wrong conditional. Accepted, and the
defect is larger than reported.** DEC-231 directed task 4 to widen
`currentValueShapes["Subject"]` "to the complete accepted SubjectType set in
`public-interface.md` section 2". Withdrawn on three counts:

1. *The recorded rationale is false.* `Subject` is referenced by exactly three
   shapes - `Condition` (`:15`), `Stop` (`:26`), `Attention` (`:27`) - and they
   serve `stop.raised`, `stop.resolved`, `attention.resolved`,
   `run.budget-observed`, and, through `stops.resolution_condition_json` in
   `validateCurrentJSON`, every `stops/insert` in any family. All are live from
   task 1 and none is fenced, so "the fence keeps every other Phase 23 family
   closed" does not hold and the widening is not monotonic-behind-the-fence.
   QA is right that it is not reachable as a defect today, because
   `stop.go:377 validSubject` still restricts raising; it is a silent loss of an
   accepted invariant, contradicting `../phase-22/public-interface.md:67-68`.
2. *No single shape can serve both domains.* `stops` and `attention_items`
   CHECK `('run','envelope','activation')` (`0015:56,82`), unchanged by 0017.
   The four judgment tables CHECK
   `('issue','envelope','run','review','finding','requirement')`
   (`schema.md` section 2). Neither set contains the other - Stop has
   `activation`, judgment has `issue`, `review`, `finding`, `requirement` - so
   "widening" is impossible in principle, not merely unwise. A union would let
   the journal admit a Stop subject its own DDL rejects.
3. *The cited source set is the wrong vocabulary entirely.* Section 2's
   `SubjectType` row is the CLI **error** subject vocabulary - `type SubjectType`
   at `internal/cli/codes.go:88` - which is why `requirement` can be "retained"
   there and why the six added values are exactly the six new entity kinds that
   appear as error subjects. It has ~19 members including `path`, `label`,
   `export`, and `idempotency`. Followed literally, task 4 would have let a Stop
   resolution condition name `type:"path"`.

**Decision:** `currentValueShapes["Subject"]` is not edited by any Phase 23
task and must be byte-identical after task 4. Task 4 instead registers
`JudgmentSubject`, `type:=issue,envelope,run,review,finding,requirement id:ID`,
used by `Finding` at task 4 and reused unchanged by `Concern` and `Deviation`
at task 8 and `HumanOutcome` at task 10. This is a new name for a domain already
accepted under DEC-222 in `contracts.md` section 8 and four identical DDL
CHECKs, not a new domain. `public-interface.md` section 2's row is renamed
`ErrorSubjectType`, and section 3 gains a four-row table separating `Subject`,
`JudgmentSubject`, `ErrorSubjectType`, and `ReasonSubjectType`, because reusing
one bare name for two domains is what produced this defect and the note is what
stops it recurring.

**R-18b, `schema.md` 6.2 misclassifies `envelope.created`. Tech Lead finding,
accepted.** Found while establishing which amended families reach cardinality.
The I-17 sentence reads "The first row is in `currentEventRules`; the last three
families live in `retainedCurrentRules`". All **four** are retained;
`envelope.created` is at `current_event_catalog.go:82`. A Developer following it
would add an `envelope.created` key to `currentEventRules`, routing a live family
onto the strict path where its `work_envelopes` and `envelope_*` Effects fail on
`unknown table` - the exact outcome I-17 existed to prevent. Corrected in
`schema.md` 6.2 and in task 9's scope, with the reason the retained stages are
still required: task 1 registers `envelope_protected_decisions` in
`currentEventTables`, and the retained branch validates every Effect on a
`currentEventTables` table against the retained rule's stages.

- **I-30, I-31.** Both applied. I-30: the idempotency-stage note now reads 15 of
  22 optional and 7 required, counted at the planning base. I-31: the DEC-231
  rule is restated once as a general obligation - *a task that unfences a family
  owns that family's entry in every closed journal registry the strict path
  consults* - above a six-row inventory in `execution-scope.md` section 6 naming
  each registry, its file and line, its exact rejection message, and its Phase 23
  owners. Task 12 owns the end-of-phase nine-family, six-registry audit, and a
  registry not on the inventory is an escalation rather than a local fix. Two
  further path registries were checked and deliberately recorded as needing
  nothing: `validateCurrentBindings`'s `entityTable` map and
  `currentResultEntity`'s map both fail open to the aggregate ID, which is
  correct for all nine families.
- **Consequences:** seven task specs (4, 5, 8, 9, 10, 12, 15),
  `execution-scope.md` section 6, `public-interface.md` sections 1, 2, 3, 4,
  `schema.md` 6.2, and `contracts.md` section 8 are edited. Task 1 is again
  unchanged. The contract stays at **79** rows; no row is added, removed,
  renumbered, reworded, or weakened. Task 4 gains two acceptance criteria that
  cite existing assertions (VAL-017, VAL-020) rather than new ones. A fresh QA
  Lead confirmation pass is required before the Runner's readiness check.
- **Evidence:** all six reproductions were run by the Tech Lead against the
  planning base in `internal/journal`, using throwaway probe files created and
  deleted within the session; the working tree is unmodified outside
  `project-docs/`. `go test ./internal/journal/` is the harness; no product file
  was touched.
- **Refs:** `qa-plan.md` sections 14.5 and 14.7 (R-17, R-18, I-30, I-31);
  `execution-scope.md` 6; `public-interface.md` 1, 2, 3; `schema.md` 6.2;
  `contracts.md` 8; `../phase-22/public-interface.md` 67-68; DEC-192, DEC-216,
  DEC-228, DEC-229, DEC-231.

### QA confirmation-pass-2 disposition table (R-17, R-18, I-30, I-31)

| ID | Disposition | Where applied |
|---|---|---|
| R-17 | Accept; owner is task 4 alone, derived independently and reproduced | DEC-232; `execution-scope.md` 6; task 4 |
| R-17b | New, Tech Lead; accept - `validateCurrentUpdate` is unowned and fails hard | DEC-232; `execution-scope.md` 6; tasks 8, 12, 15 |
| R-17c | New, Tech Lead; describe rather than reassign - file already owned | DEC-232; tasks 4, 8, 9, 10, 12; task 5 records the checked negative |
| R-18 | Accept and strengthen; `Subject` untouched, `JudgmentSubject` added | DEC-232; `execution-scope.md` 6; `public-interface.md` 1, 2, 3, 4; `contracts.md` 8; tasks 4, 8, 10 |
| R-18b | New, Tech Lead; accept - all four amended families are retained | DEC-232; `schema.md` 6.2; task 9 |
| I-30 | Apply | `execution-scope.md` 6: 15 of 22 optional, 7 required |
| I-31 | Apply, and extend from four registries to six | `execution-scope.md` 6 inventory; task 12 audit |

### QA readiness-refresh disposition table (R-14, R-15, I-17..I-26)

| ID | Disposition | Where applied |
|---|---|---|
| R-14 | Accept (a), (b), (d); accept (c) with task 1 as the owner instead of five domain tasks; reject the tasks 4/5 sub-item | DEC-229; `execution-scope.md` 6; tasks 1, 6, 9, 10, 11, 15 |
| R-15 | Accept; correct the recorded cause and adopt option (i), the causal fix | DEC-229; task 1; `execution-scope.md` 7; CF-042 in both registers; `v1-convergence-inventory.md` |
| I-17 | Apply | `schema.md` 6.2 names `retainedCurrentRules` for the three amended families |
| I-18 | Apply | `contracts.md` 2: only the relative order is accepted; renumbering the unexported ranks is permitted |
| I-19 | Apply | `public-interface.md` 5.1 step 1 excludes `obligation-unsatisfied` |
| I-20 | Apply | task 13 renames nothing in `cmd/fid/evidence_test.go`; that `reasonCodes` is already the correct sense |
| I-21 | Apply | task 4 names `idempotency-conflict` as conformance registration |
| I-22 | Apply | VAL-006 probes Obligation counts 0, 1, 16, 17 |
| I-23 | Apply | VAL-034 proves the DEC-228 unresolvable-subject rule deliberately |
| I-24 | Apply | task 10 records the blob remove/restore fixture technique |
| I-25 | Apply, and it surfaced DEC-230 | `public-interface.md` 5.1 pins six producers; the seventh code is removed |
| I-26 | Apply | task 1 records the no-single-step-down fact and the dirty `schema_migrations` failure state |

### QA confirmation-pass disposition table (R-16, I-27..I-29)

| ID | Disposition | Where applied |
|---|---|---|
| R-16 | Accept; owner is the unfencing task, not task 1 | DEC-231; `execution-scope.md` 6; tasks 3, 4, 5, 6, 8, 9, 10, 12 |
| I-27 | Apply | `qa/validation-contract.md` header records the 2026-09-19 refinement and drops the satisfied prerequisite |
| I-28 | Apply | tasks 6 and 11: the vacuous `internal/journal/` bullets are replaced with their real DEC-229/DEC-231 scope |
| I-29 | Apply | `schema.md` 6.2 `envelope.activated` owner cell now reads "rows added by task 9; also raised from this family by task 10" |

### QA first-pass disposition table (R-01..R-13, I-01..I-16)

Every Must-fix is accepted; none is rejected, deferred, or cancelled.

| ID | Disposition | Where applied |
|---|---|---|
| R-01 | Accept; protected data-model amendment | DEC-223; `schema.md` 2, 3; VAL-005; task 1 |
| R-02 | Accept; name four raising Events and add effect rows | DEC-224; `contracts.md` 10; `schema.md` 6.2; VAL-046; task 10 |
| R-03 | Accept; `outcome-not-permitted` wins for any Human Outcome subject | DEC-226; `contracts.md` 13; VAL-025 |
| R-04 | Accept; narrow reading, determined by prior accepted material | DEC-225; `contracts.md` 3; VAL-015, VAL-022, VAL-077 |
| R-05 | Accept; producer table plus denial precedence; `finding-unresolved` becomes a ReasonCode | DEC-226; `public-interface.md` 5.1; VAL-039, VAL-056; tasks 6, 13 |
| R-06 | Accept; Phase 23 denials commit nothing; no denial Event added | DEC-226; `public-interface.md` 5.2; `schema.md` 6.3; VAL-057 |
| R-07 | Accept all four parts | VAL-008, VAL-009, VAL-045, VAL-053, VAL-061; tasks 2, 6, 10, 12, 14, 15 |
| R-08 | Accept; add `internal/workflow/activation.go` to five file scopes | tasks 4, 5, 8, 9, 10; `execution-scope.md` 6 |
| R-09 | Accept; add `internal/export/archive_semantics_catalog.go` to task 12 | task 12; `v1-convergence-inventory.md` 1 |
| R-10 | Accept; task 14 generates the Human fixtures, task 18 regenerates them | VAL-070..072; tasks 14, 18 |
| R-11 | Accept; split into VAL-002 and the new VAL-075 | validation contract |
| R-12 | Accept; add VAL-076, VAL-077, VAL-078 | validation contract; tasks 6, 11, 16 |
| R-13 | Accept; explicit bound, measured budgets, and named speed ownership | DEC-227; `execution-scope.md` 7; VAL-079; every task spec |
| I-01 | Apply | VAL-001 scoped to schema 17 |
| I-02 | Apply | VAL-004 note corrected; the phantom cross-reference removed |
| I-03 | Apply | task 4 coverage cites VAL-003 |
| I-04 | Apply | VAL-035 enumerates seven negatives |
| I-05 | Apply | VAL-033 asserts both supersession directions |
| I-06 | Apply | VAL-014, 022, 038, 040, 043 gain their negatives |
| I-07 | Apply | VAL-021 records that its non-vacuous proof is VAL-026 |
| I-08 | Apply | VAL-027 import leg re-verified at task 16 |
| I-09 | Apply | VAL-010 re-verified with present evidence at task 10 |
| I-10 | Apply | VAL-069 owned by task 18; inventory section 1 corrected |
| I-11 | Apply | VAL-006 covers protected-Decision authoring |
| I-12 | Apply | VAL-061 and VAL-062 cover the implicit drops and the deletes |
| I-13 | Apply | task 1 binding constraint on edit locality; VAL-075 |
| I-14 | Apply | `contracts.md` 8 subject-resolution rule; task 8 records it |
| I-15 | Apply | `contracts.md` 2 reason codes; VAL-055, VAL-056 |
| I-16 | Apply, split | task 1 prunes stale directories now; CF-042 carries the durable fix |

The QA Lead's "no revision requested" note on `idempotency-conflict` is adopted
as binding: VAL-019 and VAL-057 require a real-process reproduction through the
built binary, and fixture inventory is not evidence.

### DEC-233: Widen the registry frame to Phase 23 values in Phase 21/22 families

- **Status:** ACCEPTED by Tech Lead 2026-09-19 under planning authority, except
  the `Condition` item, which is **escalated to the Human as ESC-23-001**. No
  protected definition, command, threshold, budget, or accepted assertion
  semantics changes. One item edits `public-interface.md` section 2 and is
  flagged to the Human with this batch as reversible: naming the implemented
  `SubjectType` set authoritative for `ErrorSubjectType`, which adds and removes
  no name (I-34).
- **Context:** the QA Lead's DEC-232 confirmation pass (`qa-plan.md` section 16)
  confirmed DEC-232 in full - the `JudgmentSubject` split is better than QA's own
  R-18 proposal, R-17's owner derivation re-derives to task 4 alone, and the
  six-row inventory is complete on the strict path - and returned one new
  Must-fix, R-19, plus I-32, I-33, and I-34. R-19 is the fifth consecutive pass
  to find a registry inside the previous pass's fix, and the first to find one
  **off** the strict path. I-33 asked me to widen the frame rather than add a
  seventh row. Doing that surfaced two further items, recorded below as Tech Lead
  findings, and one of them is not dispositionable under planning authority.

**R-19, the journal's second copy of the capability allowlist. Confirmed
independently, accepted, owners are QA's five.**
`internal/journal/current_event_value_relations.go:63` holds `const actions`,
a closed 23-name string identical to the keys of
`sortedUniqueCapabilities` at `internal/workflow/activation.go:832`. Verified in
source: `activation_components.evidence_json` is validated as `[]ActorProfile`
when `component_type` is `authority-profile` (`current_event_json.go:52-63`), on
`activation.component-staged`, which sits in `currentEventRules`
(`current_event_catalog.go:28`) behind no fence; and the integration fixtures
that build Actor profiles pass capability lists straight through Kit Activation
(`p21_conformance_test.go:96`, `p22_fake_test.go:308`, and
`p22_resume_actions_test.go:111`).
So the first Activation staging a Phase 23-capable profile is rejected before
task 4 records a Review, with its cause in a file **no** phase document names -
confirmed by search; the only prior mention anywhere is
`phase-22/task-1-done.md`.
Owners are tasks 4, 5, 8, 9, and 10, matching R-08's set for the other copy, and
each now adds its own name to both copies in one change. The split is
`public-interface.md` 4.1's: the seven grant actions in,
`requirement.record-reference` and `requirement.link` out. 23 becomes 30 in
both places.

**R-19's follow-on, accepted: task 18 audits the two copies against each
other.** Nothing cross-checks them today, and each is now widened name-by-name
across five tasks, so a single missed name leaves a capability a profile can
hold and the journal rejects. VAL-067's audit asserts set equality between the
two directly rather than each against a literal. This is the protection DEC-229
gave `semanticEventRules`, applied to the pair R-19 exposed; it needs no new
assertion, because VAL-067 already owns the thirty-name set.

**I-33, the frame. Accepted, and it is the substantive part of this
disposition.** The rule in `execution-scope.md` 6 was scoped to registries the
strict path consults *for a Phase 23 family*. Three consecutive sweeps - QA's
two and my own - walked only that path, which is exactly why R-19 survived all
three. The rule now reads: **every closed journal registry that a Phase 23 value
or Effect reaches, whichever family carries it**, with the two entry points named
- a Phase 23 family on the strict path, and a Phase 23 value inside a live or
retained Phase 21/22 family. Rows 1-6 of the inventory are the first kind; rows
7 and 8 are the second.

**I-32, the count. Accepted, prose only.** The JSON column enumeration is
**twelve**, not eleven - task 4 five, task 8 three, task 9 two, task 10 one, task
12 one. The per-task assignment was already complete and does not move. Corrected
in DEC-232's R-17c paragraph and in the confirmation brief.

**I-34, `ErrorSubjectType`. Accepted, with the implementation named
authoritative.** Measured: the phase-21 prose lists 22 names, the implemented
`SubjectType` at `internal/cli/codes.go:88` has 19. Five documented names
(`actor`, `idempotency`, `obligation`, `path`, `repository`) were never
implemented; two implemented names (`attention`, `stop`) were added by Phase 22
without amending the prose. The section 2 row already cites the implementation as
its referent, so the implementation is authoritative and the phase-end set is
**25**, not 28. Phase 23 adds the six named and nothing else; reconciling the
prose is a spec-sync item at wrap-up. The row also now states the registration
rule, which was implicit: the task that first returns an error carrying a new
subject registers its constant in the same change - 4, 5, 8, 8, 9, 10 - the same
rule DEC-229 set for Condition types.

**Tech Lead finding 1, `semanticPhase22Events` has no Phase 23 analogue.
Accepted; owner is task 1.** Found while checking whether export mirrors any
closed **value** registry. It does not - `validateSemanticMap` validates by DDL
column type and `loadSemanticTables` introspects the installed schema, so the ten
new tables arrive from migration 0017 free - but `validateEvent`
(`archive_semantics.go:328`) rejects an Event when `!found || (profile == 14 &&
semanticPhase22Events[type])`. Once task 1 mirrors the nine families into
`semanticEventRules` under DEC-229, `found` is true for them at every profile, so
a legacy projection-14 archive carrying `review.recorded` validates clean. That
is the precise hole Phase 22 added the map to close. Task 1 adds the nine in the
same file and change as the mirror and corrects the map's name. Task 1 is the
right owner for the same reason it owns `current_event_columns.go` and the
unfencing tasks do not: this is provable with a crafted archive and no Event.
VAL-065 is unaffected and already holds by the manifest check at
`archive_validation.go:130`.

**Tech Lead finding 2, the `Condition` nullable-binding relation fails open.
Accepted as inventory row 8; owners are tasks 9, 10, and 11.** Same file as
R-19, different case, found by asking what else in that file a Phase 23 value
reaches. `case "Condition"` enforces the accepted "only the field matching the
type is non-null" rule for `budget_name` and `capability_id` and nothing else.
Reproduced under overlay with the shape widened, working tree untouched: a
`human-outcome-recorded` Condition with all three new fields non-null, with none
non-null, and a `decision-satisfied` Condition carrying `outcome_class` are all
**accepted**, while the Phase 22 control still returns `condition nullable
binding mismatch`. It fails open, so it blocks nothing and is not a dispatch
blocker - it is a silent loss of an accepted invariant, the R-09 and R-18 class.
Each of tasks 9, 10, and 11 extends the branch for its own field alongside the
`validateCondition` change it already owns.

**Tech Lead finding 3, how `Condition` may be widened at all. NOT
dispositioned - raised as ESC-23-001.** Checking finding 2 showed that the
widening itself is not monotonic, and that `Condition` is the only field-adding
widening of a live journal value shape in the phase; every other section 2
amendment widens a closed enum or lands on a retained family. Two measured
consequences, neither mentioned in any Phase 23 document - a search for
`semantic_key`, `EventSchema`, and `omitempty` across the contract set and all
eighteen task specs returns nothing. First, `validateCurrentShape` matches field
count exactly, and `ReplayJournal` revalidates every historical Event against
current shapes, so after the widening no workspace holding a pre-widening Stop
can rebuild, and `UpgradeToCurrent` has no hook because Phase 23 does not bump
`EventSchema`. Second, `attentionSemanticKey` hashes the whole struct and
canonical JSON preserves nulls, so every stored `semantic_key` stops matching and
an identical Stop raises a duplicate Attention. Both remedies change an accepted
value's encoding or add unaccepted migration machinery, which is a Human
decision of the DEC-223..227 class, not a planning correction. Tasks 1-8 are
unaffected and may dispatch; tasks 9, 10, and 11 are held.

### DEC-234: Tech Lead dispositions of the sixth QA confirmation pass

- **Status:** ACCEPTED by Tech Lead 2026-09-19 under planning authority. No
  protected definition, command, threshold, budget, or accepted assertion
  semantics changes; no `VAL-NNN` is added, renumbered, narrowed, or reworded,
  and the contract stays at 79 rows. No spec file outside the phase directory is
  edited. Three of the pass's revisions are deliberately **held**, not deferred
  silently - see the last paragraph.
- **Context:** the QA Lead's sixth confirmation pass (`qa-plan.md` section 17)
  returned **READY EXCEPT ESC-23-001**: tasks 1-8 unaffected and dispatchable,
  DEC-233 confirmed correct in full, ESC-23-001 verified in every particular and
  stronger than recorded in two respects. Four Must-fix revisions (R-20..R-23)
  and four improvements (I-35..I-38). Everything below is a transcription of
  measurements the pass already recorded and cited; nothing here is newly
  derived, which is why no confirmation pass is requested for this batch.

**R-20, ESC-23-001's options were stated on wrong costs. Accepted in full and
applied; this was the priority, because the Human has not yet read the record.**
Four substantive amendments and two citation fixes, all transcribed from
`qa-plan.md` 17.2-17.4:

1. **Option A's encoding was ambiguous and one reading destroys it.** "Each new
   type carries only its own field" admits a reading that puts `omitempty` on
   `budget_name` and `capability_id` as well; measured, that reading changes
   **every** existing semantic key (`38cda009…` -> `cf865c69…`, `da2f2a11…` ->
   `ef2f1a6a…`, `28607d77…` -> `9fb954bc…`) - the exact harm Option A exists to
   prevent. The option now states `omitempty` on the **three new fields only**,
   with the six Phase 22 types keeping `budget_name` and `capability_id`
   present-and-null, so a new type carries five fields, not three. (*Corrected
   by DEC-236 / R-26: three judgment types carry five fields; `review-recorded`,
   which section 11 gives no discriminator, carries four. The five-field
   generalization was written before that was settled - do not re-derive the
   arithmetic from this line; the four authorities are the task specs,
   `contracts.md` 11, `public-interface.md` 2, and VAL-080.*) That reading
   measures byte-identical to baseline and replays the recorded pre-widening
   journal clean.
2. **Option A is one file larger than stated, and the obvious implementation
   panics.** The discriminant lives inside the value and `Condition` is also
   reached generically as a `Stop` field
   (`current_event_value_types.go:39`), so per-type dispatch needs a branch in
   `internal/journal/current_event_value_types.go`, which no task owns. The
   option now names it, requires it to join the implementing task's scope, and
   requires a **checked** type assertion: the naive `object["type"].(string)`
   panics `interface conversion: interface {} is int, not string` in code that
   validates untrusted archive bytes, where the baseline degrades cleanly to an
   invalid closed enum.
3. **Option B costs more than recorded.** Four payload sites rather than one
   (`stops.resolution_condition_json` plus the nested `Stop` in the typed
   results of `stop.raised`, `stop.resolved`, `run.budget-observed`), plus an
   `idempotency_results.request_digest` recompute, plus 11 failing tests in the
   shipped corpus, all named; against Option A's zero failures and zero fixture
   edits across `internal/journal`, `internal/workflow`, `internal/export`.
   Option B's genuine merits are now stated too, so the Human sees a fair
   comparison rather than an advocate's.
4. **A third consequence is added to the Evidence section.** `stopDigest`
   canonicalises `StopRaiseRequest`, which embeds the `Condition`, so the
   widening turns an idempotent `stop.raise` retry into
   `idempotency conflict for stop.raise:project-… key "retained-test-1"`.
   The Summary's "two independent reasons" becomes three.

Also applied: the QA Lead's harm-population framing (17.4) as a new
**"Who this reaches, and what is actually being decided"** section. It records
that `schema.md` section 1's accepted fresh-init policy plus `golang-migrate`'s
never-re-run behaviour means there is no released Phase 22 population, so the
question is "what encoding do we want for V1" rather than "how do we protect
user data" - while keeping both reasons the escalation survives that framing:
the incompatibility is in-phase at tasks 9-11 against tasks 1-8 workspaces, and
the Attention and idempotency consequences need no rebuild at all. This
strengthens the recommendation and was recorded because it changes what the
Human is choosing between, not because it changes the answer.

**R-21, the ninth registry. Accepted as an owned inventory row, not as a
positive.** `validateCurrentBudgetFacts`
(`internal/journal/current_event_budget_values.go:34-50`) returns `nil` at its
`default` for all nine Phase 23 families - measured on every one - so task 6's
`review-cycles` charge has no journal-side sample/state consistency check,
automatic-charge rule, or exhaustion-Stop derivation. QA offered either a row or
an explicit positive. **Row 9, owned by task 6**, is the better of the two: a
positive tells a later sweep not to re-derive it, which is true but insufficient
here, because this one also tells a **Developer** something they would otherwise
assume wrong. It is the third kind of inventory entry - a registry a Phase 23
family does reach and which declines to constrain it - and the frame's closing
rule already covers it ("silently fails to constrain a Phase 23 value"). The row
requires **no code** from task 6: producer-side enforcement of the charge is
accepted for this phase. What task 6 gains is the evidential statement, now in
its spec, that VAL-027 and VAL-029 are producer-side proof only and must not be
reported as journal-layer proof. It fails open, so it blocks nothing.

**I-35 accepted:** task 1's `semanticPhase22Events` addition gains an acceptance
criterion - a crafted legacy projection-14 archive carrying `review.recorded` is
rejected - with **no new VAL ID**, because VAL-065 correctly does not cover it
(the legacy branch admits profile-14 archives by design) and the phase does not
add contract rows for defensive in-task assertions.

**I-36 accepted, with one correction of the correction.** Five anchors fixed:
inventory row 8 `:48` -> `:49`, positive 2 `:76` -> `:77`, positive 3
`current_event_value_types.go:88` -> `:89`, and in ESC-23-001
`current_event_json.go:26` -> `:27` and `openAttentionBySemanticKey` ->
`loadOpenAttentionByKey` at `stop.go:557`. A sixth: the pass's own citation of
the generic shape lookup as `current_event_value_types.go:38` is off by one -
`if shape, ok := currentValueShapes[kind]; ok` is at `:39`, verified in the
tree, and `:39` is what the escalation now carries. Every anchor written into a
task-facing or Human-facing document in this batch was re-verified against the
working tree before transcription, not copied.

**I-38 accepted:** `task-10.md` grammar fixed.

**Held for the ESC-23-001 disposition batch, and recorded in the escalation's
Required resolution so they cannot be lost: R-22, R-23, and I-37.** Each is
bounded and accepted in principle; each must be **phrased against whichever
encoding the Human chooses**, so writing them now would mean writing them twice
and risking a stale phrasing surviving the rewrite. R-22 (name the proving Event
for the entries tasks 9 and 10 add, since three of DEC-224's four raising
families are retained and never exercise them) and R-23 (`task-11.md` has no
acceptance criterion citing VAL-078) touch only tasks 9, 10, and 11, which
ESC-23-001 already holds; I-37 is one sentence in `schema.md` section 1 and is
part of the same encoding question. Recorded with them: if the Human chooses
Option A, inventory row 8 is **restated, not deleted** - the invariant survives
and only its enforcer moves from the relation branch to the shape - and VAL-078
keeps its assertion while its evidence class becomes shape rejection. That batch
does need a fresh QA pass, because it amends an accepted public value and
rewrites three task specs; this one does not.

### DEC-235: Human resolves ESC-23-001 with Option A, and the batch it carries

- **Status:** **ACCEPTED by the Human 2026-09-20**, choosing Option A in the
  corrected reading R-20 established. This amends an accepted public value's
  encoding, which is Human authority; the consequent ownership, file-scope, and
  task-spec edits below are Tech Lead planning authority applied under it. One
  `VAL-NNN` row gains an evidence-class note; no assertion, method, pass or fail
  criterion, validator, or owning task changes, and the contract stays at 79
  rows.
- **Context:** ESC-23-001 established that widening the journal `Condition`
  value is not monotonic over stored history in three measured ways - journal
  replay, Attention deduplication, and `stop.raise` idempotency. Option A makes
  the three new fields absent on the types they do not apply to; Option B keeps
  the accepted present-and-null encoding and pays for it with a V3 schema bump,
  a four-site payload upgrader, two column recomputes, a mid-life Attention key
  rule, and 11 corpus repairs. The QA Lead implemented and measured Option A
  before the decision (`qa-plan.md` 17.3): byte-identical Phase 22 semantic keys,
  a clean replay of the recorded pre-widening journal, and zero failures and
  zero fixture edits across `internal/journal`, `internal/workflow`, and
  `internal/export`.

**The accepted encoding, stated once and authoritative.** `omitempty` on the
three new fields - `outcome_class`, `decision_class`, `finding_id` - **and on no
other field**. `budget_name` and `capability_id` keep present-and-null
treatment. A judgment Condition therefore carries **five** fields and a Phase 22
Condition is unchanged at four, byte-for-byte. The alternative reading, which
would also tag `budget_name` and `capability_id`, is **rejected**: measured, it
changes every existing semantic key and reproduces the exact harm the decision
prevents. `public-interface.md` section 2's `Condition` row is amended to say
this, including the superseded reading and why.

**Exclusivity moves from a relation to a shape.** The journal registers per-type
`Condition` shapes keyed on `type`, as `currentStopDetailShapes` is keyed on
`code`. Measured consequences: the wrong field rejects with `missing required
field`, an own field that is null with `invalid ID`, two fields with `closed
object has wrong shape`, and the own field alone is accepted - all four cases
inventory row 8 left open, closed by the shape itself.

**Ownership, which is where this decision costs something.**

- **Inventory row 8 is restated, not deleted.** The relation branch keeps the
  Phase 22 `budget_name` / `capability_id` binding, because both fields stay
  present-and-null on every type. (*Corrected by DEC-236 / I-40: "which
  per-type shapes cannot express" was wrong - the grammar's `Null` token could
  express it. The relation is kept by choice, not necessity, and the branch must
  still be **reached** by the new dispatch; see R-24 under DEC-236.*) **No Phase
  23 task extends it**, and all three task specs now say so explicitly - a task
  that adds a third boolean term has misread this decision. Deleting the row
  would have been the easy error: the invariant survives, only its enforcer
  moves.
- **New inventory row 10**, the per-type shape map and its `validateCurrentValue`
  dispatch. It fails **closed**, so each of tasks 9, 10, and 11 adds its own
  type's row. **Task 9 creates the map and the dispatch branch and takes
  `internal/journal/current_event_value_types.go` into its file scope** - the
  file R-20 found that no task owned. Ownership is assignable there because task
  9 is the first task that can execute it.
- **The dispatch must use a checked type assertion**, with an acceptance
  criterion in task 9. The naive `object["type"].(string)` panics on a
  non-string discriminant in code that validates untrusted journal and archive
  bytes; the baseline degrades to an invalid closed enum. A panic here is a
  crash on a crafted archive.
- **Byte-identity of the six Phase 22 shape rows is an acceptance bar**, not a
  hope: task 9 must show a recorded pre-change journal replaying clean and an
  unchanged Condition keeping its `attentionSemanticKey`. Any drift there
  reintroduces the escalation.

**VAL-078's evidence class is annotated, and nothing else about it changes.**
Its closed-schema half is now satisfied by shape rejection rather than by the
`condition nullable binding mismatch` relation. Evidence citing the relation
message would indicate the encoding was implemented wrong, so the note is
recorded in the contract row and in `task-11.md`. The assertion, method, pass
and fail criteria, validator, and owning task are untouched.

**The three held revisions are applied with this decision, as promised.**

- **R-22 accepted, and strengthened from documentation to an executable
  assertion.** QA asked that the proving Event be named, because three of
  DEC-224's four raising families are retained and return at
  `current_event_validation.go:64` before any value validation, so tasks 9 and
  10 could omit their registry entries, watch their raise tests pass, and fail
  at task 11 or tasks 15/16 with the cause outside the failing task; task 9's
  `protected-decision-required` has no strict raising path at all. Naming an
  Event is necessary but still leaves the proof outside the task. Both specs
  therefore now require a **direct strict-path append** inside the owning task,
  following `TestCurrentEventRejectsConditionalBindingCorruption`
  (`internal/journal/current_event_validation_test.go:171`), and both state that
  green raise tests are not evidence for these entries. (*Corrected by DEC-236 /
  R-27: that exemplar is wrong for the `currentStopDetailShapes` half, because
  it mutates `executionBlockedProjectionEvent` and `run.execution-blocked` is a
  `validateCurrentBudgetFacts` family. Both task specs now cite
  `TestCurrentClosedValuesRejectIncompleteAndInvalidFacts`
  (`internal/journal/current_event_closed_values_test.go:76`) and are the live
  authority; the requirement itself is unchanged.*) This is the phase's own
  rule applied honestly: ownership is only assignable where it is executable, so
  where a raise cannot prove an entry, the task asserts it directly.
- **R-23 accepted:** `task-11.md` gains the missing VAL-078 acceptance
  criterion, covering both halves and phrased against the accepted encoding.
- **I-37 accepted:** `schema.md` section 1 records that under `golang-migrate` a
  workspace migrated in place keeps its original CHECK text while reporting
  schema 17, so the first Phase 23 Stop there fails with a CHECK-constraint
  error rather than a re-init instruction. Intentional under the accepted
  baseline-replacement policy, reopening no decision, recorded so the failure
  mode is recognised.

**What this decision does not do.** It does not bump `EventSchema`, add an
upgrader, or touch migration 0017's data path - the whole point of Option A is
that none of that is needed. It does not change any Stop code, Condition type,
grant action, failure code, or evaluation reason. And it does not relieve tasks
9, 10, and 11 of `validateCondition`: producer-side exclusivity stays, with the
journal shapes as the second, now-structural line.

### DEC-236: Seventh-pass dispositions - the relation call, the six-row bar, the details route, and `review-recorded`

- **Status:** **ACCEPTED 2026-09-20 under Tech Lead planning authority.** No
  Human decision was required; the authority question R-26 raised is answered
  below and resolved on the accepted contract, not by amending it. The
  validation contract stays at **79** rows; one row (VAL-078) is refined and
  strengthened, and no assertion is weakened, narrowed, or retired.
- **Context:** QA pass 7 (`qa-plan.md` 18, commit `9b28ba2`) re-verified the
  DEC-235 batch in a fresh session under `go test -overlay` with the working
  tree untouched. It confirmed all four of the Tech Lead's asks, found Option A
  transcribed correctly in all three specs, and returned four Must-fix revisions
  and three improvements. Verdict: READY except tasks 9, 10, and 11. The Tech
  Lead independently re-derived the structural facts behind R-24, R-26, and R-27
  from the source before dispositioning them.

**R-24 accepted: the new dispatch must still reach the relation.** This is the
one finding that would have shipped a silent regression.
`validateCurrentValue`'s generic path ends
`return validateCurrentValueRelations(kind, object)`
(`current_event_value_types.go:39-44`) - verified in the source - and DEC-235
removes the `"Condition"` key from `currentValueShapes`, which deletes that
path's only call for this kind. All three specs said "do not touch
`case "Condition"`"; none said the new branch must still call it. Measured by
QA: a branch that selects a per-type shape and returns accepts a `work-ended`
Condition carrying `budget_name:"tokens"`, **appends and projects it** on the
strict path, and passes the entire ordinary suite `exit 0`. `condition nullable
binding mismatch` occurs exactly once in the repository - the `fmt.Errorf` that
produces it - so no test detects the loss, and DEC-235's own four exclusivity
criteria pass with row 8's invariant already gone. Task 9 gains requirement
**(c)** and a dedicated acceptance criterion; tasks 10 and 11 confirm the
relation still fires in their completion reports; inventory rows 8 and 10 carry
the same sentence. This is the registry-without-an-owner failure shape one level
down: correct ownership, correct prohibition, missing statement of what must
still run.

**R-25 accepted: the byte-identity bar covers all six rows, not one.**
`run-recovered` and `capability-available` reach journal Condition validation
**zero** times in the ordinary suite and are raised nowhere in the tagged
corpus; QA deleted both rows and the suite still passed `exit 0`. A dropped or
mistyped row would ship a binary rejecting those Conditions at runtime with
every gate green. Task 9's criterion becomes a six-type table - one canonical
Condition of each accepted, one mutated Condition of each rejected - retaining
the replay and `attentionSemanticKey` halves. This directly serves the phase's
stated residual risk, "a Phase 22 shape row not transcribed exactly".

**R-27 accepted: the R-22 assertion keeps its force and gains a route that
works.** The requirement itself stands - without it tasks 9 and 10 can omit a
registry entry, watch green raise tests, and fail at task 11 or 15/16 with the
cause outside the failing task. Only the cited exemplar was wrong for the
`currentStopDetailShapes` half:
`TestCurrentEventRejectsConditionalBindingCorruption` mutates
`executionBlockedProjectionEvent`, and `run.execution-blocked` is one of the
five families in `validateCurrentBudgetFacts`
(`current_event_budget_values.go:34-50`, inventory row 9 - verified in source),
so a changed Stop code is rejected `current budget boundary: unrelated Stop`
before any detail shape is consulted. Tasks 9 and 10 now cite
`TestCurrentClosedValuesRejectIncompleteAndInvalidFacts`
(`current_event_closed_values_test.go:76`), which already asserts on family
`stop.raised` for `resolution_condition_json` and on `stops.details_json` with a
`code`, needs no database, and rejects `unknown Stop code` when the entry is
missing. Two clarifications ride along: a proof that appends into the projection
depends on migration 0017's widened `stops.code` / `attention_items.code` CHECK
from task 1, and a **real `stop.raise` on the strict `stop.raised` family is
acceptable evidence** - `validateStopRaise` (`stop.go:363-374`) checks code and
condition independently. What is insufficient is a raise on the task's **own**
retained family, which the earlier blanket wording over-stated.

**R-26 accepted as branch A: `review-recorded` carries no discriminator, and
this is not a Human decision.** `task-11.md`, rewritten by the DEC-235 batch,
required `finding_id` on both `review-recorded` and `finding-resolved`.
Accepted `contracts.md` section 11 enumerates the pairing exhaustively -
`outcome_class` on `human-outcome-recorded`, `decision_class` on
`decision-satisfied`, `finding_id` on `finding-resolved`, "only the matching
field is non-null" - and gives `review-recorded` none. Four types and three
fields do not divide evenly; one type is the odd one out, and section 11 already
says which. Three reasons this is branch A:

1. **Authority.** DEC-235's Human approval was of an *encoding mechanism* -
   `omitempty` on three fields, structural exclusivity through per-type shapes -
   not of a per-type field assignment, which accepted `contracts.md` section 11
   had already fixed under DEC-222. Branch B would have added a field to a
   public value on a second type, against both section 11 and
   `public-interface.md` section 2; that would be Human authority, and it is not
   what the accepted set says. Branch A restores the accepted reading and amends
   no encoding, so it is ordinary planning correction. The pairing sentence was
   a Tech Lead drafting error introduced in `932fbea`, not a decision.
2. **Semantics.** A `review-recorded` Condition resolves on a `review.recorded`
   Event scoped by its `subject`; naming a Finding would duplicate
   `finding-resolved` rather than discriminate anything.
3. **It costs nothing structurally.** The per-type shape rows are
   field-count-exact, so the four-field `review-recorded` row is what makes a
   `review-recorded` Condition carrying `finding_id` reject with `closed object
   has wrong shape`.

**VAL-078 is refined, and the refinement strengthens it.** Its assertion said
"exactly one Condition field is non-null per type" with a fail criterion of "two
or **zero** non-null fields", which under the accepted contract was
**unsatisfiable** for `review-recorded` - every well-formed instance has zero.
The row now names the section 11 pairing explicitly and requires
`review-recorded` to **reject** any discriminator, which is an assertion the old
wording did not make. Nothing is narrowed: the resolution half, evidence method,
validator (QA Lead), and owning task (11) are unchanged, and the amendment log
records the refinement. Under the standing rule, an assertion that cannot be
satisfied is escalated rather than weakened - here the assertion was
mis-specified against its own cited authority, and correcting it against that
authority is the repair, not a relaxation.

**The 18.12 conditional blocker candidate is not raised, and no ESC identifier
is allocated.** It was explicitly conditional on branch B. Branch A is chosen,
so no accepted public-value encoding changes, no Human authority is engaged, and
no escalation exists to register. The Human is informed of the disposition
rather than asked for a decision.

**Improvements, all three accepted.**

- **I-39:** `task-11.md` gains the present-and-null sentence and the pointer to
  task 9's byte-identity bar, so all three `Condition`-editing tasks carry the
  same reason and not just the same prohibition.
- **I-40:** inventory row 8's recorded premise - "which the per-type shapes
  cannot express" - is **wrong and is corrected in place**. The grammar has a
  `Null` token (`current_event_value_types.go:59-60`), in live use at
  `current_event_value_catalog.go:21`, so rows could express the binding with no
  change to any marshaled byte. The conclusion is unchanged - keeping the
  relation is the smaller change and the right one - but it is kept **by choice,
  not by necessity**, and the corrected reason ("the accepted rows keep both
  fields `?ID` on every type, so the relation is their only enforcer") is also
  precisely what makes R-24 necessary. DEC-235 restated the same wrong premise
  and is annotated rather than rewritten.
- **I-41:** `task-18.md` now owns extending `p22ClosedFieldGrammar`'s
  `Condition` row (`cmd/fid/integration/p22_closed_values_test.go:25`) to all
  ten types with its asserted member counts. It is a self-contained wire oracle,
  so new Condition types do not fail it - it silently stops covering them, which
  is exactly the failure mode a conformance audit exists to prevent.

**What this decision does not do.** It changes no Stop code, Condition type,
grant action, failure code, or evaluation reason; adds and retires no `VAL-NNN`;
touches no required command, threshold, or budget; and reopens neither DEC-235
nor ESC-23-001. Tasks 1 to 8 were unaffected by the pass and are unaffected by
this disposition.

### DEC-237: Eighth-pass dispositions - split VAL-078, and give the two DEC-235 preservation invariants an ID

- **Status:** **ACCEPTED 2026-09-20 under Tech Lead planning authority.** No
  Human decision is required: the split preserves every case, type, evidence
  class, validator, and owning task, and the two additions assert behavior that
  three task specs already require. Creating rows under planning authority has
  precedent in this phase - DEC-228 created VAL-075 to VAL-078, and R-11 split
  VAL-002. The validation contract moves from **79** to **82** rows: 78
  deterministic and 4 Human. Nothing is removed, renumbered, or weakened.
- **Context:** two independent QA confirmations of the DEC-236 batch arrived
  against the same base `a0e8e7e`. `qa-plan.md` section 19 (commit `3907b7e`)
  returned **READY** with no Must-fix, four improvements (I-42..I-45) and one
  validation-contract proposal (P-01), and measured R-24's acceptance criterion
  as discriminating - six cases accepted by a branch that returns after shape
  validation and rejected by one that ends in the relation. `qa-plan.md` section
  20 (commit `7a865d8`) confirmed R-24 through R-27 as transcribed and returned
  **NOT READY for task 11** on one finding, R-28, a blocker candidate against
  VAL-078's granularity with no ESC allocated. Both passes converge on the same
  row from opposite directions: one says the row now covers a case it could not
  before, the other says the row was never one observable. Both are right.

**R-28 accepted: VAL-078 is split, not reworded.** The row asserted both that
the four Condition types resolve only on their own qualifying fact **and** that
each carries exactly its assigned discriminator. Those fail independently - a
shape validator can be exactly right while resolution accepts a foreign fact,
and resolution can be exactly right while malformed Conditions are stored - and
the row's own fail criterion was an `or` across the two. The practical cost was
already visible in the task pack: tasks 9 and 10 each produce encoding evidence
and are marked "evidence toward VAL-078, which task 11 completes", so the
encoding half had no recordable state until task 11 closed, and a final-review
`FAIL` would not have said which half failed. VAL-078 keeps the resolution half,
its QA Lead validator, and task 11. The new **VAL-080** carries the encoding
half with every case intact, the DEC-236 scoping that made it satisfiable for
`review-recorded`, and the DEC-235 evidence class: shape rejection, never the
relation message. Its owning-task cell records what was already true - tasks 9
and 10 evidence the type each raises, task 11 completes the four-type matrix.
This is the second bundled row in this phase; R-11 split VAL-002 on the same
pattern, a row asserting both that a thing behaves correctly and that it is
encoded correctly. The rule is now written into the contract's coverage rules so
the third one is caught at authoring time.

**P-01 accepted: the Phase 22 nullable binding becomes VAL-081.** The binding -
`budget-observed` alone carries a non-null `budget_name`, `capability-available`
alone a non-null `capability_id` - is live public behavior that this phase's
refactor can delete with every gate green. That is not a hypothesis: QA measured
a branch that returns after shape validation, which accepts and projects a
`work-ended` Condition carrying `budget_name:"tokens"` and passes the entire
ordinary suite. `condition nullable binding mismatch` appears once in the
repository, in the `fmt.Errorf` that produces it, and no other Phase 23 gate
would notice its loss: a valid Phase 22 journal replays clean with or without
the relation, so tasks 15 and 16 are blind to it, and task 18's extended
`p22ClosedFieldGrammar` is a self-contained wire oracle that never reaches the
journal. Task 9's new test will be the only detector that has ever existed. A
`VAL-NNN` makes its survival a recorded final-review outcome instead of one
checkbox in one completion report. Task 9 owns it; tasks 10 and 11 re-confirm it
after adding their rows, on the VAL-003 precedent.

**Extended by the Tech Lead: the six-row compatibility bar becomes VAL-082.**
P-01's argument does not stop at the relation. DEC-235 is a Human-approved
contract amendment, and every other Human-approved amendment in this phase
carries at least one assertion - DEC-223 (VAL-005), DEC-224 (VAL-046), DEC-225
(VAL-015/022/077), DEC-226 (VAL-025/039/056/057), DEC-227 (VAL-079). DEC-235
carried none, and the invariant it exists to protect is the larger of the two:
`omitempty` on `budget_name` or `capability_id` changes **every** stored
`semantic_key` and `request_digest` in existence, measured, with three key pairs
recorded above. It has the same profile P-01 describes - `run-recovered` and
`capability-available` reach journal Condition validation zero times in the
whole suite, so deleting both rows leaves every required gate green - and the
same single detector, task 9's six-type table. VAL-082 asserts that table plus
the replay and `attentionSemanticKey` halves as one observable: the six retained
types are unchanged. Byte identity is the property; the three legs are how it is
measured, and the latter two follow from the first.

**Two QA evidence routes declined a row, with the reason recorded.** Neither is
a coverage gap, and the boundary is deliberate: a row is warranted when the
property is publicly observable, this phase can destroy it silently, **no**
other gate would catch the loss, and the evidence is already mandated, so the
row costs no new work.

- **QA-52, a non-string `type` discriminant rejects without panic.** Declined: a
  panic fails task 9's own required gate loudly. The property has a detector, so
  it stays a task acceptance criterion. Task 9's criterion now records the
  decline rather than reading as an unexamined "No new VAL ID".
- **QA-49, each owning task's `currentStopDetailShapes` entry is reached.**
  Declined: a missing entry is caught, just late and badly - as DEC-235's R-22
  paragraph says, at task 11 or tasks 15/16, with the cause outside the failing
  task, and VAL-064's rebuild leg is its end-state backstop. The defect there is
  locality, and R-22's direct-assertion requirement is the right fix for
  locality. A contract row asserts an end state, and the end state is already
  asserted.

**Improvements, all four accepted and applied.**

- **I-44** is the one that mattered before dispatch. Task 9's six-type table now
  says **canonical means type-appropriate**: `budget-observed` carries a
  non-null `budget_name` and `capability-available` a non-null `capability_id`,
  because an all-null form of either is a `condition nullable binding mismatch`
  rejection - a VAL-081 case, not an acceptance here. Without the clause the two
  types with zero corpus coverage, which are the two the table exists for, would
  each have cost a red-test cycle on their first fixture.
- **I-45:** tasks 10 and 11 carried "confirm in the completion report that the
  relation still fires" in a file-scope bullet only, so a report that omitted it
  still ticked every box. Both now have an acceptance criterion citing VAL-081,
  which also gives the re-confirmation an ID to report against.
- **I-42:** DEC-235's R-22 paragraph still directed tasks 9 and 10 to follow
  `TestCurrentEventRejectsConditionalBindingCorruption` for the details half,
  which R-27 measured as wrong. Annotated in place, in the same form as the I-40
  annotation; the live authority is both task specs and they are correct.
- **I-43:** the over-generalized "five-field judgment Condition" survived in
  three derived places - `tl-escalations.md`'s resolution bullet, DEC-235's
  Option A description, and `contract-coverage.md`'s clause title. All three now
  state that only the three discriminated types carry five fields. This is the
  staleness that produced R-26 in the first place: a later pass re-derived the
  arithmetic from a line that predated the finding.

**What this decision does not do.** It changes no Stop code, Condition type,
discriminator pairing, encoding rule, grant action, failure code, or evaluation
reason; retires or weakens no `VAL-NNN`; touches no required command, threshold,
or budget; adds no Human checklist step; and reopens neither DEC-235 nor
ESC-23-001, which stays RESOLVED. It adds no implementation work: every new row
maps onto an acceptance criterion the task specs already carried, and the two
I-45 checkboxes promote an obligation that was already written. **Tasks 1 to 8
are untouched** - both passes checked this by name, and none of the edited
criteria, rows, or files is theirs.

**Consequences.** One fresh QA confirmation pass is required over these contract
edits before tasks 9, 10, and 11 dispatch - R-28's own requested action, and the
right bar for a Tech Lead edit to the contract. It is bounded: the three rows,
their task-spec citations, the count, and the revision log; R-24 through R-27
need no recheck, because nothing here changes them. Independently of that pass,
**tasks 9, 10, and 11 remain held until the Phase Runner receives explicit Human
clearance**; no QA verdict and no Tech Lead disposition releases that hold.

### Accepted carry-forward dispositions

| ID | Phase 23 disposition |
|---|---|
| CF-008 | Re-defer to Phase 24 with the Envelope disposition work for removed unfinished Issues. |
| CF-018 | Re-defer; no planned production edit to `internal/sqlquery/sqlquery.go` in this phase. |
| CF-023 | Re-defer to dedicated workflow maintenance. |
| CF-026a | Re-defer; unscheduled product candidate. |
| CF-026b | Re-defer; separate Human product decision. |
| CF-029 | Re-defer to dedicated workflow maintenance. |
| CF-030 | Re-defer to dedicated workflow maintenance. |
| CF-033 | Re-defer to dedicated workflow maintenance. |
| CF-036 | Re-defer to the dedicated workflow retrospective. |
| CF-037 | Re-defer to Phase 25; Phase 23 must not claim any coordinator obligation. |
| CF-038 | Re-defer; the trigger fires only on a causal edit to `cmd/fid/live_output_test.go`. |
| CF-039 | Re-defer the executable cross-check; Phase 23 adopts the manual mitigation in `contract-coverage.md` and records whether it worked in the wrap-up workflow notes. |
| CF-040 | Partially scheduled: the guidance is applied to the Human views this phase introduces (task 14, VAL-073). The polish and verbosity review of pre-existing views stays with Phase 25. |
| CF-041 | New: hosted CI runs only the ordinary suite. The tagged integration suite and the rebuild benchmark are local-only gates, so a green PR is not phase-level evidence. Owner: workflow and CI maintainer, in dedicated maintenance. |
| CF-042 | New (QA plan I-16); cause and mitigation corrected under DEC-229 (readiness refresh R-15). `cmd/fid/integration/main_test.go` leaks built-binary directories per run, not per kill: five tests re-execute the test binary, each child re-enters `runIntegrationSuite` and rebuilds its own ~16 MB binary, and the crash-process children exit abnormally by design and never reach `RemoveAll`. One clean tagged run - exit 0, no kill, 551 s - leaked 20 directories and 306 MB into a 3.9 GB RAM-backed `/tmp`; at ~26 tagged runs this phase that is ~7.8 GB, filling `/tmp` around the twelfth. Task 1's mitigation is therefore causal: publish the built binary path so a re-executing child reuses it and neither builds nor removes, keeping the 24-hour prune only as the residual safety net for a parent killed mid-run. The durable cleanup-on-kill design stays with the workflow and test-infrastructure maintainer in dedicated maintenance. |

### DEC-238: Human delegates ESC-23-002 causal repair and resume preparation

- **Status:** ACCEPTED by explicit Human delegation, 2026-09-20.
- **Authority:** The Human authorized this agent to act above the fixed role and
  workflow boundaries, complete the proposed repairs and documentation, close
  ESC-23-002, commit, and prepare a Phase Runner handoff. This is a bounded
  exception for this escalation, not a standing workflow policy change.
- **Scope:** Task 1 recovery may repair rebuild SQL preparation in
  `internal/journal` and `internal/db`, the cached execution-transition fixture
  in `cmd/fid`, and causally identified test-fixture/runtime overhead needed to
  meet the retained ordinary-package budgets in `cmd/fid` and
  `internal/workflow`. Matching tests, task/report, escalation, current-context,
  evidence, and resume documentation are included. No dependency is added.
- **Unchanged controls:** The 100k rebuild remains strictly below 30 seconds;
  ordinary packages remain below 600 seconds and tagged integration below 900
  seconds. Exact task commands, semantic assertions, expiry rejection, journal
  bytes, per-Event rollback/checkpoint semantics, and rebuild parity stand.
  Product Events, schemas, public interfaces, and migration SQL are unchanged.
- **Implementation direction:** Reuse prepared writes, including projection
  metadata, within one rebuild; close owned resources on success and failure.
  Remove cached-fixture age from the transition cases' qualifying premises.
  Keep the explicit expiry cases and all expected rejection reasons.
- **Evidence:** Preserve original failed attempts; record red/green, exact task
  gates, observed package times, integration binary counts/free space, three
  fresh reference benchmarks, and focused rollback/isolation/parity checks in
  `qa/esc-23-002-repair.md`. Update Task 1's report with current-session evidence.
- **Resume:** No later task is dispatched by this repair. The Runner resumes
  using DEC-238 and the committed evidence/repair, verifies Task 1 closure and
  ledger sequence, then continues tasks 2–8. Tasks 9–11 retain the Human hold.
  Earlier BLOCK/TRG/STOP records remain historical evidence, not deleted.
  The Runner may reconcile this completed Human-delegated recovery without
  redispatching implementation solely to satisfy role boundaries: verify the
  committed report/evidence, append the bound recovery and RESOLVED receipt,
  then record Task 1 COMPLETE. The precise binding is in
  `qa/resume-after-esc-23-002.md`; later QA review windows are unchanged.
- **Refs:** ESC-23-002; TRG-23-001/002; DEC-227; VAL-004/079;
  `../../reviews/2026-09-20-phase-23-rebuild-performance-diagnosis.md`.

### DEC-239: Resolve the VAL-015 evidence contradiction and authorize execution through Task 9 review

- **Status:** ACCEPTED 2026-09-20 under the Human's explicit delegation to
  resolve the current Phase 23 block and orchestrate through Task 9 and its
  subsequent mid-phase review. This is a bounded protected evidence-route
  decision, not permission to weaken independence or change the Run lifecycle.
- **Decision:** Retain same-Actor Run/Claim binding and durable Claim history.
  Amend VAL-015 to accept an explicitly isolated durable-history query fixture
  for the Run-holder-only direction. Retain real CLI mutating-Claim rejection,
  exact conflicting subject and Actor, no Event on rejection, and successful
  first and post-rework Reviews by the same independent Actor. The assertion,
  two conflicting directions, validator, task owner, task commands, and budgets
  remain unchanged. ESC-23-003's protected contradiction is resolved.
- **Evidence for the contradiction:** `internal/workflow/run.go`
  `validateRunStartBinding` rejects a Claim Actor unequal to the request Actor;
  `internal/review` checks retained mutating Claim history before the Run branch.
  The current `TestCheckIndependenceRejectsClaimAndRunHistoryDirections` inserts
  a different Run Actor into isolated SQL tables. That proves the query branch,
  not lifecycle reachability, journal acceptance, or import eligibility. No
  test or gate was executed by the Tech Lead for this decision.
- **Residual implementation obligations:** The retained Task 4 report mistakes
  per-Actor self-resolution rejection for accepted behavior and changes its
  fixture to two Reviewers. DEC-225 and `contracts.md` sections 3 and 5 already
  require the same independent Reviewer's later Review to be admissible. Repair
  that mismatch and prove same-Reviewer rework; do not reinterpret the new
  evidence route as satisfying it. Also close the report's VAL-003 UNVERIFIED
  parity row and the Task 4 criterion for resolving multiple prior Findings in
  one Event. These are existing Task 4 obligations, not new product decisions.
- **Task 4 resume:** Preserve TRG-23-003's BLOCK, the STOP, every failed attempt,
  and the blocked report as historical evidence. Resume the same trigger using
  DEC-239, attempt 9, task-scope in-scope implementation recovery, then verify
  the amended report and all unchanged exact task gates before completion.
  The exact runner binding and acceptance checklist are in
  `qa/resume-after-esc-23-003.md`. This decision is not a completion receipt.
- **Dispatch boundary:** The Human's current instruction supplies the explicit
  clearance formerly missing for Task 9. DEC-237's bounded QA confirmation
  remains required before that task starts. Tasks 10 and 11 remain held; the
  authorized run stops after Task 9 and the window 7-9 review, including causal
  repairs needed to close that review. Normal ordering and reviews after Tasks
  6 and 9 remain required. No phase wrap or Task 10 dispatch is authorized.
- **Spec reconciliation:** The contradiction was in VAL-015's evidence method,
  not the product independence guarantee. The validation contract, Task 4 and
  slice contract now state the executable proof boundary. Canonical
  `02_TRD.md` section 5.2 records the retained Actor binding/history explicitly;
  `07_DESIGN_PLAN_FLOWS.md` section 7 makes repeat-Reviewer resolution explicit;
  `08_QUALITY_PLAN.md` section 8 records the precise Phase 23 evidence class.
  PRD-REV-004, PRD-RUN-005, schema sections 3.5/3.7 and Glossary Claim/Run/Review
  Finding were inspected and remain consistent; no DDL, wire shape, Event, or
  public operation changes.
- **Refs:** ESC-23-003; TRG-23-003; DEC-225; DEC-237; VAL-003/015/022.

### DEC-240: Lossless split of retained Condition compatibility evidence

- **Status:** ACCEPTED 2026-09-20 under Tech Lead planning authority, applying
  bounded DEC-237 confirmation finding R-29. Independent QA confirmation of
  this final diff remains the Task 9 readiness prerequisite.
- **Finding:** VAL-082 combines retained Condition bytes/validation, recorded
  journal replay and Attention semantic-key stability. Replay or key derivation
  can fail while the per-type table passes; their outcomes must be recordable
  separately under the same single-observable rule that split VAL-078.
- **Decision:** VAL-082 retains the all-six canonical bytes and paired
  acceptance/mutation-rejection table. VAL-083 carries the recorded pre-change
  `budget-observed` journal replay with unchanged source digest. VAL-084 carries
  fixed-input `attentionSemanticKey` equality against the named pre-change
  reference. All three remain Task 9 obligations with QA Lead validation. Task
  9's combined criterion is split into matching checkboxes; QA-47/50/51 map to
  VAL-082/083/084. No case, reference, comparison or failure signal is dropped.
- **Boundaries:** The contract count is 84 (80 deterministic, 4 Human), with
  consecutive unreused IDs. Required commands, field encodings, nullable
  relation, product scope, budgets and Human checks are unchanged. This is no
  new implementation work and does not reopen R-24 through R-27. The final
  bounded QA confirmation checks rows, citations, count and revision log; no
  additional full readiness cycle is required. Task 9 has Human clearance
  under DEC-239; Tasks 10-11 stay held.
- **Spec impact:** Canonical quality plan section 8 already requires one
  observable per assertion. This split conforms to it; no canonical product
  spec edit is needed. `validation-contract.md`, Task 9, coverage pointers,
  QA-plan mapping and current context are synchronized.
- **Refs:** DEC-235/237; QA R-29; VAL-082/083/084.

### DEC-241: Assign the existing review-budget authoring validator to Task 6

- **Status:** ACCEPTED 2026-09-20 under Tech Lead task-scope authority after a
  bounded pre-dispatch inspection requested by the Human's delegate.
- **Finding:** `internal/workflow/envelope.go` `validateEnvelopeCreate`
  currently accepts only `budget.ExhaustionAction == "budget-exhausted"`.
  Task 6's VAL-030 requires accepting the matched
  `review-cycles`/`review-rework-exhausted` pair with a review Obligation,
  rejecting three invalid pairings, and accepting unbudgeted review. Its
  named budget files cannot implement that existing authoring guard alone.
- **Scope correction:** Task 6 owns the minimum authoring-validation change in
  `internal/workflow/envelope.go` and matching ordinary in-package tests for
  those five accepted cases. Retain every other Envelope and budget rule.
  `go test -race -timeout 30m -coverprofile=coverage.out ./...` and the existing
  `go test -race -timeout 30m -cover ./internal/workflow/` compile and execute
  this untagged scope. No command, assertion, threshold or contract changes.
- **Review:** The ordinary window 4-6 review checks the implemented behavior;
  this executable ownership correction does not open another readiness cycle.
  Product specifications already state the behavior and need no amendment.
- **Refs:** Task 6; VAL-030; contracts section 6; DEC-167/173.

## Spec Sync Verdict

Recorded at wrap-up. `verified-no-change` is a positive assertion after review,
not the default.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | pending | |
| `01_PRD.md` | pending | |
| `02_TRD.md` | touched: DEC-238, DEC-239 | section 6 records rebuild-scoped prepared writes; section 5.2 clarifies retained Run/Claim binding and history; task 17 still owns its later sync |
| `03_TRD_SCHEMA.md` | pending | expected `touched` at task 1 (section 8.4) |
| `04_TRD_CONFIG.md` | pending | |
| `05_DESIGN_PLAN.md` | pending | |
| `06_DESIGN_PLAN_CLI.md` | pending | expected `touched` at tasks 8, 9, 12, 13 |
| `07_DESIGN_PLAN_FLOWS.md` | touched: DEC-239 | section 7 clarifies same-independent-Reviewer resolution under DEC-225 |
| `08_QUALITY_PLAN.md` | touched: DEC-239 | section 8 records VAL-015 evidence-route boundary |
| `09_AI_WORKFLOW.md` | pending | |
| `AGENTS.md` (Commands) | pending | DEC-227: decide at wrap-up whether the explicit `-timeout` bound becomes the project default |
| `10_IMPLEMENTATION_ROADMAP.md` | pending | |
| `GLOSSARY.md` | touched: see DEC-226 | `Reason code` split into **Failure code** and **Evaluation reason**; applied by Tech Lead 2026-09-19, ahead of implementation, because the phase's task specs depend on the distinction |

## Triage-Log Audit

- [ ] Walked every TRG-23-NNN in `triage-log.md`.
- [ ] Confirmed each RESUME verdict was correct in hindsight.
- [ ] Flagged any TRG that should have escalated.
- [ ] Confirmed autonomous TL-Lite resume remains in effect (DEC-002).

## Workflow Notes

Recorded at wrap-up.
