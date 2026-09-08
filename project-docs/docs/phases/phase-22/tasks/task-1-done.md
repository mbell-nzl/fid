# Phase 22 - Task 1: Recovery schema and replay foundation - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Replaced the accepted development-baseline CHECK values in migration 0013 and added the exact schema-15 recovery up/down migrations.
- Added fresh-init, full down/up, constraint, live-apply, and replay coverage for recovery tables and rejected Effects without enabling recovery domain operations.
- Centralized active migration reporting at schema 15 for rebuild and export while preserving the separate projection-freshness schema value.
- Converged retained Phase 21 schema and Attention-boundary assertions displaced by the accepted Phase 22 schema surface.

## Files created or modified

- `migrations/0013_p21_operational_kernel.up.sql`
- `migrations/0015_p22_recovery.up.sql`
- `migrations/0015_p22_recovery.down.sql`
- `internal/db/projection_meta.go`
- `internal/db/operational_kernel_migration_test.go`
- `internal/db/recovery_migration_test.go`
- `internal/journal/project_effects.go`
- `internal/journal/recovery_schema_test.go`
- `internal/export/export.go`
- `internal/export/export_test.go`
- `cmd/fid/rebuild.go`
- `cmd/fid/rebuild_test.go`
- `cmd/fid/contract_surface_test.go`
- `cmd/fid/integration/task22_recovery_schema_test.go`
- `project-docs/docs/phases/phase-22/tasks/task-1-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-1-done.md`

## Validation Contract Coverage

- VAL-001: PASS - `TestMigration0015RecoverySchemaConstraintsAndRoundTrip` in `internal/db/recovery_migration_test.go:9` verifies all recovery tables, the non-null `run_recoveries.observation_id` column at `:69`, every accepted Stop/Attention code and both invalid direct inserts at `:94`; `TestRecoverySchemaFreshBinaryInitAndRebuild` in `cmd/fid/integration/task22_recovery_schema_test.go:13` proves the built binary installs the recovery schema.
- VAL-002: PASS - `TestMigration0015RecoverySchemaConstraintsAndRoundTrip` runs full `MigrateDown` then `MigrateUp` at `internal/db/recovery_migration_test.go:32`, verifies the same tables at `:43`, and proves every recovery table is empty at `:121`.
- VAL-003: PASS - `TestRecoverySchemaFreshInitAndReplay` in `internal/journal/recovery_schema_test.go:14` rejects the unavailable `run.recovered` Effect through live apply at `:34` and replay at `:49`; it also rejects an invalid Stop code in live apply at `:60` and replay at `:76`.

## Tests run

- `go test -race ./internal/journal -run '^TestRecoverySchemaFreshInitAndReplay$' -count=1` - RED - `FAIL` (terminated `exit_code: 1`); the first attempt was a test-harness compile failure: `internal/journal/recovery_schema_test.go:40:16: undefined: openProjectDB`.
- `go test -race ./internal/journal -run '^TestRecoverySchemaFreshInitAndReplay$' -count=1` - RED - `FAIL` (terminated `exit_code: 1`); after correcting only the test harness, the named test produced `run_recoveries table absent after fresh init` before product edits.
- `go test -race ./internal/journal -run '^TestRecoverySchemaFreshInitAndReplay$' -count=1` - GREEN - `ok  	github.com/mbell-nzl/fid/internal/journal	2.670s` (terminated `exit_code: 0`).
- `go test -race ./internal/db -run '^TestMigration0015RecoverySchemaConstraintsAndRoundTrip$' -count=1` - RED - `FAIL` (terminated `exit_code: 1`); before the migration, `table "run_runtime" does not exist`.
- `go test -race ./internal/db -run '^TestMigration0015RecoverySchemaConstraintsAndRoundTrip$' -count=1` - GREEN - `ok  	github.com/mbell-nzl/fid/internal/db	3.751s` (terminated `exit_code: 0`).
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestRecoverySchemaFreshBinaryInitAndRebuild$' -count=1` - GREEN - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.352s` (terminated `exit_code: 0`).
- `/tmp/fid-phase22-task1 --project /tmp/fid-phase22-task1.Sw07ej --json init --prefix P22` - GREEN - `{"contract_version":"fid-cli.v1","event":{"id":"EVT-20260905-000001","occurred_at":"2026-09-05T02:05:09.402846Z"},"ok":true,"operation":"project.init","result":{"journal_path":".fid/journal/2026-09-05.jsonl","project":{"id":"project-252b9f9aeebb071de1f086771f08c54f","issue_prefix":"P22","root":"/tmp/fid-phase22-task1.Sw07ej"}}}` (terminated `exit_code: 0`).
- `/tmp/fid-phase22-task1 --project /tmp/fid-phase22-task1.Sw07ej --json rebuild` - GREEN - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"project.rebuild","result":{"events_replayed":1,"files_replayed":1,"schema_version":15,"upgraded_events":0}}` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
  - Attempt 1.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
  - Attempt 2.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
  - Attempt 3.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
  - Attempt 1.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
  - Attempt 2.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
  - Attempt 3.
- `go test -race -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`).
  - Attempt 1.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	35.791s	coverage: 27.4% of statements` (terminated `exit_code: 0`)
  - Attempt 2.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 27.4% of statements` (terminated `exit_code: 0`)
  - Attempt 3.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)				61.2%` (terminated `exit_code: 0`)
  - Attempt 1.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)				61.2%` (terminated `exit_code: 0`)
  - Attempt 2.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.145s` (terminated `exit_code: 0`)
  - Attempt 1.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.933s` (terminated `exit_code: 0`)
  - Attempt 2.
- `go version` - PASS - `go version go1.26.2 linux/amd64` (terminated `exit_code: 0`)
  - Attempt 1.
- `go version` - PASS - `go version go1.26.2 linux/amd64` (terminated `exit_code: 0`)
  - Attempt 2.
- `go env GOOS GOARCH` - PASS - `amd64` (terminated `exit_code: 0`)
  - Attempt 1.
- `go env GOOS GOARCH` - PASS - `amd64` (terminated `exit_code: 0`)
  - Attempt 2.
- `uname -srm` - PASS - `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64` (terminated `exit_code: 0`)
  - Attempt 1.
- `uname -srm` - PASS - `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64` (terminated `exit_code: 0`)
  - Attempt 2.
- `nproc` - PASS - `6` (terminated `exit_code: 0`)
  - Attempt 1.
- `nproc` - PASS - `6` (terminated `exit_code: 0`)
  - Attempt 2.
- `grep -m1 '^model name' /proc/cpuinfo` - PASS - `model name	: AMD Ryzen 7 5700G with Radeon Graphics` (terminated `exit_code: 0`)
  - Attempt 1.
- `grep -m1 '^model name' /proc/cpuinfo` - PASS - `model name	: AMD Ryzen 7 5700G with Radeon Graphics` (terminated `exit_code: 0`)
  - Attempt 2.
- `grep -m1 '^MemTotal:' /proc/meminfo` - PASS - `MemTotal:        8133236 kB` (terminated `exit_code: 0`)
  - Attempt 1.
- `grep -m1 '^MemTotal:' /proc/meminfo` - PASS - `MemTotal:        8133236 kB` (terminated `exit_code: 0`)
  - Attempt 2.
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	31.083s` (terminated `exit_code: 0`)
  - Attempt 1; sample `28341 ms/rebuild`, below the unchanged `<30000 ms` assertion.
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	30.587s` (terminated `exit_code: 0`)
  - Attempt 2; sample `27805 ms/rebuild`, below the unchanged `<30000 ms` assertion.

## Deviations from spec

- No contract or command definition changed. The first required race/coverage attempt failed with `FAIL` (terminated `exit_code: 1`) because the accepted schema-15 changes invalidated two retained Phase 21 assertions: `TestDecideProjectionFreshnessRequiresPhase19SchemaVersion` had frozen the separate projection-freshness value, and `TestPhase21AttentionBoundary` prohibited any Attention schema/applier surface. The causal repair preserved projection freshness at 5, introduced a separate active migration value of 15 for rebuild/export, and constrained the retained Attention inventory to its continuing CLI/workflow no-domain boundary. The byte-identical later attempts passed with `ok  	github.com/mbell-nzl/fid/internal/workflow	35.791s	coverage: 27.4% of statements` and then the final cached workflow line, both terminated `exit_code: 0`.
- Convergence provenance: `internal/db/migrate_v01_tables_test.go` was introduced by `b63a3cd3f237d48990404d6b9c7099db69ce998e` (`test(convergence): finalize task-20 V1 fixture removal`, Refs: phase-21/task-20); `cmd/fid/contract_surface_test.go` lines 70-98 were introduced by `7277139ee5d5dde4cd3627335d39b14f616f46a1` (`feat(evidence): evaluate trusted receipt evidence`, Refs: phase-21/task-17). The accepted task-1 schema and unavailable-Event foundation displaced their older assertions. Both are executed by `go test -race -coverprofile=coverage.out ./...`.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `7647abc66d0e309b7ab056d05e80d1c501d758dd`

## Commit subject

- `fix(journal): validate current recovery Event transactions`

## Final-QA FQ-01 same-task continuation — 2026-09-09

Authority: review-7, final-review.md FQ-01, task-1.md, schema.md §§5–6 and
DEC-192 execution-scope convergence. Prior primary binding:
72c5974f6698fa08a4fe1a24099b3a93b3882f1d,
`feat(recovery): add schema and replay foundation`. All original report and
failed-attempt evidence remains above. This continuation implements the
previously undelivered current-family validation under the same task.

Root owned the repair and verification. One bounded worker changed only the
projection fixture file while root read without editing; root then resumed as
sole writer and ran the integrated evidence below. No implementation/reviewer
repair cycle or concurrent product writers were used for this continuation.

### Delivered behavior and assertion coverage

VAL-003 / QA-03, QA-79: both direct ProjectEvent and ReplayJournal now enter
the same pre-mutation current-family validator inside the existing projection
transaction. Its finite catalog covers the 22 Phase 22 replacement/extension
and incorporated activation families. It rejects unrelated or reordered
Effects, missing required stages, excess singular Effects, incomplete insert
rows, unknown columns, wrong scalar domains, non-one initial versions,
incorrect update fields/guards, missing or incorrect captured Prior, and
incorrect entity/Event/idempotency references. Current-domain tables cannot
be introduced under an unrelated retained family. Retained historical writers
are not assigned new generic policy by this repair; the existing pre-recovery
projection profile keeps its separately accepted grammar.

CommitEvent/CommitIdempotent capture missing producer-side update Prior values
under the existing write lock before dry-run/append. Received Events never get
manufactured Prior data. Captured DATETIME values retain exact stored text and
Boolean columns retain Boolean JSON values, preserving independent portable
interpretation as well as SQLite projection. Checks compare exact values,
not floating-point tolerances or reduced timestamp precision.

Bindings include predecessor/successor Run, Claim, Observation and preserved
Issue/Envelope ownership; runtime and Budget ownership; Receipt/Attempt start
and Artifact links; activation primary/prior identity, selected/completed
Event and completion direction; and ordered Stop/new-or-reused Attention/link
bundles. Stop subjects must exist, and standalone dispositions reference an
existing original idempotent Event result. These checks validate stored
bindings without replaying contemporary authorization or repository policy.
Run runtime version is deliberately not equated with Run aggregate sequence:
Attempt/Receipt Events also advance runtime under their own aggregates.

The original checkpoint relabels remain negative cases. Their positive tests
now construct actual recovered, budget-observed and execution-blocked Effects.
Closure and activation fixtures retain their original subjects with complete
rows, captured state and linked disposition facts. The execution-blocked
fixture now uses its accepted unknown-budget Stop. Conditional recovery
controls cover all four active/terminal predecessor and active/expired Claim
combinations, checking unchanged prior terminal state and full live/replay
projection equality. A migration-derived column/primary-key/nullability test
guards the static current column catalog against the accepted DDL.

TestP22CurrentEventFamiliesLiveReplayValidation obtains its input from real
built CLI journeys, including recovery, closure, abandonment, Attention reuse,
first activation, ordinary abort, and upgrade recovery in both directions.
It requires all 22 named current families. Every original Event is accepted;
mutations exercise unrelated families, missing/reordered domain Effects,
unknown columns, wrong scalar/immutable entity/result identities, initial
versions, missing captured Prior, missing Stop subjects and missing resolution
Events. Both live and replay rejection compare every table before/after.
Replay sequence-gap tests use the genuine preceding journal bytes with an
explicit StartOffset, so they reach the target Event rather than failing on
an already-applied prefix. Final fresh replay matches the entire incremental
live projection. The existing full actual recovery/activation/budget fault,
portable import/export, and executor tests remain in the required gates.

VAL-001/002 retain the original accepted schema work; this continuation changes
no migration, DDL, schema version, assertion, command, threshold or policy.
No dependency was added. No held-out scenario material was read.

### Causally displaced fixture provenance

| Surface | Introducing accepted owner/commit | Subject-preserving convergence |
|---|---|---|
| checkpoint_projection_test.go checkpoint | Phase 22/task-4, 59c62b9e1d82c5f698d1e173908446d35b64f539 | Complete captured runtime Prior/idempotency while preserving checkpoint parity |
| recovered projection | Phase 22/task-5, b7dce8e747a3f2d720002faee3db7f95b3cb4019 | Genuine successor/recovery rows replace the invalid relabel; malformed original retained as negative |
| budget/blocked projection | Phase 22/task-7, a37e7d9a18fcb8d19432e1ce5100a33db839f978 | Genuine sample/Stop transactions replace invalid relabels; negatives retained |
| closure projection | Phase 22/task-8, b965b24748f57a0f5e42b8dd3054dfe6fd5dbb74 | Exact Prior, idempotency and existing closure/Claim state assertions |
| activation projection | Phase 22/task-10, 33f85bbd6abffd0425a6a446780e8b0a37ce457a | Full version/time updates, five resources and actual linked Stop disposition |
| idempotency_test.go / event_v2_test.go | Phase 21/task-2, 94e752d25809a0e54f827b1872172f8038c0a741 | Complete current prepare fixture; same replay, concurrency, conflict and sequence-gap assertions |
| archive_semantics_test.go missing-reference fixture | Phase 22/task-13, acb569bb9447a3259991520669167fdc819c1706 | Bind idempotency so the invalid prior Run reference reaches its intended rejection; retain reference-invalid classification |
| stop_test.go helper | Phase 22/task-2, ac827959f8359473cbd6546ac6bea1ff1fd20751 | Explicit fixture idempotency and existing Run subject; same immutable Stop/shared Attention tests |
| stop_read_test.go / evaluation_test.go | Phase 22/task-9, a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460 | Seed Stop-read subject and bind prepared Event ID in result; same read-inert and Envelope-stop assertions |
| attention_test.go shared disposition fixture | Phase 22/task-2 FQ-02, f973bca7ea146f00c0e16ca5e5117a8d13f769b9 | Complete second Stop's idempotency/containing identity; retain two-Stop exact disposition and invalid underlying-fact cases |
| p22_final_journey_test.go | Phase 22/task-18, 24e3cbdb22dd9e43b55636da313dcef1ef633b6f | Extract returning fixture helper without removing any journey assertion, allowing its genuine emitted Events to feed the new validator corpus |

All introducing hashes/subjects/Refs were observed with git log. The ordinary
race gate executes the in-package converged rows; the dedicated integration
gate executes the actual final journey and new all-family corpus. New validator
helpers/tests and p22_event_validation_test.go are task-1 implementation.

### Failed attempts, diagnosis and repair

The original red admitted malformed checkpoint relabels as recovery, budget
and execution-blocked Events. Successive focused red cases demonstrated
missing/corrupt Prior, idempotency identity, conditional Stop/link and recovery
bindings, initial version >1, immutable Attempt ID mismatch and absent Stop
subjects. The exact commands and their retained tokens/terminations follow;
their later successful executions close these cases without weakening them.

The first broader convergence run exposed accepted producer forms omitted by
the draft validator: force-recover command, bare Stop request digest,
Receipt's retained idempotency header/update, first-Stop denial result identity,
and Claim expiry's actor guard. Those were reconciled to the existing producers
and contracts. Retained Issue-label spellings and historical scale fixtures
were kept outside new generic historical semantics. Incomplete projection and
idempotency fixtures were converged with the provenance above.

Temporary implementation/fixture failures are preserved: duplicate generated
catalog keys caused a compile failure; DATETIME prior values initially had the
wrong scalar; Boolean prior encoded zero instead of false and failed the
independent portable model; the new ordinary-abort corpus input initially
omitted its required expected_version. Each was repaired causally. Initial
mutation-test failures can also show subsequent changed-state errors after
the first invalid Event was accepted; those are retained consequences, not
independent passing proofs.

The first required lint failed with `2 issues:` / `* errcheck: 2`
(terminated exit_code: 1): two cleanup returns were unchecked. Explicit
discard in deferred cleanup, consistent with existing local cleanup patterns,
fixed them; the identical command later printed `0 issues.` and exited 0.

The first full required race/coverage command failed when the 3.9-GB /tmp
tmpfs filled: linker `no space left on device` and SQLite `database or disk
is full (13)` failures occurred, ending `FAIL` (terminated exit_code: 1).
111 inactive temporary integration-binary directories, 1,726,563,573 bytes,
were moved to ignored project storage with links preserving their original
paths. This restored about 1.8 GB free without deleting evidence, altering
the test environment or changing any gate command. The failed process was
allowed to terminate. Final reference probes were completed while that
already-failed attempt finished, then the identical full gate was rerun on
the final source. The relocation is housekeeping, not a gate PASS.

The first full tagged integration attempt ended `FAIL` (terminated
exit_code: 1, package409.725s): all three Receipt budget-fault variants
reported `current Event gate.receipt-recorded: missing or reordered Effect stage 2`.
The new validator had ordered budget Effects after the Receipt row, while
the already accepted producer appends budget Effects to its Artifact prefix
before receiptEvent appends the Receipt and links. Introducing producer:
a3766804f1a40ca5e827748f59550e1d72d452aa,
`feat(budget): persist gate execution charges`, Refs: phase-22/task-7,
observed with git log -S. Correcting the new catalog's order preserves that
producer, the schema extension and every fault assertion. The exact focused
nine-case budget test then passed in33.745s (exit_code: 0). The complete
twelve-command sequence was restarted after this last source correction;
the unchanged full tagged command's later receipt below governs closure.

### Own deterministic command receipts

Working directory for every command below: /home/mark/projects/fid.
Results are from this root session's returned command objects, including
observed nested termination. Prior historical receipts above remain intact.
- `go test -race ./internal/journal -run '^TestCurrentEventFamiliesRejectMalformedEffectsLiveAndReplay$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 1: current Event family corruption red.
- `go test -race ./internal/journal -run '^TestCurrentEventFamiliesRejectMalformedEffectsLiveAndReplay$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 2: current Event initial grammar repair.
- `go test -race ./internal/journal -run '^TestCurrentEventFamiliesRejectMalformedEffectsLiveAndReplay$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	34.364s` (terminated `exit_code: 0`). Attempt 3: current Event catalog compilation repair.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^(TestP22CloseRunUsesObservedRepositoryAndQualifyingProof|TestP22ActivationRecoverySurvivesAppendProjectionInterruption)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	56.208s` (terminated `exit_code: 0`). Attempt 4: current Event genuine producer controls initial.
- `go test -race ./internal/journal -run '^TestCurrentEventRequiresExactCapturedPrior$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 5: captured prior red.
- `go test -race ./internal/journal -run '^TestCurrentEventRequiresExactCapturedPrior$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	4.514s` (terminated `exit_code: 0`). Attempt 6: captured prior green.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^(TestP22CloseRunUsesObservedRepositoryAndQualifyingProof|TestP22ActivationRecoverySurvivesAppendProjectionInterruption)$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 7: genuine producer controls with captured prior.
- `go test -race ./internal/workflow -run '^TestActivationRecoverySelectionRaisesAndCompletionDisposesIncompleteStop$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 8: diagnose producer prior failure at domain boundary.
- `go test -race ./internal/workflow -run '^TestActivationRecoverySelectionRaisesAndCompletionDisposesIncompleteStop$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	3.731s` (terminated `exit_code: 0`). Attempt 9: producer prior timestamp representation repair.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^(TestP22CloseRunUsesObservedRepositoryAndQualifyingProof|TestP22ActivationRecoverySurvivesAppendProjectionInterruption)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	54.394s` (terminated `exit_code: 0`). Attempt 10: genuine producer controls after exact timestamp prior repair.
- `go test -race ./internal/journal -run '^TestCurrentEventFamiliesRejectMalformedEffectsLiveAndReplay$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 11: complete fixture and idempotency binding red.
- `go test -race ./internal/journal -run '^TestCurrentEventFamiliesRejectMalformedEffectsLiveAndReplay$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	40.780s` (terminated `exit_code: 0`). Attempt 12: current idempotency binding green.
- `go test ./internal/journal ./internal/workflow ./internal/evidence ./internal/export ./cmd/fid` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 13: current Event broader retained producer and fixture convergence check.
- `go test ./internal/journal ./internal/workflow ./internal/evidence ./internal/export ./cmd/fid` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 14: Converged accepted force/Receipt/Stop/retained grammar and test producers.
- `go test ./internal/journal -run '^TestCurrentEventRejectsConditionalBindingCorruption$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 15: Conditional binding and cardinality red.
- `go test ./internal/journal -run '^TestCurrentEventRejectsConditionalBindingCorruption$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	0.487s` (terminated `exit_code: 0`). Attempt 16: Conditional binding green attempt.
- `go test ./internal/journal ./internal/workflow ./internal/evidence ./internal/export ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	51.228s` (terminated `exit_code: 0`). Attempt 17: Complete cardinality/bindings and accepted producer corrections.
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CurrentEventFamiliesLiveReplayValidation$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 18: All current families genuine CLI corpus and live/replay mutation matrix.
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CurrentEventFamiliesLiveReplayValidation$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 19: Preserve typed boolean captured prior in actual portable replay.
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CurrentEventFamiliesLiveReplayValidation$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 20: All-family corpus with initial-version and stale-prior mutations.
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CurrentEventFamiliesLiveReplayValidation$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	19.039s` (terminated `exit_code: 0`). Attempt 21: Inserted version exact-one and complete abort fixture.
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CurrentEventFamiliesLiveReplayValidation$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 22: Closed immutable entity identity mutation red.
- `go test ./internal/journal -run '^(TestCurrentEvent|TestCurrentRecovery)' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	4.272s` (terminated `exit_code: 0`). Attempt 23: Final focused current-family conditional and real CLI proof.
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CurrentEventFamiliesLiveReplayValidation$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	20.014s` (terminated `exit_code: 0`). Attempt 24: Final focused current-family conditional and real CLI proof.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`). Attempt 25: Task1 exact gate final-QA repair.
- `golangci-lint run ./...` - BLOCKED - `* errcheck: 2` (terminated `exit_code: 1`). Attempt 26: Task1 exact gate final-QA repair.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`). Attempt 27: Exact gate after checked-cleanup lint repair.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`). Attempt 28: Exact gate after checked-cleanup lint repair.
- `go test -race -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 29: Task1 exact full race and coverage gate.
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CurrentEventFamiliesLiveReplayValidation$' -count=1` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 30: Final Stop subject and historical fact-reference red.
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CurrentEventFamiliesLiveReplayValidation$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	26.583s` (terminated `exit_code: 0`). Attempt 31: Closed Stop subject and fact references implemented.
- `go test ./internal/journal ./internal/workflow ./internal/evidence ./internal/export ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	53.814s` (terminated `exit_code: 0`). Attempt 32: Final reference conformance broad check.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`). Attempt 33: Final exact gate after reference repair and tmpfs restoration.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`). Attempt 34: Final exact gate after reference repair and tmpfs restoration.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	510.142s	coverage: 62.1% of statements` (terminated `exit_code: 0`). Attempt 35: Unchanged full race/coverage rerun after restoring tmpfs and final reference repair.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					65.6%` (terminated `exit_code: 0`). Attempt 36: Task1 final exact gate sequence.
- `go test -tags integration -race ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Attempt 37: Task1 final exact gate sequence.
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestP22BudgetBoundarySurvivesAppendAndProjectionFaults$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	33.745s` (terminated `exit_code: 0`). Attempt 38: Receipt budget stage causal repair focused verification.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`). Attempt 39: Final exact twelve-command sequence after Receipt ordering repair.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`). Attempt 40: Final exact twelve-command sequence after Receipt ordering repair.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	516.271s	coverage: 62.1% of statements` (terminated `exit_code: 0`). Attempt 41: Final exact twelve-command sequence after Receipt ordering repair.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					65.6%` (terminated `exit_code: 0`). Attempt 42: Final exact twelve-command sequence after Receipt ordering repair.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	409.209s` (terminated `exit_code: 0`). Attempt 43: Final exact twelve-command sequence after Receipt ordering repair.
- `go version` - PASS - `go version go1.26.2 linux/amd64` (terminated `exit_code: 0`). Attempt 44: Final exact twelve-command sequence after Receipt ordering repair.
- `go env GOOS GOARCH` - PASS - `amd64` (terminated `exit_code: 0`). Attempt 45: Final exact twelve-command sequence after Receipt ordering repair.
- `uname -srm` - PASS - `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64` (terminated `exit_code: 0`). Attempt 46: Final exact twelve-command sequence after Receipt ordering repair.
- `nproc` - PASS - `6` (terminated `exit_code: 0`). Attempt 47: Final exact twelve-command sequence after Receipt ordering repair.
- `grep -m1 '^model name' /proc/cpuinfo` - PASS - `model name	: AMD Ryzen 7 5700G with Radeon Graphics` (terminated `exit_code: 0`). Attempt 48: Final exact twelve-command sequence after Receipt ordering repair.
- `grep -m1 '^MemTotal:' /proc/meminfo` - PASS - `MemTotal:        8133236 kB` (terminated `exit_code: 0`). Attempt 49: Final exact twelve-command sequence after Receipt ordering repair.
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	30.393s` (terminated `exit_code: 0`). Attempt 50: Final exact twelve-command sequence after Receipt ordering repair.

Final reference result: `BenchmarkRebuildWorkspace100kEvents-6 1 27598548704 ns/op 27599 ms/rebuild 2292273456 B/op 34969623 allocs/op`, followed by `PASS` and the package result above (terminated `exit_code: 0`). Self-asserted rebuild parity and <30000 ms. Reference instance matches the six values in the accepted task; go env returned linux and amd64.

Own raw returned command objects: `.cache/phase-22-fq01-46gqnetx/command-results.json`. All50 command occurrences above terminated; final exact12 are PASS. The numbered occurrences are command receipts, not recovery attempt counts. This remains recovery attempt1 under the ledger signature.

### Final-QA repair files

- `cmd/fid/integration/p22_event_validation_test.go`
- `cmd/fid/integration/p22_final_journey_test.go`
- `internal/export/archive_semantics_test.go`
- `internal/journal/checkpoint_projection_test.go`
- `internal/journal/current_event_activation.go`
- `internal/journal/current_event_cardinality.go`
- `internal/journal/current_event_catalog.go`
- `internal/journal/current_event_columns.go`
- `internal/journal/current_event_idempotency.go`
- `internal/journal/current_event_prior.go`
- `internal/journal/current_event_references.go`
- `internal/journal/current_event_stop_references.go`
- `internal/journal/current_event_updates.go`
- `internal/journal/current_event_validation.go`
- `internal/journal/current_event_validation_test.go`
- `internal/journal/event_v2_test.go`
- `internal/journal/idempotency.go`
- `internal/journal/idempotency_test.go`
- `internal/journal/project_effects.go`
- `internal/workflow/attention_test.go`
- `internal/workflow/evaluation_test.go`
- `internal/workflow/stop_read_test.go`
- `internal/workflow/stop_test.go`

FQ-01 implementation and all exact task gates are complete. The fresh independent final QA review must determine finding closure; no Human assertion is marked passed by this report.
