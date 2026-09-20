# Phase 23 - Task 4: Immutable Review Records and computed independence - Completion Report

## Workflow outcome

- Task outcome: BLOCKED
- Category: protected-contract-conflict
- Condition: The real CLI journey now proves a different-Run Receipt rejection with no Event, independent two-Reviewer rework, and prior-Finding resolution. VAL-015's required real CLI Run-holder-only direction is impossible under the accepted Run binding: `validateRunStartBinding` requires `claim.ActorID == options.ActorID`, while durable Claim history is retained and `checkIndependence` checks that Claim before it can name the Run.
- Recovery attempts: 8
- Next action: Tech Lead must decide whether the accepted Run binding permits a distinct Run Actor, or whether VAL-015 accepts the retained direct domain-history proof for the unrepresentable Run-holder-only state. Either choice changes protected public/data-model semantics or validation acceptance.
- Human attention required: yes

## Summary of what was implemented

- Added the mandated vertical red tracer for `review.record` and retained its original absent-command failure evidence.
- Added `fid review record` and `internal/review`'s write-locked independence admission. It derives conflicts solely from durable mutating Claim and Run history for the reviewed Run's Issue and Envelope revision, and returns the exact conflicting subject and Actor.
- Removed only the `review.recorded` availability fence and registered Task 4's closed `review.record` capability in both authoritative capability allowlists.
- Completed the first atomic `review.record` slice: it locks, refreshes the projection, validates the active review Run/Claim and authority grant, re-observes repository identity, appends a `review.recorded` Event with idempotency, projects `review_records`, and returns the canonical typed result.
- Added `fid review get <review-id>` and a real CLI record/get tracer. Converged the retained executable-leaf and Attention surface inventories for the accepted new Review leaves.
- Added closed-coverage validation that rejects an item and exclusion sharing one Source digest before Event creation. Expanded the real CLI tracer to record, strictly project, and read back two non-blocking Findings in one Event.
- Added bounded `resolves_finding_ids` admission and Event effects: each requested prior Finding must be unique, belong to the same Issue and Envelope revision, have a different recording Actor, and be non-terminal; accepted IDs add atomic `resolved` disposition inserts linked to the qualifying Review.
- Completed the typed `repository-mismatch` adapter: expected and observed repository identities now serialize as the public-interface object rather than leaking Go pointer values into the failure envelope.
- Extended the real CLI Review tracer to retain no-Event snapshots for repository mismatch, duplicate Review, and idempotency digest-conflict denials; it also proves exact replay semantics.
- Updated Review read projection to expose `resolved_finding_ids`, each Finding's immutable disposition facts, and its derived disposition state. Resolution admission now rejects a resolving `blocked` Review.
- Added a fixture with two separately activated Reviewer Actors. Its CLI journey records a Receipt on a completed mutating Run, rejects that Receipt on the next review Run with unchanged journal/projection state, then completes rework and records a different Reviewer's resolving `pass` Review.
- Added the complete durable-history query matrix: Claim-holder and Run-holder conflicts name their exact mutating subject. The Run-holder row is exercised against the admissible stored-history relation because the accepted CLI lifecycle cannot create that relation.

## Files created or modified

- `cmd/fid/review_test.go`
- `cmd/fid/review.go`
- `internal/review/review.go`
- `internal/cli/codes.go`
- `cmd/fid/root.go`
- `internal/journal/{current_event_cardinality.go,current_event_json.go,current_event_result_projection.go,current_event_result_values.go,current_event_validation_test.go,current_event_value_catalog.go,current_event_value_relations.go,project_effects.go}`
- `internal/workflow/activation.go`
- `project-docs/docs/phases/phase-23/tasks/task-4-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-23/tasks/task-4-done.md`

## Validation Contract Coverage

- VAL-014: PASS - `cmd/fid/review_test.go:126` records a Receipt on the completed mutation Run; `:148` binds it to a later review Run; `:151` snapshots persistence and `:152` rejects it, with `:156` proving no Event or projection mutation. The pre-existing same-digest and immutable readback tests remain preserved.
- VAL-015: BLOCKED - `cmd/fid/review_test.go:146` through `:187` proves two independent Actors and a real rework/resolution cycle. `internal/review/review_test.go:45` proves both query directions and exact Claim/Run identifiers, but the requested real CLI Run-only Actor is incompatible with the accepted binding described under Workflow outcome.
- VAL-003: UNVERIFIED - The `review.recorded` fence is removed, while the retained Task 1 fence-parity assertion has not yet been converged to a still-fenced Event or run in this recovery.
- VAL-016: PASS - `cmd/fid/review_test.go:76` supplies a stale repository, asserts typed `repository-mismatch` with expected/observed values, and compares the database and journal snapshots before and after the rejection; the following current-identity record succeeds at `:87`.
- VAL-017: PASS - `internal/review/review_test.go:10` exercises all six valid and invalid verdict/Finding combinations; the two-Reviewer CLI journey exercises `changes-required` at `cmd/fid/review_test.go:160` and resolving `pass` at `:178`.
- VAL-018: PASS - `cmd/fid/review_test.go:103` records a second Review under a new key, asserts `duplicate-entity`, and proves the persistence snapshot is unchanged.
- VAL-019: PASS - `cmd/fid/review_test.go:98` proves same-key replay returns `replayed:true`; `:115` changes the profile under that key, asserts `idempotency-conflict`, and proves no Event or projection change.

## Tests run

- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordRejectsNonIndependentActor$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); expected red state: `review.record` is absent, so no typed `review-not-independent` rejection can be emitted.
- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordRejectsNonIndependentActor$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 8.269s` (terminated `exit_code: 0`)
- The focused green response contains `review-not-independent` and names the mutating Claim and `actor-builder`.
- `go test -race -timeout 30m ./internal/journal/ -run '^TestCurrentReviewRecordedRegistersItsDistinctJudgmentValues$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); expected red state: `Review` was absent from `currentValueShapes` after the Event fence was removed.
- `go test -race -timeout 30m ./internal/journal/ -run '^TestCurrentReviewRecordedRegistersItsDistinctJudgmentValues$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/journal 1.016s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m ./internal/journal/ -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/journal 159.995s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m ./internal/review/ -run '^TestValidateVerdictFindingAgreement$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); expected red state: Task 4 Review/Finding types and verdict validation were absent.
- `go test -race -timeout 30m ./internal/review/ -run '^TestValidateVerdictFindingAgreement$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/review 1.017s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordCreatesImmutableReview$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 9.350s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m ./internal/review/ -run '^TestValidateCoverageRejectsSourceInItemAndExclusion$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); expected red state: the closed coverage validator was absent.
- `go test -race -timeout 30m ./internal/review/ -run '^TestValidateCoverageRejectsSourceInItemAndExclusion$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/review 1.017s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordCreatesImmutableReview$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); expected red state: non-empty Finding result values carried the Event timestamp placeholder.
- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordCreatesImmutableReview$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 9.070s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m ./internal/review/ ./cmd/fid/ -run '^(TestValidateVerdictFindingAgreement|TestReviewRecordCreatesImmutableReview)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 9.078s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m ./cmd/fid/ ./internal/review/ -run '^(TestReviewRecordCreatesImmutableReview|TestValidateVerdictFindingAgreement)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); red adapter evidence: `repository-mismatch` failure construction rejected the `expected` repository struct before it could write the typed envelope.
- `go test -race -timeout 30m ./cmd/fid/ ./internal/review/ -run '^(TestReviewRecordCreatesImmutableReview|TestValidateVerdictFindingAgreement)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/review 1.018s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordCreatesImmutableReview$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 10.241s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordIndependentReworkAndReceiptBinding$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); red fixture evidence: the initial unsorted activated Actor profile was rejected before the journey began.
- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordIndependentReworkAndReceiptBinding$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); red journey evidence: self-resolution by the same Reviewer correctly rejected the prior Finding.
- `go test -race -timeout 30m ./cmd/fid/ ./internal/review/ -run '^(TestReviewRecordCreatesImmutableReview|TestReviewRecordIndependentReworkAndReceiptBinding|TestCheckIndependenceRejectsClaimAndRunHistoryDirections)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/review 1.025s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordIndependentReworkAndReceiptBinding$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 25.479s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 64.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.0%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 665.125s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 64.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 64.9%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 663.124s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`); retained executable-leaf and Attention inventories needed Task 4 convergence.
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 64.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.0%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 694.141s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 64.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 64.8%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 662.335s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 64.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 64.9%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 693.206s` (terminated `exit_code: 0`)
- Observed wall times: `cmd/fid` 542.160s failed attempt, then 473.760s, 457.990s, 442.371s, 446.494s, and 471.636s passing reruns; `internal/workflow` 230.960s first attempt, cached on later reruns including attempt 8; `cmd/fid/integration` 694.141s, 665.125s, 663.124s, 662.335s, and 693.206s. All applicable 600s/900s budgets passed.

## Deviations from spec

- Focused implementation attempt initially failed to compile with `undefined: config.DatabaseName` (terminated `exit_code: 1`). The in-scope repair used the existing `config.DatabaseFileName`; the later focused tracer passed. No task command, assertion, threshold, or policy changed.
- The current recovery's first typed repository probe failed with `FAIL` (terminated `exit_code: 1`) because `repository-mismatch` passed `review.Repository` values directly to the strict CLI error encoder. The repair serializes contract-shaped expected and observed maps; the later focused test and complete required gate pass. No command, assertion, threshold, or policy changed.
- The first new fixture failed with `FAIL` (terminated `exit_code: 1`) because authority-profile Actor IDs must be sorted; reordering the fixture preserved its identity and enabled the journey. The next run exposed intentional self-resolution rejection; the final fixture uses two activated Reviewer Actors. No command, assertion, threshold, or policy changed.

## Concerns, technical debt, or notes for the Tech Lead

- The only unresolved assertion is protected: an Actor that “held only a Run under a since-released mutating Claim” cannot be produced by the accepted CLI because `run.start` requires that actor to own the Claim, and history preserves that Claim. The direct query test confirms the intended rejection for the requested stored relation, but cannot replace the required real CLI observation without an authority decision.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `a1ae1bb61363a591bb46c9ddfd4b8df36462687b`

## Commit subject

- `test(review): cover independent reviewer rework`
