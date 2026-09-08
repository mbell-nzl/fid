# Phase 22 Task 16 completion report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Recovery attempts: 0
- Human attention required: no
- Condition: None. All12 required commands terminated successfully.
- Next action: None.

## Summary of what was implemented

Applied the exact accepted 0016 up/down DDL. Latest schema16 has no escalation
table, foreign-key columns, activity reference or default-workflow seed.
Triages, Concerns, retained lifecycle readers and Event-derived activity remain.
The current portable archive semantic profile advances to16; retained v1
validation stays pinned to14 and still cannot import.

Current assertions now prove absence and retained behavior. The full historical
Phase2 replay fixture explicitly uses trusted migrations1..5 on both projection
paths, preserving its original rows. Separate latest-schema tests reject the
removed table and both removed columns and prove preceding Effects roll back.
The full migration chain resets through zero and returns to16.

## Files created or modified

- `migrations/0016_p22_remove_escalations.up.sql`, `0016_p22_remove_escalations.down.sql`.
- `internal/db/projection_meta.go`, `escalation_removal_test.go`, `migrate_v01_tables_test.go`, `ref_triggers_test.go`.
- `internal/journal/escalation_removal_test.go`, `phase2_rebuild_test.go`.
- `internal/maintenance/check.go`, `errors.go`, `test_helpers_test.go`.
- `internal/workflow/rules.go`, `registry_test.go`.
- `internal/export/archive_semantics.go`, `archive_semantics_test.go`.
- `cmd/fid/sql_test.go`, `internal/sqlquery/sqlquery_test.go`, tagged closure snapshot and portable fake schema provenance comment.
- This report and `task-16-root-evidence.md`.

## Documentation updated

The root evidence records all21 removal-root dispositions and27 changed-row
origins, not merely historical filenames. No spec change: this implements
schema.md section4 and the accepted convergence inventory. No historical
migration was edited and no later-owned lifecycle was removed.

## Validation Contract Coverage

- VAL-066: PASS. `TestFinalSchemaHasNoEscalationsOrEscalationReferences` and reset test prove latest absence, retained tables, retired activity refusal and clean full reset.
- VAL-067: PASS. Activity remains the journal-derived inspection projection; no recovery authority is read or written there. Existing activity tests and latest negative Effect tests exercise the boundaries.
- VAL-068: PASS. Every displaced row has specific introducing commit/task and executing command in root evidence; historical5/10/13/14 fixtures retain their original subjects while latest schema tests assert absence.

## Tests run

The specified first behavioral test failed with `table "escalations" exists`
(terminated `exit_code: 1`), then passed under the byte-identical command.
Root evidence retains exact focused commands and successful termination.
The focused regression includes latest migration, historical replay, maintenance,
registry, polymorphic references and two-row ordered SQL results.

## Exact quality-gate results

Root observed all12 required command occurrences terminate successfully in
this session, in the following order, from `/home/mark/projects/fid`.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	127.569s	coverage: 52.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					63.2%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	74.847s` (terminated `exit_code: 0`)
- `go version` - PASS - `go version go1.26.2 linux/amd64` (terminated `exit_code: 0`)
- `go env GOOS GOARCH` - PASS - `amd64` (terminated `exit_code: 0`)
- `uname -srm` - PASS - `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64` (terminated `exit_code: 0`)
- `nproc` - PASS - `6` (terminated `exit_code: 0`)
- `grep -m1 '^model name' /proc/cpuinfo` - PASS - `model name	: AMD Ryzen 7 5700G with Radeon Graphics` (terminated `exit_code: 0`)
- `grep -m1 '^MemTotal:' /proc/meminfo` - PASS - `MemTotal:        8133236 kB` (terminated `exit_code: 0`)
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	30.510s` (terminated `exit_code: 0`)

The fresh100k sample measured `27821 ms/rebuild` (`27821476022 ns/op`), below
the unchanged30000ms self-assertion. It does not replace Task17's three final
samples. Observed reference: Go1.26.2 linux/amd64, Linux6.6.114.1-microsoft-standard-WSL2
x86_64, six logical CPUs, AMD Ryzen7 5700G with Radeon Graphics, MemTotal8133236kB.
All identities match CF024; no other heavy process ran during this sample.

## Deviations from spec

- None. Root directly implemented the task under the Human's delegate authority, with one product writer after independent Review5 passed. The accepted exact migration and current-profile consequence require no new product decision.
- Subject-preserving convergence is detailed in root evidence. Generic project Decision rule loading remains testable; the retired seed itself is absent. SQL LIMIT2 still asserts two sorted rows rather than a shortened result.

## Concerns, technical debt, or notes for the Tech Lead

- Task17 still owns final capability qualification, recovery maintenance matrix and three reference rebuild samples. Task18 and final independent QA/Human fixtures remain pending.
- Historical migrations and frozen readers retain their existing Phase23/24/25 expiry. Production `internal/sqlquery/sqlquery.go` is unchanged, so CF018 is not triggered.

## Dependency-budget note

- None. No dependency added or updated.

## Commit hash

- `644f30e332a9e5099f0cb31c77fc6b9699fb5db8`

## Commit subject

- `feat(db): retire escalation state from schema 16`
