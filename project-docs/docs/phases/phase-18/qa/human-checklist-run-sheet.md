# Phase 18 Human Checklist Run Sheet

> Phase: 18 - Workflow Tables Migration + Engine Data-Source Swap + Gate Seams
> Purpose: executable command sheet for
> `project-docs/docs/phases/phase-18/qa/human-checklist.md`
> Date prepared: 2026-07-25

Use this from the fid repository root. It builds the current CLI and a pinned
pre-Phase-18 v0.1 CLI, creates a small v0.1 workspace, and copies that workspace
before allowing the current CLI to migrate it. All migration, transition,
maintenance, and rebuild checks run against the disposable copy. The pristine
v0.1 source remains available as the backup and before-migration reference.

This sheet supports `HUMAN-01a` through `HUMAN-01d`. These checks cover human
product judgment and have no `VAL-NNN`; the Phase 18 final review already
records all deterministic assertions as passing. The commands produce useful
evidence, but the Human must judge behavioral continuity, clarity, and Phase 18
acceptance, then record the outcomes in `human-checklist.md`.

## Prerequisites

- Run section 1 from the fid repository root.
- Use one shell session so the exported paths, IDs, and captured output remain
  available.
- Run v0.1 setup commands only in `$FID_HV18_V01_WORKSPACE`.
- Run current-CLI migration and dogfood commands only in
  `$FID_HV18_MIGRATED_WORKSPACE`.
- Never point this sheet at the repository's own `.fid/` workspace or the sole
  copy of a real fid workspace.
- The commands require Go, Git, and the POSIX tools already used by project
  documentation: `cp`, `diff`, `mkdir`, `mktemp`, `mv`, `sed`, `sha256sum`,
  `sleep`, `tar`, and `printf`.
- The pinned v0.1 source commit is
  `3a9466fb011020d27379e357ca5bee4c40afadc1`, the parent of the first Phase 18
  implementation commit. It contains migrations 1 through 6 and no workflow
  tables.
- The fixture sets the maintenance grace to zero and waits one second before
  audit reads so the started-without-claim warning is deterministic. This
  changes only disposable fixture configuration.

## 1. Build The Current And v0.1 CLIs

Create one temporary root, export the pre-Phase-18 source into it, and build
both binaries.

```bash
set -euo pipefail

export FID_HV18_REPO=/home/mark/projects/fid
export FID_HV18_V01_COMMIT=3a9466fb011020d27379e357ca5bee4c40afadc1
export FID_HV18_ROOT="$(mktemp -d /tmp/fid-phase18-human-XXXXXX)"
export FID_HV18_V01_SOURCE="$FID_HV18_ROOT/v01-source"
export FID_HV18_V01="$FID_HV18_ROOT/fid-v01"
export FID_HV18_CURRENT="$FID_HV18_ROOT/fid-current"

mkdir -p "$FID_HV18_V01_SOURCE"
git archive "$FID_HV18_V01_COMMIT" |
  tar -x -C "$FID_HV18_V01_SOURCE"

(
  cd "$FID_HV18_V01_SOURCE"
  go build -buildvcs=false -o "$FID_HV18_V01" ./cmd/fid
)

go build -o "$FID_HV18_CURRENT" ./cmd/fid

printf 'root=%s\nv01_cli=%s\ncurrent_cli=%s\n' \
  "$FID_HV18_ROOT" "$FID_HV18_V01" "$FID_HV18_CURRENT"
ls -l "$FID_HV18_V01" "$FID_HV18_CURRENT"
```

Expected:

- Both builds exit zero.
- Both binaries exist under the disposable root.
- `-buildvcs=false` is required only because the archived v0.1 source is
  intentionally outside a Git worktree.

## 2. Create And Seed A v0.1 Source Workspace

Create a Git-backed v0.1 workspace using only the pinned old binary. The five
tasks give the transition and ready checks separate entities so one check
cannot invalidate another.

```bash
export FID_HV18_V01_WORKSPACE="$FID_HV18_ROOT/v01-workspace"
export FID_HV18_MIGRATED_WORKSPACE="$FID_HV18_ROOT/migrated-workspace"

mkdir -p "$FID_HV18_V01_WORKSPACE"
cd "$FID_HV18_V01_WORKSPACE"
git init -q
"$FID_HV18_V01" init --prefix HV18

sed \
  's/started_without_claim_grace_seconds = 300/started_without_claim_grace_seconds = 0/' \
  .fid/config.toml > .fid/config.toml.tmp
mv .fid/config.toml.tmp .fid/config.toml

EPIC_OUTPUT=$("$FID_HV18_V01" issue create \
  --kind epic \
  --title "HV18 migration dogfood" \
  --scheduling prioritised \
  --by tl 2>&1)
printf '%s\n' "$EPIC_OUTPUT"
export HV18_EPIC="$(printf '%s\n' "$EPIC_OUTPUT" |
  sed -n 's/^created issue \([^ ]*\) .*/\1/p')"

OLD_ALLOWED_OUTPUT=$("$FID_HV18_V01" issue create \
  --kind task \
  --title "HV18 pre-migration allowed transition" \
  --epic "$HV18_EPIC" \
  --priority 4 \
  --by developer 2>&1)
printf '%s\n' "$OLD_ALLOWED_OUTPUT"
export HV18_OLD_ALLOWED="$(printf '%s\n' "$OLD_ALLOWED_OUTPUT" |
  sed -n 's/^created issue \([^ ]*\) .*/\1/p')"

CURRENT_ALLOWED_OUTPUT=$("$FID_HV18_V01" issue create \
  --kind task \
  --title "HV18 post-migration allowed transition" \
  --epic "$HV18_EPIC" \
  --priority 4 \
  --by developer 2>&1)
printf '%s\n' "$CURRENT_ALLOWED_OUTPUT"
export HV18_CURRENT_ALLOWED="$(printf '%s\n' "$CURRENT_ALLOWED_OUTPUT" |
  sed -n 's/^created issue \([^ ]*\) .*/\1/p')"

REJECTED_OUTPUT=$("$FID_HV18_V01" issue create \
  --kind task \
  --title "HV18 rejected transition parity" \
  --epic "$HV18_EPIC" \
  --priority 5 \
  --by developer 2>&1)
printf '%s\n' "$REJECTED_OUTPUT"
export HV18_REJECTED="$(printf '%s\n' "$REJECTED_OUTPUT" |
  sed -n 's/^created issue \([^ ]*\) .*/\1/p')"

READY_OUTPUT=$("$FID_HV18_V01" issue create \
  --kind task \
  --title "HV18 inspect migrated presentation" \
  --epic "$HV18_EPIC" \
  --priority 1 \
  --by developer 2>&1)
printf '%s\n' "$READY_OUTPUT"
export HV18_READY="$(printf '%s\n' "$READY_OUTPUT" |
  sed -n 's/^created issue \([^ ]*\) .*/\1/p')"

SECOND_OUTPUT=$("$FID_HV18_V01" issue create \
  --kind task \
  --title "HV18 lower-ranked ready task" \
  --epic "$HV18_EPIC" \
  --priority 3 \
  --by developer 2>&1)
printf '%s\n' "$SECOND_OUTPUT"
export HV18_SECOND="$(printf '%s\n' "$SECOND_OUTPUT" |
  sed -n 's/^created issue \([^ ]*\) .*/\1/p')"

printf 'epic=%s\nold_allowed=%s\ncurrent_allowed=%s\nrejected=%s\nready=%s\nsecond=%s\n' \
  "$HV18_EPIC" "$HV18_OLD_ALLOWED" "$HV18_CURRENT_ALLOWED" \
  "$HV18_REJECTED" "$HV18_READY" "$HV18_SECOND"
```

Expected:

- All six variables print non-empty IDs beginning with `HV18-`.
- `$HV18_READY` has priority 1 and `$HV18_SECOND` has priority 3.
- The source workspace is still a migration-6 v0.1 workspace.

## 3. Capture The v0.1 Behavior And Presentation

Exercise one allowed transition in v0.1, then release its claim to create an
actionable maintenance warning. Also capture a rejected transition without
changing its task.

```bash
"$FID_HV18_V01" issue claim "$HV18_OLD_ALLOWED" \
  --by developer \
  --ttl 1h

export HV18_OLD_ALLOWED_OUTPUT="$("$FID_HV18_V01" issue transition \
  "$HV18_OLD_ALLOWED" \
  --to in_progress \
  --by developer 2>&1)"
printf '%s\n' "$HV18_OLD_ALLOWED_OUTPUT"

"$FID_HV18_V01" issue unclaim "$HV18_OLD_ALLOWED" --by developer

set +e
HV18_OLD_REJECTED_OUTPUT=$("$FID_HV18_V01" issue transition \
  "$HV18_REJECTED" \
  --to dev_complete \
  --by developer 2>&1)
HV18_OLD_REJECTED_EXIT=$?
set -e
export HV18_OLD_REJECTED_OUTPUT HV18_OLD_REJECTED_EXIT
printf 'v01_rejected_exit=%s\n%s\n' \
  "$HV18_OLD_REJECTED_EXIT" "$HV18_OLD_REJECTED_OUTPUT"

sleep 1

export HV18_OLD_GET="$("$FID_HV18_V01" issue get \
  "$HV18_READY" --format markdown)"
export HV18_OLD_LIST="$("$FID_HV18_V01" issue list \
  --epic "$HV18_EPIC" --order priority --format markdown)"
export HV18_OLD_READY="$("$FID_HV18_V01" ready \
  --epic "$HV18_EPIC" --format markdown)"
export HV18_OLD_MAINTENANCE="$("$FID_HV18_V01" maintenance check \
  --format markdown)"
export HV18_OLD_SCHEMA="$("$FID_HV18_V01" sql \
  "SELECT version, dirty FROM schema_migrations" --format markdown)"
export HV18_OLD_WORKFLOW_TABLES="$("$FID_HV18_V01" sql \
  "SELECT name FROM sqlite_schema WHERE type = 'table' AND name LIKE 'workflow%' ORDER BY name" \
  --format markdown)"

printf '%s\n' "$HV18_OLD_SCHEMA"
printf '%s\n' "$HV18_OLD_WORKFLOW_TABLES"
printf '%s\n' "$HV18_OLD_MAINTENANCE"
```

Expected:

- The allowed transition reports `todo` to `in_progress`.
- The rejected transition exits non-zero and says that `todo` to
  `dev_complete` is not allowed.
- The schema query reports version `6` with `dirty` equal to `0`.
- The workflow-table query has no data rows.
- Maintenance reports `warn` and names `$HV18_OLD_ALLOWED` in a
  `started_without_claim` warning.

## 4. Copy The v0.1 Workspace Before Migration

Copy the entire workspace, including canonical files, journal, derived
projection, and Git metadata. Hash the journals on both sides before opening
the copy with the current binary.

```bash
cd "$FID_HV18_ROOT"
cp -a "$FID_HV18_V01_WORKSPACE" "$FID_HV18_MIGRATED_WORKSPACE"

export HV18_SOURCE_JOURNAL="$(
  cd "$FID_HV18_V01_WORKSPACE"
  sha256sum .fid/journal/*.jsonl
)"
export HV18_COPY_JOURNAL_BEFORE="$(
  cd "$FID_HV18_MIGRATED_WORKSPACE"
  sha256sum .fid/journal/*.jsonl
)"

printf 'source_journal:\n%s\n' "$HV18_SOURCE_JOURNAL"
printf 'copy_journal_before:\n%s\n' "$HV18_COPY_JOURNAL_BEFORE"

if [ "$HV18_SOURCE_JOURNAL" = "$HV18_COPY_JOURNAL_BEFORE" ]; then
  printf 'source_copy_journal_parity=OK\n'
else
  printf 'source_copy_journal_parity=DIFF\n'
fi
```

Expected: `source_copy_journal_parity=OK`.

Do not run the current binary in `$FID_HV18_V01_WORKSPACE`. That directory is
the pristine backup and v0.1 reference for the remainder of the run.

## 5. Migrate The Copy And Compare Read Surfaces

The first current-CLI read migrates and refreshes the copied projection. Capture
the same four read surfaces, confirm the new tables, and check that migration
did not rewrite the journal.

```bash
cd "$FID_HV18_MIGRATED_WORKSPACE"

export HV18_NEW_GET="$("$FID_HV18_CURRENT" issue get \
  "$HV18_READY" --format markdown)"
export HV18_NEW_LIST="$("$FID_HV18_CURRENT" issue list \
  --epic "$HV18_EPIC" --order priority --format markdown)"
export HV18_NEW_READY="$("$FID_HV18_CURRENT" ready \
  --epic "$HV18_EPIC" --format markdown)"
export HV18_NEW_MAINTENANCE="$("$FID_HV18_CURRENT" maintenance check \
  --format markdown)"
export HV18_NEW_SCHEMA="$("$FID_HV18_CURRENT" sql \
  "SELECT version, dirty FROM schema_migrations" --format markdown)"
export HV18_NEW_WORKFLOW_TABLES="$("$FID_HV18_CURRENT" sql \
  "SELECT name FROM sqlite_schema WHERE type = 'table' AND name LIKE 'workflow%' ORDER BY name" \
  --format markdown)"
export HV18_COPY_JOURNAL_AFTER="$(sha256sum .fid/journal/*.jsonl)"

printf '%s\n' "$HV18_NEW_SCHEMA"
printf '%s\n' "$HV18_NEW_WORKFLOW_TABLES"

if [ "$HV18_OLD_GET" = "$HV18_NEW_GET" ]; then
  printf 'issue_get_parity=OK\n'
else
  printf 'issue_get_parity=DIFF\n'
fi

if [ "$HV18_OLD_LIST" = "$HV18_NEW_LIST" ]; then
  printf 'issue_list_parity=OK\n'
else
  printf 'issue_list_parity=DIFF\n'
fi

if [ "$HV18_OLD_READY" = "$HV18_NEW_READY" ]; then
  printf 'ready_parity=OK\n'
else
  printf 'ready_parity=DIFF\n'
fi

if [ "$HV18_OLD_MAINTENANCE" = "$HV18_NEW_MAINTENANCE" ]; then
  printf 'maintenance_parity=OK\n'
else
  printf 'maintenance_parity=DIFF\n'
fi

if [ "$HV18_COPY_JOURNAL_BEFORE" = "$HV18_COPY_JOURNAL_AFTER" ]; then
  printf 'migration_journal_parity=OK\n'
else
  printf 'migration_journal_parity=DIFF\n'
fi
```

Expected:

- The migrated schema reports version `8` with `dirty` equal to `0`.
- The workflow-table query lists `workflows`, `workflow_states`,
  `workflow_transitions`, and `workflow_transition_requires`.
- All four read-surface parity checks print `OK`.
- `migration_journal_parity=OK` prints.

## 6. HUMAN-01a - Judge Allowed And Rejected Transitions

Exercise the same allowed edge on the reserved post-migration task and the same
rejected edge on the unchanged rejection task.

```bash
"$FID_HV18_CURRENT" issue claim "$HV18_CURRENT_ALLOWED" \
  --by developer \
  --ttl 1h

export HV18_NEW_ALLOWED_OUTPUT="$("$FID_HV18_CURRENT" issue transition \
  "$HV18_CURRENT_ALLOWED" \
  --to in_progress \
  --by developer 2>&1)"
printf '%s\n' "$HV18_NEW_ALLOWED_OUTPUT"

"$FID_HV18_CURRENT" issue unclaim "$HV18_CURRENT_ALLOWED" --by developer

set +e
HV18_NEW_REJECTED_OUTPUT=$("$FID_HV18_CURRENT" issue transition \
  "$HV18_REJECTED" \
  --to dev_complete \
  --by developer 2>&1)
HV18_NEW_REJECTED_EXIT=$?
set -e
export HV18_NEW_REJECTED_OUTPUT HV18_NEW_REJECTED_EXIT

printf 'v01_allowed:\n%s\n' "$HV18_OLD_ALLOWED_OUTPUT"
printf 'migrated_allowed:\n%s\n' "$HV18_NEW_ALLOWED_OUTPUT"
printf 'v01_rejected_exit=%s\n%s\n' \
  "$HV18_OLD_REJECTED_EXIT" "$HV18_OLD_REJECTED_OUTPUT"
printf 'migrated_rejected_exit=%s\n%s\n' \
  "$HV18_NEW_REJECTED_EXIT" "$HV18_NEW_REJECTED_OUTPUT"

if [ "$HV18_OLD_REJECTED_OUTPUT" = "$HV18_NEW_REJECTED_OUTPUT" ]; then
  printf 'rejected_transition_parity=OK\n'
else
  printf 'rejected_transition_parity=DIFF\n'
fi
```

Expected:

- Both allowed outputs clearly report `todo` to `in_progress`; their issue IDs
  differ because separate tasks preserve independent fixture states.
- Both rejected commands exit non-zero.
- `rejected_transition_parity=OK` prints.
- Neither rejection sounds like a missing flag, permission error, or internal
  database error.

Record `HUMAN-01a` as `pass` only if the migrated allowed and rejected paths
feel behaviorally unchanged, clear, and unsurprising. Record `partial` or
`fail` if the mechanics pass but the migrated messages create new ambiguity.

## 7. HUMAN-01b - Judge Get, List, And Ready Presentation

Print the migrated presentation captured before the transition checks changed
fixture state.

```bash
printf '%s\n' "$HV18_NEW_GET"
printf '%s\n' "$HV18_NEW_LIST"
printf '%s\n' "$HV18_NEW_READY"
```

Expected:

- `issue get` presents `$HV18_READY` as a `todo` task in the `unstarted`
  category.
- `issue list` is readable and preserves priority ordering within
  `$HV18_EPIC`.
- `ready` selects `$HV18_READY`, titled `HV18 inspect migrated presentation`,
  at priority 1 rather than the priority-3 `$HV18_SECOND`.
- Section 5 printed `issue_get_parity=OK`, `issue_list_parity=OK`, and
  `ready_parity=OK`.

Record `HUMAN-01b` as `pass` only if the categories, ordering, and ready choice
both match v0.1 and remain sensible to the Human. A byte-identical result can
still be `partial` if the presentation is difficult to interpret.

## 8. HUMAN-01c - Judge Maintenance Clarity

Wait for the second deliberately unclaimed started task to pass the zero-second
fixture boundary, then inspect maintenance again.

```bash
sleep 1
"$FID_HV18_CURRENT" maintenance check --format markdown
```

Expected:

- The command exits zero with `**Status:** warn`.
- The warnings include `$HV18_OLD_ALLOWED` and `$HV18_CURRENT_ALLOWED`.
- Each warning is named `started_without_claim`, identifies the affected
  issue, and explains that it has no live claim.
- No workflow-pair or database-structure error appears.
- Section 5 already printed `maintenance_parity=OK` for the unchanged
  pre/post-migration fixture state.

Record `HUMAN-01c` as `pass` if the migrated maintenance output remains
understandable and points to the actionable states without creating uncertainty
about overall workspace health. Record `partial` or `fail` if the warnings are
technically present but difficult to interpret or act on.

## 9. Demonstrate Backup, Rebuild, And Journal Preservation

Keep a copy of the migrated database outside `.fid/`, rebuild the derived
projection from the journal, and compare the observable issue list and journal
hashes.

```bash
export HV18_LIST_BEFORE_REBUILD="$("$FID_HV18_CURRENT" issue list \
  --epic "$HV18_EPIC" --order priority --format markdown)"
export HV18_JOURNAL_BEFORE_REBUILD="$(sha256sum .fid/journal/*.jsonl)"

cp -p .fid/fid.sqlite "$FID_HV18_ROOT/fid.sqlite.before-rebuild"
"$FID_HV18_CURRENT" rebuild --by human

export HV18_LIST_AFTER_REBUILD="$("$FID_HV18_CURRENT" issue list \
  --epic "$HV18_EPIC" --order priority --format markdown)"
export HV18_JOURNAL_AFTER_REBUILD="$(sha256sum .fid/journal/*.jsonl)"

if [ "$HV18_LIST_BEFORE_REBUILD" = "$HV18_LIST_AFTER_REBUILD" ]; then
  printf 'issue_list_rebuild_parity=OK\n'
else
  printf 'issue_list_rebuild_parity=DIFF\n'
fi

if [ "$HV18_JOURNAL_BEFORE_REBUILD" = "$HV18_JOURNAL_AFTER_REBUILD" ]; then
  printf 'journal_rebuild_parity=OK\n'
else
  printf 'journal_rebuild_parity=DIFF\n'
fi

export HV18_SOURCE_JOURNAL_AFTER="$(
  cd "$FID_HV18_V01_WORKSPACE"
  sha256sum .fid/journal/*.jsonl
)"
if [ "$HV18_SOURCE_JOURNAL" = "$HV18_SOURCE_JOURNAL_AFTER" ]; then
  printf 'pristine_source_journal=OK\n'
else
  printf 'pristine_source_journal=DIFF\n'
fi

test -f .fid/fid.sqlite
test -f "$FID_HV18_ROOT/fid.sqlite.before-rebuild"
"$FID_HV18_CURRENT" maintenance check --format markdown
git add --dry-run .fid
```

Expected:

- Rebuild exits zero and reports `rebuild complete`.
- `issue_list_rebuild_parity=OK` and `journal_rebuild_parity=OK` print.
- `pristine_source_journal=OK` confirms that the v0.1 backup was never opened
  by the current binary.
- The rebuilt maintenance output retains the two actionable warnings.
- `git add --dry-run .fid` does not list `.fid/fid.sqlite`.
- Both the rebuilt database and the disposable pre-rebuild database copy exist.

## 10. HUMAN-01d - Record Phase Acceptance

Return to the repository root and inspect the checklist before recording all
four judgments.

```bash
cd "$FID_HV18_REPO"
sed -n '1,220p' \
  project-docs/docs/phases/phase-18/qa/human-checklist.md
```

Fill:

- `HUMAN-01a` from section 6.
- `HUMAN-01b` from section 7.
- `HUMAN-01c` from section 8.
- `HUMAN-01d` from the complete backup/copy, migration, behavior, maintenance,
  and rebuild run.

Record `HUMAN-01d` as `pass` only if the demonstrated evidence is sufficient to
accept Phase 18 for continued dogfood. Include concise rationale and name any
remaining friction that should be carried forward. Use `partial` or `fail` and
reference a follow-up decision, escalation, or carry-forward item when the
evidence is insufficient.

Update the summary totals and overall result. Overall result is `pass` only when
all four rows pass; `partial` when at least one row is partial and none fail;
`fail` when any row fails.

## 11. Cleanup

Do not clean up until outcomes and useful evidence have been recorded.

```bash
printf 'disposable_root=%s\n' "$FID_HV18_ROOT"
rm -rf "$FID_HV18_ROOT"
```

Expected: the v0.1 source, migrated copy, both binaries, archived old source,
and pre-rebuild database copy are removed together. The fid repository and any
real `.fid/` workspace are untouched.
