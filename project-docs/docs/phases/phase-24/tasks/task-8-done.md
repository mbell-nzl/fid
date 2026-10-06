# Phase 24 - Task 8: Gate run and Evidence leaves; real freshness - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: The observation and Receipt freshness slice passes focused tests, but retained fixture convergence, remaining subject/parity/exclusivity evidence and the full task gate are unfinished.
- Recovery attempts: 0
- Next action: Dispatch a same-task Task 8 continuation from this checkpoint; converge the retained failing gate/Evidence fixtures with provenance, complete the evidence matrix, and run all five exact gates in order before closing.
- Human attention required: no

OUTCOME RECOVERY_REQUIRED

## Summary of what was implemented

- Gate run accepts expected_repository and observes under the lock after keyed replay, before Attempt/process admission. The captured repository feeds Attempt, Receipt, both gate Events and the execution-blocked writer; the terminal phase reuses the Attempt observation without re-observing.
- Registration and active Run checks use Task 6 SameWorktree, including genuine import binding. Focused original/import admission and exact plain-copy refusal pass.
- Evidence link takes its repository from the Receipt, Review or Human Outcome Event, rather than caller JSON. Receipt derivation has positive runtime evidence; Review and Human Outcome positive matrices still need completion.
- Evidence evaluate observes the current worktree and compares real tree/dirty state. Real tracked edits and reverts are covered across evaluation, terminal transition and completion. Evidence get retains evaluation against the stored subject target without observing Git.
- Added the four descriptors, help/example installers and field validation. Help, example templates, recursive structural validation and removed-member tests pass; exclusivity and fake/fid parity remain unfinished.
- This is durable partial implementation, not task completion. Retained direct Go callers still compile using explicitly temporary json-ignored Repository/TargetRepository fields; those fields must be removed after the caller fixtures converge. No legacy JSON input is admitted.

## Files created or modified

- `cmd/fid/artifact.go`
- `cmd/fid/evidence.go`
- `cmd/fid/gate.go`
- `cmd/fid/input_descriptor_test.go`
- `cmd/fid/input_descriptors_evidence.go`
- `cmd/fid/integration/p24_gate_observation_test.go`
- `cmd/fid/p24_gate_evidence_test.go`
- `internal/evidence/evidence.go`
- `internal/evidence/execution.go`
- `internal/evidence/execution_gate_authority_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md`
- `project-docs/docs/phases/phase-24/tasks/task-8-done.md`

## Documentation updated

- CLI specification: current gate run and Evidence inputs, observation, freshness and descriptor authoring surface.
- Flow specification: gate admission after changes and freshness through real observations.
- This report records partial evidence, failed attempts, provenance and next actions. The runner-owned run-ledger.md and triage-log.md were already dirty at dispatch and are excluded from the task commits.

## Validation Contract Coverage

These are Task 8 slice results from focused commands in this session. PASS here is evidence for the stated matrix only; the required whole-tree gate and all acceptance criteria have not passed.

- VAL-009: PASS - TestP24GateAndEvaluateUnsafeAndReplay/advanced-HEAD in cmd/fid/integration/p24_gate_observation_test.go:19 returns the original gate result byte-for-byte except replayed:true after a real tracked-file commit, with no projection or journal change.
- VAL-011: PASS for the four adopted descriptors - TestInputDescriptorRegistryMatchesLiveTreeAndTypes in cmd/fid/input_descriptor_test.go:236 checks live registry/fence and reflected types, including seeded mismatch rejection. The final retained inventory audit is still pending.
- VAL-012: PASS - TestP24GateEvidenceDescriptorProcessCorpus in cmd/fid/p24_gate_evidence_test.go:206 walks help metadata recursively for all four leaves, including the git command meanings of expected_repository.
- VAL-013: PASS - The same process corpus compares every example byte-for-byte with its ordered, two-space descriptor template.
- VAL-014: PASS - The same process corpus decodes each example into the leaf's actual reflected Go request type with unknown fields rejected.
- VAL-015: BLOCKED - Task 8 exclusivity rows have not been added and run; the existing adopted-leaf corpus does not yet cover these four leaves.
- VAL-016: PASS - The same process corpus exercises all applicable recursive unknown/duplicate/missing/null/type/enum/format/bounds refusals, asserting exact details and no non-project directory effects.
- VAL-017: PASS for direct leaf refusals - TestP24RemovedGateEvidenceMembers in cmd/fid/p24_gate_evidence_test.go:260 rejects repository on gate.run and target_repository on both Evidence leaves before workspace effects. Shared adapter parity is pending under VAL-035.
- VAL-022: PASS for successful gate observation - TestP24GatePreSpawnObservationAfterEdit in cmd/fid/p24_gate_evidence_test.go:18 checks independent Git HEAD/tree/status/root identity and the captured observation against the Attempt effect, Receipt and both Events; the spawned process adds a file without altering the recorded pre-spawn value. A budget-blocked Event observation assertion is still pending.
- VAL-023: PASS for the two observing leaves - TestP24GateAndEvaluateUnsafeAndReplay pairs normal gate acceptance and fresh evaluation in TestP24ReceiptFreshnessRealGuardedActions with unmerged/non-Git refusals for gate.run and evidence.evaluate. It asserts repository-unavailable/4, exact reason/detail shape and unchanged projection/journal. The shared project cause-class suite was not rerun in this checkpoint.
- VAL-024: PASS - TestP24GateExpectationBeforeAttempt in cmd/fid/p24_gate_evidence_test.go:98 covers HEAD and dirty mismatch, unchanged persistence and no spawned marker; TestP24GatePreSpawnObservationAfterEdit covers matching acceptance.
- VAL-025: PASS for the gate producer - Both mismatch rows compare the complete details object to expected:{head_oid,dirty} and observed:Repository. Retained producer-inventory convergence remains pending.
- VAL-029: PASS for Task 8 - TestP24GateOriginalCopyAndImport in cmd/fid/p24_gate_evidence_test.go:215 uses active Run/Claim state in original, whole-worktree plain copy and genuine export/import workspaces. Original/import succeed; copy returns exact invalid-argument/2 fallback details and no persistence change.
- VAL-032: PASS - TestP24GatePreSpawnObservationAfterEdit compares one identical observation across Attempt, Receipt and both gate Events after committing and editing since Run start, then verifies the gate actually spawned.
- VAL-033: PASS - TestP24ReceiptFreshnessRealGuardedActions in cmd/fid/integration/p24_gate_observation_test.go:77 has unchanged and real-edit/revert branches. Evaluation is fresh unchanged/reverted, stale after editing tracked base.txt; terminal transition and completion reject with evidence-stale while edited and succeed unchanged/reverted. Retained freshness journeys still require convergence without weakening this rule.
- VAL-034: BLOCKED - TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation in cmd/fid/p24_gate_evidence_test.go:135 proves Receipt target/Event derivation after the current tree moved. Review and Human Outcome positive subject/Event derivation assertions are not yet implemented.
- VAL-035: BLOCKED - The shared fake dispatch intersects the changed Task 8 leaves at gate.run; fake and fid builders still need converged expected_repository inputs and removed-member parity. No parity PASS is claimed.
- VAL-052: BLOCKED - Final ordinary and tagged gate wall times are unavailable. Focused timings are listed below and do not substitute the budgets.
- VAL-060: PASS - TestP24GateEvidenceDescriptorProcessCorpus executes each example from an empty non-project directory, checks exit 0/empty stderr and verifies no directory change.
- VAL-066: PASS for gate unsafe retries - TestP24GateAndEvaluateUnsafeAndReplay/unmerged and /non-Git replay the original success without any persistence changes. Distinct-key controls from otherwise viable active-Run/budget state reject repository-unavailable with no Attempt/Event.

## Tests run

Working directory: /home/mark/projects/fid unless a separate entrypoint working directory is stated. Every token/termination below was observed from its own command. Earlier RED attempts remain failed evidence; later results do not erase them.

- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24GatePreSpawnObservationAfterEdit$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24GatePreSpawnObservationAfterEdit$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24GateExpectationBeforeAttempt$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24GatePreSpawnObservationAfterEdit|TestP24GateExpectationBeforeAttempt)$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation|TestP24GateExpectationBeforeAttempt)$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation|TestP24GatePreSpawnObservationAfterEdit|TestP24GateExpectationBeforeAttempt)$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.198s` (terminated `exit_code: 0`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24GateEvidenceDescriptorProcessCorpus$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24GateEvidenceDescriptorProcessCorpus$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.328s` (terminated `exit_code: 0`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame|TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestEvidenceEvaluationMultiFailureGoldenSurvivesRepeatAndRebuild)$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24GateOriginalCopyAndImport$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.173s` (terminated `exit_code: 0`)
- `go test -tags integration -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24GateAndEvaluateUnsafeAndReplay|TestP24ReceiptFreshnessRealGuardedActions)$' ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24GateAndEvaluateUnsafeAndReplay|TestP24ReceiptFreshnessRealGuardedActions)$' ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReceiptFreshnessRealGuardedActions$' ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReceiptFreshnessRealGuardedActions$' ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.537s` (terminated `exit_code: 0`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24Gate.*|TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation|TestP24RemovedGateEvidenceMembers)$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.794s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24Gate.*|TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation|TestP24RemovedGateEvidenceMembers|TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestValidateGateRunReferencesAuthorizesAgainstEveryEnvelopeGrant)$' ./cmd/fid ./internal/evidence` - PASS - `ok  	github.com/mbell-nzl/fid/internal/evidence	1.030s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24GateAndEvaluateUnsafeAndReplay|TestP24ReceiptFreshnessRealGuardedActions)$' ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.723s` (terminated `exit_code: 0`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation$' ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.991s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24Gate.*|TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation|TestP24RemovedGateEvidenceMembers|TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestValidateGateRunReferencesAuthorizesAgainstEveryEnvelopeGrant)$' ./cmd/fid ./internal/evidence` - PASS - `ok  	github.com/mbell-nzl/fid/internal/evidence	1.033s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24GateAndEvaluateUnsafeAndReplay|TestP24ReceiptFreshnessRealGuardedActions)$' ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.574s` (terminated `exit_code: 0`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24GatePreSpawnObservationAfterEdit$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.887s` (terminated `exit_code: 0`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24GateExpectationBeforeAttempt$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.802s` (terminated `exit_code: 0`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24GatePreSpawnObservationAfterEdit|TestP24GateExpectationBeforeAttempt)$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.007s` (terminated `exit_code: 0`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation|TestP24GateExpectationBeforeAttempt)$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.128s` (terminated `exit_code: 0`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation|TestP24GatePreSpawnObservationAfterEdit|TestP24GateExpectationBeforeAttempt)$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.272s` (terminated `exit_code: 0`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24GateEvidenceDescriptorProcessCorpus$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.881s` (terminated `exit_code: 0`)
- `go test -tags integration -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24GateAndEvaluateUnsafeAndReplay|TestP24ReceiptFreshnessRealGuardedActions)$' ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.045s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid-phase24-task8-checkpoint ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `git diff --check` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `gofmt -w cmd/fid/artifact.go cmd/fid/evidence.go cmd/fid/gate.go cmd/fid/input_descriptor_test.go cmd/fid/input_descriptors_evidence.go cmd/fid/integration/p24_gate_observation_test.go cmd/fid/p24_gate_evidence_test.go internal/evidence/evidence.go internal/evidence/execution.go internal/evidence/execution_gate_authority_test.go` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `goimports -w cmd/fid/artifact.go cmd/fid/evidence.go cmd/fid/gate.go cmd/fid/input_descriptor_test.go cmd/fid/input_descriptors_evidence.go cmd/fid/integration/p24_gate_observation_test.go cmd/fid/p24_gate_evidence_test.go internal/evidence/evidence.go internal/evidence/execution.go internal/evidence/execution_gate_authority_test.go` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

The final focused ordinary race result also printed cmd/fid 10.369 s; internal/evidence printed 1.033 s. The final tagged focused race result printed integration 4.574 s. These selected-test times are not VAL-052 full-package timings.

Real entrypoint, independently of go test: /tmp/fid-phase24-task8-checkpoint, Git scratch project /tmp/fid-task8-entrypoint-ZzxPiy initialized with a committed empty tree and fid prefix PSM.

- `/tmp/fid-phase24-task8-checkpoint gate run --example` - PASS - `}` (terminated `exit_code: 0`)

The example probe's working directory was /tmp. It printed the expected_repository input template.

- `/tmp/fid-phase24-task8-checkpoint --json evidence evaluate --input -` - FAIL - `{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["single UTF-8 JSON object"],"argument":"input","value":null},"message":"invalid input","retryable":false,"subject":null},"ok":false,"operation":"evidence.evaluate"}` (terminated `exit_code: 2`)
- `printf '%s\n' '{"obligation_id":null}' | /tmp/fid-phase24-task8-checkpoint --json evidence evaluate --input -` - PASS - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"evidence.evaluate","result":{"evaluations":[]}}` (terminated `exit_code: 0`)

Both evaluate probes ran in /tmp/fid-task8-entrypoint-ZzxPiy. The first mistakenly supplied no stdin and correctly failed closed; the second supplies the required JSON and succeeds without an Event.

## Exact quality-gate results

The final five-command gate was not started on this unfinished tree. Each required command is BLOCKED for this RECOVERY_REQUIRED checkpoint; no output token or termination is claimed for an unexecuted command. The focused selections above are additional evidence only.

- `gofmt -l cmd internal` - BLOCKED - Not executed as the final gate; edited Go files were formatted with gofmt and goimports, but retained convergence is unfinished. No termination observed.
- `golangci-lint run ./...` - BLOCKED - Not executed; finish the remaining in-scope implementation and convergence first. No termination observed.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - Not executed; the retained focused gate/Evidence selection is still red on removed inputs. No termination observed.
- `go tool cover -func=coverage.out` - BLOCKED - Not executed; this checkpoint produced no final whole-tree coverage profile. No termination observed.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - Not executed; retained integration/fake convergence is unfinished. No termination observed.

No conditional command is prescribed by Task 8. QA-plan task evidence requirements are the task VAL rows above; the phase-wide QA commands and checklist belong to QA, which this session did not perform.

Budget evidence pending for all four required records: cmd/fid, internal/workflow, slowest other ordinary package and tagged integration. Limits remain ordinary 600 s/tagged 900 s; warn at 480 s/720 s respectively. No full-gate headroom claim is made.

## Deviations from spec

- Session ends with a durable in-scope checkpoint rather than a completed task. No protected contract, assertion, command, threshold or policy changed.
- First observation test failed twice with token FAIL and exit_code:1: initial fixture used a clean-only oracle after a real edit; the corrected fixture reached the behavioral RED where expected_repository was rejected. The bounded product repair added observation/expectation admission and shared recorded state. Its byte-identical command now passes (exit_code:0); the exact terminal token appears in Tests run.
- HEAD/dirty mismatch first failed with token FAIL (exit_code:1), returning the generic fallback. A first renderer repair also failed with FAIL (exit_code:1) because CLI details validation requires JSON maps, not Go structs. Mapping expected and observed to their accepted closed maps fixed it; both original focused commands were rerun byte-identically and passed (Tests run).
- Receipt link/evaluate first failed with FAIL (exit_code:1) on the caller-target requirement; descriptor registration then exposed the next FAIL (exit_code:1) on freshness still reporting fresh after an edit. Deriving the subject target, observing evaluation and comparing dirty/digest fixed those in-scope failures. All original selections were rerun byte-identically and passed (Tests run).
- Descriptor corpus first failed with FAIL (exit_code:1) because help/example installation was missing and the nullable evaluation field/process arguments did not match the read leaf. Added the installers, required nullable obligation_id and read-specific process invocation while retaining every structural assertion. The byte-identical corpus passes (Tests run).
- The guarded-action journey first failed twice with FAIL (exit_code:1) because its fixture traversed array reasons using a map-only helper and then assumed terminal details had reason objects. Diagnosis inspected the real typed envelope: terminal uses unsatisfied:["evidence-stale"]. The test now checks that exact field and completion's typed reasons, preserving the freshness requirement. Both original tagged selections were rerun byte-identically and passed (Tests run).
- Stored-read regression failed with FAIL (exit_code:1): the evaluate refactor inadvertently made evidence.get observe Git and reject repository-unavailable after Git removal. EvaluateStoredEvidence restores the earlier stored-target evaluation without observation; the byte-identical Receipt test and final focused race selection pass.
- The retained focused selection still fails with FAIL (exit_code:1) at invalid input.repository. Task 8's accepted replacement of caller repository caused this displacement; it remains this task's in-scope convergence work, not a reassigned owner or a pre-existing failure claim. Recovery must preserve its original assertions, add real freshness fixtures, and rerun that exact command.
- Direct evaluate smoke initially failed with invalid input (exit_code:2) because stdin was empty; supplying {"obligation_id":null} yielded the exact success/no-Event envelope recorded above.

Retained rows already changed, with introducing-commit provenance verified from git and their owning reports:

| Retained row | Introducing commit | Phase/task | Bounded convergence and preserved assertion |
|---|---|---|---|
| TestValidateGateRunReferencesAuthorizesAgainstEveryEnvelopeGrant | 8531522e9602b0f72717d6c9d889684228cfa2ee | phase-23/task-6 | Pass observed repository separately to the updated binding check; all two granted and three denied Actor rows remain. |
| p24DescriptorProcessCorpus structural process invocation | 8d1c83855c0f681515e44e2c3a61bbe543ef515b | phase-24/task-7 | evidence.evaluate omits unsupported mutation-only flags; recursive cases, exact details, stderr/exit and no-effect checks remain. |

Displaced rows still awaiting recovery, with git introducing provenance (their fixtures were not edited here):

| Retained assertion | Introducing commit | Phase/task | Next action under Task 8 |
|---|---|---|---|
| TestGateRunPersistsAttemptAndSuccessfulReceiptBeforeTerminalFrame | b69499116e70d4c88423a30daa2010a4ff221f14 | phase-21/task-15 | Supply current expected_repository; retain Attempt-before-Receipt/terminal ordering, trusted binding and immutability. |
| TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent | 7277139ee5d5dde4cd3627335d39b14f616f46a1 | phase-21/task-17 | Derive targets, observe evaluation, use real branch/tree/dirty fixtures; retain independent validity/freshness/availability including the accepted branch-label control. |
| TestEvidenceEvaluationMultiFailureGoldenSurvivesRepeatAndRebuild | 7277139ee5d5dde4cd3627335d39b14f616f46a1 | phase-21/task-17 | Produce real stale state while keeping Artifact inputs ignored; preserve simultaneous reasons, byte-identical repetition and rebuild evidence. |

The phase-23/task-6, phase-24/task-7 and phase-21/task-17 retained reports were consulted. Every additional displaced fixture in recovery still needs its own introducing-commit/phase/task and accepted-contract audit before edits.

## Concerns, technical debt, or notes for the Tech Lead

- Continue Task 8 only. No protected contradiction was identified; unfinished in-scope work does not require Human authority.
- Converge ordinary and tagged gate/Evidence callers, including export/import and walking-skeleton/recovery/freshness helpers. Preserve ignored Artifact paths, exact failures and independent proof/trust dimensions. Never simulate a current repository by mutating the input block.
- Remove the temporary json-ignored GateRunRequest.Repository, EvidenceLinkRequest.TargetRepository and EvidenceEvaluateRequest.TargetRepository fields once direct typed callers are adapted.
- Complete positive Review/Human Outcome target/Event derivation and link-without-observation controls; keep Human Outcome freshness always fresh.
- Add execution-blocked Event repository evidence, four-leaf example exclusivity and fake/fid parity over the dispatch-derived changed-leaf set. Audit retained descriptor/mismatch-producer registries.
- Run touched-package/tagged race tiers, format edited Go files, then the five exact commands serially on the final tree. Record all required package times and any headroom warnings; run the report finalization/lint sequence afterward.
- Runner-owned dirty run-ledger.md and triage-log.md are intentionally not staged. Dispatch base verified as 9b585b73bf78188fb34cca4938734ded597afe55.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `32a9ae535dbfb34be0dff7e369c715c2a4640b90`

## Commit subject

- `feat(evidence): checkpoint gate observation and real freshness`
