# Phase 22 Task 14 completion report

## Workflow outcome

- Task outcome: COMPLETE

## Summary of what was implemented

- Added the domain import API. It snapshots and validates the supplied archive, stages configuration, journal and artifacts, rebuilds the SQLite projection with normal appliers, fsyncs the staged tree, and atomically publishes `.fid` with no replacement.
- Added durable delivery metadata and idempotent replay. A verified installed binding returns the original result only for the same archive digest and idempotency key; populated, unsafe, changed-key and malformed installed targets are rejected.
- Preserved imported Event bytes/order/identities, extended ordinal journal handling for imported history, and prevented new locally written Event IDs from colliding with imported ordinal chunks.
- Added DEC-202 destination fencing: imported observations and historical execution cannot establish destination authority; a fresh local observation and reconciliation produce destination-bound successor Claim/Run, whose Gate produces a Receipt. DEC-204 also fences CloseRun and AbandonRun against imported observations.
- Added archive semantic validation for nested observation identity and aggregate-to-immutable-Effect identity, covering Task 13 provenance findings.

## Files created or modified

- `internal/export/import.go`, `internal/export/import_test.go`, `internal/export/export_test.go`, `internal/export/archive_references.go`, and `internal/export/archive_corruption_test.go`
- `internal/config/import_binding.go`
- `internal/journal/replay.go` and `internal/journal/writer.go`
- `internal/workflow/import_binding.go`, `internal/workflow/reconciliation.go`, `internal/workflow/recovery.go`, `internal/workflow/run_close.go`, and `internal/workflow/run_inspect.go`
- `internal/evidence/execution.go` and `internal/claims/generic.go`

## Documentation updated

- Added this completion report. The DEC-202 and DEC-204 allocations were already recorded in the phase task and decision documents by the Phase Runner.

## Validation Contract Coverage

- VAL-058: `TestImportPublishesIdentityAndReplaysLostResponse` at `internal/export/import_test.go:38-59` verifies stable Project identity, source preservation, installed binding, ordered Event bytes, and no mutation on replay. `TestArchiveRejectsIntegrityValidSemanticCorruptions` at `internal/export/archive_corruption_test.go:58-77` rejects the valid-hash nested-observation and aggregate identity forgeries.
- VAL-059: `TestImportCrashBeforePublicationLeavesNoPartialProject` at `internal/export/export_test.go:74` captures the pre-publication fault boundary; `TestImportAfterPublicationFaultReplaysInstalledResult` at `internal/export/import_test.go:105-128` proves the after-publication/lost-response boundary replays the installed result.
- VAL-060: `TestImportPublishesIdentityAndReplaysLostResponse` at `internal/export/import_test.go:53-64` proves same digest/key replay with unchanged installed facts and rejects a different key. `TestImportAfterPublicationFaultReplaysInstalledResult` at `internal/export/import_test.go:116-128` proves retry after a publication response loss.
- VAL-061: `TestImportRejectsExistingAndUnsafeTargetsWithoutChangingThem` at `internal/export/import_test.go:68-103` rejects populated and symlink targets and proves the populated marker is unchanged. The identity/replay test preserves an existing target source file at `internal/export/import_test.go:32-52`.
- VAL-062: `TestImportRequiresFreshDestinationObservationBeforeRecovery` at `internal/export/import_test.go:160-220` proves fresh destination reconciliation, rejects imported old-observation recovery, CloseRun, AbandonRun and historical Gate execution, then binds successor Claim/Run to the destination and executes a real Gate with a successor Receipt.

## Tests run

- Initial compilation red: `go test ./internal/export -run '^TestImportCrashBeforePublicationLeavesNoPartialProject$' -count=1` - FAIL - `undefined: Import` (terminated `exit_code: 1`). This established the missing public domain API before implementation.
- Behavioral red: `go test ./internal/export -run '^TestImportCrashBeforePublicationLeavesNoPartialProject$' -count=1` - FAIL - `Import() error = import is not implemented` (terminated `exit_code: 1`). The staging/publish implementation repaired it.
- Staging red: the same focused import test failed while extracting staged artifacts from a closed archive handle (terminated `exit_code: 1`). Reopening the validated input for staging repaired the resource-lifetime fault.
- Semantic archive reds: `go test ./internal/export -run '^TestArchiveRejectsIntegrityValidSemanticCorruptions' -count=1` - FAIL - `observation_subject_differs_from_its_row` and `issue_aggregate_differs_from_immutable_Effect` accepted (terminated `exit_code: 1`). Archive reference validation now compares the nested observation subject and aggregate Effects.
- Destination-route green: `go test ./internal/export -run '^TestImportRequiresFreshDestinationObservationBeforeRecovery$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/export\t0.902s` (terminated `exit_code: 0`). It runs `evidence.RunGate` through the production process executor and asserts the Receipt; no CLI import adapter exists in this task by explicit Task 15 allocation.
- Focused regression green: `go test ./internal/export ./internal/workflow ./internal/evidence ./internal/journal -run '^(TestImport|TestArchiveRejectsIntegrityValidSemanticCorruptions)' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/journal\t0.005s [no tests to run]` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t91.539s\tcoverage: 52.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t(statements)\t\t\t63.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t56.257s` (terminated `exit_code: 0`)

## Deviations from spec

- The initial imported Event-byte test treated archive journal filenames as preserved physical delivery names. The archive contract preserves ordered Event bytes, not those filenames; the assertion was repaired to compare ordered bytes while retaining its identity/history strength.
- Archive semantic identity validation repairs Task 13 provenance findings for nested observation JSON and aggregate-to-Effect references. DEC-204 explicitly allocated the imported observation fencing in `internal/workflow/run_close.go`; Task 8 remains its original provenance.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes. The implementation uses the existing indirect `golang.org/x/sys` package for Linux `renameat2` no-replace publication and advisory file locking; the standard library has no equivalent no-replace rename primitive. Its BSD-3-Clause license is already present in the module graph.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
