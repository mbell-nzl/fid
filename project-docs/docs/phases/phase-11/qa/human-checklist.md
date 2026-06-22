# Phase 11 Human Verification Checklist

> Phase: 11 - Claims Layer
> Date: 2026-06-21
> Tester: MBEL
>
> Include only `(human)` items here: perception, product-owner judgment,
> dogfooding, tone, and flow coherence. Deterministic `(qa-lead)` checks belong
> in `final-review.md` with reproducible evidence. Human checks that validate a
> phase assertion cite the relevant `VAL-NNN`.

## Claim Lifecycle Dogfood

One observable per row. Run these in a real `.fid` workspace used for normal
project work, not by reading implementation details.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HV-01a | VAL-014 | Identify from CLI output which role currently holds an issue claim. | The holder role is clear from claim-related CLI output or read-back. | PASS | n/a |
| HV-01b | VAL-014 | Identify from CLI output when the current issue claim expires. | The expiry time is understandable without code inspection. | PASS | n/a |
| HV-01c | VAL-014 | Release an issue claim and confirm from CLI output or read-back that it is released. | The release state is clear and the issue no longer appears to have a current claim. | PASS | n/a |
| HV-01d | VAL-014 | Try to advance an issue held by `human` as an agent role and judge whether the block reason is clear. | The error explains that the live claim is human-held and indicates why the agent cannot advance it. | PASS | n/a |
| HV-01e | VAL-014 | Dogfood the claim lifecycle for normal work and judge whether claim, renew, unclaim, and list messages are usable for agent coordination. | Product-owner/dogfood judgment is accept or request-changes with notes. | PASS | n/a |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Claim Lifecycle Dogfood | 5 | 5 | 0 | 0 | 0 |
| **Total** | 5 | 5 | 0 | 0 | 0 |

| | |
|---|---|
| **Overall result** | PASS |
| **Date completed** | 2026-06-22 |
| **Notes** | n/a |

<!--
Outcome values: pass, fail, partial, blank.
Overall result: pass if all pass; partial if any partial and no fail; fail if
any fail.
Granularity: one observable per row. Dense comma-listed rows are imprecise and
let items get missed.
-->
