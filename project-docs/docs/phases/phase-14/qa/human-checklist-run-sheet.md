# Phase 14 Human Checklist Run Sheet

> Phase: 14 - Telemetry + Pricing
> Purpose: executable command sheet for
> `project-docs/docs/phases/phase-14/qa/human-checklist.md`
> Date prepared: 2026-06-28

Use this from the fid repository root. It creates a disposable real `.fid/`
workspace under `/tmp`, builds the current CLI to `/tmp/fid`, creates an epic
plus issues, records known-model and unknown-model agent runs, then runs the
token and cost aggregate commands needed for the human dogfood checks.

The human checklist still needs human judgment. Treat the command output as
dogfood evidence, then record the outcome and notes in `human-checklist.md`.

## Prerequisites

- Run the build command from the repo root.
- Run workspace commands only inside the disposable `/tmp` workspace created
  below, not inside this repository's own `.fid/` workspace.
- These commands require the built `fid` binary and POSIX shell tools already
  used by the project docs (`sed`, `printf`, `ls`, `mkdir`, `mktemp`).
- This sheet does not complete deterministic QA gates. Phase 14 final review
  already covers pricing defaults, `run.record` projection, unknown-model
  persistence, policy parity, stats aggregation, format parity, and phase-level
  gates. This sheet only supports `HUMAN-14a` through `HUMAN-14c` for
  `VAL-009`.

## 1. Build Current CLI

```bash
go build -o /tmp/fid ./cmd/fid
export FID=/tmp/fid
ls -l "$FID"
```

Expected: command exits zero and `/tmp/fid` exists.

## 2. Create A Fresh Test Workspace

```bash
export FID_HV14_WORKSPACE="$(mktemp -d /tmp/fid-phase14-human-XXXXXX)"
cd "$FID_HV14_WORKSPACE"
printf 'workspace=%s\n' "$FID_HV14_WORKSPACE"

"$FID" init --prefix HV14
```

Expected:

- `fid init` reports a new `.fid/` workspace.
- The workspace has a `.fid/` directory and generated `.fid/pricing.toml`.
- The rest of this sheet is run from the disposable workspace.

## 3. Create Issues For Telemetry Runs

Create one epic, one task under that epic, and one standalone task. The epic
lets the stats commands demonstrate both an epic bucket and a `(none)` bucket.

```bash
EPIC_OUTPUT=$("$FID" issue create \
  --kind epic \
  --title "HV14 telemetry epic" \
  --scheduling prioritised \
  --by tl 2>&1)
printf '%s\n' "$EPIC_OUTPUT"
export HV14_EPIC="$(printf '%s\n' "$EPIC_OUTPUT" | sed -n 's/^created issue \([^ ]*\) .*/\1/p')"
printf 'epic=%s\n' "$HV14_EPIC"

TASK_OUTPUT=$("$FID" issue create \
  --kind task \
  --title "HV14 telemetry task" \
  --epic "$HV14_EPIC" \
  --by developer 2>&1)
printf '%s\n' "$TASK_OUTPUT"
export HV14_TASK="$(printf '%s\n' "$TASK_OUTPUT" | sed -n 's/^created issue \([^ ]*\) .*/\1/p')"
printf 'task=%s\n' "$HV14_TASK"

STANDALONE_OUTPUT=$("$FID" issue create \
  --kind task \
  --title "HV14 standalone telemetry task" \
  --by developer 2>&1)
printf '%s\n' "$STANDALONE_OUTPUT"
export HV14_STANDALONE="$(printf '%s\n' "$STANDALONE_OUTPUT" | sed -n 's/^created issue \([^ ]*\) .*/\1/p')"
printf 'standalone=%s\n' "$HV14_STANDALONE"
```

Expected:

- The epic, epic-owned task, and standalone task are created.
- `HV14_EPIC`, `HV14_TASK`, and `HV14_STANDALONE` print non-empty issue IDs.

## 4. HUMAN-14a - Record A Known-Model Run With Computed Cost

Record a developer run on the epic-owned task using a model from the generated
default pricing config. Inspect the stderr confirmation as a human operator.

```bash
"$FID" run record \
  --issue "$HV14_TASK" \
  --role developer \
  --model claude-sonnet-4-6 \
  --tokens-in 1200 \
  --tokens-out 3400 \
  --by developer
```

Expected for `HUMAN-14a`:

- The command exits zero.
- The output says `recorded agent run RUN-... on HV14-...`.
- The output shows `role: developer, model: claude-sonnet-4-6`.
- The output shows `tokens: 1200 in, 3400 out`.
- The output shows a computed `cost: $...`, not `cost: unknown`.

Record `HUMAN-14a` as `pass` if the output makes it clear that an agent run was
recorded and that the run has a computed USD cost. Record `partial` or `fail` if
the cost is absent, ambiguous, or only understandable by reading implementation
files.

## 5. HUMAN-14b - Record An Unknown-Model Run With Warning

Record a QA reviewer run on the same epic-owned task using a deliberately
unknown model. Inspect the warning and confirmation.

```bash
"$FID" run record \
  --issue "$HV14_TASK" \
  --role qa_reviewer \
  --model hv14-unknown-model \
  --tokens-in 500 \
  --tokens-out 700 \
  --by qa_reviewer
```

Expected for `HUMAN-14b`:

- The command exits zero.
- The warning names `hv14-unknown-model`.
- The warning says the model was not found in `.fid/pricing.toml`.
- The command still reports `recorded agent run RUN-... on HV14-...`.
- The output shows `cost: unknown`, without implying a guessed price.

Record `HUMAN-14b` as `pass` if the output clearly communicates that the run was
recorded, the model was unknown, and no price was guessed. Record `partial` or
`fail` if the warning is missing, unclear, or the command appears to reject the
run.

## 6. Add Runs For Aggregate Inspection

Add one known-model run directly on the epic and one known-model run on the
standalone task. These make the role, model, and epic aggregate tables easier
to judge.

```bash
"$FID" run record \
  --issue "$HV14_EPIC" \
  --role tl \
  --model gemini-2.5-flash \
  --tokens-in 700 \
  --tokens-out 70 \
  --by tl

"$FID" run record \
  --issue "$HV14_STANDALONE" \
  --role developer \
  --model gpt-5.4-mini \
  --tokens-in 300 \
  --tokens-out 50 \
  --by developer
```

Expected:

- Both commands exit zero.
- Both confirmations show `recorded agent run RUN-...`.
- Both confirmations show computed `cost: $...`.

## 7. HUMAN-14c - Inspect Token Aggregates

Inspect token totals grouped by role, model, and epic.

```bash
"$FID" stats tokens --by role --format markdown
"$FID" stats tokens --by model --format markdown
"$FID" stats tokens --by epic --format markdown
```

Expected token totals:

- By role: `developer` has `1500` input, `3450` output, `2` runs;
  `qa_reviewer` has `500` input, `700` output, `1` run; `tl` has `700` input,
  `70` output, `1` run.
- By model: `claude-sonnet-4-6`, `gemini-2.5-flash`, `gpt-5.4-mini`, and
  `hv14-unknown-model` each appear with one run.
- By epic: `$HV14_EPIC` has `2400` input, `4170` output, `3` runs; `(none)` has
  `300` input, `50` output, `1` run.

## 8. HUMAN-14c - Inspect Cost Aggregates

Inspect cost totals grouped by role, model, and epic. The unknown-model run
counts in token totals and `run_count`, but contributes `0` to cost totals
because its stored cost is NULL.

```bash
"$FID" stats cost --by role --format markdown
"$FID" stats cost --by model --format markdown
"$FID" stats cost --by epic --format markdown
```

Expected cost-table shape:

- By role: `developer`, `qa_reviewer`, and `tl` each appear with token totals
  and run counts.
- By model: all four models appear, including `hv14-unknown-model`.
- By epic: `$HV14_EPIC` and `(none)` both appear.
- `qa_reviewer` and `hv14-unknown-model` show `0` cost while still showing
  their `500` input tokens, `700` output tokens, and `1` run.
- The exact known-model cost values may change if the default pricing file is
  intentionally updated; the human check is whether costs, tokens, and run
  counts are understandable from CLI output alone.

## 9. HUMAN-14c - Inspect Epic-Filtered Aggregates

Optionally scope the aggregate tables to the epic-created issue subtree.

```bash
"$FID" stats tokens --by role --epic "$HV14_EPIC" --format markdown
"$FID" stats cost --by role --epic "$HV14_EPIC" --format markdown
```

Expected:

- The filtered output excludes the standalone task run.
- The role rows are `developer`, `qa_reviewer`, and `tl`.
- The filtered token totals are `developer` `1200/3400`, `qa_reviewer`
  `500/700`, and `tl` `700/70`.
- The filtered cost table still counts the unknown-model QA reviewer run while
  showing `0` cost for that row.

Record `HUMAN-14c` as `pass` if the stats output makes token totals, cost
totals, and run counts understandable by role, model, and epic. Record
`partial` or `fail` if the grouping labels, cost handling, or run counts are
hard to interpret from CLI output.

## 10. Fill The Human Checklist

Return to the repository root and record the outcomes in the checklist.

```bash
cd /home/mark/projects/fid
sed -n '1,220p' project-docs/docs/phases/phase-14/qa/human-checklist.md
```

Fill:

- `HUMAN-14a` from section 4.
- `HUMAN-14b` from section 5.
- `HUMAN-14c` from sections 7 through 9.
- Summary totals and overall result.

Use `pass`, `partial`, or `fail` for each row. Overall result is `pass` only if
all three rows pass; `partial` if at least one row is partial and none fail;
`fail` if any row fails.

## 11. Cleanup

When the checklist evidence is recorded, remove the disposable workspace if it
is no longer needed.

```bash
printf 'workspace=%s\n' "$FID_HV14_WORKSPACE"
rm -rf "$FID_HV14_WORKSPACE"
```

Expected: the temporary `/tmp/fid-phase14-human-*` workspace is removed. Do not
run the cleanup command until any useful terminal output has been recorded in
the checklist notes.
