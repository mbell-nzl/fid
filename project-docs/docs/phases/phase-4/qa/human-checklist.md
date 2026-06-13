# Phase 4 Human Verification Checklist

> Phase: 4 - v0.1b Issue authoring and discovery
> Date: 2026-06-09
> Tester: MBEL
>
> Include only `(human)` items here. Deterministic checks are recorded in
> `final-review.md`. Every row maps to VAL-010.

## Dogfood Loop

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HV-01a | VAL-010 | Using the built binary, create one epic, one task under that epic, and one subtask under that task during a realistic local workflow. | The create flow is understandable, the required flags are discoverable enough, and the success messages make the created IDs usable for the next command. | PASS | n/a |
| HV-01b | VAL-010 | Run `fid issue get <id>` on one created issue and inspect the default output. | The issue detail output is useful for real inspection without needing to query the database or journal manually. | PASS | n/a |
| HV-01c | VAL-010 | Run `fid issue get <id> --section <section>` for one supported section that matters to the tester. | The section output is scoped, readable, and useful enough to inspect that part of the issue. | PASS | n/a |
| HV-01d | VAL-010 | Run `fid issue list` with at least one filter or ordering option against the created working set. | The list output is useful for scanning and finding the expected issue without excessive noise. | PASS | n/a |
| HV-01e | VAL-010 | Trigger one invalid hierarchy or scheduling error, such as a subtask under an epic or a sequenced task under a prioritised epic. | The error is clear, actionable, and names enough context to correct the command. | PASS | n/a |

### TL Follow-Up Notes

- During HV-01e, a deliberate invalid create was attempted: add a task with
  `--sequence 1` under prioritised epic `FID-e36ee9`. The CLI returned the
  expected constraint error, but the rejected command was still appended to the
  canonical journal as `EVT-20260610-000007` in
  `.fid/journal/2026-06-10.jsonl`.
- After that, unrelated reads such as `fid issue get FID-a816f0` failed with:
  `error: project event "EVT-20260610-000007": insert into issues: constraint failed: issues: task under prioritised epic must not have sequence (1811)`.
- Diagnosis: `issues.Create` appends the event before projecting it. When
  projection rejects the event, the append has already happened, poisoning the
  journal tail. Later `issue get` calls `EnsureProjectionFresh`, tries to
  replay the invalid event, and fails before reading the requested accepted
  issue.
- `FID-a816f0` appears to be valid; the persistent failure is caused by the
  later invalid journal event. This is a product/code defect found by the human
  checklist: rejected create commands must not advance the canonical journal.
  Add a regression test that rejected creates leave the journal line count
  unchanged, then validate scheduling/sequence before append or add an
  equivalent pre-append projection validation path.

### TL Wrap-Up Disposition

- **All checklist items passed (5/5), but testing surfaced an issue that needs
  human attention.** The HV-01e invalid-create probe revealed that a *rejected*
  `fid issue create` still appends its event to the canonical journal, poisoning
  the tail and breaking later reads (`fid issue get`) until the bad line is
  removed. This is a journal-integrity / data-availability defect.
- Recorded as **carry-forward CF-007** (target Phase 5) and **DEC-014**
  (`project-docs/docs/phases/phase-4/decisions.md`). Likely a fixup: validate
  scheduling/sequence before the journal append (or add a pre-append
  projection-dry-run), and add a regression test that a rejected create leaves
  the journal line count unchanged.
- Phase 4 still closes as **PASS** on its validation contract; CF-007 is the
  one open thread carried into Phase 5. Please confirm CF-007 priority before
  Phase 5 planning closes.

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Dogfood Loop | 5 | 5 | 0 | 0 | 0 |
| **Total** | 5 | 5 | 0 | 0 | 0 |

| | |
|---|---|
| **Overall result** | PASS but note discovered bug|
| **Date completed** | 2026-06-13 |
| **Notes** | See notes regarding bug. |

<!--
Outcome values: pass, fail, partial, blank.
Overall result: pass if all pass; partial if any partial and no fail; fail if
any fail.
-->
