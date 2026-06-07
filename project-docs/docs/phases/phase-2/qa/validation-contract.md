# Phase 2 Validation Contract - Widen The Schema Vertically

> Phase: 2 - Widen The Schema Vertically (full v0.1 schema, triggers, enriched view)
> Author: Tech Lead (draft)
> Last refined: <pending QA Lead refinement>
>
> Lives at `project-docs/docs/phases/phase-2/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 2 is complete. The validation contract specifies **what must be proven**;
the QA plan (authored next by the QA Lead) is the strategy for **how to prove
it**.

Phase 2 completes the v0.1 base SQLite schema *through the journal/projection
machinery proven in Phase 1*. The binding rule from the roadmap slicing note is
that each new table/trigger group is exercised by appending a journal event and
reading the projection (or, for triggers, by attempting a backstop-violating raw
write), **not by direct-SQL happy-path tests only**.

If a required validation assertion cannot be satisfied, stop and escalate
instead of weakening the assertion or silently adapting the spec.

## Numbering

- IDs use `VAL-NNN`, scoped to this phase, sequential from `VAL-001`.
- Retired or superseded assertions stay in the table with status and rationale;
  do not reuse IDs.

## Coverage-Evidence Carry-In (from Phase 1 DEC-007)

Where a `VAL-NNN` behavior lives in a COV-001 domain package
(`internal/issues`, `internal/journal`, `internal/policy`; thresholds in
`08_QUALITY_PLAN.md` §4.3), its evidence method **names package-level tests in
that package**, not only `cmd/fid` integration evidence. This avoids the
Phase-1 planning seam where domain packages read 0% own-package coverage.
`internal/db` is not a COV-001 package, so trigger/migration tests there do not
gate COV-001 but must still carry deterministic evidence.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | Roadmap v0.1a "SQLite schema and migrations"; TRD §4.1; TRD SCHEMA §1.2, §2, §5, §7, §9, §10, §11 | Migration `0002` up creates every remaining v0.1 base table plus the v0.1 stub tables with their specified columns and static CHECK constraints | Integration: migrate up on a temp DB; introspect `sqlite_master`/`PRAGMA table_info`; assert each table + representative columns + static CHECKs (`internal/db` test) | QA Lead | task-1 | `epic_metadata`, `acceptance_criteria`, `ac_verifications`, `gate_runs`, `task_completions`, `concerns`, `requirements`, `agent_runs`, `labels`, `issue_labels`, `links`, `claims`, `activity_log`, and stubs `decisions`, `escalations`, `triage_events` all present with specified columns and static CHECKs | Any required table, column, or static CHECK missing or misnamed | Draft | Workflow/role_permissions/comments (v0.2), QA/templates (v0.3), validation/carry_forward/retros/cycles/views (v0.4+) are excluded. |
| VAL-002 | TRD §5 (schema versioned via numbered migrations; up/down tests) | Migration `0002` down removes exactly the tables it created; up→down→up round-trips cleanly without touching Phase-1 tables | Integration: migrate up, down, up; assert post-down schema equals the Phase-1 schema and re-up restores the full schema (`internal/db` test) | QA Lead | task-1 | Down drops the new tables and leaves `projection_meta`/`roles`/`issues` intact; re-up recreates all new tables with no error | Down errors, leaves orphan tables, drops Phase-1 tables, or re-up fails | Draft | Stub-table FK ordering on down must not error even though no CLI populates them yet. |
| VAL-003 | Roadmap Phase 2 candidate ("exercised by appending an event and reading the projection, not by direct SQL only"); TRD-ARCH-001, TRD-ARCH-002 | Each new v0.1 table is reachable through the journal/projection engine: events whose insert effects target representative new tables replay and read back equal to the effect values | Package-level `internal/journal` test: build event(s), `ProjectEvent`/replay into temp DB, read rows; assert equality; FK-dependent rows (issue→completion→concern) ordered across effects/events | QA Lead | task-1 | `requirements`, `labels`, `epic_metadata`, `acceptance_criteria`, `task_completions`, and `concerns` are each populated via journal effects and read back equal; population uses journal effects, not direct SQL INSERT | Any representative table cannot be populated through journal effects, or a read-back differs from the effect values | Draft | Exhaustive per-column population is not required; representative coverage of each table group is. |
| VAL-004 | TRD §4.3 (projection_schema_version mismatch → rebuild); DEC-009; Phase-1 VAL-012 schema-version class | Growing the schema bumps `db.ProjectionSchemaVersion`; a checkpoint at the prior version is stale (full rebuild before serving); a checkpoint at the current version matching the journal serves without rebuild | `internal/db` test of `DecideProjectionFreshness` plus an `internal/journal` end-to-end rebuild triggered by a stored prior version | QA Lead | task-1 | Stored prior `projection_schema_version` → `ProjectionRebuild`; current version + matching checkpoint → `ProjectionFresh` | Prior-version checkpoint served without rebuild, or current matching checkpoint forces a needless rebuild | Draft | `ProjectionSchemaVersion` value is an internal constant, not a public interface. |
| VAL-005 | TRD SCHEMA §2.2; TRD §4.1 (`issues_enriched`); PRD-DEV-007 | `issues_enriched.open_concerns_count` equals the count of `status='open'` concerns for the issue and `concerns_raised` equals `(count > 0)` | Integration: project an issue plus open and addressed concerns via journal; query the view for both states | QA Lead | task-2 | Issue with N open concerns → `open_concerns_count=N`, `concerns_raised=true`; issue with zero open (or only addressed) concerns → `0`, `false` | Count includes non-open concerns, miscounts, or `concerns_raised` disagrees with the count | Draft | Only `status='open'` counts; addressed/escalated/wontfix excluded. |
| VAL-006 | TRD SCHEMA §2.2 | `issues_enriched.triage_status` is derived from the latest `triage_events` verdict for the issue, mapped `RESUME→'resumed'`, `BLOCK→'blocked'`, none→`NULL` | Integration: project multiple `triage_events` at distinct `created_at`; assert latest wins; cover all three states | QA Lead | task-2 | Latest `RESUME`→`'resumed'`, latest `BLOCK`→`'blocked'`, no triage→`NULL` | Not the latest row, wrong mapping, or non-null when no triage exists | Draft | "Latest" is by `created_at` per the view definition. |
| VAL-007 | Design CLI `fid issue get`; 05_DESIGN_PLAN §6.1; PRD-ISSUE-006/018; PRD-OUT-001/004 | `fid issue get [--format json]` returns `open_concerns_count`/`concerns_raised`/`triage_status` sourced from `issues_enriched`, and the values survive delete-DB + rebuild | Package-level `internal/issues` tests (COV-001 carry-in) **plus** `cmd/fid` CLI/integration evidence; both states | QA Lead | task-2 | Created issue with appended open concern(s)+triage reports matching `N`/`true`/`<status>`; an issue with none reports `0`/`false`/`null`; identical after `fid rebuild` | `fid issue get` still returns hard-coded `0`/`false`/`null`, or values diverge from the view or after rebuild | Draft | Other `--format` renderers (TOON/markdown/YAML) remain Phase 3. |
| VAL-008 | TRD-ARCH-008; TRD-PERF-004 | The `issues_enriched` concern/triage joins use `idx_concerns_issue_status` and `idx_triage_issue_created` (index lookups, not full scans of `concerns`/`triage_events`) | `EXPLAIN QUERY PLAN` assertion in a test showing index search for both joins | QA Lead | task-2 | Query plan shows `SEARCH ... USING INDEX` on both `concerns(issue_id, status)` and `triage_events(issue_id, created_at)` | Plan shows a full table `SCAN` on `concerns` or `triage_events` | Draft | Plan strings may vary by SQLite version; assert index use, not exact text. |
| VAL-009 | TRD SCHEMA §14.1 (parent_id); TRD §4.5; TRD-ARCH-006; TRD-SEC-004 | A BEFORE INSERT/UPDATE backstop trigger on `issues` rejects raw-SQL writes that violate kind hierarchy; valid hierarchies succeed | Integration: raw-SQL inserts for each violation RAISE; valid rows insert; both states | QA Lead | task-3 | Epic-with-parent, subtask-with-non-task/absent-parent, and task-with-non-epic-parent each RAISE; valid epic(no parent)/task(epic parent or none)/subtask(task parent) insert | Any invalid hierarchy inserts, or any valid hierarchy is rejected | Draft | App-layer hierarchy validation ships with `fid issue create` lifecycle (v0.1b); this is the SQL backstop. |
| VAL-010 | TRD SCHEMA §14.1 (`acceptance_criteria`, `task_completions`, `epic_metadata`) | Backstop kind-guard triggers reject raw-SQL inserts of `acceptance_criteria`/`task_completions` whose `issue_id` is not a `task` or `subtask`, and `epic_metadata` whose `issue_id` is not an `epic`; valid kinds insert | Integration: raw-SQL both states for each table | QA Lead | task-3 | Wrong-kind inserts RAISE; `acceptance_criteria`/`task_completions` on task/subtask and `epic_metadata` on epic insert | Wrong-kind insert succeeds or a valid-kind insert is rejected | Draft | `task_completions` guard follows §14.1 ("same" as AC = task or subtask); see planning note on §4.1's indicative "(task only)" wording. |
| VAL-011 | TRD SCHEMA §14.1 (`triage_events.epic_id`, `decisions.epic_id`) | Backstop kind-guard triggers reject raw-SQL inserts of `triage_events`/`decisions` whose `epic_id` is set but does not reference an `epic`; valid references (and NULL where allowed) insert | Integration: raw-SQL both states | QA Lead | task-3 | Non-epic `epic_id` RAISES on `triage_events` (NOT NULL) and on `decisions` when set; valid epic ref (and NULL `decisions.epic_id`) insert | Non-epic `epic_id` accepted or a valid reference rejected | Draft | These are stub tables (no CLI in v0.1); trigger is a pure backstop. |
| VAL-012 | TRD SCHEMA §14.7 | A BEFORE UPDATE OF kind trigger on `issues` rejects any change to `kind`; updates to other columns succeed | Integration: raw-SQL `UPDATE issues SET kind=...` RAISES; `UPDATE` of another column succeeds; both states | QA Lead | task-3 | Kind change RAISES; non-kind column update succeeds | Kind change accepted or non-kind update wrongly rejected | Draft | — |
| VAL-013 | TRD SCHEMA §14.2; TRD-SEC-004 | BEFORE INSERT polymorphic ref-integrity triggers on `links`, `claims`, and `activity_log` reject rows whose `(ref_type, ref_id)` — and link `from`/`to` endpoints — do not resolve to an existing row; valid refs insert | Integration: raw-SQL both states for each table (links: both endpoints) | QA Lead | task-4 | Dangling ref on each table RAISES; rows whose endpoints resolve insert | A dangling ref inserts or a resolvable ref is rejected | Draft | The §14.6 "activity_log inserts only during projection" guard is deferred (see Scope Deferrals). |
| VAL-014 | TRD SCHEMA §14.3, §10.1 | The `links` backstop trigger accepts only v0.1 endpoint types (`issue`, `requirement`) and relation/endpoint combinations per §10.1, rejecting `decision`/`escalation` endpoints and incompatible relation/type combos | Integration: raw-SQL covering accept and reject cases | QA Lead | task-4 | Rejects a `decision`/`escalation` endpoint and an incompatible combo (e.g. `implements` issue→issue); accepts `implements` issue→requirement, `blocks` issue→issue, `relates_to` issue↔requirement | An out-of-scope endpoint or incompatible relation/type combo inserts, or a valid v0.1 link is rejected | Draft | `fid link add` with auto-inverse is app-layer (v0.1d); only the backstop trigger is in scope here. |
| VAL-015 | TRD SCHEMA §14.4 | The `claims` backstop trigger accepts only `ref_type='issue'` for v0.1 and rejects unsupported `ref_type` values | Integration: raw-SQL both states | QA Lead | task-4 | `ref_type='issue'` (resolvable) inserts; `ref_type='qa_review'` (or other) RAISES | Unsupported `ref_type` inserts or a valid issue claim is rejected | Draft | `qa_review` claim scope ships in v0.3. |
| VAL-016 | Roadmap Phase 2 exit ("rebuild stays deterministic"); TRD-ARCH-001; PRD-OPS-003/004 | `fid rebuild` over a journal that exercises every new v0.1 table is deterministic and non-destructive, and replays cleanly with all triggers active (no false rejection of valid journal-applied rows) | Integration: build a journal touching all new tables; rebuild; compare projection state pre/post delete-DB; assert no trigger false-positive during replay | QA Lead | phase-end | Delete-DB + `fid rebuild` yields identical projection state (row/structural equality), and replay of valid effects raises no trigger error | Rebuild diverges, errors, or a trigger rejects a valid journal-applied row | Draft | Performance/benchmark tuning of rebuild remains Phase 3. |
| VAL-017 | TRD-ARCH-005; Phase-1 VAL-017 | The adapter boundary holds: no `internal/...` package imports `cmd/fid`; the enriched-read change lives in `internal/`, and `cmd/fid` stays a thin adapter | `go list` import-graph check (no `cmd/fid` in `internal/...` imports) plus a `grep` adapter-boundary source scan | QA Lead | phase-end | Import graph shows no `internal/...`→`cmd/fid` edge; `cmd/fid` delegates enriched read to `internal/issues` | Any `internal/...` imports `cmd/fid`, or business logic leaks into `cmd/fid` | Draft | Standing architectural invariant re-asserted each phase. |

## Scope Deferrals (not Phase 2 assertions)

These are intentionally excluded from Phase 2; record as carry-forward at wrap
if still open:

- **Sequence-vs-scheduling backstop trigger** (§14.1 `issues.sequence` bullet):
  the cross-table rule "tasks under sequential epics must have a sequence; tasks
  under prioritised/no-epic must not" is coupled to v0.1b scheduling semantics
  (`epic_metadata` population, tasks-under-epics, ready queue) and has an
  ambiguous missing-metadata case. The static positive-`sequence` CHECK and
  `UNIQUE(parent_id, sequence)` already ship (Phase 1). Defer the cross-table
  backstop to v0.1b when the semantics are exercised end-to-end.
- **§14.6 activity_log "inserts only during projection" guard:** Phase 2 ships
  the §14.2 polymorphic ref backstop on `activity_log`. A SQL-level
  "projection-only" guard needs a connection-scoped session flag, which is
  unreliable under the pooled multi-connection model from DEC-006, and there is
  no `activity_log` population path until v0.1b. Keep projection-only discipline
  app-layer for now; defer the trigger-level guard.
- **The `activity` effect op** (TRD §4.3 op list): `activity_log` is exercised
  via plain `insert` effects in Phase 2. The dedicated `activity` op ships with
  the lifecycle commands that actually emit activity entries (v0.1b+).

## Coverage Rules Check (Tech Lead draft)

- Every in-scope roadmap/TRD requirement is referenced by at least one
  `VAL-NNN`, or recorded under Scope Deferrals: YES.
- Every `VAL-NNN` is owned by a task or `phase-end`: YES (tasks 1-4 + phase-end).
- Conditional/binary assertions name both states: YES (VAL-002, VAL-004..015).
- COV-001 carry-in applied: VAL-003 names `internal/journal`; VAL-007 names
  `internal/issues` package-level tests.
- QA Lead still to: turn each evidence method into deterministic QA-plan steps,
  confirm trigger backstop tests assert RAISE (not silent no-op), and map every
  `VAL-NNN` to a `(qa-lead)`, mid-phase-review, or `(human)` step.

## Outcomes

Phase-end final review records every `VAL-NNN` as `PASS`, `FAIL`, `BLOCKED`, or
`CARRY-FORWARD`. `PASS` without cited evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-07 | VAL-001..017 | created (Tech Lead draft) | Phase 2 planning | N/A |
