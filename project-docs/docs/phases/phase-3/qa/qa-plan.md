# Phase 3 QA Plan - Output Formats + Foundation Hardening

> Phase: 3 - Output formats + foundation hardening
> Refreshed: 2026-06-07
> QA Lead: Codex QA Lead
> Status: READY

## 1. Scope And Source Documents

Phase 3 verifies `fid issue get --format` for JSON/TOON/markdown/YAML and
foundation hardening for rebuild performance, rebuild determinism, corruption
reporting, conflict reporting, and adapter boundaries.

Sources:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/docs/phases/phase-3/qa/validation-contract.md`
- `project-docs/docs/phases/phase-3/decisions.md`
- `project-docs/docs/phases/phase-3/tasks/task-1.md`
- `project-docs/docs/phases/phase-3/tasks/task-2.md`
- `project-docs/docs/phases/phase-3/tasks/task-3.md`
- `project-docs/docs/phases/phase-3/tasks/task-4.md`
- `project-docs/docs/phases/phase-3/tasks/task-5.md`

## 2. Validation Contract Refinement Check

- All in-scope PRD/TRD/design requirements referenced by at least one
  `VAL-NNN`: YES.
- Every `VAL-NNN` has a concrete deterministic evidence method: YES.
- Every `VAL-NNN` maps to a task/window or phase-end check: YES.
- Every `VAL-NNN` describes one observable, or the QA plan fans it out into
  one-observable checks: YES.
- Conditional or binary `VAL-NNN` assertions require both qualifying and
  non-qualifying states: YES.
- Parity or "matches reference" assertions require deterministic conformance
  proof: YES.
- Revisions requested from Tech Lead: none.

Requirement coverage:

| Requirement | Covered By |
|---|---|
| PRD-OUT-002; TRD-ARCH-012; TRD §5 TOON parity | VAL-001, VAL-002 |
| PRD-OUT-003; `06_DESIGN_PLAN_CLI.md` §21 markdown layout | VAL-003 |
| PRD-OUT-003a; DEC-012 YAML dependency decision | VAL-004 |
| PRD-OUT-004; format dispatch and rejection | VAL-005 |
| TRD-PERF-002; PRD-OPS-004 rebuild command | VAL-006 |
| TRD-ARCH-001; PRD-OPS-003/004 deterministic projection | VAL-007 |
| PRD-OPS-003/004; TRD-ARCH-002 corruption path | VAL-008 |
| PRD-OPS-012; TRD-ARCH-002/003 conflict path | VAL-009 |
| TRD-ARCH-005 adapter boundary | VAL-010 |

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | Generic render pipeline + in-tree TOON | VAL-001, VAL-002, VAL-005 | Medium | Verify parity, nested arrays, unknown format, and no TOON module dependency. |
| Task 2 | Markdown renderer | VAL-003, VAL-005 | Medium | Verify populated and empty markdown states against §21 and dispatch continuity. |
| Task 3 | YAML output | VAL-004, VAL-005 | Medium | Verify YAML decodes equal to JSON and dependency rationale is recorded. |
| Task 4 | Rebuild benchmark + determinism | VAL-006, VAL-007 | High | Capture elapsed benchmark output for ~1k and 10k; verify scale determinism. |
| Task 5 | Conflict/corruption error fidelity | VAL-008, VAL-009 | High | Verify both `fid rebuild` and read catch-up paths, plus clean states. |
| Phase end | Adapter boundary and full gate | VAL-010, no VAL | Medium | Verify import graph/source placement and run phase-level gate. |

## 4. Required Evidence Strategy

| Evidence ID | Required evidence | Covers | Owner | Produced By |
|---|---|---|---|---|
| EVID-01 | `internal/format` TOON parity tests for empty-arrays and populated/nested issue; golden TOON fixture or structural cross-check; `cmd/fid` TOON smoke. | VAL-001 | Developer, QA Lead | Task 1 report, final review |
| EVID-02 | `go list -m all` or `go.mod`/`go.sum` proof showing no third-party TOON dependency; `go build -o /tmp/fid ./cmd/fid` PASS. | VAL-002 | Developer, QA Lead | Task 1 report, final review |
| EVID-03 | Golden markdown renderer tests and `cmd/fid` markdown tests for populated and empty AC/label/link states. | VAL-003 | Developer, QA Lead | Task 2 report, final review |
| EVID-04 | YAML output decoded with `gopkg.in/yaml.v3` deep-equals JSON-decoded output for empty-arrays and populated/nested issues; dependency rationale recorded. | VAL-004 | Developer, QA Lead | Task 3 report, final review |
| EVID-05 | `cmd/fid` dispatch table covers `json`, `toon`, `markdown`, `yaml`, and unsupported `xml` non-zero error naming known formats. | VAL-005 | Developer, QA Reviewer, QA Lead | Tasks 1-3 reports, review after task 3, final review |
| EVID-06 | Benchmark output from `go test -bench . -benchmem ./...` includes rebuild cases for ~1k and 10k events; record elapsed result for each; assert ~1k < 1s and 10k < 5s. | VAL-006 | Developer, QA Lead | Task 4 report, final review |
| EVID-07 | `internal/journal` scale determinism test proves event count equals seed count, delete-DB + rebuild projection equality, and rebuild-vs-rebuild equality. | VAL-007 | Developer, QA Lead | Task 4 report, final review |
| EVID-08 | Corrupt-tail tests prove structured corruption on `fid rebuild`, structured corruption on `fid issue get` read catch-up, and clean journal success on both commands. | VAL-008 | Developer, QA Lead | Task 5 report, final review |
| EVID-09 | Precondition-conflict tests prove structured conflict on `fid rebuild`, structured conflict on `fid issue get` read catch-up before serving reads, and satisfiable journal success on both commands. | VAL-009 | Developer, QA Lead | Task 5 report, final review |
| EVID-10 | Import graph/source placement proof: no `internal/...` package imports `cmd/fid`; render logic is in `internal/format`, with `cmd/fid` limited to adapter selection/output. | VAL-010 | QA Lead | Final review |
| EVID-11 | Full task and phase gate command output, including coverage threshold evidence for `internal/journal` >= 80% and `cmd/fid` >= 70%. | no VAL | Developer, QA Reviewer, QA Lead | Completion reports, review after task 3, final review |

## 5. Deterministic Gates

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | No files listed | Every task and phase end |
| Lint/static analysis | `golangci-lint run ./...` | PASS | Every task and phase end |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS | Every task and phase end |
| Coverage report | `go tool cover -func=coverage.out` | PASS; `internal/journal` >= 80% where touched; `cmd/fid` >= 70% | Every task and phase end |
| Benchmark | `go test -bench . -benchmem ./...` | PASS; rebuild ~1k elapsed < 1s and 10k elapsed < 5s | Task 4 and phase end |
| Build | `go build -o /tmp/fid ./cmd/fid` | PASS | Tasks 1, 2, 3, 5 and phase end |
| Integration tag check | `go test -tags integration -race ./...` | PASS or explicit N/A evidence if no integration-tagged tests exist | Phase end |

Task completion reports must record every required task-spec command as `PASS`
or `BLOCKED` with reason. A missing, silently skipped, or narrower substituted
command is a blocking QA finding.

## 6. Command And Feasibility Validation

| Referenced Step | Validation Method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | Help output and path existence check. | PASS | `gofmt -h` shows `-l` and path args; `cmd/` and `internal/` exist. |
| `golangci-lint run ./...` | Help/version output. | PASS | `golangci-lint --version` => v2.12.2; `golangci-lint run --help` documents `run`. |
| `go test -race -coverprofile=coverage.out ./...` | Go test help output. | PASS | `go help testflag` documents `-coverprofile`; `go help build` documents `-race`. |
| `go tool cover -func=coverage.out` | Tool help output. | PASS | `go tool cover -h` documents `-func`. The help command exits 2 after printing usage; syntax is documented. |
| `go test -bench . -benchmem ./...` | Go test flag help output. | PASS | `go help testflag` documents `-bench .` and `-benchmem`. |
| `go build -o /tmp/fid ./cmd/fid` | Build help output and package existence check. | PASS | `go help build` documents `-o`; `go list ./cmd/fid` resolves `github.com/mbell-nzl/fid/cmd/fid`. |
| `go test -tags integration -race ./...` | Build flag help output. | PASS | `go help build` documents `-tags`; phase brief states no integration build tag is currently in use, so PASS or documented N/A is acceptable at phase end. |
| Phase packages | `go list` package existence check. | PASS | `go list ./internal/journal ./internal/format ./cmd/fid` resolves all packages. |
| Task specs and validation contract | File-state check. | PASS | `find project-docs/docs/phases/phase-3/...` shows task-1..task-5 and `qa/validation-contract.md`. |
| Benchmark fixture assumption | Task-spec feasibility check. | PASS | Task 4 requires creation of `internal/journal/rebuild_bench_test.go` and seeded valid journals; QA will fail VAL-006 if ~1k and 10k elapsed benchmark lines are missing. |
| Corrupt/conflict fixture assumption | Task-spec feasibility check. | PASS | Task 5 requires `cmd/fid` tests for corrupt journal, conflicting journal, and clean journal on both `fid rebuild` and `fid issue get`; QA will fail VAL-008/009 if either path or clean state is absent. |
| Adapter-boundary import check | Package listing feasibility check. | PASS | `go list ./internal/...` resolves internal packages; final QA uses `go list -deps ./internal/...` and fails if `github.com/mbell-nzl/fid/cmd/fid` appears. |
| Adapter-boundary source placement check | File-state feasibility check. | PASS | `find internal/format -maxdepth 1 -type f` and `rg -n "toon|markdown|yaml|render|format" internal/format cmd/fid` are available file-state/source checks. |
| Human checklist feasibility | Scope review. | PASS | Human checks are limited to dogfooding and message clarity judgment, with deterministic checks carrying VAL proof. |

No referenced command is a placeholder.

## 7. Deterministic Checklist - QA Lead

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Verify Task 1 report includes passing TOON parity tests for empty-arrays and populated/nested issue and a `cmd/fid` TOON smoke output excerpt. `(qa-lead)` | VAL-001 | Task report test output, golden/structural parity test names, CLI excerpt. |
| QA-02 | Verify TOON has no third-party dependency and the binary builds. `(qa-lead)` | VAL-002 | `go list -m all` or `go.mod`/`go.sum` excerpt; `go build -o /tmp/fid ./cmd/fid` output. |
| QA-03 | Verify Task 2 report includes golden markdown tests for populated and empty AC/label/link states and `cmd/fid issue get --format markdown` evidence. `(qa-lead)` | VAL-003 | Test names/output and CLI excerpt or golden fixture path. |
| QA-04 | Verify Task 3 report includes YAML decode-vs-JSON parity tests for empty-arrays and populated/nested issues, plus DEC-012 dependency-budget note. `(qa-lead)` | VAL-004 | Test output, `go.mod`/`go.sum` excerpt, completion-report dependency rationale. |
| QA-05 | Verify final dispatch coverage for `json`, `toon`, `markdown`, `yaml`, and unsupported `xml`, with unknown format rejected non-zero and no JSON fallback. `(qa-lead)` | VAL-005 | `cmd/fid` table-test output and stderr excerpt for `xml`. |
| QA-06 | Run or verify `go test -bench . -benchmem ./...`; capture rebuild elapsed output for ~1k events and assert elapsed/op or reported elapsed < 1s. `(qa-lead)` | VAL-006 | Benchmark output line naming the ~1k case and measured elapsed result. |
| QA-07 | Run or verify `go test -bench . -benchmem ./...`; capture rebuild elapsed output for 10k events and assert elapsed/op or reported elapsed < 5s. `(qa-lead)` | VAL-006 | Benchmark output line naming the 10k case and measured elapsed result. |
| QA-08 | Verify scale rebuild determinism test proves event count, delete-DB + rebuild equality, and rebuild-vs-rebuild equality. `(qa-lead)` | VAL-007 | `internal/journal` test output and assertion summary from task report/final run. |
| QA-09 | Verify malformed trailing journal line reports structured corruption on `fid rebuild` and clean rebuild succeeds. `(qa-lead)` | VAL-008 | `cmd/fid` test output and stderr excerpt containing file:line, parse cause, and manual repair guidance; clean-state success output. |
| QA-10 | Verify malformed trailing journal line reports structured corruption on `fid issue get` read catch-up and clean read succeeds. `(qa-lead)` | VAL-008 | `cmd/fid` test output and stderr/stdout excerpts proving non-zero corrupt read and successful clean read. |
| QA-11 | Verify failed precondition reports structured conflict on `fid rebuild` and satisfiable rebuild succeeds. `(qa-lead)` | VAL-009 | `cmd/fid` test output and stderr excerpt naming event id/type, entity table+pk, expected, actual; clean-state success output. |
| QA-12 | Verify failed precondition reports structured conflict on `fid issue get` read catch-up before serving reads, and satisfiable read succeeds. `(qa-lead)` | VAL-009 | `cmd/fid` test output and stderr/stdout excerpts proving non-zero conflict read and successful clean read. |
| QA-13 | Verify adapter boundary: no `internal/...` package imports `cmd/fid`; rendering logic lives in `internal/format`. `(qa-lead)` | VAL-010 | `go list -deps ./internal/...` output with no `github.com/mbell-nzl/fid/cmd/fid`; `find internal/format -maxdepth 1 -type f`; `rg -n "toon|markdown|yaml|render|format" internal/format cmd/fid` excerpt showing renderer files under `internal/format` and adapter-only routing in `cmd/fid`. |
| QA-14 | Run full phase gate. `(qa-lead)` | no VAL | Output for format, lint, coverage, coverage report, integration/N/A, benchmarks, and build commands. |
| QA-15 | Compare every completion-report command with the exact task-spec command list. `(qa-lead)` | no VAL | Table in final review marking each command `PASS`, `FAIL`, or `BLOCKED`; skipped/narrower commands flagged blocking. |

## 8. Mid-Phase Review Cadence - QA Reviewer

Default cadence applies: one QA Reviewer pass after every 3 completed tasks.
Phase 3 has five tasks, so the mid-phase review occurs after Task 3 and covers
the complete output-format track before hardening tasks begin. No stricter
cadence is required because Tasks 4 and 5 are independent and are covered by
task-specific reports plus final QA Lead execution.

Review artifact: `project-docs/docs/phases/phase-3/qa/review-3.md`, using
`project-docs/templates/mid-phase-review-template.md`.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QR-01 | Verify Task 1-3 completion reports cite VAL coverage and contain all exact task-spec command outcomes. `(qa_reviewer)` | VAL-001, VAL-002, VAL-003, VAL-004, VAL-005 | Review table with command/result comparison per task. |
| QR-02 | Verify output-format tests cover both required states: TOON empty/populated, markdown populated/empty, YAML empty/populated. `(qa_reviewer)` | VAL-001, VAL-003, VAL-004 | Test output excerpts and fixture/test names. |
| QR-03 | Verify dispatch is cumulative after Task 3: `json`, `toon`, `markdown`, `yaml`, and unsupported `xml` are all covered. `(qa_reviewer)` | VAL-005 | `cmd/fid` test output or completion-report evidence. |
| QR-04 | Verify no unapproved dependency was introduced beyond DEC-012 `gopkg.in/yaml.v3`, and no third-party TOON dependency exists. `(qa_reviewer)` | VAL-002, VAL-004 | `go.mod`/`go.sum` or module graph excerpt. |

QA Reviewer escalates to `tl_escalations_current.md` only for defects matching
the QA Reviewer escalation triggers; otherwise findings stay in `review-3.md`.

## 9. Human Verification Scope

Human checks are not required to prove VAL-001..010; deterministic QA carries
the validation contract. Human items are reserved for product judgment and are
copied to the phase-end human checklist.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-01 | Dogfood `fid issue get` across JSON, TOON, markdown, and YAML on a real local issue and judge whether the outputs are understandable for agent-assisted work. `(human)` | no VAL | Notes naming issue id, formats tried, and any confusing output. |
| HUMAN-02 | Review corruption and conflict error messages from QA excerpts and judge whether the corrective action is clear enough for manual repair. `(human)` | no VAL | Approval note or requested wording changes. |

## 10. Phase-End Decision Rule

The phase can proceed to final review only when:

- Every Task 1-5 completion report exists and maps evidence to its `VAL-NNN`
  IDs.
- The mid-phase QA Reviewer artifact after Task 3 exists or is explicitly
  `BLOCKED` with reason.
- Required deterministic gates pass or are recorded `BLOCKED` with escalation.
- Every `VAL-001..VAL-010` has a final-review outcome slot: `PASS`, `FAIL`,
  `BLOCKED`, or `CARRY-FORWARD`.
- `PASS` on VAL-006 includes captured elapsed benchmark results for ~1k and 10k
  events and the thresholds ~1k < 1s, 10k < 5s.
- `PASS` on VAL-008 and VAL-009 includes evidence for both `fid rebuild` and
  `fid issue get` read catch-up, plus clean/satisfiable states.
- Every `(qa-lead)` item has reproducible evidence in `final-review.md`.
- Human checks are human-judgment-only and mapped to
  `human-checklist.md`.
- No unresolved blocking defects remain.
