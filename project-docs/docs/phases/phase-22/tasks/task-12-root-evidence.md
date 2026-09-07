# Task 12 — Root delegate contribution

The Human authorized direct root implementation alongside canonical Developers.
Root owns the Human renderer, recovery adapter call, paired real-process Human
tests and the optional Human-authority integration fixture constructor. The
Developer owns domain projection, freshness, action predicates and final gates.
All commands below ran from `/home/mark/projects/fid` in the root session.

## Recovery Human result

The fixture creates ordinary agent or trusted Human authority through public
activation, modifies the approved base.txt path, observes terminated or live
executor state, and records an ordinary or force recovery with exact adopted
path and reconciliation. It compares Human replay text to the original JSON
result and proves the replay adds no Event. It requires prior/successor Run
and Claim IDs, authority, adopted paths as unproven work, known/unknown remaining
budgets, and force-specific known-live takeover plus residual filesystem risk.

`go test -tags integration ./cmd/fid/integration/... -run '^TestP22RecoveryHumanResultMatchesDurableJSON$' -count=1`

- First attempt: `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration	2.396s`
  (terminated `exit_code: 1`). Both actual recoveries succeeded; Human output
  contained only `recovered Run <successor>` and omitted the required context.
  This is semantic red evidence before the root's renderer change.
- After adding shared recovery/budget Human rendering and replacing the
  recovery adapter's one-line output, the byte-identical command printed
  `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.073s`
  (terminated `exit_code: 0`). JSON mutation semantics and replay are unchanged.

## Resume Human projection

`go test -tags integration ./cmd/fid/integration/... -run '^TestP22ResumeHumanMatchesEligibleAndStoppedRecoveryState$' -count=1`

- First attempt: `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration	1.510s`
  (terminated `exit_code: 1`). A real adopted-path recovery and trusted numeric
  budget baseline succeeded. Human resume still printed only Issue ID and
  completion status; JSON still omitted stops/attention/recoveries and the
  other pending task-12 projection fields. The test identifies those missing
  fields explicitly, rather than panicking on a type assertion.
- A second byte-identical attempt printed
  `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration	0.441s`
  (terminated `exit_code: 1`). It overlapped the Developer's announced
  ResumeIssueState signature/plumbing edit: the built CLI caller still passed
  three arguments while the domain had gained the repository path argument.
  This is a transient compilation result, not behavioral evidence. Root waits
  for the domain/plumbing handback before another invocation.

Resume Human rendering is implemented; root hands its contribution back after
the combined regression below. The Developer owns subsequent causal edits.
The Developer's five exact final-source gates are required separately; these
focused receipts do not replace them.

Combined root Human cases:
`go test -tags integration ./cmd/fid/integration/... -run '^TestP22(RecoveryHumanResultMatchesDurableJSON|ResumeHumanMatchesEligibleAndStoppedRecoveryState|ResumeHumanDistinguishesElapsedLeaseFromRecordedStatus)$' -count=1`
first terminated `exit_code: 1` with
`FAIL	github.com/mbell-nzl/fid/cmd/fid/integration [build failed]`.
The Developer's projection-ready signal preceded its own compile; concurrent
Budget serializer edits had missing delimiters in envelope.go/run.go. The
Developer acknowledged the premature signal and repaired those syntax errors.
This is compilation evidence only; root waits for the verified source signal.

After the Developer's verified compile and completed projection signal, the
byte-identical combined Human command printed
`ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.436s`
(terminated `exit_code: 0`). All three cases ran: ordinary/force recovery text,
eligible/Stop-blocked resume parity, and actual elapsed lease projection.

The effective-lease Human test creates a real five-second Claim through the
public CLI, waits for its actual stored expiry timestamp, then checks that
durable status/expiry remain unchanged while effective status becomes expired
in both JSON and Human views, with no added Event. Equality-boundary predicate
testing remains the Developer's deterministic domain obligation.

Recovery renderer provenance: `cmd/fid/run_recovery.go` was introduced by
`69941fc1b6262fdd7efaaeecfe72760390dbd511` (phase-22/task-6); DEC-196 assigned
its ordinary Human rendering to task 12. The retained resume adapter came from
`2313235ba463060538f2f0d4bcfa61f039fcdc13` (phase-21/task-18). Root freshly
verified both introducing commits with path-specific `git log --diff-filter=A`.

## Incomplete Activation And Gate Eligibility

Review 4's early independent feedback identified the activation subject as a
possible missing Stop/read-action boundary. Root added a real-process case in
`p22_resume_activation_test.go`: a previously eligible Gate, a public upgrade
prepare/recovery selection, projected activation Stop/Attention and blocked
action, followed by the actual Gate rejection and unchanged journal bytes.

`go test -tags integration ./cmd/fid/integration/... -run '^TestP22ResumeAndGateAgreeDuringIncompleteActivation$' -count=1`

- Attempt 1 printed
  `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration [build failed]`
  (terminated `exit_code: 1`): root passed an `any` value to a map helper
  without the required type assertion. This was corrected before execution.
- Attempt 2 printed
  `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration	1.625s`
  (terminated `exit_code: 1`): the retained first-activation fixture helper
  overwrote supplied predecessor/resource input with nil/defaults. Root used
  the complete accepted public prepare request directly; no product changed.
- Attempt 3 printed
  `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration	1.153s`
  (terminated `exit_code: 1`): the actual Gate correctly rejected with the
  retained `authority-unavailable` guard, while root had assumed the
  activation-specific agreement error. Fresh inspection of gate.go and
  evidence/registration.go established the accepted authority path. Root
  corrected the test's uncontracted code expectation and additionally requires
  the Resume Reason to match that actual authority guard and the failure's
  exact closed authority details. No product/error contract changed.
- Attempt 4, the byte-identical command, printed
  `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.110s`
  (terminated `exit_code: 0`). This verifies current implementation after the
  Developer's early boundary repair; it does not claim test-before-code TDD
  for that repair. All three construction failures remain above.

## Combined Handback Regression

After gofmt/goimports and the Developer's explicit stable-source window, root
ran:

`go test -tags integration ./cmd/fid/integration/... -run '^(TestP22DualAdapter(Recovery|Activation)Conformance|TestP22RecoveryHumanResultMatchesDurableJSON|TestP22ResumeHumanMatchesEligibleAndStoppedRecoveryState|TestP22ResumeHumanDistinguishesElapsedLeaseFromRecordedStatus|TestP22ResumeAndGateAgreeDuringIncompleteActivation)$' -count=1`

It printed `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	12.873s`
(terminated `exit_code: 0`). This executes the retained 37 recovery and 26
activation cases after constructor extraction, plus all four new root Human
and activation boundary tests. No root source edit follows this receipt.
Task 12's remaining domain tests and five full final-source gates belong to
the Developer and remain required before task closure.

## Retained Terminal Fixture Convergence

The current-target terminal/completion consumers displaced synthetic repository
fixtures. Root accepted only task18_completion_resume_test.go and the Stop
matrix for repair; the Stop matrix itself required no edit. The task-18 helper
now creates a fresh observed Git fixture rather than cloning a bound worktree,
and its declared generated Artifact is ignored through that fixture's local
Git exclude file. This keeps generated material outside the source freshness
target; actual retained-byte loss and every rejection assertion stay intact.
Other tests and global synthetic repository helpers remain as before.

Root's introducing-commit inventory (`git log --diff-filter=A`) identified
task18_completion_resume_test.go at
2313235ba463060538f2f0d4bcfa61f039fcdc13 (Phase 21 task 18) and the unchanged
Stop matrix at 4b690e8b863abac96c9abfd72709ac36db2235b2 (Phase 22 task 9).
The former is subject-preserving convergence under task 12's execution scope.

`go test ./cmd/fid -run '^(TestStopPredicatesAcrossIssueAndEnvelopeCommands|TestTerminalSuccessTransitionRejectsEveryOtherNonqualifyingEvidenceState)$' -count=1`

- Root first attempt: `FAIL	github.com/mbell-nzl/fid/cmd/fid	1.881s`
  (terminated `exit_code: 1`): synthetic fixtures could not supply actual Git
  observations, so guards failed before their intended Stop/Evidence predicate.
- After the bounded fixture repair, byte-identical rerun printed
  `ok  	github.com/mbell-nzl/fid/cmd/fid	11.211s`
  (terminated `exit_code: 0`). All ten Stop states and all four nonqualifying
  Evidence states execute their original assertions.

`go test ./cmd/fid -count=1` then printed
`FAIL	github.com/mbell-nzl/fid/cmd/fid	36.298s`
(terminated `exit_code: 1`). It found three further synthetic consumers:
qualified completion after Run release, incomplete Envelope all-reasons, and
the retained nonterminal-transition test's eventual terminal attempt. Root
converted the first two to existing observed variants; Developer owned the
third in task12_execution_transition_test.go. No expected rejection was weakened.

`go test ./cmd/fid -run '^(TestStopPredicatesAcrossIssueAndEnvelopeCommands|TestTerminalSuccessTransitionRejectsEveryOtherNonqualifyingEvidenceState|TestTerminalSuccessTransitionRejectsUnsuccessfulReceiptEvidence|TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease|TestEnvelopeCompleteReturnsAllCurrentReasonsWithoutMutation)$' -count=1`
printed `ok  	github.com/mbell-nzl/fid/cmd/fid	10.766s`
(terminated `exit_code: 0`). The broader byte-identical cmd/fid rerun follows
the Developer's matching fixture repair; final task-wide gates remain separate.

The byte-identical `go test ./cmd/fid -count=1` rerun printed
`ok  	github.com/mbell-nzl/fid/cmd/fid	33.890s`
(terminated `exit_code: 0`). Root hands the repaired fixture file back to the
Developer after this whole-package regression; the earlier 36.298s failed
attempt remains above.

## Budget Output And Historical Replay Contribution

Root subsequently accepted the new budget_projection_replay_test.go file and
the narrow committed-denial renderer/closed error Budget shape under DEC-201.
The previous handback receipt remains evidence for its stated source boundary;
this later contribution requires the final task-wide gates independently.
The Developer owned causal domain replay repairs. Tests were written against
the task-12 prototype; only the actual failures below claim semantic red.

`go test ./cmd/fid -run '^TestBudgetProjectionAndHistoricalReplayAcrossCLI$' -count=1`

- Attempt 1: `FAIL	github.com/mbell-nzl/fid/cmd/fid [build failed]`
  (terminated `exit_code: 1`), root used nonexistent journal.Line; corrected
  to journal.Event. This is test construction, not behavioral red.
- Attempt 2: `FAIL	github.com/mbell-nzl/fid/cmd/fid	0.574s`
  (terminated `exit_code: 1`). After three real Gate charges, Envelope create
  replay returned the current active/exhausted object instead of its original
  draft/available snapshot. Developer repaired stored-result replay paths and
  initial output derivation while retaining input-only policy Effects.
- Root byte-identical rerun: `ok  	github.com/mbell-nzl/fid/cmd/fid	0.574s`
  (terminated `exit_code: 0`). It covers create/activate, Run start, both current
  get surfaces, available/warning/exhausted transitions, historical create/start
  snapshots, replay persistence equality and canonical journal policy Effects.

`go test ./cmd/fid -run '^TestBudgetDenialIncludesDerivedOutputState$' -count=1`

- Attempt 1: `FAIL	github.com/mbell-nzl/fid/cmd/fid	0.619s`
  (terminated `exit_code: 1`): committed exhausted denial omitted Budget state.
- Attempt 2: `FAIL	github.com/mbell-nzl/fid/cmd/fid	0.537s`
  (terminated `exit_code: 1`): adding the manual renderer field exposed the
  product error serializer's retained Phase 21 closed shape, which rejected the
  accepted output. Root added exactly state to both Budget denial shapes and
  the matching contract fixture; DEC-201/task scope records that consequence.
- Byte-identical rerun: `ok  	github.com/mbell-nzl/fid/cmd/fid	0.559s`
  (terminated `exit_code: 0`). Denial replay retains counters/state and adds no
  persisted effect. No admission predicate, policy, schema or workflow control
  changed.

`go test ./cmd/fid -run '^TestBudgetUnknownProjectionPreservesNullCounters$' -count=1`

- Attempt 1: `FAIL	github.com/mbell-nzl/fid/cmd/fid [build failed]`
  (terminated `exit_code: 1`), overlapping Developer test-fixture construction:
  changed builder arity, undefined initializeObservedRepositoryForTest and
  temporary unused imports. Root waited for explicit compile-stable handback.
- Byte-identical rerun: `ok  	github.com/mbell-nzl/fid/cmd/fid	0.371s`
  (terminated `exit_code: 0`): external Budget without a sample stays unknown
  with null counters through Envelope create/activate and Run start/get;
  canonical policy Effects retain input-only values.

`go test ./internal/cli -count=1` printed
`ok  	github.com/mbell-nzl/fid/internal/cli	0.019s`
(terminated `exit_code: 0`).

Final root Budget set:
`go test ./cmd/fid -run '^TestBudget(ProjectionAndHistoricalReplayAcrossCLI|DenialIncludesDerivedOutputState|UnknownProjectionPreservesNullCounters)$' -count=1`
printed `ok  	github.com/mbell-nzl/fid/cmd/fid	1.122s`
(terminated `exit_code: 0`). The Developer receives all root Budget files after
this receipt and owns subsequent changes and the five final task-wide gates.

### Later Same-Package Size Extraction

At Developer request, root moved resumeNextActions through withSortedReasons
verbatim from resume.go to resume_actions.go, preserving behavior and importing
only sql/json/fmt. The NextAction comment now includes reads. gofmt/goimports
were applied; the files are 381 and 196 lines. Root handed the uncommitted
source back to its Developer owner for the final product commit and gates.
`go test ./internal/workflow -run '^TestResumeKeepsOpenAttentionAfterAllLinkedStopsAreResolved$' -count=1`
printed `ok  	github.com/mbell-nzl/fid/internal/workflow	0.133s`
(terminated `exit_code: 0`). This is verification of a mechanical extraction,
not a claimed behavior change or substitute for the final suite.

### Clean And Dirty Recovery, Then Rebuild

At the checkpoint reviewer's pre-final request, root accepted the existing
p22_recovery_journey_test.go from the Developer. Root table-drove its already
implemented journey with clean and dirty starting targets, using empty adoption
for the clean case and base.txt adoption for the dirty case. Both retain valid
fresh proof after recovery and branch rename, create a new successor
Gate/Receipt/Evidence, compare the full recovered Resume before/after real
rebuild, then change one byte at the same HEAD and require stale but valid
proof with retained original Evidence identity. Existing dirty assertions were
parameterized; neither state was removed. This is additional regression proof
against the completed prototype, not a retrospective red/green claim.

`go test -tags integration ./cmd/fid/integration -run '^TestP22EarlyRecoveryJourney$' -count=1`
printed `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.454s`
(terminated `exit_code: 0`). gofmt/goimports were applied. Root hands the
uncommitted shared test file back for the Developer's final commit/gates.

### Required Material And Metadata-Sufficient Artifact Loss

Root owned new artifact_material_policy_test.go. Both fixtures execute a Gate
that retains actual material, then link Evidence through the public operation
with required_artifact_ids either populated or empty. They delete retained
bytes and use artifact.verify, proving only required material blocks Evidence
evaluation. Both retain Evidence identity and provenance, Receipt facts, and
all Artifact metadata except its current availability. Receipt embeds current
Artifact availability, so that projection member alone may change; the
remaining Receipt value is compared in full. Missing bytes are not recreated.

`go test ./cmd/fid -run '^TestArtifactLossHonorsRequiredMaterialAndPreservesProofIdentity$' -count=1`

- Attempt 1: `FAIL	github.com/mbell-nzl/fid/cmd/fid [build failed]`
  (terminated `exit_code: 1`): root referenced ArtifactIDs omitted by the local
  terminal test struct. Read the existing public Receipt.artifacts instead.
- Attempt 2: `FAIL	github.com/mbell-nzl/fid/cmd/fid	0.456s`
  (terminated `exit_code: 1`): root assumed an artifact.get operation; it is
  not an accepted leaf. Use the accepted receipt.get nested Artifact view.
- Byte-identical rerun: `ok  	github.com/mbell-nzl/fid/cmd/fid	0.701s`
  (terminated `exit_code: 0`). These are construction failures and subsequent
  regression proof of existing policy, not red/green product repair. No product
  policy or serialization was changed for this test. Developer owns the new
  test after handback and includes it in the five final task-wide gates.

### Executable Ready Actions And Producer Rejection

Root added p22_resume_actions_test.go with a separate public fixture constructor
that grants renew/release/evidence.link and permits an explicit producer. It
does not change the shared fake/fid fixtures. Each of run.inspect, run.observe,
gate.run, claim.renew and claim.release gets a fresh real activated-harness
fixture, must appear ready in Resume, and then executes its actual public
operation. Reads append no Event; mutations append durable Events. A kit
producer with matching capability/grant must appear blocked for run.observe,
and its actual operation rejects authority without adding an Event.

`go test -tags integration ./cmd/fid/integration/... -run '^TestP22Resume(ReadyCandidatesExecuteTheirActualOperations|ObserveRejectsUntrustedProducerDespiteCapabilities)$' -count=1`

- Attempt 1: `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration [build failed]`
  (terminated `exit_code: 1`): root assumed an adapter obligationID field;
  corrected the request to the fixture's explicit OBL-p22-corpus identity.
- Byte-identical rerun: `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.000s`
  (terminated `exit_code: 0`). This verifies the Developer's producer repair
  after implementation and is not claimed as its original TDD red. Root
  hands back the complete new fixture/test file before final task-wide gates.

### Recovery Report Control And Git Permission Receipts

The first recovery ledger append began with a blank line and returned
`INVALID append candidate must start with an event heading`
(terminated `exit_code: 1`). Corrected input returned
`INVALID TASK-RECOVERY cites an unrecorded or malformed phase decision`
(terminated `exit_code: 1`): root had authored DEC-199..201 with H2/em-dash
headings, whereas the accepted sequence control recognizes canonical H3
`### DEC-NNN:` headings. Root corrected only those three heading formats;
their accepted authority and content are unchanged. A premature ledger commit
after the rejected append printed `nothing to commit, working tree clean`
(terminated `exit_code: 1`); the rejection had not appended any event. No
workflow control was edited. The same candidate is retried after the heading
conformance repair; that byte-identical retry printed
`APPEND_CONFIRMED state=REPAIR_PENDING` (terminated `exit_code: 0`).

The Developer produced implementation checkpoint 2339dc9 and report checkpoint
6e04604, then handed a report formatting correction to root after its Git write
failed. Root's first normal commit also printed
`fatal: Unable to create '/home/mark/projects/fid/.git/index.lock': Read-only file system`
(terminated `exit_code: 128`). Root used the required escalation surface; the
first escalated commit found the report unstaged and printed
`no changes added to commit (use "git add" and/or "git commit -a")`
(terminated `exit_code: 1`). Explicit approved staging/commit then succeeded
at 3c27a25. No Git metadata was changed through a bypass tool.

Root ran the exact task-12 report-lint command: first it printed
`FAIL report-lint found 4 failure(s)` (terminated `exit_code: 4`), because four
other focused PASS rows still ended with punctuation after termination. Root
removed that punctuation only, preserved this failed receipt in the originating
report, committed 2278953, and reran the byte-identical command. It printed
`OUTCOME RECOVERY_REQUIRED` and `PASS report-lint completed with no failures`
(terminated `exit_code: 0`). This validates recovery-report form, not task
completion. The Developer's separately reported earlier three-row formatting
failure remains attributed to that Developer; root does not claim its execution.
