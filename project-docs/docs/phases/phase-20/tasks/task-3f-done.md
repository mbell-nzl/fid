# Phase 20 - Task 3f: Dedicated integration package and gate migration - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: The MF-01 report-only repair records fresh terminating evidence for every exact Task 3f command, preserves historical failed and UNVERIFIED evidence, and completes canonical report conformance.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added a source-level classification regression that requires exactly four
  accepted real-process journeys in `cmd/fid/integration` and rejects tagged
  tests in the parent command package.
- Moved the four real-binary journeys and their shared helper into the
  dedicated package with only package-relative fixture/path adjustments.
- Removed the integration tag from 12 in-process tests so the ordinary suite
  owns their command, SQLite, journal, and rebuild assertions.
- Migrated canonical, template, active QA, and remaining Phase 20 task gate
  surfaces to `go test -tags integration -race ./cmd/fid/integration/...`.
- Ran both the former repository-wide and new dedicated commands. Both passed;
  the former still spent 152.981s rerunning `cmd/fid`, while the dedicated
  journey package itself took 10.266s.

## Files created or modified

- `cmd/fid/integration/classification_test.go`
- `cmd/fid/integration/e2e_v01_closure_helpers_test.go`
- `cmd/fid/integration/e2e_v01_closure_test.go`
- `cmd/fid/integration/e2e_v02a_workflow_authoring_test.go`
- `cmd/fid/integration/v01_parity_test.go`
- Ten reclassified `cmd/fid/*_test.go` files recorded by the anchor inventory
- `AGENTS.md`
- `project-docs/specs/02_TRD.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/09_AI_WORKFLOW.md`
- `project-docs/templates/phase-runner-prompt.md`
- `project-docs/templates/qa-plan-template.md`
- `project-docs/templates/final-review-template.md`
- `project-docs/docs/phases/phase-20/qa/qa-plan.md`
- `project-docs/docs/phases/phase-20/tasks/task-4.md` through `task-18.md`
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/tasks/task-3f-done.md`

## Documentation updated

- Recorded the test-only real-process adapter boundary in `02_TRD.md`.
- Updated all canonical and active integration gate surfaces while preserving
  completed historical task/review artifacts and the historical Review 1 row.
- Recorded every before/after test name, timing, and classification in the
  durable pack anchor.

## Validation Contract Coverage

- No VAL-NNN applies. Existing scenario assertions are unchanged. The
  classification regression and both transitional commands prove delivery
  equivalence; the task changes test ownership and gate scope only.

## Tests run

- `go test -tags integration ./cmd/fid/integration` - EXPECTED RED - `FAIL`.
- `go test -tags integration -run '^TestIntegrationClassification$' ./cmd/fid/integration` - UNVERIFIED - Historical focused-test token retained from the original report without an observed termination; this report-only continuation did not rerun the non-required focused command.
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-3f-done.md project-docs/docs/phases/phase-20/tasks/task-3f.md` - FAIL - `FAIL report-lint found 27 failure(s)` (terminated `exit_code: 27`)

## Exact quality-gate results

- `project-docs/tools/workflow-lint.sh` - PASS - `PASS workflow-lint completed with no failures` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					78.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow	(cached)` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration	(cached)` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- The before inventory contains 14 tagged files rather than the stale QA-plan
  feasibility note's 15: one helper plus 13 test files. It still contains the
  planned 16 top-level tests, classified exactly as four real-process and 12
  in-process tests. The active QA-plan row is corrected; completed historical
  artifacts are unchanged.
- The task scope named only `qa-plan-template.md` and the Phase Runner template,
  but `final-review-template.md` is also an active canonical gate surface. It
  was updated to avoid reintroducing the retired repository-wide command.
- The MF-01 pre-repair report-lint invocation failed with `FAIL report-lint
  found 27 failure(s)` (terminated `exit_code: 27`). Root cause: this
  historical report lacked the canonical typed workflow outcome and current
  evidence-row syntax. This report-only continuation adds fresh terminating
  evidence for every required command, retains the historical expected-red and
  marks the prior focused PASS without captured termination as UNVERIFIED; it
  changes no task command, product behavior, test, assertion, or validator.

## Concerns, technical debt, or notes for the Tech Lead

- The dedicated package still builds the binary once per journey. Task 3g owns
  the already-approved package-lifecycle build reuse.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `ec72fdad4cf1e21be82b22644b382825b91d3b0f`

## Commit subject

- `test(integration): isolate real-process journeys`
