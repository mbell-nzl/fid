# IMPLEMENTATION ROADMAP - <PROJECT_NAME>

> This roadmap sequences the work from first useful build to release. Keep it
> phase-oriented rather than sprint-oriented so agents can plan and execute one
> coherent phase at a time.

## Guiding Principles

1. **Dogfood early.** The first milestone should be useful to the project owner
   or a real target user.
2. **Cut scope, not quality.** Defer features rather than shipping fragile
   foundations.
3. **Stabilize contracts early.** Data models, public APIs, file formats, and
   workflow contracts are more expensive to change than UI details.
4. **Keep phases coherent.** Each phase should have a clear goal, exit
   criteria, and a small set of related tasks.
5. **Prove requirements explicitly.** Phase plans should turn scoped
   requirements into `VAL-NNN` assertions before implementation begins.

## Phase 0 - Specification And Foundation

**Goal:** Make the project implementable by agents: stack decisions, baseline
architecture, repository scaffold, quality commands, and the first thin slice.

| Item | Description | PRD/TRD Refs |
|---|---|---|
| Stack decisions | Choose language, framework, storage, test tools, build runner. | TRD-* |
| Repo scaffold | Create project skeleton, formatting, linting, test harness, CI stub. | TRD-* |
| Walking skeleton | Minimal executable or library path proving build/test/run flow. | PRD-* / TRD-* |
| Agent workflow validation | Confirm commands in `AGENTS.md` and `04_QUALITY_PLAN.md` are real, and create the first `qa/validation-contract.md` from the template. | WF-* |
| Quality maturity choice | Decide which quality layers apply before MVP and before release readiness. | QP-* |

**Exit criteria:**

- Project builds or packages through `<BUILD_COMMAND>`.
- Default tests run through `<TEST_COMMAND>`.
- Task-level gate commands are concrete and green.
- Phase planning produces task specs with `VAL-NNN` coverage and a validation
  contract with final-review ownership.
- Phase 1 can be planned without unresolved stack choices.

## Phase 1 - MVP Core

**Goal:** Deliver the smallest useful version of the core product behavior.

| Item | Description | PRD/TRD Refs |
|---|---|---|
| <Item> | <Description> | <Refs> |

**Exit criteria:**

- <Observable behavior>
- <Quality or docs criterion>

## Phase 2 - Primary Workflow

**Goal:** Make the primary user workflow complete enough for dogfooding.

| Item | Description | PRD/TRD Refs |
|---|---|---|
| <Item> | <Description> | <Refs> |

**Exit criteria:**

- <Observable behavior>
- <Quality or docs criterion>

## Phase 3 - Hardening And Release Readiness

**Goal:** Improve reliability, edge-case handling, documentation, and release
process.

| Item | Description | PRD/TRD Refs |
|---|---|---|
| <Item> | <Description> | <Refs> |
| Automated quality gates | Add stack-appropriate security, dependency, license, and secret scans. | QP-* / TRD-* |
| Release checks | Complete threat model, pre-release checklist, and scenario validation decision. | QP-* / TRD-* |

**Exit criteria:**

- <Release readiness condition>
- <Human verification condition>
- Threat model and pre-release checklist are current, or an explicit deferral is recorded.
- Held-out scenario plan is complete and protected from normal implementation context.

## Later Phases

Use this section for deferred features and expansion paths.

| Phase | Theme | Notes |
|---|---|---|
| vNext | <Theme> | <Notes> |

## Risk Register

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| <Risk> | <Low/Medium/High> | <Low/Medium/High> | <Mitigation> |

*Last updated: <YYYY-MM-DD>*
