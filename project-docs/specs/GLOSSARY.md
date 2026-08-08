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
A local-first, agent-first issue tracker and durable operational work-state
and evidence backend for AI-assisted software development. SQLite-backed CLI.
A tapered hand tool used to open up tightly-laid rope so the strands can be
worked separately. Pairs with **whet** (sharpening) as a workshop-drawer set —
fid is the tracker; whet is a separate sibling project.
Use when: naming the product that owns issue tracking and operational work
state for agents and external workflow kits.
Avoid: describing fid as the owner of engineering-workflow policy, provider
prompts, the model loop, canonical project prose/code, or Human judgment.

**Operational work state**:
The durable, queryable record of what work exists and what has happened while
executing it: issues and relationships, lifecycle state, claims and runs,
ordered events and projections, receipts and evidence, attention, reviews,
outcomes, and telemetry. fid is authoritative for this domain; its journal is
the canonical persistence representation.
Use when: distinguishing resumable workflow facts from project requirements,
source code, kit policy, or Human judgment.
Avoid: unqualified "source of truth"; claiming fid owns PRD/TRD prose, Git
objects, engineering-workflow policy, or product approval.

## Core Entities

**Issue**:
Any row in the unified `issues` table — an epic, task, or subtask. Issues and
their parent/dependency relationships are fid's only durable work graph; a Work
Envelope selects Issues rather than duplicating them as persistent Work Items
(ADR-0005). The umbrella term is not a synonym for "ticket".
Use when: referring to any durable unit of tracked work regardless of kind.
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
Use when: describing a single implementable durable unit of work.
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

**Projected workflow tables**:
The four SQLite tables that hold workflow definitions as data —
`workflows`, `workflow_states`, `workflow_transitions`, and
`workflow_transition_requires` (`03_TRD_SCHEMA.md` §3). Seeded verbatim from
the v0.1 defaults by migration `0007` and, as of Phase 18, the runtime source
of truth for states, status categories, allowed edges, and per-edge rules.
Like every SQLite table they are a derived projection; the defaults are
seeded by migration DDL because they are system reference data.
Use when: naming where workflow structure is read from at runtime.
Avoid: calling them the source of truth for *history* (the journal is);
implying custom workflows can be authored in v0.2a slice 1 (that is Phase 19).

**Workflow registry**:
The read-only in-memory structure `internal/workflow` loads from the
[[Projected workflow tables]] on each read (`workflow.Load`). It resolves a
state's category, whether an edge is allowed, and which named rules attach to
an edge. It evaluates nothing by itself — rule functions stay named Go
functions keyed by rule name.
Use when: describing how a consumer resolves workflow structure at runtime.
Avoid: "workflow definitions" when you mean the loaded registry; treating it
as a cache (it is loaded per read, not memoised).

**Frozen workflow definitions**:
The v0.1 hard-coded Go `definitions` map, retained after Phase 18 **solely as
test parity fixtures**. No production code path reads it; its only purpose is
to prove the [[Workflow registry]] stayed behavior-identical to v0.1.
Use when: describing the parity reference in tests.
Avoid: describing it as the source of truth, current behavior, or a fallback —
it is none of those; reading it from production code is a defect.

**Custom workflow**:
A user-authored workflow created by `fid workflow define`, stored as
journal-derived rows in the projected workflow tables rather than as a seeded
default.
Use when: distinguishing an authored workflow from a system reference default.
Avoid: calling it a second source of truth; its history remains journal-canonical.

**Workflow assignment**:
A `workflow_assignments` row binding a workflow to an issue kind or to the
project. A kind-scope assignment takes precedence over a project-scope
assignment when an issue is created.
Use when: describing how a workflow is selected for a new issue.
Avoid: treating an assignment as part of the workflow definition itself.

**Initial state**:
The state an issue receives at creation: the first declared `unstarted` state
in its resolved workflow, or the first declared state when none is
`unstarted`.
Use when: describing creation-time workflow resolution.
Avoid: assuming every workflow starts in `todo`.

## Execution and Audit

**Gate**:
A deterministic command that produces a pass/fail signal (test runs, linters,
coverage thresholds, format checks). From v0.2 (DEC-059) the trusted path is
a registered [[Gate definition]] executed by fid (or a trusted runner),
producing an [[Execution attempt]] and [[Execution receipt]]. The v0.1
`fid dev gate` path records a caller-reported `gate_runs` row; those rows
remain queryable but carry `trust_source=reported` and cannot satisfy a
trusted-gate obligation. A failing gate is a hard stop — TL-Lite cannot
resolve it; the Developer fixes and re-runs.
Use when: describing deterministic quality checks during implementation.
Avoid: "quality gate" when referring to subjective review; treating a
reported exit code as proof a command ran.

**Gate definition**:
A versioned value in project-owned machine-readable configuration, governed by
the kit-owned schema and activation procedure, that names one executable gate:
stable gate ID, version, structured argv array, working directory, timeout,
environment allowlist, and execution policies. fid never executes this mutable
authoring value directly; activation creates an immutable
[[Gate definition registration]]. Obligations reference gates by registered
identity/version, never by display string.
Use when: registering or referencing an executable gate.
Avoid: free-form shell command strings; "gate" for the definition when you
mean a specific run (that is an [[Execution attempt]]).

**Gate definition registration**:
The immutable operational authorization fid creates when the kit validates and
activates a complete project-approved [[Gate definition]]. It binds the exact
definition and digest to project, overlay, repository, authorization, status,
and registration version. The trusted executor resolves it by reference; a
mutable authoring-file edit cannot change it or authorize execution.
Use when: naming what fid may execute and what an [[Execution receipt]] binds.
Avoid: calling the authoring file the live registry; mutating a registration
instead of activating a new version.

**Execution attempt**:
An immutable record that fid (or a trusted runner) started executing a
[[Gate definition]] — appended as an execution-start fact *before* the
process spawns. An attempt with no terminal [[Execution receipt]] is an
abandoned start, surfaced by `fid maintenance check`; it is never a pass.
Use when: describing the start-side fact of the two-event execution model.
Avoid: treating a start fact as an outcome, or backfilling a receipt for a
crashed attempt.

**Execution receipt**:
The immutable terminal record of one [[Execution attempt]], appended only
after the process terminates. Captures producer, [[Trust source]], gate
definition ID/version, argv/cwd/timeout, exit code or signal, duration,
git branch/base/head/tree identity and dirty state, toolchain fingerprint,
bounded sanitized output excerpts with content hashes and
truncation/redaction flags, and terminal status
(`pass | fail | timed_out | cancelled | harness_error`). A crash may lose a
pass but must never fabricate one.
Use when: citing trusted deterministic evidence for an obligation.
Avoid: "gate run" (the legacy reported row); mutating or re-snapshotting
receipt content into other rows (link to the receipt instead).

**Trust source**:
The declared provenance class of a piece of deterministic evidence:
`executed` (fid ran the registered gate itself), `attested` (a separately
trusted runner produced the receipt under the versioned attestation
contract), or `reported` (a caller asserted the outcome; includes all
legacy v0.1 `gate_runs`). Only `executed`/`attested` evidence can satisfy a
trusted-gate obligation; `reported` stays queryable history. Consumers must
never silently downgrade from trusted to reported.
Use when: classifying how gate/receipt evidence came to exist.
Avoid: inferring trust from age or formatting; treating an import as
trusted because it parses.

**Unverified gate result**:
The status of a required-command result that was **not captured in the session
that records it** — attested from memory, transcribed from another agent's
report, or promised against a named future control. Recorded as `UNVERIFIED`,
which is not a passing state and never substitutes for `PASS` (DEC-102).
Use when: a report cannot cite a verbatim terminal token it produced itself.
Avoid: confusing with the `attested` [[Trust source]]. `attested` describes a
*receipt produced by a separately trusted runner under a versioned contract*
and can satisfy an obligation; `UNVERIFIED` describes a *human or agent claim
with no captured output* and cannot. A human running a gate by hand and
reporting the outcome produces an `UNVERIFIED` result, not an `attested` one.

**Evidence obligation**:
A declared requirement that a workflow transition, AC verification, QA
review, or validation outcome be backed by specific evidence — e.g. a
current [[Execution receipt]] for a named [[Gate definition]], or a Human
outcome for a judgment item. Deterministic obligations require trusted,
terminal receipts or another checkable evidence type; an unverifiable
command string cannot satisfy them. Human-judgment obligations are
satisfied by recorded manual evidence and are never mislabeled as
deterministic proof.
Use when: specifying what evidence a state change requires.
Avoid: conflating with [[Required validator]] (who must provide it).

**Evidence freshness**:
The validity dimension of evidence, separate from existence:
`current`, `stale`, `missing`, `unverifiable`, or `inapplicable`. v0.2
starts conservative: any relevant target-tree change (new head/tree
identity) makes a gate receipt `stale`. Path-aware invalidation requires
dogfood evidence that the coarse rule is too expensive, a sound
affected-path model, and a new recorded decision.
Use when: deciding whether existing evidence still proves its claim.
Avoid: treating resolution/existence (`fid evidence check` v0.1 semantics)
as freshness.

**Work envelope**:
A versioned, structured capability describing one approved slice of work:
outcome, one primary fid Issue plus an explicitly bounded related Issue set,
risk lane, approved authority, protected decisions, validation and gate
obligations, budgets, and stop conditions.
Promoted from the conditional [[Task pack]] concept (DEC-063). Kit/project
policy supplies its engineering semantics; fid persists its operational state
and exposes orchestration primitives over it — next actionable work with typed
skip reasons, claim/lease and resumable run state, retry/intervention budgets
and circuit-breaker state, machine-readable rejection/stop reasons, and
attention items — while a deterministic runner (a separate kit component)
drives the loop. The kit evaluates policy and compiles the approved envelope;
fid enforces generic integrity and its recorded preconditions without
recalculating kit policy. The Envelope is not a second work graph and creates
no persistent Work Item identity (ADR-0005).
Use when: describing the structured contract a runner executes against.
Avoid: treating it as an LLM prompt, fid as the loop host, the envelope as a
replacement Issue, or its selected-Issue references as independently mutable
work entities.

**Next action**:
A noncanonical projection of the action currently permitted or required,
derived from Issue, Work Envelope, Obligation, Claim, Run, Review, evidence,
stop, and attention state. It names the action kind, subject, required actor
role, and typed reason, but has no independent lifecycle or durable work
identity.
Use when: asking the harness what can or must happen next or why no action is
available.
Avoid: persisting it as a Work Item, using it as a second backlog, or treating
the projection as authority independent of the records from which it derives.

**Engineering-workflow policy**:
The rules that classify risk, grant authority, select roles/reviewer profiles,
route work, set budgets and obligations, and determine when Human judgment is
required. The external workflow kit and project configuration own this policy;
fid records approved operational inputs and outcomes without becoming the
method's author (ADR-0004).
Use when: distinguishing kit-owned strategic workflow rules from fid-owned
operational state, integrity checks, and evidence records.
Avoid: using it for a fid [[Workflow]], which is an issue lifecycle state
machine, or asking fid to independently recalculate kit risk or routing policy.

**Human validation brief**:
The generated wrap-up report for Human judgment (v0.4): outcomes delivered
and user-visible changes; product/architecture/data/dependency/security
decisions; deviations, unresolved risks, failed or carried-forward
obligations, and open escalations; receipt and validation roll-up with
spend, retries, and interventions; and the human-only checks still
requiring judgment. A computed view over structured state — fid renders it;
it does not author judgment.
Use when: presenting phase/epic closure to the Human.
Avoid: confusing with the [[Attention report]] (live pull-based digest of
open judgment items) or hand-authored wrap-up prose.

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
A read-only verification of cited evidence (`fid evidence check`). For each
embedded evidence ref on an AC verification, QA finding, or validation outcome
in scope, it reports `resolved`, `missing`, or `unverifiable` by a type-keyed
resolver: commits (git ref resolves, local) and file paths (file exists) are
checkable; free-form `command`, `manual`, and `url` evidence are unverifiable.
From the v0.2 evidence foundation (DEC-060/061) it additionally reports
[[Evidence freshness]] for receipt-backed evidence (`current`/`stale` from the
receipt's captured tree identity and content hashes). It never re-runs a
command inside the check — re-execution is `fid gate run`'s job — and it adds
no registry table of its own (evidence records are first-class rows from
v0.2b).
Use when: a QA Lead or Reviewer asks whether cited evidence still exists and
is still fresh before trusting a verification or outcome.
Avoid: treating it as re-execution (it does not run `go test`); calling
reported gate strings checkable (they stay `unverifiable`).

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
*(Superseded 2026-07-18 by [[Work envelope]], DEC-063.)* The former
deferred/conditional concept of a structured authoring artifact — a reviewed
phase plan ingested in one create-only transaction (`fid task-pack apply`).
Its semantics (single [[Journal event]] with many effects, source-ref→fid-id
mapping, structured YAML/JSON/TOON input, create-only, no round-trip/sync)
carry into `fid envelope apply` (PRD-OPS-021).
Use when: reading historical artifacts that predate the work envelope.
Avoid: using in new specs — write [[Work envelope]]; confusing either with a
[[Handoff pack]] (a generated prompt + manifest *output*) or with markdown
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
*advance* the work right now. Does not transfer ownership. In the V2-aligned
path, one Issue may have at most one active mutating Builder Claim; the Claim
also binds the Envelope revision and worktree. Parallel mutation requires
distinct selected Issues, separate worktrees, and provably disjoint scopes,
or the second action is `mutation-conflict` (ADR-0006).
Use when: describing who has the current working lock on an entity.
Avoid: confusing with [[Owner role]] (long-lived accountability), claiming a
Work Envelope itself, or treating separate agents as sufficient isolation.

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
completion and review state. Intentionally a cheaper model — its decisions
are state-machine deterministic, not strategic. Per DEC-063 the direction is
a deterministic runner (a separate kit component) driving fid's [[Work
envelope]] orchestration primitives, with lifecycle bookkeeping, evidence
completeness, and stop-policy enforcement in fid rather than in the runner
prompt.
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

*Last updated: 2026-08-08 (ADR-0004/0005: operational boundary and Issue-native
work identity). Previous material update: 2026-07-25 Phase 18 wrap (DEC-073:
projected workflow tables, workflow registry, frozen workflow definitions).*
