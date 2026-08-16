# Phase 20 Final Review - Standalone Issue Journey

Reviewer: QA Lead
Date: 2026-08-16
Review type: INITIAL
Prior review: none

## Verdict

REPAIR

The seven-command phase gate passes at QA basis
`33f7c669c9fde46fcf95203de43f58e8396688c0`, and 44 of 46 validation
assertions have reproducible passing evidence. The phase must not advance:

- 22 earlier completion reports fail the current canonical report checker, so
  report sequence and command-evidence conformance are not clean (MF-01; no
  VAL).
- Task 18's accepted journey is an ordinary in-process `cmd/fid` test rather
  than a committed real-process test in the dedicated integration package.
  The dedicated gate therefore passes without running the named Phase 20
  journey (MF-02; VAL-036).

Both findings have bounded, decision-free ownership in existing phase tasks.
No protected product, architecture, data-model, public-interface, dependency,
security, validation-contract, or gate-definition change is required. VAL-038
remains pending Human judgment and is not pre-answered by this review.

## Scope And Evidence

- QA basis: branch `phase-20`, starting `HEAD`
  `33f7c669c9fde46fcf95203de43f58e8396688c0`.
- Phase brief: `project-docs/docs/current/phase_brief.md`.
- Validation contract:
  `project-docs/docs/phases/phase-20/qa/validation-contract.md`.
- QA plan: `project-docs/docs/phases/phase-20/qa/qa-plan.md`.
- Accepted slice contract:
  `project-docs/docs/phases/phase-20/contracts.md`.
- Task specs and completion reports:
  `project-docs/docs/phases/phase-20/tasks/` (all 18 numbered tasks and the
  3a-3l, 9a, 16a, and 16b bounded continuations/recoveries).
- Mid-phase reviews: `review-1.md` through `review-6.md`, all bounded review
  artifacts, `infra-review-1.md`, and the infrastructure improvement pack.
- Decisions and authority: phase decisions and triage log, current and archived
  escalation registers, including DEC-130 through DEC-146.
- Current canonical specifications named by the QA plan, `AGENTS.md`,
  `project-docs/docs/playbook.md`, `CHANGELOG.md`, and the live source/command
  tree.
- Held-out scenarios and evaluator prompts were not read.

Before authoring this review, `git status --short` showed only
` M project-docs/docs/phases/phase-20/run-ledger.md`. Its diff is the
Runner-owned Task 18/review-6 append. No tracked file carried `assume-unchanged`
or `skip-worktree`. The sequence verifier returned
`VALID state=ADVANCEABLE` (terminated `exit_code: 0`). The QA Lead did not edit
the ledger.

## Phase-Level Gate Results

Every command below was run exactly from `/home/mark/projects/fid` in this QA
session. The recorded termination is the nested process result, not wrapper
completion.

| Gate | Command | Result | Evidence |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS | No output; token `exit 0 (no output)`; terminated `exit_code: 0`. |
| Lint/static analysis | `golangci-lint run ./...` | PASS | Token `0 issues.`; terminated `exit_code: 0`. |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS | Terminal token `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.6% of statements`; terminated `exit_code: 0`. |
| Coverage report | `go tool cover -func=coverage.out` | PASS | Terminal token `total:										(statements)					78.7%`; terminated `exit_code: 0`. There is no accepted numeric threshold. |
| Integration | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | Token `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.967s`; terminated `exit_code: 0`. MF-02 records that this scope omits the Phase 20 journey. |
| End-to-end | N/A | N/A | No separate end-to-end or held-out-scenario harness is accepted. The accepted Phase 20 journey is required inside the dedicated integration command above. |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS | Terminal token `ok  	github.com/mbell-nzl/fid/internal/workflow	1.977s`; terminated `exit_code: 0`. The long-running nested session was polled through its terminal result. |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | No output; token `exit 0 (no output)`; terminated `exit_code: 0`; `/tmp/fid` is executable. |
| Full phase check | The seven commands above, in QA-plan order | PASS | All seven independently terminated with the recorded statuses. This broad green gate does not close MF-01 or MF-02. |

Selected benchmark measurements from the complete output included:

- `BenchmarkInitializeWorkspace-6 ... 39668174 ns/op ... 4498311 B/op 5735 allocs/op`.
- `BenchmarkList10kIssuesWithRegistry-6 ... 42615771 ns/op`.
- `BenchmarkRebuildWorkspace1kEvents-6 ... 280633642 ns/op`.
- `BenchmarkRebuildWorkspace10kEvents-6 ... 2539097624 ns/op`.
- `BenchmarkRebuildWorkspace100kEvents-6 ... 24660620463 ns/op`.

No benchmark failed or failed to terminate. No accepted performance threshold
was exceeded because the phase contract declares no numeric threshold.

## Fresh Deterministic Evidence Index

The full suite was supplemented with focused matrices and current file-state
proof so that a broad green result is not the sole evidence.

| Evidence | Current-session result | Named tests, files, or literal evidence |
|---|---|---|
| FR-01 envelope, discovery, config, init, authority, lock | PASS | Focused package matrices terminated `exit_code: 0`; named coverage includes `TestWriteSuccessUsesExactContractKeys`, `TestWriteFailureUsesExactContractKeys`, `TestReasonCodeDetailSchemaIsExactForEveryCode`, `TestResolveFindsNearestRootAndRequiresExplicitRootToContainFID`, `TestRootFlagsRejectBeforeLegacyAndUnmigratedLeafHandlers`, `TestProjectConfigRoundTripsStrictTargetSchema`, `TestInitializeFailureBeforeActivationLeavesNoWorkspace`, `TestRebuildJSONMapsUnknownAndDeniedAuthorityWithoutMutation`, and `TestRebuildJSONReportsRetryableHeldLockAndSucceedsAfterRelease`. |
| FR-02 schema, revision, replay, reads, hierarchy | PASS | Focused package matrices terminated `exit_code: 0`; named coverage includes `TestIssuesTableExcludesRetiredPhase20Fields`, `TestReplayJournalClassifiesRetiredIssueEffectKeyAsStopped`, `TestReplayJournalRejectsStaleIssueRevisionBeforeApplyingEffect`, `TestRebuildWorkspaceRecreatesIdenticalStateAfterDeletingDatabase`, `TestIssueGetJSONMatchesClosedSchema`, `TestIssueListJSONMatchesClosedSchemaAndRetainedFilters`, `TestReadyJSONHasExclusiveClosedResultStates`, `TestCreateAllowsLegalHierarchy`, and `TestIssueMutationRequestsCannotChangeParent`. |
| FR-03 concurrency, transition, owner, attribution | PASS | Focused package matrices terminated `exit_code: 0`; named coverage includes `TestIssueTransitionExpectedRevisionMatrix`, `TestIssueTransitionCanonicalOutcomeTable`, `TestIssueTransitionJSONRejectsUnsatisfiedClaimPrecondition`, `TestIssueCreateAndUpdateRejectRemovedOwnerFlagsWithCanonicalFailures`, `TestIssueUpdateCanonicalChangesExcludeOwnerAndPreserveProjection`, `TestIssueMutationCanonicalAttributionAvoidsProtectedEffects`, and `TestPolicyShowMatchesWriteEnforcementForPermittedAndDeniedCells`. |
| FR-04 AC, label, link, workflow, maintenance | PASS | Focused package matrices terminated `exit_code: 0`; named coverage includes `TestAcceptanceCriterionAddJSONMatchesContract`, `TestIssueLabelJSONUsesSpecificAtomicFailureCodes`, `TestLinkJSONDependencyCycleHasExactPathAndNoAppend`, `TestLinkAddUsesSuppliedTypeForAmbiguousSameIDEndpoints`, `TestWorkflowAssignJSONPrecedenceAndScopedReplacement`, `TestWorkflowListEnforcesCollectionLimitWithoutTruncation`, `TestMaintenanceCheckWarningPromotionPreservesFindingsAndInfo`, and `TestCheckStartedWithoutClaimGraceBoundary`. |
| FR-05 legacy and removal | PASS | `TestPolicyShowUsesDirectLegacyJSONWithoutFormatFlag`, `TestJournalInspectCommandDoesNotExposeRevertOrResponseEnvelope`, `TestLiveCobraTreeContainsNoUniversalFormatFlags`, and `TestProductSourceContainsNoUniversalSerializers` passed. Exact no-match/module commands are recorded under Command And Feasibility Validation. |
| FR-06 exact inventories, requests, orders, limits | PASS except accepted-journey integration placement | `cmd/fid/contract_surface_test.go:21` asserts literal 21 retained and 21 legacy leaves and root behavior; `TestLiveLeafCommandsMatchContractInventories`, `TestTask10JSONRequestErrorsUseCanonicalEnvelope`, `TestLinkListJSONHonorsAtAndOverLimitWithoutTruncation`, `TestListOrdersWorkflowsByNameThenID`, `TestWorkflowListEnforcesCollectionLimitWithoutTruncation`, and `TestIssueUpdateJSONOrdersChangesByField` passed. |
| FR-07 ordinary Phase 20 journey | PASS as an in-process package test | `TestP20StandaloneIssueJourneyJSONAndRebuildParity` in `cmd/fid/e2e_p20_issue_journey_test.go:13` passed in the focused `cmd/fid` run and in the repository suite; lines 117-156 perform stale failure, maintenance, rebuild, and raw-byte comparison. This does not substitute for the required dedicated integration test. |
| FR-08 real-process QA journey | PASS as independent manual deterministic evidence | A fresh `/tmp/fid` binary and `/tmp/fid-p20-final.Kr3Aow` workspace completed the 21-leaf representative journey. Token: `P20_REAL_PROCESS workspace=/tmp/fid-p20-final.Kr3Aow retained-journey=complete typed-failure=state-conflict stale-exit=1 event-actor=developer journal-events=17 collection-order=verified rebuild-bytes=identical`. Typed stderr was `{"contract_version":"fid-cli.v1","ok":false,"operation":"issue.update","error":{"code":"state-conflict","category":"conflict","retryable":false,"details":{"actual_revision":9,"expected_revision":8}}}`. |
| FR-09 dedicated integration probe | FAIL | `go test -tags integration -race ./cmd/fid/integration/... -run '^TestP20StandaloneIssueJourneyJSONAndRebuildParity$' -count=1 -v` printed `testing: warning: no tests to run`, `PASS`, and `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.420s [no tests to run]`; terminated `exit_code: 0`. `cmd/fid/integration/classification_test.go:20-25` lists only four pre-Task-18 journeys. |
| FR-10 workflow/report controls | PARTIAL | `project-docs/tools/workflow-lint.sh` returned `PASS workflow-lint completed with no failures`; report-lint self-test returned `SELF-TEST PASS report-lint self-test completed with no failures`; both terminated `exit_code: 0`. Current report-lint passes Tasks 10-18, 16a, and 16b but rejects the 22 reports in MF-01. |
| FR-11 ledger, commits, decisions, escalations | PASS | Sequence lint returned `VALID state=ADVANCEABLE`; all report implementation hashes resolve and subjects/`Refs: phase-20/...` match Git; all reports are committed and clean. DEC-130..DEC-146 and TRG-20-001..015 preserve failed evidence and resolve their authority path. Current and archived escalation registers contain no unresolved ESC; DEC-146 resolves ESC-20-015. |
| FR-12 Human corpus readiness | READY; judgment pending | Fresh workspace `/tmp/fid-p20-human.vAjV9v`, issue `HQA-54265f`: routine commands wrote only stdout; stale update exited 1 with stderr `error: issue "HQA-54265f" revision is 2, not expected revision 1` and no stdout. The Human checklist asks three separate judgments without assigning outcomes. |

## Validation Contract Outcomes

| VAL ID | Outcome | Validator | Evidence |
|---|---|---|---|
| VAL-001 | PASS | QA Lead | FR-01; exact success/failure key tests and the FR-08 decoded real-process envelopes passed. |
| VAL-002 | PASS | QA Lead | FR-01/FR-08; typed code/detail registry, redaction, stdout/stderr separation, and non-zero stale failure were reproduced. |
| VAL-003 | PASS | QA Lead | FR-01; `TestResolveFindsNearestRootAndRequiresExplicitRootToContainFID` and root-selector tests passed. |
| VAL-004 | PASS | QA Lead | FR-01; symlink/boundary discovery and explicit-project rejection matrices passed without handler effects. |
| VAL-005 | PASS | QA Lead | FR-01; root-flag sentinel matrices across retained and legacy leaves passed. |
| VAL-006 | PASS | QA Lead | FR-01/FR-03; missing, unknown, denied, and permitted V1 authority matrices and policy/write parity passed with rejected-effect checks. |
| VAL-007 | PASS | QA Lead | FR-01; held-lock retryable failure and released-lock success tests passed with no rejected mutation. |
| VAL-008 | PASS | QA Lead | FR-01/FR-08; success stdout, typed failure stderr, exit separation, and Human corpus channels were observed. |
| VAL-009 | PASS | QA Lead | FR-02; retired-column schema, Event omission, and removed-key replay-stop tests passed. |
| VAL-010 | PASS | QA Lead | FR-02/FR-04; aggregate revision, sorted link revisions, and unchanged-on-rejection tests passed across the eleven families. |
| VAL-011 | PASS | QA Lead | FR-02/FR-04; golden absolute revision/Prior/Preconditions and stale replay enforcement tests passed. |
| VAL-012 | PASS | QA Lead | FR-02; migration/rebuild and complete projection parity tests passed. |
| VAL-013 | PASS | QA Lead | FR-02; populated/null `issue.get` closed-schema tests passed. |
| VAL-014 | PASS | QA Lead | FR-02/FR-06; list schema, filters, equal-key order, and tie-break tests passed. |
| VAL-015 | PASS | QA Lead | FR-02; exclusive ready result states, closed explanations, sequential and prioritised selectors passed. |
| VAL-016 | PASS | QA Lead | FR-02/FR-05; retired Issue read flags are absent and rejected. |
| VAL-017 | PASS | QA Lead | FR-01/FR-07/FR-08; mutation Event identities matched appended Events and reads returned `event:null`. |
| VAL-018 | PASS | QA Lead | FR-03; absent, matching, and stale expected-revision states passed with stale zero-effect evidence. |
| VAL-019 | PASS | QA Lead | FR-03; legal, unknown-edge, absent-edge, and unsatisfied-precondition transition cases passed, including retained fail-closed V1 claim authority. |
| VAL-020 | PASS | QA Lead | FR-02; legal Epic/Task/Subtask, missing/illegal parent, and no parent-mutation request tests passed. |
| VAL-021 | PASS | QA Lead | FR-04; exact AC schema, revision, order, duplicate, and limit matrices passed. |
| VAL-022 | PASS | QA Lead | FR-04; canonical `colour`, label/Issue-label schema, order, duplicate/missing, and atomicity tests passed. |
| VAL-023 | PASS | QA Lead | FR-04; typed Link add/list/remove, inverse projection, revision, and relation tests passed. |
| VAL-024 | PASS | QA Lead | FR-04; self/two/three-node cycle rejection and acyclic diamond/pre-existing-cycle termination tests passed. |
| VAL-025 | PASS | QA Lead | FR-04; self-link and missing-link removal returned their specific codes with zero append. |
| VAL-026 | PASS | QA Lead | FR-04; workflow schemas, grammar, invalid values, assignment scope/precedence, and forbidden requirement-authoring inputs passed. |
| VAL-027 | PASS | QA Lead | FR-04; closed check registry, fixed 299/300 boundary, and identical warning/promotion findings passed. |
| VAL-028 | PASS | QA Lead | FR-05; journal inspect remains valid unenveloped JSON, has no revert/Envelope contract, and rejects root `--json`. |
| VAL-029 | PASS | QA Lead | FR-05; all eight listed legacy reads decode as semantic unenveloped JSON and reject the canonical root transport. |
| VAL-030 | PASS | QA Lead | FR-05 and exact absence commands: `internal/format` and product imports are absent. |
| VAL-031 | PASS | QA Lead | FR-05/FR-06; complete Cobra walk finds no universal format flags. |
| VAL-032 | PASS | QA Lead | `go mod tidy`, clean module diff, no YAML requirement, exact DEC-145 `go mod why -vendor` stanza, full tests, and build passed. |
| VAL-033 | PASS | QA Lead | FR-05 and serializer no-match command prove no universal TOON/YAML/Markdown serializer remains. |
| VAL-034 | PASS | QA Lead | Stub/promise source no-match and unknown-command tests passed. |
| VAL-035 | PASS | QA Lead | FR-06; literal executable-leaf drift test proves a disjoint 21 retained/21 legacy partition and root transport behavior. |
| VAL-036 | FAIL | QA Lead | FR-09. FR-07 and FR-08 prove the product journey and byte-identical rebuild behavior, but the contract requires an integration journey test. The dedicated integration scope reports `[no tests to run]` for the named Phase 20 journey. |
| VAL-037 | PASS | QA Lead | FR-05/FR-11; live tree matches `AGENTS.md`, playbook lines 257-290, and `CHANGELOG.md:8-28`; Task 18 closes review-6's changelog improvement. |
| VAL-038 | CARRY-FORWARD | Human | HUMAN-038a, HUMAN-038b, and HUMAN-038c are prepared in `human-checklist.md`; all Human outcomes remain blank. Deterministic QA does not pre-answer perception or guidance quality. |
| VAL-039 | PASS | QA Lead | FR-01; exact generated project ID/`issue_prefix` key inventory and strict round-trip passed; old `prefix` is rejected. |
| VAL-040 | PASS | QA Lead | FR-03; create/update owner flags are absent, canonical changes exclude owner, and retained default-dependent safety passes. |
| VAL-041 | PASS | QA Lead | FR-03/FR-08; local attribution is present and representative mutation tests prove no Evidence, Human Outcome, or protected-authority effect. |
| VAL-042 | PASS | QA Lead | FR-01; every injected pre-activation init stage leaves no target/staging residue and success verifies rebuild before activation. |
| VAL-043 | PASS | QA Lead | FR-01; exact TOML sections/keys, signed-64 limits, path/prefix, malformed, relationship, and symlink-escape matrices passed. |
| VAL-044 | PASS | QA Lead | FR-06; exact 21-leaf request/domain tables prove valid boundaries reach handlers and invalid/removed inputs do not. |
| VAL-045 | PASS | QA Lead | FR-04/FR-06/FR-08; declared orders/tie-breaks and at/over-limit no-truncation/no-mutation states passed. |
| VAL-046 | PASS | QA Lead | FR-04; ambiguous identical ID text resolves by explicit endpoint type, typed missing endpoints differ, and failures append nothing. |

Outcome summary: 44 PASS, 1 FAIL, 0 BLOCKED, 1 CARRY-FORWARD. The
CARRY-FORWARD is solely the Human-owned VAL-038 checklist.

## Required Evidence Strategy Outcomes

| Evidence ID | Outcome | Evidence |
|---|---|---|
| EVID-01 | PASS | FR-01, FR-03, FR-04, FR-08, and full/integration suites cover envelopes, channels, redaction, Event identity, and Human corpus setup. |
| EVID-02 | PASS | FR-01 discovery trees and literal root-flag sentinel matrix. |
| EVID-03 | PASS | FR-01 authority and held/released lock tests with state/journal invariance. |
| EVID-04 | PASS | FR-02 retired-field, migration, replay-stop, and complete rebuild parity evidence. |
| EVID-05 | PASS | FR-02/FR-04 revision-family, golden tuple, stale replay, and rebuild evidence. |
| EVID-06 | PASS | FR-02 closed reads, order/filter, ready dual states, retired flags, and hierarchy matrix. |
| EVID-07 | PASS | FR-03 concurrency, transition dual states, owner removal, attribution, and forbidden-effect evidence. |
| EVID-08 | PASS | FR-04 exact AC/label schemas, `colour`, order, revision, limits, and atomic failures. |
| EVID-09 | PASS | FR-04 typed Link sequence, inverse, negatives, graph, and endpoint-resolution matrices. |
| EVID-10 | PASS | FR-04 workflow request/result/grammar, assignment, precedence, and atomic failures. |
| EVID-11 | PASS | FR-04 maintenance registry, 299/300 boundary, identical warning/promotion findings, and channels. |
| EVID-12 | PASS | FR-05 semantic legacy reads and journal-inspect no-revert/no-Envelope evidence. |
| EVID-13 | PASS | FR-05, exact absence/module commands, and full gate prove formatter/flag/YAML/serializer/stub removal. |
| EVID-14 | FAIL | Literal partition, request, order, limits, ordinary journey, FR-08 real process, and byte parity pass; required committed dedicated integration journey is absent (FR-09/MF-02). |
| EVID-15 | PASS | Live-tree/document comparison in FR-11; changelog improvement is closed. |
| EVID-16 | PASS | FR-01 strict config identity/domain and atomic-init matrices. |
| EVID-17 | CARRY-FORWARD | FR-12 corpus is ready; Human must complete three unprejudged rows. |

## Requirement Compliance

| Requirement | Status | Evidence |
|---|---|---|
| Phase 20 goal; PRD-FOUND/ISSUE/WF/OPS/REC requirements mapped by VAL-001..046 | PARTIAL | 44 deterministic assertions pass; VAL-036 fails only its required committed integration-test evidence; VAL-038 is pending Human judgment. |
| `contracts.md` sections 2-7 | MET except VAL-036 evidence placement | FR-01..FR-08 prove envelope, request, identity, schema, order, authority, and journey behavior; FR-09 identifies the dedicated-scope gap. |
| `08_QUALITY_PLAN.md` phase gate and deterministic review | PARTIAL | Seven-command gate passes; QA-01, QA-16, and QA-24 fail for the two Must-fix findings. |
| Documentation synchronization rule | MET | `AGENTS.md`, playbook, and changelog match the live tree; no current formatter/stub promise remains. |

## Deterministic Checklist Execution

| QA ID | Result | Evidence |
|---|---|---|
| QA-01 | FAIL | Exact seven-command gate passed, but fresh canonical report-lint rejects 22 reports. Task 1 ends `FAIL report-lint found 19 failure(s)` (`exit_code: 19`); Task 9a ends `FAIL report-lint found 16 failure(s)` (`exit_code: 16`). See MF-01. |
| QA-02 | PASS | FR-01/FR-03/FR-08; named envelope, reason-code, channel, Event-tail, and read-null Event tests passed. |
| QA-03 | PASS | FR-01; nearest/explicit discovery, boundary/symlink, and retained-vs-legacy root-flag sentinel tests passed. |
| QA-04 | PASS | FR-01/FR-03; policy/write permitted-denied parity, missing/unknown authority, and held/released lock tests passed with rejected-effect counts. |
| QA-05 | PASS | FR-02; retired-field PRAGMA/Event/replay tests and complete live/rebuild projection tests passed. |
| QA-06 | PASS | FR-02/FR-04; revision tests cover all eleven aggregate mutation families, rejection invariance, sorted link effects, golden tuples, stale replay, and rebuild. |
| QA-07 | PASS | FR-02; named closed get/list/ready tests plus hierarchy and no-parent-mutation tests passed. |
| QA-08 | PASS | FR-03; expected-revision matrix, transition outcome table, owner guards, attribution forbidden-effect tests, and fail-closed V1 claim check passed. |
| QA-09 | PASS | FR-04; exact AC/label key, value, `colour`, revision/order, duplicate/missing, and limit tests passed. |
| QA-10 | PASS | FR-04; typed Link add/list/remove, inverse, specific negatives, cycle/diamond, and ambiguous endpoint tests passed. |
| QA-11 | PASS | FR-04; exact workflow schema/grammar, no `--require`, assignment precedence, collection order/limit, and typed failure tests passed. |
| QA-12 | PASS | FR-04; maintenance check registry, fixed 299/300 grace, and identical findings/info across warn success and promoted failure passed. |
| QA-13 | PASS | FR-05; all eight legacy reads and journal inspect semantic/no-revert/no-Envelope/root-flag tests passed. |
| QA-14 | PASS | Exact directory/import/serializer/YAML/stub absence commands have their declared real non-zero statuses; tidy/module graph/full test/build pass. |
| QA-15 | PASS | FR-06; `TestLiveLeafCommandsMatchContractInventories` and literal source lists prove exactly 42 leaves partition 21/21. |
| QA-16 | FAIL | FR-09/MF-02: the named Task 18 test is not in `cmd/fid/integration`; a focused dedicated-package probe prints `[no tests to run]`. FR-07/FR-08 do not satisfy the accepted committed integration-test method. |
| QA-17 | PASS | `AGENTS.md`, playbook lines 257-290, changelog lines 8-28, and live source/command tree agree; review-6's Task 18 changelog improvement is closed. |
| QA-18 | PASS | FR-01; strict config round-trip/boundary/symlink tests and `TestInitializeFailureBeforeActivationLeavesNoWorkspace` passed. |
| QA-19 | PASS | FR-06; source tables at `contract_surface_test.go:21-153` plus request/domain sentinel matrices cover all 21 retained leaves and both valid/invalid states. |
| QA-20 | PASS | FR-04/FR-06; list/link/label/workflow/update order and at/over-limit tests passed with no silent truncation or over-limit mutation. |
| QA-21 | PASS | FR-01/FR-03/FR-04: redaction, path/symlink containment, lock atomicity, pre-activation residue, stale/invalid request, and rejected journal/effect tests passed. |
| QA-22 | PASS | Current benchmark session reached real `exit_code: 0`; 1k/10k/100k rebuild, 10k registry, and init benchmark tokens are recorded above. |
| QA-23 | PASS | Task reports preserve red/green and immutable deviations; focused behavior tests use `T.TempDir`/temporary DBs and cover both states. Review history records no weakened semantic assertion. The remaining report-format defect is separately QA-01/QA-24. |
| QA-24 | FAIL | All six review-window findings and ESCs are historically resolved, all hashes/subjects/Refs and dependency notes reconcile, but 22 current completion reports fail canonical report-lint (MF-01). |

Checklist summary: 21 PASS, 3 FAIL, 0 BLOCKED. The failures are QA-01,
QA-16, and QA-24.

## Command And Feasibility Validation Result

- Command and feasibility validation completed before implementation: YES.
- All seven phase-gate commands are concrete and executable.
- Invalid or unvalidated command, test scope, fixture, artifact, or Human step
  found: the QA-plan assumption at lines 159 and 201 that Task 18's journey is
  committed in the dedicated integration package is false in the current tree
  (MF-02). Human steps are feasible from the fresh FR-12 corpus.

Additional exact current-session checks:

| Command or step | Validation method | Result | Evidence |
|---|---|---|---|
| `go mod tidy` | Executed at QA basis | PASS | No output; terminated `exit_code: 0`. |
| `git diff --exit-code -- go.mod go.sum` | Executed after tidy | PASS | No output; terminated `exit_code: 0`. |
| `ls internal/format` | Expected absence | PASS | `ls: cannot access 'internal/format': No such file or directory`; expected and observed `exit_code: 2`. |
| `rg -n "internal/format" cmd internal` | Expected no-match | PASS | `exit 1 (no output)`; expected and observed `exit_code: 1`. |
| `rg -ni "writeTOON|writeYAML|writeMarkdown|yaml.Marshal" cmd internal` | Expected no-match | PASS | `exit 1 (no output)`; expected and observed `exit_code: 1`. |
| `rg -ni "yaml" go.mod` | Expected no-match | PASS | `exit 1 (no output)`; expected and observed `exit_code: 1`. |
| `go mod why -vendor -m gopkg.in/yaml.v3` | Exact DEC-145 stanza | PASS | `# gopkg.in/yaml.v3` then `(main module does not need to vendor module gopkg.in/yaml.v3)`; terminated `exit_code: 0`. |
| `rg -n "newDecisionCommand|newEscalationCommand|newTriageCommand|versionGateError|ships in v0\\.2" cmd/fid` | Expected no-match | PASS | `exit 1 (no output)`; expected and observed `exit_code: 1`. |
| `go test -tags integration -race ./cmd/fid/integration/... -run '^TestP20StandaloneIssueJourneyJSONAndRebuildParity$' -count=1 -v` | Named-test scope probe | FAIL | `[no tests to run]`; terminated `exit_code: 0`; the command succeeded without exercising the required test. |
| `project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-20/run-ledger.md` | Current append-only ledger verification | PASS | `VALID state=ADVANCEABLE`; terminated `exit_code: 0`. |
| `project-docs/tools/report-lint.sh --self-test` | Canonical checker self-test | PASS | `SELF-TEST PASS report-lint self-test completed with no failures`; terminated `exit_code: 0`. |
| Task 18 report-lint invocation | Exact report/spec pair | PASS | `OUTCOME COMPLETE` and `PASS report-lint completed with no failures`; terminated `exit_code: 0`. |

## Completion Report, Commit, And Sequence Audit

- All task specs and all 33 completion reports were read. Every report is
  committed and clean; the only pre-review dirt is Runner-owned ledger dirt.
- Every implementation hash resolves to a commit. Its reported subject and
  phase/task reference match Git. Task 18 records product commit
  `37b59d629b296b7b5383454170f9f017fdd94253`, subject
  `test(cli): add Phase 20 issue journey acceptance`.
- Every report contains a non-empty dependency-budget note. No task adds a
  dependency. Task 16 removes the direct YAML dependency and records no new or
  upgraded dependency.
- Historical failed and UNVERIFIED evidence under DEC-130/131 and later
  recoveries remains present; this review does not rewrite it.
- Current canonical report-lint passes Tasks 10-18 and support Tasks 16a/16b.
  It rejects Tasks 1, 2, 3, 3a, 3b, 3c, 3d, 3e, 3f, 3g, 3h, 3i, 3j, 3k, 3l,
  4, 5, 6, 7, 8, 9, and 9a. The typed workflow outcome requirement entered the
  canonical control mid-phase in commit
  `615eda1` (`feat(workflow): type report recovery outcomes`); the current
  checker contains no grandfather rule. This provenance explains the defect
  but does not make the current reports conformant.
- Representative fresh results: Task 1 reports 19 failures and terminates
  `exit_code: 19`; Task 9a reports 16 failures and terminates `exit_code: 16`;
  Task 18 returns `OUTCOME COMPLETE` and terminates `exit_code: 0`.
- The accepted ledger reconciliation for Task 8 is preserved under DEC-139.
  The current schema-v2 append and recovery identities are valid according to
  the sequence verifier. QA did not allocate an escalation or edit the ledger.

## Mid-Phase Review Disposition

| Review window | Outcome/history | Must-fix disposition | Fresh independent evidence | Blocker candidate / existing ESC |
|---|---|---|---|---|
| Review 1 | `review-1.md` BLOCK; `review-1b.md` BLOCK; `review-1c.md` CONTINUE | Closed by accepted planning/config/integration recoveries | `review-1c.md`; FR-01/FR-12; current full gate | ESC-20-001..003 resolved |
| Infrastructure review | Initial BLOCK; appended NON-BLOCK correction | Closed under DEC-133/134 and accepted infra improvement pack | `infra-review-1.md` correction; dedicated gate/current sequence lint pass | ESC-20-004 resolved |
| Review 2 | CONTINUE (NON-BLOCK) | No Must-fix; one low Improvement remains | `review-2.md`; FR-02; full gate | none |
| Review 3 | `review-3.md` BLOCK; `review-3b.md` CONTINUE | Task 9a closed owner/attribution/concurrency evidence | `review-3b.md`; FR-03 | ESC-20-008 resolved |
| Review 4 | `review-4.md` REPAIR | MF-11/MF-12 closed by fresh `review-4b.md` and `review-4c.md` | FR-04 plus accepted re-reviews | none open |
| Review 5 | `review-5.md` REPAIR; `review-5c.md` REPAIR | Task 13 closed in `review-5b.md`; Task 15 closed in fresh `review-5d.md` after first repair remained insufficient | FR-04/FR-05 and resolving reviews | none open |
| Review 6 | CONTINUE | No Must-fix; changelog Improvement assigned to Task 18 and now closed | `CHANGELOG.md:8-28`, FR-11 | ESC-20-015 resolved by DEC-146 |

No original mid-phase REPAIR is treated as closed by implementation evidence
alone; each cites its fresh accepted re-review. The final-review findings below
are newly identified current-conformance gaps.

## Documentation Review

- Impacted docs updated: YES. `project-docs/docs/playbook.md:257-290` describes
  the delivered JSON surface and Human mode. `CHANGELOG.md:8-28` records the
  Phase 20 journey/removals and explicitly marks the 0.1.0 formatter/stub text
  as historical. `AGENTS.md` matches the live package tree.
- Task 18 changelog improvement closure: VERIFIED.
- Specs needing Tech Lead wrap-up sync: none for product behavior. The
  `docs/current` phase status remains pre-Task-18 and is normal wrap-up context;
  the user prohibited QA edits to current context.
- Completion reports present for all tasks: YES; current canonical conformance
  is NO because of MF-01.

## Human Checklist Readiness

Generated at
`project-docs/docs/phases/phase-20/qa/human-checklist.md`.

- Contains only `(human)` items: YES.
- Uses one observable per row: YES.
- Human-only validation assertion mapped: YES, VAL-038 through HUMAN-038a,
  HUMAN-038b, and HUMAN-038c.
- Tester, outcomes, notes, summary counts, and overall result are pending.
- Product/perception risks: Human must decide whether routine output is concise,
  the stale-revision failure guidance is actionable, and the corpus avoids
  advising direct edits to fid internals.

## Coverage Caveat

The deterministic evidence proves only the asserted mechanics, data states,
command outputs, and structural conformance checks. It does not prove workflow
feel, product-owner judgment, tone, or Human guidance quality. VAL-038 remains
unconfirmed until a Human records all three outcomes.

## Release Maturity Status

| Maturity area | Status | Evidence or deferral |
|---|---|---|
| Baseline artifacts | BLOCKED | Required artifacts exist, but MF-01 and MF-02 prevent a clean phase baseline; Human VAL-038 is pending. |
| Automated gates | PARTIAL | Seven-command phase gate and deterministic security/data-loss tests pass; accepted Task 18 integration scope is incomplete. |
| Held-out scenarios | DEFERRED | The specs declare no accepted held-out harness for Phase 20; none was read. |
| Release readiness | DEFERRED | SAST, vulnerability audit, license scan, secret scan, threat model, and public-release checklist remain pre-release maturity work under `AGENTS.md`/quality plan. |

## Findings

### Must-fix

#### MF-01 - Twenty-two completion reports fail the current canonical report checker

- Category: report/control conformance; no VAL.
- Evidence: current `project-docs/tools/report-lint.sh` self-tests clean, but
  rejects the exact report/spec pairs for Tasks 1, 2, 3, 3a-3l, 4-9, and 9a.
  Task 1 terminates `exit_code: 19`; Task 9a terminates `exit_code: 16`.
  Failures include the missing typed `## Workflow outcome`, legacy evidence-row
  punctuation/termination syntax, and retained UNVERIFIED command evidence.
- Impact: QA-01 and QA-24 fail. The phase's report sequence is not clean under
  the current canonical control even though hashes, subjects, refs, report
  committed state, and historical authority are valid.
- Ownership: bounded same-task report-only continuations for each affected task.
  This is not authority to fabricate evidence or alter the checker.

#### MF-02 - Task 18's accepted journey is absent from the dedicated integration package

- Category: required test-scope/evidence omission; VAL-036, EVID-14, QA-16.
- Evidence: `task-18.md:29` calls
  `cmd/fid/e2e_p20_issue_journey_test.go` an integration artifact;
  `qa-plan.md:159` and `qa-plan.md:201` require it in the dedicated real-process
  suite. The actual test is an untagged in-process `package main` test at
  `cmd/fid/e2e_p20_issue_journey_test.go:13`.
  `cmd/fid/integration/classification_test.go:20-25` still names only four
  earlier journeys. The exact named integration probe prints
  `[no tests to run]` while exiting 0.
- Impact: the full integration gate is broadly green but cannot prove the
  Task 18 journey ran. FR-07 and FR-08 show the product behavior and rebuild
  parity work; they do not satisfy the accepted committed integration method.
- Ownership: Task 18, test and report evidence only. No product behavior,
  contract, gate command, or accepted pass/fail semantics need change.

### Improvement

#### IMP-01 - Direct invalid-input tests for `BuildIssueRevisionEffects`

- `review-2.md:99-103` identifies untested empty input, blank ID,
  non-positive revision, and duplicate-ID branches.
- Current `internal/issues/revision_test.go:5-37` still covers only the
  successful sorted two-Issue result. This is a low, deferrable test-quality
  improvement and does not invalidate the aggregate revision/replay matrices.

## Prior Finding Resolution

Not applicable: this is the initial phase-end final review. Historical
mid-phase findings and their fresh review disposition are recorded above.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-01 | Tasks 1, 2, 3, 3a-3l, 4-9, 9a | no VAL | Same-task report-only continuations for the 22 named `*-done.md` files. Add the canonical typed workflow outcome and normalize evidence syntax only where an actual observed result exists. Preserve every failed/UNVERIFIED attempt and deviation. Where a required command lacks valid session-captured termination, rerun that exact task-spec command and record the new result; never transcribe or fabricate. Do not edit task specs, product/tests, report-lint, decisions, QA contract/plan, or ledger. | Each of the 22 exact report/spec `report-lint.sh` invocations terminates `exit_code: 0` with `OUTCOME COMPLETE` and `PASS report-lint completed with no failures`; hashes/subjects/Refs and committed-clean report state still pass; any rerun command has its own current-session token and termination. | Fresh independent QA Lead reruns all 22 report-lint pairs, checks immutable failures/deviations and dependency notes, verifies report commits/subjects/Refs and sequence/ledger state, and reruns any phase gate whose evidence a repair invalidated. |
| MF-02 | Task 18 | VAL-036 | Add or move a Phase 20 real-process journey into `cmd/fid/integration/`, using the package-owned built binary and temporary workspace. Update the classifier's exact required journey inventory. The journey must retain canonical JSON steps, specified typed failure, Event identities, attribution boundary, collection order/limits, and byte-identical pre/post-rebuild Issue reads. Retain the ordinary test if useful. Update only Task 18's completion report to preserve existing red/green/manual evidence and record new integration evidence. No product behavior, contract, gate command, or pass semantics change. | The exact named integration test runs (not `[no tests to run]`) and passes; `TestIntegrationClassification` requires it; full `go test -tags integration -race ./cmd/fid/integration/...`, all seven Task 18 commands, and Task 18 report-lint pass with actual terminations. | Fresh independent QA Lead reruns the named test, classifier, complete dedicated integration gate, Task 18's seven exact commands, Task 18 report-lint, QA-16, EVID-14, and VAL-036; inspect raw-byte rebuild operands and the updated report. |

## Blocker Candidates

- None. No unresolved escalation or protected blocker remains. DEC-130 through
  DEC-146 and the escalation registers were reviewed; DEC-146 resolves the last
  registered ESC-20-015.

## Unresolved Risks And Carry-Forward Candidates

- VAL-038 is intentionally carried to Human execution; its outcome is unknown.
- IMP-01 remains a deferrable direct-test gap.
- Current-context status text is pre-Task-18 and requires ordinary Tech Lead
  wrap-up compression after repairs; it does not contradict the live product
  documentation and was outside QA write authority.
- No held-out or public-release maturity controls are accepted yet; this review
  does not infer release readiness from the phase gate.

## Final Decision

```text
Must-fix: 2  Improvement: 1  Blocker candidates: 0
VERDICT: REPAIR
```
