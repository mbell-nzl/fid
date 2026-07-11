# Phase 17 - Task 2: Add the Maintenance Grace Configuration Contract

## Objective

Add the isolated project-config contract needed by the maintenance audit:
started-category issues without a live claim receive a configurable grace
window, defaulting to five minutes. Preserve compatibility with every existing
v0.1 `config.toml` that contains only `[project]`.

## Validation Contract Coverage

- VAL-002: missing maintenance config defaults to 300 seconds; explicit valid
  values round-trip; invalid values fail clearly.

## Binding Constraints

- DEC-050 is binding; do not reinterpret its key, default, or boundaries.
- Canonical key: `[maintenance].started_without_claim_grace_seconds`.
- Missing section/key means 300 seconds; zero is valid; negative and malformed
  values are errors. No migration rewrites an existing config file.
- Use the existing in-tree config parser. No new TOML/runtime dependency.

## Scope

### Files To Create Or Modify

- `internal/config/config.go` - model, default, serialization, and parsing.
- `internal/config/config_test.go` - compatibility/boundary/round-trip tests.
- `cmd/fid/init_test.go` - assert new workspaces receive the default section.

### Files To Read For Context

- `cmd/fid/init.go`
- `scaffold/initial-spec.md` section 9.8.1 (historical detail only; canonical
  docs above must be updated in this task)
- `project-docs/docs/phases/phase-17/qa/validation-contract.md` DEC-050/VAL-002
- `project-docs/specs/01_PRD.md` PRD-OPS-008 (synced at kickoff by DEC-050)
- `project-docs/specs/04_TRD_CONFIG.md` section 2 (synced at kickoff by DEC-050)

### Out Of Scope

- Maintenance queries or CLI commands.
- Claim TTL configuration.
- A general TOML parser replacement or config migration command.

## Acceptance Criteria

- [ ] A legacy `[project]`-only config reads with a 300-second grace. (VAL-002)
- [ ] Bootstrap output writes the 300-second default; explicit zero and positive
  values serialize and read back unchanged. (VAL-002)
- [ ] Negative and malformed values return contextual errors. (VAL-002)
- [ ] PRD and config specs define the same key/default/compatibility contract.
  (VAL-002)

## First Failing Test

`TestReadProjectConfigDefaultsMissingMaintenanceGraceToFiveMinutes` fails
because `ProjectConfig` has no maintenance grace contract.

## Documentation Updates

- None expected. DEC-050 synchronized `01_PRD.md` and `04_TRD_CONFIG.md` at
  kickoff; stop if implementation proves that contract wrong.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Required-command rules live in `AGENTS.md`. Every command is reported `PASS`,
or `BLOCKED` with reason; the full gate is required before task close.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Stop if approval changes the key, default, or invalid-value semantics.
- Do not silently coerce invalid text or negative values.
- Completion report: `project-docs/docs/phases/phase-17/tasks/task-2-done.md`.
