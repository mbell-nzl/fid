# Phase 18 - Task 5: Remaining consumer swap and closure evidence - Completion Report

## Summary of what was implemented

- Completed the retained Task 5 consumer swap and closure evidence: maintenance, issue get/list/ready, and dev composites resolve workflow categories from the projected DB registry; the frozen Go map remains a parity fixture only.
- Added complete projection/repeated-rebuild, golden-v0.1 replay, journal-immutability, reserved-rule behavior-absence, and performance coverage for the Phase 18 closure assertions.
- Applied DEC-071: `rebuildWorkspaceLocked` now opens a rebuild-scoped SQLite connection with DSN pragmas `synchronous=OFF` and `journal_mode=MEMORY`, while preserving `foreign_keys=ON` and `recursive_triggers=ON`. Normal `db.Open` remains full-durability/delete-journal behavior.

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

- `project-docs/docs/phases/phase-18/tasks/task-5-done.md` - replaced the DEC-070/071 blocker report with this completion report.

## Validation Contract Coverage

- VAL-010 - PASS. `TestCheckWorkflowStatesReadsProjectedWorkflowTables` proves a healthy DB is clean and a test-only invalid projected workflow pair is reported.
- VAL-011 - PASS. `TestV01ClosureJourneyAndRebuildParity` compares every projected table, schema, and FK state before delete-and-rebuild and after a second rebuild.
- VAL-013 - PASS. `TestReplayJournalAppliesGoldenV1JournalInFileOrder` replays the golden v0.1 journal and verifies the migrated workflow/state/category projection.
- VAL-017 - PASS. `TestTransitionReservedGateRuleRemainsUnevaluated` proves the injected reserved rule remains behavior-free; the final production inventory found no receipt/gate evaluator or query.
- VAL-021 - PASS. `TestV01ClosureJourneyAndRebuildParity` compares journal files byte-for-byte across migration, stale-version refresh, delete-and-rebuild, and second rebuild.
- VAL-022 - PASS. `TestGetListAndReadyUseProjectedWorkflowCategories` plus the full closure CLI journey prove DB-registry category presentation and ready selection/order.
- VAL-023 - PASS. The projected-registry dev start, finish, block, and gate tests prove the composites retain their v0.1 outcomes/effects; the full closure CLI journey remains green.
- VAL-024 - PASS. `rg -n 'definitions' internal cmd --glob '*.go'` returns only the frozen fixture and test references; `rg -n 'workflow\\.' internal cmd --glob '*.go' --glob '!**/*_test.go'` shows DB-backed `workflow.Load` consumer paths and rule types only, paired with the VAL-022/023 runtime evidence.

## Tests run

- Red: `go test ./internal/db -run '^TestOpenForRebuildRelaxesDurabilityWithoutChangingNormalOpen$' -count=1` - expected FAIL before implementation: `undefined: OpenForRebuild`.
- Green: `go test ./internal/db -run '^TestOpenForRebuildRelaxesDurabilityWithoutChangingNormalOpen$' -count=1` - PASS. It proves rebuild pragmas (`foreign_keys=1`, `recursive_triggers=1`, `synchronous=0`, `journal_mode=memory`) and normal `Open` (`foreign_keys=1`, `recursive_triggers=1`, `synchronous=2`, `journal_mode=delete`).
- `go test ./internal/journal -run '^(TestReplayJournalAppliesGoldenV1JournalInFileOrder|TestRebuildWorkspaceRecreatesIdenticalStateAfterDeletingDatabase)$' -count=1` - PASS.
- `go test -tags integration ./cmd/fid -run '^TestV01ClosureJourneyAndRebuildParity$' -count=1` - PASS (5.750s).
- `go test ./internal/maintenance -run '^TestCheckWorkflowStatesReadsProjectedWorkflowTables$' -count=1` - PASS.
- `go test ./internal/issues -run '^(TestGetListAndReadyUseProjectedWorkflowCategories|TestTransitionReservedGateRuleRemainsUnevaluated)$' -count=1` - PASS.
- `go test ./internal/dev -run '^(TestStartRejectsEdgeMissingFromProjectedWorkflowRegistry|TestFinishRejectsEdgeMissingFromProjectedWorkflowRegistry|TestBlockRejectsEdgeMissingFromProjectedWorkflowRegistry|TestGateUsesProjectedWorkflowCategory)$' -count=1` - PASS.
- Final VAL-024 inventory commands - PASS. The two exact `rg` commands cited under VAL-024 produced only approved frozen-fixture/test map references and DB-backed production registry paths.
- Real entrypoint: in `/tmp/fid-phase18-task5-entrypoint.TtLDzj`, ran `git init`, then `/tmp/fid init`, `/tmp/fid rebuild --by workflow_runner`, and `/tmp/fid maintenance check --format json` - PASS. Rebuild replayed 1 event in 15ms and maintenance returned `{"status":"pass","errors":[],"warnings":[]...}`.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS (exit 0; no output).
- `golangci-lint run ./...` - PASS (exit 0; `0 issues.`).
- `go test -race -coverprofile=coverage.out ./...` - PASS (exit 0; `coverage.out` produced).
- `go tool cover -func=coverage.out` - PASS (exit 0; total coverage `74.1%`).
- `go test -tags integration -race ./...` - PASS (exit 0).
- `go test -bench . -benchmem ./...` - PASS (exit 0). Required benchmark rows:
  - `BenchmarkList10kIssuesWithRegistry-6    25  46807508 ns/op        46.81 ms/registry-read`
  - `BenchmarkRebuildWorkspace1kEvents-6      4  278769882 ns/op       278.8 ms/rebuild`
  - `BenchmarkRebuildWorkspace10kEvents-6     1  2647464126 ns/op      2647 ms/rebuild`
  - `BenchmarkRebuildWorkspace100kEvents-6    1  25235973633 ns/op     25236 ms/rebuild`
- `go build -o /tmp/fid ./cmd/fid` - PASS (exit 0; compiled binary used for final real-entrypoint evidence).

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `37b3fb524bc5ba3d2774f6c14a12e78cbc028e1f`

## Commit subject

- `perf(rebuild): relax projection durability`
