# Phase 21 - Task 9: Retained actor-policy removal - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: DEC-166 resolved TRG-21-010's stale ownership classification and the final task-9 gate remains green. The tagged integration source then exposed retained task-9 `--by` transports; this continuation converged those calls while preserving later-owner Issue/Event failures reached afterward.
- Recovery attempts: 12 prior task-9 recoveries; integration-fixture convergence attempt 1
- Next action: Phase Runner verifies this committed report with the two-argument report lint and routes the remaining tagged integration failures to their owning accepted scope before task-12 continuation.
- Human attention required: no

## Summary of what was implemented

- Removed public registration of the V1 `fid policy` command while retaining its unregistered helper for enumerated legacy dev/completion callers.
- Added standalone project-local actor attribution and Event-v2 idempotency for workflow definition, workflow assignment, and acceptance-criterion addition.
- Made AC verification, concern judgment mutations, and maintenance vacuum return typed capability-unavailable failures; renamed the public stats grouping to `actor` and labelled it as historical attribution.
- Added a focused CLI conformance matrix for the converted public surfaces.
- Converted label add, issue-label add/remove, link add/remove, and requirement sync to standalone actor/key Event-v2 paths; started canonical activity event normalization and transport migration.
- Registered `concern list` as an accepted project/JSON root-flag leaf and strengthened stats CLI assertions to the public `actor_attribution` envelope with `--group-by actor|epic`.
- Added an allowlisted retained-mutation fixture adapter for standalone actor/idempotency transport and made replacement-bound dev completion tests prove public AC verification is unavailable before seeding historical projection state.
- Replaced public policy-matrix fixture expectations with absent-command and no-side-effect checks; updated vacuum and converted-read root-flag fixtures to typed unavailable and accepted project/JSON behavior.
- Renamed the Phase 20 task-labelled request fixture to `retained_actor_request_contract_test.go` and converged its Task-9-owned attribution rows to canonical invalid-argument JSON, concise Human `--actor`/Actor-grammar errors, and selected-project actor/idempotency requests.
- QA Review Window 3 repair: removed the retired public `role` grouping normalization and added a CLI stats conformance matrix for explicit actor, epic, and model envelopes plus rejection/no-journal-mutation checks for role, another unknown group, and retired `--by`.
- Converged the retained tagged integration fixture transport through one allowlisted adapter: workflow define/assign, AC add, labels, links, and requirement sync now use `--actor agent:<role>` with unique idempotency keys. Issue creation, Claim, Run, AC verification, concern, dev, completion, and closure behaviors remain untouched.

## Files created or modified

- `cmd/fid/concern.go`
- `cmd/fid/capabilities_test.go`
- `cmd/fid/contract_surface_test.go`
- `cmd/fid/issue_ac.go`
- `cmd/fid/maintenance.go`
- `cmd/fid/maintenance_test.go`
- `cmd/fid/maintenance_vacuum_test.go`
- `cmd/fid/policy.go`
- `cmd/fid/policy_test.go`
- `cmd/fid/label_test.go`
- `cmd/fid/link_test.go`
- `cmd/fid/issue_label_test.go`
- `cmd/fid/requirement_test.go`
- `cmd/fid/issue_claim_test.go`
- `cmd/fid/run_test.go`
- `cmd/fid/dev_test.go`
- `cmd/fid/integration/e2e_v01_closure_helpers_test.go`
- `cmd/fid/root.go`
- `cmd/fid/root_flags_test.go`
- `cmd/fid/stats.go`
- `cmd/fid/stats_test.go`
- `cmd/fid/task9_contract_test.go`
- `cmd/fid/unavailable.go`
- `cmd/fid/workflow.go`
- `cmd/fid/workspace_fixture_test.go`
- `internal/cli/codes.go`
- `internal/issues/ac.go`
- `internal/issues/idempotency_results.go`
- `internal/stats/stats.go`
- `internal/workflow/assign.go`
- `internal/workflow/attribution.go`
- `internal/workflow/define.go`
- `internal/activity/activity.go`
- `internal/labels/attribution.go`
- `internal/labels/add.go`
- `internal/labels/types.go`
- `internal/links/add.go`
- `internal/links/attribution.go`
- `internal/links/idempotency_results.go`
- `internal/links/remove.go`
- `internal/links/types.go`
- `internal/requirements/attribution.go`
- `internal/requirements/sync.go`
- `internal/requirements/types.go`
- `internal/dev/finish_test.go`
- `cmd/fid/retained_actor_request_contract_test.go` (renamed from `cmd/fid/task10_contract_test.go`)

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-9-done.md`

## Validation Contract Coverage

- VAL-028: PASS - `TestTask9RetainedInventorySurface` (`cmd/fid/task9_contract_test.go:10`) and `TestRetainedInventoryJSONMissingActorsUseCanonicalEnvelope` (`cmd/fid/retained_actor_request_contract_test.go:79`) prove retained inventory operations reject missing/malformed attribution before journal mutation; `TestRetainedInventoryOperationsUseSelectedProjectOutsideWorkingDirectory` (`cmd/fid/retained_actor_request_contract_test.go:173`) proves selected-project actor/idempotency success.
- VAL-029: PASS - `TestPublicPolicyCommandIsAbsent` (`cmd/fid/policy_test.go:15`) and `TestPolicyShowIsReadOnlyAndDoesNotAcceptByFlag` (`cmd/fid/policy_test.go:45`) pass in the final required package gate without restoring public policy authority.
- VAL-030: PASS - `TestDeferredCapabilitySurfacesStayUnavailableWithoutSideEffects` (`cmd/fid/capabilities_test.go:102`) passes in the final required package gate with typed unavailable results and no side effects.
- VAL-031: PASS - `TestStatsTokensAcceptsOnlyClosedGroupByValues` (`cmd/fid/stats_test.go:37`) exercises exact successful `--group-by actor`, `epic`, and `model` JSON envelopes, rejects retired `role`, another invalid value, and retired `--by`, and compares journal bytes after every command. `validateStatsOptions` (`internal/stats/stats.go:132`) now accepts only actor, epic, or model. The final required package gate passes with this test.
- Integration-fixture convergence: `TestTask9IntegrationMutationRequestsUseActorAndIdempotency` (`cmd/fid/integration/e2e_v01_closure_helpers_test.go`) proves all seven retained task-9 families use actor/key transport while later-owner Issue creation, Claim, Run, AC verification, and concern calls stay byte-identical. The final tagged run reaches only task-8 Event/Issue assertions after these task-9 transports pass.

- Attempt 2 update: VAL-030 now has direct no-writer evidence in `internal/issues/judgment_unavailable_test.go` and `internal/maintenance/vacuum_unavailable_test.go`; the remaining CLI persistence-negative matrix is in scope for recovery. VAL-028, VAL-029, and VAL-031 remain PARTIAL because the required gate still fails on retained V1 fixture expectations.

- Attempt 4 update: VAL-028 internal workflow, links, requirements, and activity assertions now pass as Event-v2/actor/key and normalized-read coverage. VAL-030 continues to retain typed no-writer assertions; the incomplete CLI migration still prevents COMPLETE for VAL-028 through VAL-031.

- Attempt 5 update: VAL-028 now includes the exact retained root inventory in `TestLiveLeafCommandsMatchContractInventories` (`cmd/fid/contract_surface_test.go`) and an accepted `concern list` root capability. VAL-029 and VAL-030 retain the policy-absence and capability-unavailable coverage in `TestDeferredCapabilitySurfacesStayUnavailableWithoutSideEffects` (`cmd/fid/capabilities_test.go`). VAL-031 CLI assertions now use `TestStatsTokensDefaultsToJSONByActorAttribution`, `TestStatsTokensUsesDirectLegacyJSONAndEpicFilter`, and `TestStatsCostUsesDirectLegacyJSONAndExcludesNullCostsFromSum` (`cmd/fid/stats_test.go`) to check public envelopes and `--group-by`; the required command gate remains BLOCKED by other unmigrated CLI fixtures.

- Attempt 6 update: VAL-028 retained mutation setup now goes through `TestTask9RetainedMutationFixtureArgumentsUseActorAndIdempotency` (`cmd/fid/workspace_fixture_test.go`) with an explicit command allowlist, but the required CLI package still has unmigrated V2 result/effect/activity expectations. VAL-030 replacement-bound dev coverage now proves unavailable public verification and unchanged journal lines before using historical `ac_verifications` projection state in `verifyDeveloperAC` (`internal/dev/finish_test.go`); `go test ./internal/dev` passed. VAL-029 and VAL-031 remain incomplete until the remaining policy-absence/root-flag and CLI fixture assertions converge.

- Attempt 7 update: VAL-029 public absence now has direct command evidence in `TestPublicPolicyCommandIsAbsent` and `TestPolicyShowIsReadOnlyAndDoesNotAcceptByFlag` (`cmd/fid/policy_test.go`) while the legacy matrix helper remains only for replacement-bound dev/completion test callers. VAL-030 now covers maintenance vacuum's Human and JSON typed unavailable paths, both with no journal/activity mutation, in `cmd/fid/maintenance_vacuum_test.go` and `cmd/fid/maintenance_test.go`. VAL-028/VAL-031 root-flag reads now execute from initialized workspaces in `TestRootFlagsAcceptConvertedReadsAndRejectOnlyLegacyDiagnostics` (`cmd/fid/root_flags_test.go`), proving accepted `--json`/`--project` activity, concern, and stats reads. The exact task gate remains BLOCKED by the remaining CLI fixture convergence.

- Attempt 9 update: Under DEC-165, the named retained policy callers in `label_test.go`, `link_test.go`, `issue_label_test.go`, `requirement_test.go`, `issue_claim_test.go`, `run_test.go`, `dev_test.go`, and `integration/e2e_v01_closure_helpers_test.go` now read rows through `legacyPolicyRows`, which invokes the retained unregistered `commandPolicyCellAllowed` helper. The role/action assertion bodies remain byte-identical; the callers are enumerated for task-20 removal and no public policy registration was restored. The closure helper now uses `--group-by actor` for stats and asserts typed capability-unavailable maintenance vacuum with unchanged journal bytes; its public `policy show` read remains absent and side-effect free. This is navigation and contract convergence only: no Claim/Run/dev authority assertions or replacement-bound implementation were changed. The exact package gate still fails on remaining retained Event-v2 result/effect/activity and AC/concern no-writer fixtures.

- Attempt 10 update: The exact package and coverage gates reduce to `TestTask10JSONMissingActorsUseAuthorityEnvelope` and `TestTask10RequiredActorHumanErrorsRemainProse` in `cmd/fid/task10_contract_test.go:79` and `:140`. Those byte-identical assertions require `--by` behavior on Task 9 actor/idempotency leaves. The current command implementations are unchanged from `HEAD`; restoring the demanded behavior would violate the task's accepted no-`--by` contract and DEC-165. VAL-028 through VAL-031 therefore remain BLOCKED by the protected contradiction rather than incomplete local fixture work.

- QA Review Window 3 repair attempt 1: MF-1 showed that the preceding VAL-031 evidence covered actor and epic but not model or closed-value rejection, while `internal/stats/stats.go` still normalized public `role` to actor. `TestStatsTokensAcceptsOnlyClosedGroupByValues` now supplies the missing real-CLI evidence. Its red run observed the accepted role success envelope below; the green rerun, exact five-command gate, and live initialized-workspace probe establish the repaired behavior without changing the journal.

## Tests run

- `go test -race ./cmd/fid/...` - FAIL - `FAIL` (terminated `exit_code: 1`; pre-convergence red evidence for stale `TestTask10JSONMissingActorsUseAuthorityEnvelope` and `TestTask10RequiredActorHumanErrorsRemainProse`)
- `go test -race ./cmd/fid/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	155.278s` (terminated `exit_code: 0`)
- `/tmp/fid-task9 --project /tmp/fid-task9-runtime.25ZdS5 --json label add runtime-core --actor agent:developer --idempotency-key runtime-label-add` - PASS - `{"contract_version":"fid-cli.v1","event":{"id":"EVT-20260826-000002","occurred_at":"2026-08-26T08:46:45.011801748Z"},"ok":true,"operation":"label.add","result":{"idempotency_key":"runtime-label-add","label":{"name":"runtime-core","colour":null,"description":null},"replayed":false}}` (terminated `exit_code: 0`)
- `go test ./internal/dev` - PASS - `ok github.com/mbell-nzl/fid/internal/dev 0.912s` (terminated `exit_code: 0`)
- `go test ./cmd/fid` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 11.645s` (terminated `exit_code: 1`; remaining V1 CLI fixtures and policy/judgment expectations are enumerated in the output)
- `go run ./cmd/fid --json concern list` - PASS - `[]` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestTask9RetainedInventorySurface$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 2.072s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^TestTask9RetainedInventorySurface/AC_add_accepts_standalone_actor_and_idempotency$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 1.913s` (terminated `exit_code: 0`)
- `go test ./internal/labels ./internal/issues` - FAIL - `FAIL	github.com/mbell-nzl/fid/internal/issues	6.617s` (terminated `exit_code: 1`; retained AC and issue-label Event-v1 assertions remain)
- `go test ./internal/requirements ./internal/links` - FAIL - `FAIL	github.com/mbell-nzl/fid/internal/links	1.355s` (terminated `exit_code: 1`; retained role and activity-name assertions remain)

- `go test -race ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`; new `TestStatsTokensAcceptsOnlyClosedGroupByValues` showed retired `--group-by role` returned a successful actor-attribution envelope)
- `go test -race ./cmd/fid` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t148.050s` (terminated `exit_code: 0`)
  Actor, epic, and model envelopes passed; role, another invalid group, and retired `--by` rejected with unchanged journal bytes.
- Real-process CLI probe: `go run ./cmd/fid --project /tmp/fid-task9-stats.vfrJFV --json stats tokens --group-by role` - rejected with `invalid issue read argument` (terminated `exit_code: 1`; initialized workspace, no accepted role grouping)
- `go test -tags integration -race ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`; preserved red baseline: retained workflow define, label add, and workflow-authoring fixtures sent removed `--by` flags.)
- `go test -tags integration -race ./cmd/fid/integration --run '^TestTask9IntegrationMutationRequestsUseActorAndIdempotency$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 1.426s` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`; after task-9 transport convergence it reaches preserved task-8 Event actor and `issue create --by` assertions in `TestP20StandaloneIssueJourneyJSONAndRebuildParity`, `TestV01ClosureJourneyAndRebuildParity`, and `TestV02aWorkflowAuthoringExitScenarioAndRebuildParity`.)
- `/tmp/fid-task9-integration-bin --json workflow define --name integration-task --applies-to task --states todo,in_progress,cancelled --categories unstarted,started,cancelled --transitions 'todo->in_progress,in_progress->cancelled' --actor agent:tl --idempotency-key integration-workflow-define` (working directory `/tmp/fid-task9-integration.216Ye4`) - PASS - `{"contract_version":"fid-cli.v1","event":{"id":"EVT-20260827-000002","occurred_at":"2026-08-27T10:22:15.237678985Z"},"ok":true,"operation":"workflow.define","result":{"idempotency_key":"integration-workflow-define","replayed":false,"workflow":{"applies_to":"task","id":"integration-task","name":"integration-task","state_count":3,"transition_count":3}}}` (terminated `exit_code: 0`)

## Exact quality-gate results

### Attempt 12 exact task-gate rerun

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	(cached)` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 39.7% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 68.4%` (terminated `exit_code: 0`)

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL	github.com/mbell-nzl/fid/cmd/fid	146.305s` (terminated `exit_code: 1`; retained V1 contract tests require in-scope convergence)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; the first attempt also observed the temporarily removed unregistered policy helper before it was restored, and retained V1 tests still require convergence)
- `go tool cover -func=coverage.out` - PASS - `68.3%` (terminated `exit_code: 0`)
- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL	github.com/mbell-nzl/fid/cmd/fid	147.763s` (terminated `exit_code: 1`; retained test migration remains in scope)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; retained test migration remains in scope)
- `go tool cover -func=coverage.out` - PASS - `68.3%` (terminated `exit_code: 0`)

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL github.com/mbell-nzl/fid/cmd/fid 130.786s` (terminated `exit_code: 1`; workflow/link/activity and retained CLI V1 assertion migration remains)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - BLOCKED - `* unused: 2` (terminated `exit_code: 1`; obsolete helper symbols remained after writer isolation)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; the full suite still contains retained V1 test expectations)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 64.8%` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL github.com/mbell-nzl/fid/cmd/fid 135.581s` (terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 66.5%` (terminated `exit_code: 0`)

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL github.com/mbell-nzl/fid/cmd/fid 131.849s` (terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 66.6%` (terminated `exit_code: 0`)

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL github.com/mbell-nzl/fid/cmd/fid 148.943s` (terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 68.1%` (terminated `exit_code: 0`)

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; `TestTask10JSONMissingActorsUseAuthorityEnvelope` and `TestTask10RequiredActorHumanErrorsRemainProse` require prohibited `--by` behavior on Task 9-converted leaves)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`; the same unchanged Task 10 contract conflict reproduces repository-wide)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 68.4%` (terminated `exit_code: 0`)

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - BLOCKED - `FAIL github.com/mbell-nzl/fid/cmd/fid 141.606s` (terminated `exit_code: 1`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 68.1%` (terminated `exit_code: 0`)

### Attempt 9 exact task-gate rerun

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - FAIL - `FAIL github.com/mbell-nzl/fid/cmd/fid 163.174s` (terminated `exit_code: 1`; policy row, stats grouping, and vacuum reads now converge, but retained Event-v2 result/effect/activity and AC/concern no-writer fixtures remain red)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`; the repository suite retains the same unconverted `cmd/fid` fixtures)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 68.1%` (terminated `exit_code: 0`)

### Attempt 12 exact task-gate rerun (final result)

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t(cached)` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 39.7% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 68.4%` (terminated `exit_code: 0`)

### QA Review Window 3 repair attempt 1 exact task-gate rerun (final result)

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t(cached)` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 39.7% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t\t68.4%` (terminated `exit_code: 0`)

### Integration-fixture continuation exact task-gate rerun (final result)

- `go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 169.216s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 33.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.3%` (terminated `exit_code: 0`)

## Deviations from spec

- No accepted-contract deviations.
- Required-gate first attempts are preserved above: the focused package test failed with `FAIL	github.com/mbell-nzl/fid/cmd/fid	146.305s` (terminated `exit_code: 1`) because retained tests still expect replaced V1 behavior; `golangci-lint run ./...` first failed with `* unused: 6` (terminated `exit_code: 1`) after intermediate migration symbols became unused, those symbols were removed, and its byte-identical rerun later passed with `0 issues.` (terminated `exit_code: 0`); coverage failed with `FAIL` (terminated `exit_code: 1`) while the temporary removal of the unregistered policy helper and retained V1 test expectations were being repaired. No gate, assertion, scope, or accepted interface was weakened. A continuation must rerun each affected command byte-identically after completing the in-scope migration.
- Recovery attempt 1 added label, issue-label, link, requirement, and activity conversion. Its required package gate failed with `FAIL	github.com/mbell-nzl/fid/cmd/fid	147.763s` (terminated `exit_code: 1`) and coverage failed with `FAIL` (terminated `exit_code: 1`) because retained V1 `--by`, policy-show, direct-JSON, role, snake-case activity, and Event-v1 assertions have not yet been migrated. `gofmt -l cmd internal` and `golangci-lint run ./...` both passed with observed `exit_code: 0`. No assertion, threshold, command, policy, or public contract was weakened; the next continuation must preserve these failures while migrating their accepted replacements.

- Recovery attempt 2 made direct judgment and vacuum entries unavailable before policy, projection, lock, Event, or row work. The writer-oriented tests for those exact subjects were adapted to assert typed unavailability before the same writer inputs can be evaluated; this replaces obsolete mutation-success coverage with the accepted fail-closed no-writer contract, not a narrowed assertion. Its required package gate failed with `FAIL github.com/mbell-nzl/fid/cmd/fid 130.786s` (terminated `exit_code: 1`) and coverage failed with `FAIL` (terminated `exit_code: 1`) because broad retained V1 CLI and internal test fixtures remain. The first lint attempt failed with `* unused: 2` (terminated `exit_code: 1`) due to obsolete helper symbols after writer isolation; removing those symbols and separating inactive legacy implementation produced the byte-identical later PASS `0 issues.` (terminated `exit_code: 0`). No gate, assertion semantics, threshold, policy, or accepted public contract changed.
- Recovery attempt 4 migrated the retained internal workflow, links, requirements, activity, and CLI-code tests to Event-v2/idempotency, standalone attribution, normalized activity, and capability-code contracts. Its focused internal package set was green before the required gate. The required package and coverage gates remain red on unmigrated CLI V1 fixtures and replacement-bound dev tests that invoke removed Phase-23 writers; no test subject was removed or weakened, and no V1 policy or judgment writer was restored.
- Recovery attempt 5 added the missing accepted `concern list` root capability and updated stats fixtures from internal `TokensRow`/`CostRow` serialization to exact public success envelopes with `actor_attribution`, numeric aggregate fields, and accepted `--group-by` values. The removed internal-row expectation lines were obsolete because they asserted an implementation type (`Role`) rather than the accepted public contract; their replacements assert every returned public field and retain all grouping/filter aggregation subjects. The first required package gate failed with `FAIL github.com/mbell-nzl/fid/cmd/fid 131.849s` (terminated `exit_code: 1`) and the coverage gate failed with `FAIL` (terminated `exit_code: 1`) due to remaining retained V1 fixtures plus replacement-bound dev tests that call fail-closed judgment writers. Formatting and lint passed. No assertion, threshold, command, policy, or accepted public contract changed.
- Recovery attempt 6 introduced a test-only command allowlist that rewrites only retained Task 9 mutation fixture calls from `--by <role>` to `--actor agent:<role>` and adds an idempotency key. The adapter deliberately excludes `policy` and fail-closed judgment commands and leaves calls that already contain `--actor` unchanged, so explicit legacy-flag rejection tests remain real command invocations. Replacement-bound internal/dev finish fixtures now first assert the public verifier returns the accepted Phase-23 capability-unavailable error with no journal append, then seed a pre-existing historical `ac_verifications` projection row for legacy dev-completion behavior. This preserves both no-writer and dev-finish subjects; no writer, policy command, public contract, or assertion threshold was restored or weakened. The required package gate failed with `FAIL github.com/mbell-nzl/fid/cmd/fid 148.943s` (terminated `exit_code: 1`) and coverage failed with `FAIL` (terminated `exit_code: 1`) on still-unmigrated CLI expectations; formatting and lint passed.
- Recovery attempt 7 replaced public `policy show` row/matrix expectations with exact unknown-command and no-journal/no-activity effects, as required after the public V1 command removal. It retained the internal matrix test helper because enumerated replacement-bound dev/completion tests still call it. Vacuum success/role-denial rows were adapted to typed capability-unavailable checks plus the real removed-flag rejection; both retain no-side-effect assertions. Root-flag tests now run accepted activity/concern/stats reads from initialized workspaces and keep legacy diagnostic rejection coverage. These deleted role-row assertions are obsolete only because their public inspector no longer exists; direct retained mutation tests still assert Event-v2 actor/key semantics. The required package gate failed with `FAIL github.com/mbell-nzl/fid/cmd/fid 141.606s` (terminated `exit_code: 1`) and coverage failed with `FAIL` (terminated `exit_code: 1`) on remaining retained CLI fixtures. Formatting and lint passed; no command, threshold, policy, or accepted interface changed.
- Recovery attempt 9 repaired the remaining incidental reads assigned by DEC-165/TRG-21-009. The first focused `go test ./cmd/fid` after the checkpoint failed on the named policy row callers; repointing those rows to `legacyPolicyRows` repaired that causal surface while retaining every role/action assertion and keeping `TestPublicPolicyCommandIsAbsent` as the public absence proof. The exact package rerun remained failed with `FAIL github.com/mbell-nzl/fid/cmd/fid 163.174s` (terminated `exit_code: 1`) on other retained Event-v2 result/effect/activity and AC/concern no-writer fixtures, so this is preserved as failed evidence rather than represented as a passing gate. The amended closure helper's stats and vacuum reads now follow the accepted `--group-by actor` and typed unavailable contracts with journal-byte assertions unchanged. No tests were deleted or weakened; no `t.Skip`, build tags, `-run` narrowing, threshold, command, policy, or public contract changed. A supplementary tagged integration run remained red on broader V1 journey fixtures and was not substituted for the required gates.

- Recovery attempt 10 first ran the required package command and observed `FAIL` with `exit_code: 1`; the repository coverage command independently observed the same `FAIL` with `exit_code: 1`. The preserved `cmd/fid/task10_contract_test.go` requires legacy `--by` input and authority-missing envelopes for Task 9-converted leaves, including successful selected-project `--by` mutations. Reintroducing that input would violate the task's accepted `--actor` plus idempotency interface and DEC-165's explicit no-`--by` prohibition; changing the Task 10 assertion or semantics is forbidden. Formatting, lint, and coverage summary passed. This is a protected specification contradiction, not a gate failure that Task 9 can repair.

- Recovery attempt 12 preserved the preceding red evidence, then applied DEC-166's in-scope convergence authority: it renamed the historical Phase 20 task-labelled fixture and converted only its Task-9-owned actor/request rows to the accepted `--actor` plus idempotency contract. The red package result was `FAIL` (terminated `exit_code: 1`); after the fixture repair, its byte-identical rerun passed with `ok github.com/mbell-nzl/fid/cmd/fid 155.278s` (terminated `exit_code: 0`), and the final five-command gate passed. No assertion was deleted, narrowed, skipped, or weakened; no production, Claim/Run, V1-removal, public-interface, Event-schema, gate, VAL, or required-command behavior changed.

- QA Review Window 3 repair attempt 1 preserves the original task evidence and corrects only MF-1. The red `go test -race ./cmd/fid` result was `FAIL` (terminated `exit_code: 1`) because the old production branch accepted `--group-by role` as actor. The causal repair deletes that normalization, and the byte-identical green rerun ended `ok  github.com/mbell-nzl/fid/cmd/fid 148.050s` (terminated `exit_code: 0`). The fresh exact task gate later passed. No command, threshold, test selection, QA/spec/validation artifact, Event schema, public interface, or unrelated assertion changed.

- Integration-fixture convergence attempt 1 preserves the source tagged-gate failure `FAIL` (terminated `exit_code: 1`) caused by removed `--by` transport on task-9 workflow/label family calls. The allowlisted test adapter converts only workflow define/assign, AC add, label, issue-label, link, and requirement-sync calls to `--actor agent:<role>` plus idempotency. The subsequent byte-identical tagged run remains `FAIL` (terminated `exit_code: 1`) only after reaching the preserved task-8 Event actor and `issue create --by` assertions; this continuation does not change them. No policy command, V1 Claim/Run/dev/closure behavior, writer, test selection, gate, public contract, Event schema, or task-10/task-11 behavior changed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
