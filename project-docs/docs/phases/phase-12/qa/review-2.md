# Phase 12 - Mid-Phase Review 2 (Tasks 3-4)

Reviewer: QA Reviewer
Date: 2026-06-27
Tasks covered: phase-12/task-3 through phase-12/task-4

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-12/tasks/task-3.md`
  - `project-docs/docs/phases/phase-12/tasks/task-4.md`
- Completion reports:
  - `project-docs/docs/phases/phase-12/tasks/task-3-done.md`
  - `project-docs/docs/phases/phase-12/tasks/task-4-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-12/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-12/qa/qa-plan.md`
- Prior review and triage context:
  - `project-docs/docs/phases/phase-12/qa/review-1.md`
  - `project-docs/docs/phases/phase-12/triage-log.md`
  - `project-docs/docs/current/tl_escalations_current.md`
  - `project-docs/docs/current/decisions_active.md`
- Commits or changed files:
  - `66392c4c8d4356dd435b1ca37e30e3626aa1d2f5` (`feat(links): add link remove command`) verified with `git rev-parse --verify`.
  - `e30e40ee96c2089607dd9f3a4847406874e32a8e` (`feat(links): add link list command`) verified with `git rev-parse --verify`.
  - Spot-checked changed files under `cmd/fid`, `internal/links`, `internal/format`, and the reviewed task reports.

## Summary

- Blocking findings: 0
- Improvement findings: 2
- Escalations raised: none
- Covered validation assertions: VAL-009, VAL-010, VAL-011

Task 3 report verdict: complete for Review 2. Required commands are recorded, the CLI build is recorded as PASS, conditional integration is explicitly N/A with a reason, and evidence tracks for VAL-009.

Task 4 report verdict: complete for Review 2. Required commands are recorded, integration is recorded as PASS, the CLI build is recorded as PASS, benchmark execution is explicitly N/A with a reason, and evidence tracks for VAL-010 and VAL-011.

No blockers were found. No new open ESC is needed. VAL-012 is not a Task 3-4 deterministic task assertion; it remains a phase-end human checklist item after deterministic QA passes.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-009 | task-3 | ON TRACK | Task 3 report lines 33-39 cite declared + inverse deletion, one event, `link_removed` activity evidence, both-direction confirmation, missing-link rejection, and unknown-role rejection. Code spot-checks show `Remove` uses `EnsureProjectionFreshWithLock`, emits two `links` delete effects plus one `activity_log` effect, and rejects missing links before append. Review-run check passed: `go test ./internal/links -run 'TestRemoveDeletesDeclaredAndInverseRowsAndProjectsActivity|TestRemoveRejectsMissingLinkBeforeAppend|TestRemoveSupportsRequirementInverse'` -> `ok github.com/mbell-nzl/fid/internal/links`. |
| VAL-010 | task-4 | ON TRACK | Task 4 report lines 35-42 cite outgoing and incoming rows with `direction`, `relation`, `entity_type`, `entity_id`, resolved `entity_title`, and json/toon/markdown/yaml parity. Review-run check passed: `go test ./cmd/fid -run TestLinkListReturnsOutgoingAndIncomingLinksWithResolvedTitlesAndFormatParity` -> `ok github.com/mbell-nzl/fid/cmd/fid`. |
| VAL-011 | task-4 | ON TRACK | Task 4 report lines 43-49 cite the real-command sequence, one journal event per write, no nested-lock failure, and rebuild parity for `links` / `requirements`. Review-run check passed: `go test -tags integration ./cmd/fid -run TestV01DLinkFlowAndRebuildParity` -> `ok github.com/mbell-nzl/fid/cmd/fid`. |
| VAL-012 | phase-end human checklist | DEFERRED | The validation contract assigns VAL-012 to the Human validator at phase end. Neither Task 3 nor Task 4 claims it, and no deterministic task report evidence is required for it in this review window. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 3 format | `gofmt -l cmd internal` | PASS | Recorded in task-3-done line 55. |
| Task 3 lint | `golangci-lint run ./...` | PASS | Recorded in task-3-done line 56 as `0 issues.` |
| Task 3 race/coverage test | `go test -race -coverprofile=coverage.out ./...` | PASS | Recorded in task-3-done line 57. |
| Task 3 coverage report | `go tool cover -func=coverage.out` | PASS | Recorded in task-3-done lines 58-63 with total 76.6%, `cmd/fid` 78.3%, `internal/links` 79.5%, and named domain packages above the cited thresholds. |
| Task 3 build | `go build -o /tmp/fid ./cmd/fid` | PASS | Required for the new CLI subcommand; recorded in task-3-done line 64. |
| Task 3 conditional integration | `go test -tags integration -race ./...` | N/A | Recorded in task-3-done lines 65-66 with reason: no integration fixtures were touched. |
| Task 3 deterministic spot check | `go test ./internal/links -run 'TestRemoveDeletesDeclaredAndInverseRowsAndProjectsActivity|TestRemoveRejectsMissingLinkBeforeAppend|TestRemoveSupportsRequirementInverse'` | PASS | Review-run representative check for remove success, missing-link rejection, requirement inverse removal, and activity evidence. |
| Task 3 design/parity check | `05_DESIGN_PLAN.md` section 3.1 command tree; `06_DESIGN_PLAN_CLI.md` section 7 `fid link remove` | PASS | The design-plan command tree includes `link remove <from> <relation> <to>` as declared + inverse removal, and the CLI spec names the one-event remove plus stderr confirmation. Task 3 report lines 25-28 says no CLI spec change was required. |
| Task 4 format | `gofmt -l cmd internal` | PASS | Recorded in task-4-done line 66. |
| Task 4 lint | `golangci-lint run ./...` | PASS | Recorded in task-4-done line 67 as `0 issues.` |
| Task 4 race/coverage test | `go test -race -coverprofile=coverage.out ./...` | PASS | Recorded in task-4-done line 68. |
| Task 4 coverage report | `go tool cover -func=coverage.out` | PASS | Recorded in task-4-done lines 69-74 with total 76.2%, `cmd/fid` 77.9%, `internal/links` 79.2%, and named domain packages above the cited thresholds. |
| Task 4 integration | `go test -tags integration -race ./...` | PASS | Required by the task spec and QA plan for VAL-011; recorded in task-4-done line 75. |
| Task 4 build | `go build -o /tmp/fid ./cmd/fid` | PASS | Required for the new CLI subcommand; recorded in task-4-done line 76. |
| Task 4 conditional benchmark | `go test -bench . -benchmem ./...` | N/A | Recorded in task-4-done line 77 with reason: benchmark/scale fixtures were not touched. |
| Task 4 list parity spot check | `go test ./cmd/fid -run TestLinkListReturnsOutgoingAndIncomingLinksWithResolvedTitlesAndFormatParity` | PASS | Review-run representative check for VAL-010 four-format structural parity. |
| Task 4 e2e spot check | `go test -tags integration ./cmd/fid -run TestV01DLinkFlowAndRebuildParity` | PASS | Review-run representative check for VAL-011 real-command flow and rebuild parity. |
| Commit verification | `git rev-parse --verify <hash>^{commit}` | PASS | Both task report commit hashes resolve to commits. |
| Changed-file scope | `git show --stat --oneline <hash>` | PASS | Task 3 changed `cmd/fid/link.go`, `cmd/fid/link_test.go`, `internal/links/*`, and its report. Task 4 changed `cmd/fid/e2e_v01d_link_test.go`, `cmd/fid/link.go`, `cmd/fid/link_test.go`, `internal/format/markdown.go`, `internal/links/*`, and its report. No migration file appears in either commit stat. |
| Escalation state | `project-docs/docs/current/tl_escalations_current.md` | PASS | Open escalations section is `None`; this review raises none. |

## Per-Task Findings

### Task 3 - `fid link remove`

#### Blocking

- None.

#### Improvement

- `project-docs/docs/phases/phase-12/tasks/task-3-done.md:25` The report says the CLI spec needed no change, but it does not name the specific design sections checked in polish/parity review mode. Fix direction: future reports should cite `05_DESIGN_PLAN.md` section 3.1 and `06_DESIGN_PLAN_CLI.md` section 7 `fid link remove` when recording verified-no-change design parity.

### Task 4 - `fid link list` And End-To-End Flow

#### Blocking

- None.

#### Improvement

- `project-docs/docs/phases/phase-12/tasks/task-4-done.md:26` The report says the CLI spec needed no change, but it does not name the specific design sections checked in polish/parity review mode. Fix direction: future reports should cite `05_DESIGN_PLAN.md` section 3.1 / 5.1-5.2 and `06_DESIGN_PLAN_CLI.md` section 7 `fid link list` when recording four-format output parity.

## Cross-Cutting Observations

- Review 1 improvement notes were triaged as `RESUME` in `project-docs/docs/phases/phase-12/triage-log.md`; Review 2 found no evidence that those notes became blockers for Tasks 3-4.
- The current worktree already contains modified task done reports and untracked Phase 12 QA artifacts from prior phase work. This QA review did not modify task reports, task specs, decisions, triage, or product code. The requested `review-2.md` is the only file authored by this review.
- Task 3's objective prose says "project both removals to `activity_log`," while the validation contract, QA plan, and task acceptance criteria require activity evidence as one removal row/effect. Because this review uses the QA plan and validation contract as source of truth, the one-row `link_removed` evidence is acceptable and no ESC is raised.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-12/qa/review-2.md
Blocking: 0  Improvement: 2  Escalations: none
VERDICT: CONTINUE
```
