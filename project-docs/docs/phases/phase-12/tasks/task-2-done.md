## Summary of what was implemented

- Added the new `internal/links` domain package for Phase 12 Task 2: the
  closed v0.1 relation registry, inverse lookup, endpoint-type resolution,
  endpoint existence checks, relation/endpoint compatibility validation,
  duplicate detection, and the journaled `Add` write path that creates the
  declared row plus its inverse in one event.
- Added `fid link add <from> <relation> <to> --by <role>` with Task 2 scope
  only: flag wiring, role authorization, and stderr confirmation naming both
  directions.
- Appended the in-code `link.add` and `link.remove` matrix rows with the
  section-13 booleans for all seven roles and verified `fid policy show` parity.

## Files created or modified

- `cmd/fid/link.go`
- `cmd/fid/link_test.go`
- `cmd/fid/root.go`
- `internal/links/add.go`
- `internal/links/add_test.go`
- `internal/links/errors.go`
- `internal/links/registry.go`
- `internal/links/registry_test.go`
- `internal/links/types.go`
- `internal/policy/matrix.go`
- `internal/policy/matrix_test.go`
- `project-docs/docs/phases/phase-12/decisions.md`
- `project-docs/docs/phases/phase-12/tasks/task-2-done.md`

## Documentation updated

- Updated `project-docs/docs/phases/phase-12/decisions.md` with DEC-030, which
  records the PI-2 enforcement decision for the closed v0.1 link registry,
  app-layer prechecks, and the no-migration outcome.
- Added this completion report at
  `project-docs/docs/phases/phase-12/tasks/task-2-done.md`.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: verified no change required; the
  implemented `fid link add` confirmation text matches the existing spec.

## Validation Contract Coverage

- `VAL-001`: PASS. `internal/policy/matrix.go` now contains exactly one
  `link.add` row and one `link.remove` row, both true for all seven roles, and
  `fid policy show` exposes the same matrix entries that `Authorize` enforces
  (`TestMatrixContainsExactlyTheImplementedV01WriteActions`,
  `TestMatrixLinkRowsMatchSchemaSection13`,
  `TestLinkAddPolicyShowParityAndUnknownRoleRejection`).
- `VAL-005`: PASS. `fid link add <issue> blocks <issue> --by developer`
  appends one journal event containing two `links` inserts (declared +
  `blocked_by` inverse), projects both rows, and prints both directions; a
  duplicate add rejects before append (`TestLinkAddCreatesDeclaredAndInverseRowsInOneJournalEvent`,
  `TestAddSupportsInverseDeclaredRelationAndRequirementEndpoints/declared_inverse_relation`,
  `TestAddRejectsInvalidRequestsBeforeAppend/duplicate_link`,
  `TestLinkAddImplementsRequirementAndRejectsInvalidInputBeforeAppend/rejections_keep_journal_unchanged/duplicate_add`).
- `VAL-006`: PASS. An issue -> requirement `implements` link creates the
  `implemented_by` inverse, and an unsynced requirement endpoint rejects before
  append (`TestAddSupportsInverseDeclaredRelationAndRequirementEndpoints/issue_implements_requirement`,
  `TestLinkAddImplementsRequirementAndRejectsInvalidInputBeforeAppend/implements_requirement`,
  `TestAddRejectsInvalidRequestsBeforeAppend/missing_requirement_endpoint`,
  `TestLinkAddImplementsRequirementAndRejectsInvalidInputBeforeAppend/rejections_keep_journal_unchanged/missing_requirement_endpoint`).
- `VAL-007`: PASS. Relation/endpoint-incompatible combinations and missing
  endpoints reject before append with non-zero exit and no journal growth
  (`TestAddRejectsInvalidRequestsBeforeAppend/incompatible_relation`,
  `TestLinkAddImplementsRequirementAndRejectsInvalidInputBeforeAppend/rejections_keep_journal_unchanged/implements_issue_to_issue`,
  `.../blocks_issue_to_requirement`,
  `.../missing_requirement_endpoint`).
- `VAL-008`: PASS. The relation registry equals the documented seven relations
  plus inverses exactly, and unknown relations plus self-links reject before
  append (`TestRelationInverseMapMatchesDocumentedV01Set`,
  `TestRegistryHelpersCoverEndpointAndCompatibilityRules`,
  `TestAddRejectsInvalidRequestsBeforeAppend/unknown_relation`,
  `TestAddRejectsInvalidRequestsBeforeAppend/self_link`,
  `TestLinkAddImplementsRequirementAndRejectsInvalidInputBeforeAppend/rejections_keep_journal_unchanged/unknown_relation`,
  `.../self_link`).

## Tests run

- `go test ./cmd/fid -run TestLinkAddCreatesDeclaredAndInverseRowsInOneJournalEvent` — red: failed with `unknown command "link" for "fid"`.
- `go test ./internal/policy -run 'TestMatrixContainsExactlyTheImplementedV01WriteActions|TestMatrixLinkRowsMatchSchemaSection13'` — red: failed with missing `link.add`/`link.remove` rows.
- `go test ./internal/policy -run 'TestMatrixContainsExactlyTheImplementedV01WriteActions|TestMatrixLinkRowsMatchSchemaSection13'` — green.
- `go test ./cmd/fid -run TestLinkAddCreatesDeclaredAndInverseRowsInOneJournalEvent` — green.
- `go test ./internal/links -run 'TestRelationInverseMapMatchesDocumentedV01Set|TestAddSupportsInverseDeclaredRelationAndRequirementEndpoints|TestAddRejectsInvalidRequestsBeforeAppend'` — green.
- `go test ./cmd/fid -run 'TestLinkAddCreatesDeclaredAndInverseRowsInOneJournalEvent|TestLinkAddImplementsRequirementAndRejectsInvalidInputBeforeAppend|TestLinkAddPolicyShowParityAndUnknownRoleRejection'` — green.

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS
- `golangci-lint run ./...` — PASS (`0 issues.`)
- `go test -race -coverprofile=coverage.out ./...` — PASS
- `go tool cover -func=coverage.out` — PASS
  Coverage excerpts: total 76.6%; `cmd/fid` 78.2%; `internal/links` 81.6%;
  `internal/issues` 81.8%; `internal/claims` 81.1%; `internal/policy` 81.2%;
  `internal/journal` 80.0%; `internal/requirements` 80.4%;
  `internal/workflow` 95.5%.
- `go build -o /tmp/fid ./cmd/fid` — PASS
- `go test -tags integration -race ./...` — N/A (no integration fixtures were
  touched in Task 2)

## Deviations from spec

- The committed report was finalized before the task commit existed. Per the
  one-commit / no-amend constraint, the exact hash and subject were added here
  after commit creation without amending the commit.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `0abe93edc2d98f1adda4cd5b1bec4e15f64ae06d`

## Commit subject

- `feat(links): add link add command`
