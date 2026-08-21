# Phase 20 - Task 8: Issue read contract - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: Fresh final QA re-review.
- Human attention required: no

## Summary of what was implemented

- Migrated `issue get`, `issue list`, and `ready` to the root-`--json` `fid-cli.v1` Response Envelope with `event:null` for successful reads.
- Added canonical issue read result types with exact nullable fields, empty-array handling, revision loading, closed ready result/explanation states, retained filter validation, and read collection-limit enforcement.
- Removed the leaf `--format` flags and the Markdown-oriented `issue get --section` surface; migrated command and integration coverage to the canonical responses.
- Read commands now use the resolved project root from command context rather than resolving a leaf-local working directory.

## Files created or modified

- `cmd/fid/issue_get_list.go`, `cmd/fid/ready.go`, and their command/integration tests.
- `cmd/fid/issue_get_list_test.go` - shared envelope decoder and exact closed-schema coverage.
- `internal/issues/read.go`, `internal/issues/get.go`, `internal/issues/issue.go`, `internal/issues/list.go`, `internal/issues/ready.go`, and removed section helpers/tests.
- `internal/db/issues_workflow_fk_cli_test.go` - current-schema setup for the canonical issue read.

## Documentation updated

- `project-docs/docs/phases/phase-20/tasks/task-8-done.md` - completion evidence for Task 8.

## Validation Contract Coverage

- VAL-008: PASS - `TestIssueGetHumanSuccessUsesStdout` at `cmd/fid/issue_test.go:21` and `TestReadyHumanSuccessUsesStdout` at `cmd/fid/ready_test.go:12` assert successful human reads write only to stdout; `TestIssueReadRejectsRetiredAndMalformedInputs` at `cmd/fid/issue_get_list_test.go:163` rejects retired read flags without printing `.fid/` advice.
- VAL-013: PASS - `TestIssueGetJSONMatchesClosedSchema` at `cmd/fid/issue_get_list_test.go:88` decodes the actual root-`--json issue get` envelope, exact issue keys, required nullable keys, empty arrays, and `event:null`.
- VAL-014: PASS - `TestIssueListJSONMatchesClosedSchemaAndRetainedFilters` at `cmd/fid/issue_get_list_test.go:115` verifies the canonical list envelope, exact summary keys, retained kind/status-category/order filters, and result count; `TestIssueListFiltersAndOrdering` at `cmd/fid/issue_list_test.go:11` covers retained ordering/filter behavior.
- VAL-015: PASS - `TestReadyJSONHasExclusiveClosedResultStates` at `cmd/fid/issue_get_list_test.go:144` proves the no-ready `task:null`/closed explanation state, while `TestReadyJSONSequentialBlockerUsesClosedExplanation` at `cmd/fid/issue_get_list_test.go:194` proves the sequential blocker code and its populated nullable fields; `TestReadyEpicScopesSelectionAndBlockerStates` at `cmd/fid/ready_test.go:26` covers selection and blocker scopes.
- VAL-016: PASS - `TestIssueReadRejectsRetiredAndMalformedInputs` at `cmd/fid/issue_get_list_test.go:163`, `TestIssueGetRejectsRetiredFormatFlags` at `cmd/fid/issue_test.go:591`, and `TestV01ReadCommandsRejectRetiredFormats` at `cmd/fid/e2e_v01_format_matrix_test.go:5` reject removed `--section`/`--format` and malformed or unlisted inputs.
- VAL-017 (read half): PASS - shared `decodeReadEnvelope` at `cmd/fid/issue_get_list_test.go:26` asserts `event:null` for each actual get, list, and ready envelope exercised by the three closed-schema tests.
- VAL-044 (read rows): PASS - `TestIssueReadRejectsRetiredAndMalformedInputs` at `cmd/fid/issue_get_list_test.go:163` exercises malformed/unlisted read arguments and requires empty stdout on rejection.
- VAL-045 (read ordering/limit): PASS - `TestReadCollectionLimitRejectsOnlyOverLimit` at `cmd/fid/issue_get_list_test.go:185` accepts 10,000 and rejects 10,001 without truncation; `TestIssueListFiltersAndOrdering` at `cmd/fid/issue_list_test.go:11` covers list ordering.

## Tests run

- `go test ./cmd/fid -run '^TestIssueGetJSONMatchesClosedSchema$' -count=1` - RED - `FAIL`
- `go test ./cmd/fid -run '^TestIssueGetJSONMatchesClosedSchema$' -count=1` - UNVERIFIED - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.097s` (historical record lacked observed termination)
- `go test ./cmd/fid -run '^(TestIssueGetJSONMatchesClosedSchema|TestIssueListJSONMatchesClosedSchemaAndRetainedFilters|TestReadyJSONHasExclusiveClosedResultStates)$' -count=1` - UNVERIFIED - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.124s` (historical record lacked observed termination)
- `go test ./cmd/fid -run '^(TestIssueReadRejectsRetiredAndMalformedInputs|TestReadCollectionLimitRejectsOnlyOverLimit)$' -count=1` - UNVERIFIED - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.050s` (historical record lacked observed termination)
- `go test ./cmd/fid -run '^TestReadyJSONSequentialBlockerUsesClosedExplanation$' -count=1` - UNVERIFIED - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.052s` (historical record lacked observed termination)
- `mktemp workspace; /tmp/fid init --prefix FID; /tmp/fid issue create; /tmp/fid --json issue get <created-id>` - UNVERIFIED - `{"contract_version":"fid-cli.v1","ok":true,"operation":"issue.get",...,"event":null}` (historical record lacked observed termination)

## Exact quality-gate results

- `gofmt -l cmd internal` - UNVERIFIED - `exit 0 (no output)` (historical record lacked observed termination)
- `golangci-lint run ./...` - UNVERIFIED - `0 issues.` (historical record lacked observed termination)
- `go test -race -coverprofile=coverage.out ./...` - UNVERIFIED - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 80.5% of statements` (historical record lacked observed termination)
- `go tool cover -func=coverage.out` - UNVERIFIED - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t77.5%` (historical record lacked observed termination)
- `go test -tags integration -race ./cmd/fid/integration/...` - UNVERIFIED - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t4.971s` (historical record lacked observed termination)
- `go build -o /tmp/fid ./cmd/fid` - UNVERIFIED - `exit 0 (no output)` (historical record lacked observed termination)

Fresh direct-recovery evidence captured on 2026-08-21:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t78.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t4.915s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `4df4d66`

## Commit subject

- `feat(issue): migrate canonical issue read responses`
