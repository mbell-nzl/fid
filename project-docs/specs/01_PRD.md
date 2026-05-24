# PRODUCT REQUIREMENTS DOCUMENT - <PROJECT_NAME>

> This document defines product behavior in testable language. It should say
> what the system must do and why, not how it is implemented.

## 1. Requirement Style

Use EARS-style language where useful:

- Ubiquitous: "The system shall ..."
- Event-driven: "When <event>, the system shall ..."
- State-driven: "While <state>, the system shall ..."
- Optional feature: "Where <feature>, the system shall ..."
- Unwanted behavior: "If <error condition>, the system shall ..."

Use stable requirement IDs:

```text
PRD-<AREA>-<NNN>
```

Examples of area prefixes:

- `CORE` - core domain behavior
- `UX` - user-facing interaction
- `AUTH` - identity, access, permissions
- `DATA` - user data behavior and portability
- `INT` - integrations
- `OPS` - operational behavior

Replace these with project-specific prefixes as the product takes shape.

## 2. Personas / User Jobs

| Persona | Job To Be Done | Success Signal |
|---|---|---|
| <Persona> | <Task/outcome> | <Observable signal> |

## 3. Functional Requirements

### 3.1 <Area Name>

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-CORE-001 | The system shall <behavior>. | Must | <Notes> |
| PRD-CORE-002 | When <event>, the system shall <response>. | Should | <Notes> |

### 3.2 <Area Name>

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-UX-001 | The system shall <behavior>. | Must | <Notes> |

## 4. Non-Functional Product Requirements

These are user-visible qualities, not implementation choices.

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-NFR-001 | The system shall <reliability, accessibility, performance, privacy, or usability behavior>. | Must | <Notes> |

## 5. Error And Empty States

Define how the product behaves when data is missing, invalid, incomplete, or
ambiguous.

| ID | Condition | Required Behavior |
|---|---|---|
| PRD-ERR-001 | <Condition> | <Behavior> |

## 6. Traceability

Each roadmap item should reference PRD IDs. Each implemented PRD ID should have
tests or verification evidence.

| PRD ID | Roadmap Phase | Task(s) | Verification |
|---|---|---|---|
| PRD-CORE-001 | Phase <N> | task-<M> | <test or checklist> |

## 7. Open Product Questions

- [ ] <Question requiring human/product-owner judgment>

*Last updated: <YYYY-MM-DD>*

