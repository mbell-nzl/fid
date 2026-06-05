# Phase 0 - Task 1: Go Module And CLI Walking Skeleton

## Objective

Create the minimal Go module and `fid` CLI skeleton so the repository can
build a runnable single binary before feature work begins.

## Requirements

- PRD refs: PRD-OPS-009, PRD-NFR-005, PRD-NFR-006
- TRD refs: TRD §1, TRD §§2-3, TRD-PERF-007, TRD-SEC-001, TRD-ARCH-013
- Design refs: Design Plan §§2, 3.1
- Quality refs: Roadmap Phase 0, Quality Plan §5.1, AI Workflow §7

## Validation Contract Coverage

- VAL-001: Phase 0 scaffold inventory exists.
- VAL-002: CLI build produces a `/tmp/fid` binary under 30 MB.
- VAL-003: `/tmp/fid --help` exits 0 and displays root command help.
- VAL-004: Added dependencies are accepted and documented.
- VAL-009: Scaffold introduces no daemon, listener, remote sync, or MCP adapter.

## Scope

### Files To Create Or Modify

- `go.mod`
- `go.sum`
- `cmd/fid/`
- `internal/`
- `migrations/`
- `testdata/`
- `README.md` if the build/help command examples need updating
- `project-docs/docs/phases/phase-0/tasks/task-1-completion.md`

### Files To Read For Context

- `AGENTS.md`
- `project-docs/specs/02_TRD.md` §§1-3, §§6-10
- `project-docs/specs/05_DESIGN_PLAN.md` §§2, 3.1
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` Phase 0
- `project-docs/docs/phases/phase-0/qa/validation-contract.md`

### Out Of Scope

- Implementing `fid init`, journal writes, SQLite migrations, or `.fid/`
  workspace creation.
- Implementing v0.1 issue, claim, workflow, policy, telemetry, or report
  commands.
- Adding future Proposed dependencies for frontmatter, TUI, MCP, or remote
  sync surfaces.

## Acceptance Criteria

- [ ] `go.mod` defines the Go module and uses a Go version compatible with the local toolchain. (VAL-001)
- [ ] `cmd/fid` builds into `/tmp/fid` and the binary is under 30 MB. (VAL-002)
- [ ] `/tmp/fid --help` exits 0 and displays root help/usage for the CLI. (VAL-003)
- [ ] Any dependency added to `go.mod` is accepted by TRD §1/§9, used by the scaffold, and documented in the completion report dependency-budget note. (VAL-004)
- [ ] Empty scaffold directories needed for future phases are preserved with explicit placeholders such as `.gitkeep`, not product stubs. (VAL-001)
- [ ] No code path opens a network listener, starts a daemon, transmits project data, or implements MCP/remote sync. (VAL-009)

## TDD Plan

- First failing test: N/A for this scaffold task. Roadmap Phase 0 explicitly
  accepts a walking skeleton verified by build/help output before feature
  tests exist.
- Additional tests: Do not add product tests in this task unless needed to
  stabilize the root command behavior; if added, keep them behavior-oriented
  and record the deviation from the zero-test scaffold expectation.
- Regression risks: Overbuilding v0.1 command behavior, importing unused
  dependencies, or letting CLI setup carry business logic.

## Documentation Updates

- `README.md` - update only if the new skeleton changes the documented build
  or help command.
- `project-docs/docs/phases/phase-0/tasks/task-1-completion.md` - required
  completion report with validation coverage and gate evidence.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: None; only the CLI root exists.
- Audit or release-readiness impact: Dependency-budget note required because
  this task may add runtime dependencies.

## Exact Commands To Run

Run exactly:

```text
gofmt -l .
golangci-lint run ./...
go test ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go build -o /tmp/fid ./cmd/fid
wc -c /tmp/fid
/tmp/fid --help
```

The full task gate is required before task close. The added build, size, and
help commands are required for VAL-002 and VAL-003. Every required command must
be recorded in the completion report as `PASS`, or as `BLOCKED` with reason in
a blocker report.

Conditional commands:

```text
go test -tags integration -race ./...
go test -bench . -benchmem ./...
```

List any conditional command required by this task and why:

- None unless integration files or benchmark files are added unexpectedly.

## Notes For Developer

- Keep `cmd/fid/` as a thin adapter. Do not put domain behavior in the root
  command.
- Prefer Cobra for the root command because TRD §1 accepts Cobra and pflag for
  the command tree.
- If the local environment lacks `golangci-lint`, record a blocker report
  instead of substituting a narrower command.
