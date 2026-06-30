## Summary of what was implemented

- Extracted a reusable completion-effect builder from `internal/issues/complete.go` and moved completion-specific support helpers into `internal/issues/complete_support.go` so `fid issue complete` still uses the same completion write path while `internal/dev` can reuse it.
- Added `internal/dev/Finish` and `internal/dev/Block` as one-event composites. `dev finish` closes an open run when present, writes the shared completion, transitions to `dev_complete`, and releases the caller claim. `dev block` writes the shared failed-gate completion, transitions to `failed_gate`, and releases the caller claim without touching `agent_runs`.
- Added `fid dev finish` / `fid dev block` CLI subcommands, success confirmations, deny-path coverage, and zero-append reject coverage through domain and CLI tests.

## Files created or modified

- Modified: `internal/issues/complete.go`
- Added: `internal/issues/complete_support.go`
- Existing `internal/issues/complete_test.go` assertions remained in place and passed unchanged.
- Modified: `internal/dev/dev.go`
- Added: `internal/dev/finish.go`
- Added: `internal/dev/block.go`
- Added: `internal/dev/finish_test.go`
- Added: `internal/dev/block_test.go`
- Modified: `cmd/fid/dev.go`
- Modified: `cmd/fid/dev_test.go`

## Documentation updated

- Added this completion report: `project-docs/docs/phases/phase-15/tasks/task-3-done.md`.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` was not changed because the implemented `fid dev finish` / `fid dev block` confirmations matched the existing examples closely enough to avoid spec drift.

## Validation Contract Coverage

- `VAL-006`: completed the `dev.finish` / `dev.block` authorization coverage. CLI tests prove `policy show` rows remain visible and `qa_reviewer` is rejected before append with zero new journal lines.
- `VAL-007`: `internal/dev/finish_test.go::TestFinishClosesRunCompletesTransitionsAndReleasesClaimInOneEvent` proves one `dev.finished` event closes the run, writes the completion snapshot, transitions the issue, and releases the claim.
- `VAL-008`: `internal/dev/finish_test.go::TestFinishRejectsWhenDeveloperACsUnverifiedBeforeAppend` proves reject-before-append on unmet `all_acs_verified`; `TestFinishSucceedsWithoutPreexistingCompletionOnceDeveloperACsPass` proves `has_completion_report` is satisfied in-event.
- `VAL-009`: `internal/dev/block_test.go::TestBlockRecordsFailedGateCompletionTransitionsAndReleasesClaim` proves one `dev.blocked` event writes a `failed_gate` completion with reason, transitions the issue, and releases the claim without closing the run.
- `VAL-010`: `issues.BuildCompletionEffects` is the shared builder reused by `internal/dev`; `go test ./internal/issues -run 'TestComplete'` stayed green, preserving `fid issue complete` behavior. Function/file spot checks: `internal/issues/complete.go` 351 lines, `internal/issues/complete_support.go` 162 lines.
- Coverage from the required race+coverage gate: `internal/dev` 75.4%, `internal/issues` 81.7%, `cmd/fid` 76.9%.

## Tests run

- Red: `go test ./internal/dev ./cmd/fid`
  Expected failure before implementation: undefined `Finish` / `Block` symbols.
- Green narrow loop: `go test ./internal/dev ./cmd/fid`
- Green extraction regression: `go test ./internal/issues -run 'TestComplete'`
- Required gate: `go test -race -coverprofile=coverage.out ./...`
- Required integration gate: `go test -tags integration -race ./...`

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
  No output.
- `PASS` `golangci-lint run ./...`
  Output excerpt: `0 issues.`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
  Output excerpt: `cmd/fid 76.9%`, `internal/dev 75.4%`, `internal/issues 81.7%`.
- `PASS` `go tool cover -func=coverage.out`
  Output excerpt: `total: (statements) 75.8%`.
- `PASS` `go test -tags integration -race ./...`
  Output excerpt: `ok github.com/mbell-nzl/fid/cmd/fid 61.307s`, `ok github.com/mbell-nzl/fid/internal/dev 4.775s`, `ok github.com/mbell-nzl/fid/internal/issues 18.476s`.
- `PASS` `go build -o /tmp/fid ./cmd/fid`
  No output.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- `internal/dev/Finish` and `internal/dev/Block` are still above the repo’s 60-line non-test function target. The behavior is covered and the required gates are green, but a follow-up extraction could reduce event-assembly and validation noise in those two composites.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `b9e650d`

## Commit subject

- `feat(dev): add finish and block composites`
