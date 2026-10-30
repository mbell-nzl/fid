# Phase 24 - Mid-Phase Review 6 (tasks 7-8)

Reviewer: qa-reviewer
Date: 2026-10-07
Tasks covered: phase-24/task-7, phase-24/task-8; task-9 current-HEAD impact and remaining cadence/feasibility only
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-24/qa/review-5.md

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-24-001, ESC-24-002, ESC-24-003 (resolved history; none open)
- Covered validation assertions: VAL-009, VAL-011, VAL-012, VAL-013, VAL-014, VAL-015, VAL-016, VAL-017, VAL-018, VAL-019, VAL-022, VAL-023, VAL-024, VAL-025, VAL-029, VAL-030, VAL-031, VAL-032, VAL-033, VAL-034, VAL-035, VAL-052, VAL-060, VAL-061, VAL-062, VAL-066

**Window C PASS.** Fresh independent inspection and ordinary/real-process race verification resolve review-5 F1 and F2. Complete dirty repository expectations now come from test-local Git/config/file-image evidence. Run records and actual Events agree with those expectations; Gate's persisted Attempt, persisted Receipt, both Events and repository Effects agree with the independent pre-spawn map, with explicit non-vacuous counts. No Must-fix finding or blocker candidate remains.

This closes the Window C repair requirement. The Phase Runner owns advancement. Window D1 after task 9a, D2 after task 11, and E after task 14 remain mandatory. Task 9 receives only shared current-HEAD impact assessment here; its full acceptance remains D1. No Human judgment or phase-end acceptance is recorded.

## Scope Reviewed

Loaded AGENTS.md, the qa-evidence skill, workflow sections 2.4/2.7/8, the four requested current context documents, the phase validation contract and applicable QA-plan evidence/gates/cadence, review-5 in full, task specs/reports 7-9, phase decisions including DEC-273/278/280/281/282/284, relevant triage and resolved escalation history, the review template and quality plan. Consulted review-4 for preceding-window boundaries and evidence patterns. CLI observations use `05_DESIGN_PLAN.md` sections 2-3 and 5-7, plus the Run/Gate/Evidence CLI and flow amendments. No held-out material was read.

Reviewed branch `phase-24` in `/home/mark/projects/fid` at supplied HEAD `95280952c09d1d7996ef9f26e569dc16230e521b`. The worktree started with only the Runner-owned `project-docs/docs/phases/phase-24/run-ledger.md` dirty. It remains untouched. No product, test, report, task, decision, current document, plan, contract, tool or ledger was edited by this reviewer.

| Repair | Independently inspected commit and subject | Bounded paths verified |
|---|---|---|
| F1, task-7 | `cdb7e01f630b3e10291237d5912e63ddf0b33783`, `test(run): independently verify dirty repository observations`; report finalization `856c36d6be4457c36e718fb0bb376bf705025949` | Exactly `cmd/fid/p24_run_test.go`, `cmd/fid/integration/p24_run_inputs_test.go`, `project-docs/docs/phases/phase-24/tasks/task-7-done.md` |
| F2, task-8 | `a59fa4d1dbca96193e257d6b55d3fa67cab619fc`, `fix(test): independently verify gate pre-spawn repositories`; report finalization `95280952c09d1d7996ef9f26e569dc16230e521b` | Exactly `cmd/fid/p24_gate_evidence_test.go`, `project-docs/docs/phases/phase-24/tasks/task-8-done.md` |

Each repair has its owning `Refs: phase-24/task-N`. `git show --stat`, full test/report diffs and exact path-set comparisons against review-5's Repair Plan confirm the scope. The ranges `c690cb93..856c36d6` and `856c36d6..95280952` change exactly three and two allowed paths respectively; each is a regular `100644` blob and not a symlink. No production observer, canonicalizer, gate/control definition, dependency or accepted assertion changed. Existing qualifying/refusing/replay assertions remain in the test diffs.

Ran `git rev-parse --verify '<hash>^{commit}'` separately for the two repair commits, both report-finalization commits, prior review commit `c690cb939255c63b546180b4477da974a51b58b4`, task-9 source checkpoints `9d654e976792836cbdf7de8c114d0831b73f721e` and `69589747ca77a1a3dae7e4d6083229d0193ed52f`, task-9 report commit `d6e03c9f9ec7a445f3af50451cb9b14d4d52dd5c`, and task-9 close `e39bfcc11237e108cebb58d46aa5dc696ab44ca1`. Each returned the requested full commit hash with observed `exit_code: 0`.

## Validation Contract Coverage

ON TRACK is limited to the owned Window C slice. These results do not reopen Window A/B or certify task 9's entire slice. Test selections below exercise both qualifying and non-qualifying states; positive-only propagation equality is no longer the dirty-observation evidence.

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-009 | 7, 8 | ON TRACK | I1 Run and Gate changed-HEAD/file keyed retries preserve complete original responses and Event references with no persistence change. |
| VAL-011 | 7, 8; impact 9 | ON TRACK | O1 live-tree/type registry and seeded-mismatch rejection, seven Run/four Gate-Evidence descriptors, and current 39-leaf whole-tree corpus pass. Whole-task-9 certification remains D1. |
| VAL-012 | 7, 8 | ON TRACK | O1 descriptor process corpora inspect recursive field help, required/nullable metadata and expectation Git commands. |
| VAL-013 | 7, 8 | ON TRACK | O1 real example output matches the required-only ordered, two-space template. |
| VAL-014 | 7, 8 | ON TRACK | O1 strict typed-template/descriptor checks retain actual request-type decoding and applicable qualifying placeholder expansion. |
| VAL-015 | 7, 8 | ON TRACK | O1 built-process exclusivity tests pair standalone success with all six prohibited combinations and no effects. |
| VAL-016 | 7, 8; impact 9 | ON TRACK | O1 selected structural corpora and shared validator tests retain exact nested paths/types/keys/formats/bounds and valid controls; selected variants are exercised. |
| VAL-017 | 7, 8 | ON TRACK | O1 removed Run/Gate/Evidence member tests retain typed path rejection and no writes. |
| VAL-018 | 7 | ON TRACK | I1 Budget Source matrix accepts matching digests and rejects mismatches with independent SHA-256 details and no Event. |
| VAL-019 | 7 | ON TRACK | I1 same matrix pairs omitted-existing digest observation with omitted-missing failure and both digest fields null. |
| VAL-022 | 7, 8 | ON TRACK | F1/F2 resolved: complete independent eight-field dirty maps equal stored Run records, persisted Gate Attempt/Receipt and actual Events. Budget's retained clean Source/Event comparisons also pass I1. |
| VAL-023 | 7, 8 | ON TRACK | I1 normal/non-Git/unmerged Run-start, non-null Budget, Gate and evaluate states retain mapped failures and no effects; null Budget succeeds without Git. Shared task-4 cause-class coverage remains preceding-window evidence. |
| VAL-024 | 7, 8 | ON TRACK | O1/I1 matching expectations versus HEAD/dirty differences; Gate refuses before Attempt/process/Event. |
| VAL-025 | 7, 8 | ON TRACK | O1/I1 mismatch details compare the complete closed expected/observed maps and echo the supplied expectation. |
| VAL-029 | 7, 8 | ON TRACK for current allocation | O1 original/plain-copy/genuine-import controls retain valid prerequisites, Claim-first Run refusal and no effects. Domain-cause JSON messages remain task 10. |
| VAL-030 | 7 | ON TRACK | O1 start after commit/edit records the independent current dirty map rather than the earlier Envelope snapshot. F1 resolved. |
| VAL-031 | 7 | ON TRACK | O1 close and I1 close/abandon after commit/edit compare all terminal fields and actual Events with independent maps, base=HEAD. F1 resolved. |
| VAL-032 | 8 | ON TRACK | O1 complete independent pre-spawn map equals persisted Attempt/Receipt, both actual Events and both repository Effects; exact counts and post-spawn tree change are required. F2 resolved. |
| VAL-033 | 8 | ON TRACK | I1 unchanged/edit/revert across evaluate, terminal transition and Envelope completion; stale denials carry evidence-stale. O2 independently reproduces ordinary fresh/stale and repeat/rebuild controls. |
| VAL-034 | 8 | ON TRACK | O1/I1 Receipt, Review and Human Outcome target/Event repositories come from their immutable subjects despite later edits or Git absence; evaluate remains distinct from stored reads. |
| VAL-035 | 7, 8; impact 9 | ON TRACK for current allocation | I1 Run/Gate derived fake/fid pairs, shared amended leaves, dual-adapter recovery and current changed-leaf assembly pass. D1 still owns whole-task-9 acceptance. |
| VAL-052 | 7, 8; impact 9 | ON TRACK | Final exact Developer gate receipts and required package timings are complete and below budgets/warnings. Focused reviewer times are separate evidence. |
| VAL-060 | 7, 8 | ON TRACK | O1 standalone examples run outside a project with exit 0, empty stderr and unchanged directories; conflict controls remain. |
| VAL-061 | 7 | ON TRACK | I1 both close/abandon leaves pair expectation mismatch with unmerged/non-Git recovery-state-conflict, retaining no durable effects. |
| VAL-062 | 7 | ON TRACK | I1 real close-proof unchanged versus tracked-edit fixtures retain succeeded versus failed outcomes against observed freshness. |
| VAL-066 | 7, 8 | ON TRACK | I1 unavailable-worktree/deleted-Source retries replay original results without effects; viable new-key controls refuse for the expected repository/Source cause. |

VAL-026-028 remain Window B context, with the VAL-027-029 message changes allocated to task 10. Later grammar, Human lines, Stop/repository views, judgment gets, reasoning, take-control journey and Human judgments retain their prescribed later windows.

## Gate And Evidence Checks

### Reports against exact required commands

Compared each report's `## Tests run`, `## Exact quality-gate results` and deviations against its task specification and QA-plan section 5. Both repairs retain all five byte-identical required commands in a final ordered PASS block, with output tokens and recorded observed termination. No missing, silently skipped or narrower substituted final required gate was found.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
```

| Report | Final exact gate evidence inspected | Failed-evidence handling |
|---|---|---|
| `task-7-done.md:208` | Format silent success, lint `0 issues.`, ordinary workflow `109.901s`, total coverage `61.8%`, tagged `323.135s`; each terminated exit 0 | Original failed lint/tagged/report attempts and their causal deviations remain. Recovery REDs have `FAIL`/exit 1, then byte-identical green commands, followed by complete final gates. |
| `task-8-done.md:298` | Format silent success, lint `0 issues.`, ordinary workflow `103.339s`, total coverage `61.8%`, tagged `347.525s`; each terminated exit 0 | Original implementation/fixture failures and recovery RED remain. Partial cached gate and explicitly unrun tagged step at line 288 remain distinct from the later complete fresh sequence. Cache preparation is not evidence. |
| `task-9-done.md:190`, impact only | Five final PASS receipts plus task-spec-declared R5 absence, expected/observed exit 1 | Historical ordinary/R5/report failures and authorized shared extension remain. Fresh report-lint COMPLETE and R5 absence reproduced; full slice remains D1. |

An independent inline Python audit compared the pre-repair report blobs at `c690cb93` with current files. All **20 task-7** and **40 task-8** original failed/blocked receipt lines remain byte-for-byte. The same audit verified exact allowlists and regular-file modes. Final token `PASS original failed evidence and bounded repair scope`, observed `exit_code: 0`. Report diffs add repair authority, red/green receipts, current exact gates and deviations without rewriting those failures. Corrected task-8 evidence descriptions replace the insufficient independence claims while preserving original failure/gate history.

| Task | cmd/fid | internal/workflow | Slowest other ordinary | Tagged integration |
|---|---:|---:|---|---:|
| 7 repair | 412.793 s | 109.901 s | internal/issues 204.714 s | 323.135 s |
| 8 repair | 372.613 s | 103.339 s | internal/export 199.131 s | 347.525 s |
| 9 close, impact context | 410.123 s | 102.835 s | internal/issues 199.339 s | 278.829 s |

These are inspected Developer exact-gate receipts, not newly attested QA executions. No ordinary package reaches 600 s or the 480 s warning; tagged results are below 900 s and the 720 s warning. Full repository lint/coverage/tagged gates, benchmarks, scale controls and full-checkptr were not independently rerun in this bounded re-review; focused fresh tests below reproduce the affected behavior.

### Fresh reproduced commands

Working directory for every command: `/home/mark/projects/fid`. Every yielded nested test session was retained and polled until its own `exit_code` arrived. Complete nested result objects were preserved; silence or wrapper completion was never treated as termination. All test selections use `-count=1`.

```text
O1:
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid ./internal/cli ./internal/evidence ./internal/workflow -run '^(TestP24Run.*|TestP24AbandonRecordsObservedBase|TestP24BudgetObservesOmittedSourceDigest|TestP24RemovedRunRepositoryMembers|TestP24Gate.*|TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation|TestP24RemovedGateEvidenceMembers|TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestAdoptedExampleProcessExclusivity|TestP24WholeTreeDescriptorProcessCorpus|TestP24DescriptorVariantProcessCorpus|TestInput.*|TestValidateGateRunReferencesAuthorizesAgainstEveryEnvelopeGrant)$'

I1:
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid/integration -run '^(TestP24Run.*|TestP22CloseRunUsesObservedRepositoryAndQualifyingProof|TestP22DualAdapterRecoveryConformance|TestP24Gate.*|TestP24ReceiptFreshnessRealGuardedActions|TestP24ReviewAndHumanEvidenceSubjectRepositories|TestP24ChangedFakeFIDAssembly|TestP23SharedAmendedLeafConformance|TestP23FreshnessAndOutcomeSupersessionRealProcessJourney)$'

O2:
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid -run '^(TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestEvidenceEvaluationMultiFailureGoldenSurvivesRepeatAndRebuild)$'
```

| Fresh check | Result | Token and observed termination |
|---|---|---|
| Representative exact required command: `gofmt -l cmd internal` | PASS | `exit 0 (no output)`; terminated `exit_code: 0` |
| O1 affected ordinary/shared impact | PASS | cmd/fid `59.006s`, internal/cli `1.039s`, internal/evidence `1.051s`; final token `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t4.228s`; terminated `exit_code: 0` |
| I1 affected real-process verification | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t58.215s`; terminated `exit_code: 0` |
| O2 ordinary freshness and repeat/rebuild | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t13.646s`; terminated `exit_code: 0` |
| `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-7-done.md project-docs/docs/phases/phase-24/tasks/task-7.md` | PASS | `OUTCOME COMPLETE`; final `PASS report-lint completed with no failures`; terminated `exit_code: 0` |
| `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-8-done.md project-docs/docs/phases/phase-24/tasks/task-8.md` | PASS | `OUTCOME COMPLETE`; same final token; terminated `exit_code: 0` |
| `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-9-done.md project-docs/docs/phases/phase-24/tasks/task-9.md` | PASS, impact context | `OUTCOME COMPLETE`; same final token; terminated `exit_code: 0` |

`\t` denotes the captured package-output tab separator. No fresh reviewer gate failed.

## Task 7 - Run leaves

### Must-fix

- None. Review-5 F1 is RESOLVED; see Prior Finding Resolution.

### Improvement

- None.

Observation (design sections 2-3, 6-7): `cmd/fid/p24_run_test.go:111` and `cmd/fid/integration/p24_run_inputs_test.go:118` require the complete dirty fixture to be exactly one unstaged tracked edit outside `.fid`. They read HEAD/tree/root through Git, project ID through config, HEAD/index modes and blob bytes through `ls-tree`/`ls-files -s`/`cat-file`, and actual worktree mode/content through `Lstat`/`ReadFile`. Stdlib JSON sorts keys of the accepted `fid.repository-dirty.v1` image; stdlib SHA-256 derives the file digests, dirty digest and root identity. No expected field calls `project.ObserveRepository`, `ObserveRecordedRepository`, `testEnvelopeRepository`, a fid repository result or the production canonicalizer.

Full map equality checks dirty start (`p24_run_test.go:42`), dirty close (`:103`) and both real-process terminal leaves (`integration/p24_run_inputs_test.go:82`). Ordinary start/close compare actual journal Events; the terminal matrix additionally finds the exact response-identified Event and compares it independently at lines 86-89. These returned qualifying Run records are loaded from the committed SQLite rows: start/close through `internal/workflow/run.go:205`, abandon through `run_close.go:200`, and `run_read.go:66` reads and decodes both stored repository columns. Thus the assertions cover stored records as well as actual Events, rather than an unpersisted request echo.

O1/I1 reproduce these comparisons along with match/HEAD/dirty/unsafe controls, Source digest pairs, unavailable/deleted-source replay, genuine import versus copy and Receipt-backed closure. Business rules remain in workflow/project domains. CLI/flow documentation matches observation, replay-before-read, null Budget and unsafe-terminal behavior. This repair strengthens evidence without changing accepted product behavior.

## Task 8 - Gate run and Evidence leaves; real freshness

### Must-fix

- None. Review-5 F2 is RESOLVED; see Prior Finding Resolution.

### Improvement

- None.

Observation (design sections 2-3, 5-7): `cmd/fid/p24_gate_evidence_test.go:116` independently derives the same eight-field dirty image from Git/config/root/file bytes before execution. Its expectation uses only independently obtained HEAD/dirty. The complete returned Receipt and both actual gate Event repositories equal this map. `p24AssertGatePersistedRepositories` at line 177 queries SQLite directly, requires exactly one Attempt and one Receipt, binds each ID to the terminal frame and compares the complete stored repository JSON. Event inspection separately requires one `gate.execution-started`, one `gate.receipt-recorded`, one Attempt Effect and one Receipt Effect, comparing both Effect repository payloads. Missing rows/Events/Effects cannot pass vacuously.

The retained `/bin/touch spawned.txt` gate changes the tree after capture. Fixture ignores are `.fid/` and `*.json` (`envelope_test.go:1217`); `spawned.txt` is an actual additional untracked image, so retaining the prior digest tests pre-spawn timing. Both HEAD/dirty mismatch controls at line 204 retain unchanged DB/journal and an absent process marker. O1/I1 reproduce them and qualifying execution; O2/I1 reproduce ordinary and real-process fresh/stale/reverted proof. Receipt/Review/Human Outcome target/Event comparisons retain subject identity despite current-state edits or missing Git. Replay precedes observation, and terminal execution reuses its pre-spawn observation. CLI/flow docs match these rules; no live-frame, terminal-reobservation or dependency change occurred.

## Task 9 - Current-HEAD impact only

### Must-fix

- None within this impact scope.

### Improvement

- None.

Observation (design sections 2-3, 6-7): DEC-284's shared typed-array descriptor extension and later component string/path correction remain present at current HEAD. `git diff --stat 697e1562..e39bfcc1 -- internal/workflow internal/evidence cmd/fid/run.go cmd/fid/run_recovery.go cmd/fid/gate.go cmd/fid/evidence.go cmd/fid/artifact.go go.mod go.sum` is empty, observed `exit_code: 0`. Since review-5, the only committed changes are the two bounded test/report repairs above. O1 exercises shared descriptor validation, registry/seeded mismatch, whole-tree process conformance and both selected array/object variant corpora alongside Window C leaves. I1 exercises current fake/fid assembly and retained shared conformance. They pass. Task 9's independently linted COMPLETE report retains its original failures and exact R5 declared absence. Neither this assessment nor its green shared-impact tests substitute for D1's full review of tasks 9/9a.

## Repair Plan

- None. Both prior bounded repairs are resolved by fresh independent evidence. There is no remaining repair owner or repair dispatch authorized by this PASS review.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| **Review-5 F1:** dirty Run expectations used the production observer; terminal dirty-digest/Event comparisons were not fully independent. VAL-022, VAL-030, VAL-031. | task-7 | **RESOLVED** | Exact three-path repair at `cdb7e01f`/`856c36d6`; independent ordinary and tagged Git/config/file-image oracles; all eight fields compared with committed start/close/abandon records and actual Events, including response-ID lookup for terminal Events; O1/I1 PASS with observed exit 0. Matching/refusing/replay/Source/proof states retained. All five fresh exact Developer gates present, original failed evidence preserved and independent report-lint COMPLETE/exit 0. |
| **Review-5 F2:** Gate pre-spawn expected dirty digest was production-derived; complete independent persisted equality/count evidence was missing. VAL-022, VAL-032. | task-8 | **RESOLVED** | Exact two-path repair at `a59fa4d1`/`95280952`; independent eight-field pre-spawn oracle; direct stored Attempt/Receipt equality, terminal IDs, both actual Events and both repository Effects with explicit one-each counts; post-spawn edit and mismatch no-effects retained; O1/I1/O2 PASS with observed exit 0. All five fresh exact Developer gates present, original failures/partial cached attempt preserved and independent report-lint COMPLETE/exit 0. |

The findings are preserved as distinct historical evidence defects, not erased or relabelled as prior product failures. This new independent artifact closes them; review-5 remains unchanged. Review-4's Window B verdict and historical deferrable Runner bookkeeping observations are not reopened.

## Cross-cutting observations

### Inventory and remaining feasibility

R1-R20 were freshly regenerated using the accepted inventory commands. Markdown-escaped alternations were rendered as regex alternation. R4's `ls migrations` and schema search were separate commands. R8 used `rg -n 'definition\(' internal/cli/codes.go`; R12 listed the three named renderer paths. All query/list commands returned observed `exit_code: 0` except R5, which returned no output and observed `exit_code: 1`, the task-9-declared absence condition.

The exact section-1 recursive Bash block was extracted from `input-inventory.md` and executed unchanged through an inline Python `subprocess.run(['bash'], input=block)` from repository root. It built a fresh scratch binary and wrote labelled lists under `/tmp/tmp.gKFxKAuOT5`. Last token `271 total`, observed `exit_code: 0`.

| Inventory | Fresh result | Assessment against review-5 and remaining allocation |
|---|---|---|
| R1/R2 | 45 call sites/33 files; one shared request-digest helper | Current production inventory retained; helper result names `internal/journal/request_digest.go:14`. |
| R3/R4 | 75 matches/16 files; 38 migration files through 0019; 17 schema references/5 files | Storage/schema roots remain executable and unchanged by repairs. |
| R5 | No output; observed exit 1 | Descriptor-only structural source remains; final task-9 absence reproduced. |
| R6/R7 | Four helpers/3 files; 16 resolved-fact references/2 files | SameWorktree and observation remain; private resolved Registration/Envelope facts are not revived caller inputs. |
| R8 | 84 lines: 83 code definitions plus helper declaration | Operation-unavailable remains task 9a. |
| R9/R10 | 42 functions/32 files; 33 assertions/14 files | Task 10 must account for writeExampleConflict and all current family writers. |
| R11/R12 | 19 matches/11 files; all three Stop/Attention renderer files exist | R11 also includes guidance/JSON/help/fixtures; task 11 must classify the set. |
| R13/R14 | 15 matches/3 files; one live-leaf completeness test | Judgment retirement and three new read leaves remain later tasks. |
| R15-R18 and input list | 92 leaves, 29 groups, 39 input leaves, 50 keyed leaves, zero missing Short, 61 typed flags | 49 canonical keyed leaves after export import; typed flags are 17 int +44 bool including 39 examples. Counts match review-5. |
| R19/R20 | Seven request sites/6 files; 25 invalid-argument sites/20 files | Still three grammar +four family request fallbacks; R20 remains 18 generic +three grammar +four specific. Future four-match R19 and complete fallback coverage remain tasks 9a/10. |

All counts retain review-5's current-tree pattern; the bounded repair diff introduces no new leaf, flag, failure producer, fallback, fixture-generator requirement or command definition. The prior review's full tasks-9a-14 corpus-feasibility assessment and 1,299-call probe remain cited historical feasibility, not freshly re-attested PASS results or future acceptance. Their product and specification inputs have not changed since that review. Fresh O1/I1 additionally reproduce the changed-tree descriptor/Run/Gate/freshness foundations, and the fresh inventories close the remaining corpus sets. Unrelated windows and the full future-corpus probe were not repeated.

Both exact task-13/14 absence queries remain executable: the disclaimer query returns three matches, last token `cmd/fid/p23_human.go:176:` followed by the Outcome disclaimer, observed exit 0; `rg -n "PRIVATE_" cmd/fid/integration` returns 12 matches, last token `cmd/fid/integration/testdata/p23_conformance/operations.json:145:      "deviation.observed": "PRIVATE_OBSERVED"`, observed exit 0. These are present-day feasibility observations, not passing absence gates. Their accepted final pass condition remains exit 1 after task 13 and again at task 14.

Fresh `rg -n '^func TestP24HumanFixtures\(' cmd/fid/integration` has no output and observed exit 1. The task-11 generator still does not exist, as expected. Its selector was not run and cannot be passed through a no-tests success. D2 must run the exact prescribed command at task-11 close:

```text
FID_P24_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-24-human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration
```

### Cadence

DEC-284 explicitly preserves the late Window C after task 9; review-5's REPAIR could not advance it. Both owners have now completed bounded repair and this linked fresh review supplies its independent PASS. The current brief/progress snapshots still describe the earlier task-9 recovery, whereas committed reports and verified history establish its completion. This remains contextual snapshot lag already disclosed by review-5, not a new protected boundary; no current document or ledger was changed here.

- **D1 after task 9a:** full tasks 9/9a review, including shared descriptor extension, selected content trees/redaction, whole-tree descriptors and complete grammar/root/group/typed/string-domain corpora with escaping and message boundaries.
- **D2 after task 11:** full tasks 10/11 review, generic fallbacks and Human/JSON pairing, Stop/Attention/repository views, and QA regeneration of `TestP24HumanFixtures` at task-11 close. Record separate Human judgment in `qa/human-walkthrough-after-task-11.md`; this does not replace phase-end HUMAN-02/HUMAN-03.
- **E after task 14:** tasks 12-14, new gets, judgment reasoning/disclaimers/control escaping, take-control provenance, semantic rebuild/export/import and retirement. All prescribed phase-end deterministic and Human controls remain required.

## Blocker Candidates

- None. No protected product, architecture, data-model, public-interface, dependency, security, gate/control or external-capability authority is required by this review.

## Existing Escalation References

- ESC-24-001 and ESC-24-002: resolved Task 6 history under DEC-283.
- ESC-24-003: resolved Task 9 orchestration history under DEC-284; no verifier change.
- None open. No escalation identity or register was created or edited.

## Final Lines

Review written: project-docs/docs/phases/phase-24/qa/review-6.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: PASS
