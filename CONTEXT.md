# fid

A project-scoped CLI that provides a structured, queryable knowledge graph over the artifacts an AI-assisted development workflow produces. Rows hold metadata and identity; markdown files hold prose.

## Project family

**fid**:
A tapered hand tool used to open up tightly-laid rope so the strands can be worked separately. Pairs with **whet** (sharpening) as a workshop-drawer set — both are tools the author maintains for AI-assisted development. fid is the tracker; whet is a separate sibling project that fid was originally designed to support.

## Language

### Core entities

**Issue**:
Any row in the unified `issues` table — an epic, task, or subtask. The umbrella term, not a synonym for "ticket".
_Avoid using "issue" colloquially to mean "task"._

**Epic**:
A top-level container of related tasks with a defined scope (e.g. a feature, a tech-debt cluster, a bootstrap phase). `kind='epic'` in the `issues` table. Epics do not nest in v1.

**Sequential epic**:
An epic whose child tasks are executed in strict explicit `sequence` order, with real dependencies still expressed separately as `blocked_by` links. `epic_metadata.scheduling = 'sequential'`. `fid ready --epic <id>` returns the lowest-sequence unfinished task only when it is unblocked; if that task is blocked, the epic has no ready task. The bootstrap pattern; what whet calls a "phase". The word _Phase_ is whet's vocabulary and is not used inside fid.

**Sequence**:
A positive integer ordering value scoped to sibling issues under the same parent. For v1 it defines task order inside a [[Sequential epic]]; it is scheduling metadata, not a dependency.

**Prioritised epic**:
An epic whose child tasks are ranked by priority and pulled when ready. `epic_metadata.scheduling = 'prioritised'`. `fid ready --epic <id>` returns the highest-priority unblocked task within the epic. Suits feature work, tech-debt clusters, ongoing initiatives.
_Avoid_: "Reactive epic" (implies improvisation, but the work is still planned).

**Task**:
A unit of work — the rough equivalent of a Linear/Jira "issue". `kind='task'` in the `issues` table. Has acceptance criteria and produces task completions. May belong to an epic or stand alone.

**Subtask**:
A sub-part of a task. `kind='subtask'` in the `issues` table. Uses a simplified workflow. Belongs to exactly one task, not directly to an epic.

**Link**:
A typed traceability relationship between supported project artifacts. In v1, links are scoped to issues, requirements, decisions, and escalations; lifecycle relationships with dedicated fields are not duplicated as links.

### Execution and audit

**Gate**:
A CLI command that produces a deterministic pass/fail signal (test runs, linters, coverage thresholds, format checks). Developers run gates during a task attempt; results land in `task_completions.gate_results` as JSON `(cmd, exit_code, ran_at, output_excerpt)`. A failing gate is a hard stop — TL-Lite cannot resolve it; the Developer fixes and re-runs.

**Requirement**:
A pointer to a PRD/TRD source section, not a copy of the requirement text. Stored as `requirements` rows with source document, source section, summary label, and drift hash.

**Decision**:
A recorded senior choice that future work should respect. Stored as a `decisions` row with full prose in markdown; may resolve an [[Escalation]].

**Comment**:
A generic human- or agent-authored note attached to an [[Issue]], [[QA review]], [[Decision]], [[Escalation]], or [[Requirement]] (`comments` table with scoped polymorphic ref). Use for context, commentary, or observations that fall outside built-in structured concepts such as [[Concern]], [[Finding]], and acceptance criteria.
_Avoid_: using comments for content that should trigger workflow, carry status, be enforced, or feed reports; model that as a structured entity instead.

**Activity log**:
A projected, queryable timeline of meaningful project facts (`activity_log` table): state transitions, AC verifications, claim acquires, escalations raised, etc. Derived from journal events during projection; not independently authoritative. Distinct from [[Comment]]: comments are *narrative*; activity-log entries are *facts*.

### Roles and who-does-what

The role enum (`developer | qa_reviewer | qa_lead | tl_lite | tl | workflow_runner | human`) shows up in four semantically distinct positions. Each has its own term:

**Actor role**:
The role that performed (or is performing) an action. Stored on per-event `_role` columns in append-only history (`verifier_role`, `reviewer_role`, `developer_role`, `raised_by_role`, `author_role`) and projected as `actor` on [[Activity log]] entries. Passed via `--by <role>` on every write CLI invocation. A fact about the past.

Acceptance-criterion verification is multi-author: the same AC may have independent outcomes recorded by Developer, QA Lead, TL/Tech Lead, Human, or another configured verifier role. Workflow transitions and validation contracts specify which verifier roles are required for a given scope.

**Owner role / owner actor**:
Long-lived accountability for an [[Issue]] across its lifetime. Stored on `issues.owner_role` and `issues.owner_actor`. `owner_actor` is optional — NULL means "any actor of `owner_role`". Set at creation; changes are recorded in the journal and projected into [[Activity log]]. Distinct from [[Claim holder]].
_Avoid_: "responsible role" (RACI overloaded), "assignee" (Linear-style).

**Claim holder**:
The role currently holding the exclusive-access lease on an entity. Stored on `claims.claimed_by`. Short-lived (minutes, TTL-bounded). A right to *advance* the work right now. Does not transfer ownership.

**Required validator**:
A planned verifier role/evidence obligation for a [[Validation assertion]]. Stored as one or more `validation_assertion_validators` rows, each pairing the assertion with a required validator role and evidence method. A statement about a plan, not a recorded fact. The corresponding `val_assertion_outcomes.recorded_by_role` is the actor-role fact of who actually recorded the outcome.

### Raised quality signals

Four adjacent concepts that all flag "something is not right" but at different points in the workflow and from different roles. Keep them distinct.

**Concern**:
A structured flag raised by the Developer during a task completion. One row in the `concerns` table per concern; prose lives at `body_md_path`. Status moves `open → addressed | escalated | wontfix`. An open concern triggers [[Workflow Runner|workflow_runner]] to spawn [[TL-Lite]] for [[Triage]].
_Avoid_: "issue" (overloaded with [[Issue]]), "flag" (vague).

**Finding**:
A structured discovery made by a QA Reviewer or QA Lead during a [[QA review]]. One row in `qa_findings`. Has `severity` (`blocking` | `improvement`) and status. May escalate.
_Avoid_: using "finding" for developer-raised flags — those are [[Concern]]s.

**Escalation**:
An unresolved situation that requires a senior decision. One row in `escalations`. Open until linked to a resolving `decision`. Created from a [[Finding]], a [[Concern]] (via TL-Lite BLOCK verdict), or directly by a senior role.

**Carry-forward**:
Open work scheduled into a future epic, or promoted to a project principle. One row in `carry_forward`. Authored by TL at wrap-up. Multi-epic lifespan, unlike the per-attempt scope of a [[Concern]].

### Adjudication and intake

**Triage**:
A recorded TL-Lite adjudication of a mid-task concern. One row per verdict in `triage_events`. The verdict is `RESUME` or `BLOCK`; the source is typically a developer-raised [[Concern]]. This is *not* a state — it's an event. The word "triage" in fid always means this.
_Avoid_: using "triage" to mean the issue intake/inbox state — that's [[Inbox]].

**Inbox**:
The unsorted-incoming state of an issue. `issues.in_inbox = 1` means the issue has arrived but hasn't been assigned to an epic or prioritised yet. Linear-style intake.
_Avoid_: "triage state", "triage queue" — overloads with [[Triage]] which is a different concept here.

### Roles in the enum

**TL (Tech Lead)**:
The strategic senior role; the canonical "orchestrator" in field terms. Invoked at epic kickoff (planning, task decomposition, AC authoring, validation-contract drafting) and at epic wrap-up (QA findings review, golden-doc sync, triage-log audit, lightweight workflow retro, decision finalisation). Rarely invoked mid-flight when an escalation requires Opus judgement. Backed by a capable model. Never triages [[Concern]]s directly — those go through [[TL-Lite]]. Paired with [[TL-Lite]] the way [[QA Lead]] is paired with [[QA Reviewer]]: Opus for bookends, Sonnet for routine mid-flight work.

**TL-Lite**:
The mid-flight Sonnet-class tactical role; single-use per invocation. Two modes:
- **Triage mode** — invoked by [[Workflow Runner]] on a soft blocker (open [[Concern]], QA Blocking [[Finding]], or new Open [[Escalation]]). Ends response with literal `VERDICT: RESUME` or `VERDICT: BLOCK (escalated as ESC-<hash>)` so Workflow Runner can parse the decision.
- **Question mode** — invoked by the human for one tactical question (task-spec clarification, interpretation of an existing decision, local call inside a task or subsystem).
Either mode escalates (writes to [[Escalation]]s) when the issue crosses scope: public-interface change, schema change, contradicts standing decision, multi-subsystem, etc. Never resolves escalations — those flow to [[TL]].

**QA Lead**:
The Opus-class quality role. Authors the [[QA plan]] at epic kickoff and runs the final epic [[QA review]] with deterministic checklist execution. Generates the human verification checklist. Rarely invoked mid-flight on QA-Reviewer escalations. Paired with [[QA Reviewer]] the way [[TL]] is paired with [[TL-Lite]].

**QA Reviewer**:
The Sonnet-class quality companion to [[QA Lead]]. Single-use; invoked by [[Workflow Runner]] every three tasks for mid-epic code review. Writes findings; raises blocking findings as [[Escalation]]s when scope crosses (see §2.6 in whet's workflow doc for the seven escalation triggers).

**Developer**:
The implementing role. Single-use per task. Reads task spec, writes code via TDD, runs gates, writes a structured completion report including the mandatory `## Concerns` section. Raises [[Concern]]s on items the [[TL]] should adjudicate.

**Workflow Runner**:
The mechanical loop driver. Drives an epic's task loop, spawns [[Developer]] / [[QA Reviewer]] / [[QA Lead]] / [[TL-Lite]] agents, dispatches on completion-report inspection and review findings. Intentionally implemented as a cheaper model (Sonnet-class) — its decisions are state-machine deterministic, not strategic. May later become a deterministic harness instead of an LLM (see open thread).
_Avoid_: "Phase Runner" (legacy; renamed when "phase" → [[Sequential epic]]). "Orchestrator" is misleading because field convention puts the capable thinking model in the orchestrator role; that's the TL here.

**Human**:
The product/approval role and the escape valve. Approves epic closure, runs the human verification checklist, invokes [[TL-Lite]] in Question mode for tactical asks, owns product-level decisions. Always allowed to override role boundaries via `--by human`.
