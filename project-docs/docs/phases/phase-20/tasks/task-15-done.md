# Phase 20 - Task 15: Format shim removal and whole-tree caller convergence - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Deleted `migratedReadTestArgs`, `containsFormatJSON`, `isLegacyIssueSection`,
  `legacyMigratedReadOutput`, and the complete argument/output rewriting shim
  from `cmd/fid/issue_test.go`.
- Converted remaining ready, issue get, issue list, and legacy section callers
  to the real command surfaces: canonical `--json` envelopes for ready/get/list
  and the owning `concern`, `issue ac`, `claim`, or retained projection reads
  for section semantics.
- Preserved the behavioral assertions by checking canonical omissions against
  persisted ownership and by retaining concern, acceptance-verification, and
  completion parity assertions through their current command/DB surfaces.
- Converged the V1/V2a command and e2e parity callers, including completion
  ordering and rebuild parity. The live command tree remains free of universal
  format flags; all named retired-format rejection tests and the V1 format
  matrix remain in place.
- Left `internal/format` and the eight Task 14 leaves unchanged.

## Files created or modified

- `cmd/fid/contract_surface_test.go`
- `cmd/fid/concern_test.go`
- `cmd/fid/dev_test.go`
- `cmd/fid/e2e_v01b_test.go`
- `cmd/fid/e2e_v01c_ac_test.go`
- `cmd/fid/e2e_v01c_complete_test.go`
- `cmd/fid/e2e_v01c_concern_test.go`
- `cmd/fid/e2e_v01e_dev_test.go`
- `cmd/fid/issue_ac_test.go`
- `cmd/fid/issue_claim_test.go`
- `cmd/fid/issue_complete_test.go`
- `cmd/fid/issue_label_test.go`
- `cmd/fid/issue_list_test.go`
- `cmd/fid/issue_test.go`
- `cmd/fid/issue_transition_test.go`
- `cmd/fid/issue_update_test.go`
- `cmd/fid/maintenance_test.go`
- `cmd/fid/ready_test.go`
- `project-docs/docs/phases/phase-20/tasks/task-15-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-20/tasks/task-15-done.md` - finalized this
  same-task recovery report with deterministic gate evidence and VAL coverage.
- No pre-existing Phase Runner or user-authored documentation was changed by
  this continuation.

## Validation Contract Coverage

- VAL-029: PASS. The converged command and V1/V2a parity tests retain semantic
  assertions for concern, acceptance-criterion verification, completion, ready,
  issue get, and issue list behavior. The exact race/coverage gate and the
  integration gate both pass.
- VAL-031: PASS. `TestLiveTestsContainNoUniversalFormatInvocation` passes after
  the full caller conversion; the live source contains no shim helper names,
  and remaining `--format` literals are confined to retained rejection tests,
  including `e2e_v01_format_matrix_test.go`. The root command-tree guard is
  included in the passing full suite.
- Retained VAL-016 negative evidence remains covered by the named rejection
  tests and `TestV01ReadCommandsRejectRetiredFormats`; explicit root flags on
  legacy leaves remain rejected.

## Tests run

- `go test ./cmd/fid -run '^TestLiveTestsContainNoUniversalFormatInvocation$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`) - immutable prior focused guard failure before caller convergence; the report checkpoint records the earlier inventory.
- `go test ./cmd/fid -run '^TestLiveTestsContainNoUniversalFormatInvocation$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  0.080s` - later byte-identical guard pass (terminated `exit_code: 0`)
- `go test ./cmd/fid -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`) - initial convergence exposed retired owner, concern, acceptance-verification, and completion read assumptions.
- `go test ./cmd/fid -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`) - an intermediate completion-order repair temporarily changed the agent-run query to use a nonexistent column.
- `go test ./cmd/fid -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  10.051s` - final focused command suite (terminated `exit_code: 0`)
- `go test ./...` - FAIL - `FAIL` (terminated `exit_code: 1`) - first full non-race run exposed completion ordering after replacing the retired read section.
- `rg -n --glob '*.go' 'migratedReadTestArgs|containsFormatJSON|isLegacyIssueSection|legacyMigratedReadOutput' cmd internal` - PASS no-match inventory - `exit 1 (no output)` (terminated `exit_code: 1`); no-match is the pass condition.
- `/tmp/fid init --prefix FID` - PASS - `initialized fid project in .fid/` - real built CLI entrypoint (terminated `exit_code: 0`)
- `/tmp/fid --json issue list` - PASS - `{"contract_version":"fid-cli.v1","ok":true,"operation":"issue.list","result":{"count":0,"issues":[]},"event":null}` - real built CLI entrypoint (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - FAIL - `* unused: 2` (terminated `exit_code: 1`) - first exact attempt.
- `golangci-lint run ./...` - PASS - `0 issues.` - byte-identical rerun after removing two dead helpers (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`) - first exact attempt; the race/coverage timing crossed the RFC3339-truncated 299-second fixture boundary.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 80.6%` - byte-identical rerun after moving only the negative fixture to 298 seconds (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t77.4%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 5.136s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- Mandatory recovery deviation: the prior focused guard FAIL is preserved as
  immutable evidence, and the later byte-identical guard PASS is recorded above.
  Root cause was the in-scope remaining caller inventory; the repair converted
  callers to real command-family surfaces without weakening the guard or
  rejection tests.
- Mandatory gate deviation: the first lint attempt ended with `* unused: 2`
  and `exit_code: 1`; the two dead concern-state helpers left by section
  conversion were removed, and the unchanged lint command later ended with
  `0 issues.` and `exit_code: 0`.
- Mandatory gate deviation: the first race/coverage attempt ended with `FAIL`
  and `exit_code: 1`; under instrumentation the RFC3339-truncated 299-second
  negative maintenance fixture crossed the 300-second boundary. Only that
  fixture was moved to 298 seconds; the fixed-300-second positive assertion and
  retired-format rejection assertion were retained. The unchanged command then
  ended with the recorded passing package token and `exit_code: 0`.
- No product, architecture, public-interface, dependency, gate-definition, or
  assertion-strength decision was changed.

## Concerns, technical debt, or notes for the Tech Lead

- `internal/format` remains intentionally present for Task 16.
- No protected blocker or open escalation was created.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `f95db81f938e15d59357a35936e43a24c708b653`

## Commit subject

- `test(cli): converge format shim callers`
