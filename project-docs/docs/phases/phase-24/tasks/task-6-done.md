# Phase 24 - Task 6: Gate registration, Envelope and Claim leaves - Recovery Checkpoint

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: Focused Task 6 binding, Source, replay, import and eight-leaf descriptor checks pass, but both full test gates terminate with retained fixtures that Task 6 displaced and whose definitions are outside the exact active recovery allowlist.
- Recovery attempts: 1
- Next action: Phase Runner supplies a bounded same-task scope amendment for the 21 fixture-definition paths listed below; continue causal convergence without weakening their subjects, rerun both failed commands byte-identically within the exact five-command sequence, finalize the report, and require OUTCOME COMPLETE.
- Human attention required: no

OUTCOME RECOVERY_REQUIRED

## Summary of what was implemented

- Continued the existing implementation checkpoint `3ed4824b3f9f7d5dee0f2a50799df538a2911953` and finalized-report HEAD `c845cee5742e20b5418dd30a2780787bf4116a56`. The Runner-confirmed recovery identity, trigger and no-progress count remain unchanged; this continuation made material progress.
- Made Obligation Registration nullable as required for the existing typed review/Human Obligation shape. Removed the unused pre-descriptor Claim decoder. Claim Events now preserve the Envelope repository's base_oid and worktree_id.
- Converged allowlisted Gate, Envelope, Claim and shared integration/fake helpers to accepted identity-free inputs and real Source bytes. Kept the shared same-worktree predicate, domain reasons, fixed fallback messages, and later Run/Gate-run/Evidence interfaces intact.
- Added independently observed Registration digest/record/Event checks; all three Envelope Source positions; original-worktree success, observable whole-worktree plain-copy refusal with no Event, and real export/import destination acceptance. Claim records the local worktree while its Event retains the Envelope scope.
- Added replay after HEAD changes, non-Git and unmerged observations, and Source changes/deletion. Every same-key retry retains original result bytes and appends no Event; otherwise viable distinct-key controls refuse with no Event, then succeed when the unavailable state is restored.
- Completed the eight-leaf descriptor corpus through the built binary: examples, recursive help, closed typed templates, exclusivity, nested unknown/duplicate/missing/null/type/enum/format failures, integer/string/array bounds, sorted-unique constraints, and unchanged project persistence. The derived fake/fid Task 6 intersection is envelope.create and claim.acquire.

## Files created or modified

- `cmd/fid/claim.go`
- `cmd/fid/envelope_test.go`
- `cmd/fid/gate_test.go`
- `cmd/fid/generic_claim_test.go`
- `cmd/fid/input_descriptor_test.go`
- `cmd/fid/input_descriptors_binding.go`
- `cmd/fid/integration/p21_conformance_test.go`
- `cmd/fid/integration/p22_fake_test.go`
- `cmd/fid/integration/p23_amended_claim_oracle_test.go`
- `cmd/fid/integration/p23_amended_envelope_oracle_test.go`
- `cmd/fid/integration/p23_conformance_test.go`
- `cmd/fid/integration/p23_review_record_test.go`
- `internal/claims/generic.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md`
- `project-docs/docs/phases/phase-24/tasks/task-6-done.md`

Earlier implementation paths remain part of the same task's history; this continuation changed only the paths above. Runner-owned ledger dirt is excluded.

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: nullable Obligation Registration, exact-root imported-destination equivalence, plain-copy refusal/no Event, and local Claim worktree versus Envelope Event scope.
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md`: imported Envelope/Claim behavior and destination Source observations.
- This existing report preserves the original failed evidence and records this session's terminated commands and remaining scope gaps.

## Validation Contract Coverage

These are focused Task 6 slice results, not a passing full-task closure. All test names below executed in this session; the final full gates still fail at the listed retained fixture boundaries.

- VAL-009: PASS (Task 6 slice) - `TestP24BindingReplayBeforeObservation` (`cmd/fid/envelope_test.go:90`) covers all four observing leaves after HEAD advancement; `TestP24EnvelopeReplayAfterSourceChangeOrDeletion` (`cmd/fid/envelope_test.go:24`) covers omitted digest replay after both byte changes and deletion, with original result equality/no Event.
- VAL-011: PASS (eight leaves) - `TestP24GateEnvelopeClaimDescriptors` and `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` (`cmd/fid/input_descriptor_test.go:235`) executed in the ordinary gate; the latter includes the seeded mismatch refusal. `TestP24BindingDescriptorProcessCorpus` (`cmd/fid/input_descriptor_test.go:317`) requires exactly the eight binding descriptors and checks their field trees and typed examples.
- VAL-012, VAL-013, VAL-014: PASS (eight leaves) - `TestP24BindingDescriptorProcessCorpus` compares each recursive help field's path/type/requirement/meaning/format/enum/bounds; each real-process example's exact order and two-space formatting; and closed Go decode plus descriptor types. Optional and nullable nested objects are expanded for negative structural rows.
- VAL-015: PASS (eight leaves) - `TestAdoptedExampleProcessExclusivity` (`cmd/fid/input_descriptor_test.go:44`) includes all eight among 20 adopted leaves, with six forbidden combinations each. Examples alone succeed in the process corpus.
- VAL-016: PASS (eight leaves) - `TestP24BindingDescriptorProcessCorpus` derives structural rows from every descriptor field, including nested objects, arrays and scalar elements, and checks exact field path, closed details and no Event.
- VAL-017: PASS (Task 6 slice) - `TestP24RemovedBindingMembers` (`cmd/fid/input_descriptor_test.go:291`) rejects the four removed top-level members; `TestP24EnvelopeSourcePositions` (`cmd/fid/envelope_test.go:327`) rejects removed Source repository in every position with no Event.
- VAL-018, VAL-019: PASS (Envelope positions) - `TestP24EnvelopeSourcePositions` requires exactly risk_result.source, sources[0], obligations[0].source. Matching and omitted digests store SHA-256 of real bytes; mismatched and missing-file rows require the exact mismatch details and unchanged persistence.
- VAL-022: PASS (Task 6 slice) - `TestGateRegisterObservesRepositoryWithoutCallerIdentity` (`cmd/fid/gate_test.go:209`) independently checks Registration and Event; `TestP24EnvelopeSourcePositions` compares all recorded repositories; `TestP24BindingReplayBeforeObservation` checks successful complete Event parity; `TestP24EnvelopeAndClaimOriginalCopyAndImport` (`cmd/fid/envelope_test.go:234`) checks local Claim identity and Envelope-scope Event parity in both original and imported roots.
- VAL-023: PASS (Task 6 slice) - `TestP24BindingReplayBeforeObservation` covers normal acceptance and both non-Git/unmerged refusal for all four observing leaves, including repository-unavailable reason/detail and no Event. The phase-wide classifier's other cause classes remain Task 4 evidence.
- VAL-026: PASS - `TestGateRegisterObservesRepositoryWithoutCallerIdentity` computes project ID from config, worktree ID from root SHA-256, HEAD/tree/dirty from independent Git commands, and canonical binding digest from that repository plus scope, then compares the Registration and Event.
- VAL-027, VAL-028, VAL-029: PASS (Task 6 slice) - `TestP24EnvelopeAndClaimOriginalCopyAndImport` pairs each identical original-worktree request with an observable whole-copy refusal and a real exported/imported destination acceptance. It checks the Task 6 fixed messages, exit 2, closed details and unchanged persistence for both refused leaves. Task 10's later domain-message adoption is not claimed.
- VAL-035: PASS (Task 6 derived intersection) - `TestP23SharedAmendedLeafConformance` (`cmd/fid/integration/p23_conformance_test.go:379`) derives the changed binding leaves from the fake dispatch and requires envelope.create/claim.acquire parity. It covers valid identity-free inputs and removed scope/Source repository/worktree members on both adapters, with no Event for refusals. `TestP22DualAdapterRecoveryConformance` also passes with the converged real-Source shared builder.
- VAL-052: PASS (measured execution budgets only) - final ordinary cmd/fid 197.511s; final integration 285.231s; slowest fresh other package in the final ordinary attempt internal/export 62.311s. internal/workflow is cached there; this session's earlier complete ordinary attempt measured 102.885s, and its slowest other package internal/issues 142.117s. All are below 600s ordinary/900s tagged budgets and below warning thresholds. Correctness gates remain failed.
- VAL-060: PASS (eight leaves) - `TestP24BindingDescriptorProcessCorpus` runs all eight examples through the built binary in a fresh non-project directory, requires exit 0/empty stderr, and compares the directory tree before/after.
- VAL-066: PASS (Task 6 slice) - `TestP24BindingReplayBeforeObservation` and `TestP24EnvelopeReplayAfterSourceChangeOrDeletion` require no-Event same-key replay in unavailable states, otherwise viable distinct-key refusal with no Event, and restored-state new-key success for both repository and each deleted Source position.

## Tests run

Working directory: `/home/mark/projects/fid`, except the explicitly fresh non-project CLI directory. Every new result below was observed terminating in this session. The quoted original report at the end of this section is retained historical evidence, not a fresh attestation.

- `go test ./cmd/fid -run '^TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.777s` (terminated `exit_code: 0`)

The nullable Registration red preceded the descriptor change; the identical command then passed.

- `go test ./cmd/fid -run 'TestGateRegister|TestEnvelope|TestClaim' -count=1 -timeout 10m` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run 'TestGateRegister|TestEnvelope|TestClaim' -count=1 -timeout 10m` - FAIL - `FAIL` (terminated `exit_code: 1`)

The second identical run reached the two retained completion tests in the out-of-allowlist task18_completion_resume_test.go; earlier allowed nullable/binding fixtures were repaired.

- `go test ./cmd/fid -run '^TestP24EnvelopeAndClaimOriginalCopyAndImport$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestP24EnvelopeAndClaimOriginalCopyAndImport$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.062s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestP24BindingDescriptorProcessCorpus
- `go test ./cmd/fid -run '^TestP24BindingDescriptorProcessCorpus$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.572s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestP24BindingReplayBeforeObservation$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestP24BindingReplayBeforeObservation$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.200s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestReviewRecordCreatesImmutableReview$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestReviewRecordCreatesImmutableReview$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.166s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^(TestP24BindingDescriptorProcessCorpus|TestP24EnvelopeSourcePositions|TestP24EnvelopeReplayAfterSourceChangeOrDeletion)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^(TestP24BindingDescriptorProcessCorpus|TestP24EnvelopeSourcePositions|TestP24EnvelopeReplayAfterSourceChangeOrDeletion)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.483s` (terminated `exit_code: 0`)

The expanded corpus initially serialized a non-string invalid value with fmt instead of canonical JSON; the assertion was corrected while preserving its required exact error value. The large failed output was not rerun to obtain a different failure.

- `go test -tags integration ./cmd/fid/integration -run '^TestPhase21RealProcessJourneyStartsOneActivatedEnvelope
- `go test -tags integration ./cmd/fid/integration -run '^TestPhase21RealProcessJourneyStartsOneActivatedEnvelope$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.969s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP22DualAdapterRecoveryConformance$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.688s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP23SharedAmendedLeafConformance
- `go test -tags integration ./cmd/fid/integration -run '^TestP23SharedAmendedLeafConformance$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.114s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP23ReviewRecord' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP23ReviewRecord' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.106s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/evidence ./internal/workflow ./internal/claims` - FAIL - `FAIL` (terminated `exit_code: 1`)

Touched-package tier: cmd/fid failed at remaining fixture boundaries (203.426s); evidence 2.251s, workflow 73.468s and claims 1.033s passed.

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestPhase21RealProcessJourneyStartsOneActivatedEnvelope|TestP22DualAdapterRecoveryConformance|TestP23SharedAmendedLeafConformance|TestP23ReviewRecordRealProcessJourney)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	9.181s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^(TestP24BindingDescriptorProcessCorpus|TestP24BindingReplayBeforeObservation|TestP24EnvelopeSourcePositions|TestP24EnvelopeAndClaimOriginalCopyAndImport|TestP24EnvelopeReplayAfterSourceChangeOrDeletion|TestP24RemovedBindingMembers|TestAdoptedExampleProcessExclusivity|TestClaimAcquireAcceptsReviewMode|TestClaimAcquireKeepsOneActiveClaimAcrossModes|TestClaimAcquireBindingMatrixRejectsWithoutMutation|TestGateRegisterObservesRepositoryWithoutCallerIdentity|TestReviewRecordCreatesImmutableReview)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	27.952s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

Real entrypoint interaction: `/tmp/fid --json claim acquire --actor actor-builder --idempotency-key task6-runtime --input -` ran in a fresh non-project temporary directory with an input containing worktree_id. The child terminated with observed returncode 2, empty stdout and the exact invalid-argument path input.worktree_id; allowed keys were envelope_id, envelope_revision, issue_id, mode, scope, ttl_seconds. A subprocess wrapper checked the closed JSON error and unchanged empty directory, printed `PASS real CLI removed-member refusal: input.worktree_id; no workspace writes`, and terminated with observed exit_code 0.

Original checkpoint/report evidence, retained verbatim with quotation prefixes:

> # Phase 24 - Task 6: Gate registration, Envelope and Claim leaves - Recovery Checkpoint
>
> ## Workflow outcome
>
> - Task outcome: RECOVERY_REQUIRED
> - Category: in-scope-implementation
> - Condition: The eight leaves have partial observation and descriptor implementation, but retained Gate/Envelope/Claim fixtures are red and Task 6 parity, refusal, replay, import and Source matrices remain unfinished.
> - Recovery attempts: 0
> - Next action: Resume Task 6 on this checkpoint. Finish nullable Obligation registration input, converge retained tests and shared integration/fake fixtures with introducing-commit provenance, prove all VAL rows and real CLI behavior, run touched-package race tests, then run all five exact commands in order and complete the report.
> - Human attention required: no
>
> OUTCOME RECOVERY_REQUIRED
>
> ## Summary of what was implemented
>
> - `gate.register` now computes its recorded project/worktree and binding digest from a repository observation after replay. `claim.acquire` derives its worktree from an observation after replay. `envelope.create` observes after replay, resolves all three Source positions, and checks Registration binding through a shared same-worktree/import predicate. `envelope.complete` drops its caller repository input.
> - Adopted descriptors and `--example` for the eight Gate Registration, Envelope and Claim leaves, and added initial focused tests. The implementation and tests are unfinished; this is a checkpoint, not Task 6 acceptance.
>
> ## Files created or modified
>
> - `internal/evidence/registration.go`, `internal/workflow/envelope.go`, `internal/workflow/import_binding.go`, `internal/claims/generic.go`
> - `cmd/fid/gate.go`, `cmd/fid/envelope.go`, `cmd/fid/claim.go`, `cmd/fid/input_descriptor_command.go`, `cmd/fid/input_descriptors_binding.go`
> - `cmd/fid/gate_test.go`, `cmd/fid/envelope_test.go`, `cmd/fid/generic_claim_test.go`, `cmd/fid/input_descriptor_test.go`
> - `project-docs/specs/06_DESIGN_PLAN_CLI.md`, `project-docs/specs/07_DESIGN_PLAN_FLOWS.md`, this report
>
> ## Documentation updated
>
> - `project-docs/specs/06_DESIGN_PLAN_CLI.md` records the current input shapes and observed bindings.
> - `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` records the Gate, Envelope and Claim observation sequence.
>
> ## Validation Contract Coverage
>
> - VAL-011..016, VAL-060: PARTIAL - `TestP24GateEnvelopeClaimDescriptors` proves descriptor/type parity and `TestAdoptedExampleProcessExclusivity` passes for adopted leaves; full eight-leaf help, template, structural and side-effect corpus remains.
> - VAL-022, VAL-026..028: PARTIAL - `TestGateRegisterObservesRepositoryWithoutCallerIdentity`, `TestEnvelopeCreateObservesScopeAndSourcesWithoutCallerRepository`, and `TestClaimAcquireBindsActiveEnvelopeActorWorktreeScopeAndTTL` pass for original-workspace records. Independent full Repository/Event parity and plain-copy rejection remain.
> - VAL-009, VAL-017..019, VAL-023, VAL-029, VAL-035, VAL-066: PENDING - HEAD/replay, removed members, all Source positions, unavailable repositories, imported destinations and fake/fid paired cases still need Task 6 tests and convergence.
> - VAL-052: PENDING - final ordinary and tagged gate package timings were not measured on a green tree.
>
> ## Tests run
>
> - First red: `go test ./cmd/fid -run '^TestClaimAcquireBindsActiveEnvelopeActorWorktreeScopeAndTTL$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); omitted `worktree_id` returned the old `invalid claim request` fallback.
> - After the Claim implementation, the identical focused command - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.634s` (terminated `exit_code: 0`).
> - `go test ./cmd/fid -run '^TestGateRegisterObservesRepositoryWithoutCallerIdentity$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); the old binding validator required caller identity. After the Registration change, the identical command - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.475s` (terminated `exit_code: 0`).
> - `go test ./cmd/fid -run '^TestP24GateEnvelopeClaimDescriptors$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); all eight descriptors were behind the pending fence. After adopting them, the identical command - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.151s` (terminated `exit_code: 0`).
> - `go test ./cmd/fid -run '^TestEnvelopeCreateObservesScopeAndSourcesWithoutCallerRepository$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); the old Envelope input required `repository_scope`. After the Source/observation change, the identical command - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.556s` (terminated `exit_code: 0`).
> - `go test ./cmd/fid -run 'TestGateRegister|TestEnvelope|TestClaim' -count=1 -timeout 10m` - FAIL - `FAIL` (terminated `exit_code: 1`); retained fixtures still send removed members, copy path-bound facts, or need nullable `registration`. Preserve this failure and rerun the identical command after repair.
> - `go test ./cmd/fid -run '^(TestClaimAcquireBindsActiveEnvelopeActorWorktreeScopeAndTTL|TestGateRegisterObservesRepositoryWithoutCallerIdentity|TestEnvelopeCreateObservesScopeAndSourcesWithoutCallerRepository|TestP24GateEnvelopeClaimDescriptors|TestInputDescriptorRegistryMatchesLiveTreeAndTypes)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.676s` (terminated `exit_code: 0`)
> - `go test ./cmd/fid -run '^TestAdoptedExampleProcessExclusivity$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.577s` (terminated `exit_code: 0`)
> - `go test ./cmd/fid ./internal/workflow ./internal/evidence ./internal/claims ./internal/export -run '^$'` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	0.005s [no tests to run]` (terminated `exit_code: 0`)
> - `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> - `/tmp/fid gate register --example` - PASS - `}` (terminated `exit_code: 0`)
>
> The compile-only command ran no tests. The real `fid --example` invocation ran from `/tmp` and printed the new identity-free input template.
>
> ## Exact quality-gate results
>
> Working directory for the required gate: `/home/mark/projects/fid`. None of its commands has been run on a final tree. These are pending commands, not passing evidence.
>
> - `gofmt -l cmd internal` - BLOCKED - not run; Task 6 retained tests and implementation remain unfinished (no termination).
> - `golangci-lint run ./...` - BLOCKED - not run; Task 6 retained tests and implementation remain unfinished (no termination).
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - not run; retained Gate/Envelope/Claim tests are red (no termination).
> - `go tool cover -func=coverage.out` - BLOCKED - not run; no final coverage profile exists (no termination).
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - not run; Task 6 integration convergence is unfinished (no termination).
>
> ## Deviations from spec
>
> - No protected contract or gate definition was changed. The retained focused group ended `FAIL` (terminated `exit_code: 1`) because its inputs still contain removed repository/worktree members, its path-bound fixture caches need replacement, and the descriptor rejects a previously valid nullable Obligation Registration. In-scope convergence is pending; the identical command must be rerun after repair. No required gate attempt has been made.
> - The first report-lint check ended `FAIL report-lint found 5 failure(s)` (terminated `exit_code: 5`) because five passing test/build rows had prose or punctuation after their termination field. The report correction moved that prose to a separate paragraph and left the observed command tokens and statuses intact.
>
> ## Concerns, technical debt, or notes for the Tech Lead
>
> - Retained files already changed: `gate_test.go` was introduced by `0c4c04ea` (Phase 21 Gate task), `envelope_test.go` by `56dad8fa` (Phase 21 Envelope task), `generic_claim_test.go` by `0b67c830` (Phase 21 Claims task), and `input_descriptor_test.go` by `04b6b099` (Phase 24 task 3). The accepted Phase 24 Task 6 input and observation changes, rather than historical names, require convergence. Record introducing commits for any further retained fixture changes before editing.
> - Preserve the negative binding matrix while replacing caller `worktree_id` manipulation with a real plain-copy worktree. Do not remove the wrong-worktree assertion. Preserve the original-workspace positive control, and add imported-destination controls and no-Event checks.
> - The current descriptor's `registration` field needs the accepted nullable shape for review/Human Obligations. Its retained failure is real; changing those inputs to fake Registration objects would weaken coverage.
> - The current worktree has Runner-owned `project-docs/docs/phases/phase-24/triage-log.md` and untracked `run-ledger.md` changes. Leave them untouched during continuation and commits.
>
> ## Dependency-budget note
>
> - No dependency changes.
>
> ## Commit hash
>
> - `3ed4824b3f9f7d5dee0f2a50799df538a2911953`
>
> ## Commit subject
>
> - `feat(binding): checkpoint observed gate envelope claim inputs`

## Exact quality-gate results

Working directory for both complete five-command attempts: `/home/mark/projects/fid`. Each command was run in the specified order and watched terminate. No command was killed, timed out, substituted or silently skipped. Coverage reporting succeeded over a failed-suite profile, which does not establish task correctness.

Initial prefix before the first complete sequence:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - FAIL - `* unused: 1` (terminated `exit_code: 1`)

The lint failure was the now-unused decodeClosedClaimInput. Its removal stayed inside the allowlist.

First complete sequence after that repair:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					59.6%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)

Observed package times: cmd/fid 234.663s, internal/workflow 102.885s, slowest other ordinary package internal/issues 142.117s, tagged integration 273.774s. The real-Source digest in the shared Envelope helper was subsequently corrected; the Review positive control and expanded Task 6 race corpus passed.

Final complete sequence on the checkpoint tree:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					59.8%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)

Both BLOCKED rows describe commands that actually ran and failed, not unavailable execution. Closure is blocked by the exact fixture-definition scope gaps below. The ordinary suite ended FAIL after cmd/fid 197.511s and internal/export 62.311s; the tagged suite ended FAIL after 285.231s. No green full gate is claimed. Report-lint must return OUTCOME RECOVERY_REQUIRED for this durable checkpoint; OUTCOME COMPLETE requires the later exact gates to pass.

## Deviations from spec

- The draft-report writing operation duplicated later sections because JavaScript replacement text interpreted a literal dollar-apostrophe from a test command. The immediately following report-finalization commit replaces that duplicated draft with one copy of each required section and preserves every unique command result; implementation and observed test evidence are unchanged.

- Original failed focused/report-lint evidence is preserved in the quotation above. Historical not-run gate rows have been superseded by two observed complete five-command sequences, never rewritten as historical passes.
- This session's first lint attempt ended `* unused: 1` (terminated `exit_code: 1`); the obsolete decoder was removed and both later identical lint attempts ended `0 issues.` (terminated `exit_code: 0`).
- The nullable Obligation red ended `FAIL` (terminated `exit_code: 1`); minimum descriptor repair accepted the already specified null Registration, and the identical focused command ended `ok  	github.com/mbell-nzl/fid/cmd/fid	0.777s` (terminated `exit_code: 0`).
- The original/copy/import red ended `FAIL` (terminated `exit_code: 1`); it exposed missing Claim Event base_oid/worktree_id and a test helper assuming one journal file. Event construction was corrected and the new independent helper reads all journal files; the identical command ended `ok  	github.com/mbell-nzl/fid/cmd/fid	1.062s` (terminated `exit_code: 0`).
- Initial new-corpus fixture failures ended `FAIL` (terminated `exit_code: 1`) for help rendering assumptions, duplicate-key insertion, placeholder normalization and control-fixture assertion reuse. Each helper was corrected to exercise the specified public behavior; the identical specific commands passed before further implementation. The expanded array-value assertion correction and its identical rerun are recorded above.
- Retained integration red attempts ended `FAIL` (terminated `exit_code: 1`) for obsolete input members and fake observations; the p21 repair additionally exposed a nonexistent helper reference and missing Source directory, and the p23 fake repair additionally exposed typed-null/unknown-member parity. Allowed shared fixtures now use independent observations, real Source bytes and the accepted closed member lists; identical focused commands passed as recorded above.
- The first ordinary and tagged required commands each ended `FAIL` (terminated `exit_code: 1`). The ordinary failure exposed the shared Envelope helper's missing digest when retained Review converted its Source to external-reference; supplying the actual file digest preserved that positive assertion, and explicitly omitted/missing Source rows retained their omission coverage. The byte-identical final ordinary and tagged reruns each still ended `FAIL` (terminated `exit_code: 1`) because remaining definitions are outside the active allowlist. Their inputs were displaced by this task's accepted identity/observation contract; no ownership is assigned to a different task.
- No protected command, assertion semantics, threshold, coverage policy, domain reason/message, same-worktree predicate or later Run/Gate-run/Evidence interface changed. Task closure remains unfinished solely under the current bounded authority; do not approve this report as COMPLETE.

## Concerns, technical debt, or notes for the Tech Lead

- Recovery identity retained: attempt 1; no-progress count 0/2; trigger `task-6:session-continuation:retained-gate-envelope-claim-fixtures-and-nullable-registration`; restored capability authority `restored-capability:developer-role-spawn-01a10b08-efe6-7050-8f88-65c524e62519`; governing TRG-24-003/ESC-24-001. Runner-owned files were neither edited nor committed. This report does not select a new recovery identity or increment Runner counters.
- In-scope retained provenance was checked by introducing-commit history, not inferred from filenames: Gate finalization `0c4c04ea` (Phase 21/task 5); Envelope factory `56dad8fa0dcb9b79f7df2a31b356f1b2eb8a694e` (21/6), repository helper `2339dc9ebd8d0e6b5502ca32e1562892cf704ce6` (22/12); Claim lifecycle `0b67c8304a84e4e5bcb8df44d79f899e40a16acb` and binding matrix `59b079a5a271cc65f7964a471b6383b9c2dd6677` (21/10), review mode `35944c72aba0f2c0b918bdb494b9fa9487d27142` (23/3); descriptor corpus `04b6b0992fe20a02df402a183269a9847ddae5ab` (24/3).
- Shared integration provenance: p21 corpus `a99b4270f70e98111b5bae43c54ed92aed7078c6` and journey exposure `236540c48a2129d70480b3222ba4122a5471eb48` (21/22); p22 fake builder `f3be616f0587881ee3f5c3716ffbf6167e3f2a20` (22/6); amended Envelope/Claim oracles `2e7fd2f9fa05bce648510789453c8d2ed8633b1d` and shared p23 conformance `333514946aeab8328fa618f0f619fdce84694f9a` (23/18); p23AcquireClaim `d215e40c3fbb382d9f0bf17d54a76d50a2e76715` (23/4).
- Remaining fixture definitions below are outside the exact active recovery allowlist. This task's accepted removal of caller identities, real Source resolution and worktree observation caused these failures. A bounded Task 6 amendment should converge these definitions while preserving positive/negative, Stop, Source, import, replay and no-Event assertions. Cache repair must occur before recording path-bound history; creating fictitious import bindings or relaxing the shared predicate would change the protected behavior.

| Needed fixture-definition path | Current failure boundary | Introducing commit | Phase/task |
|---|---|---|---|
| `cmd/fid/generic_run_test.go` | runEnvelopeForTestWithDefinition copies cached observed Envelope history to another root | `d17f524ce3d7bdf7de4e715a03514cba2fe8d7a8` | 21/22 |
| `cmd/fid/evidence_test.go` | readyEvidenceEnvelopeForTest copies cached observed Envelope/Claim/Run facts to another root | `f3350bcf582aacc328c838cd4f08c8b8d9e30f7c` | 21/22 |
| `cmd/fid/review_test.go` | acquireReviewClaim supplies removed worktree_id | `a1ae1bb61363a591bb46c9ddfd4b8df36462687b` | 23/4 |
| `cmd/fid/task18_completion_resume_test.go` | completion inputs supply removed repository | `2313235` | 21/18 |
| `cmd/fid/task21_maintenance_rebuild_test.go` | completion input supplies removed repository | `8435120ac4cc31782c68a1561498a47a4e9fe71b` | 21/21 |
| `cmd/fid/stop_evaluation_matrix_test.go` | both completion Stop states supply removed repository | `4b690e8b863abac96c9abfd72709ac36db2235b2` | 22/9 |
| `internal/export/archive_fixture_test.go` | fixtureSource supplies a placeholder digest for actual project-file bytes | `1be8e882ef689f15705507c448e3b087e23d5f65` | 22/13 |
| `cmd/fid/integration/p22_resume_actions_test.go` | p22NewResumeActionAdapter supplies removed Gate/Envelope/Claim identities | `98020edaa18a744dc868121b304ec527acc303d9` | 22/12 |
| `cmd/fid/integration/p22_recovery_journey_test.go` | recovery journey supplies removed Gate/Envelope/Claim identities | `c69456fc9bf42f811e4999f1968a48c4a82d3cc5` | 22/6 |
| `cmd/fid/integration/p23_review_evidence_budget_test.go` | conflicting Claim input supplies removed worktree_id | `d05c288973df5b71d29dd8cc812374072f0736fd` | 23/6 |
| `cmd/fid/integration/p23_obligation_guard_test.go` | completion helper supplies removed repository | `613e45aad9b10445cffd6b113ee598ab493f183a` | 23/2 |
| `cmd/fid/integration/p23_human_outcome_test.go` | completion and activation fixtures supply removed identities and manufactured Source facts | `0ff2eeef66d4633309a43ae6a341aecd38def4d0` | 23/10 |
| `cmd/fid/integration/p23_human_fixtures_test.go` | p23ReviewExhaustionFixture supplies removed worktree_id; portable callers fail here | `59d2653da7cfe14699dc4d1cd85ef032bec26a5c` | 23/14 |
| `cmd/fid/integration/p23_human_failure_test.go` | Claim Stop/rework rows supply removed worktree_id | `180f5168e3f1338cb758f8388f1c26fd6f8428d3` | 23/19 |
| `cmd/fid/integration/p23_gate_grants_test.go` | Gate Registration fixture supplies removed repository/worktree_id | `8531522e9602b0f72717d6c9d889684228cfa2ee` | 23/6 |
| `cmd/fid/integration/p23_decision_reference_test.go` | completion input supplies removed repository | `828eb18f422ebbfacbe4e65c62d53a061b7e9b9f` | 23/9 |
| `cmd/fid/integration/p23_task13_precedence_test.go` | completion precedence input supplies removed repository | `d1c016daa9f5367a620a3d1908e39aed12d75110` | 23/13 |
| `cmd/fid/integration/p23_concern_deviation_test.go` | completion input supplies removed repository | `3a0aed9d35ca2d43540e686a5cbbaa245177ea0f` | 23/8 |
| `cmd/fid/integration/p24_gate_evidence_corpus_test.go` | completion corpus input supplies removed repository | `9bbbdb55feb7409c8a845d387bc714d927babd8a` | 24/2 |
| `cmd/fid/integration/p24_single_journal_test.go` | completion corpus input supplies removed repository | `95c738f23c34aef1a9e8c6630f1d88c660648591` | 24/2 |
| `cmd/fid/integration/p24_idempotency_test.go` | Claim conflict row supplies removed worktree_id | `8bddc91940dd2f0e012eabc389ffed2a7aa5e987` | 24/2 |

- Newly inspected concern for the next continuation: `internal/claims/generic.go:380` still compares stored Claim worktree directly to the Envelope worktree on renew/release/expire. A Claim acquired in an imported destination records a different local worktree. The Task 6 import acceptance test proves acquisition only; mutation of that newly acquired imported Claim is not claimed as tested. Resolve the accepted stored-authority contract before changing that check; its input shapes and domain reason must stay unchanged.
- There are no measured package-time headroom warnings. Cached workflow timing in the final attempt is explicitly distinguished from this session's earlier fresh measurement.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `1baaecc7983872b06e2a47e6df0801bfb1370ea0`

## Commit subject

- `feat(binding): converge observed envelope and claim recovery`
