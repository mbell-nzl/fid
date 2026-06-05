# Phase 0 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes,
> and wrap-up decisions. DEC numbering is global and never resets.

## Decisions

### DEC-001: Redirect Go/golangci-lint caches into gitignored in-repo dirs

- Date: 2026-06-06
- Author: Tech Lead
- Context: Task 2 Developer hit a hard blocker — `golangci-lint run ./...`
  (and `go build`/`go test`) write to caches outside the repo (`GOCACHE`
  default `~/.cache/go-build`, `GOLANGCI_LINT_CACHE` default
  `~/.cache/golangci-lint`), which are read-only under Codex
  `sandbox_mode = "workspace-write"`. Phase Runner correctly stopped (blocker
  + missing `task-2-completion.md`).
- Decision: Fix at the **environment/harness layer**, not in the Task 2
  deliverables. Added `[shell_environment_policy.set]` to `.codex/config.toml`
  pointing `GOCACHE` → `/home/mark/projects/fid/.gocache` and
  `GOLANGCI_LINT_CACHE` → `/home/mark/projects/fid/.golangci-cache`, and
  gitignored `/.gocache/`, `/.golangci-cache/`, `/.tmp/`. Mirrors the proven
  workaround in the sibling `whet` project.
- Rationale: The sandbox cache restriction is an agent-environment property,
  not a project tooling concern. Baking cache paths into `Makefile`/CI/
  `.golangci.yml` would pollute portable deliverables and contradicts Task 2's
  "thin wrappers" constraint. Task 2 scope, spec, and gate commands are
  unchanged.
- Consequences / caveats:
  - `shell_environment_policy.set` applies at session start, so the Developer
    must be **re-spawned in a fresh session** for the vars to take effect.
  - Cache path is hardcoded absolute (Codex does not expand `$PWD` inside
    `set`); needs updating if the repo moves.
  - At phase wrap, assess whether `09_AI_WORKFLOW.md` and/or `AGENTS.md` should
    document this sandbox requirement so future phases do not re-hit it.
- Files touched: `.codex/config.toml`, `.gitignore`.

### DEC-002: Retire TL-Lite warm-up; enable autonomous resume permanently

- Date: 2026-06-06
- Author: Tech Lead (Human-directed)
- Context: WF-10 introduced a warm-up safeguard requiring Phase Runner to stop
  for Human approval on every TL-Lite `VERDICT: RESUME` for the first two
  phases, as a guard while TL-Lite was new to the workflow. The Human is now
  satisfied TL-Lite functions correctly and directed permanent removal of the
  safeguard.
- Decision: Warm-up is retired permanently. Autonomous TL-Lite resume is
  ENABLED. Phase Runner records the TRG-NNN and continues on `VERDICT: RESUME`;
  `VERDICT: BLOCK` still always stops. No warm-up tally is maintained at phase
  wrap. The wrap-up triage audit (confirm RESUME verdicts in hindsight, flag any
  that should have escalated) is retained — only the stop-on-RESUME gate and the
  phase tally are removed.
- Files touched: `project-docs/docs/current/decisions_active.md`,
  `project-docs/specs/09_AI_WORKFLOW.md`, `.codex/agents/phase_runner.toml`,
  `.codex/agents/tech_lead.toml`, `.claude/agents/phase-runner.md`,
  `.claude/agents/tech-lead.md`,
  `project-docs/templates/phase-decisions-template.md`,
  `project-docs/templates/triage-log-template.md`.
- Spec sync: `09_AI_WORKFLOW.md` touched (WF-10 warm-up text removed).

### DEC-003: Standardize on golangci-lint v2; drop `gosimple` linter name

- Date: 2026-06-06
- Author: Tech Lead
- Context: Task 2 and Task 3 lint gates never ran green. Root cause is a
  three-layer version contradiction, not a local-only environment issue:
  1. `08_QUALITY_PLAN.md` LINT-002 mandates `gosimple` in the linter set — a
     golangci-lint **v1** linter name.
  2. `.golangci.yml` declares `version: "2"` (with a v2-only `formatters:`
     block) yet enables `gosimple`. In golangci-lint **v2**, `gosimple` no
     longer exists as a standalone linter (its checks are part of
     `staticcheck`), so the installed v2.12.2 correctly errors
     `unknown linters: 'gosimple'`.
  3. `.github/workflows/ci.yml` pins `golangci-lint@v1.64.8` via the **v1**
     module path (`.../cmd/golangci-lint`), which cannot correctly read a
     `version: "2"` config either. CI is mismatched in the opposite direction.
  TRG-001 RESUME ("local binary mismatch only; CI pins a compatible version")
  was incorrect in hindsight — this is a committed repo-config defect that
  blocked every task gate. Flag at wrap-up; see Workflow notes.
- Decision: Standardize the project on **golangci-lint v2**.
  1. Spec: amend LINT-002 to drop `gosimple` from the mandated set and note its
     checks are covered by `staticcheck` under v2. (Done in this wrap action.)
  2. `.golangci.yml`: remove the `- gosimple` line. Remaining names
     (`errcheck`, `govet`, `staticcheck`, `unused`, `ineffassign`, `misspell`,
     + `gofmt`/`goimports` formatters) are all valid v2 names.
  3. CI: install a v2 binary via the v2 module path, pinned to match the dev
     toolchain: `go install
     github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2`.
- Process: Task 2 is reopened — its gate was never green, so the tooling slice
  is not actually complete. A Developer applies fixes (2) and (3), runs the full
  Task 2 gate to green, and rewrites `task-2-completion.md` with real PASS
  evidence. Phase Runner then re-runs Task 3 (gate now green → `spec-audit.md`
  produced) and the phase-end QA pass.
- Spec sync: `08_QUALITY_PLAN.md` touched (LINT-002). Confirm no other spec
  pins `gosimple` (audit found only LINT-002).
- Files (TL, now): `project-docs/specs/08_QUALITY_PLAN.md`.
- Files (Developer, reopened Task 2): `.golangci.yml`,
  `.github/workflows/ci.yml`, `project-docs/docs/phases/phase-0/tasks/task-2-completion.md`.

## Carry-Forward Disposition At Kickoff

| ID | Disposition | Rationale | Follow-up |
|---|---|---|---|
| CF-002 | Closed as satisfied | The canonical spec suite now exists under `project-docs/specs/`, and Phase 0 is planned from the roadmap plus PRD/TRD/design/quality/workflow specs. | Task 3 performs a Phase 0 spec audit before Phase 1 planning. |

## Spec Sync Verdict

Pending phase wrap-up.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | pending | Phase wrap-up required |
| `01_PRD.md` | pending | Phase wrap-up required |
| `02_TRD.md` | pending | Phase wrap-up required |
| `03_TRD_SCHEMA.md` | pending | Phase wrap-up required |
| `04_TRD_CONFIG.md` | pending | Phase wrap-up required |
| `05_DESIGN_PLAN.md` | pending | Phase wrap-up required |
| `06_DESIGN_PLAN_CLI.md` | pending | Phase wrap-up required |
| `07_DESIGN_PLAN_FLOWS.md` | pending | Phase wrap-up required |
| `08_QUALITY_PLAN.md` | pending | Phase wrap-up required |
| `09_AI_WORKFLOW.md` | pending | Phase wrap-up required |
| `10_IMPLEMENTATION_ROADMAP.md` | pending | Phase wrap-up required |

## Workflow Notes

Pending phase wrap-up.

### Deep Retrospective Recommendation

Pending phase wrap-up.
