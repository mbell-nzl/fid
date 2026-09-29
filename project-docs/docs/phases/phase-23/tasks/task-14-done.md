# Phase 23 - Task 14: Human review and judgment projections - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- OUTCOME COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Replaced resume's Findings placeholder with durable open and addressed Findings; projected Reviews, Concerns, Deviations, Decision References, Human Outcomes, protected Decisions, full per-Obligation evaluations, and guarded next actions. A section exceeding 100 items returns `result-limit-exceeded`.
- Added purpose-labeled Human renderers for the available Phase 23 Review and judgment leaves, including explicit Builder-assertion wording for addressed Findings, non-proof wording for judgment facts, source labels/paths/digests instead of project prose, and Stop/Attention qualifiers.
- Added a tagged Human fixture generator and verbatim guide for mixed Finding dispositions, a protected Decision with an activated Human Outcome actor, and review-cycle exhaustion. The Human judgments remain phase-end work.
- Implemented the accepted generic Concern list read with subject, status, severity, and limit filters. The CLI now emits `{items:Concern[]}` in the canonical JSON envelope and the existing Human judgment renderer. It returns typed `result-limit-exceeded` instead of truncating. The Issue-scoped V1 read remains available until its separate removal, while the generic read excludes V1 rows.
- Preserved the original focused and tagged failures in this report, then confirmed the retained Human/JSON assertion and the five required gates pass after the TRG-23-009 repair.

## Files created or modified

- `internal/workflow/resume.go`, `internal/workflow/resume_actions.go`
- `cmd/fid/resume.go`, `cmd/fid/resume_human.go`, `cmd/fid/p23_human.go`, `cmd/fid/stop_human.go`
- `cmd/fid/review.go`, `cmd/fid/finding.go`, `cmd/fid/concern.go`, `cmd/fid/deviation.go`, `cmd/fid/decision.go`, `cmd/fid/human_outcome.go`, `cmd/fid/requirement.go`
- `cmd/fid/p23_human_test.go`, `cmd/fid/task18_completion_resume_test.go`, `cmd/fid/contract_surface_test.go`
- `cmd/fid/integration/p23_human_fixtures_test.go`, `cmd/fid/integration/p23_resume_limit_test.go`, `cmd/fid/integration/p23_concern_list_human_test.go`
- `cmd/fid/integration/p22_resume_human_test.go`, `p23_review_record_test.go`, `p23_finding_disposition_test.go`, `p23_concern_deviation_test.go`, `p23_decision_reference_test.go`, `p23_human_outcome_test.go`, `p23_review_evidence_budget_test.go`
- `cmd/fid/concern_test.go`, `internal/review/concern_list.go`; the existing `cmd/fid/concern.go` and `cmd/fid/integration/p23_concern_list_human_test.go` now include the bounded recovery.
- `project-docs/docs/phases/phase-23/qa/human-fixtures.md`, this report

## Documentation updated

- `project-docs/docs/phases/phase-23/qa/human-fixtures.md` records the retained binary SHA-256, absolute workspace and invocation arguments, verbatim Human stdout, JSON counterparts, and VAL-070/071/072 action steps.
- This report records the TRG-23-009 repair, the immutable failed attempts, and the passing final-tree evidence.

## Validation Contract Coverage

- VAL-058: PASS for task-level deterministic coverage - `TestPhase23HumanLeafTextMatchesRecordedFacts` (`cmd/fid/p23_human_test.go:11`) checks the renderer shapes and prose exclusion; `TestP23ConcernListHumanProjection` (`cmd/fid/integration/p23_concern_list_human_test.go:10`) now exercises the real `concern list --subject` entrypoint in both canonical JSON and Human modes, checking the same ID and excluding private prose. The generated Review, Finding, resume, Decision, Stop, and Attention views remain in `qa/human-fixtures.md`; the window 13-15 reviewer retains independent every-leaf fixture judgment.
- VAL-059: PASS for the implemented resume projection - `TestResumeReportsRealFindingsAndJudgmentSections` (`cmd/fid/task18_completion_resume_test.go:201`) checks the new judgment arrays; `TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess` (`cmd/fid/integration/p23_human_outcome_test.go:39`) checks all four ordered Obligation evaluations; `TestP23HumanFixtures` produces a populated mixed-Finding resume inspection recorded in `qa/human-fixtures.md`; `TestP23ResumeRejectsJudgmentSectionOverflow` (`cmd/fid/integration/p23_resume_limit_test.go:11`) records 101 Concerns and checks typed `result-limit-exceeded` at limit 100.
- VAL-060: PASS for task-level deterministic coverage - `TestPhase23HumanLeafTextMatchesRecordedFacts` rejects private prose sentinels and checks judgment non-proof language; `TestP23ConcernListHumanProjection` rejects private Concern prose in the real CLI Human output; `TestP23HumanFixtures` compares verbatim retained Human output after generator exit. Independent QA review remains pending.
- VAL-070: FIXTURE READY, HUMAN PENDING - `TestP23HumanFixtures` (`cmd/fid/integration/p23_human_fixtures_test.go:46`) seeded resolved, addressed, and still-blocking Findings; `qa/human-fixtures.md` gives absolute Review/Finding/resume invocations. The Human has not judged comprehension.
- VAL-071: FIXTURE READY, HUMAN PENDING - the same generator seeded a protected Decision and an activated `human` actor. A real `fid` invocation returned `protected-decision-required` with exit 5, Human Outcome recording exited 0, and retry exited 0. Exact action arrays are in the fixture guide. The Human has not judged clarity.
- VAL-072: FIXTURE READY, HUMAN PENDING - the same generator seeded exhausted `review-cycles`, a Stop, and Attention; a real `fid` denial returned `review-rework-exhausted` with exit 4. Exact inspection/action arrays are in the fixture guide. The Human has not judged clarity.
- VAL-073: ADVANCED, HUMAN PENDING - purpose lines and labeled sections are asserted by `TestPhase23HumanLeafTextMatchesRecordedFacts` and recorded verbatim in the fixture guide; Human acceptance is reserved for phase end.
- VAL-079: PASS for task timing - the final exact ordinary gate reported `cmd/fid 196.303s`, under 600s; it reported `internal/workflow (cached)`, with the earlier task-14 exact gate's measured `internal/workflow 117.669s` retained above the recovery boundary. The final tagged gate reported `cmd/fid/integration 229.795s`, under 900s. All final gates terminated without a timeout. No package reached 80% of its budget.

Retained-row provenance: `TestPhase22AttentionReadBoundary` was introduced before Phase 23 (`a0a93c7c`; exact source-path inventory last amended in the task-10 era). The task-14 Human renderer created a new legitimate Attention source path; its exact path inventory was extended without weakening its assertion. The retained Phase 22 resume Human test was introduced at `43ae4c7` (`Refs: phase-22/task-12`); its old Findings-placeholder expectation was converged to the accepted real array. Both execute in the ordinary or tagged required gates respectively.

The V1 Concern list tests in `cmd/fid/concern_test.go` began at `f35b3febd7a9976c21be0370d3d701c6aec8e008` (`Refs: phase-9/task-1`) and were last amended at `8cad753ad612ea5fe252afe55b25ba0874039694` (`Refs: phase-21/task-20`). The accepted Phase 23 generic `concern.list` change displaced their no-filter direct-array assertion. `TestConcernListGenericReadExcludesV1Rows` now checks that generic output is an empty canonical JSON `items` array and a purpose-labeled Human view when only V1 rows exist; `TestConcernListRetainsIssueScopedV1ReadUntilRemoval` preserves the old Issue-scoped reader's status and empty-result checks. The separate V1 removal remains outside this task.

## Tests run

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestResumeReportsRealFindingsAndJudgmentSections$' -count=1 ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestResumeReportsRealFindingsAndJudgmentSections$' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	4.261s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/workflow` - FAIL - `FAIL` (terminated `exit_code: 1`); retained Attention path inventory named no new renderer.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	79.299s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestPhase22AttentionReadBoundary$' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.259s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23ConcernListHumanProjection$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); CLI returned exit 2 for `concern list --subject`.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23ConcernListHumanProjection$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); continuation reproduced the same `--subject` exit 2 before repair.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23ConcernListHumanProjection$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.098s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestConcernList' -count=1 ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`); the Phase 9 direct-array V1 assertions were displaced by generic no-filter output.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestConcernList' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.467s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23ConcernListHumanProjection$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); the added limit case revealed `invalid-argument` where the accepted bound requires `result-limit-exceeded`.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23ConcernListHumanProjection$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.101s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23ConcernListHumanProjection$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); the new disposal fixture lacked `concern.dispose` authority and returned `authority-denied`.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23ConcernListHumanProjection$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.833s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23ConcernListHumanProjection$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); mixing V1 `--issue` with generic `--severity` silently ignored the severity filter.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23ConcernListHumanProjection$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.258s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/review` - PASS - `ok  	github.com/mbell-nzl/fid/internal/review	1.086s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/review` - PASS - `ok  	github.com/mbell-nzl/fid/internal/review	1.085s` (terminated `exit_code: 0`)
- `FID_P23_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-23-human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23HumanFixtures$' -count=1 -v ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	20.787s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23ResumeRejectsJudgmentSectionOverflow$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	17.879s` (terminated `exit_code: 0`)
- Real `fid` fixture interaction: protected action denied with `protected-decision-required` (exit 5); `human-outcome record` succeeded (exit 0); repeated action succeeded (exit 0). The invoked absolute binary, workspaces, and arguments are in `qa/human-fixtures.md` under VAL-071.
- Real `fid` fixture interaction: exhausted review-cycle claim denied with `review-rework-exhausted` (exit 4); absolute arguments are in `qa/human-fixtures.md` under VAL-072.
- Real `fid` CLI interaction through the tagged binary: `--project <temporary workspace> --json concern list --subject issue:<seeded Issue ID>` returned a canonical `concern.list` envelope with `items[0].id=CON-human-list` and exit 0; the same invocation without `--json` returned a purpose-labeled Human Concern view without private prose. The test then exercised severity, status, subject, limit, and mixed-mode rejection through the same real entrypoint (`cmd/fid/integration/p23_concern_list_human_test.go:10`).

## Exact quality-gate results

Working directory: `/home/mark/projects/fid`. Each group below was run byte-identically in task-spec order. The original lint repair changed only `internal/workflow/resume.go`; its failed attempt remains here. The first continuation group passed before the final bounded V1-route correction, and the last group passed on the final tree.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - BLOCKED - `* unused: 1` (terminated `exit_code: 1`); first attempt found the former single-Obligation helper unused after resume began returning the full array.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	117.669s	coverage: 62.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					60.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`); `TestP23ConcernListHumanProjection` failed because `concern list --subject` returned exit 2 through the old reader. The package reported 226.499s.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					60.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	243.997s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					60.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	229.795s` (terminated `exit_code: 0`)

## Deviations from spec

- The prior task-14 recovery report is preserved at commit `eb6f4f110b30a72b4203d30605217b80613df776`, following implementation checkpoint `59d2653da7cfe14699dc4d1cd85ef032bec26a5c`. It recorded RECOVERY_REQUIRED, category `gate-repair`, attempt 1, and the condition that `concern list --subject` exited 2 because the accepted generic reader lay outside the task's named source paths. The original focused and required tagged failures remain verbatim in this report's test and gate sections.
- The original task-14 checkpoint found a scope gap: the accepted `concern.list` Human/JSON leaf still called the V1 Issue-only reader, so `--subject` returned exit 2. TRG-23-009 authorized the minimum same-task extension into `internal/review/concern_list.go` and `cmd/fid/concern.go`. The new reader implements the already accepted filters and projection. The original focused `FAIL` (`exit_code: 1`) and tagged `FAIL` (`exit_code: 1`) remain above; the later focused `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.258s` and final tagged `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	229.795s` both terminated `exit_code: 0`. No contract, command, or gate definition changed.
- The first continuation touched-package run exposed the old V1 no-filter direct-array assertions in `cmd/fid/concern_test.go` with `FAIL` (terminated `exit_code: 1`). The accepted generic no-filter read displaced those assertions, so the matching in-package tests now prove generic output excludes V1 rows while retaining the V1 `--issue` status and empty-result cases. The focused rerun passed with `ok  	github.com/mbell-nzl/fid/cmd/fid	2.467s` (terminated `exit_code: 0`). The V1 reader itself was not removed or renamed.
- The extended retained integration test first exposed the missing typed limit error: `FAIL` (terminated `exit_code: 1`) with `invalid-argument`; `renderConcernError` now maps `ResultLimitExceededError`, and the byte-identical focused command passed with `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.101s` (terminated `exit_code: 0`). A later mixed `--issue`/`--severity` case first returned a silent `[]` with `FAIL` (terminated `exit_code: 1`); the CLI now rejects that mix and preserves the V1 path's existing JSON encoding and issue error handling. The focused rerun passed with `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.258s` (terminated `exit_code: 0`).
- The first touched-package race run failed with `FAIL` (terminated `exit_code: 1`) because `TestPhase22AttentionReadBoundary`'s exact source-path inventory omitted the newly added `cmd/fid/p23_human.go`. The task-14 accepted Attention renderer invalidated that historical inventory; I added only the new path. The byte-identical touched-package command then passed with `ok  	github.com/mbell-nzl/fid/internal/workflow	79.299s` (terminated `exit_code: 0`).
- The first required lint attempt failed with `* unused: 1` (terminated `exit_code: 1`) because switching resume to full evaluations left `currentEnvelopeEvaluation` unused. Resume now also calls the retained helper for its primary next-action evaluation; the byte-identical lint command then passed with `0 issues.` (terminated `exit_code: 0`).
- The first committed-report lint attempt failed with `FAIL report-lint found 17 failure(s)` (terminated `exit_code: 17`) because 17 passing command-result lines placed explanations after the required termination field. I moved those explanations into report prose and committed the corrected report. The byte-identical report-lint command then returned `PASS report-lint completed with no failures` (terminated `exit_code: 0`) and `OUTCOME COMPLETE`; the original failed attempt remains recorded here.

## Concerns, technical debt, or notes for the Tech Lead

- TRG-23-009 resolved the task-14 reader scope gap. The generic list now uses `judgment_concerns`; the explicit `--issue` compatibility route still uses the V1 reader until its separately accepted removal. The window 13-15 QA reviewer should repeat its independent Human-text inspection on the delivered binary.
- `internal/workflow/resume.go` now exceeds the 400-line non-test target because the six bounded judgment projections live in the task's named file. A later refactor needs explicit file scope; no architectural decision was made here.
- The generated Human fixture manifest records the source worktree as dirty at generation time and binds the retained binary by digest. The guide instructs regeneration after a delivered binary changes; phase-end Human judgments remain pending.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `cdcc4f01fb2060d4fa5721a4108b9fda1dde328a`

## Commit subject

- `feat(review): complete generic Concern list projection`
