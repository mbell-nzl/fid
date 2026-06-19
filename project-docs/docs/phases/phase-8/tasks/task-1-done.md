## Summary of what was implemented

Implemented `fid issue complete` as the Phase 8 Task 1 write surface.

- Added `internal/issues.Complete` with validate-before-append guards, derived attempt numbering under the write lock, append-only completion inserts, prior-attempt supersession, human-intervention storage, and epic rejection before journal append.
- Added the `fid issue complete <id>` CLI with `--by`, `--notes-md`, `--tokens`, `--commits`, `--files`, `--human-intervention`, and `--intervention-note`.
- Extended issue detail read-back so completion records expose notes path, token count, commit hashes, files touched, intervention state, and intervention note across structured output formats.
- Added Task 1 test coverage for attempt creation, supersession, human-intervention states, epic rejection, unknown-role rejection, and missing-issue rejection.

## Files created or modified

- `cmd/fid/issue.go`
- `cmd/fid/issue_complete_test.go`
- `internal/issues/complete.go`
- `internal/issues/complete_test.go`
- `internal/issues/detail.go`
- `internal/issues/issue.go`
- `internal/format/toon_test.go`
- `testdata/format/issue-empty.golden.toon`
- `project-docs/docs/phases/phase-8/tasks/task-1-done.md`

## Documentation updated

- Verified `project-docs/specs/06_DESIGN_PLAN_CLI.md` section `fid issue complete`; no spec change was required.
- Added this completion report at `project-docs/docs/phases/phase-8/tasks/task-1-done.md`.

## Validation Contract Coverage

- `VAL-001` PASS.
  Evidence:
  `TestCompleteCreatesAttemptOneAndPersistsCompletionFields`
  `TestIssueCompleteCommandCreatesAttemptOneAndReturnsJSONCompletionFields`
  Real CLI transcript shows attempt 1 creation plus JSON read-back of `notes_md_path`, `actual_tokens`, `commit_hashes`, and `files_touched`.
- `VAL-002` PASS.
  Evidence:
  `TestCompleteSupersedesPriorAttemptAndRetainsHistory`
  `TestIssueCompleteCommandSupersedesPriorAttemptAndStoresInterventionFields`
  Real CLI transcript shows attempt 2 creation; JSON read-back shows attempt 2 as `dev_complete` and attempt 1 retained as `superseded`.
- `VAL-003` PASS.
  Evidence:
  `TestCompleteStoresHumanInterventionStateWithoutActivityLogRows`
  `TestIssueCompleteCommandSupersedesPriorAttemptAndStoresInterventionFields`
  Real CLI transcript shows a human-intervention attempt with `human_intervention=true` and `human_intervention_note="paired on schema fix"`.
- `VAL-004` PASS.
  Evidence:
  `TestCompleteRejectsEpicBeforeAppend`
  `TestIssueCompleteCommandRejectsEpicBeforeAppend`
  `TestIssueCompleteCommandRejectsUnknownRoleAndMissingIssueBeforeAppend`
  Real CLI transcript shows epic rejection with non-zero exit; tests assert unchanged journal line counts and zero `task_completions` rows for rejected writes.

Real CLI transcript evidence:

```text
$ /tmp/fid init
initialized fid project in .fid/
  prefix: FID-TASK1-QRMECS
  config: .fid/config.toml
  pricing: .fid/pricing.toml
  journal: .fid/journal/2026-06-19.jsonl
  event: EVT-20260619-000001 project.initialized
  roles seeded: developer, qa_reviewer, qa_lead, tl_lite, tl, workflow_runner, human

$ /tmp/fid issue create --kind task --title "Transcript task" --by developer
created issue FID-TASK1-QRMECS-083915 (task) "Transcript task"

$ /tmp/fid issue complete FID-TASK1-QRMECS-083915 --by developer --notes-md project-docs/docs/phases/phase-8/tasks/task-1-done.md --tokens 14823 --commits abc123,def456 --files internal/issues/complete.go --files cmd/fid/issue.go
completed FID-TASK1-QRMECS-083915 (attempt 1)
  completion: TC-a4b21b
  tokens: 14823
  commits: abc123, def456

$ /tmp/fid issue complete FID-TASK1-QRMECS-083915 --by developer --tokens 15001 --human-intervention --intervention-note "paired on schema fix"
completed FID-TASK1-QRMECS-083915 (attempt 2)
  completion: TC-7ecaad
  tokens: 15001

$ /tmp/fid issue get FID-TASK1-QRMECS-083915 --format json
{"id":"FID-TASK1-QRMECS-083915","kind":"task","title":"Transcript task","status":"todo","status_category":"unstarted","workflow_id":"task-default","parent_id":null,"priority":null,"sequence":null,"spec_md_path":null,"est_tokens":null,"due_date":null,"owner_role":"developer","owner_actor":null,"labels":[],"in_inbox":false,"concerns_raised":false,"open_concerns_count":0,"triage_status":null,"acceptance_criteria":[],"created_at":"2026-06-19T22:38:39Z","updated_at":"2026-06-19T22:38:39Z","completions":[{"id":"TC-7ecaad","attempt_number":2,"developer_role":"developer","model":null,"status":"dev_complete","notes_md_path":null,"actual_tokens":15001,"commit_hashes":[],"files_touched":[],"gate_results":[],"human_intervention":true,"human_intervention_note":"paired on schema fix","completed_at":"2026-06-19T22:38:39Z"},{"id":"TC-a4b21b","attempt_number":1,"developer_role":"developer","model":null,"status":"superseded","notes_md_path":"project-docs/docs/phases/phase-8/tasks/task-1-done.md","actual_tokens":14823,"commit_hashes":["abc123","def456"],"files_touched":["internal/issues/complete.go","cmd/fid/issue.go"],"gate_results":[],"human_intervention":false,"human_intervention_note":null,"completed_at":"2026-06-19T22:38:39Z"}]}

$ /tmp/fid issue create --kind epic --title "Transcript epic" --scheduling sequential --by tl
created issue FID-TASK1-QRMECS-66d056 (epic) "Transcript epic"

$ /tmp/fid issue complete FID-TASK1-QRMECS-66d056 --by developer
error: complete issue: issue "FID-TASK1-QRMECS-66d056" is kind "epic"; task completions require task or subtask
(exit=1)
```

## Tests run

- `go test ./internal/issues ./cmd/fid -run 'Test(CompleteCreatesAttemptOneAndPersistsCompletionFields|CompleteSupersedesPriorAttemptAndRetainsHistory|CompleteStoresHumanInterventionStateWithoutActivityLogRows|CompleteRejectsEpicBeforeAppend|IssueCompleteCommandCreatesAttemptOneAndReturnsJSONCompletionFields|IssueCompleteCommandSupersedesPriorAttemptAndStoresInterventionFields|IssueCompleteCommandRejectsEpicBeforeAppend)$'`
- `go test ./internal/format ./cmd/fid -run 'Test(WriteTOONMatchesGoldenFixtureForEmptyIssue|IssueGetReturnsYAMLForEmptyAndPopulatedIssues|IssueGetDispatchesStructuredFormatsAndRejectsUnknownFormats)$'`
- `go test ./cmd/fid -run 'TestIssueCompleteCommandRejectsUnknownRoleAndMissingIssueBeforeAppend$'`
- `gofmt -l cmd internal`
- `golangci-lint run ./...`
- `go test -race -coverprofile=coverage.out ./...`
- `go tool cover -func=coverage.out`
- `go build -o /tmp/fid ./cmd/fid`

## Exact quality-gate results

- `gofmt -l cmd internal` PASS
- `golangci-lint run ./...` PASS
- `go test -race -coverprofile=coverage.out ./...` PASS
  Coverage highlights:
  `cmd/fid` 78.7%
  `internal/issues` 81.9%
  `internal/journal` 80.3%
  `internal/policy` 85.0%
  `internal/workflow` 95.5%
- `go tool cover -func=coverage.out` PASS
  Overall statement coverage: 76.9%
- `go build -o /tmp/fid ./cmd/fid` PASS
  Reason run:
  required to generate the real CLI transcript above

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- TBD

## Commit subject

- TBD
