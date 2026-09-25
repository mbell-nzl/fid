# Phase 23 - Task 10: Provenanced Human Outcomes - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 3
- Next action: None.
- Human attention required: no

## Summary of what was implemented

Provenanced Human Outcome recording: `internal/review/human_outcome.go` records
a Human Outcome bound to active Kit Activation provenance and a declared
human-outcome capability (never a role string or `--by human` flag),
`cmd/fid/human_outcome.go` wires `human-outcome record`/`human-outcome list`,
and `internal/evidence` gains the human-outcome Evidence type with the
declared `producer:"human"`, `trust_source:"activated-human"`,
`artifact_availability:"not-applicable"` triple and strict
class-must-equal-`outcome_class`/result-must-be-`approved` admission. A shared
readiness helper (`internal/workflow/human_outcome_readiness.go`,
`HumanOutcomeReadinessStopEffects` /
`humanOutcomeReadinessEffectsForEnvelope`) raises `human-outcome-required`
atomically with each of DEC-224's four named producer Events
(`envelope.activated`, `evidence.linked`, `review.recorded`,
`artifact.availability-observed`) exactly when the human-outcome Obligation
becomes the sole unsatisfied one, using an `overrides` map so each producer's
own not-yet-projected Effect (an Evidence link, a review resolution, an
Artifact-availability flip) is visible to the evaluation inside the same
committing Event. `internal/journal` gains the closed-registry, JSON-shape,
and typed-result wiring for `human-outcome.record`/`human-outcome.recorded`
(capability allowlist in `activation.go` and its second copy in
`current_event_value_relations.go`; Stop/Condition/detail shapes in
`stop.go` and `current_event_value_catalog.go`; typed result registration in
`current_event_json.go`/`current_event_result_projection.go`), and the
DEC-192 `human-outcome.recorded` fence key is deleted from
`project_effects.go` with its live-apply/replay parity re-proved. The task 7
and task 9 seeded Human Outcome projection-row fixtures are replaced with real
recorded Outcomes through the built binary.

This is a multi-session continuation: the implementation above was built and
independently code-reviewed across three prior sessions (recorded in the
run-ledger's attempt-2/3/4 entries), which also found and fixed two real
defects along the way (a `cmd/fid/artifact_test.go` fixture gap exposed by the
new readiness re-evaluation's real-git dependency, and a production
class-mismatch admission bug in `internal/evidence/evidence.go` /
`evidence_admission.go` that returned generic `invalid-argument` instead of
the spec-required `outcome-not-permitted`). This session verified the
inherited tree state, confirmed `internal/journal/current_event_result_budgets.go`'s
uncommitted change from the immediately-prior interrupted session is a
legitimate in-scope fix (`currentReviewBudgetResult` was unconditionally
gating Stop/Attention re-derivation on a `budget_samples` Effect being
present, which drops the Stop a `review.recorded` Event raises when it
resolves an Obligation with no budget sample at all - exactly the scenario the
new `TestP23HumanOutcomeRequiredRaisedByReviewRecord` fixture exercises),
resolved the VAL-046 review.recorded coverage decision (already answered by
the inherited work: a dedicated fixture, not a shared-core argument alone -
see below), ran the full required-command gate once cleanly end to end, and
closed out the report and commit.

### VAL-046 coverage decision

**Decision: a dedicated fixture for `review.recorded`, not a shared-core
argument.** This was already implemented by the immediately-prior session:
`TestP23HumanOutcomeRequiredRaisedByReviewRecord`
(`cmd/fid/integration/p23_human_outcome_test.go:252-335`) drives an
independent blocking Review, a second independent pass Review linked as the
review Obligation's Evidence while the Finding is still open, and a third
Review that resolves the Finding - proving `review.recorded` itself raises
`human-outcome-required` atomically inside that same Event when doing so
retroactively satisfies the already-linked review Obligation via
`reviewObligationReadinessOverride`. This override path is unique to
`review.recorded` (it is not exercised by the evidence.linked,
artifact.availability-observed, or envelope.activated fixtures), so
code-review confidence in the shared `HumanOutcomeReadinessStopEffects` core
alone would not have proven this producer's own override wiring behaves
correctly, or that the `currentReviewBudgetResult` no-budget-sample gap above
did not silently drop its Stop. `TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation`'s
doc comment (lines 176-189) records this reasoning and cross-references the
dedicated fixture and the other two producers' proof sites.

## Files created or modified

- `internal/review/human_outcome.go` (new) - Human Outcome recording rules,
  provenance binding, Effects.
- `internal/review/review.go` - `review.record`'s DEC-224 readiness override
  (`reviewObligationReadinessOverride`, `terminalResolvedFindingIDs`,
  `workflowRepository`) wired into `Record`.
- `internal/workflow/human_outcome_readiness.go` (new) +
  `internal/workflow/human_outcome_readiness_test.go` (new) - the shared
  cross-package readiness-Stop raising core and its four call sites.
- `internal/workflow/activation.go`, `internal/workflow/envelope.go`,
  `internal/workflow/evaluation.go`, `internal/workflow/stop.go` - the
  `human-outcome.record` capability allowlist entry, envelope-side wiring, the
  evaluation failure/reason precedence for `human-outcome-required`, and the
  closed Stop/Condition registries.
- `internal/evidence/evidence.go`, `internal/evidence/evidence_admission.go` -
  the Human Outcome Evidence type, its admission rule (class must equal
  `outcome_class`, result must be `approved`), the class-mismatch
  `outcome-not-permitted` fix, and the `evidence.linked` readiness-Stop call
  site.
- `internal/evidence/artifact.go` - the `artifact.availability-observed`
  readiness-Stop call site.
- `internal/evidence/human_outcome_admission_test.go` (new),
  `internal/evidence/human_outcome_binding_test.go` (new) - VAL-044 coverage
  and the class-mismatch regression.
- `internal/journal/current_event_catalog.go`,
  `internal/journal/current_event_value_relations.go`,
  `internal/journal/current_event_value_catalog.go`,
  `internal/journal/current_event_json.go`,
  `internal/journal/current_event_result_projection.go`,
  `internal/journal/current_event_result_budgets.go` - closed-registry,
  JSON-shape, typed-result, and budget-result wiring for
  `human-outcome.record`/`human-outcome.recorded`.
- `internal/journal/current_event_task10_values_test.go` (new) - direct
  journal-layer assertions for the Stop-details/Condition shapes and the
  Phase 22 nullable-binding relation (VAL-080, VAL-081).
- `internal/journal/project_effects.go`,
  `internal/journal/project_effects_test.go` - DEC-192 fence-key deletion and
  parity proof (VAL-003).
- `internal/export/archive_semantics_catalog.go` - mirrors the
  `evidence.linked`/`artifact.availability-observed` amendments into
  `semanticEventRules`.
- `internal/cli/codes.go` - `outcome-not-found`, `human-outcome-required`
  code registration (reuses task 6's `outcome-not-permitted`).
- `cmd/fid/human_outcome.go` (new) - `human-outcome record`/`human-outcome
  list` leaf wiring.
- `cmd/fid/root.go`, `cmd/fid/envelope.go`, `cmd/fid/evidence.go`,
  `cmd/fid/issue_mutation_response.go` - command registration and
  response/rendering wiring for the new leaf and readiness Stops.
- `cmd/fid/artifact_test.go` - added
  `observedGateReadyEvidenceEnvelopeWithSourceForTest` fixture (real-git
  dependency introduced by the new readiness re-evaluation) and the
  `artifact.availability-observed` VAL-046 raising-path assertions in
  `TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation`.
- `cmd/fid/contract_surface_test.go` - closed production-path/leaf-command
  inventory updates for the new files and commands.
- `cmd/fid/integration/p23_human_outcome_test.go` (new) - VAL-043/044/045/046/
  009/010/025/081 real-process coverage, including the dedicated
  `review.recorded` and `envelope.activated` raising fixtures.
- `cmd/fid/integration/p22_fake_test.go`,
  `cmd/fid/integration/p23_decision_reference_test.go`,
  `cmd/fid/integration/p23_freshness_supersession_test.go`,
  `cmd/fid/integration/p23_review_record_test.go` - task 7/task 9 seeded
  Human Outcome fixtures converged onto real recorded Outcomes; helper
  additions (`containsAction`).
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - `human-outcome
  record`/`human-outcome list` CLI documentation.
- `project-docs/docs/phases/phase-23/run-ledger.md` - recovery ledger entries
  from prior sessions (Runner-owned; not edited by this session beyond what
  was already committed-pending).

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - new `human-outcome`
  subcommand documentation.
- `project-docs/docs/phases/phase-23/tasks/task-10-done.md` - this report.

## Validation Contract Coverage

- VAL-043: PASS - `cmd/fid/integration/p23_human_outcome_test.go:34-65`.
  Agent-actor-type (`agentAttempt`, line 35), unregistered/unfinalized-actor
  (`unregistered`, line 42), and missing-capability
  (`p23RecordHumanOutcomeMissingCapability`, line 47) each assert
  `assertP23AuthorityMissing` (exit 5, `authority-missing`, no Event via the
  `p21JournalEventCount` before/after check). The valid case (`release`, line
  52) asserts success, exactly one appended Event, and the recorded
  `provenance.activation_id` matches the finalized Activation (line 63-65).
- VAL-044: PASS - `cmd/fid/integration/p23_human_outcome_test.go:70-86` (real
  producer: mismatched-type reject at line 70-71, matching-class-approved
  satisfies at line 72-78) plus
  `internal/evidence/human_outcome_admission_test.go:17-42`
  (`TestHumanOutcomeAdmissionSatisfiesOnlyMatchingClassWithApprovedResult`:
  matching class with `approved` admits; `rejected`/`changes-required`/
  `acknowledged` each reject `outcome-not-permitted`) and
  `internal/evidence/human_outcome_binding_test.go:23-39`
  (`TestLoadHumanOutcomeEvidenceBindingRejectsClassMismatchAsOutcomeNotPermitted`:
  same-type, different-`outcome_class` link rejects `outcome-not-permitted`,
  not the prior generic binding error).
- VAL-025: PASS - `cmd/fid/integration/p23_human_outcome_test.go` full test.
  Real recorded Human Outcome linked to its matching human-outcome Obligation
  (line 72-78, evaluation `status:"satisfied"`); mismatch to gate-evidence
  (line 70-71) and to review type via `security-vs-gate` (line 149-157, both
  asserting `outcome-not-permitted` with `obligation_type:"gate-evidence"`
  detail and no Event/projection change via the before/after
  `p21JournalEventCount` checks surrounding each call). Combined with task 6's
  receipt/review-positive and four residual mismatch cells
  (`task-6-done.md`), this closes all nine pairings per DEC-243; task 6's
  existing tests remained exercised by this task's full gate (`go test -race
  ./...` above).
- VAL-045: PASS - `cmd/fid/integration/p23_human_outcome_test.go:76,81,84`.
  Linked Evidence carries `producer:"human"` (line 81), `trust_source:
  "activated-human"` (line 84); the linked evaluation carries
  `artifact_availability:"not-applicable"` (line 76), re-confirmed at line
  137-140 (VAL-010). No field or reasoning text in `internal/review/human_outcome.go`
  or the Evidence rendering path claims deterministic or fid-executed proof.
- VAL-046: PASS - four dedicated raising-path assertions, one per DEC-224
  producer:
  - `envelope.activated`: `TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation`
    (`cmd/fid/integration/p23_human_outcome_test.go:190-250`) - activating an
    Envelope whose sole Obligation is already human-outcome raises the Stop
    and Attention immediately, with exact `obligation_id`/`outcome_class`
    details.
  - `evidence.linked`: `TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess`
    (same file, lines 88-103) - linking the qualifying gate-evidence Evidence
    raises `human-outcome-required` for the `security` class atomically
    inside that Event (`assertP23HumanOutcomeStopState(t, a, "security",
    "open")`).
  - `review.recorded`: `TestP23HumanOutcomeRequiredRaisedByReviewRecord`
    (lines 252-335, see VAL-046 decision above) - a dedicated fixture
    exercising the producer-unique `reviewObligationReadinessOverride` path,
    asserting exactly one raised Stop/Attention in the resolving Review's own
    response with correct `obligation_id`/`outcome_class` details.
  - `artifact.availability-observed`: `TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation`
    (`cmd/fid/artifact_test.go:271-`, using
    `observedGateReadyEvidenceEnvelopeWithSourceForTest`) - the documented
    I-24 missing-then-restored blob technique drives one qualifying
    `artifact.availability-observed` Event that raises the Stop.
  Routine-status non-raising is implicit in every other Obligation-satisfying
  step across these same fixtures never raising an unexpected Stop (asserted
  via exact Stop/Attention counts at each raising point, e.g. lines 234-236,
  245-248, 314).
- VAL-009: PASS - `cmd/fid/integration/p23_human_outcome_test.go:110-115`
  (`assertP23HumanOutcomeRequiredDenial`) - `envelope.complete` rejects
  `human-outcome-required` (rank 9) with `obligation_id:"OBL-human-security"`,
  `outcome_class:"security"` while that Obligation is the last unsatisfied
  one, and appends no Event. Lines 142-144 confirm `envelope.complete`
  succeeds once every Obligation (gate-evidence, review, human-outcome) is
  satisfied - the first task where that state is reachable.
- VAL-010: PASS - `cmd/fid/integration/p23_human_outcome_test.go:137-140` -
  with linked review and Human Outcome Evidence both present,
  `artifact_availability` stays `not-applicable` (not `available`/`missing`).
- VAL-003: PASS - `internal/journal/project_effects.go` (fence key deletion:
  `human-outcome.recorded` removed from `unavailableJudgmentEvents`) and
  `internal/journal/project_effects_test.go:62-71`
  (`TestConcernDeviationFenceConvergencePreservesEveryRemainingEvent`,
  updated assertion), exercised live in this session's full `go test -race
  ./...` run (`internal/journal` package: PASS).
- VAL-080: PASS - `internal/journal/current_event_value_catalog.go`
  (`currentConditionShapes["human-outcome-recorded"] = "type:=human-outcome-recorded
  subject:Subject budget_name:?ID capability_id:?ID outcome_class:ID"`) and
  `internal/journal/current_event_task10_values_test.go:55-82`
  (`TestTask10HumanOutcomeConditionShapeAndNonStringDiscriminant`: canonical
  shape accepted; `decision_class` instead of `outcome_class` rejected; null
  own field rejected; two new fields rejected). Enforced structurally at the
  journal layer, per-type; task 11 completes the four-type matrix.
- VAL-081: PASS - `internal/journal/current_event_task10_values_test.go:38-53`
  (`TestTask10HumanOutcomeConditionDispatchPreservesNullableBinding`: a
  `human-outcome-recorded` Condition carrying a stray `budget_name` is
  rejected with the exact observed message `condition nullable binding
  mismatch`) and `cmd/fid/integration/p23_human_outcome_test.go:468-485`
  (`p23AssertHumanOutcomeConditionNullableBindingFires`, invoked at line 173)
  re-confirming the real, built-binary-recorded Stop carries the canonical
  null-bound Condition form. The Phase 22 binding still fires from task 9's
  dispatch; this is not a task 9 defect.

## Tests run

- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	209.135s	coverage: 62.3% of statements` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	670.837s` (terminated `exit_code: 0`)

Real entrypoint exercise: the tagged integration suite above drives the
actual built `fid` binary as a real process
(`integrationFIDBinary`/`closureRunCommand` in
`cmd/fid/integration/p23_human_outcome_test.go`), including
`human-outcome record` and `human-outcome list` through the real CLI leaf
(`p23RecordHumanOutcome`, and the `human-outcome list` call at line 159-169:
`listed := a.p22Run(t, "human-outcome", "list", "--issue", a.issueID)`, exit
0, 2 items). This satisfies the runtime-surface entrypoint-exercise
requirement; no separate manual CLI invocation was additionally run in this
session.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	209.135s	coverage: 62.3% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					61.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	670.837s` (terminated `exit_code: 0`)

`goimports -l cmd internal` was also run (not a task-spec-required command,
applied per the instruction to apply gofmt/goimports to edited files before
the final gate): `exit 0 (no output)`, terminated `exit_code: 0`.

Wall-clock durations (VAL-079/DEC-227 budgets: tagged integration <900s,
ordinary packages <600s): `go test -race -timeout 30m -coverprofile=coverage.out
./...` reported `real 9m1.954s`, with per-package `ok` durations `cmd/fid`
539.822s (PASS, <600s) and `internal/workflow` 209.135s (PASS, <600s); no
ordinary package exceeded budget. `go test -tags integration -race -timeout
30m ./cmd/fid/integration/...` reported `real 11m12.143s`, with `ok
cmd/fid/integration` 670.837s (PASS, <900s).

Both `go test` commands above were actually invoked with `-count=1` appended
(`go test -race -timeout 30m -count=1 -coverprofile=coverage.out ./...` and
`go test -tags integration -race -timeout 30m -count=1
./cmd/fid/integration/...`) to force a genuine execution in this session
rather than Go's build/test cache; this is recorded as a deviation below. An
initial `go test -race -timeout 30m -coverprofile=coverage.out ./...` attempt
without `-count=1` returned entirely from cache (`(cached)` on every package
line, no real execution in this session) and is not recorded as evidence; it
was immediately rerun with `-count=1` added, and that rerun is the PASS
recorded above with its own genuinely observed termination and durations.

## Deviations from spec

- Both required `go test` gate commands were actually invoked with `-count=1`
  appended (see Exact quality-gate results) to force genuine execution rather
  than serve results from Go's build/test cache. This does not change test
  selection, flags, thresholds, or scope - it is a cache-bypass flag only -
  and the recorded PASS rows cite the task-spec commands verbatim with
  genuinely observed termination and duration evidence from the `-count=1`
  runs.

## Concerns, technical debt, or notes for the Tech Lead

- `internal/journal/current_event_result_budgets.go`'s
  `currentReviewBudgetResult` previously returned an empty `stops`/`attention`
  result whenever a `review.recorded` Event carried no `budget_samples`
  Effect at all, which would have silently dropped the
  `human-outcome-required` Stop/Attention this task's `review.recorded`
  raising path produces for a pass-verdict Review that satisfies the last
  blocking Finding with no budget sample involved. This was generalized (not
  narrowed) to always re-derive `stops`/`attention` from the Event's own
  Effects regardless of whether a budget sample is present, only gating the
  `budget` field on the sample's presence. This is in-scope for this task
  (task 10's own new `review.recorded` call site was not previously reachable
  without it) and is exercised by
  `TestP23HumanOutcomeRequiredRaisedByReviewRecord`.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `0ff2eeef66d4633309a43ae6a341aecd38def4d0`

## Commit subject

- `feat(review): record provenanced Human Outcomes with readiness Stops`
