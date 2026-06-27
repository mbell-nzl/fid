# Phase 12 QA Plan - Requirements + Links

> Phase: 12 - Requirements + links (v0.1d slice 3/4)
> Refreshed: 2026-06-27
> QA Lead: Codex
>
> Location: `project-docs/docs/phases/phase-12/qa/qa-plan.md`.

## 1. Scope And Source Documents

Phase 12 ships the v0.1d requirement pointer and typed-link surface:
`fid requirement sync`, `fid link add`, `fid link remove`, and
`fid link list`. Requirements ship first so issue -> requirement `implements`
links have valid endpoints. The phase brief and validation contract state that
`requirements` and `links` schema/triggers already exist; schema migration is
not expected unless implementation proves and records a spec-backed blocker.

Source artifacts used:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/01_PRD.md`
- `project-docs/specs/02_TRD.md`
- `project-docs/specs/03_TRD_SCHEMA.md`
- `project-docs/specs/05_DESIGN_PLAN.md`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/GLOSSARY.md`
- `project-docs/templates/qa-plan-template.md`
- `project-docs/templates/mid-phase-review-template.md`
- `project-docs/templates/final-review-template.md`
- `project-docs/templates/human-checklist-template.md`
- `project-docs/docs/phases/phase-12/qa/validation-contract.md`
- `project-docs/docs/phases/phase-12/tasks/task-1.md`
- `project-docs/docs/phases/phase-12/tasks/task-2.md`
- `project-docs/docs/phases/phase-12/tasks/task-3.md`
- `project-docs/docs/phases/phase-12/tasks/task-4.md`

## 2. Validation Contract Refinement Check

- All in-scope PRD/TRD/design requirements referenced by at least one
  `VAL-NNN`: YES.
- Every `VAL-NNN` has one validator role: YES.
- Every `VAL-NNN` has a concrete evidence method: YES.
- Every `VAL-NNN` maps to a task/window or phase-end check: YES.
- Every `VAL-NNN` describes one observable invariant, or human-checklist fan-out
  is explicitly required: YES. VAL-012 is human-only and is fanned out in
  Section 13. VAL-001, VAL-003, VAL-007, VAL-008, VAL-010, and VAL-011 require
  multiple deterministic subcases, but each subcase proves one bounded parent
  invariant from the contract.
- Conditional or binary `VAL-NNN` assertions require both qualifying and
  non-qualifying states: YES. VAL-002..009 and VAL-011 cannot pass on
  positive-only evidence.
- Parity or "matches reference" `VAL-NNN` assertions have deterministic
  conformance proof when available: YES. VAL-001 uses spec-vs-matrix-vs-policy
  output, VAL-008 uses registry-vs-documented-set evidence, VAL-010 uses
  four-format structural parity against JSON, and VAL-011 uses rebuild parity.
- Revisions requested from Tech Lead: none before implementation. Planning
  interpretations PI-1 and PI-2 must be recorded in Phase 12 decisions by the
  owning tasks as already required by the task specs. If they are not recorded,
  final QA treats that as a documentation defect, not as permission to edit the
  validation contract.
- Contract/spec mismatches found during QA planning: none. The contract and task
  specs consistently say no schema migration is expected.

## 3. Requirement Coverage Map

| Requirement | Phase interpretation | Covered by | Deterministic evidence |
|---|---|---|---|
| PRD-LINK-001 | Typed directional issue links | VAL-005 | Link add CLI/domain evidence for declared + inverse issue links. |
| PRD-LINK-002 | Issue -> requirement links | VAL-002, VAL-006 | Requirement sync row evidence, then `implements`/`implemented_by` link evidence. |
| PRD-LINK-003 | Auto-created inverse links | VAL-005, VAL-006 | Link add effects and DB rows prove declared + inverse in one event. |
| PRD-LINK-004 | Closed relation/inverse set | VAL-005, VAL-008 | Registry equality test and CLI rejection of unknown/self links. |
| PRD-LINK-005 | Invalid relation/entity pair rejected | VAL-007 | CLI error cases prove before-append rejection and compatible success. |
| PRD-LINK-006 | List all links from/to entity | VAL-010 | `link list` outgoing/incoming evidence in json, toon, markdown, yaml. |
| PRD-LINK-006a | Remove link + inverse and activity evidence | VAL-009 | Remove event evidence, `links` absence, `activity_log` row evidence. |
| PRD-LINK-008 | Generic links must not duplicate dedicated schema | VAL-008 | Registry equality and unknown hierarchy-style relation rejection. |
| PRD-GOV-006 | Requirement pointer records with source/hash | VAL-002, VAL-003 | Sync fixtures prove persisted pointer fields, hash change, malformed rejection. |
| TRD-ARCH-010 | Matrix code constant; unknown roles rejected | VAL-001, VAL-004 | Matrix tests, `Authorize` allowed/denied pairs, `fid policy show` parity. |
| `03_TRD_SCHEMA.md` sections 5.2, 10, 13, 14.2-14.3 | Requirements/links schema, matrix rows, endpoint/ref triggers | VAL-001..011 | DB projection checks, trigger-compatible app validation, no-migration file-state check. |
| `05_DESIGN_PLAN.md` cross-surface fields | Link list field model | VAL-010 | Link list output contains `direction`, `relation`, `entity_type`, `entity_id`, `entity_title`. |
| `06_DESIGN_PLAN_CLI.md` `fid requirement sync` | Sync command flags/output/errors | VAL-002..004 | CLI tests and final seeded command evidence. |
| `06_DESIGN_PLAN_CLI.md` `fid link add` | Add command args/output/errors | VAL-005..008 | CLI tests and final seeded command evidence. |
| `06_DESIGN_PLAN_CLI.md` `fid link remove` | Remove command args/output/errors | VAL-009 | CLI tests and final seeded command evidence. |
| `06_DESIGN_PLAN_CLI.md` `fid link list` | List command output formats | VAL-010 | Four-format parsed parity evidence. |
| `08_QUALITY_PLAN.md` COV-001/COV-002, sections 5.2, 5.6 | Coverage, phase gate, feasibility validation, human checklist split | no VAL | Required commands, coverage report, final QA artifacts, command validation evidence. |

## 4. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | Requirement pointers, `fid requirement sync`, `requirement.sync` matrix row | VAL-001 partial, VAL-002, VAL-003, VAL-004 | High | New write command, parse contract, hash behavior, no-op sync event suppression, before-append errors. |
| Task 2 | Links domain, relation registry, `fid link add`, `link.add`/`link.remove` matrix rows | VAL-001 completes, VAL-005, VAL-006, VAL-007, VAL-008 | High | Auto-inverse event atomicity, endpoint compatibility, missing endpoints, policy parity, no activity row on add. |
| Task 3 | `fid link remove` | VAL-009 | Medium | Must delete declared + inverse in one event and project removal evidence to `activity_log`. |
| Task 4 | `fid link list`, all-format parity, end-to-end flow and rebuild parity | VAL-010, VAL-011 | High | Read path must use unlocked freshness; integration test must prove event counts and rebuild parity. |
| Phase end | Human dogfood checklist | VAL-012 | Medium | Deterministic gates must pass first; human checklist contains only dogfooding/legibility observations. |

## 5. Validation Contract Evidence Strategy

| VAL ID | Evidence method refinement | Validator | Task/window | Acceptance signal | Reject/block signal |
|---|---|---|---|---|---|
| VAL-001 | Completion reports and final QA cite matrix unit tests, allowed/denied `Authorize` cases, and `fid policy show` rows for exactly `requirement.sync`, `link.add`, `link.remove`. | QA Lead | Tasks 1-2; Review 1 | Rows match section 13 and policy show/enforce parity is proven. | Missing/duplicated row, role boolean mismatch, or show/enforce drift. |
| VAL-002 | Sync fixture creates rows; unchanged rerun appends no event; changed summary updates count and hash. Final QA repeats seeded happy/no-op/update flow or cites literal task output plus DB/journal state evidence. | QA Lead | Task 1; Review 1 | Created/updated/unchanged counts and persisted fields match contract. | Field loss, wrong counts, no-op append, or stale hash. |
| VAL-003 | Error fixtures cover no parseable IDs, duplicate IDs, unknown `--type`, and missing `--id`; each records non-zero exit and unchanged journal/DB. | QA Lead | Task 1; Review 1 | All malformed states reject before append with actionable stderr. | Any malformed state appends/projects, exits zero, or hides cause. |
| VAL-004 | Role cases include tl/human allowed and denied/unknown roles rejected before append naming role + action. | QA Lead | Task 1; Review 1 | Allowed roles append one event; denied/unknown roles append none. | Denied/unknown succeeds or allowed role fails. |
| VAL-005 | Domain/CLI evidence for issue<->issue add proves two `links` rows and one journal event; duplicate add proves zero append. | QA Lead | Task 2; Review 1 | Declared + inverse rows project atomically and confirmation names both. | Wrong inverse, two events, missing confirmation, duplicate append. |
| VAL-006 | Requirement is synced first, then `implements` add creates `implemented_by`; missing requirement endpoint rejects before append. | QA Lead | Task 2; Review 1 | Existing requirement link works; missing requirement does not append. | Wrong inverse, missing endpoint accepted, or linked row absent. |
| VAL-007 | Incompatible issue->issue `implements`, issue->requirement `blocks`, and missing endpoint cases reject; compatible control case succeeds. | QA Lead | Task 2; Review 1 | Error names relation and endpoints, and journal/DB unchanged. | Incompatible/unresolvable case accepted or error lacks context. |
| VAL-008 | Registry equality test proves exactly seven semantic relations plus inverses and no dedicated-schema relation; CLI unknown relation and self-link reject before append. | QA Lead | Task 2; Review 1 | Registry and rejection behavior match docs. | Extra/missing relation, hierarchy/escalation/resolution relation, unknown/self accepted. |
| VAL-009 | Add/remove fixture proves declared + inverse delete effects in one event, zero `links` rows afterward, activity row exists, confirmation names both, nonexistent remove rejects before append. | QA Lead | Task 3; Review 2 | Remove mechanics and activity evidence match contract. | Surviving row, missing activity evidence, two events, nonexistent append. |
| VAL-010 | Fixture has outgoing and incoming links across issue and requirement endpoints; json is parsed as reference; toon/markdown/yaml parsed or structurally compared to the same rows/fields. | QA Lead | Task 4; Review 2 | All formats expose same rows and required fields with resolved titles. | Missing direction/field/title or format divergence. |
| VAL-011 | Integration test runs real command sequence in temp workspace, checks one event per write, no nested-lock failure, and rebuild parity for `links`/`requirements`. | QA Lead | Task 4; Review 2 | Sequence, event counts, and rebuild parity pass. | Lock failure, wrong event count, or rebuild divergence. |
| VAL-012 | Phase-end human checklist fans out sync counts, add inverse legibility, bidirectional list legibility, and remove inverse legibility. | Human | Phase end | Human can understand outcomes from CLI output without implementation knowledge. | Human cannot identify one or more outcomes or checklist is not run. |

## 6. Required Evidence Strategy

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | `gofmt -l cmd internal` output showing no files listed. | no VAL | Developer / QA Lead | Every task and final review |
| EVID-02 | `golangci-lint run ./...` PASS output. | no VAL | Developer / QA Lead | Every task and final review |
| EVID-03 | `go test -race -coverprofile=coverage.out ./...` PASS output. | VAL-001..011, no VAL | Developer / QA Lead | Every task and final review |
| EVID-04 | `go tool cover -func=coverage.out` package/total output, including `cmd/fid` and relevant `internal/*` package percentages. | no VAL | Developer / QA Lead | Every task and final review |
| EVID-05 | `go build -o /tmp/fid ./cmd/fid` PASS output for each CLI-task close and final review. | no VAL | Developer / QA Lead | Tasks 1-4 and final review |
| EVID-06 | `go test -tags integration -race ./...` PASS output when required by task scope and at phase end. | VAL-011, no VAL | Developer / QA Lead | Task 4 and final review |
| EVID-07 | `go test -bench . -benchmem ./...` PASS output or explicit task-level N/A where allowed. | no VAL | Developer / QA Lead | Task 4 conditional and final review |
| EVID-08 | Completion report command tables match exact task-spec commands, including conditional command PASS/N/A/BLOCKED status. | VAL-001..011, no VAL | QA Reviewer / QA Lead | Reviews 1-2 and final review |
| EVID-09 | Matrix test output, `Authorize` case output, and `fid policy show` output for Phase 12 rows. | VAL-001, VAL-004 | Developer / QA Lead | Tasks 1-2 and final review |
| EVID-10 | Requirement sync CLI/domain output plus DB/journal state for create, no-op, update, malformed inputs, and role denial. | VAL-002, VAL-003, VAL-004 | Developer / QA Lead | Task 1 and final review |
| EVID-11 | Link add CLI/domain output plus DB/journal state for issue links, requirement links, duplicate, incompatible, missing endpoint, unknown relation, self-link, and role denial. | VAL-005, VAL-006, VAL-007, VAL-008 | Developer / QA Lead | Task 2 and final review |
| EVID-12 | Link remove CLI/domain output plus DB/journal/activity state for removal and nonexistent removal. | VAL-009 | Developer / QA Lead | Task 3 and final review |
| EVID-13 | Link list output in json, toon, markdown, yaml, with parsed/structural parity evidence against JSON. | VAL-010 | Developer / QA Lead | Task 4 and final review |
| EVID-14 | End-to-end integration output and rebuild parity state comparison for requirements and links. | VAL-011 | Developer / QA Lead | Task 4 and final review |
| EVID-15 | File-state proof that no new migration file was added unless a task recorded a spec-backed blocker/decision. | VAL-002..011, no VAL | QA Lead | Final review |
| EVID-16 | Human checklist at `project-docs/docs/phases/phase-12/qa/human-checklist.md`, with HV-12a..HV-12d outcomes. | VAL-012 | Human | Phase end |

## 7. Acceptance Criteria And Expected Tests

Task 1 acceptance:

- `internal/requirements` owns parse and upsert rules; CLI owns flag/output
  wiring only.
- `requirement.sync` matrix row exists exactly once, matches section 13
  (tl/human only), and appears in `fid policy show`.
- Sync parses PRD/TRD table IDs, derives type, stores summary/source/hash/status,
  reports created/updated/unchanged, appends one event for changes, and appends
  no event for unchanged input.
- `--type` and `--id` filters behave as specified.
- Malformed input and denied/unknown roles reject before append.

Expected Task 1 tests:

- Domain tests for parser, hash determinism, upsert create/update/unchanged.
- CLI tests for sync counts, filters, malformed fixtures, role permissions, and
  policy show parity.
- Journal/DB assertions for one-event changed sync and zero-event unchanged or
  rejected sync.

Task 2 acceptance:

- `internal/links` owns relation registry, inverse lookup, endpoint resolution,
  compatibility validation, and add service.
- `link.add` and `link.remove` matrix rows exist exactly once, all seven roles
  allowed, and appear in `fid policy show`.
- Add creates declared + inverse rows in one event and does not create an
  `activity_log` row.
- Duplicate, incompatible relation/entity pairs, unresolvable endpoints, unknown
  relation, self-link, and unknown role reject before append.

Expected Task 2 tests:

- Registry equality and inverse tests for blocks/blocked_by, relates_to,
  duplicates/duplicated_by, supersedes/superseded_by, implements/implemented_by,
  caused_by/causes, discovered_from/discovered.
- Domain/CLI tests for valid issue<->issue and issue->requirement adds.
- Error-path CLI tests asserting stderr, non-zero exit, and zero journal growth.
- Matrix and `Authorize` parity tests.

Task 3 acceptance:

- Remove deletes declared + inverse links in exactly one event.
- Remove projects an `activity_log` row recording the removal.
- Confirmation names both removed directions.
- Nonexistent link and unknown role reject before append.

Expected Task 3 tests:

- Domain/CLI tests for removing issue<->issue and issue->requirement links.
- DB assertions for zero surviving link rows and one activity row.
- Journal event count assertions for successful and rejected removes.

Task 4 acceptance:

- `link list` returns outgoing and incoming rows for an entity with `direction`,
  `relation`, `entity_type`, `entity_id`, and resolved `entity_title`.
- json, toon, markdown, and yaml expose equivalent link data.
- Ordering is deterministic and rebuild-stable.
- End-to-end integration flow proves sync -> add -> list -> remove, one event
  per write, no nested-lock failure, and rebuild parity.

Expected Task 4 tests:

- Domain/CLI list tests for outgoing, incoming, mixed issue/requirement
  endpoints, empty list, title resolution, deterministic order, and all formats.
- Integration test under `-tags integration` for the full Phase 12 flow and
  rebuild parity.

## 8. Deterministic Gates

No referenced command is a placeholder. Commands are validated in Section 9.

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | No files listed | Every task and final review |
| Lint/static analysis | `golangci-lint run ./...` | PASS | Every task and final review |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS | Every task and final review |
| Coverage report | `go tool cover -func=coverage.out` | PASS and COV-001/COV-002 threshold evidence | Every task and final review |
| Integration | `go test -tags integration -race ./...` | PASS, or task-level N/A only when the task spec allows it and no integration fixture was touched | Task 4 and final review; conditional for Tasks 1-3 |
| End-to-end | `go test -tags integration -race ./...` | PASS with Phase 12 flow/rebuild test evidence | Task 4 and final review |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS | Final review; Task 4 only if benchmark/scale fixtures are touched, otherwise N/A with reason |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | Tasks 1-4 and final review |
| Full phase check | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go test -bench . -benchmem ./...`; `go build -o /tmp/fid ./cmd/fid` | PASS | Final review |

## 9. Command And Feasibility Validation

Validated on 2026-06-27 before writing this plan. Validation here proves syntax,
tool presence, file-state assumptions, and artifact feasibility; it does not
substitute for task or phase gate execution.

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | `gofmt -h`; file-state check `ls cmd internal` | PASS | Help shows `usage: gofmt [flags] [path ...]` and `-l list files whose formatting differs`; `cmd/fid` and current `internal/*` packages exist. |
| `golangci-lint run ./...` | `golangci-lint --version` | PASS | Output: `golangci-lint has version 2.12.2 built with go1.26.2`. |
| `go test -race -coverprofile=coverage.out ./...` | `go help testflag`; `go help build`; `go list ./...`; `go env GOCACHE GOMODCACHE GOPATH` | PASS | Help recognizes `-coverprofile` and `-race`; `go list ./...` listed 10 packages including `cmd/fid`; caches resolve inside repo: `.gocache`, `.gomodcache`, `.gopath`. |
| `go tool cover -func=coverage.out` | `go tool cover -h` | PASS | Help output includes `-func string output coverage profile information for each function`. |
| `go test -tags integration -race ./...` | `go help build`; `go help test`; `go list ./...` | PASS | Build help documents `-tags tag,list`; test help documents package-list mode; `./...` currently resolves to repo packages. |
| `go test -bench . -benchmem ./...` | `go help testflag`; `go list ./...` | PASS | Help documents `-bench regexp` and `-benchmem`; `./...` package scope resolves. |
| `go build -o /tmp/fid ./cmd/fid` | `go help build`; `go list ./cmd/fid` via `go list ./...` output | PASS | Build help documents `-o output`; package `github.com/mbell-nzl/fid/cmd/fid` exists; `/tmp` is writable by environment policy. |
| Task 1 command set | Task spec inspection plus command validations above | PASS | Task 1 requires standard gate plus build; integration conditional. All commands are concrete and validated. |
| Task 2 command set | Task spec inspection plus command validations above | PASS | Task 2 requires standard gate plus build; integration conditional. All commands are concrete and validated. |
| Task 3 command set | Task spec inspection plus command validations above | PASS | Task 3 requires standard gate plus build; integration conditional. All commands are concrete and validated. |
| Task 4 command set | Task spec inspection plus command validations above | PASS | Task 4 requires standard gate plus integration and build; benchmark conditional. All commands are concrete and validated. |
| Phase-level command gate | `08_QUALITY_PLAN.md` section 5.2 and AGENTS command list inspection | PASS | Phase-level gate commands are concrete; release-readiness security/dependency/license/secret commands are not required for this planned phase. |
| Coverage threshold feasibility | `08_QUALITY_PLAN.md` COV-001/COV-002 inspection; `go tool cover -h` | PASS | COV-001 requires 80% across named domain packages; COV-002 requires 70% for `cmd/fid`; cover tool can print function/package percentages. |
| No schema migration expected | File-state and spec/migration search | PASS | `migrations/0002_v01_tables.up.sql` contains `CREATE TABLE requirements` and `CREATE TABLE links`; `migrations/0005_v01_ref_triggers.up.sql` contains link endpoint/ref/compatibility triggers; validation contract says migration is not expected. |
| Phase 12 task artifacts | `find project-docs/docs/phases/phase-12 -maxdepth 3 -type f | sort` | PASS | Files present: validation contract and task-1.md through task-4.md. |
| QA artifact templates | `ls project-docs/templates`; template inspection | PASS | `qa-plan-template.md`, `mid-phase-review-template.md`, `final-review-template.md`, and `human-checklist-template.md` exist and were inspected. |
| Requirement source references | `rg` over specs for PRD/TRD/design IDs and CLI headings | PASS | Source IDs and headings found: PRD-LINK-001..006a/008, PRD-GOV-006, TRD-ARCH-010, `fid requirement sync`, `fid link add/remove/list`. |
| Human checklist feasibility for VAL-012 | Human checklist template inspection and VAL-012 fan-out in Section 13 | PASS | Template requires one observable per row with sub-IDs; HV-12a..HV-12d are defined below. |

## 10. Task-Window Feasibility Validation

| Window | Commands and artifact assumptions | Validation result | Required review evidence |
|---|---|---|---|
| Task 1 / Review 1 | Standard gate, build, conditional integration, requirement parser fixtures, `.fid` temp workspace, no migration. | PASS. Commands validated in Section 9; PRD/TRD source tables exist; requirements DDL exists in migration 0002. | Completion report must show exact command results, build PASS, integration PASS/N/A reason, sync fixture state, policy row evidence, and no-migration state. |
| Task 2 / Review 1 | Standard gate, build, conditional integration, existing requirement endpoint from Task 1, links DDL/triggers, no migration. | PASS. Commands validated; links DDL and triggers exist; Task 2 dependency on Task 1 is explicit. | Completion report must show exact command results, build PASS, integration PASS/N/A reason, registry equality, link add DB/journal state, rejection zero-append evidence, and no activity-log add evidence. |
| Task 3 / Review 2 | Standard gate, build, conditional integration, Task 2 link helpers, activity-log ref trigger, no migration. | PASS. Commands validated; `activity_log` DDL and ref trigger exist; Task 3 dependency on Task 2 is explicit. | Completion report must show exact command results, build PASS, integration PASS/N/A reason, declared+inverse delete state, `activity_log` evidence, and nonexistent remove zero-append evidence. |
| Task 4 / Review 2 | Standard gate plus integration, build, conditional benchmark, all Phase 12 commands available, rebuild parity fixture, no held-out scenarios. | PASS. Commands validated; Task 4 depends on Tasks 1-3; templates and integration command are feasible. | Completion report must show exact command results, integration PASS, build PASS, benchmark PASS/N/A reason, four-format parity output, end-to-end event counts, no nested-lock error, and rebuild parity state comparison. |

## 11. Deterministic Checklist

Every item below is executed by QA Lead at final review unless explicitly covered
by a mid-phase QA Reviewer with reproducible evidence. Each check must end with
`(qa-lead)` and must cite literal output, file-state proof, logs, or parsed
artifact evidence in `final-review.md`.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Verify all task completion reports exist and record exact required commands from their task specs, including conditional commands as PASS, N/A with reason, or BLOCKED with reason `(qa-lead)` | VAL-001..011, no VAL | File-state proof for task-1-done.md through task-4-done.md and copied command result excerpts. |
| QA-02 | Run the full phase command gate exactly as listed in Section 8 `(qa-lead)` | VAL-001..011, no VAL | Literal command output for format, lint, coverage test, coverage report, integration, benchmarks, and build. |
| QA-03 | Verify coverage thresholds from `go tool cover -func=coverage.out`: named domain packages meet COV-001 and `cmd/fid` meets COV-002 `(qa-lead)` | no VAL | Coverage output excerpts with relevant package percentages. |
| QA-04 | Verify no new migration file or migration edit exists unless a Phase 12 task recorded a spec-backed blocker/decision `(qa-lead)` | VAL-002..011, no VAL | `git diff --name-only -- migrations` and task decision/completion report state. |
| QA-05 | Verify Phase 12 matrix rows match `03_TRD_SCHEMA.md` section 13 and `fid policy show` displays the same source enforced by authorization `(qa-lead)` | VAL-001, VAL-004 | Test output plus CLI output for `requirement.sync`, `link.add`, and `link.remove`. |
| QA-06 | Execute or verify a seeded `fid requirement sync` create/no-op/update sequence, including journal event counts and persisted pointer fields `(qa-lead)` | VAL-002 | CLI output, DB query/state output, and journal line/event count evidence. |
| QA-07 | Execute or verify malformed `requirement sync` cases for no parseable IDs, duplicate IDs, unknown `--type`, and absent `--id`, each before append `(qa-lead)` | VAL-003 | Non-zero stderr excerpts and before/after journal/DB counts. |
| QA-08 | Execute or verify `requirement.sync` authorization cases for tl/human allowed and denied/unknown roles rejected before append `(qa-lead)` | VAL-004 | CLI/test output naming role + action and journal count evidence. |
| QA-09 | Execute or verify valid issue<->issue `fid link add` creates declared + inverse rows in one event and duplicate add rejects before append `(qa-lead)` | VAL-005 | CLI confirmation, DB rows, event count, duplicate stderr, unchanged journal count. |
| QA-10 | Execute or verify issue -> requirement `implements` link creates `implemented_by` inverse after requirement sync, and missing requirement rejects before append `(qa-lead)` | VAL-006 | Sync evidence, CLI output, DB rows, missing-endpoint stderr, unchanged journal count. |
| QA-11 | Execute or verify incompatible and unresolvable `link add` cases reject before append while a compatible control case succeeds `(qa-lead)` | VAL-007 | Error excerpts naming relation/endpoints, journal/DB before/after counts, compatible success output. |
| QA-12 | Verify relation registry equals the documented seven relations/inverses and rejects unknown relation plus self-link before append `(qa-lead)` | VAL-008 | Unit/CLI output and journal count evidence. |
| QA-13 | Execute or verify `fid link remove` deletes declared + inverse links in one event, writes `activity_log`, confirms both directions, and nonexistent remove rejects before append `(qa-lead)` | VAL-009 | CLI output, DB rows absent, activity row present, event count, nonexistent stderr and unchanged journal count. |
| QA-14 | Execute or verify `fid link list` for an entity with outgoing and incoming links in json, toon, markdown, and yaml, with parsed structural parity against JSON `(qa-lead)` | VAL-010 | Four output excerpts and parsed/structural comparison evidence. |
| QA-15 | Verify the Phase 12 end-to-end integration flow covers sync -> issue link -> requirement link -> list both directions -> remove, one event per write, no nested-lock failure, and rebuild parity `(qa-lead)` | VAL-011 | Integration test output and state comparison excerpts for `links` and `requirements`. |
| QA-16 | Verify impacted docs and task reports record PI-1/PI-2 decisions or explicitly record no deviation, without editing specs from QA `(qa-lead)` | VAL-001..011, no VAL | File references in task completion reports and Phase 12 decisions file, or documented finding. |
| QA-17 | Generate `human-checklist.md` from the VAL-012 human-only subset and verify it contains only HV-12a..HV-12d plus summary fields `(qa-lead)` | VAL-012 | File-state proof and checklist excerpt. |

## 12. Mid-Phase Review Cadence

Default project cadence is a QA Reviewer after every three completed tasks.
Phase 12 uses a stricter two-task cadence because Tasks 1-2 establish the
shared requirement/link domain and policy foundation that Tasks 3-4 reuse.

- Review 1: after Task 2 completes, covering Tasks 1-2 and VAL-001..008. Task 3
  should not proceed until Review 1 verdict is `CONTINUE` or a blocking finding
  is explicitly resolved/deferred by the appropriate role.
- Review 2: after Task 4 completes and before final QA, covering Tasks 3-4 and
  VAL-009..011, plus command/report consistency for Tasks 3-4.
- Review template:
  `project-docs/templates/mid-phase-review-template.md`.
- Each review compares completion-report command results against the exact task
  spec commands and this QA plan. Missing, silently skipped, or narrower
  substituted required commands are Blocking unless the task report records
  `BLOCKED` with a reason.
- Each review verifies in-window `VAL-NNN` assertions and records whether
  evidence is ON TRACK, NOT ON TRACK, or explicitly DEFERRED to final QA.

## 13. Human Verification Scope

Only VAL-012 requires human judgment. These rows are copied into
`project-docs/docs/phases/phase-12/qa/human-checklist.md` at phase end using
`project-docs/templates/human-checklist-template.md`. Each row is one
observable and ends with `(human)`.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HV-12a | Dogfood `fid requirement sync` in a real `.fid` workspace and confirm the CLI output makes the recorded pointer counts understandable without reading implementation `(human)` | VAL-012 | Human notes the command run, observed counts, and whether created/updated/unchanged meaning was clear. |
| HV-12b | Dogfood `fid link add` and confirm the CLI output makes the auto-created inverse understandable `(human)` | VAL-012 | Human notes the command run and whether both declared and inverse directions were identifiable. |
| HV-12c | Dogfood `fid link list` for both linked endpoints and confirm the output makes outgoing and incoming directions understandable `(human)` | VAL-012 | Human notes the endpoints checked and whether direction, relation, entity, and title were legible. |
| HV-12d | Dogfood `fid link remove` and confirm the CLI output makes removal of both declared and inverse directions understandable `(human)` | VAL-012 | Human notes the command run and whether both removed directions were identifiable. |

Human checklist structure at phase end:

- File:
  `project-docs/docs/phases/phase-12/qa/human-checklist.md`.
- Section: `Requirements And Link Dogfood`.
- Rows: HV-12a through HV-12d only, plus template summary tables.
- Outcome values: `pass`, `fail`, `partial`, or `blank`.
- Overall result: `pass` only if all four rows pass.

## 14. Phase-End Decision Rule

The phase can proceed to final review only when:

- Required deterministic gates pass or are recorded `BLOCKED` with escalation.
- Every VAL-001..VAL-012 has planned evidence and a final-review outcome slot.
- Every `(qa-lead)` item has reproducible evidence, not prose-only confidence.
- Human checks are human-judgment-only and mapped to the human checklist.
- No unresolved blocking defects remain.
- Required documentation updates and decision records are named in task
  completion reports.
- No schema migration was added unless a spec-backed blocker/decision proves it
  was required.

Final review records every `VAL-NNN` outcome as `PASS`, `FAIL`, `BLOCKED`, or
`CARRY-FORWARD`. `PASS` without cited evidence is invalid.
