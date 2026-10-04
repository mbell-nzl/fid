# Phase 24 - Task 5: Review leaves and repository expectations - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

OUTCOME COMPLETE

## Summary of what was implemented

- Added the shared `ExpectedRepository` pair and check. `review record` now requires that pair, compares only HEAD and dirty against one fid observation, stores the observed Repository, and returns the amended closed `repository-mismatch` details without an Event on a difference.
- Converted Review coverage items, exclusions and evidence references to `SourceInput`. Resolution checks supplied project-file digests or observes omitted ones before coverage and digest-exclusivity rules. Replays return before another repository or Source observation.
- Adopted descriptors for `review record` and `finding address`, including help, examples, exclusivity and structural errors. Converged displaced retained tests and fake/fid fixtures while preserving their assertions.

## Files created or modified

- `internal/project/expected_repository.go`, `internal/project/expected_repository_test.go`
- `internal/review/review.go`, `internal/review/review_test.go`
- `internal/cli/codes.go`, `internal/cli/codes_test.go`
- `cmd/fid/review.go`, `cmd/fid/finding.go`, `cmd/fid/input_descriptors_review.go`, `cmd/fid/input_descriptor_test.go`, `cmd/fid/review_test.go`
- `cmd/fid/integration/p23_conformance_test.go`, `cmd/fid/integration/p23_failure_contract_test.go`, `cmd/fid/integration/p23_review_leaf_oracle_test.go`, `cmd/fid/integration/p23_review_record_test.go`, `cmd/fid/integration/p24_judgment_corpus_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`, `project-docs/docs/phases/phase-24/tasks/task-5-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` documents the expected pair, both Git commands, SourceInput positions and mismatch shape.

## Validation Contract Coverage

- VAL-009: PASS for `review.record` - `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at `cmd/fid/integration/p23_review_record_test.go:121` advances HEAD, replays the original key without an Event, and checks a new key sees the changed HEAD.
- VAL-011: PASS for both leaves - `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` at `cmd/fid/input_descriptor_test.go:235` checks the live registry/type tree and seeded mismatch; `TestP24ReviewAndFindingDescriptorHelpAndStructure` at line 277 checks both leaves left the pending fence.
- VAL-012: PASS for both leaves - `TestP24ReviewAndFindingDescriptorHelpAndStructure` at `cmd/fid/input_descriptor_test.go:277` checks the help template line and both Git command meanings; `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` checks field conformance.
- VAL-013: PASS for both leaves - `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:44` runs each standalone `--example` through the built CLI; direct CLI invocations below show the ordered JSON templates.
- VAL-014: PASS for both leaves - `TestP24ReviewAndFindingDescriptorHelpAndStructure` at `cmd/fid/input_descriptor_test.go:277` validates and decodes each generated example into the closed Go request type.
- VAL-015: PASS for both leaves - `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:44` exercises six conflicting inputs per leaf, and standalone success, with empty non-project workspaces.
- VAL-016: PASS for both leaves - `TestP24ReviewAndFindingDescriptorHelpAndStructure` at `cmd/fid/input_descriptor_test.go:277` checks exact top-level and nested paths for unknown, duplicate, missing, null, type, format and enum failures through the CLI, with no created project state.
- VAL-017: PASS for `review.record` - `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at `cmd/fid/integration/p23_review_record_test.go:45` rejects top-level and all three nested caller `repository` positions with their exact paths and unchanged Event counts.
- VAL-018: PASS for `review.record` - the same test at `cmd/fid/integration/p23_review_record_test.go:77` checks mismatching supplied digests at item, exclusion and evidence positions with both digest values and no Event; its accepted Review at line 104 checks a matching supplied exclusion digest and all three stored resolved digests. All positions use the same `resolveSources` check.
- VAL-019: PASS for `review.record` - the same test at `cmd/fid/integration/p23_review_record_test.go:104` checks omitted item and evidence digests against independently hashed file bytes; at line 154 it checks missing-file `source-digest-mismatch` with both digest fields null and no Event.
- VAL-022: PASS for `review.record` - the same test at `cmd/fid/integration/p23_review_record_test.go:110` compares stored Review and all three Source repositories with the single fid observation and independently checks HEAD, dirty and tree; `TestP23ReviewRecordRealProcessJourney` at line 177 retains exact readback.
- VAL-023: PASS for `review.record` - `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at `cmd/fid/integration/p23_review_record_test.go:135` checks new-key `repository-unavailable` in a non-Git state versus normal success at line 104, with unchanged Event count.
- VAL-024: PASS for `review.record` - `TestP23ReviewRecordRealProcessJourney` at `cmd/fid/integration/p23_review_record_test.go:208` checks dirty and HEAD differences with no Events; its first matching Review succeeds at line 252.
- VAL-025: PASS - the two mismatch subtests of `TestP23ReviewRecordRealProcessJourney` at `cmd/fid/integration/p23_review_record_test.go:208` and line 231 check expected echo and full observed Repository; `TestReasonCodeDetailSchemaIsExactForEveryCode` at `internal/cli/codes_test.go:38` pins the closed shape.
- VAL-035: PASS for the two leaves - `TestP23SharedLeafConformance` at `cmd/fid/integration/p23_conformance_test.go:143` compares accepted, unknown-field and conflict fake/fid outcomes over the derived operation set, including both Review leaves.
- VAL-052: PASS - the final ordinary gate reported `cmd/fid` 267.739 s, `internal/workflow` 115.695 s and slowest other ordinary package `internal/issues` 158.407 s; the tagged gate reported 252.256 s. None reached the 480 s ordinary warning, 600 s ordinary limit or 900 s tagged limit.
- VAL-060: PASS for both leaves - `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:44` runs each standalone example from an empty non-project directory, checking exit 0, empty stderr and no created entries.
- VAL-066: PASS for `review.record` - `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at `cmd/fid/integration/p23_review_record_test.go:135` replays without an observable Git repository, and at line 149 replays after deletion of a coverage file. New-key controls return `repository-unavailable` and `source-digest-mismatch` with no Event.

## Tests run

- First TDD red: `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run 'TestP23ReviewRecordRealProcessJourney/dirty_expectation_mismatch_echoes_only_the_expected_pair' ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`); the old decoder rejected `expected_repository` before the repository check.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -timeout 30m ./internal/project ./internal/review ./internal/cli ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`) on the first retained-fixture attempt; after convergence, the identical command passed with `ok  	github.com/mbell-nzl/fid/cmd/fid	194.986s` (terminated `exit_code: 0`).
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run 'TestP24JudgmentKeyedLeafConflictAndCanonicalRetryCorpus|TestP23ReviewRecordRealProcessJourney' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.435s` (terminated `exit_code: 0`).
- `golangci-lint run ./...` - FAIL - `* staticcheck: 1` (terminated `exit_code: 1`) on the first exact-gate attempt; the corrected command passed as recorded below.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`) on its first exact-gate attempt; the corrected command passed as recorded below.
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `/tmp/fid review record --example` - PASS - `}` (terminated `exit_code: 0`).
- `/tmp/fid finding address --example` - PASS - `}` (terminated `exit_code: 0`).

Both direct CLI invocations ran from `/tmp`.

## Exact quality-gate results

Working directory: `/home/mark/projects/fid`. The five required commands were run in task-spec order. After causal repairs, only failed commands were repeated; completed passing checks were retained.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	115.695s	coverage: 62.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)				60.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	252.256s` (terminated `exit_code: 0`)

## Deviations from spec

- No accepted behavior, protected assertion, command or threshold changed. The first lint attempt ended `* staticcheck: 1` (terminated `exit_code: 1`) because a new test declared `var typed any =`; changing it to `typed :=` made the identical lint command return `0 issues.` (terminated `exit_code: 0`).
- The first exact tagged gate ended `FAIL` (terminated `exit_code: 1`). Task 5's new verdict descriptor rejected a retained corpus mutation's invalid enum, and two retained mismatch comparisons treated Go typed nil differently from JSON null. The corpus now changes `changes-required` to the accepted `blocked` verdict, preserving its idempotency-conflict assertion; the comparisons normalize the fixture to JSON while retaining full Repository equality. The identical tagged command later returned `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	252.256s` (terminated `exit_code: 0`).
- The touched-package race attempt ended `FAIL` (terminated `exit_code: 1`) because retained Review/Finding fixture inputs still used displaced shapes. Converging those inputs produced `ok  	github.com/mbell-nzl/fid/cmd/fid	194.986s` (terminated `exit_code: 0`) on the identical command.

## Concerns, technical debt, or notes for the Tech Lead

- Retained assertions changed only where Task 5 displaced their input or observation model. Introducing commits were `d215e40` and `926ef72` (phase-23/task-4 Review journeys), `d1c016d` (phase-23/task-13 failure corpus), `3335149` and `2e7fd2f` (phase-23/task-18 conformance and leaf oracle), `472ed64` (phase-23/task-13 code-schema assertion), and `9bbbdb5` (phase-24/task-2 keyed conflict corpus). Their subjects and no-Event checks remain intact.
- No package reached the 80% timing warning.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `ff6c923dcb706513e71ff7dca76723d68da6ba89`

## Commit subject

- `feat(review): check repository expectations and resolve sources`
