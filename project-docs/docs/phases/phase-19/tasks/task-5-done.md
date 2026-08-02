# Phase 19 - Task 5: `fid workflow list` in every output format - Blocker Report

## Summary of what was implemented

- Implemented the scoped Task 5 worktree changes: read-only `fid workflow
  list`, deterministic projected-workflow querying, JSON/TOON/YAML/Markdown
  rendering, `--kind` and `--active` filters, assignment shapes, empty-array
  rendering, and a committed in-test-literal Markdown golden.
- The task cannot be completed because the required full lint gate remains
  blocked by a pre-existing out-of-scope error in `internal/workflow/resolve.go`.

## Files created or modified

- `cmd/fid/workflow.go`
- `cmd/fid/workflow_list_test.go` (new)
- `internal/workflow/list.go` (new)
- `internal/workflow/list_test.go` (new)
- `internal/format/markdown.go`
- `internal/format/markdown_workflow.go` (new)
- `internal/format/markdown_workflow_test.go` (new)
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-19/tasks/task-5-done.md` (new blocker report)

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` documents list filters, ordering,
  `assigned_to` shapes, empty-array behavior, and illustrative example counts.
- `project-docs/specs/05_DESIGN_PLAN.md` §3.1 was verified: its `workflow
  list` entry matches the shipped `--kind`, `--active`, and `--format` surface;
  no correction was needed.

## Validation Contract Coverage

- VAL-018: BLOCKED - `TestWorkflowListRendersDocumentedJSONFieldSet` in
  `cmd/fid/workflow_list_test.go` and
  `TestWorkflowListFormatsFiltersAndAssignments` provide the key-set,
  projected-count, and `any`-as-null assertions, but the required full gate is
  blocked.
- VAL-019: BLOCKED - `TestWorkflowListFormatsFiltersAndAssignments` asserts
  the exact active and unfiltered `qa_review` ID sets (`alpha-any`,
  `bravo-qa`, `charlie-qa-inactive`) while excluding the seeded other-kind row;
  the required full gate is blocked.
- VAL-020: BLOCKED - `assertWorkflowListCommandShapes` in
  `cmd/fid/workflow_list_test.go` asserts the kind assignment object
  `{"scope":"kind","kind":"task"}`; the required full gate is blocked.
- VAL-060: BLOCKED - `TestWorkflowListFormatsFiltersAndAssignments` decodes
  JSON, TOON, Markdown, and YAML to one equivalent structure; the required
  full gate is blocked.
- VAL-061: BLOCKED - `TestWriteWorkflowListMarkdownMatchesCommittedLiteral`
  in `internal/format/markdown_workflow_test.go` compares exact Markdown bytes
  with a committed **in-test literal** (not a `testdata` fixture); the required
  full gate is blocked.
- VAL-062: BLOCKED - `TestWorkflowListFormatsFiltersAndAssignments` asserts
  that `--active` excludes `charlie-qa-inactive` and the unfiltered run
  includes it; the required full gate is blocked.
- VAL-063: BLOCKED - the same command test asserts the composed
  `--kind qa_review --active` intersection is exactly `alpha-any` and
  `bravo-qa`; the required full gate is blocked.
- VAL-064: BLOCKED - `assertWorkflowListCommandShapes` asserts the project
  assignment `kind` key is absent, not null; the required full gate is blocked.
- VAL-065: BLOCKED - `assertWorkflowListCommandShapes` asserts the unassigned
  `bravo-qa` value is a non-nil empty array; the required full gate is blocked.
- VAL-066: BLOCKED - `TestListFiltersCountsAssignmentsAndOrdering` asserts
  result and assignment order, while `TestWorkflowListFormatsFiltersAndAssignments`
  compares two consecutive JSON byte streams; the required full gate is blocked.
- VAL-067: BLOCKED - `TestWorkflowListFormatsFiltersAndAssignments` asserts
  `[]\n` and empty stderr for a no-match run in every format; the required full
  gate is blocked.

## Tests run

- `go test ./cmd/fid -run '^TestWorkflowListRendersDocumentedJSONFieldSet$' -count=1`
  - RED observed before registration: `FAIL`
- `go test ./cmd/fid -run '^TestWorkflowListRendersDocumentedJSONFieldSet$' -count=1`
  - GREEN: `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.024s`
- `go test ./internal/format -run '^TestWriteWorkflowListMarkdownMatchesCommittedLiteral$' -count=1`
  - RED observed before the lossless Markdown change: `FAIL`
- `go test ./internal/format -run '^TestWriteWorkflowListMarkdown' -count=1`
  - GREEN: `ok  \tgithub.com/mbell-nzl/fid/internal/format\t0.004s`
- `go test ./internal/workflow -run '^TestList' -count=1`
  - PASS: `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.049s`
- `go test ./cmd/fid -run '^TestWorkflowList' -count=1`
  - PASS: `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.095s`
- `gofmt -w cmd/fid/workflow.go cmd/fid/workflow_list_test.go internal/workflow/list.go internal/workflow/list_test.go internal/format/markdown.go internal/format/markdown_workflow.go internal/format/markdown_workflow_test.go && go test ./internal/workflow ./internal/format ./cmd/fid -run 'Test(List|WorkflowList|WriteWorkflowListMarkdown)' -count=1`
  - PASS: `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.087s`
- Real built-CLI entrypoint exercise: BLOCKED - the required build did not run
  after the lint blocker, so `/tmp/fid workflow list` was not available to
  exercise honestly.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`
- `golangci-lint run ./...` - BLOCKED - `* errcheck: 1`; after correcting the
  two in-scope `internal/workflow/list.go` findings, the byte-identical rerun
  reports the pre-existing out-of-scope `internal/workflow/resolve.go:61:18:
  Error return value of rows.Close is not checked (errcheck)`. `git show HEAD`
  verifies the same unchecked `defer rows.Close()` is already at that line.
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - not run because
  the preceding required lint command did not pass.
- `go tool cover -func=coverage.out` - BLOCKED - not run because its required
  coverage-producing test command did not run after the lint blocker.
- `go build -o /tmp/fid ./cmd/fid` - BLOCKED - not run because the preceding
  required lint command did not pass; consequently the required built-CLI
  runtime exercise was not run.

## Deviations from spec

- The requested Task 5 commit, finalized report commit, report-lint run, and
  built-CLI exercise were not performed because the mandatory lint gate is
  blocked by the verified pre-existing, out-of-scope `resolve.go` error.

## Concerns, technical debt, or notes for the Tech Lead

- Required action: fix or explicitly scope `internal/workflow/resolve.go:61`'s
  unchecked `rows.Close` error, then resume Task 5 from this preserved worktree
  and rerun the full gate from `gofmt -l cmd internal` in the specified order.

## Dependency-budget note

- No dependency changes.

## Commit hash

- No commit created - blocker report.

## Commit subject

- No commit created - blocker report.
