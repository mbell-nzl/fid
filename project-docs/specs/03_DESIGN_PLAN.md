# DESIGN PLAN - <PROJECT_NAME>

> This document defines the intended user experience across the project's
> surfaces. Delete irrelevant sections and add stack-specific details after the
> project direction is known.

## 1. Experience Principles

1. **<Principle>.** <Implication for interaction, output, or layout.>
2. **<Principle>.** <Implication.>
3. **<Principle>.** <Implication.>

## 2. Surfaces

| Surface | Primary Jobs | Constraints | Verification |
|---|---|---|---|
| <Surface> | <Jobs> | <Accessibility, latency, density, device, etc.> | <How checked> |

## 3. Information Architecture

Describe navigation, object hierarchy, commands, routes, screens, API resource
groupings, or other surface-specific structure.

```text
<top-level>
|-- <area>
|-- <area>
`-- <area>
```

## 4. Interaction Patterns

| Pattern | Use When | Required Behavior |
|---|---|---|
| <Pattern> | <Context> | <Behavior> |

## 5. Output / UI Standards

Use this section for visual, textual, or structural standards that must remain
consistent across features.

- Terminology: <preferred words and banned synonyms>
- Empty states: <rules>
- Error states: <rules>
- Confirmation / destructive actions: <rules>
- Accessibility: <rules>
- Internationalization / localization: <rules or non-goals>

## 6. Mockups / Examples

Add sketches, text output examples, API examples, wireframes, or screenshots as
the project requires. Keep examples synchronized with implemented behavior.

### 6.1 <Surface / Flow>

```text
<Example or mockup>
```

## 7. Cross-Surface Parity Matrix

Use this when multiple surfaces must represent the same concept consistently.
The active phase brief should include a narrowed phase-specific parity matrix
whenever a polish or parity phase is in progress.

| Concept | Surface A | Surface B | Surface C | Required Consistency |
|---|---|---|---|---|
| <Concept> | <Representation> | <Representation> | <Representation> | <Rule> |

## 8. Human Verification Guidance

List the kinds of judgments that require human review.

- Visual hierarchy and density.
- Tone and terminology fit.
- End-to-end flow coherence.
- Whether the product solves the intended job in real use.
- Scenario-level friction that does not appear in unit or integration tests.

## 9. Scenario Candidates

List user journeys that should become held-out scenarios before release
readiness. Do not include the hidden scenario content here; record only public
scenario themes and risk areas.

| Theme | Risk It Covers | Target Release / Phase |
|---|---|---|
| <Theme> | <Risk> | <Phase> |

*Last updated: <YYYY-MM-DD>*

