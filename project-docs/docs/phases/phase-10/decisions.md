# Phase 10 Decisions

> In-phase decisions (DEC-NNN / DEC-A..E). Per-doc spec-sync verdicts and the
> full decision history are finalized at phase wrap. Newest at top.

## Task 3 execution note (2026-06-21)

- **DEC-E executed as planned for permission retrofit fixtures.** Existing tests
  that exercised now-denied write combinations were moved to permitted roles
  without changing the behaviors under test:
  - `maintenance.rebuild`: `developer` -> `workflow_runner` in
    `cmd/fid/e2e_v01b_test.go`, `cmd/fid/e2e_v01c_ac_test.go`,
    `cmd/fid/e2e_v01c_complete_test.go`, `cmd/fid/e2e_v01c_concern_test.go`,
    `cmd/fid/init_test.go`, `cmd/fid/issue_ac_test.go`, `cmd/fid/issue_test.go`
    (all rebuild parity callers), and `cmd/fid/rebuild_test.go`.
  - Raw rebuild journal fixtures gained a seeded `workflow_runner` role so the
    permitted rebuild actor exists in the replay-only workspaces:
    `testdata/journal/replay-v1.golden.jsonl`,
    `testdata/journal/corrupt-tail.jsonl`.
  - `issue.transition.done` invalid-edge fixtures: `developer` -> `tl` in
    `cmd/fid/issue_transition_test.go`, `internal/issues/transition_test.go`,
    the synthetic precondition-conflict event in `cmd/fid/issue_test.go`, and
    `testdata/journal/conflicting-precondition.jsonl`.
  - `issue.transition.qa_review` setup fixture: `developer` -> `qa_reviewer` in
    `internal/issues/transition_test.go`.
  - Epic `issue.transition.in_progress` setup fixture: `tl` ->
    `workflow_runner` in `internal/issues/transition_test.go`.
  - `issue.transition.cancelled` task/subtask fixtures: `developer` -> `tl` in
    `internal/issues/transition_test.go`.
- **Direct-prerequisite matrix rows were added to finish enforcement coverage for
  already-implemented derived actions.** Phase 10 Task 2's initial matrix missed
  the epic transition states `issue.transition.wrapped` and
  `issue.transition.closed`, plus the denied AC-verifier roles
  `issue.ac.verify.tl_lite` and `issue.ac.verify.workflow_runner`. Task 3 added
  those rows so every existing derived write action resolves to a matrix cell and
  denied writes surface `PermissionDeniedError` instead of a missing-entry error.

## Task 1 execution note (2026-06-21)

- **DEC-D executed as planned for CF-003.** `internal/policy` role discovery no
  longer tolerates malformed or un-upgradeable trailing journal lines after
  seeded roles are found; `KnownRoles` now reports them with the same journal
  file/line context the replay path uses, so the upcoming permission matrix
  builds on a sound role set instead of a partially scanned one.
- **The TRG-002 test rewrite is intentional and decision-backed.**
  `TestKnownRolesTRG002CorruptTailTolerance` now asserts reported corruption for
  a valid role followed by malformed/un-upgradeable tail input. This closes the
  acknowledged Phase 1 debt rather than weakening coverage.

## Phase 10 planning decisions (DEC-A..E, human-confirmed 2026-06-21)

Recorded in the validation contract §Phase scoping notes and the phase brief.
Summarized here for the wrap audit; each task spec and QA artifact cites them by
DEC-A..E. The human accepted all five pre-execution.

- **DEC-A — enforcement seam in `internal/policy`, action key derived in the
  domain layer.** `policy.Authorize(workspaceRoot, role, action string) error`
  subsumes `RequireKnownRole` (preserving `MissingRoleError`/`UnknownRoleError`)
  and adds the in-code matrix lookup, returning `PermissionDeniedError{Role,
  Action}` (names role + action + "not permitted", PRD-ROLE-003). Action keys are
  derived in `internal/issues` / `cmd/fid/rebuild.go` from known
  kind/target-state/role. No new architecture boundary.
- **DEC-B — the in-code matrix is the single source of truth, shared by
  `Authorize` and `fid policy show`, and lists only implemented write actions.**
  Phase 10 seeds the v0.1a–c actions; Phases 11–13 append theirs alongside their
  commands. TRD-ARCH-010 completeness is tested against implemented actions.
- **DEC-C — `fid policy show` is a read command:** no `--by`, no permission
  check, no journal write, no `activity_log`; renders json/toon/markdown/yaml.
- **DEC-D — CF-003 reconciliation aligns role discovery with the journal
  corruption contract.** `KnownRoles`/`collectRoles` report a malformed/
  un-upgradeable journal line as corruption (non-zero exit, error naming the
  file) instead of tolerating it after seeded roles are found, matching the
  replayer (TRD-ARCH-002) and ADR-0001. The
  `TestKnownRolesTRG002CorruptTailTolerance` "tolerated" subtest is rewritten
  (decision-backed red/green), not weakened. Closes CF-003.
- **DEC-E — the retrofit moves existing test fixtures to permitted role/action
  combinations.** Enforcing the matrix changes existing command behaviour (e.g.
  `developer` may not create an epic; `qa_reviewer` may not transition to
  `dev_complete`). Fixtures using now-denied combinations move to permitted roles
  (`tl`/`human`/etc.), recorded as a deliberate, spec-mandated behaviour change,
  not test weakening. Every fixture role change is listed in the Task 3
  completion report. Validate-before-append (CF-007) holds on every denial.
