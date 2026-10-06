# Phase 24 - Mid-Phase Review 4 (Tasks 4-6)

Reviewer: qa-reviewer
Date: 2026-10-06
Tasks covered: phase-24/task-4, phase-24/task-5, phase-24/task-6
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-24/qa/review-3.md

## Verdict

PASS. Fresh inspection and affected verification resolve all eight original Window B Must-fix findings, F1-F8. No Must-fix or blocker candidate remains. Two Runner bookkeeping Improvements are deferrable. This verdict closes the Window B review requirement; the Phase Runner owns the next dispatch. Task 7 was neither dispatched nor reviewed.

## Scope Reviewed

- Reviewed committed product, tests and reports at `86f18761c72633df12ed99b2fb2eb4bcc27bf3ea`, against review-3's repair boundaries and the accepted task-4, task-5 and task-6 specifications.
- Loaded `AGENTS.md`; current `phase_brief.md`, `decisions_active.md`, `progress_current.md` and `tl_escalations_current.md`; phase-24 `qa/validation-contract.md`, `qa/qa-plan.md` (Window B and applicable evidence/gates), `qa/review-3.md`; and `project-docs/templates/mid-phase-review-template.md`.
- Inspected all three completion reports, phase decisions including DEC-283, current triage entries TRG-24-004/005, and the Runner ledger. Consulted `public-interface.md`, `input-inventory.md`, `08_QUALITY_PLAN.md`, `09_AI_WORKFLOW.md` sections 2.7/8, the `qa-evidence` skill, and relevant CLI/flow/glossary documentation. CLI parity observations use `05_DESIGN_PLAN.md` sections 1-3, 6 and 7. No held-out material was read.
- Task 4 repair: `c0115bd500c1bedd0384b13e5eb352a64c26c865`, `test(judgment): complete repository and descriptor QA evidence`; report chain `a8672406d387cb88c94c6c9532395e6c44d12a6c`, `a9a95bdfdb55e14b9445c0f36762484563192979`, ending `6352f94c564cbf248980e3056f355a40f611c95d`.
- Task 5 repair: `04246029c4488eec46eba76323dee7fbe8260c8d`, `test(review): checkpoint Window B evidence repair`; recovery report `a0ca854133ae98b77e36cf51ad596ff0e8c6e096`; causal repair `226cdc94e6460cc254d1b4e916344b8ea2e2cdb3`, `fix(cli): enforce declared Review OID formats`; final report `be56892031fd3bfc9b152598364a6e4e9079c973`.
- Task 6 report-only repair: `86f18761c72633df12ed99b2fb2eb4bcc27bf3ea`, `docs(task): add task 6 R5 inventory receipt`. Retained implementation `01e518f460d7339848c6549ecb7dbd4d58a2f6f4`, `test(binding): converge retained task 6 execution fixtures`, and earlier DEC-283 maintenance `13164582232b9b43eddac93a0ef75b6f4465783c`, `fix(workflow): reconcile STOP resume recovery tuples`, were also inspected.
- Every repair commit has its owning `Refs: phase-24/task-N`. Diffs since review-3 change seven files: the shared descriptor test, two integration test files, `internal/cli/input_validate.go`, and the three reports. The current ledger and triage log were already dirty at invocation; these are live Runner evidence, distinct from the committed product/report tree, and remain untouched by this reviewer.

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 2
- Blocker candidates: 0
- Existing escalation references: ESC-24-001 and ESC-24-002, closed by DEC-283; none open
- Covered validation assertions: VAL-006, VAL-007, VAL-008, VAL-009, VAL-010, VAL-011, VAL-012, VAL-013, VAL-014, VAL-015, VAL-016, VAL-017, VAL-018, VAL-019, VAL-020, VAL-021, VAL-022, VAL-023, VAL-024, VAL-025, VAL-026, VAL-027, VAL-028, VAL-029, VAL-035, VAL-052, VAL-060, VAL-066

## Validation Contract Coverage

ON TRACK applies to the tasks 4-6 slice, including both qualifying and refusing states where applicable. Whole-tree descriptor completion remains task 9; later observing leaves retain tasks 7/8; domain-cause Human messages retain task 10. These allocations do not defer any Window B evidence. VAL-010 is retained framework context, independently checked by the inventory walk, rather than a reopened Window A verdict.

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-006 | task-2, converged by tasks 4-6 | ON TRACK | I1 reran the judgment keyed conflict/canonical-retry corpus and shared leaf conformance; differing-field/no-Event subjects remain. |
| VAL-007 | task-2, converged by tasks 4-6 | ON TRACK | I1 keyed retries retain original result/Event equality; binding and Review replay equality also reran in C1/I1. |
| VAL-008 | task-4 Source case | ON TRACK | Judgment changed-file retry corpus checks supplied versus omitted digest conflicts naming `source`/`references`, and unchanged-input replay; I1 passed. Earlier optional-array evidence remains Window A context. |
| VAL-009 | tasks 4-6 | ON TRACK | All owned observing mutations have changed-HEAD/file retries with complete original-result equality and no Event. Review's added whole-response comparison and advanced-HEAD fresh control pass I1; judgment/binding rows pass I1/C1. |
| VAL-010 | task-3, retained context | ON TRACK | Fresh recursive live walk: 92 leaves, zero help outputs starting `Usage:`. |
| VAL-011 | tasks 4-6 | ON TRACK | C1 live-tree/Go-type reflection and seeded mismatch rejection pass; seven judgment, two Review/Finding and eight binding descriptors are adopted. |
| VAL-012 | tasks 4-6 | ON TRACK | C1 recursively checks exact paths, types/element types, requirements, meanings and applicable constraints for every adopted leaf. Corrupted-metadata guards reject; Review help names both Git commands. F3/F7 resolved. |
| VAL-013 | tasks 4-6 | ON TRACK | C1 emitted examples equal ordered templates with two-space indentation and required-only members; Review/Finding output is checked through a built process. |
| VAL-014 | tasks 4-6 | ON TRACK | Closed request-type/descriptor checks pass C1. Review/Finding additionally check recursive required-only keys and JSON types and decode with unknown fields disallowed. `<head-oid>` stays a typed template placeholder; executable controls use a real OID. |
| VAL-015 | tasks 4-6 | ON TRACK | C1 built-process exclusivity checks all six prohibited combinations per adopted leaf, exact errors and no effects, plus standalone success. |
| VAL-016 | tasks 4-6 | ON TRACK | C1 exact recursive structural corpora exercise applicable unknown/duplicate/missing/null/type/enum/format and bound/cardinality classes at top/deep paths; accepted controls reach domain dispatch without effects. Exact malformed-OID process refusals pass. F3/F7 resolved. |
| VAL-017 | tasks 4-6 | ON TRACK | Judgment and Review caller-repository paths reject in I1. Binding removed-member tests remain present, and C1 binding process/Source corpus exercises their closed descriptors and no effects. |
| VAL-018 | tasks 4-6 | ON TRACK | I1 covers judgment positions and all three Review positions with independently hashed matching/mismatching digests, exact mismatch details and no Event. C1 covers all three Envelope Source positions. F5 resolved. |
| VAL-019 | tasks 4-6 | ON TRACK | Omitted-existing acceptance and omitted-missing refusals appear for each owned Source position. Review matrix has both states for items, exclusions and evidence refs; missing details have both digests null. I1/C1 passed. |
| VAL-020 | task-4 | ON TRACK | I1/P1 cover absolute/parent paths, escaping symlink, directory and accepted in-root symlink; typed paths, null observed digests and no-Event refusals retained. |
| VAL-021 | task-4 | ON TRACK | I1/P1 reject an omitted external digest and accept a supplied digest verbatim without a local file. |
| VAL-022 | tasks 4-6 | ON TRACK | Seven judgment leaves compare independent Git/config/root observation with every applicable record/Source and actual Event repository, with pinned non-vacuous counts. Review compares all eight Repository fields on Review, three Sources and Event, including an independently hashed dirty image. Binding tests retain their independent observation comparisons. F2/F6 resolved. |
| VAL-023 | tasks 4-6 | ON TRACK | Judgment seven-leaf and Review normal/non-Git/unmerged pairs pass I1; four binding-leaf pairs pass C1. P1 reproduces all five classifier cause classes. Unsafe states have mapped reasons, exit 4 and no Event. |
| VAL-024 | task-5 Review slice | ON TRACK | I1 Review journey preserves matching, HEAD-different and dirty-different states with no Event on refusal; P1 checks only HEAD/dirty participate in the shared expectation check. |
| VAL-025 | task-5 | ON TRACK | I1 exact expected/observed mismatch comparisons and P1 registered closed-detail schema test pass. |
| VAL-026 | task-6 | ON TRACK | C1 Registration oracle independently obtains Git/config/root identity and recomputes binding digest, then compares Registration/Event. |
| VAL-027 | task-6 | ON TRACK | C1 pairs original-root Envelope creation with whole-worktree plain-copy refusal and genuine imported-destination acceptance. Fixed fallback message is correct until task 10's planned convergence. |
| VAL-028 | task-6 | ON TRACK | C1 original-root Claim acquisition versus plain-copy refusal, observed worktree identity and no-Event comparisons pass. |
| VAL-029 | task-6 slice | ON TRACK | C1 real export/import, original/copy controls and imported Claim renew/release/expire pass. Matching project/target-root binding is required. Later Run/Gate rows and task-10 messages retain their owners. |
| VAL-035 | tasks 4-6 | ON TRACK | I1 shared leaf/amended-leaf and dual-adapter recovery conformance pass; converged fake/fid rows retain acceptance and removed-member refusal parity. |
| VAL-052 | tasks 4-6 | ON TRACK | Committed gate timings are below every budget/warning; fresh and cached measurements are distinguished in the timing table below. Focused QA times do not replace full-gate measurements. |
| VAL-060 | tasks 4-6 | ON TRACK | C1 built-process standalone examples return exit 0, empty stderr and unchanged empty non-project directories; binding corpus also compares tree state. |
| VAL-066 | tasks 4-6 | ON TRACK | Judgment/binding unavailable/deleted retries retain original results and no Event. Review has normal/non-Git/unmerged/deleted-source pairs with fresh unreviewed Run/key controls, exact refusals, whole-response/Event replay equality and identical restored-request success. I1/C1 pass; F6 resolved. |

## Gate And Evidence Checks

Fresh reviewer commands ran from `/home/mark/projects/fid`. Built-process tests allocate temporary workspaces. The recursive inventory block ran from the repository root and wrote only fresh scratch files under `/tmp/tmp.tKreuYuHbD`. Every long command was polled through its retained nested session until an actual `exit_code` was returned. The results below are fresh reviewer observations. Developer gate rows discussed separately are historical receipts in the committed reports.

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Representative exact required command | `gofmt -l cmd internal` | PASS | No output; token `exit 0 (no output)`; observed `exit_code: 0`. |
| C1 affected ordinary verification | Exact command below | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t45.878s`; observed `exit_code: 0`. |
| I1 affected integration verification | Exact command below | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t29.612s`; observed `exit_code: 0`. |
| P1 Source/expectation/classifier/closed-details verification | Exact command below | PASS | Project `1.077s`, Review `1.018s`; final token `ok  \tgithub.com/mbell-nzl/fid/internal/cli\t1.125s`; observed `exit_code: 0`. |
| Task 4 final report-lint | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-4-done.md project-docs/docs/phases/phase-24/tasks/task-4.md` | PASS | Fresh `OUTCOME COMPLETE`; final token `PASS report-lint completed with no failures`; observed `exit_code: 0`. |
| Task 5 final report-lint | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-5-done.md project-docs/docs/phases/phase-24/tasks/task-5.md` | PASS | Fresh `OUTCOME COMPLETE`; same final token; observed `exit_code: 0`. |
| Task 6 final report-lint | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-6-done.md project-docs/docs/phases/phase-24/tasks/task-6.md` | PASS | Fresh `OUTCOME COMPLETE`; same final token; observed `exit_code: 0`. |
| Required command coverage | Specs/QA plan versus reports' Tests run and Exact quality-gate results | PASS | Each report has the five byte-identical final PASS rows, in order, each with its own token and recorded observed termination. Owned inventories now have separate receipts. |
| Commit/report integrity | `git rev-parse --verify <hash>^{commit}`, subjects, Refs, diffs and index flags | PASS | All 55 distinct full hashes extracted from reports resolve; every child verification returned exit 0. Primary subjects match. All reports have empty `git diff HEAD` and normal `H` index flags; fresh report-lint verifies committed clean files. |
| Failure preservation | Quoted reports versus committed prior reports; Task 6 repair diff | PASS | Task 5's 238-line prior report reconstructs byte-for-byte; Task 6's 507-line historical block does too. Task 4 preserves all 99 original lines in order with an additional trailing blank quoted line. Task 6 F8 has additions only. All earlier failed gates/tokens/terminations and root-cause/deviation notes remain. |
| Introducing provenance | Historical introducing diffs and Refs versus each converged row | PASS | Task 6's 21 retained rows and four final rows resolve to their reported paths/phases/tasks and execute under the assigned ordinary/tagged gates. Tasks 4/5 retained convergence histories and OID descriptor/validator provenance were checked against original commits, rather than filenames. |
| DEC-283 and recovery-chain control | `sh project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-24/run-ledger.md` | PASS | Fresh token `VALID state=REREVIEW_PENDING`; observed `exit_code: 0`. This correctly awaits the new independent verdict. |
| TRG-24-004/005 directive binding | Current triage versus bound TASK-RECOVERY entries | PASS | Ordered allowlists, forbidden changes, exact verification, report/base, trigger, attempt and authority are byte-identical; Repair category maps exactly to Category. Their linked repairs are RESOLVED and COMPLETE. |
| Owned R5/R6 inventories | Exact searches below and committed recording-session receipts | PASS | R5 six matches; R6 nine matches. Fresh reviewer searches each returned observed exit 0. Every report receipt's token occurs in the current set; search ordering can vary. F1/F4/F8 resolved. |
| All inventory registries | Accepted `input-inventory.md` section 1 block and R1-R20 regeneration | PASS | Walk: 92 leaves, 29 groups, 39 input leaves, 50 keyed leaves, zero no-Short leaves, 42 typed flags; final token `252 total`, child/wrapper exit 0. Registry searches all observed exit 0 with accounted matches. R8/R12 are source registries. |
| Dependency/scope/doc checks | Repair diffs; relevant production and CLI/flow/glossary text | PASS | No dependency/spec changes since review-3. The sole product repair is the two-line OID dispatch/fail-closed change authorized by TRG-24-005; Review descriptor and repaired integration file remain unchanged in that continuation. |

Exact fresh named-test commands:

```bash
# C1
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestAdoptedExampleProcessExclusivity|TestP24JudgmentDescriptorSourcePositionsAndExamples|TestP24JudgmentStructuralValidationCorpus|TestP24JudgmentEvidenceRejectsIncompleteMetadata|TestP24ReviewAndFindingDescriptorHelpAndStructure|TestP24ReviewHelpEvidenceRejectsIncompleteMetadata|TestP24ReviewMalformedOIDProcess|TestP24OIDFormatDispatch|TestP24UnknownInputFormatRejects|TestP24BindingDescriptorProcessCorpus|TestGateRegisterObservesRepositoryWithoutCallerIdentity|TestP24EnvelopeReplayAfterSourceChangeOrDeletion|TestP24BindingReplayBeforeObservation|TestP24EnvelopeAndClaimOriginalCopyAndImport|TestP24ImportedClaimStoredAuthority|TestP24EnvelopeSourcePositions|TestExecutionBound.*|TestExecutionTransitionFixtureDoesNotInheritExpiredClaim)$' ./cmd/fid
# I1
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24HumanOutcomeRejectsReferenceDigestMismatch|TestP24JudgmentRepositoryOracleUsesConfig|TestP24JudgmentSourceDigestsAndRepositoryParity|TestP24JudgmentSourceInputValidation|TestP24JudgmentLeavesRepositoryUnavailable|TestP24JudgmentLeavesReplayAfterHeadAndFileChange|TestP24ReviewSourceMatrixCompleteness|TestP24ReviewSourceMatrix|TestP24ReviewReplayEvidenceRejectsChangedHistory|TestP24ReviewReplayUnavailableAndRestoredControls|TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation|TestP23ReviewRecordRealProcessJourney|TestP24JudgmentKeyedLeafConflictAndCanonicalRetryCorpus|TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance|TestP22DualAdapterRecoveryConformance|TestP22RecoveryAdoptedTargetPermitsTerminalTransitionWithFreshSuccessorProof|TestP22RunReplayRetainsHistoricalBudgetSnapshot)$' ./cmd/fid/integration/...
# P1
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestResolveSourceInputProjectFileCases|TestResolveSourceInputExternalReference|TestExpectedRepositoryChecksOnlyHeadAndDirty|TestObserveRecordedRepositoryNotARepositoryReason|TestRepositoryUnavailableCauseFixtures|TestRepositoryCaptureChangeReason|TestRepositoryUnmergedAndSizeClassification|TestResolveReviewCoverageChecksSourcesBeforeDigestExclusivity|TestReasonCodeDetailSchemaIsExactForEveryCode)$' ./internal/project ./internal/review ./internal/cli
# R5
rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid
# R6
rg -n "func (same[A-Za-z]*Repository|validateRepositoryBinding|validEnvelopeRepository|validEnvelopeSource|validEvidenceRepository|validateRunTerminalRepository|verifyDecisionSource|importedGateDestination)" internal -g '!*_test.go'
```

R5's six current decoders are Judgment, Review, Run, Activation, Envelope and Gate. `decodeClosedClaimInput` was introduced at `0b67c8304a84e4e5bcb8df44d79f899e40a16acb` and removed by Task 6 at `1baaecc7983872b06e2a47e6df0801bfb1370ea0`; both diffs/Refs confirm the report. R6 retains Evidence validity/equality, Gate import/equality/registration binding, Run terminal validation/equality, and Envelope repository/Source validity. Judgment's Decision verifier and Review's exact snapshot comparison are gone. Remaining wrappers/comparisons have accepted later owners or binding purposes; neither inventory is an absence assertion in this window.

The other fresh registry counts remain consistent with accepted convergence: R1 45 sites, R2 one shared digest helper, R3 18 files, R4 migration/schema 19, R7 46 writer matches in six files, R9 42 renderers, R10 33 assertions, R11 15 identity matches, R13 15 pending qualifier/sentinel matches, R14 one live inventory test, R19 seven pending grammar/fallback sites and R20 25 pending fallback/specific-error sites. The 20 additional typed flags are adopted `--example` booleans. R16 excludes `export import` to leave 49 keyed corpus leaves. These observations do not certify the later retirement/grammar/Human assertions early.

The exact five-command sequence required by each reviewed task and the QA plan is:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
```

All 15 final report rows are PASS with recorded observed `exit_code: 0`. Format tokens are `exit 0 (no output)`; lint tokens are `0 issues.`. Task 4's final test token names workflow `102.588s`, coverage `62.1%`; its total is `60.9%` and tagged token `276.204s`. Task 5's final test token names workflow `101.002s`, coverage `62.1%`; total `61.0%`, tagged `313.544s`. Task 6 retains its original workflow `(cached)` token, coverage `62.1%`; total `60.7%`, tagged `274.935s`. These are inspected Developer-session receipts, not freshly rerun full-suite claims by this reviewer. No required row is missing or substituted. Full phase gates, benchmarks, archive scale, full-checkptr and Human fixtures were not rerun: their owners/timing are phase end or later windows, and affected verification above was sufficient for this bounded re-review.

VAL-052 report timing evidence:

| Task | cmd/fid | internal/workflow | Slowest other ordinary package reported | Tagged integration | Assessment |
|---|---|---|---|---|---|
| 4 repaired final gate | 374.499s | 102.588s | internal/issues 198.298s | 276.204s | All fresh, below 480s/720s warnings and 600s/900s budgets. |
| 5 repaired final gate | 366.509s | 101.002s | internal/issues 190.218s | 313.544s | All fresh, below warnings/budgets. |
| 6 retained final gate | 283.342s | Cached; no fresh duration | internal/export 116.978s among fresh rows | 274.935s | Historical final-gate receipts retained by the report-only repair. Attempt 4 workflow 114.227s is explicitly historical; no cached duration is invented. |

Task 6's pre-gate cmd/fid `304.689s` and failed ordinary/tagged attempt-4 times `316.372s`/`332.169s` remain recorded and below budgets. All earlier failed attempts retain their failed status even when a later identical command passed.

## Per-Task Findings

### Task 4 - Repository observation and SourceInput; judgment reference leaves

Observation (design sections 1-3, 6, 7; VAL-012/016/022/023): the oracle in `cmd/fid/integration/p23_judgment_repository_test.go:26` independently executes Git HEAD/tree/status, reads config identity and hashes the absolute root. The clean-state requirement explicitly excludes `.fid`; a decoy cached project ID cannot determine the expected value. Seven normal leaves compare complete actual journal Event repositories and recursively every applicable record/Source member, pinning five one-Source shapes and two source-free shapes. I1 reproduces parity and unsafe/replay pairs. Descriptor help and structural evidence at `cmd/fid/input_descriptor_test.go:933`, `:1360` and `:1624` now require exact recursive metadata and independently constructed Path/Value/Allowed expectations, accepted controls and no effects; C1 passes. SourceInput documentation agrees with the observed behavior and retains path containment and typed unavailability.

#### Must-fix

- None. F1-F3 are resolved.

#### Improvement

- None in this task; Runner bookkeeping observations appear below.

### Task 5 - Review leaves and repository expectations

Observation (design sections 1-3, 6, 7; VAL-012/014/016/018/019/022/066): `p23_review_record_test.go:104` enumerates the three Source positions through real descriptor help and checks them against the corpus. Every position has matching, mismatching, omitted-existing and omitted-missing states with independent file hashes, exact details and no Event on refusal. `:167` compares complete Review/three-Source/Event repositories; `:183` independently hashes the known dirty file image. `:207` uses a fresh unreviewed Run and the same restored key/request, so duplicate Review or stale authority cannot substitute for an observation control. I1 reproduces all states, including full original response and journal Event equality.

TRG-24-005 permits the source repair at `internal/cli/input_validate.go:306`: exactly `oid` and `OID` use the existing 40-character lowercase-hex predicate; unknown format labels reject. Synthetic controls reject malformed, empty, short, nonhex, uppercase-hex, alternate-length and placeholder values while accepting a full valid OID. Mixed-case format aliases also reject. C1 reproduces recursive Review/Finding help/structure and built-process malformed-OID failures with exact closed details `{argument:"input.expected_repository.head_oid",value:"not-an-oid",allowed:["oid"]}`, message `invalid input.expected_repository.head_oid`, invalid-input/invalid-argument, false retryability, null subject, exit 2, empty stdout and no effects. Typed `<head-oid>` templates retain their bytes, required-only keys and closed Go decoding; runtime controls substitute a qualifying OID. The descriptor file and previously repaired integration file were not changed by that continuation.

#### Must-fix

- None. F4-F7 and the causal OID dispatch defect are resolved.

#### Improvement

- None in this task; the ledger's inaccurate summary is Runner bookkeeping, not a defect in the committed Task 5 report.

### Task 6 - Gate registration, Envelope and Claim leaves

Observation (design sections 1-3, 6, 7; VAL-022/026/027/028/029/066): F8's commit adds only the exact R5 inventory/accounting and deviation note to the report. All previous gates, failures, provenance and DEC-283/TRG-24-004 evidence are retained. C1 independently reproduces observed Registration identity/digest, binding/Source replay, original/plain-copy/real-import pairs and imported Claim stored-authority lifecycle; I1 reproduces the four-path fixture continuation's successor-proof and historical Budget replay subjects. The shared `SameWorktree` predicate still binds config project and root/import identity and leaves authority in the domain layers. Documentation matches the identity-free inputs and real imported-destination rule. These fixtures also establish fresh feasibility of the unmerged, plain-copy and import-binding setup needed by later owners; no later task implementation is reviewed here.

DEC-283 maintenance changed only its two authorized workflow tools. Its integrated commit and reconciliation pair bind the monotonic Task 6 attempt-3/4/5 chain. TRG-24-004 is closed by the four-path repair and COMPLETE report; TRG-24-005 is closed by the bounded OID repair and COMPLETE report. Current escalations say none open. Fresh sequence validation returns `REREVIEW_PENDING`, awaiting this independent review rather than further repair.

#### Must-fix

- None. F8 is resolved.

#### Improvement

- None in this task.

## Prior Finding Resolution

Every original Must-fix is preserved by identity, owner and subject below and resolved from newly inspected committed evidence plus fresh affected verification. Prior review PASS assertions are not substituted for these results.

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| F1 — missing exact R5/R6 report receipts; no VAL | task-4 | RESOLVED | `task-4-done.md:88` has both exact commands, own-session output tokens, observed exit 0 and full six/nine-match accounting. Current sets match; fresh report-lint COMPLETE, exit 0; report endpoint `6352f94`. |
| F2 — production repository observer used as oracle; VAL-022 | task-4 | RESOLVED | `p23_judgment_repository_test.go:26`, `:56`, `:117`, `:294` independently derive Git/config/root identity and compare complete applicable records/Sources/Events across seven leaves with pinned counts. I1 parity/unavailable/config-decoy tests pass, exit 0. Repair `c0115bd`. |
| F3 — incomplete judgment help/structural expected values; VAL-012/016 | task-4 | RESOLVED | Recursive exact help metadata, corrupted-token guards, independently constructed exact closed error alternatives/values, valid/optional/nullable controls and no effects are present in `input_descriptor_test.go:933`, `:1002`, `:1360`, `:1624`. C1 passes, exit 0. |
| F4 — missing exact R5/R6 report receipts; no VAL | task-5 | RESOLVED | Current `task-5-done.md` Tests run captures both commands/tokens/observed exit 0 and six/nine-match accounting, retaining attempt-1 receipts. Fresh report-lint COMPLETE, exit 0; report endpoint `be56892`. |
| F5 — incomplete three-position Source matrix; VAL-018/019 | task-5 | RESOLVED | `p23_review_record_test.go:104` checks descriptor-derived items/exclusions/evidence-ref positions × four states, independent SHA-256, exact digest/null details, accepted stored digest and no-Event refusal. I1 matrix/completeness tests pass, exit 0. Repair `0424602`. |
| F6 — incomplete independent repository/replay/state/control evidence; VAL-009/022/023/066 | task-5 | RESOLVED | `:167`, `:183`, `:207`, `:298` compare full repositories and original response/Event; normal/non-Git/unmerged/deleted-source and advanced-HEAD states have otherwise viable fresh Run/key requests that succeed unchanged after restoration. I1 passes, exit 0. |
| F7 — incomplete Review/Finding help/structural corpus; VAL-012/016, causal VAL-014 template distinction | task-5 | RESOLVED | C1 recursively checks both leaves and exact deep errors/valid controls/no effects; OID/unknown-format/process tests pass, exit 0. `226cdc9` matches TRG-24-005 scope, preserves lowercase metadata, recognizes only declared OID labels, fails closed, preserves typed placeholders and rejects them at runtime. Fresh report-lint COMPLETE, exit 0. |
| F8 — missing Task 6 R5 decoder receipt; no VAL | task-6 | RESOLVED | `task-6-done.md:57` adds exact command/token/observed exit 0, accounts for six surviving decoders and verified Claim removal provenance. `86f1876` is report-only with no removals. Fresh report-lint COMPLETE; C1/I1 retain independent binding/continuation verification, exit 0. |
| I1 — inaccurate task-3 checkpoint hash; no VAL | Phase Runner bookkeeping; task-3 historical checkpoint | STILL OPEN (Improvement) | `run-ledger.md:384` retains the bad hash. Fresh exact verification returns `fatal: Needed a single revision`, exit 128; actual `545803dc2185f565d2b02258bd18ee0fd1d7d2ed` resolves with exit 0. This is the prior deferrable Improvement, outside product repair scope. |

## Repair Plan

None required. No Must-fix remains; no Developer repair or further Window B Must-fix re-review is requested by this verdict.

## Cross-cutting observations

- **I1 — Improvement; no VAL; carried forward.** `project-docs/docs/phases/phase-24/run-ledger.md:384` contains the invalid task-3 implementation hash described above. Fix direction: Phase Runner appends the correct checkpoint bookkeeping with the verified hash; preserve the historical record. This does not invalidate the clean committed reports or affected Window B results.
- **I2 — Improvement; no VAL; new.** `project-docs/docs/phases/phase-24/run-ledger.md:765` summarizes Task 5's repaired coverage total as `60.7%`; the committed `task-5-done.md:103` exact gate row records `61.0%`. Its command scope, termination and package timings agree, and coverage policy is unchanged. Fix direction: Phase Runner reconciles its summary to the committed report through append-only bookkeeping. The review relies on the committed report and fresh report-lint; it does not promote the ledger's summary into a gate receipt.
- Current brief/progress still describe the earlier DEC-283 recovery checkpoint. They provide planning context; verified current commits, reports and live ledger determine this review's state. Runner-owned ledger/triage dirt was preserved. No product, tests, task specs, existing QA artifacts, decisions, triage or ledger were edited by this reviewer.

## Blocker Candidates

- None. All original Must-fix findings have fresh closure evidence inside tasks 4-6. No protected authority question, work outside the window or exhausted recovery budget remains. No escalation ID was allocated.

## Existing Escalation References

- ESC-24-001 and ESC-24-002: already closed by DEC-283. Current `tl_escalations_current.md` has no open escalation. Maintenance/continuation commits, linked RESOLVED/COMPLETE ledger entries and fresh sequence validation corroborate closure; these historical ESCs cover no outstanding finding.

## Final Lines

Review written: project-docs/docs/phases/phase-24/qa/review-4.md
Must-fix: 0  Improvement: 2  Blocker candidates: 0
VERDICT: PASS
