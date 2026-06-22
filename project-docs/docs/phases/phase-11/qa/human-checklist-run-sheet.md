# Phase 11 Human Checklist Run Sheet

> Phase: 11 - Claims Layer
> Purpose: executable command sheet for
> `project-docs/docs/phases/phase-11/qa/human-checklist.md`
> Date prepared: 2026-06-21

Use this from the fid repository root. It creates a disposable real `.fid/`
workspace under `/tmp`, builds the current CLI to `/tmp/fid`, walks through a
claim lifecycle, then triggers one human-held-claim transition block for the
human actionability checks.

The human checklist still needs human judgment. Treat the command output as
dogfood evidence, then record the outcome and notes in `human-checklist.md`.
For `HV-01e`, prefer repeating the same lifecycle during normal project work if
you have a suitable real fid workspace; otherwise record that the judgment came
from this disposable workspace.

## Prerequisites

- Run the build command from the repo root.
- Run workspace commands only inside the disposable `/tmp` workspace created
  below, not inside this repository's own `.fid/` workspace.
- These commands require the built `fid` binary and POSIX shell tools already
  used by the project docs (`grep`, `sed`, `printf`, `ls`, `mktemp`).
- This sheet does not complete deterministic QA gates. Phase 11 final review
  already covers claim mechanics, transition gating, stale-projection behavior,
  JSON/format parity, coverage, integration, benchmarks, and build gates. This
  sheet only supports `HV-01a` through `HV-01e` for `VAL-014`.

## 1. Build Current CLI

```bash
go build -o /tmp/fid ./cmd/fid
export FID=/tmp/fid
ls -l "$FID"
```

Expected: command exits zero and `/tmp/fid` exists.

## 2. Create A Fresh Test Workspace

```bash
export FID_HV11_WORKSPACE="$(mktemp -d /tmp/fid-phase11-human-XXXXXX)"
cd "$FID_HV11_WORKSPACE"
printf 'workspace=%s\n' "$FID_HV11_WORKSPACE"

"$FID" init --prefix HV11
```

Expected:

- `fid init` reports a new `.fid/` workspace.
- The workspace has a `.fid/` directory and seeded roles.
- The rest of this sheet is run from the disposable workspace.

## 3. Create A Claimable Issue

```bash
CREATE_OUTPUT=$("$FID" issue create \
  --kind task \
  --title "HV11 claim lifecycle" \
  --by developer 2>&1)
printf '%s\n' "$CREATE_OUTPUT"

export HV11_ISSUE_ID="$(printf '%s\n' "$CREATE_OUTPUT" | sed -n 's/created issue \([^ ]*\).*/\1/p')"
printf 'issue_id=%s\n' "$HV11_ISSUE_ID"
```

Expected:

- `fid issue create` reports `created issue HV11-... (task)`.
- `issue_id` is populated with the created issue ID.

Stop and fix the setup before continuing if `issue_id` is blank.

## 4. HV-01a - Identify The Claim Holder

Acquire the issue as `developer`, then inspect the claim list as a human
operator.

```bash
CLAIM_OUTPUT=$("$FID" issue claim "$HV11_ISSUE_ID" \
  --by developer \
  --ttl 15m \
  --note "human checklist holder/expiry" 2>&1)
printf '%s\n' "$CLAIM_OUTPUT"

"$FID" claim list --by developer --format markdown
```

Expected for `HV-01a`:

- The claim command reports `claimed HV11-... by developer`.
- The markdown list shows one claim for the issue.
- The holder role `developer` is easy to identify without reading code or the
  database.

Record `HV-01a` as `pass` if the holder is clear from the command output or
read-back. Record `partial` or `fail` if the holder is hidden, ambiguous, or
requires implementation inspection.

## 5. HV-01b - Identify The Claim Expiry

Inspect the same claim output and read-back for expiry information.

```bash
printf '%s\n' "$CLAIM_OUTPUT"
"$FID" claim list --by developer --format markdown
"$FID" claim list --by developer --format json
```

Expected for `HV-01b`:

- The claim command includes an `expires` timestamp.
- The markdown list includes `ttl 900s` and an `expires` timestamp.
- The JSON list includes `ttl_seconds`, `expires_at`, and `is_stale`.
- The expiry is understandable enough for a human operator to know roughly when
  the claim stops protecting the issue.

Record `HV-01b` as `pass` if the expiry is understandable without code
inspection. Record `partial` or `fail` if the timestamp, TTL, or stale state is
too hard to interpret in practical use.

## 6. HV-01e - Exercise Renew And List Usability

Renew the claim and inspect the claim list again. This supports the overall
claim lifecycle dogfood row; the final product judgment happens in section 10.

```bash
RENEW_OUTPUT=$("$FID" issue claim "$HV11_ISSUE_ID" \
  --by developer \
  --renew \
  --ttl 20m 2>&1)
printf '%s\n' "$RENEW_OUTPUT"

"$FID" claim list --by developer --format markdown
```

Expected:

- The renew command reports `renewed claim on HV11-... by developer`.
- The expiry changes.
- The markdown list remains understandable and still shows the same holder.

Keep notes for `HV-01e` about whether the claim, renew, and list messages feel
usable for agent coordination.

## 7. HV-01c - Release The Claim And Confirm Release

Release the developer-held claim and inspect read-back.

```bash
UNCLAIM_OUTPUT=$("$FID" issue unclaim "$HV11_ISSUE_ID" --by developer 2>&1)
printf '%s\n' "$UNCLAIM_OUTPUT"

"$FID" issue get "$HV11_ISSUE_ID" --section claims --format json
"$FID" claim list --by developer --format markdown
```

Expected for `HV-01c`:

- The unclaim command reports `released claim on HV11-... by developer`.
- `issue get --section claims --format json` returns `[]`.
- `claim list --by developer --format markdown` shows no current developer
  claim for this issue.
- The release state is clear enough to trust that the issue no longer has the
  developer-held current claim.

Record `HV-01c` as `pass` if release is clear from the command output or
read-back. Record `partial` or `fail` if the operator cannot confidently tell
whether the issue is still claimed.

## 8. HV-01d - Trigger A Human-Held Transition Block

Acquire the issue as `human`, then try to advance it as `developer`.

```bash
HUMAN_CLAIM_OUTPUT=$("$FID" issue claim "$HV11_ISSUE_ID" \
  --by human \
  --ttl 15m \
  --note "human checklist block" 2>&1)
printf '%s\n' "$HUMAN_CLAIM_OUTPUT"

BLOCK_OUTPUT=$("$FID" issue transition "$HV11_ISSUE_ID" \
  --to in_progress \
  --by developer 2>&1)
BLOCK_RC=$?
printf 'block_rc=%s\n' "$BLOCK_RC"
printf '%s\n' "$BLOCK_OUTPUT"

"$FID" claim list --format markdown
```

Expected for `HV-01d`:

- The human claim command reports `claimed HV11-... by human`.
- `block_rc` is non-zero.
- The transition error says the live claim is not held by `developer`.
- The transition error names the current holder `human`.
- The transition error includes the claim expiry.
- The markdown list still shows the issue held by `human`.

Record `HV-01d` as `pass` if the block reason clearly explains that a
human-held live claim prevents the agent role from advancing the issue. Record
`partial` or `fail` if the error leaves the operator guessing whether this is a
workflow rule, permission denial, missing issue, or claim ownership problem.

## 9. Optional Human Counterpart

This is not required to fill `HV-01d`, but it can help confirm that the block is
about claim ownership rather than the transition itself.

```bash
HUMAN_TRANSITION_OUTPUT=$("$FID" issue transition "$HV11_ISSUE_ID" \
  --to in_progress \
  --by human 2>&1)
HUMAN_TRANSITION_RC=$?
printf 'human_transition_rc=%s\n' "$HUMAN_TRANSITION_RC"
printf '%s\n' "$HUMAN_TRANSITION_OUTPUT"

"$FID" claim list --format json
```

Expected for the optional counterpart:

- `human_transition_rc` is zero.
- The command reports a transition from `todo` to `in_progress`.
- The claim remains current because release happens when leaving a
  started-category state, not when entering one.

If you run this optional step, keep its output separate from the required
`HV-01d` judgment so a successful human transition does not hide whether the
developer block message was clear.

## 10. HV-01e - Overall Claim Lifecycle Dogfood Judgment

Review the outputs from sections 4 through 8, and optionally repeat the same
flow in a normal project-work `.fid` workspace if you have one available.

```bash
printf '%s\n' "$CLAIM_OUTPUT"
printf '%s\n' "$RENEW_OUTPUT"
printf '%s\n' "$UNCLAIM_OUTPUT"
printf '%s\n' "$HUMAN_CLAIM_OUTPUT"
printf '%s\n' "$BLOCK_OUTPUT"
```

Expected for `HV-01e`:

- Claim acquisition gives enough holder and expiry signal.
- Renewal reads as a continuation of the same claim lifecycle.
- Unclaim reads as a release, not as an unrelated issue update.
- `claim list` is useful enough for coordination between agents and humans.
- The human-held block message gives an actionable explanation.

Record `HV-01e` as `pass` if the product-owner/dogfood judgment is accept.
Record `partial` or `fail` with notes if the mechanics work but the messages,
flow, or coordination affordances should change.

## 11. Fill The Human Checklist

Return to the repository root and record the outcomes in the checklist.

```bash
cd /home/mark/projects/fid
sed -n '1,220p' project-docs/docs/phases/phase-11/qa/human-checklist.md
```

Fill:

- `HV-01a` from section 4.
- `HV-01b` from section 5.
- `HV-01c` from section 7.
- `HV-01d` from section 8.
- `HV-01e` from sections 6 and 10.
- Summary totals and overall result.

Use `pass`, `partial`, or `fail` for each row. Overall result is `pass` only if
all five rows pass; `partial` if at least one row is partial and none fail;
`fail` if any row fails.

## 12. Cleanup

When the checklist evidence is recorded, remove the disposable workspace if it
is no longer needed.

```bash
printf 'workspace=%s\n' "$FID_HV11_WORKSPACE"
rm -rf "$FID_HV11_WORKSPACE"
```

Expected: the temporary `/tmp/fid-phase11-human-*` workspace is removed. Do not
run the cleanup command until any useful terminal output has been recorded in
the checklist notes.
