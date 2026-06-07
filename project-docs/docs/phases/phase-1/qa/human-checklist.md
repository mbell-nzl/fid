# Phase 1 Human Verification Checklist

> Phase: 1 - End-to-End Persistence Loop (Tracer Bullet)
> Date: 2026-06-07
> Tester: MBEL
>
> This checklist contains only `(human)` items from the Phase 1 QA plan.
> Deterministic VAL-001 through VAL-017 evidence is recorded in
> `project-docs/docs/phases/phase-1/qa/final-review.md`.

## Dogfooding Flow

Run in a disposable directory with the Phase 1 binary built by QA:
`/tmp/fid`.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HV-01a | no VAL | Run `/tmp/fid init` in a new disposable directory and read the command output as a first-time user. | PASS if the output makes the created workspace, prefix, config, pricing, journal event, and seeded roles understandable without reading source code; FAIL if it is confusing or omits information needed to proceed. | PASS | n/a |
| HV-01b | no VAL | In the same directory, run `/tmp/fid issue create --kind task --title "Human smoke" --by developer`. | PASS if the command shape and success output make it clear that a task issue was created and which ID to use next; FAIL if the command or output is unclear for first dogfood use. | PASS | n/a |
| HV-01c | no VAL | Run `/tmp/fid issue get ID_FROM_HV_01B` and `/tmp/fid issue get ID_FROM_HV_01B --format json`, replacing `ID_FROM_HV_01B` with the issue ID returned by HV-01b. | PASS if the default read shape and explicit JSON read feel coherent for the tracer-bullet phase; FAIL if the output is too noisy, missing obvious fields, or hard to inspect. | PARTIAL | Output is raw json with no newlines that needs to be piped through `jq` to print nicely. Both standard command and `--format json` show exactly the same right now. Fine at this stage as long as will be addressed later. |
| HV-01d | no VAL | Run `/tmp/fid rebuild --by developer` after the issue exists. | PASS if the rebuild output communicates what happened and gives enough confidence that the projection was rebuilt from the journal; FAIL if it is ambiguous or raises unnecessary concern. | PASS | n/a |

## Error Message Judgment

Run these in a disposable initialized workspace unless the test says otherwise.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HV-02a | no VAL | Run `/tmp/fid init` a second time in a directory that already contains `.fid/`. | PASS if the rejection explains that `.fid/` already exists and what action would be needed to reinitialize; FAIL if the message is vague or misleading. | PASS | n/a |
| HV-02b | no VAL | Run `/tmp/fid issue create --kind task --title "Missing role"` without `--by`. | PASS if the missing-role error is clear and actionable; FAIL if the user would not know how to fix the command. | PASS | n/a |
| HV-02c | no VAL | Run `/tmp/fid issue create --kind task --title "Unknown role" --by stranger`. | PASS if the unknown-role error clearly names the rejected role and provides useful known-role context; FAIL if the role problem is unclear. | PASS | n/a |
| HV-02d | no VAL | Run `/tmp/fid rebuild` without `--by`, then `/tmp/fid rebuild --by stranger`. | PASS if both rebuild role errors are clear and consistent with issue-create errors; FAIL if rebuild feels inconsistent or harder to correct. | PASS | n/a |
| HV-02e | no VAL | Review the final-review evidence for corrupt-tail and precondition-conflict errors, or reproduce with the existing Phase 1 tests if desired. | PASS if the reported file/line and conflict details would be actionable for manual repair or escalation; FAIL if a human would not know what artifact needs attention. | PASS | Focused tests passed. Corrupt-tail evidence names `.fid/journal/2026-06-05.jsonl:2`; precondition-conflict evidence exposes event ID, table, and line, with the implementation error including event type, pk, expected, actual, file, and line. |
| HV-02f | no VAL | Review the final-review evidence for held-lock timeout behavior, or reproduce with the existing Phase 1 lock test if desired. | PASS if the holder and lock path information would help identify the blocking process/workspace; FAIL if it does not support a practical next step. | PASS | Focused lock timeout test passed. Error text includes `.fid/write.lock` path plus holder role, pid, hostname, command, and started_at metadata, which is enough to identify the blocking workspace/process. |

## Product Judgment

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HV-03a | no VAL | Inspect the minimal issue JSON returned by `fid issue get` for a newly created task. | PASS if JSON-only output is acceptable as the Phase 1 default while TOON, markdown, and YAML remain deferred; FAIL if another read shape is necessary before closing the phase. | PASS | n/a |
| HV-03b | no VAL | Judge whether the tracer-bullet loop `init -> issue create -> issue get -> rebuild -> issue get` is sufficient for Phase 1 dogfooding. | PASS if the loop demonstrates the intended local-first journal/projection/rebuild promise for this phase; FAIL if a missing command or output prevents meaningful dogfooding. | PASS | n/a |
| HV-03c | no VAL | Review the known deterministic final-review gap: coverage threshold failure. | PASS if the human accepts blocking Phase 1 until coverage is remediated or explicitly deferred by Tech Lead/human decision; FAIL if the final review should be revised under a different quality interpretation. | PASS | Pass from human standpoint given coverage % less meaningful with small amount of code. Request thorough TL review to judge if this is truly viable or whether the developer agents were lazy and did not follow TDD guidance. |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Dogfooding Flow | 4 | 3 | 0 | 1 | 0 |
| Error Message Judgment | 6 | 6 | 0 | 0 | 0 |
| Product Judgment | 3 | 3 | 0 | 0 | 0 |
| **Total** | 13 | 12 | 0 | 1 | 0 |

| | |
|---|---|
| **Overall result** | PARTIAL |
| **Date completed** | 2026-06-07 |
| **Notes** | See inline notes. |

<!--
Outcome values: pass, fail, partial, blank.
Overall result: pass if all pass; partial if any partial and no fail; fail if
any fail.
Granularity: one observable per row. Dense comma-listed rows are imprecise and
let items get missed.
-->
