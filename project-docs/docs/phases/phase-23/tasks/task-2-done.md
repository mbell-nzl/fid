# Phase 23 - Task 2: Typed multi-Obligation Envelope contract - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: MF-01 QA repair and the DEC-249 (review-7 BC-1) conformance repair (attempt 1 of trigger `review-7:qa-rework:task-2-BC-1`) are complete; all Task 2 required gates pass.
- Recovery attempts: 2
- Next action: Phase Runner review.
- Human attention required: no

## Summary of what was implemented

- Envelope authoring now accepts 1..16 typed gate-evidence, review, and human-outcome Obligations and 0..16 protected Decisions.
- Evaluation returns ID-ordered results with distinct validity, freshness, and Artifact availability dimensions.
- The CLI and portable integration model serialize accepted null registrations and typed authoring rows.
- DEC-249 conformance repair (2026-09-26, review-7 BC-1): `issue.transition-terminal-success` and `envelope.complete` now evaluate and deny on only the Obligations that guard them; see "DEC-249 repair" below.
- This review-2 extension is based on prior repair commit `d71618aae4d8f3f56fee98468d2400a896a7a760`; its earlier matrix evidence and failed attempts remain recorded below.

## Files created or modified

- `internal/workflow/envelope.go`, `internal/workflow/evaluation.go`, and matching workflow tests.
- `internal/journal/current_event_catalog.go` and `cmd/fid/{envelope_test.go,task18_completion_resume_test.go,integration/p22_portable_schema_test.go}`.
- `project-docs/docs/phases/phase-23/tasks/task-2-done.md`.
- DEC-249 repair: `internal/workflow/evaluation.go` (`obligationsGuarding`, terminal-success call site), `internal/workflow/envelope.go` (`evaluateEnvelopeCompletion`), new `internal/workflow/obligation_guard_test.go`, new `cmd/fid/integration/p23_obligation_guard_test.go`, `cmd/fid/integration/p22_fake_test.go` (additive fixture option `ReviewGuardsTerminalSuccess`), and the converged `cmd/fid/integration/p23_human_outcome_test.go`.

## Documentation updated

- Completion report only; implementation conforms to the accepted Phase 23 contracts.

## Validation Contract Coverage

- VAL-006: `internal/workflow/envelope_test.go:TestEnvelopeAuthoringObligationCardinalityAndTypeFieldMatrix` proves 0/1/16/17 Obligation boundaries; `TestEnvelopeAuthoringAcceptedAndForbiddenTypedFieldCombinations` covers every accepted type and each forbidden registration, discriminator, evidence-type, trust-source, freshness, Artifact-policy, and independence combination; `TestEnvelopeAuthoringProtectedDecisionMatrix` and `cmd/fid/envelope_test.go:TestEnvelopeCreateCommandCoversObligationAndDecisionBoundaries` prove 0/16/17 protected-Decision boundaries; `TestEnvelopeCreateCommandRejectsEveryForbiddenTypedFieldCombinationWithoutPersistenceChange` sends all 24 forbidden typed-field cases through `envelope create`.
- VAL-007: `internal/workflow/envelope_test.go:TestEnvelopeAuthoringProtectedDecisionMatrix` covers duplicate class, duplicate ordinal, and unguarded action; `cmd/fid/envelope_test.go:TestEnvelopeCreateCommandRejectsEveryTypedAuthoringNegativeWithoutPersistenceChange` covers all five request-level negatives and every listed protected-Decision rejection, asserting unchanged journal and projection bytes through `snapshotWorkspacePersistence`.
- VAL-008: `internal/workflow/envelope_test.go:TestEvaluateEnvelopeObligationsReturnsIndependentOrderedResults` proves Obligation-ID order, distinct validity/freshness/Artifact dimensions, and all unsatisfied reasons; `TestEvaluateEnvelopeObligationsIncludesQualifyingValidFreshAvailableState` proves the qualifying state.
- VAL-009 (task-2 blocking half, both states per DEC-249; the earlier "activation evaluates every Obligation" citation here did not prove a guarded-action probe and is replaced). Applicable means `required` and `guarded_actions` contains the action.
  - In-package, `internal/workflow/obligation_guard_test.go:TestGuardedActionsEvaluateOnlyObligationsThatGuardThem`, run through the real `EvaluateExecutionTransition` and `evaluateEnvelopeCompletion` over a seeded Envelope: subtest `unsatisfied review guarding both actions blocks both` (satisfied `gate-evidence` plus unsatisfied `review` guarding both: both actions blocked, naming `OBL-review`); subtest `unsatisfied review guarding only completion leaves terminal success unblocked` (terminal-success reasons are empty, completion still names `OBL-review`); subtest `no Obligation guards terminal success` (two unsatisfied Obligations, terminal success fully permitted); subtest `unsatisfied gate guarding only terminal success does not block completion` (symmetric); subtest `every unsatisfied applicable Obligation is named`. Every subtest also asserts the read (`currentEnvelopeEvaluations`) reports both Obligations with their true statuses.
  - `TestCompleteEnvelopeDeniedByGuardingObligationCommitsNothing`: the real `CompleteEnvelope` denial names `OBL-review`, leaves the journal directory file sizes and `work_envelopes.status` unchanged (no Event, no projection change).
  - Built binary (tagged `cmd/fid/integration`): `p23_obligation_guard_test.go:TestP23ReviewObligationGuardingBothActionsBlocksBothRealProcess` (review guarding both actions and unsatisfied, gate satisfied by a real fid-executed Receipt: `envelope.complete` exits 4 `obligation-unsatisfied` naming exactly `OBL-p23-review`; the terminal `qa_pass` exits 4 `transition-not-permitted` with the single code `obligation-evidence-absent`, which can only come from the review because the gate is satisfied; the CLI terminal denial carries reason codes, not subject IDs; journal event count and projection unchanged for both) and `TestP23UnsatisfiedObligationBlocksOnlyTheActionItGuardsRealProcess` (review guarding only completion: `envelope.complete` denied naming `OBL-p23-review`; `qa_pass` permitted and appends exactly one Event; `envelope.complete` still denied naming exactly `OBL-p23-review` afterwards; `evidence.evaluate` reads both Obligations before and after). The permitted `qa_pass` uses a dedicated single-purpose fixture and is that test's last transition.
  - Retained: `TestEnvelopeActivationEvaluationIncludesEveryObligation` (activation reports every Obligation) and `cmd/fid/envelope_test.go:TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions`.
- VAL-010: `internal/workflow/envelope_test.go:TestEvaluateEnvelopeObligationsIncludesQualifyingValidFreshAvailableState` proves gate Evidence is `available` in the qualifying valid/fresh state and review/Human metadata-sufficient rows are exactly `not-applicable`; the existing independent-state test preserves `missing` for absent gate Evidence.

## Tests run

- Focused red: `go test -race -timeout 30m ./internal/workflow/ -run '^TestEnvelopeCreateAcceptsTypedObligationSet$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- Focused green: `go test -race -timeout 30m ./cmd/fid/ -run '^TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 6.856s` (terminated `exit_code: 0`)
- Portable fixture: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/ -run '^TestP22(CurrentEventFamiliesLiveReplayValidation|CurrentEventClosedValuesAndConditionalEffects|PortableFakeAndFIDRejectUnknownEffects)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 335.024s` (terminated `exit_code: 0`)
- Required workflow pre-edit: `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 192.569s coverage: 62.6% of statements` (terminated `exit_code: 0`)
- Required workflow post-edit, uncached observation: `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 193.138s coverage: 63.8% of statements` (terminated `exit_code: 0`)
- Real entrypoint coverage: the tagged `cmd/fid/integration` gate starts the `fid` CLI process and passed in `650.206s` (terminated `exit_code: 0`).
- QA repair red: `go test -race -timeout 30m ./internal/workflow/ -run '^(TestEnvelopeAuthoringObligationCardinalityAndTypeFieldMatrix|TestEnvelopeAuthoringAcceptedAndForbiddenTypedFieldCombinations|TestEnvelopeAuthoringProtectedDecisionMatrix|TestEvaluateEnvelopeObligationsIncludesQualifyingValidFreshAvailableState)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`; test fixture initially produced 17 rows for the 16-row case).
- QA repair green: the same byte-identical focused domain command - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 3.382s` (terminated `exit_code: 0`).
- QA repair command-surface red: `go test -race -timeout 30m ./cmd/fid/ -run '^(TestEnvelopeCreateCommandCoversObligationAndDecisionBoundaries|TestEnvelopeCreateCommandRejectsEveryTypedAuthoringNegativeWithoutPersistenceChange)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`; the unguarded-action fixture still left the action guarded).
- QA repair command-surface green: the same byte-identical focused command - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 7.110s` (terminated `exit_code: 0`).
- QA repair named matrix: `go test -race -timeout 30m ./internal/workflow/ -run '^(TestEnvelopeCreateAcceptsTypedObligationSet|TestEnvelopeCreateRejectsInvalidTypedObligationAndProtectedDecisionAuthoring|TestEnvelopeCreatePersistsEveryTypedObligationAndProtectedDecision|TestEvaluateEnvelopeObligationsReturnsIndependentOrderedResults|TestEnvelopeActivationEvaluationIncludesEveryObligation|TestEnvelopeAuthoringObligationCardinalityAndTypeFieldMatrix|TestEnvelopeAuthoringAcceptedAndForbiddenTypedFieldCombinations|TestEnvelopeAuthoringProtectedDecisionMatrix|TestEvaluateEnvelopeObligationsIncludesQualifyingValidFreshAvailableState)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 3.520s` (terminated `exit_code: 0`).
- QA repair named command surface: `go test -race -timeout 30m ./cmd/fid/ -run '^(TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions|TestEnvelopeCreateCommandCoversObligationAndDecisionBoundaries|TestEnvelopeCreateCommandRejectsEveryTypedAuthoringNegativeWithoutPersistenceChange)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 7.247s` (terminated `exit_code: 0`).
- QA review-2 remaining repair command surface: `go test -race -timeout 30m ./cmd/fid/ -run '^(TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions|TestEnvelopeCreateCommandCoversObligationAndDecisionBoundaries|TestEnvelopeCreateCommandRejectsEveryTypedAuthoringNegativeWithoutPersistenceChange|TestEnvelopeCreateCommandRejectsEveryForbiddenTypedFieldCombinationWithoutPersistenceChange)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 7.687s` (terminated `exit_code: 0`). The new case table covers all 24 forbidden typed-field combinations and asserts typed rejection, empty stdout, and byte-identical persistence.

- DEC-249 pre-edit conditional baseline, `go test -race -timeout 30m -cover ./internal/workflow/` at parent commit `d6542bfc093059c8c6e656fa25c26229e60408ee` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`); the same command with `-count=1` (supplementary, to obtain an uncached time) - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	184.168s	coverage: 62.1% of statements` (terminated `exit_code: 0`).
- DEC-249 first-failing red (product unchanged, new probes present): `go test -race -timeout 30m ./internal/workflow/ -run '^(TestGuardedActionsEvaluateOnlyObligationsThatGuardThem|TestCompleteEnvelopeDeniedByGuardingObligationCommitsNothing)$' -count=1` - FAIL - `FAIL	github.com/mbell-nzl/fid/internal/workflow	5.712s` (terminated `exit_code: 1`); three subtests failed with the terminal-success denial naming `OBL-review` (and `OBL-gate`) although neither guards it, e.g. `terminal success blocked by Obligations [OBL-review], want []`.
- DEC-249 built-binary red (product reverted with `git checkout` of the two files, new and converged tagged tests present): `go test -tags integration -race -timeout 30m ./cmd/fid/integration/ -run '^(TestP23ReviewObligationGuardingBothActionsBlocksBothRealProcess|TestP23UnsatisfiedObligationBlocksOnlyTheActionItGuardsRealProcess|TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess)$' -count=1` - FAIL - `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration	12.136s` (terminated `exit_code: 1`): the journey at `p23_human_outcome_test.go:115` (`qa_pass denial reasons = [human-outcome-absent obligation-evidence-absent obligation-evidence-absent], want exactly [obligation-evidence-absent]`) and the contrast test (`qa_pass` denied `["obligation-evidence-absent"]`). The both-guarded test passed pre-change, as expected: the global evaluator also blocks that state; it is the blocking half and cannot be red.
- DEC-249 green, same in-package command after the fix - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	7.282s` (terminated `exit_code: 0`); same tagged command after the fix, with `-v` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	20.762s` (terminated `exit_code: 0`).
- DEC-249 reverted mutation checks (each edit reverted, file restored byte-identically from a saved copy, green re-confirmed `ok  	github.com/mbell-nzl/fid/internal/workflow	7.287s`): (1) completion filter flipped back to global (`if true || guarding[item.ObligationID]`) - FAIL, subtest `unsatisfied gate guarding only terminal success does not block completion`: `envelope.complete blocked by Obligations [OBL-gate OBL-review], want [OBL-review]`; (2) terminal-success call site flipped back to `envelope.Obligations` - FAIL, three subtests (`...guarding only completion leaves terminal success unblocked`, `no Obligation guards terminal success`, `unsatisfied gate guarding only terminal success...`).
- Real entrypoint for DEC-249: the two new tagged tests and the converged journey drive the built `fid` binary (`issue transition`, `envelope complete`, `evidence evaluate`, `evidence link`); observed results are in the DEC-249 green line above and in the tagged gate below.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 63.8% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.9%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 650.206s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 63.8% of statements` (terminated `exit_code: 0`)
- Observed wall times: `cmd/fid` `440.185s`; `cmd/fid/integration` `650.206s`; `internal/workflow` `193.138s` for the current-session post-edit uncached required run (the final rerun was cached).
- QA repair attempt 1 exact rerun: `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- QA repair attempt 1 exact rerun: `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- QA repair attempt 1 exact rerun: `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 203.449s coverage: 64.1% of statements` (terminated `exit_code: 0`)
- QA repair attempt 1 exact rerun: `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.9%` (terminated `exit_code: 0`)
- QA repair attempt 1 exact rerun: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 681.075s` (terminated `exit_code: 0`)
- QA repair attempt 1 exact rerun: `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 198.163s coverage: 64.1% of statements` (terminated `exit_code: 0`)
- QA repair attempt 1 wall times: `cmd/fid` `451.009s`; `cmd/fid/integration` `681.075s`; `internal/workflow` `198.163s` for the final dedicated workflow run (`203.449s` in the full race/coverage gate).
- QA review-2 remaining repair exact rerun: `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- QA review-2 remaining repair exact rerun: `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- QA review-2 remaining repair exact rerun: `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 64.1% of statements` (terminated `exit_code: 0`)
- QA review-2 remaining repair exact rerun: `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.9%` (terminated `exit_code: 0`)
- QA review-2 remaining repair exact rerun: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 666.085s` (terminated `exit_code: 0`)
- QA review-2 remaining repair exact rerun: `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 64.1% of statements` (terminated `exit_code: 0`)
- QA review-2 wall times: `cmd/fid` `448.937s`; `cmd/fid/integration` `666.085s`; `internal/workflow` retained noncached package wall time `198.163s` from the prior repair gate, with the current exact dedicated rerun cached and completing in `4.131s` process wall time.

### DEC-249 repair gate results (run serially against implementation commit `613e45aad9b10445cffd6b113ee598ab493f183a`, byte-identical spec commands, no `-count=1`)

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	209.816s	coverage: 64.0% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.0%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	691.913s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	188.151s	coverage: 64.0% of statements` (terminated `exit_code: 0`); uncached because the source changed
- DEC-249 DEC-227 wall times (all uncached, from the runs above): `cmd/fid/integration` `691.913s` (budget 900 s); `cmd/fid` `490.783s` (repository-wide race/coverage run; budget 600 s); `internal/workflow` `209.816s` in the repository-wide run and `188.151s` in the dedicated run (budget 600 s). No budget breach; the new tagged tests take about 3.7 s each.
- Gate 4 total statement coverage is `62.0%`; there is no threshold. The `65.9%` recorded by earlier repair sections was measured before later Phase 23 tasks and is not comparable.

## DEC-249 repair (2026-09-26, ESC-23-004 BC-1, attempt 1)

Authority: DEC-249 (accepted conformance ruling); order step 2 of DEC-255/DEC-256. No contract, schema, Event, gate, or public-interface change; `run.close`, readiness raising, reason codes, precedence, denial shapes and Stop/protected-Decision behavior are unchanged.

- Product change: `obligationsGuarding(obligations, action)` (`internal/workflow/evaluation.go`) keeps a `required` Obligation whose `guarded_actions` contains the action. `evaluateExecutionTransition` passes only those to `evaluateObligations` for `issue.transition-terminal-success`; `evaluateEnvelopeCompletion` (`internal/workflow/envelope.go`) keeps only the reasons of Obligations guarding `envelope.complete` while still returning every Obligation's evaluation. `run.close` still passes `envelope.Obligations` (`run_close.go:81`, unedited). The Human-Obligation readiness condition still reads every other Obligation of the revision (unedited).
- Fixture change: `cmd/fid/integration/p22_fake_test.go` gains the additive option `ReviewGuardsTerminalSuccess`, used only with `CoexistingReviewObligation`, because no existing option let a `review` Obligation guard both actions and the built-binary probe needs that state. It changes no existing fixture's behavior.

### Converged retained assertions (DEC-162 provenance)

Invalidating change for every row: this repair, implementation commit `613e45aad9b10445cffd6b113ee598ab493f183a` (DEC-249 guard filter). Introducing commits from `git log -S` in `cmd/fid/integration/p23_human_outcome_test.go`; the file was created by `0ff2eee` (2026-09-25, `feat(review): record provenanced Human Outcomes with readiness Stops`), and the rows below were introduced within Phase 23. Executing command for every row: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` (focused: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/ -run '^TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess$' -count=1`).

| Row | Introducing commit | Before | After (subject preserved) |
|---|---|---|---|
| Journey stage 1 `qa_pass` denial (`assertP23GuardedActionsDenied ... "human-outcome-stage-1"`) | `af0cbb2` (2026-09-26, `test(integration): prove VAL-009, VAL-025, and VAL-043 with the complete real-producer fixture`) | terminal success denied for `[OBL-human-security, OBL-p22-corpus, OBL-p23-review]` | `envelope.complete` still names all three; terminal success is denied for `[OBL-p22-corpus]`, the only unsatisfied Obligation that guards it; still no Event, no projection change |
| Journey stage 2 `qa_pass` denial (`... "human-outcome-stage-2"`) | `af0cbb2` | terminal success denied for `[OBL-human-security, OBL-p22-corpus]` | `envelope.complete` names `[OBL-human-security, OBL-p22-corpus]`; terminal success denied for `[OBL-p22-corpus]` |
| Journey stage 3 `qa_pass` denial (`assertP23TerminalTransitionDenied ... "human-outcome-stage-3"`) | `af0cbb2` | `qa_pass` denied for `[OBL-human-security]` | `OBL-human-security` guards only `envelope.complete`, and the existing `human-outcome-required` denial just above still proves it blocks that action. New `assertP23OnlyObligationUnsatisfied` reads `evidence evaluate` (every Obligation reported, only `OBL-human-security` unsatisfied). New `assertP23TerminalTransitionNotBlocked` asserts, with the non-committing domain evaluator `workflow.EvaluateExecutionTransition` (not a permitted `qa_pass`), that no reason blocks the terminal transition, and that journal and projection are unchanged. The one-way `qa_pass` is not performed early; the journey's real terminal `qa_pass` is unchanged |
| Header comment ("evaluates every required Obligation regardless of its own guarded_actions") | `0ff2eee` (comment added in the file's creation commit) | described the global reading | rewritten to the DEC-249 reading |
| Helper `assertP23GuardedActionsDenied` / doc comments | `af0cbb2` | one `unsatisfied` list for both actions | takes `completeUnsatisfied` and `terminalUnsatisfied` lists; the stage rows above pass them; nothing deleted or weakened |

Searches for other displaced expectations: `rg -n "qa_pass|TerminalSuccess|terminal-success" cmd internal --glob '*_test.go'` and `rg -n "ReviewObligation|CoexistingReviewObligation|HumanObligations" cmd --glob '*_test.go'` found only the rows above; every other terminal-success test uses a sole `gate-evidence` Obligation guarding both actions, or a protected Decision, and passed unchanged. The full ordinary and tagged gates then revealed no further displaced assertion.

The deeper MF-5 restructuring (a new gate-satisfied/review-unsatisfied stage, the report citation refresh and ordering of the real terminal `qa_pass`) is task-10's later step; `task-10-done.md` is not edited here.

## Deviations from spec

- First full race attempt failed on the displaced single-map completion assertion; its exact token was `FAIL` (terminated `exit_code: 1`). Task 2 converged that assertion to require one list item for its legacy fixture, then the byte-identical gate passed.
- First integration attempt failed on accepted new portable columns; its exact token was `FAIL` (terminated `exit_code: 1`). The Task 2 integration fixture now recognizes the accepted schema columns, and the byte-identical gate passed.
- QA repair attempt 1 preserved two focused red attempts: the new cardinality fixture initially returned 17 instead of 16, and the command fixture initially left the supposedly unguarded action guarded. Both test-only fixture defects were corrected within the bounded review scope; their subsequent byte-identical focused commands passed.
- DEC-249 repair: no required command failed; each ran once and passed. The two DEC-249 focused reds above are intentional first-failing evidence, not gate failures. The pre-edit baseline printed `(cached)`, so an uncached `-count=1` time was also recorded (supplementary, not a substitute). The new test files were added to the working tree only after the baseline test binary had compiled, so the baseline reflects the unmodified package.
- DEC-249 repair: `cmd/fid/integration/p22_fake_test.go` was edited additively (one fixture option) although it is a retained fixture and was not itself invalidated by the change; recorded here for the Tech Lead. It changes no existing behavior.

## Concerns, technical debt, or notes for the Tech Lead

- The built-binary terminal-success denial (`transition-not-permitted`) exposes reason codes only, so the CLI cannot name the blocking Obligation ID for `issue.transition-terminal-success`; the in-package test names it. This is existing behavior and out of scope.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `613e45aad9b10445cffd6b113ee598ab493f183a`

## Commit subject

- `fix(workflow): evaluate only guarding Obligations per guarded action`
