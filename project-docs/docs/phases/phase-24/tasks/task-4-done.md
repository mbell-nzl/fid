# Phase 24 - Task 4: QA Window B repair of F1-F3 - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

OUTCOME COMPLETE

## Summary of what was implemented

- Completed the bounded Task 4 QA Window B repair from `review-3.md`, findings F1-F3, under the same-task qa-rework recovery directive. Only the three explicitly allowed files changed.
- F1: captured the exact accepted R5/R6 inventory commands in this recording session, with each command's own final output token, observed termination and accounting for all current matches.
- F2: replaced the expected repository fixture with independent Git/config/absolute-root observations and complete applicable record/Source/Event comparisons across all seven judgment leaves; added a config-identity control.
- F3: strengthened the seven leaves' recursive help and structural evidence to require exact metadata and closed error details, including valid controls and no effects. Existing positive, negative, replay, unavailable, descriptor and no-Event assertions remain.
- Retained the original Task 4 implementation/report provenance and every historical failure/evidence row verbatim in the quoted report below. This repair changes test evidence and report documentation only.

## Files created or modified

- `cmd/fid/integration/p23_judgment_repository_test.go`
- `cmd/fid/input_descriptor_test.go`
- `project-docs/docs/phases/phase-24/tasks/task-4-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-24/tasks/task-4-done.md`: F1 registry receipts, F2/F3 assertion evidence, red/green results, fresh exact gates and immutable historical report.

## Validation Contract Coverage

- VAL-008: PASS for Task 4 Sources — `TestP24JudgmentLeavesReplayAfterHeadAndFileChange` checks supplied versus omitted digest conflicts on `source` or `references`, then unchanged-input replay.
- VAL-009: PASS for all seven leaves — the same test compares byte-identical original results and unchanged Event counts after HEAD and source-file change.
- VAL-011: PASS for seven adopted descriptors — `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` checks live membership, Go field trees and seeded mismatch rejection; `TestP24JudgmentDescriptorSourcePositionsAndExamples` checks all seven have left the pending fence.
- VAL-012: PASS for seven leaves (F3) — `p24JudgmentHelpMatches` at `cmd/fid/input_descriptor_test.go:707` recursively matches each complete help line, including exact type/element type, requirement, format, enum, applicable bounds and meaning; `TestP24JudgmentDescriptorSourcePositionsAndExamples` checks purpose and template line. `TestP24JudgmentEvidenceRejectsIncompleteMetadata` at line 776 rejects individually corrupted metadata/meaning tokens at every depth and accepts complete live help.
- VAL-013: PASS for seven leaves — `TestP24JudgmentDescriptorSourcePositionsAndExamples` compares every standalone example with the ordered, two-space `cli.InputExample` template.
- VAL-014: PASS for seven leaves — the same test applies closed Go type checking and descriptor validation to every generated template.
- VAL-015: PASS for seven leaves — `TestAdoptedExampleProcessExclusivity` checks the six conflicting flag/argument combinations per adopted judgment leaf through the built CLI, with exact error details and no workspace effects; the descriptor/example test checks standalone success.
- VAL-016: PASS for seven leaves (F3) — `TestP24JudgmentStructuralValidationCorpus` at `cmd/fid/input_descriptor_test.go:1166` preserves the original corpus and adds `p24JudgmentRecursiveStructuralDetails` at line 902. It checks exact Path/Value/Allowed for unknown/duplicate, missing/null/type, format/enum and applicable string/array bounds/sorted-unique errors at top-level, nested-object and indexed array-object depths, including element errors. Each row checks the CLI envelope and closed details, operation, invalid-input category, invalid-argument code, exit 2, empty stdout, false retryability, null subject, no Event and unchanged directory. Expanded valid inputs, optional omission and nullable null reach the same CLI family fallback after accepted structural validation. The seeded-details regression rejects wrong values and nonempty but incorrect alternatives.
- VAL-017: PASS for five changed Source leaves — `TestP24JudgmentSourceInputValidation` at `cmd/fid/integration/p23_judgment_repository_test.go:164` rejects caller `repository` members at named Source paths with exit 2 and unchanged Event counts.
- VAL-018: PASS for all five Source positions — `TestP24JudgmentSourceDigestsAndRepositoryParity` at line 117 checks supplied matching digest stored verbatim and mismatch with both digests, exit 4 and no Event; source-position enumeration remains in the descriptor test.
- VAL-019: PASS — the same parity test checks omitted digest becomes the real file SHA-256; `TestP24JudgmentLeavesReplayAfterHeadAndFileChange` checks missing-file controls with both digests null and no Event.
- VAL-020: PASS — `TestP24JudgmentSourceInputValidation` covers absolute/parent paths, escaping symlink, directory and accepted in-root symlink, with exact invalid paths and null observed digest.
- VAL-021: PASS — the same test rejects omitted external digest and accepts a supplied external digest verbatim without a local file.
- VAL-022: PASS for all seven leaves (F2) — `p24JudgmentExpectedRepository` at `cmd/fid/integration/p23_judgment_repository_test.go:26` independently executes Git HEAD/tree/status, reads config project ID and computes SHA-256 of the absolute root. The clean normal fixtures are explicitly checked through Git (excluding accepted `.fid` state), giving false dirty/null dirty digest; base is independently HEAD and capture source is Git. No production `project.ObserveRecordedRepository` is used. `TestP24JudgmentRepositoryOracleUsesConfig` at line 56 proves a decoy cached adapter ID cannot become the expected config ID. `TestP24JudgmentLeavesRepositoryUnavailable` at line 294 compares the full eight-field Event repository for every leaf and recursively every returned record/Source repository, pinning one Source repository for each of five reference-bearing records and zero for the source-free disposal/link fixtures. `TestP24JudgmentSourceDigestsAndRepositoryParity` independently compares the five stored Source repositories and Events as well.
- VAL-023: PASS — `TestP24JudgmentLeavesRepositoryUnavailable` retains all seven normal/non-Git/unmerged rows with exact reason/detail, precondition, retryable, exit 4 and no Event. `TestRepositoryUnavailableCauseFixtures`, `TestRepositoryCaptureChangeReason`, and `TestRepositoryUnmergedAndSizeClassification` in `internal/project/repository_reason_test.go` cover unborn HEAD, tracked flag/FIFO, capture change, unmerged parser and size limit in the fresh ordinary gate.
- VAL-035: PASS for the Task 4 fake/fid intersection — `TestP23SharedLeafConformance` in `cmd/fid/integration/p23_conformance_test.go` compares accepted Source and removed-member rejection cases in the fresh tagged gate.
- VAL-052: PASS — fresh final ordinary/tagged package timings and unchanged budgets are recorded below; no package reaches its 80% warning threshold.
- VAL-060: PASS for seven leaves — `TestP24JudgmentDescriptorSourcePositionsAndExamples` runs each standalone example in an empty non-project directory with exit 0, empty stderr and no created entry.
- VAL-066: PASS for seven leaves — `TestP24JudgmentLeavesRepositoryUnavailable` retains non-Git/unmerged replay with otherwise viable new-key rejection/no-Event controls; `TestP24JudgmentLeavesReplayAfterHeadAndFileChange` retains deleted-source replay, both-null new-key mismatch/no Event and restored-file success.

F1 is documentation-only inventory evidence with no separate VAL ID. All other Task 4 evidence above comes from named assertions exercised in this session's fresh final gates; the historical report is preserved separately below.

## Tests run

Working directory: `/home/mark/projects/fid`, except the direct CLI invocation explicitly marked below. Focused commands supplement the exact required gates.

- F2 first focused red: `go test -tags integration -run '^TestP24JudgmentRepositoryOracleUsesConfig$' -count=1 ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`). The previous production-observer fixture accepted the decoy cached project ID instead of config identity.
- F2 intermediate fixture red: the byte-identical command - FAIL - `FAIL` (terminated `exit_code: 1`); the independent status fixture initially counted untracked `.fid` metadata. Applying the accepted `.fid` exclusion repaired that expected-value fixture.
- F2 byte-identical green: `go test -tags integration -run '^TestP24JudgmentRepositoryOracleUsesConfig$' -count=1 ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.245s` (terminated `exit_code: 0`).
- F3 first focused red: `go test -run '^TestP24JudgmentEvidenceRejectsIncompleteMetadata$' -count=1 ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`). The previous help/details checks accepted wrong metadata and nonempty but incorrect alternatives.
- F3 byte-identical green: `go test -run '^TestP24JudgmentEvidenceRejectsIncompleteMetadata$' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.156s` (terminated `exit_code: 0`).
- Expanded corpus: `go test -run '^(TestP24JudgmentDescriptorSourcePositionsAndExamples|TestP24JudgmentStructuralValidationCorpus|TestP24JudgmentEvidenceRejectsIncompleteMetadata)$' -count=1 ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`). The expected oversize-text alternative initially ignored existing enum-before-format precedence. Repairing only that expected fixture made the byte-identical command PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.535s` (terminated `exit_code: 0`).
- Expanded array corpus: the byte-identical expanded-corpus command - FAIL - `FAIL` (terminated `exit_code: 1`), because the test generator indexed an empty baseline `scope.paths` array. Supplying a valid test element for the maximum-length fixture made the byte-identical command PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.656s` (terminated `exit_code: 0`).
- `go test -tags integration -run '^(TestP24JudgmentSourceDigestsAndRepositoryParity|TestP24JudgmentLeavesRepositoryUnavailable)$' -count=1 ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.164s` (terminated `exit_code: 0`).
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	242.901s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24Judgment' -count=1 ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	9.219s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^(TestP24JudgmentDescriptorSourcePositionsAndExamples|TestP24JudgmentStructuralValidationCorpus|TestP24JudgmentEvidenceRejectsIncompleteMetadata)$' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	4.150s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid-phase24-task4-repair ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid-phase24-task4-repair decision record-reference --help` - PASS - `      --project string   fid project root` (terminated `exit_code: 0`)
 
The direct `decision record-reference --help` entrypoint ran from `/tmp`. Its emitted help contains exact nested SourceInput and authority types, requirements, digest/ID formats, label/array limits and the standalone example line.

### First required-gate attempt

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - FAIL - `* staticcheck: 2` (terminated `exit_code: 1`)

The attempt stopped at the lint failure; no later required command was recorded as passing for that attempt. After the scoped switch refactor and focused race check, the full exact sequence below ran in order.

### F1: fresh owned registry receipts

- `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid` - PASS - `cmd/fid/gate.go:272:func decodeClosedGateInput(stdin io.Reader, path string, target any) error {` (terminated `exit_code: 0`)
- `rg -n "func (same[A-Za-z]*Repository|validateRepositoryBinding|validEnvelopeRepository|validEnvelopeSource|validEvidenceRepository|validateRunTerminalRepository|verifyDecisionSource|importedGateDestination)" internal -g '!*_test.go'` - PASS - `internal/evidence/registration.go:362:func validateRepositoryBinding(request RegisterRequest, projectID string) error {` (terminated `exit_code: 0`)

Both inventory searches legitimately return matches with observed exit 0; neither is an absence gate. Their current complete outputs are:

R5, six matches:

```text
cmd/fid/review.go:135:func decodeClosedReviewInput(stdin io.Reader, path string, target any) error {
cmd/fid/run.go:226:func decodeClosedRunInput(stdin io.Reader, path string, target any) error {
cmd/fid/activation.go:328:func decodeClosedActivationInput(stdin io.Reader, path string, target any) error {
cmd/fid/envelope.go:178:func decodeClosedEnvelopeInput(stdin io.Reader, path string, target any) error {
cmd/fid/concern.go:147:func decodeClosedJudgmentInput(cmd *cobra.Command, path string, target any) error {
cmd/fid/gate.go:272:func decodeClosedGateInput(stdin io.Reader, path string, target any) error {
```

The six remaining wrappers are Review, Run, Activation, Envelope, Judgment and Gate. Task 4's Judgment wrapper dispatches the seven adopted descriptors; its retained function name is expected. The Claim decoder is already removed by Task 6. The remaining named wrappers persist until their subsequent adoption/whole-tree decoder retirement (Task 9); this receipt does not claim zero matches.

R6, nine matches:

```text
internal/evidence/evidence.go:668:func validEvidenceRepository(repository Repository) bool {
internal/evidence/evidence.go:672:func sameEvidenceRepository(left, right Repository) bool {
internal/evidence/execution.go:346:func importedGateDestination(conn *sql.DB, workspaceRoot string, registration Registration, request GateRunRequest) bool {
internal/evidence/execution.go:431:func sameGateRepository(left, right Repository) bool {
internal/workflow/run_close.go:290:func validateRunTerminalRepository(workspaceRoot string, asserted Repository) error {
internal/workflow/run_read.go:149:func sameRepository(left, right Repository) bool {
internal/workflow/envelope.go:755:func validEnvelopeRepository(repository Repository) bool {
internal/workflow/envelope.go:759:func validEnvelopeSource(source Source) bool {
internal/evidence/registration.go:362:func validateRepositoryBinding(request RegisterRequest, projectID string) error {
```

The nine remaining matches account for Run terminal binding and Run repository comparison; Evidence repository validation/comparison; Gate import-destination matching and Gate repository comparison; Registration binding validation; and Envelope repository/Source validation. Judgment's former source verification/comparison has been replaced by the accepted shared observer/resolver, and Review's former comparator is also absent after its adoption. The remaining Run/Gate/Evidence convergence and retained binding/output validators are accounted matches, rather than evidence of absence. The shared worktree primitive is outside this historical function-name registry.


## Exact quality-gate results

Working directory: `/home/mark/projects/fid`. The five commands below ran byte-identically in the required order on the final test tree. The initial lint failure remains above and in the deviation note below. Focused checks did not replace any required gate.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	102.588s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					60.9%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	276.204s` (terminated `exit_code: 0`)

Fresh reported package wall times: `cmd/fid` 374.499 s; `internal/workflow` 102.588 s; `internal/issues` 198.298 s (slowest other ordinary package); `cmd/fid/integration` 276.204 s. Ordinary packages stay below the 600 s budget and 480 s warning; tagged integration stays below the 900 s budget and 720 s warning. None of these measurements is cached. Overall statement coverage is 60.9%; the accepted coverage policy remains unchanged.

## Deviations from spec

- No accepted Observation, SourceInput, judgment, public-interface, command, assertion semantics, coverage policy or budget changed. This recovery is limited to review-3 F1-F3 and its exact allowed files.
- Mandatory first-failure note: the initial `golangci-lint run ./...` returned `* staticcheck: 2` (terminated `exit_code: 1`). Two new test-helper `if/else` chains triggered QF1003. Rewriting those chains as tagged switches, with identical expected values and assertions, repaired the in-scope cause. The byte-identical lint command later returned `0 issues.` (terminated `exit_code: 0`), and all five final required commands passed. The first failure remains immutable failed evidence.
- Focused fixture failures remain above: the config decoy exposed the production oracle, metadata/details corruption exposed weak evidence, accepted `.fid` exclusion repaired the independent clean-status fixture, enum precedence corrected one expected alternative, and an explicit valid array element repaired the test generator's empty-array panic. No production behavior or accepted validation semantics was altered to pass these checks.

### Preserved original completion report

The following is the complete original report read before this repair, including all failures, results and provenance. Its implementation commit was `7c69c2bbb6b8b7f8688bb6d28fc72f8fb5fff61c`; its committed report endpoint was `7ac55a292ced5e0eabd6933394864da2fae01340`. Removing one quotation prefix from each physical line reconstructs that report verbatim. These are historical results; the fresh repair evidence appears above.

> # Phase 24 - Task 4: Repository observation and SourceInput; judgment reference leaves - Completion Report
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
> - Added a shared recorded-repository observation and typed unavailable reasons while preserving `errors.Is(ErrUnsafeRepository)`. Added `SourceInput` resolution for checked or observed project-file SHA-256 digests, external references, and safe relative paths.
> - Converged five judgment reference mutations to resolve and store Sources with fid's observation. `requirement link` and `concern dispose` now observe and have descriptors; all seven replay before observation. Human Outcome references are checked. Added the repository-unavailable code and nullable source-digest-mismatch details.
> - Added descriptor help, templates, exclusivity and structural error coverage for the seven leaves. Moved requirement and concern input decoding before workspace lookup so exact structural paths reach their existing family renderers. Converged retained Phase 23 fixtures and fake/fid oracles to the accepted SourceInput shape.
> 
> ## Files created or modified
> 
> - `internal/project/source.go`, `internal/project/source_test.go`, `internal/project/repository.go`, `internal/project/repository_reason_test.go`
> - `internal/review/decision.go`, `internal/review/decision_test.go`, `internal/review/requirement.go`, `internal/review/human_outcome.go`, `internal/review/concern.go`, `internal/review/concern_task8_test.go`, `internal/review/deviation.go`, `internal/review/review.go`
> - `internal/cli/codes.go`, `internal/cli/codes_test.go`
> - `cmd/fid/decision.go`, `cmd/fid/requirement.go`, `cmd/fid/human_outcome.go`, `cmd/fid/concern.go`, `cmd/fid/deviation.go`, `cmd/fid/input_descriptors_judgment.go`, `cmd/fid/input_descriptor_test.go`
> - `cmd/fid/integration/p23_judgment_repository_test.go`, `cmd/fid/integration/p23_human_outcome_test.go`, `cmd/fid/integration/p23_conformance_test.go`, `cmd/fid/integration/p23_concern_leaf_oracle_test.go`, `cmd/fid/integration/p23_decision_reference_test.go`, `cmd/fid/integration/p23_failure_contract_test.go`, `cmd/fid/integration/p23_judgment_leaf_oracle_test.go`, `cmd/fid/integration/p23_requirement_leaf_oracle_test.go`
> - `project-docs/specs/06_DESIGN_PLAN_CLI.md`, `project-docs/specs/GLOSSARY.md`, `project-docs/docs/phases/phase-24/tasks/task-4-done.md`
> 
> ## Documentation updated
> 
> - `project-docs/specs/06_DESIGN_PLAN_CLI.md` documents Task 4 SourceInput and judgment input examples.
> - `project-docs/specs/GLOSSARY.md` defines Expected repository and SourceInput.
> 
> ## Validation Contract Coverage
> 
> - VAL-008: PASS for Task 4 Sources - `TestP24JudgmentLeavesReplayAfterHeadAndFileChange` at `cmd/fid/integration/p23_judgment_repository_test.go:352` checks supplied versus omitted digest under one key conflicts on `source` or `references`; unchanged input replays after HEAD and file change.
> - VAL-009: PASS for all seven leaves - the same test replays byte-identical results with no Event after HEAD and source file change.
> - VAL-011: PASS for seven adopted descriptors - `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` at `cmd/fid/input_descriptor_test.go:234` checks registry membership and Go field trees with a seeded mismatch; `TestP24JudgmentDescriptorSourcePositionsAndExamples` at line 277 checks all seven are out of the pending fence.
> - VAL-012: PASS for seven leaves - `TestP24JudgmentDescriptorSourcePositionsAndExamples` checks every descriptor field path and meaning in help plus the `fid <leaf> --example` line.
> - VAL-013: PASS for seven leaves - the same test compares each `--example` output with the ordered, two-space `cli.InputExample` template.
> - VAL-014: PASS for seven leaves - the same test runs `cli.CheckInputType` and `cli.ValidateInput` for each generated template.
> - VAL-015: PASS for seven leaves - `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:44` exercises six conflicting flag/argument combinations per Task 4 leaf through the built CLI; `TestP24JudgmentDescriptorSourcePositionsAndExamples` checks standalone success.
> - VAL-016: PASS for seven leaves - `TestP24JudgmentStructuralValidationCorpus` at `cmd/fid/input_descriptor_test.go:361` checks top-level unknown, duplicate, missing, null, type, format, applicable enum and nested Source variants through the validator and CLI with exact `input.<path>`, nonempty allowed values, exit 2 and no Event.
> - VAL-017: PASS for five Source leaves - `TestP24JudgmentSourceInputValidation` at `cmd/fid/integration/p23_judgment_repository_test.go:126` supplies `repository` inside each Source, checks the named path, exit 2 and unchanged Event count.
> - VAL-018: PASS for all five Source positions - `TestP24JudgmentSourceDigestsAndRepositoryParity` at `cmd/fid/integration/p23_judgment_repository_test.go:75` checks matching digest stored verbatim and mismatching digest with both values, exit 4 and no Event; the descriptor source-position completeness check is at `cmd/fid/input_descriptor_test.go:277`.
> - VAL-019: PASS - the same integration test checks omitted digest stored as the file's SHA-256; `TestP24JudgmentLeavesReplayAfterHeadAndFileChange` checks missing-file controls with both digests null and no Event.
> - VAL-020: PASS - `TestP24JudgmentSourceInputValidation` covers absolute and parent paths, escaping symlink, directory and in-root symlink, including named invalid paths and null observed digest.
> - VAL-021: PASS - the same test checks omitted external digest is a named invalid argument and a supplied external digest is stored verbatim without a local file.
> - VAL-022: PASS for all seven leaves - `TestP24JudgmentLeavesRepositoryUnavailable` at `cmd/fid/integration/p23_judgment_repository_test.go:256` compares each normal Event repository with an independent observation; `TestP24JudgmentSourceDigestsAndRepositoryParity` compares all five stored Source repositories and Events.
> - VAL-023: PASS - `TestP24JudgmentLeavesRepositoryUnavailable` covers normal, non-Git and unmerged states for each of seven leaves with exact `{reason,detail}`, precondition, retryable, exit 4 and no Event. `TestRepositoryUnavailableCauseFixtures`, `TestRepositoryCaptureChangeReason`, and `TestRepositoryUnmergedAndSizeClassification` at `internal/project/repository_reason_test.go:24,56,69` cover unborn HEAD, tracked flag, tracked FIFO, capture change, unmerged parser and size limit.
> - VAL-035: PASS for Task 4 fake/fid intersection - `TestP23SharedLeafConformance` at `cmd/fid/integration/p23_conformance_test.go:143` compares the converged Source and rejection cases after the fake model was updated.
> - VAL-052: PASS - final exact ordinary gate reports `cmd/fid` 223.171 s and `internal/export` 104.915 s; the first uncached ordinary attempt reported `internal/workflow` 122.340 s and `internal/issues` 181.213 s, the slowest other ordinary package. Final tagged integration was 284.152 s. No ordinary package reached the 480 s warning or 600 s limit; tagged integration stayed below 900 s.
> - VAL-060: PASS for seven leaves - `TestP24JudgmentDescriptorSourcePositionsAndExamples` runs each standalone `--example` in an empty non-project directory and verifies exit 0, empty stderr and no created entry.
> - VAL-066: PASS for seven leaves - `TestP24JudgmentLeavesRepositoryUnavailable` replays in non-Git and unmerged states while new-key controls fail without Events; `TestP24JudgmentLeavesReplayAfterHeadAndFileChange` replays deleted-source inputs while new-key controls return both-null source-digest-mismatch and succeed when the file is restored.
> 
> ## Tests run
> 
> - First TDD red: `TestP24HumanOutcomeRejectsReferenceDigestMismatch` accepted a wrong reference digest and appended an Event before the implementation; focused result `FAIL` (terminated `exit_code: 1`). The final tagged gate passes this test.
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`); immutable first required-gate attempt before fake/fid convergence.
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/project ./internal/review ./internal/cli ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	202.774s` (terminated `exit_code: 0`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23SharedLeafConformance$' -count=1 ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	13.271s` (terminated `exit_code: 0`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^(TestP24JudgmentStructuralValidationCorpus|TestAdoptedExampleProcessExclusivity)$' -count=1 ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`); first structural run showed generic paths for requirement and concern leaves. After decoding order repair, the identical focused command passed: `ok  	github.com/mbell-nzl/fid/cmd/fid	3.786s` (terminated `exit_code: 0`).
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24Judgment' -count=1 ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.640s` (terminated `exit_code: 0`)
> - `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> - `/tmp/fid decision record-reference --example` - PASS - `}` (terminated `exit_code: 0`)
> 
> The CLI invocation ran from `/tmp`; its SourceInput template has no caller repository member.
> 
> ## Exact quality-gate results
> 
> Working directory: `/home/mark/projects/fid`. The first failed tagged attempt is retained under Tests run and Deviations from spec. The final five commands ran in the task-spec order on the final tree.
> 
> - `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> - `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.6% of statements` (terminated `exit_code: 0`)
> - `go tool cover -func=coverage.out` - PASS - `total:											(statements)				60.6%` (terminated `exit_code: 0`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	284.152s` (terminated `exit_code: 0`)
> 
> ## Deviations from spec
> 
> - No accepted behavior, assertion, command, or threshold changed. The first required tagged gate failed with last token `FAIL` (terminated `exit_code: 1`) because retained Phase 23 fake/fid conformance still modeled caller-provided Source repositories and represented a null dirty digest as a typed nil. Converging the fake to independently derive the stored Source and normalize JSON null made the byte-identical focused conformance command return `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	13.271s` (terminated `exit_code: 0`) and the later exact tagged gate return `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	284.152s` (terminated `exit_code: 0`). The failed attempt remains failed evidence.
> - The new structural corpus first returned `FAIL` (terminated `exit_code: 1`) because requirement and concern looked up the workspace before decoding closed input. Moving decoding ahead of workspace lookup under Task 4 scope made the identical focused command return `ok  	github.com/mbell-nzl/fid/cmd/fid	3.786s` (terminated `exit_code: 0`).
> 
> ## Concerns, technical debt, or notes for the Tech Lead
> 
> - Retained assertions were converged only where Task 4 displaced their input model, with their subjects and strength preserved. Introducing commits: `828eb18` (phase-23/task-9, decision reference and CLI code tests), `3a0aed` (phase-23/task-8, concern test), `d1c016` (phase-23/task-13, failure corpus), `333514` and `2e7fd2` (phase-23/task-18, fake conformance and leaf oracles), and `b71ed1` (phase-23/task-14, Human leaf pairs). Provenance was checked against Git history, the task pack and the accepted SourceInput contract.
> - No package reached the 80% timing warning. The first uncached ordinary run supplies `internal/workflow` and `internal/issues` times where later exact runs printed `(cached)`.
> 
> ## Dependency-budget note
> 
> - No dependency changes.
> 
> ## Commit hash
> 
> - `7c69c2bbb6b8b7f8688bb6d28fc72f8fb5fff61c`
> 
> ## Commit subject
> 
> - `feat(judgment): observe repositories and resolve SourceInput`
> 

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
