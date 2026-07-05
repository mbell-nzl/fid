# Phase 16 Human Verification Checklist

> Phase: 16 - Activity Timeline + Read-Only Escape Hatches
> Date: 2026-07-05
> Tester: _______________
>
> Human checks only. Deterministic QA evidence is recorded in
> `project-docs/docs/phases/phase-16/qa/final-review.md`.

## Dogfood Read Surfaces

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-16a | VAL-013 | In a real `.fid` workspace, create or update an issue, then run `fid activity list --on issue/<id>`. | The timeline output is understandable from CLI text alone and reflects the recent write. | | |
| HUMAN-16b | VAL-013 | Run a read-only query such as `fid sql "SELECT id, title FROM issues LIMIT 5"`. | The command returns readable rows without requiring implementation knowledge. | | |
| HUMAN-16c | VAL-013 | Run a mutation attempt such as `fid sql "DELETE FROM issues"`. | The command is rejected clearly and does not appear to mutate the workspace. | | |
| HUMAN-16d | VAL-013 | Pick an event ID from CLI-visible output or the journal and run `fid journal inspect <event-id>`. | The output explains the event, location, effects, hash, and revert verdict well enough to understand what happened. | | |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Dogfood Read Surfaces | 4 | 0 | 0 | 0 | 4 |
| **Total** | 4 | 0 | 0 | 0 | 4 |

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
