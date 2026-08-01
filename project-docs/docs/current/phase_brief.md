# Phase Brief - Phase 19 (In execution)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session.
>
> **Size target:** under 150 lines.

## Phase Goal

Phase 18 made workflow *structure* data. Phase 19 makes it **authorable**:
`fid workflow define` writes user workflows as journal events into the same
projected workflow tables, `fid workflow assign` binds a workflow to an issue
kind or to the project, `fid workflow list` reads them back in every supported
format, and `escalation-default` (PRD-WF-005a) is seeded as system reference
data. Success is observable: an agent defines a custom workflow, assigns it,
lists it, and creates an issue that picks it up instead of the built-in
default. v0.2a slice 2 of 3.

## Status

**Approved (DEC-093, 2026-07-26)** — mode, six-task list, refined contract
(VAL-001..070), CF-008 re-defer. QA-BLK-01..07 dispositioned by
**DEC-089..DEC-093** (`../phases/phase-19/decisions.md`); QA plan **Ready**.

**In execution. Tasks 1-4 complete. Next: Task 5, then Review 2.**

**Task 4 closed under DEC-101 (2026-08-02) with a named control.** Its gates and
VAL-059 commands were executed by the **Human without verbatim token capture**,
recorded honestly as **human-attested** rather than token-cited. **QA Lead's
Review 2 rerun is the control and must independently re-execute all six gates
and all thirteen Task 4 assertions; any failure voids the disposition and
reopens Task 4.**

**Review 1 (Tasks 1-3) verdict BLOCK, resolved (TRG-19-003 / ESC-19-003 /
DEC-098).** All 37 in-window assertions (VAL-001..013, VAL-027..050) are ON
TRACK with cited evidence and all Task 1-3 gates verified; the one blocking
finding was a **Task 4 precondition**, not a defect in shipped work — Task 4
said `git rev-parse HEAD` where QA-plan §6 requires the Task 3 *implementation*
commit, and the two demonstrably differed. **DEC-098** pins the baseline,
synchronizes Task 4 to QA-plan §6, and fixes normalization to the committed
no-op; `task-4.md` and `validation-contract.md` amended, `qa-plan.md` unchanged.

**Task 4 stop, resolved (TRG-19-004 / ESC-19-004 / DEC-099 / DEC-100).** The
first `capture.sh` run failed on a path bug in the script Task 4 was authoring.
Stopping was right — no commit, six gates unrun — but "restart from scratch" was
not: the failure wrote nothing, and DEC-098's guard is structural, not temporal.
**DEC-099** preserves the partial work and settles baseline provenance by
reproduction against the committed `SHA256SUMS`. **DEC-100** makes the general
rule standing.

**Task 4 third stop, resolved (TRG-19-005 / ESC-19-005 / DEC-101).** Step R's
`git worktree remove` / `prune` hit a read-only `.git/worktrees/` and cascaded
into four dependent commands. The stop matched the letter of Step R; the
reasoning was wrong. **Three consecutive stops, none in the product — all in
VAL-059's verification scaffolding**, each fixed narrowly so the next adjacent
brick fell. **DEC-101:** cleanup is never evidence; scratch state is unique per
attempt (`mktemp -d`); task specs declare required environment capabilities and
escalate shortfalls at step zero. Baseline provenance is settled by the six
`OK` and closed — not repeated. **CF-025** converts VAL-059 into a permanent
golden-file test in Task 6.

## Operating Mode

**Heavy** (DEC-075): two schema migrations with a projection-version bump;
three new public CLI commands plus two permissions-matrix rows; two new journal
event types that become permanent replay history. Criteria and artifact matrix:
`../../specs/09_AI_WORKFLOW.md` §4. Heavy never carries feasibility forward.

## Task Index

Specs and completion reports: `../phases/phase-19/tasks/`.

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | `workflow_assignments` schema + `escalation-default` seed | **Complete** (`a509d59`, `3b061e1`) | None |
| 2 | Workflow authoring core (validation + effect builder) | **Complete** (`10d4d24`) | 1 |
| 3 | `fid workflow define` + `workflow.defined` event | **Complete** (`50533e2`, report `697c673`) | 1, 2 |
| 4 | `fid workflow assign` + assignment-aware issue creation | **Complete** (human-attested, DEC-101; Review 2 is the control) | 1-3 |
| 5 | `fid workflow list` in every output format | Pending | 1, 3, 4 |
| 6 | Closure: integrity, rebuild parity, exit scenario, spec sync | Pending | 1-5 |

## Validation Contract

Sources, evidence, and pass/fail criteria:
`../phases/phase-19/qa/validation-contract.md`. **70 assertions, all `Ready`**
after DEC-091 (14 originals narrowed; VAL-027..070 added):

| Window | Assertions | # | Theme |
|---|---|---|---|
| task-1 | VAL-001..005, 027..031 | 10 | Assignments table + constraint enforcement, escalation seed, default parity, reversibility, projection version |
| task-2 | VAL-006..008, 032..042 | 14 | Deterministic effects, purity, rejection families, rule vocabulary, cancel-edge policy, initial state, assignment validation |
| task-3 | VAL-009..013, 043..050 | 13 | Define event/projection/load, no `activity_log`, reject-before-append, permissions, arrow aliases, output shape, `--applies-to` |
| task-4 | VAL-014..017, 051..059 | 13 | Assignment + effect shapes, precedence, assign rejections, creation resolution, v0.1 parity vs an immutable baseline |
| task-5 | VAL-018..020, 060..067 | 11 | Four-format list, markdown golden, filters, `assigned_to` shapes, ordering, empty result |
| task-6 (incl. phase-end) | VAL-021..025, 068..070 | 8 | Rebuild parity, maintenance integrity, seam + frozen-map inventories, DEC-079 boundary, exit scenario |
| phase-end (Human) | VAL-026 | 1 | Operator dogfood; QA Lead fans out into per-observable checklist rows |

## QA Artifacts

All under `../phases/phase-19/qa/`. `qa-plan.md` is **Ready**; it sets the
mid-phase cadence (`review-1` after Tasks 1-3, `review-2` after Tasks 4-5).
`final-review.md` / `human-checklist.md` are pending at phase end. Triage IDs
are allocated from `../phases/phase-19/triage-log.md`.

## Active Constraints

DEC-075..DEC-082 (planning), DEC-089..DEC-093 (QA-blocker dispositions), and
DEC-094..DEC-095 (Task 1 evidence correction, triage-log seeding) are recorded
in `../phases/phase-19/decisions.md`. Binding during execution:

- **`internal/workflow` owns authoring, assignment, and resolution** (DEC-089;
  `02_TRD.md` §3 synchronized). It still must not create issues, verify ACs,
  own migrations, or format output.
- **Issue-applicable workflows must declare exactly one `cancelled` state**
  (DEC-090, amending DEC-077): `--applies-to epic|task|subtask|any` with zero
  is rejected at define time, so PRD-WF-013 (a v0.1 Must) can never be violated
  for an issue. Zero stays legal for `escalation` / `qa_review` only.
- **Assignments are journal-canonical** (DEC-076): `workflow_assignments` is a
  projected table written by `workflow.assigned` events; seeded defaults stay
  migration DDL. **Precedence:** kind → project → built-in default; with no
  assignment rows, behavior is byte-identical to v0.1.
- **Cancel edges and initial state are derived, not hard-coded** (DEC-077).
  One `cancelled` state materializes edges for every non-cancelled state; zero
  materializes none (non-issue only, DEC-090); more than one is rejected.
  Initial state is the first `unstarted` state, else the first declared state.
- **Define accepts only the four implemented rules** (DEC-078);
  `required_gates_current`/`has_completion`/`linked_to_decision` are rejected
  with version-scoped messages. `linked_to_decision` is seeded for
  `escalation-default` only and stays known-but-unimplemented.
- **The DEC-066 seam stays behavior-free until v0.2b** (no seeded use, no
  evaluator, no query), and **no production Go path reads the frozen v0.1
  `definitions` map** (Phase 18 VAL-024; re-proven by VAL-069/070).
- **Custom state names are not transitionable until Phase 20** (DEC-079); do
  not auto-register matrix rows to work around it. **No `activity_log` rows for
  workflow authoring** (DEC-081). **Reject before append**: every invalid
  define/assign leaves journal bytes and projected rows unchanged.
- **Do not weaken TRD-PERF-003 (`<30 s`) or TRD-PERF-001 (`<500 ms`)**; CF-024
  leaves ~8-10% rebuild headroom. `internal/dev` `Finish`/`Block` stay over the
  60-line target (CF-017, v0.2e).
- **Evidence honesty and ID discipline are standing rules (DEC-084..086).** A
  non-terminating required command is never `PASS`; each result carries a
  verbatim terminal token. IDs are phase-scoped (`TRG-19-NNN` / `ESC-19-NNN`),
  allocated from `../phases/phase-19/triage-log.md` and written **before** the
  verdict.
- **No-match inventories exit 1** (DEC-094). VAL-031/042/069/070 are proven by
  `rg` finding nothing, which terminates `exit 1 (no output)`. Record that,
  never `exit 0 (no output)`.
- **`internal/workflow` may import `internal/journal`, bounded to the `Effect`
  struct** (DEC-096) — no append, projection, or replay surface. Consequently
  **no in-package (`package db`) test file may import `internal/workflow`**;
  tests needing both go in `package db_test`. Same rule for `internal/journal`
  in-package tests.
- **The VAL-059 parity baseline is pinned, never selected** (DEC-098). It is
  the Task 3 implementation commit
  `50533e2b813bf073cd75f75eb70dd09cd15e9040`, verified with `git rev-parse
  --verify` and cross-checked against `task-3-done.md`. **`git rev-parse HEAD`,
  branch names, and any other moving selector are prohibited.** Normalization is
  the committed no-op `{"json_pointers":[],"replacement":null}`; if the fixture
  cannot hold timestamps stable, Task 4 is **`BLOCKED`** — never widen it.
  **The guard is structural, not temporal** (DEC-099): the baseline is built in
  a detached worktree at the pinned commit from a committed fixed fixture, so
  working-tree edits cannot affect it and capture order is not binding.
  Provenance doubts are settled by **reproduction against `SHA256SUMS`**.
- **A task-authored evidence harness that fails while being authored is a
  deviation, not a failed gate** (DEC-100) — when the rerun is byte-identical,
  the pass is deterministic and third-party reproducible, the artifact is a task
  deliverable rather than the product under test or a pre-existing gate, and the
  report states the failure, root cause, and fix. Otherwise the strict rule
  holds. At triage, **"task incomplete, stop"** is not **"evidence tainted,
  restart"**.
- **Cleanup is never evidence** (DEC-101). Commands that only remove or prune
  scratch state are excluded from required-command evidence: never reported
  `PASS`/`BLOCKED`, never gating a sequence, and a failed cleanup is a note, not
  a stop. **Scratch state is unique per attempt** (`mktemp -d`) so no run
  depends on a prior run's cleanup — a fixed scratch path in a procedure is a
  spec defect. **Task specs declare required environment capabilities** (`/tmp`,
  `.git/`, network, elevation); a shortfall is escalated **at step zero**, never
  mid-sequence.

## Interfaces / Subsystems In Scope

`migrations/` (`0009` assignments, `0010` escalation seed);
`internal/workflow/` (authoring core, write paths, list read, resolution);
`internal/issues/create.go`; `internal/policy/matrix.go` (two rows — changes
`fid policy show`); `cmd/fid/workflow.go`; `internal/format/` (markdown);
`internal/maintenance/` (assignment integrity).

## Known Risks

- **Task 4 is the heaviest.** It changes `internal/issues.Create`, the most
  exercised path in the product. VAL-058/059 are the regression guard; capture
  the pinned baseline bytes **before** editing `create.go` (DEC-098 above).
- **Reassignment effect shape.** `applyDelete` requires exactly one affected
  row, so the event carries `delete`+`insert` only when a prior row exists.
- **Coverage.** `internal/workflow` is a COV-001 80% package gaining a lot of
  new code. **CF-023 residue (DEC-087):** no mechanical detector exists for a
  `PASS` against a non-terminating command; QA Lead's rerun is the control.
- **Review 2 carries unusual weight** (DEC-101). It is the sole independent
  verification of Task 4's six gates and thirteen assertions, not a sampling
  pass. Treat Task 4 as unverified until it clears.
- **CF-025 is due in Task 6** (DEC-101): convert VAL-059 into a permanent
  golden-file integration test. Heavy mode accepts no carry-forward past phase
  end, so Task 6 is the last legal slot.

## Open Questions

None — CF-008's re-defer to v0.2d (DEC-082) was confirmed at approval (DEC-093).

## Success / Exit Criteria

- [ ] All six tasks complete with green quality gates.
- [ ] An agent can define a custom workflow, assign it, list it in every
  format, and create an issue that picks it up instead of the built-in default
  (VAL-025); `escalation-default` is present and inspectable (VAL-003).
- [ ] VAL-001..070 all recorded with cited evidence; final review has no
  unresolved blocking finding; human checklist accepted (VAL-026 fanned out).
- [ ] Phase sign-off decision recorded.

## Next Planning Checkpoint

**Phase 20 — `role_permissions` migration + `fid role create` /
`fid policy set` + v0.2a closure**, carrying the DEC-079 unblock and the v0.2a
exit smoke.
