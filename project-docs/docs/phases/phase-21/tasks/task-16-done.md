# Phase 21 - Task 16: Safe Referenced Artifacts - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Retained declared regular files into `.fid/artifacts` by the SHA-256 digest of sanitized content, with staging-file fsync and atomic rename publication.
- Rejected missing, directory, symlink, non-regular, out-of-CWD, oversized, and changed-during-copy candidates before they can create an Artifact or Receipt link.
- Added `artifact verify`, which records loss, restoration, and unchanged availability as idempotent observation Events while retaining Artifact and Receipt identity.

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

- `project-docs/docs/phases/phase-21/tasks/task-16-done.md` - recorded task implementation and deterministic validation evidence.

## Validation Contract Coverage

- VAL-055: PASS - `TestGateRunPublishesSanitizedArtifactByDigest` in `cmd/fid/artifact_test.go:18` proves redaction-before-digest, content-addressed persistence, Artifact metadata, Receipt linkage, and no leftover staging file.
- VAL-056: PASS - `TestArtifactCandidatePathRejectsEveryUnsafeCandidate` in `internal/evidence/artifact_test.go:10` covers missing, directory, symlink, non-regular FIFO, and containment escape rejection; `TestPublishArtifactCandidateRejectsChangedDuringCopy` at `internal/evidence/artifact_test.go:52` covers revalidation during copying; `TestPublishArtifactCandidateRejectsOversizedContent` at `internal/evidence/artifact_test.go:75` covers the size limit. These rejection paths return no retained Artifact, so they create no Receipt link that could satisfy required material.
- VAL-057: PASS - `TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation` in `cmd/fid/artifact_test.go:71` proves both `missing` and restored `available` observations preserve Artifact and Receipt identity, and its check at `cmd/fid/artifact_test.go:145` proves an unchanged verification emits only the idempotency observation effect.

## Tests run

- `go test -race ./cmd/fid -run '^TestGateRunPublishesSanitizedArtifactByDigest$' -count=1` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 3.329s` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestGateRunPublishesSanitizedArtifactByDigest$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 4.485s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `printf '%s\\n' '{"artifact_id":"ART-missing"}' | /tmp/fid --project . --json artifact verify --actor actor-builder --input - --idempotency-key artifact-runtime-authority-final` (working directory `/tmp/fid-artifact-runtime-svvcFn`) - PASS - `{"contract_version":"fid-cli.v1","error":{"category":"authority","code":"authority-unavailable","details":{"actor_id":null,"operation":"artifact.verify","required_actions":["artifact.verify"],"required_capabilities":["artifact.verify"]},"message":"finalized activation authority is unavailable","retryable":false,"subject":null},"ok":false,"operation":"artifact.verify"}` (terminated `exit_code: 5`; expected typed authority rejection before activation)

## Exact quality-gate results

- `go test -race ./internal/evidence ./internal/config ./internal/journal ./internal/cli ./cmd/fid/...` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 292.644s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 32.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 64.2%` (terminated `exit_code: 0`)

## Deviations from spec

- Required scoped-gate recovery: the first `go test -race ./internal/evidence ./internal/config ./internal/journal ./internal/cli ./cmd/fid/...` attempt ended `FAIL github.com/mbell-nzl/fid/cmd/fid 285.184s` (terminated `exit_code: 1`) because the retained command inventory had not yet converged on the accepted `artifact verify` leaf, and a signaled fixture timed out in the same run. Repaired the accepted leaf inventory, reproduced the signaled fixture with `-count=5`, then reran the byte-identical required command successfully as recorded above.
- Required coverage-gate recovery: the first `go test -race -coverprofile=coverage.out ./...` attempt ended `FAIL` (terminated `exit_code: 1`) after crossing the UTC journal-date boundary from 2026-08-29 to 2026-08-30, exposing inherited single-journal-file assumptions and a transient signaled-fixture timeout. Focused reproduction passed after the date stabilized; the byte-identical rerun passed as recorded above.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
