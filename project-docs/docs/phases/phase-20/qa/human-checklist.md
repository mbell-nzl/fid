# Phase 20 Human Verification Checklist

> Phase: 20 - Standalone Issue Journey
> Date: 2026-08-16
> Tester: _______________
>
> This checklist contains only Human judgment for VAL-038. Deterministic
> command, schema, channel, and file-state checks are recorded in
> `final-review.md`. Do not infer an outcome from the prepared transcript.

## Prepared Human-Mode Corpus

The QA Lead prepared a fresh workspace at `/tmp/fid-p20-human.vAjV9v` with the
QA-basis `/tmp/fid` binary. The routine corpus is:

```text
$ /tmp/fid init --prefix HQA
initialized fid project in .fid/

$ /tmp/fid issue create --kind task --title "Human QA routine task" --by developer
created issue HQA-54265f (task) "Human QA routine task"

$ /tmp/fid issue get HQA-54265f
HQA-54265f | todo | Human QA routine task

$ /tmp/fid issue update HQA-54265f --title "Human QA routine task updated" --expected-revision 1 --by developer
updated HQA-54265f:
title "Human QA routine task" -> "Human QA routine task updated"
```

Those routine commands exited 0, wrote the displayed text to stdout, and wrote
nothing to stderr. The failure corpus is:

```text
$ /tmp/fid issue update HQA-54265f --title "stale write" --expected-revision 1 --by developer
error: issue "HQA-54265f" revision is 2, not expected revision 1
```

The failure exited 1, wrote the displayed text to stderr, and wrote nothing to
stdout. A Human may use this transcript or rerun the commands in a fresh
workspace. The three rows below are separate judgments; leave a row blank if
it was not judged.

## Human Output And Guidance Judgment

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-038a | VAL-038 | Judge the displayed routine init/create/get/update output as one interaction corpus. | The routine output is concise enough for normal terminal use. | | |
| HUMAN-038b | VAL-038 | Judge the displayed stale-revision failure message. | The message gives sufficiently actionable guidance for the legitimate next action. | | |
| HUMAN-038c | VAL-038 | Inspect the displayed routine and failure corpus for direct or indirect advice about fid storage/config internals. | No message advises editing `.fid/`, SQLite, journal files, generated config, or another fid internal directly. | | |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Human Output And Guidance Judgment | Pending | Pending | Pending | Pending | Pending |
| **Total** | **Pending** | **Pending** | **Pending** | **Pending** | **Pending** |

| | |
|---|---|
| **Overall result** | Pending |
| **Date completed** | |
| **Notes** | |

<!--
Outcome values: pass, fail, partial, blank.
Overall result: pass if all pass; partial if any partial and no fail; fail if
any fail.
Granularity: one observable per row. Dense comma-listed rows are imprecise and
let items get missed.
-->
