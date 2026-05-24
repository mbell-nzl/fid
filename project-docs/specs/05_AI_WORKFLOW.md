# AI WORKFLOW - <PROJECT_NAME>

> This document defines the reusable agentic engineering workflow: roles,
> artifact layout, orchestration, escalation, prompts, context management, and
> continuous improvement.

## 1. Principles

1. **Human defines intent, agents execute.** The human owns product direction,
   phase approval, and final judgment. Agents do implementation, test, review,
   and documentation work.
2. **Agents are disposable.** Each agent session is short-lived and focused.
   Shared state lives in repository files, not chat memory.
3. **Artifacts are the communication layer.** Agents hand off through task
   specs, completion reports, QA plans, reviews, decisions, and progress files.
4. **Deterministic gates before subjective review.** If format, lint, test,
   coverage, build, or contract checks fail, the agent fixes those before any
   review.
5. **Validation assertions connect specs to evidence.** Each phase uses
   `VAL-NNN` assertions so requirements, tasks, QA evidence, and wrap-up
   decisions stay traceable.
6. **Small tasks, tight loops.** One Developer session owns one task. If a task
   cannot fit in one session, the Tech Lead splits it.
7. **Documentation ships with the change.** Specs, README, user docs, and
   workflow docs update with the code they describe.
8. **Specs are living guardrails.** Implementation can reveal a spec is wrong.
   When that happens, update the spec and record the decision.
9. **Current context stays small.** Default-loaded current docs are summaries.
   Historical detail belongs in phase archives.
10. **Trust behavior, verify structure.** Generated tests are useful, but
   quality also needs static gates, independent audits, runtime evidence, and
   hidden scenario validation outside the implementation agent's context.
11. **Do not contaminate held-out scenarios.** Implementation agents must not
   load hidden scenario repositories, excluded scenario paths, or evaluator
   prompts unless the task is explicitly scenario validation.

## 2. Repository Artifact Layout

```text
project-docs/
|-- specs/
|   |-- 00_PROJECT_CHARTER.md
|   |-- GLOSSARY.md
|   |-- 01_PRD.md
|   |-- 02_TRD.md
|   |-- 03_DESIGN_PLAN.md
|   |-- 04_QUALITY_PLAN.md
|   |-- 05_AI_WORKFLOW.md
|   `-- 06_IMPLEMENTATION_ROADMAP.md
|-- docs/
|   |-- current/
|   |   |-- phase_brief.md
|   |   |-- decisions_active.md
|   |   |-- decisions_principles.md
|   |   |-- progress_current.md
|   |   `-- tl_escalations_current.md
|   |-- phases/
|   |   `-- phase-N/
|   |       |-- tasks/
|   |       |-- qa/
|   |       |   |-- validation-contract.md
|   |       |   |-- qa-plan.md
|   |       |   |-- review-M.md
|   |       |   |-- final-review.md
|   |       |   `-- human-checklist.md
|   |       |-- decisions.md
|   |       |-- progress.md
|   |       |-- triage-log.md
|   |       `-- tl-escalations.md
|   `-- reviews/
`-- templates/
```

## 3. Agent Roles

### 3.1 Tech Lead

**Purpose:** Plans implementation phases and wraps them up. Does not write
product code.

**Activated:** By the human at phase start and phase end.

**Phase-start responsibilities:**

- Read the relevant roadmap phase.
- Read relevant PRD/TRD/design/quality sections.
- Read `project-docs/specs/GLOSSARY.md` when naming, terminology, or domain
  language affects task planning.
- Review `docs/current/decisions_active.md`.
- Skim the prior phase's workflow notes for calibration.
- Create phase directories.
- Draft `docs/phases/phase-N/qa/validation-contract.md` from the template
  before or alongside task specs.
- Write ordered task specs under `docs/phases/phase-N/tasks/`.
- Ensure each task cites the `VAL-NNN` assertions it advances or explicitly
  explains why no validation assertion applies.
- Ensure each task can fit in one Developer session.
- Ensure each task names required documentation updates.
- Refresh `docs/current/phase_brief.md`.
- Present the task list for human review.

**Phase-end responsibilities:**

- Read QA Lead final review and human checklist outcome.
- Decide for each recommendation: do now, defer, or cancel.
- Audit final `VAL-NNN` statuses, record deviations, and disposition any
  `FAIL`, `BLOCKED`, or `CARRY-FORWARD` assertions.
- **Spec sync, mandatory and auditable.** Walk every in-phase
  `DEC-NNN`, every Developer completion-report deviation, and every QA Lead
  recommendation marked do now. For each, identify affected specs, edit them
  in this wrap-up, and record a per-doc verdict for all seven specs:
  `touched: see DEC-NNN`, `verified-no-change`, or `N/A`.
  `verified-no-change` is a positive review assertion, not the default.
- **Triage-log audit.** Walk every `TRG-NNN` in the phase triage log. Confirm
  each RESUME verdict was correct in hindsight, flag any that should have
  escalated, and update the TL-Lite warm-up/autonomous-resume state.
- **Workflow notes.** Append `## Workflow notes` and `### Deep retrospective
  recommendation` to phase decisions. Answer the section 10 questions in 1-3
  sentences each.
- Audit whether new or changed terminology requires a `GLOSSARY.md` update or
  flagged ambiguity.
- Update active decisions and progress.
- Move resolved escalations into the phase archive.
- Trim default-loaded current docs.

**Does not do:**

- Write product code or tests.
- Run quality gates.
- Spawn implementation agents.
- Make product decisions without the human when intent is ambiguous.

### 3.2 Phase Runner

**Purpose:** Mechanical orchestrator for one phase. Does not write code or make
decisions.

**Activated:** By the human after approving Tech Lead task specs.

**Responsibilities:**

- Maintain a compact run ledger: completed tasks, verified commits, QA
  reviews run, ESC IDs seen, TRG IDs logged, and warm-up/autonomous-resume
  status.
- Create or resume the phase branch.
- At phase-run start, record the current set of Open `ESC-NNN` IDs as the
  baseline. Before each task and after each review, compare against that
  baseline plus already-triaged IDs. Only newly appearing Open ESC IDs trigger
  TL-Lite triage.
- Verify the phase validation contract exists before task execution.
- Spawn QA Lead first to create the phase QA plan.
- Verify the QA plan uses `project-docs/templates/qa-plan-template.md` or
  contains equivalent required sections.
- Execute task specs in order.
- Spawn one Developer per task.
- Read and validate each completion report.
- Verify commit hash and subject, dependency-budget note, validation coverage,
  and exact quality-gate evidence before marking a task complete.
- Every third completed task, spawn QA Reviewer.
- Route soft blockers through TL-Lite triage.
- Stop on hard blockers.
- Spawn QA Lead for phase-end review after all tasks complete.
- Notify the human when the phase is ready for wrap-up.

**Hard stops:**

- Missing validation contract.
- Missing QA plan.
- Missing completion report.
- Completion report missing required sections.
- Quality gate failure.
- Ad hoc quality-gate substitution.
- Developer blocker report.
- Commit missing or unverifiable.
- Required command is still a placeholder.
- TL-Lite returns `VERDICT: BLOCK`.
- Warm-up mode treats `VERDICT: RESUME` as human-approval stop.

**Soft blockers routed through TL-Lite:**

- Developer concerns, technical debt, or notes for Tech Lead.
- QA Reviewer Blocking finding.
- New open ESC entry in `tl_escalations_current.md`.

### 3.3 Developer

**Purpose:** Implements exactly one task using red/green TDD.

**Activated:** By Phase Runner with a specific task spec.

**Responsibilities:**

- Read `AGENTS.md`, current docs, task spec, and scoped files.
- Read the phase validation contract and track the `VAL-NNN` assertions named
  by the task spec.
- Continue from any existing in-scope partial work.
- Write failing tests before implementation.
- Run the most specific relevant test command and confirm failure.
- Implement the minimum code to pass.
- Refactor with tests green.
- Update impacted docs.
- Run the task-level quality gate exactly as defined in `AGENTS.md`.
- Write a completion report immediately after a green gate.
- Commit with a Conventional Commit subject and `Refs: phase-N/task-M`.

**Completion report required headings:**

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

If there are no concerns, the concerns section must contain exactly:

```text
- None.
```

If no dependency changed, the dependency-budget section must contain exactly:

```text
- No dependency changes.
```

**Does not do:**

- Work on more than one task.
- Skip TDD.
- Delete or weaken tests to pass.
- Make architectural or product decisions.
- Modify out-of-scope files except required documentation.

### 3.4 QA Lead

**Purpose:** Creates phase QA plans, runs phase-end comprehensive review, and
generates human verification checklists. Does not write product code.

**Activated:** By Phase Runner at phase start and phase end, or by the human
for full QA review.

**Phase-start QA plan must include:**

- Output at `docs/phases/phase-N/qa/qa-plan.md`, using
  `project-docs/templates/qa-plan-template.md`.
- PRD/TRD/design requirements covered by the phase.
- Validation contract refinement and coverage strategy.
- Acceptance criteria.
- Expected tests.
- Deterministic checklist items suffixed `(qa-lead)`.
- Human checklist items suffixed `(human)`.
- Command syntax validation section with validation method, result, and
  evidence.

**Phase-end review must include:**

- Output at `docs/phases/phase-N/qa/final-review.md`, using
  `project-docs/templates/final-review-template.md`.
- Phase-level quality gate results.
- Deterministic checklist execution evidence.
- Final outcome for every `VAL-NNN` assertion.
- Acceptance criteria verification.
- Requirement compliance.
- Test quality review.
- Documentation review.
- Security/privacy review where applicable.
- Performance/reliability spot-checks where applicable.
- Findings with severity.
- Separate human checklist containing only `(human)` items.

PASS on deterministic items requires reproducible evidence.

### 3.5 QA Reviewer

**Purpose:** Single-use focused reviewer for every three completed tasks.

**Activated:** By Phase Runner every third task or by the human for a targeted
mid-phase review.

**Responsibilities:**

- Review the specified task specs, completion reports, commits, and changed
  files.
- Verify relevant in-window `VAL-NNN` assertions against the validation
  contract, QA plan, and task evidence.
- Assess architecture, behavior placement, tests, naming, error handling,
  docs, security, and quality-gate plausibility.
- Write `docs/phases/phase-N/qa/review-M.md`.
- Use `project-docs/templates/mid-phase-review-template.md` or equivalent
  sections.
- Classify findings as Blocking or Improvement.
- Escalate ambiguous or cross-cutting concerns to `tl_escalations_current.md`.

**Polish / parity review mode:**

Active when the phase brief identifies a polish, parity, consistency, audit, or
cross-surface phase, or when the reviewed task specs use language such as
`consistent`, `matches`, `parity`, `follows`, or `audit`. Additional outcome
requirements apply:

1. Cite the `03_DESIGN_PLAN.md` sections checked for each task.
2. Run at least one representative deterministic check per task when the stack
   provides one: command, test, browser interaction, API request, rendered
   artifact inspection, or equivalent. Paste or summarize reproducible evidence
   under that task.
3. Walk the phase brief's Parity Matrix where present. Verify each axis across
   the relevant surfaces or outputs. Discrepancies are Blocking by default.
4. Include at least one observation per task. Improvement-class observations are
   acceptable. Zero-finding reviews are not accepted in polish / parity mode.

Phase-specific cross-cutting rules belong in the phase brief or QA plan, not in
the durable QA Reviewer prompt.

**Escalate if the issue:**

1. Has architecture implications beyond the reviewed task window.
2. Exposes ambiguous acceptance criteria.
3. Has unclear security/privacy severity.
4. Contradicts a standing decision.
5. Spans tasks or subsystems beyond the review window.
6. Is a plausible but material spec deviation.
7. Is a quality-gate anomaly not anticipated by the QA plan.

### 3.6 TL-Lite

**Purpose:** Single-use tactical Tech Lead for mid-phase questions and soft
blocker triage. Does not write code or run commands.

**Activated:** By human for a question, or by Phase Runner with
`Triage trigger:`.

**Local-decision scope:**

- Clarifying an existing task spec.
- Resolving a local implementation question inside one task or subsystem.
- Interpreting already-recorded constraints.
- Deciding whether a soft blocker is safe to resume.

**Must escalate if the issue:**

1. Changes a public interface.
2. Changes the data model or persistence contract.
3. Overrides a standing decision.
4. Affects more than one task or subsystem.
5. Introduces a dependency, tool, or workflow.
6. Changes test strategy or coverage policy.
7. Exposes ambiguity in product intent or acceptance criteria.
8. Requires recovery of partial task work after context exhaustion and the
   recovery spans schema/data model, dependencies/platform config, public
   interfaces, canonical docs, or multiple Developer sessions.

In triage mode, TL-Lite appends a TRG entry to
`docs/phases/phase-N/triage-log.md` and ends with exactly one final line:

```text
VERDICT: RESUME
```

or:

```text
VERDICT: BLOCK (escalated as ESC-NNN)
```

## 4. Phase Lifecycle

```text
Human requests phase planning
  -> Tech Lead writes validation contract, task specs, and phase brief
  -> Human reviews and approves
  -> Phase Runner starts phase
  -> QA Lead refines validation contract and writes QA plan
  -> Developer implements task 1
  -> Phase Runner validates completion report
  -> Repeat per task
  -> QA Reviewer reviews every third task
  -> TL-Lite triages soft blockers
  -> QA Lead runs final review
  -> Human completes human checklist
  -> Tech Lead audits validation outcomes, wraps phase, and compresses context
```

## 5. Task Spec Format

Each task spec should include:

```markdown
# Phase N - Task M: <Title>

## Objective
<What this task changes and why.>

## Requirements
- PRD/TRD/design/quality refs:
  - <ID or section>

## Validation Contract Coverage
- VAL-NNN: <one-line assertion this task advances>
- <If none apply, state why explicitly.>

## Scope
### Files to create or modify
- <path>

### Files to read for context
- <path>

### Out of scope
- <explicit non-goal>

## Acceptance Criteria
- [ ] <observable/testable criterion> (VAL-NNN or no VAL)

## TDD Plan
- First failing test:
- Expected implementation path:
- Regression risks:

## Documentation Updates
- <doc path or "None expected">

## Quality Gate
- <FORMAT_COMMAND>
- <LINT_COMMAND>
- <COVERAGE_COMMAND>

## Notes For Developer
- <constraints, ordering, known pitfalls>
```

The phase validation contract lives at
`docs/phases/phase-N/qa/validation-contract.md`. Tech Lead drafts it before or
alongside task specs. QA Lead refines evidence methods and validator
assignments before implementation begins. If a required validation assertion
cannot be satisfied, the responsible agent stops and escalates rather than
silently adapting the spec.

## 6. Completion Report Format

```markdown
# Phase N - Task M: <Title> - Completion Report

## Summary of what was implemented
- <summary>

## Files created or modified
- <path>

## Documentation updated
- <path or "None required">

## Validation Contract Coverage
- VAL-NNN: <evidence produced or status>
- <If none apply, explain why.>

## Tests run
- <command and result>

## Exact quality-gate results
- <FORMAT_COMMAND>: PASS/FAIL
- <LINT_COMMAND>: PASS/FAIL
- <COVERAGE_COMMAND>: PASS/FAIL

## Deviations from spec
- None.

## Concerns, technical debt, or notes for the Tech Lead
- None.

## Dependency-budget note
- No dependency changes.

## Commit hash
- <hash>

## Commit subject
- <subject>
```

## 7. Triage And Escalation

### 7.1 ESC Entries

ESC entries live in `docs/current/tl_escalations_current.md` while open. Use
global `ESC-NNN` numbering that never resets. Move resolved items to the phase
archive at wrap-up.

### 7.2 TRG Entries

TRG entries live in `docs/phases/phase-N/triage-log.md`. Use per-phase
`TRG-NNN` numbering. Tech Lead audits all TRG entries at phase wrap-up.

### 7.3 Warm-Up Mode

For the first two phases using TL-Lite triage, treat every
`VERDICT: RESUME` as a stop pending human approval. After two consecutive
phases where the Tech Lead confirms all RESUME verdicts were correct in
hindsight, record a carry-forward decision enabling autonomous resume.

## 8. Context Loading

Default-loaded context for most agents:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/GLOSSARY.md` when terminology is in scope.

On-demand context:

- `project-docs/docs/current/decisions_principles.md`
- Specific task specs and completion reports.
- Relevant spec sections.
- `project-docs/specs/GLOSSARY.md` when terminology, naming, or ambiguous
  language affects the task.
- Closed phase archives.

Do not load full historical archives unless the task requires them.

## 9. Phase-End Compression Ritual

At every phase wrap-up, Tech Lead performs this checklist:

- [ ] Audit final `VAL-NNN` outcomes and record deviations or carry-forwards.
- [ ] Record spec sync verdict for all seven specs.
- [ ] Audit TRG entries and update warm-up state.
- [ ] Append workflow notes and deep-retro recommendation.
- [ ] Move resolved carry-forward items out of `decisions_active.md`.
- [ ] Append full phase history to `docs/phases/phase-N/progress.md`.
- [ ] Trim `progress_current.md` to current status only.
- [ ] Move resolved escalations to `docs/phases/phase-N/tl-escalations.md`.
- [ ] Update `project-docs/specs/GLOSSARY.md` or flag ambiguities when phase
      work changed project language.
- [ ] Trim default-loaded docs to their size budgets.
- [ ] Refresh `phase_brief.md` for the next phase.
- [ ] Move or delete `tl_handoff.md` if present.

Size budgets:

| File | Target Lines | Tripwire |
|---|---:|---:|
| `phase_brief.md` | 150 | 200 |
| `decisions_active.md` | 150 | 200 |
| `progress_current.md` | 100 | 133 |
| `tl_escalations_current.md` | 80 | 107 |

## 10. Continuous Improvement

### 10.1 Lightweight Notes Every Phase

At wrap-up, Tech Lead appends workflow notes to the phase decisions file:

- What worked well?
- What was inefficient?
- Were tasks the right size?
- Did the QA cadence feel right?
- Were token-efficiency opportunities missed?

### 10.2 Deep Retrospective On Trigger

Recommend a deep retrospective if:

- The same issue recurs across two consecutive phases.
- Human checklist Fail + Partial exceeds 20 percent.
- A human failure maps to a QA Lead PASS on a deterministic-adjacent item.
- The human asks.

The human decides whether to run the deep retrospective. Store it in
`docs/reviews/YYYY-MM-DD-phase-N-retrospective.md`.

## 11. Independent Quality Loops

These checks sit outside the normal Developer -> QA loop and are scheduled by
the roadmap or Tech Lead when the project reaches the relevant quality maturity
level.

- **Automated gates:** stack-appropriate CI checks for format, lint, tests,
  coverage or equivalent confidence, types where available, security scan,
  dependency audit, license check, and secret scan.
- **Held-out scenarios:** user-journey scenarios stored outside normal agent
  context. Run Mode A evaluator checks and at least one Mode B human run before
  release readiness.
- **Independent audits:** architecture, security, maintainability, performance,
  complexity, spec drift, and dependency reviews run in a fresh agent session,
  preferably with a different model family for high-stakes work.
- **Release readiness:** threat model, pre-release checklist, changelog or
  release notes, migration notes, dependency budget review, and rollback or
  recovery check where applicable.

Outputs live under `project-docs/docs/reviews/` unless a phase task spec names a
more specific location. Templates live under `project-docs/templates/`.

Glossary changes are documentation changes. Durable tradeoff decisions that
would otherwise be ADRs are recorded as phase `DEC-NNN` entries and reflected
in specs as needed; do not create a parallel ADR system unless the human
explicitly adopts one.

## 12. Model And Tooling Guidance

### 12.1 Repo Skills

Repo skills are reusable procedural helpers, not workflow authority.
`05_AI_WORKFLOW.md`, role prompts, templates, specs, and recorded decisions
remain the source of truth for role boundaries, hard stops, required artifacts,
and validation gates.

Skill authoring and distribution:

- Canonical skills live under `project-docs/skills/<skill-name>/SKILL.md`.
- Codex-discoverable generated copies live under
  `.agents/skills/<skill-name>/SKILL.md`.
- Claude-discoverable generated copies live under
  `.claude/skills/<skill-name>/SKILL.md`.
- Generated skill copies are derived files and must not be edited directly.
- Workflow-lint must fail when generated copies are missing or drift from the
  canonical skill source.

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
  behavior or external tooling is needed, and document prerequisites, safe
  defaults, non-interactive usage, and dependency rationale.

Prompt and skill boundary:

- Agent prompts keep role identity, activation rules, context loading, write
  boundaries, hard stops, required artifact paths, escalation authority, and
  provider-specific invocation mechanics.
- Skills contain repeatable procedures, template-use checklists,
  evidence-quality heuristics, examples, and references to canonical
  specs/templates.
- A skill must never be the only place that records a hard stop, role authority
  boundary, release gate, validation gate, dependency policy, public-interface
  policy, or architecture decision.

Initial skill set:

- `validation-contract` - workflow-extracted from validation-contract planning
  and refinement rules.
- `qa-evidence` - workflow-extracted from QA evidence and validation-outcome
  rules.
- `phase-wrap-up` - workflow-extracted from Tech Lead phase-end and compression
  rules.
- `tdd` - external-adapted from Matt Pocock's TDD skill, rewritten for task
  specs, validation contracts, and completion reports.
- `diagnose` - external-adapted from Matt Pocock's diagnose skill.
- `zoom-out` - external-adapted from Matt Pocock's zoom-out skill.

Defer `phase-planning` and `task-pack` until generation and workflow-lint are
stable, because those skills sit closer to Tech Lead authority and task
decomposition semantics.

Side-effectful skills require explicit invocation controls where supported.
For Claude, use `disable-model-invocation: true`. For Codex, use
`agents/openai.yaml` with `policy.allow_implicit_invocation: false`. Do not
automate a recurring workflow until it works reliably as a manual skill.

### 12.2 Role Capability Mapping

Default mapping:

| Role | Recommended Capability | Reason |
|---|---|---|
| Tech Lead | Strongest reasoning model | Planning, synthesis, decisions |
| QA Lead | Strongest reasoning model | Phase-wide verification |
| Developer | Strong coding model | Implementation and TDD |
| QA Reviewer | Medium/high reasoning model | Bounded review |
| TL-Lite | Medium/high reasoning model | Focused triage |
| Phase Runner | Fast lower-cost model | Mechanical orchestration |

The exact model names live in `.codex/agents/` and `.claude/agents/`. Update
those files when your toolchain changes.

Codex project-local `.codex/config.toml` is for repository behavior such as
agent references, sandbox defaults, and project-root markers. Provider,
profile, auth, telemetry, and notification routing belong in the user's
`~/.codex/config.toml`, not in the project-local file.

*Last updated: 2026-05-23*
