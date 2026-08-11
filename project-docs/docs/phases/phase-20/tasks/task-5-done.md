# Phase 20 - Task 5: Remove orphaned Issue fields - Completion Report

## Summary of what was implemented

- Removed `in_inbox`, `est_tokens`, and `due_date` from the Issue projection, domain and CLI surfaces, newly emitted Issue effects, format fixtures, parity fixtures, and replay allowlists.
- Added migration `0011_p20_remove_issue_fields`; its table rebuild removes the retired columns and dependent CHECK, recreates `issues_enriched`, and its down migration restores all three columns plus the former epic `est_tokens` CHECK.
- Retired Issue effect keys now stop development-journal replay with `journal-replay-stopped` instead of being ignored.
- Preserved the inherited regenerated parity fixtures and their manifest. Re-verification in this session found every `SHA256SUMS` entry valid, including the three changed Issue-get fixtures.
- Applied DEC-137's scoped verification-apparatus repair: `TestIntegrationClassification` now verifies the four Task 3f journeys are present and reports the individual missing journey; additional integration journeys are permitted. The repository-wide tagged-test walk and the binary/sleep classification tests are unchanged.

## Files created or modified

- Added `migrations/0011_p20_remove_issue_fields.up.sql`, `migrations/0011_p20_remove_issue_fields.down.sql`, and `internal/db/migrate_test.go`.
- Updated the scoped Issue model, create/update/get/list/ready paths, journal replay, format renderers, and their tests under `internal/` and `cmd/fid/` to remove the three retired fields.
- Updated `cmd/fid/integration/v01_parity_test.go` and `cmd/fid/integration/classification_test.go`.
- Updated committed regression fixtures: `cmd/fid/testdata/v01_parity/issue-get-epic.json`, `issue-get-task.json`, `issue-get-subtask.json`, `SHA256SUMS`, and `testdata/format/issue-empty.golden.toon`.
- Preserved but did not modify or stage the pre-existing out-of-scope `project-docs/docs/phases/phase-20/run-ledger.md` change.

## Documentation updated

- Updated this completion report. No product or architecture documentation change was required by the task spec.

## Validation Contract Coverage

- `VAL-009` — PASS. `TestIssuesTableExcludesRetiredPhase20Fields` in `internal/db/migrate_test.go:10` proves the up migration omits all three columns and the down migration restores the columns and epic CHECK. `TestIssueCommandsDoNotRegisterRetiredFieldFlags` in `cmd/fid/issue_section_test.go:212`, `TestIssueGetSectionsAndRetiredFieldAbsence` in `cmd/fid/issue_section_test.go:13`, and `TestCreateEventOmitsRetiredPhase20Fields` in `internal/issues/issue_test.go:851` cover CLI, decoded JSON, and newly emitted effects. `TestReplayJournalClassifiesRetiredIssueEffectKeyAsStopped` in `internal/journal/replay_test.go:31` covers the non-ignored removed-key state. The current-session race suite passed with terminal token `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 80.5% of statements`.
- `VAL-012` (field-removal half) — PASS. `TestV01ParityFixtureRebuildPreservesEveryProjectionTable` in `cmd/fid/integration/v01_parity_test.go:47` compares each live projection table to its rebuilt counterpart; the required dedicated integration command passed in this session with terminal token `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t4.716s`. Fixture-manifest verification passed: `issue-get-epic.json`, `issue-get-task.json`, and `issue-get-subtask.json` are respectively `a4c4009e5394cd1f3ab717c52419422aed2e3c7302247ccefa27331cbdd05d6d`, `ff0063400ace3435be2d95460c40bd55305bc4c4a473b5935f4636f5377b52ec`, and `0069ccb123b3e941455dad9c93263072660c20acd6a791a6f3f5544d27fff78c`; `sha256sum -c SHA256SUMS` reported every parity fixture `OK`.

## Tests run

- `sha256sum -c SHA256SUMS` from `cmd/fid/testdata/v01_parity` — PASS — `issue-list-subtask.json: OK` (terminated `exit_code: 0`); the recorded three regenerated Issue-get hashes match the inherited manifest and fixture bytes.
- RED, `go test -tags integration ./cmd/fid/integration -run '^TestIntegrationClassification$' -count=1` — expected failure after temporarily excluding the required observed journey `TestV01ParityGoldenOutputs` — `FAIL` (terminated `exit_code: 1`); diagnostic: `required dedicated integration journey "TestV01ParityGoldenOutputs" is missing; restore it to cmd/fid/integration`.
- GREEN, same focused command after restoring that observed journey — PASS — `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t0.430s` (terminated `exit_code: 0`).
- Real entrypoint: created a fresh temporary workspace from `cmd/fid/testdata/v01_parity/fixture`, then ran `/tmp/fid issue get FID-b20002 --format json` — PASS — output contained the expected Issue and omitted `in_inbox`, `est_tokens`, and `due_date` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS — `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` — PASS — `0 issues.` (terminated `exit_code: 0`).
- `go test -race -coverprofile=coverage.out ./...` — PASS — `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 80.5% of statements` (terminated `exit_code: 0`).
- `go tool cover -func=coverage.out` — PASS — `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t77.7%` (terminated `exit_code: 0`).
- `go test -tags integration -race ./cmd/fid/integration/...` — PASS — `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t4.716s` (terminated `exit_code: 0`).
- `go test -bench . -benchmem ./...` — PASS — `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.573s` (terminated `exit_code: 0` after explicit polling of session `56859`).
- `go build -o /tmp/fid ./cmd/fid` — PASS — `exit 0 (no output)` (terminated `exit_code: 0`).

## Deviations from spec

- DEC-137/DEC-138 authorized verification-apparatus repair: the first dedicated-gate failure recorded in the inherited blocker report had verbatim terminal token `FAIL` and named the exact stale inventory discrepancy: `TestV01ParityFixtureRebuildPreservesEveryProjectionTable` was present in the dedicated package but absent from `wantJourneys`. Root cause: exact equality treated Task 5's legitimate new rebuild-parity journey as a failure, although the accepted invariant is preservation of the four Task 3f journeys. Fix: replace equality with a required-presence loop that names any missing journey and permits additions; `taggedTestFilesOutside`, its repository-wide walk, `TestIntegrationBinaryOwnedByPackage`, and `TestClosureJourneyHasNoWallClockSleep` were left unchanged. No product assertion, threshold, budget, tolerance, or command changed. The byte-identical required integration command was rerun in this session and passed with its own token `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t4.716s`.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `028855b`

## Commit subject

- `feat(issue): remove retired issue fields`
