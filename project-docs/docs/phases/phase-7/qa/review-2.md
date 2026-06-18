# Phase 7 - Mid-Phase Review 2 (Tasks 3-4)

Reviewer: qa-reviewer
Date: 2026-06-14
Tasks covered: phase-7/task-3, phase-7/task-4

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-7/tasks/task-3.md`
  - `project-docs/docs/phases/phase-7/tasks/task-4.md`
- Completion reports:
  - `project-docs/docs/phases/phase-7/tasks/task-3-done.md`
  - `project-docs/docs/phases/phase-7/tasks/task-4-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-7/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-7/qa/qa-plan.md`
- Current escalation state:
  - `project-docs/docs/current/tl_escalations_current.md`
- Commits or changed files:
  - `1806aa4f45c6b53981cd33df6748a524cb139c58` - `feat(issue): add ac list command`
  - `a7cf2bdf9461918dbbe6e568488d0b2ccc418fba` - `test(e2e): add AC gating rebuild parity smoke`
  - `internal/issues/ac_list.go`
  - `internal/issues/ac_test.go`
  - `cmd/fid/issue.go`
  - `cmd/fid/issue_ac_test.go`
  - `cmd/fid/e2e_v01c_ac_test.go`
  - `internal/format/markdown.go`
  - `project-docs/specs/06_DESIGN_PLAN_CLI.md`

## Summary

- Blocking findings: 0
- Improvement findings: 2
- Escalations raised: none
- Covered validation assertions: VAL-006, VAL-007, VAL-008, VAL-009

No new open ESC was found. Baseline open escalation remains `ESC-005` only.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-006 | task-3 | ON TRACK | Task 3 report cites tests and transcript evidence for latest-per-role list output, empty list, missing issue, and parity with `issue get --section acceptance` (`task-3-done.md:24`). Code reuses `issueAcceptanceCriteria` after loading the issue owner role (`internal/issues/ac_list.go:49`). Focused test proves repeated developer fail->pass collapses to latest pass, retains `qa_lead`, orders by ordinal, and deep-equals `issue get` acceptance (`internal/issues/ac_test.go:1150`; `internal/issues/ac_test.go:1257`). |
| VAL-007 | task-3 | ON TRACK | Report includes real `/tmp/fid issue ac list` transcripts for omitted format, explicit JSON, TOON, markdown, YAML, and empty list (`task-3-done.md:28`). CLI test proves omitted `--format` equals JSON and structural parity for TOON/YAML plus markdown rendering (`cmd/fid/issue_ac_test.go:486`; `cmd/fid/issue_ac_test.go:554`). |
| VAL-008 | task-4 | ON TRACK | Task 4 report includes a real command transcript: ACs added with `fid issue ac add`, reject from `qa_review -> qa_pass` names both AC IDs and missing `qa_lead`, journal line count remains `7`, state is unchanged, then real `fid issue ac verify ... --status pass` commands allow the transition (`task-4-done.md:44`; `task-4-done.md:53`; `task-4-done.md:62`). Integration test asserts the same reject-without-append and success-after-pass path (`cmd/fid/e2e_v01c_ac_test.go:66`; `cmd/fid/e2e_v01c_ac_test.go:90`; `cmd/fid/e2e_v01c_ac_test.go:120`). |
| VAL-009 | task-4 | ON TRACK | Report snapshots `acceptance_criteria`, `ac_verifications`, `issue get --section acceptance`, and `ac list`, deletes `.fid/fid.sqlite`, runs `fid rebuild`, and records parity OK for all four comparisons (`task-4-done.md:70`; `task-4-done.md:80`; `task-4-done.md:97`). Integration test performs the same structural and latest-read-back comparisons (`cmd/fid/e2e_v01c_ac_test.go:153`; `cmd/fid/e2e_v01c_ac_test.go:168`). |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 3 required gates | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | PASS | Completion report records all required commands as PASS with package threshold excerpts: `cmd/fid 78.9%`, `internal/issues 82.0%`, `internal/journal 80.3%`, `internal/policy 85.0%`, `internal/workflow 95.5%` (`task-3-done.md:108`). |
| Task 3 conditional build | `go build -o /tmp/fid ./cmd/fid` | PASS | Triggered by real CLI transcripts and recorded as PASS (`task-3-done.md:121`). |
| Task 3 transcript obligations | Populated `ac list` per format plus empty-list case | PASS | Default JSON, explicit JSON, TOON, markdown, YAML, and empty list are pasted (`task-3-done.md:28`; `task-3-done.md:83`). Issue-not-found evidence is test-backed, not transcript-backed; see Improvement 1. |
| Task 4 required gates | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | PASS | Completion report records all required commands as PASS with required package threshold excerpts (`task-4-done.md:116`). |
| Task 4 conditional integration | `go test -tags integration -race ./...` | PASS | Triggered by `//go:build integration` test and recorded as PASS (`task-4-done.md:122`; `cmd/fid/e2e_v01c_ac_test.go:1`). |
| Task 4 conditional build | `go build -o /tmp/fid ./cmd/fid` | PASS | Triggered by the built-binary transcript and recorded as PASS (`task-4-done.md:123`). |
| Task 4 transcript obligations | Reject -> verify -> pass and delete-DB rebuild parity | PASS | Transcript includes real `ac add`, real `ac verify`, rejected transition with no journal append, successful transition after pass verifications, DB deletion, rebuild, and parity checks (`task-4-done.md:44`; `task-4-done.md:80`; `task-4-done.md:97`). |
| Commit verification | `git rev-parse --verify 1806aa4^{commit}`; `git show -s --format=%H%n%s 1806aa4` | PASS | Verified full hash `1806aa4f45c6b53981cd33df6748a524cb139c58` and subject `feat(issue): add ac list command`. |
| Commit verification | `git rev-parse --verify a7cf2bd^{commit}`; `git show -s --format=%H%n%s a7cf2bd` | PASS | Verified full hash `a7cf2bdf9461918dbbe6e568488d0b2ccc418fba` and subject `test(e2e): add AC gating rebuild parity smoke`. |
| Reviewer spot check | `go test ./internal/issues -run 'TestListAcceptanceCriteria(ReturnsOrdinalLatestPerRoleAndMatchesIssueGet\|EmptyAndNotFound)$'` | PASS | Local reviewer run returned `ok github.com/mbell-nzl/fid/internal/issues (cached)`. |
| Reviewer spot check | `go test ./cmd/fid -run 'TestIssueACListCommand(DefaultsToJSONAndMatchesIssueGetAcceptance\|DispatchesFormatsAndHandlesEmptyAndMissingIssues)$'` | PASS | Local reviewer run returned `ok github.com/mbell-nzl/fid/cmd/fid 0.119s`. |
| Reviewer spot check | `go test -tags integration ./cmd/fid -run TestV01CAllACsVerifiedGateAndRebuildParity -count=1` | PASS | Local reviewer run returned `ok github.com/mbell-nzl/fid/cmd/fid 0.129s`. |
| Escalation state | `project-docs/docs/current/tl_escalations_current.md` | PASS | `grep` found only baseline `ESC-005` under Open Escalations; no new open ESC appears. |

## Design / Parity Checks

| Task | Design sections checked | Result |
|---|---|---|
| Task 3 | `project-docs/specs/05_DESIGN_PLAN.md` §5.1 Format Strategy, §5.2 JSON Output Conventions, §5.5 Empty State Standard, §8 Cross-Surface Parity Matrix; `project-docs/specs/06_DESIGN_PLAN_CLI.md` §`fid issue ac list` and §`fid issue get`. | PASS. JSON default, four-format dispatch, empty JSON array, snake_case payload, and parity with `issue get --section acceptance` are covered by tests and transcripts. |
| Task 4 | `project-docs/specs/05_DESIGN_PLAN.md` §5.4 Error Message Standard, §5.6 Write Confirmation Standard, §8 Cross-Surface Parity Matrix; `project-docs/specs/06_DESIGN_PLAN_CLI.md` §`fid issue transition`, §`fid issue ac add`, §`fid issue ac verify`, and §`fid issue ac list`. | PASS. Reject output names rule/entities/missing role, write confirmations are present, and post-rebuild parity spans DB rows plus both read surfaces. |

## Per-Task Findings

### Task 3 - `fid issue ac list`

#### Blocking

- None.

#### Improvement

- [`project-docs/docs/phases/phase-7/tasks/task-3-done.md:83`](/home/mark/projects/fid/project-docs/docs/phases/phase-7/tasks/task-3-done.md:83) The report's pasted CLI transcripts cover all required formats and the empty-list case, while issue-not-found is covered by tests rather than a transcript. This is acceptable because the task only required transcript paste for populated formats and empty-list evidence, but final QA would be easier to audit with one pasted missing-issue stderr excerpt. Fix direction: optionally add or reproduce a missing-issue CLI excerpt in final QA evidence if public error-output evidence is consolidated there.

### Task 4 - `all_acs_verified` real-data gating + rebuild parity

#### Blocking

- None.

#### Improvement

- [`project-docs/docs/phases/phase-7/tasks/task-4-done.md:7`](/home/mark/projects/fid/project-docs/docs/phases/phase-7/tasks/task-4-done.md:7) The direct SQLite `task_completions` setup row is disclosed and stays outside the AC write evidence, but the report classifies this as no deviation without explicitly labeling it as a test fixture for bypassing Phase 8's out-of-scope completion surface. This does not block VAL-008/VAL-009 because the AC adds/verifications are real commands and the validation target is the `qa_review -> qa_pass` `all_acs_verified` edge. Fix direction: in final QA or wrap-up notes, preserve the rationale that the SQLite row is setup-only for an unrelated out-of-scope rule, not replacement evidence for AC journal writes.

## Cross-Cutting Observations

- Architecture and module boundaries look consistent with the phase scope: Task 3 keeps read behavior in `internal/issues`, routes CLI output through `internal/format`, and reuses the existing issue acceptance helper rather than duplicating latest-per-role logic.
- Security/privacy risk remains low in the reviewed slice. User-controlled IDs are passed through parameterized SQL in the read path, and Task 4's transcript uses temporary workspaces with synthetic issue/AC IDs.
- Required command coverage is plausible and complete for both reports. Task 4's integration command was correctly triggered because the new smoke is build-tagged `integration`.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-7/qa/review-2.md
Blocking: 0  Improvement: 2  Escalations: none
VERDICT: CONTINUE
```
