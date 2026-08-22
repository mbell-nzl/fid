# Phase 21 - Task 22: Real-Process Walking Skeleton And Conformance

## Objective

Complete the public Phase 21 adapter conformance corpus and the dedicated 13-step real-process one-Issue journey, including every negative counterpart and repeat after rebuild. Prove all accepted machine leaves, capability claims, removals, and canonical outputs as one coherent slice.

## Validation Contract Coverage

- VAL-058: the full real-process journey succeeds before and after rebuild.
- VAL-059: every authority/policy/idempotency/evidence/completion negative fails closed.
- VAL-060: accepted machine leaves conform to the canonical Response Envelope.
- VAL-061: delivery evidence is complete for phase-end audit.

## Binding Constraints

- Use the built CLI in isolated temporary projects; do not replace real process/Git/SQLite behavior with in-process calls.
- Do not read held-out scenarios or evaluator prompts.
- Any conformance mismatch returns to its owning task; task 22 does not change protected semantics.

## Scope

### Files To Create Or Modify

- `cmd/fid/integration/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-22-done.md`
- `project-docs/docs/playbook.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/gate-d3-2-proposal.md`
- `project-docs/docs/phases/phase-21/contracts.md`
- `project-docs/docs/phases/phase-21/qa/validation-contract.md`

### Out Of Scope

- Held-out validation, Phase 22 behavior, and baseline release-readiness claims.

## Acceptance Criteria

- [ ] The exact 13-step built-CLI journey completes and matches after rebuild. (VAL-058)
- [ ] All negative counterparts return contracted non-zero typed results with zero forbidden effects. (VAL-059)
- [ ] Every accepted non-streaming leaf passes shared Response Envelope/unknown-field/channel conformance. (VAL-060)
- [ ] Completion report links every review/report-lint/session-captured evidence artifact. (VAL-061)

## First Failing Test

Add the failing built-CLI journey through capability discovery, then extend stepwise to the full contract and negative matrix.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-22-done.md` - record red/green evidence, deviations, required command results, and spec impact.
- `project-docs/docs/playbook.md` - document the accepted operator journey and exact deferred operations.

## Exact Commands To Run

Run exactly:

```text
go test -tags integration -race ./cmd/fid/integration/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go build -o /tmp/fid ./cmd/fid
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-22-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after task 21 and QA Review 7 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-22-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
