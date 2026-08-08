# Progress - Current

> **Purpose.** In-flight phase status only. Historical progress lives under
> `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**No phase is active.** Phase 19 is **COMPLETE / WRAPPED** (2026-08-08,
DEC-106). Phase 18 is complete (2026-07-25); v0.1 remains closed for early
dogfood (DEC-057). **v0.2a workflow authoring is delivered.** Gate D1
supersedes the planned Phase 20 configurable-role slice. No replacement phase
has been planned; Gate D roadmap disposition and backend-port reslicing are
next.

Branch: `phase-19-workflow-authoring`.

## Phase 19 Outcome

Six tasks, **heavy** mode, a **70-row** validation contract, and
**70 PASS / 0 FAIL / 0 BLOCKED / 0 CARRY-FORWARD**. QA Lead's subsequent final
review returned **PASS WITH IMPROVEMENTS**; the Human checklist is **9 PASS /
0 FAIL / 0 PARTIAL**. Shipped: `fid workflow define` / `assign` / `list`, the
`workflow_assignments` projected table with kind→project→built-in precedence,
the `escalation-default` seed, and assignment-aware issue creation.

Full history, task/commit map, and carry-forward detail:
`../phases/phase-19/progress.md`. Decisions: `../phases/phase-19/decisions.md`
(DEC-075..DEC-083, DEC-089..DEC-106).

## Quality Outcome

**Zero product defects across the phase.** Eight escalations, all in
verification apparatus or evidence records:

- **Four in one observable's shell ceremony** (VAL-059): ESC-19-003 baseline
  selector, ESC-19-004 `capture.sh` path bug, ESC-19-005 sandbox-blocked
  cleanup, ESC-19-007 a manifest command that never worked from the prescribed
  directory. **CF-025 retired the ceremony** into a standing Go golden test.
- **ESC-19-006** — Task 4's human-attested lint `PASS` disproven by a regression
  Task 4 itself introduced, misclassified as "pre-existing" from `git show HEAD`.
- **ESC-19-001/002** — an evidence-token defect and a test-only import cycle.
- **ESC-19-008** — final review; five findings, none in product code.

The **irregular human-executed Task 4 session** produced **five** distinct
defects across its two commits, `32614bb` and `ce4cf5b`: a false gate
attestation, a lint regression, two missing `Refs:` footers, and a report that
never passed report-lint. **Every other Developer and QA-Lead task commit in the
phase range carries `Refs:` and produced none** (the six remaining gaps in the
range are five `docs(tl)` governance commits and one raw human commit,
`ebbcd16`). That is the phase's strongest evidence for DEC-102 item 5.

## Standing Rules Produced

**Kept** (each now a sentence in a canonical spec, not a phase decision an agent
must remember to load): DEC-094, DEC-095, DEC-096 (test-packaging half),
DEC-097, DEC-099 (reproduction principle), DEC-100, DEC-101 items 1-3,
DEC-102 items 1 and 5, DEC-103 items 4-5, DEC-104 items 1-4.

**Retired at wrap** (DEC-106 item 6): DEC-098 entire, DEC-099's baseline half,
DEC-101 item 4, DEC-103 items 1-3, and DEC-096's `Effect`-only characterization
— all procedure-specific scar tissue from the ceremony CF-025 deleted.
DEC-104 proposed retiring DEC-101 items 1-2 as well; that was **overturned**,
because both cost one sentence and `t.TempDir()` satisfies them by construction.

Scored across the eleven rule-bearing decisions (DEC-094..DEC-104): **nine
keeps, three retirements, two splits** — five retired items once the split
halves are counted individually.

## Carry-Forward

- **CF-021, CF-025** — **closed** in Task 6.
- **CF-022** — partially taken in Task 5; renderer backfill still deferred.
- **CF-023** — open, **Medium**. Gained three drill cases, the `report-lint.sh`
  `UNVERIFIED` message, mechanical `Refs:` enforcement, and the DEC-106 item 5
  prevention (spec-sync obligations must land in the owning task's
  `## Documentation Updates`).
- **CF-024** — measured in Task 6: 100k rebuild `24.749 s` against `<30 s`
  (~17.5% headroom, up from ~8-10%). Improved, not closed; fix stays v0.2b.
- **CF-026a / CF-026b** — new, from the human dogfood: the `workflow show` read
  surface, and the workflow lifecycle model (design decision before scheduling).
- **CF-027** — new: `internal/issues/create.go` and `writeMarkdown`
  decomposition. Inherited debt the phase improved; sequenced behind CF-025's
  parity test.
- **CF-008 (v0.2d), CF-017 (v0.2e), CF-018, CF-020** — re-deferred, unchanged.

## Escalations

**None open.** ESC-19-001..008 are resolved and archived in
`../phases/phase-19/tl-escalations.md`.

## Next Action

1. **Finish Gate D roadmap disposition** — explicitly classify non-port
   features and ratify the replacement vertical-slice order.
2. **Plan the first backend-port slice** only after Gate D closes. Phase 20 is
   superseded; do not schedule `role_permissions`, `fid role create`, or
   `fid policy set` without new standalone-fid dogfood evidence and a decision.
3. Disposition the recommended Phase 19 deep retrospective and outstanding
   carry-forwards during that reslice rather than treating them as implicit
   prerequisites.
