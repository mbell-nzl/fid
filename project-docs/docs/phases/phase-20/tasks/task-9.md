# Phase 20 - Task 9: Issue mutation and concurrency contract

## Objective

Migrate Issue create, update, and transition to their Response Envelopes and add
optional optimistic concurrency to update and transition.

## Validation Contract Coverage

- VAL-002, VAL-008, VAL-010, VAL-011, VAL-017, VAL-018, VAL-019, VAL-040,
  VAL-041, VAL-045, and the three Issue mutation rows of VAL-044.

## Binding Constraints

- `--expected-revision` has absent, matching, and stale states. Stale rejects
  before append with `state-conflict` and both values in `error.details`.
- Request, result, and Event shapes follow `../contracts.md` sections 3.2,
  3.4, 3.6, and 7 exactly. `--epic` and `--parent` retain their closed,
  kind-specific meanings and are mutually exclusive.
- Illegal/unknown transitions and unsupported/unsatisfied preconditions use
  `transition-not-permitted` with exact details; no message parsing.
- V1 `--by` remains required and fail-closed until Phase 21.
- Priority, sequence, expected revision, state, path, token, and text inputs
  obey the exact section 3.2.1 value domains and reject before append.
- Read the resolved project root from Task 2 command context.
- Remove public `--owner-role` and `--owner-actor` from create/update and their
  request serialization. Retain internal defaults/storage/loaders needed by
  replacement-bound V1 paths; do not expose hidden owner mutations in result changes.

## Scope

### Files To Create Or Modify

- Issue create/update/transition command adapters, requests/results, typed errors, and tests.
- Shared CLI error mapping for Issue mutation failures.

### Files To Read For Context

- `../contracts.md` sections 2.4, 3.2, 3.4, and 3.6; Task 6 revision implementation.

### Out Of Scope

- Idempotency keys, Envelope authority, labels/AC/links/workflows.

## Acceptance Criteria

- [ ] All three closed requests/results and Event identities match the appended Event. (VAL-010, VAL-017, VAL-044)
- [ ] Owned Event fixtures carry the exact aggregate revision tuple. (VAL-011)
- [ ] Absent and matching revision proceed; stale rejects without state/journal change. (VAL-018)
- [ ] Legal transition succeeds; unknown target, absent edge, and unsupported/unsatisfied precondition reject before append. (VAL-019)
- [ ] Supplied local attribution is recorded without trusted/protected side effects. (VAL-041)
- [ ] Human success is stdout-only; JSON errors use exact stderr envelopes. (VAL-002, VAL-008)
- [ ] Owner-role/actor flags are unknown and no canonical mutation result can
      report an owner field change. (VAL-040, VAL-044)
- [ ] Mutation value-domain boundaries and malformed cases reject atomically;
      `issue.update.changes` remains field-ordered. (VAL-044, VAL-045)

## First Failing Test

`cmd/fid/issue_update_test.go` - `TestIssueUpdateRejectsStaleExpectedRevision`.

## Documentation Updates

- None expected.

## Exact Commands To Run

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go build -o /tmp/fid ./cmd/fid
```

## Notes For Developer

- Read and validate expected revision while holding the mutation lock.
- Completion report: `task-9-done.md`.
