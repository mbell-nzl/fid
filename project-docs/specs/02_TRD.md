# TECHNICAL REQUIREMENTS DOCUMENT - fid

> This document records the technical shape of the project: stack choices,
> architecture, data model, integration boundaries, and non-functional
> engineering requirements. Fill stack-specific sections after the project is
> forked.

## 1. Stack Decisions

| Area | Decision | Rationale | Status |
|---|---|---|---|
| Language/runtime | Go (latest stable) | Agentic build pipeline proven on whet transfers directly. Single static binary distribution. Agent error-recovery loop more efficient than Rust during design iteration. See spec §6.16. | Accepted |
| CLI framework | `spf13/cobra` + `spf13/pflag` | Mature nested command tree (`fid issue ac verify`, etc.). POSIX-style flags. Agent-friendly patterns. | Accepted |
| Storage | SQLite via `modernc.org/sqlite` (pure Go) | No CGo cross-compile complexity. Speed adequate for single-project workload. Fallback to `mattn/go-sqlite3` only if perf becomes a real bottleneck. | Accepted |
| Migrations | `golang-migrate/migrate` with file-based SQL | Numbered up/down pairs under `migrations/`. `fid maintenance migrate` for explicit control. | Accepted |
| JSON | stdlib `encoding/json` | No external dependency needed. | Accepted |
| TOON | In-tree marshaller | Small, no third-party dependency risk from a young spec. | Accepted |
| Markdown parsing | `yuin/goldmark` | Export rendering, template section parsing, template-driven validation (v0.3). Extensible via custom node parsers. | Accepted |
| Frontmatter | `adrg/frontmatter` or equivalent | YAML frontmatter read/write for v0.4 export. | Proposed |
| TUI (v0.3) | `charmbracelet/bubbletea` + `lipgloss` | De facto Go TUI stack. Mature. | Proposed |
| Testing | stdlib `testing` + `testify` | Table-driven tests with assertion ergonomics. Integration tests use temporary SQLite per test. | Accepted |
| Build/tool runner | `go build` / `go test` / Makefile | Standard Go toolchain. No custom build system. | Accepted |
| Distribution | Single static binary (cross-compiled) | `go build` produces static binaries by default with pure-Go SQLite driver. | Accepted |

Do not start implementation tasks that depend on a stack decision still marked
`<TBD>` or `Proposed` unless the task is explicitly a discovery spike.

## 2. Architecture

fid follows a layered architecture with dependencies pointing inward. The
journal is the canonical write target; the database is a derived projection.

```text
cmd/fid/              CLI command tree (Cobra)
    │                 Thin adapter: parses flags, delegates to domain packages,
    │                 formats output.
    │
    ├──► internal/format/     Output formatting (JSON, TOON, Markdown, YAML)
    │
    ├──► internal/journal/    Journal writer + replayer
    │    │                    Append-only JSONL; event schema versioning;
    │    │                    upgrader chain; replay into SQLite.
    │    │
    │    └──► internal/db/    Schema, migrations, query helpers
    │                         SQLite is the query cache, not the source of truth.
    │
    ├──► internal/issues/     Issue, AC, completion, concern logic
    ├──► internal/claims/     Claim acquire/release/steal
    ├──► internal/workflow/   State machines + transition rules (v0.1 hard-coded)
    ├──► internal/policy/     Role × action permissions matrix
    ├──► internal/qa/         QA reviews, findings (v0.3)
    ├──► internal/validation/ Contracts, assertions, outcomes (v0.4)
    ├──► internal/templates/  Template engine + structural validation (v0.3)
    ├──► internal/export/     Markdown export with frontmatter (v0.4)
    ├──► internal/search/     Lightweight structured search (v0.2+)
    ├──► internal/stats/      Cost, tokens, gates, attempts aggregations
    ├──► internal/report/     Named reports (v0.2+)
    ├──► internal/changes/    Change digest over activity_log (v0.2+)
    ├──► internal/gitctx/     Read-only git context capture for runs (v0.2+)
    ├──► internal/handoff/    Handoff pack composition over packs/workflow/templates (v0.3+)
    └──► internal/evidence/   Read-only evidence-ref resolution check (v0.4+)
```

Rules:

- Dependencies point inward: `cmd/` depends on `internal/`; `internal/`
  packages depend on each other and on `internal/db/` and
  `internal/journal/` but never on `cmd/`.
- `internal/` boundary keeps the API surface intentional — everything is
  private by default; only `cmd/fid/` calls into the domain packages.
- Domain packages own their validation logic. `internal/db/` provides query
  helpers and connection management but does not encode business rules.
- `internal/journal/` owns the write path: every mutation flows through it.
  Domain packages prepare events; the journal appends them; the DB projects
  them.
- External systems are hidden behind package-level interfaces. No domain
  package directly imports a third-party library for persistence or I/O.

## 3. Module / Directory Boundaries

| Path / Module | Owner Layer | Allowed Responsibilities | Forbidden Responsibilities |
|---|---|---|---|
| `cmd/fid/` | Interface adapter | Parse CLI flags, delegate to domain, format output, exit codes | Business logic, direct DB queries, journal writes |
| `internal/db/` | Infrastructure | Schema definitions, migrations, connection management, query helpers, SQLite pragmas | Business rule enforcement, workflow logic, permission checks |
| `internal/journal/` | Infrastructure | JSONL append, flush, replay, event schema versioning, upgrader chain | Command logic, business validation (receives pre-validated events) |
| `internal/issues/` | Domain | Issue CRUD, AC management, completion recording, concern lifecycle | Output formatting, CLI flag parsing |
| `internal/claims/` | Domain | Claim acquire, release, force-steal, TTL enforcement | Workflow transition logic |
| `internal/workflow/` | Domain | State machine definitions, transition validation, rule evaluation | Issue creation, AC verification |
| `internal/policy/` | Domain | Role validation, permission matrix lookups, action authorization | Workflow state logic |
| `internal/format/` | Interface adapter | JSON, TOON, Markdown, YAML serialisation | Business logic, DB access |
| `internal/qa/` | Domain | QA reviews, findings, plans, plan steps (v0.3) | Issue lifecycle, claim management |
| `internal/validation/` | Domain | Contracts, assertions, required validators, outcomes (v0.4) | QA plan management |
| `internal/templates/` | Domain | Template CRUD, versioning, structural validation (v0.3) | Issue lifecycle |
| `internal/export/` | Interface adapter | Markdown tree rendering, frontmatter generation (v0.4) | Business logic |
| `internal/search/` | Domain | Literal, case-insensitive search over structured read-model fields (v0.2) | Full markdown body indexing, custom query DSL |
| `internal/stats/` | Domain | Cost, token, gate, attempt aggregations | Output formatting |
| `internal/report/` | Domain | Named report SQL composition and execution (v0.2) | Output formatting |
| `internal/changes/` | Domain | Change digest aggregation over the `activity_log` read model: anchor resolution, per-entity net field changes, added/removed child facts (v0.2) | Output formatting, journal writes, new persisted tables |
| `internal/gitctx/` | Infrastructure | Read-only capture of git context (branch, base/head commit, dirty, worktree) behind an interface; bounded, timed subprocess calls; graceful no-repo degradation (v0.2) | Any git mutation (checkout, branch, commit, stash), business logic |
| `internal/handoff/` | Domain | Compose a role handoff pack (v0.3): assemble the manifest (required reads, commands, stop conditions, expected writes, live blockers, evidence obligations) from the role's context pack, the configured workflow, role permissions, and live blocker state; render via the handoff template | Output formatting, journal writes, git operations, hard-coding role or workflow policy |
| `internal/evidence/` | Domain | Resolve cited evidence references read-only (v0.4): commit refs via `internal/gitctx`, file/source/screenshot refs via filesystem existence; classify each `resolved` / `missing` / `unverifiable` by type | Command execution, content hashing, journal writes, evidence storage (operates over existing embedded refs) |
| `migrations/` | Infrastructure | Numbered SQL up/down migration pairs | Application logic |
| `testdata/` | Test support | Fixture DBs, golden files, sample journals | Production code |

## 4. Data Model

### 4.1 Entities

The full SQL schema with DDL is in `03_TRD_SCHEMA.md`. This section
captures the entity inventory, key relationships, and design rationale. Column
lists are indicative; refer to the schema doc for exhaustive definitions.

**Core issue layer (v0.1):**

| Entity | Key Fields | Relationships | Notes |
|---|---|---|---|
| `issues` | id, kind, title, status, workflow_id, priority, sequence, spec_md_path, est_tokens, owner_role | parent_id → issues (hierarchy) | Unified table; `kind` discriminator: epic, task, subtask. v0.5 adds optional `cycle_id` when cycles ship. |
| `epic_metadata` | scheduling, brief_md_path, target_start, target_end | issue_id → issues (1:1, epic only) | Sequential or prioritised scheduling. |
| `acceptance_criteria` | ordinal, text | issue_id → issues (task/subtask only) | Verification state in `ac_verifications`. |
| `ac_verifications` | verifier_role, status, evidence_type, evidence_ref, verified_at | ac_id → acceptance_criteria, completion_id → task_completions | Append-only. Multiple per AC per role. |
| `gate_runs` | issue_id, attempt_number, command, exit_code, output_excerpt, ran_at | issue_id → issues, agent_run_id → agent_runs | Append-only gate evidence recorded before or during completion. |
| `task_completions` | attempt_number, status, actual_tokens, commit_hashes, files_touched, gate_results, human_intervention | issue_id → issues (task only) | Append-only. Previous attempts marked superseded. |
| `concerns` | summary, status, body_md_path | issue_id → issues, raised_in_completion_id → task_completions, escalated_to_esc_id → escalations, triaged_in_event_id → triage_events | Status: open → addressed / escalated / wontfix. |
| `agent_runs` | role, model, tokens_in, tokens_out, cost_usd, price_in/out_per_mtok, token_breakdown_json, tool_call_count, provider_run_id, provider_cost_usd | issue_id → issues, completion_id → task_completions | Snapshot pricing at run end plus optional provider reconciliation fields. v0.3 adds optional `review_id` when QA reviews ship. |
| `claims` | ref_type, ref_id, claimed_by, claimed_at, ttl_seconds | Polymorphic ref (v0.1: issue; v0.3: qa_review) | At most one live claim per entity. PK enforces exclusivity. |
| `roles` | name, description, is_human, active | Seeded from default profile | Existence is validity rule for all `_role` columns. |
| `labels` / `issue_labels` | name, colour; issue_id, label | issue_labels: issue_id → issues, label → labels | Free-form tags. |
| `links` | from_type, from_id, relation, to_type, to_id | Polymorphic endpoints | v0.1 scope: issue ↔ issue, issue → requirement. Auto-inverse. |
| `activity_log` | journal_event_id, ref_type, ref_id, event, details, actor | Polymorphic ref, journal_event_id → journal | Projected from journal events. Read-model only. |
| `projection_meta` | journal_file, journal_offset, last_event_id, last_event_hash, projection_schema_version, replay_status | Single-row DB projection checkpoint | Lets reads detect missing, stale, or inconsistent SQLite state before serving results. |

**Stub tables (v0.1 schema, v0.2 CLI):**

| Entity | Key Fields | Notes |
|---|---|---|
| `decisions` | summary, body_md_path, area, status | FK target for escalations.resolution_dec. |
| `escalations` | status, trigger_class, raised_by_role | FK target for concerns and triage_events. |
| `triage_events` | verdict (RESUME/BLOCK), source, source_ref_id, reasoning_md | FK target for concerns.triaged_in_event_id. |
| `requirements` | type (PRD/TRD), summary, source_doc, source_section, last_synced_hash | FK target for links (issue → requirement). |

**Workflow engine (v0.2):**

| Entity | Key Fields | Notes |
|---|---|---|
| `workflows` | id (stable slug), name, applies_to | v0.1: hard-coded in app code as enum strings. |
| `workflow_states` | workflow_id, state, category, ordinal | Category enables portable queries. |
| `workflow_transitions` | workflow_id, from_state, to_state | Explicit allowed edges. |
| `workflow_transition_requires` | workflow_id, from/to_state, rule, params_json | E.g. `all_acs_verified` with `{"roles":["developer"]}`. |
| `role_permissions` | role_name, action, conditions | v0.1: in app-code constant. v0.2: migrated to table. |

**Governance support (v0.2):**

| Entity | Key Fields | Notes |
|---|---|---|
| `comments` | ref_type, ref_id, author_role, body_md | v0.2 scope: issue, decision, escalation, requirement. v0.3 adds qa_review when QA reviews ship. |

**QA layer (v0.3):**

| Entity | Key Fields | Notes |
|---|---|---|
| `qa_reviews` | scope (epic_window/single_task/epic_full), reviewer_role, body_md_path | Polymorphic scope via epic_id / task_id columns. |
| `qa_review_task_window` | review_id, task_id | Many-to-many for epic_window scope. |
| `qa_findings` | severity (blocking/improvement), status, finding_md, evidence_excerpt | May escalate via escalated_to_esc. |
| `qa_plans` | epic_id, contract_id (nullable), status | Verification design for an epic. |
| `qa_plan_steps` | plan_id, ordinal, display_id (QA-NN), evidence_method, validator_role | Individual verification steps. |
| `templates` | kind, name, version, body_md, active | Monotonic versioning per (kind, name). |

**Validation layer (v0.4):**

| Entity | Key Fields | Notes |
|---|---|---|
| `validation_contracts` | epic_id, status (draft/active/superseded) | Assertion-level spec for an epic. |
| `validation_assertions` | contract_id, ordinal, display_id (VAL-NNN), statement | May reference carry_forward_source. |
| `validation_assertion_validators` | assertion_id, validator_role, evidence_method | Required verifier obligations. |
| `qa_plan_step_assertions` | step_id, assertion_id | Many-to-many: step ↔ assertion. |
| `val_assertion_outcomes` | assertion_id, outcome (pass/fail/carry_forward), evidence_method, evidence_ref, recorded_by_role | Append-only. |

**Remaining (v0.4–v0.5):**

| Entity | Key Fields | Notes |
|---|---|---|
| `carry_forward` | status, scheduled_epic_id, promoted_principle | Multi-epic lifespan. |
| `retrospectives` | epic_id, type (lightweight/deep), recommendation | Links to epic. |
| `cycles` | name, start/end_date, status | Optional; off by default. |
| `views` | name, query_kind, query_text, scope | Saved read-only queries. `query_kind` starts with `sql`; any future DSL kind requires a separate recorded decision. |

**Derived views:**

| View | Purpose | Notes |
|---|---|---|
| `issues_enriched` | LEFT JOINs open concern count and latest triage status onto issues | CLI reads from this view, not the base table. Requires indexes on `concerns(issue_id, status)` and `triage_events(issue_id, created_at)`. |

### 4.2 ID Generation

All generated entity IDs follow the shape `<PREFIX>-<hash>`:

- **Issue IDs** use the project prefix set at `fid init` (e.g. `FID-a1b2c3`).
- **Other entities** use fixed prefixes: AC, ACV, TC, CON, QAR, QAF, DEC,
  ESC, TRG, CF, RET, RUN, TPL, CYC, VIEW, COM, ACT, VC, ASN, QAP, QAS, VAO.
- **Hash** is `n` hex characters from `sha1(uuid4())[:n]`. Fresh UUID per
  call avoids content-derived collisions.
- **Progressive length scaling:** default `n = 6`; when row count pushes
  birthday-paradox probability above ~1%, new IDs use `n = 8`.
- **Collision handling:** on insert, retry with fresh UUID. Post-merge
  collisions detected during journal replay require human rewrite.
- **Natural keys** stay natural: `requirements.id` (PRD-/TRD-),
  `labels.name`, `roles.name`, `workflows.id` (stable slug), and all
  composite-PK junction tables.

### 4.3 Persistence Rules

**Write path (every CLI write):**

1. Acquire the checkout-scoped write lock at `.fid/write.lock`.
2. Catch up or rebuild SQLite if `projection_meta` shows the DB is stale.
3. Validate against the current DB state (permissions, workflow rules,
   kind constraints, claim state).
4. Build a fully resolved event: one semantic envelope
   (`event_type`, `event_schema`, actor metadata, command metadata) with an
   ordered `effects` array of row operations (`insert` / `update` / `delete`
   with table name, PK, and resolved values).
5. Append and flush one JSONL line to
   `.fid/journal/YYYY-MM-DD.jsonl` — this is the commit point.
6. Apply the event's effects to SQLite in a single DB transaction and update
   `projection_meta`.
7. Release `.fid/write.lock`.

If journal append fails, the DB is not modified. If DB projection fails after
journal append, the next CLI invocation catches the DB up before answering.

**Dry-run simulation:** Because validation (step 3) and event construction
(step 4) both precede the journal append (step 5, the commit point), a write
can be simulated by running steps 1–4 and returning the validation result and
the built event without appending or projecting. The v0.2 `--dry-run` flag
(PRD-OPS-015) hooks in at this seam; deferred idempotency keys (PRD-OPS-016)
would add an envelope field plus a write-time key lookup at the same seam.
Keeping prepare/validate cleanly separable from append is therefore a v0.1
write-path requirement even though the flag itself ships in v0.2.

**One logical write = one journal event — no exceptions.** Composite commands
(`fid dev finish`) produce one event with multiple effects applied atomically.
Half-applied composite state is data corruption the replay engine would
faithfully reproduce.

**Event schema versioning:**

- Each event carries `event_schema: <n>` (positive integer). Increments on
  any schema migration that changes event effect shapes.
- Journal files are immutable once written. Old events stay at their original
  schema version.
- The replay engine maintains a chain of pure upgrader functions:
  `upgrade_N_to_N+1(event) → event`. No DB reads, no side effects. Testable
  with golden-file fixtures per version boundary.
- On replay, events older than the current schema version are upgraded
  in-memory before effects are applied.

**Journal event v1 contract:**

Each JSONL line is one canonical JSON object with these top-level fields:

| Field | Required | Semantics |
|---|---|---|
| `event_id` | yes | Stable generated event ID. Used by `activity_log.journal_event_id`, replay diagnostics, and projection metadata. |
| `event_schema` | yes | Positive integer schema version. v0.1 writes `1`. |
| `event_type` | yes | Semantic command/event name, e.g. `issue.created`, `issue.transitioned`, `claim.acquired`. |
| `occurred_at` | yes | UTC ISO 8601 timestamp assigned once before append. |
| `actor` | yes | Object with `role` and optional `actor` string. `role` is the `--by` role for CLI writes. |
| `command` | yes | Sanitized command metadata: `argv`, `cwd`, `fid_version`, and no environment variables or sensitive flag values. |
| `effects` | yes | Ordered array of row-level effects applied atomically during projection. |

Effect objects use a shared shape:

| Field | Required | Semantics |
|---|---|---|
| `op` | yes | `insert`, `update`, `delete`, or `activity`. |
| `table` | yes | Target table name. `activity` effects target `activity_log`. |
| `pk` | yes | Object containing the complete primary key in schema column order. |
| `values` | for `insert`/`update`/`activity` | Resolved column values after the effect. |
| `prior` | for `update`/`delete` when available | Prior values recorded for audit and diagnostics. |
| `preconditions` | targeted | Expected current values that must match during replay before the effect applies. |

Primary keys are always represented as JSON objects, even for single-column
keys: `{"id":"FID-a1b2c3"}` rather than a bare string. Composite natural keys
use the schema column names in primary-key order.

Preconditions are required for mutable-state effects whose correctness depends
on prior state:

- issue status transitions: expected `issues.status` and `workflow_id`;
- claim acquire, renew, release, and force-acquire: expected current claim
  absence or expected holder/expiry metadata;
- update/delete of current-state rows such as `issues`, `labels`,
  `issue_labels`, `links`, `comments`, and `projection_meta`;
- task completion attempt-number allocation: expected latest attempt number
  for the issue;
- link or label removal: expected row existence.

Append-only inserts that do not depend on prior mutable state, such as
`ac_verifications`, `activity_log`, and new `task_completions` rows after
attempt-number allocation has passed, do not need prior-value preconditions.
They still rely on primary-key and uniqueness constraints during projection.

During replay, fid applies events in journal order. If an effect precondition
does not match the projected DB state, replay stops with a logical conflict
error naming `event_id`, `event_type`, table, primary key, expected value,
actual value, and the journal file/line. v0.1 does not automatically reconcile
intent; the human resolves the journal conflict or uses future merge tooling.

`activity_log.journal_event_id` is always the originating `event_id`.
Projection stores `last_event_id` and a hash of the canonical JSON event line
in `projection_meta`. Canonical JSON for hashing uses UTF-8, sorted object
keys, no insignificant whitespace, and preserves array order.

**Example v1 event:**

```json
{
  "event_id": "EVT-20260529-000001",
  "event_schema": 1,
  "event_type": "issue.transitioned",
  "occurred_at": "2026-05-29T03:12:45Z",
  "actor": {"role": "developer", "actor": null},
  "command": {
    "argv": ["fid", "issue", "transition", "FID-a1b2c3", "--to", "in_progress", "--by", "developer"],
    "cwd": ".",
    "fid_version": "0.1.0"
  },
  "effects": [
    {
      "op": "update",
      "table": "issues",
      "pk": {"id": "FID-a1b2c3"},
      "preconditions": {"status": "todo", "workflow_id": "task-default"},
      "prior": {"status": "todo"},
      "values": {"status": "in_progress", "updated_at": "2026-05-29T03:12:45Z"}
    },
    {
      "op": "activity",
      "table": "activity_log",
      "pk": {"id": "ACT-d4e5f6"},
      "values": {
        "journal_event_id": "EVT-20260529-000001",
        "ref_type": "issue",
        "ref_id": "FID-a1b2c3",
        "event": "status_changed",
        "details": {"from": "todo", "to": "in_progress"},
        "actor": "developer",
        "occurred_at": "2026-05-29T03:12:45Z"
      }
    }
  ]
}
```

**Correction model:**

Journal history is append-only. fid never mutates or deletes a committed event
as part of normal operation. Routine mistakes are corrected by appending a new
domain event, not by editing old journal lines.

v0.1 correction paths are domain-specific:

- wrong issue fields: `fid issue update`;
- wrong label assignment: `fid issue label add/remove`;
- wrong link: `fid link add/remove`;
- wrong claim state: `fid issue claim --renew`, `fid issue claim --force`, or
  `fid issue unclaim`;
- wrong AC verification: append a later `ac_verifications` row with the
  corrected status/evidence;
- wrong status: use `fid issue transition` only if the workflow allows the
  transition from the current state.

`fid journal inspect <event-id>` is a v0.1 read command that parses the
matching journal event and reports actor, command metadata, effects, activity
entries, journal file/line, hash, and revert eligibility. It does not change
state.

v0.2 adds `fid event revert <event-id> --by <role> --reason <text>` for
explicitly reversible event types. Revert writes a compensating journal event
with `event_type: "event.reverted"` and metadata:

```json
{
  "reverts_event_id": "EVT-20260529-000001",
  "reason": "Wrong label applied during triage"
}
```

Revert eligibility is intentionally narrow:

| Event/effect type | v0.1 correction | v0.2 `event revert` eligibility |
|---|---|---|
| Issue field update | `fid issue update` | Yes, if current values still match the event's `values`. |
| Issue label add/remove | `fid issue label add/remove` | Yes, if the label relationship is still in the state produced by the event. |
| Link add/remove | `fid link add/remove` | Yes, if both declared and inverse rows are still in the state produced by the event. |
| Claim release/renew | Claim command | Yes, only when current claim state still matches the event's `values`. |
| AC verification | New verification row | No. Append a corrected verification. |
| Gate run | New gate run or completion note | No. Gate evidence remains audit history. |
| Task completion | Superseding attempt or concern/comment | No. Completion history remains append-only. |
| Status transition | `fid issue transition` when allowed | No generic v0.2 revert. Status changes require workflow validation. |
| Composite developer command | Domain-specific correction | No generic v0.2 revert unless all effects are individually reversible and the command has an explicit reverser. |

If an event is not eligible, `fid event revert` exits non-zero and names the
preferred domain correction command.

**Concurrency:**

- Single-worktree: all CLI writes acquire `.fid/write.lock` before freshness
  checks, validation, event construction, journal append/fsync, and SQLite
  projection. This serializes validation against a fresh DB state before any
  journal commit is possible.
- SQLite transactions remain required inside the write-lock boundary for
  projection atomicity. SQLite locking alone is not the write serializer,
  because the journal append is the commit point.
- `.fid/write.lock` is acquired by atomically creating a lock directory that
  contains owner metadata: PID, hostname, started_at, command, and actor role
  when available. If the lock is held, fid waits for the configured short
  timeout, then exits non-zero with an error naming the lock holder and path.
- v0.1 does not silently break locks. A lock may be treated as stale only when
  fid can prove the owner process is dead on the same host; otherwise manual
  human removal is required.
- Multi-worktree: produces JSONL line-merge conflicts in the day's journal
  file, resolvable as standard text merge. Next `fid` invocation rebuilds
  the DB from the merged journal.
- Logically conflicting writes surface as JSONL conflicts on the same lines.
  No automatic intent reconciliation.

**What is committed to git:**

- `.fid/journal/YYYY-MM-DD.jsonl` — the source of truth.
- `.fid/config.toml` — project configuration.
- `.fid/pricing.toml` — model pricing data.

**What is gitignored:**

- `.fid/fid.sqlite` — derived query cache. Deletable and rebuilt.

**Rebuild:**

- `fid rebuild` replays the full journal into a fresh DB. Linear scan;
  sub-second at single-project scale.
- On every CLI invocation, the DB is checked against the journal tail. If
  behind or missing, it is caught up synchronously before the command answers.
- The SQLite read model stores one `projection_meta` row recording the last
  applied journal file, byte offset after the applied line, event ID, event
  line hash, projection schema version, and replay status.
- If `projection_meta` is absent, has the wrong projection schema version, or
  points outside the available journal files, fid rebuilds the DB from
  scratch before serving a read.
- If `projection_meta` matches the journal position and last event
  ID/hash, fid catches up by applying only events after the recorded byte
  offset.
- If the recorded event ID/hash does not match the journal at the recorded
  position, fid stops before answering and reports projection inconsistency.
  It does not silently skip, truncate, or reorder journal entries.

**Malformed journal:**

- A malformed trailing JSONL line is treated as journal corruption. Startup
  and `fid rebuild` stop and report the corrupt tail. No silent skip or
  truncation.

### 4.4 SQLite Configuration

Required pragmas on every connection:

```sql
PRAGMA foreign_keys = ON;
PRAGMA recursive_triggers = ON;
```

### 4.5 Validation Strategy

**Application-layer primary, SQLite triggers as backstop.**

Application layer (all writes):
- Workflow transition requirements with rich error messages.
- Parameterised transition rules (e.g. `all_acs_verified` reads
  `params_json.roles`).
- Permissions matrix (needs `--by <role>` context SQL can't see).
- AC verification freshness.
- Template structural validation (v0.3+).
- Claim acquire semantics (TTL, stealing).

SQLite triggers (backstop for non-CLI write paths — future web dashboard,
`fid sql`, direct DB edits):
- Kind-invariant rules: hierarchy enforcement on `issues.parent_id`,
  sequence constraints, kind guards on `acceptance_criteria`,
  `task_completions`, `epic_metadata`, `qa_reviews`, `triage_events`,
  `decisions`, `carry_forward`.
- Polymorphic ref integrity on `links`, `claims`, `comments`,
  `activity_log`.
- Link scope and relation compatibility.
- Claim scope (supported `ref_type` values per schema version).
- Comment scope (v0.2 base ref_type whitelist; v0.3 adds `qa_review`).
- Kind immutability (BEFORE UPDATE OF kind on `issues`).
- Mutually exclusive scope columns on `qa_reviews` (CHECK constraint).

What does **not** get a trigger: workflow transition rules, permissions,
AC verification freshness, template validation. These need rich context or
runtime state that SQL doesn't have.

### 4.6 Default Workflows (v0.1 Hard-Coded)

v0.1 ships with hard-coded workflow definitions in `internal/workflow/`.
`issues.workflow_id` is an enum string and `issues.status` is constrained to
the legal states for that workflow. v0.2 migrates these definitions verbatim
into the workflow tables before enabling user-defined workflows.

**Task workflow: `task-default`**

| State | Category |
|---|---|
| `todo` | `unstarted` |
| `in_progress` | `started` |
| `dev_complete` | `started` |
| `qa_review` | `started` |
| `qa_block` | `started` |
| `failed_gate` | `started` |
| `qa_pass` | `completed` |
| `done` | `completed` |
| `cancelled` | `cancelled` |

| From | To | Requirements |
|---|---|---|
| `todo` | `in_progress` | - |
| `in_progress` | `dev_complete` | `all_acs_verified{"roles":["developer"]}`, `has_completion_report` |
| `in_progress` | `failed_gate` | - |
| `failed_gate` | `in_progress` | - |
| `dev_complete` | `qa_review` | - |
| `qa_review` | `qa_pass` | `all_acs_verified{"roles":["qa_lead"]}` |
| `qa_review` | `qa_block` | - |
| `qa_block` | `in_progress` | - |
| `qa_pass` | `done` | - |
| any non-`cancelled` state | `cancelled` | - |

**Epic workflow: `epic-default`**

| State | Category |
|---|---|
| `todo` | `unstarted` |
| `in_progress` | `started` |
| `wrapped` | `started` |
| `closed` | `completed` |
| `cancelled` | `cancelled` |

| From | To | Requirements |
|---|---|---|
| `todo` | `in_progress` | - |
| `in_progress` | `wrapped` | `all_child_tasks_done_or_cancelled` |
| `wrapped` | `closed` | - |
| any non-`cancelled` state | `cancelled` | - |

**Subtask workflow: `subtask-default`**

| State | Category |
|---|---|
| `todo` | `unstarted` |
| `in_progress` | `started` |
| `done` | `completed` |
| `cancelled` | `cancelled` |

| From | To | Requirements |
|---|---|---|
| `todo` | `in_progress` | - |
| `in_progress` | `done` | - |
| any non-`cancelled` state | `cancelled` | - |

**Escalation workflow: `escalation-default` (v0.2)**

| State | Category |
|---|---|
| `open` | `started` |
| `resolved` | `completed` |

| From | To | Requirements |
|---|---|---|
| `open` | `resolved` | `linked_to_decision` |

### 4.7 Gate Run Persistence

`fid dev gate` writes append-only `gate_runs` rows. It does not create or
mutate `task_completions` directly, because completions represent finalized
attempts.

In v0.1, the active attempt number for a task is derived as:

```text
MAX(task_completions.attempt_number WHERE issue_id = <task>) + 1
```

with `0` treated as the max when no prior completion exists. `fid dev gate`
may record a gate only when the issue is in a started-category state and the
caller holds the live claim required to advance the task. The gate run may
reference the open `agent_runs.id` when one exists, but that reference is
optional so primitive workflows can still record gate evidence.

At `fid dev finish`, `fid dev block`, or `fid issue complete`, fid collects
all `gate_runs` for the issue and active attempt number, orders them by
`ran_at` then `id`, and snapshots them into `task_completions.gate_results`.
It also stores `gates_passed` and `gates_failed` counts on the completion row.
The source `gate_runs` rows remain as detailed audit evidence.

If a later attempt starts, it receives the next derived attempt number and
does not mutate prior `gate_runs` or `task_completions`.

## 5. Interfaces and Contracts

| Contract | Consumers | Compatibility Policy | Verification |
|---|---|---|---|
| CLI command surface (`fid <noun> <verb> [flags]`) | Agents, humans, shell scripts | Stable within a major version. Flag additions are non-breaking; flag removals require a deprecation cycle. | Integration tests per command. |
| Journal event format (JSONL with semantic envelope + effects) | `fid rebuild`, `fid merge --journal`, future tools | Forward-compatible via `event_schema` versioning and upgrader chain. Old events never rewritten. | Golden-file replay tests per schema version boundary. |
| SQLite schema (tables, views, indexes, triggers) | CLI queries, `fid sql`, future web dashboard | Versioned via numbered migrations. | Migration up/down tests. |
| JSON output (default for all read commands) | Agents, `jq` pipelines, downstream tools | Field additions are non-breaking; field removals require deprecation. | Snapshot tests per command output. |
| TOON output (opt-in via `--format toon`) | Agents optimising for token efficiency | Secondary format; shape mirrors JSON. | Parity tests against JSON output. |
| Search result shape (`fid search`) | Agents, humans | Additive result fields are non-breaking. Matching remains literal/case-insensitive unless a future decision commits to ranking or a query language. | Integration tests for entity filters, empty results, and JSON/TOON parity. |
| Markdown export with YAML frontmatter (v0.4) | Humans, git diffs, stakeholder review | Frontmatter schema documented; field additions non-breaking. | Export round-trip tests (structured fields survive export/re-parse). |
| `.fid/config.toml` | CLI and export | Additive sections; no removal without migration. | Config parsing tests with missing/extra keys. |
| `.fid/pricing.toml` | Cost computation in `agent_runs` | Additive model entries; user-overridable. | Test with missing model, custom model, default model. |

`fid search <term>` is a read-only discovery command, not a query language.
The v0.2 implementation performs literal, case-insensitive matching across
projected structured fields such as entity IDs, issue titles, labels, owner
fields, requirement summaries, decision summaries, escalation summaries, and
markdown path fields where present. It supports filters such as `--kind`,
`--status`, `--label`, `--epic`, `--inbox`, and `--format`. Full markdown body
search remains out of scope for fid search; callers should use shell tools such
as `rg` against markdown files when they need prose-content search.

## 6. Security and Privacy

| ID | Requirement | Verification |
|---|---|---|
| TRD-SEC-001 | The system shall not open any network listener or transmit data externally. All data stays local to the filesystem. | Code review; no `net.Listen` or HTTP client calls. Any future MCP or remote adapter requires a new recorded decision and threat model. |
| TRD-SEC-002 | The `fid sql` escape hatch shall reject any statement other than SELECT, EXPLAIN, PRAGMA query, and WITH (read-only). | Integration test: attempt INSERT/UPDATE/DELETE/DROP via `fid sql`; assert rejection. |
| TRD-SEC-003 | Journal files shall not contain secrets. CLI shall strip environment variables and sensitive flag values before recording command metadata in events. | Code review; test that event `command_metadata` contains no env vars or `--password`-style values. |
| TRD-SEC-004 | SQLite triggers shall prevent structural corruption even from non-CLI write paths (future dashboard, direct SQL). | Trigger tests: attempt invalid inserts via raw SQL; assert RAISE fires. |

## 7. Performance and Reliability

| ID | Requirement | Target | Verification |
|---|---|---|---|
| TRD-PERF-001 | Single read command latency at 10k issues | < 500ms | Benchmark test with seeded DB. |
| TRD-PERF-002 | Full journal rebuild at 10k events | < 5s | Benchmark test with seeded journal. |
| TRD-PERF-003 | Full journal rebuild at 100k events | < 30s | Benchmark test with seeded journal. |
| TRD-PERF-004 | `issues_enriched` view query with concern/triage joins | No table scan; index-only lookups | EXPLAIN QUERY PLAN assertion in test. |
| TRD-PERF-005 | Single write command (validate + journal append + DB project) | < 200ms | Benchmark test. |
| TRD-PERF-006 | Concurrent single-worktree writes from two agents | Both succeed without data loss | Integration test with two goroutines writing concurrently. |
| TRD-PERF-007 | Binary size (single static binary, no CGo) | < 30 MB | CI size check. |

## 8. Observability and Operations

- **Logging:** Structured stderr output. Verbosity controlled by `--verbose`
  / `--quiet` flags. No log files; CLI is invocation-scoped with no daemon.
- **Metrics:** No runtime metrics system. Token/cost/gate telemetry captured
  per agent run in `agent_runs` table. Aggregated via `fid stats` commands.
- **Error reporting:** Structured error output on stderr with exit code.
  JSON-formatted errors when `--format json` is active. Rich error messages
  per PRD-NFR-003.
- **Configuration:**
  - `.fid/config.toml` — project-scoped (export layout, banner toggle,
    future sync settings).
  - `.fid/pricing.toml` — model pricing data; ships with defaults,
    user-overridable.
  - Issue ID prefix set at `fid init`, stored in config.
- **Backup / recovery:** Journal is committed to git; that is the backup.
  `fid rebuild` reconstructs the DB from journal. Malformed journal tail
  is reported for manual repair.
- **Deployment/release process:** `go build` produces the binary. Cross-
  compilation for target platforms via `GOOS`/`GOARCH`. No container or
  package manager distribution in v1.
- **Rollback process:** Git revert on journal files + `fid rebuild`
  reconstructs the prior state.

## 9. Dependency Policy

Runtime dependencies (indicative; finalised during v0.1 implementation):

| Dependency | Purpose | Alternatives Considered | License | Decision Ref |
|---|---|---|---|---|
| `modernc.org/sqlite` | Pure-Go SQLite driver | `mattn/go-sqlite3` (CGo, faster but cross-compile pain) | BSD-3 | spec §6.16 |
| `spf13/cobra` | CLI command tree | `urfave/cli`, stdlib `flag` | Apache-2.0 | spec §6.16 |
| `spf13/pflag` | POSIX-style flags | stdlib `flag` | BSD-3 | spec §6.16 |
| `golang-migrate/migrate` | SQL migrations | `pressly/goose`, hand-rolled | MIT | spec §6.16 |
| `yuin/goldmark` | Markdown parsing | `gomarkdown/markdown` | MIT | spec §6.16 |
| `stretchr/testify` | Test assertions | stdlib only | MIT | spec §6.16 |
| `charmbracelet/bubbletea` | TUI framework (v0.3) | `rivo/tview`, `gdamore/tcell` | MIT | spec §6.16 |
| `charmbracelet/lipgloss` | TUI styling (v0.3) | — | MIT | spec §6.16 |

The v0.2 git-aware run features (PRD-OPS-017–019) invoke a local `git` binary
**read-only** (`rev-parse`, `status --porcelain`) through `internal/gitctx/`.
This is a local CLI tool, not a service/daemon/network dependency, and fid
already assumes a git repo (the journal is committed to git). Invocations are
bounded and timed; if `git` is absent or the cwd is not a worktree, the run is
recorded with no git context rather than failing. fid never mutates git state.

Supply-chain gates to define before release readiness:

- Vulnerability audit: `govulncheck ./...`
- License check: `go-licenses check ./...` or manual review
- Secret scan: `gitleaks detect` in CI
- Update cadence: Dependabot or manual quarterly review

## 10. Technical Requirements Traceability

| TRD ID | Requirement | Version | Verification |
|---|---|---|---|
| TRD-ARCH-001 | All writes shall flow through the journal; the DB shall be a synchronous projection of journal events. | v0.1 | Integration test: write via CLI, delete DB, rebuild, assert identical query results. |
| TRD-ARCH-002 | Journal events shall be self-contained pure data with all generated values resolved at write time. Replay shall not re-run command logic or read DB state. | v0.1 | Golden-file replay tests. |
| TRD-ARCH-003 | Event schema versioning shall use a chain of pure upgrader functions. Each upgrader shall be testable with golden-file fixtures. | v0.1 | Unit test per upgrader with captured input/output events. |
| TRD-ARCH-004 | One logical CLI write shall produce exactly one journal event. Composite commands shall use an `effects` array applied in a single SQLite transaction. | v0.1 | Integration test: `fid dev finish` produces one JSONL line; DB transaction is atomic. |
| TRD-ARCH-005 | `internal/` packages shall not import `cmd/fid/`. `cmd/fid/` shall contain no business logic. | v0.1 | `go vet` import analysis or CI lint rule. |
| TRD-ARCH-006 | SQLite triggers shall enforce kind-invariant structural rules as a backstop for non-CLI write paths. | v0.1 | Integration tests: raw SQL inserts violating kind constraints; assert trigger RAISE. |
| TRD-ARCH-007 | `PRAGMA foreign_keys = ON` and `PRAGMA recursive_triggers = ON` shall be set on every database connection. | v0.1 | Test: open connection, query pragma values, assert ON. |
| TRD-ARCH-008 | The `issues_enriched` view shall use index-only lookups on `concerns(issue_id, status)` and `triage_events(issue_id, created_at)`. | v0.1 | EXPLAIN QUERY PLAN assertion in test. |
| TRD-ARCH-009 | Workflow transition rules shall be implemented as named Go functions in v0.1, migrated to `workflow_transition_requires` table rows in v0.2. | v0.1–v0.2 | v0.1: unit tests per rule function. v0.2: migration test seeds defaults; integration test confirms identical behavior. |
| TRD-ARCH-010 | The permissions matrix shall be a code constant in v0.1, migrated to `role_permissions` table in v0.2. Unknown roles shall be hard-rejected in v0.1. | v0.1–v0.2 | v0.1: test every CLI write action has a matrix entry; test denied combinations raise errors. v0.2: migration test seeds from v0.1 constant. |
| TRD-ARCH-011 | Claim acquisition shall be atomic via `INSERT ... ON CONFLICT DO NOTHING`. Zero affected rows indicates an existing claim. | v0.1 | Integration test with concurrent acquire attempts. |
| TRD-ARCH-012 | The TOON marshaller shall be implemented in-tree with no third-party dependency. | v0.1 | Build with no TOON dependency in go.mod. |
| TRD-ARCH-013 | Agent integration shall remain CLI/skill-first until dogfooding proves a concrete MCP-only client need. Any future MCP adapter shall require a recorded product/security decision, a threat model, and an explicit transport choice before implementation. | v0.1+ | Roadmap review confirms no MCP implementation task is committed before the decision; security review confirms no listener code exists. |
| TRD-ARCH-014 | Markdown export shall include YAML frontmatter with structured fields. The frontmatter schema shall be documented and tested for round-trip stability. | v0.4 | Export → parse frontmatter → assert field values match DB. |
| TRD-ARCH-015 | Template-driven validation shall declare required sections per template kind. `fid dev finish` and similar commands shall validate files against templates and reject malformed artifacts. | v0.3 | Integration test: submit completion report missing `## Concerns` heading; assert rejection. |

## 11. Open Technical Questions

These are implementation/technical questions. Product and scope questions are
tracked in `01_PRD.md` §7.

- [ ] Should the TOON marshaller support the full TOON spec or only the
  subset needed for fid's tabular output? Start minimal and extend on demand.
- [ ] Should `fid rebuild` support partial replay (catch-up from a known
  checkpoint) or always replay from scratch? Full replay is simpler and
  fast enough at single-project scale; checkpoint is a forward-compatible
  optimisation.
- [ ] Should fid add a future MCP adapter after dogfooding? Start from the
  CLI/skill-first assumption; require a concrete MCP-only client need, a
  product/security decision, and a threat model before adding it to a roadmap.

*Last updated: 2026-06-01*
