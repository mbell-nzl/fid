# Phase 16 Human Checklist Run Sheet

> Phase: 16 - Activity Timeline + Read-Only Escape Hatches
> Purpose: executable command sheet for
> `project-docs/docs/phases/phase-16/qa/human-checklist.md`
> Date prepared: 2026-07-05

Use this from the fid repository root. It creates a disposable real `.fid/`
workspace under `/tmp`, builds the current CLI to `/tmp/fid`, creates a task,
adds an acceptance criterion, then runs the activity, SQL, and journal-inspect
commands needed for the human dogfood checks.

The human checklist still needs human judgment. Treat the command output as
dogfood evidence, then record the outcome and notes in `human-checklist.md`.

## Prerequisites

- Run the build command from the repo root.
- Run workspace commands only inside the disposable `/tmp` workspace created
  below, not inside this repository's own `.fid/` workspace.
- These commands require the built `fid` binary and POSIX shell tools already
  used by the project docs (`sed`, `printf`, `ls`, `mkdir`, `mktemp`).
- This sheet does not complete deterministic QA gates. Phase 16 final review
  already covers activity filters and format parity, CF-011 rebuild
  determinism, SQL read-only enforcement, journal inspect hash/verdict parity,
  gate history reads, and phase-level gates. This sheet only supports
  `HUMAN-16a` through `HUMAN-16d` for `VAL-013`.

## 1. Build Current CLI

```bash
go build -o /tmp/fid ./cmd/fid
export FID=/tmp/fid
ls -l "$FID"
```

Expected: command exits zero and `/tmp/fid` exists.

## 2. Create A Fresh Test Workspace

```bash
export FID_HV16_WORKSPACE="$(mktemp -d /tmp/fid-phase16-human-XXXXXX)"
cd "$FID_HV16_WORKSPACE"
printf 'workspace=%s\n' "$FID_HV16_WORKSPACE"

"$FID" init --prefix HV16
```

Expected:

- `fid init` reports a new `.fid/` workspace.
- The workspace has a `.fid/` directory.
- The rest of this sheet is run from the disposable workspace.

## 3. Create An Issue And A Second Activity Row

Create one task, then add one acceptance criterion. The task creation and AC add
give the activity timeline two issue-scoped writes to inspect.

```bash
TASK_OUTPUT=$("$FID" issue create \
  --kind task \
  --title "HV16 human dogfood task" \
  --priority 1 \
  --est-tokens 700 \
  --by developer 2>&1)
printf '%s\n' "$TASK_OUTPUT"
export HV16_TASK="$(printf '%s\n' "$TASK_OUTPUT" | sed -n 's/^created issue \([^ ]*\) .*/\1/p')"
printf 'task=%s\n' "$HV16_TASK"

AC_OUTPUT=$("$FID" issue ac add "$HV16_TASK" \
  --text "Human run sheet verifies activity output" \
  --ordinal 1 \
  --by developer 2>&1)
printf '%s\n' "$AC_OUTPUT"
export HV16_AC="$(printf '%s\n' "$AC_OUTPUT" | sed -n 's/^added AC \([^ ]*\) .*/\1/p')"
printf 'ac=%s\n' "$HV16_AC"
```

Expected:

- The task is created.
- The acceptance criterion is added to that task.
- `HV16_TASK` and `HV16_AC` print non-empty IDs.

## 4. HUMAN-16a - Inspect The Activity Timeline

Run the timeline read for the created task. Inspect the markdown output as a
human operator.

```bash
"$FID" activity list --on "issue/$HV16_TASK" --format markdown
```

Expected for `HUMAN-16a`:

- The command exits zero.
- The output has an `Activity` heading and a table.
- The table includes `issue_created` for `issue/$HV16_TASK`.
- The table includes `ac_added` for `issue/$HV16_TASK`.
- The actor column is understandable, and the `Details` column names useful
  facts such as the title or AC ID.

Record `HUMAN-16a` as `pass` if the timeline is understandable from CLI text
alone and clearly reflects the recent task creation and AC add. Record
`partial` or `fail` if the output is hard to connect to the writes without
reading implementation files or the database.

## 5. HUMAN-16b - Run A Read-Only SQL Query

Run a read-only query over the projected issue table. Inspect the output as a
human operator.

```bash
"$FID" sql "SELECT id, title, status FROM issues ORDER BY id LIMIT 5" \
  --format markdown
```

Expected for `HUMAN-16b`:

- The command exits zero.
- The output has a `SQL Query Results` heading and a table.
- The table includes columns `id`, `title`, and `status`.
- The row for `$HV16_TASK` is visible with title `HV16 human dogfood task`.
- The output is readable without knowing the implementation package layout.

Record `HUMAN-16b` as `pass` if the read-only query returns understandable rows
from the real workspace. Record `partial` or `fail` if the output is technically
successful but not useful to a human reader.

## 6. HUMAN-16c - Confirm SQL Mutation Rejection

Run a deliberately unsafe mutation attempt. This command is expected to fail;
the human check is about the clarity of the rejection and whether the workspace
still appears unchanged afterward.

```bash
"$FID" sql "DELETE FROM issues" --format markdown
```

Expected for `HUMAN-16c`:

- The command exits non-zero.
- The output says `fid sql accepts only SELECT, EXPLAIN, PRAGMA query, and WITH
  statements`.
- The output shows the rejected statement, `DELETE FROM issues`.
- The message reads like a clear rejection, not like a successful empty result.

Now confirm the issue still exists:

```bash
"$FID" sql "SELECT count(*) AS issue_count FROM issues" --format markdown
"$FID" activity list --on "issue/$HV16_TASK" --format markdown
```

Expected:

- The count query exits zero and shows at least one issue.
- The activity timeline for `$HV16_TASK` still renders.

Record `HUMAN-16c` as `pass` if the mutation is clearly rejected and the
follow-up read commands make it appear that the workspace was not mutated.
Record `partial` or `fail` if the rejection is ambiguous or the follow-up output
suggests the issue disappeared.

## 7. Capture A Journal Event ID

Run the activity list in JSON so the journal event ID is visible from CLI
output, then capture one event ID for inspection.

```bash
ACTIVITY_JSON=$("$FID" activity list --on "issue/$HV16_TASK" --format json)
printf '%s\n' "$ACTIVITY_JSON"
export HV16_EVENT="$(printf '%s\n' "$ACTIVITY_JSON" | sed -n 's/.*"journal_event_id":"\([^"]*\)".*/\1/p')"
printf 'event=%s\n' "$HV16_EVENT"
```

Expected:

- The JSON output includes `journal_event_id`.
- `HV16_EVENT` prints a non-empty `EVT-...` ID.

If `HV16_EVENT` is empty, pick an event ID directly from the journal:

```bash
export HV16_EVENT="$(sed -n '2s/.*"event_id":"\([^"]*\)".*/\1/p' .fid/journal/*.jsonl)"
printf 'event=%s\n' "$HV16_EVENT"
```

## 8. HUMAN-16d - Inspect The Journal Event

Inspect the captured journal event in markdown.

```bash
"$FID" journal inspect "$HV16_EVENT" --format markdown
```

Expected for `HUMAN-16d`:

- The command exits zero.
- The output heading names `Journal Event EVT-...`.
- The output shows the event schema, event type, journal file and line, and a
  `sha256:` hash.
- The output has `Actor`, `Command`, `Effects`, and `Revert` sections.
- The effects are connected to the task or AC write from this run.
- The revert section gives an understandable verdict or correction note.

Record `HUMAN-16d` as `pass` if the event output explains what happened well
enough from CLI text alone. Record `partial` or `fail` if the location, effects,
hash, or revert verdict are missing or too opaque to use.

## 9. Fill The Human Checklist

Return to the repository root and record the outcomes in the checklist.

```bash
cd /home/mark/projects/fid
sed -n '1,220p' project-docs/docs/phases/phase-16/qa/human-checklist.md
```

Fill:

- `HUMAN-16a` from section 4.
- `HUMAN-16b` from section 5.
- `HUMAN-16c` from section 6.
- `HUMAN-16d` from section 8.
- Summary totals and overall result.

Use `pass`, `partial`, or `fail` for each row. Overall result is `pass` only if
all four rows pass; `partial` if at least one row is partial and none fail;
`fail` if any row fails.

## 10. Cleanup

When the checklist evidence is recorded, remove the disposable workspace if it
is no longer needed.

```bash
printf 'workspace=%s\n' "$FID_HV16_WORKSPACE"
rm -rf "$FID_HV16_WORKSPACE"
```

Expected: the temporary `/tmp/fid-phase16-human-*` workspace is removed. Do not
run the cleanup command until any useful terminal output has been recorded in
the checklist notes.
