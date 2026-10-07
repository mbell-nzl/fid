# Phase 24 - Task 9a: Grammar failures and unknown commands across the live tree - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

Local red/green and gate repair attempts are preserved below. No Runner continuation was dispatched during this Developer session.

## Summary of what was implemented

- Installed one recursive grammar mechanism over 90 stable leaves and one unknown-command mechanism over the root and 28 groups. The two Human-only diagnostics and Cobra-generated commands retain their scope exceptions. Removed the five group-local validators.
- Grammar failures now name the actual flag, value or positional argument, carry exact class details and operation IDs, and reject before workspace resolution in Human and JSON modes. Moved the existing issue-create kind/title and issue-transition target presence checks ahead of workspace access.
- Added the accepted operation-unavailable code and exact detail schema. Unknown words reject without suggestions; groups without a word retain Human help and exit 0.
- Added shared control escaping and escape-then-Unicode-scalar bounding. Both modes use the same final message; details retain the full offending spelling/value. Boundary tests include Go-escaped text, non-ASCII text and an escape crossing the cut.
- Declared closed and typed flag domains once by leaf/spelling for product parsing and corpus generation, including inherited text bounds, opaque IDs/states, CSV members, actors, idempotency keys and all 61 R18 flags. Preserved native pflag integer syntax, decimal-only expected revision, empty label clearing and amended superseded/cancelled members.
- Added in-process and real-binary corpora in both modes, inside and outside a project, with filesystem snapshots and mechanical completeness checks. Converged displaced retained assertions with commit provenance and their original rejection/state-integrity subjects.

## Files created or modified

- `cmd/fid/canonical_request.go`
- `cmd/fid/concern.go`
- `cmd/fid/contract_surface_test.go`
- `cmd/fid/export.go`
- `cmd/fid/flag_domain_text.go`
- `cmd/fid/flag_domains.go`
- `cmd/fid/grammar_corpus_test.go`
- `cmd/fid/integration/e2e_v01_closure_helpers_test.go`
- `cmd/fid/integration/exit_contract_test.go`
- `cmd/fid/issue.go`
- `cmd/fid/issue_ac.go`
- `cmd/fid/issue_ac_test.go`
- `cmd/fid/issue_create.go`
- `cmd/fid/issue_mutation_response.go`
- `cmd/fid/issue_transition.go`
- `cmd/fid/issue_transition_test.go`
- `cmd/fid/issue_update_test.go`
- `cmd/fid/link_test.go`
- `cmd/fid/policy_test.go`
- `cmd/fid/retained_actor_request_contract_test.go`
- `cmd/fid/root.go`
- `cmd/fid/root_test.go`
- `cmd/fid/run.go`
- `internal/cli/codes.go`
- `internal/cli/codes_test.go`
- `internal/cli/message.go`
- `internal/cli/message_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-24/tasks/task-9a-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: delivered grammar-class table, declared domains, unknown-command failures, full detail preservation and DEC-282 escape/bound precedence. Task 10's additional Human lines remain pending.
- This report records declaration citations, live-row disagreements, R17/R18/R19 inventories, retained assertion provenance and all validation attempts.
- The Runner-owned run ledger, triage log, decisions, QA review artifacts, task spec and validation contract were not edited.

## Validation Contract Coverage

- VAL-064: PASS — `TestGrammarWholeTreeClasses` (`cmd/fid/grammar_corpus_test.go:283`) covers unknown long/shorthand flags, existing required flags and positional count on every stable leaf, plus root/group flags. `TestGrammarWholeTreeFlagDomains` (:518), `TestGrammarFlagDomainRegistryCompleteness` (:577), `TestGrammarCSVAndKeyDomains` (:604), `TestGrammarInheritedStringDomains` (:649) and `TestGrammarOpaqueStateUTF8Domain` (:693) cover declared domains, missing values and completeness. Each failure is checked in both modes and project/non-project states with filesystem snapshots. `TestGrammarBuiltBinaryCorpus` (:523) repeats class/domain/boundary rows through the real CLI with exit 2, empty stdout, exact message and typed envelope checks.
- VAL-065: PASS — `TestGrammarWholeTreeClasses` includes unknown words at all 29 R17 paths; `TestGrammarBuiltBinaryCorpus` includes these, root `isue`, short exact messages and 600-character words in both modes and both directory states. `TestGrammarGroupHelpWithoutWord` (:544) retains help/exit-0 controls at all 29 paths in-process and in the binary. `TestReasonCodeDetailSchemaIsExactForEveryCode` (`internal/cli/codes_test.go:38`) proves the new closed schema.
- VAL-064 and VAL-065 boundary evidence: `TestGrammarMessageBoundaries` (:374), the binary corpus and `TestBoundedMessage` (`internal/cli/message_test.go:8`) cover escaped 512/513-scalar sources, non-ASCII counting, already escaped Go text, dropped crossing escapes and complete details. The binary long-flag spelling and long word are each exactly 600 ASCII characters.
- VAL-052: PASS — exact-gate package output and timing table below establish all four requested measurements, including the slowest other ordinary package.
- VAL-036: ADVANCED, not a full assertion PASS — the grammar/binary checks prove JSON message/Human error-payload equality for grammar and unknown-command failures. The accepted reason/blocked-by/next lines and remaining failure families await task 10.
- QA-plan section 8 D1 evidence is ready for independent review of tasks 9 and 9a; this Developer session did not author a QA verdict.

String-domain declaration inventory (leaf paths are without the leading `fid`; enumerations are stored sorted bytewise). `cmd/fid/flag_domains.go` and `cmd/fid/flag_domain_text.go` contain the declarations and citations used by parsing and the corpus:

| Leaf(s) / spelling(s) | Declared domain | Newest accepted row / inherited constraint |
|---|---|---|
| init --prefix | ^[A-Z][A-Z0-9]{0,15}$ | Phase 21 public-interface §4.1 project.init |
| issue create --kind, --scheduling | epic/subtask/task; prioritised/sequential | P21 §4.1 issue.create |
| issue list --kind, --status-category, --order, --priority | kind enum; cancelled/completed/started/unstarted; created/priority/sequence/updated; unique CSV 0..5 | P21 §4.1 issue.list; P20 contracts §3.2.1 retained constraints |
| workflow define --name, --applies-to, --states, --transitions, --categories | ^[a-z][a-z0-9-]{0,63}$; any/epic/subtask/task; 1..64 unique printable states of 1..128 scalars; 1..4096 unique from->to pairs with bounded printable endpoints; category CSV (1..64) | P21 §4.1 workflow.define |
| workflow assign --name, --kind; workflow list --kind | workflow-name pattern; epic/subtask/task | P21 §4.1 workflow.define/assign/list |
| link add/remove --from-type, --to-type; link list --type | issue/requirement | P21 §4.1 link.add/remove/list |
| stats cost/tokens --group-by | actor/epic/model | P21 §4.1 stats.cost/tokens |
| claim list --status | active/expired/released/superseded | P22 public-interface §2 Claim, amending P21 §4.2 claim.list |
| run list --status | abandoned/active/failed/interrupted/succeeded | P22 §2 Run status and §4 run.list, amending P21 §4.2 |
| receipt list --outcome | cancelled/failed/signaled/spawn-failed/succeeded/timed-out | P22 §5 cancellation and schema §1, amending P21 §4.2 receipt.list |
| stop list --status, --subject; attention list --status, --subject | open/resolved; <type>:<id>, activation/envelope/run types | P22 §4 stop.list/attention.list |
| review list --verdict | blocked/changes-required/pass | P23 public-interface §4 review.list / §2 Review |
| finding list --state, --severity | addressed/inapplicable/open/resolved; blocker/major/minor/nit | P23 §4 finding.list / §2 Finding |
| concern list --status, --severity, --subject | disposed/open; blocker/major/minor/nit; <type>:<id> with envelope/finding/issue/requirement/review/run types | P23 §4 concern.list / §2 Concern and JudgmentSubject |
| deviation list --disposition, --subject | accepted/pending-decision/rejected; same JudgmentSubject pattern/types | P23 §4 deviation.list / §2 Deviation and JudgmentSubject |
| decision list --outcome | approved/deferred/rejected | P23 §4 decision.list / §2 DecisionRef |
| human-outcome list --result | acknowledged/approved/changes-required/rejected | P23 §4 human-outcome.list / §2 HumanOutcome |
| requirement list --status | active/superseded/withdrawn | P23 §4 requirement.list / §2 RequirementRef |
| activity list --event | all 23 amended ActivityEvent members in the declaration | P21 §3 ActivityEvent amended by P23 §2 ActivityEvent |
| issue create/update --title; issue ac add --text; label add --colour, --description; issue list --label | normalized nonempty text, maxima 512/4096/64/4096/255 UTF-8 bytes respectively | P21 §4.1 retaining P20 contracts §3.2.1 |
| issue create --label; issue update --labels | CSV/repeated normalized labels, 255-byte token bounds; update's explicit empty value clears labels | P21 §4.1 retaining P20 contracts §3.2.1 repeated-label rules |
| issue create --epic, --parent; issue list --epic; ready --epic; stats cost/tokens --epic | opaque printable UTF-8 ID, 1..128 scalars after existing trim | P21 §2 ID and §4.1 each named row |
| issue list --status; issue transition --to | printable UTF-8 opaque state, 1..128 scalars | P21 §4.1 issue.list/transition |
| claim list --issue, --envelope; receipt list --run, --obligation | opaque printable UTF-8 ID, 1..128 scalars | P21 §2 ID / §4.2 each named row |
| run list --issue, --envelope | same ID domain | P22 §4 run.list / P21 §2 ID |
| review list --issue, --envelope; finding list --issue, --review; decision list --issue, --class; human-outcome list --issue, --class; requirement list --issue | same ID domain | P23 §4 each named row / §2 Review, Finding, DecisionRef, HumanOutcome, RequirementRef |
| standalone --actor on non-activation leaves without --input | ^(agent\|human\|service):[A-Za-z0-9][A-Za-z0-9._-]{0,127}$ | P21 §2 standalone Actor |
| --actor on non-activation --input mutations | opaque activated Actor ID, printable UTF-8 1..128 scalars | P21 §2 activated Actor / §4.2 mutation rows, inherited by P22/P23 |
| all 50 --idempotency-key flags | 1..128 printable ASCII bytes, no surrounding whitespace | P21 §2 idempotency key |
| issue update/transition --expected-revision | decimal digits representing 1..9223372036854775807 | P21 §4.1 retaining P20 contracts §3.2.1 |

R18 additionally declares all 17 native integer flags and 44 booleans: priority 0..5, sequence/ordinal 1..2147483647, gate.resolve version 1..9223372036854775807, twelve list limits 1..10000 (attention, claim, concern, decision, deviation, finding, human-outcome, receipt, requirement, review, run, stop), five existing leaf booleans and 39 adopted --example booleans (P24 §8). Missing free-string/stringSlice values use their pflag type domain. Input file paths, repository paths and bootstrap activation Actor text retain their existing contextual/structural validation; input shape validation is outside this grammar task.

CF-047 live-row disagreements retained: activity.list has --on rather than the P21 --ref and has no --limit; rebuild has live --verbose although its row has no leaf flags; issue.update/transition keep optional --expected-revision although the row lists it required. Corpus scope follows the live flags. The forbidden workflow.define --require is absent and agrees with the contract; the 39 --example flags are accepted Phase 24 amendments.

Combination rules found and retained: issue.create --epic/--parent exclusivity and kind-conditional parent/scheduling/sequence rules; epic scheduling requirement; sequence's sequential-parent constraint; workflow.assign exactly one of --kind/--project-wide; workflow.define parallel state/category counts, applicable cancelled-state count and transition endpoints referring to its states; issue.update at least one supplied field and an effective change; rebuild --verbose/--json restriction; descriptor --example ordered exclusivity as task 3 specifies. These are not any of the four grammar classes. The workflow.assign combination helper keeps its existing request envelope, including wire argument:"request", using the named constant flagCombinationArgument. It is explicitly not a grammar path; the mandated literal R19 query still enumerates exactly the four family fallback map literals.

Retained assertion convergence provenance (commit introduction, phase/task and preserved subject; historical task numbers are not current ownership):

| Retained assertion | Introducing commit and Refs | Convergence / preserved strength |
|---|---|---|
| Exact reason-code schema registry | 51763124db2cd6f6c6d00d15548de37532bbc273; phase-20/task-3a | Add accepted operation-unavailable to the exhaustive registry, retaining exact-key/type rejection checks. |
| Live leaf inventory | 37b59d629b296b7b5383454170f9f017fdd94253; phase-20/task-18; current removal allocation 0dd4c57849f8be7779e1c2e43cf7b14093c945f0; phase-23/task-15 | Runnable groups are not leaves; retain the exact stable/diagnostic inventory and add independent exact 29-group coverage. |
| Removed owner-flag envelopes | 1f196a94b064554fbc8bae1b0d9238125e2ab1d4; phase-20/task-9a | Replace generic request details with offending spelling and all local allowed flags; retain rejection, envelope metadata and no-state-change checks. |
| Issue priority Human failure | da2c21e48a02b8391a05e3e25474188da0f7bc9f; phase-5/task-3 | Exact domain-source text replaces the former generic text; issue/journal immutability checks remain. |
| Link missing-flag envelopes | cf096da6391c020c24e3392996548158271bc080; phase-20/task-11 | Exact required-flag details replace request fallback details; atomic journal/projection rejection checks remain. |
| Retained Actor/request JSON and Human fixtures | bc328f9136ec44d31824d1fe2b1c63d1c8b7c921; phase-21/task-9 | Actual flag/value/count source and class details replace generic text; original rejection and journal checks remain. |
| Absent policy with --by | 026d08ad663ebc50467b5031cf1e9c9b6807aa58; phase-21/task-9 | Preserve the original invocation; unknown --by is detected before the unknown root word. Exact typed failure and no-write checks replace generic root details. |
| Absent governance commands | 95ee13d090b424ad541971d7fafe2ce295f09143; phase-20/task-17 | Keep command absence and empty stdout; require exact Human unknown-command text and exit 2. |
| Attention production-path inventory | ac827959f8359473cbd6546ac6bea1ff1fd20751; phase-22/task-2; descriptor entry 04b6b0992fe20a02df402a183269a9847ddae5ab; phase-24/task-3 | Add the actual flag-domain declaration path; retain all other footprint and no-write assertions. |
| AC ordinal/Actor Human failures | 59f02dfeccbb636887cfa62ee3e73661a1c74317; phase-7/task-1; Actor expectation 55ba57e10016ac25b5b1af4f7dd6e33402b00377; phase-21/task-9 | Exact class text replaces prior text, retaining criteria and journal atomicity. |
| Issue-transition missing --to | 0173eba6d4fcff718a76e89a902f74bbfecce1da; phase-5/task-5 | Preserve required-target rejection, now explicit and before workspace access. |
| Real-process removed leaf output | 0dd4c57849f8be7779e1c2e43cf7b14093c945f0; phase-23/task-15 | Original AC/Concern removed-call arguments and event-count checks remain; exact typed unknown-flag envelope, empty stdout and exit 2 replace silent stderr. |
| Closure journey absent-command helper | 16b15d8134e7b29f27cf996486bff997a080fcc3; phase-21/task-12 | Original requirement.sync, run.record, issue.ac.verify and concern calls remain; exact Human and companion JSON failures plus journal-byte immutability replace empty stderr. |
| Closure journey absent policy.show | 6803f04927eb28986a7926311f99b4f5f5388629; phase-21/task-9 | Root operation-unavailable in both modes replaces silence; original absence and entire-journey journal-byte checks remain. |

Provenance was obtained with git log -S on the displaced assertion text, not inferred from test names. Receipt cancellation assertions introduced by c5dc59e4bb9cf9db20a4aa157daeef68d132bb60 (phase-22/task-8) were not weakened: their first failure exposed this task's missed cancelled member, fixed in the declaration. The same audit retained Claim superseded from the P22 amendment. Rebuild/init grammar and evidence.evaluate's custom example parser remain covered and their accepted assertion subjects were preserved.

## Tests run

All commands in this section ran from /home/mark/projects/fid unless a directory is explicitly named. Edited Go files were prepared with gofmt -w and goimports -w before the required format/lint gates. For the two integration convergence edits, each of `gofmt -w cmd/fid/integration/exit_contract_test.go cmd/fid/integration/e2e_v01_closure_helpers_test.go` and `goimports -w cmd/fid/integration/exit_contract_test.go cmd/fid/integration/e2e_v01_closure_helpers_test.go` ran twice, each terminating exit_code: 0 with `exit 0 (no output)`. Before each exact gate attempt, `go clean -testcache` terminated exit_code: 0 as fresh-timing setup; cache cleanup is not PASS evidence. RED means an observed expected or diagnosed failed attempt, not passing evidence. Commands and final output tokens below retain every test termination, including fixture construction mistakes and an empty selection.

- `go test ./cmd/fid -run '^TestGrammarUnknownActivationCommand$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  First behavior red: activation bogus printed help with success.

- `go test ./cmd/fid -run '^TestGrammarUnknownActivationCommand$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.153s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestGrammarWholeTreeClasses$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Silent classes and initial envelope fixture key-count mismatch.

- `go test ./cmd/fid -run '^TestGrammarWholeTreeClasses$/002/fid_issue_list_flag$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Specific silent leaf-flag red.

- `go test ./cmd/fid -run '^TestGrammarWholeTreeClasses$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Build failed on an unused fmt import after removing a group validator.

- `go test ./cmd/fid -run '^TestGrammarWholeTreeClasses$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Custom evidence.evaluate Args path displaced; preserved its DisableFlagParsing handler.

- `go test ./cmd/fid -run '^TestGrammarWholeTreeClasses$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.017s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestGrammarMessageBoundaries$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Message escaping and bound were missing.

- `go test ./cmd/fid -run '^TestGrammarMessageBoundaries$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.173s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestGrammarFlagDomainBeforeWorkspace$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Declared domains reached workspace/family fallbacks.

- `go test ./cmd/fid -run '^TestGrammarFlagDomainBeforeWorkspace$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.231s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestGrammarWholeTreeClasses$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.685s` (terminated `exit_code: 0`)
- `go test ./internal/cli -run '^TestReasonCodeDetailSchemaIsExactForEveryCode$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	0.026s` (terminated `exit_code: 0`)

  Initial conformance check.

- `go test ./internal/cli -run '^TestReasonCodeDetailSchemaIsExactForEveryCode$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Added expected accepted code to exhaustive conformance list first; registry member was missing.

- `go test ./internal/cli -run '^(TestReasonCodeDetailSchemaIsExactForEveryCode|TestBoundedMessage)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	0.025s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^(TestLiveLeafCommandsMatchContractInventories|TestIssueCreateAndUpdateRejectRemovedOwnerFlagsWithCanonicalFailures|TestRebuildJSONRejectsClosedRequestsBeforeProjectionWork|TestInitCommandEnforcesExplicitPrefixGrammarBeforeStaging)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Retained runnable-group and owner-flag assertions displaced.

- `go test ./cmd/fid -run '^(TestLiveLeafCommandsMatchContractInventories|TestIssueCreateAndUpdateRejectRemovedOwnerFlagsWithCanonicalFailures|TestRebuildJSONRejectsClosedRequestsBeforeProjectionWork|TestInitCommandEnforcesExplicitPrefixGrammarBeforeStaging)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.399s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestGrammarWholeTreeFlagDomains$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Corpus initially did not dispatch stringSlice flags.

- `go test ./cmd/fid -run '^TestGrammarFlagDomainRegistryCompleteness$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Missing declarations for the 44 leaf booleans.

- `go test ./cmd/fid -run '^(TestGrammarFlagDomainRegistryCompleteness|TestGrammarWholeTreeFlagDomains)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.517s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestGrammarCSVAndKeyDomains$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  CSV/key checks not yet installed.

- `go test ./internal/cli -run '^TestBoundedMessage$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Literal backslash-users text was mistaken for a valid Unicode escape.

- `go test ./cmd/fid -run '^(TestGrammarCSVAndKeyDomains|TestGrammarWholeTreeFlagDomains|TestGrammarWholeTreeClasses)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	5.946s` (terminated `exit_code: 0`)
- `go test ./internal/cli -run '^TestBoundedMessage$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	0.006s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^(TestGrammarBuiltBinaryCorpus|TestGrammarGroupHelpWithoutWord)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	42.792s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^(TestGrammarIntegerDomainsPreservePflagValues|TestGrammarWholeTreeClasses)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Native pflag 0x1 input was incorrectly constrained to decimal.

- `go test ./cmd/fid -run '^(TestGrammarIntegerDomainsPreservePflagValues|TestGrammarWholeTreeClasses)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	5.125s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/cli` - FAIL - `FAIL` (terminated `exit_code: 1`)

  First touched-package tier; cmd/fid 306.975s, internal/cli 1.121s; retained output expectations and cancelled member failed.

- `go test ./cmd/fid -run '^TestGrammarAmendedFlagDomains$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Accepted P22 Claim superseded and Receipt cancelled members must remain valid.

- `go test ./cmd/fid -run '^(TestGrammarAmendedFlagDomains|TestGrammarFlagDomainBeforeWorkspace|TestIssueUpdateRejectsInvalidInputAndLeavesStateUnchanged|TestLinkJSONRequestAndTypedFailuresAreAtomic|TestRetainedInventoryJSONMissingActorsUseCanonicalEnvelope|TestRetainedInventoryJSONRequestErrorsUseCanonicalEnvelope|TestRetainedInventoryRequiredActorHumanErrorsRemainProse|TestPolicyShowIsReadOnlyAndDoesNotAcceptByFlag|TestIssueACAddCommandRejectsBeforeAppend|TestPhase22AttentionReadBoundary|TestPromisedGovernanceCommandsAreAbsent|TestReceiptReadFailuresAreTypedAndClosed)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.335s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/cli` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	1.118s` (terminated `exit_code: 0`)

  Second touched-package tier; cmd/fid 302.788s, internal/cli 1.118s.

- `go test ./cmd/fid -run '^TestGrammarWholeTreeClasses$/.*missing_(kind|title|to)$' -count=1` - NO TESTS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.226s [no tests to run]` (terminated `exit_code: 0`)

  SELECTION MISS: no tests ran; this is not behavioral PASS evidence.

- `go test ./cmd/fid -run '^TestGrammarWholeTreeClasses$/[^/]+/fid_issue_(create_missing_(kind|title)|transition_missing_to)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Existing presence checks still ran after opening the project.

- `go test ./cmd/fid -run '^TestGrammarWholeTreeClasses$/[^/]+/fid_issue_(create_missing_(kind|title)|transition_missing_to)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.458s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/cli` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	1.115s` (terminated `exit_code: 0`)

  Third touched-package tier; cmd/fid 302.732s, internal/cli 1.115s.

- `go test ./cmd/fid -run '^TestGrammarInheritedStringDomains$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Inherited text/ID checks missing and long-title fixture expectation unbounded.

- `go test ./cmd/fid -run '^(TestGrammarInheritedStringDomains|TestGrammarWholeTreeFlagDomains|TestGrammarFlagDomainRegistryCompleteness)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Product now bounded correctly; long-title fixture still expected unbounded text.

- `go test ./cmd/fid -run '^(TestGrammarInheritedStringDomains|TestGrammarOpaqueStateUTF8Domain)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  Invalid UTF-8 opaque state must reject.

- `go test ./cmd/fid -run '^(TestGrammarInheritedStringDomains|TestGrammarOpaqueStateUTF8Domain|TestGrammarWholeTreeFlagDomains|TestGrammarWholeTreeClasses|TestGrammarFlagDomainRegistryCompleteness|TestLabelAddRejectsExplicitlyEmptyOptionalFlagsBeforeAppend|TestIssueTransitionRequiresTargetStateAndKnownByRole|TestIssueUpdateRejectsInvalidInputAndLeavesStateUnchanged|TestIssueACAddCommandRejectsBeforeAppend)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	7.052s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/cli` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	1.114s` (terminated `exit_code: 0`)

  Final ordinary touched-package tier; cmd/fid 309.908s, internal/cli 1.114s.

- `go test ./cmd/fid -run '^(TestGrammarBuiltBinaryCorpus|TestGrammarGroupHelpWithoutWord)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	34.589s` (terminated `exit_code: 0`)

  Final binary corpus including exact 600-character flag spelling.

- `go test -tags integration ./cmd/fid/integration -run '^(TestProcessExitContract|TestV01ClosureJourneyAndRebuildParity)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)

  After initial required-gate failure convergence, the closure fixture reached its later policy.show silent assertion.

- `go test -tags integration ./cmd/fid/integration -run '^(TestProcessExitContract|TestV01ClosureJourneyAndRebuildParity)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.077s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestProcessExitContract|TestV01ClosureJourneyAndRebuildParity)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.688s` (terminated `exit_code: 0`)
- `git diff --check` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

  Observed before the final exact-gate attempt.


The initial required-gate integration failure is retained under the exact gates and deviation note. The first focused integration rerun reached the later policy.show silence assertion; its failure is also retained above. All original removed-command argument lists remain intact.

Real entrypoint interaction, cwd /tmp: `/tmp/tmp.kRbUbJnSBo/fid --json activation bogus` rejected with observed exit_code: 2, empty stdout and this actual stderr token (a required rejection, not an undeclared nonzero gate PASS):

```json
{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"operation-unavailable","details":{"allowed":["abort","commit","component","finalize","prepare","recover","recover-complete","resource","status"],"argument":"command","value":"bogus"},"message":"unknown command \"bogus\" for \"fid activation\"","retryable":false,"subject":null},"ok":false,"operation":"activation"}
```

Observed termination: `exit_code: 2`.

The final built-binary corpus executes every changed grammar/domain/unknown-command row through the CLI in Human and JSON modes, inside a project and outside one, and checks the observed process statuses and file snapshots.

Inventory section 1's Bash block was executed verbatim from the repository root twice, by passing its extracted block to bash. The final invocation generated:

Command block supplied unchanged on bash stdin:

```bash
S=$(mktemp -d); go build -o "$S/fid" ./cmd/fid
subs() { "$S/fid" "$@" --help | awk '/Available Commands:/{f=1;next} /^$/{f=0} f{print $1}' | grep -vx -e help -e completion; }
walk() { local s l; l=$(subs "$@"); if [ -z "$l" ]; then echo "$*"; else for s in $l; do walk "$@" "$s"; done; fi; }
groups() { local s l; l=$(subs "$@"); if [ -n "$l" ]; then echo "${*:-fid}"; for s in $l; do groups "$@" "$s"; done; fi; }
walk > "$S/leaves"                                                                  # all live leaves: 92
groups > "$S/groups"                                                                # root + command groups: 29 (R17)
while read -r l; do if "$S/fid" $l --help | grep -q -- '--input'; then echo "$l"; fi; done < "$S/leaves" > "$S/input"        # --input leaves: 39 (this section)
while read -r l; do if "$S/fid" $l --help | grep -q -- '--idempotency-key'; then echo "$l"; fi; done < "$S/leaves" > "$S/keyed" # keyed: 50 (R16)
while read -r l; do if "$S/fid" $l --help | head -1 | grep -q '^Usage:'; then echo "$l"; fi; done < "$S/leaves" > "$S/no-short"  # no Short: 20 (R15)
while read -r l; do "$S/fid" $l --help | awk -v leaf="$l" '/^Flags:/{f=1;next} /^$|^Global Flags:/{f=0} f{for(i=1;i<=NF;i++) if($i~/^--/){n=$i; t=$(i+1); if(t!~/^(string|strings|stringArray|int|int64|uint|float64|duration|count)$/) t="bool"; if(n!="--help" && t!~/^string/) print leaf "\t" n "\t" t; break}}'; done < "$S/leaves" > "$S/typed-flags"  # typed leaf flags: 22 (R18)
wc -l "$S/leaves" "$S/groups" "$S/input" "$S/keyed" "$S/no-short" "$S/typed-flags"
```

- `bash` (inventory section 1 block on stdin), attempt 1 - PASS - ` 271 total` (terminated `exit_code: 0`)
- `bash` (same block on stdin), attempt 2 - PASS - ` 271 total` (terminated `exit_code: 0`)

Final output:

```text
  92 /tmp/tmp.HP4M9fLom1/leaves
  29 /tmp/tmp.HP4M9fLom1/groups
  39 /tmp/tmp.HP4M9fLom1/input
  50 /tmp/tmp.HP4M9fLom1/keyed
   0 /tmp/tmp.HP4M9fLom1/no-short
  61 /tmp/tmp.HP4M9fLom1/typed-flags
 271 total
```

Both invocations terminated exit_code: 0 with final token ` 271 total`. R17 = 29 root/groups; R18 = 61 typed flags (17 integer, 44 boolean). The remaining section-1 counts are 92 leaves including the two diagnostics, 39 input leaves, 50 keyed leaves and zero shorthand declarations. The block's mktemp directories were fresh; scratch cleanup was not gate evidence.

## Exact quality-gate results

The six commands below were run sequentially in exactly the listed order, and each process was watched terminate. Attempt 1's gate 5 failed and its required gate 6 was still executed; after causal retained-assertion convergence and the tagged race tier, attempt 2 reran all six byte-identical commands on the final source/test tree. No narrower command substituted for a gate.

Attempt 1 (captured in this session; immutable failed integration attempt):

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	107.322s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `rg -n '"argument": *"request"' cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/issue_get_list.go:129:	return cli.NewFailure(cli.CodeInvalidArgument, "invalid issue read argument", nil, map[string]any{"argument": "request", "value": "", "allowed": []string{"accepted issue read input"}})` (terminated `exit_code: 0`)

Gate 5's BLOCKED status describes this failed attempt only; the causal repair and later complete exact attempt follow. Ordinary timings were cmd/fid 451.329s, internal/workflow 107.322s and internal/export 210.863s (slowest other package); tagged integration 367.898s. Coverage was 62.3%.

Final attempt 2 (all gates PASS):

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	115.306s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	340.634s` (terminated `exit_code: 0`)
- `rg -n '"argument": *"request"' cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/issue_get_list.go:129:	return cli.NewFailure(cli.CodeInvalidArgument, "invalid issue read argument", nil, map[string]any{"argument": "request", "value": "", "allowed": []string{"accepted issue read input"}})` (terminated `exit_code: 0`)

Fresh final gate timing (go test package wall times; no cached test results):

| Package | Observed seconds | Budget | 80% warning |
|---|---:|---:|---:|
| cmd/fid | 399.177 | 600 | 480 |
| internal/workflow | 115.306 | 600 | 480 |
| internal/issues (slowest other ordinary package) | 196.455 | 600 | 480 |
| cmd/fid/integration (tagged) | 340.634 | under 900 | 720 |

Total coverage: 62.3%. All observed package times are inside budget and below their warning levels.

R19 actual final output (observed exit_code: 0; exactly four permitted family-fallback matches):

```text
cmd/fid/link_response.go:89:		"argument": "request", "value": nil, "allowed": []string{"accepted " + operation + " input"},
cmd/fid/receipt.go:87:	return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Receipt request", nil, map[string]any{"argument": "request", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
cmd/fid/issue_mutation_response.go:221:	return cli.NewFailure(cli.CodeInvalidArgument, "invalid Issue mutation request", nil, map[string]any{"argument": "request", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
cmd/fid/issue_get_list.go:129:	return cli.NewFailure(cli.CodeInvalidArgument, "invalid issue read argument", nil, map[string]any{"argument": "request", "value": "", "allowed": []string{"accepted issue read input"}})
```

Remaining causes reaching those family fallbacks:

- issueMutationFailure: existing cross-field/kind/parent/scheduling/sequence relationships, update field-group/effective-change rules, workflow state/category/transition endpoint relationships, and otherwise unmapped workspace/path/prose/persistence errors. Existing typed project/lock/journal/idempotency/not-found cases retain their own codes.
- issueReadFailure: positional ID validity (not positional count), semantic reference/filter relationships and otherwise unmapped workspace/database/path/read failures. Typed not-found/replay/result-limit mappings remain.
- link_response.go fallback: otherwise unmapped endpoint/relation combinations and persistence/read errors. Typed unknown-relation, missing endpoints, duplicate/missing link, cycle/self/result-limit and mapped Issue failures remain.
- renderReceiptError: otherwise unmapped workspace/project/replay/path/database/read/decode failures; Receipt-not-found retains its typed mapping.

None of these residual causes is an unknown flag, missing required flag, malformed declared flag value or positional-count failure. Those classes are rejected before the family handlers.

Report-lint is the independent Runner control specified by task-9a, not a Developer gate, and was not used as evidence or a gate substitute.

## Deviations from spec

- No product, architecture, accepted command/assertion semantics, threshold, budget, policy or scope deviation. Retained assertions converged only where task-9a's accepted grammar/unknown-command change displaced their expected output; provenance and preserved rejection/state subjects are listed above.
- Mandatory failed-attempt note: the first touched-package race tier terminated exit_code: 1 with final token `FAIL`. It exposed displaced generic/silent retained expectations plus this task's omitted cancelled Receipt member. The in-scope repair named the exact grammar classes, converged those expectations with provenance, and restored the amended closed member. The byte-identical touched-package command later terminated exit_code: 0; all attempts are retained in Tests run.
- Mandatory required-gate note: attempt 1's tagged integration command terminated exit_code: 1 with final token `FAIL` after 367.898 seconds. Removed AC/Concern invocations with unsupported flags expected empty stderr; requirement.sync in the closure journey expected silence. The accepted central parser now rejects those flags with typed class details. After preserving the calls and changing only their displaced assertions, the focused rerun terminated exit_code: 1 with `FAIL` at the later policy.show assertion. Its source was the same accepted unknown-command behavior at the root. Exact Human and companion JSON assertions were added while retaining journal checks. The focused and tagged race checks then passed, and the byte-identical full six-command attempt 2 recorded above closed the required gate.
- Other observed red attempts were local TDD and fixture corrections: four-vs-five top-level envelope keys, stringSlice corpus dispatch, an unused import, preserving evidence.evaluate's DisableFlagParsing Args path, adding all boolean declarations, valid Go escape detection, preserving native integer 0x1 syntax, amended Claim/Receipt members, earlier required-flag checks, inherited string bounds, an independent bounded long-title expectation, and invalid UTF-8 state rejection. Their exact command terminations are retained; no assertion was deleted, narrowed or mined for a flaky green.
- One focused selector terminated exit_code: 0 with `[no tests to run]`; it is explicitly not behavioral evidence. The corrected selector was first red, then green, and the final whole-tree/binary/exact gates cover those cases.

## Concerns, technical debt, or notes for the Tech Lead

- CF-047 live-row disagreements and combination semantics are enumerated above for the independent D1 review. This task preserves them.
- VAL-036 is advanced only for message pairing; task 10 remains responsible for its additional Human lines. No claim of full VAL-036 completion is made.
- Headroom: cmd/fid 399.177s, internal/workflow 115.306s, internal/issues 196.455s and tagged integration 340.634s are below their 80% warnings; no headroom warning is triggered.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `2d5e9d55c83086b46268ca69793f30f65031b91b`

## Commit subject

- `feat(cli): converge grammar failures and unknown commands`
