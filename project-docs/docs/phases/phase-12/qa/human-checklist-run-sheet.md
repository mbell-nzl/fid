# Phase 12 Human Checklist Run Sheet

> Phase: 12 - Requirements + Links
> Purpose: executable command sheet for
> `project-docs/docs/phases/phase-12/qa/human-checklist.md`
> Date prepared: 2026-06-27

Use this from the fid repository root. It creates a disposable real `.fid/`
workspace under `/tmp`, builds the current CLI to `/tmp/fid`, syncs a tiny
PRD-style requirement source, creates one issue, then walks through an
issue-to-requirement link lifecycle.

The human checklist still needs human judgment. Treat the command output as
dogfood evidence, then record the outcome and notes in `human-checklist.md`.
This run sheet only supports `HV-12a` through `HV-12d` for `VAL-012`; it does
not replace the deterministic Phase 12 QA gates in `final-review.md`.

## Prerequisites

- Run the build command from the repo root.
- Run workspace commands only inside the disposable `/tmp` workspace created
  below, not inside this repository's own `.fid/` workspace.
- These commands require the built `fid` binary and POSIX shell tools already
  used by the project docs (`cat`, `sed`, `printf`, `ls`, `mktemp`).
- Phase 12 deterministic QA already covers command mechanics, journal event
  counts, rebuild parity, format parity, and rejection paths. This sheet only
  supports the human legibility checks for requirement sync, inverse creation,
  bidirectional listing, and inverse removal.

## 1. Build Current CLI

```bash
go build -o /tmp/fid ./cmd/fid
export FID=/tmp/fid
ls -l "$FID"
```

Expected: command exits zero and `/tmp/fid` exists.

## 2. Create A Fresh Test Workspace

```bash
export FID_HV12_WORKSPACE="$(mktemp -d /tmp/fid-phase12-human-XXXXXX)"
cd "$FID_HV12_WORKSPACE"
printf 'workspace=%s\n' "$FID_HV12_WORKSPACE"

"$FID" init --prefix HV12
```

Expected:

- `fid init` reports a new `.fid/` workspace.
- The workspace has a `.fid/` directory and seeded roles.
- The rest of this sheet is run from the disposable workspace.

## 3. Create A PRD-Style Source Document

```bash
export HV12_SOURCE_DOC="$FID_HV12_WORKSPACE/phase12-requirements.md"
cat > "$HV12_SOURCE_DOC" <<'EOF'
# Human Checklist Source

## Link Requirements

| ID | Summary |
|---|---|
| PRD-HV12-001 | Human checklist synced requirement pointer |
EOF

sed -n '1,20p' "$HV12_SOURCE_DOC"
```

Expected:

- The source document contains one parseable PRD requirement ID:
  `PRD-HV12-001`.
- The second table cell is the summary that should later appear as the linked
  requirement title.

## 4. HV-12a - Sync Requirement Pointers

Sync the source document once, then sync it again unchanged so both created and
unchanged counts are visible to the human tester.

```bash
SYNC_FIRST_OUTPUT=$("$FID" requirement sync \
  --source-doc "$HV12_SOURCE_DOC" \
  --by tl 2>&1)
printf '%s\n' "$SYNC_FIRST_OUTPUT"

SYNC_SECOND_OUTPUT=$("$FID" requirement sync \
  --source-doc "$HV12_SOURCE_DOC" \
  --by tl 2>&1)
printf '%s\n' "$SYNC_SECOND_OUTPUT"
```

Expected for `HV-12a`:

- The first sync reports `synced requirements from ...`.
- The first sync shows `created: 1`, `updated: 0`, and `unchanged: 0`.
- The second sync shows `created: 0`, `updated: 0`, and `unchanged: 1`.
- The counts are understandable without reading implementation or inspecting
  the database.

Record `HV-12a` as `pass` if the created, updated, and unchanged counts are
clear enough to explain what happened to the requirement pointer. Record
`partial` or `fail` if the output hides, confuses, or overstates those counts.

## 5. Create A Linkable Issue

```bash
CREATE_OUTPUT=$("$FID" issue create \
  --kind task \
  --title "HV12 primary issue" \
  --by developer 2>&1)
printf '%s\n' "$CREATE_OUTPUT"

export HV12_ISSUE_ID="$(printf '%s\n' "$CREATE_OUTPUT" | sed -n 's/created issue \([^ ]*\).*/\1/p')"
printf 'issue_id=%s\n' "$HV12_ISSUE_ID"
```

Expected:

- `fid issue create` reports `created issue HV12-... (task)`.
- `issue_id` is populated with the created issue ID.

Stop and fix the setup before continuing if `issue_id` is blank.

## 6. HV-12b - Add A Link And Identify The Inverse

Create an issue-to-requirement `implements` link. This should also create the
auto-inverse requirement-to-issue `implemented_by` row.

```bash
LINK_ADD_OUTPUT=$("$FID" link add "$HV12_ISSUE_ID" implements PRD-HV12-001 \
  --by developer 2>&1)
printf '%s\n' "$LINK_ADD_OUTPUT"
```

Expected for `HV-12b`:

- The command reports `linked HV12-... implements PRD-HV12-001`.
- The same output includes `inverse: PRD-HV12-001 implemented_by HV12-...`.
- A human tester can identify both the declared direction and the auto-created
  inverse direction from the CLI output alone.

Record `HV-12b` as `pass` if the declared and inverse directions are clear.
Record `partial` or `fail` if the inverse is hidden, ambiguous, or requires
implementation knowledge to understand.

## 7. HV-12c - List Both Linked Endpoints

Inspect the link from the issue endpoint and from the requirement endpoint.
Markdown is the normal human read-back; JSON is included as a precise field
reference if the human wants to verify the row shape.

```bash
"$FID" link list "$HV12_ISSUE_ID" --format markdown
"$FID" link list PRD-HV12-001 --format markdown

"$FID" link list "$HV12_ISSUE_ID" --format json
"$FID" link list PRD-HV12-001 --format json
```

Expected for `HV-12c`:

- Listing the issue shows an `outgoing` `implements` link to requirement
  `PRD-HV12-001`.
- Listing the requirement shows an `incoming` `implements` link from the issue.
- The link rows expose direction, relation, linked entity type, linked entity
  ID, and linked entity title.
- The issue-side linked title is the requirement summary:
  `Human checklist synced requirement pointer`.
- The requirement-side linked title is the issue title:
  `HV12 primary issue`.

Record `HV-12c` as `pass` if outgoing and incoming directions are
understandable in normal CLI use from both endpoint perspectives. Record
`partial` or `fail` if the direction, relation, endpoint type, endpoint ID, or
title is missing or hard to interpret.

## 8. HV-12d - Remove The Link And Identify Both Removed Directions

Remove the declared issue-to-requirement link, then inspect both endpoints
again.

```bash
LINK_REMOVE_OUTPUT=$("$FID" link remove "$HV12_ISSUE_ID" implements PRD-HV12-001 \
  --by developer 2>&1)
printf '%s\n' "$LINK_REMOVE_OUTPUT"

"$FID" link list "$HV12_ISSUE_ID" --format json
"$FID" link list PRD-HV12-001 --format json
```

Expected for `HV-12d`:

- The remove command reports `removed link HV12-... implements PRD-HV12-001`.
- The same output includes
  `inverse removed: PRD-HV12-001 implemented_by HV12-...`.
- The post-remove JSON lists for both endpoints contain no remaining link rows
  for this pair.
- A human tester can identify that both the declared link and inverse link were
  removed from the CLI output alone.

Record `HV-12d` as `pass` if removal of both directions is clear. Record
`partial` or `fail` if the output leaves doubt about whether the inverse was
removed or whether either endpoint still has the link.

## 9. Overall Notes For `human-checklist.md`

Review the outputs from sections 4, 6, 7, and 8 before filling the checklist.

```bash
printf '%s\n' "$SYNC_FIRST_OUTPUT"
printf '%s\n' "$SYNC_SECOND_OUTPUT"
printf '%s\n' "$LINK_ADD_OUTPUT"
printf '%s\n' "$LINK_REMOVE_OUTPUT"
```

Expected overall dogfood judgment:

- Requirement sync counts make pointer changes understandable.
- Link add explains both the declared relation and the inverse relation.
- Link list is useful from both endpoint perspectives.
- Link remove explains that both directions were removed.

Record the final `VAL-012` human evidence in
`project-docs/docs/phases/phase-12/qa/human-checklist.md`.
