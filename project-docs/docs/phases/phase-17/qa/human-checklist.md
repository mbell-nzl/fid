# Phase 17 Human Verification Checklist

> Phase: 17 - Maintenance, Stub Guards, and v0.1 Closure
> Date: 2026-07-17
> Tester: _______________
>
> This checklist contains only Human-owned dogfood judgments for VAL-025.
> Do not execute it until the deterministic VAL-024 blockers recorded in
> `final-review.md` have been cleared and the phase-end gate has been rerun.

## Real-Work v0.1 Dogfood

One observable is recorded per row. Leave the outcome blank until a Human has
performed the check during real project work. Valid outcomes are `pass`,
`fail`, or `partial`; include a disposition reference for any result other
than `pass`.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-25a | VAL-025 | Use `fid maintenance check --format markdown` during real project work and judge whether its findings explain what needs attention `(human)` | The findings are clear enough to identify the affected entity and the next useful action without reading implementation code. | | |
| HUMAN-25b | VAL-025 | Invoke one guarded v0.2 governance command and judge whether the version-gate error clearly says what ships when `(human)` | The first error line makes the unavailable entity and its shipping version immediately understandable. | | |
| HUMAN-25c | VAL-025 | Run `fid ready` on real project state and judge whether the selected work is sensible `(human)` | The selected issue, or the explanation that none is ready, agrees with the Human's understanding of current priority and blockers. | | |
| HUMAN-25d | VAL-025 | In a disposable workspace copy, follow the documented maintenance/rebuild recovery path and judge whether the guidance is understandable `(human)` | The Human can identify and follow the recovery path without risking the real workspace or needing undocumented knowledge. | | |
| HUMAN-25e | VAL-025 | After HUMAN-25a through HUMAN-25d, accept v0.1 for continued early dogfood or explicitly disposition the remaining friction `(human)` | A product-owner decision is recorded; any non-acceptance names a follow-up decision, issue, or carry-forward reference. | | |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Real-Work v0.1 Dogfood | 5 | 0 | 0 | 0 | 5 |
| **Total** | **5** | **0** | **0** | **0** | **5** |

| | |
|---|---|
| **Overall result** | |
| **Date completed** | |
| **Notes** | |

<!--
Outcome values: pass, fail, partial, blank.
Overall result: pass if all pass; partial if any partial and no fail; fail if
any fail.
-->
