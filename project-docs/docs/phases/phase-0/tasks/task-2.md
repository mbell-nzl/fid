# Phase 0 - Task 2: Quality Tooling And CI Stub

## Objective

Add concrete local and CI quality tooling so every Developer task has the same
format, lint, test, coverage, and build commands available.

## Requirements

- PRD refs: PRD-OPS-009
- TRD refs: TRD §1, TRD §3, TRD-PERF-007
- Design refs: N/A
- Quality refs: Quality Plan §§4-6, AI Workflow §7, AGENTS.md Commands,
  Roadmap Phase 0

## Validation Contract Coverage

- VAL-001: Phase 0 scaffold inventory includes Makefile, lint config, and CI stub.
- VAL-005: Quality tooling exposes concrete commands with no placeholder gates.
- VAL-006: The Phase 0 task-level quality gate is executable and green.
- VAL-009: Tooling introduces no runtime daemon, listener, remote sync, or MCP adapter.

## Scope

### Files To Create Or Modify

- `Makefile`
- `.golangci.yml`
- `.github/workflows/ci.yml`
- `project-docs/docs/phases/phase-0/tasks/task-2-completion.md`

### Files To Read For Context

- `AGENTS.md`
- `project-docs/specs/08_QUALITY_PLAN.md` §§4-6
- `project-docs/specs/09_AI_WORKFLOW.md` §7
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` Phase 0
- `project-docs/docs/phases/phase-0/qa/validation-contract.md`
- `project-docs/docs/phases/phase-0/tasks/task-1-completion.md`

### Out Of Scope

- Changing the required quality gate commands in `AGENTS.md` or the Quality
  Plan.
- Adding release-readiness tools that are explicitly deferred, such as SAST,
  `govulncheck`, `go-licenses`, or `gitleaks`.
- Adding generated code, database migrations, or product behavior.

## Acceptance Criteria

- [ ] `Makefile` provides concrete targets for format check, lint, test, coverage, and build using the commands documented in `AGENTS.md`. (VAL-005)
- [ ] `.golangci.yml` enables at minimum the linters required by Quality Plan §4.2, or records an explicit blocker if the installed `golangci-lint` version cannot support one. (VAL-005)
- [ ] CI stub runs the same concrete task-level commands and does not contain placeholder commands. (VAL-005)
- [ ] Required task-level gate commands pass in the scaffolded repository. (VAL-006)
- [ ] `go build -o /tmp/fid ./cmd/fid` remains green after tooling changes. (VAL-001)
- [ ] Tooling does not introduce a runtime service, listener, remote sync, or MCP adapter. (VAL-009)

## TDD Plan

- First failing test: N/A. This is tooling configuration, not product
  behavior.
- Additional tests: N/A.
- Regression risks: CI commands drifting from `AGENTS.md`, Makefile targets
  hiding narrower commands, or linter configuration using unsupported names.

## Documentation Updates

- `project-docs/docs/phases/phase-0/tasks/task-2-completion.md` - required
  completion report with command output and dependency-budget note.
- `README.md` - update only if build/test commands are documented there and
  need to match the new Makefile targets.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: Keep Makefile, CI, AGENTS, and Quality Plan
  command names aligned.
- Audit or release-readiness impact: This establishes the baseline automated
  gate surface; do not promote deferred release-readiness tools in this task.

## Exact Commands To Run

Run exactly:

```text
gofmt -l .
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go build -o /tmp/fid ./cmd/fid
make fmt-check
make lint
make test
make coverage
make build
```

The full task gate is required before task close. The Makefile commands are
also required because this task creates those targets. Every required command
must be recorded in the completion report as `PASS`, or as `BLOCKED` with
reason in a blocker report.

Conditional commands:

```text
go test -tags integration -race ./...
go test -bench . -benchmem ./...
```

List any conditional command required by this task and why:

- None unless integration or benchmark files are added unexpectedly.

## Notes For Developer

- Keep Makefile targets thin wrappers around documented commands.
- The CI file may be a stub, but its commands must be concrete and runnable.
- If `golangci-lint` is missing locally, do not replace it with `go vet`; write
  a blocker report.
