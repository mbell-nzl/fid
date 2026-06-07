# Phase 3 - Task 1: Generic Render Pipeline + In-Tree TOON Marshaller

## Objective

Introduce the generic render pipeline in `internal/format` (marshal the value to
JSON, decode to a generic `any` structure, render that structure) and implement
an in-tree TOON marshaller over it with no third-party dependency. Wire
`fid issue get --format toon` so the one shipped read command emits coherent TOON
that carries the same fields/values as its JSON output. This is the phase's
first vertical slice: the same code path that backs JSON now backs TOON, ready
for markdown (task 2) and YAML (task 3) to plug in.

## Requirements

- PRD refs: PRD-OUT-002 (TOON via `--format toon`), PRD-OUT-004 (`--format` on
  read commands)
- TRD refs: TRD-ARCH-012 (TOON marshaller in-tree, no third-party dependency),
  TRD §5 (TOON shape mirrors JSON; parity tests against JSON), TRD-ARCH-005
  (adapter boundary)
- Design refs: `05_DESIGN_PLAN.md` §5.1 (format strategy), `06_DESIGN_PLAN_CLI.md`
  §2 (`fid issue get` `--format`), §2 TOON list example (tabular shape)
- Decisions: DEC-013 (generic JSON→`any` render pipeline for JSON/TOON/YAML) —
  `../decisions.md`
- Quality refs: Quality Plan §4.3 (COV-002 `cmd/fid` ≥ 70%; `internal/format` is
  not a COV-001 package but needs deterministic evidence)

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-3/qa/validation-contract.md`.

- VAL-001: `fid issue get --format toon` is structurally at parity with JSON for
  empty-arrays and populated issues (nested arrays preserved).
- VAL-002: TOON marshaller is in-tree; no third-party TOON dependency.
- VAL-005: `--format` dispatch reaches the TOON renderer; an unknown format is
  still rejected (this task establishes the dispatch; markdown/yaml arms land in
  tasks 2-3).

## Scope

### Files To Create Or Modify

- `internal/format/format.go` (generic JSON→`any` structure helper; route TOON
  through the in-tree marshaller; keep JSON path intact; preserve rejection of
  genuinely-unknown formats)
- `internal/format/toon.go` (new: in-tree TOON marshaller over the generic
  structure — objects, scalars, and nested arrays)
- `internal/format/toon_test.go` (new: parity tests vs JSON + golden fixture +
  nested/empty-array cases)
- `internal/format/format_test.go` (new or extend: dispatch + unknown-format
  rejection)
- `cmd/fid/issue_test.go` (extend: `fid issue get --format toon` integration
  smoke, both states)
- `project-docs/docs/phases/phase-3/tasks/task-1-completion.md`

### Files To Read For Context

- `internal/format/format.go` (current JSON path + `UnsupportedFormatError`)
- `internal/issues/issue.go` (`Issue`/`AcceptanceCriterion` shapes the renderer
  must round-trip)
- `cmd/fid/issue.go` (how `fid issue get` calls `fidformat.Write`)
- `05_DESIGN_PLAN.md` §5.1, §5.2; `06_DESIGN_PLAN_CLI.md` §2

### Out Of Scope

- Markdown renderer (task 2) and YAML renderer (task 3).
- `fid issue list` or any new read command (v0.1b); TOON is exercised on the
  single shipped `get` command.
- Full TOON-spec coverage; implement only the minimal subset the issue shape
  needs (TRD §11 open question sanctions starting minimal).

## Acceptance Criteria

- [ ] `internal/format` exposes a generic render path: a value is marshalled to
  JSON and decoded to a generic `any` (map/slice/scalar) structure that the TOON
  renderer consumes. (VAL-001)
- [ ] `fid issue get --format toon` emits TOON whose decoded/structural content
  carries the same field names and scalar values as the JSON output, for both the
  canonical empty-arrays issue and a synthetic populated issue (labels + nested
  `acceptance_criteria`); nested arrays are preserved, not flattened. (VAL-001)
- [ ] No third-party TOON dependency is added; `go.mod`/module graph shows none
  and `go build` succeeds. (VAL-002)
- [ ] `--format toon` dispatches to the TOON renderer; an unknown `--format xml`
  is rejected with an actionable `error:` (no silent JSON fallback). (VAL-005)

## TDD Plan

- First failing test (`internal/format`): render the canonical issue value to
  TOON and assert field/value parity with its JSON encoding — red while TOON
  returns `UnsupportedFormatError`.
- Additional tests: populated issue with nested ACs + labels (parity + nesting
  preserved); empty-arrays issue; golden TOON fixture; dispatch table covering
  json/toon and unknown-format rejection; `cmd/fid` integration smoke.
- Regression risks: JSON output must stay byte-stable (existing JSON tests);
  generic decode must not reorder/lose null fields (05_DESIGN_PLAN §5.2 keeps
  nulls); keep `cmd/fid` a thin adapter (VAL-010); COV-002 `cmd/fid` ≥ 70%.

## Documentation Updates

- `project-docs/docs/phases/phase-3/tasks/task-1-completion.md` - completion
  report citing VAL-001/002/005 evidence and gate outcomes.
- If the TOON subset boundary needs documenting, add a short note to the
  completion report (and flag for spec sync at phase wrap) rather than editing
  specs mid-task.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: TOON is a compatibility surface (05_DESIGN_PLAN
  §8). It must mirror the JSON shape; do not introduce TOON-only field
  names/shapes.
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
  command (`fid issue get --format toon`); confirm the binary builds and smoke
  the TOON output. It also evidences VAL-002 (build with no TOON dependency).
- No `integration` build tag is in use in this repo; integration-style CLI tests
  run under `go test -race ./...`.

## Notes For Developer

- Keep JSON canonical and byte-stable; TOON is secondary (05_DESIGN_PLAN §5.1,
  §10.3 of `scaffold/initial-spec.md`). Render TOON from the same generic
  structure JSON produces so the shapes cannot drift.
- Implement TOON in-tree (TRD-ARCH-012). Do not add a TOON library.
- The render seam should generalise: future read commands pass any value and get
  JSON/TOON for free. Markdown (task 2) is the one renderer that needs a typed
  view, not this generic path.
- Keep render/parse logic in `internal/format`; `cmd/fid` only selects the format
  and writes the result (VAL-010).
