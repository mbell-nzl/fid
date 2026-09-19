# Phase 23 - Task 1: Review and judgment schema baseline - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: session-continuation
- Condition: The required 100k rebuild benchmark failed on the final accepted rebuild behavior, and the final repository-wide race gate also failed in an out-of-scope `cmd/fid` execution-transition assertion.
- Recovery attempts: 1
- Next action: Runner triage must preserve the benchmark and race-gate failed attempts, establish a bounded repair authority if a byte-preserving rebuild performance fix exists, and resolve the out-of-scope `cmd/fid` assertion before this task can rerun the full exact gate.
- Human attention required: no

## Summary of what was implemented

- Added migration 0017 and its link-preserving `evidence_records` rebuild, the seven accepted 0013/0015 baseline changes, schema version 17, all nine progressive-fence registrations, archive semantics, and integration-child binary reuse.
- Added acceptance tests for exact migration DDL/object inventory, the obligation discriminator, the seven-hunk baseline diff, receipt evidence plus artifact-link in-place preservation, and all nine fenced Event types.

## Files created or modified

- `migrations/0013_p21_operational_kernel.up.sql`
- `migrations/0015_p22_recovery.up.sql`
- `migrations/0017_p23_review_judgment.up.sql`
- `migrations/0017_p23_review_judgment.down.sql`
- `internal/db/`, `internal/journal/`, `internal/export/`, and `cmd/fid/integration/main_test.go`
- `project-docs/specs/03_TRD_SCHEMA.md`
- `project-docs/docs/phases/phase-23/tasks/task-1-done.md`

## Documentation updated

- `project-docs/specs/03_TRD_SCHEMA.md` adds section 8.4, incorporating the accepted Phase 23 schema authority.

## Validation Contract Coverage

- VAL-001: `TestMigration0017CreatesExactlyTheAcceptedObjects` in `internal/db/review_judgment_migration_test.go` migrates 16 to 17, requires exactly 29 newly named schema objects, checks rebuilt tables, runs `foreign_key_check` and `integrity_check`, and byte-compares both 0017 SQL files with schema.md sections 2 and 3.
- VAL-002: `TestMigration0017ObligationRegistrationDiscriminator` proves both rejected rows (gate without registration; review with registration) and both accepted rows (gate with registration; human outcome without registration).
- VAL-003: `TestProjectEventFenceRejectsEveryUnavailablePhase23EventWithoutProjectionChange` in `internal/journal/project_effects_test.go` iterates all nine fence keys, requires the exact unavailable Event envelope, and checks the projection sentinel before and after each rejection.
- VAL-004: `TestMigration0017CreatesExactlyTheAcceptedObjects` proves fresh schema-17 pragma checks; `TestPristineWorkspaceClonePreservesRebuildParityAndIsolation` passed after the rejected performance experiment was removed. The full-gate rebuild confirmation remains blocked by the benchmark and unrelated race-gate failure.
- VAL-005: `TestMigration0017PreservesReceiptEvidenceAndArtifactLinksInPlace` seeds schema 16 with a receipt Evidence row and `evidence_artifact_links.required_material = 1`, migrates in place, asserts each legacy Evidence value and null new references, restores the exact link row, runs `foreign_key_check`, and requires `schema_migrations = 17|0`.
- VAL-075: `TestPhase23BaselineMigrationDiffIsExactlySevenAcceptedHunks` compares the two migration files against planning base `0195897a347dfe1872db0f7be3800c5c3ccdefcf` and accepts only the seven declared changes.

## Tests run

- `go test ./internal/db -run '^TestMigration0017CreatesExactlyTheAcceptedObjects$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); initial red against the absent migration 0017.
- `go test ./internal/db -run '^TestMigration0017CreatesExactlyTheAcceptedObjects$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/db 0.063s` (terminated `exit_code: 0`).
- `go test ./internal/db -run '^(TestMigration0017CreatesExactlyTheAcceptedObjects|TestMigration0017ObligationRegistrationDiscriminator|TestMigration0017PreservesReceiptEvidenceAndArtifactLinksInPlace|TestPhase23BaselineMigrationDiffIsExactlySevenAcceptedHunks)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/db 0.214s` (terminated `exit_code: 0`).
- `go test -race ./internal/journal -run '^(TestProjectEventFenceRejectsEveryUnavailablePhase23EventWithoutProjectionChange|TestCurrentEventColumnDomainsMatchAcceptedDDL)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/journal 4.065s` (terminated `exit_code: 0`).
- `go test -race ./cmd/fid -run '^TestPristineWorkspaceClonePreservesRebuildParityAndIsolation$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 5.922s` (terminated `exit_code: 0`); this is the retained projection-byte-parity control after reverting the rejected exclusive-lock experiment.
- Real integration entrypoint: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` invoked the integration `TestMain`, which built the parent `fid` executable and re-executed the child suite through `FID_INTEGRATION_BINARY`; it passed with no residual integration-binary directories.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`).
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 595.167s coverage: 62.6% of statements` (terminated `exit_code: 0`); prior successful attempt. Wall times: `cmd/fid` 581.395s; `internal/workflow` 595.167s.
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.5%` (terminated `exit_code: 0`); prior successful attempt.
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 668.694s` (terminated `exit_code: 0`). Before the gate, `/tmp` had 0 `fid-integration-binary-*` directories and 3789979648 free bytes; after it had 0 directories and 3789045760 free bytes. Observed `cmd/fid/integration` wall time: 668.694s.
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - BLOCKED - `FAIL github.com/mbell-nzl/fid/internal/journal 35.245s` (terminated `exit_code: 1`); first final-behavior attempt measured 32.500159404s, above the strict 30s assertion.
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`); latest mandatory attempt after acceptance-test additions. `cmd/fid` failed after 619.168s and `internal/workflow` passed after 614.196s, both above the 600s budget. The `cmd/fid` failure includes `TestExecutionBoundTransitionPredicateMatrix/workflow_edge_invalid`, whose expected only `workflow-transition-invalid` but observed `claim-expired`, `run-inactive`, and `workflow-transition-invalid`.
- Later exact commands are BLOCKED pending recovery because command 3 did not terminate green; no narrower command is substituted for the required full gate.

## Deviations from spec

- The first repository-wide race attempt failed before retained schema-16 expectations were converged: it showed `active migration = 17, want16` and the old 0013 inventory. The accepted schema-17 replacement caused that failure; in-scope retained expectations were converged and the focused packages passed.
- Benchmark first failure: `FAIL github.com/mbell-nzl/fid/internal/journal 35.245s` (terminated `exit_code: 1`), root cause a 32.500159404s 100k rebuild after the new schema baseline. A `temp_store(MEMORY)` experiment measured 33.156833747s and was removed. A relaxed migration-open experiment measured 32.553173739s but broke the staged rebuild byte-parity control and was removed. An exclusive-lock experiment measured 27.170 ms/rebuild but broke `TestPristineWorkspaceClonePreservesRebuildParityAndIsolation`; restoring normal locking did not repair the byte difference, so the experiment was removed. The retained final behavior has no safe accepted repair in this attempt.
- Latest full-gate failure is preserved. `cmd/fid/task12_execution_transition_test.go:217` was introduced by `0ced511b` (`test(workflow): cover execution-bound predicate matrix`), which predates the Phase 23 planning base. Task 1 does not own that file, the execution-bound predicate, or its expected transition-reason semantics; no local repair was made.

## Concerns, technical debt, or notes for the Tech Lead

- Runner triage is required for the unresolved rebuild budget and the out-of-scope `cmd/fid` transition assertion before this task can close. The benchmark assertion and byte-parity control were preserved without weakening either.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `3c0b520e2d183c203356942d00a3652aef84376f`

## Commit subject

- `test(phase-23): add task-1 migration acceptance evidence`
