# Phase 19 QA Plan - Workflow Authoring CLI + Escalation Default Workflow

> Phase: 19 - `fid workflow define/assign/list` + escalation default workflow
> Refreshed: 2026-07-26
> QA Lead: Codex QA Lead
> Operating mode: Heavy
> Plan status: **BLOCKED**

Implementation must not start while this plan is BLOCKED. The command surface
is concrete and feasible, but the validation contract is not yet a valid
verification target and two canonical specifications conflict with the task
pack. Section 2 names the required Tech Lead and Human dispositions. Per
`09_AI_WORKFLOW.md` §2.6, the QA Lead has not edited
`validation-contract.md`; the Tech Lead applies accepted contract refinements.

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
- `project-docs/docs/phases/phase-19/decisions.md` (DEC-075..DEC-083)
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

- All in-scope requirements referenced by at least one `VAL-NNN`: **NO**.
- Every `VAL-NNN` has one validator role: **YES** (VAL-001..025 QA Lead;
  VAL-026 Human).
- Every `VAL-NNN` has a concrete evidence method: **NO**. VAL-003 and VAL-017
  claim parity without a fully specified immutable reference artifact, and
  several bundled rows cannot yield one outcome per observable.
- Every `VAL-NNN` maps to a task/window or phase-end check: **YES**.
- Every `VAL-NNN` describes one observable, or human-checklist fan-out is
  explicitly required: **NO**. VAL-026 is correctly marked for human fan-out;
  multiple deterministic rows are bundled.
- Conditional or binary assertions require both qualifying and
  non-qualifying states: **NO** as a contract-wide verdict. VAL-005, VAL-008,
  VAL-019, VAL-020, VAL-022, and VAL-024 name contrasting states, but
  PRD-WF-013 and DEC-077 conflict on the zero-cancelled-state branch, so
  VAL-008 does not currently have an unambiguous pass state.
- Parity or "matches reference" assertions have deterministic conformance
  proof when available: **NO**. VAL-003's count assertion cannot prove
  byte-identical defaults, and VAL-017 does not name the immutable baseline
  commit, fixture, capture path, or normalization.
- Revisions requested from Tech Lead: **required before implementation**.

### Blocking Refinement And Specification Findings

| ID | Severity | Finding | Required disposition |
|---|---|---|---|
| QA-BLK-01 | Blocking | `02_TRD.md` §3 forbids workflow authoring/mutation in `internal/workflow/`, while Tasks 2-4 explicitly create `internal/workflow/authoring*.go`, `define.go`, and `assign.go`. This is an architecture-boundary contradiction, not a test detail. | Tech Lead records the owning-boundary decision and synchronizes `02_TRD.md` §3 or revises the task pack. The QA Lead does not choose between those architecture changes. |
| QA-BLK-02 | Blocking | PRD-WF-013 requires any issue to be cancellable from any non-cancelled state. DEC-077 and VAL-008 permit a custom workflow with zero `cancelled` states and no cancel edge, and the assignment rules do not restrict that workflow away from issue kinds. A task/epic/subtask/any workflow can therefore make an issue non-cancellable. | Human/Tech Lead resolves the product contract. Examples of non-equivalent choices are requiring one cancelled state for issue-applicable workflows, limiting the zero-state branch to non-issue applicability, or formally revising PRD-WF-013. QA must not infer the choice. |
| QA-BLK-03 | Blocking | Deterministic assertions are bundled. At minimum VAL-001, VAL-003, VAL-004, VAL-007, VAL-009, VAL-010, VAL-011, VAL-014, VAL-015, VAL-016, VAL-017, VAL-018, VAL-019, and VAL-023 contain independently passable/failable observables. | Tech Lead retires/supersedes or splits the rows into sequential `VAL-NNN` assertions without weakening behavior or changing task ownership. Conditional state matrices may remain together only when they prove one conditional rule. |
| QA-BLK-04 | Blocking | The contract omits binding Phase 19 behavior: DEC-081 no-`activity_log` writes for both authoring event types; DEC-076 first-assignment `insert` versus reassignment `delete`+`insert` for kind and project scopes; DEC-080 accepted `--applies-to` vocabulary and `any`→NULL; all four implemented rule names accepted; duplicate-state, duplicate-transition, invalid-JSON, and invalid `all_acs_verified.roles` rejection cases; deterministic list ordering and all-format empty result; and the standing prohibition on production calls into the frozen definitions fixture. | Tech Lead adds narrow assertions or explicitly maps each invariant to an existing single-observable assertion with sufficient evidence. Task scope must not be widened or weakened. |
| QA-BLK-05 | Blocking | Contract source/evidence defects remain. VAL-002 cites TRD-ARCH-008, which is the `issues_enriched` index-plan requirement, not migration reversibility. VAL-009's exactly-one-event behavior needs TRD-ARCH-004 / PRD-OPS-001, not only TRD-ARCH-009. VAL-021's whole-projection replay invariant needs TRD-ARCH-001 rather than the upgrader-chain-only TRD-ARCH-003. VAL-003's default-workflow count cannot prove tuple or byte parity. VAL-012 does not explicitly require all documented aliases (`→`, `⟶`, `-->`, `->`). VAL-017 lacks a reproducible baseline artifact. | Tech Lead corrects sources and evidence methods. VAL-003 must compare all three default graphs tuple-for-tuple against a frozen reference. VAL-017 must name a Task-3 commit, frozen workspace/output artifacts, capture procedure, and normalization (preferably none). |
| QA-BLK-06 | Blocking | `05_DESIGN_PLAN.md` §3.1 is the canonical CLI command tree but has no `workflow` family, and no task owns that update. `02_TRD.md` §3 also needs synchronization if QA-BLK-01 resolves in favor of the task pack. | Tech Lead assigns same-task documentation ownership or records an explicit verified-no-change rationale consistent with DOC-001/002. |
| QA-BLK-07 | Blocking | The validation contract remains `Draft`, while the current phase brief and progress file still say Human approval is pending. A phase branch exists and the Human requested this plan, but the shared artifacts do not record approval of the mode, task list, contract, and CF-008 disposition. | Human approval and Tech Lead artifact refresh must be recorded; the refined contract must be marked `Ready` before Task 1 starts. |

### Non-Blocking Applicability Clarification

`06_DESIGN_PLAN_CLI.md` states that every v0.2 write accepts `--dry-run`, but
`10_IMPLEMENTATION_ROADMAP.md` explicitly schedules universal `--dry-run` in
v0.2e. It is therefore not Phase 19 scope and is not a missing Phase 19
assertion. The final documentation review must preserve that staged wording
and must not imply define/assign support `--dry-run` in v0.2a.

## 3. Task Inventory

| Task | Focus | Covered current VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | `workflow_assignments`, escalation seed, projection version | VAL-001..005 | High | Migration reversibility, exact seed parity, partial-index semantics, and stale/fresh projection branches. |
| Task 2 | Pure authoring and assignment validation/effect core | VAL-006..008 | High | Rejection matrix, deterministic ordering, purity, cancel policy, and 80% `internal/workflow` threshold. |
| Task 3 | Define CLI, event, projection, permissions | VAL-009..013 | High | Reject-before-append, event atomicity, format compatibility, CLI output, and architecture resolution from QA-BLK-01. |
| Task 4 | Assign CLI and assignment-aware issue creation | VAL-014..017 | Very high | Most-used create path, last-write-wins effect shape, precedence, and immutable pre-Task-4 parity evidence. Task 4 cannot start before Review 1. |
| Task 5 | Four-format workflow list | VAL-018..020 | Medium-high | Cross-format structural parity, exact markdown characterization, filters in both states, deterministic order, and empty arrays. |
| Task 6 | Integrity, rebuild, boundaries, exit scenario, docs | VAL-021..026 | High | Whole-projection parity, performance headroom, security/static inventories, no seam behavior, final dogfood, and full spec sync. |

All six task specs have a Validation Contract Coverage section and concrete
command blocks. No required task command is a placeholder.

## 4. Required Evidence Strategy

These rows map the current contract so that no existing assertion is lost
while the Tech Lead applies Section 2 refinements. Superseded/split IDs must be
mapped one-for-one into the same task windows before this plan becomes Ready.

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | Schema introspection and constraint insert results against exact DEC-076 DDL | VAL-001 | QA Lead | Task 1 tests/report; QA-01 |
| EVID-02 | Pre/up/down/up schema and row snapshots with empty diffs | VAL-002 | QA Lead | Task 1 tests/report; QA-02 |
| EVID-03 | Ordered tuples for escalation seed and frozen-reference tuples for every Phase 18 default | VAL-003 | QA Lead | Task 1 tests/report; QA-03 |
| EVID-04 | Seeded loader result with unimplemented flag plus reproducible no-evaluation/no-write state | VAL-004 | QA Lead | Task 1 tests/report; QA-04 |
| EVID-05 | Prior-version rebuild and current-version no-rebuild counterpart output | VAL-005 | QA Lead | Task 1 tests/report; QA-05 |
| EVID-06 | Golden ordered effect slice, repeated-build deep equality, and no-side-effect state | VAL-006 | QA Lead | Task 2 tests/report; QA-06 |
| EVID-07 | Table-driven invalid and valid definition outcomes with value-naming errors | VAL-007 | QA Lead | Task 2 tests/report; QA-07 |
| EVID-08 | Exact transition sets for one/zero/multiple cancelled-state fixtures | VAL-008 | QA Lead | Task 2 tests/report; QA-08 |
| EVID-09 | Built define CLI output, one journal event, ordered four-table tuples, and loader result | VAL-009 | QA Lead | Task 3 tests/report; QA-09 |
| EVID-10 | Per-case journal hashes/bytes and table counts before/after rejection, plus accepted control | VAL-010 | QA Lead | Task 3 tests/report; QA-10 |
| EVID-11 | Exact permission sets, denied/allowed CLI traces, and four-format `policy show` structures | VAL-011 | QA Lead | Task 3 tests/report; QA-11 |
| EVID-12 | Transition tuple equality for every documented arrow alias | VAL-012 | QA Lead | Task 3 tests/report; QA-12 |
| EVID-13 | Byte-exact stderr block and SQL-backed count comparison | VAL-013 | QA Lead | Task 3 tests/report; QA-13 |
| EVID-14 | First assignment/reassignment event and row snapshots, including superseded workflow presence | VAL-014 | QA Lead | Task 4 tests/report; QA-14 |
| EVID-15 | Project/kind fixture with exact assignment rows and both precedence outcomes | VAL-015 | QA Lead | Task 4 tests/report; QA-15 |
| EVID-16 | Per-case assign rejection traces with unchanged journal/rows and accepted control | VAL-016 | QA Lead | Task 4 tests/report; QA-16 |
| EVID-17 | Two initial-state fixtures, all three no-assignment kinds, and immutable pre/post output bytes | VAL-017 | QA Lead | Task 4 tests/report; QA-17 |
| EVID-18 | Four built-CLI outputs decoded to one structure, SQL count comparison, exact markdown expected bytes | VAL-018 | QA Lead | Task 5 tests/report; QA-18 |
| EVID-19 | Exact included and excluded IDs for kind, active, unfiltered, and composed filters | VAL-019 | QA Lead | Task 5 tests/report; QA-19 |
| EVID-20 | Exact JSON field presence/absence for kind, project, and empty assignments | VAL-020 | QA Lead | Task 5 tests/report; QA-20 |
| EVID-21 | Every-table normalized snapshots, schema, FK output, and journal hashes at three rebuild boundaries | VAL-021 | QA Lead | Task 6 integration evidence; QA-21 |
| EVID-22 | Healthy and incompatible maintenance outputs plus `--warn-as-error` exit | VAL-022 | QA Lead | Task 6 tests/CLI trace; QA-22 |
| EVID-23 | Reserved-rule CLI rejections, zero-row SQL output, and static production inventory | VAL-023 | QA Lead | Task 6 report/final review; QA-23 |
| EVID-24 | Custom-state rejection/no-append and v0.1-state success under a custom workflow | VAL-024 | QA Lead | Task 6 integration evidence; QA-24 |
| EVID-25 | Scripted real-workspace trajectory with outputs and unaffected control issue | VAL-025 | QA Lead | Task 6 integration evidence; QA-25 |
| EVID-26 | Nine single-judgment Human outcomes from a scratch real workspace | VAL-026 | Human | Phase-end human checklist |

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
| `go test -race ./...` | Exact execution | PASS | Terminated; terminal package token `ok github.com/mbell-nzl/fid/internal/workflow 1.777s`. |
| `go test -race -coverprofile=coverage.out ./...` | Exact execution | PASS | Terminated; terminal token `ok github.com/mbell-nzl/fid/internal/workflow 1.662s coverage: 83.2% of statements`. Current `cmd/fid` is 77.7%, `internal/issues` 81.5%, and `internal/journal` 81.0%. |
| `go tool cover -func=coverage.out` | Exact execution after the coverage command | PASS | Terminated with `total: (statements) 76.8%`. `go tool cover -help` confirms `-func`; `coverage.out` is gitignored. |
| `go test -tags integration -race ./...` | Exact execution | PASS | Terminated; terminal token `ok github.com/mbell-nzl/fid/internal/workflow 1.695s`. Existing integration build tags were found under `cmd/fid` and `internal/issues`. |
| `go test -bench . -benchmem ./...` | Exact execution | PASS | Terminated; terminal token `ok github.com/mbell-nzl/fid/internal/workflow 0.042s`. `BenchmarkList10kIssuesWithRegistry`: `46.63 ms/registry-read`; `BenchmarkRebuildWorkspace100kEvents`: `25185 ms/rebuild`. |
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
| Tasks 2-5 planned test paths | Parent/file-state checks and task specs | PASS, subject to QA-BLK-01 | `cmd/fid`, `internal/workflow`, `internal/issues`, `internal/policy`, and `internal/format` exist. The paths are creatable, but workflow ownership cannot be accepted until the TRD boundary is resolved. |
| Task 6 planned test paths | Parent/file-state checks and task spec | PASS | `internal/maintenance`, `internal/journal`, and `cmd/fid` exist; prior maintenance and E2E patterns exist. |
| VAL-017 pre-change output reference | Contract/task inspection | BLOCKED | No immutable baseline commit, fixture path, captured output path, hashes, or normalization is specified. QA-BLK-05 requires the Tech Lead to make this reproducible before Task 4. |
| VAL-003 Phase 18 default reference | Contract/task inspection | BLOCKED | The planned count assertion cannot prove byte/tuple identity. Existing frozen definitions can support a tuple comparison, but the contract does not require it. |
| VAL-018 markdown golden | Task/path inspection | PASS with phase-end precondition | Task 5 authorizes `internal/format/markdown_workflow_test.go`; expected bytes may be a committed `testdata` fixture or an explicit literal. The completion report must identify which and the final review must verify it is committed. |
| Static seam inventory | `rg` availability and current inventory | PASS | `rg` is installed. At phase end run `rg -n --glob '*.go' --glob '!**/*_test.go' '\\brequired_gates_current\\b|\\bhas_completion\\b' cmd/fid internal/issues internal/dev internal/journal internal/db internal/maintenance internal/policy`; expected result is no output / exit 1. Reserved declarations and authoring rejections inside `internal/workflow` are expected and reviewed separately. |
| Frozen-definition production-call inventory | `rg` availability and current inventory | PASS | At phase end run `rg -n --glob '*.go' --glob '!**/*_test.go' --glob '!internal/workflow/definitions_fixture.go' 'FrozenDefinitions\\(' cmd internal`; expected result is no output / exit 1. |
| Built-CLI support commands | Exact help execution | PASS | `/tmp/fid init --help`, `issue create --help`, `issue transition --help`, `policy show --help`, `maintenance check --help`, `rebuild --help`, and `sql --help` all terminated with documented flags. |
| Future `fid workflow` family | Current help plus canonical CLI/task specs | PASS with implementation precondition | Current `/tmp/fid workflow --help` exits 1 because the family does not exist yet, as expected. `06_DESIGN_PLAN_CLI.md` and Tasks 3-5 provide concrete define/assign/list syntax. Missing or divergent help after Task 5 is Blocking. |
| Review 1 representative command | `go help testflag`, package-pattern help, and directory checks | PASS | `go test -race ./internal/db/... ./internal/workflow/... ./internal/policy/... ./cmd/fid/...` is syntactically valid; all scopes exist. It is independent review evidence, never a substitute for Developer full gates. |
| Review 2 representative command | `go help testflag`, package-pattern help, and directory checks | PASS | `go test -race ./internal/issues/... ./internal/workflow/... ./internal/format/... ./cmd/fid/...` is syntactically valid; all scopes exist. It is independent review evidence, never a substitute for Developer full gates. |
| Completion reports and QA artifacts | File-state check | PASS as planned artifacts | `task-1-done.md` through `task-6-done.md`, `review-1.md`, `review-2.md`, `final-review.md`, and `human-checklist.md` do not yet exist. Their absence is expected now; each becomes a hard-stop requirement in its owning window. |
| Human scratch workspace | CLI help, design spec, and temp-directory policy | PASS with implementation precondition | `fid init` and supporting issue/transition reads are available; workflow syntax is concrete in the CLI spec. Use a fresh path returned by `mktemp -d /tmp/fid-phase19-human.XXXXXX`; do not reuse project data. |

Invalid or infeasible items are not carried forward in heavy mode. The two
blocked artifact assumptions and all Section 2 blockers must be resolved and
freshly rechecked before this plan's status changes to Ready.

## 7. Deterministic Checklist

Each current `VAL-NNN` has exactly one QA Lead step below. If the Tech Lead
splits an assertion, add one step for every resulting assertion before
implementation; do not let several new assertions inherit one ambiguous
outcome.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Verify exact assignment schema and scoped uniqueness in accepted and rejected states `(qa-lead)` | VAL-001 | Migration SQL, introspection output, duplicate-kind/project failures, distinct-kind success. |
| QA-02 | Verify migration `0009`/`0010` up-down-up reversibility `(qa-lead)` | VAL-002 | Three state snapshots and empty deterministic diffs. |
| QA-03 | Verify escalation seed conformance and frozen-default parity `(qa-lead)` | VAL-003 | Ordered tuples for all four tables and every default, not counts alone. |
| QA-04 | Verify `linked_to_decision` loads structurally as unimplemented without evaluation or mutation `(qa-lead)` | VAL-004 | Test output plus before/after DB state. |
| QA-05 | Verify stale projection rebuild and fresh projection no-rebuild branches `(qa-lead)` | VAL-005 | Both branch outputs and current schema-version state. |
| QA-06 | Verify authoring effect order and repeated-build determinism `(qa-lead)` | VAL-006 | Golden slice diff and repeated deep-equality output. |
| QA-07 | Verify every contracted core rejection and one accepted control `(qa-lead)` | VAL-007 | Named table-driven subtests and distinct value-naming errors. |
| QA-08 | Verify cancel-edge policy for exactly one, zero, and multiple cancelled states `(qa-lead)` | VAL-008 | Exact transition sets and rejection output, after QA-BLK-02 resolution. |
| QA-09 | Verify successful define event count, projected tuple conformance, and load result `(qa-lead)` | VAL-009 | Built-CLI trace, journal event, SQL tuples, loader output. |
| QA-10 | Verify every define rejection leaves journal bytes and all four table counts unchanged `(qa-lead)` | VAL-010 | Per-case before/after hashes/counts/errors and accepted control. |
| QA-11 | Verify exact define/assign permissions and four-format policy visibility `(qa-lead)` | VAL-011 | Matrix test output, denied/allowed traces, decoded outputs. |
| QA-12 | Verify all documented transition arrow aliases normalize to identical tuples `(qa-lead)` | VAL-012 | `->`, `→`, `⟶`, and `-->` tuple-set comparison. |
| QA-13 | Verify byte-exact define confirmation and SQL-backed counts `(qa-lead)` | VAL-013 | Literal stderr block and count query output. |
| QA-14 | Verify kind assignment and reassignment state `(qa-lead)` | VAL-014 | Journal effects/events, assignment rows, and superseded workflow presence. |
| QA-15 | Verify project assignment and kind-over-project precedence in both directions `(qa-lead)` | VAL-015 | Exact rows and resulting issue workflow IDs for covered/uncovered kinds. |
| QA-16 | Verify every assign rejection is pre-append and one compatible assign succeeds `(qa-lead)` | VAL-016 | Per-case journal/row before-after evidence and exact errors. |
| QA-17 | Verify assignment-aware initial state and no-assignment v0.1 parity `(qa-lead)` | VAL-017 | Two assigned fixtures, all three unassigned kinds, immutable baseline reference and byte diffs. |
| QA-18 | Verify list renders one equivalent structure in all four formats with exact markdown bytes `(qa-lead)` | VAL-018 | Decoded JSON/TOON/Markdown/YAML comparison, SQL counts, committed expected bytes. |
| QA-19 | Verify kind and active filters include and exclude exact sets in filtered/unfiltered states `(qa-lead)` | VAL-019 | Result-set IDs for every qualifying and non-qualifying branch. |
| QA-20 | Verify `assigned_to` discriminator shapes for kind, project, and unassigned states `(qa-lead)` | VAL-020 | Exact JSON key-presence and empty-array output. |
| QA-21 | Verify non-destructive deterministic whole-projection rebuild twice `(qa-lead)` | VAL-021 | Every-table snapshots, schema/FK output, and per-journal-file hashes at all boundaries. |
| QA-22 | Verify maintenance assignment audit on healthy and incompatible states `(qa-lead)` | VAL-022 | Healthy output, deterministic finding, and non-zero warn-as-error exit. |
| QA-23 | Verify gate-obligation seam remains behavior-free `(qa-lead)` | VAL-023 | Two CLI rejections, zero-row SQL, exact `rg` inventory and classification. |
| QA-24 | Verify custom-state permission boundary rejects and v0.1-name path succeeds `(qa-lead)` | VAL-024 | Exact error, unchanged journal bytes, successful control transition/state. |
| QA-25 | Execute the scripted slice-2 exit journey and unaffected control `(qa-lead)` | VAL-025 | Integration output and each command/result/projected-state checkpoint. |
| QA-26 | Verify package coverage thresholds, not only aggregate coverage `(qa-lead)` | no VAL | Fresh `go tool cover` output for `cmd/fid`, issues, claims, workflow, policy, journal. |
| QA-27 | Verify 10k-read and 100k-rebuild controls and headroom `(qa-lead)` | no VAL | Fresh benchmark lines, targets, Phase 18 and planning-baseline comparison. |
| QA-28 | Review SQL/command injection, unsafe input, data-loss, and journal-atomicity risks `(qa-lead)` | no VAL | Parameterization tests/traces, rejection/no-append evidence, security findings. |
| QA-29 | Verify file/function budgets, exported docs, dependency budget, and architecture boundaries `(qa-lead)` | no VAL | Reproducible file-state counts, lint output, dependency diff, resolved boundary decision. |
| QA-30 | Verify required documentation and help are synchronized `(qa-lead)` | no VAL | Diffs/captures for PRD, TRD, schema, design command tree, CLI spec, glossary, roadmap wrap obligations. |
| QA-31 | Audit all six completion reports against exact task commands and committed hashes/subjects `(qa-lead)` | no VAL | Task-spec/report comparison, terminal tokens, `git rev-parse` and `git log` output. |

## 8. Mid-Phase Review Cadence

Two planned QA Reviewer windows apply:

1. **Review 1: Tasks 1-3**, written to
   `project-docs/docs/phases/phase-19/qa/review-1.md` after all three
   committed completion reports exist. It must complete with `VERDICT:
   CONTINUE` **before any Task 4 work starts**. It reviews VAL-001..013,
   every Task 1-3 command result and terminal token, migration/event
   integrity, the resolved architecture decision, and Task 4 baseline
   readiness. Independent representative command:
   `go test -race ./internal/db/... ./internal/workflow/... ./internal/policy/... ./cmd/fid/...`.
2. **Review 2: Tasks 4-5**, written to
   `project-docs/docs/phases/phase-19/qa/review-2.md` after both committed
   completion reports exist and before Task 6 starts. It reviews VAL-014..020,
   every Task 4-5 command result, immutable v0.1 parity evidence, assignment
   event shape/precedence, and four-format list conformance. Independent
   representative command:
   `go test -race ./internal/issues/... ./internal/workflow/... ./internal/format/... ./cmd/fid/...`.

Task 6 is the closure window and is reviewed by the QA Lead in the final
review, not by a redundant post-Task-6 mid-phase review. Any Blocking finding
stops the next window. Each review uses
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
- execute QA-01..QA-31 and cite reproducible output/file-state evidence for
  every result;
- record a final `PASS`, `FAIL`, or `BLOCKED` for every refined assertion,
  with no evidence-free pass;
- accept **no heavy-mode carry-forward**. An unmet assertion, unavailable
  required fixture, non-terminating command, or unresolved defect blocks the
  phase and returns to Tech Lead/Human re-planning;
- review requirement compliance, tests, journal/schema compatibility,
  architecture and layer boundaries, security/data-loss/input handling,
  performance, dependencies, structure budgets, and documentation accuracy;
- verify Task 6 records the 10k-read and 100k-rebuild measurements against
  both controls and prior baselines without weakening either target;
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
- every refined `VAL-NNN` is PASS with cited evidence;
- every `(qa-lead)` item is PASS with reproducible evidence;
- VAL-026's nine `(human)` rows are all accepted;
- no Blocking finding or escalation remains open; and
- required documentation is synchronized.

Current kickoff verdict: **BLOCKED** by QA-BLK-01 through QA-BLK-07. Command
configuration itself is feasible and has no placeholder defect.
