# Phase 20 - Task 16: Delete universal formatters and YAML - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 3
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Deleted the complete universal formatter package and its TOON, YAML, and Markdown renderers, plus the obsolete format-matrix suite.
- Removed the last formatter product import and changed the retained migration parity test to use direct canonical JSON encoding.
- Added a product-source guard for the deleted package, imports, and universal serializer symbols; retained the Task 15 whole-Cobra-tree universal-flag guard.
- Removed the direct YAML requirement with `go mod tidy`, and updated the architecture adapter list to name `internal/cli` and `internal/project`.
- The third same-task continuation reran every current Task 16 command against the committed tree. All twelve commands passed under DEC-145, including the one-stanza `go mod why -vendor` build-graph proof.

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

- VAL-030: PASS - this continuation's `ls internal/format` returned the absent-directory diagnostic with declared/observed `exit_code: 2`, and `rg -n "internal/format" cmd internal` returned no output with declared/observed `exit_code: 1`; `TestProductSourceContainsNoUniversalSerializers` remains at `cmd/fid/contract_surface_test.go:76`.
- VAL-031: PASS - `TestLiveCobraTreeContainsNoUniversalFormatFlags` remains at `cmd/fid/contract_surface_test.go:69`; this continuation's built entrypoint `/tmp/fid ready --format json` rejected with observed `exit_code: 1`.
- VAL-032: PASS - this continuation's `go mod tidy`, full race/coverage suite, and build all terminated `exit_code: 0`; `rg -ni "yaml" go.mod` returned the declared no-match result, and `go mod why -vendor -m gopkg.in/yaml.v3` emitted exactly one DEC-145 stanza: `# gopkg.in/yaml.v3` followed by `(main module does not need to vendor module gopkg.in/yaml.v3)`, with observed `exit_code: 0`. The withdrawn `go.sum` text search remains immutable history and was not rerun.
- VAL-033: PASS - this continuation's `rg -ni "writeTOON|writeYAML|writeMarkdown|yaml.Marshal" cmd internal` returned no output with declared/observed `exit_code: 1`; `TestProductSourceContainsNoUniversalSerializers` at `cmd/fid/contract_surface_test.go:76` remains the source guard.
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
- `/tmp/fid ready --format json` - PASS (expected rejection) - `exit 1 (no output)` (terminated `exit_code: 1`)
  Real-process entrypoint evidence from the third same-task continuation: the removed universal format flag is rejected.

## Exact quality-gate results

- `go mod why -vendor -m gopkg.in/yaml.v3` - BLOCKED - `(main module does not need to vendor module gopkg.in/yaml.v3)` (terminated `exit_code: 0`)
  Protected literal-output mismatch: the task requires `(main module does not need module gopkg.in/yaml.v3)` exactly; this is not an import chain.
- `rg -n "gopkg.in/yaml.v3" go.mod go.sum` - BLOCKED / WITHDRAWN by DEC-143 - `go.sum:45:gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=` (terminated `exit_code: 0`)
  Immutable failed attempt. DEC-143 replaced this unreachable `go.sum` text search with the `go.mod` absence and `go mod why -vendor` build-graph checks above; it was not rerun.
- QA-plan commands beyond the Task 16 spec: none; the current spec includes every Task 16-specific absence gate.

### Second same-task continuation: immutable amended command failure (2026-08-16)

- `go mod why -vendor -m gopkg.in/yaml.v3` - BLOCKED - `(main module does not need to vendor module gopkg.in/yaml.v3)` (terminated `exit_code: 0`)
  Fresh complete output was two lines: `# gopkg.in/yaml.v3` then the recorded vendor-mode line. DEC-144 accepts that parenthesized line's wording but Task 16 requires the entire output to be one line and declares any extra line a real VAL-032 failure. No substitution or dependency/module-file change was attempted.

### Third same-task continuation: fresh amended command evidence (2026-08-16)

- `go mod tidy` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					78.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.080s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `ls internal/format` - PASS - `ls: cannot access 'internal/format': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`)
- `rg -n "internal/format" cmd internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -ni "writeTOON|writeYAML|writeMarkdown|yaml.Marshal" cmd internal` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -ni "yaml" go.mod` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `go mod why -vendor -m gopkg.in/yaml.v3` - PASS - `(main module does not need to vendor module gopkg.in/yaml.v3)` (terminated `exit_code: 0`)
  DEC-145 stanza evidence: first non-empty line `# gopkg.in/yaml.v3`; second and final non-empty line is the recorded parenthesized vendor-mode result. Exactly one stanza; no import-path or further non-empty line appeared.
- QA-plan commands beyond the Task 16 spec: none; the current spec includes every Task 16-specific absence gate.

## Deviations from spec

- No implementation or assertion-strength deviation. TDD preserved the initial red guard failure and the empty-directory follow-up failure; deleting the confirmed empty directory was the causal repair, after which the focused guard passed.
- Immutable QA-plan failure: `rg -n "gopkg.in/yaml.v3" go.mod go.sum` ended with the cited `go.sum` line and `exit_code: 0`. **WITHDRAWN by DEC-143:** it could not reach no-match while `go mod tidy` retained the accepted migration module's module-graph checksum. The withdrawn command was not rerun; DEC-143 replaces it with `rg -ni "yaml" go.mod` and `go mod why -vendor -m gopkg.in/yaml.v3`, and changes the four declared absence rows to the expected/observed PASS notation without weakening a VAL, threshold, or dependency rule.
- Immutable report-lint failures: the first run ended `FAIL report-lint found 13 failure(s)` with `exit_code: 13` because each PASS row had prose after the termination marker; the second ended `FAIL report-lint found 26 failure(s)` with `exit_code: 26` after report-only notation omitted the required backticks. The causal repair restored the validator's exact backticked termination notation with no trailing prose on PASS rows. DEC-143 then corrected the control to accept task-spec-declared expected-non-zero PASS rows; those failure tokens and terminations remain immutable evidence.
- Continuation report-lint failure: the amended control ended `FAIL report-lint found 6 failure(s)` with `exit_code: 6`. Diagnosis: `expected_nonzero_declaration_lines` in `project-docs/tools/report-lint.sh:84` returns every non-empty line under `## Expected Non-Zero Commands`, including the Task 16 section's instructional prose and wrapped descriptions; `validate_expected_nonzero_declarations` then rejects those lines as invalid grammar at `project-docs/tools/report-lint.sh:95`. This is a Task 16a control regression, not a report-row change; do not weaken, reformat, or delete the declared Task 16 evidence to evade it.
- Continuation blocker: the newly required `go mod why -vendor` command returned the real no-need vendor-mode line with `exit_code: 0`, but not the exact literal written in the task and QA plan. No source, dependency, module-file, command, threshold, or assertion change was attempted.
- Second same-task continuation blocker candidate: fresh `go mod why -vendor -m gopkg.in/yaml.v3` output was `# gopkg.in/yaml.v3` followed by `(main module does not need to vendor module gopkg.in/yaml.v3)`, with observed `exit_code: 0`. DEC-144 corrects the wording of the parenthesized line, but the current Task 16 contract separately requires the *entire* output to be exactly one parenthesized line and names any extra line as a real VAL-032 failure. This is a protected specification/output-contract contradiction, not an import chain; no dependency change, module-file hand edit, command substitution, or evidence rewording was attempted.
- DEC-145 resolved that protected complete-output contradiction without changing the command or VAL-032 assertion. This third continuation's successful `go mod why -vendor -m gopkg.in/yaml.v3` row preserves the observed header `# gopkg.in/yaml.v3` beside its parenthesized vendor-mode token, proves exactly one stanza and no import path/further non-empty line, and cites DEC-145 rather than rewording either earlier failed attempt.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- Removed the direct `gopkg.in/yaml.v3` dependency. No dependency was added or upgraded; the remaining transitive checksum is retained by `github.com/golang-migrate/migrate/v4@v4.19.1`.

## Commit hash

- `20813fc1052d0730113c3b1a2a2716a8a94a8488`

## Commit subject

- `docs(task): record Task 16 completion evidence`
