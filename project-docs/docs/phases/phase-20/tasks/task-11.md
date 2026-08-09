# Phase 20 - Task 11: Link contract and dependency invariants

## Objective

Migrate link add/remove/list and enforce self-link, missing-endpoint, and
dependency-cycle rejection over the existing links projection.

## Validation Contract Coverage

- VAL-008, VAL-010, VAL-011, VAL-017, VAL-023, VAL-024, VAL-025, VAL-041,
  VAL-044, VAL-045, and VAL-046.

## Binding Constraints

- Exact typed endpoint requests, Link schemas, codes, revision results, and
  ordering are in `../contracts.md` sections 2.4-2.5, 3, and 7.
- Require explicit `--from-type`, `--to-type`, or list `--type`. Resolve the
  supplied type through stored identity; never classify by ID prefix.
- Cycle detection applies only to normalized `blocks` direction; the stored
  inverse edge must not create false positives.
- Rejection occurs before append. Missing Issue and requirement endpoints have
  distinct codes; cycle details name endpoints and ordered path.
- One bounded traversal must terminate even if stored data is already cyclic.
- Read the resolved project root from Task 2 command context.

## Scope

### Files To Create Or Modify

- `internal/links` add/remove/list, typed errors, new cycle helper and tests.
- `cmd/fid/link.go` and focused CLI tests.

### Files To Read For Context

- Link registry/errors; PRD-ISSUE-003/004; `../contracts.md` sections 2.4 and 3.

### Out Of Scope

- Normalized `issue_dependencies`, parent/dependency combined cycles, and legacy-data repair.

## Acceptance Criteria

- [ ] Add/remove/list accept only typed endpoint requests and emit exact schemas/Event identities. (VAL-017, VAL-023, VAL-044, VAL-046)
- [ ] Link mutations advance every Issue endpoint once, return sorted revisions, and match golden Events. (VAL-010, VAL-011)
- [ ] Unsupported relation, duplicate, typed missing endpoints, missing link, and self-link use exact codes without append. (VAL-023, VAL-025, VAL-046)
- [ ] Self, two-node, three-node cycles reject without append; an acyclic diamond succeeds. (VAL-024)
- [ ] Link lists and revision arrays obey exact order/bounds. (VAL-045)
- [ ] Supplied attribution is recorded without trusted/protected effects; Human success is stdout-only. (VAL-008, VAL-041)

## First Failing Test

`internal/links/cycle_test.go` - `TestAddRejectsThreeNodeDependencyCycle`.

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

- An inability to normalize inverse edges is a data-model hard stop.
- Completion report: `task-11-done.md`.
