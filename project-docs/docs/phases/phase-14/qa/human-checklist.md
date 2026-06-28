# Phase 14 Human Verification Checklist

> Phase: 14 - Telemetry + Pricing
> Date: 2026-06-28
> Tester: MBEL
>
> Include only `(human)` items here: perception, product-owner judgment,
> dogfooding, tone, and flow coherence. Deterministic `(qa-lead)` checks belong
> in `final-review.md` with reproducible evidence. Human checks that validate a
> phase assertion cite the relevant `VAL-NNN`.

## Dogfood Observables

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-14a | VAL-009 | In a real `.fid` workspace, record an agent run with a known priced model and inspect the CLI output `(human)`. | The output makes it clear that an agent run was recorded and shows a computed USD cost. | PASS | n/a |
| HUMAN-14b | VAL-009 | In a real `.fid` workspace, record an agent run with an unknown model and inspect the CLI output and stderr warning `(human)`. | The command succeeds, records the run, and the warning names the unknown model without implying a guessed price. | PASS | n/a |
| HUMAN-14c | VAL-009 | In a real `.fid` workspace with recorded runs, inspect `fid stats tokens` and `fid stats cost` grouped by role, model, and epic `(human)`. | The aggregate output is understandable from the CLI alone and communicates token totals, cost totals, and run counts by role, model, and epic. | PASS | n/a |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Dogfood Observables | 3 | 3 | 0 | 0 | 0 |
| **Total** | 3 | 3 | 0 | 0 | 0 |

| | |
|---|---|
| **Overall result** | PASS |
| **Date completed** | 2026-06-28 |
| **Notes** | n/a |

<!--
Outcome values: pass, fail, partial, blank.
Overall result: pass if all pass; partial if any partial and no fail; fail if
any fail.
Granularity: one observable per row. Dense comma-listed rows are imprecise and
let items get missed.
-->
