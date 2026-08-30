# Phase 21 - Task 17: Trusted Evidence Evaluation - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 5
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added `evidence.link`, `evidence.get`, and `evidence.evaluate` leaves with closed inputs and structured fid-cli.v1 results.
- Linked Evidence only to persisted fid-produced Receipt chains whose Obligation, Registration, Run, Issue, and repository target are exact.
- Added the QA-56 Receipt-binding matrix: one retained qualifying link, five persisted-chain mismatch rejections, and every caller-authored proof rejection now prove no Evidence Event or record is created.
- Added pure independent validity, repository-tree freshness, required Artifact availability, and deterministic reason evaluation through the shared workflow evaluator.
- Converged the retained executable-leaf inventory to include the three accepted Evidence leaves.
- Added `TestPhase21AttentionBoundary` as the DEC-172 paired runtime/static VAL-063 proof.

## Files created or modified

- `cmd/fid/contract_surface_test.go`
- `cmd/fid/evidence.go`
- `cmd/fid/evidence_test.go`
- `cmd/fid/root.go`
- `internal/cli/codes.go`
- `internal/evidence/evidence.go`
- `internal/workflow/evaluation.go`
- `project-docs/docs/phases/phase-21/tasks/task-17-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-17-done.md` - records the preserved pre-DEC-172 failure, red/green boundary-test evidence, revised gate results, and real-entrypoint observation.

## Validation Contract Coverage

- VAL-058: PASS - `TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent` in `cmd/fid/evidence_test.go:16` retains the qualifying persisted fid-executed Receipt link and exact producer/trust source/Receipt ID. `TestEvidenceLinkRejectsPersistedReceiptBindingMismatches` in `cmd/fid/evidence_test.go:107` independently rejects altered Obligation, Registration, Run, Issue, and repository-target bindings with typed errors and unchanged journal Event/evidence-record counts. `TestEvidenceLinkRejectsEveryCallerAuthoredProof` at `cmd/fid/evidence_test.go:84` preserves pass, command, import, attestation, and role-looking caller-proof rejection and now proves the same no-Event/no-record result for every row.
- VAL-059: PASS - `TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent` at `cmd/fid/evidence_test.go:72` changes Evidence validity to `invalid` while its fresh/available baseline is separately asserted at lines 58 and 66.
- VAL-060: PASS - `TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent` asserts the original repository target is `fresh` at line 66, then changes only `tree_oid` and observes `stale` at line 80.
- VAL-061: PASS - `TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent` asserts a required retained Artifact is `available` at line 66, then sets only its availability to `missing` before observing the independent `artifact-unavailable` reason at line 80.
- VAL-062: PASS - `TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent` at line 80 asserts the complete fixed-precedence sequence `evidence-invalid`, `evidence-stale`, `artifact-unavailable`.
- VAL-063: PASS - `TestPhase21AttentionBoundary` in `cmd/fid/contract_surface_test.go:22` requires exactly one `attention.v1` row with `available:false`, limits `{"state":false}`, schema `1..1`, and trust source `phase-22-unavailable`; it rejects `fid --json attention list` with unchanged journal/projection and inventories non-test Go plus migrations for exactly that sole production match at lines 79-96 and 100-137.

## Tests run

- `go test -race ./cmd/fid -run '^TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent$' -count=1` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 3.116s` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 4.632s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestEvidenceLinkRejectsEveryCallerAuthoredProof$' -count=1` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 2.858s` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestEvidenceLinkRejectsEveryCallerAuthoredProof$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 4.033s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestEvidenceLinkRejectsPersistedReceiptBindingMismatches$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 14.126s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^(TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestEvidenceLinkRejectsEveryCallerAuthoredProof|TestEvidenceLinkRejectsPersistedReceiptBindingMismatches|TestPhase21AttentionBoundary)$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 19.925s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestLiveLeafCommandsMatchContractInventories$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 1.772s` (terminated `exit_code: 0`)
- `go test -race ./internal/cli -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/cli 1.058s` (terminated `exit_code: 0`)
- `go test -race ./internal/workflow -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow 20.641s` (terminated `exit_code: 0`)
- `go test -race ./internal/evidence -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/evidence 2.270s` (terminated `exit_code: 0`)
- Compile diagnostic preserved from the focused `TestPhase21AttentionBoundary` red run: `cmd/fid/contract_surface_test.go:94:64: missing ',' in composite literal` (terminated `exit_code: 1`).
- `go test -race ./cmd/fid -run '^TestPhase21AttentionBoundary$' -count=1` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 0.955s` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestPhase21AttentionBoundary$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 2.085s` (terminated `exit_code: 0`)
- `/tmp/fid --json attention list` - OBSERVED (non-gate) - `exit 2 (no output)` (expected unavailable no-leaf rejection; terminated `exit_code: 2`). The real compiled CLI rejected the deferred leaf; `TestPhase21AttentionBoundary` separately proves no journal or projection mutation.
- `/tmp/fid --json attention list` - OBSERVED (non-gate) - `exit 2 (no output)` (expected unavailable no-leaf rejection; terminated `exit_code: 2`). The compiled fid entrypoint remained unavailable after the MF-17-01 coverage repair; no runtime product behavior changed.

## Exact quality-gate results

- `go test -race ./internal/evidence ./internal/validation ./internal/workflow ./internal/cli ./cmd/fid/...` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 250.251s` (terminated `exit_code: 0`)
- `rg -n -i --glob '*.go' --glob '!**/*_test.go' '\battention\b|attention_items' cmd/fid internal` - BLOCKED - `internal/workflow/capabilities.go:105: {"attention.v1", false, {"state":false}, "phase-22-unavailable"},` (expected `exit_code: 1`; terminated `exit_code: 0`)
- `gofmt -l cmd internal` - BLOCKED - `not run after protected task-gate contradiction` (not executed)
- `golangci-lint run ./...` - BLOCKED - `not run after protected task-gate contradiction` (not executed)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `not run after protected task-gate contradiction` (not executed)
- `go tool cover -func=coverage.out` - BLOCKED - `not run after protected task-gate contradiction` (not executed)
- `go test -race ./internal/evidence ./internal/validation ./internal/workflow ./internal/cli ./cmd/fid/...` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 262.146s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestPhase21AttentionBoundary$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 2.059s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow 25.116s coverage: 32.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 63.0%` (terminated `exit_code: 0`)
- `go test -race ./internal/evidence ./internal/validation ./internal/workflow ./internal/cli ./cmd/fid/...` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 304.640s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestPhase21AttentionBoundary$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 2.003s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow 29.2% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 62.0%` (terminated `exit_code: 0`)

## Deviations from spec

- DEC-172 replaced only the defective broad expected-nonzero search after TRG-21-019 correctly recorded its protected contradiction. The original failed row above remains historical failed evidence; the distinct named boundary-test command is the authorized replacement gate. No production interface, validation contract, or Attention surface was changed.
- The prior committed-report lint attempt printed `OUTCOME COMPLETE` but terminated `exit_code: 10` with missing canonical workflow fields, an invalid COMPLETE disposition, the CLI exit-2 observation labelled PASS, and the preserved compile diagnostic in place of a terminal go-test token. This report-only correction restores the canonical COMPLETE fields, retains the diagnostic separately from the terminal red/green rows, and records the CLI rejection as an observed non-gate result.
- MF-17-01 required only missing evidence coverage. The new focused matrix passed on its first run with `ok github.com/mbell-nzl/fid/cmd/fid 14.126s` (terminated `exit_code: 0`), confirming the existing exact-binding implementation already rejected all six negative rows; no product-code repair was authorized or needed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `200b807bf67e1bc8113e090383a2548ffbe13a12`

## Commit subject

- `test(evidence): cover receipt binding mismatches`
