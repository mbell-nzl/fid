# Phase 20 - Task 4: Rebuild, authority, and lock failures

## Objective

Migrate `fid rebuild` to the project-selection and Response Envelope contract,
map retained fail-closed V1 authority errors, and make write-lock contention a
typed retryable failure.

## Validation Contract Coverage

- VAL-001, VAL-002, VAL-006, VAL-007, VAL-008, VAL-017, and the rebuild row of VAL-044.

## Binding Constraints

- V1 permission checks remain fail-closed and execute before journal/projection work.
- Map existing missing, unknown, and denied cases to the semantic authority codes;
  do not weaken checks or promise those general codes will be deleted in Phase 21.
- Rebuild is projection-only and returns `event: null`.
- Lock tests acquire and release the lock deterministically; timing races are invalid evidence.

## Scope

### Files To Create Or Modify

- `internal/journal/lock.go` and tests.
- `cmd/fid/rebuild.go`, shared error mapping, and rebuild/root tests.

### Files To Read For Context

- `internal/policy/authorize.go`; `../contracts.md` sections 2.2-2.5, 3.2, and 3.5.

### Out Of Scope

- Generic Actor/Envelope authority, Claims/Runs, idempotency, and Issue schema.

## Acceptance Criteria

- [ ] Rebuild accepts only its closed request and returns its exact result with `event:null`. (VAL-017, VAL-044)
- [ ] Missing, unknown, and denied V1 authority map to exact codes/details and append nothing. (VAL-002, VAL-006)
- [ ] Held write lock returns `lock-unavailable`, `retryable:true`; release allows retry. (VAL-007)
- [ ] Human success is stdout-only; Human and JSON failures are stderr-only. (VAL-001, VAL-008)

## First Failing Test

`cmd/fid/rebuild_test.go` - `TestRebuildJSONMapsMissingAuthorityWithoutMutation`.

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

- Completion report: `task-4-done.md`.
