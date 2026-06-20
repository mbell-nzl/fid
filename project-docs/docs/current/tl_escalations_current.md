# TL Escalations - Current

> **Purpose.** Unresolved escalations that exceed local-decision authority and
> need full Tech Lead or human input. Default-loaded by every agent session.
>
> **Scope.** Open items only. Resolved escalations move to
> `docs/phases/phase-N/tl-escalations.md` at phase wrap-up.
>
> **Size target:** under 80 lines.

## Standing Escalation Triggers

### TL-Lite Triggers

An issue must be escalated, not decided locally by TL-Lite, if it:

1. Changes a public interface.
2. Changes the data model or persistence contract.
3. Overrides or contradicts a standing decision.
4. Affects more than one task or subsystem.
5. Introduces a dependency, tool, quality gate, or workflow.
6. Materially changes testing strategy or coverage policy.
7. Exposes ambiguity in product intent or acceptance criteria.
8. Resumes a context-exhausted or partial task where remaining or completed
   work spans schema/data model, dependency/platform configuration, public
   interfaces, canonical docs, or multiple Developer sessions.

### QA Reviewer Triggers

Most QA findings stay in `review-M.md`. QA Reviewer also escalates here if the
finding:

1. Has architecture implications beyond the three-task review window.
2. Implies acceptance criteria are ambiguous.
3. Has unclear security or privacy severity.
4. Contradicts a standing decision.
5. Spans tasks or subsystems beyond the review window.
6. Is a plausible but material spec deviation.
7. Is a quality-gate anomaly not anticipated by the QA plan.

## Open Escalations

### ESC-006 - Phase 9 final review blockers: benchmark gate failure and concern address read-back gap

Status: DIRECTION ISSUED (fixup scheduled — Task 4); closes at Phase 9 wrap once
QA verifies Task 4.
Source: Phase 9 final review / TL-Lite TRG-002
Raised: 2026-06-20
Resolved-by (direction): Tech Lead 2026-06-20

Summary:
- The required phase benchmark gate `go test -bench . -benchmem ./...` failed in
  `cmd/fid` at `TestConcernListFiltersByIssueAndStatus/all_concerns`.
- VAL-004 / QA-12 failed because after `fid concern address`, public read-back
  showed `status="addressed"` but neither `fid concern list` nor
  `fid issue get --section concerns` exposed `addressed_in_completion_id`.

Tech Lead diagnosis:
- Read-back gap: the write path (`concern_address.go`) persists
  `addressed_in_completion_id`, but `ConcernListItem` (`concern_list.go`) and
  `IssueConcern` (`detail.go`) never select/expose it.
- Benchmark flake: both concern read queries `ORDER BY created_at DESC, id DESC`;
  `created_at` second-ties fall back to the content-hash `id`, ordering rows
  nondeterministically vs creation order (same class as DEC-021 for AC verifs).

Resolution direction (issued):
- **DEC-023** — expose `addressed_in_completion_id` (nullable) on both concern
  read models.
- **DEC-024** — order concern reads by journal append order (`rowid DESC`,
  rebuild-stable; mirrors DEC-021).
- **Task 4** (`project-docs/docs/phases/phase-9/tasks/task-4.md`) implements both
  under TDD; required benchmark gate must exit 0 and VAL-004 must PASS while
  VAL-008/009 stay PASS. See `phases/phase-9/decisions.md` (DEC-023/DEC-024) and
  the validation-contract revision log.
- VAL-011 (human dogfood) remains BLOCKED pending the human checklist; it is not
  a Task 4 deliverable.

Next action: Phase Runner executes Task 4 via Developer → QA Reviewer; Tech Lead
closes ESC-006 and wraps Phase 9 after QA verifies the gate and VAL-004.

## Resolved (archived in their phase `tl-escalations.md`)

- Phase 8: no escalations raised (both reviews CONTINUE, zero BLOCKs/triage).
- ESC-005 (Phase 7 AC verify ordering nondeterminism) → DEC-021 / Task 2b.
- ESC-004 (Phase 6 `fid ready` tie-break order) → DEC-019.
- ESC-003 (Phase 5 VAL-005 rebuild-recovery evidence gap) → Task 1b.
- ESC-002 (Phase 4 due-date read-back gap) → Task 4.
- ESC-001 (SQLite pooled-connection pragmas) → DEC-006 / Phase 1 Task 7.
- Phase 4 journal-poisoning defect tracked as CF-007 / DEC-014 (scheduled
  fixup), not an open escalation.

<!--
Add new escalations above this comment using the template at
project-docs/templates/tl-escalations-template.md. Newest at top.
Use global ESC-NNN sequential numbering that never resets.
-->
