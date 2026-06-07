# Phase 3 - Task 2: Markdown Renderer For `fid issue get`

## Objective

Add the `--format markdown` renderer for `fid issue get`, producing the
human-readable layout in `06_DESIGN_PLAN_CLI.md` §21. Markdown cannot be produced
by the generic JSON→structure pipeline (task 1) because the layout is bespoke per
command, so this task introduces a typed renderer seam in `internal/format` (a
markdown renderer that takes the typed issue view) while JSON/TOON/YAML stay on
the generic path. Wire `fid issue get --format markdown`.

## Requirements

- PRD refs: PRD-OUT-003 (markdown via `--format markdown`), PRD-OUT-004
- TRD refs: TRD-ARCH-005 (adapter boundary; render logic in `internal/format`)
- Design refs: `06_DESIGN_PLAN_CLI.md` §21 (markdown example), `05_DESIGN_PLAN.md`
  §5.1 (format strategy), §8 (cross-surface parity matrix)
- Decisions: DEC-013 (typed markdown renderer seam) — `../decisions.md`
- Quality refs: Quality Plan §4.3 (COV-002 `cmd/fid` ≥ 70%; `internal/format`
  deterministic evidence)

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-3/qa/validation-contract.md`.

- VAL-003: `fid issue get --format markdown` renders per §21 for both a populated
  issue and an issue with no ACs/labels/links.
- VAL-005: `--format markdown` dispatches to the markdown renderer (extends the
  task-1 dispatch).

## Scope

### Files To Create Or Modify

- `internal/format/markdown.go` (new: typed markdown renderer for the issue view
  per §21)
- `internal/format/format.go` (route `--format markdown` to the typed renderer
  seam; keep generic path for json/toon/yaml)
- `internal/format/markdown_test.go` (new: golden markdown fixture + both states)
- `cmd/fid/issue.go` (pass the typed issue to the markdown path; keep adapter
  thin)
- `cmd/fid/issue_test.go` (extend: `fid issue get --format markdown` integration,
  both states)
- `project-docs/docs/phases/phase-3/tasks/task-2-completion.md`

### Files To Read For Context

- `06_DESIGN_PLAN_CLI.md` §21 (exact markdown layout)
- `05_DESIGN_PLAN.md` §8 (parity matrix: status, AC verification → checkbox,
  priority → P0–P5)
- `internal/format/format.go` and `internal/format/toon.go` (task-1 seam)
- `internal/issues/issue.go` (`Issue`/`AcceptanceCriterion` fields)

### Out Of Scope

- YAML (task 3) and the hardening tasks (4-5).
- TUI board (v0.3) and markdown *export* with frontmatter (v0.4) — this is CLI
  `--format markdown` only.
- `--section` on `fid issue get` (v0.1b).
- AC/label/link *write* commands; render whatever the projection currently holds
  (today ACs/labels are typically empty — the renderer must handle both).

## Acceptance Criteria

- [ ] A populated issue renders the §21 layout: `# <id>: <title>` heading, the
  kind/status/priority line, epic/owner/sequence/spec lines, an Acceptance
  Criteria section using `- [x]`/`- [ ]` mapped from latest verification state, a
  Labels section, a Links section, and the created/updated footer. (VAL-003)
- [ ] An issue with no ACs, labels, or links renders gracefully — empty or
  omitted sections, no panic, no stray/duplicated headers. (VAL-003)
- [ ] Priority maps to `P<N>` and status/status_category match the §8 parity
  matrix. (VAL-003)
- [ ] `fid issue get --format markdown` dispatches to the markdown renderer;
  json/toon still work; an unknown format is still rejected. (VAL-005)

## TDD Plan

- First failing test (`internal/format`): render the canonical issue to markdown
  and assert it matches the §21 golden — red while markdown returns
  `UnsupportedFormatError`.
- Additional tests: populated issue (ACs verified/unverified, labels, links);
  empty-state issue; priority/status mapping; `cmd/fid` integration both states.
- Regression risks: do not regress JSON/TOON dispatch from task 1; checkbox
  mapping must reflect latest verification per role, not raw presence; keep
  `cmd/fid` thin (VAL-010); COV-002 ≥ 70%.

## Documentation Updates

- `project-docs/docs/phases/phase-3/tasks/task-2-completion.md` - completion
  report citing VAL-003/005 evidence and gate outcomes.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: markdown is a parity surface (05_DESIGN_PLAN §8).
  Keep the markdown rendering of status/AC/priority consistent with the JSON
  fields and the §8 matrix.
- Audit or release-readiness impact: None.

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

- `go build -o /tmp/fid ./cmd/fid` is required: this task changes a user-facing
  command (`fid issue get --format markdown`); confirm the binary builds and
  smoke the markdown output.

## Notes For Developer

- Markdown is a typed renderer (the §21 layout is bespoke). Take the typed issue
  view; do not try to render markdown from the generic `any` structure.
- Architecture decision already taken by the Tech Lead and recorded as DEC-013
  (`../decisions.md`): the markdown renderer lives in `internal/format` as a typed
  renderer seam. You are not introducing an unrecorded architecture decision —
  implement that seam.
- Match `06_DESIGN_PLAN_CLI.md` §21 exactly for the populated case; the empty
  case has no example, so omit empty sections cleanly and assert no panic.
- Keep render logic in `internal/format`; `cmd/fid` only selects and writes.
