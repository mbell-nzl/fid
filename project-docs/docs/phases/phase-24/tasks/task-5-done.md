# Phase 24 - Task 5: Review leaves and repository expectations - Recovery Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: qa-rework
- Condition: F7's declared Review OID-format rejection fails VAL-016 and the ordinary required gate. The production descriptor/validator repair is outside this dispatch's exact three-file scope.
- Recovery attempts: 1
- Next action: Obtain bounded same-task production scope for cmd/fid/input_descriptors_review.go and internal/cli/input_validate.go, with their qualifying controls and tests; enforce the already accepted OID contract, preserve the red regression, and rerun focused red/green and the byte-identical five-command gate before completion.
- Human attention required: no

OUTCOME RECOVERY_REQUIRED

Dispatch binding: latest TASK-RECOVERY `2026-10-06 22:12`, task-5, trigger `review-3:qa-rework:task-5`, attempt `1`, category `qa-rework`, mode `task-scope`, authority `QA review`. Its ordered Allowed files and Forbidden changes were checked against the dispatch as exact copies. No ledger or triage edits were made by this Developer.

## Summary of what was implemented

- F4: captured both accepted R5/R6 inventory commands in this session, with actual match accounting and observed exit status.
- F5: added descriptor-checked three-position Review Source matrices, each covering supplied matching/mismatching and omitted existing/missing digests. Every row has a fresh valid Review Run; acceptance checks independently hashed stored digests, refusal checks exact digest/null details and no Event.
- F6: independently compared complete clean/dirty Review, three Source and Event repositories, including base, config project identity and root identity. Added complete original-envelope and journal-Event replay equality, normal/non-Git/unmerged/deleted-source pairs, and fresh unreviewed Run/new-key refusal controls whose identical restored requests succeed. Preserved the original Review journey and all previous assertions.
- F7: strengthened recursive help metadata and structural cases only for Review/Finding. Both help guards and Finding's structural corpus pass. Review's malformed expected head OID is now a retained failing assertion; completion requires an authorized same-task production repair.

## Files created or modified

- `cmd/fid/integration/p23_review_record_test.go`
- `cmd/fid/input_descriptor_test.go`
- `project-docs/docs/phases/phase-24/tasks/task-5-done.md`

## Documentation updated

- This report records the Window B repair evidence, remaining scope gap and original report verbatim as quoted history. No accepted product documentation changed.

## Validation Contract Coverage

- VAL-009: PASS for the repaired Review evidence - `TestP24ReviewReplayUnavailableAndRestoredControls` at `cmd/fid/integration/p23_review_record_test.go:207`, head-advanced row, compares the whole original response allowing only `result.replayed`, plus the actual original Event and unchanged Event count. Its fresh control refuses the advanced HEAD and succeeds after restoration. `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at line 298 additionally changes referenced file bytes before advancing HEAD and preserves full replay equality.
- VAL-011: PASS for retained Review/Finding adoption evidence - `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` at `cmd/fid/input_descriptor_test.go:236` checks live types and seeded mismatch; `TestP24ReviewAndFindingDescriptorHelpAndStructure` at line 634 checks both pending-fence removals before the later red format row.
- VAL-012: PASS - `TestP24ReviewHelpEvidenceRejectsIncompleteMetadata` at `cmd/fid/input_descriptor_test.go:596` checks both valid help outputs and independently corrupts every recursive field's type, requirement, meaning and applicable format/member/range/length metadata; every corruption is rejected. Review help retains both Git command meanings. The identical focused command passed after the red guard.
- VAL-013, VAL-015, VAL-060: PASS for retained two-leaf evidence - `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:45` exercises built-process standalone examples and conflicting inputs in empty non-project directories; it had no failure in the full ordinary attempt. Direct examples in Tests run additionally exercise the actual binary.
- VAL-014: PASS for retained typed-template evidence - both leaves' template generation and closed Go-request decode in `TestP24ReviewAndFindingDescriptorHelpAndStructure` at `cmd/fid/input_descriptor_test.go:634` execute before the format failure. VAL-014 explicitly permits placeholders that are not valid OIDs; the continuation must preserve that template meaning while using actual qualifying values for format controls.
- VAL-016: FAIL for Review, PASS for Finding - `p24ReviewRecursiveStructuralDetails` at `cmd/fid/input_descriptor_test.go:946` derives recursive missing/null/type/unknown/duplicate/enum/format/range/cardinality cases with exact path/value/allowed and closed error metadata, valid controls and empty-workspace/no-Event checks. Finding's entire focused selection passes. Review's `recursive/input.expected_repository.head_oid/format` fails at line 969: validation is nil for `not-an-oid`, expected path `input.expected_repository.head_oid`, value `not-an-oid`, allowed `["oid"]`. Other Review structural assertions remain present; the whole leaf is not claimed green.
- VAL-017: PASS for retained Review rejection evidence - `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at `cmd/fid/integration/p23_review_record_test.go:298` retains top-level and all three nested caller-repository refusals with exact paths and unchanged Events.
- VAL-018, VAL-019: PASS - `TestP24ReviewSourceMatrix` at `cmd/fid/integration/p23_review_record_test.go:104` derives/checks all three Source paths from real descriptor help and executes all twelve position/state rows. Matching and omitted-existing each compare the stored digest with independently hashed actual file bytes; mismatching checks supplied/actual digests and omitted-missing checks both nulls, each with exact closed details and no emitted/appended Event.
- VAL-022: PASS - `p24ReviewAssertRepositories` at `cmd/fid/integration/p23_review_record_test.go:167` compares full Review, all three Source and actual journal Event repositories against independent Git/config/root observations. Clean matrix/lifecycle acceptance uses `p24JudgmentExpectedRepository`; the original dirty journey uses independent known-file canonical-image hashing in `p24ReviewDirtyRepository` at line 183. All Repository fields, including base/config/root identity and dirty digest, participate.
- VAL-023, VAL-066: PASS - `TestP24ReviewReplayUnavailableAndRestoredControls` at `cmd/fid/integration/p23_review_record_test.go:207` has normal/non-Git/unmerged/deleted-source states. Each unavailable/deleted replay preserves the complete result and original Event with no Event appended. A fresh unreviewed Run and new key fail with exact unavailable reason/detail (and precondition/retryable metadata), or missing-Source null digest details, and no Event; the identical request/Run/key succeeds in the restored observable state, appending exactly one Event with independently checked repositories.
- VAL-024, VAL-025: PASS for retained Review evidence - `TestP23ReviewRecordRealProcessJourney` at `cmd/fid/integration/p23_review_record_test.go:460` retains matching/HEAD-different/dirty-different rows and full expected/observed mismatch details with no Event on refusal; `TestReasonCodeDetailSchemaIsExactForEveryCode` at `internal/cli/codes_test.go:38` retains the closed shape assertion.
- VAL-035: PASS for retained two-leaf evidence - `TestP23SharedLeafConformance` at `cmd/fid/integration/p23_conformance_test.go:143` compares accepted/unknown/conflict fake/fid outcomes over the derived set. The complete tagged gate result below is current evidence, not the original report's attestation.
- VAL-052: PASS for timing evidence - fresh tagged integration reports 297.291 s, below its 720 s warning and 900 s budget. Fresh uncached ordinary attempt reports cmd/fid 368.515 s, internal/workflow 100.094 s and slowest other ordinary package internal/issues 193.228 s. These are below 480 s warning and 600 s budget; the ordinary behavioral gate still failed.

## Tests run

Working directory for commands below: `/home/mark/projects/fid`, except direct binary commands explicitly marked `/tmp`. Every PASS receipt below was captured with observed command termination in this session. The original report quoted under Deviations is immutable historical evidence, not re-attested PASS evidence.

- `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid` - PASS - `cmd/fid/gate.go:272:func decodeClosedGateInput(stdin io.Reader, path string, target any) error {` (terminated `exit_code: 0`)
- `rg -n "func (same[A-Za-z]*Repository|validateRepositoryBinding|validEnvelopeRepository|validEnvelopeSource|validEvidenceRepository|validateRunTerminalRepository|verifyDecisionSource|importedGateDestination)" internal -g '!*_test.go'` - PASS - `internal/workflow/envelope.go:759:func validEnvelopeSource(source Source) bool {` (terminated `exit_code: 0`)

R5 currently has six matches: Review (`review.go:135`), Judgment (`concern.go:147`), Run (`run.go:226`), Activation (`activation.go:328`), Envelope (`envelope.go:178`) and Gate (`gate.go:272`). Review/Judgment wrappers remain around adopted descriptors; Run/Activation remain for their later adoption, and Envelope/Gate wrapper names remain after Task 6 adoption. The Claim helper is gone after Task 6. This is a successful accounted-match inventory, not an absence proof.

R6 currently has nine matches: Evidence validity (`evidence.go:668`) and equality (`evidence.go:672`), Gate imported destination (`execution.go:346`) and equality (`execution.go:431`), registration binding (`registration.go:362`), Run terminal validation (`run_close.go:290`) and equality (`run_read.go:149`), Envelope repository/source validity (`envelope.go:755,759`). Review's exact Repository comparison was displaced by Task 5's expectation check; Task 4's Decision Source verifier is gone. Evidence/Gate and Run observing paths remain for Tasks 8/7; retained registration/Envelope validity functions remain after Task 6. Shared SameWorktree is outside this regex's match set. No missing-match status was converted to silent success.

F5 focused red/green (byte-identical command):

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewSourceMatrixCompleteness$' ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewSourceMatrixCompleteness$' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.516s` (terminated `exit_code: 0`)

F6 focused red/green (byte-identical command):

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewReplayEvidenceRejectsChangedHistory$' ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewReplayEvidenceRejectsChangedHistory$' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.521s` (terminated `exit_code: 0`)

F7 help focused red/green (byte-identical command; final corruption guard covers every recursive field):

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewHelpEvidenceRejectsIncompleteMetadata$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewHelpEvidenceRejectsIncompleteMetadata$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.364s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewHelpEvidenceRejectsIncompleteMetadata$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.402s` (terminated `exit_code: 0`)

F7 structural red attempts (byte-identical command; no green for Review within scope):

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewAndFindingDescriptorHelpAndStructure$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewAndFindingDescriptorHelpAndStructure$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewAndFindingDescriptorHelpAndStructure$/^finding_address$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.449s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24ReviewSourceMatrix|TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation|TestP24ReviewReplayUnavailableAndRestoredControls)$' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.516s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24ReviewSourceMatrix|TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation|TestP24ReviewReplayUnavailableAndRestoredControls)$' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	11.195s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -timeout 30m ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP(23ReviewRecordRealProcessJourney|24Review.*)$' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	13.224s` (terminated `exit_code: 0`)

The ordinary touched-package attempt reported cmd/fid 244.954 s and the same OID-format failure. The tagged touched-file selection includes the entire retained Phase 23 Review journey and all added Review tests. None of these focused runs substitutes for the exact gates.

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`)

- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid review record --example` - PASS - `}` (terminated `exit_code: 0`)
- `/tmp/fid finding address --example` - PASS - `}` (terminated `exit_code: 0`)
- `/tmp/fid --json review record --input /tmp/fid-task5-f7-25jicau1/request.json --actor ACT-1 --idempotency-key task5-format-control` - FAIL - `{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["accepted review.record input"],"argument":"input","value":null},"message":"invalid review request","retryable":false,"subject":null},"ok":false,"operation":"review.record"}` (terminated `exit_code: 2`)

The three direct binary invocations ran from `/tmp`. Both ordered examples exited 0. The malformed-OID input reached the generic family fallback (`argument:"input"`, allowed `["accepted review.record input"]`) instead of the required exact OID field failure; this is failed runtime evidence corroborating F7, not a passing non-zero command.

## Exact quality-gate results

Working directory: `/home/mark/projects/fid`. After gofmt/goimports, the five exact commands ran in the declared order, one heavy command at a time. The test cache was cleared as setup to obtain fresh package times. No test or source change followed the start of these gates.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					61.0%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	297.291s` (terminated `exit_code: 0`)

The ordinary required command's failed attempt is immutable: malformed head OID validation returned nil at input_descriptor_test.go:969. BLOCKED here is the required-command status in a RECOVERY_REQUIRED report; it is not a protected Task outcome. No whole-suite pass is claimed. Coverage reporting succeeded on that attempt's profile and does not override the failed test gate.

## Deviations from spec

- No accepted contract, gate, assertion semantics, threshold, dependency or product behavior changed. F7 is unfinished because both causal production paths lie outside the exact repair allowlist.
- F5's first red ended `FAIL` (terminated `exit_code: 1`): the original rows exposed only the mismatching state. Completing the four-state list and descriptor-checked real-process matrix made the identical guard return `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.516s` (terminated `exit_code: 0`); all matrix rows passed the broader focused run.
- F6's first red ended `FAIL` (terminated `exit_code: 1`): the old replay predicate accepted changed Event/result history. Whole-envelope equality with only replay metadata changed made the identical guard return `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.521s` (terminated `exit_code: 0`); restored viable controls and independent repositories passed the broader runs.
- F7 help's first red ended `FAIL` (terminated `exit_code: 1`): the template/Git-only checker accepted corrupted recursive type/requirement metadata. The recursive metadata checker made the identical command pass, finally `ok  	github.com/mbell-nzl/fid/cmd/fid	2.402s` (terminated `exit_code: 0`).
- F7 structural attempt 1 ended `FAIL` (terminated `exit_code: 1`): integer minimum expectations used Go ints while the JSON validator returned float64 numbers, and a malformed OID was accepted. Numeric expectations now retain the exact JSON numeric value; the OID regression uses nonempty `not-an-oid`. Attempt 2's identical command still ended `FAIL` (terminated `exit_code: 1`) solely for the OID violation. Finding's full selection passed; Review's test and final ordinary failure remain red. No green was mined or substituted.
- The required ordinary gate's first repair-window attempt ended `FAIL` (terminated `exit_code: 1`) for the same nil OID validation. Root cause: `cmd/fid/input_descriptors_review.go:24` declares lowercase `oid`, while `internal/cli/input_validate.go:306` recognizes uppercase `OID` and returns true for unknown format names. Repair requires additional same-task production scope. There is no later passing ordinary attempt in this session.

- Draft hygiene check: `git diff --check -- cmd/fid/input_descriptor_test.go cmd/fid/integration/p23_review_record_test.go project-docs/docs/phases/phase-24/tasks/task-5-done.md` first failed with `+> ` (terminated `exit_code: 2`) because blank historical quote prefixes carried a trailing space. Changing those prefixes to a bare quote marker preserves the entire unquoted original; the byte-identical check then returned `exit 0 (no output)` (terminated `exit_code: 0`).

Original Task 5 report follows verbatim with a quote prefix only. It was finalized at `010a2550ed682d6445dc2c0f79bc12814b206860`, naming implementation `ff6c923dcb706513e71ff7dca76723d68da6ba89`. Its original PASS claims, failed attempts, provenance and commit fields are retained as history rather than current completion status.

> # Phase 24 - Task 5: Review leaves and repository expectations - Completion Report
>
> ## Workflow outcome
>
> - Task outcome: COMPLETE
> - Category: none
> - Condition: None.
> - Recovery attempts: 0
> - Next action: None.
> - Human attention required: no
>
> OUTCOME COMPLETE
>
> ## Summary of what was implemented
>
> - Added the shared `ExpectedRepository` pair and check. `review record` now requires that pair, compares only HEAD and dirty against one fid observation, stores the observed Repository, and returns the amended closed `repository-mismatch` details without an Event on a difference.
> - Converted Review coverage items, exclusions and evidence references to `SourceInput`. Resolution checks supplied project-file digests or observes omitted ones before coverage and digest-exclusivity rules. Replays return before another repository or Source observation.
> - Adopted descriptors for `review record` and `finding address`, including help, examples, exclusivity and structural errors. Converged displaced retained tests and fake/fid fixtures while preserving their assertions.
>
> ## Files created or modified
>
> - `internal/project/expected_repository.go`, `internal/project/expected_repository_test.go`
> - `internal/review/review.go`, `internal/review/review_test.go`
> - `internal/cli/codes.go`, `internal/cli/codes_test.go`
> - `cmd/fid/review.go`, `cmd/fid/finding.go`, `cmd/fid/input_descriptors_review.go`, `cmd/fid/input_descriptor_test.go`, `cmd/fid/review_test.go`
> - `cmd/fid/integration/p23_conformance_test.go`, `cmd/fid/integration/p23_failure_contract_test.go`, `cmd/fid/integration/p23_review_leaf_oracle_test.go`, `cmd/fid/integration/p23_review_record_test.go`, `cmd/fid/integration/p24_judgment_corpus_test.go`
> - `project-docs/specs/06_DESIGN_PLAN_CLI.md`, `project-docs/docs/phases/phase-24/tasks/task-5-done.md`
>
> ## Documentation updated
>
> - `project-docs/specs/06_DESIGN_PLAN_CLI.md` documents the expected pair, both Git commands, SourceInput positions and mismatch shape.
>
> ## Validation Contract Coverage
>
> - VAL-009: PASS for `review.record` - `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at `cmd/fid/integration/p23_review_record_test.go:121` advances HEAD, replays the original key without an Event, and checks a new key sees the changed HEAD.
> - VAL-011: PASS for both leaves - `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` at `cmd/fid/input_descriptor_test.go:235` checks the live registry/type tree and seeded mismatch; `TestP24ReviewAndFindingDescriptorHelpAndStructure` at line 277 checks both leaves left the pending fence.
> - VAL-012: PASS for both leaves - `TestP24ReviewAndFindingDescriptorHelpAndStructure` at `cmd/fid/input_descriptor_test.go:277` checks the help template line and both Git command meanings; `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` checks field conformance.
> - VAL-013: PASS for both leaves - `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:44` runs each standalone `--example` through the built CLI; direct CLI invocations below show the ordered JSON templates.
> - VAL-014: PASS for both leaves - `TestP24ReviewAndFindingDescriptorHelpAndStructure` at `cmd/fid/input_descriptor_test.go:277` validates and decodes each generated example into the closed Go request type.
> - VAL-015: PASS for both leaves - `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:44` exercises six conflicting inputs per leaf, and standalone success, with empty non-project workspaces.
> - VAL-016: PASS for both leaves - `TestP24ReviewAndFindingDescriptorHelpAndStructure` at `cmd/fid/input_descriptor_test.go:277` checks exact top-level and nested paths for unknown, duplicate, missing, null, type, format and enum failures through the CLI, with no created project state.
> - VAL-017: PASS for `review.record` - `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at `cmd/fid/integration/p23_review_record_test.go:45` rejects top-level and all three nested caller `repository` positions with their exact paths and unchanged Event counts.
> - VAL-018: PASS for `review.record` - the same test at `cmd/fid/integration/p23_review_record_test.go:77` checks mismatching supplied digests at item, exclusion and evidence positions with both digest values and no Event; its accepted Review at line 104 checks a matching supplied exclusion digest and all three stored resolved digests. All positions use the same `resolveSources` check.
> - VAL-019: PASS for `review.record` - the same test at `cmd/fid/integration/p23_review_record_test.go:104` checks omitted item and evidence digests against independently hashed file bytes; at line 154 it checks missing-file `source-digest-mismatch` with both digest fields null and no Event.
> - VAL-022: PASS for `review.record` - the same test at `cmd/fid/integration/p23_review_record_test.go:110` compares stored Review and all three Source repositories with the single fid observation and independently checks HEAD, dirty and tree; `TestP23ReviewRecordRealProcessJourney` at line 177 retains exact readback.
> - VAL-023: PASS for `review.record` - `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at `cmd/fid/integration/p23_review_record_test.go:135` checks new-key `repository-unavailable` in a non-Git state versus normal success at line 104, with unchanged Event count.
> - VAL-024: PASS for `review.record` - `TestP23ReviewRecordRealProcessJourney` at `cmd/fid/integration/p23_review_record_test.go:208` checks dirty and HEAD differences with no Events; its first matching Review succeeds at line 252.
> - VAL-025: PASS - the two mismatch subtests of `TestP23ReviewRecordRealProcessJourney` at `cmd/fid/integration/p23_review_record_test.go:208` and line 231 check expected echo and full observed Repository; `TestReasonCodeDetailSchemaIsExactForEveryCode` at `internal/cli/codes_test.go:38` pins the closed shape.
> - VAL-035: PASS for the two leaves - `TestP23SharedLeafConformance` at `cmd/fid/integration/p23_conformance_test.go:143` compares accepted, unknown-field and conflict fake/fid outcomes over the derived operation set, including both Review leaves.
> - VAL-052: PASS - the final ordinary gate reported `cmd/fid` 267.739 s, `internal/workflow` 115.695 s and slowest other ordinary package `internal/issues` 158.407 s; the tagged gate reported 252.256 s. None reached the 480 s ordinary warning, 600 s ordinary limit or 900 s tagged limit.
> - VAL-060: PASS for both leaves - `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:44` runs each standalone example from an empty non-project directory, checking exit 0, empty stderr and no created entries.
> - VAL-066: PASS for `review.record` - `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at `cmd/fid/integration/p23_review_record_test.go:135` replays without an observable Git repository, and at line 149 replays after deletion of a coverage file. New-key controls return `repository-unavailable` and `source-digest-mismatch` with no Event.
>
> ## Tests run
>
> - First TDD red: `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run 'TestP23ReviewRecordRealProcessJourney/dirty_expectation_mismatch_echoes_only_the_expected_pair' ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`); the old decoder rejected `expected_repository` before the repository check.
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -timeout 30m ./internal/project ./internal/review ./internal/cli ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`) on the first retained-fixture attempt; after convergence, the identical command passed with `ok  	github.com/mbell-nzl/fid/cmd/fid	194.986s` (terminated `exit_code: 0`).
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run 'TestP24JudgmentKeyedLeafConflictAndCanonicalRetryCorpus|TestP23ReviewRecordRealProcessJourney' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.435s` (terminated `exit_code: 0`)
> - `golangci-lint run ./...` - FAIL - `* staticcheck: 1` (terminated `exit_code: 1`) on the first exact-gate attempt; the corrected command passed as recorded below.
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`) on its first exact-gate attempt; the corrected command passed as recorded below.
> - `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> - `/tmp/fid review record --example` - PASS - `}` (terminated `exit_code: 0`)
> - `/tmp/fid finding address --example` - PASS - `}` (terminated `exit_code: 0`)
>
> Both direct CLI invocations ran from `/tmp`.
>
> ## Exact quality-gate results
>
> Working directory: `/home/mark/projects/fid`. The five required commands were run in task-spec order. After causal repairs, only failed commands were repeated; completed passing checks were retained.
>
> - `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> - `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	115.695s	coverage: 62.6% of statements` (terminated `exit_code: 0`)
> - `go tool cover -func=coverage.out` - PASS - `total:											(statements)				60.7%` (terminated `exit_code: 0`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	252.256s` (terminated `exit_code: 0`)
>
> ## Deviations from spec
>
> - No accepted behavior, protected assertion, command or threshold changed. The first lint attempt ended `* staticcheck: 1` (terminated `exit_code: 1`) because a new test declared `var typed any =`; changing it to `typed :=` made the identical lint command return `0 issues.` (terminated `exit_code: 0`).
> - The first exact tagged gate ended `FAIL` (terminated `exit_code: 1`). Task 5's new verdict descriptor rejected a retained corpus mutation's invalid enum, and two retained mismatch comparisons treated Go typed nil differently from JSON null. The corpus now changes `changes-required` to the accepted `blocked` verdict, preserving its idempotency-conflict assertion; the comparisons normalize the fixture to JSON while retaining full Repository equality. The identical tagged command later returned `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	252.256s` (terminated `exit_code: 0`).
> - The touched-package race attempt ended `FAIL` (terminated `exit_code: 1`) because retained Review/Finding fixture inputs still used displaced shapes. Converging those inputs produced `ok  	github.com/mbell-nzl/fid/cmd/fid	194.986s` (terminated `exit_code: 0`) on the identical command.
> - The first report-lint attempt ended `FAIL report-lint found 4 failure(s)` (terminated `exit_code: 4`) because four passing command rows had punctuation after the termination field. Removing those periods preserved the recorded tokens and statuses.
>
> ## Concerns, technical debt, or notes for the Tech Lead
>
> - Retained assertions changed only where Task 5 displaced their input or observation model. Introducing commits were `d215e40` and `926ef72` (phase-23/task-4 Review journeys), `d1c016d` (phase-23/task-13 failure corpus), `3335149` and `2e7fd2f` (phase-23/task-18 conformance and leaf oracle), `472ed64` (phase-23/task-13 code-schema assertion), and `9bbbdb5` (phase-24/task-2 keyed conflict corpus). Their subjects and no-Event checks remain intact.
> - No package reached the 80% timing warning.
>
> ## Dependency-budget note
>
> - No dependency changes.
>
> ## Commit hash
>
> - `ff6c923dcb706513e71ff7dca76723d68da6ba89`
>
> ## Commit subject
>
> - `feat(review): check repository expectations and resolve sources`

## Concerns, technical debt, or notes for the Tech Lead

- Same-task scope gap: Task 5 commit `ff6c923dcb706513e71ff7dca76723d68da6ba89` introduced the lowercase expected head OID format (`git log -1 -- cmd/fid/input_descriptors_review.go`). The shared uppercase OID validator branch was introduced by Task 3 commit `04b6b0992fe20a02df402a183269a9847ddae5ab` (`git log -1 -S 'case "OID":' -- internal/cli/input_validate.go`). The Task 5 adoption exposed this mismatch; this report does not transfer repair ownership. Neither production path is in this dispatch's allowlist. Obtain bounded same-task repair scope and preserve VAL-014 placeholder typing while establishing actual valid OID structural controls and exact malformed-OID rejection.
- All new test assertions strengthen retained Review/Finding evidence. The original report's introducing-commit provenance remains quoted above. No production, ledger, triage, control, contract or gate file was modified by this child.
- No package reached its 80% warning: cmd/fid 368.515 s, internal/workflow 100.094 s, slowest other ordinary internal/issues 193.228 s (<480 s warning, <600 s budget); tagged integration 297.291 s (<720 s warning, <900 s budget). These actual timings do not override the failed ordinary gate.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
