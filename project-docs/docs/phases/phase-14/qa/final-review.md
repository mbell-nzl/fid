# Phase 14 Final Review - Telemetry + Pricing

Reviewer: QA Lead
Date: 2026-06-28

## Verdict

PASS WITH IMPROVEMENTS

Deterministic Phase 14 QA is green and no new blocker exists. Tasks 1-3 are
complete, Review 1's Task 2 integration-gate block was remediated and cleared,
and the phase is ready for Tech Lead wrap-up. VAL-009 remains a human-only
dogfood check and is carried forward to the generated human checklist.

## Scope And Evidence

- Phase brief: `project-docs/docs/current/phase_brief.md`
- Validation contract:
  `project-docs/docs/phases/phase-14/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-14/qa/qa-plan.md`
- Task specs and completion reports:
  `project-docs/docs/phases/phase-14/tasks/task-1.md`
  `project-docs/docs/phases/phase-14/tasks/task-1-done.md`
  `project-docs/docs/phases/phase-14/tasks/task-2.md`
  `project-docs/docs/phases/phase-14/tasks/task-2-done.md`
  `project-docs/docs/phases/phase-14/tasks/task-3.md`
  `project-docs/docs/phases/phase-14/tasks/task-3-done.md`
- Mid-phase review: `project-docs/docs/phases/phase-14/qa/review-1.md`
- Triage and escalations:
  `project-docs/docs/phases/phase-14/triage-log.md`,
  `project-docs/docs/current/tl_escalations_current.md`
- Current decisions and progress:
  `project-docs/docs/current/decisions_active.md`,
  `project-docs/docs/current/progress_current.md`
- Quality plan: `project-docs/specs/08_QUALITY_PLAN.md`

## Phase-Level Gate Results

| Gate | Command | Result | Evidence |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS | Exit 0 with no file output. |
| Lint/static analysis | `golangci-lint run ./...` | PASS | Output: `0 issues.` |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS | Output includes `ok github.com/mbell-nzl/fid/cmd/fid ... coverage: 77.4%`, `internal/config 80.7%`, `internal/runs 75.1%`, `internal/stats 81.6%`; all packages passed. |
| Coverage report | `go tool cover -func=coverage.out` | PASS | Output includes `total: (statements) 75.7%`; `DefaultPricingConfig 100.0%`, `computeCost 100.0%`, `internal/stats` query helpers 100.0%. |
| Integration | `go test -tags integration -race ./...` | PASS | All module packages passed; `cmd/fid`, `internal/runs`, and `internal/stats` reported `ok`. |
| End-to-end | N/A | N/A | AGENTS.md records end-to-end tests as `N/A` until an accepted harness exists. |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS | Output includes `BenchmarkRebuildWorkspace1kEvents ... 370.9 ms/rebuild` and `BenchmarkRebuildWorkspace10kEvents ... 3600 ms/rebuild`; all packages passed. |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | Exit 0, produced `/tmp/fid`. |
| Full phase check | QA plan section 8 seven-command phase gate | PASS | All required phase-end commands above passed exactly. |

## Focused Deterministic Evidence

Focused VAL run:

```text
go test -count=1 -v ./internal/config ./internal/runs ./internal/policy ./internal/stats ./cmd/fid -run 'TestDefaultPricingConfigIncludesAnthropicOpenAIAndGoogleModels|TestPricingConfigContentsRoundTripsDefaultPricingWithoutPrecisionLoss|TestWriteBootstrapFilesRoundTripsDefaultConfig|TestRecordKnownModelPersistsAgentRunAndComputedCost|TestRecordUnknownModelStoresNullCostAndWarning|TestRecordProviderFieldsPersistWithoutChangingNormalizedSummary|TestMatrixRunRecordRowMatchesSchemaSection13|TestRunRecordCreatesAgentRunAndPrintsConfirmation|TestRunRecordUnknownModelWarnsAndStillRecords|TestRunRecordRejectsUnknownActorRoleBeforeAppend|TestRunRecordPolicyShowMatchesAllowedRoles|TestTokensGroupsByRole|TestTokensGroupsByModelAndEpic|TestCostGroupsByEpicAndExcludesNullCostsFromSum|TestStatsTokensDefaultsToJSONByRole|TestStatsTokensSupportsFormatsAndEpicFilter|TestStatsCostSupportsFormatsAndExcludesNullCostsFromSum'
```

Result: PASS. Output included PASS for all named tests, including the seven
`TestRunRecordPolicyShowMatchesAllowedRoles/<role>` subtests and stats
format-parity subtests for `toon` and `yaml`; markdown parity is asserted in
the parent stats tests.

Temp `fid init` evidence:

```text
/tmp/fid init
initialized fid project in .fid/
  pricing: .fid/pricing.toml
  roles seeded: developer, qa_reviewer, qa_lead, tl_lite, tl, workflow_runner, human
rg pricing excerpt:
8:[models.claude-sonnet-4-6]
16:[models.gpt-5.4-mini]
17:in_per_mtok = 0.375
22:out_per_mtok = 0.625
24:[models.gemini-2.5-flash]
28:[models.gemini-2.5-flash-lite]
```

File-state evidence:

```text
git rev-parse --verify 2931031e0e4481f683f816cc411175ee24aaa526^{commit}
2931031e0e4481f683f816cc411175ee24aaa526
git rev-parse --verify 0cba161ed3b8bf6168af0399e229a518f23a8400^{commit}
0cba161ed3b8bf6168af0399e229a518f23a8400
git rev-parse --verify 25d3c0688ac8f6b0a79336e5ad35dbc22ebaf633^{commit}
25d3c0688ac8f6b0a79336e5ad35dbc22ebaf633
git diff --name-only -- migrations internal/journal/project_effects.go
<no output>
```

## Validation Contract Outcomes

| VAL ID | Outcome | Validator | Evidence |
|---|---|---|---|
| VAL-001 | PASS | QA Lead | Task 1 report records provider coverage, `FormatFloat(value, 'f', -1, 64)`, and `/tmp/fid init` proof. Focused tests passed: `TestDefaultPricingConfigIncludesAnthropicOpenAIAndGoogleModels`, `TestPricingConfigContentsRoundTripsDefaultPricingWithoutPrecisionLoss`, `TestWriteBootstrapFilesRoundTripsDefaultConfig`. Temp CLI evidence above shows Anthropic/OpenAI/Google entries and fractional prices preserved. |
| VAL-002 | PASS | QA Lead | Task 2 report and Review 1 cover one `run.recorded` event, one `agent_runs` row, timestamps, snapshot prices, and no `activity_log`. Focused tests passed: `TestRecordKnownModelPersistsAgentRunAndComputedCost`, `TestRunRecordCreatesAgentRunAndPrintsConfirmation`; test body queries `agent_runs`, counts journal lines before/after, and verifies `activity_log` count unchanged. |
| VAL-003 | PASS | QA Lead | `TestRecordKnownModelPersistsAgentRunAndComputedCost` passed and computes expected cost as `(1200*price_in + 3400*price_out)/1_000_000` from same-row snapshot prices with `1e-12` tolerance. |
| VAL-004 | PASS | QA Lead | `TestRecordUnknownModelStoresNullCostAndWarning` and `TestRunRecordUnknownModelWarnsAndStillRecords` passed. Tests assert unknown model `mystery-model` warns, records successfully, and stores NULL `price_in_per_mtok`, `price_out_per_mtok`, and `cost_usd`; known-model tests assert populated cost/prices and no warning path. |
| VAL-005 | PASS | QA Lead | `TestRecordProviderFieldsPersistWithoutChangingNormalizedSummary` passed. It records supplied and omitted provider-field states, verifies supplied `token_breakdown_json`, `tool_call_count`, `provider_run_id`, and `provider_cost_usd`, verifies NULLs when omitted, and compares normalized token/price/cost fields. |
| VAL-006 | PASS | QA Lead | `TestMatrixRunRecordRowMatchesSchemaSection13`, `TestRunRecordPolicyShowMatchesAllowedRoles`, and `TestRunRecordRejectsUnknownActorRoleBeforeAppend` passed. Evidence covers exactly one all-role `run.record` row, `policy show` parity, all seven default roles allowed, and unknown role `intern` rejected with unchanged journal length. |
| VAL-007 | PASS | QA Lead | `TestTokensGroupsByRole`, `TestTokensGroupsByModelAndEpic`, `TestStatsTokensDefaultsToJSONByRole`, and `TestStatsTokensSupportsFormatsAndEpicFilter` passed. Evidence covers role/model/epic grouping, `--epic` filter, hand-computed totals, JSON default, and JSON/TOON/Markdown/YAML structural parity. |
| VAL-008 | PASS | QA Lead | `TestCostGroupsByEpicAndExcludesNullCostsFromSum` and `TestStatsCostSupportsFormatsAndExcludesNullCostsFromSum` passed. Evidence covers cost grouped by epic and role, `--epic` filter, NULL-cost unknown-model run counted in tokens/run_count but excluded from cost sum, and JSON/TOON/Markdown/YAML parity. |
| VAL-009 | CARRY-FORWARD | Human | Human-only dogfood assertion. Generated `project-docs/docs/phases/phase-14/qa/human-checklist.md` with one row per observable: HUMAN-14a, HUMAN-14b, HUMAN-14c. Not marked PASS because human execution evidence is not present in this QA Lead run. |

## Requirement Compliance

| Requirement | Status | Evidence |
|---|---|---|
| PRD-NFR-010 | MET | VAL-001 PASS; Task 1 docs update in `04_TRD_CONFIG.md`; temp `pricing.toml` evidence. |
| PRD-TELEM-001 | MET | VAL-002 PASS; `fid run record` tests and completion report prove row/event fields. |
| PRD-TELEM-002 | MET | VAL-003 PASS; same-row snapshot-price formula tested. |
| PRD-TELEM-003 | MET | VAL-004 PASS; unknown-model warning and NULL pricing/cost tested alongside known-model state. |
| PRD-TELEM-003a | MET | VAL-005 PASS; provider reconciliation fields supplied/omitted states tested. |
| PRD-TELEM-004 | MET | VAL-007 and VAL-008 PASS; stats tokens/cost group by role/model/epic and support `--epic`. |
| PRD-OUT-001..004 | MET | VAL-007 and VAL-008 PASS; JSON, TOON, Markdown, and YAML parity tested for both stats commands. |
| TRD-ARCH-010 / `03_TRD_SCHEMA.md` section 13 | MET | VAL-006 PASS; in-code matrix row matches all-role `run.record` spec and show/enforce parity. |
| `03_TRD_SCHEMA.md` section 7 | MET | VAL-002..005 PASS; no new migration required or added. |
| `04_TRD_CONFIG.md` section 3 | MET | VAL-001, VAL-003, VAL-004 PASS; config doc updated for defaults, precision, and missing-model behavior. |
| `05_DESIGN_PLAN.md` section 6.4 | MET | Task 3 updated stats command tree; VAL-008 PASS. |
| `06_DESIGN_PLAN_CLI.md` `fid run record`, `stats tokens`, `stats cost` | MET | Task 2 and Task 3 completion reports cite CLI doc reconciliation; `rg` confirmed sections for `fid run record`, `fid stats tokens [--by role\|model\|epic]`, and `fid stats cost`. |
| `08_QUALITY_PLAN.md` sections 5.2, 5.6 | MET | Full phase gate passed; deterministic checklist executed; human checklist generated separately. |

## Deterministic Checklist Execution

| QA ID | Result | Evidence |
|---|---|---|
| QA-01 | PASS | Task 1, Task 2, and Task 3 completion reports exist and record exact required gates as PASS. Task 2 integration-gate omission was remediated in commit `0cba161ed3b8bf6168af0399e229a518f23a8400`; Review 1 addendum and triage TRG-002/TRG-003 mark it resolved. |
| QA-02 | PASS | Full phase gate executed: `gofmt`, `golangci-lint`, race coverage, coverage report, integration, benchmarks, and build all PASS. |
| QA-03 | PASS | Focused config tests passed; default coverage includes Anthropic/OpenAI/Google and full emit -> `ReadPricingConfig` round-trip without precision loss. |
| QA-04 | PASS | `/tmp/fid init` temp workspace evidence shows `.fid/pricing.toml` entries for `claude-sonnet-4-6`, `gpt-5.4-mini`, `gemini-2.5-flash`, `gemini-2.5-flash-lite`, with `0.375` and `0.625` preserved. |
| QA-05 | PASS | `TestRecordKnownModelPersistsAgentRunAndComputedCost` and `TestRunRecordCreatesAgentRunAndPrintsConfirmation` passed; assertions query `agent_runs`, verify role/model/tokens/timestamps/snapshot prices, and confirm one row. |
| QA-06 | PASS | `TestRecordKnownModelPersistsAgentRunAndComputedCost` passed; test records `activity_log` count before and after `Record()` and requires it unchanged. |
| QA-07 | PASS | `TestRecordKnownModelPersistsAgentRunAndComputedCost` passed; test calculates `cost_usd` from same-row snapshot prices within `1e-12`. |
| QA-08 | PASS | Unknown and known states covered by passing run tests; unknown model warns and stores NULL pricing/cost, known model records populated pricing/cost and prints cost summary. |
| QA-09 | PASS | `TestRecordProviderFieldsPersistWithoutChangingNormalizedSummary` passed; supplied and omitted provider-field rows compared. |
| QA-10 | PASS | `TestMatrixRunRecordRowMatchesSchemaSection13` and `TestRunRecordPolicyShowMatchesAllowedRoles` passed; all seven roles allowed in matrix and policy output. |
| QA-11 | PASS | `TestRunRecordPolicyShowMatchesAllowedRoles` subtests passed for developer, qa_reviewer, qa_lead, tl_lite, tl, workflow_runner, and human; `TestRunRecordRejectsUnknownActorRoleBeforeAppend` passed and asserts unchanged journal length. |
| QA-12 | PASS | `TestTokensGroupsByRole` and `TestTokensGroupsByModelAndEpic` passed; expected references include role, model, epic, null epic bucket, and `--epic` filter totals. |
| QA-13 | PASS | `TestStatsTokensSupportsFormatsAndEpicFilter` passed; JSON compared with TOON/YAML decoded structures and exact Markdown output. |
| QA-14 | PASS | `TestCostGroupsByEpicAndExcludesNullCostsFromSum` and CLI cost stats tests passed; expected references include role/epic grouping and `--epic` filtering. |
| QA-15 | PASS | `TestCostGroupsByEpicAndExcludesNullCostsFromSum` passed; unknown-model run contributes 80 in / 20 out and run_count 1 with `total_cost_usd` 0 for developer in filtered epic view. |
| QA-16 | PASS | `TestStatsCostSupportsFormatsAndExcludesNullCostsFromSum` passed; JSON compared with TOON/YAML decoded structures and exact Markdown output. |
| QA-17 | PASS | `git diff --name-only -- migrations internal/journal/project_effects.go` returned no output. |
| QA-18 | PASS | Documentation evidence present: `02_TRD.md` records `internal/runs`; `04_TRD_CONFIG.md` records pricing defaults and precision; `05_DESIGN_PLAN.md` and `06_DESIGN_PLAN_CLI.md` record stats/run CLI updates. |
| QA-19 | PASS | Phase coverage output: `internal/config 80.7%`, `internal/runs 75.1%`, `internal/policy 84.6%`, `internal/stats 81.6%`, `cmd/fid 77.4%`; total `75.7%`. |

## Command And Feasibility Validation Result

- Command and feasibility validation completed before implementation: YES.
- Invalid or unvalidated command, test scope, fixture, artifact, or human
  checklist step found: none.
- Phase-end revalidation: all referenced commands executed successfully except
  end-to-end/scenario validation, which is explicitly `N/A` until an accepted
  harness exists.

## Mid-Phase Review Disposition

| Review | Verdict | Blocking Fixed? | Improvement Disposition | ESC References |
|---|---|---|---|---|
| `review-1.md` | CONTINUE | YES | No improvement findings. Task 2 integration-gate artifact block cleared by remediation addendum. | ESC-001 resolved; TRG-002/TRG-003 resolved. |

## Documentation Review

- Impacted docs updated: YES.
- Specs needing Tech Lead wrap-up sync: none identified by QA.
- Completion reports present for all tasks: YES.
- QA artifacts generated by this review:
  `project-docs/docs/phases/phase-14/qa/final-review.md` and
  `project-docs/docs/phases/phase-14/qa/human-checklist.md`.
- Working-tree note: before this review, `tl_escalations_current.md`,
  `task-3.md`, `qa-plan.md`, `review-1.md`, and `triage-log.md` already had
  uncommitted or untracked states. This review did not modify those files.

## Human Checklist Readiness

Generated at
`project-docs/docs/phases/phase-14/qa/human-checklist.md`.

- Contains only `(human)` items: YES.
- Uses one observable per row: YES.
- Human-only `VAL-NNN` assertions mapped: YES, VAL-009 maps to HUMAN-14a,
  HUMAN-14b, and HUMAN-14c.
- Product/perception/dogfooding risks to call out: human dogfood has not been
  run in this QA Lead pass, so VAL-009 remains CARRY-FORWARD until the human
  checklist records outcomes.

## Coverage Caveat

A deterministic PASS proves only the asserted mechanics, data states, command
outputs, and structural conformance checks. It does not prove unasserted
workflow feel, product-owner judgment, or dogfooding quality; those remain
unconfirmed until the human checklist records outcomes.

## Release Maturity Status

| Maturity area | Status | Evidence or deferral |
|---|---|---|
| Baseline artifacts | READY | Phase brief, validation contract, QA plan, task specs/reports, Review 1, final review, and human checklist are present. |
| Automated gates | READY | Full phase gate passed; lint, race coverage, integration, benchmarks, build all green. |
| Held-out scenarios | DEFERRED | AGENTS.md and `08_QUALITY_PLAN.md` record scenario harness as `N/A` until accepted. No held-out scenario material was read. |
| Release readiness | DEFERRED | Security/dependency/license/secret release-readiness commands are not required before external release; no new dependency added in Phase 14. |

## Findings

### Blocking

- None.

### Major

- None.

### Improvement

- None new.

## Unresolved Risks And Carry-Forward Candidates

- VAL-009 remains CARRY-FORWARD to human dogfood checklist rows HUMAN-14a
  through HUMAN-14c.
- TRG-001 remains RESUME: Task 1 intentionally preserved existing Anthropic
  defaults rather than re-baselining them to current list prices. This does not
  block VAL-001 because exact list-price values are explicitly out of contract.
- Existing carry-forward items CF-008, CF-011, CF-013, CF-014, and CF-015 remain
  outside Phase 14 scope.
