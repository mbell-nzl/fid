# Phase Brief - Phase 6 (Ready Queue + v0.1b Closure)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** Phase 6 is **PLANNED** (validation contract + 5 task specs drafted).
> Phase 5 is CLOSED (QA PASS WITH IMPROVEMENTS; dogfood 6/6; VAL-001..013 PASS).
> Archives under `../phases/phase-5/`. DEC-017/DEC-018 **CONFIRMED by human
> 2026-06-14**. QA Lead refines the QA plan / evidence strategy, then execution
> starts at Task 1.

## Phase 6 Goal (roadmap §v0.1b, Phase 6 cut)

`fid ready` + v0.1b closure. **Exit:** the workflow runner can ask for the next
actionable task and receive the same scheduling answer a human expects.

Scope: `fid ready` (sequential epics, prioritised epics, unattached tasks,
project-wide selection, empty-result explanations, format dispatch) + an
end-to-end v0.1b smoke over create/list/get/update/transition/ready. Carry-
forward CF-009 and CF-010 are fixed first (Tasks 1–2).

Explicit deferrals (do not implement): `fid ready --explain`/`--limit`/
`--include-blocked` (v0.2, PRD-ISSUE-024); reparent / inbox promotion (v0.2,
PRD-ISSUE-025); dependency `blocked_by` link blocking and claim blocking (v0.1d);
label CRUD (later slice — use existing `labels` rows).

## Confirmed decisions (DEC-017, DEC-018) — human-approved 2026-06-14

v0.1b has **no `blocked_by` links and no claims** (both v0.1d), so readiness is
defined from workflow status category alone:

- **DEC-017 (readiness semantics):** actionable/ready = `unstarted`/`todo`;
  unfinished = not `completed`/`cancelled`. A sequential epic's lowest-sequence
  unfinished task that is in a `started` category is *not* actionable → the epic
  yields no ready task and the empty result explains it.
- **DEC-018 (empty-state contract):** bare `fid ready` finding nothing ready
  exits 0, emits `{"task": null}` on stdout (frozen shape — Workflow Runner
  depends on it), and writes a `note:` explanation to stderr.

Both confirmed; VAL-003 / VAL-007 are locked as written. (TL records DEC-017 /
DEC-018 in `../phases/phase-6/decisions.md` at wrap.)

## Carry-Forward Disposition (this phase)

| ID | Item | Disposition |
|---|---|---|
| CF-009 | `fid issue update` one-field-per-line output (DEC-016) | **Do now — Task 1.** |
| CF-010 | Split oversized files (transition/update/project) | **Do now — Task 2.** |
| CF-003 | Role-validation corrupt-tail tolerance | Re-defer to before v0.1d. |
| CF-008 | Resolution type / cancellation reason (DEC-015) | Re-defer to v0.2 (product). |

## Validation Contract Index

Full contract: `../phases/phase-6/qa/validation-contract.md`.

| VAL | Behavior (short) | Owning task | Validator |
|---|---|---|---|
| VAL-001 | Bare `fid ready` → one task in frozen §6.3 JSON shape | task-4 | QA Lead |
| VAL-002 | Sequential epic → lowest-`sequence` unfinished actionable task | tasks-3-4 | QA Lead |
| VAL-003 | Sequential current task not actionable → no task + blocker reason (both states) | tasks-3-4 | QA Lead |
| VAL-004 | Prioritised epic → highest-priority actionable task | tasks-3-4 | QA Lead |
| VAL-005 | Unattached actionable task eligible; non-actionable excluded | tasks-3-4 | QA Lead |
| VAL-006 | Project-wide → one deterministic top task; sequential contributes current only | tasks-3-4 | QA Lead |
| VAL-007 | Empty state → exit 0, `{"task": null}`, stderr note; non-empty has no note | task-4 | QA Lead |
| VAL-008 | `--format json\|toon\|markdown\|yaml` render selected + empty; JSON default | task-4 | QA Lead |
| VAL-009 | CF-009: `fid issue update` one changed field per line (both states) | task-1 | QA Lead |
| VAL-010 | CF-010: 3 files ≤400 lines, behavior unchanged + rebuild parity | task-2 | QA Lead |
| VAL-011 | E2E v0.1b smoke create→list→get→update→transition→ready + rebuild parity | task-5 | QA Lead |
| VAL-012 | Human dogfood: `fid ready` matches expectation (seq / prioritised / empty) | phase-end | Human |

PRD-ISSUE-024 (`--explain`) and PRD-ISSUE-025 (reparent/promote) have no Phase 6
assertion by design (v0.2).

## Ordered Task List

1. **Task 1 — CF-009:** `fid issue update` one-field-per-line output (VAL-009).
2. **Task 2 — CF-010:** split `transition.go`/`update.go`/`project.go` ≤400
   lines, behavior-preserving (VAL-010).
3. **Task 3 — ready engine:** `internal/issues/ready.go` selection logic +
   exhaustive unit tests (VAL-002..006, supports VAL-001).
4. **Task 4 — `fid ready` CLI:** command, `--epic`, `--format` dispatch,
   empty-state contract (VAL-001, 002..006 via CLI, 007, 008).
5. **Task 5 — v0.1b smoke:** end-to-end integration test + held-out scenario
   recommendation (VAL-011); seeds VAL-012 dogfood.

## Planning Reminders

- DEC-007 carry-in: COV-001 (`internal/issues`, `internal/journal`,
  `internal/policy`, `internal/workflow` ≥ 80%); COV-002 `cmd/fid` ≥ 70%.
  Phase 5 ended: `cmd/fid` 79.0%, `internal/issues` 82.8%, `internal/journal`
  80.3%, `internal/policy` 85.0%, `internal/workflow` 95.5%.
- Phase 4/5 deep-retro lesson: when an AC asserts recovery / journal hygiene for
  a rejected write, the task's original evidence must demonstrate that state.
- `fid ready` is read-only scheduling; product-feel surfaces in the dogfood
  (VAL-012) — that is where DEC-017/DEC-018 get a real-world sanity check.
- Phase 6 closes v0.1b: Task 5 cross-command smoke also feeds the standing
  held-out-scenario scheduled commitment (recommendation, not relocation).

## Next Checkpoint

DEC-017/DEC-018 confirmed. QA Lead: refine evidence strategy into
`../phases/phase-6/qa/qa-plan.md`, fan VAL-012 into per-observable human rows,
then Phase Runner starts Task 1.
