# Phase 20 - Task 10: Acceptance-criterion and label contracts - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Migrated `issue.ac.add/list`, `label.add/list`, and `issue.label.add/remove`
  to their exact `fid-cli.v1` Response Envelope contracts while retaining
  concise stdout-only Human success.
- Added typed missing/duplicate/limit failures, aggregate revision and Event
  identity propagation, deterministic ordering, selected-project routing, and
  the Task 10 revision golden families.
- Enforced the accepted UTF-8, ASCII-control, byte-length, ordinal, optional
  value, and 1,000/100/10,000 collection bounds in the domain layer.
- Preserved `issue ac verify` as a guarded legacy leaf and removed `--format`
  only from the two migrated list leaves.
- Fixed a restored concurrency defect by registering root-flag capabilities at
  package initialization rather than while parallel tests construct commands.
- Preserved the appended Event's RFC3339 nanosecond timestamp in every migrated
  Issue mutation response and strengthened create/update/transition and Task 10
  identity assertions against the journal tail.
- Routed arity, unknown-flag, and missing/blank required-actor failures for all
  six Task 10 leaves through exact stderr-only `invalid-argument` envelopes,
  while retaining concise Human error prose.

## Files created or modified

- `cmd/fid/issue_ac.go`, `issue_label.go`, `label.go`,
  `issue_mutation_response.go`, `label_response.go`, and focused/e2e tests.
- `internal/issues/ac.go`, `ac_list.go`, `label.go`,
  `inventory_errors.go`, and focused tests.
- `internal/labels/add.go`, `list.go`, `types.go`, and focused tests.
- `internal/validation/text.go` and `text_test.go`.
- `internal/journal/replay_test.go` and
  `testdata/journal/p20-issue-revision-effects.golden.jsonl`.
- `internal/dev/block_test.go` and `finish_test.go` for the now-observable AC
  aggregate revision advance.
- Affected `cmd/fid` e2e helpers and assertions for migrated stdout/envelope
  behavior.
- `cmd/fid/integration/e2e_v01_closure_helpers_test.go` for the accepted
  stdout-only AC success channel.
- `cmd/fid/canonical_request.go` and `task10_contract_test.go` for canonical
  request-error rendering and the exhaustive six-operation negative matrix.

## Documentation updated

- `project-docs/docs/phases/phase-20/decisions.md` - DEC-141 recovery decision.
- `project-docs/docs/phases/phase-20/run-ledger.md`, `triage-log.md`, and
  `tl-escalations.md` - preserved and reclassified ESC-20-010 recovery.
- `project-docs/docs/current/phase_brief.md` and `progress_current.md` - Phase
  20 execution state.
- `project-docs/docs/reviews/2026-08-15-workflow-reliability-overhaul-plan.md`
  - durable recovery checkpoint.

## Validation Contract Coverage

- `VAL-008`: PASS - Human-mode success channels are asserted by
  `TestIssueACAddCommandAddsAcceptanceCriterionAndPrintsConfirmation`,
  `TestIssueACListCommandMatchesClosedHumanAndJSONContracts`,
  `TestLabelAddCreatesOneRowInOneEventWithoutActivityLog`,
  `TestLabelListMatchesClosedHumanAndJSONContracts`, and the Issue-label
  add/remove projection tests.
- `VAL-010`: PASS - `TestAddAcceptanceCriterionPersistsAndReadsBackInOrdinalOrder`
  proves successive revisions; `TestIssueLabelMutationsJSONMatchContract`
  proves add/remove revisions; duplicate/missing/limit tests prove no advance.
- `VAL-011`: PASS - `TestAddAcceptanceCriterionEventRoundTripStaysCanonical`,
  the Issue-label revision-Effect assertions, and
  `TestReplayJournalAppliesPhase20IssueRevisionGolden` prove exact
  value/prior/precondition tuples and replay for all three Task 10 families.
- `VAL-017`: PASS - `TestAcceptanceCriterionAddJSONMatchesContract`,
  `TestLabelAddJSONMatchesContract`, and
  `TestIssueLabelMutationsJSONMatchContract` compare mutation Event identity
  with the journal tail; AC/label list tests use `decodeReadEnvelope`, which
  requires `event:null`.
- `VAL-021`: PASS - AC add exact shape/revision/Event is covered by
  `TestAcceptanceCriterionAddJSONMatchesContract`; add missing/duplicate
  atomic codes by `TestAcceptanceCriterionAddJSONUsesSpecificAtomicFailureCodes`;
  text/ordinal boundaries and the 1,000 cap by the two focused domain boundary
  tests; list exact shape/order/empty/missing/limit by the three AC list tests.
- `VAL-022`: PASS - Label add/list exact `colour` shapes and Event/null behavior
  are covered by the label JSON tests; duplicate/empty/boundary/10,000-cap cells
  by the label negative/domain tests; Issue-label add/remove exact
  shape/revision/Event and every missing/duplicate/limit cell by
  `TestIssueLabelMutationsJSONMatchContract`,
  `TestIssueLabelJSONUsesSpecificAtomicFailureCodes`, and
  `TestIssueLabelAddJSONReportsNestedLimitWithoutAppend`.
- `VAL-041`: PASS - AC and Issue-label Event tests assert the supplied actor and
  exact three-Effect domain/revision/activity set; Label add asserts its actor
  and sole catalog Effect. No trusted Evidence, Human Outcome, or protected
  authority Effect is present.
- `VAL-044`: PASS - `TestTask10OperationsUseSelectedProjectOutsideWorkingDirectory`
  covers all six selected-project leaves;
  `TestIssueACVerifyRemainsGuardedLegacyLeaf` covers the legacy exclusion;
  `TestTask10JSONRequestErrorsUseCanonicalEnvelope` covers arity and unknown
  flags for all six leaves plus every missing/blank required Task 10 flag;
  removed `--format`, explicit-empty option, UTF-8/control, and every
  just-inside/outside byte/ordinal cell are covered by the named request and
  domain tables.
- `VAL-045`: PASS - AC order is asserted by ordinal then ID, labels by name;
  exact-limit success and over-limit typed failure are asserted for 1,000 ACs,
  100 Issue labels, and 10,000 catalog labels without silent truncation.

DEC-140 operation-cell mapping: `issue.ac.add` is covered by its exact success,
missing/duplicate, boundary, cap, Event, and revision tests; `issue.ac.list` by
exact/read-null, order, empty, missing, removed-format, and cap tests;
`label.add` by exact success/nullability, duplicate, optional-empty, boundary,
Event, and attribution tests; `label.list` by exact/read-null, name order,
empty, removed-format, and cap tests; `issue.label.add/remove` by the shared
exact mutation test plus the per-operation missing/duplicate/token/revision/
Event/atomicity tables. No named operation is represented by sampling another.

## Tests run

- `go test ./cmd/fid -run '^TestTask10JSONRequestErrorsUseCanonicalEnvelope$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestTask10(JSONRequestErrorsUseCanonicalEnvelope|RequiredActorHumanErrorsRemainProse)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.060s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  6.577s` (terminated `exit_code: 0`)
- `go test ./internal/issues ./internal/labels -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/labels  0.598s` (terminated `exit_code: 0`)
- `go test ./internal/validation ./internal/issues ./internal/labels ./internal/journal ./cmd/fid -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid  6.564s` (terminated `exit_code: 0`)
- `go test ./... -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow  0.690s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	10.031s	coverage: 80.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					77.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.193s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- No contract deviation. Recovery preserved the first focused `FAIL` caused by
  decoding the new AC list envelope as a legacy array, then the byte-identical
  rerun's nine displaced stdout/`--format` assertions. The tests were migrated
  to the accepted contract and the byte-identical focused command now passes;
  no assertion was deleted, weakened, or re-pointed.
- The first exact format gate named `cmd/fid/root.go`; `gofmt` repaired that
  restored file and the byte-identical rerun passed silently with observed exit
  0. The first exact lint gate named only the obsolete label Markdown decoder;
  removing that unreachable helper produced `0 issues.` on the byte-identical
  rerun. No suppression was added.
- Go VCS stamping resolved the linked worktree's `.git` indirection to `/tmp`,
  so the exact integration and build commands first failed with `error obtaining
  VCS status: exit status 128` and `Use -buildvcs=false to disable VCS
  stamping.` (both terminated `exit_code: 1`). No command flag was changed.
  Every final gate above ran byte-identically against committed HEAD
  `62a3e46ebb69b76b93b5a559dc82d9bc9c0ccd3e` in a normal local clone with VCS
  stamping enabled.
- The first normal-clone race/coverage run failed with `FAIL` (terminated
  `exit_code: 1`) because `TestIssueTransitionCanonicalOutcomeTable` compared
  the whole-second Response Envelope timestamp with the journal's nanosecond
  timestamp. The TRD requires sub-second precision where available and VAL-017
  requires matching Event identity. Commit `32ce0f1` changed returned mutation
  identity timestamps to `RFC3339Nano` and strengthened create, update,
  transition, AC, catalog-label, and Issue-label comparisons. Ten repeated
  focused race runs passed before the byte-identical full gate passed on the
  repaired commit.
- The next exact integration run failed with
  `FAIL github.com/mbell-nzl/fid/cmd/fid/integration` (terminated `exit_code: 1`)
  because its closure fixture discarded AC-add stdout and parsed the legacy
  stderr channel. Commit `802e32b` changed only that fixture to require empty
  stderr and parse the accepted stdout success. The byte-identical integration
  command then passed. No product assertion, threshold, or gate was narrowed.
- Final independent integrated review at `4ae80e4` returned `REPAIR` after real
  binary probes showed empty stderr for JSON arity/unknown-flag failures and
  Human prose for AC required-flag failures. The new exhaustive request test
  initially failed all 18 original cells with `FAIL` (terminated `exit_code: 1`).
  Commit `62a3e46` added the domain-neutral canonical request renderer,
  wired all six leaves through it, enforced missing/blank required actors as
  input errors, and expanded the final matrix to 22 JSON cells plus three Human
  channel regressions. The focused rerun and the byte-identical full gate pass;
  no error schema, request rule, assertion, or command was weakened.

## Concerns, technical debt, or notes for the Tech Lead

- The broader test-performance refactors remain outside this recovery per
  DEC-141. The canonical integration gate already uses only
  `./cmd/fid/integration/...`, avoiding the former tagged `./...` duplication.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `62a3e46ebb69b76b93b5a559dc82d9bc9c0ccd3e`

## Commit subject

- `fix(cli): envelope task 10 request failures`
