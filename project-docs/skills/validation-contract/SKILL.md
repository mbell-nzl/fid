---
name: validation-contract
description: Use when drafting, refining, or auditing a phase validation contract with VAL-NNN assertions before implementation or during QA/wrap-up.
---

# Validation Contract

Provenance: Workflow-Extracted Skill.

Use this skill to turn phase scope into auditable `VAL-NNN` assertions without
making the skill the source of workflow authority.

## Procedure

1. Load the phase scope plus the relevant requirements from
   `project-docs/specs/`, the glossary, and
   `project-docs/templates/validation-contract-template.md`.
2. Express each in-scope requirement, decision, or invariant as one observable
   validation assertion. Keep each `VAL-NNN` narrow enough that evidence can
   prove or disprove it. If a requirement contains independent observables,
   split it into separate assertions; if human judgment is necessary, fan it
   out into sub-checks rather than bundling perceptions together.
3. For each assertion, fill the source requirement, evidence method, validator
   role, task/window ownership, pass signal, fail or blocked signal, and any
   explicit non-goals.
4. For conditional or binary assertions, name the qualifying and non-qualifying
   states that must be evidenced. Positive-only evidence is insufficient.
5. For parity or "matches reference" assertions, require structural
   conformance or deterministic comparison evidence where the state is
   machine-checkable; reserve `Human` for judgment that cannot be reduced to a
   reproducible check.
6. Check the planned task set: every task should cite covered `VAL-NNN` IDs or
   explain why no validation assertion applies.
7. Mark unresolved evidence, ownership, or terminology gaps where the canonical
   workflow expects them: task spec, QA plan, phase decisions, or escalation
   artifacts.
8. During QA or wrap-up, compare final evidence against the contract and keep
   outcome status tied to the original assertion language.

## References

- `project-docs/templates/validation-contract-template.md`
- `project-docs/templates/task-spec-template.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/09_AI_WORKFLOW.md`
- `project-docs/specs/GLOSSARY.md`
