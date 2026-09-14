# Phase 22 Runner Handoff

> Human execution approval: DEC-188, 2026-09-05.
> Approved refined pack: 8dee0a194a494f575b7c5f64caf6ed9ac7da70f8.
> Initial start state (historical): tasks 1..18 pending; no ledger or Open ESC.
> Current restart checkpoint: DEC-190/191; ESC-22-001 resolved; valid ledger
> ends in task-2 TASK-RECOVERY attempt 2, ready for Developer dispatch.

## Current Recovery Handoff - 2026-09-06

Resume Phase 22 on `phase-22` from the existing ledger. The protected verifier
repair is commit `3bb42b4efd391b216b752db2c31fe793f3538bac`. All 30 regression
checks, complete mechanical drills and workflow lint passed. Normal
RUN-RESUME, correcting TRIAGE TRG-22-005 and the exact task-2 recovery were
already accepted on the real ledger, each returning
`APPEND_CONFIRMED state=REPAIR_PENDING` and `exit_code: 0`. The escalation
register has no Open item. DEC-188 execution approval remains in force.

1. Read the current role/context and reconcile the final ledger entry with
   TRG-22-005, DEC-189..191, task-2's report and review-1. The old TRIAGE
   TRG-22-003 remains history; TRG-22-005 is the current authority. TRG-22-006
   retains its original BLOCK plus a separate verified resolution.
2. Append normal RUN-RESUME for task-2 with all STOP-only fields `N/A`.
   Neither rejected BLOCK nor STOP was persisted; do not invent a STOP
   reconciliation, another TRG, or a new task-2 recovery attempt.
3. Dispatch the **already recorded** task-2 TASK-RECOVERY at attempt 2 to a
   fresh `developer` child. Preserve the exact review trigger, DEC-189
   authority and two-file allowlist. The child has not yet run this attempt.
   It must confirm the coverage-child warning mechanism before repairing
   `cmd/fid/live_output_test.go`, leave both goldens unchanged, update
   `tasks/task-2-done.md`, and run every exact command in the directive/task.
4. After its typed result, follow normal REPAIR/report verification. Finish
   task-3 review-1 MF-3 in its own Developer continuation and obtain a fresh
   independent review-1 re-review before advancing to task 4.
5. Continue the previously approved sequence below from that checkpoint.
   Do not repeat initial readiness, completed tasks 1..3, or the existing
   task-2 attempt-1 review repairs merely because the original prompt follows.

Maintenance evidence: `../../reviews/2026-09-06-phase-22-workflow-recovery.md`.
The maintenance task changed no product code, product tests, or goldens and
does not claim task-2 product gates or review closure.

## Original Execution Prompt (Apply From The Current Checkpoint Above)

You are the Phase Runner for fid. Execute the approved Phase 22 failure-and-
recovery plan in /home/mark/projects/fid on branch phase-22. Human approval is
recorded in project-docs/docs/phases/phase-22/decisions.md as DEC-188. Heavy mode,
all 18 ordered tasks, all 80 VAL assertions, CF-035 and DEC-184..187 are approved.
Do not ask for the same approval again.

Load AGENTS.md and the Phase Runner role in project-docs/agents/phase-runner.md
(on Claude also CLAUDE.md); follow 09_AI_WORKFLOW.md section 2.4. Read the current
phase brief, decisions, progress and escalation register. Under
project-docs/docs/phases/phase-22/, read decisions.md, execution-scope.md,
v1-convergence-inventory.md, qa/validation-contract.md, qa/qa-plan.md,
qa/tl-refinement.md and kit-contract-mapping.md. Load the ordered task specs
and their exact context references as needed; do not load a full spec suite.

Verify branch and git history against the approved pack and DEC-188. If no
ledger exists, start a new run; if one appeared, reconcile it with committed
reports before resuming. The reported initial Open ESC baseline is empty;
verify the current register. Never assume a task ran from chat or a status label.

Perform the required QA Lead readiness refresh first, using the existing plan
and TL addendum. Preserve original QA findings and failed evidence. Have QA
check the refined DDL, digest/boundary evidence strategy, new exact commands and
task coverage; repair/triage any concrete gap through the normal workflow.
Once readiness is confirmed, dispatch tasks/task-1.md, then tasks 2..18 in order.

Use exactly the project custom child roles qa_lead, developer, qa_reviewer and
tl_lite. Set agent_type explicitly; no default/worker/explorer substitutes or
model overrides. You are the active Runner; never spawn another phase_runner.
Keep one active Developer child, one task per Developer session, with exact
task/report paths, required reads, gate commands and commit/report obligations.

Write only run-ledger.md through workflow-sequence-lint.sh --append-stdin,
following the ledger template. Each RUN-START/RUN-RESUME carries exactly one
Sequence contract: v2 and Recovery directive schema: v2. Require the real
APPEND_CONFIRMED receipt and valid sequence state before advancing. Use the
canonical workflow-route.sh, report-lint.sh and recovery protocol; do not
replace their results with prose judgment or edit controls yourself.

All tasks require ordinary race/coverage and tagged subprocess integration.
Task 1 and task 16 add one reference rebuild sample; task 17/final QA retain
three fresh samples and the unchanged <30 s limit. Tasks 2/7 capture CF-034
coverage before edits. Task 8 owns R-14's signal/timeout fixture repair and
repeated verification; if it fails earlier, use DEC-186's bounded same-task
scope-amendment route. Do not defer a failing required gate until task 8 or
retry until green. Retain failed attempts and observed termination evidence.

Review after tasks 3, 6, 9, 12, 15 and 18. The separate independent QA Lead
checkpoint after task 12 must pass before task 13. Use fresh independent
re-review after QA repair. Poll live children/commands to actual completion;
parent wait timeouts are not failures or permission to kill a running child.
Ordinary recoverable failures stay in same-task recovery; protected blockers
go through TL-Lite with durable TRG/ESC records and the prescribed STOP rules.

Continue through the approved sequence and required final QA. Obtain actual
Human outcomes for VAL-076..078; never attest them yourself. Hand the final
review/checklist and unresolved recommendations back for Tech Lead wrap-up.
Do not merge, release, claim production kit qualification, or alter held-out
scenario material. Do not implement product code/tests or write other roles'
artifacts yourself.
