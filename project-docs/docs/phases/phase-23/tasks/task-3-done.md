# Phase 23 - Task 3: Non-mutating review Claim mode and Reviewer Run - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: All Task 3 required gates pass.
- Recovery attempts: 0
- Next action: Phase Runner review.
- Human attention required: no

## Summary of what was implemented

- `claim.acquire` now accepts `mutation` and `review` modes while retaining the one-active-Claim-per-Issue rule.
- Review Claims start, close, abandon clean state, and recover clean state under the retained Phase 22 Run lifecycle.
- Review Claims reject dirty-path adoption in `run.recover` and `run.abandon` with the closed `claim-mode-invalid` failure envelope and no Event.
- The strict journal Claim grammar accepts both modes, and CLI help no longer describes acquisition as mutation-only.

## Files created or modified

- `internal/claims/generic.go`, `internal/workflow/{authority.go,recovery.go,run_close.go}`, `internal/cli/codes.go`, and `internal/journal/current_event_value_catalog.go`.
- `cmd/fid/{claim.go,run.go,generic_claim_test.go,contract_surface_test.go}` and matching workflow/CLI tests.
- `cmd/fid/integration/{p22_fake_test.go,p23_claim_mode_test.go}`.
- `project-docs/docs/phases/phase-23/tasks/task-3-done.md`.

## Documentation updated

- Completion report only; the implementation conforms to the accepted Phase 23 contracts without a contract or schema change.

## Validation Contract Coverage

- VAL-011: `cmd/fid/generic_claim_test.go:TestClaimAcquireAcceptsReviewMode` proves review acquisition and response serialization. `TestClaimAcquireKeepsOneActiveClaimAcrossModes` proves mutation-under-review rejection, mutation-after-review-release acceptance, and review-after-mutation-release acceptance.
- VAL-012: `internal/workflow/claim_mode_test.go:TestReviewClaimRejectsDirtyPathAdoption` proves the typed review-mode adoption fence. The real-process `cmd/fid/integration/p23_claim_mode_test.go:TestP23ReviewClaimRejectsDirtyPathAdoption` proves both `run.recover` and `run.abandon` return `claim-mode-invalid` at exit 4 and leave the journal Event count unchanged.
- VAL-013: `cmd/fid/integration/p23_claim_mode_test.go:TestP23ReviewRunRetainsPhase22Lifecycle` proves real-process review Run close, clean abandonment, and clean recovery; it asserts that the recovered Claim remains `mode:"review"`.

## Tests run

- Focused red: `go test -race -timeout 30m ./cmd/fid/ -run '^TestClaimAcquireAcceptsReviewMode$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); the mutation-only closed decoder rejected the new `mode` field.
- Focused green: `go test -race -timeout 30m ./cmd/fid/ -run '^(TestClaimAcquireAcceptsReviewMode|TestClaimAcquireKeepsOneActiveClaimAcrossModes)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 11.354s` (terminated `exit_code: 0`).
- Focused red: `go test -race -timeout 30m ./internal/workflow/ -run '^TestReviewClaimRejectsDirtyPathAdoption$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); the mode fence and its typed error did not exist.
- Focused green: `go test -race -timeout 30m ./internal/workflow/ -run '^TestReviewClaimRejectsDirtyPathAdoption$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 3.361s` (terminated `exit_code: 0`).
- Real-process review lifecycle: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/ -run '^TestP23Review(RunRetainsPhase22Lifecycle|ClaimRejectsDirtyPathAdoption)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 8.049s` (terminated `exit_code: 0`).
- Direct CLI entrypoint: `go build -o /tmp/fid-task3 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`); `/tmp/fid-task3 claim acquire --help` - PASS - `--project string   fid project root` (terminated `exit_code: 0`), with the observed description `Acquire one Claim.`

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 63.8% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.8%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 702.407s` (terminated `exit_code: 0`)
- Observed wall times: `cmd/fid` `438.942s` from the fresh final race gate; `cmd/fid/integration` `702.407s` from the fresh final tagged gate; `internal/workflow` `237.598s` from the current-session uncached race gate before the final CLI-help-only correction. All are within the 600 s ordinary-package and 900 s tagged-integration budgets.

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
