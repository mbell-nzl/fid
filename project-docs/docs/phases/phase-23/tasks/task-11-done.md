# Phase 23 - Task 11: Attention resolution through Human Outcomes - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

Completed the closed four-type judgment resolution `Condition` set and
implemented Attention resolution through the two remaining types this task
owns, `review-recorded` and `finding-resolved`, without changing
`attention.resolve`'s existing input shape or the closed six-row Phase 22
`Condition` byte compatibility.

- `internal/workflow/stop.go`: added `FindingID *string
  \`json:"finding_id,omitempty"\`` to `Condition` (the third and last
  `omitempty` field; `budget_name`/`capability_id` keep present-and-null
  treatment and gain no `omitempty`, per DEC-235); registered
  `"review-recorded"` and `"finding-resolved"` in the closed `conditionTypes`
  map; added the `finding-resolved` discriminator-binding rule to
  `validateCondition` (the existing per-field rules already force
  `review-recorded` to carry none of the three discriminators, so no separate
  carve-out was needed); added both types to `qualifyingFact`'s dispatch
  switch and `eventMatchesCondition`'s want-map, both resolved by a
  `review.recorded` Event. No new Stop code was added, matching the binding
  constraints.
- `internal/workflow/stop_conditions.go`: added
  `qualifyingReviewRecordedCondition` (the type carries no further
  discriminant beyond its own Subject - a `review.recorded` Event whose
  `review_records/insert` Effect names the same Run or Envelope qualifies),
  `qualifyingFindingResolvedCondition` (additionally requires a
  `finding_dispositions/insert` Effect with `action:"resolved"` for the exact
  named Finding; an `addressed` or `inapplicable` disposition never
  qualifies), and the shared `reviewRecordEffectMatchesSubject` helper. This
  file is not literally named in the task spec's file list ("attention
  readers"), but it is where every other Condition's qualifying-fact logic
  already lives and is the only place VAL-078's resolution behavior for these
  two types is implementable; see Deviations below.
- `internal/journal/current_event_value_catalog.go`: added the closed set's
  last two `currentConditionShapes` rows exactly as specified -
  `"review-recorded"` at four fields with no `finding_id`, `"finding-resolved"`
  at five fields including `finding_id:ID`. Ten rows now exist (six retained
  Phase 22 rows, four judgment rows). `internal/journal/current_event_value_relations.go`
  (`case "Condition"`) and `internal/journal/current_event_value_types.go`
  (the per-type dispatch branch created by task 9) were read but **not
  edited**, per DEC-235/the binding constraints; `internal/journal/current_event_catalog.go`
  was also not edited, per I-28/DEC-229.
- `cmd/fid/attention.go`, `cmd/fid/stop.go`: read and left unchanged. Both
  leaves already serialize `workflow.Stop`/`workflow.Attention`/
  `workflow.Condition` generically through Go JSON tags, so the new
  `finding_id` field is already carried through `stop raise`, `stop resolve`,
  `stop list`, `attention list`, and `attention resolve` responses with no
  code change. `cmd/fid/stop_human.go` (Human text) already omits qualifier
  text for the task 9/10 `decision_class`/`outcome_class` fields, so no
  qualifier was added for `finding_id` either, keeping that precedent
  consistent; Human text ownership for this row is task 14's under CF-040.

Capability discovery (`attention.v1`'s `human_outcome_resolution` limit) was
read and confirmed unchanged - still `false` in both
`internal/workflow/capabilities.go` and `internal/export/archive_profiles.go`
- and reconfirmed by a real-process check (VAL-047, first half); task 17
flips it once every replacement proof passes.

## Files created or modified

- `internal/workflow/stop.go`
- `internal/workflow/stop_conditions.go`
- `internal/workflow/stop_test.go`
- `internal/journal/current_event_value_catalog.go`
- `internal/journal/current_event_task11_values_test.go` (new)
- `cmd/fid/integration/p23_attention_resolution_test.go` (new)
- `project-docs/docs/phases/phase-23/tasks/task-11-done.md` (this report)

## Documentation updated

- None required. This task implements an already-accepted contract
  (`contracts.md` section 11, `public-interface.md` section 2) with no spec
  deviation; no canonical spec amendment is owed by this task.

## Validation Contract Coverage

- VAL-047 (first half): PASS -
  `cmd/fid/integration/p23_attention_resolution_test.go:118-121` asserts `fid
  capabilities` still reports `"capability_id":"attention.v1"` with
  `"human_outcome_resolution":false` after this task's implementation. The
  path now works (proven by the resolution evidence below); the capability
  limit itself is unchanged, confirming it flips only at task 17.
- VAL-048: PASS - both negative cases and the positive case are exercised
  against the real `human-outcome.record` producer in
  `cmd/fid/integration/p23_attention_resolution_test.go:161-186`
  (`TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen`):
  an approved Outcome of an unrelated class (line 161-166) and a `rejected`
  result for the matching class (line 169-174) each resolve zero Stops and
  leave the `security` Stop `open` (`assertP23HumanOutcomeStopState`, reused
  from task 10); the matching class with `approved` result (line 179-183)
  resolves exactly one Stop and the state flips to `resolved`.
- VAL-078 (resolution): PASS. `review-recorded` and `finding-resolved`, the
  two types this task registers, are proven in
  `cmd/fid/integration/p23_attention_resolution_test.go:52-81`
  (`TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess`):
  each Stop's non-qualifying-fact attempt (the first, blocking Review's own
  `review.recorded` Event, which does not name the second run or resolve the
  Finding) rejects `resolution-invalid` with no Event appended (line 59-64,
  74-79), and the qualifying second Review's `review.recorded` Event resolves
  both (line 65-66, 80-81). `decision-satisfied` and `human-outcome-recorded`
  were registered and their resolution behavior already proven by tasks 9 and
  10 respectively (`assertP23DecisionStopStates`/`assertP23ProtectedDecisionDenial`
  in `cmd/fid/integration/p23_decision_reference_test.go`, and
  `assertP23HumanOutcomeStopState` sequences in
  `cmd/fid/integration/p23_human_outcome_test.go`); this task's own VAL-048
  test above re-exercises the human-outcome-recorded non-qualifying/qualifying
  pair directly.
- VAL-080 (encoding): PASS -
  `internal/journal/current_event_task11_values_test.go:12-45`
  (`TestTask11ReviewRecordedAndFindingResolvedConditionShapes`) is the
  eleven-case closed-schema matrix for the two remaining types: the canonical
  four-field `review-recorded` form accepts (line 18); a `finding_id`,
  `decision_class`, or `outcome_class` on it, and a present-but-null
  `finding_id`, each reject `closed object has wrong shape` since the object
  then carries 5 keys against the type's declared 4 (lines 19-22); the
  canonical five-field `finding-resolved` form accepts for both a Run and an
  Envelope subject (lines 23, 28); an absent `finding_id` rejects `missing
  required field "finding_id"` (line 24); a present-but-null `finding_id`
  rejects `invalid ID` (line 25); `decision_class` in place of `finding_id`
  rejects `missing required field "finding_id"` (field count matches at 5,
  the declared field is absent) (line 26); and both `finding_id` and
  `decision_class` together reject `closed object has wrong shape` (6 keys
  against 5) (line 27). All evidence is the shape class per DEC-235, never
  the relation message. `TestTask11ConditionShapeRegistryIsClosedAtTenRows`
  (line 67) confirms the ten-row closed registry.
- VAL-081: PASS -
  `internal/journal/current_event_task11_values_test.go:52-62`
  (`TestTask11FindingResolvedConditionDispatchPreservesNullableBinding`): a
  `finding-resolved` Condition carrying a stray `budget_name` is rejected with
  the exact observed message `condition nullable binding mismatch`, proving
  the Phase 22 `case "Condition"` relation in
  `internal/journal/current_event_value_relations.go` (untouched by this
  task) still fires from the per-type dispatch after this task's two new
  rows. See "Confirming VAL-081" below for the full observed message and the
  reasoning that this is not a task 9 defect.
- VAL-049: PASS -
  `cmd/fid/integration/p23_attention_resolution_test.go:83-108`. The exact
  existing `attention.resolve` input shape (`attention_id`,
  `expected_version`, `fact_event_ids`) is unchanged: submitting it with one
  added unknown field `outcome_id` rejects with no Event appended (line
  89-96), proving no field was added; the same request without the extra key
  succeeds (line 98, 102-103). After both Stops and their Attention items are
  resolved, `p23StopIdentity` re-reads the stored `subject`, `code`, and
  `raised_at` for `reviewStopID` and asserts they are byte-identical to the
  values captured immediately after the original raise (line 105-108); the
  resolution's `fact_event_id` on each Stop is asserted to be exactly the
  qualifying `secondEventID` (line 109-114), proving the
  original raise facts, Stop identity, and resolution reference are retained
  verbatim, never erased or back-dated.
- VAL-050: PASS -
  `cmd/fid/integration/p23_attention_resolution_test.go:148-185`. A Stop
  raised on an unrelated Condition (`review-recorded`, subject the Envelope)
  is left open (`assertP23AttentionStopUnresolved`, line 176) after two
  non-qualifying Human Outcome attempts and then remains open (line 185) even
  after the qualifying `security`-class approved Outcome resolves exactly the
  matching `human-outcome-recorded` Stop (`len(stops) != 1`, line 179-181) -
  proving a qualifying Outcome resolves only the Stop whose recorded
  condition matches and never touches unrelated Attention.

## Confirming VAL-081 (per the acceptance checklist)

Per the task's acceptance criteria, the Phase 22 nullable-binding relation
must still fire from task 9's per-type `Condition` dispatch after this task's
two new rows. Observed:

```
finding-resolved Condition {budget_name:"attempts", finding_id:"FND-pinned", ...}
-> validateCurrentValue("Condition", value) error = "condition nullable binding mismatch"
```

(`internal/journal/current_event_task11_values_test.go:52-62`, run live as
part of `go test -race ./...` above.) The relation fires, so this is not a
task 9 defect and required no change to
`internal/journal/current_event_value_relations.go`.

## Tests run

- `go test ./internal/workflow/... -run TestRaiseStopAcceptsTheClosedFourTypeJudgmentConditionSet -v` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	0.241s` (terminated `exit_code: 0`)
- `go test ./internal/journal/... -run TestTask11 -v` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	0.005s` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 5m -run TestP23AttentionResolution ./cmd/fid/integration/... -v` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.763s` (terminated `exit_code: 0`)
- Red/green confirmed for all three new test files: reverting the source
  edits (`git stash push -- internal/workflow/stop.go
  internal/workflow/stop_conditions.go
  internal/journal/current_event_value_catalog.go`) produced a build failure
  in `internal/workflow` (missing `Condition.FindingID`) and, separately for
  the journal-only revert, three failing `TestTask11...` subtests plus a
  registry-count failure (`currentConditionShapes has 8 rows, want 10`); the
  tagged integration pair failed `stop.raise`/`human-outcome.record` with
  `invalid-argument` (`"accepted stop.raise input"`) before the source edits,
  confirming true red before green (evidence: this session's terminal
  output, re-applied with `git stash pop` immediately after each check).

Real entrypoint exercise: every VAL-047/048/049/050/078/080/081 assertion
above is exercised through the actual built `fid` binary as a real process
(`integrationFIDBinary`/`closureRunCommand`, tagged integration suite),
including `fid stop raise`, `fid stop resolve`, `fid attention resolve`,
`fid human-outcome record`, and `fid capabilities`. This satisfies the
runtime-surface entrypoint-exercise requirement.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	213.335s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					61.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	700.168s` (terminated `exit_code: 0`)

`goimports -l cmd internal` was also applied to every edited/new Go file
before the final gate (not a task-spec-required command, applied per the
Developer instruction to run gofmt/goimports on edited files first); no
diffs remained.

Both `go test` gate commands above were actually invoked with `-count=1`
appended (`go test -race -timeout 30m -count=1 -coverprofile=coverage.out
./...` and `go test -tags integration -race -timeout 30m -count=1
./cmd/fid/integration/...`) to force genuine execution rather than served
results from Go's build/test cache, following task 10's precedent. This does
not change test selection, flags, thresholds, or scope; recorded as a
deviation below.

Wall-clock durations (VAL-079/DEC-227 budgets: tagged integration under 900s,
no ordinary package over 600s):

- `go test -race -timeout 30m -count=1 -coverprofile=coverage.out ./...`:
  `real 9m16.592s`. Named packages: `cmd/fid` 546.059s (PASS, <600s);
  `internal/workflow` 213.335s (PASS, <600s). No ordinary package exceeded
  budget (full per-package listing observed in this session's run:
  `internal/export` 290.033s, `internal/issues` 214.655s, `internal/journal`
  216.519s, `internal/links` 70.284s, all also under 600s).
- `go test -tags integration -race -timeout 30m -count=1
  ./cmd/fid/integration/...`: `real 11m40.635s`, `ok
  cmd/fid/integration 700.168s` (PASS, <900s).

No gate was terminated by timeout; the workspace had no stray `go test`/
`fid.test` processes before this session began (`ps aux | grep -E "go test|
fid\.test"` returned nothing).

## Deviations from spec

- Both required `go test` gate commands were invoked with `-count=1`
  appended, matching task 10's recorded deviation, to force genuine
  execution rather than serve results from Go's build/test cache. This does
  not change test selection, flags, thresholds, or scope - it is a
  cache-bypass flag only - and the recorded PASS rows above cite the
  task-spec commands verbatim with genuinely observed termination and
  duration evidence from the `-count=1` runs.
- `internal/workflow/stop_conditions.go` was edited even though the task
  spec's file list names only `internal/workflow/resolve.go, evaluation.go,
  attention readers` for `internal/workflow`. `stop_conditions.go` holds
  every other Condition type's qualifying-fact logic
  (`qualifyingExecutorTermination`, `qualifyingRecoveryCondition`,
  `qualifyingBudgetCondition`, `qualifyingWorkEnded`, etc.) and is dispatched
  from `qualifyingFact` in `internal/workflow/stop.go` (explicitly in scope
  for "the closed Stop registries"); it is the only file where
  `review-recorded`/`finding-resolved`'s VAL-078 resolution behavior can be
  implemented, and is read as the concrete file behind the spec's "attention
  readers" file category. `resolve.go` and `evaluation.go` were read and
  required no change: neither Obligation evaluation nor `ResolveGoverningHumanOutcome`
  is affected by completing the closed Condition set. `cmd/fid/attention.go`
  and `cmd/fid/stop.go` were also read and required no change (see Summary).

## Concerns, technical debt, or notes for the Tech Lead

- None. `stop_human.go`'s Human text omits a qualifier for the task 9/10
  `decision_class`/`outcome_class` discriminators (only `budget_name`/
  `capability_id` get one today); this task kept that precedent for
  `finding_id` rather than diverging, but a future Human-text pass
  (task 14, CF-040) may want to add qualifiers for all four discriminators
  uniformly.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit

## Note for the QA Lead checkpoint

The independent QA Lead checkpoint after this task must PASS before task 12
is dispatched. Reproduction points for an independent reviewer, without this
session's context:

- `go test ./internal/journal/... -run TestTask11 -v` - the closed-schema
  matrix and nullable-binding confirmation, no database or fixtures required.
- `go test ./internal/workflow/... -run TestRaiseStopAcceptsTheClosedFourTypeJudgmentConditionSet -v`
  - domain-level Condition registration/round-trip proof.
- `go test -tags integration -race -timeout 5m -run TestP23AttentionResolution ./cmd/fid/integration/... -v`
  - the full real-process VAL-047/048/049/050/078 resolution matrix; builds
    its own temporary git workspace per subtest via `t.TempDir()`, no shared
    state with other tests.
- `internal/workflow/capabilities.go:105` and
  `internal/export/archive_profiles.go:5` both still declare
  `human_outcome_resolution:false` (VAL-047 first half; grep for the literal
  string to confirm without running the binary).
