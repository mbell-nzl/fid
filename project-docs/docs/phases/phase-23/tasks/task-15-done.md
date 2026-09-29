# Phase 23 - Task 15: Migration 0018 and V1 judgment removal - Completion Report

## Workflow outcome

- Task outcome: BLOCKED
- Category: specification-contradiction
- Condition: `schema.md` section 5 requires the byte-identical 0002 `concerns` definition in migration 0018 down. That definition contains `escalated_to_esc_id`, but the next rollback, migration 0016 down, adds `escalated_to_esc_id` to `concerns`. Implementing section 5 literally makes the full down chain fail with a duplicate column; omitting the column changes the accepted SQL authority.
- Recovery attempts: 1 session continuation after an external interruption; this attempt reproduced the missing-down failure and traced the conflicting DDL. The interrupted child returned no typed result or command termination, so none is treated as green evidence.
- Next action: Tech Lead must resolve the 0018-down `concerns` definition against migration 0016 down under protected schema authority. Then resume this same task from the committed partial tree, complete the V1 removal and VAL proofs, and run all nine unchanged required commands.
- Human attention required: yes

## Summary of what was implemented

- Preserved the interrupted task-15 changes: migration 0018 up, the generic Concern table rename through the initial readers and event catalogs, and the active migration version change to 18.
- Added `TestMigration0018DownRestoresEmptyV1Tables`. Its focused red result proves the missing down path. No down SQL was introduced because the accepted definition conflicts with the existing rollback chain.
- The partial tree is not a completed implementation. Existing production readers and retained tests still access removed tables; no required gate is claimed green.

## Files created or modified

- `migrations/0018_p23_remove_v1_judgment.up.sql`
- `internal/db/migration0018_test.go`, `internal/db/projection_meta.go`
- `internal/export/archive_semantics_catalog.go`
- `internal/journal/current_event_catalog.go`, `current_event_columns.go`, `current_event_json.go`, `current_event_result_projection.go`, `current_event_updates.go`, `current_event_validation.go`
- `internal/review/concern.go`, `concern_list.go`
- `internal/workflow/resume.go`, `resume_actions.go`
- `project-docs/docs/phases/phase-23/tasks/task-15-done.md`

## Documentation updated

- This report records the interrupted attempt, current red evidence, incomplete scope, and the protected 0018-down contradiction.

## Validation Contract Coverage

- VAL-004: BLOCKED - `TestMigration0018DownRestoresEmptyV1Tables` (`internal/db/migration0018_test.go:30`) fails because 0018 down is absent. The required post-0018 and post-rebuild pragmas have not been run.
- VAL-061: BLOCKED - `TestLatestSchemaHasNoV1JudgmentTables` (`internal/db/migration0018_test.go:8`) passed for its five table names, but the full inventory and real-process absence probes remain incomplete.
- VAL-062: BLOCKED - migration 0018 up matches the section 4 SQL in the preserved partial file, but no seeded 17-to-18 deletion or complete trigger probe has terminated in this session.
- VAL-063: BLOCKED - the touched-package test reported `check dependent row kinds: query task_completions: SQL logic error: no such table: task_completions (1)` from `internal/maintenance`; the maintenance checks have not yet been converged.
- VAL-055: BLOCKED - the requirement/decision activity producers and ten-name tagged inventory are absent from the partial tree.
- VAL-079: BLOCKED - the ordinary and tagged required gates and fresh-process benchmark have not run on a completed task-15 tree; package times cannot be claimed.

## Tests run

- `go test -run '^TestLatestSchemaHasNoV1JudgmentTables$' -count=1 ./internal/db` - PASS - `ok  	github.com/mbell-nzl/fid/internal/db	0.080s` (terminated `exit_code: 0`)
- `go test -run '^TestMigration0018DownRestoresEmptyV1Tables$' -count=1 ./internal/db` - FAIL - `FAIL` (terminated `exit_code: 1`); all five restored tables were absent because no 0018 down file exists.
- `go test -count=1 ./internal/db ./internal/issues ./internal/workflow ./internal/maintenance ./internal/review` - FAIL - `FAIL` (terminated `exit_code: 1`); the initial failure exposed the missing rollback and retained V1 readers/tests against the latest schema. This is immutable failed-attempt evidence, not a task gate.
- Prior externally interrupted child: no typed result, no observed command termination, and no accepted green result.

## Exact quality-gate results

Working directory: `/home/mark/projects/fid`. No required gate was started on a final tree because the protected schema contradiction prevents a correct 0018 down migration. Each unrun command is recorded as blocked without a fabricated output token or exit status.

- `gofmt -l cmd internal` - BLOCKED - not run; protected 0018-down contradiction leaves implementation incomplete.
- `golangci-lint run ./...` - BLOCKED - not run; protected 0018-down contradiction leaves implementation incomplete.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - not run; protected 0018-down contradiction leaves implementation incomplete.
- `go tool cover -func=coverage.out` - BLOCKED - not run; no final-tree coverage profile exists.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - not run; protected 0018-down contradiction leaves implementation incomplete.
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - BLOCKED - not run; protected 0018-down contradiction leaves implementation incomplete.
- `rg -n --glob '!*_test.go' "task_completions|ac_verifications" cmd internal` - BLOCKED - not run; production readers still exist in the partial tree, so its declared absence condition is unmet.
- `rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal` - BLOCKED - not run on a final tree; its declared expected `exit_code: 1` has not been observed.
- `rg -n --glob '!*_test.go' "has_completion_report|concerns_addressed" cmd internal` - BLOCKED - not run; the transition rule code still exists in the partial tree, so its declared absence condition is unmet.

## Deviations from spec

- No accepted contract was changed. The blocker is internal to the accepted down-migration instructions: `schema.md` section 5 says the 0018 down `concerns` definition is byte-identical to `migrations/0002_v01_tables.up.sql:125`, including `escalated_to_esc_id`; `migrations/0016_p22_remove_escalations.down.sql:13` then adds that same column. A task-scope choice between these definitions would silently revise the schema contract.
- The first focused rollback attempt ended `FAIL` (terminated `exit_code: 1`), with `restored V1 table "judgment_concerns" is missing` as a verbatim diagnostic. The root cause is the absent 0018 down file; a repair was withheld because the governing definition is contradictory. No later passing rollback result exists.

## Concerns, technical debt, or notes for the Tech Lead

- Resolve the exact 0018-down definition, including whether the restored V1 `concerns` table should match post-0016 schema without `escalated_to_esc_id`. The decision must also keep migration 0016 down executable, and should address whether its `ALTER TABLE concerns ADD COLUMN` or section 5's byte-identity requirement changes.
- The preserved partial migration removes V1 tables before the Issue and maintenance readers, CLI leaves, retained fixtures, and activity producers have converged. The touched-package command exposed these as unfinished task-15 work, not green evidence or another task's assignment.
- The external interruption consumed no no-progress budget per the Runner's TASK-RECOVERY binding. No provenance-based ownership claim is made for the retained tests before their introducing commits are checked during continuation.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `fa6355a0e221d97d7785b1e2718e134bb5609509`

## Commit subject

- `chore(db): checkpoint task 15 schema contradiction`
