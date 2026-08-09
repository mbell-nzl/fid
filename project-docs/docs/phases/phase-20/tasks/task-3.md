# Phase 20 - Task 3: Atomic project initialization

## Objective

Implement the accepted staged initialization protocol and expose its successful
result through `project.init`. A failed initialization must leave no partial fid
workspace, while a successful one must already have proven journal/projection
rebuild equivalence before `.fid/` becomes visible.

## Validation Contract Coverage

- VAL-001, VAL-002, VAL-004, VAL-008, VAL-017, VAL-039, VAL-042, VAL-043, and
  the `project.init` row of VAL-044.

## Binding Constraints

- Implement `04_TRD_CONFIG.md` section 5 and `../contracts.md` sections 2.1-2.3
  literally: validate, stage, initialize schema, append genesis, rebuild and
  compare, fsync, then atomically rename.
- Use injectable filesystem/stage hooks only at the narrow adapter boundary so
  every named failure point is deterministic in tests.
- Never delete or replace an existing `.fid/`. Cleanup owns only the unique
  staging path created by the current attempt.
- Write the exact strict target `config.toml` schema from Task 2. Temporary
  legacy `pricing.toml` remains separate and expires in Phase 25.
- Explicit/default prefixes follow the closed grammar and invalid derived
  directory names fail before staging; project identity remains generated and opaque.

## Scope

### Files To Create Or Modify

- `internal/bootstrap/init.go` and focused staging/activation helpers/tests.
- `cmd/fid/init.go` and CLI integration tests for the exact Response Envelope.
- Bootstrap journal/rebuild fixtures needed to prove the activated workspace.

### Files To Read For Context

- `04_TRD_CONFIG.md` sections 1, 2, and 5; `../contracts.md` sections 2.1-2.5.

### Out Of Scope

- General backup/restore, existing-workspace migration, capability
  advertisement, and portable export.

## Acceptance Criteria

- [ ] Every injected failure before activation leaves no `.fid/` and no staging residue. (VAL-042)
- [ ] Existing `.fid/` is never replaced; collision fails with the exact error envelope. (VAL-002, VAL-042)
- [ ] Successful init activates once only after schema, genesis Event, and byte-equivalent rebuild verification. (VAL-017, VAL-042)
- [ ] Activated config has exact target keys and round-trips through the strict parser. (VAL-039, VAL-043)
- [ ] Explicit and derived prefix boundary/malformed cases obey the exact value domain. (VAL-043, VAL-044)
- [ ] Human and JSON success/failure use the required channels and exact `project.init` request/result shape. (VAL-001, VAL-004, VAL-008, VAL-044)

## First Failing Test

`internal/bootstrap/init_test.go` - `TestInitializeFailureBeforeActivationLeavesNoWorkspace`.

## Documentation Updates

- None expected; the accepted contract and config specification already own the behavior.

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

- Do not prove cleanup with timing or process interruption; inject each failure
  synchronously and assert exact filesystem state.
- Completion report: `task-3-done.md`.
