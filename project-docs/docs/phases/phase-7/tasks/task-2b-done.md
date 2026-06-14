## Summary of what was implemented

- Switched acceptance-criteria latest-per-role resolution to journal append order by selecting the newest `ac_verifications.rowid` per `(ac_id, verifier_role)` in [detail.go](/home/mark/projects/fid/internal/issues/detail.go:156), matching DEC-021 and removing the old `verified_at`/random-id tie break.
- Kept the AC write-path split behavior-preserving by moving `VerifyAcceptanceCriterion` and its helpers into [ac_verify.go](/home/mark/projects/fid/internal/issues/ac_verify.go:45), leaving [ac.go](/home/mark/projects/fid/internal/issues/ac.go:32) at 190 lines and `ac_verify.go` at 255 lines.
- Replaced the sleep-based fail→pass coverage with deterministic no-sleep tests in [ac_test.go](/home/mark/projects/fid/internal/issues/ac_test.go:586) and [issue_ac_test.go](/home/mark/projects/fid/cmd/fid/issue_ac_test.go:289), and added rebuild-parity coverage proving latest-per-role stays identical after deleting the DB and running `fid rebuild`.
- Confirmed `git show HEAD:project-docs/docs/phases/phase-7/tasks/task-1-done.md` already contains the finalized Task 1 hash/subject, so no duplicate docs-only commit was needed.

## Files created or modified

- Modified [cmd/fid/issue_ac_test.go](/home/mark/projects/fid/cmd/fid/issue_ac_test.go:289).
- Modified [internal/issues/ac.go](/home/mark/projects/fid/internal/issues/ac.go:32).
- Created [internal/issues/ac_verify.go](/home/mark/projects/fid/internal/issues/ac_verify.go:1).
- Modified [internal/issues/ac_test.go](/home/mark/projects/fid/internal/issues/ac_test.go:586).
- Modified [internal/issues/detail.go](/home/mark/projects/fid/internal/issues/detail.go:156).
- Modified [project-docs/docs/phases/phase-7/decisions.md](/home/mark/projects/fid/project-docs/docs/phases/phase-7/decisions.md:6).
- Modified [project-docs/docs/phases/phase-7/tasks/task-2b.md](/home/mark/projects/fid/project-docs/docs/phases/phase-7/tasks/task-2b.md:79).
- Created [project-docs/docs/phases/phase-7/tasks/task-2b-done.md](/home/mark/projects/fid/project-docs/docs/phases/phase-7/tasks/task-2b-done.md:1).
- Modified [project-docs/specs/03_TRD_SCHEMA.md](/home/mark/projects/fid/project-docs/specs/03_TRD_SCHEMA.md:236).

## Documentation updated

- Updated [project-docs/specs/03_TRD_SCHEMA.md](/home/mark/projects/fid/project-docs/specs/03_TRD_SCHEMA.md:236) so `ac_verifications` latest-per-role semantics explicitly follow journal append order and keep `verified_at` as display-only.
- Updated [project-docs/docs/phases/phase-7/tasks/task-2b.md](/home/mark/projects/fid/project-docs/docs/phases/phase-7/tasks/task-2b.md:79) with the DEC-022 scope correction that limits the 400-line tripwire to `ac.go` and `ac_verify.go`.
- Replaced the blocker report with this completion report at [project-docs/docs/phases/phase-7/tasks/task-2b-done.md](/home/mark/projects/fid/project-docs/docs/phases/phase-7/tasks/task-2b-done.md:1).

## Validation Contract Coverage

- `VAL-004` PASS via `TestLatestVerificationSelectionPrefersAppendOrderWithinTimestampTie`, `TestVerifyAcceptanceCriterionPreservesHistoryAndLatestPerRole`, and `TestIssueACVerifyCommandPreservesFailToPassHistory`. These prove immediate same-role fail→pass with no `time.Sleep`, two retained verification rows, latest-per-role=`pass`, and after-fail latest=`fail`.
- `VAL-009` PASS via `TestVerifyAcceptanceCriterionLatestPerRoleSurvivesRebuild` and `TestIssueACVerifyCommandLatestPerRoleSurvivesRebuild`. These prove delete-DB + rebuild reproduces identical latest-per-role state under the new append-order rule.
- `VAL-003` remains green via `TestIssueACVerifyCommandKeepsLatestPerRoleIndependent` and the full package test run, which preserve independent latest state per verifier role.
- `VAL-005` remains green via `TestVerifyAcceptanceCriterionPersistsAllFieldsAndAppendsEvent`, `TestIssueACVerifyCommandPersistsVerificationAndPrintsConfirmation`, and the full package test run, which preserve field persistence and pre-append rejection behavior.
- Deterministic CLI transcript evidence for DEC-021 / ESC-005:

```text
$ /tmp/fid issue ac verify TMP-cf4df0 AC-e63c89 --by developer --status fail
verified AC-e63c89 on TMP-cf4df0 by developer: fail

$ /tmp/fid issue get TMP-cf4df0 --section acceptance --format json
[{"id":"AC-e63c89","ordinal":1,"text":"Latest verification follows append order","latest_verifications":{"developer":{"status":"fail","verified_at":"2026-06-14T07:41:02Z"}}}]

$ /tmp/fid issue ac verify TMP-cf4df0 AC-e63c89 --by developer --status pass
verified AC-e63c89 on TMP-cf4df0 by developer: pass

$ /tmp/fid issue ac verify TMP-cf4df0 AC-e63c89 --by qa_lead --status fail
verified AC-e63c89 on TMP-cf4df0 by qa_lead: fail

$ /tmp/fid issue get TMP-cf4df0 --section acceptance --format json
[{"id":"AC-e63c89","ordinal":1,"text":"Latest verification follows append order","latest_verifications":{"developer":{"status":"pass","verified_at":"2026-06-14T07:41:02Z"},"qa_lead":{"status":"fail","verified_at":"2026-06-14T07:41:02Z"}}}]
```

- Rebuild parity transcript evidence:

```text
$ rm /tmp/tmp.5nQ7YVwzD8/.fid/fid.sqlite

$ /tmp/fid rebuild --by developer
rebuilding .fid/fid.sqlite from journal...
  replayed 6 events from 1 journal files
  schema version: 1
  elapsed: 12ms
rebuild complete

$ /tmp/fid issue get TMP-cf4df0 --section acceptance --format json
[{"id":"AC-e63c89","ordinal":1,"text":"Latest verification follows append order","latest_verifications":{"developer":{"status":"pass","verified_at":"2026-06-14T07:41:02Z"},"qa_lead":{"status":"fail","verified_at":"2026-06-14T07:41:02Z"}}}]

$ test before_rebuild = after_rebuild
before_after_rebuild_match
```

## Tests run

- `go test ./internal/issues -run 'TestLatestVerificationSelectionPrefersAppendOrderWithinTimestampTie|TestVerifyAcceptanceCriterionPreservesHistoryAndLatestPerRole|TestVerifyAcceptanceCriterionLatestPerRoleSurvivesRebuild'`
  PASS.
- `go test ./cmd/fid -run 'TestIssueACVerifyCommandPreservesFailToPassHistory|TestIssueACVerifyCommandLatestPerRoleSurvivesRebuild|TestIssueACVerifyCommandKeepsLatestPerRoleIndependent'`
  PASS.
- `wc -l internal/issues/ac.go internal/issues/ac_verify.go`
  PASS.
  `190 internal/issues/ac.go`
  `255 internal/issues/ac_verify.go`
- `git show HEAD:project-docs/docs/phases/phase-7/tasks/task-1-done.md`
  PASS.
  Confirmed HEAD already contains:
  `59f02df`
  `feat(issue): add acceptance criteria command`

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
  Output: `0 issues.`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
  `cmd/fid`: `79.2%`
  `internal/issues`: `82.6%`
  `internal/journal`: `80.3%`
  `internal/policy`: `85.0%`
  `internal/workflow`: `95.5%`
- `PASS` `go tool cover -func=coverage.out`
  Total coverage: `77.2%`
- `PASS` `go build -o /tmp/fid ./cmd/fid`
  Required by this task to capture deterministic fail→pass and rebuild parity transcript evidence.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending git commit.

## Commit subject

- Pending git commit subject.
