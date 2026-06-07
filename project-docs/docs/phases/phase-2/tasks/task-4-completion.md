## Summary of what was implemented

Implemented Phase 2 Task 4 by adding migration `0005_v01_ref_triggers` for the remaining v0.1 SQL backstops: polymorphic ref integrity on `links`, `claims`, and `activity_log`; link endpoint scope and relation compatibility; and claim scope. Added raw-SQL trigger tests for both reject and accept states, plus a package-level `internal/journal` determinism test that replays a journal touching every new Phase 2 table and proves delete-DB rebuild parity with all triggers active.

Kept `internal/db.ProjectionSchemaVersion` at `3`. I did not bump it because `0005` adds backstop triggers only and does not change the projected row shape or view/index generation that Tasks 1-2 already tied to rebuild freshness.

## Files created or modified

- `migrations/0005_v01_ref_triggers.up.sql`
- `migrations/0005_v01_ref_triggers.down.sql`
- `internal/db/ref_triggers_test.go`
- `internal/journal/phase2_rebuild_test.go`
- `project-docs/docs/phases/phase-2/tasks/task-4-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-2/tasks/task-4-completion.md`

## Validation Contract Coverage

- `VAL-013` PASS: `internal/db/ref_triggers_test.go` proves raw-SQL dangling-ref rejection and resolvable-row acceptance for `claims`, `links` (both endpoints), and `activity_log`. `TestValidJournalProjectedRowsStillApplyWithRefTriggersActive` also proves valid journal-applied rows are not falsely rejected by the new triggers.
- `VAL-014` PASS: `TestLinkScopeAndRelationTriggersRejectInvalidCombosAndAllowValidRows` rejects `decision` and `escalation` endpoints plus an incompatible `implements issue->issue` link, and accepts `implements issue->requirement`, `blocks issue->issue`, and `relates_to issue->requirement`.
- `VAL-015` PASS: `TestClaimScopeTriggersRejectUnsupportedRefTypeAndAllowIssueRef` rejects `ref_type='qa_review'` and accepts a resolvable `ref_type='issue'` claim.
- `VAL-016` PASS: `internal/journal/phase2_rebuild_test.go` uses journal events only, touches every new Phase 2 table, snapshots schema objects plus projected rows, removes the DB, rebuilds from journal, and asserts exact parity before and after replay. Extra CLI smoke also passed with `/tmp/fid init` and `/tmp/fid rebuild --by developer` in `/tmp/fid-task4-smoke.WhoMby`.
- `VAL-017` PASS: `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/...` showed no `cmd/fid` import under `internal/...`, `rg -n "github.com/mbell-nzl/fid/cmd/fid|cmd/fid" internal` returned no matches, and the enriched read path remains in `internal/issues.Get` while `cmd/fid/issue.go` only delegates CLI inputs/outputs.

## Tests run

- Red: `go test ./internal/db -run TestClaimRefTriggerRejectsDanglingIssueRef -count=1`
- Green: `go test ./internal/db -run TestClaimRefTriggerRejectsDanglingIssueRef -count=1`
- Red: `go test ./internal/db -run 'TestPolymorphicRefTriggersRejectDanglingRefsAndAllowResolvableRows|TestLinkScopeAndRelationTriggersRejectInvalidCombosAndAllowValidRows|TestClaimScopeTriggersRejectUnsupportedRefTypeAndAllowIssueRef|TestValidJournalProjectedRowsStillApplyWithRefTriggersActive' -count=1`
- Green: `go test ./internal/db -run 'TestPolymorphicRefTriggersRejectDanglingRefsAndAllowResolvableRows|TestLinkScopeAndRelationTriggersRejectInvalidCombosAndAllowValidRows|TestClaimScopeTriggersRejectUnsupportedRefTypeAndAllowIssueRef|TestValidJournalProjectedRowsStillApplyWithRefTriggersActive' -count=1`
- Green: `go test ./internal/journal -run TestRebuildWorkspacePreservesFullPhaseTwoProjectionWithRefTriggersActive -count=1`
- Extra phase-end evidence: `/tmp/fid init`
- Extra phase-end evidence: `/tmp/fid rebuild --by developer`

## Exact quality-gate results

- PASS `gofmt -l cmd internal`
- PASS `golangci-lint run ./...`
- PASS `go test -race -coverprofile=coverage.out ./...`
- PASS `go tool cover -func=coverage.out` (`internal/journal` 80.6%, `internal/issues` 85.1%, `cmd/fid` 70.5%, total 72.6%)
- PASS `go test -tags integration -race ./...`
- PASS `go test -bench . -benchmem ./...`
- PASS `go build -o /tmp/fid ./cmd/fid`

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `457f480`

## Commit subject

- `feat(db): add ref backstop triggers`
