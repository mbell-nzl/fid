# Phase 21 - Task 17: Trusted Evidence Evaluation - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 6
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added `evidence.link`, `evidence.get`, and `evidence.evaluate` leaves with closed inputs and structured fid-cli.v1 results.
- Linked Evidence only to persisted fid-produced Receipt chains whose Obligation, Registration, Run, Issue, and repository target are exact.
- Added the QA-56 Receipt-binding matrix: one retained qualifying link, five persisted-chain mismatch rejections, and every caller-authored proof rejection now prove no Evidence Event or record is created.
- Added independent paired validity, repository-tree freshness, and required Artifact availability fixtures for Receipt-binding evaluation, including a `head_oid` branch-label control that remains fresh when the repository tree is unchanged.
- Added a canonical ordered multi-failure Evaluation golden that is byte-identical on repeat and after rebuilding the projection from its journal.
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

- `project-docs/docs/phases/phase-21/tasks/task-17-done.md` - records the preserved pre-DEC-172 failure, MF-17-01 evidence, MF-17-02 independent-dimension and repeat/rebuild evidence, revised gate results, and real-entrypoint observation.

## Validation Contract Coverage

- VAL-058: PASS - `TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent` in `cmd/fid/evidence_test.go:17` retains the qualifying persisted fid-executed Receipt link and exact producer/trust source/Receipt ID. `TestEvidenceLinkRejectsPersistedReceiptBindingMismatches` in `cmd/fid/evidence_test.go:203` independently rejects altered Obligation, Registration, Run, Issue, and repository-target bindings with typed errors and unchanged journal Event/evidence-record counts. `TestEvidenceLinkRejectsEveryCallerAuthoredProof` at `cmd/fid/evidence_test.go:180` preserves pass, command, import, attestation, and role-looking caller-proof rejection and now proves the same no-Event/no-record result for every row.
- VAL-059: PASS - `TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent` in `cmd/fid/evidence_test.go:18` evaluates the qualifying valid/fresh/available Receipt link, changes only persisted Evidence validity at line 21, then requires the sole `evidence-invalid` reason at line 22.
- VAL-060: PASS - the same test at `cmd/fid/evidence_test.go:25` retains the original tree as fresh, changes only `head_oid` at lines 29-31 while remaining fresh, then changes only `tree_oid` at lines 33-35 and requires the sole `evidence-stale` reason.
- VAL-061: PASS - the same test at `cmd/fid/evidence_test.go:38` evaluates a required Artifact as available, changes only its observed availability at line 41, then requires the sole `artifact-unavailable` reason at line 42.
- VAL-062: PASS - `TestEvidenceEvaluationMultiFailureGoldenSurvivesRepeatAndRebuild` in `cmd/fid/evidence_test.go:46` requires the full two-reason fixed-precedence object `evidence-stale`, `artifact-unavailable` at line 52, byte-compares the canonical Evaluation on repeat at lines 53-56, rebuilds at lines 58-61, and byte-compares the post-rebuild result at lines 62-64.
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
- `go test -race ./cmd/fid -run '^TestEvidenceEvaluationMultiFailureGoldenSurvivesRepeatAndRebuild$' -count=1` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 3.661s` (terminated `exit_code: 1`). The initial golden fixture decoded a bare Evaluation as a response envelope; the failure was preserved and the assertion was corrected without changing product behavior.
- `go test -race ./cmd/fid -run '^(TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestEvidenceEvaluationMultiFailureGoldenSurvivesRepeatAndRebuild)$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 14.123s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^(TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestEvidenceEvaluationMultiFailureGoldenSurvivesRepeatAndRebuild|TestEvidenceLinkRejectsEveryCallerAuthoredProof|TestEvidenceLinkRejectsPersistedReceiptBindingMismatches|TestPhase21AttentionBoundary)$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 29.303s` (terminated `exit_code: 0`)
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-17-done.md project-docs/docs/phases/phase-21/tasks/task-17-evidence-evaluation.md` - FAIL - `FAIL report-lint found 2 failure(s)` (terminated `exit_code: 2`). The report control rejected trailing punctuation after the observed termination fields on the two preceding PASS rows.

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
- `go test -race ./internal/evidence ./internal/validation ./internal/workflow ./internal/cli ./cmd/fid/...` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 337.152s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestPhase21AttentionBoundary$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 1.996s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow 29.2% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 62.0%` (terminated `exit_code: 0`)

## Deviations from spec

- DEC-172 replaced only the defective broad expected-nonzero search after TRG-21-019 correctly recorded its protected contradiction. The original failed row above remains historical failed evidence; the distinct named boundary-test command is the authorized replacement gate. No production interface, validation contract, or Attention surface was changed.
- The prior committed-report lint attempt printed `OUTCOME COMPLETE` but terminated `exit_code: 10` with missing canonical workflow fields, an invalid COMPLETE disposition, the CLI exit-2 observation labelled PASS, and the preserved compile diagnostic in place of a terminal go-test token. This report-only correction restores the canonical COMPLETE fields, retains the diagnostic separately from the terminal red/green rows, and records the CLI rejection as an observed non-gate result.
- MF-17-01 required only missing evidence coverage. The new focused matrix passed on its first run with `ok github.com/mbell-nzl/fid/cmd/fid 14.126s` (terminated `exit_code: 0`), confirming the existing exact-binding implementation already rejected all six negative rows; no product-code repair was authorized or needed.
- MF-17-02 required test/evidence completion only. The independent fixtures confirmed the existing evaluator already uses project/worktree/tree state rather than `head_oid` for freshness. No production code changed, and the deterministic golden persists the Artifact-missing state through rebuild by using the existing `artifact.verify` observation Event.
- The first post-repair report-lint attempt terminated `exit_code: 2` with `FAIL report-lint found 2 failure(s)` because two focused PASS rows ended their termination evidence with a trailing period. This report-only correction preserves the failed attempt and canonicalizes those rows; no test, command, assertion, or product behavior changed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `83706d373e2ecf0d6a668b8d0d489c185eea03c1`

## Commit subject

- `test(evidence): prove independent evaluation dimensions`
