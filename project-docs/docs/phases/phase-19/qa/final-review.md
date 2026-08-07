# Phase 19 Final Review - Workflow Authoring CLI + Escalation Default Workflow

Reviewer: Codex QA Lead (GPT-5)
Date: 2026-08-02
Operating mode: Heavy

> **Superseded 2026-08-07.** This verdict (BLOCK) is historical and is retained
> unedited for the audit trail. All five blocking findings were dispositioned by
> DEC-105 and re-reviewed in `final-review-2.md`, which carries the operative
> verdict.

## Verdict

**BLOCK**

All eight exact phase-end commands terminated successfully, all 69
deterministic validation assertions have fresh reproducible PASS evidence, all
six completion reports pass `report-lint`, and the Phase 19 performance
controls pass without relaxation. The phase is not ready because:

1. Human-only VAL-026 is BLOCKED pending PASS on all nine rows in
   `human-checklist.md`.
2. `02_TRD.md` §3.1 and DEC-096 prohibit `internal/workflow` from using journal
   append/projection surfaces, while the shipped implementation and DEC-089
   deliberately do so. The canonical architecture record is contradictory.
3. `03_TRD_SCHEMA.md` §14.11 says non-v0.1 tables are excluded from maintenance
   structural auditing, while Phase 19 ships an audit over
   `workflow_assignments`; current context and the roadmap closure record are
   also stale. Full canonical synchronization is not complete.
4. Task 4 implementation commit `32614bbbd50e683b61921e9ed6beabee73b46021`
   has the reported subject but no required `Refs: phase-19/task-4` body/footer.
5. The in-scope `internal/issues/create.go` remains 412 lines and `Create`
   spans lines 19-246, exceeding the 400-line file and 60-line function
   budgets.

Heavy mode permits no carry-forward. These findings and the Human dependency
must be resolved and re-reviewed; no assertion or finding is silently waived.

## Scope And Evidence

- Phase brief: `project-docs/docs/current/phase_brief.md`
- Validation contract:
  `project-docs/docs/phases/phase-19/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-19/qa/qa-plan.md`
- Task specs and completion reports:
  `project-docs/docs/phases/phase-19/tasks/`
- Mid-phase reviews: `project-docs/docs/phases/phase-19/qa/review-1.md` and
  `review-2.md`
- Phase decisions: `project-docs/docs/phases/phase-19/decisions.md`, especially
  DEC-079, DEC-101, DEC-102, and DEC-103
- Current decisions: `project-docs/docs/current/decisions_active.md`
- Quality and workflow rules: `project-docs/specs/08_QUALITY_PLAN.md` and
  `09_AI_WORKFLOW.md` §§2.6, 6.1, and 8
- Canonical requirement/specification surfaces reviewed:
  `01_PRD.md`, `02_TRD.md` §§3/4.6, `03_TRD_SCHEMA.md` §§3/14.11,
  `05_DESIGN_PLAN.md` §3.1, `06_DESIGN_PLAN_CLI.md`, `GLOSSARY.md`, and
  `10_IMPLEMENTATION_ROADMAP.md`

No held-out scenario repository, evaluator prompt, or excluded scenario
material was read.

### Fresh Evidence Register

Every focused test command below used `-count=1`; every command terminated
with exit 0. Assertion rows cite the specific test or subtest within these
runs rather than inheriting another assertion's outcome.

| Evidence | Fresh command or state check | Terminal evidence |
|---|---|---|
| F-01 | Task 1 focused migration/registry suite in `./internal/db ./internal/workflow` | `--- PASS: TestRegistryEqualsV01Definitions`; package terminal `ok` |
| F-02 | Full sorted tuple query against immutable Phase 18 commit `097e45cc5111be5ddcd5bf232645659045571521` and current candidate | both exits 0; `byte_identical=true`; 54 full rows; 6,248 characters on each side |
| F-03 | Task 2 focused authoring/assignment-core suite in `./internal/workflow` | all named tests and matrices PASS; package terminal `ok` |
| F-04 | Task 3 focused define/policy suite in `./internal/workflow ./internal/policy ./cmd/fid` | all named tests and subtests PASS; package terminals `ok` |
| F-05 | Task 4 focused assignment/creation suite in `./internal/workflow ./internal/issues ./cmd/fid` | all named tests and subtests PASS; package terminals `ok` |
| F-06 | `go test -tags integration -race -count=1 -v ./cmd/fid -run '^TestV01ParityGoldenOutputs$'` | parent and all six get/list subtests PASS; `ok github.com/mbell-nzl/fid/cmd/fid` |
| F-07 | Task 5 focused list/format suite in `./internal/workflow ./internal/format ./cmd/fid` | all named tests PASS; package terminals `ok` |
| F-08 | Task 6 focused integration/maintenance suite in `./internal/maintenance ./cmd/fid` | all four named tests PASS; `ok github.com/mbell-nzl/fid/cmd/fid 2.840s` |
| F-09 | Built `/tmp/fid` scratch workspace define/assign/create/list exercise | project fallback and kind precedence SQL rows correct; json/toon/markdown/yaml all exit 0 |
| F-10 | Reserved-rule SQL and production inventories | SQL `[]`; both exact `rg` inventories exit 1 with no output (the required no-match condition) |

### Immutable Tuple And Golden Evidence

VAL-003 used all columns of all four workflow tables. The fresh sorted
escalation tuples were exactly:

```text
workflows: ('escalation-default','Escalation default','escalation',1,'<NULL>')
workflow_states: ('escalation-default','open','started',1)
workflow_states: ('escalation-default','resolved','completed',2)
workflow_transitions: ('escalation-default','open','resolved')
workflow_transition_requires: ('escalation-default','open','resolved','linked_to_decision','<NULL>')
```

For VAL-030, the Phase 18 binary was built from verified immutable commit
`097e45cc5111be5ddcd5bf232645659045571521`. One ordered UNION selected every
workflow field, every state/category/ordinal, every transition, and every
requirement/rule/`params_json` with explicit `<NULL>` normalization for
`epic-default`, `task-default`, and `subtask-default`. The immutable and current
outputs both terminated, contained 54 full tuples and 6,248 characters, and
were byte-identical. This is full structural conformance, not count parity.

VAL-059 provenance and permanent evidence were reproduced without either
struck command form:

- Task 3 implementation commit resolves exactly as
  `50533e2b813bf073cd75f75eb70dd09cd15e9040`.
- At Task 4 implementation commit `32614bbbd50e683b61921e9ed6beabee73b46021`,
  the fixed fixture, `capture.sh`, all six captures, `SHA256SUMS`, and
  `normalization.json` are committed. The parity directory is unchanged
  through the reviewed candidate (`git diff --exit-code` produced no output).
- Normalization is the explicit no-op
  `{"json_pointers":[],"replacement":null}`.
- The committed capture identities are:

```text
e04e116c6df34cb85825dc4aa3b886a4edc56143409808f5ac7e4eb56fdbb5f4  issue-get-epic.json
0e71ea0d6ddbee4e9e62275d7f96aed873f23629945cbbe5eb55f109b05ed7f4  issue-get-task.json
ed7e271430d01257a3f92e8f59326e7c1ff562ce5ff16654def825c681635541  issue-get-subtask.json
8dd2f3cf2be2d9e9b78bce77a6121ada76b10b593fd75c91768881d21eecd9f4  issue-list-epic.json
febad62c2f5f9720a680689f464fd80ac6389426f6b63ba00571948ee9b3d9b5  issue-list-task.json
152bb980ee36e22ad1bbd008e0257a565810a568212c6d24afd3b48ce41570a6  issue-list-subtask.json
```

- The permanent integration test copies the committed fixture to
  `t.TempDir()`, runs the six commands against the current binary, and compares
  each result to its committed capture with raw string/byte equality. F-06
  freshly passed all six comparisons.
- Per DEC-103/CF-025, the manual VAL-059 shell ceremony is retired as a
  required procedure; its committed artifacts remain provenance. Neither the
  broken root-level manifest check nor the self-referential baseline-directory
  check was run or cited.

## Phase-Level Gate Results

The eight exact Phase End commands from QA plan §5 were run fresh in the
specified order. Each terminated before PASS was recorded.

| # | Gate | Command | Result | Exact terminal evidence |
|---:|---|---|---|---|
| 1 | Format | `gofmt -l cmd internal` | PASS | `exit 0 (no output)` |
| 2 | Lint/static analysis | `golangci-lint run ./...` | PASS | `0 issues.` |
| 3 | Unit/race | `go test -race ./...` | PASS | `ok github.com/mbell-nzl/fid/internal/workflow (cached)` |
| 4 | Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS | `ok github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 80.5% of statements` |
| 5 | Coverage report | `go tool cover -func=coverage.out` | PASS | `total: (statements) 77.0%` |
| 6 | Integration | `go test -tags integration -race ./...` | PASS | `ok github.com/mbell-nzl/fid/internal/workflow (cached)` |
| 7 | Benchmarks | `go test -bench . -benchmem ./...` | PASS | `ok github.com/mbell-nzl/fid/internal/workflow 0.378s` |
| 8 | Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | `exit 0 (no output)` |
| — | Standalone E2E | `N/A` per `AGENTS.md`; Phase 19 E2E is integration-tagged | N/A | Covered by gate 6 and F-08; no held-out harness exists |
| — | Full phase check | The eight commands above, in order | PASS | 8/8 terminated with exit 0 |

Coverage package evidence was `cmd/fid` 78.1%, `internal/issues` 81.4%,
`internal/claims` 81.1%, `internal/workflow` 80.5%, `internal/policy` 86.3%,
and `internal/journal` 81.0%; all package controls pass.

The fresh benchmark controls were:

```text
BenchmarkList10kIssuesWithRegistry-6       42847790 ns/op   42.85 ms/registry-read
BenchmarkRebuildWorkspace100kEvents-6   23934996215 ns/op   23935 ms/rebuild
```

Both controls pass without weakening (`42.85 ms < 500 ms`,
`23.935 s < 30 s`). They improve on the Phase 18 evidence
(`49.67 ms`, `27.661 s`) and planning baseline (`44.03 ms`, `24.674 s`).

## Completion Report And Commit Audit

Each report/spec pair was passed to the exact
`project-docs/tools/report-lint.sh <report> <spec>` command in this QA Lead
session. All report files were committed and clean; the pre-existing modified
`phase-19/run-ledger.md` was outside this audit and was not edited.

| Task | Implementation commit and subject | Report commit and subject | Exact-command/report-lint result |
|---|---|---|---|
| 1 | `a509d59efcc826df417f64ea72a43d743936e669` — `feat(workflow): add assignment schema and escalation seed` | `3b061e1213064554d762a0f7aa1999f4edb555d2` — `docs(task): correct VAL-031 evidence token in phase-19 task-1 report` | PASS — `PASS report-lint completed with no failures` |
| 2 | `10d4d242403851b5c676a4308cde598d251516ea` — `feat(workflow): add authoring validation and effects` | `f61f527934de0f284f4981995e19deb8b3ba4ce9` — `docs(task): finalize completion report for phase-19 task-2` | PASS — `PASS report-lint completed with no failures` |
| 3 | `50533e2b813bf073cd75f75eb70dd09cd15e9040` — `feat(workflow): add define command and journal event` | `697c673620b336d81e2d5370c2dd373b46f60202` — `docs(task): finalize completion report for phase-19 task-3` | PASS — `PASS report-lint completed with no failures` |
| 4 | `32614bbbd50e683b61921e9ed6beabee73b46021` — `feat(workflow): add assign command and assignment-aware creation` | `ba65078cdacffb6e153cfdeb4e20196084cab947` — `docs(task): record phase-19 task-4 closure evidence` | PASS — `PASS report-lint completed with no failures`; metadata defect below |
| 5 | `8f79013c7e8afcd3d4bd73f9821484d9dd9d9334` — `feat(workflow): add workflow list command` | `530a82ab8e8ed4e5d75b957b659e1e8cebbd638d` — `docs(task): record completion report lint for phase-19 task-5` | PASS — `PASS report-lint completed with no failures` |
| 6 | `09aa0a081d2ecfee025c3b4d8158bc6203d02d98` — `feat(workflow): close workflow authoring integrity` | `988737b7f68da1ef481a98b0184400248321837b` — `docs(task): finalize completion report for phase-19 task-6` | PASS — `PASS report-lint completed with no failures` |

All twelve hashes resolve and all subjects match the reports. The exact task
commands recorded by each report match its task spec, including Task 5 and
Task 6. Task 4's DEC-103 closure report is committed and lint-clean, resolving
Review 2's report-evidence blocker. Separately, `git show -s --format='%B'` on
the Task 4 implementation commit contains only its subject and omits the
required `Refs: phase-19/task-4`; therefore QA-31 is FAIL despite the green
report-lint result.

## Validation Contract Outcomes

No result is inherited and no heavy-mode carry-forward is used. Each PASS
cites a fresh test, tuple/CLI trace, or exact file-state inventory.

| VAL ID | Outcome | Validator | Distinct evidence |
|---|---|---|---|
| VAL-001 | PASS | QA Lead | F-01 `TestMigration0009WorkflowAssignmentsSchemaAndConstraints`: exact columns, FK, CHECKs, and partial indexes PASS. |
| VAL-002 | PASS | QA Lead | F-01 `TestMigrations0009And0010RoundTripWithoutResidue`: up/down/up schema and row snapshots PASS. |
| VAL-003 | PASS | QA Lead | F-01 `TestMigration0010SeedsEscalationDefaultExactlyAndPreservesV01Tuples` plus the five full sorted tuples printed above. |
| VAL-004 | PASS | QA Lead | F-01 `TestLoadSurfacesEscalationRuleAsUnimplementedWithoutWriting`: known/unimplemented rule surfaced without panic. |
| VAL-005 | PASS | QA Lead | F-01 `TestDecideProjectionFreshnessRequiresPhase19SchemaVersion`: stale and current branches PASS. |
| VAL-006 | PASS | QA Lead | F-03 `TestBuildDefinitionEffectsIsOrderedAndDeterministic`: golden order and repeat equality PASS. |
| VAL-007 | PASS | QA Lead | F-03 `TestValidateDefinitionRequestIdentifierMatrix`: four distinct value-naming rejections and accepted control PASS. |
| VAL-008 | PASS | QA Lead | F-03 `TestBuildDefinitionEffectsCancelledStatePolicy`: complete cancel edges and explicit-edge deduplication PASS. |
| VAL-009 | PASS | QA Lead | F-04 `TestDefineAppendsOneAtomicWorkflowDefinedEvent`: one current-schema event containing all effects PASS. |
| VAL-010 | PASS | QA Lead | F-04 `TestWorkflowDefineRejectsWithoutMutatingJournalOrWorkflowTables`: rejection journal hashes unchanged; accepted append control PASS. |
| VAL-011 | PASS | QA Lead | F-04 `TestMatrixWorkflowAuthoringRowsPermitOnlyTLAndHuman`: exact role sets for both actions PASS. |
| VAL-012 | PASS | QA Lead | F-04 `TestWorkflowDefineNormalizesArrowsAndAcceptsAppliesToVocabulary`: `->`, `-->`, `→`, and `⟶` tuple parity PASS. |
| VAL-013 | PASS | QA Lead | F-04 `TestWorkflowDefineProjectsDeclaredStatesTransitionsAndRules`: literal confirmation and SQL-backed counts PASS. |
| VAL-014 | PASS | QA Lead | F-05 `TestWorkflowAssignProjectsOneKindAssignment`: one kind event and one exact projected row PASS. |
| VAL-015 | PASS | QA Lead | F-05 project-scope subtest in `TestAssignEmitsScopeAwareReplacementEffects`: insert-only event and NULL-kind row PASS. |
| VAL-016 | PASS | QA Lead | F-05 `TestWorkflowAssignRejectsWithoutMutation`; permanent CLI test asserts six exact pairwise-distinct errors and valid control. |
| VAL-017 | PASS | QA Lead | F-05 `TestIssueCreateUsesAssignedWorkflowAndItsInitialState`: unstarted and first-declared initial-state branches PASS. |
| VAL-018 | PASS | QA Lead | F-07 `TestWorkflowListRendersDocumentedJSONFieldSet`: exact keys, SQL counts, and `any`→null PASS. |
| VAL-019 | PASS | QA Lead | F-07 `TestListFiltersCountsAssignmentsAndOrdering`: exact qualifying/excluded kind sets PASS. |
| VAL-020 | PASS | QA Lead | F-07 `TestWorkflowListFormatsFiltersAndAssignments`: exact kind assignment object PASS. |
| VAL-021 | PASS | QA Lead | F-08 `TestV02aWorkflowAuthoringExitScenarioAndRebuildParity`: every-table/schema/FK snapshots and journal hashes equal over two rebuilds. |
| VAL-022 | PASS | QA Lead | F-08 `TestCheckWorkflowAssignmentsFlagsIncompatibleAppliesTo`, compatible counterpart, and CLI maintenance test: both states and warn-as-error PASS. |
| VAL-023 | PASS | QA Lead | F-08 exit scenario reserved-rule cases: both version-scoped CLI rejections and unchanged journal bytes PASS. |
| VAL-024 | PASS | QA Lead | F-08 exit scenario: custom-named target rejected/no append and v0.1-named target under custom workflow succeeds. |
| VAL-025 | PASS | QA Lead | F-08 `TestV02aWorkflowAuthoringExitScenarioAndRebuildParity`: full define/assign/list/create/concern rule/control journey PASS. |
| VAL-026 | BLOCKED | Human | `human-checklist.md` has exactly nine pending rows; 0/9 Human outcomes recorded. QA judgment is not substituted. |
| VAL-027 | PASS | QA Lead | F-01 assignment-schema uniqueness subtest: duplicate kind rejected and two distinct kinds accepted. |
| VAL-028 | PASS | QA Lead | F-01 assignment-schema uniqueness subtest: first project row accepted, second rejected, one remains. |
| VAL-029 | PASS | QA Lead | F-01 constraint subtests `kind_requires_kind`, `project_rejects_kind`, and both accepted shapes PASS. |
| VAL-030 | PASS | QA Lead | F-02 immutable Phase 18/current 54-full-tuple query: both exit 0, 6,248 characters, byte-identical. |
| VAL-031 | PASS | QA Lead | F-01 no-write loader test plus exact non-test `linked_to_decision` inventory: exit 1/no output. |
| VAL-032 | PASS | QA Lead | F-03 states/categories matrix: empty, count mismatch, unknown category, and accepted control PASS. |
| VAL-033 | PASS | QA Lead | F-03 transition matrix: unknown-edge and duplicate-edge value-naming rejections PASS. |
| VAL-034 | PASS | QA Lead | F-03 `TestDefinitionRuleVocabularyAndParams`: all four implemented rule/edge/params effects PASS. |
| VAL-035 | PASS | QA Lead | F-03 rule vocabulary matrix: three version-scoped reserved errors and distinct unknown-rule error PASS. |
| VAL-036 | PASS | QA Lead | F-03 params matrix: invalid JSON and empty roles reject; valid params retained verbatim. |
| VAL-037 | PASS | QA Lead | F-03 `TestBuildDefinitionEffectsAcceptsRequirementOnMaterializedEdge`: unknown edge rejects, declared/materialized edge accepts. |
| VAL-038 | PASS | QA Lead | F-03 `TestBuildDefinitionEffectsZeroCancelledApplicability`: non-issue accepted/no edge and issue-applicable rejected. |
| VAL-039 | PASS | QA Lead | F-03 `TestBuildDefinitionEffectsRejectsMultipleCancelledStates`: both states named and no effects produced. |
| VAL-040 | PASS | QA Lead | F-03 `TestInitialState`: first-unstarted and first-declared branches PASS. |
| VAL-041 | PASS | QA Lead | F-03 `TestValidateAssignmentRequestRejectsAndAcceptsScopeMatrix`: five distinct rejects and two accepted controls PASS. |
| VAL-042 | PASS | QA Lead | F-03 no-write tests plus exact mutating-SQL inventory over core files: exit 1/no output. |
| VAL-043 | PASS | QA Lead | F-04 define projection test: declaration and all four projected table tuple sets equal. |
| VAL-044 | PASS | QA Lead | F-04 define projection/load path returns complete states, transitions, and rules. |
| VAL-045 | PASS | QA Lead | F-04 successful-define activity snapshot subtest: `activity_log` unchanged. |
| VAL-046 | PASS | QA Lead | F-04 rejection test: all four workflow-table snapshots unchanged for all rejected cases. |
| VAL-047 | PASS | QA Lead | F-04 rejection matrix: 15 named cases, non-zero exits, value-naming pairwise-distinct stderr. |
| VAL-048 | PASS | QA Lead | F-04 `TestPolicyShowIncludesWorkflowAuthoringRowsInEveryFormat`: both rows in json/toon/markdown/yaml PASS. |
| VAL-049 | PASS | QA Lead | F-04 `TestWorkflowDefineAuthorizesBeforeParsingAndDoesNotAppendOnDenial`: denied no-append and `tl` success PASS. |
| VAL-050 | PASS | QA Lead | F-04 applies-to matrix: epic/task/subtask/qa_review/escalation/any projections and invalid no-append PASS. |
| VAL-051 | PASS | QA Lead | F-05 kind reassignment subtest: assignment replaced while superseded workflow remains loadable/listed. |
| VAL-052 | PASS | QA Lead | F-05 kind effect subtests plus permanent replay test: insert, ordered delete+insert, and rebuilt row PASS. |
| VAL-053 | PASS | QA Lead | F-05 project effect subtests plus permanent replay test: insert, ordered delete+insert, and rebuilt row PASS. |
| VAL-054 | PASS | QA Lead | F-05 `TestIssueCreateResolvesKindAssignmentBeforeProjectAssignment` project-fallback subtest; F-09 SQL corroborates `qa-any`. |
| VAL-055 | PASS | QA Lead | F-05 precedence subtest; F-09 SQL corroborates task-kind `qa-task` winning over project `qa-any`. |
| VAL-056 | PASS | QA Lead | F-05 `TestWorkflowAssignRejectsWithoutMutation`: each rejected case preserves journal and assignment rows. |
| VAL-057 | PASS | QA Lead | F-05 successful-assign activity snapshot: `activity_log` unchanged. |
| VAL-058 | PASS | QA Lead | F-05 `TestCreateSupportsEpicTaskAndSubtaskAuthoring`: zero-assignment epic/task/subtask default+`todo` tuples PASS. |
| VAL-059 | PASS | QA Lead | F-06 permanent raw-byte golden test passed six subtests; verified Task 3 commit, unchanged committed fixture/captures/manifest, and no-op normalization above. |
| VAL-060 | PASS | QA Lead | F-07 four-format test plus F-09 built CLI outputs: json/toon/markdown/yaml structures equivalent. |
| VAL-061 | PASS | QA Lead | F-07 `TestWriteWorkflowListMarkdownMatchesCommittedLiteral`: byte-exact committed in-test literal PASS. |
| VAL-062 | PASS | QA Lead | F-07 active-filter subtests: inactive excluded with filter and included without it. |
| VAL-063 | PASS | QA Lead | F-07 combined-filter subtest: actual IDs equal exact kind/active intersection. |
| VAL-064 | PASS | QA Lead | F-07 project-assignment shape subtest and F-09 output: `kind` key absent, not null. |
| VAL-065 | PASS | QA Lead | F-07 unassigned-shape subtest: present `assigned_to: []`. |
| VAL-066 | PASS | QA Lead | F-07 ordering/repeat subtests: workflow IDs and assignments ordered; consecutive bytes equal. |
| VAL-067 | PASS | QA Lead | F-07 no-match and empty-markdown tests: all four formats exit 0 with empty arrays. |
| VAL-068 | PASS | QA Lead | F-10 exact ordered SQL over both reserved names returned `[]`, exit 0. |
| VAL-069 | PASS | QA Lead | F-10 exact required production `rg` inventory returned exit 1/no output; no evaluator/query match. |
| VAL-070 | PASS | QA Lead | F-10 exact required non-test `FrozenDefinitions(` inventory returned exit 1/no output. |

Assertion summary: **69 PASS, 0 FAIL, 1 BLOCKED, 0 CARRY-FORWARD**.

## Requirement Compliance

| Requirement | Status | Evidence |
|---|---|---|
| PRD-WF-005a, WF-011..015, OPS-001..007, OUT-001..004 | MET | VAL-001..025 and VAL-027..070 deterministic outcomes; VAL-026 remains Human-only. |
| `02_TRD.md` §4.6 workflow data and precedence | MET | VAL-003/030/043/044/054/055; full tuple and built-CLI evidence. |
| `02_TRD.md` §3/§3.1 architecture boundaries | NOT MET | §3 assigns write paths to `internal/workflow`, but §3.1/DEC-096 permits only `journal.Effect`; `define.go`/`assign.go` call lock, prepare, dry-run project, append, and project APIs. |
| `03_TRD_SCHEMA.md` §3 schema | MET | VAL-001..003 and VAL-027..030. |
| `03_TRD_SCHEMA.md` §14.11 maintenance parity | NOT MET | Text excludes tables absent in v0.1; `internal/maintenance/workflow_assignments.go` now audits `workflow_assignments`. |
| `05_DESIGN_PLAN.md` §3.1 | MET | Workflow define/assign/list family and limitation references are present. |
| `06_DESIGN_PLAN_CLI.md` | MET | Help/behavior and four-format output evidence F-04/F-05/F-07/F-09. |
| `01_PRD.md` and `GLOSSARY.md` sync | MET | Workflow authoring, assignment, projected registry, custom workflow, assignment, and initial-state terms present. |
| Roadmap Phase 19 closure obligation | PARTIAL | Phase 19 scope/exit text exists, but no final closure disposition is recorded yet. |
| Journal/schema compatibility and deterministic rebuild | MET | VAL-002/009/014/015/021/052/053/059; two rebuilds and raw-byte golden test pass. |
| Security, input handling, data loss, and atomic rejection | MET | QA-28; parameterized SQL, no shell execution, rejection/no-append tests, FK/schema and journal hashes. |
| Dependency budget | MET | `git diff --stat 196a31d..HEAD -- go.mod go.sum` produced no output. |
| `08_QUALITY_PLAN.md` structure budgets | NOT MET | `internal/issues/create.go` 412 lines; `Create` approximately 228 lines. |
| Phase commit metadata convention | NOT MET | Task 4 implementation commit lacks `Refs: phase-19/task-4`. |

## Deterministic Checklist Execution

| QA ID | Result | Reproducible evidence |
|---|---|---|
| QA-01 | PASS | F-01 exact schema/constraint introspection test PASS. |
| QA-02 | PASS | F-01 up/down/up residue and fingerprint test PASS. |
| QA-03 | PASS | F-01 plus the five exact escalation tuples above. |
| QA-04 | PASS | F-01 known-unimplemented loader test PASS without panic. |
| QA-05 | PASS | F-01 stale/current projection-version branches PASS. |
| QA-06 | PASS | F-03 ordered golden and repeated-effect equality PASS. |
| QA-07 | PASS | F-03 identifier matrix rejects four distinct values and accepts control. |
| QA-08 | PASS | F-03 exact cancel transition set and deduplication PASS. |
| QA-09 | PASS | F-04 one atomic current-schema event test PASS. |
| QA-10 | PASS | F-04 per-rejection no-append hashes and accepted control PASS. |
| QA-11 | PASS | F-04 policy matrix names exactly `tl` and `human`. |
| QA-12 | PASS | F-04 four documented arrow aliases normalize identically. |
| QA-13 | PASS | F-04 confirmation/count assertions PASS. |
| QA-14 | PASS | F-05 first kind event/row test PASS. |
| QA-15 | PASS | F-05 first project event/NULL-kind row test PASS. |
| QA-16 | PASS | F-05 six exact pairwise-distinct rejection messages plus accepted control PASS. |
| QA-17 | PASS | F-05 both assignment-aware initial-state branches PASS. |
| QA-18 | PASS | F-07 exact JSON key/count/null assertions PASS. |
| QA-19 | PASS | F-07 exact kind-filter included/excluded ID sets PASS. |
| QA-20 | PASS | F-07 exact kind `assigned_to` object PASS. |
| QA-21 | PASS | F-08 all-table/schema/FK/journal snapshots equal over two rebuilds. |
| QA-22 | PASS | F-08 healthy/incompatible maintenance states and warn-as-error PASS. |
| QA-23 | PASS | F-08 both reserved CLI errors retain version scope and no append. |
| QA-24 | PASS | F-08 custom-state rejected/no append and v0.1-state success PASS. |
| QA-25 | PASS | F-08 complete slice-2 operator trajectory and unaffected control PASS. |
| QA-26 | PASS | Fresh coverage: cmd 78.1%; issues 81.4%; claims 81.1%; workflow 80.5%; policy 86.3%; journal 81.0%; total 77.0%. |
| QA-27 | PASS | Fresh benchmark lines: 42.85 ms read and 23.935 s rebuild, both under unchanged controls. |
| QA-28 | PASS | Exact no-shell inventory found no production `os/exec`, `exec.Command`, `sh -c`, or `bash -c`; SQL uses bound args; F-04/F-05/F-08 prove rejection atomicity and no data loss. |
| QA-29 | FAIL | `wc -l` gives 412 for `internal/issues/create.go`; `Create` starts line 19 and next function starts line 247. Dependency diff and lint are clean, but structure and canonical boundary checks are not. |
| QA-30 | FAIL | `02_TRD.md` §3.1 contradicts implementation/DEC-089; `03_TRD_SCHEMA.md` §14.11 omits shipped assignment audit; current brief/progress and roadmap closure are stale. |
| QA-31 | FAIL | Six reports are committed/clean and 6/6 report-lint PASS, all hashes/subjects resolve, but Task 4 implementation commit body lacks mandatory `Refs: phase-19/task-4`. |
| QA-32 | PASS | F-01 duplicate-kind reject and distinct-kind accepts PASS. |
| QA-33 | PASS | F-01 first/second project uniqueness branches PASS. |
| QA-34 | PASS | F-01 both malformed and both well-formed scope/kind shapes PASS. |
| QA-35 | PASS | F-02 full 54-tuple immutable/current output is byte-identical. |
| QA-36 | PASS | F-01 no-write loader snapshots and exact no-match evaluator inventory PASS. |
| QA-37 | PASS | F-03 states/categories rejection matrix and accepted control PASS. |
| QA-38 | PASS | F-03 unknown and duplicate edge errors PASS. |
| QA-39 | PASS | F-03 all four implemented rules in effects PASS. |
| QA-40 | PASS | F-03 reserved/unknown rule vocabulary and version tokens PASS. |
| QA-41 | PASS | F-03 invalid JSON/empty roles reject and verbatim params accept PASS. |
| QA-42 | PASS | F-03 unknown-edge reject and declared/materialized-edge accept PASS. |
| QA-43 | PASS | F-03 both zero-cancel applicability states PASS. |
| QA-44 | PASS | F-03 multiple-cancel rejection names both states and yields no effects. |
| QA-45 | PASS | F-03 both initial-state branches PASS. |
| QA-46 | PASS | F-03 five assignment rejects and two compatible controls PASS. |
| QA-47 | PASS | F-03 core no-write tests and mutating-SQL inventory exit 1/no output. |
| QA-48 | PASS | F-04 all four definition projection tuple sets PASS. |
| QA-49 | PASS | F-04 newly defined workflow complete load PASS. |
| QA-50 | PASS | F-04 define activity rows unchanged. |
| QA-51 | PASS | F-04 all workflow tables unchanged for every rejection. |
| QA-52 | PASS | F-04 15 rejection exits/messages distinct and value-naming. |
| QA-53 | PASS | F-04 both permission rows in four policy formats. |
| QA-54 | PASS | F-04 denied-role no-append and `tl` accepted control. |
| QA-55 | PASS | F-04 six applicability values, NULL `any`, invalid no-append. |
| QA-56 | PASS | F-05 reassignment row replacement while old workflow remains. |
| QA-57 | PASS | F-05 kind insert/delete+insert/rebuild row shape PASS. |
| QA-58 | PASS | F-05 project insert/delete+insert/rebuild row shape PASS. |
| QA-59 | PASS | F-05/F-09 uncovered task resolves to project workflow `qa-any`. |
| QA-60 | PASS | F-05/F-09 kind assignment `qa-task` wins over project assignment. |
| QA-61 | PASS | F-05 each assign rejection preserves journal and rows. |
| QA-62 | PASS | F-05 successful assign leaves activity rows unchanged. |
| QA-63 | PASS | F-05 no-assignment epic/task/subtask defaults and `todo` PASS. |
| QA-64 | PASS | F-06 six raw-byte golden comparisons; verified Task 3 hash, manifest, fixture, no-op normalization, and clean provenance state. |
| QA-65 | PASS | F-07/F-09 four decoded workflow-list formats equivalent. |
| QA-66 | PASS | F-07 committed literal markdown bytes match exactly. |
| QA-67 | PASS | F-07 filtered/unfiltered active ID sets PASS. |
| QA-68 | PASS | F-07 combined kind/active exact intersection PASS. |
| QA-69 | PASS | F-07/F-09 project `assigned_to` omits `kind`. |
| QA-70 | PASS | F-07 unassigned `assigned_to` is present empty array. |
| QA-71 | PASS | F-07 documented ordering and repeated byte identity PASS. |
| QA-72 | PASS | F-07 four no-match outputs exit 0 as empty arrays. |
| QA-73 | PASS | F-10 exact SQL returns zero reserved rows (`[]`). |
| QA-74 | PASS | F-10 exact reserved-rule production inventory exit 1/no output. |
| QA-75 | PASS | F-10 exact production `FrozenDefinitions(` inventory exit 1/no output. |

Checklist summary: **72 PASS, 3 FAIL, 0 BLOCKED**. The deterministic FAILs
are QA-29, QA-30, and QA-31. They block the phase independently of VAL-026.

## Command And Feasibility Validation Result

- Command and feasibility validation completed before implementation: **YES**,
  as recorded in QA plan §6.
- Final QA revalidated every exact task/phase command scope, report/spec pair,
  test package, fixture path, provenance artifact, static inventory, and Human
  step: **YES**.
- Required command placeholder found: **none**.
- Invalid/infeasible historical procedure: the Review 2 root-level
  `sha256sum -c` instruction was invalid. Review 2 correctly BLOCKED it;
  DEC-103 corrected the record and CF-025 replaced the ceremony with F-06.
- Currently invalid or unvalidated final-review command/fixture: **none**.
- Human feasibility: the nine steps are concrete and shipped, but outcomes are
  intentionally pending Human execution.

## Mid-Phase Review Disposition

| Review | Historical verdict | Blocking fixed? | Improvement disposition | References |
|---|---|---|---|---|
| `review-1.md` | BLOCK | YES | Task 4 baseline selector corrected before closure; immutable Task 3 hash and provenance verified here. | DEC-098; Task 4 closure |
| `review-2.md` | BLOCK | YES | History retained: the prescribed manifest command had never passed and Task 4 report-lint failed. DEC-103 corrected the evidence rule; Task 4 closure report `ba65078c` is lint-clean. | DEC-101..103; ESC-19-007 |

Review 2's two improvements are permanent tests, not one-time observations:

1. `cmd/fid/workflow_assign_test.go` asserts the six assign rejection messages
   against their exact expected values and proves pairwise distinctness.
2. `internal/workflow/assign_test.go` rebuilds the journal projection and
   asserts replayed kind/project assignment rows and effect shapes.

CF-025 is permanently implemented by
`cmd/fid/v01_parity_test.go` under the integration tag. DEC-103 explicitly
retires manual VAL-059 execution as required procedure.

## Documentation Review

- Impacted command/help documentation updated: **YES** for PRD, design §3.1,
  CLI spec, workflow glossary terms, and TRD §4.6.
- Full canonical synchronization complete: **NO**.
- Specs/current records needing Tech Lead correction:
  - reconcile `02_TRD.md` §3/§3.1 and DEC-089/DEC-096 with the actual
    `internal/workflow` write path;
  - update `03_TRD_SCHEMA.md` §14.11 for the assignment-integrity audit;
  - refresh `phase_brief.md`, `progress_current.md`, and CF-025/current decision
    state after Task 6;
  - record the roadmap Phase 19 closure disposition.
- Completion reports present for all tasks: **YES, 6/6 committed and clean**.
- Task 4 commit metadata correction requires Tech Lead/Human disposition; QA
  did not rewrite product history.

## Human Checklist Readiness

Generated at
`project-docs/docs/phases/phase-19/qa/human-checklist.md`.

- Contains only `(human)` items: **YES**.
- Contains exactly HUMAN-26a through HUMAN-26i: **YES, nine rows**.
- Uses one product/perception judgment per row: **YES**.
- Human-only VAL-026 mapped on every row: **YES**.
- Current status: **0 PASS, 0 FAIL, 0 PARTIAL, 9 blank; VAL-026 BLOCKED**.
- Product risks to call out: help-only discoverability, four-format coherence,
  corrective error quality, and whether DEC-079 feels understandable rather
  than surprising.

## Coverage Caveat

The 69 deterministic PASS outcomes prove only the asserted mechanics, data
states, command outputs, raw-byte compatibility, static inventories, and
structural conformance. They do not prove operator coherence, discoverability,
error tone, or the understandability of DEC-079. Those remain unconfirmed
until the Human records PASS on all nine checklist rows.

## Release Maturity Status

| Maturity area | Status | Evidence or deferral |
|---|---|---|
| Baseline artifacts | BLOCKED | Artifacts exist, but VAL-026 and QA-29/30/31 are not green. |
| Automated gates | READY | Exact eight-command phase gate is 8/8 PASS; deterministic product assertions 69/69 PASS. |
| Held-out scenarios | DEFERRED | No accepted harness exists; none was read. Human dogfood is not represented as held-out evidence. |
| Release readiness | DEFERRED | Phase 19 is not an external-release phase; security/dependency/release commands remain governed by the roadmap and quality plan. |

## Findings

### Blocking

- **QA-BLK-FINAL-01 — Human VAL-026 unexecuted.** Reproduction: inspect
  `human-checklist.md`; all nine outcome cells are blank. Required closure:
  Human PASS on HUMAN-26a..HUMAN-26i, or defect handling for any non-PASS.
- **QA-BLK-FINAL-02 — Canonical architecture contradiction.** Reproduction:
  compare `02_TRD.md` §3 (workflow owns write paths) and DEC-089 with §3.1 and
  DEC-096 (journal dependency only for `Effect`), then inventory
  `internal/workflow/define.go` and `assign.go` for `WithLock`, `PrepareEvent`,
  `DryRunProjectEvent`, `AppendEvent`, and `ProjectEvent`. Required closure:
  Tech Lead records one coherent boundary and synchronizes canonical specs.
- **QA-BLK-FINAL-03 — Canonical maintenance/current/roadmap docs stale.**
  Reproduction: compare `03_TRD_SCHEMA.md` §14.11's v0.1-only exclusion with
  `internal/maintenance/workflow_assignments.go`; inspect current phase status
  and roadmap closure. Required closure: Tech Lead spec/current-context sync.
- **QA-BLK-FINAL-04 — Task 4 commit lacks phase/task reference.** Reproduction:
  `git show -s --format='%B' 32614bbbd50e683b61921e9ed6beabee73b46021`;
  no `Refs:` line is present. Required closure: Tech Lead/Human disposition
  consistent with immutable-history policy.
- **QA-BLK-FINAL-05 — In-scope structure budgets remain exceeded.**
  Reproduction: `wc -l internal/issues/create.go` → `412`; function inventory
  places `Create` at line 19 and the next function at line 247. Required
  closure: separately authorized behavior-preserving decomposition and fresh
  regression evidence; QA does not implement it.

### Major

- None beyond the blocking findings above.

### Improvement

- The 100k rebuild control now has about 20% headroom (`23.935 s` against
  `<30 s`), improved from Phase 18, but remains worth watching as v0.2b adds
  receipt/evidence volume. The target must remain unchanged.

## Unresolved Risks And Carry-Forward Candidates

- **No CARRY-FORWARD outcome is assigned.** Heavy mode requires the five
  blocking findings to return to Human/Tech Lead resolution and a subsequent
  final QA disposition.
