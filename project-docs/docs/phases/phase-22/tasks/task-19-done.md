# Phase 22 - Task 19: Explain recovery and restored work — Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None; technical repair and fixture preparation complete.
- Recovery attempts: 1
- Next action: Human completes checklist attempt 2; phase acceptance remains pending.
- Human attention required: no

## Summary of what was implemented

- Resume separates next steps and operation purposes from completion requirements, labels Actor IDs and lease facts, and explains retained uncommitted work without inventing proof.
- Run inspection displays repository comparison, paths, digest meaning and dated liveness. Verified import metadata identifies the preserved project/history/archive and destination.
- Import guidance distinguishes absent/historical/unknown liveness, an existing destination observation, and inactive/absent Runs. An active Claim never implies operation authority.
- Activation status explains resource identities, absent prior resources, both recovery directions and the activation commit versus Git; terminal states do not offer recovery.
- JSON, Events, schema, mutation rules, dependencies and gate definitions are unchanged. Historical task owners remain owners of record; DEC-219 allocates this supplemental Human guidance repair.

## Files created or modified

- cmd/fid/resume.go, resume_human.go, resume_human_guidance.go.
- cmd/fid/run_recovery.go, run_inspect_human.go.
- cmd/fid/activation.go, activation_human.go.
- internal/workflow/recovery_read_context.go and recovery_read_context_test.go.
- cmd/fid/human_recovery_guidance_test.go; integration/p22_human_guidance_test.go and p22_resume_human_test.go.

## Documentation updated

- DEC-219, public-interface section 6, task-19 spec; current/phase status and run-ledger supplement.
- Original Human checklist bookkeeping corrected without changing any Human response or outcome; separate attempt 2 has blank responses and exact fresh fixture commands.
- human-fixtures.md and independent task-19-review.md, which preserves MF-01 and appends its closure.

## Validation Contract Coverage

- VAL-053: PASS for this repair — TestP22RecoveryHumanResultMatchesDurableJSON, TestP22ResumeHumanMatchesEligibleAndStoppedRecoveryState and TestP22ResumeHumanDistinguishesElapsedLeaseFromRecordedStatus executed in the focused group and final integration gate. Independent review additionally records its own six-response byte-parity comparison.
- VAL-076: HUMAN RETEST REQUIRED — TestHumanRecoveryGuidance and TestHumanGuidanceDoesNotRecommendBlockedGate prove explanatory content and blocked-state distinction; Human comprehensibility is not inferred.
- VAL-077: HUMAN RETEST REQUIRED — TestP22HumanGuidance/activation_directions, TestHumanActivationTerminalStateDoesNotOfferRecovery and the retained TestActivationStatusHumanAndJSONDescribeRecoveryDirectionAtCommitBoundary cover both directions and terminal precedence.
- VAL-078: HUMAN RETEST REQUIRED — TestP22HumanGuidance covers missing observations, imported historical versus destination observations, invalid binding rejection, imported terminal resume/inspect and unchanged Event count. TestHumanImportedInactiveRunDoesNotRequestObservation also covers absent Run. No Human outcome is upgraded.

## Tests run

- `go test ./cmd/fid -run '^TestHumanRecoveryGuidance$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); raw `.cache/phase22-human-repair/red-renderer.log`.
- `go test ./cmd/fid -run '^TestHumanRecoveryGuidance$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.105s` (terminated `exit_code: 0`); raw `.cache/phase22-human-repair/green-renderer.log`.
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22HumanGuidance$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); raw `.cache/phase22-human-repair/red-context.log`.
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22HumanGuidance$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.695s` (terminated `exit_code: 0`); raw `.cache/phase22-human-repair/green-context.log`.
- `go test ./cmd/fid ./internal/workflow -run 'TestHuman|TestRecoveryReadContext' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); raw `.cache/phase22-human-repair/focused-unit.log`.
- `go test ./cmd/fid ./internal/workflow -run 'TestHuman|TestRecoveryReadContext' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	0.291s` (terminated `exit_code: 0`); raw `.cache/phase22-human-repair/focused-unit-2.log`.
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22(HumanGuidance|RecoveryHumanResultMatchesDurableJSON|ResumeHumanMatchesEligibleAndStoppedRecoveryState|ResumeHumanDistinguishesElapsedLeaseFromRecordedStatus)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	11.451s` (terminated `exit_code: 0`); raw `.cache/phase22-human-repair/focused-integration.log`.
- `go test ./cmd/fid -run '^TestHumanImportedFreshObservationDoesNotRequestItAgain$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); raw `.cache/phase22-human-repair/red-review.log`.
- `go test ./cmd/fid -run '^TestHumanImportedFreshObservationDoesNotRequestItAgain$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.140s` (terminated `exit_code: 0`); raw `.cache/phase22-human-repair/green-review.log`.
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22(HumanGuidance|RecoveryHumanResultMatchesDurableJSON|ResumeHumanMatchesEligibleAndStoppedRecoveryState|ResumeHumanDistinguishesElapsedLeaseFromRecordedStatus)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	15.337s` (terminated `exit_code: 0`); raw `.cache/phase22-human-repair/focused-integration-2.log`.
- `go test ./cmd/fid ./internal/workflow -run 'TestHuman|TestRecoveryReadContext' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	0.306s` (terminated `exit_code: 0`); raw `.cache/phase22-human-repair/focused-unit-3.log`.
- `go test ./cmd/fid -run '^TestActivationStatusHumanAndJSONDescribeRecoveryDirectionAtCommitBoundary$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.304s` (terminated `exit_code: 0`); raw `.cache/phase22-human-repair/green-activation-retained.log`.

## Exact quality-gate results

Working directory: `/home/mark/projects/fid`. Final passing attempts:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`); raw `.cache/phase22-human-repair/gate-format-3.log`.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`); raw `.cache/phase22-human-repair/gate-lint-3.log`.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.6% of statements` (terminated `exit_code: 0`); raw `.cache/phase22-human-repair/gate-race-2.log`.
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					65.5%` (terminated `exit_code: 0`); raw `.cache/phase22-human-repair/gate-coverage.log`.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	563.279s` (terminated `exit_code: 0`); raw `.cache/phase22-human-repair/gate-integration.log`.
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`); raw `.cache/phase22-human-repair/gate-build.log`.
- `FID_P22_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-22-human go test -tags integration ./cmd/fid/integration/... -run '^TestP22HumanFixtures$' -count=1 -v` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	18.972s` (terminated `exit_code: 0`); raw `.cache/phase22-human-repair/gate-fixtures.log`.

## Deviations from spec

- Initial behavioral red: missing adoption/next-step explanations and missing inspection/import/activation guidance each produced `FAIL` with observed exit_code:1. Their byte-identical focused commands subsequently passed; both attempts are retained above.
- New unit fixture correction: the initial focused cmd/workflow command produced `context invented history or liveness` and final token `FAIL` (terminated exit_code:1). checkpointFixture reuses recoveryFixture, which explicitly inserts OBS-terminated (recovery_test.go); the assertion incorrectly expected no liveness. It now requires the recorded terminated state and no invented import provenance. The distinct public no-observation case remains. The byte-identical focused command passed afterward.
- Independent MF-01: generic import guidance requested fresh observation even after one existed and for inactive Runs. The new regression produced `FAIL` (terminated exit_code:1). Conditioned the imperative on Run status and observation provenance; added real imported terminal and fresh-observation negative assertions. The byte-identical regression passed, and the independent reviewer closed MF-01.
- First required `go test -race -coverprofile=coverage.out ./...` failed with `FAIL` (terminated exit_code:1), cmd/fid 529.052s. TestActivationStatusHumanAndJSONDescribeRecoveryDirectionAtCommitBoundary rejects every occurrence of `rollback`, including the truthful phrase `rollback is unavailable`. Kept the test intact; changed the wording to `returning to the prior state is unavailable` and limited the absent-prior restoration footer to precommit nonterminal activations. Focused retained boundary test passed, then the byte-identical full required gate passed as recorded above. Failed log: `.cache/phase22-human-repair/gate-race.log`.
- The retained lease-rendering assertion moved from the semicolon layout to separately labeled lines while preserving exact expired/recorded-active/expiry facts and no-Event assertion. No semantic coverage was removed.
- Root used the explicitly approved single-task/direct-repair workflow with one independent reviewer and bounded closure. Protected commands, tests, tolerances and policies were not weakened.

## Independent review and Human handoff

- [task-19-review.md](../qa/task-19-review.md) preserves initial MF-01 and closes it on product source 5232d9cdfa32e71b1a54fe7a73025eac86a97a7c. Independent technical review PASS; Human judgments remain pending.
- Fresh manifest: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-386504276/manifest.json`; generated `2026-09-12T10:18:42.450198226Z`.
- Retained binary SHA-256: `sha256:e2ca68186dc96297705c519fd9849d0a041862522d8b002c96fecf3e1393fbc9`; source `ee7a5cae5fb7b572be1b8032c0ed17cb6a1afa37`, source_dirty:false.
- Root separately verified the retained binary and all twelve exact post-generation Human/JSON inspections; see `.cache/phase22-human-repair/fixture-verification.json`.
- [human-checklist-2.md](../qa/human-checklist-2.md) is the current retry. Attempt 1 remains FAIL (1 PASS, 2 FAIL, 3 PARTIAL). No phase acceptance, merge, push, release or unconditional RUN-END is claimed.

## Concerns, technical debt, or notes for the Tech Lead

- Carry CLI-wide presentation consistency forward; this task addresses the three recovery views and their accepted Human assertions.
- Full race and integration gates remain expensive. No archive/performance code changed; prior independent phase benchmark and DEC-217 scale evidence stays historical, not freshly reattested by this task.
- Fresh fixtures use the accepted 86400-second lease; regenerate if the Human review occurs after its recorded expiry.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `5232d9cdfa32e71b1a54fe7a73025eac86a97a7c`

## Commit subject

- `fix(cli): clarify forward-only activation recovery`

Initial implementation: eff28fb1b7ea0a543017bef00ccb183cc79fc396. State-guidance repair: 94afb3b27b2174078e19229d6771d89e40d740d4. Independent review checkpoint: 7d81924de15911b5008063227181422f48498095.
