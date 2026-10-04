# Phase 24 - Task 4: Repository observation and SourceInput; judgment reference leaves - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

OUTCOME COMPLETE

## Summary of what was implemented

- Added a shared recorded-repository observation and typed unavailable reasons while preserving `errors.Is(ErrUnsafeRepository)`. Added `SourceInput` resolution for checked or observed project-file SHA-256 digests, external references, and safe relative paths.
- Converged five judgment reference mutations to resolve and store Sources with fid's observation. `requirement link` and `concern dispose` now observe and have descriptors; all seven replay before observation. Human Outcome references are checked. Added the repository-unavailable code and nullable source-digest-mismatch details.
- Added descriptor help, templates, exclusivity and structural error coverage for the seven leaves. Moved requirement and concern input decoding before workspace lookup so exact structural paths reach their existing family renderers. Converged retained Phase 23 fixtures and fake/fid oracles to the accepted SourceInput shape.

## Files created or modified

- `internal/project/source.go`, `internal/project/source_test.go`, `internal/project/repository.go`, `internal/project/repository_reason_test.go`
- `internal/review/decision.go`, `internal/review/decision_test.go`, `internal/review/requirement.go`, `internal/review/human_outcome.go`, `internal/review/concern.go`, `internal/review/concern_task8_test.go`, `internal/review/deviation.go`, `internal/review/review.go`
- `internal/cli/codes.go`, `internal/cli/codes_test.go`
- `cmd/fid/decision.go`, `cmd/fid/requirement.go`, `cmd/fid/human_outcome.go`, `cmd/fid/concern.go`, `cmd/fid/deviation.go`, `cmd/fid/input_descriptors_judgment.go`, `cmd/fid/input_descriptor_test.go`
- `cmd/fid/integration/p23_judgment_repository_test.go`, `cmd/fid/integration/p23_human_outcome_test.go`, `cmd/fid/integration/p23_conformance_test.go`, `cmd/fid/integration/p23_concern_leaf_oracle_test.go`, `cmd/fid/integration/p23_decision_reference_test.go`, `cmd/fid/integration/p23_failure_contract_test.go`, `cmd/fid/integration/p23_judgment_leaf_oracle_test.go`, `cmd/fid/integration/p23_requirement_leaf_oracle_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`, `project-docs/specs/GLOSSARY.md`, `project-docs/docs/phases/phase-24/tasks/task-4-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` documents Task 4 SourceInput and judgment input examples.
- `project-docs/specs/GLOSSARY.md` defines Expected repository and SourceInput.

## Validation Contract Coverage

- VAL-008: PASS for Task 4 Sources - `TestP24JudgmentLeavesReplayAfterHeadAndFileChange` at `cmd/fid/integration/p23_judgment_repository_test.go:352` checks supplied versus omitted digest under one key conflicts on `source` or `references`; unchanged input replays after HEAD and file change.
- VAL-009: PASS for all seven leaves - the same test replays byte-identical results with no Event after HEAD and source file change.
- VAL-011: PASS for seven adopted descriptors - `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` at `cmd/fid/input_descriptor_test.go:234` checks registry membership and Go field trees with a seeded mismatch; `TestP24JudgmentDescriptorSourcePositionsAndExamples` at line 277 checks all seven are out of the pending fence.
- VAL-012: PASS for seven leaves - `TestP24JudgmentDescriptorSourcePositionsAndExamples` checks every descriptor field path and meaning in help plus the `fid <leaf> --example` line.
- VAL-013: PASS for seven leaves - the same test compares each `--example` output with the ordered, two-space `cli.InputExample` template.
- VAL-014: PASS for seven leaves - the same test runs `cli.CheckInputType` and `cli.ValidateInput` for each generated template.
- VAL-015: PASS for seven leaves - `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:44` exercises six conflicting flag/argument combinations per Task 4 leaf through the built CLI; `TestP24JudgmentDescriptorSourcePositionsAndExamples` checks standalone success.
- VAL-016: PASS for seven leaves - `TestP24JudgmentStructuralValidationCorpus` at `cmd/fid/input_descriptor_test.go:361` checks top-level unknown, duplicate, missing, null, type, format, applicable enum and nested Source variants through the validator and CLI with exact `input.<path>`, nonempty allowed values, exit 2 and no Event.
- VAL-017: PASS for five Source leaves - `TestP24JudgmentSourceInputValidation` at `cmd/fid/integration/p23_judgment_repository_test.go:126` supplies `repository` inside each Source, checks the named path, exit 2 and unchanged Event count.
- VAL-018: PASS for all five Source positions - `TestP24JudgmentSourceDigestsAndRepositoryParity` at `cmd/fid/integration/p23_judgment_repository_test.go:75` checks matching digest stored verbatim and mismatching digest with both values, exit 4 and no Event; the descriptor source-position completeness check is at `cmd/fid/input_descriptor_test.go:277`.
- VAL-019: PASS - the same integration test checks omitted digest stored as the file's SHA-256; `TestP24JudgmentLeavesReplayAfterHeadAndFileChange` checks missing-file controls with both digests null and no Event.
- VAL-020: PASS - `TestP24JudgmentSourceInputValidation` covers absolute and parent paths, escaping symlink, directory and in-root symlink, including named invalid paths and null observed digest.
- VAL-021: PASS - the same test checks omitted external digest is a named invalid argument and a supplied external digest is stored verbatim without a local file.
- VAL-022: PASS for all seven leaves - `TestP24JudgmentLeavesRepositoryUnavailable` at `cmd/fid/integration/p23_judgment_repository_test.go:256` compares each normal Event repository with an independent observation; `TestP24JudgmentSourceDigestsAndRepositoryParity` compares all five stored Source repositories and Events.
- VAL-023: PASS - `TestP24JudgmentLeavesRepositoryUnavailable` covers normal, non-Git and unmerged states for each of seven leaves with exact `{reason,detail}`, precondition, retryable, exit 4 and no Event. `TestRepositoryUnavailableCauseFixtures`, `TestRepositoryCaptureChangeReason`, and `TestRepositoryUnmergedAndSizeClassification` at `internal/project/repository_reason_test.go:24,56,69` cover unborn HEAD, tracked flag, tracked FIFO, capture change, unmerged parser and size limit.
- VAL-035: PASS for Task 4 fake/fid intersection - `TestP23SharedLeafConformance` at `cmd/fid/integration/p23_conformance_test.go:143` compares the converged Source and rejection cases after the fake model was updated.
- VAL-052: PASS - final exact ordinary gate reports `cmd/fid` 223.171 s and `internal/export` 104.915 s; the first uncached ordinary attempt reported `internal/workflow` 122.340 s and `internal/issues` 181.213 s, the slowest other ordinary package. Final tagged integration was 284.152 s. No ordinary package reached the 480 s warning or 600 s limit; tagged integration stayed below 900 s.
- VAL-060: PASS for seven leaves - `TestP24JudgmentDescriptorSourcePositionsAndExamples` runs each standalone `--example` in an empty non-project directory and verifies exit 0, empty stderr and no created entry.
- VAL-066: PASS for seven leaves - `TestP24JudgmentLeavesRepositoryUnavailable` replays in non-Git and unmerged states while new-key controls fail without Events; `TestP24JudgmentLeavesReplayAfterHeadAndFileChange` replays deleted-source inputs while new-key controls return both-null source-digest-mismatch and succeed when the file is restored.

## Tests run

- First TDD red: `TestP24HumanOutcomeRejectsReferenceDigestMismatch` accepted a wrong reference digest and appended an Event before the implementation; focused result `FAIL` (terminated `exit_code: 1`). The final tagged gate passes this test.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`); immutable first required-gate attempt before fake/fid convergence.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/project ./internal/review ./internal/cli ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	202.774s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23SharedLeafConformance$' -count=1 ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	13.271s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^(TestP24JudgmentStructuralValidationCorpus|TestAdoptedExampleProcessExclusivity)$' -count=1 ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`); first structural run showed generic paths for requirement and concern leaves. After decoding order repair, the identical focused command passed: `ok  	github.com/mbell-nzl/fid/cmd/fid	3.786s` (terminated `exit_code: 0`).
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24Judgment' -count=1 ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.640s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid decision record-reference --example` - PASS - `}` (terminated `exit_code: 0`)

The CLI invocation ran from `/tmp`; its SourceInput template has no caller repository member.

## Exact quality-gate results

Working directory: `/home/mark/projects/fid`. The first failed tagged attempt is retained under Tests run and Deviations from spec. The final five commands ran in the task-spec order on the final tree.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)				60.6%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	284.152s` (terminated `exit_code: 0`)

## Deviations from spec

- No accepted behavior, assertion, command, or threshold changed. The first required tagged gate failed with last token `FAIL` (terminated `exit_code: 1`) because retained Phase 23 fake/fid conformance still modeled caller-provided Source repositories and represented a null dirty digest as a typed nil. Converging the fake to independently derive the stored Source and normalize JSON null made the byte-identical focused conformance command return `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	13.271s` (terminated `exit_code: 0`) and the later exact tagged gate return `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	284.152s` (terminated `exit_code: 0`). The failed attempt remains failed evidence.
- The new structural corpus first returned `FAIL` (terminated `exit_code: 1`) because requirement and concern looked up the workspace before decoding closed input. Moving decoding ahead of workspace lookup under Task 4 scope made the identical focused command return `ok  	github.com/mbell-nzl/fid/cmd/fid	3.786s` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- Retained assertions were converged only where Task 4 displaced their input model, with their subjects and strength preserved. Introducing commits: `828eb18` (phase-23/task-9, decision reference and CLI code tests), `3a0aed` (phase-23/task-8, concern test), `d1c016` (phase-23/task-13, failure corpus), `333514` and `2e7fd2` (phase-23/task-18, fake conformance and leaf oracles), and `b71ed1` (phase-23/task-14, Human leaf pairs). Provenance was checked against Git history, the task pack and the accepted SourceInput contract.
- No package reached the 80% timing warning. The first uncached ordinary run supplies `internal/workflow` and `internal/issues` times where later exact runs printed `(cached)`.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `7c69c2bbb6b8b7f8688bb6d28fc72f8fb5fff61c`

## Commit subject

- `feat(judgment): observe repositories and resolve SourceInput`
