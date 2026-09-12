# Progress - Current

## Project State

Phase 22 is **in Human acceptance repair (task-19)** on `phase-22` under DEC-188. The planning base is
merged main `89d964a84dbdd297fdf472d6f78b001e2efac702`. Phase 21 completed
under DEC-179 on 2026-09-04 (87/87 VAL, 4/4 Human), and its branch has since
been merged.

All 18 planned tasks, six review windows and the independent task-12
checkpoint are complete. Original completion reports and causal allocations
remain in the phase archive and run ledger. Final-review-2 independently closed
FQ-02 through FQ-07. Final-review-3 found the remaining accepted force-recovery
intervention accounting missing; DEC-218 allocated its completion to Task 1.

Root completed that repair at c361e9e and fixture follow-up 671e4e4, with report
ab3c255 and clean review checkpoint 3832df1. All twelve unchanged task gates
passed: race CLI 414.292s / export 270.522s, coverage 65.4%, integration 557.259s, and
100000-event rebuild 27490ms below 30000ms. Every failed attempt is preserved.
The first integration run timed out at 600.834s; isolated copies of an immutable
valid prefix reduced the new force test group from 89.800s to 35.396s without removing
cases or rollback comparisons. Independent final-review-4 closes FQ-01 and
passes every technical phase check.
The overall review remains BLOCK solely for the three Human-owned assertions.
Human attempt 1 recorded 1 PASS, 2 FAIL and 3 PARTIAL outcomes. DEC-219 and
task-19 allocate one serialized guidance repair and independent review before
a fresh Human retest. No escalation is open.

The Human delegated completion of tasks 8-18, required reviews, final
deterministic QA, and runnable Human checklist preparation on 2026-09-07.
See `../phases/phase-22/delegate-handoff.md`. Human attempt 1 failed; a retest is required;
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
product/spec decision is pending.

## Phase 22 Pack

Heavy mode approved; 18 sequential tasks and 80 Ready validation assertions.
Exact slice/wire/DDL/Event contracts, executable-scope matrix, V1 convergence
inventory, QA Lead handoff, decisions and triage log are in
`../phases/phase-22/`. `phase_brief.md` indexes every task.

Execution started on 2026-09-05 under DEC-188. QA's plan and consistency
assessment remain in `qa/qa-plan.md` and `qa/consistency-review.md`; the
refinement authority is DEC-184..187 and `qa/tl-refinement.md`. The ledger
records completed tasks 1..18 and the task-8 STOP resolved by DEC-195. Product QA
receipts retain their original session attribution; the separate maintenance
report records freshly executed workflow checks only.

## Approval State

- Mode, contracts, task list and refinements: Human approved under DEC-188.
- QA plan and cross-product review: assessment recorded; TL refinements applied.
- Task statuses: original 1..18 complete; all six windows/checkpoint passed.
  The seven initial final-QA findings are independently closed: FQ-02..07
  by final-review-2 and FQ-01 by final-review-4.
- Next: the Human runs the [six checklist checks](../phases/phase-22/qa/human-checklist.md)
  for VAL-076..078, then QA/Human results are recorded before phase acceptance
  and Tech Lead wrap-up. [Final review](../phases/phase-22/qa/final-review-4.md)
  records technical qualification and the remaining Human-only block.
- Fresh fixtures: checkpoints-4132086272, clean source 3832df1, all twelve
  inspection commands independently verified. Successor lease expires
  **13 September 2026 at 12:19:09 NZST**; regenerate and requalify if needed.
- DEC-216 records the Human-approved typed snapshot contract. DEC-217 keeps
  archive ceilings unchanged while splitting small routine checks from the
  explicit production-scale phase gate. DEC-218 completes accepted force
  intervention accounting; it introduces no new product policy or threshold.
- Implementation and independent technical review are complete. Human judgments
  remain pending; no product decision or implementation repair is outstanding.

## Carry-Forward

The register below retains its pre-wrap-up dispositions. The complete trigger
audit and closure/re-targeting follow Human verification at Tech Lead wrap-up.

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
