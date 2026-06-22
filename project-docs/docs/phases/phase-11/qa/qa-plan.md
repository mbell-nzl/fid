# Phase 11 QA Plan - Claims Layer

> Phase: 11 - Claims layer (roadmap v0.1d, slice 2 of 4)
> Refreshed: 2026-06-21
> QA Lead: Codex QA Lead
>
> Location:
> `project-docs/docs/phases/phase-11/qa/qa-plan.md`.

## 1. Scope And Source Documents

Phase 11 delivers primitive issue claims and transition integration: claim
acquire, renew, force-acquire, release, list, and started-state transition
gating. The QA plan treats the validation contract as the verification target
and requires per-task QA review before the next ordered task proceeds.

Source artifacts loaded:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/docs/current/decisions_principles.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/01_PRD.md`
- `project-docs/specs/02_TRD.md`
- `project-docs/specs/03_TRD_SCHEMA.md`
- `project-docs/specs/05_DESIGN_PLAN.md`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-11/decisions.md`
- `project-docs/docs/phases/phase-11/progress.md`
- `project-docs/docs/phases/phase-11/qa/validation-contract.md`
- `project-docs/docs/phases/phase-11/tasks/task-1.md`
- `project-docs/docs/phases/phase-11/tasks/task-2.md`
- `project-docs/docs/phases/phase-11/tasks/task-3.md`
- `project-docs/docs/phases/phase-11/tasks/task-4.md`
- `project-docs/docs/phases/phase-11/tasks/task-5.md`
- `project-docs/templates/qa-plan-template.md`

Open escalation baseline: ESC-008 resolved by DEC-028 with do-now Task 5 before
final QA. Branch feasibility note: active shell branch is `phase-11`; stale
current-doc branch/status text is reconciled by the ESC-008 planning update.

## 2. Validation Contract Refinement Check

- All in-scope requirements referenced by at least one `VAL-NNN`: YES.
- Every `VAL-NNN` has one validator role: YES.
- Every `VAL-NNN` has a concrete evidence method: YES.
- Every `VAL-NNN` maps to a task/window or phase-end check: YES.
- Every `VAL-NNN` describes one observable, or human-checklist fan-out is
  explicitly required: YES. VAL-014 is human-only and is fanned out below.
- Conditional or binary `VAL-NNN` assertions require both qualifying and
  non-qualifying states: YES for VAL-003, VAL-005, VAL-006, VAL-007, VAL-009,
  VAL-011, VAL-012, and VAL-015.
- Parity or "matches reference" `VAL-NNN` assertions have deterministic
  conformance proof when available: YES for VAL-001 and VAL-008.
- Tech Lead revision after ESC-008: VAL-015 and Task 5 added; no further
  validation-contract revision is blocking.

Refined validation rows:

| VAL | Validator | Evidence method | Task/window | Refinement |
|---|---|---|---|---|
| VAL-001 | QA Lead | Unit/CLI evidence for matrix row uniqueness, allow/deny checks, and `fid policy show` parity with `Authorize`. | Task 1, reviewed in Review 1. | Compare every implemented `issue.claim.*` row against `03_TRD_SCHEMA.md` section 13; no code-reading-only pass. |
| VAL-002 | QA Lead | Domain and CLI output plus DB/read-back proof for default TTL, custom TTL, note, actor, and single current row. | Tasks 1-2, Reviews 1-2. | Require default and override TTL evidence; positive-only basic acquire does not prove note or custom TTL. |
| VAL-003 | QA Lead | CLI conflict attempt with before/after journal and DB counts plus stderr holder/expiry. | Tasks 1-2, Reviews 1-2. | Must prove both live-holder success setup and non-holder conflict rejection with zero append. |
| VAL-004 | QA Lead | Expired-claim fixture, normal acquire by new role, current-row proof, and event/activity superseded-holder evidence. | Tasks 1-2, Reviews 1-2. | Must distinguish stale takeover from force-acquire. |
| VAL-005 | QA Lead | Permitted force success and denied force rejection, with role/action permission error and zero append for denial. | Tasks 1-2, Reviews 1-2. | Must prove both permitted and denied states for `issue.claim.force_acquire`. |
| VAL-006 | QA Lead | Holder renewal with and without TTL, `activity_log` row, last event proof, and non-holder zero-append rejection. | Tasks 1-2, Reviews 1-2. | Must prove prior TTL preservation and supplied TTL override separately. |
| VAL-007 | QA Lead | CLI release success plus wrong-holder, repeat, missing role, and unknown-role reject paths with zero append. | Task 2, Review 2. | Missing/unknown role errors must follow public role policy, not silently no-op. |
| VAL-008 | QA Lead | `fid claim list` in json/toon/markdown/yaml parsed or structurally compared against JSON reference. | Task 3, Review 3. | Format parity must include all required fields and `--by` filter behavior. |
| VAL-009 | QA Lead | Live and stale fixture, `--stale` membership proof, and rowid-stable ordering before/after rebuild. | Task 3, Review 3. | Must prove live exclusion and stale inclusion; ordering cannot be inferred from code. |
| VAL-010 | QA Lead | Human claim read-back/list proof and non-holder transition reject proof. | Tasks 2-4, Reviews 2 and 4. | Must verify `human` as stored/displayed role and as blocking holder. |
| VAL-011 | QA Lead | Transition cases for no claim, wrong holder, stale holder, and matching holder with before/after journal/status proof. | Task 4, Review 4. | Must prove rejection preserves status and appends nothing. |
| VAL-012 | QA Lead | Transition out of started with one combined event, status change, release effect, and no current claim after projection. | Task 4, Review 4. | Must prove same-event atomicity; separate status and release events fail. |
| VAL-013 | QA Lead | Integration test using concurrent acquire attempts against temp SQLite workspace. | Task 1, Review 1. | Must show exactly one success, one current row, and holder/expiry conflict errors for losers. |
| VAL-014 | Human | Phase-end human dogfood checklist with one observable per row. | Phase-end. | QA Lead generates the human checklist after deterministic QA; QA Lead cannot pass this by code review. |
| VAL-015 | QA Lead | Stale-projection fixtures for claim acquire, renew, and release, plus fresh-projection and concurrent-acquire regression evidence. | Task 5, Review 5. | Must prove no nested-lock error, catch-up before validation, one-event successful writes, zero-append rejects, and unchanged public claim behavior. |

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | Claim matrix rows and domain semantics | VAL-001, VAL-002, VAL-003, VAL-004, VAL-005, VAL-006, VAL-013 | High | New package, permission matrix growth, journaled writes, validate-before-append, and concurrency proof. Review before Task 2 starts. |
| Task 2 | `fid issue claim` and `fid issue unclaim` | VAL-002, VAL-003, VAL-004, VAL-005, VAL-006, VAL-007, VAL-010 | High | Public CLI write commands. Review must compare task commands to completion report and inspect zero-append evidence. |
| Task 3 | `fid claim list` | VAL-008, VAL-009, VAL-010 | Medium | Public read command with format parity and rowid ordering. Review must verify read-only behavior and `--by` filter semantics. |
| Task 4 | Claim and transition integration | VAL-010, VAL-011, VAL-012 | High | Cross-subsystem write path. Review as its own window before final QA. |
| Task 5 | Claim write projection freshness lock ordering | VAL-013, VAL-015 | High | Do-now ESC-008 fix. Review must verify stale-projection evidence and exact benchmark/build gates before final QA. |

## 4. Required Evidence Strategy

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | Matrix row count, exact role booleans, `Authorize` allow/deny test output, and `fid policy show` parity output. | VAL-001 | Developer, QA Reviewer, QA Lead | Task 1 / Review 1 / final review |
| EVID-02 | Claim acquire command/domain output with default TTL 1800, custom TTL, note, holder, expiry, and one current row. | VAL-002 | Developer, QA Reviewer, QA Lead | Tasks 1-2 / Reviews 1-2 / final review |
| EVID-03 | Conflict reject stderr plus journal line count and DB row count before/after. | VAL-003 | Developer, QA Reviewer, QA Lead | Tasks 1-2 / Reviews 1-2 / final review |
| EVID-04 | Expired claim takeover evidence with prior holder/expiry recorded and new holder current. | VAL-004 | Developer, QA Reviewer, QA Lead | Tasks 1-2 / Reviews 1-2 / final review |
| EVID-05 | Force-acquire permitted success and denied failure naming `issue.claim.force_acquire`, with zero append on denial. | VAL-005 | Developer, QA Reviewer, QA Lead | Tasks 1-2 / Reviews 1-2 / final review |
| EVID-06 | Renewal row/event/activity evidence for holder renewal and non-holder zero-append rejection. | VAL-006 | Developer, QA Reviewer, QA Lead | Tasks 1-2 / Reviews 1-2 / final review |
| EVID-07 | `issue unclaim` holder release plus repeat/wrong/missing/unknown role rejection evidence. | VAL-007 | Developer, QA Reviewer, QA Lead | Task 2 / Review 2 / final review |
| EVID-08 | Parsed json/toon/markdown/yaml `claim list` outputs showing required fields and `--by` filter parity. | VAL-008 | Developer, QA Reviewer, QA Lead | Task 3 / Review 3 / final review |
| EVID-09 | `claim list --stale` membership proof and rowid order proof before/after rebuild. | VAL-009 | Developer, QA Reviewer, QA Lead | Task 3 / Review 3 / final review |
| EVID-10 | Human claim shown as `human` and developer transition blocked while human holds the live claim. | VAL-010 | Developer, QA Reviewer, QA Lead | Tasks 2-4 / Reviews 2 and 4 / final review |
| EVID-11 | Transition into started no-claim, wrong-holder, stale-holder, matching-holder cases with journal/status before/after. | VAL-011 | Developer, QA Reviewer, QA Lead | Task 4 / Review 4 / final review |
| EVID-12 | Transition out of started appends one combined status+claim-release event and leaves no current claim. | VAL-012 | Developer, QA Reviewer, QA Lead | Task 4 / Review 4 / final review |
| EVID-13 | Concurrent acquire integration test output proving one winner and loser holder/expiry errors. | VAL-013 | Developer, QA Reviewer, QA Lead | Task 1 / Review 1 / final review |
| EVID-14 | Human dogfood notes for holder, expiry, release, and block-reason comprehension. | VAL-014 | Human | Phase-end human checklist |
| EVID-15 | Exact task command outputs and completion-report command/result parity. | no VAL | Developer, QA Reviewer, QA Lead | Every task / every review / final review |
| EVID-16 | Coverage package percentages for `internal/issues`, `internal/claims`, `internal/workflow`, `internal/policy`, `internal/journal`, and `cmd/fid`. | no VAL | Developer, QA Reviewer, QA Lead | Every task / final review |
| EVID-17 | Stale-projection claim acquire/renew/release proof showing catch-up before validation, no nested-lock error, one-event successful writes, zero-append reject behavior, and unchanged concurrent acquire atomicity. | VAL-013, VAL-015 | Developer, QA Reviewer, QA Lead | Task 5 / Review 5 / final review |

## 5. Deterministic Gates

If any required command becomes a placeholder or cannot be validated when
execution is requested, mark the affected review BLOCKED and escalate.

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | No files listed | Every task and final review |
| Lint/static analysis | `golangci-lint run ./...` | PASS | Every task and final review |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS; writes `coverage.out` | Every task and final review |
| Coverage report | `go tool cover -func=coverage.out` | PASS; COV-001/COV-002 package percentages recorded | Every task and final review after coverage profile exists |
| Integration | `go test -tags integration -race ./...` | PASS | Every task and final review |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS | Phase-end; task-level only when task spec condition applies |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | Tasks 2-5 and phase-end; Task 1 only if its conditional trigger applies |
| End-to-end | N/A per `AGENTS.md` and `08_QUALITY_PLAN.md` until scenario harness is accepted | N/A recorded | Phase 11 |
| Full phase check | All required phase commands above | PASS or BLOCKED with reason | Phase-end final review |

Task command notes:

| Task | Required commands | Conditional commands |
|---|---|---|
| Task 1 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...` | `go build -o /tmp/fid ./cmd/fid` if command registration or CLI adapter is unexpectedly touched; `go test -bench . -benchmem ./...` if benchmark fixtures or journal scale helpers are touched. |
| Task 2 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go build -o /tmp/fid ./cmd/fid` | `go test -bench . -benchmem ./...` if benchmark fixtures or journal scale helpers are touched. |
| Task 3 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go build -o /tmp/fid ./cmd/fid` | `go test -bench . -benchmem ./...` if benchmark fixtures or list performance paths are touched. |
| Task 4 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go build -o /tmp/fid ./cmd/fid` | `go test -bench . -benchmem ./...` if transition benchmark fixtures or journal scale helpers are touched. |
| Task 5 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go test -bench . -benchmem ./...`; `go build -o /tmp/fid ./cmd/fid` | None. |

## 6. Command And Feasibility Validation

| Referenced Step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | Help output and file-state check. | PASS | `gofmt -h` shows `-l` and path operands; `find cmd internal ...` confirms both directories exist. |
| `golangci-lint run ./...` | Help output. | PASS | `golangci-lint run --help` shows `Usage: golangci-lint run [flags]`; command is installed. |
| `go test -race -coverprofile=coverage.out ./...` | Go test/build flag help. | PASS | `go help test` supports package list mode; `go help testflag` shows `-coverprofile`; `go help build` shows `-race`. |
| `go tool cover -func=coverage.out` | Tool usage output. | PASS | `go tool cover -h` shows `-func string output coverage profile information for each function`. The help command exits 2 after printing usage, which is acceptable for syntax validation. |
| `go test -tags integration -race ./...` | Go build/test help. | PASS | `go help build` shows `-tags`; `go help test` supports package list mode; `go help build` shows `-race`. |
| `go test -bench . -benchmem ./...` | Go test flag help. | PASS | `go help testflag` shows `-bench regexp` and `-benchmem`. |
| `go build -o /tmp/fid ./cmd/fid` | Go build help and writable-output assumption. | PASS | `go help build` shows `-o output`; sandbox writable roots include `/tmp`; `find cmd/fid ...` confirms `cmd/fid/main.go` exists. |
| End-to-end command | Project documentation validation. | PASS as N/A | `AGENTS.md` and `08_QUALITY_PLAN.md` state end-to-end/scenario validation is N/A until a harness is accepted. No Phase 11 task references an E2E command. |
| Placeholder command scan | Text scan of required command sources. | PASS | `rg "<TBD>|<E2E_COMMAND>|placeholder|TODO" AGENTS.md project-docs/specs/08_QUALITY_PLAN.md project-docs/docs/phases/phase-11/tasks project-docs/docs/phases/phase-11/qa/validation-contract.md` found no placeholder command in Phase 11 task specs or contract; AGENTS only contains the standing placeholder rule. |
| `rg ... project-docs/docs/phases/phase-11/tasks/task-*-done.md` completion-report scan | Ripgrep help output and glob feasibility. | PASS | `rg --help` shows `rg [OPTIONS] PATTERN [PATH ...]`; completion reports are expected task artifacts and may be absent until task completion. |
| Branch assumption `phase-11` | Git branch check. | PASS | `git branch --show-current` returned `phase-11`; current docs now name `phase-11`. |
| Phase QA artifact directory | File-state check. | PASS | `ls project-docs/docs/phases/phase-11/qa` showed `validation-contract.md`, `qa-plan.md`, and Reviews 1-3 at ESC-008 resolution time. |
| Task specs ordered 1-5 | File-state and content check. | PASS | `rg --files project-docs/docs/phases/phase-11/tasks` found `task-1.md` through `task-5.md`; phase brief task index lists dependencies in order. |
| Completion report paths | Template and task-spec check. | PASS | Task specs require `project-docs/docs/phases/phase-11/tasks/task-1-done.md` through `task-5-done.md`; Task 5 report is expected only after Task 5 completion. |
| `project-docs/templates/mid-phase-review-template.md` | File-state check. | PASS | Template exists in `project-docs/templates/`. |
| `project-docs/templates/human-checklist-template.md` | File-state check. | PASS | Template exists in `project-docs/templates/`. |
| `migrations/0002_v01_tables.up.sql` and `migrations/0005_v01_ref_triggers.up.sql` | File-state check. | PASS | Both files were present in `find migrations ...`; they support Task 1 context for existing claims table/triggers. |
| Existing command/adapter paths | File-state check. | PASS | `cmd/fid/issue.go`, `cmd/fid/root.go`, `cmd/fid/policy.go`, and existing command tests are present. |
| Existing domain/persistence paths | File-state check. | PASS | `internal/issues/transition.go`, `internal/journal/project_effects.go`, `internal/policy/matrix.go`, and related tests are present. |
| `internal/claims/` package | Task-scope validation. | PASS | Package exists from Task 1 and is intentionally touched again by Task 5. |
| Spec requirement references | Targeted spec scan. | PASS | `rg "PRD-CLAIM|TRD-ARCH-010|TRD-ARCH-011|PRD-ERR-003|PRD-ERR-006" ...` found the referenced PRD/TRD rows. |
| Claim CLI/schema references | Targeted spec scan. | PASS | `rg "issue.claim|claims|claim list|issue claim|issue unclaim" ...` found CLI and schema references. |
| Human checklist feasibility | Template and phase-end process validation. | PASS | `08_QUALITY_PLAN.md` sections 5.2, 5.3, and 5.5 require human checklist generation and reserve judgment/dogfooding for humans; template exists. |
| Human dogfood fixture assumption | Scope validation. | PASS | Human uses a real `.fid` workspace at phase-end; no held-out scenario material is required or loaded. |
| `git diff --name-only main...HEAD -- migrations` migration drift check | Git diff help, base-branch existence check, and exact dry-run. | PASS | `git diff -h` shows `--name-only` and `<commit>...<commit>` syntax; `git rev-parse --verify main^{commit}` returned `d78bea416de374ac943f974d00e1620c4d47f7af`; `git diff --name-only main...HEAD -- migrations` ran successfully with no output at plan time. |
| `git log --oneline -- project-docs/docs/phases/phase-11/tasks/task-*-done.md` completion-report commit check | Git log help output and oneline dry-run. | PASS | `git log -h` shows pathspecs are accepted after `--`; `git log --oneline -1` returned `d78bea4 Phase 11 planned`; `git log --oneline -- project-docs/docs/phases/phase-11/tasks/task-1.md` ran successfully. |
| `git rev-parse --verify HASH_FROM_COMPLETION_REPORT^{commit}` report hash check | Git rev-parse help and main-commit verification. | PASS | `git rev-parse -h` shows `git rev-parse [<options>] [<arg>...]`; `git rev-parse --verify main^{commit}` successfully verified a commit object. |

## 7. Deterministic Checklist

Every item is executed or verified by QA Lead at final review unless explicitly
covered by a mid-phase review. PASS without cited evidence is invalid.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Verify every required task command appears in the matching completion report as `PASS` or `BLOCKED`, with no narrower substitutions `(qa-lead)` | no VAL | `rg "gofmt -l cmd internal|golangci-lint run ./...|go test -race -coverprofile=coverage.out ./...|go tool cover -func=coverage.out|go test -tags integration -race ./...|go build -o /tmp/fid ./cmd/fid|go test -bench . -benchmem ./..." project-docs/docs/phases/phase-11/tasks/task-*-done.md` plus report excerpts. |
| QA-02 | Run final format, lint, coverage, integration, benchmark, and build gates exactly as listed in Section 5 `(qa-lead)` | no VAL | Literal command outputs for `gofmt -l cmd internal`, `golangci-lint run ./...`, `go test -race -coverprofile=coverage.out ./...`, `go tool cover -func=coverage.out`, `go test -tags integration -race ./...`, `go test -bench . -benchmem ./...`, and `go build -o /tmp/fid ./cmd/fid`. |
| QA-03 | Confirm coverage thresholds remain met for domain packages and `cmd/fid` `(qa-lead)` | no VAL | `go tool cover -func=coverage.out` excerpts for `internal/issues`, `internal/claims`, `internal/workflow`, `internal/policy`, `internal/journal`, and `cmd/fid`; compare to COV-001/COV-002. |
| QA-04 | Verify matrix rows for `issue.claim.acquire`, `issue.claim.renew`, `issue.claim.force_acquire`, and `issue.claim.release` are unique, match `03_TRD_SCHEMA.md`, and are shown/enforced consistently `(qa-lead)` | VAL-001 | Test output plus `fid policy show` output excerpt and allow/deny evidence from completion reports or final fixture. |
| QA-05 | Verify normal claim acquire default TTL, custom TTL, note persistence, expiry, actor, and single current row `(qa-lead)` | VAL-002 | CLI/domain test output and DB or read-back excerpts for default and custom TTL cases. |
| QA-06 | Verify live conflict rejection names holder/expiry and appends/projects nothing `(qa-lead)` | VAL-003 | Stderr excerpt plus before/after journal line count and DB row count. |
| QA-07 | Verify stale takeover succeeds without force and records prior holder/expiry `(qa-lead)` | VAL-004 | Fixture output showing expired prior claim, successful normal acquire, superseded evidence, and new current holder. |
| QA-08 | Verify force-acquire permitted and denied states `(qa-lead)` | VAL-005 | Permitted role success excerpt; denied role error naming `issue.claim.force_acquire`; zero-append proof for denial. |
| QA-09 | Verify renewal holder-only behavior, TTL preservation/override, and `activity_log` projection `(qa-lead)` | VAL-006 | Holder renewal outputs for no TTL and supplied TTL; event/activity excerpts; non-holder zero-append rejection. |
| QA-10 | Verify `fid issue unclaim` holder release and reject paths `(qa-lead)` | VAL-007 | CLI output for holder release; repeat/wrong-holder/missing-role/unknown-role errors; before/after journal and claim row proof. |
| QA-11 | Verify `fid claim list` required fields, `--by` filtering, and json/toon/markdown/yaml structural parity `(qa-lead)` | VAL-008 | Parsed output or structural comparison excerpts for all four formats with and without `--by`. |
| QA-12 | Verify `fid claim list --stale` membership and rowid-stable order before and after rebuild `(qa-lead)` | VAL-009 | CLI/DB fixture output and rebuild comparison showing live exclusion, stale inclusion, and identical order. |
| QA-13 | Verify `--by human` can hold/display a claim and blocks a non-holder agent transition `(qa-lead)` | VAL-010 | `issue get --section claims`, `claim list`, and transition rejection excerpts showing `human` holder. |
| QA-14 | Verify started-category entry claim gate for no claim, wrong holder, stale holder, and matching holder `(qa-lead)` | VAL-011 | CLI/domain test excerpts plus before/after journal/status proof for all four states. |
| QA-15 | Verify started-category exit auto-releases the caller claim in the same journal event/transaction `(qa-lead)` | VAL-012 | Event/effects excerpt showing status change plus release in one event, projection read-back showing no current claim. |
| QA-16 | Verify concurrent acquire atomicity `(qa-lead)` | VAL-013 | Integration test output showing exactly one successful acquire, one current row, and loser errors with holder/expiry. |
| QA-17 | Verify public docs and phase docs were updated or explicitly recorded as verified-no-change `(qa-lead)` | no VAL | Completion report excerpts and phase decisions/spec sync notes for `03_TRD_SCHEMA.md`, `05_DESIGN_PLAN.md`, `06_DESIGN_PLAN_CLI.md`, and task docs. |
| QA-18 | Verify no schema migration was added unless a recorded blocker/decision justifies it `(qa-lead)` | no VAL | `git diff --name-only main...HEAD -- migrations` output, plus DEC-026 reference or recorded decision. |
| QA-19 | Verify task completion reports have mechanically sourced commit hash and subject fields committed by their author `(qa-lead)` | no VAL | `git log --oneline -- project-docs/docs/phases/phase-11/tasks/task-*-done.md` and, for each copied report hash, `git rev-parse --verify HASH_FROM_COMPLETION_REPORT^{commit}` evidence. |
| QA-20 | Verify ESC-008 is closed by stale-projection claim write evidence `(qa-lead)` | VAL-015 | Review 5 excerpts showing acquire/renew/release stale-projection fixtures, absence of nested-lock errors, exact Task 5 gate results, and DEC-028 compliance. |

## 8. Mid-Phase Review Cadence

- Required cadence: per-task QA review, which is stricter than the default
  three-task cadence.
- Review 1: after Task 1 completes and before Task 2 starts.
- Review 2: after Task 2 completes and before Task 3 starts.
- Review 3: after Task 3 completes and before Task 4 starts.
- Review 4: after Task 4 completes and before final QA review.
- Review 5: after Task 5 completes and before final QA review.
- Review template:
  `project-docs/templates/mid-phase-review-template.md`.
- Each review verifies in-window `VAL-NNN` assertions, compares task completion
  command results against the task spec and this plan, and records any command
  skipped as BLOCKED with reason.
- A task may use focused red/green commands during development, but review
  acceptance requires the exact task close commands or an explicit blocker
  report.

## 9. Human Verification Scope

These items require human judgment and must be copied into
`project-docs/docs/phases/phase-11/qa/human-checklist.md` using
`project-docs/templates/human-checklist-template.md` at phase end.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-01a | In a real `.fid` workspace, identify from CLI output which role currently holds an issue claim `(human)` | VAL-014 | Human note naming command used and whether the holder was clear. |
| HUMAN-01b | In a real `.fid` workspace, identify from CLI output when the current issue claim expires `(human)` | VAL-014 | Human note stating whether expiry was understandable without code inspection. |
| HUMAN-01c | In a real `.fid` workspace, release an issue claim and confirm from CLI output/read-back that it is released `(human)` | VAL-014 | Human note with command/output summary and release comprehension verdict. |
| HUMAN-01d | In a real `.fid` workspace, try to advance an issue held by `human` as an agent role and judge whether the block reason is clear `(human)` | VAL-014 | Human note on whether the error explains the human-held claim and next action. |
| HUMAN-01e | Dogfood the claim lifecycle for normal work and judge whether claim/renew/unclaim/list messages are usable for agent coordination `(human)` | VAL-014 | Product-owner/dogfood note with accept/request-changes decision. |

## 10. Phase-End Decision Rule

The phase can proceed to final review only when:

- Required deterministic gates pass or are recorded BLOCKED with escalation.
- Reviews 1-5 exist and cover each completed task window.
- Every `VAL-NNN` has planned evidence and a final-review outcome slot.
- Every `(qa-lead)` item has reproducible evidence.
- Human checks are human-judgment-only and mapped to the human checklist.
- No unresolved blocking defects remain.
- Completion reports exist for every completed task and record exact command
  results.
- Required documentation updates or verified-no-change decisions are named.

Current QA-plan verdict: READY after Task 5 completes and Review 5 verifies
ESC-008 closure. No referenced required command is a placeholder.
