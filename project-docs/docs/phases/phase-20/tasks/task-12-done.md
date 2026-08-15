# Phase 20 - Task 12: Workflow command contract - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Migrated `workflow define`, `workflow assign`, and `workflow list` to closed `fid-cli.v1` envelopes, stdout Human confirmations, resolved-root selection, and canonical Event results.
- Removed canonical `--require` and `--format`; constrained workflow applicability, kinds, lists, edges, and assignment scope; retained deterministic kind-over-project selection and replacement.
- Added typed workflow not-found and duplicate paths without SQLite-driver message classification.

## Files created or modified

- `cmd/fid/workflow.go`; workflow command tests; tagged workflow integration scenario.
- `internal/workflow/{assign,assignment,authoring,define,list,errors}.go` and tests.
- `internal/cli/codes.go` and contract-registry test; `cmd/fid/link_test.go` obsolete Markdown helper removal.

## Documentation updated

- None required.

## Validation Contract Coverage

- VAL-008: PASS - `TestWorkflowDefineHumanSuccessUsesStdoutOnly` in `cmd/fid/workflow_define_test.go` proves Human success is stdout-only; the real `/tmp/fid ... --json workflow define` interaction returned stdout only.
- VAL-017: PASS - `TestWorkflowDefineJSONMatchesContract` asserts a non-null Event and compares its ID with the appended journal Event; `TestWorkflowListJSONMatchesContractAndUsesNullAssignmentScope` asserts read `event:null`.
- VAL-026: PASS - `TestWorkflowDefineJSONMatchesContract`, `TestWorkflowDefineRejectsRemovedAndInvalidRequestsBeforeAppend`, `TestWorkflowAssignJSONPrecedenceAndScopedReplacement`, and `TestWorkflowListJSONMatchesContractAndUsesNullAssignmentScope` prove closed schemas, removed rules, atomic invalid rejection, replacement, precedence, and null project scope.
- VAL-041: PASS - `TestWorkflowDefineJSONMatchesContract` proves supplied attribution is recorded on the workflow Event with no requirement/protected effects.
- VAL-044: PASS - define/assign/list invalid-request tests cover removed `--require`, legacy applicability, malformed values, noncanonical arrows, bounds, list kind, and exact scope selection before append.
- VAL-045: PASS - `TestWorkflowListJSONMatchesContractAndUsesNullAssignmentScope` asserts the closed collection/count shape; `internal/workflow/list.go` orders summaries by `name, id`.

## Tests run

- `go test ./cmd/fid --run '^TestWorkflowDefineJSONMatchesContract$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); first red test exposed unavailable root JSON transport.
- `go test ./cmd/fid --run '^TestWorkflowDefineJSONMatchesContract$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  0.051s` (terminated `exit_code: 0`)
- `/tmp/fid --project /tmp/fid-task12-entrypoint.bGyTm7 --json workflow define --name entry-flow --applies-to task --states todo,done,cancelled --categories unstarted,completed,cancelled --transitions 'todo->done' --by tl` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"workflow.define","result":{"workflow":{"applies_to":"task","id":"entry-flow","name":"entry-flow","state_count":3,"transition_count":3}},"event":{"id":"EVT-20260815-000002","occurred_at":"2026-08-15T05:43:27Z"}}` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow  (cached)  coverage: 80.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 77.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration  4.379s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- Initial `golangci-lint run ./...` failed with `1 issues:` (terminated `exit_code: 1`) because migration of workflow-list tests made a Markdown helper unused. Removed the obsolete helper and reran the byte-identical command successfully.
- Initial integration command failed with `FAIL` (terminated `exit_code: 1`) because its Phase 19 scenario expected the removed `--require` grammar. Replaced that assertion with the accepted unregistered-flag rejection and reran the byte-identical command successfully.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
