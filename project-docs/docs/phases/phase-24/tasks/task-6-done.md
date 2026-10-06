# Phase 24 - Task 6: Gate registration, Envelope and Claim leaves - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 5
- Next action: None.
- Human attention required: no

OUTCOME COMPLETE

## Summary of what was implemented

- Continued the same Task 6 from ending HEAD 45961a91b60bbeae7e1747642637544bbcf89fbe and task/report commit b7d1fd56b801d10378b35dff8f45684d1990804d. TRG-24-004 bounds this attempt to four retained fixture definitions and separate bookkeeping in this report. Earlier Task 6 product, descriptor, Source and import implementation remains in its committed history.
- Execution-transition fixtures now build Registration and Envelope history in the consuming root, followed by a local Claim and Run. Qualifying, wrong-worktree, holder/expiry, ordered-reason and every predicate-matrix row remain intact.
- The expired-Claim isolation regression explicitly expires and rejects a prior local Claim, then verifies a separate fresh fixture has only the expected workflow-edge denial. It uses no sleep, shared post-observation cache mutation or suite-duration assumption.
- The adopted-target journey omits only envelope.complete's removed repository input. Fresh successor proof, retained prior Receipt identity, all terminal transitions, Run close and Envelope completion remain asserted. The historical replay fixture omits only claim.acquire's removed worktree_id and retains all four replay operations, a later Budget charge, exact original/replayed result and Event equality, and no appended Event.
- All five exact gates passed in the specified order. Fresh independent QA window B remains required under the existing QA plan; this Developer report supplies no QA verdict or phase advance.

## Files created or modified

- `cmd/fid/task12_execution_transition_test.go`
- `cmd/fid/execution_transition_fixture_test.go`
- `cmd/fid/integration/p22_recovery_adopted_target_test.go`
- `cmd/fid/integration/p22_run_replay_snapshot_test.go`
- `project-docs/docs/phases/phase-24/tasks/task-6-done.md`

These are attempt 5's entire write scope. Earlier task paths and documentation updates remain in the preserved report below. Neither dirty runner file was staged or edited.

## Documentation updated

- This same completion report records attempt 5 provenance, red/green, runtime evidence and exact gate results, and preserves every earlier report line and failed attempt as a quotation. Existing Task 6 CLI/flow documentation is retained from earlier commits; this continuation changes no specification.

## Validation Contract Coverage

The named ordinary tests below ran in this session's uncached cmd/fid race/coverage gate; named integration tests ran in this session's uncached exact tagged gate. Evidence is limited to Task 6's slice of each assertion, with both qualifying and refusing states cited where applicable.

- VAL-009: PASS (Task 6 slice) - TestP24BindingReplayBeforeObservation (cmd/fid/envelope_test.go:96) covers all four observing leaves after HEAD advancement. TestP24EnvelopeReplayAfterSourceChangeOrDeletion (:30) covers changed and deleted Sources, original-result equality and no Event.
- VAL-011: PASS (eight leaves) - TestP24GateEnvelopeClaimDescriptors (cmd/fid/input_descriptor_test.go:277), TestInputDescriptorRegistryMatchesLiveTreeAndTypes (:235, including seeded mismatch rejection), and TestP24BindingDescriptorProcessCorpus (:317) require the eight binding descriptors and matching closed request types.
- VAL-012, VAL-013, VAL-014: PASS (eight leaves) - TestP24BindingDescriptorProcessCorpus compares every recursive help field, exact example key order/two-space formatting, closed Go decoding and descriptor types. It expands optional/nullable objects for structural negative rows.
- VAL-015: PASS (eight leaves) - TestAdoptedExampleProcessExclusivity (cmd/fid/input_descriptor_test.go:44) rejects all six forbidden option combinations. TestP24BindingDescriptorProcessCorpus accepts example-only invocations for every binding leaf.
- VAL-016: PASS (eight leaves) - TestP24BindingDescriptorProcessCorpus derives unknown, duplicate, missing, null, type, enum, format and bounds rows from every descriptor field, including nested positions, and checks exact field paths, closed details and unchanged persistence.
- VAL-017: PASS (Task 6 slice) - TestP24RemovedBindingMembers (cmd/fid/input_descriptor_test.go:291) rejects removed top-level identities. TestP24EnvelopeSourcePositions (cmd/fid/envelope_test.go:416) rejects Source repository at every position without an Event. The two repaired integration fixtures now omit only the respective removed member while their accepted journeys still complete.
- VAL-018, VAL-019: PASS (all Envelope Source positions) - TestP24EnvelopeSourcePositions requires risk_result.source, sources[0] and obligations[0].source. Matching/omitted digest rows store SHA-256 of actual bytes; mismatch/missing-file rows require accepted details and no state change. TestP24EnvelopeReplayAfterSourceChangeOrDeletion checks deleted-file new-key rejection and restored-file acceptance.
- VAL-022: PASS (Task 6 slice) - TestGateRegisterObservesRepositoryWithoutCallerIdentity (cmd/fid/gate_test.go:209) independently compares Registration/Event. TestP24EnvelopeSourcePositions compares Source repositories; TestP24BindingReplayBeforeObservation compares completion Event observation; TestP24EnvelopeAndClaimOriginalCopyAndImport (cmd/fid/envelope_test.go:240) compares original/imported local Claim identity with Envelope-scope Event repository.
- VAL-023: PASS (Task 6 slice) - TestP24BindingReplayBeforeObservation tests normal success against non-Git and unmerged refusal for all four observing leaves, mapped repository-unavailable details and no Event. The phase-wide classifier's other cause classes remain Task 4 evidence.
- VAL-026: PASS - TestGateRegisterObservesRepositoryWithoutCallerIdentity independently computes config project ID, root hash worktree ID, Git HEAD/tree/dirty and canonical binding digest before comparing Registration and Event.
- VAL-027, VAL-028: PASS (Task 6 slice) - TestP24EnvelopeAndClaimOriginalCopyAndImport pairs identical original-root requests with whole-worktree plain-copy refusal and genuine imported-destination acceptance. It checks the accepted fixed fallback, invalid-argument exit 2, closed details and unchanged persistence. Task 10's later domain-message adoption is not claimed.
- VAL-029: PASS (Task 6 slice) - TestP24EnvelopeAndClaimOriginalCopyAndImport proves imported Envelope creation and Claim acquisition. TestP24ImportedClaimStoredAuthority (cmd/fid/envelope_test.go:333) covers original/imported acquire then renew/release/expire, copied-root rejection for each, unchanged fallback/exit/details/no Event, replay and release of the same renewed Claim. The freshly built real CLI interaction below also passes imported lifecycle and copied-root renewal rejection.
- VAL-035: PASS (Task 6 derived intersection) - TestP23SharedAmendedLeafConformance (cmd/fid/integration/p23_conformance_test.go:379) derives the fake/fid binding intersection envelope.create/claim.acquire and checks identity-free acceptance and removed-member refusal/no Event. TestP22DualAdapterRecoveryConformance exercises the converged real-Source builder in the exact tagged gate.
- VAL-060: PASS (eight leaves) - TestP24BindingDescriptorProcessCorpus executes every binding example through the binary in a fresh non-project directory, requires exit 0/empty stderr and compares the complete directory tree before/after.
- VAL-066: PASS (Task 6 slice) - TestP24BindingReplayBeforeObservation and TestP24EnvelopeReplayAfterSourceChangeOrDeletion pair no-Event same-key replay in unavailable/deleted states with otherwise viable new-key refusal/no Event and restored-state success.
- VAL-052: PASS (observed execution budgets) - cmd/fid 283.342s; slowest other freshly timed ordinary package internal/export 116.978s; cmd/fid/integration 274.935s. internal/workflow was explicitly `(cached)` in this session's exact ordinary gate and emitted no fresh wall time; attempt 4's 114.227s remains historical evidence below, not a fresh attestation. Other cached package rows are also distinguished from fresh timing. The pre-gate cmd/fid race check took 304.689s. Fresh ordinary timings are below 600s and tagged timing below 900s; none reaches the respective 480s/720s warning. No cached duration is invented.

## Tests run

Working directory: /home/mark/projects/fid. The runtime driver uses newly allocated Git workspace roots. Focused red precedes edits for all four subjects; each focused green uses the byte-identical command. Every result below has an observed command termination.

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^(TestExecutionBound.*|TestExecutionTransitionFixtureDoesNotInheritExpiredClaim)$'` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^(TestExecutionBound.*|TestExecutionTransitionFixtureDoesNotInheritExpiredClaim)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	18.053s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestP22RecoveryAdoptedTargetPermitsTerminalTransitionWithFreshSuccessorProof|TestP22RunReplayRetainsHistoricalBudgetSnapshot)$'` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestP22RecoveryAdoptedTargetPermitsTerminalTransitionWithFreshSuccessorProof|TestP22RunReplayRetainsHistoricalBudgetSnapshot)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.483s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	304.689s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `python3 /tmp/task6-runtime.py` - PASS - `PASS real CLI imported Claim acquire/renew/release/expire: exit 0; copy renew: exit 2; replay and refusal append no Event` (terminated `exit_code: 0`)
- `rg -n "func (same[A-Za-z]*Repository|validateRepositoryBinding|validEnvelopeRepository|validEnvelopeSource|validEvidenceRepository|validateRunTerminalRepository|verifyDecisionSource|importedGateDestination)" internal -g '!*_test.go'` - PASS - `internal/workflow/envelope.go:759:func validEnvelopeSource(source Source) bool {` (terminated `exit_code: 0`)
- `rg -n "Request\.Repository|request\.Repository|Request\.TargetRepository|request\.RepositoryScope|Request\.WorktreeID|request\.WorktreeID" internal -g '!*_test.go'` - PASS - `internal/workflow/envelope.go:935:	repositoryJSON := mustEnvelopeJSON(request.RepositoryScope)` (terminated `exit_code: 0`)

The focused ordinary red terminated with FAIL, cmd/fid 8.444s: ordered-reason, worktree, holder/expiry and matrix setup, plus the isolation regression, were rejected by claim.acquire with invalid-argument/invalid claim request before their intended assertions. The focused tagged red terminated with FAIL, integration 7.643s: the successor-proof completion rejected input.repository with exit 2, and run.abandon replay setup rejected input.worktree_id with exit 2. Its other three replay operation subjects remain in the identical rerun. The green ordinary run reaches the full predicate matrix and both expired-prior and fresh-fixture controls. The green tagged run reaches successful completion and all four historical replay assertions.

The freshly built /tmp/fid entrypoint was exercised independently by python3 /tmp/task6-runtime.py using actual activation, Registration, Envelope, export/import and identity-free Claim inputs. The driver observes imported acquire/renew/release/expire exit 0 and copied-root renewal exit 2 with the accepted failure envelope; replay and refusal preserve persistence. Its own PASS token and exit_code 0 are recorded above. The two tagged fixture subjects additionally execute real subprocess CLI journeys.

Historical report retained verbatim, with one additional quotation prefix per physical line. All statuses in this quotation are historical evidence, not fresh attestations; no prior failure or deviation is rewritten:

> # Phase 24 - Task 6: Gate registration, Envelope and Claim leaves - Recovery Checkpoint
>
> ## Workflow outcome
>
> - Task outcome: RECOVERY_REQUIRED
> - Category: in-scope-implementation
> - Condition: All 21 DEC-283 retained fixture rows are converged and the imported Claim lifecycle and focused binding corpus pass. Both full test gates still fail at four additional fixture paths outside the exact 51-path allowlist: execution-transition setup copies observed Envelope history, its expiry regression depends on that cache, and two tagged fixtures retain removed input members.
> - Recovery attempts: 4
> - Next action: Obtain bounded same-task authority for cmd/fid/task12_execution_transition_test.go, cmd/fid/execution_transition_fixture_test.go, cmd/fid/integration/p22_recovery_adopted_target_test.go and cmd/fid/integration/p22_run_replay_snapshot_test.go; build path-bound setup locally and omit displaced completion/Claim inputs while preserving expiry, transition, successor-proof and historical replay assertions; rerun both failed commands byte-identically in all five exact gates before Task 6 closure.
> - Human attention required: no
>
> OUTCOME RECOVERY_REQUIRED
>
> ## Summary of what was implemented
>
> - Resumed implementation checkpoint 1baaecc7983872b06e2a47e6df0801bfb1370ea0 and finalized report 7298964a65beb375a6b50ec45c42ea12dada6c59 from phase HEAD 13164582232b9b43eddac93a0ef75b6f4465783c. This is the same Task 6 recovery attempt 4, with material progress and no change to Runner counters or ledger.
> - Converged all 21 retained fixture rows accepted by DEC-283 clause 7. Removed displaced Gate/Envelope/Claim input identities, resolved real Source bytes, and replaced in-scope post-observation caches with root-local setup. Positive, negative, import, plain-copy, Source, replay, Stop, readiness, conflict and no-Event assertions remain.
> - Preserved caller Repository snapshots for retained later Run/Gate-run/Evidence fixtures: each snapshot is independently observed in its own root before output files are written. No history crosses roots and no later runtime interface or ExpectedRepository behavior changed.
> - Reproduced imported-destination renew/release/expire failure after acquiring a Claim there, then changed validateStoredClaimAuthority to the shared workflow.SameWorktree predicate under DEC-283 clause 8. Original-root controls, copied-destination refusal, replay/no-Event, and release of the same renewed Claim pass.
> - Confirmed the checkpoint's nullable Obligation Registration, independently observed Registration/Source/Event parity, all three Source positions, repository/Source replay and viable new-key refusal matrices, original/copy/import acceptance, fake/fid parity and complete eight-leaf descriptor corpus through the focused and required commands recorded below.
> - The ordinary full gate reaches the execution-transition copied-history scope gap. The tagged full gate reaches two additional retained input members. All four defining paths are outside the exact allowlist and remain untouched.
>
> ## Files created or modified
>
> - `cmd/fid/envelope_test.go`
> - `cmd/fid/evidence_test.go`
> - `cmd/fid/generic_run_test.go`
> - `cmd/fid/integration/p22_recovery_journey_test.go`
> - `cmd/fid/integration/p22_resume_actions_test.go`
> - `cmd/fid/integration/p23_concern_deviation_test.go`
> - `cmd/fid/integration/p23_decision_reference_test.go`
> - `cmd/fid/integration/p23_gate_grants_test.go`
> - `cmd/fid/integration/p23_human_failure_test.go`
> - `cmd/fid/integration/p23_human_fixtures_test.go`
> - `cmd/fid/integration/p23_human_outcome_test.go`
> - `cmd/fid/integration/p23_obligation_guard_test.go`
> - `cmd/fid/integration/p23_review_evidence_budget_test.go`
> - `cmd/fid/integration/p23_review_record_test.go`
> - `cmd/fid/integration/p23_task13_precedence_test.go`
> - `cmd/fid/integration/p24_gate_evidence_corpus_test.go`
> - `cmd/fid/integration/p24_idempotency_test.go`
> - `cmd/fid/integration/p24_single_journal_test.go`
> - `cmd/fid/review_test.go`
> - `cmd/fid/stop_evaluation_matrix_test.go`
> - `cmd/fid/task18_completion_resume_test.go`
> - `cmd/fid/task21_maintenance_rebuild_test.go`
> - `internal/claims/generic.go`
> - `internal/export/archive_fixture_test.go`
> - `project-docs/docs/phases/phase-24/tasks/task-6-done.md`
> - `project-docs/specs/06_DESIGN_PLAN_CLI.md`
> - `project-docs/specs/07_DESIGN_PLAN_FLOWS.md`
>
> These are this continuation's edits. Earlier Task 6 implementation paths remain in its committed history and in the quoted prior report. Runner-owned ledger changes are excluded.
>
> ## Documentation updated
>
> - project-docs/specs/06_DESIGN_PLAN_CLI.md and project-docs/specs/07_DESIGN_PLAN_FLOWS.md document imported Claim renew/release/expire authority through the same-worktree rule and copied-root refusal.
> - This same completion report preserves all earlier evidence and records the attempt-4 red/green results, exact gates, retained provenance and remaining bounded scope gap.
>
> ## Validation Contract Coverage
>
> These are Task 6 slice results. A passing slice does not close the failed full ordinary gate.
>
> - VAL-009: PASS (Task 6 slice) - TestP24BindingReplayBeforeObservation (cmd/fid/envelope_test.go:96) tests all four observing leaves after HEAD advancement; TestP24EnvelopeReplayAfterSourceChangeOrDeletion (:30) tests Source byte changes and deletion, original result equality and no Event. Both passed in the focused race command.
> - VAL-011: PASS (eight leaves) - TestP24BindingDescriptorProcessCorpus (cmd/fid/input_descriptor_test.go:317) checks the exact eight binding descriptors and closed template types. TestP24GateEnvelopeClaimDescriptors (:277) and TestInputDescriptorRegistryMatchesLiveTreeAndTypes (:235, including seeded mismatch refusal) executed in the ordinary command.
> - VAL-012, VAL-013, VAL-014: PASS (eight leaves) - TestP24BindingDescriptorProcessCorpus compares every recursive help field's path/type/requirement/meaning/format/enum/bounds, template key order and two-space formatting, and closed Go decoding. Optional/nullable objects remain expanded for structural negative rows. Its focused race command terminated successfully.
> - VAL-015: PASS (eight leaves) - TestAdoptedExampleProcessExclusivity (cmd/fid/input_descriptor_test.go:44) executed all six forbidden option combinations per adopted leaf in the ordinary command; TestP24BindingDescriptorProcessCorpus proves example-only success for all eight through the binary.
> - VAL-016: PASS (eight leaves) - TestP24BindingDescriptorProcessCorpus derives nested unknown, duplicate, missing, null, type, enum, format and bounds rows from every descriptor field, and checks the exact path, closed details and unchanged persistence. The focused race command passed.
> - VAL-017: PASS (Task 6 slice) - TestP24RemovedBindingMembers (cmd/fid/input_descriptor_test.go:291) rejects the four removed top-level identities; TestP24EnvelopeSourcePositions (cmd/fid/envelope_test.go:416) rejects removed Source repository in every position without an Event. Both passed in the focused race command.
> - VAL-018, VAL-019: PASS (all Envelope positions) - TestP24EnvelopeSourcePositions requires exactly risk_result.source, sources[0] and obligations[0].source. Matching/omitted digest rows store SHA-256 of real bytes; mismatch/missing-file rows require the accepted details and no state change. TestP24EnvelopeReplayAfterSourceChangeOrDeletion separately checks deleted Source new-key refusal and restored-file acceptance.
> - VAL-022: PASS (Task 6 slice) - TestGateRegisterObservesRepositoryWithoutCallerIdentity (cmd/fid/gate_test.go:209) independently compares Registration/Event; TestP24EnvelopeSourcePositions compares Source repositories; TestP24BindingReplayBeforeObservation compares complete Event observation; TestP24EnvelopeAndClaimOriginalCopyAndImport (cmd/fid/envelope_test.go:240) compares original/imported local Claim identity with Envelope-scope Event repository. The focused commands and ordinary execution cover these named tests.
> - VAL-023: PASS (Task 6 slice) - TestP24BindingReplayBeforeObservation tests normal acceptance versus non-Git and unmerged refusal for all four leaves, mapped repository-unavailable details and no Event. The phase-wide classifier's other cause classes remain Task 4 evidence.
> - VAL-026: PASS - TestGateRegisterObservesRepositoryWithoutCallerIdentity independently computes project ID, root hash worktree ID, Git HEAD/tree/dirty and canonical binding digest, then compares Registration and Event. The fresh GateRegister/Envelope/Claim selection passed.
> - VAL-027, VAL-028: PASS (Task 6 slice) - TestP24EnvelopeAndClaimOriginalCopyAndImport pairs identical requests in the original root with whole-worktree plain-copy refusal and real imported-destination acceptance. It requires the accepted fixed fallback, invalid-argument exit 2, closed details and no persistence change. The focused race command passed; Task 10 domain-message adoption is not claimed.
> - VAL-029: PASS (Task 6 slice) - TestP24EnvelopeAndClaimOriginalCopyAndImport proves imported Envelope creation and Claim acquisition against genuine export/import history. TestP24ImportedClaimStoredAuthority (cmd/fid/envelope_test.go:333) adds original/imported acquire then renew/release/expire, copied-destination refusal for each with exact unchanged fallback/exit/details/no Event, same-key replay/no Event, and release of the same renewed Claim. All passed under the race detector and through the real CLI interaction below.
> - VAL-035: PASS (Task 6 derived intersection) - TestP23SharedAmendedLeafConformance (cmd/fid/integration/p23_conformance_test.go:379) derives the fake/fid binding intersection envelope.create/claim.acquire, checks identity-free success and removed-member refusal/no Event. TestP22DualAdapterRecoveryConformance uses the converged real-Source builder. These named tests executed in the exact tagged gate; its actual result is recorded below.
> - VAL-052: PASS (measured execution budgets only) - cmd/fid 316.372s; internal/workflow 114.227s; slowest fresh other ordinary package internal/export 144.968s; tagged integration 332.169s. All are below the 600s ordinary/900s tagged budgets and below 80% headroom thresholds. Correctness of the ordinary suite remains failed.
> - VAL-060: PASS (eight leaves) - TestP24BindingDescriptorProcessCorpus executes every binding example through the binary in a fresh non-project directory, requires exit 0/empty stderr and compares the complete directory tree before/after.
> - VAL-066: PASS (Task 6 slice) - TestP24BindingReplayBeforeObservation and TestP24EnvelopeReplayAfterSourceChangeOrDeletion pair no-Event same-key replay in unavailable/deleted states with distinct viable new-key refusal/no Event and restored-state success. The focused race command passed.
>
> ## Tests run
>
> Working directory: /home/mark/projects/fid, except fresh scratch roots used by real CLI interactions. Every new PASS below has an observed nested command termination. Failed attempts are retained individually; the quoted prior report remains historical evidence.
>
> - `go test ./cmd/fid -run '^TestP24ImportedClaimStoredAuthority$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test ./cmd/fid -run '^TestP24ImportedClaimStoredAuthority$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test ./cmd/fid -run '^TestP24ImportedClaimStoredAuthority$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.549s` (terminated `exit_code: 0`)
> - `go test ./cmd/fid ./internal/export -run '^(TestRunStartBindsActiveClaimEnvelopeActorRepositoryAndBudgets|TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestReviewRecordIndependentReworkAndReceiptBinding|TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease|TestFinalSchemaNamedReadsAndExportMatchAfterRebuild|TestStopPredicatesAcrossIssueAndEnvelopeCommands|TestCreateRecoveryManifestV2ConfigAndOrdinalJournal)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test ./cmd/fid ./internal/export -run '^(TestRunStartBindsActiveClaimEnvelopeActorRepositoryAndBudgets|TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestReviewRecordIndependentReworkAndReceiptBinding|TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease|TestFinalSchemaNamedReadsAndExportMatchAfterRebuild|TestStopPredicatesAcrossIssueAndEnvelopeCommands|TestCreateRecoveryManifestV2ConfigAndOrdinalJournal)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test ./cmd/fid ./internal/export -run '^(TestRunStartBindsActiveClaimEnvelopeActorRepositoryAndBudgets|TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestReviewRecordIndependentReworkAndReceiptBinding|TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease|TestFinalSchemaNamedReadsAndExportMatchAfterRebuild|TestStopPredicatesAcrossIssueAndEnvelopeCommands|TestCreateRecoveryManifestV2ConfigAndOrdinalJournal)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test ./cmd/fid ./internal/export -run '^(TestRunStartBindsActiveClaimEnvelopeActorRepositoryAndBudgets|TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestReviewRecordIndependentReworkAndReceiptBinding|TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease|TestFinalSchemaNamedReadsAndExportMatchAfterRebuild|TestStopPredicatesAcrossIssueAndEnvelopeCommands|TestCreateRecoveryManifestV2ConfigAndOrdinalJournal)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test -tags integration ./cmd/fid/integration -run '^(TestP22ResumeReadyCandidatesExecuteTheirActualOperations|TestP22EarlyRecoveryJourney|TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation|TestP23GateRunAuthorizesEveryEnvelopeGrantThroughBuiltBinary)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test ./cmd/fid ./internal/export -run '^(TestRunStartBindsActiveClaimEnvelopeActorRepositoryAndBudgets|TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestReviewRecordIndependentReworkAndReceiptBinding|TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease|TestFinalSchemaNamedReadsAndExportMatchAfterRebuild|TestStopPredicatesAcrossIssueAndEnvelopeCommands|TestCreateRecoveryManifestV2ConfigAndOrdinalJournal)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	0.314s` (terminated `exit_code: 0`)
> The following compile failure ended with a compiler diagnostic as its last printed line. Its exact command, failed status, token and observed termination are retained verbatim:
>
> > - `go test -tags integration ./cmd/fid/integration -run '^(TestP22ResumeReadyCandidatesExecuteTheirActualOperations|TestP22EarlyRecoveryJourney|TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation|TestP23GateRunAuthorizesEveryEnvelopeGrantThroughBuiltBinary)$' -count=1` - FAIL - `cmd/fid/integration/p23_human_outcome_test.go:329:62: undefined: sha256` (terminated `exit_code: 1`)
> - `go test -tags integration ./cmd/fid/integration -run '^(TestP22ResumeReadyCandidatesExecuteTheirActualOperations|TestP22EarlyRecoveryJourney|TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation|TestP23GateRunAuthorizesEveryEnvelopeGrantThroughBuiltBinary)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test -tags integration ./cmd/fid/integration -run '^(TestP22ResumeReadyCandidatesExecuteTheirActualOperations|TestP22EarlyRecoveryJourney|TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation|TestP23GateRunAuthorizesEveryEnvelopeGrantThroughBuiltBinary)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test -tags integration ./cmd/fid/integration -run '^(TestP22ResumeReadyCandidatesExecuteTheirActualOperations|TestP22EarlyRecoveryJourney|TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation|TestP23GateRunAuthorizesEveryEnvelopeGrantThroughBuiltBinary)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.071s` (terminated `exit_code: 0`)
> - `rg -n "func (same[A-Za-z]*Repository|validateRepositoryBinding|validEnvelopeRepository|validEnvelopeSource|validEvidenceRepository|validateRunTerminalRepository|verifyDecisionSource|importedGateDestination)" internal -g '!*_test.go'` - PASS - `internal/workflow/envelope.go:759:func validEnvelopeSource(source Source) bool {` (terminated `exit_code: 0`)
> - `rg -n "Request\.Repository|request\.Repository|Request\.TargetRepository|request\.RepositoryScope|Request\.WorktreeID|request\.WorktreeID" internal -g '!*_test.go'` - PASS - `internal/workflow/run.go:233:	if envelope.Status != "active" || envelope.PrimaryIssueID != request.IssueID || envelope.CurrentRevision != request.EnvelopeRevision || !sameRepository(envelope.RepositoryScope, request.Repository) {` (terminated `exit_code: 0`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/evidence ./internal/workflow ./internal/claims ./internal/export` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test ./cmd/fid -run '^(TestGateRunPublishesSanitizedArtifactByDigest|TestArtifactLossHonorsRequiredMaterialAndPreservesProofIdentity|TestGateRunUnsafeArtifactsEmitControlsAndRemainRequiredUnavailable)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.941s` (terminated `exit_code: 0`)
> - `python3 /tmp/task6-runtime.py` - FAIL - `AssertionError: (['/tmp/fid', '--project', '/tmp/fid-task6-runtime-p_2eqz3l/original', '--json', 'activation', 'prepare', '--input', '-', '--idempotency-key', 'prepare'], 2, '', '{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["accepted activation input"],"argument":"input","value":null},"message":"invalid activation input","retryable":false,"subject":null},"ok":false,"operation":"activation.prepare"}\n')` (terminated `exit_code: 1`)
> - `python3 /tmp/task6-runtime.py` - FAIL - `AssertionError: (['/tmp/fid', '--project', '/tmp/fid-task6-runtime-re1ip4wp/original', '--json', 'activation', 'prepare', '--input', '/tmp/fid-task6-runtime-re1ip4wp/request-0.json', '--idempotency-key', 'prepare'], 2, '', '{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["accepted activation input"],"argument":"input","value":null},"message":"invalid activation input","retryable":false,"subject":null},"ok":false,"operation":"activation.prepare"}\n')` (terminated `exit_code: 1`)
> - `python3 /tmp/task6-runtime.py` - FAIL - `AssertionError: (['/tmp/fid', '--project', '/tmp/fid-task6-runtime-7bf6i6_m/imported', '--json', 'export', 'import', '/tmp/fid-task6-runtime-7bf6i6_m/archive.zip', '--repository', '/tmp/fid-task6-runtime-7bf6i6_m/imported', '--idempotency-key', 'import'], 7, '', 'error: --project is not available for export import until its fid-cli.v1 transport is migrated\n')` (terminated `exit_code: 1`)
> - `python3 /tmp/task6-runtime.py` - PASS - `PASS real CLI imported Claim acquire/renew/release/expire: exit 0; copy renew: exit 2; replay and refusal append no Event` (terminated `exit_code: 0`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestP22EarlyRecoveryJourney|TestP22RecoveryConformanceCorpusHasSharedCases|TestP22ResumeObserveRejectsUntrustedProducerDespiteCapabilities|TestP22ResumeReadyCandidatesExecuteTheirActualOperations|TestP23ConcernDeviationRealProcessJourney|TestP23DecisionReferenceProtectedLifecycleRealProcess|TestP23FailurePrecedenceWithProtectedDecisionAndObligation|TestP23FailurePrecedenceWithReadinessStopAndObligation|TestP23GateRunAuthorizesEveryEnvelopeGrantThroughBuiltBinary|TestP23HumanFixtureStepsRemainRunnable|TestP23HumanFixtures|TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess|TestP23HumanOutcomeRequiredRaisedByArtifactAvailability|TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation|TestP23HumanOutcomeRequiredRaisedByReviewRecord|TestP23ReviewCycleChargeSurvivesLinkedRecoveryAndRebuild|TestP23ReviewCyclesRealProcessJourney|TestP23ReviewEvidenceRealProcessJourney|TestP23ReviewObligationGuardingBothActionsBlocksBothRealProcess|TestP23ReviewReworkManualResolutionWithPriorOutcome|TestP23StopFailureCodeMatrix|TestP23StopFailureHumanAndJSONParity|TestP23UnsatisfiedObligationBlocksOnlyTheActionItGuardsRealProcess|TestP24ActivationPrepareConflictNamesChangedKit|TestP24All49KeyedLeavesRebuildInOneJournal|TestP24ClaimAcquireConflictNamesChangedLease|TestP24GateEvidenceAndEnvelopeKeyedLeafCorpus|TestP24GateRunCommittedDenialPersistsNamedFields|TestP24GateRunPersistsNamedFieldsThroughReceipt|TestP24HumanOutcomeRejectsReferenceDigestMismatch|TestP24IssueCreateConflictNamesChangedTitle|TestP24ObservingJudgmentLeavesNameChangedPublicField|TestP24ObservingJudgmentReplayNeverObserves|TestP24ObservingJudgmentRetriesAfterSourceAndHEADChange|TestP24ReviewOptionalPresenceAndEmptyException|TestP24RunStartConflictNamesChangedClaimVersion|TestP24StopRaiseConflictNamesChangedDetails)$'` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test -tags integration ./cmd/fid/integration -run '^TestP23DecisionReferenceProtectedLifecycleRealProcess$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `go test -tags integration ./cmd/fid/integration -run '^TestP23DecisionReferenceProtectedLifecycleRealProcess$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.267s` (terminated `exit_code: 0`)
> - `go test ./cmd/fid -run 'TestGateRegister|TestEnvelope|TestClaim' -count=1 -timeout 10m` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	4.941s` (terminated `exit_code: 0`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^(TestP24ImportedClaimStoredAuthority|TestP24BindingDescriptorProcessCorpus|TestP24BindingReplayBeforeObservation|TestP24EnvelopeSourcePositions|TestP24EnvelopeAndClaimOriginalCopyAndImport|TestP24EnvelopeReplayAfterSourceChangeOrDeletion|TestP24RemovedBindingMembers|TestGateRunPublishesSanitizedArtifactByDigest|TestArtifactLossHonorsRequiredMaterialAndPreservesProofIdentity|TestGateRunUnsafeArtifactsEmitControlsAndRemainRequiredUnavailable|TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	59.959s` (terminated `exit_code: 0`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestP22EarlyRecoveryJourney|TestP22RecoveryConformanceCorpusHasSharedCases|TestP22ResumeObserveRejectsUntrustedProducerDespiteCapabilities|TestP22ResumeReadyCandidatesExecuteTheirActualOperations|TestP23ConcernDeviationRealProcessJourney|TestP23DecisionReferenceProtectedLifecycleRealProcess|TestP23FailurePrecedenceWithProtectedDecisionAndObligation|TestP23FailurePrecedenceWithReadinessStopAndObligation|TestP23GateRunAuthorizesEveryEnvelopeGrantThroughBuiltBinary|TestP23HumanFixtureStepsRemainRunnable|TestP23HumanFixtures|TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess|TestP23HumanOutcomeRequiredRaisedByArtifactAvailability|TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation|TestP23HumanOutcomeRequiredRaisedByReviewRecord|TestP23ReviewCycleChargeSurvivesLinkedRecoveryAndRebuild|TestP23ReviewCyclesRealProcessJourney|TestP23ReviewEvidenceRealProcessJourney|TestP23ReviewObligationGuardingBothActionsBlocksBothRealProcess|TestP23ReviewReworkManualResolutionWithPriorOutcome|TestP23StopFailureCodeMatrix|TestP23StopFailureHumanAndJSONParity|TestP23UnsatisfiedObligationBlocksOnlyTheActionItGuardsRealProcess|TestP24ActivationPrepareConflictNamesChangedKit|TestP24All49KeyedLeavesRebuildInOneJournal|TestP24ClaimAcquireConflictNamesChangedLease|TestP24GateEvidenceAndEnvelopeKeyedLeafCorpus|TestP24GateRunCommittedDenialPersistsNamedFields|TestP24GateRunPersistsNamedFieldsThroughReceipt|TestP24HumanOutcomeRejectsReferenceDigestMismatch|TestP24IssueCreateConflictNamesChangedTitle|TestP24ObservingJudgmentLeavesNameChangedPublicField|TestP24ObservingJudgmentReplayNeverObserves|TestP24ObservingJudgmentRetriesAfterSourceAndHEADChange|TestP24ReviewOptionalPresenceAndEmptyException|TestP24RunStartConflictNamesChangedClaimVersion|TestP24StopRaiseConflictNamesChangedDetails)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	62.169s` (terminated `exit_code: 0`)
>
> - `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
>
> The build result was observed directly in this session. The real-entrypoint wrapper python3 /tmp/task6-runtime.py invokes /tmp/fid with fresh Git repositories and actual activation, Gate Registration, Envelope, export/import, and identity-free Claim acquire inputs. It checks subprocess return codes: imported acquire/renew/release/expire exit 0; copied-destination renew exit 2 with invalid-argument, invalid claim request and the accepted closed details; replay/copy refusal preserve DB and every journal file. It renews and releases the same imported Claim, and expires a second Claim. The wrapper's printed PASS and observed exit_code 0 are recorded above. Its three preceding failed attempts are also retained: Phase 24 activation is not accepted by the existing Phase 21 activation surface, and export import has not adopted --project. Correcting the scratch driver to the actual existing interfaces changed no product behavior.
>
> The first stored-authority test attempt exposed a one-journal-file helper assumption after import. The next red (before any product change) failed renew, release and expire specifically in imported destinations with the existing invalid-argument fallback. The identical command passed after the one-line shared-predicate repair. The expanded race selection then passed, including same renewed-Claim release.
>
> The touched-package tier ended FAIL at execution-transition copied history plus Artifact fixture snapshot mismatches. Root-local retained snapshots repaired the in-scope Artifact rows; the three Artifact tests and the expanded binding/Artifact race selection passed. The tagged 37-test race selection initially panicked at p23ConflictInput on the removed completion repository; changing its conflict variation to expected_revision preserved its same-key conflict/no-Event subject. Both the focused reproduction and byte-identical 37-test rerun then passed. Package times in the failed touched-package tier: cmd/fid 326.779s, evidence 2.349s, workflow 110.402s, claims 1.055s, export 149.226s.
>
> Complete previous report, retained verbatim with one additional quotation prefix per physical line:
>
> > # Phase 24 - Task 6: Gate registration, Envelope and Claim leaves - Recovery Checkpoint
> >
> > ## Workflow outcome
> >
> > - Task outcome: RECOVERY_REQUIRED
> > - Category: in-scope-implementation
> > - Condition: Focused Task 6 binding, Source, replay, import and eight-leaf descriptor checks pass, but both full test gates terminate with retained fixtures that Task 6 displaced and whose definitions are outside the exact active recovery allowlist.
> > - Recovery attempts: 1
> > - Next action: Phase Runner supplies a bounded same-task scope amendment for the 21 fixture-definition paths listed below; continue causal convergence without weakening their subjects, rerun both failed commands byte-identically within the exact five-command sequence, finalize the report, and require OUTCOME COMPLETE.
> > - Human attention required: no
> >
> > OUTCOME RECOVERY_REQUIRED
> >
> > ## Summary of what was implemented
> >
> > - Continued the existing implementation checkpoint `3ed4824b3f9f7d5dee0f2a50799df538a2911953` and finalized-report HEAD `c845cee5742e20b5418dd30a2780787bf4116a56`. The Runner-confirmed recovery identity, trigger and no-progress count remain unchanged; this continuation made material progress.
> > - Made Obligation Registration nullable as required for the existing typed review/Human Obligation shape. Removed the unused pre-descriptor Claim decoder. Claim Events now preserve the Envelope repository's base_oid and worktree_id.
> > - Converged allowlisted Gate, Envelope, Claim and shared integration/fake helpers to accepted identity-free inputs and real Source bytes. Kept the shared same-worktree predicate, domain reasons, fixed fallback messages, and later Run/Gate-run/Evidence interfaces intact.
> > - Added independently observed Registration digest/record/Event checks; all three Envelope Source positions; original-worktree success, observable whole-worktree plain-copy refusal with no Event, and real export/import destination acceptance. Claim records the local worktree while its Event retains the Envelope scope.
> > - Added replay after HEAD changes, non-Git and unmerged observations, and Source changes/deletion. Every same-key retry retains original result bytes and appends no Event; otherwise viable distinct-key controls refuse with no Event, then succeed when the unavailable state is restored.
> > - Completed the eight-leaf descriptor corpus through the built binary: examples, recursive help, closed typed templates, exclusivity, nested unknown/duplicate/missing/null/type/enum/format failures, integer/string/array bounds, sorted-unique constraints, and unchanged project persistence. The derived fake/fid Task 6 intersection is envelope.create and claim.acquire.
> >
> > ## Files created or modified
> >
> > - `cmd/fid/claim.go`
> > - `cmd/fid/envelope_test.go`
> > - `cmd/fid/gate_test.go`
> > - `cmd/fid/generic_claim_test.go`
> > - `cmd/fid/input_descriptor_test.go`
> > - `cmd/fid/input_descriptors_binding.go`
> > - `cmd/fid/integration/p21_conformance_test.go`
> > - `cmd/fid/integration/p22_fake_test.go`
> > - `cmd/fid/integration/p23_amended_claim_oracle_test.go`
> > - `cmd/fid/integration/p23_amended_envelope_oracle_test.go`
> > - `cmd/fid/integration/p23_conformance_test.go`
> > - `cmd/fid/integration/p23_review_record_test.go`
> > - `internal/claims/generic.go`
> > - `project-docs/specs/06_DESIGN_PLAN_CLI.md`
> > - `project-docs/specs/07_DESIGN_PLAN_FLOWS.md`
> > - `project-docs/docs/phases/phase-24/tasks/task-6-done.md`
> >
> > Earlier implementation paths remain part of the same task's history; this continuation changed only the paths above. Runner-owned ledger dirt is excluded.
> >
> > ## Documentation updated
> >
> > - `project-docs/specs/06_DESIGN_PLAN_CLI.md`: nullable Obligation Registration, exact-root imported-destination equivalence, plain-copy refusal/no Event, and local Claim worktree versus Envelope Event scope.
> > - `project-docs/specs/07_DESIGN_PLAN_FLOWS.md`: imported Envelope/Claim behavior and destination Source observations.
> > - This existing report preserves the original failed evidence and records this session's terminated commands and remaining scope gaps.
> >
> > ## Validation Contract Coverage
> >
> > These are focused Task 6 slice results, not a passing full-task closure. All test names below executed in this session; the final full gates still fail at the listed retained fixture boundaries.
> >
> > - VAL-009: PASS (Task 6 slice) - `TestP24BindingReplayBeforeObservation` (`cmd/fid/envelope_test.go:90`) covers all four observing leaves after HEAD advancement; `TestP24EnvelopeReplayAfterSourceChangeOrDeletion` (`cmd/fid/envelope_test.go:24`) covers omitted digest replay after both byte changes and deletion, with original result equality/no Event.
> > - VAL-011: PASS (eight leaves) - `TestP24GateEnvelopeClaimDescriptors` and `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` (`cmd/fid/input_descriptor_test.go:235`) executed in the ordinary gate; the latter includes the seeded mismatch refusal. `TestP24BindingDescriptorProcessCorpus` (`cmd/fid/input_descriptor_test.go:317`) requires exactly the eight binding descriptors and checks their field trees and typed examples.
> > - VAL-012, VAL-013, VAL-014: PASS (eight leaves) - `TestP24BindingDescriptorProcessCorpus` compares each recursive help field's path/type/requirement/meaning/format/enum/bounds; each real-process example's exact order and two-space formatting; and closed Go decode plus descriptor types. Optional and nullable nested objects are expanded for negative structural rows.
> > - VAL-015: PASS (eight leaves) - `TestAdoptedExampleProcessExclusivity` (`cmd/fid/input_descriptor_test.go:44`) includes all eight among 20 adopted leaves, with six forbidden combinations each. Examples alone succeed in the process corpus.
> > - VAL-016: PASS (eight leaves) - `TestP24BindingDescriptorProcessCorpus` derives structural rows from every descriptor field, including nested objects, arrays and scalar elements, and checks exact field path, closed details and no Event.
> > - VAL-017: PASS (Task 6 slice) - `TestP24RemovedBindingMembers` (`cmd/fid/input_descriptor_test.go:291`) rejects the four removed top-level members; `TestP24EnvelopeSourcePositions` (`cmd/fid/envelope_test.go:327`) rejects removed Source repository in every position with no Event.
> > - VAL-018, VAL-019: PASS (Envelope positions) - `TestP24EnvelopeSourcePositions` requires exactly risk_result.source, sources[0], obligations[0].source. Matching and omitted digests store SHA-256 of real bytes; mismatched and missing-file rows require the exact mismatch details and unchanged persistence.
> > - VAL-022: PASS (Task 6 slice) - `TestGateRegisterObservesRepositoryWithoutCallerIdentity` (`cmd/fid/gate_test.go:209`) independently checks Registration and Event; `TestP24EnvelopeSourcePositions` compares all recorded repositories; `TestP24BindingReplayBeforeObservation` checks successful complete Event parity; `TestP24EnvelopeAndClaimOriginalCopyAndImport` (`cmd/fid/envelope_test.go:234`) checks local Claim identity and Envelope-scope Event parity in both original and imported roots.
> > - VAL-023: PASS (Task 6 slice) - `TestP24BindingReplayBeforeObservation` covers normal acceptance and both non-Git/unmerged refusal for all four observing leaves, including repository-unavailable reason/detail and no Event. The phase-wide classifier's other cause classes remain Task 4 evidence.
> > - VAL-026: PASS - `TestGateRegisterObservesRepositoryWithoutCallerIdentity` computes project ID from config, worktree ID from root SHA-256, HEAD/tree/dirty from independent Git commands, and canonical binding digest from that repository plus scope, then compares the Registration and Event.
> > - VAL-027, VAL-028, VAL-029: PASS (Task 6 slice) - `TestP24EnvelopeAndClaimOriginalCopyAndImport` pairs each identical original-worktree request with an observable whole-copy refusal and a real exported/imported destination acceptance. It checks the Task 6 fixed messages, exit 2, closed details and unchanged persistence for both refused leaves. Task 10's later domain-message adoption is not claimed.
> > - VAL-035: PASS (Task 6 derived intersection) - `TestP23SharedAmendedLeafConformance` (`cmd/fid/integration/p23_conformance_test.go:379`) derives the changed binding leaves from the fake dispatch and requires envelope.create/claim.acquire parity. It covers valid identity-free inputs and removed scope/Source repository/worktree members on both adapters, with no Event for refusals. `TestP22DualAdapterRecoveryConformance` also passes with the converged real-Source shared builder.
> > - VAL-052: PASS (measured execution budgets only) - final ordinary cmd/fid 197.511s; final integration 285.231s; slowest fresh other package in the final ordinary attempt internal/export 62.311s. internal/workflow is cached there; this session's earlier complete ordinary attempt measured 102.885s, and its slowest other package internal/issues 142.117s. All are below 600s ordinary/900s tagged budgets and below warning thresholds. Correctness gates remain failed.
> > - VAL-060: PASS (eight leaves) - `TestP24BindingDescriptorProcessCorpus` runs all eight examples through the built binary in a fresh non-project directory, requires exit 0/empty stderr, and compares the directory tree before/after.
> > - VAL-066: PASS (Task 6 slice) - `TestP24BindingReplayBeforeObservation` and `TestP24EnvelopeReplayAfterSourceChangeOrDeletion` require no-Event same-key replay in unavailable states, otherwise viable distinct-key refusal with no Event, and restored-state new-key success for both repository and each deleted Source position.
> >
> > ## Tests run
> >
> > Working directory: `/home/mark/projects/fid`, except the explicitly fresh non-project CLI directory. Every new result below was observed terminating in this session. The quoted original report at the end of this section is retained historical evidence, not a fresh attestation.
> >
> > - `go test ./cmd/fid -run '^TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go test ./cmd/fid -run '^TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.777s` (terminated `exit_code: 0`)
> >
> > The nullable Registration red preceded the descriptor change; the identical command then passed.
> >
> > - `go test ./cmd/fid -run 'TestGateRegister|TestEnvelope|TestClaim' -count=1 -timeout 10m` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go test ./cmd/fid -run 'TestGateRegister|TestEnvelope|TestClaim' -count=1 -timeout 10m` - FAIL - `FAIL` (terminated `exit_code: 1`)
> >
> > The second identical run reached the two retained completion tests in the out-of-allowlist task18_completion_resume_test.go; earlier allowed nullable/binding fixtures were repaired.
> >
> > - `go test ./cmd/fid -run '^TestP24EnvelopeAndClaimOriginalCopyAndImport$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go test ./cmd/fid -run '^TestP24EnvelopeAndClaimOriginalCopyAndImport$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.062s` (terminated `exit_code: 0`)
> > - `go test ./cmd/fid -run '^TestP24BindingDescriptorProcessCorpus$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go test ./cmd/fid -run '^TestP24BindingDescriptorProcessCorpus$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go test ./cmd/fid -run '^TestP24BindingDescriptorProcessCorpus$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go test ./cmd/fid -run '^TestP24BindingDescriptorProcessCorpus$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.572s` (terminated `exit_code: 0`)
> > - `go test ./cmd/fid -run '^TestP24BindingReplayBeforeObservation$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go test ./cmd/fid -run '^TestP24BindingReplayBeforeObservation$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.200s` (terminated `exit_code: 0`)
> > - `go test ./cmd/fid -run '^TestReviewRecordCreatesImmutableReview$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go test ./cmd/fid -run '^TestReviewRecordCreatesImmutableReview$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.166s` (terminated `exit_code: 0`)
> > - `go test ./cmd/fid -run '^(TestP24BindingDescriptorProcessCorpus|TestP24EnvelopeSourcePositions|TestP24EnvelopeReplayAfterSourceChangeOrDeletion)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go test ./cmd/fid -run '^(TestP24BindingDescriptorProcessCorpus|TestP24EnvelopeSourcePositions|TestP24EnvelopeReplayAfterSourceChangeOrDeletion)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.483s` (terminated `exit_code: 0`)
> >
> > The expanded corpus initially serialized a non-string invalid value with fmt instead of canonical JSON; the assertion was corrected while preserving its required exact error value. The large failed output was not rerun to obtain a different failure.
> >
> > - `go test -tags integration ./cmd/fid/integration -run '^TestPhase21RealProcessJourneyStartsOneActivatedEnvelope$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go test -tags integration ./cmd/fid/integration -run '^TestPhase21RealProcessJourneyStartsOneActivatedEnvelope$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go test -tags integration ./cmd/fid/integration -run '^TestPhase21RealProcessJourneyStartsOneActivatedEnvelope$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go test -tags integration ./cmd/fid/integration -run '^TestPhase21RealProcessJourneyStartsOneActivatedEnvelope$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.969s` (terminated `exit_code: 0`)
> > - `go test -tags integration ./cmd/fid/integration -run '^TestP22DualAdapterRecoveryConformance$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.688s` (terminated `exit_code: 0`)
> > - `go test -tags integration ./cmd/fid/integration -run '^TestP23SharedAmendedLeafConformance$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go test -tags integration ./cmd/fid/integration -run '^TestP23SharedAmendedLeafConformance$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go test -tags integration ./cmd/fid/integration -run '^TestP23SharedAmendedLeafConformance$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.114s` (terminated `exit_code: 0`)
> > - `go test -tags integration ./cmd/fid/integration -run '^TestP23ReviewRecord' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go test -tags integration ./cmd/fid/integration -run '^TestP23ReviewRecord' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.106s` (terminated `exit_code: 0`)
> > - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/evidence ./internal/workflow ./internal/claims` - FAIL - `FAIL` (terminated `exit_code: 1`)
> >
> > Touched-package tier: cmd/fid failed at remaining fixture boundaries (203.426s); evidence 2.251s, workflow 73.468s and claims 1.033s passed.
> >
> > - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestPhase21RealProcessJourneyStartsOneActivatedEnvelope|TestP22DualAdapterRecoveryConformance|TestP23SharedAmendedLeafConformance|TestP23ReviewRecordRealProcessJourney)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	9.181s` (terminated `exit_code: 0`)
> > - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^(TestP24BindingDescriptorProcessCorpus|TestP24BindingReplayBeforeObservation|TestP24EnvelopeSourcePositions|TestP24EnvelopeAndClaimOriginalCopyAndImport|TestP24EnvelopeReplayAfterSourceChangeOrDeletion|TestP24RemovedBindingMembers|TestAdoptedExampleProcessExclusivity|TestClaimAcquireAcceptsReviewMode|TestClaimAcquireKeepsOneActiveClaimAcrossModes|TestClaimAcquireBindingMatrixRejectsWithoutMutation|TestGateRegisterObservesRepositoryWithoutCallerIdentity|TestReviewRecordCreatesImmutableReview)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	27.952s` (terminated `exit_code: 0`)
> > - `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> >
> > Real entrypoint interaction: `/tmp/fid --json claim acquire --actor actor-builder --idempotency-key task6-runtime --input -` ran in a fresh non-project temporary directory with an input containing worktree_id. The child terminated with observed returncode 2, empty stdout and the exact invalid-argument path input.worktree_id; allowed keys were envelope_id, envelope_revision, issue_id, mode, scope, ttl_seconds. A subprocess wrapper checked the closed JSON error and unchanged empty directory, printed `PASS real CLI removed-member refusal: input.worktree_id; no workspace writes`, and terminated with observed exit_code 0.
> >
> > Original checkpoint/report evidence, retained verbatim with quotation prefixes:
> >
> > > # Phase 24 - Task 6: Gate registration, Envelope and Claim leaves - Recovery Checkpoint
> > >
> > > ## Workflow outcome
> > >
> > > - Task outcome: RECOVERY_REQUIRED
> > > - Category: in-scope-implementation
> > > - Condition: The eight leaves have partial observation and descriptor implementation, but retained Gate/Envelope/Claim fixtures are red and Task 6 parity, refusal, replay, import and Source matrices remain unfinished.
> > > - Recovery attempts: 0
> > > - Next action: Resume Task 6 on this checkpoint. Finish nullable Obligation registration input, converge retained tests and shared integration/fake fixtures with introducing-commit provenance, prove all VAL rows and real CLI behavior, run touched-package race tests, then run all five exact commands in order and complete the report.
> > > - Human attention required: no
> > >
> > > OUTCOME RECOVERY_REQUIRED
> > >
> > > ## Summary of what was implemented
> > >
> > > - `gate.register` now computes its recorded project/worktree and binding digest from a repository observation after replay. `claim.acquire` derives its worktree from an observation after replay. `envelope.create` observes after replay, resolves all three Source positions, and checks Registration binding through a shared same-worktree/import predicate. `envelope.complete` drops its caller repository input.
> > > - Adopted descriptors and `--example` for the eight Gate Registration, Envelope and Claim leaves, and added initial focused tests. The implementation and tests are unfinished; this is a checkpoint, not Task 6 acceptance.
> > >
> > > ## Files created or modified
> > >
> > > - `internal/evidence/registration.go`, `internal/workflow/envelope.go`, `internal/workflow/import_binding.go`, `internal/claims/generic.go`
> > > - `cmd/fid/gate.go`, `cmd/fid/envelope.go`, `cmd/fid/claim.go`, `cmd/fid/input_descriptor_command.go`, `cmd/fid/input_descriptors_binding.go`
> > > - `cmd/fid/gate_test.go`, `cmd/fid/envelope_test.go`, `cmd/fid/generic_claim_test.go`, `cmd/fid/input_descriptor_test.go`
> > > - `project-docs/specs/06_DESIGN_PLAN_CLI.md`, `project-docs/specs/07_DESIGN_PLAN_FLOWS.md`, this report
> > >
> > > ## Documentation updated
> > >
> > > - `project-docs/specs/06_DESIGN_PLAN_CLI.md` records the current input shapes and observed bindings.
> > > - `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` records the Gate, Envelope and Claim observation sequence.
> > >
> > > ## Validation Contract Coverage
> > >
> > > - VAL-011..016, VAL-060: PARTIAL - `TestP24GateEnvelopeClaimDescriptors` proves descriptor/type parity and `TestAdoptedExampleProcessExclusivity` passes for adopted leaves; full eight-leaf help, template, structural and side-effect corpus remains.
> > > - VAL-022, VAL-026..028: PARTIAL - `TestGateRegisterObservesRepositoryWithoutCallerIdentity`, `TestEnvelopeCreateObservesScopeAndSourcesWithoutCallerRepository`, and `TestClaimAcquireBindsActiveEnvelopeActorWorktreeScopeAndTTL` pass for original-workspace records. Independent full Repository/Event parity and plain-copy rejection remain.
> > > - VAL-009, VAL-017..019, VAL-023, VAL-029, VAL-035, VAL-066: PENDING - HEAD/replay, removed members, all Source positions, unavailable repositories, imported destinations and fake/fid paired cases still need Task 6 tests and convergence.
> > > - VAL-052: PENDING - final ordinary and tagged gate package timings were not measured on a green tree.
> > >
> > > ## Tests run
> > >
> > > - First red: `go test ./cmd/fid -run '^TestClaimAcquireBindsActiveEnvelopeActorWorktreeScopeAndTTL$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); omitted `worktree_id` returned the old `invalid claim request` fallback.
> > > - After the Claim implementation, the identical focused command - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.634s` (terminated `exit_code: 0`).
> > > - `go test ./cmd/fid -run '^TestGateRegisterObservesRepositoryWithoutCallerIdentity$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); the old binding validator required caller identity. After the Registration change, the identical command - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.475s` (terminated `exit_code: 0`).
> > > - `go test ./cmd/fid -run '^TestP24GateEnvelopeClaimDescriptors$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); all eight descriptors were behind the pending fence. After adopting them, the identical command - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.151s` (terminated `exit_code: 0`).
> > > - `go test ./cmd/fid -run '^TestEnvelopeCreateObservesScopeAndSourcesWithoutCallerRepository$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); the old Envelope input required `repository_scope`. After the Source/observation change, the identical command - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.556s` (terminated `exit_code: 0`).
> > > - `go test ./cmd/fid -run 'TestGateRegister|TestEnvelope|TestClaim' -count=1 -timeout 10m` - FAIL - `FAIL` (terminated `exit_code: 1`); retained fixtures still send removed members, copy path-bound facts, or need nullable `registration`. Preserve this failure and rerun the identical command after repair.
> > > - `go test ./cmd/fid -run '^(TestClaimAcquireBindsActiveEnvelopeActorWorktreeScopeAndTTL|TestGateRegisterObservesRepositoryWithoutCallerIdentity|TestEnvelopeCreateObservesScopeAndSourcesWithoutCallerRepository|TestP24GateEnvelopeClaimDescriptors|TestInputDescriptorRegistryMatchesLiveTreeAndTypes)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.676s` (terminated `exit_code: 0`)
> > > - `go test ./cmd/fid -run '^TestAdoptedExampleProcessExclusivity$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.577s` (terminated `exit_code: 0`)
> > > - `go test ./cmd/fid ./internal/workflow ./internal/evidence ./internal/claims ./internal/export -run '^$'` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	0.005s [no tests to run]` (terminated `exit_code: 0`)
> > > - `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> > > - `/tmp/fid gate register --example` - PASS - `}` (terminated `exit_code: 0`)
> > >
> > > The compile-only command ran no tests. The real `fid --example` invocation ran from `/tmp` and printed the new identity-free input template.
> > >
> > > ## Exact quality-gate results
> > >
> > > Working directory for the required gate: `/home/mark/projects/fid`. None of its commands has been run on a final tree. These are pending commands, not passing evidence.
> > >
> > > - `gofmt -l cmd internal` - BLOCKED - not run; Task 6 retained tests and implementation remain unfinished (no termination).
> > > - `golangci-lint run ./...` - BLOCKED - not run; Task 6 retained tests and implementation remain unfinished (no termination).
> > > - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - not run; retained Gate/Envelope/Claim tests are red (no termination).
> > > - `go tool cover -func=coverage.out` - BLOCKED - not run; no final coverage profile exists (no termination).
> > > - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - not run; Task 6 integration convergence is unfinished (no termination).
> > >
> > > ## Deviations from spec
> > >
> > > - No protected contract or gate definition was changed. The retained focused group ended `FAIL` (terminated `exit_code: 1`) because its inputs still contain removed repository/worktree members, its path-bound fixture caches need replacement, and the descriptor rejects a previously valid nullable Obligation Registration. In-scope convergence is pending; the identical command must be rerun after repair. No required gate attempt has been made.
> > > - The first report-lint check ended `FAIL report-lint found 5 failure(s)` (terminated `exit_code: 5`) because five passing test/build rows had prose or punctuation after their termination field. The report correction moved that prose to a separate paragraph and left the observed command tokens and statuses intact.
> > >
> > > ## Concerns, technical debt, or notes for the Tech Lead
> > >
> > > - Retained files already changed: `gate_test.go` was introduced by `0c4c04ea` (Phase 21 Gate task), `envelope_test.go` by `56dad8fa` (Phase 21 Envelope task), `generic_claim_test.go` by `0b67c830` (Phase 21 Claims task), and `input_descriptor_test.go` by `04b6b099` (Phase 24 task 3). The accepted Phase 24 Task 6 input and observation changes, rather than historical names, require convergence. Record introducing commits for any further retained fixture changes before editing.
> > > - Preserve the negative binding matrix while replacing caller `worktree_id` manipulation with a real plain-copy worktree. Do not remove the wrong-worktree assertion. Preserve the original-workspace positive control, and add imported-destination controls and no-Event checks.
> > > - The current descriptor's `registration` field needs the accepted nullable shape for review/Human Obligations. Its retained failure is real; changing those inputs to fake Registration objects would weaken coverage.
> > > - The current worktree has Runner-owned `project-docs/docs/phases/phase-24/triage-log.md` and untracked `run-ledger.md` changes. Leave them untouched during continuation and commits.
> > >
> > > ## Dependency-budget note
> > >
> > > - No dependency changes.
> > >
> > > ## Commit hash
> > >
> > > - `3ed4824b3f9f7d5dee0f2a50799df538a2911953`
> > >
> > > ## Commit subject
> > >
> > > - `feat(binding): checkpoint observed gate envelope claim inputs`
> >
> > ## Exact quality-gate results
> >
> > Working directory for both complete five-command attempts: `/home/mark/projects/fid`. Each command was run in the specified order and watched terminate. No command was killed, timed out, substituted or silently skipped. Coverage reporting succeeded over a failed-suite profile, which does not establish task correctness.
> >
> > Initial prefix before the first complete sequence:
> >
> > - `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> > - `golangci-lint run ./...` - FAIL - `* unused: 1` (terminated `exit_code: 1`)
> >
> > The lint failure was the now-unused decodeClosedClaimInput. Its removal stayed inside the allowlist.
> >
> > First complete sequence after that repair:
> >
> > - `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> > - `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
> > - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`)
> > - `go tool cover -func=coverage.out` - PASS - `total:											(statements)					59.6%` (terminated `exit_code: 0`)
> > - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
> >
> > Observed package times: cmd/fid 234.663s, internal/workflow 102.885s, slowest other ordinary package internal/issues 142.117s, tagged integration 273.774s. The real-Source digest in the shared Envelope helper was subsequently corrected; the Review positive control and expanded Task 6 race corpus passed.
> >
> > Final complete sequence on the checkpoint tree:
> >
> > - `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> > - `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
> > - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
> > - `go tool cover -func=coverage.out` - PASS - `total:											(statements)					59.8%` (terminated `exit_code: 0`)
> > - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
> >
> > Both BLOCKED rows describe commands that actually ran and failed, not unavailable execution. Closure is blocked by the exact fixture-definition scope gaps below. The ordinary suite ended FAIL after cmd/fid 197.511s and internal/export 62.311s; the tagged suite ended FAIL after 285.231s. No green full gate is claimed. Report-lint must return OUTCOME RECOVERY_REQUIRED for this durable checkpoint; OUTCOME COMPLETE requires the later exact gates to pass.
> >
> > ## Deviations from spec
> >
> > - The draft-report writing operation duplicated later sections because JavaScript replacement text interpreted a literal dollar-apostrophe from a test command. The immediately following report-finalization commit replaces that duplicated draft with one copy of each required section and preserves every unique command result; implementation and observed test evidence are unchanged.
> >
> > - Original failed focused/report-lint evidence is preserved in the quotation above. Historical not-run gate rows have been superseded by two observed complete five-command sequences, never rewritten as historical passes.
> > - This session's first lint attempt ended `* unused: 1` (terminated `exit_code: 1`); the obsolete decoder was removed and both later identical lint attempts ended `0 issues.` (terminated `exit_code: 0`).
> > - The nullable Obligation red ended `FAIL` (terminated `exit_code: 1`); minimum descriptor repair accepted the already specified null Registration, and the identical focused command ended `ok  	github.com/mbell-nzl/fid/cmd/fid	0.777s` (terminated `exit_code: 0`).
> > - The original/copy/import red ended `FAIL` (terminated `exit_code: 1`); it exposed missing Claim Event base_oid/worktree_id and a test helper assuming one journal file. Event construction was corrected and the new independent helper reads all journal files; the identical command ended `ok  	github.com/mbell-nzl/fid/cmd/fid	1.062s` (terminated `exit_code: 0`).
> > - Initial new-corpus fixture failures ended `FAIL` (terminated `exit_code: 1`) for help rendering assumptions, duplicate-key insertion, placeholder normalization and control-fixture assertion reuse. Each helper was corrected to exercise the specified public behavior; the identical specific commands passed before further implementation. The expanded array-value assertion correction and its identical rerun are recorded above.
> > - Retained integration red attempts ended `FAIL` (terminated `exit_code: 1`) for obsolete input members and fake observations; the p21 repair additionally exposed a nonexistent helper reference and missing Source directory, and the p23 fake repair additionally exposed typed-null/unknown-member parity. Allowed shared fixtures now use independent observations, real Source bytes and the accepted closed member lists; identical focused commands passed as recorded above.
> > - The first ordinary and tagged required commands each ended `FAIL` (terminated `exit_code: 1`). The ordinary failure exposed the shared Envelope helper's missing digest when retained Review converted its Source to external-reference; supplying the actual file digest preserved that positive assertion, and explicitly omitted/missing Source rows retained their omission coverage. The byte-identical final ordinary and tagged reruns each still ended `FAIL` (terminated `exit_code: 1`) because remaining definitions are outside the active allowlist. Their inputs were displaced by this task's accepted identity/observation contract; no ownership is assigned to a different task.
> > - No protected command, assertion semantics, threshold, coverage policy, domain reason/message, same-worktree predicate or later Run/Gate-run/Evidence interface changed. Task closure remains unfinished solely under the current bounded authority; do not approve this report as COMPLETE.
> >
> > ## Concerns, technical debt, or notes for the Tech Lead
> >
> > - Recovery identity retained: attempt 1; no-progress count 0/2; trigger `task-6:session-continuation:retained-gate-envelope-claim-fixtures-and-nullable-registration`; restored capability authority `restored-capability:developer-role-spawn-01a10b08-efe6-7050-8f88-65c524e62519`; governing TRG-24-003/ESC-24-001. Runner-owned files were neither edited nor committed. This report does not select a new recovery identity or increment Runner counters.
> > - In-scope retained provenance was checked by introducing-commit history, not inferred from filenames: Gate finalization `0c4c04ea` (Phase 21/task 5); Envelope factory `56dad8fa0dcb9b79f7df2a31b356f1b2eb8a694e` (21/6), repository helper `2339dc9ebd8d0e6b5502ca32e1562892cf704ce6` (22/12); Claim lifecycle `0b67c8304a84e4e5bcb8df44d79f899e40a16acb` and binding matrix `59b079a5a271cc65f7964a471b6383b9c2dd6677` (21/10), review mode `35944c72aba0f2c0b918bdb494b9fa9487d27142` (23/3); descriptor corpus `04b6b0992fe20a02df402a183269a9847ddae5ab` (24/3).
> > - Shared integration provenance: p21 corpus `a99b4270f70e98111b5bae43c54ed92aed7078c6` and journey exposure `236540c48a2129d70480b3222ba4122a5471eb48` (21/22); p22 fake builder `f3be616f0587881ee3f5c3716ffbf6167e3f2a20` (22/6); amended Envelope/Claim oracles `2e7fd2f9fa05bce648510789453c8d2ed8633b1d` and shared p23 conformance `333514946aeab8328fa618f0f619fdce84694f9a` (23/18); p23AcquireClaim `d215e40c3fbb382d9f0bf17d54a76d50a2e76715` (23/4).
> > - Remaining fixture definitions below are outside the exact active recovery allowlist. This task's accepted removal of caller identities, real Source resolution and worktree observation caused these failures. A bounded Task 6 amendment should converge these definitions while preserving positive/negative, Stop, Source, import, replay and no-Event assertions. Cache repair must occur before recording path-bound history; creating fictitious import bindings or relaxing the shared predicate would change the protected behavior.
> >
> > | Needed fixture-definition path | Current failure boundary | Introducing commit | Phase/task |
> > |---|---|---|---|
> > | `cmd/fid/generic_run_test.go` | runEnvelopeForTestWithDefinition copies cached observed Envelope history to another root | `d17f524ce3d7bdf7de4e715a03514cba2fe8d7a8` | 21/22 |
> > | `cmd/fid/evidence_test.go` | readyEvidenceEnvelopeForTest copies cached observed Envelope/Claim/Run facts to another root | `f3350bcf582aacc328c838cd4f08c8b8d9e30f7c` | 21/22 |
> > | `cmd/fid/review_test.go` | acquireReviewClaim supplies removed worktree_id | `a1ae1bb61363a591bb46c9ddfd4b8df36462687b` | 23/4 |
> > | `cmd/fid/task18_completion_resume_test.go` | completion inputs supply removed repository | `2313235` | 21/18 |
> > | `cmd/fid/task21_maintenance_rebuild_test.go` | completion input supplies removed repository | `8435120ac4cc31782c68a1561498a47a4e9fe71b` | 21/21 |
> > | `cmd/fid/stop_evaluation_matrix_test.go` | both completion Stop states supply removed repository | `4b690e8b863abac96c9abfd72709ac36db2235b2` | 22/9 |
> > | `internal/export/archive_fixture_test.go` | fixtureSource supplies a placeholder digest for actual project-file bytes | `1be8e882ef689f15705507c448e3b087e23d5f65` | 22/13 |
> > | `cmd/fid/integration/p22_resume_actions_test.go` | p22NewResumeActionAdapter supplies removed Gate/Envelope/Claim identities | `98020edaa18a744dc868121b304ec527acc303d9` | 22/12 |
> > | `cmd/fid/integration/p22_recovery_journey_test.go` | recovery journey supplies removed Gate/Envelope/Claim identities | `c69456fc9bf42f811e4999f1968a48c4a82d3cc5` | 22/6 |
> > | `cmd/fid/integration/p23_review_evidence_budget_test.go` | conflicting Claim input supplies removed worktree_id | `d05c288973df5b71d29dd8cc812374072f0736fd` | 23/6 |
> > | `cmd/fid/integration/p23_obligation_guard_test.go` | completion helper supplies removed repository | `613e45aad9b10445cffd6b113ee598ab493f183a` | 23/2 |
> > | `cmd/fid/integration/p23_human_outcome_test.go` | completion and activation fixtures supply removed identities and manufactured Source facts | `0ff2eeef66d4633309a43ae6a341aecd38def4d0` | 23/10 |
> > | `cmd/fid/integration/p23_human_fixtures_test.go` | p23ReviewExhaustionFixture supplies removed worktree_id; portable callers fail here | `59d2653da7cfe14699dc4d1cd85ef032bec26a5c` | 23/14 |
> > | `cmd/fid/integration/p23_human_failure_test.go` | Claim Stop/rework rows supply removed worktree_id | `180f5168e3f1338cb758f8388f1c26fd6f8428d3` | 23/19 |
> > | `cmd/fid/integration/p23_gate_grants_test.go` | Gate Registration fixture supplies removed repository/worktree_id | `8531522e9602b0f72717d6c9d889684228cfa2ee` | 23/6 |
> > | `cmd/fid/integration/p23_decision_reference_test.go` | completion input supplies removed repository | `828eb18f422ebbfacbe4e65c62d53a061b7e9b9f` | 23/9 |
> > | `cmd/fid/integration/p23_task13_precedence_test.go` | completion precedence input supplies removed repository | `d1c016daa9f5367a620a3d1908e39aed12d75110` | 23/13 |
> > | `cmd/fid/integration/p23_concern_deviation_test.go` | completion input supplies removed repository | `3a0aed9d35ca2d43540e686a5cbbaa245177ea0f` | 23/8 |
> > | `cmd/fid/integration/p24_gate_evidence_corpus_test.go` | completion corpus input supplies removed repository | `9bbbdb55feb7409c8a845d387bc714d927babd8a` | 24/2 |
> > | `cmd/fid/integration/p24_single_journal_test.go` | completion corpus input supplies removed repository | `95c738f23c34aef1a9e8c6630f1d88c660648591` | 24/2 |
> > | `cmd/fid/integration/p24_idempotency_test.go` | Claim conflict row supplies removed worktree_id | `8bddc91940dd2f0e012eabc389ffed2a7aa5e987` | 24/2 |
> >
> > - Newly inspected concern for the next continuation: `internal/claims/generic.go:380` still compares stored Claim worktree directly to the Envelope worktree on renew/release/expire. A Claim acquired in an imported destination records a different local worktree. The Task 6 import acceptance test proves acquisition only; mutation of that newly acquired imported Claim is not claimed as tested. Resolve the accepted stored-authority contract before changing that check; its input shapes and domain reason must stay unchanged.
> > - There are no measured package-time headroom warnings. Cached workflow timing in the final attempt is explicitly distinguished from this session's earlier fresh measurement.
> >
> > ## Dependency-budget note
> >
> > - No dependency changes.
> >
> > ## Commit hash
> >
> > - `1baaecc7983872b06e2a47e6df0801bfb1370ea0`
> >
> > ## Commit subject
> >
> > - `feat(binding): converge observed envelope and claim recovery`
>
> ## Exact quality-gate results
>
> Working directory: /home/mark/projects/fid. These five byte-identical required commands ran sequentially in the specified order on the final implementation/test tree. No command was substituted, interrupted, killed, silently skipped or run concurrently with other heavy tests. The QA plan's every-task commands are the same five; final-review-only commands are not Developer Task 6 gates.
>
> - `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> - `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
> - `go tool cover -func=coverage.out` - PASS - `total:											(statements)					60.7%` (terminated `exit_code: 0`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
>
> Both BLOCKED rows record commands that executed and failed, each with an observed exit_code 1. Every ordinary failure reaches Claim acquisition through executionTransitionFixture before the intended transition assertion. That cache copies observed Envelope history to a different root, correctly refused under Task 6. The tagged failures are TestP22RecoveryAdoptedTargetPermitsTerminalTransitionWithFreshSuccessorProof at p22_recovery_adopted_target_test.go:123 (input.repository is rejected by envelope.complete), and TestP22RunReplayRetainsHistoricalBudgetSnapshot/run.abandon at p22_run_replay_snapshot_test.go:26 via p22StartAfterAbandon (input.worktree_id is rejected by claim.acquire). Their bodies are outside DEC-283 clause 7. These later failure boundaries are progress after the 21 accepted fixture repairs. Neither full gate is claimed green. Coverage reporting over the failed-suite profile is a successful reporting command only.
>
> Observed Go package wall times: cmd/fid 316.372s; internal/workflow 114.227s; slowest fresh other ordinary package internal/export 144.968s; tagged integration 332.169s. Other cached ordinary packages are explicitly reported as cached by Go; no fresh timing is claimed for them. There are no package-time headroom warnings.
>
> ## Deviations from spec
>
> - The first committed-report lint attempt ended `FAIL report-lint found 1 failure(s)` (terminated `exit_code: 1`). It rejected the actual compiler-diagnostic last line of the earlier failed focused command as a Go test token. That complete failed row is preserved verbatim as quoted terminal evidence; no output token, termination, command, assertion or control was changed. The identical report-lint command must be rerun against this committed correction.
>
> - All original failed evidence is retained in the quotation above. Earlier failed required attempts stay failed. The current ordinary and tagged commands are byte-identical reruns of those retained failed gates, and their independently observed results above supersede no historical evidence.
> - DEC-283 clause 7 authorizes the 21 retained paths; all are causally converged for Task 6's accepted removal of caller identity and real Source observation. The failed retained selections first ended FAIL (terminated exit_code: 1), including obsolete input members, copied observed history, placeholder Source digests, compile errors from fixture edits, missing Source directories, and an optional non-null review_profile passed as null. The fixture corrections use real local setup and accepted optional omission, preserving nullable Registration and each behavioral assertion. Identical selections later passed with their own tokens above.
> - DEC-283 clause 8 red ended FAIL (terminated exit_code: 1) for imported renew/release/expire after successful local acquisition. Replacing the direct stored worktree equality with workflow.SameWorktree preserves the invalid-binding reason, public inputs, fallback code/message/exit and no-Event behavior. The identical focused command later ended ok with 1.549s (terminated exit_code: 0), and the expanded race command ended ok with 59.959s (terminated exit_code: 0).
> - The tagged 37-test red and focused p23 decision-reference red ended FAIL (terminated exit_code: 1), because p23ConflictInput tried to mutate removed repository.head_oid. The accepted expected_revision variation preserves the same-key idempotency conflict and no-state-change assertions. The identical focused rerun ended ok with 5.267s and the identical tagged race rerun ended ok with 62.169s, each terminated exit_code: 0.
> - The pre-gate touched-package command ended FAIL (terminated exit_code: 1); this task's accepted observation displaced cached execution-transition Envelope history and later fixture snapshots. The in-scope Artifact snapshot repair passed its specific command and expanded race command. The required ordinary command nevertheless ended FAIL (terminated exit_code: 1) after 316.372s in cmd/fid, because the execution-transition definition remains outside the exact allowlist. This is a same-Task-6 scope gap, not ownership assigned to any historical numbered task.
> - The byte-identical full tagged rerun ended `FAIL` (terminated `exit_code: 1`) after 332.169s. This Task 6's accepted removal of Envelope completion repository and Claim acquisition worktree_id invalidates the two later retained fixture inputs. Their paths are outside the accepted allowlist, so no repair was attempted there. Preserve the successor-proof and historical Budget replay/no-Event assertions in a bounded same-task continuation. The historical full tagged failures remain failed; the 37-test selection's later green result does not replace this required gate.
> - No accepted input shape beyond the Task 6 contract, domain reason/message, public interface outside Gate/Envelope/Claim, shared predicate semantics, gate command, assertion semantics, coverage policy, dependency or later Run/Gate-run/Evidence/ExpectedRepository behavior changed. No tests were deleted, weakened, re-pointed or mined for a flaky green. The required full ordinary gate must pass after a bounded same-task continuation before closure.
> - The initial staging wrapper terminated with observed exit_code 1 because Git could not create .git/index.lock on the sandbox's read-only Git metadata filesystem. The same 27 allowlisted paths were staged through the authorized direct git add command, which terminated with observed exit_code 0 and no output. The runner-owned ledger was excluded in both attempts.
>
> ## Concerns, technical debt, or notes for the Tech Lead
>
> - All 21 rows below were checked against introducing-commit history and Ref footers, not filenames or prose labels. O3 means the exact ordinary command above; I5 means the exact tagged command above. Both execute the retained rows; focused red/green evidence is additionally recorded. The invalidating change for every row is this Task 6's accepted identity-free input and repository/Source observation contract.
>
> | Converged fixture-definition path | Preserved subject and repair | Introducing commit | Phase/task | Executing required command |
> |---|---|---|---|---|
> | `cmd/fid/generic_run_test.go` | Run binding; build observed Envelope history locally | `d17f524ce3d7bdf7de4e715a03514cba2fe8d7a8` | 21/22 | O3 |
> | `cmd/fid/evidence_test.go` | Evidence trust/freshness/Artifact/no-Event; build local Envelope/Claim/Run setup | `f3350bcf582aacc328c838cd4f08c8b8d9e30f7c` | 21/22 | O3 |
> | `cmd/fid/review_test.go` | Review authority; omit displaced caller worktree | `a1ae1bb61363a591bb46c9ddfd4b8df36462687b` | 23/4 | O3 |
> | `cmd/fid/task18_completion_resume_test.go` | Completion qualification and denial; omit displaced repository | `2313235ba463060538f2f0d4bcfa61f039fcdc13` | 21/18 | O3 |
> | `cmd/fid/task21_maintenance_rebuild_test.go` | Rebuild/export parity; omit displaced completion repository | `8435120ac4cc31782c68a1561498a47a4e9fe71b` | 21/21 | O3 |
> | `cmd/fid/stop_evaluation_matrix_test.go` | Both completion Stop states; omit displaced repository | `4b690e8b863abac96c9abfd72709ac36db2235b2` | 22/9 | O3 |
> | `internal/export/archive_fixture_test.go` | Archive Source and recovery parity; digest actual fixture file bytes | `1be8e882ef689f15705507c448e3b087e23d5f65` | 22/13 | O3 |
> | `cmd/fid/integration/p22_resume_actions_test.go` | Runnable resume operations; local identity-free binding and real Source file | `98020edaa18a744dc868121b304ec527acc303d9` | 22/12 | I5 |
> | `cmd/fid/integration/p22_recovery_journey_test.go` | Recovery freshness journey; local binding and real Source file | `c69456fc9bf42f811e4999f1968a48c4a82d3cc5` | 22/6 | I5 |
> | `cmd/fid/integration/p23_review_evidence_budget_test.go` | Review Evidence budget conflicts; omit displaced Claim worktree | `d05c288973df5b71d29dd8cc812374072f0736fd` | 23/6 | I5 |
> | `cmd/fid/integration/p23_obligation_guard_test.go` | Guard-specific completion denial/success; omit displaced repository | `613e45aad9b10445cffd6b113ee598ab493f183a` | 23/2 | I5 |
> | `cmd/fid/integration/p23_human_outcome_test.go` | Human readiness and activation; local real Source and accepted optional fields | `0ff2eeef66d4633309a43ae6a341aecd38def4d0` | 23/10 | I5 |
> | `cmd/fid/integration/p23_human_fixtures_test.go` | Runnable exhaustion/portable fixtures; omit displaced Claim worktree | `59d2653da7cfe14699dc4d1cd85ef032bec26a5c` | 23/14 | I5 |
> | `cmd/fid/integration/p23_human_failure_test.go` | Claim Stop/rework failure parity; omit displaced worktree | `180f5168e3f1338cb758f8388f1c26fd6f8428d3` | 23/19 | I5 |
> | `cmd/fid/integration/p23_gate_grants_test.go` | Every Gate grant and denial; identity-free binding with real Source input | `8531522e9602b0f72717d6c9d889684228cfa2ee` | 23/6 | I5 |
> | `cmd/fid/integration/p23_decision_reference_test.go` | Protected decision lifecycle; omit displaced completion repository | `828eb18f422ebbfacbe4e65c62d53a061b7e9b9f` | 23/9 | I5 |
> | `cmd/fid/integration/p23_task13_precedence_test.go` | Failure precedence; omit displaced completion repository | `d1c016daa9f5367a620a3d1908e39aed12d75110` | 23/13 | I5 |
> | `cmd/fid/integration/p23_concern_deviation_test.go` | Concern/deviation completion; omit displaced repository | `3a0aed9d35ca2d43540e686a5cbbaa245177ea0f` | 23/8 | I5 |
> | `cmd/fid/integration/p24_gate_evidence_corpus_test.go` | Gate Evidence keyed completion corpus; omit displaced repository | `9bbbdb55feb7409c8a845d387bc714d927babd8a` | 24/2 | I5 |
> | `cmd/fid/integration/p24_single_journal_test.go` | All49 single-journal rebuild corpus; omit displaced completion repository | `95c738f23c34aef1a9e8c6630f1d88c660648591` | 24/2 | I5 |
> | `cmd/fid/integration/p24_idempotency_test.go` | Claim idempotency conflicts; omit displaced worktree | `8bddc91940dd2f0e012eabc389ffed2a7aa5e987` | 24/2 | I5 |
>
> - Additional converged assertion: p23ConflictInput envelope.complete same-key conflict/no-Event, cmd/fid/integration/p23_review_record_test.go:513, introduced by d1c016daa9f5367a620a3d1908e39aed12d75110 (Refs: phase-23/task-13). The removed repository field is replaced by accepted expected_revision as the conflict variation. The asserted conflict code, identity, exit and no Event remain intact.
> - Retained snapshot helper provenance: runRepositoryForTest introduced by 5ca732ff163ee1d31fc155d3d9ebb3dac1501838 (21/11); buildFinalizedEnvelopeWorkspaceWithRepository introduced by 2339dc9ebd8d0e6b5502ca32e1562892cf704ce6 (22/12). Their current local snapshots are independent observations of the using root, preserving later caller-input fixtures rather than changing later runtime behavior.
> - Remaining ordinary failure scope gap: cmd/fid/task12_execution_transition_test.go:320 executionTransitionFixture copies observed Envelope history through cloneReusableWorkspaceFixture. Its cache was introduced by d17f524ce3d7bdf7de4e715a03514cba2fe8d7a8 (Refs: phase-21/task-22). It blocks the qualifying, wrong-worktree, holder/expiry, ordered-reason and full predicate matrix controls before their intended assertions. Task 6's observation caused the invalidation.
> - Dependent path needing the same bounded authority: cmd/fid/execution_transition_fixture_test.go:5 TestExecutionTransitionFixtureDoesNotInheritExpiredClaim, introduced by 730128cd8842e2943c9ab2b21f8f2b154148c10a (Refs: phase-23/task-1). It directly accesses the execution-transition cached template to age it. Preserve the isolation/expired-Claim subject when replacing the invalid copied-history setup; do not remove that regression.
> - Tagged scope gap: cmd/fid/integration/p22_recovery_adopted_target_test.go:123 retains completion repository in its positive successor-proof journey. Introducing commit c2d80980b9bfa9c1c20723cc7751526eb8baf2bf (Refs: phase-22/task-12), checked by git log -S p22-adopted-terminal-complete-b. Task 6's accepted completion input removal displaced this assertion's setup; preserve its fresh successor-proof, prior Receipt identity and successful completion subject.
> - Tagged scope gap: cmd/fid/integration/p22_run_replay_snapshot_test.go p22StartAfterAbandon retains Claim worktree_id. Introducing commit 9301b105bc6b456b90b6c5365997fc5e90830561 (Refs: phase-22/task-12), checked by git log -S after-abandon-claim and file-add provenance. Task 6's accepted Claim input removal displaced this setup; preserve all four replay operations, later Budget charge, exact original/replayed Event and result equality and no-Event assertion.
> - All four paths are absent from the exact 51-path allowlist and remain unchanged. No shared-helper workaround, copied-history rebinding, fictitious import or predicate relaxation was used. A bounded Task 6 continuation can repair them without a protected contract change; this report does not route them to another task.
> - Runner ledger bytes were checked unchanged. Its attempt-4 recovery fields below are exact copies, not normalized or paraphrased. Allowed files was checked byte-for-byte against DEC-283 clause 7 before editing (51 paths). No recovery identity or counter was selected by this Developer.
>
> > - Source: project-docs/docs/phases/phase-24/tasks/task-6-done.md typed RECOVERY_REQUIRED and report-lint OUTCOME RECOVERY_REQUIRED; DEC-283
> > - Category: in-scope-implementation
> > - Repair mode: task-scope
> > - Repair authority: task spec; DEC-283
> > - Trigger signature: task-6:session-liveness:developer-role-capacity-after-bounded-retry
> > - Attempt count: 4
> > - Consecutive no-progress attempts: 0/2
> > - Checkpoint: implementation 1baaecc7983872b06e2a47e6df0801bfb1370ea0; finalized report and ending HEAD 7298964a65beb375a6b50ec45c42ea12dada6c59 with report-lint OUTCOME RECOVERY_REQUIRED; phase HEAD 13164582232b9b43eddac93a0ef75b6f4465783c after DEC-283 verifier maintenance; Runner-owned run-ledger.md is the only dirty path
>
> - No ordinary package reached 480s and no tagged package reached 720s. No dependency changes or protected authority issue were encountered.
>
> ## Dependency-budget note
>
> - No dependency changes.
>
> ## Commit hash
>
> - `b7d1fd56b801d10378b35dff8f45684d1990804d`
>
> ## Commit subject
>
> - `fix(binding): converge task 6 fixtures and imported claims`

## Exact quality-gate results

Working directory: /home/mark/projects/fid. These commands ran sequentially in the exact task/QA-plan order on the final fixture tree. Focused selections and the touched-package tier were additional checks and replaced no required command. Every PASS has a directly observed nested exit_code 0.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					60.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	274.935s` (terminated `exit_code: 0`)

Observed package wall times: cmd/fid 283.342s; slowest other freshly timed ordinary package internal/export 116.978s; cmd/fid/integration 274.935s. internal/workflow was explicitly `(cached)` in this session's exact ordinary gate and emitted no fresh wall time; attempt 4's 114.227s remains historical evidence below, not a fresh attestation. Other cached package rows are also distinguished from fresh timing. The pre-gate cmd/fid race check took 304.689s.

## Deviations from spec

- Historical required failures remain FAIL/BLOCKED evidence in the complete quotation above. In particular, attempt 4's ordinary and tagged gates each ended `FAIL` (terminated `exit_code: 1`), at 316.372s and 332.169s respectively. TRG-24-004 grants the four-path bounded convergence without retroactively widening attempt 4.
- The first focused ordinary failure ended `FAIL` (terminated `exit_code: 1`), caused by copying observed Envelope history into another root, which the accepted Task 6 same-worktree rule correctly rejects. Building Registration/Envelope history in the consuming root and explicitly separating expired-prior from fresh-Claim controls repairs that setup. The identical command later ended `ok  	github.com/mbell-nzl/fid/cmd/fid	18.053s` (terminated `exit_code: 0`). The byte-identical full ordinary gate subsequently passed as recorded above.
- The first focused tagged failure ended `FAIL` (terminated `exit_code: 1`), caused by Task 6's accepted removal of completion repository and acquisition worktree_id. Omitting only those displaced members preserves successor-proof/Receipt/completion and Budget/result/Event/no-Event assertions. The identical command later ended `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.483s` (terminated `exit_code: 0`). The byte-identical full tagged gate subsequently passed as recorded above.
- Provenance was confirmed through the introducing rows in git show and each commit's Refs footer, rather than filenames or embedded task numbers. The invalidating change for all four is this Phase 24 Task 6's accepted observation/input convergence, not a later task's ownership.

| Converged fixture definition | Preserved subject and accepted invalidating change | Introducing commit | Verified Refs |
|---|---|---|---|
| cmd/fid/task12_execution_transition_test.go | Cache displaced by observed root identity; all transition predicate subjects retained with local history | d17f524ce3d7bdf7de4e715a03514cba2fe8d7a8 | phase-21/task-22 |
| cmd/fid/execution_transition_fixture_test.go | Expired-Claim isolation depends on that displaced cache; expired-prior and fresh-root controls retained | 730128cd8842e2943c9ab2b21f8f2b154148c10a | phase-23/task-1 |
| cmd/fid/integration/p22_recovery_adopted_target_test.go | Positive successor completion setup sends removed repository; proof and Receipt identity retained | c2d80980b9bfa9c1c20723cc7751526eb8baf2bf | phase-22/task-12 |
| cmd/fid/integration/p22_run_replay_snapshot_test.go | Acquisition after abandonment sends removed worktree_id; all four replay operations and Budget equality retained | 9301b105bc6b456b90b6c5365997fc5e90830561 | phase-22/task-12 |

- No product or protected contract change was needed. Assertion subjects and strength, failure envelopes, same-worktree/import semantics, later Run/Gate-run/Evidence interfaces, input/Event schemas and gate definitions remain as accepted. Plain-copy negative controls still explicitly reject; no copied-history rebinding, invented import or identity workaround was used.

## Concerns, technical debt, or notes for the Tech Lead

- Fresh independent QA window B after Task 6 remains mandatory. This continuation makes no QA, TL-Lite or Phase Runner decision and does not select a new trigger or counter.
- No measured package reached 80% of its budget: ordinary warning 480s/limit 600s; tagged warning 720s/limit 900s. internal/workflow and other cached packages emitted no fresh duration; their historical timings are retained only as history. No fresh timing is claimed for a cached package.
- At invocation both run-ledger.md and triage-log.md were dirty, although the supplied checkpoint described ledger-only dirt. Both were snapshotted before edits, checked byte-identical, excluded from staging, and preserved. This does not widen Developer authority.
- Attempt-5 ledger Category/Repair category mapping, authority, trigger, attempt, allowlist, forbidden changes, report, exact verification and control base were verified as exact byte-for-byte copies of TRG-24-004. The following latest TASK-RECOVERY values are copied verbatim, not paraphrased:

> - Source: TRG-24-004; project-docs/docs/phases/phase-24/tasks/task-6-done.md; DEC-283
> - Category: in-scope-implementation
> - Repair mode: task-scope
> - Repair authority: task spec; DEC-283
> - Triage disposition: N/A
> - Trigger signature: task-6:session-liveness:developer-role-capacity-after-bounded-retry
> - Attempt count: 5
> - Control contract: N/A
> - Control mismatch: N/A
> - Allowed files: cmd/fid/task12_execution_transition_test.go,cmd/fid/execution_transition_fixture_test.go,cmd/fid/integration/p22_recovery_adopted_target_test.go,cmd/fid/integration/p22_run_replay_snapshot_test.go
> - Forbidden changes: No product files, shared helpers outside the four allowed paths, specs, task specs, phase/current decisions, public-interface or input-inventory documents, QA contracts/plans, dependencies, controls, prompts, templates, generated files or runner-owned run-ledger.md. Do not change Run, gate-run, Evidence or ExpectedRepository behavior; same-worktree/import predicate, domain reasons, failure codes/messages/exits, input or Event schemas, persistence, gate commands, assertion semantics, coverage policy, accepted status/scope, thresholds/budgets/tolerances or policy. Never copy observed path-bound history to another root, rebind copied history, invent an import binding, relax the shared predicate, remove/weaken/narrow/re-point tests or replace required gates with focused selections. A deliberately copied negative fixture must retain explicit rejection under the accepted plain-copy rule, not become accepted through a transport or identity workaround.
> - Originating report: project-docs/docs/phases/phase-24/tasks/task-6-done.md
> - Exact verification: Working directory /home/mark/projects/fid, or the Phase Runner-recorded same-task worktree repository root declared explicitly in the report; sequentially rerun all five byte-identical Task 6 commands: `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...`. Record each command's own output token and observed termination, package wall times and existing VAL-052 budgets/warnings. Phase Runner independently runs `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-6-done.md project-docs/docs/phases/phase-24/tasks/task-6.md` against the committed report; this is a Runner control, not a substituted Developer gate.
> - Control repair base: N/A
> - Consecutive no-progress attempts: 0/2
> - Checkpoint: attempt 4 implementation/report ending HEAD 45961a91b60bbeae7e1747642637544bbcf89fbe with report commit b7d1fd56b801d10378b35dff8f45684d1990804d and report-lint OUTCOME RECOVERY_REQUIRED; TRG-24-004 bounded the four newly exposed fixture definitions; phase HEAD 45961a91b60bbeae7e1747642637544bbcf89fbe; Runner-owned run-ledger.md is the only dirty path
> - Phase worktree: /home/mark/projects/fid on branch phase-24; Runner-owned run-ledger.md is the only dirty path
> - Temporary repair worktree: N/A
> - Hidden-index flags: none; assume-unchanged and skip-worktree are forbidden
> - Resolution evidence: PENDING - retain each four-path assertion's introducing provenance, run focused red before edits and byte-identical focused green after edits, build root-local observed history and expired-Claim isolation, preserve successor-proof/Receipt identity and historical Budget replay/result/no-Event assertions, run all five exact Task 6 commands in order with observed termination, finalize and commit the same Task 6 report with Refs: phase-24/task-6, rerun report-lint, and require OUTCOME COMPLETE before task closure and fresh independent QA window B
> - Next action: dispatch one same-task Developer continuation for Task 6 attempt 5 using only the exact four-file allowlist plus separate completion-report bookkeeping, preserving attempt 4 evidence and the 0/2 no-progress budget; then independently run report-lint
> - Human attention required: no

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
