# fid - Project Context

## What Is This Project?

fid is a local-first, agent-first issue tracker for AI-assisted software
development. It provides a structured, queryable knowledge graph over the
artifacts an agentic development workflow produces. Rows hold metadata and
identity; markdown files hold prose. SQLite-backed CLI.

## Tech Stack

Fill this before implementation starts.

- Primary language/runtime: Go
- Application framework(s): <TBD>
- Data/storage: SQLite (journal-canonical, DB derived — see `docs/adr/0001`)
- Interface surfaces: CLI
- Test tools: <TBD>
- Build/tool runner: <TBD>
- Deployment/runtime environment: local-first, single-user

If any item remains `<TBD>`, agents may still help refine specifications, but
they must not begin implementation that depends on that choice.

## Architecture

Describe the intended module boundaries in project-specific terms.

- User-facing adapters: <paths or modules>
- Core/domain/business logic: <paths or modules>
- Persistence/integration adapters: <paths or modules>
- Shared utilities/design system: <paths or modules>

Rules:

- Business rules belong in the core/domain layer, not in interface adapters.
- Persistence and external-service details stay behind thin adapters.
- New boundaries must be recorded in `project-docs/specs/02_TRD.md`.

## Key Conventions

- Use red/green TDD for behavior changes: failing test first, then code.
- Tests must assert behavior, not implementation trivia.
- Do not delete or weaken tests to make a failing gate pass.
- Keep changes small and scoped to one task.
- Update impacted documentation in the same task as the code change.
- Each planned phase uses `VAL-NNN` validation assertions in
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

Project-specific conventions to fill:

- Naming: <TBD>
- Error handling: <TBD>
- Formatting style: <TBD>
- File/function size budgets: <TBD or "none">
- Test style: <TBD>
- Public interface compatibility policy: <TBD>

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

Use the maturity path in `project-docs/specs/04_QUALITY_PLAN.md`:

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

## Commands

Replace every placeholder before implementation starts. Agents must stop and
report a configuration defect if a required command is still a placeholder.

- `<FORMAT_COMMAND>` - format code and docs where applicable.
- `<LINT_COMMAND>` - run static analysis or linting.
- `<TEST_COMMAND>` - run the default unit test suite.
- `<COVERAGE_COMMAND>` - run tests with coverage or equivalent confidence gate.
- `<INTEGRATION_COMMAND>` - run integration tests, if applicable.
- `<E2E_COMMAND>` - run end-to-end or browser tests, if applicable.
- `<BUILD_COMMAND>` - build/package the project.
- `<CHECK_COMMAND>` - run the full phase-level quality gate.

Release-readiness commands are stack-specific and may be `N/A` early in a
project, but must be concrete before a public or user-impacting release:

- `<SECURITY_SCAN_COMMAND>` - SAST or equivalent security scan.
- `<DEPENDENCY_AUDIT_COMMAND>` - dependency vulnerability audit.
- `<LICENSE_CHECK_COMMAND>` - license compatibility check.
- `<SECRET_SCAN_COMMAND>` - committed-secret scan.
- `<E2E_OR_SCENARIO_COMMAND>` - end-to-end or scenario validation command.

Recommended task-level gate:

```text
<FORMAT_COMMAND>
<LINT_COMMAND>
<COVERAGE_COMMAND>
```

Recommended phase-level gate:

```text
<CHECK_COMMAND>
```

## Playbook

See `project-docs/docs/playbook.md` for a concise operator runbook.

## Specification Documents

- `project-docs/specs/00_PROJECT_CHARTER.md` - vision, scope, constraints.
- `project-docs/specs/GLOSSARY.md` - canonical project language.
- `project-docs/specs/01_PRD.md` - product requirements.
- `project-docs/specs/02_TRD.md` - technical requirements.
- `project-docs/specs/03_DESIGN_PLAN.md` - user experience and surfaces.
- `project-docs/specs/04_QUALITY_PLAN.md` - testing and verification.
- `project-docs/specs/05_AI_WORKFLOW.md` - agent workflow.
- `project-docs/specs/06_IMPLEMENTATION_ROADMAP.md` - phase plan.

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
