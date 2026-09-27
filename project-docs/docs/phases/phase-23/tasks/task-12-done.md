# Phase 23 - Task 12: Requirement references and a single traceability model - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: session-continuation
- Condition: The real-process requirement journey is rejected while creating its Activation authority profile. The closed capability registry at `internal/workflow/activation.go:832` omits `requirement.record-reference` and `requirement.link`. That file is outside task 12's named write scope, and the task explicitly routes an unowned closed-registry rejection to Runner triage.
- Recovery attempts: 1
- Next action: Runner triages the registry scope gap and supplies authority for a same-task continuation. Continue task 12, converge retained tests and fixtures, complete the V1 removals and documentation, then run every exact gate in order.
- Human attention required: no

## Summary of what was implemented

- Added a reference writer with Source digest verification under the journal write lock, an effect plan that supersedes the prior active row, a link writer, and read paths. This is partial: the real-process journey has not passed.
- Registered the two requirement result shapes, Source JSON column, projection result branches, and supersession update fields in the current Event validators; removed the two remaining judgment fence keys and the empty guard. Removed the V1 sync Event rules from the current and archive catalogs.
- Replaced the `requirement sync` command tree with `record-reference`, `link`, `get`, and `list`. The V1 package and its retained tests remain pending convergence.

## Files created or modified

- `internal/review/requirement.go`, `internal/review/requirement_test.go`.
- `cmd/fid/requirement.go`, `cmd/fid/requirement_surface_test.go`, `cmd/fid/integration/p23_requirement_reference_test.go`.
- `internal/journal/current_event_catalog.go`, `current_event_json.go`, `current_event_result_projection.go`, `current_event_updates.go`, `current_event_value_catalog.go`, `project_effects.go`.
- `internal/export/archive_semantics_catalog.go`.
- `project-docs/docs/phases/phase-23/tasks/task-12-done.md`.

## Documentation updated

- This checkpoint report only. The task's two named specs remain to be updated in the continuation.

## Validation Contract Coverage

- VAL-003: BLOCKED - The fence keys and guard were removed from `internal/journal/project_effects.go`; the retained parity test and full replay proof remain to be converged.
- VAL-051: BLOCKED - `TestRequirementRecordReferenceSupersedesPriorActiveRow` passes for the effect plan, but `TestP23RequirementReferenceAndLinkRealProcess` fails during Activation setup before the writer can run. Digest rejection, readable supersession, and exactly one active row remain unverified.
- VAL-052: BLOCKED - The link writer and typed relation set are present; the real-process journey stops before its first reference, so accepted and rejected link cases remain unverified.
- VAL-053: BLOCKED - The sync leaf was removed from `cmd/fid/requirement.go` and its Event rules from both catalogs. `internal/requirements` and retained corpus fixtures still require removal; the declared absence search has not run.
- VAL-054: BLOCKED - The `links` endpoint registry and retained tests have not yet been converged to Issue-only endpoints.

## Tests run

Working directory for all commands: `/home/mark/projects/fid`. Each result below is from the command's observed termination.

- `go test -run '^TestRequirementRecordReferenceSupersedesPriorActiveRow$' ./internal/review` - FAIL - `FAIL` (terminated `exit_code: 1`); compiler detail: `undefined: requirementReferenceEffects`.
- `go test -run '^TestRequirementRecordReferenceSupersedesPriorActiveRow$' ./internal/review` - PASS - `ok  	github.com/mbell-nzl/fid/internal/review	0.005s` (terminated `exit_code: 0`)
- `go test -run '^TestRequirementCommandHasReferenceLeaves$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`); all four new leaves were missing and `sync` remained.
- `go test -run '^TestRequirementCommandHasReferenceLeaves$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.131s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP23RequirementReferenceAndLinkRealProcess$' ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); first authority-profile setup returned `invalid activation input`.
- `go test -tags integration -run '^TestP23RequirementReferenceAndLinkRealProcess$' ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); the expanded baseline action set still returned `invalid activation input`. The closed registry at `internal/workflow/activation.go:832` omits both requirement actions.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestRequirementRecordReferenceSupersedesPriorActiveRow|TestRequirementCommandHasReferenceLeaves)$' ./internal/review ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.041s` (terminated `exit_code: 0`)

The focused race command also reported `ok  github.com/mbell-nzl/fid/internal/review 1.021s`.
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-12-done.md project-docs/docs/phases/phase-23/tasks/task-12.md` - FAIL - `FAIL report-lint found 4 failure(s)` (terminated `exit_code: 4`). Three were report punctuation or trailing commentary and are corrected here. The fourth is the task spec's expected-nonzero declaration, whose two-line form does not match the control's one-line grammar at `project-docs/tools/report-lint.sh:104`.
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-12-done.md project-docs/docs/phases/phase-23/tasks/task-12.md` - FAIL - `FAIL report-lint found 1 failure(s)` (terminated `exit_code: 1`). The report evidence format passed; the remaining failure is the task-spec declaration grammar.

## Exact quality-gate results

The implementation is unfinished and the closed-registry rejection prevents the required real-process acceptance journey. No exact task gate was run or claimed as a pass in this checkpoint. All seven commands remain mandatory for the continuation, in this order:

- `gofmt -l cmd internal` - BLOCKED - not executed; task incomplete and scope gap pending Runner triage (no termination observed).
- `golangci-lint run ./...` - BLOCKED - not executed; task incomplete and scope gap pending Runner triage (no termination observed).
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - not executed; task incomplete and scope gap pending Runner triage (no termination observed).
- `go tool cover -func=coverage.out` - BLOCKED - not executed; no task coverage profile exists (no termination observed).
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - not executed; the focused real-process journey fails during Activation setup (no termination observed).
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - BLOCKED - not executed; task incomplete and scope gap pending Runner triage (no termination observed).
- `rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal` - BLOCKED - not executed; `internal/requirements` still contains the V1 writer, so the expected exit-1 absence proof cannot yet pass (no termination observed).

DEC-227 package times for `cmd/fid/integration`, `cmd/fid`, and `internal/workflow` are not measured by this checkpoint. The continuation must run the exact DEC-259 race commands and enforce the 900-second tagged and 600-second ordinary budgets.

## Deviations from spec

- The declared real-process test exposed an unowned closed capability registry. No registry or control was changed. The task remains incomplete and requires Runner triage before its exact gate.
- Report-lint attempt 1 ended `FAIL report-lint found 4 failure(s)` (terminated `exit_code: 4`). Three failures came from report-line punctuation and trailing commentary; the report was corrected and committed. The byte-identical command then ended `FAIL report-lint found 1 failure(s)` (terminated `exit_code: 1`), retaining only the task-spec grammar mismatch. No task spec or control definition was changed.

## Concerns, technical debt, or notes for the Tech Lead

- The accepted contract requires the two requirement capabilities, while `internal/workflow/activation.go:832` rejects them and is absent from task 12's named scope. The QA plan's earlier scope review listed that file for other Phase 23 capability additions, but not task 12. This is a scope gap, not an owner assignment by this report.
- Report lint cannot pass against the current task spec: `FAIL expected non-zero declaration has invalid grammar` (first attempt, terminated `exit_code: 4`). Its parser requires the command and expected exit code on one bullet; task 12 puts the expected exit code on the following bullet. The task spec and control are outside this session's write scope.
- The current checkpoint is intentionally red outside the two focused green tests. In particular, retained sync and link assertions, replay fixtures, and the task-1 fence test need the accepted convergence work. No DEC-227 headroom claim is possible yet.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `0e4bb33cbe32fa9d4d83955b3121a4135edf78cd`

## Commit subject

- `feat(requirement): checkpoint task 12 reference convergence`
