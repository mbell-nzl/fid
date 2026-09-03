# Phase 21 - Task 13: Portable Executor Policy - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: qa-rework
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added the standard-library executor adapter for validated GateDefinitions. It invokes structured argv through `exec.CommandContext`, resolves the declared cwd, assembles only the declared literal, inherited, and `FID_CREDENTIAL_<NAME>` environment entries, and enforces `timeout_ms`.
- Returned local stdout, stderr, and observed exit status only. Live framing, redaction application, Receipt persistence, Artifact retention, OS isolation, and platform-specific containment remain unimplemented as assigned to later tasks.
- Added real-child-process policy tests for literal argv metacharacters, exact cwd, closed environment, credential transfer, declared redaction/Artifact inputs, timeout, and missing credential rejection before spawn.
- Final-review MF-2 recovery removed the two unsupported executor-limit keys from the section 2 `config.toml` example. The example now names exactly the implemented config keys, `max_output_bytes` and `max_excerpt_bytes`, without changing rules, bounds, product code, or tests.

## Files created or modified

- `internal/evidence/executor.go`
- `internal/evidence/executor_test.go`
- `project-docs/specs/04_TRD_CONFIG.md`
- `project-docs/docs/phases/phase-21/tasks/task-13-done.md`

## Documentation updated

- `project-docs/specs/04_TRD_CONFIG.md` section 2 example block only
- `project-docs/docs/phases/phase-21/tasks/task-13-done.md`

## Validation Contract Coverage

- VAL-044: PASS - `TestExecuteAppliesAcceptedGatePolicy` at `internal/evidence/executor_test.go:16` runs a real child process and proves structured literal argv, resolved cwd, the complete closed child environment, credential transfer, accepted `sandbox:none`/`network:allowed`/`write:allowed`, configured bounds, redaction input, and safe Artifact declaration. `TestExecuteEnforcesTimeoutLimit` at `internal/evidence/executor_test.go:101` separately proves the timeout bound terminates with `context.DeadlineExceeded`.
- VAL-045: PASS - `TestGateRegisterRejectsUnsupportedPolicyWithoutPersistence` at `cmd/fid/gate_test.go:71` is the one-field Registration matrix for argv, shell, sandbox, network, write, environment, path, bounds, credentials, redaction, and Artifact values; each row asserts `executor-policy-unsupported` and `registrationCount == 0`. The final required package gate includes this test and terminated `exit_code: 0`.
- VAL-046: PASS - the final source inventory returned `exit 1 (no output)` with expected and observed `exit_code: 1`, and the final Phase 21 contract inventory returned the same. The executable policy test at `internal/evidence/executor_test.go:16` carries only the explicit permissive values `sandbox:none`, `network:allowed`, and `write:allowed`.
- No additional VAL applies to the final-review MF-2 documentation recovery; the amended task specification identifies it as no-VAL QA-91 evidence. VAL-002 and its initial final-review BLOCK remain unchanged and are not part of this task-scope repair.

## Tests run

- `go test ./internal/evidence` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`): `Execute` was undefined after adding the positive policy process probe and missing-credential test.
- `go test ./internal/evidence` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t0.014s` (terminated `exit_code: 0`)
- `go test ./internal/evidence` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`): the timeout probe received `execute gate process: signal: killed` before the adapter exposed the deadline cause.
- `go test ./internal/evidence` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t0.023s` (terminated `exit_code: 0`)
- `go test -race ./internal/evidence` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t2.055s` (terminated `exit_code: 0`)
- Real process adapter exercise: `TestExecuteAppliesAcceptedGatePolicy` at `internal/evidence/executor_test.go:16` spawned the test executable through `Execute`, observed exit code `0`, and decoded the child-reported argv, cwd, and complete environment. There is intentionally no `fid gate run` CLI surface in this task; its framing and Receipt lifecycle are task 14/15 scope.
- Documented-example red reproduction (the exact task-spec command against the uncorrected section 2 example) - FAIL (expected red) - `{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["accepted issue.create input"],"argument":"request","value":null},"message":"invalid Issue mutation request","retryable":false,"subject":null},"ok":false,"operation":"issue.create"}` (terminated `exit_code: 2`)
- Documented-example green rerun (the byte-identical command after the bounded section 2 edit) - PASS - `{"contract_version":"fid-cli.v1","event":{"id":"EVT-20260903-000002","occurred_at":"2026-09-03T19:40:45.577017366Z"},"ok":true,"operation":"issue.create","result":{"idempotency_key":"k1","issue":{"id":"FID-7014a5","kind":"task","revision":1,"status":"todo","title":"T"},"replayed":false}}` (terminated `exit_code: 0`)

## Exact quality-gate results

- `bash -c 'set -euo pipefail; repo=$(pwd); scratch=$(mktemp -d); go build -o "$scratch/fid" ./cmd/fid; mkdir -p "$scratch/workspace"; cd "$scratch/workspace"; "$scratch/fid" init --prefix FID >/dev/null; awk "/^## 2[.] /{section=1; next} section && /^...toml$/{fence=1; next} fence && /^...$/{exit} fence{print}" "$repo/project-docs/specs/04_TRD_CONFIG.md" > .fid/config.toml; "$scratch/fid" --json issue create --kind task --title T --actor agent:dev --idempotency-key k1'` - PASS - `{"contract_version":"fid-cli.v1","event":{"id":"EVT-20260903-000002","occurred_at":"2026-09-03T19:41:02.336233642Z"},"ok":true,"operation":"issue.create","result":{"idempotency_key":"k1","issue":{"id":"FID-08e100","kind":"task","revision":1,"status":"todo","title":"T"},"replayed":false}}` (terminated `exit_code: 0`)
- `go test -race ./internal/evidence ./internal/config ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t244.539s` (terminated `exit_code: 0`)
- `rg -n -i --glob '*.go' --glob '!**/*_test.go' -e 'sandboxed' -e 'seccomp' -e 'chroot' -e 'namespace isolation' -e 'privilege (drop|separation)' -e 'network (blocked|denied|disabled|isolated)' -e 'filesystem (isolation|containment)' -e 'isolated (process|environment|execution)' -e 'contained (process|execution)' cmd internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n -i -e 'sandboxed' -e 'seccomp' -e 'chroot' -e 'namespace isolation' -e 'privilege (drop|separation)' -e 'network (blocked|denied|disabled|isolated)' -e 'filesystem (isolation|containment)' -e 'isolated (process|environment|execution)' -e 'contained (process|execution)' project-docs/docs/phases/phase-21/contracts.md project-docs/docs/phases/phase-21/public-interface.md` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 32.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t\t66.0%` (terminated `exit_code: 0`)

Recovery attempt 1 required-command results, executed in task-spec order:

- `go test -race ./internal/evidence ./internal/config ./internal/cli ./cmd/fid/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	356.757s` (terminated `exit_code: 0`)
- `rg -n -i --glob '*.go' --glob '!**/*_test.go' -e 'sandboxed' -e 'seccomp' -e 'chroot' -e 'namespace isolation' -e 'privilege (drop|separation)' -e 'network (blocked|denied|disabled|isolated)' -e 'filesystem (isolation|containment)' -e 'isolated (process|environment|execution)' -e 'contained (process|execution)' cmd internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n -i -e 'sandboxed' -e 'seccomp' -e 'chroot' -e 'namespace isolation' -e 'privilege (drop|separation)' -e 'network (blocked|denied|disabled|isolated)' -e 'filesystem (isolation|containment)' -e 'isolated (process|environment|execution)' -e 'contained (process|execution)' project-docs/docs/phases/phase-21/contracts.md project-docs/docs/phases/phase-21/public-interface.md` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 27.4% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)				61.2%` (terminated `exit_code: 0`)

## Deviations from spec

- The first required package gate attempt failed with `FAIL` (terminated `exit_code: 1`): the positive child-process fixture used a 1,000 ms accepted timeout and a race-instrumented child test executable exceeded it. Root cause was test-process startup overhead, not an adapter policy failure. The fixture-only policy timeout was raised to 5,000 ms; the separate 10 ms timeout assertion remains unchanged. The byte-identical required command later passed with `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t262.606s` (terminated `exit_code: 0`), and the final-tree rerun passed with `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t244.539s` (terminated `exit_code: 0`).
- Final-review MF-2 recovery first reproduced `{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["accepted issue.create input"],"argument":"request","value":null},"message":"invalid Issue mutation request","retryable":false,"subject":null},"ok":false,"operation":"issue.create"}` (terminated `exit_code: 2`) from the exact documented-example command. Root cause: the section 2 example named unsupported `max_timeout_ms` and `max_artifact_bytes` while the shipped parser accepts only `max_output_bytes` and `max_excerpt_bytes`. Bounded repair: remove exactly those two example lines. The byte-identical green rerun produced `{"contract_version":"fid-cli.v1","event":{"id":"EVT-20260903-000002","occurred_at":"2026-09-03T19:40:45.577017366Z"},"ok":true,"operation":"issue.create","result":{"idempotency_key":"k1","issue":{"id":"FID-7014a5","kind":"task","revision":1,"status":"todo","title":"T"},"replayed":false}}` (terminated `exit_code: 0`); the same command also passed as the first final required gate after the durable config commit.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `abf6614e0d02e8e713d5df478284a26ca759c260`

## Commit subject

- `docs(config): correct executor limits example`
