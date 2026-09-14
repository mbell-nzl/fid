# Phase 23 - Task 7: Review freshness and Human Outcome supersession

## Objective

Stale target-bound Review evidence when the target tree changes, keep Human
judgment out of tree-based staling, and make explicit Human Outcome
supersession the only way a recorded Outcome stops governing evaluation.

## Validation Contract Coverage

- VAL-031: A target tree change stales linked Review evidence and retains the
  Review Record verbatim.
- VAL-032: `human-outcome` Obligations are not staled by a tree change.
- VAL-033: A later Human Outcome on the same class with a different result
  supersedes the earlier one for evaluation while both records remain readable.

## Binding Constraints

- `contracts.md` section 7 is the exact authority. Freshness reuses the existing
  `internal/evidence` deterministic repository-tree evaluation; add no new
  freshness mechanism, branch-label behavior, or external freshness source.
- Staling never deletes, rewrites, or invalidates a Review Record. Identity,
  coverage, verdict, and Findings are retained verbatim.
- `freshness_target:"none"` on `human-outcome` Obligations must read `fresh`,
  not a special-cased null; the wire value stays in its accepted domain.
- Supersession uses recorded Event order, never a caller-supplied timestamp,
  and reports reason `outcome-superseded`. Both records remain readable and
  unmodified.
- The three evidence dimensions stay independent; a stale record is still a
  valid record with an unavailable-or-not-applicable Artifact dimension.
- Human Outcome recording itself is task 10; this task implements only the
  evaluation rules over Outcomes and may use a seeded projection fixture until
  task 10 supplies the writer. Do not add a temporary writer.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/evidence/policy.go`, `evidence.go`.
- `internal/workflow/evaluation.go`, `resolve.go`.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change.
- `project-docs/docs/phases/phase-23/tasks/task-7-done.md`.

### Files To Read For Context

- `contracts.md` sections 2, 7, 13; `public-interface.md` section 2.
- `project-docs/docs/phases/phase-21/public-interface.md` freshness grammar.
- `task-6-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Human Outcome recording, authority, and provenance.
- Any new freshness target, branch labelling, or external artifact freshness.

## Acceptance Criteria

- [ ] After a target tree change the review evaluation reports `stale` with
  reason `evidence-stale` and the Review Record is unchanged. (VAL-031)
- [ ] A `human-outcome` Obligation stays `fresh` and satisfied across the same
  tree change. (VAL-032)
- [ ] A later differing Outcome of the same class governs with reason
  `outcome-superseded` and both records remain readable. (VAL-033)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestReviewEvidenceStalesOnTargetTreeChange` fails against the current
receipt-only freshness evaluation.

## Documentation Updates

- `project-docs/docs/phases/phase-23/tasks/task-7-done.md` in the full format.
- Flag spec impact; only Tech Lead accepts a contract or schema change.
## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
```

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

- This task belongs to QA Review window 3.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-7`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
