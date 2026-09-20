# Phase 23 - Task 4: Immutable Review Records and computed independence - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 9
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Implemented atomic `review.record` and immutable `review.get`, including repository re-observation, coverage, verdict and Finding validation, Receipt-to-Run binding, one Review per Run, strict typed Event/result projection, idempotent replay, and digest-conflict rejection.
- Enforced Reviewer independence from durable mutating Claim and Run history. Under DEC-239, the unreachable Run-holder-only relation is proved only through the isolated domain query; the real CLI still proves the Claim-holder rejection and same independent Reviewer's later Review after Builder rework.
- Corrected Finding resolution to prohibit only a Review resolving a Finding it created. A later Review by the same independent Actor can now resolve two Findings from that Actor's earlier Review in one Event.
- Added built-binary real-process coverage for typed independence, repository mismatch, replay, digest conflict, duplicate Review, same-Reviewer rework, two-Finding resolution, and immutable readback.
- Converged VAL-003 by retaining every still-unavailable Phase 23 fence key while `review.recorded` follows the strict live/replay path. Added direct cardinality evidence that only `review.recorded` admits repeated Finding-disposition inserts.
- Added all three real CLI `review-verdict-invalid` negatives with no Event or projection mutation and retained the six-case domain validity matrix.

## Files created or modified

- `cmd/fid/review.go`, `cmd/fid/review_test.go`, and `cmd/fid/root.go`
- `cmd/fid/integration/p23_review_record_test.go` and `cmd/fid/integration/p22_fake_test.go`
- `internal/review/review.go` and `internal/review/review_test.go`
- `internal/cli/codes.go`
- `internal/journal/current_event_cardinality.go`, `internal/journal/current_event_json.go`, `internal/journal/current_event_result_projection.go`, `internal/journal/current_event_result_values.go`, `internal/journal/current_event_validation_test.go`, `internal/journal/current_event_value_catalog.go`, `internal/journal/current_event_value_relations.go`, `internal/journal/project_effects.go`, and `internal/journal/project_effects_test.go`
- `internal/workflow/activation.go`
- `project-docs/docs/phases/phase-23/tasks/task-4-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-23/tasks/task-4-done.md`

## Validation Contract Coverage

- VAL-014: PASS - `TestP23ReviewRecordRealProcessJourney` records and reads immutable Run, Issue, Envelope revision, repository, verdict, and two-Finding bindings through the built binary; `TestReviewRecordIndependentReworkAndReceiptBinding` rejects a Receipt from another Run without persistence change; `TestValidateCoverageRejectsSourceInItemAndExclusion` rejects the shared-Source case.
- VAL-015: PASS - `TestP23ReviewRecordRealProcessJourney/mutating_claim_actor_is_rejected_without_an_event` names the exact Claim and Actor and proves no Event, then the same `actor-denied` Reviewer succeeds before and after Builder rework. `TestCheckIndependenceRejectsClaimAndRunHistoryDirections` proves the exact Run and Actor over isolated stored history. Per DEC-239, that fixture proves only the domain query and makes no CLI, journal-replay, or import reachability claim.
- VAL-003: PASS - the fence originated in `4aa23200789e0ce65172d6e42c49704849854a4a` and Task 4 removed only `review.recorded` in `b500ab45406dd1cbd79153232d2add3ca8a74021`. `TestReviewRecordedFenceConvergencePreservesEveryRemainingEvent` fixes the exact eight remaining keys, `TestProjectEventFenceRejectsEveryUnavailablePhase23EventWithoutProjectionChange` proves their no-effect rejection, and `TestReviewRecordCreatesImmutableReview` proves byte-identical `review.get` output after rebuild. The focused executing command is recorded under Tests run.
- VAL-016: PASS - `TestP23ReviewRecordRealProcessJourney/repository_mismatch_is_typed_and_non_mutating` exercises the built binary, requires typed expected/observed identities, and proves no Event; `TestReviewRecordCreatesImmutableReview` also compares journal and projection snapshots.
- VAL-017: PASS - `TestReviewRecordCreatesImmutableReview` sends all three invalid combinations through the CLI, requires `review-verdict-invalid`, and proves no persistence change. `TestValidateVerdictFindingAgreement` covers all three valid counterparts and requires typed `*VerdictInvalidError` values for the negatives. The real-process journey records valid `changes-required` and `pass` Reviews.
- VAL-018: PASS - `TestP23ReviewRecordRealProcessJourney/replay_duplicate_and_digest_conflict_preserve_history` rejects a second Review for the Run as `duplicate-entity` and proves no Event.
- VAL-019: PASS - that same built-binary subtest proves `replayed:true` without an Event and `idempotency-conflict` for a changed digest without an Event.
- Acceptance evidence toward VAL-020: PASS - `TestReviewRecordIndependentReworkAndReceiptBinding` and `TestP23ReviewRecordRealProcessJourney` each resolve two prior Findings in one `review.recorded` Event and read back both dispositions. `TestReviewRecordedAloneAdmitsRepeatedFindingDispositionEffects` proves `finding.addressed` retains its maximum of one disposition insert. `TestCurrentReviewRecordedRegistersItsDistinctJudgmentValues` also proves the retained `Subject` shape and separate `JudgmentSubject` registration.

## Tests run

- Historical first tracer: `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordRejectsNonIndependentActor$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); `review.record` was absent.
- Historical strict-result tracer: `go test -race -timeout 30m ./internal/journal/ -run '^TestCurrentReviewRecordedRegistersItsDistinctJudgmentValues$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); Review value shapes were absent.
- Historical verdict tracer: `go test -race -timeout 30m ./internal/review/ -run '^TestValidateVerdictFindingAgreement$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); Review/Finding validation was absent.
- Historical coverage tracer: `go test -race -timeout 30m ./internal/review/ -run '^TestValidateCoverageRejectsSourceInItemAndExclusion$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); closed coverage validation was absent.
- Historical Finding projection tracer: `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordCreatesImmutableReview$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); non-empty Finding results retained an Event timestamp placeholder.
- Historical typed adapter tracer: `go test -race -timeout 30m ./cmd/fid/ ./internal/review/ -run '^(TestReviewRecordCreatesImmutableReview|TestValidateVerdictFindingAgreement)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); repository mismatch passed Go structs to the strict encoder.
- Historical rework fixture attempts: `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordIndependentReworkAndReceiptBinding$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); first the Actor profile was unsorted, then the Actor-based Finding guard rejected the same Reviewer's later Review.
- Historical required gate: `go test -race -timeout 30m -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`); retained executable-leaf and Attention inventories had not yet converged.
- Same-Reviewer regression RED: `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordIndependentReworkAndReceiptBinding$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); restoring the Actor-based guard rejected the later resolving Review.
- Same-Reviewer regression GREEN: `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordIndependentReworkAndReceiptBinding$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 25.318s` (terminated `exit_code: 0`).
- Focused acceptance audit: `go test -race -timeout 30m ./internal/journal/ ./internal/review/ ./cmd/fid/ -run '^(TestProjectEventFenceRejectsEveryUnavailablePhase23EventWithoutProjectionChange|TestReviewRecordedFenceConvergencePreservesEveryRemainingEvent|TestReviewRecordedAloneAdmitsRepeatedFindingDispositionEffects|TestValidateVerdictFindingAgreement|TestCheckIndependenceRejectsClaimAndRunHistoryDirections|TestReviewRecordCreatesImmutableReview|TestReviewRecordIndependentReworkAndReceiptBinding)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 36.117s` (terminated `exit_code: 0`); preceding package tokens were `internal/journal 2.243s` and `internal/review 1.028s`.
- Integration helper first attempt: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/... -run '^TestP23ReviewRecordRealProcessJourney$'` - FAIL - `FAIL` (terminated `exit_code: 1`); the new Envelope grant fixture was unsorted.
- Integration helper focused evidence, attributed to `/root/task4_process_journey`: the byte-identical command above - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 4.509s` (terminated `exit_code: 0`).
- Integration helper shared-fixture evidence, attributed to `/root/task4_process_journey`: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/... -run '^TestP22FinalActivatedActionSetAcceptsProfilesAndGrants$'` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 3.078s` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 217.082s coverage: 64.1% of statements` (terminated `exit_code: 0`); `cmd/fid` was 496.639s and every ordinary package stayed below 600s.
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 64.9%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 659.237s` (terminated `exit_code: 0`); below the 900s budget.

## Deviations from spec

- The initial implementation attempt used the nonexistent `config.DatabaseName` and failed to compile; the repair used the existing `config.DatabaseFileName`. No contract or gate changed.
- The first typed repository probe failed because the adapter passed `review.Repository` structs to the strict error encoder; the repair emits contract-shaped maps. No contract or gate changed.
- The earlier report described same-Actor resolution as an "intentional self-resolution rejection." That claim was wrong: DEC-225 prohibits self-resolution by Review identity, while permitting a later Review by the same independent Actor. The regression reproduced the failure, `resolveFindingEffects` now compares the creating and resolving Review IDs, and both in-process and built-binary journeys pass.
- The first real-process journey attempt failed because its second Envelope grant was not sorted by Actor. Sorting the fixture preserved the accepted request and the byte-identical focused command then passed. No product behavior or protected definition changed.
- DEC-239 amended only VAL-015's evidence route. The isolated Run-holder history fixture is not claimed to be reachable through `run.start`, valid canonical journal replay, or import.

## Concerns, technical debt, or notes for the Tech Lead

- The Run-holder-only VAL-015 branch is intentionally a domain-query fixture under DEC-239 because the accepted Run/Claim Actor binding cannot construct that relation through the CLI.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `81651df858143933c6e9c4501d0147700f9b6704`

## Commit subject

- `fix(review): admit same reviewer rework`
