# Phase 23 - Task 5: Finding dispositions and resolution authority - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added atomic, immutable Finding creation inside `review.record`, an active-mutation-Claim-only `finding address` command, and Finding list/read projection with derived `open`, `addressed`, `resolved`, and `inapplicable` states.
- Added later-Review resolution and reasoned inapplicability under the recorded independence and verdict rules. The implementation rejects creating-Review self-resolution, a mutating Claim holder, blocked verdicts, duplicate terminal dispositions, foreign-scope targets, malformed collections, and unknown Findings without appending a partial Event.
- Preserved DEC-225 behavior: a different independent Actor and the same independent Reviewer in a later Review may both make terminal dispositions. Mixed resolved and inapplicable Findings are committed atomically, while `Review.resolved_finding_ids` remains resolved-only.
- Added strict DEC-242 input handling: omission and `[]` canonicalize alike, explicit `null` and duplicate or unknown keys reject, reasons remain verbatim and are bounded to nonblank valid UTF-8 of at most 4096 bytes, and at most 64 unique/disjoint terminal targets are accepted.
- Converged only the `finding.addressed` fence key, registered its existing typed result/projection shape, and added `finding.address` to both sorted closed action allowlists. Task 4's `Finding` and `FindingDisposition` value shapes remain unchanged.
- Added in-process and built-binary coverage for address, terminal disposition authority, input boundaries, strict projection and replay parity, derived state, exact error details, idempotent replay, and changed-reason conflict.

## Files created or modified

- `cmd/fid/finding.go`, `cmd/fid/finding_test.go`, `cmd/fid/review.go`, `cmd/fid/review_test.go`, `cmd/fid/root.go`, `cmd/fid/envelope_test.go`, and `cmd/fid/contract_surface_test.go`
- `cmd/fid/integration/p22_fake_test.go` and `cmd/fid/integration/p23_finding_disposition_test.go`
- `internal/review/finding.go`, `internal/review/disposition.go`, `internal/review/review.go`, and `internal/review/review_test.go`
- `internal/cli/codes.go` and `internal/cli/codes_test.go`
- `internal/journal/current_event_json.go`, `internal/journal/current_event_result_projection.go`, `internal/journal/current_event_result_values.go`, `internal/journal/current_event_value_relations.go`, `internal/journal/project_effects.go`, and `internal/journal/project_effects_test.go`
- `internal/workflow/activation.go`
- `project-docs/docs/phases/phase-23/tasks/task-5-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-23/tasks/task-5-done.md`

## Validation Contract Coverage

- VAL-003: PASS - the `finding.addressed` fence originated in `4aa23200789e0ce65172d6e42c49704849854a4a`. `TestFindingAddressedFenceConvergencePreservesEveryRemainingEvent` fixes the remaining fence inventory, `TestProjectEventFenceRejectsEveryUnavailablePhase23EventWithoutProjectionChange` preserves rejection parity for every retained key, and `TestFindingAddressDoesNotResolveOrUnblock` proves the converged family applies and replays with the strict result projection.
- VAL-020: PASS - `TestReviewRecordCreatesImmutableReview` and `TestP23FindingDispositionRealProcessJourney` create Findings only as atomic effects of `review.recorded`, read them back unchanged, and exercise replay without a duplicate Event. There is no independent Finding create or update leaf; the contract inventory test admits only `finding address` and `finding list`.
- VAL-021: TASK 5 EVIDENCE PASS; FINAL ASSERTION DEPENDS ON TASK 6 VAL-026 - `TestFindingAddressDoesNotResolveOrUnblock` rejects no-Claim and review-Claim requests, requires `claim-mode-invalid` with `required_mode:"mutation"`, then proves a mutating Claim produces `addressed` while `resolved` remains false and `blocking` remains true. The final non-vacuous proof that the review Obligation stays unsatisfied is intentionally deferred to Task 6's VAL-026 sequence because Obligation evaluation is outside this task.
- VAL-022: PASS - `TestTerminalFindingDispositionAuthorityMatrix` runs both terminal actions across creating Review, mutating holder, blocked verdict, different independent Actor, same independent Reviewer later, and already-terminal cases. `TestReviewRecordMixedTerminalDispositions` and `TestP23FindingDispositionRealProcessJourney` prove atomic mixed terminal dispositions, exact inapplicability reason, qualifying Review and Actor, resolved-only Review output, typed negatives, replay, and changed-reason idempotency conflict.
- VAL-023: PASS - `TestFindingAddressDoesNotResolveOrUnblock`, `TestReviewRecordMixedTerminalDispositions`, and `TestP23FindingDispositionRealProcessJourney` prove supplied state rejection and derived `open` to `addressed`, `resolved`, and `inapplicable` transitions through live reads and rebuild.

## Tests run

- Required first RED: `go test -race -timeout 30m ./cmd/fid/ -run '^TestFindingAddressDoesNotResolveOrUnblock$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); the active mutation Claim could not yet address a Finding.
- Required first GREEN: the byte-identical command above - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 13.475s` (terminated `exit_code: 0`).
- Mixed-terminal fixture attempt: `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordMixedTerminalDispositions$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); the new fixture supplied unsorted Finding IDs.
- Mixed-terminal product RED: the byte-identical command above - FAIL - `FAIL` (terminated `exit_code: 1`); the optional inapplicability collection was not yet admitted.
- Mixed-terminal GREEN: the byte-identical command above - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 15.807s` (terminated `exit_code: 0`).
- Focused acceptance matrix: `go test -race -timeout 30m ./cmd/fid/ ./internal/review/ ./internal/cli/ ./internal/journal/ -run '^(TestFindingAddressDoesNotResolveOrUnblock|TestReviewRecordMixedTerminalDispositions|TestTerminalFindingDispositionAuthorityMatrix|TestValidateTerminalFindingInputBoundaries|TestDecodeClosedReviewInputRejectsMalformedInapplicability|TestReviewRecordCreatesImmutableReview|TestFindingAddressedFenceConvergencePreservesEveryRemainingEvent|TestProjectEventFenceRejectsEveryUnavailablePhase23EventWithoutProjectionChange)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/journal 2.211s` (terminated `exit_code: 0`); preceding package tokens were `cmd/fid 37.690s`, `internal/review 1.245s`, and `internal/cli 1.102s`.
- Focused fence and executable-leaf inventory: `go test -race -timeout 30m ./cmd/fid/ ./internal/journal/ -run '^(TestLiveLeafCommandsMatchContractInventories|TestFindingAddressedFenceConvergencePreservesEveryRemainingEvent|TestProjectEventFenceRejectsEveryUnavailablePhase23EventWithoutProjectionChange)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/journal 2.237s` (terminated `exit_code: 0`); the preceding package token was `cmd/fid 3.406s`.
- Integration helper compile attempt, attributed to `/root/task4_process_journey`: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/... -run '^TestP23FindingDispositionRealProcessJourney$'` - UNVERIFIED - reported `FAIL` with `exit_code: 1`; the test-only fixture contained an unused variable.
- Integration helper first running attempt, attributed to `/root/task4_process_journey`: the byte-identical command above - UNVERIFIED - reported `FAIL` with `exit_code: 1`; five fixture Finding IDs violated retained lexical ordering.
- Integration helper focused evidence, attributed to `/root/task4_process_journey`: the byte-identical command above - UNVERIFIED - reported `ok  github.com/mbell-nzl/fid/cmd/fid/integration 7.858s` with `exit_code: 0`; the committed helper change is `1942d8cd029e286d7a40ea4ea330a2315b80c86b`.
- Independently captured real-binary journey: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/... -run '^TestP23FindingDispositionRealProcessJourney$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 8.289s` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 237.239s coverage: 64.1% of statements` (terminated `exit_code: 0`)
- Ordinary package budgets: `cmd/fid` was 588.268s and `internal/workflow` was 237.239s; every package stayed below 600s.
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 64.6%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 669.496s` (terminated `exit_code: 0`)
- Tagged integration budget: `cmd/fid/integration` was 669.496s, below 900s.

## Deviations from spec

- The first mixed-terminal fixture attempt used Finding IDs in non-lexical order and failed before reaching the intended product tracer. Sorting the fixture preserved the accepted input contract; the next byte-identical attempt exposed the absent DEC-242 input, and the later attempt passed.
- The integration helper's initial compile attempt retained an unused variable. After that test-only repair, its first running attempt exposed the same fixture ordering requirement across five Findings; reordering those fixture IDs preserved product behavior. Two assertion repairs then matched the accepted contract: an addressed Finding has JSON `null` for `qualifying_review_id`, and the fixture releases the same Reviewer's Claim before testing independent later-Review authority.
- No product, schema, Event, protected gate, threshold, or validation-contract deviation was made.

## Concerns, technical debt, or notes for the Tech Lead

- VAL-021's Task 5 behavior is complete, but its final non-vacuous Obligation proof remains the Task 6 VAL-026 dependency recorded in the validation contract.
- `cmd/fid` completed in 588.268s, leaving 11.732s below the 600s ordinary-package budget. The gate passed, but Task 6 should avoid adding repeated full CLI fixtures where a shared snapshot can preserve isolation.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `0be596e9068795ceb3a0a89bdac50d538f8e8791`

## Commit subject

- `feat(finding): implement disposition lifecycle`
