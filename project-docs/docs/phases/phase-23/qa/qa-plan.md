# Phase 23 QA Plan - Review And Human Judgment

> Phase: 23 - Review and Human judgment
> Refreshed: 2026-09-15
> QA Lead: Claude Opus 5 (activated directly by the Human; no Tech Lead handoff)
> Operating mode: Heavy. Every command and feasibility result below was produced
> fresh in this session on `phase-23` at `0c28936a6ad3980d2b3c4439b42ec70c83313d4c`
> (planning base `0195897a347dfe1872db0f7be3800c5c3ccdefcf` plus documentation
> commits only). Nothing carries forward from Phase 22.
>
> **Status: NOT READY - pending Tech Lead refinement.** This is not a
> placeholder `BLOCKED`: every required command is concrete and was executed or
> feasibility-proven in section 6. The plan cannot become Ready because
> thirteen Must-fix revisions (section 2.3) name contract contradictions,
> infeasible owning-task evidence, unowned scope, and unasserted clauses, three
> of which are blocker candidates for protected authority (section 2.4). No
> task may be dispatched until the Tech Lead dispositions R-01..R-13, applies
> the accepted edits, and a fresh QA Lead readiness refresh (section 11) passes.

## 1. Scope And Source Documents

Phase 23 delivers the Review and Human judgment slice: independent Review
Records with coverage and Findings bound to one Envelope revision and repository
identity; Builder `addressed` assertions versus independent resolution;
review-cycle budgets and rework exhaustion; generic Concerns and append-only
Deviations; Decision References and protected Decisions; provenanced Human
Outcomes and Attention resolution through them; requirement references replacing
V1 requirement authoring; and replacement-bound removal of the V1 judgment
surface. 18 sequential tasks, 74 assertions (70 deterministic, 4 Human), six
three-task review windows, one independent QA Lead checkpoint after task 11.

Sources read for this plan:

- `AGENTS.md`, `CLAUDE.md`
- `project-docs/docs/current/phase_brief.md`, `decisions_active.md`,
  `progress_current.md`, `tl_escalations_current.md`
- `project-docs/specs/08_QUALITY_PLAN.md` sections 4-8
- `project-docs/specs/09_AI_WORKFLOW.md` sections 2.6, 4.3, 4.4, 6.1
- `project-docs/specs/01_PRD.md` requirement table (IDs and PRD-REV/AUTH/HUM rows)
- `project-docs/docs/phases/phase-23/qa/validation-contract.md`
- `project-docs/docs/phases/phase-23/contracts.md`, `public-interface.md`,
  `schema.md`, `execution-scope.md`, `v1-convergence-inventory.md`,
  `contract-coverage.md`, `decisions.md`, `triage-log.md`
- `project-docs/docs/phases/phase-23/tasks/task-1.md` .. `task-18.md`
- `project-docs/docs/phases/phase-22/public-interface.md` (committed-denial rule,
  lines 252-257) and `qa/qa-plan.md` (format precedent only)
- `project-docs/templates/qa-plan-template.md`
- `project-docs/skills/qa-evidence/SKILL.md`,
  `project-docs/skills/validation-contract/SKILL.md`
- Source files inspected read-only are cited inline where they are evidence.

No held-out scenario, evaluator, or excluded-scenario material was read. Mode B
scenario execution is **N/A**: `AGENTS.md` records end-to-end and scenario
validation as `N/A` until a held-out harness location and command are recorded,
and none is. The `(human)` checks in section 9 remain required.

## 2. Validation Contract Refinement Check

- All in-scope requirements referenced by at least one `VAL-NNN`: **YES**. The
  contract's Coverage Rules map PRD-ISSUE-008, AUTH-005/007/008, EVID-001..003,
  EVID-005/006, REV-001..009, HUM-001..003, PROJ-001/002/003/005, REC-001..004;
  every one of those IDs exists in `01_PRD.md`. PRD-AUTH-007 is mapped to
  VAL-043 but only asserted in an out-of-scope note (I-06).
- Every `VAL-NNN` has one validator role: **YES**. 36 QA Lead, 34 QA Reviewer,
  4 Human; no row names two.
- Every `VAL-NNN` has a concrete evidence method: **YES as written; NO as
  executable** for VAL-008, VAL-009, VAL-045, VAL-053 at their owning task
  (R-07) and for VAL-046, VAL-057 under the accepted Event contract (R-02, R-06).
- Every `VAL-NNN` maps to a task/window or phase-end check: **YES**, verified
  mechanically. All 74 IDs appear in the contract; VAL-001..069, VAL-073, and
  VAL-074 are cited by task spec coverage sections; no task cites an ID absent
  from the contract. VAL-070..072 are phase-end Human rows cited by **no** task
  spec, which violates `09_AI_WORKFLOW.md` section 6.1 ("Every VAL-NNN is
  referenced by at least one task spec acceptance criterion") and leaves their
  fixtures unowned (R-10).
- Every `VAL-NNN` describes one observable: **NO for VAL-002** (textual diff plus
  CHECK enforcement, R-11). Matrix rows VAL-006, VAL-017, VAL-025, VAL-035,
  VAL-056, VAL-061 are acceptable only because section 7 enumerates their cases.
  Human rows VAL-070..073 are fanned out in section 9.
- Conditional or binary assertions require both states: **NO for five rows.**
  VAL-009 and VAL-046 cannot produce their qualifying state at the owning task
  or under the accepted effect rules (R-07, R-02); VAL-021's "still unsatisfied"
  half is vacuous at task 5 (I-07); VAL-033 tests only one supersession
  direction (I-05); VAL-015's non-qualifying state is ambiguous (R-04).
- Parity or "matches reference" assertions have deterministic proof: **YES**,
  with conditions. VAL-013 ("identically"), VAL-064, VAL-066, VAL-068 require
  structural comparison; section 4 pins the comparison inputs.
- Revisions requested from Tech Lead: R-01..R-13 Must-fix, I-01..I-16
  Improvement (section 2.3).

### 2.1 Independent contract-coverage re-derivation

`contract-coverage.md` was re-derived independently from `contracts.md`
sections 1-15, `public-interface.md` sections 2-6, and `schema.md` sections 1-7,
not by editing the Tech Lead's matrix. Result: every task 1..18 owns at least
one clause, and every schema table, index, and trigger has a consumer. The
following clauses have **no assertion that evidences them**, although the
matrix maps some of them to a VAL:

| Clause | Matrix mapping | Gap |
|---|---|---|
| `schema.md` 7: a corrupted or unknown Phase 23 Event payload is rejected rather than partially applied | "schema 7" -> VAL-004/005/064 | None of those rows asserts payload rejection (R-12a) |
| `contracts.md` 13: linking a Review that fails independence rejects `review-not-independent` | none | No row; a Reviewer who later takes a mutating Claim on the same Issue and revision is the reachable case (R-12b) |
| `contracts.md` 11 / `public-interface.md` 2: `review-recorded` and `finding-resolved` Condition types | "11: four Condition types" -> VAL-048 | VAL-048 exercises only `human-outcome-recorded` (R-12c) |
| `contracts.md` 10: `human-outcome-required` Stop producer | "10" -> VAL-046 | No accepted effect row can raise it (R-02) |
| `contracts.md` 9: 0..16 protected Decisions, unique class, affected-action subset | "9: 0..16" -> VAL-039 | VAL-039 tests blocking only; authoring validation unasserted (I-11) |
| `contracts.md` 4: referenced Receipts belong to the recording Run; `public-interface.md` 3: coverage items and exclusions may not share a digest | "4" -> VAL-014 | Negative cases absent from VAL-014 (I-06) |
| `schema.md` 4: `DELETE FROM links` / `DELETE FROM activity_log` for non-Issue rows; implicit drop of two kind-guard triggers | "15: owned removals" -> VAL-061 | Not enumerated; section 7 QA-28 adds them (I-12) |
| `public-interface.md` 2: ActivityEvent additions; ReasonCode `human-outcome-absent` | none | Unasserted (I-15) |

`contract-coverage.md` carries a duplicated row for "12: `links` restricted"
(lines 84 and 93); harmless.

### 2.2 Clauses verified sound

Applied DDL (section 6.4) confirms `schema.md` sections 1, 2, 4 and the literal
head of section 5 are executable as written on a fresh chain; the 0013/0015
hunks are exactly the seven accepted edits; the discriminator, one-gate-evidence
index, `claims.mode`, ten-member Stop and Attention code lists, and the
`concerns` self-foreign-key rewrite all behave as accepted. Documentation
amendment targets named by tasks 1, 8, 9, 12, 13, 17 all exist (section 6.8).

### 2.3 Revisions requested from Tech Lead

Must-fix (M) blocks Ready. Improvement (I) may be applied at the Tech Lead's
discretion. The Tech Lead owns every edit to `validation-contract.md`,
`contracts.md`, `public-interface.md`, `schema.md`, `execution-scope.md`, and the
task specs; this plan proposes and never applies them.

| ID | Sev | Where | Finding and proposal |
|---|---|---|---|
| R-01 | M | `schema.md` 2 (and 3); VAL-005; task 1 | **0017 cannot migrate a schema-16 workspace whose receipt Evidence has Artifact links.** The runner migrates on a `foreign_keys(1)` connection (`internal/db/db.go:17`) with one transaction per file (golang-migrate v4.19.1 `database/sqlite/sqlite.go:199-203`). `evidence_artifact_links.evidence_id` is an immediate FK to `evidence_records(id)` (`0013:293`), so `DROP TABLE evidence_records` performs an implicit DELETE that violates it. Measured on a consistent seed (FK check clean beforehand): accepted 0017 -> `Runtime error near line 265: FOREIGN KEY constraint failed (19)`, rc=1; prepending `PRAGMA defer_foreign_keys = ON` -> fails at COMMIT (line 274); staging the link rows in a TEMP table, deleting them, rebuilding, and restoring them -> rc=0, Evidence row identical by named column, link restored, `foreign_key_check` empty, `integrity_check` ok, link DDL byte-identical. Without link rows the accepted 0017 succeeds. VAL-005's seeded workspace is exactly the failing case in real use. Amend `schema.md` 2 (and the symmetric 3 down) with a link-preserving rebuild; do not narrow VAL-005 to link-free seeds. Blocker candidate: see 2.4. |
| R-02 | M | `contracts.md` 10; `schema.md` 6.2; VAL-046, VAL-056 | **No Event can raise the `human-outcome-required` Stop.** It must appear when the Human Obligation becomes the only unsatisfied one, which for any multi-Obligation Envelope happens at `evidence.linked`. `schema.md` 6.2 keeps `evidence.linked` "unchanged effect rows" (`internal/journal/current_event_catalog.go:88`: `evidence_records`, `evidence_artifact_links`, `idempotency_results` only). Only `review.recorded` and `envelope.activated` gain `stops/insert`. Name the raising Event(s) and add their effect rows (DEC-195 applied to an effect family), or restate the trigger. Blocker candidate: see 2.4. |
| R-03 | M | `contracts.md` 10 vs 13; VAL-025 vs VAL-044 | **Two codes for one input.** Linking a Human Outcome to a `gate-evidence` or `review` Obligation must reject `outcome-not-permitted` (section 10, VAL-044) and also `evidence-type-not-accepted` (section 13, VAL-025 "the six mismatches"). Two of VAL-025's nine cells contradict VAL-044. State precedence and amend whichever row loses. |
| R-04 | M | `contracts.md` 3; task 4; VAL-015, VAL-022, VAL-026, VAL-070 | **Independence is ambiguous for review Runs.** "must not equal any Actor that held a mutating Claim or a Run on the reviewed Issue within the same Envelope revision" does not say whether "a Run" includes a Run under a `review` Claim. Literally, a Reviewer whose first Review returned `changes-required` can never record the second Review that resolves the rework, and every cycle needs a new Actor; read narrowly, only Builder Runs conflict. The two readings produce opposite VAL-015/VAL-022 fixtures and different product behavior. Blocker candidate: see 2.4. Also clarify VAL-015 "in both directions". |
| R-05 | M | `public-interface.md` 5; VAL-039, VAL-056 | **New failure codes have no named producer or precedence.** `finding-unresolved`, `human-outcome-required`, and `protected-decision-required` are added both as failure codes and as ReasonCodes. The retained guarded-action denial is `obligation-unsatisfied` with a `reasons` array (`cmd/fid/envelope.go:236`, `internal/cli/codes.go:285`). No document says which operation returns each new code, or which code wins when a protected Decision and an unsatisfied Obligation apply together. VAL-056's "one reproduction per code" and VAL-039's "returns `protected-decision-required`" cannot be written deterministically. Name the producing operation(s) per code and the precedence order. |
| R-06 | M | `public-interface.md` 5; `schema.md` 6.3; VAL-057 | **The committed-denial replay has no producer.** Phase 22's rule (`phase-22/public-interface.md:252-257`) stores and replays a denial only when the denial itself durably raises the Stop. The `review-rework-exhausted` Stop is raised by the successful exhausting `review.recorded`; the `protected-decision-required` Stops are raised by `envelope.activated`. The denied `claim.acquire`, `run.start`, `issue.transition-terminal-success`, and `envelope.complete` requests have no denial Event in `schema.md` 6.2, so they commit nothing, return null EventRef, and store no snapshot. VAL-057's "denials replay their original Stop, Attention, and EventRef" is unproducible. Either add denial Events with effect rows, or restate VAL-057 for non-committing denials (null EventRef, no Event, existing Stop/Attention IDs, same result on retry). |
| R-07 | M | VAL-008, VAL-009, VAL-045, VAL-053; tasks 2, 10, 12 | **Owning task cannot produce the named evidence.** (a) VAL-008 (task 2) uses `evidence.evaluate` with no Obligation ID, but task 6 says the leaf "gains its optional Obligation ID here" and `public-interface.md` 4 lists it "null / 6,13"; task 2's own spec says `cmd/fid/evidence.go` "request and result shapes only". Evidence task 2 through the domain evaluator and `envelope.activate`'s `Evaluation[]` (row "2,9"), and add an explicit VAL-008 leaf re-check at task 6, or move the leaf change to task 2. (b) VAL-009's qualifying state (two or more applicable Obligations, all satisfied) is impossible at task 2: at most one `gate-evidence` Obligation exists and review or Human evidence links only from tasks 6/10. Add a named re-verification with gate + review + Human all satisfied at task 10 or the checkpoint. (c) VAL-045's Human-text clause is a `contracts.md` 14 obligation owned by task 14, and task 13 forbids pinning text task 14 replaces. Split: record/evaluation triple at task 10, text at task 14 (overlaps VAL-060). (d) VAL-053 asserts `last_synced_hash` absent from the latest schema at task 12, but it is a column of the V1 `requirements` table that migration 0018 drops in task 15 (task 12: "Until then the table exists"). Its "source search finds nothing" criterion also has no declared command and task 12 has no `## Expected Non-Zero Commands` section, so its non-zero PASS would be undeclared. Move column absence to task 15 (VAL-061) and declare the search in task 12 with expected `exit_code: 1` (feasible once R-09 is applied; current matches are only `internal/requirements/sync.go`, `internal/journal/current_event_catalog.go:101`, and `internal/export/archive_semantics_catalog.go:82`). |
| R-08 | M | tasks 4, 5, 8, 9, 10 `Files To Create Or Modify` | **Grant actions are added in a file no task names.** Each task lists `internal/workflow/authority.go` for its new grant action, but the closed 23-name allowlist is `sortedUniqueCapabilities` in `internal/workflow/activation.go:832`, and Envelope grants reuse it through `sortedUniqueEnvelopeActions` (`internal/workflow/envelope.go:540-546`). `authority.go` holds no action list. Under DEC-167/173, add `internal/workflow/activation.go` to each of the five file lists before dispatch. |
| R-09 | M | task 12 `Files To Create Or Modify` | **The removed Event's archive rule is out of scope.** `internal/export/archive_semantics_catalog.go:82` carries the `requirement.synced` semantic rule; task 12 removes the Event but lists no `internal/export` file, and task 15's declared search `rg -n --glob '!*_test.go' "requirement\.synced\|last_synced_hash" cmd internal` (expected exit 1) would then fail with its cause outside task 15. Add the file to task 12 (DEC-167/173). |
| R-10 | M | VAL-070..072; tasks 14/18 | **Human rows have no task reference and no fixture owner.** No task spec cites VAL-070..072 (section 6.1 coverage rule), and no task produces a runnable Human workspace. VAL-071 requires the Human to run a blocked action, record an Outcome as an activated `human`-type Actor (accepted by `internal/workflow/activation.go:817`), and rerun. Phase 22's precedent is a tagged generator (`cmd/fid/integration/p22_human_fixtures_test.go`) emitting a manifest and absolute invocation arguments. Assign an equivalent generator and instructions for VAL-070/071/072 to task 14 or 18, cite those IDs there, and keep judgment with the Human. |
| R-11 | M | VAL-002 | **Bundled assertion with asymmetric criteria.** It asserts (a) the 0013/0015 textual changes are exactly the accepted ones (file state, DEC-178) and (b) the discriminator enforces both states (behavior). The pass criterion covers only (b); the fail criterion adds (a). Split into two rows, e.g. VAL-002 (discriminator both states) and a new VAL-075 (baseline hunks exactly equal `schema.md` 1: measured seven hunks at `0013` lines 155-157, inserted before 159, 181 and `0015` lines 61, 87). |
| R-12 | M | new rows | **Clauses without an assertion** (section 2.1): (a) corrupted/unknown Phase 23 Event payload rejection with no partial effect (`schema.md` 7), owner tasks 4/5/8/9/10/12 or 16; (b) link-time `review-not-independent` (`contracts.md` 13), owner task 6; (c) `review-recorded` and `finding-resolved` Condition resolution and the one-non-null-field rule for all four Condition types (`contracts.md` 11), owner task 11. Add rows or extend VAL-003/VAL-024/VAL-048 with both states. |
| R-13 | M | `execution-scope.md` 3, 5; every task | **Test-time headroom against Go's default 10-minute per-binary timeout.** Fresh measurements: tagged integration `ok ... 560.243s` (93% of 600 s; `t.Parallel()` appears in 0 of 58 files in `cmd/fid/integration`), `cmd/fid` 471.955 s, `internal/workflow` 458.353 s under the race gate. No Makefile target, CI workflow, or tool script passes `-timeout`. Phase 23 adds real-process journeys at tasks 4, 6, 11, 14, 17, 18 to exactly these packages. A breach is a failed required command whose only non-repair cure, adding `-timeout`, is a protected command change (DEC-221 chose performance repair over a timeout increase). Record a disposition before dispatch: which task owns any in-scope speed repair, and whether the Human pre-authorizes anything further. Each completion report should record the package times so a regression is attributed to its causing task. |
| I-01 | I | VAL-001 | At phase end the latest schema is 18, where `judgment_concerns` and `judgment_concerns_by_subject` are renamed. State that VAL-001 evidence is scoped to migration 17 (migrate-to-17 introspection), so task 15 need not converge it. |
| I-02 | I | VAL-004 out-of-scope note | Cites "VAL-064's benchmark note", which does not exist. Rebuild performance is gated by the 100k benchmark's 30 s self-assert and DEC-158's three fresh phase-end samples (section 5). Fresh headroom: worst of four samples 27222 ms, three dedicated samples 26714/26723/26831 ms, so about 9-11% headroom before Phase 23 adds 10 tables and 9 Events. |
| I-03 | I | task 4 coverage | Task 4 converges its `review.recorded` fence key but omits VAL-003 from its coverage section; tasks 5, 8, 9, 10, 12 include it. |
| I-04 | I | VAL-035 | The evidence method lists five negatives; the pass criterion says "all four". The fail criterion mentions an action outside the three, which the method does not exercise; self-supersession is also unexercised. Enumerate seven negatives. |
| I-05 | I | VAL-033 | Add the reverse direction (rejected then approved -> satisfied) so both states are evidenced. |
| I-06 | I | VAL-014, 022, 038, 040, 043 | Add: foreign-Run Receipt and item/exclusion shared digest rejects (VAL-014); a `blocked` resolving Review rejects (VAL-022); missing source file (VAL-038); approved reference with non-matching authority capabilities (VAL-040); standalone Issue-inventory attribution rejects (VAL-043, PRD-AUTH-007, required by task 10). |
| I-07 | I | VAL-021 | At task 5 the review Obligation is unsatisfied anyway, so "still unsatisfied" is vacuous. Its non-vacuous proof is VAL-026's sequence at task 6; section 7 treats VAL-021's final PASS as requiring both. |
| I-08 | I | VAL-027 | Import rebuilds from the journal (`internal/export/import.go:125`), so the no-double-charge import leg is feasible once `review.recorded` replays, but the archive profile is extended only in task 16. Re-verify the import leg at task 16. |
| I-09 | I | VAL-010 | "Present" review and Human evidence exist only from tasks 6/10; the matrix is re-verified at window 3 and task 10 (contract-coverage already says "2 defines, 7 confirms"). |
| I-10 | I | VAL-069 | `v1-convergence-inventory.md` 1 says task 17 proves the cancelled families absent; the contract and task 18 say task 18. Align. |
| I-11 | I | VAL-039 / task 2 | Protected-Decision authoring (0, 16, 17 entries; duplicate class; duplicate ordinal; action outside the guarded set) is written by task 2 but asserted nowhere. Enumerate under VAL-006 or VAL-039. |
| I-12 | I | VAL-061, VAL-062 | 0018 implicitly drops `decisions_epic_kind_guard_insert` and `task_completions_issue_kind_guard_insert` (measured), and deletes non-Issue `links` and `activity_log` rows. VAL-061's inventory must include both triggers; the deletes need a seeded 17 -> 18 probe. `issues_enriched` loses no column (measured identical column list), so the "view column" clause has an empty removal set. |
| I-13 | I | task 1 note | `registration_id TEXT NOT NULL,` and `registration_version INTEGER NOT NULL,` also occur in `execution_attempts` (`0013:222-223`). The edit must stay inside `envelope_obligations`, or DEC-178 fails. |
| I-14 | I | `contracts.md` 8 | Whether a Concern or Deviation subject must resolve to an existing row is unstated; a `requirement` subject cannot resolve before task 12. Record the rule; task 8 records its interpretation. |
| I-15 | I | `public-interface.md` 2, 5 | The ten added ActivityEvent names and ReasonCode `human-outcome-absent` are unasserted; cover them in VAL-055/056 or task 18's corpus. |
| I-16 | I | environment (carry-forward candidate) | `cmd/fid/integration/main_test.go` removes its binary directory on normal paths (lines 37, 42) but not after a timeout panic or kill. 148 stale `/tmp/fid-integration-binary-*` directories (2.3G, all 2026-09-12) had filled the RAM-backed 3.9G `/tmp` to 1010M free; the Human removed them during this session (now 3.2G free). Propose a workflow/test-infra carry-forward; not Phase 23 product scope. |

Idempotency conflict evidence note (no revision requested): VAL-019 and VAL-057
name `idempotency-conflict`, an accepted Phase 21 code
(`phase-21/public-interface.md:467`). No non-test source in `cmd/` or
`internal/cli` maps `journal.IdempotencyConflictError`
(`internal/journal/idempotency.go:32`) to it, and the p21 conformance failure
fixture is generated from the declared code list
(`cmd/fid/integration/p21_conformance_test.go:513`), not reproduced. This session
could not complete a live probe (standalone probes stopped at request
validation before idempotency). Therefore VAL-019 and VAL-057 PASS require a
real-process reproduction through the built binary for each Phase 23 mutation;
fixture inventory is not evidence. If the mapping is absent, adding it for the
Phase 23 leaves is implementation of an accepted code inside tasks 4 and 13.

### 2.4 Blocker candidates

QA emits candidates only; no `ESC` identity is allocated. TL-Lite or the Tech
Lead decides whether each is resolvable within existing authority.

| Candidate | Protected category | Evidence | Authority boundary | Requested action |
|---|---|---|---|---|
| R-01 | data-model (accepted DDL, DEC-178) | section 6.4 S2 and remediation runs | `schema.md` is accepted under DEC-222; shipped SQL must equal it exactly, so no task may repair it | Tech Lead amends `schema.md` 2/3 with a link-preserving rebuild (a working variant is measured) and obtains the approval DEC-222 amendments require |
| R-02 | Event/effect contract | `current_event_catalog.go:88`; `schema.md` 6.2; `contracts.md` 10 | Effect rows are accepted contract; no task may add one | Tech Lead names the raising Event(s) and effect rows, applying DEC-195 by analogy or escalating |
| R-04 | product and public-interface semantics | `contracts.md` 3 text; rework-cycle consequence | Determines who may review after `changes-required`; not decidable from the accepted text | Tech Lead states the reading if prior accepted material determines it; otherwise Human decision |

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| 1 | Schema baseline, 0017, fence | 001-005 | high | R-01 blocks VAL-005 on linked seeds. I-01 scope VAL-001 to 17. I-13 edit locality. Runs the 100k benchmark (early warning, I-02). |
| 2 | Typed multi-Obligation Envelope | 006-010 | high | Wide-reach. R-07a/b change VAL-008/009 evidence. Baselines measured: `internal/workflow` 62.6%. |
| 3 | Review Claim mode | 011-013 | medium | VAL-013 parity needs normalized JSON diff between modes. |
| 4 | Review Records, independence | 014-019 (003) | high | R-04 decides VAL-015 fixtures. R-08 file scope. Real-process idempotency conflict required. |
| 5 | Finding dispositions | 003, 020-023 | high | R-08. I-07 vacuous half of VAL-021. |
| 6 | Review evidence, cycles, exhaustion | 024-030 | high | Wide-reach. R-03 (VAL-025 cells), R-06 (denial replay), R-12b. Baseline `internal/evidence` 23.6%. |
| 7 | Freshness, supersession | 031-033 | medium | Seeded Outcome fixtures until task 10 (no temporary writer). I-05. |
| 8 | Concerns, Deviations | 003, 034-037, 074 | medium | R-08. I-04, I-14. Three spec amendments. |
| 9 | Decision References, protected Decisions | 003, 038-042 | high | R-05 precedence, R-06 denial. R-08. Seeded Outcome fixture. |
| 10 | Human Outcomes | 003, 043-046 | high | R-02 blocks VAL-046. R-07c. R-08. Replaces task 7/9 seeded fixtures. |
| 11 | Attention via Outcomes | 047-050 | high | R-12c. Checkpoint follows; must leave reproducible state. |
| 12 | Requirement references, single model | 051-054 | high | R-07d, R-09. Removes the last fence (precedent: `33f85bb` converted the Phase 22 fence test into a live/replay parity test). Benchmark. |
| 13 | CLI surface, codes, idempotency | 055-057 | high | R-05, R-06, real-process idempotency. |
| 14 | Human projections | 058-060, 073 | high | Authors `qa/human-fixtures.md`; R-10 proposes the Human generator here or in 18. |
| 15 | 0018, V1 removal | 004, 061-063 | high | Only task with expected non-zero commands. I-12 inventory. Benchmark. |
| 16 | Portable recovery | 004, 064-065 | high | Every Phase 23 table named in equivalence evidence. I-08. Benchmark. |
| 17 | Capability matrix | 047, 066 | medium | Four surfaces; baseline matrix recorded in 6.5. |
| 18 | Conformance and wire audit | 067-069 | high | Independently pinned 30-name corpus; per-task introduction proven from git history. |

## 4. Required Evidence Strategy

Rows state what QA requires; they do not prescribe implementation. "Named
tests" means tests the completion report names and whose output it records.

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | Migrate-to-17 object inventory (type and name, excluding `sqlite_autoindex_*`) equal to `schema.md` 2: 10 tables, 18 indexes, 1 trigger, plus the two `envelope_obligations` columns and the `evidence_records` rebuild; nothing removed | VAL-001 | Developer | task 1 |
| EVID-02 | Four direct INSERT probes on the discriminator (two rejected by the CHECK, two accepted) and the one-gate-evidence index | VAL-002 | Developer | task 1 |
| EVID-03 | `git diff 0195897 -- migrations/0013_* migrations/0015_*` equal to the seven accepted hunks (after R-11 split) | VAL-002 / R-11 row | Developer | task 1 |
| EVID-04 | Per-Event fence rejection with byte-identical projection before and after; at each owner, live-apply versus replay projection equality for that Event; last owner converts the fence test to parity as `33f85bb` did | VAL-003 | Developer | tasks 1, 4, 5, 8, 9, 10, 12 |
| EVID-05 | `PRAGMA foreign_key_check` (no rows) and `integrity_check` (`ok`) at migration 17, at 18, and after a full journal rebuild | VAL-004 | Developer | tasks 1, 15, 16 |
| EVID-06 | Seeded schema-16 workspace with receipt Evidence **and** Artifact links migrated to 17: named-column row equality, links preserved, new columns null | VAL-005 | Developer | task 1 (after R-01) |
| EVID-07 | Authoring matrix: 0, 1, 16, 17 Obligations; every per-type field combination; five duplicate/subject negatives; each rejection with journal line count and projection digest unchanged | VAL-006, VAL-007 | Developer | task 2 |
| EVID-08 | Multi-Obligation evaluation in ID order with independent absence/validity-freshness/availability dimensions and the full reason set; `not-applicable` exactly for the two new types | VAL-008, VAL-010 | Developer | tasks 2, 6, 7, 10 |
| EVID-09 | Paired guarded-action probes: blocked with any applicable unsatisfied, permitted with all satisfied including gate + review + Human | VAL-009 | Developer | tasks 2, 10 (R-07b) |
| EVID-10 | Claim acquire matrix across modes; review-mode `run.recover`/`run.abandon` with adopted paths rejecting `claim-mode-invalid` with no Event versus mutating-mode Phase 22 behavior | VAL-011, VAL-012 | Developer | task 3 |
| EVID-11 | The same real-process Run lifecycle under both modes with normalized JSON results diffed, differing only in adoption | VAL-013 | Developer | task 3 |
| EVID-12 | `review.record` read back through `review.get` and raw projection; no update leaf, writer, or effect rule; repository mismatch returns expected and observed values; second Review per Run rejects; replay returns `replayed:true` with no second journal line; different digest returns `idempotency-conflict` through the built binary | VAL-014, 016, 018, 019 | Developer | task 4 |
| EVID-13 | Independence both states per the R-04 reading, naming the conflicting subject and Actor | VAL-015 | Developer | task 4 |
| EVID-14 | Verdict matrix: three invalid, three valid | VAL-017 | Developer | task 4 |
| EVID-15 | Findings only in `review.recorded` effects; address under mutating, review, and no Claim; four resolution cases plus a `blocked` resolver; supplied `state` rejected; state sequence open -> addressed -> resolved/inapplicable | VAL-020-023 | Developer | task 5 |
| EVID-16 | Link `pass`/`changes-required`/`blocked`; 3x3 type matrix per the R-03 precedence; addressed-but-unresolved blocker keeps `finding-unresolved` after a later `pass` | VAL-024-026 | Developer | task 6 |
| EVID-17 | Budget sample inspection after record, replay, linked recovery, and import; harness observation rejected; exhausting Review with charge, Stop, Attention in one Event; later `claim.acquire`/`run.start` refused with no Claim or Run; five authoring pairings | VAL-027-030 | Developer | tasks 6, 16 |
| EVID-18 | Tree change stales review Evidence with the Review Record byte-identical; Human Obligation stays fresh; supersession both directions with both records readable | VAL-031-033 | Developer | tasks 7, 10 |
| EVID-19 | Concern on every subject type; open blocker-severity Concern leaves both guarded actions unaffected; three dispositions and seven negatives; self-disposal with both Actors in JSON and Human text; Deviation supersession; Concern/Deviation link rejection | VAL-034-037, 074 | Developer | task 8 |
| EVID-20 | No body column; digest match, mismatch, missing file; blocking exactly the declared action per the R-05 precedence; six satisfaction states; two Stops then one resolved in one Event; `linked_to_decision` both states | VAL-038-042 | Developer | task 9 |
| EVID-21 | Five authority negatives plus the valid case; Obligation-type by result link matrix; Evidence triple and Human text; Stop only in the middle state per the R-02 producer | VAL-043-046 | Developer | tasks 10, 14 |
| EVID-22 | `capabilities` output from a binary built at the task-11 close commit (false) and at HEAD (true) | VAL-047 | Developer, QA Lead | tasks 11, 17 |
| EVID-23 | Wrong-class and wrong-result resolution attempts; unchanged `attention.resolve` input shape; raise facts retained; two Stops, one resolved; plus R-12c conditions | VAL-048-050 | Developer | task 11 |
| EVID-24 | Reference record/mismatch/supersede with one active row; three relations plus unknown requirement and relation; sync leaf unknown and Event absent; `links` endpoint and relation rejects with retained relations working | VAL-051-054 | Developer | tasks 12, 15 |
| EVID-25 | Real-process conformance over every new and amended leaf with unknown-field, wrong-nullability, trailing-JSON, duplicate-key negatives; one reproduction per added code against the accepted table; replay per mutation | VAL-055-057 | Developer | task 13 |
| EVID-26 | Verbatim Human fixtures in both suites recorded in `qa/human-fixtures.md`, with a JSON-field-set comparison and deterministic no-proof/no-prose checks | VAL-058, 060 | Developer; QA Reviewer check | task 14, window 13-15 |
| EVID-27 | Resume with mixed judgment state and a section past its bound returning `result-limit-exceeded` | VAL-059 | Developer | task 14 |
| EVID-28 | Latest-schema, Event-catalog, and real-process leaf probes for every inventory row including the two implicit triggers; rebuilt view and triggers with seeded data; maintenance check-name inventory | VAL-061-063 | Developer | task 15 |
| EVID-29 | Live, rebuilt, and imported projections compared per named Phase 23 table; validate mutation-free; prior projection schema rejected with `schema-unsupported` | VAL-064, 065 | Developer | task 16 |
| EVID-30 | Discovery, resume, export manifest, and shared fake/fid fixtures compared row by row with the accepted table | VAL-066 | Developer | task 17 |
| EVID-31 | 30-name allow/reject corpus pinned independently; per-name introducing commit from `git log -S` on `internal/workflow/activation.go` | VAL-067 | Developer, QA Lead | task 18, final review |
| EVID-32 | Shared case corpus run against both adapters with identical results, including negative authority, stale-version, duplicate-delivery | VAL-068 | Developer | task 18 |
| EVID-33 | Latest-schema table inventory and full real-process command tree with no cancelled family under any name | VAL-069 | Developer | task 18 |
| EVID-34 | Disposable Human workspace generator, manifest, and absolute invocation arguments | VAL-070-072 | Developer fixture, Human judgment | owner per R-10 |
| EVID-35 | Human review of every new view against the CF-040 points | VAL-073 | Human | phase end |
| EVID-36 | Full phase gate plus archive-scale and three fresh 100k samples with instance identity and observed termination | no VAL | QA Lead | final review |
| EVID-37 | Security review: Human Outcome provenance, independence computation, digest verification, no prose storage, import boundary, authority allowlist | no VAL | QA Lead | final review |

## 5. Deterministic Gates

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | no output, `exit_code: 0` | every task, final review |
| Lint | `golangci-lint run ./...` | `0 issues.`, `exit_code: 0` | every task, final review |
| Race + coverage | `go test -race -coverprofile=coverage.out ./...` | all packages `ok`, `exit_code: 0` | every task, final review |
| Coverage report | `go tool cover -func=coverage.out` | `total:` line recorded, `exit_code: 0` | every task, final review |
| Integration | `go test -tags integration -race ./cmd/fid/integration/...` | `ok  github.com/mbell-nzl/fid/cmd/fid/integration`, `exit_code: 0` | every task from task 1, final review. Hosted CI does not run it (CF-041) |
| Package coverage | `go test -race -cover ./internal/workflow/` (tasks 2, 6); `go test -race -cover ./internal/evidence/` (task 6), before and after edits | both values recorded | tasks 2, 6 |
| Reference rebuild | `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` | `ok`, parity self-assert and `< 30000 ms/rebuild` | tasks 1, 12, 15, 16 once; final review three consecutive fresh processes (DEC-158) |
| Reference instance | `go version`; `go env GOOS GOARCH`; `uname -srm`; `nproc`; `grep -m1 '^model name' /proc/cpuinfo`; `grep -m1 '^MemTotal:' /proc/meminfo` | values recorded with the samples | final review |
| Absence searches | the three `rg -n --glob '!*_test.go' ...` commands in task 15 | `exit_code: 1`, no output (declared expected non-zero) | task 15, final review |
| Benchmarks | `go test -bench . -benchmem ./...` | all packages `ok`, `exit_code: 0` | final review |
| Archive scale | `go test -tags archive_scale -race ./internal/export -run '^TestArchive(ExpandedByteBoundaries\|CumulativeBoundary\|EntryCountBoundaries)$' -count=1` | `ok`, `exit_code: 0` | final review (DEC-217 phase-end gate); also any task that changes the archive limiter |
| Build | `go build -o /tmp/fid ./cmd/fid` | `exit 0 (no output)` | final review |
| Full phase gate | the seven-command sequence in `AGENTS.md` / `08_QUALITY_PLAN.md` 7 | every command `PASS` with observed termination | final review |
| End-to-end | N/A | no held-out harness recorded | - |

No gate is narrowed, substituted, or deferred. Only task 15 declares expected
non-zero commands (R-07d adds one to task 12). Cleanup commands are never gate
results. Report lint (`project-docs/tools/report-lint.sh <task-N-done.md>
<task-N.md>`) is a Phase Runner control after the report commit, never a
Developer or QA gate. Heavy gates run serially.

## 6. Command And Feasibility Validation

Heavy mode: nothing carries forward. Working directory for every command is
`/home/mark/projects/fid` unless stated. Termination status is the command's own
exit status, captured by a serial wrapper around each command (`bash -c "<cmd>";
rc=$?`). `go clean -testcache` ran first as setup (not a gate result) so no
result is a cached replay. Light read-only probes (SQLite CLI, `fid --help`) ran
concurrently with gate steps 3-5; no step failed.

### 6.1 Toolchain and reference instance

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Go toolchain | executed | PASS | `go version go1.26.2 linux/amd64`; `go env GOOS GOARCH` -> `linux` `amd64` |
| Kernel / CPU / memory | executed | PASS | `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64`; `nproc` 6; `AMD Ryzen 7 5700G with Radeon Graphics`; `MemTotal: 8133232 kB` (Phase 22 recorded 8133236 kB) |
| Linter | executed | PASS | `golangci-lint has version 2.12.2 built with go1.26.2` (`exit_code: 0`) |
| Git, ripgrep, SQLite CLI | executed | PASS | `git version 2.53.0`; `ripgrep 15.1.0`; `sqlite3 3.50.6` |
| Runtime SQLite and migrator | `go.mod` inspection | PASS | `modernc.org/sqlite v1.51.0`; `github.com/golang-migrate/migrate/v4 v4.19.1` |
| `goimports` | `which` | PASS | `/home/mark/go/bin/goimports` |

### 6.2 Phase gate commands on the planning base

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | executed | PASS | `exit 0 (no output)` (terminated `exit_code: 0`) |
| `golangci-lint run ./...` | executed | PASS | `0 issues.` (terminated `exit_code: 0`, 48 s) |
| `go test -race -coverprofile=coverage.out ./...` | executed | PASS | last line `ok github.com/mbell-nzl/fid/internal/workflow 458.353s coverage: 62.6% of statements`; 0 `FAIL` lines (terminated `exit_code: 0`, 609 s wall). Slowest: `cmd/fid` 471.955 s, `internal/workflow` 458.353 s, `internal/export` 244.463 s |
| `go tool cover -func=coverage.out` | executed | PASS | `total: (statements) 65.5%` (terminated `exit_code: 0`); `coverage.out` is gitignored (`.gitignore:16`) |
| `go test -tags integration -race ./cmd/fid/integration/...` | executed | PASS | `ok github.com/mbell-nzl/fid/cmd/fid/integration 560.243s` (terminated `exit_code: 0`, 568 s wall). Drives R-13 |
| `go test -bench . -benchmem ./...` | executed | PASS | last line `ok github.com/mbell-nzl/fid/internal/workflow 32.886s`; 0 `FAIL` lines; in-run `BenchmarkRebuildWorkspace100kEvents-6 1 27222038134 ns/op 27222 ms/rebuild` (terminated `exit_code: 0`, 160 s) |
| `go build -o /tmp/fid ./cmd/fid` | executed | PASS | `exit 0 (no output)` (terminated `exit_code: 0`); 15995021-byte binary |

### 6.3 Task-specific commands

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| 100k rebuild (tasks 1, 12, 15, 16; final x3) | executed three times, fresh processes | PASS | `26714 ms/rebuild` / `26723 ms/rebuild` / `26831 ms/rebuild`; `ok ... internal/journal 29.336s` / `29.335s` / `29.400s` (each terminated `exit_code: 0`). The command carries no `-timeout`; each run is ~30 s against the 10 min default |
| `go test -race -cover ./internal/workflow/` | executed | PASS | `ok ... 403.037s coverage: 62.6% of statements` (terminated `exit_code: 0`) - pre-phase baseline |
| `go test -race -cover ./internal/evidence/` | executed | PASS | `ok ... 2.271s coverage: 23.6% of statements` (terminated `exit_code: 0`) - pre-phase baseline |
| `rg -n --glob '!*_test.go' "task_completions\|ac_verifications" cmd internal` | executed | PASS (syntax); matches expected before task 15 | `exit_code: 0` with matches only in `internal/issues/detail.go` (3), `transition_rules.go` (1), `internal/maintenance/errors.go` (1), `operational.go` (1), `warnings.go` (3), all in task 15's file list; no `_test.go` path returned |
| `rg -n --glob '!*_test.go' "requirement\.synced\|last_synced_hash" cmd internal` | executed | PASS (syntax); **scope gap** | `exit_code: 0`; matches `internal/requirements/sync.go` (6), `internal/journal/current_event_catalog.go:101`, `internal/export/archive_semantics_catalog.go:82`. The last file is in no task scope (R-09) |
| `rg -n --glob '!*_test.go' "has_completion_report\|concerns_addressed" cmd internal` | executed | PASS (syntax) | `exit_code: 0`; matches only `internal/workflow/rules.go` (2), in task 15's scope |
| Archive-scale gate | executed | PASS | `ok github.com/mbell-nzl/fid/internal/export 193.874s` (terminated `exit_code: 0`, 196 s) |
| `project-docs/tools/report-lint.sh` | usage output | PASS | `Usage: project-docs/tools/report-lint.sh <completion-report.md> <task-spec.md> [--repo <git-dir>]`; matches the two-argument control in `execution-scope.md` 5 |
| First Failing Test names (18) | `grep "func <name>("` over `cmd internal` | PASS | none of the 18 names exists today, so each red test is new |
| Integration build and tag | `cmd/fid/integration/main_test.go` | PASS | `//go:build integration`; `TestMain` builds `./cmd/fid` into `os.MkdirTemp("", "fid-integration-binary-")` |
| Rebuild benchmark source | `internal/journal/rebuild_bench_test.go` | PASS | `BenchmarkRebuildWorkspace100kEvents` at line 19 |

### 6.4 Accepted DDL - applied, not read

`schema.md` sections 1-5 were extracted programmatically and applied to a copy
of the migration chain in session scratch storage, never the repository. Each
file ran with `PRAGMA foreign_keys=ON` inside one `BEGIN ... COMMIT`, mirroring
the runner (`internal/db/db.go:17`; golang-migrate `executeQuery` ->
`Begin()`).

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| `schema.md` 1 edits in 0013 | scripted exact-match substitution, scoped to the `envelope_obligations` block | PASS | hunks at `0013` lines 155-157, insert before 159, line 181; `registration_*` lines also occur at 222-223 (`execution_attempts`) and were untouched (I-13) |
| `schema.md` 1 edits in 0015 | substitution, count asserted 2 | PASS | `stops.code` (line 61) and `attention_items.code` (line 87) lists gain the three codes |
| `claims.mode` after 0014 rename | schema inspection at 17 | PASS | `0014` uses `ALTER TABLE work_claims RENAME TO claims`, so the 0013 CHECK edit reaches the latest schema |
| Fresh chain 0001-0017 | applied | PASS | `foreign_key_check` 0 rows, `integrity_check` ok. Added exactly 10 tables, 18 indexes, 1 trigger as `schema.md` 2 declares; nothing removed |
| Fresh 0018 | applied | PASS | `foreign_key_check` 0 rows, `integrity_check` ok. Added `concerns_by_subject_status`; removed `ac_verifications`, `decisions`, `requirements`, `task_completions`, `judgment_concerns` (renamed), `idx_concerns_issue_status`, `judgment_concerns_by_subject`, and implicitly `decisions_epic_kind_guard_insert`, `task_completions_issue_kind_guard_insert` (I-12). `agent_runs DROP COLUMN completion_id` succeeded |
| `issues_enriched` columns 17 vs 18 | `pragma_table_info` diff | PASS | identical column lists |
| `concerns` self-FK after rename | inspection plus inserts at 18 | PASS | stored as `supersedes_concern_id TEXT REFERENCES "concerns"(id)`; valid superseding row accepted; dangling `C-404` rejected `FOREIGN KEY constraint failed`; `foreign_key_check` 0 rows |
| Discriminator and index probes at 17 | direct INSERT | PASS | gate-evidence with null registration and review with registration both rejected by the CHECK; the two valid rows accepted; second gate-evidence row rejected `UNIQUE constraint failed: envelope_obligations.envelope_id, envelope_obligations.revision` |
| 0017 down from 17 | applied | PASS | object inventory identical to 16; `evidence_records` and `envelope_obligations` `table_info` identical to 16 (the stored SQL text differs only by rename quoting) |
| 0018 down, `schema.md` 5 literal head | applied from 18 | PASS | rc=0; the remainder is specified by reference to `0002`/`0003`/`0005`/`0016` definitions, all present |
| Seeded 16 -> 17 without Artifact links | applied | PASS | named-column Evidence row identical before/after; `review_id`, `human_outcome_id` NULL |
| Seeded 16 -> 17 **with** an Artifact link (VAL-005) | applied | **FAIL - drives R-01** | `Runtime error near line 265: FOREIGN KEY constraint failed (19)`, rc=1, no Phase 23 table created. `defer_foreign_keys` variant fails at COMMIT; link-staging variant rc=0 with row, link, FK and integrity all clean |

### 6.5 CLI surface assumptions

Probe project: `fid init --prefix P23` in a fresh `git init` directory (without
`--prefix` the default prefix derived from a directory name containing `.` is
rejected as `invalid project initialization argument`; Human fixture
directories must use `--prefix` or a plain name).

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| New command groups absent before their tasks | `fid review\|finding\|deviation\|decision\|human-outcome --help` | PASS | each `exit_code: 2`, unknown command |
| V1 leaves present until removal | `fid concern --help`, `fid requirement --help`, `fid issue ac --help` | PASS | `address`, `list`, `raise`; `sync`; `add`, `list`, `verify` |
| Retained surfaces | `fid stop\|attention\|export\|maintenance\|link --help` | PASS | `stop list/raise/resolve`; `attention list/resolve`; `export create/import/validate`; `maintenance check/vacuum`; `link add/list/remove` |
| Mutation flags | `fid claim acquire --help` | PASS | `--actor`, `--idempotency-key`, `--input`; no mode yet (task 3) |
| `evidence evaluate` shape | `fid evidence evaluate --help` | PASS | "Evaluate one Obligation without mutation", `--input` only (R-07a) |
| `resume` | `fid resume --help` | PASS | `fid resume <issue-id> [--repository]` |
| Capability baseline (VAL-047/066 before state) | `fid --json capabilities` | PASS | 12 rows: `attention.v1` true `{"human_outcome_resolution":false,"state":true}`; `human.outcomes.v1` false `{"record":false}` `phase-23-unavailable`; `reviews.coverage.v1` false `{"record":false}` `phase-23-unavailable`; `work.claims.v1` false `phase-24-unavailable`; `telemetry.trajectory.v1` false `phase-25-unavailable` |
| Grant allowlist size | source inspection | PASS | exactly 23 names at `internal/workflow/activation.go:832` (R-08) |
| Existing codes | `internal/cli/codes.go` | PASS | none of the 20 new codes exists; `requirement-not-found`, `issue-not-found`, `authority-missing`, `result-limit-exceeded`, `archive-invalid`, `capability-unavailable` exist; `obligation-unsatisfied` is the retained guarded-action code (R-05) |
| `idempotency-conflict` emission | source search, fixture inspection, live probe | **UNVERIFIED** | no production mapping found; fixture is inventory-only; live probe did not reach idempotency. Real-process evidence required (section 2.3 note) |

### 6.6 Scope and gate reach

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Gate reach covers every owned path | `execution-scope.md` 3 against task file lists | PASS with R-08, R-09 | both gates compile every `cmd/`/`internal/` package; two production files an owning task must edit are unlisted |
| `internal/review` new, `internal/requirements` removed | `ls` | PASS | `internal/review` absent; `internal/requirements` present |
| Import replays the journal | source inspection | PASS | `internal/export/import.go:125` rebuilds the staged projection with `journal.RebuildWorkspace` (I-08) |
| Last-fence convergence precedent | `git show 33f85bb` | PASS | `TestRemainingRecoveryEventsRemainFenced` became `TestActivationRecoveryEventsProjectIdenticallyLiveAndReplay`; unknown-effect became unknown-Event rejection. Task 12's instruction is executable |

### 6.7 Human checklist feasibility

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Human can run a built binary on a disposable project | executed | PASS | `git init` + `fid init --prefix P23` -> `initialized fid project in .fid/` |
| Human-type Actor can be activated | source inspection | PASS | profile `actor_type` accepts `agent`, `human`, `service` (`activation.go:817`); `fid activation prepare/component/commit/finalize` present |
| Runnable fixtures for VAL-070..072 | task specs | **FAIL - drives R-10** | no task owns a generator; Phase 22 precedent `p22_human_fixtures_test.go` exists |
| Verbatim Human view fixtures for VAL-073 and window 13-15 | task 14 spec | PASS | task 14 authors `qa/human-fixtures.md` |
| Judgment not delegable | role rules | PASS | VAL-070..073 are `Human`; no `(qa-lead)` item substitutes |

### 6.8 Documentation targets

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| `03_TRD_SCHEMA.md` 3.2, 3.7, 5, 8.4 | header inspection | PASS | 3.2 line 74, 3.7 line 134, 5 line 178, 8 line 223 with 8.3 at 701 (8.4 follows) |
| `06_DESIGN_PLAN_CLI.md` 3, 5 | header inspection | PASS | line 90, line 217 |
| `02_TRD.md` 4 | header inspection | PASS | line 129 |
| Phase 21/22 read-context sections | header inspection | PASS | Phase 22 `contracts.md` 2, 4, 5, 6, 10 at lines 28, 138, 174, 235, 359; Phase 21 freshness grammar at `public-interface.md:132, 251` |

### 6.9 Storage and time budgets

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| `/tmp` filesystem | `findmnt`, `df` | CONSTRAINT | tmpfs 3.9G on an 8G host; 1010M free at assessment (148 stale integration binary dirs, I-16), 3.2G free after the Human's cleanup. Integration and archive-scale runs use it |
| Disk-backed alternative | `df -h /` | PASS | 908G available |
| Per-binary test time | measured | CONSTRAINT - drives R-13 | integration 93%, `cmd/fid` 79%, `internal/workflow` 76% of the 600 s default |

## 7. Deterministic Checklist

Executed or verified by QA Lead at final review. "Rerun named tests" means
`go test -race -count=1 -run '^(<names>)$' <packages>` using the exact names in
the completion report, with output recorded. Real-process probes use the
`/tmp/fid` built by the phase gate against a fresh `fid init --prefix` project.
PASS without evidence is invalid.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Full seven-command phase gate, each with token and observed termination `(qa-lead)` | no VAL | command transcript |
| QA-02 | Archive-scale gate plus three consecutive fresh 100k samples, each `< 30000 ms/rebuild`, with the reference instance captured in the same session `(qa-lead)` | no VAL | transcript and instance values |
| QA-03 | Migrate-to-17 object inventory equals `schema.md` 2 (10 tables, 18 indexes, 1 trigger, two columns); rerun task 1 named test `(qa-lead)` | VAL-001 | inventory diff output |
| QA-04 | `git diff 0195897 -- migrations/0013_p21_operational_kernel.up.sql migrations/0015_p22_recovery.up.sql` shows exactly the seven accepted hunks; four discriminator INSERT probes behave as accepted `(qa-lead)` | VAL-002 | diff and probe output |
| QA-05 | For each of the nine keys, `git log -S'"<event>"' --format='%h %s%n%b' -- internal/journal/project_effects.go` shows addition by the task-1 commit and removal by its owner's commit (`Refs: phase-23/task-N`); `rg -n unavailableJudgmentEvents cmd internal` returns `exit_code: 1`; rerun the owner parity tests `(qa-lead)` | VAL-003 | git output, rg status, test output |
| QA-06 | FK and integrity pragmas at migration 17, on a fresh 18 init, and after `fid maintenance` rebuild/reset per task 16 report `(qa-lead)` | VAL-004 | three pragma outputs |
| QA-07 | Seeded 16 -> 17 with receipt Evidence and an Artifact link: rows and links equal by named column, new columns null `(qa-lead)` | VAL-005 | before/after row dumps |
| QA-08 | Authoring matrix via real process: 0/1/16/17 Obligations, per-type field combinations, five duplicate/subject negatives, protected-Decision authoring (I-11); journal line count and projection digest unchanged on every reject `(qa-lead)` | VAL-006, VAL-007 | response JSON, line counts |
| QA-09 | `evidence.evaluate` without Obligation ID on a gate + review + Human Envelope in mixed states: ID order, three dimensions, full reason set, `not-applicable` for the two new types only `(qa-lead)` | VAL-008, VAL-010 | JSON output |
| QA-10 | Both guarded actions blocked with one applicable Obligation unsatisfied, permitted with all three satisfied `(qa-lead)` | VAL-009 | paired outputs |
| QA-11 | Claim acquire matrix; review-mode adoption rejects with no Event; same lifecycle under both modes diffed after normalizing IDs and timestamps `(qa-lead)` | VAL-011-013 | outputs and normalized diff |
| QA-12 | Record, `review.get`, raw projection row, no update path in the command tree or effect rules; stale repository rejects with expected/observed; second Review per Run rejects; replay `replayed:true` with journal line count unchanged; different digest `idempotency-conflict` through the binary `(qa-lead)` | VAL-014, 016, 018, 019 | outputs, line counts, catalog query |
| QA-13 | Independence both states per the R-04 disposition, with the conflicting subject and Actor in details `(qa-lead)` | VAL-015 | outputs |
| QA-14 | Six verdict cases `(qa-lead)` | VAL-017 | outputs |
| QA-15 | Finding creation only inside `review.recorded`; address under three Claim states; five resolution cases; supplied `state` rejects; state sequence; VAL-021 closed only with QA-16's addressed-blocker evidence `(qa-lead)` | VAL-020-023 | outputs, Event effect listing |
| QA-16 | Link three verdicts; 3x3 type matrix per R-03; `finding-unresolved` after an unresolving `pass` `(qa-lead)` | VAL-024-026 | outputs |
| QA-17 | `budget_samples` rows after record, replay, recovery, import; harness observation rejects; exhausting Event effects list charge, Stop, Attention together; next `claim.acquire`/`run.start` refused with no new row; five pairings `(qa-lead)` | VAL-027-030 | SQL row dumps, outputs |
| QA-18 | Tree change stales review Evidence with Review row unchanged; Human Obligation stays fresh; supersession both directions with both rows readable `(qa-lead)` | VAL-031-033 | outputs, row dumps |
| QA-19 | Concern on six subject types; blocker Concern leaves both actions unaffected; three dispositions, seven negatives; self-disposal shows two Actor fields; Deviation supersession; Concern/Deviation link rejects `(qa-lead)` | VAL-034-037, VAL-074 | outputs |
| QA-20 | `decision_refs` column list has no body; digest match/mismatch/missing file; blocking exactly the declared action; six satisfaction states; two Stops then one resolved by one Event; `linked_to_decision` both states `(qa-lead)` | VAL-038-042 | schema query, outputs |
| QA-21 | Five authority negatives and the valid case; link matrix across types and four results; Evidence triple; Stop only in the middle state `(qa-lead)` | VAL-043-046 | outputs |
| QA-22 | `git worktree add <scratch> <task-11 close commit>`, build, `fid --json capabilities` -> `human_outcome_resolution:false`; HEAD -> `true`; working path demonstrated by QA-23 `(qa-lead)` | VAL-047 | both outputs |
| QA-23 | Wrong-class and wrong-result attempts leave the Stop open; `attention.resolve` input unchanged against the Phase 22 shape; raise facts retained; two Stops, one resolved; `review-recorded` and `finding-resolved` conditions (R-12c) `(qa-lead)` | VAL-048-050 | outputs |
| QA-24 | Requirement reference record/mismatch/supersede, one active row; three relations and two unknown rejects; `fid requirement sync` unknown; `links` requirement endpoints and removed relations reject while retained relations work `(qa-lead)` | VAL-051-054 | outputs, row dumps |
| QA-25 | Rerun task 13 real-process conformance; spot-reproduce every added code against `public-interface.md` 5 (category, exit, subject, closed detail keys); replay each mutation with journal line count unchanged `(qa-lead)` | VAL-055-057 | test output, probe outputs |
| QA-26 | Diff `qa/human-fixtures.md` against live Human output of every new leaf; compare each fixture's facts with the JSON field set; search fixtures for proof claims on judgment facts and for prose beyond label, path, digest `(qa-lead)` | VAL-058, VAL-060 | diffs, search results |
| QA-27 | Resume with mixed judgment state; one section seeded past the bound returns `result-limit-exceeded` `(qa-lead)` | VAL-059 | outputs |
| QA-28 | Latest-schema query proves every inventory table, index, trigger (including both implicit kind guards), Event, and rule absent; real-process probes of removed leaves; three `rg` searches `exit_code: 1`; seeded 17 -> 18 shows non-Issue `links`/`activity_log` rows deleted; rebuilt view and trigger probes; `fid maintenance check` clean with check-name inventory `(qa-lead)` | VAL-061-063 | query output, probes, rg status |
| QA-29 | Export, validate (mutation-free by before/after digest), import into empty target; per-table live/rebuilt/imported comparison naming every Phase 23 table; prior projection schema rejects `schema-unsupported` `(qa-lead)` | VAL-064, VAL-065 | comparison output |
| QA-30 | Capability rows from discovery, resume, export manifest, and fixtures compared with `contracts.md` 15 and the 6.5 baseline `(qa-lead)` | VAL-066 | four-way diff |
| QA-31 | 30-name corpus test rerun; per-name introducing commit from `git log -S'"<name>"' -- internal/workflow/activation.go` matches its owning task `(qa-lead)` | VAL-067 | test and git output |
| QA-32 | Rerun the shared fake/fid corpus; confirm every case runs against both adapters and no case is one-sided `(qa-lead)` | VAL-068 | test output, case inventory |
| QA-33 | Latest-schema table list and full `fid --help` tree walk: no escalation, QA-plan, validation-contract, template, retrospective, or carry-forward family `(qa-lead)` | VAL-069 | outputs |
| QA-34 | Every completion report lists every exact task-spec command with token and observed termination; declared non-zero commands recorded as declared; converged rows carry file/line, introducing commit, invalidating change, executing command `(qa-lead)` | no VAL | report-to-spec comparison |
| QA-35 | `git diff 0195897 -- go.mod go.sum` empty; package list adds only `internal/review` and removes only `internal/requirements` `(qa-lead)` | no VAL | diff, `ls` |
| QA-36 | Security review: Outcome provenance, computed independence, digest verification under lock, no prose columns, import boundary, allowlist size `(qa-lead)` | no VAL | findings with file references |
| QA-37 | Six accepted spec amendments applied by their owning tasks; spec sync inputs for wrap-up `(qa-lead)` | no VAL | spec diffs |
| QA-38 | All six window reviews and the task-11 checkpoint end `PASS`, each earlier `REPAIR` closed by a fresh re-review `(qa-lead)` | no VAL | review files |
| QA-39 | Per-task package times recorded; no gate relied on hosted CI (CF-041) `(qa-lead)` | no VAL | report excerpts |

## 8. Mid-Phase Review Cadence

- QA Reviewer after windows 1-3, 4-6, 7-9, 10-12, 13-15, 16-18, using
  `project-docs/templates/mid-phase-review-template.md`. Each review reruns the
  in-window named tests, verifies in-window `VAL-NNN` both states, compares each
  report's commands with its spec, and records gate evidence or why a gate was
  not run.
- **Independent QA Lead checkpoint after task 11** (`qa/checkpoint-after-task-11.md`),
  PASS required before task 12. It reproduces, through the built binary:
  independence both states (VAL-015), resolution authority (VAL-022), cycle
  charging and exhaustion (VAL-027, VAL-029), protected-Decision blocking and
  satisfaction (VAL-039-041), Human Outcome provenance (VAL-043), Attention
  resolution through an Outcome (VAL-048-050), and the VAL-009 all-satisfied
  state (R-07b). It re-runs fresh feasibility for every non-default command
  and fixture tasks 12-18 prescribe, including the R-13 time measurements.
- **Window 13-15** additionally performs the deterministic Human-text check
  against `contracts.md` 14 using `qa/human-fixtures.md` (VAL-058, VAL-060).
- Window 2 confirms that tasks 4 and 6 introduced their real-process journeys;
  window 6 confirms task 18 only assembled cases.
- Each review records exactly one outcome: `PASS`, `REPAIR`, or `BLOCK`.
  `REPAIR` names the owning task, affected assertion, exact repair boundary,
  required evidence, and fresh re-review scope; the window does not advance.
  After repair, a fresh independent QA Reviewer writes the next unused review
  path, cites the prior review, and verifies each original finding.
- `BLOCK` is reserved for protected authority, work outside the reviewed
  task/window scope, or an exhausted no-progress budget. QA records a blocker
  candidate and never allocates an `ESC` ID.

## 9. Human Verification Scope

Fixtures come from the owner assigned under R-10 (VAL-070..072) and from task
14 (VAL-073). The Human records judgment only.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-01a | From fid output alone, which Findings are resolved, and how can you tell? `(human)` | VAL-070 | answer and quoted output |
| HUMAN-01b | Which Findings are addressed but not resolved, and what tells you they are not resolved? `(human)` | VAL-070 | answer and quoted output |
| HUMAN-01c | What still blocks completion, and why? `(human)` | VAL-070 | answer and quoted output |
| HUMAN-02a | Before and after recording the Human Outcome, was the blocked action refused and then permitted? `(human)` | VAL-071 | observation |
| HUMAN-02b | State, from fid output alone, why the action became permitted. `(human)` | VAL-071 | answer |
| HUMAN-03a | From the rework-exhaustion denial, what caused the block? `(human)` | VAL-072 | answer |
| HUMAN-03b | What exact action would resolve it, and does the output avoid suggesting you edit internals? `(human)` | VAL-072 | answer |
| HUMAN-04a | Does every new view present labeled sections? `(human)` | VAL-073 | accept/reject per view |
| HUMAN-04b | Does every new view state its operation purpose? `(human)` | VAL-073 | accept/reject per view |
| HUMAN-04c | Is every new view free of routine-state narration? `(human)` | VAL-073 | accept/reject per view |

## 10. Phase-End Decision Rule

The phase can proceed to final review only when:

- Every required deterministic gate in section 5 passes with observed
  termination evidence, including the three fresh 100k samples and the
  archive-scale gate.
- Every `VAL-NNN` (all 74, plus any rows added under R-11/R-12) has planned
  evidence and a final-review outcome slot; every outcome is `PASS`, `FAIL`,
  `BLOCKED`, or `CARRY-FORWARD` with cited evidence.
- Every `(qa-lead)` item has reproducible evidence captured in the final-review
  session.
- Human checks are judgment-only and mapped to `human-checklist.md`.
- Every mid-phase review and the task-11 checkpoint is `PASS`; each earlier
  `REPAIR` has a fresh independent `PASS` re-review.
- No unresolved Must-fix finding or blocker candidate remains.
- The six accepted spec amendments are applied and named for wrap-up sync.

An in-scope failed gate or incomplete evidence triggers repair; recording
`BLOCKED` does not make it acceptable. A protected missing or contradictory
contract prevents final review and is emitted as a blocker candidate for
TL-Lite or the Tech Lead, which alone may allocate an `ESC` identity.

## 11. Readiness Conditions

This plan becomes Ready only after:

1. The Tech Lead dispositions R-01..R-13 and records whether R-01, R-02, and
   R-04 need Human decision; Improvement items are applied or declined.
2. Accepted edits are applied by the Tech Lead to the contract, the three slice
   documents, and the affected task specs (at minimum tasks 1, 2, 4, 5, 6, 8,
   9, 10, 12, 14 or 18).
3. A fresh QA Lead readiness refresh revalidates every changed command, DDL
   row, fixture owner, and assertion in a new section of this file, reusing no
   result above for anything that changed.
4. The Phase Runner's readiness check passes.
