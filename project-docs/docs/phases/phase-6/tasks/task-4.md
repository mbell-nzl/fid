# Phase 6 - Task 4: `fid ready` CLI command + formats + empty state

## Objective

Wire the Task 3 selection engine to a new `fid ready` Cobra command: the frozen
single-task JSON shape (design §6.3), the `--epic` scope flag, `--format`
dispatch (json/toon/markdown/yaml), and the empty-state contract (exit 0,
`{"task": null}`, explanatory stderr `note:`). This is the agent-runnable exit
for the ready queue.

## Requirements

- PRD refs: PRD-ISSUE-008, 009, 010, 011, 012; PRD-OUT-001–004 (formats)
- TRD refs: TRD §4.3 (clean read path)
- Design refs: `05_DESIGN_PLAN.md` §6.3; `06_DESIGN_PLAN_CLI.md` §`fid ready`
  (flags, empty-state stderr note)
- Quality refs: `08_QUALITY_PLAN.md` STRUCT-001/002

## Validation Contract Coverage

- VAL-001: bare `fid ready` returns one actionable task in the frozen §6.3 JSON
  shape.
- VAL-002..006: CLI surface of the engine selection cases (sequential,
  prioritised, unattached, project-wide).
- VAL-007: empty-state contract — exit 0, `{"task": null}`, stderr `note:`;
  non-empty returns the task with no note.
- VAL-008: `--format json|toon|markdown|yaml` render selected + empty states;
  JSON is default.

## Scope

### Files To Create Or Modify

- `cmd/fid/ready.go` — new `fid ready [--epic <id>] [--format ...]` command
  calling the Task 3 engine and dispatching through the existing format engine.
- `cmd/fid/root.go` — register the `ready` command.
- `cmd/fid/ready_test.go` — CLI tests: §6.3 JSON shape, `--epic` scope, all four
  formats, empty-state (exit 0 + `{"task": null}` + stderr note), default=JSON.

### Files To Read For Context

- `cmd/fid/issue.go` (existing read-command + `--format` dispatch patterns)
- `cmd/fid/issue_list_test.go` (format/CLI test patterns)
- `internal/issues/ready.go` (Task 3 engine API)
- `internal/format/` (format engine entry point)

### Out Of Scope

- `--explain`, `--limit`, `--include-blocked` (v0.2, PRD-ISSUE-024) — do not add
  these flags this phase.
- Any change to selection logic (lives in Task 3).

## Acceptance Criteria

- [ ] `fid ready` returns the selected task in the frozen §6.3 JSON shape
  (VAL-001).
- [ ] `fid ready --epic <id>` scopes selection to that epic for sequential and
  prioritised cases (VAL-002, 003, 004).
- [ ] Project-wide `fid ready` returns the single top task across prioritised +
  unattached, sequential contributing only its current task (VAL-005, 006).
- [ ] Empty result: exit 0, stdout `{"task": null}`, stderr `note:` explaining
  why (e.g. all completed/cancelled, or current sequential task in flight)
  (VAL-007).
- [ ] Non-empty result emits no stderr `note:` (VAL-007).
- [ ] `--format json|toon|markdown|yaml` render selected and empty states;
  omitting `--format` yields JSON (VAL-008).

## TDD Plan

- First failing test: `fid ready --format json` over a ready fixture asserts the
  §6.3 shape and the expected task id.
- Additional tests: empty fixture → exit 0 + `{"task": null}` + stderr note;
  `--epic` scoping; each non-JSON format renders; default omitted format = JSON.
- Regression risks: stderr vs stdout separation (note must be stderr, payload
  stdout); empty-state exit code drift (must be 0); markdown/toon/yaml rendering
  of a null task.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` §`fid ready` — confirm/illustrate
  the v0.1b default + empty-state behavior (exit 0, `{"task": null}`, stderr
  note) so the frozen contract is explicit. Keep `--explain` clearly marked
  v0.2.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None directly (smoke is Task 5).
- Cross-surface parity impact: the bare `fid ready` JSON shape is frozen
  (Workflow Runner depends on it); JSON/TOON/markdown/YAML must all represent
  the same selected task / empty state.
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Hold COV-002 (`cmd/fid` ≥ 70%).

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
```

- `go build -o /tmp/fid ./cmd/fid`: build the binary to confirm the new command
  registers and runs for manual/CLI evidence.

## Notes For Developer

- The empty-state shape (`{"task": null}`) and exit code 0 are fixed by DEC-018
  (validation-contract scoping note). The Workflow Runner depends on the bare
  shape — do not add fields or change the key.
- Render the `note:` explanation on stderr only; the JSON/markdown/yaml/toon
  payload (including the null task) goes to stdout.
- Reuse the existing `--format` dispatch helper rather than hand-rolling format
  switches.
