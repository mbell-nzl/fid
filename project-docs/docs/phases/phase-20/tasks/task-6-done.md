# Phase 20 - Task 6: Replayable Issue aggregate revision - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: Fresh final QA re-review.
- Human attention required: no

## Summary of what was implemented

- Added migration `0012_p20_issue_revision`, which adds the non-null Issue `revision` column with default `1` and removes it on down migration.
- Added the reusable `issues.BuildIssueRevisionEffects` helper. It derives absolute next revisions from observed locked state, rejects invalid/duplicate inputs, and sorts multi-Issue Effects by Issue ID.
- Applied revision Effects to Issue creation, update, transition, and legacy `dev start`, `dev finish`, and `dev block` status mutations. Creation records revision 1; update-like Events append one separate absolute revision Effect with matching prior/precondition values.
- Added focused regression, golden-replay, migration, live/rebuild-parity, and stale-replay tests. No expected-revision request handling, JSON response reshaping, or later aggregate families were added.

## Files created or modified

- Added `migrations/0012_p20_issue_revision.up.sql` and `migrations/0012_p20_issue_revision.down.sql`.
- Added `internal/issues/revision.go` and `internal/issues/revision_test.go`; updated scoped Issue create, update, transition state, and regression tests under `internal/issues/`.
- Updated legacy status mutation paths and their focused regressions under `internal/dev/`.
- Updated `internal/db/migrate_test.go`, `internal/journal/replay_test.go`, and added `testdata/journal/p20-issue-revision-effects.golden.jsonl`.
- Preserved without staging the pre-existing `project-docs/docs/phases/phase-20/run-ledger.md` modification and pre-existing untracked developer-stoppage review artifact.

## Documentation updated

- `project-docs/docs/phases/phase-20/tasks/task-6-done.md`

## Validation Contract Coverage

- VAL-010: PASS - `TestCreateEventOmitsRetiredPhase20Fields` at `internal/issues/issue_test.go:851` proves new Issue revision 1 with no revision prior/precondition; `TestUpdateEventRecordsReplayableRevision` at `internal/issues/update_test.go:80` and `TestTransitionAllowsWorkflowEdgeProjectsActivityAndRebuilds` at `internal/issues/transition_test.go:15` prove update/transition 1→2; `TestStartClaimsTransitionsAndOpensRunInOneEvent` at `internal/dev/start_test.go:20`, `TestFinishClosesRunCompletesTransitionsAndReleasesClaimInOneEvent` at `internal/dev/finish_test.go:14`, and `TestBlockRecordsFailedGateCompletionTransitionsAndReleasesClaim` at `internal/dev/block_test.go:13` prove legacy dev status Effects. `TestUpdateRejectsInvalidChangesWithoutAppending` at `internal/issues/update_test.go:131` and `TestUpdateRejectsNoOpWithoutAppending` at `internal/issues/update_test.go:457` prove rejected updates append nothing. `TestBuildIssueRevisionEffectsSortsObservedRevisionsAndDerivesResults` at `internal/issues/revision_test.go:5` proves sorted multi-Issue helper output without caller-supplied results.
- VAL-011: PASS - `testdata/journal/p20-issue-revision-effects.golden.jsonl` and `TestReplayJournalAppliesPhase20IssueRevisionGolden` at `internal/journal/replay_test.go:123` cover exact create/update/transition/dev-start/dev-finish/dev-block revision Effects and replay them to revision 6. `TestUpdateEventRecordsReplayableRevision` at `internal/issues/update_test.go:80` proves the live update tuple. `TestReplayJournalRejectsStaleIssueRevisionBeforeApplyingEffect` at `internal/journal/replay_test.go:283` proves stale prior/precondition replay rejects while the stored revision remains 2.
- VAL-012 (revision half): PASS - `TestIssueRevisionMigrationAddsAndRemovesRevision` at `internal/db/migrate_test.go:44` proves up adds default 1 and down removes the column. `TestUpdateAppendsOneEventProjectsActivityAndRebuilds` at `internal/issues/update_test.go:205` compares live and rebuilt Issue revisions at `internal/issues/update_test.go:446`; the required current-session integration gate passed with `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t4.744s`.

## Tests run

- RED: `go test ./internal/issues -run '^TestUpdateEventRecordsReplayableRevision$' -count=1` - expected FAIL - `FAIL`; diagnostic: `update_test.go:115: missing Issue revision effect` (terminated `exit_code: 1`).
- GREEN: `go test ./internal/issues -run '^TestUpdateEventRecordsReplayableRevision$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/issues\t0.054s` (terminated `exit_code: 0`)
- `go test ./internal/issues ./internal/dev ./internal/db ./internal/journal -run '^(TestBuildIssueRevisionEffectsSortsObservedRevisionsAndDerivesResults|TestUpdateEventRecordsReplayableRevision|TestUpdateAppendsOneEventProjectsActivityAndRebuilds|TestCreateEventOmitsRetiredPhase20Fields|TestTransitionAllowsWorkflowEdgeProjectsActivityAndRebuilds|TestStartClaimsTransitionsAndOpensRunInOneEvent|TestFinishClosesRunCompletesTransitionsAndReleasesClaimInOneEvent|TestBlockRecordsFailedGateCompletionTransitionsAndReleasesClaim|TestIssueRevisionMigrationAddsAndRemovesRevision|TestReplayJournalRejectsStaleIssueRevisionBeforeApplyingEffect)$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/journal\t0.024s` (terminated `exit_code: 0`)
- `go test ./internal/journal -run '^TestReplayJournalAppliesPhase20IssueRevisionGolden$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/journal\t0.026s` (terminated `exit_code: 0`)
- `go test ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.702s` (terminated `exit_code: 0`)
- Real entrypoint: in a fresh temporary workspace, ran `/tmp/fid --project <workspace> init --prefix FID`, then `/tmp/fid issue create --kind task --title 'CLI revision target' --by developer`, and `/tmp/fid issue update FID-517105 --title 'CLI revised target' --by developer` - PASS - final journal Event contained `"values":{"revision":2}` with `"prior":{"revision":1}` and `"preconditions":{"revision":1}` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 80.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t77.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t4.744s` (terminated `exit_code: 0`)
- `go test -bench . -benchmem ./...` - UNVERIFIED - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.551s` (historical wrapper-polling record was not in canonical termination form)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

Fresh direct-recovery evidence captured on 2026-08-21:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t78.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t4.915s` (terminated `exit_code: 0`)
- `go test -bench . -benchmem ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t1.970s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `29c52469be0277ac9e3e84cd0c5f63cd16a06066`

## Commit subject

- `feat(issue): add replayable aggregate revisions`
