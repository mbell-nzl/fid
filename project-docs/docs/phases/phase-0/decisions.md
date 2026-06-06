# Phase 0 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes,
> and wrap-up decisions. DEC numbering is global and never resets.

## Decisions

### DEC-001: Redirect Go/golangci-lint caches into gitignored in-repo dirs

- Date: 2026-06-06
- Author: Tech Lead
- Context: Task 2 Developer hit a hard blocker — `golangci-lint run ./...`
  (and `go build`/`go test`) write to caches outside the repo (`GOCACHE`
  default `~/.cache/go-build`, `GOMODCACHE` default `~/go/pkg/mod`,
  `GOLANGCI_LINT_CACHE` default `~/.cache/golangci-lint`), which are read-only
  under Codex `sandbox_mode = "workspace-write"`. Phase Runner correctly
  stopped (blocker + missing `task-2-completion.md`).
- Decision: Fix at the **environment/harness layer**, not in the Task 2
  deliverables. Added `[shell_environment_policy.set]` to `.codex/config.toml`
  pointing `GOCACHE` -> `/home/mark/projects/fid/.gocache`,
  `GOMODCACHE` -> `/home/mark/projects/fid/.gomodcache`, `GOPATH` ->
  `/home/mark/projects/fid/.gopath`, `GOLANGCI_LINT_CACHE` ->
  `/home/mark/projects/fid/.golangci-cache`, and `XDG_CACHE_HOME` ->
  `/home/mark/projects/fid/.cache`. The directories are gitignored, including
  `/.tmp/`. Mirrors the proven cache-directory workaround in the sibling
  `whet` project and extends it for module/cache paths that Go uses during
  dependency fetches or tool installs.
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
  - `go install` still needs normal approval when it downloads from the network
    or changes tool dependencies; this decision only makes the resulting caches
    and module writes sandbox-writable.
  - `AGENTS.md` now documents the operational rule. At phase wrap, assess
    whether `09_AI_WORKFLOW.md` should also carry the project-local cache
    convention.
- Files touched: `.codex/config.toml`, `.gitignore`, `AGENTS.md`.

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

### DEC-004: Bump Cobra patch version to avoid a stale pin

- Date: 2026-06-06
- Author: Tech Lead
- Context: The scaffolded direct dependency `github.com/spf13/cobra v1.10.1`
  was an incidental Phase 0 pin, not a project requirement. The local tool
  environment already had Cobra `v1.10.2` available, and there was no recorded
  reason to keep the older patch version.
- Decision: Upgrade the accepted Cobra runtime dependency from `v1.10.1` to
  `v1.10.2`, keeping the same framework and the same CLI surface. This is a
  patch-level dependency update, not a new runtime dependency.
- Dependency-budget rationale:
  - Need: avoid carrying an unnecessary older patch pin while staying on the
    already accepted Cobra framework from TRD §1/§9.
  - Alternatives considered: keep `v1.10.1` and require an initial networked
    `go mod download` into `.gomodcache/`; vendor dependencies; or allow
    agents to use `~/go/pkg/mod`. The patch bump is narrower than vendoring and
    the cache policy still keeps writes inside the workspace.
  - License check: Cobra `v1.10.2` includes an Apache-2.0 license, which is
    compatible with the default dependency policy.
- Consequences / caveats:
  - `go.mod` and `go.sum` now record Cobra `v1.10.2`.
  - No product behavior, CLI contract, or dependency family changes.
  - A clean workspace still needs to populate `.gomodcache/` once; DEC-001
    controls where that write goes.
- Files touched: `go.mod`, `go.sum`.

### DEC-005: Keep the sandbox cache-redirection rule in AGENTS.md, not in `09_AI_WORKFLOW.md`

- Date: 2026-06-06
- Author: Tech Lead
- Context: A Phase 0 scheduled commitment asked, at wrap, whether DEC-001
  (sandbox cache redirection into gitignored repo-local dirs) also belongs in
  `09_AI_WORKFLOW.md`. The operational rule already lives in `AGENTS.md`
  (lines 170-179), which `CLAUDE.md` imports, so every agent session already
  receives it.
- Decision: Do **not** copy the cache-redirection rule into
  `09_AI_WORKFLOW.md`. It is an agent-environment / host-sandbox property, not a
  workflow-process step. `09_AI_WORKFLOW.md` describes the canonical agent
  workflow and is consumed as-is by the build pipeline; host caching belongs in
  the grounding file (`AGENTS.md`). This closes the scheduled commitment.
- Consequences: The cache convention has a single home (`AGENTS.md` +
  `.codex/config.toml`, recorded in DEC-001). `09_AI_WORKFLOW.md` spec-sync for
  this item is `verified-no-change`; the file was still touched separately for
  DEC-002 (warm-up removal).
- Files touched: none (decision to make no change; recorded here and in the
  spec-sync table).

## Carry-Forward Disposition At Kickoff

| ID | Disposition | Rationale | Follow-up |
|---|---|---|---|
| CF-002 | Closed as satisfied | The canonical spec suite now exists under `project-docs/specs/`, and Phase 0 is planned from the roadmap plus PRD/TRD/design/quality/workflow specs. | Task 3 performs a Phase 0 spec audit before Phase 1 planning. |

## Spec Sync Verdict

Recorded at phase wrap-up (2026-06-06). `verified-no-change` is a positive
review assertion: the spec was walked against in-phase DEC-NNN, deviations, and
do-now QA recommendations and confirmed to need no edit.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | verified-no-change | No Phase 0 decision or deviation touched charter scope. |
| `01_PRD.md` | verified-no-change | VAL-008 spec audit found PRD §7 open questions out of Phase 0 scope; no requirement changed. |
| `02_TRD.md` | verified-no-change | Stack (Cobra) was accepted, not changed; Proposed frontmatter/TUI deps stay deferred. DEC-004 is a patch bump within the accepted dependency. |
| `03_TRD_SCHEMA.md` | verified-no-change | No schema work in Phase 0 (scaffold only). |
| `04_TRD_CONFIG.md` | verified-no-change | No config-surface work in Phase 0. |
| `05_DESIGN_PLAN.md` | verified-no-change | Root CLI help skeleton matches Design §§2, 3.1; no design change required. |
| `06_DESIGN_PLAN_CLI.md` | verified-no-change | No v0.1 subcommand surface implemented; CLI design untouched. |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | No flow implementation in Phase 0. |
| `08_QUALITY_PLAN.md` | touched: see DEC-003 | LINT-002 amended to drop `gosimple` (folded into `staticcheck` under golangci-lint v2). Verified in spec lines 173-178. |
| `09_AI_WORKFLOW.md` | touched: see DEC-002 | WF-10 warm-up safeguard removed; autonomous TL-Lite resume documented. DEC-005 confirms the DEC-001 cache rule stays in `AGENTS.md` (no further 09 change). |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 0 met roadmap items and exit criteria as written; no roadmap edit needed. |

## Workflow Notes

**What worked well.**

- The validation contract → task spec → completion report → final review chain
  held: all nine VAL assertions reached `PASS` with reproducible command
  evidence, and the final review verified all three task commit hashes.
- Failures were caught and root-caused rather than papered over. The lint gate
  never going green forced a real diagnosis (DEC-003) instead of a narrowed
  command.
- Recovering by reopening Task 2 (rather than bolting fixes onto Task 3) kept
  the "tooling slice is not done until its gate is green" invariant honest.

**What was inefficient.**

- The golangci-lint v1/v2 contradiction was a committed, three-layer defect
  (spec LINT-002 + `.golangci.yml` + CI pin) that masqueraded as a local
  environment mismatch for a full triage cycle before being correctly
  diagnosed. That cost a reopen and a re-run of Task 3.
- The sandbox cache blocker (DEC-001) and the cobra pin (DEC-004) were
  environment papercuts that surfaced mid-phase rather than during scaffold
  setup.

**Were tasks the right size?**

- Yes. Three ordered tasks (scaffold → tooling → audit) were appropriately
  one-session each. The Task 2 reopen was caused by a latent repo defect, not
  by the task being too large.

**Did QA cadence feel right?**

- Adequate for a planning/foundation phase. One mid-phase review (review-1,
  CONTINUE) plus the phase-end final review was proportionate to three tasks.
  The mid-phase review did not catch the v1/v2 defect because the gate had not
  yet been forced green — a deterministic "the gate must actually run green"
  check earlier would have surfaced it sooner.

**Were token-efficiency opportunities missed?**

- Minor. The validation contract, brief, and decisions stayed within budget.
  The main avoidable cost was re-triaging TRG-001 instead of treating a
  persistently red gate as a repo defect on first encounter.

### TRG / Triage Audit

- **TRG-001 (RESUME) — mis-triage in hindsight, flagged.** It treated a
  committed repo-config defect (golangci-lint v2 rejecting `gosimple`, with a
  mismatched CI v1 pin) as a local-only installed-binary mismatch. Under the
  TL-Lite escalation triggers this should have escalated: it ultimately changed
  a quality gate, spanned multiple files/subsystems (`.golangci.yml`, CI, spec
  LINT-002), and was a quality-gate anomaly not anticipated by the QA plan
  (triggers 4, 5, and QA-anomaly). The repo is now green via DEC-003; no
  carry-forward remains, but the heuristic lesson stands (see deep retro).
- Autonomous TL-Lite resume confirmed in effect (DEC-002); no warm-up tally to
  maintain.

### Deep Retrospective Recommendation

Recommend a **lightweight** retrospective note rather than a deep retro. Phase 0
is small and its single substantive lesson is already captured. Carry one
heuristic into Phase 1 triage guidance: **a quality-gate command that stays red
across runs should be treated as a candidate repo/config defect (escalate),
not assumed to be a local environment mismatch (resume), unless the failure is
demonstrably tied to an unsupported locally-installed tool version.** Consider
encoding this as a triage check when v0.2 governance tooling formalizes triage.
