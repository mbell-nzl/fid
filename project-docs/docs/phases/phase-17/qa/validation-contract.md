# Phase 17 Validation Contract - Maintenance, Stub Guards, and v0.1 Closure

> Phase: 17 - Maintenance + stub guard + v0.1 closure (roadmap v0.1e,
> slice 4 of 4)
> Author: Tech Lead (draft)
> Last refined: 2026-07-10 by Tech Lead (kickoff approval pending)

## Purpose

This file is the finite set of assertions that must be true before v0.1 can be
closed. The contract specifies what must be proven; the QA plan specifies how.
If an assertion cannot be satisfied, stop and escalate rather than weakening it.

## Scope Decisions Pending Human Approval

- **DC-A - maintenance warning semantics.** Restore the detail retained in
  `scaffold/initial-spec.md` section 9.8.1 but missing from the canonical
  extracts: started-without-claim uses a project setting
  `[maintenance].started_without_claim_grace_seconds`, defaults to 300 seconds,
  permits zero, and rejects negative values. Age is measured from the newest
  append-ordered `status_changed` activity, falling back to `issues.created_at`
  when no transition exists; non-status updates do not reset it. Stale AC
  verification means a row whose non-NULL `completion_id` references a task
  completion with `status = 'superseded'`. A telemetry gap is a completion with no
  `agent_runs.completion_id` reference to that completion.
- **DC-B - maintenance result semantics.** `status` is `pass` with no findings,
  `warn` for warnings-only without `--warn-as-error`, and `fail` for any error
  or promoted warning. Promotion changes status/exit only; findings stay in the
  `warnings` array. Findings order by check, entity ID, then message. The CLI
  supports json/toon/markdown/yaml, following the canonical read-format policy.
- **DC-C - exact stub-guard set.** Guard every documented v0.2 command for the
  three v0.1 schema-only governance entities: `decision create/get/list`,
  `escalation create/resolve/get/list`, and `triage record/start/verdict`.
  Each returns `commands for <entity> ship in v0.2` non-zero before flag, role,
  DB, or journal processing. Shipped `requirement sync` remains operational;
  later requirement reads are not treated as stub-table entities.

Task execution is blocked until the Human approves these decisions and Heavy
Mode. Approval promotes every row below from `Draft` to `Ready` and the
decisions are recorded as phase decisions.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | CF-019; `08_QUALITY_PLAN.md` section 4.4 | The markdown formatter is split into focused non-test files at or below 400 lines with byte-for-byte unchanged behavior for every pre-Phase-17 markdown surface. | Source-state line counts plus existing and golden formatter tests before/after the split. | QA Reviewer | task-1 / review-1 | Every affected non-test file is <=400 lines and all existing markdown outputs are unchanged. | A file remains over budget or any existing markdown output changes. | Draft | New maintenance rendering belongs to task-5. |
| VAL-002 | PRD-OPS-008; `04_TRD_CONFIG.md` section 2; DC-A | Existing configs without `[maintenance]` read with a 300-second started-without-claim grace; explicit zero/non-negative values round-trip; negative or malformed values fail clearly. | Config unit tests over missing, zero, positive, negative, and malformed values plus bootstrap round-trip. | QA Reviewer | task-2 / review-1 | Missing=300, explicit values persist, invalid values error, and existing prefix parsing is unchanged. | Backward compatibility breaks, an invalid value is accepted, or a valid value is lost. | Draft | No general config framework or third-party TOML dependency. |
| VAL-003 | PRD-OPS-005/006; TRD-ARCH-006; `03_TRD_SCHEMA.md` sections 14.1/14.3 | The audit lists every active v0.1 issue hierarchy, sequence/scheduling, and dependent-row kind violation as an error, while a valid hierarchy yields none. | Table-driven domain tests seeded through a trigger-bypassing fixture, one invalid and one valid state per rule family. | QA Reviewer | task-3 / review-2 | Every seeded violation is listed with check/message/entity ID; valid counterparts are excluded. | A seeded violation is missed, misclassified, duplicated, or a valid row is flagged. | Draft | Later-version tables not present in v0.1 are excluded. |
| VAL-004 | PRD-OPS-005/006; `03_TRD_SCHEMA.md` sections 10 and 14.2 | The audit lists unresolved v0.1 polymorphic link endpoints, claims, activity refs, and activity journal-event IDs as errors, while resolved refs yield none. | Domain tests for each supported ref type and journal present/missing states. | QA Reviewer | task-3 / review-2 | All orphan states are listed and all resolved states are excluded. | An orphan is missed or a valid reference is flagged. | Draft | v0.2+ polymorphic types remain version-gated. |
| VAL-005 | PRD-OPS-005/006; TRD-ARCH-007 | The audit reports every row returned by `PRAGMA foreign_key_check` and reports none for a valid DB. | SQLite integration test with foreign keys disabled only for fixture seeding, then audit under normal connection settings. | QA Reviewer | task-3 / review-2 | Each seeded FK violation is listed with table/row/parent identity; clean DB has none. | FK violations are missed, collapsed, or clean rows are flagged. | Draft | The audit does not repair rows. |
| VAL-006 | PRD-OPS-005/006; PRD-WF-005/006; `02_TRD.md` section 4.6 | The audit checks issue workflow/state pairs against `internal/workflow` definitions and lists invalid pairs, while every valid default pair yields none. | Domain tests over all default workflow states plus unknown workflow and unknown state fixtures. | QA Reviewer | task-3 / review-2 | All valid pairs pass; unknown workflow/state pairs are errors. | State lists are duplicated inconsistently or a valid/invalid pair is misclassified. | Draft | Configurable workflows are v0.2. |
| VAL-007 | PRD-OPS-007/008; PRD-CLAIM-001/007 | A claim at or past TTL is a `stale_claims` warning; a live claim is not. | Time-injected domain tests at before and on/after expiry boundaries. | QA Reviewer | task-4 / review-2 | Live state excluded; boundary and expired states included once. | A live claim is warned or an expired claim is missed. | Draft | Claim repair/force release is not performed. |
| VAL-008 | PRD-OPS-007/008; PRD-CLAIM-005; DC-A | A started-category issue with no live claim is excluded within the configured grace and warned at/after it; a started issue with a live claim is excluded. | Time-injected tests for live claim, no-claim before boundary, at boundary, after boundary, and grace=0. | QA Reviewer | task-4 / review-2 | Both qualifying and non-qualifying states match the configured boundary. | Grace is ignored, measured from the wrong state time, or a live-claim issue is warned. | Draft | The command does not acquire a claim. |
| VAL-009 | Roadmap Phase 17 maintenance scope; AC verification freshness; DC-A | An AC verification tied to a superseded completion is a `stale_verifications` warning; one tied to a current completion or no completion is not. | Domain tests with superseded/current/NULL completion references. | QA Reviewer | task-4 / review-2 | Superseded reference included once; current and NULL references excluded. | A stale reference is missed or current/NULL evidence is warned. | Draft | Content-hash evidence drift is not checked. |
| VAL-010 | PRD-OPS-007/008; PRD-TELEM-001..004; DC-A | A task completion without an `agent_runs.completion_id` reference is a `telemetry_gaps` warning; a linked run suppresses it. | Domain tests over completion without/with a linked run, including a run linked only by issue ID. | QA Reviewer | task-4 / review-2 | Only completions lacking the completion-level run link are listed. | Issue-only linkage suppresses the warning or a linked completion is warned. | Draft | `dev finish` cost capture remains v0.2c. |
| VAL-011 | `05_DESIGN_PLAN.md` section 6.5; DC-B | Audit info reports `issue_count`, `claim_count_live`, and `recent_activity_7d` from one consistent snapshot, and findings are deterministic by check/entity/message. | Seeded domain test with known counts and repeated-run equality. | QA Reviewer | task-4 / review-2 | Counts are exact and repeated results are deeply equal in deterministic order. | Counts drift, ordering changes, or repeated runs differ. | Draft | Broader entity-count inventory is deferred unless specs add it. |
| VAL-012 | PRD-OPS-005/006; PRD-OUT-001..004; DC-B | `fid maintenance check` renders the documented result shape in json/toon/markdown/yaml; a healthy workspace exits zero with `pass`, and structural errors exit non-zero with `fail` while listing every error. | CLI tests plus built-binary smoke over healthy and seeded-corrupt workspaces in all formats. | QA Reviewer | task-5 / review-2 | All formats render; healthy=zero/pass; corrupt=non-zero/fail with complete errors. | A format fails, status/exit disagree, or violations are omitted. | Draft | Automatic repair is not provided. |
| VAL-013 | PRD-OPS-007/008; DC-B | A warnings-only audit exits zero with `warn` by default and non-zero with `fail` under `--warn-as-error`, without moving or rewriting warnings. | CLI test runs the identical warning fixture with and without the flag. | QA Reviewer | task-5 / review-2 | Default and promoted states both match, with identical warning payloads. | Promotion is ignored, changes payload classification, or default exits non-zero. | Draft | Individual warning suppression flags are out of scope. |
| VAL-014 | PRD-OPS-003; ADR-0001 | Maintenance check refreshes a missing/stale projection before reading and never changes journal bytes or appends an event. | Integration test over stale/missing projection plus before/after journal hash and line count. | QA Reviewer | task-5 / review-2 | Fresh result returned; journal remains byte-identical and event count unchanged. | Stale data is served or the audit mutates the journal. | Draft | Scratch rebuild/diff audit is v0.2+. |
| VAL-015 | Roadmap `fid maintenance vacuum`; `02_TRD.md` section 4.3; `03_TRD_SCHEMA.md` section 13 | `fid maintenance vacuum --by <role>` runs SQLite VACUUM under the projection lock for `tl`/`human`, preserves projected data and journal bytes, reports before/after DB sizes, and rejects missing/unknown/denied roles before mutation; policy show matches enforcement. | Domain/CLI integration tests with authorized and denied roles, lock contention, state snapshots, and journal hashes. | QA Reviewer | task-6 / review-3 | Authorized vacuum completes with unchanged data/journal; all deny states are non-zero and non-mutating; matrix parity holds. | Authorization, locking, output, or preservation fails. | Draft | `maintenance migrate` remains unimplemented. |
| VAL-016 | PRD-ERR-012; roadmap stub-table guard; DC-C | Every approved decision/escalation/triage stub command returns the exact entity/version gate non-zero before validation and appends no journal/DB state; shipped `requirement sync` still succeeds. | Table-driven root CLI test over the exact command matrix plus before/after journal/DB counts and one requirement sync control. | QA Reviewer | task-7 / review-3 | All guarded commands fail with the right entity/version and no mutation; control succeeds. | Any stub is absent, validates first, exits zero, mutates state, or breaks requirement sync. | Draft | Actual governance CRUD remains v0.2b. |
| VAL-017 | PRD-OUT-001..004; v0.1 exit criteria | Every implemented v0.1 read command accepts json/toon/markdown/yaml; JSON, TOON, and YAML preserve the same machine-checkable data and markdown exposes the same primary identities without errors. | Integration-tagged command matrix over all v0.1 read surfaces with structural decode/conformance checks. | QA Lead | task-8 | Every matrix row supports four formats and deterministic parity checks pass. | A read surface/format is missing, fails, or loses required data. | Draft | Version-gated commands are not read surfaces. |
| VAL-018 | Roadmap lifecycle exit; PRD-AC-003/004; PRD-CLAIM-005/006; PRD-WF-001..004,007..009; PRD-ROLE-002/003 | A real workspace completes create -> claim -> transition -> complete -> unclaim with permission denial before append, AC verification, and concern enforcement observable. | Integration-tagged built-CLI journey with journal line counts and state assertions. | QA Lead | task-9 | Full primitive lifecycle succeeds; denied and unmet-rule paths are non-mutating. | A step cannot run, a guard is bypassed, or a rejection appends. | Draft | v0.2 governance workflow is excluded. |
| VAL-019 | v0.1 exit criteria; PRD-ISSUE-009..014 | Sequential and prioritised epics return the documented next ready task and exclude blocked/finished/non-current tasks. | Same closure journey with both scheduling modes and positive/negative ready states. | QA Lead | task-9 | Ready selection matches both scheduling contracts. | Wrong or blocked work is returned, or eligible work is omitted. | Draft | `ready --explain` is v0.2c. |
| VAL-020 | v0.1 exit criteria; PRD-LINK-001..008; PRD-ISSUE-015/023; PRD-GOV-006 | The closure journey exercises labels plus issue/requirement links, including auto-inverse and relation rejection, with observable read-back. | Same journey with label attach, requirement sync, valid link/inverse, and invalid relation rejection. | QA Lead | task-9 | Labels and both link directions read back; invalid relation is non-mutating. | Read-back/parity fails or an invalid relation writes. | Draft | Governance entity links remain version-gated. |
| VAL-021 | v0.1 exit criteria; PRD-TELEM-001..004 | The closure journey records priced run telemetry and cost/tokens stats report the expected non-NULL computed USD values. | Same journey using shipped pricing config, `run record`, and stats reads. | QA Lead | task-9 | Run snapshot and aggregate values match the pricing formula. | Telemetry is absent, NULL unexpectedly, or cost is wrong. | Draft | Provider reconciliation fields are optional. |
| VAL-022 | v0.1 exit criteria; PRD-OPS-003..008; ADR-0001 | The closure workspace passes maintenance, then delete-and-rebuild reproduces identical user-visible/projected state and maintenance results. | Snapshot all exercised state and maintenance output, delete DB, rebuild, and compare structural equality. | QA Lead | task-9 | Rebuild restores roles and all exercised state with equal post-rebuild audit. | Data is lost, rebuilt state differs, or maintenance changes unexpectedly. | Draft | Scratch journal-vs-DB diff is v0.2+. |
| VAL-023 | v0.1 exit criteria; PRD-INIT-001..004; PRD-OPS-001/004 | `fid init` creates a journal that can be tracked by git while `.fid/fid.sqlite` is ignored, and maintenance/vacuum/stub/read commands do not alter journal append-only history. | Temp git-repo integration test using generated `.fid/.gitignore`, git ignore/track checks, and before/after journal hashes. | QA Lead | task-9 | Journal is trackable, DB ignored, rebuild restores roles, and non-journal commands leave bytes unchanged. | Journal is ignored/untrackable, DB trackable, roles missing, or journal bytes change. | Draft | No automatic `git add` or commit. |
| VAL-024 | Roadmap v0.1 exit criteria; `01_PRD.md` every requirement row targeted to v0.1; `08_QUALITY_PLAN.md` sections 5.2 and 9 | Every v0.1-targeted PRD row has reproducible passing test/evidence mapping; full unit/integration/benchmark/build gates pass and the CGO-disabled binary is under 30 MB. | QA Lead phase-end traceability audit plus exact phase gate, static build, and size evidence. | QA Lead | phase-end | No v0.1 PRD row is unmapped; every required gate passes; binary <30 MB. | Any row is unmapped, a gate fails/skips, or the binary exceeds the limit. | Draft | Public-release threat model/scans are deferred to v0.3 per charter/quality plan. |
| VAL-025 | Charter early dogfood outcome; v0.1 exit criteria; `08_QUALITY_PLAN.md` section 5.3 | The Human uses the built CLI for real project work and judges the maintenance findings, version-gate errors, ready choice, and recovery path clear enough to accept v0.1. | Human checklist, one row per observable, after deterministic QA passes. | Human | phase-end | Every human row is completed and accepted or explicitly dispositioned before wrap. | A row fails/partial without disposition, or dogfood is not run. | Draft | No held-out scenario exists; this is ordinary real-work dogfood, not a hidden evaluator run. |

## Coverage Audit

- PRD-OPS-005/006 -> VAL-003..006, VAL-012; PRD-OPS-007/008 ->
  VAL-002, VAL-007..013; PRD-ERR-012 -> VAL-016.
- Roadmap maintenance vacuum -> VAL-015; roadmap read-format closure ->
  VAL-017; all v0.1 exit criteria -> VAL-018..025.
- `01_PRD.md` every v0.1-targeted requirement row -> VAL-024's phase-end
  traceability audit; Tasks 1-9 add narrower behavioral evidence where changed.
- CF-019 is scheduled because task-5 necessarily reopens markdown rendering;
  CF-008, CF-017, and CF-018 remain outside the packages this phase changes.
- Every task cites its owned IDs. VAL-024 and VAL-025 are honestly phase-end.
- Conditional assertions prove both states: VAL-002..010, VAL-012..016, and
  VAL-018..023. Parity assertions use deterministic deep/structural comparison:
  VAL-001, VAL-011, VAL-017, VAL-022, and VAL-023.
- Validator ownership is singular: review-1 owns VAL-001/002; review-2 owns
  VAL-003..014; review-3 owns VAL-015/016; QA Lead owns VAL-017..024; Human owns
  VAL-025.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-07-10 | VAL-001..025 | Created as Draft | Phase 17 planning | Human approval pending for DC-A/DC-B/DC-C and Heavy Mode |
