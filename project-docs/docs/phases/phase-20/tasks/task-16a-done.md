# Phase 20 - Task 16a: Accept declared expected-non-zero commands in report-lint - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added task-spec parsing and validation for optional declared expected-non-zero command rows.
- Allowed only the DEC-143 exact expected/observed non-zero `PASS` grammar while retaining the prior exit-zero rule for every other `PASS` row.
- Added one accepted fixture and six rejection fixtures to the report-lint self-test.

## Files created or modified

- `project-docs/tools/report-lint.sh`
- `project-docs/docs/phases/phase-20/tasks/task-16a-done.md`

## Documentation updated

- `project-docs/tools/report-lint.sh` - updated its header comment to describe declared expected-non-zero verification.

## Validation Contract Coverage

- No VAL applies: this task changes workflow-control tooling only and unblocks Task 16 evidence; it changes no product behavior.

## Tests run

- `project-docs/tools/report-lint.sh --self-test` - PASS - `SELF-TEST PASS report-lint self-test completed with no failures` (terminated `exit_code: 0`)

## Exact quality-gate results

- `sh -n project-docs/tools/report-lint.sh` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `project-docs/tools/report-lint.sh --self-test` - PASS - `SELF-TEST PASS report-lint self-test completed with no failures` (terminated `exit_code: 0`)
- `project-docs/tools/workflow-lint.sh` - PASS - `PASS workflow-lint completed with no failures` (terminated `exit_code: 0`)
- Recorded gate exception: no Go gate was run because Task 16a changes no Go source; Task 16's same-task continuation owns its full Go gate.

## Deviations from spec

- Red-first required self-test attempt: `SELF-TEST FAIL report-lint self-test found 4 failure(s)` (terminated `exit_code: 1`). Root cause: the existing validator only accepted a parenthesized zero exit termination and rejected the newly seeded truthful DEC-143 row. Repair: added declaration validation and the exact declared expected-non-zero row grammar; the unchanged command later passed with `SELF-TEST PASS report-lint self-test completed with no failures` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit

OUTCOME COMPLETE
