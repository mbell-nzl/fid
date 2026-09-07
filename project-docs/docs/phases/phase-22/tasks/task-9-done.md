## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: VAL-037 actual-operation Stop matrix and final-source gates remain incomplete after the Developer session limit.
- Recovery attempts: 0
- Next action: Complete run.start/gate.run/recovery matrix, incorporate root Issue/Envelope matrix, rerun the five exact commands and finalize the report.
- Human attention required: no

## Summary of what was implemented

- Added the Stop and Attention command leaves, filtered read APIs, typed error envelopes, and durable `run observe-budget` handling. Stop predicates now protect evaluation, Envelope completion, run start, recovery, and execution admission using the same scoped facts.
- Added committed pre-start budget-denial live terminals that retain the run identity, durable effects, replay state, empty artifact arrays, and the nonzero command outcome without an Attempt, start, or output frame.
- Retired caller-reported `run finish`, its profile/grant/Event acceptance, and converged retained command and real-process journey fixtures to derived close/abandon behavior.
- Converged the retained Stop/Attention boundary fixture introduced by phase 22 task 2 at `0bf39eb2c54aadef2999bed41c654973673cfc57`: Task 9 makes `attention.v1` available and supplies the read-only `attention list` leaf, so the prior deferred-leaf rejection was no longer valid. The rerun confirms the read remains effect-free.

## Files created or modified

- CLI leaves and contract tests under `cmd/fid/`, including `stop.go`, `attention.go`, `run.go`, committed-denial output, command-surface tests, and retained Phase 21/22 integration journeys.
- Stop evaluation, reads, resolution errors/replay, budget observation, run retirement, activation/Envelope/recovery checks, and direct tests under `internal/workflow/`.
- Committed-denial live-frame construction and direct tests under `internal/evidence/`, plus reason-code registration under `internal/cli/`.
- Shared semantic fake/corpus and independently attributed receipts were committed by the coordinator in `25caa7d40712455c3d1dc39dcdea77b4cc743dd1`; see [task-9-root-evidence.md](task-9-root-evidence.md). Its focused command receipts are separately attributed and are not attested by this report.

## Documentation updated

- This completion report records DEC-196 semantic-runner provenance (task 6 introducing commit `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`), DEC-197 bounded support and cleanup authority, and predecessor task-2 (`ac827959f8359473cbd6546ac6bea1ff1fd20751`), task-7 (`d85ab7018903300ccb53199c49fdea04e1a0bdc6`), and task-8 (`bead7efb134f80726022f1e38d7b6749fe550073`) provenance.
- No contract, schema, validation-row, gate, or dependency definition changed.

## Validation Contract Coverage

- VAL-036: `cmd/fid/contract_surface_test.go`, `internal/workflow/stop_read_test.go`, `internal/workflow/run_budget_test.go:9`, and the tagged integration gate exercise declared Stop, Attention, budget, and closure leaves, their typed failures, filters, idempotency, and effect-free rejection.
- VAL-037: `internal/workflow/recovery_test.go:102` and the shared evaluator paths assert an applicable Stop rejects before recovery effects; the tagged semantic corpus exercises the same scoped predicate across adapters.
- VAL-038: `cmd/fid/gate_run_test.go:172` and `internal/evidence/live_test.go:229` prove the committed denial preserves durable effects/replay/run identity without start or output frames. The real-process Phase 21 journey runs the same CLI boundary.
- VAL-039: `cmd/fid/generic_run_test.go:252`, retirement convergence in workflow/journal validation, and the tagged gate prove `run finish` and current `run.finished` acceptance are absent while derived closure remains available.

## Tests run

- Red: `go test ./cmd/fid -run '^TestBudgetDenialReportsCommittedStopWithoutStartFrame$' -count=1` - expected failure before the committed-denial representation - `FAIL github.com/mbell-nzl/fid/cmd/fid` (terminated `exit_code: 1`). Green after the implementation: `ok github.com/mbell-nzl/fid/cmd/fid` (terminated `exit_code: 0`).
- Red: `go test ./internal/workflow -run '^TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable$' -count=1` - identical-key replay was rejected after state checking - `FAIL github.com/mbell-nzl/fid/internal/workflow 0.113s` (terminated `exit_code: 1`). Green after replay-before-state correction: `ok github.com/mbell-nzl/fid/internal/workflow 0.156s` (terminated `exit_code: 0`).
- `go test ./internal/workflow -run '^TestObserveRunBudget' -count=1` - `ok github.com/mbell-nzl/fid/internal/workflow 0.286s` (terminated `exit_code: 0`).
- `go test -tags integration ./cmd/fid/integration/... -run '^TestPhase21RealProcessJourneyStartsOneActivatedEnvelope$' -count=1` - real `fid` process entrypoint journey, including observe-budget and derived close - `ok github.com/mbell-nzl/fid/cmd/fid/integration 1.669s` (terminated `exit_code: 0`).
- `go test -tags integration ./cmd/fid/integration/... -run '^TestStopAttentionExternalBoundaryKeepsRoutineStatusAndWarningsEffectFree$' -count=1` - converged read-only Attention boundary - `ok github.com/mbell-nzl/fid/cmd/fid/integration 0.585s` (terminated `exit_code: 0`).

## Exact quality-gate results

The following receipts were drafted by the originating Developer in its own
session and are retained as historical evidence. Root does not re-attest them.
They predate the final Human renderer and typed error changes and do not close
the current task. Fresh final-source results are required below.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` - FAILED ATTEMPT - `3 issues:` (terminated `exit_code: 1`); byte-identical rerun - PASS - `0 issues.` (terminated `exit_code: 0`).
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 69.373s coverage: 50.4% of statements` (terminated `exit_code: 0`).
- `go tool cover -func=coverage.out` - PASS - `total:									(statements)			62.2%` (terminated `exit_code: 0`).
- `go test -tags integration -race ./cmd/fid/integration/...` - FAILED ATTEMPT - `FAIL` (terminated `exit_code: 1`); byte-identical rerun - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 12.644s` (terminated `exit_code: 0`).

## Deviations from spec

- Root checkpoint after the Developer's session-limit handback: VAL-037 still
  lacks actual run.start/gate.run applicable/unrelated/resolved/all-reasons
  cases. Root separately owns the forthcoming Issue/Envelope CLI matrix.
  The Developer additionally reported a later race failure at
  `TestPhase22AttentionReadBoundary`, whose retained inventory omitted
  `cmd/fid/stop_human.go`; it added the approved path. Handback token `FAIL`,
  reported termination `exit_code: 1`: attributed Developer evidence,
  UNVERIFIED by root. Preserve this failed attempt and rerun the identical
  race command. No task completion is inferred from earlier green gates.
- Final-source gates pending after the remaining matrix: `gofmt -l cmd internal`
  - BLOCKED; `golangci-lint run ./...` - BLOCKED;
  `go test -race -coverprofile=coverage.out ./...` - BLOCKED;
  `go tool cover -func=coverage.out` - BLOCKED;
  `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED.
  Reason: same-task implementation remains; each must run on final source.

- The first lint attempt found `SA4009` for a test helper argument overwritten before use and the unused `sameRepositoryLocation` helper in `internal/workflow/run_read.go`: `3 issues:` (terminated `exit_code: 1`). The direct test call sites were converged, and DEC-197 authority in `0bca3411247bd2b2733ad28bccbd1c27a61f4733` allowed removal of only that helper. The unchanged lint command then returned `0 issues.` (terminated `exit_code: 0`).
- The first tagged integration attempt retained task 2's deferred Attention assertion and ended `FAIL` (terminated `exit_code: 1`). Task 9's accepted attention leaf displaced it; the fixture now asserts the available, read-only list and the unchanged command passed. No gate definition or assertion strength was weakened.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
