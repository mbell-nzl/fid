## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: gate-repair
- Condition: Task 12 has converged the retained V1 Claim observation to the accepted standalone/execution-bound boundary, but its tagged integration gate now stops at Task 9's `issue ac verify` transport and repository coverage stops at the out-of-scope `internal/dev` V1 Claim expectation.
- Recovery attempts: 1
- Next action: Phase Runner must route the tagged AC-verification fixture and `internal/dev/start_test.go` to their owning accepted scopes, then dispatch the same Task 12 continuation to rerun the six exact commands.
- Human attention required: no

## Summary of what was implemented

- Added atomic invalid-input validation to Issue revision Effect construction.
- Added the execution-bound nonterminal transition predicate for Envelope, trusted Actor authority, generic Claim/Run, repository, Issue revision, and workflow edge qualification.
- Added the shared fail-closed Obligation evaluator. Until Group 2 Evidence exists, terminal-success reports `obligation-evidence-absent` and commits no Event.
- Replaced the retained non-bound Issue transition's V1 Claim precondition with the accepted boundary: standalone structural transitions use attribution only, while Envelope-bound transitions continue to read generic `work_claims` through the execution evaluator.
- Converged only the directly invalidated V1 Claim/role closure fixtures and retained no-mutation/rebuild-parity assertions outside that obsolete precondition.

## Files created or modified

- `cmd/fid/issue_mutation_response.go`
- `cmd/fid/task12_execution_transition_test.go`
- `internal/issues/revision.go`
- `internal/issues/revision_test.go`
- `internal/issues/transition.go`
- `internal/issues/transition_test.go`
- `cmd/fid/issue_transition_test.go`
- `cmd/fid/integration/e2e_v01_closure_test.go`
- `cmd/fid/e2e_v01d_closure_test.go`
- `internal/workflow/evaluation.go`
- `internal/workflow/evaluation_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-12-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-12-done.md`

## Validation Contract Coverage

- VAL-038: PASS - `TestExecutionBoundNonterminalTransitionRequiresActiveEnvelopeClaimRunAndRevision` at `cmd/fid/task12_execution_transition_test.go:9` advances one all-qualified nonterminal transition through the CLI root command.
- VAL-039: PASS - `TestExecutionBoundTransitionReturnsEveryOrderedReasonWithoutMutation` at `cmd/fid/task12_execution_transition_test.go:56` proves ordered all-reasons output and an unchanged Event count.
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

## Exact quality-gate results

- `go test -race ./internal/workflow ./internal/issues ./internal/policy ./internal/validation ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t(cached)` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`): retained fixtures use removed `--by` flags; provenance is `e712f85d08a0b3259b23badd34d9306bf0a10f9b`, `7119d5b38d79663f544f093d55f676c375746419`, and `09aa0a081d2ecfee025c3b4d8158bc6203d02d98`, all before Phase 21.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t24.081s\tcoverage: 33.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `65.3%` (terminated `exit_code: 0`)

### Continuation attempt 1

- `go test -race ./internal/workflow ./internal/issues ./internal/policy ./internal/validation ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t171.603s` (terminated `exit_code: 0`). The first post-repair package attempt failed at retained V1 Claim-precondition tests; those exact Task 12 Claim-observation tests were converged before this passing rerun.
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`): `TestV01ClosureJourneyAndRebuildParity` now reaches `issue ac verify ... --status pass --by qa_lead`, which exits 1 with `error: unknown flag: --status`. This is Task 9's accepted fail-closed AC verification surface, not Claim-observation collateral.
- `gofmt -l cmd internal` - FAIL - `internal/issues/transition_test.go` (terminated `exit_code: 0`): the command terminated successfully while reporting an unformatted edited file. After `gofmt -w internal/issues/transition_test.go`, the byte-identical rerun PASSed with `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` - FAIL - `1 issues:` (terminated `exit_code: 1`): `entersStartedCategory` became unused after removal of the V1 Claim fallback. Deleting that obsolete helper produced the byte-identical rerun PASS token `0 issues.` (terminated `exit_code: 0`).
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`): `internal/dev/TestStartRejectsInvalidWorkflowEdgeBeforeAppend` still expects the removed V1 Claim precondition. `internal/dev` is outside Task 12 scope and its Dev assertion is expressly forbidden from this repair.
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t\t65.3%` (terminated `exit_code: 0`).

## Deviations from spec

- Required tagged integration attempt before repair: `FAIL` (terminated `exit_code: 1`) at the retained V1 Claim observation. The Task 12 repair removed only the non-bound V1 Claim read and converged its direct retained Claim assertions; it preserves the Envelope-bound generic `work_claims` evaluator and all VAL-038..041 behavior.
- The first post-repair format and lint failures are retained above. Their root causes were an unformatted edited test and an obsolete helper left after the bounded deletion; neither changed a command, assertion semantics, threshold, or public contract.
- The remaining tagged and coverage failures are not repaired here: the tagged `issue ac verify` failure belongs to Task 9's accepted fail-closed AC surface, and the coverage `internal/dev` failure is forbidden Dev collateral outside Task 12 scope.

## Concerns, technical debt, or notes for the Tech Lead

- Route `cmd/fid/integration/e2e_v01_closure_test.go` AC verification to Task 9 and `internal/dev/start_test.go` to the Dev-owning accepted scope before rerunning Task 12's exact gate.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `c407c305f404c2b418b021c415323972d5044ebb`

## Commit subject

- `fix(issues): converge standalone transition claims`
