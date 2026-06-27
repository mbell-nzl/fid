# Phase Brief - Phase 12 COMPLETE / Phase 13 planning pending

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.

## Status

**Phase 12 (Requirements + links, v0.1d slice 3/4) is COMPLETE - PASS**, wrapped
2026-06-27 on branch `phase-12`. No phase is currently active. The next action
is **Tech Lead planning of Phase 13 (Labels + v0.1d closure)**; this brief is
replaced when Phase 13 planning begins.

## Phase 12 Outcome (for reference)

Phase 12 delivered requirement pointers and the typed-link surface:
`fid requirement sync` (minimal), `fid link add` (auto-inverse + relation/
endpoint validation), `fid link remove` (declared + inverse in one event +
`activity_log` evidence), and `fid link list` (outgoing/incoming, four formats,
resolved titles), plus the `requirement.sync`/`link.add`/`link.remove` rows in
the Phase 10 permissions matrix. No schema migration was needed.

- All four tasks complete with recorded gates; Reviews 1-2 CONTINUE.
- Final review PASS; full deterministic gate green (coverage total 76.2%,
  domain packages >=80%, `cmd/fid` 77.9%, build OK).
- VAL-001..011 PASS; VAL-012 PASS via human checklist (overall PASS, HV-12a..d).
- DEC-029 (`requirement sync` hashes the parsed summary cell only), DEC-030
  (closed v0.1 relation registry + app-layer prechecks), DEC-031 (one
  `link_removed` activity event records both directions for PRD-LINK-006a).
- No escalations; TRG-001/TRG-002 audited (both RESUME correct).
- Spec-sync verdicts recorded for all specs; no spec edits required.
- Archives: `../phases/phase-12/` (tasks, qa, decisions, progress, triage-log,
  tl-escalations).

## Carry-Forward Open At Wrap

| ID | Item | Target |
|---|---|---|
| CF-008 | Resolution-reason field + cancel `--reason` | v0.2 |
| CF-011 | Broad `activity_log` reconciliation (link-removal slice done) | v0.1e |
| CF-013 | `complete.go` orchestration cleanup | v0.1e |
| CF-014 | Split `cmd/fid/issue.go` (984) + `replay.go` (485) under 400-line target | v0.1e |

## Next Phase Target - Phase 13 (Labels + v0.1d closure)

Scope from roadmap v0.1d Phase 13 (plan in detail at phase start):

- `fid label add` / `fid label list`: additive label CRUD.
- `fid issue label add` / `fid issue label remove`: per-issue assignment that
  does **not** replace the full set (`fid issue update --labels` stays
  replace-all).
- End-to-end v0.1d closure smoke over the full lifecycle: create -> claim ->
  transition -> complete -> unclaim with permission enforcement, links, and
  labels (the v0.1 headline exit criterion).
- Scope refs: PRD-ISSUE-015, PRD-ISSUE-023. Label schema/triggers already exist
  (v0.1a); expect CLI surface + matrix rows only, no migration.

## Planning Reminders For Phase 13

- Create `../phases/phase-13/qa/validation-contract.md` with `VAL-NNN`
  assertions before any task execution.
- Review the Carry-Forward Register (CF-008/011/013/014) at kickoff and
  schedule, re-defer, or cancel each.
- Add label-action rows to the in-code permissions matrix per-implemented-action
  (DEC-026 pattern); copy the `03_TRD_SCHEMA.md` §13 role booleans, edit the
  spec only if implementation proves it wrong.
- Every new write command refreshes the projection with the held-lock freshness
  helper inside `journal.WithLock` (DEC-028); one logical write = one event.
- Pin any genuine spec-wording ambiguity as a phase decision at kickoff rather
  than letting it recur through QA rounds (Phase 12 PRD-LINK-006a lesson,
  DEC-031).
