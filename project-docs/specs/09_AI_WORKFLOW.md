# AI Workflow - fid

This document is the canonical agentic engineering workflow for
`fid`. Agent prompts, `AGENTS.md`, `CLAUDE.md`, Codex config, and
templates are derivative and must align to this file.

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
6. Tasks stay small enough for one Developer session.
7. Documentation ships with the change.
8. Specs are guides and guardrails. When implementation proves a spec wrong,
   update the affected spec and record the decision.
9. Prompt and config files do not create workflow policy by implication.

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
   `project-docs/docs/phases/phase-<id>/qa/`.
6. For normal and heavy modes, draft
   `project-docs/docs/phases/phase-<id>/qa/validation-contract.md` using
   `project-docs/templates/validation-contract-template.md`. The contract is a
   finite checklist of testable behavioral assertions, written before or
   alongside task specs. Each assertion has a per-phase sequential `VAL-NNN` ID
   (resets each phase), requirement source, user-visible behavior, evidence
   method, validator role, task window, pass and fail criteria, and an
   out-of-scope note. The contract is the **spec**; the QA plan is the
   **test plan**. See section 6.1.
7. Break the phase into ordered tasks that each fit one Developer session and
   either satisfy a subset of `VAL-NNN` assertions or explicitly state why no
   validation assertion applies.
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
  `RESUME` or `BLOCK`, append a TRG-NNN entry, and end with a literal verdict
  line.

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
8. Resumes a context-exhausted or partial task where remaining or completed
   work spans schema/data model, dependency/platform configuration, public
   interfaces, canonical docs, or multiple Developer sessions.

When in doubt, escalate. Do not optimize for answering immediately when a
trigger is crossed.

Triage rules:
- Every triage event appends TRG-NNN to
  `project-docs/docs/phases/phase-<id>/triage-log.md`, creating it from
  `project-docs/templates/triage-log-template.md` if needed.
- Triage RESUME is allowed only when the issue is local-decision-eligible and
  no escalation is needed.
- Triage BLOCK appends ESC-NNN to
  `project-docs/docs/current/tl_escalations_current.md` when no existing ESC
  covers the issue, then logs TRG-NNN with BLOCK.
- If the source is an existing ESC-NNN, do not create a duplicate. Log TRG-NNN
  with BLOCK pointing to the existing ESC unless the new issue is clearly
  unrelated.
- For QA Reviewer Blocking findings, RESUME is allowed only when the finding is
  demonstrably misclassified, already resolved by a verified commit, or covered
  by an accepted decision. Otherwise BLOCK.
- TL-Lite cannot run commands and never verifies commits itself. A commit
  counts as verified only when the invoking prompt attaches
  commit-verification evidence (hash, `git rev-parse --verify` result, and
  subject match, or report-lint output). Without attached evidence, do not
  RESUME on the already-resolved-by-a-verified-commit ground.
- The final line in triage mode must be exactly `VERDICT: RESUME` or
  `VERDICT: BLOCK (escalated as ESC-NNN)`.

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
6. Execute task specs in order by spawning exactly one Developer per task,
   unless that Developer produces a blocker report.
7. Use explicit child prompts. Each prompt must restate the required reads,
   task path, report path, exact gate commands, stop conditions, and
   commit/report requirements.
   Codex child spawns must set `spawn_agent.agent_type` exactly to the intended
   project custom agent (`qa_lead`, `developer`, `qa_reviewer`, or `tl_lite`).
   Do not omit `agent_type`, do not use `default`, `worker`, or `explorer`, and
   do not set model overrides for these project roles.
8. After each Developer run, run
   `project-docs/tools/report-lint.sh <report> <task spec>`. It
   deterministically verifies required headings, sentinel sections, task-spec
   command coverage, commit hash validity via `git rev-parse --verify
   <hash>^{commit}`, and commit subject match against `git log -1 --format=%s`.
   A nonzero exit is a hard stop; do not re-derive these checks by prose
   parsing. A hash that does not resolve is a hard stop even when its
   abbreviated prefix is valid.
9. Read the report for the judgment items report-lint cannot check: blocker
   semantics, deviation content, concerns needing triage, and QA-plan commands
   required for that task window beyond the task spec. Every QA-plan required
   command must be recorded as `PASS`, or as `BLOCKED` with reason in a
   blocker report. Missing, silently skipped, or narrower substituted required
   commands are hard stops.
10. Append the TASK entry to the run ledger after each verified task: report
    path, report-lint result, verified commit.
11. Re-read `project-docs/docs/current/tl_escalations_current.md` before each
    task and after each review.
12. Maintain the append-only run ledger at
    `project-docs/docs/phases/phase-<id>/run-ledger.md` using
    `project-docs/templates/run-ledger-template.md`: one appended entry per
    event (RUN-START, TASK complete, REVIEW, TRIAGE, STOP, RUN-END). The
    ledger, not chat notes, is the durable run state; a resumed runner reads
    it first and reconciles it against completion reports and git history. It
    is the only file the Phase Runner writes.
13. Every third completed task, or at the cadence defined by the QA plan,
    spawn QA Reviewer for the review window with exact task, report, and
    output paths, then append the REVIEW entry to the run ledger.
14. Route soft blockers through TL-Lite triage. When a triage may turn on
    "already resolved by a verified commit", verify the commit first and
    attach the commit-verification evidence (hash, `git rev-parse --verify`
    result, subject match) to the triage trigger; TL-Lite judges evidence and
    never verifies commits itself. After TL-Lite returns, read the triage log
    to confirm the TRG-NNN entry exists, then parse the literal final
    `VERDICT:` line and append a TRIAGE entry to the run ledger.
15. After all tasks complete, spawn QA Lead for final review and human
    checklist.
16. Append the RUN-END entry to the run ledger and notify the Human with the
    final review path and current stop/completion state. When stopping early,
    append a STOP entry naming the trigger and the next human action.

Hard stop conditions:
- Developer quality gate fails. Do not send red tests, lint/static-analysis,
  formatting, or coverage failure to TL-Lite.
- Completion report fails `project-docs/tools/report-lint.sh`.
- Developer completion report omits a required task-spec or QA-plan command,
  records it as skipped without a `BLOCKED` reason, or substitutes a narrower
  command for a required gate.
- Developer fails to produce a completion or blocker report.
- Completion report lacks required headings, commit hash, or commit subject.
- Recorded commit hash does not resolve via `git rev-parse --verify
  <hash>^{commit}` to a single commit matching the report subject.
- A blocker report says the task could not be completed, the gate cannot pass,
  a commit cannot be produced, or scope conflicts prevent progress.
- Scope conflict, missing artifact, child-agent unavailability, or branch issue
  prevents safe execution.
- TL-Lite returns `VERDICT: BLOCK`.

Soft blocker conditions:
- Developer report contains concerns, technical debt, or notes for Tech Lead.
- QA Reviewer reports Blocking findings.
- A new Open ESC-NNN appears in
  `project-docs/docs/current/tl_escalations_current.md` after the run-start
  baseline or after the last triage.

TL-Lite resume:
- Autonomous TL-Lite resumption is permanently enabled (DEC-002). On
  `VERDICT: RESUME`, Phase Runner records the TRG-NNN and continues without a
  Human stop. `VERDICT: BLOCK` still always stops the run.
- The former WF-10 warm-up safeguard is retired; no warm-up tally is maintained.

### 2.5 Developer

Purpose:
- Implement exactly one task per session using strict red/green TDD.

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
7. If any gate fails, fix the issue and rerun the exact failed command. Do not
   substitute ad hoc coverage commands for the required gate.
8. Where the task changes behavior with a runtime surface, exercise the
   changed behavior once through the real entrypoint (the `fid` CLI or other
   real interface) before writing the report, and record the command or
   interaction and its observed result. Tests alone do not prove the
   entrypoint works.
9. Create or update the completion report immediately after gate outcome is
   known.
10. Commit with a Conventional Commit subject and `Refs: phase-<id>/task-<n>`.
11. Record the commit hash and subject by copying them verbatim from `git`
    output: `git rev-parse --short HEAD` or `git rev-parse HEAD` for the hash,
    and `git log -1 --format=%s` for the subject. Never type, pad, complete, or
    recall a hash from memory.
12. Stop after report and commit are complete.

Required completion report headings, in order:
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

Blocker behavior:
- If the task cannot be completed, the quality gate cannot pass, a commit
  cannot be produced, or a scope conflict prevents progress, still write
  `project-docs/docs/phases/phase-<id>/tasks/task-<n>-done.md` as a blocker
  report before stopping.
- The blocker report must include current state, changed files, commands run,
  exact blocker, gate status, whether a commit exists, and all required
  headings.
- Never stop silently without either a completion report or blocker report.

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
  be evidenced is a Blocking finding unless the assertion is explicitly
  scheduled for a later task window.
- Compare each reviewed task's completion report against the task spec's exact
  commands and the QA plan's required commands. Missing required gates, skipped
  gates without a `BLOCKED` reason, or narrower substituted commands are
  Blocking findings.
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
- Classify findings as Blocking or Improvement.
- Each finding should include severity, file:line where possible, what is
  wrong, and fix direction.
- Write exactly one review artifact:
  `project-docs/docs/phases/phase-<id>/qa/review-<n>.md`.
- End with `Blocking: <count>  Improvement: <count>  Escalations: <list or
  none>` and `VERDICT: CONTINUE | BLOCK`.

CLI/TUI review mode:
- Active when the phase brief identifies CLI UX or TUI work, or task language
  emphasizes consistency, matches, parity, follows, or audit across surfaces.
- Cite `project-docs/specs/05_DESIGN_PLAN.md` sections per reviewed task.
- If the phase brief has a Parity Matrix, walk it. Parity discrepancies are
  Blocking by default unless an accepted decision says otherwise.
- Use focused test output, source-linked behavior proof, CLI invocation
  evidence, or another concrete evidence source. Code reading alone is not
  enough when runnable evidence exists.
- Include at least one observation per reviewed task. Zero-observation reviews
  are not accepted in this mode.

Standing QA escalation rule:
- Most findings stay in `review-<n>.md`.
- Also append ESC-NNN to
  `project-docs/docs/current/tl_escalations_current.md` if the finding:
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

When in doubt, escalate. Do not stretch a focused review to make a call that
warrants QA Lead or Tech Lead judgment.

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
  -> Phase Runner verifies report, gate status, and commit
  -> Phase Runner routes soft blockers to TL-Lite
  -> QA Reviewer reviews every 3 tasks or planned checkpoint
  -> Phase Runner routes review blockers to TL-Lite
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

### 3.5 Documentation Layout And Size Governance

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
- New ESC entries use globally sequential ESC-NNN numbering.
- Resolved entries move to
  `project-docs/docs/phases/phase-<id>/tl-escalations.md` at phase wrap-up.

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
| TL-Lite | `project-docs/docs/current/*`, trigger source | principles, specific task, referenced archives | ESC append, TRG append/create |
| Phase Runner | `project-docs/docs/current/*`, phase task list, validation contract, QA plan after creation | triage log, reports, reviews | git branch operations only |
| Developer | `project-docs/docs/current/*`, assigned task, validation contract, scoped files | principles, referenced archives | code, tests, docs, report, commit |
| QA Lead | `project-docs/docs/current/*`, quality plan, validation contract, phase tasks/reports | PRD/TRD sections, principles, archives | QA plan (refines contract), final review (records VAL outcomes), human checklist |
| QA Reviewer | `project-docs/docs/current/*`, QA plan, validation contract, named tasks/reports, quality plan | Design Plan, principles, archives | review artifact, ESC append |

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

Required-command rules are stated once in `AGENTS.md` (Commands); the task
spec's `## Exact Commands To Run` section points there rather than restating
them at length. Scenario, parity, and audit obligations are phase-level
concerns declared in the phase brief (Parity Matrix) and QA plan; parity or
heavy-mode phases may add task-spec sections for them when the phase brief
requires it.

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
- After changing any of these files, start a fresh Codex or Claude session
  before relying on the updated definitions.

Model guidance:

| Role | Codex | Reasoning | Notes |
|---|---|---:|---|
| Tech Lead | `gpt-5.5` | high | phase planning, spec sync, wrap-up |
| Phase Runner | `gpt-5.4-mini` | medium | mechanical orchestration, ledger discipline |
| Developer | `gpt-5.5` | high | implementation and TDD |
| QA Lead | `gpt-5.5` | high | phase-wide QA synthesis |
| QA Reviewer | `gpt-5.5` | medium | bounded mid-phase review |
| TL-Lite | `gpt-5.5` | medium | focused tactical triage |

Note: this table is Codex-specific. Claude Code agent model choices are
configured separately, per agent, in `.claude/agents/*.md`, and are not
duplicated here to avoid drift. The **Reasoning** column (high / medium) is the
tool-agnostic signal — map it to the equivalent capability tier in whichever
toolchain runs the role.

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

`project-docs/tools/report-lint.sh <report> <task spec>` performs the
mechanical subset of completion-report verification deterministically:
required headings, sentinel sections, task-spec command coverage, hash
validity, and subject match. Orchestrators (Phase Runner or human) run it
instead of prose-parsing those checks; reviewers may run it as independent
evidence. Judgment checks - QA-plan command applicability, deviation content,
blocker acceptability - remain with agents and humans.

`project-docs/tools/drills/run-drills.sh` runs the evidence-chain drills:
seeded-defect fixtures that verify the QA chain itself still catches
fabricated hashes, wrong hash prefixes, skipped or narrower substituted
commands, weakened tests, unrecorded deviations, stale doc claims, and
positive-only `VAL` passes. Run quarterly, after any agent prompt change, and
before trusting a new model in an orchestrator or QA role.

## 9. Stop And Escalate

Stop and surface the issue instead of guessing when:

- Required task, QA, current-state, or canonical docs are missing.
- Task boundaries are wrong or need re-slicing.
- Required format, vet, test, codegen, or lint gates fail and stay unresolved.
- Another session's in-progress edits conflict with assigned scope.
- Implementation reveals a spec conflict or product decision.
- A change would alter SQLite schema, journal event contracts, CLI command
  contracts, output format contracts, dependency policy, or workflow policy.
- A child agent required by Phase Runner cannot be spawned.

Hard failures such as red tests, vet failures, missing reports, and scope
conflicts stop immediately. Soft blockers route through TL-Lite only in Phase
Runner mode.

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

*Last updated: 2026-07-05*
