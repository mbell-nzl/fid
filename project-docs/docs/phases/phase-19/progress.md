# Phase 19 Progress

> Phase: 19 - `fid workflow define/assign/list` + escalation default workflow
> (roadmap v0.2a, slice 2 of 3)
> Status: COMPLETE / WRAPPED (2026-08-08). All six tasks closed, QA Lead
> subsequent final review **PASS WITH IMPROVEMENTS**, Human checklist 9/9 PASS,
> VAL-001..070 all PASS. No Developer task is active.

## Plan

Phase 18 made workflow *structure* data. Phase 19 made it **authorable**:
`fid workflow define` writes user workflows as journal events into the same
projected workflow tables, `fid workflow assign` binds a workflow to an issue
kind or to the project, `fid workflow list` reads them back in all four output
formats, and `escalation-default` (PRD-WF-005a) is seeded as system reference
data. Operating mode: **Heavy** (DEC-075) — two migrations with a
projection-version bump, three new public CLI commands, two permissions-matrix
rows, and two new journal event types that become permanent replay history.

## Task Index

| Task | Title | Status |
|---|---|---|
| 1 | `workflow_assignments` schema + `escalation-default` seed | **Complete** (`a509d59`, report `3b061e1`) |
| 2 | Workflow authoring core (validation + effect builder) | **Complete** (`10d4d24`; DEC-096 test-packaging repair) |
| 3 | `fid workflow define` + `workflow.defined` event | **Complete** (`50533e2`, report `697c673`) |
| 4 | `fid workflow assign` + assignment-aware issue creation | **Complete** (`32614bb`, report `ce4cf5b`, closure evidence session `ba65078` per DEC-103) |
| 5 | `fid workflow list` in every output format | **Complete** (`8f79013`) |
| 6 | Closure: integrity, rebuild parity, exit scenario, spec sync | **Complete** (`09aa0a0`, report `988737b`) — CF-025 first and mandatory, CF-021 closed, CF-024 measured |

## Validation

`qa/validation-contract.md` — **70 assertions, 70 PASS / 0 FAIL / 0 BLOCKED /
0 CARRY-FORWARD.** 69 QA Lead-owned deterministic assertions plus VAL-026
(Human dogfood, fanned out into HUMAN-26a..26i, 9 PASS). The contract was
refined from 26 to 70 `Ready` assertions at planning (DEC-091) to unbundle
independent observables and close coverage gaps. **No assertion, threshold,
fixture, or gate command was weakened at any point.**

Performance controls (not `VAL-NNN`): 10k-issue registry read `44.29 ms`
(<500 ms, 91% headroom); 100k-event rebuild `24.749 s` (<30 s, ~17.5% headroom,
2.2-2.9 s faster than the Phase 18 baseline — CF-024 improved but not closed).
Coverage total `76.8%` with every enforced package threshold met.

## Progress History

- **Planning (2026-07-26):** Heavy mode declared (DEC-075); DEC-076..DEC-083
  recorded; QA plan returned QA-BLK-01..07, dispositioned by DEC-089..DEC-093;
  contract refined to 70 `Ready` assertions; six-task pack and triage log seeded;
  Human approval recorded (DEC-093). CF-008 re-deferred to v0.2d.
- **Task 1 (2026-07-26):** `workflow_assignments` migration + `escalation-default`
  seed. TRG-19-001 → **DEC-094/095**: the report recorded a no-match `rg` as
  `exit 0 (no output)`; `rg` exits 1 on no match. Observable held; report
  corrected. DEC-095 makes the Tech Lead seed the triage log at phase start.
- **Task 2 (2026-07-29):** authoring core. TRG-19-002 → **DEC-096/097**: the
  binding `[]journal.Effect` return type makes `workflow → journal → db`, so an
  in-package `internal/db` test closed a cycle. Repaired by moving two tests to
  `package db_test`; no production change. DEC-097 root-causes TL-Lite's missing
  triage record as a Codex provider-prompt defect and fixes it.
- **Task 3 (2026-07-29/30):** `fid workflow define` + `workflow.defined`.
- **Review 1 (2026-07-30):** BLOCK. All 37 in-window assertions ON TRACK; the one
  blocking finding was a Task 4 precondition. TRG-19-003 → **DEC-098**: the
  VAL-059 parity baseline is pinned to `50533e2`, ambient `HEAD` prohibited,
  Task 4 synchronized to QA-plan §6.
- **Task 4, attempt 1 (2026-08-01):** BLOCK. TRG-19-004 → **DEC-099/100**:
  `capture.sh` resolved a relative output root after `cd`. Partial work
  preserved; baseline provenance settled by reproduction against `SHA256SUMS`.
  **DEC-100** makes a task-authored harness failure a deviation, not a failed
  gate.
- **Task 4, attempt 2 (2026-08-01/02):** BLOCK. TRG-19-005 → **DEC-101**: Step R
  `git worktree remove` hit a read-only `.git/`, and one failed cleanup gated
  four dependent commands. **Cleanup is never evidence**; scratch state is
  unique per attempt; task specs declare environment capabilities. Task 4's
  gates were then human-attested without token capture, with Review 2 named as
  the control. **CF-025** raised.
- **Task 5 (2026-08-02):** BLOCK. TRG-19-006 → **DEC-102**: the lint failure at
  `internal/workflow/resolve.go:61` was declined as "pre-existing" on the
  strength of `git show HEAD`; `git log --diff-filter=A` returns `32614bb`,
  Task 4's own commit. **Task 4's attested lint `PASS` could not have passed.**
  Uncaptured gate results are `UNVERIFIED`, never `PASS`; "pre-existing" is a
  provenance claim. Task 5 then completed at `8f79013`.
- **Review 2 (2026-08-02):** BLOCK. Independently re-ran all six Task 4 gates
  green and all five Task 5 gates; 25 of 26 in-window assertions ON TRACK.
  TRG-19-007 → **DEC-103**: the prescribed VAL-059 `sha256sum -c` never worked
  from the repo root and the obvious repair was self-referential; `task-4-done.md`
  had never passed report-lint. VAL-059 = **PASS** on executed evidence; Task 4
  closes via a Developer evidence session (`ba65078`); **CF-025 promoted to first
  and mandatory in Task 6**.
- **Prevention program (2026-08-03): DEC-104.** The Tech Lead gains a shell and
  executes any command it authors before publishing it; specify observables, not
  bespoke mechanism; second-escalation tripwire; evidence is an automated test
  unless the contract says why not; wrap-up records keep/retire per rule.
- **Task 6 (2026-08-03/06):** CF-025 executed first — VAL-059 is now
  `TestV01ParityGoldenOutputs`. Assignment integrity check, exit scenario,
  double-rebuild parity, reserved-rule seam checks, CF-021 closed, CF-024
  measured.
- **Final review (2026-08-02/06):** BLOCK on five findings. TRG-19-008 →
  **ESC-19-008** → **DEC-105**: VAL-026 completed by the Human (9 PASS, raising
  CF-026); the architecture contradiction resolved *against the decision, not
  the code*; stale canonical records accepted for wrap; the missing `Refs:`
  accepted as a deviation with enforcement moved to CF-023; the structure-budget
  finding **inverted on provenance** into CF-027. **None required a product code
  change.**
- **Human checklist (2026-08-07):** HUMAN-26a..26i all **PASS**. Raised the one
  product gap no deterministic assertion could find — workflow inspection depth
  and lifecycle management — split into **CF-026a** and **CF-026b**.
- **Subsequent final review (2026-08-07):** **PASS WITH IMPROVEMENTS.** Four
  findings closed on verified evidence; QA-BLK-FINAL-03 sustained in part with
  two documentary preconditions on the wrap commit; **MAJOR-01** raised — a
  false generalisation in `02_TRD.md` §3.1 that this phase introduced. Two of
  my DEC-105 evidence claims corrected while reaching the same outcomes.
- **Wrap (2026-08-08): DEC-106.** PRECON-01/02 and MAJOR-01 landed; MINOR-01/02
  onto the register; **a spec-sync obligation Task 6 never received was found
  and landed** (DEC-101/102/103 rules into `09_AI_WORKFLOW.md` and `AGENTS.md`);
  keep/retire verdicts recorded for the eleven rule-bearing decisions
  (DEC-094..DEC-104) — nine keeps, three retirements, two splits, i.e. five
  retired items counting split halves; full per-doc spec-sync verdict table
  completed; deep retrospective recommended.

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

The irregular human-executed Task 4 session produced **five** distinct defects
across two commits (`32614bb`, `ce4cf5b`): a false gate attestation, a lint
regression, two missing `Refs:` footers, and a report that never passed
report-lint. **Every other Developer and QA task commit in the phase range
carries `Refs:` and produced none.** That is the phase's strongest evidence for
DEC-102 item 5.

## Carry-Forward (final)

- **CF-021** — **CLOSED** in Task 6 (exported `maintenance.Vacuum` /
  `VacuumResult` documentation).
- **CF-025** — **CLOSED** in Task 6. VAL-059 is a permanent integration golden
  test; the capture script, manifest, normalization file, fixture, and goldens
  remain committed as provenance only.
- **CF-022** — partially taken in Task 5 (byte-exact golden for the new
  workflow-list markdown renderer); the backfill for pre-existing renderers
  stays deferred to the next `internal/format` reopen.
- **CF-024** — measured in Task 6: 100k rebuild `24.749 s` against `<30 s`
  (~17.5% headroom, up from ~8-10%). Improved, **not closed**; the batching fix
  stays v0.2b. Do not weaken the target.
- **CF-023** — open and now **Medium**. Gained three drill cases this phase
  (DEC-100 harness handling, DEC-102 unverified "pre-existing", and "a count is
  not a rate until the population is stated"), the `report-lint.sh`
  `UNVERIFIED`/`FAIL` message fix, mechanical `Refs:` enforcement, and the
  DEC-106 item 5 prevention (spec-sync obligations land in the owning task's
  `## Documentation Updates`).
- **CF-026a NEW** — `fid workflow show`: a named-workflow read surface rendering
  states, categories, ordered transitions, requirements, and assignments. Bounded,
  additive, no schema change. Target: v0.2a closure / Phase 20 planning candidate.
- **CF-026b NEW** — workflow lifecycle model (revision, deactivation,
  retirement, slug reuse). **Design decision before scheduling**; related to
  CF-008. Target: v0.2d governance work.
- **CF-027 NEW** — behavior-preserving decomposition of
  `internal/issues/create.go` (412 lines; `Create` ~228 lines) and
  `internal/format/markdown.go` `writeMarkdown` (68 lines). **Inherited debt the
  phase improved** — `create.go` was 419 lines at phase start. Sequenced
  **after** CF-025's parity test, which is the regression guard for that path.
- **CF-008 (v0.2d), CF-017 (v0.2e), CF-018, CF-020** — re-deferred, unchanged.

## Escalations

Eight escalations, **ESC-19-001..008**, all resolved and archived in
`tl-escalations.md`. Every one carried a real, unique, phase-scoped identifier
(DEC-084) — the ID-collision defect that affected Phases 17 and 18 did not
recur.

Triage: eight entries, all `BLOCK`, all correct to stop; three (004/005/006)
were wrong or incomplete in their stated reason, and all three traced to
defects in Tech Lead artifacts. Three of eight entries required a Tech Lead
backfill. Full audit in `decisions.md` → Triage-Log Audit.
