# Phase 24 - Task 9: Activation descriptors; whole-tree help and template conformance

## Objective

Add descriptors for the eight `activation` leaves (shapes unchanged),
including the discriminated `component` content. Delete the now-empty
pending-descriptor fence. Make the help, template, typed-template,
exclusivity and structural-validation conformance run over all 39 `--input`
leaves. Assemble the shared fake/fid conformance evidence for every changed
leaf.

## Validation Contract Coverage

- VAL-011..016 (completion: all 39 leaves, no fence).
- VAL-035 (assembly across tasks 4-8).
- VAL-052.

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Contract section 8; DEC-276.
- `activation component` `content` is `ActorProfile[]` when
  `component_type` is `authority-profile`, and `GateDefinition[]` when it is
  `gate-definitions`. The descriptor uses discriminated variants. Structural
  errors name the path inside the selected variant. A mismatched variant
  names `input.content[i]`.
- `GateDefinition` literal `environment[].value` and `redact_literals` are
  `sensitive`: their values never appear in failure details.
- Deleting the fence leaves one invariant: registry equals the live
  `--input` set (DEC-192 pattern). No other leaf changes shape.

## Scope

### Files To Create Or Modify

- Descriptor files for the activation leaves (new); `cmd/fid/activation.go`
- `cmd/fid/input_descriptor_test.go` (fence removal; whole-tree corpora)
- Fake/fid conformance assembly in `cmd/fid/integration/` (the Phase 22 fake
  corpus files)

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` section 8
- `project-docs/docs/phases/phase-24/input-inventory.md` sections 1, 2 (R5)
- `project-docs/docs/phases/phase-24/decisions.md` DEC-276
- `project-docs/docs/phases/phase-21/public-interface.md` section 2.2;
  `project-docs/docs/phases/phase-22/public-interface.md` section 3
  (activation inputs)

### Out Of Scope

- Any input shape change. Human failure lines (task 10).

## Acceptance Criteria

- [ ] Reflection conformance covers all 39 leaves. The fence and its test
  are gone, and registry equals the live set. (VAL-011)
- [ ] Help conformance covers all 39 leaves. Full text is pinned for
  `activation component`. (VAL-012)
- [ ] `--example` and typed-template corpora cover all 39 leaves.
  (VAL-013, VAL-014)
- [ ] Exclusivity and structural-validation corpora cover all 39 leaves,
  including both `content` variants and a sensitive-field case with
  `value:null`. (VAL-015, VAL-016)
- [ ] The fake/fid conformance table lists every changed leaf the fake
  implements, with identical acceptance and rejection. (VAL-035)
- [ ] Wall times inside budget. (VAL-052)

## First Failing Test

The live `--input` set equals the descriptor registry with no fence. Today 8 activation leaves remain fenced.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: the activation input help.

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
rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid
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

- The `rg` line regenerates R5. Report whether any family-local decoder
  remains and why (every structural check must go through descriptors).

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-24/tasks/task-9-done.md`.
  Commit with `Refs: phase-24/task-9`.
- Red first: record the first failing test before changing product code.
- Retained-test convergence follows DEC-162/165/166. Record each displaced
  assertion with its introducing commit, phase and task. Never infer
  ownership from a filename or task label. Never delete or weaken an
  assertion to pass a gate.
- Any missing registry member, contract ambiguity, or need to change an
  accepted contract shape or assertion is a `BLOCKED` escalation candidate,
  not a local decision (DEC-195).
