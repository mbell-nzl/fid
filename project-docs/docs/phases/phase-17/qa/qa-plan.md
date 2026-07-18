# Phase 17 QA Plan - Maintenance, Stub Guards, and v0.1 Closure

> Phase: 17 - Maintenance + stub guard + v0.1 closure
> Refreshed: 2026-07-14
> QA Lead: Codex
> Operating mode: Heavy (DEC-049)
> Status: **READY FOR TASK 0; TASK 1 BLOCKED PENDING ENTRY REVIEW**
> Status note (Tech Lead, 2026-07-15): The banner above predates Task 0
> completion and is stale. Task 0 is committed and the Entry Review is
> committed with `VERDICT: CONTINUE` (see `qa/review-entry.md` and the run
> ledger), so QA-BLK-01 is closed and Task 1 is authorized — currently paused
> only at the ESC-004 report-commit fix, not at entry review. QA Lead owns the
> authoritative banner and `Refreshed` update at the next QA touch and the
> QA-18 phase-end spec sync.

## Blocking Entry Findings

The evidence design, commands, fixtures, artifacts, and human checks are
feasible. DEC-055's required model qualification is now committed and closes
QA-BLK-02. QA-BLK-01 remains open at the evidenced 78.9% baseline, but DEC-054
authorizes tests-only Task 0 as the remediation path. This plan authorizes Task
0 only; it does not claim that COV-001 already passes or authorize Task 1.

| ID | Finding and status | Evidence | Required resolution / disposition |
|---|---|---|---|
| QA-BLK-01 | **OPEN; Task 0 authorized.** `internal/journal` remains below COV-001: 78.9% against the unchanged 80% package threshold. Task 0 is the tests-only entry remediation; Tasks 1-9 do not authorize journal coverage work and Task 9 remains audit-only. | Fresh focused feasibility on 2026-07-14: `go test -race -coverprofile=/tmp/fid-journal-coverage.out ./internal/journal/...` passed tests but reported `coverage: 78.9%`; `go tool cover -func=/tmp/fid-journal-coverage.out` reported `total ... 78.9%`. The unchanged full-suite feasibility run on 2026-07-11 also reported `internal/journal ... coverage: 78.9%`. | Task 0 must add meaningful `internal/journal/*_test.go` coverage only, report every exact command, prove `internal/journal >=80.0%` in both focused and unchanged full-suite runs, and prove production files unchanged. `review-entry.md` must independently accept test quality and evidence before QA-BLK-01 closes or Task 1 starts. No production change, threshold change, or waiver is authorized. |
| QA-BLK-02 | **CLOSED; ESC-009 qualification condition satisfied.** DEC-053/055 model qualification is satisfied for QA Lead `gpt-5.6-sol` xhigh and QA Reviewer `gpt-5.6-terra` high. | Commit `2d5e536fd523682e096909e1c2defc03b25fbc55` (`docs(qa): record phase 17 model qualification`) commits `qa/model-audit.md` and five verbatim evidence files. The Human scorecard records PASS/DETECTED for mechanical drills 00..04, Reviewer drills 05..08, and QA Lead drill 08, with no MISS; the audit explicitly verifies seeded-defect detection, deterministic evidence, blocker/VAL disposition, cross-cutting escalation, and preservation of both-state hard stops. | Satisfied by the committed audit, which explicitly states that it is sufficient to close QA-BLK-02/ESC-009. DEC-055 permits no routing fallback, waiver, or inferred qualification; none is used here. Fresh qualified role sessions remain required by DEC-053. |

### Entry Criteria And Authorization

**Task 0 is Ready** once this revised QA plan is committed. Its spec is
concrete, its focused and full-suite commands are feasible, QA-BLK-02 is
closed and ESC-009's qualification condition is satisfied by commit `2d5e536`,
and DEC-054 expressly authorizes only meaningful
`internal/journal/*_test.go` changes. The current 78.9% result is the red entry
condition, not a passing threshold result.

**Task 1 remains blocked** until all of the following are true:

- Task 0's test commit and finalized `task-0-done.md` are committed.
- Every Task 0 exact command is recorded `PASS`; both focused and unchanged
  full-suite evidence report `internal/journal >=80.0%`.
- Commit/file-state evidence proves no non-test Go file or compatibility
  surface changed and no threshold or coverage policy changed.
- The exact Task 0 report-lint command passes.
- A fresh QA Reviewer commits `review-entry.md` with no unresolved Blocking
  finding and `VERDICT: CONTINUE` after independently reviewing test quality,
  reproducing representative coverage evidence, and confirming both focused
  and full-suite threshold proof.

Only that evidence closes QA-BLK-01. A partial increase below 80%, a
focused-only pass, or positive test output without the full coverage evidence
keeps Task 1 blocked.

The stale roadmap wording `--format json|markdown` is a non-blocking
documentation defect because DEC-051, the PRD output contract, design spec, CLI
spec, validation contract, and Task 5 all unambiguously require
json/toon/markdown/yaml. QA-18 requires that wording to be reconciled during
the phase's full-scope spec sync.

## 1. Scope And Source Documents

Phase 17 closes v0.1 early dogfood by adding a non-repairing maintenance audit,
locked derived-DB VACUUM, deterministic v0.2 governance guards, all-read-format
conformance, and a full lifecycle/rebuild/telemetry proof.

Sources loaded for this plan:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/docs/phases/phase-17/progress.md`
- `project-docs/docs/phases/phase-17/decisions.md`
- `project-docs/docs/phases/phase-17/qa/validation-contract.md`
- `project-docs/docs/phases/phase-17/qa/model-audit.md` and committed adjacent
  evidence files (commit `2d5e536`)
- `project-docs/docs/phases/phase-17/tasks/task-0.md` through `task-9.md`
- `project-docs/specs/08_QUALITY_PLAN.md` sections 3.4, 4.1-4.5, 5.1-5.9,
  7, 9, and 10
- `project-docs/specs/09_AI_WORKFLOW.md` sections 2.6, 4, 6.1, and 8
- Requirement excerpts named by the contract from `01_PRD.md`, `02_TRD.md`,
  `03_TRD_SCHEMA.md`, `04_TRD_CONFIG.md`, `05_DESIGN_PLAN.md`,
  `06_DESIGN_PLAN_CLI.md`, `07_DESIGN_PLAN_FLOWS.md`,
  `10_IMPLEMENTATION_ROADMAP.md`, and `GLOSSARY.md`
- `project-docs/templates/qa-plan-template.md`,
  `mid-phase-review-template.md`, `final-review-template.md`,
  `human-checklist-template.md`, and `completion-report-template.md`

Public-release threat modelling, release scans, pre-release checklist, and
held-out scenarios remain deferred to v0.3 by DEC-049. This phase still runs
ordinary real-work human dogfood for VAL-025.

## 2. Validation Contract Refinement Check

- All in-scope PRD/TRD requirements referenced by at least one `VAL-NNN`: YES.
- Every `VAL-NNN` has one validator role: YES.
- Every `VAL-NNN` has a concrete evidence method: YES.
- Every `VAL-NNN` maps to a task/window or phase-end check: YES.
- Every `VAL-NNN` describes one observable, or human fan-out is required: YES;
  VAL-025 fans out into HUMAN-25a..e.
- Conditional/binary assertions name qualifying and non-qualifying states:
  YES for VAL-002..010, VAL-012..016, and VAL-018..023.
- Parity assertions have deterministic conformance proof: YES for VAL-001,
  VAL-011, VAL-017, VAL-022, and VAL-023.
- Every Task 0-9 cites its owned assertions: YES; Task 0 supports VAL-024 and
  does not claim a terminal VAL-024 outcome.
- Every VAL-001..025 has a QA-plan verification step: YES; see section 7 and
  the coverage rollup in section 11.
- One-line rollup: **all in-scope requirements are referenced by at least one
  VAL assertion: YES.**
- Revisions requested from Tech Lead: none to assertion wording or ownership.
  DEC-054 provides the tests-only execution path for QA-BLK-01 without
  weakening VAL-024. QA-BLK-02 is closed by DEC-055's required committed audit
  evidence; no assertion refinement is needed.

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA notes |
|---|---|---|---|---|
| 0 | Restore journal COV-001 coverage with tests only | VAL-024 | High | Entry remediation only: meaningful `internal/journal/*_test.go` changes, focused and unchanged full-suite coverage >=80%, and production-byte/file-state proof; mandatory entry review before Task 1. |
| 1 | Behavior-preserving formatter split | VAL-001 | Medium | Review byte parity and every resulting non-test file <=400 lines. |
| 2 | Maintenance grace config | VAL-002 | High | Review legacy/missing, zero, positive, negative, malformed, and bootstrap round-trip states. |
| 3 | Structural maintenance errors | VAL-003..006 | High | Trigger/FK bypass must be test-only; valid counterpart required for every violation family. |
| 4 | Warnings and info | VAL-007..011 | High | Inject time; prove boundary and excluded states; deep-equality ordering. |
| 5 | Maintenance check CLI/formats | VAL-012..014 | High | Structured stdout must survive expected non-zero outcomes; journal stays byte-identical. |
| 6 | Locked VACUUM | VAL-015 | High | Authorization before mutation, one lock, data/journal preservation, policy parity. |
| 7 | Governance guards | VAL-016 | High | Exact command matrix and precedence; no validation or persistence access before gate. |
| 8 | Read-format conformance | VAL-017 | High | Exact 15-command matrix; structural machine-format parity and markdown identities. |
| 9 | v0.1 closure journey | VAL-018..023 | High | Audit-only; any product defect stops for triage. Temp git and delete/rebuild are mandatory. |
| Phase end | Traceability, full gates, size | VAL-024 | High | Every v0.1 PRD row maps to reproducible evidence; Task 0 and entry-review evidence are reverified, but terminal VAL-024 remains QA Lead-owned. |
| Human | Real-work dogfood | VAL-025 | High | Run only after deterministic QA; one judgment per HUMAN row. |

## 4. Required Evidence Strategy

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-00 | Meaningful failing-capable journal behavioral tests; exact focused and unchanged full-suite outputs both showing `internal/journal >=80.0%`; commit/file-state proof that production files, compatibility surfaces, and thresholds are unchanged; exact report-lint PASS | VAL-024 | Developer, QA Reviewer, QA Lead | Task 0 / entry review / final review |
| EVID-01 | Before/after golden and existing formatter outputs plus `wc -l` of all split files | VAL-001 | Developer, QA Reviewer | Task 1 / review 1 |
| EVID-02 | Config table tests for missing, zero, positive, negative, malformed, bootstrap, and legacy prefix-only input | VAL-002 | Developer, QA Reviewer | Task 2 / review 1 |
| EVID-03 | One corrupt and one valid fixture per hierarchy/scheduling family, with full ordered finding payload | VAL-003 | Developer, QA Reviewer | Task 3 / review 2 |
| EVID-04 | Orphan/resolved polymorphic refs and present/missing journal-event IDs | VAL-004 | Developer, QA Reviewer | Task 3 / review 2 |
| EVID-05 | All rows from `PRAGMA foreign_key_check` plus clean-DB zero rows under normal pragmas | VAL-005 | Developer, QA Reviewer | Task 3 / review 2 |
| EVID-06 | Every default workflow/state pair plus unknown workflow and unknown state | VAL-006 | Developer, QA Reviewer | Task 3 / review 2 |
| EVID-07 | Injected-time before/at/after TTL evidence | VAL-007 | Developer, QA Reviewer | Task 4 / review 2 |
| EVID-08 | Live claim, before/at/after grace, grace zero, fallback creation time, and non-status-update cases | VAL-008 | Developer, QA Reviewer | Task 4 / review 2 |
| EVID-09 | Superseded/current/NULL completion verification states | VAL-009 | Developer, QA Reviewer | Task 4 / review 2 |
| EVID-10 | Missing completion link, issue-only link, and completion-linked run states | VAL-010 | Developer, QA Reviewer | Task 4 / review 2 |
| EVID-11 | Known counts and two deeply equal ordered audit results | VAL-011 | Developer, QA Reviewer | Task 4 / review 2 |
| EVID-12 | Healthy and corrupt built-CLI results in four formats with status/exit/payload agreement | VAL-012 | Developer, QA Reviewer | Task 5 / review 2 |
| EVID-13 | Same warning fixture with/without `--warn-as-error`; identical warnings payload | VAL-013 | Developer, QA Reviewer | Task 5 / review 2 |
| EVID-14 | Missing/stale projection refresh with journal hash and line count unchanged | VAL-014 | Developer, QA Reviewer | Task 5 / review 2 |
| EVID-15 | Authorized/denied/unknown/missing-role, lock contention, snapshots, sizes, policy parity, and no journal/activity delta | VAL-015 | Developer, QA Reviewer | Task 6 / review 3 |
| EVID-16 | Exact 10-command guard matrix, precedence probes, mutation counts, help visibility, and requirement-sync control | VAL-016 | Developer, QA Reviewer | Task 7 / review 3 |
| EVID-17 | Exact 15-read x 4-format matrix with structural decode, primary identities, stable order, and empty/non-empty parity | VAL-017 | Developer, QA Lead | Task 8 / final review |
| EVID-18 | Built-CLI lifecycle with permission and unmet-rule zero-append branches | VAL-018 | Developer, QA Lead | Task 9 / final review |
| EVID-19 | Sequential and prioritised ready positive/negative state matrix | VAL-019 | Developer, QA Lead | Task 9 / final review |
| EVID-20 | Label/link/inverse read-back and invalid-relation zero-append proof | VAL-020 | Developer, QA Lead | Task 9 / final review |
| EVID-21 | Priced run snapshot and independently calculated non-NULL cost/token aggregates | VAL-021 | Developer, QA Lead | Task 9 / final review |
| EVID-22 | Pre/post delete-and-rebuild structural snapshots and equal maintenance output | VAL-022 | Developer, QA Lead | Task 9 / final review |
| EVID-23 | Temp-git ignore/track proof and journal hashes across check, vacuum, guards, reads, and rebuild | VAL-023 | Developer, QA Lead | Task 9 / final review |
| EVID-24 | Every v0.1 PRD row mapped to task/test/review evidence; Task 0 and entry-review evidence reverified; exact phase gate, COV-001/COV-002 values, static size | VAL-024 | QA Lead | Final review |
| EVID-25 | Five separately recorded real-work judgments | VAL-025 | Human | Human checklist |
| EVID-26 | Exact task command coverage and valid commit hash/subject from all ten committed completion reports | no VAL | QA Reviewer, QA Lead | Entry review, each review, and final review |

## 5. Deterministic Gates

No required command contains a placeholder. Every task report must record every
command from its task spec verbatim as `PASS`, or as `BLOCKED` in a blocker
report. Narrower substitutions do not count.

| Gate | Command | Expected result | Required when |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | No files listed | Every task and phase end |
| Lint | `golangci-lint run ./...` | `0 issues.` | Every task and phase end |
| Focused journal coverage execution | `go test -race -coverprofile=/tmp/fid-journal-coverage.out ./internal/journal/...` | Tests PASS and package reports >=80.0% | Task 0 and entry review |
| Focused journal coverage report | `go tool cover -func=/tmp/fid-journal-coverage.out` | Exit 0 and total reports >=80.0% | Task 0 and entry review |
| Coverage execution | `go test -race -coverprofile=coverage.out ./...` | All packages PASS | Every task and phase end |
| Coverage report | `go tool cover -func=coverage.out` | Exit 0; COV-001 named packages >=80%, `cmd/fid` >=70%, no threshold decrease | Every task and phase end |
| Unit suite | `go test -race ./...` | All packages PASS | Phase end, per quality plan 5.2 |
| Integration/E2E | `go test -tags integration -race ./...` | All packages PASS | Tasks 5-9 as specified; phase end |
| Benchmarks | `go test -bench . -benchmem ./...` | All benchmarks PASS; compare TRD-PERF targets | Task 9 and phase end |
| Build | `go build -o /tmp/fid ./cmd/fid` | Exit 0 and binary exists | Tasks 5-9 as specified; phase end |
| Static build | `CGO_ENABLED=0 go build -o /tmp/fid-static ./cmd/fid` | Exit 0 | Task 9 and phase end |
| Static size | `stat -c '%s %n' /tmp/fid-static` | Size evidence emitted | Task 9 and phase end |
| Size assertion | `test "$(stat -c %s /tmp/fid-static)" -lt 31457280` | Exit 0 | Task 9 and phase end |
| Formatter size | `wc -l internal/format/markdown.go internal/format/markdown_issue.go internal/format/markdown_lists.go internal/format/markdown_stats.go` | Every row <=400 | Task 1 / review 1 / final |
| Completion report lint | Ten exact `project-docs/tools/report-lint.sh` commands listed in section 8 | Exit 0 for every task report/spec pair | Entry review, each review, and final |

Security review at phase end covers SQL construction/input binding,
test-only corruption bypass isolation, guard-before-input behavior, lock
handling, path handling, and journal/data-loss risk. Public-release scanners
remain deferred by DEC-049.

## 6. Command And Feasibility Validation

Fresh Heavy Mode validation was performed on 2026-07-11, with the newly added
Task 0 and model-audit paths validated on 2026-07-14; no prior-phase command
result was carried forward.

### 6.1 Command Feasibility

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | Exact command | PASS | Exit 0; no files listed. `gofmt -h` confirms `-l` and path arguments. |
| `golangci-lint run ./...` | Exact command and help | PASS | Exit 0, `0 issues.`; `golangci-lint run --help` confirms syntax. |
| `go test -race -coverprofile=/tmp/fid-journal-coverage.out ./internal/journal/...` | Exact command on 2026-07-14 | **BLOCKED for threshold, command feasible** | Tests passed and the profile was produced, but the package reported `coverage: 78.9% of statements`. This is Task 0's red entry condition, not a pass. |
| `go tool cover -func=/tmp/fid-journal-coverage.out` | Exact command on 2026-07-14 | **BLOCKED for threshold, command feasible** | Exit 0 and function rows were emitted; `total ... 78.9%`. Task 0 and entry review require >=80.0%. |
| `go test -race -coverprofile=coverage.out ./...` | Exact command | **BLOCKED for threshold, command feasible** | All tests passed; `cmd/fid` 77.5%; `internal/journal` 78.9%. The latter creates QA-BLK-01. |
| `go tool cover -func=coverage.out` | Exact command and help | **BLOCKED for threshold, command feasible** | Exit 0; total 75.6%. `go tool cover -help` confirms `-func`. COV-001 package assessment must use the package totals emitted by the preceding `go test`. |
| `go test -race ./...` | `go help testflag` plus successful stricter coverage invocation | PASS | Help confirms `-race`; coverage invocation ran the same suite with race instrumentation. Exact command remains mandatory at phase end. |
| `go test -tags integration -race ./...` | Exact command | PASS | All packages passed; `cmd/fid` completed in 67.006s. Existing integration-tagged v0.1 journeys were exercised. |
| `go test -bench . -benchmem ./...` | Exact command and help | PASS | All packages passed; journal 1k and 10k benchmarks ran at 408.7 ms and 4047 ms per rebuild in this feasibility run. |
| `go build -o /tmp/fid ./cmd/fid` | Exact command and `go help build` | PASS | Exit 0 and `/tmp/fid` was produced. |
| `CGO_ENABLED=0 go build -o /tmp/fid-static ./cmd/fid` | Exact command | PASS | Exit 0. |
| `stat -c '%s %n' /tmp/fid-static` | Exact command and `stat --help` | PASS | `12738389 /tmp/fid-static`; `%s` and `%n` supported. |
| `test "$(stat -c %s /tmp/fid-static)" -lt 31457280` | Exact command and shell `help test` | PASS | Exit 0; current static binary is below 30 MiB. |
| Task 1 `wc -l` command | `wc` availability, existing source check, task output-path check | PASS | `internal/format/markdown.go` exists; the other three paths are explicit Task 1 outputs and must exist before its exact command runs. |
| Ten exact report-lint commands in section 8 | Script inspection and `project-docs/tools/report-lint.sh --self-test` on 2026-07-14 | PASS | Self-test passed required-command, narrower-substitution, blocker-report, hash/subject, uncommitted, and dirty-report cases; all ten task/report paths are declared by Task 0-9 specs. |
| CLI human prerequisites | Built CLI help | PASS | Root, `init`, `ready`, `rebuild`, and `requirement sync` help all exit 0. Planned `maintenance` and governance syntax is fixed by DEC-051/052 and Tasks 5/7. |

### 6.2 Fixture Feasibility

| Fixture or scope | Validation method | Result | Evidence |
|---|---|---|---|
| Trigger-bypassed structural fixture | Disposable temp workspace plus SQLite trigger/FK probe | PASS | Dropped only fixture-copy triggers, inserted `issues` row with missing parent, and `PRAGMA foreign_key_check` returned `issues|1|issues|0`. Normal connection contract remains `foreign_keys=1`, `recursive_triggers=1`. |
| Required migrations | File existence checks | PASS | `0001_v01a_base`, `0004_v01_kind_triggers`, `0005_v01_ref_triggers`, and `0006_v01b_sequence_scheduling` up migrations exist. |
| Valid-state controls | Schema/workflow inspection and current tests | PASS | Active CHECK/trigger rules and `internal/workflow` definitions exist; Task 3 explicitly requires one valid counterpart per family. |
| Time boundaries | Existing claim/config patterns plus task-owned injected clock helper | PASS | Current tests use temp DBs and deterministic option/time seams; Task 4 owns `test_helpers_test.go` for injected time. No wall-clock sleep is accepted as boundary proof. |
| Lock contention | Existing lock tests and integration run | PASS | `internal/journal/lock_test.go` proves held-lock timeout and delayed release; Task 6 can reuse the established held-lock pattern. |
| Projection freshness/rebuild | Existing journal and e2e tests plus integration run | PASS | `EnsureProjectionFreshWithLock`, delete/rebuild snapshots, and temp-workspace CLI helpers already exist and passed integration. |
| Format decode/parity | Existing format modules and CLI tests | PASS | JSON, TOON, YAML, and markdown formatters/tests exist; command tests already decode machine formats. Task 8 adds test-only cross-command structural helpers. |
| Journal hash/line evidence | Existing journal helpers and closure tests | PASS | Canonical SHA-256, journal byte reads, line counts, and rebuild equality helpers exist. |
| Temp git fixture | Disposable git repo and built CLI | PASS | `fid init` produced a trackable journal and ignored `.fid/fid.sqlite`; `git add --dry-run .fid` listed `.gitignore`, config, pricing, and journal only. |
| Static binary fixture | Exact static build and size assertion | PASS | `/tmp/fid-static` built at 12,738,389 bytes, below 31,457,280. |

### 6.3 Artifact Feasibility

| Artifact assumption | Validation method | Result | Evidence |
|---|---|---|---|
| Validation target | File and VAL audit | PASS | `validation-contract.md` exists with VAL-001..025, singular owners, evidence, and terminal criteria. |
| Task evidence inputs | File inventory | PASS | Ten Task 0-9 specs exist; each has exact commands and a completion-report path. No completion report is expected before its task runs. |
| Task 0 tests-only boundary | Task spec/file-scope inspection plus Git/report-lint capability check | PASS | Task 0 permits only `internal/journal/*_test.go`; its verified task commit file list and finalized report can prove that no production file, compatibility surface, or threshold changed. Any product change is a hard stop. |
| Mid-phase reviews | Template/path validation | PASS | `mid-phase-review-template.md` exists; `review-entry.md`, `review-1.md`, `review-2.md`, and `review-3.md` are feasible at the checkpoints in section 8. |
| Final review | Template/path validation | PASS | `final-review-template.md` includes gate results, every VAL outcome, deterministic execution, and maturity disposition. |
| Human checklist | Template/path validation | PASS | `human-checklist-template.md` supports one observable per row and human-only outcomes. |
| Completion report verification | Template/tool validation | PASS | Completion template and executable report-lint tool exist; self-test is green. |
| Commit evidence | Git/tool validation | PASS | Git 2.53.0 is available; report-lint verifies hash, subject, committed state, and exact command coverage. |
| QA artifact ownership | Repository/write check | PASS | This QA plan can be committed independently with `docs(qa): ...` and `Refs: phase-17`; unrelated `run-ledger.md` remains untouched. |
| DEC-053/055 model-audit evidence | `git rev-parse`, commit metadata, committed-tree inventory, and audit/scorecard inspection on 2026-07-14 | PASS | `2d5e536fd523682e096909e1c2defc03b25fbc55` is a valid commit with subject `docs(qa): record phase 17 model qualification`; it commits the audit and five referenced outputs, and the Human scorecard records no MISS. QA-BLK-02 is closed. |

### 6.4 Human-Check Feasibility

| Human check | Validation method | Result | Evidence |
|---|---|---|---|
| HUMAN-25a maintenance finding clarity | Approved CLI/design contract and Task 5 syntax | PASS, contingent on deterministic QA | Human can run built `fid maintenance check --format markdown` during real work after VAL-012..014 pass. |
| HUMAN-25b version-gate clarity | DEC-052 exact first line and Task 7 matrix | PASS, contingent on deterministic QA | A guard invocation is safe/non-mutating once VAL-016 passes. |
| HUMAN-25c ready-choice sensibility | Current built `fid ready --help` and existing real-work command | PASS | Help exits 0 and supports project/epic scope. Human judges the selected work, not mechanics already covered by VAL-019. |
| HUMAN-25d recovery-path clarity | Current built `fid rebuild --help`, temp-git fixture, Task 9 disposable-copy journey | PASS, contingent on deterministic QA | Human uses a disposable copy/temp workspace; no request to delete the real dogfood DB. |
| HUMAN-25e overall v0.1 acceptance | Charter/quality-plan dogfood gate and four prior rows | PASS, contingent on rows a-d | Separate product-owner disposition is feasible after each discrete judgment. |

## 7. Deterministic Checklist

Every item is executed at final review. Review-owned VALs may cite the committed
mid-phase review, but QA Lead still verifies the evidence is reproducible and
records a terminal outcome for every VAL.

| ID | Check | Covers | Evidence required |
|---|---|---|---|
| QA-00 | Reverify Task 0's meaningful journal tests, focused and unchanged full-suite `internal/journal >=80.0%` results, tests-only commit scope, exact command coverage, and accepted entry review `(qa-lead)` | VAL-024 | `task-0-done.md`, exact focused/full outputs, verified commit file list, report-lint PASS, and committed `review-entry.md`; 78.9% or focused-only evidence cannot pass. |
| QA-01 | Verify byte-identical pre-Phase-17 markdown outputs and all split formatter files <=400 lines `(qa-lead)` | VAL-001 | Review 1, golden/existing test output, final `wc -l`. |
| QA-02 | Verify missing/zero/positive/negative/malformed config and bootstrap/legacy round-trip states `(qa-lead)` | VAL-002 | Focused test output naming every state. |
| QA-03 | Verify every structural violation family and valid counterpart, including all FK-check rows `(qa-lead)` | VAL-003..006 | Review 2 and focused domain test output; fixture setup excerpt. |
| QA-04 | Verify every warning qualifying/non-qualifying state and exact boundary `(qa-lead)` | VAL-007..010 | Review 2 and focused injected-time/domain test output. |
| QA-05 | Verify exact info counts and repeated deep-equality ordering `(qa-lead)` | VAL-011 | Focused test output and deterministic result excerpt. |
| QA-06 | Verify maintenance healthy/error/warning promotion in four formats and fresh/read-only journal behavior `(qa-lead)` | VAL-012..014 | Review 2, built CLI/integration outputs, decoded structures, before/after journal hash/line evidence. |
| QA-07 | Verify vacuum authorization, policy parity, lock contention, size output, and data/journal/activity preservation `(qa-lead)` | VAL-015 | Review 3, focused tests, snapshots, hashes, size output. |
| QA-08 | Verify all governance guards gate first without mutation and requirement sync still writes `(qa-lead)` | VAL-016 | Review 3, exact matrix output, count/hash deltas, control output. |
| QA-09 | Execute/verify the exact 15-read x 4-format conformance matrix `(qa-lead)` | VAL-017 | Integration output, structural parity summary, markdown identity rows. |
| QA-10 | Verify full primitive lifecycle and zero-append rejection branches `(qa-lead)` | VAL-018 | Closure test output and journal line deltas. |
| QA-11 | Verify sequential/prioritised ready positive and negative states `(qa-lead)` | VAL-019 | Closure test subcase output and selected IDs. |
| QA-12 | Verify labels, links, inverse read-back, and invalid relation no-op `(qa-lead)` | VAL-020 | Closure snapshots and journal line evidence. |
| QA-13 | Verify completion-linked telemetry and independently calculated cost/token aggregates `(qa-lead)` | VAL-021 | Pricing fixture, formula, run snapshot, stats output. |
| QA-14 | Verify maintenance before/after delete-and-rebuild structural equality `(qa-lead)` | VAL-022 | Deep-equality snapshots and equal audit output. |
| QA-15 | Verify journal trackability/immutability, DB ignore, and role restoration `(qa-lead)` | VAL-023 | Git commands, journal hashes, seeded role output. |
| QA-16 | Audit every v0.1-targeted PRD row to reproducible task/test/review evidence `(qa-lead)` | VAL-024 | Complete PRD traceability table with no unmapped row. |
| QA-17 | Run the exact phase gates, enforce COV-001/COV-002, and prove static binary <30 MiB `(qa-lead)` | VAL-024 | Literal command excerpts, package percentages, benchmark output, build/stat/test outputs. |
| QA-18 | Verify impacted docs and full seven-spec sync, including four maintenance formats and exact guard matrix `(qa-lead)` | VAL-012, VAL-016, VAL-017, VAL-024 | Diff/file-state review and rendered command examples. |
| QA-19 | Verify all ten completion reports are committed, clean, exact-command complete, and hash/subject valid `(qa-lead)` | VAL-024 | Ten report-lint PASS outputs and git evidence. |

## 8. Mid-Phase Review Cadence

The Heavy Mode cadence is mandatory and overrides the default every-three-task
rule:

1. **Entry Review immediately after Task 0 and before Task 1.** Write
   `project-docs/docs/phases/phase-17/qa/review-entry.md`; cover Task 0's
   tests-only progress toward VAL-024 without assigning a terminal VAL-024
   outcome. Run
   `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-0-done.md project-docs/docs/phases/phase-17/tasks/task-0.md`,
   independently rerun at least the focused journal coverage command and its
   report, verify the unchanged full-suite output also reports
   `internal/journal >=80.0%`, inspect the verified task commit file list for
   `*_test.go`-only scope, and assess whether every new test asserts observable
   behavior and can fail. Any production change, threshold/policy change,
   focused-only threshold evidence, result below 80.0%, invalid/missing command
   evidence, or weak line-coverage-only test is Blocking. Only a committed
   review with `VERDICT: CONTINUE` closes QA-BLK-01 and authorizes Task 1.
2. **Review 1 immediately after Task 2 and before Task 3.** Cover Tasks 1-2 and
   VAL-001..002. Run report-lint for:
   `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-1-done.md project-docs/docs/phases/phase-17/tasks/task-1.md`
   and
   `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-2-done.md project-docs/docs/phases/phase-17/tasks/task-2.md`.
   Recheck formatter line counts, byte-parity evidence, all config states, exact
   task commands, committed report hashes, and QA-BLK disposition.
3. **Review 2 immediately after Task 5 and before Task 6.** Cover Tasks 3-5 and
   VAL-003..014. Run:
   `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-3-done.md project-docs/docs/phases/phase-17/tasks/task-3.md`,
   `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-4-done.md project-docs/docs/phases/phase-17/tasks/task-4.md`, and
   `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-5-done.md project-docs/docs/phases/phase-17/tasks/task-5.md`.
   Re-run the integration gate, inspect test-only corruption isolation,
   both-state warning boundaries, all four formats, structured non-zero stdout,
   and journal immutability.
4. **Review 3 immediately after Task 7 and before Task 8.** Cover Tasks 6-7 and
   VAL-015..016. Run:
   `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-6-done.md project-docs/docs/phases/phase-17/tasks/task-6.md`
   and
   `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-7-done.md project-docs/docs/phases/phase-17/tasks/task-7.md`.
   Re-run the integration gate and verify lock, authorization/policy parity,
   exact guard precedence, and no-mutation proof.

After Task 9 and before final-review synthesis, QA Lead also runs:
`project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-8-done.md project-docs/docs/phases/phase-17/tasks/task-8.md`
and
`project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-9-done.md project-docs/docs/phases/phase-17/tasks/task-9.md`.

Each review uses `project-docs/templates/mid-phase-review-template.md`, is
committed by its QA Reviewer, and records `CONTINUE` or `BLOCK`. Missing,
narrower, silently skipped, or uncommitted task evidence is Blocking. A review
does not advance while an applicable blocking entry or in-window finding
remains unresolved. Review 1 cannot substitute for or absorb the Entry Review.

## 9. Human Verification Scope

Generate only these rows in `human-checklist.md` after deterministic QA passes:

| ID | Check | Covers | Human evidence expected |
|---|---|---|---|
| HUMAN-25a | Use maintenance during real project work and judge whether its findings explain what needs attention `(human)` | VAL-025 | Observation and pass/fail/partial notes. |
| HUMAN-25b | Invoke one governance stub and judge whether the version-gate error clearly says what ships when `(human)` | VAL-025 | Exact observed message plus clarity judgment. |
| HUMAN-25c | Run `fid ready` on real project state and judge whether the selected work is sensible `(human)` | VAL-025 | Selected ID/context plus judgment. |
| HUMAN-25d | In a disposable workspace copy, follow the documented maintenance/rebuild recovery path and judge whether the guidance is understandable `(human)` | VAL-025 | Recovery notes; no real-workspace DB deletion. |
| HUMAN-25e | After rows a-d, accept v0.1 for continued early dogfood or explicitly disposition remaining friction `(human)` | VAL-025 | Product-owner decision and follow-up reference if not accepted. |

## 10. Phase-End Decision Rule

The phase can proceed to final review only when:

- QA-BLK-02 remains closed by committed audit `2d5e536`; QA-BLK-01 is closed
  only by Task 0's focused/full >=80.0% proof and accepted Entry Review.
- Tasks 0-9 and their completion reports are committed and report-lint clean.
- Reviews after Tasks 0, 2, 5, and 7 are committed with no unresolved Blocking
  finding; `review-entry.md` must precede Task 1.
- Every required task and phase command passed, or the phase is stopped with a
  blocker report; substitutions are not accepted.
- COV-001 named domain packages are each >=80% and COV-002 `cmd/fid` is >=70%.
- QA-00..19 have reproducible evidence.
- VAL-001..024 each receive `PASS`, `FAIL`, `BLOCKED`, or `CARRY-FORWARD` in
  `final-review.md`; every PASS cites evidence.
- `human-checklist.md` contains only HUMAN-25a..e and VAL-025 is accepted or
  explicitly dispositioned before wrap.
- Full seven-spec sync is recorded, including correction of stale roadmap
  format wording.
- No unresolved blocking defect, spec deviation, or data-loss/security risk
  remains.

## 11. QA-Step To Validation Coverage Rollup

| VAL range | QA plan coverage | Mid/final owner |
|---|---|---|
| VAL-001 | QA-01, EVID-01 | Review 1 / QA Lead |
| VAL-002 | QA-02, EVID-02 | Review 1 / QA Lead |
| VAL-003..006 | QA-03, EVID-03..06 | Review 2 / QA Lead |
| VAL-007..010 | QA-04, EVID-07..10 | Review 2 / QA Lead |
| VAL-011 | QA-05, EVID-11 | Review 2 / QA Lead |
| VAL-012..014 | QA-06, EVID-12..14; QA-18 for docs | Review 2 / QA Lead |
| VAL-015 | QA-07, EVID-15 | Review 3 / QA Lead |
| VAL-016 | QA-08, EVID-16; QA-18 for docs | Review 3 / QA Lead |
| VAL-017 | QA-09, EVID-17; QA-18 for docs | QA Lead |
| VAL-018 | QA-10, EVID-18 | QA Lead |
| VAL-019 | QA-11, EVID-19 | QA Lead |
| VAL-020 | QA-12, EVID-20 | QA Lead |
| VAL-021 | QA-13, EVID-21 | QA Lead |
| VAL-022 | QA-14, EVID-22 | QA Lead |
| VAL-023 | QA-15, EVID-23 | QA Lead |
| VAL-024 | QA-00, QA-16..19, EVID-00/24/26 | Entry Review / QA Lead |
| VAL-025 | HUMAN-25a..e, EVID-25 | Human |

Coverage planning is complete. This plan is **Ready for Task 0 only** because
QA-BLK-02 is closed by commit `2d5e536` and DEC-054 authorizes the tests-only
remediation. QA-BLK-01 and Task 1 remain blocked until Task 0 proves the
unchanged 80% threshold in both required coverage runs and the committed Entry
Review accepts its tests, scope, commands, and evidence. No validation
assertion or hard stop is waived.
