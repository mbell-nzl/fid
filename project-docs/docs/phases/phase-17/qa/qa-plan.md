# Phase 17 QA Plan - Maintenance, Stub Guards, and v0.1 Closure

> Phase: 17 - Maintenance + stub guard + v0.1 closure
> Refreshed: 2026-07-11
> QA Lead: Codex
> Operating mode: Heavy (DEC-049)
> Status: **BLOCKED - not Ready**

## Blocking Entry Findings

The evidence design, commands, fixtures, artifacts, and human checks are
feasible, but Developer execution must not start until both prerequisites below
have durable resolution evidence.

| ID | Blocking finding | Evidence | Required resolution |
|---|---|---|---|
| QA-BLK-01 | The fresh COV-001 baseline is below the required threshold: `internal/journal` is 78.9%, while `08_QUALITY_PLAN.md` COV-001 requires each named domain package to remain at or above 80%. Tasks 1-9 do not authorize a journal-coverage remediation; Task 9 is explicitly audit-only. VAL-024 therefore has no presently authorized path to PASS. | Fresh `go test -race -coverprofile=coverage.out ./...` on 2026-07-11: `internal/journal ... coverage: 78.9%`; Phase 16 final review records the same 78.9%. Earlier DEC-007 and Phase 1/10/12 QA evidence establish package-level enforcement. | Tech Lead adds a scoped remediation task or explicitly expands a task spec to restore `internal/journal >= 80%`, then QA reruns the exact coverage commands. A quality-policy waiver would require a recorded Tech Lead/Human decision; QA will not infer one from Phase 16's contradictory PASS. |
| QA-BLK-02 | DEC-053 says Phase 17 must not trust the new QA Lead/QA Reviewer model routing until model-audit evidence confirms seeded-defect detection and hard-stop behavior. No model-audit artifact or evidence reference exists in the required context or Phase 17 run ledger. | `project-docs/docs/phases/phase-17/decisions.md` DEC-053; repository search found no model-audit evidence beyond the decision itself. `run-ledger.md` contains only `RUN-START`. | Tech Lead/Human records and cites the required model-audit evidence, or records an approved proven-model routing/fallback. Start fresh role sessions after any routing change, as DEC-053 requires. |

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
- `project-docs/docs/phases/phase-17/tasks/task-1.md` through `task-9.md`
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
- Every task cites its owned assertions: YES.
- Every VAL-001..025 has a QA-plan verification step: YES; see section 7 and
  the coverage rollup in section 11.
- One-line rollup: **all in-scope requirements are referenced by at least one
  VAL assertion: YES.**
- Revisions requested from Tech Lead: none to assertion wording or ownership.
  QA-BLK-01 is an execution-scope prerequisite, not a proposal to weaken
  VAL-024. QA-BLK-02 is a DEC-053 evidence prerequisite.

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA notes |
|---|---|---|---|---|
| 1 | Behavior-preserving formatter split | VAL-001 | Medium | Review byte parity and every resulting non-test file <=400 lines. |
| 2 | Maintenance grace config | VAL-002 | High | Review legacy/missing, zero, positive, negative, malformed, and bootstrap round-trip states. |
| 3 | Structural maintenance errors | VAL-003..006 | High | Trigger/FK bypass must be test-only; valid counterpart required for every violation family. |
| 4 | Warnings and info | VAL-007..011 | High | Inject time; prove boundary and excluded states; deep-equality ordering. |
| 5 | Maintenance check CLI/formats | VAL-012..014 | High | Structured stdout must survive expected non-zero outcomes; journal stays byte-identical. |
| 6 | Locked VACUUM | VAL-015 | High | Authorization before mutation, one lock, data/journal preservation, policy parity. |
| 7 | Governance guards | VAL-016 | High | Exact command matrix and precedence; no validation or persistence access before gate. |
| 8 | Read-format conformance | VAL-017 | High | Exact 15-command matrix; structural machine-format parity and markdown identities. |
| 9 | v0.1 closure journey | VAL-018..023 | High | Audit-only; any product defect stops for triage. Temp git and delete/rebuild are mandatory. |
| Phase end | Traceability, full gates, size | VAL-024 | High | Every v0.1 PRD row maps to reproducible evidence; COV-001 blocker must be closed. |
| Human | Real-work dogfood | VAL-025 | High | Run only after deterministic QA; one judgment per HUMAN row. |

## 4. Required Evidence Strategy

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
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
| EVID-24 | Every v0.1 PRD row mapped to task/test/review evidence; exact phase gate, COV-001/COV-002 values, static size | VAL-024 | QA Lead | Final review |
| EVID-25 | Five separately recorded real-work judgments | VAL-025 | Human | Human checklist |
| EVID-26 | Exact task command coverage and valid commit hash/subject from committed completion reports | no VAL | QA Reviewer, QA Lead | Each review / final review |

## 5. Deterministic Gates

No required command contains a placeholder. Every task report must record every
command from its task spec verbatim as `PASS`, or as `BLOCKED` in a blocker
report. Narrower substitutions do not count.

| Gate | Command | Expected result | Required when |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | No files listed | Every task and phase end |
| Lint | `golangci-lint run ./...` | `0 issues.` | Every task and phase end |
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
| Completion report lint | Nine exact `project-docs/tools/report-lint.sh` commands listed in section 8 | Exit 0 for every task report/spec pair | Each review and final |

Security review at phase end covers SQL construction/input binding,
test-only corruption bypass isolation, guard-before-input behavior, lock
handling, path handling, and journal/data-loss risk. Public-release scanners
remain deferred by DEC-049.

## 6. Command And Feasibility Validation

Fresh validation was performed on 2026-07-11; no prior-phase command result was
carried forward.

### 6.1 Command Feasibility

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | Exact command | PASS | Exit 0; no files listed. `gofmt -h` confirms `-l` and path arguments. |
| `golangci-lint run ./...` | Exact command and help | PASS | Exit 0, `0 issues.`; `golangci-lint run --help` confirms syntax. |
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
| Nine exact report-lint commands in section 8 | Script inspection and `project-docs/tools/report-lint.sh --self-test` | PASS | Self-test passed clean/skipped/narrower/fabricated-hash/dirty-report cases; all nine task/report paths are declared by task specs. |
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
| Task evidence inputs | File inventory | PASS | Nine task specs exist; each has exact commands and a completion-report path. No completion report is expected before its task runs. |
| Mid-phase reviews | Template/path validation | PASS | `mid-phase-review-template.md` exists; `review-1.md`, `review-2.md`, and `review-3.md` are feasible at the checkpoints in section 8. |
| Final review | Template/path validation | PASS | `final-review-template.md` includes gate results, every VAL outcome, deterministic execution, and maturity disposition. |
| Human checklist | Template/path validation | PASS | `human-checklist-template.md` supports one observable per row and human-only outcomes. |
| Completion report verification | Template/tool validation | PASS | Completion template and executable report-lint tool exist; self-test is green. |
| Commit evidence | Git/tool validation | PASS | Git 2.53.0 is available; report-lint verifies hash, subject, committed state, and exact command coverage. |
| QA artifact ownership | Repository/write check | PASS | This QA plan can be committed independently with `docs(qa): ...` and `Refs: phase-17`; unrelated `run-ledger.md` remains untouched. |
| DEC-053 model-audit evidence | Repository search and run-ledger inspection | BLOCKED | No evidence artifact/reference found; see QA-BLK-02. |

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
| QA-19 | Verify all nine completion reports are committed, clean, exact-command complete, and hash/subject valid `(qa-lead)` | VAL-024 | Nine report-lint PASS outputs and git evidence. |

## 8. Mid-Phase Review Cadence

The Heavy Mode cadence is mandatory and overrides the default every-three-task
rule:

1. **Review 1 immediately after Task 2 and before Task 3.** Cover Tasks 1-2 and
   VAL-001..002. Run report-lint for:
   `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-1-done.md project-docs/docs/phases/phase-17/tasks/task-1.md`
   and
   `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-2-done.md project-docs/docs/phases/phase-17/tasks/task-2.md`.
   Recheck formatter line counts, byte-parity evidence, all config states, exact
   task commands, committed report hashes, and QA-BLK disposition.
2. **Review 2 immediately after Task 5 and before Task 6.** Cover Tasks 3-5 and
   VAL-003..014. Run:
   `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-3-done.md project-docs/docs/phases/phase-17/tasks/task-3.md`,
   `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-4-done.md project-docs/docs/phases/phase-17/tasks/task-4.md`, and
   `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-17/tasks/task-5-done.md project-docs/docs/phases/phase-17/tasks/task-5.md`.
   Re-run the integration gate, inspect test-only corruption isolation,
   both-state warning boundaries, all four formats, structured non-zero stdout,
   and journal immutability.
3. **Review 3 immediately after Task 7 and before Task 8.** Cover Tasks 6-7 and
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
does not advance while either QA-BLK prerequisite remains unresolved.

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

- QA-BLK-01 and QA-BLK-02 have durable resolution evidence.
- Tasks 1-9 and their completion reports are committed and report-lint clean.
- Reviews after Tasks 2, 5, and 7 are committed with no unresolved Blocking
  finding.
- Every required task and phase command passed, or the phase is stopped with a
  blocker report; substitutions are not accepted.
- COV-001 named domain packages are each >=80% and COV-002 `cmd/fid` is >=70%.
- QA-01..19 have reproducible evidence.
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
| VAL-024 | QA-16..19, EVID-24/26 | QA Lead |
| VAL-025 | HUMAN-25a..e, EVID-25 | Human |

Coverage is complete, but this plan remains **BLOCKED** until QA-BLK-01 and
QA-BLK-02 are resolved. QA Lead must refresh the status and cite both
resolutions before the Phase Runner starts Task 1.
