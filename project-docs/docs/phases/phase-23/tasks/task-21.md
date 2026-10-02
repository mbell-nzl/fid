# Phase 23 - Task 21: Resolvable review-rework exhaustion

> Supplemental task under DEC-268 / DEC-269 Part 1 (Human checklist finding A1).
> Dispatch only after DEC-269's status line reads ACCEPTED. If the Human chose
> an option other than option 1, follow the amended DEC-269 text and this
> spec's Planning Amendment, if any.

## Objective

The `review-rework-exhausted` Stop currently names `budget-observed` as its
resolution condition. No Phase 23 Event can produce that fact for the
fid-owned `review-cycles` unit, so the Stop can never be resolved. This task
makes an authorized Human Outcome of class `review-rework-extension` resolve
it. It also stops the exhausted rework budget from leaking into generic
budget admission at `gate.run` and in resume next actions. Finally, it gives
the Human-mode denial a guidance line naming the qualifying action.

## Validation Contract Coverage

- VAL-086: the Stop resolves only through an approved Human Outcome of class
  `review-rework-extension` scoped to the Envelope, after which rework is
  admitted.
- VAL-087: an exhausted `review-cycles` budget never produces a generic
  `budget-exhausted` or `budget-data-unavailable` block or next-action
  reason, while every other budget keeps its Phase 22 behavior.
- VAL-088: after resolution, later `changes-required` Reviews are each charged
  once and raise no further rework Stop, and the Human denial names the
  qualifying action while the Stop is open.
- VAL-072 (Human, re-checked at phase end as HUMAN-03a and both HUMAN-03b
  rows): the regenerated rework fixture.
- Keep VAL-027, VAL-028, VAL-029, VAL-030, VAL-046, VAL-048, VAL-050, VAL-078,
  VAL-080, VAL-081, VAL-082, VAL-083, and VAL-084 green.

## Binding Constraints

- `contracts.md` section 6 as amended by DEC-269. The rework Stop's condition
  is `{type:"human-outcome-recorded", subject:{type:"envelope", id:<envelope>},
  outcome_class:"review-rework-extension"}`. It uses the DEC-235 encoding:
  `budget_name` and `capability_id` are present and null, and
  `outcome_class` is set. Unchanged:
  - Stop code, StopDetails, blocked actions (`claim.acquire`, `run.start`),
    and resolver capability;
  - the `review-rework-exhausted` failure details;
  - the crossing rule (raise only when this charge exhausts the budget);
  - one charge per `changes-required` Review (VAL-027);
  - the `run.observe-budget` rejection of `review-cycles` (VAL-028).
- Resolution is DEC-250's shared predicate and the existing atomic path inside
  `human-outcome.recorded`, keyed on the Condition. Do **not** change the
  predicate, `QualifyingJudgmentStops`, or the Attention resolution path.
  Manual `stop resolve` with the Outcome's Event ID must also qualify.
- `budget-data-unavailable` and every other Stop keep their Phase 22
  conditions. Only the rework Stop changes.
- A budget whose `exhaustion_action` is `review-rework-exhausted` is excluded
  from generic `BudgetAdmissionBlocks` use at `gate.run`
  (`internal/evidence/execution.go`) and in resume next-action reasons
  (`internal/workflow/resume_actions.go`). Exclude it with one shared helper,
  not two copies. Every other unit's denial and Stop are unchanged.
- Human-mode denial for `review-rework-exhausted` only (on top of task 19's
  `error:` line) is exactly two lines:

  ```text
  error: review rework budget is exhausted
  next: an activated Human records an approved Human Outcome of class review-rework-extension for this Envelope's Issue; see fid stop list --subject envelope:<envelope-id>
  ```

  JSON is unchanged.
- The shared fake (`cmd/fid/integration/p22_fake_test.go`) mirrors the new
  condition and the admission exclusion so fake/fid conformance stays
  identical (VAL-068).
- Phase 23 adds no denial Event (DEC-226). Envelope revision stays 1.

## Scope

### Files To Create Or Modify

- `internal/workflow/budget.go` (`BudgetStopEffects` rework condition; the
  shared admission-exclusion helper)
- `internal/evidence/execution.go` (gate admission)
- `internal/workflow/resume_actions.go` (next-action reasons)
- `cmd/fid/stop.go` (rework guidance line only)
- Tests and fixtures: `internal/workflow/*_test.go`, `internal/evidence/*_test.go`,
  `internal/claims/generic_test.go`, `internal/review/review_test.go`,
  `cmd/fid/integration/p23_review_evidence_budget_test.go`,
  `cmd/fid/integration/p22_fake_test.go`,
  `cmd/fid/integration/p23_human_fixtures_test.go` (rework fixture),
  `cmd/fid/integration/p23_human_failure_test.go` (rework guidance row)
- `project-docs/docs/phases/phase-23/qa/human-fixtures.md` (rework section
  only: the resolution step)
- Retained assertions displaced by this change converge here with per-row
  provenance (`execution-scope.md` section 9)

### Files To Read For Context

- `project-docs/docs/phases/phase-23/decisions.md` DEC-269 Part 1, DEC-250,
  DEC-235
- `project-docs/docs/phases/phase-23/contracts.md` sections 6, 10, 11
- `internal/workflow/stop_conditions.go`, `internal/review/review.go`
  (rework charge), `internal/review/human_outcome.go`
- `project-docs/docs/phases/phase-23/qa/human-checklist.md`, finding A1

### Out Of Scope

- Option 2 re-trip semantics and any ordering rule in the predicate.
- Stop-details rendering (finding A4) and generic denial reasons (finding A5),
  which are on the wrap-up agenda.
- Envelope revision and fresh budgets (Phase 24).

## Acceptance Criteria

- [ ] Exhausting a `review-cycles` budget raises the rework Stop with the
  amended condition, encoded per DEC-235. (VAL-086)
- [ ] Non-qualifying Outcomes leave the Stop open, and `claim.acquire` /
  `run.start` still return `review-rework-exhausted` with no Claim or Run.
  The cases are: another class; `rejected`; an out-of-scope Issue; and a
  non-human Actor, which is rejected `authority-missing` with no Event.
  (VAL-086)
- [ ] The qualifying Outcome resolves the Stop atomically in its own Event, and
  both actions then succeed. On a second fixture, manual `stop resolve` with
  that Outcome's Event ID resolves it. (VAL-086)
- [ ] A Builder Run active when the exhausting Review commits can still
  `gate run` without a generic `budget-exhausted` from the rework budget, and
  resume lists no generic reason for it. An `attempts`-budget control still
  denies and raises its Phase 22 Stop. (VAL-087)
- [ ] After resolution, two further `changes-required` Reviews each charge
  exactly one cycle, `remaining` stays 0, and no new rework Stop is raised.
  (VAL-088)
- [ ] The Human-mode denial prints exactly the two lines above, and JSON is
  byte-identical to before. (VAL-088)
- [ ] The regenerated rework Human fixture seeds the exhausted state. It ships
  an activated human Actor and a ready qualifying Outcome input, and its
  instructions show the denial, the Stop, the resolution step, and the
  admitted retry. (VAL-072)

## First Failing Test

A built-binary journey exhausts the budget, records an approved
`review-rework-extension` Outcome as the activated human Actor, and expects
`claim acquire` to succeed. Today it gets `review-rework-exhausted`.

## Documentation Updates

- `qa/human-fixtures.md`, rework section: add the resolution step. Canonical
  specs are synced by the Tech Lead at wrap-up.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
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

- Report path: `project-docs/docs/phases/phase-23/tasks/task-21-done.md`.
  Commit with `Refs: phase-23/task-21`.
- Record red evidence for both defects separately: the unresolvable Stop, and
  the generic `gate.run` denial.
- If the VAL-087 state (a Builder Run still active when the exhausting Review
  commits) cannot be reached through the CLI under the accepted Claim and Run
  rules, cover it in-package through the shared exclusion helper at both call
  sites. Report it as an isolated-fixture route, as DEC-239 did for VAL-015.
  Do not weaken the assertion.
- Any need to touch the DEC-250 predicate, a closed journal registry, or an
  Effect row is a `BLOCKED` escalation candidate. DEC-269 found none
  necessary.
