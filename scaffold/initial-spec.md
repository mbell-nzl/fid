# `fid` — Specification v0.2

> **Status:** Draft v0.2 — supersedes v0.1 (TRACKER_SPEC_v0.1.md, 2026-05-09).
> Captures naming decision, entity reshape (epic-as-issue), templates,
> configurable workflows, sync/export, lightweight feature borrows from
> Linear/Jira, and TOON/JSON output strategy.
>
> **Tool name:** `fid`. A fid is a tapered hand tool used to open up
> tightly-laid rope so the strands can be worked separately. The metaphor
> fits: this tool opens up tightly-coupled work into separable parts you
> can manipulate. Pairs with `whet` (sharpening) as a workshop-drawer set.
>
> **Audience:** the author, future TL agents bootstrapping projects with
> this tool, and (eventually) team members.
>
> **Purpose of this document:** lock down the design intent and schema
> shape before any implementation work begins, so subsequent task
> breakdown has a stable target.

---

## 0. Changes from v0.1

This is a substantial revision. Major changes:

- **Naming and IDs.** Tool name fixed as `fid`. Issue IDs use a
  per-project prefix plus a hash suffix (`WHET-a1b2c3`, `FID-7c3e9f`).
  Prefix is configurable at `fid init`; hash length scales
  progressively (6 → 8 hex chars) as issue count grows.
  Hash-based IDs reduce coordination needs when multiple agents or
  branches create issues concurrently; rare collisions are detected
  locally or during journal replay. Other generated entities follow
  fixed prefixes documented in §6.0.
- **Schema reshape.** `epics` and `tasks` collapse into a unified `issues`
  table with `kind` discriminator (`epic` | `task` | `subtask`).
  Kind-specific structured data lives in side tables that FK to `issues`
  (task_completions, acceptance_criteria, epic_metadata).
- **Sub-tasks.** Same `issues` table, `kind = 'subtask'`. Simpler status
  set inherited via workflow assignment.
- **Workflows.** v0.1 ships three hard-coded default workflows
  (epic / task / subtask) implemented in code. v0.2 introduces the
  configurable workflow engine (new tables `workflows`,
  `workflow_states`, `workflow_transitions`,
  `workflow_transition_requires`) which migrates the defaults into rows
  and enables user customisation. `status_category` mapping is present
  from v0.1 (as code) to enable portable queries.
- **Templates.** New `templates` table with monotonic versioning.
- **Sync / export.** v0.4 ships markdown export (`fid export`) with
  the frontmatter contract; markdown is downstream of the DB (which
  is itself downstream of the journal — see ADR-0001). Bidirectional
  sync (`fid import`, `fid sync`, MD-canonical mode) is deferred to
  v0.5+ if a real use case emerges.
- **Linking.** Formalised relation enum; CLI auto-creates the inverse.
- **Lightweight Linear/Jira borrows.** Labels, estimates, inbox state,
  saved views, activity log, comments, optional cycles, due dates.
- **Output formats.** JSON is the default for all read commands; TOON
  is opt-in via `--format toon` and recommended for agent invocations
  of `list` commands where token efficiency compounds. Markdown for
  human inspection. Always overridable via `--format`.
- **Kanban / dashboard.** TUI board ships in v0.3; web dashboard is a
  separate package shipped in v0.5.

---

## 1. Overview

`fid` is a project-scoped CLI tool that provides a structured, queryable
knowledge graph over the artifacts an AI-assisted development workflow
produces. It is designed to work for both **bootstrapping (sequential
epics walked in order)** and **steady-state (prioritised epics, pulled
by readiness)** styles without changing schema. See §4.

The default profile encodes an opinionated multi-role agentic workflow
with three-tier quality (validation contracts → QA plans → QA reviews).
Both the role/permissions surface (ADR-0002) and the QA depth (ADR-0003)
are layered and optional: use the parts that fit your project, ignore
the rest.

It is **not** a replacement for the agent workflow protocol (the role
definitions, escalation triggers, deterministic gates, and epic wrap-up
discipline). It is the storage and query layer underneath that protocol.

The key design move is the **row / document split**: rows hold structured
metadata and identity; markdown files hold the prose. Both are
first-class inputs and outputs.

---

## 2. Goals and Non-Goals

### 2.1 Goals

1. **Queryable history across epics.** Cross-cutting questions
   ("every DEC that touched the data model", "all ESC entries with
   trigger=schema-change", "tasks where actual_tokens > 1.5× estimate")
   become single SQL queries rather than `grep` plus agent reading.
2. **Schema-enforced workflow invariants.** Required fields on
   completion reports, required `VERDICT:` line on triage events,
   required AC verification before final epic pass — all enforced at
   write time via configurable workflow transition rules.
3. **First-class non-task entities.** Decisions, escalations, triage
   events, QA reviews, retrospectives, requirements, carry-forward
   items are peer entities to issues, not "issues with a label".
4. **Bridge to existing markdown discipline.** Prose lives in markdown
   files referenced by row-level pointers. Rows index metadata, not text.
5. **Token / cost telemetry per task and per agent run** to feed the
   §9 retrospective discipline with computed signals rather than
   judgement calls.
6. **Versatility across bootstrap and v1+ modes** without schema changes
   — only usage patterns differ.
7. **Multi-format output: JSON by default for universal pipeability;
   TOON opt-in for token-efficient agent reads of tabular data;
   markdown for humans.** Same CLI command, format flag selects.
8. **Project-scoped database; common CLI across projects.** One
   `.fid/` directory per project, single binary on the user's path.
9. **Workflow flexibility.** v0.1 ships hard-coded default state
   machines per `kind` with category mapping for portable queries.
   v0.2 makes these configurable per kind or per project via the
   workflow engine.
10. **Markdown export with a frontmatter contract.** Project state can
    be exported as a markdown tree matching whet's structure (or any
    configured layout). Exports include YAML frontmatter so structured
    fields are self-describing and machine-readable. Bidirectional
    sync (importing edited markdown back into the DB) is deferred to
    v0.5+ pending evidence of a real use case.

### 2.2 Non-Goals

- Not a replacement for git, PRDs/TRDs, or the agent workflow protocol.
- Not a multi-tenant server. Local-first, single project per database.
- Not a real-time collaboration tool. Last-write-wins assumed; no CRDTs.
- Not a hosting product. No auth, no remote sync (initially).
- Not a replacement for prose. Decision records, retrospectives, and
  completion reports remain prose-shaped artifacts; the tool indexes
  them.
- Not a custom-fields platform. The schema is opinionated; user-defined
  fields are out of scope for v1.

---

## 3. Core Principles

1. **Rows hold structure; markdown holds prose.** The schema captures
   identity, status, relationships, and metadata. Prose stays in
   markdown files referenced by `*_md_path` columns. Don't shred prose
   into columns.
2. **Append-only for completion reports and reviews.** A task can be
   re-attempted; a feature can be re-reviewed. Each attempt is a child
   row, not an edit. The history is the value.
3. **Pointers, not copies, for requirements.** PRD/TRD rows store
   `(id, summary, source_doc, source_section, last_synced_hash)`. The
   summary is a label, not a copy. Drift is detected, not papered over.
4. **Scoped typed edges in a single `links` table.** Traceability
   relationships between supported artifact types share one shape.
   Lifecycle relationships with dedicated schema stay out of `links`.
5. **State transitions are constrained.** Each kind has a defined
   state machine; the CLI rejects invalid jumps. Per-transition
   validation rules (e.g., `all_acs_verified` before `dev_complete`)
   are enforced at write time. In v0.1 the state machines and rules
   are implemented in code; v0.2 moves them into configurable tables.
6. **Field-level reads happen via on-demand markdown parsing.**
   `fid issue get task-3 --section acceptance` parses the spec file's
   `## Acceptance Criteria` heading on demand. Don't pre-shred sections
   into rows.
7. **Scheduling-agnostic schema.** The same tables serve sequential
   epics (sequence-order walk) and prioritised epics (priority-ranked
   pull). Scheduling is a field on `epic_metadata`, not a schema
   branch. See §4.
8. **The tool is the index, not the canon, for prose artifacts.**
   Decision records, retrospectives, and PRDs remain authoritative in
   their files. The DB makes them findable.
9. **Unified core, kind-specific edges.** Issues (epic / task / subtask)
   share core fields in one table; kind-specific structured data lives
   in side tables that FK back. Avoids both the parallel-table
   duplication problem and the nullable-field swamp.
10. **Categories beat labels for portability.** Workflow states have a
    `category` (`unstarted` / `started` / `completed` / `cancelled`).
    Dashboards, queries, and integrations target categories so they
    keep working when users rename states.
11. **IDs are hashes, not sequences.** Generated entity IDs use a short
    hex hash suffix rather than a sequential integer; §6.0 defines the
    prefix vocabulary and natural-key exceptions. Two agents on two
    branches can both create issues without coordination and with low
    collision risk; local collisions are retried and merge-time
    collisions surface during journal replay. The cost is lost
    ordering-by-ID; use `created_at` for chronology.
12. **Journal is canonical; DB is a derived cache; markdown export is
    a derived view.** `.fid/journal/YYYY-MM-DD.jsonl` is the source of
    truth, committed to git. `.fid/fid.sqlite` is gitignored and
    materialized from the journal — the journal append is the commit
    point, and DB projection is synchronous but replayable if
    interrupted. There is no background daemon. The DB can be deleted
    and rebuilt non-destructively. See ADR-0001 and §3.1.
13. **Opinionated default profile; configurable surface.** Roles and
    permissions ship as a built-in `default` profile, seeded into the
    `roles` table at `fid init`. From v0.2 the permissions matrix
    migrates to data and is per-project configurable. The tool is
    useful out of the box and customisable when needed; it does not
    dictate one workflow. See ADR-0002.

### 3.1 Durability and concurrency model

Per ADR-0001, the journal is canonical and the DB is a derived cache.
This drives what gets committed to git and how branches reconcile.

**What is committed:**

- `.fid/journal/YYYY-MM-DD.jsonl` — one file per day, append-only,
  one self-contained JSON event per line. The full source of truth.
  Every event captures the resolved row(s) being written, including
  generated values (hash IDs, timestamps, snapshot pricing, gate
  output) so replay is deterministic.
- The markdown export (when v0.4 ships) is optionally committed for
  human review and stakeholder consumption; per §11.1 it carries a
  "generated by fid" banner. Many projects will gitignore it.

**What is NOT committed:**

- `.fid/fid.sqlite` — gitignored. It is a query cache, materialized
  from the journal. `rm -rf .fid/fid.sqlite` is non-destructive; the
  next CLI invocation rebuilds it.

**Write semantics:**

- Every CLI write validates against the current DB, builds a fully
  resolved event, appends and flushes one JSONL line to the journal,
  then applies that event to SQLite in a DB transaction. The journal
  append is the commit point; DB projection is synchronous but
  replayable if interrupted. No background daemon, no sync window.
- If journal append fails, the DB is not modified and the write is not
  committed. If DB projection fails after the journal append, the next
  CLI invocation catches the DB up from the journal before answering.
- A malformed trailing JSONL line is treated as journal corruption, not
  a partial commit. Startup / `fid rebuild` stops and reports the
  corrupt tail so a human can truncate or repair it.
- On every CLI invocation, the DB is checked against the journal's
  current tail. If the DB is behind (or missing), it is rebuilt or
  caught up synchronously before the command answers. Rebuild is a
  linear scan of the journal; sub-second at single-project scale.
- `fid rebuild` is a v0.1 primitive; it is what `fid` itself runs on
  startup when it detects a stale cache.

**Journal event contract:**

- **One logical CLI write creates one journal event — no exceptions.**
  A primitive command like `fid issue create` is one event. A composite
  command like `fid dev finish` is also one event: its `effects` array
  contains inserts/updates across multiple tables (e.g., `agent_runs`,
  `task_completions`, `concerns`, `issues` status update,
  `activity_log` entries), applied atomically in a single SQLite
  transaction. Half-applied composite state is data corruption the
  replay engine would faithfully reproduce, so all-or-nothing is
  mandatory. If any step of a composite fails (e.g., concern parsing
  rejects the completion report), the entire event is rejected and
  neither the journal nor the DB is modified.
- Events use a semantic envelope plus generic ordered row effects. The
  envelope has `event_type` such as `issue.created` or `dev.finished`,
  an `event_schema` version, actor metadata, and command metadata for
  audit.
- `effects` is an ordered list of row operations: `insert`, `update`,
  or `delete`, with table name, primary key, and resolved values /
  patch. Effects include generated IDs, timestamps, snapshot pricing,
  gate output, token counts, and any activity-log rows needed for the
  projected DB.
- Replay upgrades the event to the current `event_schema`, then applies
  `effects` in order inside a SQLite transaction. Replay does not
  re-run command handlers or ask what the command would do under the
  current code.
- The journal is the canonical audit stream. `activity_log` is a
  projected, curated timeline table for queries and UI; it is populated
  only by applying journal event effects, never by a standalone write
  API.

**Event schema versioning:**

- Each event carries `event_schema: <n>` (a positive integer). The
  version increments whenever a schema migration changes the shape of
  any table that appears in event effects — new columns, renamed
  columns, removed columns, changed semantics.
- **Journal files are immutable once written.** Old events stay at the
  schema version they were created with. The replay engine upgrades
  them in memory; it never rewrites journal files.
- The replay engine maintains a chain of registered **upgrader
  functions**: `upgrade_1_to_2(event) → event`,
  `upgrade_2_to_3(event) → event`, etc. When replaying an event whose
  `event_schema` is older than the current version, the engine runs it
  through the chain before applying effects. Upgraders transform the
  `effects` array — adding columns with default values, renaming
  fields, adjusting values — without re-running command logic or
  reading DB state.
- Each upgrader is a **pure function**: old event in, new event out.
  No DB reads, no side effects. Testable with golden-file fixtures
  per version boundary: capture a real event at version N, run it
  through `upgrade_N_to_N+1`, assert the output matches the expected
  version N+1 shape.
- The upgrader chain grows linearly with schema versions. At fid's
  scale (single project, hundreds to low thousands of events) the
  cost is negligible. If it ever matters, a compaction step that
  rewrites a journal archive to the current schema is a
  forward-compatible optimisation — but not required for v1.

**Concurrency:**

- Single-worktree concurrent writes are safe via SQLite's default
  locking; agents in the same checkout serialise correctly.
- Multi-worktree concurrent writes against the same branch produce
  JSONL line-merge conflicts in today's journal file, resolvable as
  any other text merge. After resolution, the next `fid` invocation
  rebuilds the DB from the merged journal. The binary-blob merge
  problem does not arise.
- Logically-conflicting writes (e.g., both worktrees transition the
  same issue to incompatible states) surface as JSONL conflicts on the
  same line(s), where a human resolves. There is no automatic
  reconciliation of intent.

**Explicit non-goals for v1:** real-time concurrent writes with
automatic conflict resolution, CRDT-style merge, distributed claim
coordination across worktrees.

---

## 4. Epic Scheduling

Epics differ in *how their child tasks are picked up*, not in whether
they're planned. The `epic_metadata.scheduling` field captures this
on each epic.

### 4.1 Sequential epics

- Tasks are executed in explicit `sequence` order.
- Task `sequence` is a positive integer scoped to the parent epic.
  It is scheduling metadata, not a dependency edge.
- `fid ready --epic <id>` examines the lowest-sequence unfinished task.
  If that task is unblocked, it is returned. If that task is blocked,
  no later task is skipped to; the command returns no ready task and
  reports the blocker.
- The bootstrap pattern; what whet calls a "phase".
- Best when strict progression matters because each task builds on the
  previous, context needs to accumulate in order, or a wrap-up gate
  expects serial completion.

### 4.2 Prioritised epics

- Tasks are ranked by priority and pulled when ready.
- Task `sequence` is null; hard dependencies still use
  `blocks` / `blocked_by` links.
- `fid ready --epic <id>` returns the highest-priority unblocked task
  within the epic.
- Suits feature work, tech-debt clusters, ongoing initiatives —
  anywhere the scope is bounded but tasks aren't pre-ordered.

### 4.3 Declaration and unattached tasks

Scheduling is set at epic creation:

```
fid issue create --kind epic --scheduling sequential --title "Web UI bootstrap"
fid issue create --kind task --epic <id> --sequence 10 --title "Create shell"
fid issue create --kind epic --scheduling prioritised --title "Auth feature"
fid issue create --kind task --title "Fix off-by-one in pagination"   # epic optional
```

Tasks without a parent epic are scheduled by priority globally.
`fid ready` (no `--epic`) returns the highest-priority unblocked task
across all prioritised epics and unattached tasks. Sequential epics
contribute only their current lowest-sequence unfinished task, and only
when that task is unblocked.

---

## 5. Data Model — Entity Overview

| Entity | Purpose | Notes |
|---|---|---|
| `issues` | Unified table for epic / task / subtask. | `kind` discriminator. Holds shared fields. |
| `epic_metadata` | Epic-only structured fields (scheduling, brief path, target dates). | One-to-one with epic-kind issues. |
| `acceptance_criteria` | Per-task ACs; verification state lives in `ac_verifications`. | Replaces `- [ ]` markdown checkboxes. |
| `ac_verifications` | Append-only per-AC, per-role verification records. | Allows Developer / QA Reviewer / QA Lead / TL / Human to each verify independently. |
| `task_completions` | Append-only child of task-kind issues. One row per dev attempt. | Preserves history across re-runs. |
| `concerns` | Developer-raised concern attached to a task completion. | Triggers TL-Lite triage. `status` tracks lifecycle: open → addressed / escalated / wontfix. |
| `qa_reviews` | A QA pass. Polymorphic scope (`epic_window` / `single_task` / `epic_full`). | Linked to issue or epic via FK. May execute against a `qa_plan` (v0.3+). |
| `qa_findings` | Individual finding within a review. | `severity`, optional task subject, may escalate. |
| `qa_review_task_window` | Many-to-many: a review can cover N tasks. | Used when scope = `epic_window`. |
| `qa_plans` | QA verification design for an epic. (v0.3+) | Optionally cites a `validation_contract`. |
| `qa_plan_steps` | Individual verification steps within a plan. (v0.3+) | Display IDs `QA-NN`; evidence method per step. |
| `qa_plan_step_assertions` | Many-to-many: step ↔ assertion. (v0.4+) | A step may verify multiple assertions; assertions may have multiple covering steps. |
| `validation_contracts` | Assertion-level spec for an epic. (v0.4+) | Authored at planning; refined by QA Lead. Optional. |
| `validation_assertions` | Individual `VAL-NNN` assertions in a contract. (v0.4+) | Assertion text; required validators live in `validation_assertion_validators`. |
| `validation_assertion_validators` | Required verifier-role/evidence obligations per assertion. (v0.4+) | Allows one assertion to require Developer + QA Lead + TL + Human sign-off without duplicating the assertion. |
| `val_assertion_outcomes` | Append-only outcomes per assertion and validator role. (v0.4+) | `pass` / `fail` / `carry_forward`; required evidence or rationale. |
| `requirements` | PRD/TRD pointers (not copies). | `last_synced_hash` for drift detection. |
| `decisions` | `DEC-<hash>` records. | Body in markdown; row holds metadata. |
| `escalations` | `ESC-<hash>` records. Open or resolved. | Linked to triage / findings / decisions. |
| `triage_events` | `TRG-<hash>` records. TL-Lite verdicts. | `verdict` (RESUME / BLOCK), `esc_id?`. |
| `carry_forward` | CF items, with promotion-to-principle flag. | Active and scheduled commitments. |
| `retrospectives` | Epic / project retros. | `type` (lightweight / deep), recommendation. |
| `agent_runs` | Token telemetry per agent invocation. | `role`, `model`, in/out tokens, durations. |
| `templates` | Reusable templates with monotonic versioning. | `kind` enum, body in markdown with placeholders. |
| `workflows` | Named state machines. (v0.2+) | Per-kind or per-project assignment. In v0.1, workflows are hard-coded. |
| `workflow_states` | States within a workflow. (v0.2+) | `category` enables portable queries. |
| `workflow_transitions` | Allowed state-to-state transitions. (v0.2+) | One row per allowed edge. |
| `workflow_transition_requires` | Per-transition validation rules. (v0.2+) | E.g. `all_acs_verified` with `roles` params before `dev_complete → qa_review`. |
| `cycles` | Optional time-boxed sprint containers. | Off by default; opt-in per project. |
| `views` | Saved queries / filters. | DSL stored as text. |
| `comments` | Generic authored notes on issues, QA reviews, decisions, escalations, or requirements. | Scoped polymorphic ref. For context outside structured entities. |
| `activity_log` | Projected timeline of meaningful state transitions, assignments, etc. | Derived from journal events; query/read model only. |
| `labels` | Free-form tags applied to issues. | Multi-value via `issue_labels`. |
| `issue_labels` | Many-to-many between issues and labels. | |
| `links` | Scoped typed-edge table for traceability relationships. | v1 endpoints: issue, requirement, decision, escalation. Not for lifecycle relationships with dedicated columns. |
| `claims` | Active exclusive-access claim on supported entities. | v0.1: issues; v0.3+: QA reviews. At most one live claim per entity. TTL for orphan recovery. |
| `roles` | Set of valid actor roles for the project. | Seeded from the default profile at `fid init`. Existence here is the validity rule for all `_role` columns. |
| `role_permissions` | Per-role × action permissions matrix. (v0.2+) | In v0.1, permissions live in app code; migrated to this table in v0.2 (see ADR-0002). |

---

## 6. Schema Sketch

Indicative schema; column lists are illustrative and not exhaustive.

### 6.0 ID generation

Generated entity IDs follow the shape `<PREFIX>-<hash>` where:

- **Prefix** is per-entity-type: project prefix (set at `fid init`) for
  issues (`WHET-`, `FID-`); fixed prefixes for other generated
  entities (see table below); `PRD-`/`TRD-` retain their external
  document-derived identifiers (not generated by `fid`).
- **Hash** is `n` hex characters derived from `sha1(uuid4())[:n]`. Using
  a fresh UUID per call (not content-derived) avoids the "two issues
  with the same title collide" failure mode.
- **Progressive length scaling.** Default `n = 6`. When the count of
  rows in the relevant entity reaches a threshold that pushes
  birthday-paradox collision probability above ~1%, new IDs use
  `n = 8`. Old IDs remain valid.
- **Hierarchy is via `parent_id`, not encoded in the ID.** A subtask
  under task `WHET-a1b2c3` has its own hash like `WHET-7c3e9f`, not
  `WHET-a1b2c3.T01`. v1 hierarchy is strict: epics are top-level,
  tasks may belong to an epic or stand alone, and subtasks belong to
  a task.

Drop-in collision handling: on insert, if the generated ID already
exists, retry with a fresh UUID at the configured length. If a
collision appears only after journal merge, `fid rebuild` reports the
conflicting events and requires a human to rewrite one event's ID plus
its same-event references.

Generated prefixes:

| Prefix | Entity |
|---|---|
| `<project>` | Issue (`epic`, `task`, `subtask`) |
| `AC` | Acceptance criterion |
| `ACV` | Acceptance criterion verification |
| `TC` | Task completion |
| `CON` | Concern |
| `QAR` | QA review |
| `QAF` | QA finding |
| `DEC` | Decision |
| `ESC` | Escalation |
| `TRG` | Triage event |
| `CF` | Carry-forward item |
| `RET` | Retrospective |
| `RUN` | Agent run |
| `TPL` | Template |
| `CYC` | Cycle |
| `VIEW` | Saved view |
| `COM` | Comment |
| `ACT` | Activity-log entry |
| `VC` | Validation contract |
| `ASN` | Validation assertion storage row (`VAL-NNN` remains the display ID) |
| `QAP` | QA plan |
| `QAS` | QA plan step (`QA-NN` remains the display ID) |
| `VAO` | Validation assertion outcome |

Natural and composite keys stay natural: `requirements.id`
(`PRD-*` / `TRD-*`) comes from the source document, `labels.name`,
`roles.name`, `workflows.id` (stable slug, e.g. `task-default`),
`issue_labels(issue_id, label)`, `links(...)`, `claims(ref_type,
ref_id)`, `workflow_states(...)`, `workflow_transitions(...)`,
`workflow_transition_requires(...)`, `qa_review_task_window(...)`,
`qa_plan_step_assertions(...)`, `validation_assertion_validators(...)`,
and `role_permissions(...)`.

```sql
-- 6.1 Unified issues table

CREATE TABLE issues (
  id                  TEXT PRIMARY KEY,        -- e.g. WHET-a1b2c3 (prefix + hash; see §6.0)
  kind                TEXT NOT NULL,           -- 'epic' | 'task' | 'subtask'
  parent_id           TEXT REFERENCES issues(id),
                                               -- subtask's parent task,
                                               -- task's epic (optional),
                                               -- NULL for epics in v1
  title               TEXT NOT NULL,
  workflow_id         TEXT NOT NULL,           -- v0.1: enum string ('task-default' |
                                               --   'epic-default' | 'subtask-default').
                                               -- v0.2: FK to workflows(id).
  status              TEXT NOT NULL,           -- v0.1: CHECK-constrained per workflow_id.
                                               -- v0.2: FK to workflow_states(workflow_id, state).
  priority            INTEGER,                 -- 0..5 (0 = critical)
  sequence            INTEGER,                 -- positive sibling order; required for tasks
                                               -- under sequential epics, optional for subtasks,
                                               -- NULL for prioritised/unattached tasks
  spec_md_path        TEXT,                    -- canonical spec / brief
  est_tokens          INTEGER,                 -- task/subtask only
  due_date            DATE,
  in_inbox            BOOLEAN DEFAULT 0,       -- unsorted intake, awaiting assignment
  -- Note: concerns_raised and triage_status are NOT stored here.
  -- They are derived via the issues_enriched view (see §6.1.1).
  cycle_id            TEXT REFERENCES cycles(id),  -- optional
  owner_role          TEXT,                    -- long-lived ownership (≠ short-lived claim).
                                               -- Defaults: 'tl' for epics, 'developer' for tasks.
                                               -- Set at creation; mutable only via explicit
                                               -- override (journaled; projected to activity_log).
  owner_actor         TEXT,                    -- optional: 'agent' | 'human:<name>'.
                                               -- NULL means "any actor of owner_role".
  created_at          DATETIME, updated_at DATETIME,
  CHECK (kind IN ('epic', 'task', 'subtask')),
  CHECK (kind != 'epic' OR est_tokens IS NULL),  -- example kind-specific rule
  CHECK (sequence IS NULL OR sequence > 0),
  UNIQUE (parent_id, sequence)
);

CREATE TABLE epic_metadata (
  issue_id            TEXT PRIMARY KEY REFERENCES issues(id),
  scheduling          TEXT NOT NULL,           -- 'sequential' | 'prioritised'
  brief_md_path       TEXT,                    -- epic brief (whet's "phase brief" when sequential)
  target_start        DATE,
  target_end          DATE
);
```

#### 6.1.1 Derived columns via `issues_enriched`

Two values that callers frequently want alongside an issue —
"are there open concerns?" and "what's the latest triage verdict?"
— are not stored on the `issues` table. They genuinely live in
`concerns` and `triage_events` respectively, and storing them on
`issues` would introduce drift risk for no real benefit.

Instead, a SQL view `issues_enriched` LEFT JOINs the open-concerns
count and the latest `triage_events` row that targets each issue,
exposing the derived values as columns. The CLI's `fid issue get` and
`fid issue list` read from the view, not the base table.

Indicative shape (refined during implementation):

```sql
CREATE VIEW issues_enriched AS
WITH open_concerns AS (
  SELECT issue_id, COUNT(*) AS n
  FROM concerns
  WHERE status = 'open'
  GROUP BY issue_id
),
latest_triage AS (
  SELECT te.*
  FROM triage_events te
  WHERE te.created_at = (
      SELECT MAX(created_at)
      FROM triage_events
      WHERE issue_id = te.issue_id
    )
)
SELECT
  i.*,
  (COALESCE(oc.n, 0) > 0) AS concerns_raised,
  COALESCE(oc.n, 0)       AS open_concerns_count,
  CASE lt.verdict
    WHEN 'RESUME' THEN 'resumed'
    WHEN 'BLOCK'  THEN 'blocked'
    ELSE NULL
  END AS triage_status
FROM issues i
LEFT JOIN open_concerns oc ON oc.issue_id = i.id
LEFT JOIN latest_triage  lt ON lt.issue_id = i.id;
```

Filters like `fid issue list --has-concerns` translate to
`WHERE concerns_raised = 1` on the view. SQLite expands the view inline
and uses indexes on `concerns(issue_id, status)` and
`triage_events(issue_id, created_at)`, so the cost is two index
lookups per row in the result — not a table scan.

Indexes required for the view to perform well:
- `idx_concerns_issue_status` on `concerns(issue_id, status)`.
- `idx_triage_issue_created` on `triage_events(issue_id,
  created_at)`.

Because nothing is stored, there is no recompute or drift case. The
view is always consistent with its sources by construction. If view
performance ever becomes a bottleneck (it won't at single-project
scale), materialising via triggers is a forward-compatible upgrade
path.

#### 6.1.2 Remaining issue-side tables

```sql
CREATE TABLE acceptance_criteria (
  id                  TEXT PRIMARY KEY,        -- AC-<hash> (see §6.0)
  issue_id            TEXT NOT NULL REFERENCES issues(id),
  ordinal             INTEGER NOT NULL,
  text                TEXT NOT NULL
  -- Verification state lives in ac_verifications (append-only).
  -- Validation: only valid when issues.kind in ('task', 'subtask').
);

CREATE TABLE ac_verifications (
  id                  TEXT PRIMARY KEY,        -- ACV-<hash> (see §6.0)
  ac_id               TEXT NOT NULL REFERENCES acceptance_criteria(id),
  verifier_role       TEXT NOT NULL,           -- 'developer' | 'qa_reviewer' | 'qa_lead' | 'tl' | 'human'
  verifier_actor      TEXT,                    -- 'agent' | 'human:<name>' (optional)
  status              TEXT NOT NULL,           -- 'pass' | 'fail' | 'blocked' | 'not_applicable'
  evidence_type       TEXT,                    -- 'test' | 'source' | 'screenshot' |
                                               -- 'manual' | 'command' | 'commit'
  evidence_ref        TEXT,                    -- path / commit hash / command id / URL
  notes               TEXT,                    -- short, free-form
  completion_id       TEXT REFERENCES task_completions(id),  -- nullable; ties to attempt if relevant
  verified_at         DATETIME NOT NULL
);

-- Append-only by convention; UI never overwrites. A failed verification followed by a
-- passing one from the same role results in two rows. Latest-per-role queries use
-- (ac_id, verifier_role, MAX(verified_at)). Queries that need "verified by every required
-- role" (e.g. all_acs_verified with params_json roles) require at least one 'pass'
-- from each required role since the last completion attempt.

CREATE TABLE task_completions (
  id                  TEXT PRIMARY KEY,        -- TC-<hash> (see §6.0)
  issue_id            TEXT NOT NULL REFERENCES issues(id),
  attempt_number      INTEGER NOT NULL,
  developer_role      TEXT,                    -- 'developer' (for now)
  model               TEXT,                    -- e.g. 'sonnet-4.6'
  status              TEXT,                    -- 'dev_complete' | 'failed_gate' | 'superseded'
  notes_md_path       TEXT,                    -- prose: "what was implemented"
  -- Concerns raised during this attempt live as rows in the `concerns`
  -- table (see §6.4a), each with its own body_md_path and lifecycle.
  actual_tokens       INTEGER,
  commit_hashes       TEXT,                    -- JSON array
  files_touched       TEXT,                    -- JSON array
  expected_files_touched TEXT,                  -- JSON array (hint for v0.5+ file-scope claims)
  gates_passed        INTEGER DEFAULT 0,       -- count of passing gate runs in this attempt
  gates_failed        INTEGER DEFAULT 0,       -- count of failing gate runs in this attempt
  gate_results        TEXT,                    -- JSON array: [{cmd, exit_code, ran_at, output_excerpt}]
  human_intervention  BOOLEAN DEFAULT 0,       -- true if a human had to step in mid-attempt
  completed_at        DATETIME
);

-- 6.2 Workflows (v0.2+; in v0.1 these are hard-coded in application code)

CREATE TABLE workflows (
  id                  TEXT PRIMARY KEY,        -- stable slug, e.g. 'task-default'
  name                TEXT NOT NULL,
  applies_to          TEXT,                    -- 'epic' | 'task' | 'subtask' | 'qa_review' | NULL (any)
  active              BOOLEAN DEFAULT 1,
  created_at          DATETIME
);

CREATE TABLE workflow_states (
  workflow_id         TEXT NOT NULL REFERENCES workflows(id),
  state               TEXT NOT NULL,           -- 'todo' | 'in_progress' | 'qa_pass' | ...
  category            TEXT NOT NULL,           -- 'unstarted' | 'started' | 'completed' | 'cancelled'
  ordinal             INTEGER,
  PRIMARY KEY (workflow_id, state)
);

CREATE TABLE workflow_transitions (
  workflow_id         TEXT NOT NULL REFERENCES workflows(id),
  from_state          TEXT NOT NULL,
  to_state            TEXT NOT NULL,
  PRIMARY KEY (workflow_id, from_state, to_state)
);

CREATE TABLE workflow_transition_requires (
  workflow_id         TEXT NOT NULL,
  from_state          TEXT NOT NULL,
  to_state            TEXT NOT NULL,
  rule                TEXT NOT NULL,           -- 'all_acs_verified' | 'has_completion_report' |
                                               -- 'concerns_addressed' | ...
  params_json         TEXT,                    -- e.g. {"roles":["developer"]};
                                               -- NULL for non-parameterised rules
  FOREIGN KEY (workflow_id, from_state, to_state)
    REFERENCES workflow_transitions(workflow_id, from_state, to_state)
);

-- 6.3 QA

CREATE TABLE qa_reviews (
  id              TEXT PRIMARY KEY,            -- QAR-<hash> (see §6.0)
  scope           TEXT NOT NULL,               -- 'epic_window' | 'single_task' | 'epic_full'
  epic_id         TEXT REFERENCES issues(id),  -- when scope ∈ {'epic_window', 'epic_full'}
  task_id         TEXT REFERENCES issues(id),  -- when scope = 'single_task'
  reviewer_role   TEXT NOT NULL,               -- 'qa_reviewer' | 'qa_lead'
  model           TEXT,
  body_md_path    TEXT,                        -- review document
  created_at      DATETIME
);

CREATE TABLE qa_review_task_window (
  review_id       TEXT NOT NULL REFERENCES qa_reviews(id),
  task_id         TEXT NOT NULL REFERENCES issues(id),
  PRIMARY KEY (review_id, task_id)
);

CREATE TABLE qa_findings (
  id                  TEXT PRIMARY KEY,        -- QAF-<hash> (see §6.0)
  review_id           TEXT NOT NULL REFERENCES qa_reviews(id),
  severity            TEXT NOT NULL,           -- 'blocking' | 'improvement'
  task_id_subject     TEXT REFERENCES issues(id),
  status              TEXT NOT NULL,           -- 'open' | 'resolved' | 'wontfix' | 'escalated'
  escalated_to_esc    TEXT REFERENCES escalations(id),
  finding_md          TEXT,                    -- short prose; long-form in review body
  evidence_excerpt    TEXT                     -- literal CLI output / file excerpt as proof
);

-- 6.4a Concerns (developer-raised, structured)

CREATE TABLE concerns (
  id                          TEXT PRIMARY KEY,    -- CON-<hash> (see §6.0)
  issue_id                    TEXT NOT NULL REFERENCES issues(id),
  raised_in_completion_id     TEXT NOT NULL REFERENCES task_completions(id),
  raised_by_role              TEXT NOT NULL,       -- usually 'developer'
  summary                     TEXT NOT NULL,       -- one-line label
  body_md_path                TEXT,                -- prose: details, rationale
  status                      TEXT NOT NULL,       -- 'open' | 'addressed' | 'escalated' | 'wontfix'
  addressed_in_completion_id  TEXT REFERENCES task_completions(id),
  escalated_to_esc_id         TEXT REFERENCES escalations(id),
  triaged_in_event_id         TEXT REFERENCES triage_events(id),
  created_at                  DATETIME,
  updated_at                  DATETIME
);

-- A concern is the trigger for TL-Lite triage. Workflow Runner spawns
-- TL-Lite when issues_enriched.concerns_raised is true (i.e. one or
-- more concerns on the issue with status='open'). TL-Lite's verdict
-- (recorded in triage_events) updates the concern: RESUME flips it to
-- 'addressed' (with rationale on the triage event); BLOCK flips it to
-- 'escalated' and creates the escalation. A subsequent task_completion
-- that addresses the concern sets addressed_in_completion_id.
--
-- The 'concerns_addressed' workflow transition rule is satisfied when
-- no rows exist with (issue_id = ?, status = 'open').

-- 6.4 Decisions / escalations / triage

CREATE TABLE requirements (
  id                  TEXT PRIMARY KEY,        -- e.g. PRD-AUTH-001 / TRD-DB-003
  type                TEXT NOT NULL,           -- 'PRD' | 'TRD'
  summary             TEXT NOT NULL,           -- short label, NOT a copy
  source_doc          TEXT NOT NULL,
  source_section      TEXT,
  last_synced_hash    TEXT,                    -- hash of source section
  status              TEXT                     -- 'active' | 'deprecated'
);

CREATE TABLE decisions (
  id              TEXT PRIMARY KEY,            -- DEC-<hash> (see §6.0)
  epic_id         TEXT REFERENCES issues(id),
  author_role     TEXT NOT NULL,               -- 'tl' | 'human'
  summary         TEXT NOT NULL,
  body_md_path    TEXT NOT NULL,               -- full prose
  area            TEXT,                        -- 'schema' | 'cli' | 'workflow' | ...
  status          TEXT NOT NULL,               -- 'active' | 'superseded'
  created_at      DATETIME
);

CREATE TABLE escalations (
  id              TEXT PRIMARY KEY,            -- ESC-<hash> (see §6.0)
  status          TEXT NOT NULL,               -- 'open' | 'resolved'
  trigger_class   TEXT,                        -- enum from §2.5/§2.6 escalation rules
  raised_by_role  TEXT,                        -- 'tl_lite' | 'qa_reviewer' | ...
  raised_at       DATETIME,
  resolved_at     DATETIME,
  resolution_dec  TEXT REFERENCES decisions(id)
);

CREATE TABLE triage_events (
  id              TEXT PRIMARY KEY,            -- TRG-<hash> (see §6.0)
  epic_id         TEXT NOT NULL REFERENCES issues(id),
  issue_id        TEXT NOT NULL REFERENCES issues(id),
                                               -- the issue this triage pertains to;
                                               -- denormalized for the issues_enriched view
                                               -- (the join path from source_ref_id varies by
                                               -- source type, so a direct FK is cleaner)
  triager_role    TEXT NOT NULL,               -- 'tl_lite' | 'human'
  source          TEXT NOT NULL,               -- 'dev_concern' | 'qa_block' | 'esc_open'
  source_ref_id   TEXT,                        -- concern / finding / esc id (the specific trigger)
  verdict         TEXT NOT NULL,               -- 'RESUME' | 'BLOCK'
  esc_id          TEXT REFERENCES escalations(id),
  reasoning_md    TEXT,
  created_at      DATETIME
);

-- 6.5 Carry-forward / retrospectives

CREATE TABLE carry_forward (
  id                  TEXT PRIMARY KEY,        -- CF-<hash> (see §6.0)
  status              TEXT NOT NULL,           -- 'open' | 'scheduled' | 'closed' | 'principled'
  scheduled_epic_id   TEXT REFERENCES issues(id),
  promoted_principle  BOOLEAN DEFAULT 0,
  body_md             TEXT,
  created_at DATETIME, updated_at DATETIME
);

CREATE TABLE retrospectives (
  id              TEXT PRIMARY KEY,            -- RET-<hash> (see §6.0)
  epic_id         TEXT REFERENCES issues(id),
  type            TEXT NOT NULL,               -- 'lightweight' | 'deep'
  notes_md_path   TEXT NOT NULL,
  recommendation  TEXT,                        -- 'NOT_RECOMMENDED' | 'RECOMMENDED — <reason>'
  created_at      DATETIME
);

-- 6.6 Telemetry

CREATE TABLE agent_runs (
  id              TEXT PRIMARY KEY,            -- RUN-<hash> (see §6.0)
  issue_id        TEXT REFERENCES issues(id),
  completion_id   TEXT REFERENCES task_completions(id),
  review_id       TEXT REFERENCES qa_reviews(id),
  role            TEXT NOT NULL,               -- 'developer' | 'qa_reviewer' | 'qa_lead' |
                                               -- 'tl' | 'tl_lite' | 'workflow_runner'
  model           TEXT NOT NULL,
  tokens_in       INTEGER,
  tokens_out      INTEGER,
  price_in_per_mtok   REAL,                    -- USD per million input tokens at run end
  price_out_per_mtok  REAL,                    -- USD per million output tokens at run end
  cost_usd        REAL,                        -- computed at run end:
                                               -- (tokens_in  * price_in_per_mtok +
                                               --  tokens_out * price_out_per_mtok) / 1e6
  started_at      DATETIME,
  ended_at        DATETIME
);

-- Pricing source: .fid/pricing.toml maps model → (in_per_mtok, out_per_mtok).
-- Ships with defaults for current Anthropic/OpenAI/Google models; user-overridable.
-- If a model is missing from the config at run-close, cost_usd is left NULL
-- and a warning is logged (preferable to silently guessing).
-- Snapshot pricing on each row preserves audit trail across price changes.

-- 6.7 Templates

CREATE TABLE templates (
  id              TEXT PRIMARY KEY,            -- TPL-<hash> (see §6.0)
  kind            TEXT NOT NULL,               -- 'task_spec' | 'qa_review' | 'qa_plan' |
                                               -- 'human_checklist' | 'decision_record' |
                                               -- 'completion_report' | 'epic_brief' | ...
  name            TEXT NOT NULL,
  version         INTEGER NOT NULL,            -- monotonic per (kind, name)
  body_md         TEXT NOT NULL,               -- template content with {{placeholders}}
  active          BOOLEAN DEFAULT 1,
  created_at      DATETIME,
  UNIQUE (kind, name, version)
);

-- 6.8 Cycles, views, comments, activity log, labels

CREATE TABLE cycles (
  id              TEXT PRIMARY KEY,            -- CYC-<hash> (see §6.0)
  name            TEXT NOT NULL,               -- 'Cycle 12', 'Sprint 7', etc.
  start_date      DATE,
  end_date        DATE,
  status          TEXT NOT NULL                -- 'planned' | 'active' | 'completed'
);

CREATE TABLE views (
  id              TEXT PRIMARY KEY,            -- VIEW-<hash> (see §6.0)
  name            TEXT NOT NULL,               -- 'my-blocking-findings'
  description     TEXT,
  query_dsl       TEXT NOT NULL,               -- DSL or SQL
  scope           TEXT,                        -- 'project' | 'user'
  created_at      DATETIME
);

CREATE TABLE comments (
  id              TEXT PRIMARY KEY,            -- COM-<hash> (see §6.0)
  ref_type        TEXT NOT NULL,               -- 'issue' | 'qa_review' | 'decision' |
                                               -- 'escalation' | 'requirement'
  ref_id          TEXT NOT NULL,
  author_role     TEXT,                        -- 'human' | 'developer' | 'qa_reviewer' | ...
  body_md         TEXT NOT NULL,
  created_at      DATETIME
);

-- Comments are authored project state for notes that fall outside
-- structured concepts such as concerns, findings, and acceptance
-- criteria. If the content should trigger workflow, carry status, be
-- enforced, or feed reports, it should be modeled as a structured
-- entity instead. A comment-added activity row may be projected for
-- timelines, but the comment body itself is the authored state.
-- v1 comment scope is deliberately narrow: issue, qa_review, decision,
-- escalation, requirement. Low-level child rows (ACs, verifications,
-- links, labels, workflow rows, activity-log rows, etc.) use comments
-- on their parent entity instead.

CREATE TABLE activity_log (
  id              TEXT PRIMARY KEY,            -- ACT-<hash> (see §6.0)
  journal_event_id TEXT NOT NULL,              -- source journal event; no standalone writes
  ref_type        TEXT NOT NULL,
  ref_id          TEXT NOT NULL,
  event           TEXT NOT NULL,               -- 'status_changed' | 'ac_verified' | 'comment_added' | ...
  details         TEXT,                        -- JSON of before/after, or freeform
  actor           TEXT,                        -- role or 'human'
  occurred_at     DATETIME
);

CREATE TABLE labels (
  name            TEXT PRIMARY KEY,
  colour          TEXT,
  description     TEXT
);

CREATE TABLE issue_labels (
  issue_id        TEXT NOT NULL REFERENCES issues(id),
  label           TEXT NOT NULL REFERENCES labels(name),
  PRIMARY KEY (issue_id, label)
);

-- 6.9 Generic typed-edge

CREATE TABLE links (
  from_type       TEXT NOT NULL,               -- v1: 'issue' | 'requirement' |
                                               -- 'decision' | 'escalation'
  from_id         TEXT NOT NULL,
  relation        TEXT NOT NULL,               -- see relation enum below
  to_type         TEXT NOT NULL,
  to_id           TEXT NOT NULL,
  created_at      DATETIME,
  PRIMARY KEY (from_type, from_id, relation, to_type, to_id)
);
```

### 6.10 Link relation enum

| Relation | Inverse | Semantics |
|---|---|---|
| `blocks` | `blocked_by` | Issue → Issue. Hard dependency; target can't progress until source is done. |
| `relates_to` | `relates_to` | Soft association, no semantic. |
| `duplicates` | `duplicated_by` | Issue → Issue. This issue is the same as the target. |
| `supersedes` | `superseded_by` | Issue → Issue or Decision → Decision. This artifact replaces the target. |
| `implements` | `implemented_by` | Issue → Requirement. |
| `caused_by` | `causes` | Issue → Issue. One issue caused another (root-cause post-mortem links). |
| `discovered_from` | `discovered` | Issue → Issue. This issue was spotted while working on the target (context of discovery, not cause). Canonical for tech-debt-found-during-work. |

v1 link endpoints are deliberately scoped to `issue`, `requirement`,
`decision`, and `escalation`. The core use case is issue graphing
(`blocks`, `duplicates`, `supersedes`, `caused_by`,
`discovered_from`), issue → requirement implementation traceability,
issue/decision context, decision → decision supersession, and
issue/escalation context when an escalation affects work beyond the
row that created it.

Relation compatibility is part of validation, not convention. The CLI
and trigger backstop reject relation/entity pairs outside the table
above. `relates_to` is the only relation allowed across any two v1 link
endpoint types.

The generic link table is not a second encoding for lifecycle
relationships that already have dedicated schema. Do not use `links`
for hierarchy (`issues.parent_id`), QA review scope
(`qa_reviews.*` / `qa_review_task_window`), concern/finding escalation
(`concerns.escalated_to_esc_id`, `qa_findings.escalated_to_esc`),
escalation resolution (`escalations.resolution_dec`), AC ↔ assertion
coverage, or QA plan step ↔ assertion coverage.

Deferred endpoints: `concern`, `qa_finding`, `qa_review`,
`carry_forward`, and `validation_assertion`. Add them only when the
owning feature ships with precise relation names and no better
dedicated relationship table.

Hierarchy (subtask → task → epic) is expressed via `issues.parent_id`,
not through `links`. Keeping a single hierarchy mechanism avoids the
two-encodings-of-the-same-thing problem. v1 does not support nested
epics or initiatives above epics; those would need a deliberate
migration and explicit rollup semantics later.

Execution order for sequential epics is expressed via `issues.sequence`,
not through a chain of `blocked_by` links. A `blocked_by` link means the
target is genuinely invalid or unsafe to advance until the source is
complete; routine planned order belongs in `sequence`. Sequential epic
readiness is strict: a blocked lower-sequence task prevents later tasks
in the same epic from being returned by `fid ready`.

The CLI auto-creates the inverse on every link write so callers only
declare one direction.

### 6.11 Claims (exclusive-access leases)

```sql
CREATE TABLE claims (
  ref_type        TEXT NOT NULL,               -- v0.1: 'issue'; v0.3+: 'qa_review'
  ref_id          TEXT NOT NULL,
  claimed_by      TEXT NOT NULL,               -- role string: 'developer' | 'qa_reviewer' | 'human' | ...
  claimed_at      DATETIME NOT NULL,
  ttl_seconds     INTEGER NOT NULL DEFAULT 1800,  -- default 30 minutes
  note            TEXT,                        -- optional, e.g. "OIDC discovery flow"
  PRIMARY KEY (ref_type, ref_id)
);

-- A claim is "live" while (now - claimed_at) < ttl_seconds.
-- "Stale" otherwise; stale claims can be force-released by another caller.
-- Acquisition is atomic: INSERT ... ON CONFLICT DO NOTHING; zero rows = already claimed.
-- Polymorphic ref means no FK cascade; application/triggers clean on parent delete.
```

### 6.12 Claim semantics

Claims are short-lived exclusive-access leases on individual entity
rows. They sit alongside `status` and `activity_log`, not in place of
them:

- **Status** records what state the work is in.
- **Activity log** projects who did what and when for query/display.
- **Claim** records who currently holds the right to advance the work.
- **Ownership** (separate, on `issues.owner_role` /
  `issues.owner_actor`) records who is accountable for the work
  across its lifetime. Claims are short-lived (minutes); ownership is
  long-lived (an epic's full duration). An agent holding a claim does
  not become an owner — accountability stays with the owner (often
  the TL or a named human). This mirrors Linear's model: agents
  delegate; humans (or senior roles) remain accountable.

Rules:

- An entity may have at most one live claim at a time.
- Acquiring a claim is atomic. Concurrent `claim acquire` calls on the
  same entity: exactly one wins; the rest see "already claimed by X
  until T."
- Transitions into a `started`-category state require a live claim by
  the caller. Transitions out clear the claim automatically.
- A claim expires `ttl_seconds` after `claimed_at`. Stale claims do
  not block re-acquisition; the new caller's `claim acquire --force`
  (or implicit on stale) takes over and the prior holder is recorded
  in the journal and projected into `activity_log` as superseded.
- Humans can claim entities too (`--by human`) to signal "agent hands
  off while I think about this."

Entities supporting claims in v0.1: `issues`. v0.3 adds `qa_reviews`
when the QA review tables ship. The polymorphic `claims` table can
extend to other entities (`decisions`, `escalations`) without schema
change if a use case emerges.

### 6.13 Roles and permissions (actor policy)

The default workflow has hard role boundaries — Developer cannot create
decisions; QA Reviewer cannot mark ACs verified by QA Lead; TL-Lite
cannot rewrite task specs; Human owns product decisions. Claims
prevent concurrent advancement of the same entity; permissions prevent
the wrong role from advancing at all.

Per ADR-0002, the role list and permissions matrix together form a
project profile, seeded at `fid init` from a built-in default.

**v0.1 model.** The `roles` table is seeded at init with the default
profile's role list; existence in `roles` is the validity rule for any
`_role` column. The permissions matrix lives in app code as a constant.
Unknown roles passed via `--by` are hard-rejected with a clear "role X
not in profile; v0.2 will add `fid role create`."

```sql
CREATE TABLE roles (
  name            TEXT PRIMARY KEY,
  description     TEXT,
  is_human        BOOLEAN DEFAULT 0,           -- 'human' is conventionally true; others false
  active          BOOLEAN DEFAULT 1,
  created_at      DATETIME
);
```

**v0.2 model.** The matrix migrates to a `role_permissions` table
alongside the workflow-engine migration. CLI gains `fid role create`,
`fid policy set`, `fid policy show`. The default profile's matrix is
seeded into `role_permissions` verbatim during the schema migration;
existing v0.1 projects pick up the data form without change.

```sql
-- v0.2+
CREATE TABLE role_permissions (
  role_name       TEXT NOT NULL REFERENCES roles(name),
  action          TEXT NOT NULL,               -- e.g. 'issue.create.task', 'issue.transition.dev_complete'
  conditions      TEXT,                        -- optional restriction string, e.g. 'discovered_from'
  PRIMARY KEY (role_name, action)
);
```

Every write CLI invocation declares `--by <role>`. The CLI looks up
`(role, action)` (v0.1: in the matrix constant; v0.2: in the table);
if absent or denied, the command is rejected. A `fid policy show`
read command prints the active matrix in JSON or markdown so agents
can introspect their permitted actions.

**The default profile's role pairing pattern.** TL/TL-Lite and
QA Lead/QA Reviewer are intentional pairs: an Opus-class role
(`tl`, `qa_lead`) for bookends (epic kickoff + wrap-up; QA plan +
final epic review) and a Sonnet-class companion (`tl_lite`,
`qa_reviewer`) for routine mid-flight work (concern triage; every-
three-tasks reviews). The pairing keeps the capable model's context
budget for synthesis work and lets cheaper agents handle pattern-
matching adjudication. Custom profiles can collapse, split, or
rename these pairs.

**Default profile matrix (illustrative, not exhaustive):**

| Action | developer | qa_reviewer | qa_lead | tl_lite | tl | workflow_runner | human |
|---|---|---|---|---|---|---|---|
| issue create (task) | ✓ (discovered_from) | ✓ (from finding) | – | – | ✓ | ✓ | ✓ |
| issue create (epic) | – | – | – | – | ✓ | – | ✓ |
| issue transition → in_progress | ✓ | – | – | – | – | ✓ | ✓ |
| issue transition → dev_complete | ✓ | – | – | – | – | – | ✓ |
| issue transition → qa_pass | – | – | ✓ | – | – | – | ✓ |
| issue transition → done | – | – | – | – | ✓ | – | ✓ |
| issue ac verify --by developer | ✓ | – | – | – | – | – | – |
| issue ac verify --by qa_reviewer | – | ✓ | – | – | – | – | – |
| issue ac verify --by qa_lead | – | – | ✓ | – | – | – | – |
| issue ac verify --by tl | – | – | – | – | ✓ | – | – |
| issue ac verify --by human | – | – | – | – | – | – | ✓ |
| qa-review create | – | – | – | – | – | ✓ | ✓ |
| qa-finding add | – | ✓ | ✓ | – | – | – | – |
| concern raise | ✓ | – | – | – | – | – | ✓ |
| concern address | ✓ | – | – | – | – | – | ✓ |
| escalation create | – | ✓ | ✓ | ✓ | ✓ | – | ✓ |
| escalation resolve | – | – | – | – | ✓ | – | ✓ |
| triage record | – | – | – | ✓ | – | – | ✓ |
| decision create | – | – | – | – | ✓ | – | ✓ |
| retrospective create | – | – | – | – | ✓ | – | ✓ |
| comment add | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |

The `--by` flag is mandatory on every write command. The default
profile's matrix is the v0.1 source of truth; tests assert that every
CLI write action has a row in it and every denied combination raises
the right error.

### 6.14 Validation strategy

Schema constraints alone (CHECK, FK, UNIQUE) can't express every
invariant the workflow needs. Two examples:

- `acceptance_criteria.issue_id` FKs to `issues(id)` but can't constrain
  that the referenced row's `kind` is `task` or `subtask`.
- Workflow transition rules (`all_acs_verified`, `concerns_addressed`,
  etc.) need to read multiple tables and produce specific error
  messages.

The strategy is **application-layer primary, SQLite triggers as
backstop** (see OQ7 resolution). Each layer has a different role:

**Application layer (Go/TS) — primary validation, all writes:**

- Workflow transition requirements (`all_acs_verified`,
  `has_completion_report`, `concerns_addressed`,
  `all_child_tasks_done_or_cancelled`). These need rich error
  messages ("can't transition to dev_complete: 3 of 5 ACs not yet
  verified by developer: ac-2, ac-4, ac-5") that a flat trigger
  RAISE can't produce.
- Parameterised transition rules: `all_acs_verified` reads
  `params_json.roles` and can require any configured verifier-role set
  (`developer`, `qa_lead`, `tl`, `human`, etc.) without minting new
  rule names.
- Permissions matrix (§6.13) — needs `--by <role>` context that
  SQL can't see.
- AC verification freshness (verification belongs to current
  non-superseded attempt).
- Template structural validation (v0.3+) — needs markdown parsing.
- Claim acquire semantics — the table PK on `(ref_type, ref_id)`
  already enforces exclusivity; application handles TTL and stealing.

**SQLite triggers — backstop for writes that bypass the CLI:**

Why a backstop exists: future write paths will exist that don't go
through the CLI — the MCP server (v0.3), web dashboard (v0.5), ad-hoc
`fid sql` invocations, and direct DB edits during recovery. Triggers
guarantee that even those paths can't violate kind invariants and
silently corrupt the project state. The cost is small (~10 trigger
definitions); the failure mode without them ("future tool inserts
garbage; only `fid maintenance check` finds it later") is worse.

Trigger scope is deliberately narrow — only static structural rules
that SQL can express cleanly:

- **Kind-invariant rules.** BEFORE INSERT / UPDATE triggers on:
  - `issues.parent_id` — epics must have no parent; tasks may have no
    parent or an epic parent; subtasks must have a task parent
  - `issues.sequence` — tasks under sequential epics must have a
    sequence; tasks under prioritised epics or no epic must not;
    subtasks may use sequence for sibling display order, but v1 does
    not assign readiness semantics to subtask ordering
  - `acceptance_criteria` — `issue_id` must reference a `task` or
    `subtask`
  - `task_completions` — same
  - `epic_metadata` — `issue_id` must reference an `epic`
  - `qa_reviews` — `epic_id` references an epic when set; `task_id`
    references a task when set
  - `qa_review_task_window.task_id` — references a task
  - `triage_events.epic_id` — references an epic
  - `decisions.epic_id` — references an epic when set
  - `carry_forward.scheduled_epic_id` — references an epic when set
- **Polymorphic ref integrity.** BEFORE INSERT triggers on
  `links`, `claims`, `comments`, `activity_log` validating that
  `(ref_type, ref_id)` resolves to an existing row in the named
  table.
- **Link scope and relation compatibility.** BEFORE INSERT triggers on
  `links` reject endpoint types outside the v1 link scope (`issue`,
  `requirement`, `decision`, `escalation`) and reject relation/entity
  combinations that duplicate dedicated schema relationships.
- **Claim scope.** BEFORE INSERT / UPDATE triggers on `claims` reject
  unsupported `ref_type` values for the active schema version. v0.1
  supports only `issue`; v0.3 adds `qa_review`.
- **Comment scope.** BEFORE INSERT / UPDATE triggers on `comments`
  reject `ref_type` values outside the v1 comment scope (`issue`,
  `qa_review`, `decision`, `escalation`, `requirement`).
- **Activity-log projection discipline.** `activity_log` inserts are
  allowed only while applying journal event effects. It is a read-model
  table, not an independent write surface.
- **Kind immutability.** A BEFORE UPDATE OF kind trigger on `issues`
  rejects any attempt to change `kind` after creation. Trying to flip
  a row from `task` to `epic` post-hoc breaks every dependent row;
  forbid it at the schema level.
- **Mutually-exclusive scope columns on `qa_reviews`** — expressed as
  a CHECK constraint at column level rather than a trigger (CHECK is
  simpler when no cross-row lookup is needed):
  `CHECK ((scope = 'single_task' AND task_id IS NOT NULL AND epic_id
  IS NULL) OR (scope IN ('epic_window', 'epic_full') AND epic_id IS
  NOT NULL AND task_id IS NULL))`.

Trigger error messages are short and structural (`"acceptance_criteria.
issue_id must reference a task or subtask, got epic"`). The CLI catches
these and reformats with user-friendly framing; direct-SQL callers see
the raw message, which is acceptable because direct SQL is an escape
hatch.

What does **not** get a trigger backup: workflow transition rules,
permissions, AC verification freshness, template validation. These all
need rich context, multi-table reads with specific results, or runtime
state (calling role) that SQL doesn't have. They live in app-layer
code exclusively, with `fid maintenance check` as the periodic audit
(§9.8).

SQLite pragmas required at every connection:

```sql
PRAGMA foreign_keys = ON;       -- enforce FK constraints
PRAGMA recursive_triggers = ON; -- ensure triggers fire on cascaded updates
```

### 6.15 Validation contracts, QA plans, and assertion outcomes

Per ADR-0003, fid ships a layered quality model. The schema and CLI
support all three layers; a project picks the depth it needs.

#### 6.15.1 Layering and optionality

The three layers, from lightest to deepest:

1. **ACs + ad-hoc QA reviews** (v0.1 model). Per-task acceptance
   criteria with `ac_verifications`; QA reviews with `qa_findings`.
   No plans, no contracts. The lightest path; sufficient for most
   personal-scale projects.
2. **Add QA plans** (v0.3). `qa_plans` + `qa_plan_steps` capture the
   intended verification approach for an epic. Plans cite ACs;
   reviews cite a plan. Adds structured "how we'll verify" alongside
   "what we'll verify".
3. **Add validation contracts** (v0.4). `validation_contracts` +
   `validation_assertions` + `validation_assertion_validators` +
   `val_assertion_outcomes` capture the assertion-level specification
   independent of task decomposition. Plan steps cite assertions; ACs
   optionally cite assertions; per-assertion, per-validator outcomes
   feed pass/fail/carry-forward traceability.

**Layers are independently optional.** A project may use ACs only
(coverage rules inactive); ACs + plans (plans cite ACs directly);
ACs + plans + contracts (full traceability with coverage rules
enforced). Plans without contracts work — the plan effectively
becomes its own assertion+verification artifact. Contracts without
plans work — the assertions alone serve as the scope target. Reviews
work at every layer.

Coverage rules and traceability queries only activate when both a
contract and a plan are `active` for the same epic. A reader who only
needs ACs + reviews can stop reading §6.15 here.

#### 6.15.2 Why three entities

Production workflows that take quality seriously tend to converge on a
three-tier model:

- **Validation contract** — assertion-level specification. "What must
  be true for this scope of work." Authored at planning time, before
  task decomposition or test design. Stable target.
- **QA plan** — verification design. "How we prove the contract
  holds." Authored by QA Lead; cites contract assertions; declares
  evidence methods.
- **QA review** — execution. Findings, outcomes per assertion, polish
  evidence. Cites the plan it executes against.

This pattern shows up in Factory Missions (validation contract →
workers → validators), Cucumber/Gherkin (feature → scenarios → step
defs), regulated-industry workflows (requirements → verification
procedures → test results), and any project that needs traceability.
It generalises lighter workflows as a degenerate case: when no
contract exists, task ACs serve as the assertion layer directly.

#### 6.15.3 Schema

```sql
-- Validation contracts: the assertion-level spec for an epic.
CREATE TABLE validation_contracts (
  id              TEXT PRIMARY KEY,        -- VC-<hash> (see §6.0)
  epic_id         TEXT NOT NULL REFERENCES issues(id),
  body_md_path    TEXT,                    -- prose: rationale, scope notes, etc.
  status          TEXT NOT NULL,           -- 'draft' | 'active' | 'superseded'
  created_at      DATETIME
);

-- Individual assertions within a contract.
CREATE TABLE validation_assertions (
  id              TEXT PRIMARY KEY,        -- ASN-<hash> (storage ID; see §6.0)
  contract_id     TEXT NOT NULL REFERENCES validation_contracts(id),
  ordinal         INTEGER NOT NULL,        -- 1, 2, 3... per contract
  display_id      TEXT NOT NULL,           -- 'VAL-001' (composed, presentation-only)
  statement       TEXT NOT NULL,           -- the assertion text
  carry_forward_source TEXT,               -- prior display_id if carried forward
  out_of_scope_note TEXT,
  created_at      DATETIME,
  UNIQUE (contract_id, display_id)
);

-- Required validator/evidence obligations for an assertion.
CREATE TABLE validation_assertion_validators (
  assertion_id    TEXT NOT NULL REFERENCES validation_assertions(id),
  validator_role  TEXT NOT NULL,           -- 'developer' | 'qa_reviewer' |
                                           -- 'qa_lead' | 'tl' | 'human'
  evidence_method TEXT NOT NULL,           -- 'widget_test' | 'repository_test' |
                                           -- 'manual_checklist' | 'screenshot' |
                                           -- 'command' | 'source_state' |
                                           -- 'source_review' | ...
  PRIMARY KEY (assertion_id, validator_role, evidence_method)
);

-- QA plan header.
CREATE TABLE qa_plans (
  id              TEXT PRIMARY KEY,        -- QAP-<hash> (see §6.0)
  epic_id         TEXT NOT NULL REFERENCES issues(id),
  contract_id     TEXT REFERENCES validation_contracts(id),  -- nullable
  body_md_path    TEXT,
  status          TEXT NOT NULL,           -- 'draft' | 'active' | 'superseded'
  created_at      DATETIME
);

-- Steps within a plan.
CREATE TABLE qa_plan_steps (
  id              TEXT PRIMARY KEY,        -- QAS-<hash> (see §6.0)
  plan_id         TEXT NOT NULL REFERENCES qa_plans(id),
  ordinal         INTEGER NOT NULL,
  display_id      TEXT NOT NULL,           -- 'QA-01' composed
  description     TEXT NOT NULL,
  evidence_method TEXT,                    -- same enum as required validator rows
  validator_role  TEXT NOT NULL,
  UNIQUE (plan_id, display_id)
);

-- Many-to-many: a step may verify multiple assertions; an assertion may
-- be verified by multiple steps (defensive depth).
CREATE TABLE qa_plan_step_assertions (
  step_id         TEXT NOT NULL REFERENCES qa_plan_steps(id),
  assertion_id    TEXT NOT NULL REFERENCES validation_assertions(id),
  PRIMARY KEY (step_id, assertion_id)
);

-- Append-only outcomes: per-assertion, per-validator verdict at review or wrap-up time.
CREATE TABLE val_assertion_outcomes (
  id              TEXT PRIMARY KEY,        -- VAO-<hash> (see §6.0)
  assertion_id    TEXT NOT NULL REFERENCES validation_assertions(id),
  review_id       TEXT REFERENCES qa_reviews(id),
  outcome         TEXT NOT NULL,           -- 'pass' | 'fail' | 'carry_forward'
  evidence_method TEXT,                    -- actual submitted evidence method
  evidence_ref    TEXT,                    -- path/commit/command-result; required when outcome='pass'
  rationale       TEXT,                    -- required when outcome='carry_forward'
  recorded_by_role TEXT NOT NULL,
  recorded_at     DATETIME NOT NULL
);

-- Acceptance criteria optionally cite assertions (v0.4 column addition).
ALTER TABLE acceptance_criteria
  ADD COLUMN val_assertion_id TEXT REFERENCES validation_assertions(id);

-- QA reviews optionally execute against a plan (v0.3 column addition).
ALTER TABLE qa_reviews
  ADD COLUMN plan_id TEXT REFERENCES qa_plans(id);
```

#### 6.15.4 Coverage rules (enforced when contract+plan present)

When a contract is `active` and a plan is `active` for the same epic,
fid enforces these invariants at plan-finalisation time
(`fid qa-plan finalise`):

- Every requirement linked to the epic (via `links` with relation
  `implements`) is referenced by at least one assertion.
- Every assertion has at least one referencing AC OR an explicit
  `out_of_scope_note` for non-implementation-bound assertions.
- Every assertion has at least one required validator row.
- Every required validator row has at least one linked `qa_plan_step`
  for the same assertion with matching `validator_role` and
  `evidence_method`.

`PASS` outcomes require non-null `evidence_method` and `evidence_ref`.
`CARRY_FORWARD` outcomes require non-null `rationale`. A required
validator row is satisfied by a latest `PASS` outcome whose
`assertion_id`, `recorded_by_role`, and `evidence_method` match the
required row. These are trigger-enforced on `val_assertion_outcomes`
insert (§6.14).

#### 6.15.5 Traceability queries

Coverage and trace queries are named reports rather than ad-hoc DSL
(see §9.5):

- `fid report contract-coverage <epic-id>` — gaps in coverage:
  unreferenced requirements, ACs without assertions, assertions
  without ACs, assertions without required validators, required
  validators without matching plan steps.
- `fid report assertion-trace <assertion-id>` — full trace: source
  requirements → assertion → required validators → implementing ACs
  → covering plan steps → recorded outcomes.
- `fid report contract-outcomes <epic-id>` — pass/fail/carry-forward
  rollup by assertion, required validator role, and evidence method;
  the §13 step 7 wrap-up signal.

### 6.16 Implementation stack

Resolved per OQ6: **Go**. The agentic-build pipeline proven on whet
transfers directly, distribution is a single static binary across
all targets, and the agent error-recovery loop is more efficient than
in Rust for a design-iteration period.

Indicative library picks (final choices may shift during v0.1 build):

- **SQLite driver:** `modernc.org/sqlite` (pure Go, no CGo). Avoids
  the CGo cross-compile complexity; speed is more than adequate for
  fid's workload. Fall back to `mattn/go-sqlite3` if performance ever
  becomes a real bottleneck.
- **CLI framework:** `spf13/cobra` for the nested command tree (`fid
  issue ac verify`, `fid qa-plan step add`, etc.) plus `spf13/pflag`
  for POSIX-style flags. Mature, well-known, agent-friendly patterns.
- **Migrations:** `golang-migrate/migrate` with file-based SQL
  migrations under `migrations/`. Schema evolves via numbered
  up/down pairs; CLI exposes `fid maintenance migrate` for explicit
  control.
- **JSON / TOON:** stdlib `encoding/json` for JSON; in-tree TOON
  marshaller (small, no third-party dependency that might break with
  spec revisions — see §10.3).
- **Markdown parsing:** `yuin/goldmark` for §11 export rendering,
  template section parsing, and the v0.3 template-driven validation.
  Extensible via custom node parsers if needed.
- **Frontmatter:** `adrg/frontmatter` or equivalent for YAML
  frontmatter read/write on v0.4 export.
- **TUI (v0.3 board):** `charmbracelet/bubbletea` with `lipgloss`
  for styling. The de facto Go TUI stack; mature and a pleasure to
  write.
- **MCP server (v0.3):** community Go MCP SDK; specific package
  chosen at v0.3 implementation time based on what's most stable
  then.
- **Testing:** stdlib `testing` + `testify` for assertions and
  table-driven test ergonomics. Integration tests use a temporary
  SQLite DB per test.

**Project layout (indicative):**

```
fid/
├── cmd/fid/                 # main package; cobra command tree
├── internal/
│   ├── db/                  # schema, migrations, query helpers
│   ├── issues/              # issue, AC, completion logic
│   ├── claims/              # claim acquire/release/steal
│   ├── workflow/            # state machines + transition rules (v0.1 hard-coded)
│   ├── policy/              # role × action permissions matrix (§6.13)
│   ├── qa/                  # qa_reviews, qa_findings (v0.3)
│   ├── validation/          # contracts, assertions, outcomes (v0.4)
│   ├── templates/           # template engine + structural validation (v0.3)
│   ├── export/              # markdown export with frontmatter (v0.4)
│   ├── journal/             # JSONL writer + replayer (both v0.1)
│   ├── stats/               # cost, tokens, gates, attempts aggregations
│   ├── report/              # named reports (v0.2+)
│   └── format/              # JSON / TOON / markdown / YAML output
├── migrations/              # numbered SQL migrations
├── skills/                  # built-in agent skills (v0.3+)
├── templates/               # built-in artifact templates (v0.3+)
└── testdata/                # fixture DBs and golden files

```

The `internal/` boundary keeps the API surface intentional —
everything is private by default; only `cmd/fid/` calls into the
domain packages.

---

## 7. Workflows and State Transitions

### 7.1 Workflow concept

A **workflow** is a named state machine that applies to issues of a given
`kind` (or to QA reviews, etc.). Each workflow has:

- **States** with a `category` mapping (`unstarted` / `started` /
  `completed` / `cancelled`). Categories enable portable queries.
- **Transitions** — explicit allowed edges in the state graph.
- **Transition requirements** — per-edge validation rules enforced at
  write time.

The CLI rejects invalid jumps. Re-attempts of a task create a new
`task_completions` row; the previous one is marked `superseded`.

**v0.1 implementation.** The workflows in §7.2–§7.5 are hard-coded in
the application. `issues.workflow_id` is an enum string
(`'task-default'`, `'epic-default'`, `'subtask-default'`), and `status`
is CHECK-constrained against the legal states for that workflow.
Transition rules are functions in code (one per rule name listed in
§16 v0.2 prep).

**v0.2 implementation.** The four workflow tables (`workflows`,
`workflow_states`, `workflow_transitions`,
`workflow_transition_requires`) replace the hard-coded approach. A
schema migration seeds them with the v0.1 defaults verbatim, then
`fid workflow define` and `fid workflow assign` (§7.6) become
available. The category mapping is preserved across the migration, so
existing queries against `status_category` keep working.

### 7.2 Default workflow — task

Ships with `fid` and applies to `kind='task'` unless overridden.

```
States:
  todo            (unstarted)
  in_progress     (started)
  dev_complete    (started)
  qa_review       (started)
  qa_block        (started)
  failed_gate     (started)
  qa_pass         (completed)
  done            (completed)
  cancelled       (cancelled)

Transitions:
  todo            → in_progress
  in_progress     → dev_complete  [requires: all_acs_verified{roles:[developer]}, has_completion_report]
  in_progress     → failed_gate
  failed_gate     → in_progress
  dev_complete    → qa_review
  qa_review       → qa_pass       [requires: all_acs_verified{roles:[qa_lead]}]
  qa_review       → qa_block
  qa_block        → in_progress
  qa_pass         → done
  *               → cancelled
```

### 7.3 Default workflow — epic

```
States:
  todo            (unstarted)
  in_progress     (started)
  wrapped         (started)        [TL has completed wrap-up]
  closed          (completed)      [human approval recorded]
  cancelled       (cancelled)

Transitions:
  todo        → in_progress
  in_progress → wrapped   [requires: all_child_tasks_done_or_cancelled]
  wrapped     → closed
  *           → cancelled
```

### 7.4 Default workflow — subtask

Subtasks use a simplified workflow by default:

```
States:
  todo            (unstarted)
  in_progress     (started)
  done            (completed)
  cancelled       (cancelled)

Transitions:
  todo        → in_progress
  in_progress → done
  *           → cancelled
```

### 7.5 Default workflow — escalation

```
States:
  open      (started)
  resolved  (completed) [requires: linked_to_decision]
```

### 7.6 Customisation (v0.2+)

Users can replace any default workflow per project, per kind, or per
project-and-kind:

```
fid workflow define --name dev-flow-strict --applies-to task \
    --states "todo,in_progress,dev_complete,qa_review,qa_pass,qa_block,done" \
    --categories "unstarted,started,started,started,completed,started,completed" \
    --transitions "todo→in_progress,in_progress→dev_complete,..." \
    --require 'in_progress→dev_complete:all_acs_verified{"roles":["developer"]}'
fid workflow assign --name dev-flow-strict --kind task
```

---

## 8. Templates

Templates encode reusable artifact shells (task spec format, QA review
shell, decision record skeleton, etc.) with monotonic versioning.

```
fid template create --kind task_spec --name default-feature-task \
    --from path/to/template.md
fid template list
fid template show task_spec/default-feature-task --version latest

fid issue create --kind task --from-template default-feature-task \
    --title "Add OIDC support"
```

Versioning is monotonic (v1, v2, v3…). Old artifacts retain a reference
to the template version they were created from, so subsequent template
changes don't retroactively alter history. Set `active=0` to retire a
template without breaking historical references.

Default templates ship with the binary (good defaults); project
templates override.

### 8.1 Provenance (v0.3+)

When an entity is created from a template, the template kind, name,
and version are recorded on the entity row so later queries can
distinguish artifacts produced by different template revisions.

Provenance columns are added at v0.3 (the release that ships templates
themselves) to the following entities only:

- `issues` (from `task_spec` / `subtask_spec` / `epic_brief`)
- `task_completions` (from `completion_report`)
- `qa_reviews` (from `qa_review` / `qa_plan`)
- `decisions` (from `decision_record`)

Lower-volume entities (`retrospectives`, `triage_events`) do not carry
provenance columns; senior-role authors can include the template
version in the body if it ever matters.

Column shape on each of the entities above:

```sql
created_from_template_kind     TEXT,     -- e.g. 'task_spec'
created_from_template_name     TEXT,     -- e.g. 'default-feature-task'
created_from_template_version  INTEGER,  -- e.g. 2
```

These are nullable (entities created without a template — direct CLI
calls — leave them NULL). They're soft references, not FK to
`templates(kind, name, version)`, so retiring a template (setting
`active=0`) or even deleting it does not cascade.

Built-in template kinds:
- `task_spec` — matches §5 of the whet AI workflow doc
- `subtask_spec` — lighter shape for subtasks
- `qa_plan` — QA Lead epic plan
- `qa_review` — mid-epic or final epic review
- `human_checklist` — human verification subset
- `decision_record` — decision body template
- `completion_report` — Developer completion report including the
  required `## Concerns` section
- `epic_brief` — epic brief
- `triage_log_entry` — triage event entry skeleton
- `retro_lightweight` — §9.1 workflow notes template

---

## 9. CLI Surface (sketch)

Indicative; final command shapes will be refined during implementation.

### 9.1 Initialisation

```
fid init                                     # in current project root
fid init --prefix WHET                       # set issue ID prefix (default: project name)
fid init --import whet-private-docs/         # v0.5+ conditional; see §16
```

### 9.2 Reads (agent-friendly)

```
fid issue get <id>                                    # full record + ACs
fid issue get <id> --section acceptance               # parsed section only
fid issue get <id> --section concerns
fid issue get <id> --format json|toon|markdown
fid issue get <id> --context-pack developer           # v0.2+: role-tuned bundle
fid issue get <id> --context-pack qa_reviewer         # v0.2+
fid issue get <id> --context-pack tl                  # v0.2+

fid issue list --status-category started --kind task
fid issue list --label tech-debt --priority 0,1
fid issue list --epic <id> --order sequence|priority

fid ready                                             # cross-project ready queue
fid ready --epic <id>                                 # next-in-sequence or top-priority per epic scheduling

fid requirement get PRD-AUTH-001
fid decision get DEC-a1b2c3
fid esc list --status open
fid qa-review list --epic <id>
fid qa-finding list --epic <id> --severity blocking --status open

# v0.3+ — QA plans
fid qa-plan get <plan-id>
fid qa-plan list --epic <id>
fid qa-plan step list <plan-id>

# v0.4+ — Validation contracts and assertions
fid contract get <epic-id>                           # the active contract for an epic
fid assertion get VAL-001 --epic <epic-id>
fid assertion list --epic <epic-id> [--validator-role qa_lead]
fid assertion list --uncovered --epic <epic-id>      # no AC ref, no required validator,
                                                     # no matching plan step, etc.
fid outcome list --assertion VAL-001
fid outcome list --epic <epic-id> --status fail      # epic wrap-up signal

fid view run my-blocking-findings                     # saved query
```

#### 9.2.1 Context packs (v0.2+)

`fid issue get <id> --context-pack <role>` returns a role-tuned bundle
of the issue plus everything that role needs and nothing more.
Reduces context window pressure and removes the "what should I read
first?" question.

Built-in packs (final composition refined during implementation):

- `developer` — issue + ACs + spec_md + linked requirements + open
  comments + relevant template + active claim status + parent epic
  brief excerpt
- `qa_reviewer` — issue + ACs (with latest verifications per role) +
  spec_md + completion notes + linked findings from prior reviews
- `qa_lead` — same as qa_reviewer plus human-checklist generation
  hooks and any stricter evidence requirements declared in the
  qa_review template for this project
- `tl_lite` — concern excerpt + originating task spec + active
  escalations referencing this epic + relevant decisions
- `tl` — full epic view: child task statuses, all decisions, open
  escalations, all triage events, retro draft if any

Custom packs (`fid issue get <id> --context-pack-spec path.toml`) are
a v0.3+ consideration; v0.2 ships the built-ins only.

### 9.3 Writes (agent-driven)

```
fid issue create --kind task --title "..." --epic <id> --priority 1 \
    --sequence 10 --spec-md path/to/spec.md --label bug --label auth

fid issue ac add <id> --text "..." --ordinal 1
fid issue ac verify <id> <ac-id> --by developer --status pass \
    --evidence-type test --evidence-ref test/foo_test.dart::test_x
fid issue ac verify <id> <ac-id> --by qa_lead --status pass \
    --evidence-type command --evidence-ref "flutter test --tags=auth"
fid issue ac list <id>                                # all ACs + latest-per-role verifications

fid issue transition <id> --to dev_complete           # validates rules
fid issue complete <id> \                             # convenience composite
    --notes-md path/to/notes.md \
    --tokens 14823 \
    --commits abc123,def456
# Concerns raised during this attempt are recorded separately:
#   fid concern raise --on <issue-id> --by developer \
#       --in-completion <completion-id> --summary "..." --body-md path
# v0.3+ template-driven: `fid dev finish` parses the `## Concerns`
# section of the completion report and creates concern rows
# automatically.

fid qa-review create --scope epic_window --epic <id> \
    --tasks t1,t2,t3 --reviewer qa_reviewer \
    --plan <plan-id>                                  # v0.3+: cite plan if any
fid qa-finding add <review-id> --severity blocking --task <task-id> \
    --finding-md path

# v0.3+ — QA plan authoring
fid qa-plan create --epic <id> --by qa_lead \
    --body-md path \
    --contract <contract-id>                          # nullable; v0.4+ once contracts exist
fid qa-plan step add <plan-id> --description "..." \
    --evidence-method widget_test --validator-role qa_reviewer
fid qa-plan step link <step-id> --assertion VAL-001   # v0.4+: cite assertions
fid qa-plan finalise <plan-id>                        # asserts coverage rules pass (§6.15.4)

# v0.4+ — Validation contracts
fid contract create --epic <id> --by tl --body-md path
fid contract activate <contract-id>
fid assertion add <contract-id> \
    --statement "Logged-in user can see their orders"
fid assertion validator add <assertion-id> \
    --role developer --evidence-method repository_test
fid assertion validator add <assertion-id> \
    --role qa_lead --evidence-method widget_test
fid assertion validator add <assertion-id> \
    --role tl --evidence-method source_review
fid assertion validator add <assertion-id> \
    --role human --evidence-method manual_checklist
fid assertion link-ac <assertion-id> <ac-id>          # link an AC to the assertion it implements
fid outcome record <assertion-id> --review <review-id> \
    --status pass --evidence-method widget_test \
    --evidence-ref "test/orders_test.dart::shows_orders" \
    --by qa_lead
fid outcome record <assertion-id> --review <review-id> \
    --status pass --evidence-method source_review \
    --evidence-ref "reviews/orders-source-review.md" \
    --by tl
fid outcome record <assertion-id> --review <review-id> \
    --status carry_forward --rationale "Drift migration needed; deferred to a future epic" \
    --by tl

fid concern raise --on <issue-id> --by developer \
    --in-completion <completion-id> --summary "..." \
    --body-md path/to/details.md
fid concern address <concern-id> --in-completion <completion-id>
fid concern list --issue <issue-id> [--status open]

fid escalation create --from-finding <finding-id> \
    --trigger-class arch_violation
fid escalation create --from-concern <concern-id> \
    --trigger-class arch_violation
fid triage record --epic <id> --source dev_concern --ref <concern-id> \
    --verdict RESUME --reasoning-md path
# --ref points at the source's row id:
#   --source dev_concern → concerns.id
#   --source qa_block    → qa_findings.id
#   --source esc_open    → escalations.id
fid decision create --summary "..." --body-md path --area schema

fid link add WHET-a1b2c3 blocks WHET-7c3e9f            # auto-creates inverse
fid link add WHET-a1b2c3 implements PRD-AUTH-001
fid link list WHET-a1b2c3                              # all links from/to WHET-a1b2c3

fid comment add --on issue/WHET-a1b2c3 --body-md path

# Claims (exclusive-access leases on issues; qa_reviews in v0.3+)
fid issue claim WHET-a1b2c3 --by developer --ttl 30m   # acquire
fid issue unclaim WHET-a1b2c3                          # release (live or stale)
fid issue unclaim WHET-a1b2c3 --force                  # take over a live claim
fid qa-review claim QAR-7c3e9f --by qa_reviewer        # v0.3+
fid claim list                                         # all live claims, any entity
fid claim list --stale                                 # past-TTL claims (recoverable)
fid claim list --by developer                          # mine
```

### 9.3a Role-namespace composite commands

The `fid issue ...` namespace exposes primitives. Role-namespace
commands compose those primitives into the right sequence for a given
role's workflow, so agents do one call instead of remembering five.
Permissions (§6.13) apply automatically — each composite asserts
`--by <role>`.

```
# Developer
fid dev start <issue-id> --by developer                # claim + transition→in_progress + open agent_run
fid dev gate <issue-id> --cmd "flutter test --tags=auth" \
    --by developer                                     # record gate result on current attempt
fid dev finish <issue-id> --by developer \             # close agent_run + create task_completion +
    --completion-md path \                             # parse `## Concerns` into concerns rows (v0.3+);
    --tokens 14823 --commits abc123 --files lib/x.dart # transition→dev_complete (validates rules)
fid dev block <issue-id> --by developer \              # transition→failed_gate + completion w/ reason
    --reason-md path

# QA Reviewer / QA Lead
fid qa start <epic-id> --scope epic_window --by qa_reviewer
fid qa finding <review-id> --by qa_reviewer \
    --severity blocking --task <task-id> --finding-md path
fid qa verify <issue-id> <ac-id> --by qa_lead \
    --status pass --evidence-type command \
    --evidence-ref "flutter test --tags=auth"
fid qa finish <review-id> --by qa_lead                 # transition the review to completed

# TL-Lite
fid triage start <epic-id> --source dev_concern --ref <concern-id> --by tl_lite
fid triage verdict --verdict RESUME --by tl_lite \
    --reasoning-md path
fid triage verdict --verdict BLOCK --by tl_lite \
    --escalation-trigger arch_violation --reasoning-md path

# Tech Lead
fid tl wrapup <epic-id> --by tl                        # generate retrospective row +
                                                       # close resolved esc + transition→wrapped
```

Each composite is a thin wrapper over `fid issue ...` / `fid qa-review
...` / etc., with the right argument plumbing and an atomic
transaction so partial failures don't leave half-applied state.
Composites are convenience, not gate: agents *may* call the primitives
directly when needed.

```
fid template create --kind task_spec --name default-feature-task \
    --from path.md
fid template list
fid template show task_spec/default-feature-task

# v0.2+ — workflow engine
fid workflow define --name ...
fid workflow assign --name ... --kind task
fid workflow list
```

### 9.5 Named reports, stats, and queries

```
# Named reports (v0.2+) — opinionated wrappers over hand-coded SQL.
# These cover the high-frequency cross-cutting needs.
fid report epic-status [--epic <id>]                # tasks/ACs/gates/findings/escalations summary
fid report open-escalations
fid report qa-blockers [--epic <id>]
fid report ac-matrix [--epic <id>]                  # AC × verifier-role grid
fid report contract-coverage <epic-id>              # v0.4+: gaps in VAL ↔ validator ↔ AC ↔ QA step coverage
fid report assertion-trace <assertion-id>           # v0.4+: requirement → assertion → validators → ACs → steps → outcomes
fid report contract-outcomes <epic-id>              # v0.4+: outcomes by assertion/validator/evidence
fid report epic-wrapup-pack <epic-id>               # everything TL needs for wrap-up
fid report task-hygiene <task-id>                   # missing headings, commit refs, AC verifications
fid report stale-claims                             # past-TTL active claims
fid report list                                     # all available reports

# Stats — quantitative aggregations
fid stats tokens --by role --by model
fid stats cost   --by role --by model               # USD; uses snapshot pricing on rows
fid stats cost   --by epic --include est-vs-actual
fid stats tasks --by epic --include est-vs-actual
fid stats attempts --by epic                        # avg attempts per task, rework rate
fid stats gates --by epic                           # pass/fail counts, common failures
fid stats interventions --by epic                   # human intervention frequency
fid stats findings --by epic --by severity

# Ad-hoc query (v0.3+) — small filter DSL for cases reports don't cover.
fid query "decisions where area=schema since WHET-a1b2c3"
fid query "tasks where actual_tokens > 1.5 * est_tokens"
fid query "blocking findings still open across all epics"

# SQL escape hatch.
#   Read-only by default: any statement other than SELECT/EXPLAIN/PRAGMA
#   query/WITH is rejected. The escape hatch is for inspection, not
#   editing — writes go through CLI commands so the journal records
#   the event and the schema rules run.
#   Allowed for any role (read-only is harmless). SQL reads are not
#   activity-log events; command audit metadata lives in the journal.
#   Direct DB edits during recovery (sqlite3 .fid/fid.sqlite) bypass
#   the journal and are a "you broke it, you fix it" path; fid
#   maintenance check is the safety net.
fid sql "SELECT ..."
```

**Why reports first, DSL later.** Named reports deliver more value
with less effort: they're opinionated, well-tested SQL behind a
memorable name, and they cover the common workflow needs (epic
status, wrap-up pack, hygiene checks). A generic query DSL is its own
product — designing the grammar, parser, error messages, and docs
takes real effort that only pays off for the *uncovered* tail of
queries. Build the reports first; the DSL can land in v0.3 once you
know which queries don't fit a report shape.

### 9.6 Export (and future sync)

```
# v0.4+ — export only (markdown is downstream of the DB)
fid export --to ./docs/                         # full project tree
fid export --to ./docs/ --epic WHET-a1b2c3      # one epic
fid export --to ./docs/ --since 2026-04-01      # incremental

# v0.5+ — deferred; ship only if a real use case emerges
fid import --from ./docs/                       # reverse direction
fid sync                                        # bidirectional reconcile (config-aware)
```

### 9.7 Visualisation

```
fid board                                       # TUI kanban (v0.3+)
fid board --epic WHET-a1b2c3 --by category
fid roadmap                                     # epic timeline (when target dates set)
```

The web dashboard ships separately as `fid-dashboard` (v0.5+); reads the
same SQLite database.

### 9.8 Maintenance

```
fid maintenance vacuum                          # SQLite VACUUM
fid maintenance check                           # consistency audit; non-zero exit on errors
fid maintenance check --warn-as-error           # promote warnings to errors (CI mode)
fid maintenance check --format json|markdown    # output format (default JSON per §10)
```

#### 9.8.1 `fid maintenance check` scope

`check` is the periodic audit that complements the trigger backstop
(§6.14). It looks for inconsistencies that should never exist if all
writes went through proper paths, and serves as both the recovery
tool and the pre-merge / CI safety net. Findings are categorised:

**Errors** (non-zero exit code; the DB is in an inconsistent state):

- **Kind invariant violations.** Issue hierarchy violations
  (`epic.parent_id` set, task parent not an epic, subtask parent not a
  task), sequence violations (task missing sequence under a sequential
  epic, task sequence present under prioritised/no epic), plus rows in
  `acceptance_criteria`, `task_completions`, `epic_metadata`,
  `qa_reviews`, `qa_review_task_window`,
  `triage_events`, `decisions`, `carry_forward` whose parent issue has
  the wrong `kind`. Triggers prevent these on write; the check is
  belt-and-suspenders for any historical or bypassed data.
- **Orphaned polymorphic refs.** Rows in `links`, `claims`,
  `comments`, `activity_log` whose `(ref_type, ref_id)` doesn't
  resolve to an existing row in the indicated table, plus
  `activity_log` rows whose `journal_event_id` is not present in the
  journal.
- **FK reference integrity.** Trust-but-verify pass over all FK
  columns — should never fail given `PRAGMA foreign_keys = ON`, but
  worth reporting if it does.
- **Workflow state validity.** Any `issues.status` that's not in the
  allowed state set for its `workflow_id`.

**Warnings** (informational; likely problems, not necessarily errors):

- **Stale claims past TTL.** Equivalent to `fid claim list --stale`;
  surfaced here for the periodic audit context. Recoverable via
  `fid claim --steal` or `fid claim unclaim --force`.
- **Issues in a `started`-category state with no live claim.** With
  a configurable grace window (default 5 minutes after last
  status transition) to allow legitimate brief gaps.
- **AC verifications referencing superseded completions.** Indicates
  stale evidence pointing at abandoned attempts.
- **Task completions with no corresponding `agent_run`.**
  Telemetry gap; doesn't break workflow but means cost/token stats
  will under-count.

**Info** (housekeeping; reportable but never an error):

- Counts of each entity type, claim activity in the last N days,
  recent activity log volume — useful for an `epic wrap-up` audit
  step.

v0.2+ additions (deferred):

- **Activity log projection gaps.** Every journal event that declares
  activity effects should have matching `activity_log` rows. Rebuild
  from the journal and diff against the live DB projection.
- **Journal consistency.** Every DB row should be reproducible from
  the journal. Rebuild into a scratch DB and diff against the live
  DB; any discrepancy indicates a write path that bypassed the
  journal (a real bug) or a hand-edit to the DB.

`--warn-as-error` is intended for pre-merge hooks: if a CI run wants
to block merges on stale claims or unclaimed in-progress issues, this
flag promotes warnings to non-zero exit codes.

---

## 10. Output Formats

### 10.1 Format strategy

| Format | Default for | Use case |
|---|---|---|
| JSON | all read commands (`get`, `list`, `query`) | Universal pipe format; `jq`-pipeable; supported by every library and downstream tool. |
| TOON | opt-in via `--format toon` | Token-efficient (~30-60% saving) for uniform-shape tabular reads. Recommended for agent invocations of `list` commands. |
| Markdown | inspection by humans (`fid issue show`, `fid export`) | Human-readable rendering. |
| YAML | secondary | Optional. Useful for config / template authoring. |

Every read command supports `--format json|toon|markdown|yaml` to
override the default.

### 10.2 Why JSON default, TOON opt-in

JSON is the universal pipe format. Any human or external tool that
shells `fid issue list | jq ...` works without surprises. Tutorials,
shell scripts, CI workflows, and ad-hoc debugging all assume JSON
unless told otherwise — defaulting to it means the discoverable path
is the working path.

TOON is genuinely token-efficient for tabular agent reads — that's its
design goal — but loses its edge on heterogeneous / deeply-nested data,
has near-zero ecosystem tooling, and is a young format whose spec may
still shift. Capturing its savings should be a deliberate choice, not
a default that surprises non-agent callers.

Practical pattern for agent integration: built-in agent role templates
(shipping with v0.3 templates) include `--format toon` as a default
arg for `list`-style commands. Roles inherit the optimisation without
each prompt needing to spell it out.

### 10.3 TOON support boundary

- Treat TOON as a **secondary** wire format. JSON is canonical for any
  inter-tool contract.
- Don't pin behaviour or downstream consumers to TOON-specific shapes
  that differ from the JSON shape.
- Implement the TOON marshaller in-tree (it's small) rather than taking
  a third-party dependency that may break with a spec revision.
- If TOON usage measurements over time show overwhelming opt-in for
  `list` commands, flipping the default in a later version is a small
  change. The reverse — flipping a TOON default back to JSON later —
  breaks shell pipelines written against TOON output, so the
  conservative direction is JSON-first now.

---

## 11. Export to Markdown (and future Sync)

### 11.1 Scope: v0.4 export, v0.5+ optional import

v0.4 ships **one-way export only**: the DB is upstream of the markdown
tree, and `fid export` renders `./docs/` as a derived view. (The
journal is upstream of the DB per ADR-0001; for export purposes the
DB is the immediate source.) Exported files carry a banner so humans don't mistake them for
the source of truth:

```markdown
<!-- generated by fid; edits will be lost on next export -->
```

The exported tree is suitable for git diffs, PR review, Confluence-
style stakeholder consumption, and human reading. It is typically
gitignored or committed under an `_exports/` path treated as
generated content.

Bidirectional sync — importing edited markdown back into the DB, or
running with the markdown tree as the source of truth — is **deferred
to v0.5+** pending evidence of a real use case. The articulated need
("hand-write a Jira-like ticket and import it without typing CLI
calls") is currently better served by an LLM agent with a `fid` skill
that parses loose markdown into the right CLI invocations; that path
doesn't require fid itself to ship a strict parser.

What is preserved regardless: the **frontmatter contract** (§11.3).
Every exported file is self-describing in YAML frontmatter, which keeps
the door open for round-trip without committing to it.

Config persisted in `.fid/config.toml`:

```toml
[export]
markdown_root = "./docs/"
layout = "whet"               # 'whet' | 'flat' | 'custom'
banner = true                 # emit the "generated by fid" banner
```

A future `[sync]` section (v0.5+) would add `canonical = "db" | "md"`
and `auto_export = true | false`. Not present in v0.4 config.

### 11.2 Default `whet` layout

Mirrors the existing `whet-private-docs/` structure. Directory and file
names use the entity's ID (since hash IDs replace the sequential
`phase-N` / `task-M` numbering of the older whet layout):

```
docs/
├── current/
│   ├── phase_brief.md            (rendered from active epic + epic_metadata)
│   ├── decisions_active.md       (rendered from decisions WHERE status='active')
│   ├── decisions_principles.md   (rendered from carry_forward WHERE promoted_principle)
│   ├── progress_current.md       (computed view over status fields)
│   ├── tl_escalations_current.md (rendered from escalations WHERE status='open')
│   └── tl_handoff.md             (optional, per-session)
├── epics/
│   └── <epic-id>/                (e.g. WHET-a1b2c3/)
│       ├── tasks/
│       │   ├── <task-id>.md      (issues row + spec_md_path + ACs as `- [x]`)
│       │   └── <task-id>-done.md (latest task_completions row)
│       ├── qa/
│       │   ├── qa-plan.md
│       │   ├── review-<review-id>.md
│       │   ├── final-review.md
│       │   └── human-checklist.md
│       ├── decisions.md
│       ├── triage-log.md
│       └── progress.md
└── reviews/
    └── YYYY-MM-DD-<epic-id>-retrospective.md
```

A `slug` field on issues (added if/when human-friendly directory names
become a priority) could later substitute for the raw ID in paths.
v0.1 uses the ID directly.

### 11.3 Frontmatter contract

Every exported file has YAML frontmatter carrying the structured fields:

```markdown
---
fid:
  id: WHET-a1b2c3
  kind: task
  status: dev_complete
  workflow: dev-flow
  epic: WHET-4f7c8d
  priority: 1
  sequence: 10
  labels: [bug, auth]
  est_tokens: 8000
  actual_tokens: 9120
  acs:
    - id: ac-1
      text: "..."
      verified: true
      verifier: developer
  links:
    - blocks: WHET-7c3e9f
    - implements: PRD-AUTH-001
---

# WHET-a1b2c3: Add OIDC support

## Objective
...
```

In v0.4 (export only) the frontmatter is informational — it makes
exported files self-describing and machine-readable by future tooling
without requiring fid itself to consume them back. If/when v0.5+ ships
import, the frontmatter contract is exactly what that import path
reads.

Round-tripping prose perfectly is impossible; the frontmatter contract
guarantees only that *structured fields* are reliably represented.
User-facing docs are explicit about this.

---

## 12. Mapping `whet`'s Existing Artifacts

For the bootstrap-of-`fid`-on-`whet` case:

| `whet` artifact | `fid` entity | Notes |
|---|---|---|
| Phase | `issues` (kind='epic') + `epic_metadata` (scheduling='sequential') | Phase brief becomes `epic_metadata.brief_md_path`. |
| Task spec (`task-M.md`) | `issues` (kind='task') + `acceptance_criteria` rows + `spec_md_path` | ACs extracted into rows; rest stays markdown. |
| Task completion report (`task-M-done.md`) | `task_completions` row + `notes_md_path`; `## Concerns` section parsed into `concerns` rows | Append-only; each concern is a structured row with its own lifecycle. |
| QA Reviewer mid-phase (`review-M.md`) | `qa_reviews` (scope='epic_window') + `qa_findings` rows | Whet mid-phase artifact; findings become queryable rows. |
| QA Lead phase-end (`final-review.md`) | `qa_reviews` (scope='epic_full') + `qa_findings` rows | Whet phase-end artifact; literal CLI evidence stored per finding via `evidence_excerpt`. |
| QA plan (`qa-plan.md`) | `qa_plans` table | First-class QA verification design for the epic. |
| `decisions_active.md` | `decisions` rows (status='active') + `carry_forward` rows | Replaces the file; rendered on demand. |
| `decisions_principles.md` | `carry_forward` rows with `promoted_principle = 1` | Promotion is a column flip. |
| `tl_escalations_current.md` | `escalations` rows (status='open') | Replaces the file; rendered on demand. |
| `triage-log.md` | `triage_events` rows | Per-epic, queryable across epics. |
| `progress_current.md` | Computed view over status fields | No primary entity; rendered on read. |
| Phase decisions (`phases/phase-N/decisions.md`) | `decisions` rows linked to epic + `retrospectives` row for "Workflow notes" | Workflow notes section becomes a `retrospectives` row, type='lightweight'. |
| Deep retro (`docs/reviews/...`) | `retrospectives` row, type='deep' | |
| `tl_handoff.md` | Out of scope for v0.1; remain markdown | Optional handoff between TL sessions. |

Note: `AGENTS.md`, `CLAUDE.md`, golden docs (`01_PRD.md`, `02_TRD.md`,
etc.), and the workflow doc itself (`09_AI_WORKFLOW.md`) stay as
markdown. They are not project state — they are project canon.

---

## 13. End-to-End Task Flow

For a sequential epic (bootstrap pattern):

1. **TL planning.** TL creates the epic-kind issue, then for each task:
   `fid issue create --kind task ...`, then adds ACs and `links` rows
   for requirement implementation and dependencies. Writes
   `spec_md_path` pointing to the task spec markdown.
2. **Workflow Runner spawn.** Runner calls `fid ready --epic <id>` to get
   the next task. Spawns Developer with the task ID. Developer acquires
   a claim via `fid issue claim <id> --by developer --ttl 30m`, then
   calls `fid issue transition <id> --to in_progress` (which requires
   a live claim by the caller). Opens an `agent_runs` row.
3. **Developer execution.** Developer reads the task via
   `fid issue get <id>`, optionally drilling to specific sections. As
   ACs are met, calls `fid issue ac verify` per AC. On completion, calls
   `fid issue complete` (composite — closes `agent_runs`, creates
   `task_completions`, sets `concerns_raised` if applicable, transitions
   to `dev_complete` with rule validation).
4. **Workflow Runner inspection.** Runner reads task via
   `fid issue get <id> --json`, inspects `concerns_raised`.
   - If true → spawns TL-Lite in **Triage mode** → TL-Lite calls
     `fid triage record ... --verdict RESUME|BLOCK`. (TL-Lite also
     has a separate **Question mode**, invoked by the human directly
     for a single tactical clarification — not part of the runner
     loop, but uses the same role and escalation rules.)
   - If false → continues.
5. **Mid-epic QA review (every 3 tasks, v0.3+).** Runner calls
   `fid qa-review create --scope epic_window --tasks t1,t2,t3`, spawns
   QA Reviewer who acquires the review via `fid qa-review claim <id>
   --by qa_reviewer`. QA Reviewer adds findings; blocking findings may
   escalate via `fid escalation create --from-finding ...`. With
   claims per QA review, two QA Reviewer agents can work independent
   review pools in parallel without conflict.
6. **Final epic QA Lead review (v0.3+).** After all tasks complete, Runner calls
   `fid qa-review create --scope epic_full`. QA Lead executes the
   deterministic checklist, verifies ACs with `verifier_role=qa_lead`,
   generates the human checklist as a queryable subset
   (`fid view run human-checklist-WHET-4f7c8d`).
7. **TL wrap-up.** TL queries the epic to pull all decisions,
   escalations, triage events, reviews, findings, and AC verification
   state. Updates `decisions` rows, closes resolved escalations, writes
   the `retrospectives` (lightweight) row, transitions the epic to
   `wrapped`. Human transitions to `closed` on approval.

For a prioritised epic or unattached task (steady-state):

- Same lifecycle, minus the batched QA. QA review (if any) is
  scoped `single_task`. Optional retrospective.
- The `ready` queue uses priority + dependencies, not sequence order.
- Single-task tickets may use a lighter workflow (e.g., `quick-flow`)
  that skips formal `qa_review` for trivial bugs.

---

## 14. Lightweight Linear/Jira borrows — included

The following are integrated as first-class because the cost is small
and the value compounds:

- **Labels / tags** — multi-value, free-form, project-scoped. Cross-cutting
  filter. (`labels`, `issue_labels` tables.)
- **Estimates** — `est_tokens` on issues. Surfaces oversized tasks; feeds
  retrospective signals.
- **Status categories** — `unstarted` / `started` / `completed` /
  `cancelled` on `workflow_states`. Queries target categories so they
  stay portable when users rename states.
- **Inbox state** — `in_inbox` boolean on issues. Items land in the
  inbox before being assigned to an epic / cycle (Linear-style intake).
  Distinct from `triage_events`, which records TL-Lite adjudications
  of mid-task concerns — same word in real-world usage, different
  concepts in fid.
- **Cycles / sprints** (optional) — opt-in per project. `cycles` table;
  `cycle_id` on issues.
- **Saved views / filters** — `views` table; reusable named queries.
- **Activity log** — projected timeline of meaningful state
  transitions, AC verifications, comment additions, escalations. Derived from
  journal events for querying and display.
- **ID prefix per project** — set at `fid init`; ID generator respects it.
- **Comments** — generic authored notes on issues, QA reviews,
  decisions, escalations, or requirements. Use for context outside
  built-in structured entities; not for workflow-triggering or
  reportable facts.
- **Due dates** — optional; nullable on issues.

### 14.1 Lightweight features explicitly **not** included

- SLA tracking
- Built-in burn-up / burn-down charts (insights queries cover this;
  charts belong in the dashboard package)
- Customer requests
- Multiple assignees / subscribers / watchers
- Initiatives above epics (extra hierarchy level)
- User-defined custom fields
- OKRs / goal tracking
- Time tracking in hours (token telemetry replaces this)

---

## 15. Open and Resolved Questions

### 15.1 Open questions

1. **Should `workflow_runner` become a deterministic harness rather
   than an LLM role?** (Opened 2026-05-23 during spec grilling;
   deferred.) The role's decisions are state-machine-mechanical: pick
   next task, spawn worker, inspect result, dispatch. An LLM isn't
   required. Deferring because the v0.1 prompt-based approach is
   faster to build and iterate, and the harness option stays open
   without schema impact (drop the role from the enum and add
   `fid run --epic` when ready). Revisit when: routine workflow_runner
   cost / latency becomes meaningful, agent prompt drift causes a
   real failure, or v0.3 MCP server work forces the spawn mechanism
   to be explicit.

> Historical note: open questions are now tracked in the formal spec suite —
> product/scope in `project-docs/specs/01_PRD.md` §7, technical in
> `project-docs/specs/02_TRD.md` §11. New questions go there, not here. This
> section is retained as the origin record.

### 15.2 Resolved questions

1. **Where do QA plans live? — distinct, three-tier model (resolved
   during v0.2 spec).** Validation contracts (assertion-level spec),
   QA plans (verification design), and QA reviews (execution) are
   three first-class entities, linked but separate. All three layers
   are optional — a project may have a QA plan without a contract, or
   task ACs without any assertion layer. When present, fid enforces
   coverage rules and answers traceability queries. See §6.15 for the
   schema and §9.5 `fid report contract-coverage` /
   `fid report assertion-trace`. Sequencing: `qa_plans` +
   `qa_plan_steps` in v0.3 alongside `qa_reviews`;
   `validation_contracts` + `validation_assertions` +
   `validation_assertion_validators` + `val_assertion_outcomes` +
   linkage tables in v0.4.

2. **Concerns heading enforcement on completion notes — template-
   driven in v0.3+.** Workflow Runner continues to perform the check in
   v0.1/v0.2. When templates ship in v0.3, the `completion_report`
   template declares required sections (including the concerns
   section) and `fid dev finish` validates the file against the
   template structurally. Removes the duplicated check from Workflow
   Runner. The general template-driven validation pattern then applies
   to other templated artifacts (task specs, QA reviews, decision
   records) the same way.

3. **Implementation language / runtime — Go.** The agentic build
   pipeline proven on `whet` (also Go) transfers directly; reference
   implementation Beads is also Go, providing a baseline of
   architectural patterns to learn from. The correctness wins Rust
   offers (sum-type modelling of polymorphic refs, sqlx
   compile-time-checked queries) are real but not sufficient to offset
   the iteration tax and the cost of simultaneously learning Rust
   while designing fid. Library and project-layout sub-decisions
   captured at §6.16.

4. **Validation timing for `kind`-specific rules — hybrid.**
   Application-layer validation is primary; SQLite triggers are a
   backstop for writes that bypass the CLI (future MCP tools, web
   dashboard, `fid sql` escape hatch). See §6.14 for the rule
   distribution and §9.8.1 for `fid maintenance check` scope.

5. **Frontmatter dialect — YAML.** Markdown frontmatter is YAML by
   convention across every static-site generator and editor plugin.
   The Codex-config TOML-match isn't compelling because frontmatter
   solves a different problem (document metadata, YAML's home turf)
   than config (TOML's home turf). YAML footguns (Norway problem,
   `yes`/`no` parsing) don't apply because `fid export` writes from a
   typed schema, not from arbitrary hand input.

### 15.3 Closed by upstream decisions

These were live in earlier drafts but became moot once larger decisions
landed. Listed for traceability.

- **DSL for `fid query`** — superseded by §9.5 named reports (v0.2)
  with DSL deferred to v0.3 conditional on uncovered queries appearing
  in practice. See OQ2 in earlier drafts.

- **`workflow_transition_requires` rule library shape** — closed by
  the Decision-1 outcome that defers the configurable workflow engine
  to v0.2. Rules live as named Go functions in v0.1 (hard-coded);
  the v0.2 migration moves them to data while keeping the named
  references. A DSL for rules is not on the roadmap. See OQ3 in earlier
  drafts.

- **Migration tool from existing whet markdown** — closed by
  Decision-9 (export-only in v0.4, import deferred to v0.5+
  conditional). The agent-with-skill alternative covers the use case
  without fid-side parser engineering. See OQ4 in earlier drafts.

---

## 16. Implementation Roadmap

Build order, by value-per-effort:

### v0.1 — Core issue layer

- Schema for `issues`, `epic_metadata`, `acceptance_criteria`,
  `ac_verifications`, `task_completions`, `concerns`, `links`,
  `claims`, `agent_runs`, `labels`, `issue_labels`, `activity_log`,
  `roles`.
- **Stub tables created in v0.1 (schema only, no CLI commands):**
  `decisions`, `escalations`, `triage_events`, `requirements`. These
  entities ship CLI surface in v0.2, but their tables must exist in
  v0.1 so that FK columns on `concerns` (`escalated_to_esc_id`,
  `triaged_in_event_id`) and `links` (requirement endpoints) are
  valid from day one. The CLI returns a clear "this entity's commands
  ship in v0.2" message if a user tries to write to them directly.
- `links` v0.1 scope: issue ↔ issue and issue → requirement. Other
  endpoints (`decision`, `escalation`) activate in v0.2 when their
  CLI surface ships.
- `roles` table seeded at `fid init` from the built-in `default`
  profile (the seven roles in §6.13). Permissions matrix lives in app
  code in v0.1 (see ADR-0002); unknown roles via `--by` are
  hard-rejected with a clear error.
- Hard-coded default workflows for `epic`, `task`, `subtask`
  (implemented in application code, not in tables). `issues.status`
  is CHECK-constrained per `workflow_id` enum.
- Hard-coded transition rule functions: `all_acs_verified`
  (parameterised by required verifier roles), `has_completion_report`,
  `concerns_addressed`, `all_child_tasks_done_or_cancelled`.
- Claim primitives: atomic acquire with TTL, force-release for
  stale takeover; transitions into `started`-category require live
  claim by caller.
- CLI: `init`, `issue create/get/list/transition/complete`, `issue ac
  add/verify`, `issue claim/unclaim`, `claim list`, `ready`,
  `link add/list`, `label add/list`, basic `activity` log read.
- Output formats: markdown, JSON, TOON.
- Markdown section parsing (`--section` flag).
- Status transition enforcement with hard-coded rules.
- ID prefix configuration; hash-suffix ID generation.
- Cost telemetry: `agent_runs` captures snapshot pricing + `cost_usd`
  at run end; `.fid/pricing.toml` ships with defaults for current
  Anthropic / OpenAI / Google models. `fid stats cost` available from
  v0.1.
- Journal: the canonical artifact per ADR-0001. Every CLI write
  appends and flushes `.fid/journal/YYYY-MM-DD.jsonl` as the commit
  point, then applies the event to `.fid/fid.sqlite` as a synchronous
  projection. Journal committed to git; DB gitignored.
- Journal events: one logical CLI write = one semantic event envelope
  with generic ordered row effects. Replay applies effects verbatim
  after event-schema upgrade; it does not re-run command handlers.
- `fid rebuild`: replays the journal from scratch into a fresh DB.
  Foundational primitive — `fid` itself invokes it on startup when
  the DB is missing or behind the journal tail. Events are versioned
  (`event_schema: <n>`) so the replay engine can upgrade old events
  into the current schema.

This alone replaces ~70% of daily interaction with whet's task
artifacts.

### v0.2 — Configurable workflow engine, decisions, escalations, triage, requirements

- Schema migration: introduce `workflows`, `workflow_states`,
  `workflow_transitions`, `workflow_transition_requires` tables.
  Seed them with the v0.1 defaults verbatim. Migrate
  `issues.workflow_id` to a real FK and `issues.status` to a composite
  FK on `(workflow_id, state)`.
- CLI: `fid workflow define`, `fid workflow assign`, `fid workflow list`.
- Schema migration: introduce `role_permissions` table (see ADR-0002).
  Seed from the v0.1 app-code matrix verbatim. CLI gains `fid role
  create`, `fid policy set`, `fid policy show`. After migration,
  `--by <unknown-role>` becomes a config issue rather than a hard
  rejection.
- CLI surface for `decisions`, `escalations`, `triage_events`,
  `requirements` (tables already exist as v0.1 stubs).
- `links` endpoint scope expands to `decision` and `escalation` for
  contextual traceability. Escalation resolution remains
  `escalations.resolution_dec`, not a generic link.
- Drift detection for `requirements.last_synced_hash`.
- `comments` on issues, decisions, escalations, and requirements.
  The `qa_review` comment scope activates in v0.3 when QA review
  tables ship.
- Named reports (`fid report epic-status`, `open-escalations`,
  `qa-blockers`, `ac-matrix`, `epic-wrapup-pack`, `task-hygiene`,
  `stale-claims`). Opinionated wrappers over hand-coded SQL.
- Context packs: `fid issue get <id> --context-pack <role>` returns
  role-tuned bundles (see §9.2.1). Built-in packs for developer,
  qa_reviewer, qa_lead, tl_lite, tl.
- `fid merge --journal <path>` (replay foreign journal with conflict
  detection on entity / claim / status). Mostly redundant with
  `git merge` on textual JSONL, but useful when integrating a
  journal from outside git (e.g., a recovered worktree). `fid
  rebuild` itself ships in v0.1 (see above).

This is where cross-epic analytical value compounds.

### v0.3 — QA reviews and findings, QA plans, templates, TUI board, MCP server

- Schema and CLI for `qa_reviews`, `qa_findings`,
  `qa_review_task_window`. Polymorphic scope handling. Literal
  evidence excerpts on findings (`qa_findings.evidence_excerpt`).
- Claims expand to `qa_reviews`; `fid qa-review claim` / unclaim
  commands become available.
- Schema and CLI for `qa_plans`, `qa_plan_steps` (see §6.15).
  `qa_reviews` gains nullable `plan_id` FK. QA plan finalisation
  is a no-op for coverage rules until v0.4 ships contracts.
- `templates` table and CLI (`template create/list/show`); built-in
  template kinds. Template versioning.
- Template-driven artifact validation: templates declare required
  sections (e.g., `completion_report` requires `## Concerns`); the
  CLI validates files against the template at write time and rejects
  malformed artifacts. Removes structural checks from Workflow Runner.
- Schema migration: add `created_from_template_{kind,name,version}`
  columns to `issues`, `task_completions`, `qa_reviews`, `decisions`
  (see §8.1). Populated on creation when `--from-template` is used.
- TUI kanban board (`fid board`), grouped by status category, filtered
  by epic / labels / cycle.
- **Conditional:** `fid query` DSL, only if v0.2 named reports leave
  meaningful gaps. Grammar specced based on uncovered queries.
- **MCP server (`fid mcp start`).** Standard
  resources / prompts / tools split:
  - **Resources** — current epic brief, active task spec, ready queue,
    active escalations, open concerns for the active issue, AC matrix,
    QA plan.
  - **Prompts** — start developer task, run QA review, TL-Lite triage,
    epic wrap-up. Source of truth is the role-namespace commands +
    permissions matrix.
  - **Tools** — narrow set: claim, transition, ac verify, gate record,
    finish, concern raise, finding add, triage record. Lean by
    default; concern address comes via the template-driven completion
    parser, not an explicit tool. The full CLI surface is *not*
    exposed via MCP, to keep agent context small (Task Master's
    lesson).

Enables full TL wrap-up workflow; replaces individual `review-M.md`
file reading.

### v0.4 — Retrospectives, carry-forward, validation contracts, telemetry, export

- Schema and CLI for `retrospectives`, `carry_forward`.
- Schema and CLI for `validation_contracts`, `validation_assertions`,
  `validation_assertion_validators`, `qa_plan_step_assertions`,
  `val_assertion_outcomes` (see §6.15).
  Adds nullable `acceptance_criteria.val_assertion_id` FK. Coverage
  rule enforcement at `fid qa-plan finalise` time.
- Computed retro signals (token-budget overruns, finding density, AC
  verification gap, contract pass-rate, carry-forward count).
- Named reports: `contract-coverage`, `assertion-trace`,
  `contract-outcomes` (§6.15.5).
- `fid export` with the `whet` layout. One-way; markdown is
  downstream of the DB. Frontmatter contract per §11.3. Exported
  files carry a "generated by fid; edits will be lost" banner.
- `views` (saved queries).

The §9 deep-retro trigger conditions become computed; validation
contract outcomes become the primary wrap-up signal.

### v0.5 — Polish (tentative)

- Web dashboard package (`fid-dashboard`); reads the same SQLite DB.
- `cycles` opt-in support.
- `due_date` UI in board.
- File-scope claims / `expected_files_touched` enforcement (see §6.12).
- **Conditional:** `fid import` + `fid sync` + MD-canonical mode, if
  use-case evidence accumulated during v0.4 export usage justifies it.
  Alternative path: an LLM agent with a `fid` skill that translates
  loose human-authored markdown into CLI calls (no fid changes needed).
- **Conditional:** `fid init --import` from arbitrary whet markdown
  tree, dependent on the above.

---

## 17. Out of Scope (v1)

- Multi-project queries across separate databases.
- Real-time multi-worktree coordination (e.g., locking across
  worktrees, distributed claims). Per §3.1, multi-worktree writes
  produce textual JSONL conflicts at merge time, which is well-
  defined but not automatic. Automatic intent reconciliation is
  out of scope for v1.
- Real-time multi-agent write conflict resolution beyond
  last-write-wins.
- Remote sync, auth, or hosting.
- Time tracking beyond agent_runs token telemetry.
- Notifications, webhooks, integration with external trackers (Jira /
  Linear / GitHub Issues).
- Custom field definitions per project.
- Multi-tenant support.

---

*Draft v0.2 — 2026-05-09. Captures naming decision, entity reshape
(epic-as-issue), templates, configurable workflows, sync/export, and
lightweight Linear/Jira borrows. Open questions in §15 to be resolved
before v0.1 implementation kicks off.*
