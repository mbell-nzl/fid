# Phase Brief - Phase 24 (executing; Task 14 resume under DEC-290)

## Phase Goal

"Operating experience and input convergence" (DEC-272, CF-045). A Human
takes control of a stuck or diverging journey using fid output alone. They
read why each action is refused, see readable Stop, Attention and judgment
details, and write a valid input for any leaf from fid's own help and
template. fid observes repository identity itself, so callers send only the
expectations a check uses. Request digests cover the input as written.

## Status

**Executing on branch `phase-24`.** Tasks 1-13 are complete (Task 12 at
`fc4bdff7`, Task 13 at `fd6783f9` under DEC-289). QA windows A, B, C (late,
DEC-284), D1 and D2 have passed; review-10 (`3f920f31`) closed D2.

Task 14 attempt 1 returned BLOCKED at `4b6a3837` / `98d45617`. The journey
had no permitted source for the terminal-success `--to` state (`qa_pass`),
and none for the Issue revision `--expected-revision` needs. TRG-24-016
raised ESC-24-008, and the ledger is `STOPPED` (2026-10-11 08:18).
**DEC-290 (2026-10-11)** closes it. `resume` gains `Revision:` and
`Next states:` lines (JSON `issue.next_states`), and the `issue transition`
help points there. Task 14 owns the change (bounded scope, VAL-068). The
journey takes both values as `fid-output`; DEC-278's classes are unchanged.
The Human may override the choice by DEC before Task 14 closes.

**Next:**

1. Phase Runner, DEC-290 clause 8: append RUN-RESUME (TRG-24-016 /
   ESC-24-008, attempt 2, `Resolution authority: DEC-290`), then
   TASK-RECOVERY task-14 (`in-scope-implementation`, `task-scope`,
   `task-14:public-interface-decision:terminal-state-provenance-qa-pass`,
   attempt 2, `0/2`). Dispatch one same-task continuation, then REPAIR
   `RESOLVED` and TASK task-14 COMPLETE.
2. QA window E (tasks 12-14, including VAL-068), then the final review and
   the Human checklist.

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
  reruns its registry commands.
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
  revision and next states; Task 14 scope). See
  `../phases/phase-24/decisions.md`.

## Task Index

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Idempotency field-digest storage (migration 0019) | Complete | None |
| 2 | Input-as-written digests and named conflict fields, every keyed leaf | Complete | 1 |
| 3 | Input descriptor framework, help, templates, field-named validation | Complete (window A PASS) | None |
| 4 | Repository observation and SourceInput; judgment reference leaves | Complete (window B PASS) | 2, 3 |
| 5 | Review leaves and repository expectations | Complete (window B PASS) | 4 |
| 6 | Gate registration, Envelope and Claim leaves | Complete (window B PASS) | 4 |
| 7 | Run leaves (DEC-273) | Complete (window C PASS) | 5, 6 |
| 8 | Gate run and Evidence leaves; real freshness (DEC-273) | Complete (window C PASS) | 7 |
| 9 | Activation descriptors; whole-tree input and help conformance | Complete (window D1 PASS) | 3-8 |
| 9a | Grammar failures and unknown commands across the live tree (DEC-280) | Complete (window D1 PASS) | 9 |
| 10 | Human failure lines in every family; generic-fallback messages (DEC-281) | Complete (`a6b6203`) | 9a |
| 11 | Readable Stop, Attention and repository lines; Human fixture generator | Complete (`93a1f792`; window D2 PASS, review-10) | 10 |
| 12 | Concern, Deviation and Human Outcome `get` leaves (DEC-274) | Complete (`fc4bdff7`) | 4, 10 |
| 13 | Judgment disclaimers once per section; recorded reasoning shown | Complete (`fd6783f9`, DEC-289) | 11, 12 |
| 14 | Human-takes-control journey and phase conformance; `resume` Issue position (DEC-290) | Resume attempt 2 (DEC-290) | 1-13, 9a |

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
| Mid-phase reviews | A after 3 (review-2 PASS); B after 6 (review-4 PASS); C late, after 9 (review-6 PASS, DEC-284); D1 after 9a (review-8 PASS); D2 after 11 (review-9 BLOCK; review-10 PASS, DEC-288); E after 14 (also VAL-068, DEC-290) | Six windows per QA plan section 8 (DEC-280/282). Human walkthrough after task 11 |
| Human fixtures | `TestP24HumanFixtures` (task 11, extended by task 14) | D2 interim walkthrough recorded (`a23b9f33`); phase-end pending |
| Final review / Human checklist | `../phases/phase-24/qa/` | Pending |

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
  budget 600 s) and 337.6 s at Task 14 attempt 1; reports watch it.
- **Scratch space:** `/tmp` is a 3.9 GB tmpfs (1.3 GB free at DEC-290),
  mostly held by other projects. A gate can fail for environmental reasons
  (CF-042).

## Next Planning Checkpoint

Task 14 closes under DEC-290, then window E and the phase-end review.
Phase 25 "Parallel and evolving work" follows.
