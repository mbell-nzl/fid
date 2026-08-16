# Phase 20 QA Plan - Clean Standalone Issue Journey

> Phase: 20 - Clean standalone Issue journey
> Refreshed: 2026-08-09
> QA Lead: Codex QA Lead
> Operating mode: Heavy
> Readiness: **READY**

## 1. Scope And Source Documents

Phase 20 establishes the `fid-cli.v1` Response Envelope and one coherent
standalone Issue journey: initialize, create, relate, inspect, update,
transition, explain a typed failure, rebuild, and prove unchanged Issue state.
It also removes the universal formatter surface, YAML dependency, retired
Issue fields, and promised-command stubs while preserving fail-closed V1
authority until Phase 21.

This is a fresh heavy-mode QA pass. No command or feasibility result is carried
forward. No phase outcome is pre-awarded: all `VAL-001..VAL-046` rows remain
subject to evidence and a phase-end outcome.

Sources read:

- `AGENTS.md`.
- `project-docs/docs/current/phase_brief.md`.
- `project-docs/docs/current/decisions_active.md`.
- `project-docs/docs/current/progress_current.md`.
- `project-docs/docs/current/tl_escalations_current.md`.
- `project-docs/specs/08_QUALITY_PLAN.md`.
- `project-docs/specs/09_AI_WORKFLOW.md`, especially sections 2.6, 4.3,
  6.1, and 8.
- `project-docs/docs/phases/phase-20/qa/validation-contract.md`.
- `project-docs/docs/phases/phase-20/contracts.md`.
- `project-docs/docs/phases/phase-20/tasks/task-1.md` through `task-18.md`.
- `project-docs/templates/qa-plan-template.md`.
- Scoped requirement excerpts named by the contract: `01_PRD.md` sections
  1-4 and 9-10; `02_TRD.md` sections 4 and 6; `03_TRD_SCHEMA.md` sections
  7-8.1; `04_TRD_CONFIG.md` sections 2 and 5; `05_DESIGN_PLAN.md` section 2
  rule 6 and section 6; `06_DESIGN_PLAN_CLI.md` sections 1-2, 5, and 8;
  `10_IMPLEMENTATION_ROADMAP.md` sections 5-7; and `GLOSSARY.md` Machine
  Interface terms.

Out of scope for this plan are product implementation, tests, task-spec or
decision edits, validation-contract edits, phase-runner ledger entries,
held-out scenario material, and Phase 21 capability/Work Envelope work.

## 2. Validation Contract Refinement Check

- All in-scope PRD/TRD requirements referenced by at least one `VAL-NNN`:
  **YES**. The contract also explicitly dispositions requirements deferred to
  Phases 21-25.
- Every `VAL-NNN` has one validator role: **YES**. Inventory: 37 QA Reviewer,
  8 QA Lead, and 1 Human row.
- Every `VAL-NNN` has a concrete evidence method: **YES**.
- Every `VAL-NNN` maps to a task/window or phase-end check: **YES**.
- Every `VAL-NNN` is cited by at least one task-spec acceptance criterion:
  **YES**; the acceptance-criteria inventory is exactly `VAL-001..VAL-046`
  with no gap.
- Every `VAL-NNN` describes one observable or an explicitly required Human
  fan-out: **YES**. Cross-family rows such as VAL-010, VAL-011, VAL-035,
  VAL-044, and VAL-045 are finite structural-conformance invariants over
  literal registries, not open-ended bundles. VAL-038 explicitly requires
  Human sub-IDs.
- Conditional or binary assertions require qualifying and non-qualifying
  evidence: **YES**. The evidence plan names both states for discovery,
  authority, lock release, read/mutation Event identity, ready results,
  hierarchy, concurrency, transitions, cycles, warning promotion, retained vs
  legacy transport, public-owner absence vs internal safety retention, atomic
  initialization, request validity, limits, and typed endpoint resolution.
- Parity or reference-matching assertions have deterministic structural proof:
  **YES**. Exact decoded key sets, literal leaf/request registries, golden
  Event tuples, complete table dumps, identical warning findings, and
  byte-level pre/post-rebuild Issue JSON are required where applicable.
- Revisions requested from Tech Lead: **none**.
- Escalations: **none**.

VAL-038 is fanned out as `HUMAN-038a..c`. Its first two checks are perception
and product-judgment checks. The third asks the Human to judge semantic advice,
including indirect wording; deterministic output/source assertions that catch
literal `.fid/` guidance remain supporting `(qa-lead)` evidence and do not
replace the Human judgment.

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | Response contract foundation | VAL-001, VAL-002 | High | Exact success/failure key sets, closed reason registry, channels, redaction; no domain/Cobra coupling. |
| Task 2 | Project selection and strict config | VAL-003, VAL-005, VAL-026, VAL-027, VAL-039, VAL-043, VAL-044 | High | Widest task; review one config model, nearest-boundary discovery, one root-flag registry, fixed maintenance grace, and workflow flag collision. |
| Task 3 | Atomic project initialization | VAL-001, VAL-002, VAL-004, VAL-008, VAL-017, VAL-039, VAL-042, VAL-043, VAL-044 | High | Inject every named pre-activation failure; prove no target/staging residue and verified rebuild before one activation. |
| Task 4 | Rebuild, authority, and locks | VAL-001, VAL-002, VAL-006, VAL-007, VAL-008, VAL-017, VAL-044 | High | Both denied and permitted authority states; held and released lock; no mutation on rejection. |
| Task 5 | Retired Issue field removal | VAL-009, VAL-012 | High | Migration/table/view/fixture/replay churn; removed keys must fail loudly, not be ignored. |
| Task 6 | Replayable aggregate revision | VAL-010, VAL-011, VAL-012 | High | Six owned Event families, legacy dev safety, exact prior/precondition tuple, and live/rebuild parity. |
| Task 7 | Hierarchy invariant | VAL-020 | Medium | Legal, illegal, missing-parent, and absent parent-mutation surface all required. |
| Task 8 | Issue read contract | VAL-008, VAL-013, VAL-014, VAL-015, VAL-016, VAL-017, VAL-044, VAL-045 | High | Exact nullability/key sets, both ready-result states, all filters/tie-breaks, and bounded collections. |
| Task 9 | Issue mutation/concurrency | VAL-002, VAL-008, VAL-010, VAL-011, VAL-017, VAL-018, VAL-019, VAL-040, VAL-041, VAL-044, VAL-045 | High | Absent/matching/stale revisions; all transition states; owner removal with retained safety; attribution without trust. |
| Task 10 | Acceptance criteria and labels | VAL-008, VAL-010, VAL-011, VAL-017, VAL-021, VAL-022, VAL-041, VAL-044, VAL-045 | High | Six operations, canonical `colour`, exact revisions/Events, typed duplicate/missing behavior, order/bounds. |
| Task 11 | Links and dependency invariants | VAL-008, VAL-010, VAL-011, VAL-017, VAL-023, VAL-024, VAL-025, VAL-041, VAL-044, VAL-045, VAL-046 | High | Explicit typed endpoints, opaque IDs, inverse projection, self/2/3-node cycles, acyclic diamond, bounded termination. |
| Task 12 | Workflow command contract | VAL-008, VAL-017, VAL-026, VAL-041, VAL-044, VAL-045 | High | Closed grammar, no `--require`, project/kind precedence, no root-flag shadow, retained fail-closed legacy preconditions. |
| Task 13 | Maintenance inspection | VAL-002, VAL-008, VAL-017, VAL-027, VAL-044, VAL-045 | High | Closed registry, exact 299/300 boundary, identical warn findings across success/failure branches, no Attempt-config alias. |
| Task 14 | Legacy JSON conversion | VAL-028, VAL-029, VAL-031 | High | Eight unenveloped semantic JSON reads; journal inspect remains legacy and has no revert/Envelope promise. |
| Task 15 | Remaining format callers/tests | VAL-029, VAL-031 | Medium | Preserve semantic assertions while removing all live/test `--format` callers. |
| Task 16 | Formatter/YAML deletion | VAL-030, VAL-031, VAL-032, VAL-033, VAL-037 | High | Sharp deletion, module tidy, no-match evidence, full gate, architecture-doc update. |
| Task 17 | Promised-command stub deletion | VAL-034 | Medium | Unknown-command behavior, source absence, and retained maintenance tables. |
| Task 18 | Journey acceptance/docs | VAL-008, VAL-017, VAL-035, VAL-036, VAL-037, VAL-038, VAL-041, VAL-044, VAL-045 | High | Acceptance evidence only; exact 21/21 inventory, complete JSON journey, byte-identical rebuild, docs, and Human-check support. Product mismatches block. |

## 4. Required Evidence Strategy

All automated evidence is produced from committed product/test artifacts. A
completion report attestation alone is not a pass; the relevant review or
phase-end session must cite reproducible output/file state. Rejected operations
must include journal/projection state proving zero effect where the contract
requires atomicity.

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | Decoded success/failure key-set matrices plus CLI stdout/stderr/exit captures, redaction cases, and journal-tail Event comparison for every retained section-3 leaf; both Human success and failure channels. | VAL-001, VAL-002, VAL-008, VAL-017 | Developer; QA Reviewer; QA Lead | Tasks 1, 3-13; reviews 1-5; Task 18; final review |
| EVID-02 | Temporary root/child/boundary cases and literal Cobra sentinel matrix proving nearest discovery, exact explicit selection, canonical-leaf acceptance, and legacy/not-yet-migrated rejection before handler execution. | VAL-003, VAL-004, VAL-005 | Developer; QA Reviewer | Tasks 2-3; review 1 |
| EVID-03 | Permitted/missing/unknown/denied authority table with journal/projection counts, and deterministic held-lock then released-lock retry evidence. | VAL-006, VAL-007 | Developer; QA Reviewer | Task 4; review 2 |
| EVID-04 | PRAGMA/decoded output/Event checks for all retired fields; migration up/down evidence; removed-key replay failure; complete deterministic live vs rebuilt table dumps. | VAL-009, VAL-012 | Developer; QA Reviewer; QA Lead | Tasks 5-6; review 2; final review |
| EVID-05 | Revision sequence across all eleven Event families, sorted one/two-Issue link revisions, rejection invariance, exact golden Values/Prior/Preconditions, stale replay failure, and identical rebuild. | VAL-010, VAL-011 | Developer; QA Reviewer | Tasks 6, 9-11; reviews 2-4 |
| EVID-06 | Fully populated/all-null read schemas, equal-key filter/order fixtures, both exclusive ready states and both explanation codes, retired-flag rejection, legal/illegal/missing hierarchy, and no parent-mutation request. | VAL-013, VAL-014, VAL-015, VAL-016, VAL-020 | Developer; QA Reviewer | Tasks 7-8; review 3 |
| EVID-07 | Absent/matching/stale expected-revision cases; legal/unknown-edge/absent-edge/unsupported-or-unsatisfied transition cases; public owner-input/result absence plus retained safety regression; attribution Event and forbidden trusted/protected effect checks. | VAL-018, VAL-019, VAL-040, VAL-041 | Developer; QA Reviewer; QA Lead | Tasks 9-12; reviews 3-4; Task 18; final review |
| EVID-08 | Exact AC/label key/type/revision/order matrices and typed duplicate/missing atomic failures, including canonical `colour`. | VAL-021, VAL-022 | Developer; QA Reviewer | Task 10; review 4 |
| EVID-09 | Exact typed-Link add/list/remove sequence, inverse projection, specific negative codes, zero append, self/2/3-node cycles, acyclic diamond, pre-existing-cycle termination, ambiguous-prefix and typed missing-endpoint matrix. | VAL-023, VAL-024, VAL-025, VAL-046 | Developer; QA Reviewer | Task 11; review 4 |
| EVID-10 | Workflow request/result/grammar matrix; valid and invalid values; `--require` and requirement-count absence; project/kind assignment states and precedence; atomic typed failures. | VAL-026 | Developer; QA Reviewer | Tasks 2 and 12; reviews 1 and 4 |
| EVID-11 | Exact maintenance check-name registry; config-independence and 299/300-second cases; same fixture decoded in warn success and promoted failure with identical findings/info; channel/exit evidence. | VAL-027 | Developer; QA Reviewer | Tasks 2 and 13; reviews 1 and 5 |
| EVID-12 | Decoded semantic JSON for all eight legacy reads; journal-inspect source/command-tree absence of revert and Envelope; explicit `--json` rejection before handler. | VAL-028, VAL-029 | Developer; QA Reviewer | Tasks 14-15; review 5 |
| EVID-13 | Directory/module/source/Cobra no-match or expected-failure checks plus full test/build evidence for formatters, flags, YAML, serializers, and promised stubs. | VAL-030, VAL-031, VAL-032, VAL-033, VAL-034 | Developer; QA Reviewer; QA Lead | Tasks 14-17; reviews 5-6; final review |
| EVID-14 | Literal 42-leaf disjoint partition test, all 21 closed request/domain matrices with valid and invalid handler sentinels, all collection/change ordering and at/over-limit fixtures, complete integration journey, and byte-identical pre/post-rebuild Issue JSON. | VAL-035, VAL-036, VAL-044, VAL-045 | Developer; QA Lead | Task 18; final review |
| EVID-15 | Live-tree comparison with architecture documentation, playbook, and changelog; no current formatter/stub contradiction. | VAL-037 | Developer; QA Lead | Tasks 16 and 18; review 6; final review |
| EVID-16 | Strict config round-trip and key inventory; generated opaque identity; prefix/path/signed-64-bit just-inside/outside/malformed/symlink cases; injected init failures and successful rebuild/activation proof. | VAL-039, VAL-042, VAL-043 | Developer; QA Reviewer | Tasks 2-3; review 1 |
| EVID-17 | Human directly judges concise routine output, actionable failure guidance, and absence of advice to edit fid internals, one observation per sub-ID. | VAL-038 | Human | Phase-end human checklist |

## 5. Deterministic Gates

Unless stated otherwise, commands run from repository root
`/home/mark/projects/fid`. Task specs remain authoritative for per-task
commands. Commands added here for Task 16/17 absence proof are also required
and must appear in those completion reports with their real exit behavior.

**Absence-command notation (DEC-143).** A search or inventory command that
proves an absence terminates non-zero. Its owning task spec declares the
expected status under `## Expected Non-Zero Commands`, and the report records
the real observed status in the declared form
``- `<cmd>` - PASS - `<token>` (expected `exit_code: N`; terminated `exit_code: M`)``.
`go.sum` checksum lines retained for the module graph of an accepted retained
dependency are module-graph bookkeeping, not fid dependencies, so no gate here
searches `go.sum`. The withdrawn row `rg -n "gopkg.in/yaml.v3" go.mod go.sum`
is superseded by the two YAML rows below, which prove requirement-level and
build-graph absence instead of textual absence.

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | Exit 0; no files listed | Every task; phase end |
| Lint/static analysis | `golangci-lint run ./...` | Exit 0 | Every task; phase end |
| Race/coverage | `go test -race -coverprofile=coverage.out ./...` | Exit 0 and current-session `coverage.out` | Every task; phase end |
| Coverage report | `go tool cover -func=coverage.out` | Exit 0; per-function and total report from the immediately preceding profile | Every task; phase end |
| Integration | `go test -tags integration -race ./cmd/fid/integration/...` | Exit 0 | Tasks 4-18 after DEC-133; phase end |
| Accepted journey | `go test -tags integration -race ./cmd/fid/integration/...` | Task 18's committed Phase 20 journey passes inside the dedicated real-process suite | Task 18; phase end; no separate E2E placeholder |
| Benchmarks | `go test -bench . -benchmem ./...` | Exit 0 with benchmark timing/allocation output | Tasks 5, 6, and 18; phase end |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | Exit 0 and executable produced | Every task; phase end |
| Module convergence | `go mod tidy` | Exit 0; only intended `go.mod`/`go.sum` changes | Task 16 |
| Formatter directory absence | `ls internal/format` | Expected failure: non-zero with directory absent | Task 16; review 6; phase end |
| Formatter import absence | `rg -n "internal/format" cmd internal` | Expected no-match: exit 1, no output | Task 16; review 6; phase end |
| Serializer absence | `rg -ni "writeTOON|writeYAML|writeMarkdown|yaml.Marshal" cmd internal` | Expected no-match: exit 1, no output | Task 16; review 6; phase end |
| YAML module requirement absence | `rg -ni "yaml" go.mod` | Expected no-match: exit 1, no output | Task 16; review 6; phase end |
| YAML build-graph absence | `go mod why -vendor -m gopkg.in/yaml.v3` | Exit 0 whose whole output is one parenthesized `does not need` line naming the module - on this toolchain `(main module does not need to vendor module gopkg.in/yaml.v3)`, with the non-vendor wording equally lawful. An import chain, an extra line, or a non-zero exit fails (DEC-144) | Task 16; review 6; phase end |
| Stub/promise absence | `rg -n "newDecisionCommand|newEscalationCommand|newTriageCommand|versionGateError|ships in v0\\.2" cmd/fid` | Expected no-match: exit 1, no output | Task 17; review 6; phase end |
| Full phase gate | The seven-command block below | Every command terminates and passes with its own current-session token | Phase end |

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

There is no accepted numeric coverage threshold; coverage is a diagnostic and
must not substitute for the scenario matrices above. Security SAST,
dependency-vulnerability, license, committed-secret, and held-out scenario
commands remain `N/A` before release readiness under `AGENTS.md`; Phase 20
still receives the deterministic path/redaction/atomicity review in QA-21.

## 6. Command And Feasibility Validation

Fresh validation was performed on 2026-08-09. `PASS` here means the command or
future step is concrete and executable as planned; it is not evidence that the
not-yet-implemented Phase 20 behavior passes. Current positive matches for
removal inventories are the expected pre-Task-16/17 state; the finished-work
pass condition remains real exit 1 with no output.

| Referenced Step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | Executed exactly from repository root | PASS | Exit 0, no output; both target directories exist. |
| `golangci-lint run ./...` | `golangci-lint run --help`; version/config inventory | PASS | Help exit 0 with `Usage: golangci-lint run [flags]`; version 2.12.2; `.golangci.yml` exists. |
| `go test -race -coverprofile=coverage.out ./...` | `go help test`, `go help testflag`, platform/cache inspection | PASS | Go 1.26.2 on linux/amd64; help documents `-race` and `-coverprofile`; module is repository `go.mod`; writable repo-local caches are configured. |
| `go tool cover -func=coverage.out` | `go tool cover -help`; producer/consumer sequence inspection | PASS | Help printed `go tool cover -func=c.out`; plan requires the profile from the immediately preceding successful test command, preventing stale-profile evidence. |
| `go test -tags integration -race ./cmd/fid/integration/...` | `go help build` plus committed classification regression | PASS | Help documents `-tags`; `cmd/fid/integration/classification_test.go` requires exactly four current real-process journeys and no integration-tagged parent-package tests. Task 18 adds its accepted journey to this dedicated scope. |
| `go test -bench . -benchmem ./...` | `go help testflag` plus benchmark inventory | PASS | Help documents `-bench .` and `-benchmem`; four existing benchmarks cover rebuild and Issue registry reads. |
| `go build -o /tmp/fid ./cmd/fid` | `go help build`, package/file state, target state | PASS | Help documents `-o`; `cmd/fid/main.go` contains `package main`/`func main`; `/tmp` is writable and the existing `/tmp/fid` is owned by the current user. Command exit remains the gate, never stale binary existence. |
| `go mod tidy` | `go help mod tidy`; module-file state | PASS | Help exit 0 with exact usage; `go.mod` and `go.sum` exist; Task 16 owns their intended dependency-removal diff. |
| `ls internal/format` | Executed exactly against current tree | PASS | Current exit 0 and directory listing prove a real existing deletion target; finished-work expectation is non-zero because Task 16 deletes it. |
| `rg -n "internal/format" cmd internal` | Executed exactly against current tree | PASS | Current exit 0 listed 16 live imports, proving concrete targets; finished-work expectation is exit 1/no output. |
| `rg -ni "writeTOON|writeYAML|writeMarkdown|yaml.Marshal" cmd internal` | Executed exactly against current tree | PASS | Current exit 0 listed serializer definitions/calls/tests; finished-work expectation is exit 1/no output. |
| ~~`rg -n "gopkg.in/yaml.v3" go.mod go.sum`~~ | Withdrawn 2026-08-16 by DEC-143 | WITHDRAWN | Planning validated only that the pattern currently matched. It could never reach exit 1: `go mod tidy` deterministically restores the `go.sum` checksum required by the retained `github.com/golang-migrate/migrate/v4@v4.19.1`. Superseded by the two rows below. |
| `rg -ni "yaml" go.mod` | Read the committed `go.mod` after Task 16's tidy | PASS | `go.mod` contains no yaml requirement, direct or indirect, so the finished-work expectation exit 1/no output is reachable; the broader pattern also rejects any other yaml module. |
| `go mod why -vendor -m gopkg.in/yaml.v3` | Observed run of the exact command in the Task 16 continuation (DEC-144); previously `go help mod why` semantics only | PASS | `-vendor` excludes tests of dependencies, which is the asserted boundary; a module outside that graph exits 0 and prints only `(main module does not need to vendor module <path>)` under `-vendor` (`(main module does not need module <path>)` without it). An import chain instead is a real VAL-032 failure. The original row inferred the non-vendor wording from documentation and was corrected by DEC-144 after the tool was observed. |
| Stub/promise absence search | Executed the exact planned `rg` against `cmd/fid` | PASS | Current exit 0 listed `versionGateError`, three constructors, and root registrations; Task 17 owns removal and future exit 1/no output. |
| Task-spec command completeness/placeholders | Parsed all 18 `## Exact Commands To Run` blocks and searched phase contract/task command text | PASS | Every command is concrete; no required command placeholder. Angle-bracket tokens found only in contract JSON examples, not command blocks. |
| VAL/task acceptance coverage | Compared contract IDs with IDs inside every task's `## Acceptance Criteria` | PASS | Both inventories are exactly sequential `VAL-001..VAL-046`. |
| Validator ownership | Parsed all 46 assertion rows | PASS | Every row has exactly one of `QA Reviewer`, `QA Lead`, or `Human`; no blank/multiple validator field. |
| Planned package/directories | File-state inventory of `cmd`, `internal`, `migrations`, module files, and task scope | PASS | Existing owning packages/files are present. `internal/cli`, `internal/project`, migrations 0011/0012, and Task 18 acceptance files are explicitly create-scoped deliverables with existing writable parents. |
| Temporary repository/fixture setup | Existing test-pattern inventory and filesystem check | PASS | Existing suites use `testing.T.TempDir()` extensively; `/tmp` is writable; integration-tag and temporary SQLite/journal patterns already exist. No fixed project fixture or cleanup-dependent repository is required. |
| Atomic-init fault injection | Accepted contract/task inspection plus adapter/test patterns | PASS | Task 3 explicitly owns narrow injectable stage/filesystem hooks, one synchronous case per named failure, unique staging ownership, target/residue assertions, and success rebuild comparison. |
| Exact structural conformance | Contract/task inspection | PASS | Closed key schemas, reason/detail registries, literal 21/21 leaf lists, request/value tables, Event families, ordering rules, and limits are all finite machine-checkable references. |
| Rebuild parity | Contract/task inspection | PASS | VAL-012 uses complete table dumps; VAL-036 uses canonical `issue.get` byte comparison. Both specify the live and rebuilt operands and do not rely on visual/code-reading parity. |
| Task 18 journey artifact | Parent/file convention inspection | PASS | `cmd/fid/` exists; integration build tags and `T.TempDir` are established; task spec names `cmd/fid/e2e_p20_issue_journey_test.go` and the exact integration command. |
| Documentation artifacts | File existence check | PASS | `AGENTS.md`, `project-docs/docs/playbook.md`, and `CHANGELOG.md` exist and Tasks 16/18 name ownership. Historical archives are excluded from current-claim searches. |
| Human concise-output judgment | Accepted Human-channel contract and Task 18 setup review | PASS | Sections 2.2/3.2 define a runnable Human-mode surface; final QA will provide a fresh-project interaction sequence and captured output from committed behavior for direct Human judgment. |
| Human actionable-guidance judgment | Accepted typed-failure and Human error contract review | PASS | Task 18 must provoke a specified typed failure; the same state is runnable without `--json`, allowing the Human to judge next-action guidance. |
| Human no-internal-edit judgment | Design rule and output-corpus feasibility review | PASS | `05_DESIGN_PLAN.md` section 6 supplies the rule; final QA supplies the Human failure corpus and deterministic literal/internal-path support evidence for semantic Human review. |
| Six review windows and phase-end order | Phase brief, task dependencies, and workflow inspection | PASS | Concrete windows close after Tasks 3, 6, 9, 12, 15, and 17; Task 18 follows review 6; final QA follows verified Task 18. |

No command or step is invalid, unavailable, or dependent on an unaccepted
fixture. No configuration defect or escalation was found.

## 7. Deterministic Checklist

At phase end, the QA Lead runs or independently verifies every row. A broad
green suite is paired with the named test/file-state evidence; code reading
alone is not a pass. Reviewer-owned assertions retain their contract validator
but are rechecked for complete evidence and assigned final outcomes by QA Lead.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Verify every task/review required command is present with exact command text, real terminal status/token, and no skipped or narrower substitute; run the full phase gate `(qa-lead)` | no VAL | Completion reports, review logs, report-lint records, and current-session full-gate output. |
| QA-02 | Verify exact success/failure envelopes, channel/exit separation, redaction, and mutation/read Event identity across all retained leaves `(qa-lead)` | VAL-001, VAL-002, VAL-008, VAL-017 | Decoded key-set/code tests, CLI captures, journal-tail comparisons, full/integration test output. |
| QA-03 | Verify discovery/explicit-init boundaries and both global flags across retained and legacy/not-yet-migrated handlers `(qa-lead)` | VAL-003, VAL-004, VAL-005 | Temporary-tree CLI tests, literal Cobra sentinel table, review-1 reproduced evidence. |
| QA-04 | Verify permitted and all rejected V1 authority states plus held/released lock behavior without rejected effects `(qa-lead)` | VAL-006, VAL-007 | State/journal counts, typed code/detail assertions, deterministic lock test output. |
| QA-05 | Verify retired field absence/replay rejection and complete live/rebuild table equivalence `(qa-lead)` | VAL-009, VAL-012 | PRAGMA/CLI/Event tests, migration up/down results, complete deterministic table dumps. |
| QA-06 | Verify all eleven Issue aggregate mutation families, exact revision tuples, rejection invariance, sorted affected revisions, stale replay failure, and rebuild `(qa-lead)` | VAL-010, VAL-011 | Named golden fixtures/tests and current test output covering every family. |
| QA-07 | Verify closed Issue read schemas, filters/tie-breaks, both ready states/codes, retired flags, and every legal/illegal/missing hierarchy state `(qa-lead)` | VAL-013, VAL-014, VAL-015, VAL-016, VAL-020 | Decoded schema and table-driven behavior output plus request-tree absence proof. |
| QA-08 | Verify absent/matching/stale concurrency, all transition outcomes, public owner removal with retained safety, and attribution without trusted/protected side effects `(qa-lead)` | VAL-018, VAL-019, VAL-040, VAL-041 | Journal/state counts, exact error details, flag/result guards, forbidden table/Event-effect assertions. |
| QA-09 | Verify exact AC and label contracts, canonical vocabulary, revisions/order, and atomic missing/duplicate behavior `(qa-lead)` | VAL-021, VAL-022 | Key/type/value-domain/order tests and mutation/rejection Event counts. |
| QA-10 | Verify typed Link schemas/directions, specific negatives, bounded cycle cases plus diamond, and opaque typed endpoint resolution `(qa-lead)` | VAL-023, VAL-024, VAL-025, VAL-046 | Add/list/remove sequence, graph table, ambiguous-prefix/missing-endpoint tests, journal counts. |
| QA-11 | Verify exact workflow schemas/grammar/assignment precedence and absence of V1 rule authoring or flag collision `(qa-lead)` | VAL-026 | Request/result/value-domain, forbidden-flag/result, precedence, and typed-failure output. |
| QA-12 | Verify closed maintenance registry, fixed 299/300 behavior, and identical warn findings/info across success and promoted failure `(qa-lead)` | VAL-027 | Config-independent clock fixtures, decoded branch comparison, channel/exit output. |
| QA-13 | Verify all eight legacy reads retain semantic valid JSON and journal inspect remains unenveloped, no-revert, and rejects `--json` `(qa-lead)` | VAL-028, VAL-029 | Decoded outputs/semantic fields, command-tree sentinel, source-absence tests. |
| QA-14 | Verify formatter directory/imports/serializers/YAML, universal flags, and governance promises are absent while module/test/build gates pass `(qa-lead)` | VAL-030, VAL-031, VAL-032, VAL-033, VAL-034 | Required `ls`/`rg` real exit statuses in declared expected-non-zero form, `go mod why -vendor` build-graph absence, command-tree tests, tidy diff, full gate. |
| QA-15 | Verify exactly 42 executable leaves form one disjoint 21-retained/21-legacy partition with correct root-flag transport behavior `(qa-lead)` | VAL-035 | Literal-list Cobra drift test output; no family ellipses or inferred membership. |
| QA-16 | Verify the complete canonical JSON journey, specified typed failure, rebuild, and byte-identical Issue reads `(qa-lead)` | VAL-036 | Task 18 integration output and captured pre/post byte comparison artifact. |
| QA-17 | Verify architecture, playbook, and changelog describe the delivered tree and no current formatter/stub promise survives `(qa-lead)` | VAL-037 | Live command/source tree, rendered/file review, and exact cited doc locations. |
| QA-18 | Verify exact config key/domain round-trip and all atomic-init failure/success stages `(qa-lead)` | VAL-039, VAL-042, VAL-043 | Exact TOML key inventory, parser boundary matrix, filesystem snapshots, rebuild comparison. |
| QA-19 | Verify every retained leaf's exact closed args/flags/value domains, with valid boundaries reaching and invalid/removed inputs not reaching handlers `(qa-lead)` | VAL-044 | Literal request-schema/Cobra table and just-inside/outside/malformed handler-sentinel output for all 21 leaves. |
| QA-20 | Verify every declared collection/change order and explicit bound, including equal keys and at/over-limit states with no truncation/mutation `(qa-lead)` | VAL-045 | Complete ordering registry fixtures, at-limit success, typed over-limit and unchanged-state output. |
| QA-21 | Review security/data-loss risks for redaction, config/path containment, symlink escape, lock atomicity, pre-activation residue, and rejected-operation effects `(qa-lead)` | VAL-001, VAL-002, VAL-003, VAL-004, VAL-042, VAL-043, VAL-044 | Adversarial test output, filesystem/journal state, and no-secret/internal-path response captures. |
| QA-22 | Run performance spot-check and compare rebuild/Issue-registry benchmark completion for regressions or nontermination `(qa-lead)` | no VAL | Current-session benchmark output; unexplained material regression is a finding, not silently accepted. |
| QA-23 | Verify TDD/test quality: first-failing-test intent, behavioral assertions, both-state coverage, temporary DB/repos, no weakened/deleted semantic assertions `(qa-lead)` | no VAL | Completion-report red/green evidence, committed tests, review observations, reproduced commands. |
| QA-24 | Verify all six review windows, completion reports, contract evidence, deviations, dependency notes, and independent consistency findings are resolved before final verdict `(qa-lead)` | no VAL | `review-1.md` through `review-6.md`, task reports, verified commits, and escalation state. |

## 8. Mid-Phase Review Cadence

All six windows are mandatory. A later task does not start until the preceding
window's review artifact records `VERDICT: CONTINUE` and every Blocking finding
or escalation is resolved. Reviewers compare each completion report's exact
commands with both its task spec and this plan; missing, silently skipped, or
narrower substituted commands are Blocking unless the report is a blocker
report with a reason.

| Review | Closes after | Reviewed tasks/reports | In-window VAL focus | Minimum independently reproduced command | Additional focus |
|---|---|---|---|---|---|
| Review 1 | Task 3 | Tasks/reports 1-3 | VAL-001..005, VAL-008, VAL-017, VAL-026, VAL-027, VAL-039, VAL-042, VAL-043, VAL-044 | `go test -tags integration -race ./...` | One config model, one root-flag registry, discovery boundary, strict parser, fixed maintenance grace, atomic init stages. |
| Review 2 | Task 6 | Tasks/reports 4-6 | VAL-006, VAL-007, VAL-009, VAL-010, VAL-011, VAL-012 plus shared VAL-001, VAL-002, VAL-008, VAL-017, VAL-044 | `go test -bench . -benchmem ./...` | Fail-closed authority, deterministic locks, migration/down path, all owned revision families, live/rebuild progress. |
| Review 3 | Task 9 | Tasks/reports 7-9 | VAL-013..020, VAL-040, VAL-041, VAL-044, VAL-045 plus shared VAL-002, VAL-008, VAL-010, VAL-011, VAL-017 | `go test -tags integration -race ./cmd/fid/integration/...` | Legal/illegal hierarchy, exact reads, both ready states, concurrency and transition atomicity, owner boundary. |
| Review 4 | Task 12 | Tasks/reports 10-12 | VAL-021..026, VAL-046 plus shared VAL-008, VAL-010, VAL-011, VAL-017, VAL-041, VAL-044, VAL-045 | `go test -tags integration -race ./cmd/fid/integration/...` | AC/label/link/workflow structural contracts, cycles and diamond, opaque endpoints, V1 authoring absence. |
| Review 5 | Task 15 | Tasks/reports 13-15 | VAL-027, VAL-028, VAL-029, VAL-031 plus shared VAL-002, VAL-008, VAL-017, VAL-044, VAL-045 | `go test -tags integration -race ./cmd/fid/integration/...` | Same-fixture maintenance branches, semantic legacy JSON, journal boundary, no weakened caller tests. |
| Review 6 | Task 17 | Tasks/reports 16-17 | VAL-030..034 and documentation progress for VAL-037 | `go test -race -coverprofile=coverage.out ./...` | Mandatory independent consistency review of accepted CLI/config/schema/Event/removal contracts against `contracts.md`, canonical spec amendments, live tree, and Tasks 1-17; execute all Task 16/17 absence checks. |

Review artifacts use `project-docs/templates/mid-phase-review-template.md` and
are written as `review-1.md` through `review-6.md`. Each review records an
outcome or explicit later owner for every in-window VAL it examines. Partial
coverage may continue only when the original contract names a later close
task/window; it may not be relabelled PASS.

After Review 6 continues, Task 18 runs. Phase-end QA begins only after Task 18's
committed completion report and exact-command evidence are verified. Task 18
cannot fix behavior from Tasks 1-17; any discovered mismatch blocks and returns
to its owning scope through escalation/triage.

## 9. Human Verification Scope

Only Human judgment appears here. Deterministic shape, channel, path, and
literal-output assertions remain in section 7. The final Human checklist uses
the committed Phase 20 behavior in a fresh disposable project and provides the
exact setup/interaction sequence generated during final QA; no held-out
scenario material is required or read in this phase.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-038a | Judge whether the routine Human-mode standalone Issue journey is concise enough for direct terminal use `(human)` | VAL-038 | Pass/Fail/Partial observation with the interaction transcript or notes and the command step that felt excessive, if any. |
| HUMAN-038b | Provoke the planned Human-mode failure and judge whether its guidance identifies a legitimate next action `(human)` | VAL-038 | Pass/Fail/Partial observation quoting or identifying the actionable/opaque part without treating prose as machine evidence. |
| HUMAN-038c | Judge whether any Human journey or failure message advises editing fid storage/config internals, directly or indirectly `(human)` | VAL-038 | Pass/Fail/Partial observation identifying the offending message if present; deterministic literal-path checks are supporting evidence only. |

The Human records each sub-ID independently. One passed observation cannot
mask another failed or partial observation. Colour, TUI/dashboard hierarchy,
and other deferred visual judgments are out of scope.

## 10. Phase-End Decision Rule

The phase can proceed to final review only when:

- Tasks 1-17 and all six mandatory review windows completed in the stated
  order, Review 6 recorded `CONTINUE`, and Task 18 then completed.
- Every required task-spec and QA-plan command is recorded with exact command
  text and current-session terminal evidence as `PASS`, or a blocker report
  records it `BLOCKED` with reason. Missing, skipped, hung, interrupted,
  killed, or narrower substituted commands are not passes.
- The full phase gate and every `(qa-lead)` checklist item have reproducible
  evidence; code reading alone does not satisfy deterministic checks.
- Every `VAL-001..VAL-046` receives exactly one final-review outcome:
  `PASS`, `FAIL`, `BLOCKED`, or `CARRY-FORWARD`, with cited evidence. `PASS`
  without evidence is invalid.
- Both qualifying and non-qualifying states are evidenced for every
  conditional/binary assertion, and all reference/parity claims use the
  deterministic comparisons named in this plan.
- VAL-038 is copied to the phase-end Human checklist only as
  `HUMAN-038a..c`; each observation is recorded independently.
- Architecture documentation, playbook, and changelog updates named by Tasks
  16/18 are committed and agree with the live command/source tree.
- Review findings, dependency-budget questions, contract deviations, and
  security/data-loss/journal-integrity concerns have explicit dispositions;
  no unresolved Blocking issue or escalation remains.

Final readiness at kickoff: **READY**. All 46 assertions have feasible
evidence, every task acceptance section covers its assigned assertions, all
required commands and artifact assumptions are concrete and validated, and no
ESC is required before Task 1.

## 11. DEC-133 Infrastructure Pack Addendum

The Human-approved Tasks 3e-3l are test/delivery infrastructure and cover no
product VAL. Their binding scope, evidence, and checkpoint state live in
`../infra-improvement-pack.md` and their task specs. They do not change the
phase validation contract or six accepted product review windows.

Task 3f runs both the former repository-wide integration command and the new
dedicated-package command, proves the complete tagged-test inventory, then
updates active gate surfaces to
`go test -tags integration -race ./cmd/fid/integration/...`. Reviews 3-5 and
remaining Task 4-18 specs use that command after migration; completed historical
artifacts remain unchanged.

A fresh QA Reviewer performs one additional bounded review after Task 3l and
final measurements, writing `qa/infra-review-1.md`. It independently checks
test inventory, gate equivalence, real-binary reuse, fixture parity/isolation,
working-directory semantics, race stability, clock behavior, documentation
sync, and timing reproducibility. Task 4 is blocked until this review is
non-BLOCK and any corrections pass the final gate.
