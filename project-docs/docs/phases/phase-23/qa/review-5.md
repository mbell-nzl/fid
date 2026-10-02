# Phase 23 - Mid-Phase Review 5 (Tasks 7-9)

Reviewer: QA Reviewer
Date: 2026-09-22
Tasks covered: phase-23/task-7 through phase-23/task-9
Review type: INITIAL
Prior review: none

## Verdict

PASS

Tasks 7-9 supply the required freshness/supersession, generic Concern and
Deviation, and Decision Reference/protected-Decision behavior. Task 9 in
particular satisfies the heavy DEC-235/236/237/240 Condition-dispatch
constraints exactly: the per-type dispatch keeps the six retained Phase 22
`Condition` types byte-identical (only `decision_class` carries `omitempty`),
still calls `validateCurrentValueRelations("Condition", object)` so the Phase
22 nullable-binding relation is not silently dropped, uses a checked type
assertion for the discriminant instead of an unchecked panic-prone one, and
DEC-248's bounded CLI-adapter addition renders `protected-decision-required`
at exit 5 ahead of the generic Obligation denial. No Must-fix finding remains.

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-23/tasks/task-7.md`
  - `project-docs/docs/phases/phase-23/tasks/task-8.md`
  - `project-docs/docs/phases/phase-23/tasks/task-9.md`
- Completion reports:
  - `project-docs/docs/phases/phase-23/tasks/task-7-done.md`
  - `project-docs/docs/phases/phase-23/tasks/task-8-done.md`
  - `project-docs/docs/phases/phase-23/tasks/task-9-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-23/qa/validation-contract.md` rows
    VAL-003, VAL-031-033, VAL-034-037, VAL-074, VAL-038-042, VAL-079-084.
- QA plan:
  - `project-docs/docs/phases/phase-23/qa/qa-plan.md` sections 3 (task
    inventory), 4 (EVID-04, EVID-18-20), 8 (mid-phase cadence), and 21
    (bounded DEC-237 confirmation / DEC-239 evidence sync, QA-43/46-52).
- Decisions and supporting records:
  - `project-docs/docs/phases/phase-23/decisions.md` DEC-248 (Task 9 owns
    protected-Decision denial adapters).
  - `project-docs/docs/phases/phase-23/qa/task-9-supervisor-execution.md`
    (root-owned orchestration record, inspected for scope only).
  - `project-docs/docs/phases/phase-23/run-ledger.md` (Runner-owned; read
    only to confirm the Task 9 TASK-RECOVERY/REPAIR sequence, not changed).
- Commits or changed files:
  - Task 7 implementation `e4fe99111e035ebd72c547eda3a9dbf828cafa14`, report
    `3dd192ef9a34b0dffe7bc5ee293d50f2c779e81d`.
  - Task 8 implementation `3a0aed9d35ca2d43540e686a5cbbaa245177ea0f`, reports
    `1e7ba38aae1e2b2ec480fb0360ad6c36348de405` and
    `acf1a7b0a7bab1c836bdd4d121798935a924a718`.
  - Task 9 implementation `828eb18f422ebbfacbe4e65c62d53a061b7e9b9f` (checked
    in as `chore(task-9)` during an in-progress recovery pause) and DEC-248
    repair `fec3bdf8a1665f3b79bfdcd0146ca16b51527a51`, reports
    `733cacc42bef206be2176185acfbf82a2d43937a` and
    `99a4d0e4af421e6da5b5685ea169f6bcbb63c509`.

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 2
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-003, VAL-031, VAL-032, VAL-033,
  VAL-034, VAL-035, VAL-036, VAL-037, VAL-074, VAL-038, VAL-039, VAL-040,
  VAL-041, VAL-042, VAL-079, VAL-080 (Task 9/10 partial contribution),
  VAL-081, VAL-082, VAL-083, VAL-084

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-003 | tasks 8, 9 | ON TRACK | `TestConcernDeviationFenceConvergencePreservesEveryRemainingEvent` (`internal/journal/project_effects_test.go`) independently reproduced `PASS`; it proves task 8 deletes only its three fence keys and task 9 deletes only `decision.referenced`, with every remaining fenced key preserved. `unavailableJudgmentEvents` no longer contains `decision.referenced` (`internal/journal/project_effects.go`). |
| VAL-031 | task 7 | ON TRACK | `TestReviewEvidenceStalesOnTargetTreeChange` and the reproduced `TestP23FreshnessAndOutcomeSupersessionRealProcessJourney` (independently rerun, `PASS`) observe `stale`/`evidence-stale` after a target tree change with the Review Record unchanged. |
| VAL-032 | task 7 | ON TRACK | `freshnessForTarget` in `internal/evidence/evidence.go` returns `fresh` whenever `FreshnessTarget=="none"`; `TestAbsentHumanOutcomeFreshnessRemainsFresh` and the built-binary journey cover the absent and linked cases across the same tree change. |
| VAL-033 | task 7 | ON TRACK | `ResolveGoverningHumanOutcome` (`internal/workflow/resolve.go`) orders strictly on `created_event_id`, never a caller timestamp; `TestHumanOutcomeSupersessionFollowsRecordedEventOrderBothDirections` uses deliberately reversed timestamps and both directions pass. |
| VAL-034 | task 8 | ON TRACK | `TestConcernSubjectAndDispositionMatrix` and the reproduced real-process journey cover all six closed subject types, including the deliberately unresolved `requirement` subject, and prove no completion column exists and both guarded actions succeed with an open blocker Concern present. |
| VAL-035 | task 8 | ON TRACK | `validateConcernDispositionTargets` in `internal/review/concern.go` enforces FK-checked identities only where required (`carried-forward` Issue, `superseded` Concern/Finding); `ConcernAlreadyDisposedError` blocks a second disposal before target validation runs. The built-binary journey exercises one success per action plus all seven negatives with their exact codes. |
| VAL-036 | task 8 | ON TRACK | `deviation.go` exposes no update leaf; `TestLiveLeafCommandsMatchContractInventories` (`cmd/fid/contract_surface_test.go`) proves `deviation record`/`deviation list` exist with no update path; the journey records a two-row superseding chain with both rows readable after rebuild. |
| VAL-037 | task 8 | ON TRACK | `evidence.link`'s closed type domain only recognizes `receipt`/`review`/`human-outcome` (`internal/evidence/evidence.go:245`), so a Concern or Deviation is rejected structurally; the journey confirms `invalid-argument` with no Event for both. |
| VAL-074 | task 8 | ON TRACK | `TestConcernSelfDisposalKeepsRaisingAndDisposingActors` and the real-process journey store `actor_id` (raising) and `disposed_actor_id` separately, proving equal values for self-disposal and distinct values for a second Actor. |
| VAL-038 | task 9 | ON TRACK | `verifyDecisionSource` (`internal/review/decision.go`) is called before any journal write and rejects a missing file or digest mismatch with `SourceDigestMismatchError`; no prose column exists in `decision_refs`. `TestDecisionReferenceSourceMustExistAndMatchDigest` and the reproduced real-process journey confirm both negatives append no Event. |
| VAL-039 | task 9 | ON TRACK | `sortEvaluationReasons` ranks `protected-decision-required` (precedence 8) ahead of `obligation-evidence-absent` (10), so the correct denial wins even when both are present; DEC-248's `cmd/fid/envelope.go`/`issue_mutation_response.go` adapters map it to exit 5 with the exact four immutable details. The reproduced `TestP23DecisionReferenceProtectedLifecycleRealProcess` exercises both guarded adapters and the simultaneous-reasons case. |
| VAL-040 | task 9 | ON TRACK | `TestProtectedDecisionSatisfactionRequiresBothApprovedFactsAndExactAuthority` covers neither/reference-only/outcome-only/both, `rejected`/`deferred` reference outcomes, and a capability mismatch; only `both-approved` with exact authority satisfies. |
| VAL-041 | task 9 | ON TRACK | `protectedDecisionStopEffects` raises one Stop per unsatisfied Decision; `plannedDecisionResolutions` resolves only the Stop whose `ResolutionCondition.DecisionClass` matches the newly satisfied class, in the same Event. The reproduced journey starts two open Stops, satisfies one, and confirms only that Stop/Attention resolve. |
| VAL-042 | task 9 | ON TRACK | `internal/issues/transition_rules.go` adds `case workflow.RuleLinkedToDecision` querying `decision_refs` scoped to the Issue with `outcome='approved'`; `LinkedToDecision` in `internal/workflow/rules.go` carries no role parameter. `TestLinkedToDecisionRequiresApprovedReferenceInIssueScope` proves the before/after states. |
| VAL-080 (Task 9/10 partial) | task 9 | ON TRACK | `TestTask9DecisionConditionShapeAndNonStringDiscriminant` rejects wrong-field, null-own-field, two-new-field, and non-string-discriminant `decision-satisfied` Conditions and accepts only the canonical five-field form; reproduced `PASS`. Task 11 completes the four-type matrix per the contract. |
| VAL-081 | task 9 | ON TRACK | `internal/journal/current_event_value_relations.go`'s `case "Condition"` is untouched (byte-for-byte identical to the pre-change reference `74c83b7`); the new dispatch's final call is `return validateCurrentValueRelations("Condition", object)` (`current_event_value_types.go`), so the relation still fires. `TestTask9ConditionDispatchPreservesNullableBindings` reproduces all three negatives with the exact `condition nullable binding mismatch` message. |
| VAL-082 | task 9 | ON TRACK | `TestTask9ConditionDispatchRetainsPinnedPhase22Bytes` (reproduced `PASS`) accepts all six pinned canonical pre-change type-appropriate byte forms (including non-null `budget_name`/`capability_id` per QA I-44) and rejects one same-type mutation per type. |
| VAL-083 | task 9 | ON TRACK | `TestTask9PrechangeBudgetObservedJournalReplaysWithoutMutation` (reproduced `PASS`) replays a compressed, base64-embedded pre-change fixture referencing commit `74c83b7` through `RebuildWorkspace` and confirms the source-journal SHA-256 digest is byte-identical before and after. |
| VAL-084 | task 9 | ON TRACK | `TestTask9RetainedAttentionKeyMatchesPrechangeReference` pins `attentionSemanticKey` for a fixed `budget-observed` input. Independently confirmed: `attentionSemanticKey`'s body and the `Condition` struct's marshaled bytes for a `budget-observed` value (no `decision_class` emitted, `omitempty`) are unchanged between `74c83b7` and the Task 9 implementation, so the pinned key is not tautological. |
| VAL-079 | tasks 7-9 | ON TRACK | Each report records the exact ordinary and tagged gate durations below the 600 s/900 s budgets: task 7 `cmd/fid` 499.695 s / `internal/workflow` 211.628 s / integration 648.604 s; task 8 `cmd/fid` 538.754 s / `internal/workflow` 211.559 s / integration 651.834 s; task 9 `cmd/fid` 513.222 s / `internal/workflow` 431.647 s / integration 713.194 s. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Exact task-command audit | Each report's `## Exact quality-gate results` against its task spec's `## Exact Commands To Run` | PASS | All five required commands recorded `PASS` with observed termination in all three reports; every attempted failure (task 8's two lint iterations, task 9's claims-fixture and V01-fixture RED attempts and the concurrently loaded overlap) is retained with its exact FAIL token before a byte-identical rerun. |
| Report/spec mechanical control | `project-docs/tools/report-lint.sh` for tasks 7, 8, 9 | PASS | All three returned `OUTCOME COMPLETE` / `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| Commit identity and patch integrity | `git rev-parse --verify` and `git show --check` on all six implementation/report commits | PASS | All six hashes resolve to exactly the cited commits; `git show --check` produced no whitespace errors; every commit carries its Task 7/8/9 `Refs:` footer. |
| Format | `gofmt -l cmd internal` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`), independently rerun at current HEAD. |
| Lint | `golangci-lint run ./...` | PASS | `0 issues.` (terminated `exit_code: 0`), independently rerun at current HEAD. |
| Window real-process journeys | `go test -tags integration -race -timeout 30m ./cmd/fid/integration/... -run '^(TestP23FreshnessAndOutcomeSupersessionRealProcessJourney\|TestP23ConcernDeviationRealProcessJourney\|TestP23DecisionReferenceProtectedLifecycleRealProcess)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid/integration 13.898s` (terminated `exit_code: 0`); all three built-binary journeys covering VAL-031-033, VAL-034-037/074, and VAL-038-042 pass independently. |
| Task 9 Condition-dispatch direct assertions | `go test -race -timeout 10m ./internal/journal/... -run '^(TestTask9ConditionDispatchRetainsPinnedPhase22Bytes\|TestTask9ConditionDispatchPreservesNullableBindings\|TestTask9DecisionConditionShapeAndNonStringDiscriminant\|TestTask9PrechangeBudgetObservedJournalReplaysWithoutMutation\|TestConcernDeviationFenceConvergencePreservesEveryRemainingEvent)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/internal/journal 2.419s` (terminated `exit_code: 0`); reproduces VAL-081, VAL-082, VAL-083, VAL-003 (task 9 half), and the `decision-satisfied` shape matrix directly. |
| Task 9 Stop-detail/Condition strict shapes | `go test -race -timeout 5m ./internal/journal/... -run '^TestTask9ProtectedDecisionStopDetailsAndConditionAreStrict$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/internal/journal 1.020s` (terminated `exit_code: 0`); confirms the acceptance criterion that `currentStopDetailShapes`/per-type `Condition` entries are proved by a direct journal-layer assertion, not the task's own retained-family raise path. |
| File-scope audit | `git diff --name-only` for each task's commit range against its spec's `## Files To Create Or Modify` | PASS | Task 7, 8, and 9 file sets stay inside their declared scope (plus ordinary matching-test/fixture convergence explicitly permitted by each spec); no unlisted production file outside the accepted scope was touched. |
| Closed-allowlist consistency | `grep decision.record-reference` across `internal/workflow/activation.go` and `internal/journal/current_event_value_relations.go` | PASS | Both copies of the closed 27/28-name capability allowlist stay sorted and identical after task 9's addition, matching the R-19/DEC-233 dual-registry requirement. |
| DEC-248 scope conformance | `project-docs/docs/phases/phase-23/decisions.md` DEC-248 vs. `fec3bdf`'s changed files | PASS | The repair touches exactly `cmd/fid/envelope.go`, `cmd/fid/issue_mutation_response.go`, and matching tests/fixtures, as DEC-248 authorizes; no domain, gate, or interface rule changed. |

## Per-Task Findings

### Task 7 - Review freshness and Human Outcome supersession

#### Must-fix

- None.

#### Improvement

- None.

### Task 8 - Generic Concerns and append-only Deviations

#### Must-fix

- None.

#### Improvement

- None.

### Task 9 - Decision References and protected Decisions

#### Must-fix

- None.

#### Improvement

- [`project-docs/docs/phases/phase-23/tasks/task-9-done.md` Concerns section] The
  report retains an unresolved `TestLiveOutputGoldenFixtures` failure observed
  only under a root-caused concurrently loaded gate invocation (documented
  separately in `qa/task-9-supervisor-execution.md`); it did not recur in
  isolation or in the serial full ordinary gate this review independently
  reproduced. This is the same file CF-036/CF-038 already tracks
  (`internal/workflow` was already the note's subject in CF-038's carried-forward
  entry for `cmd/fid/live_output_test.go` coalescing hardening). Fix direction:
  when CF-038 is next picked up, fold in this task-9 concurrent-load
  observation as an additional reproduction data point rather than opening a
  new carry-forward.
- [`internal/workflow/protected_decision_test.go:118`] `TestTask9RetainedAttentionKeyMatchesPrechangeReference`
  pins the pre-change `attentionSemanticKey` value as a literal constant rather
  than computing it from a checked-out copy of the reference commit at test
  time. This review independently confirmed the pinned value is not
  tautological (the function body and the `Condition` struct's marshaled bytes
  are unchanged between the reference and Task 9), so this is not a Must-fix,
  but a future hardening could compute the reference value from an embedded
  historical fixture the way `TestTask9ConditionDispatchRetainsPinnedPhase22Bytes`
  and the VAL-083 replay fixture already do, removing the need for a reviewer
  to re-derive the proof by hand.

## Prior Finding Resolution

Not applicable: this is the initial window 7-9 review.

## Repair Plan

Not required: no Must-fix finding remains.

## Cross-cutting observations

- DEC-248 is a properly bounded, bookkept task-scope correction (two named CLI
  adapter files plus matching tests) that completes an already-accepted
  interface contract (VAL-039, `public-interface.md` section 5.1); it changes
  no domain rule, gate, or public-interface shape, and the Runner's
  TASK-RECOVERY/REPAIR ledger entries for `task-9:in-scope-implementation:protected-decision-cli-adapter-mapping`
  correctly bracket it.
- Task 9's `internal/journal/current_event_value_types.go` dispatch
  demonstrably satisfies all three hard requirements the task spec calls out:
  (a) the six Phase 22 shapes are byte-identical strings in
  `currentConditionShapes`; (b) the discriminant read uses a checked
  `object["type"].(string)` assertion, verified panic-free by
  `TestTask9DecisionConditionShapeAndNonStringDiscriminant/non-string-type`;
  and (c) the branch ends with `return validateCurrentValueRelations("Condition", object)`,
  independently reproduced through `TestTask9ConditionDispatchPreservesNullableBindings`.

## Blocker Candidates

- None. DEC-248 supplies bounded, decision-free authority for the one scope
  addition made in this window, and no other finding crosses task or
  authority boundaries.

## Existing Escalation References

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-23/qa/review-5.md
Must-fix: 0  Improvement: 2  Blocker candidates: 0
VERDICT: PASS
```
