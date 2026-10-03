# Phase 24 - Task 10: Human failure lines in every family

## Objective

Replace the Phase 21 "error plus two guidance lines" rule with contract
section 9. Every Human-mode failure prints `error:` with the exact JSON
message, then `reason:`, then `blocked by:` and `next:` where they apply.
Every line derives from the JSON failure. One shared renderer is adopted by
every family.

## Validation Contract Coverage

- VAL-036: per code, the Human lines repeat the JSON failure.
- VAL-037: `blocked by:` appears exactly when details carry blocking IDs.
- VAL-038: `next:` appears exactly for next-table codes.
- VAL-039: `reason:` key-detail rendering, escaping and the 200-scalar cut.
- VAL-052.
- Advances VAL-054 (Human walk-through at the task-11 window).

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Contract section 9 is exact. It replaces Phase 21 `public-interface.md`
  section 1's guidance-line rule. JSON failures, exit codes and stdout are
  unchanged.
- **One renderer.** A single function in `internal/cli` takes the
  `fid-cli.v1` failure (code, message, subject, details) plus the operation
  and leaf path, and writes the lines. Every family writer (inventory R9)
  calls it in Human mode instead of `renderIssueError(cause)` or ad hoc
  text. This includes the standalone families (issue, label, link,
  workflow, init, rebuild, maintenance, export, capabilities and
  unavailable). A family that formats Human failures without a JSON failure
  constructs the same failure it would emit in JSON mode.
- **`reason:`** is `<code>`, optionally followed by `; <key detail>`.
  Scalar details render as `name=value` in key order. Scalar arrays are
  comma-joined. The closed overrides are exactly those in contract section
  9. Control characters are escaped, and the line is cut to 200 Unicode
  scalars with `…`.
- **`blocked by:`** lists `stop_ids` then `attention_ids` when either is
  nonempty.
- **`next:`** uses the closed table in contract section 9 with values from
  details. `<leaf>` is the CLI path (for example `fid run close`). The
  Phase 23 task 21 rework guidance line is replaced by its table row.
- **`gate run`.** Pre-frame failures use these lines. A terminal non-success
  in Human mode prints them after the streamed output.
- `fid journal inspect` and `fid sql` are unstable Human-only diagnostics
  outside `fid-cli.v1` and are excluded.

## Scope

### Files To Create Or Modify

- `internal/cli/` Human failure renderer (new), with tests
- Every function listed by inventory R9 (41 in `cmd/fid`)
- Retained Human stderr assertions (R10: 32 in 13 files): convergence with
  provenance
- A built-binary Human/JSON pairing corpus (new), enumerated from the codes
  with CLI producers

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` section 9
- `project-docs/docs/phases/phase-24/input-inventory.md` R8, R9, R10
- `project-docs/docs/phases/phase-24/decisions.md` DEC-272
- `project-docs/docs/phases/phase-21/public-interface.md` sections 1, 7;
  `project-docs/specs/05_DESIGN_PLAN.md` section 6
- `cmd/fid/stop.go` `writeStopFailure` (the task 19/21 guidance)

### Out Of Scope

- Stop, Attention and repository views (task 11). Judgment views (task 13).
- Changing any failure message, code or detail.

## Acceptance Criteria

- [ ] The pairing corpus is enumerated from every failure code with a CLI
  producer: Human stderr is `error: <JSON message>` then `reason: <code>…`,
  with only the permitted labels, in order, once each. stdout is empty and
  the exit matches JSON. It includes the Phase 23 `issue transition`
  protected-Decision case. (VAL-036)
- [ ] Stop-blocked denials print `blocked by:` with exactly their IDs; other
  denials print none. (VAL-037)
- [ ] Each next-table code prints its row; at least five other codes print
  no `next:`. (VAL-038)
- [ ] Unit fixtures cover scalar flattening, every override, a
  control-character detail and a 201-scalar detail. (VAL-039)
- [ ] `rg` over R9 shows no Human-mode branch that bypasses the renderer
  (the report lists each R9 function and its call).
- [ ] Wall times inside budget. (VAL-052)

## First Failing Test

`issue transition` blocked by a protected Decision prints `error: protected Decision requires an approved reference and Human Outcome` in Human mode. Today it prints a different sentence.

## Documentation Updates

- `project-docs/specs/05_DESIGN_PLAN.md` section 6, and
  `06_DESIGN_PLAN_CLI.md`: the Human failure-line form with examples.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root the Runner records. No command is a placeholder.
Apply gofmt and goimports to edited Go files before the final gate. Every
command below is required before task close. The tagged integration gate is
mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
rg -n "func (write[A-Za-z]*Failure|render[A-Za-z]*Error)\(" cmd/fid -g '!*_test.go'
```

The `-timeout 30m` bounds follow DEC-227/259. Record the wall time `go test`
reports for `cmd/fid/integration`, `cmd/fid` and `internal/workflow`. The
budgets are tagged integration under 900 s and no ordinary package over 600 s
(VAL-052). A breach caused by this task's own new journeys is a bounded speed
repair owned by this task; otherwise stop and escalate. Never record a
timed-out gate as `PASS`.

Required-command rules live in `AGENTS.md` (Commands). Every required command
is reported `PASS` or `BLOCKED` with reason, with its own output token and
observed termination. Focused red/green selections never replace these
commands. Report lint is a Runner control, not a Developer gate.

List any conditional command required by this task and why:

- The `rg` line regenerates R9 for the report's adoption table.

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-24/tasks/task-10-done.md`.
  Commit with `Refs: phase-24/task-10`.
- Red first: record the first failing test before changing product code.
- Retained-test convergence follows DEC-162/165/166. Record each displaced
  assertion with its introducing commit, phase and task. Never infer
  ownership from a filename or task label. Never delete or weaken an
  assertion to pass a gate.
- Any missing registry member, contract ambiguity, or need to change an
  accepted contract shape or assertion is a `BLOCKED` escalation candidate,
  not a local decision (DEC-195).
