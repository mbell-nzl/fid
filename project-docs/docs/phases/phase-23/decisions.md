# Phase 23 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes,
> and wrap-up decisions. DEC numbering is global and never resets.
>
> Status: accepted. Phase 23 is approved and not yet started. The QA Lead's
> first-pass plan returned thirteen Must-fix revisions and the readiness
> refresh returned two more; all fifteen are dispositioned and applied. All
> five contract amendments, DEC-223 through DEC-227, are approved by the Human
> on 2026-09-19, DEC-226 together with the part 5 vocabulary fix added during
> that review. DEC-229 and DEC-230 are Tech Lead planning dispositions of the
> readiness refresh. No decision remains open. Next: a fresh QA Lead
> confirmation of the R-14/R-15 edits, then the Runner's readiness check.

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

## Spec Sync Verdict

Recorded at wrap-up. `verified-no-change` is a positive assertion after review,
not the default.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | pending | |
| `01_PRD.md` | pending | |
| `02_TRD.md` | pending | expected `touched` at task 17 |
| `03_TRD_SCHEMA.md` | pending | expected `touched` at task 1 (section 8.4) |
| `04_TRD_CONFIG.md` | pending | |
| `05_DESIGN_PLAN.md` | pending | |
| `06_DESIGN_PLAN_CLI.md` | pending | expected `touched` at tasks 8, 9, 12, 13 |
| `07_DESIGN_PLAN_FLOWS.md` | pending | |
| `08_QUALITY_PLAN.md` | pending | |
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
