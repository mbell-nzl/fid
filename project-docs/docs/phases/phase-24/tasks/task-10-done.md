# Phase 24 - Task 10: Human failure lines - Recovery checkpoint

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: RegisteredCodes pairing and R20 domain-process completeness are unfinished; the R9 caller audit found the stable activity.list raw renderer call outside the named function scope.
- Recovery attempts: 0
- Next action: Resume task-10 from this checkpoint, complete both closed process corpora and Human plain-copy pairing, and obtain the task-scope disposition for cmd/fid/activity.go before changing its handler; finish the no-bypass caller audit and rerun the exact gate.
- Human attention required: no

## Summary of what was implemented

- Added one internal/cli Human failure renderer: verbatim bounded JSON message, escaped scalar/key-detail reason, exact Stop/Attention ordering, and the closed next table. Reason fixtures cover the whole-line 200/201-scalar boundary, Unicode, and whole-escape cuts; nullable object values render null, while nested objects and object arrays are omitted outside the overrides.
- Adopted it in the scoped family writers, grammar/unknown-command transport, standalone init/rebuild/maintenance/export/capabilities/unavailable failures, and gate post-frame refusals. Streamed gate output remains on stdout.
- Changed all 18 final R20 fallbacks to use task-9a's BoundedMessage on the actual non-whitespace cause, retaining each fixed empty-cause message and its code/details/exit. Converged retained JSON plain-copy messages without changing their valid Claim or original/copy/import controls.
- Added real-binary protected-Decision, unsuccessful terminal gate, and sensitive duplicate-component refusal checks. Converged retained Human assertions with provenance and retained the original rejection/state assertions. This is partial task-10 work, not task completion.

## Files created or modified

- `cmd/fid/activation.go`
- `cmd/fid/canonical_request.go`
- `cmd/fid/capabilities.go`
- `cmd/fid/claim.go`
- `cmd/fid/concern.go`
- `cmd/fid/decision.go`
- `cmd/fid/deviation.go`
- `cmd/fid/e2e_v01d_link_test.go`
- `cmd/fid/envelope.go`
- `cmd/fid/envelope_test.go`
- `cmd/fid/evidence.go`
- `cmd/fid/export.go`
- `cmd/fid/finding.go`
- `cmd/fid/gate.go`
- `cmd/fid/generic_run_test.go`
- `cmd/fid/grammar_corpus_test.go`
- `cmd/fid/human_failure_test.go`
- `cmd/fid/human_outcome.go`
- `cmd/fid/init.go`
- `cmd/fid/init_test.go`
- `cmd/fid/input_descriptor_test.go`
- `cmd/fid/integration/e2e_v01_closure_test.go`
- `cmd/fid/integration/exit_contract_test.go`
- `cmd/fid/integration/p23_human_failure_test.go`
- `cmd/fid/integration/p24_human_failure_test.go`
- `cmd/fid/issue_ac_test.go`
- `cmd/fid/issue_get_list.go`
- `cmd/fid/issue_label_test.go`
- `cmd/fid/issue_mutation_response.go`
- `cmd/fid/issue_transition_test.go`
- `cmd/fid/label_response.go`
- `cmd/fid/label_test.go`
- `cmd/fid/link_response.go`
- `cmd/fid/link_test.go`
- `cmd/fid/maintenance.go`
- `cmd/fid/p24_gate_evidence_test.go`
- `cmd/fid/p24_run_test.go`
- `cmd/fid/rebuild.go`
- `cmd/fid/rebuild_test.go`
- `cmd/fid/receipt.go`
- `cmd/fid/retained_actor_request_contract_test.go`
- `cmd/fid/review.go`
- `cmd/fid/root_test.go`
- `cmd/fid/run.go`
- `cmd/fid/stop.go`
- `cmd/fid/unavailable.go`
- `cmd/fid/workflow.go`
- `internal/cli/human_failure.go`
- `internal/cli/human_failure_test.go`
- `project-docs/specs/05_DESIGN_PLAN.md`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-24/tasks/task-10-done.md`

## Documentation updated

- project-docs/specs/05_DESIGN_PLAN.md section 6: exact Human failure labels, ordering, cut behavior, and gate-streaming exception.
- project-docs/specs/06_DESIGN_PLAN_CLI.md: bounded domain-cause messages, the four plain-copy sentences, and section 9 examples.
- This report records the incomplete corpus and adoption criteria explicitly. No ledger, triage, decisions, QA artifact, task contract, dependency, or task-11/task-13 view was edited.

## Validation Contract Coverage

- VAL-036: FAIL (unfinished evidence). TestP24HumanProtectedDecisionFailure and TestP24HumanGateTerminalFailure in cmd/fid/integration/p24_human_failure_test.go exercise the actual binary; TestGrammarBuiltBinaryCorpus in cmd/fid/grammar_corpus_test.go preserves short/long/512/513/Unicode/escape-cut and near-miss rows with exact Human expectations. The required complete RegisteredCodes pairing registry with pinned exclusions is not implemented. The R9 activity caller gap also remains.
- VAL-037: PASS for the deterministic renderer states. TestHumanFailureLines/stops in internal/cli/human_failure_test.go pins stop_ids then attention_ids; its scalars/controls/no-input rows pin absence. TestP23StopFailureCodeMatrix and TestP23StopFailureHumanAndJSONParity in cmd/fid/integration/p23_human_failure_test.go compare exact actual-process lines, including both blocking-ID arrays and nonblocking not-found/version/invalid rows.
- VAL-038: PASS for the closed rendering table and absence states. TestHumanFailureLines pins every next-table entry and the non-input invalid-argument absence. Its state-conflict, transition-not-permitted, resolution-invalid and issue-not-found rows, together with the process matrix's stop-not-found, attention-not-found and attention-version-conflict rows, pin more than five outside-table codes with no next line. TestGrammarBuiltBinaryCorpus compares input-leaf grammar guidance and non-input grammar absence, including stop raise and issue create.
- VAL-039: PASS. TestHumanFailureLines pins scalar sorting, scalar arrays, every closed override, control details, empty nested arrays and nullable-object null. TestHumanFailureReasonCut pins exactly 200 uncut, 201 retaining 200 plus ellipsis, non-ASCII scalar counting, and an ESC escape crossing the cut dropped whole.
- VAL-067: FAIL (unfinished evidence). TestP24GenericFallbackMessage in cmd/fid/human_failure_test.go pins 512/513 escaped scalars, Unicode, ESC, escaped boundaries, whole-escape cuts and empty/whitespace causes. TestP24GenericFallbackSensitiveProcess proves a real activation.component generic domain failure: idempotency_results/insert sorted/unique rejection, unchanged fallback details/exit, no sentinel in either stream and no appended Event. Seventeen remaining R20 domain-process rows and the pinned-unreachability completeness check are not implemented.
- VAL-027: FAIL (message pairing unfinished). TestP24EnvelopeAndClaimOriginalCopyAndImport in cmd/fid/envelope_test.go now pins the envelope JSON message prior registration binding does not match envelope and the retained details. Its existing original/copy/import controls remain. The built-binary Human counterpart is still required.
- VAL-028: FAIL (message pairing unfinished). The same retained original/copy/import test pins claim acquire binding is invalid in JSON; a built-binary Human counterpart remains required.
- VAL-029: FAIL (message pairing unfinished). TestP24RunStartOriginalCopyAndImport in cmd/fid/p24_run_test.go preserves the otherwise valid Claim and original/copy/import controls and pins run claim binding is invalid. TestP24GateOriginalCopyAndImport in cmd/fid/p24_gate_evidence_test.go pins run gate: repository binding is invalid, retaining its controls and code/details. Built-binary Human pairing of both copied-worktree refusals remains required.
- VAL-052: PASS. Final exact gate reports cmd/fid 409.529s, internal/workflow 112.664s, slowest other ordinary package internal/export 198.511s, and tagged integration 362.561s. Every ordinary package is below 600s and integration below 900s.
- VAL-054: UNVERIFIED Human judgment, advanced only by the emitted denials. The required Human walkthrough remains the D2 checkpoint after task-11; this Developer run does not attest it.

R9 adoption inventory (regenerated by the exact sixth gate command). Calls below are source traces; they do not substitute for the unfinished process completeness criterion. Every non-diagnostic scoped Human branch reaches writeHumanFailure, which calls cli.WriteHumanFailure; the noted activity call prevents claiming no bypass.

| R9 function / current source | Shared-renderer call or trace |
|---|---|
| `renderActivationError` / `cmd/fid/activation.go:323` | `writeHumanFailure` |
| `renderCanonicalRequestError` / `cmd/fid/canonical_request.go:98` | `writeHumanFailure` |
| `writeGrammarFailure` / `cmd/fid/canonical_request.go:162` | `writeHumanFailure` |
| `writeHumanFailure` / `cmd/fid/canonical_request.go:180` | `cli.WriteHumanFailure` |
| `writeGenericFailure` / `cmd/fid/canonical_request.go:196` | `writeHumanFailure` |
| `renderCapabilitiesError` / `cmd/fid/capabilities.go:75` | `writeHumanFailure` |
| `renderClaimError` / `cmd/fid/claim.go:174` | `renderJudgmentObservationError`, `renderStopError`, `writeClaimFailure`, `writeGenericFailure` |
| `writeClaimFailure` / `cmd/fid/claim.go:237` | `writeHumanFailure` |
| `renderConcernError` / `cmd/fid/concern.go:102` | `renderJudgmentObservationError`, `writeClaimFailure`, `writeGenericFailure` |
| `renderDecisionError` / `cmd/fid/decision.go:120` | `renderJudgmentObservationError`, `writeClaimFailure`, `writeGenericFailure` |
| `renderJudgmentObservationError` / `cmd/fid/decision.go:143` | `writeClaimFailure` |
| `renderDeviationError` / `cmd/fid/deviation.go:93` | `renderJudgmentObservationError`, `writeClaimFailure`, `writeGenericFailure` |
| `renderEnvelopeError` / `cmd/fid/envelope.go:189` | `renderJudgmentObservationError`, `writeEnvelopeFailure`, `writeGenericFailure`, `writeHumanFailure` |
| `writeEnvelopeFailure` / `cmd/fid/envelope.go:245` | `writeHumanFailure` |
| `renderEvidenceError` / `cmd/fid/evidence.go:183` | `renderJudgmentObservationError`, `writeGenericFailure`, `writeRunFailure` |
| `renderExportError` / `cmd/fid/export.go:96` | `writeHumanFailure` |
| `renderFindingError` / `cmd/fid/finding.go:99` | `renderJudgmentObservationError`, `writeClaimFailure`, `writeGenericFailure` |
| `renderGateError` / `cmd/fid/gate.go:286` | `renderJudgmentObservationError`, `renderStopError`, `writeGateFailure`, `writeGenericFailure` |
| `writeGateFailure` / `cmd/fid/gate.go:342` | `writeHumanFailure` |
| `renderHumanOutcomeError` / `cmd/fid/human_outcome.go:121` | `renderJudgmentObservationError`, `writeClaimFailure`, `writeGenericFailure` |
| `renderIssueReadError` / `cmd/fid/issue_get_list.go:98` | `writeHumanFailure` |
| `renderIssueError` / `cmd/fid/issue_helpers.go:25` | `renderJournalReplayError`; raw activity.list call remains a scope gap; diagnostics excluded |
| `renderIssueMutationError` / `cmd/fid/issue_mutation_response.go:99` | `writeHumanFailure` |
| `renderJournalInspectError` / `cmd/fid/journal.go:82` | `renderIssueError`; excluded diagnostic helper/path (activity caller audit still pending) |
| `renderJournalReplayError` / `cmd/fid/journal_errors.go:12` | Excluded diagnostic helper/path; activity caller audit still pending. |
| `renderLabelError` / `cmd/fid/label_response.go:11` | `renderIssueMutationError`, `writeHumanFailure` |
| `renderLinkError` / `cmd/fid/link_response.go:11` | `writeHumanFailure` |
| `renderMaintenanceCheckError` / `cmd/fid/maintenance.go:165` | `writeHumanFailure` |
| `renderRebuildProjectError` / `cmd/fid/rebuild.go:137` | `renderRebuildFailure`, `renderRebuildOperationError` |
| `renderRebuildOperationError` / `cmd/fid/rebuild.go:155` | `writeHumanFailure` |
| `renderReceiptError` / `cmd/fid/receipt.go:78` | `writeGenericFailure`, `writeRunFailure` |
| `renderRequirementError` / `cmd/fid/requirement.go:136` | `renderDecisionError`, `writeClaimFailure` |
| `renderResumeError` / `cmd/fid/resume.go:53` | `renderEnvelopeError`, `writeClaimFailure` |
| `renderReviewListError` / `cmd/fid/review.go:57` | `writeClaimFailure`, `writeGenericFailure` |
| `renderReviewError` / `cmd/fid/review.go:129` | `renderJudgmentObservationError`, `writeClaimFailure`, `writeGenericFailure` |
| `renderRunError` / `cmd/fid/run.go:258` | `renderJudgmentObservationError`, `renderStopError`, `writeGenericFailure`, `writeRunFailure` |
| `writeRunFailure` / `cmd/fid/run.go:305` | `writeHumanFailure` |
| `renderBudgetSampleError` / `cmd/fid/run.go:316` | `renderRunError`, `writeRunFailure` |
| `renderRunRecoveryError` / `cmd/fid/run_recovery.go:99` | `renderRunError`, `writeRunFailure` |
| `renderSQLCommandError` / `cmd/fid/sql.go:51` | `renderIssueError`; excluded diagnostic helper/path (activity caller audit still pending) |
| `renderStopError` / `cmd/fid/stop.go:125` | `writeGenericFailure`, `writeStopFailure` |
| `writeStopFailure` / `cmd/fid/stop.go:180` | `writeHumanFailure` |
| `renderWorkflowError` / `cmd/fid/workflow.go:252` | `renderIssueMutationError`, `writeHumanFailure` |

Additional non-R9-shaped adapters call writeHumanFailure: renderInitFailure (init.go), renderRebuildFailure (rebuild.go), renderCapabilitiesInvalidRequest (capabilities.go), renderMaintenanceDivergence (maintenance.go), writeCapabilityUnavailable (unavailable.go), and gate admission/terminal/infrastructure branches (gate.go). The old raw renderInitError/renderRebuildError helpers were removed after their owning typed writers adopted the renderer. The activity handler still calls renderIssueError; stats.go's explicit JSON-port calls also require caller-level contract audit before declaring adoption complete. Neither caller file was modified.

R20 process coverage disposition. Each listed final mapper uses genericFailure directly or via writeGenericFailure. Unit family checks are deliberately not labelled domain-process corpus coverage, and no unreachability exclusion is invented.

| Final fallback | Process row / remaining work |
|---|---|
| `cmd/fid/activation.go` / `activationFailure` | TestP24GenericFallbackSensitiveProcess: actual duplicate-component domain rejection in both modes |
| `cmd/fid/claim.go` / `renderClaimError` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/concern.go` / `renderConcernError` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/decision.go` / `renderDecisionError` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/deviation.go` / `renderDeviationError` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/envelope.go` / `renderEnvelopeError` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/evidence.go` / `renderEvidenceError` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/finding.go` / `renderFindingError` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/gate.go` / `renderGateError` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/human_outcome.go` / `renderHumanOutcomeError` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/issue_get_list.go` / `issueReadFailure` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/issue_mutation_response.go` / `issueMutationFailure` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/link_response.go` / `linkFailure` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/receipt.go` / `renderReceiptError` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/review.go` / `renderReviewListError` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/review.go` / `renderReviewError` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/run.go` / `renderRunError` | Pending domain-process row; not excluded or claimed unreachable |
| `cmd/fid/stop.go` / `renderStopError` | Pending domain-process row; not excluded or claimed unreachable |

The other five exact R20 matches are not generic fallbacks: canonical_request.go's CF-047 combination failure, init.go's initialization argument, link_response.go's typed endpoint error, and finding.go/review.go's typed Finding target errors. The fixed fallback strings remain for empty causes, as required. issue_update_test.go:331 already pins unknown flag text following task-9a (2d5e9d55c83086b46268ca69793f30f65031b91b); its pre-convergence fallback assertion was displaced by task-9a, so no additional task-10 weakening or rewrite was made there.

Retained convergence provenance. Locations below are the pre-task HEAD's displaced assertion locations (0181fc75d6fde6e6f0aa8b156c7f6125b09c0921), not current line numbers. Introducing commits were checked with git blame/log, including relocated fixtures. Task numbers in filenames supplied no authority. Assertions now pin the canonical Human representation or accepted actual cause; original commands, stdout, exit/metadata, DB/journal/no-Event and original/copy/import controls remain. Transition-edge and rebuild conflict data omitted by the new closed Human form are additionally retained on the propagated underlying cause.

| Displaced assertion location | Introducing commit | Recorded phase/task |
|---|---|---|
| `cmd/fid/e2e_v01d_link_test.go:32` | `1679230ee4dca4d1cbffbb5b2732b7d1a9dbd06d` | `phase-23/task-12` |
| `cmd/fid/envelope_test.go:273` | `1baaecc7983872b06e2a47e6df0801bfb1370ea0` | `phase-24/task-6` |
| `cmd/fid/envelope_test.go:275` | `1baaecc7983872b06e2a47e6df0801bfb1370ea0` | `phase-24/task-6` |
| `cmd/fid/envelope_test.go:385` | `b7d1fd56b801d10378b35dff8f45684d1990804d` | `phase-24/task-6` |
| `cmd/fid/generic_run_test.go:164` | `8d1c83855c0f681515e44e2c3a61bbe543ef515b` | `phase-24/task-7` |
| `cmd/fid/generic_run_test.go:400` | `8c763e9cddee1d4ff0df3379e08b286ef74d5300` | `phase-21/task-11` |
| `cmd/fid/generic_run_test.go:411` | `8c763e9cddee1d4ff0df3379e08b286ef74d5300` | `phase-21/task-11` |
| `cmd/fid/grammar_corpus_test.go:235` | `2d5e9d55c83086b46268ca69793f30f65031b91b` | `phase-24/task-9a` |
| `cmd/fid/init_test.go:427` | `42f5d05244eb67aaa9b389e2c612a5d667a16636` | `phase-20/task-3` |
| `cmd/fid/init_test.go:499` | `42f5d05244eb67aaa9b389e2c612a5d667a16636` | `phase-20/task-3` |
| `cmd/fid/input_descriptor_test.go:94` | `545803dc2185f565d2b02258bd18ee0fd1d7d2ed` | `phase-24/task-3` |
| `cmd/fid/input_descriptor_test.go:2615` | `29c8682c98c612443ed85df503d818dacc70552c` | `phase-24/task-9` |
| `cmd/fid/integration/p23_human_failure_test.go:25` | `c6e070b1ea0b3aff339564c47bb1cd39cfccbcd0` | `phase-23/task-21` |
| `cmd/fid/integration/p23_human_failure_test.go:182` | `c6e070b1ea0b3aff339564c47bb1cd39cfccbcd0` | `phase-23/task-21` |
| `cmd/fid/issue_ac_test.go:282` | `59f02dfeccbb636887cfa62ee3e73661a1c74317` | `phase-7/task-1` |
| `cmd/fid/issue_ac_test.go:501` | `62a3e46ebb69b76b93b5a559dc82d9bc9c0ccd3e` | `phase-20/task-10` |
| `cmd/fid/issue_label_test.go:291` | `16f4f8d2d2b5870d1502e6713300f5671dcfdf4e` | `phase-20/task-10` |
| `cmd/fid/issue_label_test.go:302` | `16f4f8d2d2b5870d1502e6713300f5671dcfdf4e` | `phase-20/task-10` |
| `cmd/fid/issue_label_test.go:493` | `16f4f8d2d2b5870d1502e6713300f5671dcfdf4e` | `phase-20/task-10` |
| `cmd/fid/issue_transition_test.go:616` | `0173eba6d4fcff718a76e89a902f74bbfecce1da` | `phase-5/task-5` |
| `cmd/fid/issue_transition_test.go:617` | `0173eba6d4fcff718a76e89a902f74bbfecce1da` | `phase-5/task-5` |
| `cmd/fid/label_test.go:219` | `62a3e46ebb69b76b93b5a559dc82d9bc9c0ccd3e` | `phase-20/task-10` |
| `cmd/fid/label_test.go:275` | `d5ad0ea2f5b4dbbaf6f468dd20523636d48a0d6a` | `phase-13/task-1` |
| `cmd/fid/link_test.go:131` | `1679230ee4dca4d1cbffbb5b2732b7d1a9dbd06d` | `phase-23/task-12` |
| `cmd/fid/link_test.go:163` | `0abe93edc2d98f1adda4cd5b1bec4e15f64ae06d` | `phase-12/task-2` |
| `cmd/fid/link_test.go:171` | `0abe93edc2d98f1adda4cd5b1bec4e15f64ae06d` | `phase-12/task-2` |
| `cmd/fid/link_test.go:180` | `0abe93edc2d98f1adda4cd5b1bec4e15f64ae06d` | `phase-12/task-2` |
| `cmd/fid/link_test.go:188` | `1679230ee4dca4d1cbffbb5b2732b7d1a9dbd06d` | `phase-23/task-12` |
| `cmd/fid/link_test.go:197` | `1679230ee4dca4d1cbffbb5b2732b7d1a9dbd06d` | `phase-23/task-12` |
| `cmd/fid/link_test.go:205` | `1679230ee4dca4d1cbffbb5b2732b7d1a9dbd06d` | `phase-23/task-12` |
| `cmd/fid/link_test.go:354` | `66392c4c8d4356dd435b1ca37e30e3626aa1d2f5` | `phase-12/task-3` |
| `cmd/fid/p24_gate_evidence_test.go:363` | `32a9ae535dbfb34be0dff7e369c715c2a4640b90` | `phase-24/task-8` |
| `cmd/fid/p24_run_test.go:201` | `8d1c83855c0f681515e44e2c3a61bbe543ef515b` | `phase-24/task-7` |
| `cmd/fid/rebuild_test.go:485` | `f1b18fb837ba4aa20920911ffbf7579ad6540206` | `phase-3/task-5` |
| `cmd/fid/rebuild_test.go:486` | `f1b18fb837ba4aa20920911ffbf7579ad6540206` | `phase-3/task-5` |
| `cmd/fid/root_test.go:33` | `2d5e9d55c83086b46268ca69793f30f65031b91b` | `phase-24/task-9a` |
| `cmd/fid/rebuild_test.go:456` | `dc49d3d5f73a3850b900e76828c3a044240e48e6` | `phase-1/task-4` |
| `cmd/fid/retained_actor_request_contract_test.go:167` | `bc328f9136ec44d31824d1fe2b1c63d1c8b7c921` | `phase-21/task-9` |
| `cmd/fid/integration/exit_contract_test.go:105` | `2d5e9d55c83086b46268ca69793f30f65031b91b` | `phase-24/task-9a` |
| `cmd/fid/integration/e2e_v01_closure_test.go:48` | `1679230ee4dca4d1cbffbb5b2732b7d1a9dbd06d` | `phase-23/task-12` |

## Tests run

All tests/gates ran from /home/mark/projects/fid unless explicitly stated. All listed terminations were observed from the nested command result; empty polling output was never treated as success. Red attempts remain failed evidence. The test-result cache was cleared once as setup before gate attempt 2 to obtain fresh required package times; cleanup is not PASS evidence.

- `go test -tags integration ./cmd/fid/integration -run '^TestP24HumanProtectedDecisionFailure$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Specified first behavior red: Human printed the old execution-bound-transition cause rather than the protected-Decision JSON sentence; package 2.518s. The identical selection subsequently passed in 2.470s in this session.
- `go test ./internal/cli -run '^TestHumanFailure' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). New renderer API was absent; the identical selection subsequently passed in 0.004s.
- `go test ./cmd/fid -run '^TestP24(FamilyFailureRenderer|GenericFallbackMessage)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Initial fixture compile failure used a nonexistent FailureError.Failure field; corrected to the actual Code field before the behavior red.
- `go test ./cmd/fid -run '^TestP24(FamilyFailureRenderer|GenericFallbackMessage)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Genuine renderer/fallback red: old single-line branches and fixed generic messages; package 0.160s.
- `go test ./cmd/fid -run '^TestP24(FamilyFailureRenderer|GenericFallbackMessage)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Enhanced generic-message expectation remained red after renderer adoption; package 0.158s. Minimum shared cause construction then made the identical selection green in 0.163s.
- `go test -tags integration ./cmd/fid/integration -run '^TestP24HumanGateTerminalFailure$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Actual terminal gate refusal had empty Human stderr; package 1.568s. The protected-Decision/terminal-gate combined selection then passed in 2.422s.
- `go test ./cmd/fid -run '^(TestGrammar.*|TestInput.*|TestAdoptedExampleProcessExclusivity|TestInit.*|TestIssue.*|TestLabel.*|TestLink.*|TestRetained.*|TestRunStartBindingMatrixRejectsWithoutMutation|TestP24RunStartOriginalCopyAndImport|TestP24Envelope.*|TestP24Claim.*|TestP24Gate.*|TestRebuild.*)$' -count=1 > /tmp/task10-retained-red.log` - FAIL (terminated `exit_code: 1`). Terminal output was redirected, so no terminal token is invented. The captured diagnostic log's last line is `FAIL`, and it reports package 43.185s. Displaced Human/fallback expectations were converged with provenance; the scratch log is not required durable evidence.
- `go test ./cmd/fid -run '^(TestGrammarMessageBoundaries|TestRetainedInventoryRequiredActorHumanErrorsRemainProse|TestAdoptedExampleProcessExclusivity|TestRunStartBindingMatrixRejectsWithoutMutation)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Fixture corrections: Actor domain pattern, example allowed flag and the earlier mismatched-Envelope-revision cause; package 3.353s.
- `go test ./cmd/fid -run '^(TestGrammar.*|TestRetainedInventoryRequiredActorHumanErrorsRemainProse|TestAdoptedExampleProcessExclusivity|TestRunStartBindingMatrixRejectsWithoutMutation|TestP24RunStartOriginalCopyAndImport|TestP24GateOriginalCopyAndImport|TestInitCommandRequiresExplicitPrefixForInvalidDirectoryBasename|TestInitCommandRejectsReinitWithoutChanges)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	45.069s` (terminated `exit_code: 0`). Retained grammar/binding convergence selection.
- `go test ./cmd/fid -run '^(TestIssue.*|TestLabel.*|TestLink.*|TestRebuildCommandReports.*)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Package 3.341s; the preceding scratch edit had failed with Python IndentationError and made no edits.
- `go test ./cmd/fid -run '^(TestIssue.*|TestLabel.*|TestLink.*|TestRebuildCommandReports.*)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Package 3.251s; a fixture's errors.As expectation incorrectly assumed FailureError exposes Unwrap. Original conflict fields were instead asserted on its propagated cause text without changing exit.go or removing fields.
- `go test ./cmd/fid -run '^(TestIssue.*|TestLabel.*|TestLink.*|TestRebuildCommandReports.*)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.184s` (terminated `exit_code: 0`). Original conflict/edge data and state checks retained.
- `go test -tags integration ./cmd/fid/integration -run '^(TestP23StopFailureHumanAndJSONParity|TestP23StopFailureCodeMatrix|TestProcessExitContract)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Retained single-line/old rework next text; package 5.314s. Extended pairing selection including TestP24Human.*Failure subsequently passed in 5.628s.
- `go test ./internal/cli -run '^TestHumanFailureLines/repository$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Seven-hex repository reason disagreed with the accepted twelve-hex short HEAD; package 0.003s. Identical selection green after the minimum length correction.
- `go test ./internal/cli -run '^TestHumanFailureLines/empty_nested_array$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Empty object array was incorrectly flattened; package 0.003s. Identical selection green after consulting the typed detail definition.
- `go test -tags integration ./cmd/fid/integration -run '^TestP24GenericFallbackSensitiveProcess$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Fixture used unsupported activation phase 24; package 1.450s.
- `go test -tags integration ./cmd/fid/integration -run '^TestP24GenericFallbackSensitiveProcess$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Fixture omitted environment.source; package 0.695s. Corrected the fixture to the accepted phase/literal environment form; identical selection green in 1.686s, reaching the intended domain fallback.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/cli` - FAIL - `FAIL` (terminated `exit_code: 1`). First ordinary pre-gate tier: cmd/fid 308.286s, cli 1.113s; retained governance and removed endpoint Human text.
- `go test ./cmd/fid -run '^(TestPromisedGovernanceCommandsAreAbsent|TestV01DLinkFlowAndRebuildParity)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Focused reproduction, package 0.340s.
- `go test ./cmd/fid -run '^(TestPromisedGovernanceCommandsAreAbsent|TestV01DLinkFlowAndRebuildParity)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.441s` (terminated `exit_code: 0`). Ran successfully twice during retained root-allowed expectation refinement.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/cli` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	1.115s` (terminated `exit_code: 0`). cmd/fid 307.357s; observed termination, no silent-success inference.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestP23StopFailureHumanAndJSONParity|TestP23StopFailureCodeMatrix|TestProcessExitContract|TestP24Human.*Failure|TestP24GenericFallbackSensitiveProcess)$'` - FAIL - `FAIL` (terminated `exit_code: 1`). Typed nested filter suppressed actual event=null; package 6.555s.
- `go test ./internal/cli -run '^TestHumanFailureLines/nullable_object$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Specific null-value regression, package 0.003s.
- `go test ./internal/cli -run '^TestHumanFailureLines/nullable_object$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	0.003s` (terminated `exit_code: 0`). Scalar null retained; non-null nested values still omitted.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestP23StopFailureHumanAndJSONParity|TestP23StopFailureCodeMatrix|TestProcessExitContract|TestP24Human.*Failure|TestP24GenericFallbackSensitiveProcess)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.413s` (terminated `exit_code: 0`). Byte-identical recovery of the preceding red command.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/cli` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	1.122s` (terminated `exit_code: 0`). Race tier after the last CLI implementation change.
- `go test -tags integration ./cmd/fid/integration -run '^TestV01ClosureJourneyAndRebuildParity$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Focused reproduction of required-gate failure: removed requirement.sync single-line stderr expectation; package 0.715s.
- `go test -tags integration ./cmd/fid/integration -run '^TestV01ClosureJourneyAndRebuildParity$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Progress to the later unsupported implements relation's raw Human wording; package 0.950s.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestV01ClosureJourneyAndRebuildParity|TestP23StopFailureHumanAndJSONParity|TestP23StopFailureCodeMatrix|TestProcessExitContract|TestP24Human.*Failure|TestP24GenericFallbackSensitiveProcess)$'` - FAIL - `FAIL` (terminated `exit_code: 1`). Same later relation assertion; package 6.607s. This failure remains evidence and was not routed elsewhere.
- `go test -tags integration ./cmd/fid/integration -run '^TestV01ClosureJourneyAndRebuildParity$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.207s` (terminated `exit_code: 0`). Full closure journey now reaches migration, stale refresh and repeated rebuild checks.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestV01ClosureJourneyAndRebuildParity|TestP23StopFailureHumanAndJSONParity|TestP23StopFailureCodeMatrix|TestProcessExitContract|TestP24Human.*Failure|TestP24GenericFallbackSensitiveProcess)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.958s` (terminated `exit_code: 0`). Final tagged pre-gate tier on the checkpoint tree.

Scratch-only diagnostic failures were also retained as limitations, not gates: an initial provenance script tried to blame an untracked new file and exited 1; it was restricted to retained tracked assertions. A later heading-extraction read assumed DEC headings were H2 and exited 1; the actual H3 authority sections were read. Neither failed diagnostic supplies PASS evidence. No control, task contract or evaluator material was edited or used to repair these diagnostics.


Real entrypoint observations, cwd /tmp: go build -o /tmp/fid-task10 ./cmd/fid had already terminated exit_code: 0 with no output before the direct probe. These expected refusals are runtime observations, not undeclared non-zero gate PASS claims.

- `/tmp/fid-task10 stop raise --bogus` - OBSERVED REFUSAL - `next: see fid stop raise --help; template: fid stop raise --example` (terminated `exit_code: 2`). Its preceding lines were exactly error: unknown flag: --bogus and reason: invalid-argument; argument=--bogus, allowed=--actor|--example|--idempotency-key|--input.
- `/tmp/fid-task10 --json stop raise --bogus` - OBSERVED REFUSAL (terminated `exit_code: 2`): empty stdout; operation stop.raise, code invalid-argument, invalid-input, retryable false, subject null, full spelling in argument/value, and the same four allowed flags. The decoded message was unknown flag: --bogus.
- TestP24HumanProtectedDecisionFailure, TestP24HumanGateTerminalFailure and TestP24GenericFallbackSensitiveProcess exercised freshly built actual fid subprocesses in both modes; the focused tagged and full tagged commands below observed their terminations.

QA-plan section 8 assigns D2 review and Human fixture regeneration after task-11. Those are independent Reviewer/Human work, not additional task-10 Developer commands. No held-out scenario or evaluator was read. Report-lint is the independent Runner control required after the report commit; it was not executed or claimed PASS here.

## Exact quality-gate results

Both full attempts executed all seven byte-identical task commands in order from the repository root. Attempt 1 is immutable failed evidence; attempt 2 is the final checkpoint gate. The rg commands both expect and observed exit_code: 0 with declared matches; no narrower command substituted for any gate.

Attempt 1:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.6%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `rg -n "func (write[A-Za-z]*Failure|render[A-Za-z]*Error)\(" cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/label_response.go:11:func renderLabelError(cmd *cobra.Command, operation string, cause error) error {` (terminated `exit_code: 0`)
- `rg -n 'CodeInvalidArgument, "invalid [^"]*"' cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/finding.go:144:	return writeGenericFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Finding request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})` (terminated `exit_code: 0`)

Attempt 2:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	112.664s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.6%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	362.561s` (terminated `exit_code: 0`)
- `rg -n "func (write[A-Za-z]*Failure|render[A-Za-z]*Error)\(" cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/sql.go:51:func renderSQLCommandError(writer io.Writer, err error) error {` (terminated `exit_code: 0`)
- `rg -n 'CodeInvalidArgument, "invalid [^"]*"' cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/finding.go:144:	return writeGenericFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Finding request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})` (terminated `exit_code: 0`)

Final fresh package times: cmd/fid 409.529s; internal/workflow 112.664s; slowest other ordinary package internal/export 198.511s; cmd/fid/integration 362.561s. Attempt 1's ordinary cmd/fid time was 350.670s; its workflow/issues and other package lines were cached, so no fresh time is claimed for those cached rows. Attempt 1 tagged integration reported 357.104s with exit_code: 1.

## Deviations from spec

- The task is unfinished: complete RegisteredCodes and R20 domain-process completeness registries, Human plain-copy process pairing, and the no-bypass caller audit remain. These are RECOVERY_REQUIRED work; a green automated gate does not establish these missing acceptance criteria.
- Mandatory required-command failed-attempt note: gate attempt 1's exact tagged integration command ended with `FAIL` (terminated `exit_code: 1`, 357.104s). Task-10's accepted second Human line displaced assertRemovedCommandFailure's single-line expectation (introduced by task-9a). Preserving the invocation and typed JSON/journal checks, the in-scope assertion was converged. The focused rerun then reached a later unsupported-link assertion introduced by phase-23/task-12; task-10's accepted fixed JSON sentence displaced its raw-domain Human substring. Its exact canonical error/reason expectation now retains the relation value, while the original no-append checks remain. The focused journey and combined tagged race checks passed; gate attempt 2 reran the byte-identical full tagged command with exit_code: 0. No failure was routed to another task.
- The first ordinary touched-package race attempt ended with `FAIL` (terminated `exit_code: 1`, cmd/fid 308.286s); task-10 displaced the retained governance single-line assertion and removed-endpoint raw wording. Provenance was checked, canonical lines added, and the original cause/empty-stdout/journal checks retained. The byte-identical touched-package command passed (cmd/fid 307.357s; cli 1.115s).
- The first focused tagged race attempt ended with `FAIL` (terminated `exit_code: 1`, 6.555s). The schema-based nested-detail filter incorrectly suppressed an actual scalar null. The nullable-object unit was first red, then the one-line in-scope value filter repair restored event=null without displaying non-null nested objects/arrays. The identical tagged selection passed in 8.413s; cli's race tier then passed in 1.122s. This changed no failure schema or assertion contract.
- Other red attempts retained in Tests run cover the specified protected-Decision first failure, missing renderer API, generic cause construction, empty nested array and 12-hex repository reason, terminal gate stderr, displaced retained output, and fixture/script construction mistakes. No test was deleted, command narrowed, pass condition changed, or flaky green mined.

## Concerns, technical debt, or notes for the Tech Lead

- Same task-10 continuation is required. Complete the two finite process corpora and their honest row/exclusion completeness checks, including exact both-mode plain-copy messages with the valid Claim prerequisite; retain all checkpoint evidence.
- Scope gap: cmd/fid/activity.go:46 calls the R9 renderIssueError(io.Writer,error) without operation/leaf context; its handler is not a named R9 function. The uniform task-10 accepted change exposes this stable caller's raw failure path. Obtain the bounded task-scope disposition before editing that caller. The call was introduced by 977974c0ec00d3d9a63f7d7515e633683c9b9ac5 (phase-16/task-1), and git merge-base --is-ancestor confirmed it predates phase-24's approved task-pack commit 304aa45f12cf3966567231e6986e3070df409c68. This is not assigned to task-11 or task-13. Also audit stats.go's explicit JSON-port calls against their accepted caller contract before claiming no bypass.
- Headroom warning: none; every freshly timed ordinary package is below 480s.
- The Runner-owned dirty run ledger was excluded from all edits/staging. Report-lint remains Runner-owned and must run after this committed report.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
