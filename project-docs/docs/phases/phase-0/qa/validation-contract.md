# Phase 0 Validation Contract - Specification And Foundation

> Phase: 0 - Specification And Foundation
> Author: Tech Lead (draft)
> Last refined: Not yet refined by QA Lead

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 0 is complete. The validation contract is the specification of what must
be proven; the QA plan is the strategy for how to prove it.

If a required validation assertion cannot be satisfied, stop and escalate
instead of weakening the assertion or silently adapting the spec.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | Roadmap Phase 0 "Repo scaffold"; TRD §§2-3 | The repository contains the Phase 0 scaffold inventory: Go module, `cmd/fid/`, `internal/`, `migrations/`, `testdata/`, Makefile, `.golangci.yml`, and CI stub. | File-state proof from `find`/source tree plus task completion report. | QA Lead | tasks 1-2 | All required scaffold paths exist with no unrelated v0.1 feature directories omitted from the planned architecture. | Any required path is missing, misplaced, or replaced with an unrecorded boundary. | Draft | Runtime `.fid/` workspace creation is v0.1 `fid init`, not Phase 0 scaffold. |
| VAL-002 | Roadmap Phase 0 "Walking skeleton"; PRD-OPS-009; TRD-PERF-007 | The CLI build produces a runnable `/tmp/fid` binary under 30 MB. | Literal output from `go build -o /tmp/fid ./cmd/fid` and `wc -c /tmp/fid`. | QA Lead | task 1 | Build exits 0 and byte count is below 31,457,280. | Build fails, binary is missing, or binary size is 31,457,280 bytes or larger. | Draft | Cross-compilation and release packaging are out of scope. |
| VAL-003 | Roadmap Phase 0 "Walking skeleton"; Design Plan §§2, 3.1 | `/tmp/fid --help` exits 0 and displays root command help for the `fid` CLI. | Literal command output excerpt from `/tmp/fid --help`. | QA Lead | task 1 | Command exits 0 and output identifies `fid` with available help/usage text. | Command exits non-zero, panics, or emits no usable root help. | Draft | No v0.1 subcommand behavior is required beyond the root skeleton. |
| VAL-004 | TRD §1; TRD §9; AGENTS.md Dependency Budget | Phase 0 introduces only dependencies accepted for Phase 0 use, and any added dependency has documented budget rationale. | `go.mod`/`go.sum` review plus dependency-budget note in completion report. | QA Lead | tasks 1 and 3 | Added dependencies are limited to accepted stack decisions used by the scaffold, with rationale and license noted; future Proposed dependencies are not added. | An unapproved dependency is added, an unused future dependency is added without rationale, or license/rationale is missing. | Draft | Future frontmatter and TUI dependencies remain out of scope until their phases. |
| VAL-005 | Roadmap Phase 0 "Agent workflow validation"; Quality Plan §§4-6; AI Workflow §7; AGENTS.md Commands | Quality tooling files expose concrete Go commands matching the documented task gate, with no placeholder gate commands. | File-state proof from Makefile, `.golangci.yml`, CI workflow, and AGENTS command review. | QA Lead | task 2 | Tooling references concrete commands for format, lint, tests, coverage, and build; no required command is a placeholder. | Any required gate is missing, placeholder-only, or narrower than the documented task gate without recorded exception. | Draft | Release-readiness tools such as SAST, govulncheck, license checks, and gitleaks may remain deferred as documented. |
| VAL-006 | Roadmap Phase 0 exit criteria; Quality Plan §5.1; AI Workflow §7; AGENTS.md Commands | The Phase 0 task-level quality gate is executable and green in the scaffolded repository. | Literal command output for `gofmt -l .`, `golangci-lint run ./...`, `go test -race -coverprofile=coverage.out ./...`, and `go tool cover -func=coverage.out`. | QA Lead | task 2 / phase-end | All four commands exit 0; `gofmt -l .` lists no files. | Any command is unavailable, exits non-zero, or is skipped without a blocker report. | Draft | Integration, benchmark, and release-readiness commands are phase-end or future gates unless a task explicitly requires them. |
| VAL-007 | Quality Plan §5.4; AI Workflow §2.2 and §5.1 | Phase planning artifacts map validation assertions to ordered task specs and the current phase brief. | File-state proof from validation contract, task specs, and `project-docs/docs/current/phase_brief.md`. | QA Lead | task 3 / phase-start | Every `VAL-NNN` is owned by a task/window, every task cites covered `VAL-NNN` IDs, and the current phase brief includes the contract index. | Any task lacks validation coverage, any assertion lacks an owner, or the brief omits the contract index. | Draft | QA plan authoring is the QA Lead's next step, not a Tech Lead planning deliverable. |
| VAL-008 | Roadmap Phase 0 "Spec audit" and exit criteria; TRD §1; PRD §7 | Phase 1 can be planned without unresolved Phase 0 stack choices or command placeholders. | `spec-audit.md` plus review of PRD/TRD open questions and command text. | QA Lead | task 3 | Audit records no blocking Phase 0 `<TBD>`, Proposed stack dependency, or placeholder command; future questions are classified out of scope. | A Phase 0 stack choice, required command, or spec dependency remains unresolved. | Draft | Future product questions in PRD §7 are not resolved unless they block Phase 1 planning. |
| VAL-009 | PRD-NFR-005; PRD-NFR-006; TRD-SEC-001; TRD-ARCH-013 | The Phase 0 scaffold does not introduce a daemon, network listener, remote sync path, or MCP adapter. | Source review and file-state proof. | QA Lead | tasks 1-3 / phase-end | No scaffolded code or CI step opens a network listener, transmits project data, or implements MCP/remote sync. | Any listener, remote transport, or unapproved external service integration is added. | Draft | Dependency downloads by Go tooling during build are not fid runtime behavior. |

## Coverage Rules

- Every Phase 0 roadmap item is referenced by at least one `VAL-NNN`.
- Every `VAL-NNN` has one validator role and one owning task/window.
- Every Developer task spec cites the `VAL-NNN` IDs it advances.
- QA Lead must refine evidence into a QA plan before task execution begins.

## Outcomes

Phase-end final review records every `VAL-NNN` outcome as `PASS`, `FAIL`,
`BLOCKED`, or `CARRY-FORWARD`. `PASS` without evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-05 | VAL-001..VAL-009 | Created | Phase 0 planning kickoff | N/A |
