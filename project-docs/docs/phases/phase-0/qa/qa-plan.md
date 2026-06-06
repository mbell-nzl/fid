# Phase 0 QA Plan - Specification And Foundation

> Phase: 0 - Specification And Foundation
> Refreshed: 2026-06-05
> QA Lead: Codex

## 1. Scope And Source Documents

This plan verifies Phase 0's scaffold, quality tooling, and spec-audit readiness using the phase brief, validation contract, task specs, and quality/workflow specs.

- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/phases/phase-0/qa/validation-contract.md`
- `project-docs/docs/phases/phase-0/tasks/task-1.md`
- `project-docs/docs/phases/phase-0/tasks/task-2.md`
- `project-docs/docs/phases/phase-0/tasks/task-3.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/09_AI_WORKFLOW.md`

## 2. Validation Contract Refinement Check

| Check | Result | Notes |
|---|---|---|
| In-scope requirements referenced by at least one `VAL-NNN` | YES | Phase-0 roadmap, TRD, PRD, quality, and workflow claims are covered. |
| Every `VAL-NNN` has one validator role | YES | All assertions assign `QA Lead`. |
| Every `VAL-NNN` has a concrete evidence method | YES | Each row names command output, file-state proof, or source review. |
| Every `VAL-NNN` maps to a task/window or phase-end check | YES | Ownership is assigned in the contract. |
| Every `VAL-NNN` is one observable | YES | No row bundles unrelated behaviors. |
| Conditional/binary assertions name both states | YES | Build/help and placeholder/no-placeholder cases are explicit. |
| Parity assertions have deterministic proof | N/A | None are present. |
| Revisions requested from Tech Lead | None | Contract is ready for execution. |

## 3. Task Inventory

| Task | Focus | Covered `VAL-NNN` IDs | QA Notes |
|---|---|---|---|
| Task 1 | Go module and CLI walking skeleton | VAL-001, VAL-002, VAL-003, VAL-004, VAL-009 | Evidence is build/help output plus source-tree review. |
| Task 2 | Quality tooling and CI stub | VAL-001, VAL-005, VAL-006, VAL-009 | Evidence is Makefile/CI config plus task-gate command output. |
| Task 3 | Spec audit and phase readiness docs | VAL-004, VAL-007, VAL-008, VAL-009 | Evidence is audit output and readiness-doc file state. |

## 4. Planned Evidence By VAL

| VAL ID | Planned evidence | Validator | Window |
|---|---|---|---|
| VAL-001 | Source-tree proof of scaffold paths plus task completion reports. | QA Lead | tasks 1-2 |
| VAL-002 | `go build -o /tmp/fid ./cmd/fid` and `wc -c /tmp/fid` literal output. | QA Lead | task 1 |
| VAL-003 | `/tmp/fid --help` literal output showing root help/usage. | QA Lead | task 1 |
| VAL-004 | `go.mod`/`go.sum` review plus dependency-budget note in task 1 and task 3 completion reports. | QA Lead | tasks 1 and 3 |
| VAL-005 | `Makefile`, `.golangci.yml`, and `.github/workflows/ci.yml` inspection plus `make -n` dry-runs for the listed targets. | QA Lead | task 2 |
| VAL-006 | Task-2 gate command output: format, lint, coverage, build, and Makefile target runs. | QA Lead | task 2 / phase-end |
| VAL-007 | Validation-contract, task-spec, and phase-brief cross-reference plus mid-phase review file. | QA Lead | task 3 / phase-start |
| VAL-008 | `spec-audit.md` classifications plus `rg` audit output against PRD/TRD/spec/current-docs. | QA Lead | task 3 |
| VAL-009 | Source review of `cmd/`, `internal/`, `.github/`, and `Makefile` for no daemon/listener/MCP/remote-sync surface. | QA Lead | tasks 1-3 / phase-end |

## 5. Command Evidence

| Window | Commands to evidence | Expected result |
|---|---|---|
| Task 1 | `gofmt -l .`; `golangci-lint run ./...`; `go test ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go build -o /tmp/fid ./cmd/fid`; `wc -c /tmp/fid`; `/tmp/fid --help` | PASS, with literal stdout/stderr captured in the completion report. |
| Task 2 | `gofmt -l .`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go build -o /tmp/fid ./cmd/fid`; `make fmt-check`; `make lint`; `make test`; `make coverage`; `make build` | PASS, with Makefile targets shown to be thin wrappers. |
| Task 3 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go build -o /tmp/fid ./cmd/fid`; `rg -n "<TBD>|PLACEHOLDER|placeholder|Proposed" AGENTS.md project-docs/specs project-docs/docs/current`; `rg -n "net.Listen|http.ListenAndServe|grpc|mcp|remote sync|daemon" cmd internal .github Makefile` | PASS or expected no-match exit on the no-listener audit. |
| Phase end | `go test -race ./...`; `go test -tags integration -race ./...`; `go test -bench . -benchmem ./...`; `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go build -o /tmp/fid ./cmd/fid` | PASS, or BLOCKED with a documented environment reason. |

## 6. Command And Feasibility Validation

| Referenced step | Validation method | Result to record | Evidence |
|---|---|---|---|
| Task 1 build/help commands | Task-spec inspection plus repo root command design in TRD/Design. | PASS | Completion report output. |
| Task 2 Makefile targets and CI stub | `make -n` dry-run and file inspection once the files exist. | PASS or BLOCKED | Dry-run output and file references. |
| Task 3 audit commands | `rg` against repo-owned paths only; no held-out sources. | PASS or BLOCKED | `rg` output and `spec-audit.md`. |
| Human checklist feasibility | No Phase 0 `(human)` assertions are planned; the human still performs the final sign-off review and manual `/tmp/fid --help` smoke. | PASS | Final-review human-checklist readiness. |

## 7. Mid-Phase Review Cadence

- One scheduled mid-phase review after Task 2 completes, before Task 3 starts.
- Review window: Tasks 1-2 only.
- Review artifact: `project-docs/docs/phases/phase-0/qa/review-1.md`.
- Trigger earlier only if Task 1 or Task 2 produces a blocker that crosses task boundaries or changes the validation contract.
- No second scheduled mid-phase review; Task 3 is the audit/readiness task and phase-end review closes the phase.

## 8. Deterministic Checklist

| ID | Check | Covers | Evidence required |
|---|---|---|---|
| QA-01 | Confirm scaffold inventory exists and no runtime surface was introduced `(qa-lead)` | VAL-001, VAL-009 | File-state proof from tree/source review. |
| QA-02 | Capture build size and binary existence `(qa-lead)` | VAL-002 | `go build` and `wc -c /tmp/fid` output. |
| QA-03 | Capture root help output `(qa-lead)` | VAL-003 | Literal `/tmp/fid --help` output. |
| QA-04 | Confirm dependency budget notes are present `(qa-lead)` | VAL-004 | `go.mod`/`go.sum` review plus completion-report note. |
| QA-05 | Confirm Makefile/CI commands are concrete and runnable `(qa-lead)` | VAL-005, VAL-006 | File inspection plus `make -n` and task-gate output. |
| QA-06 | Confirm spec audit and phase-readiness docs are complete `(qa-lead)` | VAL-007, VAL-008 | `spec-audit.md`, `rg` audit output, and brief/contract cross-check. |

## 9. Human Verification Scope

- No Phase 0 `(human)` assertions are planned.
- The human sign-off gate still includes a manual `/tmp/fid --help` smoke and review of `final-review.md`.
- If a human-only judgment is added, it must be one observable per row, cite the parent `VAL-NNN`, and stay outside the deterministic checklist.

## 10. Phase-End Decision Rule

- Every `VAL-NNN` has planned evidence and a final-review outcome slot.
- Every required command is either run or recorded BLOCKED with reason.
- Every `(qa-lead)` item has reproducible evidence.
- Any missing or placeholder command, unsupported target, or unresolved spec gap blocks final review until escalated.
