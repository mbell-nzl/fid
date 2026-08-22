# Phase 20 - Task 15: Format shim removal and whole-tree caller convergence - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 2
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
- Added the missing dynamic full-Cobra-tree guard in
  `cmd/fid/contract_surface_test.go`. It recursively inspects every live
  command's local non-persistent, persistent, and inherited flag sets for
  `format`, `toon`, `yaml`, and `markdown`. Its inherited-root negative
  control injects a root persistent `format` flag and requires the violation
  on the descendant path `fid issue --format`; the local-leaf control remains
  in place.
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
- VAL-031: PASS. `TestLiveCobraTreeContainsNoUniversalFormatFlags` recursively
  walks the complete `newRootCommandAt` Cobra tree and checks local,
  persistent, and inherited flag sets for all four universal-format names.
  `TestLiveCobraTreeUniversalFormatGuardDetectsInjectedFlags` now proves the
  inherited-root injection is observed on the descendant path `fid issue
  --format`, as well as proving a local leaf flag is detected. The existing
  `TestLiveTestsContainNoUniversalFormatInvocation` AST guard also passes; no
  shim helper names remain, and remaining `--format` literals are confined to
  retained rejection tests including `e2e_v01_format_matrix_test.go`.
- Retained VAL-016 negative evidence remains covered by the named rejection
  tests and `TestV01ReadCommandsRejectRetiredFormats`; explicit root flags on
  legacy leaves remain rejected.

## Tests run

- `go test ./cmd/fid -run '^TestLiveCobraTree' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`) - red-first compile failure before the dynamic walker was added; the new tests referenced the not-yet-defined guard helper.
- `go test ./cmd/fid -run '^TestLiveCobraTree' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  0.050s` - dynamic live-tree guard and injected inherited/local negative controls (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^(TestLiveCobraTreeContainsNoUniversalFormatFlags|TestLiveCobraTreeUniversalFormatGuardDetectsInjectedFlags|TestLiveTestsContainNoUniversalFormatInvocation|TestIssueGetRejectsRetiredFormatFlags|TestIssueGetRejectsAllFormatValues|TestIssueListRejectsRetiredFormatFlag|TestLabelListRejectsLegacyFormatFlag|TestReadyRejectsRetiredFormatFlag|TestIssueACListCommandRejectsLegacyFormatsAndHandlesEmptyAndMissingIssues|TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat|TestV01ReadCommandsRejectRetiredFormats|TestIssueReadRejectsRetiredAndMalformedInputs)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  0.292s` - dynamic guard, AST source guard, and retained rejection matrix (terminated `exit_code: 0`)
- Review 5c MF-02 category evidence (immutable red finding): the prior
  inherited negative control accepted only the root violation `fid
  --format`, which did not prove that descendant inherited flags were
  inspected. The repaired control requires `fid issue --format`.
- `go test ./cmd/fid -run '^TestLiveCobraTreeUniversalFormatGuardDetectsInjectedFlags$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  0.050s` - strengthened inherited-root descendant control and retained local-leaf control (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^(TestLiveCobraTreeContainsNoUniversalFormatFlags|TestLiveCobraTreeUniversalFormatGuardDetectsInjectedFlags|TestLiveTestsContainNoUniversalFormatInvocation|TestIssueGetRejectsRetiredFormatFlags|TestIssueGetRejectsAllFormatValues|TestIssueListRejectsRetiredFormatFlag|TestLabelListRejectsLegacyFormatFlag|TestReadyRejectsRetiredFormatFlag|TestIssueACListCommandRejectsLegacyFormatsAndHandlesEmptyAndMissingIssues|TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat|TestV01ReadCommandsRejectRetiredFormats|TestIssueReadRejectsRetiredAndMalformedInputs)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  0.294s` - strengthened dynamic guard, AST source guard, and retained rejection matrix (terminated `exit_code: 0`)
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
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` - MF-02 repair rerun (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` - MF-02 repair rerun (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 80.6%` - MF-02 repair rerun (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t77.4%` - MF-02 repair rerun (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 4.734s` - MF-02 repair rerun (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` - MF-02 repair rerun (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` - Review 5c inherited-control repair rerun (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` - Review 5c inherited-control repair rerun (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 80.6%` - Review 5c inherited-control repair rerun (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t(statements)\t77.4%` - Review 5c inherited-control repair rerun (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 4.776s` - Review 5c inherited-control repair rerun (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` - Review 5c inherited-control repair rerun (terminated `exit_code: 0`)

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
- Review 5 MF-02 recovery deviation: the red-first dynamic guard test initially
  ended with `FAIL` and `exit_code: 1` because the walker helper was not yet
  defined. The repair added the test-only recursive Cobra traversal and its
  inherited/local injected-flag negative controls. The existing AST guard and
  retained rejection tests were not weakened or changed; the focused matrix
  and all unchanged Task 15 gates then passed.
- Review 5c MF-02 repair deviation: the prior root-only inherited negative
  control was insufficient because it could pass from the root persistent flag
  set without proving descendant inherited inspection. The test-only control
  was strengthened to require the injected flag at `fid issue --format`; the
  complete recursive walker, local-leaf control, AST guard, retained rejection
  tests, root-flag tests, and Task 14 guards were preserved. The focused
  strengthened control, retained matrix, and unchanged Task 15 gates ended
  with the passing evidence recorded above.
- No product, architecture, public-interface, dependency, gate-definition, or
  assertion-strength decision was changed.

## Concerns, technical debt, or notes for the Tech Lead

- `internal/format` remains intentionally present for Task 16.
- No protected blocker or open escalation was created.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `24f50aef13ae7a23a080a87920f396242b30f911`

## Commit subject

- `test(cli): strengthen inherited format guard`
