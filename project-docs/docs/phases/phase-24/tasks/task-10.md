# Phase 24 - Task 10: Human failure lines in every family

## Objective

Replace the Phase 21 "error plus two guidance lines" rule with contract
section 9. Every Human-mode failure prints `error:` with the exact JSON
message, then `reason:`, then `blocked by:` and `next:` where they apply.
Every line derives from the JSON failure. One shared renderer is adopted by
every family.

The Human decided D6 = A (DEC-281), which resolves its decision prerequisite.
Dispatch follows QA acceptance and Runner readiness controls in the phase brief.
Each generic `invalid-argument` fallback now carries its domain cause as the
JSON message, so the `error:` line still shows why the request was refused.

## Validation Contract Coverage

- VAL-036: per code, the Human lines repeat the JSON failure.
- VAL-037: `blocked by:` appears exactly when details carry blocking IDs.
- VAL-038: `next:` appears exactly for next-table codes.
- VAL-039: `reason:` key-detail rendering, escaping and the 200-scalar cut.
- VAL-067: each generic fallback's message is the domain cause it caught
  (DEC-281).
- VAL-027, VAL-028, VAL-029 (message half): the plain-copy refusals carry
  their domain reasons.
- VAL-052.
- Advances VAL-054 (Human walk-through at the task-11 window).

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Contract section 9 is exact. It replaces Phase 21 `public-interface.md`
  section 1's guidance-line rule. JSON failures, exit codes and stdout are
  unchanged, except the generic-fallback messages below.
- **Generic fallbacks (DEC-281).** Each of the 18 generic fallbacks
  (inventory R20) sets `message` to the caught error's text,
  `cause.Error()`, the text Human mode printed after `error: ` at
  `62dddab`.
  - Reuse task 9a's shared message construction (contract section 9,
    DEC-282): escape, then bound to 512 scalars. Text over 512 keeps the
    first 511 plus `…`, dropping any escape crossing the cut whole. An
    empty or all-whitespace cause keeps the fixed message.
  - Code, category, retryable, subject, details and exit stay unchanged.
  - A cause never carries a `sensitive` descriptor value.
  - The plain-copy messages become `prior registration binding does not
    match envelope`, `claim acquire binding is invalid`, `run claim
    binding is invalid` and `run gate: repository binding is invalid`
    (VAL-027..029; DEC-282 corrects the Run cause to the earlier Claim check).
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
  9. Control characters are escaped. The cut follows contract section 9 as
  clarified by DEC-278: count the whole line, label included, after
  escaping. Over 200 scalars, keep the first 200 and add `…`. Exactly 200
  stays uncut. Never split an escape sequence.
- **`blocked by:`** lists `stop_ids` then `attention_ids` when either is
  nonempty.
- **`next:`** uses the closed table in contract section 9 with values from
  details. `<leaf>` is the CLI path (for example `fid run close`). The
  Phase 23 task 21 rework guidance line is replaced by its table row. The
  `invalid-argument` row applies to every `invalid-argument` on an `--input`
  leaf, grammar failures and `--example` exclusivity failures included
  (contract section 9; DEC-280 item 1). Failures on leaves without `--input`
  print no `next:`, and neither does `operation-unavailable`.
- **Grammar failures and unknown commands** (task 9a) render through the
  same renderer: the `error:` payload is their already escaped and bounded
  JSON `message` (DEC-282), verbatim. Full offending spelling/value stays
  in details even if the message is cut. Do not escape or truncate the
  final message again.
- **Escaping (DEC-280 item 6).** The renderer uses the `internal/cli` escape
  function from task 9a for the `reason:` and `blocked by:` values. The
  `error:` line is the message verbatim, and a message never holds a
  control character.
- **`gate run`.** Pre-frame failures use these lines. A terminal non-success
  in Human mode prints them on stderr after the streamed output. The
  streamed output stays on stdout, as Phase 21 section 1's streaming
  exception allows (DEC-280 item 7).
- `fid journal inspect` and `fid sql` are unstable Human-only diagnostics
  outside `fid-cli.v1` and are excluded.

## Scope

### Files To Create Or Modify

- `internal/cli/` Human failure renderer (new), with tests
- Every function listed by inventory R9 (41 in `cmd/fid`)
- Retained Human stderr assertions (R10: 32 in 13 files): convergence with
  provenance
- The 18 generic fallbacks (inventory R20). Retained fallback-message
  assertions converge with provenance: the plain-copy message assertions
  from tasks 6-8, and `cmd/fid/issue_update_test.go` and
  `cmd/fid/generic_run_test.go` (one line each at `62dddab`)
- A built-binary generic-fallback corpus (new), one row per R20 site, with
  a pinned unreachability list
- A built-binary Human/JSON pairing corpus (new), enumerated from
  `cli.RegisteredCodes()` with a pinned exclusion list

### DEC-285 Amendment (recovery attempt 2)

- Three stable callers outside the R9 functions adopt the renderer.
  TRG-24-009 added `cmd/fid/activity.go`; DEC-285 clause 3 adds
  `cmd/fid/stats.go` and `cmd/fid/issue_ac.go`. The ordinary race/coverage
  gate, `golangci-lint run ./...` and the tagged integration gate below
  already execute all three files.
- `activity.list` reads JSON mode first. It sends the unwrapped
  workspace-resolution error and every `activity.List` error to
  `renderIssueReadError(cmd, jsonMode, "activity.list", err)`. That is the
  retained-read mapper its Phase 21 section 4.1 siblings use. The exact
  expected failures are in the DEC-285 clause 2 table: `--since`/`--on`
  give `invalid-argument`/2 with the domain cause as message; no project
  gives `project-not-found`/3; replay stopped gives
  `journal-replay-stopped`/6. Today's raw exit 7 is a Phase 21
  nonconformance, not an accepted contract.
- In `stats.go` and `issue_ac.go`, the six forced `true` arguments become the
  live JSON mode. JSON failures stay byte-identical. Human success output
  is unchanged (stats prints JSON there today).
- Do not edit `internal/activity`, `issueReadFailure`, `genericFailure` or
  the flag-domain declarations. Retained assertions these callers displace
  converge with provenance (DEC-162/165/166) in the test file that holds
  them.
- Recovery identity: `task-10:public-interface:activity-list-failure-envelope-mapping`
  at attempt 2 (DEC-285 clause 7).

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` section 9
- `project-docs/docs/phases/phase-24/input-inventory.md` R8, R9, R10
- `project-docs/docs/phases/phase-24/decisions.md` DEC-272, DEC-278,
  DEC-280, DEC-281 (D6) and DEC-282
- `project-docs/docs/phases/phase-21/public-interface.md` sections 1, 7;
  `project-docs/specs/05_DESIGN_PLAN.md` section 6
- `cmd/fid/stop.go` `writeStopFailure` (the task 19/21 guidance)

### Out Of Scope

- Stop, Attention and repository views (task 11). Judgment views (task 13).
- Changing any failure code or detail. No message changes, other than the
  generic fallbacks (DEC-281) and the grammar messages that task 9a already
  set.

## Acceptance Criteria

- [ ] The pairing corpus has one row per code in `cli.RegisteredCodes()`. A
  pinned exclusion list in the test gives each excluded code its reason (for
  example, no reachable CLI producer). A code with neither a row nor an
  exclusion fails the test. Rows for grammar failures (task 9a), unknown
  commands (including the near miss `fid isue`) and a post-frame `gate run`
  failure (`gate-unsuccessful`) are kept apart from the `--input`
  structural rows. In every row, Human stderr is `error: <JSON message>`
  then `reason: <code>…`, with only the permitted labels, in order, once
  each. No message contains a control character, and the exit matches JSON.
  stdout is empty in every row except the post-frame `gate run` row: there
  stdout carries the streamed Human output, and the JSON side is the
  failure inside the `--json` run's terminal frame. The corpus includes
  the Phase 23 `issue transition` protected-Decision case. (VAL-036)
- [ ] The pairing corpus includes task 9a's short exact grammar and
  unknown-command messages, the long flag and long word, the 512/513
  escaped-text boundaries, non-ASCII scalars and an escape crossing the
  cut. Human `error:` matches the bounded JSON message byte for byte;
  details retain the full offending values and metadata/exit are unchanged.
  (VAL-036; rechecks task 9a's VAL-064/065 construction)
- [ ] Stop-blocked denials print `blocked by:` with exactly their IDs; other
  denials print none. (VAL-037)
- [ ] Each next-table code prints its row. So does a grammar failure on an
  `--input` leaf (`fid stop raise --bogus` prints
  `next: see fid stop raise --help; template: fid stop raise --example`).
  At least five other codes print no `next:`, and neither does
  `invalid-argument` from a leaf without `--input` (for example a grammar
  failure on `fid issue create`). (VAL-038)
- [ ] Unit fixtures cover scalar flattening, every override, a
  control-character detail, a `reason:` line of exactly 200 scalars (uncut)
  and one of 201 (first 200 plus `…`), and an escape sequence at the cut
  (dropped whole). (VAL-039)
- [ ] `rg` over R9 shows no Human-mode branch that bypasses the renderer
  (the report lists each R9 function and its call, plus the `activity.go`,
  `stats.go` and `issue_ac.go` callers).
- [ ] Built-binary rows, both modes (DEC-285 clause 6): `activity list
  --since invalid`, `activity list --on invalid` and `activity list` with no
  project, each matching the DEC-285 clause 2 table; and `stats tokens`,
  `stats cost` and `issue ac list <id>` with no project. In every row
  Human stderr is the section 9 lines for the JSON failure, the exits
  match, stdout is empty and nothing in the filesystem changes. (VAL-036,
  VAL-067)
- [ ] The generic-fallback corpus has one row per R20 site, driven by a
  reachable domain cause in both modes. A fallback with neither a row nor a
  pinned unreachability reason fails the test.
  - JSON `message` is the escaped and bounded cause text, and the Human
    `error:` line is the same text.
  - Code, details and exit are unchanged.
  - Unit rows cover escaped causes of 512 scalars (uncut) and 513 scalars
    (first 511 plus `…`), non-ASCII scalar counting, a cause holding ESC,
    an escape crossing the cut (dropped whole), and empty and all-whitespace
    causes (the fixed message stays).
  - A guard row: an `activation component` domain failure whose
    `gate-definitions` input holds a sentinel `environment[].value`. The
    sentinel appears in neither stream. Prove the row reaches the generic
    fallback, rather than an earlier structural or typed failure (QA plan
    section 12.6.2 records a reachable duplicate-component case).

  (VAL-067)
- [ ] The plain-copy refusals of `envelope.create`, `claim.acquire`,
  `run.start` and `gate.run` carry the four DEC-281/282 messages, with code
  and details unchanged. In both modes the plain-copy `run.start` row says
  `run claim binding is invalid`; preserve task 7's valid Claim and its
  original/copy/import controls. (VAL-027, VAL-028, VAL-029)
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
rg -n 'CodeInvalidArgument, "invalid [^"]*"' cmd/fid -g '!*_test.go'
```

The `-timeout 30m` bounds follow DEC-227/259. Record the wall time `go test`
reports for `cmd/fid/integration`, `cmd/fid`, `internal/workflow` and the
slowest other ordinary package. The budgets are tagged integration under
900 s and no ordinary package over 600 s (VAL-052); report any ordinary
package over the 480 s warning. A breach caused by this task's own new journeys is a bounded speed
repair owned by this task; otherwise stop and escalate. Never record a
timed-out gate as `PASS`.

Required-command rules live in `AGENTS.md` (Commands). Every required command
is reported `PASS` or `BLOCKED` with reason, with its own output token and
observed termination. Focused red/green selections never replace these
commands. Report lint is a Runner control, not a Developer gate.

List any conditional command required by this task and why:

- The first `rg` line regenerates R9 for the report's adoption table. Family
  writers remain as callers of the shared renderer, so it is expected to
  terminate `exit_code: 0`.
- The second `rg` line regenerates R20. The fixed fallback messages remain
  for empty causes, so it is expected to terminate `exit_code: 0`. The
  report maps each generic-fallback match to its corpus row or its
  unreachability reason.

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
