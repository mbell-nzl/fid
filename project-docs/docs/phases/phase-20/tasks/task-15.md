# Phase 20 - Task 15: Format shim removal and whole-tree caller convergence

> **Amended 2026-08-16 by DEC-142 (resolves ESC-20-011).** Task 14 now converges
> the callers of the eight leaves it converts, so this task no longer carries
> another task's gate. Its real remaining scope — the `migratedReadTestArgs`
> compatibility shim, the callers routed through it, the parity suites, and the
> two guards — is now stated explicitly.

## Objective

Delete the `--format json` test-harness compatibility shim, rewrite the callers
that still depend on it, and lock the result with repository and command-tree
guards. Keep behavioral assertions intact so Task 16 can delete the formatter
framework against a green, caller-free tree.

## Validation Contract Coverage

- The whole-tree test/caller portion of VAL-031; VAL-029 regression coverage
  must remain green while callers are rewritten.

## Binding Constraints

- Delete `migratedReadTestArgs`, `containsFormatJSON`, `isLegacyIssueSection`,
  `legacyMigratedReadOutput`, and the rest of the `--format json` → `--json`
  rewriting shim in `cmd/fid/issue_test.go`. Callers must invoke the real
  surface. Do not relocate or rename the shim to keep it alive.
- Update command-family tests and remaining V1/e2e invocations; do not delete
  behavioral assertions merely to remove a flag.
- **Retain every test whose assertion is that a retired format flag is
  rejected** — for example `TestIssueGetRejectsRetiredFormatFlags`,
  `TestIssueGetRejectsAllFormatValues`, `TestIssueListRejectsRetiredFormatFlag`,
  `TestLabelListRejectsLegacyFormatFlag`, `TestReadyRejectsRetiredFormatFlag`,
  `TestIssueACListCommandRejectsLegacyFormatsAndHandlesEmptyAndMissingIssues`,
  `TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat`, and
  `cmd/fid/e2e_v01_format_matrix_test.go`. These are required VAL-016 and
  VAL-031 negative evidence. The source guard must exempt them by construction,
  not by deleting them.
- Leave `internal/format` in place until Task 16 so source guards and the
  compiler identify missed callers.

## Scope

### Files To Create Or Modify

- `cmd/fid/issue_test.go` - delete the shim and its helpers.
- The `ready`, `issue get`, `issue list`, and legacy `--section` callers that
  the shim was carrying, across command and e2e tests.
- `internal/db/issues_workflow_fk_cli_test.go` and retained V1/V2a parity suites.
- Repository source guard `TestLiveTestsContainNoUniversalFormatInvocation` and
  the root command-tree guard test for zero runtime format flags.

### Files To Read For Context

- `../contracts.md` sections 5.2, 5.3, and 6.
- `../decisions.md` DEC-142 - caller-ownership rule and the Task 14/15 boundary.

### Out Of Scope

- Callers of the eight legacy leaves converted in Task 14; that task owns them
  and left the tree green.
- Deleting `internal/format`, YAML dependency removal,
  `cmd/fid/e2e_v01_format_matrix_test.go` deletion, legacy result redesign, and
  command stubs.

## Acceptance Criteria

- [ ] Live Cobra tree contains no universal format flag. (VAL-031)
- [ ] The `--format json` → `--json` test shim is absent and no equivalent argument-rewriting helper replaces it. (VAL-031)
- [ ] All remaining tests invoke the new surface without weakening semantic assertions. (VAL-029, VAL-031)
- [ ] Source/test search finds no live `--format` invocation outside historical docs and retained rejection assertions. (VAL-031)
- [ ] Retained rejection tests still assert non-zero exit for every retired format value. (VAL-016, VAL-031)
- [ ] Explicit root flags remain rejected on legacy leaves in retained regression coverage. (VAL-029)

## First Failing Test

`cmd/fid/contract_surface_test.go` -
`TestLiveTestsContainNoUniversalFormatInvocation` (repository source guard,
exempting invocations asserted to be rejected). It must be red on the tree Task
14 leaves green, because the shim and its callers still exist.

## Documentation Updates

- None expected.

## Exact Commands To Run

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go build -o /tmp/fid ./cmd/fid
```

## Notes For Developer

- A discovered live product caller belongs to its owning Task 8-14. Report it;
  it is not permission to widen this cleanup task. Under DEC-142 an in-scope
  test caller that your own changes break is yours, not the next task's.
- Verify the shim's true reach before deleting it; the planning-time reading is
  `ready`, `issue get`, and `issue list` plus the legacy `--section` unwrap.
  Report any divergence rather than assuming the plan is exhaustive.
- Completion report: `task-15-done.md`.
