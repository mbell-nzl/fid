# Phase Brief - Phase 24 (readiness-refresh findings applied; D6 open)

## Phase Goal

"Operating experience and input convergence" (DEC-272, CF-045). A Human
takes control of a stuck or diverging journey using fid output alone. They
read why each action is refused, see readable Stop, Attention and judgment
details, and write a valid input for any leaf from fid's own help and
template. fid observes repository identity itself, so callers send only the
expectations a check uses. Request digests cover the input as written.

## Status

**Readiness-refresh findings applied (DEC-280).** Branch `phase-24`. The
QA plan's first pass was dispositioned in DEC-278 and D5 in DEC-279. Its
readiness refresh (`87f1071`) was NOT READY, with nine Must-fix items
(RR-01..09) and 13 improvements (RI-01..13). The Tech Lead reproduced each
one and dispositioned all of them in DEC-280:
- task 9 is split, and new task 9a owns grammar failures and unknown
  commands (VAL-064/065);
- inventory R17-R19 are added;
- a control-character defect found during the check is fixed in tasks 9a
  and 10;
- RR-06 goes to the Human as **D6** (contract section 13).

- **Next:** the Human decides D6, recorded as DEC-281. Then the QA Lead
  re-checks everything DEC-280/281 changed, then the Phase Runner's
  readiness check, then task 1. Nothing is dispatched before the re-check
  passes, and task 10 also waits for DEC-281.

## Operating Mode

**Heavy** (DEC-277). Entry criteria: a public-interface change (17 inputs
reshaped, `--help`/`--example`, Human failure lines, three new read leaves)
and a persistence change (migration 0019). The full artifact set applies.
Mode B held-out runs are `N/A` as in Phases 21-23 (no harness recorded). The
Human checklist's "take control" scenario is the Human run.

## Authority

- Contract: `../phases/phase-24/public-interface.md` (ACCEPTED, DEC-277;
  clarified by DEC-278 and DEC-280; D6 open in section 13).
- Inventory: `../phases/phase-24/input-inventory.md` (R1-R19). Every task
  reruns its registry commands.
- Decisions: DEC-272 (scope), DEC-273 (real-worktree run and gate checks),
  DEC-274 (`get` leaves), DEC-275 (per-field digests, migration 0019),
  DEC-276 (help and template rules), DEC-277 (plan approval), DEC-278
  (QA-plan dispositions), DEC-279 (unknown-command `operation`) and DEC-280
  (readiness-refresh dispositions). See `../phases/phase-24/decisions.md`.

## Task Index

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Idempotency field-digest storage (migration 0019) | Ready (after the QA re-check) | None |
| 2 | Input-as-written digests and named conflict fields, every keyed leaf | Pending | 1 |
| 3 | Input descriptor framework, help, templates, field-named validation | Pending | None |
| 4 | Repository observation and SourceInput; judgment reference leaves | Pending | 2, 3 |
| 5 | Review leaves and repository expectations | Pending | 4 |
| 6 | Gate registration, Envelope and Claim leaves | Pending | 4 |
| 7 | Run leaves (DEC-273) | Pending | 5, 6 |
| 8 | Gate run and Evidence leaves; real freshness (DEC-273) | Pending | 7 |
| 9 | Activation descriptors; whole-tree input and help conformance | Pending | 3-8 |
| 9a | Grammar failures and unknown commands across the live tree (DEC-280) | Pending | 9 |
| 10 | Human failure lines in every family (waits for D6, DEC-281) | Pending | 9a |
| 11 | Readable Stop, Attention and repository lines; Human fixture generator | Pending | 10 |
| 12 | Concern, Deviation and Human Outcome `get` leaves (DEC-274) | Pending | 4, 10 |
| 13 | Judgment disclaimers once per section; recorded reasoning shown | Pending | 11, 12 |
| 14 | Human-takes-control journey and phase conformance | Pending | 1-13, 9a |

## Validation Contract

Full contract: `../phases/phase-24/qa/validation-contract.md`. It has 66
rows (61 deterministic, 5 Human), all `Ready`. VAL-058..063 were split from
bundled rows, and VAL-064..066 are new. D6 may amend VAL-027..029, VAL-036
and VAL-054.

| VAL range | Subject | Owning tasks | Validator |
|---|---|---|---|
| 001-005, 058 | Field-digest storage, migration up/down, schema 19, rebuild budget | 1, 2 | QA Reviewer (005 QA Lead) |
| 006-009, 066 | Digest as written; named conflict fields; replay never observes | 2, 4-8 | QA Reviewer |
| 010-016, 059, 060 | Short lines, descriptors, help, `--example`, field paths | 3-9 | QA Reviewer |
| 017-023 | Removed members, SourceInput, recorded observation, `repository-unavailable` | 4-8 | QA Reviewer |
| 024-035, 061, 062 | Expectations, same worktree, Runs, gate runs, freshness, fake parity | 5-9 | QA Reviewer |
| 064, 065 | Grammar failures; unknown commands (Phase 21 conformance) | 9a | QA Reviewer |
| 036-039 | Human failure lines | 10 | QA Reviewer |
| 040-041 | Stop, Attention and repository Human lines | 11 | QA Reviewer |
| 042-043 | `get` leaves | 12 | QA Reviewer |
| 044-048, 063 | Disclaimers, recorded reasoning, sentinel guard, glossary, `PRIVATE_` removal | 13 | QA Reviewer (048 QA Lead) |
| 049-052 | Take-control journey, parity, absence, wall-time budgets | 14 (052 every task) | QA Reviewer (051 QA Lead) |
| 053-057 | Human: take control, failure lines, views, guidance, authoring | phase-end (054/055 also at window D) | Human |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| QA plan | `../phases/phase-24/qa/qa-plan.md` | Refresh NOT READY; dispositions applied (DEC-280); QA re-check after D6 |
| Mid-phase reviews | windows A-E after tasks 3, 6, 8, 11 (plus the Human walk-through) and 14 | Accepted (DEC-278). Window D now spans tasks 9, 9a, 10 and 11 |
| Human fixtures | `TestP24HumanFixtures` (task 11, extended by task 14) | Pending |
| Final review / Human checklist | `../phases/phase-24/qa/` | Pending |

## Active Constraints

- No compatibility with pre-Phase-24 journals (ADR-0011, DEC-270).
- Retained-test convergence needs introducing-commit provenance
  (DEC-162/165/166). A gap in a closed registry needs the missing authority,
  not a weakened assertion (DEC-195).
- Task 8 makes freshness real. A displaced journey row is convergence of the
  accepted change, never a reason to weaken freshness.
- The journey's provenance classes are closed (DEC-278). A missing value is a
  product gap to escalate, never a reason to widen them.
- Agent prompts and `09_AI_WORKFLOW.md` gain no fid product features.

## Known Risks

- **Real freshness (task 8)** and **imported destinations** (tasks 6-8):
  see inventory section 4.
- **Rebuild headroom:** about 19.6% (QA plan, worst fresh sample 24128 ms).
  Tasks 1 and 2 run the 100k benchmark (VAL-005).
- **Breadth:** tasks 2, 9, 9a and 10 touch every family or the whole tree.
  Their corpora are enumerated from registries, not hand-picked.
- **Package time:** `cmd/fid` takes 330.8 s (55% of 600 s). Tasks 3-10 add
  corpora there, so reports watch the 480 s warning.
- **Scratch space:** `/tmp` is a 3.9 GB tmpfs with 0.7 GB free at DEC-280.
  The QA refresh measured 1.2 GB. Most of the space is held by other
  projects' directories. A gate can fail for environmental reasons
  (CF-042).

## Next Planning Checkpoint

The Human decides D6. The QA Lead then re-checks everything DEC-280 and
DEC-281 changed, in a new subsection of QA plan section 12. Phase 25
"Parallel and evolving work" follows Phase 24.
