# Phase 6 Human Verification Checklist

> Phase: 6 - Ready Queue + v0.1b Closure
> Date: 2026-06-14
> Tester: MBEL
>
> Include only human dogfood judgment here. Deterministic command gates and
> structural assertions are recorded in `final-review.md`.

## `fid ready` Dogfood

Use a real `.fid/` workspace and the public v0.1b `fid ready` surface only:
`fid ready [--epic <id>] [--format json|toon|markdown|yaml]`. Do not use hidden
scenario material, evaluator prompts, or non-public test fixtures.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HV-012a | VAL-012 | In a real workspace with a sequential epic, run `fid ready --epic <sequential-epic-id> --format json` for the epic's current scheduling state. Judge the result against the task a human would choose next for that sequence. | The returned task is the lowest-sequence unfinished actionable task. If the lowest unfinished task is already in flight, `task` is null and the stderr `note:` explains that current task rather than suggesting a later task. | pass | Human reported the fresh-workspace dogfood run selected the expected sequence-2 task and did not skip to the later higher-priority sequential task. |
| HV-012b | VAL-012 | In a real workspace with a prioritised epic containing multiple actionable candidates, run `fid ready --epic <prioritised-epic-id> --format json`. Judge whether the selected task matches human priority expectation. | The returned task is the highest-priority actionable task. Started tasks and lower-priority actionable tasks are not chosen ahead of it. | pass | Human reported the fresh-workspace dogfood run selected the expected priority-1 actionable task, excluding the in-flight priority-0 task and lower-priority task. |
| HV-012c | VAL-012 | In a real workspace or epic where no task should be ready because all tasks are completed, cancelled, or already in flight, run `fid ready --format json` or `fid ready --epic <epic-id> --format json`. Judge whether the empty result is understandable. | The command exits successfully, stdout is the frozen empty `task` shape, and the stderr `note:` gives a clear, expected reason that no ready task exists. | pass | Human reported the scoped empty-epic dogfood run returned `{"task":null}` and a clear `note:` explaining no ready tasks existed in the epic. |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| `fid ready` Dogfood | 3 | 3 | 0 | 0 | 0 |
| **Total** | 3 | 3 | 0 | 0 | 0 |

| | |
|---|---|
| **Overall result** | pass |
| **Date completed** | 2026-06-14 |
| **Notes** | Human reported all three public `fid ready` dogfood checks passed in a fresh real `.fid/` workspace using the phase-end command sequence. |

<!--
Outcome values: pass, fail, partial, blank.
Overall result: pass if all pass; partial if any partial and no fail; fail if
any fail.
-->
