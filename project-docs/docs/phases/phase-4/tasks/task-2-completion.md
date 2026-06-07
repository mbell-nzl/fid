## Summary of what was implemented

- Expanded `internal/issues.Create` beyond task-only authoring so it now builds single-event journal writes for `epic`, `task`, and `subtask` creation.
- Added domain validation for hierarchy, epic scheduling usage, owner-role overrides, blank/unknown labels, and core create-option normalization before append/project.
- Added Phase 4 defaults for workflow IDs, initial `todo` status, and owner-role defaults (`tl` for epics, `developer` for tasks/subtasks).
- Added epic metadata projection for epics, issue-label assignment effects for existing labels, and label read-back in `internal/issues.Get`.
- Added package-level authoring tests that prove multi-kind round trip behavior, invalid hierarchy rejection, optional field persistence, and delete-DB rebuild parity.

## Files created or modified

- `internal/issues/issue.go`
- `internal/issues/issue_test.go`
- `project-docs/docs/phases/phase-4/tasks/task-2-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-4/tasks/task-2-completion.md`

## Validation Contract Coverage

- `VAL-002` PASS.
  Exact evidence:
  `internal/issues/issue_test.go:171` `TestCreateSupportsEpicTaskAndSubtaskAuthoring` proves epic/task/subtask creation succeeds with generated IDs, legal hierarchy, `epic-default` / `task-default` / `subtask-default`, initial `todo`, `unstarted` status category on read-back, and correct owner defaults (`tl` on epic, `developer` on subtask).
  `internal/issues/issue_test.go:342` `TestCreateRejectsInvalidHierarchyAndUnknownLabels` proves invalid hierarchy states are rejected with explicit constraint errors.
- `VAL-003` PASS.
  Exact evidence:
  `internal/issues/issue_test.go:171` `TestCreateSupportsEpicTaskAndSubtaskAuthoring` proves priority, sequence, `spec_md_path`, `est_tokens`, `due_date` (via `SELECT date(due_date)`), existing-label assignment, inbox flag, owner-role override, and owner actor persist through create and read-back.
  `internal/issues/issue_test.go:475` `TestCreatePreservesAuthoringFieldsAcrossRebuild` proves those persisted fields survive delete-DB rebuild and read back identically after `journal.RebuildWorkspace`.
  `internal/issues/issue_test.go:342` `TestCreateRejectsInvalidHierarchyAndUnknownLabels` proves unknown labels are rejected and no auto-create path exists.
- `VAL-004` PASS.
  Exact evidence:
  `internal/issues/issue_test.go:171` `TestCreateSupportsEpicTaskAndSubtaskAuthoring` queries `epic_metadata` and proves epic creation writes exactly one row with `scheduling = 'sequential'`.
  `internal/issues/issue_test.go:342` `TestCreateRejectsInvalidHierarchyAndUnknownLabels` proves non-epic scheduling metadata is rejected by domain validation with `create issue: scheduling is only valid for epics`.
- `VAL-005` PASS.
  Exact evidence:
  `internal/issues/issue_test.go:171` `TestCreateSupportsEpicTaskAndSubtaskAuthoring` reads the journal file and proves each create call appends one logical `issue.created` event line even when the event contains multiple effects.
  `internal/issues/issue_test.go:475` `TestCreatePreservesAuthoringFieldsAcrossRebuild` proves delete-DB + rebuild reconstructs the same issue state and epic metadata.
  `go test -race -coverprofile=coverage.out ./...` PASS with `internal/issues` coverage `86.5%`, `internal/journal` coverage `80.7%`, `internal/policy` coverage `85.0%`, and `cmd/fid` coverage `74.3%`, satisfying the Phase 4 COV carry-in thresholds cited by the QA plan.

## Tests run

- `go test ./internal/issues -run 'TestCreateSupportsEpicTaskAndSubtaskAuthoring|TestCreateRejectsInvalidHierarchyAndUnknownLabels|TestCreatePreservesAuthoringFieldsAcrossRebuild|TestCreateRejectsUnsupportedKind'`
  Red: build failed before implementation because `CreateOptions` did not yet expose the Phase 4 authoring fields.
- `go test ./internal/issues -run 'TestCreateSupportsEpicTaskAndSubtaskAuthoring|TestCreateRejectsInvalidHierarchyAndUnknownLabels|TestCreatePreservesAuthoringFieldsAcrossRebuild|TestCreateRejectsUnsupportedKind'`
  Green: PASS after implementing Phase 4 authoring behavior.
- `gofmt -l cmd internal`
  PASS.
- `golangci-lint run ./...`
  PASS (`0 issues.`).
- `go test -race -coverprofile=coverage.out ./...`
  PASS.
- `go tool cover -func=coverage.out`
  PASS.

## Exact quality-gate results

- `gofmt -l cmd internal`: PASS.
- `golangci-lint run ./...`: PASS (`0 issues.`).
- `go test -race -coverprofile=coverage.out ./...`: PASS.
  Coverage highlights from command output:
  `cmd/fid` `74.3%`
  `internal/issues` `86.5%`
  `internal/journal` `80.7%`
  `internal/policy` `85.0%`
- `go tool cover -func=coverage.out`: PASS (`total: (statements) 76.3%`).

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit creation at report-write time; copied from `git rev-parse` after the required Task 2 commit.

## Commit subject

- Pending commit creation at report-write time; copied from `git log -1 --format=%s` after the required Task 2 commit.
