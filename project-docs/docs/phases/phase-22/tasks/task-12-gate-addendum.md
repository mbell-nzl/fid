# Phase 22 - Task 12: Developer gate receipt addendum

This addendum records only termination results observed by the Task 12
Developer in the recovery continuation. It does not replace the completion
report or re-attest Root-observed commands.

## Final ordered source-bound gate

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 52.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					62.4%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration	51.743s` (terminated `exit_code: 0`)

## Retained continuation attempts

- Session `90704`: `go test -race -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`). The command exposed retained synthetic terminal fixtures after current-target observation; its output included `cmd/fid` failures in Stop and terminal Evidence matrix cases.
- Session `53773`: `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow	89.215s	coverage: 52.9% of statements` (terminated `exit_code: 0`). This receipt was invalidated by the subsequent accepted DEC-201 replay-source handback.
- Session `29973`: `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow	85.521s	coverage: 52.9% of statements` (terminated `exit_code: 0`). This receipt was invalidated by the subsequent retained Phase 21 real-Git fixture repair.
- Session `83059`: the final source-bound `go test -race -coverprofile=coverage.out ./...` receipt is the cached 52.9% token in the final ordered gate above.
- `golangci-lint run ./...` - FAIL - `* unused: 1` (terminated `exit_code: 1`). The observed-Git fixture convergence left an unused test-only wrapper; removing that obsolete helper preserved every assertion. The byte-identical lint rerun passed with `0 issues.` (terminated `exit_code: 0`).
- `go test -tags integration -race ./cmd/fid/integration/...` - FAIL - `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration	49.665s` (terminated `exit_code: 1`). `TestPhase21RealProcessJourneyStartsOneActivatedEnvelope` retained a synthetic repository target, which current observation correctly rejected.
- The byte-identical tagged causal rerun after that fixture used a fresh Git repository with ignored generated inputs and artifacts - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration	51.320s` (terminated `exit_code: 0`).

## Focused behavior receipts

- `go test -tags integration ./cmd/fid/integration -run '^TestP22RecoveryAdoptsChangedTargetForSuccessorEvidence$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): the target-A Gate charge left the fixture's Run version stale for recovery.
- The same targeted command - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration	1.381s` (terminated `exit_code: 0`).
- `go test -tags integration ./cmd/fid/integration -run '^TestP22RecoveryAdoptedTargetPermitsTerminalTransitionWithFreshSuccessorProof$' -count=1` - FAIL - `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration	1.708s` (terminated `exit_code: 1`): actual target-B `issue transition` returned `run-inactive` because the active Run was compared with the original Envelope target.
- `go test ./cmd/fid -run '^TestExecutionBoundTransitionPredicateMatrix/run_repository_project_mismatched$' -count=1` - FAIL - `FAIL	github.com/mbell-nzl/fid/cmd/fid	0.407s` (terminated `exit_code: 1`): project identity was not retained by the interim worktree-only predicate.
- `go test ./cmd/fid -run '^TestExecutionBoundTransitionPredicateMatrix/(run_repository_mismatched|run_repository_project_mismatched)$' -count=1 && go test -tags integration ./cmd/fid/integration -run '^TestP22RecoveryAdoptedTargetPermitsTerminalTransitionWithFreshSuccessorProof$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration	1.971s` (terminated `exit_code: 0`).

## Construction-only attempts

The successor target-B journey had later test-construction failures around an
invalid manually supplied observation time and an unnecessary generic Claim
release of the recovery Claim. They were removed while preserving the actual
CLI terminal Issue transition, Run close, and Envelope completion. Their
details are UNVERIFIED in this addendum because no retained receipt token is
needed to support the final accepted behavior.
