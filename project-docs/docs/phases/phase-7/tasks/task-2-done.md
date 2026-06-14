## Summary of what was implemented

- Added `fid issue ac verify` as the Phase 7 append-only acceptance-criterion verification write path.
- Implemented `internal/issues.VerifyAcceptanceCriterion` with validate-before-append checks for known role, allowed `--status`, allowed `--evidence-type`, existing AC, and AC-to-issue membership.
- Persisted each verification as exactly one `issue.ac_verified` journal event with one `ac_verifications` insert effect, preserving full history and leaving `ac_verifications.completion_id` NULL.
- Wired `cmd/fid/issue.go` with the new `issue ac verify` CLI surface, required-flag stderr behavior, and the design-specified confirmation output including the optional `evidence:` line.
- Added engine and CLI tests for stored fields, multi-role independence, same-role fail->pass history retention, and rejection paths with unchanged journal state.

## Files created or modified

- Modified `internal/issues/ac.go`.
- Modified `internal/issues/ac_test.go`.
- Modified `cmd/fid/issue.go`.
- Modified `cmd/fid/issue_ac_test.go`.
- Added `project-docs/docs/phases/phase-7/tasks/task-2-done.md`.

## Documentation updated

- No change to `project-docs/specs/06_DESIGN_PLAN_CLI.md`; the implemented `fid issue ac verify` flags and stderr confirmation matched the existing Phase 7 design text.
- Added this completion report at `project-docs/docs/phases/phase-7/tasks/task-2-done.md`.

## Validation Contract Coverage

- `VAL-003` PASS via `TestVerifyAcceptanceCriterionPreservesHistoryAndLatestPerRole` and `TestIssueACVerifyCommandKeepsLatestPerRoleIndependent`, which prove independent latest statuses for `developer` and `qa_lead` on the same AC.
- `VAL-004` PASS via `TestVerifyAcceptanceCriterionPreservesHistoryAndLatestPerRole` and `TestIssueACVerifyCommandPreservesFailToPassHistory`, which prove same-role `fail` then `pass` retains two rows and updates latest-per-role to `pass`.
- `VAL-005` PASS via `TestVerifyAcceptanceCriterionPersistsAllFieldsAndAppendsEvent`, `TestVerifyAcceptanceCriterionRejectsBeforeAppend`, `TestIssueACVerifyCommandPersistsVerificationAndPrintsConfirmation`, and `TestIssueACVerifyCommandRejectsBeforeAppend`, which prove field persistence, `completion_id` NULL, and invalid-status rejection before append.
- Transcript excerpt, pass verify:

```text
$ /tmp/fid issue ac verify FID-TASK2-V21AYI-048311 AC-07c6d4 --by developer --status pass --evidence-type test --evidence-ref test/journal_test.go::TestAppendFlush
verified AC-07c6d4 on FID-TASK2-V21AYI-048311 by developer: pass
  evidence: test → test/journal_test.go::TestAppendFlush
```

- Transcript excerpt, multi-role latest-per-role:

```text
$ /tmp/fid issue ac verify FID-TASK2-V21AYI-048311 AC-07c6d4 --by qa_lead --status fail --notes Needs\ explicit\ QA\ repro
verified AC-07c6d4 on FID-TASK2-V21AYI-048311 by qa_lead: fail

$ /tmp/fid issue get FID-TASK2-V21AYI-048311 --section acceptance --format json
[{"id":"AC-07c6d4","ordinal":1,"text":"Developer and QA can verify independently","latest_verifications":{"developer":{"status":"pass","verified_at":"2026-06-14T03:39:26Z"},"qa_lead":{"status":"fail","verified_at":"2026-06-14T03:39:26Z"}}},{"id":"AC-b36060","ordinal":2,"text":"Developer can move from fail to pass","latest_verifications":{"developer":null}}]
```

- Transcript excerpt, fail then pass retains latest pass:

```text
$ /tmp/fid issue ac verify FID-TASK2-V21AYI-048311 AC-b36060 --by developer --status fail --notes Initial\ failure
verified AC-b36060 on FID-TASK2-V21AYI-048311 by developer: fail

$ /tmp/fid issue ac verify FID-TASK2-V21AYI-048311 AC-b36060 --by developer --status pass --notes Regression\ fixed
verified AC-b36060 on FID-TASK2-V21AYI-048311 by developer: pass

$ /tmp/fid issue get FID-TASK2-V21AYI-048311 --section acceptance --format json
[{"id":"AC-07c6d4","ordinal":1,"text":"Developer and QA can verify independently","latest_verifications":{"developer":{"status":"pass","verified_at":"2026-06-14T03:39:26Z"},"qa_lead":{"status":"fail","verified_at":"2026-06-14T03:39:26Z"}}},{"id":"AC-b36060","ordinal":2,"text":"Developer can move from fail to pass","latest_verifications":{"developer":{"status":"pass","verified_at":"2026-06-14T03:39:27Z"}}}]
```

- Transcript excerpt, invalid status rejected:

```text
$ /tmp/fid issue ac verify FID-TASK2-V21AYI-048311 AC-b36060 --by developer --status bogus
error: verify acceptance criterion: status "bogus" must be one of pass, fail, blocked, not_applicable
[exit 1]
```

## Tests run

- `go test ./internal/issues -run TestVerifyAcceptanceCriterionPersistsAllFieldsAndAppendsEvent`
  FAIL first: `undefined: VerifyAcceptanceCriterion` and `undefined: VerifyAcceptanceCriterionOptions`.
- `go test ./internal/issues -run TestVerifyAcceptanceCriterionPersistsAllFieldsAndAppendsEvent`
  PASS after implementing the engine write path.
- `go test ./internal/issues -run 'TestVerifyAcceptanceCriterion(PersistsAllFieldsAndAppendsEvent|PreservesHistoryAndLatestPerRole|RejectsBeforeAppend)'`
  PASS.
- `go test ./cmd/fid -run TestIssueACVerifyCommandPersistsVerificationAndPrintsConfirmation`
  FAIL first: `unknown flag: --by`.
- `go test ./cmd/fid -run TestIssueACVerifyCommandPersistsVerificationAndPrintsConfirmation`
  PASS after wiring the CLI subcommand.
- `go test ./cmd/fid -run 'TestIssueACVerifyCommand(PersistsVerificationAndPrintsConfirmation|KeepsLatestPerRoleIndependent|PreservesFailToPassHistory|RejectsBeforeAppend)'`
  PASS.
- `go test ./internal/issues ./cmd/fid`
  PASS.
- `gofmt -w internal/issues/ac.go internal/issues/ac_test.go cmd/fid/issue.go cmd/fid/issue_ac_test.go`
  PASS.

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
  Triggered because this report includes real CLI transcript excerpts for the public-output change.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit.

## Commit subject

- Pending commit.
