# Phase 20 - Task 16a: Accept declared expected-non-zero commands in report-lint

## Objective

`report-lint.sh` fails every truthful `PASS` row whose command's pass condition
is a non-zero exit, so no honest Task 16 report can pass it. Teach the control
the notation that `AGENTS.md` and `09_AI_WORKFLOW.md` §8 already require: a task
spec pre-declares an expected non-zero status, and the report records the real
observed status beside it. Nothing else about the control relaxes.

## Validation Contract Coverage

- No VAL applies. This is workflow-control tooling, not product behavior. It
  unblocks the VAL-030/031/032/033/037 evidence Task 16 already produced.

## Binding Constraints

- DEC-143 defines this rule; implement it literally. Any divergence escalates.
- The expectation is **declared by the task spec**, never invented by the
  report. A report may not self-authorize a non-zero `PASS`.
- The observed termination stays verbatim. This change makes a true statement
  recordable; it never lets a failed gate read as `PASS`.
- Every other current check keeps its exact behavior and message: missing
  termination, non-zero without declaration, invalid status words, go-test
  token plausibility, command coverage, commit identity, committed/clean state.
- Red first: add each self-test case and observe it fail before implementing.

## Scope

### Files To Create Or Modify

- `project-docs/tools/report-lint.sh` (parser, `validate_pass_termination`,
  header comment, and self-test).

### Files To Read For Context

- `../decisions.md` DEC-143 (decision 5 is the rule; decision 6 is the notation).
- `AGENTS.md` (Commands) and `project-docs/specs/09_AI_WORKFLOW.md` §8 for the
  canonical row grammar this control must accept.
- `project-docs/templates/task-spec-template.md`
  (`## Expected Non-Zero Commands`) and
  `project-docs/templates/completion-report-template.md` (example row).

### Out Of Scope

- Task 16 product code, `go.mod`/`go.sum`, `internal/format`, and any Go test.
- Any other report-lint rule, message, or exit-code semantic.
- `workflow-lint.sh`, `workflow-sequence-lint.sh`, agent prompts, and skills.

## Acceptance Criteria

- [ ] The task-spec parser reads an optional `## Expected Non-Zero Commands`
      section whose rows read ``- `<command>` - expected `exit_code: <N>` ``,
      with trailing prose after the expected field allowed and ignored.
- [ ] A declared `N` of `0` fails with a distinct message naming the command.
- [ ] A declared command that is absent from the spec's
      `## Exact Commands To Run` block fails with a distinct message.
- [ ] A `PASS` row ending exactly in
      ``(expected `exit_code: N`; terminated `exit_code: M`)`` passes when the
      command is declared with the same `N` and `M` equals `N`.
- [ ] The same row fails when `M` differs from `N`, when the command is not
      declared, when the reported `N` differs from the declared value, and when
      any prose follows the closing parenthesis.
- [ ] A `PASS` row with a plain non-zero termination and no declaration still
      fails with today's message, and a `PASS` row with
      ``(terminated `exit_code: 0`)`` still passes unchanged.
- [ ] The verbatim-token check still applies to the new form: the row's token
      field is still the fourth backtick field and may not be empty.
- [ ] `--self-test` covers one acceptance and six rejections and reports no
      failures.

## First Failing Test

`project-docs/tools/report-lint.sh --self-test` - add the acceptance case first
(a declared `rg -n absence-fixture cmd` row in the clean fixture report and its
declaration in the fixture task spec) and observe the clean-report case fail
before touching the validator.

## Implementation Notes

Rule, per DEC-143 decision 5, for each `- \`<cmd>\` - PASS - ` line:

1. Ends in ``(terminated `exit_code: 0`)`` - unchanged behavior.
2. Ends in ``(expected `exit_code: N`; terminated `exit_code: M`)`` - accept
   only when `N > 0`, `M == N`, and the spec declares that exact command with
   that exact `N`.
3. Has a termination marker but neither shape - today's
   `did not terminate successfully with exit_code 0` failure.
4. Has no termination marker - today's `missing observed termination` failure.

Both end-of-line matches stay anchored (trailing whitespace only), because
prose after the marker already broke two Task 16 report attempts. Explanatory
notes belong on the following indented line, which the parser ignores.

Declarations come from the spec file, never from the report. Suggested
self-test fixture edits, all `sed` one-liners over the existing clean fixture:

| Fixture | Seeded defect | Expected message fragment |
|---|---|---|
| acceptance | declared row, `N` = `M` = 1 | clean report still passes |
| mismatched observation | terminated `exit_code: 2` | observed termination does not match |
| undeclared command | row command renamed | not declared in the task spec |
| reported/declared mismatch | expected `exit_code: 2` in the report | does not match the declared |
| trailing prose | text appended after `)` | did not terminate successfully |
| zero declaration | spec declares `exit_code: 0` | expected exit_code must be non-zero |
| undeclared-in-gate | spec declares a command absent from the command block | not a required command |

`expect()` currently hardcodes `$tmp/task-spec.md`; give it an optional spec
argument so the last two cases can use a variant spec fixture.

## Documentation Updates

- `project-docs/tools/report-lint.sh` header comment: add the declared
  expected-non-zero rule to the numbered check list.
- No other doc. `AGENTS.md`, `09_AI_WORKFLOW.md` §8, and both templates already
  carry the canonical notation (DEC-143 decision 6).

## Exact Commands To Run

Run exactly:

```text
sh -n project-docs/tools/report-lint.sh
project-docs/tools/report-lint.sh --self-test
project-docs/tools/workflow-lint.sh
```

**Recorded gate exception.** This task changes no Go source, so the standard Go
gate is not required here; the Task 16 continuation reruns the full Go gate in
its own session immediately afterwards. `workflow-lint.sh` is required because
DEC-143 amended `AGENTS.md`, `09_AI_WORKFLOW.md` §8, and both templates; if it
reports a failure caused by those doc edits, that is a Tech Lead defect - stop
and report it rather than editing a spec to satisfy the linter.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Completion report: `task-16a-done.md`.
- Your own report has no expected-non-zero rows; it must pass under the amended
  script exactly as it would have under the current one.
- If the self-test cannot express a case without changing an unrelated check,
  stop and report it - do not weaken a neighbouring rule to make room.
