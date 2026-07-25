# Phase Runner Prompt Template

Copy the relevant block below and replace placeholders before invoking the
Phase Runner agent.

## Start-Of-Phase Variant

```text
Start Phase <PHASE_NUMBER> execution.

Read:
- AGENTS.md
- project-docs/docs/current/phase_brief.md
- project-docs/docs/current/decisions_active.md
- project-docs/docs/current/progress_current.md
- project-docs/docs/current/tl_escalations_current.md
- all Phase <PHASE_NUMBER> task specs in project-docs/docs/phases/phase-<PHASE_NUMBER>/tasks/
- existing Phase <PHASE_NUMBER> QA artifacts if any
- project-docs/docs/phases/phase-<PHASE_NUMBER>/qa/validation-contract.md

Current state:
- Phase <PHASE_NUMBER> is approved for execution.
- No tasks have been executed in this phase yet.

Execution rules:
- Verify that all quality commands in AGENTS.md are concrete, not placeholders.
- Verify that the validation contract exists before task execution.
- Maintain the append-only run ledger at
  project-docs/docs/phases/phase-<PHASE_NUMBER>/run-ledger.md using
  project-docs/templates/run-ledger-template.md. Append one entry per event
  (RUN-START, TASK complete, REVIEW, TRIAGE, STOP, RUN-END); never edit or
  delete prior entries. The ledger is the only file you write and it, not chat
  notes, is the durable run state.
- Append a RUN-START entry recording the current Open ESC-NNN IDs before task
  execution; only newly appearing Open ESC IDs trigger TL-Lite triage. If the
  source is already an ESC, do not create a duplicate; reference the existing
  escalation unless the issue is clearly unrelated.
- Create or resume the correct phase branch.
- Codex delegation must use the exact custom-agent selector:
  `spawn_agent(agent_type = "qa_lead" | "developer" | "qa_reviewer" | "tl_lite")`.
  Do not omit `agent_type`, do not use `default`, `worker`, or `explorer`, and
  do not set a model override for these project roles.
- Spawn qa_lead first with `agent_type = "qa_lead"` to create project-docs/docs/phases/phase-<PHASE_NUMBER>/qa/qa-plan.md.
- Verify the QA plan exists, read it before task execution, and confirm command
  and feasibility validation evidence plus the mid-phase review cadence.
- Execute task specs strictly in order.
- Before spawning each Developer, audit the QA plan's review cadence against the
  run ledger and existing QA review artifacts. If a scheduled review is due or
  overdue, spawn `qa_reviewer` for that review first, append a REVIEW entry, and
  apply the normal TL-Lite/STOP rules before continuing implementation.
- Spawn exactly one developer per task with `agent_type = "developer"`.
- Require the developer to follow .codex/agents/developer.toml or the equivalent local agent definition.
- Require the exact commands listed in the task spec's
  `## Exact Commands To Run` section. The full task gate is required before
  task close unless the task spec explicitly says otherwise.
- Do not substitute ad hoc commands for required gates. Placeholder, missing, or
  failing commands are hard stops.
- After each developer run, run
  project-docs/tools/report-lint.sh <report path> <task spec path>.
  It deterministically verifies required headings, sentinel sections, task-spec
  command coverage, the commit hash via `git rev-parse --verify
  <hash>^{commit}`, the subject against `git log -1 --format=%s`, and that the
  report file itself is committed with no uncommitted modifications. A nonzero
  exit is a hard stop. Do not re-derive these checks by prose parsing.
- Read the report for judgment items report-lint cannot check: blocker
  semantics, deviation content, concerns needing triage, and QA-plan commands
  required for that task beyond the task spec. Every QA-plan required command
  must be present as `PASS`, or `BLOCKED` with a reason in a blocker report.
- Before appending a final STOP for report-lint failure, missing report, blocker
  report, quality-gate failure, placeholder command, missing command, skipped
  command, narrower substituted command, or ad hoc gate substitution, spawn
  `tl_lite` as the first-step triager. Append a TRIAGE entry for its verdict.
  Continue only on `VERDICT: RESUME`; stop on `VERDICT: BLOCK`, and the STOP
  entry must reference the TRG/ESC.
- Append a TASK entry to the run ledger after each verified task.
- Route soft blockers through tl_lite triage. When a triage may turn on
  "already resolved by a verified commit", verify the commit yourself first and
  attach the commit-verification evidence (hash, git rev-parse --verify result,
  subject match) to the triage trigger; TL-Lite judges evidence and never
  verifies commits itself.
- Every third completed task, spawn qa_reviewer with
  `agent_type = "qa_reviewer"` unless the QA plan sets a stricter cadence, then
  append a REVIEW entry to the run ledger.
- For QA Blocking findings, TL-Lite may RESUME only when the finding was
  misclassified, already resolved by a verified commit, or accepted by a
  recorded decision.
- Parse only the literal final TL-Lite verdict line:
  VERDICT: RESUME or VERDICT: BLOCK (escalated as ESC-<phase>-<NNN>).
  Identifiers are phase-scoped and allocated from the phase triage log. Record
  the concrete IDs in the TRIAGE entry's `Triage log:` field; if TL-Lite
  returned a placeholder or wrote no triage-log entry, record that verbatim.
- After all tasks complete, spawn qa_lead with `agent_type = "qa_lead"` for
  final review and human checklist generation.
- Verify final-review.md records every VAL-NNN outcome.
- Append the RUN-END entry to the run ledger.
- When stopping early, append a STOP entry naming the trigger and the next
  human action.
- Do not do developer, QA, or TL-Lite work yourself.
```

## Resume Variant

```text
Resume Phase <PHASE_NUMBER> execution from Task <START_TASK> on branch <PHASE_BRANCH>.

Read:
- AGENTS.md
- project-docs/docs/current/phase_brief.md
- project-docs/docs/current/decisions_active.md
- project-docs/docs/current/progress_current.md
- project-docs/docs/current/tl_escalations_current.md
- project-docs/docs/phases/phase-<PHASE_NUMBER>/run-ledger.md (authoritative
  run state; read before anything below)
- project-docs/docs/phases/phase-<PHASE_NUMBER>/qa/qa-plan.md
- project-docs/docs/phases/phase-<PHASE_NUMBER>/qa/validation-contract.md
- project-docs/docs/phases/phase-<PHASE_NUMBER>/triage-log.md if it exists
- completed task reports in project-docs/docs/phases/phase-<PHASE_NUMBER>/tasks/
- remaining task specs in project-docs/docs/phases/phase-<PHASE_NUMBER>/tasks/

Current state:
- Branch already exists: <PHASE_BRANCH>
- Completed tasks before this run: <summary>
- Start at Task <START_TASK>

Execution rules:
- Reconcile the run ledger against completion reports and git history before
  executing anything. If the ledger and reports disagree, stop and report the
  discrepancy. If the ledger is missing for an in-progress phase, reconstruct
  a RUN-START entry from task reports and git history and note the gap.
- Audit the QA plan's mid-phase review cadence against the run ledger and QA
  review artifacts before spawning any Developer. If a scheduled review is
  overdue, spawn `qa_reviewer` for that review first, append a REVIEW entry, and
  apply the normal TL-Lite/STOP rules before continuing implementation.
- If resuming from an existing STOP caused by a task blocker or failed required
  gate and no TRIAGE entry exists for that STOP, spawn `tl_lite` first to
  classify it, append the TRIAGE entry, and continue only on `VERDICT: RESUME`.
- Do not recreate the branch.
- Do not rerun completed tasks.
- Verify completed-task commits are present in current history.
- Resume strictly from Task <START_TASK>.
- For a QA Blocking resume, continue only when the issue was misclassified,
  already resolved by a verified commit, or accepted by a recorded decision,
  with commit-verification evidence attached to the triage trigger.
- Apply all start-of-phase execution rules from this template, including
  run-ledger appends and report-lint verification.
```
