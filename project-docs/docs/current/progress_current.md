# Progress — Current

Updated 2026-09-15. **Phase 22 closed; Phase 23 approved, QA-planned, not
started; five contract amendments await Human approval.**
Branch: `phase-23` at `0195897a347dfe1872db0f7be3800c5c3ccdefcf`; the `phase-22`
PR remains open with green CI.

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
integration binary directories, with bounded pruning added to task 1.

Execution has not started and nothing is dispatched. Next: the QA Lead
readiness refresh, then the Runner's readiness check. See
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
