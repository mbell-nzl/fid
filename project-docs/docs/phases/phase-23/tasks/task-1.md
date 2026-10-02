# Phase 23 - Task 1: Review and judgment schema baseline

## Objective

Land the Phase 23 persistence baseline: the accepted textual changes to
migrations 0013 and 0015, the new migration `0017_p23_review_judgment` with its
down, the `evidence_records` rebuild, and the progressive availability fence
that rejects every Phase 23 Event until its domain exists.

## Validation Contract Coverage

- VAL-001: Migration 0017 creates exactly the declared objects and nothing else,
  introspected at accepted schema 17.
- VAL-002: The Obligation registration discriminator enforces both states.
- VAL-075: The 0013 and 0015 baseline diff is exactly the seven accepted hunks.
- VAL-003: The progressive fence rejects every unimplemented Phase 23 Event.
- VAL-004: `foreign_key_check` and `integrity_check` are clean after 0017 and
  after a full rebuild.
- VAL-005: The `evidence_records` rebuild preserves every prior receipt row
  **and its Artifact links** on a seeded workspace migrated in place.

Full sources, pass and fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-23/qa/validation-contract.md`.

## Binding Constraints

- APPROVED under the accepted slice contract; execute after the Runner's
  required QA readiness check.
- `schema.md` sections 1, 2, 3, 6.1, and 7 are the exact authority. Do not add,
  rename, reorder, or re-type a single column, index, trigger, or CHECK.
- DEC-223: the `evidence_artifact_links` staging in the `evidence_records`
  rebuild is accepted DDL, not an optimization. The runner migrates on a
  `foreign_keys(1)` connection with one transaction per file, so
  `DROP TABLE evidence_records` performs an implicit `DELETE` that violates
  `evidence_artifact_links.evidence_id` whenever a link row exists;
  `PRAGMA defer_foreign_keys` moves the same failure to `COMMIT`. Ship the
  accepted staging order - stage, delete, rebuild, rename, restore, drop the
  staging table - in both the up and the down direction, and do not rebuild
  `evidence_artifact_links` itself.
- I-13 edit locality: `registration_id TEXT NOT NULL,` and
  `registration_version INTEGER NOT NULL,` also occur in `execution_attempts`
  (`0013` lines 222-223). The replacement applies **only** inside
  `envelope_obligations`; a global substitution fails DEC-178 and VAL-075.
- Two facts measured on the real chain during QA readiness (I-26). The product
  exposes no single-step down - `internal/db.MigrateDown` runs `m.Down()` to
  zero - so the 0017 down evidence is either the full teardown through that API
  or a test-local `Steps(-1)`; name which one the report uses. A failed
  migration leaves `schema_migrations` dirty (observed `16|1`), which is
  precisely the state VAL-005 rules out: after the accepted up and the accepted
  down, `schema_migrations` must be `17|0` and `16|0` respectively, never dirty.
- DEC-178: shipped migration SQL and the canonical slice DDL must agree exactly.
  Every line of 0013 and 0015 outside the accepted replacement table stays
  byte-identical.
- This is a deliberate pre-adoption baseline replacement requiring fresh init
  and rebuild. It is not an upgrade of old fid data and needs no converter.
- DEC-192 fence protocol: add `unavailableJudgmentEvents` to
  `internal/journal/project_effects.go` rejecting all nine Phase 23 Event types
  listed in `schema.md` section 6. Delete no key in this task. Later owners
  delete exactly their own key; task 12 removes the empty map and its guard.
- Register the nine Event types and their operation mappings in
  `internal/journal/current_event_catalog.go` with the effect rules in
  `schema.md` section 6.2, so a fenced type is rejected by the fence and not by
  an unknown-type path.
- Registration is not complete until all three closed catalogs agree (R-14).
  The same nine rules require their ten new effect tables registered in
  `internal/journal/current_event_columns.go` - `validateCurrentEffect` returns
  `unknown table` for an effect on an unregistered table - and their mirrored
  rows in `internal/export/archive_semantics_catalog.go`, which
  `export validate` and `import` key on. Column domains are copied from the
  0017 DDL this task ships, and
  `TestCurrentEventColumnDomainsMatchAcceptedDDL` proves that agreement here.
- Manifest `projection_schema` becomes 17 in this task. It becomes 18 in task
  15; do not anticipate that value.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this accepted
  change, and every owned path executes under the commands below.

## Scope

### Files To Create Or Modify

- `migrations/0013_p21_operational_kernel.up.sql` - only the accepted rows.
- `migrations/0015_p22_recovery.up.sql` - only the accepted rows.
- `migrations/0017_p23_review_judgment.up.sql`, `.down.sql` - new.
- `internal/db/` - migration loading, latest-schema expectations, tests.
- `internal/journal/project_effects.go`, `current_event_catalog.go` and tests.
- `internal/journal/current_event_columns.go` - register all ten new Phase 23
  tables in `currentEventTables` with the column domains of the 0017 DDL:
  `review_records`, `review_findings`, `finding_dispositions`,
  `judgment_concerns`, `deviations`, `decision_refs`, `human_outcomes`,
  `envelope_protected_decisions`, `requirement_refs`, and
  `requirement_issue_links` (R-14). Task 15 renames the `judgment_concerns`
  key with the table; no other task edits this map.
- `internal/export/archive_semantics_catalog.go` - mirror the nine new
  `currentEventRules` rows into `semanticEventRules` from the same
  `schema.md` section 6.2 source (R-14). Standing rule in
  `execution-scope.md` section 6: a task that changes a journal Event rule
  mirrors it here in the same task.
  In the same file and the same change, add the nine families to
  `semanticPhase22Events` at `:85` and correct that map's name to say what it
  means - families that do not exist at projection 14 (DEC-233). Mirroring the
  rules makes `found` true for the nine at **every** profile, so
  `validateEvent` (`archive_semantics.go:328`) would otherwise accept
  `review.recorded` inside a legacy projection-14 archive: the exact hole Phase
  22 added that map to close. This is provable with a crafted legacy archive and
  no Event, which is why it is task 1's and not an unfencing task's. VAL-065 is
  unaffected - the manifest check at `archive_validation.go:130` already rejects
  a cross-version archive with `schema-unsupported` before any Event is read.
- `internal/export/` - manifest projection schema and its fixtures.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change.
- `cmd/fid/integration/main_test.go` - stop the per-run binary-directory leak
  (test infrastructure only, no product behavior; CF-042, corrected under
  R-15). Measured cause: five tests re-execute the test binary as a child
  process, each child re-enters `runIntegrationSuite`, creates its own
  directory and rebuilds the same ~16 MB binary, and the crash-process children
  exit abnormally by design and never reach `RemoveAll`. One clean tagged run -
  exit 0, no kill, no timeout - leaked **20 directories and 306 MB** into a
  3.9 GB RAM-backed `/tmp`. At roughly 26 tagged runs this phase that is about
  7.8 GB, so `/tmp` fills around the twelfth run and the gates become
  environment-blocked rather than failing. Required change: have
  `runIntegrationSuite` publish the built binary path in an environment
  variable after a successful build, and have a process that finds that
  variable set and pointing at an existing executable reuse it - building
  nothing and removing nothing. Only the building process removes its own
  directory, exactly as today. All five re-executing call sites already spawn
  with `append(os.Environ(), ...)`, so they inherit it and need no edit; a
  variable pointing at a missing path must fall back to building. Keep the
  bounded 24-hour prune of stale `fid-integration-binary-*` directories before
  creating a directory, as the residual safety net for a parent killed mid-run;
  age is the only safe bound, because a count- or newest-based prune would
  delete a concurrent run's live binary. Record in the completion report the
  `fid-integration-binary-*` directory count and `/tmp` free space immediately
  before and after the tagged gate; a clean run must leak none. The durable
  cleanup-on-kill design stays with CF-042.
- `project-docs/docs/phases/phase-23/tasks/task-1-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-23/schema.md` - sections 1, 2, 3, 6, 7.
- `project-docs/docs/phases/phase-23/contracts.md` - sections 1 and 2 only.
- `project-docs/docs/phases/phase-23/execution-scope.md`,
  `v1-convergence-inventory.md`, `contract-coverage.md`.
- `migrations/0015_p22_recovery.up.sql` and `0016_p22_remove_escalations.up.sql`
  as the closest precedent for shape and down behavior.

### Out Of Scope

- Any domain rule, CLI leaf, Human text, or capability advertisement.
- Migration 0018 and every V1 removal.
- Deleting any fence key.

## Acceptance Criteria

- [ ] Latest-schema object inventory at schema 17 equals the accepted DDL
  exactly. (VAL-001)
- [ ] Both invalid and both valid registration-discriminator inserts behave as
  accepted. (VAL-002)
- [ ] The `0013` and `0015` diff against the planning base is exactly the seven
  accepted hunks, all inside `envelope_obligations`, `work_claims`, `stops`, and
  `attention_items`; every other line is byte-identical. (VAL-075)
- [ ] Every one of the nine Phase 23 Event types is rejected by the fence with
  a byte-identical projection afterwards. (VAL-003)
- [ ] `PRAGMA foreign_key_check` and `PRAGMA integrity_check` are clean after
  0017 and after a full journal rebuild. (VAL-004)
- [ ] On a seeded schema-16 workspace whose receipt Evidence carries
  `evidence_artifact_links` rows, migrated in place through the ordinary
  catch-up path, the migration succeeds, every Evidence row is identical column
  by column with the two added columns (`review_id`, `human_outcome_id`) null
  (count corrected by DEC-267), every link row is restored with its
  exact `required_material`, and `PRAGMA foreign_key_check` is empty. A
  link-free seed alone does not satisfy this criterion. (VAL-005)
- [ ] A crafted legacy archive at `ProjectionSchema` 14 carrying a
  `review.recorded` Event is rejected by `validateEvent`, and the same archive
  at the current projection schema is not rejected for that reason. (No VAL ID:
  VAL-065 correctly does not cover this - the legacy branch at
  `archive_validation.go:127-129` admits profile-14 archives by design, which is
  why `semanticPhase22Events` has to do the work. I-35, DEC-234.)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit, phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

`TestMigration0017CreatesExactlyTheAcceptedObjects` fails against the absent
migration. Capture the failure before product edits, then implement the
smallest behavior that satisfies it. A build failure alone does not prove the
intended rejection.

## Documentation Updates

- `project-docs/specs/03_TRD_SCHEMA.md`: add section 8.4 incorporating
  `schema.md` by reference, per the accepted contract.
- `project-docs/docs/phases/phase-23/tasks/task-1-done.md`: implementation,
  red/green, every exact command with output token and observed termination,
  assertion evidence, convergence provenance, and deviations.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory from this task.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -timeout 30m ./cmd/fid/integration/...
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/
```

Two gates run for minutes and are explicitly bounded (DEC-227): the `-timeout
30m` above replaces Go's implicit 10-minute per-binary default, which the tagged
suite already approaches at 560 s. Record in the completion report the wall time
`go test` reports for `cmd/fid/integration`, `cmd/fid`, and `internal/workflow`,
so a regression is attributed to the task that caused it. Budgets: tagged
integration under 900 s and no ordinary package over 600 s (VAL-079). A breach is
a bounded speed repair owned by this task when this task's own new journeys
caused it; otherwise stop and escalate. Never record a timed-out gate as `PASS`.

No exception to the full task gate. Commands whose expected pass is non-zero:
none. Focused red/green selections never replace these commands. Report each
attempt with the command's own output token and observed termination, as
required by `AGENTS.md`. An interrupted or hung command is never `PASS`;
preserve all failed attempts and rerun the byte-identical command after an
authorized causal repair. Report lint is a Runner control after the report
commit, not a Developer gate.

## Notes For Developer

- This task belongs to QA Review window 1.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-1`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.

## ESC-23-002 recovery scope amendment (DEC-238)

The Human's 2026-09-20 delegation assigns this recovery the rebuild statement
reuse and cached execution-transition fixture repair described in DEC-238,
plus causal test-runtime repair in `cmd/fid` and `internal/workflow` if needed.
Matching tests and repair evidence are in scope. All exact commands and
thresholds above remain unchanged. This authorizes causal cross-task repair;
it does not authorize weakening assertions or dispatching later tasks.
