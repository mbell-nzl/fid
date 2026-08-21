# Phase 20 - Task 9a: Mutation concurrency, transition, and attribution evidence - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: Fresh final QA re-review.
- Human attention required: no

## Summary of what was implemented

- Closed Review 3 Blocking finding 1 with one canonical JSON matrix covering absent, matching, and stale `--expected-revision` cases for both Issue update and transition.
- Closed Review 3 Blocking finding 2 with the complete canonical transition outcome table, journal-tail Event identity comparison, and contract-correct handling of unknown target states and stored unsupported rules.
- Closed Review 3 Blocking finding 3 with four owner-flag rejection cells, canonical update owner-invariance evidence, retained default-owner regression coverage, and create/update/transition attribution checks that exclude protected effects and projection writes.

## Files created or modified

- `cmd/fid/issue_transition_test.go`
- `cmd/fid/issue_update_test.go`
- `internal/issues/transition.go`
- `internal/issues/transition_rules.go`
- `internal/issues/transition_test.go`
- `project-docs/docs/phases/phase-20/tasks/task-9a-done.md`

## Documentation updated

- Added this completion report.

## Validation Contract Coverage

- VAL-002: PASS - `TestIssueTransitionExpectedRevisionMatrix` at `cmd/fid/issue_transition_test.go:151` and `TestIssueCreateAndUpdateRejectRemovedOwnerFlagsWithCanonicalFailures` at `cmd/fid/issue_update_test.go:311` assert non-zero canonical rejections with empty stdout; `TestIssueTransitionCanonicalOutcomeTable` at `cmd/fid/issue_transition_test.go:338` does the same for every transition rejection.
- VAL-008: PASS - the success branches at `cmd/fid/issue_transition_test.go:128` and `:307`, owner/update success at `cmd/fid/issue_update_test.go:380`, and attribution success at `cmd/fid/issue_transition_test.go:422` each require stdout-only canonical responses.
- VAL-010: PASS - `TestIssueTransitionExpectedRevisionMatrix` at `cmd/fid/issue_transition_test.go:142` asserts both update and transition successes return prior revision plus one; `TestIssueTransitionCanonicalOutcomeTable` at `:325` covers the legal transition row. The real `/tmp/fid` matching transition returned revision `2` from revision `1`.
- VAL-011: PASS - `TestIssueTransitionCanonicalOutcomeTable` at `cmd/fid/issue_transition_test.go:331` compares the successful transition response with the appended journal tail; retained reserved-rule coverage at `internal/issues/transition_test.go:406` asserts unsupported stored rules reject without appending.
- VAL-017: PASS - `TestIssueTransitionCanonicalOutcomeTable` at `cmd/fid/issue_transition_test.go:331` compares the successful transition envelope Event ID and timestamp to the journal tail.
- VAL-018: PASS - Review 3 Blocking finding 1 closed by `TestIssueTransitionExpectedRevisionMatrix` at `cmd/fid/issue_transition_test.go:65`: the declared six-cell update/transition table at `:71`, success append/revision assertions at `:142`, and stale `state-conflict`, unchanged projection, and zero-append assertions at `:170`, `:176`, and `:179` prove every required state.
- VAL-019: PASS - Review 3 Blocking finding 2 closed by `TestIssueTransitionCanonicalOutcomeTable` at `cmd/fid/issue_transition_test.go:186`: it names legal, unknown target, absent edge, stored unsupported precondition, and unsatisfied-precondition rows at `:197`, `:207`, `:214`, `:238`, and `:264`; rejection details, unchanged status, and zero append are asserted at `:358`, `:361`, `:364`, and `:367`.
- VAL-040: PASS - Review 3 Blocking finding 3 closed by `TestIssueCreateAndUpdateRejectRemovedOwnerFlagsWithCanonicalFailures` at `cmd/fid/issue_update_test.go:266`, which covers both owner flags for both create and update with canonical failure details and zero append. `TestIssueUpdateCanonicalChangesExcludeOwnerAndPreserveProjection` at `:350` asserts all mutable change fields omit owner and projected owner values remain unchanged at `:393` and `:403`. Retained default-owner regression: `TestIssueCreateSupportsHierarchyAndAuthoringFlags` in `cmd/fid/issue_test.go` passed in the focused run and continues to assert task/subtask default `owner_role` values.
- VAL-041: PASS - Review 3 Blocking finding 3 closed by `TestIssueMutationCanonicalAttributionAvoidsProtectedEffects` at `cmd/fid/issue_transition_test.go:374`: create/update/transition rows are at `:380`, `:390`, and `:399`; journal actor/Event identity and allowed-only effects are asserted at `:442` and `:446`, projected activity attribution at `:451`, and every other projection-table count is invariant at `:454`.
- VAL-044: PASS - transition closed-request rows are exercised by the six-cell matrix at `cmd/fid/issue_transition_test.go:65`; the four removed owner request cells, exact canonical failure shape, and zero append are exercised at `cmd/fid/issue_update_test.go:266` through `:344`.

## Tests run

- First run: `go test ./cmd/fid -run '^TestIssueTransitionExpectedRevisionMatrix$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.147s` (terminated `exit_code: 0`).
- First run: `go test ./cmd/fid -run '^TestIssueTransitionCanonicalOutcomeTable$' -count=1` - expected FAIL - `FAIL` (terminated `exit_code: 1`). The observed failure named `Code:invalid-argument` for both unknown-target and stored-unsupported-precondition cells.
- GREEN: `go test ./cmd/fid -run '^TestIssueTransitionCanonicalOutcomeTable$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.148s` (terminated `exit_code: 0`).
- First run: `go test ./cmd/fid -run '^(TestIssueCreateAndUpdateRejectRemovedOwnerFlagsWithCanonicalFailures|TestIssueUpdateCanonicalChangesExcludeOwnerAndPreserveProjection)$' -count=1` - expected FAIL - `FAIL` (terminated `exit_code: 1`). The observed owner-test decoder showed `Version:` empty because it read `version` instead of the contract's `contract_version`; the all-fields update test itself printed no failure.
- GREEN: `go test ./cmd/fid -run '^(TestIssueCreateAndUpdateRejectRemovedOwnerFlagsWithCanonicalFailures|TestIssueUpdateCanonicalChangesExcludeOwnerAndPreserveProjection)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.076s` (terminated `exit_code: 0`).
- First run: `go test ./cmd/fid -run '^TestIssueMutationCanonicalAttributionAvoidsProtectedEffects$' -count=1` - expected FAIL - `FAIL` (terminated `exit_code: 1`). The observed transition fixture failure was `projection table "claims" count after transition = 1, want unchanged 0`.
- GREEN: `go test ./cmd/fid -run '^TestIssueMutationCanonicalAttributionAvoidsProtectedEffects$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.210s` (terminated `exit_code: 0`).
- RED: `go test ./cmd/fid -run '^TestIssueTransitionCanonicalOutcomeTable$' -count=1` after correcting the public Event decoder - expected FAIL - `FAIL` (terminated `exit_code: 1`). The last two failing rows were unknown target and stored unsupported precondition, each reporting `invalid-argument`.
- GREEN: `go test ./internal/issues -run '^TestTransitionTreatsReservedGateRuleAsUnsatisfied$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/issues	0.056s` (terminated `exit_code: 0`).
- `go test ./cmd/fid ./internal/issues -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/issues	4.672s` (terminated `exit_code: 0`)
- Real entrypoint in `/tmp/fid-task-9a.xP3V7M`: `/tmp/fid init --prefix FID`, `/tmp/fid --json issue create --kind task --title 'Task 9a real entrypoint' --by developer`, `/tmp/fid issue claim FID-9cd7a5 --by developer`, and `/tmp/fid --json issue transition FID-9cd7a5 --to in_progress --expected-revision 1 --by developer` - PASS - returned `{"id":"FID-9cd7a5","revision":2}` with Event `EVT-20260814-000004` (each terminated `exit_code: 0`).
- Real entrypoint stale rejection: `/tmp/fid --json issue transition FID-9cd7a5 --to dev_complete --expected-revision 1 --by developer` - expected rejection - `{"actual_revision":2,"expected_revision":1}` (terminated `exit_code: 1`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					77.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.100s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

Fresh direct-recovery evidence captured on 2026-08-21:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					78.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.915s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- `TestIssueTransitionCanonicalOutcomeTable` first failed with verbatim token `FAIL` (terminated `exit_code: 1`), including `failure = {Operation:issue.transition Error:{Code:invalid-argument Details:{FromState: ToState: Unsatisfied:[]}}}, want issue.transition transition-not-permitted`. Root causes: authority lookup ran before declared-state validation, and stored unknown workflow rules returned an unclassified error. Fixed in `internal/issues/transition.go` by validating the target state before authorization and in `internal/issues/transition_rules.go` by treating an unknown stored rule as an unmet precondition; the byte-identical focused command then passed.
- The same first outcome-table run exposed a test decoder mismatch: it decoded public `{id, occurred_at}` response data as the internal journal Event schema. The test now decodes the public envelope Event shape; the assertion itself was not narrowed.
- The owner test's first run used `json:"version"` rather than the contract field `json:"contract_version"`; corrected the test decoder without changing its envelope assertion. The attribution test's first run measured before a journal-seeded claim had been projected; it now runs the supported rebuild command before snapshotting counts. Neither correction changed a product assertion, threshold, or accepted cell.
- The broader focused domain run first failed because `TestTransitionReservedGateRuleRemainsUnevaluated` expected the pre-contract untyped unsupported-rule error. It is strengthened as `TestTransitionTreatsReservedGateRuleAsUnsatisfied` to require `BlockedTransitionError` with the exact stored rule and retain the existing no-append assertion; this aligns its verdict with VAL-019.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `1f196a94b064554fbc8bae1b0d9238125e2ab1d4`

## Commit subject

- `fix(issue): recover mutation contract evidence`
