# Phase 20 - Task 7: Issue hierarchy invariant

## Objective

Enforce legal Epic/Task/Subtask parent relationships in the Issue domain before
any Event append, with a typed `hierarchy-invalid` failure.

## Validation Contract Coverage

- VAL-020.

## Binding Constraints

- Legal shape: Epic has no parent; Task may have no parent or an Epic parent;
  Subtask requires a Task parent. No deeper nesting is introduced.
- Parent is assigned only when creating a new Issue and must reference an
  existing Issue; no update/transition flag changes parent. Together with the
  hierarchy rule, parent cycles are therefore unrepresentable. Preserve and
  test that absence rather than adding cycle code for an impossible mutation.
- Domain validation is primary; SQLite triggers remain corruption backstops.
- Missing parent uses `issue-not-found`; an existing but illegal parent uses
  `hierarchy-invalid`. Both reject before append.

## Scope

### Files To Create Or Modify

- `internal/issues/create.go`, relevant request inventory, typed errors and tests.
- Shared CLI error mapping tests only where needed to prove the code.

### Files To Read For Context

- PRD-ISSUE-002/004; kind triggers; `../contracts.md` section 2.4.

### Out Of Scope

- Dependency links/cycles, normalized dependency storage, and new parent-update behavior.

## Acceptance Criteria

- [ ] Epic-with-parent, Subtask-under-Epic, and Subtask-without-parent reject without append. (VAL-020)
- [ ] A missing parent rejects as `issue-not-found` without append. (VAL-020)
- [ ] Task-under-Epic, root Task, and Subtask-under-Task succeed. (VAL-020)
- [ ] Existing-but-illegal parents map to `hierarchy-invalid` with exact details. (VAL-020)
- [ ] No Issue mutation surface can change parent after creation, preserving the
      construction-time proof that parent cycles cannot be formed. (VAL-020)

## First Failing Test

`internal/issues/create_test.go` - `TestCreateRejectsIllegalHierarchyBeforeAppend`.

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

- Completion report: `task-7-done.md`.
