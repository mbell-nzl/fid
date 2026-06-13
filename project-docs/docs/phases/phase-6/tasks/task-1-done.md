## Summary of what was implemented

- Changed `fid issue update` confirmation output to keep the `updated <id>:` header on its own line and render each changed field on a separate following line.
- Preserved update behavior, changed-field ordering, and the existing journal/activity audit path.
- Updated the CLI design examples to match the new one-field-per-line rendering.

## Files created or modified

- Modified `cmd/fid/issue.go`.
- Modified `cmd/fid/issue_update_test.go`.
- Modified `project-docs/specs/06_DESIGN_PLAN_CLI.md`.
- Created `project-docs/docs/phases/phase-6/tasks/task-1-done.md`.

## Documentation updated

- Updated `project-docs/specs/06_DESIGN_PLAN_CLI.md` so the `fid issue update` confirmation example and the v0.2 inbox-promote example both use one changed field per line.

## Validation Contract Coverage

- `VAL-009` PASS via `cmd/fid/issue_update_test.go`.
- Multi-field evidence: `updated <id>:` header line, then `title "Implement issue update" -> "Ship issue update"` and `priority 2 -> 1` on separate lines.
- Single-field evidence: `updated <id>:` header line, then exactly one `priority 2 -> 1` field line.
- Regression evidence: the required full `go test -race -coverprofile=coverage.out ./...` gate stayed green, so update behavior and audit coverage remained intact.

## Tests run

- `go test ./cmd/fid -run 'TestIssueUpdate(ChangesFieldsAndPrintsConfirmation|SingleFieldPrintsExactlyOneChangedFieldLine)$'` during red/green TDD.
- `gofmt -l cmd internal`
- `golangci-lint run ./...`
- `go test -race -coverprofile=coverage.out ./...`
- `go tool cover -func=coverage.out`

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
  `cmd/fid` coverage: `79.0%`
  `internal/issues` coverage: `82.8%`
  `internal/journal` coverage: `80.3%`
  `internal/policy` coverage: `85.0%`
  `internal/workflow` coverage: `95.5%`
- `PASS` `go tool cover -func=coverage.out`
  Total coverage: `77.7%`

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
