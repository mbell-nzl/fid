# Phase 20 - Task 7: Issue hierarchy invariant - Completion Report

## Summary of what was implemented

- Added domain-level `HierarchyInvalidError` with exact `kind`, `parent_id`, `parent_kind`, and `rule` values for illegal existing-parent and missing-parent-shape cases.
- `Create` now resolves every supplied parent before hierarchy evaluation: a missing ID returns typed `NotFoundError`, while an existing parent of an illegal kind returns typed hierarchy-invalid failure before Event preparation or append.
- Preserved parent assignment as a create-only operation. Update and transition requests expose no `ParentID`, and their Cobra leaves register neither `--epic` nor `--parent`.

## Files created or modified

- `internal/issues/create.go` — resolves supplied parents, distinguishes missing parents, and builds typed hierarchy-invalid errors.
- `internal/issues/issue.go` — adds `HierarchyInvalidError`.
- `internal/issues/create_test.go` — adds legal/illegal/missing hierarchy, typed-detail, and zero-append coverage.
- `cmd/fid/issue_parent_surface_test.go` — verifies update and transition have no parent mutation flags.
- `project-docs/docs/phases/phase-20/tasks/task-7-done.md` — completion report.

## Documentation updated

- `project-docs/docs/phases/phase-20/tasks/task-7-done.md`

## Validation Contract Coverage

- `VAL-020` — PASS. `TestCreateAllowsLegalHierarchy` at `internal/issues/create_test.go:9` proves Epic-without-parent, root Task, Task-under-Epic, and Subtask-under-Task creation. `TestCreateRejectsIllegalHierarchyBeforeAppend` at `internal/issues/create_test.go:38` covers Epic-with-parent, Subtask-under-Epic, Subtask-without-parent, Task-under-Subtask, and a missing parent; lines 151-170 assert the exact typed error fields/`NotFoundError.IssueID` and unchanged journal line count for every rejection. `TestIssueMutationRequestsCannotChangeParent` at `internal/issues/create_test.go:176` and `TestIssueUpdateAndTransitionDoNotRegisterParentMutationFlags` at `cmd/fid/issue_parent_surface_test.go:9` prove the domain request and Cobra mutation surfaces contain no parent input. The built CLI also rejected a Subtask under a real Epic while the script verified unchanged journal line count: `error: issues: subtask parent_id must reference a task`.

## Tests run

- RED: `go test ./internal/issues -run '^TestCreateRejectsIllegalHierarchyBeforeAppend$' -count=1` — expected FAIL before the typed error existed — `FAIL` (terminated `exit_code: 1`); diagnostic: `undefined: HierarchyInvalidError`.
- Focused implementation check: same command after adding the typed error — expected FAIL — `FAIL` (terminated `exit_code: 1`); diagnostic: `Create() error = *errors.errorString, want *HierarchyInvalidError`, identifying the remaining untyped Epic option guard.
- GREEN: same command after routing the Epic case through parent validation — PASS — `ok  \tgithub.com/mbell-nzl/fid/internal/issues\t0.073s` (terminated `exit_code: 0`).
- `go test ./internal/issues ./cmd/fid -run '^(TestCreateAllowsLegalHierarchy|TestCreateRejectsIllegalHierarchyBeforeAppend|TestIssueMutationRequestsCannotChangeParent|TestIssueUpdateAndTransitionDoNotRegisterParentMutationFlags)$' -count=1` — PASS — `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.046s` (terminated `exit_code: 0`).
- `go test ./internal/issues ./cmd/fid -count=1` — PASS — `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t5.509s` (terminated `exit_code: 0`).
- Real entrypoint: from a fresh temporary project working directory, ran `/tmp/fid init --prefix FID`, created an Epic, then attempted `/tmp/fid issue create --kind subtask --title 'Illegal child' --parent <epic-id> --by developer`; the script required the create to fail and compared pre/post journal line counts — PASS — `error: issues: subtask parent_id must reference a task` (terminated `exit_code: 0` for the assertion script).

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS — `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` — PASS — `0 issues.` (terminated `exit_code: 0`).
- `go test -race -coverprofile=coverage.out ./...` — PASS — `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 80.5% of statements` (terminated `exit_code: 0`).
- `go tool cover -func=coverage.out` — PASS — `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t77.7%` (terminated `exit_code: 0`).
- `go test -tags integration -race ./cmd/fid/integration/...` — PASS — `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t4.767s` (terminated `exit_code: 0`).
- `go build -o /tmp/fid ./cmd/fid` — PASS — `exit 0 (no output)` (terminated `exit_code: 0`).

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
