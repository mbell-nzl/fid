# Phase 8 Validation Contract - Task Completions

> Phase: 8 - Task completions (v0.1c)
> Author: Tech Lead (draft)
> Last refined: 2026-06-18 by Tech Lead (awaiting QA Lead evidence refinement)
>
> Location: `project-docs/docs/phases/phase-8/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 8 is complete. The contract specifies **what must be proven**; the QA
plan specifies **how**.

If a required assertion cannot be satisfied, stop and escalate rather than
weakening it.

## Phase scoping notes

Phase 8 ships the **task-completion write surface** (`fid issue complete`) on top
of the existing `task_completions` / `gate_runs` schema (v0.1a schema-widening,
Phase 2) and read path (`internal/issues/detail.go` `issueCompletions`). It
replaces the Phase 5 journal-fixture evidence for the `has_completion_report`
transition rule with real write commands. The generic journal effect machinery
already inserts/updates arbitrary tables, so no projection allowlist change is
expected. Decisions confirmed by the human (2026-06-20):

- **DEC-A (confirmed) — new-attempt completion status.** A standalone
  `fid issue complete` records a new attempt with `status = 'dev_complete'`;
  prior attempts on the same issue move to `status = 'superseded'`. Gate
  pass/fail is recorded independently via `gates_passed`/`gates_failed`; the
  `failed_gate` status and the paired workflow transition are the job of the
  v0.1e composites (`fid dev finish`/`fid dev block`), not this primitive.
- **DEC-B (confirmed) — attempt numbering.** `attempt_number` is derived as
  `MAX(attempt_number) + 1` over the issue's `task_completions` (1 when none),
  computed under `.fid/write.lock`. Deterministic; not user-supplied.
- **DEC-C (confirmed) — gate snapshot is fixture-seeded in Phase 8.**
  `fid dev gate` (the `gate_runs` writer) ships in v0.1e, so the gate-snapshot
  mechanism (VAL-005) is exercised by seeding `gate_runs` through a journal
  fixture (Phase 5 / Phase 7 forward-wrinkle pattern). In real Phase 8 usage no
  gate writer exists yet, so a real `fid issue complete` snapshots an empty set
  (`gate_results = []`, `gates_passed = gates_failed = 0`) — which is correct.
  Real-write linkage is re-verified in v0.1e.
- **DEC-D (confirmed) — no `activity_log` for completions (DEC-020 precedent).**
  Completions are append-only fact creations and follow the `fid issue create` /
  AC-write precedent: no `activity_log` row. Consistent timeline population is
  deferred to v0.1e via **CF-011**. No activity assertion in this contract by
  design.
- **DEC-E (confirmed) — completions reject on epics (CF-007).** `fid issue
  complete` targets tasks/subtasks only; an epic target is rejected **before**
  any journal append (non-zero exit, clear error, zero new event/row). The
  `task_completions_issue_kind_guard_insert` trigger is a backstop; the CLI must
  reject first. Permissions are `RequireKnownRole` only (matrix is v0.1d). Epics
  never need a completion: epic closure runs through the epic workflow
  (`in_progress→wrapped` gated by `all_child_tasks_done_or_cancelled`, then
  `wrapped→closed`), not `task_completions`, so the rejection cannot strand an
  epic. The epic `wrapped`/`closed` path is out of scope for Phase 8.

CF-012 (`internal/issues/issue.go` 677-line behavior-preserving split) is folded
into this phase as Task 4 and is covered by VAL-009.

## Numbering

- IDs use `VAL-NNN`, scoped to Phase 8, sequential from `VAL-001`.
- Carry-forward assertions are re-issued with a new `VAL-NNN` and a source note.
- Retired/superseded assertions stay in the table; IDs are never reused.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | PRD-DEV-001; design §2 `fid issue complete` | `fid issue complete <task-or-subtask> --by` creates a `TC-<hash>` `task_completions` row storing `developer_role`, `actual_tokens`, `commit_hashes` (JSON array), `files_touched` (JSON array), and `notes_md_path`; read-back via `fid issue get` (and JSON) shows the completion with those fields. | CLI: complete a task with `--tokens`/`--commits`/`--files`/`--notes-md`; read back via `fid issue get --format json`; assert `TC-` ID shape and each stored field. | QA Lead | task-1 | Completion row persists with `TC-` ID and all named fields stored/read-back correctly. | A named field is dropped/mis-stored, wrong ID shape, or no row persisted. | Draft | `expected_files_touched` (v0.5 file-scope claims) and `model` capture beyond `--by` are out of scope. |
| VAL-002 | PRD-DEV-002; PRD-DEV-003 | The first completion on a task is `attempt_number = 1`; a second completion on the same task is `attempt_number = 2` and the prior attempt's `status` becomes `superseded` while its row is **retained** (append-only). | CLI: complete the same task twice; assert attempt numbers 1 then 2, prior row status `superseded`, both rows still present. | QA Lead | task-1 | Both attempt numbers correct; prior attempt retained with `superseded` status; new attempt active. | Attempt not incremented, prior row deleted/overwritten, or prior status not set to `superseded`. | Draft | Per-attempt AC re-verification windowing ("since last completion") is not implemented here. |
| VAL-003 | PRD-DEV-009; design §2 | `--human-intervention` sets `human_intervention = 1` and `--intervention-note` stores the note; omitting the flag yields `human_intervention = 0` and no note. | CLI both states: (a) complete with `--human-intervention --intervention-note`; (b) complete without; assert flag and note in each. | QA Lead | task-1 | Flagged case stores `human_intervention=1` + note; unflagged case stores `0` + no note. | Flag ignored, note dropped, or unflagged case records an intervention. | Draft | Intervention analytics (`fid stats interventions`) is v0.2c. |
| VAL-004 | PRD-DEV-001; CF-007; DEC-E | `fid issue complete` on an **epic** is rejected before any journal append (non-zero exit, clear error naming the kind constraint, zero new event/row); the same command on a **task/subtask** succeeds. | CLI over two targets: (a) epic → non-zero exit + error + zero new journal lines; (b) task → success + one new completion. | QA Lead | task-1 | Epic rejected with clear error and zero append; task/subtask succeeds. | Epic appends a journal line/row, succeeds, or fails opaquely; or task case is rejected. | Draft | Kind-guard trigger is a backstop; the CLI must reject first (no journal poisoning). |
| VAL-005 | PRD-DEV-004; PRD-DEV-010; design §3 `fid dev gate` snapshot note; DEC-C | Completing an attempt that has matching `gate_runs` rows (seeded via journal fixture, since `fid dev gate` is v0.1e) snapshots them into `task_completions.gate_results` as `[{cmd, exit_code, ran_at, output_excerpt}]` and derives `gates_passed` (exit_code 0) / `gates_failed` (exit_code ≠ 0); the source `gate_runs` rows are preserved. A completion with no matching `gate_runs` yields `gate_results = []` and `0`/`0`. | Integration: seed `gate_runs` fixture for `(issue, attempt)` with a pass and a fail, complete, assert snapshot array + derived counts + source rows preserved; plus a no-gate completion case. | QA Lead | task-2 | Both states correct: snapshot matches seeded gates with right counts and preserves source; no-gate case is empty 0/0. | Snapshot omits/duplicates rows, wrong `gates_passed`/`gates_failed`, mutates/deletes source `gate_runs`, or no-gate case is non-empty. | Draft | Real `fid dev gate` writer is v0.1e; this proves the snapshot mechanism via fixtures. Gate→attempt matching only on `(issue_id, attempt_number)`. |
| VAL-006 | PRD-DEV-001/010; PRD-WF-008 (`has_completion_report`); roadmap §v0.1c; Phase 5 fixture replacement | A real `has_completion_report`-gated transition is **rejected** when no completion exists (error names the issue, zero state change) and **succeeds** once `fid issue complete` records a real completion — driven by real write commands, no journal fixtures. | Integration/CLI: create task → attempt gated transition (reject, naming the issue, zero state change) → `fid issue complete` → re-attempt (success). | QA Lead | task-3 | Reject path names the unmet issue with no state change; success path transitions after a real completion. | Reject path omits the id or mutates state, or the transition fails after a real completion exists. | Draft | `all_acs_verified` (Phase 7) and `concerns_addressed` (Phase 9) gating are out of scope here. |
| VAL-007 | ADR-0001; roadmap §v0.1c exit | Deleting the SQLite DB and running `fid rebuild` reproduces identical `task_completions` (including `attempt_number`, `superseded` status, `gate_results` snapshot, `gates_passed`/`gates_failed`, `human_intervention`) and `gate_runs` projected state from the journal. | Integration: snapshot completion + gate state, delete DB, `fid rebuild`, assert identical projected state. | QA Lead | task-3 | Post-rebuild completion + gate state is identical to pre-delete state. | Any divergence in `task_completions`, `gate_runs`, snapshot, derived counts, or supersession after rebuild. | Draft | Held-out scenario harness location remains a separate scheduled commitment. |
| VAL-008 | Phase brief; roadmap §v0.1c | Human dogfood confirms the completion lifecycle (record an attempt; observe attempt numbering + supersession; see the gate snapshot rolled into the completion; flag a human intervention; satisfy a real `has_completion_report`-gated transition) reads and behaves as expected in a real `.fid/` workspace. | Human checklist against a real `.fid/` workspace; record outcomes per observable. | Human | phase-end | Human confirms attempt numbering/supersession, gate snapshot, human-intervention flag, and the gated transition all match expectation. | Human finds incorrect attempt/supersession, a missing/incorrect gate snapshot, a dropped intervention flag, or behavior that contradicts expectation. | Draft | QA Lead fans this into one human-checklist row per observable (record+fields / attempt+supersede / gate snapshot / intervention / gated transition). |
| VAL-009 | DEC-007 size tripwire; CF-012 | `internal/issues/issue.go` is split into cohesive same-package files, each ≤400 non-test lines, with **no behavior change** (no public API change, no test deleted/weakened) and the full task gate green. | Command: `wc -l` on `issue.go` + new split files; full `go test -race ./...` green; diff shows function moves only. | QA Lead | task-4 | `issue.go` and every new split file ≤400 non-test lines; suite green; no behavior/API change. | Any resulting non-test file >400 lines, a behavior/API change, or a weakened/deleted test. | Draft | Behavior-preserving only; mirrors the Phase 6 CF-010 transition/update/project split. |

## Coverage Rules

Confirmed before task execution begins:

- Phase-scope requirements each map to ≥1 `VAL-NNN`: PRD-DEV-001→001/004/006,
  DEV-002→002, DEV-003→002, DEV-004→005, DEV-009→003, DEV-010→005/006;
  PRD-WF-008 real-data gating→006; journal/rebuild canonical (ADR-0001)→007;
  human completion-lifecycle feel→008; CF-012 size split→009.
- PRD-DEV-005–008 (concerns, composite `fid dev *`) have **no Phase 8 assertion
  by design** — concerns are Phase 9; composites/`fid dev gate` are v0.1e. The
  phase brief records the boundary.
- Every `VAL-NNN` is referenced by a task spec or marked `phase-end` (VAL-008).
- Binary/conditional assertions name both states — VAL-002, 003, 004, 005, 006.
- DEC-D (no `activity_log` for completions; matches `create`/AC) is proposed for
  human confirmation; consistent timeline population is deferred to v0.1e via
  CF-011. No activity assertion in this contract by design.
- VAL-005 evidence is fixture-seeded by design (DEC-C): `fid dev gate` is v0.1e;
  the contract proves the snapshot mechanism, and real-write linkage is
  re-verified in v0.1e.
- QA Lead maps each `VAL-NNN` to deterministic `(qa-lead)`, mid-phase-review, or
  `(human)` verification and fans VAL-008 into per-observable human rows.

## Outcomes

Phase-end final review records every `VAL-NNN` outcome: `PASS`, `FAIL`,
`BLOCKED`, or `CARRY-FORWARD`. `PASS` without evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-18 | VAL-001..009 | created | Phase 8 planning draft | phase brief; roadmap §v0.1c; PRD-DEV-001..004/009/010; PRD-WF-008; CF-012 |
