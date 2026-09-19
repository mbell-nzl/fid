# Phase 23 - Task 18: Conformance assembly and wire-surface audit

## Objective

Assemble the shared fake and fid conformance corpus for every Phase 23
operation, audit the thirty-name grant action set, and prove no cancelled
product family reappeared.

## Validation Contract Coverage

- VAL-067: The activated capability and grant action set is exactly thirty
  names with allow and reject cases.
  The audit covers **both** copies of the closed set and compares them against
  each other: `sortedUniqueCapabilities` in `internal/workflow/activation.go`
  and the `const actions` string in
  `internal/journal/current_event_value_relations.go:63`. Nothing cross-checks
  them today, and each was widened name-by-name across tasks 4, 5, 8, 9, and 10
  (R-19, DEC-233); a single missed name leaves a capability that a profile can
  hold but the journal rejects, or the reverse. Assert set equality between the
  two directly, not each against a literal, so the pair cannot drift later.
- VAL-068: Shared fake and fid adapters pass identical conformance cases for
  every Phase 23 operation.
- VAL-069: No product escalation, QA-plan, validation-contract, template,
  retrospective, or carry-forward table or command exists. This task owns the
  row; `v1-convergence-inventory.md` section 1 is corrected to match (I-10).
- VAL-070, VAL-071, VAL-072, VAL-073: regenerate task 14's Human fixture
  manifest against the final binary and confirm every recorded invocation still
  runs, so the Human judgment at phase end is exercised against the shipped
  surface rather than a mid-phase one (R-10).
- VAL-079: the phase-end wall-time budgets hold with the explicit `-timeout`,
  and every completion report carries its three recorded package times
  (R-13, DEC-227).

## Binding Constraints

- `public-interface.md` sections 4, 4.1, and 5 are the exact authority. The set
  is exactly thirty names; an action added before its owning task or missing at
  this audit both fail.
- The audit must compare the explicit name set with an independently pinned
  corpus. Parsing one table to produce both expected and actual values cannot
  pass (the Phase 22 rule, retained).
- Conformance cases exist for every new and amended leaf, amended read shape,
  failure code, and Human view, with negative unknown-field, authority,
  stale-version, and duplicate-delivery cases. Fixture authorship belongs to the
  earlier owning tasks; this task assembles and audits completeness, not first
  behavior. A missing case is a finding against its owning task, recorded with
  provenance, not silently authored here without recording that fact.
- New source helpers and fake adapters must be exercised by real assertions.
  Compiling an unused helper is not behavioral evidence.
- The cancelled-family absence proof covers tables and the real-process command
  tree, under any name or alias. `project-docs` delivery artifacts are outside
  product scope and must not be confused with product state.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `cmd/fid/integration/` - Phase 23 conformance corpus, fake and fid adapters,
  the audit tests, and regeneration of `p23_human_fixtures_test.go`'s manifest
  against the final binary.
- `cmd/fid/integration/testdata/` - the Phase 23 operation fixtures.
- `internal/workflow/authority.go` - audit support only, no new action.
- `project-docs/docs/phases/phase-23/tasks/task-18-done.md`.

### Files To Read For Context

- `public-interface.md` in full; `contracts.md` sections 13, 14, 15.
- `project-docs/docs/phases/phase-22/tasks/task-18.md` for the accepted audit
  pattern.
- `task-17-done.md`; `execution-scope.md`; `contract-coverage.md`;
  `v1-convergence-inventory.md`.

### Out Of Scope

- New product behavior, new codes, and any contract change.
- The production kit adapter, which stays Phase 25 under CF-037.

## Acceptance Criteria

- [ ] Exactly thirty grant actions are accepted, each at or after its owning
  task, with every unknown action rejected, proven against an independently
  pinned corpus; and the workflow and journal copies of the closed set are
  asserted equal to each other. (VAL-067)
- [ ] Fake and fid adapters produce identical accepted results and identical
  rejections for every Phase 23 case. (VAL-068)
- [ ] No cancelled product family is present under any name or alias. (VAL-069)
- [ ] Any missing fixture is recorded as a finding against its owning task with
  provenance.
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestGrantActionSetIsExactlyThirtyNames` fails against the incomplete audit.

## Documentation Updates

- `project-docs/docs/phases/phase-23/tasks/task-18-done.md` in the full format,
  including the corpus inventory and the per-operation fake and fid comparison.
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

- This task belongs to QA Review window 6.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-18`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
