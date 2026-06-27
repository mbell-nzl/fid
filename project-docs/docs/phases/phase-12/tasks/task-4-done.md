## Summary of what was implemented

- Added `fid link list <id>` on the existing `link` command tree, keeping the
  read query in `internal/links` and leaving `cmd/fid` responsible only for
  flag wiring, error rendering, and format dispatch.
- Added the `internal/links` list read path so reads use the unlocked
  projection-freshness helper, return outgoing and incoming declared links with
  deterministic ordering, and resolve `entity_title` from issues and
  requirements.
- Added Phase 12 end-to-end integration coverage for requirement sync -> link
  add -> link list -> link remove -> rebuild parity, including one-event-per-
  write checks and no nested-lock failure assertions.

## Files created or modified

- `cmd/fid/e2e_v01d_link_test.go`
- `cmd/fid/link.go`
- `cmd/fid/link_test.go`
- `internal/format/markdown.go`
- `internal/links/list.go`
- `internal/links/list_test.go`
- `internal/links/types.go`
- `project-docs/docs/phases/phase-12/tasks/task-4-done.md`

## Documentation updated

- Added this completion report at
  `project-docs/docs/phases/phase-12/tasks/task-4-done.md`.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: verified no change required; the
  implemented `fid link list` output matches the existing spec shape and field
  set.

## Validation Contract Coverage

- `VAL-010`: PASS. `fid link list <id>` now returns outgoing and incoming
  declared links with `direction`, `relation`, `entity_type`, `entity_id`, and
  resolved `entity_title`, and `json`, `toon`, `markdown`, and `yaml` carry the
  same data. Evidence:
  `TestListReturnsOutgoingAndIncomingLinksWithResolvedTitlesInStableOrder`,
  `TestListReturnsEmptySliceForEntityWithNoLinks`,
  `TestLinkListReturnsOutgoingAndIncomingLinksWithResolvedTitlesAndFormatParity`,
  `TestLinkListEmptyResult`.
- `VAL-011`: PASS. The end-to-end flow runs requirement sync -> issue link add
  -> requirement link add -> `link list` (issue and requirement directions) ->
  link remove -> `fid rebuild`, with one journal event per write, no nested-lock
  failure, and identical `links` / `requirements` state after rebuild.
  Evidence: `TestV01DLinkFlowAndRebuildParity`.

## Tests run

- `go test ./cmd/fid -run TestLinkListReturnsOutgoingAndIncomingLinksWithResolvedTitlesAndFormatParity`
  — red: failed with `unknown flag: --format` because `link list` was not
  wired.
- `go test ./internal/links -run TestListReturnsOutgoingAndIncomingLinksWithResolvedTitlesInStableOrder`
  — red: build failed with `undefined: List`, `ListOptions`, and `ListItem`.
- `go test ./internal/links -run 'TestListReturnsOutgoingAndIncomingLinksWithResolvedTitlesInStableOrder|TestListReturnsEmptySliceForEntityWithNoLinks'`
  — green.
- `go test ./cmd/fid -run 'TestLinkListReturnsOutgoingAndIncomingLinksWithResolvedTitlesAndFormatParity|TestLinkListEmptyResult'`
  — green.
- `go test -tags integration ./cmd/fid -run TestV01DLinkFlowAndRebuildParity`
  — green.
- `go test ./internal/links` — green.
- `go test ./cmd/fid -run 'TestLink(Add|List|Remove)'` — green.

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS
- `golangci-lint run ./...` — PASS (`0 issues.`)
- `go test -race -coverprofile=coverage.out ./...` — PASS
- `go tool cover -func=coverage.out` — PASS
  Coverage excerpts: total `76.2%`; `cmd/fid` `77.9%`; `internal/links`
  `79.2%`; `internal/requirements` `80.4%`; `internal/issues` `81.8%`;
  `internal/claims` `81.1%`; `internal/journal` `80.0%`;
  `internal/policy` `81.2%`.
- `go test -tags integration -race ./...` — PASS
- `go build -o /tmp/fid ./cmd/fid` — PASS
- `go test -bench . -benchmem ./...` — N/A (`benchmark/scale fixtures were not touched`)

## Deviations from spec

- The report is written before the required task commit exists. Per the
  one-commit / no-amend constraint, the exact hash and subject are filled in
  after commit creation without amending the commit.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `e30e40ee96c2089607dd9f3a4847406874e32a8e` (`e30e40e`)

## Commit subject

- `feat(links): add link list command`
