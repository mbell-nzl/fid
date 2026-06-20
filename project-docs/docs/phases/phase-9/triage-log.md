# Phase 9 Triage Log

## TRG-001 - Task 3 expanded write set spec deviation

Trigger: Task 3 report notes a spec deviation for the expanded write set (`internal/format/markdown.go` and a test-only stability adjustment).

Disposition: RESUME

Reasoning:
- The `internal/format/markdown.go` change is task-local because Task 3 explicitly required `fid concern list` coverage across `json`, `toon`, `markdown`, and `yaml`; extending markdown formatting to satisfy VAL-008 implements that required output surface rather than changing the public interface beyond the task acceptance target.
- The `internal/issues/transition_rules_test.go` adjustment is test-only and removes an existing concern-ID ordering assumption surfaced by the required integration gate. Based on the report context, it does not alter persistence, domain behavior, CLI contract, or phase testing policy.
- The report records the full gate as PASS and no remaining concerns, tech debt, or notes, so this does not meet the blocking/escalation triggers for public interface, data model, standing decision, cross-subsystem behavior, dependency, testing strategy, or unresolved product acceptance ambiguity.

References:
- project-docs/docs/phases/phase-9/tasks/task-3-done.md

## TRG-002 - Final review blocking findings for Task 3 and phase closeout

Trigger: final review has blocking findings for Task 3 and phase closeout:
`go test -bench . -benchmem ./...` fails in `cmd/fid` at
`TestConcernListFiltersByIssueAndStatus/all_concerns`; VAL-004 / QA-12 fails
because after `fid concern address`, public read-back shows `status="addressed"`
but does not expose `addressed_in_completion_id` on either `fid concern list` or
`fid issue get --section concerns`.

Disposition: BLOCK

Escalation: ESC-006

Reasoning:
- The benchmark failure is a required phase-level quality gate failure, not a
  misclassified finding. The final review records the command exit as `1`, and
  no later verified commit is recorded as resolving it.
- The VAL-004 / QA-12 failure is a real acceptance-contract failure for the
  concern address lifecycle. The phase brief requires address to record the
  completion, and the final review confirms the status mutates to `addressed`
  while both public read-back surfaces omit `addressed_in_completion_id`.
- No accepted recorded decision permits omitting the addressed completion from
  public concern read-back. DEC-D says address records an in-place
  `open`->`addressed` mutation; it does not waive the VAL-004 read-back
  requirement.
- Under the QA Blocking triage rule, RESUME is allowed only when the finding is
  misclassified, already resolved by a verified commit, or covered by an
  existing accepted decision. None of those conditions apply.

References:
- project-docs/docs/phases/phase-9/qa/final-review.md: Verdict, Findings,
  Phase-Level Gate Results, Validation Contract Outcomes, Deterministic
  Checklist Execution
- project-docs/docs/current/phase_brief.md: Validation contract index, Planning
  decisions DEC-D
- project-docs/docs/current/tl_escalations_current.md: ESC-006
