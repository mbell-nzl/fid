# Phase 5 Human Verification Checklist

> Phase: 5 - Mutation + Workflow Transitions
> Date: 2026-06-13
> Tester: _______________
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
| HUMAN-013a | VAL-013 | In a real `.fid/` workspace, update issue fields used in normal work and judge whether the command flow and confirmation are clear enough to trust. | The update command feels understandable, the confirmation is clear, and any confusing output is recorded. | | |
| HUMAN-013b | VAL-013 | In the same workspace, advance at least one allowed workflow edge and judge whether the transition flow communicates the new state clearly. | The transition output makes the new status/category understandable, and the edge tried is recorded. | | |
| HUMAN-013c | VAL-013 | Attempt an invalid transition and judge whether the rejection names the current/requested states and gives an actionable next step. | The error identifies both states and gives enough information to decide what to do next. | | |
| HUMAN-013d | VAL-013 | Attempt a rejected update, such as invalid priority or unknown owner role, and judge whether the rejection explains what to fix. | The error explains the invalid field or role and what correction is needed. | | |
| HUMAN-013e | VAL-013 | After the rejected update and rejected transition, run a later valid read or unrelated command and judge whether the workspace still behaves normally. | Later reads or unrelated commands still work; no rejected write appears to poison the journal or block work. | | |
| HUMAN-013f | VAL-013 | Inspect dogfood friction from the combined mutate/transition session and decide whether any error wording or workflow behavior should block phase close despite deterministic PASS results. | Product-owner judgment is pass, fail, or partial with specific wording or workflow concerns recorded. | | |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Dogfood Session | 6 | 0 | 0 | 0 | 6 |
| **Total** | 6 | 0 | 0 | 0 | 6 |

| | |
|---|---|
| **Overall result** | |
| **Date completed** | |
| **Notes** | |

<!--
Outcome values: pass, fail, partial, blank.
Overall result: pass if all pass; partial if any partial and no fail; fail if
any fail.
Granularity: one observable per row. Dense comma-listed rows are imprecise and
let items get missed.
-->
