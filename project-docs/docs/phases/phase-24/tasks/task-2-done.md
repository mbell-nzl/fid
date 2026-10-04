# Phase 24 - Task 2: Input-as-written digests and named conflict fields - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: session-continuation
- Condition: Task 2 is partially implemented. The shared field-digest helper and one `run.start` path are green, but the other keyed leaves, observation ordering, completeness corpus, documentation, and exact gate remain unfinished.
- Recovery attempts: 0
- Next action: Resume Phase 24 Task 2 on this checkpoint; complete all 43 digest sites, remove the 11 family helpers, cover all 49 keyed leaves and seven observing leaves, then run the exact ordered gate and finalize this report.
- Human attention required: no

## Summary of what was implemented

- Added a shared request and per-field SHA-256 helper with sorted field-difference calculation. The helper hashes the family's supplied digest components and derives public field digests from the canonical request object plus `--actor`.
- Adopted the helper for `run.start` through `mutateRun`. Its same-key `claim_version` conflict now carries the stored and supplied field maps and renders an exit-4 `idempotency-conflict` naming `claim_version`.
- Added the accepted `differing_fields` member to the typed conflict details definition. Existing explicit conflict renderers call one details builder. The remaining family renderers and digest sites still need adoption.

## Files created or modified

- `internal/journal/request_digest.go`; `internal/journal/request_digest_test.go`
- `internal/workflow/run.go`; `internal/cli/codes.go`
- `cmd/fid/idempotency_conflict.go`; `cmd/fid/integration/p24_idempotency_test.go`
- `cmd/fid/claim.go`; `cmd/fid/concern.go`; `cmd/fid/decision.go`; `cmd/fid/deviation.go`; `cmd/fid/envelope.go`; `cmd/fid/evidence.go`; `cmd/fid/finding.go`; `cmd/fid/human_outcome.go`; `cmd/fid/review.go`; `cmd/fid/run.go`
- `project-docs/docs/phases/phase-24/tasks/task-2-done.md`

## Documentation updated

- The Task 2 completion report records the partial implementation. The specified `project-docs/specs/06_DESIGN_PLAN_CLI.md` update remains for continuation.

## Validation Contract Coverage

- VAL-002: BLOCKED - The every-keyed-leaf rebuild parity corpus is unfinished.
- VAL-005: BLOCKED - The required 100k benchmark has not run on a final tree.
- VAL-006: BLOCKED - `TestP24RunStartConflictNamesChangedClaimVersion` proves the `run.start` one-field case and no Event, but the other 48 keyed leaves and the 20 family multi-field and Actor rows are unfinished.
- VAL-007: BLOCKED - Canonical retry rows for every keyed leaf are unfinished.
- VAL-008: BLOCKED - `TestComputeRequestDigestOptionalPresence` covers nonempty versus omitted maps at the shared helper; the real `review.record` nonempty/omitted and `[]`/omitted cases are unfinished.
- VAL-009: BLOCKED - Changed-observation retries for the seven existing observing leaves are unfinished.
- VAL-052: BLOCKED - Required ordinary and tagged integration wall times have not been captured on the final tree.
- VAL-066: BLOCKED - Unobservable and deleted-file retries with new-key controls for the seven observing leaves are unfinished.

## Tests run

- First red: `go test -tags integration -run '^TestP24RunStartConflictNamesChangedClaimVersion$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`); `run.start` returned exit 2 `invalid-argument`, rather than the required exit-4 conflict.
- Helper red: `go test -run '^TestComputeRequestDigest' -count=1 ./internal/journal` - FAIL - `FAIL` (terminated `exit_code: 1`); `ComputeRequestDigest` and `DifferingRequestFields` were undefined.
- `go test -run '^TestComputeRequestDigest' -count=1 ./internal/journal` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	0.004s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP24RunStartConflictNamesChangedClaimVersion$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.213s` (terminated `exit_code: 0`)
- `go test -run 'Conflict|Idempotency' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.742s` (terminated `exit_code: 0`)
- After formatting, `go test -tags integration -run '^TestP24RunStartConflictNamesChangedClaimVersion$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.560s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - BLOCKED - Final-tree gate deferred until the remaining Task 2 implementation and corpus are complete.
- `golangci-lint run ./...` - BLOCKED - Final-tree gate deferred until the remaining Task 2 implementation and corpus are complete.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - Final-tree gate deferred until the remaining Task 2 implementation and corpus are complete.
- `go tool cover -func=coverage.out` - BLOCKED - Coverage gate depends on the unfinished ordinary test gate.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - Final-tree gate deferred until the remaining Task 2 implementation and corpus are complete.
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - BLOCKED - Final-tree benchmark deferred until the remaining Task 2 implementation and corpus are complete.
- `rg -n "RequestDigest:" internal cmd -g '!*_test.go'` - BLOCKED - Final registry check deferred; 42 other sites remain to adopt.
- `rg -n "func [a-zA-Z]*(RequestDigest|requestDigest|judgmentDigest)\(" internal -g '!*_test.go'` - BLOCKED - Final registry check deferred; family-local helpers remain.

## Deviations from spec

- No accepted-contract deviation. This is a same-task checkpoint, not task completion.

## Concerns, technical debt, or notes for the Tech Lead

- The typed conflict details definition now requires `differing_fields`, but most write sites still store `{}`. The continuation must adopt real maps everywhere before running the complete gate. Additional family renderers currently fall through to generic failures and need the conflict mapping required by Task 2.
- No retained assertion has been changed in this checkpoint. Before converging any displaced retained assertion, inspect its introducing commit and current accepted contract, as required by DEC-162/165/166.
- The first red exposed a missing `run.start` CLI conflict mapping. The integration test and focused green prove that one path only; they do not prove the 49-leaf or observation contracts.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
