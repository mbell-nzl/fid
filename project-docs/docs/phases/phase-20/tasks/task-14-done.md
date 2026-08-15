# Phase 20 - Task 14: Legacy JSON renderer conversion - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Converted `policy show`, `claim list`, `activity list`, `concern list`, `stats tokens`, `stats cost`, `sql`, and `journal inspect` to direct `encoding/json` output with `SetEscapeHTML(false)` and one trailing newline.
- Preserved unenveloped legacy result shapes and tested semantic fields; retained root `--json` and `--project` rejection before each converted handler.
- Removed formatter imports and `--format` declarations from exactly the eight converted handlers without adding a shared renderer or forward contract.
- Removed `journal inspect`'s `revert` field and the unused `Verdict`/`RevertVerdict` implementation and test family.
- Rewrote every remaining test invocation targeting those eight leaves, including unit, e2e, and `cmd/fid/integration` callers, and deleted formatter helpers orphaned by those rewrites. `migratedReadTestArgs`, ready/issue get/issue list callers, the format matrix, YAML, and `internal/format` remain unchanged.

## Files created or modified

- `cmd/fid/{activity,claim,concern,journal,policy,sql,stats}.go`
- `cmd/fid/{activity,claim,concern,journal,policy,root_flags,sql,stats}_test.go`
- `cmd/fid/{dev,e2e_v01c_concern,e2e_v01e_activity,issue_claim,issue_label,issue_test,label,link,maintenance_vacuum,requirement,run}_test.go`
- `cmd/fid/integration/e2e_v01_closure_helpers_test.go`
- `internal/journal/inspect.go`
- `internal/journal/inspect_test.go`
- `project-docs/docs/phases/phase-20/tasks/task-14-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-20/tasks/task-14-done.md` - completion report with DEC-142 caller-ownership evidence and gate history.

## Validation Contract Coverage

- VAL-028: PASS - `TestJournalInspectCommandUsesDirectLegacyJSONWithoutRevert` asserts retained event/location fields and one trailing newline at `cmd/fid/journal_test.go:41`; `TestJournalInspectCommandDoesNotExposeRevertOrResponseEnvelope` rejects `revert`, `contract_version`, and envelope `event` at `cmd/fid/journal_test.go:65`; `TestRootFlagsRejectBeforeEachConvertedLegacyJSONRead` exercises both root flags for `journal inspect` at `cmd/fid/root_flags_test.go:143`.
- VAL-029: PASS - `TestPolicyShowUsesDirectLegacyJSONWithoutFormatFlag` decodes direct, unenveloped policy rows at `cmd/fid/policy_test.go:15`; the focused green command and full race suite also exercised semantic direct-JSON behavior for claim, activity, concern, tokens, cost, SQL, and journal inspect.
- VAL-031: PASS - `TestRemainingLegacyJSONReadsDoNotDeclareFormatFlags` checks all seven remaining converted constructors at `cmd/fid/policy_test.go:41`, while the policy test checks the eighth at `cmd/fid/policy_test.go:20`; root-flag rejection covers all eight at `cmd/fid/root_flags_test.go:143`, and the caller sweep leaves no live invocation of those eight with `--format`.

## Tests run

- `go test ./cmd/fid -run '^TestPolicyShowUsesDirectLegacyJSONWithoutFormatFlag$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestPolicyShowUsesDirectLegacyJSONWithoutFormatFlag$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  0.046s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^(TestPolicyShowUsesDirectLegacyJSONWithoutFormatFlag|TestRemainingLegacyJSONReadsDoNotDeclareFormatFlags|TestClaimListUsesDirectLegacyJSONAndFiltersByStoredRole|TestActivityListUsesDirectLegacyJSON|TestConcernListUsesDirectLegacyJSON|TestStatsTokensUsesDirectLegacyJSONAndEpicFilter|TestStatsCostUsesDirectLegacyJSONAndExcludesNullCostsFromSum|TestSQLCommandAcceptsReadStatementsAsDirectLegacyJSON|TestJournalInspectCommandUsesDirectLegacyJSONWithoutRevert|TestJournalInspectCommandDoesNotExposeRevertOrResponseEnvelope|TestRootFlagsRejectBeforeEachConvertedLegacyJSONRead)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  0.232s` (terminated `exit_code: 0`)
- `go test ./internal/journal -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/journal  1.679s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  10.052s` (terminated `exit_code: 0`)
- `/tmp/fid policy show` (working directory: `/home/mark/projects/fid`) - PASS - `]` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - BLOCKED - `20 issues:` (terminated `exit_code: 1`) - The newly unused legacy formatter test helpers are owned by the deferred Task 15 caller cleanup.
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`) - Deferred policy/activity/concern `--format json` callers in unrelated and e2e tests now reject as required by the converted handlers.
- `go tool cover -func=coverage.out` - PASS - `77.4%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`) - Deferred integration callers still pass `--format json` to `fid sql`.
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `77.4%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration  4.556s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- The three BLOCKED rows above are immutable first-attempt records from the preserved checkpoint; the later byte-identical PASS rows are the closing results. DEC-142 moved ownership of every invalidated caller and orphaned formatter helper into Task 14, so the repairs below follow the amended boundary.
- `golangci-lint run ./...`: first failure token `20 issues:` with observed `exit_code: 1`; root cause was the formatter helper/parser family orphaned by the Task 14 caller conversion. Repair deleted only those now-unreferenced helpers and their unused imports. Later unchanged result was `0 issues.` with observed `exit_code: 0`.
- `go test -race -coverprofile=coverage.out ./...`: first failure token `FAIL` with observed `exit_code: 1`; root cause was remaining unit/e2e invocations targeting converted policy, activity, and concern leaves with the retired flag. Repair removed only those per-invocation flags while retaining their behavioral assertions, as required by DEC-142's moved caller boundary. Later unchanged result ended with `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 80.6% of statements` and observed `exit_code: 0`.
- `go test -tags integration -race ./cmd/fid/integration/...`: first failure token `FAIL` with observed `exit_code: 1`; root cause was remaining closure-helper invocations targeting converted SQL, stats, and policy leaves with the retired flag. Repair removed only those target invocations and retained the integration assertions. Later unchanged result was `ok  github.com/mbell-nzl/fid/cmd/fid/integration  4.556s` with observed `exit_code: 0`.
- `go test ./cmd/fid -count=1`: an intermediate repair attempt failed with token `FAIL` and observed `exit_code: 1` because `encoding/json` became unused after deleting the orphaned policy formatter helper. Removing that import was the causal in-scope repair; the unchanged command later passed with `ok  github.com/mbell-nzl/fid/cmd/fid  10.052s` and observed `exit_code: 0`.
- No validation assertion, gate command, threshold, protected decision, dependency, or task boundary was weakened or substituted.

## Concerns, technical debt, or notes for the Tech Lead

- None. The remaining `internal/format`, YAML, `migratedReadTestArgs`, ready/issue get/issue list callers, and format-matrix coverage remain intentionally assigned to later tasks.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending task commit finalization.

## Commit subject

- Pending task commit finalization.
