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
    ├──► internal/stats/      Cost, tokens, gates, attempts aggregations
    └──► internal/report/     Named reports (v0.2+)
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
| `internal/stats/` | Domain | Cost, token, gate, attempt aggregations | Output formatting |
| `internal/report/` | Domain | Named report SQL composition and execution (v0.2) | Output formatting |
| `migrations/` | Infrastructure | Numbered SQL up/down migration pairs | Application logic |
| `testdata/` | Test support | Fixture DBs, golden files, sample journals | Production code |

## 4. Data Model

### 4.1 Entities

The full schema is defined in `scaffold/initial-spec.md` §6. This section
captures the entity inventory, key relationships, and design rationale. Column
lists are indicative; refer to the spec for exhaustive definitions.

**Core issue layer (v0.1):**

| Entity | Key Fields | Relationships | Notes |
|---|---|---|---|
| `issues` | id, kind, title, status, workflow_id, priority, sequence, spec_md_path, est_tokens, owner_role | parent_id → issues (hierarchy), cycle_id → cycles | Unified table; `kind` discriminator: epic, task, subtask. |
| `epic_metadata` | scheduling, brief_md_path, target_start, target_end | issue_id → issues (1:1, epic only) | Sequential or prioritised scheduling. |
| `acceptance_criteria` | ordinal, text | issue_id → issues (task/subtask only) | Verification state in `ac_verifications`. |
| `ac_verifications` | verifier_role, status, evidence_type, evidence_ref, verified_at | ac_id → acceptance_criteria, completion_id → task_completions | Append-only. Multiple per AC per role. |
| `task_completions` | attempt_number, status, actual_tokens, commit_hashes, files_touched, gate_results | issue_id → issues (task only) | Append-only. Previous attempts marked superseded. |
| `concerns` | summary, status, body_md_path | issue_id → issues, raised_in_completion_id → task_completions, escalated_to_esc_id → escalations, triaged_in_event_id → triage_events | Status: open → addressed / escalated / wontfix. |
| `agent_runs` | role, model, tokens_in, tokens_out, cost_usd, price_in/out_per_mtok | issue_id → issues, completion_id → task_completions, review_id → qa_reviews | Snapshot pricing at run end. |
| `claims` | ref_type, ref_id, claimed_by, claimed_at, ttl_seconds | Polymorphic ref (v0.1: issue; v0.3: qa_review) | At most one live claim per entity. PK enforces exclusivity. |
| `roles` | name, description, is_human, active | Seeded from default profile | Existence is validity rule for all `_role` columns. |
| `labels` / `issue_labels` | name, colour; issue_id, label | issue_labels: issue_id → issues, label → labels | Free-form tags. |
| `links` | from_type, from_id, relation, to_type, to_id | Polymorphic endpoints | v0.1 scope: issue ↔ issue, issue → requirement. Auto-inverse. |
| `activity_log` | journal_event_id, ref_type, ref_id, event, details, actor | Polymorphic ref, journal_event_id → journal | Projected from journal events. Read-model only. |

**Stub tables (v0.1 schema, v0.2 CLI):**

| Entity | Key Fields | Notes |
|---|---|---|
| `decisions` | summary, body_md_path, area, status | FK target for concerns.escalated_to_esc_id and escalations.resolution_dec. |
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
| `views` | name, query_dsl, scope | Saved queries. |
| `comments` | ref_type, ref_id, author_role, body_md | Polymorphic ref scoped to issue, qa_review, decision, escalation, requirement. |

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

1. Validate against the current DB state (permissions, workflow rules,
   kind constraints, claim state).
2. Build a fully resolved event: one semantic envelope
   (`event_type`, `event_schema`, actor metadata, command metadata) with an
   ordered `effects` array of row operations (`insert` / `update` / `delete`
   with table name, PK, and resolved values).
3. Append and flush one JSONL line to
   `.fid/journal/YYYY-MM-DD.jsonl` — this is the commit point.
4. Apply the event's effects to SQLite in a single DB transaction.

If journal append fails, the DB is not modified. If DB projection fails after
journal append, the next CLI invocation catches the DB up before answering.

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

**Concurrency:**

- Single-worktree: safe via SQLite's default locking.
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

SQLite triggers (backstop for non-CLI write paths — MCP server, web
dashboard, `fid sql`, direct DB edits):
- Kind-invariant rules: hierarchy enforcement on `issues.parent_id`,
  sequence constraints, kind guards on `acceptance_criteria`,
  `task_completions`, `epic_metadata`, `qa_reviews`, `triage_events`,
  `decisions`, `carry_forward`.
- Polymorphic ref integrity on `links`, `claims`, `comments`,
  `activity_log`.
- Link scope and relation compatibility.
- Claim scope (supported `ref_type` values per schema version).
- Comment scope (v1 ref_type whitelist).
- Kind immutability (BEFORE UPDATE OF kind on `issues`).
- Mutually exclusive scope columns on `qa_reviews` (CHECK constraint).

What does **not** get a trigger: workflow transition rules, permissions,
AC verification freshness, template validation. These need rich context or
runtime state that SQL doesn't have.

## 5. Interfaces and Contracts

| Contract | Consumers | Compatibility Policy | Verification |
|---|---|---|---|
| CLI command surface (`fid <noun> <verb> [flags]`) | Agents, humans, shell scripts | Stable within a major version. Flag additions are non-breaking; flag removals require a deprecation cycle. | Integration tests per command. |
| Journal event format (JSONL with semantic envelope + effects) | `fid rebuild`, `fid merge --journal`, future tools | Forward-compatible via `event_schema` versioning and upgrader chain. Old events never rewritten. | Golden-file replay tests per schema version boundary. |
| SQLite schema (tables, views, indexes, triggers) | CLI queries, `fid sql`, MCP server, web dashboard | Versioned via numbered migrations. | Migration up/down tests. |
| JSON output (default for all read commands) | Agents, `jq` pipelines, downstream tools | Field additions are non-breaking; field removals require deprecation. | Snapshot tests per command output. |
| TOON output (opt-in via `--format toon`) | Agents optimising for token efficiency | Secondary format; shape mirrors JSON. | Parity tests against JSON output. |
| Markdown export with YAML frontmatter (v0.4) | Humans, git diffs, stakeholder review | Frontmatter schema documented; field additions non-breaking. | Export round-trip tests (structured fields survive export/re-parse). |
| MCP server resources/prompts/tools (v0.3) | AI agents via Model Context Protocol | Experimental until v0.4. | Integration tests against MCP SDK. |
| `.fid/config.toml` | CLI, export, future MCP server | Additive sections; no removal without migration. | Config parsing tests with missing/extra keys. |
| `.fid/pricing.toml` | Cost computation in `agent_runs` | Additive model entries; user-overridable. | Test with missing model, custom model, default model. |

## 6. Security and Privacy

| ID | Requirement | Verification |
|---|---|---|
| TRD-SEC-001 | The system shall not open any network listener or transmit data externally. All data stays local to the filesystem. | Code review; no `net.Listen` or HTTP client calls outside MCP server (v0.3, localhost only). |
| TRD-SEC-002 | The `fid sql` escape hatch shall reject any statement other than SELECT, EXPLAIN, PRAGMA query, and WITH (read-only). | Integration test: attempt INSERT/UPDATE/DELETE/DROP via `fid sql`; assert rejection. |
| TRD-SEC-003 | Journal files shall not contain secrets. CLI shall strip environment variables and sensitive flag values before recording command metadata in events. | Code review; test that event `command_metadata` contains no env vars or `--password`-style values. |
| TRD-SEC-004 | SQLite triggers shall prevent structural corruption even from non-CLI write paths (MCP, dashboard, direct SQL). | Trigger tests: attempt invalid inserts via raw SQL; assert RAISE fires. |
| TRD-SEC-005 | The MCP server (v0.3) shall bind to localhost only and shall not expose the full CLI surface. | Config default; integration test that remote connections are refused. |

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
| TRD-ARCH-013 | The MCP server shall expose a curated subset of CLI tools, not the full surface. Resources, prompts, and tools shall be the three MCP primitives. | v0.3 | Code review; tool inventory test. |
| TRD-ARCH-014 | Markdown export shall include YAML frontmatter with structured fields. The frontmatter schema shall be documented and tested for round-trip stability. | v0.4 | Export → parse frontmatter → assert field values match DB. |
| TRD-ARCH-015 | Template-driven validation shall declare required sections per template kind. `fid dev finish` and similar commands shall validate files against templates and reject malformed artifacts. | v0.3 | Integration test: submit completion report missing `## Concerns` heading; assert rejection. |

## 11. Open Technical Questions

- [ ] Should the TOON marshaller support the full TOON spec or only the
  subset needed for fid's tabular output? Start minimal and extend on demand.
- [ ] Should `fid rebuild` support partial replay (catch-up from a known
  checkpoint) or always replay from scratch? Full replay is simpler and
  fast enough at single-project scale; checkpoint is a forward-compatible
  optimisation.
- [ ] What is the MCP server's Go SDK choice for v0.3? Depends on community
  maturity at implementation time.

*Last updated: 2026-05-24*
