# Phase Brief - Phase 24 (QA recheck corrections applied; verification pending)

## Phase Goal

"Operating experience and input convergence" (DEC-272, CF-045). A Human
takes control of a stuck or diverging journey using fid output alone. They
read why each action is refused, see readable Stop, Attention and judgment
details, and write a valid input for any leaf from fid's own help and
template. fid observes repository identity itself, so callers send only the
expectations a check uses. Request digests cover the input as written.

## Status

**QA recheck corrections applied (DEC-282).** Branch `phase-24`. The
first QA pass and readiness refresh were dispositioned in DEC-278..281,
including task 9a for grammar/unknown commands and D6=A for generic-fallback
messages. Recheck `4819b30` confirmed those dispositions but remains
**NOT READY** with RC-01/RC-02. DEC-282 addresses both:

- plain-copy `run.start` preserves the earlier Claim refusal,
  `run claim binding is invalid` (task 7; message convergence task 10);
- grammar and unknown-command messages use the shared escape-then-bound
  rule, at most 512 scalars, retaining full details (task 9a; pairing task 10).

**Next:** QA Lead rechecks RC-01/RC-02 and DEC-282's changed text, then the
Phase Runner checks readiness, then task 1. Nothing is dispatched before
the QA recheck passes. No Phase 24 implementation has started.

## Operating Mode

**Heavy** (DEC-277). Entry criteria: a public-interface change (17 inputs
reshaped, `--help`/`--example`, Human failure lines, three new read leaves)
and a persistence change (migration 0019). The full artifact set applies.
Mode B held-out runs are `N/A` as in Phases 21-23 (no harness recorded). The
Human checklist's "take control" scenario is the Human run.

## Authority

- Contract: `../phases/phase-24/public-interface.md` (ACCEPTED, DEC-277;
  clarified by DEC-278/280/282; D6 decided by DEC-281).
- Inventory: `../phases/phase-24/input-inventory.md` (R1-R20). Every task
  reruns its registry commands.
- Decisions: DEC-272 (scope), DEC-273 (real-worktree run and gate checks),
  DEC-274 (`get` leaves), DEC-275 (per-field digests, migration 0019),
  DEC-276 (help and template rules), DEC-277 (plan approval), DEC-278
  (QA-plan dispositions), DEC-279 (unknown-command `operation`), DEC-280
  (readiness-refresh dispositions), DEC-281 (generic-fallback messages),
  DEC-282 (QA recheck corrections). See `../phases/phase-24/decisions.md`.

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
| 10 | Human failure lines in every family; generic-fallback messages (DEC-281) | Pending | 9a |
| 11 | Readable Stop, Attention and repository lines; Human fixture generator | Pending | 10 |
| 12 | Concern, Deviation and Human Outcome `get` leaves (DEC-274) | Pending | 4, 10 |
| 13 | Judgment disclaimers once per section; recorded reasoning shown | Pending | 11, 12 |
| 14 | Human-takes-control journey and phase conformance | Pending | 1-13, 9a |

## Validation Contract

Full contract: `../phases/phase-24/qa/validation-contract.md`. It has 67
rows (62 deterministic, 5 Human), all `Ready`. VAL-058..063 were split from
bundled rows. VAL-064..066 are new (DEC-278), and so is VAL-067 (DEC-281).
DEC-282 refines 029/036/064/065/067; `Ready` is planning status, not QA PASS.

| VAL range | Subject | Owning tasks | Validator |
|---|---|---|---|
| 001-005, 058 | Field-digest storage, migration up/down, schema 19, rebuild budget | 1, 2 | QA Reviewer (005 QA Lead) |
| 006-009, 066 | Digest as written; named conflict fields; replay never observes | 2, 4-8 | QA Reviewer |
| 010-016, 059, 060 | Short lines, descriptors, help, `--example`, field paths | 3-9 | QA Reviewer |
| 017-023 | Removed members, SourceInput, recorded observation, `repository-unavailable` | 4-8 | QA Reviewer |
| 024-035, 061, 062 | Expectations, same worktree, Runs, gate runs, freshness, fake parity | 5-9 (027-029 messages: 10) | QA Reviewer |
| 064, 065 | Grammar failures; unknown commands (Phase 21 conformance) | 9a | QA Reviewer |
| 036-039, 067 | Human failure lines; generic-fallback messages | 10 | QA Reviewer |
| 040-041 | Stop, Attention and repository Human lines | 11 | QA Reviewer |
| 042-043 | `get` leaves | 12 | QA Reviewer |
| 044-048, 063 | Disclaimers, recorded reasoning, sentinel guard, glossary, `PRIVATE_` removal | 13 | QA Reviewer (048 QA Lead) |
| 049-052 | Take-control journey, parity, absence, wall-time budgets | 14 (052 every task) | QA Reviewer (051 QA Lead) |
| 053-057 | Human: take control, failure lines, views, guidance, authoring | phase-end (054/055 also at window D2) | Human |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| QA plan | `../phases/phase-24/qa/qa-plan.md` | Recheck NOT READY (`4819b30`); RC-01/02 addressed by DEC-282; fresh QA recheck next |
| Mid-phase reviews | A after 3; B after 6; C after 8; D1 after 9a; D2 after 11; E after 14 | Six windows per QA plan section 8 (DEC-280/282). Human walkthrough after task 11 |
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

The QA Lead rechecks RC-01/RC-02 and DEC-282's changed text in a new
subsection of QA plan section 12, retaining the prior review. Phase 25
"Parallel and evolving work" follows Phase 24.
