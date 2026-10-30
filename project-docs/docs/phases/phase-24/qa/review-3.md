# Phase 24 - Mid-Phase Review 3 (Tasks 4-6)

Reviewer: qa-reviewer
Date: 2026-10-06
Tasks covered: phase-24/task-4, phase-24/task-5, phase-24/task-6
Review type: INITIAL
Prior review: project-docs/docs/phases/phase-24/qa/review-2.md

## Verdict

REPAIR. Window B does not advance to Task 7. Eight Must-fix findings have bounded, decision-free ownership in Tasks 4-6. The reproduced tests and report-lint controls pass, but several required evidence rows are absent or assert less than the validation contract requires. No protected decision or exhausted recovery budget is established.

## Scope Reviewed

- Task specs: `project-docs/docs/phases/phase-24/tasks/task-4.md`, `task-5.md`, and `task-6.md`, including the DEC-283 amendment.
- Completion reports: the committed `task-4-done.md`, `task-5-done.md`, and final `task-6-done.md`, including its preserved earlier reports and failed attempts.
- Validation contract: `project-docs/docs/phases/phase-24/qa/validation-contract.md`; QA plan: `project-docs/docs/phases/phase-24/qa/qa-plan.md`, especially Window B and EVID-06..18, EVID-21..25, EVID-33, EVID-37, EVID-38 and EVID-40.
- Context and authority: `AGENTS.md`; the four requested current-context documents; phase decisions through DEC-283; `public-interface.md`; `input-inventory.md`; `triage-log.md` through TRG-24-004; prior review-2; `08_QUALITY_PLAN.md`; `09_AI_WORKFLOW.md` sections 2.7 and 8; and the `qa-evidence` skill. CLI parity was assessed against `05_DESIGN_PLAN.md` sections 1-3, 6 and 7 and the amended CLI/flow/glossary documentation. Held-out material was not opened.
- Implementation commits: task 4 `7c69c2bbb6b8b7f8688bb6d28fc72f8fb5fff61c`; task 5 `ff6c923dcb706513e71ff7dca76723d68da6ba89`; task 6's checkpoint/continuation chain through `1baaecc7983872b06e2a47e6df0801bfb1370ea0`, `b7d1fd56b801d10378b35dff8f45684d1990804d` and `01e518f460d7339848c6549ecb7dbd4d58a2f6f4`. Workflow verifier maintenance: `13164582232b9b43eddac93a0ef75b6f4465783c`.
- Committed report endpoints: task 4 `7ac55a292ced5e0eabd6933394864da2fae01340`; task 5 `010a2550ed682d6445dc2c0f79bc12814b206860`; task 6/review base `7f64268a9081d8ea97532b917ddf36d4b8e8093b`. Commit diffs, subject/Refs footers, production observation/resolution/binding code, descriptor and domain tests, command tests, integration journeys and documentation were inspected. Runner-owned dirty ledger and triage files were read as current evidence and left untouched.

## Summary

- Review outcome: REPAIR
- Must-fix findings: 8
- Improvement findings: 1
- Blocker candidates: 0
- Existing escalation references: ESC-24-001 and ESC-24-002, both closed by DEC-283; none open or covering a new finding
- Covered validation assertions: VAL-006, VAL-007, VAL-008, VAL-009, VAL-011, VAL-012, VAL-013, VAL-014, VAL-015, VAL-016, VAL-017, VAL-018, VAL-019, VAL-020, VAL-021, VAL-022, VAL-023, VAL-024, VAL-025, VAL-026, VAL-027, VAL-028, VAL-029, VAL-035, VAL-052, VAL-060, VAL-066

## Validation Contract Coverage

Statuses apply to the reviewed leaves. Whole-tree completion, later Run/Gate/Evidence producers, Task 10's domain messages and phase-end checks retain their explicitly scheduled owners; their later timing does not defer missing evidence for Tasks 4-6. Passing reruns below prove the assertions actually present in the tests, not missing corpus rows.

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-006 | task-2, converged by tasks 4-6 | ON TRACK | Named judgment keyed conflict/canonical retry corpus and shared leaf conformance reran in I1; reviewed convergence retains differing-field and no-Event subjects. Later leaf convergence remains due. |
| VAL-007 | task-2, converged by tasks 4-6 | ON TRACK | Same-key canonical retries and fake/fid conformance remain asserted in I1; binding replay fixtures retain result/Event identity in C1. This is the Window B slice, not a new whole-tree certification. |
| VAL-008 | task-4 | ON TRACK | `TestP24JudgmentLeavesReplayAfterHeadAndFileChange` compares supplied versus omitted Source digest and names `source`/`references`; unchanged input replays. I1 passed. Review optional-array presence remains preceding-window context. |
| VAL-009 | tasks 4-6 | NOT ON TRACK | Task 4 compares original results after HEAD/file change; Task 6 compares binding and Source replay results. Task 5 only checks `replayed` and Event count after HEAD advancement, without complete original-result equality: F6. |
| VAL-011 | tasks 4-6 | ON TRACK | Live-tree/type reflection and seeded mismatch rejection reran in C1; seven judgment, two Review and eight binding leaves leave the descriptor fence. Completion of all 39 leaves remains Task 9. |
| VAL-012 | tasks 4-6 | NOT ON TRACK | Task 6's recursive help checks types, requirements, meanings and constraints. Task 4 checks field paths/meanings but omits type/requirement/constraint assertions; Task 5 checks only template line and Git commands: F3/F7. |
| VAL-013 | tasks 4-6 | ON TRACK | C1 exercises standalone examples for all adopted leaves; judgment/binding tests compare emitted templates. Fresh scratch-binary Review/Finding examples were independently inspected for required-only keys, documented order and two-space indentation; both terminated exit 0. |
| VAL-014 | tasks 4-6 | ON TRACK | Closed Go decoding, descriptor validation/type reflection and emitted-template checks in C1 pass. |
| VAL-015 | tasks 4-6 | ON TRACK | Built-process exclusivity test checks the six prohibited combinations and standalone successes for all 17 Window B leaves with exact details and no workspace effects; C1 passed. |
| VAL-016 | tasks 4-6 | NOT ON TRACK | Task 6 derives recursive structural failures with exact details. Task 4 checks only nonempty `allowed`; Task 5 lacks applicable nested missing/null/type/format rows and top-level verdict enum row: F3/F7. |
| VAL-017 | tasks 4-6 | ON TRACK | Judgment Source repositories, Review top-level/all three Source repositories and binding identity members are explicitly rejected at named paths with no Event; C1/I1 passed. |
| VAL-018 | tasks 4-6 | NOT ON TRACK | Task 4's five positions and Task 6's three positions have matching/mismatch rows. Task 5 lacks accepted supplied-digest rows for all three descriptor-enumerated positions: F5. Shared resolver success does not replace the missing leaf/position rows. |
| VAL-019 | tasks 4-6 | NOT ON TRACK | Task 4 and Task 6 have existing/missing omitted-digest controls. Task 5 lacks omitted-existing exclusion and omitted-missing exclusion/evidence rows: F5. |
| VAL-020 | task-4 | ON TRACK | P1/I1 cover absolute/parent paths, escaping symlink, directory and accepted in-root symlink; typed paths, null observed digest and no-Event checks are present. |
| VAL-021 | task-4 | ON TRACK | P1/I1 cover omitted external digest refusal and supplied digest stored verbatim without a local file. |
| VAL-022 | tasks 4-6 | NOT ON TRACK | Task 4 uses production `ObserveRecordedRepository` as its oracle: F2. Task 5 checks only independent HEAD/tree/dirty, without independent base/root/config identity or Event repository: F6. Task 6 independently observes Git/config/root identity and checks the applicable stored/Event repositories. |
| VAL-023 | tasks 4-6 | NOT ON TRACK | Task 4's seven normal/non-Git/unmerged pairs and five classifier causes pass P1/I1; Task 6's four observing-leaf pairs pass C1. Task 5 has no unmerged Review row: F6. |
| VAL-024 | task-5 | ON TRACK | `TestExpectedRepositoryChecksOnlyHeadAndDirty` and real-process Review match/HEAD-differs/dirty-differs rows pass P1/I1; refusals append no Event. Tasks 7/8 own later expectation producers. |
| VAL-025 | task-5 | ON TRACK | Real-process mismatch checks echo the input pair and full observed repository; D1 pins exact registered detail schema. Later producers remain Tasks 7/8. |
| VAL-026 | task-6 | ON TRACK | `TestGateRegisterObservesRepositoryWithoutCallerIdentity` independently computes config project ID, root hash, HEAD/tree/dirty and binding digest, then compares Registration/Event; C1 passed. |
| VAL-027 | task-6 | ON TRACK | Original-root acceptance, plain-copy refusal/no Event and genuine imported-destination acceptance for Envelope creation pass C1. Task 10 owns the later domain-cause message. |
| VAL-028 | task-6 | ON TRACK | Claim acquisition original-root acceptance versus plain-copy rejection, observed worktree ID and unchanged persistence pass C1. Task 10 owns the later domain-cause message. |
| VAL-029 | task-6 slice | ON TRACK | Genuine import binding/original/copy controls for Envelope/Claim and stored Claim renew/release/expire pass C1; wrong binding/root and copied destination reject. Later Run/Gate-run rows remain Tasks 7/8, messages Task 10. |
| VAL-035 | tasks 4-6 | ON TRACK | I1 reruns shared leaf/amended-leaf and dual-adapter recovery conformance; accepted new inputs and removed-member refusals retain fake/fid parity. Later assembly remains Tasks 7-9. |
| VAL-052 | tasks 4-6 | ON TRACK | All reported fresh ordinary/tagged package timings remain below warning/limit; cached workflow rows are explicitly distinguished from historical uncached measurements. Details below. No Window B benchmark is required. |
| VAL-060 | tasks 4-6 | ON TRACK | Built-binary standalone examples run in empty non-project directories with exit 0, empty stderr and unchanged directory state; C1 passed. |
| VAL-066 | tasks 4-6 | NOT ON TRACK | Task 4 and Task 6 have unmerged/deleted replay and otherwise viable new-key/refusal/restored-success pairs. Task 5 lacks unmerged replay and reuses an already-reviewed Run for new-key controls, so observable-state viability is unproved: F6. |

## Gate And Evidence Checks

Working directory for independent checks: `/home/mark/projects/fid`, except the two direct examples, run from `/tmp` using the freshly built inventory binary `/tmp/tmp.sLCU14p3d8/fid`. Every result below reached an observed command termination; long nested sessions were polled until an actual `exit_code` arrived. Focused reruns supplement, rather than replace, the Developer's full required commands. No full phase-level gate was run.

The five exact task commands are the same for Tasks 4-6, in this order:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
```

All three committed reports record these byte-identical commands in order, each PASS with its own output token and observed `exit_code: 0`. Task 4 records `exit 0 (no output)`, `0 issues.`, an `ok ... internal/workflow (cached)` coverage row, `total: ... 60.6%`, and `ok ... cmd/fid/integration 284.152s`; Tasks 5/6 record coverage totals of 60.7% and integration tokens of 252.256s/274.935s respectively. Task 6's final workflow row explicitly says `(cached)` and coverage 62.1%. Earlier failures remain immutable failed evidence. The mechanical five-command comparison passes; it does not cover the extra accepted inventory obligation: the brief at lines 49-50 and inventory at lines 6-9 require each task to rerun and record its owned registries. R5 is assigned to Tasks 3-9 and R6 to Tasks 4-8. Their omissions produce F1/F4/F8 even though report-lint returns COMPLETE.

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Required command comparison | Exact Commands To Run versus Tests run/Exact quality-gate results, tasks 4-6 | PASS | Five commands each; order, scope, tokens and recorded observed exit 0 present. No placeholder or narrower substitute. Additional owned registry omissions are separate failures below. |
| Representative exact gate | `gofmt -l cmd internal` | PASS | `exit 0 (no output)`; observed termination `exit_code: 0`. |
| P1: Source/observation/expectation tests | Named race command below, `./internal/project` | PASS | `ok  \tgithub.com/mbell-nzl/fid/internal/project\t1.074s`; observed termination `exit_code: 0`. |
| C1: descriptor/binding/lifecycle/converged ordinary fixtures | Named race command below, `./cmd/fid` | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t43.925s`; observed termination `exit_code: 0`. |
| I1: judgment/Review/fake-fid/final tagged fixtures | Named tagged race command below, `./cmd/fid/integration/...` | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t20.326s`; observed termination `exit_code: 0`. |
| D1: coverage resolution/closed details | Named race command below, `./internal/review ./internal/cli` | PASS | `ok ... internal/review 1.016s`; final `ok  \tgithub.com/mbell-nzl/fid/internal/cli\t1.125s`; observed termination `exit_code: 0`. |
| Task 4 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-4-done.md project-docs/docs/phases/phase-24/tasks/task-4.md` | PASS | Fresh `OUTCOME COMPLETE`; `PASS report-lint completed with no failures`; observed termination `exit_code: 0`. |
| Task 5 report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-5-done.md project-docs/docs/phases/phase-24/tasks/task-5.md` | PASS | Fresh `OUTCOME COMPLETE`; `PASS report-lint completed with no failures`; observed termination `exit_code: 0`. |
| Task 6 final report control | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-6-done.md project-docs/docs/phases/phase-24/tasks/task-6.md` | PASS | Fresh `OUTCOME COMPLETE`; `PASS report-lint completed with no failures`; observed termination `exit_code: 0`. |
| Report/implementation provenance | `git rev-parse --verify '<hash>^{commit}'`; `git show` subjects, Refs and diffs | PASS | All 47 full commit hashes extracted from the reports resolve with child exit 0; primary subjects match. Report endpoints are committed, with no working/index report diff or hidden index flag. Introducing diffs/Refs checked for all 21 retained Task 6 rows and its four final rows. |
| Task 6 failure preservation | Final report quotation versus committed report at `45961a91b60bbeae7e1747642637544bbcf89fbe` | PASS | Removing one added quotation prefix reconstructs all 507 prior report lines byte-for-byte. Earlier failed ordinary/tagged and imported-lifecycle attempts retain tokens and terminations; later passing evidence does not overwrite them. |
| Task 6 bounded continuation | `git show 01e518f460d7339848c6549ecb7dbd4d58a2f6f4`; TRG-24-004/DEC-283/ledger | PASS | Exactly four allowed fixture definitions plus the report. Root-local history replaces positive copied history; expired-prior/fresh controls, successor proof/Receipt and all four replay/Budget/Event equality subjects remain. C1/I1 reproduce affected verification. |
| DEC-283 recovery-chain evidence | `sh project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-24/run-ledger.md` | PASS | `VALID state=ADVANCEABLE`; observed termination `exit_code: 0`. Maintenance commit changes only the authorized verifier and drill script. Current ledger records Task 6 attempt 5 RESOLVED and COMPLETE. This check is evidence about the chain, not authority to advance past this review. |
| R1-R20 and recursive live inventory | Accepted `input-inventory.md` section 1 block and section 2 regeneration commands | PASS | Fresh scratch build/walk: 92 leaves, 29 groups, 39 input leaves, 50 keyed leaves (49 excluding import), no missing Short, 42 typed flags. Walk ends `252 total`, observed exit 0. Individual searches observed exit 0 with accounted matches. |
| Owned registry report coverage | Task 4/5 R5/R6; Task 6 R5 | FAIL | No captured command/token/termination in the respective reports: F1/F4/F8. QA's inventory rerun cannot be transcribed as Developer-session evidence. |
| Direct examples | `/tmp/tmp.sLCU14p3d8/fid review record --example`; `/tmp/tmp.sLCU14p3d8/fid finding address --example` | PASS | Each own final token `}`; observed termination `exit_code: 0`. Inspected complete objects: optional digest/inapplicable findings absent, required nullable Source members present, documented order/two-space indentation. |

Exact independent named-test commands:

```bash
# P1
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestResolveSourceInputProjectFileCases|TestResolveSourceInputExternalReference|TestExpectedRepositoryChecksOnlyHeadAndDirty|TestObserveRecordedRepositoryNotARepositoryReason|TestRepositoryUnavailableCauseFixtures|TestRepositoryCaptureChangeReason|TestRepositoryUnmergedAndSizeClassification)$' ./internal/project
# C1
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestAdoptedExampleProcessExclusivity|TestP24JudgmentDescriptorSourcePositionsAndExamples|TestP24JudgmentStructuralValidationCorpus|TestP24ReviewAndFindingDescriptorHelpAndStructure|TestP24GateEnvelopeClaimDescriptors|TestP24RemovedBindingMembers|TestP24BindingDescriptorProcessCorpus|TestGateRegisterObservesRepositoryWithoutCallerIdentity|TestP24EnvelopeReplayAfterSourceChangeOrDeletion|TestP24BindingReplayBeforeObservation|TestP24EnvelopeAndClaimOriginalCopyAndImport|TestP24ImportedClaimStoredAuthority|TestP24EnvelopeSourcePositions|TestExecutionBound.*|TestExecutionTransitionFixtureDoesNotInheritExpiredClaim)$' ./cmd/fid
# I1
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24HumanOutcomeRejectsReferenceDigestMismatch|TestP24JudgmentSourceDigestsAndRepositoryParity|TestP24JudgmentSourceInputValidation|TestP24JudgmentLeavesRepositoryUnavailable|TestP24JudgmentLeavesReplayAfterHeadAndFileChange|TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation|TestP23ReviewRecordRealProcessJourney|TestP24JudgmentKeyedLeafConflictAndCanonicalRetryCorpus|TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance|TestP22DualAdapterRecoveryConformance|TestP22RecoveryAdoptedTargetPermitsTerminalTransitionWithFreshSuccessorProof|TestP22RunReplayRetainsHistoricalBudgetSnapshot)$' ./cmd/fid/integration/...
# D1
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestResolveReviewCoverageChecksSourcesBeforeDigestExclusivity|TestReasonCodeDetailSchemaIsExactForEveryCode)$' ./internal/review ./internal/cli
```

Inventory deltas are explained by accepted work: R1 now has 45 call sites; R2 has the shared digest helper; R3 has 18 owning files and migration 0019; R5 has six remaining closed decoders after Claim convergence; R6 has nine surviving comparison/validation functions after judgment/Review convergence, with the new shared `SameWorktree` inspected directly; R7's remaining caller writers include later Run/Gate/Evidence owners. R9 has 42 renderers, R10 33 assertion matches, R11 15 identity matches, R13 15 pending sentinel matches, R19 seven pending grammar sites and R20 25 pending fallback sites. No-match Short inventory and the extra 20 boolean example flags are expected descriptor adoption effects. Later-owner matches are not retired-surface absence failures in Window B.

VAL-052: task 4's final cmd/fid is 223.171s and tagged integration 284.152s; its uncached historical workflow/issues measurements are 122.340s/181.213s. Task 5 records cmd/fid 267.739s, workflow 115.695s, issues 158.407s and tagged 252.256s. Task 6's final cmd/fid is 283.342s, export 116.978s and tagged 274.935s; pre-gate cmd/fid is 304.689s. Its final workflow is cached, with attempt 4's 114.227s retained only as history. Even Task 6's failed 316.372s ordinary and 332.169s tagged attempts remain below the 480s/720s warning levels and 600s/900s limits. No budget breach or invented cached duration was found; the review's focused times do not substitute for full-suite timing.

## Per-Task Findings

### Task 4 - Repository observation and SourceInput; judgment reference leaves

Observation: Source resolution belongs in `internal/project`, and judgment domain mutations resolve/store it after authority and replay checks. CLI descriptors and updated SourceInput examples follow the accepted authoring flow. Path containment, typed unavailable reasons and no-Event refusals have positive/negative fixtures. The remaining weaknesses are in the proof of exact contracts below.

#### Must-fix

- **F1 — Must-fix; no VAL.** [`task-4-done.md:58`](../tasks/task-4-done.md#tests-run) records no R5 decoder or R6 repository-comparison inventory command. The accepted inventory header requires each owning task to rerun and record its commands; R5 owns Tasks 3-9, R6 Tasks 4-8. Owning task: task-4. Repair boundary: capture the exact R5/R6 commands and current match accounting in the task-4 report, retaining all prior gate evidence. Required evidence: each command's own token and observed exit status from the recording session, committed report and fresh COMPLETE report-lint. Re-review scope: both registry rows and report integrity; QA's own reruns are not a substitute.
- **F2 — Must-fix; VAL-022.** [`cmd/fid/integration/p23_judgment_repository_test.go:83`](../../../../../cmd/fid/integration/p23_judgment_repository_test.go) and line 266 call the production `project.ObserveRecordedRepository` to construct expected repositories, then compare the production mutation's Source/Event against that same implementation. This is not the contract's independent Git/root/config observation and can pass a shared wrong observation. Owning task: task-4. Repair boundary: replace only these expected-value fixtures in this integration file with independently obtained Git HEAD/tree/dirty, absolute-root SHA-256 and config project ID; compare every applicable record/Source/Event member for all seven leaves. Required evidence: the named repository-parity/unavailable tests, exact task gates and updated report prove base=HEAD and complete identity/capture fields without a production observer as oracle. Re-review scope: independence of the oracle and all seven normal-worktree comparisons, plus affected integration rerun.
- **F3 — Must-fix; VAL-012, VAL-016.** [`cmd/fid/input_descriptor_test.go:754`](../../../../../cmd/fid/input_descriptor_test.go) checks judgment help paths/meanings but not each field's type, requirement or declared constraints; line 846's structural error checks accept any nonempty `allowed` array instead of the exact accepted keys/type/members. Thus wrong help metadata or wrong nonempty alternatives still pass the claimed conformance. Owning task: task-4. Repair boundary: strengthen only the seven judgment leaves' help/structural evidence in this test file and task-4 report. Required evidence: recursive descriptor-derived help metadata and exact closed CLI error details for applicable failure classes/depths, with valid controls and no effects, followed by exact gates and COMPLETE report-lint. Re-review scope: judgment descriptor conformance and the new expected values; rerun affected C1 tests. Shared generator correctness does not replace per-leaf error mapping checks.

#### Improvement

- None in this task. The prior ledger bookkeeping Improvement is preserved below.

### Task 5 - Review leaves and repository expectations

Observation: the shared expectation check compares only HEAD/dirty in `internal/project`; Review domain code resolves Sources and applies existing coverage rules. The CLI mismatch shape, Git guidance and updated documentation agree with the accepted pair. The passing test fixture does not yet establish all required Review states.

#### Must-fix

- **F4 — Must-fix; no VAL.** [`task-5-done.md:55`](../tasks/task-5-done.md#tests-run) likewise has no owned R5/R6 inventory command or captured result. Owning task: task-5. Repair boundary: task-5 report only for these registry rows and match accounting; retain earlier failures and exact gates. Required evidence: exact command, own output token, observed termination, committed report and fresh COMPLETE report-lint. Re-review scope: both rows and committed evidence integrity.
- **F5 — Must-fix; VAL-018, VAL-019.** [`cmd/fid/integration/p23_review_record_test.go:30`](../../../../../cmd/fid/integration/p23_review_record_test.go) constructs three manually chosen Source positions. Lines 74-101 mismatch each position, but the accepted Review at line 104 supplies only the exclusion digest and omits item/evidence digests. No accepted supplied-digest evidence row, omitted-existing exclusion row, or omitted-missing exclusion/evidence row completes the required position-by-state matrix. The helper unit test's matching item does not prove stored CLI acceptance at every position. Owning task: task-5. Repair boundary: this Review integration file and task-5 report; derive/check the complete three-position set against the descriptor, and use otherwise valid fresh Review Runs for each state. Required evidence: every position × matching, mismatching, omitted-existing, omitted-missing, with independently hashed stored digests, exact mismatch/null details and no Event on refusal; exact gates and COMPLETE report-lint. Re-review scope: enumeration, all matrix rows and affected real-process rerun.
- **F6 — Must-fix; VAL-009, VAL-022, VAL-023, VAL-066.** [`cmd/fid/integration/p23_review_record_test.go:118`](../../../../../cmd/fid/integration/p23_review_record_test.go) checks only independent HEAD/tree/dirty and Source=Review; it never compares base/config/root identity and the Event repository to a complete independent observation. Lines 127-151 assert only `replayed:true`/Event count, without comparing the complete original result/Event. The unavailable fixture only removes `.git`, omitting required unmerged paths. Its new-key controls reuse a Run already carrying a Review; the non-Git control also retains an outdated HEAD expectation, and no restored observable-success control proves these requests otherwise succeed. Owning task: task-5. Repair boundary: this integration file and task-5 report, preserving production behavior and expectation contract. Required evidence: complete independent record/Source/Event comparison; original-result equality allowing only replay metadata; normal/non-Git/unmerged and deleted-source pairs; fresh unreviewed Run/new-key controls with matching expectations that succeed in the restored observable state and fail with the exact required code/details/no Event in the target state; exact gates and COMPLETE report-lint. Re-review scope: complete observation, replay equality, both unavailable states and control viability in a fresh affected integration rerun.
- **F7 — Must-fix; VAL-012, VAL-016.** [`cmd/fid/input_descriptor_test.go:610`](../../../../../cmd/fid/input_descriptor_test.go) checks Review/Finding help only for the template line and Review Git commands. The structural cases at lines 619-642 cover the first top-level field and selected nested unknown/enum rows; nested missing/non-nullable null/wrong type/declared format rows and the applicable top-level verdict enum are absent. This is narrower than EVID-11/EVID-14 and the report's claimed complete field/class conformance. Owning task: task-5. Repair boundary: only the two Review/Finding leaves' evidence in this test file and task-5 report. Required evidence: recursive help paths/types/requirements/meanings/constraints and all applicable structural classes at top level and two nested levels where available, exact path/value/allowed/error metadata, valid controls and no workspace/Event effects; exact gates and COMPLETE report-lint. Re-review scope: both adopted leaves and affected descriptor process rerun.

#### Improvement

- None in this task.

### Task 6 - Gate registration, Envelope and Claim leaves

Observation: Registration/Envelope/Claim binding uses fid's observation and the shared project/worktree/import predicate. The tests preserve original-root acceptance, plain-copy rejection and real imported-destination acceptance. Renew/release/expire use the same stored-authority rule. CLI/flow documentation matches the adopted inputs; later Run/Gate-run/Evidence expectation interfaces and Task 10 messages remain assigned to their planned owners.

DEC-283 verifier maintenance and TRG-24-004 were checked only as QA evidence. The final four-file convergence matches the allowed paths and preserves predicate-matrix, expired-Claim isolation, successor-proof/Receipt and historical replay/Budget subjects. The imported Claim lifecycle report preserves the same named command's initial fixture failure, true renew/release/expire product red and later green; `TestP24ImportedClaimStoredAuthority` independently passes in C1. All five final gates, prior introducing rows and fresh report-lint COMPLETE are present. TRG-24-004's fixture repair is closed; F8 is a distinct missing inventory receipt, not a reopening of its four-file runtime repair.

#### Must-fix

- **F8 — Must-fix; no VAL.** [`task-6-done.md:56`](../tasks/task-6-done.md#tests-run) records R6/R7 but no R5 decoder inventory command, although R5 assigns ownership to Tasks 3-9 and this task removes `decodeClosedClaimInput`. Neither the final nor preserved report contains that required command/token/termination. Owning task: task-6. Repair boundary: report-only capture of exact R5 and accounting for the removed decoder/current six remaining decoders; retain all prior failures, provenance and passing gates. Required evidence: fresh recording-session output and observed exit status, committed report and fresh report-lint COMPLETE. Re-review scope: R5 receipt/report integrity, with representative binding verification retained; no product, fixture, control or gate-definition change is authorized by this finding.

#### Improvement

- None in this task.

## Prior Finding Resolution

This is an INITIAL Window B review. Review-2's PASS and its Window A F1-F5 resolutions are preserved as preceding-window context, not reused as current deterministic evidence or re-certified outside this scope. The only carried Improvement was independently checked again.

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| Review-2 records review-1 F1/F2 storage/archive evidence resolved | task-1 | RESOLVED in prior window; outside this initial review | Current prior artifact and retained commits/reports were inspected for context. No new Window A verdict is issued. |
| Review-2 records review-1 F3 unavailable/deleted judgment replay evidence resolved | task-2 | RESOLVED in prior window; outside this initial review | Current judgment unavailable/deleted tests reran in I1 for Task 4's adopted behavior; prior repair attribution remains unchanged. |
| Review-2 records review-1 F4 exclusivity and F5 structural evidence resolved | task-3 | RESOLVED in prior window; outside this initial review | Live descriptor framework and built-process exclusivity reran in C1 for adopted Window B leaves. F3/F7 here concern later leaf-specific evidence, not a reversal of Window A. |
| **I1 — Improvement; no VAL:** incorrect task-3 implementation hash in ledger checkpoint | task-3 checkpoint; Phase Runner bookkeeping | STILL OPEN | `run-ledger.md:384` still says `545803d9f55b7a1e957a4d27b0a9a7703b6e5c2b`; fresh `git rev-parse --verify '<bad-hash>^{commit}'` emits `fatal: Needed a single revision`, observed exit 128. Actual `545803dc2185f565d2b02258bd18ee0fd1d7d2ed` resolves with observed exit 0. Fix direction: Phase Runner corrects its checkpoint as bookkeeping; the task-3 committed report has the valid hash. No ledger edit is made by this review, and this deferrable carry-forward does not cause REPAIR. |

## Repair Plan

Repair order: task-4, then task-5, then task-6. Repairs remain within these existing tasks; accepted behavior, schemas, dependency choices, authority boundaries, commands, pass semantics, coverage policy and budgets must remain unchanged. Test-evidence changes for Tasks 4/5 require their five exact task gates above and committed updated reports; report-only inventory capture for Task 6 retains already recorded immutable gate evidence and requires fresh report-lint. Each Developer captures its own evidence; none may attest from this review's output. A fresh independent Window B reviewer must check every F1-F8 row, rerun representative affected commands and write the next unused artifact before Task 7 dispatch.

Exact omitted registry commands, from the accepted inventory:

```bash
# R5, all three owners
rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid
# R6, tasks 4 and 5
rg -n "func (same[A-Za-z]*Repository|validateRepositoryBinding|validEnvelopeRepository|validEnvelopeSource|validEvidenceRepository|validateRunTerminalRepository|verifyDecisionSource|importedGateDestination)" internal -g '!*_test.go'
```

Both currently return matches and exit 0; record their real observed status and explain remaining later-owner matches. No absence gate or undeclared non-zero PASS is introduced.

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| F1 | task-4 | no VAL | `task-4-done.md`: exact R5/R6 receipts and current match accounting only | Fresh own-session command/token/termination; committed report; COMPLETE report-lint; retain all failures/gates | Verify both command receipts and report provenance |
| F2 | task-4 | VAL-022 | `cmd/fid/integration/p23_judgment_repository_test.go` expected-value fixtures and task-4 report | Independent Git/config/root oracle, complete applicable record/Source/Event comparison for seven leaves; named tests; all five exact gates; COMPLETE report-lint | Inspect oracle independence and seven comparisons; rerun parity/unavailable integration tests |
| F3 | task-4 | VAL-012, VAL-016 | Seven judgment leaves' help/structural tests in `cmd/fid/input_descriptor_test.go` and task-4 report | Recursive exact help metadata and exact error alternatives/closed details; qualifying controls; affected tests; five gates; COMPLETE report-lint | Inspect expected values and rerun judgment help/structural corpus |
| F4 | task-5 | no VAL | `task-5-done.md`: exact R5/R6 receipts and match accounting only | Fresh own-session command/token/termination; committed report; COMPLETE report-lint; retained history | Verify both command receipts and report provenance |
| F5 | task-5 | VAL-018, VAL-019 | Review Source matrix in `cmd/fid/integration/p23_review_record_test.go` and task-5 report | Descriptor-complete three positions × four required states; actual SHA-256, exact errors/nulls/no Event; five gates; COMPLETE report-lint | Inspect enumeration/all rows; rerun affected Review integration test |
| F6 | task-5 | VAL-009, VAL-022, VAL-023, VAL-066 | Observation/replay/unavailable/control fixtures in the same Review integration file and task-5 report | Independent complete repositories including Event; original-result/Event equality; normal/non-Git/unmerged/deleted states; viable fresh controls and restored success; five gates; COMPLETE report-lint | Independently inspect all state pairs/control viability; rerun affected Review integration tests |
| F7 | task-5 | VAL-012, VAL-016 | Two Review/Finding leaves' conformance tests in `cmd/fid/input_descriptor_test.go` and task-5 report | Complete recursive help metadata and applicable top/deep structural classes with exact closed details/no effects and valid controls; five gates; COMPLETE report-lint | Inspect both leaves and rerun affected descriptor process test |
| F8 | task-6 | no VAL | `task-6-done.md` only: R5 receipt/accounting; preserve all final/historical evidence | Fresh exact R5 token/observed termination; committed report; fresh COMPLETE report-lint | Verify decoder receipt and report integrity; retain affected binding/lifecycle representative rerun |

## Blocker Candidates

- None. Every Must-fix has existing task ownership and requires evidence completion under accepted contracts. Failed historical gates and the earlier Task 6 RECOVERY_REQUIRED/BLOCKED labels do not establish an open protected blocker. No new escalation ID is allocated.

## Existing Escalation References

- ESC-24-001 and ESC-24-002: closed by DEC-283, independently corroborated by the maintenance/continuation commits, current escalation register and valid recovery ledger. They are historical recovery context and cover none of F1-F8. No open ESC exists.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-24/qa/review-3.md
Must-fix: 8  Improvement: 1  Blocker candidates: 0
VERDICT: REPAIR
```
