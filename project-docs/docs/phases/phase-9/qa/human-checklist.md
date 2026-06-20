# Phase 9 Human Verification Checklist

> Phase: 9 - Concerns + v0.1c closure
> Date: 2026-06-20
> Tester: _______________
>
> Include only `(human)` items here: perception, product-owner judgment,
> dogfooding, tone, and flow coherence. Deterministic `(qa-lead)` checks belong
> in `final-review.md` with reproducible evidence. These rows validate VAL-011
> only after a human records outcomes from a real `.fid/` workspace.

## Concern Lifecycle Dogfood

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-011a | VAL-011 | Raise a concern against a real completion and confirm the visible open concern has the expected summary, role, issue link, and completion link `(human)` | The concern is visible as open and reads as the concern raised against the intended issue/completion. | | |
| HUMAN-011b | VAL-011 | After raising the concern, confirm issue reads make the open concern obvious enough through `concerns_raised` and `open_concerns_count` `(human)` | `fid issue get` and `fid issue list` make unresolved concern work visible before/after raise. | | |
| HUMAN-011c | VAL-011 | Address the concern from a later completion and confirm the lifecycle reads as an in-place resolution, not a new unrelated concern `(human)` | The address flow reads as resolving the same concern and does not imply a duplicate concern was created. | | |
| HUMAN-011d | VAL-011 | After addressing the last open concern, confirm issue reads no longer imply unresolved concern work remains `(human)` | `fid issue get` and `fid issue list` no longer suggest open concern work remains. | | |
| HUMAN-011e | VAL-011 | Use `fid concern list` with issue and status filters and confirm the filtered views match the way a developer would look up concern state during dogfooding `(human)` | Filters are useful, empty-result behavior is understandable, and output does not confuse the concern state. | | |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Concern Lifecycle Dogfood | 5 | 0 | 0 | 0 | 5 |
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
