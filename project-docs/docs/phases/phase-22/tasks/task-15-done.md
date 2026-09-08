# Phase 22 Task 15 completion report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Recovery attempts: 1
- Human attention required: no
- Condition: None. Root observed all final required gates terminate successfully.
- Next action: None.

## Summary of what was implemented

Exposed `fid export import` with required destination and idempotency key,
null Event success and the exact non-retryable import-target-conflict shape.
Added real built-CLI export/validate/import/rebuild/resume and destination
successor execution journeys. The portable corpus exercises eight states:
active, interrupted and abandoned Runs; deduplicated Attention; exhausted and
unknown Budgets; partial activation; and missing Artifacts.

Every portable fixture compares complete named live/rebuilt/imported views and
canonical re-export entry bytes. The same expected semantic state is checked
against fid and a separate in-memory logical Effect interpreter. Valid-hash
negative cases reject unknown table/column/operation, retired Event type and
duplicate Event identity on both implementations. The fake reads fixture
archive bytes, never fid storage or production appliers.

## Files created or modified

- `cmd/fid/export.go`, `export_test.go`, `capabilities_test.go`, `contract_surface_test.go`.
- `internal/cli/codes.go`, `codes_test.go`.
- `cmd/fid/integration/p22_import_recovery_test.go`, the five `p22_portable_*_test.go` helpers, shared fake helpers and `testdata/p22_conformance/`.
- This report and `task-15-root-evidence.md`.

## Documentation updated

This report and root evidence preserve implementation, interpretation, failed
attempts and convergence provenance. No spec change: DEC-202/203 and Task 15
already assign the accepted behavior. Root separately maintains the ledger.

## Validation Contract Coverage

- VAL-063: PASS. `TestP22PortableRecoveryConformance` compares full named views for all eight states; ordered Event bytes and entity identities are retained. The shared fake independently reconstructs semantic state.
- VAL-064: PASS. Every canonical archive entry is byte-compared across live/rebuilt/imported re-exports. Actual destination observation is separate from controlled-input durable parity.
- VAL-065: PASS. `TestRecoveryLiveRebuildImportParity` and the eight-state runner execute the built CLI; destination recovery produces a new Claim/Run and a real Gate Receipt.

## Tests run

Root's directly observed failures/repairs and focused greens are retained in
`task-15-root-evidence.md`. No final required command is inferred from them.

Developer hand-back attribution (not root fresh execution): the first named
parity test failed behaviorally because import was absent (CLI exit 2; test
exit 1). Its implementation then passed the active parity/destination journey
and shared 41-case corpus. The Developer also observed an open conflict-reason
domain red, repaired its enum check, and converged two displaced CLI fixtures.
It handed back without a report, commit or final five-command gate. Root's
same-task continuation completes the omitted portable state corpus and gates.

## Exact quality-gate results

Root directly observed every command below terminate, in this order, from
`/home/mark/projects/fid`. The first failed full race attempt, its verbatim
failure token/exit status and causal fixture repair remain in root evidence.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 52.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					63.2%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	72.834s` (terminated `exit_code: 0`)

## Deviations from spec

Root took over directly under the Human's super-delegate authorization after
an incomplete Developer hand-back. No concurrent product writer was active.
The fake Budget decoding and partial-activation fixture preparation failures
were test-construction errors, repaired without changing expected outcomes or
product policy. Every failed receipt remains in root evidence.

The earlier Developer did not leave the required report; its observations are
explicitly attributed rather than reconstructed as root gate evidence. Root
records all final gates in this session. Subject-preserving convergence and
introducing commits are in root evidence; ordinary and tagged gates execute
the affected assertions. No dependency, VAL, command or threshold changed.

## Concerns, technical debt, or notes for the Tech Lead

Independent Review 5 follows this task. Human judgments remain pending. The
portable fake pins schema 15 columns; Task 16 owns schema-16 convergence.

## Dependency-budget note

- None. No dependency was added or updated.
- The initial root report-lint attempt failed with `FAIL required heading missing: ## Dependency-budget note` (terminated `exit_code: 2`). The report already stated no dependency change; root added this required heading and sentinel. Product and gate evidence are unchanged.

## Commit hash

- `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`

## Commit subject

- `feat(export): expose import and prove portable recovery parity`
