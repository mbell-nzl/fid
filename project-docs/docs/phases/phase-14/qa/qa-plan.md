# Phase 14 QA Plan - Telemetry + Pricing

> Phase: 14 - Telemetry + pricing (v0.1e slice 1/4)
> Refreshed: 2026-06-28
> QA Lead: Codex QA Lead
>
> Location:
> `project-docs/docs/phases/phase-14/qa/qa-plan.md`.

## 1. Scope And Source Documents

Phase 14 adds the first telemetry and pricing slice for v0.1e:

- Task 1: expand `.fid/pricing.toml` shipped defaults to Anthropic, OpenAI,
  and Google models, and preserve sub-cent price precision.
- Task 2: add `fid run record`, including one `agent_runs` insert, snapshot
  price and cost computation, unknown-model warning behavior, provider
  reconciliation fields, and the `run.record` policy row.
- Task 3: add `fid stats tokens` and `fid stats cost` grouped by role, model,
  and epic in json, toon, markdown, and yaml.

Explicit scope boundaries from the phase brief and validation contract:

- No schema migration is required. The `agent_runs` table and columns already
  exist in `migrations/0002_v01_tables.up.sql` and
  `03_TRD_SCHEMA.md` section 7.
- No projector change is expected. Generic insert projection already covers
  the `agent_runs` insert effect.
- Phase 14 does not add `activity_log` rows for telemetry; CF-011 remains
  Phase 16 scope.
- Phase 14 does not implement run list/show, run status lifecycle, git context,
  provider-billed cost reports, `fid sql`, or stats expansions beyond tokens
  and cost.

Source artifacts used:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/01_PRD.md`
- `project-docs/specs/02_TRD.md`
- `project-docs/specs/03_TRD_SCHEMA.md`
- `project-docs/specs/04_TRD_CONFIG.md`
- `project-docs/specs/05_DESIGN_PLAN.md`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md`
- `project-docs/docs/phases/phase-14/qa/validation-contract.md`
- `project-docs/docs/phases/phase-14/tasks/task-1.md`
- `project-docs/docs/phases/phase-14/tasks/task-2.md`
- `project-docs/docs/phases/phase-14/tasks/task-3.md`
- `project-docs/templates/qa-plan-template.md`

## 2. Requirement Coverage

| Requirement | Phase coverage | VAL IDs |
|---|---|---|
| PRD-NFR-010 | Shipped pricing defaults cover current Anthropic, OpenAI, and Google models and remain user-overridable. | VAL-001 |
| PRD-TELEM-001 | `fid run record` captures role, model, tokens, timestamps, and snapshot pricing. | VAL-002 |
| PRD-TELEM-002 | `fid run record` computes USD cost from row snapshot prices. | VAL-003 |
| PRD-TELEM-003 | Unknown models leave cost/prices NULL and emit a warning while still recording the run. | VAL-004 |
| PRD-TELEM-003a | Optional provider token/cost reconciliation fields persist without replacing normalized summary fields. | VAL-005 |
| PRD-TELEM-004 | Token and cost statistics aggregate by role, model, and epic. | VAL-007, VAL-008 |
| PRD-OUT-001..004 | Stats read commands support json default plus toon, markdown, and yaml with equivalent data. | VAL-007, VAL-008 |
| TRD-ARCH-010; `03_TRD_SCHEMA.md` section 13 | `run.record` appears once in the in-code permission matrix, show and enforce use the same source, and unknown roles reject before append. | VAL-006 |
| `03_TRD_SCHEMA.md` section 7 | `agent_runs` is the existing target projection table for run recording, including provider reconciliation fields. | VAL-002, VAL-003, VAL-004, VAL-005 |
| `04_TRD_CONFIG.md` section 3 | `.fid/pricing.toml` drives snapshot pricing and missing-model cost NULL behavior. | VAL-001, VAL-003, VAL-004 |
| `05_DESIGN_PLAN.md` section 6.4 | `fid stats cost` exposes cost totals, tokens, and run counts. | VAL-008 |
| `06_DESIGN_PLAN_CLI.md` `fid run record` | Run-record command spelling, required flags, optional provider fields, and stderr behavior. | VAL-002, VAL-003, VAL-004, VAL-005 |
| `06_DESIGN_PLAN_CLI.md` `fid stats tokens` / `stats cost` | Stats command spellings, grouping flags, epic filter, and format flags. | VAL-007, VAL-008 |
| `08_QUALITY_PLAN.md` sections 5.2, 5.6 | Deterministic phase evidence, command feasibility evidence, and human checklist separation. | VAL-009, no VAL gates |

## 3. Validation Contract Refinement Check

| Check | Result | Notes |
|---|---|---|
| All in-scope requirements referenced by at least one `VAL-NNN` | YES | See requirement coverage table. |
| Every `VAL-NNN` has one validator role | YES | VAL-001..008: QA Lead; VAL-009: Human. |
| Every `VAL-NNN` has a concrete evidence method | YES | Refined in sections 5 and 10. |
| Every `VAL-NNN` maps to a task/window or phase-end check | YES | Task/window mapping in section 4. |
| Every `VAL-NNN` describes one observable, or human-checklist fan-out is explicitly required | YES | VAL-001 is pricing-default round-trip; VAL-002 is the run-record write contract; VAL-003 is cost formula; VAL-004 is known/unknown model binary behavior; VAL-005 is provider-field supplied/omitted binary behavior; VAL-006 is `run.record` policy parity and rejection; VAL-007 and VAL-008 are stats aggregate parity; VAL-009 fans out to HUMAN-14a..HUMAN-14c. |
| Conditional or binary `VAL-NNN` assertions require both qualifying and non-qualifying states | YES | VAL-001, VAL-004, VAL-005, VAL-006, VAL-007, and VAL-008 require both states or all named dimensions/formats. |
| Parity or "matches reference" assertions have deterministic conformance proof when available | YES | VAL-001 requires emit/read round-trip; VAL-006 requires section-13 matrix/show/enforce comparison; VAL-007 and VAL-008 require format parity against JSON/reference totals. |
| Revisions requested from Tech Lead | None | No uncovered VAL or infeasible evidence method found. Composite deterministic assertions are decomposed into evidence facets in section 10; partial evidence cannot pass the parent VAL. |

## 4. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | Pricing defaults plus precision | VAL-001 | Medium | Must prove all three providers and exact emit -> `ReadPricingConfig` round-trip, including sub-cent values. |
| Task 2 | `fid run record` write path and `run.record` policy row | VAL-002, VAL-003, VAL-004, VAL-005, VAL-006 | High | First telemetry write surface; must prove one event, no `activity_log`, lock-safe projection refresh, unknown-model warning, and validate-before-append for unknown role. |
| Task 3 | `fid stats tokens` and `fid stats cost` | VAL-007, VAL-008 | High | Must prove role/model/epic grouping, `--epic` filter, NULL-cost handling, deterministic ordering, and four-format parity. |
| Phase end | Final QA and human dogfood checklist | VAL-009; final outcomes for VAL-001..009 | Medium | Human checklist contains only dogfooding/product judgment rows; deterministic VALs remain QA Lead evidence obligations. |

## 5. Required Evidence Strategy

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | Unit output proving `DefaultPricingConfig` has at least one Anthropic, OpenAI, and Google model, and that every default entry round-trips through emit -> `ReadPricingConfig` with `in_per_mtok` and `out_per_mtok` unchanged. | VAL-001 | Developer, QA Lead | Task 1 report; review-1; final review |
| EVID-02 | Build/CLI output from `fid init` in a temp workspace plus parsed `.fid/pricing.toml` evidence showing generated defaults match `DefaultPricingConfig`. | VAL-001 | Developer, QA Lead | Task 1 report; review-1; final review |
| EVID-03 | Domain/CLI output proving known-model `fid run record` appends exactly one `run.recorded` event, projects one `agent_runs` row with role/model/tokens/timestamps/snapshot prices, and projects zero `activity_log` rows. | VAL-002 | Developer, QA Lead | Task 2 report; review-1; final review |
| EVID-04 | Cost evidence showing `cost_usd == (tokens_in * price_in_per_mtok + tokens_out * price_out_per_mtok) / 1e6` using snapshot prices stored on the same row, within stated float tolerance. | VAL-003 | Developer, QA Lead | Task 2 report; review-1; final review |
| EVID-05 | CLI output for both unknown and known model states: unknown exits 0, warns on stderr naming the model, stores NULL cost/prices; known stores populated cost/prices and emits no unknown-model warning. | VAL-004 | Developer, QA Lead | Task 2 report; review-1; final review |
| EVID-06 | Domain/CLI output for provider-field supplied and omitted states, proving supplied provider columns persist, omitted provider columns are NULL, and normalized token/price/cost fields are unchanged for identical inputs. | VAL-005 | Developer, QA Lead | Task 2 report; review-1; final review |
| EVID-07 | Matrix, policy-show, and CLI evidence proving exactly one `run.record` row matches `03_TRD_SCHEMA.md` section 13 all-role booleans, all seven default roles may record, and an unknown role rejects before append with a non-zero exit and an error naming role plus `run.record`. | VAL-006 | Developer, QA Lead | Task 2 report; review-1; final review |
| EVID-08 | CLI output or parsed test evidence proving `fid stats tokens --by role`, `--by model`, `--by epic`, and `--epic <id>` match hand-computed token totals and run counts. | VAL-007 | Developer, QA Lead | Task 3 report; final review |
| EVID-09 | Structural parity evidence proving `stats tokens` json, toon, markdown, and yaml carry the same grouped data as the JSON reference. | VAL-007 | Developer, QA Lead | Task 3 report; final review |
| EVID-10 | CLI output or parsed test evidence proving `fid stats cost --by role`, `--by model`, `--by epic`, and `--epic <id>` match hand-computed cost totals, token totals, and run counts, with NULL-cost runs excluded from cost sums but included in counts/tokens. | VAL-008 | Developer, QA Lead | Task 3 report; final review |
| EVID-11 | Structural parity evidence proving `stats cost` json, toon, markdown, and yaml carry the same grouped data as the JSON reference. | VAL-008 | Developer, QA Lead | Task 3 report; final review |
| EVID-12 | Human dogfood checklist rows confirming the user can identify run recording with computed USD cost, unknown-model warning success, and role/model/epic aggregates from CLI output without reading implementation. | VAL-009 | Human | Phase-end human checklist |
| EVID-13 | Exact task-gate command outputs, conditional command dispositions, and package coverage percentages requested by each task spec. | no VAL | Developer, QA Reviewer, QA Lead | Every task report; review-1; final review |
| EVID-14 | File-state evidence proving no migration or projector change was added unless a task records a spec-backed blocker/decision. | no VAL | QA Lead | review-1; final review |
| EVID-15 | Documentation evidence proving `02_TRD.md`, `04_TRD_CONFIG.md`, and `06_DESIGN_PLAN_CLI.md` updates required by the task specs are present or explicitly verified no-change. | no VAL | Developer, QA Lead | Task reports; final review |

## 6. Acceptance Criteria

- VAL-001 through VAL-008 have deterministic QA Lead evidence slots and cannot
  pass from code-reading alone.
- VAL-009 is human-only and fans out into HUMAN-14a through HUMAN-14c.
- Every task completion report records the exact required commands as `PASS`,
  or `BLOCKED` with a blocker reason. Missing, silently skipped, or narrower
  substituted commands are blocking QA findings.
- Task 1 records the chosen pricing emit precision format and the evidence for
  `fid init` generated `pricing.toml`.
- Task 2 records the `internal/runs` boundary decision in `02_TRD.md`, proves
  `run.record` policy parity, and proves unknown-role rejection before append.
- Task 3 records four-format parity for both stats commands and the `--by epic`
  documentation reconciliation in `06_DESIGN_PLAN_CLI.md`.
- Final QA confirms no schema migration, no projector change, and no
  out-of-task-scope implementation unless an escalation or decision explicitly
  records the deviation.

## 7. Expected Tests

| Area | Expected tests | Covers |
|---|---|---|
| `internal/config` | Tests for three-provider default coverage, stable existing Anthropic entries, and full default emit -> `ReadPricingConfig` round-trip including sub-cent prices. | VAL-001 |
| `cmd/fid/init.go` / init CLI tests | Test or CLI evidence that `fid init` writes `.fid/pricing.toml` whose parsed entries match `DefaultPricingConfig`. | VAL-001 |
| `internal/runs` | Record tests for known-model row fields, cost formula, one event, no `activity_log`, unknown-model NULL cost/prices, provider-field supplied/omitted states, and lock-safe append/project flow. | VAL-002, VAL-003, VAL-004, VAL-005 |
| `cmd/fid/run.go` | CLI tests for required flags, known-model success, unknown-model warning with exit 0, provider flags, all permitted roles, and unknown-role reject-before-append. | VAL-002, VAL-003, VAL-004, VAL-005, VAL-006 |
| `internal/policy` | Matrix completeness test for exactly one `run.record` row and role booleans matching `03_TRD_SCHEMA.md` section 13. | VAL-006 |
| `cmd/fid/policy.go` / policy tests | Policy-show output includes `run.record` from the same matrix enforced by `Authorize`. | VAL-006 |
| `internal/stats` | Aggregation tests for tokens and cost by role/model/epic, `--epic` filter, null/`(none)` epic bucket, deterministic ordering, and NULL-cost exclusion from sums while included in counts/tokens. | VAL-007, VAL-008 |
| `cmd/fid/stats.go` | CLI tests for `stats tokens` and `stats cost` flags plus json/toon/markdown/yaml parity against a JSON/reference dataset. | VAL-007, VAL-008 |

## 8. Exact Task Gates

Every command below is exact. Focused commands may be used during TDD, but they
do not substitute for the close gate.

| Task | Required gate commands | Conditional commands |
|---|---|---|
| Task 1 | `gofmt -l cmd internal`<br>`golangci-lint run ./...`<br>`go test -race -coverprofile=coverage.out ./...`<br>`go tool cover -func=coverage.out` | `go build -o /tmp/fid ./cmd/fid` is required for the `fid init` temp-workspace evidence. |
| Task 2 | `gofmt -l cmd internal`<br>`golangci-lint run ./...`<br>`go test -race -coverprofile=coverage.out ./...`<br>`go tool cover -func=coverage.out` | `go test -tags integration -race ./...` if any test is integration-tagged; `go build -o /tmp/fid ./cmd/fid` is required for CLI build evidence. |
| Task 3 | `gofmt -l cmd internal`<br>`golangci-lint run ./...`<br>`go test -race -coverprofile=coverage.out ./...`<br>`go tool cover -func=coverage.out` | `go test -tags integration -race ./...` if any test is integration-tagged; `go build -o /tmp/fid ./cmd/fid` is required for CLI build evidence. |

Conditional commands must be recorded as `PASS`, or as `N/A`/`BLOCKED` with the
task-specific reason. A required build used to produce CLI evidence is not
optional for that task's completion report.

Phase-end QA gate:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

## 9. Command And Feasibility Validation

| Referenced Step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | Help output and path existence. | PASS | `gofmt -h` shows `usage: gofmt [flags] [path ...]` and `-l`; `find cmd internal -maxdepth 0 -type d -print` returned `cmd` and `internal`. |
| `golangci-lint run ./...` | Help output and tool version. | PASS | `golangci-lint run --help` shows usage `golangci-lint run [flags]`; `golangci-lint version` returned version `2.12.2`, matching the v2 expectation in `08_QUALITY_PLAN.md`. |
| `go test -race -coverprofile=coverage.out ./...` | Go help output and package-list validation. | PASS | `go help testflag` documents `-coverprofile`; `go help build` documents `-race`; `go help packages` documents `./...`; `go list ./...` listed the current module packages. |
| `go tool cover -func=coverage.out` | Tool help output. | PASS | `go tool cover -h` documents `-func string` for coverage percentage output. The help command exits non-zero by convention after printing usage, but the flag is valid. |
| `go test -tags integration -race ./...` | Go help output. | PASS | `go help build` documents `-tags` and `-race`; `go help packages` documents `./...`. Existing integration-tag files are present under `cmd/fid/`, so the scope is meaningful. |
| `go test -bench . -benchmem ./...` | Go help output. | PASS | `go help testflag` documents `-bench regexp` and `-benchmem`; `go help packages` documents `./...`. |
| `go build -o /tmp/fid ./cmd/fid` | Go help output, package-list validation, and actual build. | PASS | `go help build` documents `-o output`; `go list ./cmd/fid` returned `github.com/mbell-nzl/fid/cmd/fid`; `go build -o /tmp/fid ./cmd/fid` completed with exit 0. |
| Task specs | File-state check. | PASS | `git ls-files` returned `task-1.md`, `task-2.md`, and `task-3.md` under `project-docs/docs/phases/phase-14/tasks/`. |
| QA templates | File-state check. | PASS | `find project-docs/templates -maxdepth 1 -type f -print` returned `qa-plan-template.md`, `mid-phase-review-template.md`, `final-review-template.md`, and `human-checklist-template.md`. |
| Validation contract artifact | File-state and contract inspection. | PASS | `git ls-files` returned `phase-14/qa/validation-contract.md`; the contract defines VAL-001 through VAL-009. |
| Source document artifact sizes | File-state check. | PASS | `wc -l` found the phase brief, validation contract, three task specs, and QA template present and readable. |
| Placeholder command scan | Spec/task inspection. | PASS | Phase 14 task specs use concrete Go/golangci-lint/build commands. `08_QUALITY_PLAN.md` contains release-readiness `N/A` items, but no placeholder command is required for Phase 14 execution. |
| Temp-workspace CLI evidence | Build target and writable path validation. | PASS | `/tmp` is a writable root; `go build -o /tmp/fid ./cmd/fid` succeeded. Final QA can create temp `.fid` workspaces without writing outside allowed roots. |
| Pricing fixture assumption | Existing source path check. | PASS | `internal/config/config.go` and `internal/config/config_test.go` exist; Task 1 names those as the pricing default and round-trip test surface. |
| `agent_runs` table assumption | Schema/task inspection. | PASS | Validation contract cites `migrations/0002_v01_tables.up.sql` and `03_TRD_SCHEMA.md` section 7 as existing table/column sources; no migration task is planned. Final QA must verify no new migration was added unless a blocker/decision records why. |
| `run.record` matrix reference | Contract/task inspection. | PASS | Task 2 and VAL-006 require the in-code row to copy `03_TRD_SCHEMA.md` section 13 booleans verbatim and prove show/enforce parity. |
| Stats format parity assumption | Existing format package check. | PASS | `go list ./...` includes `github.com/mbell-nzl/fid/internal/format`; existing format tests are present under `internal/format/`. Task 3 must reuse this surface. |
| Human checklist feasibility | Template and VAL inspection. | PASS | `human-checklist-template.md` exists; VAL-009 explicitly requires one row per human-observable dogfood outcome after deterministic QA passes. |

## 10. Deterministic Checklist

Every item is executed by QA Lead during final review unless assigned to the
mid-phase review window. PASS requires reproducible evidence in the review
artifact, such as command output, file state, parsed CLI output, or DB/query
state. Code-reading alone is not sufficient.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Verify each task completion report records every exact required command for that task as `PASS`, or `BLOCKED` with a reason; reject missing, silently skipped, or narrower substituted commands `(qa-lead)` | no VAL | File-state proof from Task 1-3 completion reports plus copied command-result excerpts. |
| QA-02 | Run the full phase-end gate exactly as listed in section 8 `(qa-lead)` | no VAL | Literal command output for all seven phase-end commands. |
| QA-03 | Verify `DefaultPricingConfig` covers at least one Anthropic, one OpenAI, and one Google model, and that all default prices round-trip through emit -> `ReadPricingConfig` without precision loss `(qa-lead)` | VAL-001 | Unit-test output and parsed round-trip evidence. |
| QA-04 | Verify `/tmp/fid init` in a temp workspace writes `.fid/pricing.toml` whose parsed entries match `DefaultPricingConfig` `(qa-lead)` | VAL-001 | Build output, temp-workspace command output, and parsed/file-state excerpt. |
| QA-05 | Verify known-model `fid run record` appends exactly one `run.recorded` journal event and projects exactly one `agent_runs` row with role, model, tokens, timestamps, and snapshot prices `(qa-lead)` | VAL-002 | CLI/test output plus journal-line and DB-state excerpt. |
| QA-06 | Verify `fid run record` writes zero `activity_log` rows for telemetry `(qa-lead)` | VAL-002 | DB/query or test output excerpt showing no new `activity_log` row. |
| QA-07 | Verify known-model `cost_usd` equals the formula over same-row snapshot prices within the documented tolerance `(qa-lead)` | VAL-003 | Test output or DB/query calculation excerpt. |
| QA-08 | Verify unknown-model and known-model behavior both states: unknown model exits 0, warns on stderr naming the model, stores NULL cost/prices; known model stores populated cost/prices and emits no unknown-model warning `(qa-lead)` | VAL-004 | CLI/test output for both states plus DB-state excerpts. |
| QA-09 | Verify optional provider fields persist when supplied, are NULL when omitted, and do not alter normalized tokens/prices/cost for identical inputs `(qa-lead)` | VAL-005 | CLI/test output and DB-state comparison for supplied and omitted states. |
| QA-10 | Verify the in-code permissions matrix has exactly one `run.record` row matching `03_TRD_SCHEMA.md` section 13 all-role booleans, and `fid policy show` displays the same row enforced by `Authorize` `(qa-lead)` | VAL-006 | Test output, policy-show CLI excerpt, and section-13 comparison evidence. |
| QA-11 | Verify all seven default-profile roles may record a run, and an unknown role is rejected before append with non-zero exit and an error naming role plus `run.record` `(qa-lead)` | VAL-006 | CLI/test output and journal-line count excerpt for allowed and rejected states. |
| QA-12 | Verify `fid stats tokens` grouped by role, model, and epic, plus `--epic <id>`, matches a hand-computed reference for `total_tokens_in`, `total_tokens_out`, and `run_count` `(qa-lead)` | VAL-007 | Parsed CLI/test output and reference totals. |
| QA-13 | Verify `fid stats tokens` json, toon, markdown, and yaml carry the same grouped data as the JSON/reference output `(qa-lead)` | VAL-007 | Structural parity comparison or literal parsed excerpts for all four formats. |
| QA-14 | Verify `fid stats cost` grouped by role, model, and epic, plus `--epic <id>`, matches a hand-computed reference for `total_cost_usd`, token totals, and `run_count` `(qa-lead)` | VAL-008 | Parsed CLI/test output and reference totals. |
| QA-15 | Verify NULL-cost runs are counted in `run_count` and token totals but excluded from `total_cost_usd` sums `(qa-lead)` | VAL-008 | Seeded fixture output showing known-cost and unknown-model rows in the same reference set. |
| QA-16 | Verify `fid stats cost` json, toon, markdown, and yaml carry the same grouped data as the JSON/reference output `(qa-lead)` | VAL-008 | Structural parity comparison or literal parsed excerpts for all four formats. |
| QA-17 | Verify no new migration file and no `internal/journal/project_effects.go` change landed unless a task records a spec-backed blocker or decision `(qa-lead)` | no VAL | `git diff --name-only`/file-state output and blocker/decision references if any. |
| QA-18 | Verify required documentation updates are present: `04_TRD_CONFIG.md` for Task 1, `02_TRD.md` and `06_DESIGN_PLAN_CLI.md` for Task 2, and `06_DESIGN_PLAN_CLI.md` for Task 3 `(qa-lead)` | no VAL | File-state proof, task-report excerpts, and spec diff references. |
| QA-19 | Verify completion reports record package coverage percentages requested by their task specs: `internal/config`; `internal/runs`, `internal/policy`, `cmd/fid`; `internal/stats`, `cmd/fid` `(qa-lead)` | no VAL | `go tool cover -func=coverage.out` excerpts from task reports and final review. |

## 11. Mid-Phase Review Cadence

- Cadence for this 3-task phase: one mandatory mid-phase review after Task 2
  completes.
- Review file: `project-docs/docs/phases/phase-14/qa/review-1.md`.
- Review window: Tasks 1-2.
- Covered assertions in the review window: VAL-001 through VAL-006.
- Required review checks:
  - Compare Task 1 and Task 2 completion-report command results with the exact
    required gates in section 8.
  - Verify Task 1 produced VAL-001 evidence, including `fid init` generated
    `pricing.toml` proof.
  - Verify Task 2 produced VAL-002 through VAL-006 evidence, including
    no-`activity_log`, unknown-model warning, provider-field supplied/omitted
    states, and `run.record` policy parity.
  - Verify no migration or projector change appears unless a task recorded a
    spec-backed blocker/decision.
  - Update or create the Phase 14 triage log if the review finds any findings;
    record CONTINUE/BLOCK per the mid-phase review template.
- No separate mid-phase review is required after Task 1 because Task 2 is the
  larger write-surface risk and the phase brief requires at least the Task 2
  checkpoint.
- Task 3 is reviewed at phase-end final review, where VAL-007, VAL-008, VAL-009,
  and the full phase gate are evaluated.

Missing, skipped, narrower, or silently substituted required commands are
blocking unless the task report records them as `BLOCKED` with a reason in a
blocker report.

## 12. Human Verification Scope

Items here require human dogfooding or product judgment and must be generated
into `project-docs/docs/phases/phase-14/qa/human-checklist.md` at phase end.
Each item is one observable from VAL-009.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-14a | Confirm from real CLI output in a real `.fid` workspace that an agent run was recorded with a visible computed USD cost `(human)` | VAL-009 | Human note or pasted CLI excerpt. |
| HUMAN-14b | Confirm from real CLI output that recording a run with an unknown model still succeeds and emits a warning that names the model `(human)` | VAL-009 | Human note or pasted command transcript with stderr warning. |
| HUMAN-14c | Confirm from real CLI output that `fid stats tokens` and `fid stats cost` report understandable aggregates by role, model, and epic `(human)` | VAL-009 | Human note or pasted CLI excerpts for the aggregate views. |

## 13. Phase-End Decision Rule

The phase can proceed to final sign-off only when:

- All Task 1-3 completion reports exist and record exact gate results.
- Review 1 exists after Task 2 and records outcomes for VAL-001 through
  VAL-006 or explicit blockers.
- The deterministic phase-end gate passes or is explicitly BLOCKED with
  escalation.
- Final review records `PASS`, `FAIL`, `BLOCKED`, or `CARRY-FORWARD` for every
  VAL-001 through VAL-009.
- VAL-001 through VAL-008 have cited deterministic evidence, not code-reading
  only claims.
- VAL-009 has a generated human checklist containing only HUMAN-14a through
  HUMAN-14c or narrower human-only rows.
- No schema migration, projector change, public-interface drift, or
  out-of-task-scope work is present without a recorded spec-backed decision.
- No unresolved blocking defects remain.
