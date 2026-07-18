# Phase 17 - Mid-Phase Review 2 (Tasks 3-5)

Reviewer: qa-reviewer
Date: 2026-07-15
Tasks covered: phase-17/task-3 through phase-17/task-5

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-17/tasks/task-3.md`
  - `project-docs/docs/phases/phase-17/tasks/task-4.md`
  - `project-docs/docs/phases/phase-17/tasks/task-5.md`
- Completion reports:
  - `project-docs/docs/phases/phase-17/tasks/task-3-done.md`
  - `project-docs/docs/phases/phase-17/tasks/task-4-done.md`
  - `project-docs/docs/phases/phase-17/tasks/task-5-done.md`
- Validation contract: `project-docs/docs/phases/phase-17/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-17/qa/qa-plan.md`
- Commits or changed files:
  - `d2a087d858fe12a1b33d9380e03e4497d9194350` — `feat(maintenance): audit structural projection errors`
  - `4257bcec983a0ea664c54d453af8d8e68e8f026c` — `feat(maintenance): add audit warnings and info`
  - `0dec198f67dae3814af2a8811606e4945fda3fc9` — `feat(maintenance): add check command formats`

## Summary

- Blocking findings: 0
- Improvement findings: 0
- Escalations raised: none
- Covered validation assertions: VAL-003, VAL-004, VAL-005, VAL-006, VAL-007, VAL-008, VAL-009, VAL-010, VAL-011, VAL-012, VAL-013, VAL-014

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-003 | task-3 | ON TRACK | `TestCheckErrorsReportsAllSeededStructuralViolationFamilies` passed uncached; it asserts the complete ordered structural-check list, actionable identity, repeat equality, and no journal/DB write. `TestCheckErrorsExcludesEveryResolvedV01Rule` passed with no errors. |
| VAL-004 | task-3 | ON TRACK | The same structural test includes orphaned and resolved link, claim, activity-reference, and journal-event states; the valid fixture contains all supported activity reference types. |
| VAL-005 | task-3 | ON TRACK | The structural test asserts two distinct `foreign_key_integrity` findings with table/row/parent identity; the valid fixture returns no errors. Trigger/FK bypass helpers are confined to `internal/maintenance/test_helpers_test.go`. |
| VAL-006 | task-3 | ON TRACK | Invalid workflow and state pairs are asserted as findings, while `seedAllWorkflowStates` supplies the valid default-workflow control; the implementation delegates classification to `internal/workflow.CategoryForState`. |
| VAL-007 | task-4 | ON TRACK | Injected-time warning test includes a live claim and a claim at TTL expiry, and passed uncached. |
| VAL-008 | task-4 | ON TRACK | The warning fixture asserts live-claim exclusion, before/at/after grace, created-at fallback, non-status update preservation, and the separate zero-grace test passed. |
| VAL-009 | task-4 | ON TRACK | The warning fixture contains superseded, current, and NULL completion-reference controls; only the superseded verification is expected. |
| VAL-010 | task-4 | ON TRACK | The warning fixture distinguishes missing, issue-only, and completion-linked agent-run states; only missing and issue-only cases remain telemetry gaps. |
| VAL-011 | task-4 | ON TRACK | `TestCheckInfoCountsAndResultAreDeterministic` passed uncached with exact counts and a deeply equal second audit; findings are sorted by check/entity/message. |
| VAL-012 | task-5 | ON TRACK | Focused CLI execution passed healthy and failing subtests for JSON, TOON, Markdown, and YAML. Each failing state retained structured stdout and reported `fail`. |
| VAL-013 | task-5 | ON TRACK | Focused CLI execution passed the default warning (`warn`, zero) and promoted warning (`fail`, non-zero) control, with decoded warning payloads asserted identical and errors empty. |
| VAL-014 | task-5 | ON TRACK | The uncached integration test passed both `missing_projection_rebuilds` and `stale_projection_catches_up`; it compares journal bytes before/after and verifies the caught-up projection row. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 3 report evidence | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-3-done.md project-docs/docs/phases/phase-17/tasks/task-3.md` | PASS | All four exact task commands are recorded PASS; committed report, hash, and subject verified. |
| Task 4 report evidence | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-4-done.md project-docs/docs/phases/phase-17/tasks/task-4.md` | PASS | All four exact task commands are recorded PASS; committed report, hash, and subject verified. |
| Task 5 report evidence | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-5-done.md project-docs/docs/phases/phase-17/tasks/task-5.md` | PASS | All six exact task commands, including integration and build, are recorded PASS; committed report, hash, and subject verified. |
| Required integration gate | `go test -tags integration -race ./...` | PASS | Reproduced exact gate; all packages passed, including `cmd/fid` and `internal/maintenance`. |
| Structural and warning evidence | `go test -race ./internal/maintenance -count=1 -v` | PASS | Reproduced uncached: structural invalid/valid controls, deterministic info result, warning qualifying/non-qualifying controls, and zero-grace control all passed. |
| Format and non-zero stdout evidence | `go test -race ./cmd/fid -run '^TestMaintenanceCheckCommand(RendersHealthyAndFailingResults|RendersAllFormats|PromotesWarningsWithoutRewritingThem)$' -count=1 -v` | PASS | Reproduced healthy/failing JSON, TOON, Markdown, and YAML subtests and warning-promotion payload preservation. |
| Freshness and journal immutability evidence | `go test -tags integration -race ./cmd/fid -run '^TestMaintenanceCheckRefreshesProjectionWithoutChangingJournal$' -count=1 -v` | PASS | Reproduced both missing-projection rebuild and stale-projection catch-up subtests; each compares journal bytes before and after. |
| Test-only corruption isolation | `internal/maintenance/test_helpers_test.go:92-207`; `cmd/fid/maintenance_test.go:35-45,108-118` | PASS | Trigger dropping and relaxed SQLite constraint/FK pragmas occur only in `_test.go` temp-fixture setup. The production maintenance path contains no bypass and reads through normal projection refresh/audit code. |
| Structured non-zero behavior | `cmd/fid/maintenance.go:66-75`; `cmd/fid/maintenance_test.go:47-69,120-147,172-199` | PASS | The command renders the selected structured view before returning the non-zero audit sentinel; tests decode stdout and require empty stderr for structural and promoted-warning outcomes. |
| Documentation and format parity | `project-docs/specs/05_DESIGN_PLAN.md:344-383`; `project-docs/specs/06_DESIGN_PLAN_CLI.md:1438-1466`; `project-docs/docs/playbook.md:117-133` | PASS | CLI syntax, four formats, warning promotion, stdout, and exit semantics match the implementation and focused tests. |
| Commit verification | `git rev-parse --verify <reported-hash>^{commit}` and `git show --no-patch --format=fuller <reported-hash>` for Tasks 3-5 | PASS | All three reported hashes resolve to commits with the reported subjects and `Refs: phase-17/task-3`, `task-4`, or `task-5` respectively. |
| Entry blocker / phase-end ownership | `qa/review-entry.md:28-37,76-79`; `qa-plan.md:8-14,107-108,132,299-309`; `tl_escalations_current.md:24-29` | PASS | QA-BLK-01 is closed by the accepted Entry Review and 81.0% focused/full journal coverage evidence; QA-BLK-02 remains closed by committed model-audit evidence. VAL-024 remains a terminal phase-end outcome owned by the QA Lead. |

## Per-Task Findings

### Task 3 - Maintenance structural errors

#### Blocking

- None.

#### Improvement

- None.

#### Observation

- Test-only corruption paths are isolated to temporary test databases (`internal/maintenance/test_helpers_test.go:42-70,92-207`); the paired valid fixture and read-only assertions (`internal/maintenance/check_test.go:58-70,93-104`) provide both required states without weakening production trigger defenses.

### Task 4 - Maintenance warnings and info

#### Blocking

- None.

#### Improvement

- None.

#### Observation

- Warning classification uses injected time in the tests (`internal/maintenance/warnings_test.go:9-68`) and compares precise qualifying/non-qualifying identities, while the implementation reads the latest append-ordered `status_changed` timestamp (`internal/maintenance/warnings.go:95-141`) and retains deterministic sorting (`internal/maintenance/check.go:86-100`).

### Task 5 - `fid maintenance check` CLI + formats

#### Blocking

- None.

#### Improvement

- None.

#### Observation

- The command renders all structured output before its expected non-zero audit return (`cmd/fid/maintenance.go:66-75`), and the four-format CLI test exercises both healthy and structural-error states (`cmd/fid/maintenance_test.go:72-148`). The integration test separately proves missing and stale projections refresh without changing journal bytes (`cmd/fid/maintenance_integration_test.go:13-86`).

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-17/qa/review-2.md
Blocking: 0  Improvement: 0  Escalations: none
VERDICT: CONTINUE
```
