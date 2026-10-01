# Phase 23 - Task 19: Human-mode failures from the Stop-error family

> Supplemental task under DEC-268 / DEC-269 Part 2 (Human checklist finding A2).
> Dispatch only after DEC-269's status line reads ACCEPTED.

## Objective

Every Human-mode failure that goes through `renderStopError` exits non-zero
with nothing on stdout or stderr. `writeStopFailure` and the fallback return a
`FailureError` without writing, and `main` only maps the exit code. This task
makes that family print `error: <message>` on stderr, as Phase 21
`public-interface.md` section 1 already requires, with no contract change.

## Validation Contract Coverage

- VAL-085: every Human-mode failure returned through the Stop-error family
  writes exactly `error: <message>` (the JSON `message`) to stderr, with empty
  stdout and the same exit status as `--json`.
- VAL-072 (Human, re-checked at phase end): removes the silent
  `claim acquire` denial the Human reported under HUMAN-03a. Task 21 adds the
  rework guidance line.

## Binding Constraints

- Phase 21 `public-interface.md` section 1: Human failure goes to stderr only,
  as `error: <message>` plus at most two non-semantic guidance lines. This
  task adds **no** guidance line. Task 21 owns the only one (VAL-088).
- The Human message is byte-identical to the failure's JSON `message`. The
  following are unchanged: JSON envelopes, failure codes, categories, exit
  classes, `details`, and subjects.
- Phase 23 adds no denial Event (DEC-226). A Stop-blocked request still
  commits nothing.
- `gate run` keeps the `fid.live.v1` contract. Before the first frame, failures
  use the ordinary failure channel. This task changes only that pre-frame
  Human rendering.
- Other leaves already render their own Human failures. Do not refactor them.

## Scope

### Files To Create Or Modify

- `cmd/fid/stop.go` (`writeStopFailure`, `renderStopError` fallback)
- `cmd/fid/stop_error_test.go`
- `cmd/fid/integration/p23_human_failure_test.go` (new)
- Any retained Human-mode assertion that expected empty stderr for these paths
  converges here, with per-row provenance (`execution-scope.md` section 9)

### Files To Read For Context

- `project-docs/docs/phases/phase-23/decisions.md` DEC-269 Part 2
- `project-docs/docs/phases/phase-21/public-interface.md` section 1
- `cmd/fid/claim.go`, `cmd/fid/run.go`, `cmd/fid/gate.go`, `cmd/fid/attention.go`
  (callers), `internal/cli/exit.go`
- `project-docs/docs/phases/phase-23/qa/human-checklist.md` section "Tech Lead
  analysis", finding A2

### Out Of Scope

- Guidance lines, including the rework guidance (task 21).
- Stop-details rendering in `stop list`, `attention list`, and `resume`
  (finding A4, wrap-up agenda).
- Human-mode failures outside the Stop-error family. If the table finds one,
  record it as a concern in the report; do not fix it.

## Acceptance Criteria

- [ ] A built-binary table covers each failure code `renderStopError` can
  return, per calling leaf, in Human and `--json` mode:
  - `claim acquire` and `run start` denied by an open
    `review-rework-exhausted` Stop, and by a `stop-unresolved` Stop;
  - the `gate run` Stop denial;
  - `stop resolve` and `attention resolve` not-found, version conflict,
    `resolution-invalid`, and invalid argument;
  - `attention list` invalid argument.

  Every Human case writes stderr beginning `error: ` followed by the JSON
  `message`, writes empty stdout, and exits with the same status as its JSON
  twin. (VAL-085)
- [ ] The table is enumerated from the codes `renderStopError` constructs, not
  hand-picked. Adding a new code to that function without a table row fails a
  test. (VAL-085)
- [ ] JSON output for every case is byte-identical before and after this
  change. (VAL-085)

## First Failing Test

The built-binary `claim acquire` in Human mode against an exhausted rework
fixture expects stderr `error: review rework budget is exhausted`. It gets
empty stderr today.

## Documentation Updates

- None expected. No contract changes. The Tech Lead records spec impact at
  wrap-up.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
```

Two gates run for minutes and are explicitly bounded (DEC-227): the `-timeout
30m` above replaces Go's implicit 10-minute per-binary default, which the tagged
suite already approaches at 560 s. Record in the completion report the wall time
`go test` reports for `cmd/fid/integration`, `cmd/fid`, and `internal/workflow`,
so a regression is attributed to the task that caused it. Budgets: tagged
integration under 900 s and no ordinary package over 600 s (VAL-079). A breach is
a bounded speed repair owned by this task when this task's own new journeys
caused it; otherwise stop and escalate. Never record a timed-out gate as `PASS`.

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. No exception to the full task
gate. Focused red/green selections never replace these commands. Report each
attempt with the command's own output token and observed termination. An
interrupted or hung command is never `PASS`; preserve all failed attempts and
rerun the byte-identical command after an authorized causal repair. Report lint
is a Runner control after the report commit, not a Developer gate.

List any conditional command required by this task and why:

- None.
## Notes For Developer

- Report path: `project-docs/docs/phases/phase-23/tasks/task-19-done.md`.
  Commit with `Refs: phase-23/task-19`.
- Red first: record the failing built-binary case before changing
  `stop.go`.
- Any missing registry, contract ambiguity, or need for a guidance line is a
  `BLOCKED` escalation candidate, not a local decision.
