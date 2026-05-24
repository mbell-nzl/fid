# Phase <N> - Task <M>: <Title>

## Objective

<What this task changes and why.>

## Requirements

- PRD refs: <IDs or sections>
- TRD refs: <IDs or sections>
- Design refs: <sections or "N/A">
- Quality refs: <sections or "N/A">

## Validation Contract Coverage

List the `VAL-NNN` assertions this task advances. Assertions live in
`project-docs/docs/phases/phase-<N>/qa/validation-contract.md`.

- VAL-NNN: <one-line restatement of behavior or invariant>

If this task does not advance a validation assertion, state why explicitly
(for example, `documentation-only task; no behavior assertion`).

## Scope

### Files To Create Or Modify

- `<path>`

### Files To Read For Context

- `<path>`

### Out Of Scope

- <explicit non-goal>

## Acceptance Criteria

- [ ] <observable/testable criterion> (VAL-NNN or no VAL)
- [ ] <observable/testable criterion> (VAL-NNN or no VAL)

## TDD Plan

- First failing test:
- Additional tests:
- Regression risks:

## Documentation Updates

- `<doc path>` - <why>

## Scenario / Parity / Audit Notes

- Held-out scenario impact: <None / scenario theme only; do not paste hidden scenario content>
- Cross-surface parity impact: <None / Phase Brief Parity Matrix axes to check>
- Audit or release-readiness impact: <None / template or review to update>

## Quality Gate

Run exactly:

```text
<FORMAT_COMMAND>
<LINT_COMMAND>
<COVERAGE_COMMAND>
```

If any command is still a placeholder, stop and report a configuration defect.

## Notes For Developer

- <constraints, ordering, known pitfalls>
