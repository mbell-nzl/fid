# Phase 17 - Task 9: Add the Full v0.1 Closure Journey - Completion Report

## Summary of what was implemented

- Added the integration-tagged, built-CLI v0.1 closure journey. It covers the
  primitive lifecycle and its rejection paths, both ready scheduling modes,
  labels, requirement and issue links/inverses, telemetry, maintenance,
  authorized vacuum, a governance version guard, Git journal layout, and
  delete-and-rebuild recovery.
- Ran the DEC-056 no-rebuild diagnostic in the completed journey. Two ordinary
  `fid maintenance check --format json` calls, four seconds apart and with no
  DB deletion or rebuild, reported `recent_activity_7d` values of 17 and 20.
  This selects Branch A: the movement is clock-relative and does not establish a
  rebuild defect.
- Replaced the former aggregate snapshot with exact, separately named snapshots
  for all 19 projected tables. Rebuild also compares only the state-derived
  maintenance subset (`status`, `errors`, `warnings`, and `info.issue_count`),
  while clock-relative counters are bounded by the no-rebuild control.

## Files created or modified

- `cmd/fid/e2e_v01_closure_test.go`
- `cmd/fid/e2e_v01_closure_helpers_test.go`
- `project-docs/docs/phases/phase-17/tasks/task-9-done.md`

## Documentation updated

- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` - verified no shipped-flow drift;
  no modification required.
- `project-docs/docs/phases/phase-17/tasks/task-9-done.md` - recorded Task 9
  implementation and validation evidence.

## Validation Contract Coverage

- VAL-018: PASS. The real compiled-CLI journey proves permission and
  no-live-claim rejections append nothing, then completes create -> claim ->
  transition -> complete -> unclaim with AC verification and concern address.
- VAL-019: PASS. The journey selects the next sequential task after cancellation
  of the current task and the highest-priority task in a prioritised epic.
- VAL-020: PASS. The journey reads back an issue label, a requirement link, an
  issue-link inverse, and proves an invalid relation is non-mutating.
- VAL-021: PASS. The journey records completion-linked priced runs, verifies the
  snapshot cost formula, and reads the expected cost and token aggregates.
- VAL-022: PASS. The DEC-056 diagnostic's no-rebuild control observed
  `recent_activity_7d` 17 then 20. After delete-and-rebuild, all 19 projected
  tables are exactly equal individually; `status`, `errors`, `warnings`, and
  `info.issue_count` are exactly equal; `recent_activity_7d` and
  `claim_count_live` move no more than their no-rebuild control permits.
- VAL-023: PASS. The temp Git workspace tracks the journal while ignoring only
  `fid.sqlite`; all non-journal commands preserve journal bytes; rebuild restores
  the seeded active roles.
- VAL-025: Prerequisites supplied, not claimed. The test uses a built CLI and a
  disposable recovery path; the Human remains the sole validator of dogfood
  judgment.

## Tests run

- `go test -tags integration -run '^TestV01ClosureJourneyAndRebuildParity$' -v ./cmd/fid` - FAIL (red): the newly added DEC-056 control referenced the not-yet-implemented maintenance-result decoder.
- `go test -tags integration -run '^TestV01ClosureJourneyAndRebuildParity$' -v ./cmd/fid` - PASS: real compiled `fid` CLI journey; the retained no-rebuild diagnostic logged `recent_activity_7d` 17 then 20 across four seconds, and the amended parity checks passed.
- `go test -tags integration -race ./...` - PASS: full integration-tagged suite, including the compiled-CLI closure journey.
- `go build -o /tmp/fid ./cmd/fid` - PASS: built CLI entrypoint.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS (`0 issues.`)
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS (total statements: 76.4%)
- `go test -tags integration -race ./...` - PASS
- `go test -bench . -benchmem ./...` - PASS
- `go build -o /tmp/fid ./cmd/fid` - PASS
- `CGO_ENABLED=0 go build -o /tmp/fid-static ./cmd/fid` - PASS
- `stat -c '%s %n' /tmp/fid-static` - PASS (`12869586 /tmp/fid-static`)
- `test "$(stat -c %s /tmp/fid-static)" -lt 31457280` - PASS

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- CF-020 remains deferred to v0.2 under DEC-056. The closure control confirms
  that burst-write timestamps can make `recent_activity_7d` advance without a
  rebuild; no product behavior was changed in this audit task.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `7119d5b38d79663f544f093d55f676c375746419`

## Commit subject

- `test(closure): add v0.1 closure journey`
