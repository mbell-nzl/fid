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
- Maintain a compact run ledger: completed tasks, verified commits, QA reviews,
  ESC IDs seen, TRG IDs logged, and warm-up/autonomous-resume state.
- Record the current Open ESC-NNN IDs before task execution; only newly appearing
  Open ESC IDs trigger TL-Lite triage. If the source is already an ESC, do not
  create a duplicate; reference the existing escalation unless the issue is
  clearly unrelated.
- Create or resume the correct phase branch.
- Spawn qa_lead first to create project-docs/docs/phases/phase-<PHASE_NUMBER>/qa/qa-plan.md.
- Verify the QA plan exists, read it before task execution, and confirm command
  and feasibility validation evidence plus the mid-phase review cadence.
- Execute task specs strictly in order.
- Spawn exactly one developer per task.
- Require the developer to follow .codex/agents/developer.toml or the equivalent local agent definition.
- Require the exact commands listed in the task spec. The full task gate is
  required before task close unless the task spec explicitly says otherwise.
- Do not substitute ad hoc commands for required gates. Placeholder, missing, or
  failing commands are hard stops.
- After each developer run, verify the completion report exists at the expected
  path, required headings are present, Validation Contract Coverage is present,
  the concerns section is non-empty, the dependency-budget note is present, and
  commit hash/subject are present.
- Verify the commit hash with `git rev-parse --verify <hash>^{commit}` and
  verify the report subject matches `git log -1 --format=%s <hash>`.
- Compare `## Tests run` and `## Exact quality-gate results` with the task
  spec's `## Exact Commands To Run` and every QA-plan command required for that
  task. Every required command must be present as `PASS`, or `BLOCKED` with a
  reason in a blocker report.
- Stop immediately on missing report, blocker report, quality-gate failure,
  placeholder command, missing command, skipped command, narrower substituted
  command, missing dependency-budget note, or unverifiable commit.
- Route soft blockers through tl_lite triage.
- Every third completed task, run qa_reviewer unless the QA plan sets a stricter
  cadence.
- For QA Blocking findings, TL-Lite may RESUME only when the finding was
  misclassified, already resolved by a verified commit, or accepted by a
  recorded decision.
- Parse only the literal final TL-Lite verdict line:
  VERDICT: RESUME or VERDICT: BLOCK (escalated as ESC-NNN).
- After all tasks complete, run qa_lead final review and human checklist generation.
- Verify final-review.md records every VAL-NNN outcome.
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
- Do not recreate the branch.
- Do not rerun completed tasks.
- Verify completed-task commits are present in current history.
- Resume strictly from Task <START_TASK>.
- For a QA Blocking resume, continue only when the issue was misclassified,
  already resolved by a verified commit, or accepted by a recorded decision.
- Apply all start-of-phase execution rules from this template.
```
