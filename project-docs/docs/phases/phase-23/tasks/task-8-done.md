# Phase 23 - Task 8: Generic Concerns and append-only Deviations - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added generic Concern recording over the six closed `JudgmentSubject` types and one irreversible disposal transition with exactly `carried-forward`, `rejected`, and `superseded`. Subject references remain typed without requiring a matching row, including the deliberately unresolved `requirement` case; only disposition target identities must resolve.
- Removed completion linkage from the generic Concern shape and kept the retained V1 `concern raise`, `concern address`, and `concern list` leaves intact for Task 15 convergence. Raising and disposing Actors are stored in separate fields, including equal values for self-disposal and distinct values for a second activated Actor.
- Added append-only Deviation recording and filtered listing. A changed judgment creates a new row with `supersedes_deviation_id`; missing and self predecessors reject, both rows remain readable, and the live command inventory contains no update leaf.
- Added canonical JSON CLI writers, typed failures, exact idempotent replay, and the three grant actions in both sorted closed allowlists. Concern and Deviation remain outside the Evidence type domain.
- Unfenced exactly `concern.recorded`, `concern.disposed`, and `deviation.recorded`; added their strict result grammars, value shapes, JSON-column registrations, projection fill, aggregate bindings, exact Concern update admission, and live/rebuild result parity. `JudgmentSubject` was reused unchanged and `Subject` was not changed.
- Added the explicit `open` to `disposed` terminal guard in `validateCurrentTerminalValues`; the update also binds the disposal Event and disposing Actor. Open blocker Concerns do not participate in guarded-action evaluation.
- Updated the accepted CLI and TRD schema specifications with the disposal leaf, append-only Deviation read path, Concern disposition columns, and `Concern disposed` Event family.

## Files created or modified

- `internal/review/concern.go`, `deviation.go`, and `judgment.go`, with `concern_task8_test.go`
- `internal/journal/current_event_json.go`, `current_event_result_projection.go`, `current_event_updates.go`, `current_event_validation.go`, `current_event_value_catalog.go`, `current_event_value_relations.go`, `project_effects.go`, and matching Task 8/project-effect tests
- `internal/workflow/activation.go`
- `internal/cli/codes.go` and `codes_test.go`
- `cmd/fid/concern.go`, `deviation.go`, `root.go`, `contract_surface_test.go`, and integration fixture helpers
- `cmd/fid/integration/p23_concern_deviation_test.go`
- `project-docs/specs/03_TRD_SCHEMA.md`, `project-docs/specs/06_DESIGN_PLAN_CLI.md`, and `project-docs/docs/phases/phase-23/tasks/task-8-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 3 now includes `fid concern dispose` and records the append-only `fid deviation list` semantics.
- `project-docs/specs/03_TRD_SCHEMA.md` section 3.7 now records the Concern disposition identity, reason, Actor, time, and Event columns; section 5 includes `Concern disposed` and separates the Deviation recorded family.
- `project-docs/docs/phases/phase-23/tasks/task-8-done.md` records the implementation, validation, gate history, and fence convergence evidence.

## Validation Contract Coverage

- VAL-003: PASS - `TestConcernDeviationFenceConvergencePreservesEveryRemainingEvent` in `internal/journal/project_effects_test.go:50` proves the three Task 8 keys are absent and retains exact coverage of the four remaining fence keys. `TestP23ConcernDeviationRealProcessJourney` in `cmd/fid/integration/p23_concern_deviation_test.go:17` records, disposes, rebuilds, and rereads the unfenced families through the built binary.
- VAL-034: PASS - `TestConcernSubjectAndDispositionMatrix` in `internal/review/concern_task8_test.go:42` covers all six closed subjects. The real-process journey at `cmd/fid/integration/p23_concern_deviation_test.go:17` records every type, proves the unresolved `requirement` subject succeeds, directly inspects `judgment_concerns` for no `completion_id`, leaves open blocker Concerns present, and then successfully performs both `issue.transition` to terminal `qa_pass` and `envelope.complete`.
- VAL-035: PASS - `TestConcernDisposeCarriedForwardBindsExistingIssue` and `TestConcernSubjectAndDispositionMatrix` at `internal/review/concern_task8_test.go:14` and `:42` cover the target rules. The built-binary journey exercises one success for each action, idempotent replay with no Event, and all seven required negative/no-Event paths. Exact codes are `concern-already-disposed`, `invalid-argument`, `issue-not-found`, `concern-not-found`, `finding-not-found`, `invalid-argument`, and `invalid-argument` in the contract order.
- VAL-074: PASS - `TestConcernSelfDisposalKeepsRaisingAndDisposingActors` at `internal/review/concern_task8_test.go:89` proves the two stored fields do not overwrite each other. The real-process journey proves both the equal self-disposal case and a different activated disposing Actor while the raising Actor remains unchanged in JSON readback.
- VAL-036: PASS - `TestDeviationAcceptsExplicitPredecessorReference` at `internal/review/concern_task8_test.go:101` covers the closed input. The real-process journey rejects missing and self predecessors without Events, records and replays a two-row chain, lists both rows with the exact predecessor, rebuilds, and lists the same values. `TestLiveLeafCommandsMatchContractInventories` at `cmd/fid/contract_surface_test.go:128` proves the exact command tree contains `deviation record` and `deviation list` with no update path.
- VAL-037: PASS - the real-process journey invokes `fid evidence link` once with a Concern and once with a Deviation, receives `invalid-argument` from the closed Evidence type domain, and proves neither attempt appends an Event.
- VAL-079: PASS - the exact ordinary gate reported `cmd/fid` at 538.754s and `internal/workflow` at 211.559s, both below 600s. The exact tagged integration gate completed in 651.834s, below 900s. Compared with Task 7, these are +39.059s, -0.069s, and +3.230s respectively.

## Tests run

- Required first RED: `go test ./internal/review -run '^TestConcernDisposeCarriedForwardBindsExistingIssue$' -count=1` - FAIL - `undefined: ConcernDisposeRequest` / `FAIL` (terminated `exit_code: 1`).
- Required first GREEN: the byte-identical command above - PASS - `ok  github.com/mbell-nzl/fid/internal/review 0.005s` (terminated `exit_code: 0`).
- Strict journal RED: `go test ./internal/journal -run '^(TestConcernDisposedRequiresExactFieldsAndOpenGuard|TestTask8ValueShapesReuseJudgmentSubject|TestConcernDeviationFenceConvergencePreservesEveryRemainingEvent)$' -count=1` - FAIL - `FAIL` after `current Event update table not admitted`, a missing Concern shape, and a still-fenced `concern.recorded` (terminated `exit_code: 1`).
- Strict journal GREEN: the byte-identical command above - PASS - `ok  github.com/mbell-nzl/fid/internal/journal 0.004s` (terminated `exit_code: 0`).
- Focused convergence attempt: `go test ./internal/review ./internal/journal ./internal/cli ./internal/workflow ./cmd/fid -count=1` - FAIL - `FAIL` after the error registry remained at 67 entries and the live leaf inventory omitted the four Task 8 leaves (terminated `exit_code: 1`). The accepted registries were converged in scope.
- Focused convergence GREEN: the byte-identical command above - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 43.796s` (terminated `exit_code: 0`); preceding tokens included `internal/journal 10.414s` and `internal/workflow 22.453s`.
- Focused replay RED: `go test -tags integration ./cmd/fid/integration -run '^TestP23ConcernDeviationRealProcessJourney$' -count=1 -timeout 10m` - FAIL - `concern-already-disposed` on a same-key disposal replay / `FAIL` (terminated `exit_code: 1`). Moving replay resolution ahead of terminal-state rejection fixed the product defect.
- Focused real-process GREEN after replay, predecessor, workflow-state, and multi-Actor fixture repairs: the byte-identical command above - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 3.466s` (terminated `exit_code: 0`). The journey executes the built `fid` binary for every public mutation/read path, both guarded actions, and rebuild.
- Post-refactor focused package attempt: `go test ./internal/review ./internal/journal ./internal/cli ./internal/workflow ./cmd/fid -count=1` - FAIL - `internal/review/judgment.go:7:2: "fmt" imported and not used` / `FAIL` (terminated `exit_code: 1`).
- Post-refactor byte-identical GREEN: the command above - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 44.560s` (terminated `exit_code: 0`); preceding tokens included `internal/journal 10.171s` and `internal/workflow 22.354s`.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- First `golangci-lint run ./...` attempt - FAIL - `* staticcheck: 3` (terminated `exit_code: 1`)
- Second byte-identical `golangci-lint run ./...` attempt - FAIL - `* staticcheck: 1` (terminated `exit_code: 1`)
- Final byte-identical `golangci-lint run ./...` attempt - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 211.559s coverage: 64.0% of statements` (terminated `exit_code: 0`)
- Ordinary package budgets: `cmd/fid` was 538.754s and `internal/workflow` was 211.559s; every reported ordinary package remained below 600s.
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 63.2%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 651.834s` (terminated `exit_code: 0`)
- Tagged integration budget: `cmd/fid/integration` was 651.834s, below 900s.

## Deviations from spec

- No product, schema, protected gate, threshold, validation-contract, or dependency deviation was made.
- The first exact lint attempt found three capitalized internal error strings; after those were lowercased, the byte-identical second attempt exposed the remaining capitalized Deviation self-reference error. All four strings were repaired without changing typed failure behavior, and the third byte-identical attempt passed with `0 issues.`. Both failed attempts and their real exit statuses remain recorded above.
- Focused real-process development progressed through expected sequential fixture failures: the first gate run lacked an observed Budget sample; the next used stale workflow state names; replay then exposed terminal-state ordering; the predecessor negative initially expected the wrong not-found exit; and the first multi-Actor fixture revisions widened authority grants and denied the later gate. Each cause was repaired in the same in-scope journey without changing its command, assertion, or accepted behavior; the final byte-identical command passed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
