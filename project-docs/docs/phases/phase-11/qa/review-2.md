# Phase 11 - Mid-Phase Review 2 (tasks 2-2)

Reviewer: qa-reviewer
Date: 2026-06-21
Tasks covered: phase-11/task-2

## Verdict

CONTINUE

## Scope Reviewed

- Task spec:
  - `project-docs/docs/phases/phase-11/tasks/task-2.md`
- Completion report:
  - `project-docs/docs/phases/phase-11/tasks/task-2-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-11/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-11/qa/qa-plan.md`
- Current escalation register:
  - `project-docs/docs/current/tl_escalations_current.md`
- Design-plan sections checked:
  - `project-docs/specs/05_DESIGN_PLAN.md` section 3.1 CLI Command Tree, lines 45-58 and 101-104.
  - `project-docs/specs/05_DESIGN_PLAN.md` sections 5.3, 5.4, and 5.6, lines 190-255.
  - `project-docs/specs/05_DESIGN_PLAN.md` section 8 Cross-Surface Parity Matrix, lines 406-415.
- Commit reviewed:
  - `fbb8c2f0f10dbdb865c99ffa7e1fadae29fcfd41` `feat(issue): add claim and unclaim commands`

## Summary

- Blocking findings: 0
- Improvement findings: 1
- Escalations raised: none
- Covered validation assertions: VAL-002, VAL-003, VAL-004, VAL-005, VAL-006, VAL-007, VAL-010

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-002 | task-2 | ON TRACK | Completion report claims default/custom TTL and note evidence at `task-2-done.md` lines 27 and 37-41. CLI tests assert default TTL `1800`, one current row, expiry, and empty stdout in `cmd/fid/issue_claim_test.go` lines 13-57; custom TTL and note read-back in lines 186-216. |
| VAL-003 | task-2 | ON TRACK | Completion report claims live conflict rejection at `task-2-done.md` line 28. Test coverage asserts non-zero error, empty stdout, holder/expiry text, and unchanged journal line count in `cmd/fid/issue_claim_test.go` lines 267-300. |
| VAL-004 | task-2 | ON TRACK | Completion report claims CLI stale takeover at `task-2-done.md` line 29. Task 2 CLI evidence proves stale takeover without `--force` and current holder read-back in `cmd/fid/issue_claim_test.go` lines 422-447; domain evidence for superseded previous holder and expiry remains covered by `internal/claims/claims_test.go` lines 397-435 and was included in Task 2's package test runs. |
| VAL-005 | task-2 | ON TRACK | Completion report claims force-acquire allow/deny and policy parity at `task-2-done.md` line 30. CLI tests cover previous-holder confirmation and new holder in `cmd/fid/issue_claim_test.go` lines 155-184, denied force with `issue.claim.force_acquire` and zero append in lines 302-333, and `fid policy show` parity in lines 449-488. |
| VAL-006 | task-2 | ON TRACK | Completion report claims renew preserve/override/non-holder evidence at `task-2-done.md` line 31. CLI tests assert TTL preservation, note preservation, claimed/expires timestamp changes, and `activity_log` increment in `cmd/fid/issue_claim_test.go` lines 59-123; TTL override and non-holder zero-append rejection are covered in lines 336-392. Domain event/activity shape is covered in `internal/claims/claims_test.go` lines 620-650. |
| VAL-007 | task-2 | ON TRACK | Completion report claims release and reject-path evidence at `task-2-done.md` line 32. CLI tests assert holder release removes the claim and appends one journal event in `cmd/fid/issue_claim_test.go` lines 125-152; wrong-holder, repeat, missing-role, and unknown-role rejection paths preserve journal line counts in lines 490-617. |
| VAL-010 | task-2 | DEFERRED PARTIAL | Task 2 covers the human claim persistence/read-back slice: completion report line 33 and `cmd/fid/issue_claim_test.go` lines 218-237. QA plan assigns the transition-blocking half to Tasks 2-4 and Reviews 2 and 4; Task 4 must still prove non-holder agent transition rejection while `human` holds the claim. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task spec exact commands | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go build -o /tmp/fid ./cmd/fid` | PASS | Required by task spec lines 106-123 and QA plan lines 126-132; completion report records all six as PASS at `task-2-done.md` lines 54-59. |
| Conditional benchmark command | `go test -bench . -benchmem ./...` | N/A ACCEPTED | Task spec lines 125-132 require this only if benchmark fixtures or journal scale helpers are touched. Completion report line 60 records N/A because those areas were not touched; changed-file list confirms no benchmark fixtures or journal scale helpers. |
| Coverage reporting | `go tool cover -func=coverage.out` package percentages | PASS | Completion report line 57 records `cmd/fid` 78.1%, `internal/issues` 81.9%, `internal/claims` 83.4%, `internal/workflow` 95.5%, `internal/policy` 81.2%, and `internal/journal` 80.3%, satisfying COV-001/COV-002 thresholds from `08_QUALITY_PLAN.md`. |
| Commit hash verification | `git rev-parse --verify fbb8c2f^{commit}` | PASS | Verified as `fbb8c2f0f10dbdb865c99ffa7e1fadae29fcfd41`; completion report lines 74-80 match the short hash and subject. |
| Commit subject/body | `git show --no-patch --format=%H%n%s%n%b fbb8c2f` | PASS | Subject is `feat(issue): add claim and unclaim commands`; body includes `Refs: phase-11/task-2`. |
| Changed files | `git show --name-only --format= fbb8c2f` | PASS | Files are limited to `cmd/fid/issue.go`, `cmd/fid/issue_claim_test.go`, `internal/claims/*`, and `task-2-done.md`, matching the task scope. |
| Representative deterministic spot-check | `go test ./cmd/fid -run 'TestIssueClaimCommand|TestIssueUnclaimCommand' -count=1` | PASS | Local rerun returned `ok github.com/mbell-nzl/fid/cmd/fid 1.512s`. |
| Escalation baseline | Re-read `project-docs/docs/current/tl_escalations_current.md` | PASS | Open escalations remain `None` at lines 41-43. No new cross-task, architecture, security, or acceptance ambiguity was found. |

## Per-Task Findings

### Task 2 - `fid issue claim` And `fid issue unclaim`

#### Blocking

- None.

#### Improvement

- [cmd/fid/issue_claim_test.go:429] The CLI stale-takeover test waits on real time with `time.Sleep(1100 * time.Millisecond)` after a `1s` TTL. This is acceptable for this task because the focused test and recorded gates passed, but it is slower and more timing-sensitive than the quality plan's clock-control preference for time-dependent behavior. Fix direction: when the claims test helpers next move, add a deterministic stale-claim fixture or test clock path for CLI stale-takeover evidence.

## Cross-cutting observations

- Task 2 satisfies the QA plan's per-task Review 2 cadence. The completion report's exact gate list matches both the task spec and QA plan, and the conditional benchmark omission is explicitly reasoned rather than silently skipped.
- VAL-010 remains intentionally partial in this window. Task 4 must still provide the `claim list`/transition rejection proof that a non-holder agent cannot advance a human-held issue.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-11/qa/review-2.md
Blocking: 0  Improvement: 1  Escalations: none
VERDICT: CONTINUE
```
