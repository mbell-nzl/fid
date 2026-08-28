## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: DEC-167 attempt 6 converged the stale retained Dev precondition and every accepted-Phase-21-invalidated tagged integration row. All six unchanged required commands pass.
- Recovery attempts: 6
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added atomic invalid-input validation to Issue revision Effect construction.
- Added the execution-bound nonterminal transition predicate for Envelope, trusted Actor authority, generic Claim/Run, repository, Issue revision, and workflow edge qualification.
- Added the shared fail-closed Obligation evaluator. Until Group 2 Evidence exists, terminal-success reports `obligation-evidence-absent` and commits no Event.
- Replaced the retained non-bound Issue transition's V1 Claim precondition with the accepted boundary: standalone structural transitions use attribution only, while Envelope-bound transitions continue to read generic `work_claims` through the execution evaluator.
- Enforced the Envelope repository worktree match when resolving an active generic Claim for an execution-bound transition.
- Restored the retained closure coverage as accepted-contract equivalents: standalone attribution commits without a Claim or role-derived authority, rejected workflow edges leave the journal unchanged, and Claim/rebuild parity remains asserted.
- Restored the removed Claim-holder and Claim-expiry transition obligations at their accepted generic boundary: an Envelope-bound transition reports the Claim-specific reason plus `run-inactive` and appends no Event when its Claim belongs to another Actor or is expired.
- Converged the tagged closure journey's unavailable V1 AC and Concern mutations into typed `capability-unavailable`, no-journal-append assertions; its removed V1 `run record` leaf now rejects without a journal append and maintenance asserts the resulting telemetry-gap warning.
- Removed only the stale standalone-transition Claim-rejection probe from the Dev invalid-edge setup. The `Start()` subject, typed invalid-transition assertion, current-state assertion, and journal-unchanged assertion remain unchanged.

### DEC-167 retained-row provenance

| Retained row | Current assertion | Invalidating accepted change | Task |
|---|---|---|---|
| `internal/dev/start_test.go:225` | `TestStartRejectsInvalidWorkflowEdgeBeforeAppend` retains the post-Claim setup and all `Start()` assertions. | Non-bound Issue transitions no longer read the V1 Claim precondition. | task 12, `c407c305f404c2b418b021c415323972d5044ebb` |
| `cmd/fid/integration/e2e_v01_closure_test.go:38` | AC verification is a typed `capability-unavailable` rejection with no journal append. | VAL-030 makes `issue ac verify` unconditionally unavailable. | task 9, `cmd/fid/issue_ac.go:109-111` |
| `cmd/fid/integration/e2e_v01_closure_helpers_test.go:260` | Concern raise/address are typed `capability-unavailable` rejections with no journal append. | VAL-030 makes Concern mutation unavailable. | task 9, `cmd/fid/concern.go:83,101` |
| `cmd/fid/integration/e2e_v01_closure_test.go:61` | Legacy `run record` rejects with empty output and no journal append. | The generic Run lifecycle removed the `run record` leaf. | task 11, `5ca732ff163ee1d31fc155d3d9ebb3dac1501838`; surviving equivalent `cmd/fid/run_test.go:5` |
| `cmd/fid/integration/e2e_v01_closure_helpers_test.go:350` | Maintenance reports the deterministic `telemetry_gaps` warning without a journal append. | The removed `run record` writer leaves completion telemetry unavailable. | task 11, `5ca732ff163ee1d31fc155d3d9ebb3dac1501838` |

The retained Task 8 Event-actor and `issue create --by` assertions remain in `TestP20StandaloneIssueJourneyJSONAndRebuildParity`, `TestV01ClosureJourneyAndRebuildParity`, and `TestV02aWorkflowAuthoringExitScenarioAndRebuildParity` as named by `task-9-done.md:122`.

## Files created or modified

- `cmd/fid/issue_mutation_response.go`
- `cmd/fid/task12_execution_transition_test.go`
- `internal/issues/revision.go`
- `internal/issues/revision_test.go`
- `internal/issues/transition.go`
- `internal/issues/transition_test.go`
- `cmd/fid/issue_transition_test.go`
- `cmd/fid/integration/e2e_v01_closure_test.go`
- `cmd/fid/integration/e2e_v01_closure_helpers_test.go`
- `cmd/fid/e2e_v01d_closure_test.go`
- `internal/dev/start_test.go`
- `internal/workflow/evaluation.go`
- `internal/workflow/evaluation_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-12-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-12-done.md`

## Validation Contract Coverage

- VAL-038: PASS - `TestExecutionBoundNonterminalTransitionRequiresActiveEnvelopeClaimRunAndRevision` at `cmd/fid/task12_execution_transition_test.go:11` advances one all-qualified nonterminal transition through the CLI root command; `TestExecutionBoundTransitionRejectsMismatchedClaimWorktreeWithoutMutation` at `cmd/fid/task12_execution_transition_test.go:78` rejects a Claim whose worktree differs from the Envelope scope.
- VAL-039: PASS - `TestExecutionBoundTransitionReturnsEveryOrderedReasonWithoutMutation` at `cmd/fid/task12_execution_transition_test.go:58` proves ordered all-reasons output and an unchanged Event count; `TestExecutionBoundTransitionRejectsMismatchedClaimWorktreeWithoutMutation` at `cmd/fid/task12_execution_transition_test.go:78` proves `claim-missing,run-inactive` and no Event for a mismatched worktree.
- VAL-040: PASS - `TestEvaluateObligationsReturnsIdenticalUnsatisfiedReasonsForBothGuards` at `internal/workflow/evaluation_test.go:8` compares the terminal-success and `envelope.complete` evaluator outputs byte-for-byte; `cmd/fid/task12_execution_transition_test.go:40` proves terminal transition remains atomic before Evidence exists.
- VAL-041: PASS - `TestBuildIssueRevisionEffectsRejectsInvalidInputsWithoutEffects` at `internal/issues/revision_test.go:42` covers whitespace IDs, normalized duplicates, and overflow with nil Effects.

## Tests run

- `go test ./internal/issues --run TestBuildIssueRevisionEffectsRejectsInvalidInputsWithoutEffects -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/issues\t0.004s` (terminated `exit_code: 0`)
- `go test ./cmd/fid --run TestExecutionBoundNonterminalTransitionRequiresActiveEnvelopeClaimRunAndRevision -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.313s` (terminated `exit_code: 0`)
- `go test ./cmd/fid --run TestExecutionBoundTransitionReturnsEveryOrderedReasonWithoutMutation -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.208s` (terminated `exit_code: 0`)
- `go test ./internal/workflow --run TestEvaluateObligationsReturnsIdenticalUnsatisfiedReasonsForBothGuards -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.003s` (terminated `exit_code: 0`)
- `go test ./internal/issues ./internal/workflow ./cmd/fid --run 'Test(BuildIssueRevisionEffectsRejectsInvalidInputsWithoutEffects|ExecutionBound)' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.452s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestIssueTransitionJSONAllowsStandaloneActorWithoutGenericClaim$' -count=1` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`): the retained V1 `live_claim_held` precondition rejected the standalone transition.
- `go test -race ./cmd/fid -run '^TestIssueTransitionJSONAllowsStandaloneActorWithoutGenericClaim$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t1.942s` (terminated `exit_code: 0`)
- `go test -race ./internal/issues ./cmd/fid -run '^(TestTransitionIntoStartedAllowsStandaloneActorWithoutGenericClaim|TestTransitionIntoStartedIgnoresRetainedV1Claims|TestIssueTransitionJSONAllowsStandaloneActorWithoutGenericClaim|TestIssueTransitionIgnoresRetainedV1ClaimsForStandaloneStates|TestV01DClosureLifecycleSmokeAndRebuildParity)$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t5.281s` (terminated `exit_code: 0`)
- Real entrypoint: `go build -o /tmp/fid-task12 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`); `/tmp/fid-task12 --project /tmp/fid-task12-entry.tDHf2P issue transition FID-7770e6 --to in_progress --actor agent:developer --idempotency-key task12-entry-transition` - PASS - `transitioned FID-7770e6: todo → in_progress` (terminated `exit_code: 0`).
- `go test ./cmd/fid --run '^TestExecutionBoundTransitionRejectsMismatchedClaimWorktreeWithoutMutation$' -count=1` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`): a Claim mutated to `other-worktree` still transitioned the Envelope-bound Issue and appended an Event.
- `go test ./cmd/fid --run '^TestExecutionBoundTransitionRejectsMismatchedClaimWorktreeWithoutMutation$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.210s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid --run '^(TestV01DClosureLifecycleSmokeAndRebuildParity|TestExecutionBoundTransitionRejectsMismatchedClaimWorktreeWithoutMutation)$' -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t6.235s` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestV01ClosureJourneyAndRebuildParity$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`): the restored closure journey reaches Task 9's `issue ac verify ... --status pass` and receives `error: unknown flag: --status`.
- Historical expected-red tagged integration evidence: `go test -tags integration -race ./cmd/fid/integration/...` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`): the retained AC success call passed removed `--status` and stopped before the VAL-030 unavailable leaf.
- Historical expected-red tagged integration evidence: `go test -tags integration -race ./cmd/fid/integration/...` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`): after AC and Concern convergence, the journey reached removed `run record`.
- Historical expected-red tagged integration evidence: `go test -tags integration -race ./cmd/fid/integration/...` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`): after the `run record` rejection, maintenance exposed the expected `telemetry_gaps` warning rather than the stale `pass` status.
- Historical expected-red coverage evidence: `go test -race -coverprofile=coverage.out ./...` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`): the initial Dev convergence retained two `in_progress` transitions; the second rejected as an invalid edge.
- `go test -race ./internal/dev` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/dev\t10.125s` (terminated `exit_code: 0`)
- Real entrypoint: `go build -o /tmp/fid-task12-attempt6 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`); `/tmp/fid-task12-attempt6 --json issue ac verify FID-1 AC-1 --actor agent:qa_lead` in `/tmp/fid-task12-attempt6.dHeJ9O` - PASS (expected capability-unavailable rejection) - `{"contract_version":"fid-cli.v1","error":{"category":"authority","code":"capability-unavailable","details":{"available":false,"capability_id":"reviews.coverage.v1","owner_phase":23},"message":"capability is unavailable","retryable":false,"subject":{"id":"reviews.coverage.v1","type":"capability"}},"ok":false,"operation":"issue.ac.verify"}` (terminated `exit_code: 1`).
- Real entrypoint: `go build -o /tmp/fid-task12-attempt2 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`); `/tmp/fid-task12-attempt2 issue transition FID-70fcfa --to in_progress --actor agent:qa_reviewer --idempotency-key task12-entry-transition` in `/tmp/fid-task12-entry.y8p4B8` - PASS - `transitioned FID-70fcfa: todo → in_progress` (terminated `exit_code: 0`).
- `go test ./cmd/fid` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t14.266s` (terminated `exit_code: 0`)
- `TestExecutionBoundTransitionRetainsClaimHolderAndExpiryNoMutationGuards` preserves wrong-holder and expired generic Claim denials with no journal append.

## Exact quality-gate results

- Attempt 6 final exact task-gate rerun:
- `go test -race ./internal/workflow ./internal/issues ./internal/policy ./internal/validation ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t175.213s` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t5.597s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 33.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t\t65.3%` (terminated `exit_code: 0`)

- Attempt 6 first coverage attempt:
- Historical Attempt 6 first coverage evidence: `go test -race -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`): `TestStartRejectsInvalidWorkflowEdgeBeforeAppend` exposed the duplicate setup transition; the stale V1 Claim-rejection probe was removed, leaving the original post-Claim transition as the only setup transition.

- Attempt 5 (this continuation):
- `go test -race ./internal/workflow ./internal/issues ./internal/policy ./internal/validation ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t173.260s` (terminated `exit_code: 0`)
- Historical Attempt 5 tagged integration evidence: `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`): `TestV01ClosureJourneyAndRebuildParity` reaches Task 9's `issue ac verify ... --status pass --by qa_lead`, which exits 1 with `error: unknown flag: --status`; Task 12 must preserve this external failure.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- Historical Attempt 5 coverage evidence: `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`): `internal/dev/TestStartRejectsInvalidWorkflowEdgeBeforeAppend` still expects the removed V1 Claim precondition; later Task 20 owns that fixture and Task 12 must not repair it.
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.3%` (terminated `exit_code: 0`)

- Attempt 3 (this continuation):
- `go test -race ./internal/workflow ./internal/issues ./internal/policy ./internal/validation ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t172.951s` (terminated `exit_code: 0`)
- Historical Attempt 3 tagged integration evidence: `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`): `TestV01ClosureJourneyAndRebuildParity` reaches Task 9's `issue ac verify ... --status pass --by qa_lead`, which exits 1 with `error: unknown flag: --status`; Task 12 must preserve this external failure.
- `gofmt -l cmd internal` - FAIL - `cmd/fid/task12_execution_transition_test.go` (terminated `exit_code: 0`): the new assertion required formatting.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- Historical Attempt 3 coverage evidence: `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`): `internal/dev/TestStartRejectsInvalidWorkflowEdgeBeforeAppend` still expects the removed V1 Claim precondition; later Task 20 owns that fixture and Task 12 must not repair it.
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.3%` (terminated `exit_code: 0`)

- Historical continuation attempt 2:

- `go test -race ./internal/workflow ./internal/issues ./internal/policy ./internal/validation ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t171.975s` (terminated `exit_code: 0`)
- Historical continuation attempt 2 tagged integration evidence: `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`): `TestV01ClosureJourneyAndRebuildParity` reaches Task 9's `issue ac verify ... --status pass --by qa_lead`, which exits 1 with `error: unknown flag: --status`; Task 12 must preserve this external failure.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- Historical continuation attempt 2 coverage evidence: `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`): `internal/dev/TestStartRejectsInvalidWorkflowEdgeBeforeAppend` still expects the removed V1 Claim precondition; later Task 20 owns that fixture and Task 12 must not repair it.
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.3%` (terminated `exit_code: 0`)

- Historical continuation attempt 1:

- `go test -race ./internal/workflow ./internal/issues ./internal/policy ./internal/validation ./internal/cli ./cmd/fid/...` - FAIL - `FAIL` (terminated `exit_code: 1`): retained V1 Claim-precondition tests still expected the obsolete standalone Claim gate.
- `go test -race ./internal/workflow ./internal/issues ./internal/policy ./internal/validation ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t171.603s` (terminated `exit_code: 0`)
- Historical continuation attempt 1 tagged integration evidence: `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`): `TestV01ClosureJourneyAndRebuildParity` now reaches `issue ac verify ... --status pass --by qa_lead`, which exits 1 with `error: unknown flag: --status`. This is Task 9's accepted fail-closed AC verification surface, not Claim-observation collateral.
- `gofmt -l cmd internal` - FAIL - `internal/issues/transition_test.go` (terminated `exit_code: 0`): the command reported an unformatted edited file.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - FAIL - `1 issues:` (terminated `exit_code: 1`): `entersStartedCategory` became unused after removal of the V1 Claim fallback.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- Historical continuation attempt 1 coverage evidence: `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`): `internal/dev/TestStartRejectsInvalidWorkflowEdgeBeforeAppend` still expects the removed V1 Claim precondition. `internal/dev` is outside Task 12 scope and its Dev assertion is expressly forbidden from this repair.
- `go tool cover -func=coverage.out` - PASS - `65.3%` (terminated `exit_code: 0`)

## Deviations from spec

- Required tagged integration attempt before repair: `FAIL` (terminated `exit_code: 1`) at the retained V1 Claim observation. The Task 12 repair removed only the non-bound V1 Claim read and converged its direct retained Claim assertions; it preserves the Envelope-bound generic `work_claims` evaluator and all VAL-038..041 behavior.
- The first post-repair format and lint failures are retained above. `gofmt -w internal/issues/transition_test.go` corrected the formatting, and deletion of the obsolete `entersStartedCategory` helper corrected lint; neither changed a command, assertion semantics, threshold, or public contract.
- The remaining tagged and coverage failures are not repaired here: the tagged `issue ac verify` failure belongs to Task 9's accepted fail-closed AC surface, and the coverage `internal/dev` failure is forbidden Dev collateral outside Task 12 scope.
- Attempt 3 corrected the prior assertion-boundary defect without restoring V1 `issue claim|unclaim` aliases or `--by`: the direct generic evaluator test retains the historic Claim-holder/expiry and atomic no-mutation obligations at the only accepted execution-bound surface. Task 9 owns the AC fixture transport; Task 20 owns the retained Dev fixture.
- Attempt 3's focused coverage restoration is committed at `2a584fc125ba16f90973951dccda731c04c357a7`; after recording the six exact command terminations above, its child session was shut down under MISSING_TERMINATION before a terminal response. This report-only continuation preserves that partial-work evidence and does not rerun or reclassify either remaining owner-routed failure.
- Attempt 5 made no implementation or assertion change. It preserves the immutable Task 9 and Task 20 gate failures above, with their observed `exit_code: 1` terminations, pending their separately authorized owner repairs.
- Attempt 6 initially retained two transitions while removing the stale Dev probe, so the coverage gate correctly rejected an `in_progress` to `in_progress` setup edge. The repair removed that duplicate setup only; `Start()`'s `*issues.InvalidTransitionError`, `CurrentState`, and no-journal-append assertions are unchanged. The later exact coverage gate passed.
- DEC-167 supersedes the prior owner-routing conclusions for the tagged integration and Dev rows. This report preserves their failed evidence but records the accepted-change provenance above instead.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `16b15d8134e7b29f27cf996486bff997a080fcc3`

## Commit subject

- `test(integration): converge retained closure journey`
