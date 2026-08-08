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
    ├──► internal/workflow/   Registry loaded from the workflow tables; from v0.2a
    │                         also workflow authoring, assignment, and resolution
    ├──► internal/maintenance/ Read-only structural audit over projection + journal
    ├──► internal/policy/     Role × action permissions matrix
    ├──► internal/dev/        Developer composites over claim/transition/run/gate writes
    ├──► internal/gates/      Gate registry, registered execution, receipts (v0.2b)
    ├──► internal/envelope/   Work envelopes + orchestration primitives (v0.3a)
    ├──► internal/runs/       Agent run recording with snapshot pricing
    ├──► internal/qa/         QA reviews, findings (v0.3)
    ├──► internal/validation/ Contracts, assertions, outcomes (v0.4)
    ├──► internal/templates/  Template engine + structural validation (v0.3)
    ├──► internal/export/     Markdown export with frontmatter (v0.4)
    ├──► internal/search/     Lightweight structured search (v0.2+)
    ├──► internal/stats/      Cost, tokens, gates, attempts aggregations
    ├──► internal/activity/   Projected activity timeline reads over activity_log
    ├──► internal/sqlquery/   Read-only SQL escape hatch over the projected DB
    ├──► internal/report/     Named reports (v0.2+)
    ├──► internal/changes/    Change digest over activity_log (v0.2+)
    ├──► internal/gitctx/     Read-only git context capture for runs (v0.2+)
    ├──► internal/handoff/    Handoff pack composition over packs/workflow/templates (v0.3+)
    └──► internal/evidence/   Evidence records, obligations, freshness + ref checks (v0.2b+)
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

### 2.1 Cross-Product Ownership Boundary (ADR-0004)

fid is the durable operational work-state and evidence backend as well as the
issue tracker. It owns issue/work identity and lifecycle, journal-canonical
events, projections, claims and runs, receipts and evidence, attention,
outcomes, telemetry, rebuild, and portable recovery.

The external workflow kit owns engineering-workflow policy: risk
classification, authority grants, role/reviewer routing, budgets and
obligation selection, deterministic orchestration, and provider projections.
Project files and Git remain canonical for requirements, architecture,
decisions, code, and repository identity; Human actors remain authoritative
for protected decisions and product judgment. Cross-product records refer to
those authorities by versioned identity/digest rather than copying their
truth.

The kit evaluates engineering-workflow policy and compiles its approved result
into the versioned Work Envelope. fid does not recalculate risk, routing,
authority, budget, or obligation-selection policy. It enforces generic
operational integrity and the recorded envelope's preconditions: identity and
expected revision, idempotency, claim/lease ownership, approval authority,
immutable history, non-weakening, trusted evidence and freshness, and legal
lifecycle changes. Shared conformance fixtures guard this public contract.

fid Issues and their relationships are the only durable work graph. Work
Envelopes reference one primary Issue and an explicitly bounded related Issue
set; they do not contain persistent Work Item identities. The harness derives
noncanonical Next actions from durable state. Claims and Runs bind Issue,
Envelope revision, actor role, and optional Obligation. The detailed joint
contract must cover review rework and multi-Issue Envelope revision before
implementation; parallel Claim isolation and interrupted-Run recovery are
resolved by ADR-0006/0007 (ADR-0005).

ADR-0006 resolves parallel claims: fid atomically allows at most one active
mutating Builder Claim per Issue and rejects unproven isolation with
`mutation-conflict`. Parallel mutation requires distinct selected Issues,
separate worktree identities, and declared disjoint scopes established by the
kit/harness. Review Runs may coexist only as repository-identity-bound,
read-only executions whose results stale when the target changes. ADR-0007
separately resolves interrupted-Run recovery; review rework and multi-Issue
Envelope behavior remain open.

ADR-0007 resolves interrupted Runs. A new session receives a new linked Run;
the prior Run becomes `interrupted`. Recovery requires confirmed executor
termination, Claim release/expiry/authorized supersession, and repository-port
reconciliation before mutation. Unknown liveness and unsafe repository state
fail closed as `executor-liveness-unknown` and `recovery-state-conflict`.
Recovery persists facts and references but performs no Git mutation. Budgets,
obligations, concerns, and evidence freshness carry forward. Two joint
scenarios remain open.

fid is the primary trusted local executor for registered gates: callers supply
an immutable registration reference, fid records the start before process
spawn, observes termination, and creates the receipt.
External attested runners are a secondary CI/Assured/release extension;
reported results cannot satisfy trusted obligations. Live execution output is
required, with exact interactive/machine framing still to be specified. No
adapter may bridge contract gaps by reading SQLite, journal files, or `.fid/`
internals.

## 3. Module / Directory Boundaries

| Path / Module | Owner Layer | Allowed Responsibilities | Forbidden Responsibilities |
|---|---|---|---|
| `cmd/fid/` | Interface adapter | Parse CLI flags, delegate to domain, format output, exit codes | Business logic, direct DB queries, journal writes |
| `internal/db/` | Infrastructure | Schema definitions, migrations, connection management, query helpers, SQLite pragmas | Business rule enforcement, workflow logic, permission checks |
| `internal/journal/` | Infrastructure | JSONL append, flush, replay, event schema versioning, upgrader chain | Command logic, business validation (receives pre-validated events) |
| `internal/issues/` | Domain | Issue CRUD, AC management, completion recording, concern lifecycle | Output formatting, CLI flag parsing |
| `internal/claims/` | Domain | Claim acquire, release, force-steal, TTL enforcement | Workflow transition logic |
| `internal/workflow/` | Domain | Registry loading from the workflow tables (§4.6), transition validation, rule evaluation; from v0.2a also workflow authoring (validation, cancel-edge materialization, deterministic effect construction), workflow-assignment validation and write paths, and assignment resolution consumed by `internal/issues` | Issue creation, AC verification, output formatting, schema/migration ownership |
| `internal/maintenance/` | Domain | Read-only consistency audit over the projected SQLite data and journal references; structural errors, later warnings, and summary counts | CLI flag parsing, output formatting, journal writes, repair, deletion, or migration |
| `internal/policy/` | Domain | Role validation, permission matrix lookups, action authorization | Workflow state logic |
| `internal/dev/` | Domain | Compose developer write workflows (`dev start`, `dev gate`, later `dev finish`/`dev block`) into one validated journal event by orchestrating existing claim, workflow, run, and gate effect shapes | CLI flag parsing, output formatting, direct schema/migration ownership |
| `internal/runs/` | Domain | Record `agent_runs` rows with snapshot pricing, computed cost, and provider reconciliation fields | Output formatting, CLI flag parsing, read aggregations |
| `internal/format/` | Interface adapter | JSON, TOON, Markdown, YAML serialisation | Business logic, DB access |
| `internal/qa/` | Domain | QA reviews, findings, plans, plan steps (v0.3) | Issue lifecycle, claim management |
| `internal/validation/` | Domain | Contracts, assertions, required validators, outcomes (v0.4) | QA plan management |
| `internal/templates/` | Domain | Template CRUD, versioning, structural validation (v0.3) | Issue lifecycle |
| `internal/export/` | Interface adapter | Markdown tree rendering, frontmatter generation (v0.4) | Business logic |
| `internal/search/` | Domain | Literal, case-insensitive search over structured read-model fields (v0.2) | Full markdown body indexing, custom query DSL |
| `internal/stats/` | Domain | Cost, token, gate, attempt aggregations | Output formatting |
| `internal/activity/` | Domain | Projected activity timeline reads over the `activity_log` read model: chronological listing, entity/date/event filtering, and details deserialisation | Output formatting, journal writes, new persisted tables |
| `internal/sqlquery/` | Domain | Read-only SQL escape hatch over the projected SQLite read model: statement classification, mutation rejection, read-only connection opening, and generic row scanning for `fid sql` | Output formatting, journal writes, schema ownership, arbitrary write access |
| `internal/report/` | Domain | Named report SQL composition and execution (v0.2) | Output formatting |
| `internal/changes/` | Domain | Change digest aggregation over the `activity_log` read model: anchor resolution, per-entity net field changes, added/removed child facts (v0.2) | Output formatting, journal writes, new persisted tables |
| `internal/gitctx/` | Infrastructure | Read-only capture of git context (branch, base/head commit, dirty, worktree) behind an interface; bounded, timed subprocess calls; graceful no-repo degradation (v0.2) | Any git mutation (checkout, branch, commit, stash), business logic |
| `internal/handoff/` | Domain | Compose a role handoff pack (v0.3): assemble the manifest (required reads, commands, stop conditions, expected writes, live blockers, evidence obligations) from the role's context pack, the configured workflow, role permissions, and live blocker state; render via the handoff template | Output formatting, journal writes, git operations, hard-coding role or workflow policy |
| `internal/gates/` | Domain | Immutable Gate Definition Registration ingestion/resolution, registered gate execution (argv exec with cwd/timeout/env allowlist), execution-start and terminal-receipt event preparation, output bounding/sanitization/hashing, attested-receipt import validation, trust-source classification | Parsing kit/project authoring files, shell string evaluation, running agents or models, git mutation, job scheduling, CLI flag parsing, output formatting |
| `internal/envelope/` | Domain | Work envelope ingestion/validation and orchestration primitives: next-actionable with typed skip reasons, budget/circuit-breaker state, machine-readable rejection/stop reasons, attention items (v0.3a) | Hosting an LLM loop, spawning agents, scheduling, prompt authoring |
| `internal/evidence/` | Domain | Evidence records and links, obligation evaluation, freshness classification (`current`/`stale` from receipt tree identity and hashes), and read-only resolution of cited refs (commit refs via `internal/gitctx`, file/source/screenshot refs via filesystem existence; `resolved` / `missing` / `unverifiable` by type) (v0.2b+) | Command execution (that is `internal/gates/`), journal-bypassing writes, re-proving claims by re-running tools |
| `migrations/` | Infrastructure | Numbered SQL up/down migration pairs | Application logic |
| `testdata/` | Test support | Fixture DBs, golden files, sample journals | Production code |

### 3.1 Dependency Direction Constraints

Domain packages that produce journal events depend on `internal/journal` for
the event/effect data contract **and** for the journal write API — `WithLock`,
`PrepareEvent`, `DryRunProjectEvent`, `AppendEvent`, `ProjectEvent` — and read
paths additionally use the `EnsureProjectionFresh*` freshness helpers.
`internal/journal` depends on `internal/db`. **That chain is one-way and must
stay one-way**; the one-way rule, not the width of the surface, is the
constraint this section exists to protect.

Ten non-test domain packages call the write API today — `bootstrap`, `claims`,
`dev`, `issues`, `labels`, `links`, `maintenance`, `requirements`, `runs`, and
`workflow` (e.g. `internal/issues/ac.go`). An `Effect`-only dependency is
therefore **not** the project norm and never was; a package that only builds
effects is the narrow case, not the rule. Corrected 2026-08-08 (DEC-106): the
earlier "`Effect` data contract only" phrasing was introduced in Phase 19 with
the DEC-096 sync, was false for those ten packages on the day it was written,
and made an auditor reading `internal/issues` against this section find a
violation that does not exist.

- **`internal/workflow` → `internal/journal` is permitted and includes the
  journal write API** (v0.2a; DEC-089, corrected by DEC-105). The authoring core
  returns `[]journal.Effect` so the effect shape is never duplicated away from
  the event schema it serialises into, and the authoring **write paths**
  (`define.go`, `assign.go`) additionally use `WithLock`,
  `EnsureProjectionFreshWithLock`, `PrepareEvent`, `DryRunProjectEvent`,
  `AppendEvent`, and `ProjectEvent`. The read path (`list.go`) uses
  `EnsureProjectionFresh`. This follows directly from DEC-089:
  `internal/workflow` **owns** workflow authoring, assignment, and resolution,
  and owning a write path means appending the event. Pushing append and
  projection up into `cmd/fid` would be the worse design.
  - **Superseded wording (DEC-096, corrected by DEC-105):** DEC-096 described
    this edge as "bounded to the `Effect` struct — no append, no projection, no
    replay." That was an accurate description of the package **at Task 2**, when
    it only built effects, and it was wrongly stated as a permanent invariant.
    Tasks 3-4 added the write paths DEC-089 had already authorized. **The
    boundary that still holds** is the one in §3: no issue creation, no AC
    verification, no output formatting, and no schema or migration ownership.
  - **Replay stays out.** `internal/workflow` must not add resolution logic to
    the projection or replay path — replay re-resolving assignments would make
    history non-deterministic (DEC-076).
- **Consequence for tests:** because `internal/journal` imports `internal/db`,
  **no in-package (`package db`) test file in `internal/db` may import
  `internal/workflow`**, and no in-package `internal/journal` test file may
  either. Go compiles in-package test files into the package itself, so such an
  import closes a cycle and fails `typecheck` for the whole repository. Tests
  needing both belong in the external `package db_test` / `package journal_test`
  packages, which are not compiled into the package under test.

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
| `task_completions` | attempt_number, status, actual_tokens, commit_hashes, files_touched, gate_results, human_intervention | issue_id → issues (task/subtask only) | Append-only. Previous attempts marked superseded. |
| `concerns` | summary, status, body_md_path | issue_id → issues, raised_in_completion_id → task_completions, escalated_to_esc_id → escalations, triaged_in_event_id → triage_events | Status: open → addressed / escalated / wontfix. |
| `agent_runs` | role, model, tokens_in, tokens_out, cost_usd, price_in/out_per_mtok, token_breakdown_json, tool_call_count, provider_run_id, provider_cost_usd | issue_id → issues, completion_id → task_completions | Snapshot pricing at run end plus optional provider reconciliation fields. v0.3 adds optional `review_id` when QA reviews ship. |
| `claims` | ref_type, ref_id, claimed_by, claimed_at, ttl_seconds | Polymorphic ref (v0.1: issue; v0.3: qa_review) | At most one live claim per entity. PK enforces exclusivity. |
| `roles` | name, description, is_human, active | Seeded from default profile | Existence is validity rule for all `_role` columns. |
| `labels` / `issue_labels` | name, colour; issue_id, label | issue_labels: issue_id → issues, label → labels | Free-form tags. |
| `links` | from_type, from_id, relation, to_type, to_id | Polymorphic endpoints | v0.1 scope: issue ↔ issue, issue → requirement. Auto-inverse. |
| `activity_log` | journal_event_id, ref_type, ref_id, event, details, actor | Polymorphic ref, journal_event_id → journal | Projected from journal events. Read-model only. |
| `projection_meta` | journal_file, journal_offset, last_event_id, last_event_hash, projection_schema_version, replay_status | Single-row DB projection checkpoint | Lets reads detect missing, stale, or inconsistent SQLite state before serving results. |

**Governance support tables (v0.1 schema, v0.2 CLI unless noted):**

| Entity | Key Fields | Notes |
|---|---|---|
| `decisions` | summary, body_md_path, area, status | FK target for escalations.resolution_dec. |
| `escalations` | status, trigger_class, raised_by_role | FK target for concerns and triage_events. |
| `triage_events` | verdict (RESUME/BLOCK), source, source_ref_id, reasoning_md | FK target for concerns.triaged_in_event_id. |
| `requirements` | type (PRD/TRD), summary, source_doc, source_section, last_synced_hash | FK target for links (issue → requirement). v0.1 ships minimal `fid requirement sync` so requirement link targets can be created before v0.2 governance reads. |

**Evidence foundation (v0.2b, DEC-059/060/061):**

| Entity | Key Fields | Notes |
|---|---|---|
| Gate Definition Registrations (exact projection schema finalized with the backend port) | registration id/version/status; definition id/version/digest and complete definition snapshot; project/overlay/repository binding digest; authorization provenance | Immutable operational authorization created from project-owned authoring configuration through the kit registration contract. fid resolves registrations by reference and never reads mutable authoring configuration at execution time. |
| `gate_executions` | id, gate_id, gate_version, issue_id, attempt_number, agent_run_id, trust_source, argv_json, cwd, timeout_seconds, started_at, started_by_role | Append-only execution-start facts, written before process spawn. A row with no `gate_receipts` child is an abandoned start (maintenance warning), never a pass. |
| `gate_receipts` | id, execution_id (unique), receipt_schema, producer, trust_source, status (pass/fail/timed_out/cancelled/harness_error), exit_code, signal, duration_ms, branch, base_commit, head_commit, tree_hash, dirty, capture_source, toolchain_fingerprint, stdout_excerpt, stderr_excerpt, stdout_hash, stderr_hash, stdout_bytes, stderr_bytes, truncated, redacted, log_artifact_ref, ended_at | Immutable terminal receipts; at most one per execution. Legacy `gate_runs` rows remain and project as `trust_source='reported'`. |
| `evidence_records` | id, evidence_type (receipt/commit/source_state/screenshot/artifact/human_outcome), ref (receipt id, commit hash, path, ...), recorded_by_role, human_actor, human_outcome, note, recorded_at | Typed first-class evidence. Human outcomes record actor, checklist item, outcome, timestamp, optional note. |
| `evidence_links` | evidence_id, ref_type (ac_verification/qa_review/val_outcome/workflow_gate), ref_id | Connects evidence to the obligations/facts it supports. |

**Work envelope (v0.3a, DEC-063):**

| Entity | Key Fields | Notes |
|---|---|---|
| `work_envelopes` | id, version, primary_issue_id, selected_issue_refs, risk_lane, authority_json, protected_decisions_json, budgets_json, stop_conditions_json, status | Versioned execution contract referencing the existing Issue graph; no Work Item table or identity. Obligations reference registered gates and selected Issue subjects. Exact normalized projection and revision rules follow joint scenario design (ADR-0005). |

**Workflow engine (v0.2):**

| Entity | Key Fields | Notes |
|---|---|---|
| `workflows` | id (stable slug), name, applies_to | v0.1: hard-coded in app code as enum strings. Seeded into the tables and made the runtime source of truth in v0.2a (§4.6). |
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

**Bootstrap contract:**

`fid init` is the first journaled write in a workspace. It creates the
workspace directories, writes canonical `.fid/config.toml` and
`.fid/pricing.toml`, applies the current SQLite schema, then appends a
`project.initialized` event to the current UTC journal file. The event's
effects seed fid entity state that must survive delete-and-rebuild, including
the seven default `roles` rows. The projection step writes the initial
`projection_meta` checkpoint for that event. Config and pricing files are
committed canonical configuration files; they are not reconstructed from the
journal.

After a successful init, an otherwise empty project has exactly one journal
event and no issues. `fid rebuild` starts from an empty schema, replays
`project.initialized`, and must recreate the roles table and projection
checkpoint without consulting the old SQLite file. A workspace with no
`project.initialized` event is not a valid v0.1 fid workspace.

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

**Registered gate execution is deliberately two events (v0.2b, DEC-060).**
External process execution and journal append cannot share one filesystem
transaction, and fid models that honestly rather than pretending atomicity:

1. `gate.execution_started` — appended (one event, one insert effect) *before*
   the process spawns. This is a normal single logical write.
2. The process runs to termination under fid's timeout.
3. `gate.receipt_recorded` — appended only after termination, carrying the
   terminal status, repository identity, and bounded sanitized output.

Each of the two appends individually obeys the one-write-one-event rule. If
fid crashes between them, the start fact survives with no receipt; `fid
maintenance check` reports it as an abandoned execution, and it satisfies
nothing. A crash may lose a pass; it must never fabricate one. Receipts are
immutable once appended — corrections are new executions, never edits.

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
| `event_id` | yes | Stable generated event ID with shape `EVT-<UTC YYYYMMDD>-<sequence>`. Used by `activity_log.journal_event_id`, replay diagnostics, and projection metadata. |
| `event_schema` | yes | Positive integer schema version. v0.1 writes `1`. |
| `event_type` | yes | Semantic command/event name, e.g. `project.initialized`, `issue.created`, `issue.transitioned`, `claim.acquired`. |
| `occurred_at` | yes | UTC ISO 8601 timestamp assigned once before append. |
| `actor` | yes | Object with `role` and optional `actor` string. `role` is the `--by` role for CLI writes. |
| `command` | yes | Sanitized command metadata: `argv`, `cwd`, `fid_version`, and no environment variables or sensitive flag values. |
| `effects` | yes | Ordered array of row-level effects applied atomically during projection. |

**Event ID and journal file allocation:**

The writer reads the UTC clock exactly once while holding `.fid/write.lock`.
That instant supplies:

- `occurred_at`;
- the journal file date, `.fid/journal/<UTC YYYY-MM-DD>.jsonl`;
- the event ID date component, `EVT-<UTC YYYYMMDD>-<sequence>`.

The sequence component is a zero-padded decimal number allocated per UTC journal
file. Under the lock, fid scans existing events in the target file, takes the
largest sequence already present for that date, and writes the next sequence
(`000001` for the first event in that UTC day). If a day exceeds six digits,
the width may grow; readers treat the suffix as decimal, not fixed-width text.

If the chosen UTC date is earlier than the greatest existing journal file date,
fid exits non-zero with a clock-skew error and does not append. If the clock
crosses midnight after the single timestamp read, the event still goes to the
file chosen from that read. During replay, journal files are processed in
ascending `YYYY-MM-DD` filename order, then line order within each file.

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

**Bootstrap event example:**

`project.initialized` is the only event type that may be written before a
workspace already has valid projection metadata. Its effects include the
default role rows. The projection step writes the initial `projection_meta`
row for the event after applying its effects. The project prefix is stored in
`.fid/config.toml`; it may also appear in event metadata for diagnostics, but
config remains the canonical source for project configuration.

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
- Any path that mutates the SQLite projection acquires the same
  `.fid/write.lock`, including automatic catch-up before a read, explicit
  `fid rebuild`, and explicit migration/vacuum commands. A read command whose
  `projection_meta` already matches the journal may serve from SQLite without
  the lock; a read command that detects a stale or missing projection waits for
  the lock, catches up or rebuilds, and only then serves output. It never
  serves stale data while another process is projecting a committed event.
- SQLite transactions remain required inside the write-lock boundary for
  projection atomicity. SQLite locking alone is not the write serializer,
  because the journal append is the commit point.
- `fid maintenance vacuum --by <role>` is an explicit derived-DB maintenance
  operation, not a domain write. After authorizing the role, it holds the same
  `.fid/write.lock`, refreshes the projection through the held-lock path, then
  runs SQLite `VACUUM` against `.fid/fid.sqlite`. It appends no journal event,
  creates no activity row, and reports the database byte sizes before and
  after compaction.
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
- Replay processes `.fid/journal/YYYY-MM-DD.jsonl` files in ascending UTC date
  filename order, then line order within each file.
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
  position, the projection checkpoint is treated as diverged from the canonical
  journal: fid rebuilds the read model from the journal before serving the read
  rather than trusting the stale projection. (Malformed or
  precondition-conflicting journal entries encountered during that replay still
  stop before answering — see Malformed journal below.) fid never silently
  skips, truncates, or reorders journal entries.

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

**Rebuild-connection exception (Phase 18, DEC-071).** The connection used by
`fid rebuild` additionally sets `synchronous = OFF` and `journal_mode = MEMORY`
(applied through the DSN `_pragma` mechanism so every pooled connection carries
them). `foreign_keys` and `recursive_triggers` stay ON, and the normal
command-path connection keeps full durability — this relaxation is scoped to
rebuild only. It is safe because rebuild reconstructs a *derived* store from the
canonical journal, from scratch and under the workspace write lock: a crash
mid-rebuild discards the partial DB and re-rebuilds (ADR-0001, TRD-ARCH-001), so
durability during reconstruction has no contract value. Pragmas change durability
and journaling, not row data, so the rebuilt projection stays byte-identical.
`journal_mode = MEMORY` (rather than `OFF`) preserves transactional rollback, so
per-event error surfacing is unchanged.

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
`fid sql`, direct DB edits) are created only for tables present in the active
schema migration:
- Kind-invariant rules: hierarchy enforcement on `issues.parent_id`,
  sequence constraints, kind guards on `acceptance_criteria`,
  `task_completions`, `epic_metadata`, `triage_events`, `decisions`,
  `carry_forward`, and `qa_reviews` when the QA layer ships.
- Polymorphic ref integrity on `links`, `claims`, `comments`,
  `activity_log`.
- Link scope and relation compatibility.
- Claim scope (supported `ref_type` values per schema version).
- Comment scope (v0.2 base ref_type whitelist; v0.3 adds `qa_review`).
- Kind immutability (BEFORE UPDATE OF kind on `issues`).
- Mutually exclusive scope columns on `qa_reviews` when v0.3 ships
  (CHECK constraint).
- Validation outcome row-local evidence/rationale requirements
  (CHECK constraints).

What does **not** get a trigger: workflow transition rules, permissions,
AC verification freshness, required-validator satisfaction for validation
outcomes, template validation. These need rich context, ordering,
multi-table reads, or runtime state that SQL doesn't have.

### 4.6 Default Workflows (migrated to the workflow tables in v0.2a)

v0.1 shipped these workflow definitions hard-coded in `internal/workflow/`.
**As of Phase 18 (v0.2a slice 1) they are seeded verbatim into the workflow
tables (`03_TRD_SCHEMA.md` §3) by migration `0007`, and the tables are the
runtime source of truth.** `internal/workflow` now loads a read-only registry
from those tables on each read; the Go `definitions` map is retained solely as
frozen parity fixtures for tests, and no production code path reads it. The
tables below therefore document the seeded default data, not a hard-coded
runtime structure.

The defaults are seeded by migration DDL rather than journal events: they are
system reference data, which keeps v0.1 history unchanged and satisfies foreign
key ordering on rebuild. Migration `0010` also seeds the v0.2
`escalation-default` reference workflow; user-defined workflow definitions and
assignments remain journal events when `fid workflow define/assign/list` ships.

`issues.workflow_id` and `issues.status` are foreign keys into `workflows` and
`workflow_states` (migration `0008`), so an issue cannot hold a workflow/state
pair the tables do not define. Cancel edges (`<state>` → `cancelled` for every
non-cancelled state) are seeded explicitly rather than resolved by a code
special-case.

`workflow_assignments` is a journal-projected table (migration `0009`), not
seed data. At issue creation, workflow resolution is kind-scope assignment,
then project-scope assignment, then the built-in kind default
(`epic-default`, `task-default`, or `subtask-default`). With no assignment
rows, resolution remains byte-identical to v0.1.

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

### 4.7 Gate Run And Receipt Persistence

**v0.1 (legacy reported path).** `fid dev gate` writes append-only `gate_runs`
rows. It does not create or mutate `task_completions` directly, because
completions represent finalized attempts.

**v0.2b (trusted path, DEC-059/060).** `fid gate run <gate-id>` executes the
registered definition and persists a `gate_executions` start fact plus a
terminal `gate_receipts` row (two events; see §4.3). Receipts snapshot the
definition fields they ran with (gate ID, version, argv, cwd, timeout) and the
repository identity at capture. `gate_runs` remains queryable as
`trust_source='reported'` history; the deprecated `fid dev gate --cmd
--exit-code` surface keeps appending to it as an explicit reported-evidence
import until a recorded removal decision. Trusted obligations
(`required_gates_current`, PRD-WF-015) read only `gate_receipts` with
`trust_source IN ('executed','attested')`, status `pass`, and current
freshness (PRD-EVID-003). Trusted-mode completions link receipts by ID; the
v0.1 `task_completions.gate_results` JSON snapshot continues only for the
legacy reported path.

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
| Gate Definition Registration contract | Kit adapter, `fid gate run`, workflow obligations, external runners | Versioned immutable registrations bind complete definition, authorization, project/overlay/repository identity, and digests. Authoring-file changes require a new activation; callers cannot substitute fields at execution. | Shared conformance fixtures; registration mismatch/staleness/substitution tests; receipt binding tests. |
| Receipt JSON/TOON shape (v0.2b) | Workflow kits, `fid evidence check`, reports, attested-runner import | Versioned via `receipt_schema`. Field additions non-breaking; consumers reject unknown-major schemas. Import path validates the attestation contract version. | Golden receipt fixtures per schema version; import rejection tests. |
| Capability discovery (`fid capabilities`, v0.2b) | Workflow kits and integrations | JSON/TOON; named capability flags only added, never silently removed. Consumers detect missing capabilities and stop; no silent fallback from trusted receipts to reported evidence. | Contract snapshot test; downgrade-detection test. |

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
| TRD-SEC-005 | Project-owned machine-readable configuration is the Gate Definition authoring source, but fid shall execute only active immutable registrations produced through the kit-owned validation and activation contract. A mutable file change alone never changes runtime authority; unknown, stale, mismatched, unauthorized, or field-substituted registrations fail closed. | Shared contract tests for unknown/stale/mismatched registration, unauthorized change, field substitution, and source change without activation. |
| TRD-SEC-006 | Registered gate execution shall use structured argv with no implicit shell expansion. If an explicit shell-mode escape is retained (PRD-GATE-008), it shall be policy-controlled, visibly flagged in the registry and on every receipt, and off by default. (v0.2b) | Test: argv containing shell metacharacters is passed literally; shell-mode without the policy opt-in is rejected. |
| TRD-SEC-007 | Gate execution shall apply the definition's environment allowlist, scrub sensitive flag values and environment data before journal persistence, bound captured stdout/stderr, and record redaction/truncation metadata on the receipt. (v0.2b) | Tests: secret-like env/args do not appear in journal bytes; oversized output is truncated with correct hashes/flags. |
| TRD-SEC-008 | The Git-committed journal shall carry only safe metadata, excerpts, and hashes for gate output. Optional full logs shall be stored as content-addressed local (or CI) artifacts under an explicit retention policy, outside the journal. (v0.2b) | Journal byte inspection test; artifact store retention config test. |
| TRD-SEC-009 | A receipt imported from an external producer shall declare its trust source and attestation contract version; import shall never silently satisfy a trusted-execution obligation with reported or unversioned evidence. The no-network, no-daemon, local-first default is preserved — import reads local files only. (v0.2b) | Import tests: missing/`reported` trust source cannot satisfy `required_gates_current`; no listener code. |

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
| `gopkg.in/yaml.v3` | `fid issue get --format yaml` output and JSON-parity decode in tests (v0.1) | In-tree YAML emitter (more correctness/maintenance burden, still needs a decoder for parity tests); defer YAML | MIT | DEC-012 |
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
| TRD-ARCH-015 | Template-driven validation shall declare required sections per template kind. Where a workflow profile still requires completion Markdown, `fid dev finish` shall validate files against templates and reject malformed artifacts; structured-completion profiles (PRD-WF-014) do not require the report at all. | v0.3 | Integration test: submit completion report missing `## Concerns` heading; assert rejection. |
| TRD-ARCH-016 | The executor boundary: fid shall execute only active immutable Gate Definition Registrations approved from project configuration and supplied by reference. It shall not parse the kit's authoring source, run agents, host models, mutate git, evaluate shell strings by default, or schedule jobs. External attested runners are secondary; fid still owns registration resolution, receipt schema/validation, journal events, linkage, and state-transition enforcement. | V2-aligned backend port | Code review + conformance tests: unregistered/substituted command rejected; no git-mutating or listener code paths in `internal/gates/`. |
| TRD-ARCH-017 | Gate receipts shall be immutable, identified rows linked to their execution-start facts, produced via the two-event model in §4.3. Rebuild shall reproduce receipt, evidence, and freshness state exactly from the journal. | v0.2b | Golden replay fixtures including an abandoned start; delete-and-rebuild parity test over receipts/evidence. |
| TRD-ARCH-018 | Machine-readable integration shall go through documented JSON/TOON contracts (receipts, capability discovery, transition rejections, envelope primitives), never through SQLite internals or `.fid/` file layout assumptions. | v0.2b+ | Contract snapshot tests; docs list the supported surfaces. |

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
- [ ] Attestation contract detail for companion runners (DEC-059): key/identity
  model for `producer`, contract version negotiation, and whether attested
  receipts require a local co-signature. Decide at v0.2b phase planning.
- [ ] Toolchain fingerprint contents (PRD-GATE-003): start with Go toolchain +
  linter versions from `go version` / tool output, or a generic hash of a
  declared tool list? Decide at v0.2b phase planning.
- [ ] Full-log artifact store defaults (TRD-SEC-008): retention window, size
  caps, and whether the store lives under `.fid/artifacts/` (gitignored) only
  or may also target a CI artifact path. Decide at v0.2b phase planning.

*Last updated: 2026-07-23*
