## Summary of what was implemented

- Added the read-only `internal/issues.ListAcceptanceCriteria` path in `internal/issues/ac_list.go`, using `journal.EnsureProjectionFresh`, a read-only SQLite DSN, and the existing `issueAcceptanceCriteria` / `latestVerificationsByRole` helpers so `ac list` matches `issue get --section acceptance`.
- Wired `fid issue ac list <issue-id>` into `cmd/fid/issue.go` with shared format-engine dispatch and `json` as the default `--format`.
- Extended the markdown formatter in `internal/format/markdown.go` so `issues.AcceptanceCriteriaList` renders through the existing acceptance-criteria markdown helper.
- Added Task 3 engine and CLI coverage for ordinal ordering, latest-per-role output, default JSON parity, `toon` / `markdown` / `yaml` dispatch, empty-list behavior, and missing-issue rejection.
- Aligned the `fid issue ac list` JSON example in `project-docs/specs/06_DESIGN_PLAN_CLI.md` with the actual output shape.

## Files created or modified

- Created `internal/issues/ac_list.go`.
- Modified `internal/issues/ac_test.go`.
- Modified `cmd/fid/issue.go`.
- Modified `cmd/fid/issue_ac_test.go`.
- Modified `internal/format/markdown.go`.
- Modified `project-docs/specs/06_DESIGN_PLAN_CLI.md`.
- Created `project-docs/docs/phases/phase-7/tasks/task-3-done.md`.

## Documentation updated

- Updated `project-docs/specs/06_DESIGN_PLAN_CLI.md` so the `fid issue ac list` JSON example matches the implemented `latest_verifications` payload.
- Added this completion report at `project-docs/docs/phases/phase-7/tasks/task-3-done.md`.

## Validation Contract Coverage

- `VAL-006` PASS via `TestListAcceptanceCriteriaReturnsOrdinalLatestPerRoleAndMatchesIssueGet`, `TestListAcceptanceCriteriaEmptyAndNotFound`, `TestIssueACListCommandDefaultsToJSONAndMatchesIssueGetAcceptance`, and `TestIssueACListCommandDispatchesFormatsAndHandlesEmptyAndMissingIssues`. These prove ordinal ordering, parity with `issue get --section acceptance`, latest-per-role output keyed by role, empty-list success, and missing-issue rejection.
- `VAL-007` PASS via `TestIssueACListCommandDefaultsToJSONAndMatchesIssueGetAcceptance` and `TestIssueACListCommandDispatchesFormatsAndHandlesEmptyAndMissingIssues`. These prove omitted `--format` equals JSON and that `json`, `toon`, `markdown`, and `yaml` all render successfully.
- Real CLI transcript excerpt, default JSON (`--format` omitted):

```text
$ /tmp/fid issue ac list FID-TASK3-RJG9NC-cfd9a0
[{"id":"AC-823ed4","ordinal":1,"text":"Journal file is flushed to disk before returning success","latest_verifications":{"developer":{"status":"pass","verified_at":"2026-06-14T08:26:40Z"},"qa_lead":{"status":"blocked","verified_at":"2026-06-14T08:26:45Z"}}},{"id":"AC-163ae3","ordinal":2,"text":"Partial write leaves journal intact","latest_verifications":{"developer":null}}]
```

- Real CLI transcript excerpt, `--format json`:

```text
$ /tmp/fid issue ac list FID-TASK3-RJG9NC-cfd9a0 --format json
[{"id":"AC-823ed4","ordinal":1,"text":"Journal file is flushed to disk before returning success","latest_verifications":{"developer":{"status":"pass","verified_at":"2026-06-14T08:26:40Z"},"qa_lead":{"status":"blocked","verified_at":"2026-06-14T08:26:45Z"}}},{"id":"AC-163ae3","ordinal":2,"text":"Partial write leaves journal intact","latest_verifications":{"developer":null}}]
```

- Real CLI transcript excerpt, `--format toon`:

```text
$ /tmp/fid issue ac list FID-TASK3-RJG9NC-cfd9a0 --format toon
[{id:AC-823ed4,latest_verifications:{developer:{status:pass,verified_at:"2026-06-14T08:26:40Z"},qa_lead:{status:blocked,verified_at:"2026-06-14T08:26:45Z"}},ordinal:1,text:"Journal file is flushed to disk before returning success"},{id:AC-163ae3,latest_verifications:{developer:null},ordinal:2,text:"Partial write leaves journal intact"}]
```

- Real CLI transcript excerpt, `--format markdown`:

```text
$ /tmp/fid issue ac list FID-TASK3-RJG9NC-cfd9a0 --format markdown
# Acceptance Criteria

- [x] AC-823ed4: Journal file is flushed to disk before returning success
  - developer: pass (2026-06-14T08:26:40Z)
  - qa_lead: blocked (2026-06-14T08:26:45Z)
- [ ] AC-163ae3: Partial write leaves journal intact
  - developer: (not verified)
```

- Real CLI transcript excerpt, `--format yaml`:

```text
$ /tmp/fid issue ac list FID-TASK3-RJG9NC-cfd9a0 --format yaml
- id: AC-823ed4
  latest_verifications:
    developer:
        status: pass
        verified_at: "2026-06-14T08:26:40Z"
    qa_lead:
        status: blocked
        verified_at: "2026-06-14T08:26:45Z"
  ordinal: 1
  text: Journal file is flushed to disk before returning success
- id: AC-163ae3
  latest_verifications:
    developer: null
  ordinal: 2
  text: Partial write leaves journal intact
```

- Real CLI transcript excerpt, empty list:

```text
$ /tmp/fid issue ac list FID-TASK3-RJG9NC-5d6bcf --format json
[]
```

## Tests run

- Resumed Task 3 verification: the existing focused tests were already green when this session resumed, so no additional behavior patch was required before running the close-out gate.
- `go test ./internal/issues -run 'TestListAcceptanceCriteria(ReturnsOrdinalLatestPerRoleAndMatchesIssueGet|EmptyAndNotFound)$'`
  PASS.
- `go test ./cmd/fid -run 'TestIssueACListCommand(DefaultsToJSONAndMatchesIssueGetAcceptance|DispatchesFormatsAndHandlesEmptyAndMissingIssues)$'`
  PASS.
- `gofmt -l cmd internal`
  PASS.
- `golangci-lint run ./...`
  PASS.
- `go test -race -coverprofile=coverage.out ./...`
  PASS.
- `go tool cover -func=coverage.out`
  PASS.
- `go build -o /tmp/fid ./cmd/fid`
  PASS.

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
  Output: `0 issues.`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
  `cmd/fid`: `78.9%`
  `internal/issues`: `82.0%`
  `internal/journal`: `80.3%`
  `internal/policy`: `85.0%`
  `internal/workflow`: `95.5%`
- `PASS` `go tool cover -func=coverage.out`
  Total coverage: `76.7%`
- `PASS` `go build -o /tmp/fid ./cmd/fid`
  Required by this task to capture real `fid issue ac list` transcript evidence in all supported formats plus the empty-list case.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `1806aa4`

## Commit subject

- `feat(issue): add ac list command`
