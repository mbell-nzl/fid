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

## TDD Plan

- First failing test (`internal/issues`): populated detail read returns labels,
  AC latest verification by role, and links without helper SQL in `cmd/fid`.
- Additional tests: empty detail state; section outputs for all supported names;
  unknown section rejection; markdown renderer still receives the data it needs.
- Regression risks: latest-verification tie-breaking, markdown renderer drift,
  `cmd/fid` retaining SQL helpers, and ambiguous `--section` semantics.

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
- Completion report must cite VAL-006 and VAL-007.

## Notes For Developer

- Treat the validation contract's Phase 4 planning assumptions as binding unless
  the human revises them before execution.
- Keep the detail shape reusable by `fid issue list` where appropriate, but do
  not force list rows to carry the full detail payload.
- Do not implement `--context-pack`.
