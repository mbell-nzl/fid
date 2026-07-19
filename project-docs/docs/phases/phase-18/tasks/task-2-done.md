# Phase 18 - Task 2: DB-backed workflow registry loader + v0.1 parity fixtures - Completion Report

## Summary of what was implemented

- Added a read-only DB-backed `workflow.Registry` loader for the four workflow tables.
- Decoded states, edges, ordered per-edge rules, and `params_json.roles` into the existing definition shape while retaining the frozen v0.1 map as the parity fixture.
- Added a structural `Unimplemented` marker for the two reserved rule names; no evaluator, journal append, DB write, consumer swap, or issues-schema change was added.

## Files created or modified

- `internal/workflow/registry.go`
- `internal/workflow/registry_test.go`
- `internal/workflow/workflow.go`
- `project-docs/docs/phases/phase-18/tasks/task-2-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-18/tasks/task-2-done.md`

## Validation Contract Coverage

- VAL-014: PASS — `TestRegistryEqualsV01Definitions` migrates a fresh SQLite database, loads all three seeded workflows, and deep-compares their IDs, states/categories, edges, ordered rules, and parameterised role lists against `FrozenDefinitions()`.
- VAL-015: PASS — `TestLoadSurfacesReservedRulesAsUnimplementedWithoutWriting` inserts test-only `required_gates_current` and `has_completion` rows, turns on SQLite `query_only`, and verifies that the loader returns both as ordered `RuleReference`s with `Unimplemented: true`. The read-only loader performs no evaluation or write.
- VAL-002/VAL-003/VAL-004: advanced — the same registry parity proof independently rechecks the seeded state/category, edge, and rule/role-parameter data against the frozen v0.1 reference.

## Tests run

- `go test ./internal/workflow -run '^TestRegistryEqualsV01Definitions$' -count=1` - expected RED: failed to build because `workflow.Load` was undefined.
- `go test ./internal/workflow -run '^TestRegistryEqualsV01Definitions$' -count=1` - PASS after implementation.
- `go test ./internal/workflow -run '^TestLoadSurfacesReservedRulesAsUnimplementedWithoutWriting$' -count=1` - expected RED: `RuleReference.Unimplemented` was not yet defined.
- `go test ./internal/workflow -run '^TestLoadSurfacesReservedRulesAsUnimplementedWithoutWriting$' -count=1` - PASS after implementation.
- `go test ./internal/workflow -count=1` - PASS.
- Runtime entrypoint exercise - N/A: the new registry is additive and no CLI/runtime consumer is repointed in this task by specification.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS (exit 0; no files listed).
- `golangci-lint run ./...` - PASS (`0 issues.`).
- `go test -race -coverprofile=coverage.out ./...` - PASS (exit 0).
- `go tool cover -func=coverage.out` - PASS (total statement coverage: 76.3%).
- `go test -tags integration -race ./...` - PASS (exit 0).

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `799acf0943be467b9f00424f7cb0ee8f63267dd7`

## Commit subject

- `feat(workflow): load workflow registry from database`
