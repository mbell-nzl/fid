# Phase 20 - Task 8: Issue read contract

## Objective

Migrate `issue get`, `issue list`, and `ready` to their exact Response Envelope
schemas and delete the Markdown-oriented `issue get --section` surface.

## Validation Contract Coverage

- VAL-008, VAL-013, VAL-014, VAL-015, VAL-016, the read half of VAL-017,
  read rows of VAL-044, and relevant ordering/limit cases in VAL-045.

## Binding Constraints

- Implement `../contracts.md` sections 3.1-3.3, 3.5, and 3.7 exactly; no provisional
  owner/concern fields or `latest_verifications` appear.
- Required nullable keys are present; arrays are empty, never null.
- `ready` returns exactly one non-null member of `task` and `explanation`.
- Ready explanations use only the two closed codes and their code-specific
  nullable fields; no free-form code reaches canonical JSON.
- Remove `--format` from these leaves; root `--json` activates their envelope.
- Read the resolved project root from Task 2 command context; do not reintroduce
  leaf-local `os.Getwd` resolution.

## Scope

### Files To Create Or Modify

- Issue get/list/ready domain result types, loaders, and command adapters/tests.
- Delete section helpers and `cmd/fid/issue_section_test.go`.

### Files To Read For Context

- `../contracts.md` sections 3.1-3.3, 3.5, 3.7, and 5.5; `05_DESIGN_PLAN.md` CLI guidance.

### Out Of Scope

- Mutations, per-Issue actionable-reason projection, and role-specific AC verification.

## Acceptance Criteria

- [ ] Exact decoded key/type/nullability assertions cover get/list and every read returns `event:null`. (VAL-013, VAL-014, VAL-017)
- [ ] Retained filters and their exact value domains/tie-breaks work;
      removed, malformed, and unlisted inputs reject. (VAL-014, VAL-016, VAL-044)
- [ ] Ready proves sequential/prioritized selection plus both exclusive result
      states and the closed explanation-code/nullability registry. (VAL-015)
- [ ] Read collections pass at limits and reject over-limit without truncation. (VAL-045)
- [ ] Human output is stdout-only on success and offers no `.fid/` editing advice. (VAL-008)

## First Failing Test

`cmd/fid/issue_get_list_test.go` - `TestIssueGetJSONMatchesClosedSchema`.

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

- Update shared test decoders rather than copying key assertions per command.
- Completion report: `task-8-done.md`.
