# Phase 19 Human Checklist Run Sheet

> Phase: 19 - Workflow Authoring CLI + Escalation Default Workflow
> Purpose: executable command sheet for
> `project-docs/docs/phases/phase-19/qa/human-checklist.md`
> Date prepared: 2026-08-07

Use this from the fid repository root. It builds the current CLI once, then
uses a new temporary, Git-backed workspace. Every authoring command below runs
only in that disposable workspace. The sheet captures the observable output
needed for the nine Human judgments; it does not substitute deterministic
evidence or product judgments already recorded in `final-review.md`.

This sheet supports `HUMAN-26a` through `HUMAN-26i`, all of which cover
Human-only `VAL-026`. Record the outcome and concise judgment for each row in
`human-checklist.md`. `VAL-026` passes only when all nine rows pass.

## Prerequisites

- Run section 1 from the fid repository root in one shell session. Later
  sections rely on its exported paths and IDs.
- Use only `$FID_HV19_WORKSPACE` after section 1. Never point this sheet at
  the repository's `.fid/` workspace or a real fid workspace.
- The commands require Go, Git, `mktemp`, `mkdir`, `printf`, and `sed`.
- For `HUMAN-26a` and `HUMAN-26b`, do not read Go source, tests, or phase task
  specifications before recording the judgment. The shipped command help and
  the shipped CLI specification excerpt provided below are the intended
  operator guidance.
- Preserve the command output in the terminal or copy the relevant fragments
  into the checklist Notes cells, especially the two correction errors and the
  DEC-079 boundary error.

## 1. Build A Disposable Workspace

Build the current CLI, create a scratch workspace, initialize it, and expose
the shipped command guidance. Do not alter the generated workspace outside
the commands in this sheet.

```bash
set -euo pipefail

export FID_HV19_REPO=/home/mark/projects/fid
export FID_HV19_ROOT="$(mktemp -d /tmp/fid-phase19-human-XXXXXX)"
export FID_HV19_BIN="$FID_HV19_ROOT/fid"
export FID_HV19_WORKSPACE="$FID_HV19_ROOT/workspace"

go build -o "$FID_HV19_BIN" ./cmd/fid
mkdir -p "$FID_HV19_WORKSPACE"
cd "$FID_HV19_WORKSPACE"
git init -q
"$FID_HV19_BIN" init --prefix HV19

printf 'root=%s\nbinary=%s\nworkspace=%s\n' \
  "$FID_HV19_ROOT" "$FID_HV19_BIN" "$FID_HV19_WORKSPACE"
"$FID_HV19_BIN" workflow define --help
sed -n '2133,2223p' "$FID_HV19_REPO/project-docs/specs/06_DESIGN_PLAN_CLI.md"
```

Expected:

- The build and initialization exit zero.
- The help text and the CLI-spec excerpt describe `workflow define`, including
  required flags, state/category ordering, transitions, and `--require`.
- The Human makes the `HUMAN-26a` judgment now, based only on that guidance.

## 2. Define And Assign A Valid Workflow

Define a valid task workflow using v0.1 state names, then inspect assignment
guidance and bind the workflow to kind `task`.

```bash
"$FID_HV19_BIN" workflow define \
  --name hv19-strict \
  --applies-to task \
  --states todo,in_progress,dev_complete,done,cancelled \
  --categories unstarted,started,started,completed,cancelled \
  --transitions 'todo->in_progress,in_progress->dev_complete,dev_complete->done' \
  --require 'in_progress->dev_complete:concerns_addressed' \
  --by tl

"$FID_HV19_BIN" workflow assign --help
"$FID_HV19_BIN" workflow assign --name hv19-strict --kind task --by tl
```

Expected:

- Definition reports `defined workflow hv19-strict for task` and reports its
  state, transition, and requirement counts.
- Assignment reports `assigned workflow hv19-strict to kind task`.
- The Human records `HUMAN-26b` after judging whether the assignment help makes
  kind assignment discoverable and understandable without source-reading.

## 3. Inspect The Four List Formats

Run each documented read format separately. Judge the relevant output instead
of comparing formatting mechanically; the deterministic cross-format checks
are already recorded in the final review.

```bash
export FID_HV19_LIST_JSON="$("$FID_HV19_BIN" workflow list --format json)"
export FID_HV19_LIST_TOON="$("$FID_HV19_BIN" workflow list --format toon)"
export FID_HV19_LIST_MARKDOWN="$("$FID_HV19_BIN" workflow list --format markdown)"
export FID_HV19_LIST_YAML="$("$FID_HV19_BIN" workflow list --format yaml)"

printf '%s\n' "$FID_HV19_LIST_JSON"
printf '%s\n' "$FID_HV19_LIST_TOON"
printf '%s\n' "$FID_HV19_LIST_MARKDOWN"
printf '%s\n' "$FID_HV19_LIST_YAML"
```

Expected:

- Each output includes `hv19-strict` and its `task` applicability.
- The JSON result includes a kind assignment for `task`.
- The Human records `HUMAN-26c` through `HUMAN-26f` independently: JSON for
  operator inspection, TOON for agent-oriented inspection, Markdown for
  human scanability, and YAML for configuration-oriented inspection.

## 4. Judge Definition And Assignment Error Quality

Trigger one invalid definition and one invalid assignment. Each command is
intentionally non-zero; the shell captures the exit status while keeping this
runsheet moving. Copy each exact error into its requested Notes cell.

```bash
set +e
FID_HV19_BAD_DEFINE_OUTPUT=$("$FID_HV19_BIN" workflow define \
  --name 'Invalid Workflow' \
  --applies-to task \
  --states todo,done,cancelled \
  --categories unstarted,completed,cancelled \
  --transitions 'todo->done' \
  --by tl 2>&1)
FID_HV19_BAD_DEFINE_EXIT=$?

FID_HV19_BAD_ASSIGN_OUTPUT=$("$FID_HV19_BIN" workflow assign \
  --name hv19-strict \
  --by tl 2>&1)
FID_HV19_BAD_ASSIGN_EXIT=$?
set -e

printf 'invalid_definition_exit=%s\n%s\n' \
  "$FID_HV19_BAD_DEFINE_EXIT" "$FID_HV19_BAD_DEFINE_OUTPUT"
printf 'invalid_assignment_exit=%s\n%s\n' \
  "$FID_HV19_BAD_ASSIGN_EXIT" "$FID_HV19_BAD_ASSIGN_OUTPUT"
```

Expected:

- Both commands exit non-zero.
- The definition error identifies the offending workflow name.
- The assignment error identifies the missing assignment scope.
- The Human records `HUMAN-26g` and `HUMAN-26h` independently, including the
  exact error and whether it supports self-correction.

## 5. Exercise The DEC-079 Custom-State Boundary

Author and assign a second task workflow with the custom `backlog` state,
create an issue under it, then try to enter that custom state. The phase ships
this limitation intentionally: Phase 20 will make it configurable through
`fid policy set`.

```bash
"$FID_HV19_BIN" workflow define \
  --name hv19-custom-state \
  --applies-to task \
  --states todo,backlog,done,cancelled \
  --categories unstarted,started,completed,cancelled \
  --transitions 'todo->backlog,backlog->done' \
  --by tl
"$FID_HV19_BIN" workflow assign --name hv19-custom-state --kind task --by tl

FID_HV19_CREATE_OUTPUT=$("$FID_HV19_BIN" issue create \
  --kind task \
  --title 'HV19 custom-state boundary' \
  --by developer 2>&1)
printf '%s\n' "$FID_HV19_CREATE_OUTPUT"
export FID_HV19_CUSTOM_ISSUE="$(printf '%s\n' "$FID_HV19_CREATE_OUTPUT" |
  sed -n 's/^created issue \([^ ]*\) .*/\1/p')"
printf 'custom_issue=%s\n' "$FID_HV19_CUSTOM_ISSUE"

set +e
FID_HV19_CUSTOM_STATE_OUTPUT=$("$FID_HV19_BIN" issue transition \
  "$FID_HV19_CUSTOM_ISSUE" \
  --to backlog \
  --by developer 2>&1)
FID_HV19_CUSTOM_STATE_EXIT=$?
set -e

printf 'custom_state_transition_exit=%s\n%s\n' \
  "$FID_HV19_CUSTOM_STATE_EXIT" "$FID_HV19_CUSTOM_STATE_OUTPUT"
```

Expected:

- The second definition and assignment exit zero; the created issue ID is
  non-empty and begins with `HV19-`.
- The transition exits non-zero and names the missing
  `issue.transition.backlog` permissions-matrix entry.
- The Human records `HUMAN-26i`, copying the exact error and judging whether
  the error plus the DEC-079 note above makes the temporary boundary
  understandable rather than surprising.

## 6. Record The Checklist

Update
`project-docs/docs/phases/phase-19/qa/human-checklist.md` directly after the
run. Enter `pass`, `fail`, or `partial` for every `HUMAN-26a` through
`HUMAN-26i`, include the requested error text, recalculate the summary, and
set the overall result according to the checklist's stated rules. A blank row
keeps `VAL-026` blocked.

The disposable material remains at `$FID_HV19_ROOT` for inspection during the
review. Remove it only after the Human has finished recording the outcomes.
