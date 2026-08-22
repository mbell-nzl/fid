# Phase 20 - Task 12: Workflow command contract - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: Review 4 MF-02 is resolved within the accepted Task 12 repair boundary.
- Recovery attempts: 2
- Next action: Fresh independent QA re-review of MF-02.
- Human attention required: no

## Summary of what was implemented

- The original Task 12 migration remains recorded at `f081929cfed7a7be737010c2fb14403c61dac7da`.
- Resolved Review 4 MF-02: `workflow.list` reads at most 10,001 ordered rows, returns a closed typed `result-limit-exceeded` failure for 10,001 workflows, and never silently truncates the collection.
- Added the bounded workflow evidence matrix: 10,000/10,001 list states, `name,id` ordering, define/assign journal-tail Event identity, list `event:null`, root `--project <path>` plus `--project-wide`, local-flag absence, closed grammar/scope rejection, and no-append rejection checks.

## Files created or modified

- `internal/workflow/list.go`
- `internal/workflow/errors.go`
- `internal/workflow/list_test.go`
- `cmd/fid/workflow.go`
- `cmd/fid/workflow_define_test.go`
- `cmd/fid/workflow_assign_test.go`
- `cmd/fid/workflow_list_test.go`
- `project-docs/docs/phases/phase-20/tasks/task-12-done.md`

## Documentation updated

- None required.

## Validation Contract Coverage

- VAL-008: PASS - `TestWorkflowDefineHumanSuccessUsesStdoutOnly` in `cmd/fid/workflow_define_test.go` proves Human success is stdout-only; the preserved real `/tmp/fid ... workflow define` invocation also terminated `exit_code: 0` with JSON on stdout only.
- VAL-017: PASS - `TestWorkflowDefineJSONMatchesContract` and `TestWorkflowAssignJSONPrecedenceAndScopedReplacement` compare both response Event ID and RFC3339 timestamp with the journal tail; `TestWorkflowListJSONMatchesContractAndUsesNullAssignmentScope` and `TestWorkflowListEnforcesCollectionLimitWithoutTruncation` assert `event:null` for successful reads.
- VAL-026: PASS - `TestWorkflowDefineJSONMatchesContract`, `TestWorkflowDefineRejectsRemovedAndInvalidRequestsBeforeAppend`, `TestWorkflowAssignJSONPrecedenceAndScopedReplacement`, and `TestWorkflowAssignCombinesRootProjectSelectorAndProjectWideScope` prove the closed schemas, removed V1 rule authoring, kind-over-project precedence, scoped replacement, root-project/project-wide combination, and absence of a local `--project` boolean.
- VAL-041: PASS - `TestWorkflowDefineJSONMatchesContract` verifies the journal-tail actor is `tl` and the Event has only the expected workflow effects; no trusted/protected effect is present.
- VAL-044: PASS - `TestWorkflowDefineRejectsRemovedAndInvalidRequestsBeforeAppend` covers `--require`, legacy applicability, malformed CSV/arrow, and state bounds without append; `TestWorkflowAssignRequiresExactlyOneScopeAndMapsUnknownWorkflow` covers absent/both/invalid scope and unknown workflow without append; `TestWorkflowListRejectsNonIssueKind` covers invalid list kind.
- VAL-045: PASS - `TestListOrdersWorkflowsByNameThenID` proves the `name,id` tie-break; `TestListRejectsCollectionAboveLimit` and `TestWorkflowListEnforcesCollectionLimitWithoutTruncation` prove exactly 10,000 succeeds, 10,001 returns `{collection:"workflows",limit:10000,observed:10001}`, the journal is unchanged, and the projection still holds all 10,001 rows.

## Tests run

- `go test ./cmd/fid --run '^TestWorkflowDefineJSONMatchesContract$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); original Task 12 red test exposed unavailable root JSON transport.
- `go test ./cmd/fid --run '^TestWorkflowDefineJSONMatchesContract$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  0.051s` (terminated `exit_code: 0`)
- `go test ./cmd/fid --run '^TestWorkflowListEnforcesCollectionLimitWithoutTruncation$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); the new command-boundary test received a successful 10,001-item `workflow.list` envelope before the cap existed.
- `go test ./cmd/fid --run '^TestWorkflowListEnforcesCollectionLimitWithoutTruncation$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  1.410s` (terminated `exit_code: 0`)
- `go test ./internal/workflow --run '^(TestListFiltersCountsAssignmentsAndOrdering|TestListReturnsEmptySliceForNoMatch|TestListOrdersWorkflowsByNameThenID|TestListRejectsCollectionAboveLimit)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); the initial ordering fixture omitted bootstrap-seeded workflow rows.
- `go test ./internal/workflow --run '^(TestListFiltersCountsAssignmentsAndOrdering|TestListReturnsEmptySliceForNoMatch|TestListOrdersWorkflowsByNameThenID|TestListRejectsCollectionAboveLimit)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow  1.493s` (terminated `exit_code: 0`)
- The internal fixture repair clears bootstrap-seeded workflow rows before asserting the isolated `name,id` order.
- `go test ./cmd/fid --run '^(TestWorkflowDefineJSONMatchesContract|TestWorkflowDefineRejectsRemovedAndInvalidRequestsBeforeAppend|TestWorkflowAssignJSONPrecedenceAndScopedReplacement|TestWorkflowAssignRequiresExactlyOneScopeAndMapsUnknownWorkflow|TestWorkflowAssignCombinesRootProjectSelectorAndProjectWideScope|TestWorkflowListJSONMatchesContractAndUsesNullAssignmentScope|TestWorkflowListRejectsNonIssueKind|TestWorkflowListEnforcesCollectionLimitWithoutTruncation)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  1.487s` (terminated `exit_code: 0`)
- `/tmp/fid --project /tmp/fid-task12-mf02.95sInD --json workflow assign --name entry-flow --project-wide --by tl` - FAIL - `{"contract_version":"fid-cli.v1","ok":false,"operation":"workflow.assign","error":{"code":"invalid-argument","category":"invalid-input","retryable":false,"message":"invalid Issue mutation request","subject":null,"details":{"allowed":["accepted workflow.assign input"],"argument":"request","value":null}}}` (terminated `exit_code: 1`); a task-only workflow is not a valid project default.
- `/tmp/fid --project /tmp/fid-task12-mf02.95sInD --json workflow assign --name entry-any --project-wide --by tl` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"workflow.assign","result":{"assignment":{"scope_type":"project","scope_value":null,"workflow_id":"entry-any"}},"event":{"id":"EVT-20260815-000004","occurred_at":"2026-08-15T08:18:12Z"}}` (terminated `exit_code: 0`)
- `/tmp/fid --project /tmp/fid-task12-mf02.95sInD --json workflow list` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"workflow.list","result":{"count":6,"workflows":[{"active":true,"applies_to":"epic","assigned_to":[],"id":"epic-default","name":"Epic default","states":5,"transitions":7},{"active":true,"applies_to":"escalation","assigned_to":[],"id":"escalation-default","name":"Escalation default","states":2,"transitions":1},{"active":true,"applies_to":"subtask","assigned_to":[],"id":"subtask-default","name":"Subtask default","states":4,"transitions":5},{"active":true,"applies_to":"task","assigned_to":[],"id":"task-default","name":"Task default","states":9,"transitions":17},{"active":true,"applies_to":null,"assigned_to":[{"scope_type":"project","scope_value":null}],"id":"entry-any","name":"entry-any","states":3,"transitions":3},{"active":true,"applies_to":"task","assigned_to":[],"id":"entry-flow","name":"entry-flow","states":3,"transitions":3}]},"event":null}` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow  19.677s  coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:                                  (statements)                         77.8%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration  4.601s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- No specification deviations. The original Task 12 required-gate retries remain preserved: lint initially ended `1 issues:` (terminated `exit_code: 1`) after an obsolete Markdown helper became unused, and integration initially ended `FAIL` (terminated `exit_code: 1`) because its Phase 19 scenario expected removed `--require` grammar; both were repaired inside original Task 12 scope and rerun byte-identically.
- MF-02's focused red test exposed the missing collection cap and the initial internal ordering fixture omitted seeded defaults; both repairs preserve the accepted contract, test intent, command, and assertion strength.
- Initial `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-12-done.md project-docs/docs/phases/phase-20/tasks/task-12.md` ended `FAIL report-lint found 1 failure(s)` (terminated `exit_code: 1`) because the passing internal-test row had explanatory prose after its termination marker. The evidence was retained and the prose moved to its own line before rerunning the byte-identical command.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `49c179add40724871bae582fc3812ffc0aff04fc`

## Commit subject

- `fix(workflow): enforce workflow list limit`
