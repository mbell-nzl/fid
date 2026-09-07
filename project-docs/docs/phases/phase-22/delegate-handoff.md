# Phase 22 Delegated Completion

## Human Instruction — 2026-09-07

The Human authorized the root delegate to orchestrate the rest of Phase 22,
including final QA, so the Human checklist is ready for the following day.
The Human explicitly permitted canonical subagents and direct root edits to
repair routine workflow friction without repeated role handoffs. This is
execution authority from the current request, not a new product decision.

Preserve accepted product contracts, deterministic evidence, failed attempts,
and independent review. Human VAL-076..078 remain for the Human. No release,
merge, or completed Human sign-off is implied.

## Starting Checkpoint

Root inspected branch `phase-22` at
`3438f8cc0b445247069611f2d74fb827a2187eb7` (DEC-195 docs resolution).
Only `run-ledger.md` was dirty, containing earlier execution history through
the task-8 STOP. Tasks 1-7 and reviews 1b/2b are recorded complete. This is
historical artifact inspection, not a fresh attestation of their test gates.

The Phase Runner child appended the task-8 resume/recovery but reported that
its tool surface could not spawn children. Root took over orchestration with
a normal RUN-RESUME, retaining the bound attempt 4. Root owns the append-only
ledger and direct canonical task/review dispatch for this delegated run.

Task 8 completed at product `c5dc59e4bb9cf9db20a4aa157daeef68d132bb60`,
report `622089e59211ac5612674651ef917a7002e41ba9`. Root observed report-lint
`OUTCOME COMPLETE` with `exit_code: 0`, verified commit identity/subject,
appended REPAIR RESOLVED/TASK COMPLETE, and observed `VALID state=ADVANCEABLE`.
Task 9 is dispatched. All product gate results retain Developer-session
attribution in the report; root's fresh checks are orchestration evidence.

## Remaining Sequence And Delivery

- Complete task 8, then 9; review window 3.
- Complete tasks 10-12; review window 4 and independent QA Lead checkpoint
  after task 12 before dispatching task 13.
- Complete tasks 13-15; review window 5.
- Complete tasks 16-18; review window 6.
- Complete fresh final deterministic QA, including all 77 QA-owned assertions,
  the exact phase gates, reference rebuild samples, and comparison obligations.
- Provide verified, reproducible Human fixtures and the six separately
  markable Human checklist rows covering VAL-076..078.
- Refresh current context to distinguish deterministic readiness from pending
  Human validation. Full phase closure and post-Human wrap-up follow the
  actual Human outcomes; no pending outcome is upgraded to PASS.

The QA plan's full phase-end decision rule still governs full phase closure.
Preparing the final deterministic review before Human execution implements
the Human's requested handoff order; it does not waive the Human criteria.

## Coordination

DEC-196 records the preflight allocation repairs for tasks 9/12/18: semantic
fake assertions, Human resume rendering, and retained fixture generation.
Review 3 must independently check the fake repair; the task-12 checkpoint
checks both renderings; final QA verifies the retained fixture instructions.

The root delegate handles routine scope/documentation/control friction when
accepted behavior determines the repair. Preserve the failed attempt and
record any bounded scope amendment before a canonical Developer executes it.
Any new product or architecture choice remains explicit, never an implicit
substitution to obtain a green gate.

Root refreshed stale current progress and approval text at dispatch. The
ledger and actual completion/review artifacts remain the detailed execution
record. Root does not edit active child-owned implementation files or the
Runner ledger without coordinating ownership first.

During task 9, root owned the semantic corpus while the Developer owned
product code. Root handed back and committed its separately attributed
evidence before the Developer's final gates. Root resolved omitted causal file
scope directly in DEC-197 without creating a protected-control repair cycle.
The Developer's brief scope-request return did not abandon or replace the
active task: root supplied the exact cleanup authority and continued the same
agent. Failed lint and integration receipts remain in that task's report.

Review 3 may prepare by inspecting completed tasks 7/8 and the committed corpus
while task 9 finishes. Its final source verification, representative gate,
verdict and commit wait for root's verified task-9 completion checkpoint. No
review outcome or next-task dispatch is advanced on mutable partial work.

Review 3 passed at `3305edea74862a752be6376b49aed5e2f809128e` and is
recorded in the ledger. Root initially invoked the route helper with the
unsupported source `qa-review` (usage output, terminated `exit_code: 2`),
then corrected the source to `qa`; `workflow-route.sh qa PASS` returned
`ACTION ADVANCE` (terminated `exit_code: 0`). This was an invocation error;
no routing control or accepted outcome changed.

## Preparatory Final QA Dispatch

Task 10 root assistance: by explicit ownership agreement, root added
`cmd/fid/activation_export_test.go` and handed it back before the Developer's
final gates. The test invokes export create/validate during a prepared upgrade
and checks that the activation projection and journal Event count remain
unchanged after each command. Root ran, from `/home/mark/projects/fid`,
`go test ./cmd/fid -run '^TestIncompleteActivationAllowsExportWithoutChangingActivation$' -count=1`:
PASS — `ok  	github.com/mbell-nzl/fid/cmd/fid	0.293s`
(terminated `exit_code: 0`). This is a regression check of retained behavior,
not a claimed red/green product change. The Developer's own final gates execute
it again; root's focused receipt retains root attribution.

Root also accepted ownership of the unstarted per-resource matrix in
`cmd/fid/activation_resource_matrix_test.go`, then handed it back before source
repair. From `/home/mark/projects/fid`, root ran
`go test ./cmd/fid -run '^TestActivation(RequiresEveryResourceAtCommitAndFinalize|CompleteResourceAgreementFinalizesTrust)$' -count=1`
twice. Attempt 1 ended with `FAIL` (package result 1.463s; terminated
`exit_code: 1`). Each missing row and each NULL observed identity incorrectly
reached commit/finalize. Pending and mismatched cases did reject, but the test
also incorrectly required an uncontracted message substring. Root corrected
that construction assumption to structured rejection, retaining unchanged
state/Event count and denied-trust assertions; typed adapter errors are task
11's separate obligation. Attempt 2, the byte-identical command, ended with
`FAIL` (package result 1.330s; terminated `exit_code: 1`), exposing only the
real missing-row/NULL-observation acceptance. The causal query counted
disagreements among present rows, so missing rows and SQL NULL comparison were
not counted. No product source changed before either root red run. The matrix
has 40 negative states plus an all-agree positive control; its direct projection
fixtures isolate transition guards and do not claim producer or replay parity.
Developer owns causal repair and subsequent green receipts after handback.

Before task-10 close, root found that VAL-040's cited prepare test did not
actually inspect predecessor/resource bindings, and the other cited tests did
not reject mixed internal digest selections. The Developer confirmed the gap.
Root added `cmd/fid/activation_binding_test.go`, covering first/upgrade frozen
bindings across staging/observation/restore, new-key immutable-content rejection,
missing/mismatched/mixed/prior-only digest negatives and the exact new-set
positive upgrade through finalization. From `/home/mark/projects/fid`, root ran
`go test ./cmd/fid -run '^Test(ActivationFreezesFirstAndUpgradeBindingsAcrossStagingAndRecovery|ActivationCommitRejectsMissingMismatchedAndMixedComponentDigests|StagedActivationContentCannotBeReplacedUnderANewKey)$' -count=1`.
It ended with `FAIL` (package result 0.655s; terminated `exit_code: 1`): the
positive upgrade committed but failed to finalize. The other binding and
rejection cases passed. Root handed this test back for Developer diagnosis and
causal repair, including the same upgrade ordering in forward recovery. The
earlier five PASS gates remain evidence for the earlier product checkpoint;
the Developer must rerun them on the final source before task completion.

After task 18 and review 6, the QA Lead's scope is a fresh deterministic final
review plus a feasible Human checklist. Its report must explicitly distinguish
that requested readiness verdict from the plan's full phase-end sign-off,
which still awaits Human judgment. Record all 80 VAL rows; use `BLOCKED` with
the precise reason “Human judgment not yet performed” for VAL-076..078 in
this preparatory report, never PASS or an unexplained carry-forward. Pending
scheduled Human judgment is not a new product/spec escalation. All 77
QA-owned rows and QA-01..84 still require the QA Lead's own reproducible
evidence, including reference identity, three fresh consecutive rebuild
samples, both full benchmark runs, source-currency comparison and security
review. The root supplies no proxy pass for those checks.

Keep Human-facing prompts neutral and the six answers independently markable.
Put expected judgments in separately inspectable QA material, rather than
providing the answer in the first inspection instructions. Validate the exact
retained binary, manifest IDs, commands and fixture roots after generator exit.
If source changed after the prior retained generation, regenerate from final
committed source and preserve both receipts. Release gates and held-out
scenarios remain explicitly outside this phase per the accepted QA plan.
