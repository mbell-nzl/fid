# Phase 7 Human Checklist Run Sheet

> Phase: 7 - Acceptance Criteria Surface
> Purpose: executable command sheet for
> `project-docs/docs/phases/phase-7/qa/human-checklist.md`
> Date prepared: 2026-06-15

Use this from the fid repository root. It creates a disposable `.fid/`
workspace under `/tmp`, builds the current CLI to `/tmp/fid`, seeds one task and
one epic, then walks the human checklist rows plus a few lightweight regression
checks.

The human checklist still needs human judgment. Treat the command output as
evidence, then record the outcome and any friction in
`human-checklist.md`.

## Prerequisites

- Run these commands from the repo root.
- Do not run the workspace commands inside this repository's own `.fid/`
  workspace.
- These commands require only the built `fid` binary and POSIX shell tools
  already used by the project docs (`sed`, `wc`, `cmp`-free shell comparison).
- The rebuild check renames the disposable SQLite DB to
  `.fid/fid.sqlite.before-rebuild`; it does not delete any repo data.

## 1. Build Current CLI

```bash
go build -o /tmp/fid ./cmd/fid
export FID=/tmp/fid
```

Expected: command exits zero and `/tmp/fid` exists.

## 2. Create A Fresh Test Workspace

```bash
set -euo pipefail

export FID_HV7_WORKSPACE="$(mktemp -d /tmp/fid-phase7-human-XXXXXX)"
cd "$FID_HV7_WORKSPACE"
printf 'workspace=%s\n' "$FID_HV7_WORKSPACE"

"$FID" init --prefix HV7

TASK_CREATE=$("$FID" issue create \
  --kind task \
  --title "HV7 acceptance criteria dogfood" \
  --by developer 2>&1)
TASK_ID=$(printf '%s\n' "$TASK_CREATE" | sed -n 's/^created issue \([^ ]*\).*/\1/p')
printf 'task_id=%s\n' "$TASK_ID"

EPIC_CREATE=$("$FID" issue create \
  --kind epic \
  --title "HV7 epic rejection dogfood" \
  --scheduling sequential \
  --by tl 2>&1)
EPIC_ID=$(printf '%s\n' "$EPIC_CREATE" | sed -n 's/^created issue \([^ ]*\).*/\1/p')
printf 'epic_id=%s\n' "$EPIC_ID"
```

Expected:

- `fid init` reports a new `.fid/` workspace.
- `task_id` starts with `HV7-`.
- `epic_id` starts with `HV7-`.

## 3. HV-010a - Add Ordered ACs And Read Back

Add the ordinal 2 row first, then ordinal 1, so the read-back proves ordering is
by ordinal rather than append order.

```bash
ADD_SECOND=$("$FID" issue ac add "$TASK_ID" \
  --text "Second AC should remain second" \
  --ordinal 2 \
  --by developer 2>&1)
printf '%s\n' "$ADD_SECOND"
AC2_ID=$(printf '%s\n' "$ADD_SECOND" | sed -n 's/^added AC \([^ ]*\).*/\1/p')

ADD_FIRST=$("$FID" issue ac add "$TASK_ID" \
  --text "First AC should read first" \
  --ordinal 1 \
  --by developer 2>&1)
printf '%s\n' "$ADD_FIRST"
AC1_ID=$(printf '%s\n' "$ADD_FIRST" | sed -n 's/^added AC \([^ ]*\).*/\1/p')

printf 'ac1_id=%s\nac2_id=%s\n' "$AC1_ID" "$AC2_ID"

"$FID" issue ac list "$TASK_ID" --format json
"$FID" issue get "$TASK_ID" --section acceptance --format json
```

Expected for `HV-010a`:

- Both add confirmations are clear.
- `ac1_id` and `ac2_id` start with `AC-`.
- The list/read-back output shows ordinal `1` before ordinal `2`.
- The text values match the two commands.

Record `HV-010a` as `pass` if that is clear and the ordering behaves as
expected. Record `partial` or `fail` with notes if the wording or ordering is
confusing.

## 4. HV-010b - Multi-Role Verify And Latest Per Role

```bash
"$FID" issue ac verify "$TASK_ID" "$AC1_ID" \
  --by developer \
  --status pass \
  --evidence-type manual \
  --evidence-ref human-checklist:HV-010b \
  --notes "Developer role accepts the first AC"

"$FID" issue ac verify "$TASK_ID" "$AC1_ID" \
  --by qa_lead \
  --status fail \
  --evidence-type manual \
  --evidence-ref human-checklist:HV-010b \
  --notes "QA role records an independent failure"

"$FID" issue ac list "$TASK_ID" --format json
```

Expected for `HV-010b`:

- Both verification confirmations are clear.
- The first AC shows `latest_verifications.developer.status = pass`.
- The same AC also shows `latest_verifications.qa_lead.status = fail`.
- The two roles are understandable as independent latest statuses.

Record `HV-010b` as `pass` if the role labels and latest-per-role result match
your expectation. Record friction or ambiguity in the notes.

## 5. HV-010c - Epic Rejection Clarity

This also checks the validate-before-append behavior by comparing journal line
counts before and after the rejected command.

```bash
JOURNAL_LINES_BEFORE_EPIC_REJECT=$(wc -l .fid/journal/*.jsonl | sed -n 's/^ *\([0-9][0-9]*\).*/\1/p')

set +e
EPIC_REJECT_OUTPUT=$("$FID" issue ac add "$EPIC_ID" \
  --text "This AC must be rejected because the target is an epic" \
  --ordinal 1 \
  --by developer 2>&1)
EPIC_REJECT_EXIT=$?
set -e

JOURNAL_LINES_AFTER_EPIC_REJECT=$(wc -l .fid/journal/*.jsonl | sed -n 's/^ *\([0-9][0-9]*\).*/\1/p')

printf 'epic_reject_exit=%s\n' "$EPIC_REJECT_EXIT"
printf '%s\n' "$EPIC_REJECT_OUTPUT"
printf 'journal_lines_before=%s\njournal_lines_after=%s\n' \
  "$JOURNAL_LINES_BEFORE_EPIC_REJECT" \
  "$JOURNAL_LINES_AFTER_EPIC_REJECT"
```

Expected for `HV-010c`:

- `epic_reject_exit=1`.
- The error says the issue is kind `epic` and that ACs require a task or
  subtask.
- `journal_lines_before` equals `journal_lines_after`.

Record `HV-010c` as `pass` if the rejection is clear and actionable. Record the
confusing text or behavior if not.

## 6. Additional Regression Checks

These are not human-only checklist rows, but they are useful to run while the
fresh workspace is already populated.

### Default JSON Matches Explicit JSON

```bash
AC_LIST_DEFAULT=$("$FID" issue ac list "$TASK_ID")
AC_LIST_JSON=$("$FID" issue ac list "$TASK_ID" --format json)

if [ "$AC_LIST_DEFAULT" = "$AC_LIST_JSON" ]; then
  printf 'default_json_parity=OK\n'
else
  printf 'default_json_parity=DIFF\n'
  printf 'default:\n%s\njson:\n%s\n' "$AC_LIST_DEFAULT" "$AC_LIST_JSON"
fi
```

Expected: `default_json_parity=OK`.

### `ac list` Formats Render

```bash
"$FID" issue ac list "$TASK_ID" --format json
"$FID" issue ac list "$TASK_ID" --format toon
"$FID" issue ac list "$TASK_ID" --format markdown
"$FID" issue ac list "$TASK_ID" --format yaml
```

Expected: all four commands exit zero and render readable output.

### `issue get --section acceptance` Mirrors `ac list`

```bash
AC_LIST_JSON=$("$FID" issue ac list "$TASK_ID" --format json)
ISSUE_ACCEPTANCE_JSON=$("$FID" issue get "$TASK_ID" --section acceptance --format json)

if [ "$AC_LIST_JSON" = "$ISSUE_ACCEPTANCE_JSON" ]; then
  printf 'acceptance_section_parity=OK\n'
else
  printf 'acceptance_section_parity=DIFF\n'
  printf 'ac list:\n%s\nissue get acceptance:\n%s\n' \
    "$AC_LIST_JSON" \
    "$ISSUE_ACCEPTANCE_JSON"
fi
```

Expected: `acceptance_section_parity=OK`.

### Invalid Status Rejects Before Append

```bash
JOURNAL_LINES_BEFORE_BAD_STATUS=$(wc -l .fid/journal/*.jsonl | sed -n 's/^ *\([0-9][0-9]*\).*/\1/p')

set +e
BAD_STATUS_OUTPUT=$("$FID" issue ac verify "$TASK_ID" "$AC1_ID" \
  --by developer \
  --status bogus 2>&1)
BAD_STATUS_EXIT=$?
set -e

JOURNAL_LINES_AFTER_BAD_STATUS=$(wc -l .fid/journal/*.jsonl | sed -n 's/^ *\([0-9][0-9]*\).*/\1/p')

printf 'bad_status_exit=%s\n' "$BAD_STATUS_EXIT"
printf '%s\n' "$BAD_STATUS_OUTPUT"
printf 'journal_lines_before=%s\njournal_lines_after=%s\n' \
  "$JOURNAL_LINES_BEFORE_BAD_STATUS" \
  "$JOURNAL_LINES_AFTER_BAD_STATUS"
```

Expected:

- `bad_status_exit=1`.
- Error names the allowed statuses: `pass`, `fail`, `blocked`,
  `not_applicable`.
- Journal line counts are unchanged.

### Rebuild Preserves AC Read-Back

```bash
AC_LIST_BEFORE_REBUILD=$("$FID" issue ac list "$TASK_ID" --format json)
ISSUE_ACCEPTANCE_BEFORE_REBUILD=$("$FID" issue get "$TASK_ID" --section acceptance --format json)

mv .fid/fid.sqlite .fid/fid.sqlite.before-rebuild
"$FID" rebuild --by developer

AC_LIST_AFTER_REBUILD=$("$FID" issue ac list "$TASK_ID" --format json)
ISSUE_ACCEPTANCE_AFTER_REBUILD=$("$FID" issue get "$TASK_ID" --section acceptance --format json)

if [ "$AC_LIST_BEFORE_REBUILD" = "$AC_LIST_AFTER_REBUILD" ]; then
  printf 'ac_list_rebuild_parity=OK\n'
else
  printf 'ac_list_rebuild_parity=DIFF\n'
fi

if [ "$ISSUE_ACCEPTANCE_BEFORE_REBUILD" = "$ISSUE_ACCEPTANCE_AFTER_REBUILD" ]; then
  printf 'issue_acceptance_rebuild_parity=OK\n'
else
  printf 'issue_acceptance_rebuild_parity=DIFF\n'
fi
```

Expected:

- `fid rebuild` reports `rebuild complete`.
- `ac_list_rebuild_parity=OK`.
- `issue_acceptance_rebuild_parity=OK`.

## 7. Optional Repo-Level Regression Commands

Run these from the repo root, not from the disposable workspace. They are not
needed to fill the human checklist, but they provide useful final confidence
while Phase 7 is being closed.

```bash
cd /home/mark/projects/fid

gofmt -l cmd internal
go test -tags integration -race ./cmd/fid -run TestV01CAllACsVerifiedGateAndRebuildParity -count=1
go build -o /tmp/fid ./cmd/fid
```

Expected:

- `gofmt -l cmd internal` prints nothing.
- The integration smoke passes.
- The build command exits zero.

For the full phase-level gate, use the command list in `AGENTS.md` or
`project-docs/docs/phases/phase-7/qa/qa-plan.md`.

## 8. Recording Results In The Human Checklist

Suggested notes to copy into `human-checklist.md` after running:

- `HV-010a`: `pass` if the add confirmations were clear and read-back ordered
  ordinal 1 before ordinal 2.
- `HV-010b`: `pass` if the developer and qa_lead latest statuses were both
  visible and understandable on the same AC.
- `HV-010c`: `pass` if the epic rejection was clear/actionable and the journal
  line count did not change.

If all three rows pass, update the summary totals to `Pass = 3`, `Blank = 0`,
and `Overall result = pass`. If any row is `partial` or `fail`, keep the exact
friction in that row's notes and set the overall result according to the
checklist footer rules.
