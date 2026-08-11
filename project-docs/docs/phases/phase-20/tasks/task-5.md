# Phase 20 - Task 5: Remove orphaned Issue fields

## Recovery Directive (DEC-137)

The first session completed this task's implementation and passed format, lint,
race/coverage, and every focused test, then stopped because
`TestIntegrationClassification` compares the dedicated package against a
hardcoded journey list that omits the newly added
`TestV01ParityFixtureRebuildPreservesEveryProjectionTable`. Scope, acceptance
criteria, VAL coverage, and the command list below are unchanged. This session:

1. Inherits the existing uncommitted working tree, including the regenerated
   `cmd/fid/testdata/v01_parity/*.json`, `SHA256SUMS`, and
   `testdata/format/issue-empty.golden.toon`. Do not discard or re-derive it,
   and re-verify the fixture hashes recorded in `task-5-done.md` rather than
   trusting them.
2. Changes `wantJourneys` in `cmd/fid/integration/classification_test.go` from
   exact equality to **required presence**: the four journeys Task 3f migrated
   must still exist and their loss must still fail, while additional journeys
   are allowed. Make the failure message name the missing journey. Leave
   `taggedTestFilesOutside` and its repository-wide walk, along with
   `TestIntegrationBinaryOwnedByPackage` and
   `TestClosureJourneyHasNoWallClockSleep`, exactly as they are — DEC-134's
   classification invariant is not being relaxed.
3. Proves the relaxed check still fails on a real regression before relying on
   it: red/green the presence assertion by removing a required journey name from
   the expected set, observing the failure, and restoring it.
4. Runs the full command list below and completes the task normally.

Task 6 owns revision and concurrency; nothing here changes that boundary.

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
