# Phase 13 Human Verification Checklist

> Phase: 13 - Labels + v0.1d Closure
> Date: 2026-06-27
> Tester: _______________
>
> Include only `(human)` items here: perception, product-owner judgment,
> dogfooding, tone, and flow coherence. Deterministic `(qa-lead)` checks belong
> in `final-review.md` with reproducible evidence. Human checks that validate a
> phase assertion cite the relevant validation ID.

## Dogfood Observables

One observable per row. Use a real `.fid` workspace and judge only what can be
understood from CLI output, not implementation files.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HV-13a | VAL-010 | Confirm from real CLI output that a label was created and listed with its name, colour, and description visible. | The output makes the created label and its metadata visible enough for a user to identify them. | | |
| HV-13b | VAL-010 | Confirm from real CLI output that adding one label to an issue left the issue's existing labels in place. | Before/after output shows the new label plus the pre-existing labels; the command does not look like a replace-all operation. | | |
| HV-13c | VAL-010 | Confirm from real CLI output that removing one label from an issue left the remaining labels in place. | Before/after output shows only the selected label removed while the other labels remain visible. | | |
| HV-13d | VAL-010 | Confirm from real CLI output that create -> claim -> transition -> complete -> unclaim ran end-to-end with a link and label on the issue. | The transcript is understandable as one complete lifecycle and shows the issue retained its link and label through the flow. | | |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Dogfood Observables | 4 | | | | 4 |
| **Total** | 4 | | | | 4 |

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
