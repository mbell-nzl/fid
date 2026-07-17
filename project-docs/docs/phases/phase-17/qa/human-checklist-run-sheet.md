# Phase 17 Human Checklist Run Sheet

> Phase: 17 - Maintenance, Stub Guards, and v0.1 Closure
> Purpose: executable command sheet for
> `project-docs/docs/phases/phase-17/qa/human-checklist.md`
> Date prepared: 2026-07-18

Use this from the fid repository root. It creates a disposable Git-backed
workspace under `/tmp`, builds the current CLI to `/tmp/fid`, and seeds a small
project state for the five human-owned Phase 17 checks.

This sheet supports `HUMAN-25a` through `HUMAN-25e` for `VAL-025`. It does not
replace the deterministic Phase 17 QA gates or the final review. The commands
produce evidence, but the Human must judge clarity, usefulness, and product
acceptance, then record the outcome in `human-checklist.md`.

## Prerequisites

- Run the build command from the fid repository root.
- Run all workspace commands inside the disposable workspace created below.
- Do not run the recovery section against the repository's own `.fid/`
  workspace.
- These commands require Go, Git, and the POSIX tools already used by the
  project docs: `sed`, `mktemp`, `mv`, `sha256sum`, and `printf`.
- The fixture changes the maintenance grace to zero seconds so the warning can
  be observed immediately. This is disposable fixture configuration only.

## 1. Build The Current CLI

```bash
go build -o /tmp/fid ./cmd/fid
export FID=/tmp/fid
ls -l "$FID"
```

Expected: the command exits zero and `/tmp/fid` exists.

## 2. Create A Disposable Test Workspace

```bash
set -euo pipefail

export FID_HV17_WORKSPACE="$(mktemp -d /tmp/fid-phase17-human-XXXXXX)"
cd "$FID_HV17_WORKSPACE"
printf 'workspace=%s\n' "$FID_HV17_WORKSPACE"

git init -q
"$FID" init --prefix HV17
```

Expected:

- `fid init` reports a new `.fid/` workspace with prefix `HV17`.
- The workspace is a Git repository, but no commit is required.
- Keep the shell in this workspace for the rest of the fixture run.

## 3. Make The Maintenance Warning Immediate

The normal grace is 300 seconds. Set it to zero in the disposable config so a
started issue with no live claim is visible to `maintenance check` without
waiting.

```bash
sed 's/started_without_claim_grace_seconds = 300/started_without_claim_grace_seconds = 0/' \
  .fid/config.toml > .fid/config.toml.tmp
mv .fid/config.toml.tmp .fid/config.toml
sed -n '/\[maintenance\]/,$p' .fid/config.toml
```

Expected: the last line is
`started_without_claim_grace_seconds = 0`.

## 4. Seed A Small Project State

Create one started task without a live claim for the maintenance finding. The
claim/transition/unclaim sequence is required because the workflow requires a
claim for the transition, while the final fixture state must have no live
claim.

```bash
STALE_OUTPUT=$("$FID" issue create \
  --kind task \
  --title "HV17 maintenance follow-up" \
  --priority 0 \
  --by developer 2>&1)
printf '%s\n' "$STALE_OUTPUT"
export HV17_STALE_TASK="$(printf '%s\n' "$STALE_OUTPUT" | sed -n 's/^created issue \([^ ]*\) .*/\1/p')"

"$FID" issue claim "$HV17_STALE_TASK" --by developer --ttl 1h
"$FID" issue transition "$HV17_STALE_TASK" --to in_progress --by developer
"$FID" issue unclaim "$HV17_STALE_TASK" --by developer

EPIC_OUTPUT=$("$FID" issue create \
  --kind epic \
  --title "HV17 v0.1 dogfood work" \
  --scheduling prioritised \
  --by tl 2>&1)
printf '%s\n' "$EPIC_OUTPUT"
export HV17_EPIC="$(printf '%s\n' "$EPIC_OUTPUT" | sed -n 's/^created issue \([^ ]*\) .*/\1/p')"

READY_OUTPUT=$("$FID" issue create \
  --kind task \
  --title "Review maintenance findings" \
  --epic "$HV17_EPIC" \
  --priority 1 \
  --by developer 2>&1)
printf '%s\n' "$READY_OUTPUT"
export HV17_READY_TASK="$(printf '%s\n' "$READY_OUTPUT" | sed -n 's/^created issue \([^ ]*\) .*/\1/p')"

SECOND_OUTPUT=$("$FID" issue create \
  --kind task \
  --title "Document recovery procedure" \
  --epic "$HV17_EPIC" \
  --priority 3 \
  --by developer 2>&1)
printf '%s\n' "$SECOND_OUTPUT"
export HV17_SECOND_TASK="$(printf '%s\n' "$SECOND_OUTPUT" | sed -n 's/^created issue \([^ ]*\) .*/\1/p')"

printf 'stale_task=%s\nepic=%s\nready_task=%s\nsecond_task=%s\n' \
  "$HV17_STALE_TASK" "$HV17_EPIC" "$HV17_READY_TASK" "$HV17_SECOND_TASK"
```

Expected:

- All four variables print non-empty IDs beginning with `HV17-`.
- `HV17_STALE_TASK` is `in_progress` with no live claim.
- `HV17_READY_TASK` has priority 1 and `HV17_SECOND_TASK` has priority 3.

## 5. HUMAN-25a - Inspect Maintenance Findings

Run the documented human-facing audit output and judge whether it explains the
fixture's actionable problem.

```bash
"$FID" maintenance check --format markdown
```

Expected:

- The command exits zero with `**Status:** warn`.
- The output contains a `started_without_claim` warning for
  `$HV17_STALE_TASK`.
- The warning names the affected issue and says that it has no live claim.
- The Markdown is readable without knowing the internal table or package names.

Record `HUMAN-25a` as `pass` only if the finding makes the affected entity and
the next useful action clear from CLI output alone. Record `partial` or `fail`
if the warning is technically present but hard to interpret or act on.

Optional promotion check:

```bash
set +e
"$FID" maintenance check --warn-as-error --format markdown
HV17_MAINTENANCE_PROMOTED_EXIT=$?
set -e
printf 'warn_as_error_exit=%s\n' "$HV17_MAINTENANCE_PROMOTED_EXIT"
```

Expected: the output still contains the same warning, reports `**Status:**
fail`, and the captured exit code is non-zero. This is supporting context for
the human judgment, not a separate checklist row.

## 6. HUMAN-25b - Check A v0.2 Governance Guard

Pass a future-looking flag deliberately. The guard should win before flag,
actor, workspace, database, or journal validation.

```bash
set +e
HV17_GUARD_OUTPUT=$("$FID" decision create --future-flag value 2>&1)
HV17_GUARD_EXIT=$?
set -e

printf 'guard_exit=%s\n%s\n' "$HV17_GUARD_EXIT" "$HV17_GUARD_OUTPUT"
```

Expected:

- `guard_exit` is non-zero.
- The first error line is exactly:
  `error: commands for decision ship in v0.2`.
- No Cobra unknown-command or unknown-flag message appears first.

Record `HUMAN-25b` as `pass` if the error immediately explains both what is
unavailable and when it ships. Record `partial` or `fail` if the version gate
is ambiguous, buried, or sounds like an unrelated argument error.

## 7. HUMAN-25c - Judge The Ready Choice

Inspect the ready task selected from the fixture's prioritised epic.

```bash
"$FID" ready --epic "$HV17_EPIC" --format markdown
```

Expected:

- The command exits zero and renders `# Ready Task`.
- The selected task is `$HV17_READY_TASK`, titled `Review maintenance
  findings`, with priority `P1`.
- The priority-3 `Document recovery procedure` task is not selected while the
  P1 task remains ready.

Record `HUMAN-25c` from the Human's product judgment, not only from the
priority comparison: decide whether reviewing the maintenance finding is the
sensible next piece of work for this small project state. A mechanically
correct selection can still receive `partial` if the displayed context is not
enough to make the choice understandable.

## 8. HUMAN-25d - Exercise The Disposable Recovery Path

First inspect the current state and compact the derived database as an
authorized Human. Then move the disposable database outside `.fid/`, rebuild it
from the journal, and inspect the restored state.

```bash
export HV17_ISSUES_BEFORE="$("$FID" issue list --format markdown)"
export HV17_JOURNAL_BEFORE="$(sha256sum .fid/journal/*.jsonl)"
printf '%s\n' "$HV17_ISSUES_BEFORE"
printf 'journal_before=%s\n' "$HV17_JOURNAL_BEFORE"

"$FID" maintenance vacuum --by human

mv .fid/fid.sqlite "$FID_HV17_WORKSPACE/fid.sqlite.before-rebuild"
"$FID" rebuild --by human

export HV17_ISSUES_AFTER="$("$FID" issue list --format markdown)"
export HV17_JOURNAL_AFTER="$(sha256sum .fid/journal/*.jsonl)"
printf '%s\n' "$HV17_ISSUES_AFTER"
printf 'journal_after=%s\n' "$HV17_JOURNAL_AFTER"

if [ "$HV17_ISSUES_BEFORE" = "$HV17_ISSUES_AFTER" ]; then
  printf 'issue_list_rebuild_parity=OK\n'
else
  printf 'issue_list_rebuild_parity=DIFF\n'
fi

if [ "$HV17_JOURNAL_BEFORE" = "$HV17_JOURNAL_AFTER" ]; then
  printf 'journal_hash_parity=OK\n'
else
  printf 'journal_hash_parity=DIFF\n'
fi

test -f .fid/fid.sqlite
"$FID" maintenance check --format markdown
git add --dry-run .fid
```

Expected:

- VACUUM reports a completion line with before and after database byte counts.
- Rebuild exits zero and reports `rebuild complete`.
- `issue_list_rebuild_parity=OK` and `journal_hash_parity=OK` print.
- The rebuilt maintenance check still shows the same warning for
  `$HV17_STALE_TASK`.
- `git add --dry-run .fid` lists the journal and canonical configuration files,
  but does not list `.fid/fid.sqlite`.

Record `HUMAN-25d` as `pass` only if a Human can follow the sequence, identify
what is safe to move or rebuild, and confirm that the journal-backed state was
restored without risking the real workspace. Record `partial` or `fail` for
unclear recovery guidance, surprising output, or any temptation to operate on
the wrong workspace.

## 9. HUMAN-25e - Record The v0.1 Dogfood Decision

Return to the repository root and record the first four judgments before making
the phase-end product decision.

```bash
cd /home/mark/projects/fid
sed -n '1,220p' project-docs/docs/phases/phase-17/qa/human-checklist.md
```

Fill `HUMAN-25a` through `HUMAN-25d` with `pass`, `partial`, or `fail` and
include concise notes tied to the observed output. Then fill `HUMAN-25e`:

- Record whether v0.1 is accepted for continued early dogfood.
- If accepted, name any friction that should still be carried forward.
- If not accepted, reference the follow-up issue, decision, escalation, or
  carry-forward item that disposes of the remaining friction.
- Update the summary totals and overall result. Overall result is `pass` only
  when all five rows pass; `partial` when at least one row is partial and none
  fail; `fail` when any row fails.

## 10. Cleanup

Do not clean up until the notes and outcomes have been recorded.

```bash
printf 'workspace=%s\n' "$FID_HV17_WORKSPACE"
rm -rf "$FID_HV17_WORKSPACE"
```

Expected: the disposable `/tmp/fid-phase17-human-*` workspace and its
temporary pre-rebuild database are removed. No repository `.fid/` data is
touched.
