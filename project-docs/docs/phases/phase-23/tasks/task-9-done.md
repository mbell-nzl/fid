# Phase 23 - Task 9: Decision References and protected Decisions - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 2; attempt 1 was interrupted by the Human-requested quota pause before implementation began, and attempt 2 completed the DEC-248 adapter repair and gates.
- Next action: Pause before the independent window 7-9 QA review, as directed by the Human.
- Human attention required: no

## Summary of what was implemented

- Added Decision Reference recording and listing, Source-digest verification, protected-Decision satisfaction over the approved Reference plus approved Human Outcome, the `decision.record-reference` capability, and `linked_to_decision` transition evaluation.
- Added the strict journal effects, value shapes, JSON registrations, retained-family rules, Stop details, Condition discriminator shape, archive semantics, and replay paths needed for Decision References and protected Decisions.
- Preserved all six Phase 22 Condition encodings and nullable relation bindings through the new per-type dispatch, including pre-change journal replay and Attention-key stability checks.
- Applied DEC-248’s bounded CLI adapter repair: `envelope.complete` and terminal Issue transition render the immutable `protected-decision-required` failure at exit 5 before an unsatisfied Obligation’s generic denial.
- Converged Task 9-invalidated retained fixtures: the claims fake projection includes the two current Decision tables, and the V01 closure journey reads the new `fid decision list` leaf without appending a journal Event.

## Files created or modified

- `internal/review/decision.go` and Decision Reference tests.
- `internal/workflow/` Decision, evaluation, activation, Stop, registry, rule, and authority paths with matching tests.
- `internal/journal/` retained-family effects, typed result/value grammar, JSON registration, replay fixture, and Condition-dispatch tests.
- `internal/issues/transition_rules.go` and its Decision-link test.
- `internal/export/archive_semantics_catalog.go` and `internal/cli/codes.go` with tests.
- `cmd/fid/decision.go`, `envelope.go`, `issue_mutation_response.go`, command-surface tests, and Task 9 integration journeys.
- `internal/claims/generic_test.go` and `cmd/fid/integration/e2e_v01_closure_helpers_test.go` retained fixtures converged by the Task 9 command/schema additions.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` and this completion report.

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` documents `fid decision record-reference`, `fid decision list`, and protected-Decision completion behavior.
- `project-docs/docs/phases/phase-23/tasks/task-9-done.md` records validation, recovery, gate, and commit evidence.

## Validation Contract Coverage

- VAL-003: PASS - `internal/journal/project_effects_test.go:50` (`TestConcernDeviationFenceConvergencePreservesEveryRemainingEvent`) retains the progressive-fence subject while Task 9 deletes only `decision.referenced`; `cmd/fid/integration/p23_decision_reference_test.go:16` records, rebuilds, and rereads the unfenced family through the built binary.
- VAL-038: PASS - `internal/review/decision_test.go:13` (`TestDecisionReferenceSourceMustExistAndMatchDigest`) proves matching Source recording and both missing/mismatched digest rejection. The real-process journey at `cmd/fid/integration/p23_decision_reference_test.go:16` asserts `source-digest-mismatch` and unchanged Event count for both negatives.
- VAL-039: PASS - `internal/workflow/protected_decision_test.go:11` proves only declared actions are guarded. The real-process journey at `cmd/fid/integration/p23_decision_reference_test.go:16` observes both CLI adapters return exit 5, `protected-decision-required`, the exact four immutable Decision details, and no Event; its terminal Issue transition carries simultaneous protected-Decision and Obligation reasons before DEC-248’s adapter precedence.
- VAL-040: PASS - `internal/workflow/protected_decision_test.go:45` covers neither/reference-only/outcome-only/both states and exact capability matching. The built-binary journey records rejected, deferred, wrong-authority, reference-only, and both-approved states and checks the corresponding Stop remains open or resolves.
- VAL-041: PASS - `cmd/fid/integration/p23_decision_reference_test.go:16` starts with `release` and `security` Stops, verifies both are open, then records each satisfying pair and verifies exactly the matching Stop and Attention resolve.
- VAL-042: PASS - `internal/issues/transition_rules_task9_test.go:11` (`TestLinkedToDecisionRequiresApprovedReferenceInIssueScope`) proves the implemented rule rejects before and accepts after a qualifying scoped Reference.
- VAL-080: PASS (Task 9 contribution) - `internal/journal/current_event_task9_values_test.go:5` rejects incomplete, extra, wrong discriminator, and wrong-subject protected-Decision Stop values while accepting only the `decision_class` Condition row; Task 11 completes the four-type matrix.
- VAL-081: PASS - `internal/journal/current_event_closed_values_test.go:140` (`TestTask9ConditionDispatchPreservesNullableBindings`) proves the new dispatch still reaches the nullable binding relation for retained and Decision Conditions.
- VAL-082: PASS - `internal/journal/current_event_closed_values_test.go:115` (`TestTask9ConditionDispatchRetainsPinnedPhase22Bytes`) accepts canonical pre-change bytes and rejects mutations for all six retained Condition types.
- VAL-083: PASS - `internal/journal/task9_prechange_replay_test.go:21` (`TestTask9PrechangeBudgetObservedJournalReplaysWithoutMutation`) replays the named pre-change compressed journal and compares the source digest before and after.
- VAL-084: PASS - `internal/workflow/protected_decision_test.go:118` (`TestTask9RetainedAttentionKeyMatchesPrechangeReference`) compares the fixed-input retained Attention key with reference `74c83b7033729785eb99403ecec9dcde30ac4280`.
- VAL-079: PASS - the serial ordinary gate reports `cmd/fid` at 513.222s and the earlier observed `internal/workflow` package time was 431.647s, both below 600s. The final tagged integration gate was 713.194s, below 900s.

## Tests run

- Required DEC-248 RED: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/... -run '^TestP23DecisionReferenceProtectedLifecycleRealProcess$'` - FAIL - `FAIL` after the inner CLI returned exit 4 `transition-not-permitted` with `protected-decision-required` and `obligation-evidence-absent` (terminated `exit_code: 1`).
- Required DEC-248 GREEN: the byte-identical command above - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 6.439s` (terminated `exit_code: 0`). It executes the built `fid` binary through both protected Decision denials.
- Claims fixture RED: `go test -race -timeout 30m ./internal/claims -run '^TestClaimAcquireRejectsOpenReviewReworkStop$' -count=1` - FAIL - `FAIL` after `no such table: envelope_protected_decisions` (terminated `exit_code: 1`).
- Claims fixture GREEN: the byte-identical command above - PASS - `ok  github.com/mbell-nzl/fid/internal/claims 1.038s` (terminated `exit_code: 0`).
- V01 Decision-leaf RED: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/... -run '^TestV01ClosureJourneyAndRebuildParity$'` - FAIL - `FAIL` because the retained fixture expected the new `decision` command to be absent (terminated `exit_code: 1`).
- V01 Decision-leaf GREEN: the byte-identical command above - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 6.315s` (terminated `exit_code: 0`).
- Optional broad focused command: `go test -race -timeout 30m ./cmd/fid ./internal/review ./internal/workflow ./internal/issues ./internal/journal` - PASS - `ok  github.com/mbell-nzl/fid/internal/journal 191.336s` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- First `golangci-lint run ./...` attempt - FAIL - `* staticcheck: 3` (terminated `exit_code: 1`).
- Second byte-identical `golangci-lint run ./...` attempt - FAIL - `* staticcheck: 1` (terminated `exit_code: 1`).
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`).
- First `go test -race -timeout 30m -coverprofile=coverage.out ./...` attempt - FAIL - `FAIL` (terminated `exit_code: 1`): Task 9’s new Stop ordering query exposed the outdated claims fake schema; the concurrently loaded `TestLiveOutputGoldenFixtures` also failed, but its isolated race rerun passed and no causal conclusion is claimed.
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 63.0% of statements` (terminated `exit_code: 0`).
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 62.5%` (terminated `exit_code: 0`).
- First `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` attempt - FAIL - `FAIL` (terminated `exit_code: 1`): the retained V01 closure fixture asserted that Task 9’s new Decision command was unknown.
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 713.194s` (terminated `exit_code: 0`). This final invocation followed `gofmt -w cmd/fid/integration/e2e_v01_closure_helpers_test.go` as required formatting preparation; the required Go command itself is byte-identical and was the final shell command.

## Deviations from spec

- No product, public-interface, validation-contract, threshold, dependency, or protected gate definition changed.
- DEC-248 added the two missing CLI adapter files and matching tests to Task 9 scope. The domain admission rules and immutable Stop details remain unchanged; only the accepted CLI mapping and precedence were completed.
- The initial lint failures were Task 9’s new Condition-dispatch error literals. The repair parameterizes the existing `Condition` kind so the accepted emitted bytes remain unchanged while satisfying staticcheck.
- The first ordinary gate ran while a root verification invocation overlapped. Its `cmd/fid` live-output failure is preserved as failed evidence; isolated retry and the serial ordinary rerun passed, but this does not establish a causal explanation for the intermittent observation.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
