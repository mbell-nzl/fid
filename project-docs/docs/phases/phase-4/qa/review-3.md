# Phase 4 - Mid-Phase Review 3 (Tasks 1-3)

Reviewer: QA Reviewer
Date: 2026-06-07
Tasks covered: phase-4/task-1 through phase-4/task-3

## Verdict

BLOCK

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-4/tasks/task-1.md`
  - `project-docs/docs/phases/phase-4/tasks/task-2.md`
  - `project-docs/docs/phases/phase-4/tasks/task-3.md`
- Completion reports:
  - `project-docs/docs/phases/phase-4/tasks/task-1-completion.md`
  - `project-docs/docs/phases/phase-4/tasks/task-2-completion.md`
  - `project-docs/docs/phases/phase-4/tasks/task-3-completion.md`
- Validation contract:
  - `project-docs/docs/phases/phase-4/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-4/qa/qa-plan.md`
- Related tracking:
  - `project-docs/docs/phases/phase-4/triage-log.md`
  - `project-docs/docs/current/tl_escalations_current.md`
- Commits or changed files:
  - `bce37eb7e676346cde0d134b8d557c9a0c4dc452` `feat(db): enforce issue sequence scheduling`
  - `397cb774d636c1606ab465ff6d4aa0617936708b` `feat(issues): support phase 4 issue authoring`
  - `2bfe94f24846ac09a24e94eb4f89d0fd468f9ec0` `feat(cli): wire issue create authoring flags`

## Summary

- Blocking findings: 2
- Improvement findings: 0
- Escalations raised: ESC-002
- Covered validation assertions: VAL-001, VAL-002, VAL-003, VAL-004, VAL-005

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-1, task-3 | ON TRACK | Task 1 report cites raw SQL accept/reject coverage, migration down-path coverage, and journal projection coverage at `task-1-completion.md:28`; focused rerun passed: `go test ./internal/db -run TestIssueSequenceSchedulingTriggerRejectsInvalidTaskStatesAndAllowsValidRows -count=1`. Task 3 report cites actionable CLI sequence/scheduling errors at `task-3-completion.md:18`; focused rerun passed: `go test ./cmd/fid -run 'TestIssueCreateRejectsInvalidFlagsAndConstraintViolations|TestIssueCreateUsesOneJournalEventPerCommandAndRebuildsRepresentativeHierarchy' -count=1`. |
| VAL-002 | task-2, task-3 | ON TRACK | Task 2 report cites internal issue authoring tests for epic/task/subtask defaults and invalid hierarchy at `task-2-completion.md:21`. Task 3 report cites CLI create/get smoke for epic/task/subtask and unattached task at `task-3-completion.md:19`. Focused rerun passed: `go test ./internal/issues -run 'TestCreateSupportsEpicTaskAndSubtaskAuthoring|TestCreateRejectsInvalidHierarchyAndUnknownLabels|TestCreatePreservesAuthoringFieldsAcrossRebuild' -count=1`. |
| VAL-003 | task-2, task-3 | NOT ON TRACK | Most authoring fields are evidenced, including priority, sequence, `spec_md_path`, `est_tokens`, labels, inbox, owner role, owner actor, unknown-label rejection, and rebuild. Due date is only checked by direct SQLite queries in Task 2 and Task 3, while `internal/issues.Issue` and `fid issue get --format json` do not expose `due_date`; see Blocking findings and ESC-002. |
| VAL-004 | task-2, task-3 | ON TRACK | Task 2 report cites epic metadata row creation and non-epic scheduling rejection at `task-2-completion.md:30`. Task 3 report cites CLI epic metadata proof and non-epic/missing-metadata errors at `task-3-completion.md:21`. The focused `internal/issues` and `cmd/fid` reruns passed. |
| VAL-005 | task-2, task-3 | ON TRACK | Task 2 report cites one logical `issue.created` event and delete-DB rebuild parity at `task-2-completion.md:34`. Task 3 report cites one journal line per CLI create plus `fid rebuild` parity at `task-3-completion.md:22`. Focused `cmd/fid` rerun passed. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 1 required gate | `gofmt -l cmd internal` | PASS | Completion report records exact command as PASS at `task-1-completion.md:58`. |
| Task 1 required gate | `golangci-lint run ./...` | PASS | Completion report records exact command as PASS at `task-1-completion.md:59`. |
| Task 1 required gate | `go test -race -coverprofile=coverage.out ./...` | PASS | Completion report records exact command as PASS at `task-1-completion.md:60`. |
| Task 1 required gate | `go tool cover -func=coverage.out` | PASS | Completion report records total coverage `75.2%` at `task-1-completion.md:61`. |
| Task 2 required gate | `gofmt -l cmd internal` | PASS | Completion report records exact command as PASS at `task-2-completion.md:57`. |
| Task 2 required gate | `golangci-lint run ./...` | PASS | Completion report records exact command as PASS at `task-2-completion.md:58`. |
| Task 2 required gate | `go test -race -coverprofile=coverage.out ./...` | PASS | Completion report records exact command as PASS and COV highlights at `task-2-completion.md:59`. |
| Task 2 required gate | `go tool cover -func=coverage.out` | PASS | Completion report records total coverage `76.3%` at `task-2-completion.md:65`. |
| Task 3 required gate | `gofmt -l cmd internal` | PASS | Completion report records exact command as PASS at `task-3-completion.md:34`. |
| Task 3 required gate | `golangci-lint run ./...` | PASS | Completion report records exact command as PASS at `task-3-completion.md:35`. |
| Task 3 required gate | `go test -race -coverprofile=coverage.out ./...` | PASS | Completion report records exact command as PASS and COV highlights at `task-3-completion.md:36`. |
| Task 3 required gate | `go tool cover -func=coverage.out` | PASS | Completion report records total coverage `77.0%` at `task-3-completion.md:37`. |
| Task 3 required build | `go build -o /tmp/fid ./cmd/fid` | PASS | Completion report records required CLI build as PASS at `task-3-completion.md:38`. |
| Commit verification | `git rev-parse --verify <hash>^{commit}` for all three reported hashes | PASS | All three short hashes resolve to commits: `bce37eb7...`, `397cb774...`, `2bfe94f2...`. No placeholder commit fields found. |
| Scope check | `git show --stat --oneline --no-renames bce37eb 397cb77 2bfe94f` | PASS | Changed files align to sequence scheduling, `internal/issues` authoring, CLI create, tests, migrations, and completion reports. No Task 4/5 implementation scope was found in these commit stats. |
| Representative Task 1 check | `go test ./internal/db -run TestIssueSequenceSchedulingTriggerRejectsInvalidTaskStatesAndAllowsValidRows -count=1` | PASS | Output: `ok github.com/mbell-nzl/fid/internal/db 0.020s`. |
| Representative Task 2 check | `go test ./internal/issues -run 'TestCreateSupportsEpicTaskAndSubtaskAuthoring|TestCreateRejectsInvalidHierarchyAndUnknownLabels|TestCreatePreservesAuthoringFieldsAcrossRebuild' -count=1` | PASS | Output: `ok github.com/mbell-nzl/fid/internal/issues 0.090s`. |
| Representative Task 3 check | `go test ./cmd/fid -run 'TestIssueCreateRejectsInvalidFlagsAndConstraintViolations|TestIssueCreateUsesOneJournalEventPerCommandAndRebuildsRepresentativeHierarchy' -count=1` | PASS | Output: `ok github.com/mbell-nzl/fid/cmd/fid 0.162s`. |

## Per-Task Findings

### Task 1 - Sequence Scheduling Backstop

Design plan sections checked: `05_DESIGN_PLAN.md` section 3.2 for issue hierarchy context and section 5.4 for actionable error-message expectations. Task 1's primary design reference is `07_DESIGN_PLAN_FLOWS.md` section 1.1; no `05_DESIGN_PLAN.md` parity matrix row directly governs the database-only backstop.

#### Blocking

- None.

#### Improvement

- None.

Observation: Task 1 evidence tracks VAL-001 for the persistence window. The report covers the full raw SQL matrix, projection behavior, down migration, exact required commands, and a verified commit hash.

### Task 2 - Issue Authoring Domain Service

Design plan sections checked: `05_DESIGN_PLAN.md` section 3.2 for hierarchy, section 5.1/5.2 for stable structured read shapes, and section 8 for cross-surface parity concepts relevant to read-back fields.

#### Blocking

- `internal/issues/issue.go:73` `Issue` has no `DueDate` field, and `internal/issues/issue.go:321` selects from `issues_enriched` without `due_date`. Task 2 reports VAL-003 PASS for due-date persistence at `task-2-completion.md:25`, but its proof is a direct SQLite query rather than service read-back. VAL-003 requires supported populated fields to survive journal append, projection, read-back, and rebuild. Fix direction: expose and test `due_date` through the central `internal/issues.Get` read model, or record a Tech Lead-approved VAL/spec narrowing. Escalated as ESC-002.

#### Improvement

- None.

Observation: Outside the due-date gap, Task 2 provides strong VAL-002/003/004/005 evidence: internal package tests cover valid and invalid hierarchy, unknown-label rejection, epic metadata, one-event writes, and rebuild parity.

### Task 3 - `fid issue create` CLI Surface

Design plan sections checked: `05_DESIGN_PLAN.md` section 3.2 for hierarchy, section 5.3 for stdout/stderr channel separation, section 5.4 for error wording, section 5.6 for write confirmations, and section 8 for cross-surface parity.

#### Blocking

- `cmd/fid/issue_test.go:216` defines the CLI JSON read-back shape without `due_date`, and `cmd/fid/issue_test.go:311` verifies due date by direct SQLite query after `fid issue create --due-date`. Task 3 reports VAL-003 PASS at `task-3-completion.md:20`, but the QA plan and task acceptance require CLI JSON read-back evidence for supported authoring fields. Fix direction: add `due_date` to the issue get JSON/read model and assert it through `fid issue get --format json`, or record an accepted carry-forward/narrowing before treating VAL-003 as passed. Escalated as ESC-002.

#### Improvement

- None.

Observation: Task 3 otherwise records the exact required command gate including `go build`, verifies unknown-label rejection without auto-create, checks brief stderr write confirmation, and proves one journal line per CLI create plus rebuild parity.

## Cross-cutting observations

- The Phase 4 triage log contains TRG-001 for Task 1's post-commit completion-report state and disposes it as RESUME; this review found no remaining Task 1 acceptance blocker from that note.
- Completion reports for Tasks 1-3 have concrete commit fields and all reported hashes verify. No placeholder commit fields were found.
- Required command coverage is complete for Tasks 1-3. The blocking issue is not a missing command; it is an evidence/contract gap for VAL-003 read-back of one supported field.

## Escalations

- ESC-002 (QA Reviewer triggers 5 and 6): Phase 4 due-date create field is persisted but not exposed through `internal/issues.Get` or `fid issue get --format json`, so VAL-003 read-back evidence is incomplete across Tasks 2-3.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-4/qa/review-3.md
Blocking: 2  Improvement: 0  Escalations: ESC-002
VERDICT: BLOCK
```
