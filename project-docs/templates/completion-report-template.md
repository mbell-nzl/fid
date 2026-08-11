# Phase <N> - Task <M>: <Title> - Completion Report

## Summary of what was implemented

- <summary>

## Files created or modified

- `<path>`

## Documentation updated

- `<path>` or `None required.`

## Validation Contract Coverage

- VAL-NNN: <PASS/FAIL/BLOCKED> - <the specific test name, `file:line`, or
  command output that proves it>
- <If none apply, explain why explicitly.>

Cite the artifact that proves the assertion, not the belief that it holds.
"Covered by the existing suite", "behavior is unchanged", or a restatement of
the assertion is not evidence. If the assertion names two states (a
conditional, a both-directions parity, a valid/invalid matrix), cite evidence
for **each** state; positive-only evidence is not a PASS.

## Tests run

- `<command>` - PASS/FAIL - `<verbatim final output line>` (terminated `exit_code: <N>`)

## Exact quality-gate results

Compare this section against the task spec's `## Exact Commands To Run`.

- `gofmt -l cmd internal` - PASS/FAIL/BLOCKED - `<verbatim final output line>` (terminated `exit_code: <N>`)
- `golangci-lint run ./...` - PASS/FAIL/BLOCKED - `<verbatim final output line>` (terminated `exit_code: <N>`)
- `go test -race -coverprofile=coverage.out ./...` - PASS/FAIL/BLOCKED - `<verbatim final output line>` (terminated `exit_code: <N>`)
- `go tool cover -func=coverage.out` - PASS/FAIL/BLOCKED - `<verbatim final output line>` (terminated `exit_code: <N>`)
- `<conditional command from task spec>` - PASS/FAIL/BLOCKED/N/A - `<verbatim final output line>` (terminated `exit_code: <N>`)

A real example row:

```text
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
```

Every required command from the task spec and QA plan must appear here or in
`## Tests run`. If a command is blocked, record the reason. Do not substitute a
narrower command for a required gate without a blocker report; silently skipped
or narrower substituted commands are failed gates.

Each result carries one short verbatim token copied from the command's own
terminal output — the last line it printed, or `exit 0 (no output)` for a
silent success. Copy it; do not summarize or reconstruct it.

Each result also carries the **observed termination result**: the command's
real `exit_code`, or the explicit signal/termination result your execution
surface reported when it exposes no exit code. A printed last line proves the
command produced output, not that the process exited; for a silent success the
exit status is the only evidence that separates it from a discarded session.
Copy the termination result from the command's own result object. Wrapper
completion, an empty output field, and absence of new output are not
termination, and an `exit_code` you assert rather than observed is fabricated
evidence exactly like a fabricated token.

**A command you did not watch terminate is never `PASS`.** A command that hung,
was interrupted, was killed, or whose output you never saw complete has no
final line to copy, and inventing one is fabricated evidence. Record it as
`BLOCKED` with the reason in a blocker report and stop. A gate recorded `PASS`
without terminating is invisible to every downstream review until someone
re-runs it independently.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `<hash or "No commit created - blocker report">`

Draft this section as `- Pending commit` until the task commit exists. Then,
while `HEAD` is still the task commit, copy the hash verbatim from
`git rev-parse --short HEAD` or `git rev-parse HEAD`. Never type, pad,
complete, or recall a hash from memory.

## Commit subject

- `<subject or "No commit created - blocker report">`

Copy the subject verbatim from `git log -1 --format=%s`.

When replacing either placeholder, retain the surrounding Markdown backticks.
The finalized fields must use this exact shape:

```text
## Commit hash
- `0123456789abcdef`

## Commit subject
- `feat(scope): subject copied from git`
```

After filling in the hash and subject, commit the finalized report in an
immediately following `docs(task): finalize completion report for phase-N
task-M` commit. The hash above is the task commit's, never that docs commit's.
The report counts as evidence only once committed; report-lint fails a report
file with uncommitted modifications.
