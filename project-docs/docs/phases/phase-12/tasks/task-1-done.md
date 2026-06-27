## Summary of what was implemented

- Added the new `internal/requirements` domain package for Phase 12 Task 1:
  markdown-table parsing, deterministic summary hashing, requirement upsert
  planning, and the journaled `Sync` write path.
- Added `fid requirement sync` as the minimal CLI surface for v0.1
  requirement-pointer ingestion, with `--source-doc`, `--type`, `--id`, and
  `--by` wired through the existing error-rendering path.
- Appended the `requirement.sync` in-code policy row with the section-13
  booleans (`tl` + `human` only) and verified `fid policy show` parity.

## Files created or modified

- `cmd/fid/root.go`
- `cmd/fid/requirement.go`
- `cmd/fid/requirement_test.go`
- `internal/policy/matrix.go`
- `internal/policy/matrix_test.go`
- `internal/requirements/errors.go`
- `internal/requirements/parser.go`
- `internal/requirements/sync.go`
- `internal/requirements/sync_test.go`
- `internal/requirements/types.go`
- `project-docs/docs/phases/phase-12/decisions.md`
- `project-docs/docs/phases/phase-12/tasks/task-1-done.md`

## Documentation updated

- Added `project-docs/docs/phases/phase-12/decisions.md` with DEC-029, which
  records the PI-1 parse contract and the exact `last_synced_hash` input.
- Added this completion report at
  `project-docs/docs/phases/phase-12/tasks/task-1-done.md`.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: verified no change required; the
  implemented flags and sync-count output match the existing `fid requirement
  sync` spec text.

## Validation Contract Coverage

- `VAL-001` (partial): PASS. `internal/policy/matrix.go` now contains exactly
  one `requirement.sync` row with `tl=true`, `human=true`, all other roles
  false, and `fid policy show` exposes the same row (`TestMatrixRequirementSyncRowMatchesSchemaSection13`,
  `TestRequirementSyncAuthorizationAndPolicyShowParity/policy show includes requirement sync row`).
- `VAL-002`: PASS. `fid requirement sync` persists `id`, `type`, `summary`,
  `source_doc`, `source_section`, `last_synced_hash`, and `status='active'`;
  initial sync creates rows in one event, identical re-sync appends no event,
  summary changes update the stored hash, and `--type` / `--id` restrict the
  synced rows (`TestRequirementSyncCreatesPointersInOneJournalEvent`,
  `TestSyncCreatesThenSkipsUnchangedThenUpdatesChangedRows`,
  `TestSyncSupportsTypeAndIDFilters`).
- `VAL-003`: PASS. No-parseable-ID, duplicate-ID, unknown-`--type`, and
  missing-`--id` cases all reject before append with zero journal-line delta
  and zero persisted requirement rows (`TestRequirementSyncRejectsMalformedInputBeforeAppend`,
  `TestSyncRejectsMalformedInputAndDeniedRolesBeforeAppend`).
- `VAL-004`: PASS. `tl` and `human` succeed; denied and unknown roles reject
  before append while naming the role/action in stderr (`TestRequirementSyncCreatesPointersInOneJournalEvent`,
  `TestRequirementSyncAuthorizationAndPolicyShowParity`,
  `TestSyncRejectsMalformedInputAndDeniedRolesBeforeAppend`).

## Tests run

- `go test ./cmd/fid -run TestRequirementSyncCreatesPointersInOneJournalEvent`
  — red: failed with `unknown command "requirement" for "fid"`.
- `go test ./cmd/fid -run TestRequirementSyncCreatesPointersInOneJournalEvent`
  — green.
- `go test ./internal/requirements`
- `go test ./cmd/fid -run 'TestRequirementSync|TestPolicyShow'`
- `go test ./internal/policy`

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS
- `golangci-lint run ./...` — PASS (`0 issues.`)
- `go test -race -coverprofile=coverage.out ./...` — PASS
- `go tool cover -func=coverage.out` — PASS
  Coverage excerpts: total 76.5%; `cmd/fid` 78.1%; `internal/requirements`
  80.4%; `internal/issues` 81.8%; `internal/claims` 81.4%;
  `internal/policy` 81.2%; `internal/journal` 80.0%;
  `internal/workflow` 95.5%.
- `go build -o /tmp/fid ./cmd/fid` — PASS
- `go test -tags integration -race ./...` — N/A (no integration fixtures were
  touched in Task 1)

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `e236591de43d994b2b816e4398c1470b4fc4d4c9` (`e236591`)

## Commit subject

- `feat(requirements): add requirement sync`
