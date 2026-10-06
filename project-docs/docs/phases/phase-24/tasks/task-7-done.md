# Phase 24 - Task 7: Run leaves - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

OUTCOME COMPLETE

## Summary of what was implemented

- Start, close and abandon accept required expected_repository and record fid's observation under the write lock after keyed replay. Mismatches have the amended closed repository-mismatch details; unsafe terminal observations retain recovery-state-conflict.
- Start checks Claim then Envelope through the accepted same-worktree predicate. Commits and edits are accepted, genuine import destinations are accepted, and plain copies keep the exact fixed fallback and domain reason run claim binding is invalid, without an Event.
- Close/abandon record terminal_repository with base_oid equal to observed HEAD. Close evaluates Obligations against that observation, retaining Receipt-backed success and stale-Evidence failure.
- Budget measurement_ref is nullable SourceInput. Non-null measurements observe, resolve/check project-file digests, and store the observed Source/Event repository. Null measurements do not observe. Equal keyed retries replay before HEAD or Source reads.
- Adopted all seven Run descriptors, generated field help and typed editable examples, recursive structural rejection, exclusivity, removed-member rejection, and fake/fid convergence. Added process tests for unsafe repositories, replay, deletion and viable new-key controls.
- Converged retained Run callers and assertions with introducing-commit provenance below. Gate-run and Evidence product behavior is outside this implementation; their fixtures use the repository actually recorded by the new Run.

## Files created or modified

- `cmd/fid/gate_run_test.go`
- `cmd/fid/generic_run_test.go`
- `cmd/fid/input_descriptor_test.go`
- `cmd/fid/input_descriptors_run.go`
- `cmd/fid/integration/p21_conformance_test.go`
- `cmd/fid/integration/p22_abandon_authority_test.go`
- `cmd/fid/integration/p22_close_proof_test.go`
- `cmd/fid/integration/p22_fake_test.go`
- `cmd/fid/integration/p22_recovery_adopted_target_test.go`
- `cmd/fid/integration/p22_recovery_journey_test.go`
- `cmd/fid/integration/p22_resume_actions_test.go`
- `cmd/fid/integration/p22_run_replay_snapshot_test.go`
- `cmd/fid/integration/p22_semantic_fake_test.go`
- `cmd/fid/integration/p23_claim_mode_test.go`
- `cmd/fid/integration/p23_failure_contract_test.go`
- `cmd/fid/integration/p23_gate_grants_test.go`
- `cmd/fid/integration/p23_human_failure_test.go`
- `cmd/fid/integration/p23_review_evidence_budget_test.go`
- `cmd/fid/integration/p23_review_record_test.go`
- `cmd/fid/integration/p24_flag_corpus_test.go`
- `cmd/fid/integration/p24_idempotency_test.go`
- `cmd/fid/integration/p24_run_inputs_test.go`
- `cmd/fid/p24_run_test.go`
- `cmd/fid/review_test.go`
- `cmd/fid/run.go`
- `cmd/fid/task12_execution_transition_test.go`
- `internal/export/archive_fixture_test.go`
- `internal/export/import_test.go`
- `internal/workflow/recovery_test.go`
- `internal/workflow/resume.go`
- `internal/workflow/run.go`
- `internal/workflow/run_budget.go`
- `internal/workflow/run_budget_test.go`
- `internal/workflow/run_close.go`
- `internal/workflow/run_close_authority_test.go`
- `internal/workflow/run_close_test.go`
- `internal/workflow/stop_conditions_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md`
- `project-docs/docs/phases/phase-24/tasks/task-7-done.md`

## Documentation updated

- project-docs/specs/06_DESIGN_PLAN_CLI.md: current Run leaves, expectation inputs, observation and digest/replay rules, same-worktree/import semantics and authoring surface.
- project-docs/specs/07_DESIGN_PLAN_FLOWS.md: review/rework Run admission after changes and terminal Obligation evaluation.
- This completion report records evidence, convergence provenance and failed attempts. Runner-owned run-ledger.md and triage-log.md were neither edited nor included in either task commit.

## Validation Contract Coverage

Scope is the Task 7 Run slice of the phase assertions; whole-phase completion is not claimed. Named tests are covered by this session's final ordinary or tagged command results. Unchanged packages marked (cached) in the output are identified as such; touched packages have fresh results.

- VAL-009: PASS - TestP24RunReplayBeforeObservationAndSourceReads (cmd/fid/integration/p24_run_inputs_test.go), all four leaves / HEAD-and-file: retry after a commit and Source edit returns the original result/EventRef except replayed:true, with unchanged projection and journal.
- VAL-011: PASS - TestInputDescriptorRegistryMatchesLiveTreeAndTypes (cmd/fid/input_descriptor_test.go) checks the live tree versus registry/fence and rejects a seeded type mismatch; TestP24RunDescriptorsAdoptSevenLeaves checks all seven registrations and types.
- VAL-012: PASS - TestP24RunDescriptorProcessCorpus checks every field recursively for path/type/requirement/meaning/constraints and the template command, including expected_repository git command meanings.
- VAL-013: PASS - TestP24RunDescriptorProcessCorpus runs all seven --example leaves through the built binary and compares stdout byte-for-byte with the ordered, two-space InputExample template.
- VAL-014: PASS - TestP24RunDescriptorProcessCorpus decodes each raw example into its actual reflected request type with unknown fields disallowed; TestP24JudgmentDescriptorSourcePositionsAndExamples retains strict typed decoding and validates expanded format placeholders.
- VAL-015: PASS - TestAdoptedExampleProcessExclusivity covers six prohibited combinations for each of the 27 adopted leaves (162 combinations, including all seven Runs). Example-only success is separately asserted by TestP24RunDescriptorProcessCorpus.
- VAL-016: PASS - TestP24RunDescriptorProcessCorpus derives unknown/duplicate/missing/null/type/enum/format/bounds rows at every applicable descriptor level, including authority.grant_ref children. It checks exact failure details and no directory effects.
- VAL-017: PASS - TestP24RemovedRunRepositoryMembers rejects repository on start/close/abandon and measurement_ref.repository on observe-budget, with exact paths and no writes. TestP24RunChangedFakeFIDInputParity checks the shared changed-leaf intersection.
- VAL-018: PASS - TestP24RunBudgetSourceInputMatrix/matching stores the supplied file digest and independent repository; /mismatching requires both digests in closed source-digest-mismatch details and unchanged durable state.
- VAL-019: PASS - TestP24RunBudgetSourceInputMatrix/omitted stores SHA-256 of real bytes; /missing requires expected_digest:null and observed_digest:null with no Event. TestP24BudgetObservesOmittedSourceDigest (internal/workflow/run_budget_test.go) also preserves the stored digest when the file is deleted before replay.
- VAL-022: PASS - TestP24RunStartAfterCommitAndEdit independently checks Git HEAD/tree/status, root hash and start/Event observation; TestP24RunTerminalExpectationObservationAndUnsafe/match independently checks both terminal records and Events with base_oid=HEAD; TestP24RunBudgetSourceInputMatrix/matching and /omitted compare Source/Event to the independent Git/config oracle.
- VAL-023: PASS - TestP24RunReplayBeforeObservationAndSourceReads/start proves initial acceptance and fresh unmerged/non-Git control refusals; TestP24RunBudgetSourceInputMatrix compares matching/omitted acceptance with unsafe/missing-Git refusal and checks mapped reason/detail, no Event. Both null unsafe cases succeed. Shared classifier cause fixtures in internal/project/repository_reason_test.go are included in the unchanged internal/project result (cached) reported by the ordinary gate.
- VAL-024: PASS - TestP24RunStartExpectationMatrix checks match/HEAD/dirty; TestP24RunTerminalExpectationObservationAndUnsafe checks the same three states separately for close and abandon, preserving all durable state on mismatch.
- VAL-025: PASS - The two expectation matrices compare the entire details object to exactly {expected:{head_oid,dirty},observed:Repository}, echoing input and independently observed repository. The retained failure-producer registry includes all three Run producers.
- VAL-029: PASS - TestP24RunStartOriginalCopyAndImport uses a valid active Claim snapshot: original and genuine export/import destinations succeed and record the local observation; a whole-worktree plain copy returns the exact fallback, exit 2, and no durable changes. The direct domain call requires run claim binding is invalid, preserving Claim-before-Envelope.
- VAL-030: PASS - TestP24RunStartAfterCommitAndEdit creates Envelope/Claim, commits a new tracked file, edits it and starts successfully with the full current observation rather than the Envelope snapshot.
- VAL-031: PASS - TestP24RunTerminalExpectationObservationAndUnsafe/close/match and /abandon/match commit and edit after start and compare every terminal/Event repository field with the independent observation and base_oid=HEAD.
- VAL-035: PASS - TestP24RunChangedFakeFIDInputParity derives close/abandon/observe-budget from the shared dispatch intersected with changed leaves; both adapters accept converged inputs and reject removed members with exit 2. TestP22DualAdapterRecoveryConformance retains the full semantic corpus. Start is outside that fake dispatch and has independent real CLI coverage.
- VAL-052: PASS - Final ordinary and tagged package times are recorded below. Every ordinary package is below 600 s and tagged integration below 900 s; the 480/720 s warning thresholds were checked.
- VAL-060: PASS - TestP24RunDescriptorProcessCorpus executes all seven examples in fresh empty non-project directories, requires exit 0/empty stderr and verifies directories remain empty. Explicit CLI probes below independently verify the same surface.
- VAL-061: PASS - TestP24RunTerminalExpectationObservationAndUnsafe pairs both terminal leaves with HEAD/dirty mismatch -> repository-mismatch and unmerged/non-Git -> recovery-state-conflict, exit 4 and no projection/journal change; the retained Phase 22 failure contract checker verifies typed metadata/details.
- VAL-062: PASS - TestP22CloseRunUsesObservedRepositoryAndQualifyingProof/qualifying closes a Receipt-backed Obligation as succeeded; the separate /stale fixture edits tracked base.txt and closes failed against the actual terminal observation.
- VAL-066: PASS - TestP24RunReplayBeforeObservationAndSourceReads covers all four observing leaves in unmerged and non-Git worktrees, plus budget deleted-source. Retries return original results/no Event; distinct controls have fresh Claim/Run or sample-chain prerequisites and fail with repository-unavailable, recovery-state-conflict, or source-digest-mismatch without an Event.

## Tests run

Working directory: /home/mark/projects/fid unless stated otherwise. RED attempts remain failed evidence. Every recorded result was watched through the nested command's observed termination.

- `go test ./cmd/fid -run '^TestP24RunStartAfterCommitAndEdit$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestP24RunCloseAfterCommitAndEdit$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./internal/workflow -run '^TestP24BudgetObservesOmittedSourceDigest$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./internal/workflow -run '^TestP24AbandonRecordsObservedBase$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./internal/workflow ./cmd/fid -run '^TestP24(Run(Start|Close)AfterCommitAndEdit|AbandonRecordsObservedBase|BudgetObservesOmittedSourceDigest)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.006s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestP24RunDescriptorsAdoptSevenLeaves$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestP24RunDescriptorsAdoptSevenLeaves$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.156s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestP24RunDescriptorProcessCorpus$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestP24RunDescriptorProcessCorpus$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.484s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP22DualAdapterRecoveryConformance$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid ./internal/workflow ./internal/export -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid ./internal/export -run '^(TestRunStartBindingMatrixRejectsWithoutMutation|TestEvidenceEvaluationMultiFailureGoldenSurvivesRepeatAndRebuild|TestExecutionBoundTransitionPredicateMatrix|TestP24Run.*|TestImport.*)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^(TestRunStartBindingMatrixRejectsWithoutMutation|TestEvidenceEvaluationMultiFailureGoldenSurvivesRepeatAndRebuild|TestExecutionBoundTransitionPredicateMatrix)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.844s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestP24Run(StartOriginalCopyAndImport|StartExpectationMatrix)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.232s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration -run '^(TestP24RunTerminalExpectationObservationAndUnsafe|TestP24RunBudgetSourceInputMatrix|TestP22DualAdapterRecoveryConformance|TestP22CloseRunUsesObservedRepositoryAndQualifyingProof)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP24RunReplayBeforeObservationAndSourceReads$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration -run '^(TestP24Run.*|TestP22DualAdapterRecoveryConformance|TestP22CloseRunUsesObservedRepositoryAndQualifyingProof)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP24RunChangedFakeFIDInputParity$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.796s` (terminated `exit_code: 0`)
- `go test ./cmd/fid ./internal/workflow ./internal/export -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestGateRunPublishesSanitizedArtifactByDigest$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.829s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/workflow ./internal/export` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestP24JudgmentDescriptorSourcePositionsAndExamples$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestP24JudgmentDescriptorSourcePositionsAndExamples$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.162s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP22AbandonUsesExactHumanForceAndTypedObservationFailures/(ordinary|force)/forged-target$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP22AbandonUsesExactHumanForceAndTypedObservationFailures/(ordinary|force)/forged-target$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.567s` (terminated `exit_code: 0`)
- `golangci-lint run ./internal/workflow` - FAIL - `* unused: 1` (terminated `exit_code: 1`)
- `golangci-lint run ./internal/workflow` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test ./internal/workflow -run '^Test.*Budget' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	1.237s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/workflow ./internal/export` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	128.636s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	333.953s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	57.949s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24RunBudgetSourceInputMatrix$' -count=1 ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.913s` (terminated `exit_code: 0`)

First-start test construction also failed once on the missing runGitForTest helper (exit_code:1, token FAIL); its next behavioral RED rejected expected_repository with invalid Run request. A subsequent test comparison failed on pointer versus decoded base_oid representation (exit_code:1, token FAIL); normalizing the comparison passed the same start command with token ok  	github.com/mbell-nzl/fid/cmd/fid	0.759s (exit_code:0). These test-construction failures preceded the combined green selection above.

Real entrypoint, independently of go test: scratch binary /tmp/tmp.wF590Cuq4E/fid, working directory /tmp/fid-task7-cli-QIoFOx.

- `/tmp/tmp.wF590Cuq4E/fid run start --example` - PASS - `}` (terminated `exit_code: 0`)
- `/tmp/tmp.wF590Cuq4E/fid run observe-budget --example` - PASS - `}` (terminated `exit_code: 0`)
- `ls -A /tmp/fid-task7-cli-QIoFOx` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- Direct --json run start --actor ACT-probe --input - --idempotency-key task7-entrypoint with the expanded template plus repository:{} returned observed CLI exit_code:2, invalid-argument and input.repository, empty stdout and unchanged directory. The Python probe asserted those observations and terminated exit_code:0 with token `PASS real CLI removed repository: exit_code: 2; input.repository; directory unchanged`.

Inventory regeneration (accepted input-inventory.md syntax; Markdown-escaped alternation rendered as regex alternation):

- `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid` - PASS - `cmd/fid/gate.go:272:func decodeClosedGateInput(stdin io.Reader, path string, target any) error {` (terminated `exit_code: 0`)
- `rg -n "func (same[A-Za-z]*Repository|validateRepositoryBinding|validEnvelopeRepository|validEnvelopeSource|validEvidenceRepository|validateRunTerminalRepository|verifyDecisionSource|importedGateDestination)" internal -g '!*_test.go'` - PASS - `internal/workflow/envelope.go:759:func validEnvelopeSource(source Source) bool {` (terminated `exit_code: 0`)
- `rg -n "Request\.Repository|request\.Repository|Request\.TargetRepository|request\.RepositoryScope|Request\.WorktreeID|request\.WorktreeID" internal -g '!*_test.go'` - PASS - `internal/evidence/execution.go:768:	event := receiptEvent(binding.actor, receipt, string(repositoryJSON), string(toolchainJSON), string(rejectedJSON), options.Request.Repository, effects)` (terminated `exit_code: 0`)
- Executed the exact section 1 recursive inventory block through bash /tmp/fid-task7-inventory-3am4bsih/walk.sh: PASS, token `Inventory artifacts: /tmp/tmp.wF590Cuq4E` (terminated `exit_code: 0`). Counts: leaves 92, groups 29, input 39, keyed 50, no-short 0, typed-flags 49. The 27 adopted example flags account for the increase from 22 typed flags. Run no longer has a closed family decoder or caller-repository writer; terminal validation now observes/checks the expectation. Remaining decoder/comparison/writer rows belong to unchanged phase surfaces, including private resolved Envelope/Registration facts.
- Applied gofmt -w and goimports -w to all edited Go files before final gates; both terminated exit_code:0 with exit 0 (no output).
- `git diff --check` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).

## Exact quality-gate results

Attempt 1 stopped for in-scope causal repair after lint; its failed result is retained:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - FAIL - `* unused: 1` (terminated `exit_code: 1`)

Exact task-spec and task-applicable QA-plan gate order, final tree / attempt 2:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	121.978s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					61.1%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	344.128s` (terminated `exit_code: 0`)

Reported final gate times: cmd/fid **342.375 s**; internal/workflow **121.978 s**; slowest other ordinary package with a fresh reported time, internal/issues **177.502 s** (internal/export 171.061 s); tagged cmd/fid/integration **344.128 s**. All measured ordinary packages are below 600 s and tagged integration below 900 s. None reaches the 480 s ordinary or 720 s tagged warning threshold. Unchanged ordinary packages marked (cached) in the captured output have no new wall time; no historical time is substituted. Overall statement coverage: **61.1%**.

## Deviations from spec

- Tagged race attempt 1 ended token FAIL (terminated exit_code:1; 305.005 s): ordinary/force forged-target abandon rows still expected recovery-state-conflict, although DEC-273 / Task 7 expressly maps an expectation mismatch to repository-mismatch. Introducing commit c0003e0c5df7 (phase-22/task-8) was verified with git log -S. The focused reproduction failed with token FAIL (exit_code:1); converged only that expected code, retaining both authority modes, exit 4, typed details and no-state-change checks. The identical focused command passed, and the byte-identical full tagged race tier rerun passed.
- Required lint attempt 1 failed with verbatim token `* unused: 1` (terminated `exit_code: 1`): SourceInput conversion removed the last caller-side validEnvelopeSource use. Focused package lint reproduced the same token/exit. The existing stored-Source invariant (introduced by `56dad8fa0dcb`, phase-21/task-6) now checks the resolved measurement after digest resolution and observation, in the allowlisted run_budget.go. Its type, label, digest and repository checks are preserved; caller inputs and accepted failure contracts are unchanged. No out-of-scope definition was edited. Focused lint returned `0 issues.` (exit_code:0), budget tests and post-repair race checks passed, and the byte-identical required lint returned `0 issues.` (exit_code:0) in the complete restarted ordered gate.
- No accepted contract, gate, budget, dependency or product-scope deviation. The first failed attempts remain immutable evidence below and in Tests run; the final exact five-command gate passes.
- Expected contract REDs: start/close rejected the new expectation, budget rejected digest omission and abandon retained caller base_oid. Implemented observation/checking and separate resolved facts, then the combined behavioral selection passed.
- Descriptor REDs: seven missing registrations, followed by absent command help installation. Registered descriptors and installed them on Run mutation leaves; the identical focused commands passed.
- Retained convergence failures (token FAIL, terminated exit_code:1): the first broad ordinary run exposed old request fields and an export package identifier shadow. Converged in-scope Run fixtures and used the existing project adapter with an unambiguous alias. The next focused binding run exposed a repository-mismatch renderer passing structs to a map-validated failure shape, resulting in empty stderr; emitted the accepted closed maps and the focused binding selection passed.
- Tagged focused failures (each token FAIL, terminated exit_code:1) exposed the retained failure-producer inventory omitting the newly accepted Run producers. Added run.observe-budget to source mismatch and start/close/abandon to repository mismatch, preserving every shape check. The replay selection reached that same inventory failure. The next combined selection exposed a shallow-copy parity fixture corrupting the accepted Source through its rejection input; deep-copying that test input made the parity selection pass. Both semantic families then passed in the full tagged race tier and final exact integration gate.
- The second broad ordinary attempt (token FAIL, terminated exit_code:1) exposed an Artifact test starting a new Run from a retained Envelope snapshot. Updated only the retained activeRunIDForGateTest fixture to use the actual new Run observation, preserving the Artifact sanitization/digest assertion and existing-Run fixture snapshots. Its focused test passed and the full touched ordinary race tier subsequently passed.
- Touched ordinary race attempt 1 ended token FAIL (terminated exit_code:1; cmd/fid 294.390 s, workflow 95.969 s, export 125.817 s). TestP24JudgmentDescriptorSourcePositionsAndExamples/run_observe-budget treated the intentionally editable timestamp placeholder as a concrete timestamp. The focused reproduction failed with token FAIL (exit_code:1). Kept strict decoding of each raw example into its actual request type, expanded declared format placeholders for structural validation, and preserved all fields/constraints/no-effects assertions. The focused command passed, followed by the byte-identical tier rerun (cmd/fid 300.717 s, workflow 93.827 s, export 128.636 s).

Retained assertion provenance was checked using introducing-commit git log -S and blame of displaced lines at dispatch HEAD. The accepted invalidating change is Phase 24 Task 7 / DEC-273 (expected_repository, observed terminal HEAD, SourceInput and Run descriptors); DEC-278 defines null/replay behavior and DEC-282 preserves Claim-first refusal. Each row retains its subject and rejection/no-mutation strength. Execution is by the final ordinary gate (O) or exact tagged integration gate (I), whose commands/results are above. Shared fixture rows converge their callers rather than assigning ownership by file/task labels.

| Retained assertion or fixture | Introducing commit / phase/task | Convergence and executing gate |
|---|---|---|
| `cmd/fid/generic_run_test.go`: TestRunStartBindingMatrixRejectsWithoutMutation | `8c763e9cddee` phase-21/task-11 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; O |
| `cmd/fid/generic_run_test.go`: runStartInput | `8c763e9cddee` phase-21/task-11; displaced row 7e419bfd8c23 phase-22/task-3 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; O |
| `cmd/fid/generic_run_test.go`: closeRunForTest | `a0a93c7ce4dd` phase-22/task-9 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; O |
| `cmd/fid/input_descriptor_test.go`: TestAdoptedExampleProcessExclusivity | `545803dc2185` phase-24/task-3; displaced row 1baaecc79838 phase-24/task-6 | Seven-leaf adoption / expanded declared placeholders; existing metadata and exclusivity assertions retained; O |
| `cmd/fid/input_descriptor_test.go`: p24ExpandDescriptorObjects | `1baaecc79838` phase-24/task-6 | Seven-leaf adoption / expanded declared placeholders; existing metadata and exclusivity assertions retained; O |
| `cmd/fid/integration/p21_conformance_test.go`: TestPhase21RealProcessJourneyStartsOneActivatedEnvelope | `236540c48a21` phase-21/task-22; displaced row 7e419bfd8c23 phase-22/task-3 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; I |
| `cmd/fid/integration/p22_abandon_authority_test.go`: TestP22AbandonUsesExactHumanForceAndTypedObservationFailures | `c0003e0c5df7` phase-22/task-8 | Expected repository; forged-target now repository-mismatch in ordinary/force; no-mutation and authority assertions retained; I |
| `cmd/fid/integration/p22_close_proof_test.go`: TestP22CloseRunUsesObservedRepositoryAndQualifyingProof | `c0003e0c5df7` phase-22/task-8 | Caller base rows now reject removed repository; forged old target mismatches; qualifying/stale proof outcomes retained; I |
| `cmd/fid/integration/p22_fake_test.go`: p22NewFIDRecoveryAdapterWithClaimModeAt | `35944c72aba0` phase-23/task-3; displaced row f3be616f0587 phase-22/task-6 | Converged real adapter request builder and SourceInput; I |
| `cmd/fid/integration/p22_fake_test.go`: adapter p22BuildArgs | `25caa7d40712` phase-22/task-9 | Converged real adapter request builder and SourceInput; I |
| `cmd/fid/integration/p22_recovery_adopted_target_test.go`: TestP22RecoveryAdoptedTargetPermitsTerminalTransitionWithFreshSuccessorProof | `c2d80980b9bf` phase-22/task-12 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; I |
| `cmd/fid/integration/p22_recovery_journey_test.go`: p22RecoveryFreshnessJourney | `2339dc9ebd8d` phase-22/task-12; displaced row c69456fc9bf4 phase-22/task-6 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; I |
| `cmd/fid/integration/p22_resume_actions_test.go`: p22NewResumeActionAdapter | `98020edaa18a` phase-22/task-12 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; I |
| `cmd/fid/integration/p22_run_replay_snapshot_test.go`: p22StartAfterAbandon | `9301b105bc6b` phase-22/task-12 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; I |
| `cmd/fid/integration/p22_semantic_fake_test.go`: adapter execute | `25caa7d40712` phase-22/task-9 | Accepted expectation/nullable SourceInput shape; removed repository rejected; I |
| `cmd/fid/integration/p23_claim_mode_test.go`: TestP23ReviewClaimRejectsDirtyPathAdoption | `35944c72aba0` phase-23/task-3 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; I |
| `cmd/fid/integration/p23_claim_mode_test.go`: TestP23ReviewRunRetainsPhase22Lifecycle | `35944c72aba0` phase-23/task-3 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; I |
| `cmd/fid/integration/p23_gate_grants_test.go`: TestP23GateRunAuthorizesEveryEnvelopeGrantThroughBuiltBinary | `8531522e9602` phase-23/task-6 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; I |
| `cmd/fid/integration/p23_human_failure_test.go`: TestP23StopFailureCodeMatrix | `180f5168e3f1` phase-23/task-19 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; I |
| `cmd/fid/integration/p23_human_failure_test.go`: p23StopFailureReworkFixture | `180f5168e3f1` phase-23/task-19 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; I |
| `cmd/fid/integration/p23_review_evidence_budget_test.go`: TestP23ReviewCyclesRealProcessJourney | `d05c288973df` phase-23/task-6; displaced row d1c016daa9f5 phase-23/task-13 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; I |
| `cmd/fid/integration/p23_review_record_test.go`: p23StartRun | `d215e40c3fbb` phase-23/task-4 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; I |
| `cmd/fid/integration/p23_review_record_test.go`: p23CloseAndRelease | `d215e40c3fbb` phase-23/task-4 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; I |
| `cmd/fid/integration/p24_flag_corpus_test.go`: TestP24FlagAndClaimKeyedLeafCorpus | `9bbbdb55feb7` phase-24/task-2 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; I |
| `cmd/fid/integration/p24_idempotency_test.go`: TestP24RunStartConflictNamesChangedClaimVersion | `8bddc91940dd` phase-24/task-2 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; I |
| `cmd/fid/review_test.go`: closeReviewRun | `a1ae1bb61363` phase-23/task-4 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; O |
| `cmd/fid/task12_execution_transition_test.go`: startExecutionTransitionFixture | `730128cd8842` phase-23/task-1; displaced row 7e419bfd8c23 phase-22/task-3 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; O |
| `internal/export/archive_fixture_test.go`: addJournalBackedExportFixtureFacts | `1be8e882ef68` phase-22/task-13 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; O |
| `internal/export/archive_fixture_test.go`: fixtureRepository | `c9b46655d880` phase-22/task-13 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; O |
| `internal/export/import_test.go`: TestImportRequiresFreshDestinationObservationBeforeRecovery | `cc73b092bf5e` phase-22/task-14 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; O |
| `internal/workflow/recovery_test.go`: TestRecoverRunForceRequiresCurrentHumanAuthorityAndFencesPriorRun | `b7dce8e747a3` phase-22/task-5; displaced row a0a93c7ce4dd phase-22/task-9 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; O |
| `internal/workflow/run_close_authority_test.go`: TestAbandonRunRequiresExactCurrentAuthorityAndObservation | `c0003e0c5df7` phase-22/task-8 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; O |
| `internal/workflow/run_close_authority_test.go`: assertAbandonProvenanceAndReplay | `c0003e0c5df7` phase-22/task-8 | Recorded start binding replaces displaced request repository; provenance/replay assertions retained; O |
| `internal/workflow/run_close_authority_test.go`: TestCloseRunRejectsUnqualifiedObservationVersionAndAuthority | `c0003e0c5df7` phase-22/task-8 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; O |
| `internal/workflow/run_close_test.go`: TestCloseRunDerivesTerminalStateFromTrustedObservation | `b965b24748f5` phase-22/task-8 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; O |
| `internal/workflow/run_close_test.go`: TestAbandonRunRequiresTerminalProofOrExactHumanForceAuthority | `b965b24748f5` phase-22/task-8 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; O |
| `internal/workflow/stop_conditions_test.go`: TestStopResolutionWorkEndedRequiresCurrentTerminalSubject | `f973bca7ea14` phase-22/task-2 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; O |
| `internal/workflow/stop_conditions_test.go`: TestStopResolutionWorkEndedAcceptsCompletedEnvelopeFact | `f973bca7ea14` phase-22/task-2 | Expected repository for Run caller; lifecycle/authority/replay/assertions retained; O |
| cmd/fid/gate_run_test.go: activeRunIDForGateTest | `f3350bcf582a` phase-21/task-22 | New Run fixture observes current tree; existing Run and Artifact contract preserved; O |
| cmd/fid/input_descriptor_test.go: TestP24JudgmentDescriptorSourcePositionsAndExamples | `7c69c2bbb6b8` phase-24/task-4 | Added Run Source position; raw examples decode strictly, expanded declared timestamp/digest/OID placeholders validate; O |
| cmd/fid/integration/p23_failure_contract_test.go: source-digest-mismatch / repository-mismatch producer rows | `d1c016daa9f5` phase-23/task-13; Source row amended by `ff6c923dcb70` phase-24/task-5 | Added accepted Run producers without changing typed failure assertions; I |

The Run-start fixture line in p22NewFIDRecoveryAdapterWithClaimModeAt originated at `f3be616f0587` phase-22/task-6; the renamed/claim-mode constructor signature originated at `35944c72aba0` phase-23/task-3. Likewise the original Run-start line in startExecutionTransitionFixture came from `7e419bfd8c23` phase-22/task-3 before its phase-23/task-1 signature change. Provenance includes these displaced-line introductions instead of attributing earlier assertions to later renames.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit

