## Summary of what was implemented

- Centralized the enriched issue detail read model in `internal/issues`, including `due_date`, parent title, labels, acceptance criteria with latest verification per role, concern/triage summaries, links, completions, and claims.
- Added `fid issue get --section` for the Phase 4 structured sections `acceptance`, `concerns`, `completions`, `links`, and `claims`, with stable empty-array responses and actionable unknown-section errors.
- Removed SQL-backed detail assembly from `cmd/fid/issue.go`; markdown issue reads now consume the same `internal/issues.Get` shape as JSON issue reads.
- Closed ESC-002 by proving due-date read-back through both `internal/issues.Get` and `fid issue get --format json`.

## Files created or modified

- `cmd/fid/issue.go`
- `cmd/fid/issue_section_test.go`
- `internal/format/markdown.go`
- `internal/format/markdown_test.go`
- `internal/issues/detail.go`
- `internal/issues/detail_test.go`
- `internal/issues/issue.go`
- `project-docs/docs/phases/phase-4/tasks/task-4-completion.md`
- `testdata/format/issue-empty.golden.toon`

## Documentation updated

- `project-docs/docs/phases/phase-4/tasks/task-4-completion.md`

## Validation Contract Coverage

- `VAL-006` PASS: `TestGetReturnsCentralizedDetailShape` proves the central `internal/issues.Get` read model returns populated and empty detail states for labels, acceptance criteria, latest verification per role, open concern count, latest triage status, links, completions, claims, parent title, and stable empty arrays. `TestIssueGetReturnsMarkdownForEmptyAndPopulatedIssues` and `TestWriteMarkdownMatchesGoldenFixtureForPopulatedIssue` prove markdown reads the same centralized shape. Adapter-boundary evidence: `rg -n 'SELECT|QueryRow|Query\(' cmd/fid/issue.go` returned no matches, while `cmd/fid/issue.go` now routes section reads via `issues.SectionValue`.
- `VAL-007` PASS: `TestIssueGetSectionsAndDueDateReadBack` proves `fid issue get --section acceptance|concerns|completions|links|claims --format json` returns only the requested section, returns `[]` for empty sections, and rejects `--section unknown` with `supported sections: acceptance, concerns, completions, links, claims`.
- `VAL-003` PASS: due-date read-back closeout is evidenced through both `internal/issues.Get` and `fid issue get --format json`. `TestGetReturnsCentralizedDetailShape` proves a populated `due_date` reads back as `2026-06-30` and an omitted due date reads back as `nil` through `internal/issues.Get`. `TestIssueGetSectionsAndDueDateReadBack` proves `fid issue create --due-date 2026-06-30` followed by `fid issue get --format json` returns `"due_date":"2026-06-30"`, and an omitted due date reads back as `null`.

## Tests run

- `go test ./internal/issues -run TestGetReturnsCentralizedDetailShape -count=1`
- `go test ./cmd/fid -run TestIssueGetSectionsAndDueDateReadBack -count=1`
- `go test ./internal/format -run TestWriteMarkdown -count=1`
- `go test ./cmd/fid -run 'TestIssueGetReturnsMarkdownForEmptyAndPopulatedIssues|TestIssueCreateThenGetReturnsJSONByDefaultAndWithFormatFlag|TestIssueCreateSupportsHierarchyAndAuthoringFlags' -count=1`
- `go test ./internal/issues -run 'TestCreateSupportsEpicTaskAndSubtaskAuthoring|TestGetMapsPopulatedNullableFields' -count=1`
- `go test ./internal/format -run TestWriteTOONMatchesGoldenFixtureForEmptyIssue -count=1`
- `gofmt -l cmd internal`
- `golangci-lint run ./...`
- `go test -race -coverprofile=coverage.out ./...`
- `go tool cover -func=coverage.out`
- `go build -o /tmp/fid ./cmd/fid`

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal` - no output.
- `PASS` `golangci-lint run ./...` - `0 issues.`
- `PASS` `go test -race -coverprofile=coverage.out ./...` - `cmd/fid` `79.7%`; `internal/issues` `81.7%`; `internal/journal` `80.7%`; `internal/policy` `85.0%`.
- `PASS` `go tool cover -func=coverage.out` - total statements `76.8%`; `cmd/fid/issue.go newIssueGetCommand 84.6%`; `internal/issues/detail.go issueConcerns 82.4%`; `internal/issues/detail.go issueCompletions 83.3%`; `internal/issues/detail.go issueLinks 80.0%`; `internal/issues/detail.go issueClaims 83.3%`.
- `PASS` `go build -o /tmp/fid ./cmd/fid` - exited 0.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `c62c6aa`

## Commit subject

- `feat(issue): centralize detail reads`
