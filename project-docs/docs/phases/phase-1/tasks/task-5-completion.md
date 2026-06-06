## Summary of what was implemented

- Added the `fid init` CLI command as a thin Cobra adapter that delegates bootstrap work to `internal/bootstrap` and renders the documented stderr confirmation and re-init error.
- Added bootstrap services to create the `.fid/` workspace, write canonical `config.toml`, `pricing.toml`, and `.gitignore`, run schema migrations, append exactly one `project.initialized` event through the existing lock-backed writer, and project it immediately into SQLite.
- Added config helpers to write and parse the Task 5 bootstrap files, plus tests covering file creation, default and explicit prefix handling, seeded roles, re-init rejection with no mutation, and delete-DB plus `fid rebuild` parity.

## Files created or modified

- `cmd/fid/init.go`
- `cmd/fid/init_test.go`
- `cmd/fid/root.go`
- `internal/bootstrap/init.go`
- `internal/config/config.go`
- `internal/config/config_test.go`
- `project-docs/docs/phases/phase-1/tasks/task-5-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-1/tasks/task-5-completion.md`

## Validation Contract Coverage

- `VAL-001` covered by `TestInitCommandCreatesWorkspaceFiles`, which runs `fid init` in a temp workspace and asserts `.fid/journal/`, `config.toml`, `pricing.toml`, `fid.sqlite`, and `.gitignore` exist, that `.gitignore` includes `fid.sqlite`, and that config/pricing parse successfully.
- `VAL-002` covered by `TestInitCommandStoresConfiguredOrDefaultPrefix`, which asserts `fid init --prefix ABC` stores `ABC` and a flagless init stores the uppercased directory basename.
- `VAL-003` covered by `TestInitCommandSeedsExactlyOneBootstrapEventAndSevenRoles`, which asserts exactly one journal line exists, that its event type is `project.initialized`, and that the projected roles are exactly `developer`, `qa_reviewer`, `qa_lead`, `tl_lite`, `tl`, `workflow_runner`, and `human`.
- `VAL-004` covered by `TestInitCommandRejectsReinitWithoutChanges`, which asserts a second `fid init` exits non-zero with `error: .fid/ already exists` and leaves the entire `.fid/` tree byte-identical.
- `VAL-005` covered by `TestInitCommandRebuildRoundTripRecreatesRolesAndProjectionMeta`, which deletes `.fid/fid.sqlite`, runs `fid rebuild --by developer`, and asserts the rebuilt roles, `projection_meta` checkpoint, and journal bytes match the post-init state.

## Tests run

- `go test ./cmd/fid -run TestInitCommandCreatesWorkspaceFiles`
- `go test ./cmd/fid -run 'TestInitCommand|TestRebuildCommand'`
- `go test ./internal/config -run TestWriteBootstrapFilesRoundTripsDefaultConfig`
- `gofmt -l cmd internal`
- `golangci-lint run ./...`
- `go test -race -coverprofile=coverage.out ./...`
- `go tool cover -func=coverage.out`
- `go build -o /tmp/fid ./cmd/fid`

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
- `PASS` `go tool cover -func=coverage.out`
- `PASS` `go build -o /tmp/fid ./cmd/fid`
- `Not triggered` `go test -tags integration -race ./...` because the init/rebuild tests were kept in the default test packages and not placed behind the `integration` build tag.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit creation.

## Commit subject

- Pending commit creation.
