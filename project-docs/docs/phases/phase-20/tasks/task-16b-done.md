# Phase 20 - Task 16b: report-lint parses declaration lines, not prose - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Restricted expected-non-zero parsing to column-0 backtick bullets, leaving explanatory section prose unparsed.
- Extended the report-lint self-test with the real instructional, wrapped-continuation, and prose-bullet shapes, plus malformed-declaration rejection.

## Files created or modified

- `project-docs/tools/report-lint.sh`
- `project-docs/docs/phases/phase-20/tasks/task-16b-done.md`

## Documentation updated

- `project-docs/tools/report-lint.sh` header and declaration-selector comment.

## Validation Contract Coverage

- No VAL applies. This task repairs workflow-control tooling and unblocks Task 16 evidence; it does not change product behavior.

## Tests run

- `project-docs/tools/report-lint.sh --self-test` - FAIL - `SELF-TEST FAIL report-lint self-test found 6 failure(s)` (terminated `exit_code: 1`) — red attempt after seeding the real declaration-section shapes.
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-16-done.md project-docs/docs/phases/phase-20/tasks/task-16.md` - DIAGNOSTIC - `PASS report-lint completed with no failures` (terminated `exit_code: 0`) — observed `OUTCOME BLOCKED` as expected for Task 16's unfinalized blocker report; its six former declaration-grammar failures were absent.

## Exact quality-gate results

- `sh -n project-docs/tools/report-lint.sh` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `project-docs/tools/report-lint.sh --self-test` - PASS - `SELF-TEST PASS report-lint self-test completed with no failures` (terminated `exit_code: 0`)
- `project-docs/tools/workflow-lint.sh` - PASS - `PASS workflow-lint completed with no failures` (terminated `exit_code: 0`)

## Deviations from spec

- Required-command attempt retained: `project-docs/tools/report-lint.sh --self-test` first failed with `SELF-TEST FAIL report-lint self-test found 6 failure(s)` (terminated `exit_code: 1`). Root cause: before the selector repair, every non-empty expected-non-zero-section line was parsed as a declaration. Repair: select only column-0 backtick bullets, preserving the existing declaration validator. The byte-identical rerun passed with `SELF-TEST PASS report-lint self-test completed with no failures` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `d8388d5c605750694004d446ffdc32f676e1521f`

## Commit subject

- `fix(workflow): parse expected-nonzero declarations`
