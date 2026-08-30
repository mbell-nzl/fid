# Phase 21 - Task 18: Guarded Terminal Transition, Completion, And Resume - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Wired the shared workflow evaluator to persisted Receipt, Evidence validity, repository freshness, and required Artifact availability before terminal-success transitions.
- Added Envelope completion guarded by the same evaluator plus terminal-success Issue state, with one completion Event only on success.
- Added bounded read-only `fid resume <issue-id>` projection with closed next actions.

## Files created or modified

- `cmd/fid/contract_surface_test.go`
- `cmd/fid/envelope.go`
- `cmd/fid/evidence_test.go`
- `cmd/fid/resume.go`
- `cmd/fid/root.go`
- `cmd/fid/task18_completion_resume_test.go`
- `internal/cli/codes.go`
- `internal/workflow/envelope.go`
- `internal/workflow/evaluation.go`
- `internal/workflow/resume.go`
- `project-docs/docs/phases/phase-21/tasks/task-18-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-18-done.md` - draft records red/green evidence and awaits the required task gate.

## Validation Contract Coverage

- VAL-064: PASS - `TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease` in `cmd/fid/task18_completion_resume_test.go` drives a successful Receipt link through the terminal `qa_pass` transition before completion.
- VAL-065: PASS - `TestTerminalSuccessTransitionRejectsUnsuccessfulReceiptEvidence` in `cmd/fid/task18_completion_resume_test.go` asserts a linked failed Receipt returns `receipt-unsuccessful` with no Issue/Event mutation.
- VAL-066: PASS - `TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease` compares the successful completion `evaluation` object byte-for-byte as a JSON map to `evidence.evaluate` and completes only after `qa_pass`.
- VAL-067: PASS - `TestEnvelopeCompleteReturnsAllCurrentReasonsWithoutMutation` asserts `obligation-evidence-absent` and `issue-not-terminal-success`, then compares Envelope, Issue, Evidence, legacy-row, and Event snapshots.
- VAL-068: PASS - `TestResumeReadsDurableFactsWithoutMutation` reads the active Envelope/Issue projection from a new CLI command invocation and asserts an unchanged Event count.

## Tests run

- `go test -race ./cmd/fid -run '^TestTerminalSuccessTransitionRejectsUnsuccessfulReceiptEvidence$' -count=1` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 5.231s` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestTerminalSuccessTransitionRejectsUnsuccessfulReceiptEvidence$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 6.244s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease$' -count=1` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 6.611s` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 8.050s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestResumeReadsDurableFactsWithoutMutation$' -count=1` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 1.673s` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestResumeReadsDurableFactsWithoutMutation$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 2.777s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestEnvelopeCompleteReturnsAllCurrentReasonsWithoutMutation$' -count=1` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 2.188s` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestEnvelopeCompleteReturnsAllCurrentReasonsWithoutMutation$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 3.483s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^(TestLiveLeafCommandsMatchContractInventories|TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease|TestEnvelopeCompleteReturnsAllCurrentReasonsWithoutMutation|TestResumeReadsDurableFactsWithoutMutation|TestTerminalSuccessTransitionRejectsUnsuccessfulReceiptEvidence)$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 15.395s` (terminated `exit_code: 0`)
- `/tmp/fid-task18 --project /tmp/FIDtask18vy8P2y --json resume FID-missing` - OBSERVED (non-gate) - `{"contract_version":"fid-cli.v1","error":{"category":"not-found","code":"envelope-not-found","details":{"id":"FID-missing"},"message":"envelope was not found","retryable":false,"subject":{"id":"FID-missing","type":"envelope"}},"ok":false,"operation":"projection.resume"}` (expected read-only absent-Envelope response; terminated `exit_code: 3`)

## Exact quality-gate results

- `go test -race ./internal/workflow ./internal/issues ./internal/evidence ./internal/validation ./internal/cli ./cmd/fid/...` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 270.363s` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid/integration 6.201s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow 25.538s coverage: 29.2% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 62.0%` (terminated `exit_code: 0`)

## Deviations from spec

- First required command failure retained: `go test -race ./internal/workflow ./internal/issues ./internal/evidence ./internal/validation ./internal/cli ./cmd/fid/...` ended with `FAIL` (terminated `exit_code: 1`) because `TestReasonCodeDetailSchemaIsExactForEveryCode` still expected the pre-task-18 closed registry length of 37. The in-scope repair added the accepted `obligation-unsatisfied` schema expectation in `internal/cli/codes_test.go`; the byte-identical rerun passed with `ok github.com/mbell-nzl/fid/cmd/fid 263.427s` (terminated `exit_code: 0`), and the final complete task block passed from the first command again.
- Lint first failed with `* unused: 1` (terminated `exit_code: 1`) because `task18DecodeResult` was unused in the new task test. The repair removed that helper; the byte-identical rerun passed with `0 issues.` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `2313235ba463060538f2f0d4bcfa61f039fcdc13`

## Commit subject

- `feat(workflow): guard terminal completion and resume`
