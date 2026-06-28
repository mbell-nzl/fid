## Summary of what was implemented

Implemented Task 2 only: a new `fid issue label` command group with
`fid issue label add` and `fid issue label remove`, backed by new
`internal/issues` services that mutate exactly one `issue_labels` row per
command and project exactly one `activity_log` row per successful add/remove.

The add/remove services now validate and authorize before append, refresh the
projection under the held write lock, reuse the existing issue + label lookup
helpers, and reject missing labels, already-present adds, and absent removals
without appending journal lines. The in-code permissions matrix now includes
the `issue.label.add` and `issue.label.remove` rows from schema section 13, and
`fid policy show` exposes those same rows.

## Files created or modified

- `cmd/fid/issue.go`
- `cmd/fid/issue_label.go`
- `cmd/fid/issue_label_test.go`
- `internal/issues/label.go`
- `internal/issues/label_test.go`
- `internal/issues/policy_actions.go`
- `internal/policy/authorize_test.go`
- `internal/policy/matrix.go`
- `internal/policy/matrix_test.go`
- `project-docs/docs/phases/phase-13/decisions.md`
- `project-docs/docs/phases/phase-13/tasks/task-2-done.md`

## Documentation updated

- Updated `project-docs/docs/phases/phase-13/decisions.md` with DEC-033 to
  record PI-2 and PI-3: reject-before-append semantics, the `activity_log` row
  shape (`event` `issue_label_added` / `issue_label_removed`; `details`
  carrying `label` + `direction`), and the `cmd/fid/issue_label.go` file
  placement.
- Reviewed `project-docs/specs/06_DESIGN_PLAN_CLI.md`; no change required.
  Task 2 matches the documented command spellings and confirmation output.
- Created this completion report at
  `project-docs/docs/phases/phase-13/tasks/task-2-done.md`.

## Validation Contract Coverage

- `VAL-001` complete: the in-code matrix now contains exactly one
  `issue.label.add` row and one `issue.label.remove` row, both matching
  `03_TRD_SCHEMA.md` section 13, and `fid policy show` exposes those rows.
  Evidence: `TestMatrixContainsExactlyTheImplementedV01WriteActions`,
  `TestMatrixIssueLabelRowsMatchSchemaSection13`,
  `TestAuthorizeIssueLabelActionsMatchMatrixExpectations`, and
  `TestIssueLabelCommandsPolicyShowIncludesMatrixRows`.
- `VAL-005`: `fid issue label add` adds exactly one label without replacing the
  rest, emits one event, and rejects nonexistent/already-present labels before
  append. Evidence:
  `TestAddLabelAddsOneRowWithoutReplacingOthersAndProjectsActivity`,
  `TestAddLabelRejectsMissingOrExistingLabelBeforeAppend`,
  `TestIssueLabelAddPreservesExistingLabelsAndProjectsActivity`, and
  `TestIssueLabelCommandsRejectBeforeAppendAndHonorPermissions/add rejects missing and duplicate label`.
- `VAL-006`: `fid issue label remove` removes exactly one label without
  disturbing the rest, emits one event, preserves the `labels` table row, and
  rejects absent-label removal before append. Evidence:
  `TestRemoveLabelRemovesOneRowPreservesLabelTableAndProjectsActivity`,
  `TestRemoveLabelRejectsAbsentLabelBeforeAppend`,
  `TestIssueLabelRemovePreservesRemainingLabelsAndProjectsActivity`, and
  `TestIssueLabelRemoveRejectsAbsentLabelBeforeAppend`.
- `VAL-007`: each successful add/remove projects exactly one `activity_log`
  row with the issue, label, and direction recorded. Task 1 already proved
  `fid label add` writes no `activity_log` row; Task 2 adds the issue-label
  side of the assertion. Evidence:
  `TestAddLabelAddsOneRowWithoutReplacingOthersAndProjectsActivity`,
  `TestRemoveLabelRemovesOneRowPreservesLabelTableAndProjectsActivity`,
  `TestIssueLabelAddPreservesExistingLabelsAndProjectsActivity`, and
  `TestIssueLabelRemovePreservesRemainingLabelsAndProjectsActivity`.
- `VAL-008`: `issue.label.add` and `issue.label.remove` authorize only
  `developer`, `tl`, `workflow_runner`, and `human`; denied and unknown roles
  reject before append naming role + action. Evidence:
  `TestAuthorizeIssueLabelActionsMatchMatrixExpectations` and
  `TestIssueLabelCommandsRejectBeforeAppendAndHonorPermissions/permissions reject denied and unknown roles`.

## Tests run

- `go test ./cmd/fid -run TestIssueLabelAddPreservesExistingLabelsAndProjectsActivity -count=1`
- `go test ./cmd/fid ./internal/issues ./internal/policy -count=1`
- `gofmt -w cmd/fid/issue.go cmd/fid/issue_label.go cmd/fid/issue_label_test.go internal/issues/label.go internal/issues/label_test.go internal/issues/policy_actions.go internal/policy/authorize_test.go internal/policy/matrix.go internal/policy/matrix_test.go`
- `go test ./cmd/fid ./internal/issues ./internal/policy -count=1`
- `gofmt -l cmd internal`
- `golangci-lint run ./...`
- `go test -race -coverprofile=coverage.out ./...`
- `go tool cover -func=coverage.out`
- `go build -o /tmp/fid ./cmd/fid`

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
- `PASS` `go tool cover -func=coverage.out`
  `cmd/fid` coverage: `78.1%`
  `internal/issues` coverage: `81.3%`
  `internal/policy` coverage: `84.6%`
  overall coverage: `76.1%`
- `PASS` `go build -o /tmp/fid ./cmd/fid`
- `N/A` `go test -tags integration -race ./...`
  Reason: no integration fixture or integration-tagged test file was touched in
  Task 2; the task spec allows this command to be skipped in that case.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `f62b6093d11e47073af8a1336cd34af6ae6da475` (`f62b609`)

## Commit subject

- `feat(issue): add issue label commands`
