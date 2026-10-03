# Phase Brief - Phase 24 (Planning approved; QA plan pending)

## Phase Goal

"Operating experience and input convergence" (DEC-272, CF-045). A Human
takes control of a stuck or diverging journey using fid output alone. They
read why each action is refused, see readable Stop, Attention and judgment
details, and write a valid input for any leaf from fid's own help and
template. fid observes repository identity itself, so callers send only the
expectations a check uses. Request digests cover the input as written.

## Status

**Planning approved (DEC-277).** Branch `phase-24`. The contract, inventory,
validation contract (57 rows, `Draft`) and 14 task specs are written. Next:
the QA Lead writes `qa/qa-plan.md` with fresh command and feasibility
validation (Heavy) and proposes assertion refinements. The Tech Lead applies
them and marks rows `Ready`. Then the Phase Runner starts task 1. Nothing is
dispatched before the QA plan is ready.

## Operating Mode

**Heavy** (DEC-277). Entry criteria: a public-interface change (17 inputs
reshaped, `--help`/`--example`, Human failure lines, three new read leaves)
and a persistence change (migration 0019). The full artifact set applies.
Mode B held-out runs are `N/A` as in Phases 21-23 (no harness recorded). The
Human checklist's "take control" scenario is the Human run.

## Authority

- Contract: `../phases/phase-24/public-interface.md` (ACCEPTED, DEC-277).
- Inventory: `../phases/phase-24/input-inventory.md`. Every task reruns its
  registry commands.
- Decisions: DEC-272 (scope), DEC-273 (real-worktree run and gate checks),
  DEC-274 (`get` leaves), DEC-275 (per-field digests, migration 0019),
  DEC-276 (help and template rules), DEC-277 (plan approval). See
  `../phases/phase-24/decisions.md`.

## Task Index

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Idempotency field-digest storage (migration 0019) | Ready (after QA plan) | None |
| 2 | Input-as-written digests and named conflict fields, all families | Pending | 1 |
| 3 | Input descriptor framework, help, templates, field-named validation | Pending | None |
| 4 | Repository observation and SourceInput; judgment reference leaves | Pending | 2, 3 |
| 5 | Review leaves and repository expectations | Pending | 4 |
| 6 | Gate registration, Envelope and Claim leaves | Pending | 4 |
| 7 | Run leaves (DEC-273) | Pending | 5, 6 |
| 8 | Gate run and Evidence leaves; real freshness (DEC-273) | Pending | 7 |
| 9 | Activation descriptors; whole-tree help and template conformance | Pending | 3-8 |
| 10 | Human failure lines in every family | Pending | 9 |
| 11 | Readable Stop, Attention and repository lines | Pending | 10 |
| 12 | Concern, Deviation and Human Outcome `get` leaves (DEC-274) | Pending | 4 |
| 13 | Judgment disclaimers once per section; recorded reasoning shown | Pending | 11, 12 |
| 14 | Human-takes-control journey and phase conformance | Pending | 1-13 |

## Validation Contract

Full contract: `../phases/phase-24/qa/validation-contract.md`. There are 57
rows (52 deterministic, 5 Human), all `Draft`.

| VAL range | Subject | Owning tasks | Validator |
|---|---|---|---|
| 001-005 | Field-digest storage, schema 19, rebuild budget | 1, 2 | QA Reviewer (005 QA Lead) |
| 006-009 | Digest as written; named conflict fields | 2, 4-8 | QA Reviewer |
| 010-016 | Short lines, descriptors, help, `--example`, field paths | 3-9 | QA Reviewer |
| 017-023 | Removed members, SourceInput, recorded observation, `repository-unavailable` | 4-8 | QA Reviewer |
| 024-035 | Expectations, same worktree, Runs, gate runs, freshness, fake parity | 5-9 | QA Reviewer |
| 036-039 | Human failure lines | 10 | QA Reviewer |
| 040-041 | Stop, Attention and repository Human lines | 11 | QA Reviewer |
| 042-043 | `get` leaves | 12 | QA Reviewer |
| 044-048 | Disclaimers, recorded reasoning, sentinel guard, glossary | 13 | QA Reviewer (048 QA Lead) |
| 049-052 | Take-control journey, parity, absence, wall-time budgets | 14 (052 every task) | QA Reviewer (051 QA Lead) |
| 053-057 | Human: take control, failure lines, views, guidance, authoring | phase-end (054/055 also at the task-11 window) | Human |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| QA plan | `../phases/phase-24/qa/qa-plan.md` | Missing (QA Lead next) |
| Mid-phase reviews | windows after tasks 3, 8, 11 (plus a Human walk-through of new failure and Stop output), 14 | Proposed in DEC-277 |
| Final review / Human checklist | `../phases/phase-24/qa/` | Pending |

## Active Constraints

- No compatibility with pre-Phase-24 journals (ADR-0011, DEC-270).
- Retained-test convergence needs introducing-commit provenance
  (DEC-162/165/166). A gap in a closed registry needs the missing authority,
  not a weakened assertion (DEC-195).
- Task 8 makes freshness real. A displaced journey row is convergence of the
  accepted change, never a reason to weaken freshness.
- Agent prompts and `09_AI_WORKFLOW.md` gain no fid product features.

## Known Risks

- **Real freshness (task 8)** and **imported destinations** (tasks 6-8):
  see inventory section 4.
- **Rebuild headroom:** about 10% at Phase 23 end. Tasks 1 and 2 run the 100k
  benchmark (VAL-005).
- **Breadth:** tasks 2 and 10 touch every family. Their corpora are
  enumerated from registries, not hand-picked.

## Next Planning Checkpoint

The QA Lead writes the QA plan; the Tech Lead applies refinements and marks
rows `Ready`. Phase 25 "Parallel and evolving work" follows Phase 24.
