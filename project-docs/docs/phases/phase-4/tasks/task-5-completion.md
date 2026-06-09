## Summary of what was implemented

Implemented `fid issue list` for Phase 4 Task 5 over the `issues_enriched`
read model. The task added domain-level list querying in `internal/issues`
with filters for kind, exact status, status category, epic parent, existing
label, priority set, inbox flag, and open-concern flag; deterministic ordering
for `sequence`, `priority`, `created`, and `updated`; a thin Cobra adapter in
`cmd/fid`; and markdown list rendering through the existing format engine.

## Files created or modified

- `cmd/fid/issue.go`
- `cmd/fid/issue_list_test.go`
- `internal/issues/list.go`
- `internal/issues/list_test.go`
- `internal/format/markdown.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-4/tasks/task-5-completion.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
  Added the public markdown example for `fid issue list --format markdown`.

## Validation Contract Coverage

- `VAL-008` PASS.
  Exact evidence:
  `go test ./internal/issues -run TestList`
  `TestListFiltersByKind` proved `--kind task` includes only tasks and excludes
  epic/subtask rows.
  `TestListFiltersAndOrdering/filters` proved exact include/exclude subsets for
  `status`, `status_category`, `epic`, `label`, `priority`, `inbox`, and
  `open concerns`, plus the empty-result case (`Label: "missing"` => `[]`).
  `go test ./cmd/fid -run TestIssueListFiltersAndOrdering`
  proved representative CLI flag mapping for `--kind`, `--status`,
  `--status-category`, `--epic`, `--label`, `--priority`, `--inbox`,
  `--has-concerns`, and empty-result exit-0 JSON output.

- `VAL-009` PASS.
  Exact evidence:
  `go test ./internal/issues -run TestList`
  `TestListFiltersAndOrdering/ordering` proved exact stable ID order for
  `sequence`, `priority`, `created`, and `updated`.
  `go test ./cmd/fid -run TestIssueListDispatchesFormatsAndRejectsUnknownFormats`
  proved JSON baseline parity with TOON and YAML, an exact markdown string
  assertion for the list renderer, and actionable `--format xml` rejection with
  the known format list.
  `go test -race -coverprofile=coverage.out ./...` and
  `go tool cover -func=coverage.out`
  reported `cmd/fid` coverage at `79.5%` and `internal/issues` coverage at
  `79.5%`, satisfying the task's COV references.

## Tests run

- `go test ./internal/issues -run TestList`
- `go test ./cmd/fid -run TestIssueListFiltersAndOrdering`
- `go test ./cmd/fid -run TestIssueListDispatchesFormatsAndRejectsUnknownFormats`
- `go test ./internal/issues ./internal/format ./cmd/fid`

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS
- `golangci-lint run ./...` — PASS
- `go test -race -coverprofile=coverage.out ./...` — PASS
- `go tool cover -func=coverage.out` — PASS
- `go build -o /tmp/fid ./cmd/fid` — PASS

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit.

## Commit subject

- Pending commit.
