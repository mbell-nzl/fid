# Phase 23 - Independent QA Lead Checkpoint After Task 11 - Re-review 3

Reviewer: Codex QA Lead, fresh independent session; no involvement in review-7,
its repairs, the delegated performance repair, or the task reports.
Date: 2026-09-27
Review type: FRESH INDEPENDENT RE-REVIEW
Re-review of: [checkpoint-after-task-11-2.md](checkpoint-after-task-11-2.md)
(Runner review-7, `9cb9420516c461ddd9b4d22d1944759a416f7bb7`).
Authority: DEC-255, ordered by DEC-256, with DEC-258/259.
Reviewed HEAD: `92ea1e509198466b9c0a497b82f80710ca2a8d4d` on `phase-23`.

## Verdict

PASS

Review-7's Must-fix findings are resolved by the named owners and independently
verified below. Its three protected questions have accepted decisions and
matching implementation/evidence. There are **0 Must-fix findings, 1 Improvement,
and 0 unresolved blocker candidates**. ESC-23-004 and ESC-23-005 have recorded
resolutions; this review creates no escalation identity.

This closes the DEC-255 checkpoint review only. It does not execute a Runner
transition, dispatch task 12, or accept the remaining phase work.

## Scope Reviewed

This is the checkpoint required by QA plan section 8, including the retained
review-7 findings and the prospective command feasibility check for tasks
12-18. It is not phase-end acceptance. Task 12 remains held throughout this
review; this session neither dispatches nor implements it.

Context read: `AGENTS.md`, `CLAUDE.md`, the four current-context files,
`08_QUALITY_PLAN.md` (especially sections 5-8), `09_AI_WORKFLOW.md` sections
2.6, 6.1 and 8, and the mid-phase review template. Applied the
`validation-contract` and `qa-evidence` skills.

Phase sources: validation contract; QA plan sections 4-8, 20-22 and the
referenced readiness amendments; review-7 including its Repair Plan; task
2/6/10/11 specifications and completion reports; predecessor reports and
task 4/5/9 evidence for the retained checkpoint areas; tasks 12-18 command,
scope and fixture requirements; contracts sections 2-7 and 9-11; public
interface sections 4/4.1; execution-scope section 6; contract-coverage mapping;
the corresponding PRD and glossary entries; DEC-249 through DEC-259, with
DEC-227/235/236/239/240/242/243 retained as applicable authority.

Read [esc-23-005-repair.md](esc-23-005-repair.md), all three committed gate
logs in `esc-23-005-evidence/`, the escalation dispositions, and the Runner
ledger as evidence only. No ledger, original review, contract, task, decision,
current-context file, product source, or test was edited. Scratch CLI work uses
`/tmp/fid-qa23-11-3.vzUzxV`; no held-out scenario material was read.

## Summary

- Outcome: PASS; Must-fix 0, Improvement 1, unresolved blocker candidates 0.
- Review-7's BC-1/2/3, MF-4/5 and I-5/6 are resolved, with retained earlier
  findings rechecked. I-7 is the sole new, deferrable report-wording finding.
- Contract accounting: 31 checkpoint PASS outcomes and 53 CARRY-FORWARD
  outcomes, each assigned below. This is not phase-end acceptance.
- Existing escalation references: ESC-23-004 and ESC-23-005, both resolved.
- Task 12 remains undispatched by this session.

## Commit And Report Verification

Hashes and subjects below were resolved from Git in this session. The repair
order is task 6, task 2, task 10, task 11. The intervening report commits are
ancestors of the reviewed HEAD.

| Owner | Implementation commit and subject | Committed report |
|---|---|---|
| task-6 | `8531522e9602b0f72717d6c9d889684228cfa2ee` — `fix(evidence): authorize gate.run against every Envelope grant` | `d6542bfc093059c8c6e656fa25c26229e60408ee` |
| task-2 | `613e45aad9b10445cffd6b113ee598ab493f183a` — `fix(workflow): evaluate only guarding Obligations per guarded action` | `7077776135370ac121c658d80880c983cf365b77` |
| task-10 | `a425ffe22ec12558e953bbbb5cb265d22f28779e` — `test(integration): prove VAL-009 gate-satisfied review-unsatisfied stage under DEC-249` | `bc7e6f393450780137f3d89ad9763e382028ccd3` |
| task-11 | `466745f2dc4505d60eb67ac50193e004532b7239` — `fix(workflow): resolve judgment Stops by recorded Condition through one shared predicate (DEC-250)` | Final delegated closure: `21f60f2e691b7ddc1a9f2440d23883c1a4f0eb8c` |
| DEC-258 | `3dff4574a71438900100449e36d7044969ba9162` — `perf(test): parallelize serial test suites for ESC-23-005` | Cited by task 11's final report |
| DEC-259 | `92ea1e509198466b9c0a497b82f80710ca2a8d4d` — `build(test): build modernc SQLite without checkptr in race gates` | Prospective command amendment |

Scope inspection used `git diff --numstat` over `7a0243e..d6542bf`,
`d6542bf..7077776`, and `7077776..bc7e6f3`, and `git show 466745f`.
Task 6 changes four authorized source/test files plus its report; task 2
changes its two evaluators, matching tests and fixture convergence plus its
report; task 10 changes only its journey and report. Task 11's shared
predicates and producer call sites are within DEC-250/253's ratified scope.
No closed journal registry was silently widened by that repair.

DEC-258 changes 215 Go files, all test files. Comparing each parent/child file
after removing `t.Parallel()` lines, comments and inline `paralleltest`
annotations leaves zero differences. The three initially nonidentical files
were inspected individually: `live_output_test.go`, `maintenance_test.go` and
`id_test.go` contain the documented serial exceptions. Assertions are retained.
`git diff 3dff457..HEAD -- cmd internal migrations go.mod go.sum` is empty:
the delegated closure and this review use identical Go sources and dependencies.

### Exact-command audit

The original DEC-227 command text is retained in tasks 2/6/10/11. Their latest
repair sections contain all required commands: respectively 6, 7, 5 and 5.
These include the full ordinary race/coverage command, coverage report, full
tagged integration command, and task 2/6's dedicated package coverage commands.
No focused selection substitutes for a required command. The original task-10
`-count=1` deviation is labelled historical; its repair rows use the exact
specification commands.

| Owner | Latest reported ordinary `cmd/fid` / `internal/workflow`; tagged integration | Additional evidence checked |
|---|---|---|
| task-6 | 478.295 / 185.338; 681.447 s | Workflow's repository-wide result was cached; the dedicated exact run supplies its fresh time. Evidence coverage 29.2%, workflow 62.1%. DEC-246's explicitly reconstructed pre-change baselines remain labelled as such. |
| task-2 | 490.783 / 209.816; 691.913 s | Dedicated workflow after-run 188.151 s, 64.0%; pre-edit baseline and supplementary uncached measurement are recorded separately. |
| task-10 | 468.708 / 192.775; 689.526 s | Workflow's exact full-suite result was cached; the separately disclosed uncached diagnostic supplies 192.775 s. The actual required full command is present. |
| task-11 | 531.533 / 165.920; 326.421 s | Delegated DEC-258 evidence on the original DEC-227 commands, as detailed below. |

These are audited report receipts, not executions made by this QA session.
The ledger records Runner report-lint `OUTCOME COMPLETE` and termination
`exit_code: 0` for all four owners. Task 11's ledger explicitly says the Runner
did not rerun the long gates. Fresh QA command results are recorded separately.

## Prior Finding Resolution

The original review remains unchanged. Closure here depends on the accepted
decisions, inspected implementation and tests, and this session's executions.

| Prior finding | Owner and authority | Disposition and independent evidence route |
|---|---|---|
| BC-1 — applicable Obligations, VAL-009 | task-2; DEC-249 | `obligationsGuarding` filters by `required` and the requested action. Terminal success evaluates that subset; completion filters denial reasons while retaining every evaluation in reads. The five-case domain matrix and two built-binary tests cover a satisfied gate with an unsatisfied review guarding both actions, the contrast where it guards only completion, and a no-guard state. The contrast actually commits `qa_pass`, then confirms completion remains blocked. |
| BC-2 — Human Outcome grant contradiction, VAL-043 | Human DEC-251 Option B; no implementation dispatch | Public-interface section 4.1 now agrees with contracts section 10: finalized human profile, capability and activation provenance, with no Envelope grant. The decision journey records real Outcomes as `actor-human-outcome`, which has no Envelope grant. The Human Outcome matrix still executes all five authority cases, bracketing rejected requests with Event counts. |
| BC-3 — judgment Condition qualification, VAL-048/050/078 | task-11; DEC-250/253 | `qualifyingJudgmentPredicate` is called by both manual resolution and the producers' `QualifyingJudgmentStops` scan. The scan selects recorded Condition types, not raising codes. The two new built-binary journeys cover producer and operator Stops, both Decision-pair orders, the earlier-half rejection, wrong class/result/type, supersession, and run/envelope/activation scoping. Activation-subject `decision-satisfied` rejects at raise with no Event. |
| MF-4 — missing linked-recovery evidence, VAL-027 | task-6; review-7 Repair Plan, DEC-256 | `TestP23ReviewCycleChargeSurvivesLinkedRecoveryAndRebuild` records one Review-keyed sample with value 1, recovers into a distinct successor linked to the prior Run, and compares sample count, consumed sum and every sampled field before/after recovery and rebuild. Its test-only repair needed no budget-product change. Import remains task 16. |
| MF-5 — missing mixed state, VAL-009 | task-10 after task-2; DEC-249/255/256 | `TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess` stage 2 now links the real gate receipt before the review. The read requires exactly review and security Human Obligations unsatisfied; completion names both and changes neither Event count nor the full projection. The non-committing evaluator proves no terminal blocker; the single real terminal transition remains at the all-satisfied end as DEC-255 requires. |
| I-5 — multi-grant gate authorization | task-6; DEC-252/256 | `authorizeGateRunGrants` examines every grant at the bound revision. The real-process test succeeds for first- and last-sorted granted Actors; an Actor with Claim/Run authority but no `gate.run` grant is denied without Event or gate-projection change. The five-Actor domain matrix also rejects wrong Issue and missing required capability. The existing failure code/message is retained. |
| I-6 — report citations and evidence descriptions | task-10; folded into MF-5 | The dated repair section and corrected VAL-009/010/025/043-046/080/081 entries resolve to the named assertions at `a425ffe`; later parallelism adds only line offsets. Historical false claims are explicitly superseded. The report distinguishes terminal reason-code multisets, domain evaluation and real CLI transitions, and labels its original `-count=1` history. |
| I-1 — time headroom | DEC-254, superseded for its speed step by DEC-258/259 | DEC-258 supplies the authorized parallelism repair; DEC-259 supplies the separately approved compiler flag. Fresh measurements below check the remaining headroom. DEC-254 Part 1 attribution remains binding; task 13 must not repeat the already-satisfied Part 2 step. |
| I-2 — `stop_conditions.go` scope | task-11; DEC-253 | The decision and task amendment explicitly ratify the file, its tests and the producer change in `decision.go`; full gates cover them. |
| Earlier MF-1/2/3 and I-4 | task-10; retained review-7 rows | Rechecked coexistence and all-satisfied completion; actual authority requests and consumed before/after Event brackets; real Human Outcome positive/mismatch links; and `not-applicable` with linked review/Human evidence. The repaired journeys and classifier matrix execute again below. |
| Earlier I-3 | task-11; superseded by BC-3 | Its previously missing manual qualification route is now covered by the BC-3 implementation and executions. |

### Two evidence distinctions

MF-5's preference for both Human Obligations already satisfied was not a
mandatory isolation requirement. Keeping security unsatisfied preserves the
existing readiness-Stop raise/resolve proof. The new stage explicitly requires
the review ID in the denial. Task 2 separately isolates the review as the sole
blocker and drives the actual terminal transition. This review does not describe
the task-10 stage-2 domain call as a real-process transition or claim to rerun
the Developer's reverted mutation experiments.

VAL-025 retains all nine cells. Durable receipt/gate and review/review links
come from task 6 and the mixed journey. The four non-Human mismatches use the
production classifier and check `evidence-type-not-accepted` with the accepted
list. Task 10 records and durably links the approved Human Outcome; its two
wrong-type requests check `outcome-not-permitted`, exact Obligation type, and
unchanged journal/projection. Seeded classifier acceptance is not used as the
Human positive-link evidence.

## Delegated ESC-23-005 Closure And DEC-259

DEC-258 is the Human-delegated authority for the test-only repair beyond task
11's ordinary file list. Its committed logs bind `3dff457` and record the exact
five DEC-227 task commands with observed termination:

| Command | Delegated closing result in `exact-gates-attempt-2.log` |
|---|---|
| `gofmt -l cmd internal` | Silent success; recorded `exit_code: 0` (line 4). |
| `golangci-lint run ./...` | `0 issues.`; recorded `exit_code: 0` (lines 6-7). |
| `go test -race -timeout 30m -coverprofile=coverage.out ./...` | `cmd/fid` 531.533 s, `internal/workflow` 165.920 s; recorded `exit_code: 0` (lines 9, 30-31). |
| `go tool cover -func=coverage.out` | `total:` 62.4%; recorded `exit_code: 0` (lines 1838-1839). |
| `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` | `ok` integration 326.421 s; recorded `exit_code: 0` (lines 1841-1842). |

These are delegated execution receipts independently inspected here. They are
not fresh Runner or QA executions of task 11's closing gates. The reported
shuffle/count-two flake guards have no retained raw logs; their evidence class
remains delegated attestation, and they are not re-labelled as this session's
verified gate results.

The following failed attempts remain failed:

- Task-11 ordinary gates: `cmd/fid` 686.386 s and 630.927 s, each process
  `exit_code: 0` but **FAIL against the 600 s budget**. Their workflow times
  were 277.603/252.947 s; tagged times 778.278/766.070 s. The attribution
  parent's 614.601 s also breached the ordinary budget.
- DEC-258 attempt 1: ordinary `FAIL .../cmd/fid 525.379s`, recorded
  `exit_code: 1`; 14 `journal files = 2, want 1` failures. The successful
  335.029 s tagged result does not turn that ordinary attempt into PASS.
- The repair's earlier maintenance-grace diagnostic failed and led to the
  documented serial exception; no assertion was changed.

The midnight assumption predates this phase: `git log -S 'journal files = '
-- cmd/fid/init_test.go` identifies `bdd393a61b64581d6dc41268b5d875802eab6f34`
(2026-06-06, `Refs: phase-1/task-5`), before Phase 23 planning on 2026-09-14.
The existing disposition is CF-044. I-5's single-grant defect likewise traces
to `b694991` (2026-08-30, phase-21/task-15), and has now been repaired.

DEC-259 subsequently changes the two exact race-gate command strings for
**tasks 12-18**, adding `-gcflags=modernc.org/...=-d=checkptr=0` while retaining
race detection, package scope, assertions, 600/900 s budgets and `-timeout 30m`.
The committed `dec-259-gates.log` records unchanged Go sources at `21f60f2`:
`cmd/fid` 299.562 s, workflow 96.706 s, integration 239.187 s, both commands
`exit_code: 0`. Those figures do not retroactively replace task 11's receipts.

QA plan section 22 also requires `go test -race -timeout 30m ./...` without
the flag at phase end, and when a task changes a modernc version. That control
must PASS and has no VAL-079 timing budget. It is not due at this checkpoint;
no modernc version changed. The phase-end control remains outstanding.

## Gate And Evidence Checks

All commands below ran serially from `/home/mark/projects/fid`, except
the explicitly named disposable CLI directory. Each termination was observed
in the nested command result; every live session was polled until it returned
`exit_code`. No discarded or merely silent session is treated as completed.
The test cache was cleared once before the gates, as setup. No gate was retried
to obtain a passing sample.

Reference instance: Go 1.26.2, linux/amd64; Linux
6.6.114.1-microsoft-standard-WSL2 x86_64; 6 available CPUs; AMD Ryzen 7 5700G;
`MemTotal: 8133236 kB`. Repo-local Go caches were used and `GOFLAGS` was empty.
The initial `/tmp` check showed 2.7G free on its 3.9G tmpfs.

### Fresh deterministic results

These are this review's executions, distinct from the delegated closing
receipts above. G3/G5 use the current section-22 command text for the required
forward-feasibility check; they do not rewrite earlier task requirements.
G6/G7 retain full checkptr and add independently named behavioral evidence.

| ID | Exact command | Result and observed token |
|---|---|---|
| G1 | `gofmt -l cmd internal` | PASS — `exit 0 (no output)` (terminated `exit_code: 0`). |
| G2 | `golangci-lint run ./...` | PASS — `0 issues.` (terminated `exit_code: 0`). |
| G3 | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` | PASS — final workflow package tokens `104.092s` and `coverage: 64.6% of statements`; 21 packages `ok`, no cached results or FAIL lines; `cmd/fid` 300.650 s (terminated `exit_code: 0`). |
| G4 | `go tool cover -func=coverage.out` | PASS — total coverage token `62.4%` (terminated `exit_code: 0`). |
| G5 | `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` | PASS — integration package `ok`, token `235.103s` (terminated `exit_code: 0`). |
| G6 | Complete named CLI command below | PASS — integration package `ok`, token `23.647s`; all 17 named top-level tests and nine subtests PASS (terminated `exit_code: 0`). |
| G7 | Complete named domain/journal command below | PASS — review 1.164 s, evidence 1.045 s, workflow 7.115 s; final journal package `ok`, token `2.467s`; all 23 named tests PASS (terminated `exit_code: 0`). |
| G8 | `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` | PASS — `22860 ms/rebuild`, below 30000; journal package `ok`, token `25.597s` (terminated `exit_code: 0`). One checkpoint sample, not the three phase-end samples. |
| G9 | `go build -o /tmp/fid ./cmd/fid` | PASS — `exit 0 (no output)` (terminated `exit_code: 0`). |

G3/G5 timing margins are 299.350 s and 664.897 s against the unchanged 600 s
and 900 s budgets. Their 50.1% and 26.1% utilization are below DEC-258's 80%
warning line. These are single-session measurements, not a claim that future
test additions cannot exhaust the budgets.

G6, copied in full so the selection is reproducible:

```bash
go test -tags integration -race -timeout 30m -count=1 -run '^(TestP23ReviewRecordRealProcessJourney|TestP23FindingDispositionRealProcessJourney|TestP23ReviewEvidenceRealProcessJourney|TestP23ReviewCyclesRealProcessJourney|TestP23ReviewCycleChargeSurvivesLinkedRecoveryAndRebuild|TestP23GateRunAuthorizesEveryEnvelopeGrantThroughBuiltBinary|TestP23ReviewObligationGuardingBothActionsBlocksBothRealProcess|TestP23UnsatisfiedObligationBlocksOnlyTheActionItGuardsRealProcess|TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess|TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation|TestP23HumanOutcomeRequiredRaisedByReviewRecord|TestP23DecisionReferenceProtectedLifecycleRealProcess|TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess|TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen|TestP23StopResolveHumanOutcomeConditionKeysOnConditionNotRaisingCode|TestP23StopResolveDecisionSatisfiedConditionKeysOnConditionNotRaisingCode|TestP23FreshnessAndOutcomeSupersessionRealProcessJourney)$' ./cmd/fid/integration/... -v
```

G7, copied in full:

```bash
go test -race -count=1 -run '^(TestCheckIndependenceRejectsClaimAndRunHistoryDirections|TestTerminalFindingDispositionAuthorityMatrix|TestChangesRequiredReviewChargesExactlyOneReviewCycle|TestValidateGateRunReferencesAuthorizesAgainstEveryEnvelopeGrant|TestEvidenceAdmissionTask6PureProductionClassifierMatrix|TestGuardedActionsEvaluateOnlyObligationsThatGuardThem|TestCompleteEnvelopeDeniedByGuardingObligationCommitsNothing|TestHumanOutcomeConditionQualifyingFactRules|TestDecisionSatisfiedConditionQualifyingFactRules|TestManualAndAtomicResolutionShareOnePredicate|TestQualifyingJudgmentStopsKeysOnConditionNotRaisingCode|TestJudgmentFactFromEventReadsTheCommittedEffectRow|TestEventMatchesJudgmentConditionTypes|TestStopRaiseRejectsDecisionSatisfiedOnActivationSubject|TestTask11ReviewRecordedAndFindingResolvedConditionShapes|TestTask11ConditionShapeRegistryIsClosedAtTenRows|TestTask11FindingResolvedConditionDispatchPreservesNullableBinding|TestTask10HumanOutcomeConditionShapeAndNonStringDiscriminant|TestTask10HumanOutcomeConditionDispatchPreservesNullableBinding|TestTask9ConditionDispatchPreservesNullableBindings|TestTask9ConditionDispatchRetainsPinnedPhase22Bytes|TestTask9PrechangeBudgetObservedJournalReplaysWithoutMutation|TestTask9RetainedAttentionKeyMatchesPrechangeReference)$' ./internal/review/ ./internal/evidence/ ./internal/workflow/ ./internal/journal/ -v
```

Named tests are supplementary evidence. Neither selection substitutes for G3,
G5, or an owning task's exact command. G6 builds and invokes the real binary;
its assertions inspect responses, SQL projection state and journal effects.
G7 supplies the isolated durable-Run-history test permitted by DEC-239 and the
structural Condition tests that cannot be proven by CLI help or code reading.

### Checkpoint checklist results

| Check `(qa-lead)` | Result | Fresh evidence |
|---|---|---|
| QA-10 / VAL-009: denied and permitted guarded actions, including the repaired mixed state | PASS | G6's two `ReviewObligation`/ `UnsatisfiedObligation` journeys plus the Human Outcome journey; G7's five-case guarding matrix and non-mutating completion denial. |
| QA-13 / VAL-015: independence in both durable-history directions | PASS | G6 `TestP23ReviewRecordRealProcessJourney`: exact conflicting Claim/Actor, no Event, same independent reviewer before/after rework; G7 `TestCheckIndependenceRejectsClaimAndRunHistoryDirections`: isolated Run-history query. |
| QA-15 / VAL-022: terminal Finding disposition authority | PASS | G6 Finding journey and G7 `TestTerminalFindingDispositionAuthorityMatrix`: six cases for each terminal action; wrong/creating/non-independent/blocked Reviews rejected, later qualifying independent Reviews accepted. |
| QA-16/21 / VAL-025: nine evidence-type cells | PASS | G6 Review Evidence and Human Outcome journeys plus G7 `TestEvidenceAdmissionTask6PureProductionClassifierMatrix`; durable Human positive and two typed non-mutating mismatches retained. |
| QA-17 / VAL-027: charge, replay, linked recovery and rebuild | PASS for the due legs; import CARRY-FORWARD | G6 cycles and linked-recovery journeys plus G7 `TestChangesRequiredReviewChargesExactlyOneReviewCycle`; the import leg is explicitly task 16. |
| QA-17 / VAL-029: atomic exhaustion and refused next rework | PASS | G6 `TestP23ReviewCyclesRealProcessJourney` verifies charge/Stop/Attention in the recording Event and both denied next-Run/Claim paths without new durable rows. |
| QA-20 / VAL-039/040/041: protected Decisions | PASS | G6 Decision lifecycle and Decision-condition journeys: affected-action isolation, missing/wrong/rejected halves, successful pairs in both orders, correct resolution Event, other Stops remain open. |
| QA-21 / VAL-043: real Human provenance | PASS | G6 Human Outcome authority matrix: four rejected provenance/capability/type states and valid finalized Human; Decision journey proves DEC-251's no-grant state. |
| QA-23/43 / VAL-048/049/050/078: qualified Attention/Stop resolution and retained input | PASS | G6's two Attention journeys and two recorded-Condition journeys, plus G7's shared predicate/fact/scope matrices. Extra input is rejected without Event; wrong facts keep Stops open; matched facts resolve only their target conditions. |
| QA-46/48 / VAL-080/081: closed discriminator shapes and nullable-binding relation | PASS | G7's task-9/10/11 journal matrices, G3's complete journal suite, and strict-append assertions. Relation failures are checked as `condition nullable binding mismatch`, separately from shape errors. |
| QA-47/50/51 / VAL-082/083/084: retained Phase 22 bytes, replay, semantic key | PASS | G7's three named pinned-byte/replay/key tests; all six retained type fixtures, unchanged source journal and fixed-input key equality. |
| QA-22 / VAL-047: capability-advertisement first half | PASS for false state; true CARRY-FORWARD | Fresh G9 binary in the disposable project returns `"human_outcome_resolution":false`; task 17 owns the true state and phase-end two-version comparison. |
| Section 8 / R-13: prospective feasibility | PASS | G1-G5, G8-G9 and the complete command/fixture audit below. |

### Report-lint control audit

This is a supplementary audit of the Runner's control, not a new Developer or
QA gate. The documented two-argument usage was read before execution. Each
command returned `OUTCOME COMPLETE` and ended with the verbatim token
`PASS report-lint completed with no failures`, observed termination
`exit_code: 0`:

```bash
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-6-done.md project-docs/docs/phases/phase-23/tasks/task-6.md
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-2-done.md project-docs/docs/phases/phase-23/tasks/task-2.md
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-10-done.md project-docs/docs/phases/phase-23/tasks/task-10.md
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-11-done.md project-docs/docs/phases/phase-23/tasks/task-11.md
```

The audits confirm exact command strings, resolvable implementation hashes,
matching subjects, committed clean reports, and typed completion. The separate
claim-by-claim inspection above establishes whether the evidence supports the
reports; lint alone does not establish behavior.

## Command And Feasibility Validation

All nine unique command strings prescribed by tasks 12-18 are concrete.
Their required-block command counts are 7, 5, 5, 9, 6, 5 and 5. Each has exactly
the two DEC-259 race forms. None needs an inferred flag, renamed path, narrower
package scope, dependency change or new protected decision.

| Command or step | Validation method | Result and evidence |
|---|---|---|
| Go flags, package scopes, coverage output | `go help testflag`, `go help build`; full G3-G7 execution | PASS: timeout/count/run/coverprofile and package-targeted gcflags are supported; coverage file is generated and readable. |
| Named checkpoint selections | Exact function declaration existence plus G6/G7 `-v` output | PASS: 17 CLI and 23 domain/journal functions actually ran; no empty selection. |
| Reference benchmark (tasks 12/15/16) | File existence at `internal/journal/rebuild_bench_test.go:19`; G8 | PASS: real 100k rebuild with semantic-parity assertion, 22860 ms. |
| Tagged harness and archive fixtures | File/directory checks and G5 | PASS: `cmd/fid/integration/main_test.go`, `testdata/`, and `p22_closed_values_test.go` exist; the harness builds disposable binaries and workspaces. Task 16 extends the accepted archive profile; current fixtures are not claimed as future Phase 23 round-trip closure. |
| Task-12 requirement-removal search | Exact command A below executed | Feasible inventory: `internal/journal/current_event_catalog.go:119:` is the final match; 8 matching lines, terminated `exit_code: 0`. Product matches remain as expected before task 12. |
| Task-15 removed-table search | Exact command B below executed | Feasible inventory: final match `internal/issues/detail.go:208:`; 9 lines, terminated `exit_code: 0`. All matches are in the owning task's convergence scope. |
| Task-15 removed-rule search | Exact command C below executed | Feasible inventory: final match `internal/workflow/rules.go:8:`; 2 lines, terminated `exit_code: 0`. |
| Expected-non-zero declarations | Read task 12/15 sections against A/B/C | PASS: eventual absence requires no output and observed `exit_code: 1`. The current `exit_code: 0` inventories are not absence PASS results. |
| Disposable CLI setup | `git init -q`, `/tmp/fid init --help`, `/tmp/fid init --prefix QAC` | PASS: each terminated `exit_code: 0`; initialization token `initialized fid project in .fid/`. Cwd `/tmp/fid-qa23-11-3.vzUzxV/cli.pMsdCD`. Explicit prefix avoids directory-name constraints. |
| Capability first half | `/tmp/fid --json capabilities` in that directory | PASS: `"capability_id":"attention.v1","limits":{"human_outcome_resolution":false,"state":true}` (terminated `exit_code: 0`). |
| Human read/act steps | Built-binary help and G6 actual invocations | PASS for feasibility: `review get`, `finding list`, `envelope get`, `resume`, `human-outcome record`, `issue transition`, `attention list`, `attention resolve`, `stop resolve` each has usable help, terminated `exit_code: 0`. Mutation leaves expose actor/idempotency/input or the declared transition flags. |
| Human fixture generation and provenance | Read task 14's scope and Phase 22 precedent; file existence | Feasible, due later: `p22_human_fixtures_test.go` and phase-22 `qa/human-fixtures.md` exist. Task 14 owns `p23_human_fixtures_test.go` and phase-23 `qa/human-fixtures.md`, which do not exist yet; task 18 regenerates them. Manifest binds binary digest, source commit, cwd and absolute argument vectors. |
| Human judgment | QA plan section 9 against task 14/18 fixture ownership | CARRY-FORWARD: HUMAN-01a/b/c, 02a/b, 03a/b and 04a/b/c remain Human questions. G6 proves mechanics; it cannot answer comprehension, obviousness or presentation judgments. |
| Phase-end controls | QA plan sections 5/22, existing source/fixture paths and Go help | Applicable later: full benchmark sweep, archive-scale gate, three consecutive reference samples, full-checkptr control, completed migration/capability/contract audits and Human judgments. No archive limiter or modernc version changed in these repairs. |

A, B and C, exactly as executed:

```bash
rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal
rg -n --glob '!*_test.go' "task_completions|ac_verifications" cmd internal
rg -n --glob '!*_test.go' "has_completion_report|concerns_addressed" cmd internal
```

Exploratory help errors are retained separately: the guessed
`/tmp/fid transition --help` returned no output, `exit_code: 2`;
`/tmp/fid attention get --help` returned the parent Attention help,
`exit_code: 0`, so it did not prove a get leaf. The root command tree and the
accepted journey calls identify `issue transition` and `attention list`;
their correct help commands above both exited 0. A source search that guessed
`cmd/fid/transition.go` likewise exited 2 because that path does not exist.
These are corrected exploratory probes, not failed accepted commands or
invented Human checklist steps. No test or product change followed them.

## Validation Contract Coverage

The file-state audit found exactly 84 unique IDs, with nonempty source,
evidence, validator, owner and pass/fail fields. The DEC-237/240 split remains
intact: resolution VAL-078; shape VAL-080; relation VAL-081; retained-byte
VAL-082; replay VAL-083; semantic-key VAL-084. No contract refinement is needed
for this checkpoint, and no contract file was edited.

Every ID has an outcome below. PASS is a fresh checkpoint conclusion with the
cited evidence. CARRY-FORWARD identifies a scheduled later leg or an assertion
outside this checkpoint's deep review; it does not waive phase-end validation
or reopen completed tasks. Full gate execution supplies regression coverage,
but does not replace a missing specialized phase-end proof.

| VAL ID | Owning task | Outcome | Evidence or remaining scope |
|---|---|---|---|
| VAL-001, VAL-002, VAL-005, VAL-075 | 1 | CARRY-FORWARD | Earlier schema window; G3 regression passes. Phase-end QA-03/04/07 inventory, exact hunk and seeded preservation proofs remain scheduled. |
| VAL-003 | 1/4/5/8/9/10/12 | CARRY-FORWARD | Progressive fence contributors through task 10 audited; task 12's family and final fence-removal provenance remain due. |
| VAL-004 | 1/15/16 | CARRY-FORWARD | Final 0018 and rebuild pragma states require tasks 15/16. |
| VAL-006, VAL-007 | 2 | CARRY-FORWARD | Task-2 exact-command/report audit and G3 pass; full authoring matrix is outside the required section-8 re-review and remains phase-end QA-08. |
| VAL-008 | 2/6 | PASS | G6 mixed-Obligation reads retain the complete ordered evaluation and reasons while G7 filters only the guarded action. |
| VAL-009 | 2/10 | PASS | BC-1/MF-5 closure; G6 paired denied/permitted states and G7 requested-action matrix. |
| VAL-010 | 2/10 | PASS | G6 Human Outcome journey's present gate/review/Human evidence and exact availability assertions; report citations repaired under I-6. |
| VAL-011, VAL-012, VAL-013 | 3 | CARRY-FORWARD | Earlier Claim window; G3/G5 regression passes. Normalized two-mode lifecycle parity remains phase-end QA-11. |
| VAL-014, VAL-016, VAL-017, VAL-018, VAL-019 | 4 | CARRY-FORWARD | G6 Review journey passes; broader immutable-record, binding, verdict and complete replay checks remain phase-end QA-12/14. |
| VAL-015 | 4 | PASS | G6 real Claim conflict/same-reviewer rework and G7 isolated Run-history query, as DEC-239 specifies. |
| VAL-020, VAL-021, VAL-023 | 5/6 | CARRY-FORWARD | G6 Finding and Review Evidence journeys pass; full original creation/address/state matrix remains phase-end QA-15/16. |
| VAL-022 | 5 | PASS | G6 Finding journey; G7 twelve terminal authority cases. |
| VAL-024 | 6 | PASS | G6 Review Evidence journey's three verdict link states. |
| VAL-025 | 6/10 | PASS | G6 durable gate/review/Human positives and Human negatives; G7 production mismatch classifications, nine cells reconciled above. |
| VAL-026 | 6 | PASS | G6 Review Evidence journey retains an unresolved blocker after a later pass Review is linked. |
| VAL-027 | 6/16 | CARRY-FORWARD | MF-4's recovery gap is RESOLVED by G6/G7; record/replay/recovery/rebuild pass. Import leg explicitly awaits task 16. |
| VAL-028, VAL-029 | 6 | PASS | G6 cycle journey: fid-owned observation rejects; atomic exhaustion and next-attempt denials preserve journal/entity counts. |
| VAL-030 | 6 | PASS | G3 executes `internal/workflow/envelope_test.go`'s `TestEnvelopeAuthoringReviewCycleBudgetMatrix` and `TestEnvelopeAuthoringReviewCycleAlongsideRetainedBudgets`: both mismatches, missing Review Obligation, accepted pair and retained ordinary budgets. |
| VAL-031, VAL-032, VAL-033 | 7/10 | PASS | G6 freshness/supersession journey: changed tree, retained Review row, Human freshness, both supersession directions with records retained. |
| VAL-034, VAL-035, VAL-036, VAL-037, VAL-074 | 8 | CARRY-FORWARD | Earlier Concern/Deviation window; G3/G5 regression passes. Full QA-19 matrix remains phase-end scope. |
| VAL-038, VAL-042 | 9 | CARRY-FORWARD | G6 Decision lifecycle passes; full reference storage/source/transition-rule audit remains phase-end QA-20. |
| VAL-039, VAL-040, VAL-041 | 9 | PASS | G6 protected lifecycle and repaired Decision-condition journey, with G7 wrong/earlier-half/absent-declaration negatives. |
| VAL-043 | 10 | PASS | BC-2 closure plus G6 finalized-human provenance matrix and actual no-grant Human recording. |
| VAL-044, VAL-045 | 10 | PASS | G6 Human Outcome mixed journey and Attention journey cover matching result/type and judgment-only Evidence. |
| VAL-046 | 10 | PASS | G6 covers activation, evidence-link and Review readiness producers plus deduplication; G3 executes `cmd/fid/artifact_test.go`'s `TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation` for the fourth, `artifact.availability-observed`, producer. |
| VAL-047 | 11/17 | CARRY-FORWARD | Fresh capability probe proves false at this checkpoint; task 17 owns true and the final paired comparison. |
| VAL-048, VAL-049, VAL-050 | 11 | PASS | G6 wrong-class/result and input-shape negatives, matching atomic resolution and retained request; G7 qualification matrices. |
| VAL-051, VAL-052, VAL-053, VAL-054 | 12 | CARRY-FORWARD | Task 12 held. Its commands and removal scope are feasible; no implementation or completion is claimed. |
| VAL-055, VAL-056, VAL-057 | 13 | CARRY-FORWARD | Full leaf/error/idempotency corpus due at task 13. DEC-258 already satisfied its preliminary speed step. |
| VAL-058, VAL-059, VAL-060 | 14 | CARRY-FORWARD | Human projections, bounded resume and deterministic text fixtures due at task 14/window 13-15. |
| VAL-061, VAL-062, VAL-063 | 15 | CARRY-FORWARD | 0018 removals and maintenance convergence remain due; current inventory matches are expected. |
| VAL-064, VAL-065, VAL-076 | 16 | CARRY-FORWARD | Current-schema archive round trip, exact import fence and corrupt-payload rejection due at task 16. |
| VAL-066 | 17 | CARRY-FORWARD | Final capability matrix due at task 17. |
| VAL-067, VAL-068, VAL-069 | 18 | CARRY-FORWARD | Exact thirty-action conformance, adapter parity and cancelled-family audit due at task 18. |
| VAL-070, VAL-071, VAL-072, VAL-073 | 14/18; Human | CARRY-FORWARD | Fixture generation/regeneration followed by the ten Human-only subquestions; no QA judgment substitution. |
| VAL-077 | 6 | PASS | G6 Review Evidence journey repeats link-time independence denial; G3 production tests pass. |
| VAL-078 | 11 | PASS | BC-3 closure; G6 all four qualifying/non-qualifying Condition routes, G7 manual/atomic/shared qualification and scope matrices. |
| VAL-079 | 18/phase end | CARRY-FORWARD | Fresh G3/G5 and G8 satisfy checkpoint feasibility/headroom. Final full gate and three samples still required; section-22 full-checkptr control is additional and has no timing budget. |
| VAL-080, VAL-081 | 9/10/11 | PASS | G7 closed Condition matrix and exact nullable-binding failures; full G3 journal suite. |
| VAL-082, VAL-083, VAL-084 | 9 | PASS | G7 exact pinned six-type bytes, retained journal replay/digest and same-input semantic key. |

## Per-Task Findings

### Task 6 — Review Evidence, cycle charges and grants

Must-fix: none. MF-4 and I-5 are RESOLVED. The fresh linked-recovery and
first/later/ungranted Actor cases pass. The task-16 import leg remains under
its original owner. Improvements: none newly opened.

### Task 2 — Guarded-action evaluation

Must-fix: none. BC-1 is RESOLVED under DEC-249. The repair is confined to
requested-action applicability and matching evidence; it preserves full read
evaluation and does not change `run.close`. Improvements: none newly opened.

### Task 10 — Human Outcomes and mixed-state evidence

Must-fix: none. MF-5 and I-6 are RESOLVED; BC-2 is RESOLVED by DEC-251 and the
actual Human recording tests. The mandatory gate-satisfied,
review-unsatisfied evidence exists, and the all-satisfied terminal transition
remains intact. Improvements: none newly opened.

### Task 11 — Recorded judgment Conditions and delegated closure

Must-fix: none. BC-3 and I-2 are RESOLVED under DEC-250/253. ESC-23-005's
delegated closing evidence is accepted under DEC-258, with the original
failed attempts preserved. I-1's immediate headroom concern is RESOLVED by
the current measured command forms and DEC-259; CF-044 is retained.

**I-7 — label historical status prose (Improvement; task-11 report; no VAL).**
In `tasks/task-11-done.md`, the final bullet under
`## Concerns, technical debt, or notes for the Tech Lead` still says
“The exact ordinary gate remains over budget” and asks for scope disposition.
That describes the pre-DEC-258 attempt. The current summary, exact gate rows,
raw receipts, commit and final `OUTCOME COMPLETE` establish the later closure
unambiguously. On the next authorized documentation pass, label that paragraph
historical and identify DEC-258/259 as the later disposition, retaining all
failed evidence. This wording cleanup is deferrable; it does not require a
repair dispatch, gate change or new authority. This session leaves the report
untouched.

## Repair Plan

None required. Both review-7 Must-fix rows have their named owning tasks,
bounded changes, required evidence and fresh independent re-review here:
MF-4 → task 6 → linked recovery/rebuild assertions → G6/G7;
MF-5 → task 10 after task 2 → mixed-state assertions and report citations →
G6/G7 plus exact-command/report-lint audit.

The protected rows followed their recorded authority before those repairs:
BC-1 → DEC-249/task 2; BC-2 → DEC-251/specification resolution with task-10
evidence; BC-3 → DEC-250/253/task 11. No finding is self-resolved in review-7.

## Blocker Candidates

None unresolved. Existing **ESC-23-004** is resolved by DEC-249/250/251;
DEC-252/253/255/256 specify the associated scope and ordered verification.
Existing **ESC-23-005** is resolved by the Human-delegated DEC-258 repair and
the original-form closing gates. DEC-259 is the separately approved
prospective command amendment. DEC-257's Runner-ledger handling was inspected
as authority, not edited or executed by QA.

No new protected product, architecture, public-interface, data-model,
dependency, security, external-infrastructure or specification decision is
needed for the checkpoint. No escalation ID is allocated.

## Artifact Integrity

The pre-commit file-state check terminated `exit_code: 0` and reported
`84 unique assertion outcomes:` with 31 PASS and 53 CARRY-FORWARD entries;
all 84 IDs occur exactly once in the outcome table. It also checked the exact
three final typed lines, a single typed verdict, balanced fences and no
trailing whitespace. The index was empty and HEAD still matched the reviewed
commit before staging this one artifact.

The Runner ledger's SHA-256 remained
`667e68478ebd86eb25bb6e8b3f7a608fac698a99442606d897dfc5b006ac0281`,
identical to the initial checkout measurement. Only this new review and the
pre-existing ledger modification appeared in `git status --short`.

## Final Disposition

PASS for the DEC-255 independent checkpoint. Must-fix 0; Improvement 1 (I-7);
unresolved blocker candidates 0. Every affected deterministic checklist item
has a fresh execution or an explicitly identified delegated receipt; later
phase obligations and Human judgments remain scheduled with their owners.

Only this review is authored and committed by QA. The original reviews,
reports, contracts and Runner ledger remain unchanged. Task 12 stayed held
throughout this session and was neither dispatched nor implemented. The
Runner owns any subsequent transition using this committed review.

Review written: project-docs/docs/phases/phase-23/qa/checkpoint-after-task-11-3.md
Must-fix: 0  Improvement: 1  Blocker candidates: 0
VERDICT: PASS
