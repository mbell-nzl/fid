# Phase 21 Human Verification Checklist

> Phase: 21 - One-Issue Operational Walking Skeleton
> Date: 2026-09-03
> Tester: _______________
>
> Include only `(human)` items here: perception, product-owner judgment,
> dogfooding, tone, and flow coherence. Deterministic `(qa-lead)` checks belong
> in `final-review.md` with reproducible evidence. Human checks that validate a
> phase assertion cite the relevant `VAL-NNN`.

`VAL-036` is the only Human assertion in phase 21. It is one perception
assertion fanned out into four discrete observables. QA has deliberately not
described what the recorded session contains, so that each answer below is a
genuine first reading rather than a confirmation of QA's summary.

## Setup

Perform this once, yourself, in a **real interactive terminal** at the
repository root `/home/mark/projects/fid`:

```text
cat cmd/fid/testdata/val-036/human-session.txt
```

Notes on the setup:

- The file is the byte-exact recorded Human-mode stdout of one `fid gate run`.
- `cat` is used deliberately, with no pager, no `-v`, and no `cat -A`. Any
  unescaped control byte in the file would reach your terminal directly, which
  is what HUMAN-01d observes. Using `less`, `bat`, an editor, or `cat -A`
  would suppress that effect and invalidate HUMAN-01d.
- Answer HUMAN-01a, HUMAN-01b, and HUMAN-01c from that displayed session alone.
  Do not open the source, the tests, the live-frame fixture, the phase
  contracts, or `final-review.md` first.
- Answer each row `yes` or `no` and write the requested note. A `no` on any row
  makes `VAL-036` fail; the finding is then owned by task 14.

## Gate Session Comprehension (VAL-036)

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-01a | VAL-036 | Reading the recorded session, tell which lines came from the child process's stdout and which came from its stderr, without consulting any other artifact | yes - stream origin is unambiguous for every process-output line | PASS | n/a |
| HUMAN-01b | VAL-036 | Reading the same session, identify each typed control notice as a distinct, named event rather than as process output | yes - every control notice is recognisable as fid's own typed notice, and you can name each one you found | PASS | n/a |
| HUMAN-01c | VAL-036 | Reading the same session, state the terminal result: whether a Receipt exists, the outcome, the exit or signal or timeout, Artifact IDs and rejections, and the Obligation status | yes - you can restate the complete terminal result in your own words from the session alone | PASS | n/a |
| HUMAN-01d | VAL-036 | While the file is displayed with `cat` in a real terminal, watch for any raw control byte taking effect: colour change, cursor movement, bell, or line overwrite | yes - no display change occurs; the terminal renders plain text only | PASS | n/a |

Note for HUMAN-01c: the recorded session is a **non-success** execution. Judge
whether the terminal frame reports that honestly and legibly, not whether the
gate passed.

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Gate Session Comprehension (VAL-036) | 4 | 4 | 0 | 0 | 0 |
| **Total** | 4 | 4 | 0 | 0 | 0 |

| | |
|---|---|
| **Overall result** | PASS |
| **Date completed** | 2026-09-04 |
| **Notes** | n/a |

<!--
Outcome values: pass, fail, partial, blank.
Overall result: pass if all pass; partial if any partial and no fail; fail if
any fail.
Granularity: one observable per row. Dense comma-listed rows are imprecise and
let items get missed.
VAL-036 passes only when all four rows are answered yes. Any no is a Must-fix
finding owned by task 14.
-->
