# Phase 23 - Task 16: Portable recovery for the review and judgment domain - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- The archive semantic validator accepts the existing boolean representation of `review_findings.blocking`. This lets `export validate` and `export import` replay real Phase 23 Review Events against projection schema 18.
- A built-binary integration journey exports, validates, rebuilds, and imports four populated Phase 23 workspaces. It compares every application table row and every canonical archive entry between live, rebuilt, and imported state, and verifies the schema-18 projection integrity. The four workspaces collectively populate every table named below and emit all nine Phase 23 Event types.
- The journey verifies non-mutating validation, exact schema-18 import and schema-17 rejection, one `review-cycles` sample keyed to its Review ID, and rejection with no projection change for three corrupt payloads per Phase 23 Event type plus an unknown type.
- Review-13 MF-1 repair: each of those 28 corrupt Event probes now clones the same clean `.fid` workspace into its own temporary directory, opens a separate projection connection, verifies its complete starting projection against the source, and checks the complete projection again after current-Event rejection. A named test runs all 28 probes independently and logs each workspace, validator error, and before/after projection digest.
- Final-review-1 MF-3 repair: the shared archive fixture now acquires its source Claim for one hour, long enough for the retained import/recovery journey under its unchanged 30-minute test timeout. `TestImportRequiresFreshDestinationObservationBeforeRecovery` reads the persisted acquisition and expiry timestamps and asserts that lease length before exercising its historical rejection, destination observation, release, recovery, successor Gate, and import-retry checks. Claim expiry rules and Event content are unchanged.
- MF-3 provenance: `qa/final-review-1.md` traces the retained import recovery test to phase-22/task-14 commit `cc73b092bf5ecd3fdba62a1362062a6d543697ab` and the 60-second archive fixture lease to phase-22/task-13 commit `1be8e882ef689f15705507c448e3b087e23d5f65`. The accepted task-16 import scope owns this bounded repair, as recorded in the final-review finding and the latest `TASK-RECOVERY task-16` ledger entry.
- The original implementation commit was `d02614e6c26757076a91d8a12e08dfa328fe699d` (`fix(export): validate portable Phase 23 review recovery`); this repair preserves that implementation and its prior report evidence.

## Files created or modified

- `internal/export/archive_semantics_values.go`: accept `review_findings.blocking` as a JSON boolean in the archive semantic replay.
- `cmd/fid/integration/p23_portable_recovery_test.go`: portable parity, corruption, schema, and charge journey; review-13 MF-1 isolation repair and named corruption test.
- `internal/export/archive_fixture_test.go`: one-hour source Claim lease for archive fixture consumers.
- `internal/export/import_test.go`: persisted lease-length regression assertion in the retained import/recovery journey.
- `project-docs/docs/phases/phase-23/tasks/task-16-done.md`: this report.

## Documentation updated

- This report records the Phase 23 table and Event inventory, validation evidence, exact command results, and provenance. It also preserves the final-review-1 MF-3 full-checkptr failure and the focused red/green and gate reruns. No governing contract or Runner-owned ledger was edited.

## Validation Contract Coverage

- VAL-004: PASS - `TestMigration0018RoundTripIntegrity` checks `foreign_key_check` and `integrity_check` on the final schema during the ordinary gate; `TestPortableRoundTripCoversEveryPhase23Table` calls `p23AssertRebuiltProjectionIntegrity` after each real CLI rebuild and import (`p23_portable_recovery_test.go:233` and `:246`). The helper checks no foreign-key rows and `integrity_check = ok`.
- VAL-027 (import leg): PASS - `TestPortableRoundTripCoversEveryPhase23Table/review-exhaustion` asserts one `review-cycles` row, consumed value 1, and a `source_sample_id` matching `review_records.id` (`p23_portable_recovery_test.go:88`). Its exact `budget_samples` row is compared after rebuild and import (`:219`); the imported state has no duplicate charge. Task 6's recording, replay, and linked-recovery legs remain covered by `TestP23ReviewCycleChargeSurvivesLinkedRecoveryAndRebuild` in the full tagged gate.
- VAL-064: PASS - `TestPortableRoundTripCoversEveryPhase23Table` compares every application table's complete sorted rows (`p23_portable_recovery_test.go:254`, `:288`) and canonical archive entries including ordered Event bytes (`:219`) across live, rebuilt, and imported workspaces. It checks `.fid` file bytes before and after `export validate` (`:300`), and asserts every named table has populated rows (`:22`, `:34`). Archive entry equality also preserves project, entity, and Event identities and history. The direct CLI create, validate, and import below additionally returned the same project ID and terminal Event hash.
- VAL-065: PASS - the current archive imports and reports `projection_schema:18` (`p23_portable_recovery_test.go:219`); `p23PortableAssertPriorSchemaRejected` changes only the manifest projection schema to 17, then checks `archive-invalid` with `schema-unsupported` and absence of a published `.fid` (`:173`). The current manifest remains `fid.export-manifest.v2`.
- VAL-076: PASS after review-13 MF-1 repair - the original `p23PortableAssertCorruptEventsRejected` (`p23_portable_recovery_test.go:104` in the original task commit) used real exported Events for all nine named types, with unknown Effect table, missing required Effect, and malformed scalar body cases, plus one unknown type. It asserted current-Event rejection and complete projection equality, but reused one workspace and connection, so its isolation evidence was insufficient. The repaired helper (`p23_portable_recovery_test.go:129`) clones a clean workspace and opens a new connection in each subtest (`:147-158`), verifies each clone equals the source before injection (`:159-162`), retains the validator assertion (`:187-190`), and compares every table after rejection (`:191-198`). `TestPortableAssertCorruptEventsRejected` (`:88`) supplied 28 passing subtests with distinct workspace paths and per-probe validator output and digest pairs; the 28-case record below names each case. The complete 59-table snapshot SHA-256 was `6a279ff389b35e1d3689116706d3d24677bf9a2fe76942f69a8505277cb853f9` both before and after **each** probe in the verbose run.
- VAL-079: PASS - original ordinary `cmd/fid` took `246.504s` and `internal/workflow` took `100.101s` (each under 600s); tagged `cmd/fid/integration` took `269.692s` (under 900s). The review-13 repair gate reported uncached `cmd/fid` `196.206s`, uncached `internal/workflow` `57.597s` (the ordinary gate's workflow result was cached), and tagged `cmd/fid/integration` `232.525s`. All original and repair measurements are below 80% of their respective budgets; neither gate timed out.

Final-review-1 MF-3 has **no VAL** of its own: it repairs the mandatory QA plan §22 full-checkptr control. The current unchanged task gates still execute the VAL-004/027/064/065/076 evidence above. This attempt's ordinary gate reported `cmd/fid` `215.384s` and `internal/workflow` `(cached)`; the earlier uncached workflow measurement remains `57.597s`. Tagged integration reported `281.806s`. These are below the 600s ordinary and 900s tagged budgets; none reaches 80% headroom warning. The separate full-checkptr control reported `cmd/fid` `356.407s` and `internal/export` `163.526s` and is not subject to VAL-079's budget.

The table provenance is the accepted Phase 23 task that introduced or changed the projected row, as recorded in `contract-coverage.md`. Each row below is executed by `TestPortableRoundTripCoversEveryPhase23Table` under `go test -tags integration -run '^TestPortableRoundTripCoversEveryPhase23Table$' -count=1 -v ./cmd/fid/integration/` (PASS, `exit_code: 0`) and the mandatory tagged race gate. The named subcase supplies non-empty rows; `p23PortableRoundTrip` compares the entire row set across all three projections, and compares the complete archive.

| Table | Row provenance | Executing subcase |
|---|---|---|
| `review_records` | Phase 23 task 4, Review record | `mixed-findings` |
| `review_findings` | Phase 23 task 4, Review Findings | `mixed-findings` |
| `finding_dispositions` | Phase 23 task 5, Finding disposition | `mixed-findings` |
| `concerns` | Phase 23 task 8 generic Concern; task 15 rename from `judgment_concerns` | `references-and-judgment` |
| `deviations` | Phase 23 task 8, Deviation | `references-and-judgment` |
| `decision_refs` | Phase 23 task 9, Decision Reference | `protected-decision` |
| `human_outcomes` | Phase 23 task 10, Human Outcome | `protected-decision` |
| `envelope_protected_decisions` | Phase 23 tasks 2 and 9, protected Decision binding | `protected-decision` |
| `requirement_refs` | Phase 23 task 12, Requirement Reference | `references-and-judgment` |
| `requirement_issue_links` | Phase 23 task 12, traceability link | `references-and-judgment` |
| `envelope_obligations` | Phase 23 task 2, Review and judgment Obligation generalization | `mixed-findings` |
| `evidence_records` | Phase 23 tasks 1, 6, and 10, three Evidence subject types | `protected-decision` |
| `budget_samples` | Phase 23 task 6, Review-cycle charge | `review-exhaustion` |
| `activity_log` | Phase 23 tasks 8, 13, and 15, Review and judgment Activity Effects | `references-and-judgment` |

The same executing command checks these nine exported Event types in `p23PortableJournalEvents` (`p23_portable_recovery_test.go:324`) and the three corruption subcases in `p23PortableAssertCorruptEventsRejected` (`:104`): `review.recorded` (task 4), `finding.addressed` (task 5), `concern.recorded` and `concern.disposed` (task 8, rename converged by task 15), `deviation.recorded` (task 8), `decision.referenced` (task 9), `human-outcome.recorded` (task 10), `requirement.referenced` and `requirement.linked` (task 12). An `unknown.phase23` type is rejected by that command too. No retained test row was changed or converged; this task adds new rows and preserves the existing assertions.

Review-13 MF-1 repair evidence: the named verbose test below printed one distinct `workspace=/tmp/TestPortableAssertCorruptEventsRejected...` path, the quoted `current Event` validator error, `tables=59`, and the full before/after SHA-256 pair stated above for **each** populated cell. `stage 1` and `stage 2` abbreviate the verbatim `missing or reordered Effect stage 1/2` errors; `invalid id` abbreviates the verbatim `invalid column or scalar "id"` error with the table named in the cell. Each cell is a separately cloned workspace and connection, checked against the same clean baseline before injection; the complete row snapshot was identical afterwards.

| Event type | Unknown effect table | Missing required effect | Malformed body |
|---|---|---|---|
| `review.recorded` | `stage 1` | `stage 1` | `review_records/insert: invalid id` |
| `finding.addressed` | `stage 1` | `stage 1` | `finding_dispositions/insert: invalid id` |
| `concern.recorded` | `stage 1` | `stage 1` | `concerns/insert: invalid id` |
| `concern.disposed` | `stage 1` | `stage 1` | `concerns/update: invalid id` |
| `deviation.recorded` | `stage 1` | `stage 1` | `deviations/insert: invalid id` |
| `decision.referenced` | `stage 1` | `stage 1` | `decision_refs/insert: invalid id` |
| `human-outcome.recorded` | `stage 1` | `stage 1` | `human_outcomes/insert: invalid id` |
| `requirement.referenced` | `stage 2` | `stage 2` | `requirement_refs/insert: invalid id` |
| `requirement.linked` | `stage 1` | `stage 1` | `requirement_issue_links/insert: invalid id` |
| `unknown.phase23` | `unsupported Event family "unknown.phase23"` | N/A | N/A |

## Tests run

- Focused red: `go test -tags integration -run '^TestPortableRoundTripCoversEveryPhase23Table$' -count=1 -v ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`). The diagnostic was `review_findings.blocking requires an integer, not a boolean` in `export validate`; this is the first expected TDD failure.
- Focused test-fixture attempt: the same command - FAIL - `FAIL` (terminated `exit_code: 1`). The new `references-and-judgment` fixture had not created its Git directory. The fixture now creates it.
- Focused corruption-test attempt: the same command - FAIL - `FAIL` (terminated `exit_code: 1`). A zero-Effect test Event could not be hashed; the test now removes only its required Effect, leaving a serializable malformed Event.
- Focused final: `go test -tags integration -run '^TestPortableRoundTripCoversEveryPhase23Table$' -count=1 -v ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	9.125s` (terminated `exit_code: 0`)
- Pre-gate `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/export/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	74.008s` (terminated `exit_code: 0`)
- Pre-gate `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestPortableRoundTripCoversEveryPhase23Table$' ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	13.085s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid-task16 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- Direct CLI `/tmp/fid-task16 --project /tmp/fid-task16-cli-KjxBpO/source --json export create /tmp/fid-task16-cli-KjxBpO/archive.zip` - PASS - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"export.create","result":{"export":{"archive_path":"/tmp/fid-task16-cli-KjxBpO/archive.zip","archive_digest":"sha256:bfc37ca738c91111eaae92d7b200ac2622e50584a2e3d9f238a14721304bb0b0","manifest_schema":"fid.export-manifest.v2","project_id":"project-0b7bf012d95afc763c9b5a706b4c8a47","event_count":1,"terminal_event_hash":"sha256:826a110a565da6901e0f173e7d85164f044c60d2c6f0ba0a63d944116627a8bd","projection_schema":18,"entry_count":2}}}` (terminated `exit_code: 0`)
- Direct CLI `/tmp/fid-task16 --project /tmp/fid-task16-cli-KjxBpO/source --json export validate /tmp/fid-task16-cli-KjxBpO/archive.zip` - PASS - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"export.validate","result":{"validation":{"archive_path":"/tmp/fid-task16-cli-KjxBpO/archive.zip","archive_digest":"sha256:bfc37ca738c91111eaae92d7b200ac2622e50584a2e3d9f238a14721304bb0b0","status":"valid","manifest_schema":"fid.export-manifest.v2","entry_count":2}}}` (terminated `exit_code: 0`)
- Direct CLI `/tmp/fid-task16 --json export import /tmp/fid-task16-cli-KjxBpO/archive.zip --repository /tmp/fid-task16-cli-KjxBpO/target --idempotency-key task16-direct-cli` - PASS - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"export.import","result":{"import":{"project_id":"project-0b7bf012d95afc763c9b5a706b4c8a47","archive_digest":"sha256:bfc37ca738c91111eaae92d7b200ac2622e50584a2e3d9f238a14721304bb0b0","event_count":1,"terminal_event_hash":"sha256:826a110a565da6901e0f173e7d85164f044c60d2c6f0ba0a63d944116627a8bd","projection_schema":18,"target_root":"/tmp/fid-task16-cli-KjxBpO/target","replayed":false}}}` (terminated `exit_code: 0`)

Review-13 MF-1 repair tests, working directory `/home/mark/projects/fid`:

- Initial named command before adding its test: `go test -tags integration -run '^TestPortableAssertCorruptEventsRejected$' -count=1 ./cmd/fid/integration/` - PASS as a command, no test executed - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.970s [no tests to run]` (terminated `exit_code: 0`). This was not behavior evidence.
- Red after adding the named test and the connection-isolation assertion: the same command - FAIL - `FAIL` (terminated `exit_code: 1`); 27 subtests reported `corruption probe reused a projection connection` after the first probe.
- First green after cloning per case: the same command - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.134s` (terminated `exit_code: 0`).
- Per-probe evidence: `go test -tags integration -run '^TestPortableAssertCorruptEventsRejected$' -count=1 -v ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.197s` (terminated `exit_code: 0`). Its 28 named subtests each logged a distinct workspace path, current-Event rejection, and equal full-projection digest pair over 59 tables; the case matrix above records the validator output for all 28.
- Final byte-identical named command on the logged implementation: `go test -tags integration -run '^TestPortableAssertCorruptEventsRejected$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.156s` (terminated `exit_code: 0`).
- Tagged pre-gate race check: `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestPortableAssertCorruptEventsRejected$' ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	11.158s` (terminated `exit_code: 0`).
- Uncached VAL-079 workflow measurement after the six exact gates: `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/workflow/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	57.597s` (terminated `exit_code: 0`). The ordinary gate reported workflow as cached, so this supplies its current uncached package time.
- Report-lint expectation after the finalized report commit: `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-16-done.md project-docs/docs/phases/phase-23/tasks/task-16.md` must terminate with `OUTCOME COMPLETE`. The post-commit observed result is returned to the Phase Runner.
- Report-lint attempt 1 on finalized report commit `e25fa2d`: the same command - FAIL - `FAIL report-lint found 2 failure(s)` (terminated `exit_code: 2`). Two new gate result lines had explanatory prose after their termination fields. The entries below now end at the observed termination, with those notes on separate lines; the original failed control evidence is retained here.

Final-review-1 MF-3 and this bounded repair (working directory `/home/mark/projects/fid`):

- QA attempt 1, mandatory `go test -race -timeout 30m ./...` - FAIL - `FAIL` (terminated `exit_code: 1`). The failing package reported `--- FAIL: TestImportRequiresFreshDestinationObservationBeforeRecovery (109.65s)` and `import_test.go:263: release imported Claim: claim expired`; see `qa/final-review-1.md` MF-3 and its linked `final-evidence-1/full-checkptr.json`. This is QA's observed failed attempt, retained here.
- Focused red after adding the persisted lease assertion: `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestImportRequiresFreshDestinationObservationBeforeRecovery$' ./internal/export/` - FAIL - `FAIL` (terminated `exit_code: 1`). The test reported `source Claim lease = 1m0s, want at least 1h for import recovery under the 30m test timeout`; package result `FAIL github.com/mbell-nzl/fid/internal/export 4.417s`.
- Focused green after extending the fixture lease: the same command - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	11.465s` (terminated `exit_code: 0`). The complete test retains its negative historical-observation and Gate checks and its positive destination release, recovery, successor Gate, and idempotent retry checks.
- Pre-gate touched package: `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/export/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	77.432s` (terminated `exit_code: 0`)
- Later byte-identical full-checkptr control: `go test -race -timeout 30m ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)` (terminated `exit_code: 0`)

The pre-gate touched package wall time was `83.364s`; the full-checkptr control wall time was `374.439s`. Full-checkptr also reported `internal/export 163.526s` and `cmd/fid 356.407s`.

## Exact quality-gate results

Working directory for all six commands: `/home/mark/projects/fid`. The six original task-spec commands below ran in order on the original task implementation tree. Their evidence is retained unchanged.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	100.101s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:									(statements)			60.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	269.692s` (terminated `exit_code: 0`)
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	25.376s` (terminated `exit_code: 0`)

Benchmark row: `22689 ms/rebuild` (reported by the same command).

Review-13 MF-1 repair gate: the same six commands ran byte-identically, in the task-spec order, on the repaired fixture tree. No QA-plan command beyond the task spec was required of this Developer; QA-41 and QA-44 are phase-end QA Lead checks.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:									(statements)			60.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	232.525s` (terminated `exit_code: 0`)
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	25.218s` (terminated `exit_code: 0`)

The repair gate reported `cmd/fid` `196.206s` uncached; the benchmark row was `22542 ms/rebuild`.

Final-review-1 MF-3 repair gate: all six unchanged task-spec commands ran in order from `/home/mark/projects/fid` after the focused proof. Wall time is measured from invocation through observed termination.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.7% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					60.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	281.806s` (terminated `exit_code: 0`)
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	27.845s` (terminated `exit_code: 0`)

Observed wall times in the same command order: `0.219s`, `4.425s`, `232.800s`, `0.277s`, `299.608s`, and `29.133s`.

The ordinary gate's uncached `cmd/fid` package took `215.384s`, and the benchmark reported `24311 ms/rebuild`. The full-checkptr control result is in `## Tests run` above.

## Deviations from spec

- None. The initial focused failure was the missing archive boolean allowance; adding `review_findings.blocking` to the accepted semantic boolean columns produced the passing focused and full gates. The two intermediate focused failures were new test-fixture defects and were repaired in the same test file. All failed attempts above remain recorded; the exact same focused command later passed. Review-13 MF-1 identified that the original corruption assertion reused one projection: the new named test failed with the copied terminal token `FAIL` (terminated `exit_code: 1`) because 27 cases reused the connection. Cloning the clean workspace and opening a new connection inside each case repaired that cause; the byte-identical command later returned `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.156s` (terminated `exit_code: 0`). The first report-lint attempt failed with `FAIL report-lint found 2 failure(s)` (terminated `exit_code: 2`) because two PASS lines had trailing prose; moving those notes after the six-command list lets the unchanged control parse each observed termination. No accepted contract, command, threshold, or product behavior changed.

- Final-review-1 MF-3 caused a required-command deviation: QA's first byte-identical `go test -race -timeout 30m ./...` ended `FAIL` (terminated `exit_code: 1`) at `release imported Claim: claim expired`. The retained source fixture acquired a Claim for only 60 seconds, while the import/recovery journey performs export, import, and several rejection checks before release. The new persisted-timestamp assertion deterministically failed with a `1m0s` lease and `FAIL` (terminated `exit_code: 1`). Extending only that fixture lease to one hour made the same focused command pass with `ok  	github.com/mbell-nzl/fid/internal/export	11.465s` (terminated `exit_code: 0`) and the byte-identical full-checkptr command pass with `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)` (terminated `exit_code: 0`). Claim expiry enforcement, negative and positive recovery assertions, Event identity, gate definitions, and budgets remain intact.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
