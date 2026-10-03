# Phase 24 - Task 11: Readable Stop, Attention and repository lines

## Objective

Render Stops and Attention as labelled lines, with a `Next:` command keyed on
each Stop's resolution Condition, in `stop list`, `attention list` and
`resume`, and never print raw JSON. Shorten every Human repository line to a
12-digit HEAD and `clean`/`dirty`.

## Validation Contract Coverage

- VAL-040: labelled Stop and Attention rendering with `Next:` per Condition.
- VAL-041: short Human repository lines.
- VAL-052.
- Advances VAL-054 and VAL-055 (Human walk-through at this task's window).

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Contract section 10 (A4, repository lines). JSON is unchanged.
- **Each Stop** shows subject, code, blocked actions, one labelled line per
  StopDetails member of its closed variant (10 codes), resolution
  condition, resolver capability, status (open, or resolved with its
  Event), and the `Next:` line from the contract section 10 Condition table.
  An open Stop with no lifecycle disposition also names `fid stop resolve`,
  then `fid attention resolve`.
- **Attention** shows its status, code and linked Stops, rendered as above.
- **Repository lines** in every Human view (inventory R11) show
  `HEAD <12 hex>` with `clean` or `dirty`. Recorded-versus-current
  comparisons say `same worktree` or `different worktree`, using the task 6
  predicate. No Human view prints `worktree_id`, `tree_oid`, `dirty_digest`,
  `capture_source` or a nil marker.
- Labelled sections and `Purpose:` lines stay (DEC-219 pattern).

## Scope

### Files To Create Or Modify

- `cmd/fid/{stop_human,resume_human,resume_human_guidance,run_inspect_human,p23_human}.go`
  (repository lines only in `p23_human.go`), plus every other view R11 finds
- Human fixtures in ordinary and tagged suites: convergence with provenance

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` section 10
- `project-docs/docs/phases/phase-24/input-inventory.md` R11, R12
- `project-docs/docs/phases/phase-24/decisions.md` DEC-272
- `project-docs/docs/phases/phase-22/public-interface.md` section 2
  (StopDetails); `project-docs/docs/phases/phase-23/public-interface.md`
  section 2 (three added variants)

### Out Of Scope

- Judgment disclaimers and reasoning (task 13). Failure lines (task 10).
- Layout and spacing redesign (CF-040).

## Acceptance Criteria

- [ ] `stop list`, `attention list` and `resume` with one Stop of each of
  the 10 codes print every StopDetails member as a labelled line, plus the
  correct `Next:` line, and contain no raw-JSON `{`. (VAL-040)
- [ ] Clean and dirty fixtures across every R11 view print only the short
  repository form, and the forbidden tokens are absent. (VAL-041)
- [ ] Wall times inside budget. (VAL-052)

## First Failing Test

`stop list` with a `budget-exhausted` Stop prints `Budget: <name>` as a labelled line. Today it prints the details as JSON.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: Stop/Attention Human examples.

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
rg -n -i "worktree|dirty digest|HeadOID|head_oid" cmd/fid -g '!*_test.go'
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

- The `rg` line regenerates R11. The report accounts for every remaining match (JSON or input code only).

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-24/tasks/task-11-done.md`.
  Commit with `Refs: phase-24/task-11`.
- Red first: record the first failing test before changing product code.
- Retained-test convergence follows DEC-162/165/166. Record each displaced
  assertion with its introducing commit, phase and task. Never infer
  ownership from a filename or task label. Never delete or weaken an
  assertion to pass a gate.
- Any missing registry member, contract ambiguity, or need to change an
  accepted contract shape or assertion is a `BLOCKED` escalation candidate,
  not a local decision (DEC-195).
- The QA window after this task includes a short Human walk-through of the
  new failure and Stop output (VAL-054, VAL-055).
