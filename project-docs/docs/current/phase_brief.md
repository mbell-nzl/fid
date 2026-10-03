# Phase Brief - No Active Phase

## Status

Phase 23 "Review and Human judgment" is complete (DEC-271, 2026-10-03) and
merged to `main` as PR #25 (`a46ec77`). DEC-272 inserts a new Phase 24,
which the Tech Lead plans next on branch `phase-24`.

## Last Phase Summary

- Heavy mode; 18 planned tasks plus supplemental tasks 19-21 (DEC-269).
- 89/89 validation assertions (85 deterministic, 4 Human). Final review 2
  PASS; Human checklist 1 FAIL (55/2/6), checklist 2 7/7 PASS.
- Delivered: independent Review coverage and Findings, Finding dispositions,
  review-cycle budgets and a resolvable `review-rework-exhausted` Stop,
  Concerns, Deviations, Decision References with protected-Decision blocking,
  provenanced Human Outcomes, Attention resolution, requirement references
  with stable `requirement_id`, migrations 0017/0018, and V1 judgment removal.
- Qualified capabilities: `reviews.coverage.v1`, `human.outcomes.v1`.
- Records: [final acceptance](../phases/phase-23/qa/final-acceptance.md),
  [wrap audit](../phases/phase-23/wrap-up-audit.md),
  [decisions](../phases/phase-23/decisions.md).

## Next Phase: 24 "Operating experience and input convergence" (DEC-272)

Goal: a Human can take control of a stuck or diverging journey using fid
output alone. Agreed scope (full text in DEC-272):

1. **Input convergence.** fid observes repository identity itself. Inputs
   drop the `Repository` block, `project_id`, `worktree_id` and
   `capture_source`. Callers send only checked expectations: expected HEAD
   and dirty state on Review and Run inputs, and an optional file digest on
   references. Request digests cover the input as written.
2. **Human authoring (option C).** `--help` documents each `--input`;
   `--example` prints a template; guidance names the resolving command.
3. **Human output.** A4 readable Stop/Attention details; A5 `error:` plus
   `reason:`/`blocked by:`/`next:`; A6 disclaimers once per section; A7
   recorded reasoning shown (truncated in lists, full in detail views).

Expected Heavy mode (public-interface change), roughly 7-10 tasks; the mode
and task list go to the Human for approval at planning.

## Later Phases (Renumbered)

- Phase 25 "Parallel and evolving work": parallel Claims in disjoint
  worktrees, multi-Issue Envelope revision, `retry-exhausted`,
  `work.claims.v1`.
- Phase 26 "Qualification and baseline convergence", then the full
  Human-output design revisit (CF-040).

Optional, separate: the deep workflow retrospective (CF-036/039).
