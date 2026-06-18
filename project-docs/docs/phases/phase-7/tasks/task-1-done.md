## Summary of what was implemented

- Added `fid issue ac add` as the first Phase 7 acceptance-criteria write path.
- Implemented `internal/issues.AddAcceptanceCriterion` with validate-before-append guards for unknown roles, missing issues, invalid ordinals, and epic rejection before any journal append.
- Projected each accepted write as exactly one `issue.ac_added` journal event with one `acceptance_criteria` insert effect and no `activity_log` effect, preserving DEC-020.
- Wired `cmd/fid/issue.go` with the new `issue ac add` CLI surface and required-flag stderr behavior.
- Added engine and CLI tests for ordered persistence/read-back, epic rejection, required flags, invalid input, no-append rejection paths, and public confirmation output.

## Files created or modified

- Modified `cmd/fid/issue.go`.
- Created `cmd/fid/issue_ac_test.go`.
- Created `internal/issues/ac.go`.
- Created `internal/issues/ac_test.go`.
- Created `project-docs/docs/phases/phase-7/tasks/task-1-done.md`.

## Documentation updated

- No change to `project-docs/specs/06_DESIGN_PLAN_CLI.md`; the implemented `fid issue ac add` flags and stderr confirmation matched the existing Phase 7 design text.
- Added this completion report at `project-docs/docs/phases/phase-7/tasks/task-1-done.md`.

## Validation Contract Coverage

- `VAL-001` PASS via `TestAddAcceptanceCriterionPersistsAndReadsBackInOrdinalOrder`, `TestAddAcceptanceCriterionTieBreaksByID`, and `TestIssueACAddCommandAddsAcceptanceCriterionAndPrintsConfirmation`, which prove `AC-` ID creation, persisted text/ordinal, ordinal-then-id read-back, and CLI confirmation output.
- `VAL-002` PASS via `TestAddAcceptanceCriterionRejectsBeforeAppend`, `TestAddAcceptanceCriterionNotFoundErrorShape`, `TestAddAcceptanceCriterionInputValidation`, and `TestIssueACAddCommandRejectsBeforeAppend`, which prove epic rejection, unknown-role rejection, invalid-ordinal rejection, required-flag rejection, clear not-found behavior, and unchanged journal lines / zero projected AC rows on rejected writes.
- Transcript excerpt, successful add:

```text
$ /tmp/fid issue create --kind task --title "CLI transcript task" --by developer
created issue TMP.WQNDUDLWZH-da63e9 (task) "CLI transcript task"

$ /tmp/fid issue ac add TMP.WQNDUDLWZH-da63e9 --text "Journal file is flushed to disk before returning success" --ordinal 1 --by developer
added AC AC-18bda2 on TMP.WQNDUDLWZH-da63e9 (ordinal 1): "Journal file is flushed to disk before returning success"

$ /tmp/fid issue get TMP.WQNDUDLWZH-da63e9 --section acceptance --format json
[{"id":"AC-18bda2","ordinal":1,"text":"Journal file is flushed to disk before returning success","latest_verifications":{"developer":null}}]
```

- Transcript excerpt, epic reject before append:

```text
$ /tmp/fid issue create --kind epic --title "CLI transcript epic" --scheduling prioritised --by tl
created issue TMP.WQNDUDLWZH-cdfff2 (epic) "CLI transcript epic"

$ wc -l /tmp/tmp.WqnDUDLWzh/.fid/journal/2026-06-14.jsonl
4 /tmp/tmp.WqnDUDLWzh/.fid/journal/2026-06-14.jsonl

$ /tmp/fid issue ac add TMP.WQNDUDLWZH-cdfff2 --text "Epic AC should fail" --ordinal 1 --by developer
error: add acceptance criterion: issue "TMP.WQNDUDLWZH-cdfff2" is kind "epic"; acceptance criteria require task or subtask

$ wc -l /tmp/tmp.WqnDUDLWzh/.fid/journal/2026-06-14.jsonl
4 /tmp/tmp.WqnDUDLWzh/.fid/journal/2026-06-14.jsonl
```

## Tests run

- `go test ./internal/issues ./cmd/fid`
  Red: failed first because `AddAcceptanceCriterion` and the `issue ac add` command were not implemented.
- `go test ./internal/issues ./cmd/fid`
  PASS after implementing the engine path and CLI wiring.
- `gofmt -w internal/issues/ac.go internal/issues/ac_test.go cmd/fid/issue.go cmd/fid/issue_ac_test.go`
  PASS.
- `go test ./internal/issues ./cmd/fid`
  PASS after formatting.

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
  Output: `0 issues.`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
  `cmd/fid`: `79.1%`
  `internal/issues`: `82.8%`
  `internal/journal`: `80.3%`
  `internal/policy`: `85.0%`
  `internal/workflow`: `95.5%`
- `PASS` `go tool cover -func=coverage.out`
  Total coverage: `77.0%`
- `PASS` `go build -o /tmp/fid ./cmd/fid`
  Triggered because the report includes real CLI transcript excerpts for the public-output change.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `59f02df`

## Commit subject

- `feat(issue): add acceptance criteria command`
