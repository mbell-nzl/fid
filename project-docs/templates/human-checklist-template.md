# Phase <N> Human Verification Checklist

> Phase: <N> - <Phase Name>
> Date: <YYYY-MM-DD>
> Tester: _______________
>
> Include only `(human)` items here: perception, product-owner judgment,
> dogfooding, tone, and flow coherence. Deterministic `(qa-lead)` checks belong
> in `final-review.md` with reproducible evidence. Human checks that validate a
> phase assertion cite the relevant `VAL-NNN`.

## <Section Name>

One observable per row. If a `(human)` `VAL-NNN` spans several surfaces,
affordances, CLI flows, or output formats, split it into one row per discrete
observable. Use sub-IDs (`HV-01a`, `HV-01b`, ...) that all cite the same parent
`VAL-NNN`.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HV-01a | VAL-NNN or no VAL | <single manual observable> | <expected observation> | | |
| HV-01b | VAL-NNN or no VAL | <next manual observable> | <expected observation> | | |

## <Section Name>

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| YY-01 | VAL-NNN or no VAL | <Manual step> | <Expected observation> | | |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| <Section Name> | | | | | |
| **Total** | | | | | |

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
