# Phase 22 Human Verification Checklist

> Phase: 22 — Recovery and portability
> Prepared: 2026-09-08
> Tester: _______________
> Status: **PENDING task-18 fixture manifest and independent readiness verification.**
> Human responses and outcomes have not been recorded.

## Before starting

Use a fresh session without prior chat or an explanation of the intended
answers. Inspect the supplied views and write your own response for each row
before recording an outcome. Record anything unclear in the notes, including
the wording or view that led to your interpretation.

Task 18 will provide disposable projects, a retained binary, exact identities
and neutral inspection commands in its generated manifest. The preparation
record below will be completed from that manifest after QA verifies that the
commands work after the generator exits. There are no runnable fixture
instructions in this draft yet. Do not use an existing personal project.

| Fixture preparation record | Verified location or value |
|---|---|
| Fixture guide | PENDING task-18 `qa/human-fixtures.md` |
| Generated manifest | PENDING generation with `FID_P22_HUMAN_FIXTURE_ROOT` |
| Retained binary | PENDING manifest |
| Interrupted Run / adopted changes project and inspection commands | PENDING manifest |
| Pre-commit activation project and inspection commands | PENDING manifest |
| Post-commit activation project and inspection commands | PENDING manifest |
| Imported project and inspection commands | PENDING manifest |
| Independent command and retained-fixture verification | PENDING task-18 completion |

## Recovery view

Open the supplied resume view for the interrupted Run using its manifest
instructions.

| Test ID | Covers | Prompt | Your response / evidence | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-01a (human) | VAL-076 | What single recovery action would you take next? State what in the view led you to that choice. | | | |
| HUMAN-01b (human) | VAL-076 | What does the view tell you about the proof status of the adopted changes? Explain your interpretation. | | | |

## Activation recovery

Inspect each supplied activation checkpoint using its own manifest
instructions.

| Test ID | Covers | Prompt | Your response / evidence | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-02a (human) | VAL-077 | In the pre-commit checkpoint, which recovery direction is offered, and what in the view explains it? | | | |
| HUMAN-02b (human) | VAL-077 | In the post-commit checkpoint, is rollback possible? Explain what you understand from the view. | | | |

## Imported work

Open the supplied imported project using its manifest instructions.

| Test ID | Covers | Prompt | Your response / evidence | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-03a (human) | VAL-078 | Which restored Issue, Run and history are you looking at? Record the identities and what helped you recognize the work. | | | |
| HUMAN-03b (human) | VAL-078 | What intervention, if any, remains before this imported work can be mutated? State what in the views led you to that conclusion. | | | |

## Human results

All six outcomes remain blank. Agents may prepare and verify the fixtures;
only the Human records the judgments above. Deterministic QA evidence belongs
in the separate final review.

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Recovery view | 2 | 0 | 0 | 0 | 2 |
| Activation recovery | 2 | 0 | 0 | 0 | 2 |
| Imported work | 2 | 0 | 0 | 0 | 2 |
| **Total** | **6** | **0** | **0** | **0** | **6** |

| Result record | Human entry |
|---|---|
| Overall result | |
| Date completed | |
| Notes | |

Outcome values are `pass`, `fail` or `partial`; leave an unperformed item
blank. Record the overall result only after completing the six responses.
