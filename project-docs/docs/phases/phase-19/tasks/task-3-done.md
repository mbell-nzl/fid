# Phase 19 - Task 3: `fid workflow define` command + `workflow.defined` event - Completion Report

## Summary of what was implemented

- Registered the shipped `fid workflow` command family and implemented
  `fid workflow define` with the specified flags, authorization-first parsing,
  four arrow aliases, exact applicability vocabulary, rule parsing, and the
  documented stderr confirmation.
- Added the workflow write path: Task 2 effects are validated before the write
  lock, then appended as one current-schema `workflow.defined` event and
  projected atomically without an `activity_log` effect.
- Added the `workflow.define` and `workflow.assign` policy rows, permitting
  only `tl` and `human`, and exposed them in every `fid policy show` format.

## Files created or modified

- `cmd/fid/root.go`
- `cmd/fid/workflow.go`
- `cmd/fid/workflow_define_test.go`
- `internal/workflow/define.go`
- `internal/workflow/define_test.go`
- `internal/policy/matrix.go`
- `internal/policy/matrix_test.go`
- `cmd/fid/policy_test.go`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` — shipped define flags,
  applicability, arrows, rules, cancellation policy, initial-state rule, and
  error contract.
- `project-docs/specs/01_PRD.md` — PRD-WF-013 DEC-077/DEC-090 scoping note.
- Verified `project-docs/specs/05_DESIGN_PLAN.md` §3.1 already matches the
  shipped `workflow define` flag surface; no change was required.

## Validation Contract Coverage

- VAL-009: PASS — QA plan EVID-09 / QA-09; `TestWorkflowDefineProjectsDeclaredStatesTransitionsAndRules` asserts one appended line, `workflow.defined`, `journal.CurrentEventSchema`, and all 12 effects at `cmd/fid/workflow_define_test.go:47`.
- VAL-043: PASS — QA plan EVID-43 / QA-48; tuple-for-tuple workflow, states,
  transitions (including cancel edges), and requirements assertions are at
  `cmd/fid/workflow_define_test.go:250`.
- VAL-044: PASS — QA plan EVID-44 / QA-49; the newly projected definition is
  loaded from `workflow.Load` at `cmd/fid/workflow_define_test.go:63`.
- VAL-045: PASS — QA plan EVID-45 / QA-50; before/after `activity_log` count
  equality is asserted at `cmd/fid/workflow_define_test.go:25` and `:70`.
- VAL-010: PASS — QA plan EVID-10 / QA-10; accepted control appends one event
  at `cmd/fid/workflow_define_test.go:47`, while every enumerated rejection
  byte-compares the journal at `:140` and `:190`.
- VAL-046: PASS — QA plan EVID-46 / QA-51; every rejection snapshots and
  deep-compares all four workflow tables at `cmd/fid/workflow_define_test.go:170`
  and `:193`.
- VAL-047: PASS — QA plan EVID-47 / QA-52; the rejection matrix asserts
  non-zero command errors, value-naming stderr, and distinct messages at
  `cmd/fid/workflow_define_test.go:176`.
- VAL-011: PASS — QA plan EVID-11 / QA-11; exact `tl`/`human`-only rows for
  both actions are asserted by
  `TestMatrixWorkflowAuthoringRowsPermitOnlyTLAndHuman` at
  `internal/policy/matrix_test.go:84`.
- VAL-048: PASS — QA plan EVID-48 / QA-53; JSON, TOON, YAML, and Markdown
  each assert both rows in
  `TestPolicyShowIncludesWorkflowAuthoringRowsInEveryFormat` at
  `cmd/fid/policy_test.go:196`.
- VAL-049: PASS — QA plan EVID-49 / QA-54; the denied developer command has
  the standard permission error and unchanged journal at
  `cmd/fid/workflow_define_test.go:75`, while the successful `tl` control is
  at `:27`.
- VAL-012: PASS — QA plan EVID-12 / QA-12; all `->`, `-->`, `→`, and `⟶`
  cases assert identical projected tuples at
  `cmd/fid/workflow_define_test.go:203`.
- VAL-013: PASS — QA plan EVID-13 / QA-13; the byte-exact confirmation and
  SQL-backed projected counts are asserted at
  `cmd/fid/workflow_define_test.go:43` and `:265`.
- VAL-050: PASS — QA plan EVID-50 / QA-55; the six accepted values and `any`
  SQL `NULL` are asserted at `cmd/fid/workflow_define_test.go:215`, while the
  invalid value is rejected without mutation in the matrix at `:146` and
  `:190`.

## Tests run

- `go test ./cmd/fid --run '^TestWorkflowDefineProjectsDeclaredStatesTransitionsAndRules$' -count=1` — RED expected failure before implementation — `FAIL`.
- `go test ./internal/policy --run '^TestMatrixContainsExactlyTheImplementedV01WriteActions$' -count=1` — RED expected failure before matrix-test update — `FAIL`.
- `go test ./cmd/fid --run '^TestWorkflowDefineAuthorizesBeforeParsingAndDoesNotAppendOnDenial$' -count=1` — RED expected failure before authorization ordering repair — `FAIL`.
- `go test ./cmd/fid --run '^TestWorkflowDefinePreservesJSONRequirementParametersAfterTheEdgeColon$' -count=1` — RED expected failure before JSON-array parsing repair — `FAIL`.
- `go test ./cmd/fid ./internal/workflow ./internal/policy -count=1` — PASS — `ok github.com/mbell-nzl/fid/internal/policy 0.081s`.
- Built CLI exercise: `mktemp -d /tmp/fid-workflow-define-XXXXXX`; then
  `/tmp/fid init` and `/tmp/fid workflow define --name runtime-flow --applies-to task --states todo,in_progress,done,cancelled --categories unstarted,started,completed,cancelled --transitions 'todo→in_progress,in_progress-->done' --require 'in_progress->done:all_acs_verified{"roles":["developer"]}' --require 'in_progress->done:has_completion_report' --by tl` — PASS — `requirements: 2` (observed block: `defined workflow runtime-flow for task`, `states: 4`, `transitions: 5`, `requirements: 2`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`.
- `golangci-lint run ./...` - PASS - `0 issues.`
- `go test -race -coverprofile=coverage.out ./...` - PASS - `exit 0 (no output)`.
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 77.1%`.
- `go test -tags integration -race ./...` - PASS - `exit 0 (no output)`.
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)`.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
