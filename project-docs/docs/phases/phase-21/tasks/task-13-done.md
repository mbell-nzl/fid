# Phase 21 - Task 13: Portable Executor Policy - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added the standard-library executor adapter for validated GateDefinitions. It invokes structured argv through `exec.CommandContext`, resolves the declared cwd, assembles only the declared literal, inherited, and `FID_CREDENTIAL_<NAME>` environment entries, and enforces `timeout_ms`.
- Returned local stdout, stderr, and observed exit status only. Live framing, redaction application, Receipt persistence, Artifact retention, OS isolation, and platform-specific containment remain unimplemented as assigned to later tasks.
- Added real-child-process policy tests for literal argv metacharacters, exact cwd, closed environment, credential transfer, declared redaction/Artifact inputs, timeout, and missing credential rejection before spawn.

## Files created or modified

- `internal/evidence/executor.go`
- `internal/evidence/executor_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-13-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-13-done.md`

## Validation Contract Coverage

- VAL-044: PASS - `TestExecuteAppliesAcceptedGatePolicy` at `internal/evidence/executor_test.go:16` runs a real child process and proves structured literal argv, resolved cwd, the complete closed child environment, credential transfer, accepted `sandbox:none`/`network:allowed`/`write:allowed`, configured bounds, redaction input, and safe Artifact declaration. `TestExecuteEnforcesTimeoutLimit` at `internal/evidence/executor_test.go:101` separately proves the timeout bound terminates with `context.DeadlineExceeded`.
- VAL-045: PASS - `TestGateRegisterRejectsUnsupportedPolicyWithoutPersistence` at `cmd/fid/gate_test.go:71` is the one-field Registration matrix for argv, shell, sandbox, network, write, environment, path, bounds, credentials, redaction, and Artifact values; each row asserts `executor-policy-unsupported` and `registrationCount == 0`. The final required package gate includes this test and terminated `exit_code: 0`.
- VAL-046: PASS - the final source inventory returned `exit 1 (no output)` with expected and observed `exit_code: 1`, and the final Phase 21 contract inventory returned the same. The executable policy test at `internal/evidence/executor_test.go:16` carries only the explicit permissive values `sandbox:none`, `network:allowed`, and `write:allowed`.

## Tests run

- `go test ./internal/evidence` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`): `Execute` was undefined after adding the positive policy process probe and missing-credential test.
- `go test ./internal/evidence` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t0.014s` (terminated `exit_code: 0`)
- `go test ./internal/evidence` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`): the timeout probe received `execute gate process: signal: killed` before the adapter exposed the deadline cause.
- `go test ./internal/evidence` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t0.023s` (terminated `exit_code: 0`)
- `go test -race ./internal/evidence` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t2.055s` (terminated `exit_code: 0`)
- Real process adapter exercise: `TestExecuteAppliesAcceptedGatePolicy` at `internal/evidence/executor_test.go:16` spawned the test executable through `Execute`, observed exit code `0`, and decoded the child-reported argv, cwd, and complete environment. There is intentionally no `fid gate run` CLI surface in this task; its framing and Receipt lifecycle are task 14/15 scope.

## Exact quality-gate results

- `go test -race ./internal/evidence ./internal/config ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t244.539s` (terminated `exit_code: 0`)
- `rg -n -i --glob '*.go' --glob '!**/*_test.go' -e 'sandboxed' -e 'seccomp' -e 'chroot' -e 'namespace isolation' -e 'privilege (drop|separation)' -e 'network (blocked|denied|disabled|isolated)' -e 'filesystem (isolation|containment)' -e 'isolated (process|environment|execution)' -e 'contained (process|execution)' cmd internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n -i -e 'sandboxed' -e 'seccomp' -e 'chroot' -e 'namespace isolation' -e 'privilege (drop|separation)' -e 'network (blocked|denied|disabled|isolated)' -e 'filesystem (isolation|containment)' -e 'isolated (process|environment|execution)' -e 'contained (process|execution)' project-docs/docs/phases/phase-21/contracts.md project-docs/docs/phases/phase-21/public-interface.md` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 32.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t\t66.0%` (terminated `exit_code: 0`)

## Deviations from spec

- The first required package gate attempt failed with `FAIL` (terminated `exit_code: 1`): the positive child-process fixture used a 1,000 ms accepted timeout and a race-instrumented child test executable exceeded it. Root cause was test-process startup overhead, not an adapter policy failure. The fixture-only policy timeout was raised to 5,000 ms; the separate 10 ms timeout assertion remains unchanged. The byte-identical required command later passed with `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t262.606s` (terminated `exit_code: 0`), and the final-tree rerun passed with `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t244.539s` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `d2feab0be5c53f46553cc18b79e605e6586da434`

## Commit subject

- `feat(evidence): add portable executor policy`
