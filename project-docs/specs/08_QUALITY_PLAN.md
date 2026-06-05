# QUALITY PLAN — fid

> This document defines the testing strategy, quality gates, code standards,
> independent verification, and release-readiness practices for fid. Designed
> for execution by AI coding agents under human oversight, with deterministic
> checks as the primary quality enforcement mechanism.

---

## 1. Quality Philosophy

1. **Trust behavior, verify structure.** Agents can make their own tests pass;
   quality requires deterministic execution evidence, static analysis, and
   periodic review from outside the implementation session.
2. **Tests are executable specifications.** Behavior changes start with a
   failing test unless the task is explicitly documentation-only or exploratory.
3. **Deterministic checks first.** Formatting, linting, tests, coverage, build,
   contract, security, dependency, license, and secret checks run before
   subjective review when they apply.
4. **Agents cannot weaken tests.** Existing tests may be refactored only when
   the task requires it and behavior-preserving intent is clear.
5. **Human judgment is reserved for intent.** Agents execute deterministic
   checks; humans validate product fit, perception, and dogfooding.
6. **Specs are quality infrastructure.** Ambiguous requirements are defects.
   Fix the spec or escalate before implementing around the gap.
7. **Validation contracts make proof explicit.** Phase-level `VAL-NNN`
   assertions connect requirements to task evidence, QA review, and final
   sign-off.
8. **Independent checks compound.** Audit prompts, held-out scenarios, and CI
   gates are reusable controls that should improve across projects.

---

## 2. Quality Layer Framework

| Layer | What It Catches | Default Owner | When Required |
|---|---|---|---|
| 1. Spec precision | Ambiguity, missing acceptance criteria, scope drift | Human + Tech Lead | Before implementation |
| 2. Validation contract | Unproven requirements, vague acceptance, missing evidence ownership | Tech Lead + QA Lead | Every planned phase |
| 3. Red/green TDD | Functional regressions and edge cases | Developer | Every behavior task |
| 4. Agent QA review | Spec compliance, test quality, architecture smells | QA Lead / QA Reviewer | Every phase |
| 5. Automated static gates | Mechanical defects, vulnerable patterns, dependency and license risk, leaked secrets | CI / QA Lead | Before external use; earlier when cheap |
| 6. Held-out scenario validation | Spec gaming and end-to-end product misses | Human / evaluator outside implementation context | Before release readiness |
| 7. Independent audits | Structural drift, security gaps, complexity, stale specs | Fresh agent/model + human triage | Periodic or pre-release |
| 8. Manual behavioral validation | Perception, product judgment, domain fit | Human | Every phase sign-off |

---

## 3. Test Strategy

### 3.1 Test Types

| Type | Scope | Location | Runner | When Run |
|---|---|---|---|---|
| Unit | Single function or method; mocked dependencies | `*_test.go` alongside source | `go test ./...` | Every task |
| Integration | Multiple components with real SQLite (temp DB per test) | `*_test.go` with `//go:build integration` tag | `go test -tags integration ./...` | Phase end and targeted tasks |
| CLI | End-to-end command execution; verifies output, exit codes, DB state | `cmd/fid/*_test.go` or `internal/*/cmd_*_test.go` | `go test ./...` | Every CLI task |
| Journal | Replay correctness; golden-file comparison per schema version boundary | `internal/journal/*_test.go` | `go test ./internal/journal/...` | Every journal/schema task |
| Contract | Public interfaces, JSON output shapes, event schema shapes | `*_test.go` with snapshot/golden files | `go test ./...` | Any interface task |
| Benchmark | Performance targets (read latency, rebuild speed, write latency) | `*_test.go` with `Benchmark*` functions | `go test -bench . -benchmem ./...` | Targeted tasks; phase end |

### 3.2 Test File Convention

Every source file with logic shall have a corresponding test file:

```text
internal/journal/writer.go        → internal/journal/writer_test.go
internal/issues/service.go        → internal/issues/service_test.go
internal/claims/service.go        → internal/claims/service_test.go
internal/workflow/engine.go       → internal/workflow/engine_test.go
internal/policy/matrix.go         → internal/policy/matrix_test.go
internal/format/json.go           → internal/format/json_test.go
internal/format/toon.go           → internal/format/toon_test.go
internal/db/schema.go             → internal/db/schema_test.go
internal/stats/aggregator.go      → internal/stats/aggregator_test.go
cmd/fid/issue.go                  → cmd/fid/issue_test.go
```

### 3.3 Test Patterns (Go-Specific)

**TDD-001** [Ubiquitous]
All implementation code shall be written using red/green TDD: write a failing
test first, then write the minimum code to make it pass, then refactor.

**TDD-002** [Ubiquitous]
Unit tests shall use table-driven test patterns where a function is tested
with multiple inputs and expected outputs.

```go
func TestWorkflow_Transition(t *testing.T) {
    tests := []struct {
        name      string
        from      string
        to        string
        wantErr   bool
        errContains string
    }{
        {name: "valid: todo to in_progress", from: "todo", to: "in_progress", wantErr: false},
        {name: "invalid: todo to dev_complete", from: "todo", to: "dev_complete", wantErr: true, errContains: "not allowed"},
        {name: "cancel from any state", from: "in_progress", to: "cancelled", wantErr: false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ... test body
        })
    }
}
```

**TDD-003** [Ubiquitous]
Tests shall use the testify library (`assert` and `require` packages) for
assertions. `require` for preconditions that must succeed for the test to be
meaningful. `assert` for assertions that should be checked even if earlier
assertions fail.

**TDD-004** [Ubiquitous]
Database tests shall use a temporary SQLite database with migrations applied
in a test helper function. Each test shall start with a clean database state.
Use `t.TempDir()` for the `.fid/` directory to ensure isolation.

**TDD-005** [Ubiquitous]
Tests that interact with the filesystem (journal files, spec markdown,
export) shall use `t.TempDir()` for isolated temporary directories.

**TDD-006** [Ubiquitous]
Tests shall not depend on external services, network access, or system clock.
Time-dependent logic shall accept a clock interface for test injection.

**TDD-007** [Ubiquitous]
Journal replay tests shall use golden-file fixtures: capture a real event at
schema version N, run it through the upgrader chain, and assert the output
matches the expected current-version shape. Golden files live in `testdata/`.

**TDD-008** [Ubiquitous]
Trigger tests shall attempt invalid inserts via raw SQL and assert that the
trigger fires with the expected RAISE message.

### 3.4 Test Data

**TD-001** [Ubiquitous]
The `testdata/` directory shall contain fixture files for journal replay,
golden-file comparisons, and sample markdown specs. Test fixtures shall be
committed to the repository.

**TD-002** [Ubiquitous]
Test helpers for creating common entities (issues, ACs, completions, claims,
roles) shall be defined in a `internal/testutil/` package to avoid
duplication across test files.

**TD-003** [Ubiquitous]
Benchmark tests shall use a seeded database with a configurable entity count
(default: 10,000 issues, 100,000 journal events) to verify performance
targets.

---

## 4. Code Quality Standards

### 4.1 Formatting

**FMT-001** [Ubiquitous]
All Go code shall be formatted with `gofmt`. No exceptions.

**FMT-002** [Ubiquitous]
All Go imports shall be organised by `goimports` into three groups: standard
library, external packages, internal packages.

### 4.2 Linting

**LINT-001** [Ubiquitous]
All Go code shall pass `golangci-lint run` with zero errors before any commit.

**LINT-002** [Ubiquitous]
The golangci-lint configuration (`.golangci.yml`) shall enable at minimum:
`errcheck`, `govet`, `staticcheck`, `unused`, `gosimple`, `ineffassign`,
`misspell`, `gofmt`, `goimports`.

### 4.3 Coverage

**COV-001** [Ubiquitous]
The project shall maintain a minimum test coverage of 80% across the domain
packages (`internal/issues/`, `internal/claims/`, `internal/workflow/`,
`internal/policy/`, `internal/journal/`).

**COV-002** [Ubiquitous]
The project shall maintain a minimum test coverage of 70% across the
CLI layer (`cmd/fid/`).

**COV-003** [Ubiquitous]
Coverage shall be measured and reported on every test run using
`go test -coverprofile=coverage.out ./...` and visualised with
`go tool cover -func=coverage.out`.

**COV-004** [Ubiquitous]
Coverage gates shall be enforced in CI. A change that reduces overall
coverage below the thresholds shall not be merged.

### 4.4 Code Structure

**STRUCT-001** [Ubiquitous]
No function shall exceed 60 lines (excluding tests). Functions longer than
60 lines shall be refactored into smaller, named functions.

**STRUCT-002** [Ubiquitous]
No file shall exceed 400 lines (excluding tests). Files longer than 400 lines
shall be split into logically cohesive files.

**STRUCT-003** [Ubiquitous]
All exported functions and types shall have doc comments following Go
conventions.

**STRUCT-004** [Ubiquitous]
Error handling shall follow Go conventions: errors are returned, not panicked.
Errors shall be wrapped with context using `fmt.Errorf("operation: %w", err)`.

### 4.5 Documentation

**DOC-001** [Ubiquitous]
Every task that changes user-visible behaviour, CLI commands, configuration,
or repository conventions shall update the impacted documentation in the same
change set.

**DOC-002** [Ubiquitous]
Documentation updates shall be reviewed as part of the Definition of Done.
Outdated examples, commands, and file references shall be treated as defects.

---

## 5. Quality Gates

### 5.1 Task-Level Gate (Run by Developer Agent)

Before completing any task, the developer agent shall run the exact commands
listed in the task spec. If the task spec does not define commands, run:

```bash
gofmt -l .
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

All required commands must pass, or be recorded as `BLOCKED` with reason in a
blocker report. If any fail, the developer agent must fix the issue and rerun
the exact failed command before committing. Missing, silently skipped, or
narrower substituted commands are failed gates. Coverage thresholds must not
decrease.

### 5.2 Phase-Level Gate (Run by QA Lead Agent)

At the end of each implementation phase, the QA Lead shall:

1. Run all unit tests: `go test -race ./...`
2. Run integration tests: `go test -tags integration -race ./...`
3. Run linting: `golangci-lint run ./...`
4. Check coverage: verify thresholds are met
5. Run benchmarks against performance targets: `go test -bench . -benchmem ./...`
6. PRD compliance: verify each requirement in the phase plan has been
   implemented and tested
7. Security review: check for SQL injection risks, command injection,
   unsafe input handling
8. Performance spot-check: verify CLI commands complete within TRD-PERF
   thresholds
9. Documentation review: verify impacted docs, examples, and usage guidance
   are current
10. **Execute the deterministic verification checklist on a seeded fixture.**
    Run every `(qa-lead)` item from the QA plan against a `fid` instance
    initialised in a temp directory. Paste literal CLI output excerpts into
    `final-review.md` under each item. **PASS verdicts without an output
    excerpt are not allowed** — if the evidence cannot be shown, the verdict
    is FAIL or BLOCKED.
11. Generate the **human verification checklist** (the `(human)` subset —
    perceptual, dogfooding, and product-owner judgment items).

### 5.3 Human Verification Gate

At the end of each phase, the human shall:

1. Review the QA Lead's assessment and human verification checklist
2. Run the tool manually for real project work (dogfooding)
3. Verify error messages are clear and actionable
4. Verify the ready queue returns sensible results
5. Approve or request changes before the phase is marked complete

### 5.4 Validation Contract Checks

Each phase uses
`project-docs/docs/phases/phase-N/qa/validation-contract.md`, created from
`project-docs/templates/validation-contract-template.md`.

Required checks:

- Tech Lead drafts `VAL-NNN` assertions before or alongside task specs.
- Each assertion states one observable behavior, invariant, or release-readiness
  claim. Multi-observable human assertions must be split or explicitly marked
  for one-row-per-observable human checklist fan-out.
- Conditional or binary assertions must define evidence for both the qualifying
  and non-qualifying states; positive-only evidence is insufficient.
- Parity or "matches reference" assertions need a structural conformance or
  parity proof when the reference can be checked deterministically.
- QA Lead refines evidence methods, validator roles, and task windows before
  implementation starts, and bounces under-specified contracts to Tech Lead
  before finalizing the QA plan.
- Developer task specs and completion reports cite covered `VAL-NNN` IDs or
  explicitly state why none apply.
- QA Reviewer verifies relevant `VAL-NNN` IDs for each reviewed task window.
- QA Lead records final `PASS`, `FAIL`, `BLOCKED`, or `CARRY-FORWARD` outcome
  for every `VAL-NNN`.

If a required validation assertion cannot be satisfied, the responsible agent
stops and escalates. Silent spec adaptation or evidence-free `PASS` is not
allowed.

### 5.5 Deterministic And Human Checklists

QA plans must split verification by type of judgment:

- Deterministic items end with `(qa-lead)`. QA Lead executes these and records
  literal command output, file-state evidence, screenshots, logs, or other
  reproducible proof in `final-review.md`.
- Human items end with `(human)`. The human executes these when judgment,
  perception, product intent, or dogfooding is required. Human items that
  verify a validation assertion cite the relevant `VAL-NNN`.

PASS verdicts on deterministic items require evidence. If evidence cannot be
shown, mark the item FAIL or BLOCKED. QA Lead must not PASS a perceptual or
product-owner judgment by reading code.

### 5.6 Command And Feasibility Validation

Before a QA plan is finalized, QA Lead validates every referenced command, test
scope, fixture setup, artifact assumption, and human-checklist feasibility using
help output, dry-run mode, documentation, existing file-state checks, or
equivalent proof. Invalid or unvalidated steps are QA-plan defects.

### 5.7 Required QA Artifacts

Each phase stores QA artifacts under `project-docs/docs/phases/phase-N/qa/`.
Use the templates in `project-docs/templates/`:

- `validation-contract-template.md` → `validation-contract.md`
- `qa-plan-template.md` → `qa-plan.md`
- `mid-phase-review-template.md` → `review-M.md`
- `final-review-template.md` → `final-review.md`
- `human-checklist-template.md` → `human-checklist.md`

QA artifacts must cite evidence, not prose-only confidence. Missing required
artifacts are hard stops for the phase runner or phase wrap-up.

### 5.8 Release-Readiness Gate

Before a public release, define and run:

```bash
govulncheck ./...                    # vulnerability audit
go-licenses check ./...              # license compatibility
gitleaks detect                      # secret scan
go test -tags integration -race ./...  # full test suite
go build -o /dev/null ./cmd/fid      # build verification
```

If any command is intentionally deferred, record why and what compensating
control applies.

### 5.9 Workflow-Lint And Skill Drift

When repo skills are present, workflow-lint should verify:

- Canonical skills exist under `project-docs/skills/`.
- Generated Codex copies exist under `.agents/skills/`.
- Generated Claude copies exist under `.claude/skills/`.
- Generated copies match canonical skill sources, except for deliberate
  provider-specific metadata.
- Every skill has valid Agent Skills frontmatter and a `Provenance:` line in
  the Markdown body.
- Side-effectful skills disable implicit invocation where the provider supports
  that control.

---

## 6. CI / Automation

### 6.1 Pipeline Stages

```text
push/PR → Format Check → Lint → Unit Tests → Integration Tests → Coverage Gate → Build
```

**CI-001** [Event-driven]
When a commit is pushed to any branch, the CI pipeline shall run: `gofmt`
check, `golangci-lint`, unit tests with race detection, and coverage
reporting.

**CI-002** [Event-driven]
When a pull request is opened or updated against `main`, the CI pipeline
shall additionally run integration tests and enforce coverage thresholds.

**CI-003** [Event-driven]
When a Git tag matching `v*` is pushed, the CI pipeline shall run the full
test suite, build cross-platform binaries, and create a GitHub release.

### 6.2 CI Configuration

The project shall use GitHub Actions:

```yaml
# Simplified structure — full workflow defined during v0.1
jobs:
  test:
    - Install Go (latest stable)
    - Run golangci-lint
    - Run go test -race -coverprofile=coverage.out ./...
    - Run go test -tags integration -race ./...
    - Check coverage thresholds
    - Upload coverage report
  build:
    - Build binaries for linux/amd64, darwin/amd64, darwin/arm64, windows/amd64
    - Verify binary size < 30 MB (TRD-PERF-007)
    - Run smoke test: fid --help exits 0
```

---

## 7. Held-Out Scenario Validation

Held-out scenarios are user-journey checks that implementation agents do not
see. They are designed to catch shared blind spots between generated code and
generated tests.

Rules:

- Store scenarios outside normal implementation context, ideally in a private
  repository or an ignored/excluded path that agents do not load.
- Treat a scenario as contaminated if an implementation agent reads it; retire
  or rewrite it.
- Scenarios describe persona, goal, starting state, trajectory questions,
  satisfaction criteria, and forbidden behavior.
- Mode A: an evaluator agent drives the built product and grades the trace.
- Mode B: the human plays the persona and logs friction. Run at least one Mode B
  scenario before a release-readiness sign-off.
- Scenario failures become normal roadmap/task work; do not patch around them
  inside the evaluator run.

Use `project-docs/templates/held-out-scenario-template.yaml` as the starting
point.

---

## 8. Independent Audits

Use a fresh agent session, and preferably a different model family from the
implementation loop, for high-stakes audits. Save outputs under
`project-docs/docs/reviews/`.

| Audit | Cadence | Template |
|---|---|---|
| Architecture | Every major phase block or release train | `templates/audits/audit-prompts.md` |
| Security | Before external release; periodically for sensitive data | `templates/audits/audit-prompts.md` |
| Maintainability | When change velocity slows or before handoff | `templates/audits/audit-prompts.md` |
| Performance | Before scale-sensitive release or after hot-path changes | `templates/audits/audit-prompts.md` |
| Complexity | When files/modules grow or refactors accumulate | `templates/audits/audit-prompts.md` |
| Spec drift | End of major phase blocks | `templates/audits/audit-prompts.md` |
| Dependency / supply chain | Quarterly or before release | `templates/audits/audit-prompts.md` |

---

## 9. Release Readiness

Before external use, the release candidate should have:

- [ ] All deterministic gates green on the release branch.
- [ ] Threat model written or refreshed.
- [ ] Pre-release checklist completed.
- [ ] Changelog or release notes updated.
- [ ] Migration path documented if data model or public contracts changed.
- [ ] No new dependency without recorded budget rationale and license check.
- [ ] At least one independent audit current for the risk profile.
- [ ] Held-out scenario run completed; at least one Mode B run since the last
  release-readiness sign-off.
- [ ] Rollback, backup, or recovery path verified where applicable.
- [ ] Binary size verified < 30 MB.
- [ ] Performance benchmarks pass TRD-PERF targets.

Use `project-docs/templates/threat-model-template.md` and
`project-docs/templates/pre-release-checklist.md`.

---

## 10. Definition Of Done

A task is done when:

- [ ] Failing tests were written first or an explicit exception was recorded.
- [ ] Implementation satisfies task acceptance criteria.
- [ ] Covered `VAL-NNN` assertions are cited with evidence, or the task records
      why no validation assertion applies.
- [ ] Impacted documentation was updated.
- [ ] Task-level quality gate passed (`gofmt`, `golangci-lint`, `go test` with
      coverage).
- [ ] Dependency-budget note is present if new dependencies were added.
- [ ] Completion report was written.
- [ ] Conventional Commit created with phase/task reference.

A phase is done when:

- [ ] All planned tasks are complete or explicitly deferred/cancelled.
- [ ] Phase-level quality gate passed.
- [ ] Every `VAL-NNN` has a final evidenced outcome or recorded carry-forward.
- [ ] QA Lead final review has no unresolved blocking findings.
- [ ] Human checklist was completed and accepted or follow-up work was planned.
- [ ] Specs were synchronized where needed.
- [ ] Phase decisions and progress were archived.
- [ ] Current docs were compressed.

*Last updated: 2026-05-27*
