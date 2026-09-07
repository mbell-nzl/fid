# Progress - Current

## Project State

Phase 22 is **executing** on `phase-22` under DEC-188. The planning base is
merged main `89d964a84dbdd297fdf472d6f78b001e2efac702`. Phase 21 completed
under DEC-179 on 2026-09-04 (87/87 VAL, 4/4 Human), and its branch has since
been merged.

Tasks 1-12 are complete; review-1b and review-2b passed. Task 8 completed the
DEC-195 cancellation DDL/filter convergence with task-1 provenance, all eight
required gates and independent report lint. Task 9 completed DEC-196/197 semantic
conformance, Human Stop text and actual-operation Stop guards with all five
final gates and independent report lint. Review 3 passed independently. Task 10
completed under DEC-198 at 03cc6c7 with exact binding/agreement/upgrade/replay
coverage, all five final gates and independent report lint. Task 11 completed
under DEC-199/200 at 30dd780 with independent fake/fid cases, real response-loss
boundary tests, all five final gates and root report lint. Task 12 completed under DEC-196/201 at c2d8098, with final gate addendum
833b4fc and preserved report f57b4aa. Review 4 and the independent checkpoint
are executing; task 13 is prepared but has not started. No escalation is open.

The Human delegated completion of tasks 8-18, required reviews, final
deterministic QA, and runnable Human checklist preparation on 2026-09-07.
See `../phases/phase-22/delegate-handoff.md`. Human judgments remain pending;
this delegation does not claim phase sign-off or release readiness.

Earlier recovery history: task 6 entered a same-task
continuation: its real built-binary journey reached the absent `fid run inspect`
leaf and returned `RECOVERY_REQUIRED`, raising TRG-22-008 / ESC-22-003. DEC-194
resolves it. Task 4 closed COMPLETE without writing the `run.checkpoint`
producer or `run.inspect` read API that its own file scope, `schema.md` section
5, `public-interface.md` section 3, and `contracts.md` section 3.1 all assign to
it, so the owner was never ambiguous and this is undelivered in-scope work, not
a protected boundary. Task 4 is closed with two successors on top and the
Runner's only recovery route is a same-task continuation, so DEC-194 moves
execution of exactly four `internal/workflow` files into task 6 with task-4
provenance recorded. Task 4 stays COMPLETE and unedited. ESC-22-001 (DEC-190/
191) and ESC-22-002 (DEC-192/193) closed earlier. No escalation is open and no
Human decision is pending.

## Phase 22 Pack

Heavy mode approved; 18 sequential tasks and 80 Ready validation assertions.
Exact slice/wire/DDL/Event contracts, executable-scope matrix, V1 convergence
inventory, QA Lead handoff, decisions and triage log are in
`../phases/phase-22/`. `phase_brief.md` indexes every task.

Execution started on 2026-09-05 under DEC-188. QA's plan and consistency
assessment remain in `qa/qa-plan.md` and `qa/consistency-review.md`; the
refinement authority is DEC-184..187 and `qa/tl-refinement.md`. The ledger
records completed tasks 1..12 and the task-8 STOP resolved by DEC-195. Product QA
receipts retain their original session attribution; the separate maintenance
report records freshly executed workflow checks only.

## Approval State

- Mode, contracts, task list and refinements: Human approved under DEC-188.
- QA plan and cross-product review: assessment recorded; TL refinements applied.
- Task statuses: 1..12 complete; reviews 1b/2b/3 passed; review 4/checkpoint executing; 13..18 pending.
- Next: complete review window 4 and the independent task-12 checkpoint
  before task 13.

## Carry-Forward

CF-024: task 17/final QA, unchanged reference rebuild rule; early samples at 1/16.
CF-034: direct workflow/evidence branch tests and before/after coverage in
tasks 2/7. CF-035: config-consuming-path proposal in task 11.
CF-033: kickoff inventory recorded; AGENTS edit awaits separate maintenance.
CF-018: exact production sqlquery.go trigger retained.
Other CFs re-deferred with rationale in DEC-182; CF-036 schedules separate
workflow identity/registration maintenance planning from DEC-183.
CF-037: external kit contract/adapter duties required for Phase 25 qualification.
R-14: task 8 owns process fixture repair; earlier-failure route is DEC-186.
CF-038: live-output golden hardening and post-truncation `output_bytes`
wording, both from DEC-189; neither blocks this phase.
CF-039: a task file scope is not self-enforcing (DEC-194). Planning must give
every named production path an acceptance criterion or VAL citation that fails
when it is absent; wrap-up should decide whether report-lint or QA review should
also diff a report's delivered files against its accepted file scope.

## References

- Phase 22 planning/history: `../phases/phase-22/decisions.md`, `progress.md`.
- Phase 21 closure: `../phases/phase-21/decisions.md` DEC-179 and
  `qa/final-review-2.md`.
