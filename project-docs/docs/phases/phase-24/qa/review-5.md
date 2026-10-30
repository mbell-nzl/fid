# Phase 24 - Mid-Phase Review 5 (tasks 7-8)

Reviewer: qa-reviewer
Date: 2026-10-07
Tasks covered: phase-24/task-7, phase-24/task-8; task-9 current-HEAD impact and tasks 9a-14 feasibility only
Review type: INITIAL
Prior review: none for Window C; review-4 is preceding Window B history

## Summary

- Review outcome: REPAIR
- Must-fix findings: 2
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-24-001, ESC-24-002, ESC-24-003 (resolved history; none open)
- Covered validation assertions: VAL-009, VAL-011-019, VAL-022-025, VAL-029-035, VAL-052, VAL-060-062, VAL-066. Later-window assertions receive feasibility assessment, not acceptance.

The focused ordinary and real-process tests pass at both recorded task close commits and at the supplied phase HEAD. The reports contain the exact required final gates and observed termination receipts, and independently pass report-lint. Two test-evidence gaps remain: the dirty Run and Gate repository expectations are partly produced by the same production observer being verified. Independent Git checks omit `dirty_digest`; the terminal Run matrix also compares its Event with the returned record rather than a complete independent expected repository. These gaps leave the complete independent-observation assertions unproved despite green tests. Both repairs have decision-free ownership inside the reviewed tasks.

Window C does not advance. Repair task 7, then task 8, and obtain a fresh independent Window C review in the next unused artifact. D1 after task 9a, D2 after task 11, and E after task 14 remain mandatory. No Human judgment is recorded here.

## Scope Reviewed

Read AGENTS.md; workflow sections 2.4, 2.7 and 8; the four current context documents; quality-plan sections governing evidence and phase checks; the phase contract, QA plan, executable inventory and public-interface excerpts; task specs/reports 7-9; task specs 9a-14; phase decisions including DEC-284; triage and escalation status; the review template; and preceding review-3/review-4 evidence. CLI observations use design-plan sections 2-3 and 5-7. No full PRD/TRD/roadmap or held-out material was loaded.

| Scope | Verified commit | Inspection |
|---|---|---|
| Task 7 close | `9b585b73bf78188fb34cca4938734ded597afe55` | Isolated committed spec/report/source; implementation `8d1c83855c0f681515e44e2c3a61bbe543ef515b`; ordinary and tagged reruns |
| Task 8 close | `697e1562acf1103f7239e7804b291a3deaf27a97` | Isolated committed spec/report/source; final implementation `9f0fc3577128857c8f5523870b965daded1b3f1c`; ordinary and tagged reruns |
| Supplied current HEAD | `e39bfcc11237e108cebb58d46aa5dc696ab44ca1` | Task 9 report and shared-descriptor impact; source checkpoints `9d654e976792836cbdf7de8c114d0831b73f721e` and `69589747ca77a1a3dae7e4d6083229d0193ed52f`; current ordinary/tagged reruns and feasibility probes |

Each close hash and the cited report/implementation hashes was checked using `git rev-parse --verify '<hash>^{commit}'`, with the requested full hash returned and observed `exit_code: 0`. Commit subjects, task references, diffs, retained-test convergence and report deviations were inspected. Task 9's report-cited `d6e03c9f9ec7a445f3af50451cb9b14d4d52dd5c` also resolves; the ending supplied commit finalizes that report.

The main worktree began on `phase-24` at the supplied HEAD with only `project-docs/docs/phases/phase-24/run-ledger.md` dirty. Detached inspection worktrees were created at the two close commits. Their Go VCS-stamp builds encountered the unrelated `/tmp/.git` directory because their own `.git` is a worktree file. Local shared clones with real `.git` directories, checked out at the identical commits, provided clean executable inspection roots:

- C7: `/tmp/fid-p24-review5-clone7` at task 7 close.
- C8: `/tmp/fid-p24-review5-clone8` at task 8 close.
- H: `/home/mark/projects/fid` at supplied HEAD, preserving runner ledger dirt.

No product, task, decision, workflow, dependency, generated or runner-owned artifact was edited by this reviewer.

## Validation Contract Coverage

ON TRACK applies only to the owned Window C slice. The separate task-9 whole-tree review remains D1. Earlier task-4 path/classifier and task-6 binding assertions remain preceding-window context; Window C does not re-certify all their leaves.

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-009 | 7, 8 | ON TRACK | C7I/C8I/H2 compare committed replay responses after HEAD changes and retain journal/projection state; unavailable/deleted-source controls exercise the replay ordering. |
| VAL-011 | 7, 8 | ON TRACK | C7U/C8U enumerate live adopted descriptors; seven Run and four Gate/Evidence/Artifact leaves leave the pending fence. H1 repeats live registry equality after task 9. |
| VAL-012 | 7, 8 | ON TRACK | Descriptor process corpora inspect typed fields/help, required/nullable metadata, repository expectations and Source positions. |
| VAL-013 | 7, 8 | ON TRACK | Typed examples and required-only templates exercised at each close; process corpora and current whole-tree corpus pass. |
| VAL-014 | 7, 8 | ON TRACK | Example conflicts exercised in both modes, including input/actor/key/local-flag combinations and no effects. |
| VAL-015 | 7, 8 | ON TRACK | Shared-reader stdin/file paths and descriptor adoption survive close-commit and current process selections. |
| VAL-016 | 7, 8 | ON TRACK | Structural field-path/type/key/depth alternatives and valid controls retained in adopted process corpora; H1 exercises shared validator changes. |
| VAL-017 | 7, 8 | ON TRACK | C7R/H1 reject removed Run repository members; C8U/H1 reject removed Gate/Evidence members with exact paths and no durable effects. |
| VAL-018 | 7 | ON TRACK | C7I budget Source matrix: supplied matching digest accepted and stored; mismatching digest rejected with exact independent SHA-256 details and no Event. |
| VAL-019 | 7 | ON TRACK | Same matrix: omitted existing-file digest observed; missing file yields both digest fields null and no Event. |
| VAL-022 | 7, 8 | NOT ON TRACK | Clean independent binding/Budget comparisons and stored/Event equality exist, but the dirty Run/Gate rows use the production observer for part of the expected repository. F1/F2. |
| VAL-023 | 7, 8 | ON TRACK | Normal, non-Git and unmerged states for start, non-null Budget, gate and evaluate; exact unavailable reason/details and no effects. Null Budget remains accepted without Git. Terminal leaves use VAL-061. |
| VAL-024 | 7, 8 | ON TRACK | Match/HEAD-differs/dirty-differs matrices; refusals leave persistence unchanged. Gate mismatches additionally leave no Attempt/process marker. |
| VAL-025 | 7, 8 | ON TRACK | Exact closed expected/observed mismatch details and expected input echo; C7U/C7I/C8U pass. |
| VAL-026-028 | 6; message half 10 | DEFERRED / context | Registration/Envelope/Claim acceptance belongs to Window B. Four plain-copy message changes remain task 10; no waiver or new finding against task 6. |
| VAL-029 | 7, 8 | ON TRACK for current allocation | Original/plain-copy/real-import controls preserve valid Claim/Registration prerequisites, refusal ordering, fixed pre-task-10 JSON messages and no effects. Final domain-cause messages remain task 10. |
| VAL-030 | 7 | NOT ON TRACK | Run admission after commit/edit is reproduced; full independence of recorded dirty repository is missing. F1. |
| VAL-031 | 7 | NOT ON TRACK | Close/abandon base=HEAD, identity and terminal/Event consistency are checked, but complete independent dirty repository equality is missing. F1. |
| VAL-032 | 8 | NOT ON TRACK | Gate records consistent pre-spawn Attempt/Receipt/Event maps and its process dirties the tree afterward. The expected dirty digest is production-derived. F2. |
| VAL-033 | 8 | ON TRACK | Unchanged/edited/reverted Receipt freshness and real terminal-transition/Envelope-completion denial/acceptance states pass, including archive round-trip fixtures. |
| VAL-034 | 8 | ON TRACK | Receipt, Review and Human Outcome subject repositories drive target/Event repository even after current state changes; linking works without Git. Stored get remains distinct from current evaluate. |
| VAL-035 | 7, 8; shared impact 9 | ON TRACK for this slice | Derived fake/fid Run/Gate sets and current changed-fake assembly pass; unrelated final whole-tree acceptance remains D1/E. |
| VAL-052 | 7, 8; impact 9 | ON TRACK | Exact-gate timing receipts below are within ordinary/tagged budgets and warning thresholds. Fresh focused reruns also terminate; they do not replace full-gate timings. |
| VAL-060 | 7, 8 | ON TRACK | Each close's adopted example process corpus checks non-project success, empty stderr and unchanged tree; current whole-tree descriptor process corpus passes. |
| VAL-061 | 7 | ON TRACK | Close/abandon mismatch versus non-Git/unmerged pairs return repository-mismatch versus recovery-state-conflict without Events. |
| VAL-062 | 7 | ON TRACK | C7I includes separate unchanged and tracked-edit close-proof fixtures: succeeded versus failed against Receipt freshness. |
| VAL-066 | 7, 8 | ON TRACK | C7I/C8I unavailable/advanced-HEAD/deleted-file replay plus otherwise viable new-key controls; complete result replay and no effects are asserted. |

## Gate And Evidence Checks

### Exact required gates versus reports

All three specs require the following five byte-identical commands from repository root; task 9 adds R5 absence. Both `## Tests run` and `## Exact quality-gate results` were checked, including each task's required inventory commands and preserved failed attempts.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
```

| Check | Report evidence inspected | Independent result / limits |
|---|---|---|
| Task 7 final gates | `task-7-done.md:167-173`: five final PASS receipts, tokens `exit 0 (no output)`, `0 issues.`, ordinary workflow `121.978s`, total `61.1%`, tagged `344.128s`; each reports terminated exit 0 | Close report-lint COMPLETE; exact format gate independently reproduced, silent success with observed `exit_code: 0`. R5/R6/R7 receipts are present at lines 147-151. |
| Task 8 final gates | `task-8-done.md:256-267`: five final PASS receipts, same format/lint tokens, ordinary workflow `103.583s`, total `61.1%`, tagged `340.478s`; each reports terminated exit 0 | Close report-lint COMPLETE; exact format gate independently reproduced in C8, silent success with observed `exit_code: 0`. R5/R6/R7 receipts are present at lines 207-210. |
| Task 9 current gates | `task-9-done.md:191-200`: five final PASS receipts, ordinary workflow `102.835s`, total `61.8%`, tagged `278.829s`; each reports terminated exit 0; declared R5 absence reports expected/observed exit 1 | Fresh report-lint COMPLETE; fresh R5 search has no matches and observed `exit_code: 1`. Current focused shared/caller tests pass. |
| Immutable earlier failures | Task 7 unused-helper lint/report formatting failures; task 8 initial incomplete checkpoint, typecheck/unused and retained-fixture failures; task 9 ordinary/R5 failure and bounded shared extension | Failures remain distinct from final evidence with causes, repairs and later identical-command results. Task 8's unverified earlier attestation is explicitly displaced by continuation receipts. No missing final required command found. |
| Report lint | Exact two-argument `project-docs/tools/report-lint.sh <report> <task spec>` for task 7 in its close worktree, task 8 in its close worktree, task 9 in H | Each prints `OUTCOME COMPLETE` and last token `PASS report-lint completed with no failures`, observed `exit_code: 0`. This checks report mechanics, not behavioral acceptance. |

The full ordinary coverage, repository lint, full tagged gate, benchmarks, scale checks and phase-level checks were not independently rerun. The report receipts above are report evidence, not a fresh QA attestation of those executions. Focused race/real-process reruns address the window's risk without rerunning full phase gates.

### VAL-052 timing receipts

These are the final exact-gate measurements captured by the Developers and inspected in the committed reports, not substituted timings from focused QA selections. Cached packages supply no fresh time.

| Task | cmd/fid | internal/workflow | Slowest other ordinary | Tagged integration | Assessment |
|---|---:|---:|---|---:|---|
| 7 | 342.375 s | 121.978 s | internal/issues 177.502 s | 344.128 s | Ordinary <600 s; tagged <900 s; below 480/720 s warnings |
| 8 | 389.949 s | 103.583 s | internal/issues 198.966 s | 340.478 s | Same |
| 9 | 410.123 s | 102.835 s | internal/issues 199.339 s | 278.829 s | Same |

### Fresh reproduced commands

The following commands ran exactly as printed. C7/C8/H are the declared directories above. Every long nested session was polled until a real exit code arrived; empty output and wrapper completion were never accepted as termination. Tokens below retain the output's package/time suffix; `\t` represents the captured tab separator.

```text
C7U (C7):
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid ./internal/workflow -run '^(TestP24Run.*|TestP24AbandonRecordsObservedBase|TestP24BudgetObservesOmittedSourceDigest|TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestAdoptedExampleProcessExclusivity|TestP24JudgmentDescriptorSourcePositionsAndExamples)$'

C7I (C7):
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid/integration -run '^(TestP24Run.*|TestP22CloseRunUsesObservedRepositoryAndQualifyingProof|TestP22DualAdapterRecoveryConformance)$'

C7R (C7):
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid -run '^TestP24RemovedRunRepositoryMembers$'

C8U (C8):
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid ./internal/evidence -run '^(TestP24Gate.*|TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation|TestP24RemovedGateEvidenceMembers|TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestAdoptedExampleProcessExclusivity|TestValidateGateRunReferencesAuthorizesAgainstEveryEnvelopeGrant)$'

C8I (C8):
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid/integration -run '^(TestP24GateAndEvaluateUnsafeAndReplay|TestP24ReceiptFreshnessRealGuardedActions|TestP24ReviewAndHumanEvidenceSubjectRepositories|TestP24GateChangedFakeFIDInputParity|TestP23SharedAmendedLeafConformance|TestP23FreshnessAndOutcomeSupersessionRealProcessJourney)$'

H1 (H):
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid ./internal/cli ./internal/project -run '^(TestP24Run.*|TestP24RemovedRunRepositoryMembers|TestP24Gate.*|TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation|TestP24RemovedGateEvidenceMembers|TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestAdoptedExampleProcessExclusivity|TestP24WholeTreeDescriptorProcessCorpus|TestP24DescriptorVariantProcessCorpus|TestInput.*|TestRepository.*|TestStopHumanListsShowOpenResolvedAndEmptyState|TestPhase23HumanLeafTextMatchesRecordedFacts)$'

H2 (H):
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid/integration -run '^(TestP24Run.*|TestP24Gate.*|TestP24ReceiptFreshnessRealGuardedActions|TestP24ReviewAndHumanEvidenceSubjectRepositories|TestP24ChangedFakeFIDAssembly|TestP24ExecutedKeyedCorpusMatches49LeafInventory|TestP24All49KeyedLeavesRebuildInOneJournal|TestP23FailureCorpusHasEveryAcceptedCode|TestP23SharedAmendedLeafConformance|TestP23StopFailure.*|TestP23ConditionWireGrammarMatchesJournalAndFixtures|TestP23HumanFixtures|TestP23HumanFixtureStepsRemainRunnable|TestP23Concern.*|TestP23FindingDispositionRealProcessJourney|TestP23FreshnessAndOutcomeSupersessionRealProcessJourney|TestP23DecisionReferenceProtectedLifecycleRealProcess|TestP23AttentionResolution.*|TestP23AmendedEnvelopeReplayThroughBuiltBinary)$'

H3 (H):
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./internal/workflow -run '^(TestStopResolution.*|TestStopRaiseRejectsDecisionSatisfiedOnActivationSubject|TestProtectedDecision.*|TestStopAndAttentionListsFilterLimitAndRemainReadOnly)$'

H4 (H):
FID_P23_HUMAN_FIXTURE_ROOT=/tmp/fid-review5-human-foundation go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23HumanFixtures$' -count=1 -v ./cmd/fid/integration

O (H):
python3 /tmp/fid-review5-oracle-audit.py
```

| ID | Result | Fresh output token and observed termination |
|---|---|---|
| C7U | PASS | `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t3.138s`; cmd/fid 13.262 s; terminated `exit_code: 0` |
| C7I | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t35.504s`; terminated `exit_code: 0` |
| C7R | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t2.374s`; terminated `exit_code: 0` |
| C8U | PASS | `ok  \tgithub.com/mbell-nzl/fid/internal/evidence\t1.030s`; cmd/fid 13.805 s; terminated `exit_code: 0` |
| C8I | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t7.421s`; terminated `exit_code: 0` |
| H1 | PASS | `ok  \tgithub.com/mbell-nzl/fid/internal/project\t1.078s`; cmd/fid 34.013 s, internal/cli 1.016 s; terminated `exit_code: 0` |
| H2 | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t48.765s`; terminated `exit_code: 0` |
| H3 | PASS | `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t15.220s`; terminated `exit_code: 0` |
| H4 | PASS, predecessor foundation only | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t18.264s`; terminated `exit_code: 0` |
| O | Evidence gap reproduced | `ORACLE AUDIT: two bounded evidence gaps reproduced at close commits and preserved at reviewed HEAD`; terminated `exit_code: 0` |

O reads the exact close-commit blobs with `git show`, verifies that `testEnvelopeRepository` calls `project.ObserveRepository` and copies its DirtyDigest, verifies both dirty Run expected maps call that helper, verifies the terminal matrix omits an independent dirty-digest assertion, and verifies the Gate pre-spawn expected map has the same dependency. Findings are source/evidence defects; these tests' green status does not contradict them.

The first C7U attempt in `/tmp/fid-p24-review5-task7-9b585b7` terminated `exit_code: 1`, token `FAIL`, with `error obtaining VCS status: exit status 128` during process-corpus binary builds. `go build -x` showed Git running in `/tmp`; explicit worktree Git commands succeeded. The byte-identical C7U command then passed in the same-commit clone above. No VCS flag, test, gate definition or product change was used to obtain that result. The failed reviewer attempt remains evidence of the inspection environment.

## Task 7 - Run leaves

### Must-fix

- **F1 — Must-fix; VAL-022, VAL-030, VAL-031.** `cmd/fid/p24_run_test.go:31` and `:96` construct dirty start/close expectations using `testEnvelopeRepository`, whose close-commit implementation calls the production `project.ObserveRepository` and copies its DirtyDigest (`cmd/fid/envelope_test.go`, helper definition). The start test independently checks HEAD/base/tree/dirty/root, but never independently derives the dirty digest or config identity for its full expected map. The close test compares its returned map with that same production-derived expected map. `cmd/fid/integration/p24_run_inputs_test.go:72` independently checks the terminal identity fields for close and abandon, omits `dirty_digest`, and at `:79` compares the Event to the returned record. A wrong shared dirty digest can therefore satisfy all these comparisons. Owning task: task-7. Repair boundary: expected-value construction and complete record/Event comparisons in `cmd/fid/p24_run_test.go` and `cmd/fid/integration/p24_run_inputs_test.go`, plus `task-7-done.md`; test-local helpers may be added inside those files. Derive the complete repository from independent Git/config/root/file-image evidence, without calling production repository observers for expected values. Required evidence: all eight repository fields, including independently hashed dirty image and project ID from config, compared with start/close/abandon records and actual Events after the prescribed commit/edit; retain match/refusal/replay/Source pairs; rerun the affected named tests, five exact task gates and report-lint with fresh captured termination. Re-review scope: inspect the independent oracle at the repair commit and rerun affected ordinary and real-process Run tests, including both terminal leaves. Fresh independent Window C review required.

### Improvement

- None.

Observation (design sections 2-3, 6-7): the accepted identity-free inputs are placed in descriptors while expectation/worktree rules stay in workflow/project domains. Start accepts an edited later tree, terminal actions observe under the lock after replay, and Budget observes only for a non-null Source. Unsafe terminal requests preserve recovery-state-conflict. Source digest, import/plain-copy, no-effect and Receipt-proof pairs assert visible behavior rather than helper implementation. The changed CLI/flow documentation matches those rules; F1 concerns the independence of evidence, not a demonstrated runtime failure.

## Task 8 - Gate run and Evidence leaves; real freshness

### Must-fix

- **F2 — Must-fix; VAL-022, VAL-032.** `cmd/fid/p24_gate_evidence_test.go:32` obtains the complete pre-spawn expected map through `testEnvelopeRepository`; its expected `dirty_digest` is therefore returned by the production observer. The independent Git checks at `:45` do not recompute that digest or independently check config identity. Receipt/Attempt-effect/Event equalities at `:62-87` consequently prove consistent propagation of a production-derived map, rather than equality to the contract's complete independent pre-spawn observation. Owning task: task-8. Repair boundary: this pre-spawn fixture and its test-local oracle/comparisons in `cmd/fid/p24_gate_evidence_test.go`, plus `task-8-done.md`; preserve Gate/Evidence production behavior. Required evidence: independently obtain all eight repository fields for the tracked-file edit, compare the persisted Attempt start repository, Receipt terminal repository and both actual Events, with non-vacuous record/Event counts; retain post-spawn tree change and mismatch no-Attempt/no-process/no-Event controls; rerun the named Gate/Evidence tests, five exact task gates and report-lint with fresh captured termination. Re-review scope: verify oracle independence and complete stored record/Event equality, and rerun affected ordinary/real-process Gate/Evidence/freshness tests. Fresh independent Window C review required.

### Improvement

- None.

Observation (design sections 2-3, 5-7): execution observes once before spawn and reuses that snapshot at terminal recording; replay precedes observation. Binding checks use shared worktree identity rather than HEAD/tree equality. Evidence linking derives the target from each supported subject and works without current Git, while evaluate observes current state. Fresh/edited/reverted guarded-action tests and no-effect mismatch checks pass. Help and CLI/flow docs agree with the accepted read-versus-observe distinction. No terminal re-observation, live-frame change or new dependency was introduced.

## Task 9 - Current-HEAD impact only

### Must-fix

- None within Window C's impact scope.

### Improvement

- None.

Observation (design sections 2-3, 6-7): the DEC-284 four-path shared extension adds selected typed raw-array descriptor validation and sensitive-value handling, then adopts the last Activation component descriptor and removes the obsolete decoder fence. The diff from task 8 close to supplied HEAD contains no change to `internal/workflow`, `internal/evidence`, Run/Gate/Evidence/Artifact command implementations, `go.mod` or `go.sum` (exact scoped `git diff --stat` was empty, observed exit 0). Shared descriptor metadata and retained adapter fixtures do change, so H1/H2 exercise their actual impact on the reviewed families. They pass and preserve F1/F2's evidence limitations. Task 9's COMPLETE report retains its failed ordinary/R5 attempts and the authorized repair. D1 must still review task 9 in full with task 9a.

## R1-R20 Inventory Reproduction

All registry queries ran freshly from H. The inventory's Markdown-escaped alternation was rendered as regex alternation. Commands were executed separately; the two R4 operations each have their own result. R15-R18 and the input set came from the exact accepted recursive section-1 block, extracted unchanged into `/tmp/fid-review5-inventory-1ti4a98y/walk.sh` and executed with `bash /tmp/fid-review5-inventory-1ti4a98y/walk.sh`. Its last token was `271 total`, observed `exit_code: 0`; scratch binary and labelled lists are `/tmp/tmp.k5Rlm9mZax/`. Counts are 92 leaves, 29 groups, 39 input, 50 keyed, 0 no-short, 61 typed flags.

```text
R1: rg -n "RequestDigest:" internal cmd -g '!*_test.go'
R2: rg -n "func [a-zA-Z]*(RequestDigest|requestDigest|judgmentDigest)\(" internal -g '!*_test.go'
R3: rg -n "idempotency_results" migrations internal -g '!*_test.go'
R4a: ls migrations
R4b: rg -n "ProjectionSchemaVersion|projection_schema" internal -g '!*_test.go'
R5: rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid
R6: rg -n "func (same[A-Za-z]*Repository|validateRepositoryBinding|validEnvelopeRepository|validEnvelopeSource|validEvidenceRepository|validateRunTerminalRepository|verifyDecisionSource|importedGateDestination)" internal -g '!*_test.go'
R7: rg -n "Request\.Repository|request\.Repository|Request\.TargetRepository|request\.RepositoryScope|Request\.WorktreeID|request\.WorktreeID" internal -g '!*_test.go'
R8: rg -n 'definition\(' internal/cli/codes.go
R9: rg -n "func (write[A-Za-z]*Failure|render[A-Za-z]*Error)\(" cmd/fid -g '!*_test.go'
R10: rg -n '"error: ' cmd internal -g '*_test.go'
R11: rg -n -i "worktree|dirty digest|HeadOID|head_oid" cmd/fid -g '!*_test.go'
R12: ls cmd/fid/stop_human.go cmd/fid/resume_human.go cmd/fid/resume_human_guidance.go
R13: rg -n "p23CheckJudgmentQualifiers|PRIVATE_" cmd
R14: rg -n "func TestLiveLeafCommandsMatchContractInventories" cmd/fid
R19: rg -n '"argument": *"request"' cmd/fid -g '!*_test.go'
R20: rg -n 'CodeInvalidArgument, "invalid [^"]*"' cmd/fid -g '!*_test.go'
```

| ID | Fresh set / accounting | Verbatim output token; observed termination |
|---|---|---|
| R1 | 45 matches / 33 files. Compared with planning: execution adds one named preflight row; journal/idempotency adds the supplied/stored digest-conflict row. These account for 43/32 →45/33. | `internal/workflow/run_checkpoint.go:64:` prefix of its matching idempotency line; terminated `exit_code: 0` |
| R2 | One consolidated helper, ComputeRequestDigest; prior per-family helpers removed under task 2. | `internal/journal/request_digest.go:14:func ComputeRequestDigest`; terminated `exit_code: 0` |
| R3 | 126 matches / 18 files, including both new 0019 DDL files and closed validation/projection/update registries. | `internal/maintenance/operational.go:128:` prefix of its SELECT match; terminated `exit_code: 0` |
| R4a/b | 38 migration files, latest 0019; 17 schema matches / 5 files. Physical active migration is 19; internal ProjectionSchemaVersion is 5 (`internal/db/projection_meta.go:10-11`), not a claim that those distinct constants are equal. | `0019_p24_idempotency_fields.up.sql`; `internal/journal/projection.go:67:` prefix; each terminated `exit_code: 0` |
| R5 | No family decoder remains at H. Task 7 close had five, task 8 close four; final removal is task 9's owned absence. | `exit 1 (no output)`; terminated `exit_code: 1`, task 9 declares expected `exit_code: 1` |
| R6 | Four retained helpers in execution, run_read and envelope. SameWorktree replaces caller-state validators/import helper; resolved internal validators remain. | `internal/workflow/envelope.go:759:func validEnvelopeSource(source Source) bool {`; terminated `exit_code: 0` |
| R7 | 16 matches / 2 files, registration and envelope. These private resolved facts are assigned from observation; corresponding removed public JSON fields are rejected. No Run/Gate/Evidence caller repository writer remains. | `internal/workflow/envelope.go:935:` prefix of repository serialization; terminated `exit_code: 0` |
| R8 | 84 query lines = 83 code definition calls plus the helper declaration; 82 at planning, with repository-unavailable added. Operation-unavailable remains task 9a. | `361:` prefix of CodeImportTargetConflict definition; terminated `exit_code: 0` |
| R9 | 42 functions / 32 files, including task 3's writeExampleConflict in addition to the planning 41. Task 10 must account for that extra caller. | `cmd/fid/label_response.go:11:func renderLabelError`; terminated `exit_code: 0` |
| R10 | 33 matches / 14 files, including the adopted example conflict assertion beyond planning 32/13. | `cmd/fid/link_test.go:352:` prefix of retained error assertion; terminated `exit_code: 0` |
| R11 | 19 matches / 11 files. Four existing repository renderers plus guidance; remaining matches include JSON/input/help and integration JSON fixtures, which the glob does not exclude. Task 11 must classify rather than pretend all hits are renderers. | `cmd/fid/integration/testdata/p21_conformance/manifest.json:144:    "worktree"`; terminated `exit_code: 0` |
| R12 | All three named renderers exist. | `cmd/fid/stop_human.go`; terminated `exit_code: 0` |
| R13 | 15 matches / 3 files. Existing guard and recorded-reasoning PRIVATE values are still present for task 13 convergence. | `cmd/fid/integration/testdata/p23_conformance/operations.json:145:` prefix; terminated `exit_code: 0` |
| R14 | One live-leaf completeness test. Three new get leaves are still pending task 12. | `cmd/fid/contract_surface_test.go:128:func TestLiveLeafCommandsMatchContractInventories(t *testing.T) {`; terminated `exit_code: 0` |
| R15 | 92 leaves; no missing Short. | Section-1 output `0 /tmp/tmp.k5Rlm9mZax/no-short`; walk terminated `exit_code: 0` |
| R16 | 50 keyed, minus export import =49 canonical digest leaves; H2 checks executed-set completeness and single-journal rebuild. | Section-1 keyed count 50; walk last token `271 total`, terminated `exit_code: 0` |
| R17 | Root +28 groups =29; help/completion filtered from traversal. | Section-1 groups count 29; walk last token `271 total`, terminated `exit_code: 0` |
| R18 | 61 =17 int +44 bool, including all 39 adopted examples. Original 22+39 explains the increase; task 9a must enumerate live flags, not pin 22. | Section-1 typed-flags count 61; walk last token `271 total`, terminated `exit_code: 0` |
| R19 | Seven sites / six files: three grammar and four family fallbacks. Exactly four is task 9a's future close condition, not today's feasibility pass condition. | `cmd/fid/issue_get_list.go:129:` prefix of fallback; terminated `exit_code: 0` |
| R20 | 25 sites /20 files =18 generic, three grammar, four specific. | `cmd/fid/gate.go:312:` prefix of generic fallback; terminated `exit_code: 0` |

The R20 generic set is activation, claim, concern, decision, deviation, envelope, evidence, finding, gate, human_outcome, issue_get_list, issue_mutation_response, link_response, receipt, run, stop and both review mappers. A future row or explicit unreachability reason is required for every site. The three grammar and four specific branches are separate controls, not substitutes for generic-fallback coverage.

## Tasks 9a-14 Command And Corpus Feasibility

### Non-default commands

Every non-default command in the task specs is concrete and executable against existing roots. R19, R9/R20 and R11 were run byte-identically above for tasks 9a, 10 and 11. Task 12 has no extra command. Tasks 13 and 14 prescribe the same two absence queries, both run freshly from H:

```text
rg -n "This (Concern|Deviation) is recorded judgment|This Outcome is authorized Human judgment" cmd/fid -g '!*_test.go'
rg -n "PRIVATE_" cmd/fid/integration
```

Current observations: the first returns three matches, last token `cmd/fid/p23_human.go:176:` followed by the Outcome disclaimer, observed `exit_code: 0`; the second returns 12 matches in three files, last token `cmd/fid/integration/testdata/p23_conformance/operations.json:145:` followed by PRIVATE_OBSERVED, observed `exit_code: 0`. Both commands are feasible. Their future declared success is absence with exit 1 after task 13, re-proved by task 14; neither current result is falsely recorded as a passing absence gate.

The exact D2 generator command is:

```text
FID_P24_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-24-human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration
```

**NOT RUN:** TestP24HumanFixtures is a task-11 deliverable and is absent at the reviewed HEAD. Fresh `rg -n '^func TestP24HumanFixtures\(' cmd/fid/integration` returned no output, observed `exit_code: 1`. Running that Go selector now would risk a misleading no-tests success. H4 instead proves the actual predecessor generator mechanism works; it does not satisfy the future D2 command. H4 retained `/tmp/fid-review5-human-foundation/checkpoints-156568860/manifest.json`: source commit is the supplied HEAD, source_dirty is true because the runner ledger is dirty, three checkpoints and 18 executable Human/JSON leaf pairs, binary digest `sha256:899bfd971b0d270823322c72c9d9c76c22abc7284d91fd3dca9b470b7a64c24c`. This is deterministic fixture infrastructure, not a Human walkthrough or a Phase 24 fixture manifest.

### Fresh process corpus probe

Executed `python3 /tmp/fid-review5-feasibility.py` from H against the freshly built inventory binary `/tmp/tmp.k5Rlm9mZax/fid`. The final probe allocates a fresh project and a fresh non-project child directory per non-project invocation, captures each real subprocess argv/cwd/stdout/stderr/exit, and exercises both Human and JSON modes. It terminates `exit_code: 0` with exact token:

```text
FEASIBLE: 90 leaves, 29 groups, 61 typed flags; project/nonproject and both modes; sensitive guard reaches generic fallback without leaking sentinel
```

There are 1,299 receipts at `/tmp/fid-review5-feasibility-b06l3v5f/receipts.json`, SHA-256 `d7c809d65eef04490b229a96212badf40a64ce1902112a9e55f4d9df87591aee`. This probe establishes executable fixtures and baseline states, not final task-9a/10 grammar or rendering conformance.

| Corpus | Exact argv recipe / scope | Observed baseline termination |
|---|---|---|
| Every leaf unknown flag | `fid [--json] <leaf> --bogus`, 90 live leaves excluding journal inspect/sql ×two modes ×project/non-project | 360 calls, exit 2 each |
| Every leaf positional refusal | `fid [--json] <leaf> extra extra`, same set | 360 calls: 358 exit 2; initialized-project init in both modes exit 4. Grammar-before-workspace remains task 9a. |
| R17 unknown commands | `fid [--json] <group> bogus`, root included, both directories | 116 calls: 92 exit 0/help, 24 exit 2. This reproduces the accepted 23 versus six group baseline, not operation-unavailable acceptance. |
| R17 unknown flags | `fid [--json] <group> --bogus` | 116 calls, exit 2 each |
| R18 malformed values | `fid [--json] <leaf> <flag>=maybe` for every live typed flag | 244 calls: 232 exit 2; 12 exit 7 for has-concerns, rebuild verbose and requirement-list limit across directories/modes. Task 9a owns normalization. |
| R17 no-word controls | `fid <group>`, Human, both directories | 58 calls, exit 0 each |
| Domain/missing/near-miss controls | Both modes: `issue create`, `stop raise`, `issue list --order bogus`, `run list --status bogus`, `issue list --limit`, `isue` | 12 calls: ten exit 2, two issue-order non-project calls exit 3/project-not-found. Accepted future pre-workspace validation remains task 9a. |
| Message boundaries | Both modes: `stop raise --<word>` with bogus, x×496/497/600, 界×496, x×494+ESC+z; root word x×484/485/600, 界×484, x×482+ESC+z | 12 flag and ten word calls, exit 2 each. Inputs are executable; final bounded/escaped pair equality remains future implementation evidence. |
| Sensitive generic fallback | `activation prepare --example`, substitute accepted setup; prepare and first gate-definitions component succeed; same component identity/version under fresh JSON/Human keys, input environment literal contains a distinct sentinel | Template exit 0; setup/component exit 0; both duplicate controls exit 2, no sentinel in either stream, JSON generic `argument:input`. |

The sensitive row genuinely reaches activation's domain/event validation, not descriptor validation. Human token is `error: current Event activation.component-staged idempotency_results/insert: invalid JSON column "result_json": activation: array must be sorted and unique`. JSON retains `invalid activation input`, allowed accepted activation input, value null, invalid-argument/exit 2. The planning snapshot's UNIQUE-constraint cause no longer describes this probe, but the accepted guard requirement is generic-fallback reachability, which remains feasible.

Probe history is retained: first attempt used unsupported network/write `deny` and terminated exit 1 at an assertion after the CLI returned `input.content[0].network`, allowed `["allowed"]`, exit 2. After correcting only scratch fixture values, the second probe terminated exit 1 because it still expected the planning UNIQUE text; its receipts are `/tmp/fid-review5-feasibility-zwxn9vve/receipts.json`. A third execution completed with exit 0 in `/tmp/fid-review5-feasibility-tcxh426g/`, but its reused non-project directory could be initialized by an init baseline. The final execution above corrects that isolation by allocating every non-project cwd fresh. These are reviewer harness attempts, not failed product gates or changes to accepted pass semantics.

### Coverage of every prescribed later corpus

| Owner / corpus | Fresh feasibility evidence | What was not run and why |
|---|---|---|
| 9a: live-leaf grammar, root/groups, typed and declared string domains, missing values, positional counts, no-project/no-write controls, operation IDs | Full live lists and baseline process matrices above; accepted Phase 21-24 row references and live flag registrations exist. Representative free-text/missing and closed-domain probes execute. Current extra example flags and baseline exit/read differences are recorded. | The future complete declared string-domain registry, all per-leaf required-flag rows and final no-effect/metadata assertions are not implemented yet. Task 9a must enumerate them from accepted row citations, retain CF-047 combination exclusions and repeat each in isolated non-project state. Feasibility is not a narrower substitute for that corpus. |
| 9a/10: escaped/bounded messages, 512/513, Unicode, escape crossing, already-escaped %q | Fresh boundary inputs execute; QA-plan section 12.7.2 supplies deterministic scalar recipes. Existing live-text escape implementation is inspectable. | Shared message-construction unit assertions and final Human/JSON equality await task 9a, then task 10 adoption. No bounded-message PASS is claimed today. |
| 10: RegisteredCodes Human/JSON rows, grammar versus structural versus terminal frame, exclusions | R8/R9/R10 inventories; H2 registered-code completeness, stop-failure pairing and Gate/real-guard foundations pass. Existing real-process adapters reach typed producers and streamed Gate outcomes. | New renderer pairing/exclusion corpus is task 10 work. It must cover all 83 current codes plus task 9a's new code, every R9 writer, five no-next codes and both blocked-by states; old pairings do not prove the new line form. |
| 10: all R20 domain fallbacks, plain-copy messages, sensitive guard, 200/201 reason cuts | Eighteen exact generic sites mapped above; sensitive domain fallback freshly reached in both modes without leaking; C7U/C8U original/copy/import controls and H2 real guard/Stop fixtures pass. Contract gives closed overrides and next table. | Each final fallback row or pinned unreachability reason, empty/whitespace causes, scalar/array overrides and reason boundaries are task 10's new assertions. Preserve JSON metadata/details/exit and Claim-first Run refusal. |
| 11: ten StopDetails variants ×ten Conditions, open/resolved, next only open, all repository views clean/dirty and same/copy | R11/R12 close the renderer set. H3 runs actual resolution predicate and read-only tests; H2 runs closed Condition grammar and Stop/Attention real-process tests. C7U/C8U prove copy/import setup; existing resume/review read paths accept repository context. | New labelled-line and short-HEAD rendering corpus is not implemented yet. Every closed member and both statuses require task 11 evidence, including no raw JSON/no nil/internal identity tokens. |
| 11: Phase 24 fixture generator and D2 Human checkpoint | H4 executes predecessor generator with persistent binary/manifest and 18 re-executed pairs; H2 independently runs the predecessor six-step fixtures. All required fixture capabilities have real producing leaves. | Exact P24 generator and five-denial/four-view manifest are absent until task 11. D2 must regenerate at task-11 close, then collect Human judgment separately. |
| 12: three get/list-equality, event:null, known/unknown IDs, missing/extra/unknown flags | Current Concern/Deviation/Outcome records and list fixtures are produced by H2/H4; existing review/requirement get patterns and R14 live-tree test are present. | The three new get leaves do not exist yet. Their exact reads/not-found/grammar pairs and new live-tree/pairing rows remain task 12 implementation evidence; task has no non-default command. |
| 13: every judgment view with 0/1/3 items; all closed reasoning fields 80/81/multiline; ESC/CR/BEL/NUL; source-content sentinel; glossary distinction | R13 closes existing guards/fixture references. H1/H2/H4 reproduce current human fact projections, Concern/disposition/Deviation and 18 leaf pairs with actual referenced files. Required text/control recipes are ordinary JSON strings and need no external facility. | New section-disclaimer function, enlarged detail views, reasoning boundaries/controls and replacement sentinel sweep await tasks 12/13. Current PRIVATE values are recorded reasoning; they must be renamed with provenance, separately from project-file sentinels. Absence queries currently match as recorded above. |
| 14: protected-Decision journey to completion, five provenance classes, no Receipt obligation, input files outside tree | H2 reproduces protected Decision lifecycle, executable Human fixture steps and amended Envelope/replay foundations. H4 retains actionable protected-Decision workspace/steps. Task-9 descriptors now supply templates across the tree; the five allowed classes are explicit. | The all-Human journey with per-value source-line provenance needs task 10/11 guidance and task 14's test. No current run is relabelled as that journey. A missing command/value source would require the owning authority, not widened provenance. |
| 14: semantic rebuild/export/validate/import, schema 19, full forbidden-member scan, generator extension | H2 executes all-49 single-journal rebuild and freshness round-trip foundations; R3/R4 verify accepted field-digest storage/schema roots. H1 checks removed members and complete adopted descriptors. | Full take-control workspace parity, three new get reads, final retirement scan and extended P24 manifest require future product/fixture completion. Same two absence queries are runnable but not yet absent. No full phase-level scale/checkptr/scenario gate or Human judgment was run. |

## Repair Plan

Repair order is task-7, then task-8. Both findings concern deterministic evidence required by the already accepted contract. No product, architecture, public-interface, dependency, gate command, threshold, policy or assertion-definition change is authorized by this review.

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| F1 | task-7 | VAL-022, VAL-030, VAL-031 | `cmd/fid/p24_run_test.go`, `cmd/fid/integration/p24_run_inputs_test.go`, and `project-docs/docs/phases/phase-24/tasks/task-7-done.md`; test-local oracle/comparison work only | Full independent Git/config/root/dirty-image expected repository for dirty start, close and abandon; equality with stored record and actual Event; retained qualifying/refusing/replay states; affected named tests, all five exact task gates and fresh COMPLETE report-lint, with observed termination | Fresh independent reviewer checks repair diff/oracle independence, both terminal leaves and complete field/Event comparisons, and reruns representative affected ordinary/tagged verification |
| F2 | task-8 | VAL-022, VAL-032 | `cmd/fid/p24_gate_evidence_test.go` and `project-docs/docs/phases/phase-24/tasks/task-8-done.md`; pre-spawn test-local oracle and persisted comparisons only | Independent full dirty pre-spawn map against persisted Attempt/Receipt/both Events, explicit non-vacuous counts, post-spawn change retained, mismatch no-effects retained; affected Gate/Evidence tests, five exact gates and fresh COMPLETE report-lint, with observed termination | Fresh independent reviewer verifies oracle and complete persisted equality, reruns representative affected Gate/Evidence/freshness verification, and checks shared current-HEAD impact |

The fresh review must cite review-5, preserve each finding as RESOLVED/STILL OPEN/protected blocker candidate with fresh evidence, and write a new artifact. This reviewer does not fix tests, update reports or resolve its own findings.

## Prior Finding Resolution

This is the first Window C review. Review-4's Window B PASS and its F1-F8 resolutions remain intact; no prior C finding is being resolved. In particular, its independent dirty Review oracle provides a useful existing fixture pattern, not acceptance evidence for the distinct dirty Run/Gate rows. Earlier Runner bookkeeping Improvements remain historical context. No earlier review was appended, rewritten or superseded by self-correction.

## Cross-cutting observations

DEC-284 explicitly authorizes this late Window C after task 9 COMPLETE; it records that the Runner missed the planned after-task-8 checkpoint. The current report/commit chain confirms task 9 has since completed even though the current brief/progress snapshot still describes its preceding recovery. That snapshot is context, not an extra protected blocker. This review does not edit it or the ledger.

The late checkpoint changes neither C's acceptance contract nor the remaining cadence. D1 reviews tasks 9/9a after 9a; D2 reviews tasks 10/11 after 11 and regenerates the task-11-close Human fixture output for its walkthrough; E reviews tasks 12-14 after 14. D2 judgment also does not replace phase-end HUMAN-02/HUMAN-03. This REPAIR must be closed by fresh independent QA before Window C advances.

## Blocker Candidates

- None. F1/F2 are bounded in-scope evidence repairs; no protected decision or external capability is needed.

## Existing Escalation References

- ESC-24-001 and ESC-24-002: resolved under DEC-283; historical Task 6 recovery only.
- ESC-24-003: resolved by DEC-284, correcting Task 9 TRIAGE without a verifier change and assigning the late Window C.
- Current escalation status: none open. No escalation ID was allocated and no register was edited.

## Final Lines

Review written: project-docs/docs/phases/phase-24/qa/review-5.md
Must-fix: 2  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
