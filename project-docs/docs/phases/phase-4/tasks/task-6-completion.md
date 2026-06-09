## Summary of what was implemented

Added behavior-focused `internal/issues` tests only, covering the helper
branches left uncovered after Task 5. The new tests exercise priority-filter
parsing, list-order validation, section selection and unknown-section
rejection, parent-title lookup for nil/missing/existing parents, and due-date
normalization for nil/blank/date-only/RFC3339/malformed inputs.

The coverage closeout raised `internal/issues` from the phase-end blocking
`79.5%` to `85.5%` without modifying any non-test product code. `cmd/fid`
remained at `79.5%`, so COV-002 stayed green.

## Files created or modified

- `internal/issues/list_test.go`
- `internal/issues/detail_test.go`
- `project-docs/docs/phases/phase-4/tasks/task-5-completion.md`
- `project-docs/docs/phases/phase-4/tasks/task-6-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-4/tasks/task-5-completion.md`
  Corrected the inaccurate Task 5 COV claim to state that `internal/issues`
  `79.5%` did not satisfy COV-001 and that Task 6 owns the closeout.
- `project-docs/docs/phases/phase-4/tasks/task-6-completion.md`
  Recorded the coverage increase, branch evidence, exact task-gate results, and
  diff-scope confirmation.

## Validation Contract Coverage

- No new `VAL-NNN` was advanced. This task restores the COV-001 deterministic
  gate that Phase 4 final review blocked on.
- The new tests strengthen existing evidence for `VAL-003`, `VAL-006`,
  `VAL-007`, `VAL-008`, and `VAL-009` by covering previously untested helper
  branches in `internal/issues`.
- Coverage evidence:
  Phase 4 final review recorded `internal/issues` at `79.5%` and blocked
  COV-001.
  `go test -race -coverprofile=coverage.out ./...` after Task 6 reported
  `internal/issues` `85.5%`, `cmd/fid` `79.5%`, `internal/journal` `80.7%`,
  and `internal/policy` `85.0%`.
  `go tool cover -func=coverage.out` confirmed `internal/issues >= 80%`
  (COV-001 PASS) and `cmd/fid >= 70%` (COV-002 PASS).
- Specific uncovered branches now covered:
  `ParsePriorityFilter` blank input, valid comma-separated input, embedded empty
  value rejection, and non-integer rejection.
  `issueListOrderClause` blank defaulting, supported orders, and unknown-order
  rejection through `UnsupportedListOrderError`.
  `SectionValue` supported section dispatch and unknown-section rejection
  through `UnsupportedSectionError`.
  `normalizeDateOnlyPointer` nil, blank, date-only, RFC3339, and malformed
  inputs.
  `issueParentTitle` nil parent, no-row parent, and existing-parent lookup.

## Tests run

- `go test -race -coverprofile=/tmp/issues-list.cover ./internal/issues -run 'TestListFiltersByKind|TestListFiltersAndOrdering|TestParsePriorityFilter|TestIssueListOrderClause'`
- `go tool cover -func=/tmp/issues-list.cover`
- `go test -race -coverprofile=/tmp/issues-helpers.cover ./internal/issues -run 'TestGetReturnsCentralizedDetailShape|TestSectionValue|TestIssueParentTitle|TestNormalizeDateOnlyPointer|TestListFiltersByKind|TestListFiltersAndOrdering|TestParsePriorityFilter|TestIssueListOrderClause'`
- `go tool cover -func=/tmp/issues-helpers.cover`
- `go test -race -coverprofile=/tmp/issues-full.cover ./internal/issues`
- `go tool cover -func=/tmp/issues-full.cover`
- `gofmt -l cmd internal`
- `golangci-lint run ./...`
- `go test -race -coverprofile=coverage.out ./...`
- `go tool cover -func=coverage.out`
- `git diff --stat`
- `git diff --stat --cached`

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS. No output.
- `golangci-lint run ./...` - PASS. Output: `0 issues.`
- `go test -race -coverprofile=coverage.out ./...` - PASS.
  Coverage lines relevant to the task: `cmd/fid 79.5%`, `internal/issues 85.5%`,
  `internal/journal 80.7%`, `internal/policy 85.0%`.
- `go tool cover -func=coverage.out` - PASS.
  `internal/issues/detail.go` helper coverage now includes `SupportedSections 100.0%`,
  `SectionValue 100.0%`, `issueParentTitle 87.5%`, and
  `normalizeDateOnlyPointer 100.0%`.
  `internal/issues/list.go` helper coverage now includes `SupportedListOrders 100.0%`,
  `issueListOrderClause 100.0%`, and `ParsePriorityFilter 100.0%`.
  COV-001 PASS: `internal/issues 85.5% >= 80%`.
  COV-002 PASS: `cmd/fid 79.5% >= 70%`.
- `git diff --stat` - BLOCKED for task-local proof.
  The repo-wide worktree already contained unrelated concurrent edits from
  other agents, so plain `git diff --stat` showed additional files outside Task
  6:

  ```text
   internal/issues/detail_test.go                     | 175 +++++++++++++++++++++
   internal/issues/list_test.go                       | 141 +++++++++++++++++
   .../docs/current/tl_escalations_current.md         |  23 ++-
   .../docs/phases/phase-4/qa/validation-contract.md  |   3 +-
   .../docs/phases/phase-4/tasks/task-1-completion.md |   4 +-
   .../docs/phases/phase-4/tasks/task-2-completion.md |   4 +-
   .../docs/phases/phase-4/tasks/task-3-completion.md |   4 +-
   .../docs/phases/phase-4/tasks/task-4-completion.md |   4 +-
   project-docs/docs/phases/phase-4/tasks/task-4.md   |  22 ++-
   .../docs/phases/phase-4/tasks/task-5-completion.md |   8 +-
   project-docs/docs/phases/phase-4/triage-log.md     |  29 ++++
   11 files changed, 401 insertions(+), 16 deletions(-)
  ```

- `git diff --stat --cached` - PASS.
  Staging only the four task files isolates the Task 6 commit scope and proves
  the diff is test-only plus the two named documentation files:

  ```text
   internal/issues/detail_test.go                     | 175 +++++++++++++++++++++
   internal/issues/list_test.go                       | 141 +++++++++++++++++
   .../docs/phases/phase-4/tasks/task-5-completion.md |   8 +-
   .../docs/phases/phase-4/tasks/task-6-completion.md | 135 ++++++++++++++++
   4 files changed, 456 insertions(+), 3 deletions(-)
  ```

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending post-commit verification.

## Commit subject

- Pending post-commit verification.
