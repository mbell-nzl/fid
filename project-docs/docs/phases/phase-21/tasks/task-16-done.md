# Phase 21 - Task 16: Safe Referenced Artifacts - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Retained declared regular files into `.fid/artifacts` by the SHA-256 digest of sanitized content, with staging-file fsync and atomic rename publication.
- Rejected missing, directory, symlink, non-regular, out-of-CWD, oversized, and changed-during-copy candidates before they can create an Artifact or Receipt link.
- Added `artifact verify`, which records loss, restoration, and unchanged availability as idempotent observation Events while retaining Artifact and Receipt identity.
- QA rework MF-16-01 adds observable gate CLI coverage for every unsafe candidate plus an Evidence-linked loss/restore lifecycle snapshot.

## Files created or modified

- `cmd/fid/artifact.go`
- `cmd/fid/artifact_test.go`
- `cmd/fid/contract_surface_test.go`
- `cmd/fid/gate.go`
- `cmd/fid/root.go`
- `internal/evidence/artifact.go`
- `internal/evidence/artifact_test.go`
- `internal/evidence/execution.go`
- `internal/evidence/live.go`
- `internal/evidence/policy.go`
- `internal/evidence/policy_test.go`
- `internal/journal/event.go`
- `internal/journal/idempotency.go`
- `internal/journal/writer.go`
- `project-docs/docs/phases/phase-21/tasks/task-16-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-16-done.md` - retains the original task evidence and records QA-rework red/green, CLI, and required-gate evidence.

## Validation Contract Coverage

- VAL-055: PASS - `TestGateRunPublishesSanitizedArtifactByDigest` in `cmd/fid/artifact_test.go:21` proves redaction-before-digest, content-addressed persistence, Artifact metadata, Receipt linkage, and no leftover staging file.
- VAL-056: PASS - `TestGateRunUnsafeArtifactsEmitControlsAndRemainRequiredUnavailable` in `cmd/fid/artifact_test.go:74` executes the real `gate run` CLI path for `changed-during-copy`, `directory`, `missing`, `non-regular`, `outside-cwd`, `oversized`, and `symlink`. Its assertions at lines 179-184 require the exact `artifact-rejected` control, zero Receipt Artifact links, and an `unsatisfied` Evidence evaluation with missing availability for every row.
- VAL-057: PASS - `TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation` at `cmd/fid/artifact_test.go:252` links qualifying Evidence before availability changes (lines 274-288), snapshots Artifact/Receipt/Evidence identity (lines 289-291 and 386-398), and proves `missing` availability with `artifact-unavailable` before restoration to `available` (lines 307-328).

## Tests run

- `go test -race ./cmd/fid -run '^TestGateRunPublishesSanitizedArtifactByDigest$' -count=1` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 3.329s` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestGateRunPublishesSanitizedArtifactByDigest$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 4.485s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `printf '%s\n' '{"artifact_id":"ART-missing"}' | /tmp/fid --project . --json artifact verify --actor actor-builder --input - --idempotency-key artifact-runtime-authority-final` (working directory `/tmp/fid-artifact-runtime-svvcFn`) - OBSERVED (non-gate) - `{"contract_version":"fid-cli.v1","error":{"category":"authority","code":"authority-unavailable","details":{"actor_id":null,"operation":"artifact.verify","required_actions":["artifact.verify"],"required_capabilities":["artifact.verify"]},"message":"finalized activation authority is unavailable","retryable":false,"subject":null},"ok":false,"operation":"artifact.verify"}` (expected typed authority rejection before activation; terminated `exit_code: 5`)
- `go test -race ./cmd/fid -run '^TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation$' -count=1` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 3.048s` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation$' -count=1` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 3.408s` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 5.236s` (terminated `exit_code: 0`)
- Compile diagnostic preserved from the unsafe-matrix red run: `cmd/fid/artifact_test.go:181:4: undefined: assertArtifactRejectedControl` (terminated `exit_code: 1`).
- `go test -race ./cmd/fid -run '^TestGateRunUnsafeArtifactsEmitControlsAndRemainRequiredUnavailable$' -count=1` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 18.686s` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestGateRunUnsafeArtifactsEmitControlsAndRemainRequiredUnavailable$' -count=1` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 20.232s` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestGateRunUnsafeArtifactsEmitControlsAndRemainRequiredUnavailable$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 29.409s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^(TestGateRunUnsafeArtifactsEmitControlsAndRemainRequiredUnavailable|TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation)$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 32.641s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid --project /tmp/fid-artifact-runtime-aRh3Yi --json artifact verify --actor actor-builder --input /tmp/fid-artifact-runtime-aRh3Yi/artifact-verify.json --idempotency-key artifact-runtime-authority-rework` - OBSERVED (non-gate) - `{"contract_version":"fid-cli.v1","error":{"category":"authority","code":"authority-unavailable","details":{"actor_id":null,"operation":"artifact.verify","required_actions":["artifact.verify"],"required_capabilities":["artifact.verify"]},"message":"finalized activation authority is unavailable","retryable":false,"subject":null},"ok":false,"operation":"artifact.verify"}` (expected typed authority rejection before activation; terminated `exit_code: 5`)

## Exact quality-gate results

Historical original task-close gate:

- `go test -race ./internal/evidence ./internal/config ./internal/journal ./internal/cli ./cmd/fid/...` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 292.644s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 32.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 64.2%` (terminated `exit_code: 0`)

QA-rework attempt 1 task gate:

- `go test -race ./internal/evidence ./internal/config ./internal/journal ./internal/cli ./cmd/fid/...` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 290.801s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 29.2% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 62.0%` (terminated `exit_code: 0`)

## Deviations from spec

- Original required scoped-gate recovery retained: the first `go test -race ./internal/evidence ./internal/config ./internal/journal ./internal/cli ./cmd/fid/...` attempt ended `FAIL github.com/mbell-nzl/fid/cmd/fid 285.184s` (terminated `exit_code: 1`) because the retained command inventory had not yet converged on the accepted `artifact verify` leaf, and a signaled fixture timed out in the same run. The accepted leaf inventory was repaired, the signaled fixture was reproduced with `-count=5`, and the byte-identical command later passed as recorded above.
- Original required coverage-gate recovery retained: the first `go test -race -coverprofile=coverage.out ./...` attempt ended `FAIL` (terminated `exit_code: 1`) after crossing the UTC journal-date boundary from 2026-08-29 to 2026-08-30, exposing inherited single-journal-file assumptions and a transient signaled-fixture timeout. Focused reproduction passed after the date stabilized; the byte-identical rerun passed as recorded above.
- No product, public-interface, validation-contract, or dependency deviation. The QA-rework changes only observable test evidence within task-16 scope.
- Focused lifecycle red evidence first ended with `FAIL github.com/mbell-nzl/fid/cmd/fid 3.048s` (terminated `exit_code: 1`) because the retained fixture had no Evidence row. The repair links qualifying Evidence before deletion. A second focused attempt ended with `FAIL github.com/mbell-nzl/fid/cmd/fid 3.408s` (terminated `exit_code: 1`) because the ordinary envelope fixture lacks `evidence.link` authority; the causal fixture repair uses the existing Evidence-enabled envelope, and the byte-identical focused command then passed.
- The unsafe-matrix red attempts exposed missing test helpers, a pre-retention escaping-symlink policy rejection, and a helper file-creation defect. The test-only helper now creates the candidate or escaping parent after definition validation and mutates the regular file during retention; the byte-identical matrix command passed with `ok github.com/mbell-nzl/fid/cmd/fid 29.409s` (terminated `exit_code: 0`). No production code changed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `457366d63e9e365e90918237224d54e512050a74`

## Commit subject

- `test(evidence): cover artifact rejection lifecycle`
