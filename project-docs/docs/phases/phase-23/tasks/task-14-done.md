# Phase 23 - Task 14: Human review and judgment projections - Recovery checkpoint

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- OUTCOME RECOVERY_REQUIRED
- Category: gate-repair
- Condition: The accepted `concern.list` leaf still uses the V1 Issue-only reader. `TestP23ConcernListHumanProjection` fails on `--subject`, and the required tagged integration gate exits 1. Task 14 names the CLI Human renderer but does not name the missing `internal/review` Concern list reader in its source scope.
- Recovery attempts: 1
- Next action: Obtain a bounded same-task scope directive for the generic Concern list reader, implement the accepted `public-interface.md` section 2 filter and projection, keep the red assertion, rerun focused red/green and all five exact commands, then finish this report. Do not change the accepted contract or route this test to another task.
- Human attention required: no

## Summary of what was implemented

- Replaced resume's Findings placeholder with durable open and addressed Findings; projected Reviews, Concerns, Deviations, Decision References, Human Outcomes, protected Decisions, full per-Obligation evaluations, and guarded next actions. A section exceeding 100 items returns `result-limit-exceeded`.
- Added purpose-labeled Human renderers for the available Phase 23 Review and judgment leaves, including explicit Builder-assertion wording for addressed Findings, non-proof wording for judgment facts, source labels/paths/digests instead of project prose, and Stop/Attention qualifiers.
- Added a tagged Human fixture generator and verbatim guide for mixed Finding dispositions, a protected Decision with an activated Human Outcome actor, and review-cycle exhaustion. The Human judgments remain phase-end work.
- Preserved the red generic Concern list integration test and the failing required tagged gate as the scope-gap evidence.

## Files created or modified

- `internal/workflow/resume.go`, `internal/workflow/resume_actions.go`
- `cmd/fid/resume.go`, `cmd/fid/resume_human.go`, `cmd/fid/p23_human.go`, `cmd/fid/stop_human.go`
- `cmd/fid/review.go`, `cmd/fid/finding.go`, `cmd/fid/concern.go`, `cmd/fid/deviation.go`, `cmd/fid/decision.go`, `cmd/fid/human_outcome.go`, `cmd/fid/requirement.go`
- `cmd/fid/p23_human_test.go`, `cmd/fid/task18_completion_resume_test.go`, `cmd/fid/contract_surface_test.go`
- `cmd/fid/integration/p23_human_fixtures_test.go`, `cmd/fid/integration/p23_resume_limit_test.go`, `cmd/fid/integration/p23_concern_list_human_test.go`
- `cmd/fid/integration/p22_resume_human_test.go`, `p23_review_record_test.go`, `p23_finding_disposition_test.go`, `p23_concern_deviation_test.go`, `p23_decision_reference_test.go`, `p23_human_outcome_test.go`, `p23_review_evidence_budget_test.go`
- `project-docs/docs/phases/phase-23/qa/human-fixtures.md`, this report

## Documentation updated

- `project-docs/docs/phases/phase-23/qa/human-fixtures.md` records the retained binary SHA-256, absolute workspace and invocation arguments, verbatim Human stdout, JSON counterparts, and VAL-070/071/072 action steps.
- This report records the scope gap, validation limits, and immutable failed attempts.

## Validation Contract Coverage

- VAL-058: BLOCKED - `TestPhase23HumanLeafTextMatchesRecordedFacts` (`cmd/fid/p23_human_test.go:11`) covers the renderer shapes and prose exclusion; the generated Review, Finding, resume, Decision, Stop, and Attention views are in `qa/human-fixtures.md`. `TestP23ConcernListHumanProjection` (`cmd/fid/integration/p23_concern_list_human_test.go:10`) fails because the generic `concern.list` entrypoint remains unavailable, so every-leaf coverage is incomplete.
- VAL-059: PASS for the implemented resume projection - `TestResumeReportsRealFindingsAndJudgmentSections` (`cmd/fid/task18_completion_resume_test.go:201`) checks the new judgment arrays; `TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess` (`cmd/fid/integration/p23_human_outcome_test.go:39`) checks all four ordered Obligation evaluations; `TestP23HumanFixtures` produces a populated mixed-Finding resume inspection recorded in `qa/human-fixtures.md`; `TestP23ResumeRejectsJudgmentSectionOverflow` (`cmd/fid/integration/p23_resume_limit_test.go:11`) records 101 Concerns and checks typed `result-limit-exceeded` at limit 100.
- VAL-060: PASS for the implemented views - `TestPhase23HumanLeafTextMatchesRecordedFacts` rejects private prose sentinels and checks judgment non-proof language; `TestP23HumanFixtures` compares verbatim retained Human output after generator exit. Full every-leaf validation remains limited by VAL-058.
- VAL-070: FIXTURE READY, HUMAN PENDING - `TestP23HumanFixtures` (`cmd/fid/integration/p23_human_fixtures_test.go:46`) seeded resolved, addressed, and still-blocking Findings; `qa/human-fixtures.md` gives absolute Review/Finding/resume invocations. The Human has not judged comprehension.
- VAL-071: FIXTURE READY, HUMAN PENDING - the same generator seeded a protected Decision and an activated `human` actor. A real `fid` invocation returned `protected-decision-required` with exit 5, Human Outcome recording exited 0, and retry exited 0. Exact action arrays are in the fixture guide. The Human has not judged clarity.
- VAL-072: FIXTURE READY, HUMAN PENDING - the same generator seeded exhausted `review-cycles`, a Stop, and Attention; a real `fid` denial returned `review-rework-exhausted` with exit 4. Exact inspection/action arrays are in the fixture guide. The Human has not judged clarity.
- VAL-073: ADVANCED, HUMAN PENDING - purpose lines and labeled sections are asserted by `TestPhase23HumanLeafTextMatchesRecordedFacts` and recorded verbatim in the fixture guide; Human acceptance is reserved for phase end.
- VAL-079: PASS for timing only - the exact ordinary gate reported `cmd/fid 264.201s` and `internal/workflow 117.669s`, both under 600s; the exact tagged gate reported `cmd/fid/integration 226.499s`, under 900s. All terminated without a timeout. The tagged gate remains red for VAL-058, so task completion is not claimed.

Retained-row provenance: `TestPhase22AttentionReadBoundary` was introduced before Phase 23 (`a0a93c7c`; exact source-path inventory last amended in the task-10 era). The task-14 Human renderer created a new legitimate Attention source path; its exact path inventory was extended without weakening its assertion. The retained Phase 22 resume Human test was introduced at `43ae4c7` (`Refs: phase-22/task-12`); its old Findings-placeholder expectation was converged to the accepted real array. Both execute in the ordinary or tagged required gates respectively.

## Tests run

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestResumeReportsRealFindingsAndJudgmentSections$' -count=1 ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestResumeReportsRealFindingsAndJudgmentSections$' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	4.261s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/workflow` - FAIL - `FAIL` (terminated `exit_code: 1`); retained Attention path inventory named no new renderer.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	79.299s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestPhase22AttentionReadBoundary$' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.259s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23ConcernListHumanProjection$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); CLI returned exit 2 for `concern list --subject`.
- `FID_P23_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-23-human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23HumanFixtures$' -count=1 -v ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	20.787s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23ResumeRejectsJudgmentSectionOverflow$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	17.879s` (terminated `exit_code: 0`)
- Real `fid` fixture interaction: protected action denied with `protected-decision-required` (exit 5); `human-outcome record` succeeded (exit 0); repeated action succeeded (exit 0). The invoked absolute binary, workspaces, and arguments are in `qa/human-fixtures.md` under VAL-071.
- Real `fid` fixture interaction: exhausted review-cycle claim denied with `review-rework-exhausted` (exit 4); absolute arguments are in `qa/human-fixtures.md` under VAL-072.

## Exact quality-gate results

Working directory: `/home/mark/projects/fid`. The commands below were run byte-identically in task-spec order. The lint repair changed only `internal/workflow/resume.go`; the failed lint attempt remains here.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - BLOCKED - `* unused: 1` (terminated `exit_code: 1`); first attempt found the former single-Obligation helper unused after resume began returning the full array.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	117.669s	coverage: 62.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					60.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`); `TestP23ConcernListHumanProjection` failed because `concern list --subject` returned exit 2 through the old reader. The package reported 226.499s.

## Deviations from spec

- The accepted `concern.list` Human/JSON leaf remains incomplete because its current `cmd/fid/concern.go` entrypoint calls V1 `issues.ListConcerns` with only `--issue` and `--status`. The required generic reader belongs in `internal/review`, outside task 14's named source files. `public-interface.md` section 2 requires `--subject`, `--status`, `--severity`, and `--limit`. No assertion or gate was narrowed; the red test and full tagged failure are retained for a bounded same-task continuation.
- The first touched-package race run failed with `FAIL` (terminated `exit_code: 1`) because `TestPhase22AttentionReadBoundary`'s exact source-path inventory omitted the newly added `cmd/fid/p23_human.go`. The task-14 accepted Attention renderer invalidated that historical inventory; I added only the new path. The byte-identical touched-package command then passed with `ok  	github.com/mbell-nzl/fid/internal/workflow	79.299s` (terminated `exit_code: 0`).
- The first required lint attempt failed with `* unused: 1` (terminated `exit_code: 1`) because switching resume to full evaluations left `currentEnvelopeEvaluation` unused. Resume now also calls the retained helper for its primary next-action evaluation; the byte-identical lint command then passed with `0 issues.` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- The task-14 file scope omits the generic Concern list reader needed by the accepted every-leaf Human contract. The test's subject and strength are unchanged; the scope gap needs a bounded same-task directive before repair.
- `internal/workflow/resume.go` now exceeds the 400-line non-test target because the six bounded judgment projections live in the task's named file. A later refactor needs explicit file scope; no architectural decision was made here.
- The generated Human fixture manifest records the source worktree as dirty at generation time and binds the retained binary by digest. The guide instructs regeneration after a delivered binary changes; phase-end Human judgments remain pending.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
