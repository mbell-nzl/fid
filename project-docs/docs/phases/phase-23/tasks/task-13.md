# Phase 23 - Task 13: Phase 23 CLI surface, failure codes, and idempotency

## Objective

Complete the real-process CLI and JSON surface for every Phase 23 operation,
land every added failure code with its exact shape, and make every Phase 23
mutation replay its original typed snapshot, including committed denials.

## Validation Contract Coverage

- VAL-055: Every Phase 23 leaf is reachable through the real process with
  canonical JSON and closed request validation, and each of the ten added
  ActivityEvent names is emitted by its producing leaf with no name outside the
  accepted set.
- VAL-056: Every added failure code returns its exact category, exit, subject,
  and closed details from the producing operation named in
  `public-interface.md` section 5.1, and the two precedence probes return the
  single code the accepted order fixes (DEC-226). `finding-unresolved` is a
  ReasonCode, not a failure code, and is not registered in `internal/cli`.
- VAL-057: Every Phase 23 mutation is idempotent and replays its original typed
  snapshot; every Stop-blocked denial commits nothing - no Event, no idempotency
  result, null `EventRef`, `replayed:false` - and returns the same immutable Stop
  and Attention IDs on repeat (DEC-226). `idempotency-conflict` must be
  reproduced through the built binary, never claimed from a fixture inventory.

## Binding Constraints

- `GLOSSARY.md` (DEC-226) distinguishes a **Failure code** (names one refused
  operation; exactly one per failed response; carries category, retryability,
  and exit status) from an **Evaluation reason** (describes why one Obligation
  is unsatisfied; zero or more; no category, retryability, or exit status).
  Register an identifier in the set it belongs to. An identifier may belong to
  both sets only where the same fact both refuses an action and describes a
  state, as `protected-decision-required` does.

- `public-interface.md` sections 1, 4, and 5 are the exact authority. Add no
  leaf, flag, field, or code that is not listed there.
- Every mutation takes `--input`, `--actor`, and `--idempotency-key` and no
  positional argument. Every read uses its declared read grammar. All request
  objects stay closed recursively.
- List leaves keep default limit 100, accepted range 1..10000, and
  `result-limit-exceeded` on a bounded complete projection.
- DEC-216: store and replay the complete typed snapshot for each mutation.
  A committed `review-rework-exhausted` or `protected-decision-required` denial
  replays its original Budget or protected-Decision state, Stop and Attention
  IDs, and EventRef without re-deriving current state. A denial that commits
  nothing has a null EventRef.
- Reuse existing adapters and error rendering; do not create a second response
  construction path or a generic formatter.
- Human text obligations belong to task 14. This task may emit minimal Human
  output, but must not pin fixtures that task 14 will immediately replace.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `cmd/fid/review.go`, `finding.go`, `concern.go`, `deviation.go`,
  `decision.go`, `human_outcome.go`, `requirement.go`, `evidence.go`,
  `envelope.go`, `claim.go`, `root.go`.
- `internal/cli/codes.go`, `internal/cli` response and failure adapters.
- Vocabulary rename (DEC-226): `reasonCodeDefinitions` names the **failure
  code** registry, not the evaluation-reason set, and that collision is what let
  `finding-unresolved` be filed as a failure code. Rename it to
  `failureCodeDefinitions` at its declaration and its three other uses
  (`internal/cli/codes.go`, `internal/cli/envelope.go:38`), plus the test-local
  `reasonCodeExpectation` in `internal/cli/codes_test.go`. Rename nothing in
  `cmd/fid/evidence_test.go`: its `reasonCodes` is the `assertEvaluation`
  parameter naming evaluation reasons, which is the correct sense (I-20). This is a pure
  identifier rename of unexported names: no behavior, no JSON key, and no
  public value changes. `GLOSSARY.md` defines both terms.
- Idempotency snapshot writers in the owning domain packages.
- Matching in-package tests and `cmd/fid/integration/` fixtures.
- `project-docs/docs/phases/phase-23/tasks/task-13-done.md`.

### Files To Read For Context

- `public-interface.md` sections 1, 2, 3, 4, 5; `contracts.md` section 13;
  `schema.md` section 6.3.
- `task-12-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Human text content, resume projections, capability advertisement, and the
  final conformance audit.

## Acceptance Criteria

- [ ] Every Phase 23 leaf returns its declared result through the real process,
  and unknown-field, wrong-nullability, trailing-JSON, and duplicate-key
  negatives reject with no Event. (VAL-055)
- [ ] Every added code matches its accepted category, exit, subject, and closed
  details with no extra key and no reused generic code. (VAL-056)
- [ ] Every mutation replays its original typed result with `replayed:true` and
  no second Event, and both committed denials replay their original Stop,
  Attention, and EventRef. (VAL-057)
- [ ] The failure-code registry and its uses are renamed to name the failure
  set, with no behavior, JSON key, or public value change. (DEC-226)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestPhase23LeavesRejectUnknownRequestFields` fails against the incomplete leaf
surface.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 5: record the full Phase 23
  closed code additions accepted by the slice contract.
- `project-docs/docs/phases/phase-23/tasks/task-13-done.md` in the full format.
## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -timeout 30m ./cmd/fid/integration/...
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

- This task belongs to QA Review window 5.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-13`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.

## Planning Amendment - DEC-249 (2026-09-26)

- In VAL-056's second precedence probe, the unsatisfied Obligation that
  accompanies the open readiness Stop must guard the probed action. An
  Obligation that does not guard it no longer contributes a reason.
- Also read DEC-254 before dispatch. Part 1's `cmd/fid` breach-attribution
  rule applies here, and Part 2 may add a first speed step if the Human
  approves it.
