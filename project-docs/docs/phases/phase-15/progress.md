# Phase 15 Progress

> Phase: 15 - `fid dev *` composites + orchestration cleanup (roadmap v0.1e,
> slice 2 of 4)
> Status: COMPLETE - PASS. All 4 tasks done; final review PASS WITH IMPROVEMENTS
> (0 Blocking/Major); human checklist PASS; wrapped 2026-06-29.

## Plan

Phase 15 ships the developer **composite** commands on top of the existing
claim/transition/completion/`agent_runs`/`gate_runs` surfaces, plus the three
scheduled carry-forward cleanups on the same paths.

- `fid dev start`: one event = acquire caller claim + transition
  `todo → in_progress` + open one `agent_runs` row (`--model` required, DN-B).
- `fid dev gate`: one event = append one `gate_runs` row on the derived active
  attempt; records a reported result, does **not** execute the command
  (`--exit-code` required, `--output-excerpt` optional, DN-A).
- `fid dev finish`: one event = close open run + completion (`dev_complete`,
  gate snapshot) + transition `in_progress → dev_complete` + release claim.
- `fid dev block`: one event = completion (`failed_gate`, gate snapshot, reason)
  + transition `in_progress → failed_gate` + release claim.
- CF-013/014/015 cleanups on the same paths.

No schema migration and no projector change: the tables, triggers, generic
effect projector, and the `dev.*` section-13 matrix keys already exist.

## Task Index

| Task | Title | Status |
|---|---|---|
| 1 | Cleanup CF-014 (split `issue.go`/`replay.go`) + CF-015 (label fn extraction) | Complete; Review 1 CONTINUE |
| 2 | `internal/dev` + `dev start` + `dev gate` + four matrix rows | Complete; Review 1 CONTINUE |
| 3 | CF-013 shared builder + `dev finish` + `dev block` | Complete; validated at final review |
| 4 | E2E composite smoke + rebuild determinism + dogfood prep | Complete; validated at final review |

## Validation

Validation contract at `qa/validation-contract.md` includes VAL-001..012.
VAL-012 is the phase-end human dogfood of start/gate/finish/block observability.
QA plan at `qa/qa-plan.md`; one mid-phase review (`qa/review-1.md`) over a 2-task
window (Tasks 3-4 closed at final review).

## Progress History

- Kickoff: drafted VAL-001..012 and Tasks 1-4; confirmed DN-A/DN-B/PI-1..PI-6
  with the human; scheduled CF-013/014/015 into the phase; re-deferred
  CF-008 (v0.2) and CF-011 (Phase 16).
- Task 1 complete: behavior-preserving split — `cmd/fid/issue.go` 984→22 (largest
  sibling 228), `internal/journal/replay.go` 485→205/192/128; label add/remove
  write flows extracted ≤60 lines (VAL-001/002 PASS).
- Task 2 complete: new `internal/dev` boundary; `dev start` (claim + transition +
  open run, `--model`) and `dev gate` (records-not-executes, `--exit-code`/
  `--output-excerpt`); four `dev.*` matrix rows; `02_TRD.md` /
  `06_DESIGN_PLAN_CLI.md` synced (DEC-037/038/039; VAL-003/004/005, VAL-006
  partial).
- Review 1 (Tasks 1-2): CONTINUE; 0 Blocking. Two Improvement notes (missing
  template H1 in two reports; `internal/dev`/`cmd/fid/dev.go` function-size debt)
  → TRG-001 RESUME.
- Task 3 complete: CF-013 shared `issues.BuildCompletionEffects` extracted
  (`fid issue complete` unchanged); `dev finish` (close run + completion +
  transition + release) and `dev block` (`failed_gate` completion + reason +
  transition + release); VAL-006 finish/block enforcement completed
  (VAL-006..010 PASS). `internal/dev` `Finish`/`Block` left above 60 lines →
  TRG-002 RESUME.
- Task 4 complete: `cmd/fid/e2e_v01e_dev_test.go` drives both composite paths,
  proves one event each, `qa_reviewer` denial, and delete-and-rebuild parity;
  `07_DESIGN_PLAN_FLOWS.md` reconciled (VAL-011 PASS; VAL-012 prep).
- Final review: PASS WITH IMPROVEMENTS, 0 Blocking/0 Major, 2 Improvement
  (carry-forward) notes. VAL-001..011 PASS; VAL-012 CARRY-FORWARD to the human
  checklist; green phase gate (total coverage 75.8%; integration + e2e smoke +
  benchmarks + build all PASS).
- Human checklist (HV-15a..HV-15d): complete, overall PASS (2026-06-29) →
  VAL-012 final outcome PASS.
- Wrap (2026-06-29): VAL-012 dispositioned PASS; DN-A/DN-B/PI-1..PI-6 recorded as
  DEC-037/038/039; spec-sync verdicts recorded for all specs; TRG-001/002
  audited; CF-013/014/015 closed; CF-016 (gate read surface) and CF-017
  (`internal/dev` function-size cleanup) opened. Phase COMPLETE.

## Carry-Forward

- CF-008 re-deferred to v0.2 (resolution-reason field).
- CF-011 re-deferred to Phase 16 (broad `activity_log` reconciliation); the
  composites add only the per-command rows their transition/claim writes already
  emit. Not worsened.
- CF-013 **CLOSED** (shared completion-effect builder, VAL-010 PASS).
- CF-014 **CLOSED** (`issue.go`/`replay.go` split ≤400 lines, VAL-001 PASS).
- CF-015 **CLOSED** (label write-path functions ≤60 lines, VAL-002 PASS).
- CF-016 **OPENED** (Phase 16 / v0.2c): first-class issue-scoped gate read
  surface (`fid issue get --section gates` recommended); `fid sql` /
  `fid journal inspect` are the Phase 16 interim path.
- CF-017 **OPENED** (v0.2c): behavior-preserving function extraction for
  `internal/dev` `Finish` (109) / `Block` (97) under the 60-line target.

## Escalations

- No escalations raised in Phase 15. No open escalations at wrap.
