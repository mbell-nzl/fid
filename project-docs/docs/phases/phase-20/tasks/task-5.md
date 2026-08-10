# Phase 20 - Task 5: Remove orphaned Issue fields

## Objective

Remove `in_inbox`, `est_tokens`, and `due_date` from the projection, Issue
domain/CLI, emitted Event effects, fixtures, and rebuild path using migration
`0011_p20_remove_issue_fields`.

## Validation Contract Coverage

- VAL-009 and the field-removal half of VAL-012.

## Binding Constraints

- Follow `../contracts.md` sections 4, 5.5, and 7.
- Rebuild the SQLite table where necessary to remove the dependent CHECK and
  recreate `issues_enriched` without the columns.
- Existing development journals are not migrated. Repository fixtures are
  regenerated and documented; unknown removed effect keys fail replay loudly.
- Do not add or increment `revision`; Task 6 owns that separate migration.

## Scope

### Files To Create Or Modify

- `migrations/0011_p20_remove_issue_fields.up.sql` and `.down.sql`.
- Issue model/create/update/list/ready code and corresponding CLI flags.
- Journal projection column allowlists and affected committed fixtures/tests.

### Files To Read For Context

- Migration `0008_v02a_issues_workflow_fk`; `03_TRD_SCHEMA.md` sections 3 and 8.

### Out Of Scope

- Revision/concurrency, Issue JSON reshaping unrelated to these fields, and Phase 25 view removal.

## Acceptance Criteria

- [ ] Schema, CLI tree, domain structs, JSON, and new Events contain none of the three fields. (VAL-009)
- [ ] Down migration restores the previous logical shape and CHECK behavior. (VAL-009)
- [ ] Regenerated fixtures rebuild to a projection equal to the live projection. (VAL-012)
- [ ] A removed-key development Event fails with `journal-replay-stopped` rather than being ignored. (VAL-009)

## First Failing Test

`internal/db/migrate_test.go` - `TestIssuesTableExcludesRetiredPhase20Fields`.

## Documentation Updates

- None expected.

## Exact Commands To Run

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

## Notes For Developer

- Record the exact fixture regeneration procedure and resulting committed files.
- Completion report: `task-5-done.md`.
