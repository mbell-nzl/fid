# Phase 8 Human Checklist Run Sheet

> Phase: 8 - Task Completions
> Purpose: executable command sheet for
> `project-docs/docs/phases/phase-8/qa/human-checklist.md`
> Date prepared: 2026-06-20

Use this from the fid repository root. It creates a disposable `.fid/`
workspace under `/tmp`, builds the current CLI to `/tmp/fid`, creates one task,
then walks the human checklist rows for `VAL-008`.

The human checklist still needs human judgment. Treat the command output as
evidence, then record the outcome and any friction in
`human-checklist.md`.

## Prerequisites

- Run the build command from the repo root.
- Run workspace commands only inside the disposable `/tmp` workspace created
  below, not inside this repository's own `.fid/` workspace.
- These commands require the built `fid` binary and POSIX shell tools already
  used by the project docs (`sed`, `wc`, `ls`, `sort`, `tail`, `mv`, `printf`).
- The gate snapshot step appends a fixture journal event for `gate_runs`, then
  rebuilds the disposable SQLite DB. This matches Phase 8 DEC-C: the real
  `fid dev gate` writer is v0.1e, so Phase 8 proves snapshotting with
  fixture-seeded gates.
- The rebuild steps rename disposable SQLite DB files inside `/tmp`; they do
  not touch repo data.

## 1. Build Current CLI

```bash
go build -o /tmp/fid ./cmd/fid
export FID=/tmp/fid
```

Expected: command exits zero and `/tmp/fid` exists.

## 2. Create A Fresh Test Workspace

```bash
export FID_HV8_WORKSPACE="$(mktemp -d /tmp/fid-phase8-human-XXXXXX)"
cd "$FID_HV8_WORKSPACE"
printf 'workspace=%s\n' "$FID_HV8_WORKSPACE"

"$FID" init --prefix HV8

TASK_CREATE=$("$FID" issue create \
  --kind task \
  --title "HV8 completion lifecycle dogfood" \
  --by developer 2>&1)
TASK_ID=$(printf '%s\n' "$TASK_CREATE" | sed -n 's/^created issue \([^ ]*\).*/\1/p')
printf 'task_id=%s\n' "$TASK_ID"
```

Expected:

- `fid init` reports a new `.fid/` workspace.
- `task_id` starts with `HV8-`.

## 3. HV-08e Reject Path - Gated Transition Before Completion

Move the task into progress, then try the `has_completion_report`-gated
`dev_complete` transition before any completion exists.

```bash
"$FID" issue transition "$TASK_ID" --to in_progress --by developer

JOURNAL=$(ls .fid/journal/*.jsonl | sort | tail -n 1)
JOURNAL_LINES_BEFORE_REJECT=$(wc -l < "$JOURNAL")
ISSUE_BEFORE_REJECT=$("$FID" issue get "$TASK_ID" --format json)

REJECT_EXIT=0
REJECT_OUTPUT=$("$FID" issue transition "$TASK_ID" \
  --to dev_complete \
  --by developer 2>&1) || REJECT_EXIT=$?

JOURNAL_LINES_AFTER_REJECT=$(wc -l < "$JOURNAL")
ISSUE_AFTER_REJECT=$("$FID" issue get "$TASK_ID" --format json)

printf 'reject_exit=%s\n' "$REJECT_EXIT"
printf '%s\n' "$REJECT_OUTPUT"
printf 'journal_lines_before=%s\njournal_lines_after=%s\n' \
  "$JOURNAL_LINES_BEFORE_REJECT" \
  "$JOURNAL_LINES_AFTER_REJECT"

if [ "$ISSUE_BEFORE_REJECT" = "$ISSUE_AFTER_REJECT" ]; then
  printf 'reject_state_parity=OK\n'
else
  printf 'reject_state_parity=DIFF\n'
fi
```

Expected for the first half of `HV-08e`:

- `reject_exit=1`.
- The error names `has_completion_report` and the task ID.
- `journal_lines_before` equals `journal_lines_after`.
- `reject_state_parity=OK`.
- The rejection reads as understandable and actionable.

## 4. HV-08a - Record Attempt 1 And Read Back Fields

```bash
ATTEMPT1_OUTPUT=$("$FID" issue complete "$TASK_ID" \
  --by developer \
  --notes-md project-docs/docs/phases/phase-8/tasks/task-1-done.md \
  --tokens 14823 \
  --commits abc123,def456 \
  --files internal/issues/complete.go \
  --files cmd/fid/issue.go 2>&1)
printf '%s\n' "$ATTEMPT1_OUTPUT"

"$FID" issue get "$TASK_ID" --section completions --format json
```

Expected for `HV-08a`:

- Completion confirmation says `completed <task> (attempt 1)`.
- Completion ID starts with `TC-`.
- Read-back shows `developer_role`, `actual_tokens`, `commit_hashes`,
  `files_touched`, and `notes_md_path`.
- Field names and values are understandable without needing to inspect the DB.

Record `HV-08a` as `pass` if the confirmation and read-back fields are clear.
Record `partial` or `fail` with notes if a field is missing, confusing, or
presented poorly.

## 5. HV-08e Pass Path - Gated Transition After Completion

```bash
TRANSITION_AFTER_COMPLETION=$("$FID" issue transition "$TASK_ID" \
  --to dev_complete \
  --by developer 2>&1)
printf '%s\n' "$TRANSITION_AFTER_COMPLETION"

"$FID" issue get "$TASK_ID" --format json
```

Expected for the second half of `HV-08e`:

- Transition confirmation says `in_progress -> dev_complete` or the CLI's
  arrow-rendered equivalent.
- The task read-back shows `status = dev_complete`.
- The reject-before-completion / pass-after-completion lifecycle feels coherent.

Record `HV-08e` as `pass` if both halves of the lifecycle behave clearly.

## 6. Seed Gate Runs For Attempt 2

There is no real `fid dev gate` writer in Phase 8. This step appends two
fixture `gate_runs` rows for attempt 2, rebuilds the disposable projection, and
then the next completion snapshots those rows.

```bash
JOURNAL=$(ls .fid/journal/*.jsonl | sort | tail -n 1)
JOURNAL_DATE=$(printf '%s\n' "$JOURNAL" | sed -n 's|.*/\([0-9][0-9-]*\)\.jsonl$|\1|p')
EVENT_DATE=$(printf '%s' "$JOURNAL_DATE" | sed 's/-//g')
EVENT_SEQ=$(( $(wc -l < "$JOURNAL") + 1 ))
EVENT_ID=$(printf 'EVT-%s-%06d' "$EVENT_DATE" "$EVENT_SEQ")
OCCURRED_AT="${JOURNAL_DATE}T12:00:00Z"

printf '{"event_id":"%s","event_schema":1,"event_type":"gate.seeded","occurred_at":"%s","actor":{"role":"developer","actor":null},"command":{"argv":["fid","maintenance","seed"],"cwd":".","fid_version":"0.1.0"},"effects":[{"op":"insert","table":"gate_runs","pk":{"id":"GATE-HV8-001"},"values":{"issue_id":"%s","attempt_number":2,"recorded_by_role":"developer","command":"go test -race ./...","exit_code":0,"output_excerpt":"ok github.com/mbell-nzl/fid/cmd/fid 0.321s","ran_at":"2026-06-20T09:00:00Z"}},{"op":"insert","table":"gate_runs","pk":{"id":"GATE-HV8-002"},"values":{"issue_id":"%s","attempt_number":2,"recorded_by_role":"developer","command":"golangci-lint run ./...","exit_code":1,"output_excerpt":"cmd/fid/e2e_v01c_complete_test.go: line exceeds budget","ran_at":"2026-06-20T09:01:00Z"}}]}\n' \
  "$EVENT_ID" \
  "$OCCURRED_AT" \
  "$TASK_ID" \
  "$TASK_ID" >> "$JOURNAL"

mv .fid/fid.sqlite .fid/fid.sqlite.before-gate-seed
"$FID" rebuild --by developer
```

Expected:

- `fid rebuild` reports `rebuild complete`.
- No repo data is touched; all files are under `$FID_HV8_WORKSPACE`.

## 7. HV-08b, HV-08c, HV-08d - Record Attempt 2 With Gates And Intervention

```bash
ATTEMPT2_OUTPUT=$("$FID" issue complete "$TASK_ID" \
  --by developer \
  --notes-md project-docs/docs/phases/phase-8/tasks/task-3-done.md \
  --tokens 1700 \
  --human-intervention \
  --intervention-note "paired with qa on rebuild parity" 2>&1)
printf '%s\n' "$ATTEMPT2_OUTPUT"

"$FID" issue get "$TASK_ID" --section completions --format json
```

Expected for `HV-08b`:

- Completion confirmation says `completed <task> (attempt 2)`.
- Read-back shows two retained completion rows.
- Attempt 2 is the latest active `dev_complete` attempt.
- Attempt 1 remains present with `status = superseded`.
- Attempt numbering and supersession are understandable.

Expected for `HV-08c`:

- Attempt 2 shows `gates_passed = 1` and `gates_failed = 1`.
- Attempt 2 `gate_results` includes both fixture commands:
  `go test -race ./...` and `golangci-lint run ./...`.
- Each gate result shows command, exit code, run time, and output excerpt.
- The rolled-up gate counts and details read as expected.

Expected for `HV-08d`:

- Attempt 2 shows `human_intervention = true`.
- Attempt 2 shows `human_intervention_note = "paired with qa on rebuild parity"`.
- Attempt 1 still shows no intervention flag or note.
- The intervention presentation is clear and expected.

Record `HV-08b`, `HV-08c`, and `HV-08d` separately. A gate snapshot can pass
even if intervention wording has friction, and vice versa.

## 8. Optional Rebuild Read-Back Check

This is not required to fill the human checklist, but it is a useful confidence
check while the disposable workspace is already populated.

```bash
COMPLETIONS_BEFORE_REBUILD=$("$FID" issue get "$TASK_ID" --section completions --format json)

mv .fid/fid.sqlite .fid/fid.sqlite.before-final-rebuild
"$FID" rebuild --by developer

COMPLETIONS_AFTER_REBUILD=$("$FID" issue get "$TASK_ID" --section completions --format json)

if [ "$COMPLETIONS_BEFORE_REBUILD" = "$COMPLETIONS_AFTER_REBUILD" ]; then
  printf 'completion_rebuild_parity=OK\n'
else
  printf 'completion_rebuild_parity=DIFF\n'
  printf 'before:\n%s\nafter:\n%s\n' \
    "$COMPLETIONS_BEFORE_REBUILD" \
    "$COMPLETIONS_AFTER_REBUILD"
fi
```

Expected:

- `fid rebuild` reports `rebuild complete`.
- `completion_rebuild_parity=OK`.

## 9. Optional Repo-Level Regression Commands

Run these from the repo root, not from the disposable workspace. They are not
needed to fill the human checklist, but they provide useful final confidence
while Phase 8 is being closed.

```bash
cd /home/mark/projects/fid

gofmt -l cmd internal
go test -tags integration -race ./cmd/fid -run TestV01CHasCompletionReportGateAndRebuildParity -count=1
go build -o /tmp/fid ./cmd/fid
```

Expected:

- `gofmt -l cmd internal` prints nothing.
- The integration smoke passes.
- The build command exits zero.

For the full phase-level gate, use the command list in `AGENTS.md` or
`project-docs/docs/phases/phase-8/qa/qa-plan.md`.

## 10. Recording Results In The Human Checklist

Suggested notes to copy into `human-checklist.md` after running:

- `HV-08a`: `pass` if attempt 1 creation and read-back clearly show role,
  tokens, commits, files, and notes path.
- `HV-08b`: `pass` if attempt 2 and attempt 1 supersession are clear and both
  rows remain visible.
- `HV-08c`: `pass` if gate counts and gate details on attempt 2 read clearly.
- `HV-08d`: `pass` if the human-intervention flag and note are visible and
  understandable.
- `HV-08e`: `pass` if the transition rejects before completion, mutates no
  state, and succeeds after the real completion exists.

If all five rows pass, update the summary totals to `Pass = 5`, `Blank = 0`,
and `Overall result = pass`. If any row is `partial` or `fail`, keep the exact
friction in that row's notes and set the overall result according to the
checklist footer rules.
