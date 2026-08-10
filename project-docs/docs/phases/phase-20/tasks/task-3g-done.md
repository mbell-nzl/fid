# Phase 20 - Task 3g: Reuse one integration binary - Completion Report

## Summary of what was implemented

- Added a package-owned `TestMain` lifecycle that builds one immutable `fid`
  binary, reports build output on failure, runs the package, and removes its
  temporary binary root.
- Replaced all four journey build calls with a checked read-only accessor while
  preserving four distinct `t.TempDir()` workspaces and real subprocesses.
- Extended the source classifier to reject any per-test build helper and require
  `TestMain`; improved its journey scan to recognize real `*testing.T`
  signatures rather than every `Test`-prefixed function.
- Reduced a fresh dedicated race run from Task 3f's 10.266s to 8.228s.

## Files created or modified

- `cmd/fid/integration/main_test.go`
- `cmd/fid/integration/classification_test.go`
- `cmd/fid/integration/e2e_v01_closure_helpers_test.go`
- `cmd/fid/integration/e2e_v01_closure_test.go`
- `cmd/fid/integration/e2e_v02a_workflow_authoring_test.go`
- `cmd/fid/integration/v01_parity_test.go`
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/tasks/task-3g-done.md`

## Documentation updated

- `infra-improvement-pack.md` records the one-build source inventory, four
  isolated workspaces, and before/after dedicated-package timings.

## Validation Contract Coverage

- No VAL-NNN applies. Existing journey assertions and subprocess behavior are
  unchanged. `TestIntegrationBinaryOwnedByPackage` and the four workspace
  source matches provide delivery-harness evidence.

## Tests run

- `go test -tags integration -run '^TestIntegrationBinaryOwnedByPackage$' ./cmd/fid/integration` - EXPECTED RED - `FAIL`.
- `go test -tags integration -run '^TestIntegration(Classification|BinaryOwnedByPackage)$' ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.409s`.
- Source build/workspace inventory - PASS - `workspace_inventory_status=0`.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`.
- `golangci-lint run ./...` - PASS - `0 issues.`.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.5% of statements`.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					77.7%`.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS on byte-identical rerun - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.228s`.
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)`.

## Deviations from spec

- The first full dedicated integration gate failed with `FAIL`. Root cause:
  Task 3f's source classifier identified every `Test`-prefixed function as a
  journey, so the new harness test and `TestMain(*testing.M)` appeared as two
  unexpected journeys. The classifier now recognizes only `*testing.T` test
  signatures and explicitly excludes its two infrastructure tests. No journey,
  product assertion, command, tag, or validator was removed or narrowed; the
  byte-identical full gate passed deterministically after the harness repair.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `1edb689b2df560a1f7ce60875458a3e639e60910`

## Commit subject

- `test(integration): reuse one package binary`
