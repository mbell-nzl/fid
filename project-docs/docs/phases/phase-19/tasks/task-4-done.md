# Phase 19 - Task 4: `fid workflow assign` + assignment-aware issue creation - Completion Report

## Evidence provenance — read first (DEC-101)

The gates and VAL-059 commands below were executed **by the Human directly**,
outside a Developer session, during the ESC-19-005 recovery. **Verbatim terminal
tokens were not captured.** Results are recorded as **human-attested**, not as
token-cited evidence, and no token has been reconstructed or invented — TRG-19-001
in this phase was exactly that failure and it is not repeated here.

This does not meet DEC-084/086. It is dispositioned by **DEC-101**, which names
**QA Lead's independent re-execution at Review 2** as the control. QA Lead owns
all thirteen Task 4 assertions and re-runs every command. **If any command fails
under QA Lead, DEC-101's disposition is void and this task reopens.**

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

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` — `fid workflow assign`
  compatibility rules, error list, reassignment semantics (DEC-076, DEC-080)
- `project-docs/specs/03_TRD_SCHEMA.md` §2.1 — `workflow_assignments` named as
  the creation-time resolution source
- `project-docs/specs/05_DESIGN_PLAN.md` §3.1 — `workflow assign` flag surface
  verified against what shipped (DEC-092)

## Validation Contract Coverage

All thirteen recorded **PASS (human-attested)** under the provenance note above;
each is re-validated independently by QA Lead at Review 2.

- VAL-014 — first kind assignment: one event, one kind-scope row
- VAL-051 — reassignment replaces the row; superseded workflow stays defined
- VAL-052 — kind-scope effect shape: `insert` alone, then ordered
  `delete`+`insert`; replay from empty reproduces one row
- VAL-053 — same effect-shape rule for project scope, PK `{"scope":"project"}`
- VAL-015 — `--project` writes one project-scope row with NULL `kind`
- VAL-054 — kind without kind-scope assignment resolves to the project one
- VAL-055 — kind-scope assignment wins over an existing project assignment
- VAL-016 — six enumerated rejections each produce a distinct value-naming
  error; the compatible case succeeds
- VAL-056 — every rejection leaves journal bytes and row counts unchanged
- VAL-057 — successful assign writes no `activity_log` row (DEC-081)
- VAL-058 — with no assignment rows, all three kinds yield the v0.1 workflow
  and `status='todo'`
- VAL-017 — creation under an assignment yields the assigned workflow and the
  DEC-077 initial state, for a workflow with an `unstarted` state and one
  without
- VAL-059 — with no assignment rows, all six candidate captures are
  byte-identical to the committed baselines; see below

## VAL-059 evidence

**Baseline provenance — established and complete.** Reproduced from the pinned
Task 3 implementation commit
`50533e2b813bf073cd75f75eb70dd09cd15e9040` in a detached worktree, captured
through the committed `capture.sh`, and checked against the committed
`SHA256SUMS`: **all six entries `OK`**, confirmed twice (the authorized
2026-08-01 Step R rerun and the 2026-08-02 human re-verification). Under DEC-099
reproduction from pinned inputs *is* the provenance test. Per DEC-101 this
one-time establishment is now closed and is not repeated.

**Candidate parity — the standing assertion.** The candidate binary
(`go build -o /tmp/fid ./cmd/fid`) captured through the same committed
`capture.sh` into a fresh scratch root, checked against the committed manifest:
**all six entries `OK`**. Byte-identical, so the no-assignment path is unchanged
from v0.1.

Committed digests:

```text
e04e116c6df34cb85825dc4aa3b886a4edc56143409808f5ac7e4eb56fdbb5f4  issue-get-epic.json
0e71ea0d6ddbee4e9e62275d7f96aed873f23629945cbbe5eb55f109b05ed7f4  issue-get-task.json
ed7e271430d01257a3f92e8f59326e7c1ff562ce5ff16654def825c681635541  issue-get-subtask.json
8dd2f3cf2be2d9e9b78bce77a6121ada76b10b593fd75c91768881d21eecd9f4  issue-list-epic.json
febad62c2f5f9720a680689f464fd80ac6389426f6b63ba00571948ee9b3d9b5  issue-list-task.json
152bb980ee36e22ad1bbd008e0257a565810a568212c6d24afd3b48ce41570a6  issue-list-subtask.json
```

Normalization is the committed no-op `{"json_pointers":[],"replacement":null}`;
it was never widened, and no expectation was hand-written.

## Exact quality-gate results

All six recorded **PASS (human-attested)** — see the provenance note; tokens
were not captured, and QA Lead's Review 2 rerun is the control.

- `gofmt -l cmd internal` — PASS (human-attested), no files listed
- `golangci-lint run ./...` — PASS (human-attested)
- `go test -race -coverprofile=coverage.out ./...` — PASS (human-attested)
- `go tool cover -func=coverage.out` — PASS (human-attested)
- `go test -tags integration -race ./...` — PASS (human-attested)
- `go build -o /tmp/fid ./cmd/fid` — PASS (human-attested)

## Deviations from spec

1. **`capture.sh` path bug (DEC-100 harness deviation).** The first invocation
   failed with `capture.sh: 20: cannot create
   cmd/fid/testdata/v01_parity/issue-get-epic.json: Directory nonexistent`.
   Root cause: the script resolved a relative output root *after* `cd`-ing into
   its temp workspace. Fix: absolutize `output_root` and create it with
   `mkdir -p`. The repaired script reruns byte-identically and deterministically;
   recorded `PASS` with this note, per DEC-100.
2. **Step R cleanup was sandbox-blocked (ESC-19-005 → DEC-101).** `git worktree
   remove /tmp/fid-phase19-v01-baseline` and `git worktree prune` failed with
   `error: failed to delete '.git/worktrees/fid-phase19-v01-baseline':
   Read-only file system`, leaving the path registered-but-missing and
   cascading into four dependent commands. Retained verbatim here as required by
   TRG-19-005. DEC-101 rules cleanup out of the evidence set and requires unique
   per-attempt scratch paths; the recovery used `mktemp -d` accordingly.
3. **Evidence executed by the Human without token capture.** See the provenance
   note; dispositioned by DEC-101.
4. **Step R described the baselines as "committed" when they were untracked.**
   `cmd/fid/testdata/v01_parity/` was `??` throughout the recovery. The
   reproduction was unaffected (it regenerates from the pinned commit), but the
   spec described a state that did not exist. Resolved by this task's
   implementation commit, which commits the directory.

## Concerns, technical debt, or notes for the Tech Lead

- **CF-025 (DEC-101):** VAL-059 should become a permanent golden-file
  integration test in `cmd/fid` rather than a manual ceremony, reusing the
  `buildClosureFIDBinary` / `closureRunFID` pattern. Scheduled for Task 6; it is
  the last legal slot, as heavy mode accepts no carry-forward past phase end.
- **Review 2 carries unusual weight.** It is the sole independent verification
  of all thirteen Task 4 assertions and all six gates, not just a sampling pass.
- Unassign / clear-assignment remains unscheduled. `fid issue transition` into
  custom-named states remains blocked until Phase 20 (DEC-079).

## Dependency-budget note

- No dependency changes.

## Commit hash

- `32614bbbd50e683b61921e9ed6beabee73b46021`

## Commit subject

- `feat(workflow): add assign command and assignment-aware creation`
