# Phase 22 - Task 13: Recovery archive format and validation - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Manifest v2 carries closed identity-bound project configuration and ordinal journal chunks with unchanged canonical Event bytes. Available Artifact bodies stream from content-addressed storage; unavailable Artifact facts remain in the journal without invented bytes. No source prose, SQLite image, SQL or import delivery metadata is exported.
- Validation bounds actual ZIP entry count before allocation, actual per-member and cumulative decompression, and fragmented Event lines. It rejects unsafe/special/duplicate paths before expansion and stops at the first byte overflow.
- V2 validates canonical identity/sequence/hash and closed Event/ordered Effect families, table/key/column/value shape, immutable operational update fields, Event references, Recovery ownership and polymorphic subjects. Trusted installed migrations build a private in-memory projection; the normal Event applier enforces SQLite references and captured preconditions. The archive supplies no SQL.
- Manifest v1 keeps its pinned schema14/capability profile and accepted read-only integrity semantics. New full semantic replay belongs to v2; no failed v2 archive falls back to v1. The implementation is split by writing, reading, scanning, profiles and semantic responsibilities.
- This is the root delegate's completion record after serial takeover, with earlier assistant contributions explicitly attributed below.

## Files created or modified

- `cmd/fid/export.go`
- `cmd/fid/export_test.go`
- `internal/config/project_config.go`
- `internal/export/archive_availability_test.go`
- `internal/export/archive_config.go`
- `internal/export/archive_consistency.go`
- `internal/export/archive_corruption_test.go`
- `internal/export/archive_directory.go`
- `internal/export/archive_fixture_test.go`
- `internal/export/archive_journal.go`
- `internal/export/archive_limits_test.go`
- `internal/export/archive_paths_test.go`
- `internal/export/archive_profiles.go`
- `internal/export/archive_references.go`
- `internal/export/archive_scan.go`
- `internal/export/archive_semantics.go`
- `internal/export/archive_semantics_catalog.go`
- `internal/export/archive_semantics_json.go`
- `internal/export/archive_semantics_test.go`
- `internal/export/archive_semantics_values.go`
- `internal/export/archive_updates.go`
- `internal/export/archive_v2_test.go`
- `internal/export/archive_validation.go`
- `internal/export/archive_writer.go`
- `internal/export/archive_writer_test.go`
- `internal/export/export.go`
- `internal/export/export_boundary_test.go`
- `internal/export/export_test.go`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-13-done.md`
- `project-docs/docs/phases/phase-22/tasks/task-13-root-evidence.md` retains original/red/construction and gate-attempt history.
- Runner-owned `delegate-handoff.md` and `run-ledger.md` record the renewed serial execution authorization; they are administrative artifacts, not product validation.

## Validation Contract Coverage

All five final-source gates passed with observed termination. Assertion evidence:

- VAL-054: PASS - `TestCreateRecoveryManifestV2ConfigAndOrdinalJournal`, `TestCreateIsByteIdenticalWithCompleteDeterministicInventory`, `TestCreateAndValidatePreserveUnavailableReceiptArtifactWithoutBody`, `TestJournalWriterChunksAtLineBoundaryPreservingBytes`: configuration/identity, original Event lines, complete Source/Artifact inventory, both availability states, deterministic repeated export, and exact chunk boundary.
- VAL-055: PASS - `TestArchiveWithValidHashesButInvalidRecoveryReferenceRejects` first validates a real recovered-history control then changes only its observation reference with recomputed hashes. `TestArchiveRejectsIntegrityValidSemanticCorruptions` adds unknown family/table/column, wrong scalar/boolean types, unresolved idempotency Event, foreign Recovery aggregate, duplicate immutable Artifact, retired caller outcome and executor rewrite. `TestSemanticValidatorRejectsNestedDuplicateJSONKey` and `TestSemanticValidatorRejectsRecoveryEffectFamilyTampering` cover duplicate nested keys, wrong operations and ordering.
- VAL-056: PASS - 056-a `TestArchiveEntryCountBoundaries` creates 999999/1000000/1000001 actual ZIP records and runs both pre-allocation record counting and directory inspection. 056-b/c `TestArchiveExpandedByteBoundaries` streams valid-CRC ZIPs at 64MiB and 1GiB minus/equal/plus one. 056-d `TestArchiveCumulativeBoundary` streams 8GiB minus/equal/plus one through individually legal members and checks measured counts. 056-e `TestReadEventLinesUsesPerLineLimitAcrossFragments` covers 1MiB minus/equal/plus one, fragmented input and multiple legal lines. 056-f..k `TestArchiveActualZIPPathTypeAndDuplicateMatrix` names real safe, traversal, absolute/drive, symlink, device, directory and duplicate cases. `TestArchiveStopsAtFirstExpandedOverflow` proves exactly L+1 consumed and no later member opened. Large fixtures are generated in fresh disk-backed .cache children and cleaned after each test; no limit was reduced or test skipped.
- VAL-057: PASS - `TestValidateHasNoSourceOrDestinationEffects` snapshots both initialized source and uninitialized destination with operator sentinel for valid/invalid archives. Original negative-check and special-file controls also snapshot source bytes. V1 has a retained read-only control in `TestValidateRetainsPinnedV1Profile`.

## Tests run

- `go test ./internal/export -run '^TestArchiveWithValidHashesButInvalidRecoveryReferenceRejects$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	0.795s` (terminated `exit_code: 0`)
- `go test ./internal/export -run '^(TestArchiveEntryCountBoundaries|TestArchiveStopsAtFirstExpandedOverflow)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	10.135s` (terminated `exit_code: 0`)
- `go test ./internal/export -run '^(TestArchiveRejectsIntegrityValidSemanticCorruptions|TestArchiveWithValidHashesButInvalidRecoveryReferenceRejects|TestSemantic|TestCreate|TestValidate)' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	9.340s` (terminated `exit_code: 0`)
- `go test ./internal/export -run '^TestArchiveActualZIPPathTypeAndDuplicateMatrix$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	0.006s` (terminated `exit_code: 0`)

- `go test -race ./internal/export -run '^(TestArchiveExpandedByteBoundaries|TestArchiveCumulativeBoundary|TestArchiveStopsAtFirstExpandedOverflow)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	168.333s` (terminated `exit_code: 0`)

These are root-observed focused receipts. Original assistant-run large-boundary
and semantic red receipts are attributed in the contribution history, not
presented as fresh root commands. The final full gates execute every test.

The built `/tmp/fid` CLI initialized fresh scratch, created a v2 archive, and validated it from `/tmp` outside a workspace. All three commands terminated `exit_code: 0`; validation returned `status:valid` and the matching archive digest. Exact commands/results are retained in `task-13-root-evidence.md`.

## Exact quality-gate results

Working directory: `/home/mark/projects/fid`.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)

- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 52.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					63.2%` (terminated `exit_code: 0`)

- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	56.215s` (terminated `exit_code: 0`)

## Deviations from spec

- Tagged integration attempt1 failed with `FAIL` (terminated `exit_code: 1`), package55.349s. The retained full CLI journey reached a `run.execution-blocked` idempotency row missing from the new semantic catalog. Added the existing writer's required final stage. Its error also exposed the Phase21 adapter's string-pointer/null conversion bug; new public CLI negative tests captured both missing bodies before repair. Focused journey and error rendering passed afterward. The temporary direct-Validate diagnostic was removed, leaving the original journey assertion intact. The final-source full race/coverage rerun passed, and the byte-identical final tagged rerun passed `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	56.215s` (terminated `exit_code: 0`).
- No product contract, migration, Event schema, gate definition, VAL threshold or dependency changed.
- Initial test-quality/TDD deviation: the originally named first test only corrupted Source metadata. It is preserved under an accurate Source name. The semantic assistant then captured the genuine public-Validate acceptance red before its semantic replay implementation, after separate scanner/writer edits had begun. Root subsequently strengthened the final test to full real Recovery history. No earlier result is relabeled as stronger proof.
- Required race/coverage attempt1 failed with `FAIL` (terminated `exit_code: 1`) and `panic: test timed out after 10m0s`; the stack identified repeated zero-fixture compression before the first cumulative-limit assertion. Root changed fixture generation to reuse stdlib-generated byte-aligned non-final DEFLATE blocks, with real full-size expansion and CRC checking unchanged. Every original boundary remains. The byte-identical second full gate passed after the fixture repair: `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 52.9% of statements` (terminated `exit_code: 0`). Export passed at367.930s and the CLI at461.712s; no timeout was raised.
- Required lint attempt1 failed with `* staticcheck: 4` (terminated `exit_code: 1`); attempt2 with `* staticcheck: 3` (terminated `exit_code: 1`). Interrupted code had uppercase error strings and a simplifiable boolean expression. All strings were corrected together after the second result and the byte-identical third command passed `0 issues.` (terminated `exit_code: 0`). Earlier failures remain in root-evidence.
- Additional semantic/path reds and construction failures are retained in `task-13-root-evidence.md`, including the initial valid-file controls, malformed test field names, unsupported fixture release reason and interrupted Go assignment. The causes were repaired without weakening assertions. Initial parallel ownership created avoidable construction failures; the renewed Human instruction changed execution to a single root writer.
- Root recorded the administrative RUN-RESUME after initial reconciliation/repair rather than before them; the ledger explicitly preserves that ordering deviation.

## Concerns, technical debt, or notes for the Tech Lead

- Retained convergence provenance: `35a75baf7b07d8f5eaace36ac8142ea6ac7c4ddd`, `feat(export): create deterministic archive validation`, Phase21/task19, introduced the direct-SQL Source/Artifact fixture and dated-file inventory helper. V2's journal-backed semantic facts and ordinal chunks invalidate those constructions. The complete inventory and every existing rejection remain covered by the full ordinary gate; CLI/export retained journeys execute under the tagged gate.
- Compatibility interpretation is independently supported by Phase21 contracts section11: v1 validation explicitly does not establish identity/history parity. Phase22 task13/public-interface section5 retain it; no new decision or weakened v2 check was needed.
- Task14/15 still own atomic publication, imported destination qualification, and every-family import/re-export parity. Task16 owns schema16 adoption and Task17 the final capability values. Task18/finalQA still own the full corpus and retained Human delivery.
- Fixture contributions committed at 1be8e88 and c9b4665 are earlier same-task work. Their author may not independently approve those contributions in Review5. Root's final exact gates cover the integrated source.

## Dependency-budget note

- No dependency changes. Uses the standard library and existing SQLite, migration and journal packages.

## Commit hash

- `acb569bb9447a3259991520669167fdc819c1706`

## Commit subject

- `feat(export): validate bounded v2 recovery archives`
