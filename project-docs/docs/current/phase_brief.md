# Phase Brief - Phase 9 (Concerns + v0.1c closure)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** Phase 9 (Concerns + derived flag + v0.1c closure, roadmap §v0.1c)
> is **PLANNED — APPROVED** (validation contract + task list + DEC-A..E
> human-confirmed 2026-06-20). Phase 8 is CLOSED (VAL-001..009 PASS, human
> dogfood 5/5). Phase Runner may start Task 1.

## Phase 9 goal (roadmap §v0.1c — Concerns + v0.1c closure)

Ship the **concern write surface** (`fid concern raise`/`address`/`list`) on top
of the already-shipped concern schema, derived flags, and read path; prove the
derived flag + the `concerns_addressed` rule with real writes; and run the
end-to-end v0.1c closure smoke.

**Exit (roadmap):** an agent can raise a concern against a completion, see open
concerns light up the `issues_enriched` derived flag on issue reads, address it
and watch the flag clear, and close out the v0.1c lifecycle end-to-end over
ac add/verify/list → complete → concern raise/address/list.

## Key planning finding (what already exists vs Phase 9 work)

- **Already shipped:** the `concerns` table + index + `issues_enriched`
  `concerns_raised`/`open_concerns_count` (Phase 2); those derived columns
  already surface on `fid issue get` and `fid issue list` (incl. `--has-concerns`)
  and `fid issue get --section concerns` reads concern rows (Phase 4); the
  `concerns_addressed` rule function + `loadConcernStatus` DB loader exist.
- **Phase 9 adds only:** the CLI write/read commands (`fid concern
  raise/address/list`) and proves the existing read path + rule with **real
  writes**. Mirrors Phase 7 (ACs) and Phase 8 (completions).

## Planning decisions (human-confirmed 2026-06-20)

- **DEC-A — PRD-WF-009 proven at the rule-evaluation seam; no default edge.**
  No v0.1 default workflow edge gates on `concerns_addressed` (TRD §4.6;
  `internal/workflow/workflow.go` matches). The rule is supported but unreachable
  via `fid issue transition`. Phase 9 proves it via `evaluateTransitionRule` →
  `loadConcernStatus` → `ConcernsAddressed` with real concern data. Attaching a
  workflow edge is **out of scope** (changes the TRD §4.6 default graph; v0.2
  configurable workflows).
- **DEC-B — concern domain lives in `internal/issues`; no new boundary.**
  New files (`concern_raise.go`/`concern_address.go`/`concern_list.go`) +
  `cmd/fid/concern*.go`. The concern read model + rule loader already live in
  `internal/issues`, so no architecture-boundary decision is needed.
- **DEC-C — validate-before-append linkage (CF-007).** raise rejects (before any
  append) on unknown role, nonexistent `--on`, nonexistent `--in-completion`, or
  completion-not-on-issue. address rejects on unknown role, nonexistent concern,
  concern-not-`open`, or nonexistent `--in-completion`. Every rejection = non-zero
  exit + clear error + zero new journal lines/rows.
- **DEC-D — address is an in-place status mutation; raise→address only.**
  `open`→`addressed` updates the same row. `escalated`/`wontfix` (PRD-DEV-006
  statuses) have no Phase 9 writer (deferred v0.1e+).
- **DEC-E — no `activity_log` for concern writes (DEC-020/CF-011 precedent).**
  Matches create/AC/completion; consistent timeline population is v0.1e (CF-011).

## Validation contract index

`project-docs/docs/phases/phase-9/qa/validation-contract.md` (VAL-001..011):

- VAL-001 raise creates `CON-<hash>` open concern w/ all fields (task-1)
- VAL-002 raise validate-before-append rejections + valid append (task-1)
- VAL-003 raise lights derived flag false→true on get + list (task-1)
- VAL-004 address open→addressed in place, records completion (task-2)
- VAL-005 address validate-before-append rejections + valid mutation (task-2)
- VAL-006 addressing last open concern clears flag; partial keeps it (task-2)
- VAL-007 `concerns_addressed` rule seam: unsatisfied(named)→satisfied (task-2)
- VAL-008 `fid concern list` issue/status filters, all formats, empty=`[]` (task-3)
- VAL-009 delete-DB + `fid rebuild` parity for concerns + flags (task-3)
- VAL-010 end-to-end v0.1c closure smoke (task-3)
- VAL-011 human dogfood of concern lifecycle (phase-end; fan per-observable)

## Ordered task list (one session each)

1. **Task 1 — `fid concern raise`** (VAL-001/002/003). Write command +
   validate-before-append + derived-flag round-trip; read back via
   `fid issue get --section concerns`.
2. **Task 2 — `fid concern address` + rule proof** (VAL-004/005/006/007).
   In-place status mutation + flag clearing + `concerns_addressed` rule-seam proof.
3. **Task 3 — `fid concern list`, rebuild parity, v0.1c smoke** (VAL-008/009/010).
   Top-level list command + rebuild parity + end-to-end closure; closes v0.1c.

## Planning inputs honored

- CF-007 validate-before-append on every reject path (VAL-002/005).
- `RequireKnownRole` permissions only (matrix is v0.1d), matching Phase 7/8.
- ADR-0001 rebuild parity for concerns + derived flags (VAL-009).
- Human dogfood VAL for the concern lifecycle feel, per observable (VAL-011).
- Coverage (DEC-007): COV-001 `internal/issues|journal|policy|workflow` ≥ 80%;
  COV-002 `cmd/fid` ≥ 70% — enforced at the gate each task (Phase 8 ended
  `cmd/fid` 78.7%, `internal/issues` 81.9%).
- Task-spec nudges (carried): scope file-size acceptance to the task's write set;
  require finalized self-contained report evidence in the code-behavior commit.
- Keep concern logic out of `internal/issues/detail.go` (near the 400-line budget).

## Carry-Forward Register (open)

| ID | Item | Disposition |
|---|---|---|
| CF-003 | Role-validation corrupt-tail tolerance | Before v0.1d role/permission work. |
| CF-008 | Resolution type / cancellation reason (DEC-015) | v0.2 (product). |
| CF-011 | Consistent `activity_log` timeline population (DEC-020) | v0.1e (with `fid activity list`). |
| CF-013 | `internal/issues/complete.go` `Complete` orchestration size | v0.1e, alongside `fid dev finish`/`block`. |

## Next Checkpoint

Plan approved and DEC-A..E confirmed. Phase Runner starts Task 1
(`fid concern raise`, VAL-001/002/003). QA Lead refines the evidence strategy in
`qa/qa-plan.md` before execution.
