# Phase 19 QA Plan - Workflow Authoring CLI + Escalation Default Workflow

> Phase: 19 - `fid workflow define/assign/list` + escalation default workflow
> Refreshed: 2026-07-26
> QA Lead: Codex QA Lead
> Operating mode: Heavy
> Plan status: **Ready**

DEC-089..DEC-093 dispositioned QA-BLK-01..QA-BLK-07 and the Tech Lead refined
the contract to 70 `Ready` assertions without changing task ownership or
dropping an observable. This refresh re-verifies that contract, adds
one-to-one evidence and checklist coverage for VAL-027..070, and resolves the
two formerly blocked evidence foundations (VAL-003/030 tuple parity and
VAL-059 immutable output parity). Per `09_AI_WORKFLOW.md` §2.6, the QA Lead
has not edited `validation-contract.md`.

## 1. Scope And Source Documents

Phase 19 makes projected workflow data authorable. It adds two migrations, the
`workflow.defined` and `workflow.assigned` journal event types, the public
`fid workflow define/assign/list` CLI, assignment-aware issue creation, the
`escalation-default` seed, assignment integrity checks, and closure evidence
for deterministic rebuild, performance, and operator dogfooding.

Sources loaded for this plan:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/docs/phases/phase-19/decisions.md` (DEC-075..DEC-093)
- `project-docs/docs/phases/phase-19/qa/validation-contract.md`
- `project-docs/docs/phases/phase-19/tasks/task-1.md` through `task-6.md`
- `project-docs/specs/01_PRD.md` §§3.5, 3.7, 3.13, 3.14
- `project-docs/specs/02_TRD.md` §§3, 4.3, 4.6, 7, 10
- `project-docs/specs/03_TRD_SCHEMA.md` §§2.1, 3, 13, 14.11
- `project-docs/specs/05_DESIGN_PLAN.md` §§1, 3.1, 4, 5, 8, 9
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` Conventions and
  `fid workflow define/assign/list`, `issue transition`, `rebuild`, and
  `maintenance check`
- `project-docs/specs/08_QUALITY_PLAN.md`, especially §§3-5 and §10
- `project-docs/specs/09_AI_WORKFLOW.md` §§2.6, 4, 6.1, 8
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` v0.2a and v0.2e slices
- `project-docs/specs/GLOSSARY.md` Projected workflow tables, Workflow
  registry, Frozen workflow definitions, and workflow terms
- `project-docs/templates/qa-plan-template.md`,
  `validation-contract-template.md`, `mid-phase-review-template.md`,
  `final-review-template.md`, and `human-checklist-template.md`
- `project-docs/docs/phases/phase-18/qa/final-review.md` benchmark evidence
  excerpt and `phase-18/tasks/task-5-done.md` benchmark excerpt

No Phase 19 completion report or mid-phase review exists yet, as expected
before Task 1. Their planned paths were checked in Section 6.

## 2. Validation Contract Refinement Check

- All in-scope requirements referenced by at least one `VAL-NNN`: **YES**.
- Every `VAL-NNN` has one validator role: **YES** (69 QA Lead assertions;
  VAL-026 Human).
- Every `VAL-NNN` has a concrete evidence method: **YES**. Section 4 contains
  one evidence row per assertion, and Section 6 records the reproducible
  foundations for VAL-003/030 and VAL-059.
- Every `VAL-NNN` maps to a task/window or phase-end check: **YES**.
- Every `VAL-NNN` describes one observable, or human-checklist fan-out is
  explicitly required: **YES**. DEC-091 narrowed the 14 identified bundled
  parents and added split assertions in the same task windows; VAL-026 is the
  sole Human assertion and remains fanned out in Section 9.
- Conditional or binary assertions require both qualifying and
  non-qualifying states: **YES**. In particular, DEC-090 makes VAL-038
  explicit: zero-cancelled-state definitions are accepted for
  `escalation`/`qa_review` and rejected for `epic`/`task`/`subtask`/`any`.
- Parity or "matches reference" assertions have deterministic conformance
  proof when available: **YES**. VAL-003 and VAL-030 use sorted tuple
  deep-equality, never counts; VAL-059 uses a verified immutable commit,
  committed fixed fixture/captures, SHA-256 manifest, and raw-byte comparison.
- Revisions requested from Tech Lead: **none**.

### Blocking Refinement And Specification Findings

| ID | Prior severity | Recorded closure | Decision |
|---|---|---|---|
| QA-BLK-01 | Blocking | **Closed:** `internal/workflow` owns authoring, assignment, and resolution; `02_TRD.md` §3 is synchronized and retains explicit must-nots. | DEC-089 |
| QA-BLK-02 | Blocking | **Closed:** PRD-WF-013 remains intact; issue-applicable definitions require exactly one `cancelled` state, while zero is legal only for non-issue applicability. | DEC-090 |
| QA-BLK-03 | Blocking | **Closed:** the 14 bundled parents were narrowed and split-out observables received new IDs in their original task windows. | DEC-091 |
| QA-BLK-04 | Blocking | **Closed:** uncovered binding behavior received one-for-one assertions through VAL-070. | DEC-091 |
| QA-BLK-05 | Blocking | **Closed:** sources were corrected; arrow aliases are explicit; VAL-003/030 and VAL-059 now require deterministic, reproducible references. Fresh feasibility proof is in Section 6. | DEC-091 |
| QA-BLK-06 | Blocking | **Closed:** `05_DESIGN_PLAN.md` §3.1 now contains `workflow define/assign/list`; Tasks 3-5 verify their entries and Task 6 records the final spec-sync verdict. | DEC-092 |
| QA-BLK-07 | Blocking | **Closed:** Human approval of heavy mode, six tasks, the 70-row `Ready` contract, and CF-008 disposition is recorded. | DEC-093 |

Contract reconciliation result: **70/70 assertions retained**. VAL-001..026
remain present (with the DEC-091 narrowing applied), VAL-027..070 each have
their own evidence row and QA step, task ownership remains 10/14/13/13/11/9,
and no assertion is retired, inherited, or omitted.

### Non-Blocking Applicability Clarification

`06_DESIGN_PLAN_CLI.md` states that every v0.2 write accepts `--dry-run`, but
`10_IMPLEMENTATION_ROADMAP.md` explicitly schedules universal `--dry-run` in
v0.2e. It is therefore not Phase 19 scope and is not a missing Phase 19
assertion. The final documentation review must preserve that staged wording
and must not imply define/assign support `--dry-run` in v0.2a.

## 3. Task Inventory

| Task | Focus | Covered current VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | `workflow_assignments`, escalation seed, projection version | VAL-001..005, VAL-027..031 (10) | High | Migration reversibility, exact seed/default tuple parity, constraint branches, behavior-free load, and stale/fresh projection branches. |
| Task 2 | Pure authoring and assignment validation/effect core | VAL-006..008, VAL-032..042 (14) | High | Rejection matrix, deterministic ordering, purity, cancel policy, initial state, assignment validation, and 80% `internal/workflow` threshold. |
| Task 3 | Define CLI, event, projection, permissions | VAL-009..013, VAL-043..050 (13) | High | Reject-before-append, event atomicity, projection/load, no activity row, permissions, aliases, output, and applicability vocabulary. |
| Task 4 | Assign CLI and assignment-aware issue creation | VAL-014..017, VAL-051..059 (13) | Very high | Most-used create path, last-write-wins effect shapes, precedence, rejection atomicity, and immutable pre-Task-4 parity evidence. Task 4 cannot start before Review 1. |
| Task 5 | Four-format workflow list | VAL-018..020, VAL-060..067 (11) | Medium-high | Cross-format structural parity, exact markdown characterization, filters in both states, assignment shapes, deterministic order, and empty arrays. |
| Task 6 | Integrity, rebuild, boundaries, exit scenario, docs | VAL-021..026, VAL-068..070 (9) | High | Whole-projection parity, performance headroom, security/static inventories, seam/frozen-map boundaries, final dogfood, and full spec sync. |

All six task specs have a Validation Contract Coverage section and concrete
command blocks. No required task command is a placeholder.

## 4. Required Evidence Strategy

The 70 rows below are one-to-one with the Ready contract. Each assertion has
its own evidence ID and outcome source; no split assertion inherits its
parent's result.

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | Schema introspection proving the exact DEC-076 columns, types, FK, CHECK text, and partial-index definitions | VAL-001 | QA Lead | Task 1 tests/report; QA-01 |
| EVID-02 | Pre/up/down/up schema and row snapshots with empty diffs | VAL-002 | QA Lead | Task 1 tests/report; QA-02 |
| EVID-03 | Sorted full-column tuples for `escalation-default` across all four workflow tables against the literal TRD §4.6 reference | VAL-003 | QA Lead | Task 1 tests/report; QA-03 |
| EVID-04 | Seeded loader result surfacing `linked_to_decision` as known-but-unimplemented without panic | VAL-004 | QA Lead | Task 1 tests/report; QA-04 |
| EVID-05 | Prior-version rebuild and current-version no-rebuild counterpart output | VAL-005 | QA Lead | Task 1 tests/report; QA-05 |
| EVID-06 | Golden ordered effect slice and repeated-build deep equality | VAL-006 | QA Lead | Task 2 tests/report; QA-06 |
| EVID-07 | Identifier-validation matrix covering invalid/duplicate slug, invalid/duplicate state, and accepted identifiers with distinct value-naming errors | VAL-007 | QA Lead | Task 2 tests/report; QA-07 |
| EVID-08 | Exact transition set for one cancelled state, including deduplication of an explicit cancel edge | VAL-008 | QA Lead | Task 2 tests/report; QA-08 |
| EVID-09 | Built define run proving exactly one current-schema `workflow.defined` event carrying every effect | VAL-009 | QA Lead | Task 3 tests/report; QA-09 |
| EVID-10 | Per-case journal size/hash before and after each rejected define plus one accepted append control | VAL-010 | QA Lead | Task 3 tests/report; QA-10 |
| EVID-11 | Exact permitted role sets for `workflow.define` and `workflow.assign` | VAL-011 | QA Lead | Task 3 tests/report; QA-11 |
| EVID-12 | Transition tuple equality for every documented arrow alias | VAL-012 | QA Lead | Task 3 tests/report; QA-12 |
| EVID-13 | Byte-exact stderr block and SQL-backed count comparison | VAL-013 | QA Lead | Task 3 tests/report; QA-13 |
| EVID-14 | First kind assignment's single event and exact one-row projected state | VAL-014 | QA Lead | Task 4 tests/report; QA-14 |
| EVID-15 | First project assignment's single event and exact NULL-kind one-row projected state | VAL-015 | QA Lead | Task 4 tests/report; QA-15 |
| EVID-16 | Six distinct value-naming assign rejection traces plus one compatible accepted control | VAL-016 | QA Lead | Task 4 tests/report; QA-16 |
| EVID-17 | Two assigned-workflow creation fixtures proving both DEC-077 initial-state branches | VAL-017 | QA Lead | Task 4 tests/report; QA-17 |
| EVID-18 | Exact JSON key set, SQL-backed state/transition counts, and `any` as null | VAL-018 | QA Lead | Task 5 tests/report; QA-18 |
| EVID-19 | Exact qualifying and excluded workflow IDs for the `--kind` filter | VAL-019 | QA Lead | Task 5 tests/report; QA-19 |
| EVID-20 | Exact kind-scope `assigned_to` object shape | VAL-020 | QA Lead | Task 5 tests/report; QA-20 |
| EVID-21 | Every-table normalized snapshots, schema, FK output, and journal hashes at three rebuild boundaries | VAL-021 | QA Lead | Task 6 integration evidence; QA-21 |
| EVID-22 | Healthy and incompatible maintenance outputs plus `--warn-as-error` exit | VAL-022 | QA Lead | Task 6 tests/CLI trace; QA-22 |
| EVID-23 | Built-CLI rejection traces for both reserved rules with version-scoped stderr and unchanged journal bytes | VAL-023 | QA Lead | Task 6 report/final review; QA-23 |
| EVID-24 | Custom-state rejection/no-append and v0.1-state success under a custom workflow | VAL-024 | QA Lead | Task 6 integration evidence; QA-24 |
| EVID-25 | Scripted real-workspace trajectory with outputs and unaffected control issue | VAL-025 | QA Lead | Task 6 integration evidence; QA-25 |
| EVID-26 | Nine single-judgment Human outcomes from a scratch real workspace | VAL-026 | Human | Phase-end human checklist |
| EVID-27 | Duplicate-kind rejection, two distinct-kind acceptances, and resulting rows | VAL-027 | QA Lead | Task 1 tests/report; QA-32 |
| EVID-28 | First project-row acceptance, second project-row rejection, and final one-row state | VAL-028 | QA Lead | Task 1 tests/report; QA-33 |
| EVID-29 | Four scope/kind INSERT cases proving both rejected and both accepted shapes | VAL-029 | QA Lead | Task 1 tests/report; QA-34 |
| EVID-30 | Sorted full-column tuples for all three Phase 18 defaults before and after `0010`, deep-equal to the immutable Phase 18 reference | VAL-030 | QA Lead | Task 1 tests/report; QA-35 |
| EVID-31 | Before/after journal, schema, and workflow-table snapshots around load plus non-test evaluator inventory | VAL-031 | QA Lead | Task 1 tests/report; QA-36 |
| EVID-32 | States/categories rejection matrix with empty, length-mismatch, unknown-category, and accepted controls | VAL-032 | QA Lead | Task 2 tests/report; QA-37 |
| EVID-33 | Unknown-edge and duplicate-transition rejection outputs naming each edge | VAL-033 | QA Lead | Task 2 tests/report; QA-38 |
| EVID-34 | Effect rows proving all four implemented rule names are accepted on their declared edges with params | VAL-034 | QA Lead | Task 2 tests/report; QA-39 |
| EVID-35 | Four distinct reserved/unknown rule rejection messages with v0.2b/v0.2c/v0.2d tokens | VAL-035 | QA Lead | Task 2 tests/report; QA-40 |
| EVID-36 | Invalid JSON and empty-roles rejection outputs plus verbatim valid params in the effect | VAL-036 | QA Lead | Task 2 tests/report; QA-41 |
| EVID-37 | Unknown rule-edge rejection and declared-edge acceptance traces | VAL-037 | QA Lead | Task 2 tests/report; QA-42 |
| EVID-38 | Non-issue zero-cancel acceptance/no-edge state and issue-applicable rejection with exact DEC-090 error | VAL-038 | QA Lead | Task 2 tests/report; QA-43 |
| EVID-39 | Two-cancelled-state rejection naming both states and producing no effect slice | VAL-039 | QA Lead | Task 2 tests/report; QA-44 |
| EVID-40 | Initial-state helper output for first-`unstarted` and no-`unstarted` branches | VAL-040 | QA Lead | Task 2 tests/report; QA-45 |
| EVID-41 | Five distinct assignment-validator rejections and two compatible accepted controls | VAL-041 | QA Lead | Task 2 tests/report; QA-46 |
| EVID-42 | Core-suite before/after journal and row snapshots plus static no-mutation SQL inventory | VAL-042 | QA Lead | Task 2 tests/report; QA-47 |
| EVID-43 | Sorted projected tuples from all four workflow tables deep-equal to the definition | VAL-043 | QA Lead | Task 3 tests/report; QA-48 |
| EVID-44 | `workflow.Load` result for the newly defined workflow with complete states, edges, and rules | VAL-044 | QA Lead | Task 3 tests/report; QA-49 |
| EVID-45 | Before/after `activity_log` row sets around successful define | VAL-045 | QA Lead | Task 3 tests/report; QA-50 |
| EVID-46 | Four workflow-table row snapshots unchanged after every rejected define | VAL-046 | QA Lead | Task 3 tests/report; QA-51 |
| EVID-47 | Non-zero exits and pairwise-distinct value-naming stderr for every define rejection | VAL-047 | QA Lead | Task 3 tests/report; QA-52 |
| EVID-48 | Decoded JSON/TOON/Markdown/YAML `policy show` outputs containing both new rows | VAL-048 | QA Lead | Task 3 tests/report; QA-53 |
| EVID-49 | Denied-role define error/no-append trace and successful `tl` control | VAL-049 | QA Lead | Task 3 tests/report; QA-54 |
| EVID-50 | Six accepted applicability projections, explicit `any IS NULL`, and invalid-value rejection/no-append trace | VAL-050 | QA Lead | Task 3 tests/report; QA-55 |
| EVID-51 | Reassignment row state plus continued `workflows` and list presence for the superseded workflow | VAL-051 | QA Lead | Task 4 tests/report; QA-56 |
| EVID-52 | Kind first-assign insert event, ordered delete+insert reassignment event, and replayed one-row state | VAL-052 | QA Lead | Task 4 tests/report; QA-57 |
| EVID-53 | Project first-assign insert event, ordered delete+insert reassignment event, and replayed one-row state | VAL-053 | QA Lead | Task 4 tests/report; QA-58 |
| EVID-54 | Creation trace showing an uncovered kind resolves to the project assignment | VAL-054 | QA Lead | Task 4 tests/report; QA-59 |
| EVID-55 | Creation trace showing kind assignment wins when project and kind rows coexist | VAL-055 | QA Lead | Task 4 tests/report; QA-60 |
| EVID-56 | Per-rejection journal hashes and assignment-row snapshots proving no append or projection mutation | VAL-056 | QA Lead | Task 4 tests/report; QA-61 |
| EVID-57 | Before/after `activity_log` row sets around successful assign | VAL-057 | QA Lead | Task 4 tests/report; QA-62 |
| EVID-58 | Three no-assignment creations proving built-in workflow IDs and `todo` status | VAL-058 | QA Lead | Task 4 tests/report; QA-63 |
| EVID-59 | Verified Task-3 commit, baseline binary hash, committed fixed fixture/captures/SHA-256 manifest, empty normalization allow-list, and raw-byte candidate diffs | VAL-059 | QA Lead | Task 4 tests/report and Review 1 baseline gate; QA-64 |
| EVID-60 | Four workflow-list formats decoded to one equivalent structure | VAL-060 | QA Lead | Task 5 tests/report; QA-65 |
| EVID-61 | Byte-exact markdown test against identified committed expected bytes | VAL-061 | QA Lead | Task 5 tests/report; QA-66 |
| EVID-62 | Filtered/unfiltered exact ID sets proving inactive exclusion and inclusion | VAL-062 | QA Lead | Task 5 tests/report; QA-67 |
| EVID-63 | Combined-filter result set deep-equal to the exact qualifying intersection | VAL-063 | QA Lead | Task 5 tests/report; QA-68 |
| EVID-64 | Project-scope `assigned_to` key-presence proof with `kind` absent | VAL-064 | QA Lead | Task 5 tests/report; QA-69 |
| EVID-65 | Unassigned workflow JSON proof showing `assigned_to: []`, not null or missing | VAL-065 | QA Lead | Task 5 tests/report; QA-70 |
| EVID-66 | Ordered results/assignments plus byte-identical consecutive outputs | VAL-066 | QA Lead | Task 5 tests/report; QA-71 |
| EVID-67 | Four no-match list runs proving empty arrays and exit 0 | VAL-067 | QA Lead | Task 5 tests/report; QA-72 |
| EVID-68 | Ordered SQL query returning zero seeded rows for both reserved rule names | VAL-068 | QA Lead | Task 6 report/final review; QA-73 |
| EVID-69 | Exact non-test production `rg` inventory returning no evaluator/query matches | VAL-069 | QA Lead | Task 6 report/final review; QA-74 |
| EVID-70 | Exact non-test `FrozenDefinitions(` inventory returning no production matches | VAL-070 | QA Lead | Task 6 report/final review; QA-75 |

Every Developer completion report must cite the specific test name, file path
or output proving each covered assertion; restating the assertion is not
evidence. Every required command must be recorded with its verbatim terminal
token. A non-terminating command is `BLOCKED`, never `PASS`.

## 5. Deterministic Gates And Applicability

### Required Command Matrix

`Y` means the exact command is required by that task spec. The Phase End
column is the QA Lead's fresh heavy-mode gate; no result carries forward from
Phase 18 or from this planning-time feasibility run.

| Exact command | T1 | T2 | T3 | T4 | T5 | T6 | Phase End |
|---|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| `gofmt -l cmd internal` | Y | Y | Y | Y | Y | Y | Y |
| `golangci-lint run ./...` | Y | Y | Y | Y | Y | Y | Y |
| `go test -race ./...` | - | - | - | - | - | - | Y |
| `go test -race -coverprofile=coverage.out ./...` | Y | Y | Y | Y | Y | Y | Y |
| `go tool cover -func=coverage.out` | Y | Y | Y | Y | Y | Y | Y |
| `go test -tags integration -race ./...` | Y | - | Y | Y | - | Y | Y |
| `go test -bench . -benchmem ./...` | - | - | - | - | - | Y | Y |
| `go build -o /tmp/fid ./cmd/fid` | Y | - | Y | Y | Y | Y | Y |

For each task, a missing command, narrower substitute, silently skipped
command, or command reported without a terminal token is Blocking unless the
completion report is a committed blocker report that records it `BLOCKED`
with a reason.

### Expected Results

| Gate | Expected result |
|---|---|
| Format | No file paths; terminating exit 0. |
| Lint/static analysis | `0 issues.` and terminating exit 0. |
| Unit/race | Every package terminates successfully; no race report. |
| Coverage/confidence | Every package terminates successfully and writes a fresh `coverage.out`. |
| Coverage report | `cmd/fid` remains at least 70%; `internal/issues`, `internal/claims`, `internal/workflow`, `internal/policy`, and `internal/journal` each remain at least 80%; terminal total line recorded. |
| Integration | Integration-tagged packages terminate successfully; no skipped Phase 19 integration path. |
| Benchmarks | Command terminates; 10k-issue read `<500 ms`; 100k-event rebuild `<30 s`; margin compared with Phase 18 `49.67 ms` / `27.661 s` and the planning baseline in Section 6. |
| Build | `/tmp/fid` is created by a terminating exit-0 command. |

### Command Applicability

- Standalone end-to-end command: **N/A**. `AGENTS.md` records no accepted
  standalone E2E/held-out harness. VAL-024/025 must be integration-tagged Go
  tests under `cmd/fid`, and the exact
  `go test -tags integration -race ./...` command is required.
- `go generate ./...`: **N/A**. No task names generated inputs.
- `govulncheck ./...`, `go-licenses check ./...`, and `gitleaks detect`:
  **N/A for the Phase 19 gate**, because this is not an external-release or
  release-readiness phase. Security review still applies under
  `08_QUALITY_PLAN.md` §5.2. Any new dependency would independently trigger
  the dependency-budget hard stop.
- Universal `--dry-run`: **N/A until v0.2e** per
  `10_IMPLEMENTATION_ROADMAP.md`; see Section 2.
- Held-out scenario command: **N/A** because no accepted harness location or
  command exists. VAL-026 remains a normal Human dogfood check, not a
  represented-as-held-out scenario.

Assertion applicability is complete across the 70-row contract:

| Assertion window | Applicable command/evidence surface |
|---|---|
| VAL-001..005, VAL-027..031 | Task 1 coverage gate plus integration and build; migration/schema SQL, tuple snapshots, loader state, static inventory, and stale/fresh projection tests |
| VAL-006..008, VAL-032..042 | Task 2 coverage gate; pure-core unit matrices, golden effects, state snapshots, and static mutation inventory |
| VAL-009..013, VAL-043..050 | Task 3 coverage, integration, and build; built-CLI, journal, four-table, permissions, help/output, and `activity_log` evidence |
| VAL-014..017, VAL-051..059 | Task 4 coverage, integration, and build; assignment events/replay, issue creation, rejection atomicity, and committed v0.1 parity captures |
| VAL-018..020, VAL-060..067 | Task 5 coverage and build; four-format built-CLI outputs, SQL-backed counts, key-presence checks, filters, ordering, and committed markdown bytes |
| VAL-021..025, VAL-068..070 | Task 6 and fresh phase-end eight-command gate; rebuild/maintenance/E2E traces, exact SQL, and exact static inventories |
| VAL-026 | Nine Human-only rows in Section 9 after deterministic readiness is established |

## 6. Command And Feasibility Validation

Heavy-mode validation was performed fresh on 2026-07-26. The exact commands
were executed on the pre-implementation baseline to prove syntax, tool/cache
configuration, package scope, termination, and fixture applicability. These
are feasibility results only; they do not satisfy any task or phase-end gate
after product changes.

### Exact Command Validation

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | Exact execution | PASS | Terminated; `exit 0 (no output)`. `gofmt -h` confirms `-l` and path arguments. |
| `golangci-lint run ./...` | Exact execution | PASS | Terminated with `0 issues.`; `golangci-lint run --help` confirms command/flags. Installed version: 2.12.2 built with Go 1.26.2. |
| `go test -race ./...` | Exact execution | PASS | Terminated; terminal token `ok github.com/mbell-nzl/fid/internal/workflow (cached)`. |
| `go test -race -coverprofile=coverage.out ./...` | Exact execution | PASS | Terminated; terminal token `ok github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 83.2% of statements`. Current `cmd/fid` is 77.7%, `internal/issues` 81.5%, and `internal/journal` 81.0%. |
| `go tool cover -func=coverage.out` | Exact execution after the coverage command | PASS | Terminated with `total: (statements) 76.8%`. `go tool cover -help` confirms `-func`; `coverage.out` is gitignored. |
| `go test -tags integration -race ./...` | Exact execution | PASS | Terminated; terminal token `ok github.com/mbell-nzl/fid/internal/workflow (cached)`. Existing integration build tags were found under `cmd/fid` and `internal/issues`. |
| `go test -bench . -benchmem ./...` | Exact execution | PASS | Terminated; terminal token `ok github.com/mbell-nzl/fid/internal/workflow 0.043s`. `BenchmarkList10kIssuesWithRegistry`: `44.03 ms/registry-read`; `BenchmarkRebuildWorkspace100kEvents`: `24674 ms/rebuild`. |
| `go build -o /tmp/fid ./cmd/fid` | Exact execution | PASS | Terminated; `exit 0 (no output)`; `/tmp/fid` then supplied the help evidence below. |

Tool/cache feasibility:

- Go is `go1.26.2 linux/amd64`; race mode is supported for linux/amd64.
- `.codex/config.toml` points Go, module, GOPATH, lint, and XDG caches at
  gitignored workspace paths. All exact commands ran without escalation.
- `cmd/`, `internal/`, `go.mod`, and `cmd/fid/main.go` exist.
- `coverage.out`, `.gocache`, `.gomodcache`, `.gopath`,
  `.golangci-cache`, and `.cache` are ignored.

### Test Scope, Fixture, Artifact, And Review-Step Validation

| Referenced step or assumption | Validation method | Result | Evidence |
|---|---|---|---|
| Integration scope | File-state inventory | PASS | Integration-tagged tests exist in `cmd/fid/e2e_v01*_test.go`, `cmd/fid/sql_integration_test.go`, and `internal/issues/transition_rules_test.go`; exact integration command passed. |
| Performance controls | Test inventory plus exact benchmark run | PASS | `internal/issues/registry_read_bench_test.go` contains `BenchmarkList10kIssuesWithRegistry`; `internal/journal/rebuild_bench_test.go` contains 1k/10k/100k rebuild benchmarks. Both required controls produced measurements. |
| Whole-projection rebuild fixture | Existing test-pattern inspection | PASS | `internal/journal/rebuild_scale_test.go` already performs before/after/second rebuild snapshots; `cmd/fid/e2e_v01_closure_helpers_test.go` has every-table, `sqlite_schema`, and `PRAGMA foreign_key_check` helpers. Task 6 extends this proven pattern. |
| Task 1 planned test/migration paths | Parent/file-state checks and task spec | PASS | `migrations/`, `internal/db/`, and `internal/workflow/` exist; Task 1 explicitly authorizes the new files. Absence of migrations `0009`/`0010` before Task 1 is expected. |
| Tasks 2-5 planned test paths | Parent/file-state checks, task specs, and DEC-089 | PASS | `cmd/fid`, `internal/workflow`, `internal/issues`, `internal/policy`, and `internal/format` exist. `02_TRD.md` §3 now assigns authoring, assignment, and resolution to `internal/workflow` with the required must-nots. |
| Task 6 planned test paths | Parent/file-state checks and task spec | PASS | `internal/maintenance`, `internal/journal`, and `cmd/fid` exist; prior maintenance and E2E patterns exist. |
| VAL-003/030 tuple references | Immutable commit verification, fresh migrated fixture, exact SQL capture | PASS | Phase 18 closure commit `097e45cc5111be5ddcd5bf232645659045571521` resolves and contains migration `0007` plus `definitions_fixture.go`. Four ordered queries on a fresh migrated workspace returned full default tuples. Candidate comparisons are sorted deep equality, never counts. |
| VAL-059 pre-change output reference | Contract/task inspection, `git worktree add -h`, detached-worktree build, `go -C` dry syntax check, SHA-256 tool help | PASS with mandatory Review-1/Task-4 capture condition | A detached worktree at immutable commit `097e45cc5111be5ddcd5bf232645659045571521` built successfully as a procedure surrogate; binary SHA-256 was `db94dad55a89815ccd1ac0077612e8567883024abf85249ce119ee17cbe3ddfc`. `go -C /home/mark/projects/fid version` terminated with `go version go1.26.2 linux/amd64`. The actual source is necessarily the later Task 3 implementation commit and must be captured before `create.go` changes using the exact procedure below. |
| VAL-018 markdown golden | Task/path inspection | PASS with phase-end precondition | Task 5 authorizes `internal/format/markdown_workflow_test.go`; expected bytes may be a committed `testdata` fixture or an explicit literal. The completion report must identify which and the final review must verify it is committed. |
| Static seam inventory | `rg` availability and exact current inventory below | PASS | The exact command terminated with no output / exit 1. Reserved declarations and authoring rejections inside `internal/workflow` are expected and reviewed separately. |
| Frozen-definition production-call inventory | `rg` availability and exact current inventory below | PASS | The exact command terminated with no output / exit 1. |
| Built-CLI support commands | Exact help execution | PASS | `/tmp/fid init --help`, `issue create --help`, `issue transition --help`, `policy show --help`, `maintenance check --help`, `rebuild --help`, and `sql --help` all terminated with documented flags. |
| Future `fid workflow` family | Current help plus canonical CLI/task specs | PASS with implementation precondition | Current `/tmp/fid workflow --help` exits 1 because the family does not exist yet, as expected. `06_DESIGN_PLAN_CLI.md` and Tasks 3-5 provide concrete define/assign/list syntax. Missing or divergent help after Task 5 is Blocking. |
| Review 1 representative command | `go help testflag`, package-pattern help, and directory checks | PASS | `go test -race ./internal/db/... ./internal/workflow/... ./internal/policy/... ./cmd/fid/...` is syntactically valid; all scopes exist. It is independent review evidence, never a substitute for Developer full gates. |
| Review 2 representative command | `go help testflag`, package-pattern help, and directory checks | PASS | `go test -race ./internal/issues/... ./internal/workflow/... ./internal/format/... ./cmd/fid/...` is syntactically valid; all scopes exist. It is independent review evidence, never a substitute for Developer full gates. |
| Completion reports and QA artifacts | File-state check | PASS as planned artifacts | `task-1-done.md` through `task-6-done.md`, `review-1.md`, `review-2.md`, `final-review.md`, and `human-checklist.md` do not yet exist. Their absence is expected now; each becomes a hard-stop requirement in its owning window. |
| Human scratch workspace | CLI help, design spec, and temp-directory policy | PASS with implementation precondition | `fid init` and supporting issue/transition reads are available; workflow syntax is concrete in the CLI spec. Use a fresh path returned by `mktemp -d /tmp/fid-phase19-human.XXXXXX`; do not reuse project data. |

The exact phase-end static inventory commands are:

```bash
rg -n --glob '*.go' --glob '!**/*_test.go' '\brequired_gates_current\b|\bhas_completion\b' cmd/fid internal/issues internal/dev internal/journal internal/db internal/maintenance internal/policy
rg -n --glob '*.go' --glob '!**/*_test.go' --glob '!internal/workflow/definitions_fixture.go' 'FrozenDefinitions\(' cmd internal
```

Both were executed during this fresh feasibility pass and returned no output /
exit 1. Any phase-end match fails VAL-069 or VAL-070 respectively; the
expected no-match exit 1 is not a command failure for these two inventory
assertions.

#### VAL-003 And VAL-030 Frozen Tuple Procedure

The frozen full-tuple source for the three Phase 18 defaults is commit
`097e45cc5111be5ddcd5bf232645659045571521`
(`docs(phase): wrap phase 18 and reconcile workflow specs`):

- `migrations/0007_v02a_workflow_tables.up.sql` supplies every persisted
  workflow, state/ordinal, explicit transition (including cancel edges), rule,
  and NULL `params_json` tuple.
- `internal/workflow/definitions_fixture.go` independently supplies the frozen
  v0.1 semantic graph. The Task 1 test must compare the loaded semantic graph
  to this fixture and compare the persisted tuples to migration `0007`; one
  reference cannot stand in for the other.
- QA verified the commit with
  `git rev-parse --verify 097e45cc5111be5ddcd5bf232645659045571521^{commit}`.
  A fresh pre-Phase-19 migrated workspace returned all expected full rows from
  the four queries below, including NULLs and ordinals.

Capture each table as ordered JSON before `0010`, after `0010`, and from the
immutable worktree, then deep-equal the bytes:

```text
SELECT id,name,applies_to,active,COALESCE(created_at,'<NULL>') AS created_at
FROM workflows
WHERE id IN ('epic-default','task-default','subtask-default')
ORDER BY id;

SELECT workflow_id,state,category,ordinal
FROM workflow_states
WHERE workflow_id IN ('epic-default','task-default','subtask-default')
ORDER BY workflow_id,ordinal,state;

SELECT workflow_id,from_state,to_state
FROM workflow_transitions
WHERE workflow_id IN ('epic-default','task-default','subtask-default')
ORDER BY workflow_id,from_state,to_state;

SELECT workflow_id,from_state,to_state,rule,
       COALESCE(params_json,'<NULL>') AS params_json
FROM workflow_transition_requires
WHERE workflow_id IN ('epic-default','task-default','subtask-default')
ORDER BY workflow_id,from_state,to_state,rule,params_json;
```

The Task 1 migration harness must write those four query results under each of
these exact capture roots:

```text
/tmp/fid-phase19-val030/immutable-phase18/
/tmp/fid-phase19-val030/before-0010/
/tmp/fid-phase19-val030/after-0010/
```

Each root contains `workflows.json`, `workflow-states.json`,
`workflow-transitions.json`, and `workflow-transition-requires.json`. The
immutable capture comes from a clean workspace initialized by a binary built
from the verified Phase 18 commit; the other two come from the Task 1
migration harness immediately before and after applying `0010`. The harness
uses the shown column order, row order, JSON encoding, and NULL sentinel for
all three captures, then records `sha256sum` lines and byte diffs for every
same-named file. This procedure was feasibility-checked with the built
binary's `fid sql` command and `--format json`, using each of the four complete
SQL blocks above as its positional query argument, in the fresh
`/tmp/fid-phase19-val003-feasibility-qa` workspace; every query terminated and
returned complete tuples.

VAL-003 uses the same column/order discipline against this literal reference,
with no count substitute:

```text
workflows:
  ('escalation-default','Escalation default','escalation',1,'<NULL>')
workflow_states:
  ('escalation-default','open','started',1)
  ('escalation-default','resolved','completed',2)
workflow_transitions:
  ('escalation-default','open','resolved')
workflow_transition_requires:
  ('escalation-default','open','resolved','linked_to_decision','<NULL>')
```

Any missing, extra, differently ordered after sorting, or differing tuple
fails VAL-003 or VAL-030. Row counts alone cannot pass either assertion.

#### VAL-059 Immutable Baseline Procedure

The actual baseline commit is late-bound by design: it is the Task 3
**implementation** commit copied from the committed `task-3-done.md`, not the
following report commit and not ambient `HEAD`. Review 1 sets the task-specific
shell variable `PHASE19_TASK3_IMPLEMENTATION_COMMIT` to the exact copied
40-hex value, runs
`git rev-parse --verify "$PHASE19_TASK3_IMPLEMENTATION_COMMIT^{commit}"`, and
records the exact 40-hex result before Task 4 starts.

Before editing `internal/issues/create.go`, Task 4 must:

1. Create a detached worktree at that verified hash at
   `/tmp/fid-phase19-v01-baseline`, build its `./cmd/fid` as
   `/tmp/fid-phase19-v01-baseline.bin`, and record the binary's
   `sha256sum` line in `task-4-done.md`. The exact commands are:

   ```bash
   git worktree add --detach /tmp/fid-phase19-v01-baseline "$PHASE19_TASK3_IMPLEMENTATION_COMMIT"
   go -C /tmp/fid-phase19-v01-baseline build -o /tmp/fid-phase19-v01-baseline.bin ./cmd/fid
   sha256sum /tmp/fid-phase19-v01-baseline.bin
   ```

2. Commit `cmd/fid/testdata/v01_parity/capture.sh` and a fixed fixture rooted
   at `cmd/fid/testdata/v01_parity/fixture/`. The fixture contains no
   `workflow_assignments` event/row and fixes IDs (`FID-a10001`,
   `FID-b20002`, `FID-c30003`), titles, epic/task/subtask kinds, parent links,
   and timestamps in its journal.
3. Run the same committed capture script first with the baseline binary and
   later with the Task 4 candidate. Baseline outputs are exactly:
   `issue-get-epic.json`, `issue-get-task.json`,
   `issue-get-subtask.json`, `issue-list-epic.json`,
   `issue-list-task.json`, and `issue-list-subtask.json`, all directly under
   `cmd/fid/testdata/v01_parity/`. Candidate outputs go under
   `/tmp/fid-phase19-v01-parity-candidate/`.
   `capture.sh` accepts exactly two positional arguments — binary path and
   output root — recreates the workspace from the committed fixture for each
   run, and executes these six command shapes:

   ```text
   fid issue get FID-a10001 --format json
   fid issue get FID-b20002 --format json
   fid issue get FID-c30003 --format json
   fid issue list --kind epic --format json
   fid issue list --kind task --format json
   fid issue list --kind subtask --format json
   ```

   The baseline and candidate captures are invoked exactly as:

   ```bash
   cmd/fid/testdata/v01_parity/capture.sh /tmp/fid-phase19-v01-baseline.bin cmd/fid/testdata/v01_parity
   cmd/fid/testdata/v01_parity/capture.sh /tmp/fid /tmp/fid-phase19-v01-parity-candidate
   ```

4. Commit `cmd/fid/testdata/v01_parity/normalization.json` with exactly
   `{"json_pointers":[],"replacement":null}`. Normalization is therefore an
   explicit no-op and all six files are raw-byte compared. If the fixed
   fixture cannot keep timestamps stable, Task 4 is `BLOCKED`; the allow-list
   is not silently widened.
5. Commit `cmd/fid/testdata/v01_parity/SHA256SUMS` containing a `sha256sum`
   line for each of the six baseline outputs. Every manifest line must contain
   the actual 64-character lowercase SHA-256, exactly two spaces, and the bare
   relative output filename; placeholders are forbidden. `task-4-done.md`
   copies the six lines verbatim. After Task 4 is committed, the immutable
   baseline-artifact identity is the reported Task 4 implementation commit plus
   those six manifest hashes.

   **Corrected 2026-08-02 by DEC-103.** This step previously directed Review 2
   and final QA to run `sha256sum -c cmd/fid/testdata/v01_parity/SHA256SUMS`
   from the repository root. **That command never worked** — `sha256sum -c`
   resolves listed filenames against the working directory, and the manifest
   correctly holds bare names so it can be checked against any capture
   directory. Review 2 executed it as written and it exited 1 with
   `sha256sum: WARNING: 6 listed files could not be read` (ESC-19-007). The
   authoritative procedure is:

   ```bash
   manifest="$(git rev-parse --show-toplevel)/cmd/fid/testdata/v01_parity/SHA256SUMS"
   capture_root="$(mktemp -d)"
   go build -o "$capture_root/fid" ./cmd/fid
   cmd/fid/testdata/v01_parity/capture.sh "$capture_root/fid" "$capture_root"
   ( cd "$capture_root" && sha256sum -c "$manifest" )
   ```

   All six entries must report `OK`. Running the manifest from
   `cmd/fid/testdata/v01_parity/` is **also prohibited**: it checks the
   committed baselines against their own digests and proves nothing about the
   candidate. The baseline is established and closed (DEC-101) and is never
   regenerated. **CF-025 retires this whole procedure** — once Task 6 lands the
   Go golden test, the shell ceremony stays committed as provenance only and no
   further manual run is authorized.

The detached-worktree feasibility run terminated successfully. `git worktree
add -h` confirms `--detach <path> <commit-ish>`; `sha256sum --help` confirms
manifest creation/checking. The temporary worktree was removed after the
proof. No Phase 18 result is carried as Phase 19 assertion evidence.

Invalid or infeasible items are not carried forward in heavy mode. This fresh
run found no open feasibility blocker and no placeholder command.

## 7. Deterministic Checklist

Each of the 69 QA Lead-validated assertions has exactly one step below;
VAL-026 has the Human-only verification steps in Section 9. Every step has one
validator role and produces an independent outcome.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Verify the exact DEC-076 assignment-table structure `(qa-lead)` | VAL-001 | Migration SQL plus table, FK, CHECK, and partial-index introspection output. |
| QA-02 | Verify migration `0009`/`0010` up-down-up reversibility `(qa-lead)` | VAL-002 | Three state snapshots and empty deterministic diffs. |
| QA-03 | Verify `escalation-default` seed tuple conformance `(qa-lead)` | VAL-003 | Full sorted tuples for all four tables against the literal Section 6 reference; no counts-only proof. |
| QA-04 | Verify `linked_to_decision` loads structurally as unimplemented `(qa-lead)` | VAL-004 | Loader output identifying the rule and unimplemented flag, with no panic. |
| QA-05 | Verify stale projection rebuild and fresh projection no-rebuild branches `(qa-lead)` | VAL-005 | Both branch outputs and current schema-version state. |
| QA-06 | Verify authoring effect order and repeated-build determinism `(qa-lead)` | VAL-006 | Golden slice diff and repeated deep-equality output. |
| QA-07 | Verify identifier validation and accepted control `(qa-lead)` | VAL-007 | Named identifier subtests and distinct offending-value errors. |
| QA-08 | Verify one-cancelled-state edge materialization and explicit-edge deduplication `(qa-lead)` | VAL-008 | Exact expected transition set. |
| QA-09 | Verify successful define appends exactly one current-schema event carrying every effect `(qa-lead)` | VAL-009 | Built-CLI trace and decoded journal event. |
| QA-10 | Verify every define rejection leaves journal bytes unchanged `(qa-lead)` | VAL-010 | Per-case before/after size/hash plus accepted append control. |
| QA-11 | Verify exact define/assign permission sets `(qa-lead)` | VAL-011 | Matrix test output naming exactly `tl` and `human` for both actions. |
| QA-12 | Verify all documented transition arrow aliases normalize to identical tuples `(qa-lead)` | VAL-012 | `->`, `→`, `⟶`, and `-->` tuple-set comparison. |
| QA-13 | Verify byte-exact define confirmation and SQL-backed counts `(qa-lead)` | VAL-013 | Literal stderr block and count query output. |
| QA-14 | Verify first kind assignment event and projected row `(qa-lead)` | VAL-014 | Event count/type and exact one-row kind assignment state. |
| QA-15 | Verify first project assignment event and NULL-kind projected row `(qa-lead)` | VAL-015 | Event count/type and exact one-row project assignment state. |
| QA-16 | Verify every assign rejection message and one compatible success `(qa-lead)` | VAL-016 | Six non-zero, pairwise-distinct value-naming errors and accepted control. |
| QA-17 | Verify assignment-aware workflow and initial-state resolution `(qa-lead)` | VAL-017 | Two assigned fixtures proving both initial-state branches. |
| QA-18 | Verify exact workflow-list JSON shape and SQL-backed counts `(qa-lead)` | VAL-018 | Exact key set, SQL count output, and `any` null proof. |
| QA-19 | Verify `--kind` filter qualifying and excluded sets `(qa-lead)` | VAL-019 | Exact included/excluded IDs. |
| QA-20 | Verify kind-scope `assigned_to` shape `(qa-lead)` | VAL-020 | Exact JSON object and key-presence output. |
| QA-21 | Verify non-destructive deterministic whole-projection rebuild twice `(qa-lead)` | VAL-021 | Every-table snapshots, schema/FK output, and per-journal-file hashes at all boundaries. |
| QA-22 | Verify maintenance assignment audit on healthy and incompatible states `(qa-lead)` | VAL-022 | Healthy output, deterministic finding, and non-zero warn-as-error exit. |
| QA-23 | Verify both reserved gate-obligation rules remain CLI-rejected `(qa-lead)` | VAL-023 | Exact version-scoped errors and unchanged journal bytes. |
| QA-24 | Verify custom-state permission boundary rejects and v0.1-name path succeeds `(qa-lead)` | VAL-024 | Exact error, unchanged journal bytes, successful control transition/state. |
| QA-25 | Execute the scripted slice-2 exit journey and unaffected control `(qa-lead)` | VAL-025 | Integration output and each command/result/projected-state checkpoint. |
| QA-26 | Verify package coverage thresholds, not only aggregate coverage `(qa-lead)` | no VAL | Fresh `go tool cover` output for `cmd/fid`, issues, claims, workflow, policy, journal. |
| QA-27 | Verify 10k-read and 100k-rebuild controls and headroom `(qa-lead)` | no VAL | Fresh benchmark lines, targets, Phase 18 and planning-baseline comparison. |
| QA-28 | Review SQL/command injection, unsafe input, data-loss, and journal-atomicity risks `(qa-lead)` | no VAL | Parameterization tests/traces, rejection/no-append evidence, security findings. |
| QA-29 | Verify file/function budgets, exported docs, dependency budget, and architecture boundaries `(qa-lead)` | no VAL | Reproducible file-state counts, lint output, dependency diff, resolved boundary decision. |
| QA-30 | Verify required documentation and help are synchronized `(qa-lead)` | no VAL | Diffs/captures for PRD, TRD, schema, design command tree, CLI spec, glossary, roadmap wrap obligations. |
| QA-31 | Audit all six completion reports against exact task commands and committed hashes/subjects `(qa-lead)` | no VAL | Task-spec/report comparison, terminal tokens, `git rev-parse` and `git log` output. |
| QA-32 | Verify kind-scope uniqueness rejects duplicate kind and accepts distinct kinds `(qa-lead)` | VAL-027 | Duplicate INSERT error, two accepted INSERTs, exact resulting rows. |
| QA-33 | Verify project-scope uniqueness accepts first row and rejects second `(qa-lead)` | VAL-028 | Both INSERT outcomes and final one-row query. |
| QA-34 | Verify both malformed and both well-formed scope/kind shapes `(qa-lead)` | VAL-029 | Four-case INSERT output and accepted-row state. |
| QA-35 | Verify all three Phase 18 defaults remain frozen tuple-for-tuple `(qa-lead)` | VAL-030 | Before/after `0010` full sorted tuple bytes deep-equal to immutable Phase 18 capture. |
| QA-36 | Verify seeded load neither writes nor gains a `linked_to_decision` evaluator `(qa-lead)` | VAL-031 | Before/after journal/schema/table snapshots and exact non-test inventory. |
| QA-37 | Verify states/categories structural validation `(qa-lead)` | VAL-032 | Empty, mismatch, unknown-category rejection outputs with accepted control. |
| QA-38 | Verify unknown-state edge and duplicate-transition rejection `(qa-lead)` | VAL-033 | Distinct edge-naming errors. |
| QA-39 | Verify all four implemented rule names reach the effect set `(qa-lead)` | VAL-034 | Rule/edge/params effect tuples for all four names. |
| QA-40 | Verify reserved and unknown rule rejection vocabulary `(qa-lead)` | VAL-035 | Distinct messages with v0.2b, v0.2c, v0.2d, and offending unknown name. |
| QA-41 | Verify `--require` params rejection and verbatim accepted params `(qa-lead)` | VAL-036 | Invalid-JSON and empty-roles errors plus accepted effect value. |
| QA-42 | Verify rule-edge validation rejects unknown and accepts declared edge `(qa-lead)` | VAL-037 | Two-case outcome and edge-naming error. |
| QA-43 | Verify DEC-090 zero-cancel conditional in both applicability states `(qa-lead)` | VAL-038 | Non-issue accepted/no-edge output and issue-applicable exact rejection. |
| QA-44 | Verify multiple cancelled states reject before effect construction `(qa-lead)` | VAL-039 | Both offending state names and no effect slice. |
| QA-45 | Verify both initial-state helper branches `(qa-lead)` | VAL-040 | First-`unstarted` and first-declared outputs. |
| QA-46 | Verify assignment-validator rejection matrix and compatible controls `(qa-lead)` | VAL-041 | Five distinct errors and two accepted request shapes. |
| QA-47 | Verify the Task 2 core is write-free `(qa-lead)` | VAL-042 | Before/after journal/row snapshots and static mutating-SQL inventory. |
| QA-48 | Verify defined workflow projection tuple-for-tuple `(qa-lead)` | VAL-043 | Full sorted tuples from all four workflow tables against declaration. |
| QA-49 | Verify a newly defined workflow resolves completely through `workflow.Load` `(qa-lead)` | VAL-044 | Loaded states, transitions, and rules. |
| QA-50 | Verify successful define writes no `activity_log` row `(qa-lead)` | VAL-045 | Before/after row-set equality. |
| QA-51 | Verify define rejections leave all workflow tables unchanged `(qa-lead)` | VAL-046 | Per-case four-table snapshots and empty diffs. |
| QA-52 | Verify define rejection stderr is non-zero, distinct, and value-naming `(qa-lead)` | VAL-047 | Captured exits/messages and pairwise comparison. |
| QA-53 | Verify both permission rows appear in all four `policy show` formats `(qa-lead)` | VAL-048 | Four decoded outputs with both exact actions. |
| QA-54 | Verify denied-role define is no-append and `tl` define succeeds `(qa-lead)` | VAL-049 | Standard permission error, journal hashes, accepted trace. |
| QA-55 | Verify exact `--applies-to` vocabulary, projections, NULL mapping, and rejection `(qa-lead)` | VAL-050 | Six accepted SQL values, `IS NULL` result, invalid-value error/hash. |
| QA-56 | Verify kind reassignment replaces only the assignment and preserves the old workflow `(qa-lead)` | VAL-051 | Exact assignment row, workflow row, and list presence. |
| QA-57 | Verify first/reassigned kind effect shapes and replay `(qa-lead)` | VAL-052 | Decoded insert-only and delete+insert events plus replayed row. |
| QA-58 | Verify first/reassigned project effect shapes and replay `(qa-lead)` | VAL-053 | Decoded insert-only and delete+insert events plus replayed row. |
| QA-59 | Verify project assignment is the fallback for an uncovered kind `(qa-lead)` | VAL-054 | Created issue's exact `workflow_id`. |
| QA-60 | Verify kind assignment wins over a project assignment `(qa-lead)` | VAL-055 | Created issue's exact `workflow_id`. |
| QA-61 | Verify every assign rejection leaves journal and assignment rows unchanged `(qa-lead)` | VAL-056 | Per-case hashes, row snapshots, and empty diffs. |
| QA-62 | Verify successful assign writes no `activity_log` row `(qa-lead)` | VAL-057 | Before/after row-set equality. |
| QA-63 | Verify no-assignment creation parity for all three kinds `(qa-lead)` | VAL-058 | Three exact `(kind, workflow_id, status)` tuples. |
| QA-64 | Verify immutable no-assignment `issue get/list` output parity `(qa-lead)` | VAL-059 | Verified Task-3 commit, binary/artifact SHA-256, clean committed files, reproducible captures, and six empty raw-byte diffs. |
| QA-65 | Verify four workflow-list formats carry equivalent data `(qa-lead)` | VAL-060 | Four decoded structures and deep-equality output. |
| QA-66 | Verify workflow-list markdown bytes against a committed expectation `(qa-lead)` | VAL-061 | Golden test output and committed-file/literal proof. |
| QA-67 | Verify active filter in filtered and unfiltered states `(qa-lead)` | VAL-062 | Exact ID sets showing exclusion and inclusion. |
| QA-68 | Verify kind and active filters compose as exact intersection `(qa-lead)` | VAL-063 | Expected and actual intersection sets. |
| QA-69 | Verify project-scope `assigned_to` omits `kind` `(qa-lead)` | VAL-064 | JSON key-presence proof. |
| QA-70 | Verify unassigned `assigned_to` is an empty array `(qa-lead)` | VAL-065 | JSON type/value and key-presence proof. |
| QA-71 | Verify deterministic result/assignment ordering and repeat bytes `(qa-lead)` | VAL-066 | Ordered IDs/entries and consecutive byte diff. |
| QA-72 | Verify all-format empty results are arrays with exit 0 `(qa-lead)` | VAL-067 | Four stdout captures and exit codes. |
| QA-73 | Verify neither reserved seam rule is seeded `(qa-lead)` | VAL-068 | Ordered zero-row SQL output. |
| QA-74 | Verify no production evaluator/query exists for reserved seam rules `(qa-lead)` | VAL-069 | Exact recorded `rg` command with no output/exit 1 and classification. |
| QA-75 | Verify no production path calls `FrozenDefinitions` `(qa-lead)` | VAL-070 | Exact recorded `rg` command with no output/exit 1. |

## 8. Mid-Phase Review Cadence

Two planned QA Reviewer windows apply:

1. **Review 1: Tasks 1-3**, written to
   `project-docs/docs/phases/phase-19/qa/review-1.md` after all three
   committed completion reports exist. It must complete with `VERDICT:
   CONTINUE` **before any Task 4 work starts**. It reviews all 37 Task 1-3
   assertions (VAL-001..013 and VAL-027..050), every Task 1-3 command result
   and terminal token, migration/event integrity, the resolved architecture
   decision, and Task 4 baseline readiness. It copies and verifies the Task 3
   implementation commit that VAL-059 will use; Task 4 must capture the
   baseline before editing `internal/issues/create.go`. Independent
   representative command:
   `go test -race ./internal/db/... ./internal/workflow/... ./internal/policy/... ./cmd/fid/...`.
2. **Review 2: Tasks 4-5**, written to
   `project-docs/docs/phases/phase-19/qa/review-2.md` after both committed
   completion reports exist and before Task 6 starts. It reviews all 24
   Task 4-5 assertions (VAL-014..020 and VAL-051..067), every Task 4-5 command
   result, the committed immutable v0.1 parity artifacts and SHA-256 check,
   assignment event shape/precedence, and four-format list conformance.
   Independent representative command:
   `go test -race ./internal/issues/... ./internal/workflow/... ./internal/format/... ./cmd/fid/...`.

Task 6 is the closure window and is reviewed by the QA Lead in the final
review: nine assertions (VAL-021..026 and VAL-068..070), not by a redundant
post-Task-6 mid-phase review. Any Blocking finding stops the next window. Each
review uses
`project-docs/templates/mid-phase-review-template.md`, verifies every in-window
assertion, compares the completion report with the task spec's exact commands,
and records one reproducible independent command. Targeted reviewer commands
do not replace the Developer's broad required commands.

## 9. Human Verification Scope

VAL-026 is the only Human-validated assertion. The phase-end
`human-checklist.md` must contain only these rows, one judgment per row. The
Human uses a scratch `.fid/` workspace and the shipped help/CLI spec, not
source code.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-26a | Judge whether `fid workflow define --help` and the CLI spec are sufficient to define one valid custom workflow without source-reading `(human)` | VAL-026 | Outcome and friction note. |
| HUMAN-26b | Judge whether `fid workflow assign --help` makes kind assignment discoverable and understandable `(human)` | VAL-026 | Outcome and friction note. |
| HUMAN-26c | Judge whether JSON workflow-list output is coherent for operator inspection `(human)` | VAL-026 | Outcome and optional captured excerpt. |
| HUMAN-26d | Judge whether TOON workflow-list output is coherent for agent-oriented inspection `(human)` | VAL-026 | Outcome and optional captured excerpt. |
| HUMAN-26e | Judge whether Markdown workflow-list output is scannable for a human `(human)` | VAL-026 | Outcome and optional captured excerpt. |
| HUMAN-26f | Judge whether YAML workflow-list output is coherent for configuration-oriented inspection `(human)` | VAL-026 | Outcome and optional captured excerpt. |
| HUMAN-26g | Trigger one invalid definition and judge whether the error supports self-correction `(human)` | VAL-026 | Outcome and exact error noted. |
| HUMAN-26h | Trigger one invalid assignment and judge whether the error supports self-correction `(human)` | VAL-026 | Outcome and exact error noted. |
| HUMAN-26i | Encounter the DEC-079 custom-state transition limitation and judge whether the shipped message/documentation makes the temporary boundary understandable rather than surprising `(human)` | VAL-026 | Outcome, exact error, and comprehension note. |

The checklist is feasible only after Tasks 3-6 provide the command family,
documentation, and DEC-079 note. A missing help surface, unreadable format, or
undocumented limitation is a failed Human row, not a reason to collapse rows
or substitute QA Lead judgment.

## 10. Phase-End Final Review Requirements And Decision Rule

The QA Lead creates and commits
`project-docs/docs/phases/phase-19/qa/final-review.md` and
`human-checklist.md` from their templates. The final review must:

- verify all six completion reports are committed, clean, and match their task
  specs' exact commands;
- copy commit hashes and subjects from `git` output and verify each reported
  hash resolves to a commit with `git rev-parse`;
- run the eight exact Phase End commands in Section 5 after Task 6, record a
  terminal token for each, and never treat this planning baseline as phase
  evidence;
- execute QA-01..QA-75 and cite reproducible output/file-state evidence for
  every result;
- record a final `PASS`, `FAIL`, `BLOCKED`, or `CARRY-FORWARD` for every
  VAL-001..070 assertion with its own cited evidence; no assertion inherits
  another row's result and no evidence-free pass is valid;
- assign **no heavy-mode `CARRY-FORWARD` outcome**. An unmet assertion,
  unavailable
  required fixture, non-terminating command, or unresolved defect blocks the
  phase and returns to Tech Lead/Human re-planning; record it `BLOCKED` until
  a new Human/Tech Lead disposition establishes a different phase plan;
- review requirement compliance, tests, journal/schema compatibility,
  architecture and layer boundaries, security/data-loss/input handling,
  performance, dependencies, structure budgets, and documentation accuracy;
- verify Task 6 records the 10k-read and 100k-rebuild measurements against
  both controls and prior baselines without weakening either target;
- reproduce VAL-003/030 from the immutable Phase 18 commit and full sorted
  tuples, and reproduce VAL-059 from the verified Task 3 commit, committed
  fixed fixture, six captures, `SHA256SUMS`, explicit no-op normalization,
  and six raw-byte candidate comparisons;
- verify full canonical spec-sync disposition at wrap, including
  `05_DESIGN_PLAN.md` §3.1, `02_TRD.md` §3/§4.6,
  `03_TRD_SCHEMA.md`, `06_DESIGN_PLAN_CLI.md`, `01_PRD.md`, `GLOSSARY.md`,
  and the roadmap closure note;
- generate the nine-row Human checklist above and keep all deterministic
  checks out of it; and
- record VAL-026 `PASS` only after the Human records PASS for all nine
  sub-checks.

The phase may proceed to final review only after Task 6 has a committed,
clean completion report and both mid-phase reviews are `CONTINUE` with every
Blocking finding fixed. The phase verdict may be PASS only when:

- this plan was unblocked before implementation and the refined contract was
  `Ready`;
- all exact task and phase commands terminated and passed;
- all 70 assertions are PASS with distinct cited evidence (69 QA Lead,
  VAL-026 Human);
- every `(qa-lead)` item is PASS with reproducible evidence;
- VAL-026's nine `(human)` rows are all accepted;
- no Blocking finding or escalation remains open; and
- required documentation is synchronized.

Current kickoff verdict: **Ready**. QA-BLK-01 through QA-BLK-07 are closed by
DEC-089..DEC-093; fresh heavy-mode command/fixture feasibility is green, no
new QA-BLK is open, and no required command is a placeholder.
