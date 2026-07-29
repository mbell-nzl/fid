# Phase 19 - Task 1: `workflow_assignments` schema + `escalation-default` seed - Completion Report

## Summary of what was implemented

- Added reversible migrations `0009` for journal-projected workflow-assignment storage and `0010` for the exact `escalation-default` system reference workflow.
- Added the `linked_to_decision` known-but-unimplemented loader rule and bumped the projection schema version to 5.
- Preserved the frozen Phase 18 defaults and updated the schema/TRD documentation for the new table, seed, vocabulary, and future assignment precedence.
- This hash-bearing report is finalized in the immediately following `docs(task)` commit.

## Files created or modified

- `migrations/0009_v02a_workflow_assignments.up.sql`
- `migrations/0009_v02a_workflow_assignments.down.sql`
- `migrations/0010_v02a_escalation_default_workflow.up.sql`
- `migrations/0010_v02a_escalation_default_workflow.down.sql`
- `internal/db/workflow_assignments_test.go`
- `internal/db/escalation_workflow_test.go`
- `internal/db/workflow_tables_test.go`
- `internal/db/migrate_v01_tables_test.go`
- `internal/db/projection_meta.go`
- `internal/workflow/rules.go`
- `internal/workflow/registry.go`
- `internal/workflow/registry_test.go`

## Documentation updated

- `project-docs/specs/03_TRD_SCHEMA.md` — assignment DDL/indexes, `escalation` applicability, and the known-but-unimplemented rule set.
- `project-docs/specs/02_TRD.md` — migration `0010` seed and DEC-076 assignment-resolution precedence.

## Validation Contract Coverage

- VAL-001: PASS - `TestMigration0009WorkflowAssignmentsSchemaAndConstraints` introspects columns, CHECK text, FK, and both exact partial-index definitions.
- VAL-027: PASS - `TestMigration0009CreatesWorkflowAssignmentsWithScopedUniqueness` rejects the duplicate `task` kind row and accepts the distinct `epic` kind row.
- VAL-028: PASS - `TestMigration0009WorkflowAssignmentsSchemaAndConstraints` accepts one project row and rejects the second.
- VAL-029: PASS - `TestMigration0009WorkflowAssignmentsSchemaAndConstraints` has two malformed rejected cases and two well-formed accepted cases.
- VAL-002: PASS - `TestMigrations0009And0010RoundTripWithoutResidue` deep-compares pre-up, down, and second-up `sqlite_schema` and every table-row snapshot.
- VAL-003: PASS - `TestMigration0010SeedsEscalationDefaultExactlyAndPreservesV01Tuples` compares the four ordered escalation table tuple sets to the literal TRD §4.6 values.
- VAL-030: PASS - `TestMigration0010SeedsEscalationDefaultExactlyAndPreservesV01Tuples` captures all four default-table tuple sets before/after `0010`; `TestRegistryEqualsV01Definitions` compares each loaded frozen default with `FrozenDefinitions()`.
- VAL-004: PASS - `TestLoadSurfacesEscalationRuleAsUnimplementedWithoutWriting` loads `escalation-default` under `PRAGMA query_only` and asserts `RuleLinkedToDecision` has `Unimplemented: true`.
- VAL-031: PASS (no match) - exit 1 (no output) — the expected no-match result for this inventory per VAL-031 / QA-36 / DEC-094; `TestLoadSurfacesEscalationRuleAsUnimplementedWithoutWriting` succeeds under SQLite `query_only`, and the production inventory `rg -n --glob '*.go' --glob '!**/*_test.go' 'RuleLinkedToDecision|linked_to_decision' internal/issues cmd/fid` returned no matches.
- VAL-005: PASS - `TestDecideProjectionFreshnessRequiresPhase19SchemaVersion` pins version 5, asserts stored version 4 rebuilds, and asserts a current clean version remains fresh.

## Tests run

- `go test ./internal/db -run '^TestMigration0009CreatesWorkflowAssignmentsWithScopedUniqueness$' -count=1` - FAIL (expected red) - `FAIL`
- `go test ./internal/db -run '^TestMigration0009CreatesWorkflowAssignmentsWithScopedUniqueness$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/db 0.021s`
- `go test ./internal/workflow -run '^TestLoadSurfacesEscalationRuleAsUnimplementedWithoutWriting$' -count=1` - FAIL (expected red) - `FAIL`
- `go test ./internal/workflow -run '^TestLoadSurfacesEscalationRuleAsUnimplementedWithoutWriting$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow 0.021s`
- `go test ./internal/db -run '^TestDecideProjectionFreshnessRequiresPhase19SchemaVersion$' -count=1` - FAIL (expected red) - `FAIL`
- `go test ./internal/db -run '^TestDecideProjectionFreshnessRequiresPhase19SchemaVersion$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/db 0.004s`
- `go test ./internal/db ./internal/workflow -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow 0.058s`
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'RuleLinkedToDecision|linked_to_decision' internal/issues cmd/fid` - PASS (no match) - exit 1 (no output) — exit 1 is the expected no-match result for this inventory per VAL-031 / QA-36 / DEC-094.
- QA-plan applicability: QA-01..05 and QA-32..36 apply to this task; their deterministic evidence is supplied by the cited Task 1 tests and mandatory gates.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`
- `golangci-lint run ./...` - PASS - `0 issues.`
- `go test -race -coverprofile=coverage.out ./...` - PASS - `exit 0 (no output)`
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 76.6%`
- `go test -tags integration -race ./...` - PASS - `exit 0 (no output)`
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)`

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- Correction (TRG-19-001 / DEC-094): The first finalized report recorded the inventory as `exit 0 (no output)`; `rg` exits 1 on no match, so that token was not copied from the command. The command was re-executed and corrected under TRG-19-001 / DEC-094; implementation commit `a509d59` is unchanged.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `a509d59efcc826df417f64ea72a43d743936e669`

## Commit subject

- `feat(workflow): add assignment schema and escalation seed`
