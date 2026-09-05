# Phase 22 Runner Handoff

> Human execution approval: DEC-188, 2026-09-05.
> Approved refined pack: 8dee0a194a494f575b7c5f64caf6ed9ac7da70f8.
> Start state: branch phase-22; tasks 1..18 pending; no run ledger or Open ESC.

## Prompt

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
