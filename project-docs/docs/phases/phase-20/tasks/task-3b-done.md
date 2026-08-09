# Phase 20 - Task 3b: Strict-config and init boundary evidence - Completion Report

## Summary of what was implemented

- Added the table-driven strict project-config boundary matrix required to close
  Review 1 Task 2's Blocking VAL-043 evidence gap, including schema structure,
  project ID, prefix, artifact-root, executor, maintenance, symlink, and
  byte-stable writer cases.
- Added the end-to-end `fid init` prefix/derived-prefix matrix and exact
  activated-config check.
- Added the committed dangling `.fid` symlink collision regression that closes
  Review 1 Task 3's Medium improvement: the JSON conflict remains typed, the
  symlink target is unchanged, and no staging directory remains.
- Moved the strict project-config parser/writer from `config.go` into
  `project_config.go` without an API or behavior change. Non-test files are now
  230 and 295 lines, respectively.

## Files created or modified

- `internal/config/config_test.go` - exhaustive strict-config boundary and
  writer/symlink tests.
- `cmd/fid/init_test.go` - init prefix/derived-prefix, activated-config, and
  dangling-symlink collision tests.
- `internal/config/config.go` - removed the moved strict project-config code.
- `internal/config/project_config.go` - new owning file for strict
  project-config parsing, validation, and writing.
- `project-docs/docs/phases/phase-20/tasks/task-3b-done.md` - this report.

## Documentation updated

- No contract documentation changes were required: this is DEC-130 evidence
  recovery against the already accepted config/init contract.
- Added this completion report.

## Validation Contract Coverage

- VAL-043: PASS - `TestParseProjectConfigBoundaryDomains` at
  `internal/config/config_test.go:97` proves the exact schema plus every
  requested qualifying/rejecting ID, prefix, path, executor, and maintenance
  pair; `TestProjectConfigWriterRoundTripsByteStableAtInt64Boundaries` at
  `internal/config/config_test.go:189` proves byte-stable normal/MaxInt64
  round-trips; `TestReadProjectConfigRejectsArtifactRootSymlinkComponentEscape`
  at `internal/config/config_test.go:215` proves the outside-symlink rejection;
  and `TestInitCommandPrefixBoundaryMatrix` at `cmd/fid/init_test.go:304`
  proves accepted/rejected explicit and derived init-prefix states before
  staging.
- VAL-042: PASS (supporting) -
  `TestInitCommandPreservesDanglingWorkspaceSymlinkCollision` at
  `cmd/fid/init_test.go:358` asserts the `project-already-initialized` JSON
  failure, unchanged dangling symlink, and no staging residue; the successful
  and rejected init rows at `cmd/fid/init_test.go:304` assert activation versus
  no `.fid`/no residue.
- VAL-002: PASS (supporting) - the dangling-symlink regression at
  `cmd/fid/init_test.go:358` decodes the failed `project.init` JSON response and
  asserts reason code `project-already-initialized` at `:373`, so the collision
  is exposed through the typed response surface rather than a filesystem-only
  error.
- VAL-039: PASS (supporting) - the writer round-trips at
  `internal/config/config_test.go:189`, and every successful init row calls
  `assertActivatedConfigRoundTrip` from `cmd/fid/init_test.go:339`, which reads
  the generated config through the strict parser and compares its full exact
  key/value bytes.

## Tests run

- `go test ./internal/config -run 'TestParseProjectConfigBoundaryDomains|TestProjectConfigWriterRoundTripsByteStableAtInt64Boundaries|TestReadProjectConfigRejectsArtifactRootSymlinkComponentEscape' -count=1` - RED - `FAIL` (the first 4,096-byte fixture used one host-invalid filename component; the assertion was retained using a normalized 4,096-byte path with valid components).
- `go test ./internal/config -run 'TestParseProjectConfigBoundaryDomains|TestProjectConfigWriterRoundTripsByteStableAtInt64Boundaries|TestReadProjectConfigRejectsArtifactRootSymlinkComponentEscape' -count=1` - PASS - `ok   github.com/mbell-nzl/fid/internal/config	0.010s`
- `go test ./cmd/fid -run 'TestInitCommandPrefixBoundaryMatrix|TestInitCommandPreservesDanglingWorkspaceSymlinkCollision' -count=1` - RED - `FAIL` (the new helper captured output before executing Cobra; it was corrected without changing a behavior assertion).
- `go test ./cmd/fid -run 'TestInitCommandPrefixBoundaryMatrix|TestInitCommandPreservesDanglingWorkspaceSymlinkCollision' -count=1` - PASS - `ok   github.com/mbell-nzl/fid/cmd/fid	0.141s`
- `go test ./internal/config ./cmd/fid -run 'TestParseProjectConfigBoundaryDomains|TestProjectConfigWriterRoundTripsByteStableAtInt64Boundaries|TestReadProjectConfigRejectsArtifactRootSymlinkComponentEscape|TestInitCommandPrefixBoundaryMatrix|TestInitCommandPreservesDanglingWorkspaceSymlinkCollision' -count=1` - PASS - `  295 internal/config/project_config.go`
- `/tmp/fid --project <fresh>/FID --json init` (after creating a fresh `<fresh>/FID` directory) - PASS - stdout contained one successful `project.init` Response Envelope with `issue_prefix:"FID"`, stderr was empty, and exit status was 0.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`
- `golangci-lint run ./...` - PASS - `0 issues.`
- `go test -race -coverprofile=coverage.out ./...` - PASS - `exit 0 (no output)`
- `go tool cover -func=coverage.out` - PASS - `76.2%`
- `go test -tags integration -race ./...` - PASS - `exit 0 (no output)`
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)`

## Deviations from spec

- None. The two initial focused-test failures were test-fixture/helper repairs
  during normal red/green work; the matrix assertions and accepted contract
  remained unchanged.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
