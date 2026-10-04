# Phase 24 - Task 12: Concern, Deviation and Human Outcome get leaves

## Objective

Add the read leaves `fid concern get <concern-id>`,
`fid deviation get <deviation-id>` and `fid human-outcome get <outcome-id>`
(DEC-274), so each judgment record has a detail view.

## Validation Contract Coverage

- VAL-042: each `get` returns the same record as its list item.
- VAL-043: unknown IDs and bad grammar are rejected.
- VAL-052.

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Contract section 10 (detail views) and DEC-274. The operations are
  `concern.get`, `deviation.get` and `human-outcome.get`. Each takes exactly
  one positional ID and only the global flags. Results are
  `{concern:Concern}`, `{deviation:Deviation}` and
  `{outcome:HumanOutcome}`, with `event:null`.
- Not-found reuses `concern-not-found`, `deviation-not-found` and
  `outcome-not-found`: exit 3, subject ID = the requested ID, `{id}`. A
  grammar error is `invalid-argument`, exit 2, with the task 9 grammar
  details (DEC-278) and Human lines through the task 10 renderer.
- The Human view uses the existing item rendering in this task. Task 13
  adds full reasoning to detail views.
- No capability row changes. These are read operation IDs only, outside the
  activated-action allowlist (Phase 23 section 4.1 treatment).

## Scope

### Files To Create Or Modify

- `internal/review/{concern_list,deviation,human_outcome}.go` (read
  functions)
- `cmd/fid/{concern,deviation,human_outcome}.go`
- `cmd/fid/contract_surface_test.go` (live leaf inventory, R14); code
  producer conformance

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` section 10
- `project-docs/docs/phases/phase-24/input-inventory.md` R14
- `project-docs/docs/phases/phase-24/decisions.md` DEC-274
- `cmd/fid/review.go` (`review get` pattern); `cmd/fid/requirement.go`
  (`requirement get`)

### Out Of Scope

- Reasoning display and disclaimers (task 13).

## Acceptance Criteria

- [ ] Each `get` returns the record equal to its list item, with
  `event:null` and exit 0. (VAL-042)
- [ ] Unknown ID returns the matching `*-not-found` (exit 3, `{id}`). A
  missing ID, an extra argument or an unknown flag returns
  `invalid-argument` (exit 2). (VAL-043)
- [ ] The live leaf inventory includes the three leaves, each with a
  one-line description (VAL-010 stays green). The whole-tree grammar corpus
  and the pairing corpus cover them with no new exclusion (VAL-064 and
  VAL-036 stay green).
- [ ] Wall times inside budget. (VAL-052)

## First Failing Test

`fid --json concern get CON-1` returns `{"concern":{...}}`. Today the leaf does not exist.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: the three leaves.
- Phase 24 `public-interface.md` needs no change, because DEC-274 already
  states them.

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

- None.

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-24/tasks/task-12-done.md`.
  Commit with `Refs: phase-24/task-12`.
- Red first: record the first failing test before changing product code.
- Retained-test convergence follows DEC-162/165/166. Record each displaced
  assertion with its introducing commit, phase and task. Never infer
  ownership from a filename or task label. Never delete or weaken an
  assertion to pass a gate.
- Any missing registry member, contract ambiguity, or need to change an
  accepted contract shape or assertion is a `BLOCKED` escalation candidate,
  not a local decision (DEC-195).
