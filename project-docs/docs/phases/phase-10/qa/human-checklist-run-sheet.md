# Phase 10 Human Checklist Run Sheet

> Phase: 10 - Permissions Matrix + Policy
> Purpose: executable command sheet for
> `project-docs/docs/phases/phase-10/qa/human-checklist.md`
> Date prepared: 2026-06-21

Use this from the fid repository root. It creates a disposable real `.fid/`
workspace under `/tmp`, builds the current CLI to `/tmp/fid`, displays the
permissions matrix, then triggers one denied write for the human actionability
checks.

The human checklist still needs human judgment. Treat the command output as
dogfood evidence, then record the outcome and notes in `human-checklist.md`.

## Prerequisites

- Run the build command from the repo root.
- Run workspace commands only inside the disposable `/tmp` workspace created
  below, not inside this repository's own `.fid/` workspace.
- These commands require the built `fid` binary and POSIX shell tools already
  used by the project docs (`grep`, `sed`, `printf`, `ls`, `mktemp`).
- This sheet does not complete deterministic QA gates. Phase 10 final review
  already covers JSON validity, format parity, read-only behavior,
  show-vs-enforce parity, and validate-before-append mechanics. This sheet only
  supports `HUMAN-09a` through `HUMAN-09d` for `VAL-009`.

## 1. Build Current CLI

```bash
go build -o /tmp/fid ./cmd/fid
export FID=/tmp/fid
ls -l "$FID"
```

Expected: command exits zero and `/tmp/fid` exists.

## 2. Create A Fresh Test Workspace

```bash
export FID_HV10_WORKSPACE="$(mktemp -d /tmp/fid-phase10-human-XXXXXX)"
cd "$FID_HV10_WORKSPACE"
printf 'workspace=%s\n' "$FID_HV10_WORKSPACE"

"$FID" init --prefix HV10
```

Expected:

- `fid init` reports a new `.fid/` workspace.
- The workspace has a `.fid/` directory and seeded roles.
- The rest of this sheet is run from the disposable workspace.

## 3. HUMAN-09a - Inspect The Human-Readable Policy Matrix

Run the markdown view and inspect it as a human operator.

```bash
"$FID" policy show --format markdown
```

Optionally narrow the terminal output to the denied/permitted row used later:

```bash
"$FID" policy show --format markdown | grep 'issue.create.epic'
```

Expected for `HUMAN-09a`:

- The output reads as a role x action matrix.
- You can find a specific action without ambiguity.
- You can compare allowed roles for `issue.create.epic`.
- The `developer` cell is not permitted and the `tl` and `human` cells are
  permitted.

Record `HUMAN-09a` as `pass` if the human-readable matrix is easy enough to scan
and compare. Record `partial` or `fail` if action rows, role columns, or boolean
cells are hard to find or interpret.

## 4. HUMAN-09b - Inspect The JSON Policy Matrix

Run the JSON view and judge scanability only. Do not treat this as a JSON
validity test; deterministic QA already covered that.

```bash
"$FID" policy show --format json
```

Optionally narrow the terminal output to the same row:

```bash
"$FID" policy show --format json | grep -o '{[^}]*"action":"issue.create.epic"[^}]*}'
```

Expected for `HUMAN-09b`:

- Action names are discoverable.
- Role booleans are understandable enough for an operator or agent to read.
- The `issue.create.epic` row makes `developer:false`, `tl:true`, and
  `human:true` visible enough for practical use.

Record `HUMAN-09b` as `pass` if the JSON is understandable enough to locate an
action and read true/false role cells. Record `partial` or `fail` if the JSON is
valid but too awkward for practical operator use.

## 5. HUMAN-09c - Trigger A Denied Write And Inspect The Error

Use a denied write in the disposable workspace. `developer` is a known role, but
the Phase 10 matrix denies `developer` for `issue.create.epic`.

```bash
DENIED_OUTPUT=$("$FID" issue create \
  --kind epic \
  --title "HV10 denied epic" \
  --scheduling sequential \
  --by developer 2>&1)
DENIED_RC=$?
printf 'denied_rc=%s\n' "$DENIED_RC"
printf '%s\n' "$DENIED_OUTPUT"
```

Expected for `HUMAN-09c`:

- `denied_rc` is non-zero.
- The error names role `developer`.
- The error names action `issue.create.epic`.
- The error says the role is not permitted.

Record `HUMAN-09c` as `pass` if the error clearly identifies the denied role,
the denied action, and not-permitted status. Record `partial` or `fail` if any
of those facts are missing, buried, or ambiguous.

## 6. HUMAN-09d - Judge Whether The Error Gives A Usable Next Step

Read the denied-write output from the previous step again and decide whether it
points the operator toward a useful next action.

```bash
printf '%s\n' "$DENIED_OUTPUT"
```

Expected for `HUMAN-09d`:

- The output gives enough signal to choose a next step.
- A reasonable next step is to inspect `fid policy show`, switch to a permitted
  role, or escalate a role/action mismatch.
- The message should not leave the operator guessing whether this is a missing
  role, invalid issue kind, malformed command, or permissions denial.

To confirm the "switch to a permitted role" next step while the workspace is
already set up, run the permitted counterpart:

```bash
PERMITTED_OUTPUT=$("$FID" issue create \
  --kind epic \
  --title "HV10 permitted epic" \
  --scheduling sequential \
  --by tl 2>&1)
PERMITTED_RC=$?
printf 'permitted_rc=%s\n' "$PERMITTED_RC"
printf '%s\n' "$PERMITTED_OUTPUT"
```

Expected for the optional counterpart:

- `permitted_rc` is zero.
- The command reports `created issue HV10-... (epic)`.
- This confirms that changing to a matrix-permitted role is one viable path.

Record `HUMAN-09d` as `pass` if the denial message is actionable enough to pick
a next step. Record `partial` or `fail` if the message is technically accurate
but does not make the next move clear.

## 7. Optional Format Sweep

This is not required to fill the human checklist, but it can help if the
markdown view is not the tester's preferred human-readable format.

```bash
"$FID" policy show --format toon
"$FID" policy show --format yaml
```

Expected:

- Each format exits zero.
- Each format exposes the same role/action matrix.
- At least one human-readable format is comfortable enough for future
  operator-facing policy inspection.

## 8. Fill The Human Checklist

Return to the repository root and record the outcomes in the checklist.

```bash
cd /home/mark/projects/fid
sed -n '1,220p' project-docs/docs/phases/phase-10/qa/human-checklist.md
```

Fill:

- `HUMAN-09a` from section 3.
- `HUMAN-09b` from section 4.
- `HUMAN-09c` from section 5.
- `HUMAN-09d` from section 6.
- Summary totals and overall result.

Use `pass`, `partial`, or `fail` for each row. Overall result is `pass` only if
all four rows pass; `partial` if at least one row is partial and none fail;
`fail` if any row fails.

## 9. Cleanup

When the checklist evidence is recorded, remove the disposable workspace if it
is no longer needed.

```bash
printf 'workspace=%s\n' "$FID_HV10_WORKSPACE"
rm -rf "$FID_HV10_WORKSPACE"
```

Expected: the temporary `/tmp/fid-phase10-human-*` workspace is removed. Do not
run the cleanup command until any useful terminal output has been recorded in
the checklist notes.
