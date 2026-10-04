# Phase 24 - Task 8: Gate run and Evidence leaves; real freshness

## Objective

Converge execution and Evidence:
- `gate run` replaces `repository` with a required `expected_repository`
  and records one pre-spawn observation as Attempt, Receipt and Event
  repository;
- `evidence link` derives `target_repository` from its subject;
- `evidence evaluate` evaluates against the current observation;
- `artifact verify` gets its descriptor (four leaves leave the fence).

Freshness becomes a comparison of real observations.

## Validation Contract Coverage

- VAL-009, VAL-066, VAL-022 (`gate.run`, the one observing mutation of this
  task).
- VAL-023 (`gate.run` and `evidence.evaluate`, the observing leaves).
- VAL-011..016, VAL-060 (descriptors for all four leaves), VAL-017
  (`gate.run`, `evidence.link`, `evidence.evaluate`).
- VAL-025 (`gate.run` producer): the amended mismatch details.
- VAL-024 (`gate.run`): expectation check, with no Attempt or process on
  mismatch.
- VAL-029 (`gate.run`): imported-destination equivalence.
- VAL-032: one pre-spawn observation everywhere.
- VAL-033: real freshness in both directions.
- VAL-034: Evidence targets come from the subject.
- VAL-035, VAL-052.

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Contract sections 4, 5 and 7; DEC-273.
- **`gate.run`**
  - Closed input, then replay, then the existing checks. The Registration
    and Run repository comparisons use the task 6 *same worktree*
    predicate, which subsumes `importedGateDestination`; never stored state.
  - The expectation check runs before any Attempt append. One observation
    is taken before spawn and recorded as Attempt `start_repository`,
    Receipt `terminal_repository`, and the `repository` of
    `gate.execution-started`, `gate.receipt-recorded` and any committed
    `run.execution-blocked` Event.
  - There is no re-observation after the process exits. The `fid.live.v1`
    contract is unchanged.
- **`evidence.link`** does not observe (contract section 2, DEC-278).
  `target_repository` and the `evidence.linked` Event `repository` are the
  Receipt's `terminal_repository`, the Review's `repository`, or the
  repository of the `human-outcome.recorded` Event.
- **`evidence.evaluate`** input is `{obligation_id:?ID}`. It is the one
  observing read: it observes once, after closed-input validation, and an
  unobservable worktree returns `repository-unavailable` with the mapped
  reason (DEC-278). Freshness compares each Evidence target with the current
  observation, exactly as terminal transition and completion do. Human
  Outcome Evidence stays `fresh` (Phase 23).
- **Plain-copy refusal (DEC-278).** `gate.run` from a different worktree
  without a matching import binding returns the existing fallback failure:
  `invalid-argument`, exit 2, message `invalid gate request`, details
  `{argument:"input",value:null,allowed:["accepted gate registration input"]}`,
  and no Attempt, process or Event. The plain-copy fixture is defined in
  task 6.
- **Real freshness.** Phase 21-23 journeys that relied on copied blocks may
  now see `evidence-stale` after editing tracked files. A displaced
  assertion is convergence of this accepted change: fix the fixture so it
  reflects the intended state, recording provenance. Never weaken the
  freshness rule. Artifact paths must stay git-ignored in fixtures
  (inventory section 3).

## Scope

### Files To Create Or Modify

- `internal/evidence/{execution,evidence,registration}.go`,
  `internal/workflow/evaluation.go` (only if the evaluate path needs the
  shared observation)
- `cmd/fid/{gate,evidence,artifact}.go` and descriptor files for the four
  leaves
- Retained tests and integration journeys (walking skeleton, recovery,
  freshness supersession, fake adapters): convergence with provenance

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` sections 4, 5
- `project-docs/docs/phases/phase-24/input-inventory.md` sections 1-4
- `project-docs/docs/phases/phase-24/decisions.md` DEC-273
- `project-docs/docs/phases/phase-21/contracts.md` sections 8-10;
  `project-docs/docs/phases/phase-23/contracts.md` sections 7, 13

### Out Of Scope

- Terminal re-observation after gate exit; changes to live frames.
- Human rendering (tasks 10, 11).

## Acceptance Criteria

- [ ] Start a Run, edit files, then `gate.run` with a matching expectation:
  admitted. The Attempt, Receipt and both Events carry one identical
  repository equal to the independent pre-spawn observation. (VAL-032,
  VAL-022)
- [ ] HEAD-differs and dirty-differs return `repository-mismatch` with the
  amended details and no Attempt, process or Event. (VAL-024, VAL-025)
- [ ] A Receipt-backed Obligation is fresh with the tree unchanged, stale
  (`evidence-stale`) in `evidence.evaluate`, the terminal-transition denial
  and the `envelope.complete` denial after a tracked edit, and fresh again
  after reverting. (VAL-033)
- [ ] Receipt, Review and Human Outcome Evidence each record the subject's
  repository as `target_repository` and as the Event `repository`.
  (VAL-034)
- [ ] `gate.run` in an imported destination is admitted. The same request in
  a plain copy returns the exact refusal above. (VAL-029)
- [ ] Removed members are rejected for `gate.run`, `evidence.link` and
  `evidence.evaluate`. (VAL-017)
- [ ] `gate.run` and `evidence.evaluate` with an unobservable repository
  return `repository-unavailable` with the mapped reason (and, for
  `gate.run`, no Attempt, process or Event). (VAL-023)
- [ ] `gate.run`: a retry after HEAD advanced replays (VAL-009). A retry in an
  unobservable worktree replays with no Event, and a new-key control fails
  with `repository-unavailable` (VAL-066).
- [ ] The four leaves leave the fence with VAL-011..016 and VAL-060 rows;
  fake/fid agree over the derived set. (VAL-011..016, VAL-060, VAL-035)
- [ ] Wall times inside budget. (VAL-052)

## First Failing Test

`gate run` after editing a tracked file since Run start, with a matching expectation, records the edited tree's dirty digest on its Receipt. Today the Receipt repeats the caller's block.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` and `07_DESIGN_PLAN_FLOWS.md`:
  the gate run and Evidence inputs, and freshness from observations.

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

- Report path: `project-docs/docs/phases/phase-24/tasks/task-8-done.md`.
  Commit with `Refs: phase-24/task-8`.
- Red first: record the first failing test before changing product code.
- Retained-test convergence follows DEC-162/165/166. Record each displaced
  assertion with its introducing commit, phase and task. Never infer
  ownership from a filename or task label. Never delete or weaken an
  assertion to pass a gate.
- Any missing registry member, contract ambiguity, or need to change an
  accepted contract shape or assertion is a `BLOCKED` escalation candidate,
  not a local decision (DEC-195).
- Remove this task's leaves from the pending-descriptor fence (task 3). Each
  adopted leaf gets a descriptor, help, `--example`, typed template,
  exclusivity and structural-validation corpus rows (VAL-011..016).
- This is the highest convergence-risk task. If a retained journey's intent
  cannot be preserved without changing the freshness rule, stop and
  escalate.
