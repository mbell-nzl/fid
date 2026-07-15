# Phase <N> - Task <M>: <Title> - Completion Report

## Summary of what was implemented

- <summary>

## Files created or modified

- `<path>`

## Documentation updated

- `<path>` or `None required.`

## Validation Contract Coverage

- VAL-NNN: <evidence produced, status, or reason this task does not satisfy it>
- <If none apply, explain why explicitly.>

## Tests run

- `<command>` - PASS/FAIL

## Exact quality-gate results

Compare this section against the task spec's `## Exact Commands To Run`.

- `gofmt -l cmd internal` - PASS/FAIL/BLOCKED
- `golangci-lint run ./...` - PASS/FAIL/BLOCKED
- `go test -race -coverprofile=coverage.out ./...` - PASS/FAIL/BLOCKED
- `go tool cover -func=coverage.out` - PASS/FAIL/BLOCKED
- `<conditional command from task spec>` - PASS/FAIL/BLOCKED/N/A

Every required command from the task spec and QA plan must appear here or in
`## Tests run`. If a command is blocked, record the reason. Do not substitute a
narrower command for a required gate without a blocker report; silently skipped
or narrower substituted commands are failed gates.

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
