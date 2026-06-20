# Phase 9 Human Checklist Run Sheet

> Phase: 9 - Concerns + v0.1c Closure
> Purpose: executable command sheet for
> `project-docs/docs/phases/phase-9/qa/human-checklist.md`
> Date prepared: 2026-06-21

Use this from the fid repository root. It creates a disposable `.fid/`
workspace under `/tmp`, builds the current CLI to `/tmp/fid`, creates one task,
records two completions, raises a concern against the first completion, addresses
it from the second completion, then walks the human checklist rows for
`VAL-011`.

The human checklist still needs human judgment. Treat the command output as
evidence, then record the outcome and any friction in
`human-checklist.md`.

## Prerequisites

- Run the build command from the repo root.
- Run workspace commands only inside the disposable `/tmp` workspace created
  below, not inside this repository's own `.fid/` workspace.
- These commands require the built `fid` binary and POSIX shell tools already
  used by the project docs (`sed`, `wc`, `ls`, `sort`, `tail`, `mv`, `printf`).
- Task 4 fixed the earlier Phase 9 closeout blockers, so read-back is expected
  to include `addressed_in_completion_id` and concern list ordering is expected
  to be deterministic newest-first append order.
- This sheet does not complete deterministic QA gates. It is only the human
  dogfood aid for `HUMAN-011a` through `HUMAN-011e`.

## 1. Build Current CLI

```bash
go build -o /tmp/fid ./cmd/fid
export FID=/tmp/fid
ls -l "$FID"
```

Expected: command exits zero and `/tmp/fid` exists.

## 2. Create A Fresh Test Workspace

```bash
export FID_HV9_WORKSPACE="$(mktemp -d /tmp/fid-phase9-human-XXXXXX)"
cd "$FID_HV9_WORKSPACE"
printf 'workspace=%s\n' "$FID_HV9_WORKSPACE"

"$FID" init --prefix HV9

TASK_CREATE=$("$FID" issue create \
  --kind task \
  --title "HV9 concern lifecycle dogfood" \
  --by developer 2>&1)
printf '%s\n' "$TASK_CREATE"

TASK_ID=$(printf '%s\n' "$TASK_CREATE" | sed -n 's/^created issue \([^ ]*\).*/\1/p')
printf 'task_id=%s\n' "$TASK_ID"

printf '# HV9 concern body\n\nRaised during human dogfood.\n' > concern-body.md
```

Expected:

- `fid init` reports a new `.fid/` workspace.
- `task_id` starts with `HV9-`.
- `concern-body.md` exists inside the disposable workspace.

## 3. Baseline Issue Reads Before Raising A Concern

```bash
"$FID" issue get "$TASK_ID" --format json
"$FID" issue list --format json
```

Expected for the baseline side of `HUMAN-011b`:

- The task is visible.
- `concerns_raised` is `false`.
- `open_concerns_count` is `0`.
- The fields are easy enough to find in both `issue get` and `issue list`.

Do not record `HUMAN-011b` yet; it also needs the after-raise view.

## 4. HUMAN-011a - Raise A Concern Against A Real Completion

Create the completion that the concern will be raised against, then raise the
concern.

```bash
COMPLETE1=$("$FID" issue complete "$TASK_ID" --by developer 2>&1)
printf '%s\n' "$COMPLETE1"

RAISE_COMPLETION_ID=$(printf '%s\n' "$COMPLETE1" | sed -n 's/^[[:space:]]*completion: \(TC-[^ ]*\).*/\1/p')
printf 'raise_completion_id=%s\n' "$RAISE_COMPLETION_ID"

RAISE_OUTPUT=$("$FID" concern raise \
  --on "$TASK_ID" \
  --by developer \
  --in-completion "$RAISE_COMPLETION_ID" \
  --summary "Human dogfood concern" \
  --body-md "$FID_HV9_WORKSPACE/concern-body.md" 2>&1)
printf '%s\n' "$RAISE_OUTPUT"

CONCERN_ID=$(printf '%s\n' "$RAISE_OUTPUT" | sed -n 's/^raised concern \([^ ]*\).*/\1/p')
printf 'concern_id=%s\n' "$CONCERN_ID"

"$FID" concern list --issue "$TASK_ID" --format json
"$FID" issue get "$TASK_ID" --section concerns --format json
```

Expected for `HUMAN-011a`:

- Completion output includes a `completion: TC-...` line.
- Concern raise output says `raised concern CON-... on <task>`.
- `concern_id` starts with `CON-`.
- Both read-back commands show the concern as `status = open`.
- The visible row has the expected summary, role, issue link, raised completion
  link, and body path.
- `addressed_in_completion_id` is present and `null` while the concern is open.
- The output reads as the concern raised against the intended issue/completion.

Record `HUMAN-011a` as `pass` if the concern is visible and understandable.
Record `partial` or `fail` with notes if the concern is hard to identify, a
field is missing, or the output suggests the wrong issue/completion.

## 5. HUMAN-011b - Confirm The Open Concern Is Obvious On Issue Reads

```bash
"$FID" issue get "$TASK_ID" --format json
"$FID" issue list --format json
```

Expected for the after-raise side of `HUMAN-011b`:

- `concerns_raised` is `true`.
- `open_concerns_count` is `1`.
- The change from the baseline false/0 state is clear on both `issue get` and
  `issue list`.
- The output makes unresolved concern work obvious enough during dogfooding.

Record `HUMAN-011b` as `pass` if both issue reads make the open concern state
clear. Record `partial` or `fail` if the state is technically present but easy
to miss or confusing.

## 6. HUMAN-011c - Address The Concern From A Later Completion

Create a later completion, then address the same concern from that completion.

```bash
COMPLETE2=$("$FID" issue complete "$TASK_ID" --by developer 2>&1)
printf '%s\n' "$COMPLETE2"

ADDRESS_COMPLETION_ID=$(printf '%s\n' "$COMPLETE2" | sed -n 's/^[[:space:]]*completion: \(TC-[^ ]*\).*/\1/p')
printf 'address_completion_id=%s\n' "$ADDRESS_COMPLETION_ID"

ADDRESS_OUTPUT=$("$FID" concern address "$CONCERN_ID" \
  --in-completion "$ADDRESS_COMPLETION_ID" \
  --by developer 2>&1)
printf '%s\n' "$ADDRESS_OUTPUT"

"$FID" concern list --issue "$TASK_ID" --format json
"$FID" issue get "$TASK_ID" --section concerns --format json
```

Expected for `HUMAN-011c`:

- Address output says `addressed concern <CON> in completion <TC>`.
- The same concern ID remains visible; it is not replaced by a new unrelated
  concern row.
- `status` changes from `open` to `addressed`.
- `addressed_in_completion_id` is the second completion ID.
- The lifecycle reads as an in-place resolution of the original concern.

Record `HUMAN-011c` as `pass` if the address flow reads like a clear resolution.
Record `partial` or `fail` if it feels like a duplicate concern was created, the
addressing completion is hard to see, or the status transition is unclear.

## 7. HUMAN-011d - Confirm The Last Open Concern Clears From Issue Reads

```bash
"$FID" issue get "$TASK_ID" --format json
"$FID" issue list --format json
"$FID" concern list --issue "$TASK_ID" --status open --format json
```

Expected for `HUMAN-011d`:

- `issue get` shows `concerns_raised = false`.
- `issue get` shows `open_concerns_count = 0`.
- `issue list` shows the same cleared state for the task.
- `concern list --status open` returns an empty collection, not an error.
- The issue reads no longer imply unresolved concern work remains.

Record `HUMAN-011d` as `pass` if the cleared state is obvious. Record `partial`
or `fail` if the issue still appears to have unresolved concerns or the empty
open-list behavior is confusing.

## 8. HUMAN-011e - Exercise Concern List Filters

Use issue and status filters against the dogfood concern.

```bash
"$FID" concern list --issue "$TASK_ID" --format json
"$FID" concern list --issue "$TASK_ID" --status addressed --format json
"$FID" concern list --issue "$TASK_ID" --status open --format json
"$FID" concern list --status addressed --format json
```

Expected for the main filter views:

- `--issue "$TASK_ID"` returns the concern for this task.
- `--issue "$TASK_ID" --status addressed` returns the addressed concern.
- `--issue "$TASK_ID" --status open` returns an empty collection with exit 0.
- `--status addressed` returns the addressed concern and includes `issue_id`,
  making the row useful outside an issue-specific view.
- The filtered views match the way a developer would look up concern state while
  dogfooding.

Optionally inspect the same filtered view in every supported format:

```bash
"$FID" concern list --issue "$TASK_ID" --status addressed --format json
"$FID" concern list --issue "$TASK_ID" --status addressed --format toon
"$FID" concern list --issue "$TASK_ID" --status addressed --format markdown
"$FID" concern list --issue "$TASK_ID" --status addressed --format yaml
```

Expected for the format sweep:

- Each format exits zero.
- Each format clearly identifies the same addressed concern.
- Empty-result behavior remains understandable if you repeat the sweep with
  `--status open`.

Record `HUMAN-011e` as `pass` if the filters and empty-result behavior are useful
and not misleading. Record `partial` or `fail` with the confusing command/output
if lookup friction remains.

## 9. Optional Rebuild Read-Back Check

This is not required to fill the human checklist, but it is a useful confidence
check while the disposable workspace is already populated.

```bash
CONCERNS_BEFORE_REBUILD=$("$FID" concern list --issue "$TASK_ID" --format json)
ISSUE_BEFORE_REBUILD=$("$FID" issue get "$TASK_ID" --format json)
ISSUE_CONCERNS_BEFORE_REBUILD=$("$FID" issue get "$TASK_ID" --section concerns --format json)

mv .fid/fid.sqlite .fid/fid.sqlite.before-human-rebuild
"$FID" rebuild --by developer

CONCERNS_AFTER_REBUILD=$("$FID" concern list --issue "$TASK_ID" --format json)
ISSUE_AFTER_REBUILD=$("$FID" issue get "$TASK_ID" --format json)
ISSUE_CONCERNS_AFTER_REBUILD=$("$FID" issue get "$TASK_ID" --section concerns --format json)

if [ "$CONCERNS_BEFORE_REBUILD" = "$CONCERNS_AFTER_REBUILD" ] &&
   [ "$ISSUE_BEFORE_REBUILD" = "$ISSUE_AFTER_REBUILD" ] &&
   [ "$ISSUE_CONCERNS_BEFORE_REBUILD" = "$ISSUE_CONCERNS_AFTER_REBUILD" ]; then
  printf 'human_rebuild_parity=OK\n'
else
  printf 'human_rebuild_parity=DIFF\n'
  printf 'concerns before:\n%s\nconcerns after:\n%s\n' \
    "$CONCERNS_BEFORE_REBUILD" \
    "$CONCERNS_AFTER_REBUILD"
  printf 'issue before:\n%s\nissue after:\n%s\n' \
    "$ISSUE_BEFORE_REBUILD" \
    "$ISSUE_AFTER_REBUILD"
  printf 'issue concerns before:\n%s\nissue concerns after:\n%s\n' \
    "$ISSUE_CONCERNS_BEFORE_REBUILD" \
    "$ISSUE_CONCERNS_AFTER_REBUILD"
fi
```

Expected:

- `fid rebuild` reports `rebuild complete`.
- `human_rebuild_parity=OK`.

## 10. Optional Repo-Level Regression Commands

Run these from the repo root, not from the disposable workspace. They are not
needed to fill the human checklist, but they provide useful final confidence
while Phase 9 is being closed.

```bash
cd /home/mark/projects/fid

gofmt -l cmd internal
go test -tags integration -race ./cmd/fid -run 'TestV01C(ConcernRebuildParity|ClosureSmoke)' -count=1
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

Expected:

- `gofmt -l cmd internal` prints nothing.
- The targeted v0.1c concern integration tests pass.
- The benchmark command exits zero, including the previously blocked concern
  list test path.
- The build command exits zero.

## 11. Fill The Human Checklist

Open `project-docs/docs/phases/phase-9/qa/human-checklist.md` and record one
outcome per row:

- `HUMAN-011a`: use Section 4.
- `HUMAN-011b`: use Sections 3 and 5.
- `HUMAN-011c`: use Section 6.
- `HUMAN-011d`: use Section 7.
- `HUMAN-011e`: use Section 8.

For each row, record `pass`, `partial`, or `fail` plus enough notes to explain
the product judgment. `VAL-011` can pass only after these human outcomes are
complete.
