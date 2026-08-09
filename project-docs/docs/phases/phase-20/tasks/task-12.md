# Phase 20 - Task 12: Workflow command contract

## Objective

Migrate workflow define, assign, and list to their exact Response Envelope
contracts and typed error mapping.

## Validation Contract Coverage

- VAL-008, VAL-017, VAL-026, VAL-041, the workflow rows of VAL-044, and the
  workflow ordering/limit cases of VAL-045.

## Binding Constraints

- Remove canonical `--require` authoring and its role/assurance-shaped request
  parser. Phase 20 workflow definitions contain states and transitions only.
- Bootstrap-seeded V1 requirement rows remain internal compatibility state for
  legacy `fid dev` until Phases 21/23; do not expose them in canonical results,
  let retained commands author them, or rename them into future semantics.
- A stored unknown/unsupported transition precondition still rejects an Issue
  transition before append; it is never ignored.
- Implement exact WorkflowSummary and assignment/definition result schemas.
- Preserve Task 2's `--project-wide` collision fix; exactly one of it and
  `--kind` is required. Global `--project <path>` remains the root selector.
- Kind assignment deterministically takes precedence over the project default;
  duplicate assignments never create ambiguous selection.
- Slugs, states, list sizes, categories, and edges follow the exact closed value
  domains in contract section 3.2.1.
- Unknown workflow or scope maps to `workflow-not-found`; duplicates and invalid
  definitions map to the appropriate existing closed code without driver-message parsing.
- Remove this family’s `--format` selector.
- Read the resolved project root from Task 2 command context.

## Scope

### Files To Create Or Modify

- `cmd/fid/workflow.go` and define/assign/list tests.
- `internal/workflow` result objects and typed error paths required by the adapter.

### Files To Read For Context

- Phase 19 DEC-078/080; `../contracts.md` sections 2.4 and 3.

### Out Of Scope

- Generic precondition authoring, execution-bound transitions, generic
  authority, internal V1 requirement-row removal, and workflow schema changes.

## Acceptance Criteria

- [ ] Define, assign, and list accept only closed requests and match exact schemas/Event/null rules. (VAL-017, VAL-026, VAL-044)
- [ ] Project/kind assignment nullability, list order/bound, and kind-over-project precedence are deterministic. (VAL-026, VAL-045)
- [ ] Global `--project <path>` combines with Task 2's `--project-wide` without
      flag shadowing; the removed local boolean spelling remains unregistered. (VAL-026, VAL-044)
- [ ] `--require` is unregistered; no role/assurance rule or requirement count
      appears in canonical requests/results. (VAL-026, VAL-044)
- [ ] Slug/state/category/edge boundaries, duplicates, and invalid cases return
      typed failures without append or historical version promises. (VAL-026, VAL-044)
- [ ] Attribution is recorded without trusted/protected effects; Human success is stdout-only. (VAL-008, VAL-041)

## First Failing Test

`cmd/fid/workflow_define_test.go` - `TestWorkflowDefineJSONMatchesContract`.

## Documentation Updates

- None expected.

## Exact Commands To Run

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

## Notes For Developer

- Completion report: `task-12-done.md`.
