# Phase 3 - Task 3: YAML Output For `fid issue get`

## Objective

Add the `--format yaml` renderer for `fid issue get` over the generic
JSON→structure pipeline from task 1, so the issue serialises to YAML whose
decoded structure equals its JSON output. This completes the four-format
`--format` surface.

> **YAML approach is decided (DEC-012, human-approved 2026-06-07): adopt
> `gopkg.in/yaml.v3`** (MIT) for the YAML renderer (emit) and YAML parity tests
> (decode). No in-tree emitter; no further escalation needed. Record the
> dependency in `go.mod`/`go.sum` and in the completion report per the
> AGENTS.md dependency budget (the need/alternatives/license/reuse rationale is
> already captured in DEC-012 — restate it briefly in the report).

## Requirements

- PRD refs: PRD-OUT-003a (YAML via `--format yaml`, Could/v0.1), PRD-OUT-004
- TRD refs: TRD-ARCH-005 (adapter boundary); TRD §9 / `AGENTS.md` dependency
  budget (DEC-012 adds `gopkg.in/yaml.v3`)
- Design refs: `05_DESIGN_PLAN.md` §5.1 (format strategy; YAML "secondary"),
  `scaffold/initial-spec.md` §10.1
- Decisions: DEC-012 (adopt `gopkg.in/yaml.v3`), DEC-013 (generic render
  pipeline) — `../decisions.md`
- Quality refs: Quality Plan §4.3 (COV-002 `cmd/fid` ≥ 70%; `internal/format`
  deterministic evidence)

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-3/qa/validation-contract.md`.

- VAL-004: `fid issue get --format yaml` decodes to a structure equal to the JSON
  output, for both empty-arrays and populated issues.
- VAL-005: `--format yaml` dispatches to the YAML renderer (completes the
  json/toon/markdown/yaml dispatch set).

## Scope

### Files To Create Or Modify

- `internal/format/yaml.go` (new: YAML renderer over the generic structure via
  `gopkg.in/yaml.v3`)
- `internal/format/format.go` (route `--format yaml` through the generic path)
- `internal/format/yaml_test.go` (new: YAML decode-and-deep-equal-JSON parity,
  both states)
- `cmd/fid/issue_test.go` (extend: `fid issue get --format yaml` integration,
  both states)
- `go.mod` / `go.sum` (add `gopkg.in/yaml.v3`, DEC-012)
- `project-docs/docs/phases/phase-3/tasks/task-3-completion.md`

### Files To Read For Context

- `internal/format/format.go`, `internal/format/toon.go` (task-1 generic seam)
- `05_DESIGN_PLAN.md` §5.1, §5.2
- `internal/issues/issue.go` (`Issue` shape)

### Out Of Scope

- Markdown (task 2) and the hardening tasks (4-5).
- v0.4 markdown export frontmatter (`adrg/frontmatter`) — that is a separate
  later surface; do not pull it in here.
- Any YAML *input*/parsing for config; this is output rendering only.

## Acceptance Criteria

- [ ] `gopkg.in/yaml.v3` is added as a direct dependency (DEC-012) and recorded
  with need/alternatives/license in the completion report; `AGENTS.md` dependency
  budget is honoured. (VAL-004)
- [ ] `fid issue get --format yaml` output decodes (via `gopkg.in/yaml.v3`) to a
  structure equal to the JSON output for both the canonical empty-arrays issue
  and a synthetic populated issue (labels + nested ACs). (VAL-004)
- [ ] `--format yaml` dispatches to the YAML renderer; json/toon/markdown still
  work; an unknown format is still rejected. (VAL-005)

## TDD Plan

- First failing test (`internal/format`): render the canonical issue to YAML,
  decode it, and assert deep equality with the JSON-decoded value — red while
  YAML returns `UnsupportedFormatError`.
- Additional tests: populated + empty states; nested ACs/labels preserved;
  `cmd/fid` integration both states.
- Regression risks: do not regress task-1/2 dispatch; YAML must keep null fields
  (05_DESIGN_PLAN §5.2); if using an in-tree emitter, cover string escaping and
  nested sequences; COV-002 ≥ 70%.

## Documentation Updates

- `project-docs/docs/phases/phase-3/tasks/task-3-completion.md` - completion
  report citing VAL-004/005 evidence and the `gopkg.in/yaml.v3` dependency-budget
  rationale (need/alternatives/license/reuse, per DEC-012).

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: YAML mirrors the JSON shape (05_DESIGN_PLAN §5.1).
  Do not introduce YAML-only field names.
- Audit or release-readiness impact: `gopkg.in/yaml.v3` is the first direct
  dependency beyond migrate/cobra/sqlite; record it for the dependency audit /
  license check (`govulncheck`, `go-licenses`).

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
  command (`fid issue get --format yaml`); confirm the binary builds and smoke
  the YAML output.

## Notes For Developer

- YAML is decided (DEC-012): use `gopkg.in/yaml.v3`. Do not write an in-tree
  emitter.
- Render YAML from the generic JSON→structure pipeline (DEC-013) so YAML cannot
  drift from JSON; do not marshal the `Issue` struct directly with yaml tags.
- The parity test decodes the YAML output with `gopkg.in/yaml.v3` and deep-equals
  it to the JSON-decoded value — the same decoder library both emits and verifies.
- Keep render logic in `internal/format`; `cmd/fid` only selects and writes
  (VAL-010).
