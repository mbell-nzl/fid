# Phase 23 - Task 10: Provenanced Human Outcomes

## Objective

Record authorized Human judgment with active Kit Activation provenance, allow it
to satisfy only a matching `human-outcome` Obligation, and never present it as
deterministic proof.

## Validation Contract Coverage

- VAL-043: Recording requires a human Actor type in a finalized Kit Activation
  profile with the capability, and standalone Issue-inventory attributed
  authority rejects (PRD-AUTH-007).
- VAL-044: A Human Outcome satisfies only a matching class with result
  `approved`.
- VAL-025: complete Task 6's staged nine-cell compatibility matrix with a real
  recorded approved Human Outcome linked durably to its matching human-outcome
  Obligation; re-prove the Human row's two mismatch denials with that real
  producer. Final PASS requires Task 6's other cells too (DEC-243).
- VAL-045: the Human Outcome Evidence record and its evaluation carry
  `producer:"human"`, `trust_source:"activated-human"`, and
  `artifact_availability:"not-applicable"`. The Human-text half of
  `contracts.md` section 14 belongs to task 14 (VAL-058/VAL-060); do not pin text
  task 14 replaces.
- VAL-046: A `human-outcome-required` Stop is raised by each of its four named
  Events exactly when the Human Obligation becomes the only unsatisfied one, at
  most once while unresolved, and never for routine status (DEC-224).
- VAL-009: the completing half - a guarded action is permitted only once the
  gate-evidence, review, and human-outcome Obligations are all satisfied. This
  is the first task where that state is reachable.
- VAL-010: re-verified here with linked review and Human Outcome evidence
  present, not only with evidence absent.
- VAL-003: task 1's fence assertion for `human-outcome.recorded` converges
  here without losing its parity subject.
- VAL-080: this task's `human-outcome-recorded` row carries `outcome_class` and
  no other discriminator, enforced by the shape at the journal layer. Task 11
  completes the four-type matrix.
- VAL-081: re-confirm, after adding this task's row, that the Phase 22 nullable
  binding still fires from task 9's dispatch. Task 9 owns the assertion; a
  failure here is a task 9 defect (DEC-237).


## Binding Constraints

- `contracts.md` section 10 and `public-interface.md` sections 3 and 4 are the
  exact authority; `schema.md` section 2 fixes `human_outcomes`.
- Authority is established by active-version Activation provenance and the
  declared capability, never by a role string or a free-form `--by human` flag.
  A non-human actor type, an unfinalized Activation, or a missing capability all
  reject `authority-missing` with no Event.
- Standalone Issue inventory attribution can never reach this leaf
  (PRD-AUTH-007). Prove the negative case explicitly.
- A Human Outcome may satisfy only a `human-outcome` Obligation whose
  `outcome_class` equals its class, and only with result `approved`. Any other
  link rejects `outcome-not-permitted` naming the Obligation type.
- DEC-243: Task 6's Human-row evidence may use a production classifier or seeded
  fixture because it cannot record Human Outcomes. This task closes that
  evidence boundary with actual `human-outcome.record` and `evidence.link`
  requests. Reuse Task 6's `outcome-not-permitted` registration, preserve all
  nine pairings and the DEC-226 precedence, and reconcile the per-cell evidence
  against `task-6-done.md`; seeded/classifier acceptance alone never closes the
  positive Human Outcome link.
- The Evidence record for a Human Outcome retains exactly `producer:"human"`,
  `trust_source:"activated-human"`, and
  `artifact_availability:"not-applicable"`. No field, reason, or text may
  present it as deterministic or fid-executed proof (PRD-HUM-001).
- Raise a `human-outcome-required` Stop exactly when the Human Obligation is the
  only unsatisfied Obligation. Routine status raises no Attention (PRD-HUM-002).
- The `artifact.availability-observed` raising sequence has a proven fixture
  technique (I-24, reproduced live at QA readiness): remove the retained
  content-addressed blob under `.fid/artifacts/<digest>` and observe `missing`,
  then restore it and observe `available`. `internal/evidence/artifact.go`
  derives availability by digest-verifying that blob and writes one
  `referenced_artifacts/update` effect with a precondition on the prior value,
  so `missing -> available` is exactly one qualifying Event and the Obligation
  transitions unsatisfied -> satisfied inside it.
- DEC-192: delete exactly the `"human-outcome.recorded"` key and prove its
  live-apply and replay parity; converge task 1's fence assertion for that key.
- Add the `human-outcome.record` grant action exactly here. Replace the seeded
  Outcome fixtures used by tasks 7 and 9 with real recorded Outcomes and record
  that convergence with provenance.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/review/` - Human Outcome rules, provenance, effects.
- `internal/evidence/evidence.go`, `policy.go` - the Human Outcome evidence type.
- `internal/workflow/authority.go`, `activation.go`, `evaluation.go`,
  `resolve.go` - the new grant action belongs to the closed
  `sortedUniqueCapabilities` allowlist in `activation.go` (R-08, DEC-167/173).
- `internal/journal/current_event_value_relations.go` - the **second** copy of
  that same closed allowlist (R-19, DEC-233). `validateCurrentValueRelations`'s
  `case "ActorProfile"` holds a `const actions` string with the identical 23
  names at `:63`, and nothing cross-checks the two. Add
  `human-outcome.record`
  there in the same change that adds it to `sortedUniqueCapabilities`, keeping
  both copies sorted and identical. This is not on any Phase 23 family's strict
  path: it is reached through `activation.component-staged`, a live Phase 21/22
  family, so without it the first Activation that stages a profile carrying
  `human-outcome.record`
  is rejected with `invalid closed enum` before this task can append a single
  Event, with its cause outside this task's scope.
- `internal/workflow/` readiness-Stop raising shared by the four producers, and
  its call sites for `envelope.activated`, `evidence.linked`, `review.recorded`,
  and `artifact.availability-observed` (DEC-224).
- `internal/journal/project_effects.go`, `current_event_catalog.go` - the fence
  key, plus the optional `stops/insert`, `attention_items/insert`, and
  `attention_stop_links/insert` effect rows added to `evidence.linked` and
  `artifact.availability-observed` under `schema.md` section 6.2 (DEC-224).
- `internal/workflow/stop.go` - the closed Stop registries (R-14). Add
  `human-outcome-required` to `stopCodes` with its exact `validStopDetails` key
  set, `human-outcome-recorded` to `conditionTypes`, and the `outcome_class`
  field to the `Condition` struct together with its exactly-one-non-null-field
  rule in `validateCondition`, all from `public-interface.md` section 2.
  **Tag the new field `omitempty`, and add `omitempty` to nothing else**
  (DEC-235). `budget_name` and `capability_id` keep present-and-null treatment;
  tagging either changes every stored `semantic_key` and `request_digest` in
  existence. Task 9 set this pattern - follow it exactly.
- `internal/journal/current_event_json.go` and
  `internal/journal/current_event_result_projection.go` - the typed-result
  registration for `human-outcome.recorded` (R-16, DEC-231). Deleting the fence
  key routes this family onto the journal's strict result path, so both are
  required here or the first successful `human-outcome record` is rejected at
  append time. Add the `human-outcome.recorded` entry in the `fields` grammar
  map, derived from the `human-outcome.record` result row in
  `public-interface.md` section 4, and the matching `currentExpectedResult`
  branch. Also, in `current_event_json.go` and distinct from the `fields` map,
  register `human_outcomes.references_json` as `[]Source` in the `effect.Table +
  "." + column` switch at `:14`, which otherwise returns `unregistered current
  JSON column human_outcomes.references_json` on the first Outcome (DEC-232).
- `internal/journal/current_event_value_relations.go` - **do not touch
  `case "Condition"`** (DEC-235, restating inventory row 8). ESC-23-001 is
  resolved: the new fields' exclusivity is structural, enforced by the per-type
  shapes task 9 created, so there is no relation work for this field. The branch
  stays as it is because it still carries the Phase 22 `budget_name` /
  `capability_id` binding, and the accepted rows keep both fields `?ID` on every
  type, so this relation is their only enforcer (I-40). **Confirm in the
  completion report that the relation still fires** from task 9's dispatch: a
  `human-outcome-recorded` Condition carrying a stray `budget_name` is rejected
  `condition nullable binding mismatch`. No test in the repository asserts that
  message, so a task 9 branch that returned after shape validation would have
  deleted the binding silently (QA R-24). If it does not fire, that is a task 9
  defect - stop and escalate rather than adding the call here. This task still
  edits this file, but only for the capability allowlist above.
- `internal/journal/current_event_value_catalog.go` - also add the
  `HumanOutcome` entry to `currentValueShapes` from `public-interface.md`
  section 3 (R-16, DEC-231). Reuse task 9's `ProtectedDecision` unchanged; if it
  is missing or wrong, that is a task 9 defect. `HumanOutcome.subject` resolves
  through task 4's `JudgmentSubject`, also reused unchanged; do **not** edit
  `currentValueShapes["Subject"]`, which Phase 23 leaves as it is (DEC-232,
  correcting DEC-231). The same file carries the
  matching
  `currentStopDetailShapes` entry, and the per-type `Condition` shape map task 9
  created, which gains exactly one row here:
  `type:=human-outcome-recorded subject:Subject budget_name:?ID
  capability_id:?ID outcome_class:ID` (R-14, DEC-235). Add a row; do not widen
  task 9's rows, and do not touch the six Phase 22 rows.
- `internal/export/archive_semantics_catalog.go` - mirror this task's
  `retainedCurrentRules` amendments for `evidence.linked` and
  `artifact.availability-observed` into `semanticEventRules` (R-14). Without
  the mirror, `export validate` and `import` reject the journal of any
  workspace whose readiness Stop was raised by either Event.
- `internal/cli/codes.go` - reuse Task 6's `outcome-not-permitted` registration
  (DEC-243); add `outcome-not-found`,
  `human-outcome-required`.
- `cmd/fid/human_outcome.go` - new leaf wiring for `record` and `list`.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change, including the task 7 and task 9 seeded fixtures.
- `project-docs/docs/phases/phase-23/tasks/task-10-done.md`.

### Files To Read For Context

- `contracts.md` sections 7, 9, 10, 13; `public-interface.md` sections 2, 3, 4,
  4.1, 5; `schema.md` sections 1, 2, 6.
- `task-9-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Attention resolution through a Human Outcome, which is task 11.
- The `human_outcome_resolution` capability limit, which is task 17.

## Acceptance Criteria

- [ ] Agent actor type, unfinalized Activation, and missing capability each
  reject `authority-missing` with no Event; the valid case records with its
  provenance bound. (VAL-043)
- [ ] Only the matching class with `approved` satisfies; every other link
  attempt rejects `outcome-not-permitted`. (VAL-044)
- [ ] A real recorded approved Human Outcome links to the matching
  human-outcome Obligation with its immutable Evidence/Event and readback;
  attempts to link it to gate-evidence and review reject `outcome-not-permitted`
  with exact closed details and no Event/projection change. Combine these
  observed cells with Task 6's receipt/review positives and other four mismatch
  results, accounting for all nine pairings before claiming VAL-025 PASS.
  Existing Task 6 tests remain exercised by this task's full gate. (VAL-025;
  DEC-243; complements VAL-044)
- [ ] The Evidence record and evaluation carry exactly the declared triple and
  no field or text claims deterministic proof. (VAL-045)
- [ ] The `human-outcome-required` Stop appears only when the Human Obligation
  is the single unsatisfied one, and never for routine status. (VAL-046)
- [ ] Exclusivity is enforced by the shape, asserted at the journal layer: a
  `human-outcome-recorded` Condition carrying `decision_class` instead of
  `outcome_class` rejects, one with its own field null rejects, one carrying two
  new fields rejects, and one carrying only its own field is accepted.
  (Evidence toward VAL-080, which task 11 completes.)
- [ ] The `currentStopDetailShapes` and per-type `Condition` entries this task
  adds are proved by a **direct assertion inside this task**, not by this task's
  own retained-family raise tests. Three of DEC-224's four raising families
  (`envelope.activated`, `evidence.linked`, `artifact.availability-observed`)
  are retained and return at `current_event_validation.go:64` before any value
  validation, so the raises this task owns do not exercise these entries.
  **Follow `TestCurrentClosedValuesRejectIncompleteAndInvalidFacts`**
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
  insufficient is a raise on one of this task's **own** retained families. If
  the chosen proof appends into the projection it also depends on migration
  0017's widened `stops.code` and `attention_items.code` CHECK from task 1
  (`schema.md` section 1); the `validateCurrentJSON` route touches no database
  and has no such dependency. Omitting the entries would pass here and fail at
  task 11 or tasks 15/16 with the cause outside the failing task. (QA
  R-22/R-27, DEC-235.)
- [ ] The Phase 22 nullable binding still fires from task 9's dispatch after
  this task's row is added: a `human-outcome-recorded` Condition carrying a
  stray `budget_name` is rejected `condition nullable binding mismatch`, and the
  completion report records the observed message. A failure here is a task 9
  defect, not a local fix. (VAL-081; QA R-24, I-45)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row, including the seeded task 7 and task 9
  fixtures, records its provenance and executing command.

## First Failing Test

`TestHumanOutcomeRequiresFinalizedActivationProvenance` fails against the
absent Human Outcome writer.

## Documentation Updates

- `project-docs/docs/phases/phase-23/tasks/task-10-done.md` in the full format,
  including the fence-key deletion and its parity proof.
- Flag spec impact; only Tech Lead accepts a contract or schema change.
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

- This task belongs to QA Review window 4.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-10`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
