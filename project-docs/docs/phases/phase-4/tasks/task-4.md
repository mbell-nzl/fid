# Phase 4 - Task 4: Centralised Issue Detail Read Model + Sections

## Objective

Move the enriched issue detail assembly out of `cmd/fid` and into
`internal/issues`, then add `fid issue get --section` for the named structured
sections in the CLI spec. This closes CF-006 before a second read command reuses
the same issue shape.

## Requirements

- PRD refs: PRD-ISSUE-006, PRD-ISSUE-018, PRD-OUT-004
- TRD refs: TRD-ARCH-005, TRD-ARCH-008
- Design refs: `05_DESIGN_PLAN.md` section 6.1; `06_DESIGN_PLAN_CLI.md`
  `fid issue get`
- Quality refs: `08_QUALITY_PLAN.md` section 5.4, DEC-007 COV-001 carry-in
- Carry-forward: CF-006 in `project-docs/docs/current/decisions_active.md`
- Dependencies: Tasks 1-3

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-4/qa/validation-contract.md`.

- VAL-006: enriched issue detail read model is centralised in `internal/issues`.
- VAL-007: `fid issue get --section` returns named structured sections and
  rejects unknown sections.
- VAL-003 (read-back closeout, ESC-002): the centralised read model exposes
  `due_date` through `internal/issues.Get` and `fid issue get --format json`.
  Tasks 2-3 proved due-date persistence via direct SQLite query only; this task
  closes the public read-back clause of VAL-003. No schema change is required:
  `issues_enriched` already projects `i.*` (migration 0003) and `issues.due_date`
  exists (migration 0001).

## Scope

### Files To Create Or Modify

- `internal/issues/issue.go` or new focused read-model files under
  `internal/issues/`
- `internal/issues/issue_test.go` or new package-level tests
- `internal/format/markdown.go` only if the typed markdown view needs to move
  with the centralised read model
- `cmd/fid/issue.go`
- `cmd/fid/issue_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` if section output shape needs an
  example or clarification
- `project-docs/docs/phases/phase-4/tasks/task-4-completion.md`

### Files To Read For Context

- `cmd/fid/issue.go` (current markdown SQL assembly to remove)
- `internal/issues/issue.go`
- `internal/format/markdown.go`
- `project-docs/specs/05_DESIGN_PLAN.md` section 6.1
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` `fid issue get`

### Out Of Scope

- Full markdown-body heading extraction from `spec_md_path`; Phase 4 implements
  named structured sections only.
- `--context-pack` (v0.2).
- Claim/activity timeline enrichment beyond rows that already exist in the
  projection.

## Acceptance Criteria

- [ ] `internal/issues` exposes a single enriched detail read shape used by both
  JSON and markdown issue get paths. (VAL-006)
- [ ] Populated and empty issue detail tests cover labels, acceptance criteria,
  latest verifications per role, open concern count, latest triage status, and
  links. (VAL-006)
- [ ] `cmd/fid/issue.go` no longer performs direct SQL assembly for parent
  title, labels, acceptance criteria/latest verifications, or links. (VAL-006)
- [ ] `fid issue get <id> --section acceptance|concerns|completions|links|claims`
  returns only the requested structured section and handles empty states.
  (VAL-007)
- [ ] `fid issue get <id> --section unknown` exits non-zero and names supported
  sections. (VAL-007)
- [ ] The centralised read model carries `due_date`, and a populated due date set
  via `fid issue create --due-date` reads back through both `internal/issues.Get`
  and `fid issue get --format json`; an omitted due date reads back as null.
  (VAL-003, closes ESC-002)

## TDD Plan

- First failing test (`internal/issues`): populated detail read returns labels,
  AC latest verification by role, and links without helper SQL in `cmd/fid`.
- Additional tests: empty detail state; section outputs for all supported names;
  unknown section rejection; markdown renderer still receives the data it needs.
- Due-date read-back (VAL-003 / ESC-002): a failing `internal/issues` test
  asserting populated `due_date` returns through `Get`, plus a `cmd/fid` test
  asserting `fid issue get --format json` exposes `due_date` after
  `fid issue create --due-date`; omitted due date reads back as null. These
  replace the direct-SQLite due-date checks from Tasks 2-3 as the VAL-003
  read-back evidence.
- Regression risks: latest-verification tie-breaking, markdown renderer drift,
  `cmd/fid` retaining SQL helpers, ambiguous `--section` semantics, and the new
  `due_date` field drifting between JSON and markdown read shapes.

## Documentation Updates

- `project-docs/docs/phases/phase-4/tasks/task-4-completion.md` - record
  VAL-006/007 evidence and exact command outcomes.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - update only if needed to document
  the exact structured section output shape implemented. If the PRD/CLI section
  semantics appear contradictory, stop and escalate before implementation.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None; do not read held-out scenarios.
- Cross-surface parity impact: JSON and markdown issue get must consume the same
  enriched detail source.
- Audit or release-readiness impact: TRD-ARCH-005 adapter boundary is explicitly
  rechecked in this task.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Every required command
must be recorded in the completion report as `PASS`, or as `BLOCKED` with reason
in a blocker report. Missing, silently skipped, or narrower substituted commands
are failed gates.

Conditional commands:

```text
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go build -o /tmp/fid ./cmd/fid` is required because this task changes the
  public `fid issue get --section` surface.

## Output Requirements

- Write `project-docs/docs/phases/phase-4/tasks/task-4-completion.md`.
- Completion report must cite VAL-006, VAL-007, and VAL-003 (due-date read-back,
  closing ESC-002).

## Notes For Developer

- Treat the validation contract's Phase 4 planning assumptions as binding unless
  the human revises them before execution.
- Keep the detail shape reusable by `fid issue list` where appropriate, but do
  not force list rows to carry the full detail payload.
- Do not implement `--context-pack`.
