# Phase Brief - Phase 7 (Acceptance Criteria Surface, v0.1c)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** Phase 7 is **IN PROGRESS** — Tasks 1-2 implemented; Review 1 was
> **BLOCK** (ESC-005 + Task 1 report finalization). Way forward: **DEC-021**
> (append-order latest-per-role) implemented via fixup **Task 2b**. Task 2b hit a
> scope BLOCK (over-worded file-size criterion swept in pre-existing
> `issue.go`); resolved by **DEC-022** — criterion scoped to the ac.go split,
> `issue.go` deferred to **CF-012**. Task 2b is unblocked: Developer runs the
> full gate for the in-scope split + Task 1 report finalization commit. Phase 6
> (v0.1b) is CLOSED; archives under `../phases/phase-6/`.

## Phase 7 goal (roadmap §v0.1c — first slice)

Ship the **acceptance-criteria write surface** on the existing
`acceptance_criteria` / `ac_verifications` schema and read path, and replace the
Phase 5 journal-fixture evidence for the `all_acs_verified` transition rule with
real write commands:

- `fid issue ac add` — ordered ACs on tasks/subtasks; reject on epics (PRD-AC-001/002).
- `fid issue ac verify` — multi-role, append-only verification with evidence/status (PRD-AC-003/004/005).
- `fid issue ac list` — latest verification status per role, all read formats (PRD-AC-006, PRD-OUT-001–004).
- End-to-end: real AC verification drives `all_acs_verified`-gated transition + rebuild parity (PRD-WF-007, ADR-0001).

Out of scope (later v0.1c phases): `fid issue complete` + `task_completions`
(Phase 8); `fid concern raise/address/list` + derived concern flag (Phase 9).

## Ordered task list

| Task | Title | VAL coverage | Notes |
|---|---|---|---|
| 1 | `fid issue ac add` | VAL-001, VAL-002 | New write path; epic rejection before append (CF-007). |
| 2 | `fid issue ac verify` | VAL-003, VAL-004, VAL-005 | Append-only multi-role; status/evidence enums; completion_id NULL. |
| 2b | Deterministic latest-per-role fixup (ESC-005/DEC-021) | VAL-004, VAL-009 | Append-order (`rowid`) latest; de-sleep tests; commit Task 1 report finalization; split 431-line ac.go. |
| 3 | `fid issue ac list` | VAL-006, VAL-007 | Read-only; reuse latest-per-role helpers; format dispatch. |
| 4 | `all_acs_verified` real-data + rebuild parity (closure) | VAL-008, VAL-009 | E2E reject→verify→pass; delete-DB + `fid rebuild` parity. |
| — | Human dogfood | VAL-010 | phase-end checklist; QA Lead fans into per-observable rows. |

## Validation contract index

`../phases/phase-7/qa/validation-contract.md` (VAL-001..010):

- VAL-001 ordered `ac add`; VAL-002 epic rejection (both states).
- VAL-003 multi-role independence; VAL-004 append-only history; VAL-005 field recording + invalid-status rejection.
- VAL-006 latest-per-role `ac list`; VAL-007 four-format dispatch + JSON default.
- VAL-008 real-data `all_acs_verified` gating (reject→pass); VAL-009 rebuild parity.
- VAL-010 human AC-lifecycle dogfood.

## Planning decision (DEC-020, confirmed by human)

AC writes (`ac add`, `ac verify`) **do not emit `activity_log`** in Phase 7 —
consistent with `fid issue create` (creations don't log to the derived timeline
table; only `update`/`transition` mutations do, per Phase 5 CF-005). The journal
still records every AC write (rebuild is unaffected; VAL-009 proves it). Broad,
consistent `activity_log` timeline population is deferred to v0.1e via **CF-011**
(reconcile create + AC writes + any other unlogged creations when `fid activity
list` / PRD-OUT-017 is designed). No activity VAL in Phase 7.

## Carry-Forward Register (open)

| ID | Item | Disposition |
|---|---|---|
| CF-003 | Role-validation corrupt-tail tolerance | Before v0.1d role/permission work (not Phase 7). |
| CF-008 | Resolution type / cancellation reason (DEC-015) | v0.2 (product). |
| CF-011 | Consistent `activity_log` timeline population across all writes (DEC-020) | v0.1e (with `fid activity list`). |
| CF-012 | `internal/issues/issue.go` 677-line tripwire split (DEC-022) | Next v0.1c slot; behavior-preserving (à la CF-010). |

CF-003 and CF-008 are not in Phase 7 scope; CF-011 is opened by Phase 7 DEC-020
and targets v0.1e. CF-012 is opened by Phase 7 DEC-022 (Task 2b scope correction)
and is **not** Task 2b work — `issue.go` is pre-existing debt, out of scope here.

## Standing commitments to honor

- Coverage (DEC-007): COV-001 `internal/issues|journal|policy|workflow` ≥ 80%;
  COV-002 `cmd/fid` ≥ 70%. Phase 6 ended `cmd/fid` 78.3%, `internal/issues`
  82.5%, `internal/journal` 80.3%, `internal/policy` 85.0%, `internal/workflow`
  95.5%.
- Validate-before-append (CF-007): all AC rejections occur before journal append
  (non-zero exit, clear error, zero new lines/rows).
- Before release-readiness: decide where held-out scenarios live and how
  implementation agents are prevented from loading them
  (`../phases/phase-6/held-out-recommendation.md` recommends the v0.1b smoke seed;
  Task 4's AC smoke may extend it — recommend only).
- Process nudge (Phase 6): require pasted exact command transcripts for any
  public-output change in completion reports — folded into each task spec.

## Next Checkpoint

QA Lead: refine `../phases/phase-7/qa/validation-contract.md` into a QA plan
mapping every VAL to deterministic `(qa-lead)` / mid-phase-review / `(human)`
evidence, fan VAL-010 into per-observable human rows, then start the dev/QA loop
on Task 1.
</content>
