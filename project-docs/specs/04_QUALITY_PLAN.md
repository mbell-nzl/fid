# QUALITY PLAN - <PROJECT_NAME>

> This document defines testing strategy, quality gates, code standards,
> independent verification, and release-readiness practices for agentic
> implementation work. Replace all command placeholders after the stack is
> chosen.

## 1. Quality Philosophy

1. **Trust behavior, verify structure.** Agents can make their own tests pass;
   quality requires deterministic execution evidence, static analysis, and
   periodic review from outside the implementation session.
2. **Tests are executable specifications.** Behavior changes start with a
   failing test unless the task is explicitly documentation-only or exploratory.
3. **Deterministic checks first.** Formatting, linting, tests, coverage, build,
   contract, security, dependency, license, and secret checks run before
   subjective review when they apply.
4. **Agents cannot weaken tests.** Existing tests may be refactored only when
   the task requires it and behavior-preserving intent is clear.
5. **Human judgment is reserved for intent.** Agents execute deterministic
   checks; humans validate product fit, perception, and dogfooding.
6. **Specs are quality infrastructure.** Ambiguous requirements are defects.
   Fix the spec or escalate before implementing around the gap.
7. **Validation contracts make proof explicit.** Phase-level `VAL-NNN`
   assertions connect requirements to task evidence, QA review, and final
   sign-off.
8. **Independent checks compound.** Audit prompts, held-out scenarios, and CI
   gates are reusable controls that should improve across projects.

## 2. Quality Layer Framework

Use these layers as the maturity path. Do not skip a layer silently; record
intentional deferrals in the roadmap or phase decisions.

| Layer | What It Catches | Default Owner | When Required |
|---|---|---|---|
| 1. Spec precision | Ambiguity, missing acceptance criteria, scope drift | Human + Tech Lead | Before implementation |
| 2. Validation contract | Unproven requirements, vague acceptance, missing evidence ownership | Tech Lead + QA Lead | Every planned phase |
| 3. Red/green TDD | Functional regressions and edge cases | Developer | Every behavior task |
| 4. Agent QA review | Spec compliance, test quality, architecture smells | QA Lead / QA Reviewer | Every phase |
| 5. Automated static gates | Mechanical defects, vulnerable patterns, dependency and license risk, leaked secrets | CI / QA Lead | Before external use; earlier when cheap |
| 6. Held-out scenario validation | Spec gaming and end-to-end product misses | Human / evaluator outside implementation context | Before release readiness |
| 7. Independent audits | Structural drift, security gaps, complexity, stale specs | Fresh agent/model + human triage | Periodic or pre-release |
| 8. Manual behavioral validation | Perception, product judgment, domain fit | Human | Every phase sign-off |

## 3. Test Strategy

Replace locations and tooling with stack-specific choices.

| Type | Scope | Location | Runner | When Run |
|---|---|---|---|---|
| Unit | Single function/module/component behavior | `<TBD>` | `<TEST_COMMAND>` | Every task |
| Integration | Multiple components or real adapters | `<TBD>` | `<INTEGRATION_COMMAND>` | Phase end and targeted tasks |
| Contract | Public interfaces, schemas, commands, events, API shapes | `<TBD>` | `<TBD>` | Any interface task |
| End-to-end | Critical user journeys | `<TBD>` | `<E2E_COMMAND>` | Relevant phases |
| Accessibility / usability | Human or automated surface checks | `<TBD>` | `<TBD>` | Surface phases |
| Scenario | Held-out user journeys | Outside normal agent context | `<E2E_OR_SCENARIO_COMMAND>` or manual Mode B | Release readiness |

## 4. Test Conventions

Fill these after stack selection.

- Test naming: <TBD>
- Fixture strategy: <TBD>
- Test data isolation: <TBD>
- Time/randomness control: <TBD>
- External services: <mock/fake/local/test account policy>
- Snapshot/golden files: <allowed or not, update process>

Default rules:

- Tests must be deterministic and isolated.
- Tests must not depend on network access unless explicitly marked and gated.
- Time-dependent logic must be injectable or controllable in tests.
- Public interface behavior requires regression coverage.
- Tests should assert behavior and contracts, not implementation trivia.
- Do not modify a test to match a broken implementation; change the code or
  record a spec decision.

## 5. Code Quality Standards

### 5.1 Formatting

All code and generated artifacts must pass:

```text
<FORMAT_COMMAND>
```

### 5.2 Linting / Static Analysis

All code must pass:

```text
<LINT_COMMAND>
```

Before external release, define stack-appropriate checks for:

- Type checking where the stack supports it.
- SAST or language-specific security analysis.
- Dependency vulnerability audit.
- License compatibility check.
- Secret scanning.

### 5.3 Coverage / Confidence Gate

The project shall maintain the following minimum confidence gates:

| Area | Gate | Target |
|---|---|---|
| Core/domain behavior | <coverage, mutation, scenario, or review gate> | <TBD> |
| Interface adapters | <gate> | <TBD> |
| Critical workflows | <gate> | <TBD> |

Default command:

```text
<COVERAGE_COMMAND>
```

### 5.4 Code Structure

Fill in stack-specific limits. Defaults until replaced:

- Keep functions/components/modules small enough to review in one pass.
- Split files when unrelated responsibilities accumulate.
- Keep business behavior out of adapters.
- Prefer explicit error handling and useful context.
- Avoid hidden global state unless explicitly justified.

## 6. Quality Gates

### 6.1 Task-Level Gate

Run by the Developer before writing the completion report and commit.

```text
<FORMAT_COMMAND>
<LINT_COMMAND>
<COVERAGE_COMMAND>
```

All commands must pass. If a command is not applicable, replace it with a real
project command or record a Tech Lead decision before execution begins.

### 6.2 Phase-Level Gate

Run by the QA Lead at phase end.

```text
<CHECK_COMMAND>
```

The phase-level review also verifies:

1. Acceptance criteria from the phase QA plan.
2. Final outcome for every `VAL-NNN` in the phase validation contract.
3. PRD/TRD requirement compliance.
4. Test quality and meaningful coverage.
5. Documentation accuracy.
6. Security and privacy risks.
7. Performance or reliability targets touched by the phase.
8. Deterministic checklist execution with evidence.
9. Human checklist generation for perception and dogfooding.

### 6.3 Validation Contract Checks

Each phase uses
`project-docs/docs/phases/phase-N/qa/validation-contract.md`, created from
`project-docs/templates/validation-contract-template.md`.

Required checks:

- Tech Lead drafts `VAL-NNN` assertions before or alongside task specs.
- QA Lead refines evidence methods, validator roles, and task windows before
  implementation starts.
- Developer task specs and completion reports cite covered `VAL-NNN` IDs or
  explicitly state why none apply.
- QA Reviewer verifies relevant `VAL-NNN` IDs for each reviewed task window.
- QA Lead records final `PASS`, `FAIL`, `BLOCKED`, or `CARRY-FORWARD` outcome
  for every `VAL-NNN`.

If a required validation assertion cannot be satisfied, the responsible agent
stops and escalates. Silent spec adaptation or evidence-free `PASS` is not
allowed.

### 6.4 Release-Readiness Gate

Before a public, customer, or user-impacting release, define and run:

```text
<SECURITY_SCAN_COMMAND>
<DEPENDENCY_AUDIT_COMMAND>
<LICENSE_CHECK_COMMAND>
<SECRET_SCAN_COMMAND>
<E2E_OR_SCENARIO_COMMAND>
```

If any command is intentionally deferred, record why and what compensating
control applies.

### 6.5 Deterministic And Human Checklists

QA plans must split verification by type of judgment:

- Deterministic items end with `(qa-lead)`. QA Lead executes these and records
  literal command output, file-state evidence, screenshots, logs, API traces,
  or other reproducible proof in `final-review.md`.
- Human items end with `(human)`. The human executes these when judgment,
  perception, product intent, or dogfooding is required. Human items that
  verify a validation assertion cite the relevant `VAL-NNN`.

PASS verdicts on deterministic items require evidence. If evidence cannot be
shown, mark the item FAIL or BLOCKED. QA Lead must not PASS a perceptual or
product-owner judgment by reading code.

### 6.6 Command Syntax Validation

Before a QA plan is finalized, QA Lead validates every command referenced in
either checklist using help output, dry-run mode, documentation, or equivalent
proof. Invalid command syntax is a QA-plan defect.

### 6.7 Required QA Artifacts

Each phase stores QA artifacts under `project-docs/docs/phases/phase-N/qa/`.
Use the templates in `project-docs/templates/`:

- `validation-contract-template.md` -> `validation-contract.md`
- `qa-plan-template.md` -> `qa-plan.md`
- `mid-phase-review-template.md` -> `review-M.md`
- `final-review-template.md` -> `final-review.md`
- `human-checklist-template.md` -> `human-checklist.md`

QA artifacts must cite evidence, not prose-only confidence. Missing required
artifacts are hard stops for the phase runner or phase wrap-up.

### 6.8 Workflow-Lint And Skill Drift

Workflow-lint is a stack-neutral repository check for agent workflow drift. It
should run without project runtime dependencies and should be safe to invoke
from hooks or manual quality gates.

When repo skills are present, workflow-lint should verify:

- Canonical skills exist under `project-docs/skills/`.
- Generated Codex copies exist under `.agents/skills/`.
- Generated Claude copies exist under `.claude/skills/`.
- Generated copies match canonical skill sources, except for deliberate
  provider-specific metadata.
- Every skill has valid Agent Skills frontmatter and a `Provenance:` line in
  the Markdown body.
- Side-effectful skills disable implicit invocation where the provider supports
  that control.

Prompt slimming must not depend on skills until workflow-lint can prove the
generated copies are present and current.

## 7. Held-Out Scenario Validation

Held-out scenarios are user-journey checks that implementation agents do not
see. They are designed to catch shared blind spots between generated code and
generated tests.

Rules:

- Store scenarios outside normal implementation context, ideally in a private
  repository or an ignored/excluded path that agents do not load.
- Treat a scenario as contaminated if an implementation agent reads it; retire
  or rewrite it.
- Scenarios describe persona, goal, starting state, trajectory questions,
  satisfaction criteria, and forbidden behavior.
- Mode A: an evaluator agent drives the built product and grades the trace.
- Mode B: the human plays the persona and logs friction. Run at least one Mode B
  scenario before a release-readiness sign-off.
- Scenario failures become normal roadmap/task work; do not patch around them
  inside the evaluator run.

Use `project-docs/templates/held-out-scenario-template.yaml` as the starting
point.

## 8. Independent Audits

Use a fresh agent session, and preferably a different model family from the
implementation loop, for high-stakes audits. Save outputs under
`project-docs/docs/reviews/`.

Recommended audit set:

| Audit | Cadence | Template |
|---|---|---|
| Architecture | Every major phase block or release train | `templates/audits/audit-prompts.md` |
| Security | Before external release; periodically for sensitive data | `templates/audits/audit-prompts.md` |
| Maintainability | When change velocity slows or before handoff | `templates/audits/audit-prompts.md` |
| Performance | Before scale-sensitive release or after hot-path changes | `templates/audits/audit-prompts.md` |
| Complexity | When files/modules grow or refactors accumulate | `templates/audits/audit-prompts.md` |
| Spec drift | End of major phase blocks | `templates/audits/audit-prompts.md` |
| Dependency / supply chain | Quarterly or before release | `templates/audits/audit-prompts.md` |

## 9. Release Readiness

Before external use, the release candidate should have:

- [ ] All deterministic gates green on the release branch.
- [ ] Threat model written or refreshed.
- [ ] Pre-release checklist completed.
- [ ] Changelog or release notes updated.
- [ ] Migration path documented if data model or public contracts changed.
- [ ] No new dependency without recorded budget rationale and license check.
- [ ] At least one independent audit current for the risk profile.
- [ ] Held-out scenario run completed; at least one Mode B run since the last
  release-readiness sign-off.
- [ ] Rollback, backup, or recovery path verified where applicable.

Use `project-docs/templates/threat-model-template.md` and
`project-docs/templates/pre-release-checklist.md`.

## 10. Definition Of Done

A task is done when:

- [ ] Failing tests were written first or an explicit exception was recorded.
- [ ] Implementation satisfies task acceptance criteria.
- [ ] Covered `VAL-NNN` assertions are cited with evidence, or the task records
      why no validation assertion applies.
- [ ] Impacted documentation was updated.
- [ ] Task-level quality gate passed.
- [ ] Dependency-budget note is present.
- [ ] Completion report was written.
- [ ] Conventional Commit created with phase/task reference.

A phase is done when:

- [ ] All planned tasks are complete or explicitly deferred/cancelled.
- [ ] Phase-level quality gate passed.
- [ ] Every `VAL-NNN` has a final evidenced outcome or recorded carry-forward.
- [ ] QA Lead final review has no unresolved blocking findings.
- [ ] Human checklist was completed and accepted or follow-up work was planned.
- [ ] Specs were synchronized where needed.
- [ ] Phase decisions and progress were archived.
- [ ] Current docs were compressed.

## 11. CI / Automation

Fill in once the project has automation.

| Event | Required Checks |
|---|---|
| Push | `<TBD>` |
| Pull request | `<TBD>` |
| Release | `<TBD>` |

At minimum before external use, CI should block on format, lint, tests,
coverage or equivalent confidence gate, security scan, dependency audit, license
check, and secret scan. Where the stack lacks a command, record the chosen
compensating control.

*Last updated: 2026-05-23*
