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
  TRIAGE, STOP, RUN-END) only through workflow-sequence-lint.sh
  --append-stdin; never edit or delete prior entries. The ledger is the only
  file you write and it, not chat notes, is the durable run state.
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
- Route child-agent liveness separately. A parent polling wait timeout is not
  child termination: route `session WAIT_TIMEOUT`, follow `ACTION WAIT`, and
  poll the same child again. A running child routes as `session LIVE`; never
  interrupt it because a parent wait expired or a command budget elapsed at the
  session level. A returned child routes as `session RETURNED` before its typed
  result is evaluated. A child externally interrupted before returning routes
  as `session INTERRUPTED <consecutive-count>`: first to recovery, second to
  orchestration session-liveness triage. Record ledger `Outcome: INTERRUPTED`
  with interruption cause, elapsed session duration, and command in flight;
  use the exact command, `none observed`, or `unknown after external
  termination`, never an invented value.
- Put exactly one `Sequence contract: v2` and exactly one
  `Recovery directive schema: v2` in every newly appended
  RUN-START/RUN-RESUME. Both become effective in that event and schema may
  never downgrade. A standalone historical v1 ledger remains valid, with
  general identity, STOP, authority, budget, and state checks still enforced;
  its next appended RUN uses v2. Every mandatory field in every schema-v2 event
  occurs exactly once, using `N/A` only where the ledger template permits it.
  Field order and contiguity do not carry authority. Pipe each exact
  candidate entry, beginning with its
  single ### event heading, to project-docs/tools/workflow-sequence-lint.sh
  <run-ledger> --append-stdin. The tool validates and atomically appends at EOF;
  INVALID leaves the ledger byte-for-byte unchanged. Require exit zero and
  literal APPEND_CONFIRMED state=<state>. If termination or receipt is missing,
  retry the byte-identical candidate; an exact EOF match confirms it without a
  duplicate. Do not advance, submit a different event, or dispatch work while
  that append result is ambiguous. A terminated INVALID rejection is conclusive:
  validate the unchanged ledger, correct the candidate against that state, and
  resubmit through the append command. Preserve the rejected output and exit
  status in the corrected event's reconciliation evidence; a rejected candidate
  with a valid ledger does not require a stop. Never separately mutate the
  ledger with apply_patch, redirection, or another editor. workflow-route.sh
  still chooses the action.
  Run sequence lint without --append-stdin to read current state and require
  literal VALID state=ADVANCEABLE before actual advancement. After a v2 schema
  marker, the validator binds every typed control directive value byte-for-byte
  to the governing phase triage-log entry, including runner-recorded base and
  originating report. It also binds the external TRG and ledger TRIAGE's exact
  `Originating task`, `Trigger signature`, `Attempt count`, `Repair category`,
  and `Repair authority`; direct REPAIR carries all five into TASK-RECOVERY.
  That recovery's Source references the governing TRG, while its exact-bound
  ordinary Repair authority is not changed merely to repeat the TRG ID.
  When direct schema-v2 TRIAGE REPAIR classifies a blocker candidate returned
  by an active recovery awaiting REPAIR, its bound recovery replaces that
  pending dispatch only for the same task; other overlaps remain invalid.
  It rejects STOP without a preceding TRIAGE BLOCK and concrete ESC, validates
  every DEC citation in an authority field, and binds a STOP resume to its
  exact next recovery. It requires the same task across recovery, a recorded
  REPAIR before task completion, the exact report-repair after control
  RESOLVED, and a differently identified QA PASS review whose Re-review of
  field names the original REPAIR review. A linked re-review may transfer repair
  ownership only after the current owner completes; it cannot overwrite an
  unresolved owner. An INVALID candidate is recoverable;
  a missing receipt is resolved only by byte-identical retry. Retain and stop
  only if standalone validation fails after a confirmed append.
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
- An externally interrupted child that returned no typed result never consumes
  the no-progress budget and does not reset it. Preserve the prior count in the
  next TASK-RECOVERY. No budget bounds the child session lifetime; only
  declared per-command budgets apply, and a parent polling wait timeout means
  poll again. Two consecutive external interruptions require orchestration
  session-liveness triage, never `recovery-budget-exhausted`. For long
  verification, checkpoint durable implementation before running commands and
  append each result to the report as it terminates. A repair verdict from that
  triage advances total attempt count, preserves the no-progress count, and
  starts a new interruption streak.
- Keep attempt count and consecutive no-progress count monotonic for a trigger
  signature across runner sessions, branches, and worktrees. Material progress
  resets only consecutive no-progress to 0/2. A second consecutive no-progress
  outcome exhausts the budget and returns to triage; never start another
  continuation by resetting or renaming the trigger.
- Distinguish a protected gate/control definition change from control
  conformance. Definition changes alter a command, assertion's accepted
  pass/fail semantics, accepted status/scope, threshold/budget/tolerance, or
  policy. A canonical verification/workflow control implementation that
  violates an already accepted, uniquely determining contract is eligible for
  REPAIR; it is not BLOCK merely because the control is canonical or outside
  the originating task's ordinary files. A failure alone never authorizes edits
  to a test, gate, or control.
- Route a proposed cross-scope control-conformance repair to TL-Lite. Require a
  written directive with these exact one-line labels: Originating task, Control
  trigger signature, Control attempt count, Control contract, Control mismatch,
  Allowed files, Forbidden changes, Exact verification, Originating report,
  and Control repair base.
  Require each label exactly once and each value on one physical line; label
  order and contiguity do not matter. `Allowed files` is one exact
  comma-separated ordered list with no spaces around commas.
  Accept only the exact verdict `VERDICT: REPAIR / control-conformance` for
  this route. Ambiguity or a protected definition change is BLOCK.
- For `REPAIR / control-conformance`, first verify that the phase worktree has
  no dirt except the run ledger and no hidden index flags. Record and verify the
  full phase HEAD yourself before TL-Lite and Developer dispatch; supply it as
  Control repair base with the other nine exact values above. Confirm the
  triage-log entry, append ledger TRIAGE referencing it, then append
  TASK-RECOVERY by copying all directive values byte-for-byte. Ledger Trigger
  signature and initial Attempt count map to Control trigger signature and
  Control attempt count. Leave the ledger dirt in the phase worktree.
  Create a dedicated clean temporary repair branch and worktree at the recorded
  starting commit, and supply its path, branch, base, allowlist, and directive
  verbatim to a dedicated same-task control-repair continuation. The Developer
  works only in that repair worktree, verifies HEAD equals Control repair base,
  never rebaselines, uses red/green, and commits only the tools-only repair.
  Do not advance or ask the Human. Never authorize normal control-conformance
  self-repair of workflow-control-scope-lint.sh or workflow-sequence-lint.sh;
  either verifier requires separate protected decision authority.
- Never use `git update-index --assume-unchanged`, `--skip-worktree`, or an
  equivalent index flag to conceal phase- or repair-worktree dirt. Before and
  after the repair, the phase worktree may be dirty only at the run ledger; the
  repair worktree must be clean at its checked-out ending commit.
- Before appending the REPAIR ledger event for that commit, materialize
  workflow-control-scope-lint.sh from the trusted Control repair base into
  fresh scratch storage, never from repair HEAD, execute that materialized blob,
  and require PASS with observed termination. Pass the repair worktree,
  runner-recorded base, Developer ending HEAD, and individual allowed paths in
  their CSV order. Never accept a Developer-selected base. Require the PASS
  receipt's comma-separated ordered allowlist to be byte-identical to `Allowed
  files` in the exact receipt `PASS control repair scope: <base> -> <head>; allowed=<Allowed files> (terminated exit_code: 0)`; reject missing,
  normalized, reordered, or different receipt data.
  Product, spec, prompt, template, dependency, task,
  QA-contract, generated, and report files are forbidden. The claimed repair
  commit must be checked-out HEAD, the index/worktree clean, and every
  allowlisted artifact a regular file; later commits, dirty state, symlinks,
  and gitlinks fail. Provider sync is a separate authorized maintenance step.
- After every control scope PASS, fast-forward the phase branch to repair HEAD
  without a merge commit or cherry-pick, verify the phase worktree still has no
  dirt except the run ledger, append REPAIR, and clean up after its receipt.
  After MATERIAL-PROGRESS or NO-PROGRESS, issue the next control TASK-RECOVERY
  with Attempt count incremented and Control repair base equal to that prior
  integrated HEAD; preserve authority, contract, mismatch, allowlist, forbidden
  set, originating report, and exact verification byte-for-byte.
- Control RESOLVED permits only a linked Category: report-repair / Repair
  mode: task-scope continuation whose Allowed files is exactly Originating
  report and whose Exact verification and Control repair base equal the
  resolved attempt's values. It edits only that report.
  Direct TASK complete is invalid; only the report-repair's later RESOLVED
  permits completion.
- Treat OUTCOME BLOCKED as a blocker candidate, never authoritative state.
  Before STOP, spawn `tl_lite` as the first-step triager. RESUME advances
  only when cleared; VERDICT: REPAIR dispatches recovery without advancing;
  VERDICT: BLOCK appends STOP referencing the TRG/ESC.
- For every schema-v2 TRIAGE, require the external TRG and ledger event to
  contain exactly one `Originating task`, `Trigger signature`, `Attempt count`,
  `Repair category`, and `Repair authority`. Direct REPAIR binds all five;
  RESUME/BLOCK uses `N/A` for repair category/authority when none applies.
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
  VERDICT: RESUME, VERDICT: REPAIR,
  VERDICT: REPAIR / control-conformance, or
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
  executing anything. An active recovery can have newer implementation commits
  than its report; continue that recovery to finish verification and reporting.
  Stop only if branch/worktree state cannot be reconciled non-destructively.
  If the ledger is missing for an in-progress phase, reconstruct
  a RUN-START entry from task reports and git history and note the gap.
- Audit the QA plan's mid-phase review cadence against the run ledger and QA
  review artifacts before spawning any Developer. If a scheduled review is
  overdue, spawn `qa_reviewer` for that review first, append a REVIEW entry, and
  apply the normal TL-Lite/STOP rules before continuing implementation.
- Choose the resume form from standalone sequence validation and the ledger,
  ahead of stale current-context summaries. A historical STOP with an accepted
  RUN-RESUME is already consumed. For a normal session resume, use N/A for
  Governing TRG, Governing ESC, Trigger signature, Attempt count, Resolution
  authority, and Next recovery; omit Repair owners entirely. Keep Continuing
  task and checkpoint concrete. An existing TASK-RECOVERY keeps its task,
  trigger, attempt, authority, and review-owner progress. Reconcile its commits,
  report, and child-session state, then continue it without a duplicate
  TASK-RECOVERY or rerunning completed owners. REPAIR_PENDING permits continuing
  the active repair, not advancing to the next task.
- If the current ledger state is STOPPED and that STOP has a later decision,
  verified commit, or restored
  capability, verify it and append RUN-RESUME retaining the
  STOP's governing phase-scoped TRG/ESC, trigger, next attempt, continuing task,
  and exact Next recovery tuple. Record `Resolution authority` as exactly one of
  `verified-commit:<40hex>`, `restored-capability:<stable-id>`, or `DEC-<NNN>`,
  with separate Runner-observed evidence. Copy that exact authority into the
  next recovery's Source and Repair authority. A protected decision change or
  resolution requires DEC authority. The next TASK-RECOVERY must match the
  tuple. Validate every boundary-valid DEC citation in an authority field
  against phase decisions; one real DEC does not mask a nonexistent or embedded
  pseudo-ID. Do not allocate a new TRG for a resolved BLOCK. If an existing DEC
  converts it to control-conformance repair, first record Control repair base
  and have TL-Lite append one unique separate H2 amendment naming the existing
  governing TRG. Parse it independently from the original TRG's retained BLOCK
  verdict; then RUN-RESUME retains that TRG/ESC and binds the amended recovery.
  A normal non-STOP RUN-RESUME records `N/A` governing IDs and `Next recovery`.
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
