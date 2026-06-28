# Phase 13 Human Checklist Run Sheet

> Phase: 13 - Labels + v0.1d Closure
> Purpose: executable command sheet for
> `project-docs/docs/phases/phase-13/qa/human-checklist.md`
> Date prepared: 2026-06-28

Use this from the fid repository root. It creates a disposable real `.fid/`
workspace under `/tmp`, builds the current CLI to `/tmp/fid`, creates labels and
issues, then runs the label and lifecycle commands needed for the human
dogfood checks.

The human checklist still needs human judgment. Treat the command output as
dogfood evidence, then record the outcome and notes in `human-checklist.md`.

## Prerequisites

- Run the build command from the repo root.
- Run workspace commands only inside the disposable `/tmp` workspace created
  below, not inside this repository's own `.fid/` workspace.
- These commands require the built `fid` binary and POSIX shell tools already
  used by the project docs (`grep`, `sed`, `printf`, `ls`, `mkdir`, `mktemp`).
- This sheet does not complete deterministic QA gates. Phase 13 final review
  already covers matrix parity, validate-before-append behavior, activity-log
  projection, rebuild parity, and phase-level gates. This sheet only supports
  `HV-13a` through `HV-13d` for `VAL-010`.

## 1. Build Current CLI

```bash
go build -o /tmp/fid ./cmd/fid
export FID=/tmp/fid
ls -l "$FID"
```

Expected: command exits zero and `/tmp/fid` exists.

## 2. Create A Fresh Test Workspace

```bash
export FID_HV13_WORKSPACE="$(mktemp -d /tmp/fid-phase13-human-XXXXXX)"
cd "$FID_HV13_WORKSPACE"
printf 'workspace=%s\n' "$FID_HV13_WORKSPACE"

"$FID" init --prefix HV13
mkdir -p notes
printf '# HV13 completion notes\n\nHuman checklist lifecycle smoke.\n' > notes/hv13-completion.md
```

Expected:

- `fid init` reports a new `.fid/` workspace.
- The workspace has a `.fid/` directory.
- The rest of this sheet is run from the disposable workspace.

## 3. HUMAN-13a - Create And List A Label

Create labels used by the later checks. Inspect the list output as a human
operator.

```bash
"$FID" label add core --colour '#2f6fed' --description 'Core label for human label check' --by developer
"$FID" label add legacy --colour '#6b7280' --description 'Existing label that should stay on the issue' --by developer
"$FID" label add api --colour '#10b981' --description 'Added label for additive check' --by developer

"$FID" label list --format markdown
```

Optionally narrow the terminal output to the label rows:

```bash
"$FID" label list --format markdown | grep -E 'core|legacy|api|Name|Colour|Description'
```

Expected for `HV-13a`:

- The output names the created labels.
- The `core` row shows name, colour, and description.
- The labels are visible enough for a user to identify them without reading
  implementation files.

Record `HV-13a` as `pass` if the label list makes the created label and its
metadata easy enough to identify. Record `partial` or `fail` if the metadata is
missing, hard to locate, or ambiguous.

## 4. HUMAN-13b - Add One Label Without Replacing Existing Labels

Create an issue with two existing labels, then add exactly one more label.
Inspect the before/after issue output.

```bash
CREATE_OUTPUT=$("$FID" issue create \
  --kind task \
  --title "HV13 additive label issue" \
  --label core \
  --label legacy \
  --by developer 2>&1)
printf '%s\n' "$CREATE_OUTPUT"
export HV13_LABEL_ISSUE="$(printf '%s\n' "$CREATE_OUTPUT" | sed -n 's/^created issue \([^ ]*\) .*/\1/p')"
printf 'label_issue=%s\n' "$HV13_LABEL_ISSUE"

"$FID" issue get "$HV13_LABEL_ISSUE" --format markdown

"$FID" issue label add "$HV13_LABEL_ISSUE" api --by developer
"$FID" issue get "$HV13_LABEL_ISSUE" --format markdown
```

Optionally narrow the before/after checks to the label section:

```bash
"$FID" issue get "$HV13_LABEL_ISSUE" --format json | grep '"labels"'
```

Expected for `HV-13b`:

- Before the add, the issue shows `core` and `legacy`.
- The add command reports `added label "api" to HV13-...`.
- After the add, the issue shows `api`, `core`, and `legacy`.
- The output reads as an additive operation, not a replace-all operation.

Record `HV-13b` as `pass` if the before/after output clearly shows the new
label plus the pre-existing labels. Record `partial` or `fail` if the output
does not make preservation of the existing labels clear.

## 5. HUMAN-13c - Remove One Label Without Disturbing Remaining Labels

Remove one label from the same issue and inspect the before/after issue output.

```bash
"$FID" issue get "$HV13_LABEL_ISSUE" --format markdown

"$FID" issue label remove "$HV13_LABEL_ISSUE" legacy --by developer
"$FID" issue get "$HV13_LABEL_ISSUE" --format markdown
```

Optionally narrow the after check to the label section:

```bash
"$FID" issue get "$HV13_LABEL_ISSUE" --format json | grep '"labels"'
```

Expected for `HV-13c`:

- Before the remove, the issue shows `api`, `core`, and `legacy`.
- The remove command reports `removed label "legacy" from HV13-...`.
- After the remove, the issue still shows `api` and `core`.
- The selected label is gone and the remaining labels are still visible.

Record `HV-13c` as `pass` if the before/after output clearly shows only the
selected label removed. Record `partial` or `fail` if another label appears to
be disturbed or the remaining-label state is hard to verify.

## 6. HUMAN-13d - Run The v0.1d Lifecycle With A Link And Label

Run the closure smoke through real commands: create -> claim -> transition ->
complete -> unclaim, with one label and one typed link on the primary issue.

```bash
PRIMARY_OUTPUT=$("$FID" issue create \
  --kind task \
  --title "HV13 primary lifecycle issue" \
  --label core \
  --by developer 2>&1)
printf '%s\n' "$PRIMARY_OUTPUT"
export HV13_PRIMARY_ISSUE="$(printf '%s\n' "$PRIMARY_OUTPUT" | sed -n 's/^created issue \([^ ]*\) .*/\1/p')"
printf 'primary_issue=%s\n' "$HV13_PRIMARY_ISSUE"

PEER_OUTPUT=$("$FID" issue create \
  --kind task \
  --title "HV13 peer lifecycle issue" \
  --by developer 2>&1)
printf '%s\n' "$PEER_OUTPUT"
export HV13_PEER_ISSUE="$(printf '%s\n' "$PEER_OUTPUT" | sed -n 's/^created issue \([^ ]*\) .*/\1/p')"
printf 'peer_issue=%s\n' "$HV13_PEER_ISSUE"

"$FID" issue claim "$HV13_PRIMARY_ISSUE" --by developer
"$FID" issue transition "$HV13_PRIMARY_ISSUE" --to in_progress --by developer
"$FID" link add "$HV13_PEER_ISSUE" discovered_from "$HV13_PRIMARY_ISSUE" --by developer
"$FID" issue complete "$HV13_PRIMARY_ISSUE" \
  --by developer \
  --notes-md notes/hv13-completion.md \
  --tokens 1600
"$FID" issue unclaim "$HV13_PRIMARY_ISSUE" --by developer
```

Inspect the final issue, link list, and claims section.

```bash
"$FID" issue get "$HV13_PRIMARY_ISSUE" --format markdown
"$FID" link list "$HV13_PRIMARY_ISSUE" --format json
"$FID" issue get "$HV13_PRIMARY_ISSUE" --section claims --format json
```

Expected for `HV-13d`:

- The command transcript reads as one lifecycle:
  `created issue`, `claimed`, `transitioned`, `completed`, and
  `released claim`.
- The final issue still shows label `core`.
- The link list shows an incoming `discovered_from` link from the peer issue.
- The claims section is empty after `unclaim`.
- The output is understandable as one complete end-to-end flow.

Record `HV-13d` as `pass` if the transcript and final inspection make the full
lifecycle, retained label, retained link, and released claim clear. Record
`partial` or `fail` if the flow is hard to follow or any final state is
unclear from CLI output.

## 7. Optional Lifecycle Rebuild Check

This is not required to fill the human checklist because deterministic QA
already covered rebuild parity, but it can help if the tester wants extra
confidence from the same disposable workspace.

```bash
"$FID" rebuild --by workflow_runner
"$FID" issue get "$HV13_PRIMARY_ISSUE" --format markdown
"$FID" link list "$HV13_PRIMARY_ISSUE" --format json
```

Expected:

- `fid rebuild` reports completion.
- The primary issue still shows label `core`.
- The link list still shows the `discovered_from` link.

## 8. Fill The Human Checklist

Return to the repository root and record the outcomes in the checklist.

```bash
cd /home/mark/projects/fid
sed -n '1,220p' project-docs/docs/phases/phase-13/qa/human-checklist.md
```

Fill:

- `HV-13a` from section 3.
- `HV-13b` from section 4.
- `HV-13c` from section 5.
- `HV-13d` from section 6.
- Summary totals and overall result.

Use `pass`, `partial`, or `fail` for each row. Overall result is `pass` only if
all four rows pass; `partial` if at least one row is partial and none fail;
`fail` if any row fails.

## 9. Cleanup

When the checklist evidence is recorded, remove the disposable workspace if it
is no longer needed.

```bash
printf 'workspace=%s\n' "$FID_HV13_WORKSPACE"
rm -rf "$FID_HV13_WORKSPACE"
```

Expected: the temporary `/tmp/fid-phase13-human-*` workspace is removed. Do not
run the cleanup command until any useful terminal output has been recorded in
the checklist notes.
