# Phase 17 - Task 8: Prove v0.1 Read-Format Conformance

## Objective

Add one integration-tagged conformance matrix over every implemented v0.1 read
command. Prove json/toon/yaml structural parity and markdown identity coverage,
closing the cross-surface output exit criterion without changing healthy
commands opportunistically.

## Validation Contract Coverage

- VAL-017: every v0.1 read surface supports all four formats with deterministic
  data/identity conformance.

## Binding Constraints

- Exact command matrix: `issue get`, `issue list`, `issue ac list`, `ready`,
  `concern list`, `claim list`, `link list`, `label list`, `policy show`,
  `stats tokens`, `stats cost`, `activity list`, `sql`, `journal inspect`, and
  `maintenance check`.
- Seed one real temp workspace through shipped CLI writes. Do not directly seed
  normal domain state except where a command already has a documented test
  fixture helper.
- Compare JSON/TOON/YAML structurally after decode where a parser exists;
  markdown must render successfully and contain the same primary IDs/labels.
- Command-specific shape remains canonical. Do not force unlike results into a
  generic list or weaken existing markdown presentation.
- If conformance exposes a broad product/output ambiguity, stop and escalate;
  fix only a narrow documented omission inside this task.
- Do not modify `internal/sqlquery`. If SQL conformance requires that package to
  change, stop so Tech Lead can schedule CF-018 as a separate cleanup first.

## Scope

### Files To Create Or Modify

- `cmd/fid/e2e_v01_format_matrix_test.go` - integration command/format matrix
  and structural conformance helpers.
- `cmd/fid/*` - only a narrowly identified missing `--format` dispatch.
- `internal/format/*` - only a narrowly identified missing renderer/decoder
  needed for a documented v0.1 read shape; preserve task-1 budgets.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - reconcile only if the canonical
  matrix is missing a shipped format declaration.

### Files To Read For Context

- All command files named by the matrix under `cmd/fid/`
- `internal/format/format.go`, `toon.go`, `yaml.go`, and markdown files
- `cmd/fid/e2e_v01b_test.go`
- `cmd/fid/e2e_v01e_activity_test.go`
- `project-docs/specs/01_PRD.md` PRD-OUT-001..004

### Out Of Scope

- Write-command confirmation output.
- Version-gated governance commands.
- A universal public output schema or new format.
- New product behavior unrelated to format conformance.

## Acceptance Criteria

- [ ] Every exact matrix row accepts json/toon/markdown/yaml. (VAL-017)
- [ ] JSON, TOON, and YAML outputs preserve the same required fields, row count,
  identities, and empty/non-empty classification. (VAL-017)
- [ ] Markdown succeeds for every row and contains the same primary identities
  without unsupported-type errors. (VAL-017)
- [ ] The matrix is deterministic and hermetic: temp workspace, no network,
  stable row order, and no wall-clock ordering assumptions. (VAL-017)

## First Failing Test

`TestV01ReadCommandsSupportAllFormatsWithStructuralParity` begins as a complete
table and fails on the first unsupported or divergent command/format pair.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - only for a discovered declaration
  gap; otherwise record verified-no-change in the completion report.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

Required-command rules live in `AGENTS.md`. Every command is reported `PASS`,
or `BLOCKED` with reason; the full gate, integration suite, and build are
required before task close.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Depends on task-5; task-7 is not part of the read matrix.
- Prefer shared test-only decode/assert helpers. Do not add production parsing
  solely to compare test output.
- Completion report: `project-docs/docs/phases/phase-17/tasks/task-8-done.md`.
