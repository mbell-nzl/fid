# Phase 24 - Task 6: Gate registration, Envelope and Claim leaves

## Objective

Converge the leaves that bind repository identity before a Run exists:
- `gate register` drops `worktree_id` and `repository`;
- `envelope create` drops `repository_scope` and every Source `repository`;
- `envelope complete` drops `repository`;
- `claim acquire` drops `worktree_id`.

Implement the contract section 5 *same worktree* rule once, including
imported-destination equivalence. Descriptors cover all eight gate-register,
envelope and claim leaves.

## Validation Contract Coverage

- VAL-009, VAL-011..016, VAL-017, VAL-018, VAL-019 (`envelope.create`
  Sources), VAL-022, VAL-023 (these leaves).
- VAL-026: `gate.register` binds the observed repository.
- VAL-027: `envelope.create` records the observed scope; plain-copy
  Registration refused.
- VAL-028: `claim.acquire` records the observed worktree; plain copy
  refused.
- VAL-029 (`claim.acquire`, `envelope.create`): imported-destination
  equivalence.
- VAL-035, VAL-052.

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Contract sections 2-5 (*same worktree*) and 7.
- **Same worktree.** One shared predicate: equal `project_id` and
  `worktree_id`, or an imported destination whose import binding names this
  project with `target_root` equal to the workspace root. It generalizes
  `importedDestinationBinding` and `importedGateDestination`. Tasks 7 and 8
  reuse it; do not duplicate it.
- **`gate.register`** records the observed `project_id` and `worktree_id`.
  `binding_digest` is computed over the observed repository plus scope.
- **`envelope.create`** sets `repository_scope` to the observation. Every
  Source (`risk_result.source`, `sources[]`, `obligations[].source`) resolves
  per contract section 3. The existing Registration-worktree check uses the
  predicate.
- **`envelope.complete`** already evaluates against an observation; it only
  drops the input member.
- **`claim.acquire`** records the observed `worktree_id`. The Envelope-scope
  check uses the predicate. The `renew`, `release` and `expire` inputs are
  unchanged.
- Events keep `repository` from fid values only. Claim and Finding Events
  keep the Envelope scope.

## Scope

### Files To Create Or Modify

- `internal/evidence/registration.go`, `internal/workflow/envelope.go`,
  `internal/workflow/import_binding.go` (shared predicate),
  `internal/claims/generic.go`
- `cmd/fid/{gate,envelope,claim}.go` and descriptor files for the eight
  leaves
- Retained tests and helpers (for example `testEnvelopeRepository` and
  `finalizedEnvelopeWorkspaceWithRepository`); integration corpora
  (`p21_*`, `p22_*`, `p23_*`, fake adapters): convergence with provenance

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` sections 2-5
- `project-docs/docs/phases/phase-24/input-inventory.md` sections 1-4 (R6, R7)
- `project-docs/docs/phases/phase-24/decisions.md` DEC-273, DEC-277
- `project-docs/docs/phases/phase-22/contracts.md` section 10 (import)
- `internal/evidence/execution.go` `importedGateDestination`

### Out Of Scope

- The Run, gate-run and Evidence leaves (tasks 7, 8); `ExpectedRepository`
  on these leaves (they take none).
- Parallel worktrees, Envelope revision and naming worktrees by path
  (Phase 25).

## Acceptance Criteria

- [ ] The Registration's `project_id`, `worktree_id` and `binding_digest`
  derive from the independent observation. (VAL-026)
- [ ] For every observing leaf of this task, the stored record repository,
  each Source repository and the Event `repository` equal the independent
  observation. (VAL-022)
- [ ] `envelope.create`: `repository_scope` equals the observation and every
  Source position passes VAL-017/018/019. A Registration from a plain copy is
  refused with no Event. (VAL-027, VAL-017, VAL-018, VAL-019)
- [ ] `claim.acquire` records the observed worktree. A plain copy is refused
  with no Event. (VAL-028)
- [ ] After export/import into a new directory, `claim.acquire` and
  `envelope.create` against imported history are accepted, and refused in a
  plain copy. (VAL-029)
- [ ] An unobservable repository returns `repository-unavailable` for the
  observing leaves. Retries after HEAD advanced replay. (VAL-023, VAL-009)
- [ ] The eight leaves leave the fence; fake/fid agree. (VAL-011..016,
  VAL-035)
- [ ] Wall times inside budget. (VAL-052)

## First Failing Test

`claim acquire` without `worktree_id` succeeds and records the observed worktree. Today it is rejected as missing input.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` and
  `07_DESIGN_PLAN_FLOWS.md`: the gate register, envelope create and claim
  acquire inputs.

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

- None.

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-24/tasks/task-6-done.md`.
  Commit with `Refs: phase-24/task-6`.
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
- Phase 21-23 journeys build these inputs through shared helpers. Converge
  the helper, not each caller.
