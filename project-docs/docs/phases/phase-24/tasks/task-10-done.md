# Phase 24 - Task 10: Human failure lines - Completion report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 2
- Next action: None.
- Human attention required: no

## Summary of what was implemented

Same-task recovery under `TRG-24-011 / DEC-286`. Entry HEAD was `7a8a2968fbcca321e2777ecfa57da0fe1649fc3a`; implementation checkpoint `6607143b68ab2d08444cf5d5ecf9ddfd84ff7f49` and report checkpoint `e0495b50769623b0dc819113d7894b76b1dbd10a` were preserved. The latest RUN-RESUME/TASK-RECOVERY values were compared exactly with the directive: task `task-10`, trigger `task-10:public-interface:unavailable-root-flag-failure-mapping`, attempt `2`, category `in-scope-implementation`, mode `task-scope`, authority `TRG-24-011 / DEC-286`. Allowed/Forbidden fields and checkpoint values are exact directive copies, not paraphrased authority. No second task is implemented.

- The stable unavailable-root-flag guard now returns the DEC-286 grammar failure through the existing writeGrammarFailure path. `export import --project` gives the exact invalid-argument envelope/exit 2, exact two Human lines, empty stdout and no filesystem change, before workspace access. Non-project and nonexistent-project cases are both covered. Diagnostic refusal bytes and retained diagnostic assertions remain unchanged.
- canonical_request.go has only the two authorized pure extractions: sorted local long-flag alternatives and the sql/journal inspect predicate. The capability registry, failure helpers, shared renderer, failure definitions and flag domains are untouched in this continuation.
- Completed all seventeen missing registered-code process rows and activation/projection/budget fixtures. The 80-member registry has 74 ordinary pairs, five special pairs and one pinned unsupported-kind exclusion. A supplemental reachable idempotency-conflict pair is retained because that defined code is absent from RegisteredCodes. No exclusion was added.
- Completed the caller audit, including the root guard, diagnostic-only raw branch and negative controls. Preserved DEC-285's six exact process pairs, eighteen R20 process pairs, four plain-copy pairs and retained assertion provenance. Protected-Decision and unsuccessful-gate tests now compare the entire independent expected Human text rather than a prefix.

## Files created or modified

This continuation changes only:

- `cmd/fid/root.go`
- `cmd/fid/canonical_request.go`
- `cmd/fid/root_flags_test.go`
- `cmd/fid/human_failure_caller_audit_test.go`
- `cmd/fid/integration/p24_retained_read_failure_test.go`
- `cmd/fid/integration/p24_code_pairing_test.go`
- `cmd/fid/integration/p24_code_pairing_fixtures_test.go`
- `cmd/fid/integration/p24_judgment_failure_pairs_test.go` (new)
- `cmd/fid/integration/p24_budget_failure_pair_test.go` (new)
- `cmd/fid/integration/p24_human_failure_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-24/tasks/task-10-done.md`

The preserved checkpoint's complete implementation inventory appears in the historical evidence below. Runner-owned run-ledger.md was the sole entry dirt and is neither edited nor staged. No task-11/task-13 surface, protected artifact, dependency or generated file is changed.

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: exact DEC-286 refusal, sorted alternatives, lack of next line, workspace precedence and retained diagnostics.
- This report: current acceptance evidence, complete caller audit, red/green and fresh exact gates, retained historical report/control failures and convergence provenance.
- The checkpoint's `05_DESIGN_PLAN.md` section 6 and earlier `06_DESIGN_PLAN_CLI.md` changes remain part of Task 10; their inventory and evidence are preserved below.

## Validation Contract Coverage

- VAL-036: PASS. `TestP24RegisteredCodePairingCompleteness` (`cmd/fid/integration/p24_code_pairing_test.go:127`) binds all 80 RegisteredCodes members to 79 real-process pairs plus the unchanged pinned unsupported-kind exclusion. `TestP24RegisteredCodePairs` (:155) checks operation/code, no control characters, matching exits and exact independently derived Human text, including labels/order/optional lines. All 74 ordinary rows require empty stdout in both modes. The five special pairs are protected-decision-required, review-rework-exhausted, gate-unsuccessful and the two budget terminal denials. `TestP24HumanProtectedDecisionFailure` (:12 in p24_human_failure_test.go), `TestP24HumanGateTerminalFailure` (:31) and `TestP23StopFailureHumanAndJSONParity` pin complete text and typed exits. Budget rows pin the existing singleton fid.live.v1 terminal denial and exact Human stderr with empty Human stdout (Phase 22 section 4; no new transport exception). `TestGrammarBuiltBinaryCorpus` (`cmd/fid/grammar_corpus_test.go:523`) retains short exact and near-miss isue rows, 600-character flag/word, 512/513 escaped-message boundaries, Unicode and whole-escape cuts with full uncut details. `TestP24RetainedReadFailurePairs` (`p24_retained_read_failure_test.go:15`) retains all six DEC-285 pairs. `TestP24UnavailableImportRootFlagFailurePair` (:104) pins literal DEC-286 JSON bytes and exactly two Human lines/exit 2, empty stdout and whole-filesystem equality for both non-project and nonexistent-project cases. In-process root rows and unchanged diagnostic rows pass in the ordinary gate. No reachable stable leaf is excluded.
- VAL-037: PASS. `TestHumanFailureLines/stops` (`internal/cli/human_failure_test.go:10`) pins stop_ids then attention_ids, with absence in scalars/controls/no-input cases. `TestP23StopFailureHumanAndJSONParity` and the registered `human-outcome-required`, `budget-data-unavailable`, `budget-exhausted`, `stop-unresolved` and protected-Decision process pairs derive the exact blocking line from the actual failure. `issue-not-found`, `finding-not-found`, `review-not-found`, `state-conflict`, `transition-not-permitted` and the DEC-286 refusal pin absence through full-text equality.
- VAL-038: PASS. The real-process registry plus supplemental `idempotency-conflict` covers every next-table code: repository-mismatch, source-digest-mismatch, stop-unresolved, review-rework-exhausted, protected-decision-required, human-outcome-required, obligation-unsatisfied and invalid-argument. `TestGrammarBuiltBinaryCorpus` pins `stop raise --bogus` with next/help/example and `issue create` grammar with no next. The DEC-286 export.import pair also has no next. Full-text `issue-not-found`, `finding-not-found`, `concern-not-found`, `deviation-not-found`, `review-not-found`, `state-conflict` and `operation-unavailable` pairs prove at least five outside-table absences. `TestHumanFailureLines` checks every closed table row and non-input absence.
- VAL-039: PASS. `TestHumanFailureLines` (`internal/cli/human_failure_test.go:10`) pins sorted scalar details, scalar arrays, every override, escaped control values, omitted object arrays and retained scalar null. `TestHumanFailureReasonCut` (:49) pins exactly 200 uncut scalars, 201 retaining 200 plus ellipsis, non-ASCII scalar counting and an escape crossing the cut dropped whole. The fresh full ordinary gate executes both tests; renderer code is unchanged in this continuation.
- VAL-067: PASS. `TestP24GenericFallbackCorpusCompleteness` (`p24_generic_fallback_corpus_test.go:49`) derives the actual final fallback sites with an AST inventory and rejects missing, stale or duplicate rows. All eighteen sites have actual-binary rows in `TestP24GenericFallbackDomainPairs` (:115); exclusion map is empty. Each pins cause message, original typed metadata/details/exit, exact Human text and full filesystem equality. `TestP24GenericFallbackMessage` (`cmd/fid/human_failure_test.go:72`) pins 512/513 escaped scalars, Unicode, ESC, whole-escape cuts, empty and whitespace causes. `TestP24GenericFallbackSensitiveProcess` (`p24_human_failure_test.go:116`) reaches activation.component's duplicate-component domain fallback, retaining failure details/exit/no-Event assertions and proving the environment sentinel absent from both streams. The site-to-process mapping below identifies every producer.
- VAL-027: PASS for Task 10's message half. `TestP24PlainCopyFailurePairs/envelope.create` (`p24_plain_copy_failure_test.go:19`) pins `prior registration binding does not match envelope` in both modes with typed metadata/details/exit and complete filesystem equality. `TestP24EnvelopeAndClaimOriginalCopyAndImport` retains original/copy/import and no-effect controls in the fresh ordinary gate.
- VAL-028: PASS for the message half. `TestP24PlainCopyFailurePairs/claim.acquire` pins `claim acquire binding is invalid` in both modes with unchanged metadata/details/exit/filesystem; `TestP24EnvelopeAndClaimOriginalCopyAndImport` retains the accepted original and imported cases. No Claim prerequisite or retained assertion is removed.
- VAL-029: PASS for the message half. `TestP24PlainCopyFailurePairs/run.start` and `/gate.run` pin `run claim binding is invalid` and `run gate: repository binding is invalid` in both modes with a valid Claim and unchanged filesystem. `TestP24RunStartOriginalCopyAndImport` and `TestP24GateOriginalCopyAndImport` retain the original/copy/import controls in the fresh ordinary gate. These message claims do not replace the earlier tasks' full relocation evidence.
- VAL-052: PASS. Fresh numeric package times from the exact count-cache-cleared gates are recorded below; all ordinary packages are under 600 seconds and integration under 900 seconds. No cached timing is claimed as fresh.
- VAL-054: UNVERIFIED Human judgment. This task advances the deterministic denial fixtures only. QA Window D2 and the Human walkthrough after Task 11 remain mandatory; no Human PASS is claimed.

R9 caller audit. Every stable Human trace below reaches writeHumanFailure -> cli.WriteHumanFailure. `TestP24HumanFailureCallerAudit` (`cmd/fid/human_failure_caller_audit_test.go:15`) rejects raw helper callers outside the two diagnostics, forced-true read-mapper calls, and raw root failures outside the positive diagnostic branch. `TestP24RootFailureAuditNegativeControl` (:91) accepts that diagnostic branch and rejects raw stable-else, unconditional stable and negated-diagnostic branches. The exact sixth gate inventories all 43 named functions:

| R9 function / current source | Shared-renderer call or trace |
|---|---|
| `renderActivationError` / `cmd/fid/activation.go:323` | `writeHumanFailure` |
| `renderCanonicalRequestError` / `cmd/fid/canonical_request.go:109` | `writeHumanFailure` |
| `writeGrammarFailure` / `cmd/fid/canonical_request.go:173` | `writeHumanFailure` |
| `writeHumanFailure` / `cmd/fid/canonical_request.go:191` | `cli.WriteHumanFailure` |
| `writeGenericFailure` / `cmd/fid/canonical_request.go:207` | `writeHumanFailure` |
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
| `renderIssueError` / `cmd/fid/issue_helpers.go:25` | renderJournalReplayError; called only by sql/journal inspect diagnostic paths |
| `renderIssueMutationError` / `cmd/fid/issue_mutation_response.go:99` | `writeHumanFailure` |
| `renderJournalInspectError` / `cmd/fid/journal.go:82` | renderIssueError; journal inspect diagnostic only |
| `renderJournalReplayError` / `cmd/fid/journal_errors.go:12` | Diagnostic helper reached from renderIssueError; stable journal-replay-stopped uses typed family transport |
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
| `renderSQLCommandError` / `cmd/fid/sql.go:51` | renderIssueError; sql diagnostic only |
| `renderStopError` / `cmd/fid/stop.go:125` | `writeGenericFailure`, `writeStopFailure` |
| `writeStopFailure` / `cmd/fid/stop.go:180` | `writeHumanFailure` |
| `renderWorkflowError` / `cmd/fid/workflow.go:252` | `renderIssueMutationError`, `writeHumanFailure` |

| Additional stable caller | Live trace |
|---|---|
| activity.list / `cmd/fid/activity.go:39,50` | live jsonMode -> renderIssueReadError -> writeHumanFailure |
| stats.tokens / `cmd/fid/stats.go:42,52` | live jsonMode -> renderIssueReadError -> writeHumanFailure |
| stats.cost / `cmd/fid/stats.go:81,91` | live jsonMode -> renderIssueReadError -> writeHumanFailure |
| issue.ac.list / `cmd/fid/issue_ac.go:113,122` | live jsonMode -> renderIssueReadError -> writeHumanFailure |
| rejectUnavailableRootFlags / `cmd/fid/root.go:106` | stable leaf -> writeGrammarFailure -> writeHumanFailure; diagnosticLeaf branch preserves only sql/journal inspect raw bytes |

The other standalone adapters also reach writeHumanFailure: renderInitFailure, renderRebuildFailure, renderCapabilitiesInvalidRequest, renderMaintenanceDivergence, writeCapabilityUnavailable, and gate admission/terminal/infrastructure branches. Diagnostics are the task's two established non-fid-cli.v1 paths, not newly excluded stable callers. Stats success output is unchanged. The live predicate is the unchanged sql/journal inspect test extracted from installLeafGrammar.

R20 mapping. Each final mapper below is the identically named subtest of `TestP24GenericFallbackDomainPairs`; the AST inventory has zero exclusions:

| Final mapper site | Built-process domain request |
|---|---|
| activation.go:activationFailure | activation status ACT-missing, no project |
| claim.go:renderClaimError | claim list, no project |
| concern.go:renderConcernError | concern list, no project; pinned lock/mkdir error |
| decision.go:renderDecisionError | decision list, no project |
| deviation.go:renderDeviationError | deviation list, no project |
| envelope.go:renderEnvelopeError | envelope get ENV-missing, no project |
| evidence.go:renderEvidenceError | evidence get EVD-missing, no project |
| finding.go:renderFindingError | finding list, no project |
| gate.go:renderGateError | gate resolve REG-missing --version 1, no project |
| human_outcome.go:renderHumanOutcomeError | human-outcome list, no project |
| issue_get_list.go:issueReadFailure | activity list --since invalid, initialized project |
| issue_mutation_response.go:issueMutationFailure | issue create --kind epic, missing required scheduling |
| link_response.go:linkFailure | link list with a whitespace entity ID |
| receipt.go:renderReceiptError | receipt list, no project |
| review.go:renderReviewListError | review list, no project |
| review.go:renderReviewError | review get REV-missing, no project; retained operation review.record |
| run.go:renderRunError | run list, no project |
| stop.go:renderStopError | stop list, no project |

The exact seventh gate finds 23 sites: these eighteen final generic fallbacks and five typed/non-final invalid-argument definitions (Finding targets in review/finding, init arguments, canonical combination and link endpoint type). The five are not unreachable exclusions. The four plain-copy messages are pinned separately as above.

Retained assertion provenance. The checkpoint's introducing-commit inventory below is preserved as historical provenance, not reconstructed from filenames or task numbers. The locations are displaced pre-task locations at `0181fc75d6fde6e6f0aa8b156c7f6125b09c0921`. The accepted Task 10 Human representation/domain-cause change caused those convergences; their subjects and original commands, exit/metadata, DB/journal/no-Event and original/copy/import controls remain. Fresh full gates execute the retained assertions. This continuation changes only its own corpus/assertions and adds root tests; no external retained assertion is displaced by the two product files.

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

Fresh ownership checks for this continuation: code-pairing corpus, fixture helpers, caller audit and unavailable-import regression were introduced by `6607143b68ab2d08444cf5d5ecf9ddfd84ff7f49` (`Refs: phase-24/task-10`); protected-Decision/terminal-gate tests by `d574cf4587d790a83bbd14de22ff918edbe1e115` (`Refs: phase-24/task-10`). DEC-286 clauses 5 and 10 authorize completion of these own fixtures. Existing root-flag assertions originate at `cd53c43319a3301d3a8cc0089b6785d9afbb163f` (`Refs: phase-20/task-2`); their bodies remain byte-identical, as does the complete canonical_request.go suffix starting at writeGrammarFailure. The export.import registry row originates at `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee` (`Refs: phase-22/task-15`); its capability values are unchanged. No ownership was inferred from labels.

## Tests run

- `git diff --check` - FAIL - `+> ` (terminated `exit_code: 2`)

- `git diff --check` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

Current session commands ran from `/home/mark/projects/fid` unless a working directory is specified. Every recorded command's own terminal result was watched to a real exit; no empty polling response or outer-wrapper completion was used as termination. Focused checks and touched-package race tiers precede and do not replace the full gates. FAIL rows are immutable red evidence; later greens do not erase them.

- `go test ./cmd/fid -run '^TestP24(UnavailableRootFlagMapping|DiagnosticRootFlagBytesPreserved)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

- `go test ./cmd/fid -run '^TestP24(UnavailableRootFlagMapping|DiagnosticRootFlagBytesPreserved)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.173s` (terminated `exit_code: 0`)

- `go test -tags integration ./cmd/fid/integration -run '^TestP24UnavailableImportRootFlagFailurePair$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

- `go test ./cmd/fid -run '^TestP24(HumanFailureCallerAudit|RootFailureAuditNegativeControl)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.168s` (terminated `exit_code: 0`)

initial missing rows/fixtures 13.575s:

- `go test -tags integration ./cmd/fid/integration -run '^TestP24(UnavailableImportRootFlagFailurePair|RegisteredCodePairingCompleteness|RegisteredCodePairs)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

duplicate fixture name build failure:

- `go test -tags integration ./cmd/fid/integration -run '^TestP24(UnavailableImportRootFlagFailurePair|RegisteredCodePairingCompleteness|RegisteredCodePairs)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

terminal decode panic 1.416s:

- `go test -tags integration ./cmd/fid/integration -run '^TestP24(UnavailableImportRootFlagFailurePair|RegisteredCodePairingCompleteness|RegisteredCodePairs)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

domain prerequisite failures 22.570s:

- `go test -tags integration ./cmd/fid/integration -run '^TestP24(UnavailableImportRootFlagFailurePair|RegisteredCodePairingCompleteness|RegisteredCodePairs)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

human obligation prerequisites 27.060s:

- `go test -tags integration ./cmd/fid/integration -run '^TestP24(UnavailableImportRootFlagFailurePair|RegisteredCodePairingCompleteness|RegisteredCodePairs)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

- `go test -tags integration ./cmd/fid/integration -run '^TestP24RegisteredCodePairs/(activation-incomplete|human-outcome-required|review-changes-required|evidence-type-not-accepted)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

- `go test -tags integration ./cmd/fid/integration -run '^TestP24RegisteredCodePairs/human-outcome-required$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

- `go test -tags integration ./cmd/fid/integration -run '^TestP24RegisteredCodePairs/human-outcome-required$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.059s` (terminated `exit_code: 0`)

- `go test -tags integration ./cmd/fid/integration -run '^TestP24(UnavailableImportRootFlagFailurePair|RegisteredCodePairingCompleteness|RegisteredCodePairs)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	20.582s` (terminated `exit_code: 0`)

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	314.734s` (terminated `exit_code: 0`)

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestP24(UnavailableImportRootFlagFailurePair|RegisteredCodePairingCompleteness|RegisteredCodePairs|RetainedReadFailurePairs|GenericFallbackDomainPairs|GenericFallbackCorpusCompleteness|GenericFallbackSensitiveProcess|PlainCopyFailurePairs|HumanProtectedDecisionFailure|HumanGateTerminalFailure)|TestP23StopFailureHumanAndJSONParity)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	29.515s` (terminated `exit_code: 0`)

Real entrypoint evidence, outside a project:

- `go build -o /tmp/fid-task10-dec286 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid-task10-dec286 --project /tmp/fid-task10-dec286-runtime-kkAymw/nonexistent export import missing.zip --repository /tmp/fid-task10-dec286-runtime-kkAymw --idempotency-key probe` - OBSERVED REFUSAL - `reason: invalid-argument; argument=--project, allowed=--idempotency-key|--repository` (terminated `exit_code: 2`)
- `/tmp/fid-task10-dec286 --json --project /tmp/fid-task10-dec286-runtime-kkAymw/nonexistent export import missing.zip --repository /tmp/fid-task10-dec286-runtime-kkAymw --idempotency-key probe` - OBSERVED REFUSAL - `{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["--idempotency-key","--repository"],"argument":"--project","value":"--project"},"message":"--project is not available for export import","retryable":false,"subject":null},"ok":false,"operation":"export.import"}` (terminated `exit_code: 2`)

Both probes ran in `/tmp/fid-task10-dec286-runtime-kkAymw`; the directory remained empty. The process regression separately captures stdout/stderr and pins unchanged filesystem bytes for both working-directory/path cases. The ordinary no-output build success is not used as a go-test token.

QA-plan section 5 assigns report-lint to the independent Phase Runner; section 8 assigns D2 fixture regeneration/review and the Human walkthrough after Task 11. No held-out scenario or evaluator was read. A Developer report check, if run under the role's report-finalization sequence, is not the Runner's independent receipt and never advances the ledger.

Historical evidence preservation. The complete incoming 86,395-byte report (SHA-256 `6e3854634443638428063eda4cffa37e2c6ba42b8d5e485cabd617583986cca7`) is retained verbatim below by prefixing each nonempty physical line with `> ` and each blank line with `>`. Removing the quote prefixes reconstructs the incoming bytes exactly. Historical PASS receipts are UNVERIFIED as fresh-session evidence; they are retained receipts from their original sessions, not renewed attestations. All historical BLOCKED/RECOVERY_REQUIRED outcomes, failed commands/tokens/terminations, six-pair DEC-285 evidence, original probe text, full-gate attempts, report-lint failures and provenance remain unchanged.

Additional preserved Runner control evidence at `e0495b50769623b0dc819113d7894b76b1dbd10a`: TRG-24-011 and its ledger record report-lint token `OUTCOME BLOCKED` (terminated `exit_code: 23`). This was not executed or re-attested in the current Developer session. The incoming report also retains prior Runner tokens `OUTCOME RECOVERY_REQUIRED` / exit_code 10 and `OUTCOME BLOCKED` / exit_code 11. None is rewritten as green or substituted for the Runner's pending independent final control.

<details>
<summary>Complete preserved report at e0495b5 (historical evidence)</summary>

> # Phase 24 - Task 10: Human failure lines - Recovery attempt 2 checkpoint
>
> ## Workflow outcome
>
> - Task outcome: BLOCKED
> - Category: public-interface-decision
> - Condition: The stable export.import root-flag rejection in rejectUnavailableRootFlags has no fid-cli.v1 JSON failure and bypasses the section-9 renderer in both modes; root.go's guard is outside the R9 functions and DEC-285 extensions, and Task 10 does not authorize choosing its missing code/details/exit mapping.
> - Recovery attempts: 2
> - Next action: Obtain protected authority for the root-flag failure mapping and bounded caller scope, then resume the same Task 10 checkpoint to complete the remaining registry rows and fixture repairs and rerun the exact gates. No other task is assigned this failure.
> - Human attention required: yes
>
> ## Summary of what was implemented
>
> Recovery attempt 2, dispatched under `TRG-24-010 / DEC-285`, continues only phase-24/task-10. The latest RUN-RESUME and TASK-RECOVERY were checked against the dispatch: task `task-10`, trigger `task-10:public-interface:activity-list-failure-envelope-mapping`, attempt `2`, category `in-scope-implementation`, repair mode `task-scope`, repair authority `TRG-24-010 / DEC-285`. Entry HEAD was `ea2f7e79119bc6e187bbe04ff68e603d84205f4a`; implementation `d574cf4587d790a83bbd14de22ff918edbe1e115` and report `9c3e003e36d6547c56593c384a2b9e32d4db6095` were retained. Ledger/directive values were compared as exact values, not paraphrased authority.
>
> - Implemented DEC-285: activity.list reads JSON mode before workspace resolution and routes unwrapped errors through the existing issue-read mapper. Stats tokens/cost and issue.ac.list use live JSON mode for their six failure calls. Stats success rendering is untouched.
> - Six exact DEC-285 built-binary pairs pass: both modes pin message, reason, envelope metadata/details, exit, empty stdout and unchanged complete filesystem snapshots.
> - All eighteen R20 domain sites now have actual-binary both-mode rows, pinned causes and metadata/details/exit checks; the exclusion map is honestly empty. Completeness checks reject missing/stale/duplicate rows.
> - Four plain-copy built-binary pairs pass with the exact DEC-281/282 messages, accepted Claim prerequisites, unchanged filesystem snapshots and original-worktree controls. Copy fixtures normalize Git's stat index before observation; no production code changed for this fixture correction.
> - Added registry-driven process pairs with one audited unsupported-kind exclusion and a supplemental idempotency-conflict pair. This corpus is partial: seventeen members remain missing and several fixture/transport rows remain red. No reachable code was called unreachable to obtain green.
> - Preserved red caller-audit and process regressions for the protected root-flag guard. The original activity mapping blocker is resolved by DEC-285; the current protected candidate is a different stable failure path, independently reproduced through the real CLI.
>
> Prior attempt-1 Workflow outcome is preserved as historical evidence, not the current authority question:
>
> ```text
> - Task outcome: BLOCKED
> - Category: public-interface-decision
> - Condition: activity.list has no fid-cli.v1 failure producer on its raw domain-error path; the authorized repair cannot preserve existing JSON code/details/exit while deriving Human lines from an absent JSON failure without choosing a protected producer mapping.
> - Recovery attempts: 1
> - Next action: Return the observed activity.list transport evidence to the Phase Runner for Tech Lead authority identifying the accepted producer mapping or explicitly resolving its code/details/exit; obtain bounded dispositions for the audited stats/AC-list caller gaps, then resume the same task's remaining corpora, pairing, renderer adoption and exact gates.
> - Human attention required: yes
> ```
>
> Prior finalized report commit fields are preserved:
>
> ```text
> - `cc99b7325041d933c770bbbbb11af25565c1ced6`
>
> Commit subject (historical)
>
> - `docs(cli): record task 10 activity transport authority blocker`
> ```
>
> The remaining original report text below is historical evidence from earlier sessions; its PASS receipts are not re-attestations in this session. Current attempt-2 evidence appears first under each heading.
>
> Continuation attempt 1, 2026-10-07, remains phase-24/task-10 under TRG-24-009, trigger `task-10:in-scope-implementation:human-failure-corpus-and-caller-audit`, repair category `in-scope-implementation`, repair authority `task spec`. The ledger's five binding values (task identity from the event heading, trigger, attempt, category and authority) and Forbidden changes were checked against the governing TRG without paraphrasing its values. The extension is exactly `cmd/fid/activity.go`.
>
> - Reproduced the activity.list raw transport in a freshly built fid binary, in both modes, with unchanged filesystem snapshots. The JSON envelope prerequisite check is red; no product mapping was invented to obtain green.
> - Completed the requested stats.go assessment and identified the same forced-JSON failure path in issue.ac.list during the caller audit. These are concrete additional caller scope gaps, not assignments to another task.
> - Reran all seven exact Task 10 commands in order and watched each terminate. No product, test, renderer, contract or control was edited in this continuation; the only tracked edit is this report.
> - The RegisteredCodes corpus, complete R20 domain-process corpus and four-leaf built-binary plain-copy pairing remain unfinished. These are ordinary incomplete work, not the protected blocker. The blocker is the independently reproduced absence of the activity JSON producer the directive requires preserving.
> - Prior implementation checkpoint: `d574cf4587d790a83bbd14de22ff918edbe1e115`, `feat(cli): checkpoint task 10 Human failure rendering`; prior finalized report/entry HEAD: `f1ad50fa18eb45b308cc52ffb3ea0c65b6dd198d`. Its typed outcome was RECOVERY_REQUIRED.
>
> Historical implementation summary follows unchanged:
>
> - Added one internal/cli Human failure renderer: verbatim bounded JSON message, escaped scalar/key-detail reason, exact Stop/Attention ordering, and the closed next table. Reason fixtures cover the whole-line 200/201-scalar boundary, Unicode, and whole-escape cuts; nullable object values render null, while nested objects and object arrays are omitted outside the overrides.
> - Adopted it in the scoped family writers, grammar/unknown-command transport, standalone init/rebuild/maintenance/export/capabilities/unavailable failures, and gate post-frame refusals. Streamed gate output remains on stdout.
> - Changed all 18 final R20 fallbacks to use task-9a's BoundedMessage on the actual non-whitespace cause, retaining each fixed empty-cause message and its code/details/exit. Converged retained JSON plain-copy messages without changing their valid Claim or original/copy/import controls.
> - Added real-binary protected-Decision, unsuccessful terminal gate, and sensitive duplicate-component refusal checks. Converged retained Human assertions with provenance and retained the original rejection/state assertions. This is partial task-10 work, not task completion.
>
> ## Files created or modified
>
> Recovery attempt 2 modifies `cmd/fid/activity.go`, `cmd/fid/stats.go`, `cmd/fid/issue_ac.go` and this report. New tests: `cmd/fid/human_failure_caller_audit_test.go`, `cmd/fid/integration/p24_retained_read_failure_test.go`, `cmd/fid/integration/p24_generic_fallback_corpus_test.go`, `cmd/fid/integration/p24_plain_copy_failure_test.go`, `cmd/fid/integration/p24_code_pairing_test.go`, `cmd/fid/integration/p24_code_pairing_fixtures_test.go`. Entry dirt was only the Runner-owned run-ledger.md; it is not edited or staged. All product changes are within DEC-285's exact three extensions; no internal/activity, mapper, flag-domain, root guard, dependency, task-11/task-13 or protected artifact changes.
>
> Historical attempt-1 inventory and earlier checkpoint inventory follow unchanged:
>
> Current continuation: only `project-docs/docs/phases/phase-24/tasks/task-10-done.md`. The existing dirty run-ledger.md and triage-log.md were present on entry and were not edited or staged. No second task or agent was started.
>
> Historical checkpoint file inventory follows unchanged:
>
> - `cmd/fid/activation.go`
> - `cmd/fid/canonical_request.go`
> - `cmd/fid/capabilities.go`
> - `cmd/fid/claim.go`
> - `cmd/fid/concern.go`
> - `cmd/fid/decision.go`
> - `cmd/fid/deviation.go`
> - `cmd/fid/e2e_v01d_link_test.go`
> - `cmd/fid/envelope.go`
> - `cmd/fid/envelope_test.go`
> - `cmd/fid/evidence.go`
> - `cmd/fid/export.go`
> - `cmd/fid/finding.go`
> - `cmd/fid/gate.go`
> - `cmd/fid/generic_run_test.go`
> - `cmd/fid/grammar_corpus_test.go`
> - `cmd/fid/human_failure_test.go`
> - `cmd/fid/human_outcome.go`
> - `cmd/fid/init.go`
> - `cmd/fid/init_test.go`
> - `cmd/fid/input_descriptor_test.go`
> - `cmd/fid/integration/e2e_v01_closure_test.go`
> - `cmd/fid/integration/exit_contract_test.go`
> - `cmd/fid/integration/p23_human_failure_test.go`
> - `cmd/fid/integration/p24_human_failure_test.go`
> - `cmd/fid/issue_ac_test.go`
> - `cmd/fid/issue_get_list.go`
> - `cmd/fid/issue_label_test.go`
> - `cmd/fid/issue_mutation_response.go`
> - `cmd/fid/issue_transition_test.go`
> - `cmd/fid/label_response.go`
> - `cmd/fid/label_test.go`
> - `cmd/fid/link_response.go`
> - `cmd/fid/link_test.go`
> - `cmd/fid/maintenance.go`
> - `cmd/fid/p24_gate_evidence_test.go`
> - `cmd/fid/p24_run_test.go`
> - `cmd/fid/rebuild.go`
> - `cmd/fid/rebuild_test.go`
> - `cmd/fid/receipt.go`
> - `cmd/fid/retained_actor_request_contract_test.go`
> - `cmd/fid/review.go`
> - `cmd/fid/root_test.go`
> - `cmd/fid/run.go`
> - `cmd/fid/stop.go`
> - `cmd/fid/unavailable.go`
> - `cmd/fid/workflow.go`
> - `internal/cli/human_failure.go`
> - `internal/cli/human_failure_test.go`
> - `project-docs/specs/05_DESIGN_PLAN.md`
> - `project-docs/specs/06_DESIGN_PLAN_CLI.md`
> - `project-docs/docs/phases/phase-24/tasks/task-10-done.md`
>
> ## Documentation updated
>
> Attempt 2 updates only this completion report. Prior task-10 specification documentation remains committed and retained. Decisions, task specs, QA plans/contracts, triage, ledger and workflow controls were not edited.
>
> Historical documentation record follows:
>
> Current continuation: this report adds the protected authority question, actual-binary caller audit, own-session red observation and fresh exact-gate evidence. No specification document changed in this continuation.
>
> Historical documentation record follows unchanged:
>
> - project-docs/specs/05_DESIGN_PLAN.md section 6: exact Human failure labels, ordering, cut behavior, and gate-streaming exception.
> - project-docs/specs/06_DESIGN_PLAN_CLI.md: bounded domain-cause messages, the four plain-copy sentences, and section 9 examples.
> - This report records the incomplete corpus and adoption criteria explicitly. No ledger, triage, decisions, QA artifact, task contract, dependency, or task-11/task-13 view was edited.
>
> ## Validation Contract Coverage
>
> Attempt-2 status:
>
> - VAL-036: BLOCKED for full adoption. `TestP24RetainedReadFailurePairs` proves all six DEC-285 pairs; `TestP24UnavailableImportRootFlagFailurePair` proves the root guard still lacks JSON transport/section-9 lines. `TestP24HumanFailureCallerAudit` detects the exact bypass. `TestP24RegisteredCodePairingCompleteness` remains red for seventeen missing members; process fixtures remain partial. These incomplete rows are separate from the protected mapping condition.
> - VAL-067: eighteen domain-process sites and completeness checks pass in focused selection (`TestP24GenericFallbackDomainPairs`, `TestP24GenericFallbackCorpusCompleteness`). `TestP24GenericFallbackMessage` (cmd/fid/human_failure_test.go:72) covers retained 512/513, empty/whitespace, Unicode and escape/cut units in the ordinary suite, whose only failing test is the new caller audit. `TestP24GenericFallbackSensitiveProcess` passes in the current tagged count=1 selection, reaching the duplicate-component fallback without leaking the environment sentinel. These VAL-067 component successes are distinct from the incomplete registry pairing.
> - VAL-027, VAL-028, VAL-029 message half: `TestP24PlainCopyFailurePairs/envelope.create`, `/claim.acquire`, `/run.start`, `/gate.run` pass both modes with exact four messages, matching metadata/details/exits, empty stdout and complete filesystem equality. The Run fixture retains a valid Claim; original accepted Envelope/Claim/Run invocations remain controls. Original/copy/import retained controls are not deleted or weakened and are exercised by the ordinary suite. Full relocation assertions remain the owning earlier tasks' evidence.
> - VAL-037, VAL-038, VAL-039: current gate results below cover the retained shared-renderer/grammar tests. Partial registry pairing cannot be recorded as full completion of these assertions.
> - VAL-052: observed ordinary times are cmd/fid 340.419s and slowest other fresh package internal/export 112.885s. Internal/workflow reports `(cached)` in the exact gate; a fresh time is UNVERIFIED, not copied from prior reports. Tagged integration reports 352.186s in the final gate below. No cached or prior receipt is presented as a fresh time; timing acceptance does not close the failing gates.
> - VAL-054: UNVERIFIED Human judgment; D2 at task-11 close remains required. This blocked task does not advance the phase.
>
> Attempt-2 caller extension audit:
>
> | Caller | Mapper calls using live JSON mode |
> |---|---|
> | activity.list | `activity.go:39,50` → renderIssueReadError |
> | stats.tokens | `stats.go:42,52` → renderIssueReadError |
> | stats.cost | `stats.go:81,91` → renderIssueReadError |
> | issue.ac.list | `issue_ac.go:113,122` → renderIssueReadError |
>
> `rejectUnavailableRootFlags` (`root.go:106-128`) remains outside those extensions and the named R9 functions. Its raw `fmt.Fprintf` rejection fails the supplemental audit. Historical R9 table entries for named functions remain valid as function adoption evidence; their earlier whole-call-graph no-bypass conclusion is superseded by this actual counterexample.
>
> Current R20 corpus mapping (each site is the identically named subtest of `TestP24GenericFallbackDomainPairs`; no exclusions):
>
> | Final mapper site | Built-process domain request |
> |---|---|
> | activation.go:activationFailure | activation status ACT-missing, no project |
> | claim.go:renderClaimError | claim list, no project |
> | concern.go:renderConcernError | concern list, no project; pinned lock/mkdir error |
> | decision.go:renderDecisionError | decision list, no project |
> | deviation.go:renderDeviationError | deviation list, no project |
> | envelope.go:renderEnvelopeError | envelope get ENV-missing, no project |
> | evidence.go:renderEvidenceError | evidence get EVD-missing, no project |
> | finding.go:renderFindingError | finding list, no project |
> | gate.go:renderGateError | gate resolve REG-missing --version 1, no project |
> | human_outcome.go:renderHumanOutcomeError | human-outcome list, no project |
> | issue_get_list.go:issueReadFailure | activity list --since invalid, initialized project |
> | issue_mutation_response.go:issueMutationFailure | issue create --kind epic, missing required scheduling |
> | link_response.go:linkFailure | link list with a whitespace entity ID |
> | receipt.go:renderReceiptError | receipt list, no project |
> | review.go:renderReviewListError | review list, no project |
> | review.go:renderReviewError | review get REV-missing, no project; retained operation review.record |
> | run.go:renderRunError | run list, no project |
> | stop.go:renderStopError | stop list, no project |
>
> The raw R20 rg inventory also matches five typed, non-final invalid-argument definitions (review's Finding target, init arguments, canonical request, link endpoint type and finding's Finding target). They are not generic fallbacks and are not unreachable exclusions. The AST completeness test identifies only actual final genericFailure/writeGenericFailure calls. Each row pins the real domain message and original typed details/exit and checks both full filesystem snapshots.
>
>
> Historical validation evidence follows unchanged:
>
> Current continuation evidence and limits:
>
> - VAL-036: BLOCKED for activity.list adoption. `cmd/fid/activity.go:35` returns workspace resolution errors without transport; `:46` renders List errors through the raw helper, and `:49` returns the untyped error. Actual-binary `activity list --since invalid` and `--on invalid` in both modes produce prose only with exit 7; missing-project calls are silent with exit 7. There is no JSON failure to pair or supply to the shared renderer. The existing RegisteredCodes completeness obligation is also unfinished.
> - VAL-037/038/039: the retained named tests in the historical coverage rows below were included in this continuation's fresh ordinary and tagged gates, both exit 0. No new corpus completeness claim is made.
> - VAL-067 and VAL-027/028/029 message pairing: FAIL (unfinished acceptance work, not the protected category). The complete R20 registry and built-binary Human/JSON plain-copy pairs are still absent. Existing unit boundary/sensitive-process and original/copy/import assertions remain intact and pass the current full suites; those retained tests do not establish the missing process corpus.
> - VAL-052: PASS from this continuation's exact gates: cmd/fid 412.776s; internal/workflow 120.190s; slowest other ordinary package internal/export 197.915s; tagged cmd/fid/integration 361.027s. All ordinary packages are below 600s and 480s warning; integration is below 900s.
> - VAL-054: UNVERIFIED Human judgment. D2 after task 11, fixture regeneration and the Human walkthrough remain required; this report does not advance to task 11.
>
> Current caller-level audit supplements the unchanged R9 function/adoption table below:
>
> | Caller / source | Current trace and own-session observation | Disposition |
> |---|---|---|
> | activity.list / activity.go:35,46,49 | Resolution returns an untyped error with no output. List failure calls renderIssueError -> raw error text and returns the untyped error. JSON mode therefore has no code/details; both exits are 7. | Protected mapping authority required under TRG-24-009's explicit boundary. |
> | stats.tokens / stats.go:38,48 | Forces renderIssueReadError(cmd,true,...); Human and JSON missing-project requests both emit the same fid-cli.v1 project-not-found envelope, exit 3, with no effects. | Audit only; no stats.go edit authority. Human-mode renderer adoption remains unsatisfied. |
> | stats.cost / stats.go:73,83 | Same forced-JSON read mapper and observed project-not-found envelope/exit 3 in both modes. | Same audit-only scope gap. |
> | issue.ac.list / issue_ac.go:109,118 | Same forced-JSON read mapper. Both missing-project modes emit project-not-found, exit 3, with no effects. | Additional exact caller scope gap; issue_ac.go is not named by R9 and has no extension here. No edit or alternative-task assignment. |
> | journal inspect, sql / journal.go:70,90; sql.go:39,59 | Remaining raw-helper callers are the two explicitly excluded diagnostics. | Honest diagnostic exclusion already stated by the task; no stable-caller exclusion invented. |
> | root-flag rejection / root.go:121 | Raw rejection is guarded by the capability registry. Stable leaves are registered for both global transports; journal inspect/sql keep their excluded diagnostic flag rejection. | No new stable renderer exception claimed. |
>
> The exact sixth/seventh commands regenerated 43 R9 functions and 23 R20 matches. R20 still has the 18 final generic mappers and five non-generic sites already accounted for below. This is inventory/audit evidence, not a completed no-bypass or process-completeness gate.
>
> Historical coverage, pending-row dispositions and retained-convergence provenance follow unchanged. Historical PASS receipts are preserved evidence from the original session; the current gate receipts below are independently captured in this continuation:
>
> - VAL-036: FAIL (unfinished evidence). TestP24HumanProtectedDecisionFailure and TestP24HumanGateTerminalFailure in cmd/fid/integration/p24_human_failure_test.go exercise the actual binary; TestGrammarBuiltBinaryCorpus in cmd/fid/grammar_corpus_test.go preserves short/long/512/513/Unicode/escape-cut and near-miss rows with exact Human expectations. The required complete RegisteredCodes pairing registry with pinned exclusions is not implemented. The R9 activity caller gap also remains.
> - VAL-037: PASS for the deterministic renderer states. TestHumanFailureLines/stops in internal/cli/human_failure_test.go pins stop_ids then attention_ids; its scalars/controls/no-input rows pin absence. TestP23StopFailureCodeMatrix and TestP23StopFailureHumanAndJSONParity in cmd/fid/integration/p23_human_failure_test.go compare exact actual-process lines, including both blocking-ID arrays and nonblocking not-found/version/invalid rows.
> - VAL-038: PASS for the closed rendering table and absence states. TestHumanFailureLines pins every next-table entry and the non-input invalid-argument absence. Its state-conflict, transition-not-permitted, resolution-invalid and issue-not-found rows, together with the process matrix's stop-not-found, attention-not-found and attention-version-conflict rows, pin more than five outside-table codes with no next line. TestGrammarBuiltBinaryCorpus compares input-leaf grammar guidance and non-input grammar absence, including stop raise and issue create.
> - VAL-039: PASS. TestHumanFailureLines pins scalar sorting, scalar arrays, every closed override, control details, empty nested arrays and nullable-object null. TestHumanFailureReasonCut pins exactly 200 uncut, 201 retaining 200 plus ellipsis, non-ASCII scalar counting, and an ESC escape crossing the cut dropped whole.
> - VAL-067: FAIL (unfinished evidence). TestP24GenericFallbackMessage in cmd/fid/human_failure_test.go pins 512/513 escaped scalars, Unicode, ESC, escaped boundaries, whole-escape cuts and empty/whitespace causes. TestP24GenericFallbackSensitiveProcess proves a real activation.component generic domain failure: idempotency_results/insert sorted/unique rejection, unchanged fallback details/exit, no sentinel in either stream and no appended Event. Seventeen remaining R20 domain-process rows and the pinned-unreachability completeness check are not implemented.
> - VAL-027: FAIL (message pairing unfinished). TestP24EnvelopeAndClaimOriginalCopyAndImport in cmd/fid/envelope_test.go now pins the envelope JSON message prior registration binding does not match envelope and the retained details. Its existing original/copy/import controls remain. The built-binary Human counterpart is still required.
> - VAL-028: FAIL (message pairing unfinished). The same retained original/copy/import test pins claim acquire binding is invalid in JSON; a built-binary Human counterpart remains required.
> - VAL-029: FAIL (message pairing unfinished). TestP24RunStartOriginalCopyAndImport in cmd/fid/p24_run_test.go preserves the otherwise valid Claim and original/copy/import controls and pins run claim binding is invalid. TestP24GateOriginalCopyAndImport in cmd/fid/p24_gate_evidence_test.go pins run gate: repository binding is invalid, retaining its controls and code/details. Built-binary Human pairing of both copied-worktree refusals remains required.
> - VAL-052: PASS. Final exact gate reports cmd/fid 409.529s, internal/workflow 112.664s, slowest other ordinary package internal/export 198.511s, and tagged integration 362.561s. Every ordinary package is below 600s and integration below 900s.
> - VAL-054: UNVERIFIED Human judgment, advanced only by the emitted denials. The required Human walkthrough remains the D2 checkpoint after task-11; this Developer run does not attest it.
>
> R9 adoption inventory (regenerated by the exact sixth gate command). Calls below are source traces; they do not substitute for the unfinished process completeness criterion. Every non-diagnostic scoped Human branch reaches writeHumanFailure, which calls cli.WriteHumanFailure; the noted activity call prevents claiming no bypass.
>
> | R9 function / current source | Shared-renderer call or trace |
> |---|---|
> | `renderActivationError` / `cmd/fid/activation.go:323` | `writeHumanFailure` |
> | `renderCanonicalRequestError` / `cmd/fid/canonical_request.go:98` | `writeHumanFailure` |
> | `writeGrammarFailure` / `cmd/fid/canonical_request.go:162` | `writeHumanFailure` |
> | `writeHumanFailure` / `cmd/fid/canonical_request.go:180` | `cli.WriteHumanFailure` |
> | `writeGenericFailure` / `cmd/fid/canonical_request.go:196` | `writeHumanFailure` |
> | `renderCapabilitiesError` / `cmd/fid/capabilities.go:75` | `writeHumanFailure` |
> | `renderClaimError` / `cmd/fid/claim.go:174` | `renderJudgmentObservationError`, `renderStopError`, `writeClaimFailure`, `writeGenericFailure` |
> | `writeClaimFailure` / `cmd/fid/claim.go:237` | `writeHumanFailure` |
> | `renderConcernError` / `cmd/fid/concern.go:102` | `renderJudgmentObservationError`, `writeClaimFailure`, `writeGenericFailure` |
> | `renderDecisionError` / `cmd/fid/decision.go:120` | `renderJudgmentObservationError`, `writeClaimFailure`, `writeGenericFailure` |
> | `renderJudgmentObservationError` / `cmd/fid/decision.go:143` | `writeClaimFailure` |
> | `renderDeviationError` / `cmd/fid/deviation.go:93` | `renderJudgmentObservationError`, `writeClaimFailure`, `writeGenericFailure` |
> | `renderEnvelopeError` / `cmd/fid/envelope.go:189` | `renderJudgmentObservationError`, `writeEnvelopeFailure`, `writeGenericFailure`, `writeHumanFailure` |
> | `writeEnvelopeFailure` / `cmd/fid/envelope.go:245` | `writeHumanFailure` |
> | `renderEvidenceError` / `cmd/fid/evidence.go:183` | `renderJudgmentObservationError`, `writeGenericFailure`, `writeRunFailure` |
> | `renderExportError` / `cmd/fid/export.go:96` | `writeHumanFailure` |
> | `renderFindingError` / `cmd/fid/finding.go:99` | `renderJudgmentObservationError`, `writeClaimFailure`, `writeGenericFailure` |
> | `renderGateError` / `cmd/fid/gate.go:286` | `renderJudgmentObservationError`, `renderStopError`, `writeGateFailure`, `writeGenericFailure` |
> | `writeGateFailure` / `cmd/fid/gate.go:342` | `writeHumanFailure` |
> | `renderHumanOutcomeError` / `cmd/fid/human_outcome.go:121` | `renderJudgmentObservationError`, `writeClaimFailure`, `writeGenericFailure` |
> | `renderIssueReadError` / `cmd/fid/issue_get_list.go:98` | `writeHumanFailure` |
> | `renderIssueError` / `cmd/fid/issue_helpers.go:25` | `renderJournalReplayError`; raw activity.list call remains a scope gap; diagnostics excluded |
> | `renderIssueMutationError` / `cmd/fid/issue_mutation_response.go:99` | `writeHumanFailure` |
> | `renderJournalInspectError` / `cmd/fid/journal.go:82` | `renderIssueError`; excluded diagnostic helper/path (activity caller audit still pending) |
> | `renderJournalReplayError` / `cmd/fid/journal_errors.go:12` | Excluded diagnostic helper/path; activity caller audit still pending. |
> | `renderLabelError` / `cmd/fid/label_response.go:11` | `renderIssueMutationError`, `writeHumanFailure` |
> | `renderLinkError` / `cmd/fid/link_response.go:11` | `writeHumanFailure` |
> | `renderMaintenanceCheckError` / `cmd/fid/maintenance.go:165` | `writeHumanFailure` |
> | `renderRebuildProjectError` / `cmd/fid/rebuild.go:137` | `renderRebuildFailure`, `renderRebuildOperationError` |
> | `renderRebuildOperationError` / `cmd/fid/rebuild.go:155` | `writeHumanFailure` |
> | `renderReceiptError` / `cmd/fid/receipt.go:78` | `writeGenericFailure`, `writeRunFailure` |
> | `renderRequirementError` / `cmd/fid/requirement.go:136` | `renderDecisionError`, `writeClaimFailure` |
> | `renderResumeError` / `cmd/fid/resume.go:53` | `renderEnvelopeError`, `writeClaimFailure` |
> | `renderReviewListError` / `cmd/fid/review.go:57` | `writeClaimFailure`, `writeGenericFailure` |
> | `renderReviewError` / `cmd/fid/review.go:129` | `renderJudgmentObservationError`, `writeClaimFailure`, `writeGenericFailure` |
> | `renderRunError` / `cmd/fid/run.go:258` | `renderJudgmentObservationError`, `renderStopError`, `writeGenericFailure`, `writeRunFailure` |
> | `writeRunFailure` / `cmd/fid/run.go:305` | `writeHumanFailure` |
> | `renderBudgetSampleError` / `cmd/fid/run.go:316` | `renderRunError`, `writeRunFailure` |
> | `renderRunRecoveryError` / `cmd/fid/run_recovery.go:99` | `renderRunError`, `writeRunFailure` |
> | `renderSQLCommandError` / `cmd/fid/sql.go:51` | `renderIssueError`; excluded diagnostic helper/path (activity caller audit still pending) |
> | `renderStopError` / `cmd/fid/stop.go:125` | `writeGenericFailure`, `writeStopFailure` |
> | `writeStopFailure` / `cmd/fid/stop.go:180` | `writeHumanFailure` |
> | `renderWorkflowError` / `cmd/fid/workflow.go:252` | `renderIssueMutationError`, `writeHumanFailure` |
>
> Additional non-R9-shaped adapters call writeHumanFailure: renderInitFailure (init.go), renderRebuildFailure (rebuild.go), renderCapabilitiesInvalidRequest (capabilities.go), renderMaintenanceDivergence (maintenance.go), writeCapabilityUnavailable (unavailable.go), and gate admission/terminal/infrastructure branches (gate.go). The old raw renderInitError/renderRebuildError helpers were removed after their owning typed writers adopted the renderer. The activity handler still calls renderIssueError; stats.go's explicit JSON-port calls also require caller-level contract audit before declaring adoption complete. Neither caller file was modified.
>
> R20 process coverage disposition. Each listed final mapper uses genericFailure directly or via writeGenericFailure. Unit family checks are deliberately not labelled domain-process corpus coverage, and no unreachability exclusion is invented.
>
> | Final fallback | Process row / remaining work |
> |---|---|
> | `cmd/fid/activation.go` / `activationFailure` | TestP24GenericFallbackSensitiveProcess: actual duplicate-component domain rejection in both modes |
> | `cmd/fid/claim.go` / `renderClaimError` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/concern.go` / `renderConcernError` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/decision.go` / `renderDecisionError` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/deviation.go` / `renderDeviationError` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/envelope.go` / `renderEnvelopeError` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/evidence.go` / `renderEvidenceError` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/finding.go` / `renderFindingError` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/gate.go` / `renderGateError` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/human_outcome.go` / `renderHumanOutcomeError` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/issue_get_list.go` / `issueReadFailure` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/issue_mutation_response.go` / `issueMutationFailure` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/link_response.go` / `linkFailure` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/receipt.go` / `renderReceiptError` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/review.go` / `renderReviewListError` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/review.go` / `renderReviewError` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/run.go` / `renderRunError` | Pending domain-process row; not excluded or claimed unreachable |
> | `cmd/fid/stop.go` / `renderStopError` | Pending domain-process row; not excluded or claimed unreachable |
>
> The other five exact R20 matches are not generic fallbacks: canonical_request.go's CF-047 combination failure, init.go's initialization argument, link_response.go's typed endpoint error, and finding.go/review.go's typed Finding target errors. The fixed fallback strings remain for empty causes, as required. issue_update_test.go:331 already pins unknown flag text following task-9a (2d5e9d55c83086b46268ca69793f30f65031b91b); its pre-convergence fallback assertion was displaced by task-9a, so no additional task-10 weakening or rewrite was made there.
>
> Retained convergence provenance. Locations below are the pre-task HEAD's displaced assertion locations (0181fc75d6fde6e6f0aa8b156c7f6125b09c0921), not current line numbers. Introducing commits were checked with git blame/log, including relocated fixtures. Task numbers in filenames supplied no authority. Assertions now pin the canonical Human representation or accepted actual cause; original commands, stdout, exit/metadata, DB/journal/no-Event and original/copy/import controls remain. Transition-edge and rebuild conflict data omitted by the new closed Human form are additionally retained on the propagated underlying cause.
>
> | Displaced assertion location | Introducing commit | Recorded phase/task |
> |---|---|---|
> | `cmd/fid/e2e_v01d_link_test.go:32` | `1679230ee4dca4d1cbffbb5b2732b7d1a9dbd06d` | `phase-23/task-12` |
> | `cmd/fid/envelope_test.go:273` | `1baaecc7983872b06e2a47e6df0801bfb1370ea0` | `phase-24/task-6` |
> | `cmd/fid/envelope_test.go:275` | `1baaecc7983872b06e2a47e6df0801bfb1370ea0` | `phase-24/task-6` |
> | `cmd/fid/envelope_test.go:385` | `b7d1fd56b801d10378b35dff8f45684d1990804d` | `phase-24/task-6` |
> | `cmd/fid/generic_run_test.go:164` | `8d1c83855c0f681515e44e2c3a61bbe543ef515b` | `phase-24/task-7` |
> | `cmd/fid/generic_run_test.go:400` | `8c763e9cddee1d4ff0df3379e08b286ef74d5300` | `phase-21/task-11` |
> | `cmd/fid/generic_run_test.go:411` | `8c763e9cddee1d4ff0df3379e08b286ef74d5300` | `phase-21/task-11` |
> | `cmd/fid/grammar_corpus_test.go:235` | `2d5e9d55c83086b46268ca69793f30f65031b91b` | `phase-24/task-9a` |
> | `cmd/fid/init_test.go:427` | `42f5d05244eb67aaa9b389e2c612a5d667a16636` | `phase-20/task-3` |
> | `cmd/fid/init_test.go:499` | `42f5d05244eb67aaa9b389e2c612a5d667a16636` | `phase-20/task-3` |
> | `cmd/fid/input_descriptor_test.go:94` | `545803dc2185f565d2b02258bd18ee0fd1d7d2ed` | `phase-24/task-3` |
> | `cmd/fid/input_descriptor_test.go:2615` | `29c8682c98c612443ed85df503d818dacc70552c` | `phase-24/task-9` |
> | `cmd/fid/integration/p23_human_failure_test.go:25` | `c6e070b1ea0b3aff339564c47bb1cd39cfccbcd0` | `phase-23/task-21` |
> | `cmd/fid/integration/p23_human_failure_test.go:182` | `c6e070b1ea0b3aff339564c47bb1cd39cfccbcd0` | `phase-23/task-21` |
> | `cmd/fid/issue_ac_test.go:282` | `59f02dfeccbb636887cfa62ee3e73661a1c74317` | `phase-7/task-1` |
> | `cmd/fid/issue_ac_test.go:501` | `62a3e46ebb69b76b93b5a559dc82d9bc9c0ccd3e` | `phase-20/task-10` |
> | `cmd/fid/issue_label_test.go:291` | `16f4f8d2d2b5870d1502e6713300f5671dcfdf4e` | `phase-20/task-10` |
> | `cmd/fid/issue_label_test.go:302` | `16f4f8d2d2b5870d1502e6713300f5671dcfdf4e` | `phase-20/task-10` |
> | `cmd/fid/issue_label_test.go:493` | `16f4f8d2d2b5870d1502e6713300f5671dcfdf4e` | `phase-20/task-10` |
> | `cmd/fid/issue_transition_test.go:616` | `0173eba6d4fcff718a76e89a902f74bbfecce1da` | `phase-5/task-5` |
> | `cmd/fid/issue_transition_test.go:617` | `0173eba6d4fcff718a76e89a902f74bbfecce1da` | `phase-5/task-5` |
> | `cmd/fid/label_test.go:219` | `62a3e46ebb69b76b93b5a559dc82d9bc9c0ccd3e` | `phase-20/task-10` |
> | `cmd/fid/label_test.go:275` | `d5ad0ea2f5b4dbbaf6f468dd20523636d48a0d6a` | `phase-13/task-1` |
> | `cmd/fid/link_test.go:131` | `1679230ee4dca4d1cbffbb5b2732b7d1a9dbd06d` | `phase-23/task-12` |
> | `cmd/fid/link_test.go:163` | `0abe93edc2d98f1adda4cd5b1bec4e15f64ae06d` | `phase-12/task-2` |
> | `cmd/fid/link_test.go:171` | `0abe93edc2d98f1adda4cd5b1bec4e15f64ae06d` | `phase-12/task-2` |
> | `cmd/fid/link_test.go:180` | `0abe93edc2d98f1adda4cd5b1bec4e15f64ae06d` | `phase-12/task-2` |
> | `cmd/fid/link_test.go:188` | `1679230ee4dca4d1cbffbb5b2732b7d1a9dbd06d` | `phase-23/task-12` |
> | `cmd/fid/link_test.go:197` | `1679230ee4dca4d1cbffbb5b2732b7d1a9dbd06d` | `phase-23/task-12` |
> | `cmd/fid/link_test.go:205` | `1679230ee4dca4d1cbffbb5b2732b7d1a9dbd06d` | `phase-23/task-12` |
> | `cmd/fid/link_test.go:354` | `66392c4c8d4356dd435b1ca37e30e3626aa1d2f5` | `phase-12/task-3` |
> | `cmd/fid/p24_gate_evidence_test.go:363` | `32a9ae535dbfb34be0dff7e369c715c2a4640b90` | `phase-24/task-8` |
> | `cmd/fid/p24_run_test.go:201` | `8d1c83855c0f681515e44e2c3a61bbe543ef515b` | `phase-24/task-7` |
> | `cmd/fid/rebuild_test.go:485` | `f1b18fb837ba4aa20920911ffbf7579ad6540206` | `phase-3/task-5` |
> | `cmd/fid/rebuild_test.go:486` | `f1b18fb837ba4aa20920911ffbf7579ad6540206` | `phase-3/task-5` |
> | `cmd/fid/root_test.go:33` | `2d5e9d55c83086b46268ca69793f30f65031b91b` | `phase-24/task-9a` |
> | `cmd/fid/rebuild_test.go:456` | `dc49d3d5f73a3850b900e76828c3a044240e48e6` | `phase-1/task-4` |
> | `cmd/fid/retained_actor_request_contract_test.go:167` | `bc328f9136ec44d31824d1fe2b1c63d1c8b7c921` | `phase-21/task-9` |
> | `cmd/fid/integration/exit_contract_test.go:105` | `2d5e9d55c83086b46268ca69793f30f65031b91b` | `phase-24/task-9a` |
> | `cmd/fid/integration/e2e_v01_closure_test.go:48` | `1679230ee4dca4d1cbffbb5b2732b7d1a9dbd06d` | `phase-23/task-12` |
>
> ## Tests run
>
> Attempt-2 focused commands, cwd `/home/mark/projects/fid` (each own nested terminal result was watched through termination):
>
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24RetainedReadFailurePairs$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); 3.813s. Written before product edits; all six behavior rows exposed raw activity/forced JSON failures.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24RetainedReadFailurePairs$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.432s` (terminated `exit_code: 0`). Byte-identical green after DEC-285's bounded repair.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24GenericFallbackCorpusCompleteness$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); 0.520s, missing issueMutationFailure/linkFailure rows.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24GenericFallback(CorpusCompleteness|DomainPairs)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.781s` (terminated `exit_code: 0`). All eighteen real domain sites; empty exclusion map.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24PlainCopyFailurePairs$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); 1.684s, Git index stat refresh changed copied filesystem snapshots.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24PlainCopyFailurePairs/run.start$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); 1.427s, diagnostic identified only `.git/index` bytes changed.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24PlainCopyFailurePairs$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.786s` (terminated `exit_code: 0`). Copy normalization occurs before observation; no snapshot assertion removed.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24RegisteredCodePairingCompleteness$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); 0.512s, missing members and initial stale idempotency-conflict row. That reachable defined code is now a supplemental pair because RegisteredCodes does not enumerate it.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24RegisteredCodePairs$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); first process run 7.120s, workflow fixture used a nonexistent get leaf. Accepted workflow.assign fixture preserves the not-found producer assertion.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24RegisteredCodePairs$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); later expanded run 18.860s, structural/domain fixture mistakes, nested-array oracle and JSON gate budget terminal transport. Ordinary fixture fixes made causal progress but this corpus has no green receipt.
> - `go test ./cmd/fid -run '^TestP24HumanFailureCallerAudit$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.169s` (terminated `exit_code: 0`). Initial narrow raw-helper/forced-true check; does not prove the later full audit.
> - `go test ./cmd/fid -run '^TestP24HumanFailureCallerAudit$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); 0.212s. Strengthened root guard audit detects the separate protected path; not repaired outside scope.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24UnavailableImportRootFlagFailurePair$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); 0.686s. Human has no reason line; JSON parsing reaches raw error text. No typed mapping is invented.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24(GenericFallbackCorpusCompleteness|GenericFallbackDomainPairs|RetainedReadFailurePairs|PlainCopyFailurePairs|RegisteredCodePairingCompleteness|RegisteredCodePairs|UnavailableImportRootFlagFailurePair)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); 20.773s. Missing registry members, protected guard and partial activation/projection/budget rows remain red; R20, six adopted paths and four plain-copy pairs pass within this run.
>
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24GenericFallbackCorpusCompleteness$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.511s` (terminated `exit_code: 0`). Byte-identical rerun of the missing-sites red command.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24PlainCopyFailurePairs/run.start$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.486s` (terminated `exit_code: 0`). Byte-identical rerun of the index-refresh diagnosis selection.
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`); cmd/fid 304.188s. All retained ordinary tests passed; the only failing test is the new root-guard audit.
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^TestP24(GenericFallbackCorpusCompleteness|GenericFallbackDomainPairs|GenericFallbackSensitiveProcess|RetainedReadFailurePairs|PlainCopyFailurePairs|RegisteredCodePairingCompleteness|RegisteredCodePairs|UnavailableImportRootFlagFailurePair)$'` - FAIL - `FAIL` (terminated `exit_code: 1`); 18.526s. The protected root guard and partial registry rows/fixtures remain red; the sensitive domain row, eighteen R20 pairs, six DEC-285 pairs and four plain-copy pairs pass under race. This pre-gate selection does not substitute for the required full tagged command.
>
> Real entrypoint evidence:
>
> - `go build -o /tmp/fid-task10-attempt2 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`; observed command wait wall time 0.439552s). Built before corpus-only fixture changes; product tree matches the runtime observations below.
> - `/tmp/fid-task10-attempt2 --project /tmp/fid-task10-runtime-hJuzYn activity list --since invalid`, cwd `/tmp/fid-task10-runtime-hJuzYn` - OBSERVED REFUSAL - `reason: invalid-argument; argument=request, allowed=accepted issue read input` (terminated `exit_code: 2`). Direct exec_command result; preceding error line is exactly `error: activity list: --since must be YYYY-MM-DD or RFC3339`.
> - `/tmp/fid-task10-attempt2 --json --project /tmp/fid-task10-runtime-hJuzYn activity list --since invalid`, same cwd - OBSERVED REFUSAL - `{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["accepted issue read input"],"argument":"request","value":""},"message":"activity list: --since must be YYYY-MM-DD or RFC3339","retryable":false,"subject":null},"ok":false,"operation":"activity.list"}` (terminated `exit_code: 2`). This is direct runtime evidence, not an undeclared non-zero gate PASS.
> - `/tmp/fid-task10-attempt2 --project /tmp/fid-task10-rootguard-BisAwW export import missing.zip --repository /tmp/fid-task10-rootguard-BisAwW --idempotency-key probe`, cwd `/tmp/fid-task10-rootguard-BisAwW` - OBSERVED REFUSAL - `error: --project is not available for export import until its fid-cli.v1 transport is migrated` (terminated `exit_code: 7`). Direct exec_command, not a wrapper status.
> - `/tmp/fid-task10-attempt2 --json --project /tmp/fid-task10-rootguard-BisAwW export import missing.zip --repository /tmp/fid-task10-rootguard-BisAwW --idempotency-key probe`, same cwd - OBSERVED REFUSAL - `error: --project is not available for export import until its fid-cli.v1 transport is migrated` (terminated `exit_code: 7`). Direct exec_command; JSON transport is absent.
> - Independent built-process root-guard probes: `/tmp/fid-task10-attempt2 [--json] --project <fresh-empty-root> export import missing.zip --repository <fresh-empty-root> --idempotency-key probe`, with working directory that fresh root. Both subprocesses returned `exit_code: 7`, empty stdout, unchanged filesystem snapshots, stderr `error: --project is not available for export import until its fid-cli.v1 transport is migrated
> `. The observation program's own wrapper terminated `exit_code: 0`; that wrapper result is not substituted for the child refusals. The permanent TestP24UnavailableImportRootFlagFailurePair reproduces it deterministically.
>
>
> Additional preserved Runner control receipt: run-ledger.md:1300 records the independent report-lint of prior finalized report `9c3e003e36d6547c56593c384a2b9e32d4db6095` as `OUTCOME BLOCKED` (observed `exit_code: 11`). This is historical Runner evidence, UNVERIFIED as a newly executed Developer command; it is never a current PASS. Its absent activity mapping root cause was resolved by DEC-285, with later six-pair green and direct typed activity refusals recorded above. No diagnostic list is invented, and the earlier `OUTCOME RECOVERY_REQUIRED` / `exit_code: 10` evidence also remains intact.
>
> Historical tests and failed report-lint evidence follow unchanged:
>
> Current continuation, cwd /home/mark/projects/fid unless a fixture cwd is stated:
>
> - `go build -o /tmp/fid-task10-recovery ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`). Its yielded exec session was polled to its own exit before any binary observation.
> - `python3 -` (the envelope-prerequisite program below on stdin) - FAIL - `FAIL activity.list JSON failure envelope absent; stdout=''; stderr='error: activity list: --since must be YYYY-MM-DD or RFC3339\n'; observed exit_code=7` (terminated `exit_code: 1`). This red check asserts accepted transport presence and operation only; it selects no failure code/details/exit. Green was not attempted by inventing a mapping.
> - `python3 -` (14 subprocess/snapshot observations summarized below) - PASS - `OBSERVED 14 caller transport probes; all filesystem snapshots unchanged` (terminated `exit_code: 0`). The script records each subprocess returncode and separate stdout/stderr; expected refusals retain their actual non-zero codes rather than becoming gate PASS assertions.
> - `/tmp/fid-task10-recovery --json activity list --since invalid` and `/tmp/fid-task10-recovery activity list --since invalid`, cwd /tmp/fid-task10-audit-D5s0O5 - OBSERVED REFUSALS - `error: activity list: --since must be YYYY-MM-DD or RFC3339` (each terminated `exit_code: 7`). These were also executed directly through exec_command before the snapshot matrix.
>
> The prerequisite program was executed byte-for-byte as:
>
> ```python
> import json,subprocess,sys
> p=subprocess.run(['/tmp/fid-task10-recovery','--json','activity','list','--since','invalid'],cwd='/tmp/fid-task10-audit-D5s0O5',text=True,capture_output=True)
> try:
>     envelope=json.loads(p.stderr)
> except json.JSONDecodeError:
>     print('FAIL activity.list JSON failure envelope absent; stdout='+repr(p.stdout)+'; stderr='+repr(p.stderr)+'; observed exit_code='+str(p.returncode))
>     sys.exit(1)
> assert envelope['contract_version']=='fid-cli.v1' and envelope['ok'] is False and envelope['operation']=='activity.list'
> print('PASS activity.list JSON failure envelope exists')
> ```
>
> Actual-binary observations: build as above; initialize the project fixture with `/tmp/fid-task10-recovery init --prefix AUD --json` (observed exit 0), then run each argv below both without and with the root `--json` flag. All stdout is empty on refusal. The fixture snapshot consists of every regular file's relative path, mode and SHA-256 content hash, including SQLite and journal. All fourteen before/after snapshots were equal.
>
> | Fixture cwd | Leaf argv | Human stderr / exit | JSON stderr / exit | Success/no-effects control |
> |---|---|---|---|---|
> | /tmp/fid-task10-audit-D5s0O5 | activity list --since invalid | `error: activity list: --since must be YYYY-MM-DD or RFC3339` plus LF / 7 | Identical prose plus LF / 7 | Both no effects. |
> | Same initialized fixture | activity list --on invalid | `error: activity list: --on must be <type>/<id>` plus LF / 7 | Identical prose plus LF / 7 | Both no effects. |
> | Same initialized fixture | activity list | Empty / 0 | Empty / 0 | Human stdout empty for empty timeline; JSON stdout exactly fid-cli.v1 success, activity.list, event:null, result:{items:[]}; both no effects. |
> | /tmp/fid-task10-no-project-fjh2k3i8 | stats tokens | project-not-found JSON envelope / 3 | Same envelope / 3 | Both no effects. |
> | Same non-project fixture | stats cost | project-not-found JSON envelope / 3 | Same envelope / 3 | Both no effects. |
> | Same non-project fixture | issue ac list AUD-missing | project-not-found JSON envelope / 3 | Same envelope / 3 | Both no effects. |
> | Same non-project fixture | activity list | Empty / 7 | Empty / 7 | Both no effects. |
>
> Each stats/AC failure is exactly fid-cli.v1, ok:false, operation matching its leaf, code:project-not-found, category:not-found, retryable:false, message:"fid project was not found", subject:{id:"",type:"project"}, details:{path:""}.
>
> Caller provenance was established with git blame and git show:
>
> - activity.go raw List caller: `977974c0ec00d3d9a63f7d7515e633683c9b9ac5`, phase-16/task-1.
> - stats.go forced read transport: `e51afe3c3e0ea157ee960fb643996453250a37c6`, phase-21/task-9.
> - issue_ac.go forced read transport: `16f4f8d2d2b5870d1502e6713300f5671dcfdf4e`, phase-20/task-10.
>
> For each full hash, its individual `git merge-base --is-ancestor <hash> 304aa45f12cf3966567231e6986e3070df409c68` command was executed and terminated `exit_code: 0` with `exit 0 (no output)`. Earlier provenance does not excuse the task-10 no-bypass requirement.
>
> Historical Runner report-lint evidence is preserved and explicitly attributed: the 20:15 TASK-RECOVERY records `OUTCOME RECOVERY_REQUIRED`, observed `exit_code: 10`, and ten failures against the prior checkpoint. The 20:17 bound recovery repeats that non-complete result. These are Runner observations, not a PASS captured in this Developer continuation; individual diagnostics are not supplied in those records and are not reconstructed. Root cause recorded by the checkpoint/Runner is incomplete acceptance evidence and caller adoption. This continuation narrows the caller failure to an absent typed activity producer plus forced-JSON callers; the corpus and pairing findings remain unresolved. Finalized-report lint remains the Runner's independent control and was not executed or claimed PASS here.
>
> No Go source/test was edited; this report-only protected hand-back does not substitute a green implementation for the red transport observation. The test-result cache was cleared as setup before the seven-command continuation gate to collect fresh times; setup is not gate evidence.
>
> Historical tests, runtime observations and all failed/passing receipts follow unchanged:
>
> All tests/gates ran from /home/mark/projects/fid unless explicitly stated. All listed terminations were observed from the nested command result; empty polling output was never treated as success. Red attempts remain failed evidence. The test-result cache was cleared once as setup before gate attempt 2 to obtain fresh required package times; cleanup is not PASS evidence.
>
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24HumanProtectedDecisionFailure$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Specified first behavior red: Human printed the old execution-bound-transition cause rather than the protected-Decision JSON sentence; package 2.518s. The identical selection subsequently passed in 2.470s in this session.
> - `go test ./internal/cli -run '^TestHumanFailure' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). New renderer API was absent; the identical selection subsequently passed in 0.004s.
> - `go test ./cmd/fid -run '^TestP24(FamilyFailureRenderer|GenericFallbackMessage)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Initial fixture compile failure used a nonexistent FailureError.Failure field; corrected to the actual Code field before the behavior red.
> - `go test ./cmd/fid -run '^TestP24(FamilyFailureRenderer|GenericFallbackMessage)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Genuine renderer/fallback red: old single-line branches and fixed generic messages; package 0.160s.
> - `go test ./cmd/fid -run '^TestP24(FamilyFailureRenderer|GenericFallbackMessage)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Enhanced generic-message expectation remained red after renderer adoption; package 0.158s. Minimum shared cause construction then made the identical selection green in 0.163s.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24HumanGateTerminalFailure$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Actual terminal gate refusal had empty Human stderr; package 1.568s. The protected-Decision/terminal-gate combined selection then passed in 2.422s.
> - `go test ./cmd/fid -run '^(TestGrammar.*|TestInput.*|TestAdoptedExampleProcessExclusivity|TestInit.*|TestIssue.*|TestLabel.*|TestLink.*|TestRetained.*|TestRunStartBindingMatrixRejectsWithoutMutation|TestP24RunStartOriginalCopyAndImport|TestP24Envelope.*|TestP24Claim.*|TestP24Gate.*|TestRebuild.*)$' -count=1 > /tmp/task10-retained-red.log` - FAIL (terminated `exit_code: 1`). Terminal output was redirected, so no terminal token is invented. The captured diagnostic log's last line is `FAIL`, and it reports package 43.185s. Displaced Human/fallback expectations were converged with provenance; the scratch log is not required durable evidence.
> - `go test ./cmd/fid -run '^(TestGrammarMessageBoundaries|TestRetainedInventoryRequiredActorHumanErrorsRemainProse|TestAdoptedExampleProcessExclusivity|TestRunStartBindingMatrixRejectsWithoutMutation)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Fixture corrections: Actor domain pattern, example allowed flag and the earlier mismatched-Envelope-revision cause; package 3.353s.
> - `go test ./cmd/fid -run '^(TestGrammar.*|TestRetainedInventoryRequiredActorHumanErrorsRemainProse|TestAdoptedExampleProcessExclusivity|TestRunStartBindingMatrixRejectsWithoutMutation|TestP24RunStartOriginalCopyAndImport|TestP24GateOriginalCopyAndImport|TestInitCommandRequiresExplicitPrefixForInvalidDirectoryBasename|TestInitCommandRejectsReinitWithoutChanges)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	45.069s` (terminated `exit_code: 0`). Retained grammar/binding convergence selection.
> - `go test ./cmd/fid -run '^(TestIssue.*|TestLabel.*|TestLink.*|TestRebuildCommandReports.*)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Package 3.341s; the preceding scratch edit had failed with Python IndentationError and made no edits.
> - `go test ./cmd/fid -run '^(TestIssue.*|TestLabel.*|TestLink.*|TestRebuildCommandReports.*)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Package 3.251s; a fixture's errors.As expectation incorrectly assumed FailureError exposes Unwrap. Original conflict fields were instead asserted on its propagated cause text without changing exit.go or removing fields.
> - `go test ./cmd/fid -run '^(TestIssue.*|TestLabel.*|TestLink.*|TestRebuildCommandReports.*)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.184s` (terminated `exit_code: 0`). Original conflict/edge data and state checks retained.
> - `go test -tags integration ./cmd/fid/integration -run '^(TestP23StopFailureHumanAndJSONParity|TestP23StopFailureCodeMatrix|TestProcessExitContract)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Retained single-line/old rework next text; package 5.314s. Extended pairing selection including TestP24Human.*Failure subsequently passed in 5.628s.
> - `go test ./internal/cli -run '^TestHumanFailureLines/repository$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Seven-hex repository reason disagreed with the accepted twelve-hex short HEAD; package 0.003s. Identical selection green after the minimum length correction.
> - `go test ./internal/cli -run '^TestHumanFailureLines/empty_nested_array$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Empty object array was incorrectly flattened; package 0.003s. Identical selection green after consulting the typed detail definition.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24GenericFallbackSensitiveProcess$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Fixture used unsupported activation phase 24; package 1.450s.
> - `go test -tags integration ./cmd/fid/integration -run '^TestP24GenericFallbackSensitiveProcess$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Fixture omitted environment.source; package 0.695s. Corrected the fixture to the accepted phase/literal environment form; identical selection green in 1.686s, reaching the intended domain fallback.
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/cli` - FAIL - `FAIL` (terminated `exit_code: 1`). First ordinary pre-gate tier: cmd/fid 308.286s, cli 1.113s; retained governance and removed endpoint Human text.
> - `go test ./cmd/fid -run '^(TestPromisedGovernanceCommandsAreAbsent|TestV01DLinkFlowAndRebuildParity)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Focused reproduction, package 0.340s.
> - `go test ./cmd/fid -run '^(TestPromisedGovernanceCommandsAreAbsent|TestV01DLinkFlowAndRebuildParity)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.441s` (terminated `exit_code: 0`). Ran successfully twice during retained root-allowed expectation refinement.
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/cli` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	1.115s` (terminated `exit_code: 0`). cmd/fid 307.357s; observed termination, no silent-success inference.
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestP23StopFailureHumanAndJSONParity|TestP23StopFailureCodeMatrix|TestProcessExitContract|TestP24Human.*Failure|TestP24GenericFallbackSensitiveProcess)$'` - FAIL - `FAIL` (terminated `exit_code: 1`). Typed nested filter suppressed actual event=null; package 6.555s.
> - `go test ./internal/cli -run '^TestHumanFailureLines/nullable_object$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Specific null-value regression, package 0.003s.
> - `go test ./internal/cli -run '^TestHumanFailureLines/nullable_object$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	0.003s` (terminated `exit_code: 0`). Scalar null retained; non-null nested values still omitted.
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestP23StopFailureHumanAndJSONParity|TestP23StopFailureCodeMatrix|TestProcessExitContract|TestP24Human.*Failure|TestP24GenericFallbackSensitiveProcess)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.413s` (terminated `exit_code: 0`). Byte-identical recovery of the preceding red command.
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/cli` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	1.122s` (terminated `exit_code: 0`). Race tier after the last CLI implementation change.
> - `go test -tags integration ./cmd/fid/integration -run '^TestV01ClosureJourneyAndRebuildParity$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Focused reproduction of required-gate failure: removed requirement.sync single-line stderr expectation; package 0.715s.
> - `go test -tags integration ./cmd/fid/integration -run '^TestV01ClosureJourneyAndRebuildParity$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). Progress to the later unsupported implements relation's raw Human wording; package 0.950s.
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestV01ClosureJourneyAndRebuildParity|TestP23StopFailureHumanAndJSONParity|TestP23StopFailureCodeMatrix|TestProcessExitContract|TestP24Human.*Failure|TestP24GenericFallbackSensitiveProcess)$'` - FAIL - `FAIL` (terminated `exit_code: 1`). Same later relation assertion; package 6.607s. This failure remains evidence and was not routed elsewhere.
> - `go test -tags integration ./cmd/fid/integration -run '^TestV01ClosureJourneyAndRebuildParity$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.207s` (terminated `exit_code: 0`). Full closure journey now reaches migration, stale refresh and repeated rebuild checks.
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestV01ClosureJourneyAndRebuildParity|TestP23StopFailureHumanAndJSONParity|TestP23StopFailureCodeMatrix|TestProcessExitContract|TestP24Human.*Failure|TestP24GenericFallbackSensitiveProcess)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.958s` (terminated `exit_code: 0`). Final tagged pre-gate tier on the checkpoint tree.
>
> Scratch-only diagnostic failures were also retained as limitations, not gates: an initial provenance script tried to blame an untracked new file and exited 1; it was restricted to retained tracked assertions. A later heading-extraction read assumed DEC headings were H2 and exited 1; the actual H3 authority sections were read. Neither failed diagnostic supplies PASS evidence. No control, task contract or evaluator material was edited or used to repair these diagnostics.
>
>
> Real entrypoint observations, cwd /tmp: go build -o /tmp/fid-task10 ./cmd/fid had already terminated exit_code: 0 with no output before the direct probe. These expected refusals are runtime observations, not undeclared non-zero gate PASS claims.
>
> - `/tmp/fid-task10 stop raise --bogus` - OBSERVED REFUSAL - `next: see fid stop raise --help; template: fid stop raise --example` (terminated `exit_code: 2`). Its preceding lines were exactly error: unknown flag: --bogus and reason: invalid-argument; argument=--bogus, allowed=--actor|--example|--idempotency-key|--input.
> - `/tmp/fid-task10 --json stop raise --bogus` - OBSERVED REFUSAL (terminated `exit_code: 2`): empty stdout; operation stop.raise, code invalid-argument, invalid-input, retryable false, subject null, full spelling in argument/value, and the same four allowed flags. The decoded message was unknown flag: --bogus.
> - TestP24HumanProtectedDecisionFailure, TestP24HumanGateTerminalFailure and TestP24GenericFallbackSensitiveProcess exercised freshly built actual fid subprocesses in both modes; the focused tagged and full tagged commands below observed their terminations.
>
> QA-plan section 8 assigns D2 review and Human fixture regeneration after task-11. Those are independent Reviewer/Human work, not additional task-10 Developer commands. No held-out scenario or evaluator was read. Report-lint is the independent Runner control required after the report commit; it was not executed or claimed PASS here.
>
> ## Exact quality-gate results
>
> Attempt 2, final source tree; every exact command executed sequentially from `/home/mark/projects/fid`, with its own termination observed. BLOCKED rows preserve failed attempts rather than treating an incomplete or red test as the protected condition itself.
>
> - `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
> - `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`).
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Root guard audit is red; mapping/caller repair needs protected authority.
> - `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.6%` (terminated `exit_code: 0`).
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Root guard and incomplete registry corpus are red; fixture/completeness work remains distinct from protected authority.
> - `rg -n "func (write[A-Za-z]*Failure|render[A-Za-z]*Error)\(" cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/journal_errors.go:12:func renderJournalReplayError(writeLine func(string, ...any) error, err error, stopMessage string) (bool, error) {` (terminated `exit_code: 0`).
> - `rg -n 'CodeInvalidArgument, "invalid [^"]*"' cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/finding.go:144:	return writeGenericFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Finding request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})` (terminated `exit_code: 0`).
>
> Observed package wall times: cmd/fid 340.419s; cmd/fid/integration 352.186s; slowest other fresh ordinary package internal/export 112.885s. Internal/workflow and the remaining internal packages reported `(cached)`; no fresh numeric times are claimed for cached rows. Both fresh task-budget observations are under 600s/900s respectively; none reaches its 80% warning threshold. Full task gates remain failed despite those times. Coverage is 62.6%; the coverage-report command's successful parse does not change either failing test-gate status.
>
> Historical exact attempts follow unchanged:
>
> Continuation gate (third full sequence across the durable Task 10 history), cwd /home/mark/projects/fid. All seven commands ran byte-identically, in this order, without concurrent heavy test work. Each receipt below was captured in this session and its nested result supplied the termination. The source/test tree remains the original checkpoint.
>
> - `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> - `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	120.190s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
> - `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.6%` (terminated `exit_code: 0`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	361.027s` (terminated `exit_code: 0`)
> - `rg -n "func (write[A-Za-z]*Failure|render[A-Za-z]*Error)\(" cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/activation.go:323:func renderActivationError(cmd *cobra.Command, operation string, cause error) error {` (terminated `exit_code: 0`)
> - `rg -n 'CodeInvalidArgument, "invalid [^"]*"' cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/finding.go:144:	return writeGenericFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Finding request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})` (terminated `exit_code: 0`)
>
> Fresh required timings: cmd/fid 412.776s; internal/workflow 120.190s; slowest other ordinary package internal/export 197.915s; cmd/fid/integration 361.027s. Both rg commands observed the declared matches and exit_code: 0. No required command failed in this continuation. Passing these existing gates does not close the independently red activity transport or missing acceptance corpus.
>
> Historical gate attempts 1 and 2 follow unchanged:
>
> Both full attempts executed all seven byte-identical task commands in order from the repository root. Attempt 1 is immutable failed evidence; attempt 2 is the final checkpoint gate. The rg commands both expect and observed exit_code: 0 with declared matches; no narrower command substituted for any gate.
>
> Attempt 1:
>
> - `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> - `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`)
> - `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.6%` (terminated `exit_code: 0`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
> - `rg -n "func (write[A-Za-z]*Failure|render[A-Za-z]*Error)\(" cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/label_response.go:11:func renderLabelError(cmd *cobra.Command, operation string, cause error) error {` (terminated `exit_code: 0`)
> - `rg -n 'CodeInvalidArgument, "invalid [^"]*"' cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/finding.go:144:	return writeGenericFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Finding request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})` (terminated `exit_code: 0`)
>
> Attempt 2:
>
> - `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> - `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	112.664s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
> - `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.6%` (terminated `exit_code: 0`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	362.561s` (terminated `exit_code: 0`)
> - `rg -n "func (write[A-Za-z]*Failure|render[A-Za-z]*Error)\(" cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/sql.go:51:func renderSQLCommandError(writer io.Writer, err error) error {` (terminated `exit_code: 0`)
> - `rg -n 'CodeInvalidArgument, "invalid [^"]*"' cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/finding.go:144:	return writeGenericFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Finding request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})` (terminated `exit_code: 0`)
>
> Final fresh package times: cmd/fid 409.529s; internal/workflow 112.664s; slowest other ordinary package internal/export 198.511s; cmd/fid/integration 362.561s. Attempt 1's ordinary cmd/fid time was 350.670s; its workflow/issues and other package lines were cached, so no fresh time is claimed for those cached rows. Attempt 1 tagged integration reported 357.104s with exit_code: 1.
>
> ## Deviations from spec
>
> Attempt-2 deviations:
>
> - DEC-285 resolves the prior activity-producer question; its original failed transport observations and report-lint non-complete evidence are preserved, with no deletion or rewritten verdict. First six-pair focused command printed `FAIL` (terminated `exit_code: 1`); root cause was raw activity errors plus six forced-JSON caller arguments. The bounded three-file repair made the byte-identical six-pair command pass (`1.432s`, terminated `exit_code: 0`). Direct CLI activity --since observations now exit 2 in both modes with the accepted transport.
> - R20 completeness initially printed `FAIL` (terminated `exit_code: 1`) for two omitted sites; adding real epic-scheduling and blank-link-ID domain rows retained the assertion and completed all eighteen sites. First plain-copy pair run printed `FAIL` (terminated `exit_code: 1`) because copied Git stat data refreshed .git/index on observation. A specific Run rerun confirmed the index-only change. Normalizing copied index stat data before the snapshot made the byte-identical full four-pair command pass (`1.786s`, terminated `exit_code: 0`); snapshot assertions retain their full strength.
> - The initially narrow no-bypass audit passed but was incomplete. Its strengthened root-guard check is red (`FAIL`, terminated `exit_code: 1`); actual-binary Human/JSON refusals both exit 7 with raw text. No mapping was chosen and no forbidden caller was edited to turn it green. Protected authority for this distinct gap is required.
> - Registry missing rows and unfinished fixtures remain red rather than being assigned exclusions or another task. Proper export.import fixtures use accepted --repository without --project, while a separate regression explicitly preserves the invalid root-flag failure. The two budget admission rows currently fail the empty-stdout expectation because JSON emits a fid.live.v1 terminal denial; accepted streaming transport is unchanged and no new stdout exception is invented.
> - Report lint is expressly the Runner's control in Task 10 and QA-plan section 5. It is not run or claimed green by this Developer. Prior failed control evidence remains verbatim in the historical sections. No required task command is substituted by a focused selection.
> - The preserved prior envelope-prerequisite attempt ended with `FAIL activity.list JSON failure envelope absent; stdout=''; stderr='error: activity list: --since must be YYYY-MM-DD or RFC3339\n'; observed exit_code=7` (terminated `exit_code: 1`). Its root cause was the raw activity.list path with absent JSON transport. DEC-285 uniquely supplied the accepted mapping; only the three authorized callers were repaired. The byte-identical six-row focused command now passes (1.432s, exit_code: 0), and direct Human/JSON activity --since entrypoint observations both show exit_code: 2 with the accepted failure. Prior report-lint OUTCOME RECOVERY_REQUIRED / exit_code: 10 and OUTCOME BLOCKED / exit_code: 11 remain attributed to the Runner and are not rewritten as green.
> - Mandatory current required-command failed-attempt note: exact ordinary race/coverage command printed `FAIL` (terminated `exit_code: 1`; cmd/fid 340.419s), caused only by the new audit of the retained raw root guard. Exact tagged integration command printed `FAIL` (terminated `exit_code: 1`; 352.186s), caused by that guard plus seventeen missing registry rows and unfinished activation/projection/budget fixtures. The bounded DEC-285 caller changes pass their six rows; no root mapping, gate/control definition, assertion semantics or stdout exception was changed. No later green is claimed for either full command. Root repair requires protected authority; the distinct corpus remainder is ordinary same-task work, not a separate protected condition or a routed task.
>
> Historical deviations follow unchanged:
>
> Current continuation:
>
> - Protected boundary: TRG-24-009's Recovery scope says, "If the activity repair requires inventing a failure mapping or changing the accepted JSON contract, return to triage rather than inventing it." The observed JSON request emits raw prose with exit 7, and a missing-project request emits nothing with exit 7. activity.go reads jsonMode only after List succeeds, and its failure path constructs no cli.Failure. Existing renderIssueReadError/issueReadFailure would select invalid-argument with request-level details and exit 2 for the filter cause, project-not-found/exit 3 for missing project, or journal-replay-stopped/exit 6 for replay failures. Those are changed exits and newly selected producer mappings. Constructing internal-error solely from the raw exit 7 would also invent a producer mapping and message/details contract. The shared renderer requires a typed Failure; an empty fabricated Failure cannot supply accepted semantics.
> - Resolution needed: identify the authoritative existing activity failure mapping, or record precise protected authority reconciling the missing JSON producer and its accepted code/details/exit with the preservation boundary. No accepted mapping was guessed, code changed, assertion weakened, or success output altered.
> - The two finite corpora and four-leaf built-binary plain-copy pairing remain incomplete. They are not labelled BLOCKED merely because they are unfinished; the protected outcome comes from the separate activity producer constraint above.
> - The prerequisite check's last token is preserved in Tests run with exit_code: 1. Its cause is the actual raw JSON transport, and no later green is claimed. This is supplemental failed evidence, not a required-command failure or an authorized change to the task's assertions.
> - The original report-lint non-complete result and all historical red/green/gate receipts remain preserved. This session's exact gates all passed; there is no new failed required-command attempt to erase.
> - Developer report-lint was not run: Task 10 explicitly assigns that control to the Runner, and the invoking instructions require the Runner to run it after this finalized report commit. D2 re-review and the Human walkthrough remain mandatory after task 11; neither is waived.
>
> Historical deviations follow unchanged:
>
> - The task is unfinished: complete RegisteredCodes and R20 domain-process completeness registries, Human plain-copy process pairing, and the no-bypass caller audit remain. These are RECOVERY_REQUIRED work; a green automated gate does not establish these missing acceptance criteria.
> - Mandatory required-command failed-attempt note: gate attempt 1's exact tagged integration command ended with `FAIL` (terminated `exit_code: 1`, 357.104s). Task-10's accepted second Human line displaced assertRemovedCommandFailure's single-line expectation (introduced by task-9a). Preserving the invocation and typed JSON/journal checks, the in-scope assertion was converged. The focused rerun then reached a later unsupported-link assertion introduced by phase-23/task-12; task-10's accepted fixed JSON sentence displaced its raw-domain Human substring. Its exact canonical error/reason expectation now retains the relation value, while the original no-append checks remain. The focused journey and combined tagged race checks passed; gate attempt 2 reran the byte-identical full tagged command with exit_code: 0. No failure was routed to another task.
> - The first ordinary touched-package race attempt ended with `FAIL` (terminated `exit_code: 1`, cmd/fid 308.286s); task-10 displaced the retained governance single-line assertion and removed-endpoint raw wording. Provenance was checked, canonical lines added, and the original cause/empty-stdout/journal checks retained. The byte-identical touched-package command passed (cmd/fid 307.357s; cli 1.115s).
> - The first focused tagged race attempt ended with `FAIL` (terminated `exit_code: 1`, 6.555s). The schema-based nested-detail filter incorrectly suppressed an actual scalar null. The nullable-object unit was first red, then the one-line in-scope value filter repair restored event=null without displaying non-null nested objects/arrays. The identical tagged selection passed in 8.413s; cli's race tier then passed in 1.122s. This changed no failure schema or assertion contract.
> - Other red attempts retained in Tests run cover the specified protected-Decision first failure, missing renderer API, generic cause construction, empty nested array and 12-hex repository reason, terminal gate stderr, displaced retained output, and fixture/script construction mistakes. No test was deleted, command narrowed, pass condition changed, or flaky green mined.
>
> ## Concerns, technical debt, or notes for the Tech Lead
>
> Attempt-2 protected candidate and checkpoint notes:
>
> - Resolution needed: identify the accepted typed root-flag failure mapping (JSON code/details/exit and section-9 derivation) and authorize bounded repair of rejectUnavailableRootFlags. Task 10's accepted universal Human failure change exposes this retained raw guard. DEC-285 only authorizes three named caller extensions and cannot settle the absent root mapping; no authority is inferred from a failed test.
> - Provenance is confirmed independently: root guard introduced by `cd53c43319a3301d3a8cc0089b6785d9afbb163f`, `feat(project): add strict project selection config`, `Refs: phase-20/task-2`; export.import capability introduced by `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, `feat(export): expose import and prove portable recovery parity`, `Refs: phase-22/task-15`. Both individual `git merge-base --is-ancestor <hash> 304aa45f12cf3966567231e6986e3070df409c68` commands silently terminated exit_code: 0 in this session. Historical provenance does not excuse convergence or assign the failure to another task.
> - RegisteredCodes omits the defined reachable idempotency-conflict member (codes.go definition versus registeredCodes slice). A supplemental process pair retains its actual next-table coverage without changing the protected enumeration or calling it a registry member. This inventory concern needs the owning authority, not a silent registry edit in the Human-renderer scope.
> - Remaining decision-free work: seventeen registry members need domain rows (claim-mode-invalid, obligation-unsatisfied, finding-not-found, human-outcome-required, outcome-not-permitted, review-changes-required, review-verdict-invalid, result-limit-exceeded, source-digest-mismatch, evidence-type-not-accepted, finding-already-disposed, concern-already-disposed, outcome-not-found, authority-missing, review-not-independent, concern-not-found, deviation-not-found); activation-conflict and activation-incomplete fixtures need accepted descriptor fields, projection warning-array oracle needs its accepted typed nested shape, and budget pairs need a contract-consistent frame treatment. Partial tests are committed as red checkpoint evidence, not passed acceptance. Complete these in the same Task 10 after authority disposition; no other task is nominated.
> - No retained assertion was displaced by this session's three product edits; full ordinary/tagged gates reveal any additional displacement, and the prior convergence/provenance table remains unchanged. No retained test was deleted, narrowed or weakened.
> - Headroom warning: none for current fresh observations; cmd/fid 340.419s is below 480s and tagged integration 352.186s is below 720s. Internal/workflow timing is UNVERIFIED fresh because the exact gate reported (cached).
> - Only this task's files are staged. Runner-owned run-ledger.md was the sole entry dirt and remains unchanged by the Developer. Committed task/report paths must be clean; the full worktree still includes that prohibited-to-edit Runner dirt.
>
> Historical concerns follow unchanged; DEC-285 supersedes their activity/stats/AC-list authority questions:
>
> Current protected candidate and scope notes:
>
> - Activity mapping authority is needed for the exact raw producer paths at activity.go:35,46,49. This is the explicit return-to-triage condition in TRG-24-009, not a request to expand its one-file extension independently. Its file-scope question is resolved; its runtime JSON-producer prerequisite is absent.
> - Stats tokens/cost at stats.go:38,48,73,83 and issue.ac.list at issue_ac.go:109,118 have verified Human-mode JSON failure bypasses. The shared mapper is adopted, but these callers force its JSON branch. Both exact caller files lie outside the named R9-function scope and only activity.go has the authorized extension. Record bounded scope disposition within the same task; no other task is assigned and no file scope is expanded by this report.
> - Remaining decision-free work after authority resolution: complete RegisteredCodes row/exclusion completeness, the eighteen R20 domain-process rows/honest unreachable exclusions, and actual-binary envelope.create/claim.acquire/run.start/gate.run plain-copy Human/JSON pairs while retaining the valid Claim and original/copy/import controls. The closed contract is unchanged.
> - Headroom warning: none for this continuation; maximum ordinary package is 412.776s, below the 480s warning. Tagged integration is 361.027s.
> - Entry branch and HEAD were phase-24 at f1ad50fa18eb45b308cc52ffb3ea0c65b6dd198d. Both run-ledger.md and triage-log.md were dirty on entry, despite the recovery prose saying ledger-only dirt. Their initial contents were retained; only the report is staged/committed. No protected control, generated artifact, dependency, task-11/task-13 surface or held-out material was edited/read.
> - This is a Developer BLOCKED candidate for the Phase Runner's triage decision. It does not rewrite the ledger state or allocate any TRG/ESC.
>
> Historical checkpoint notes follow unchanged (the previously requested activity file extension is now granted by TRG-24-009, with the protected mapping condition above still unresolved):
>
> - Same task-10 continuation is required. Complete the two finite process corpora and their honest row/exclusion completeness checks, including exact both-mode plain-copy messages with the valid Claim prerequisite; retain all checkpoint evidence.
> - Scope gap: cmd/fid/activity.go:46 calls the R9 renderIssueError(io.Writer,error) without operation/leaf context; its handler is not a named R9 function. The uniform task-10 accepted change exposes this stable caller's raw failure path. Obtain the bounded task-scope disposition before editing that caller. The call was introduced by 977974c0ec00d3d9a63f7d7515e633683c9b9ac5 (phase-16/task-1), and git merge-base --is-ancestor confirmed it predates phase-24's approved task-pack commit 304aa45f12cf3966567231e6986e3070df409c68. This is not assigned to task-11 or task-13. Also audit stats.go's explicit JSON-port calls against their accepted caller contract before claiming no bypass.
> - Headroom warning: none; every freshly timed ordinary package is below 480s.
> - The Runner-owned dirty run ledger was excluded from all edits/staging. Report-lint remains Runner-owned and must run after this committed report.
>
> ## Dependency-budget note
>
> - No dependency changes.
>
> ## Commit hash
>
> - `6607143b68ab2d08444cf5d5ecf9ddfd84ff7f49`
>
> ## Commit subject
>
> - `feat(cli): checkpoint task 10 caller recovery and failure corpora`

</details>

## Exact quality-gate results

All seven byte-identical Task 10 commands ran sequentially from `/home/mark/projects/fid` after gofmt/goimports. Each receipt is fresh output from its own completed nested command. There was no concurrent heavy test work. The test-result cache was cleared before this sequence to obtain fresh package wall times; cache cleanup is not gate evidence.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	116.440s	coverage: 62.1% of statements` (terminated `exit_code: 0`)

- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.6%` (terminated `exit_code: 0`)

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	335.908s` (terminated `exit_code: 0`)

- `rg -n "func (write[A-Za-z]*Failure|render[A-Za-z]*Error)\(" cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/label_response.go:11:func renderLabelError(cmd *cobra.Command, operation string, cause error) error {` (terminated `exit_code: 0`)

- `rg -n 'CodeInvalidArgument, "invalid [^"]*"' cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/gate.go:327:		return writeGenericFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid gate request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted gate registration input"}})` (terminated `exit_code: 0`)

Fresh recorded package wall times: cmd/fid 484.734s; internal/workflow 116.440s; internal/issues 216.748s; cmd/fid/integration 335.908s. Headroom warning: cmd/fid 484.734s.

## Deviations from spec

- The final report diff check first ended `+> ` (terminated `exit_code: 2`) because consecutive historical blank-line quote prefixes retained one trailing space. Blank quote lines now use `>`; decoding the quotation still reconstructs every incoming report byte. This documentation-only correction changes no source, test, required gate or historical evidence.
- DEC-286 resolves the prior protected root mapping; it does not erase the correct prior BLOCKED outcome or raw exit-7 probes. Current in-process root red printed `FAIL` (terminated `exit_code: 1`, 0.170s), and the process root red printed `FAIL` (terminated `exit_code: 1`, 1.260s). The cause was the retained raw root guard. The bounded root.go change and pure shared extractions make the byte-identical in-process command green (0.173s, exit_code 0), and the combined root/registry command green (20.582s, exit_code 0). Exact JSON/Human bytes, diagnostics and no-effect checks remain explicit.
- Mandatory prior required-command failed-attempt note: the incoming exact ordinary race/coverage gate ended `FAIL` (terminated `exit_code: 1`, cmd/fid 340.419s), caused by the new root guard audit; the exact tagged gate ended `FAIL` (terminated `exit_code: 1`, 352.186s), caused by that guard plus seventeen missing registry members and unfinished activation/projection/budget fixtures. DEC-286 supplies the missing mapping and allows the bounded repair. This continuation completes those own corpus fixtures and reruns both byte-identical full commands; their fresh PASS receipts above close the failed attempts without changing any command/threshold/assertion semantics or assigning a different task.
- Current combined corpus reds are preserved: `FAIL`/exit_code 1 at 13.575s for missing rows and invalid activation descriptor/array oracle/frame fixtures; a build `FAIL`/exit_code 1 for a duplicate helper name; `FAIL`/exit_code 1 at 1.416s for reading terminal data as its outer frame; `FAIL`/exit_code 1 at 22.570s for invalid projection state and unmet domain prerequisites; and `FAIL`/exit_code 1 at 27.060s for two still-unsatisfied Human obligations. The focused four-row red (`FAIL`, exit_code 1, 2.263s) exposed missing --input on the real gate invocation; the focused Human row red (`FAIL`, exit_code 1, 1.758s) exposed the need to satisfy release approval first. Minimal fixture repairs use accepted fields, pending resource state, real committed terminal decoding, actual activated authority/Claim/Review/evidence and real gate input. The single Human row is green at 2.059s, and the identical combined command is green at 20.582s. No earlier red is recast as an exclusion.
- activation-conflict corrects the own fixture's descriptor field to expected_prior_version. activation-incomplete now reaches the same domain subject through a real committed activation with lost resource agreement at finalization; the previous empty component-digests request failed structurally before it could test that subject. DEC-286 clause 10 authorizes completing this fixture. The full pairing assertion remains exact, without stubbing or weakening domain checks. Projection reason computation follows the accepted object-array detail type, including empty warnings/errors arrays; product renderer and details are unchanged.
- Budget-data-unavailable and budget-exhausted reach gate.run's existing committed denials. Phase 22 public-interface section 4 (lines 252-260) explicitly requires the singleton fid.live.v1 terminal fid-error carrying the same response, no start/output/spawn and durable replay; Phase 21 section 1 preserves live transport. These rows therefore compare that terminal response to Human stderr, assert no Human stdout, no JSON stderr, exact exit 4, terminal sequence/outcome/null Attempt/Receipt and unchanged replay filesystem. This is retained terminal-contract convergence, not a newly invented stdout exception. Ordinary non-live pairs still require empty stdout in both modes; gate-unsuccessful retains its streamed Human stdout assertion.
- DEC-285's prior six-pair red `FAIL`/exit_code 1 at 3.813s, its bounded activity/live-jsonMode repair and green 1.432s receipt remain verbatim in historical evidence; the same six pairs pass freshly in the tagged race selection and full integration gate. Historical R20 omission and copied Git-index snapshot failures/repairs and every prior mandatory gate deviation remain verbatim below the historical evidence summary. No accepted caller or remaining failure is routed elsewhere.
- The historical report-lint failures are independent Runner evidence, retained with their original tokens and exits; no current Runner PASS is claimed. No command in the seven-command gate is narrowed or substituted. The role-required Developer report check is separate from the Runner's control.

## Concerns, technical debt, or notes for the Tech Lead

- RegisteredCodes still omits the defined reachable idempotency-conflict member. Its supplemental actual-process pair covers the next-table contract without changing the protected registry. The omission remains an inventory concern for the owning authority.
- QA Window D2 after Task 11, including all DEC-285/DEC-286 rows and grammar recheck, plus the Human walkthrough remain mandatory. VAL-054 is not attested by this Developer.
- Headroom warning: cmd/fid 484.734s.
- All Developer task/report changes are committed at handoff. Runner-owned run-ledger.md is intentionally retained as its sole preexisting dirt and is neither modified, staged nor hidden; this is the explicit dispatch boundary, so the Phase Runner must finish its own ledger transaction before the entire phase worktree can be clean.
- The complete original failure/report history is retained; current COMPLETE describes only the fully verified Task 10 acceptance and exact gates. The Phase Runner independently verifies the finalized committed report before TASK COMPLETE. No ledger transition or phase advance is claimed here.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `d11bec713c470596131190c1ca9c33274fd6133c`

## Commit subject

- `fix(cli): complete task 10 root-flag failure mapping and corpus`
