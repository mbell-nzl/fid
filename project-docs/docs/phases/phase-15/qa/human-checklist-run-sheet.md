# Phase 15 Human Checklist Run Sheet

> Phase: 15 - `fid dev *` Composites + Cleanup
> Purpose: executable command sheet for
> `project-docs/docs/phases/phase-15/qa/human-checklist.md`
> Date prepared: 2026-06-29

Use this from the fid repository root. It creates a disposable real `.fid/`
workspace under `/tmp`, builds the current CLI to `/tmp/fid`, creates one task
for the finish path and one task for the block path, then runs the developer
composite commands needed for the human dogfood checks.

The human checklist still needs human judgment. Treat the command output as
dogfood evidence, then record the outcome and notes in `human-checklist.md`.

## Prerequisites

- Run the build command from the repo root.
- Run workspace commands only inside the disposable `/tmp` workspace created
  below, not inside this repository's own `.fid/` workspace.
- These commands require the built `fid` binary and POSIX shell tools already
  used by the project docs (`sed`, `printf`, `ls`, `mkdir`, `mktemp`).
- This sheet does not complete deterministic QA gates. Phase 15 final review
  already covers one-event composite writes, zero-append rejection paths,
  permission parity, completion snapshots, rebuild determinism, and phase-level
  gates. This sheet only supports `HV-15a` through `HV-15d` for `VAL-012`.

## 1. Build Current CLI

```bash
go build -o /tmp/fid ./cmd/fid
export FID=/tmp/fid
ls -l "$FID"
```

Expected: command exits zero and `/tmp/fid` exists.

## 2. Create A Fresh Test Workspace

```bash
export FID_HV15_WORKSPACE="$(mktemp -d /tmp/fid-phase15-human-XXXXXX)"
cd "$FID_HV15_WORKSPACE"
printf 'workspace=%s\n' "$FID_HV15_WORKSPACE"

"$FID" init --prefix HV15
mkdir -p reports
printf '# HV15 finish report\n' > reports/dev-finish.md
printf '# HV15 block reason\n\nGate failed during dogfood.\n' > reports/dev-block.md
```

Expected:

- `fid init` reports a new `.fid/` workspace.
- The workspace has a `.fid/` directory.
- The rest of this sheet is run from the disposable workspace.
- The `reports/` directory contains `dev-finish.md` and `dev-block.md`.

## 3. Create Issues For Composite Runs

Create one task for the `start -> gate -> finish` path and one task for the
`start -> block` path. Add one acceptance criterion to the finish-path task so
`dev finish` can satisfy the developer AC rule.

```bash
FINISH_OUTPUT=$("$FID" issue create \
  --kind task \
  --title "HV15 finish path" \
  --priority 1 \
  --est-tokens 800 \
  --by developer 2>&1)
printf '%s\n' "$FINISH_OUTPUT"
export HV15_FINISH_TASK="$(printf '%s\n' "$FINISH_OUTPUT" | sed -n 's/^created issue \([^ ]*\) .*/\1/p')"
printf 'finish_task=%s\n' "$HV15_FINISH_TASK"

BLOCK_OUTPUT=$("$FID" issue create \
  --kind task \
  --title "HV15 block path" \
  --priority 2 \
  --est-tokens 500 \
  --by developer 2>&1)
printf '%s\n' "$BLOCK_OUTPUT"
export HV15_BLOCK_TASK="$(printf '%s\n' "$BLOCK_OUTPUT" | sed -n 's/^created issue \([^ ]*\) .*/\1/p')"
printf 'block_task=%s\n' "$HV15_BLOCK_TASK"

AC_OUTPUT=$("$FID" issue ac add "$HV15_FINISH_TASK" \
  --text "Developer gate passes before finish" \
  --ordinal 1 \
  --by developer 2>&1)
printf '%s\n' "$AC_OUTPUT"
export HV15_FINISH_AC="$(printf '%s\n' "$AC_OUTPUT" | sed -n 's/^added AC \([^ ]*\) .*/\1/p')"
printf 'finish_ac=%s\n' "$HV15_FINISH_AC"
```

Expected:

- Both tasks are created.
- `HV15_FINISH_TASK`, `HV15_BLOCK_TASK`, and `HV15_FINISH_AC` print non-empty
  IDs.

## 4. HV-15a - Start A Developer Composite Run

Run `fid dev start` on the finish-path task. Inspect the stderr confirmation as
a human operator.

```bash
"$FID" dev start "$HV15_FINISH_TASK" \
  --model claude-sonnet-4-6 \
  --by developer
```

Expected for `HV-15a`:

- The command exits zero.
- The output starts with `dev start HV15-...`.
- The output says `claimed by developer`.
- The output says `transitioned: todo -> in_progress` with the CLI's arrow
  formatting.
- The output shows `agent run: RUN-... opened`.

Record `HV-15a` as `pass` if the output makes it clear that the task was
claimed by the developer, moved to `in_progress`, and an agent run was opened.
Record `partial` or `fail` if any of those outcomes is missing or only
understandable by inspecting implementation files or the database.

## 5. HV-15b - Record A Reported Gate Result

Run `fid dev gate` on the same task. The command string is deliberately a
reported command; fid should record the result supplied by the operator.

```bash
"$FID" dev gate "$HV15_FINISH_TASK" \
  --cmd "go test ./..." \
  --exit-code 0 \
  --output-excerpt "ok" \
  --by developer
```

Expected for `HV-15b`:

- The command exits zero.
- The output says `gate recorded on HV15-... (attempt 1)`.
- The output shows `cmd: go test ./...`.
- The output shows `exit_code: 0`.
- The output shows `gates_passed: 1, gates_failed: 0`.
- The output reads like a record of the reported result, not like fid executed
  the gate command itself.

Record `HV-15b` as `pass` if the output makes it clear that fid recorded the
reported gate result and does not appear to have executed the command. Record
`partial` or `fail` if the output is ambiguous about recording versus executing
the command.

## 6. Verify Developer ACs For The Finish Path

This prerequisite is not a human checklist row. It lets `dev finish` pass the
developer AC transition rule.

```bash
"$FID" issue ac verify "$HV15_FINISH_TASK" "$HV15_FINISH_AC" \
  --status pass \
  --by developer \
  --evidence-type command \
  --evidence-ref "fid dev gate $HV15_FINISH_TASK --cmd 'go test ./...' --exit-code 0"
```

Expected:

- The command exits zero.
- The output says `verified AC-... on HV15-... by developer: pass`.

## 7. HV-15c - Finish The Developer Composite Run

Run `fid dev finish` on the finish-path task using the disposable completion
report created in section 2. Inspect the stderr confirmation as a human
operator.

```bash
"$FID" dev finish "$HV15_FINISH_TASK" \
  --completion-md reports/dev-finish.md \
  --tokens 14823 \
  --commits abc123 \
  --files project-docs/docs/phases/phase-15/qa/human-checklist-run-sheet.md \
  --by developer
```

Expected for `HV-15c`:

- The command exits zero.
- The output starts with `dev finish HV15-...`.
- The output shows `agent run: RUN-... closed`, with the token count if
  provided.
- The output shows `completion: TC-... (attempt 1)`.
- The output says `transitioned: in_progress -> dev_complete` with the CLI's
  arrow formatting.
- The output says `claim released`.

Record `HV-15c` as `pass` if the output makes it clear that the task was
completed, the open run was closed, and the developer claim was released.
Record `partial` or `fail` if any of those outcomes is missing or unclear from
CLI output alone.

## 8. Start The Block-Path Task

Start a separate task for the block-path check.

```bash
"$FID" dev start "$HV15_BLOCK_TASK" \
  --model claude-sonnet-4-6 \
  --by developer
```

Expected:

- The command exits zero.
- The output says `claimed by developer`, `transitioned: todo -> in_progress`
  with the CLI's arrow formatting, and `agent run: RUN-... opened`.

This start output can also be used as supporting context for `HV-15d`, but the
human checklist row for `HV-15d` is about the later block command.

## 9. Add A Failed Gate For The Block Path

Record a failed gate so the block output has a concrete failed-gate context.

```bash
"$FID" dev gate "$HV15_BLOCK_TASK" \
  --cmd "golangci-lint run ./..." \
  --exit-code 1 \
  --output-excerpt "human dogfood simulated lint failure" \
  --by developer
```

Expected:

- The command exits zero.
- The output says `gate recorded on HV15-... (attempt 1)`.
- The output shows `exit_code: 1`.
- The output shows `gates_passed: 0, gates_failed: 1`.

## 10. HV-15d - Block The Developer Composite Run

Run `fid dev block` on the block-path task using the disposable blocker reason
created in section 2. Inspect the stderr confirmation as a human operator.

```bash
"$FID" dev block "$HV15_BLOCK_TASK" \
  --reason-md reports/dev-block.md \
  --by developer
```

Expected for `HV-15d`:

- The command exits zero.
- The output starts with `dev block HV15-...`.
- The output says `transitioned: in_progress -> failed_gate` with the CLI's
  arrow formatting.
- The output shows `completion: TC-... (attempt 1, status: failed_gate)`.
- The output says `claim released`.

Record `HV-15d` as `pass` if the output makes it clear that a failed-gate
completion/reason was recorded and the developer claim was released. Record
`partial` or `fail` if the output hides the failed-gate status, reason linkage,
or claim-release outcome.

## 11. Fill The Human Checklist

Return to the repository root and record the outcomes in the checklist.

```bash
cd /home/mark/projects/fid
sed -n '1,220p' project-docs/docs/phases/phase-15/qa/human-checklist.md
```

Fill:

- `HV-15a` from section 4.
- `HV-15b` from section 5.
- `HV-15c` from section 7.
- `HV-15d` from section 10.
- Summary totals and overall result.
