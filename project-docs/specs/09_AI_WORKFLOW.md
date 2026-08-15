# AI Workflow - fid

**Status:** Binding V1 repository delivery workflow
**Last reviewed:** 2026-08-09

This document is the canonical agentic engineering workflow for
`fid`. Agent prompts, `AGENTS.md`, `CLAUDE.md`, Codex config, and
templates are derivative and must align to this file.

This document describes how fid is scaffolded and built **now**. It is not a
target V2 kit workflow, a fid product model, or a backend-port integration
contract. ADR-0011/0012 and the current PRD, TRD, and roadmap define that
future product direction without changing this workflow. The V1 process
remains binding until the Human approves a separately recorded adoption.

Tooling: Claude Code and Codex.
Orchestration: hybrid - Human at phase gates and true escalations, automated
dev/QA loop available through Phase Runner, artifact-based communication, and
deterministic quality gates.

## 1. Principles

1. Human defines intent, validates outcomes, and resolves product decisions.
2. Agents are disposable; all shared state lives in files.
3. Artifacts are the communication layer. Agents hand off through structured
   files, not chat memory.
4. Default context is intentionally small and lives under
   `project-docs/docs/current/`.
5. Deterministic gates run before subjective review.
6. A Developer session handles at most one task. A task may span same-task
   continuation sessions when durable recovery is required.
7. Documentation ships with the change.
8. Specs are guides and guardrails. When implementation proves a spec wrong,
   update the affected spec and record the decision.
9. Prompt and config files do not create workflow policy by implication.
10. Detection severity and workflow disposition are separate: a real failure
    may require repair without requiring Human attention.
11. Failed attempts remain immutable evidence. Causal repair creates a later
    attempt; it never rewrites failure into PASS.

## 2. Agent Roles

### 2.1 Human

Purpose:
- Own product intent, phase selection, and final approval.
- Approve task packs and resolve product-level questions.
- Run human verification checklists.

Does not do:
- Use chat history as project memory.
- Skip final human verification for planned user-facing work.

### 2.2 Tech Lead

Purpose:
- Plan phases and wrap them up. Does not write implementation code or tests.

When activated:
- By the Human at phase start for planning.
- By the Human at phase end after QA and human verification.
- By the Human for targeted ESC resolution when a mid-phase escalation blocks
  progress.

Phase start:
1. Read `AGENTS.md`, `project-docs/docs/current/*`, relevant specs, and the
   roadmap.
2. Review `project-docs/docs/current/decisions_active.md`; schedule, re-defer,
   or close carry-forward items.
3. Skim the prior phase's `## Workflow notes` in
   `project-docs/docs/phases/phase-<previous-id>/decisions.md` if present. Use
   it for calibration signals such as task sizing, QA cadence, or model
   selection; it binds nothing unless promoted to carry-forward.
4. Propose the phase operating mode (light, normal, or heavy) per section 4
   and record it in the phase brief; the Human approves it with the task list.
   Light mode skips the validation contract and QA plan per the artifact
   matrix.
5. Create `project-docs/docs/phases/phase-<id>/tasks/` and
   `project-docs/docs/phases/phase-<id>/qa/`, and seed an empty
   `project-docs/docs/phases/phase-<id>/triage-log.md` from
   `project-docs/templates/triage-log-template.md`. Seeding it at phase start
   keeps ID allocation off the critical path of a mid-phase agent that must
   otherwise create the file before it can mint an identifier.
6. For normal and heavy modes, draft
   `project-docs/docs/phases/phase-<id>/qa/validation-contract.md` using
   `project-docs/templates/validation-contract-template.md`. The contract is a
   finite checklist of testable behavioral assertions, written before or
   alongside task specs. Each assertion has a per-phase sequential `VAL-NNN` ID
   (resets each phase), requirement source, user-visible behavior, evidence
   method, validator role, task window, pass and fail criteria, and an
   out-of-scope note. The contract is the **spec**; the QA plan is the
   **test plan**. See section 6.1.
7. Break the phase into ordered tasks sized to fit one Developer session. A
   task may use a same-task continuation when an estimate is exceeded or a
   recoverable verification state remains. Each task either satisfies a subset
   of `VAL-NNN` assertions or explicitly states why no validation assertion
   applies.
8. Write task specs to `project-docs/docs/phases/phase-<id>/tasks/task-<n>.md`
   using the task contract in section 6. Each task spec lists the `VAL-NNN`
   assertions it satisfies in its Acceptance Criteria.
9. Ensure each task lists required docs, exact commands, acceptance criteria,
   files that are explicitly out of scope, and the `VAL-NNN` assertions it
   satisfies. In `## Exact Commands To Run`, state when the full task gate is
   required. Unless the task explicitly says otherwise, the full task gate is
   required before task close because the Developer does not load the QA plan.
10. Verify coverage rules: every PRD/TRD requirement in scope is referenced by
   at least one `VAL-NNN`; every `VAL-NNN` is referenced by at least one task
   spec acceptance criterion; every `VAL-NNN` has exactly one validator role.
11. Refresh `project-docs/docs/current/phase_brief.md` using the phase brief
    template. Keep it as an index and execution brief, not a duplicate of full
    task specs. Record the operating mode and reference the validation
    contract path when the mode requires one.
12. Present the ordered task list, operating mode, and validation contract when
    required to the Human before execution.

Phase end:
1. Read `project-docs/docs/phases/phase-<id>/qa/final-review.md` and the
   completed human checklist.
2. Decide each QA finding or recommendation: fix now, carry forward, reject
   with reason, or escalate to the Human.
3. Validation contract audit: confirm every `VAL-NNN` in
   `project-docs/docs/phases/phase-<id>/qa/validation-contract.md` has a
   recorded outcome (`PASS`, `FAIL`, or `CARRY-FORWARD`) in the final review.
   Disposition each `FAIL` and `CARRY-FORWARD`: fix now, carry forward to a
   future phase (re-issue under a new VAL-NNN later), reject with reason, or
   escalate. Record outcomes summary in
   `project-docs/docs/phases/phase-<id>/decisions.md`.
4. Spec sync is mandatory. Walk every in-phase DEC-NNN, every Developer
   completion-report deviation, and every QA Lead "do now" recommendation. For
   the full canonical spec suite, record `touched`, `verified-no-change`, or
   `N/A` in `project-docs/docs/phases/phase-<id>/decisions.md`.
   `verified-no-change` is a positive assertion after review, not a default.
5. Walk `project-docs/docs/phases/phase-<id>/triage-log.md` if present.
   Confirm each TL-Lite RESUME disposition was correct in hindsight. If any TRG
   should have escalated, record the pattern in phase decisions and tighten
   `decisions_active.md`, `decisions_principles.md`, or `AGENTS.md` before the
   next phase.
6. Confirm each TL-Lite RESUME verdict in the triage log was correct in
   hindsight and flag any that should have escalated (the audit in step 5).
   Autonomous TL-Lite resumption is permanently enabled (DEC-002); there is no
   warm-up tally to maintain.
7. Append lightweight `## Workflow notes` to this phase's decisions archive.
   Answer the section 10 calibration questions in 1-3 sentences each, compare
   with the prior phase when available, and end with
   `### Deep retrospective recommendation`.
8. Move resolved carry-forward items and resolved escalations out of current
   files. History lives in the phase archive.
9. Update `project-docs/docs/current/progress_current.md` and append the full
   phase history to `project-docs/docs/phases/phase-<id>/progress.md`.
10. Refresh `project-docs/docs/current/phase_brief.md` for the next active
    phase or mark that no phase is active.
11. Delete `project-docs/docs/current/tl_handoff.md` after it has been
    consumed, or archive it under
    `project-docs/docs/phases/phase-<id>/handoffs/`.

Compression ritual at every phase end:
- Record validation contract outcomes (`VAL-NNN` PASS / FAIL / CARRY-FORWARD).
- Record spec sync verdicts for the full canonical spec suite.
- Audit every TRG-NNN entry.
- Append workflow notes and deep-retro recommendation.
- Move resolved CFs out of `decisions_active.md`.
- Move resolved ESCs out of `tl_escalations_current.md`.
- Append phase progress history and trim `progress_current.md`.
- Trim current-state files to their size targets.
- Refresh the next phase brief.

Does not do:
- Write implementation code or tests.
- Run tests or linting.
- Activate other agents; Phase Runner handles orchestration.
- Make product-level decisions without Human approval.

### 2.3 TL-Lite

Purpose:
- Single-use tactical Tech Lead for one mid-phase question or one Phase Runner
  triage event. Does not write code, tests, task specs, decisions, or phase
  wrap-up artifacts. Does not run commands.

Modes:
- Question mode: answer a local tactical question that does not change specs,
  public interfaces, schema, dependencies, workflow, or product intent.
- Triage mode: invoked by Phase Runner with `Triage trigger:`. Decide
  `RESUME`, `REPAIR`, or `BLOCK`, append a TRG-NNN entry, and end
  with a literal verdict line.

Default context:
- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- the specific file, error, report, review, or excerpt named in the prompt

On-demand context:
- `project-docs/docs/current/decisions_principles.md`
- the specific task spec for the question
- closed phase archives when explicitly referenced
- never full PRD/TRD/roadmap by default

Must escalate if the issue:
1. Changes a public interface, CLI command contract, output format, journal
   event schema, or user-visible workflow.
2. Changes the data model or persistence contract. This includes SQLite
   schema, migrations, journal replay semantics, and storage invariants.
3. Contradicts `decisions_active.md`, `decisions_principles.md`, or
   `AGENTS.md`.
4. Affects more than one task or subsystem.
5. Introduces a dependency, tool, quality gate, or workflow.
6. Changes testing strategy or coverage policy.
7. Exposes ambiguity in product intent or acceptance criteria. Design
   ambiguity that changes behavior also escalates.
8. Requires an unapproved protected change to recover partial task work, such
   as a schema/data-model, dependency/platform, public-interface, or canonical
   workflow change. Context exhaustion and same-task continuation alone are
   `REPAIR`, for any number of continuation sessions.

When in doubt, escalate. Do not optimize for answering immediately when a
trigger is crossed.

Triage rules:
- Identifiers are **phase-scoped**: `TRG-<phase>-<NNN>` and
  `ESC-<phase>-<NNN>`, for example `TRG-19-001` and `ESC-19-001`. The phase
  triage log `project-docs/docs/phases/phase-<id>/triage-log.md` is the single
  allocation source for both series: take the highest `<NNN>` already recorded
  there in the relevant series and add one, starting at `001` in a new phase.
  Phase scoping exists because a global counter has no readable source of
  truth — closed escalations leave the current register, so a global sequence
  cannot be resumed by an agent that cannot scan every phase archive.
- Every triage event appends a `TRG-<phase>-<NNN>` entry to the
  phase-start seeded
  `project-docs/docs/phases/phase-<id>/triage-log.md`. The entry is
  written **before** the verdict is reported; the verdict quotes the
  identifier that was written. A verdict carrying a placeholder identifier
  (`ESC-NNN`, `ESC-NEW`) or unaccompanied by a log entry is a defective
  triage: the event has no durable record, cannot be resumed mechanically, and
  forces a full Tech Lead reconstruction. The Tech Lead seeds the file at phase
  start, so an absent triage log is a phase-setup defect: the Phase Runner
  records it and stops rather than accepting a verdict with no identifier. When
  a triage genuinely cannot be written, the Phase Runner logs
  `TRIAGE (identifier unavailable)` and the full Tech Lead backfills both the
  TRG and the ESC, labelled as reconstructions.
- Triage `RESUME` means the trigger is demonstrably misclassified,
  already resolved by verified evidence, or covered by an accepted decision.
  The Phase Runner may advance from the current checkpoint.
- Triage `REPAIR` means the trigger is real but the causal work is local,
  decision-free, inside the existing task authority, and changes no protected
  decision, gate definition, threshold, command, or assertion strength. The
  Phase Runner must not advance: it appends `TASK-RECOVERY`, dispatches a
  same-task Developer continuation, reruns the required verification, and
  repeats any affected independent review.
- Triage BLOCK appends a concrete `ESC-<phase>-<NNN>` to
  `project-docs/docs/current/tl_escalations_current.md` when no existing ESC
  covers the issue, then logs the TRG entry with BLOCK.
- If the source is an existing `ESC-<phase>-<NNN>`, do not create a duplicate.
  Log the TRG entry with BLOCK pointing to the existing ESC unless the new
  issue is clearly unrelated.
- For a QA `REPAIR` outcome, verify that every must-fix finding names an
  owning task, bounded repair scope, unchanged acceptance criterion, and
  required re-review evidence. If so, return `REPAIR`. If the finding is
  cleared, return `RESUME`. If protected authority is required or
  cross-task work cannot be separated without a new decision, return `BLOCK`.
- A report-lint failure, missing report, typed `RECOVERY_REQUIRED` result,
  in-scope quality-gate failure, missing/skipped command, or incomplete
  implementation is normally `REPAIR`; incomplete work is not a blocker.
- `BLOCK` requires a protected category: product, architecture, data
  model, public interface, dependency approval, security/data loss,
  unrepairable external infrastructure, specification contradiction, or
  exhausted no-progress recovery budget. Red tests, failed gates, report
  defects, and file count are not protected categories by themselves.
- TL-Lite cannot run commands and never verifies commits itself. A commit
  counts as verified only when the invoking prompt attaches
  commit-verification evidence (hash, `git rev-parse --verify` result, and
  subject match, or report-lint output). Without attached evidence, do not
  RESUME on the already-resolved-by-a-verified-commit ground.
- The final line in triage mode must be exactly `VERDICT: RESUME`,
  `VERDICT: REPAIR`, or
  `VERDICT: BLOCK (escalated as ESC-<phase>-<NNN>)`, with the concrete
  identifier substituted.

Write boundaries:
- May append ESC entries to
  `project-docs/docs/current/tl_escalations_current.md`.
- May append/create TRG entries in
  `project-docs/docs/phases/phase-<id>/triage-log.md`.
- Must not write task specs, `decisions_active.md`,
  `decisions_principles.md`, or phase decisions.

### 2.4 Phase Runner

Purpose:
- Mechanical orchestrator for one approved phase. Does not write code, tests,
  task specs, QA artifacts, decisions, or product docs.

Naming:
- "Phase Runner" is this build pipeline's orchestrator agent (`@phase-runner`).
  It corresponds to fid's `workflow_runner` role (see `GLOSSARY.md`) — the same
  orchestration responsibility seen from the build-tool side vs. the fid product
  side. "Phase Runner" is current, not a legacy alias for "Workflow Runner";
  the pipeline will converge on "Workflow Runner" once fid is dogfooded.

Default context:
- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- task specs under `project-docs/docs/phases/phase-<id>/tasks/`
- `project-docs/docs/phases/phase-<id>/decisions.md`
- `project-docs/docs/phases/phase-<id>/qa/validation-contract.md`
- QA plan after QA Lead creates it

Responsibilities:
1. Verify task specs exist and are approved.
2. Create or resume the phase branch. When resuming, verify the repo is on the
   expected branch and previously completed task commits are present.
3. Verify `project-docs/docs/phases/phase-<id>/qa/validation-contract.md`
   exists and is non-empty. If missing, stop and tell the Human the contract is
   required before phase execution. Do not author the contract yourself.
4. Spawn QA Lead first to create or refresh
   `project-docs/docs/phases/phase-<id>/qa/qa-plan.md`. If required child
   agents cannot be spawned, stop with a phase-runner configuration error. Do
   not do child work yourself.
5. Read the validation contract and the QA plan before any task execution.
6. Execute task specs in order with one active Developer child at a time. A
   Developer session handles at most one task; the same task may use a
   same-task continuation after a recoverable return. Never advance while the
   current task outcome is `RECOVERY_REQUIRED`.
7. Use explicit child prompts. Each prompt must restate the required reads,
   task path, report path, exact gate commands, stop conditions, and
   commit/report requirements.
   Codex child spawns must set `spawn_agent.agent_type` exactly to the intended
   project custom agent (`qa_lead`, `developer`, `qa_reviewer`, or `tl_lite`).
   Do not omit `agent_type`, do not use `default`, `worker`, or `explorer`, and
   do not set model overrides for these project roles.
8. After each Developer run, run
   `project-docs/tools/report-lint.sh <report> <task spec>`. It
   deterministically verifies the typed outcome, required headings, task-spec
   command coverage, observed termination, commit identity, and
   committed/clean state. Parse its literal `OUTCOME COMPLETE`,
   `OUTCOME RECOVERY_REQUIRED`, or `OUTCOME BLOCKED` line. A nonzero
   exit is `report-repair`; do not re-derive failed checks by prose parsing.
9. Read the report for the judgment items report-lint cannot check: blocker
   evidence, deviation content, protected concerns, and QA-plan commands
   required for that task window beyond the task spec. Every QA-plan required
   command must be recorded as `PASS`, or as `BLOCKED` with reason in a
   non-complete report. Missing, silently skipped, or narrower substituted
   required commands prevent completion and enter recovery.
10. Route parsed Developer, report-lint, QA, TL-Lite, command-liveness, and
    recovery-budget results through
    `project-docs/tools/workflow-route.sh <source> <outcome> [attempts]` and
    follow its literal `ACTION` result. `WAIT` retains and polls a live command
    session; `EVALUATE` inspects its observed terminated result. Invalid router
    input is an orchestration defect to repair, not permission to invent a
    route.
    Every new RUN-START/RUN-RESUME carries `Sequence contract: v2`; older events
    remain legacy-tolerant. Pipe each exact candidate entry, beginning with its
    single `###` event heading, to
    `project-docs/tools/workflow-sequence-lint.sh <run-ledger> --append-stdin`.
    The tool validates and atomically appends the entry at EOF; an `INVALID`
    result leaves the ledger byte-for-byte unchanged. Require exit zero and
    literal `APPEND_CONFIRMED state=<state>`. If termination or receipt is
    missing, retry the byte-identical candidate; an exact EOF match confirms it
    without another append. Do not advance, submit a different entry, or perform
    a separate `apply_patch`, redirect, or editor mutation. This validator and
    append transaction does not choose
    the next action; `workflow-route.sh` does. Run it without `--append-stdin`
    to read current state, and require literal `VALID state=ADVANCEABLE` before
    an actual advance. It preserves task identity and requires recovery
    completion plus an explicitly linked, differently identified QA re-review.
    A linked re-review may transfer repair ownership to another task only after
    the current owner completes; it cannot overwrite an unresolved owner. An
    `INVALID` candidate is recoverable while the ledger remains valid. A
    missing receipt is an ambiguous result resolved only by byte-identical
    retry; retain and stop only if standalone validation fails after a confirmed
    append.
11. On `OUTCOME COMPLETE`, append TASK complete with report path, lint
    result, and verified commit. On `OUTCOME RECOVERY_REQUIRED` or a
    recoverable lint, report, gate, or session defect, append
    `TASK-RECOVERY` with trigger signature, attempt count, checkpoint, and
    next action, then dispatch a same-task Developer continuation. On
    `OUTCOME BLOCKED`, treat the report as a blocker candidate and triage
    it; never accept a child-authored blocker as authoritative lifecycle state.
12. Re-read `project-docs/docs/current/tl_escalations_current.md` before each
    task and after each review.
13. Maintain the append-only run ledger at
    `project-docs/docs/phases/phase-<id>/run-ledger.md` using
    `project-docs/templates/run-ledger-template.md`: one appended entry per
    event (RUN-START, `RUN-RESUME`, TASK complete, `TASK-RECOVERY`,
    REPAIR, REVIEW, TRIAGE, STOP, RUN-END). The
    ledger, not chat notes, is the durable run state; a resumed runner reads
    it first and reconciles it against completion reports and git history. It
    is the only file the Phase Runner writes.
14. Every third completed task, or at the cadence defined by the QA plan,
    spawn QA Reviewer for the review window with exact task, report, and
    output paths, then append the REVIEW entry. `PASS` advances.
    `REPAIR` dispatches the named owning task to a Developer continuation
    and requires a fresh independent re-review of the same window. `BLOCK`
    is a blocker candidate for TL-Lite.
15. Route blocker candidates and ambiguous recovery classifications through
    TL-Lite triage.
    When a triage may turn on "already resolved by a verified commit", verify
    the commit first and attach the commit-verification evidence (hash,
    `git rev-parse --verify` result, subject match) to the triage trigger;
    TL-Lite judges evidence and never verifies commits itself. After TL-Lite
    returns, read the triage log to confirm the TRG-NNN entry exists, then
    parse the literal final `VERDICT:` line and append a TRIAGE entry to the
    run ledger. `VERDICT: RESUME` advances only when the trigger is
    cleared. `VERDICT: REPAIR` appends `TASK-RECOVERY` and dispatches
    the same task; it never advances. `VERDICT: BLOCK` appends STOP with
    the concrete ESC.
16. Apply a default budget of two consecutive Phase Runner-mediated
    continuation attempts for the same trigger signature with no material
    progress. Same-session red/green iterations do not consume this budget. A
    relevant in-scope change tied to the trigger is progress only when it newly
    satisfies or removes an obligation, or verification demonstrably narrows
    the failure. Formatting and metadata-only or unrelated changes do not
    qualify; no-op churn is not progress. On exhaustion, triage the protected
    `recovery-budget-exhausted` candidate.
17. After all tasks complete, spawn QA Lead for final review and human
    checklist.
18. Append the RUN-END entry to the run ledger and notify the Human with the
    final review path and current stop/completion state. When stopping early,
    append a STOP entry naming the trigger and the next human action.

Recovery conditions:

- `OUTCOME RECOVERY_REQUIRED`.
- Report-lint failure, missing report, incomplete commit/report sequence, or
  missing termination evidence.
- In-scope focused/full-gate failure or obsolete assertion caused by the task.
- Missing/skipped/narrower command record that can be corrected without
  changing the accepted command.
- Child session ending while durable in-scope work remains.
- QA `REPAIR` with an owning task and bounded decision-free repair.
- Recoverable local sandbox, cache, or scratch-state failure.

Protected blocker candidates:

- `OUTCOME BLOCKED` with a closed protected category and supporting
  evidence.
- Placeholder or contradictory required command/specification.
- Product, architecture, data-model, public-interface, dependency, security,
  data-loss, or gate-definition authority is required.
- Unrepairable external infrastructure/capability failure.
- Recovery budget exhausted for an identical no-progress trigger.
- QA `BLOCK`.
- A new Open ESC-NNN appears after the run-start baseline.

Before appending STOP for a blocker candidate, spawn TL-Lite as the first-step triager,
append TRIAGE for its verdict, and stop only on
`VERDICT: BLOCK`. `VERDICT: REPAIR` keeps the current task open and
dispatches recovery; `VERDICT: RESUME` advances only when the trigger is
already cleared.

Immediate safety stops are limited to a branch/worktree state that cannot be
reconciled non-destructively, a missing validation contract that requires
authoring authority, or an inability to spawn any required role after bounded
retry. Record the exact resolution condition. Missing QA output by itself is a
retryable QA-session failure, not a Human blocker.

Autonomous triage and resolved resume:

- Autonomous TL-Lite `RESUME` and `REPAIR` are enabled. Only
  `BLOCK` creates a Human-visible STOP.
- When a prior STOP has a later decision, verified commit, restored capability,
  or other recorded resolution, verify that reference, append `RUN-RESUME`
  linking STOP, TRG/ESC, and resolution, and continue at the unresolved task.
  A prior TRIAGE entry does not prohibit a resolved resume.
- If no TRIAGE entry exists for an existing STOP, invoke TL-Lite before
  resuming.
- Ordinary concerns and technical debt remain in the report for Tech Lead
  wrap-up; they do not trigger triage unless explicitly identified as a
  protected blocker candidate.

### 2.5 Developer

Purpose:
- Implement at most one task per session using strict red/green TDD. Continue
  an assigned task from any durable in-scope checkpoint.

Default context:
- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- the assigned task spec
- `project-docs/docs/phases/phase-<id>/qa/validation-contract.md` (the
  `VAL-NNN` assertions the task is expected to satisfy)
- scoped source, test, and doc files
- for a continuation: the latest `TASK-RECOVERY` ledger entry and owning QA
  review, when applicable

Responsibilities:
1. Confirm task scope, allowed files, and the `VAL-NNN` assertions listed in
   the task spec acceptance criteria.
2. Inspect existing in-scope partial work and continue from that state.
3. Write failing tests first for behavior changes.
4. Implement minimum code to pass, then refactor only while tests stay green.
5. Update docs named in the task spec.
6. Run required quality gates in order. If the task spec lists both a focused
   loop and a full task gate, the full task gate is required before task close
   unless the task explicitly says otherwise. The focused loop is the iteration
   loop and does not substitute for the full gate at close. If the task spec
   does not define a gate, use the standard fallback gate from `AGENTS.md`.
7. If any gate fails, preserve its token and observed termination, diagnose the
   cause, make a causal in-scope repair, and rerun the byte-identical failed
   command. Product defects, stale expectations, harness defects, and
   bookkeeping drift are all repairable when the task owns the change and no
   protected decision or gate definition changes. Never weaken, narrow,
   re-point, substitute, or mine a flaky green.
8. Where the task changes behavior with a runtime surface, exercise the
   changed behavior once through the real entrypoint (the `fid` CLI or other
   real interface) before writing the report, and record the command or
   interaction and its observed result. Tests alone do not prove the
   entrypoint works.
9. Create or update the completion report immediately after gate outcome is
   known, leaving `## Commit hash` and `## Commit subject` as
   `- Pending commit`. Set `Task outcome` to `COMPLETE` only when every
   obligation and gate is satisfied; use `RECOVERY_REQUIRED` for unfinished
   in-scope continuation work and `BLOCKED` only for a protected category.
10. Create the task commit (implementation, tests, docs, and the report draft)
    with a Conventional Commit subject and `Refs: phase-<id>/task-<n>`.
11. While `HEAD` is still the task commit, record the commit hash and subject
    by copying them verbatim from `git` output: `git rev-parse --short HEAD` or
    `git rev-parse HEAD` for the hash, and `git log -1 --format=%s` for the
    subject. Replace the report placeholders with those values. Never type,
    pad, complete, or recall a hash from memory. Retain the Markdown backticks
    when replacing either placeholder; the finalized fields must use this exact
    shape:

    ```text
    ## Commit hash
    - `0123456789abcdef`

    ## Commit subject
    - `feat(scope): subject copied from git`
    ```
12. Commit the finalized report in an immediately following
    `docs(task): finalize completion report for phase-<id> task-<n>` commit.
    The hash in the report is the task commit's, never this docs commit's; a
    report cannot carry the hash of the same commit that introduces its final
    content, so this two-step pattern is mandatory.
13. Run `project-docs/tools/report-lint.sh <report> <task spec>` against the
    committed report. If it fails because the report needs correction, update
    the report, create another `docs(task): ...` commit, and rerun the command
    against the newly committed file.
14. Verify the report path is committed and clean:
    `git status --porcelain -- <report path>` prints nothing. Then stop.

Required completion report headings, in order:
- `## Workflow outcome`
- `## Summary of what was implemented`
- `## Files created or modified`
- `## Documentation updated`
- `## Validation Contract Coverage`
- `## Tests run`
- `## Exact quality-gate results`
- `## Deviations from spec`
- `## Concerns, technical debt, or notes for the Tech Lead`
- `## Dependency-budget note`
- `## Commit hash`
- `## Commit subject`

If there are no concerns, the concerns section contains exactly `- None.`. If no
dependencies changed, dependency-budget note contains exactly
`- No dependency changes.`.

Recovery and blocker behavior:
- Stay in the red/green/repair loop while causal progress is possible. Same
  session iterations do not consume the Phase Runner continuation budget.
- If the session must end with real in-scope work remaining, create a durable
  checkpoint, set `Task outcome: RECOVERY_REQUIRED`, choose the recovery
  category, name the condition and next action, record every failed/not-run
  required command as `BLOCKED`, finalize and lint the report, and return it for
  a same-task continuation. Incomplete work is not a blocker.
- Use `Task outcome: BLOCKED` only for product/architecture/data-model/public-
  interface decision, dependency approval, security/data loss, unrepairable
  external infrastructure, specification contradiction, or exhausted recovery
  budget. Name why current task authority cannot resolve it and the exact
  resolution action.
- Never stop silently. Return a typed report even when no task commit can be
  produced; the Phase Runner classifies the observation and decides lifecycle
  state.

Constraints:
- Never skip TDD for logic changes.
- Never delete or weaken tests to pass.
- Never modify files outside task scope unless the task requires related docs.
- Never make architectural or product decisions; record concerns instead.
- Never hand edit generated files.

### 2.6 QA Lead

Purpose:
- Create phase QA plans and run phase-end comprehensive review. Does not write
  implementation code or fix bugs.

Default context:
- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/docs/phases/phase-<id>/qa/validation-contract.md`
- phase task specs and completion reports
- `qa-plan.md` if present
- mid-phase reviews if present

Validation contract refinement (before finalizing the QA plan):
- Confirm every PRD/TRD requirement in scope is referenced by at least one
  `VAL-NNN`.
- Confirm every `VAL-NNN` has a clear validator role and a concrete evidence
  method that can be turned into a `QA-NN` step or `(human)` item.
- Confirm each `VAL-NNN` describes one observable behavior or invariant. If a
  human assertion necessarily spans several observables, require one
  human-checklist row per observable using sub-IDs that cite the parent
  `VAL-NNN`.
- Confirm conditional or binary assertions require both qualifying and
  non-qualifying states. Positive-only evidence and source reads such as "uses
  flag X" are not sufficient proof.
- Confirm parity or "matches reference" assertions have a deterministic
  conformance proof when the reference can be checked mechanically.
- If gaps exist, return revisions to Tech Lead or raise an ESC. Do not proceed
  with the QA plan until the contract is sufficient.
- Refinements are proposals recorded in the QA plan; the Tech Lead applies all
  edits to the validation-contract file. QA Lead never rewrites the contract
  file itself.

QA plan must include:
- Source list including `validation-contract.md`.
- One-line coverage rollup: all in-scope PRD/TRD requirements referenced by
  >=1 `VAL-NNN` (YES/NO).
- Task inventory (operational risk and review focus).
- Required automated coverage table; rows cite the `VAL-NNN` they cover.
- Required commands (quality gates) table.
- Deterministic checklist items suffixed `(qa-lead)` with reproducible evidence
  expectations. Each `QA-NN` step cites the `VAL-NNN` it verifies.
- Human checklist items suffixed `(human)` only for perception, dogfooding,
  product-owner judgment, tone, or real-device behavior. Each `(human)` item
  cites the `VAL-NNN` it judges.
- One observable per human checklist row. When a `(human)` `VAL-NNN` covers
  several surfaces or affordances, fan it out into one row per discrete
  observable, with sub-IDs that all cite the same parent `VAL-NNN`.
- `## Command And Feasibility Validation` listing each referenced command,
  test scope, fixture setup, artifact assumption, and human-checklist
  feasibility check. Invalid or unvalidated checklist steps are QA-plan
  defects that must be fixed before task execution begins.
- Phase-end decision rule including: every `VAL-NNN` has a recorded outcome
  with evidence.

Routine mid-phase review:
- Handled by QA Reviewer. QA Lead handles mid-phase review only when the Human
  explicitly asks, an ESC requires QA Lead judgment, or the phase plan says so.

Phase-end review must:
- Run or verify every `(qa-lead)` item from the QA plan.
- Include a deterministic checklist execution log in `final-review.md`.
- Record `VAL-NNN` outcomes in `## Validation Contract Outcomes` of
  `final-review.md`. Each row is `PASS` / `FAIL` / `BLOCKED` with cited
  evidence. `PASS` without evidence is invalid.
- Use evidence from test output, source/file state, CLI command output,
  generated artifact state, journal/DB state, or another reproducible proof as
  appropriate.
- Treat PASS without evidence as invalid. If evidence cannot be shown, mark the
  item FAIL or BLOCKED.
- Review requirement compliance, gate results, test quality, architecture and
  layer boundaries, documentation, security/data-loss/privacy risk, and
  carry-forward candidates.
- Generate `project-docs/docs/phases/phase-<id>/qa/final-review.md`.
- Generate `project-docs/docs/phases/phase-<id>/qa/human-checklist.md` with
  only the `(human)` subset, using the human checklist template.
- End with `Outcome: PASS | REPAIR | BLOCK`. `REPAIR` findings must
  name an owning task, bounded decision-free scope, required evidence, and
  fresh independent re-review. `BLOCK` is reserved for a protected
  blocker candidate or unresolved Human-only validation.

### 2.7 QA Reviewer

Purpose:
- Single-use focused mid-phase reviewer. Does not write or fix implementation
  code.

Default context:
- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/docs/phases/phase-<id>/qa/qa-plan.md`
- `project-docs/docs/phases/phase-<id>/qa/validation-contract.md`
- the task specs and completion reports named in the prompt

On-demand:
- `project-docs/specs/05_DESIGN_PLAN.md` for CLI UX or TUI work.
- `project-docs/docs/current/decisions_principles.md` for standing rule checks.
- closed phase archives only when a finding cross-references them.
- never full PRD/TRD/roadmap by default.

Responsibilities:
- Review the named task specs, completion reports, tests, and committed code.
- Inspect code each task changed with `git log`, `git show`, and file reads.
- Verify in-window `VAL-NNN` assertions: from each reviewed task, list the
  `VAL-NNN` it claims to satisfy and assess whether evidence in the commit,
  tests, and report demonstrably tracks toward the pass criterion. Cite the
  `VAL-NNN` and evidence in the per-task section. An assertion that cannot
  be evidenced is a Must-fix finding unless the assertion is explicitly
  scheduled for a later task window.
- Compare each reviewed task's completion report against the task spec's exact
  commands and the QA plan's required commands. Missing required gates, skipped
  gates without a `BLOCKED` reason, or narrower substituted commands are
  Must-fix findings.
- Check architecture, layer boundaries, business logic placement, test quality,
  test-pair compliance, naming, docs, dependency-budget compliance,
  quality-gate reality, and commit/report hygiene.
- Rerun at least one representative required command from the window's task
  specs and record its result in the review as reproduced evidence; expand
  where risk warrants. This is the only independent command execution between
  the Developer and phase end. Targeted commands such as
  `go test ./internal/<scope>/...` or package-scoped linting are appropriate;
  do not run broad mutating or full-gate commands unless explicitly asked.
  `project-docs/tools/report-lint.sh <report> <task spec>` may be run per task
  as independent mechanical evidence.
- Classify findings as Must-fix or Improvement. Severity says whether work must
  change before sign-off; it does not decide whether Human authority is needed.
- Each finding should include severity, file:line where possible, what is
  wrong, and fix direction.
- Write exactly one review artifact:
  `project-docs/docs/phases/phase-<id>/qa/review-<n>.md`.
- For every Must-fix finding, record its owning task, whether the repair stays
  inside accepted scope, the required verification, and the independent
  re-review condition.
- End with `Must-fix: <count>  Improvement: <count>  Blocker candidates:
  <list or none>` and `Outcome: PASS | REPAIR | BLOCK`.
  `PASS` means no Must-fix finding remains. `REPAIR` means every
  Must-fix finding is real, bounded, decision-free, and owned by an existing
  task. `BLOCK` means at least one finding needs protected authority or
  cannot be bounded within the review window. A fresh Reviewer must perform
  re-review after repair; the original Reviewer never resolves its own finding
  by assertion.

CLI/TUI review mode:
- Active when the phase brief identifies CLI UX or TUI work, or task language
  emphasizes consistency, matches, parity, follows, or audit across surfaces.
- Cite `project-docs/specs/05_DESIGN_PLAN.md` sections per reviewed task.
- If the phase brief has a Parity Matrix, walk it. Parity discrepancies are
  Must-fix by default unless an accepted decision says otherwise. They produce
  `REPAIR` when bounded and decision-free, and `BLOCK` only when protected
  authority is required.
- Use focused test output, source-linked behavior proof, CLI invocation
  evidence, or another concrete evidence source. Code reading alone is not
  enough when runnable evidence exists.
- Include at least one observation per reviewed task. Zero-observation reviews
  are not accepted in this mode.

Standing QA escalation rule:
- Most findings stay in `review-<n>.md`.
- QA roles never allocate or append ESC identities. Mark a blocker candidate
  for TL-Lite if the finding:
  1. Is a suspected architecture or layer-boundary violation with implications
     beyond the review window.
  2. Is a test gap that implies the acceptance criterion itself is ambiguous.
  3. Is a security, data-loss, journal-integrity, or privacy concern where
     severity is genuinely unclear.
  4. Contradicts a standing decision in `decisions_principles.md` or
     `AGENTS.md`.
  5. Spans tasks or subsystems beyond the review window.
  6. Is a plausible but material spec deviation.
  7. Is a coverage, vet, lint, or codegen anomaly the phase QA plan did not
     anticipate.

When in doubt, return `BLOCK` with the candidate category and evidence.
Do not stretch a focused review to make a call that warrants QA Lead or Tech
Lead judgment.

### 2.8 Product Designer

Optional specialist for TUI or dashboard work. Use only when visual or
interaction design needs separate attention. Output must feed back into task
specs, QA plan, Design Plan, or review artifacts.

## 3. Orchestration Flow

### 3.1 Phase Lifecycle

```text
Human selects phase
  -> Tech Lead proposes operating mode, plans phase, and writes task specs
  -> Human reviews and approves task pack
  -> Phase Runner starts or resumes the approved phase
  -> Phase Runner spawns QA Lead for QA plan
  -> Phase Runner executes each task through Developer
  -> Phase Runner verifies typed report outcome, gate status, and commit
  -> Recoverable result loops to same-task Developer continuation
  -> QA Reviewer reviews every 3 tasks or planned checkpoint
  -> QA REPAIR loops through owning task and fresh independent re-review
  -> Protected blocker candidates route to TL-Lite
  -> QA Lead runs final review and writes human checklist
  -> Human runs human checklist
  -> Tech Lead wraps phase, syncs specs, archives state
```

Minimum Human time per phase:
- Review task specs before execution.
- Review final QA report and run human checklist.
- Resolve true escalations and product decisions.

### 3.2 Communication Protocol

All inter-agent communication happens via files. The directory layout is:

```text
project-docs/
|-- specs/                             # canonical specs (golden docs)
|   |-- 00_PROJECT_CHARTER.md
|   |-- GLOSSARY.md
|   |-- 01_PRD.md
|   |-- 02_TRD.md
|   |-- 03_TRD_SCHEMA.md
|   |-- 04_TRD_CONFIG.md
|   |-- 05_DESIGN_PLAN.md
|   |-- 06_DESIGN_PLAN_CLI.md
|   |-- 07_DESIGN_PLAN_FLOWS.md
|   |-- 08_QUALITY_PLAN.md
|   |-- 09_AI_WORKFLOW.md
|   `-- 10_IMPLEMENTATION_ROADMAP.md
|-- docs/
|   |-- current/                       # default-loaded current state
|   |   |-- phase_brief.md
|   |   |-- decisions_active.md
|   |   |-- decisions_principles.md    # on demand
|   |   |-- progress_current.md
|   |   `-- tl_escalations_current.md
|   |-- phases/
|   |   `-- phase-<id>/
|   |       |-- tasks/
|   |       |   |-- task-<n>.md
|   |       |   `-- task-<n>-done.md
|   |       |-- qa/
|   |       |   |-- validation-contract.md
|   |       |   |-- qa-plan.md
|   |       |   |-- review-<n>.md
|   |       |   |-- final-review.md
|   |       |   `-- human-checklist.md
|   |       |-- decisions.md
|   |       |-- progress.md
|   |       |-- triage-log.md
|   |       `-- tl-escalations.md
|   `-- reviews/
`-- templates/
```

### 3.3 Git Workflow

- Each phase uses a phase branch such as `phase-1/core-bootstrap`.
- Each task should produce a single Conventional Commit.
- Commit body or footer includes `Refs: phase-<id>/task-<n>`.
- Phase Runner verifies the recorded commit before continuing.
- Phase completion is recorded in
  `project-docs/docs/current/progress_current.md` and archived in
  `project-docs/docs/phases/phase-<id>/progress.md`.

### 3.4 Codex Subagent Triggering

- Codex subagents are explicit; Human or parent agent must request delegation.
- Canonical role prompt sources live in `project-docs/agents/<role>.md`.
  `.claude/agents/` and `.codex/agents/` are generated provider files. Edit the
  canonical source and run `project-docs/tools/sync-agents.sh`; workflow-lint
  fails when generated prompts drift.
- Project Codex agents live in `.codex/agents/`.
- `.codex/config.toml` must set `[features].multi_agent = true`.
- `.codex/config.toml` must set `[agents].max_depth = 3`.
- Phase Runner delegates only to `qa_lead`, `qa_reviewer`, `developer`, and
  `tl_lite`.
- In Codex, Phase Runner must call `spawn_agent` with `agent_type` set exactly
  to that role name. An omitted `agent_type` spawns the default agent, which
  inherits the parent runtime instead of the project custom agent's model and
  reasoning settings.
- Child prompts must be explicit and self-contained. A bare role handoff is not
  enough.
- Project-scoped Codex agent files live in `.codex/agents/*.toml` and carry
  each role's model choices, reasoning effort, and stop conditions. Do not add
  `[profiles.*]` to project-local `.codex/config.toml`; current Codex ignores
  project-local `profiles` keys. Personal profile files, if needed, live under
  `~/.codex/` and are selected with `--profile`.
- Start a fresh Codex session after changing Codex config or agent files.

### 3.5 Repo Skills

Repo skills are reusable procedural helpers, not workflow authority.
`09_AI_WORKFLOW.md`, role prompts, templates, specs, and recorded decisions
remain the source of truth for role boundaries, hard stops, required artifacts,
and validation gates.

Skill authoring and distribution:

- Canonical skills live under `project-docs/skills/<skill-name>/SKILL.md`.
- Codex-discoverable generated copies live under
  `.agents/skills/<skill-name>/SKILL.md`. This is not a courtesy copy: Codex
  natively discovers `.agents/skills/` by searching the working directory and
  parent directories up to the repo root, then `$HOME/.agents/skills`, then
  `/etc/codex/skills` with symlinks followed. The generated tree is the live
  Codex skill surface.
- Claude-discoverable generated copies live under
  `.claude/skills/<skill-name>/SKILL.md`.
- Generated skill copies are derived files and must not be edited directly.
- Workflow-lint must fail when generated copies are missing or drift from the
  canonical skill source.

Skill invocation:

- Codex invokes a skill implicitly when the task matches the skill
  description, or explicitly via `$skill-name` or the `/skills` picker. Codex
  budgets the in-context skill list, name plus description per skill, at 2% of
  the context window or 8,000 characters, so descriptions must front-load
  trigger language.
- Claude invokes skills through its skill surface, matching on the same
  descriptions.
- Per-skill disabling for Codex is user-level `~/.codex/config.toml`
  `[[skills.config]]`, never repository config.

Skill file rules:

- Canonical skills should follow the Agent Skills format: a skill directory
  with `SKILL.md`, required `name`, required `description`, and Markdown body.
- Keep `SKILL.md` focused on one coherent job. Move long references into
  supporting files only when the skill genuinely needs them.
- Descriptions must say what the skill does and when to use it, with clear
  trigger language and boundaries.
- Put `Provenance:` in the Markdown body, not custom frontmatter. Use
  `Workflow-Extracted Skill`, `External-Adapted Skill`, or
  `Reference-Only Pattern` terminology from `GLOSSARY.md`.
- First versions should be instruction-only. Add scripts only when deterministic
  and safe to run non-interactively.

### 3.6 Documentation Layout And Size Governance

Default-loaded current files:

| File | Purpose | Target | Tripwire |
|---|---|---:|---:|
| `phase_brief.md` | Execution brief for active phase | 150 lines | 200 |
| `decisions_active.md` | Active CF register and scheduled commitments | 150 lines | 200 |
| `progress_current.md` | Current phase status only | 100 lines | 133 |
| `tl_escalations_current.md` | Unresolved escalations only | 80 lines | 107 |
| `tl_handoff.md` | Optional incoming TL handoff | 120 lines | 160 |

If a file exceeds its tripwire, Tech Lead prunes it at the next phase boundary.
When a phase completes, full history moves to
`project-docs/docs/phases/phase-<id>/`.

Per-phase run state:
- `project-docs/docs/phases/phase-<id>/run-ledger.md` is the Phase Runner's
  append-only run ledger (see 2.4), created from
  `project-docs/templates/run-ledger-template.md`. It stays in the phase
  archive; it is never loaded as default context.

Escalations:
- `tl_escalations_current.md` carries unresolved escalations only; it is not a
  running log.
- New ESC entries use phase-scoped `ESC-<phase>-<NNN>` numbering allocated from
  that phase's seeded `triage-log.md`. TL-Lite or the full Tech Lead is the
  allocator; QA roles emit blocker candidates only.
- Resolved entries move to
  `project-docs/docs/phases/phase-<id>/tl-escalations.md`; the current
  projection retains unresolved entries only.

Phase brief required sections:
- Goal
- Phase status
- Operating mode
- Task index pointing to full specs
- Active constraints
- Interfaces in scope
- Known risks
- Open questions
- Success/exit criteria
- Next checkpoint
- Parity Matrix when work involves CLI UX polish, consistency, matching, or
  parity across surfaces

Triage log:
- Created on first TL-Lite triage event in a phase.
- Append-only audit trail.
- TRG numbering is per phase and resets each phase.
- Walked by Tech Lead at phase wrap-up.

Decision numbering:
- DEC numbers are globally sequential across all phases. A phase decisions file
  is a filing shard, not a namespace.

Carry-forward classification:
- `decisions_active.md` contains live CF items and scheduled commitments.
- `decisions_principles.md` contains standing rules enforced across multiple
  phases.
- Start strict: keep items in `decisions_active.md` until they clearly deserve
  promotion.

On-demand loading:
- Phase archives are not loaded by default. Read them only when a task, review,
  decision, or escalation points there.

## 4. Operating Modes

Every phase declares one operating mode - light, normal, or heavy - proposed
by the Tech Lead at planning time, approved by the Human with the task list,
and recorded in the phase brief. Modes vary the required artifact set, never
the evidence standard.

Mode-invariant rules never turn off:
- Exact-command gates and failed-gate semantics.
- Evidence rules (`PASS`/`BLOCKED` recording, both-states proof for
  conditional assertions).
- Commit discipline (Conventional Commits, `Refs:`, hashes copied from git
  output).
- Glossary discipline.
- Every hard stop.

Promotion rule: if a hard-stop condition or a violated entry criterion
surfaces mid-phase, promote the phase one mode and re-plan the remainder;
never push through in the lower mode. Record the promotion in the phase brief
and phase decisions.

### 4.1 Light Mode - Supervised Session

Entry criteria (all must hold):
- Documentation, tooling, or configuration change only. This is a deliberate
  initial restriction: widen to other half-day, single-subsystem work only by
  recorded decision after observed clean light-mode usage.
- Scope fits in about half a day and touches a single subsystem.
- No schema, public-interface, dependency, or security/privacy change.
- The Human is present and orchestrating directly, not AFK.

Required: a task list in `phase_brief.md`, or task specs for non-trivial
items; one completion note per task with summary, files changed, commands run
with results, and commit hash plus subject; commits with `Refs:`. The full
task gate still runs, and the Human reviews the diff or dogfoods immediately.
Specs affected by the change are updated with the change itself.

Skipped: validation contract, QA plan, mid-phase reviews, final review, and
Phase Runner.

### 4.2 Normal Mode - Default For Feature Phases

Entry criteria (any): multi-task phase; behavior changes carrying product
intent; the Human intends AFK execution for some tasks.

Required: validation contract; one task spec per task; QA plan - command and
feasibility validation results may carry forward from the prior phase when the
commands are unchanged, and changed or new commands must be revalidated;
completion reports in the full format; final review with `VAL-NNN` outcomes;
human checklist before wrap-up. Mid-phase review when the phase exceeds four
tasks or touches shared surfaces; otherwise the QA plan may record a lighter
cadence.

Per-phase spec sync stays full-scope in normal mode: at wrap-up, Tech Lead
records `touched`, `verified-no-change`, or `N/A` for the full canonical spec
suite.

### 4.3 Heavy Mode - Release-Adjacent, Security-Sensitive, Novel-Domain

Entry criteria (any):
- Data migration, public-interface, or persistence-contract change.
- Security or privacy surface.
- Release readiness.
- The first 2-3 phases of a new project, or the first phase in a new domain.
  Command and feasibility validation earns most of its value here, so results
  never carry forward in heavy mode.
- Anything the threat model flags.

Required: the full workflow - validation contract, QA plan with fresh command
and feasibility validation, mid-phase cadence per the QA plan, final review,
human checklist, spec sync across all seven specs, plus the maturity-level
artifacts (threat model, held-out scenarios, independent audits) that the
quality plan schedules. Include at least one Mode B human scenario run where
held-out scenarios exist.

### 4.4 Artifact Matrix

| Artifact | Light | Normal | Heavy |
|---|---|---|---|
| Validation contract | - | required | required |
| QA plan | - | required; feasibility may carry forward when commands unchanged | required; fresh feasibility validation |
| Task specs | task list in phase brief, or per-task specs | one spec per task | one spec per task, full detail |
| Completion reports | trimmed note | full format | full format |
| Mid-phase review | - | if more than 4 tasks or shared surfaces | per QA plan |
| Final review | - | required | required |
| Human checklist | immediate diff review or dogfood | required | required, plus Mode B scenario run |
| Spec sync | with the change itself | all seven specs at wrap-up | all seven specs at wrap-up |
| Orchestration | Human direct | Phase Runner or Human | Phase Runner |

## 5. Context Loading

Every agent starts with `AGENTS.md` or `CLAUDE.md` in Claude Code. Additional
loads:

| Agent | Default Load | On Demand | Writes |
|---|---|---|---|
| Tech Lead | `project-docs/docs/current/*`, workflow doc | Roadmap, PRD/TRD, principles, phase archives | task specs, phase brief, validation contract, decisions, current state, specs, handoff |
| TL-Lite | `project-docs/docs/current/*`, trigger source | principles, specific task, referenced archives | append TRG; append ESC on BLOCK |
| Phase Runner | `project-docs/docs/current/*`, phase task list, validation contract, QA plan after creation | triage log, reports, reviews | append-only run ledger; git branch operations |
| Developer | `project-docs/docs/current/*`, assigned task, validation contract, scoped files | principles, referenced archives | code, tests, docs, report, commit |
| QA Lead | `project-docs/docs/current/*`, quality plan, validation contract, phase tasks/reports | PRD/TRD sections, principles, archives | QA plan (refines contract), final review (records VAL outcomes), human checklist |
| QA Reviewer | `project-docs/docs/current/*`, QA plan, validation contract, named tasks/reports, quality plan | Design Plan, principles, archives | review artifact only; blocker candidates are routed to TL-Lite |

## 6. Task Specification Contract

Task specs live at `project-docs/docs/phases/phase-<id>/tasks/task-<n>.md` and
include:

- Objective
- Validation contract assertions satisfied (`VAL-NNN` list)
- Binding constraints
- Scope: files to create/modify, files to read, and out of scope
- Acceptance criteria (each criterion cites the `VAL-NNN` it advances)
- First failing test
- Documentation updates
- Exact commands to run
- Notes for Developer

Requirement references live in the validation contract rows; task specs cite
`VAL-NNN` instead of duplicating PRD/TRD reference lists. Carry-forward items
closed by the task, dependencies on prior tasks, interpretation decisions, and
output requirements are recorded under `## Binding Constraints`,
`## Objective`, or `## Notes For Developer` as appropriate.

Older task specs may contain expanded sections. New task specs use the slim
template unless the phase brief explicitly requires an additional section for
parity, audit, or release-readiness work.

Template: use `project-docs/templates/task-spec-template.md`, whose headings
are (exactly as written there):

```markdown
# Phase <N> - Task <M>: <Title>

## Objective

<2-3 sentences: what changes and why.>

## Validation Contract Coverage

- VAL-NNN: <one-line assertion this task advances>
- <If none apply, state why explicitly.>

## Binding Constraints

- <interpretation decisions, ordering rules, contract semantics that bind
  this implementation>

## Scope

### Files To Create Or Modify

### Files To Read For Context

### Out Of Scope

## Acceptance Criteria

- [ ] <observable/testable criterion> (VAL-NNN or no VAL)

## First Failing Test

<one line; non-obvious regression risks only>

## Documentation Updates

- <doc path or "None expected">

## Exact Commands To Run

- <FORMAT_COMMAND>
- <LINT_COMMAND>
- <TEST_COMMAND>
- <COVERAGE_COMMAND>

## Notes For Developer

- <constraints, ordering, known pitfalls, escalation tripwires>
```

A task spec that prescribes an evidence procedure also states two things the
Developer cannot infer:

- **Required environment capabilities** — `/tmp`, a writable `.git/`, network,
  elevation. A shortfall is discovered and escalated at step zero, never
  halfway through a sequence.
- **Scratch state is unique per attempt.** Temporary directories, worktrees,
  and binaries are allocated with `mktemp -d` or an equivalent fresh path, so no
  run can depend on a prior run's cleanup. A fixed scratch path in a procedure
  is a spec defect.

Prefer expressing an observable as an automated test over a manual procedure.
A validation contract that records a manual evidence procedure states why a
test cannot carry the assertion.

Required-command rules are stated once in `AGENTS.md` (Commands); the task
spec's `## Exact Commands To Run` section points there rather than restating
them at length. Scenario, parity, and audit obligations are phase-level
concerns declared in the phase brief (Parity Matrix) and QA plan; parity or
heavy-mode phases may add task-spec sections for them when the phase brief
requires it.

The canonical integration command targets the dedicated real-process package,
`go test -tags integration -race ./cmd/fid/integration/...`. Ordinary
in-process command and persistence tests run under the default race/coverage
gate. Expanding the integration command to `./...` is not additional evidence;
it duplicates the ordinary suite under a build tag and is a gate-definition
regression.

### 6.1 Validation Contract

`project-docs/docs/phases/phase-<id>/qa/validation-contract.md` is a sibling
artifact to task specs and the QA plan. The contract is the **spec** (what must
be true); the QA plan is the **test plan** (how we prove it). Tech Lead authors
the contract during phase planning before task specs are finalized. QA Lead
refines it during QA plan creation. Developer task specs cite the `VAL-NNN`
they satisfy. QA Reviewer and QA Lead verify outcomes against the contract.

Assertion fields and numbering live in
`project-docs/templates/validation-contract-template.md`. IDs are `VAL-NNN`
per-phase sequential and reset each phase. A carried-forward assertion is
re-issued under a new `VAL-NNN` in the receiving phase with a
`Carry-forward source` note.

Coverage rules (verified before phase execution):
- Every PRD/TRD requirement in scope is referenced by at least one `VAL-NNN`.
- Every `VAL-NNN` is referenced by at least one task spec acceptance criterion.
- Every `VAL-NNN` has at least one QA plan verification step (`QA-NN` or
  `(human)`) and a single validator role.

Authoring rules:
- Each `VAL-NNN` should describe one observable behavior, invariant, or
  release-readiness claim.
- Conditional or binary assertions must name both qualifying and
  non-qualifying states in their pass/fail criteria. Positive-only evidence is
  not sufficient.
- Parity or "matches reference" assertions should include a structural
  conformance or parity proof whenever the reference can be checked
  deterministically.
- Human assertions are reserved for perception, product-owner judgment,
  dogfooding, tone, or flow coherence. If a human assertion necessarily spans
  several observables, the contract notes that QA Lead must fan it out into one
  human-checklist row per observable.

Outcomes (recorded at phase end in `final-review.md`):
- `PASS` - assertion satisfied with cited evidence.
- `FAIL` - assertion not satisfied; blocking unless explicitly carried forward.
- `CARRY-FORWARD` - assertion deferred with Tech Lead rationale at wrap-up.

`PASS` without evidence is invalid.

## 7. Agent Prompt Materialization

The role definitions in section 2 are canonical behavior for all tools.

- Claude Code materializes them in `.claude/agents/*.md`.
- Codex materializes them in `.codex/agents/*.toml`.
- Codex custom agents in `.codex/agents/*.toml` must stay behaviorally aligned
  with this workflow, including model choices, role boundaries, stop
  conditions, and artifact requirements. Project-local `.codex/config.toml`
  must not define `[profiles.*]`; personal profile files belong under
  `~/.codex/` when an operator needs a named top-level launch profile.
- Templates in `project-docs/templates/` must not contradict the workflow.
- **A provider-shared prompt body must not assert a capability limit that only
  one provider's frontmatter enforces.** A shared sentence such as "you cannot
  run commands" reads on the provider that does enforce it as a description of
  its own tool restrictions, and on the provider that does not as a broader
  prohibition — including "cannot write files". State capability limits where
  they are enforced; keep the shared body about role behavior.
- After changing any of these files, start a fresh Codex or Claude session
  before relying on the updated definitions.

Model guidance:

| Role | Codex | Reasoning | Notes |
|---|---|---:|---|
| Tech Lead | `gpt-5.6-sol` | xhigh | phase planning, spec sync, wrap-up |
| Phase Runner | `gpt-5.4-mini` | xhigh | mechanical orchestration, ledger discipline |
| Developer | `gpt-5.6-terra` | high | implementation and TDD |
| QA Lead | `gpt-5.6-sol` | xhigh | phase-wide QA synthesis |
| QA Reviewer | `gpt-5.6-terra` | high | bounded mid-phase review |
| TL-Lite | `gpt-5.6-terra` | medium | focused tactical triage |

Note: this table is Codex-specific. Claude Code agent model choices are
configured separately, per agent, in `.claude/agents/*.md`, and are not
duplicated here to avoid drift. The **Reasoning** column (xhigh / high / medium)
is the tool-agnostic signal — map it to the equivalent capability tier in
whichever toolchain runs the role.

## 8. Git And Quality Gates

Use the task-spec command list when it exists. Otherwise use the `AGENTS.md`
standard commands.

Quality gate commands are defined in `AGENTS.md`. Replace every placeholder
before implementation starts.

Task-level default gate:

```bash
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Conditional:
- `go generate ./...` when generated inputs change.

Run codegen when generated inputs change. Run additional linting when requested,
when lint debt is in scope, or when the phase QA plan enforces it.

Every required command listed by a task spec or QA plan must be recorded in the
completion report as `PASS`, or as `BLOCKED` with reason in a blocker report.
Missing, silently skipped, or narrower substituted commands are failed gates.

Each recorded result carries one short verbatim token copied from the
command's own terminal output — its last printed line, or `exit 0 (no output)`
for a silent success. **A command that did not terminate is never `PASS`.** A
hung, interrupted, or killed command produces no final line to copy; recording
one is fabricated evidence, and recording `PASS` without it makes the false
claim invisible to every downstream review until someone re-runs the gate
independently. Such a command is `BLOCKED` with the reason in a blocker
report. This rule is mode-invariant.

For an execution tool that can yield a long-running command, wrapper
completion is not command completion. The executing role preserves the full
result object, retains any nested session handle, polls that session until the
nested result reports an exit code or explicit termination signal, and records
the output token only then. A downstream reviewer or orchestrator rejects a
termination claim that is supported only by completion of the outer wrapper,
empty output, or absence of new output.

The recorded result carries that termination result alongside the token —
`(terminated exit_code: 0)`, or the explicit signal where the execution
surface exposes no exit code — so termination is visible in the artifact
rather than trusted. A printed last line proves the command produced output;
for a silent success the observed exit status is the only evidence separating
it from a discarded session. An `exit_code` asserted rather than copied from
the command's own result is fabricated evidence, exactly like a fabricated
token. A control that requires more of a completion report than this contract
and the completion-report template do is a defect in the control: the
canonical contract governs, and the Tech Lead fixes the divergence rather than
failing a compliant report.

The token is the command's real exit behavior, not the expected one. A search
or inventory command (`rg`, `grep`) used to prove an *absence* — no production
evaluator, no forbidden call site — terminates with **exit 1 and no output**
when it finds nothing, which is the pass condition for that assertion. Record
`PASS (no match) - exit 1 (no output)` and name the assertion; `exit 0 (no
output)` is reserved for a silent success and is fabricated evidence here. A
task spec or QA plan that relies on such an inventory states its expected exit
status, so the Developer never has to infer one.

A required-command failure is an immutable failed attempt, not automatically a
failed task or Human blocker. When its causal repair is inside accepted task
scope, changes no protected decision, gate definition, threshold, command, or
assertion strength, the Developer repairs it and reruns the byte-identical
command. This general rule includes product defects, stale expectations,
task-authored harnesses, and bookkeeping inventories. Record a mandatory
deviation note with the first attempt's token and termination, root cause,
repair, and later result. The final pass must be deterministic and reproducible
from committed artifacts.

Changing what a gate checks, deleting or weakening a test, narrowing or
re-pointing a command, altering a threshold/budget/tolerance, or mining a flaky
green is not repair. It requires the authority that owns the gate or
specification. DEC-138's bookkeeping-only carve-out is superseded by this
general causal-repair rule; its anti-gaming constraints remain.

The distinction matters at triage. **"The task is incomplete, stop"** and
**"the evidence is tainted, restart"** are different verdicts with different
remedies, and a harness bug is normally the first. A required-command list
states what must be true of a task's **final, reproducible evidence**; it is not
a prohibition on iteration, and reading it as one would demand that any task
shipping its own evidence tooling get that tooling right on the first attempt.
Where provenance of a captured artifact is genuinely in doubt, settle it by
**reproducing the artifact from its pinned inputs and comparing digests** —
reproduction, not chronology or narrative, is what establishes provenance.

**A gate result not captured in the session that records it is `UNVERIFIED`,
never `PASS`.** A result attested from memory, transcribed from another agent's
report, or promised against a named future control is `UNVERIFIED`, and
`UNVERIFIED` is not a passing state. Naming a later verifier is not a
substitute: a control that fires late enough to be tripped over by accident is
not a control. Copy the token in the session that runs the command, or record
`UNVERIFIED` and enter recovery. A later session must run the exact
command itself; transcription never becomes PASS.

**"Pre-existing" is a provenance claim and needs provenance evidence.**
Declining an out-of-scope repair is often correct, but the report must cite the
introducing commit — `git log -S`, `git log --diff-filter=A`, or `git blame` —
and show that it predates the phase-start commit. `git show HEAD` proves a line
is committed, not that it is inherited; HEAD contains the current phase.
An unverified "pre-existing" claim travels through triage and escalation intact
and can conceal a regression the same task introduced.

**Cleanup is never evidence.** A command whose only purpose is removing or
pruning scratch state is excluded from required-command evidence: never
reported `PASS` or `BLOCKED`, never gating a sequence, and a failed cleanup is
a note rather than a stop. Cleanup proves nothing about the observable it sits
next to; fencing it inside an evidence set lets one failed `rm` block a whole
procedure.

**A required command is executed by its author before it is published.** A
command written into a task spec or QA plan as required evidence and never once
run successfully by anyone is not a specification, it is a guess. Where a
command's result depends on the working directory, the spec names that
directory or the command is written to be directory-independent; bare relative
paths in a committed manifest are correct and must be paired with a `cd` into
the directory they resolve against.

`project-docs/tools/report-lint.sh <report> <task spec>` performs the
mechanical subset of completion-report verification deterministically:
typed workflow outcome, required headings, task-spec command coverage, observed
termination, hash validity, subject match, and committed/clean report state.

A completion report cannot contain the hash of the exact commit that
introduces the report's final content - changing the report changes the hash.
The finalized, hash-bearing report is therefore committed in an immediately
following `docs(task): ...` commit after the task commit it describes. A report
is evidence only once committed; lint verdicts recorded against an uncommitted
worktree file do not count toward the Definition of Done.

`project-docs/tools/report-lint.sh <report> <task spec>` also verifies that
the linted report file is committed with no uncommitted modifications.
Orchestrators (Phase Runner or human) run report-lint instead of prose-parsing
these checks; reviewers may run it as independent evidence. Judgment checks -
QA-plan command applicability, deviation content, blocker acceptability -
remain with agents and humans.

`project-docs/tools/drills/run-drills.sh` runs the evidence-chain drills:
seeded-defect fixtures that verify the QA chain itself still catches
fabricated hashes, wrong hash prefixes, skipped or narrower substituted
commands, weakened tests, unrecorded deviations, stale doc claims, and
positive-only `VAL` passes. Recovery drills also cover typed reports,
premature blocker claims, same-task continuation, QA rework, and protected
stops. Run quarterly, after any agent prompt change, and before trusting a new
model in an orchestrator or QA role.

## 9. Stop And Escalate

Enter `BLOCK` and surface the issue instead of guessing only when:

- Required task, validation, or canonical authority is missing or contradictory
  and cannot be authored inside the current role.
- Task boundaries require protected re-slicing or cross-task authority.
- Required gates cannot be repaired inside scope without weakening or changing
  their accepted definition.
- Another session's edits cannot be reconciled non-destructively.
- Implementation reveals a spec conflict or product decision.
- A change would alter SQLite schema, journal event contracts, CLI command
  contracts, output format contracts, dependency policy, or workflow policy.
- A required external capability or child role remains unavailable after
  bounded recovery.
- Two consecutive same-trigger continuation attempts make no material progress.

Red tests, vet/lint failures, missing reports, interrupted sessions, incomplete
in-scope work, and bounded QA findings are `REPAIR` by default. In Phase
Runner mode they dispatch a same-task continuation and preserve failed
evidence. They become `BLOCK` only when a separate protected category or
exhausted recovery budget is demonstrated. `RESUME` means the trigger is
already cleared and it is safe to advance; `VERDICT: REPAIR` means do not
advance.

## 10. Continuous Improvement

Workflow improvement runs on a two-tier cadence: lightweight every phase, deep
on trigger.

### 10.1 Lightweight Workflow Notes

At phase wrap-up, Tech Lead appends `## Workflow notes` to
`project-docs/docs/phases/phase-<id>/decisions.md`. Answer each question in
1-3 sentences:

- What worked well?
- What was inefficient?
- Were tasks the right size?
- Did the QA review cadence feel right?
- Did TL-Lite triage make correct calls?
- Did QA evidence match the plan?
- Did specs drift?
- Should any recurring decision move from `decisions_active.md` to
  `decisions_principles.md` or `AGENTS.md`?
- Were there token-efficiency or context-loading opportunities missed?

`No notable change` is acceptable when true. These notes are observations, not
commitments. If a note must bind future work, promote it to
`project-docs/docs/current/decisions_active.md`,
`project-docs/docs/current/decisions_principles.md`, or `AGENTS.md`.

### 10.2 Deep Retrospective Recommendation

End the workflow notes with:

```text
### Deep retrospective recommendation
NOT RECOMMENDED
```

or:

```text
### Deep retrospective recommendation
RECOMMENDED - <reason>
```

Recommend a deep retrospective when any trigger applies:

- The same workflow theme appears in this phase and the prior phase.
- Human verification returns Fail + Partial on more than 20% of items.
- Any `(human)` failure maps to a `(qa-lead)` item marked PASS.
- The Human asks for one.

The Tech Lead only recommends. The Human decides whether to run the deep
retrospective and who participates. Deep retrospective artifacts live under
`project-docs/docs/reviews/YYYY-MM-DD-phase-<id>-retrospective.md`.

### 10.3 Why Two Tiers

A deep retrospective on every phase produces rote output and trains readers to
skim past real signal. Lightweight notes every phase generate data; deep
retrospectives happen only when there is something concrete to learn from.

## 11. Planned Enforcement Migration (Future — Not Yet In Effect)

**This section changes nothing about the current pipeline.** The workflow
above remains the markdown-based process used to build fid, exactly as
written, until fid is adopted as this project's state backend and the Human
approves the switch.

Recorded direction (DEC-059..DEC-064, 2026-07-18; reaffirmed by
ADR-0011/0012): the mechanical
enforcement this document currently places in role prompts and shell tools —
completion-report commit discipline and report-lint, gate-evidence citation
rules, "required command missing = failed gate" checks, escalation-ID
uniqueness, stop-policy bookkeeping — is the class of rule that migrates into
fid/runner behavior (execution receipts, structured completion obligations,
machine-readable rejections, envelope stop conditions) once the corresponding
fid capabilities ship and are dogfooded. Protected decisions — product,
architecture, data-model, public-interface, and release judgments — remain
with the Human and are never delegated to automated enforcement. When
adoption happens, it will be a separate recorded decision with its own
workflow document; this file is not updated incrementally toward it.

*Last updated: 2026-08-09*
