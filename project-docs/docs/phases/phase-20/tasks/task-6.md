# Phase 20 - Task 6: Replayable Issue aggregate revision

## Objective

Add Issue aggregate revision through migration `0012_p20_issue_revision`,
implement the reusable revision-Effect helper, and apply it to create, update,
transition, and the legacy `dev` start/finish/block status mutations. Tasks
10-11 apply the same helper to AC, Issue-label, and link Events so every family
listed in the contract advances the aggregate.

## Validation Contract Coverage

- VAL-010, VAL-011, and the revision half of VAL-012.

## Binding Constraints

- The caller never supplies the resulting revision.
- For update-like effects, `Values.revision = prior + 1`, while `Prior.revision`
  and `Preconditions.revision` both equal the observed prior revision.
- Issue creation writes revision 1. The helper can add one ordered revision
  Effect per affected Issue; rejected operations append nothing and change none.
- Until Phase 21 deletes `fid dev`, `dev.started`, `dev.finished`, and
  `dev.blocked` use the same helper whenever they mutate Issue status. They do
  not gain expected-revision flags or canonical Response Envelopes.
- Use the existing absolute `Effect.Values` projection semantics; do not add an
  increment operation or a new Event family.

## Scope

### Files To Create Or Modify

- `migrations/0012_p20_issue_revision.up.sql` and `.down.sql`.
- `internal/issues` create/update/transition write paths, read models, and the
  narrowly shared Issue-revision Effect helper used later by AC/label/link.
- `internal/dev` start/finish/block status-Effect call sites and focused regressions.
- Golden Event fixtures and replay/rebuild tests for the six owned families.

### Files To Read For Context

- `internal/journal/event.go`, `project_effects.go`; `../contracts.md` sections 3.6, 4, and 7.

### Out Of Scope

- `--expected-revision` request handling and JSON response reshaping (Task 9).

## Acceptance Criteria

- [ ] New Issues start at revision 1; owned update/transition and legacy dev
      status mutations advance exactly one. (VAL-010)
- [ ] Six owned golden Event families carry exact absolute/prior/precondition values. (VAL-011)
- [ ] Shared helper produces sorted one-or-more Issue revision Effects without caller-supplied results. (VAL-010, VAL-011)
- [ ] Replay matches live revision and stale preconditions fail closed. (VAL-011, VAL-012)
- [ ] Down migration removes the column cleanly. (VAL-012)

## First Failing Test

`internal/issues/update_test.go` - `TestUpdateEventRecordsReplayableRevision`.

## Documentation Updates

- None expected.

## Exact Commands To Run

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

## Notes For Developer

- A need for new Effect semantics is a hard stop, not a local adaptation.
- Completion report: `task-6-done.md`.
