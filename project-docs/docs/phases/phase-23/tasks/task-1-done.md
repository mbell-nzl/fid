# Phase 23 - Task 1: Review and judgment schema baseline - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: Accepted schema baseline and all exact Task 1 gates pass after the Human-delegated DEC-238 causal repair; ESC-23-002 is resolved without changing any threshold or assertion.
- Recovery attempts: 3
- Next action: Phase Runner verifies the committed report and reconciles the retained STOP using qa/resume-after-esc-23-002.md, then continues the ordered tasks 2-8; tasks 9-11 remain held.
- Human attention required: no

## Summary of what was implemented

- Migration 0017, the seven accepted 0013/0015 baseline changes, schema version 17, all nine progressive fences, archive semantics, and integration-child binary reuse remain in place with their acceptance tests.
- Full rebuild now reuses at most 128 prepared SQL shapes on one owned connection. Each Event retains its transaction, validation, rollback and checkpoint. The cache and connection close before staged installation; live apply and catch-up retain their existing path.
- The execution-transition fixture caches durable Envelope setup and acquires a fresh real claim and Run for each clone. Expected rejection reasons and the 60-second TTL are unchanged.
- Recovery fixtures and 25 Issue/maintenance setup sites clone the existing real initialized template. Each case retains its own data, commands and assertions; dedicated init and real-process journeys still initialize normally.

## Files created or modified

- Original Task 1 implementation: migrations 0013/0015/0017, `internal/db/`, `internal/journal/`, `internal/export/`, `cmd/fid/integration/main_test.go`, and schema spec section 8.4.
- Recovery production files: `internal/db/projection_meta.go`; `internal/journal/{project.go,project_effects.go,projection.go,replay.go,rebuild_statements.go}`.
- Recovery tests: `internal/journal/rebuild_statements_test.go`; `cmd/fid/{task12_execution_transition_test.go,execution_transition_fixture_test.go,issue_test.go,issue_ac_test.go,issue_get_list_test.go,maintenance_vacuum_test.go,maintenance_integration_test.go}`; `internal/workflow/{recovery_test.go,recovery_fixture_test.go}`.
- TRD, DEC-238, Task 1 scope/report, repair evidence, escalation/triage/current-context records, and Runner handoff. The original run ledger is committed with its STOP retained.

## Documentation updated

- `project-docs/specs/03_TRD_SCHEMA.md` section 8.4 retains the accepted Phase 23 schema authority.
- `project-docs/specs/02_TRD.md` section 6 records bounded rebuild statement reuse and unchanged transaction/resource semantics.
- DEC-238 records the Human's bounded cross-role authority. `qa/esc-23-002-repair.md` records causation, failed attempts, current-session evidence, and direct implementation review. `qa/resume-after-esc-23-002.md` gives the exact resume binding.
- ESC-23-002 is closed in the phase/current registers; earlier TRG/BLOCK/STOP and failed evidence are preserved. Normal independent QA review windows remain required.

## Validation Contract Coverage

- VAL-001: `TestMigration0017CreatesExactlyTheAcceptedObjects` verifies exact object inventory, accepted SQL bytes, schema-17 pragmas and rebuilt tables in the passing full race suite.
- VAL-002: `TestMigration0017ObligationRegistrationDiscriminator` verifies both accepted and rejected discriminator states in that suite.
- VAL-003: `TestProjectEventFenceRejectsEveryUnavailablePhase23EventWithoutProjectionChange` verifies all nine unavailable Events and unchanged projection in that suite.
- VAL-004: schema-17 pragma checks, retained `TestPristineWorkspaceClonePreservesRebuildParityAndIsolation`, deterministic rebuild tests, and new failure rollback/installed-projection preservation tests pass. The real CLI smoke also proves unchanged Issue output/journal bytes and clean FK/integrity checks after rebuild.
- VAL-005: `TestMigration0017PreservesReceiptEvidenceAndArtifactLinksInPlace` verifies seeded receipt Evidence and Artifact-link preservation, null new references, and clean schema state in the full suite.
- VAL-075: `TestPhase23BaselineMigrationDiffIsExactlySevenAcceptedHunks` verifies exactly the accepted baseline changes against planning base in the full suite.
- VAL-079: fresh package timings below satisfy DEC-227. Rebuild's separate strict 30-second assertion also passes in three consecutive fresh benchmark processes.

## Tests run

All new PASS claims below were observed in this repair session, working directory
`/home/mark/projects/fid`. The verbatim earlier report is preserved in
`qa/esc-23-002-evidence/task-1-before-repair.md` as historical evidence; its
prior-session PASS rows are not re-attested here.

- Focused race command recorded verbatim in `qa/esc-23-002-repair.md`: PASS, cmd/fid `33.064s` and journal `14.879s` (terminated `exit_code: 0`). It covers the aged-fixture regression, retained expiry/holder rejection, projection byte parity, Event rollback/connection release, failed staged rebuild preservation, and deterministic rebuild.
- `go test -race ./internal/workflow -run '^TestRecoveryFixturesKeepRepositoryAndProjectionIsolation$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	4.078s` (terminated `exit_code: 0`). Distinct real Git bindings, sibling projection isolation, and pristine template bytes are checked.
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `python3 project-docs/docs/phases/phase-23/qa/esc-23-002-evidence/rebuild-smoke.py` - PASS - `PASS real CLI rebuild: read parity, unchanged journal, schema 17, FK/integrity clean` (terminated `exit_code: 0`). Raw subprocess results are in `qa/esc-23-002-evidence/rebuild-smoke.log`.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	231.010s	coverage: 62.6% of statements` (terminated `exit_code: 0`); fresh CLI **527.099s**, workflow **231.010s**, every package below 600s. Raw output: `qa/esc-23-002-evidence/race-attempt-3.log`.
- `go tool cover -func=coverage.out` - PASS - `65.6%` (terminated `exit_code: 0`); total statement coverage. Raw output: `qa/esc-23-002-evidence/coverage.log`.
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	669.586s` (terminated `exit_code: 0`); below 900s. Binary directories **0 before / 0 after**; `/tmp` free bytes **3696472064 before / 3696328704 after**. Raw output and resource snapshots: `qa/esc-23-002-evidence/integration.log` and `integration-resources.jsonl`.
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - PASS - `PASS` (terminated `exit_code: 0`); three consecutive fresh process runs **22.728s / 22.623s / 22.489s**, all strictly below 30s. Each process terminated `exit_code: 0`; logs `qa/esc-23-002-evidence/benchmark-{1,2,3}.log`.

## Deviations from spec

- DEC-238 authorizes this bounded cross-role and cross-task repair at the Human's explicit request. It changes ownership for this recovery, not the commands, accepted status, semantic assertions, thresholds, dependencies, schema, Events or public interface.
- Original schema-expectation failures and original benchmark `32.500159404s` / `FAIL github.com/mbell-nzl/fid/internal/journal 35.245s` (terminated `exit_code: 1`) remain in the historical report. This session reproduced `100k-event rebuild = 33.049607193s, want <30s` (terminated `exit_code: 1`). Repeated preparation of FK-aware writes dominated the regression; scoped prepared reuse repairs it. The exact rerun passed at `23338 ms/rebuild` and the final three measurements above pass. Earlier pragma/locking experiments are not present.
- The original retained transition assertion failure (`FAIL`, terminated `exit_code: 1`) came from `0ced511b`, predating Phase 23. A deterministic regression reproduced the cached expired claim. Fresh per-clone claims/Runs repair its premise while preserving expiry behavior and expected reasons. The focused rerun and the full exact race gate pass.
- A new failure-preservation test initially failed with `no such file or directory` / `FAIL` (terminated `exit_code: 1`): its helper created directories, not the original DB. Creating and closing that DB before its snapshot repaired setup; the byte-identical focused command then passed.
- This recovery's full race attempts 1 and 2 returned `62.6% of statements` (both terminated `exit_code: 0`) but FAILED the runtime control: CLI/workflow 628.145s/623.362s, then 604.046s/233.338s. Neither is counted as a passing gate. Removing repeated bootstrap setup, supported by before/after measurements and isolation tests, led to the byte-identical fresh attempt 3 above. Raw failed attempts are retained.

- The additional smoke harness first omitted the required idempotency key and then used Python SQLite 3.46.1, which cannot parse the existing dynamic RAISE trigger. Both script invocations terminated `exit_code: 1` and their logs are retained. The corrected invocation uses actual create JSON and the installed SQLite CLI 3.50.6; the byte-identical script command passed without product changes.

## Concerns, technical debt, or notes for the Tech Lead

- Keep 30s as the accepted reference-machine regression target; this repair restores headroom without relaxing it. The measured environment is Linux/amd64 WSL2, Go 1.26.2, six logical CPUs, AMD Ryzen 7 5700G, approximately 8 GB RAM; timings are not a promise for every machine.
- The canonical ledger deliberately still ends at its historical STOP. Follow `qa/resume-after-esc-23-002.md` to append the verified DEC-238 resume/recovery/resolution, then Task 1 COMPLETE. No later task was dispatched; tasks 9-11 remain held.
- Direct repair review is documented; this is not an independent-agent QA claim. Normal Heavy-mode review windows and phase gates remain required.

## Dependency-budget note

- No new dependency. Runtime uses the existing Go SQL API and SQLite driver; smoke evidence uses Python's standard library and the already-installed SQLite CLI 3.50.6 (schema requires 3.47+); no package was installed.

## Commit hash

- `730128cd8842e2943c9ab2b21f8f2b154148c10a`

## Commit subject

- `perf(phase-23): repair rebuild and task gate runtime`
