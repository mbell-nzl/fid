# Progress — Current

Updated 2026-09-20. **Phase 23 Tasks 1-3 complete; Task 4 recovery authorized
under DEC-239 after ESC-23-003's protected evidence contradiction is resolved.**
Branch: `phase-23`. The Runner must reconcile the retained Task 4 STOP; Task 4
is not complete. The Human authorizes execution through Task 9 and the window
7-9 review. Task 9 still needs bounded DEC-237 QA confirmation; Tasks 10-11 stay
held.

## Task 4 recovery authority - 2026-09-20

DEC-239 preserves Run/Claim Actor binding and durable history, and explicitly
amends VAL-015's unreachable Run-only CLI fixture to isolated query evidence.
Real CLI Claim rejection and same-Reviewer rework remain required. Inspection
also found the current implementation incorrectly rejects the same independent
Actor resolving their earlier Review's Finding. One bounded recovery must fix
that, prove multiple prior Findings resolved atomically, provide the required
built-binary evidence and close the report's VAL-003 UNVERIFIED row before
recording completion. Historical failures and the eight attempts remain intact.
See [resume handoff](../phases/phase-23/qa/resume-after-esc-23-003.md).

DEC-240 dispositions the bounded Task 9 readiness finding R-29: split retained
Condition validation, recorded-history replay and Attention-key stability into
VAL-082/083/084. All cases and task commands remain; the contract now has 84
assertions. Independent QA confirms this final diff before Task 9 dispatch.

## Task 1 recovery — 2026-09-20

ESC-23-002 is **RESOLVED 2026-09-20 under DEC-238**. The Human delegated
bounded repair authority above the fixed agent roles. Rebuild statement reuse
and fresh/isolated test setup repair the causes without changing any gate or
assertion. Fresh exact gates pass: CLI **527.099s**, workflow **231.010s**,
tagged integration **669.586s**; rebuild samples **22.728s / 22.623s / 22.489s**.
Integration binary directories stay **0 → 0**. Task 1 is complete in its report;
the Runner must reconcile the retained STOP before advancing. Tasks 9-11 stay
held. See [repair evidence](../phases/phase-23/qa/esc-23-002-repair.md) and
[resume handoff](../phases/phase-23/qa/resume-after-esc-23-002.md).

The original benchmark, expired-fixture assertion and runtime-budget failures
remain durable evidence. Thresholds stay 30s / 600s / 900s. Prepared rebuild
writes retain per-Event transactions and rollback; cached test setup creates
fresh expiring state and preserves per-case isolation. No later task dispatched.

## Phase 23 planning — 2026-09-14

Tech Lead planned the roadmap's "Review and Human judgment" slice in Heavy mode
and the Human accepted it under DEC-222: slice contract, wire amendment, DDL and
Event contract, execution scope, V1 convergence inventory, a clause-to-task
coverage matrix, a validation contract, and an ordered 18-task pack,
plus six canonical spec amendments applied by their owning tasks.

Both open planning questions were resolved before acceptance. Concern
dispositions are exactly `carried-forward`, `rejected`, and `superseded`, with
self-disposal permitted and both Actors recorded; `links` is restricted to Issue
endpoints so `requirement_issue_links` is the single traceability model, keeping
the new `verifies` relation. Planning also found a third spec gap: the canonical
Event family list had no `Concern disposed` family, added under the DEC-195 rule.

New carry-forward CF-041 records that hosted CI does not run the tagged
integration suite or the rebuild benchmark.

## Phase 23 QA first pass — 2026-09-15

QA Lead authored the Heavy-mode plan with fresh command and feasibility
validation and returned it **NOT READY** with thirteen Must-fix revisions and
sixteen improvements: a migration that fails wherever Evidence has Artifact
links, a Stop no Event could raise, one input with two accepted codes, an
ambiguous independence rule, unnamed code producers, a denial replay with no
denial Event, four assertions whose owning task could not produce their
evidence, two unexecutable file scopes, unowned Human fixtures, a bundled
assertion, three unasserted clauses, and 560 s of test time against an implicit
600 s limit.

Tech Lead accepted all thirteen and applied them. DEC-223..227 amend accepted
contracts; DEC-228 carries the rest under planning authority. The Human approved
all five on 2026-09-19.

Reviewing DEC-226 surfaced a root cause worth recording: the project had one
name for two closed sets. `GLOSSARY.md` defined `Reason code` as "why an
operation failed" while the slice contracts used `ReasonCode` for the
evaluation-reason set, and the implementation registry carries the same
collision - which is how `finding-unresolved` came to be filed as a failure
code. The glossary now defines **Failure code** and **Evaluation reason**
separately; task 13 renames the registry to match. The validation contract grew from 74 to 79 assertions
and every row is now `Ready`. New carry-forward CF-042 records the leaked
integration binary directories.

The QA Lead's 2026-09-19 readiness refresh re-proved all five approved
decisions against fresh evidence - DEC-223 on the real migration chain under
`modernc.org/sqlite`, DEC-227's budgets re-measured at the new command text,
and all four DEC-224 raising sequences reachable - and returned two Must-fix
revisions. DEC-229 dispositions both: four closed registries
(`internal/workflow/stop.go`, `current_event_value_catalog.go`,
`current_event_columns.go`, `archive_semantics_catalog.go`) had no owning task
and now have one, with a standing rule that a journal Event rule change is
mirrored in the archive catalog in the same task; and CF-042's recorded cause
was wrong - the leak is per-run, not per-kill, so task 1's mitigation becomes
binary reuse rather than age-based pruning. Pinning the `*-not-found` producers
then showed `decision-not-found` has no producing leaf at all, and DEC-230
removes it. Two QA findings were not accepted as written: task 1 already
registers the nine Event rules, so tasks 4 and 5 gain no catalog file, and the
column catalog goes to task 1 rather than five domain tasks.

Six further readiness passes followed. DEC-231/232/233 closed the journal
strict-path registry set after three independent re-derivations; ESC-23-001
raised the `Condition` encoding to the Human, who chose Option A (DEC-235); and
DEC-236 dispositioned pass 7, whose R-24 caught a correct prohibition with no
statement of what must still **run** - the literal reading deleted the Phase 22
nullable binding with every gate green. Pass 8 ran twice independently and both
runs converged on VAL-078. DEC-237 splits it into resolution (VAL-078) and
Condition encoding (VAL-080), and gives the two DEC-235 preservation invariants
an ID each (VAL-081, VAL-082): live behavior this refactor could delete with no
other detector. The contract is 82 rows and nothing was weakened.

Task 1's delegated repair was completed and the Runner advanced through Task 3.
The current Task 4 STOP is resolved into recovery authority under DEC-239 above.
The Human has cleared Task 9; DEC-237's bounded QA confirmation still applies,
and Tasks 10-11 remain held. See
[phase brief](phase_brief.md) and [Phase 23 decisions](../phases/phase-23/decisions.md).

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
post-wrap PR CI follow-up is task 20 under DEC-221; Phase 23 is approved above
with all five amendments approved.
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
