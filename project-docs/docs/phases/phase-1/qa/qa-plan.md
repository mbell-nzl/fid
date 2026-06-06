# Phase 1 QA Plan - End-to-End Persistence Loop (Tracer Bullet)

> Phase: 1 - End-to-End Persistence Loop (Tracer Bullet)
> Refreshed: 2026-06-06
> QA Lead: Codex QA Lead
>
> Status: READY WITH NOTED DOCUMENTATION MISMATCH. No required command is a
> placeholder.

## 1. Scope And Source Documents

Phase 1 proves the thinnest end-to-end persistence path through journal,
projection, and rebuild, exercised by real CLI commands: `fid init`,
`fid rebuild`, `fid issue create`, and `fid issue get --format json`.

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
- `project-docs/specs/09_AI_WORKFLOW.md`
- `project-docs/docs/phases/phase-1/qa/validation-contract.md`
- `project-docs/docs/phases/phase-1/tasks/task-1.md`
- `project-docs/docs/phases/phase-1/tasks/task-2.md`
- `project-docs/docs/phases/phase-1/tasks/task-3.md`
- `project-docs/docs/phases/phase-1/tasks/task-4.md`
- `project-docs/docs/phases/phase-1/tasks/task-5.md`
- `project-docs/docs/phases/phase-1/tasks/task-6.md`

Noted mismatch:

- `project-docs/docs/current/phase_brief.md` says Phase 1 is planned and
  awaiting this QA plan. `project-docs/docs/current/progress_current.md` still
  says Phase 1 is not yet planned. This is a documentation-state mismatch only;
  it does not change Phase 1 scope.

## 2. Requirements Coverage

| Requirement Source | Phase 1 Coverage | VAL IDs |
|---|---|---|
| PRD-INIT-001..004 | Workspace init, prefix handling, seeded roles, re-init rejection. | VAL-001..004 |
| PRD-OPS-001..004, PRD-OPS-012 | Atomic event writes, journal-derived projection, stale/missing DB rebuild, explicit rebuild, replay conflicts. | VAL-005..012 |
| PRD-ROLE-001..002 | Phase 1 write commands require `--by`; unknown roles are rejected. | VAL-016 |
| PRD-NFR-008 | Same-checkout concurrent writes are serialized by `.fid/write.lock`. | VAL-013 |
| PRD-OUT-001, PRD-OUT-004 | Issue reads are JSON by default and support `--format json`. | VAL-006 |
| TRD §1, §3, §4.2..4.4; TRD CONFIG §1..4; TRD SCHEMA §1, §2.1, §9.5, §9.7, §15 | Go/Cobra/SQLite persistence substrate, minimal migrations, event v1, canonical JSON, ID shapes, write lock, replay, checkpoint, config layout. | VAL-001..015 |
| TRD-ARCH-001..005, TRD-ARCH-007, TRD-ARCH-012 | Journal as source of truth, pure data events, upgrader chain, one write equals one event, adapter boundaries, required pragmas, format seam. | VAL-006..009, VAL-014, VAL-017 |
| TRD-PERF-006 | Two same-checkout writes both succeed without data loss. | VAL-013 |
| Design CLI §1, §2, §18; Design Plan §6.1 | User-visible init, rebuild, issue create/get, JSON shape, and error text expectations. | VAL-001, VAL-004, VAL-006, VAL-010, VAL-013 |

Scope explicitly deferred by the phase brief is not treated as a Phase 1 QA
failure: remaining v0.1 tables/triggers/views, full lifecycle, TOON/markdown/YAML
renderers, full permission matrix, held-out scenarios, and release-readiness
SAST/license/secret gates.

## 3. Validation Contract Refinement Check

- All in-scope requirements referenced by at least one `VAL-NNN`: YES.
- Every `VAL-NNN` has one validator role: YES.
- Every `VAL-NNN` has a concrete evidence method: YES, refined below.
- Every `VAL-NNN` maps to a task/window or phase-end check: YES.
- Every `VAL-NNN` describes one observable, or deterministic sub-evidence is
  required for each bundled state: YES WITH SUB-EVIDENCE.
- Conditional or binary assertions require both qualifying and non-qualifying
  states: YES. See VAL-002, VAL-004, VAL-010, VAL-011, VAL-012, VAL-013,
  VAL-015, VAL-016.
- Parity or "matches reference" assertions have deterministic conformance proof
  when available: YES. See VAL-005, VAL-006, VAL-008, VAL-009.
- Revisions requested from Tech Lead: none blocking. Record the
  `progress_current.md` mismatch at phase wrap or before Developer execution if
  the human wants current-state docs perfectly synchronized.

Refinement notes:

- VAL-012 bundles clean-read plus four mismatch classes. Final review must
  report each class separately and can pass VAL-012 only if all classes pass.
- VAL-013 bundles concurrent-write serialization and held-lock rejection. Final
  review must report both sub-evidence rows.
- VAL-016 bundles missing role flag, unknown role, and valid seeded-role states
  across `fid issue create` and `fid rebuild`. Final review must report every
  command/state combination.
- VAL-017 contains one fully deterministic import-direction check and one
  structural adapter-boundary check. Final review must include machine evidence
  for import direction and adapter-boundary smoke checks; QA source review is
  supporting evidence, not the sole pass basis.

## 4. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | SQLite connection pragmas and minimal migrations. | VAL-014 | Medium | Dependency additions require completion-report budget notes; later VALs depend on schema substrate. |
| Task 2 | Journal event v1, canonical JSON, ID helpers, upgrader chain. | VAL-007, VAL-009, VAL-015 | High | Golden v1 fixture becomes later replay evidence; no I/O should enter this task. |
| Task 3 | Write lock and append+flush journal writer. | VAL-007, VAL-013, VAL-015 | High | Race gate is mandatory; lock and sequence tests need both happy and contention paths. |
| Task 4 | Replay, projection, checkpoint, catch-up, `fid rebuild`. | VAL-008, VAL-009, VAL-010, VAL-011, VAL-012 | High | First user-facing rebuild command; build command is required by the task's conditional list. |
| Task 5 | `fid init` bootstrap. | VAL-001, VAL-002, VAL-003, VAL-004, VAL-005 | High | First complete write+projection loop; re-init must prove no file/journal mutation. |
| Task 6 | Minimal issue create/get, JSON format seam, role existence checks. | VAL-006, VAL-007, VAL-008, VAL-016 | High | Closes tracer bullet; role checks must include `fid rebuild` as well as issue create. |
| Phase end | Architecture boundary and final evidence audit. | VAL-001..017 | High | QA Lead reruns deterministic checklist and full phase gate before final verdict. |

## 5. Required Task Command Evidence

The phase runner can verify task reports against this table. Missing, silently
skipped, reordered into a narrower substitute, or unreported required commands
are blocking QA findings unless the completion report marks the command
`BLOCKED` with a reason.

| Task | Required exact commands before task close | Conditional exact commands and trigger |
|---|---|---|
| Task 1 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | `go test -tags integration -race ./...` only if migration tests are behind the `integration` build tag. |
| Task 2 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | None expected. |
| Task 3 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | `go test -tags integration -race ./...` if concurrency/lock tests are behind the `integration` build tag. |
| Task 4 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | `go build -o /tmp/fid ./cmd/fid` is required because the task adds `fid rebuild`; `go test -tags integration -race ./...` if replay/rebuild tests are behind the `integration` build tag. |
| Task 5 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go build -o /tmp/fid ./cmd/fid` | `go test -tags integration -race ./...` if init/rebuild round-trip tests are behind the `integration` build tag. |
| Task 6 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go build -o /tmp/fid ./cmd/fid` | `go test -tags integration -race ./...` if create/get/rebuild round-trip tests are behind the `integration` build tag. |

Phase-end QA Lead gate:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

## 6. VAL Evidence Matrix

Every final `PASS` requires task completion-report evidence plus QA Lead
phase-end rerun evidence. Test names below are expected intent; exact names may
vary, but the task report must make the matching test or command evidence
traceable.

| VAL | Deterministic Verification Evidence | Required Command Evidence | Feasibility Evidence |
|---|---|---|---|
| VAL-001 | Integration/CLI evidence from a temp dir that `fid init` creates `.fid/journal/`, `.fid/config.toml`, `.fid/pricing.toml`, `.fid/fid.sqlite`, and `.fid/.gitignore`; file-state proof with `find`/`grep`; config parse test. | Task 5 required commands; phase-end full gate; `/tmp/fid init` smoke evidence. | Task 5 names required files; `cmd`, `internal`, `migrations`, and `testdata` parents exist; `go build -o /tmp/fid ./cmd/fid` syntax validated by `go help build`. |
| VAL-002 | Integration/CLI evidence for `fid init --prefix X` and default init in a directory whose basename has a known uppercase form; parsed config assertions for both states. | Task 5 required commands; phase-end full gate. | `fid init --prefix` is specified in Task 5 and Design CLI §1; binary build command feasible. |
| VAL-003 | Integration evidence that the journal has exactly one `project.initialized` JSONL line and the projected roles set equals developer, qa_reviewer, qa_lead, tl_lite, tl, workflow_runner, human. | Task 5 required commands; phase-end full gate. | Task 5 creates bootstrap event builder and roles projection; minimal `roles` table comes from Task 1. |
| VAL-004 | Integration evidence for both states: fresh dir init succeeds; re-init exits non-zero with clear `.fid/ already exists` error and file/journal checksums or listing before/after are unchanged. | Task 5 required commands; phase-end full gate. | Task 5 acceptance criterion explicitly requires no-change proof; `grep` and file-state checks are available. |
| VAL-005 | Init, copy journal bytes, delete `.fid/fid.sqlite`, run `fid rebuild --by <seeded-role>`, assert roles/checkpoint match post-init state and journal bytes are identical by `cmp` or test assertion. | Task 5 required commands; Task 4 build evidence for `fid rebuild`; phase-end full gate. | `fid rebuild` is introduced in Task 4; Task 5 explicitly depends on Task 4 outputs and requires delete-DB rebuild evidence. |
| VAL-006 | CLI/integration evidence: create one issue through `fid issue create --kind task --title <t> --by developer`, capture ID, read via `fid issue get <id>` and `fid issue get <id> --format json`, validate well-formed JSON and field equality. | Task 6 required commands; phase-end full gate; `/tmp/fid` smoke evidence. | Task 6 specifies `cmd/fid/issue.go`, `internal/issues`, and `internal/format`; `go test` and build syntax validated. |
| VAL-007 | Unit/integration evidence that a single writer append and each single CLI write adds exactly one JSONL line with v1 fields `event_id`, `event_schema`, `event_type`, `occurred_at`, `actor`, `command`, and ordered `effects`. | Task 2, Task 3, and Task 6 required commands; phase-end full gate. | Task 2 defines v1 envelope; Task 3 defines append+flush; Task 6 proves full CLI path. |
| VAL-008 | Delete-DB rebuild tests for created issue: `fid issue get` JSON before deletion equals JSON after rebuild; replay tests prove effects use stored resolved values, not command logic. | Task 4 and Task 6 required commands; phase-end full gate. | Task 4 owns replay/projection; Task 6 closes issue create/get path; parity proof can be deterministic JSON byte or structural comparison. |
| VAL-009 | Golden-file v1 journal replay test and unit evidence that the upgrader chain dispatches by `event_schema` and returns v1 unchanged. | Task 2 and Task 4 required commands; phase-end full gate. | `testdata/` exists; Task 2 creates `testdata/journal/` golden v1 fixtures for Task 4 reuse. |
| VAL-010 | Integration evidence for both states: malformed trailing JSONL line exits non-zero and names journal file/line; clean journal replays without corruption error. | Task 4 required commands; phase-end full gate. | Task 4 creates corrupt-tail fixture; Go test package pattern and integration tag syntax validated. |
| VAL-011 | Integration evidence for both states: crafted precondition mismatch stops replay with structured fields `event_id`, `event_type`, `table`, `pk`, `expected`, `actual`, file, and line; matching precondition applies cleanly. | Task 4 required commands; phase-end full gate. | Task 4 creates conflicting-precondition fixture and matching control fixture. |
| VAL-012 | Freshness matrix evidence: matching `projection_meta` serves reads without rebuild; missing meta, schema-version mismatch, position mismatch, and last-event-hash mismatch each trigger catch-up/rebuild before serving. | Task 4 required commands; phase-end full gate. | Task 4 owns `projection_meta` helpers and freshness decision; each mismatch class is independently testable in temp DB/journal fixtures. |
| VAL-013 | Concurrency evidence that two same-checkout writes both persist as distinct events with monotonic sequences; held-lock evidence that a waiting caller exits non-zero naming holder and lock path; lock metadata file-state proof includes PID, hostname, started_at, command, and role when available. | Task 3 required commands with mandatory `-race`; conditional integration command if tagged; phase-end full gate. | Task 3 acceptance criteria include goroutine concurrency and held lock cases; `.fid/write.lock` is a lock directory whose metadata is file-state inspectable. |
| VAL-014 | Unit/integration evidence opening every connection path and querying `PRAGMA foreign_keys` and `PRAGMA recursive_triggers`, both returning ON. | Task 1 required commands; phase-end full gate. | `gofmt`, `golangci-lint`, `go test`, and coverage commands are syntactically valid; Task 1 owns `internal/db/` helper. |
| VAL-015 | Unit evidence for issue ID shape with default 6 hex and supported 8 hex, event ID format `EVT-<UTC YYYYMMDD>-<sequence>`, per-day sequence increments under lock, and forced collision retry succeeds. | Task 2 and Task 3 required commands; phase-end full gate. | Task 2 owns pure ID helpers and collision predicate; Task 3 owns file-backed sequence allocation. |
| VAL-016 | For `fid issue create` and `fid rebuild`, integration evidence for missing `--by`, unknown role, and valid seeded role. Positive-only evidence is insufficient. `fid init` is excluded as bootstrap. | Task 6 required commands; Task 4/5 rebuild command evidence; phase-end full gate. | Task 6 owns minimal policy role existence; seeded roles are created by Task 5. |
| VAL-017 | Import graph evidence that no `internal/...` package imports `cmd/fid`; grep/source-structure evidence that `cmd/fid` avoids persistence/domain operations and delegates to `internal/`; lint/test/build gate evidence. | Phase-end full gate plus `go list`/`grep` architecture checks. | `go help list` validates import graph command; `grep --help` validates source scan; exact checks are in QA-17. |

## 7. Command And Feasibility Validation

| Referenced Step | Validation Method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | Help output and directory existence check. | PASS | `gofmt -h` shows `-l` and path arguments; `ls -d cmd internal` found both dirs. |
| `golangci-lint run ./...` | Help output. | PASS | `golangci-lint run --help` shows `run [flags]`; `golangci-lint --help` lists `run`. |
| `go test -race -coverprofile=coverage.out ./...` | Help output for test/build flags and package patterns; cache location check. | PASS | `go help build` documents `-race`; `go help testflag` documents `-coverprofile`; `go help packages` documents `./...`; `go env` shows writable repo-local caches. |
| `go test -race ./...` | Help output for `go test`, build flags, and package patterns; cache location check. | PASS | `go help test` documents package-list test mode; `go help build` documents `-race`; `go help packages` documents `./...`; `go env` shows `GOCACHE` under repo. |
| `go tool cover -func=coverage.out` | Help output. | PASS | `go tool cover -h` documents `-func` for coverage profile output. The help command exits non-zero after printing usage, but validates supported syntax. |
| `go test -tags integration -race ./...` | Help output for build tags, race flag, test command, and package patterns. | PASS | `go help build` documents `-tags` and `-race`; `go help test` documents test package mode; `go help packages` documents `./...`. |
| `go test -bench . -benchmem ./...` | Help output for benchmark flags and package patterns. | PASS | `go help testflag` documents `-bench` and `-benchmem`; `go help packages` documents `./...`. |
| `go build -o /tmp/fid ./cmd/fid` | Help output and output path feasibility. | PASS | `go help build` documents `-o`; `/tmp` is writable by project policy; `cmd/fid` exists. |
| `go list` import graph checks for VAL-017 | Help output. | PASS | `go help list` documents `-f`, `Imports`, and package patterns. |
| `grep`/file-state checks in deterministic checklist | Help output and existing parent dirs. | PASS | `grep --help` documents fixed/recursive search and exit codes; `find cmd internal migrations testdata` found parent dirs. |
| Phase 1 task specs | File existence check. | PASS | `find project-docs/docs/phases/phase-1/tasks -maxdepth 1 -type f -name 'task-*.md'` found task-1 through task-6. |
| Golden fixtures under `testdata/journal/` | Parent directory and task scope check. | PASS | `testdata` exists now; Task 2 and Task 4 explicitly create `testdata/journal/` fixtures. Final review must verify committed fixtures exist. |
| Completion reports `task-N-completion.md` | Task spec artifact assumption. | PASS | Every task lists its completion report path. Final review must treat missing reports as blocking. |
| Human checklist generation | Template existence and quality plan requirement. | PASS | `ls -l project-docs/templates/human-checklist-template.md project-docs/templates/mid-phase-review-template.md` found both templates; Quality Plan §5.2 requires human checklist generation at phase end. |
| Required command placeholder scan | Text scan. | PASS | `rg "<TBD>|<E2E_COMMAND>|placeholder|N/A"` found policy text and accepted `N/A` release/e2e deferrals, not placeholder task commands. |

## 8. Deterministic Checklist

Items are executed or verified by QA Lead and must be evidenced in
`final-review.md`. PASS without output excerpts, file-state proof, or command
logs is invalid.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Verify Task 1 completion report includes exact required command outcomes and dependency-budget notes for `modernc.org/sqlite` and `golang-migrate/migrate`. `(qa-lead)` | VAL-014, no VAL | Completion report excerpt with all required commands PASS/BLOCKED and license/rationale note. |
| QA-02 | Verify Task 2 completion report includes exact required command outcomes, golden v1 fixture path, canonical JSON/hash tests, upgrader dispatch tests, and ID/collision tests. `(qa-lead)` | VAL-007, VAL-009, VAL-015 | Completion report excerpt plus committed fixture file-state proof. |
| QA-03 | Verify Task 3 completion report includes exact required command outcomes, race-enabled concurrency evidence, held-lock evidence, and sequence allocation evidence. `(qa-lead)` | VAL-007, VAL-013, VAL-015 | Completion report excerpt and test output naming concurrency and held-lock cases. |
| QA-04 | Verify Task 4 completion report includes exact required command outcomes, required build command outcome, replay/checkpoint/corruption/conflict/freshness evidence, and any integration command trigger. `(qa-lead)` | VAL-008, VAL-009, VAL-010, VAL-011, VAL-012 | Completion report excerpt with test output or BLOCKED reason for every acceptance criterion. |
| QA-05 | Verify Task 5 completion report includes exact required command outcomes, init file-state evidence, prefix evidence, one bootstrap event, re-init no-change proof, and delete-DB rebuild proof. `(qa-lead)` | VAL-001, VAL-002, VAL-003, VAL-004, VAL-005 | Completion report excerpt, file-state proof, and command output. |
| QA-06 | Verify Task 6 completion report includes exact required command outcomes, create/get JSON evidence, one-event evidence, rebuild parity evidence, and role validation matrix for `fid issue create` and `fid rebuild`. `(qa-lead)` | VAL-006, VAL-007, VAL-008, VAL-016 | Completion report excerpt and JSON output comparison evidence. |
| QA-07 | Run `gofmt -l cmd internal`; expect no files listed. `(qa-lead)` | no VAL | Literal command output in final review. |
| QA-08 | Run `golangci-lint run ./...`; expect exit 0. `(qa-lead)` | no VAL, VAL-017 support | Literal command output in final review. |
| QA-09 | Run `go test -race ./...`; expect exit 0. `(qa-lead)` | no VAL, VAL support by regression coverage | Literal command output in final review. |
| QA-10 | Run `go test -race -coverprofile=coverage.out ./...`; expect exit 0 and coverage profile produced. `(qa-lead)` | VAL-001..VAL-016 support | Literal command output and `ls -l coverage.out` or equivalent file-state proof. |
| QA-11 | Run `go tool cover -func=coverage.out`; verify total and relevant package coverage do not violate Quality Plan thresholds. `(qa-lead)` | no VAL | Literal coverage output and threshold assessment. |
| QA-12 | Run `go test -tags integration -race ./...`; expect exit 0 or no tagged tests with a clear `go test` success line. `(qa-lead)` | VAL-001, VAL-003, VAL-005, VAL-006, VAL-008, VAL-010, VAL-011, VAL-012, VAL-013, VAL-016 | Literal command output. |
| QA-13 | Run `go test -bench . -benchmem ./...`; expect benchmarks run or packages without benchmarks report success; investigate failures. `(qa-lead)` | VAL-013 support, no VAL | Literal command output and performance spot-check notes. |
| QA-14 | Run `go build -o /tmp/fid ./cmd/fid`; expect exit 0 and `/tmp/fid` exists. `(qa-lead)` | VAL-001, VAL-006, VAL-016 support | Literal command output and file-state proof. |
| QA-15 | Execute an end-to-end temp-workspace smoke path with `/tmp/fid init`, `/tmp/fid rebuild --by developer`, `/tmp/fid issue create --kind task --title <t> --by developer`, and `/tmp/fid issue get <id> --format json`. `(qa-lead)` | VAL-001, VAL-003, VAL-005, VAL-006, VAL-007, VAL-008, VAL-016 | CLI output excerpts, JSON output, journal line counts, and temp workspace path. |
| QA-16 | Verify negative paths in a temp workspace: re-init rejection, unknown role, missing `--by`, corrupt-tail rebuild, and precondition conflict fixture. `(qa-lead)` | VAL-004, VAL-010, VAL-011, VAL-016 | CLI/test output excerpts proving non-zero exits and required error fields. |
| QA-17 | Verify architecture boundary: run an import graph check showing no `internal/...` import includes `cmd/fid`; scan `cmd/fid` for direct persistence/domain operations; review adapters delegate to `internal/`. `(qa-lead)` | VAL-017 | `go list`/`grep` output excerpts plus concise source-review notes with file references. |
| QA-18 | Verify final review records an outcome for every VAL-001 through VAL-017 as PASS, FAIL, BLOCKED, or CARRY-FORWARD. `(qa-lead)` | VAL-001..VAL-017 | Final-review outcome table excerpt. |
| QA-19 | Verify no held-out scenarios or evaluator prompts were read or used in Phase 1 task evidence. `(qa-lead)` | no VAL | Completion report excerpts or absence statement. |

Suggested VAL-017 import check for final review:

```text
go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/...
```

The output must contain no `cmd/fid` import. Use a separate `grep` or `rg`
source scan over `cmd/fid` to support the adapter-boundary review.

## 9. Mid-Phase Review Cadence

- Review 1: after Task 3 is complete, write
  `project-docs/docs/phases/phase-1/qa/review-1.md`.
- Review 2: after Task 6 is complete and before final review, write
  `project-docs/docs/phases/phase-1/qa/review-2.md`.
- Template:
  `project-docs/templates/mid-phase-review-template.md`.

This uses the phase runner's expected three-task cadence. No stricter cadence is
required because Tasks 1-3 form the lower persistence substrate and Tasks 4-6
complete the end-to-end user-command loop. Review 2 is still required even
though it lands at the end of the task list; it verifies Task 4-6 evidence
before QA Lead final review and human checklist generation.

Review 1 focus:

- Command coverage for Tasks 1-3 exactly as listed.
- Dependency-budget note for Task 1.
- v1 event envelope, canonical hash, golden fixture, upgrader identity, ID shape
  and collision retry evidence.
- Lock metadata, held-lock error, race-enabled concurrency, one-line append, and
  sequence evidence.
- Early architecture boundary drift in `cmd/fid` and `internal/`.

Review 2 focus:

- Command coverage for Tasks 4-6 exactly as listed, including required build
  commands and triggered integration commands.
- Replay purity, golden replay, corruption, precondition conflict, checkpoint
  freshness matrix.
- Init workspace, prefix, seeded roles, re-init no-change, delete-DB rebuild.
- Issue create/get JSON, one-event-per-write, role validation for every Phase 1
  write command.
- Documentation updates named by task specs.

## 10. Human Verification Scope

Human checks are judgment-only and do not replace deterministic VAL evidence.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-01 | Run the completed tracer-bullet flow on a disposable local directory and judge whether `fid init`, `fid issue create`, `fid issue get`, and `fid rebuild` feel understandable for first dogfood use. `(human)` | no VAL | Dogfooding notes and approval/requested changes. |
| HUMAN-02 | Judge whether user-facing error messages for re-init, missing/unknown role, corrupt journal, precondition conflict, and held lock are clear and actionable. `(human)` | no VAL | Notes naming any confusing text. |
| HUMAN-03 | Judge whether the minimal JSON issue output is sufficient as the Phase 1 default read shape while TOON/markdown/YAML remain deferred. `(human)` | no VAL | Product-owner judgment. |

## 11. Phase-End Decision Rule

The phase can proceed to final review only when:

- Required deterministic gates pass or are recorded `BLOCKED` with escalation.
- Every task completion report exists and records every exact required command
  as `PASS` or `BLOCKED`.
- Every conditional command trigger is explicitly evaluated in the relevant
  completion report.
- Every `VAL-NNN` has reproducible deterministic evidence and a final-review
  outcome slot.
- Every `(qa-lead)` item has command output, file-state proof, logs, or another
  reproducible evidence excerpt.
- Human checks are human-judgment-only and mapped to the human checklist.
- No unresolved blocking QA findings, missing required artifacts, or placeholder
  required commands remain.
- Required documentation updates from the task specs are present or explicitly
  recorded as blockers.
