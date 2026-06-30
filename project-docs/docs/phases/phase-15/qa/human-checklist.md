# Phase 15 Human Verification Checklist

> Phase: 15 - `fid dev *` Composites + Cleanup
> Date: 2026-06-29
> Tester: MBEL
>
> This checklist contains only human-observable dogfood checks for VAL-012.
> Deterministic QA evidence is recorded in `final-review.md`.

## CLI Dogfood Observability

One observable per row. Run against a real `.fid` workspace after the
deterministic phase gate has passed.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HV-15a | VAL-012 | Run `fid dev start <task-id> --model <model> --by developer` and read the CLI output without inspecting implementation or the database. `(human)` | The output makes it clear that the task was claimed by the developer, moved to `in_progress`, and an agent run was opened. | PASS | n/a |
| HV-15b | VAL-012 | Run `fid dev gate <task-id> --cmd "<gate command>" --exit-code <n> --output-excerpt "<excerpt>" --by developer` and read the CLI output without inspecting implementation or the database. `(human)` | The output makes it clear that fid recorded the reported gate result and does not appear to have executed the command. | PASS | n/a |
| HV-15c | VAL-012 | After developer ACs are verified, run `fid dev finish <task-id> --by developer --completion-md <path>` with any applicable optional completion flags, then read the CLI output without inspecting implementation or the database. `(human)` | The output makes it clear that the task was completed, the open run was closed, and the developer claim was released. | PASS | n/a |
| HV-15d | VAL-012 | On a separate in-progress task, run `fid dev block <task-id> --by developer --reason-md <path>` and read the CLI output without inspecting implementation or the database. `(human)` | The output makes it clear that a failed-gate completion/reason was recorded and the developer claim was released. | PASS | n/a |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| CLI Dogfood Observability | 4 | 4 | 0 | 0 | 0 |
| **Total** | 4 | 4 | 0 | 0 | 0 |

| | |
|---|---|
| **Overall result** | PASS |
| **Date completed** | 2026-06-29 |
| **Notes** | See notes below -- possible carry-forwards. |

<!--
Outcome values: pass, fail, partial, blank.
Overall result: pass if all pass; partial if any partial and no fail; fail if
any fail.
Granularity: one observable per row. Dense comma-listed rows are imprecise and
let items get missed.
-->

## TL Wrap Review Notes

- `fid dev gate` records append-only `gate_runs`, but there is no direct
  issue-level read surface for those recorded gates while a task is still in
  progress. `fid issue get` currently supports `acceptance`, `concerns`,
  `completions`, `links`, and `claims`; it does not support a `gates` section.
- After `fid dev finish` or `fid dev block`, gate evidence may be visible only
  indirectly through the completion snapshot (`task_completions.gate_results`),
  not as the raw gate history. This makes it unclear how a developer, QA
  reviewer, or TL should inspect gates recorded by `fid dev gate` before a
  completion exists.
- Suggested TL wrap review question: should Phase 16 or a carry-forward item add
  a public read surface such as `fid issue get <task-id> --section gates` or
  `fid dev gate list <task-id>`? The former seems to fit the existing
  issue-section model and would keep completion snapshots distinct from raw
  append-only gate history.
