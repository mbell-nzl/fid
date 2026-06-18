# Phase 7 - Mid-Phase Review 1 (Tasks 1-2)

Reviewer: QA Reviewer
Date: 2026-06-14
Tasks covered: phase-7/task-1 through phase-7/task-2

## Verdict

BLOCK

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-7/tasks/task-1.md`
  - `project-docs/docs/phases/phase-7/tasks/task-2.md`
- Completion reports:
  - `project-docs/docs/phases/phase-7/tasks/task-1-done.md`
  - `project-docs/docs/phases/phase-7/tasks/task-2-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-7/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-7/qa/qa-plan.md`
- Commits or changed files:
  - `59f02df` / `59f02dfeccbb636887cfa62ee3e73661a1c74317` - `feat(issue): add acceptance criteria command`
  - `437d093` / `437d0935d7a9f88422731a045158d12ca70fe262` - `feat(issue): add ac verify command`
  - `ddc6059` / `ddc6059427bce4418bfae15b353b7944bc500e09` - `docs(phase-7): finalize task-2 report`
  - `cmd/fid/issue.go`
  - `cmd/fid/issue_ac_test.go`
  - `internal/issues/ac.go`
  - `internal/issues/ac_test.go`

## Summary

- Blocking findings: 2
- Improvement findings: 1
- Escalations raised: ESC-005
- Covered validation assertions: VAL-001, VAL-002, VAL-003, VAL-004, VAL-005

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-1 | ON TRACK | Task 1 report cites engine/CLI tests and includes real CLI transcript for `ac add` plus `issue get --section acceptance --format json` showing `AC-` id, text, ordinal, and read-back. Focused reviewer check `go test ./internal/issues ./cmd/fid` passed. |
| VAL-002 | task-1 | ON TRACK | Task 1 report includes task success and epic rejection transcript with unchanged journal line count; engine/CLI tests cover unknown role, invalid ordinal, issue-not-found, epic rejection, unchanged journal lines, and zero AC rows. |
| VAL-003 | task-2 | ON TRACK | Task 2 report includes `developer=pass` and `qa_lead=fail` transcript/read-back; tests cover independent latest-per-role retention. |
| VAL-004 | task-2 | NOT ON TRACK | Same-role fail->pass evidence depends on sleeping past a one-second timestamp boundary. Implementation records `verified_at` as RFC3339 seconds and latest selection breaks ties by random `ACV-` id, so fast fail->pass writes can select the wrong latest row. See ESC-005. |
| VAL-005 | task-2 | ON TRACK | Task 2 report cites field-persistence and reject-before-append tests, includes pass verify and invalid-status CLI transcript, records `go build` for public-output transcript generation, and tests cover invalid status/evidence type, AC not found, AC/issue mismatch, unknown role, unchanged journal lines, and zero verification rows. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 1 required headings | Completion report sections | PASS | Core required sections exist: summary, files, docs, VAL coverage, tests, exact gates, deviations, concerns, dependency note, commit hash, commit subject. |
| Task 2 required headings | Completion report sections | PASS | Core required sections exist: summary, files, docs, VAL coverage, tests, exact gates, deviations, concerns, dependency note, commit hash, commit subject. |
| Task 1 exact commands | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | PASS | All required commands are recorded as `PASS`; no narrower substitute is used for the exact gate. |
| Task 2 exact commands | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | PASS | All required commands are recorded as `PASS`; no narrower substitute is used for the exact gate. |
| Conditional build command | `go build -o /tmp/fid ./cmd/fid` | PASS | Both reports include real public-output transcripts and record the conditional build as `PASS`. |
| Coverage thresholds | Task 1 and Task 2 `go tool cover` excerpts | PASS | Task 1 reports `cmd/fid` 79.1% and `internal/issues` 82.8%; Task 2 reports `cmd/fid` 79.2% and `internal/issues` 82.6%, above the task thresholds. |
| Commit verification | `git rev-parse --verify 59f02df^{commit}`; `git rev-parse --verify 437d093^{commit}` | PASS | Hashes resolve to `59f02dfeccbb636887cfa62ee3e73661a1c74317` and `437d0935d7a9f88422731a045158d12ca70fe262`; subjects match the reports and both commits include `Refs: phase-7/task-*`. |
| Task 2 report finalization commit | `git show ddc6059` | PASS | Follow-up docs commit finalizes Task 2 report hash/subject and includes `Refs: phase-7/task-2`. |
| Task 1 report finalization | `git status --short`; `git show 59f02df:.../task-1-done.md` | FAIL | Current report has the hash/subject filled in, but that edit is uncommitted; the committed report still says `Pending git commit.` |
| Public-output transcript evidence | Task 1 and Task 2 completion reports | PASS | Task 1 includes add and epic-reject transcripts; Task 2 includes pass verify, multi-role, fail->pass, and invalid-status transcripts. |
| Validate-before-append evidence | Task 1 and Task 2 tests/reports | PASS | Task evidence covers rejected writes with unchanged journal lines and zero projected rows for relevant AC/verification tables. |
| Append-only evidence | Task 2 tests/report | FAIL | Rows are retained, but latest selection is not robust for same-second same-role writes; this blocks VAL-004. |
| Concerns/debt/notes | Completion report sections | PASS | Both reports include non-empty sections with `None.` |
| Dependency budget | Completion report sections | PASS | Both reports include `No dependency changes.` |
| Open escalations before review | `tl_escalations_current.md` | PASS | No open ESCs existed at review start. ESC-005 was raised by this review. |
| Reviewer spot-check | `go test ./internal/issues ./cmd/fid` | PASS | Output: `ok github.com/mbell-nzl/fid/internal/issues (cached)` and `ok github.com/mbell-nzl/fid/cmd/fid 2.403s`. |

## Per-Task Findings

### Task 1 - `fid issue ac add`

#### Blocking

- [task-1-done.md:94] The current report's commit hash/subject edit is not committed; `git show 59f02df:project-docs/docs/phases/phase-7/tasks/task-1-done.md` still shows `Pending git commit`, while the worktree has the finalized hash/subject as an uncommitted modification. This violates the completion-report commit requirement and leaves the authoritative repo artifact stale. Fix direction: commit the Task 1 report finalization or otherwise make the committed report contain the verified hash and subject before phase sign-off.

#### Improvement

- [internal/issues/ac.go:431] `internal/issues/ac.go` is now 431 lines, above the 400-line non-test file tripwire and contrary to Task 1's note to keep the new file under 400 lines. This is not blocking VAL-001..VAL-005, but it should be split before the file becomes the shared home for Task 3/4 read and workflow additions. Fix direction: move AC verification helpers or write-path helpers into a second cohesive `internal/issues` file.

### Task 2 - `fid issue ac verify`

#### Blocking

- [internal/issues/ac.go:256] `VerifyAcceptanceCriterion` records `verified_at` with `time.RFC3339` second precision, while [detail.go:165] selects latest verification by `verified_at DESC, id DESC`. Because `ACV-` ids are generated hashes rather than append-order keys, two same-role verifications in the same second can return the lexicographically larger id instead of the later append. The tests avoid this by sleeping past the timestamp boundary at [ac_test.go:654] and [issue_ac_test.go:310], so VAL-004 is not reliably proven and later `ac list` / `all_acs_verified` consumers inherit the ambiguity. Fix direction: define and implement a deterministic append-order tie-break or higher-resolution timestamp, then add fail->pass coverage without `time.Sleep`.

#### Improvement

- None.

## Cross-cutting observations

- The code keeps AC write behavior in `internal/issues` and CLI rendering in `cmd/fid`, matching the project boundary guidance.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` §3 matches the implemented `ac add` and `ac verify` command syntax and stderr examples reviewed for Tasks 1-2.
- DEC-020 scope discipline is preserved in Task 1 evidence: AC add writes journal facts without adding `activity_log` rows.

## Escalations

- ESC-005 (QA Reviewer triggers 5 and 6): AC verification latest ordering can be nondeterministic for same-second same-role writes, affecting Task 2 VAL-004 and later Task 3/4 consumers.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-7/qa/review-1.md
Blocking: 2  Improvement: 1  Escalations: ESC-005
VERDICT: BLOCK
```
