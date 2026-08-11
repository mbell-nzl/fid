# Phase 20 — Task 5: Remove orphaned Issue fields — Blocker Report

## Summary of what was implemented

- Prepared the Task 5 implementation in the working tree: migration `0011_p20_remove_issue_fields` removes `in_inbox`, `est_tokens`, and `due_date` from the Issue projection; the down migration restores the prior logical columns and the epic `est_tokens` CHECK behavior.
- Removed the three fields from the scoped Issue domain, CLI flags and JSON, emitted Issue event effects, format fixtures, parity fixtures, and replay path. Replay now classifies a retired Issue effect key as a stopped journal replay instead of ignoring it.
- BLOCKED before completion: the required integration gate failed because `TestIntegrationClassification` has a fixed inventory that does not include the new `TestV01ParityFixtureRebuildPreservesEveryProjectionTable` journey. Per the task instructions, no completion commit was created and the remaining required gates were not run.

## Files created or modified

- Added `migrations/0011_p20_remove_issue_fields.up.sql` and `migrations/0011_p20_remove_issue_fields.down.sql`.
- Added `internal/db/migrate_test.go` and updated scoped Issue, journal replay, format, CLI, parity-fixture, and test files under `internal/` and `cmd/fid/`.
- Regenerated prepared fixture files (all remain uncommitted because the required gate blocked task completion):
  - `cmd/fid/testdata/v01_parity/issue-get-epic.json`
  - `cmd/fid/testdata/v01_parity/issue-get-subtask.json`
  - `cmd/fid/testdata/v01_parity/issue-get-task.json`
  - `cmd/fid/testdata/v01_parity/SHA256SUMS`
  - `testdata/format/issue-empty.golden.toon`
- Preserved the pre-existing unrelated modification to `project-docs/docs/phases/phase-20/run-ledger.md` without editing it.

## Documentation updated

- This blocker report documents the incomplete task state and fixture-regeneration procedure. No product or architecture documentation change was required by the task spec.

## Validation Contract Coverage

- `VAL-009` — BLOCKED at the required integration gate. Focused evidence passed before the gate: `TestIssuesTableExcludesRetiredPhase20Fields` in `internal/db/migrate_test.go:10` verifies the up migration removes all three Issue projection columns and the down migration restores them plus the epic `est_tokens` CHECK; `TestIssueCommandsDoNotRegisterRetiredFieldFlags` in `cmd/fid/issue_section_test.go` verifies the removed CLI flags; `TestIssueGetSectionsAndRetiredFieldAbsence` in the same file verifies JSON omits all three fields.
- `VAL-012` field-removal half — BLOCKED at the required integration gate. Focused evidence passed before the gate: `TestCreateEventOmitsRetiredPhase20Fields` in `internal/issues/issue_test.go` verifies new issue-created effects omit all three keys; `TestReplayJournalClassifiesRetiredIssueEffectKeyAsStopped` in `internal/journal/replay_test.go:31` verifies a development Event with a removed key becomes a stopped replay; `TestV01ParityFixtureRebuildPreservesEveryProjectionTable` in `cmd/fid/integration/v01_parity_test.go:47` compares every rebuilt projection table with the live projection.

## Tests run

- RED: `go test ./internal/db -run '^TestIssuesTableExcludesRetiredPhase20Fields$' -count=1` — expected failure before migration `0011`: `FAIL`.
- GREEN: `go test ./internal/db -run '^TestIssuesTableExcludesRetiredPhase20Fields$' -count=1` — `ok  github.com/mbell-nzl/fid/internal/db 0.026s`.
- RED: `go test ./internal/journal -run '^TestReplayJournalClassifiesRetiredIssueEffectKeyAsStopped$' -count=1` — expected failure before replay classification: `FAIL`.
- GREEN: `go test ./internal/journal -run '^TestReplayJournalClassifiesRetiredIssueEffectKeyAsStopped$' -count=1` — `ok github.com/mbell-nzl/fid/internal/journal 0.023s`.
- `go test ./internal/issues ./cmd/fid -run '^(TestCreateEventOmitsRetiredPhase20Fields|TestIssueCommandsDoNotRegisterRetiredFieldFlags|TestIssueGetSectionsAndRetiredFieldAbsence)$' -count=1` — `ok github.com/mbell-nzl/fid/cmd/fid 0.125s`.
- `go test ./internal/db ./internal/issues ./internal/journal ./cmd/fid ./internal/format -count=1 && go test -tags integration ./cmd/fid/integration -run '^(TestV01ParityGoldenOutputs|TestV01ParityFixtureRebuildPreservesEveryProjectionTable)$' -count=1` — `ok github.com/mbell-nzl/fid/cmd/fid/integration 0.536s`.
- Fixture regeneration procedure: `go build -o /tmp/fid-task5-fixture ./cmd/fid && fixture_out=$(mktemp -d) && cmd/fid/testdata/v01_parity/capture.sh /tmp/fid-task5-fixture "$fixture_out" && git diff --no-index -- cmd/fid/testdata/v01_parity "$fixture_out"`. The capture completed; the final `git diff --no-index` returned the expected non-zero result because it displayed the three regenerated `issue-get-*.json` field removals and `SHA256SUMS` changes. The prepared fixture hashes are `a4c4009e5394cd1f3ab717c52419422aed2e3c7302247ccefa27331cbdd05d6d` (epic), `0069ccb123b3e941455dad9c93263072660c20acd6a791a6f3f5544d27fff78c` (subtask), and `ff0063400ace3435be2d95460c40bd55305bc4c4a473b5935f4636f5377b52ec` (task).
- Real entrypoint: `go build -o /tmp/fid-task5-runtime ./cmd/fid && task5_workspace=$(mktemp -d) && cp -R cmd/fid/testdata/v01_parity/fixture/. "$task5_workspace" && cd "$task5_workspace" && /tmp/fid-task5-runtime issue get FID-b20002 --format json` — emitted Issue JSON omitted `in_inbox`, `est_tokens`, and `due_date`.

## Exact quality-gate results

- `PASS` — `gofmt -l cmd internal` — `exit 0 (no output)`.
- `PASS` — `golangci-lint run ./...` — `0 issues.`
- `PASS` — `go test -race -coverprofile=coverage.out ./...` — `ok  	github.com/mbell-nzl/fid/internal/workflow	9.525s	coverage: 80.5% of statements`.
- `PASS` — `go tool cover -func=coverage.out` — `total: (statements) 77.7%`.
- `BLOCKED` — `go test -tags integration -race ./cmd/fid/integration/...` — `FAIL`. `TestIntegrationClassification` reported its journey inventory as `[TestMaintenanceCheckReportsIncompatibleWorkflowAssignment TestV01ClosureJourneyAndRebuildParity TestV01ParityFixtureRebuildPreservesEveryProjectionTable TestV01ParityGoldenOutputs TestV02aWorkflowAuthoringExitScenarioAndRebuildParity]` while its expected list omitted `TestV01ParityFixtureRebuildPreservesEveryProjectionTable`.
- `NOT RUN` — `go test -bench . -benchmem ./...` — stopped after the preceding required gate failed.
- `NOT RUN` — `go build -o /tmp/fid ./cmd/fid` — stopped after the preceding required gate failed.

## Deviations from spec

- None. The task was not completed because the required quality gate failed; no exception or narrower substitute was used.

## Concerns, technical debt, or notes for the Tech Lead

- The integration-classification inventory must be reconciled with the newly added parity-fixture rebuild journey before Task 5 can complete. This report leaves that repair uncommitted because the required gate failure requires a blocker stop.

## Dependency-budget note

- No dependency changes.

## Commit hash

- No commit created — blocker report.

## Commit subject

- No commit created — blocker report.
