# Phase 19 - Task 4: `fid workflow assign` + assignment-aware issue creation - Completion Report

## Summary of what was implemented

- `fid workflow assign --name <slug> [--kind <k>] [--project] --by <role>`
  appends one `workflow.assigned` event and projects into
  `workflow_assignments`. Exactly one of `--kind`/`--project` is required.
- Compatibility validation reuses the Task 2 assignment validator (DEC-080):
  `--kind <k>` requires `applies_to` equal to `<k>` or `any`; `--project`
  requires `any`.
- Reassignment is last-write-wins expressed literally (DEC-076): `insert` alone
  when no row exists for the scope, ordered `delete` + `insert` when one does.
  State is read inside the write lock so no speculative delete reaches
  `applyDelete`.
- `internal/workflow/resolve.go` resolves kind-scope → project-scope → built-in
  kind default. `internal/issues/create.go` no longer calls
  `defaultWorkflowID(kind)` and no longer writes the `"todo"` literal; initial
  state comes from the Task 2 helper (first `unstarted` state, else first
  declared state, DEC-077).
- Resolved workflow and status are written into the event effect, as before.
  No resolution logic was added to the projection or replay path.
- `cmd/fid/testdata/v01_parity/` committed: `capture.sh`, `fixture/`, the six
  baseline `.json` outputs, the no-op `normalization.json`, and `SHA256SUMS`.

## Files created or modified

- `cmd/fid/workflow.go` — `assign` subcommand
- `cmd/fid/workflow_assign_test.go` (new)
- `internal/workflow/assign.go` (new) — write path
- `internal/workflow/assign_test.go` (new)
- `internal/workflow/resolve.go` (new) — assignment resolution
- `internal/issues/create.go` — resolved workflow + initial state
- `internal/issues/issue_test.go` — assigned and unassigned creation coverage
- `cmd/fid/testdata/v01_parity/` (new) — parity fixture, harness, and goldens
- `project-docs/docs/phases/phase-19/tasks/task-4-done.md` — DEC-103 closure
  evidence only; no product code was changed in this session.

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` — `fid workflow assign`
  compatibility rules, error list, reassignment semantics (DEC-076, DEC-080)
- `project-docs/specs/03_TRD_SCHEMA.md` §2.1 — `workflow_assignments` named as
  the creation-time resolution source
- `project-docs/specs/05_DESIGN_PLAN.md` §3.1 — `workflow assign` flag surface
  verified against what shipped (DEC-092)
- This completion report — fresh DEC-103 closure evidence and report-lint-ready
  gate records.

## Validation Contract Coverage

Closure-session scope is fresh command evidence, not a replacement for QA
Lead's assertion-level validation. The fresh `go test -race -coverprofile=coverage.out
./...` gate passed with terminal token `ok github.com/mbell-nzl/fid/internal/workflow
(cached) coverage: 80.4% of statements`; its tests include the cited Task 4
coverage below. The earlier human-attested results remain history only, and no
QA Review 2 terminal token is transcribed here.

- VAL-014 — PASS — `TestWorkflowAssignProjectsOneKindAssignment` verifies the
  built CLI's one appended `workflow.assigned` event and one kind assignment at
  `cmd/fid/workflow_assign_test.go:23-50`.
- VAL-051 — PASS — `TestAssignEmitsScopeAwareReplacementEffects` covers the
  second kind assignment's replacement effect at
  `internal/workflow/assign_test.go:54-75`; superseded-workflow listability was
  independently validated in the QA window and is not re-attested by this
  closure session.
- VAL-052 — PASS — `TestAssignEmitsScopeAwareReplacementEffects` checks the
  first kind `insert` and ordered kind `delete` + `insert` with the exact PK at
  `internal/workflow/assign_test.go:54-75`. DEC-103 records the missing
  permanent per-scope replay assertion for Task 6; this report does not claim a
  new assertion-level replay trace.
- VAL-053 — PASS — the `project` subtest checks the first project `insert` and
  ordered project `delete` + `insert` with PK `{"scope":"project"}` at
  `internal/workflow/assign_test.go:42-75`. The same DEC-103 Task 6 permanent
  replay-test follow-up applies.
- VAL-015 — PASS — `TestWorkflowAssignProjectsOneKindAssignment` is the built
  CLI event/row pattern at `cmd/fid/workflow_assign_test.go:23-50`; project-scope
  row evidence remains in the Task 4 QA validation window rather than being
  newly re-run here.
- VAL-054 — PASS — project fallback is asserted by
  `TestIssueCreateResolvesKindAssignmentBeforeProjectAssignment` at
  `internal/issues/issue_test.go:281-320`.
- VAL-055 — PASS — kind-over-project precedence is asserted by the same test at
  `internal/issues/issue_test.go:281-320`.
- VAL-016 — PASS — all six CLI rejection inputs, their non-zero stderr, and an
  accepted control are covered at `cmd/fid/workflow_assign_test.go:66-92`.
  That test does not yet retain pairwise-distinct/value-naming assertions;
  DEC-103 assigns that permanent-test strengthening to Task 6, so no stronger
  closure-session claim is made.
- VAL-056 — PASS — each VAL-016 rejection preserves journal bytes and leaves
  zero assignment rows at `cmd/fid/workflow_assign_test.go:74-91`.
- VAL-057 — PASS — the successful built-CLI assign leaves `activity_log` empty
  at `cmd/fid/workflow_assign_test.go:52-58`.
- VAL-017 — PASS — both initial-state branches under an assignment are asserted
  in `TestIssueCreateUsesAssignedWorkflowAndItsInitialState` at
  `internal/issues/issue_test.go:176-258`.
- VAL-058 — PASS — the no-assignment all-kind regression coverage remains in
  `internal/issues/issue_test.go`'s create fixtures and is additionally checked
  at the CLI boundary by this session's VAL-059 candidate capture below.
- VAL-059 — PASS — this closure session executed the corrected DEC-103
  candidate procedure and observed all six raw-byte manifest checks `OK` in
  `## Tests run`; baseline provenance remains the closed one-time Task 3
  reproduction recorded by DEC-101.

## Tests run

- `go test -race -coverprofile=coverage.out ./...` — PASS — `ok github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.4% of statements`
- Corrected DEC-103 VAL-059 candidate procedure — PASS:
  - `issue-get-epic.json: OK`
  - `issue-get-task.json: OK`
  - `issue-get-subtask.json: OK`
  - `issue-list-epic.json: OK`
  - `issue-list-task.json: OK`
  - `issue-list-subtask.json: OK`

  Executed exactly:

  ```bash
  manifest="$(git rev-parse --show-toplevel)/cmd/fid/testdata/v01_parity/SHA256SUMS"
  capture_root="$(mktemp -d)"
  go build -o "$capture_root/fid" ./cmd/fid
  cmd/fid/testdata/v01_parity/capture.sh "$capture_root/fid" "$capture_root"
  ( cd "$capture_root" && sha256sum -c "$manifest" )
  ```

  The manifest was checked from the fresh candidate capture root by absolute
  path. The prohibited repository-root and self-referential artifact-directory
  SHA commands were not run.

## Exact quality-gate results

This DEC-103 closure session re-established every gate with its own terminal
token. The original human-attested `golangci-lint run ./...` PASS was **false**
(DEC-102): Task 5 later found the Task 4-introduced unchecked `rows.Close`.
That history is retained below as a deviation; it is not evidence for this
report. The fresh result is the PASS recorded here.

- `gofmt -l cmd internal` — PASS — `exit 0 (no output)`
- `golangci-lint run ./...` — PASS — `0 issues.`
- `go test -race -coverprofile=coverage.out ./...` — PASS — `ok github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.4% of statements`
- `go tool cover -func=coverage.out` — PASS — `total:									(statements)					77.0%`
- `go test -tags integration -race ./...` — PASS — `ok github.com/mbell-nzl/fid/internal/workflow	(cached)`
- `go build -o /tmp/fid ./cmd/fid` — PASS — `exit 0 (no output)`

## Deviations from spec

1. **`capture.sh` path bug (DEC-100 harness deviation).** The first invocation
   failed with `capture.sh: 20: cannot create
   cmd/fid/testdata/v01_parity/issue-get-epic.json: Directory nonexistent`.
   Root cause: the script resolved a relative output root *after* `cd`-ing into
   its temp workspace. Fix: absolutize `output_root` and create it with
   `mkdir -p`. The repaired script reruns byte-identically and deterministically;
   the harness correction is retained as required history.
2. **Step R cleanup was sandbox-blocked (ESC-19-005 → DEC-101).** `git worktree
   remove /tmp/fid-phase19-v01-baseline` and `git worktree prune` failed with
   `error: failed to delete '.git/worktrees/fid-phase19-v01-baseline':
   Read-only file system`, leaving the path registered-but-missing and cascading
   into four dependent commands. DEC-101 rules cleanup out of the evidence set,
   requires unique per-attempt scratch paths, and the closure candidate capture
   therefore uses `mktemp -d` with no cleanup evidence.
3. **Human-attested evidence episode (DEC-101, voided by DEC-102).** The six
   required gates, baseline reproduction, and candidate parity check were first
   executed by the Human outside a Developer session with no captured terminal
   tokens. They were correctly recorded as human-attested rather than
   reconstructed. DEC-102 subsequently established that uncaptured results are
   `UNVERIFIED`, never `PASS`.
4. **Original attested lint PASS was false (DEC-102).** The unchecked
   `rows.Close` in Task 4's new `internal/workflow/resolve.go` caused
   `golangci-lint run ./...` to fail; Task 5 repaired it in scope. This session
   re-ran the exact gate and establishes the fresh `0 issues.` PASS token above.
5. **Untracked-baselines wording.** Step R described
   `cmd/fid/testdata/v01_parity/` as committed while it was `??`. Reproduction
   was unaffected because it regenerated from the pinned Task 3 commit, but the
   description was inaccurate. Task 4's implementation commit subsequently
   committed the directory.

## Concerns, technical debt, or notes for the Tech Lead

- **CF-025 (DEC-101/DEC-103):** convert VAL-059 into the permanent Go
  golden-file test scheduled first and mandatorily in Task 6. This manual shell
  procedure remains provenance only after that test lands; no further manual
  VAL-059 run is authorized.
- DEC-103 also assigns permanent assertion strengthening to Task 6: pairwise
  distinct/value-naming CLI rejection checks for VAL-016 and per-scope replay
  row assertions for VAL-052/VAL-053. The behavior is QA-validated, but those
  durable tests must replace the review trace.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `32614bbbd50e683b61921e9ed6beabee73b46021`

## Commit subject

- `feat(workflow): add assign command and assignment-aware creation`
