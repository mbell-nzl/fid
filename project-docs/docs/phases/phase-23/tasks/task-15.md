# Phase 23 - Task 15: Migration 0018 and V1 judgment removal

## Objective

Execute the owned replacement-bound removal: drop the V1 judgment tables,
rename the generic Concern table, rebuild the affected view and triggers,
delete the role-bound transition rules and their evaluators, and prove every
removed surface absent.

## Validation Contract Coverage

- VAL-004: `foreign_key_check` and `integrity_check` are clean after 0018 and
  after a full rebuild.
- VAL-061: Every removed table, index, trigger, view column, column, Event type,
  rule, and CLI leaf is absent from the latest schema and the built binary. The
  inventory must include the `requirements.last_synced_hash` column (moved here
  from VAL-053) and the two triggers 0018 drops implicitly with their tables,
  `decisions_epic_kind_guard_insert` and
  `task_completions_issue_kind_guard_insert`. Record the measured
  `issues_enriched` column list: it is identical across 0018, so the accepted
  view-column removal set is empty and that clause is satisfied vacuously.
- VAL-062: The rebuilt `issues_enriched` view and link and activity triggers
  behave exactly as accepted, and a seeded 17-to-18 probe holding both Issue and
  non-Issue `links` and `activity_log` rows proves 0018's deletes remove exactly
  the non-Issue rows.
- VAL-063: Maintenance checks over removed tables are removed and the remaining
  checks pass.

## Binding Constraints

- `schema.md` sections 4 and 5 are the exact SQL authority; the removal targets
  are exactly the rows in `v1-convergence-inventory.md` section 1.
- Removal is replacement-bound: every removed behavior already has a passing
  replacement from tasks 2 through 14. Do not remove anything whose replacement
  is not yet green, and do not defer an owned removal to Phase 25 because a test
  breaks.
- `agent_runs.completion_id` is dropped because its foreign key targets the
  removed `task_completions`. The rest of `agent_runs`, `triage_events`,
  `epic_metadata`, and historical stats stay for Phase 24 and Phase 25.
- `acceptance_criteria` add and list are retained target behavior; only the
  role-keyed verification projection and the `verify` leaf are removed.
- Rename `judgment_concerns` to `concerns` and change exactly the two concern
  effect-rule table names. No other effect row moves.
- Manifest `projection_schema` becomes 18 in this task.
- Absence must be proven by a projection or real-process probe, not by a source
  search alone. The declared absence searches below are supporting evidence.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `migrations/0018_p23_remove_v1_judgment.up.sql`, `.down.sql` - new.
- `internal/issues/` - delete `ac_verify.go`, `concern_raise.go`,
  `concern_address.go`, `concern_list.go`, `capability_unavailable.go`, and the
  completion and verification readers in `detail.go`, `get.go`, `list.go`,
  `issue.go`, `read.go`, `ac.go`, `transition_rules.go`.
- `internal/workflow/rules.go`, `registry.go`, `workflow.go` - delete the three
  role-bound rules, their evaluators, and the roles parameter.
- `internal/maintenance/errors.go`, `warnings.go`, `operational.go`, `check.go`.
- `internal/stats/stats.go` - the dropped column only.
- `internal/journal/current_event_catalog.go`,
  `internal/journal/current_event_columns.go`, and
  `internal/export/archive_semantics_catalog.go` - carry the
  `judgment_concerns` -> `concerns` rename through all three closed catalogs in
  this task: the two concern effect rules, the `currentEventTables` key with
  its column domains, and the mirrored `semanticEventRules` rows (R-14). No
  other rule or table moves.
- `internal/journal/`, `internal/db/` - latest-schema expectations and rebuild
  parity.
- `internal/export/` - projection schema 18 and its fixtures.
- `cmd/fid/issue_ac.go`, `concern.go` - delete the removed leaves.
- Matching in-package tests and `cmd/fid/integration/` fixtures, including the
  four pinned fail-closed conformance fixtures.
- `project-docs/docs/phases/phase-23/tasks/task-15-done.md`.

### Files To Read For Context

- `schema.md` sections 4, 5, 7; `contracts.md` section 15;
  `v1-convergence-inventory.md` in full.
- `migrations/0014_p21_retire_v1_authority_execution.up.sql` and
  `migrations/0016_p22_remove_escalations.up.sql` as the removal precedents.
- `task-14-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- `triage_events`, the triage-status projection, `agent_runs` beyond its dropped
  column, `epic_metadata`, historical stats behavior, and old migrations.
- Capability advertisement, which is task 17.

## Acceptance Criteria

- [ ] `PRAGMA foreign_key_check` and `PRAGMA integrity_check` are clean after
  0018 and after a full journal rebuild. (VAL-004)
- [ ] Every inventory removal row is proven absent by a projection or
  real-process probe. (VAL-061)
- [ ] Open-Concern counts derive from generic Concerns for Issue subjects only
  and each rebuilt trigger rejects exactly its declared invalid input.
  (VAL-062)
- [ ] No maintenance check references a removed table and the command reports no
  unexpected finding. (VAL-063)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestLatestSchemaHasNoV1JudgmentTables` fails against the present V1 tables.

## Documentation Updates

- `project-docs/docs/phases/phase-23/tasks/task-15-done.md` in the full format,
  including per-row absence evidence for the whole inventory.
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
go test -race -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -timeout 30m ./cmd/fid/integration/...
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/
rg -n --glob '!*_test.go' "task_completions|ac_verifications" cmd internal
rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal
rg -n --glob '!*_test.go' "has_completion_report|concerns_addressed" cmd internal
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

- the named fresh-process rebuild benchmark, because this task changes the latest schema and the replay corpus.

## Expected Non-Zero Commands

These three commands prove absence, so their pass condition is a non-zero exit.
Record the real observed status; never convert it to `exit 0 (no output)`.

- `rg -n --glob '!*_test.go' "task_completions|ac_verifications" cmd internal` - expected `exit_code: 1` - no production reader or writer of the removed completion and verification tables remains.
- `rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal` - expected `exit_code: 1` - the removed requirement sync Event and column have no remaining production producer or consumer.
- `rg -n --glob '!*_test.go' "has_completion_report|concerns_addressed" cmd internal` - expected `exit_code: 1` - the two removed transition rules have no remaining evaluator or registration.

The `!*_test.go` exclusion is deliberate and is not a narrowed gate. Tests that
are explicitly scoped to earlier migrations retain those historical names by
design, and `migrations/` legitimately contains them in the files that create
and drop them. The authoritative absence proof is the latest-schema projection
query and the real-process leaf probe recorded under VAL-061; these searches are
supporting evidence for the production source tree only. If any search returns a
match, it is a real finding: repair the production reader or return
`RECOVERY_REQUIRED`. Never widen or narrow the declared command to make it pass.

## Notes For Developer

- This task belongs to QA Review window 5.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-15`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
