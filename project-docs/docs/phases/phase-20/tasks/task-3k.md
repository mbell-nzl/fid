# Phase 20 - Task 3k: Deterministic claims time

## Objective

Replace claim-expiry wall-clock sleeps with explicit deterministic time in the
claims domain boundary. Preserve the zero-value production path as real UTC
time and keep claim ownership and expiry semantics unchanged.

## Validation Contract Coverage

- No new VAL-NNN applies. Existing claim lifecycle and CLI assertions retain
  their coverage; deterministic boundary tests add evidence without changing
  the contract.

## Binding Constraints

- Add an optional time value or clock only to the narrow claims option/domain
  seam that evaluates claim timestamps.
- A zero value obtains `time.Now().UTC()` at the same operation boundary as
  today.
- Tests use fixed instants and explicit boundary advances; no shortened sleep,
  polling, or environment-based clock.
- CLI claim tests may age persisted fixture state through existing storage
  helpers when command dependency injection would broaden production design.
- Do not change TTL defaults, comparison operators, error text, event fields,
  or journal/projection behavior.

## Scope

### Files To Create Or Modify

- `internal/claims/` implementation and tests
- `cmd/fid/claim_test.go`
- `cmd/fid/issue_claim_test.go`
- Narrow test helpers needed to age persisted claim state
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/tasks/task-3k-done.md`

### Files To Read For Context

- Claims acquire, renew, release, list, and live-claim validation paths
- Existing claim event/projection test helpers
- Task 3j report and DEC-133

### Out Of Scope

- A repository-wide clock abstraction.
- Changes to maintenance time handling; Task 3l owns it.

## Acceptance Criteria

- [ ] Fixed-time domain tests cover just-before, exact-boundary, and
      just-after claim expiry deterministically. (no VAL)
- [ ] Acquire, renew, release, list, and live-claim validation use one supplied
      operation time consistently where applicable.
- [ ] Zero-value option tests prove production defaults to a current UTC time.
- [ ] Claim command tests contain no `time.Sleep` and retain their original
      ownership/expiry assertions.
- [ ] Repeated claims and command-package race tests pass.

## First Failing Test

Add a fixed-time exact-boundary takeover test in `internal/claims`; confirm it
fails because the current operation reads the wall clock directly.

## Documentation Updates

- `infra-improvement-pack.md` - clock seam, eliminated sleep inventory, and
  repeated-race evidence.

## Exact Commands To Run

Run exactly:

```text
go test -race -count=10 -timeout 30m ./internal/claims ./cmd/fid
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go build -o /tmp/fid ./cmd/fid
```

Required-command rules live in `AGENTS.md` (Commands).

**2026-08-18, DEC-147 (Tech Lead).** The first command gained `-timeout 30m`.
`-count` repeats inside one test-binary invocation, so Go's default
`-timeout 10m` was an unstated wall-clock budget rather than a chosen
assertion. The repeat count, race detector, package scope, and pass/fail
semantics are unchanged. Expect about 18 minutes for `cmd/fid`; poll it to a
real termination.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Prefer an explicit operation timestamp in existing option structs over a
  package-global replaceable function.
- Completion report: `tasks/task-3k-done.md`.
