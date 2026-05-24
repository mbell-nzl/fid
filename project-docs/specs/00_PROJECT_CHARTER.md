# PROJECT CHARTER - <PROJECT_NAME>

> This document records the durable intent of the project: why it exists, who
> it serves, what is in and out of scope, and which decisions should survive
> beyond individual implementation phases.

## 1. Vision

<One to three paragraphs describing the long-term product or system vision.>

## 2. Problem Statement

<Describe the user, business, technical, or operational problem. Include what
people do today, why it is insufficient, and what meaningful improvement would
look like.>

## 3. Target Users / Consumers

| Segment | Description | Primary Needs |
|---|---|---|
| Primary | <Who this is mainly for> | <Jobs to be done> |
| Secondary | <Adjacent users or operators> | <Needs> |
| Not target | <Explicit non-audience> | <Reason excluded> |

## 4. Design Principles

These principles guide tradeoffs when requirements conflict.

1. **<Principle>.** <Implication for product and engineering decisions.>
2. **<Principle>.** <Implication.>
3. **<Principle>.** <Implication.>

## 5. Scope

### 5.1 Core Concepts

Define domain concepts without committing to implementation details. Canonical
term definitions live in `project-docs/specs/GLOSSARY.md`; keep this section as
a brief concept inventory, not a duplicate glossary.

| Concept | Description |
|---|---|
| <Concept> | <Meaning in the product/domain> |

### 5.2 In Scope

- <Capability or constraint>
- <Capability or constraint>

### 5.3 Out of Scope

- <Explicitly deferred feature, integration, platform, or behavior>
- <Explicit non-goal>

## 6. Interfaces / Surfaces

List the user-facing or integration surfaces the project will expose.

| Surface | Purpose | Notes |
|---|---|---|
| <Surface> | <What it is for> | <Constraints or lifecycle> |

## 7. Success Criteria

Define objective or observable outcomes.

| Horizon | Criteria |
|---|---|
| Early dogfood | <What must be true before daily internal use> |
| Public release | <What must be true before external use> |
| Long-term | <Business, adoption, reliability, or impact target> |

## 8. Constraints

- **Time / staffing:** <Constraints>
- **Budget:** <Constraints>
- **Security / privacy / compliance:** <Constraints>
- **Operational:** <Constraints>
- **Legal / licensing:** <Constraints>

## 9. Durable Decisions Log

Keep only project-shaping decisions here. Phase-local decisions belong in
`project-docs/docs/phases/phase-N/decisions.md`.

| Decision | Rationale | Date / Reference |
|---|---|---|
| <Decision> | <Why> | <DEC-NNN or date> |

## 10. Document Suite

This charter is part of the core specification suite:

- `00_PROJECT_CHARTER.md` - vision, scope, constraints, durable decisions.
- `GLOSSARY.md` - canonical project language and flagged ambiguities.
- `01_PRD.md` - product requirements.
- `02_TRD.md` - technical requirements.
- `03_DESIGN_PLAN.md` - user experience and interface design.
- `04_QUALITY_PLAN.md` - testing strategy and quality gates.
- `05_AI_WORKFLOW.md` - agent roles and orchestration.
- `06_IMPLEMENTATION_ROADMAP.md` - implementation phases.

*Last updated: <YYYY-MM-DD>*
