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
- project-docs/docs/phases/phase-<PHASE_NUMBER>/decisions.md
- all Phase <PHASE_NUMBER> task specs in project-docs/docs/phases/phase-<PHASE_NUMBER>/tasks/
- existing Phase <PHASE_NUMBER> QA artifacts if any
- project-docs/docs/phases/phase-<PHASE_NUMBER>/qa/validation-contract.md

Current state:
- Phase <PHASE_NUMBER> is approved for execution.
- No tasks have been executed in this phase yet.

Execution rules:
- Verify that all quality commands in AGENTS.md are concrete, not placeholders.
- Preserve the dedicated real-process integration package named in AGENTS.md;
  do not widen its command to `./...` and duplicate the ordinary suite.
- Verify that the validation contract exists before task execution.
- Maintain the append-only run ledger at
  project-docs/docs/phases/phase-<PHASE_NUMBER>/run-ledger.md using
  project-docs/templates/run-ledger-template.md. Append one entry per event
  (RUN-START, RUN-RESUME, TASK complete, TASK-RECOVERY, REPAIR, REVIEW,
  TRIAGE, STOP, RUN-END); never edit or
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
- Keep exactly one active Developer at a time. A Developer session handles at
  most one task, but a recoverable task may use a same-task continuation.
- Require the developer to follow .codex/agents/developer.toml or the equivalent local agent definition.
- Require the exact commands listed in the task spec's
  `## Exact Commands To Run` section. The full task gate is required before
  task close unless the task spec explicitly says otherwise.
- Do not substitute ad hoc commands for required gates. A failing in-scope
  command is REPAIR: preserve the failed attempt, make the causal repair, and
  require a byte-identical rerun. Incomplete work is not a blocker.
- Accept a completion report's claim that a command terminated only when its
  evidence includes the nested command exit code or an explicit termination
  signal. Completion of an outer execution cell proves only that its wrapper
  returned; empty output or absence of new output is not command completion.
- After each developer run, run
  project-docs/tools/report-lint.sh <report path> <task spec path>.
  Parse its literal OUTCOME COMPLETE, OUTCOME RECOVERY_REQUIRED, or OUTCOME
  BLOCKED line. It verifies typed state, command coverage, observed
  termination, commit identity, and that the report is committed and clean.
  A nonzero exit is report-repair. Do not re-derive failed checks by prose.
- Route every parsed Developer, report-lint, QA, TL-Lite, command-liveness, and
  recovery-budget result through
  project-docs/tools/workflow-route.sh <source> <outcome> [attempts]. Follow
  its literal ACTION. WAIT retains and polls the live command session;
  EVALUATE inspects the observed terminated result. Invalid router input is a
  recoverable orchestration defect, not permission to invent a route.
- Put `Sequence contract: v2` in every new RUN-START/RUN-RESUME; older entries
  stay legacy-tolerant. Before each append, pipe the exact proposed entry to
  project-docs/tools/workflow-sequence-lint.sh <run-ledger> --append-stdin and
  require exit zero. Append only a valid candidate, then run the tool on the
  ledger as a postcondition. workflow-route.sh still chooses the action.
  Require literal VALID state=ADVANCEABLE before actual advancement. The
  validator requires the same task across recovery, a recorded REPAIR before
  task completion, and a differently identified QA PASS review whose Re-review
  of field names the original REPAIR review. Never append an INVALID candidate
  or rewrite ledger history; retain and stop on a post-append integrity mismatch.
- Read the report for judgment items report-lint cannot check: blocker
  semantics, deviation content, concerns needing triage, and QA-plan commands
  required for that task beyond the task spec. Every QA-plan required command
  must be present as `PASS`, or `BLOCKED` with a reason in a blocker report.
- On RECOVERY_REQUIRED, report-lint failure, missing report, in-scope gate
  failure, silently skipped command, missing termination, session end with
  partial work, or QA REPAIR,
  append TASK-RECOVERY with trigger signature, attempt count, checkpoint,
  evidence, progress, and next action; then dispatch a same-task continuation.
  Use two consecutive no-progress continuation attempts per trigger signature
  as the default budget. Same-session red/green iterations do not consume it;
  incomplete work is not a blocker. Count progress only when a relevant
  in-scope change tied to the trigger satisfies/removes an obligation or
  verification demonstrably narrows the failure. Formatting, metadata-only,
  unrelated, or no-op churn is not progress.
- Treat OUTCOME BLOCKED as a blocker candidate, never authoritative state.
  Before STOP, spawn `tl_lite` as the first-step triager. RESUME advances
  only when cleared; VERDICT: REPAIR dispatches recovery without advancing;
  VERDICT: BLOCK appends STOP referencing the TRG/ESC.
- Append a TASK entry to the run ledger after each verified task.
- Route only protected blocker candidates or ambiguous classification through
  tl_lite triage. Ordinary concerns remain for wrap-up. When triage may turn on
  "already resolved by a verified commit", verify the commit yourself first and
  attach the commit-verification evidence (hash, git rev-parse --verify result,
  subject match) to the triage trigger; TL-Lite judges evidence and never
  verifies commits itself.
- Every third completed task, spawn qa_reviewer with
  `agent_type = "qa_reviewer"` unless the QA plan sets a stricter cadence, then
  append a REVIEW entry to the run ledger.
- QA Outcome: PASS advances. QA Outcome: REPAIR must name an owning task,
  bounded repair, required evidence, and fresh independent re-review; dispatch
  it directly. QA Outcome: BLOCK is a blocker candidate.
- Parse only the literal final TL-Lite verdict line:
  VERDICT: RESUME, VERDICT: REPAIR, or
  VERDICT: BLOCK (escalated as ESC-<phase>-<NNN>).
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
- project-docs/docs/phases/phase-<PHASE_NUMBER>/decisions.md
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
- If a prior STOP has a later decision, verified commit, restored capability,
  or other recorded resolution, verify it and append RUN-RESUME linking STOP,
  TRG/ESC, resolution, branch/HEAD, and next task. A prior TRIAGE BLOCK does not
  prohibit resolved resume.
- If no TRIAGE entry exists for an existing STOP, spawn `tl_lite` first. RESUME
  advances only if cleared; REPAIR dispatches the same task; BLOCK stops.
- Do not recreate the branch.
- Do not rerun completed tasks.
- Verify completed-task commits are present in current history.
- Resume strictly from Task <START_TASK>, including a same-task continuation
  when the latest durable outcome is RECOVERY_REQUIRED.
- For QA REPAIR, dispatch the owning task and require fresh independent
  re-review before advancing.
- Apply all start-of-phase execution rules from this template, including
  run-ledger appends and report-lint verification.
```
