# Phase 16 - Task 6 Completion Report

## Summary of what was implemented

- Added `fid issue get <id> --section gates` as a read-only issue section over
  raw `gate_runs` history, ordered by attempt number and run time, with empty
  list behavior when no gates exist.
- Kept the new section distinct from `task_completions.gate_results` by reading
  append-only `gate_runs` rows directly and preserving the completion snapshot
  separately.
- Closed CF-016 for Phase 16 by satisfying VAL-012 across the domain reader,
  CLI section surface, all four output formats, and a built-CLI smoke run.

## Files created or modified

- `cmd/fid/issue_section_test.go`
- `internal/format/markdown.go`
- `internal/format/markdown_test.go`
- `internal/issues/detail.go`
- `internal/issues/detail_test.go`
- `internal/issues/get.go`
- `internal/issues/issue.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-16/tasks/task-6-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
  - Added `gates` to the `fid issue get --section` enum and documented the raw
    `gate_runs` row shape plus empty-list behavior.
- `project-docs/docs/phases/phase-16/tasks/task-6-done.md`
  - Recorded CF-016 as closed by this task with VAL-012 evidence.

## Validation Contract Coverage

- `VAL-012` PASS.
  - `internal/issues/detail_test.go` proves `gates` is a supported section,
    returns ordered raw `gate_runs` rows, preserves nullable `agent_run_id` and
    `output_excerpt`, and stays distinct from completion snapshots.
  - `cmd/fid/issue_section_test.go` proves `fid issue get --section gates`
    returns the raw history for populated issues, renders in `json`, `toon`,
    `yaml`, and `markdown`, and returns `[]` for an issue with no gates.
  - Built CLI smoke evidence:
    `/tmp/fid issue get <id> --section gates --format json` returned the
    recorded gate row after `dev start` -> `dev gate` -> `dev finish`, proving
    finished-task history still reads from `gate_runs`.
- Coverage percentages from the required coverage gate:
  - `internal/issues`: `81.5%`
  - `cmd/fid`: `77.5%`

## Tests run

- `go test ./internal/issues -run 'Test(GetReturnsCentralizedDetailShape|SectionValue)'`
  - Red: build failed because `Issue` had no hidden gates section yet.
  - Green: passed after wiring `gate_runs` into the issue detail assembly and
    section lookup.
- `go test ./cmd/fid -run TestIssueGetSectionsAndDueDateReadBack/sections`
  - Red: failed with `unsupported markdown value type []issues.IssueGate`.
  - Green: passed after adding markdown rendering for gate section rows.
- `go test ./internal/format -run TestWriteMarkdownRendersIssueGateLists`
  - PASS.
- Required task gate:
  - `gofmt -l cmd internal`
  - `golangci-lint run ./...`
  - `go test -race -coverprofile=coverage.out ./...`
  - `go tool cover -func=coverage.out`
  - `go build -o /tmp/fid ./cmd/fid`
- Real entrypoint exercise:
  - `/tmp/fid init`
  - `/tmp/fid issue create --kind task --title "Smoke gates" --by developer`
  - `/tmp/fid dev start <id> --model gpt-5 --by developer`
  - `/tmp/fid dev gate <id> --cmd "go test ./..." --exit-code 0 --output-excerpt ok --by developer`
  - `/tmp/fid dev finish <id> --completion-md report.md --by developer`
  - `/tmp/fid issue get <id> --section gates --format json`
  - Observed result: returned the raw gate row with `attempt_number`,
    `agent_run_id`, `command`, `exit_code`, `output_excerpt`, and `ran_at`
    after the task was finished, confirming the section remains distinct from
    the completion snapshot.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS
- `go test -tags integration -race ./...` - Not required
  - Reason: Task 6 added no integration-tagged tests.
- `go build -o /tmp/fid ./cmd/fid` - PASS

## Deviations from spec

- The task-scoped file list omitted the minimal hidden-field and formatter
  wiring needed to keep `--section gates` off the full `fid issue get` JSON
  surface while still supporting markdown output. This task therefore also
  touched `internal/issues/get.go`, `internal/issues/issue.go`,
  `internal/format/markdown.go`, and `internal/format/markdown_test.go`.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `31fb3dd`

## Commit subject

- `feat: add issue gates section`
