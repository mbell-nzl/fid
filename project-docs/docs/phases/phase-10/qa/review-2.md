# Phase 10 - Mid-Phase Review 2 (Tasks 2-2)

Reviewer: QA Reviewer
Date: 2026-06-21
Tasks covered: phase-10/task-2

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-10/tasks/task-2.md`
- Completion reports:
  - `project-docs/docs/phases/phase-10/tasks/task-2-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-10/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-10/qa/qa-plan.md`
- Design sections checked:
  - `project-docs/specs/05_DESIGN_PLAN.md` sections 1 and 3.1: machine-readable default, terminal/error principles, and command tree entry for `fid policy show`.
  - `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 11: `fid policy show` formats and JSON row shape.
- Commits or changed files:
  - `6f4ce6b03b823b40f7f3b4ea7d6c0f3a59167df3` - `feat(policy): add permissions matrix and show command`
  - `cmd/fid/policy.go`
  - `cmd/fid/policy_test.go`
  - `cmd/fid/root.go`
  - `internal/format/markdown.go`
  - `internal/policy/authorize.go`
  - `internal/policy/authorize_test.go`
  - `internal/policy/matrix.go`
  - `internal/policy/matrix_test.go`

## Summary

- Blocking findings: 0
- Improvement findings: 1
- Escalations raised: none
- Covered validation assertions: VAL-002, VAL-003

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-002 | task-2 | ON TRACK | `cmd/fid/policy.go` renders `internal/policy.MatrixRows()` through the shared format engine for json/toon/markdown/yaml. `TestPolicyShowDefaultsToJSONAndDispatchesFormats` covers default json, explicit json, toon, yaml structural equivalence, markdown output, and action presence. `TestPolicyShowIsReadOnlyAndDoesNotAcceptByFlag` covers no `--by`, unchanged journal line count, unchanged `activity_log` row count, and non-empty output. Reviewer reran `go test ./cmd/fid -run 'TestPolicyShow'` with PASS. |
| VAL-003 | task-2 | ON TRACK | `internal/policy/matrix.go` contains the Phase 10 v0.1a-c implemented action rows from TRD §13 and omits future claim/link/label rows. `TestMatrixContainsExactlyTheImplementedV01WriteActions` asserts exactly one entry for each implemented action and no unexpected action. `TestAuthorizePermitsAllowedCombinationAndRejectsDeniedCombination` covers one allowed and one denied cell with `PermissionDeniedError{Role, Action}`; `TestAuthorizePreservesMissingAndUnknownRoleErrors` covers preserved role validation. Reviewer reran `go test ./internal/policy -run 'Test(MatrixContainsExactlyTheImplementedV01WriteActions|AuthorizePermitsAllowedCombinationAndRejectsDeniedCombination|AuthorizePreservesMissingAndUnknownRoleErrors)'` with PASS. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Exact task command | `gofmt -l cmd internal` | PASS | Completion report records PASS. Required by task spec and QA-plan EVID-10. |
| Exact task command | `golangci-lint run ./...` | PASS | Completion report records PASS. |
| Exact task command | `go test -race -coverprofile=coverage.out ./...` | PASS | Completion report records PASS. |
| Exact task command | `go tool cover -func=coverage.out` | PASS | Completion report records COV-001: `internal/issues` 81.8%, `internal/journal` 80.3%, `internal/policy` 81.2%, `internal/workflow` 95.5%; COV-002: `cmd/fid` 77.8%. Reviewer spot-check of current `coverage.out` is consistent with `internal/policy` 81.2% and `cmd/fid` 77.8%. |
| Conditional task command | `go build -o /tmp/fid ./cmd/fid` | PASS | Completion report records PASS; required because Task 2 added a root subcommand. |
| Conditional task command | `go test -tags integration -race ./...` | N/A | Completion report records N/A because Task 2 added no integration-tagged tests. This matches the task conditional rule. |
| QA-plan evidence | EVID-03 / QA-03 | PASS | Task report and tests cover json/toon/yaml structural equivalence, markdown rendering, default json, and action coverage. |
| QA-plan evidence | EVID-04 / QA-05 | PASS | Task report and tests cover exact action-entry count plus one allowed and one denied `Authorize` cell. |
| QA-plan evidence | EVID-10 | PASS | Completion report lists every required exact task command as PASS or justified N/A for the conditional integration command. |
| QA-plan evidence | EVID-11 | PASS | Completion report includes COV-001 and COV-002 package threshold excerpts. |
| Commit evidence | `git rev-parse --verify 6f4ce6b03b823b40f7f3b4ea7d6c0f3a59167df3^{commit}` | PASS | Hash resolves. `git show --format=fuller --no-patch` confirms subject and `Refs: phase-10/task-2`. |
| Dependency budget | Completion report dependency-budget note | PASS | Report says no dependency changes; commit stat and code inspection showed no dependency file changes. |
| Documentation / decisions | Completion report docs section and `git diff HEAD~1..HEAD` | PASS | Report records verified-no-change for `06_DESIGN_PLAN_CLI.md` §11 and Phase 10 decisions. Commit did not modify migrations, design spec, or decisions. |
| Representative deterministic check | Focused Task 2 tests | PASS | Reviewer reran `go test ./internal/policy -run 'Test(MatrixContainsExactlyTheImplementedV01WriteActions|AuthorizePermitsAllowedCombinationAndRejectsDeniedCombination|AuthorizePreservesMissingAndUnknownRoleErrors)'` and `go test ./cmd/fid -run 'TestPolicyShow'`; both passed. |

## Per-Task Findings

### Task 2 - In-code permissions matrix + Authorize seam + `fid policy show`

#### Blocking

- None.

#### Improvement

- `internal/policy/matrix_test.go:8` The matrix completeness test hard-codes the expected action keys, but it does not hard-code the expected role booleans for every TRD §13 v0.1a-c row. Current source inspection shows the implemented cells match TRD §13, and `Authorize` samples one allowed and one denied cell, so this is not a Task 2 blocker. Fix direction: add a table of expected `MatrixRow` values for every implemented action when this matrix next changes, so future role-cell drift fails without relying on review inspection.

## Cross-cutting observations

- DEC-B is structurally followed in this task: `Authorize` looks up `MatrixEntry`, and `fid policy show` writes `MatrixRows()`, so display and enforcement share the same in-code matrix object rather than maintaining parallel copies.
- The Task 2 window intentionally stops before write-command enforcement. VAL-008 show-vs-enforce parity remains a Task 3 responsibility; this review did not assess Task 3.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-10/qa/review-2.md
Blocking: 0  Improvement: 1  Escalations: none
VERDICT: CONTINUE
```

CONTINUE
