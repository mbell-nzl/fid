# Phase 12 Human Verification Checklist

> Phase: 12 - Requirements + Links
> Date: 2026-06-27
> Tester: MBEL
>
> Include only `(human)` items here: perception, product-owner judgment,
> dogfooding, tone, and flow coherence. Deterministic `(qa-lead)` checks belong
> in `final-review.md` with reproducible evidence. Human checks that validate a
> phase assertion cite the relevant `VAL-NNN`.

## Requirements And Link Dogfood

One observable per row. Outcome values: `pass`, `fail`, `partial`, or `blank`.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HV-12a | VAL-012 | In a real `.fid` workspace, dogfood `fid requirement sync` against a PRD/TRD-style source document and judge whether the CLI output makes the recorded pointer counts understandable without reading implementation `(human)` | The tester can identify created, updated, and unchanged counts from the output and can explain what happened to requirement pointers. | | |
| HV-12b | VAL-012 | In the same workspace, dogfood `fid link add` for a valid link and judge whether the CLI output makes the auto-created inverse understandable `(human)` | The tester can identify both the declared direction and the inverse direction from the output. | | |
| HV-12c | VAL-012 | Dogfood `fid link list` for both linked endpoints and judge whether outgoing and incoming directions are understandable in normal CLI use `(human)` | The tester can identify direction, relation, linked entity ID/type, and title for links from both endpoint perspectives. | | |
| HV-12d | VAL-012 | Dogfood `fid link remove` for the declared link and judge whether the CLI output makes removal of both directions understandable `(human)` | The tester can identify that both the declared link and inverse link were removed. | | |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Requirements And Link Dogfood | 4 | | | | |
| **Total** | 4 | | | | |

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
