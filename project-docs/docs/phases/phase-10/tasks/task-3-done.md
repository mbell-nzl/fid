## Summary of what was implemented

- Retrofitted `policy.Authorize` onto every scoped write path: `issue create`, `issue update`, `issue transition`, `issue complete`, `issue ac add`, `issue ac verify`, `concern raise`, `concern address`, and `rebuild`.
- Derived permission actions in the domain layer via `internal/issues/policy_actions.go`, keeping CLI adapters thin and preserving plain known-role validation for `owner_role`.
- Added actionable permission-denied CLI rendering that names the role and action and points operators at `fid policy show`.
- Added explicit coverage for denied-before-append writes, state-keyed transition enforcement, role-keyed AC verification enforcement, human verifier attribution, and `fid policy show` show-vs-enforce parity.
- Updated DEC-E fixtures to permitted role/action combinations and refreshed raw journal fixtures so rebuild tests use a seeded permitted role.
- Added the direct-prerequisite matrix/spec rows needed to cover already-implemented derived actions: `issue.ac.verify.tl_lite`, `issue.ac.verify.workflow_runner`, `issue.transition.wrapped`, and `issue.transition.closed`.

## Files created or modified

- `internal/issues/create.go`, `internal/issues/update.go`, `internal/issues/transition.go`, `internal/issues/complete.go`, `internal/issues/ac.go`, `internal/issues/ac_verify.go`, `internal/issues/concern_raise.go`, `internal/issues/concern_address.go`, `internal/issues/policy_actions.go`
- `cmd/fid/issue.go`, `cmd/fid/rebuild.go`
- `internal/policy/matrix.go`, `internal/policy/matrix_test.go`
- `cmd/fid/issue_test.go`, `cmd/fid/issue_ac_test.go`, `cmd/fid/issue_transition_test.go`, `cmd/fid/policy_test.go`, `cmd/fid/rebuild_test.go`, `cmd/fid/init_test.go`
- `cmd/fid/e2e_v01b_test.go`, `cmd/fid/e2e_v01c_ac_test.go`, `cmd/fid/e2e_v01c_complete_test.go`, `cmd/fid/e2e_v01c_concern_test.go`, `internal/issues/transition_test.go`
- `testdata/journal/replay-v1.golden.jsonl`, `testdata/journal/corrupt-tail.jsonl`, `testdata/journal/conflicting-precondition.jsonl`
- `project-docs/docs/phases/phase-10/decisions.md`, `project-docs/specs/03_TRD_SCHEMA.md`, `project-docs/docs/phases/phase-10/tasks/task-3-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-10/decisions.md` — added the Task 3 DEC-E execution note, including every fixture role change and the direct-prerequisite matrix additions.
- `project-docs/specs/03_TRD_SCHEMA.md` section 13 — added the missing direct-prerequisite rows for `issue.ac.verify.tl_lite`, `issue.ac.verify.workflow_runner`, `issue.transition.wrapped`, and `issue.transition.closed`.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` sections 2, 3, 6, and 11 — `verified-no-change`; the existing CLI spec already matched the enforced flag surface and command behavior for Task 3.
- `project-docs/specs/01_PRD.md` — `verified-no-change`; PRD-ROLE-001 through PRD-ROLE-004 remain MET as written.

## Validation Contract Coverage

- `VAL-004` — covered by `TestIssueCreateAndRebuildValidateRoleExistence/issue create denied role leaves journal unchanged`, `TestPolicyShowMatchesWriteEnforcementForPermittedAndDeniedCells`, and the existing successful create/rebuild tests. Evidence: denied epic create now exits non-zero with a role/action-named permission error and zero journal growth; permitted epic create succeeds.
- `VAL-005` — covered by the existing missing/unknown-role command tests that remained green after the retrofit: `TestIssueCreateAndRebuildValidateRoleExistence`, `TestIssueTransitionRequiresTargetStateAndKnownByRole`, `TestIssueACAddCommandRejectsBeforeAppend`, `TestIssueACVerifyCommandRejectsBeforeAppend`, and `TestIssueUpdateCommandRequiresKnownByRole`. Evidence: missing `--by` and unknown roles still reject before append.
- `VAL-006` — covered by `TestIssueACVerifyCommandUsesRoleSpecificPermissionsAndHumanAttribution`. Evidence: `workflow_runner` is denied for AC verify, `human` succeeds on the same write path, and acceptance read-back records the latest `human` verification under the `human` role.
- `VAL-007` — covered by `TestIssueACVerifyCommandUsesRoleSpecificPermissionsAndHumanAttribution` and `TestIssueTransitionUsesStateSpecificPermissions`. Evidence: AC verify is keyed on `issue.ac.verify.<role>` (`qa_reviewer` allowed, `workflow_runner` denied) and transition is keyed on `issue.transition.<target-state>` (`developer` allowed for `dev_complete`, `qa_reviewer` denied).
- `VAL-008` — covered by `TestPolicyShowMatchesWriteEnforcementForPermittedAndDeniedCells`. Evidence: `fid policy show` reports `issue.create.epic` as denied for `developer` and permitted for `tl`, and the actual command outcomes match those booleans.

## Tests run

- `go test ./cmd/fid -run TestIssueCreateAndRebuildValidateRoleExistence/issue_create_denied_role_leaves_journal_unchanged` — FAIL (expected red: epic create was still allowed for any known role)
- `go test ./cmd/fid -run TestIssueCreateAndRebuildValidateRoleExistence/issue_create_denied_role_leaves_journal_unchanged` — PASS
- `go test ./cmd/fid` — FAIL (expected DEC-E fallout: rebuild and invalid-transition fixtures still used now-denied roles)
- `go test ./cmd/fid` — PASS
- `go test ./internal/issues` — FAIL (expected DEC-E and matrix-completeness fallout on transition fixtures)
- `go test ./internal/issues ./internal/policy` — PASS
- `go test ./cmd/fid -run 'TestIssueACVerifyCommandUsesRoleSpecificPermissionsAndHumanAttribution|TestIssueTransitionUsesStateSpecificPermissions|TestPolicyShowMatchesWriteEnforcementForPermittedAndDeniedCells|TestMatrixContainsExactlyTheImplementedV01WriteActions'` — FAIL (expected direct-prerequisite gap: missing AC verifier matrix row for `workflow_runner`)
- `go test ./cmd/fid ./internal/policy -run 'TestIssueACVerifyCommandUsesRoleSpecificPermissionsAndHumanAttribution|TestIssueTransitionUsesStateSpecificPermissions|TestPolicyShowMatchesWriteEnforcementForPermittedAndDeniedCells|TestMatrixContainsExactlyTheImplementedV01WriteActions'` — PASS
- `go test ./cmd/fid ./internal/issues ./internal/policy` — PASS

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS
- `golangci-lint run ./...` — PASS
- `go test -race -coverprofile=coverage.out ./...` — PASS
- `go tool cover -func=coverage.out` — PASS; COV-001 satisfied: `internal/issues` 81.9%, `internal/journal` 80.3%, `internal/policy` 81.2%, `internal/workflow` 95.5%. COV-002 satisfied: `cmd/fid` 77.8%.
- `go test -tags integration -race ./...` — PASS
- `go build -o /tmp/fid ./cmd/fid` — PASS
- `go test -bench . -benchmem ./...` — PASS; required because `internal/journal` benchmark coverage includes the rebuild path, and Task 3 changed rebuild fixtures/roles under DEC-E.

## Deviations from spec

- Added direct-prerequisite matrix/spec rows for `issue.ac.verify.tl_lite`, `issue.ac.verify.workflow_runner`, `issue.transition.wrapped`, and `issue.transition.closed`. Without those rows, existing write paths would fail with missing-entry errors instead of enforced permit/deny outcomes. This stayed within the task's allowed prerequisite scope.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `5293fd2674b0afdc902feb2f7079d12b45e82551`

## Commit subject

- `feat(policy): enforce write permissions on commands`
