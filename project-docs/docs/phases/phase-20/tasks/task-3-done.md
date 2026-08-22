# Phase 20 - Task 3: Atomic project initialization - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 2
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Replaced direct `.fid/` creation with a unique sibling staging protocol: validate, stage strict bootstrap files, initialize the schema, append/project genesis, rebuild and byte-compare the projection, fsync files/directories, and atomically activate `.fid/`.
- Added deterministic, narrow stage-boundary failure hooks. Every injected pre-activation failure removes only its own staging directory and leaves the target absent.
- Migrated `fid init` to `project.init`: it accepts the global `--project`/`--json` transport, validates derived and explicit prefix grammar before staging, emits the exact JSON Response Envelope, and uses stdout/stderr correctly in Human mode.

## Files created or modified

- `internal/bootstrap/init.go`
- `internal/bootstrap/init_test.go`
- `cmd/fid/init.go`
- `cmd/fid/init_test.go`
- `cmd/fid/root_flags_test.go`
- `project-docs/docs/phases/phase-20/tasks/task-3-done.md`

## Documentation updated

- Completion report only; the accepted Phase 20 contract and configuration specification already document this behavior.

## Validation Contract Coverage

- VAL-001: PASS — `TestInitCommandUsesExactProjectInitResponseEnvelope` at `cmd/fid/init_test.go:125` asserts the exact five-key success envelope and stdout-only JSON channel; the final `/tmp/fid --project /tmp/fid-task3-final-cli-qIt4Jx --json init --prefix TASK3` entrypoint run emitted one such envelope.
- VAL-002: PASS — the collision half of `TestInitCommandUsesExactProjectInitResponseEnvelope` at `cmd/fid/init_test.go:125` asserts empty stdout plus the closed failure envelope, `project-already-initialized`, conflict category, project subject, and `{path}` details on stderr.
- VAL-004: PASS — `TestInitCommandUsesExactProjectInitResponseEnvelope` at `cmd/fid/init_test.go:125` initializes only its explicit temporary `--project` root; the real CLI command above reports that same absolute root in `result.project.root`.
- VAL-008: PASS (Task 3 portion) — `TestInitCommandCreatesWorkspaceFiles` at `cmd/fid/init_test.go:21` asserts Human success is stdout-only, while `TestInitCommandRejectsReinitWithoutChanges` at `cmd/fid/init_test.go:425` asserts Human collision failure is stderr-only.
- VAL-017: PASS (Task 3 mutation) — `TestInitCommandUsesExactProjectInitResponseEnvelope` at `cmd/fid/init_test.go:125` compares the returned Event ID and RFC3339-nanosecond timestamp with the appended `project.initialized` journal Event.
- VAL-039: PASS — `TestInitCommandCreatesWorkspaceFiles` at `cmd/fid/init_test.go:21` inventories strict config keys, rejects retired keys, and round-trips the generated opaque project ID and issue prefix with `config.ReadProjectConfig`.
- VAL-042: PASS — `TestInitializeFailureBeforeActivationLeavesNoWorkspace` at `internal/bootstrap/init_test.go:12` injects failures after validation, staging, schema initialization, genesis append, rebuild, comparison, fsync, and immediately before activation; every case asserts no `.fid/` and no staging residue. `TestInitCommandCreatesWorkspaceFiles` at `cmd/fid/init_test.go:21` proves successful activation leaves no stage residue, and `TestInitCommandRejectsReinitWithoutChanges` at `cmd/fid/init_test.go:425` snapshots the existing workspace before/after collision.
- VAL-043: PASS (Task 3 activation portion) — `TestInitCommandCreatesWorkspaceFiles` at `cmd/fid/init_test.go:21` reads the activated target config through the strict parser; `TestInitCommandEnforcesExplicitPrefixGrammarBeforeStaging` at `cmd/fid/init_test.go:267` proves the 16-character boundary activates and malformed explicit input returns `invalid-argument` with no target.
- VAL-044: PASS (Task 3 `project.init` row) — `TestInitCommandEnforcesExplicitPrefixGrammarBeforeStaging` at `cmd/fid/init_test.go:267` proves valid/invalid prefix boundaries reach or reject before staging; `TestInitCommandUsesExactProjectInitResponseEnvelope` at `cmd/fid/init_test.go:125` proves the closed `--project`, `--json`, and `--prefix` init request/result transport. Complete all-leaf closure remains Task 18.

## Tests run

Historical Task 3 test and entrypoint evidence is preserved below. Its original
direct-test, unchanged-surface, and gate records predate the required
termination syntax and are not used as current gate evidence.

> - `go test ./internal/bootstrap -run '^TestInitializeFailureBeforeActivationLeavesNoWorkspace$'` — EXPECTED RED — `FAIL`
> - `go test ./cmd/fid -run '^TestInitCommandUsesExactProjectInitResponseEnvelope$'` — EXPECTED RED — `FAIL`
> - `go test ./cmd/fid -run '^TestInitCommandEnforcesExplicitPrefixGrammarBeforeStaging$'` — EXPECTED RED — `FAIL`
> - `go test ./internal/bootstrap ./cmd/fid -run '^(TestInitializeFailureBeforeActivationLeavesNoWorkspace|TestInitCommand|TestRootFlagsRejectBeforeLegacyAndUnmigratedLeafHandlers$)'` — PASS — `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.297s`
> - `go test ./...` — PASS — `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.534s`
> - `go test ./cmd/fid -run '^TestInitCommandUsesExactProjectInitResponseEnvelope$'` — PASS — `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.043s`
> - `/tmp/fid --project /tmp/fid-task3-final-cli-qIt4Jx --json init --prefix TASK3` — PASS — `{"contract_version":"fid-cli.v1","ok":true,"operation":"project.init","result":{"journal_path":".fid/journal/2026-08-09.jsonl","project":{"id":"project-c76cb9c33e5fc1a2210fa32a6ea2ff13","issue_prefix":"TASK3","root":"/tmp/fid-task3-final-cli-qIt4Jx"}},"event":{"id":"EVT-20260809-000001","occurred_at":"2026-08-09T04:32:32.893716987Z"}}`
> - QA-plan applicability: QA-01 is covered by the exact gate below; QA-02, QA-03, QA-18, QA-19, and QA-21 have the envelope, explicit-target, strict-config/prefix, atomic-stage, and collision/no-residue evidence cited above. The Task-3-required integration command passed in the exact gate below.

## Exact quality-gate results

> - `gofmt -l cmd internal` — PASS — `exit 0 (no output)`
> - `golangci-lint run ./...` — PASS — `0 issues.`
> - `go test -race -coverprofile=coverage.out ./...` — UNVERIFIED — see Gate-evidence correction (DEC-131).
> - `go tool cover -func=coverage.out` — PASS — `total:									(statements)					77.1%`
> - `go test -tags integration -race ./...` — UNVERIFIED — see Gate-evidence correction (DEC-131).
> - `go build -o /tmp/fid ./cmd/fid` — PASS — `exit 0 (no output)`

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:									(statements)					78.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- Historical Task 3 implementation had no product or specification deviations.
- DEC-131 gate-evidence correction remains immutable historical context: the
  original full-suite records were UNVERIFIED because their silent tokens could
  not prove termination. This continuation preserves those records and adds
  fresh, session-captured PASS results without rewriting them.
- MF-01 report/control-conformance repair: the initial exact report-lint
  attempt failed with `FAIL report-lint found 21 failure(s)` (terminated
  `exit_code: 21`) because the report lacked the canonical typed workflow
  outcome and used legacy evidence punctuation/termination syntax. This
  report-only repair adds the required typed fields, preserves historical
  evidence verbatim as quoted context, and records fresh terminal evidence for
  every exact Task 3 command.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `42f5d05244eb67aaa9b389e2c612a5d667a16636`

## Commit subject

- `feat(project): stage atomic fid initialization`

## Gate-evidence correction (DEC-131)

- The full-suite `go test -race -coverprofile=coverage.out ./...` and `go test -tags integration -race ./...` records are UNVERIFIED. Each original recorded token `exit 0 (no output)` could not have come from a completed `go test` run: long runs may print nothing until their packages finish, while a killed or timed-out run is BLOCKED, never PASS with a silent token.
- Preserved original record: `go test -race -coverprofile=coverage.out ./...` — PASS — original recorded token `exit 0 (no output)`.
- Preserved original record: `go test -tags integration -race ./...` — PASS — original recorded token `exit 0 (no output)`.
- `task-3c-done.md`'s current-tree full-gate run is the Phase 20 verified baseline; these historical commands were not re-run at Task 3's commit.
- Task 3 remains complete, but its own full-suite gate evidence is UNVERIFIED as of its task commit.
