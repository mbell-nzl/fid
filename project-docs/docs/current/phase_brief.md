# Phase Brief - Phase 8 (Task completions, v0.1c)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** Phase 8 (Task completions, roadmap §v0.1c) is **PLANNED AND
> APPROVED** — validation contract + task list + DEC-A..E confirmed by human
> (2026-06-20). Ready for Phase Runner to start Task 1. Phase 7 (Acceptance
> Criteria Surface) is CLOSED (VAL-001..010 all PASS). No task work is in flight.

## Phase 8 goal (roadmap §v0.1c — Task completions)

Ship the **task-completion write surface** (`fid issue complete`) on the existing
`task_completions` / `gate_runs` schema and read path, and replace the Phase 5
journal-fixture evidence for `has_completion_report` with real write commands.

**Exit:** an agent can record a completion attempt, observe attempt numbering +
supersession of prior attempts, see the gate snapshot rolled into the completion,
flag a human intervention, and satisfy a real `has_completion_report`-gated
transition.

## Confirmed decisions (human, 2026-06-20)

- **DEC-A — new-attempt status.** `fid issue complete` records the new attempt as
  `status = 'dev_complete'`; prior attempts move to `'superseded'`. `failed_gate`
  + the paired transition are the v0.1e `fid dev finish`/`block` composites' job.
- **DEC-B — attempt numbering.** Derived `MAX(attempt_number)+1` (1 when none)
  under `.fid/write.lock`; not user-supplied.
- **DEC-C — gate snapshot is fixture-seeded.** `fid dev gate` (the `gate_runs`
  writer) is v0.1e, so VAL-005 seeds `gate_runs` via a journal fixture. Real
  `fid issue complete` with no gate writer yet snapshots an empty set (correct);
  real-write linkage re-verified in v0.1e.
- **DEC-D — no `activity_log` for completions** (DEC-020 precedent; matches
  `create`/AC). Consistent timeline population deferred to v0.1e via CF-011.
- **DEC-E — completions reject on epics** before any journal append (CF-007;
  trigger is backstop only). Permissions = `RequireKnownRole` (matrix is v0.1d).
  Epics never need a completion: epic closure runs through the **epic workflow**
  (`todo→in_progress→wrapped→closed`, `in_progress→wrapped` gated by
  `all_child_tasks_done_or_cancelled`), not `task_completions`, so the rejection
  cannot strand an epic. The epic `wrapped`/`closed` path is verified in a later
  slice, not Phase 8.

## Validation contract index (`../phases/phase-8/qa/validation-contract.md`)

| VAL | Behavior | Owning task |
|---|---|---|
| VAL-001 | `fid issue complete` creates `TC-<hash>` row with tokens/commits/files/notes/role | task-1 |
| VAL-002 | Attempt 1→2 numbering; prior attempt retained as `superseded` | task-1 |
| VAL-003 | `--human-intervention`/`--intervention-note` both states | task-1 |
| VAL-004 | Epic target rejected before append; task/subtask succeeds | task-1 |
| VAL-005 | Gate snapshot `gate_runs`→`gate_results` + derived counts (fixture-seeded); no-gate = empty 0/0 | task-2 |
| VAL-006 | Real `has_completion_report` gating reject→pass (no fixtures) | task-3 |
| VAL-007 | Delete-DB + `fid rebuild` parity for completions + gate snapshot | task-3 |
| VAL-008 | Human dogfood: completion lifecycle reads/behaves as expected | phase-end |
| VAL-009 | CF-012: `issue.go` split into ≤400-line cohesive files, no behavior change | task-4 |

## Ordered tasks

1. **Task 1 — `fid issue complete` core write.** Attempt numbering, supersession,
   fields (tokens/commits/files/notes), human-intervention, epic rejection.
   (VAL-001/002/003/004)
2. **Task 2 — Gate snapshot.** `gate_runs`→`gate_results` JSON + `gates_passed`/
   `gates_failed`, fixture-seeded; no-gate = empty/0/0. (VAL-005)
3. **Task 3 — `has_completion_report` real-data gating + rebuild parity.**
   Replace Phase 5 fixtures with real commands; delete-DB rebuild parity. (VAL-006/007)
4. **Task 4 — CF-012 `issue.go` split.** Behavior-preserving split into ≤400-line
   files (à la Phase 6 CF-010). Do last, after the feature lands. (VAL-009)

## Carry-Forward Register (open)

| ID | Item | Disposition |
|---|---|---|
| CF-003 | Role-validation corrupt-tail tolerance | Before v0.1d role/permission work. |
| CF-008 | Resolution type / cancellation reason (DEC-015) | v0.2 (product). |
| CF-011 | Consistent `activity_log` timeline population (DEC-020) | v0.1e (with `fid activity list`). |
| CF-012 | `internal/issues/issue.go` 677-line split (DEC-022) | **Folded into Phase 8 Task 4.** |

## Standing commitments to honor

- Coverage (DEC-007): COV-001 `internal/issues|journal|policy|workflow` ≥ 80%;
  COV-002 `cmd/fid` ≥ 70%. Phase 7 ended `cmd/fid` 78.9%, `internal/issues`
  82.0%, `internal/journal` 80.3%, `internal/policy` 85.0%, `internal/workflow`
  95.5%.
- Validate-before-append (CF-007): all rejections occur before journal append
  (non-zero exit, clear error, zero new lines/rows).
- DEC-020/DEC-D: completions do not emit `activity_log` (CF-011 reconciles in v0.1e).
- Before release-readiness: decide where held-out scenarios live and how
  implementation agents are prevented from loading them (recommend only).
- Process nudge: paste exact command transcripts for any public-output change in
  completion reports — folded into each task spec.
- Task-spec authoring nudge (Phase 7): scope file-size acceptance to the task's
  allowed write set, not the whole package (avoids DEC-022-style scope BLOCKs).

## Next Checkpoint

Human: review and approve the Phase 8 validation contract + ordered task list
(and the five proposed decisions DEC-A..E). On approval, Phase Runner starts
Task 1.
