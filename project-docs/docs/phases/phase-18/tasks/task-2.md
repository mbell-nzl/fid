# Phase 18 - Task 2: DB-backed workflow registry loader + v0.1 parity fixtures

## Objective

Add a DB-backed workflow registry that loads workflow definitions from the
seeded tables (Task 1) into the existing `workflow.Definition` shape, and prove
it is behavior-equivalent to the frozen v0.1 Go definitions. This is additive:
no consumer is repointed yet, so behavior is unchanged.

## Validation Contract Coverage

- VAL-014: the DB-backed registry loaded from the seeded tables is behavior-equivalent to the frozen v0.1 Go definitions (states, categories, edges, rules) for all three workflows.
- VAL-015: the loader recognizes `required_gates_current` and `has_completion` as known-but-unimplemented rule names — surfaced structurally with no panic, no evaluation, and no write.
- Advances VAL-002/003/004 (the loader reading the seeded rows is a second, independent check that the seed matches v0.1).

## Binding Constraints

- **Keep the pure logic; swap only the source.** Add a loader
  (e.g. `workflow.Load(conn) (*Registry, error)` or `LoadDefinition(conn,
  workflowID)`) that reads `workflows`/`workflow_states`/
  `workflow_transitions`/`workflow_transition_requires` and builds the same
  `Definition`/`RuleReference` structures the package already uses. Reuse the
  existing pure helpers (`RulesForEdge`, `CategoryForState`, `IsAllowedEdge`)
  by operating them on the loaded definition rather than the package-global map.
- **Do not repoint any consumer in this task.** `internal/issues/transition.go`
  and every other caller keep using the current package-level functions (Go
  map) until Task 3+. This task only introduces the loader and its equality
  proof, so the tree stays green with behavior unchanged.
- **Retain the v0.1 Go `definitions` map as frozen parity fixtures.** It is the
  reference the registry is proven against; it is not deleted in this phase.
- **Parse `params_json` back into `RuleReference.RequiredRoles`** for
  parameterised rules; ignore/return-untouched any reserved gate-obligation
  params (they are never seeded in Phase 18). Unknown/unimplemented rule names
  must not crash the loader — surface them structurally without evaluating them.
- The registry loader is read-only over the projection; it must not append
  journal events or write to the DB.

## Scope

### Files To Create Or Modify

- `internal/workflow/registry.go` (or similar) — the DB-backed loader + registry type
- `internal/workflow/registry_test.go` — equality/characterization tests vs frozen v0.1 definitions
- Minimal additions to `internal/workflow/workflow.go` if methods must operate
  on a loaded `Definition` (do not change existing exported behavior)

### Files To Read For Context

- `internal/workflow/workflow.go`, `internal/workflow/rules.go`
- `migrations/0007_v02a_workflow_tables.up.sql` (Task 1 seed)
- `internal/db/db.go` (connection open / pragmas)
- `internal/issues/transition_rules.go` (how rules are consumed — for shape parity)

### Out Of Scope

- Repointing `transition.go` or any other consumer (Tasks 3, 5).
- Any `issues` schema change (Task 4).
- Removing the Go `definitions` map (kept as parity fixtures this phase).

## Acceptance Criteria

- [ ] `workflow.Load`/registry reads the four seeded tables and builds
  `Definition`s for `epic-default`, `task-default`, `subtask-default`
- [ ] Registry states, categories, allowed edges, and per-edge rules
  (incl. parameterised roles) equal the frozen v0.1 Go definitions for all
  three workflows (VAL-014)
- [ ] Loader parses `params_json` roles back into `RequiredRoles`; reserved
  gate-obligation params and unimplemented rule names do not crash the loader
- [ ] Loading a test-only row carrying `required_gates_current` or
  `has_completion` surfaces the name as known-but-unimplemented with no panic,
  no evaluation, and no write (VAL-015)
- [ ] No consumer behavior changes; full suite stays green

## First Failing Test

`TestRegistryEqualsV01Definitions` (or similar): load `task-default` from the
DB and assert its edges/categories/rules equal the frozen Go fixture — fails
until the loader exists.

## Documentation Updates

- None expected. Note any loader-shape decision in the completion report.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates.

Conditional command required by this task:

- `go test -tags integration -race ./...` — keeps migration/rebuild paths green
  while the loader is introduced.

## Notes For Developer

- Depends on Task 1 (seeded tables). The loader is the de-risking step before
  the transition-path swap in Task 3; keep it small and thoroughly characterized.
- Do not change error messages or evaluation semantics — this task must be a
  no-op for runtime behavior.
- Completion report at `project-docs/docs/phases/phase-18/tasks/task-2-done.md`.
