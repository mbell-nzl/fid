# Phase 20 - Task 16: Delete universal formatters and YAML - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: gate-repair
- Condition: The QA-plan YAML absence command cannot reach its required no-match state: `go mod tidy` restores `gopkg.in/yaml.v3` checksum entries required transitively by the retained `github.com/golang-migrate/migrate/v4@v4.19.1` module.
- Recovery attempts: 1
- Next action: Obtain direction on whether VAL-032/QA-14 treats the removed direct `go.mod` requirement as sufficient, or authorize an explicitly scoped migration-dependency change; then rerun the exact YAML absence command and the full Task 16 gate.
- Human attention required: no

## Summary of what was implemented

- Deleted the complete universal formatter package and its TOON, YAML, and Markdown renderers, plus the obsolete format-matrix suite.
- Removed the last formatter product import and changed the retained migration parity test to use direct canonical JSON encoding.
- Added a product-source guard for the deleted package, imports, and universal serializer symbols; retained the Task 15 whole-Cobra-tree universal-flag guard.
- Removed the direct YAML requirement with `go mod tidy`, and updated the architecture adapter list to name `internal/cli` and `internal/project`.

## Files created or modified

- `cmd/fid/contract_surface_test.go`
- `cmd/fid/issue_helpers.go`
- `internal/db/issues_workflow_fk_cli_test.go`
- `go.mod`, `go.sum`
- Deleted `internal/format/` (all formatter source and tests) and `cmd/fid/e2e_v01_format_matrix_test.go`.
- `project-docs/docs/phases/phase-20/tasks/task-16-done.md`

## Documentation updated

- `AGENTS.md` - replaces `internal/format` with `internal/cli` and `internal/project` in the user-facing adapter list.

## Validation Contract Coverage

- VAL-030: PASS - `TestProductSourceContainsNoUniversalSerializers` at `cmd/fid/contract_surface_test.go:76` checks the `internal/format` directory and imports; `ls internal/format` produced `No such file or directory` with expected `exit_code: 2`, and `rg -n "internal/format" cmd internal` produced no output with expected `exit_code: 1`.
- VAL-031: PASS - `TestLiveCobraTreeContainsNoUniversalFormatFlags` at `cmd/fid/contract_surface_test.go:69` traverses the live Cobra tree, while its inherited/local injected-flag controls remain in the same file; the focused guard test passed and `/tmp/fid ready --help` lists only `--epic`, `--json`, and `--project`.
- VAL-032: BLOCKED - the direct `gopkg.in/yaml.v3` requirement is removed from `go.mod`, and `go mod tidy`, the race suite, and build pass, but required QA-plan command `rg -n "gopkg.in/yaml.v3" go.mod go.sum` still finds the two transitive `go.sum` lines recreated by tidy. `go mod graph` identifies `github.com/golang-migrate/migrate/v4@v4.19.1` as their source.
- VAL-033: PASS - `TestProductSourceContainsNoUniversalSerializers` at `cmd/fid/contract_surface_test.go:76` scans product ASTs, and `rg -ni "writeTOON|writeYAML|writeMarkdown|yaml.Marshal" cmd internal` produced no output with expected `exit_code: 1`.
- VAL-037: PASS (architecture-doc portion) - `AGENTS.md:32` names `cmd/fid/`, `internal/cli/`, `internal/project/`, and `internal/export/` as the user-facing adapters; the removed formatter package is no longer listed.

## Tests run

- `go test ./cmd/fid -run '^TestProductSourceContainsNoUniversalSerializers$' -count=1` - EXPECTED RED - `FAIL` (terminated exit_code: 1); the new guard reported the live formatter imports, serializers, and YAML implementation before deletion.
- `go test ./cmd/fid -run '^Test(ProductSourceContainsNoUniversalSerializers|LiveCobraTreeContainsNoUniversalFormatFlags|LiveCobraTreeUniversalFormatGuardDetectsInjectedFlags)$' -count=1` - EXPECTED RED - `FAIL` (terminated exit_code: 1); after file deletion it correctly found the remaining empty `internal/format` directory, which was then removed.
- `go test ./cmd/fid -run '^Test(ProductSourceContainsNoUniversalSerializers$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 0.072s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^Test(ProductSourceContainsNoUniversalSerializers|LiveCobraTreeContainsNoUniversalFormatFlags|LiveCobraTreeUniversalFormatGuardDetectsInjectedFlags)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 0.071s` (terminated `exit_code: 0`)
- `/tmp/fid ready --help` - PASS - `      --project string   fid project root` (terminated `exit_code: 0`)
  Real-process entrypoint evidence: the retained leaf exposes no universal format flag.
- `go mod graph | rg 'gopkg\.in/yaml\.v3'` - DIAGNOSTIC - `github.com/golang-migrate/migrate/v4@v4.19.1 gopkg.in/yaml.v3@v3.0.1` (terminated exit_code: 0).
- `go mod why -m gopkg.in/yaml.v3` - DIAGNOSTIC - `gopkg.in/yaml.v3` (terminated exit_code: 0); the path passes through the retained migration module's test graph.

## Exact quality-gate results

- `go mod tidy` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 78.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 4.658s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `rg -n "internal/format" cmd internal` - PASS - `exit 1 (no output)` (terminated `exit_code: 1`)
  Expected no-match condition.
- `rg -ni "writeTOON|writeYAML|writeMarkdown|yaml.Marshal" cmd internal` - PASS - `exit 1 (no output)` (terminated `exit_code: 1`)
  Expected no-match condition.
- `ls internal/format` - PASS - `ls: cannot access 'internal/format': No such file or directory` (terminated `exit_code: 2`)
  Expected absent-directory condition required by the QA plan.
- `rg -n "gopkg.in/yaml.v3" go.mod go.sum` - BLOCKED - `go.sum:45:gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=` (terminated `exit_code: 0`)
  Expected no-match condition is incompatible with the retained migration module after `go mod tidy`.
- QA-plan commands beyond the Task 16 spec: the two absence commands above were run; no other Task 16-specific QA-plan command exists.

## Deviations from spec

- No implementation or assertion-strength deviation. TDD preserved the initial red guard failure and the empty-directory follow-up failure; deleting the confirmed empty directory was the causal repair, after which the focused guard passed.
- Immutable QA-plan failure: `rg -n "gopkg.in/yaml.v3" go.mod go.sum` ended with the cited `go.sum` line and `exit_code: 0`. Root cause is the retained migration module's explicit YAML requirement; no in-scope repair can remove it while keeping `go mod tidy` stable. No gate, threshold, command, dependency, or assertion was changed.
- Immutable report-lint failures: the first run ended `FAIL report-lint found 13 failure(s)` with `exit_code: 13` because each PASS row had prose after the termination marker; the second ended `FAIL report-lint found 26 failure(s)` with `exit_code: 26` after report-only notation omitted the required backticks. The causal repair restored the validator's exact backticked termination notation with no trailing prose on PASS rows. Its rule still rejects the Task 16/QA-plan expected non-zero absence checks even though the task defines those exits as PASS, so no truthful report-only repair can make that validator exit zero.

## Concerns, technical debt, or notes for the Tech Lead

- VAL-032's direct-dependency intent is satisfied in `go.mod`, but QA plan section 5 requires no YAML checksum in `go.sum`, which contradicts the tidy result for the retained migration module. Resolve this before same-task completion.

## Dependency-budget note

- Removed the direct `gopkg.in/yaml.v3` dependency. No dependency was added or upgraded; the remaining transitive checksum is retained by `github.com/golang-migrate/migrate/v4@v4.19.1`.

## Commit hash

- `afb262ddc3636f6a4c9a05d6b9e77b4aeac276eb`

## Commit subject

- `refactor(cli): remove universal formatter framework`
