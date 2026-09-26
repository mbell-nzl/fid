# Phase 23 - Task 11: Attention resolution through Human Outcomes - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: gate-repair
- Condition: The exact ordinary race/coverage command exited 0 twice but `cmd/fid` took 686.386s and then 630.927s, both above the accepted 600s package budget. A fresh same-session run at task 11's pre-repair parent also breached at 614.601s. No task-11 causal speed repair is established; this is a scope gap for the still-failed budget gate.
- Recovery attempts: 2
- Next action: The Phase Runner obtains bounded Tech Lead disposition of the `cmd/fid` speed scope gap while retaining all three budget failures. A same-task continuation may repair only a demonstrated task-11 cause under its accepted file scope and must rerun the unchanged five-command gate before closure.
- Human attention required: no

## Summary of what was implemented

The original task-11 commit `c5bf0f89c69c57003ece353ef005a7596dd266e4` completed the four judgment Condition encodings and the `review-recorded` and `finding-resolved` resolution paths. The accepted DEC-250 repair commit `466745f2dc4505d60eb67ac50193e004532b7239` adds one shared qualifying-fact predicate for each of `human-outcome-recorded` and `decision-satisfied`. Manual `stop resolve` and atomic resolution in `human-outcome.recorded` and `decision.referenced` use the same predicates. Resolution scans the recorded Condition, so producer-raised Stops and `operator-stop` Stops qualify alike. A `decision-satisfied` Condition on an activation subject now rejects through the existing invalid-Condition path with no Event. No contract, schema, Event, failure code, or public-interface shape changed.

The review-7 BC-3 red observation remains preserved at `qa/checkpoint-after-task-11-2.md:744-755`: a matching approved Outcome resolved the readiness Stop but left the matching `operator-stop` open, and manual resolution rejected `EVT-20260925-000026` with `resolution-invalid` and Event count `26->26`. The current focused built-binary tests are green (Tests run). This report replaces the earlier COMPLETE claim, whose pre-repair implementation was `c5bf0f8` and whose recorded Go gates appended `-count=1`.

## Files created or modified

- Original task: `internal/workflow/stop.go`, `stop_conditions.go`, `stop_test.go`; `internal/journal/current_event_value_catalog.go` and `current_event_task11_values_test.go`; `cmd/fid/integration/p23_attention_resolution_test.go`.
- DEC-250 repair: `internal/workflow/stop.go`, `stop_conditions.go`, `stop_judgment_conditions_test.go`, `evaluation.go`, `resolve.go`; `internal/review/human_outcome.go` and `decision.go`; `cmd/fid/integration/p23_attention_resolution_test.go` and `p23_stop_resolution_conditions_test.go`.
- This attempt-2 continuation: `project-docs/docs/phases/phase-23/tasks/task-11-done.md` only; no product or test change.

## Documentation updated

- This report records DEC-250/253 repair evidence, both exact task-11 budget failures, and the same-session parent timing comparison. Accepted specifications and QA artifacts remain unchanged.

## Validation Contract Coverage

- VAL-047 (task-11 half): `TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess` checks `attention.v1` reports `human_outcome_resolution:false` at `cmd/fid/integration/p23_attention_resolution_test.go:118`; the built-binary resolution tests below establish the path before task 17 flips the limit.
- VAL-048: `TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen` at `cmd/fid/integration/p23_attention_resolution_test.go:169-207` records wrong-class and wrong-result Outcomes, then `assertP23StopResolveRejected` checks exit 4 `resolution-invalid`, unchanged Event count, and open producer and operator Stops. The approved matching Outcome resolves both.
- VAL-049: `TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess` at `cmd/fid/integration/p23_attention_resolution_test.go:83-114` rejects an extra `outcome_id` input with no Event, accepts the existing `attention_id`/`expected_version`/`fact_event_ids` shape, and re-reads the original Stop identity and exact resolution fact IDs.
- VAL-050: `TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen` at `cmd/fid/integration/p23_attention_resolution_test.go:156-208` resolves the producer Stop and a matching-Condition `operator-stop` in the same Outcome Event; the different `review-recorded` Stop remains open.
- VAL-078: `TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess` at `cmd/fid/integration/p23_attention_resolution_test.go:52-81` covers non-qualifying and qualifying Review facts for `review-recorded` and `finding-resolved`. `TestP23StopResolveHumanOutcomeConditionKeysOnConditionNotRaisingCode` at `cmd/fid/integration/p23_stop_resolution_conditions_test.go:26-128` covers non-qualifying and qualifying Outcome facts on producer and operator Stops, including envelope, run, and activation subjects. `TestP23StopResolveDecisionSatisfiedConditionKeysOnConditionNotRaisingCode` at `cmd/fid/integration/p23_stop_resolution_conditions_test.go:130-228` covers rejected/deferred/earlier-half facts and both orders of the qualifying pair on producer and operator Stops, plus the activation-subject raise rejection with no Event. `TestHumanOutcomeConditionQualifyingFactRules` and `TestDecisionSatisfiedConditionQualifyingFactRules` at `internal/workflow/stop_judgment_conditions_test.go:75` and `:129` cover scope, supersession, and pair semantics.
- VAL-080: `TestTask11ReviewRecordedAndFindingResolvedConditionShapes` at `internal/journal/current_event_task11_values_test.go:12` covers canonical, extra, wrong, missing, and null discriminator shapes; `TestTask11ConditionShapeRegistryIsClosedAtTenRows` at `:67` checks the ten-row registry. The original task-11 evidence for the task-9/10 types remains in their committed tests. Malformed shapes reject through shape validation; no relation rule was added.
- VAL-081: `TestTask11FindingResolvedConditionDispatchPreservesNullableBinding` at `internal/journal/current_event_task11_values_test.go:52-62` observes the exact error `condition nullable binding mismatch` for a `finding-resolved` Condition with stray `budget_name` after task-9 dispatch. `TestTask9ConditionDispatchPreservesNullableBindings` at `internal/journal/current_event_closed_values_test.go:140` retains the three Phase 22 negatives and canonical positives. The relation was unchanged.
- VAL-079: FAIL. The exact ordinary race/coverage command's `cmd/fid` package line was `630.927s` in this attempt, exceeding 600s; the earlier task-11 attempt was `686.386s`. The exact tagged integration gate took `766.070s`, under 900s. A fresh pre-repair parent run took `614.601s` for `cmd/fid`, so this attempt does not establish a task-11 causal repair.

## Tests run

- Focused built-binary green: `go test -tags integration -race -timeout 5m -run 'TestP23StopResolve(HumanOutcome|DecisionSatisfied)ConditionKeysOnConditionNotRaisingCode|TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen' ./cmd/fid/integration/... -v` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	22.659s` (terminated `exit_code: 0`). Its three named tests invoke the built `fid` process and exercise `stop raise`, `stop resolve`, `human-outcome record`, and `decision record-reference`.
- Timing diagnostic at current HEAD: `go test -race -timeout 30m ./cmd/fid -run 'TestStop' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	112.324s` (terminated `exit_code: 0`).
- Same timing diagnostic at the pre-repair parent `bc7e6f393450780137f3d89ad9763e382028ccd3` in temporary worktree `/tmp/fid-task11-base.LMrrHe` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	109.999s` (terminated `exit_code: 0`).
- Exact ordinary gate at that parent, for attribution only: `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	238.620s	coverage: 64.0% of statements` (terminated `exit_code: 0`); `cmd/fid` was 593.812s. This does not replace the task-11 gate.
- Historical original-task red/green evidence: the earlier report recorded `currentConditionShapes has 8 rows, want 10` before the original implementation and green focused tests after it; review-7 BC-3 records the later DEC-250 red observation above. Neither is presented as a new termination observed in this continuation.
- Attempt-2 parent attribution diagnostic, working directory `/tmp/fid-task11-parent-attribution.N4P6ls/worktree` at `bc7e6f393450780137f3d89ad9763e382028ccd3`, after clearing the Go test-result cache: `go test -race -timeout 30m -coverprofile=coverage.out ./...` - FAIL on the same 600s package budget - `ok  	github.com/mbell-nzl/fid/internal/workflow	247.706s	coverage: 64.0% of statements` (terminated `exit_code: 0`); `cmd/fid` was 614.601s. This diagnostic does not replace a task-11 gate.
- Task-11 production diff from that parent changed `internal/workflow` and `internal/review` plus two `cmd/fid/integration` test files; it added no ordinary `cmd/fid` test. The previous focused `TestStop` parent/current times were 109.999s/112.324s. The new full parent/current comparison is 614.601s/630.927s for `cmd/fid`, 247.706s/252.947s for `internal/workflow`, 333.871s/345.281s for `internal/export`, and 251.922s/262.567s for `internal/journal`. The parent itself breached and unrelated packages rose in the same direction; the task-11 share of the 30.927s current overage remains unproven.

## Exact quality-gate results

Working directory for these five ordered commands: `/home/mark/projects/fid`.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - FAIL - `ok  	github.com/mbell-nzl/fid/internal/workflow	277.603s	coverage: 64.6% of statements` (terminated `exit_code: 0`). The same output reports `cmd/fid` 686.386s, above 600s; all tests themselves passed.
- Attempt 1's ordinary gate remains failed on budget; its `BLOCKED` disposition was pending attribution, not a second execution.
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.4%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	778.278s` (terminated `exit_code: 0`)

Attempt 2 ran the five exact task-spec commands in order at `/home/mark/projects/fid`. Its first ordinary run was cached. After the five-command sequence, `go clean -testcache` invalidated only Go test results and the byte-identical ordinary command ran again for a measured package time.

- Attempt 2 `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- Attempt 2 `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- Attempt 2 `go test -race -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 64.6% of statements` (terminated `exit_code: 0`). Every package was cached, so the run supplied no new wall-time evidence.
- Attempt 2 `go test -race -timeout 30m -coverprofile=coverage.out ./...` - FAIL - `ok  	github.com/mbell-nzl/fid/internal/workflow	252.947s	coverage: 64.6% of statements` (terminated `exit_code: 0`). Fresh `cmd/fid` time: 630.927s, over 600s.
- Attempt 2 `go test -race -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `ok  	github.com/mbell-nzl/fid/internal/workflow	252.947s	coverage: 64.6% of statements` (terminated `exit_code: 0`). This is the unresolved budget disposition of the immediately preceding observed attempt, not another execution; the 630.927s `cmd/fid` result remains FAIL.
- Attempt 2 `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.4%` (terminated `exit_code: 0`)
- Attempt 2 `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	766.070s` (terminated `exit_code: 0`)

DEC-227 times from this continuation: `cmd/fid/integration` 778.278s (<900s), `cmd/fid` 686.386s (>600s), `internal/workflow` 277.603s (<600s). The parent comparison had `cmd/fid` 593.812s and `internal/workflow` 238.620s. Other unrelated packages were also slower in the task-11 attempt (`internal/export` 387.634s versus 329.302s; `internal/journal` 295.151s versus 250.583s). These observations suggest broad timing variation, but do not establish a causal repair or pass the budget. The earlier task-11 report's 546.059s/213.335s/700.168s run appended `-count=1` and remains historical evidence, not the current exact gate.

Attempt-2 exact task-11 times: `cmd/fid/integration` 766.070s (<900s), `cmd/fid` 630.927s (>600s), `internal/workflow` 252.947s (<600s). The fresh same-session pre-repair parent took 614.601s for `cmd/fid`, itself above budget. The first attempt-2 ordinary result was cached and had no measured package time. The non-cached attempt still fails VAL-079; an unmodified rerun does not erase the earlier failed attempt.

## Deviations from spec

- Report-lint first attempt: `FAIL report-lint found 4 failure(s)` (terminated `exit_code: 4`). The recovery-attempt value contained trigger text, and the two ordinary-gate status fields contained reasons rather than the exact `FAIL` and `BLOCKED` values. I kept the trigger in the summary and the reasons after the output tokens, then reran the byte-identical report-lint command: `PASS report-lint completed with no failures` (terminated `exit_code: 0`), with `OUTCOME RECOVERY_REQUIRED`.
- No command, assertion, threshold, or budget was changed in this recovery. The original report's `-count=1` command substitution is retained as historical deviation; this continuation ran all five task-spec commands byte-for-byte.
- The first exact ordinary gate's last line was `ok  	github.com/mbell-nzl/fid/internal/workflow	277.603s	coverage: 64.6% of statements` (terminated `exit_code: 0`), but its `cmd/fid` package time breached the 600s budget. Root cause remains unproven; no repair or later passing task-11 gate is claimed.
- Attempt 2 repeated the exact ordinary command after clearing the test cache and observed `cmd/fid` at 630.927s (terminated `exit_code: 0`), still over budget. The same-session pre-repair parent was 614.601s (terminated `exit_code: 0`), also over budget. The 16.326s current-minus-parent difference is smaller than the current 30.927s overage; unrelated package increases and the prior focused Stop comparison prevent attributing the breach to task 11. No test, command, threshold, or source file was altered to seek a green timing result.
- Attempt-2 report-lint first result: `FAIL report-lint found 1 failure(s)` (terminated `exit_code: 1`) because the final required ordinary-command row was `FAIL` without its unresolved `BLOCKED` disposition. The observed FAIL attempt remains intact; this report correction adds the required disposition without changing a gate or its evidence.
- `internal/workflow/stop_conditions.go` was ratified into task-11 scope by DEC-253. The retained `p23_attention_resolution_test.go` test was introduced by `c5bf0f89c69c57003ece353ef005a7596dd266e4` (verified with `git log --diff-filter=A`); `466745f2dc4505d60eb67ac50193e004532b7239` preserved its original wrong-class/wrong-result subjects and strengthened the same row with explicit reject/no-Event/open-Stop and matching operator-Stop assertions. No retained assertion was narrowed.

## Concerns, technical debt, or notes for the Tech Lead

- The exact ordinary gate remains over budget. A same-session parent run also breached at 614.601s; task 11's accepted changes add no ordinary `cmd/fid` test file, and the measured Stop subset rose only 2.325s in the prior diagnostic. No task-11 causal speed repair can be justified from this evidence. The remaining speed work lies outside the named task-11 source/test scope unless a future diagnostic proves a task-11 path; the Phase Runner needs bounded scope disposition. This report does not assign that work to another task.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `0061dde96f7e1c9809994f71abdbd3e3a09558c4`

## Commit subject

- `docs(task): checkpoint phase-23 task-11 timing scope gap`

OUTCOME RECOVERY_REQUIRED
