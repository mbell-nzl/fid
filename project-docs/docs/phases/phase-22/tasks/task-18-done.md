# Phase 22 - Task 18: Final recovery conformance and Human validation fixture - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: Independent Review6, then finalQA and Human checklist readiness.
- Human attention required: no

## Summary of what was implemented

- Added the independently pinned 23-action profile/grant acceptance and rejection matrix, an operation/error/frame inventory linked to executable behavioral cases, and actual failure requests with closed error-shape checks.
- Completed the same-Issue public journey through committed denial/replay, Budget observation, embedded Stop resolution/replay, Attention disposition, proof, interruption, reconciliation/checkpoint, linked recovery, adoption, stale then fresh proof, export/validate/import, rebuild and resume.
- Pinned all ten portable semantic-summary fields independently and compare complete fake/fid summaries for eight retained state cases. The fake filters current-Run Stop/Attention views by explicit subject, while full canonical entry-byte and named fid-view parity preserve historical facts.
- Implemented DEC-205's canonical embedded Stop disposition sequence, DEC-206's committed budget-stop live control and Human notice, and DEC-207's three already accepted error mappings. No contract or gate was weakened.
- Added a child-process Human fixture generator with retained binary/hash, source provenance, exact identities, neutral commands, lease expiry, fresh directories and parent-process inspection after generator exit. Four checkpoints cover recovery/adoption, both activation boundaries and imported intervention. Human outcomes are deliberately unrecorded.

## Files created or modified

- `cmd/fid/contract_surface_test.go`
- `cmd/fid/integration/p22_activation_recovery_test.go`
- `cmd/fid/integration/p22_fake_test.go`
- `cmd/fid/integration/p22_portable_conformance_test.go`
- `cmd/fid/integration/p22_portable_fake_test.go`
- `cmd/fid/integration/testdata/p22_conformance/manifest.json`
- `cmd/fid/integration/testdata/p22_conformance/portable/abandoned.json`
- `cmd/fid/integration/testdata/p22_conformance/portable/active.json`
- `cmd/fid/integration/testdata/p22_conformance/portable/deduplicated-attention.json`
- `cmd/fid/integration/testdata/p22_conformance/portable/exhausted-budget.json`
- `cmd/fid/integration/testdata/p22_conformance/portable/interrupted.json`
- `cmd/fid/integration/testdata/p22_conformance/portable/missing-artifact.json`
- `cmd/fid/integration/testdata/p22_conformance/portable/partial-activation.json`
- `cmd/fid/integration/testdata/p22_conformance/portable/unknown-budget.json`
- `cmd/fid/live_output.go`
- `cmd/fid/run.go`
- `cmd/fid/run_recovery.go`
- `internal/cli/codes.go`
- `internal/cli/codes_test.go`
- `internal/evidence/execution.go`
- `internal/workflow/stop.go`
- `cmd/fid/integration/p22_action_conformance_test.go`
- `cmd/fid/integration/p22_budget_frames_test.go`
- `cmd/fid/integration/p22_conformance_test.go`
- `cmd/fid/integration/p22_final_journey_test.go`
- `cmd/fid/integration/p22_human_fixtures_test.go`
- `cmd/fid/integration/p22_wire_failure_test.go`
- `cmd/fid/integration/testdata/p22_conformance/wire-inventory.json`
- `cmd/fid/live_budget_test.go`
- `internal/evidence/budget_frames.go`
- `project-docs/docs/phases/phase-22/qa/human-fixtures.md`

## Documentation updated

- `qa/human-fixtures.md`: neutral generator/manifest usage and expiry/regeneration instructions.
- `tasks/task-18.md` and `decisions.md`: preceding docs-only DEC-205/206/207 allocations; these did not change any task command, accepted contract, VAL assertion or validator.
- This completion report records original-task provenance and captured execution receipts.

## Validation Contract Coverage

- VAL-053: PASS - `TestP22RecoveryHumanResultMatchesDurableJSON`, `TestP22ResumeHumanMatchesEligibleAndStoppedRecoveryState`, `TestP22ResumeHumanDistinguishesElapsedLeaseFromRecordedStatus` in `p22_resume_human_test.go`; generated recovery and activation Human views also have paired JSON assertions before neutral instructions are written.
- VAL-073: PASS - `TestP22DualAdapterRecoveryConformance`, `TestP22DualAdapterActivationConformance`, and `TestP22PortableRecoveryConformance` execute shared accepted/rejected cases against fake and fid. The latter pins ten fields, eight states and full summary equality; `p22PortableThreeWay` separately preserves full named fid reads and canonical entry bytes across source/imported/rebuilt states. Fake code receives semantic inputs and never opens the backend.
- VAL-074: PASS - `TestP22CompleteOneIssueRecoveryJourney` connects every required step on one Issue, including committed denial and same-key replay, durable Stop disposition/replay and separate Attention resolution, interrupted predecessor/successor identity, adopted work staying stale until fresh Gate/Receipt/Evidence, and the full portable roundtrip.
- VAL-075: PASS, subject to independent QA sufficiency review - `TestP22ConformanceCorpusHasEveryAcceptedOperationAndFailure` binds 32 operations,22 failure codes,10 frame variants and12 quoted closed-value rows to actual behavioral tests. `TestP22FinalActivatedActionSetAcceptsProfilesAndGrants` independently pins all23 names and executes actual profile/grant acceptance plus unknown/retired/bootstrap exclusions. `TestP22AmendedFailureWireCorpus` supplies actual closed failures including non-null observation rejection and recovery duplicate subject. `TestP22BudgetCrossingFrameOrder` proves no control before crossing, committed control before normal terminal and terminal-only replay with no Event. Inventory reference presence alone is not behavioral evidence; its referenced tests run under the full ordinary/tagged gates.
- VAL-076, VAL-077, VAL-078: Human judgments PENDING. Readiness only: `TestP22HumanFixtures` generates four final-path checkpoints in a subprocess; after it exits, the parent validates retained binary SHA256 and executes every Human/JSON inspection. The exact retained-root gate below records the durable manifest. No response/outcome has been filled for the Human.
- VAL-080: PENDING independent finalQA; this implementation report does not self-certify the independent cross-product comparison. DEC-187 and the accepted kit pin/dispositions remain unchanged.

## Tests run

- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22ConformanceCorpusHasEveryAcceptedOperationAndFailure

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 52.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					63.0%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	101.661s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `FID_P22_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-22-human go test -tags integration ./cmd/fid/integration/... -run '^TestP22HumanFixtures

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- Retained manifest: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-2335511884/manifest.json`. Generated from clean source commit `24e3cbdb22dd9e43b55636da313dcef1ef633b6f`; final root observed all seven exact commands terminate exit_code0. Aggregate ordinary coverage:63.0%.
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `24e3cbdb22dd9e43b55636da313dcef1ef633b6f`

## Commit subject

- `feat(recovery): complete conformance and retained Human fixtures`
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22FinalActivatedActionSetAcceptsProfilesAndGrants

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.295s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22PortableRecoveryConformance

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.074s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22HumanFixtures

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1 -v` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./internal/workflow -run '^Test(EquivalentStopsShareAttentionButKeepDistinctFacts|AttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable)

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	0.282s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22HumanFixtures

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1 -v` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.385s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22BudgetCrossingFrameOrder

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22BudgetCrossingFrameOrder

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22BudgetCrossingFrameOrder

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.719s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestBudgetStopControlHumanNamesDurableIntervention

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestBudgetStopControlHumanNamesDurableIntervention

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.096s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22AmendedFailureWireCorpus

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22AmendedFailureWireCorpus

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22AmendedFailureWireCorpus

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.604s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22ConformanceCorpusHasEveryAcceptedOperationAndFailure

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./internal/cli -run '^TestReasonCodeDetailSchemaIsExactForEveryCode

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./internal/cli -run '^TestReasonCodeDetailSchemaIsExactForEveryCode

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	0.021s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22ConformanceCorpusHasEveryAcceptedOperationAndFailure

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.507s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22(BudgetCrossingFrameOrder|CompleteOneIssueRecoveryJourney|PortableRecoveryConformance|HumanFixtures|AmendedFailureWireCorpus|ConformanceCorpusHasEveryAcceptedOperationAndFailure|FinalActivated)' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22(BudgetCrossingFrameOrder|CompleteOneIssueRecoveryJourney|PortableRecoveryConformance|HumanFixtures|AmendedFailureWireCorpus|ConformanceCorpusHasEveryAcceptedOperationAndFailure|FinalActivated)' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	30.124s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22PortableRecoveryConformance

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	12.863s` (terminated `exit_code: 0`)

Preserved first required-gate attempt (superseded after the documented inventory repair):

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`)

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1 -v` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	17.507s` (terminated `exit_code: 0`)

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22FinalActivatedActionSetAcceptsProfilesAndGrants

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.295s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22PortableRecoveryConformance

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22CompleteOneIssueRecoveryJourney

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.074s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22HumanFixtures

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1 -v` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./internal/workflow -run '^Test(EquivalentStopsShareAttentionButKeepDistinctFacts|AttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable)

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	0.282s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22HumanFixtures

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1 -v` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.385s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22BudgetCrossingFrameOrder

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22BudgetCrossingFrameOrder

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22BudgetCrossingFrameOrder

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.719s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestBudgetStopControlHumanNamesDurableIntervention

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestBudgetStopControlHumanNamesDurableIntervention

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.096s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22AmendedFailureWireCorpus

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22AmendedFailureWireCorpus

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22AmendedFailureWireCorpus

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	10.604s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22ConformanceCorpusHasEveryAcceptedOperationAndFailure

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./internal/cli -run '^TestReasonCodeDetailSchemaIsExactForEveryCode

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./internal/cli -run '^TestReasonCodeDetailSchemaIsExactForEveryCode

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	0.021s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22ConformanceCorpusHasEveryAcceptedOperationAndFailure

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.507s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22(BudgetCrossingFrameOrder|CompleteOneIssueRecoveryJourney|PortableRecoveryConformance|HumanFixtures|AmendedFailureWireCorpus|ConformanceCorpusHasEveryAcceptedOperationAndFailure|FinalActivated)' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22(BudgetCrossingFrameOrder|CompleteOneIssueRecoveryJourney|PortableRecoveryConformance|HumanFixtures|AmendedFailureWireCorpus|ConformanceCorpusHasEveryAcceptedOperationAndFailure|FinalActivated)' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	30.124s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22PortableRecoveryConformance

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
 -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	12.863s` (terminated `exit_code: 0`)

Preserved first required-gate attempt (superseded after the documented inventory repair):

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`)

## Exact quality-gate results

GATE_ROWS

## Deviations from spec

- No semantic, schema, Event, authority, dependency or quality-gate deviation. DEC-205..207 are minimum execution allocations for accepted behavior missed by earlier tasks. Root implemented this task directly under the Human's delegate authority, with no concurrent product writer or subagent implementation.
- Required race attempt1 failed `TestPhase22AttentionReadBoundary`, terminal `FAIL` (terminated `exit_code: 1`): the exact production Attention-path inventory omitted DEC-206's newly required `cmd/fid/live_output.go` notice. Under the task's incorporated execution-scope.md, root added that one literal path while retaining the exact-set comparison and read-only assertions. Introducing row: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9. The unchanged ordinary race command executes it. This causal convergence changes no protected gate or boundary semantics; all first-attempt receipts remain preserved and the byte-identical gate was restarted from formatting.
- Initial intended inventory red: `amended operation claim.expire has no conformance fixture`, terminal `FAIL` / exit_code1. Literal operation/error/frame membership and executable fixture references repaired it; byte-identical focused command later passed.
- Portable first expansion failed `fid attention_count = 1, want 0`, terminal `FAIL` / exit_code1. Complete expectations now include the contract-required Budget Stop/Attention/link and the missing-artifact fixture's existing token baseline. No semantic dimension was deleted. The byte-identical portable command later passed.
- Journey setup failures are preserved individually above: missing Go type assertion; stale local Run version after committed denial; three runs diagnosing the same embedded Stop error; then a fake summary that counted predecessor Stops outside the current-Run query. These were fixture setup/model corrections except DEC-205's diagnosed `aggregate stop/STP-70ef70 sequence 2, want 1` product gap. The sequence now comes from canonical journal history, and explicit Stop controls plus embedded resolution/replay pass. Historical records remain in full archive parity.
- Human fixture first red was a child-process failure from a nonexistent `project get` metadata read, terminal `FAIL` / exit_code1. Metadata now comes from the existing typed project config; setup and inspection remain actual public processes. Same exact generator command later passed.
- Budget-frame fixture first incorrectly expected a start stream (the accepted stream grammar has no such stream); correction preserved the actual intended red `admitted attempt3 budget-stop frames = 0, want 1`, terminal `FAIL` / exit_code1. DEC-206 binds controls only to this execution's committed Attempt/Receipt crossing facts, retaining counters and normal terminal semantics. Human notice red was `[control] unknown`, terminal `FAIL` / exit_code1. Both byte-identical commands passed after the repair. A later replay check used `response.replayed` instead of accepted `response.result.replayed`; the combined focused command's failed receipt is retained and its corrected byte-identical rerun passed.
- Failure corpus initially supplied an empty disposition Event array, rejected before the intended lookup/version branch; it now supplies a nonempty missing Event ID. This exposed the three DEC-207 defects with otherwise valid requests. The intended failed corpus, terminal `FAIL` / exit_code1, is preserved separately from setup failure; the byte-identical corpus passed after typed mappings were repaired. The retained Run-start duplicate remains null-subject; recovery names its predecessor Run.
- Inventory parser initially counted the table header as a nineteenth error; registry test initially placed the added accepted code at the wrong ordinal. Both test-only mistakes have captured terminal `FAIL` / exit_code1 and later byte-identical PASS. No protected enum or expectation was weakened.

### Retained-row and transferred execution provenance

- `internal/workflow/stop.go`: original Stop producer/loader `ac827959f8359473cbd6546ac6bea1ff1fd20751`, phase22/task2. DEC-205 transfers only canonical aggregate sequence for embedded disposition; DEC-207 transfers only missing-Stop mapping. Original report remains COMPLETE. Ordinary race executes retained explicit Stop tests; tagged race executes embedded same-Issue disposition/replay and missing-Stop wire failure.
- `internal/evidence/execution.go` and budget effects: `a3766804f1a40ca5e827748f59550e1d72d452aa`, phase22/task7, introduced durable crossing Effects; task8 owned execution/closure/live convergence. DEC-206 transfers only committed budget-stop projection and consequential Human notice. Both ordinary and tagged exact gates execute matching controls and terminal regressions.
- `cmd/fid/run_recovery.go` / duplicate failure renderer: `7e419bfd8c233f933680d4a3b3957814c29dbd5e`, phase22/task3, inherited by task6 recovery. `cmd/fid/run.go` and adjacent error registry assertion: `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, phase22/task9, consumed task7 typed Budget sample conflict. DEC-207 preserves closed domains/nullability and maps the missing cases; ordinary exact race executes registry positive/negative schema tests and retained Run-start null case, tagged exact race executes real failures.
- `p22_fake_test.go`: `f3be616f0587881ee3f5c3716ffbf6167e3f2a20`, phase22/task6, introduced the fixture helper and default actions. Task18 adds final-location setup and optional explicit action set without changing default adapter cases. `p22_activation_recovery_test.go`: `30dd780e8d610962d304cf93ddab933a060e8f03`, phase22/task11, introduced the boundary helper; final-location setup retains previous wrapper behavior. All retained and new callers execute under the unchanged tagged race command.
- `p22_portable_conformance_test.go`, `p22_portable_fake_test.go`, all eight `portable/*.json` rows and their manifest: `5ec69e6ef31dd93229c3d42fc76a4c02d60084ee`, phase22/task15. Task18's accepted complete field inventory displaces partial expectations. Each row now pins all ten fields; current-Run filtering preserves the public query subject while full entry bytes retain prior history. Exact tagged race executes all eight cases, their shared fake and source/imported/rebuilt fid views.

## Concerns, technical debt, or notes for the Tech Lead

- RETAINED_MANIFEST
- FinalQA must independently audit the as-built kit mapping and all80 VAL dispositions. Only77 assertions are deterministic; Human076–078 remain blank. Production kit adapter qualification and accepted coordinator archive/terminalState differences remain Phase25 work under DEC-187/CF-037.
- The retained successor lease is24hours and its exact expiry is in the manifest. After expiry regenerate into a new fresh directory before judging an eligible successor view. Prior fixtures are never overwritten.
- Three accepted behaviors escaped earlier task evidence and were caught only by the connected/final corpus. Wrap-up should carry this concrete observation alongside DEC-194/195's allocation and executable acceptance lessons.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
