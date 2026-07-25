# Phase 18 - Mid-Phase Review 2 (Task 4)

Reviewer: qa-reviewer
Date: 2026-07-19
Tasks covered: phase-18/task-4

## Verdict

BLOCK

Task 5 must not start. The Phase Runner must route this localized Blocking
finding through TL-Lite before resuming the phase.

## Scope Reviewed

- Task specification:
  `project-docs/docs/phases/phase-18/tasks/task-4.md`
- Completion report:
  `project-docs/docs/phases/phase-18/tasks/task-4-done.md`
- Prior review and remediation review:
  `project-docs/docs/phases/phase-18/qa/review-1.md`,
  `project-docs/docs/phases/phase-18/qa/review-1-remediation.md`
- Phase decisions: `project-docs/docs/phases/phase-18/decisions.md`
  (DEC-067 and DEC-068)
- Validation contract and refreshed QA plan:
  `project-docs/docs/phases/phase-18/qa/validation-contract.md`,
  `project-docs/docs/phases/phase-18/qa/qa-plan.md`
- Implementation commit:
  `0b4b8aaa14498d43b23381bd4f83669468c15e10`
  (`feat(db): enforce issue workflow foreign keys`)
- Report-finalization commit:
  `b79aa9d702fdc532f6cb46e6b6fe7c82fa6fb817`
  (`docs(task): finalize completion report for phase-18 task-4`)

## Summary

- Blocking findings: 1
- Improvement findings: 0
- Escalations raised: none
- Covered validation assertions: VAL-006, VAL-007, VAL-018, VAL-019

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-006 | task-4 | NOT ON TRACK | The independent non-cached migration test passed, and it proves unknown-workflow and invalid-known-pair INSERT rejection, a valid INSERT, two invalid UPDATE rejections, both declared FKs, and a clean `foreign_key_check`. It does not exercise a valid-pair UPDATE, required by QA-09's INSERT/UPDATE valid-and-invalid matrix. See the Blocking finding. |
| VAL-007 | task-4 | ON TRACK | `TestMigration0008PreservesIssuesAndDependentObjectsThroughDownUp` migrates a version-7 epic/task/subtask fixture, deep-compares ordered quoted snapshots of every `issues` column before/after 0008 and after down/up, and passed independently. |
| VAL-018 | task-4 | ON TRACK | The same independently passed test compares the complete `sqlite_schema` fingerprint (except the intentionally rebuilt table), validates the five recreated `ON issues` triggers and view, two retained named indexes, all five retained kind triggers, all five retained ref-integrity triggers, both structural autoindexes, and a stable down/up fingerprint. |
| VAL-019 | task-4 | ON TRACK | `TestMigration0008PreservesIssueGetAndListRenderedBytes` creates a version-7 fixture, captures JSON rendered through the shared get/list entry paths, migrates to 0008, and requires byte-identical get and list output; it passed in the independent scoped run. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 4 report hygiene | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-18/tasks/task-4-done.md project-docs/docs/phases/phase-18/tasks/task-4.md` | PASS | All required headings and all five verbatim task commands are recorded `PASS`; the task hash resolves, subject matches, and the report is committed and clean. |
| Exact command coverage | Task 4 spec and QA plan §5 comparison | PASS | The report records each required command exactly: `gofmt -l cmd internal`, `golangci-lint run ./...`, `go test -race -coverprofile=coverage.out ./...`, `go tool cover -func=coverage.out`, and `go test -tags integration -race ./...`. QA plan §5 adds no additional Task 4 close command. |
| Commit metadata and scope | `git rev-parse --verify 0b4b8aaa14498d43b23381bd4f83669468c15e10^{commit}`; `git show --check`; `git diff-tree --name-status -r` | PASS | Hash resolves; subject and `Refs: phase-18/task-4` footer match. Whitespace check passes. The implementation diff is confined to migration 0008, migration/FK/parity tests, test-only maintenance fixtures, and the draft report; the final report commit has the matching footer and only finalizes its metadata. |
| Reproduced migration/FK/object/read-parity scope | `go test -race -count=1 ./internal/db -run '^TestMigration0008(EnforcesWorkflowAndStateForeignKeys|PreservesIssuesAndDependentObjectsThroughDownUp|PreservesIssueGetAndListRenderedBytes)$'` | PASS with evidence gap | Non-cached result: `ok github.com/mbell-nzl/fid/internal/db 1.907s`. It confirms row preservation, invalid FK INSERT/UPDATE rejection, valid INSERT, clean `foreign_key_check`, named-object fingerprint/down-up stability, retained-trigger probes, and get/list byte parity. It exposes the missing valid UPDATE qualifying state recorded below. |
| Reproduced exact required command | `go test -tags integration -race ./...` | PASS | Exact Task 4 command exited 0; all listed packages passed. |
| Migration DDL and dependent-object audit | `migrations/0008_v02a_issues_workflow_fk.up.sql` and focused test assertions | PASS | The migration replaces the workflow/status enum checks with `workflow_id -> workflows(id)` and `(workflow_id,status) -> workflow_states(workflow_id,state)` FKs, preserves the non-workflow checks/UNIQUE, recreates the view and five `ON issues` triggers, and leaves surviving dependents in place. |

## Per-Task Findings

### Task 4 - Migrate issues.workflow_id / status to real FKs (table rebuild)

#### Blocking

- [internal/db/issues_workflow_fk_test.go:61] VAL-006 and QA-09 require qualifying and non-qualifying INSERT **and UPDATE** FK evidence. The test accepts a valid custom pair only through `execIssueInsert` at lines 61-69; its only two `UPDATE issues` probes at lines 71-78 are invalid and require rejection. There is no successful UPDATE to a valid `(workflow_id, status)` pair in the Task 4 committed scope, so the qualifying UPDATE state is unproven despite the report's overall VAL-006 PASS. Fix direction: add a focused valid-pair UPDATE probe (and assert the persisted updated pair) alongside the existing invalid workflow and invalid known-pair UPDATE probes, correct the report evidence if needed, rerun the exact Task 4 gate, and route the result through TL-Lite before Task 5.

#### Improvement

- None.

#### Observation

- The table rebuild otherwise follows the DEC-067 exact-inventory correction rather than the superseded generic-index wording: the test distinguishes recreated `ON issues` objects from retained cross-table objects, uses structural PRAGMA checks for generated autoindexes, and verifies the retained trigger families functionally. The version-targeted rendered-read fixture also avoids the prohibited auto-migrating-only “before” comparison.

## Escalations

- None. The missing qualifying UPDATE test is a direct, local Task 4 acceptance/evidence repair. It has no cross-task, architecture, product, security, or validation-contract ambiguity requiring an `ESC-NNN`; Phase Runner triage through TL-Lite is still required because the review verdict is BLOCK.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-18/qa/review-2.md
Blocking: 1  Improvement: 0  Escalations: none
VERDICT: BLOCK
```
