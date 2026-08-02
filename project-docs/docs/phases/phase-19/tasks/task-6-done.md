# Phase 19 - Task 6: Closure — integrity, rebuild parity, exit scenario, spec sync - Completion Report

## Summary of what was implemented

- Added a read-only workflow-assignment integrity check to `fid maintenance
  check`, including deterministic kind- and project-scope compatibility
  findings.
- Converted VAL-059/CF-025 into a permanent integration golden test: the
  current binary copies the committed v0.1 fixture and byte-compares all six
  committed JSON outputs. The capture script, manifest, normalization file,
  fixture, and goldens remain committed provenance only.
- Made Review 2's two improvement findings permanent: explicit pairwise
  value-naming assignment errors and per-scope assignment replay assertions.
- Added the scripted v0.2a exit scenario, DEC-079 custom-state boundary,
  reserved-rule SQL/CLI seam checks, and full-projection double-rebuild parity.
- Added exported Vacuum documentation and synchronized the CLI limitation and
  Phase 19 glossary terms.

## Files created or modified

- `cmd/fid/v01_parity_test.go` — CF-025 permanent v0.1 golden parity test
- `cmd/fid/e2e_v02a_workflow_authoring_test.go` — rebuild, maintenance,
  reserved-rule, DEC-079, and exit-scenario integration evidence
- `cmd/fid/e2e_v01_closure_helpers_test.go` — whole-projection snapshot now
  includes `workflow_assignments`
- `cmd/fid/workflow_assign_test.go` — explicit six-message VAL-016 coverage
- `internal/maintenance/check.go`, `internal/maintenance/workflow_assignments.go`,
  `internal/maintenance/workflow_assignments_test.go` — assignment integrity check
- `internal/maintenance/vacuum.go` — exported Go documentation
- `internal/workflow/assign_test.go` — kind/project replay-row assertions
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`, `project-docs/specs/GLOSSARY.md`
- `project-docs/docs/phases/phase-19/tasks/task-6-done.md`

## Documentation updated

- `06_DESIGN_PLAN_CLI.md` documents DEC-079: custom state-name transitions are
  denied until Phase 20's explicit `fid policy set` permission configuration.
- `GLOSSARY.md` now defines Custom workflow, Workflow assignment, and Initial
  state.
- Verified without correction: `02_TRD.md` §3 already assigns workflow
  authoring/assignment/resolution to `internal/workflow` and keeps its declared
  must-nots (DEC-089); `05_DESIGN_PLAN.md` §3.1 already matches the shipped
  `workflow define`, `workflow assign`, and `workflow list` command family
  (DEC-092).

## Validation Contract Coverage

- VAL-021: PASS — `TestV02aWorkflowAuthoringExitScenarioAndRebuildParity`
  (`cmd/fid/e2e_v02a_workflow_authoring_test.go`) snapshots every projected
  table (including `workflow_assignments`), `sqlite_schema`, clean
  `PRAGMA foreign_key_check`, and every journal file before deletion/rebuild,
  after the first rebuild, and after the second rebuild.
- VAL-022: PASS — `TestCheckWorkflowAssignmentsFlagsIncompatibleAppliesTo` and
  `TestCheckWorkflowAssignmentsAcceptsCompatibleKindAndProjectScopes`
  (`internal/maintenance/workflow_assignments_test.go`) prove the incompatible
  and healthy states; `TestMaintenanceCheckReportsIncompatibleWorkflowAssignment`
  (`cmd/fid/e2e_v02a_workflow_authoring_test.go`) proves the real CLI exits
  non-zero with `--warn-as-error` and names `workflow_assignment` / `kind:epic`.
- VAL-023: PASS — `assertReservedWorkflowRulesRejected`
  (`cmd/fid/e2e_v02a_workflow_authoring_test.go`) proves both distinct CLI
  rejections (`required_gates_current`/`v0.2b`, `has_completion`/`v0.2c`) and
  unchanged journal bytes.
- VAL-024: PASS — `assertCustomStateTransitionBoundary`
  (`cmd/fid/e2e_v02a_workflow_authoring_test.go`) proves a custom `backlog`
  state fails with the documented no-matrix-entry error and no append, while a
  `task-strict` workflow using `in_progress` transitions successfully.
- VAL-025: PASS — `TestV02aWorkflowAuthoringExitScenarioAndRebuildParity`
  defines and assigns `task-strict`, lists it as JSON/TOON/Markdown/YAML,
  creates the selected task, blocks `in_progress -> dev_complete` on the named
  open concern, addresses it, then succeeds while the earlier `task-default`
  control remains `todo` under `task-default`.
- VAL-026: PENDING (Human-only) — no Developer substitute was run; the final
  phase human checklist must record the operator-coherence judgment.
- VAL-068: PASS — `assertReservedRulesAreNotSeeded`
  (`cmd/fid/e2e_v02a_workflow_authoring_test.go`) runs the ordered SQL query
  and asserts `[]`; a built `/tmp/fid` entrypoint rerun also printed `[]`.
- VAL-069: PASS — exact production inventory returned no output with exit 1:
  `rg -n --glob '*.go' --glob '!**/*_test.go' '\\brequired_gates_current\\b|\\bhas_completion\\b' cmd/fid internal/issues internal/dev internal/journal internal/db internal/maintenance internal/policy`.
- VAL-070: PASS — exact production inventory returned no output with exit 1:
  `rg -n --glob '*.go' --glob '!**/*_test.go' --glob '!internal/workflow/definitions_fixture.go' 'FrozenDefinitions\\(' cmd internal`.

## Tests run

- CF-025 first: `go test -tags integration ./cmd/fid -run '^TestV01ParityGoldenOutputs$' -count=1` — PASS — `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.650s`; six current-binary stdout values byte-matched their committed JSON goldens.
- Red: `go test ./internal/maintenance -run '^TestCheckWorkflowAssignments' -count=1` — expected FAIL — `want workflow_assignment finding for kind:epic` before `checkWorkflowAssignments` existed.
- Green: `go test ./internal/maintenance -run '^TestCheckWorkflowAssignments' -count=1` — PASS — `ok  \tgithub.com/mbell-nzl/fid/internal/maintenance\t0.040s`.
- Review 2 improvement test: `go test ./cmd/fid -run '^TestWorkflowAssignRejectsWithoutMutation$' -count=1` — PASS — `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.034s`; all six expected errors are explicit and pairwise distinct.
- Review 2 improvement test: `go test ./internal/workflow -run '^TestAssignEmitsScopeAwareReplacementEffects$' -count=1` — PASS — `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.092s`; both kind and project replay rows survive rebuild.
- Integration scenario: `go test -tags integration ./cmd/fid -run '^(TestV02aWorkflowAuthoringExitScenarioAndRebuildParity|TestMaintenanceCheckReportsIncompatibleWorkflowAssignment)$' -count=1` — PASS — `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t1.749s`.
- Built CLI entrypoint: `/tmp/fid init && /tmp/fid sql "SELECT workflow_id, from_state, to_state, rule FROM workflow_transition_requires WHERE rule IN ('required_gates_current', 'has_completion') ORDER BY workflow_id, from_state, to_state, rule" --format json` in a fresh workspace — PASS — `[]`.
- CF-024 control: `go test -run '^$' -bench '^BenchmarkList10kIssuesWithRegistry$' -benchtime=1x -benchmem ./internal/issues` — PASS — `ok  \tgithub.com/mbell-nzl/fid/internal/issues\t2.044s`; `44.29 ms/registry-read`, 455.71 ms (91.14%) under the <500 ms control.
- CF-024 control: `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime=1x -benchmem ./internal/journal` — PASS — `ok  \tgithub.com/mbell-nzl/fid/internal/journal\t26.826s`; `24749 ms/rebuild`, 5.251 s (17.50%) under the <30 s control and 2.241–2.911 s faster than the Phase 18 26.99–27.66 s baseline.
- VAL-069 inventory — PASS (no match) — `exit 1 (no output)`.
- VAL-070 inventory — PASS (no match) — `exit 1 (no output)`.

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS — `exit 0 (no output)`
- `golangci-lint run ./...` — PASS — `0 issues.`
- `go test -race -coverprofile=coverage.out ./...` — PASS — `exit 0 (no output)`
- `go tool cover -func=coverage.out` — PASS — `total:\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t\t\t76.8%`
- `go test -tags integration -race ./...` — PASS — `exit 0 (no output)`
- `go test -bench . -benchmem ./...` — PASS — `PASS`; the exact gate completed and the focused CF-024 controls above retain the terminating benchmark measurements.
- `go build -o /tmp/fid ./cmd/fid` — PASS — `exit 0 (no output)`

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
