# Phase 20 - Task 3d: Evidence-token control and gate-record correction

## Objective

Make the fabricated-token failure mode mechanically detectable instead of
review-dependent, and correct the Phase 20 completion reports that recorded a
full-suite `go test` gate with an impossible token (DEC-131).

## Validation Contract Coverage

- No `VAL-NNN` applies. This is V1 delivery-process tooling and evidence-record
  correction, not product behavior. `AGENTS.md` (Commands) is the governing
  rule and phase-19 DEC-102 is the governing precedent.

## Binding Constraints

- `project-docs/tools/report-lint.sh` is workflow tooling, not product code. It
  stays POSIX `sh`, dependency-free, and non-interactive, and its exit code
  stays the number of failed checks.
- The new check is **narrow and mechanical**. It must not attempt to judge
  whether a command truly ran; it detects only a token that cannot have come
  from the recorded command.
- Extend the existing `--self-test` with cases that fail before the change and
  pass after, covering both a rejected impossible token and an accepted real
  token. A control without its own regression test is not a control.
- Correction notes **restate history; they never rewrite it.** Do not edit the
  original recorded lines out of existence and do not re-run historical gates at
  old commits. Mark the affected lines `UNVERIFIED` and point at task-3c's
  current-tree baseline.
- Do not weaken any existing report-lint check to make the new one pass.

## Scope

### Files To Create Or Modify

- `project-docs/tools/report-lint.sh` (check plus `--self-test` cases).
- `project-docs/docs/phases/phase-20/tasks/task-1-done.md`,
  `task-2-done.md`, `task-3-done.md`, `task-3a-done.md`, `task-3b-done.md`
  (appended correction notes only).
- `AGENTS.md` only if the new check needs its rule stated more explicitly.

### Files To Read For Context

- `project-docs/docs/phases/phase-20/decisions.md` DEC-131.
- `AGENTS.md` (Commands), especially the evidence-token and `UNVERIFIED` rules.
- `project-docs/docs/phases/phase-19/decisions.md` DEC-102 for the correction
  precedent.

### Out Of Scope

- QA review artifacts. Extending the control to reviews is CF-029.
- Any product code, test, or contract.
- Re-running or re-litigating task-3c's gate results.

## Acceptance Criteria

- [ ] `report-lint.sh` fails a report that records `exit 0 (no output)` as the
      token for any command containing `go test`, naming the offending line.
- [ ] `report-lint.sh` fails a `go test` line whose token does not look like
      `go test` output; a token containing `ok `, `FAIL`, `PASS`, `---`, or
      `no test files` is accepted. `task-3b-done.md:70`'s
      `  295 internal/config/project_config.go` is a worked example that must
      now fail.
- [ ] Commands that legitimately print nothing (`gofmt -l`, `go build`) still
      accept `exit 0 (no output)`; no existing check is weakened.
- [ ] `--self-test` gains cases for the rejected impossible token, the rejected
      mismatched token, and the accepted real token, and passes.
- [ ] Each of the five Phase 20 completion reports carries an appended
      correction note marking its full-suite `go test` line(s) `UNVERIFIED`,
      stating why the recorded token could not have come from that command, and
      citing task-3c's current-tree run as the phase's verified baseline.
- [ ] Reports whose task-level verdict changes are stated plainly: the task
      remains complete, but with gates unverified as of its own commit.

## First Failing Test

`project-docs/tools/report-lint.sh --self-test` - the new case asserting that a
`go test -race ./...` line recorded `exit 0 (no output)` fails lint.

## Documentation Updates

- `AGENTS.md` (Commands) only if the token rule needs sharpening to match the
  implemented check. Keep it to the existing rule's wording; add no fid product
  feature to workflow docs.

## Exact Commands To Run

Run exactly:

```text
project-docs/tools/report-lint.sh --self-test
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. The Go gates are required here
because this task edits committed repository files; record their real tokens.

List any conditional command required by this task and why:

- `project-docs/tools/report-lint.sh <report> <spec>` re-run against each
  corrected report, to confirm the corrections themselves lint clean.

## Notes For Developer

- Match the token against the recorded command text, not against a hardcoded
  list of task names, so the check keeps working for future phases.
- Prefer a clear failure message that quotes the offending command and token;
  this control exists to be read by a human under time pressure.
- If the check would flag a legitimate pattern you did not anticipate, record it
  as a deviation and stop rather than adding an escape hatch.
- Completion report: `task-3d-done.md`.
