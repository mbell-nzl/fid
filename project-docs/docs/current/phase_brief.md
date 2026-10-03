# Phase Brief - Phase 24 (Planning)

## Phase Goal

"Operating experience and input convergence" (DEC-272, CF-045). A Human
takes control of a stuck or diverging journey using fid output alone. They
read why each action is refused, see readable Stop, Attention and judgment
details, and write a valid input for any leaf from fid's own help and
template. fid observes repository identity itself, so callers send only the
expectations a check uses. Request digests cover the input as written.

## Status

**Planning.** Branch `phase-24` (from `main` after PR #25). The contract draft
and executable inventory are written. D1-D3 are decided (DEC-273/274/275).
D4 awaits the Human (`../phases/phase-24/public-interface.md` section 13),
then approval of the mode and task list. The validation contract, task specs and QA plan follow
approval. Nothing is dispatched until then.

## Operating Mode

**Heavy (proposed).** Entry criteria: a public-interface change (17 inputs
reshaped, `--example`, Human failure lines, three new read leaves under
DEC-274) and a persistence change (migration 0019). The full artifact set applies. Mode B
held-out runs are `N/A` as in Phases 21-23 (no harness recorded); the
"Human takes control" checklist scenario is the Human run.

## Task Index (proposed; specs written after approval)

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Idempotency field-digest storage (migration 0019) | Proposed | None |
| 2 | Input-as-written digests and named conflict fields, all families | Proposed | 1 |
| 3 | Input descriptor framework, `--help`, `--example`, field-named validation; Stop/Attention leaves; all Short lines | Proposed | None |
| 4 | Observation and SourceInput foundation; judgment reference leaves | Proposed | 2, 3 |
| 5 | Review leaves; `repository-mismatch` details | Proposed | 4 |
| 6 | Gate registration, Envelope and Claim leaves | Proposed | 4 |
| 7 | Run leaves (DEC-273) | Proposed | 5, 6 |
| 8 | Gate run and Evidence leaves; real freshness (DEC-273) | Proposed | 7 |
| 9 | Activation descriptors; remove pending fence; whole-tree help conformance | Proposed | 3-8 |
| 10 | Human failure lines in every family (A5) | Proposed | 9 |
| 11 | Readable Stop, Attention and repository lines (A4) | Proposed | 10 |
| 12 | Concern, Deviation and Human Outcome `get` leaves (DEC-274) | Proposed | 4 |
| 13 | Disclaimers once per section; recorded reasoning (A6/A7) | Proposed | 11, 12 |
| 14 | "Human takes control" real-process journey; phase conformance | Proposed | 1-13 |

## Validation Contract

Drafted after D1-D4 at `../phases/phase-24/qa/validation-contract.md`. VAL-060
and VAL-073 from Phase 23 are restated, not weakened (contract section 11).

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| Contract draft | `../phases/phase-24/public-interface.md` | Draft |
| Executable inventory | `../phases/phase-24/input-inventory.md` | Draft |
| QA plan | `../phases/phase-24/qa/qa-plan.md` | Missing |
| Mid-phase reviews | windows after tasks 3, 8, 11 (+ Human walk-through of new failure and Stop output), 14 | Proposed |
| Final review / Human checklist | `../phases/phase-24/qa/` | Pending |

## Active Constraints

- Contract first: no task spec before the Human approves the contract.
- No compatibility with pre-Phase-24 journals (ADR-0011, DEC-270).
- Agent prompts and `09_AI_WORKFLOW.md` gain no fid product features.
- Split any task that would touch more than one closed journal registry
  (Phase 23 workflow notes).

## Known Risks

- **Real freshness (task 8).** Phase 21-23 journeys stay fresh only because
  the caller's block is copied. Displaced rows are convergence; see inventory
  section 4.
- **Registry misses.** Every task reruns the inventory commands it owns.

## Next Planning Checkpoint

The Human decides D1-D4. The Tech Lead then drafts the validation contract
and the 14 task specs, and the QA Lead writes the QA plan.
