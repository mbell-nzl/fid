# Phase 21 - Task 14: Live Framing, Redaction, And Controls - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: Phase Runner runs the two-argument report-lint control and dispatches the Human VAL-036 check.
- Human attention required: no

## Summary of what was implemented

- Added the `fid.live.v1` evidence streamer with one global sequence across stdout, stderr, control, and caller-bound terminal frames.
- Added streaming redaction for credential and registered literal sources before frame emission, bounded UTF-8 chunks, safe control-byte rendering, one truncation control while draining both pipes, and ordered redaction/Artifact-rejection controls.
- Added Human rendering for labeled process streams, typed controls, and the terminal summary.
- Added and byte-compared the required VAL-036 machine and Human fixtures. They include interleaved streams, all three controls, a failed terminal outcome, and safely rendered ESC, CR, and BEL data.

## Files created or modified

- `internal/evidence/live.go`
- `internal/evidence/live_test.go`
- `cmd/fid/live_output.go`
- `cmd/fid/live_output_test.go`
- `cmd/fid/testdata/val-036/human-session.txt`
- `cmd/fid/testdata/val-036/live-frames.jsonl`
- `project-docs/docs/phases/phase-21/tasks/task-14-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-14-done.md`

## Validation Contract Coverage

- VAL-036: HUMAN PENDING - `TestLiveOutputGoldenFixtures` at `cmd/fid/live_output_test.go:14` regenerates and byte-compares both prescribed fixtures, including stream identity, every typed control, failed terminal data, and safe ESC/CR/BEL rendering. The required HUMAN-01a..01d verdict remains exclusively for the Human.
- VAL-047: PASS - `TestRunLiveFramesInterleavedOutputControlsAndFailure` at `internal/evidence/live_test.go:11` asserts `fid.live.v1` headers and consecutive global sequence values for a real stdout/stderr child process; `TestLiveOutputGoldenFixtures` byte-compares its interleaved NDJSON fixture.
- VAL-048: PASS - `TestRunLiveFramesInterleavedOutputControlsAndFailure` at `internal/evidence/live_test.go:76` asserts the exact closed control sequence `output-truncated,redaction-applied,artifact-rejected`; the golden test verifies each variant in both rendered forms.
- VAL-049: PASS - `TestRunLiveFramesInterleavedOutputControlsAndFailure` at `internal/evidence/live_test.go:62` rejects either source canary in every serialized frame, and `TestLiveOutputGoldenFixtures` at `cmd/fid/live_output_test.go:31` independently scans both machine and Human fixture output for the credential and literal canaries.
- VAL-050: PASS - `TestRunLiveDrainsBothPipesAfterTruncationAndPreservesOutcome` at `internal/evidence/live_test.go:97` runs a 256 KiB two-pipe child, asserts one truncation control, bounded emitted data, and the independently observed failed exit `24`.

## Tests run

- `go test ./internal/evidence` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`): `RunLive` was undefined before the first framing implementation.
- `go test ./internal/evidence` - PASS - `ok  	github.com/mbell-nzl/fid/internal/evidence	0.022s` (terminated `exit_code: 0`)
- `go test ./internal/evidence ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	16.981s` (terminated `exit_code: 0`)
- `go test -race ./internal/evidence ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	208.601s` (terminated `exit_code: 0`)
- Real process exercise: `TestRunLiveFramesInterleavedOutputControlsAndFailure` at `internal/evidence/live_test.go:11` invokes the test executable through `RunLive`, observes interleaved process pipes and exit `23`; `TestRunLiveDrainsBothPipesAfterTruncationAndPreservesOutcome` at `internal/evidence/live_test.go:97` drains 256 KiB across both pipes and observes exit `24`.

## Exact quality-gate results

- `go test -race ./internal/evidence ./internal/cli ./cmd/fid/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	208.464s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 32.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:						(statements)				66.3%` (terminated `exit_code: 0`)

## Deviations from spec

- The first `golangci-lint run ./...` attempt failed with `* staticcheck: 1` (terminated `exit_code: 1`): `visibleLiveText` used `WriteString(fmt.Sprintf(...))`. Replaced it with `fmt.Fprintf` to the existing builder; the byte-identical command later returned `0 issues.` (terminated `exit_code: 0`).
- The first `go test -race -coverprofile=coverage.out ./...` attempt failed with `FAIL` (terminated `exit_code: 1`) because frame redaction depended on pipe-read boundaries and the terminal byte count continued changing after truncation. A second byte-identical attempt also failed with `FAIL` (terminated `exit_code: 1`), isolating the remaining post-truncation counter change. The redactor now retains potential source prefixes and the bounded count freezes at the first truncation observation; the later byte-identical command passed with `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 32.9% of statements` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `f0eec233bba8695c2b33b5a74237fc0eea591ff7`

## Commit subject

- `feat(evidence): add live output framing`
