# Phase 11 Final Review - Claims Layer

Reviewer: QA Lead
Date: 2026-06-21

## Verdict

PASS WITH HUMAN CHECKLIST PENDING.

Deterministic QA passed. `VAL-001` through `VAL-013` and `VAL-015` are PASS.
`VAL-014` is CARRY-FORWARD to the generated human checklist because it requires
human dogfooding and product judgment. ESC-008 is resolved, not open: DEC-028
required Task 5 and VAL-015 evidence before final QA, and Review 5 plus the
final gate confirm the stale-projection claim-write fix.

## Scope And Evidence

- Phase brief: `project-docs/docs/current/phase_brief.md`
- Validation contract:
  `project-docs/docs/phases/phase-11/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-11/qa/qa-plan.md`
- Task completion reports:
  `project-docs/docs/phases/phase-11/tasks/task-1-done.md` through
  `task-5-done.md`
- Mid-phase reviews:
  `project-docs/docs/phases/phase-11/qa/review-1.md` through `review-5.md`
- Decisions and escalations:
  `project-docs/docs/phases/phase-11/decisions.md`,
  `project-docs/docs/phases/phase-11/tl-escalations.md`,
  `project-docs/docs/current/tl_escalations_current.md`
- Disposable CLI fixture workspace:
  `/tmp/fid-phase11-final-qa.ujFUN4`

## Phase-Level Gate Results

| Gate | Command | Result | Evidence |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS | Command exited 0 with no output. |
| Lint/static analysis | `golangci-lint run ./...` | PASS | Output: `0 issues.` |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS | Output included `ok github.com/mbell-nzl/fid/cmd/fid (cached) coverage: 77.8%`, `internal/claims 81.4%`, `internal/issues 81.8%`, `internal/journal 80.0%`, `internal/policy 81.2%`, `internal/workflow 95.5%`. |
| Coverage report | `go tool cover -func=coverage.out` | PASS | Output ended with `total: (statements) 76.2%`. |
| Integration | `go test -tags integration -race ./...` | PASS | Output included `ok github.com/mbell-nzl/fid/cmd/fid (cached)` and PASS/ok for `internal/claims`, `internal/issues`, `internal/journal`, `internal/policy`, and `internal/workflow`. |
| End-to-end | N/A | N/A | `AGENTS.md` and `08_QUALITY_PLAN.md` define end-to-end/scenario validation as N/A until a scenario harness is accepted. |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS | Output included `BenchmarkRebuildWorkspace1kEvents-6 377.9 ms/rebuild`, `BenchmarkRebuildWorkspace10kEvents-6 3616 ms/rebuild`, and PASS for all packages. |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | Command exited 0 with no output. Binary size check: `wc -c /tmp/fid` returned `12134506 /tmp/fid`, below TRD-PERF-007 30 MB. |
| Full phase check | Full command list from QA plan section 5 | PASS | Every exact phase command above passed in this final review run. |

## Validation Contract Outcomes

| VAL ID | Outcome | Validator | Evidence |
|---|---|---|---|
| VAL-001 | PASS | QA Lead | `go test -v ./internal/policy ./cmd/fid -run 'TestMatrixContainsExactlyTheImplementedV01WriteActions|TestMatrixClaimRowsMatchSchemaSection13|TestIssueClaimCommandMatchesPolicyShowForForceAcquire|TestPolicyShow' -count=1` passed. Output included subtests for `issue.claim.acquire`, `issue.claim.renew`, `issue.claim.force_acquire`, and `issue.claim.release`. CLI fixture `fid policy show --format markdown` printed the four rows with schema-matching booleans. |
| VAL-002 | PASS | QA Lead | Combined Task 1 domain and Task 2 CLI evidence. Focused tests passed: `TestAcquireCreatesLiveIssueClaimWithDefaultTTL`, `TestAcquireSupportsCustomTTLAndNote`, `TestIssueClaimCommandCreatesLiveClaimWithDefaultTTL`, and `TestIssueClaimCommandSupportsCustomTTLNoteAndHumanReadBack/custom_ttl_and_note`. CLI read-back excerpt: `[{"ref_type":"issue","ref_id":"FID-PHASE11-FINAL-QA.UJFUN4-74e2a2","claimed_by":"developer","claimed_at":"2026-06-21T05:59:34Z","ttl_seconds":120,"expires_at":"2026-06-21T06:01:34Z","note":"qa fixture"}]`. |
| VAL-003 | PASS | QA Lead | Combined Task 1 domain and Task 2 CLI evidence. Focused tests passed: `TestAcquireRejectsLiveConflictBeforeAppend` and `TestIssueClaimCommandRejectsLiveConflictBeforeAppend`. CLI conflict excerpt: `error: issue "FID-PHASE11-FINAL-QA.UJFUN4-74e2a2" already claimed by "developer" until 2026-06-21T06:01:34Z`. Reviews 1-2 cite unchanged journal line and DB row counts. |
| VAL-004 | PASS | QA Lead | Combined Task 1 domain superseded-holder evidence and Task 2 CLI stale-takeover evidence. Focused tests passed: `TestAcquireTakesOverStaleClaimWithoutForcePermission` and `TestIssueClaimCommandTakesOverStaleClaimWithoutForce`. Review 1 cites `claim_superseded` activity evidence naming prior holder and expiry; Review 2 cites CLI current-holder read-back after stale takeover. |
| VAL-005 | PASS | QA Lead | Focused tests passed: `TestAcquireForceSupersedesLiveClaimForPermittedRole`, `TestAcquireForceRejectsDeniedRoleBeforeAppend`, `TestIssueClaimCommandForceAcquireReportsPreviousHolder`, `TestIssueClaimCommandRejectsDeniedForceBeforeAppend`, and policy-show parity tests. CLI excerpts: denied force returned `role "developer" not permitted for action "issue.claim.force_acquire"`; permitted force returned `force-acquired claim ... by tl (previously held by developer until 2026-06-21T06:01:34Z)`. |
| VAL-006 | PASS | QA Lead | Focused tests passed: `TestRenewUpdatesClaimTimePreservesOrOverridesTTLAndProjectsActivity`, `TestRenewRejectsNonHolderBeforeAppend`, `TestIssueClaimCommandRenewsHolderClaimAndPreservesTTL`, and `TestIssueClaimCommandRenewsWithOverrideTTLAndRejectsNonHolder`. CLI renewal excerpt: `renewed claim on FID-PHASE11-FINAL-QA.UJFUN4-74e2a2 by tl (expires 2026-06-21T06:01:39Z)`. Reviews 1-2 cite activity-log projection and zero-append non-holder rejection. |
| VAL-007 | PASS | QA Lead | Focused tests passed: `TestReleaseDeletesHolderClaim`, `TestReleaseRejectsNonHolderBeforeAppend`, `TestIssueUnclaimCommandReleasesHolderClaim`, and `TestIssueUnclaimCommandRejectsWrongHolderRepeatAndRoleErrorsBeforeAppend` with wrong-holder, repeat, missing-role, and unknown-role subtests. CLI excerpts: `released claim on ... by tl`; post-release claims read-back returned `[]`. |
| VAL-008 | PASS | QA Lead | Focused CLI tests passed: `TestClaimListDefaultsToJSONAndListsLiveClaims` and `TestClaimListDispatchesFormatsAndFiltersByStoredRole` with toon/yaml subtests; Task 3 report and Review 3 cite markdown parity. CLI JSON excerpt included all required fields: `ref_type`, `ref_id`, `claimed_by`, `claimed_at`, `ttl_seconds`, `expires_at`, `is_stale`, and `note`. Markdown excerpt showed `stale: false` plus note text. |
| VAL-009 | PASS | QA Lead | Focused tests passed: `TestListFiltersClaimsByRoleAndStalenessAndOrdersByRowID`, `TestClaimListStaleFilterIsReadOnlyAndUsesRowIDOrder`, and `TestClaimListOrderIsStableAcrossRebuild`. Review 3 cites live exclusion, stale inclusion, rowid order, and identical order after `fid rebuild`. |
| VAL-010 | PASS | QA Lead | Evidence combined across Tasks 2-4. Task 2 proves `--by human` persists/read-backs as `human`; Task 3 proves human-held claims list as `claimed_by: human`; Task 4 proves developer transition is blocked while human holds the live claim. CLI excerpt: `claimed ... by human`; developer transition rejected with `live claim not held by "developer"; held by "human" until 2026-06-21T06:05:30Z`. |
| VAL-011 | PASS | QA Lead | Focused tests passed: `TestTransitionIntoStartedRequiresLiveCallerClaim`, `TestTransitionIntoStartedRejectsWrongHolderAndStaleClaims`, and CLI `TestIssueTransitionRequiresLiveCallerClaimForStartedStates` with no-claim, wrong-holder, stale-caller, and human-holder subtests. Reviews 4 cites unchanged journal/status proof for rejects and matching-holder success. |
| VAL-012 | PASS | QA Lead | Focused tests passed: `TestTransitionOutOfStartedReleasesCallerClaimInSameEvent` and CLI `TestIssueTransitionReleasesCallerClaimWhenLeavingStartedState`. CLI fixture transitioned `human` from `todo` to `in_progress`, then `in_progress` to `cancelled`; post-transition claim read-back returned `[]`. Review 4 cites one event with issue update, claim delete, and activity effect. |
| VAL-013 | PASS | QA Lead | Task 1 and Task 5 both preserve concurrent acquire atomicity. Focused test passed: `TestAcquireIsAtomicAcrossConcurrentAttempts`; Task 5 report also records `go test -race ./internal/claims -run TestAcquireIsAtomicAcrossConcurrentAttempts -count=10` PASS. Reviews 1 and 5 cite exactly one success, loser holder/expiry conflict, and one current row. |
| VAL-014 | CARRY-FORWARD | Human | Human-only dogfood validation. Generated checklist at `project-docs/docs/phases/phase-11/qa/human-checklist.md` fans this into HV-01a through HV-01e. QA Lead did not mark this PASS because product comprehension and dogfooding require human judgment. |
| VAL-015 | PASS | QA Lead | ESC-008 closure evidence. Task 5 and Review 5 show acquire, renew, and release use lock-aware projection freshness before validation. Focused tests passed: `TestAcquireRejectsLiveConflictWhenJournalIsAheadOfProjection`, `TestRenewCatchesUpJournalAheadClaimBeforeValidation`, and `TestReleaseCatchesUpJournalAheadClaimBeforeValidation`. Review 5 cites zero append on stale-projection acquire conflict, exactly one event for stale-projection renew/release successes, no nested-lock error, and unchanged concurrent acquire behavior. DEC-028 is satisfied and ESC-008 is resolved, not open. |

## Requirement Compliance

| Requirement | Status | Evidence |
|---|---|---|
| PRD-CLAIM-001 | MET | VAL-002 PASS: default TTL, custom TTL, note, holder, expiry, single current row. |
| PRD-CLAIM-002; PRD-ERR-003 | MET | VAL-003 PASS: live conflict rejects with holder/expiry and no append/projection mutation. |
| PRD-CLAIM-003 | MET | VAL-004 PASS: stale takeover succeeds without force and records prior holder/expiry. |
| PRD-CLAIM-004; TRD-ARCH-010 | MET | VAL-005 and VAL-007 PASS: force-acquire permission and release enforcement. |
| PRD-CLAIM-005 | MET | VAL-011 PASS: started entry requires caller-held live claim. |
| PRD-CLAIM-006 | MET | VAL-012 PASS: started exit auto-releases claim in same event/transaction. |
| PRD-CLAIM-007 | MET | VAL-008 and VAL-009 PASS: claim listing fields, formats, `--by`, stale filter, and rowid ordering. |
| PRD-CLAIM-008 | MET deterministically; human judgment pending | VAL-010 PASS for mechanics; VAL-014 CARRY-FORWARD for human comprehension and dogfooding. |
| PRD-CLAIM-008a; CF-011 Phase 11 slice | MET | VAL-006 PASS: renewal projects `activity_log`; broader CF-011 remains deferred by decision. |
| PRD-OPS-003; PRD-NFR-008; `02_TRD.md` sections 4.3 and 4.9 | MET | VAL-015 PASS and DEC-028 satisfied. |
| TRD-ARCH-011 | MET | VAL-013 and VAL-015 PASS: concurrent acquire atomicity and stale-projection lock ordering. |
| 08_QUALITY_PLAN.md phase gate | MET with human checklist pending | Full deterministic gate passed; human checklist generated from `(human)` subset only. |

## Deterministic Checklist Execution

| QA ID | Result | Evidence |
|---|---|---|
| QA-01 | PASS | `rg "gofmt -l cmd internal|golangci-lint run ./...|go test -race -coverprofile=coverage.out ./...|go tool cover -func=coverage.out|go test -tags integration -race ./...|go build -o /tmp/fid ./cmd/fid|go test -bench . -benchmem ./..." project-docs/docs/phases/phase-11/tasks/task-*-done.md` found each required task command. Task 1 and Task 2 benchmark/build omissions were explicitly N/A per conditional scope; Task 4 benchmark was N/A with reason; Task 5 includes benchmark and build PASS. |
| QA-02 | PASS | Final exact gates all passed: `gofmt -l cmd internal` no output; `golangci-lint run ./...` `0 issues.`; `go test -race -coverprofile=coverage.out ./...` PASS; `go tool cover -func=coverage.out` total 76.2%; integration PASS; benchmarks PASS; build PASS. |
| QA-03 | PASS | Package coverage from final race coverage command: `cmd/fid` 77.8% (>=70), `internal/issues` 81.8%, `internal/claims` 81.4%, `internal/workflow` 95.5%, `internal/policy` 81.2%, `internal/journal` 80.0% (domain packages all >=80). |
| QA-04 | PASS | Policy focused tests passed and `fid policy show --format markdown` printed `issue.claim.acquire`, `issue.claim.renew`, `issue.claim.force_acquire`, and `issue.claim.release` rows. |
| QA-05 | PASS | Focused acquire tests passed; CLI fixture read-back showed `ttl_seconds":120`, `claimed_by":"developer"`, expiry, and `"note":"qa fixture"`. Task reports cover default TTL 1800. |
| QA-06 | PASS | CLI conflict returned holder/expiry error. Reviews 1-2 cite before/after journal line count and DB row count. |
| QA-07 | PASS | Focused stale takeover tests passed; Review 1 cites `claim_superseded` activity with prior holder/expiry; Review 2 cites CLI stale takeover current-holder read-back. |
| QA-08 | PASS | Denied force CLI names `issue.claim.force_acquire`; permitted force CLI reports previous holder and new `tl` holder; focused allow/deny tests passed. |
| QA-09 | PASS | Renewal focused tests passed for preserve/override TTL and non-holder rejection; CLI renewal output captured; Reviews 1-2 cite `activity_log` evidence. |
| QA-10 | PASS | Release focused tests passed for holder release and repeat/wrong/missing/unknown role rejects; CLI release output and `[]` claim read-back captured. |
| QA-11 | PASS | Claim-list focused tests passed for json/toon/markdown/yaml parity and `--by`; CLI JSON and markdown excerpts captured. |
| QA-12 | PASS | Stale/list ordering focused tests passed, including stale-only membership, live exclusion, rowid order, and rebuild-stable order. |
| QA-13 | PASS | Human claim mechanics proven across CLI fixture and focused tests; developer transition blocked by human-held claim with holder/expiry in the error. |
| QA-14 | PASS | Transition focused tests passed for no claim, wrong holder, stale holder, and matching holder; Reviews 4 cites unchanged journal/status proof. |
| QA-15 | PASS | Started-exit focused tests passed; CLI fixture post-exit claim read-back returned `[]`; Review 4 cites same-event issue update plus claim delete. |
| QA-16 | PASS | `TestAcquireIsAtomicAcrossConcurrentAttempts` passed in final focused run; Task 5 report records `-race -count=10` PASS. |
| QA-17 | PASS | Completion reports record docs updated or verified-no-change. Task 5 verified no changes required for `02_TRD.md`, `03_TRD_SCHEMA.md`, and `06_DESIGN_PLAN_CLI.md`; Task 2/3/4 verified no changes needed for claim CLI/design docs. Phase wrap-up spec-sync remains outside this final QA request. |
| QA-18 | PASS | `git diff --name-only main...HEAD -- migrations` exited 0 with no output. No schema migration was added; DEC-026 and DEC-028 both say no migration planned. |
| QA-19 | PASS | `git log --oneline -- project-docs/docs/phases/phase-11/tasks/task-*-done.md` showed each implementation commit and reconciliation commit. `git rev-parse --verify` succeeded for reported hashes: `2a8a3aac...`, `fbb8c2f...`, `47c0f86...`, `75826699...`, and `de6a73c...`. `git show --no-patch` confirmed subjects and `Refs: phase-11/task-N` bodies. |
| QA-20 | PASS | ESC-008 is closed by VAL-015 evidence. `project-docs/docs/phases/phase-11/tl-escalations.md` marks ESC-008 Resolved; DEC-028 requires Task 5 and VAL-015; Review 5 records stale-projection acquire/renew/release proof and no nested-lock regression. |

## Command And Feasibility Validation Result

- Command and feasibility validation completed before implementation: YES.
- Invalid or unvalidated command, test scope, fixture, artifact, or human
  checklist step found: none.
- Placeholder command scan:
  `rg "<TBD>|<E2E_COMMAND>|placeholder|TODO" AGENTS.md project-docs/specs/08_QUALITY_PLAN.md project-docs/docs/phases/phase-11/tasks project-docs/docs/phases/phase-11/qa/validation-contract.md`
  found only the standing AGENTS.md placeholder rule, not an unresolved Phase
  11 command placeholder.
- Branch check: `git branch --show-current` returned `phase-11`.

## Mid-Phase Review Disposition

| Review | Verdict | Blocking Fixed? | Improvement Disposition | ESC References |
|---|---|---|---|---|
| `review-1.md` | CONTINUE | N/A | Policy-show evidence gap closed by final policy tests and CLI `policy show` excerpt. | none |
| `review-2.md` | CONTINUE | N/A | Timing-sensitive CLI stale test noted as future improvement; deterministic tests passed. | none |
| `review-3.md` | CONTINUE | N/A | Sparse task report output supplemented by final gate and focused test output. | none |
| `review-4.md` | CONTINUE | N/A | VAL-010 split-window wording corrected by combined final evidence from Tasks 2-4. | ESC-008 recent resolution only |
| `review-5.md` | CONTINUE | N/A | No improvements. ESC-008 resolved via DEC-028/VAL-015. | ESC-008 resolved |

## Documentation Review

- Impacted docs updated or verified-no-change: YES.
- Completion reports present for all tasks: YES.
- QA plan and Reviews 1-5 present: YES.
- Specs needing Tech Lead wrap-up sync: phase wrap-up spec-sync table in
  `project-docs/docs/phases/phase-11/decisions.md` is still pending by design;
  this final QA did not perform phase wrap-up duties.
- Documentation caveat: the current worktree had pre-existing uncommitted and
  untracked Phase 11 documentation artifacts before this review. This final QA
  added only `final-review.md` and `human-checklist.md`.

## Human Checklist Readiness

Generated at
`project-docs/docs/phases/phase-11/qa/human-checklist.md`.

- Contains only `(human)` items: YES.
- Uses one observable per row: YES.
- Human-only `VAL-NNN` assertions mapped: YES, all rows cite VAL-014.
- Product/perception/dogfooding risks to call out: human must judge whether
  holder, expiry, release, block reason, and claim lifecycle messages are clear
  in a real `.fid` workspace.

## Coverage Caveat

A deterministic PASS proves the asserted mechanics, data states, command
outputs, and structural conformance checks. It does not prove claim lifecycle
message clarity, product-owner judgment, or dogfooding quality; those remain
unconfirmed until the human checklist records outcomes.

## Release Maturity Status

| Maturity area | Status | Evidence or deferral |
|---|---|---|
| Baseline artifacts | READY WITH HUMAN PENDING | Validation contract, QA plan, task reports, reviews, final review, and human checklist exist. VAL-014 remains human-owned. |
| Automated gates | READY | Format, lint, race coverage, coverage report, integration, benchmarks, and build passed. |
| Held-out scenarios | DEFERRED | N/A per AGENTS.md and `08_QUALITY_PLAN.md` until a scenario harness is accepted. |
| Release readiness | DEFERRED | Phase 11 is not a public release-readiness phase; threat model, release checklist, security scan, dependency audit, license check, and secret scan remain release-readiness controls. |

## Findings

### Blocking

- None.

### Major

- None.

### Improvement

- Code-size tripwire remains worth tracking outside this final QA task:
  `cmd/fid/issue.go` is 984 lines and `internal/journal/replay.go` is 485 lines
  by `wc -l`. This did not block Phase 11 because the QA plan did not define a
  file-size gate, behavior gates passed, and no product-code edits are allowed
  in this QA Lead task.

## Unresolved Risks And Carry-Forward Candidates

- VAL-014 remains pending human execution in
  `project-docs/docs/phases/phase-11/qa/human-checklist.md`.
- Existing carry-forward items CF-008, broad CF-011, and CF-013 remain governed
  by `project-docs/docs/current/decisions_active.md`; this review did not
  perform wrap-up disposition.
