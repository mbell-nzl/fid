# Phase 19 - Task 5: `fid workflow list` in every output format - Completion Report

## Summary of what was implemented

- Added the read-only `fid workflow list` command with `--kind`, `--active`,
  and JSON/TOON/Markdown/YAML output.
- Added deterministic projected-workflow listing with SQL-backed state and
  transition counts, ordered assignment summaries, and non-nil empty arrays.
- Added the dedicated Markdown renderer and a byte-exact committed in-test
  literal golden.
- Applied DEC-102's sole authorized out-of-scope repair to the unchecked
  `rows.Close()` defer in workflow resolution.

## Files created or modified

- `cmd/fid/workflow.go`
- `cmd/fid/workflow_list_test.go` (new)
- `internal/workflow/list.go` (new)
- `internal/workflow/list_test.go` (new)
- `internal/workflow/resolve.go` (DEC-102 authorized one-line repair)
- `internal/format/markdown.go`
- `internal/format/markdown_workflow.go` (new)
- `internal/format/markdown_workflow_test.go` (new)
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-19/tasks/task-5-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` documents the list filters,
  deterministic result and assignment ordering, assignment JSON shapes,
  empty-array behavior, and illustrative count examples.
- Verified `project-docs/specs/05_DESIGN_PLAN.md` §3.1 already matches the
  shipped `workflow list` flag surface; no correction was required.

## Validation Contract Coverage

- VAL-018: PASS - `TestWorkflowListRendersDocumentedJSONFieldSet`
  (`cmd/fid/workflow_list_test.go:13-40`) asserts the exact JSON key set;
  `assertWorkflowListCommandShapes` (`cmd/fid/workflow_list_test.go:153-174`)
  asserts `any` as null and the projected state/transition counts.
- VAL-060: PASS - `TestWorkflowListFormatsFiltersAndAssignments`
  (`cmd/fid/workflow_list_test.go:43-77`) decodes JSON, TOON, Markdown, and
  YAML and compares every result to the same structure.
- VAL-061: PASS - `TestWriteWorkflowListMarkdownMatchesCommittedLiteral`
  (`internal/format/markdown_workflow_test.go:10-34`) byte-compares the
  renderer output to its explicit in-test `const want` literal. The expected
  Markdown is a committed **in-test literal**, not a `testdata` fixture.
- VAL-019: PASS - `TestWorkflowListFormatsFiltersAndAssignments`
  (`cmd/fid/workflow_list_test.go:50-56,83-85`) asserts the exact qualifying
  active and unfiltered `qa_review` ID sets, excluding `delta-epic`.
- VAL-062: PASS - `TestWorkflowListFormatsFiltersAndAssignments`
  (`cmd/fid/workflow_list_test.go:79-85`) proves the inactive workflow is
  excluded with `--active` and included without it.
- VAL-063: PASS - `TestWorkflowListFormatsFiltersAndAssignments`
  (`cmd/fid/workflow_list_test.go:50-56`) asserts the exact
  `--kind qa_review --active` intersection.
- VAL-020: PASS - `assertWorkflowListCommandShapes`
  (`cmd/fid/workflow_list_test.go:160-163`) asserts the exact kind assignment
  object `{"scope":"kind","kind":"task"}`.
- VAL-064: PASS - `assertWorkflowListCommandShapes`
  (`cmd/fid/workflow_list_test.go:164-169`) asserts project scope and proves
  that the `kind` key is absent.
- VAL-065: PASS - `assertWorkflowListCommandShapes`
  (`cmd/fid/workflow_list_test.go:171-174`) asserts the unassigned value is a
  present, non-nil empty array.
- VAL-066: PASS - `TestListFiltersCountsAssignmentsAndOrdering`
  (`internal/workflow/list_test.go:26-35`) asserts workflow and assignment
  order; `TestWorkflowListFormatsFiltersAndAssignments`
  (`cmd/fid/workflow_list_test.go:79-81`) asserts consecutive JSON bytes are
  identical.
- VAL-067: PASS - `TestWorkflowListFormatsFiltersAndAssignments`
  (`cmd/fid/workflow_list_test.go:88-95`) asserts `[]\n` and empty stderr for
  a no-match result in JSON, TOON, Markdown, and YAML.

## Tests run

- `go test ./internal/workflow ./internal/format ./cmd/fid -run 'Test(List|WorkflowList|WriteWorkflowListMarkdown)' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.086s`
- Built-CLI entrypoint: in fresh `/tmp/fid-phase19-task5.siuHqS`,
  `/tmp/fid init && /tmp/fid workflow list --format markdown` - PASS - rendered
  the four seeded workflow rows, including `| task-default | Task default | task | true | 9 | 17 | [] |`.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`
- `golangci-lint run ./...` - PASS - `0 issues.` This repo-wide green run is
  the independent re-execution and verification of Task 4's formerly
  falsely-attested lint gate 2 of 6, as required by DEC-102.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t5.477s\tcoverage: 80.4% of statements`
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t77.0%`
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)`

All five PASS tokens above were observed from terminating commands in this
session. No previous blocker-session result is recorded as current PASS.

## Deviations from spec

- ESC-19-006 was initially misclassified as pre-existing because the prior
  session relied on `git show HEAD`; DEC-102 established that
  `internal/workflow/resolve.go` was introduced by Task 4 commit `32614bb`.
  DEC-102 therefore authorized the exact one-line repair at `resolve.go:61`:
  `defer rows.Close()` became `defer func() { _ = rows.Close() }()`. No other
  out-of-scope repair was made; behavior is unchanged because rows are drained
  and `rows.Err()` is checked.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `8f79013c7e8afcd3d4bd73f9821484d9dd9d9334`

## Commit subject

- `feat(workflow): add workflow list command`
