# Phase Brief - Phase 2: Widen The Schema Vertically

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** PLANNED (drafted 2026-06-07). Validation contract + ordered task
> specs drafted; awaiting QA Lead QA plan and human approval before execution.
> Phase 1 is CLOSED (see `../phases/phase-1/`).

## Phase Goal

Complete the **v0.1 base SQLite schema** — the remaining v0.1 tables, the stub
tables, the v0.1 triggers, and the `issues_enriched` view — *through the same
journal → projection → rebuild machinery proven in Phase 1* (Roadmap "Phase 2 —
Widen the schema vertically"). Each new table/trigger group is exercised by
appending a journal event and reading the projection (or, for triggers, by a
backstop-violating raw write), **not by direct-SQL happy-path tests only**.

**Headline end-to-end increment:** `fid issue get` stops returning hard-coded
zeros and surfaces real derived `concerns_raised` / `open_concerns_count` /
`triage_status` from `issues_enriched`.

## Scope (in)

- Migration `0002`: remaining v0.1 base tables (`epic_metadata`,
  `acceptance_criteria`, `ac_verifications`, `gate_runs`, `task_completions`,
  `concerns`, `requirements`, `agent_runs`, `labels`, `issue_labels`, `links`,
  `claims`, `activity_log`) + stub tables (`decisions`, `escalations`,
  `triage_events`). DDL only; verified through journal effects.
- Migration `0003`: `issues_enriched` view + `idx_concerns_issue_status` +
  `idx_triage_issue_created`; wire `fid issue get` to read the view.
- Migration `0004`: kind-invariant + kind-immutability backstop triggers.
- Migration `0005`: polymorphic ref integrity + link scope/relation + claim
  scope backstop triggers.
- `ProjectionSchemaVersion` bump so existing derived DBs rebuild on upgrade.

## Scope (out — deferred)

- v0.2+ tables (workflows, role_permissions, comments), v0.3 QA/templates,
  v0.4+ validation/carry_forward/retros/cycles/views.
- Lifecycle commands (`issue create` kinds hierarchy, transitions, claims, ACs,
  list filters, `--has-concerns`) → v0.1b+.
- TOON/markdown/YAML renderers, `--section`, perf tuning → Phase 3 / v0.1b.
- **Sequence-vs-scheduling backstop trigger** (§14.1 `issues.sequence`) →
  v0.1b (coupled to scheduling semantics; static positive-`sequence` CHECK +
  `UNIQUE(parent_id, sequence)` already ship).
- **§14.6 activity_log "projection-only" guard trigger** and the **`activity`
  effect op** → v0.1b (no population path yet; pooled-connection model makes a
  SQL session guard unreliable). `activity_log` is exercised via `insert`
  effects + the §14.2 polymorphic ref backstop now.

## Validation Contract Index

Contract: `../phases/phase-2/qa/validation-contract.md` (VAL-001..VAL-017).

| VAL | Behavior | Owner |
|---|---|---|
| VAL-001 | Migration 0002 creates all remaining v0.1 + stub tables | task 1 |
| VAL-002 | Migration 0002 down clean; up→down→up round-trips | task 1 |
| VAL-003 | New tables reachable via journal effects (not direct SQL) | task 1 |
| VAL-004 | `ProjectionSchemaVersion` bump forces stale-DB rebuild | task 1 |
| VAL-005 | `issues_enriched` open-concern count + `concerns_raised` | task 2 |
| VAL-006 | `issues_enriched.triage_status` from latest verdict | task 2 |
| VAL-007 | `fid issue get` returns enriched cols; survives rebuild | task 2 |
| VAL-008 | Enriched joins use concern/triage indexes (no scan) | task 2 |
| VAL-009 | `issues.parent_id` kind-hierarchy backstop trigger | task 3 |
| VAL-010 | Kind guards: ac/completion (task\|subtask), epic_metadata (epic) | task 3 |
| VAL-011 | `epic_id` kind guards: triage_events, decisions | task 3 |
| VAL-012 | Kind immutability (BEFORE UPDATE OF kind) | task 3 |
| VAL-013 | Polymorphic ref integrity: links, claims, activity_log | task 4 |
| VAL-014 | Link endpoint scope + relation compatibility (§10.1) | task 4 |
| VAL-015 | Claim scope (`ref_type='issue'` only) | task 4 |
| VAL-016 | Deterministic full-schema rebuild with triggers active | phase-end |
| VAL-017 | `internal/` never imports `cmd/fid/`; thin adapters | phase-end |

## Ordered Tasks

1. **Complete the v0.1 base schema (migration 0002)** — remaining v0.1 + stub
   tables; exercise via journal effects; bump projection schema version.
   (VAL-001..004)
2. **`issues_enriched` view + indexes + enriched `fid issue get` (0003)** —
   headline end-to-end increment. (VAL-005..008)
3. **Kind-invariant + kind-immutability backstop triggers (0004)** — raw-SQL
   backstop verification. (VAL-009..012)
4. **Polymorphic ref + link/claim scope triggers (0005) + phase-end** —
   determinism + adapter boundary. (VAL-013..017)

## Exit Criteria

- Full v0.1 base schema present; migrations up/down round-trip cleanly.
- `issues_enriched` returns concern count + triage status; `fid issue get`
  surfaces them, identical after rebuild.
- Triggers enforce kind hierarchy, kind immutability, polymorphic ref integrity,
  link scope, and claim scope (raw-SQL RAISE evidence).
- `fid rebuild` over a journal touching every new table is deterministic and
  non-destructive with all triggers active.
- All VAL-001..VAL-017 PASS with deterministic evidence; QA final review has no
  unresolved blocking findings; human checklist PASS.

## Planning Notes / Carry-In

- Cut a new phase branch (e.g. `phase-2/schema-widening`).
- Apply Phase-1 DEC-007 carry-in: VALs whose behavior lives in COV-001 domain
  packages name **package-level** tests — VAL-003 (`internal/journal`), VAL-007
  (`internal/issues`). Hold COV-001 ≥ 80% (issues/journal) and COV-002 ≥ 70%
  (cmd/fid).
- The generic projection engine already inserts into arbitrary tables, so the
  "exercise via journal" requirement needs no engine change — only well-formed
  events.
- **Spec ambiguity to resolve at execution:** `task_completions` kind guard
  follows §14.1 (task **or** subtask); §4.1's "(task only)" is the indicative
  entity note. Record the resolution in the Task 3 completion report; revisit at
  spec sync.
- Carry-Forward Register: CF-003 (role-validation corrupt-tail tolerance)
  remains open but is not exercised by Phase 2. Two new deferrals
  (sequence-scheduling trigger; §14.6 guard / `activity` op) to be formalized as
  carry-forward at wrap if still open.

## Next Checkpoint

QA Lead writes `../phases/phase-2/qa/qa-plan.md` mapping every VAL to
deterministic verification; Tech Lead presents contract + tasks to the human;
Developer execution begins only after the human approves.
