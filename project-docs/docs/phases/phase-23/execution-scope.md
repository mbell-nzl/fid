# Phase 23 Execution Scope And Gate Reach

> Tech Lead planning inspection, 2026-09-14. ACCEPTED under DEC-222.
> No build, test, lint, or feasibility gate was run by Tech Lead.
> Base `0195897a347dfe1872db0f7be3800c5c3ccdefcf`.
> Amended 2026-09-15: section 7 adds the test-time headroom disposition
> (DEC-227, R-13), and three file scopes are corrected under DEC-167/173
> (R-08, R-09) in the affected task specs.

## 1. Ownership Rule

Every task owns its named implementation paths, the matching in-package tests,
and only those retained assertions or fixtures invalidated by its own accepted
change. The default gate compiles every `cmd/` and `internal/` package; its
required tagged companion also builds the CLI and runs real-process journeys.
Accordingly `cmd/`, `internal/`, their `testdata/`, root `testdata/`, and the
accepted migration chain are in scope for minimum causal convergence only, not
for unrelated behavior, cleanup, or redesign. This is the task pack's explicit
scope for those gate roots, never a waiver of narrow task ownership.

For each displaced retained row, the completion report records file and line,
introducing commit and phase/task, the accepted invalidating change, how the
row's subject is preserved or replaced, and the exact executing command. Labels
and filenames never assign ownership. A latent production defect first observed
by a new gate follows DEC-175: the Runner obtains the minimum causal scope
amendment before repair. This document grants no unrestricted repair authority.

## 2. Planning Inspection

Read-only evidence used to design the commands below:

- `cmd/fid/integration/main_test.go` carries `//go:build integration`; its
  `TestMain` builds `./cmd/fid` into fresh temporary storage and runs the
  tagged suite. Plain `go test ./...` excludes those files entirely.
- `internal/db/migrate.go` loads migrations from an explicit directory, so
  database and fresh-init tests exercise the actual SQL rather than a
  duplicated test-only schema. Every `schema.md` DDL row is therefore
  executable evidence, not documentation.
- `internal/journal/rebuild_bench_test.go` exposes
  `BenchmarkRebuildWorkspace1kEvents`, `...10kEvents`, and `...100kEvents`;
  the 100k case self-fails at 30 seconds or more. Tasks 1, 15, and 16 preserve
  that command and its accepted contract.
- Live packages at the base are `activity`, `bootstrap`, `claims`, `cli`,
  `config`, `db`, `evidence`, `export`, `issues`, `journal`, `labels`, `links`,
  `maintenance`, `policy`, `project`, `requirements`, `sqlquery`, `stats`,
  `validation`, and `workflow`. This phase adds exactly one, `internal/review`,
  and removes exactly one, `internal/requirements`.
- Hosted CI (`.github/workflows/ci.yml`) runs format, lint, build, the shard
  verifier, and four exhaustive ordinary race-and-coverage shards. It does
  **not** run the tagged integration suite or any benchmark. Those remain
  mandatory local task and phase gates and are never substituted by a green PR.
  Closing that hosted gap is workflow maintenance recorded as CF-041, not
  Phase 23 product scope.

This inspection confirms command selection and tag reach, not runtime success.
QA Lead performs fresh feasibility validation in Heavy mode before Ready.

## 3. Task-To-Gate Matrix

| Tasks | Owned behavior | Executing gate |
|---|---|---|
| 1, 15, 16 | Migration, DDL, rebuild, reset, current schema, portable schema | ordinary race and coverage, full tagged integration, plus the named fresh-process rebuild benchmark |
| 2, 3, 5, 6, 7, 8, 9, 10, 11 | Domain rules, Effects, SQLite projections, authority and failure branches | ordinary race and coverage plus full tagged integration; the tagged gate catches displaced adapter fixtures |
| 4, 13, 14, 17, 18 | CLI, JSON, Human text, shared semantic fixtures, process journeys, availability matrix | both gates; tagged integration is mandatory and never substituted |
| 12 | Requirement reference replacement and link-model removal | both gates plus the fresh-process rebuild benchmark, because it removes an Event family from the replay corpus, plus its own declared expected-non-zero absence search |
| 2, 6 | New branch coverage in `internal/workflow` and `internal/evidence` | named package `go test -race -cover` before and after edits |
| 15 | V1 removal absence proof | both gates plus the declared expected-non-zero absence searches in that task spec |

The tag selection is case-sensitive and remains exactly
`go test -tags integration -race -timeout 30m ./cmd/fid/integration/...`
(the `-timeout` is the DEC-227 bound in section 7; the tag selection itself is
unchanged). No other environment-gated suite is assigned. New helpers and fake
adapters must be used by real assertions; compiling an unused helper is not
behavioral evidence.

Every task runs the tagged integration gate from task 1. There is no delayed
first integration gate and therefore no planned convergence debt. If
implementation reveals a build tag or environment switch outside these
commands, amend the task before assigning that scope; never narrow a gate.

## 4. Ordered Delivery And Fence Ownership

`schema.md` section 6.1 assigns each progressive fence key to exactly one task:
task 4 `review.recorded`; task 5 `finding.addressed`; task 8
`concern.recorded`, `concern.disposed`, `deviation.recorded`; task 9
`decision.referenced`; task 10 `human-outcome.recorded`; task 12
`requirement.referenced`, `requirement.linked` plus the now-empty map and its
guard branch. Every other key stays fenced. Task 1's fence assertion is
scheduled convergence for each later owner, substituting a still-fenced Event
and preserving the parity subject; it is never deleted.

`public-interface.md` section 4.1 assigns each new grant action to its first
owning task. Adding an action before its owning task, or leaving one
unimplemented at task 18, both fail the final audit.

## 5. Review And Command Accounting

Three-task review windows: 1-3, 4-6, 7-9, 10-12, 13-15, 16-18. QA Reviewer
reviews each completed window.

An independent QA Lead checkpoint after task 11 must PASS before task 12. It
checks real Review independence, Finding resolution authority, review-cycle
charging and exhaustion, protected-Decision blocking and satisfaction, Human
Outcome provenance, and Attention resolution through a Human Outcome.

Acting on the Phase 22 workflow note that Human-readable guidance was checked
too late, the window 13-15 review additionally performs a deterministic
Human-text conformance check against `contracts.md` section 14, using the
verbatim fixtures task 14 records in
`project-docs/docs/phases/phase-23/qa/human-fixtures.md`. A Human-owned
validation assertion still covers product judgment in the human checklist.

Runner report lint is a post-report-commit, pre-next-dispatch control using
`project-docs/tools/report-lint.sh` with exactly two positional arguments: the
completed `task-N-done.md` followed by its `task-N.md`, both under
`project-docs/docs/phases/phase-23/tasks/`. Concrete N is 1..18. This is Runner
orchestration and never a Developer `Exact Commands` obligation. The first
concrete invocation is:

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-1-done.md project-docs/docs/phases/phase-23/tasks/task-1.md
```

Required commands report their own output token and observed termination. Tasks
12 and 15 declare expected non-zero commands; task 12's declaration was added
2026-09-15 with its `requirement.synced` absence search (R-07d). Red tests retain their failing
evidence and never authorize weakening a gate. No cleanup is a gate. Do not
impose a parent session lifetime as a command timeout. External interruption
retains no-progress state and follows current schema-v2 recovery.

## 6. Corrected File Scopes (2026-09-15, extended 2026-09-19)

Three assignments were incomplete under DEC-167/173 and are corrected in the
task specs before dispatch, not after a failure:

- The closed 23-name activated-capability and grant-action allowlist is
  `sortedUniqueCapabilities` in `internal/workflow/activation.go`; Envelope
  grants reuse it through `sortedUniqueEnvelopeActions` in `envelope.go`, and
  `authority.go` holds no action list. Tasks 4, 5, 8, 9, and 10 each add
  `internal/workflow/activation.go` to their file scope for their own new
  action (R-08).
- `internal/export/archive_semantics_catalog.go` carries the
  `requirement.synced` archive semantic rule. Task 12 removes the Event, so that
  file joins task 12's scope; without it, task 12's own declared absence search
  fails with its cause outside the task (R-09).
- Task 10 owns the readiness-Stop effect rows added to `evidence.linked` and
  `artifact.availability-observed` in `internal/journal/current_event_catalog.go`
  (DEC-224), alongside its fence key.

Four further closed registries had no owning task. Corrected under R-14
(DEC-229); same rule, same reason - a task owns only what it names, and a
registry nobody names is rejected at runtime with its cause outside the failing
task's scope:

- `internal/workflow/stop.go` holds the closed `stopCodes`, `conditionTypes`,
  `Condition`, and `validStopDetails` registries, and
  `internal/journal/current_event_value_catalog.go` holds their journal-side
  grammar (`currentStopDetailShapes`, which `StopCode` validation itself keys
  on, and the `Condition` value shape). Both files join the scope of every task
  that adds to them: task 6 (`review-rework-exhausted`), task 9
  (`protected-decision-required`, `decision-satisfied`, `decision_class`),
  task 10 (`human-outcome-required`, `human-outcome-recorded`,
  `outcome_class`), and task 11 (`review-recorded`, `finding-resolved`,
  `finding_id`). A Condition type is registered by the task that first raises a
  Stop carrying it, never deferred to task 11; task 11 completes the closed
  four-type set and owns the resolution behavior of all four (VAL-078).
- `internal/journal/current_event_columns.go` joins task 1, which registers all
  ten new Phase 23 tables in `currentEventTables` alongside the nine effect
  rules it already registers. `validateCurrentEffect` returns `unknown table`
  for an effect on an unregistered table, and
  `TestCurrentEventColumnDomainsMatchAcceptedDDL` proves the registration
  against the 0017 DDL in the same task that ships it. Task 15 renames the
  `judgment_concerns` key with its table; no other task edits this map.
- `internal/export/archive_semantics_catalog.go` mirrors the journal catalogs
  for `export validate` and `import`, and no test asserts parity between them.
  Standing rule: **a task that adds, amends, or removes a journal Event rule
  mirrors that same change in `semanticEventRules` in the same task.** Task 1
  mirrors the nine new rules, task 9 its two amendments, task 10 its two, task
  12 the `requirement.synced` deletion (R-09, already assigned), and task 15
  the concerns rename. Without this, a divergence stays silent until VAL-064 at
  task 16, where the cause sits outside that task's ordinary scope.

A fifth gap, larger than the four above, was found while confirming them.
Corrected under R-16 (DEC-231):

- Registering the nine families in `currentEventRules` also routes them onto the
  journal's **strict** result path. `validateCurrentAppliedResults`
  (`current_event_result_projection.go:12`, called from `project_effects.go:20`)
  gates on `currentEventRules` membership, and `schema.md` 6.3 makes a typed
  snapshot mandatory for every Phase 23 mutation, so each family must also clear
  three further closed registries: the inline `fields` grammar map in
  `internal/journal/current_event_json.go` (else `unregistered current result
  family`), `currentExpectedResult`'s switch in
  `internal/journal/current_event_result_projection.go` (else `unknown current
  result`), and `currentValueShapes` in
  `internal/journal/current_event_value_catalog.go`, which
  `currentProjectedValue` reads the projected field list from.
- **Standing rule: a family's typed-result grammar, its `currentExpectedResult`
  branch, its value shapes, and any `currentProjectedValue` fill helper are
  added by the task that deletes that family's fence key - never by task 1.**
  All three registries execute only when an Event of the family is appended, so
  the unfencing task is the earliest point at which any of them is provable, and
  it is the point at which the task's own first successful Event proves all
  three at once. This is the same principle DEC-229 applied to Condition types,
  which are registered by the task that first raises a Stop carrying them.
  Task 1 is deliberately **not** the owner here, unlike
  `current_event_columns.go`: that catalog is proved by
  `TestCurrentEventColumnDomainsMatchAcceptedDDL`, which introspects the 0017
  DDL task 1 ships and needs no Event at all, whereas these three would sit
  behind task 1's own fence unexecuted. Ownership is only assignable where it is
  executable.
- Owners, one-to-one with the nine fence keys: task 4 `review.recorded`; task 5
  `finding.addressed`; task 8 `concern.recorded`, `concern.disposed`,
  `deviation.recorded`; task 9 `decision.referenced`; task 10
  `human-outcome.recorded`; task 12 `requirement.referenced`,
  `requirement.linked`. A value shape is registered once, by the first of those
  tasks to need it, and reused unchanged afterwards - so task 5 adds no shape,
  and task 10 reuses task 9's `ProtectedDecision`.
- Two shared shapes are widened outside that set, each by the task that makes
  the wider domain reachable: `currentValueShapes["Claim"]`'s `mode:=mutation`
  becomes `mutation,review` at **task 3**, reachable because `run.abandoned` and
  `run.recovered` both project a `Claim`; and `["Budget"]` / `["ImmutableBudget"]`
  gain `review-cycles` in `unit` and `review-rework-exhausted` in
  `exhaustion_action` at **task 6**. Task 4 additionally widens the `Subject`
  shape's `type:=run,envelope,activation` to the complete accepted SubjectType
  set in `public-interface.md` section 2, in one monotonic edit, as the first
  task to register a value shape referencing `Subject`; no later task re-narrows
  it.
- Note for every owner: `idempotency_results` is encoded as an optional stage in
  `currentEventRules` for nearly every family, but `validateCurrentAppliedResults`
  requires the result row unconditionally. "Optional" in `schema.md` 6.2 is the
  stage encoding, not permission to omit the typed snapshot.

Not corrected, and deliberately so: tasks 4 and 5 do **not** gain
`internal/journal/current_event_catalog.go`. Task 1's binding constraints
already register all nine Phase 23 Event types, their operation mappings, and
their complete `schema.md` section 6.2 effect rules, so `review.recorded` and
`finding.addressed` need no catalog edit at tasks 4 and 5. Those tasks own only
their fence key in `internal/journal/project_effects.go`. That rejection is
unaffected by R-16: the three registries above are different files, and tasks 4
and 5 gain those, not `current_event_catalog.go`.

## 7. Test-Time Headroom (DEC-227)

Fresh QA Lead measurements at the planning base: the tagged integration suite
runs 560.2 s against Go's implicit 600 s per-binary default, with `cmd/fid` at
472.0 s and `internal/workflow` at 458.4 s under the race gate, and
`t.Parallel()` appears in none of the 58 `cmd/fid/integration` files. Phase 23
adds real-process journeys to exactly those packages at tasks 4, 6, 11, 14, 17,
and 18, so an unbounded default would convert ordinary growth into a killed
command - which is `BLOCKED`, never `PASS`, and whose only cure inside a task
would be a protected command change.

The disposition:

- Every Phase 23 task runs the two long gates with an explicit `-timeout 30m`.
  The bound is deliberately generous: it makes the harness limit a declared
  constant instead of an accidental cliff, and it weakens no assertion, changes
  no scope, and removes no gate.
- Wall time stays a measured budget, not an aspiration. Every completion report
  records the times `go test` reports for `cmd/fid/integration`, `cmd/fid`, and
  `internal/workflow`. Budgets: tagged integration under 900 s, no ordinary
  package over 600 s, asserted at phase end by VAL-079.
- Speed repair follows causation. A task whose own new journeys push a package
  past its budget owns the bounded repair - reuse of seeded workspaces, fewer
  redundant end-to-end setups, and `t.Parallel()` in **new** Phase 23 journey
  files where the fixture is self-contained. A retrofit of the 58 existing files
  is not in any task's scope; if the cause is pre-existing, stop and escalate
  rather than absorbing it.
- Task 1 additionally stops the per-run integration binary leak in
  `cmd/fid/integration/main_test.go` (CF-042, cause corrected under R-15 /
  DEC-229). The dominant cause is not timeout kills: five tests re-execute the
  test binary, each child re-enters `runIntegrationSuite` and rebuilds its own
  ~16 MB binary, and the crash-process children never reach `RemoveAll`. One
  clean run - exit 0, no kill - leaks 20 directories and 306 MB into a 3.9 GB
  RAM-backed `/tmp`, so this phase's ~26 tagged runs would need about 7.8 GB and
  fill `/tmp` around the twelfth. Because the leak is per-run, not per-kill, an
  age bound prunes nothing within a working day. The mitigation is therefore
  causal: `runIntegrationSuite` publishes the built binary path, a process that
  finds it reuses it and neither builds nor removes, and the 24-hour prune stays
  only as the residual safety net for a parent killed mid-run. Reuse also
  removes about twenty redundant `go build ./cmd/fid` invocations per run, which
  helps the budgets above. A full `/tmp` is a gate-blocking environment failure,
  not a test failure.

## 8. Incorporating Earlier Delivery Lessons

- `contract-coverage.md` maps every binding contract clause to its schema row,
  owning task, and validation assertion. It exists because Phase 22 shipped a
  cancellation clause with no schema row (DEC-195) and an undelivered checkpoint
  allocation (DEC-194). QA Lead re-derives that matrix independently before
  execution approval; a clause with no schema row or no owning task is a
  planning defect, not an implementation surprise.
- Cross-cutting validators reach many producers. The Obligation generalization
  in task 2 and the Evidence widening in task 6 are the two known wide-reach
  changes; both are isolated to their own task rather than bundled with a leaf.
- End-to-end journeys enter with each implemented wire surface at tasks 4, 6,
  11, and 14; task 18 audits their completeness, not their first behavior.
- Before a BLOCK is emitted, its allocated ESC must already exist in the
  current escalation register. Mis-scoped directives and externally interrupted
  sessions never consume a no-progress recovery budget.
