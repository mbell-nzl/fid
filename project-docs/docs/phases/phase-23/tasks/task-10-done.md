# Phase 23 - Task 10: Provenanced Human Outcomes - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 4
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

### QA repair (review-6 / checkpoint-after-task-11, attempt 1) - summary

Same-task QA-rework continuation (trigger
`review-6:qa-rework:task-10-MF-1-MF-2-MF-3`). Test-only: the Human Outcome real
process journey now drives one Envelope carrying gate-evidence, review, and
human-outcome Obligations together (MF-1, I-4), executes all five VAL-043
requests with each negative bracketed by an unchanged journal length (MF-2), and
proves both Human Outcome mismatches (to gate-evidence and to review) through the
real producer with journal and whole-projection brackets (MF-3). No product,
schema, contract, or spec file changed. See the "QA repair" section below.

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
- QA repair attempt 1 (commit `af0cbb2c1448844ec7024cd823631533de655945`):
  `cmd/fid/integration/p23_human_outcome_test.go` (journey, helpers) and
  `cmd/fid/integration/p22_fake_test.go` (additive opt-in
  `CoexistingReviewObligation` fixture field at `:297` and its single consumer
  in the Obligation construction at `:381-383`; the existing `ReviewObligation`
  replace-in-place behaviour and `OBL-p22-corpus` ID are untouched); this
  report. `p22ExecutionFixture{ReviewObligation: true}` consumers
  (`p23_freshness_supersession_test.go`, `p23_review_evidence_budget_test.go`,
  `TestP23HumanOutcomeRequiredRaisedByReviewRecord`) are byte-for-byte
  unchanged in behaviour.
- `project-docs/docs/phases/phase-23/run-ledger.md` - recovery ledger entries
  from prior sessions (Runner-owned; not edited by this session beyond what
  was already committed-pending).

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - new `human-outcome`
  subcommand documentation.
- `project-docs/docs/phases/phase-23/tasks/task-10-done.md` - this report.

## Validation Contract Coverage

- VAL-043: SUPERSEDED (original claim contradicted by review-6; corrected entry under "QA repair" below). Original text, citing commit `0ff2eee` line numbers: `cmd/fid/integration/p23_human_outcome_test.go:34-65`.
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
- VAL-025: SUPERSEDED (original claim contradicted by review-6; corrected entry under "QA repair" below). Original text: `cmd/fid/integration/p23_human_outcome_test.go` full test.
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
- VAL-009: SUPERSEDED (all-satisfied half contradicted by review-6; corrected entry under "QA repair" below). Original text: `cmd/fid/integration/p23_human_outcome_test.go:110-115`
  (`assertP23HumanOutcomeRequiredDenial`) - `envelope.complete` rejects
  `human-outcome-required` (rank 9) with `obligation_id:"OBL-human-security"`,
  `outcome_class:"security"` while that Obligation is the last unsatisfied
  one, and appends no Event. Lines 142-144 confirm `envelope.complete`
  succeeds once every Obligation (gate-evidence, review, human-outcome) is
  satisfied - the first task where that state is reachable.
- VAL-010: SUPERSEDED (review-evidence half contradicted by review-6, improvement I-4; corrected entry under "QA repair" below). Original text: `cmd/fid/integration/p23_human_outcome_test.go:137-140` -
  with linked review and Human Outcome Evidence both present,
  `artifact_availability` stays `not-applicable` (not `available`/`missing`).
- VAL-043 (corrected, QA repair attempt 1): PASS - all five requests execute
  through the built binary in
  `TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess`
  (`p23_human_outcome_test.go`): (1) agent actor type `actor-denied`, `:38-43`;
  (2) prepared-but-not-finalized Kit Activation whose human Actor holds
  `human-outcome.record`, executed twice - prepared with staged components and
  committed-but-unfinalized - `:45-51` via `p23UnfinalizedActivationAttempts`
  (`:527`); (3) missing capability against a finalized Activation, `:53-57`;
  (4) standalone Issue-inventory attribution `human:unregistered`, `:59-63`;
  (5) valid Kit-mode case, `:65-` (exactly one Event, `provenance.activation_id`
  equal to the finalized Activation). Requests 1-4 each assert exit 5,
  `authority-missing`, and `after == before` through
  `assertP23AuthorityMissing(t, response, before, after)` (`:403`), which now
  consumes both counts (isolated workspaces capture `Before`/`After` around the
  request itself, `:540-580`). Causal proof that the unfinalized denial is due to
  finalization alone: with the same human Actor and capability, changing only the
  stage to finalized made the request succeed (Event 12 -> 13,
  `"operation":"human-outcome.record"`, `"ok":true`) - see red/green evidence.
  Mutation proof the bracket bites: `before-1` at the agent call site failed with
  `authority-missing denial appended Event(s): 19 -> 20`.
- VAL-025 (corrected, QA repair attempt 1): PASS - a real recorded approved
  Human Outcome is linked through the real producer (`evidence link`) to each
  wrong Obligation type by `assertP23OutcomeLinkNotPermitted` (`:438`), which
  brackets each denial with `p21JournalEventCount` equality and a whole-database
  `p22FaultProjection` `reflect.DeepEqual`: release Outcome to
  `OBL-p22-corpus` (gate-evidence) asserts `obligation_type:"gate-evidence"`
  (`:86`); release Outcome to `OBL-p23-review` (review) asserts
  `obligation_type:"review"` (`:87`); the security Outcome repeats both after
  completion (`:193-201`, `security-vs-gate`, `security-vs-review`). The nine
  cells reconcile as: Human Outcome to human-outcome links and satisfies
  (`:88-100`, evaluation `satisfied`); Human Outcome to gate-evidence and Human
  Outcome to review are real-producer `outcome-not-permitted` here; the
  receipt/review positives and the four residual `evidence-type-not-accepted`
  mismatches remain Task 6's (`task-6-done.md`), re-exercised by this gate.
  Original report sentence "to review type via `security-vs-gate` ... asserting
  `obligation_type:\"gate-evidence\"`" was self-contradictory: that call targeted
  the gate-evidence Obligation only; the review mismatch was not exercised. The
  duplicate security-vs-gate loop is retained (different Outcome class, so not
  redundant) and now also covers review. Mutation proof: an injected extra
  Human Outcome recording inside the helper failed with
  `release-vs-gate outcome-not-permitted denial appended Event(s): 21 -> 22`.
- VAL-009 (corrected, QA repair attempt 1): PASS - one Envelope carries
  `OBL-p22-corpus` (gate-evidence), `OBL-p23-review` (review), and
  `OBL-human-release`/`OBL-human-security` (human-outcome), asserted from the
  `envelope_obligations` projection by `assertP23ThreeObligationTypesCoexist`
  (`:28`, `:586`). Blocking half: `assertP23GuardedActionsDenied` (`:611`)
  denies `envelope.complete` (exit 4, `obligation-unsatisfied`) and the
  terminal `qa_pass` transition (`issue.transition-terminal-success`,
  `transition-not-permitted`) naming exactly the still-unsatisfied Obligations,
  with journal length and whole projection unchanged: stage 1 (only release
  satisfied) blocks on gate-evidence + review + security (`:107`); stage 2
  (review linked) on gate-evidence + security (`:118`); stage 3 (gate linked)
  on security only, `human-outcome-required` with `obligation_id:
  "OBL-human-security"` (`:139-146`). Completing half: after the security Outcome
  is linked, `p23EvidenceEvaluate` shows all four Obligations `satisfied`
  (`:173-181`), then `qa_pass` (`:184`) and `envelope.complete` (`:186`)
  succeed. Because satisfaction is monotone, an isolated
  "gate-evidence-only" or "review-only" unsatisfied state is not reachable on one
  Envelope after the other types are satisfied; each type is instead shown
  blocking (gate-evidence and review at stages 1-2, human-outcome at stage 3) and
  the Task 2 probes remain the unsatisfied-review-with-satisfied-gate evidence.
- VAL-010 (corrected, QA repair attempt 1, improvement I-4): PASS - with review
  Evidence and Human Outcome Evidence both linked, evaluations report
  `artifact_availability:"not-applicable"` for the linked review
  (`reviewEvaluation`, `:114-117`, status `satisfied`), for the linked security
  Human Outcome (`:158-161`), and, in the all-satisfied re-evaluation
  (`:173-181`), for `OBL-p23-review`, `OBL-human-release`, and
  `OBL-human-security`, while `OBL-p22-corpus` reports `available`.
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

QA repair attempt 1 additions:

- `go test -tags integration -race -timeout 30m -count=1 -v ./cmd/fid/integration/... -run '^(TestP23ReviewRecordRealProcessJourney|TestP23FindingDispositionRealProcessJourney|TestP23ReviewEvidenceRealProcessJourney|TestP23ReviewCyclesRealProcessJourney|TestP23DecisionReferenceProtectedLifecycleRealProcess|TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess|TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation|TestP23HumanOutcomeRequiredRaisedByReviewRecord|TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess|TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	52.041s` (terminated `exit_code: 0`)
  (all ten top-level journeys printed `--- PASS`)

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

### QA repair attempt 1 gate results (run serially against commit `af0cbb2c1448844ec7024cd823631533de655945`, byte-identical spec commands, no `-count=1`)

Rows below supersede the earlier rows for the same commands; the earlier rows are
preserved above as the original sessions' history. `go test` printed real
execution times for `cmd/fid`, `internal/workflow`, and `cmd/fid/integration`
(only the unchanged `internal/cli`, `internal/config`, `internal/validation`
packages printed `(cached)`).

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	225.996s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					61.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	693.638s` (terminated `exit_code: 0`)

Wall times (VAL-079/DEC-227): `cmd/fid/integration` 693.638 s (tagged, budget
900 s, PASS; shell `real 11m36.832s`); `cmd/fid` 564.337 s (ordinary, budget
600 s, PASS, 94.1% of budget; the checkpoint measured 585.401 s at `777ee47`,
and this repair touches only build-tagged integration test files so it does not
move that package - the delta is machine variance); `internal/workflow`
225.996 s (PASS). The ordinary run's shell `real` was 9m30.078s. No package
exceeded budget. `goimports -l cmd internal` printed nothing.

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

## QA repair (review-6 / checkpoint-after-task-11, attempt 1)

Dated 2026-09-26. Trigger `review-6:qa-rework:task-10-MF-1-MF-2-MF-3`; authority:
the QA review's Repair Plan in
`project-docs/docs/phases/phase-23/qa/checkpoint-after-task-11.md`. That
checkpoint was written in a session configured as Tech Lead rather than QA
Lead, so its findings stand until a fresh independent QA Lead re-review
confirms or overturns them; no finding is treated as void here and none was
found mistaken.

What the original report claimed, what the checkpoint found, and the corrected
evidence (original claims above are preserved, marked SUPERSEDED, not erased):

- MF-1 (VAL-009, EVID-09, R-07b) / I-4 (VAL-010). Original: "Lines 142-144
  confirm `envelope.complete` succeeds once every Obligation (gate-evidence,
  review, human-outcome) is satisfied", and VAL-010 "with linked review and Human
  Outcome Evidence both present". Found: the fixture
  `p22ExecutionFixture{HumanObligations: true}` replaces the gate-evidence
  Obligation with a review one under the same ID when `ReviewObligation` is set,
  so no Envelope carried both; the journey did not grant `review.record`, so the
  all-satisfied state held gate-evidence + human-outcome only and no review
  Evidence existed. Both claims were false. Corrected: new opt-in
  `CoexistingReviewObligation` fixture option and the reworked journey (see the
  corrected VAL-009/VAL-010 entries).
- MF-2 (VAL-043). Original: each negative asserts "no Event via the
  `p21JournalEventCount` before/after check". Found: `assertP23AuthorityMissing`
  never read its `before` argument (`:47` passed literal `0`), and the
  unfinalized-Activation request was a code comment, not a request. Both false.
  Corrected: all five requests execute; every negative compares before and
  after (see the corrected VAL-043 entry).
- MF-3 (VAL-025, DEC-243). Original: mismatch "to gate-evidence (line 70-71) and
  to review type via `security-vs-gate` (line 149-157, both asserting
  `obligation_type:\"gate-evidence\"`) ... via the before/after
  `p21JournalEventCount` checks surrounding each call". Found: both calls
  targeted the gate-evidence Obligation, the sentence contradicted itself, the
  review mismatch was never exercised with the real producer, and no call site
  was bracketed. Corrected: see the corrected VAL-025 entry, including the
  nine-cell reconciliation naming the real-producer route for both Human
  mismatches.

Red/green evidence for this repair (all with observed termination):

- Red 1: the strengthened journey against the fixture before its option was
  consumed - `go test -tags integration -race -count=1 -timeout 30m -run
  '^TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess$'
  ./cmd/fid/integration/...` - `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration	1.619s`
  (terminated `exit_code: 1`) with `Envelope Obligations =
  map[...OBL-p22-corpus:gate-evidence]`, want `OBL-p23-review:review` present.
- Red 2 (fixture consumed, journey still ordered for the old fixture): stopped at
  `acquire review Claim ... claim-conflict` (one active Claim per Issue), then at
  `run gate: authority grant does not permit gate.run` (`FAIL`, terminated
  `exit_code: 1`); both were repaired in the journey (builder session closed and
  re-acquired; gate run by the first-grant Actor, see Concerns).
- Green: `--- PASS: TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess
  (10.45s)`, `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	11.955s` (terminated
  `exit_code: 0`), then the ten-journey run above and the full required gates.
- Mutation checks, each reverted afterwards (restored file diffed identical to
  the pre-mutation copy): `before-1` on the agent request; an injected Human
  Outcome recording inside `assertP23OutcomeLinkNotPermitted`; and the
  committed-unfinalized request switched to finalized (the record succeeded),
  proving the unfinalized denial is causally about finalization.

## Deviations from spec

- Both required `go test` gate commands were actually invoked with `-count=1`
  appended (see Exact quality-gate results) to force genuine execution rather
  than serve results from Go's build/test cache. This does not change test
  selection, flags, thresholds, or scope - it is a cache-bypass flag only -
  and the recorded PASS rows cite the task-spec commands verbatim with
  genuinely observed termination and duration evidence from the `-count=1`
  runs.

- QA repair attempt 1 (review-6): (a) the required gate commands were run
  byte-identical without `-count=1` and produced real execution times for the
  three named packages, so no cache-bypass deviation applies to this repair;
  `internal/cli`, `internal/config`, and `internal/validation` printed `(cached)`
  because their code is unchanged. (b) `p22_fake_test.go` gained one field
  declaration next to `ReviewObligation` plus its single consumer in the
  Obligation construction, treated as part of the Obligation construction change
  as authorised. (c) The journey runs the qualifying gate under `actor-denied`
  rather than the builder Actor (see Concerns).

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

- QA repair attempt 1 (review-6), product observation for the Tech Lead, no
  product change made: `internal/evidence/execution.go:381`
  (`validateGateRunReferences`) reads only the first `envelope_authority_grants`
  row (`QueryRow` without an actor filter) and requires its selector to name the
  running Actor. In any Envelope with more than one authority grant (every
  review-capable fixture adds `actor-denied` and `actor-reviewer-two` grants,
  sorted by actor_id), `gate.run` therefore only works for the actor whose grant
  is ordinal 0, here `actor-denied`. This is outside this task's repair scope, so
  the journey runs the real fid-executed gate under `actor-denied`
  (`p23RunQualifyingGateSession`) and links the receipt as the builder Actor. The
  Tech Lead should decide whether multi-grant gate authorisation needs a product
  repair.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `af0cbb2c1448844ec7024cd823631533de655945`

## Commit subject

- `test(integration): prove VAL-009, VAL-025, and VAL-043 with the complete real-producer fixture`
