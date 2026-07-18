# Phase 17 - Task 1: Close CF-019 by Splitting the Markdown Formatter - Completion Report

## Summary of what was implemented

- Split the private markdown formatter into focused issue, list, and ready/stats files while retaining dispatch and the shared priority helper in `markdown.go`.
- Preserved every pre-Phase-17 markdown rendering branch and added byte-for-byte characterization coverage for the previously uncovered pointer `ReadyView` dispatch branch.

## Files created or modified

- `internal/format/markdown.go`
- `internal/format/markdown_issue.go`
- `internal/format/markdown_lists.go`
- `internal/format/markdown_stats.go`
- `internal/format/markdown_test.go`
- `project-docs/docs/phases/phase-17/tasks/task-1-done.md`

## Documentation updated

- Added this required Task 1 completion report only. No product, policy, compatibility, or specification documentation changed.

## Validation Contract Coverage

- **VAL-001 — PASS:** Existing golden and CLI Markdown-output tests passed unchanged in the full race suite. The added `ReadyView` pointer characterization asserts exact output with all optional display fields. The required line-count command confirms `markdown.go` (88), `markdown_issue.go` (164), `markdown_lists.go` (321), and `markdown_stats.go` (121) are all at or below 400 lines.

## Tests run

- `go test -race ./internal/format -run '^TestWriteMarkdownRendersReadyViewPointer$' -count=1` — PASS; pre-split characterization baseline for the previously uncovered pointer `ReadyView` dispatch branch.
- `go test -race ./internal/format -run '^(TestWriteMarkdownMatchesGoldenFixtureForPopulatedIssue|TestWriteMarkdownOmitsEmptySectionsForEmptyIssue|TestWriteMarkdownRendersIssueGateLists|TestWriteMarkdownRendersReadyViewPointer)$' -count=1` — PASS; existing exact issue/gate output checks and the new ready-view output check after extraction.
- `go build -o /tmp/fid ./cmd/fid && /tmp/fid --help` — PASS; the built `fid` CLI displayed the unchanged root command tree.

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS; no output.
- `golangci-lint run ./...` — PASS; `0 issues.`
- `go test -race -coverprofile=coverage.out ./...` — PASS; all packages passed with race instrumentation and produced `coverage.out`.
- `go tool cover -func=coverage.out` — PASS; coverage report emitted successfully; total statement coverage was 76.1%.
- `wc -l internal/format/markdown.go internal/format/markdown_issue.go internal/format/markdown_lists.go internal/format/markdown_stats.go` — PASS; 88, 164, 321, and 121 lines respectively.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- The Phase 17 QA-plan banner still says Task 1 is blocked pending entry review, although the committed Task 0 report and committed `review-entry.md` satisfy every stated prerequisite and the review records `VERDICT: CONTINUE`. Refreshing that stale phase-status text is outside this task's allowed documentation scope.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `65d2f94587310668b252ebfd4141ed9ef26c8ec0`

## Commit subject

- `refactor(format): split markdown renderer`
