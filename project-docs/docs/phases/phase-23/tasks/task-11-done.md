# Phase 23 - Task 11: Attention resolution through Human Outcomes - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 4
- Next action: None.
- Human attention required: no

## Summary of what was implemented

The original task-11 commit `c5bf0f89c69c57003ece353ef005a7596dd266e4` completed the four judgment Condition encodings and the `review-recorded` and `finding-resolved` resolution paths. The accepted DEC-250 repair commit `466745f2dc4505d60eb67ac50193e004532b7239` adds one shared qualifying-fact predicate for each of `human-outcome-recorded` and `decision-satisfied`. Manual `stop resolve` and atomic resolution in `human-outcome.recorded` and `decision.referenced` use the same predicates. Resolution scans the recorded Condition, so producer-raised Stops and `operator-stop` Stops qualify alike. A `decision-satisfied` Condition on an activation subject now rejects through the existing invalid-Condition path with no Event. No contract, schema, Event, failure code, or public-interface shape changed.

The review-7 BC-3 red observation remains preserved at `qa/checkpoint-after-task-11-2.md:744-755`: a matching approved Outcome resolved the readiness Stop but left the matching `operator-stop` open, and manual resolution rejected `EVT-20260925-000026` with `resolution-invalid` and Event count `26->26`. The current focused built-binary tests are green (Tests run). This report replaces the earlier COMPLETE claim, whose pre-repair implementation was `c5bf0f8` and whose recorded Go gates appended `-count=1`.

ESC-23-005 closure (DEC-258, Human-delegated repair above role boundaries): the budget breach was not a task-11 cause. It came from serial test structure plus measurement spread (`qa/esc-23-005-repair.md`). Repair commit `3dff4574a71438900100449e36d7044969ba9162` adds `t.Parallel()` to existing serial tests and enables `paralleltest`, and changes no assertion, command, or budget. The five exact commands then passed on that commit, with `cmd/fid` at 531.533s. The earlier failed budget attempts below remain failed evidence.

## Files created or modified

- Original task: `internal/workflow/stop.go`, `stop_conditions.go`, `stop_test.go`; `internal/journal/current_event_value_catalog.go` and `current_event_task11_values_test.go`; `cmd/fid/integration/p23_attention_resolution_test.go`.
- DEC-250 repair: `internal/workflow/stop.go`, `stop_conditions.go`, `stop_judgment_conditions_test.go`, `evaluation.go`, `resolve.go`; `internal/review/human_outcome.go` and `decision.go`; `cmd/fid/integration/p23_attention_resolution_test.go` and `p23_stop_resolution_conditions_test.go`.
- This attempt-2 continuation: `project-docs/docs/phases/phase-23/tasks/task-11-done.md` only; no product or test change.
- DEC-258 repair (outside task-11's named scope, authorized by DEC-258): test files across `cmd/fid`, `cmd/fid/integration`, and `internal/*` (`t.Parallel()` only, plus serial annotations), and `.golangci.yml`. It also updates `AGENTS.md`, `09_AI_WORKFLOW.md`, the Developer prompt and its generated copies, `decisions.md` (DEC-258), and `tasks/task-13.md`. No product source changed.

## Documentation updated

- This report records DEC-250/253 repair evidence, both exact task-11 budget failures, and the same-session parent timing comparison.
- DEC-258 documentation: `qa/esc-23-005-repair.md` and its `esc-23-005-evidence/` gate logs, the tiered test loop in `09_AI_WORKFLOW.md` §2.5 step 6 and the Developer prompt, and the `AGENTS.md` test-style rule.

## Validation Contract Coverage

- VAL-047 (task-11 half): `TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess` checks `attention.v1` reports `human_outcome_resolution:false` at `cmd/fid/integration/p23_attention_resolution_test.go:118`; the built-binary resolution tests below establish the path before task 17 flips the limit.
- VAL-048: `TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen` at `cmd/fid/integration/p23_attention_resolution_test.go:169-207` records wrong-class and wrong-result Outcomes, then `assertP23StopResolveRejected` checks exit 4 `resolution-invalid`, unchanged Event count, and open producer and operator Stops. The approved matching Outcome resolves both.
- VAL-049: `TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess` at `cmd/fid/integration/p23_attention_resolution_test.go:83-114` rejects an extra `outcome_id` input with no Event, accepts the existing `attention_id`/`expected_version`/`fact_event_ids` shape, and re-reads the original Stop identity and exact resolution fact IDs.
- VAL-050: `TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen` at `cmd/fid/integration/p23_attention_resolution_test.go:156-208` resolves the producer Stop and a matching-Condition `operator-stop` in the same Outcome Event; the different `review-recorded` Stop remains open.
- VAL-078: `TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess` at `cmd/fid/integration/p23_attention_resolution_test.go:52-81` covers non-qualifying and qualifying Review facts for `review-recorded` and `finding-resolved`. `TestP23StopResolveHumanOutcomeConditionKeysOnConditionNotRaisingCode` at `cmd/fid/integration/p23_stop_resolution_conditions_test.go:26-128` covers non-qualifying and qualifying Outcome facts on producer and operator Stops, including envelope, run, and activation subjects. `TestP23StopResolveDecisionSatisfiedConditionKeysOnConditionNotRaisingCode` at `cmd/fid/integration/p23_stop_resolution_conditions_test.go:130-228` covers rejected/deferred/earlier-half facts and both orders of the qualifying pair on producer and operator Stops, plus the activation-subject raise rejection with no Event. `TestHumanOutcomeConditionQualifyingFactRules` and `TestDecisionSatisfiedConditionQualifyingFactRules` at `internal/workflow/stop_judgment_conditions_test.go:75` and `:129` cover scope, supersession, and pair semantics.
- VAL-080: `TestTask11ReviewRecordedAndFindingResolvedConditionShapes` at `internal/journal/current_event_task11_values_test.go:12` covers canonical, extra, wrong, missing, and null discriminator shapes; `TestTask11ConditionShapeRegistryIsClosedAtTenRows` at `:67` checks the ten-row registry. The original task-11 evidence for the task-9/10 types remains in their committed tests. Malformed shapes reject through shape validation; no relation rule was added.
- VAL-081: `TestTask11FindingResolvedConditionDispatchPreservesNullableBinding` at `internal/journal/current_event_task11_values_test.go:52-62` observes the exact error `condition nullable binding mismatch` for a `finding-resolved` Condition with stray `budget_name` after task-9 dispatch. `TestTask9ConditionDispatchPreservesNullableBindings` at `internal/journal/current_event_closed_values_test.go:140` retains the three Phase 22 negatives and canonical positives. The relation was unchanged.
- VAL-079: PASS on the final exact gates at `3dff457` (DEC-258 attempt 2): `cmd/fid` `531.533s` (<600s), `internal/workflow` `165.920s` (<600s), and tagged `cmd/fid/integration` `326.421s` (<900s). None terminated by timeout. Earlier failures are retained. Task-11 attempts measured `cmd/fid` at `686.386s` and `630.927s` (FAIL), and the pre-repair parent at `614.601s`. DEC-258 attempt 1 failed on the midnight journal-file flake recorded below. `cmd/fid` sits at 88.6% of its budget, which triggers DEC-258's headroom warning (see Concerns).

## Tests run

- Focused built-binary green: `go test -tags integration -race -timeout 5m -run 'TestP23StopResolve(HumanOutcome|DecisionSatisfied)ConditionKeysOnConditionNotRaisingCode|TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen' ./cmd/fid/integration/... -v` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	22.659s` (terminated `exit_code: 0`). Its three named tests invoke the built `fid` process and exercise `stop raise`, `stop resolve`, `human-outcome record`, and `decision record-reference`.
- Timing diagnostic at current HEAD: `go test -race -timeout 30m ./cmd/fid -run 'TestStop' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	112.324s` (terminated `exit_code: 0`).
- Same timing diagnostic at the pre-repair parent `bc7e6f393450780137f3d89ad9763e382028ccd3` in temporary worktree `/tmp/fid-task11-base.LMrrHe` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	109.999s` (terminated `exit_code: 0`).
- Exact ordinary gate at that parent, for attribution only: `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	238.620s	coverage: 64.0% of statements` (terminated `exit_code: 0`); `cmd/fid` was 593.812s. This does not replace the task-11 gate.
- Historical original-task red/green evidence: the earlier report recorded `currentConditionShapes has 8 rows, want 10` before the original implementation and green focused tests after it; review-7 BC-3 records the later DEC-250 red observation above. Neither is presented as a new termination observed in this continuation.
- Attempt-2 parent attribution diagnostic, working directory `/tmp/fid-task11-parent-attribution.N4P6ls/worktree` at `bc7e6f393450780137f3d89ad9763e382028ccd3`, after clearing the Go test-result cache: `go test -race -timeout 30m -coverprofile=coverage.out ./...` - FAIL on the same 600s package budget - `ok  	github.com/mbell-nzl/fid/internal/workflow	247.706s	coverage: 64.0% of statements` (terminated `exit_code: 0`); `cmd/fid` was 614.601s. This diagnostic does not replace a task-11 gate.
- Task-11 production diff from that parent changed `internal/workflow` and `internal/review` plus two `cmd/fid/integration` test files; it added no ordinary `cmd/fid` test. The previous focused `TestStop` parent/current times were 109.999s/112.324s. The new full parent/current comparison is 614.601s/630.927s for `cmd/fid`, 247.706s/252.947s for `internal/workflow`, 333.871s/345.281s for `internal/export`, and 251.922s/262.567s for `internal/journal`. The parent itself breached and unrelated packages rose in the same direction; the task-11 share of the 30.927s current overage remains unproven.

## Exact quality-gate results

### Final DEC-258 verification (repair commit `3dff4574a71438900100449e36d7044969ba9162`)

Working directory: `/home/mark/projects/fid`. Each attempt ran the five commands in order after `go clean -testcache`. Full logs are in `qa/esc-23-005-evidence/exact-gates-attempt-1.log` and `exact-gates-attempt-2.log`.

- DEC-258 attempt 1 `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- DEC-258 attempt 1 `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- DEC-258 attempt 1 `go test -race -timeout 30m -coverprofile=coverage.out ./...` - FAIL - `FAIL	github.com/mbell-nzl/fid/cmd/fid	525.379s` (terminated `exit_code: 1`). All 14 failures were `journal files = 2, want 1` from `mustJournalPath`. The run crossed 00:00 UTC, and journal files are per UTC day. This is a pre-existing midnight fixture flake, independent of the repair.
- DEC-258 attempt 1 `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.4%` (terminated `exit_code: 0`)
- DEC-258 attempt 1 `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	335.029s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	165.920s	coverage: 64.6% of statements` (terminated `exit_code: 0`)
- The same attempt-2 ordinary output reports `cmd/fid` at 531.533s.
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.4%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	326.421s` (terminated `exit_code: 0`)

Final DEC-227 times: `cmd/fid/integration` 326.421s (<900s), `cmd/fid` 531.533s (<600s), `internal/workflow` 165.920s (<600s).

### Earlier task-11 attempts (retained)

Working directory for these five ordered commands: `/home/mark/projects/fid`.

- Task-11 attempt 1 `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- Task-11 attempt 1 `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- Task-11 attempt 1 `go test -race -timeout 30m -coverprofile=coverage.out ./...` - FAIL - `ok  	github.com/mbell-nzl/fid/internal/workflow	277.603s	coverage: 64.6% of statements` (terminated `exit_code: 0`). The same output reports `cmd/fid` 686.386s, above 600s; all tests themselves passed.
- Attempt 1's ordinary gate remains failed on budget; its `BLOCKED` disposition was pending attribution, not a second execution.
- Task-11 attempt 1 `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.4%` (terminated `exit_code: 0`)
- Task-11 attempt 1 `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	778.278s` (terminated `exit_code: 0`)

Attempt 2 ran the five exact task-spec commands in order at `/home/mark/projects/fid`. Its first ordinary run was cached. After the five-command sequence, `go clean -testcache` invalidated only Go test results and the byte-identical ordinary command ran again for a measured package time.

- Attempt 2 `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- Attempt 2 `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- Attempt 2 `go test -race -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 64.6% of statements` (terminated `exit_code: 0`). Every package was cached, so the run supplied no new wall-time evidence.
- Attempt 2 `go test -race -timeout 30m -coverprofile=coverage.out ./...` - FAIL - `ok  	github.com/mbell-nzl/fid/internal/workflow	252.947s	coverage: 64.6% of statements` (terminated `exit_code: 0`). Fresh `cmd/fid` time: 630.927s, over 600s.
- Task-11 attempt 2 disposition `go test -race -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `ok  	github.com/mbell-nzl/fid/internal/workflow	252.947s	coverage: 64.6% of statements` (terminated `exit_code: 0`). This is the unresolved budget disposition of the immediately preceding observed attempt, not another execution; the 630.927s `cmd/fid` result remains FAIL.
- Attempt 2 `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.4%` (terminated `exit_code: 0`)
- Attempt 2 `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	766.070s` (terminated `exit_code: 0`)

DEC-227 times from this continuation: `cmd/fid/integration` 778.278s (<900s), `cmd/fid` 686.386s (>600s), `internal/workflow` 277.603s (<600s). The parent comparison had `cmd/fid` 593.812s and `internal/workflow` 238.620s. Other unrelated packages were also slower in the task-11 attempt (`internal/export` 387.634s versus 329.302s; `internal/journal` 295.151s versus 250.583s). These observations suggest broad timing variation, but do not establish a causal repair or pass the budget. The earlier task-11 report's 546.059s/213.335s/700.168s run appended `-count=1` and remains historical evidence, not the current exact gate.

Attempt-2 exact task-11 times: `cmd/fid/integration` 766.070s (<900s), `cmd/fid` 630.927s (>600s), `internal/workflow` 252.947s (<600s). The fresh same-session pre-repair parent took 614.601s for `cmd/fid`, itself above budget. The first attempt-2 ordinary result was cached and had no measured package time. The non-cached attempt still fails VAL-079; an unmodified rerun does not erase the earlier failed attempt.

## Deviations from spec

- DEC-258: the delegated repair edits files outside task-11's named scope. It adds only `t.Parallel()` and serial annotations to tests, plus lint config and workflow docs. No assertion, command, budget, or product source changed.
- `TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat` was already parallel. It failed once during repair measurement (`maintenance warnings differs across fixed-grace branches`) because its 299/300-second fixture ages against the wall clock across two command runs. It now runs serial with its assertions unchanged. The deterministic boundary is covered in `internal/maintenance`.

- Report-lint first attempt: `FAIL report-lint found 4 failure(s)` (terminated `exit_code: 4`). The recovery-attempt value contained trigger text, and the two ordinary-gate status fields contained reasons rather than the exact `FAIL` and `BLOCKED` values. I kept the trigger in the summary and the reasons after the output tokens, then reran the byte-identical report-lint command: `PASS report-lint completed with no failures` (terminated `exit_code: 0`), with `OUTCOME RECOVERY_REQUIRED`.
- No command, assertion, threshold, or budget was changed in this recovery. The original report's `-count=1` command substitution is retained as historical deviation; this continuation ran all five task-spec commands byte-for-byte.
- The first exact ordinary gate's last line was `ok  	github.com/mbell-nzl/fid/internal/workflow	277.603s	coverage: 64.6% of statements` (terminated `exit_code: 0`), but its `cmd/fid` package time breached the 600s budget. Root cause remains unproven; no repair or later passing task-11 gate is claimed.
- Attempt 2 repeated the exact ordinary command after clearing the test cache and observed `cmd/fid` at 630.927s (terminated `exit_code: 0`), still over budget. The same-session pre-repair parent was 614.601s (terminated `exit_code: 0`), also over budget. The 16.326s current-minus-parent difference is smaller than the current 30.927s overage; unrelated package increases and the prior focused Stop comparison prevent attributing the breach to task 11. No test, command, threshold, or source file was altered to seek a green timing result.
- Attempt-2 report-lint first result: `FAIL report-lint found 1 failure(s)` (terminated `exit_code: 1`) because the final required ordinary-command row was `FAIL` without its unresolved `BLOCKED` disposition. The observed FAIL attempt remains intact; this report correction adds the required disposition without changing a gate or its evidence. The second report-lint attempt also returned `FAIL report-lint found 1 failure(s)` (terminated `exit_code: 1`): its parser did not recognize the `Attempt 2` prefix on that disposition row. The disposition now uses the exact command at the row start.
- `internal/workflow/stop_conditions.go` was ratified into task-11 scope by DEC-253. The retained `p23_attention_resolution_test.go` test was introduced by `c5bf0f89c69c57003ece353ef005a7596dd266e4` (verified with `git log --diff-filter=A`); `466745f2dc4505d60eb67ac50193e004532b7239` preserved its original wrong-class/wrong-result subjects and strengthened the same row with explicit reject/no-Event/open-Stop and matching operator-Stop assertions. No retained assertion was narrowed.

## Concerns, technical debt, or notes for the Tech Lead

- DEC-258 headroom warning: `cmd/fid` recorded 531.533s, 88.6% of its 600s budget and above the 480s warning line. With every package parallel, the ordinary gate is CPU-throughput bound, so package times now track total suite CPU. A profile attributes 33% of `cmd/fid` CPU to SQLite re-parsing the schema on each connection open (several opens per command, mainly `journal.ensureProjectionFresh` and `workflow.openWorkflowDB`), and 32% to race-implied `checkptr` inside modernc SQLite. Tasks 13 and 14 add `cmd/fid` tests. Schedule a speed item first; the options and authority needed are in `qa/esc-23-005-repair.md`.
- Carry-forward: `mustJournalPath` (`cmd/fid/init_test.go:698`) assumes one journal file, but journals are split per UTC day. Any gate that crosses 00:00 UTC (13:00 NZDT) can fail with `journal files = 2, want 1`. This is pre-existing and was not changed here.

- The exact ordinary gate remains over budget. A same-session parent run also breached at 614.601s; task 11's accepted changes add no ordinary `cmd/fid` test file, and the measured Stop subset rose only 2.325s in the prior diagnostic. No task-11 causal speed repair can be justified from this evidence. The remaining speed work lies outside the named task-11 source/test scope unless a future diagnostic proves a task-11 path; the Phase Runner needs bounded scope disposition. This report does not assign that work to another task.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `3dff4574a71438900100449e36d7044969ba9162`

## Commit subject

- `perf(test): parallelize serial test suites for ESC-23-005`

OUTCOME COMPLETE
