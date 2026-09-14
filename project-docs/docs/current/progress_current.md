# Progress — Current

Updated 2026-09-14. **Phase 22 closed; Phase 23 planned and awaiting approval.**
Branch: `phase-22` (PR open, CI green); planning base for Phase 23 is
`0195897a347dfe1872db0f7be3800c5c3ccdefcf` on `main`.

## Phase 23 planning — 2026-09-14

Tech Lead drafted the full Heavy-mode plan for the roadmap's "Review and Human
judgment" slice: slice contract, wire amendment, DDL and Event contract,
execution scope, V1 convergence inventory, a clause-to-task coverage matrix, a
73-assertion validation contract, and an ordered 18-task pack. Proposed DEC-222
requests acceptance of that set as one unit. Nothing is approved; no product
code, migration, or test for the slice may start yet. Two open questions go to
the Human: adding `fid concern dispose` to the accepted CLI surface, and
restricting `links` to Issue endpoints so requirement traceability has one
authoritative model. New carry-forward CF-041 records that hosted CI does not
run the tagged integration suite or the rebuild benchmark. See
[phase brief](phase_brief.md) and
[Phase 23 decisions](../phases/phase-23/decisions.md).

## Completed

- All 18 planned tasks and supplemental Human-guidance task 19.
- Six review windows, task-12 checkpoint, all seven final-QA findings and the
  independent task-19 review, including its MF-01 repair.
- 77 QA-owned assertions plus three Human-owned assertions: 80/80 PASS.
- 84 deterministic QA items; second Human checklist 6/6 and overall PASS.
- DEC-220 sign-off, canonical spec sync, decision/deviation/TRG/ESC audit,
  carry-forward trigger audit and current-context compression.
- CF-024/034/035 closed; CF-038 documentation half closed. CF-039 formally
  indexed; CF-040 records the Human's nonblocking CLI polish/verbosity request.

## Evidence and remaining boundaries

[Final acceptance](../phases/phase-22/qa/final-acceptance.md) attributes the final
outcomes; [wrap audit](../phases/phase-22/wrap-up-audit.md) records dispositions.
[Phase progress](../phases/phase-22/progress.md) retains execution history and
[decisions](../phases/phase-22/decisions.md) retain all authorities and workflow notes.
Phase 22's original implementation and Human acceptance remain closed. The
post-wrap PR CI follow-up is task 20 under DEC-221; the Phase 23 plan above is
drafted but unapproved.
Separate workflow retrospective remains CF-036/039. Claims and external adapter
qualification remain in their approved Phase 24/25 owners.

## Post-wrap performance repair — 2026-09-14

PR #24's original test step failed after 31m33s, with two package timeouts.
The supplied failure and earlier unapproved timeout proposal remain in the
[diagnosis](../phases/phase-22/qa/pr-ci-diagnosis.md). DEC-221/task 20 implements
the chosen performance repair; ESC-22-005's decision request is resolved with
no timeout increase.

All seven local gates passed, including full race coverage (65.5%) and build.
Matching uncached before/after suites retain all 2,653 original run cases and
add two isolation regressions. CLI 541.652s -> 483.994s; export 368.615s ->
264.146s; workflow 503.550s -> 464.948s (unchanged tests, less contention).
Four exhaustive CI shards replace duplicate full runs. Controller tests include
real Go execution and coverage reconciliation; independent task-20 review PASS.
See [task 20 report](../phases/phase-22/tasks/task-20-done.md).
Hosted run34805256766 on6b8fe50 passed all six jobs in13m53s. All724 targets
and12967 statement blocks reconciled;65.5% coverage. Summed job time47m48s
(about50 rounded minutes). Task20 and CI repair COMPLETE; PR ready for Human
merge. [Hosted evidence](../phases/phase-22/qa/task-20-hosted-ci.md) binds the
verified PR head. Closure records are local documentation only; root does not
push or merge.
