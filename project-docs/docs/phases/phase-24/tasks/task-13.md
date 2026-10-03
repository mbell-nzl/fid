# Phase 24 - Task 13: Judgment disclaimers once per section; recorded reasoning shown

## Objective

Apply A6 and A7 (contract section 10):
- each judgment section prints its disclaimer once;
- Human views show recorded reasoning, cut to one line in lists and in full
  in detail views, with control characters escaped;
- a sentinel guard proving that referenced project-file content never
  appears replaces the `PRIVATE_` leak checks;
- the glossary gains the distinction between recorded reasoning and
  project-owned prose.

## Validation Contract Coverage

- VAL-044: one disclaimer per section (restates Phase 23 VAL-060, proof
  half).
- VAL-045: recorded reasoning, cut in lists and full in detail views.
- VAL-046: control characters are escaped.
- VAL-047: no referenced project-file content (restates Phase 23 VAL-060,
  prose half).
- VAL-048: glossary terms.
- VAL-052.
- Advances VAL-056 (restated Phase 23 VAL-073, Human, phase-end).

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Contract sections 10 and 11. The closed recorded-reasoning fields are:
  - Concern `summary`, `reasoning` and disposition `reason`;
  - Finding `summary` and every FindingDisposition `reason`;
  - Deviation `expected`, `observed` and `reasoning`;
  - Human Outcome `reasoning`.
- **List views** (every `list` leaf and `resume` sections): one line,
  whitespace collapsed, cut to 80 Unicode scalars with `…`; 80 or fewer are
  uncut.
- **Detail views** (the task 12 `get` leaves, mutation outputs and
  `review get`): full text, with line breaks kept and indented.
- **Disclaimer:** exactly one per section that lists Human Outcomes,
  Concerns or Deviations, before the items. None on items. An empty section
  prints `none` and no disclaimer. The `Purpose:` lines stay. The on/off
  switch is Phase 26.
- Control characters in recorded reasoning are escaped (PRD-OPS-007).
- **Sentinel guard:** fixtures seed referenced project files with sentinel
  content. A test scans every Human output of records that reference them.
  The `PRIVATE_` checks (`p23_human_leaf_pairs_test.go:123`,
  `p23_human_fixtures_test.go:143`) are removed.

## Scope

### Files To Create Or Modify

- `cmd/fid/p23_human.go` (split if it exceeds 400 lines), plus the judgment
  sections of `resume_human.go`
- `cmd/fid/p23_human_test.go`;
  `cmd/fid/integration/p23_human_leaf_pairs_test.go`,
  `p23_human_fixtures_test.go` (guard replacement, with provenance)
- `project-docs/specs/GLOSSARY.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` sections 10, 11
- `project-docs/docs/phases/phase-24/input-inventory.md` R13
- `project-docs/docs/phases/phase-24/decisions.md` DEC-272, DEC-274
- `project-docs/docs/phases/phase-23/contracts.md` section 14;
  `project-docs/docs/phases/phase-23/public-interface.md` section 7

### Out Of Scope

- The disclaimer on/off switch (Phase 26). Layout redesign (CF-040).
- Stop and repository lines (task 11).

## Acceptance Criteria

- [ ] Every judgment view with zero, one and three items: one disclaimer
  when items exist, none on items, and `none` when empty. (VAL-044)
- [ ] Every closed field with 80-scalar, 81-scalar and multi-line text
  renders cut in lists and in full in detail views. (VAL-045)
- [ ] Reasoning containing ESC, CR, BEL and NUL is escaped in every view
  that shows it. (VAL-046)
- [ ] The sentinel guard passes over every Human fixture, and the `PRIVATE_`
  checks are gone. (VAL-047)
- [ ] `GLOSSARY.md` defines **Recorded reasoning** and **Project-owned
  prose** consistently with contract section 10. (VAL-048)
- [ ] Wall times inside budget. (VAL-052)

## First Failing Test

`fid concern list` prints a `Summary:` line cut to 80 scalars with `…`. Today no summary is printed.

## Documentation Updates

- `project-docs/specs/GLOSSARY.md` (VAL-048).
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: judgment view examples.

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
rg -n "This (Concern|Deviation) is recorded judgment|This Outcome is authorized Human judgment" cmd/fid -g '!*_test.go'
rg -n "PRIVATE_" cmd/fid/integration
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

- The two `rg` lines prove the per-item disclaimers and the `PRIVATE_` checks are removed.

## Expected Non-Zero Commands

- `rg -n "This (Concern|Deviation) is recorded judgment|This Outcome is authorized Human judgment" cmd/fid -g '!*_test.go'` - expected `exit_code: 1` - no per-item disclaimer remains
- `rg -n "PRIVATE_" cmd/fid/integration` - expected `exit_code: 1` - the old leak checks are replaced by the sentinel guard

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-24/tasks/task-13-done.md`.
  Commit with `Refs: phase-24/task-13`.
- Red first: record the first failing test before changing product code.
- Retained-test convergence follows DEC-162/165/166. Record each displaced
  assertion with its introducing commit, phase and task. Never infer
  ownership from a filename or task label. Never delete or weaken an
  assertion to pass a gate.
- Any missing registry member, contract ambiguity, or need to change an
  accepted contract shape or assertion is a `BLOCKED` escalation candidate,
  not a local decision (DEC-195).
