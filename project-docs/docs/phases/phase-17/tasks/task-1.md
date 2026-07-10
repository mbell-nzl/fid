# Phase 17 - Task 1: Close CF-019 by Splitting the Markdown Formatter

## Objective

Split the oversized markdown renderer into focused package-local files before
Phase 17 adds another renderer. Preserve every existing public output exactly
and close the formatter file-size debt without changing format dispatch.

## Validation Contract Coverage

- VAL-001: all pre-Phase-17 markdown behavior remains byte-identical and every
  affected non-test file is at or below 400 lines.

## Binding Constraints

- This is behavior-preserving extraction only. Do not rename headings, reorder
  rows, alter empty states, or broaden generic markdown support.
- Keep the package API private and unchanged; task-5 owns maintenance output.
- CF-019 remains open until Tech Lead wrap accepts VAL-001 evidence.

## Scope

### Files To Create Or Modify

- `internal/format/markdown.go` - retain dispatch and shared helpers only.
- `internal/format/markdown_issue.go` - issue detail/list rendering.
- `internal/format/markdown_lists.go` - AC, concern, claim, gate, label, link,
  and policy list rendering.
- `internal/format/markdown_stats.go` - ready and stats rendering.
- `internal/format/markdown_test.go` and existing focused formatter tests - add
  characterization coverage only where a moved branch is currently uncovered.

### Files To Read For Context

- `internal/format/format.go`
- `project-docs/docs/current/decisions_active.md` (CF-019)
- `project-docs/specs/08_QUALITY_PLAN.md` section 4.4

### Out Of Scope

- New maintenance result rendering.
- Generic reflection-based markdown rendering.
- Any JSON, TOON, YAML, command, or domain behavior change.

## Acceptance Criteria

- [ ] Existing markdown tests prove byte-for-byte output equality for every
  pre-Phase-17 dispatched type. (VAL-001)
- [ ] `markdown.go`, `markdown_issue.go`, `markdown_lists.go`, and
  `markdown_stats.go` are each at or below 400 lines. (VAL-001)
- [ ] No exported API or non-markdown format behavior changes. (VAL-001)

## First Failing Test

Before moving code, add a characterization test for the first currently
uncovered dispatch branch; if all branches are covered, record that the existing
suite was the pre-refactor green baseline and use the first move-induced failure
as red evidence.

## Documentation Updates

- None expected beyond `task-1-done.md`; Tech Lead closes CF-019 at phase wrap.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
wc -l internal/format/markdown.go internal/format/markdown_issue.go internal/format/markdown_lists.go internal/format/markdown_stats.go
```

Required-command rules live in `AGENTS.md`. Every command is reported `PASS`,
or `BLOCKED` with reason; the full gate is required before task close.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Task execution is blocked until kickoff approval promotes VAL-001 to Ready.
- Preserve helper ownership inside `internal/format`; do not create a shared
  utility package for this split.
- Completion report: `project-docs/docs/phases/phase-17/tasks/task-1-done.md`.
