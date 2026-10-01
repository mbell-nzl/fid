# Phase 23 - Task 19: Human-mode failures from the Stop-error family - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- `renderStopError` now writes one `error: <JSON message>` line to stderr for Human failures. Its invalid-argument fallback uses the same writer. JSON rendering, typed errors, exit classes, and failure details remain on the existing path.
- Added a built-binary Human/JSON matrix for all seven codes constructed by `renderStopError` across claim, run, gate, Stop, and Attention leaves. The matrix checks exact Human text, empty stdout, matching exit status, repeatable JSON bytes, and no new Event. It compares its code list against the renderer source so a new code needs a row.

## Files created or modified

- `cmd/fid/stop.go`
- `cmd/fid/integration/p23_human_failure_test.go`
- `project-docs/docs/phases/phase-23/tasks/task-19-done.md`

## Documentation updated

- None required; the accepted public interface already requires this Human error line. This completion report records the task evidence.

## Validation Contract Coverage

- VAL-085: PASS for task 19. `TestP23StopFailureCodeMatrix` runs 15 built-binary Human/JSON pairs, including both open-Stop codes for `claim acquire` and `run start`, the pre-frame `gate run` denial, Stop/Attention resolution failures, and invalid list/raise requests. `p23AssertStopFailureCodeInventory` checks the renderer's seven code constants against the table. `p23AssertStopFailurePair` checks exact stderr, empty stdout, equal exit status, repeated byte-identical JSON, and unchanged Event count. The focused tagged run passed (`ok .../cmd/fid/integration 5.638s`, exit 0); the full tagged gate passed at 215.530s.
- VAL-072: The task-19 contribution to HUMAN-03a is complete: `TestP23StopFailureHumanAndJSONParity` confirms the previously silent rework denial now prints `error: review rework budget is exhausted` and exits 4. Full VAL-072 remains a phase-end Human judgment after task 21 adds the qualifying-resolution guidance and repairs the resolution path.

## Tests run

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/... -run '^TestP23StopFailureHumanAndJSONParity$' -count=1` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`). The built `claim acquire` exited 4 with stdout `""` and stderr `""`; the test expected `error: review rework budget is exhausted\n`.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/... -run '^TestP23StopFailureCodeMatrix$' -count=1` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`). The table exposed silent Human failures for the Stop-error family; its first gate fixture also exposed budget admission precedence and was corrected before the green run.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/... -run '^TestP23StopFailureCodeMatrix/gate/open-stop$' -count=1` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`). With an observed budget baseline, the real pre-frame gate Stop denial exited 4 with empty Human stderr.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/... -run '^TestP23StopFailure(HumanAndJSONParity|CodeMatrix)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.638s` (terminated `exit_code: 0`).
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	196.778s` (terminated `exit_code: 0`).
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	55.842s` (terminated `exit_code: 0`). This uncached measurement supplies the workflow wall time because the exact ordinary gate reported it as cached.
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- Built CLI in fresh `/tmp/fid-task19-aD2m7t`: `/tmp/fid --project /tmp/fid-task19-aD2m7t attention list --subject issue:ISS-1` - observed `error: invalid Stop request` (terminated `exit_code: 2`); `--json` twin returned code `invalid-argument`, message `invalid Stop request`, and the same `exit_code: 2`.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`).
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`). `cmd/fid` reported 215.387s; `internal/workflow` was cached in this command and reported no wall time. Its uncached package run above reported 55.842s. No ordinary package exceeded 600s.
- `go tool cover -func=coverage.out` - PASS - `60.5%` (last-line token; terminated `exit_code: 0`).
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	215.530s` (terminated `exit_code: 0`), below the 900s budget.

## Deviations from spec

- None. The gate fixture needed a baseline budget observation so `gate run` reached the intended Stop denial; the initial table red result is preserved above.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `180f5168e3f1338cb758f8388f1c26fd6f8428d3`

## Commit subject

- `fix(stop): render Human failures on stderr`
