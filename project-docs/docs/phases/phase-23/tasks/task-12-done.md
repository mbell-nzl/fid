# Phase 23 - Task 12: Requirement references and a single traceability model - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 2
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Replaced V1 markdown requirement authoring with digest-verified Requirement References and requirement-to-Issue links. A trusted finalized-Activation Actor can record and link without either removed capability name or an Envelope grant (DEC-260). Supersession retains the old row, points it to the successor, and keeps one active row.
- Removed the V1 `internal/requirements` parser and writer, `requirement sync` leaf, and `requirement.synced` rules from both Event catalogs. The `links` adapter now admits Issue endpoints only and excludes `implements` and `implemented_by`; traceability uses `requirement_issue_links`.
- Converged retained command, journal, and real-process fixtures to the accepted rejection and no-state-change behavior. The requirement journey observes digest mismatch, supersession, all three stored relations, and both unknown-link rejections.

## Files created or modified

- `internal/review/requirement.go`; `internal/journal/project_effects.go`, `project_effects_test.go`, `current_event_catalog.go`; `internal/export/archive_semantics_catalog.go`.
- `internal/links/add.go`, `list.go`, `registry.go`, `errors.go` and their matching `*_test.go` files.
- `cmd/fid/link_response.go`, `link_test.go`, `link_limit_test.go`, `e2e_v01d_link_test.go`, `requirement_test.go`, `contract_surface_test.go`, `task9_contract_test.go`; `cmd/fid/integration/p23_requirement_reference_test.go`, `e2e_v01_closure_test.go`, `e2e_v01_closure_helpers_test.go`.
- Deleted `internal/requirements/` with its V1 parser, sync writer, and package tests. The new reference journey and converged retained tests cover the accepted replacement behavior.
- `project-docs/docs/phases/phase-23/tasks/task-12-done.md`.

## Documentation updated

- `project-docs/specs/03_TRD_SCHEMA.md` section 3.2 now states Issue-only `links` endpoints and the three `requirement_issue_links` relations.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 3 names `requirement record-reference|link|get|list`.

## Validation Contract Coverage

- VAL-003: PASS - `TestProjectEventFenceRejectsMalformedPhase23EventsWithoutProjectionChange` at `internal/journal/project_effects_test.go:14` rejects each of the nine Phase 23 families plus an unknown family and compares projection state before/after; `TestConcernDeviationFenceConvergencePreservesEveryRemainingEvent` at `:53` asserts all nine strict current-Event rules remain registered. The focused journal command passed.
- VAL-051: PASS - `TestP23RequirementReferenceAndLinkRealProcess` at `cmd/fid/integration/p23_requirement_reference_test.go:16` records twice through the real binary; `:40` checks `source-digest-mismatch` with unchanged Event count; `:45` reads the old row and successor pointer; `:55` reads exactly one active row with the new ID. The focused journey and tagged gate passed.
- VAL-052: PASS - The same real-process test at `cmd/fid/integration/p23_requirement_reference_test.go:67` records `implements`, `verifies`, and `relates-to`; `:83` reads back all three stored values; `:106` rejects an unknown requirement and relation with no new Event.
- VAL-053: PASS - `TestRequirementCommandHasReferenceLeaves` at `cmd/fid/requirement_surface_test.go:5` asserts the `sync` leaf absent; `TestV01ClosureJourneyAndRebuildParity` at `cmd/fid/integration/e2e_v01_closure_test.go:14` runs the removed leaf through the built binary and checks unchanged journal bytes; `internal/requirements/` is deleted; `internal/journal/current_event_catalog.go:34` and `internal/export/archive_semantics_catalog.go:46` contain only the replacement rules. The exact declared source search returned `exit 1 (no output)` (expected and terminated `exit_code: 1`). The V1 table column remains until task 15 as the contract states.
- VAL-054: PASS - `TestLinkAddRejectsRequirementEndpointsAndInvalidInputBeforeAppend` at `cmd/fid/link_test.go:119` checks endpoint and removed-relation rejection without rows or Events; `TestV01ClosureJourneyAndRebuildParity` at `cmd/fid/integration/e2e_v01_closure_test.go:14` checks a real-process requirement-endpoint rejection and a retained Issue link with rebuild parity; the requirement journey at `cmd/fid/integration/p23_requirement_reference_test.go:111` checks both removed relation names through the binary with unchanged Event counts.
- VAL-079 (task-12 report part): PASS - Exact ordinary gate reported `cmd/fid` 189.843s (first passing attempt 285.237s) and `internal/workflow` 100.258s on the first passing attempt, then `(cached)` on the final tree. Exact tagged integration reported 208.862s. All observed times are below the ordinary 600s and tagged 900s budgets; neither gate timed out.

## Tests run

Working directory for every command below: `/home/mark/projects/fid`. Each termination is copied from that command's result.

- `go test -run '^TestRequirementRecordReferenceSupersedesPriorActiveRow$' ./internal/review` - FAIL - `FAIL` (terminated `exit_code: 1`); earlier compiler detail: `undefined: requirementReferenceEffects`.
- `go test -run '^TestRequirementRecordReferenceSupersedesPriorActiveRow$' ./internal/review` - PASS - `ok  	github.com/mbell-nzl/fid/internal/review	0.005s` (terminated `exit_code: 0`)
- `go test -run '^TestRequirementCommandHasReferenceLeaves$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`); earlier four leaves were missing and `sync` remained.
- `go test -run '^TestRequirementCommandHasReferenceLeaves$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.131s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP23RequirementReferenceAndLinkRealProcess$' ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); earlier Activation profile returned `invalid activation input`.
- `go test -tags integration -run '^TestP23RequirementReferenceAndLinkRealProcess$' ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); earlier expanded profile still returned `invalid activation input` because both requirement names were absent from the closed registry.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestRequirementRecordReferenceSupersedesPriorActiveRow|TestRequirementCommandHasReferenceLeaves)$' ./internal/review ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.041s` (terminated `exit_code: 0`)
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-12-done.md project-docs/docs/phases/phase-23/tasks/task-12.md` - FAIL - `FAIL report-lint found 4 failure(s)` (terminated `exit_code: 4`); earlier report punctuation plus the task-spec grammar mismatch.
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-12-done.md project-docs/docs/phases/phase-23/tasks/task-12.md` - FAIL - `FAIL report-lint found 1 failure(s)` (terminated `exit_code: 1`); the remaining task-spec grammar mismatch was resolved by DEC-260 in the spec.
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-12-done.md project-docs/docs/phases/phase-23/tasks/task-12.md` - FAIL - `FAIL report-lint found 21 failure(s)` (terminated `exit_code: 21`); this continuation's first finalized report put prose or punctuation after PASS termination receipts.
- `go test -tags integration -run '^TestP23RequirementReferenceAndLinkRealProcess$' ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); this session's red test reached `authority-denied` without either requirement capability.
- `go test -tags integration -run '^TestP23RequirementReferenceAndLinkRealProcess$' ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); after removing the capability checks, supersession reached `missing or reordered Effect stage 1`.
- `go test -tags integration -run '^TestP23RequirementReferenceAndLinkRealProcess$' ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); after correcting stage order, the supersession pointer hit `FOREIGN KEY constraint failed (787)`.
- `go test -tags integration -run '^TestP23RequirementReferenceAndLinkRealProcess$' ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.668s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP23RequirementReferenceAndLinkRealProcess$' ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.539s` (terminated `exit_code: 0`)
- `go test -run '^TestRegistryHelpersCoverEndpointAndCompatibilityRules$' ./internal/links` - FAIL - `FAIL` (terminated `exit_code: 1`); old endpoint registry admitted Requirement.
- `go test -run '^TestRegistryHelpersCoverEndpointAndCompatibilityRules$' ./internal/links` - PASS - `ok  	github.com/mbell-nzl/fid/internal/links	0.005s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestV01ClosureJourneyAndRebuildParity$' ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`); the first converged fixture invoked `sync` without its removed flags and the parent command returned success/help.
- `go test -tags integration -run '^TestV01ClosureJourneyAndRebuildParity$' ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.427s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/review ./internal/journal ./internal/links ./internal/export ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	221.383s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP23RequirementReferenceAndLinkRealProcess$' ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.300s` (terminated `exit_code: 0`)

## Exact quality-gate results

Earlier checkpoint: the following seven BLOCKED rows were retained as failed-attempt history. None had a termination result in that checkpoint.

- `gofmt -l cmd internal` - BLOCKED - not executed; task incomplete and scope gap pending Runner triage (no termination observed).
- `golangci-lint run ./...` - BLOCKED - not executed; task incomplete and scope gap pending Runner triage (no termination observed).
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - not executed; task incomplete and scope gap pending Runner triage (no termination observed).
- `go tool cover -func=coverage.out` - BLOCKED - not executed; no task coverage profile existed (no termination observed).
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - not executed; focused journey rejected at Activation (no termination observed).
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - BLOCKED - not executed; task incomplete and scope gap pending Runner triage (no termination observed).
- `rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal` - BLOCKED - not executed; V1 writer still existed (no termination observed).

First exact attempt in this continuation:

- `gofmt -l cmd internal` - FAIL - `internal/journal/project_effects_test.go` (terminated `exit_code: 0`); gofmt reports unformatted files with a successful process exit.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	100.258s	coverage: 64.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					61.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`); `TestV01ClosureJourneyAndRebuildParity` still called removed `requirement sync`.
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	24.893s` (terminated `exit_code: 0`)
- `rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)

Final exact attempt on the repaired tree, in task-spec order:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 64.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					61.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	208.862s` (terminated `exit_code: 0`)
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	24.913s` (terminated `exit_code: 0`)
- `rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)

## Deviations from spec

- DEC-260 resolved ESC-23-006: requirement writes use a trusted Actor from finalized Kit Activation without an Actor-profile capability or Envelope grant. The closed thirty-name registry and both allowlist copies remained unchanged. DEC-260 also corrected the task spec's one-line expected-nonzero declaration; this continuation changed neither the task spec nor `report-lint.sh`.
- Earlier report-lint attempt 1 ended `FAIL report-lint found 4 failure(s)` (terminated `exit_code: 4`); report punctuation was corrected. Attempt 2 ended `FAIL report-lint found 1 failure(s)` (terminated `exit_code: 1`), caused by the task-spec grammar that DEC-260 corrected. Both failures remain above.
- The first exact tagged gate ended `FAIL` (terminated `exit_code: 1`): retained V1 closure fixture called removed `requirement sync`. The fixture was converged within task 12 scope to require rejection and no journal change; its focused rerun and the byte-identical tagged gate then passed. The first format check printed `internal/journal/project_effects_test.go` (terminated `exit_code: 0`); `gofmt -w` corrected it, and the byte-identical check returned no output.
- The first finalized report-lint attempt ended `FAIL report-lint found 21 failure(s)` (terminated `exit_code: 21`): PASS rows carried prose or punctuation after their observed termination receipts. This report correction keeps receipts as the final text on each PASS row; the byte-identical control is rerun after committing this correction.

## Concerns, technical debt, or notes for the Tech Lead

- Provenance checked against introducing commits and the accepted task pack: the task-12 fixture was introduced by `0e4bb33`; its only capability-related assertion was the Actor-profile action list, now removed. The displaced V1 requirement tests originated in `e236591`; retained `links` tests trace to `0abe93e` and `e30e40e`; the Phase 23 fence parity test originated in `3c0b520`. Convergence follows inventory rows 21-22 and task 12's VAL-003 instruction. No historical task number was used as ownership authority.
- VAL-079 headroom: the final tagged 208.862s is below 80% of 900s; ordinary `cmd/fid` 189.843s and earlier measured `internal/workflow` 100.258s are below 80% of 600s. No headroom warning is triggered.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `1679230ee4dca4d1cbffbb5b2732b7d1a9dbd06d`

## Commit subject

- `feat(requirement): complete reference and traceability convergence`
