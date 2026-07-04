# fid - Project Context

## What Is This Project?

fid is a local-first, agent-first issue tracker for AI-assisted software
development. It provides a structured, queryable knowledge graph over the
artifacts an agentic development workflow produces. Rows hold metadata and
identity; markdown files hold prose. SQLite-backed CLI.

## Tech Stack

Concrete choices are recorded in `project-docs/specs/02_TRD.md` §1.

- Primary language/runtime: Go
- Application framework(s): `spf13/cobra` + `spf13/pflag` for the CLI command
  tree; no server framework in v1.
- Data/storage: SQLite via `modernc.org/sqlite` (pure Go), with
  journal-canonical persistence and a derived DB projection — see
  `docs/adr/0001` and `project-docs/specs/02_TRD.md`.
- Interface surfaces: CLI
- Test tools: stdlib `testing` + `stretchr/testify`; integration tests use
  temporary SQLite databases; benchmarks use `go test -bench`.
- Build/tool runner: `go build`, `go test`, and Makefile targets added during
  Phase 0 scaffolding.
- Deployment/runtime environment: local-first, single-user

If any item remains `<TBD>`, agents may still help refine specifications, but
they must not begin implementation that depends on that choice.

## Architecture

- User-facing adapters: `cmd/fid/`, `internal/format/`,
  `internal/export/`.
- Core/domain/business logic: `internal/issues/`, `internal/claims/`,
  `internal/workflow/`, `internal/policy/`, `internal/qa/`,
  `internal/validation/`, `internal/templates/`, `internal/search/`,
  `internal/stats/`, `internal/report/`, `internal/changes/`,
  `internal/handoff/`, `internal/evidence/`.
- Persistence/integration adapters: `internal/journal/`, `internal/db/`,
  `internal/gitctx/`, `migrations/`.
- Shared utilities/design system: no shared design system in v1. Keep shared
  helpers inside the narrowest owning `internal/*` package until a broader
  boundary is recorded in `project-docs/specs/02_TRD.md`.

Rules:

- Business rules belong in the core/domain layer, not in interface adapters.
- Persistence and external-service details stay behind thin adapters.
- New boundaries must be recorded in `project-docs/specs/02_TRD.md`.

## Key Conventions

- Use red/green TDD for behavior changes: failing test first, then code.
- Tests must assert behavior, not implementation trivia.
- Do not delete or weaken tests to make a failing gate pass.
- Keep changes small and scoped to one task.
- A Developer session implements exactly one task.
- Use file artifacts as shared memory. Do not rely on chat history.
- Update impacted documentation in the same task as the code change.
- Normal and heavy planned phases use `VAL-NNN` validation assertions in
  `project-docs/docs/phases/phase-N/qa/validation-contract.md`; task specs and
  completion reports must cite covered assertions or state why none apply.
- Use Conventional Commits for commit subjects.
- Include `Refs: phase-N/task-M` in commit body or footer for phase work.
- Prefer existing local patterns over new abstractions.
- Add abstractions only when they remove real complexity or match an existing
  project pattern.
- Use `project-docs/specs/GLOSSARY.md` as the canonical project language
  source. Update it or flag ambiguity when terms are unclear.
- Record spec deviations in the completion report and phase decisions.
- Treat specifications as the primary engineering artifact. If a requirement
  cannot be satisfied, stop, report the gap, and do not stub, mock, or silently
  adapt the spec.
- If a required validation assertion cannot be satisfied, stop and escalate
  rather than weakening the assertion or passing it without evidence.
- Deterministic claims require evidence: command output, test logs, file state,
  screenshots, API traces, or equivalent reproducible proof.
- Keep held-out scenario material outside normal implementation context. Do not
  read hidden scenario repositories, excluded scenario folders, or evaluator
  prompts unless the task is explicitly scenario validation.

Project-specific conventions:

- Naming: follow Go package/function naming conventions; package names are
  short, lowercase, and domain-specific. fid entity IDs use
  `<PREFIX>-<hash>`; phase validation assertions use `VAL-NNN`; QA plan steps
  use `QA-NN`.
- Error handling: return errors rather than panicking for expected failures;
  wrap lower-level errors with context using `fmt.Errorf("operation: %w", err)`;
  present structured stderr output and JSON errors when `--format json` is
  active.
- Formatting style: Go code uses `gofmt` and `goimports`; Markdown stays
  concise and keeps command examples current with the CLI specs.
- File/function size budgets: non-test functions target 60 lines or less;
  non-test files target 400 lines or less, per
  `project-docs/specs/08_QUALITY_PLAN.md`.
- Test style: use red/green TDD for behavior changes; prefer table-driven
  tests; test observable behavior and contracts, not implementation trivia;
  use temporary SQLite databases for integration tests and golden fixtures for
  journal/event compatibility.
- Public interface compatibility policy: CLI flags, JSON output shapes, journal
  event schemas, migrations, and documented markdown/frontmatter contracts are
  compatibility surfaces. Breaking or versioned changes require a spec update
  and matching tests before implementation.

## Repo Skills

Canonical repo skills live under `project-docs/skills/`. Provider-discoverable
copies are generated from canonical skills into `.agents/skills/` for Codex and
`.claude/skills/` for Claude.

Rules:

- Edit canonical skills only; do not edit generated provider copies directly.
- Every skill must include a `Provenance:` line in the Markdown body.
- Skills are procedural helpers, not the source of truth for role authority,
  hard stops, required artifacts, validation gates, dependency policy, public
  interface policy, or architecture decisions.
- First versions should be instruction-only. Scripts require documented
  prerequisites, safe non-interactive behavior, and dependency-budget rationale
  when they add dependencies.
- Run `project-docs/tools/sync-skills.sh` after editing canonical skills.
- Run `project-docs/tools/workflow-lint.sh` after skill generation to check
  skill, prompt, template, and config drift.

## Dependency Budget

Do not add a new runtime dependency without a spec update, task spec, decision
entry, or completion-report note recording:

- Need: what problem it solves and why the standard library or an existing
  dependency is insufficient.
- Alternatives considered.
- License check.

Default license policy: prefer permissive licenses such as MIT, BSD,
Apache-2.0, or MPL-2.0. Avoid GPL, AGPL, SSPL, unknown, or custom licenses
unless the human explicitly approves and records the decision.

Test-only and dev-only dependencies follow the same rule with a lower bar:
record the rationale in the task completion report.

## Quality And Release Readiness

Use the maturity path in `project-docs/specs/08_QUALITY_PLAN.md`:

1. Baseline: specs, TDD, task reports, QA plan, deterministic QA evidence, and
   human checklist.
2. Automated gates: stack-appropriate format, lint, test, coverage or equivalent
   confidence gate, type checks where available, SAST, dependency audit, license
   check, and secret scan.
3. Held-out scenarios: user-journey scenarios stored outside normal agent
   context and run before release readiness; Mode B human runs are required for
   product judgment.
4. Release readiness: threat model, pre-release checklist, changelog, migration
   notes where applicable, backup/rollback path, and periodic independent audits.

Record any deliberate deferral in a task spec, phase decision, or roadmap item.

Each phase declares an operating mode (light, normal, or heavy) in the phase
brief. Modes vary required artifacts, never evidence rules or hard stops. See
`project-docs/specs/09_AI_WORKFLOW.md` section 4.

## Commands

Replace every placeholder before implementation starts. Agents must stop and
report a configuration defect if a required command is still a placeholder.
Use the task-spec command list when it exists; otherwise use these defaults.
Every required command listed by a task spec or QA plan must be recorded in the
completion report as `PASS`, or as `BLOCKED` with reason in a blocker report.
A missing, silently skipped, or narrower substituted required command is a
failed gate, not an acceptable omission.

Sandboxed Codex agents use `.codex/config.toml` to redirect Go and tool caches
into gitignored repo-local directories (`.gocache/`, `.gomodcache/`,
`.gopath/`, `.golangci-cache/`, `.cache/`). Do not request unsandboxed
execution just to write Go build, module, or lint caches. A fresh agent session
is required after changing that config. Network downloads and dependency
installs still require the normal dependency-budget and approval path.

- Format check: `gofmt -l cmd internal`. Apply `gofmt -w` and `goimports` to
  edited Go files before running the check. Do not use `gofmt -l .`; it walks
  ignored repo-local tool caches.
- Lint: `golangci-lint run ./...`.
- Default unit test suite: `go test -race ./...`.
- Coverage gate:

  ```bash
  go test -race -coverprofile=coverage.out ./...
  go tool cover -func=coverage.out
  ```

- Integration tests: `go test -tags integration -race ./...`.
- End-to-end tests: `N/A` until an end-to-end or held-out scenario harness is
  accepted in the specs.
- Build verification: `go build -o /tmp/fid ./cmd/fid`.
- Full phase-level command gate:

  ```bash
  gofmt -l cmd internal
  golangci-lint run ./...
  go test -race -coverprofile=coverage.out ./...
  go tool cover -func=coverage.out
  go test -tags integration -race ./...
  go test -bench . -benchmem ./...
  go build -o /tmp/fid ./cmd/fid
  ```

  Phase-level QA also performs the PRD compliance, security review,
  performance spot-check, documentation review, deterministic QA checklist, and
  human checklist steps in `project-docs/specs/08_QUALITY_PLAN.md` §5.2.

Release-readiness commands are stack-specific and may be `N/A` early in a
project, but must be concrete before a public or user-impacting release:

- Security scan: `N/A` before external release; a SAST or equivalent command
  must be selected and recorded before release readiness.
- Dependency vulnerability audit: `govulncheck ./...`.
- License compatibility check: `go-licenses check ./...`.
- Committed-secret scan: `gitleaks detect`.
- End-to-end or scenario validation: `N/A` until the held-out scenario harness
  location and command are recorded.

Recommended task-level gate:

```bash
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Recommended phase-level gate:

```bash
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

## Playbook

See `project-docs/docs/playbook.md` for a concise operator runbook.

## Specification Documents

- `project-docs/specs/00_PROJECT_CHARTER.md` - vision, scope, constraints.
- `project-docs/specs/GLOSSARY.md` - canonical project language.
- `project-docs/specs/01_PRD.md` - product requirements.
- `project-docs/specs/02_TRD.md` - technical requirements.
- `project-docs/specs/03_TRD_SCHEMA.md` - full SQL schema, triggers, validation.
- `project-docs/specs/04_TRD_CONFIG.md` - `.fid/` directory structure and config.
- `project-docs/specs/05_DESIGN_PLAN.md` - user experience and surfaces.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - CLI flag specs and output examples.
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` - end-to-end flows, whet mapping, export.
- `project-docs/specs/08_QUALITY_PLAN.md` - testing and verification.
- `project-docs/specs/09_AI_WORKFLOW.md` - agent workflow.
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` - phase plan.

## Agent Working Documents

- `project-docs/docs/current/` - default-loaded current context.
- `project-docs/docs/phases/phase-N/` - archived per-phase task, QA, progress,
  decisions, triage, and escalation artifacts.
- `project-docs/docs/phases/phase-N/qa/validation-contract.md` - phase
  `VAL-NNN` assertions that connect requirements to task evidence and final QA
  outcomes.
- `project-docs/docs/reviews/` - retrospectives, audits, and external reviews.
- `project-docs/templates/` - reusable artifact templates.
- `project-docs/skills/` - canonical skill authoring source.

## Hard Stops

Agents must stop and notify the human when:

- A deterministic quality gate fails and cannot be fixed inside the assigned
  task scope.
- A required project command is still a placeholder.
- A task requires an architectural, product, data-model, or public-interface
  decision not already recorded in specs.
- A task depends on ambiguous terminology that is not resolved or flagged in
  `project-docs/specs/GLOSSARY.md`.
- A dependency addition lacks documented budget rationale.
- A required completion report, QA plan, review, or human checklist is missing.
- A required phase validation contract is missing, has uncovered `VAL-NNN`
  assertions, or needs a product/architecture decision before evidence can be
  produced.
- A normal implementation task tries to read held-out scenarios or evaluator
  prompts that are meant to stay outside agent context.
- A release-readiness phase lacks the required threat model, pre-release
  checklist, or scenario/audit decision for its maturity level.
