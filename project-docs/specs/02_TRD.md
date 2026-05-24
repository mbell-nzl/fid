# TECHNICAL REQUIREMENTS DOCUMENT - <PROJECT_NAME>

> This document records the technical shape of the project: stack choices,
> architecture, data model, integration boundaries, and non-functional
> engineering requirements. Fill stack-specific sections after the project is
> forked.

## 1. Stack Decisions

| Area | Decision | Rationale | Status |
|---|---|---|---|
| Language/runtime | <TBD> | <Why> | Proposed / Accepted |
| Frameworks | <TBD> | <Why> | Proposed / Accepted |
| Storage | <TBD> | <Why> | Proposed / Accepted |
| Test tooling | <TBD> | <Why> | Proposed / Accepted |
| Build/tool runner | <TBD> | <Why> | Proposed / Accepted |
| Deployment | <TBD> | <Why> | Proposed / Accepted |

Do not start implementation tasks that depend on a stack decision still marked
`<TBD>` or `Proposed` unless the task is explicitly a discovery spike.

## 2. Architecture

Describe the intended layers and dependency direction.

```text
<interface adapters>
        |
        v
<application / core / domain services>
        |
        v
<repositories / integration adapters>
        |
        v
<external systems / storage / platform>
```

Rules:

- Dependencies point inward toward domain rules where possible.
- Interface adapters stay thin.
- External systems are hidden behind project-owned interfaces or modules.
- Shared libraries must not become a dumping ground for business behavior.

## 3. Module / Directory Boundaries

| Path / Module | Owner Layer | Allowed Responsibilities | Forbidden Responsibilities |
|---|---|---|---|
| `<path>` | <Layer> | <Allowed> | <Forbidden> |

## 4. Data Model

Fill this with the conceptual and physical model once storage is chosen.

### 4.1 Entities

| Entity | Fields / Attributes | Relationships | Notes |
|---|---|---|---|
| <Entity> | <Fields> | <Links> | <Notes> |

### 4.2 Persistence Rules

- <Transaction, consistency, migration, retention, or serialization rule>
- <Backup/export/import rule, if any>

## 5. Interfaces And Contracts

| Contract | Consumers | Compatibility Policy | Verification |
|---|---|---|---|
| <API, command, schema, event, file format, UI route, library API> | <Who/what uses it> | <Stable/experimental/internal> | <How tested> |

## 6. Security And Privacy

| ID | Requirement | Verification |
|---|---|---|
| TRD-SEC-001 | The system shall <security behavior>. | <Test/review/check> |

## 7. Performance And Reliability

| ID | Requirement | Target | Verification |
|---|---|---|---|
| TRD-PERF-001 | <Operation> | <Target> | <Benchmark/test/manual check> |

## 8. Observability And Operations

- Logging: <TBD>
- Metrics: <TBD>
- Error reporting: <TBD>
- Configuration: <TBD>
- Backup / recovery: <TBD or N/A>
- Deployment/release process: <TBD>
- Rollback process: <TBD or N/A>

Before external use, create or link a threat model and release checklist from
`project-docs/templates/`.

## 9. Dependency Policy

Dependencies must follow the budget in `AGENTS.md`. Record accepted runtime
dependencies here.

| Dependency | Purpose | Alternatives Considered | License | Decision Ref |
|---|---|---|---|---|
| <name> | <need> | <alternatives> | <license> | <DEC-NNN> |

Supply-chain gates to define before release readiness:

- Vulnerability audit: <command or compensating control>
- License check: <command or manual process>
- Secret scan: <command or manual process>
- Update cadence: <Dependabot/Renovate/manual cadence or N/A>

## 10. Technical Requirements Traceability

Use stable IDs:

```text
TRD-<AREA>-<NNN>
```

| TRD ID | Requirement | Phase | Verification |
|---|---|---|---|
| TRD-ARCH-001 | <Requirement> | <Phase> | <Test/review> |

## 11. Open Technical Questions

- [ ] <Question requiring Tech Lead or human decision>

*Last updated: <YYYY-MM-DD>*

