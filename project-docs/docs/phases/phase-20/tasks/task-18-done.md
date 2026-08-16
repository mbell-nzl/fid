# Phase 20 - Task 18: Journey acceptance and documentation - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: Phase-end QA and the Human checklist.
- Human attention required: no

## Summary of what was implemented

- Added a literal Cobra leaf-inventory guard for the exact 21 retained and 21 legacy Phase 20 commands. It also confirms the retained root-flag registry is exact and that both root flags reject before every legacy handler.
- Added the standalone Issue acceptance journey. It uses only `--json` for initialization, hierarchy/relationship mutations, collection reads, revisioned mutation and typed failure, maintenance inspection, rebuild, and byte-identical canonical `issue.get` reads.
- Updated the operator playbook and the Unreleased changelog to describe the delivered JSON surface and removals while preserving the historical 0.1.0 record.

## Files created or modified

- `cmd/fid/contract_surface_test.go`
- `cmd/fid/e2e_p20_issue_journey_test.go`
- `project-docs/docs/playbook.md`
- `CHANGELOG.md`
- `project-docs/docs/phases/phase-20/tasks/task-18-done.md`
- Scope check: the pre-existing modified `project-docs/docs/phases/phase-20/run-ledger.md` was left unstaged and unmodified by this task.

## Documentation updated

- `project-docs/docs/playbook.md:257` documents the retained Response Envelope surface, legacy root-flag rejection, and current audit commands.
- `CHANGELOG.md:10` records the delivered Phase 20 JSON journey and removals, while `CHANGELOG.md:26` marks the 0.1.0 formatter/stub claims as historical context.
- `AGENTS.md` required no change: the Task 16 architecture-drift correction already lists the current adapter boundaries and no formatter package.

## Validation Contract Coverage

- VAL-035: PASS - `TestLiveLeafCommandsMatchContractInventories` at `cmd/fid/contract_surface_test.go:21` pins every literal retained leaf (`:22-44`) and legacy leaf (`:45-67`), requires a 21/21 partition (`:69-76`), matches the exact root-flag capability registry (`:79-85`), and executes both root-flag rejections for every legacy invocation (`:87-152`). The qualifying retained transport and non-qualifying legacy rejection states are both covered.
- VAL-036: PASS - `TestP20StandaloneIssueJourneyJSONAndRebuildParity` at `cmd/fid/e2e_p20_issue_journey_test.go:13` performs the complete JSON journey. It asserts hierarchy/mutations and Event identities (`:17-89`), ordered issue/label/AC/link collections and ready selection (`:91-115`), revisioned update/transition and typed stale failure (`:117-145`), then compares raw canonical `issue.get` bytes for every created Issue before and after rebuild (`:147-156`). The real built CLI independently produced `P20 real JSON journey passed`.
- VAL-037: PASS - the delivered-surface/removal statements are at `project-docs/docs/playbook.md:257-290` and `CHANGELOG.md:10-28`; neither file retains a current universal-format or promised-command claim. Review 6's bounded changelog improvement is therefore closed without rewriting the historical 0.1.0 entry.
- VAL-008 and VAL-017: PASS - the journey helper requires stdout-only successful `fid-cli.v1` envelopes with a trailing newline (`cmd/fid/e2e_p20_issue_journey_test.go:184-213`), distinguishes mutating Event identities from read `event:null` (`:166-182`), and verifies a stderr-only typed failure (`:265-290`).
- VAL-041: PASS - the update's Event ID is located in the journal and asserted to retain the supplied local `developer` attribution (`cmd/fid/e2e_p20_issue_journey_test.go:127-129,293-323`); the journey creates no trusted Evidence, Human Outcome, or protected-authority surface.
- VAL-044: PASS - the exact retained registry and legacy pre-handler rejection matrix are exercised by `TestLiveLeafCommandsMatchContractInventories` (`cmd/fid/contract_surface_test.go:21-152`), while the journey reaches valid canonical requests across init, Issue, label, link, ready, workflow, maintenance, and rebuild operations (`cmd/fid/e2e_p20_issue_journey_test.go:17-151`).
- VAL-045: PASS - the acceptance journey proves deterministic Issue sequence order, label name order, AC ordinal order, and link direction/relation order with matching collection counts (`cmd/fid/e2e_p20_issue_journey_test.go:91-105`). The repository race gate also executes the existing at/over-limit no-truncation fixtures, including `TestLinkListJSONHonorsAtAndOverLimitWithoutTruncation` in `cmd/fid/link_limit_test.go`.
- VAL-038: SUPPORTED ONLY - the required discrete Human observations remain exactly `HUMAN-038a`, `HUMAN-038b`, and `HUMAN-038c` in `project-docs/docs/phases/phase-20/qa/qa-plan.md:303-305`. This task supplies deterministic journey setup and does not pre-answer any Human judgment.

## Tests run

- `go test ./cmd/fid -run '^TestLiveLeafCommandsMatchContractInventories$' -count=1` - EXPECTED RED - `FAIL` (terminated `exit_code: 1`); first red test referenced the not-yet-defined `liveExecutableLeafCommands` helper.
- `go test ./cmd/fid -run '^TestLiveLeafCommandsMatchContractInventories$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.056s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestLiveLeafCommandsMatchContractInventories$' -count=1` - EXPECTED RED - `FAIL` (terminated `exit_code: 1`); the permitted temporary local test-harness probe appended a twenty-second legacy item and the guard reported `contract inventory counts = retained 21, legacy 22; want 21 each`. The probe was reverted before the following green run and before any commit.
- `go test ./cmd/fid -run '^TestLiveLeafCommandsMatchContractInventories$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.055s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestP20StandaloneIssueJourneyJSONAndRebuildParity$' -count=1` - EXPECTED RED - `FAIL` (terminated `exit_code: 1`); first red journey test referenced the not-yet-defined `p20JourneyJSON` helper.
- `go test ./cmd/fid -run '^TestP20StandaloneIssueJourneyJSONAndRebuildParity$' -count=1` - REPAIR ATTEMPT - `FAIL` (terminated `exit_code: 1`); the first helper implementation had two local test-harness compile errors: `undefined: bytesNewReader` and `undefined: scanner`.
- `go test ./cmd/fid -run '^TestP20StandaloneIssueJourneyJSONAndRebuildParity$' -count=1` - REPAIR ATTEMPT - `FAIL` (terminated `exit_code: 1`); the first link expectation reported `link length = 3, want 2`, preserving the observed inverse `relates_to` entry before the assertion was corrected to the current ordered projection.
- `go test ./cmd/fid -run '^TestP20StandaloneIssueJourneyJSONAndRebuildParity$' -count=1` - REPAIR ATTEMPT - `FAIL` (terminated `exit_code: 1`); the journal decoder looked for `id` rather than the event's `event_id` field and reported `journal has no event "EVT-20260816-000016"`.
- `go test ./cmd/fid -run '^TestP20StandaloneIssueJourneyJSONAndRebuildParity$' -count=1` - REPAIR ATTEMPT - `FAIL` (terminated `exit_code: 1`); the initial fixture attempted a retained V1-guarded `todo -> in_progress` transition without a live Claim and received `transition-not-permitted` with `live_claim_held`. The accepted journey instead uses the valid revisioned `todo -> cancelled` transition; no product behavior was changed.
- `go test ./cmd/fid -run '^(TestLiveLeafCommandsMatchContractInventories|TestP20StandaloneIssueJourneyJSONAndRebuildParity)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.294s` (terminated `exit_code: 0`)
- Built CLI interaction in a fresh disposable `/tmp` project: `/tmp/fid --json init`, legal hierarchy creation, AC/label/link mutation, list/get/ready, revisioned update/transition, typed stale update, maintenance check, rebuild, and byte comparison of canonical `issue.get` - PASS - `P20 real JSON journey passed` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					78.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.591s` (terminated `exit_code: 0`)
- `go test -bench . -benchmem ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	1.964s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- Permitted harness deviation: the exact inventory guard was temporarily made red by appending `temporary-harness-probe` to the local legacy list. The observed result was `contract inventory counts = retained 21, legacy 22; want 21 each` with `exit_code: 1`; the edit was reverted before the green run and will not be committed.
- No product, architecture, data-model, public-interface, dependency, command, threshold, or assertion-strength deviation. The acceptance-test repairs above corrected only test-harness assumptions and preserved the required JSON journey, typed failure, and byte-level rebuild comparison.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
