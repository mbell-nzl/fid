# Phase 23 - Task 11: Attention resolution through Human Outcomes - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: gate-repair
- Condition: The exact ordinary race/coverage command exited 0, but `cmd/fid` took 686.386s and breached the accepted 600s package budget. Attribution and a passing budget gate remain unresolved.
- Recovery attempts: 1
- Next action: Continue task 11 with bounded timing attribution and an authorized causal speed repair if the DEC-250 implementation is responsible; retain the 686.386s failure and rerun the byte-identical required gate after repair. If the cause is outside task-11 scope, report the scope gap for Tech Lead disposition.
- Human attention required: no

## Summary of what was implemented

The original task-11 commit `c5bf0f89c69c57003ece353ef005a7596dd266e4` completed the four judgment Condition encodings and the `review-recorded` and `finding-resolved` resolution paths. The accepted DEC-250 repair commit `466745f2dc4505d60eb67ac50193e004532b7239` adds one shared qualifying-fact predicate for each of `human-outcome-recorded` and `decision-satisfied`. Manual `stop resolve` and atomic resolution in `human-outcome.recorded` and `decision.referenced` use the same predicates. Resolution scans the recorded Condition, so producer-raised Stops and `operator-stop` Stops qualify alike. A `decision-satisfied` Condition on an activation subject now rejects through the existing invalid-Condition path with no Event. No contract, schema, Event, failure code, or public-interface shape changed.

The review-7 BC-3 red observation remains preserved at `qa/checkpoint-after-task-11-2.md:744-755`: a matching approved Outcome resolved the readiness Stop but left the matching `operator-stop` open, and manual resolution rejected `EVT-20260925-000026` with `resolution-invalid` and Event count `26->26`. The current focused built-binary tests are green (Tests run). This report replaces the earlier COMPLETE claim, whose pre-repair implementation was `c5bf0f8` and whose recorded Go gates appended `-count=1`.

## Files created or modified

- Original task: `internal/workflow/stop.go`, `stop_conditions.go`, `stop_test.go`; `internal/journal/current_event_value_catalog.go` and `current_event_task11_values_test.go`; `cmd/fid/integration/p23_attention_resolution_test.go`.
- DEC-250 repair: `internal/workflow/stop.go`, `stop_conditions.go`, `stop_judgment_conditions_test.go`, `evaluation.go`, `resolve.go`; `internal/review/human_outcome.go` and `decision.go`; `cmd/fid/integration/p23_attention_resolution_test.go` and `p23_stop_resolution_conditions_test.go`.
- This continuation: `project-docs/docs/phases/phase-23/tasks/task-11-done.md` only.

## Documentation updated

- This report records DEC-250/253 repair evidence and the current budget-gate failure. Accepted specifications and QA artifacts remain unchanged.

## Validation Contract Coverage

- VAL-047 (task-11 half): `TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess` checks `attention.v1` reports `human_outcome_resolution:false` at `cmd/fid/integration/p23_attention_resolution_test.go:118`; the built-binary resolution tests below establish the path before task 17 flips the limit.
- VAL-048: `TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen` at `cmd/fid/integration/p23_attention_resolution_test.go:169-207` records wrong-class and wrong-result Outcomes, then `assertP23StopResolveRejected` checks exit 4 `resolution-invalid`, unchanged Event count, and open producer and operator Stops. The approved matching Outcome resolves both.
- VAL-049: `TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess` at `cmd/fid/integration/p23_attention_resolution_test.go:83-114` rejects an extra `outcome_id` input with no Event, accepts the existing `attention_id`/`expected_version`/`fact_event_ids` shape, and re-reads the original Stop identity and exact resolution fact IDs.
- VAL-050: `TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen` at `cmd/fid/integration/p23_attention_resolution_test.go:156-208` resolves the producer Stop and a matching-Condition `operator-stop` in the same Outcome Event; the different `review-recorded` Stop remains open.
- VAL-078: `TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess` at `cmd/fid/integration/p23_attention_resolution_test.go:52-81` covers non-qualifying and qualifying Review facts for `review-recorded` and `finding-resolved`. `TestP23StopResolveHumanOutcomeConditionKeysOnConditionNotRaisingCode` at `cmd/fid/integration/p23_stop_resolution_conditions_test.go:26-128` covers non-qualifying and qualifying Outcome facts on producer and operator Stops, including envelope, run, and activation subjects. `TestP23StopResolveDecisionSatisfiedConditionKeysOnConditionNotRaisingCode` at `cmd/fid/integration/p23_stop_resolution_conditions_test.go:130-228` covers rejected/deferred/earlier-half facts and both orders of the qualifying pair on producer and operator Stops, plus the activation-subject raise rejection with no Event. `TestHumanOutcomeConditionQualifyingFactRules` and `TestDecisionSatisfiedConditionQualifyingFactRules` at `internal/workflow/stop_judgment_conditions_test.go:75` and `:129` cover scope, supersession, and pair semantics.
- VAL-080: `TestTask11ReviewRecordedAndFindingResolvedConditionShapes` at `internal/journal/current_event_task11_values_test.go:12` covers canonical, extra, wrong, missing, and null discriminator shapes; `TestTask11ConditionShapeRegistryIsClosedAtTenRows` at `:67` checks the ten-row registry. The original task-11 evidence for the task-9/10 types remains in their committed tests. Malformed shapes reject through shape validation; no relation rule was added.
- VAL-081: `TestTask11FindingResolvedConditionDispatchPreservesNullableBinding` at `internal/journal/current_event_task11_values_test.go:52-62` observes the exact error `condition nullable binding mismatch` for a `finding-resolved` Condition with stray `budget_name` after task-9 dispatch. `TestTask9ConditionDispatchPreservesNullableBindings` at `internal/journal/current_event_closed_values_test.go:140` retains the three Phase 22 negatives and canonical positives. The relation was unchanged.

## Tests run

- Focused built-binary green: `go test -tags integration -race -timeout 5m -run 'TestP23StopResolve(HumanOutcome|DecisionSatisfied)ConditionKeysOnConditionNotRaisingCode|TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen' ./cmd/fid/integration/... -v` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	22.659s` (terminated `exit_code: 0`). Its three named tests invoke the built `fid` process and exercise `stop raise`, `stop resolve`, `human-outcome record`, and `decision record-reference`.
- Timing diagnostic at current HEAD: `go test -race -timeout 30m ./cmd/fid -run 'TestStop' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	112.324s` (terminated `exit_code: 0`).
- Same timing diagnostic at the pre-repair parent `bc7e6f393450780137f3d89ad9763e382028ccd3` in temporary worktree `/tmp/fid-task11-base.LMrrHe` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	109.999s` (terminated `exit_code: 0`).
- Exact ordinary gate at that parent, for attribution only: `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	238.620s	coverage: 64.0% of statements` (terminated `exit_code: 0`); `cmd/fid` was 593.812s. This does not replace the task-11 gate.
- Historical original-task red/green evidence: the earlier report recorded `currentConditionShapes has 8 rows, want 10` before the original implementation and green focused tests after it; review-7 BC-3 records the later DEC-250 red observation above. Neither is presented as a new termination observed in this continuation.

## Exact quality-gate results

Working directory for these five ordered commands: `/home/mark/projects/fid`.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - FAIL - `ok  	github.com/mbell-nzl/fid/internal/workflow	277.603s	coverage: 64.6% of statements` (terminated `exit_code: 0`). The same output reports `cmd/fid` 686.386s, above 600s; all tests themselves passed.
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `ok  	github.com/mbell-nzl/fid/internal/workflow	277.603s	coverage: 64.6% of statements` (terminated `exit_code: 0`). This row records pending budget status from the prior attempt, not a second execution; task-11 speed attribution and repair remain open.
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.4%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	778.278s` (terminated `exit_code: 0`)

DEC-227 times from this continuation: `cmd/fid/integration` 778.278s (<900s), `cmd/fid` 686.386s (>600s), `internal/workflow` 277.603s (<600s). The parent comparison had `cmd/fid` 593.812s and `internal/workflow` 238.620s. Other unrelated packages were also slower in the task-11 attempt (`internal/export` 387.634s versus 329.302s; `internal/journal` 295.151s versus 250.583s). These observations suggest broad timing variation, but do not establish a causal repair or pass the budget. The earlier task-11 report's 546.059s/213.335s/700.168s run appended `-count=1` and remains historical evidence, not the current exact gate.

## Deviations from spec

- Report-lint first attempt: `FAIL report-lint found 4 failure(s)` (terminated `exit_code: 4`). The recovery-attempt value contained trigger text, and the two ordinary-gate status fields contained reasons rather than the exact `FAIL` and `BLOCKED` values. I kept the trigger in the summary and the reasons after the output tokens, then reran the byte-identical report-lint command: `PASS report-lint completed with no failures` (terminated `exit_code: 0`), with `OUTCOME RECOVERY_REQUIRED`.
- No command, assertion, threshold, or budget was changed in this recovery. The original report's `-count=1` command substitution is retained as historical deviation; this continuation ran all five task-spec commands byte-for-byte.
- The first exact ordinary gate's last line was `ok  	github.com/mbell-nzl/fid/internal/workflow	277.603s	coverage: 64.6% of statements` (terminated `exit_code: 0`), but its `cmd/fid` package time breached the 600s budget. Root cause remains unproven; no repair or later passing task-11 gate is claimed.
- `internal/workflow/stop_conditions.go` was ratified into task-11 scope by DEC-253. The retained `p23_attention_resolution_test.go` test was introduced by `c5bf0f89c69c57003ece353ef005a7596dd266e4` (verified with `git log --diff-filter=A`); `466745f2dc4505d60eb67ac50193e004532b7239` preserved its original wrong-class/wrong-result subjects and strengthened the same row with explicit reject/no-Event/open-Stop and matching operator-Stop assertions. No retained assertion was narrowed.

## Concerns, technical debt, or notes for the Tech Lead

- The exact ordinary gate is not budget-green. A parent-commit run was 593.812s for `cmd/fid`, only 6.188s below the limit, while the task-11 run was 686.386s. The focused Stop comparison was 109.999s versus 112.324s and unrelated packages slowed too. A continuation must establish cause and repair under task scope or return the exact scope gap for authority; a cache hit or unmodified rerun alone cannot erase the failed attempt.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `15399a45e9c47a96af838e62aed3083e43b5967f`

## Commit subject

- `docs(task): checkpoint phase-23 task-11 gate recovery`

OUTCOME RECOVERY_REQUIRED
