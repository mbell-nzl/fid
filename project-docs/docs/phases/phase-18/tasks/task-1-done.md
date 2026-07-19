# Phase 18 - Task 1: Workflow tables migration + seed defaults + gate-obligation seam - Completion Report

## Summary of what was implemented

- Added forward/down migration 0007 with the four v0.2 workflow tables and the three active v0.1 default workflows as static reference data.
- Seeded exact state/category mappings, allowed transitions with explicit non-cancelled-to-cancelled edges, and existing per-edge rules/role parameters.
- Added frozen v0.1 definitions for migration-parity tests, reserved the two future gate-obligation rule names without evaluators, and bumped the projection schema version from 3 to 4.

## Files created or modified

- `migrations/0007_v02a_workflow_tables.up.sql`
- `migrations/0007_v02a_workflow_tables.down.sql`
- `internal/db/projection_meta.go`
- `internal/db/migrate_v01_tables_test.go`
- `internal/db/workflow_tables_test.go`
- `internal/workflow/definitions_fixture.go`
- `internal/workflow/workflow.go`
- `internal/workflow/rules.go`

## Documentation updated

- None required.

## Validation Contract Coverage

- VAL-001: PASS — `TestMigration0007CreatesWorkflowSchema` verifies all four tables, specified columns, the three prescribed PKs, the composite requirement-to-transition FK, and absence of a PK/UNIQUE on `workflow_transition_requires`.
- VAL-002: PASS — `TestMigration0007SeedsFrozenV01WorkflowDefinitions` deep-compares every seeded state/category tuple with `workflow.FrozenDefinitions`; `TestWorkflowDefinitionsMatchTRDExactly` preserves the fixture-to-TRD characterization.
- VAL-003: PASS — the same seed-parity test compares the complete transition tuple set, including every explicit non-cancelled-to-cancelled edge and absence of `cancelled` to `cancelled`.
- VAL-004: PASS — the seed-parity test compares every requirement tuple, including canonical role `params_json` values and null parameters for non-parameterised rules.
- VAL-005: PASS — `TestMigration0007ReservesGateObligationShapeWithoutSeedingReservedRules` inserts, loads, parses, and marshals `{"gates":[{"id":"gate-test","version":1}]}` without loss.
- VAL-012: PASS — `TestDecideProjectionFreshnessRequiresPhase18SchemaVersion` verifies prior version 3 rebuilds and current version 4 remains fresh; `TestEnsureProjectionFreshRebuildsStoredPriorProjectionSchemaVersion` exercises a real rebuild/catch-up and stores the current version.
- VAL-016: PASS — `TestMigration0007ReservesGateObligationShapeWithoutSeedingReservedRules` queries both reserved names and finds zero seeded rows.

## Tests run

- `go test ./internal/db -run '^TestMigrateUpSeedsTaskDefaultWorkflowStates$' -count=1` - expected RED: failed with `no such table: workflow_states` before migration 0007.
- `gofmt -w internal/db/workflow_tables_test.go internal/workflow/definitions_fixture.go internal/workflow/workflow.go internal/workflow/rules.go && go test ./internal/db -run '^TestMigrateUpSeedsTaskDefaultWorkflowStates$' -count=1` - PASS.
- `go test ./internal/db ./internal/workflow ./internal/journal -count=1` - PASS.
- `/tmp/fid-phase18 init --prefix P18` in disposable `/tmp/fid-phase18-cli-9YKAd9` workspace - PASS; initialized the workspace through the real CLI entrypoint.
- `/tmp/fid-phase18 sql "SELECT name FROM sqlite_master WHERE type = 'table' AND name LIKE 'workflow%' ORDER BY name"` - PASS; returned `workflow_states`, `workflow_transition_requires`, `workflow_transitions`, and `workflows`.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS (exit 0; no files listed).
- `golangci-lint run ./...` - PASS (`0 issues.`).
- `go test -race -coverprofile=coverage.out ./...` - PASS.
- `go tool cover -func=coverage.out` - PASS (total statement coverage: 76.3%).
- `go test -tags integration -race ./...` - PASS.

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
