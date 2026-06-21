# Phase 10 Decisions

> In-phase decisions (DEC-NNN / DEC-A..E). Per-doc spec-sync verdicts and the
> full decision history are finalized at phase wrap. Newest at top.

## DEC-025 — ESC-007 resolution: reconcile uncommitted task-report commit fields (2026-06-21)

- **Context.** QA Review 3 raised ESC-007 (TRG-001 BLOCK): the Task 1/2/3
  completion reports carry filled commit-hash/subject fields in the working tree
  but are **uncommitted**, so the committed reports still read "Pending commit".
  `08_QUALITY_PLAN.md` §10 requires the completion report to be committed by its
  author in the task commit or an immediately-following `docs(task): ...` commit;
  a task is not done while its report is uncommitted.
- **Verification (TL).** The working-tree commit fields are **accurate** — each
  matches git history by short hash and subject:
  - `task-1-done.md` → `716329c` `fix(policy): report corrupt tails during role discovery`
  - `task-2-done.md` → `6f4ce6b` `feat(policy): add permissions matrix and show command`
  - `task-3-done.md` → `5293fd2` `feat(policy): enforce write permissions on commands`
  No report content changes are needed; the only defect is that the three docs
  are unstaged/uncommitted.
- **Decision (do-now).** Reconcile via the §10-sanctioned alternative: a single
  immediately-following `docs(task): reconcile phase-10 task report commit fields`
  commit staging `task-1-done.md`, `task-2-done.md`, and `task-3-done.md`. This is
  a mechanical git action; the Phase Runner (or human) executes it — the Tech Lead
  does not run git. Once it lands, QA re-verifies Review 3's single blocking
  finding (committed report no longer says "Pending commit"); VAL-004..008 then
  close on evidence already gathered in Review 3.
- **Spec impact.** None. §10 already prescribes the committed-report rule; this is
  a compliance gap, not a spec change. No spec edits.
- **Recurrence note (for wrap deep-retro).** The drift recurred across all three
  tasks: the loop fills the commit-hash field in the working tree but skipped the
  follow-up `docs(task):` reconciliation commit each time. Flag at wrap whether the
  Phase Runner should always pair "fill commit hash" with an immediate
  reconciliation commit (or adopt a deferred-hash placeholder convention inside the
  task commit), so this stops surfacing as a phase-end blocker.

## Task 3 execution note (2026-06-21)

- **DEC-E executed as planned for permission retrofit fixtures.** Existing tests
  that exercised now-denied write combinations were moved to permitted roles
  without changing the behaviors under test:
  - `maintenance.rebuild`: `developer` -> `workflow_runner` in
    `cmd/fid/e2e_v01b_test.go`, `cmd/fid/e2e_v01c_ac_test.go`,
    `cmd/fid/e2e_v01c_complete_test.go`, `cmd/fid/e2e_v01c_concern_test.go`,
    `cmd/fid/init_test.go`, `cmd/fid/issue_ac_test.go`, `cmd/fid/issue_test.go`
    (all rebuild parity callers), and `cmd/fid/rebuild_test.go`.
  - Raw rebuild journal fixtures gained a seeded `workflow_runner` role so the
    permitted rebuild actor exists in the replay-only workspaces:
    `testdata/journal/replay-v1.golden.jsonl`,
    `testdata/journal/corrupt-tail.jsonl`.
  - `issue.transition.done` invalid-edge fixtures: `developer` -> `tl` in
    `cmd/fid/issue_transition_test.go`, `internal/issues/transition_test.go`,
    the synthetic precondition-conflict event in `cmd/fid/issue_test.go`, and
    `testdata/journal/conflicting-precondition.jsonl`.
  - `issue.transition.qa_review` setup fixture: `developer` -> `qa_reviewer` in
    `internal/issues/transition_test.go`.
  - Epic `issue.transition.in_progress` setup fixture: `tl` ->
    `workflow_runner` in `internal/issues/transition_test.go`.
  - `issue.transition.cancelled` task/subtask fixtures: `developer` -> `tl` in
    `internal/issues/transition_test.go`.
- **Direct-prerequisite matrix rows were added to finish enforcement coverage for
  already-implemented derived actions.** Phase 10 Task 2's initial matrix missed
  the epic transition states `issue.transition.wrapped` and
  `issue.transition.closed`, plus the denied AC-verifier roles
  `issue.ac.verify.tl_lite` and `issue.ac.verify.workflow_runner`. Task 3 added
  those rows so every existing derived write action resolves to a matrix cell and
  denied writes surface `PermissionDeniedError` instead of a missing-entry error.

## Task 1 execution note (2026-06-21)

- **DEC-D executed as planned for CF-003.** `internal/policy` role discovery no
  longer tolerates malformed or un-upgradeable trailing journal lines after
  seeded roles are found; `KnownRoles` now reports them with the same journal
  file/line context the replay path uses, so the upcoming permission matrix
  builds on a sound role set instead of a partially scanned one.
- **The TRG-002 test rewrite is intentional and decision-backed.**
  `TestKnownRolesTRG002CorruptTailTolerance` now asserts reported corruption for
  a valid role followed by malformed/un-upgradeable tail input. This closes the
  acknowledged Phase 1 debt rather than weakening coverage.

## Phase 10 planning decisions (DEC-A..E, human-confirmed 2026-06-21)

Recorded in the validation contract §Phase scoping notes and the phase brief.
Summarized here for the wrap audit; each task spec and QA artifact cites them by
DEC-A..E. The human accepted all five pre-execution.

- **DEC-A — enforcement seam in `internal/policy`, action key derived in the
  domain layer.** `policy.Authorize(workspaceRoot, role, action string) error`
  subsumes `RequireKnownRole` (preserving `MissingRoleError`/`UnknownRoleError`)
  and adds the in-code matrix lookup, returning `PermissionDeniedError{Role,
  Action}` (names role + action + "not permitted", PRD-ROLE-003). Action keys are
  derived in `internal/issues` / `cmd/fid/rebuild.go` from known
  kind/target-state/role. No new architecture boundary.
- **DEC-B — the in-code matrix is the single source of truth, shared by
  `Authorize` and `fid policy show`, and lists only implemented write actions.**
  Phase 10 seeds the v0.1a–c actions; Phases 11–13 append theirs alongside their
  commands. TRD-ARCH-010 completeness is tested against implemented actions.
- **DEC-C — `fid policy show` is a read command:** no `--by`, no permission
  check, no journal write, no `activity_log`; renders json/toon/markdown/yaml.
- **DEC-D — CF-003 reconciliation aligns role discovery with the journal
  corruption contract.** `KnownRoles`/`collectRoles` report a malformed/
  un-upgradeable journal line as corruption (non-zero exit, error naming the
  file) instead of tolerating it after seeded roles are found, matching the
  replayer (TRD-ARCH-002) and ADR-0001. The
  `TestKnownRolesTRG002CorruptTailTolerance` "tolerated" subtest is rewritten
  (decision-backed red/green), not weakened. Closes CF-003.
- **DEC-E — the retrofit moves existing test fixtures to permitted role/action
  combinations.** Enforcing the matrix changes existing command behaviour (e.g.
  `developer` may not create an epic; `qa_reviewer` may not transition to
  `dev_complete`). Fixtures using now-denied combinations move to permitted roles
  (`tl`/`human`/etc.), recorded as a deliberate, spec-mandated behaviour change,
  not test weakening. Every fixture role change is listed in the Task 3
  completion report. Validate-before-append (CF-007) holds on every denial.

## Validation Contract Outcome Audit (at wrap)

Final review recorded VAL-001..008 PASS (deterministic) and VAL-009
CARRY-FORWARD to the human checklist. The human checklist (`qa/human-checklist.md`)
is now complete: HUMAN-09a..HUMAN-09d all PASS, overall PASS (2026-06-21).

- **VAL-001..008 — PASS.** Confirmed against final-review evidence (QA-01..QA-17),
  Review 1/2/3 + addendum, and the three committed task reports. No FAIL/BLOCKED.
- **VAL-009 — PASS (was CARRY-FORWARD).** Resolved by the completed human dogfood
  checklist; matrix readability and denial-error actionability both confirmed by
  the human. No remaining carry-forward from the contract.

No `VAL-NNN` ended FAIL, BLOCKED, or unresolved CARRY-FORWARD; nothing to
disposition forward.

## Spec Sync Verdict (at wrap)

Recorded for every spec. `verified-no-change` is a positive review assertion, not
a default skip.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | verified-no-change | Phase 10 added a cross-cutting permission gate + `fid policy show` over existing scope; no identity / scope / constraint change. |
| `01_PRD.md` | verified-no-change | PRD-ROLE-001..005 MET as written (final review; Task 3 report records `verified-no-change`). No requirement reworded; configurable `role_permissions` / `fid policy set` (PRD-ROLE-006) remain deferred to v0.2. |
| `02_TRD.md` | verified-no-change | `policy.Authorize` subsumes `RequireKnownRole` in `internal/policy`; action keys derived in `internal/issues` / `cmd/fid/rebuild.go` (DEC-A). No new architecture boundary; CF-003 reconciliation aligns role discovery with the TRD-ARCH-002 corruption contract (Task 1 report). |
| `03_TRD_SCHEMA.md` | touched: see Task 3 execution note | §13 permissions matrix gained the four direct-prerequisite rows for already-implemented derived actions: `issue.ac.verify.tl_lite`, `issue.ac.verify.workflow_runner`, `issue.transition.wrapped`, `issue.transition.closed`. Verified present at §13 (lines 886/888/898/899). No table/column/trigger change; `role_permissions` table stays v0.2. |
| `04_TRD_CONFIG.md` | N/A | Phase 10 added no `.fid/` layout or config surface; the matrix is an in-code constant (DEC-B), not a config file. |
| `05_DESIGN_PLAN.md` | verified-no-change | `fid policy show` and the enforcement gate follow the documented command map and read/mutation flows; no flow or surface change. |
| `06_DESIGN_PLAN_CLI.md` | verified-no-change | §11 `fid policy show [--format json\|toon\|markdown\|yaml]` already matched the implemented read-only command and JSON row shape (Task 2/Task 3 reports record `verified-no-change`; verified at §11 line 1024). |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | Permission denials follow the validate-before-append reject flow (CF-007); permitted writes follow the journal→projection→read flow unchanged. |
| `08_QUALITY_PLAN.md` | verified-no-change | COV-001 (`internal/issues` 81.9%, `internal/journal` 80.3%, `internal/policy` 81.2%, `internal/workflow` 95.5%) and COV-002 (`cmd/fid` 77.8%) held at close. New policy files within §4.4 file/function budgets. No plan reword. |
| `09_AI_WORKFLOW.md` | verified-no-change | Workflow ran as specified (contract → QA plan → review-after-each-task loop → final review → human checklist). The Review 3 BLOCK→ESC-007→DEC-025→reconciliation-commit→addendum re-verify path is the supported hard-stop behavior. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 10 delivered the §v0.1d phase-slicing "Phase 10 — Permissions matrix + policy" cut; the remaining 11→12→13 slices are unchanged. |

## Triage-Log Audit (at wrap)

One TRG entry in `triage-log.md`:

- **TRG-001 (BLOCK → ESC-007) — task-report commit-field evidence uncommitted.**
  Confirmed correct in hindsight. Review 3 verified that committed Task 3 report
  still read "Pending commit" while the working tree carried the filled hash; per
  `08_QUALITY_PLAN.md` §10 a task is not done while its completion report is
  uncommitted. This is a real required-artifact gap, not a misclassification,
  pre-resolution, or accepted-decision case, so BLOCK (not autonomous RESUME) was
  correct. It escalated to ESC-007, the Tech Lead issued DEC-025, the Phase Runner
  landed the `docs(task): reconcile phase-10 task report commit fields` commit
  (`eeaf9fa`), and the Review 3 addendum re-verified clearance. No entry should
  have escalated that did not; none over-escalated.

Autonomous TL-Lite resume is permanently enabled (DEC-002); there is no warm-up
tally to maintain.

## Carry-Forward Disposition (at wrap)

- **CF-003 — CLOSED.** Role-discovery corrupt-tail tolerance reconciled by Phase 10
  Task 1 (DEC-D): `KnownRoles`/`collectRoles` now report a malformed or
  un-upgradeable trailing journal line as corruption with the same journal
  file/line context the replayer uses, instead of tolerating it after seeded roles
  are found. VAL-001 PASS; the `TestKnownRolesTRG002CorruptTailTolerance` rewrite is
  decision-backed, not weakening. Moved out of the open register.
- **CF-008 — re-deferred** to v0.2 (resolution type/reason product enhancement).
- **CF-011 — re-deferred** to v0.1e. Phase 10 added no new `activity_log`-projecting
  write; the v0.1d writes the spec projects to `activity_log` (claim renewal, link
  removal, label changes) arrive in Phases 11–13 and follow the existing
  per-command discipline. CF-011's broader reconciliation stays v0.1e.
- **CF-013 — re-deferred** to v0.1e, alongside the `fid dev finish`/`block`
  composites that also touch `complete.go`.

No new carry-forward items opened by Phase 10.

## QA Recommendation Disposition (at wrap)

Final review verdict was **PASS for deterministic QA** with 0 blocking / 0 major /
0 improvement findings. Mid-phase and final follow-ups dispositioned:

- **Review 1 improvement — reproduce corrupt-tail CLI evidence.** **Done** —
  final review QA-02 reproduced the corrupt-tail CLI denial; no further action.
- **Review 2 improvement — matrix row-cell hard-coding in show tests.**
  **Deferred** — non-blocking; revisit when the matrix gains rows in Phases 11–13.
- **Review 3 blocking — task-report commit-field drift (ESC-007).** **Done** —
  resolved by DEC-025 and the `eeaf9fa` reconciliation commit; addendum verified.
- **Final-review follow-up — refresh current docs + archive ESC-007.** **Done now**
  — `phase_brief.md`, `progress_current.md`, `decisions_active.md`, and
  `tl_escalations_current.md` refreshed in this wrap; ESC-007 moved to the phase
  archive; CF-003 closed.
- **DEC-025 recurrence note — commit-field drift recurred on all three tasks.**
  **Carried to deep-retro recommendation** below as a Phase Runner process nudge.

## Workflow notes

- **What worked well.** Settling DEC-A..E with the human before execution again
  removed mid-phase decision churn on a high-blast-radius change (a cross-cutting
  gate retrofitted onto every existing write surface). The CF-003 → matrix →
  retrofit ordering was correct: reconciling role discovery first (Task 1) gave the
  matrix a sound role set to build on. The single-source matrix (DEC-B) made the
  show-vs-enforce parity check (VAL-008) structural rather than a runtime guess.
  Review-after-each-task cadence matched the per-task validation windows and caught
  the Task 3 commit-field gap before final review.
- **What was inefficient.** The Review 3 BLOCK was a process/artifact gap, not a
  code defect: the loop filled the task-report commit-hash field in the working
  tree but skipped the follow-up `docs(task):` reconciliation commit on all three
  tasks, so a mechanical omission surfaced as a phase-end blocker and cost an
  escalation + re-review round-trip. The underlying code and VAL evidence were
  already sound.
- **Were tasks the right size?** Yes. Task 1 (CF-003) and Task 2 (matrix + Authorize
  + show) were cleanly bounded; Task 3 was correctly flagged high-blast-radius and
  was the right unit despite touching many fixture files — the fixture moves were
  mechanical DEC-E consequences, not scope creep. The Task 3 direct-prerequisite
  matrix rows were a legitimate in-scope prerequisite, not a missed task.
- **Did QA cadence feel right?** Yes. Review-after-each-task (over the default
  after-three) suited three distinct validation windows with a high-risk Task 3.
  The addendum review to clear the commit-field blocker was the right lightweight
  follow-up shape.
- **Token-efficiency opportunities missed?** The commit-field round-trip was
  avoidable: a deferred-hash placeholder convention or a runner step that always
  pairs "fill commit hash" with an immediate reconciliation commit would have kept
  it out of phase-end QA. DEC-A..E up-front remained the big efficiency win.

### Deep retrospective recommendation

**Light, process-only — no deep multi-session retro.** Phase 10's only friction was
the recurring task-report commit-field drift (DEC-025: it recurred on all three
tasks and escalated as ESC-007). One focused fix: give the Phase Runner / Developer
loop a standing convention so the completion report's commit hash is either filled
in the task commit itself (deferred-hash placeholder reconciled in the same commit)
or always paired with an immediate `docs(task): reconcile ...` commit — so this
stops resurfacing as a phase-end blocker. The Phase 9 deep-retro item is now
discharged: the append-order read-ordering rule is promoted to
`decisions_principles.md` at this wrap. No code or contract retro needed.
