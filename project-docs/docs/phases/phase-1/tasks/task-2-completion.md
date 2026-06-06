# Phase 1 - Task 2: Journal Event v1 Model, ID Generation, Canonical JSON, Upgrader Chain - Completion Report

## Summary of what was implemented

- Added the pure-data `internal/journal` package for the v1 journal envelope:
  `Event`, `Actor`, `Command`, and ordered `Effect` types with schema
  constants.
- Added canonical JSON serialization with sorted object keys, preserved array
  order, no insignificant whitespace, and deterministic event hashing with a
  `sha256:` prefix.
- Added generated ID helpers for `<PREFIX>-<hash>` IDs, event ID formatting as
  `EVT-<UTC YYYYMMDD>-<sequence>`, and collision retry via a caller-supplied
  existence predicate.
- Added a schema-upgrader chain with dispatch by `event_schema` and a verified
  v1 identity upgrader at the current schema.
- Added a reusable golden v1 event fixture for later replay tests.

## Files created or modified

- `internal/journal/event.go`
- `internal/journal/canonical.go`
- `internal/journal/id.go`
- `internal/journal/upgrade.go`
- `internal/journal/canonical_test.go`
- `internal/journal/id_test.go`
- `internal/journal/upgrade_test.go`
- `testdata/journal/event-v1.golden.json`
- `project-docs/docs/phases/phase-1/tasks/task-2-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-1/tasks/task-2-completion.md`

## Validation Contract Coverage

- `VAL-007` - PASS. `TestMarshalCanonicalJSONProducesStableV1EventBytes`
  asserts the v1 envelope serializes with `event_id`, `event_schema`,
  `event_type`, `occurred_at`, `actor`, `command`, and an ordered `effects`
  array in canonical JSON. `TestV1GoldenFixtureRoundTrips` proves the committed
  v1 fixture round-trips through unmarshal and canonical marshal unchanged.
- `VAL-009` - PASS for Task 2 scope. `TestUpgradeToCurrentLeavesV1EventUnchanged`
  verifies the v1 no-op path at the current schema, and
  `TestUpgradeChainDispatchesByEventSchema` verifies dispatch by
  `event_schema`.
- `VAL-015` - PASS for Task 2 scope. `TestGenerateIDUsesDefaultSixHexChars`,
  `TestGenerateIDWithLengthSupportsEightHexChars`,
  `TestGenerateIDRetriesOnCollision`, and
  `TestFormatEventIDUsesUTCDateAndZeroPaddedSequence` cover ID shape, 6/8 hex
  support, collision retry, and event ID formatting. Per-day file-backed
  sequence allocation remains Task 3 scope.

## Tests run

- `go test ./internal/journal -run 'TestMarshalCanonicalJSONProducesStableV1EventBytes|TestHashEventReturnsStableDigest'` - PASS
- `go test ./internal/journal -run 'TestGenerateID|TestFormatEventID|TestUpgrade'` - PASS
- `go test ./internal/journal` - PASS
- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `735b2dc`

## Commit subject

- `feat(journal): add canonical event model`
