# Phase 23 - Task 9: Decision References and protected Decisions

## Objective

Record operational Decision References that point at canonical project prose,
and make an unsatisfied protected Decision block exactly its declared actions
until a matching approved reference and Human Outcome both exist.

## Validation Contract Coverage

- VAL-038: `decision.record-reference` stores a reference only and rejects both
  a digest mismatch and a `Source` whose file does not exist.
- VAL-039: An unsatisfied protected Decision blocks exactly its declared
  affected actions and returns `protected-decision-required` from them, also
  when an Obligation is unsatisfied at the same time (DEC-226 precedence).
- VAL-040: Satisfaction requires the approved reference and the approved Human
  Outcome together, including a negative for an approved reference whose
  recording authority does not match the declared required capabilities.
- VAL-041: Activation raises one Stop per unsatisfied protected Decision and the
  satisfying pair resolves that Stop and its Attention in one Event.
- VAL-042: `linked_to_decision` becomes an implemented transition rule.
- VAL-003: task 1's fence assertion for `decision.referenced` converges here
  without losing its parity subject.
- VAL-080: this task's `decision-satisfied` row carries `decision_class` and no
  other discriminator, enforced by the shape at the journal layer. Task 11
  completes the four-type matrix.
- VAL-081: the Phase 22 `budget_name` / `capability_id` nullable binding still
  fires **from the new per-type dispatch this task creates**. This task owns the
  assertion; tasks 10 and 11 re-confirm it (DEC-237, from QA P-01).
- VAL-082: the six retained Phase 22 Condition types stay byte-compatible across
  that dispatch with their per-type validation behavior retained (DEC-237/240).
- VAL-083: recorded pre-change `budget-observed` Stop history replays against an
  unchanged source journal digest (split from VAL-082 under DEC-240).
- VAL-084: the fixed subject/code/Condition derives the same Attention key at
  the named pre-change reference and after the change (DEC-240).


## Binding Constraints

- `contracts.md` section 9 is the exact authority; `schema.md` section 2 fixes
  `decision_refs`, and task 2 already landed `envelope_protected_decisions`.
- fid stores no decision body. Verify the `Source` digest against the observed
  project file under the write lock and reject `source-digest-mismatch` with no
  Event on any mismatch.
- A `rejected` or `deferred` reference never satisfies a protected Decision.
  Authority matching uses the declared required capabilities, never a role
  string or an actor-type heuristic.
- Envelope activation is permitted with an unsatisfied protected Decision;
  execution of the affected actions is not. Raise exactly one Stop per
  unsatisfied protected Decision, subject Envelope, blocked actions exactly the
  declared affected actions, resolver capability `stop.resolve`, resolution
  condition `decision-satisfied` naming that class.
- Satisfying the pair resolves that Stop and its Attention in the same Event,
  leaving every other Stop untouched.
- `linked_to_decision` moves from reserved to implemented and is evaluated
  against `decision_refs` in the Issue's scope. Do not reintroduce the removed
  escalation seed or any role-bound rule parameter.
- Human Outcome recording is task 10. Implement the satisfaction predicate over
  `human_outcomes` here and use a seeded projection fixture until task 10
  supplies the writer; do not add a temporary writer.
- DEC-192: delete exactly the `"decision.referenced"` key and prove its
  live-apply and replay parity; converge task 1's fence assertion for that key.
- Add the `decision.record-reference` grant action exactly here.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/review/` - Decision Reference rules, effects, satisfaction.
- `internal/workflow/envelope.go`, `evaluation.go`, `resolve.go`, `rules.go`,
  `registry.go`, `authority.go`, `activation.go` - the new grant action belongs
  to the closed `sortedUniqueCapabilities` allowlist in `activation.go`
  (R-08, DEC-167/173).
- `internal/journal/current_event_value_relations.go` - the **second** copy of
  that same closed allowlist (R-19, DEC-233). `validateCurrentValueRelations`'s
  `case "ActorProfile"` holds a `const actions` string with the identical 23
  names at `:63`, and nothing cross-checks the two. Add
  `decision.record-reference`
  there in the same change that adds it to `sortedUniqueCapabilities`, keeping
  both copies sorted and identical. This is not on any Phase 23 family's strict
  path: it is reached through `activation.component-staged`, a live Phase 21/22
  family, so without it the first Activation that stages a profile carrying
  `decision.record-reference`
  is rejected with `invalid closed enum` before this task can append a single
  Event, with its cause outside this task's scope.
- `internal/issues/transition_rules.go` - the `linked_to_decision` evaluator.
- `internal/journal/project_effects.go`, `current_event_catalog.go` - the
  `decision.referenced` fence key, and this task's two amendments to existing
  families. Both amended families are in **`retainedCurrentRules`**, not
  `currentEventRules` (`current_event_catalog.go:82-83`; `schema.md` 6.2 as
  corrected by DEC-232). Add the `envelope_protected_decisions/insert` stage to
  the retained `envelope.created` rule and the Stop stages to the retained
  `envelope.activated` rule. Moving either family into `currentEventRules` would
  route a live family onto the strict path and fail immediately on `unknown
  table`; the retained path needs the stages only because task 1 registers
  `envelope_protected_decisions` in `currentEventTables`.
- `internal/workflow/stop.go` - the closed Stop registries (R-14). Add
  `protected-decision-required` to `stopCodes` with its exact
  `validStopDetails` key set, `decision-satisfied` to `conditionTypes`, and the
  `decision_class` field to the `Condition` struct together with its
  exactly-one-non-null-field rule in `validateCondition`, all from
  `public-interface.md` section 2. **Tag the new field `omitempty`, and add
  `omitempty` to nothing else** (DEC-235, resolving ESC-23-001). `budget_name`
  and `capability_id` keep present-and-null treatment: putting `omitempty` on
  either changes every stored `attention_items.semantic_key` and
  `idempotency_results.request_digest` in existence - measured - which is the
  exact harm the decision exists to prevent. A `decision-satisfied` Condition
  therefore marshals five fields, a Phase 22 Condition still four, byte-for-byte
  as today.
- `internal/journal/current_event_json.go` and
  `internal/journal/current_event_result_projection.go` - the typed-result
  registration for `decision.referenced` (R-16, DEC-231). Deleting the fence key
  routes this family onto the journal's strict result path, so both are required
  here or the first successful `decision record-reference` is rejected at append
  time. Add the `decision.referenced` entry in the `fields` grammar map, derived
  from the `decision.record-reference` result row in `public-interface.md`
  section 4, and the matching `currentExpectedResult` branch. Also, in
  `current_event_json.go` and distinct from the `fields` map, register
  `decision_refs.source_json` as `Source` and
  `decision_refs.authority_capabilities_json` as `[]ID` in the `effect.Table +
  "." + column` switch at `:14`, which otherwise returns `unregistered current
  JSON column decision_refs.source_json` on the first Decision Reference
  (DEC-232).
- `internal/journal/current_event_value_types.go` - **new to this task's scope,
  and this task creates the pattern the next two follow** (DEC-235, inventory
  row 10). ESC-23-001 is resolved: `Condition` is widened by **per-type shapes
  keyed on `type`**, not by widening one shape. Replace the single
  `currentValueShapes["Condition"]` entry with a per-type map in
  `current_event_value_catalog.go` and add the selecting `Condition` branch to
  `validateCurrentValue` here. Three hard requirements. **(a)** The six Phase 22
  type rows keep their exact current four-field shape -
  `type:=<type> subject:Subject budget_name:?ID capability_id:?ID` - so
  historical Events replay unchanged and no stored semantic key moves; this was
  measured byte-identical and it is the acceptance bar. **(b)** The branch reads
  the discriminant **before** `validateCurrentShape` has proved it is a string,
  so it must use a **checked** type assertion and fall through to the ordinary
  rejection: `object["type"].(string)` panics `interface conversion: interface
  {} is int, not string`, and this code validates untrusted journal and archive
  bytes during replay and `import`. **(c)** The branch must **end by calling
  the relation**: `return validateCurrentValueRelations("Condition", object)`,
  exactly as the generic `currentValueShapes` path it replaces does at
  `current_event_value_types.go:39-44`. Removing the `"Condition"` key from
  `currentValueShapes` deletes that path's only call for this kind, so a branch
  that validates the shape and returns leaves the Phase 22 `budget_name` /
  `capability_id` binding unenforced. Measured (QA R-24): such a branch accepts
  a `work-ended` Condition carrying `budget_name:"tokens"`, **appends and
  projects it** on the strict path, and still passes the entire ordinary suite
  `exit 0` - `condition nullable binding mismatch` is asserted by no test in the
  repository, tagged or not, so nothing else detects the loss. "Do not touch
  `case "Condition"`" below means do not change the relation; it does not mean
  stop reaching it. Your own type's row is
  `type:=decision-satisfied subject:Subject budget_name:?ID capability_id:?ID
  decision_class:ID`.
- `internal/journal/current_event_value_relations.go` - **do not touch
  `case "Condition"`** (DEC-235, restating inventory row 8). Under the accepted
  encoding the new fields' exclusivity is structural, enforced by the shapes
  above; adding a third boolean term here has misread the decision. The branch
  stays as it is because it still carries the Phase 22 `budget_name` /
  `capability_id` binding, and the accepted rows deliberately keep both fields
  `?ID` on **every** type, so this relation is their only enforcer - which is
  exactly why requirement (c) above exists (I-40 corrects the earlier reason,
  "the shapes cannot express it"; the grammar's `Null` token could express it,
  but the accepted rows do not use it). This task still edits this file, but
  only for the capability allowlist below.
- `internal/journal/current_event_value_catalog.go` - also add the `DecisionRef`
  and `ProtectedDecision` entries to `currentValueShapes` from
  `public-interface.md` section 3 (R-16, DEC-231). This task owns
  `ProtectedDecision` because it is the first to project one; task 10 reuses it
  unchanged. The same file carries the matching
  `currentStopDetailShapes` entry, and the per-type `Condition` shape map
  described above, into which this task writes the six retained Phase 22 rows
  unchanged plus its own `decision-satisfied` row (R-14, DEC-235).
- `internal/export/archive_semantics_catalog.go` - mirror this task's own
  catalog amendments (`envelope.created` gaining
  `envelope_protected_decisions/insert`, `envelope.activated` gaining the
  optional Stop rows) into `semanticEventRules` (R-14). Without the mirror,
  `export validate` and `import` reject the journal of any workspace that
  recorded those effects, and the divergence stays silent until task 16.
- `internal/cli/codes.go` - `protected-decision-required` and
  `source-digest-mismatch` only. Do **not** add `decision-not-found`: it is
  removed from the accepted set under DEC-230 because no leaf takes a Decision
  Reference ID, so nothing could return it and VAL-056 could not reproduce it.
- `cmd/fid/decision.go` - new leaf wiring for `record-reference` and `list`.
- `cmd/fid/envelope.go`, `cmd/fid/issue_mutation_response.go`, and matching
  tests - map the domain's protected-Decision denial to the accepted
  `protected-decision-required` code, exit status, and immutable Decision
  details for `envelope.complete` and `issue.transition-terminal-success`
  (DEC-248). Preserve the required precedence when an Obligation is also
  unsatisfied; other completion and transition error mappings stay unchanged.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change.
- `project-docs/docs/phases/phase-23/tasks/task-9-done.md`.

### Files To Read For Context

- `contracts.md` sections 2, 9, 10; `public-interface.md` sections 2, 3, 4, 4.1,
  5; `schema.md` sections 1, 2, 6.2.
- `task-8-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Human Outcome recording, authority, and provenance.
- Attention resolution through a Human Outcome, which is task 11.

## Acceptance Criteria

- [ ] No prose body column exists; a matching digest records and a mismatch
  rejects with no Event. (VAL-038)
- [ ] Only the declared affected action is blocked, with the exact code.
  (VAL-039)
- [ ] Only the both-approved state satisfies; neither fact alone and no
  `rejected` or `deferred` reference does. (VAL-040)
- [ ] Two unsatisfied protected Decisions raise exactly two Stops, and
  satisfying one resolves exactly that Stop and Attention in one Event.
  (VAL-041)
- [ ] `linked_to_decision` is listed as implemented and is satisfied only with a
  qualifying reference. (VAL-042)
- [ ] The per-type dispatch accepts the exact pinned canonical pre-change bytes
  of **all six** retained Phase 22 Condition types, not left to corpus
  traffic. A table over `executor-terminated`, `run-recovered`,
  `activation-recovered`, `capability-available`, `budget-observed`, and
  `work-ended`: for each, one canonical four-field Condition from a named
  pre-change reference is accepted unchanged at the journal layer and one
  mutated Condition of that same type is rejected.
  **Canonical means type-appropriate**, not all-null: `budget-observed` carries
  a non-null `budget_name` and `capability-available` a non-null
  `capability_id`, because those two types are the ones the Phase 22 relation
  binds. An all-null form of either is rejected `condition nullable binding
  mismatch` - measured (QA I-44) - so it is a VAL-081 case, not an acceptance
  here. Getting this wrong costs a red-test cycle on exactly the two
  zero-coverage types this table exists for.
  The per-type table is required because
  `run-recovered` and `capability-available` reach journal Condition validation
  **zero** times in the whole ordinary suite and are raised nowhere in the
  tagged corpus - measured (QA R-25): deleting both rows outright leaves every
  gate this phase requires green, so a dropped, misspelled, or mistyped row
  ships a binary that rejects those Conditions at runtime with `Condition: type:
  invalid closed enum` and no gate notices. This is the acceptance bar of
  DEC-235; a failure here means the per-type rows were not transcribed exactly.
  (VAL-082)
- [ ] A recorded pre-change journal carrying a `budget-observed` Stop replays
  clean under the new binary with an unchanged source-journal digest. Record
  the reference commit, fixture path/digest, replay result and post-replay
  source digest. (VAL-083; DEC-240, split from VAL-082)
- [ ] The `attentionSemanticKey` for the same fixed subject, code and retained
  Condition is unchanged across the named pre-change reference and this task.
  Record identical inputs and both key strings. (VAL-084; DEC-240, split from
  VAL-082)
- [ ] Exclusivity is enforced by the shape, asserted at the journal layer and
  not through the raise: a `decision-satisfied` Condition carrying
  `outcome_class` instead of `decision_class` rejects, one with its own field
  null rejects, one carrying two new fields rejects, and one carrying only its
  own field is accepted. (Evidence toward VAL-080, which task 11 completes.)
- [ ] The Phase 22 nullable binding still fires **from the new dispatch**: a
  `work-ended` Condition carrying a non-null `budget_name`, a `budget-observed`
  Condition carrying `budget_name:null`, and a `decision-satisfied` Condition
  carrying a stray `budget_name` are each rejected `condition nullable binding
  mismatch` on the strict path, and the projection is unchanged. This is the
  executable proof of requirement (c) - that the branch still reaches
  `validateCurrentValueRelations` - and it is the only assertion of that message
  anywhere in the repository. (VAL-081; DEC-235 inventory row 8, QA R-24.)
- [ ] A Condition whose `type` is not a string is **rejected, not a panic**,
  proving the checked assertion. (No new VAL ID, declined with reason under
  DEC-237: a panic fails this task's own gate loudly, so unlike VAL-081 and
  VAL-082 the property already has a detector; DEC-235.)
- [ ] The `currentStopDetailShapes` and per-type `Condition` entries this task
  adds are proved by a **direct assertion inside this task**, not by this task's
  own retained-family raise tests. `protected-decision-required` is raised at
  `envelope.activated`, a retained family that returns at
  `current_event_validation.go:64` before any value validation, so that raise
  path never exercises these entries. **Follow
  `TestCurrentClosedValuesRejectIncompleteAndInvalidFacts`**
  (`internal/journal/current_event_closed_values_test.go:76`): it already
  asserts on family `stop.raised` for `stops.resolution_condition_json`, and on
  `stops.details_json` with a `code` in its `values` map, needs no database, and
  rejects `unknown Stop code` when the detail entry is missing. Do **not**
  follow `TestCurrentEventRejectsConditionalBindingCorruption` for the details
  half: it mutates `executionBlockedProjectionEvent`, and `run.execution-blocked`
  is one of the five families in `validateCurrentBudgetFacts`
  (`current_event_budget_values.go:34`, inventory row 9), so a changed Stop code
  there is rejected `current budget boundary: unrelated Stop` before any detail
  shape is consulted - measured, QA R-27. A **real `stop.raise` on the strict
  `stop.raised` family is acceptable evidence too**: `validateStopRaise`
  (`stop.go:363-374`) checks the code and the condition independently. What is
  insufficient is a raise on this task's **own** retained family. If the chosen
  proof appends into the projection it also depends on migration 0017's widened
  `stops.code` and `attention_items.code` CHECK from task 1 (`schema.md` section
  1); the `validateCurrentJSON` route touches no database and has no such
  dependency. Omitting the entries would pass here and fail at task 11 or tasks
  15/16 with the cause outside the failing task. (QA R-22/R-27, DEC-235.)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestProtectedDecisionBlocksOnlyDeclaredAffectedActions` fails against the
absent protected-Decision evaluation.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 3 and section 5: add
  `fid decision list` and the Phase 23 closed codes owned by this task.
- `project-docs/docs/phases/phase-23/tasks/task-9-done.md` in the full format,
  including the fence-key deletion and its parity proof.
## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -timeout 30m ./cmd/fid/integration/...
```

Two gates run for minutes and are explicitly bounded (DEC-227): the `-timeout
30m` above replaces Go's implicit 10-minute per-binary default, which the tagged
suite already approaches at 560 s. Record in the completion report the wall time
`go test` reports for `cmd/fid/integration`, `cmd/fid`, and `internal/workflow`,
so a regression is attributed to the task that caused it. Budgets: tagged
integration under 900 s and no ordinary package over 600 s (VAL-079). A breach is
a bounded speed repair owned by this task when this task's own new journeys
caused it; otherwise stop and escalate. Never record a timed-out gate as `PASS`.

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. No exception to the full task
gate. Focused red/green selections never replace these commands. Report each
attempt with the command's own output token and observed termination. An
interrupted or hung command is never `PASS`; preserve all failed attempts and
rerun the byte-identical command after an authorized causal repair. Report lint
is a Runner control after the report commit, not a Developer gate.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- This task belongs to QA Review window 3.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-9`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
