## Summary of what was implemented

Implemented Task 1 only: a new `fid label` command group with `fid label add`
and `fid label list`, a new `internal/labels` package that owns the label
add/list business rules, markdown rendering for label lists, and the in-code
`label.add` permissions row required by schema section 13.

`fid label add` now validates and authorizes before append, refreshes the
projection under the held write lock, emits exactly one `label.added` event
with one `labels` insert effect, rejects duplicate names before append, stores
optional `colour` / `description` values, and projects no `activity_log` row.
`fid label list` reads labels from the existing `labels` table and returns the
same `name` / `colour` / `description` data in json, toon, markdown, and yaml.

## Files created or modified

- `cmd/fid/label.go`
- `cmd/fid/label_test.go`
- `cmd/fid/root.go`
- `internal/format/markdown.go`
- `internal/labels/add.go`
- `internal/labels/labels_test.go`
- `internal/labels/list.go`
- `internal/labels/types.go`
- `internal/policy/matrix.go`
- `internal/policy/matrix_test.go`
- `project-docs/docs/phases/phase-13/decisions.md`
- `project-docs/docs/phases/phase-13/tasks/task-1-done.md`

## Documentation updated

- Created `project-docs/docs/phases/phase-13/decisions.md` with DEC-032 to
  record PI-1 label-add semantics: duplicate rejection before append, optional
  metadata handling, and no `activity_log` projection.
- Created this completion report at
  `project-docs/docs/phases/phase-13/tasks/task-1-done.md`.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` was reviewed and required no
  change; implementation matches section 9 for `fid label add` and
  `fid label list`.

## Validation Contract Coverage

- `VAL-001` partial: `internal/policy/matrix.go` now includes exactly one
  `label.add` row matching `03_TRD_SCHEMA.md` section 13
  (`developer`, `tl`, `workflow_runner`, `human` true; all others false).
  Evidence: `TestMatrixContainsExactlyTheImplementedV01WriteActions`,
  `TestMatrixLabelAddRowMatchesSchemaSection13`, and
  `TestLabelAddPolicyShowIncludesMatrixRow`.
- `VAL-002`: label creation succeeds with one row in one event and duplicate
  add rejects before append. Evidence:
  `TestLabelAddCreatesOneRowInOneEventWithoutActivityLog`,
  `TestAddStoresOptionalFieldsAndRejectsDuplicateBeforeAppend`, and
  `TestLabelAddRejectsDuplicateDeniedAndUnknownBeforeAppend/duplicate add`.
- `VAL-003`: label list returns the same `name` / `colour` / `description`
  data across json, toon, markdown, and yaml, preserving null optional fields.
  Evidence: `TestLabelListDefaultsToJSON`,
  `TestLabelListFormatsMatchJSONReference`, and
  `TestListReturnsLabelsInNameOrderWithNilOptionalFields`.
- `VAL-004`: `label.add` authorizes only permitted roles, and denied / unknown
  roles reject before append naming role + action. Evidence:
  `TestLabelAddRejectsDuplicateDeniedAndUnknownBeforeAppend/denied and unknown reject before append`.
- `VAL-007` partial: successful `fid label add` projects no `activity_log` row.
  Evidence: `TestLabelAddCreatesOneRowInOneEventWithoutActivityLog`.

## Tests run

- `go test ./cmd/fid -run TestLabelAddCreatesOneRowInOneEventWithoutActivityLog -count=1`
- `go test ./cmd/fid -run TestLabelListDefaultsToJSON -count=1`
- `go test ./cmd/fid -run TestLabelListFormatsMatchJSONReference -count=1`
- `go test ./internal/labels -count=1`
- `go test ./cmd/fid -run 'TestLabel(AddCreatesOneRowInOneEventWithoutActivityLog|ListDefaultsToJSON|ListFormatsMatchJSONReference|AddRejectsDuplicateDeniedAndUnknownBeforeAppend|AddPolicyShowIncludesMatrixRow)$' -count=1`
- `go test ./internal/policy -count=1`
- `go test ./internal/labels -count=1`
- `gofmt -l cmd internal`
- `golangci-lint run ./...`
- `go test -race -coverprofile=coverage.out ./...`
- `go tool cover -func=coverage.out`
- `go build -o /tmp/fid ./cmd/fid`

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
- `PASS` `go tool cover -func=coverage.out`
  `cmd/fid` coverage: `77.8%`
  `internal/claims` coverage: `81.4%`
  `internal/issues` coverage: `81.8%`
  `internal/journal` coverage: `80.0%`
  `internal/labels` coverage: `81.1%`
  `internal/policy` coverage: `81.2%`
  `internal/workflow` coverage: `95.5%`
  overall coverage: `76.0%`
- `PASS` `go build -o /tmp/fid ./cmd/fid`
- `N/A` `go test -tags integration -race ./...`
  Reason: no integration fixture or integration-tagged test file was touched in
  Task 1; the task spec allows this command to be skipped in that case.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `d5ad0ea2f5b4dbbaf6f468dd20523636d48a0d6a` (`d5ad0ea`)

## Commit subject

- `feat(labels): add label create and list commands`
