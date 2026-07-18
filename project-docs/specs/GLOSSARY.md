# GLOSSARY - fid

> This document defines project language. It is not an architecture spec,
> product requirements document, implementation guide, or decision log.

## Purpose

Use this glossary to keep domain, workflow, and operational language
consistent across specs, task specs, code, tests, QA artifacts, and human
review. Prefer the terms here when writing new artifacts.

When terminology changes during a phase, update this file or add a flagged
ambiguity before implementation relies on the term.

## Usage Rules

- Define project-specific terms, not general programming vocabulary.
- Prefer one term for one concept. Put rejected aliases in `Avoid`.
- Keep definitions short and implementation-free.
- `Definition`, `Use when`, and `Avoid` are required. `Related` and
  `Examples` are optional — include them when they clarify boundaries
  between nearby terms.
- Use `project-docs/specs/02_TRD.md` for architecture and module boundaries.
- Use `DEC-NNN` entries in phase decisions for durable tradeoff decisions.
- Pre-scaffold design decisions live in `docs/adr/` and predate the phase
  system.
- If a term is ambiguous and cannot be resolved locally, flag it below and
  escalate through the normal Tech Lead/TL-Lite path.

## Project Identity

**fid**:
A local-first, agent-first issue tracker for AI-assisted software
development. SQLite-backed CLI. A tapered hand tool used to open up
tightly-laid rope so the strands can be worked separately. Pairs with
**whet** (sharpening) as a workshop-drawer set — fid is the tracker; whet
is a separate sibling project.

## Core Entities

**Issue**:
Any row in the unified `issues` table — an epic, task, or subtask. The
umbrella term, not a synonym for "ticket".
Use when: referring to any tracked work item regardless of kind.
Avoid: "ticket", using "issue" colloquially to mean "task".

**Epic**:
A top-level container of related tasks with a defined scope (e.g. a
feature, a tech-debt cluster, a bootstrap phase). `kind='epic'` in the
`issues` table. Epics do not nest in v1.
Use when: grouping related tasks under a shared scope.
Avoid: "project", "milestone".

**Sequential epic**:
An epic whose child tasks are executed in strict explicit `sequence` order,
with real dependencies still expressed separately as `blocked_by` links.
`epic_metadata.scheduling = 'sequential'`. `fid ready --epic <id>` returns
the lowest-sequence unfinished task only when it is unblocked; if that task
is blocked, the epic has no ready task. The bootstrap pattern; what whet
calls a "phase".
Use when: describing ordered task execution within an epic.
Avoid: "Phase" (whet's vocabulary, not used inside fid).

**Sequence**:
A positive integer ordering value scoped to sibling issues under the same
parent. For v1 it defines task order inside a [[Sequential epic]]; it is
scheduling metadata, not a dependency.
Use when: describing explicit ordering of tasks within an epic.
Avoid: confusing with [[Link]] dependencies.

**Prioritised epic**:
An epic whose child tasks are ranked by priority and pulled when ready.
`epic_metadata.scheduling = 'prioritised'`. `fid ready --epic <id>` returns
the highest-priority unblocked task within the epic. Suits feature work,
tech-debt clusters, ongoing initiatives.
Use when: describing priority-based task selection within an epic.
Avoid: "Reactive epic" (implies improvisation, but the work is still
planned).

**Task**:
A unit of work — the rough equivalent of a Linear/Jira "issue".
`kind='task'` in the `issues` table. Has acceptance criteria and produces
task completions. May belong to an epic or stand alone.
Use when: describing a single implementable work item.
Avoid: "issue" when you specifically mean a task.

**Subtask**:
A sub-part of a task. `kind='subtask'` in the `issues` table. Uses a
simplified workflow. Belongs to exactly one task, not directly to an epic.
Use when: decomposing a task into smaller tracked parts.
Avoid: "child task" (ambiguous with task-under-epic).

**Link**:
A typed traceability relationship between supported project artifacts. In
v1, links are scoped to issues, requirements, decisions, and escalations;
lifecycle relationships with dedicated fields are not duplicated as links.
Use when: describing cross-entity relationships outside parent-child
hierarchy.
Avoid: duplicating lifecycle fields (owner, status) as links.

## Execution and Audit

**Gate**:
A CLI command that produces a deterministic pass/fail signal (test runs,
linters, coverage thresholds, format checks). Developers run gates during a
task attempt; standalone `fid dev gate` writes append-only `gate_runs` rows.
Finish/block commands snapshot the attempt's gate runs into
`task_completions.gate_results` as JSON `(cmd, exit_code, ran_at,
output_excerpt)`. A failing gate is a hard stop — TL-Lite cannot resolve it;
the Developer fixes and re-runs.
Use when: describing deterministic quality checks during implementation.
Avoid: "quality gate" when referring to subjective review.

**Started-without-claim**:
An issue in a `started` status category with no live claim after the configured
maintenance grace window. The window starts at the latest append-ordered
`status_changed` activity, falling back to issue creation when no transition
exists. A maintenance warning and handoff-attention signal, not structural
corruption.
Use when: identifying started work that has remained unclaimed beyond the
allowed handoff gap.
Avoid: warning during the grace window or resetting age after a non-status
issue update.

**Stale verification**:
An acceptance-criterion verification whose non-NULL `completion_id` references
a superseded task completion. A maintenance warning that the evidence belongs
to an abandoned attempt.
Use when: auditing attempt-bound AC evidence after a newer completion supersedes
the referenced attempt.
Avoid: using the term for NULL/current completion references, old timestamps,
or content-hash drift.

**Telemetry gap**:
A task completion with no `agent_runs` row explicitly linked through
`agent_runs.completion_id`. A maintenance warning that completion-level cost or
token statistics may under-count; an issue-only run is not sufficient linkage.
Use when: auditing whether a completion has attributable run telemetry.
Avoid: treating the warning as DB corruption or inferring linkage from matching
issue ID/timestamp proximity.

**Clock-relative read**:
A read whose value is a function of `(projected state, wall-clock now)` rather
than of derived state alone, e.g. `maintenance check` `info.recent_activity_7d`
(a `now`-anchored window) and `info.claim_count_live` (TTL liveness at `now`).
Because `now` is an input the journal does not reproduce, a clock-relative read
is not a rebuild-parity observable and is excluded from delete-and-rebuild
equality checks.
Use when: reasoning about which outputs rebuild must reproduce exactly.
Avoid: asserting equality of a clock-relative read across two invocations
separated by real elapsed time, or treating its movement as lost/duplicated
projection.

**Requirement**:
A pointer to a PRD/TRD source section, not a copy of the requirement text.
Stored as `requirements` rows with source document, source section, summary
label, and drift hash.
Use when: tracing work back to a source specification.
Avoid: duplicating requirement text into fid's storage.

**Decision**:
A recorded senior choice that future work should respect. Stored as a
`decisions` row with full prose in markdown; may resolve an [[Escalation]].
Use when: recording a durable choice that constrains future work.
Avoid: "ADR" as a separate document system (pre-scaffold ADRs live in
`docs/adr/`; new decisions use `DEC-NNN`).

**Comment**:
A generic human- or agent-authored note attached to an [[Issue]],
[[QA review]], [[Decision]], [[Escalation]], or [[Requirement]] (`comments`
table with scoped polymorphic ref). Use for context, commentary, or
observations that fall outside built-in structured concepts such as
[[Concern]], [[Finding]], and acceptance criteria.
Use when: adding narrative context that doesn't trigger workflow.
Avoid: using comments for content that should carry status, be enforced, or
feed reports — model that as a structured entity instead.

**Activity log**:
A projected, queryable timeline of meaningful project facts (`activity_log`
table): state transitions, AC verifications, claim acquires, escalations
raised, etc. Derived from journal events during projection; not
independently authoritative.
Use when: querying what happened and when.
Avoid: confusing with [[Comment]] (narrative) or journal (raw events).

**Change digest**:
The aggregated, entity-grouped delta between an anchor point and now
(`fid changes --since <anchor>`). For each entity changed since the anchor it
reports net field changes (`from → to`, with intermediate transitions
collapsed) and newly-added child facts (AC verifications, concerns,
completions, links, and later findings, comments, and governance facts).
Derived from [[Activity log]] entries inside the anchor window; not an
independent source of truth. The anchor is a [[Journal event]] ID or a
timestamp in v0.2, extensible to review IDs (v0.3) and commits (with run
state) without changing the output contract.
Use when: reviewing what changed for a task or epic since a checkpoint — a QA
reviewer picking up a task window, a human resuming a phase.
Avoid: confusing with [[Activity log]] (the raw chronological stream this
aggregates) or treating it as a stored entity (it is a computed read).

**Search**:
A lightweight read command for literal, case-insensitive discovery across
structured projected fields such as IDs, titles, labels, owners, summaries,
and markdown paths. It is not a query language and does not scan full markdown
body text.
Use when: finding candidate entities before a focused `get`, `report`, or
read-only SQL command.
Avoid: using as a synonym for `fid sql`, saved views, or shell `rg`.

**Evidence check**:
A read-only verification that cited evidence references still resolve
(`fid evidence check`, v0.4). For each embedded evidence ref on an AC
verification, QA finding, or validation outcome in scope, it reports `resolved`,
`missing`, or `unverifiable` by a type-keyed resolver: commits (git ref resolves,
local) and file paths (file exists) are checkable; `command`, `manual`, and `url`
evidence are unverifiable. It verifies *resolution/existence* only — never
re-running a command or re-proving a claim (that would make fid an executor) and
never comparing content hashes (no capture-time hash is stored). A computed read
over existing embedded evidence; it adds no registry table.
Use when: a QA Lead or Reviewer asks whether cited tests, commits, or files still
exist before trusting a verification or outcome.
Avoid: treating it as re-execution (it does not run `go test`) or as content
drift detection (it checks existence, not whether the content still supports the
claim).

**Attention report**:
The pull-based, entity-aggregating digest of items currently needing human — or,
by declared audience, senior — judgment (`fid report attention`). Items are
grouped by reason: epics awaiting closure approval, open escalations without a
resolving decision (flagged stale past a window), issues blocked by a cancelled
(dead-end) blocker, and live human-held claims. Each group declares an
`audience`; further reasons (blocking findings v0.3, validation outcomes v0.4)
light up as their source entities ship, keyed by reason without a contract
change. A computed read, not a stored entity, and intentionally low-noise: an
empty result is the healthy state. Replaces the human grepping across current
files; it is not a watcher or notification system (those are out of scope).
Use when: the human (or, later, a senior role) wants the single low-noise list of
what needs judgment now before approving, interrupting, or resuming a phase.
Avoid: confusing with [[Change digest]] (what changed since a checkpoint) or with
notifications/watchers; treating it as a stored queue (it is recomputed each
call).

**Pattern report**:
A deterministic read that groups recurring quality signals — findings, concerns,
gate runs, triage events, and validation-outcome failures — by structured
dimensions (trigger class, label, area, gate command, finding severity, affected
file-path prefix, repeated `VAL-NNN` failures) and ranks the recurring groups as
candidates with counts and exemplars (`fid report recurring-findings` /
`recurring-concerns` / `recurring-patterns`, v0.4). It **surfaces** candidates for
a human or TL to turn into a [[Decision]] or principle; it never authors the
recommendation itself. Deterministic grouping only — no ML, embeddings, or
semantic theme extraction (fid is a local-first analyzer, not an AI service). A
computed read over existing signals; adds no schema.
Use when: a TL or QA Lead at wrap-up wants recurring failure themes turned into
candidate principles, refactors, or templates.
Avoid: expecting fid to write the recommendation prose (it ranks structured
patterns; a human authors the principle); confusing it with the [[Attention
report]] (open items needing judgment now) or with single-metric computed retro
signals (not grouped recurrences).

**Handoff pack**:
The generated bundle that equips one role to act on one entity (`fid handoff
--role <role> --issue <id>`): a rendered prompt plus a machine-readable
**manifest** of required reads, commands, stop conditions, expected writes, live
blockers, and (from v0.4) evidence obligations. The manifest is *derived* from
configurable sources — the role's context pack (the `--context-pack` data
bundle), the configured workflow, and role permissions — not hard-coded, so it
retracks when a project changes its workflow. The target role is a parameter,
never an enumerated subcommand. A read-only computed view: generating a handoff
claims nothing and writes nothing.
Use when: a [[Workflow Runner]] (or human) spawns a role against a task or epic
and needs the exact prompt plus obligations without hand-authoring them.
Avoid: confusing the handoff pack (prompt + orchestration manifest for spawning a
role) with the context pack (the role-tuned *data* bundle it composes over);
treating it as a write (it is a view and acquires no claim).

**Task pack**:
A structured authoring artifact — a reviewed phase plan listing an epic and its
tasks with acceptance criteria, dependencies, labels, estimates, and
sequence/priority — ingested in one create-only transaction (`fid task-pack
apply`, deferred/conditional). Apply produces a single [[Journal event]] with
many effects and returns a source-ref→fid-id mapping. Input is structured
(YAML/JSON/TOON), create-only, with no round-trip or sync — which is what keeps
it distinct from markdown import (`fid import`, the separately deferred v0.5
capability).
Use when: describing bulk creation of a planned phase from one reviewed
transaction, if the dogfood gate is met.
Avoid: confusing a task pack (a create-only authoring *input*) with a
[[Handoff pack]] (a generated prompt + manifest *output*), or with markdown
import/sync (round-trip editing of an existing markdown tree).

**Journal event**:
One durable JSON object appended as one line to `.fid/journal/YYYY-MM-DD.jsonl`.
It is the source-of-truth record for a logical write and contains event
metadata plus ordered projection effects. Mutable-state effects carry
preconditions so replay can detect logical conflicts.
Use when: describing persisted write history, replay, merge conflict
detection, or activity projection.
Avoid: using for projected rows in SQLite; those are derived read-model state.

**Effect precondition**:
An expected current value attached to a mutable journal effect. Replay checks
it before applying the effect and stops with a logical conflict if the
projected DB state differs.
Use when: describing stale-write and cross-branch conflict detection.
Avoid: using for append-only audit rows that do not depend on prior mutable
state.

**Compensating event**:
A new journal event that corrects the effect of an earlier valid-but-wrong
event while preserving both events in history (`fid event revert`). It applies
one of two modes by the target event's type: **inverse** compensation restores
the prior value of a mutable-state effect (e.g. a transition's prior status),
guarded by an effect precondition that the current state still matches what the
event produced; **retraction** marks an append-only audit fact (an AC
verification, concern, comment, or triage verdict) as retracted so it drops out
of current-state derivation while its row is preserved. Eligible only for
explicitly reversible event types — creations, completions, claims, telemetry,
and composites are not, and revert returns the forward domain-correction path
instead. Never a mutation or deletion of old journal lines.
Use when: describing first-class correction via `fid event revert`, or the
`revert.eligible`/`mode` verdict on [[Journal event]] inspection.
Avoid: "undo" when it implies deleting or rewriting audit history; assuming every
event is reversible (eligibility is event-type-keyed).

**Write lock**:
The checkout-scoped exclusive lock at `.fid/write.lock` that serializes fid
write commands before validation, journal append, and SQLite projection. It
protects the append-only journal as the source of truth; SQLite transactions
still protect projection atomicity inside the lock.
Use when: describing same-checkout write serialization.
Avoid: treating SQLite's default locking as the write serializer for journal
commits.

**Dry-run**:
A simulated write. `--dry-run` on any write command runs validation and builds
the would-be [[Journal event]], then stops before the append/commit point:
nothing is persisted. The result reports whether the write would apply, the
validation outcome, and a preview of the event effects; the exit code mirrors
the simulated outcome. IDs that are minted only at append render as
placeholders.
Use when: previewing or gating a write — especially a composite — before
committing it.
Avoid: assuming a dry-run reserves an ID or a claim; it changes no state.

**Agent run**:
The record of one agent execution against an entity (`agent_runs`,
`RUN-<hash>`): token/cost telemetry plus, from v0.2, optional git context
(branch, base/head commit, dirty, worktree, `git_source`) and a lifecycle
(`open` → `closed` | `abandoned`). Opened by `fid dev start` (or recorded
post-hoc by `fid run record`), closed by `fid dev finish`/`fid dev block`,
abandoned by `fid dev abandon`, and continued by `fid dev resume` (which opens
a new run linked via `resumed_from_run_id`).
Use when: referring to a single agent execution attempt and its telemetry, git
context, and lifecycle.
Avoid: treating it as telemetry only (its role broadened in v0.2), or as a
[[Claim holder|Claim]] (a run does not itself grant exclusive access).

## Roles and Who-Does-What

The role enum (`developer | qa_reviewer | qa_lead | tl_lite | tl |
workflow_runner | human`) shows up in four semantically distinct positions.
Each has its own term:

**Actor role**:
The role that performed (or is performing) an action. Stored on per-event
`_role` columns in append-only history (`verifier_role`, `reviewer_role`,
`developer_role`, `raised_by_role`, `author_role`, `triager_role`) and
projected as `actor` on [[Activity log]] entries. Passed via `--by <role>`
on every write CLI invocation. A fact about the past.
Use when: recording who did something.
Avoid: confusing with [[Owner role]] or [[Claim holder]].

Acceptance-criterion verification is multi-author: the same AC may have
independent outcomes recorded by Developer, QA Reviewer, QA Lead, TL/Tech
Lead, Human, or another configured verifier role. Workflow transitions and
validation contracts specify which verifier roles are required for a given
scope.

**Owner role / owner actor**:
Long-lived accountability for an [[Issue]] across its lifetime. Stored on
`issues.owner_role` and `issues.owner_actor`. `owner_actor` is optional —
NULL means "any actor of `owner_role`". Set at creation; changes are
recorded in the journal and projected into [[Activity log]]. Distinct from
[[Claim holder]].
Use when: describing who is accountable for an issue's lifecycle.
Avoid: "responsible role" (RACI overloaded), "assignee" (Linear-style).

**Claim holder**:
The role currently holding the exclusive-access lease on an entity. Stored
on `claims.claimed_by`. Short-lived (minutes, TTL-bounded). A right to
*advance* the work right now. Does not transfer ownership.
Use when: describing who has the current working lock on an entity.
Avoid: confusing with [[Owner role]] (long-lived accountability).

**Required validator**:
A planned verifier role/evidence obligation for a [[Validation assertion]].
Stored as one or more `validation_assertion_validators` rows, each pairing
the assertion with a required validator role and evidence method. A
statement about a plan, not a recorded fact. The corresponding
`val_assertion_outcomes.recorded_by_role` is the actor-role fact of who
actually recorded the outcome.
Use when: specifying who must provide evidence for a validation assertion.
Avoid: confusing with [[Actor role]] (what happened) vs required validator
(what should happen).

## Raised Quality Signals

Four adjacent concepts that all flag "something is not right" but at
different points in the workflow and from different roles. Keep them
distinct.

**Concern**:
A structured flag raised by the Developer during a task completion. One row
in the `concerns` table per concern; prose lives at `body_md_path`. Status
moves `open → addressed | escalated | wontfix`. An open concern triggers
[[Workflow Runner|workflow_runner]] to spawn [[TL-Lite]] for [[Triage]].
Use when: a developer flags something during implementation.
Avoid: "issue" (overloaded with [[Issue]]), "flag" (vague).

**Finding**:
A structured discovery made by a QA Reviewer or QA Lead during a
[[QA review]]. One row in `qa_findings`. Has `severity` (`blocking` |
`improvement`) and status. May escalate.
Use when: QA discovers something during review.
Avoid: using "finding" for developer-raised flags — those are [[Concern]]s.

**Escalation**:
An unresolved situation that requires a senior decision. One row in
`escalations`. Open until linked to a resolving `decision`. Created from a
[[Finding]], a [[Concern]] (via TL-Lite BLOCK verdict), or directly by a
senior role.
Use when: an issue exceeds local authority and needs TL or human input.
Avoid: using for routine quality signals that don't need senior judgement.

**Carry-forward**:
Open work scheduled into a future epic, or promoted to a project principle.
One row in `carry_forward`. Authored by TL at wrap-up. Multi-epic lifespan,
unlike the per-attempt scope of a [[Concern]].
Use when: deferring work to a future epic at wrap-up.
Avoid: confusing with [[Concern]] (per-attempt, developer-raised).

## Adjudication and Intake

**Triage**:
A recorded TL-Lite adjudication of a mid-task concern. One row per verdict
in `triage_events`. The verdict is `RESUME` or `BLOCK`; the source is
typically a developer-raised [[Concern]]. This is *not* a state — it's an
event. The word "triage" in fid always means this.
Use when: describing TL-Lite's adjudication of a soft blocker.
Avoid: using "triage" to mean the issue intake/inbox state — that's
[[Inbox]].

**Inbox**:
The unsorted-incoming state of an issue. `issues.in_inbox = 1` means the
issue has arrived but hasn't been assigned to an epic or prioritised yet.
Linear-style intake.
Use when: describing unsorted incoming issues.
Avoid: "triage state", "triage queue" — overloads with [[Triage]] which is
a different concept here.

## Roles in the Enum

**TL (Tech Lead)**:
The strategic senior role. Invoked at epic kickoff and at epic wrap-up.
Rarely invoked mid-flight when an escalation requires Opus judgement. Never
triages [[Concern]]s directly — those go through [[TL-Lite]].
Use when: describing the strategic planning and wrap-up role.
Avoid: confusing with [[TL-Lite]] (tactical mid-flight role).

**TL-Lite**:
The mid-flight tactical role; single-use per invocation. Two modes: triage
mode (invoked by [[Workflow Runner]] on a soft blocker, ends with literal
`VERDICT: RESUME` or `VERDICT: BLOCK`) and question mode (invoked by the
human for one tactical question).
Use when: describing mid-flight tactical decisions or triage.
Avoid: confusing with [[TL]] (strategic bookend role).

**QA Lead**:
The quality planning role. Authors the [[QA plan]] at epic kickoff and runs
the final epic [[QA review]]. Generates the human verification checklist.
Use when: describing the quality planning and final review role.
Avoid: confusing with [[QA Reviewer]] (mid-flight review role).

**QA Reviewer**:
The mid-flight quality companion to [[QA Lead]]. Single-use; invoked by
[[Workflow Runner]] every three tasks for mid-epic code review. Writes
findings; raises blocking findings as [[Escalation]]s.
Use when: describing mid-flight code review.
Avoid: confusing with [[QA Lead]] (bookend quality role).

**Developer**:
The implementing role. Single-use per task. Reads task spec, writes code via
TDD, runs gates, writes a structured completion report including the
mandatory `## Concerns` section.
Use when: describing the implementation role.
Avoid: "engineer", "coder".

**Workflow Runner**:
The mechanical loop driver — fid's `workflow_runner` role (one of the seven
seeded roles). Drives an epic's task loop, spawns agents, dispatches on
completion-report inspection and review findings. Intentionally a cheaper
model — its decisions are state-machine deterministic, not strategic. May
later become a deterministic harness.
Use when: describing fid's mechanical orchestration role (`workflow_runner`).
Avoid: "Orchestrator" (misleading; the capable model is the TL).
Distinct from: "Phase Runner" — the Claude Code build agent (`@phase-runner`;
`09_AI_WORKFLOW.md` §2.4) that corresponds to this role while building fid.
Current, not a legacy alias. The build pipeline will converge on "Workflow
Runner" once fid is dogfooded.

**Human**:
The product/approval role and the escape valve. Approves epic closure, runs
the human verification checklist, invokes [[TL-Lite]] in Question mode,
owns product-level decisions. Always allowed to override role boundaries via
`--by human`, while role-attributed evidence still records `human` rather than
impersonating another role.
Use when: describing the human product owner and approval role.
Avoid: "user" in workflow context (ambiguous with end-user).

## Workflow Terms

### Operating Mode

Definition: The per-phase declaration - light, normal, or heavy - of which
workflow artifacts the phase must produce. Modes vary the artifact set, never
the evidence standard: exact-command gates, evidence rules, commit discipline,
and hard stops apply in every mode.
Use when: Planning a phase, choosing orchestration, or deciding whether a
phase artifact is required.
Avoid: Process level, rigor tier; TL-Lite warm-up mode, which is a
trust-calibration state, not an operating mode.
Related: Phase Brief, Validation Contract, QA Plan.
Examples: "This config-only phase runs in light mode." "The migration phase
is heavy mode, so command feasibility is revalidated, not carried forward."

### Validation Contract

Definition: A per-phase artifact that lists the `VAL-NNN` assertions that
must be proven before the phase is complete.
Use when: planning or reviewing phase-level proof obligations.
Avoid: acceptance checklist, QA plan, test plan.

### Validation Assertion

Definition: One testable behavior, invariant, or release-readiness claim
inside a validation contract, identified as `VAL-NNN`.
Use when: connecting a source requirement to task work and QA evidence.
Avoid: acceptance criterion when referring to the phase-level assertion.

### Triage Entry

Definition: A per-phase `TRG-NNN` record of TL-Lite's disposition for a
soft blocker or escalation trigger.
Use when: the Workflow Runner routes a soft blocker through TL-Lite.
Avoid: decision entry when no durable choice was made.

### Workflow-Lint

Definition: A stack-neutral repository check that verifies agent workflow
artifacts, templates, config, and generated skill copies have not drifted.
Use when: checking prompt/template/config/skill consistency.
Avoid: quality gate when referring to product code checks.

## Flagged Ambiguities

| Term or phrase | Why it is ambiguous | Owner | Resolution status |
|---|---|---|---|
| `task_completions` kind guard ("task only" vs "task or subtask") | TRD SCHEMA §14.1 specifies the backstop kind guard is the "same" as `acceptance_criteria` (task or subtask); TRD §4.1's entity note says "(task only)". §4.1 explicitly states its column lists are indicative and §14.1 is the trigger reference. | Tech Lead | Phase 2 resolved the backstop trigger to §14.1 (task **or** subtask). Reconcile the §4.1 wording (or confirm the divergence is intentional for app-layer completions) at Phase 2 spec sync. |

*Last updated: 2026-07-11*
