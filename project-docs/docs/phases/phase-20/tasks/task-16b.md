# Phase 20 - Task 16b: report-lint parses declaration lines, not prose

## Objective

`report-lint.sh` treats **every** non-empty line under a task spec's
`## Expected Non-Zero Commands` as a declaration, so an instructional paragraph
or a wrapped explanation is rejected as invalid grammar. `task-16.md` produced
exactly six such failures; `task-17.md` and the task-spec template have the same
shape. Teach the parser which lines are declarations. Nothing about the
declaration rule itself changes.

## Validation Contract Coverage

- No VAL applies. This is workflow-control tooling, not product behavior. It
  unblocks the VAL-030/031/032/033/037 evidence Task 16 already produced.

## Binding Constraints

- DEC-144 decision 4 defines the rule; implement it literally. Any divergence
  escalates.
- **Only line selection changes.** The set of accepted *declarations* is
  identical: same grammar, `N` still greater than zero, declared command still
  required to appear in `## Exact Commands To Run`, same failure messages.
- A malformed line that *looks* like a declaration must still fail. Silently
  ignoring a typo'd declaration would let an undeclared non-zero `PASS` through
  a later check with a misleading message.
- Every other current check keeps its exact behavior and message: missing
  termination, non-zero without declaration, invalid status words, go-test token
  plausibility, command coverage, commit identity, committed/clean state.
- Red first: add each self-test case and observe it fail before implementing.

## Scope

### Files To Create Or Modify

- `project-docs/tools/report-lint.sh`
  (`expected_nonzero_declaration_lines`, header comment, and self-test).

### Files To Read For Context

- `../decisions.md` DEC-144 (decision 4 is the rule; decision 7 is why the
  fixture must match real shapes) and DEC-143 decision 5 (the unchanged
  declaration rule).
- `task-16.md` and `task-17.md` `## Expected Non-Zero Commands` - the two real
  sections in the repository, and the shapes the fixture must reproduce.
- `project-docs/templates/task-spec-template.md`
  (`## Expected Non-Zero Commands`).

### Out Of Scope

- Task 16 product code, `go.mod`/`go.sum`, `internal/format`, and any Go test.
- `task-16-done.md`, `task-16.md`, `task-17.md`, and the templates - the parser
  is wrong, not the documents. Do not reformat a real section to suit the
  script.
- Any other report-lint rule, message, or exit-code semantic.
- `workflow-lint.sh`, `workflow-sequence-lint.sh`, agent prompts, and skills.

## Acceptance Criteria

- [ ] Inside `## Expected Non-Zero Commands`, only lines beginning at column 0
      with ``- ` `` (hyphen, space, backtick) are parsed as declarations.
- [ ] Blank lines, prose paragraphs, indented continuation lines, and bullets
      that do not open with a backtick are ignored without a failure.
- [ ] A declaration line that does not match the DEC-143 grammar still fails
      with today's invalid-grammar message naming the line.
- [ ] A declared `N` of `0` and a declared command absent from
      `## Exact Commands To Run` still fail with today's distinct messages.
- [ ] Every accepted and rejected `PASS`-row case from `task-16a` behaves
      exactly as it does today.
- [ ] The self-test's fixture task spec reproduces all three real shapes: a
      leading instructional paragraph, a declaration whose trailing prose wraps
      onto an indented next line, and a prose bullet that does not open with a
      backtick.
- [ ] `--self-test` reports no failures and covers the shapes above plus the new
      malformed-declaration rejection.

## First Failing Test

`project-docs/tools/report-lint.sh --self-test` - first extend the fixture task
spec's `## Expected Non-Zero Commands` section with the instructional paragraph,
the wrapped continuation line, and the non-backtick prose bullet, then observe
the clean-report case fail with `expected non-zero declaration has invalid
grammar` before touching the parser.

## Implementation Notes

- The change belongs in `expected_nonzero_declaration_lines`: emit only lines
  matching `^- ` followed by a backtick. `validate_expected_nonzero_declarations`
  keeps its current body, so a selected-but-malformed line still reaches the
  existing failure.
- Awk pattern shape: `in_section && /^- `/ { print }`. Watch quoting of the
  backtick inside the single-quoted awk program.
- Suggested self-test additions, over the existing fixture:

| Fixture | Seeded shape | Expected result |
|---|---|---|
| acceptance | instructional paragraph, wrapped continuation, prose bullet in the fixture spec | clean report still passes, `OUTCOME COMPLETE` |
| malformed declaration | a backtick bullet whose expected field loses its backticks | fails with `invalid grammar` |

- The `zero-declaration` and `undeclared-in-gate` variant specs already exist;
  keep them working with the new selection rule.

## Documentation Updates

- `project-docs/tools/report-lint.sh` header comment: state that declarations
  are the backtick-bullet lines in that section and that other lines are prose.
- No other doc. `09_AI_WORKFLOW.md` §8 and the task-spec template already carry
  the declaration-line rule (DEC-144).

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
DEC-144 amended `09_AI_WORKFLOW.md` §8 and the task-spec template; if it reports
a failure caused by those doc edits, that is a Tech Lead defect - stop and report
it rather than editing a spec to satisfy the linter.

List any conditional command required by this task and why:

- None.

No `## Expected Non-Zero Commands` section: every required command above must
exit 0, so the section is omitted rather than filled with a placeholder.

## Notes For Developer

- Completion report: `task-16b-done.md`.
- Record one **DIAGNOSTIC** row in `## Tests run` for
  `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-16-done.md project-docs/docs/phases/phase-20/tasks/task-16.md`
  with its verbatim last line and observed termination. This is the real case
  that failed; the six declaration-grammar failures must be gone. It is a
  diagnostic, not a gate: `OUTCOME BLOCKED` is expected there because that report
  is still an unfinalized blocker report, and any residual non-grammar failure
  belongs to the Task 16 continuation, so record it and note it rather than
  fixing it here.
- Your own report has no expected-non-zero rows; it must pass under the amended
  script exactly as it would have under the current one.
- If the self-test cannot express a case without changing an unrelated check,
  stop and report it - do not weaken a neighbouring rule to make room.
