# Phase 20 - Task 3c: Integration-journey convergence and honest gate baseline - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added explicit rejected `executor_limits.max_excerpt_bytes` rows for zero, negative, over-`MaxInt64`, float, and quoted values.
- Converged all three binary-driven integration initializations on the accepted explicit `--prefix FID` form, preserving their scenarios and assertions.
- No product code changed. Derived-prefix rejection remains contract-correct, and the accepted `runInitCommand` `--prefix FID` injection remains unchanged.
- Report-only continuation: added the canonical workflow disposition and current-session gate termination evidence required by MF-01, without changing product, tests, task scope, or historical evidence.

## Files created or modified

- `internal/config/config_test.go` — added five rejected `max_excerpt_bytes` boundary/type rows.
- `cmd/fid/e2e_v01_closure_test.go` — `TestV01ClosureJourneyAndRebuildParity` now supplies `--prefix FID` to binary `fid init`.
- `cmd/fid/e2e_v02a_workflow_authoring_test.go` — `TestV02aWorkflowAuthoringExitScenarioAndRebuildParity` and `TestMaintenanceCheckReportsIncompatibleWorkflowAssignment` now supply `--prefix FID` to binary `fid init`.
- `project-docs/docs/phases/phase-20/tasks/task-3c-done.md` — this completion report.

## Documentation updated

- None required; DEC-131 already records the fixture-convergence disposition.

## Validation Contract Coverage

- VAL-043: PASS — `TestParseProjectConfigBoundaryDomains` has accepted minimum/maximum/equal-to-output rows at `internal/config/config_test.go:151-152,163` and explicit rejected zero, negative, overflow, float, and quoted excerpt rows at `internal/config/config_test.go:158-162`; the focused run ended `ok   github.com/mbell-nzl/fid/internal/config 1.033s`.
- VAL-044: PASS — the binary integration entrypoint reaches all preserved scenarios with explicit conforming prefixes at `cmd/fid/e2e_v01_closure_test.go:19` and `cmd/fid/e2e_v02a_workflow_authoring_test.go:17,79`; the focused real-binary run passed all three tests, ending `ok   github.com/mbell-nzl/fid/cmd/fid 10.675s`, and the repository-wide binary-fixture sweep found only these three `init` calls, each with `--prefix FID`.

## Tests run

- Pre-repair integration-tag reproduction: `go test -tags integration -race ./cmd/fid -run 'TestV01ClosureJourneyAndRebuildParity|TestV02aWorkflowAuthoringExitScenarioAndRebuildParity|TestMaintenanceCheckReportsIncompatibleWorkflowAssignment' -count=1 -v` — FAIL — real output tail:

  ```text
  --- FAIL: TestMaintenanceCheckReportsIncompatibleWorkflowAssignment (0.56s)
  FAIL
  FAIL	github.com/mbell-nzl/fid/cmd/fid	1.808s
  FAIL
  ```

  Each of the three failures named a `t.TempDir()` `.../001/fid init` and printed `error: invalid project initialization argument`.
- Provenance: `git log --oneline -S 'TestInitCommandRequiresExplicitPrefixForInvalidDirectoryBasename' -- cmd/fid` — PASS — `cd53c43 feat(project): add strict project selection config`. `git show --no-patch --format=... cd53c43`, task-3a, and task-3b shows introducing commit `cd53c43319a3301d3a8cc0089b6785d9afbb163f` (`Refs: phase-20/task-2`, `2026-08-09T16:16:51+12:00`) predates task-3a (`51763124...`, `19:46:24+12:00`) and task-3b (`4f308c86...`, `19:57:51+12:00`).
- Boundary coverage: `go test -race ./internal/config -run '^TestParseProjectConfigBoundaryDomains$' -count=1` — PASS — `ok   github.com/mbell-nzl/fid/internal/config 1.033s`.
- Repaired real-entrypoint journeys: `go test -tags integration -race ./cmd/fid -run 'TestV01ClosureJourneyAndRebuildParity|TestV02aWorkflowAuthoringExitScenarioAndRebuildParity|TestMaintenanceCheckReportsIncompatibleWorkflowAssignment' -count=1 -v` — PASS — `ok   github.com/mbell-nzl/fid/cmd/fid 10.675s`.
- Binary-init fixture sweep: `rg -n -U '(?:RunFID|runFID|closureRunFID)\\([\\s\\S]{0,180}"init"' cmd/fid --glob '*_test.go'` — PASS — final matching line: `cmd/fid/e2e_v01_closure_test.go:19:\tclosureRunFID(t, binary, workspace, "init", "--prefix", "FID")`.

- MF-01 initial report conformance check: `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-3c-done.md project-docs/docs/phases/phase-20/tasks/task-3c.md` - FAIL - `FAIL report-lint found 21 failure(s)` (terminated `exit_code: 21`).
- Current provenance verification: `git log --oneline -S 'TestInitCommandRequiresExplicitPrefixForInvalidDirectoryBasename' -- cmd/fid` - PASS - `cd53c43 feat(project): add strict project selection config` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`).
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.6% of statements` (terminated `exit_code: 0`).
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					78.7%` (terminated `exit_code: 0`).
- `go test -tags integration -race ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)` (terminated `exit_code: 0`).
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).

## Deviations from spec

- None. The initial canonical report-lint failure is preserved above as MF-01 repair evidence; this continuation changed only the completion report and reran the exact Task 3c commands.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `f978b6d382b3ae90ff44a7f4577799bb048329d4`

## Commit subject

- `test(config): converge init integration fixtures`
