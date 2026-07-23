# Phase 18 - Task 5: Remaining consumer swap and closure evidence - Completion Report

## Summary of what was implemented

- Completed the retained Task 5 consumer swap and closure evidence: maintenance, issue get/list/ready, and dev composites resolve workflow categories from the projected DB registry; the frozen Go map remains a parity fixture only.
- Added complete projection/repeated-rebuild, golden-v0.1 replay, journal-immutability, reserved-rule behavior-absence, and performance coverage for the Phase 18 closure assertions.
- Applied DEC-071: `rebuildWorkspaceLocked` now opens a rebuild-scoped SQLite connection with DSN pragmas `synchronous=OFF` and `journal_mode=MEMORY`, while preserving `foreign_keys=ON` and `recursive_triggers=ON`. Normal `db.Open` remains full-durability/delete-journal behavior.
- Applied DEC-072: `workflow.Load` now materializes the ordered workflow IDs and closes the outer query before loading each definition, preventing nested-query deadlock on a one-connection pool without changing registry contents, ordering, or normal DB pooling.

## Files created or modified

- `cmd/fid/e2e_v01_closure_helpers_test.go`, `cmd/fid/e2e_v01_closure_test.go`
- `internal/db/db.go`, `internal/db/db_test.go`
- `internal/dev/block.go`, `internal/dev/block_test.go`, `internal/dev/dev.go`, `internal/dev/finish.go`, `internal/dev/finish_test.go`, `internal/dev/gate.go`, `internal/dev/gate_test.go`, `internal/dev/start.go`, `internal/dev/start_test.go`
- `internal/issues/get.go`, `internal/issues/issue.go`, `internal/issues/list.go`, `internal/issues/ready.go`, `internal/issues/registry_consumers_test.go`, `internal/issues/registry_read_bench_test.go`, `internal/issues/transition_test.go`
- `internal/journal/projection.go`, `internal/journal/rebuild_bench_test.go`, `internal/journal/rebuild_scale_test.go`, `internal/journal/replay_test.go`
- `internal/maintenance/check.go`, `internal/maintenance/check_test.go`, `internal/maintenance/warnings.go`
- `internal/workflow/definitions_fixture.go`, `internal/workflow/registry.go`, `internal/workflow/registry_test.go`, `internal/workflow/workflow.go`
- `project-docs/docs/phases/phase-18/tasks/task-5-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-18/tasks/task-5-done.md` - corrected the DEC-072 remediation evidence with the fresh terminating phase gate, coverage profile, benchmark metrics, focused maintenance evidence, and VAL-024 inventory.

## Validation Contract Coverage

- VAL-010 - PASS. Fresh `go test -count=1 -v ./cmd/fid -run '^(TestMaintenanceCheckCommandRendersHealthyAndFailingResults|TestMaintenanceCheckCommandPromotesWarningsWithoutRewritingThem)$'` passed: it proves a healthy maintenance result, a seeded-invalid workflow pair finding, and `--warn-as-error` promotion to a non-zero result without changing the warning payload. The built real entrypoint also passed in `/tmp/fid-phase18-task5-remediation.iT4wj1`: after `git init`, `/tmp/fid init`, and `/tmp/fid rebuild --by workflow_runner`, `/tmp/fid maintenance check --format json` returned `{"status":"pass","errors":[],"warnings":[],"info":{"issue_count":0,"claim_count_live":0,"recent_activity_7d":0}}`.
- VAL-011 - PASS. `TestV01ClosureJourneyAndRebuildParity` compares every projected table, schema, and FK state before delete-and-rebuild and after a second rebuild.
- VAL-013 - PASS. `TestReplayJournalAppliesGoldenV1JournalInFileOrder` replays the golden v0.1 journal and verifies the migrated workflow/state/category projection.
- VAL-017 - PASS. `TestTransitionReservedGateRuleRemainsUnevaluated` proves the injected reserved rule remains behavior-free; the final production inventory found no receipt/gate evaluator or query.
- VAL-021 - PASS. `TestV01ClosureJourneyAndRebuildParity` compares journal files byte-for-byte across migration, stale-version refresh, delete-and-rebuild, and second rebuild.
- VAL-022 - PASS. `TestGetListAndReadyUseProjectedWorkflowCategories` plus the full closure CLI journey prove DB-registry category presentation and ready selection/order.
- VAL-023 - PASS. The projected-registry dev start, finish, block, and gate tests prove the composites retain their v0.1 outcomes/effects; the full closure CLI journey remains green.
- VAL-024 - PASS. `rg -n 'definitions' internal cmd --glob '*.go'` returns only the frozen fixture and test references; `rg -n 'workflow\\.' internal cmd --glob '*.go' --glob '!**/*_test.go'` shows DB-backed `workflow.Load` consumer paths and rule types only, paired with the VAL-022/023 runtime evidence.

## Tests run

- Red (DEC-072 reproduction): `go test -race -count=1 -v -timeout 5s ./internal/maintenance -run '^TestCheckWorkflowStatesReadsProjectedWorkflowTables$'` - expected FAIL before the fix: `panic: test timed out after 5s`, blocked in `database/sql.(*DB).conn` from `workflow.loadStates` via `workflow.Load`.
- Red (new bounded regression): `go test -count=1 -v ./internal/workflow -run '^TestLoadReleasesWorkflowRowsBeforeNestedQueries$'` - expected FAIL before the fix in 1.01s: `Load() did not return with one open database connection`.
- Green: `go test -count=1 -v ./internal/workflow -run '^TestLoadReleasesWorkflowRowsBeforeNestedQueries$'` - PASS (0.02s); `Load` returns with `SetMaxOpenConns(1)`.
- Green: `go test -race -count=1 -v -timeout 5s ./internal/maintenance -run '^TestCheckWorkflowStatesReadsProjectedWorkflowTables$'` - PASS (1.234s), proving the formerly hung maintenance path completes under the original one-connection fixture.
- Green: `go test -count=1 -v ./internal/workflow -run '^TestRegistryEqualsV01Definitions$'` - PASS (0.02s), preserving frozen v0.1 registry contents.
- Red: `go test ./internal/db -run '^TestOpenForRebuildRelaxesDurabilityWithoutChangingNormalOpen$' -count=1` - expected FAIL before implementation: `undefined: OpenForRebuild`.
- Green: `go test ./internal/db -run '^TestOpenForRebuildRelaxesDurabilityWithoutChangingNormalOpen$' -count=1` - PASS. It proves rebuild pragmas (`foreign_keys=1`, `recursive_triggers=1`, `synchronous=0`, `journal_mode=memory`) and normal `Open` (`foreign_keys=1`, `recursive_triggers=1`, `synchronous=2`, `journal_mode=delete`).
- `go test ./internal/journal -run '^(TestReplayJournalAppliesGoldenV1JournalInFileOrder|TestRebuildWorkspaceRecreatesIdenticalStateAfterDeletingDatabase)$' -count=1` - PASS.
- `go test -tags integration ./cmd/fid -run '^TestV01ClosureJourneyAndRebuildParity$' -count=1` - PASS (5.750s).
- `go test -count=1 -v ./cmd/fid -run '^(TestMaintenanceCheckCommandRendersHealthyAndFailingResults|TestMaintenanceCheckCommandPromotesWarningsWithoutRewritingThem)$'` - PASS. Healthy, seeded-invalid, and `--warn-as-error` non-zero maintenance evidence for VAL-010.
- `go test -count=1 -v ./internal/issues -run '^(TestGetListAndReadyUseProjectedWorkflowCategories|TestTransitionReservedGateRuleRemainsUnevaluated)$'` - PASS.
- `go test -count=1 -v ./internal/dev -run '^(TestStartRejectsEdgeMissingFromProjectedWorkflowRegistry|TestFinishRejectsEdgeMissingFromProjectedWorkflowRegistry|TestBlockRejectsEdgeMissingFromProjectedWorkflowRegistry|TestGateUsesProjectedWorkflowCategory)$'` - PASS.
- `go test -tags integration -count=1 -v ./cmd/fid -run '^TestV01ClosureJourneyAndRebuildParity$'` - PASS (5.844s), preserving full projection and journal parity.
- Final VAL-024 inventory commands - PASS. The two exact `rg` commands cited under VAL-024 produced only approved frozen-fixture/test map references and DB-backed production registry paths; the focused issue, dev, and closure tests above provide runtime parity.
- Real entrypoint: in `/tmp/fid-phase18-task5-remediation.iT4wj1`, ran `git init`, then `/tmp/fid init`, `/tmp/fid rebuild --by workflow_runner`, and `/tmp/fid maintenance check --format json` - PASS. Rebuild replayed 1 event in 16ms and maintenance returned `{"status":"pass","errors":[],"warnings":[],"info":{"issue_count":0,"claim_count_live":0,"recent_activity_7d":0}}`.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS (exit 0; no output).
- `golangci-lint run ./...` - PASS (exit 0; `0 issues.`).
- `go test -race -coverprofile=coverage.out ./...` - PASS (exit 0; fresh `coverage.out` produced). Key package results: `cmd/fid` 77.7%, `internal/issues` 81.5%, `internal/claims` 81.4%, `internal/journal` 81.0%, `internal/policy` 86.3%, and `internal/workflow` 83.2%.
- `go tool cover -func=coverage.out` - PASS (exit 0; fresh total coverage `76.8%`).
- `go test -tags integration -race ./...` - PASS (exit 0; `cmd/fid` 106.462s and all packages terminated successfully).
- `go test -bench . -benchmem ./...` - PASS (exit 0). Required benchmark rows:
  - `BenchmarkList10kIssuesWithRegistry-6        22  47590455 ns/op        47.59 ms/registry-read  9841780 B/op  191571 allocs/op`
  - `BenchmarkRebuildWorkspace1kEvents-6          4  290188734 ns/op       290.2 ms/rebuild      18868438 B/op  339379 allocs/op`
  - `BenchmarkRebuildWorkspace10kEvents-6         1  2825204196 ns/op      2825 ms/rebuild        176844800 B/op  3384977 allocs/op`
  - `BenchmarkRebuildWorkspace100kEvents-6        1  26986552792 ns/op     26987 ms/rebuild       1756626200 B/op 33842237 allocs/op`
- `go build -o /tmp/fid ./cmd/fid` - PASS (exit 0; compiled binary used for final real-entrypoint evidence).

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
