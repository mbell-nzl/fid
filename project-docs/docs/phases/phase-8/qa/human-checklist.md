# Phase 8 Human Verification Checklist

> Phase: 8 - Task Completions
> Date: 2026-06-20
> Tester: _______________
>
> Include only `(human)` items here: perception, product-owner judgment,
> dogfooding, tone, and flow coherence. Deterministic `(qa-lead)` checks belong
> in `final-review.md` with reproducible evidence. Human checks that validate a
> phase assertion cite the relevant `VAL-NNN`.

## Completion Lifecycle Dogfood

One observable per row for `VAL-008`.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HV-08a | VAL-008 | In a real `.fid/` workspace, record a completion attempt and inspect the shown/read-back fields: role, tokens, commits, files, and notes path. | The completion record and detail output read as expected, with no confusing or missing field presentation. | | |
| HV-08b | VAL-008 | Record a second completion attempt for the same task and inspect attempt numbering plus prior-attempt supersession. | Attempt 2 and the prior `superseded` attempt are understandable and match the expected lifecycle. | | |
| HV-08c | VAL-008 | Inspect a completion that includes a gate snapshot. | The rolled-up gate result, pass/fail counts, and gate detail presentation read as expected. | | |
| HV-08d | VAL-008 | Record a completion with human intervention and an intervention note, then inspect the completion output/read-back. | The intervention flag and note presentation are clear and expected. | | |
| HV-08e | VAL-008 | Exercise a real `has_completion_report`-gated transition before and after recording a completion. | The reject/pass lifecycle behaves as expected and the rejection is understandable before the completion exists. | | |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Completion Lifecycle Dogfood | 5 | 0 | 0 | 0 | 5 |
| **Total** | 5 | 0 | 0 | 0 | 5 |

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
