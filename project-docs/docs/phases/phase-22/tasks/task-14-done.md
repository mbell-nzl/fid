# Phase 22 Task 14 completion report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Recovery attempts: 1
- Human attention required: no
- Condition: None. Root repaired the retained checkpoint defects; all final exact gates passed.
- Next action: None.

## Summary of what was implemented

Root completed the same-task continuation at `8c322b53c749a2a7f0dd0be3e03dbb5034d70e40`. The failed contender/history tests, repairs, actual process-exit matrix and final exact gate receipts are in `task-14-root-evidence.md`. Developer history below retains its original attribution.

- Added the domain import API. It snapshots and validates the supplied archive, stages configuration, journal and artifacts, rebuilds the SQLite projection with normal appliers, fsyncs the staged tree, and atomically publishes `.fid` with no replacement.
- Added durable delivery metadata and idempotent replay. A verified installed binding returns the original result only for the same archive digest and idempotency key; populated, unsafe, changed-key and malformed installed targets are rejected.
- Preserved imported Event bytes/order/identities, extended ordinal journal handling for imported history, and prevented new locally written Event IDs from colliding with imported ordinal chunks.
- Added DEC-202 destination fencing: imported observations and historical execution cannot establish destination authority; a fresh local observation and reconciliation produce destination-bound successor Claim/Run, whose Gate produces a Receipt. DEC-204 also fences CloseRun and AbandonRun against imported observations.
- Added archive semantic validation for nested observation identity and aggregate-to-immutable-Effect identity, covering Task 13 provenance findings.

## Files created or modified

- `internal/export/import_publication.go`, `internal/export/import_boundary_test.go`, `internal/export/archive_v2_test.go`, `internal/export/import.go`, `internal/export/import_test.go`, `internal/export/export_test.go`, `internal/export/archive_references.go`, and `internal/export/archive_corruption_test.go`
- `internal/config/import_binding.go`, `internal/config/import_history.go`
- `internal/journal/replay.go` and `internal/journal/writer.go`
- `internal/workflow/import_binding.go`, `internal/workflow/reconciliation.go`, `internal/workflow/recovery.go`, `internal/workflow/run_close.go`, and `internal/workflow/run_inspect.go`
- `internal/evidence/execution.go` and `internal/claims/generic.go`

## Documentation updated

- Added this completion report and root continuation evidence `task-14-root-evidence.md`. The DEC-202 and DEC-204 allocations were already recorded in the phase task and decision documents by the Phase Runner.

## Validation Contract Coverage

Final root outcomes, backed by the five final exact gates:

- VAL-058: PASS — ordered journal identity comparison, unchanged source files, semantic corruption rejections and verified imported journal bytes/IDs; see `TestImportPublishesIdentityAndReplaysLostResponse`, `TestImportBindingRejectsAlteredHistoryIdentity` and `TestArchiveRejectsIntegrityValidSemanticCorruptions`.
- VAL-059: PASS — `TestImportProcessExitAtPublicationBoundaries` abruptly exits real subprocesses at five boundaries; before publication no target `.fid` exists, after publication a complete verified installation replays. `TestImportLockExcludesRepeatedContenders` proves exclusion survives repeated contenders.
- VAL-060: PASS — same digest/key replay leaves facts unchanged, including after successor Events; changed digest/key, corrupt result and history reject. `TestImportRejectsChangedDigestAndInstalledHistory` and the destination journey cover both sides.
- VAL-061: PASS — populated and symlinked targets reject unchanged, Git/source files and even a preexisting similarly named lock file survive. Snapshot source replacement occurs before staging; original validated content is installed.
- VAL-062: PASS — old observation key replays without a new qualifying fact; imported observation Close/Abandon/recovery and historical Gate reject; new local liveness/reconciliation binds successor Claim/Run to the actual destination, and a real Gate creates its Receipt. Ordinary foreign-Worktree rejection remains in the required suites. See `TestImportRequiresFreshDestinationObservationBeforeRecovery` and retained recovery tests.

### Retained Developer checkpoint claims

- VAL-058 (Developer checkpoint; root final validation pending): `TestImportPublishesIdentityAndReplaysLostResponse` at `internal/export/import_test.go:38-59` verifies stable Project identity, source preservation, installed binding, ordered Event bytes, and no mutation on replay. `TestArchiveRejectsIntegrityValidSemanticCorruptions` at `internal/export/archive_corruption_test.go:58-77` rejects the valid-hash nested-observation and aggregate identity forgeries.
- VAL-059 (Developer checkpoint; root final validation pending): `TestImportCrashBeforePublicationLeavesNoPartialProject` at `internal/export/export_test.go:74` captures the pre-publication fault boundary; `TestImportAfterPublicationFaultReplaysInstalledResult` at `internal/export/import_test.go:105-128` proves the after-publication/lost-response boundary replays the installed result. `TestImportSnapshotsInputReclaimsStaleLockAndRejectsTamperedMarker` verifies snapshot staging after pathname replacement and that a stale lock path is reusable.
- VAL-060 (Developer checkpoint; root final validation pending): `TestImportPublishesIdentityAndReplaysLostResponse` at `internal/export/import_test.go:53-64` proves same digest/key replay with unchanged installed facts and rejects a different key. `TestImportAfterPublicationFaultReplaysInstalledResult` at `internal/export/import_test.go:116-128` proves retry after a publication response loss.
- VAL-061 (Developer checkpoint; root final validation pending): `TestImportRejectsExistingAndUnsafeTargetsWithoutChangingThem` at `internal/export/import_test.go:68-103` rejects populated and symlink targets and proves the populated marker is unchanged. The identity/replay test preserves an existing target source file at `internal/export/import_test.go:32-52`.
- VAL-062 (Developer checkpoint; root final validation pending): `TestImportRequiresFreshDestinationObservationBeforeRecovery` at `internal/export/import_test.go:160-220` proves fresh destination reconciliation, rejects imported old-observation recovery, CloseRun, AbandonRun and historical Gate execution, then binds successor Claim/Run to the destination and executes a real Gate with a successor Receipt.

## Tests run

The following receipts remain attributed to the original Developer session.
Root continuation red/green and final receipts are separately recorded in
`task-14-root-evidence.md`; no historical receipt is presented as a root rerun.

- Initial compilation red: `go test ./internal/export -run '^TestImportCrashBeforePublicationLeavesNoPartialProject$' -count=1` - FAIL - `undefined: Import` (terminated `exit_code: 1`). This established the missing public domain API before implementation.
- Behavioral red: `go test ./internal/export -run '^TestImportCrashBeforePublicationLeavesNoPartialProject$' -count=1` - FAIL - `Import() error = import is not implemented` (terminated `exit_code: 1`). The staging/publish implementation repaired it.
- Staging red: the same focused import test failed while extracting staged artifacts from a closed archive handle (terminated `exit_code: 1`). Reopening the validated input for staging repaired the resource-lifetime fault.
- Semantic archive reds: `go test ./internal/export -run '^TestArchiveRejectsIntegrityValidSemanticCorruptions' -count=1` - FAIL - `observation_subject_differs_from_its_row` and `issue_aggregate_differs_from_immutable_Effect` accepted (terminated `exit_code: 1`). Archive reference validation now compares the nested observation subject and aggregate Effects.
- Destination-route green: `go test ./internal/export -run '^TestImportRequiresFreshDestinationObservationBeforeRecovery$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/export\t0.902s` (terminated `exit_code: 0`). It runs `evidence.RunGate` through the production process executor and asserts the Receipt; no CLI import adapter exists in this task by explicit Task 15 allocation.
- Import hardening green: `go test ./internal/export -run '^TestImportSnapshotsInputReclaimsStaleLockAndRejectsTamperedMarker$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/export\t0.578s` (terminated `exit_code: 0`). It replaces the source pathname after the immutable snapshot, reuses a stale post-exit lock pathname, and rejects marker delivery metadata whose Project binding does not match the installed config.
- Focused regression green: `go test ./internal/export ./internal/workflow ./internal/evidence ./internal/journal -run '^(TestImport|TestArchiveRejectsIntegrityValidSemanticCorruptions)' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/journal\t0.005s [no tests to run]` (terminated `exit_code: 0`).

## Exact quality-gate results

Root observed every final-source result below terminate in this continuation, working directory `/home/mark/projects/fid`.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	101.200s	coverage: 52.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					63.2%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	52.659s` (terminated `exit_code: 0`)

### Preserved Developer intermediate cycle

These superseded receipts are retained verbatim with child-session attribution, including their original escaped whitespace. They are not the final closure receipts.

> - `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> - `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
> - `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t91.539s\tcoverage: 52.9% of statements` (terminated `exit_code: 0`)
> - `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t(statements)\t\t\t63.3%` (terminated `exit_code: 0`)
> - `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t56.257s` (terminated `exit_code: 0`)

## Deviations from spec

- Root used the Human-authorized direct same-task repair route after checkpoint 0a4d00d. No new product decision or gate change. Convergence provenance, including Task8/13 and the journal writer, is recorded in `task-14-root-evidence.md`.

- The initial imported Event-byte test treated archive journal filenames as preserved physical delivery names. The archive contract preserves ordered Event bytes, not those filenames; the assertion was repaired to compare ordered bytes while retaining its identity/history strength.
- Archive semantic identity validation repairs Task 13 provenance findings for nested observation JSON and aggregate-to-Effect references. DEC-204 explicitly allocated the imported observation fencing in `internal/workflow/run_close.go`; Task 8 remains its original provenance.

## Concerns, technical debt, or notes for the Tech Lead

- The Developer COMPLETE claim was superseded by checkpoint 0a4d00d RECOVERY_REQUIRED. Root repaired its advisory-lock unlink race and missing installed-history verification; all five final exact commands now pass. Independent Review5 remains due after Task15. This does not judge Human outcomes.

### Root takeover — 2026-09-08

Root received clean checkpoint `0a4d00de3581c05f5f1c1b9748da87916e28bc1c`
with no live commands. The Developer's later race/coverage command terminated
exit_code 0 in its session 33222; this is child-attributed intermediate
evidence, not root's final-source gate. The retained historical tokens above
contain literal escaped whitespace and are not silently reconstructed as
fresh root receipts. Root has recorded its own final exact commands above.
The report's recovery fields were normalized to the existing report schema
before control execution. No product gate or acceptance definition changed.

## Dependency-budget note

- No dependency changes. The implementation uses the existing indirect `golang.org/x/sys` package for Linux `renameat2` no-replace publication and advisory file locking; the standard library has no equivalent no-replace rename primitive. Its BSD-3-Clause license is already present in the module graph.

## Commit hash

- `8c322b53c749a2a7f0dd0be3e03dbb5034d70e40`

## Commit subject

- `fix(export): verify imported history and preserve exclusive publication`
