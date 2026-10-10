# Phase Brief - Phase 24 (phase-end; final review BLOCK, ESC-24-010 under DEC-292)

## Phase Goal

"Operating experience and input convergence" (DEC-272, CF-045). A Human
takes control of a stuck or diverging journey using fid output alone. They
read why each action is refused, see readable Stop, Attention and judgment
details, and write a valid input for any leaf from fid's own help and
template. fid observes repository identity itself, so callers send only the
expectations a check uses. Request digests cover the input as written.

## Status

**Phase-end on branch `phase-24`.** Tasks 1-14 are complete. QA windows
A-E have passed (review-12 closed E). The final review `qa/final-review.md`
(`6d633351`) returned BLOCK: 63 deterministic VAL PASS and every gate
passing. VAL-053..057 are BLOCKED because the Human checklist is blank
(B1-H). QA-36 fails on two report-only findings: Task 1 lacks its owned R4
receipt (F1), and Task 10 its owned R10 receipt (F2). TRG-24-019 raised
ESC-24-010; the ledger is `STOPPED` (12:12). **DEC-292** disposes it;
ESC-24-010 stays OPEN until the Human records the checklist.

**Next:**

1. Human: record every row of `../phases/phase-24/qa/human-checklist.md`
   (143 rows) and commit it (DEC-292 clause 1). Live workspaces
   (HUMAN-01, HUMAN-05) expire about 11:50 NZDT 2026-10-12; regenerate
   them first after that.
2. Tech Lead: append DEC-292 clause 9 and close ESC-24-010, or disposition
   any FAIL/PARTIAL row in a new decision (clause 2).
3. Phase Runner, DEC-292 clause 6: evidence-only Task 14 recovery, then
   report-only Task 1 (R4) and Task 10 (R10) recoveries.
4. Fresh QA Lead `qa/final-review-1.md` (review-13) closes B1-H, F1, F2
   and QA-36 (clause 7); a PASS leads to RUN-END and wrap-up.

## Operating Mode

**Heavy** (DEC-277). Entry criteria: a public-interface change (17 inputs
reshaped, `--help`/`--example`, Human failure lines, three new read leaves)
and a persistence change (migration 0019). The full artifact set applies.
Mode B held-out runs are `N/A` as in Phases 21-23 (no harness recorded). The
Human checklist's "take control" scenario is the Human run.

## Authority

- Contract: `../phases/phase-24/public-interface.md` (ACCEPTED, DEC-277;
  clarified by DEC-278/280/282/285/286/287; D6 by DEC-281; DEC-290 adds).
- Inventory: `../phases/phase-24/input-inventory.md` (R1-R20). Every task
  reruns its registry commands; a table command is its rendered text (DEC-291).
- Decisions: DEC-272 (scope), DEC-273 (real-worktree run and gate checks),
  DEC-274 (`get` leaves), DEC-275 (per-field digests, migration 0019),
  DEC-276 (help and template rules), DEC-277 (plan approval), DEC-278
  (QA-plan dispositions), DEC-279 (unknown-command `operation`), DEC-280
  (readiness-refresh dispositions), DEC-281 (generic-fallback messages),
  DEC-282 (QA recheck corrections), DEC-283 (Task 6 resume reconciliation),
  DEC-284 (Task 9 correcting TRIAGE; late window C), DEC-285 (`activity.list`
  failure mapping; Task 10 caller scope), DEC-286 (root-flag refusal
  mapping), DEC-287 (plain-copy `run.inspect` keeps its refusal), DEC-288
  (window D2 Human walkthrough; evidence-only resume), DEC-289 (Task 13
  retained-predicate scope; `None.` marker), DEC-290 (`resume` Issue
  revision and next states; Task 14 scope), DEC-291 (table commands read
  as rendered; Task 13 R13 receipt resume), DEC-292 (phase-end
  disposition; Human checklist, F1/F2 route). See
  `../phases/phase-24/decisions.md`.

## Task Index

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Idempotency field-digest storage (migration 0019) | Complete; R4 receipt pending (final F1, DEC-292) | None |
| 2 | Input-as-written digests and named conflict fields, every keyed leaf | Complete | 1 |
| 3 | Input descriptor framework, help, templates, field-named validation | Complete (window A PASS) | None |
| 4 | Repository observation and SourceInput; judgment reference leaves | Complete (window B PASS) | 2, 3 |
| 5 | Review leaves and repository expectations | Complete (window B PASS) | 4 |
| 6 | Gate registration, Envelope and Claim leaves | Complete (window B PASS) | 4 |
| 7 | Run leaves (DEC-273) | Complete (window C PASS) | 5, 6 |
| 8 | Gate run and Evidence leaves; real freshness (DEC-273) | Complete (window C PASS) | 7 |
| 9 | Activation descriptors; whole-tree input and help conformance | Complete (window D1 PASS) | 3-8 |
| 9a | Grammar failures and unknown commands across the live tree (DEC-280) | Complete (window D1 PASS) | 9 |
| 10 | Human failure lines in every family; generic-fallback messages (DEC-281) | Complete (`a6b6203`); R10 receipt pending (final F2, DEC-292) | 9a |
| 11 | Readable Stop, Attention and repository lines; Human fixture generator | Complete (`93a1f792`; window D2 PASS, review-10) | 10 |
| 12 | Concern, Deviation and Human Outcome `get` leaves (DEC-274) | Complete (`fc4bdff7`) | 4, 10 |
| 13 | Judgment disclaimers once per section; recorded reasoning shown | Complete (`a9aa99f9`; R13 receipt DEC-291, review-12 PASS) | 11, 12 |
| 14 | Human-takes-control journey and phase conformance; `resume` Issue position (DEC-290) | Complete (`86e0c26a`, DEC-290) | 1-13, 9a |

## Validation Contract

Full contract: `../phases/phase-24/qa/validation-contract.md`: 68 rows (63
deterministic, 5 Human), all `Ready` (planning status, not QA PASS). Split:
058..063. Added: 064..066 (DEC-278), 067 (DEC-281), 068 (DEC-290).

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
| 049-052, 068 | Take-control journey, parity, absence, wall-time budgets; `resume` Issue position | 14 (052 every task) | QA Reviewer (051 QA Lead) |
| 053-057 | Human: take control, failure lines, views, guidance, authoring | phase-end (054/055 also at window D2) | Human |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| QA plan | `../phases/phase-24/qa/qa-plan.md` | READY (`7e41e6b`, section 12.7); RC-01/02 closed |
| Mid-phase reviews | A after 3 (review-2 PASS); B after 6 (review-4 PASS); C late, after 9 (review-6 PASS, DEC-284); D1 after 9a (review-8 PASS); D2 after 11 (review-9 BLOCK; review-10 PASS, DEC-288); E after 14 (also VAL-068; review-11 REPAIR F1, review-12 PASS under DEC-291) | Six windows per QA plan section 8 (DEC-280/282). Human walkthrough after task 11 |
| Human fixtures | `TestP24HumanFixtures` (task 11, extended by task 14) | D2 interim walkthrough `a23b9f33`; phase-end set `fixtures-2145711523` (11:50 2026-10-11) |
| Final review / Human checklist | `../phases/phase-24/qa/` | `final-review.md` BLOCK (`6d633351`); checklist awaits the Human; `final-review-1.md` next (DEC-292) |

## Active Constraints

- No compatibility with pre-Phase-24 journals (ADR-0011, DEC-270).
- Retained-test convergence needs introducing-commit provenance
  (DEC-162/165/166, Phase 21: method only; ledger authority fields cite
  Phase 24 decisions, DEC-289). A gap in a closed registry needs the missing authority,
  not a weakened assertion (DEC-195).
- Task 8 makes freshness real. A displaced journey row is convergence of the
  accepted change, never a reason to weaken freshness.
- The journey's provenance classes are closed (DEC-278). A missing value is a
  product gap to escalate, never a reason to widen them. DEC-290 is the only
  product change Task 14 may make.
- Agent prompts and `09_AI_WORKFLOW.md` gain no fid product features.

## Known Risks

- **Real freshness (task 8)** and **imported destinations** (tasks 6-8):
  see inventory section 4.
- **Rebuild headroom:** about 19.6% (QA plan, worst fresh sample 24128 ms).
  Tasks 1 and 2 run the 100k benchmark (VAL-005).
- **Breadth:** tasks 2, 9, 9a and 10 touch every family or the whole tree.
  Their corpora are enumerated from registries, not hand-picked.
- **Package time:** `cmd/fid` took 484.7 s at Task 10 close (warning 480 s,
  budget 600 s) and 488.2 s at Task 14 close (warning recorded).
- **Scratch space:** `/tmp` is a 3.9 GB tmpfs (1.3 GB free at DEC-290),
  mostly held by other projects. A gate can fail for environmental reasons
  (CF-042).

## Next Planning Checkpoint

DEC-292 clause 9 after the Human checklist, then final-review-1 and
wrap-up. Phase 25 "Parallel and evolving work" follows.
