# Phase 5 Human Verification Checklist

> Phase: 5 - Mutation + Workflow Transitions
> Date: 2026-06-13
> Tester: MBEL
>
> Include only `(human)` items here: perception, product-owner judgment,
> dogfooding, tone, and flow coherence. Deterministic `(qa-lead)` checks belong
> in `final-review.md` with reproducible evidence. Human checks that validate a
> phase assertion cite the relevant `VAL-NNN`.

## Dogfood Session

One observable per row. `VAL-013` is fanned out from the QA plan as
`HUMAN-013a` through `HUMAN-013f`.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-013a | VAL-013 | In a real `.fid/` workspace, update issue fields used in normal work and judge whether the command flow and confirmation are clear enough to trust. | The update command feels understandable, the confirmation is clear, and any confusing output is recorded. | PASS | update fields tried: title, priority, owner_role, owner_actor, est_tokens, due_date. All ok, but the output formatting would be clearer and easier to read if each updated field was formatted on a newline rather than a single long comma-separated list like current: `updated FID-f009f8: title "Phase 5 human dogfood phase-5-human-20260613T192719Z" -> "Phase 5 human dogfood updated phase-5-human-20260613T192719Z", priority 2 -> 1, est-tokens (none) -> 1234, due-date (none) -> "2026-06-30", owner-role "developer" -> "qa_lead", owner-actor (none) -> "human:MBEL"`|
| HUMAN-013b | VAL-013 | In the same workspace, advance at least one allowed workflow edge and judge whether the transition flow communicates the new state clearly. | The transition output makes the new status/category understandable, and the edge tried is recorded. | PASS | allowed edge tried: todo -> in_progress |
| HUMAN-013c | VAL-013 | Attempt an invalid transition and judge whether the rejection names the current/requested states and gives an actionable next step. | The error identifies both states and gives enough information to decide what to do next. | PASS | invalid edge tried: todo -> done. @tech-lead: does the spec support the ability to mark a task as done/closed without completing it? For example, a backlog for a product may include low priority tasks that don't get worked on for long enough that some other repo refactor etc. occurs making these tasks redundant. In this case, we do need the ability for the workflow to allow `todo`->`done` but we'd need some `resolution_type=wont_complete` or something of that kind (I know Jira has something like this.) |
| HUMAN-013d | VAL-013 | Attempt a rejected update, such as invalid priority or unknown owner role, and judge whether the rejection explains what to fix. | The error explains the invalid field or role and what correction is needed. | PASS | rejected update tried: priority 8, optionally owner-role not_a_role |
| HUMAN-013e | VAL-013 | After the rejected update and rejected transition, run a later valid read or unrelated command and judge whether the workspace still behaves normally. | Later reads or unrelated commands still work; no rejected write appears to poison the journal or block work. | PASS | later reads/list command still worked: yes/no |
| HUMAN-013f | VAL-013 | Inspect dogfood friction from the combined mutate/transition session and decide whether any error wording or workflow behavior should block phase close despite deterministic PASS results. | Product-owner judgment is pass, fail, or partial with specific wording or workflow concerns recorded. | PASS | But see notes above about improvements. |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Dogfood Session | 6 | 6 | 0 | 0 | 0 |
| **Total** | 6 | 6 | 0 | 0 | 0 |

| | |
|---|---|
| **Overall result** | PASS |
| **Date completed** | 2026-06-14 |
| **Notes** | See notes above for improvement items. |

<!--
Outcome values: pass, fail, partial, blank.
Overall result: pass if all pass; partial if any partial and no fail; fail if
any fail.
Granularity: one observable per row. Dense comma-listed rows are imprecise and
let items get missed.
-->
