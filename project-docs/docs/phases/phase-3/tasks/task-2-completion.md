# Phase 3 - Task 2: Markdown Renderer For `fid issue get` - Completion Report

## Summary of what was implemented

- Added a typed markdown renderer seam in `internal/format` for issue-shaped output, separate from the generic JSON/TOON pipeline.
- Wired `fid issue get --format markdown` to build a markdown-specific view and render the Phase 3 §21 layout, including populated Acceptance Criteria, Labels, Links, and footer metadata.
- Added deterministic markdown golden and empty-state tests in `internal/format`, plus CLI markdown coverage for populated and empty issues and cumulative format-dispatch coverage.

## Files created or modified

- `cmd/fid/issue.go`
- `cmd/fid/issue_test.go`
- `internal/format/format.go`
- `internal/format/format_test.go`
- `internal/format/markdown.go`
- `internal/format/markdown_test.go`
- `testdata/format/issue-populated.golden.md`
- `project-docs/docs/phases/phase-3/tasks/task-2-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-3/tasks/task-2-completion.md`

## Validation Contract Coverage

- `VAL-003` - PASS. `internal/format/markdown_test.go` locks the populated markdown rendering to `testdata/format/issue-populated.golden.md` and proves the empty-state issue omits empty Acceptance Criteria / Labels / Links sections without panicking. `cmd/fid/issue_test.go` exercises `fid issue get --format markdown` for both a populated seeded issue and an empty seeded issue, matching the required heading, metadata lines, checkbox mapping, labels, links, and footer layout.
- `VAL-005` - PASS for Task 2 scope. `internal/format/format_test.go` verifies markdown now dispatches successfully while YAML remains explicitly unsupported. `cmd/fid/issue_test.go` proves the CLI renders `json`, `toon`, and `markdown` through their distinct paths, and still rejects unknown `xml` with an actionable known-formats error instead of silently falling back.

## Tests run

- `go test ./internal/format -run TestWriteMarkdownMatchesGoldenFixtureForPopulatedIssue` - FAIL (expected red: markdown renderer types/seam not implemented yet)
- `go test ./internal/format -run TestWriteMarkdownMatchesGoldenFixtureForPopulatedIssue` - PASS
- `go test ./cmd/fid -run TestIssueGetReturnsMarkdownForEmptyAndPopulatedIssues` - FAIL (expected red: populated markdown output was missing the unverified role line)
- `go test ./cmd/fid -run TestIssueGetReturnsMarkdownForEmptyAndPopulatedIssues` - PASS
- `go test ./internal/format ./cmd/fid` - PASS

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS
- `go build -o /tmp/fid ./cmd/fid` - PASS

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- Markdown currently assembles parent-title, label, acceptance-criteria, verification, and link data in `cmd/fid/issue.go` for the markdown path only because the shared `issues.Get` read model does not yet expose all of that data. This keeps Task 2 within scope, but it leaves some read-shape duplication that should be centralised if future read commands need the same enriched issue view.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `5ff796b`

## Commit subject

- `feat(format): add markdown issue renderer`
