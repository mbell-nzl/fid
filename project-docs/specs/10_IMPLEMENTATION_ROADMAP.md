# IMPLEMENTATION ROADMAP - fid

> This roadmap sequences the work from first useful build to release. Keep it
> phase-oriented rather than sprint-oriented so agents can plan and execute one
> coherent phase at a time.

## Guiding Principles

1. **Dogfood early.** v0.1 should replace ~70% of daily interaction with whet's
   task artifacts. The author uses fid for real project work before v0.2 starts.
2. **Cut scope, not quality.** Defer features rather than shipping fragile
   foundations. Stub tables are acceptable; half-implemented workflows are not.
3. **Stabilize contracts early.** The journal event format, CLI command surface,
   JSON output shapes, and SQLite schema are harder to change than feature
   additions. Lock them down in v0.1.
4. **Keep phases coherent.** Each phase has a clear goal, exit criteria, and a
   small set of related tasks. The Tech Lead slices phases from within each
   roadmap version.
5. **Prove requirements explicitly.** Phase plans turn scoped requirements into
   `VAL-NNN` assertions before implementation begins.

## Phase 0 - Specification And Foundation

**Goal:** Make the project implementable by agents: stack decisions, repository
scaffold, quality commands, walking skeleton, and validated specs.

| Item | Description | PRD/TRD Refs |
|---|---|---|
| Stack decisions | Go, Cobra, pure-Go SQLite, golang-migrate, goldmark, testify. All accepted. | TRD §1 |
| Repo scaffold | Go module, `cmd/fid/`, `internal/`, `migrations/`, `testdata/`, Makefile, `.golangci.yml`, CI stub. | TRD §2–3 |
| Walking skeleton | `fid --help` exits 0; `go test ./...` passes; `go build` produces a binary. | TRD-PERF-007, PRD-OPS-009 |
| Agent workflow validation | Quality gate commands in `AGENTS.md` are concrete and green. First `validation-contract.md` created from template. | QP §5, WF §7 |
| Spec audit | Specs reviewed for completeness and consistency before implementation. | — |

**Exit criteria:**

- `go build ./cmd/fid` produces a static binary under 30 MB.
- `go test ./...` passes with zero test files (no false failures).
- `gofmt -l .` and `golangci-lint run ./...` are green.
- Quality gate commands are concrete and documented in `AGENTS.md`.
- Phase 1 can be planned without unresolved stack choices.

---

## v0.1 — Core Issue Layer

**Goal:** Ship the foundational persistence, issue lifecycle, claims, links,
telemetry, and output formats. This alone replaces ~70% of daily interaction
with whet's task artifacts.

### v0.1a — Persistence Foundation

| Item | Description | PRD/TRD Refs |
|---|---|---|
| `.fid/` directory structure | Create `.fid/`, `.fid/journal/`, `.fid/.gitignore` (ignores `fid.sqlite`). Store `config.toml` and `pricing.toml`. | PRD-INIT-001, TRD §4.3 |
| SQLite schema and migrations | All v0.1 tables: `issues`, `epic_metadata`, `acceptance_criteria`, `ac_verifications`, `gate_runs`, `task_completions`, `concerns`, `links`, `claims`, `agent_runs`, `labels`, `issue_labels`, `activity_log`, `projection_meta`, `roles`, `requirements`. Stub tables (schema only, no CLI): `decisions`, `escalations`, `triage_events`. | TRD §4.1, spec §6 |
| SQLite triggers | v0.1 triggers only for v0.1 tables: issue kind hierarchy/sequence guards, AC/completion/epic metadata kind guards, triage/decision stub kind guards, polymorphic ref integrity for `links`/`claims`/`activity_log`, link scope validation, claim scope, and issue kind immutability. Comment scope ships with the v0.2 comments table; `qa_reviews` checks ship with the v0.3 QA layer. `PRAGMA foreign_keys = ON; PRAGMA recursive_triggers = ON`. | TRD §4.4–4.5, spec §6.14 |
| `issues_enriched` view | LEFT JOIN open concern count and latest triage status. Indexes on `concerns(issue_id, status)` and `triage_events(issue_id, created_at)`. | TRD-ARCH-008, spec §6.1.1 |
| Write/projection lock | Checkout-scoped `.fid/write.lock` acquired before freshness checks, validation, event construction, journal append/fsync, SQLite projection, automatic read catch-up, explicit rebuild, migrate, and vacuum. Lock metadata records PID, hostname, started_at, command, and actor role when available. Clean reads whose `projection_meta` matches the journal do not need the lock. | PRD-NFR-008, PRD-OPS-003, TRD §4.3 |
| Journal writer | Append-only JSONL to `.fid/journal/YYYY-MM-DD.jsonl` using UTC dates. Flush before returning. Event schema v1 includes `event_id` (`EVT-<UTC YYYYMMDD>-<sequence>` allocated per UTC journal file under the lock), `event_schema`, `event_type`, timestamp, actor, sanitized command metadata, and ordered effects. Mutable-state effects include targeted preconditions. One logical write = one event. Runs inside `.fid/write.lock`; journal append is the commit point. | PRD-OPS-001–002, PRD-NFR-008, TRD-ARCH-001–004 |
| Journal replayer | Replay journal into SQLite. Event schema versioning with pure upgrader functions. Catch-up on startup when DB is behind or missing. Malformed trailing line = corruption, reported not skipped. Failed effect preconditions stop replay with a structured logical-conflict error. | PRD-OPS-003–004, PRD-OPS-012, TRD-ARCH-002–003 |
| `fid rebuild` | Replay from scratch into fresh DB. Sub-second at single-project scale. Foundational — invoked on startup when DB is stale. | PRD-OPS-003–004, TRD-PERF-002 |
| `fid init` | Create `.fid/`, write canonical `config.toml` and `pricing.toml`, apply schema, append the first `project.initialized` journal event, and project it to seed the `roles` table from the default profile (7 roles) plus initial `projection_meta`. Accept `--prefix` for issue ID prefix. Reject re-init. | PRD-INIT-001–004 |
| ID generation | `<PREFIX>-<hash>` with progressive length scaling (6 → 8 hex). Collision retry on insert. | TRD §4.2, spec §6.0 |
| Output format engine | JSON (default), TOON (opt-in), markdown, YAML. `--format` flag on all read commands. | PRD-OUT-001–004, TRD-ARCH-012 |

### v0.1b — Issue Lifecycle

| Item | Description | PRD/TRD Refs |
|---|---|---|
| `fid issue create` | Create epic/task/subtask. `--kind`, `--title`, `--parent-id`/`--epic`, `--priority`, `--sequence`, `--spec-md`, `--label`, `--est-tokens`, `--due-date`, `--inbox`, `--owner-role`, `--owner-actor`. Enforce kind hierarchy (epic: no parent; task: optional epic parent; subtask: required task parent). Assign default workflow per kind. | PRD-ISSUE-001–005, 013–016, 019–020 |
| `fid issue get` | Return issue + ACs + latest verifications per role + open concern count + triage status (from `issues_enriched`). Support `--section` for on-demand markdown parsing. Support `--format`. | PRD-ISSUE-006, 018 |
| `fid issue list` | Filter by kind, status category, epic, label, priority, inbox state. `--order sequence\|priority`. | PRD-ISSUE-007 |
| `fid issue update` | Update title, priority, labels, spec_md_path, est_tokens, due_date, owner_role, owner_actor on existing issues. Journal and project to activity_log. | PRD-ISSUE-015, 021 |
| `fid issue transition` | Validate against workflow state machine. Enforce transition rules. `--to <state> --by <role>`. Rich error messages naming unmet conditions and entity IDs. | PRD-WF-001–010, 013 |
| Epic metadata | `epic_metadata` table with scheduling (sequential/prioritised), brief_md_path, target dates. Created alongside epic-kind issue. | PRD-ISSUE-001, spec §6.1 |
| `fid ready` | Return next actionable task. Sequential epics: lowest-sequence unfinished, only when unblocked. Prioritised epics: highest-priority unblocked. Cross-project: top task across all. Empty result with explanation when nothing ready. | PRD-ISSUE-008–012 |
| Hard-coded workflows | Default state machines for epic, task, subtask from `02_TRD.md` §4.6 implemented in code. `workflow_id` as enum string. Status CHECK-constrained per workflow. Status category mapping for portable queries. | PRD-WF-005–006, TRD-ARCH-009 |
| Hard-coded transition rules | `all_acs_verified` (parameterised by roles), `has_completion_report`, `concerns_addressed`, `all_child_tasks_done_or_cancelled`. Named Go functions. | PRD-WF-007–010 |

### v0.1c — ACs, Completions, Concerns

| Item | Description | PRD/TRD Refs |
|---|---|---|
| `fid issue ac add` | Add ordered acceptance criteria to tasks and subtasks. Reject on epics. `--text`, `--ordinal`. | PRD-AC-001–002 |
| `fid issue ac verify` | Multi-role verification. `--by <role> --status pass\|fail\|blocked\|not_applicable --evidence-type --evidence-ref`. Append-only records. | PRD-AC-003–005 |
| `fid issue ac list` | List ACs with latest verification status per verifier role. | PRD-AC-006 |
| `fid issue complete` | Create append-only `task_completions` row. Attempt number, tokens, commits, files touched, gate results snapshot from `gate_runs`, optional human-intervention flag/note. Mark previous attempts superseded. | PRD-DEV-001–004, PRD-DEV-009–010 |
| `fid concern raise` | Create structured concern linked to issue and completion. `--on <issue-id> --by developer --in-completion <id> --summary --body-md`. | PRD-DEV-005 |
| `fid concern address` | Transition concern to addressed. `--in-completion <id>`. | PRD-DEV-006 |
| `fid concern list` | List concerns filtered by issue and status. | PRD-DEV-006–007 |
| Derived concern flag | `issues_enriched` view exposes `concerns_raised` and `open_concerns_count`. | PRD-DEV-007 |

### v0.1d — Claims, Permissions, Links

| Item | Description | PRD/TRD Refs |
|---|---|---|
| `fid issue claim` | Atomic acquire with TTL (default 30m). `--by <role> --ttl`, optional `--renew` for the current holder, optional `--force` for permitted force-acquire. Reject when live claim exists unless renewing, stale, or force-acquiring. | PRD-CLAIM-001–004, PRD-CLAIM-008a, TRD-ARCH-011 |
| `fid issue unclaim` | Release caller's current claim. Requires `--by <role>`. Force-acquire is handled by `fid issue claim --force`. | PRD-CLAIM-004 |
| `fid claim list` | List all live claims. `--stale` for past-TTL. `--by <role>`. | PRD-CLAIM-007 |
| Claim-transition integration | Transitions into started-category require live claim by caller. Transitions out release claim. | PRD-CLAIM-005–006 |
| Human claims | `--by human` signals agents should not advance. | PRD-CLAIM-008 |
| Permissions matrix | In-code constant for v0.1. Enforce `--by <role>` on every write. Reject unknown roles. Reject denied role/action combinations with actionable error. `fid policy show` displays the read-only in-code matrix in v0.1 so agents can inspect permitted actions before v0.2 configurability. | PRD-ROLE-001–005, TRD-ARCH-010 |
| `fid requirement sync` (minimal) | Parse canonical PRD/TRD requirement IDs from a source doc and upsert pointer rows (`id`, type, summary, source, hash) so v0.1 issue → requirement links have valid endpoints. No requirement get/list or drift report until v0.2. | PRD-GOV-006, PRD-LINK-002 |
| `fid link add/remove` | Create or remove typed links with auto-inverse projection. v0.1 scope: issue ↔ issue, issue → requirement. Validate relation/entity-type compatibility. Reject links that duplicate dedicated schema. | PRD-LINK-001–006a, 008 |
| `fid link list` | List all links from and to an entity. | PRD-LINK-006 |
| `fid label add/list` and `fid issue label add/remove` | Create labels, list labels, and add/remove issue label assignments without replacing the full label set. `fid issue update --labels` remains replace-all. `issue_labels` junction for multi-value. | PRD-ISSUE-015, PRD-ISSUE-023 |

### v0.1e — Telemetry, Composites, Maintenance

| Item | Description | PRD/TRD Refs |
|---|---|---|
| `fid run record` | Record `agent_runs` row: role, model, tokens in/out, snapshot pricing, computed cost_usd, and optional provider reconciliation fields. Warning when model not in pricing config. | PRD-TELEM-001–003a |
| `fid stats cost` | Aggregate cost by role, by model, by epic. | PRD-TELEM-004 |
| `fid stats tokens` | Aggregate tokens by role, by model. | PRD-TELEM-004 |
| `fid dev start` | Composite: claim + transition → in_progress + open agent_run. `--by developer`. | PRD-DEV-008, PRD-CLAIM-005 |
| `fid dev gate` | Record append-only `gate_runs` row on the derived active attempt before completion exists. | PRD-DEV-004, PRD-DEV-010 |
| `fid dev finish` | Composite: close agent_run + create task_completion with `gate_runs` snapshot + transition → dev_complete (validates rules). Supports optional human-intervention flag/note. `--by developer`. | PRD-DEV-008–010 |
| `fid dev block` | Transition → failed_gate with reason and create completion with `gate_runs` snapshot. | PRD-DEV-008, PRD-DEV-010 |
| `fid sql` | Read-only escape hatch. Reject non-SELECT/EXPLAIN/PRAGMA/WITH. | PRD-OUT-008, TRD-SEC-002 |
| `fid journal inspect` | Read-only journal event inspection by `event_id`. Shows actor, command metadata, effects, journal location, event hash, and revert eligibility. | PRD-OPS-013 |
| `fid maintenance check` | Consistency audit: kind-hierarchy violations, orphaned polymorphic refs, FK integrity, workflow state validity. Warnings: stale claims, started-without-claim, stale verifications, telemetry gaps. `--warn-as-error` for CI. `--format json\|markdown`. | PRD-OPS-005–008 |
| `fid maintenance vacuum` | SQLite VACUUM. | — |
| Activity log | Projected timeline from journal events. Read-model only. Basic `fid activity list` read with entity, event, and date filters. | PRD-OUT-017 |
| Pricing config | `.fid/pricing.toml` ships with defaults for current Anthropic, OpenAI, Google models. User-overridable. | PRD-NFR-010 |
| Stub table CLI guard | Commands targeting v0.2 entities (`decision create`, `escalation create`, etc.) return "commands for <entity> ship in <version>" with non-zero exit code. | PRD-ERR-012 |

**v0.1 exit criteria:**

- `fid init` creates a working project. `fid rebuild` reconstructs the DB from
  journal, including the roles seeded by `project.initialized`.
  Delete-and-rebuild is non-destructive.
- Full issue lifecycle: create → claim → transition → complete → unclaim works
  with workflow validation and permission enforcement.
- AC verification, concerns, and claim semantics all enforced.
- Ready queue returns correct next task for sequential and prioritised epics.
- Links with auto-inverse and relation validation work.
- JSON, TOON, and markdown output work for all read commands.
- Journal is append-only, committed to git. DB is gitignored and derived.
- `fid maintenance check` catches structural violations.
- Cost telemetry captures per-run pricing with computed USD.
- All v0.1 PRD requirements have passing tests.
- The author uses fid for real project work (dogfooding).

---

## v0.2 — Governance, Workflow Engine, Reports

**Goal:** Ship the configurable workflow engine, the full governance entity CLI
surface (decisions, escalations, triage, requirements), named reports, context
packs, and comments. This is where cross-epic analytical value compounds.

### v0.2a — Workflow Engine Migration

| Item | Description | PRD/TRD Refs |
|---|---|---|
| Workflow tables | Schema migration: `workflows`, `workflow_states`, `workflow_transitions`, `workflow_transition_requires`. Seed with v0.1 defaults verbatim. Migrate `issues.workflow_id` to real FK, `issues.status` to composite FK `(workflow_id, state)`. | PRD-WF-011–012, TRD-ARCH-009 |
| `fid workflow define/assign/list` | Define custom workflows with states, transitions, and per-transition rules. Assign per kind or per project. | PRD-WF-011–012 |
| Escalation workflow | Default workflow for escalation kind: open → resolved (requires linked_to_decision). | PRD-WF-005a |
| Role permissions table | Schema migration: `role_permissions` table. Seed from v0.1 code constant. Add `fid role create` and `fid policy set`; `fid policy show` keeps the same CLI contract but reads from `role_permissions` after this migration. | PRD-ROLE-006, TRD-ARCH-010 |

### v0.2b — Governance Entities

| Item | Description | PRD/TRD Refs |
|---|---|---|
| `fid decision create/get/list` | Decision records with summary, body_md_path, area, active/superseded status. `--by tl\|human`. | PRD-GOV-001 |
| `fid escalation create/get/list` | Escalation records with open/resolved status. `--from-finding`, `--from-concern`, `--trigger-class`. Resolution requires linked decision. | PRD-GOV-002–003, 005 |
| `fid triage record` | Triage events with RESUME/BLOCK verdict, source classification, reasoning. | PRD-GOV-004 |
| `fid triage start/verdict` | Composite TL-Lite triage commands. | PRD-GOV-012 |
| `fid tl wrapup` | Composite: capture wrap-up notes, close already-resolved escalations, and transition → wrapped. Retrospective row creation waits for v0.4 retrospective support. | PRD-GOV-013 |
| `fid requirement get/list` and sync drift reporting | Requirement pointer reads and filters. `fid requirement sync` already exists from v0.1; v0.2 adds drift reporting from `last_synced_hash` and `fid requirement list --drifted`. | PRD-GOV-006–007 |
| `fid comment add` | Comments on issues, decisions, escalations, requirements. `--on <type>/<id> --body-md --by <role>`. | PRD-GOV-008 |
| Links expansion | Link endpoints expand to `decision` and `escalation`. | PRD-LINK-007 |

### v0.2c — Reports and Context Packs

| Item | Description | PRD/TRD Refs |
|---|---|---|
| Named reports | `fid report epic-status`, `open-escalations`, `ac-matrix`, `epic-wrapup-pack`, `task-hygiene`, `stale-claims`. `fid report list`. `qa-blockers` waits for v0.3 QA findings. | PRD-OUT-006–007 |
| Attention report | `fid report attention`: pull-based human-judgment digest grouped by reason (open/stale escalations, epics awaiting approval, dead-end blocked, human claims). Audience-tagged; reasons light up as source entities ship. Read-only, no schema. | PRD-OUT-020 |
| Lightweight search | `fid search <term>` over structured projected fields with filters for kind, status, label, epic, inbox, and format. Literal/case-insensitive only; no custom DSL. | PRD-OUT-018 |
| Context packs | `fid issue get <id> --context-pack <role>`. Built-in packs: developer, qa_reviewer, qa_lead, tl_lite, tl. Role-tuned bundles. | PRD-OUT-005 |
| `fid changes` | Entity-grouped change digest: net field changes plus newly-added child facts since an anchor (event ID or timestamp), scoped by `--epic`/`--entity`. Aggregated from `activity_log`; distinct from chronological `fid activity list`. Extensible anchor resolver; role shaping deferred. | PRD-OUT-019 |
| Ready-queue explain | `fid ready --explain [--limit N] [--include-blocked]`: ranked candidates with deterministic `order_factors` + typed skip reasons (`blocked_by`/`claimed`/`sequence_wait`/`not_actionable`) + stale-claim flags. Default single-task output frozen; surfaces the existing rule, no new ranking semantics. | PRD-ISSUE-024 |
| Inbox promote primitive | `fid issue update --epic`/`--inbox`: assign/clear a task's epic parent (null↔epic; sequence required for sequential) and toggle the inbox flag, completing the capture→promote flow. Formal `fid intake` verbs deferred. | PRD-ISSUE-025 |
| `--dry-run` (all writes) | Universal write-preview flag: validate and build the would-be event, return the result plus an effect preview, append nothing. Exit code mirrors the simulated outcome. Relies on the v0.1 write-path seam (TRD §4.3) that separates validate/build from append. | PRD-OPS-015 |
| Run git context | Auto-captured read-only git context (branch, base/head commit, dirty, worktree) on `agent_runs`, with flag override and graceful degradation outside a repo. New `internal/gitctx/`. Additive migration to `agent_runs`. | PRD-OPS-017 |
| Run lifecycle + `fid run list/show` | `agent_runs.status` (open/closed/abandoned); `fid run list --open\|--abandoned` and `fid run show <run-id>`. | PRD-OPS-018 |
| `fid dev resume` / `fid dev abandon` | Resume re-acquires the claim, opens a new run linked via `resumed_from_run_id`, returns git resume hints (no git operations). Abandon releases the claim and marks the run abandoned, preserving the row. | PRD-OPS-019 |
| Stats expansion | `fid stats tasks --include est-vs-actual`, `fid stats attempts` (rework rate), `fid stats gates`, `fid stats interventions`. | PRD-TELEM-005–006 |
| `fid merge --journal` | Replay foreign journal with conflict detection using v0.1 event preconditions. Useful for recovered worktrees. | PRD-OPS-011 |
| `fid event revert` | Compensating-event correction: **inverse** mode restores mutable-state effects (transitions, field/label/link changes) under an effect precondition; **retraction** mode withdraws append-only facts (AC verifications, concerns, comments, triage) preserving the row. Creations/completions/claims/composites ineligible → forward domain guidance. Type-keyed, extensible. No bespoke domain-correction verbs. | PRD-OPS-014 |

**v0.2 exit criteria:**

- Custom workflows can be defined and assigned. v0.1 defaults migrate cleanly.
- Decisions, escalations, triage events, and requirements have full CLI surface.
- Comments work on all scoped entity types.
- Named reports answer common cross-cutting questions in single commands.
- `fid search` handles common discovery without requiring SQL.
- Context packs reduce agent context window pressure.
- Role permissions are data-configurable per project.
- All v0.2 PRD requirements have passing tests.

---

## v0.3 — QA Layer, Templates, TUI

**Goal:** Ship QA reviews and findings, QA plans, the template engine, TUI
kanban board, and CLI/skill-first agent workflow guidance. Enables the full agentic task loop
(TL plan → Developer execute → QA review → TL wrap-up) through fid.

### v0.3a — QA Reviews and Findings

| Item | Description | PRD/TRD Refs |
|---|---|---|
| QA tables | `qa_reviews`, `qa_findings`, `qa_review_task_window`. | PRD-QA-001–005 |
| `fid qa-review create/get/list` | Polymorphic scope: epic_window, single_task, epic_full. Record which tasks are covered for window scope. | PRD-QA-001–002 |
| `fid qa-finding add/list` | Findings with severity, status, task subject, evidence excerpt. Escalation from blocking findings. | PRD-QA-003–005 |
| `fid report qa-blockers` | Blocking QA findings report; ships with `qa_findings` so it has a first-class source table. | PRD-OUT-006 |
| QA review claims | Claims expand to `qa_reviews`. `fid qa-review claim/unclaim`. | PRD-CLAIM-009 |
| QA composites | `fid qa start`, `fid qa finding`, `fid qa verify`, `fid qa finish`. | PRD-QA-008 |
| Comment scope expansion | Comments on `qa_review` entities. | PRD-GOV-008 |

### v0.3b — QA Plans

| Item | Description | PRD/TRD Refs |
|---|---|---|
| QA plan tables | `qa_plans`, `qa_plan_steps`. `qa_reviews` gains nullable `plan_id` FK. | PRD-QA-006 |
| `fid qa-plan create/get/list` | QA plan per epic with optional contract reference (nullable until v0.4). | PRD-QA-006 |
| `fid qa-plan step add/list/link` | Ordered steps with evidence method and validator role. Display IDs `QA-NN`. Link to assertions deferred to v0.4. | PRD-QA-006 |
| `fid qa-plan finalise` | Coverage rule enforcement (no-op for contract rules until v0.4). | PRD-QA-007 |

### v0.3c — Templates

| Item | Description | PRD/TRD Refs |
|---|---|---|
| Templates table and CLI | `templates` table. `fid template create/list/show`. Monotonic versioning per (kind, name). Built-in kinds: task_spec, subtask_spec, qa_plan, qa_review, human_checklist, decision_record, completion_report, epic_brief, triage_log_entry, retro_lightweight. | PRD-OUT-012 |
| Template provenance | Add `created_from_template_{kind,name,version}` columns to `issues`, `task_completions`, `qa_reviews`, `decisions`. Populated when `--from-template` is used. | PRD-OUT-013 |
| Template-driven validation | Templates declare required sections. `fid dev finish` validates completion reports against the `completion_report` template (e.g., requires `## Concerns` heading). Rejects malformed artifacts. | PRD-OUT-014, TRD-ARCH-015 |

### v0.3d — TUI and Agent Skills

| Item | Description | PRD/TRD Refs |
|---|---|---|
| `fid board` | TUI kanban board. Columns grouped by status category. Cards show ID, title, priority, flags. Keyboard navigation, filtering by epic/labels/cycle. Bubbletea + Lipgloss. | PRD-OUT-015 |
| Agent skills | Fid-aware skills and examples for coding agents that use `fid --help`, JSON/TOON output, and shell commands. No MCP server, daemon, or network listener is committed in v0.3. | PRD-OPS-010, TRD-ARCH-013 |
| `fid handoff` | Handoff pack composition: `fid handoff --role <role> --issue <id>` composes context pack + workflow/permission-derived manifest (required reads, commands, stop conditions, expected writes, live blockers) + handoff template into a prompt and a machine-readable manifest. Role is a parameter; read-only. New `internal/handoff/`. | PRD-OUT-021 |
| **Conditional:** `fid query` DSL | Only if v0.2 search, named reports, and read-only SQL leave meaningful gaps. Grammar specced based on uncovered queries from real usage. | — |
| **Conditional:** `fid task-pack apply` | Bulk create-only ingestion of a reviewed structured task pack (epic + tasks + ACs + links) as one journal event with a source-ref→fid-id map. Only if v0.1–v0.2 dogfooding proves TL planning ingestion is real friction. Structured input only; not markdown import. Schema locked deliberately at build. | PRD-OPS-020 |

**v0.3 exit criteria:**

- Full QA review lifecycle: create → add findings → complete, with escalation.
- QA plans with ordered steps work for epic-level verification design.
- Templates with versioning and structural validation prevent malformed artifacts.
- TUI board provides useful kanban visualisation.
- Fid-aware agent skills make CLI-first workflows discoverable to Codex/Claude-style agents.
- The full agentic task loop runs through fid with enforced invariants.
- All v0.3 PRD requirements have passing tests.

---

## v0.4 — Validation Contracts, Retros, Export

**Goal:** Ship the validation contract layer, retrospectives, carry-forward
items, markdown export, saved views, and computed retrospective signals.
Validation contract outcomes become the primary wrap-up signal.

### v0.4a — Validation Contracts

| Item | Description | PRD/TRD Refs |
|---|---|---|
| Contract tables | `validation_contracts`, `validation_assertions`, `validation_assertion_validators`, `qa_plan_step_assertions`, `val_assertion_outcomes`. Add nullable `acceptance_criteria.val_assertion_id` FK. | PRD-VAL-001–008 |
| `fid contract create/get/activate` | Contracts per epic with draft/active/superseded status. | PRD-VAL-001 |
| `fid assertion add/get/list` | Assertions with `VAL-NNN` display IDs, statements, carry-forward source. | PRD-VAL-001 |
| `fid assertion validator add` | Required validator/evidence obligations per assertion. | PRD-VAL-002 |
| `fid assertion link-ac` | Link ACs to assertions for traceability. | PRD-VAL-006 |
| `fid qa-plan step link` | Link plan steps to assertions. | PRD-VAL-007 |
| `fid outcome record/list` | Append-only outcomes: pass/fail/carry_forward. Pass requires evidence. Carry_forward requires rationale. | PRD-VAL-003–005 |
| Coverage rule enforcement | At `fid qa-plan finalise`: every requirement has an assertion, every assertion has a validator and a plan step, every pass has evidence. | PRD-VAL-008, PRD-QA-007 |
| Contract reports | `fid report contract-coverage`, `assertion-trace`, `contract-outcomes`. | PRD-OUT-011 |
| `fid evidence check` | Read-only evidence-integrity check: resolves cited refs on AC verifications, findings, and validation outcomes (commits via local git, files via filesystem; command/manual/url unverifiable). `--warn-as-error`. Verifier not executor; no schema, no registry. New `internal/evidence/`. | PRD-QA-009 |

### v0.4b — Retrospectives and Carry-Forward

| Item | Description | PRD/TRD Refs |
|---|---|---|
| `carry_forward` table and CLI | Status: open, scheduled, closed, principled. `scheduled_epic_id` for future work. `promoted_principle` flag for principles rendering. | PRD-GOV-009–010 |
| `retrospectives` table and CLI | Per-epic, lightweight or deep type. Recommendation field. `notes_md_path`. | PRD-GOV-011 |
| Computed retro signals | Token-budget overruns, finding density, AC verification gap, contract pass-rate, carry-forward count. | spec §16 v0.4 |
| Recurring-pattern reports | `fid report recurring-findings`/`recurring-concerns`/`recurring-patterns`: deterministic grouping of findings/concerns/gates/triage/VAL failures by structured dimension, ranked as principle candidates. Surfaces candidates only — no ML, no fid-authored prose. | PRD-OUT-023 |

### v0.4c — Export and Views

| Item | Description | PRD/TRD Refs |
|---|---|---|
| `fid export` | One-way markdown export. Default `whet` layout (docs/current/, docs/epics/<id>/tasks/, docs/epics/<id>/qa/, docs/reviews/). Frontmatter contract with YAML structured fields. "Generated by fid" banner. Config: `.fid/config.toml` `[export]` section. `--to`, `--epic`, `--since` flags. | PRD-OUT-009–010, TRD-ARCH-014 |
| `views` table and CLI | Saved read-only queries. Start with SQL-backed views; any DSL-backed view kind requires the conditional query decision to be accepted first. `fid view create/run/list`. | PRD-OUT-016 |

**v0.4 exit criteria:**

- Validation contracts with assertions, required validators, and outcomes provide
  full traceability from requirements through to evidence.
- Coverage rules enforce that every assertion has evidence.
- Contract-coverage, assertion-trace, and contract-outcomes reports work.
- Retrospectives and carry-forward items support epic wrap-up discipline.
- Markdown export produces a human-readable project tree with frontmatter.
- All v0.4 PRD requirements have passing tests.

---

## v0.5 — Polish (Tentative)

**Goal:** Improve the user experience with a web dashboard, cycles support,
and conditional features based on usage evidence.

| Item | Description | PRD/TRD Refs |
|---|---|---|
| Web dashboard | Separate `fid-dashboard` package. Reads the same SQLite DB. Browser-based project visualisation. | — |
| Cycles support | `cycles` table opt-in. `fid cycle create/list`. `due_date` UI in board. | PRD-ISSUE-017 |
| File-scope claims | `expected_files_touched` enforcement on task completions. | spec §6.12 |
| **Conditional:** MCP adapter evaluation | Only if dogfooding identifies a concrete MCP-only client need. Start from stdio as the default candidate; any HTTP/listener transport requires a separate threat model and product/security decision. | PRD-OPS-010, TRD-ARCH-013 |
| **Conditional:** Idempotency keys | Caller-supplied keys for retry-safe programmatic writes; a repeated key echoes the original event and result. Only with a concrete programmatic/MCP retry need. Additive to the write-path seam; pairs with the MCP evaluation. | PRD-OPS-016 |
| **Conditional:** `fid context estimate` | Estimate-only context-pack/handoff token metering, by section, as a model-agnostic heuristic. Only if context packs prove heavy enough to overload sessions. No auto-trim; no per-model tokenizer. | PRD-OUT-022 |
| **Conditional:** `fid import` + `fid sync` | Bidirectional markdown sync + MD-canonical mode. Only if v0.4 export usage generates evidence of a real use case. Alternative: LLM agent with fid skill. | — |
| **Conditional:** `fid init --import` | Import from arbitrary whet markdown tree. Depends on import/sync. | — |

**v0.5 exit criteria:**

- Dashboard reads and renders project state from SQLite.
- Conditional features shipped only with evidence of demand.

---

## Later Phases

| Phase | Theme | Notes |
|---|---|---|
| v0.6+ | Agent skill library | Built-in agent skills for Claude Code / Codex that encode fid-aware workflows |
| v0.6+ | Multi-project | Cross-database queries, project federation | Out of scope for v1 |
| v0.6+ | Performance tier | Checkpoint-based partial rebuild, journal compaction | Forward-compatible optimisations |

---

## Risk Register

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| Schema migrations break journal replay | Medium | High | Event schema versioning with pure upgrader functions (TRD-ARCH-003). Golden-file replay tests per version boundary. |
| Clean textual git merge creates logical journal conflict | Medium | High | v1 journal events include targeted preconditions. Replay stops with structured conflict details; v0.2 `fid merge --journal` builds on the same data. |
| Pure-Go SQLite driver performance at scale | Low | Medium | Benchmarks at 10k issues / 100k events (TRD-PERF-001–003). Fallback to `mattn/go-sqlite3` if perf is real issue. |
| TOON spec drift breaks the in-tree marshaller | Low | Low | Subset implementation. JSON is canonical; TOON is secondary. |
| Journal merge conflicts on high-traffic branches | Medium | Medium | Same-checkout writes are serialized by `.fid/write.lock`; cross-worktree JSONL line-merge is textual and resolvable. `fid rebuild` reconciles after merge. |
| Trigger backstop misses a corruption path | Low | High | `fid maintenance check` as periodic audit (PRD-OPS-005). Trigger + app-layer defense in depth. |
| Scope creep from v0.1 into later version features | Medium | Medium | Stub tables with "ships in vN" CLI guard. PRD traceability table enforces version boundaries. |
| Agent implementation quality degrades for under-designed commands | High | High | Design Plan output examples for every v0.1 command (audit recommendation #2). |
| Premature MCP adapter duplicates the CLI and weakens no-listener guarantees | Medium | Medium | Keep v0.3 CLI/skill-first. Evaluate MCP only after dogfooding proves a concrete MCP-only client need and a threat model exists. |

---

## Dependency Map

```text
v0.1a Persistence Foundation
  └──► v0.1b Issue Lifecycle
        └──► v0.1c ACs, Completions, Concerns
              └──► v0.1d Claims, Permissions, Links
                    └──► v0.1e Telemetry, Composites, Maintenance

v0.1 ──► v0.2a Workflow Engine Migration
           └──► v0.2b Governance Entities
                 └──► v0.2c Reports and Context Packs

v0.2 ──► v0.3a QA Reviews ──► v0.3b QA Plans
         v0.3c Templates (parallel with QA)
         v0.3d TUI + Agent Skills (parallel, after QA tables exist)

v0.3 ──► v0.4a Validation Contracts
         v0.4b Retros + Carry-Forward (parallel)
         v0.4c Export + Views (parallel, after contracts)

v0.4 ──► v0.5 Polish (conditional items)
```

*Last updated: 2026-06-01*
