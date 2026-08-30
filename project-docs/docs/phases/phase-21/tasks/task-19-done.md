# Phase 21 - Task 19: Deterministic Export Create And Validate - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: All task-19 implementation, CLI, and exact gate requirements passed.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added deterministic `fid export create` archives with fixed ZIP metadata, canonical manifest data, journal prefix, canonical source metadata, capability snapshot, and available Artifact blobs.
- Added non-mutating `fid export validate` integrity checks for ZIP safety, inventory, manifest canonicality, digests and sizes, journal sequencing/schema/terminal hash, Artifact references, and capability snapshot compatibility.
- Registered only `export create` and `export validate`; `export import` remains an unknown command and `export.portable.v1` remains exactly unavailable with `import:false`.

## Files created or modified

- `internal/export/export.go`
- `internal/export/export_test.go`
- `cmd/fid/export.go`
- `cmd/fid/export_test.go`
- `cmd/fid/root.go`
- `cmd/fid/capabilities_test.go`
- `cmd/fid/contract_surface_test.go`
- `internal/cli/codes.go`
- `internal/cli/codes_test.go`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-19-done.md`

## Validation Contract Coverage

- VAL-069: PASS - `TestCreateIsByteIdenticalWithCompleteDeterministicInventory` at `internal/export/export_test.go:20` creates two archives from the same persisted journal/source/Artifact facts, byte-compares them, validates one, and asserts its complete entry inventory and source metadata.
- VAL-070: PASS - `TestValidateRejectsJournalAggregateSequenceCollisionWithoutMutation` at `internal/export/export_test.go:70` constructs an otherwise manifest-consistent aggregate-sequence collision, receives `journal-invalid`, and byte-compares project state before and after validation. `Validate` in `internal/export/export.go` emits the complete closed check set for path safety, duplicate/extra/missing inventory, manifest, digest, size, journal, terminal hash, Artifact-reference, and capability-snapshot defects.
- VAL-071: PASS - `TestExportCreateValidateAndImportRemainsAbsent` at `cmd/fid/export_test.go:9` proves create/validate JSON success, `export import` unknown-command rejection, and the exact false capability row with `import:false`; the required source absence command also passed.

## Tests run

- `go test ./internal/export` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	0.114s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run TestExportCreateValidateAndImportRemainsAbsent` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.060s` (terminated `exit_code: 0`)
- `go test ./cmd/fid ./internal/export ./internal/cli` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	(cached)` (terminated `exit_code: 0`)
- `/tmp/fid --project /tmp/tmp.f9uDkl9Ton export create /tmp/tmp.f9uDkl9Ton/recovery.zip --json` - PASS - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"export.create","result":{"export":{"archive_path":"/tmp/tmp.f9uDkl9Ton/recovery.zip","archive_digest":"sha256:7d0b6337ed251c49effd7dab67c60a436a2a1a3bbb080c588c9fc59bfaa78ee4","manifest_schema":"fid.export-manifest.v1","project_id":"project-918ca2004c738264387b828de5c3ab1f","event_count":1,"terminal_event_hash":"sha256:f96fe3cd377f4ad8d243ca05475bcd87be6666426416a59a3632fc2bf2c37b0b","projection_schema":14,"entry_count":2}}}` (terminated `exit_code: 0`)
- `/tmp/fid export validate /tmp/tmp.f9uDkl9Ton/recovery.zip --json` - PASS - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"export.validate","result":{"validation":{"archive_path":"/tmp/tmp.f9uDkl9Ton/recovery.zip","archive_digest":"sha256:7d0b6337ed251c49effd7dab67c60a436a2a1a3bbb080c588c9fc59bfaa78ee4","status":"valid","manifest_schema":"fid.export-manifest.v1","entry_count":2}}}` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race ./internal/export ./internal/evidence ./internal/journal ./internal/cli ./cmd/fid/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	333.981s` (terminated `exit_code: 0`)
- `rg -n --glob '*.go' --glob '!**/*_test.go' 'Use:\s*"import|newExportImport|export\.import' cmd/fid internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 29.2% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					62.1%` (terminated `exit_code: 0`)

## Deviations from spec

- Initial `golangci-lint run ./...` attempt failed with `1 issues:` (terminated `exit_code: 1`) because `validateJournalBytes` was an unused wrapper. Removed that in-scope wrapper in `2518931`; the byte-identical later lint command passed with `0 issues.` (terminated `exit_code: 0`). The six-command gate was restarted from command 1 after the repair.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
