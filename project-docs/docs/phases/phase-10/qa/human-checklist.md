# Phase 10 Human Verification Checklist

> Phase: 10 - Permissions Matrix + Policy
> Date: 2026-06-21
> Tester: _______________
>
> Include only `(human)` items here: perception, product-owner judgment,
> dogfooding, tone, and flow coherence. Deterministic `(qa-lead)` checks belong
> in `final-review.md` with reproducible evidence. Human checks that validate a
> phase assertion cite the relevant `VAL-NNN`.

## Policy Matrix Dogfood

Run these checks in a real `.fid/` workspace. QA final review already validated
deterministic JSON validity, row parity, write outcomes, and DB/journal counts;
these rows ask for human readability and actionability judgment only.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-09a | VAL-009 | Run `fid policy show --format markdown` or the preferred human-readable format and inspect the role x action matrix `(human)`. | You can find a specific action and compare allowed roles without ambiguity. | PASS | n/a |
| HUMAN-09b | VAL-009 | Run `fid policy show --format json` and inspect action names plus role booleans for scanability `(human)`. | The JSON is understandable enough for an operator or agent to locate an action and read true/false role cells. | PASS | n/a |

## Denial Error Dogfood

Use a denied write in the same real `.fid/` workspace, for example
`fid issue create --kind epic --title "Human checklist denied" --scheduling sequential --by developer`.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-09c | VAL-009 | Trigger a permission-denied write and judge whether the error names the role, names the action, and says the role is not permitted `(human)`. | The message clearly identifies the denied role and action and includes not-permitted language. | PASS | n/a |
| HUMAN-09d | VAL-009 | Judge whether the denied-write error gives enough signal to choose a next step, such as using a permitted role, checking `fid policy show`, or escalating a permission mismatch `(human)`. | The message is actionable enough that the operator knows what to try or where to look next. | PASS | n/a |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Policy Matrix Dogfood | 2 | 2 | 0 | 0 | 0 |
| Denial Error Dogfood | 2 | 2 | 0 | 0 | 0 |
| **Total** | 4 | 4 | 0 | 0 | 0 |

| | |
|---|---|
| **Overall result** | PASS |
| **Date completed** | 2026-06-21 |
| **Notes** | n/a |

<!--
Outcome values: pass, fail, partial, blank.
Overall result: pass if all pass; partial if any partial and no fail; fail if
any fail.
Granularity: one observable per row. Dense comma-listed rows are imprecise and
let items get missed.
-->
