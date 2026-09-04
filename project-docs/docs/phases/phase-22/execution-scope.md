# Phase 22 Execution Scope And Gate Reach

> Tech Lead planning inspection, 2026-09-05; DEC-181.
> No build, test, lint, or feasibility gate was run by Tech Lead.

## 1. Ownership Rule

Every task owns its named implementation paths, matching package tests, and
only those retained assertions/fixtures invalidated by its accepted change.
The default gate compiles all `cmd/` and `internal/` packages. Its required
integration companion also builds the CLI and runs tagged subprocess journeys.
Accordingly, `cmd/`, `internal/`, their `testdata/`, root `testdata/`, and the
accepted migration chain are in scope for minimum causal convergence, not for
unrelated behavior, cleanup, or redesign. This is the task pack's explicit
scope for those gate roots, not a waiver of narrow task ownership.

For each displaced retained row, record file/line, introducing commit and
phase/task, the accepted invalidating change, preservation/replacement of its
subject, and the exact executing command. Labels and filenames never assign
ownership. A latent production defect first observed by a new gate follows
DEC-175: the Runner obtains the minimum causal scope amendment before repair;
this document grants no unrestricted product-code repair authority.

All tasks run both ordinary race/coverage and tagged integration gates. There
is no delayed first integration gate and hence no planned convergence debt.
If implementation reveals a build tag or environment switch outside those
commands, amend the task before assigning that scope. Do not narrow a gate.

## 2. Planning Inspection

Read-only evidence used to design commands:

- `cmd/fid/integration/main_test.go` has `//go:build integration`; TestMain
  builds `./cmd/fid` into fresh temporary storage and runs the suite. Ordinary
  `go test ./...` alone excludes those files.
- `internal/db/migrate.go` loads migration files from an explicit directory;
  database integration tests and process fresh-init tests exercise the actual
  SQL, rather than a duplicated test-only schema.
- Existing `internal/journal/rebuild_bench_test.go` exposes
  `BenchmarkRebuildWorkspace100kEvents`; one-iteration invocations self-fail
  at >=30 seconds. Task 17 preserves that command and its accepted contract.
- The baseline has live packages `activity`, `bootstrap`, `claims`, `cli`,
  `config`, `db`, `evidence`, `export`, `issues`, `journal`, `labels`, `links`,
  `maintenance`, `policy`, `project`, `requirements`, `sqlquery`, `stats`,
  `validation`, and `workflow`. `internal/gitctx` is not currently present.
  Repository observation is assigned to the existing `internal/project` adapter.

This inspection confirms command selection and tag reach, not runtime success.
QA Lead must perform fresh feasibility validation in Heavy mode before Ready.

## 3. Task-To-Gate Matrix

| Tasks | Owned behavior | Executing gate |
|---|---|---|
| 1,16,17 | Migration/DDL/reset/current schema | ordinary race/coverage plus full tagged integration |
| 2-5,7-8,10,12 | Domain rules, Effects, SQLite, authority and failure branches | ordinary race/coverage; tagged gate catches displaced adapter fixtures |
| 6,9,11,15,18 | CLI/JSON/live outputs, shared semantic fixtures, process journeys | both gates; tagged integration is mandatory, never substituted |
| 13-14 | Archive parser, effect validation, staging/publication fault injection | ordinary race/coverage; process corpus advances in task 15 |
| 2,7 | CF-034 workflow/evidence branch coverage | named package `go test -race -cover` before and after edits |
| 17 | CF-024 reference rebuild and phase performance | three exact fresh-process benchmark commands, then general benchmark gate |

The tag selection is case-sensitive and remains exactly:
`go test -tags integration -race ./cmd/fid/integration/...`.
No other environment-gated suite is assigned. New source helpers and fake
adapters must be used by real assertions; merely compiling unused helpers is
not behavioral evidence.

## 4. Review And Command Accounting

Three-task windows: 1-3, 4-6, 7-9, 10-12, 13-15, 16-18. QA Reviewer reviews
each completed window. Independent QA Lead checkpoint after task 12 must PASS
before task 13; it checks actual recovery, Stop/budget effects, activation
direction, freshness, and the feasibility of import fault injection.

Runner report lint is a post-report-commit/pre-next-dispatch control using
`project-docs/tools/report-lint.sh` with exactly two positional arguments:
the completed task's `task-N-done.md`, followed by its `task-N.md`, both under
`project-docs/docs/phases/phase-22/tasks/`. The task index gives concrete N=1..18;
this is Runner orchestration, never a Developer Exact Commands obligation.
The first concrete invocation is:

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-22/tasks/task-1-done.md project-docs/docs/phases/phase-22/tasks/task-1.md
```

Required commands report their own output token and observed termination.
No expected-nonzero gate is planned. Red tests retain their expected failing
evidence; they do not authorize weakening a gate. No cleanup is a gate.
Do not impose a parent session lifetime as a command timeout. External
interruption retains no-progress state and follows current schema-v2 recovery.

## 5. Incorporating Earlier Delivery Lessons

The early task-6 journey is extended at 9, 11/12, and 15, rather than first
executed in task 18. Conformance cases/fakes enter with each implemented wire
surface; task 18 audits their completeness, not their first behavior.

Before a BLOCK is emitted, its allocated ESC must exist in the current
escalation register. Phase-scoped/cross-task recovery identity limitations
remain governed by the existing workflow; this product phase does not rewrite
the sequence verifier. A dedicated workflow-maintenance proposal is recorded
as CF-036. Mis-scoped directives and externally interrupted sessions never
consume a no-progress recovery budget.
