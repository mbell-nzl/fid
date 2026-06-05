# Phase Brief - Phase 0 (Planning)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.

## Phase Goal

Phase 0 makes fid implementable by agents: a Go/Cobra repository scaffold,
walking CLI binary, concrete quality tooling, and a final spec audit proving
Phase 1 can be planned without unresolved stack or command choices.

## Status

**Planning.** Tech Lead has drafted the validation contract and ordered task
specs. Next action: human review, then QA Lead creates the Phase 0 QA plan
before any Developer task execution.

## Task Index

Full task specs and completion reports live under
`../phases/phase-0/tasks/`.

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Go Module And CLI Walking Skeleton | Pending | None |
| 2 | Quality Tooling And CI Stub | Pending | Task 1 |
| 3 | Spec Audit And Phase Readiness Docs | Pending | Tasks 1-2 |

## Validation Contract

Full contract:
`../phases/phase-0/qa/validation-contract.md`.

| VAL ID | Source | Owning Task/Window | Validator | Status |
|---|---|---|---|---|
| VAL-001 | Roadmap Phase 0, TRD §§2-3 | tasks 1-2 | QA Lead | Draft |
| VAL-002 | PRD-OPS-009, TRD-PERF-007 | task 1 | QA Lead | Draft |
| VAL-003 | Roadmap Phase 0, Design §§2, 3.1 | task 1 | QA Lead | Draft |
| VAL-004 | TRD §§1, 9; dependency budget | tasks 1 and 3 | QA Lead | Draft |
| VAL-005 | QP §§4-6, WF §7, AGENTS commands | task 2 | QA Lead | Draft |
| VAL-006 | Roadmap exit, QP §5.1, WF §7 | task 2 / phase-end | QA Lead | Draft |
| VAL-007 | QP §5.4, WF §§2.2, 5.1 | task 3 / phase-start | QA Lead | Draft |
| VAL-008 | Roadmap spec audit, TRD §1, PRD §7 | task 3 | QA Lead | Draft |
| VAL-009 | PRD-NFR-005/006, TRD-SEC-001 | tasks 1-3 / phase-end | QA Lead | Draft |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| Validation contract | `../phases/phase-0/qa/validation-contract.md` | Draft |
| QA plan | `../phases/phase-0/qa/qa-plan.md` | Missing |
| Mid-phase reviews | `../phases/phase-0/qa/review-M.md` | QA cadence pending |
| Final review | `../phases/phase-0/qa/final-review.md` | Pending |
| Human checklist | `../phases/phase-0/qa/human-checklist.md` | Pending |

## Active Constraints

- Do not implement v0.1 `fid init`, journal, SQLite, issue lifecycle, policy,
  workflow, telemetry, reports, TUI, MCP, or remote sync in Phase 0.
- Every Developer task must run the exact commands listed in its task spec.
- Do not add future Proposed dependencies for frontmatter or TUI in Phase 0.
- Stop with a blocker report if a required command or Phase 0 stack choice is
  not concrete.

## Interfaces / Subsystems In Scope

- `cmd/fid/` - root CLI skeleton only.
- `internal/`, `migrations/`, `testdata/` - scaffold directories only.
- Makefile, `.golangci.yml`, CI stub - concrete quality gate surface.
- Phase 0 docs under `project-docs/docs/phases/phase-0/`.

## Known Risks

- **Tool availability.** `golangci-lint` may be missing locally; record BLOCKED
  rather than substituting a narrower lint command.
- **Scope creep.** The skeleton can easily drift into v0.1 init or storage
  work; task specs explicitly forbid that.

## Open Questions

- None blocking Phase 0. PRD §7 future product questions remain deferred.

## Success / Exit Criteria

- [ ] All three tasks complete with required command evidence.
- [ ] `/tmp/fid` builds under 30 MB and `/tmp/fid --help` exits 0.
- [ ] Quality gate commands are concrete, documented, and green.
- [ ] Spec audit finds no Phase 1 planning blocker.
- [ ] Final QA review has no unresolved blocking findings.

## Next Planning Checkpoint

After Phase 0 QA and human acceptance, Tech Lead wraps Phase 0 and plans the
first v0.1 persistence foundation slice.
