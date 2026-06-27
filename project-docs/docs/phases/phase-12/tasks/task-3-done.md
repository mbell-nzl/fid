## Summary of what was implemented

- Added `fid link remove <from> <relation> <to> --by <role>` on the existing
  `link` command tree, keeping business rules in `internal/links`.
- Added the `internal/links` remove write path so one journal event deletes the
  declared link row, deletes its inverse row, and projects one
  `activity_log` removal row keyed to the declared `from` entity.
- Reused the Task 2 relation registry and endpoint-type helpers for remove-time
  validation, including before-append rejection of unknown roles, invalid
  relation/type combinations, and missing links.

## Files created or modified

- `cmd/fid/link.go`
- `cmd/fid/link_test.go`
- `internal/links/add.go`
- `internal/links/errors.go`
- `internal/links/registry_test.go`
- `internal/links/remove.go`
- `internal/links/remove_test.go`
- `internal/links/types.go`
- `project-docs/docs/phases/phase-12/tasks/task-3-done.md`

## Documentation updated

- Added this completion report at
  `project-docs/docs/phases/phase-12/tasks/task-3-done.md`.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: verified no change required; the
  implemented `fid link remove` stderr confirmation matches the existing spec.

## Validation Contract Coverage

- `VAL-009`: PASS. `fid link remove` now removes the declared row and its
  inverse in one event, projects a `link_removed` activity row, prints both
  removed directions, rejects a missing link before append, and rejects an
  unknown role before append. Evidence:
  `TestRemoveDeletesDeclaredAndInverseRowsAndProjectsActivity`,
  `TestRemoveRejectsMissingLinkBeforeAppend`,
  `TestRemoveSupportsRequirementInverse`,
  `TestLinkRemoveDeletesDeclaredAndInverseRowsInOneJournalEvent`,
  `TestLinkRemoveRejectsInvalidInputBeforeAppend`.
- Coverage gate evidence for Task 3 scope: `go tool cover -func=coverage.out`
  reported total coverage `76.6%`, `cmd/fid` `78.3%`, and `internal/links`
  `79.5%`.

## Tests run

- `go test ./internal/links -run 'TestRemoveDeletesDeclaredAndInverseRowsAndProjectsActivity|TestRemoveRejectsMissingLinkBeforeAppend'` — red: build failed with `undefined: Remove`, `undefined: RemoveOptions`, and `undefined: MissingLinkError`.
- `go test ./cmd/fid -run 'TestLinkRemoveDeletesDeclaredAndInverseRowsInOneJournalEvent|TestLinkRemoveRejectsInvalidInputBeforeAppend'` — red: failed with `unknown flag: --by` because `link remove` was not wired.
- `go test ./internal/links -run 'TestRemoveDeletesDeclaredAndInverseRowsAndProjectsActivity|TestRemoveRejectsMissingLinkBeforeAppend|TestRemoveSupportsRequirementInverse|TestLinkErrorsRenderActionableMessages'` — green.
- `go test ./cmd/fid -run 'TestLinkRemoveDeletesDeclaredAndInverseRowsInOneJournalEvent|TestLinkRemoveRejectsInvalidInputBeforeAppend'` — green.
- `go test ./cmd/fid -run 'TestLink(AddCreatesDeclaredAndInverseRowsInOneJournalEvent|AddImplementsRequirementAndRejectsInvalidInputBeforeAppend|AddPolicyShowParityAndUnknownRoleRejection|RemoveDeletesDeclaredAndInverseRowsInOneJournalEvent|RemoveRejectsInvalidInputBeforeAppend)'` — green.

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS
- `golangci-lint run ./...` — PASS (`0 issues.`)
- `go test -race -coverprofile=coverage.out ./...` — PASS
- `go tool cover -func=coverage.out` — PASS
  Coverage excerpts: total `76.6%`; `cmd/fid` `78.3%`; `internal/links`
  `79.5%`; `internal/issues` `81.8%`; `internal/claims` `81.1%`;
  `internal/journal` `80.0%`; `internal/policy` `81.2%`;
  `internal/requirements` `80.4%`.
- `go build -o /tmp/fid ./cmd/fid` — PASS
- `go test -tags integration -race ./...` — N/A (no integration fixtures were
  touched in Task 3)

## Deviations from spec

- The report is written before the required task commit exists. Per the
  one-commit / no-amend constraint, the commit hash and subject are filled in
  after commit creation without amending the commit.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending final task commit.

## Commit subject

- Pending final task commit.
