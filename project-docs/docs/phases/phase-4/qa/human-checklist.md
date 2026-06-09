# Phase 4 Human Verification Checklist

> Phase: 4 - v0.1b Issue authoring and discovery
> Date: 2026-06-09
> Tester: _______________
>
> Include only `(human)` items here. Deterministic checks are recorded in
> `final-review.md`. Every row maps to VAL-010.

## Dogfood Loop

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HV-01a | VAL-010 | Using the built binary, create one epic, one task under that epic, and one subtask under that task during a realistic local workflow. | The create flow is understandable, the required flags are discoverable enough, and the success messages make the created IDs usable for the next command. | | |
| HV-01b | VAL-010 | Run `fid issue get <id>` on one created issue and inspect the default output. | The issue detail output is useful for real inspection without needing to query the database or journal manually. | | |
| HV-01c | VAL-010 | Run `fid issue get <id> --section <section>` for one supported section that matters to the tester. | The section output is scoped, readable, and useful enough to inspect that part of the issue. | | |
| HV-01d | VAL-010 | Run `fid issue list` with at least one filter or ordering option against the created working set. | The list output is useful for scanning and finding the expected issue without excessive noise. | | |
| HV-01e | VAL-010 | Trigger one invalid hierarchy or scheduling error, such as a subtask under an epic or a sequenced task under a prioritised epic. | The error is clear, actionable, and names enough context to correct the command. | | |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Dogfood Loop | 5 | 0 | 0 | 0 | 5 |
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
-->
