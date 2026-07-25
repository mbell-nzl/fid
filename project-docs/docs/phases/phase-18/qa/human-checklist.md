# Phase 18 Human Verification Checklist

> Phase: 18 - Workflow Tables Migration + Engine Data-Source Swap + Gate Seams
> Date: 2026-07-23
> Tester: MBEL
>
> Run on a disposable copy of a v0.1 workspace, never the sole real workspace.
> Final deterministic QA passed on 2026-07-23. Execute this Human-only
> checklist to complete Phase 18 product/dogfood sign-off.

## Behavior-Preservation Dogfood

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-01a | no VAL | Judge whether representative allowed and rejected transitions feel behaviorally unchanged after migration `(human)`. | Allowed transitions succeed and rejected transitions remain equally clear and unsurprising. | PASS | n/a |
| HUMAN-01b | no VAL | Judge whether issue get/list/ready presentation and ordering feel unchanged and sensible `(human)`. | Categories, selected issues, and ordering match the tester's v0.1 expectations. | PASS | n/a |
| HUMAN-01c | no VAL | Judge whether `fid maintenance check` remains clear and actionable on the migrated disposable copy `(human)`. | Output is understandable, points to actionable states, and does not create uncertainty about workspace health. | PASS | n/a |

## Phase Acceptance

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-01d | no VAL | Decide whether the demonstrated backup/copy, migration, rebuild, and remediation evidence is sufficient to accept Phase 18 for continued dogfood `(human)`. | Tester records an explicit accept or request-changes decision with rationale. | PASS | n/a |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Behavior-Preservation Dogfood | 3 | 3 | 0 | 0 | 0 |
| Phase Acceptance | 1 | 1 | 0 | 0 | 0 |
| **Total** | **4** | **4** | **0** | **0** | **0** |

| | |
|---|---|
| **Overall result** | PASS |
| **Date completed** | 2026-07-25 |
| **Notes** | n/a |

<!--
Outcome values: pass, fail, partial, blank.
Overall result: pass if all pass; partial if any partial and no fail; fail if
any fail.
Granularity: one observable per row. Dense comma-listed rows are imprecise and
let items get missed.
-->
