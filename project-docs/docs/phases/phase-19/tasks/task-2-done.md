# Phase 19 - Task 2: Workflow authoring core (validation + effect builder) - Completion Report

## Summary of what was implemented

- Added the pure workflow-definition validator and deterministic journal-effect builder, including ordered state, transition, and requirement effects; rule and parameter validation; DEC-090 cancelled-state policy; materialized cancel edges; and the exported initial-state helper.
- Added the read-only assignment-request validator for Task 4's scope matrix.
- Applied DEC-096's test-only package repair: moved the two workflow-importing migration parity tests and their private assertion closure into external `package db_test`, with workflow-scoped local DB helpers. The production `[]journal.Effect` API is unchanged.

## Files created or modified

- `internal/workflow/authoring.go`
- `internal/workflow/authoring_effects.go`
- `internal/workflow/assignment.go`
- `internal/workflow/rules.go`
- `internal/workflow/authoring_test.go`
- `internal/workflow/authoring_effects_test.go`
- `internal/workflow/assignment_test.go`
- `internal/db/workflow_tables_test.go`
- `internal/db/workflow_seed_parity_test.go`
- `project-docs/docs/phases/phase-19/tasks/task-2-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-19/tasks/task-2-done.md` — completion evidence and DEC-096 authorization record.

## Validation Contract Coverage

- VAL-006: PASS - `TestBuildDefinitionEffectsIsOrderedAndDeterministic` in `internal/workflow/authoring_effects_test.go:10` deep-compares the golden ordered effect slice and a second build of the same request.
- VAL-007: PASS - `TestValidateDefinitionRequestIdentifierMatrix` in `internal/workflow/authoring_test.go:15` covers invalid and duplicate slugs, invalid and duplicate states, plus the accepted control with distinct value-naming errors.
- VAL-008: PASS - `TestBuildDefinitionEffectsCancelledStatePolicy` in `internal/workflow/authoring_test.go:125` asserts every materialized cancel edge and explicit-edge deduplication.
- VAL-032: PASS - `TestValidateDefinitionRequestStatesAndTransitions` in `internal/workflow/authoring_test.go:45` covers empty states, states/categories count mismatch, and unknown category.
- VAL-033: PASS - `TestValidateDefinitionRequestStatesAndTransitions` in `internal/workflow/authoring_test.go:45` covers an undeclared transition state and a duplicate transition with edge-naming errors.
- VAL-034: PASS - `TestDefinitionRuleVocabularyAndParams` in `internal/workflow/authoring_test.go:68` verifies all four implemented rules are accepted and appear on their declared `workflow_transition_requires` effects with their params.
- VAL-035: PASS - `TestDefinitionRuleVocabularyAndParams` in `internal/workflow/authoring_test.go:68` exercises the three reserved rule names with v0.2b/v0.2c/v0.2d tokens and an unknown rule with its offending value.
- VAL-036: PASS - `TestDefinitionRuleVocabularyAndParams` in `internal/workflow/authoring_test.go:68` rejects invalid JSON and empty roles, and asserts verbatim valid `all_acs_verified` params in the effect.
- VAL-037: PASS - `TestValidateDefinitionRequestRequirementEdges` in `internal/workflow/authoring_test.go:109` rejects an undeclared edge and accepts a declared edge; `TestBuildDefinitionEffectsAcceptsRequirementOnMaterializedEdge` in `internal/workflow/authoring_effects_test.go:64` proves an allowed materialized edge.
- VAL-038: PASS - `TestBuildDefinitionEffectsZeroCancelledApplicability` in `internal/workflow/authoring_test.go:144` proves accepted zero-cancel `escalation` and `qa_review` definitions without cancel edges, and rejected `epic`, `task`, `subtask`, and `any` definitions naming PRD-WF-013.
- VAL-039: PASS - `TestBuildDefinitionEffectsRejectsMultipleCancelledStates` in `internal/workflow/authoring_test.go:167` names both conflicting cancelled states and asserts no effect slice.
- VAL-040: PASS - `TestInitialState` in `internal/workflow/authoring_test.go:183` proves the first-`unstarted` and first-declared branches.
- VAL-041: PASS - `TestValidateAssignmentRequestRejectsAndAcceptsScopeMatrix` in `internal/workflow/assignment_test.go:11` covers five distinct rejections and compatible kind/project controls.
- VAL-042: PASS - `TestAuthoringCoreDoesNotWrite` in `internal/workflow/authoring_test.go:192` and `TestValidateAssignmentRequestDoesNotWrite` in `internal/workflow/assignment_test.go:70` run with `PRAGMA query_only` and compare workflow-table snapshots plus journal bytes. Fresh inventory `rg -n --glob '*.go' 'INSERT|UPDATE|DELETE' internal/workflow/authoring.go internal/workflow/authoring_effects.go internal/workflow/assignment.go` was PASS (no match) - `exit 1 (no output)`.

## Tests run

- `go test ./internal/db -count=1` - FAIL (expected red before DEC-096 repair) - `imports github.com/mbell-nzl/fid/internal/db from project.go: import cycle not allowed in test`
- `go test ./internal/db -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/db\t0.464s`
- `go test ./internal/workflow -count=1` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t0.248s`
- `go test -race -cover -count=1 ./internal/workflow` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t3.785s\tcoverage: 90.2% of statements`
- `rg -n --glob '*.go' 'INSERT|UPDATE|DELETE' internal/workflow/authoring.go internal/workflow/authoring_effects.go internal/workflow/assignment.go` - PASS (no match) - `exit 1 (no output)` — VAL-042 expected no-match result under DEC-094.
- Runtime entrypoint: N/A — Task 2 explicitly excludes CLI and write-path surfaces; it supplies the pure core consumed by later tasks.
- QA-plan applicability: QA-06..08 and QA-37..47 are evidenced by the named focused core tests above.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`
- `golangci-lint run ./...` - PASS - `0 issues.`
- `go test -race -coverprofile=coverage.out ./...` - PASS - `exit 0 (no output)`
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t77.1%`

## Deviations from spec

- Authorized scope extension: DEC-096 permits this test-only split of the two `internal/workflow`-importing migration parity tests into external `package db_test`. No production file under `internal/db` or `internal/journal` changed, `driverName` remains unexported, and the `[]journal.Effect` API is unchanged.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `10d4d242403851b5c676a4308cde598d251516ea`

## Commit subject

- `feat(workflow): add authoring validation and effects`
