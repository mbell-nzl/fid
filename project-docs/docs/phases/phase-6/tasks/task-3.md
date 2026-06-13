# Phase 6 - Task 3: `fid ready` selection engine (internal)

## Objective

Implement the read-only ready-queue **selection engine** in `internal/issues`
that, given the projected DB state, computes the next actionable task for: a
sequential epic, a prioritised epic, unattached tasks, and project-wide
selection — plus a structured empty/blocked explanation when nothing is ready.
This task is library-only; Task 4 wires it to the `fid ready` CLI. Exhaustive
deterministic unit tests live here.

## Requirements

- PRD refs: PRD-ISSUE-008, 009, 010, 011, 012
- TRD refs: TRD §4.3 (clean reads need no write lock; matched `projection_meta`);
  TRD §4.6 (workflow status categories)
- Design refs: `05_DESIGN_PLAN.md` §6.3 (bare ready task shape);
  `06_DESIGN_PLAN_CLI.md` §`fid ready`; GLOSSARY "Sequential epic",
  "Prioritised epic", "Sequence"
- Quality refs: `08_QUALITY_PLAN.md` STRUCT-001/002

## Validation Contract Coverage

- VAL-002: sequential epic returns lowest-`sequence` unfinished actionable task,
  never a higher-sequence task.
- VAL-003: sequential epic with non-actionable current task returns no task +
  structured blocker reason; actionable case returns the task.
- VAL-004: prioritised epic returns highest-priority actionable task; started /
  lower-priority not chosen ahead.
- VAL-005: unattached actionable task is a candidate; non-actionable unattached
  task excluded.
- VAL-006: project-wide selection returns one deterministic top task across
  prioritised epics + unattached, sequential epics contributing only their
  current actionable task.

## Scope

### Files To Create Or Modify

- `internal/issues/ready.go` — selection engine: query projected issues +
  `epic_metadata`, apply DEC-017 readiness rules, return the selected task (in
  the §6.3 field set) plus a structured empty/blocker explanation.
- `internal/issues/ready_test.go` — table-driven unit tests over fixtures for
  sequential, prioritised, unattached, and project-wide cases, both ready and
  empty/blocked.

### Files To Read For Context

- `internal/issues/list.go`, `internal/issues/detail.go` (existing read/query
  patterns against the projection and `issues_enriched`)
- `internal/issues/issue.go` (issue/epic-metadata types, status category helpers)
- `internal/workflow/workflow.go` (status → category mapping, if exposed)

### Out Of Scope

- CLI wiring, `--format` dispatch, stderr `note` rendering, exit codes (Task 4).
- `fid ready --explain` diagnostic shape (v0.2, PRD-ISSUE-024).
- Dependency-link (`blocked_by`) and claim blocking (v0.1d).

## Acceptance Criteria

- [ ] Readiness follows DEC-017: actionable = `unstarted`/`todo`; unfinished =
  not `completed`/`cancelled` (VAL-002..006).
- [ ] Sequential epic: lowest-`sequence` unfinished actionable task selected;
  higher-sequence never chosen ahead (VAL-002).
- [ ] Sequential epic: lowest-`sequence` unfinished task that is `started`
  yields no task and a structured reason naming the blocking current task; an
  actionable lowest task is returned (VAL-003).
- [ ] Prioritised epic: highest-priority actionable task selected; started /
  lower-priority not chosen ahead (VAL-004).
- [ ] Unattached actionable task eligible; non-actionable unattached excluded
  (VAL-005).
- [ ] Project-wide: one deterministic top task across prioritised + unattached,
  sequential contributing only its current actionable task; tie-break order is
  documented and deterministic (VAL-006).
- [ ] Engine returns the §6.3 field set (`id`, `title`, `status`, `epic_id`,
  `epic_title`, `sequence`, `priority`, `est_tokens`) and a structured
  empty-result explanation value (VAL-001 support for Task 4).

## TDD Plan

- First failing test: sequential fixture (seq 1 done, seq 2 todo, seq 3 todo)
  selects seq 2.
- Additional tests: sequential blocked (seq 2 in_progress → no task + reason);
  prioritised highest-priority selection with a started higher-priority task
  excluded; unattached eligibility; project-wide deterministic winner across
  mixed epics; all-done/all-cancelled empty.
- Regression risks: status-category classification mismatches; sequential
  "unfinished" vs "actionable" confusion; non-deterministic ordering without an
  explicit tie-break.

## Documentation Updates

- `project-docs/specs/GLOSSARY.md` — only if a new precise term is needed for
  v0.1b status-category readiness; otherwise none (flag, do not invent terms).

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None (engine logic; the v0.1b smoke is Task 5).
- Cross-surface parity impact: the engine's selected-task field set must match
  design §6.3 exactly so Task 4's JSON shape is correct.
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Hold COV-001 (`internal/issues` ≥ 80%).

Conditional commands:

- None.

## Notes For Developer

- Readiness semantics for v0.1b are fixed by DEC-017 (see the validation
  contract scoping note): there are no `blocked_by` links or claims yet, so
  "blocked" means the current task is in a `started` category. Do not invent
  dependency-link handling.
- Define and document the project-wide tie-break order explicitly (e.g.
  priority, then `created_at`, then `id`) so VAL-006 is deterministic. If the
  PRD does not pin a tie-break, record your choice in the completion report for
  TL confirmation rather than leaving it implicit.
- Return a structured explanation (not a pre-formatted string) so Task 4 owns
  the `note:` wording and format dispatch.
- Reads should follow the clean-read path (no write lock when `projection_meta`
  matches); reuse existing list/detail query helpers.
