# Phase 24 - Task 5: Review leaves and repository expectations - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 2
- Next action: None.
- Human attention required: no

OUTCOME COMPLETE

Developer completion under TRG-24-005. Dispatch binding: TASK-RECOVERY `2026-10-06 22:46`, task-5, trigger `review-3:qa-rework:task-5`, attempt `2`, category `qa-rework`, mode `task-scope`, authority `QA review`. The five identity bindings, ordered Allowed files, Forbidden changes, Exact verification, Originating report and Control repair base were compared byte-for-byte to TRG-24-005. This is ordinary product task recovery; control fields remain N/A. Starting HEAD was `a0ca854133ae98b77e36cf51ad596ff0e8c6e096`, preserving attempt 1's implementation/report commit `04246029c4488eec46eba76323dee7fbe8260c8d`. Independent QA re-review remains the Runner's separate next checkpoint.

## Summary of what was implemented

- The exact `oid` and `OID` format labels now dispatch to the existing 40-character lowercase hexadecimal OID predicate. Unknown format labels reject. No case folding, normalization, additional OID shape or runtime placeholder exemption was added.
- Preserved Review's lowercase `oid` descriptor, help/meaning/requirement/placeholder and error metadata without editing the descriptor file. Added valid, malformed, empty, short, nonhex, uppercase-hex, alternate-length and placeholder controls for both dispatch labels, with exact InputError and Details equality. Synthetic unknown and mixed-case format names fail closed.
- Review/Finding examples retain their original bytes and placeholders. Their evidence checks required members, closed keys and recursive JSON types independently of runtime formats, decodes into their Go requests with unknown fields disallowed, and exercises the built CLI outside a project without state changes. Runtime Review structural controls substitute the real checkpoint commit OID.
- Retained the malformed-OID regression and complete recursive help/structural corpus. Added built-process Review OID refusals checking the precise closed failure envelope, field/value/allowed metadata, exit 2 and empty workspace. Attempt 1's F4-F6 fixtures and all earlier F4-F7 evidence/failures remain intact.

## Files created or modified

- `internal/cli/input_validate.go`
- `cmd/fid/input_descriptor_test.go`
- `project-docs/docs/phases/phase-24/tasks/task-5-done.md`

## Documentation updated

- This report records the bounded repair, own-session evidence and completion; the full prior report is preserved under Deviations as quoted historical evidence. No spec, task, QA, control, ledger or triage artifact was edited.

## Validation Contract Coverage

Current evidence below is from this session's passing ordinary and tagged gates, the focused selections and direct CLI observations. Quoted historical PASS rows remain attributable to their original sessions.

- VAL-009: PASS - `TestP24ReviewReplayUnavailableAndRestoredControls` at `cmd/fid/integration/p23_review_record_test.go:207` compares whole original response/Event and unchanged Event count after HEAD advances; the fresh control refuses, then succeeds after restoration. `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at line 298 additionally changes referenced bytes before the advanced-HEAD retry. Both ran in the exact tagged gate.
- VAL-011: PASS - `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` at `cmd/fid/input_descriptor_test.go:236` checks registry/live-tree/request-type agreement and seeded mismatch rejection. `TestP24ReviewAndFindingDescriptorHelpAndStructure` at line 731 checks both adopted leaves left the pending fence. The exact ordinary gate passed.
- VAL-012: PASS - `TestP24ReviewHelpEvidenceRejectsIncompleteMetadata` at `cmd/fid/input_descriptor_test.go:596` checks both valid help outputs and rejects independently corrupted recursive type, requirement, meaning and applicable constraints. The Review/Finding structure test at line 731 checks exact help metadata and the two Git commands. The final focused selection and exact ordinary gate passed.
- VAL-013, VAL-014, VAL-060: PASS - the two-leaf structure test at `cmd/fid/input_descriptor_test.go:731` compares built-process stdout byte-for-byte with ordered templates, recursively checks required-only closed keys and JSON types via `p24ReviewTemplateObject` at line 874, decodes into the Go requests with `DisallowUnknownFields`, preserves `<head-oid>`, and checks exit 0/empty stderr/empty non-project directories. Direct binary examples below also preserve the original emitted form. Executable structure uses a real OID; synthetic runtime controls reject the placeholder.
- VAL-015: PASS - `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:45` retains both leaves' standalone success and all six conflicting-input process rows, exact details/exit 2 and no state. It ran in the exact ordinary gate.
- VAL-016: PASS - `p24ReviewRecursiveStructuralDetails` at `cmd/fid/input_descriptor_test.go:1128`, exercised by the two-leaf structure test, retains recursive missing/null/type/unknown/duplicate/enum/format/range/cardinality cases, exact details and valid controls with no effects. Its Review baseline now contains a real qualifying OID; `not-an-oid` still has argument `input.expected_repository.head_oid`, value `not-an-oid`, allowed `["oid"]`. `TestP24ReviewMalformedOIDProcess` at line 693 checks four actual process refusals with the closed envelope, invalid-input/invalid-argument, retryable false, subject null, exact message, exit 2 and no workspace/Event effects. `TestP24OIDFormatDispatch` at line 650 proves both labels' valid and invalid values, and `TestP24UnknownInputFormatRejects` at line 634 proves unknown-label rejection with exact InputError details. Focused green and the exact ordinary gate passed.
- VAL-017: PASS - `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at `cmd/fid/integration/p23_review_record_test.go:298` retains all four removed caller-repository positions, exact paths and unchanged Events. It ran in the exact tagged gate.
- VAL-018, VAL-019: PASS - `TestP24ReviewSourceMatrix` at `cmd/fid/integration/p23_review_record_test.go:104` derives/checks all three Source paths from real help and exercises twelve position/state rows. Matching/omitted-existing acceptance independently hashes stored bytes; mismatching/omitted-missing refusals check exact supplied/observed or null digest details and no Event. The exact tagged gate passed without fixture edits.
- VAL-022: PASS - `p24ReviewAssertRepositories` at `cmd/fid/integration/p23_review_record_test.go:167` compares complete Review, all three Source and journal Event repositories with independent clean observations in the matrix and restored-control tests. `p24ReviewDirtyRepository` at line 183 independently hashes known dirty file images and Git/config/root identity for the retained dirty journey at line 298. Base, HEAD, tree, dirty, dirty digest, capture source, project and root identity all participate. The exact tagged gate passed.
- VAL-023, VAL-066: PASS - `TestP24ReviewReplayUnavailableAndRestoredControls` at `cmd/fid/integration/p23_review_record_test.go:207` retains normal/non-Git/unmerged/deleted-source pairs, whole-response/original-Event replay equality, exact unavailable or missing-Source refusals with no Event, and identical fresh Run/key/request success after restoration. The exact tagged gate passed.
- VAL-024, VAL-025: PASS - `TestP23ReviewRecordRealProcessJourney` at `cmd/fid/integration/p23_review_record_test.go:460` retains matching, dirty-different and HEAD-different expectations, complete expected/observed details and no Event on refusal; the exact tagged gate passed. `TestReasonCodeDetailSchemaIsExactForEveryCode` at `internal/cli/codes_test.go:38` retains the closed shape and ran in the exact ordinary gate.
- VAL-035: PASS - `TestP23SharedLeafConformance` at `cmd/fid/integration/p23_conformance_test.go:143` compares fake/fid acceptance, unknown-field rejection and conflict cases over the derived set, including both Task 5 leaves. The exact tagged gate passed.
- VAL-052: PASS - fresh ordinary times: cmd/fid 366.509 s, internal/workflow 101.002 s, slowest other internal/issues 190.218 s. Tagged integration 313.544 s. All are below 600/900 s budgets and 480/720 s warnings. These are gate-reported durations, not cached or focused substitutes.

## Tests run

Working directory: `/home/mark/projects/fid` unless a direct binary working directory is explicitly stated. Every new receipt below was captured in this attempt with observed termination. All earlier failure/pass receipts remain in the preserved attempt-1 report below.

F7 structural TDD, byte-identical red/green:

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewAndFindingDescriptorHelpAndStructure$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewAndFindingDescriptorHelpAndStructure$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	4.524s` (terminated `exit_code: 0`)

Unknown-format TDD, byte-identical red/green:

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24UnknownInputFormatRejects$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24UnknownInputFormatRejects$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.355s` (terminated `exit_code: 0`)

OID-dispatch TDD, byte-identical red/green:

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24OIDFormatDispatch$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24OIDFormatDispatch$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.373s` (terminated `exit_code: 0`)

Final focused corpus and touched-package pre-gate:

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24ReviewAndFindingDescriptorHelpAndStructure|TestP24ReviewHelpEvidenceRejectsIncompleteMetadata|TestP24ReviewMalformedOIDProcess|TestP24OIDFormatDispatch|TestP24UnknownInputFormatRejects)$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	4.740s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/cli ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	245.024s` (terminated `exit_code: 0`)

The touched-package command also reported internal/cli 1.115 s. No tagged test was edited in this continuation. Focused selections do not substitute for the five exact gates.

Fresh F4 inventory receipts (no VAL):

- `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid` - PASS - `cmd/fid/gate.go:272:func decodeClosedGateInput(stdin io.Reader, path string, target any) error {` (terminated `exit_code: 0`)
- `rg -n "func (same[A-Za-z]*Repository|validateRepositoryBinding|validEnvelopeRepository|validEnvelopeSource|validEvidenceRepository|validateRunTerminalRepository|verifyDecisionSource|importedGateDestination)" internal -g '!*_test.go'` - PASS - `internal/workflow/envelope.go:759:func validEnvelopeSource(source Source) bool {` (terminated `exit_code: 0`)

R5 has the same six matches: Review, Judgment, Run, Activation, Envelope and Gate. Adopted wrappers remain, later Run/Activation owners have not removed theirs, and Task 6 removed Claim. R6 has the same nine matches: Evidence validity/equality, Gate imported destination/equality/registration binding, Run terminal validity/equality, and Envelope repository/Source validity. Task 5 displaced Review's exact Repository comparison with expectations; Task 4 removed its Decision Source verifier. Remaining observing paths belong to the accepted later slices or retained binding checks. These are accounted-match commands with observed exit 0, not absence proofs. Attempt 1's full accounting is preserved below.

Real entrypoint evidence:

- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid review record --example` - PASS - `}` (terminated `exit_code: 0`)
- `/tmp/fid finding address --example` - PASS - `}` (terminated `exit_code: 0`)

Both examples ran in the fresh non-project directory `/tmp/fid-task5-attempt2-Fs7Kir`. The malformed-OID invocation ran in its empty `workspace` subdirectory, with the input file outside that workspace. Its observed refusal is runtime evidence, not a non-zero gate PASS:

- `/tmp/fid --json review record --input /tmp/fid-task5-attempt2-Fs7Kir/request.json --actor ACT-1 --idempotency-key task5-format-control` - observed structural refusal - `{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["oid"],"argument":"input.expected_repository.head_oid","value":"not-an-oid"},"message":"invalid input.expected_repository.head_oid","retryable":false,"subject":null},"ok":false,"operation":"review.record"}` (terminated `exit_code: 2`)

A Python file-state assertion observed an empty workspace before the invocation (`workspace empty before Review refusal`, terminated `exit_code: 0`) and again afterward (`workspace unchanged: no files or Events`, terminated `exit_code: 0`). The permanent built-process test above asserts the same state and exact envelope for malformed/empty/short/nonhex values.

## Exact quality-gate results

Working directory: `/home/mark/projects/fid`. Gofmt/goimports were applied before the gates. The test cache was cleared as setup for fresh package timings. All five byte-identical Task 5 commands ran in this exact order, with one heavy command at a time, on the final source/test tree. No source or test change followed the start of the gates. The QA plan requires no additional task-applicable gate beyond these and the owned inventories above.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	101.002s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					61.0%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	313.544s` (terminated `exit_code: 0`)

Ordinary gate: cmd/fid 366.509 s, internal/workflow 101.002 s, slowest other internal/issues 190.218 s; next slowest internal/export 177.391 s and internal/journal 173.714 s. Tagged gate: 313.544 s. No package reached its 80% warning or budget. Coverage reporting does not replace either test gate; both terminated successfully here.

## Deviations from spec

- No protected contract, assertion semantics, accepted status/scope, command, coverage policy, threshold/budget, dependency or domain rule changed. TRG-24-005 authorizes the exact bounded source and evidence repair; its four-path allowlist was observed, with the unchanged Review descriptor requiring no edit.
- Immutable first ordinary-gate failure from attempt 1: `FAIL` (terminated `exit_code: 1`) at input_descriptor_test.go:969, because validation returned nil for `not-an-oid`. Introducing provenance: Task 5's lowercase descriptor was introduced by `ff6c923dcb706513e71ff7dca76723d68da6ba89` (phase-24/task-5); the uppercase-only branch was introduced by `04b6b0992fe20a02df402a183269a9847ddae5ab` (phase-24/task-3). Own-session `git log -1 -- cmd/fid/input_descriptors_review.go` and `git log -1 -S 'case "OID":' -- internal/cli/input_validate.go` confirmed both. Task 5's accepted adoption exposed the mismatch; no repair ownership is transferred to another task.
- Attempt 2 first captured the mandated unchanged focused structural red, ending `FAIL` (terminated `exit_code: 1`) for the exact same malformed-OID acceptance, before production edits. Adding the lowercase label to the existing uppercase predicate restored rejection; preserving typed template checks while substituting a real OID for executable controls made the byte-identical structural command return `ok  	github.com/mbell-nzl/fid/cmd/fid	4.524s` (terminated `exit_code: 0`).
- Separate unknown-format red ended `FAIL` (terminated `exit_code: 1`) because synthetic unknown-format/Oid/oID descriptors accepted a valid-looking OID. Changing only matchesInputFormat's unknown-format fallback to false made the identical command return `ok  	github.com/mbell-nzl/fid/cmd/fid	2.355s` (terminated `exit_code: 0`).
- Separate OID-dispatch red ended `FAIL` (terminated `exit_code: 1`) for lowercase malformed/empty/short/nonhex/uppercase-hex/alternate-length/placeholder values while uppercase controls already held. The shared existing predicate made the identical command return `ok  	github.com/mbell-nzl/fid/cmd/fid	2.373s` (terminated `exit_code: 0`). No permissive aliases or value expansion were added.
- The byte-identical required ordinary gate now ends `ok  	github.com/mbell-nzl/fid/internal/workflow	101.002s	coverage: 62.1% of statements` (terminated `exit_code: 0`); the five exact current gates above are passing close evidence. Earlier failed attempts remain immutable below. The direct malformed-OID CLI now returns the exact structural envelope above with observed exit 2, replacing its prior generic fallback observation.

Attempt 1 report preserved verbatim with one quote prefix per line, including all earlier F4-F7 evidence, failures, timing, introducing provenance and original Task 5 history. Its PASS rows are historical receipts, not re-attested results for this attempt:

> # Phase 24 - Task 5: Review leaves and repository expectations - Recovery Report
>
> ## Workflow outcome
>
> - Task outcome: RECOVERY_REQUIRED
> - Category: qa-rework
> - Condition: F7's declared Review OID-format rejection fails VAL-016 and the ordinary required gate. The production descriptor/validator repair is outside this dispatch's exact three-file scope.
> - Recovery attempts: 1
> - Next action: Obtain bounded same-task production scope for cmd/fid/input_descriptors_review.go and internal/cli/input_validate.go, with their qualifying controls and tests; enforce the already accepted OID contract, preserve the red regression, and rerun focused red/green and the byte-identical five-command gate before completion.
> - Human attention required: no
>
> OUTCOME RECOVERY_REQUIRED
>
> Dispatch binding: latest TASK-RECOVERY `2026-10-06 22:12`, task-5, trigger `review-3:qa-rework:task-5`, attempt `1`, category `qa-rework`, mode `task-scope`, authority `QA review`. Its ordered Allowed files and Forbidden changes were checked against the dispatch as exact copies. No ledger or triage edits were made by this Developer.
>
> ## Summary of what was implemented
>
> - F4: captured both accepted R5/R6 inventory commands in this session, with actual match accounting and observed exit status.
> - F5: added descriptor-checked three-position Review Source matrices, each covering supplied matching/mismatching and omitted existing/missing digests. Every row has a fresh valid Review Run; acceptance checks independently hashed stored digests, refusal checks exact digest/null details and no Event.
> - F6: independently compared complete clean/dirty Review, three Source and Event repositories, including base, config project identity and root identity. Added complete original-envelope and journal-Event replay equality, normal/non-Git/unmerged/deleted-source pairs, and fresh unreviewed Run/new-key refusal controls whose identical restored requests succeed. Preserved the original Review journey and all previous assertions.
> - F7: strengthened recursive help metadata and structural cases only for Review/Finding. Both help guards and Finding's structural corpus pass. Review's malformed expected head OID is now a retained failing assertion; completion requires an authorized same-task production repair.
>
> ## Files created or modified
>
> - `cmd/fid/integration/p23_review_record_test.go`
> - `cmd/fid/input_descriptor_test.go`
> - `project-docs/docs/phases/phase-24/tasks/task-5-done.md`
>
> ## Documentation updated
>
> - This report records the Window B repair evidence, remaining scope gap and original report verbatim as quoted history. No accepted product documentation changed.
>
> ## Validation Contract Coverage
>
> - VAL-009: PASS for the repaired Review evidence - `TestP24ReviewReplayUnavailableAndRestoredControls` at `cmd/fid/integration/p23_review_record_test.go:207`, head-advanced row, compares the whole original response allowing only `result.replayed`, plus the actual original Event and unchanged Event count. Its fresh control refuses the advanced HEAD and succeeds after restoration. `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at line 298 additionally changes referenced file bytes before advancing HEAD and preserves full replay equality.
> - VAL-011: PASS for retained Review/Finding adoption evidence - `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` at `cmd/fid/input_descriptor_test.go:236` checks live types and seeded mismatch; `TestP24ReviewAndFindingDescriptorHelpAndStructure` at line 634 checks both pending-fence removals before the later red format row.
> - VAL-012: PASS - `TestP24ReviewHelpEvidenceRejectsIncompleteMetadata` at `cmd/fid/input_descriptor_test.go:596` checks both valid help outputs and independently corrupts every recursive field's type, requirement, meaning and applicable format/member/range/length metadata; every corruption is rejected. Review help retains both Git command meanings. The identical focused command passed after the red guard.
> - VAL-013, VAL-015, VAL-060: PASS for retained two-leaf evidence - `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:45` exercises built-process standalone examples and conflicting inputs in empty non-project directories; it had no failure in the full ordinary attempt. Direct examples in Tests run additionally exercise the actual binary.
> - VAL-014: PASS for retained typed-template evidence - both leaves' template generation and closed Go-request decode in `TestP24ReviewAndFindingDescriptorHelpAndStructure` at `cmd/fid/input_descriptor_test.go:634` execute before the format failure. VAL-014 explicitly permits placeholders that are not valid OIDs; the continuation must preserve that template meaning while using actual qualifying values for format controls.
> - VAL-016: FAIL for Review, PASS for Finding - `p24ReviewRecursiveStructuralDetails` at `cmd/fid/input_descriptor_test.go:946` derives recursive missing/null/type/unknown/duplicate/enum/format/range/cardinality cases with exact path/value/allowed and closed error metadata, valid controls and empty-workspace/no-Event checks. Finding's entire focused selection passes. Review's `recursive/input.expected_repository.head_oid/format` fails at line 969: validation is nil for `not-an-oid`, expected path `input.expected_repository.head_oid`, value `not-an-oid`, allowed `["oid"]`. Other Review structural assertions remain present; the whole leaf is not claimed green.
> - VAL-017: PASS for retained Review rejection evidence - `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at `cmd/fid/integration/p23_review_record_test.go:298` retains top-level and all three nested caller-repository refusals with exact paths and unchanged Events.
> - VAL-018, VAL-019: PASS - `TestP24ReviewSourceMatrix` at `cmd/fid/integration/p23_review_record_test.go:104` derives/checks all three Source paths from real descriptor help and executes all twelve position/state rows. Matching and omitted-existing each compare the stored digest with independently hashed actual file bytes; mismatching checks supplied/actual digests and omitted-missing checks both nulls, each with exact closed details and no emitted/appended Event.
> - VAL-022: PASS - `p24ReviewAssertRepositories` at `cmd/fid/integration/p23_review_record_test.go:167` compares full Review, all three Source and actual journal Event repositories against independent Git/config/root observations. Clean matrix/lifecycle acceptance uses `p24JudgmentExpectedRepository`; the original dirty journey uses independent known-file canonical-image hashing in `p24ReviewDirtyRepository` at line 183. All Repository fields, including base/config/root identity and dirty digest, participate.
> - VAL-023, VAL-066: PASS - `TestP24ReviewReplayUnavailableAndRestoredControls` at `cmd/fid/integration/p23_review_record_test.go:207` has normal/non-Git/unmerged/deleted-source states. Each unavailable/deleted replay preserves the complete result and original Event with no Event appended. A fresh unreviewed Run and new key fail with exact unavailable reason/detail (and precondition/retryable metadata), or missing-Source null digest details, and no Event; the identical request/Run/key succeeds in the restored observable state, appending exactly one Event with independently checked repositories.
> - VAL-024, VAL-025: PASS for retained Review evidence - `TestP23ReviewRecordRealProcessJourney` at `cmd/fid/integration/p23_review_record_test.go:460` retains matching/HEAD-different/dirty-different rows and full expected/observed mismatch details with no Event on refusal; `TestReasonCodeDetailSchemaIsExactForEveryCode` at `internal/cli/codes_test.go:38` retains the closed shape assertion.
> - VAL-035: PASS for retained two-leaf evidence - `TestP23SharedLeafConformance` at `cmd/fid/integration/p23_conformance_test.go:143` compares accepted/unknown/conflict fake/fid outcomes over the derived set. The complete tagged gate result below is current evidence, not the original report's attestation.
> - VAL-052: PASS for timing evidence - fresh tagged integration reports 297.291 s, below its 720 s warning and 900 s budget. Fresh uncached ordinary attempt reports cmd/fid 368.515 s, internal/workflow 100.094 s and slowest other ordinary package internal/issues 193.228 s. These are below 480 s warning and 600 s budget; the ordinary behavioral gate still failed.
>
> ## Tests run
>
> Working directory for commands below: `/home/mark/projects/fid`, except direct binary commands explicitly marked `/tmp`. Every PASS receipt below was captured with observed command termination in this session. The original report quoted under Deviations is immutable historical evidence, not re-attested PASS evidence.
>
> - `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid` - PASS - `cmd/fid/gate.go:272:func decodeClosedGateInput(stdin io.Reader, path string, target any) error {` (terminated `exit_code: 0`)
> - `rg -n "func (same[A-Za-z]*Repository|validateRepositoryBinding|validEnvelopeRepository|validEnvelopeSource|validEvidenceRepository|validateRunTerminalRepository|verifyDecisionSource|importedGateDestination)" internal -g '!*_test.go'` - PASS - `internal/workflow/envelope.go:759:func validEnvelopeSource(source Source) bool {` (terminated `exit_code: 0`)
>
> R5 currently has six matches: Review (`review.go:135`), Judgment (`concern.go:147`), Run (`run.go:226`), Activation (`activation.go:328`), Envelope (`envelope.go:178`) and Gate (`gate.go:272`). Review/Judgment wrappers remain around adopted descriptors; Run/Activation remain for their later adoption, and Envelope/Gate wrapper names remain after Task 6 adoption. The Claim helper is gone after Task 6. This is a successful accounted-match inventory, not an absence proof.
>
> R6 currently has nine matches: Evidence validity (`evidence.go:668`) and equality (`evidence.go:672`), Gate imported destination (`execution.go:346`) and equality (`execution.go:431`), registration binding (`registration.go:362`), Run terminal validation (`run_close.go:290`) and equality (`run_read.go:149`), Envelope repository/source validity (`envelope.go:755,759`). Review's exact Repository comparison was displaced by Task 5's expectation check; Task 4's Decision Source verifier is gone. Evidence/Gate and Run observing paths remain for Tasks 8/7; retained registration/Envelope validity functions remain after Task 6. Shared SameWorktree is outside this regex's match set. No missing-match status was converted to silent success.
>
> F5 focused red/green (byte-identical command):
>
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewSourceMatrixCompleteness$' ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewSourceMatrixCompleteness$' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.516s` (terminated `exit_code: 0`)
>
> F6 focused red/green (byte-identical command):
>
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewReplayEvidenceRejectsChangedHistory$' ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewReplayEvidenceRejectsChangedHistory$' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.521s` (terminated `exit_code: 0`)
>
> F7 help focused red/green (byte-identical command; final corruption guard covers every recursive field):
>
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewHelpEvidenceRejectsIncompleteMetadata$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewHelpEvidenceRejectsIncompleteMetadata$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.364s` (terminated `exit_code: 0`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewHelpEvidenceRejectsIncompleteMetadata$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.402s` (terminated `exit_code: 0`)
>
> F7 structural red attempts (byte-identical command; no green for Review within scope):
>
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewAndFindingDescriptorHelpAndStructure$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewAndFindingDescriptorHelpAndStructure$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReviewAndFindingDescriptorHelpAndStructure$/^finding_address$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.449s` (terminated `exit_code: 0`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24ReviewSourceMatrix|TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation|TestP24ReviewReplayUnavailableAndRestoredControls)$' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.516s` (terminated `exit_code: 0`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24ReviewSourceMatrix|TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation|TestP24ReviewReplayUnavailableAndRestoredControls)$' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	11.195s` (terminated `exit_code: 0`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -timeout 30m ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP(23ReviewRecordRealProcessJourney|24Review.*)$' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	13.224s` (terminated `exit_code: 0`)
>
> The ordinary touched-package attempt reported cmd/fid 244.954 s and the same OID-format failure. The tagged touched-file selection includes the entire retained Phase 23 Review journey and all added Review tests. None of these focused runs substitutes for the exact gates.
>
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`)
>
> - `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> - `/tmp/fid review record --example` - PASS - `}` (terminated `exit_code: 0`)
> - `/tmp/fid finding address --example` - PASS - `}` (terminated `exit_code: 0`)
> - `/tmp/fid --json review record --input /tmp/fid-task5-f7-25jicau1/request.json --actor ACT-1 --idempotency-key task5-format-control` - FAIL - `{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["accepted review.record input"],"argument":"input","value":null},"message":"invalid review request","retryable":false,"subject":null},"ok":false,"operation":"review.record"}` (terminated `exit_code: 2`)
>
> The three direct binary invocations ran from `/tmp`. Both ordered examples exited 0. The malformed-OID input reached the generic family fallback (`argument:"input"`, allowed `["accepted review.record input"]`) instead of the required exact OID field failure; this is failed runtime evidence corroborating F7, not a passing non-zero command.
>
> ## Exact quality-gate results
>
> Working directory: `/home/mark/projects/fid`. After gofmt/goimports, the five exact commands ran in the declared order, one heavy command at a time. The test cache was cleared as setup to obtain fresh package times. No test or source change followed the start of these gates.
>
> - `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> - `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
> - `go tool cover -func=coverage.out` - PASS - `total:											(statements)					61.0%` (terminated `exit_code: 0`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	297.291s` (terminated `exit_code: 0`)
>
> The ordinary required command's failed attempt is immutable: malformed head OID validation returned nil at input_descriptor_test.go:969. BLOCKED here is the required-command status in a RECOVERY_REQUIRED report; it is not a protected Task outcome. No whole-suite pass is claimed. Coverage reporting succeeded on that attempt's profile and does not override the failed test gate.
>
> ## Deviations from spec
>
> - No accepted contract, gate, assertion semantics, threshold, dependency or product behavior changed. F7 is unfinished because both causal production paths lie outside the exact repair allowlist.
> - F5's first red ended `FAIL` (terminated `exit_code: 1`): the original rows exposed only the mismatching state. Completing the four-state list and descriptor-checked real-process matrix made the identical guard return `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.516s` (terminated `exit_code: 0`); all matrix rows passed the broader focused run.
> - F6's first red ended `FAIL` (terminated `exit_code: 1`): the old replay predicate accepted changed Event/result history. Whole-envelope equality with only replay metadata changed made the identical guard return `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.521s` (terminated `exit_code: 0`); restored viable controls and independent repositories passed the broader runs.
> - F7 help's first red ended `FAIL` (terminated `exit_code: 1`): the template/Git-only checker accepted corrupted recursive type/requirement metadata. The recursive metadata checker made the identical command pass, finally `ok  	github.com/mbell-nzl/fid/cmd/fid	2.402s` (terminated `exit_code: 0`).
> - F7 structural attempt 1 ended `FAIL` (terminated `exit_code: 1`): integer minimum expectations used Go ints while the JSON validator returned float64 numbers, and a malformed OID was accepted. Numeric expectations now retain the exact JSON numeric value; the OID regression uses nonempty `not-an-oid`. Attempt 2's identical command still ended `FAIL` (terminated `exit_code: 1`) solely for the OID violation. Finding's full selection passed; Review's test and final ordinary failure remain red. No green was mined or substituted.
> - The required ordinary gate's first repair-window attempt ended `FAIL` (terminated `exit_code: 1`) for the same nil OID validation. Root cause: `cmd/fid/input_descriptors_review.go:24` declares lowercase `oid`, while `internal/cli/input_validate.go:306` recognizes uppercase `OID` and returns true for unknown format names. Repair requires additional same-task production scope. There is no later passing ordinary attempt in this session.
>
> - Draft hygiene check: `git diff --check -- cmd/fid/input_descriptor_test.go cmd/fid/integration/p23_review_record_test.go project-docs/docs/phases/phase-24/tasks/task-5-done.md` first failed with `+> ` (terminated `exit_code: 2`) because blank historical quote prefixes carried a trailing space. Changing those prefixes to a bare quote marker preserves the entire unquoted original; the byte-identical check then returned `exit 0 (no output)` (terminated `exit_code: 0`).
>
> Original Task 5 report follows verbatim with a quote prefix only. It was finalized at `010a2550ed682d6445dc2c0f79bc12814b206860`, naming implementation `ff6c923dcb706513e71ff7dca76723d68da6ba89`. Its original PASS claims, failed attempts, provenance and commit fields are retained as history rather than current completion status.
>
> > # Phase 24 - Task 5: Review leaves and repository expectations - Completion Report
> >
> > ## Workflow outcome
> >
> > - Task outcome: COMPLETE
> > - Category: none
> > - Condition: None.
> > - Recovery attempts: 0
> > - Next action: None.
> > - Human attention required: no
> >
> > OUTCOME COMPLETE
> >
> > ## Summary of what was implemented
> >
> > - Added the shared `ExpectedRepository` pair and check. `review record` now requires that pair, compares only HEAD and dirty against one fid observation, stores the observed Repository, and returns the amended closed `repository-mismatch` details without an Event on a difference.
> > - Converted Review coverage items, exclusions and evidence references to `SourceInput`. Resolution checks supplied project-file digests or observes omitted ones before coverage and digest-exclusivity rules. Replays return before another repository or Source observation.
> > - Adopted descriptors for `review record` and `finding address`, including help, examples, exclusivity and structural errors. Converged displaced retained tests and fake/fid fixtures while preserving their assertions.
> >
> > ## Files created or modified
> >
> > - `internal/project/expected_repository.go`, `internal/project/expected_repository_test.go`
> > - `internal/review/review.go`, `internal/review/review_test.go`
> > - `internal/cli/codes.go`, `internal/cli/codes_test.go`
> > - `cmd/fid/review.go`, `cmd/fid/finding.go`, `cmd/fid/input_descriptors_review.go`, `cmd/fid/input_descriptor_test.go`, `cmd/fid/review_test.go`
> > - `cmd/fid/integration/p23_conformance_test.go`, `cmd/fid/integration/p23_failure_contract_test.go`, `cmd/fid/integration/p23_review_leaf_oracle_test.go`, `cmd/fid/integration/p23_review_record_test.go`, `cmd/fid/integration/p24_judgment_corpus_test.go`
> > - `project-docs/specs/06_DESIGN_PLAN_CLI.md`, `project-docs/docs/phases/phase-24/tasks/task-5-done.md`
> >
> > ## Documentation updated
> >
> > - `project-docs/specs/06_DESIGN_PLAN_CLI.md` documents the expected pair, both Git commands, SourceInput positions and mismatch shape.
> >
> > ## Validation Contract Coverage
> >
> > - VAL-009: PASS for `review.record` - `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at `cmd/fid/integration/p23_review_record_test.go:121` advances HEAD, replays the original key without an Event, and checks a new key sees the changed HEAD.
> > - VAL-011: PASS for both leaves - `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` at `cmd/fid/input_descriptor_test.go:235` checks the live registry/type tree and seeded mismatch; `TestP24ReviewAndFindingDescriptorHelpAndStructure` at line 277 checks both leaves left the pending fence.
> > - VAL-012: PASS for both leaves - `TestP24ReviewAndFindingDescriptorHelpAndStructure` at `cmd/fid/input_descriptor_test.go:277` checks the help template line and both Git command meanings; `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` checks field conformance.
> > - VAL-013: PASS for both leaves - `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:44` runs each standalone `--example` through the built CLI; direct CLI invocations below show the ordered JSON templates.
> > - VAL-014: PASS for both leaves - `TestP24ReviewAndFindingDescriptorHelpAndStructure` at `cmd/fid/input_descriptor_test.go:277` validates and decodes each generated example into the closed Go request type.
> > - VAL-015: PASS for both leaves - `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:44` exercises six conflicting inputs per leaf, and standalone success, with empty non-project workspaces.
> > - VAL-016: PASS for both leaves - `TestP24ReviewAndFindingDescriptorHelpAndStructure` at `cmd/fid/input_descriptor_test.go:277` checks exact top-level and nested paths for unknown, duplicate, missing, null, type, format and enum failures through the CLI, with no created project state.
> > - VAL-017: PASS for `review.record` - `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at `cmd/fid/integration/p23_review_record_test.go:45` rejects top-level and all three nested caller `repository` positions with their exact paths and unchanged Event counts.
> > - VAL-018: PASS for `review.record` - the same test at `cmd/fid/integration/p23_review_record_test.go:77` checks mismatching supplied digests at item, exclusion and evidence positions with both digest values and no Event; its accepted Review at line 104 checks a matching supplied exclusion digest and all three stored resolved digests. All positions use the same `resolveSources` check.
> > - VAL-019: PASS for `review.record` - the same test at `cmd/fid/integration/p23_review_record_test.go:104` checks omitted item and evidence digests against independently hashed file bytes; at line 154 it checks missing-file `source-digest-mismatch` with both digest fields null and no Event.
> > - VAL-022: PASS for `review.record` - the same test at `cmd/fid/integration/p23_review_record_test.go:110` compares stored Review and all three Source repositories with the single fid observation and independently checks HEAD, dirty and tree; `TestP23ReviewRecordRealProcessJourney` at line 177 retains exact readback.
> > - VAL-023: PASS for `review.record` - `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at `cmd/fid/integration/p23_review_record_test.go:135` checks new-key `repository-unavailable` in a non-Git state versus normal success at line 104, with unchanged Event count.
> > - VAL-024: PASS for `review.record` - `TestP23ReviewRecordRealProcessJourney` at `cmd/fid/integration/p23_review_record_test.go:208` checks dirty and HEAD differences with no Events; its first matching Review succeeds at line 252.
> > - VAL-025: PASS - the two mismatch subtests of `TestP23ReviewRecordRealProcessJourney` at `cmd/fid/integration/p23_review_record_test.go:208` and line 231 check expected echo and full observed Repository; `TestReasonCodeDetailSchemaIsExactForEveryCode` at `internal/cli/codes_test.go:38` pins the closed shape.
> > - VAL-035: PASS for the two leaves - `TestP23SharedLeafConformance` at `cmd/fid/integration/p23_conformance_test.go:143` compares accepted, unknown-field and conflict fake/fid outcomes over the derived operation set, including both Review leaves.
> > - VAL-052: PASS - the final ordinary gate reported `cmd/fid` 267.739 s, `internal/workflow` 115.695 s and slowest other ordinary package `internal/issues` 158.407 s; the tagged gate reported 252.256 s. None reached the 480 s ordinary warning, 600 s ordinary limit or 900 s tagged limit.
> > - VAL-060: PASS for both leaves - `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:44` runs each standalone example from an empty non-project directory, checking exit 0, empty stderr and no created entries.
> > - VAL-066: PASS for `review.record` - `TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation` at `cmd/fid/integration/p23_review_record_test.go:135` replays without an observable Git repository, and at line 149 replays after deletion of a coverage file. New-key controls return `repository-unavailable` and `source-digest-mismatch` with no Event.
> >
> > ## Tests run
> >
> > - First TDD red: `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run 'TestP23ReviewRecordRealProcessJourney/dirty_expectation_mismatch_echoes_only_the_expected_pair' ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`); the old decoder rejected `expected_repository` before the repository check.
> > - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -timeout 30m ./internal/project ./internal/review ./internal/cli ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`) on the first retained-fixture attempt; after convergence, the identical command passed with `ok  	github.com/mbell-nzl/fid/cmd/fid	194.986s` (terminated `exit_code: 0`).
> > - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run 'TestP24JudgmentKeyedLeafConflictAndCanonicalRetryCorpus|TestP23ReviewRecordRealProcessJourney' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.435s` (terminated `exit_code: 0`)
> > - `golangci-lint run ./...` - FAIL - `* staticcheck: 1` (terminated `exit_code: 1`) on the first exact-gate attempt; the corrected command passed as recorded below.
> > - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`) on its first exact-gate attempt; the corrected command passed as recorded below.
> > - `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> > - `/tmp/fid review record --example` - PASS - `}` (terminated `exit_code: 0`)
> > - `/tmp/fid finding address --example` - PASS - `}` (terminated `exit_code: 0`)
> >
> > Both direct CLI invocations ran from `/tmp`.
> >
> > ## Exact quality-gate results
> >
> > Working directory: `/home/mark/projects/fid`. The five required commands were run in task-spec order. After causal repairs, only failed commands were repeated; completed passing checks were retained.
> >
> > - `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> > - `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
> > - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	115.695s	coverage: 62.6% of statements` (terminated `exit_code: 0`)
> > - `go tool cover -func=coverage.out` - PASS - `total:											(statements)				60.7%` (terminated `exit_code: 0`)
> > - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	252.256s` (terminated `exit_code: 0`)
> >
> > ## Deviations from spec
> >
> > - No accepted behavior, protected assertion, command or threshold changed. The first lint attempt ended `* staticcheck: 1` (terminated `exit_code: 1`) because a new test declared `var typed any =`; changing it to `typed :=` made the identical lint command return `0 issues.` (terminated `exit_code: 0`).
> > - The first exact tagged gate ended `FAIL` (terminated `exit_code: 1`). Task 5's new verdict descriptor rejected a retained corpus mutation's invalid enum, and two retained mismatch comparisons treated Go typed nil differently from JSON null. The corpus now changes `changes-required` to the accepted `blocked` verdict, preserving its idempotency-conflict assertion; the comparisons normalize the fixture to JSON while retaining full Repository equality. The identical tagged command later returned `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	252.256s` (terminated `exit_code: 0`).
> > - The touched-package race attempt ended `FAIL` (terminated `exit_code: 1`) because retained Review/Finding fixture inputs still used displaced shapes. Converging those inputs produced `ok  	github.com/mbell-nzl/fid/cmd/fid	194.986s` (terminated `exit_code: 0`) on the identical command.
> > - The first report-lint attempt ended `FAIL report-lint found 4 failure(s)` (terminated `exit_code: 4`) because four passing command rows had punctuation after the termination field. Removing those periods preserved the recorded tokens and statuses.
> >
> > ## Concerns, technical debt, or notes for the Tech Lead
> >
> > - Retained assertions changed only where Task 5 displaced their input or observation model. Introducing commits were `d215e40` and `926ef72` (phase-23/task-4 Review journeys), `d1c016d` (phase-23/task-13 failure corpus), `3335149` and `2e7fd2f` (phase-23/task-18 conformance and leaf oracle), `472ed64` (phase-23/task-13 code-schema assertion), and `9bbbdb5` (phase-24/task-2 keyed conflict corpus). Their subjects and no-Event checks remain intact.
> > - No package reached the 80% timing warning.
> >
> > ## Dependency-budget note
> >
> > - No dependency changes.
> >
> > ## Commit hash
> >
> > - `ff6c923dcb706513e71ff7dca76723d68da6ba89`
> >
> > ## Commit subject
> >
> > - `feat(review): check repository expectations and resolve sources`
>
> ## Concerns, technical debt, or notes for the Tech Lead
>
> - Same-task scope gap: Task 5 commit `ff6c923dcb706513e71ff7dca76723d68da6ba89` introduced the lowercase expected head OID format (`git log -1 -- cmd/fid/input_descriptors_review.go`). The shared uppercase OID validator branch was introduced by Task 3 commit `04b6b0992fe20a02df402a183269a9847ddae5ab` (`git log -1 -S 'case "OID":' -- internal/cli/input_validate.go`). The Task 5 adoption exposed this mismatch; this report does not transfer repair ownership. Neither production path is in this dispatch's allowlist. Obtain bounded same-task repair scope and preserve VAL-014 placeholder typing while establishing actual valid OID structural controls and exact malformed-OID rejection.
> - All new test assertions strengthen retained Review/Finding evidence. The original report's introducing-commit provenance remains quoted above. No production, ledger, triage, control, contract or gate file was modified by this child.
> - No package reached its 80% warning: cmd/fid 368.515 s, internal/workflow 100.094 s, slowest other ordinary internal/issues 193.228 s (<480 s warning, <600 s budget); tagged integration 297.291 s (<720 s warning, <900 s budget). These actual timings do not override the failed ordinary gate.
>
> ## Dependency-budget note
>
> - No dependency changes.
>
> ## Commit hash
>
> - `04246029c4488eec46eba76323dee7fbe8260c8d`
>
> ## Commit subject
>
> - `test(review): checkpoint Window B evidence repair`

## Concerns, technical debt, or notes for the Tech Lead

- Fresh independent Window B QA re-review of F4-F7 and the shared format dispatch remains required under TRG-24-005. This Developer issues no QA, TL-Lite or Phase Runner verdict and edits none of their artifacts.
- The old same-task scope gap is resolved by TRG-24-005. Review's lowercase descriptor metadata and the earlier integration fixtures remain unchanged; all source changes are inside matchesInputFormat's authorized dispatch/fallback.
- VAL-052 headroom: cmd/fid 366.509 s, internal/workflow 101.002 s, slowest other ordinary internal/issues 190.218 s are below 480 s warning/600 s budget. Tagged integration 313.544 s is below 720 s warning/900 s budget. No headroom warning.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `226cdc94e6460cc254d1b4e916344b8ea2e2cdb3`

## Commit subject

- `fix(cli): enforce declared Review OID formats`
