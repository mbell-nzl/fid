# Phase 23 QA Plan - Review And Human Judgment

> Phase: 23 - Review and Human judgment
> First pass: 2026-09-15. Readiness refresh: 2026-09-19 (section 12).
> Readiness confirmation of DEC-229/DEC-230: 2026-09-19 (section 13).
> Readiness confirmation of DEC-231: 2026-09-19 (section 14).
> Readiness confirmation of DEC-232: 2026-09-19 (section 15).
> QA Lead: Claude Opus 5 (activated directly by the Human; no Tech Lead handoff)
> Operating mode: Heavy. Sections 1-11 were produced fresh on `phase-23` at
> `0c28936a6ad3980d2b3c4439b42ec70c83313d4c`; section 12 was produced fresh at
> `6057066` (planning base `0195897a347dfe1872db0f7be3800c5c3ccdefcf` plus
> documentation commits only). Nothing carries forward from Phase 22, and
> section 12 reuses no earlier result for anything the Tech Lead changed.
>
> **Status: NOT READY - one Must-fix revision remains (section 15.6).** This is
> not a placeholder `BLOCKED`: every required command is concrete and was
> executed or feasibility-proven in sections 6 and 12. Every disposition so far
> is confirmed: R-01..R-13 and I-01..I-16 under DEC-223..228 (section 12), R-14,
> R-15 and DEC-230 under DEC-229 (section 13), R-16 under DEC-231 (section 14),
> and R-17, R-17b, R-17c, R-18 and R-18b under DEC-232 (section 15). The plan
> cannot become Ready because the journal holds one further closed registry with
> no owning task: **R-19** - `current_event_value_relations.go` carries a second,
> byte-identical copy of the 23-name capability allowlist that R-08 assigns in
> `internal/workflow/activation.go`, and it rejects all seven Phase 23 grant
> actions, blocking the first Activation that stages a Phase 23-capable Actor
> profile. It is bounded and decision-free, with no blocker candidate. Read
> section 12 before sections 2-11: it supersedes the earlier statements listed
> in 12.8. Then section 13 (13.4 supersedes statements in section 12), then
> sections 14 and 15.

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

## 12. Readiness Refresh - 2026-09-19

Readiness condition 3. A fresh QA Lead session revalidated every changed
command, DDL row, fixture owner, and assertion. **No result in sections 1-11 is
reused for anything that changed.** Conditions 1 and 2 are satisfied: the Tech
Lead dispositioned R-01..R-13 and I-01..I-16 (all accepted), the Human approved
DEC-223..227 on 2026-09-19, DEC-228 covers the remaining dispositions, and the
edits are applied across `contracts.md`, `public-interface.md`, `schema.md`,
`execution-scope.md`, `contract-coverage.md`, `v1-convergence-inventory.md`,
`qa/validation-contract.md`, `GLOSSARY.md`, and all eighteen task specs.

Working directory for every command: `/home/mark/projects/fid`. The gate
sequence ran serially through one wrapper that recorded each command's own exit
status and wall time; `go clean -testcache` ran first as setup, so no result is
a cached replay. Scratch state was allocated fresh (`mktemp -d`) under the
session scratch directory; the repository working tree stayed clean throughout
and no product file was written.

### 12.1 Toolchain And Reference Instance (fresh)

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Go toolchain | executed | PASS | `go version go1.26.2 linux/amd64`; `go env GOOS GOARCH` -> `linux` `amd64` |
| Kernel / CPU / memory | executed | PASS | `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64`; `nproc` 6; `AMD Ryzen 7 5700G with Radeon Graphics`; `MemTotal: 8133236 kB` |
| Linter | executed | PASS | `golangci-lint has version 2.12.2 built with go1.26.2` |
| Git, ripgrep, SQLite CLI | executed | PASS | `git version 2.53.0`; `ripgrep 15.1.0`; `sqlite3 3.50.6` |
| Runtime SQLite and migrator | `go.mod` inspection | PASS | `modernc.org/sqlite v1.51.0`; `github.com/golang-migrate/migrate/v4 v4.19.1` |
| Phase HEAD under test | `git rev-parse HEAD` | PASS | `6057066`; `git status --porcelain` empty |

### 12.2 DEC-223 - Amended 0017 On The Real Chain

The Tech Lead's own validation note is explicit that its sandbox proof used the
`sqlite3` 3.50.6 CLI on a four-table minimal schema, and that the shipped
migration under `modernc.org/sqlite` was still unproven. That gap is now closed.

**Method.** A real schema-16 workspace was built end to end **through the
shipped `fid` binary only** - `git init`, `fid init --prefix P23`,
`issue create`, `activation prepare/component/resource/commit/finalize`,
`gate register`, `envelope create/activate`, `claim acquire`, `run start`,
`gate run` (artifact path `base.txt`), `evidence link` with
`required_artifact_ids` - so the receipt Evidence and its `evidence_artifact_links`
row are ordinary durable projection state, not hand-written SQL. The accepted
`schema.md` sections 2 and 3 were extracted programmatically into
`0017_p23_review_judgment.{up,down}.sql` in a copy of the real migration chain,
and applied **in place** to that workspace's `.fid/fid.sqlite` through
`internal/db.MigrateUp`, which is exactly what `catchUpWorkspaceLocked`
(`internal/journal/projection.go:182-185`) calls. A second copy of the chain
carried the pre-amendment (DEC-222) text with the six staging lines removed.

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Seed reachability through the product | 16 real CLI calls against the built binary | PASS | `EVIDENCE=EVD-672103 STATUS=satisfied`, `RECEIPT=RCP-2f0abd`, `ARTIFACTS=["ART-fd0274"]`; projection at `schema_migrations` `16\|0` with one `evidence_records` row, one `evidence_artifact_links` row (`EVD-672103\|ART-fd0274\|1`), `foreign_key_check` empty, `integrity_check` `ok` |
| Pre-amendment 0017 up, real runner, linked seed | executed | **FAIL as predicted** | `MIGRATE-ERROR: run up migrations: constraint failed: FOREIGN KEY constraint failed (787) in line 0: CREATE TABLE review_records ...`, `rc=1`. R-01's premise reproduces under `modernc.org/sqlite`, not only under the CLI |
| Pre-amendment 0017 down, real runner, linked state | executed (`Steps(-1)`) | **FAIL as predicted** | `STEP-ERROR: constraint failed: FOREIGN KEY constraint failed (787) in line 0: DROP INDEX evidence_records_one_outcome; ...`, `rc=1`, and the projection is left `16\|1` - **dirty**. The down direction needed the amendment too |
| **Accepted (amended) 0017 up, in place** | executed | **PASS** | `MIGRATE-OK`, `rc=0`; `schema_migrations` `17\|0` |
| Evidence row preserved by named column | before/after dumps compared | PASS | `EVD-672103` identical in `id, obligation_id, evidence_type, receipt_id, producer, trust_source, target_repository_json, validity, created_at`; `review_id IS NULL` and `human_outcome_id IS NULL` both `1` |
| Artifact links preserved exactly | before/after dumps diffed | PASS | `LINKS-IDENTICAL` - `EVD-672103\|ART-fd0274\|1` with its exact `required_material` |
| `evidence_artifact_links` never rebuilt | `sqlite_schema.sql` diffed | PASS | `DDL-BYTE-IDENTICAL` |
| No staging table leaked | schema query | PASS | `SELECT count(*) FROM sqlite_schema WHERE name LIKE '%_carry'` -> `0` |
| Integrity after up | pragmas | PASS | `foreign_key_check` empty; `integrity_check` `ok` |
| Object inventory 16 -> 17 (VAL-001, EVID-01) | `sqlite_schema` diff excluding `sqlite_autoindex_*` | PASS | exactly **10 tables, 18 indexes, 1 trigger** added, **nothing removed**; `envelope_obligations.review_profile` and `.outcome_class` present and null |
| **Accepted 0017 down, single step, with a `review` Evidence row present** | executed | **PASS** | `STEP-OK`, `rc=0`; `schema_migrations` `16\|0` |
| Receipt-only filter in the down direction | row dumps | PASS | `EVD-qa-review` and its link are gone; `EVD-672103` and `EVD-672103\|ART-fd0274\|1` survive - exactly the explicit loss `schema.md` 3 declares |
| Down restores schema 16 | object inventory diff | PASS | `INVENTORY-IDENTICAL-TO-16`; `foreign_key_check` empty; `integrity_check` `ok`; no `%_carry` leak |
| Shipped full-teardown API | `internal/db.MigrateDown` on the linked workspace | PASS | `MIGRATE-OK`, `rc=0` - the amended down text is also executable inside the whole-chain teardown the product actually exposes |

**Verdict.** DEC-223 is correct as shipped, in both directions, on the real
chain, under the real driver, in place, with links present. VAL-005 is feasible
exactly as written and must not be narrowed. Two observations for the owning
task, carried as improvements below: the product exposes no single-step down
(`MigrateDown` runs `m.Down()` to zero), so task 1's down evidence is either the
full teardown or a test-local `Steps(-1)`; and a failed migration leaves
`schema_migrations` **dirty**, which is worth recording as the failure mode
VAL-005 rules out.

### 12.3 DEC-223 Follow-on - 0018 Re-verified On The Same Real Workspace

VAL-061 and VAL-062 were both restated, so their subjects were re-measured on
the migrated real workspace rather than on the first pass's synthetic chain.

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| The two implicitly dropped kind guards exist at 17 | schema query | PASS | `task_completions_issue_kind_guard_insert`, `decisions_epic_kind_guard_insert` |
| `requirements.last_synced_hash` exists at 17 | `pragma_table_info` | PASS | count `1` |
| 0018 up on the real 17 workspace | executed (`Steps(1)`) | PASS | `STEP-OK`; `schema_migrations` `18\|0` |
| Non-Issue `links` and `activity_log` deletes (VAL-062) | seeded `issue->requirement` link and a `requirement` activity row, then migrated | PASS | after 0018 only `issue->issue relates_to` and `activity:issue` remain; every Issue row retained |
| Both kind guards dropped implicitly (VAL-061) | schema query | PASS | count `0` |
| `requirements` table and its column gone (VAL-061) | schema query | PASS | count `0` |
| `issues_enriched` column list across 0018 (VAL-061) | `pragma_table_info` diff | PASS | `VIEW-COLUMNS-IDENTICAL (17 columns)` - the accepted view-column removal set is empty and the measured list is 17 columns |
| Rename and new index | schema query | PASS | `concerns` present, `judgment_concerns` absent, `concerns_by_subject_status` present |
| `agent_runs.completion_id` dropped | `pragma_table_info` | PASS | count `0` |
| Pragmas at 18 (VAL-004) | pragmas | PASS | `foreign_key_check` empty; `integrity_check` `ok` |

### 12.4 DEC-227 - New Command Text, Fresh Feasibility And Budgets

Every task spec now carries `-timeout 30m` on the two long gates, and tasks 2
and 6 on their named package runs. That is new command text, so every result
here is new. All eighteen task specs were checked mechanically against
the byte-exact bounded command lines: 18/18 carry both, task 2 carries the
bounded `internal/workflow` run, task 6 carries both bounded package runs, and
no task retains an unbounded form of either long gate (zero unbounded
occurrences across the pack).

| Referenced step (byte-identical to the task specs) | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | executed | PASS | `exit 0 (no output)` (terminated `exit_code: 0`, 0 s) |
| `golangci-lint run ./...` | executed | PASS | `0 issues.` (terminated `exit_code: 0`, 2 s) |
| `go test -race -timeout 30m -coverprofile=coverage.out ./...` | executed | PASS | last line `ok github.com/mbell-nzl/fid/internal/workflow 443.528s coverage: 62.6% of statements`; 0 `FAIL` lines (terminated `exit_code: 0`, 562 s wall) |
| `go tool cover -func=coverage.out` | executed | PASS | `total: (statements) 65.5%` (terminated `exit_code: 0`) |
| `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` | executed | PASS | `ok github.com/mbell-nzl/fid/cmd/fid/integration 551.480s` (terminated `exit_code: 0`, 554 s wall) |
| `go test -race -timeout 30m -cover ./internal/evidence/` (task 6) | executed | PASS | `ok github.com/mbell-nzl/fid/internal/evidence 2.273s coverage: 23.6% of statements` (terminated `exit_code: 0`) |
| `go test -race -timeout 30m -cover ./internal/workflow/` (tasks 2, 6) | executed | PASS | `ok github.com/mbell-nzl/fid/internal/workflow 405.408s coverage: 62.6% of statements` (terminated `exit_code: 0`, 409 s) |
| `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` (unchanged text, re-run) | executed | PASS | `26580 ms/rebuild`; `ok github.com/mbell-nzl/fid/internal/journal 29.108s` (terminated `exit_code: 0`) |
| `rg -n --glob '!*_test.go' "requirement\.synced\|last_synced_hash" cmd internal` (task 12, new declaration) | executed | PASS (syntax); matches expected before task 12 | `exit_code: 0`; matches exactly the three files task 12's `## Expected Non-Zero Commands` names - `internal/export/archive_semantics_catalog.go:82`, `internal/requirements/sync.go` (6), `internal/journal/current_event_catalog.go:101`; all three are in task 12's own scope |
| Task 15's two declared absence searches (unchanged text, re-run) | executed | PASS (syntax); matches expected before task 15 | `exit_code: 0`; matches only `internal/maintenance/{operational,errors,warnings}.go`, `internal/issues/{detail,transition_rules}.go`, `internal/workflow/rules.go` - all in task 15's scope |

**VAL-079 budgets, measured under the new command text.**

| Budget | Limit | Measured | Margin |
|---|---|---|---|
| Tagged integration gate | < 900 s | **551.480 s** | 39% |
| `cmd/fid` under the race gate | < 600 s | **473.689 s** | 21% |
| `internal/workflow` under the race gate | < 600 s | **443.528 s** | 26% |

DEC-227's baselines are confirmed: 551.5 s against the recorded 560.2 s, 473.7 s
against 472.0 s, and 443.5 s against 458.4 s - all within run-to-run variance on
the same instance, and all inside budget. The explicit `-timeout 30m` is
effective and no gate approached it. The headroom the decision buys is real:
without the bound, `cmd/fid` at 473.7 s sits 21% under an *undeclared* 600 s
cliff that a killed command would turn into `BLOCKED`.

**One measured fact contradicts the accepted cause of CF-042 - see R-15.** A
single clean tagged run (exit 0, no kill, no timeout, 551 s) left **20**
`fid-integration-binary-*` directories and **306 MB** in a 3.9 GB RAM-backed
`/tmp`. `/tmp` held 3.6 MB before the run. The directories are created between
09:44:07 and 09:52:04 - during the run - each holding a complete 16 MB `fid`
binary. `TestMain` creates exactly one directory
(`cmd/fid/integration/main_test.go:22`) and removes it on a normal exit
(`:42`), so these are not its leak: five tests re-execute the test binary as a
child process (`exec.Command(os.Args[0], ...)` in `p22_activation_fault_test.go:134`,
`p22_budget_fault_test.go:208`, `p22_force_budget_fault_test.go:34`,
`p22_recovery_fault_test.go:52`, and `os.Executable()` in
`p22_human_fixtures_test.go:64`), each child runs `runIntegrationSuite` again,
rebuilds the binary, and - for the crash-process children, whose whole purpose
is to die abnormally - never reaches its `RemoveAll`.

### 12.5 DEC-224 - All Four Raising Sequences Reachable At Task 10

| Sequence | Reachability at task 10 | Result | Evidence |
|---|---|---|---|
| `envelope.activated` on a Human-only Envelope | `contracts.md` 2 accepts 1..16 Obligations with no `gate-evidence` requirement; the `human-outcome` Obligation type and its authoring land at task 2; `stops.code` and `attention_items.code` gain `human-outcome-required` at task 1 (`schema.md` 1) | PASS | `envelope_obligations_one_gate_evidence` is a partial unique index, not a presence requirement - verified on the applied schema at 17; a revision whose only Obligation is `human-outcome` is authorable and activates unsatisfied |
| `evidence.linked` satisfying the last non-Human Obligation | tasks 2 and 6 precede | PASS | reproduced end to end in 12.2: `evidence link` returned `STATUS=satisfied` for the `gate-evidence` Obligation through the built binary |
| `review.recorded` resolving the last blocking Finding | tasks 4, 5, 6 precede; `schema.md` 6.2 already gives `review.recorded` optional `stops/insert`, `attention_items/insert`, `attention_stop_links/insert` | PASS | the new-rule row is present in `schema.md` 6.2 as the decision states; no additional effect-row amendment is needed for this producer |
| **`artifact.availability-observed` flipping the last required Artifact** | the one sequence with no prior evidence anywhere | **PASS - reproduced live** | On the real seeded workspace: `artifact verify` -> `available`, `changed:false`, `EVT-...000021`; retained blob removed -> `missing`, `changed:true`, `EVT-...000022`, and `evidence evaluate` returns `status:"unsatisfied"`, `artifact_availability:"missing"`, reason `artifact-unavailable`; blob restored -> `available`, `changed:true`, `EVT-...000023`, and evaluation returns `status:"satisfied"` with `reasons:[]`. The Obligation genuinely transitions unsatisfied -> satisfied inside that one Event |

The fourth sequence - the one the Human specifically asked about - is real, not
theoretical. `internal/evidence/artifact.go:428-444` derives availability by
digest-verifying the retained content-addressed blob and writes a
`referenced_artifacts/update` effect with a precondition on the prior value, so
`missing -> available` is exactly one qualifying Event. The fixture technique it
requires (remove and restore the blob under `.fid/artifacts/<digest>`) is
recorded nowhere; see I-24.

Effect-row ownership for the amendment is assigned and consistent:
`execution-scope.md` 6 and task 10's file list both give task 10 the optional
`stops/insert`, `attention_items/insert`, `attention_stop_links/insert` rows on
`evidence.linked` and `artifact.availability-observed` in
`internal/journal/current_event_catalog.go`, and task 9 already owns the
`envelope.activated` rows. Two mechanical gaps in that assignment are R-14 and
I-17 below.

### 12.6 DEC-226 - Failure Codes, Producers, Precedence, Vocabulary

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Every code in `public-interface.md` 5 has a producer in 5.1 | list comparison | PASS | all 14 table rows are covered; the grouped `*-not-found` row maps to "the read and mutation leaves naming that entity" (category, not operation - see I-25) |
| `finding-unresolved` is no longer a failure code | source search | PASS | absent from `internal/cli/codes.go`; the section 5 table no longer lists it; `contracts.md` 2 places it in the evaluation reason precedence instead |
| Count of genuinely new failure codes | comparison against `internal/cli/codes.go` | PASS | 20 codes declared in section 5; exactly one (`requirement-not-found`) already exists, so **19** are new. This supersedes section 6.5's "20 new codes" |
| Codes VAL-056 must reproduce, per producing task | per-code search across all eighteen task specs | PASS | all 19 new codes have an owning task that names `internal/cli/codes.go`: task 3 (`claim-mode-invalid`, named in its body at `task-3.md:21` rather than in the file bullet), task 4 (4), task 5 (2), task 6 (4, including `review-not-independent` again for the VAL-077 link-time check), task 8 (3), task 9 (3), task 10 (3); task 12 retains `requirement-not-found` and task 13 owns the adapters |
| `evidence.link` single-code precedence | `contracts.md` 13 against DEC-226 part 1 | PASS | identical order in both: `outcome-not-permitted`, `review-changes-required`, `review-not-independent`, `evidence-type-not-accepted`; the nine-cell matrix resolves to 3 accepted, 2 `outcome-not-permitted`, 4 `evidence-type-not-accepted`, which is exactly VAL-025's restated pass criterion |
| Guarded-action denial precedence | `public-interface.md` 5.1 against VAL-039, VAL-056 | PASS with one wording note | the five-step order is determinate and both VAL-056 precedence probes are constructible at task 13 (tasks 9 and 10 precede). Step 1 says "every retained Phase 21/22 precondition and authority failure" while step 5 names `obligation-unsatisfied`, itself a retained failure - determinate only because step 5 names it explicitly (I-19) |
| Denials commit nothing | `schema.md` 6.3, `public-interface.md` 5.2, `contracts.md` 6 | PASS | mutually consistent: no denial Event, no `idempotency_results` row, null `EventRef`, `replayed:false`, immutable Stop/Attention IDs; Phase 22 `run.execution-blocked` explicitly retained. VAL-057's restated criterion matches all four clauses |
| `idempotency-conflict` real-process emission | source search, live probe through the built binary | **UNVERIFIED (unchanged)** | replay under the original key works (`ok:true`, `replayed:true`, same `EVD-672103` and `EVT-...000020`, no second Event); the same key with a different payload returns the adapter's generic `invalid-argument` / "invalid Evidence request" (`cmd/fid/evidence.go:125`). `"idempotency-conflict"` exists in `internal/journal/idempotency.go` as `IdempotencyConflictError` but appears as a CLI code only in fake-adapter test files. It **is** an accepted Phase 21 code (`phase-21/public-interface.md:467`, `06_DESIGN_PLAN_CLI.md:230`), so registering it is conformance, not a new code - see I-21 |
| Evaluation reason precedence insertion | `internal/workflow/evaluation.go:380-387` | PASS with an ambiguity | `precedence` is an integer map where `obligation-evidence-absent` is 8 and seven other codes also hold 8; `receipt-unsuccessful` is 9. "Immediately after `obligation-evidence-absent` ... Phase 21/22 ranks are unchanged" is satisfiable only as *relative* order (I-18). No asserted observable depends on the finer reading |
| `GLOSSARY.md` split | read in full | PASS | `Reason code` is replaced by **Failure code** and **Evaluation reason**, each with definition, use-when, avoid, the both-sets rule with `protected-decision-required` as the worked example, and an explicit retirement note naming DEC-226. `Last updated` is 2026-09-19 |
| The task 13 rename targets | source verification | PASS with one false example | `reasonCodeDefinitions` is declared at `internal/cli/codes.go:200` with exactly three other uses - `codes.go:328`, `codes.go:339`, `envelope.go:38` - matching the spec's "declaration and its three other uses". `reasonCodeExpectation` (`internal/cli/codes_test.go:30`) does name the failure set. `reasonCodes` in `cmd/fid/evidence_test.go:429` does **not**: it is the parameter of `assertEvaluation` compared against `evaluation reasons` (I-20) |

**Term usage in this plan.** Checked against the new GLOSSARY entries. Sections
4, 5, 7 and this section use *failure code* only for the closed refusal set
(`review-not-independent`, `outcome-not-permitted`, `protected-decision-required`,
`human-outcome-required`, `idempotency-conflict`, ...) and *evaluation reason*
only for the state set (`finding-unresolved`, `human-outcome-absent`,
`outcome-superseded`, `artifact-unavailable`, `review-changes-required` when it
appears in a `reasons` array). Two earlier statements needed correction and are
superseded in 12.8: section 6.5's "20 new codes", and EVID-16/QA-16's phrasing
that an Obligation "keeps `finding-unresolved`" - correct as an evaluation
reason, and now explicitly labelled as one.

### 12.7 Contract, File Scopes, And Fixture Owners

**Mechanical contract checks, all re-run on the amended file.**

| Check | Result | Evidence |
|---|---|---|
| Row count | PASS | 79 `VAL-NNN` rows; header claims 79 (75 deterministic + 4 Human) |
| Human rows | PASS | exactly 4 (`VAL-070`, `071`, `072`, `073`) |
| Status | PASS | 0 rows not `Ready` |
| One validator per row | PASS | **40** QA Lead + **35** QA Reviewer + **4** Human = 79; no row names two (the first pass counted 36/34/4 = 74) |
| Every `VAL-NNN` cited by a task spec | PASS | 0 uncited; VAL-070..072 are now cited by tasks 14 and 18, closing the first pass's `09_AI_WORKFLOW.md` 6.1 violation |
| No task cites an undefined `VAL-NNN` | PASS | the cited set and the defined set are identical, 79 members each |

**New rows.**

| Row | Evidence method feasible? | Basis |
|---|---|---|
| VAL-075 | PASS | `git cat-file -t 0195897a347...` is a commit and `git diff --stat 0195897..HEAD -- migrations/` is empty, so the baseline diff shows exactly task 1's edits. All seven hunk anchors verified in place: `0013:155` `obligation_type` CHECK, `0013:156-157` registration columns, insertion before `0013:159` `UNIQUE (envelope_id, revision, id)`, `0013:181` `work_claims.mode`, `0015:58` `stops.code`, `0015:84` `attention_items.code`. I-13's collision is real: `registration_id TEXT NOT NULL,` and `registration_version INTEGER NOT NULL,` recur byte-identically at `0013:222-223` in `execution_attempts` |
| VAL-076 | PASS | nine Phase 23 Event types exist in `schema.md` 6.2's new-rules table; all nine are implemented by task 12, so task 16 is the first point where the per-family corruption proof is executable, exactly as the row's note says |
| VAL-077 | PASS | `evidence.link` already loads its binding and authority before committing (`internal/evidence/evidence.go:107-121`), so a link-time independence re-check is an ordinary added precondition at task 6; both states are named |
| VAL-078 | PASS with a scope gap | the four Condition types and the "exactly one non-null field" rule are enforceable, but the closed registries that must change are unowned - see R-14 |
| VAL-079 | PASS | measured in 12.4; both budgets hold with margin, and "every completion report carries its three recorded package times" is a report-lint-checkable fact |

**Restated rows.** All 74 pre-existing rows moved `Draft` -> `Ready`. Comparing
every row's ten columns against `f7811c6` with the `Status` column excluded,
**33** are substantively restated, not the 27 the refresh brief names: `VAL-001`,
`VAL-004`, `VAL-006`, `VAL-010`, `VAL-014` and `VAL-024` also changed, under
I-01, I-02, I-11, I-09 and I-06 respectively, and were re-read here too. All 33
were re-read against their sources, and every restatement is a strengthening or
a note clarification, never a weakening: VAL-002 sheds the file
diff to VAL-075 and gains the two valid counterparts; VAL-005 gains the link
requirement and the explicit "a link-free seed alone is not evidence"; VAL-008,
009, 010 and 027 gain their second point of verification; VAL-015 goes from two
recordings to four and names the mutating-side-only reading; VAL-021 names the
`claim-mode-invalid` code and `required_mode`; VAL-022 goes from four cases to
six; VAL-025 goes from one code to the full precedence; VAL-033 adds the reverse
supersession direction; VAL-035 goes from four negatives to seven; VAL-038, 040,
043 gain their negatives; VAL-039 gains the with-an-Obligation-also-unsatisfied
probe; VAL-045 fences off task 14's text; VAL-046 goes from one producer to four
plus the at-most-once rule and both non-qualifying states; VAL-053 moves
`last_synced_hash` to VAL-061 and names both Event catalogs; VAL-055 adds the
activity-log inventory; VAL-056 adds the two precedence probes; VAL-057 adds the
denial half; VAL-061 adds the implicit triggers and the column; VAL-062 adds the
0018 delete probe; VAL-069 moves to task 18; VAL-070..072 gain a fixture owner;
VAL-001 is scoped to schema 17 with the 0018 rename excluded; VAL-006 gains the
four protected-Decision authoring negatives; VAL-010 gains its task-10
present-evidence leg; VAL-014 gains the foreign-Receipt and shared-digest
negatives; and VAL-004 and VAL-024 gain corrected cross-references only.
No row's pass criterion was loosened and no negative case was dropped.

**Amended file scopes and fixture owners.**

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| R-08: `internal/workflow/activation.go` added to tasks 4, 5, 8, 9, 10 | task-spec inspection plus source check | PASS | all five name it with the correct rationale; `sortedUniqueCapabilities` is at `internal/workflow/activation.go:832` and holds **exactly 23** names; `sortedUniqueEnvelopeActions` reuses it (`envelope.go:505,546`); `authority.go` holds no action list |
| Grant-action arithmetic (VAL-067) | count | PASS | 23 + 7 added (`review.record`, `finding.address`, `concern.record`, `concern.dispose`, `deviation.record`, `decision.record-reference`, `human-outcome.record`) = **30**, matching `public-interface.md` 4.1 and task 18; `requirement.record-reference` and `requirement.link` are correctly excluded |
| R-09: `internal/export/archive_semantics_catalog.go` added to task 12 | task-spec inspection plus search | PASS | named with the rationale, and the declared absence search does match `archive_semantics_catalog.go:82` today, so without the file that command could not pass |
| R-07d: task 12 gains `## Expected Non-Zero Commands` | section inspection plus execution | PASS | the section exists, declares `expected exit_code: 1`, names the exact recording form, lists the three current matches, and the command itself is in task 12's `Exact Commands To Run` block and its conditional-command list |
| R-10: task 14 owns the VAL-070..072 fixture generator | task-spec inspection plus precedent check | PASS | task 14 names `cmd/fid/integration/p23_human_fixtures_test.go` and all three seeded scenarios; the `p22_human_fixtures_test.go` precedent already emits exactly the prescribed manifest (`SourceCommit`, `BinarySHA256`, `Checkpoints[].Inspections[].{CWD,Args}`) and already supports a durable output root via an environment variable, so a Human-runnable fixture is producible. Task 18 names the regeneration against the final binary and cites all four Human rows |
| I-16 / CF-042: task 1 prunes stale integration binary directories | task-spec inspection plus measurement | **FAIL - drives R-15** | the instruction is present and bounded, but its stated cause and its 24-hour age bound do not match the measured leak (12.4) |
| I-10: VAL-069 moves to task 18 | cross-document check | PASS | `v1-convergence-inventory.md` 1, task 18, and the contract now all say task 18 |
| DEC-228 subject-resolution rule | `contracts.md` 8 against VAL-034 | PASS, implicitly | `contracts.md` 8 states it and task 8 records it; VAL-034 proves it only as a side effect of recording a `requirement`-subject Concern at task 8, where no requirement row can exist (I-23) |

### 12.8 Statements In Sections 1-11 Superseded By This Section

| Where | Superseded statement | Replacement |
|---|---|---|
| Header, section 1, section 10 | "74 assertions (70 deterministic, 4 Human)"; "all 74" | **79 assertions, 75 deterministic and 4 Human** |
| Section 2 and 2.3 | R-01..R-13 open; three blocker candidates | all thirteen accepted and applied under DEC-223..228; **no blocker candidate is open**; the open items are R-14 and R-15 in 12.10 |
| Section 2.1 gap table | the eight mapped-but-unasserted clauses; the duplicated `links` row | all eight now carry an evidencing row - `schema.md` 7 payload rejection -> **VAL-076**; `contracts.md` 13 link-time independence -> **VAL-077**; `contracts.md` 11 Condition types -> **VAL-078**; `contracts.md` 10 Stop producers -> **VAL-046** under DEC-224; `contracts.md` 9 protected-Decision authoring -> **VAL-006**; `contracts.md` 4 and `public-interface.md` 3 negatives -> **VAL-014**; `schema.md` 4 deletes and implicit trigger drops -> **VAL-061**, **VAL-062**; `public-interface.md` 2 ActivityEvent and ReasonCode additions -> **VAL-055** and the VAL-008/026/033/046 set - and the duplicate row is gone |
| Section 3 | task coverage columns | task 1 also covers VAL-075; task 6 VAL-077; task 11 VAL-078; task 16 VAL-076; task 18 VAL-079 and the VAL-070..073 regeneration; task 14 owns the VAL-070..072 fixtures; tasks 12 **and** 15 declare expected non-zero commands |
| Section 4, EVID-03 | "VAL-002 / R-11 row" | **VAL-075** |
| Section 4, EVID-09 / EVID-21 / EVID-34 | "(R-07b)", "the R-02 producer", "owner per R-10" | task 10's all-satisfied state; the **four** DEC-224 producers; **task 14** generates and **task 18** regenerates |
| Section 5 | the unbounded `go test -race -coverprofile=coverage.out ./...`, `go test -tags integration -race ./cmd/fid/integration/...`, and the two package runs | the byte-identical bounded forms with `-timeout 30m` in 12.4; all four re-measured |
| Section 5 note | "Only task 15 declares expected non-zero commands (R-07d adds one to task 12)" | tasks **12 and 15** declare them; task 12's is applied and validated |
| Section 6.2, 6.3 | the first pass's gate timings | the 12.4 measurements; the first pass's 609 s / 568 s wall figures are superseded by 562 s / 554 s under the new command text |
| Section 6.5 | "none of the 20 new codes exists" | **19** are new; `requirement-not-found` already exists; `finding-unresolved` is no longer among them |
| Section 6.7 | "Runnable fixtures for VAL-070..072 - FAIL - no task owns a generator" | **PASS** - task 14 owns it, task 18 regenerates |
| Section 6.9 | "148 stale integration binary dirs" attributed to timeout kills | **20 directories and 306 MB leak per clean run**; see 12.4 and R-15 |
| Section 7, QA-04 | one item covering VAL-002 including the baseline diff | split: QA-04 keeps the four discriminator probes for VAL-002; **QA-40** owns the baseline diff for VAL-075 |
| Section 7, QA-13 / QA-16 / QA-21 | "per the R-04 disposition", "per R-03", "Stop only in the middle state per the R-02 producer" | **DEC-225** four recordings; **DEC-226** precedence; **DEC-224** four sequences plus the repeat case and every non-qualifying state |
| Section 9 preamble | "Fixtures come from the owner assigned under R-10" | task 14 generates all four Human rows' fixtures; task 18 regenerates them against the final binary |

### 12.9 Deterministic Checklist Additions

Appended to section 7. Same rules: executed by QA Lead at final review, `PASS`
without cited evidence is invalid.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-40 | `git diff 0195897a347dfe1872db0f7be3800c5c3ccdefcf -- migrations/0013_p21_operational_kernel.up.sql migrations/0015_p22_recovery.up.sql` is exactly the seven accepted hunks, every one inside `envelope_obligations`, `work_claims`, `stops`, or `attention_items`, with `execution_attempts` lines 222-223 untouched `(qa-lead)` | VAL-075 | full diff output |
| QA-41 | For each of the nine Phase 23 Event types, an unknown effect table, a missing required effect, and a malformed body are rejected, plus one unknown Event type; projection digest byte-identical before and after each `(qa-lead)` | VAL-076 | 28 probe outputs and digest pairs |
| QA-42 | A Review recorded independently, whose Actor then takes a mutating Claim on the same Issue and revision, cannot be linked (`review-not-independent`, no Event), while a Review whose Actor never took one links `(qa-lead)` | VAL-077 | both outputs, journal line counts |
| QA-43 | Each of the four Condition types resolves only on its own qualifying fact Event and rejects a foreign one; a Condition with two or with zero non-null fields is rejected with no Event `(qa-lead)` | VAL-078 | per-type paired outputs, rejection outputs |
| QA-44 | The phase-end run of both bounded gates records per-package times; tagged integration under 900 s, no ordinary package over 600 s, no gate terminated by timeout; every completion report carries its three package times `(qa-lead)` | VAL-079 | gate transcript with times, report excerpts |
| QA-45 | `/tmp` free space recorded before and after the phase-end tagged gate, and the `fid-integration-binary-*` directory count recorded at both points `(qa-lead)` | no VAL | `df` and `ls` output pairs |

### 12.10 Revisions Requested From The Tech Lead (second pass)

Must-fix (M) blocks Ready. These are bounded, decision-free planning
corrections; neither is a blocker candidate and neither needs Human authority.

| ID | Sev | Where | Finding and proposal |
|---|---|---|---|
| R-14 | M | task specs 4, 5, 6, 8, 9, 10, 11, 12, 16; `execution-scope.md` 6; `schema.md` 6.2 | **Four closed Phase 23 registries have no owning task.** `execution-scope.md` 1 is explicit that "every task owns its **named** implementation paths", and the Tech Lead accepted exactly this finding class twice already (R-08, R-09). Four production files that an owning task must edit are named by no task: (a) **`internal/workflow/stop.go`** - `stopCodes` is a closed 7-member map (`:114`), `conditionTypes` a closed 6-member map (`:120`), `Condition` a closed 4-field struct (`:26`), and `validStopDetails` a closed per-code key map (`:409`). Tasks 6, 9, 10 and 11 each raise a new Stop code or a new Condition type; without this file `validateStopRaise` (`:364`) and `validateCondition` (`:396`) reject every Phase 23 Stop. No task names it; task 11 names only `cmd/fid/stop.go` "result shapes only". (b) **`internal/journal/current_event_value_catalog.go`** - line 15 holds the closed `Condition` value grammar with its six accepted types, and lines 43-51 `currentStopDetailShapes`, the closed per-Stop-code details grammar that `StopCode` validation itself keys on (`current_event_value_types.go:89-91`). Named by no task. (c) **`internal/journal/current_event_columns.go`** - `currentEventTables` (`:5`) registers 22 tables with their column types; `validateCurrentEffect` returns `unknown table` for any effect on an unregistered table (`current_event_validation.go:160-164`), so all ten new Phase 23 tables must be registered. Named by no task. (d) **`internal/export/archive_semantics_catalog.go`** - `semanticEventRules` is the archive-side mirror of the journal catalogs; its 44 keys are exactly the union of `currentEventRules` (22) and `retainedCurrentRules` (22), and `validateEffects` (`archive_semantics.go:337-361`) rejects any Event whose family is unknown or whose effect table/op is not accepted. The nine new Phase 23 rules and the three amended ones (`envelope.activated` at task 9, `evidence.linked` and `artifact.availability-observed` at task 10 under DEC-224) must be mirrored there, or `export validate` and `import` reject the journal of any workspace that recorded a Phase 23 Event or raised a readiness Stop. The file is named only by task 12, for a **deletion**. No test asserts parity between the two catalogs, so the divergence is silent until VAL-064 at task 16 - by which point the cause sits outside task 16's ordinary scope, which is the precise failure mode R-09 was accepted to prevent. Additionally `internal/journal/current_event_catalog.go` is named by tasks 8, 9, 10 and 12 but **not** by tasks 4 and 5, which own the new `review.recorded` and `finding.addressed` rules and their `currentEventCommands` entries. **Proposal:** add `internal/workflow/stop.go`, `internal/journal/current_event_value_catalog.go` and `internal/journal/current_event_columns.go` to the file scopes of tasks 4, 5, 6, 8, 9, 10, 11 and 12 as each one's own additions require; add `internal/journal/current_event_catalog.go` to tasks 4 and 5; and give `internal/export/archive_semantics_catalog.go` an explicit owner for the additions and the three amendments - either each Event's own task, or task 16 as a single named owner, stated in `execution-scope.md` 6 either way. No assertion changes. |
| R-15 | M | task 1; `execution-scope.md` 7; CF-042 | **The accepted CF-042 mitigation does not address the measured cause and does not stop the measured leak.** CF-042 and task 1 record the cause as "`TestMain` removes its directory only on normal exits, so a timeout kill or interrupt leaks a ~16 MB directory", and the mitigation as pruning directories "older than 24 hours". Fresh measurement (12.4): one **clean** tagged run - `exit_code: 0`, 551.480 s, no kill, no timeout - left **20** directories and **306 MB** in a `/tmp` that held 3.6 MB beforehand. The dominant cause is not timeout kills: five tests deliberately re-execute the test binary as a child process, each child re-enters `runIntegrationSuite`, creates its own directory and rebuilds the 16 MB binary, and the crash-process children by design never reach `RemoveAll`. Because the leak is per-run rather than per-kill, an age bound of 24 hours prunes **nothing** during a working day: Phase 23 runs the tagged gate at least 18 times for tasks plus 6 window reviews, the task-11 checkpoint and final review, so roughly 26 runs at 306 MB is about 7.8 GB against a 3.9 GB RAM-backed `/tmp`. On the measured rate `/tmp` fills around the twelfth run - well before any timeout - and a full `/tmp` is a gate-blocking environment failure, not a test failure. **Proposal:** keep the CF-042 split, but correct the recorded cause and make the task 1 mitigation match it. Either (i) have `runIntegrationSuite` publish the built binary path in an environment variable and have a child that finds it set reuse it, creating and removing nothing - which also removes about twenty redundant `go build ./cmd/fid` invocations per run and so helps the DEC-227 budget; or (ii) keep pruning but bound it by "not the current process's directory" rather than by age. Option (i) is the causal fix and stays inside `cmd/fid/integration/main_test.go` plus the five re-executing tests. Record the measured 20-directory, 306 MB per-run figure in `execution-scope.md` 7 in place of the "148 directories / 2.3 GB" framing, which described an accumulated symptom rather than the rate. |

Improvements (the Tech Lead may apply or decline; none blocks Ready).

| ID | Where | Proposal |
|---|---|---|
| I-17 | `schema.md` 6.2 | The section says "`currentEventRules` gains exactly these rows", but the three amended families (`envelope.activated`, `evidence.linked`, `artifact.availability-observed`) live in **`retainedCurrentRules`** (`current_event_catalog.go:77`), which has its own validation path (`current_event_validation.go:39-60`). Name both maps so the DEC-224 rows are added in the right one. |
| I-18 | `contracts.md` 2 | "each new code takes a rank ... immediately after `obligation-evidence-absent`; the Phase 21/22 ranks are unchanged" is unsatisfiable read numerically: `precedence` (`evaluation.go:380`) is an integer map where seven codes already share rank 8 and `receipt-unsuccessful` is 9. State that only the relative order is accepted and that renumbering the unexported ranks is permitted. |
| I-19 | `public-interface.md` 5.1 | Precedence step 1 says "every retained Phase 21/22 precondition and authority failure", and step 5 names `obligation-unsatisfied`, which is itself retained. Write step 1 as "every retained Phase 21/22 precondition and authority failure **other than `obligation-unsatisfied`**". |
| I-20 | task 13 | The rename list names `reasonCodes` in `cmd/fid/evidence_test.go` "where they name the failure set". It does not: at `:429` it is the `assertEvaluation` parameter compared against `evaluation reasons` (`:438`). The qualifier makes the instruction vacuous rather than wrong, but naming a real identifier invites the wrong rename - exactly the collision DEC-226 removes. Drop it, or rename it to `evaluationReasons` as an explicit evaluation-set clarification. |
| I-21 | task 4 | Name `idempotency-conflict` in task 4's `internal/cli/codes.go` list. It is an accepted Phase 21 code with **no** entry in `reasonCodeDefinitions` and no adapter mapping today (12.6), while VAL-019 and VAL-057 require a real-process reproduction. Registering an already-accepted code is conformance and is in task 4's declared scope; naming it removes any doubt that it is. |
| I-22 | VAL-006 | The evidence method now enumerates the protected-Decision counts (0, 16, 17) but no longer names the Obligation counts. Add 0, 1, 16 and 17 so the row's "accepts 1..16 Obligations" clause has both boundary states, matching EVID-07. |
| I-23 | VAL-034 | The row proves the DEC-228 "a subject need not resolve" rule only as a side effect: a `requirement`-subject Concern at task 8 cannot resolve, because task 12 has not run. Say so in the pass criterion so the proof is deliberate rather than incidental. |
| I-24 | task 10 | Record the fixture technique for the `artifact.availability-observed` raising sequence: remove the retained content-addressed blob under `.fid/artifacts/<digest>`, observe `missing`, restore it, observe `available`. Proven reachable in 12.5; without it the sequence reads as unreachable and the Developer has no accepted way to produce the qualifying state. |
| I-25 | `public-interface.md` 5.1 | The grouped `*-not-found` row names a category ("the read and mutation leaves naming that entity") rather than an operation, while VAL-056 requires "one reproduction per code through the producer named in section 5.1". Pin the exact leaf per code - `deviation-not-found`, for instance, has no `deviation get`; its only producer is `deviation.record` with an unknown `supersedes_deviation_id`. |
| I-26 | task 1 | Two facts measured in 12.2 that the task should know: the product exposes no single-step down (`internal/db.MigrateDown` runs `m.Down()` to zero), so the 0017 down evidence is either the full teardown or a test-local `Steps(-1)`; and a failed migration leaves `schema_migrations` **dirty** (observed `16\|1`), which is the exact failure state VAL-005 rules out and is worth naming in the fail criterion. |

### 12.11 Blocker Candidates

**None.** The first pass's three candidates (R-01, R-02, R-04) were resolved by
DEC-223, DEC-224 and DEC-225, all Human-approved on 2026-09-19, and section 12
confirms each against fresh evidence. R-14 and R-15 change no protected command,
no assertion semantics, no accepted status or scope, no threshold, budget or
tolerance, and no policy: R-14 adds file paths to task scopes that already own
the behavior, and R-15 corrects a recorded cause and a test-infrastructure
bound inside a file task 1 already owns. No `ESC` identity is requested and none
is minted here.

### 12.12 Readiness Verdict

**NOT READY.** Readiness condition 3 is executed and complete, and conditions 1
and 2 are satisfied, but two Must-fix revisions remain open.

What is proven ready:

- The five Human-approved decisions hold under fresh, independent evidence. In
  particular DEC-223 is now proven on the real migration chain, under
  `modernc.org/sqlite`, in place, on a workspace whose receipt Evidence carries
  Artifact links produced by the product itself - and the pre-amendment text is
  proven to fail there, in **both** directions, leaving the projection dirty.
- Every command in every task spec is concrete, syntactically valid, and was
  executed or feasibility-proven at its new text. Both DEC-227 budgets hold with
  21-39% margin.
- All four DEC-224 raising sequences are reachable at task 10, including
  `artifact.availability-observed`, which was reproduced live.
- The contract is internally consistent at 79 rows: all Ready, one validator
  each, every row cited by a task spec, no task citing an undefined row, and
  every restatement a strengthening.

What blocks Ready: **R-14** and **R-15**. Both are bounded and decision-free.
No task may be dispatched until the Tech Lead dispositions them, applies the
accepted edits, and a fresh QA Lead readiness refresh confirms them; then the
Phase Runner's readiness check (condition 4) can run.

No assertion was weakened to reach this verdict, and no `VAL-NNN` was narrowed.

## 13. Readiness Confirmation - DEC-229 / DEC-230 (2026-09-19)

Readiness condition 3, second cycle. A fresh QA Lead session confirmed the
edits the Tech Lead applied in `5ff523a` under DEC-229 and DEC-230, per
`tl-qa-confirmation-brief.md`. Base: branch `phase-23` at `c0a0876`,
`git status --porcelain` empty before and after this pass.

Scope bounds honoured. R-01..R-13 and DEC-223..228 were not re-opened. The
DEC-227 budgets, the DEC-223 migration-chain proof, and the DEC-224 raising
sequences were not re-measured: no product file changed in `5ff523a` (all 19
edited files are under `project-docs/docs/`, confirmed by `git show --stat`),
so section 12's measurements stand except where restated in 13.5 below.

Method. Every claim below was checked against the product source or the task
specs themselves, never against the brief's summary. Two throwaway probes were
executed against the real packages through `go test -overlay`, which injects a
test file without writing to the working tree; both probes and their overlays
live in the session scratch directory and the tree stayed clean.

### 13.1 The Four Judged Items

| Item | Verdict | Evidence |
|---|---|---|
| **3.1 Condition-type sequencing** | **CONFIRMED - the split is complete and correctly ordered** | Phase 23 adds exactly three StopCodes and four Condition types (`public-interface.md` 2). Registration per the task specs: task 6 `review-rework-exhausted` and *no* Condition type, because its resolution condition is the retained `budget-observed` (`task-6.md:64-72`); task 9 `protected-decision-required` + `decision-satisfied` + the `decision_class` field *with its exactly-one-non-null rule in `validateCondition`* (`task-9.md:65-73`); task 10 `human-outcome-required` + `human-outcome-recorded` + `outcome_class` with the same rule (`task-10.md:83-90`); task 11 the remaining `review-recorded` and `finding-resolved` plus `finding_id`, and no Stop code (`task-11.md:55-65`). That is 3/3 codes, 4/4 types, and 3/3 new `Condition` fields, each field's rule added by the task that adds the field. No task raises a Stop whose code or condition type is registered later: `budget-observed` is already in the 6-member `conditionTypes` (`stop.go:120`); tasks 9 and 10 register their own at the raising task; tasks 6, 9, 10 all precede 11. A search of all eighteen specs for the three new codes and four new types returns nothing outside tasks 6, 9, 10, 11 except two descriptive mentions in `task-13.md:34,44`. The Tech Lead's correction of the refresh's under-specified "as each one's own additions require" is right, and is the only reading compatible with `validateStopRaise` (`stop.go:363-370`), which rejects an unregistered code and calls `validateCondition` before any commit. Two supporting facts also hold: every Phase 23 Stop takes subject Envelope, which `validSubject` (`stop.go:377`) already accepts, and every declared blocked action is already in `blockedActionSet` (`stop.go:125`), so neither registry needs an addition. |
| **3.2 Task 1 as owner of `current_event_columns.go`** | **CONFIRMED - the executable proof lands at task 1, and no effect can precede registration** | `TestCurrentEventColumnDomainsMatchAcceptedDDL` exists today at `internal/journal/current_event_validation_test.go:228`. It iterates `currentEventTables` and introspects the live migrated database via `PRAGMA table_info`, checking column presence, integer-vs-text type, nullability and primary-key membership, and it fails closed: a registered table absent from the DDL yields zero rows, so `len(found) != len(definition.columns)` fires `t.Fatal`. Task 1 ships the 0017 DDL and registers the ten tables in the same task, so the proof genuinely executes there, inside `internal/journal`, which task 1's `go test -race -timeout 30m -coverprofile=coverage.out ./...` runs. No effect on a new table can precede registration: task 1 is the first task, and its own fence rejects all nine Event types until a later owner deletes its key. The ten names in `task-1.md:87-93` are exactly the new tables reachable from `schema.md` 6.2 - the nine required/optional new-rule tables plus `envelope_protected_decisions`, which only task 9's `envelope.created` amendment writes; registering it at task 1 is earlier than its first effect, not later. Assigning this to one task rather than five is also the safer choice for a single shared map. Task 15's rename of the `judgment_concerns` key is correctly the sole exception (`task-15.md:63`). |
| **3.3 The rejection at tasks 4 and 5** | **CONFIRMED - the rejection is correct; the refresh's sub-item was wrong** | Verified directly against `task-1.md:60-63`: "Register the nine Event types **and their operation mappings** in `internal/journal/current_event_catalog.go` with the effect rules in `schema.md` section 6.2", and `current_event_catalog.go` is in task 1's file list (`task-1.md:86`). `schema.md` 6.2 gives `review.recorded` and `finding.addressed` their complete rows, including every optional effect those families ever take - `review.recorded` already carries `budget_samples/insert`, `stops/insert`, `attention_items/insert`, `attention_stop_links/insert`, `stop_resolutions/insert` and `attention_items/update`, which are what tasks 6 and 10 later need. The two maps tasks 4 and 5 would have touched are `currentEventRules` and `currentEventCommands` (`current_event_catalog.go:26,50`); task 1's wording covers both. So tasks 4 and 5 need no catalog edit and the file would only have created a shared-map conflict. My own first-pass inference from file lists alone was the error. |
| **3.4 DEC-230** | **CONFIRMED - the premise holds; I agree with the removal over the leaf** | Verified independently of the brief. `public-interface.md` 4 gives Decision References exactly `decision.record-reference` and `decision.list`; there is no `decision get` row. `decision.record-reference`'s `id:ID` is the caller-supplied identity of the record being created, not a lookup key, so it cannot return not-found. A repository-wide search of the phase documents for `decision_ref_id` returns exactly one hit, `public-interface.md:108`, and it is a **result** field of `ProtectedDecision` (`ProtectedDecisionInput + {status,decision_ref_id:?ID,human_outcome_id:?ID}`), never an input; `ProtectedDecisionInput` has no such field. No other leaf's input names a Decision Reference ID. `linked_to_decision` is confirmed a transition rule "evaluated against Decision References in the Issue's scope" (`contracts.md:359-361`), reporting satisfied/unsatisfied - not a not-found error. The row was therefore unsatisfiable for VAL-056 exactly as DEC-230 states. I also agree with the choice: adding a `decision.get` leaf would expand the public surface to rescue dead vocabulary, and `deviation-not-found` and `outcome-not-found` are precedents in this very phase for a not-found code reachable only through another leaf's input field. The removal is the smaller change and is one line to reverse. |

### 13.2 Mechanical Re-checks

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| `public-interface.md` 5 and 5.1 well formed | full read plus set comparison | PASS | Section 5 declares **19** codes (13 single rows plus the grouped not-found row's six); section 5.1 has **19** producer rows; the two sets are identical. The grouped row now lists six codes, not seven. The table is well formed - no ragged row, every code `retryable:false`, every details shape closed. |
| Each pinned producer exists | leaf and input-field check against `public-interface.md` 3 and 4 | PASS | All six pinned not-found producers resolve: `review.get` is a declared leaf; `requirement.get` is a declared leaf; `finding.address` takes `finding_id`; `review.record` takes `resolves_finding_ids`; `concern.dispose` takes `concern_id`, `supersedes_concern_id` and `supersedes_finding_id`; `deviation.record` takes `supersedes_deviation_id`; `evidence.link` takes `review_id` and `human_outcome_id`. The two "there is no ... get leaf" notes are true of section 4. The other thirteen producers are existing or newly declared leaves; `source-digest-mismatch`'s two producers both take `source:Source`. |
| New-code count | re-derivation against `internal/cli/codes.go` | PASS - **restated to 18** | 19 codes declared in section 5, of which `requirement-not-found` already exists, so **18 are new**. This supersedes 12.6's "19" exactly as 12.8 superseded 6.5's "20". DEC-230's "the added set becomes 18" uses the same sense and agrees. |
| 12.7 file-scope table, four new rows | re-derived from the task specs, not the brief | PASS | `internal/workflow/stop.go` -> tasks 6, 9, 10, 11. `internal/journal/current_event_value_catalog.go` -> tasks 6, 9, 10, 11. `internal/journal/current_event_columns.go` -> tasks 1, 15. `internal/export/archive_semantics_catalog.go` -> tasks 1, 9, 10, 12, 15. Each derivation is from the task's own `Files To Create Or Modify` bullet and matches `execution-scope.md` 6. The R-08 and R-09 rows are untouched by this commit and stand. |
| Contract row count and integrity | mechanical re-run on the amended file | PASS | **79** `VAL-NNN` rows, **79** unique IDs. Validator column extracted per row: **40** QA Lead + **35** QA Reviewer + **4** Human = 79, and zero rows whose validator column is anything else, so one validator per row holds exactly. **0** rows not `Ready`. The four Human rows are VAL-070..073. |
| Contract citation closure | defined-vs-cited set comparison across all eighteen specs | PASS | Defined set and cited set are identical, 79 members each; zero uncited, zero undefined. |
| Contract diff is two rows, no weakening | `git show 5ff523a -- qa/validation-contract.md` | PASS | Exactly two rows changed and nothing else. **VAL-006** gains "at Obligation counts 0, 1, 16, and 17" in the evidence method and "1 and 16 Obligations are accepted and 0 and 17 are rejected" in the pass criterion (I-22) - a strengthening that gives the "accepts 1..16" clause both boundary states in both directions. **VAL-034** makes the DEC-228 unresolvable-subject proof deliberate in the pass criterion and adds "or the unresolvable `requirement` subject is rejected" to the fail criterion (I-23) - a strengthening that names both states. No row added, removed, renumbered, reworded elsewhere, or weakened. |
| Gate reach for every file gained | `go list` plus per-task command-block read | PASS | `go list ./...` returns `internal/workflow`, `internal/journal` and `internal/export`, and each has in-package tests (`stop_test.go`, `current_event_validation_test.go`, `archive_semantics_test.go`). Tasks 1, 6, 9, 10, 11 and 15 each carry the byte-identical `go test -race -timeout 30m -coverprofile=coverage.out ./...`, which executes all three. `cmd/fid/integration/main_test.go` is `//go:build integration`; `go list -tags integration ./cmd/fid/integration/...` resolves the package, and task 1 carries `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` and declares the tagged gate mandatory from that task. Confirmed, not assumed. |
| I-17 applied | `schema.md` 6.2 read | PASS | The amended-rules table is now explicitly prefaced "the last three families live in `retainedCurrentRules`, which has its own validation path, so add their rows there and not to `currentEventRules`", and the standing mirror rule is stated in the same place. |
| I-18 applied | `contracts.md` 2 read | PASS | `contracts.md:83-89` now states the four-code order and then "Only the relative order is accepted (I-18)", naming the integer `precedence` map and permitting renumbering. Matches `evaluation.go:380`. |
| I-19 applied | `public-interface.md` 5.1 read | PASS | Step 1 reads "every retained Phase 21/22 precondition and authority failure **other than `obligation-unsatisfied`**", so the five-step order is determinate without relying on step 5's mention. |
| I-20 applied | `task-13.md` read | PASS | `task-13.md:67-69` now says "Rename nothing in `cmd/fid/evidence_test.go`: its `reasonCodes` is the `assertEvaluation` ..." - the vacuous instruction is gone and the correct sense is recorded. |
| I-21 applied | `task-4.md` read | PASS | `internal/cli/codes.go` list now names `idempotency-conflict` with the conformance rationale and the VAL-019 / VAL-057 justification. |
| I-24 applied | `task-10.md` read | PASS | The blob remove/restore technique is recorded with the `.fid/artifacts/<digest>` path, the `referenced_artifacts/update` precondition, and the note that `missing -> available` is exactly one qualifying Event. |
| I-25 applied | `public-interface.md` 5.1 read | PASS | Six pinned producer rows replace the grouped category row; see the two checks above. |
| I-26 applied | `task-1.md` read | PASS | `task-1.md:48-55` records both facts: no single-step down (`internal/db.MigrateDown` runs `m.Down()` to zero), and the dirty `schema_migrations` failure state with the required `17\|0` / `16\|0` end states. |
| `archive_semantics_catalog.go` mirror coverage | task-spec derivation against the standing rule | PASS | Every journal Event-rule change this phase has a named mirroring task: task 1 the nine new rules, task 9 `envelope.created` and `envelope.activated`, task 10 `evidence.linked` and `artifact.availability-observed`, task 12 the `requirement.synced` deletion, task 15 the concerns rename. One arithmetic correction to 12.10: `semanticEventRules` does hold **44** keys (measured), but it is not the union of `currentEventRules` (22) and `retainedCurrentRules` (**24**, not 22), which would be 46. The measured key count and the disposition are unaffected. |

### 13.3 R-15 Specifically

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| The fix stays inside `main_test.go` | source check of all five re-executing call sites | PASS | All five spawn with `append(os.Environ(), ...)`: `p22_budget_fault_test.go:209`, `p22_activation_fault_test.go:135`, `p22_human_fixtures_test.go:69`, `p22_force_budget_fault_test.go:35`, `p22_recovery_fault_test.go:53`. Each re-executes the test binary itself (`os.Args[0]` or `os.Executable()`), so each child re-enters `TestMain` -> `runIntegrationSuite`, confirming the per-run mechanism. A variable published by the parent is inherited without editing any of the five. |
| The causal claim holds | `main_test.go` read | PASS | `runIntegrationSuite` unconditionally `MkdirTemp`s `fid-integration-binary-*`, builds a ~16 MB binary into it, and calls `RemoveAll` only after `m.Run()` returns - so an abnormally exiting child leaks by construction. Task 1's instruction (publish the path; a process that finds it set and pointing at an existing executable reuses it and neither builds nor removes; only the builder removes; a stale path falls back to building) is the matching causal fix and is fully specified. |
| Option (ii) correctly rejected | reasoning check against the same source | PASS | A "not the current process's directory" bound would delete a concurrent run's live binary, because `integrationBinary` is per-process and nothing distinguishes a sibling run's directory from a stale one. Age remains the only safe bound. Keeping the 24-hour prune as a residual net for a parent killed mid-run is right; note that today's `main_test.go` contains no prune at all, so task 1's "Keep" means "retain the I-16 instruction", and the required change list does state it. |
| Cause corrected in all four places | cross-document check | PASS | `task-1.md:104-127`, `execution-scope.md` 7, `v1-convergence-inventory.md:72-75`, and CF-042 in both `decisions.md` and `decisions_active.md:159`. All five quote the measured 20 directories / 306 MB per clean run against a 3.9 GB `/tmp`, the ~26-run / ~7.8 GB projection, and the per-run-not-per-kill cause. The superseded "148 directories / 2.3 GB" framing appears nowhere. |
| Task 1 records the measurement | `task-1.md` read | PASS | It requires the `fid-integration-binary-*` directory count and `/tmp` free space immediately before and after the tagged gate, with "a clean run must leak none". This overlaps QA-45 by design, as the brief says. |
| **Ownership question** | judgment | **CONFIRMED at task 1** | Task 1 already owns `cmd/fid/integration/` fixtures, is the first task to run the tagged gate, and declares that gate mandatory from itself. Any later owner leaves every intervening task exposed to the same leak, and section 12 measured `/tmp` filling around the twelfth run. The change is decision-free, test-infrastructure only, confined to one file, and needs no edit to the five fault tests. It is a defensible addition to a schema baseline task precisely because it is the earliest point that protects the rest of the phase. |

### 13.4 Statements Superseded By This Section

| Where | Superseded statement | Replacement |
|---|---|---|
| Section 12.6, 12.8 | "20 codes declared in section 5 ... so **19** are new"; "all 19 new codes have an owning task" | **19** declared, **18** new (DEC-230 removes `decision-not-found`). The owning-task check re-runs clean at 18: task 9 now names two codes, not three. |
| Section 12.7 | `semanticEventRules`' "44 keys are exactly the union of `currentEventRules` (22) and `retainedCurrentRules` (22)" | 44 keys measured; `retainedCurrentRules` holds **24**, so the union is 46 and the two catalogs are not key-identical today. The mirror-ownership finding is unaffected. |
| Section 12.10, 12.12 | R-14 and R-15 open | Both dispositioned under DEC-229 and confirmed here. The open item is **R-16** in 13.5. |
| Header status block | "two Must-fix revisions remain (section 12.10)" | one Must-fix remains: **R-16** (section 13.5). |

### 13.5 Revisions Requested From The Tech Lead (third pass)

| ID | Sev | Where | Finding and proposal |
|---|---|---|---|
| R-16 | M | `internal/journal/current_event_json.go`; `internal/journal/current_event_result_projection.go`; `internal/journal/current_event_value_catalog.go`; task specs 1, 4, 5, 8, 9, 10, 12; `execution-scope.md` 6 | **The typed-result registries for the nine new Event families have no owning task, and the first Phase 23 Event to clear its fence is rejected at append time.** This is the same class as R-08, R-09 and R-14, found while confirming R-14, and it is larger than any of them. Mechanism, verified in source: `project_effects.go:20` calls `validateCurrentAppliedResults` on every append, and `validateCurrentEvent` (`current_event_validation.go:37-63`) routes a family to the **strict** path if and only if it is in `currentEventRules`. Task 1 registers all nine Phase 23 types there (`task-1.md:60-63`, `schema.md` 6.2), and `schema.md` 6.3 requires a typed snapshot for every Phase 23 mutation, so each of the nine must clear three further closed registries that no task names: (a) the inline `fields` map in **`current_event_json.go:135-141`**, which maps Event type to its typed-result grammar and fails with `unregistered current result family`; (b) `currentExpectedResult`'s switch in **`current_event_result_projection.go:69-121`**, whose `default` branch fails with `unknown current result` for anything not prefixed `activation.`; and (c) **`currentValueShapes`** in `current_event_value_catalog.go`, from which `currentProjectedValue` (`:175-203`) reads the projected field list - an unregistered kind yields an empty list and a guaranteed snapshot mismatch. **Executable proof** (`go test -overlay`, `internal/journal`, working tree untouched): for all nine types, `validateCurrentResultJSON` returns `unregistered current result family` and `currentExpectedResult` returns `unknown current result`; and `currentValueShapes` (35 entries) contains none of `Review`, `Finding`, `FindingDisposition`, `Concern`, `ConcernDisposition`, `Deviation`, `DecisionRef`, `HumanOutcome`, `ProtectedDecision`, `RequirementRef`, `RequirementLink`, `Coverage`. A repository-wide search of all eighteen task specs returns **zero** mentions of `current_event_json.go` or `current_event_result_projection.go`, and the four tasks that do name `current_event_value_catalog.go` (6, 9, 10, 11) name it only for the `currentStopDetailShapes` entry and the `Condition` value grammar. Task 1 stays green because its fence rejects all nine and VAL-003 asserts exactly that; the failure therefore first lands at **task 4**, when `review.recorded` is unfenced, with its cause in files outside task 4's declared scope - the precise failure mode R-09 and R-14 were accepted to prevent, and here it blocks every remaining domain task rather than one archive check. Two consequential sub-items in the same files: `currentValueShapes["Claim"]` is `mode:=mutation`, which `public-interface.md` 2 widens to `mutation\|review` (task 3), and `currentValueShapes["Budget"]` / `["ImmutableBudget"]` close `unit` without `review-cycles` and `exhaustion_action` without `review-rework-exhausted`, which `public-interface.md` 2 adds (task 6); and if any new shape is written as `subject:Subject`, the closed `Subject` enum `type:=run,envelope,activation` needs the six added SubjectTypes. **Proposal:** apply the R-14 standing-rule pattern. Give `current_event_json.go` and `current_event_result_projection.go` explicit owners - a family's typed-result grammar, expected-result branch and value shapes are added by the task that unfences that family (tasks 4, 5, 8, 9, 10, 12), or all nine at task 1 alongside the `currentEventRules` registration it already owns - and state the choice in `execution-scope.md` 6 as a standing rule, as DEC-229 did for `semanticEventRules`. Add the `Claim` domain to task 3's scope and the two Budget domains to task 6's. Note that `currentProjectedValue`'s `switch kind` carries per-kind fill helpers for composite values (`Run`, `Activation`, `Stop`/`Attention`, `Receipt`); Phase 23's composite values - `Review` with `findings:[]Finding` and `coverage:Coverage`, `Finding` with `dispositions:[]FindingDisposition`, `Concern` with `disposition:?ConcernDisposition` - will need the same, so the disposition should say which task owns that too. **No assertion changes, no protected command, threshold, budget or public surface changes**: every value shape is derived from the already-accepted `public-interface.md` 3, and the typed-snapshot requirement is already accepted under DEC-216 and `schema.md` 6.3. This is an ownership and file-scope correction under planning authority, exactly like R-14. |

Improvements (the Tech Lead may apply or decline; none blocks Ready).

| ID | Where | Proposal |
|---|---|---|
| I-27 | `qa/validation-contract.md` header | The provenance block still reads "Last refined: 2026-09-15 ... applying the QA Lead's first-pass revisions R-01..R-13 and improvements I-01..I-16", and still says "A fresh QA Lead readiness refresh (`qa-plan.md` section 11) is still required". Both are now stale: VAL-006 and VAL-034 were refined on 2026-09-19 under DEC-229 (I-22, I-23), and the refresh is section 12. Record the second refinement and drop the satisfied prerequisite. The rows themselves are correct; only the header misstates their history. Contract edits are Tech Lead-owned, so this is recorded here rather than applied. |
| I-28 | tasks 6 and 11 | Both still carry a directory-level bullet - `internal/journal/` "the `review.recorded` effect-rule optional rows" (`task-6.md:63`) and "the `human-outcome.recorded` optional effect rows" (`task-11.md:54`). Under DEC-229 task 1 registers each family's **complete** `schema.md` 6.2 rule, including every optional row those two tasks would have added, so both bullets are now no-ops that invite a Developer to edit a map task 1 owns. Delete them, or restate them as "no catalog edit; task 1 registered the complete rule". |
| I-29 | `schema.md` 6.2 amended-rules table | The `envelope.activated` owner cell reads "task 9 protected Decisions; task 10 readiness Stop", but the three optional Stop rows are identical for both and only task 9's spec adds them; task 10's catalog bullet correctly limits itself to `evidence.linked` and `artifact.availability-observed`. No gap exists - task 9 precedes task 10 - but the cell reads as though task 10 also edits that row. Say "rows added by task 9; also raised from this family by task 10". |

### 13.6 Blocker Candidates

**None.** R-16 changes no protected definition, command, threshold, budget,
tolerance, assertion semantics, accepted status, public surface, or policy. It
assigns existing files to tasks that already own the behavior, and every value
shape it needs is derived from the already-accepted `public-interface.md`
section 3 under the already-accepted DEC-216 / `schema.md` 6.3 snapshot rule. It
is a bounded, decision-free planning correction within Tech Lead planning
authority, as R-08, R-09 and R-14 were. No `ESC` identity is requested and none
is minted here. DEC-230's `decision.get` alternative remains a live Human-facing
choice, and this pass supports the Tech Lead's recommendation to leave the code
removed.

### 13.7 Readiness Verdict

**NOT READY.**

**Every edit this pass was asked to confirm is confirmed correct.** DEC-229 and
DEC-230 are both sound and faithfully applied. Specifically: the R-14 (a), (b)
and (d) dispositions are complete; the Condition-type sequencing the Tech Lead
imposed on the refresh's under-specified proposal is right and is the only
reading compatible with `validateStopRaise`; task 1 is the correct owner of
`current_event_columns.go` and its executable proof genuinely lands there; the
rejection at tasks 4 and 5 is correct and my own first-pass sub-item was wrong;
DEC-230's premise is independently verified and its choice is the better one;
R-15's cause is corrected in all four places, option (i) is the causal fix,
option (ii) is correctly rejected, and task 1 is the right owner. The contract
is unchanged at 79 rows with both edits strengthening. I-17 through I-26 are all
applied. Nothing in DEC-229 or DEC-230 needs to be revisited.

What blocks Ready is a **new** Must-fix surfaced while confirming R-14:
**R-16**. Registering the nine Phase 23 Event families in `currentEventRules` -
which task 1 must do, and which DEC-229 correctly reasoned about for the column
and archive catalogs - also routes them onto the journal's strict validation
path, and three further closed registries on that path have no owning task. The
first unfenced Phase 23 Event, `review.recorded` at task 4, is rejected at
append time with its cause outside task 4's scope. This is proven by execution,
not by reading. It is bounded and decision-free, and it belongs to the Tech Lead
under the same planning authority as R-14.

No task may be dispatched until the Tech Lead dispositions R-16, applies the
accepted edits, and a fresh QA Lead confirmation pass verifies them; then the
Phase Runner's readiness check (condition 4) can run.

No assertion was weakened to reach this verdict, and no `VAL-NNN` was narrowed,
re-scoped, or restated to reach Ready.

## 14. Readiness Confirmation - DEC-231 (2026-09-19)

Readiness condition 3, third cycle. A fresh confirmation of the edits applied in
`8cd613e` under DEC-231, per the superseding `tl-qa-confirmation-brief.md`.
Base: branch `phase-23` at `8cd613e`, `git status --porcelain` empty before and
after this pass.

Scope bounds honoured. Sections 12 and 13 are not reopened. No product file
changed in `8cd613e` (all 15 edited files are under `project-docs/docs/`,
confirmed by `git show --stat`), so the DEC-227 budgets, the DEC-223 migration
chain, the DEC-224 raising sequences, the failure-code sets and the contract row
counts were not re-measured. Two throwaway probes ran against the real
`internal/journal` package through `go test -overlay`, which injects a test file
without writing to the working tree.

### 14.1 The Judged Item - Owner Is The Unfencing Task, Not Task 1

**CONFIRMED. The disposition is correct, and it is the right reading of
DEC-229 rather than a departure from it.**

The executability test is the correct test, and it separates the two cases
cleanly. `TestCurrentEventColumnDomainsMatchAcceptedDDL` introspects the 0017
DDL through `PRAGMA table_info`, needs no Event appended, and fails closed -
section 13.1 confirmed exactly that, which is why task 1 is right for
`current_event_columns.go`. The three registries here execute only inside
`validateCurrentEvent`'s strict branch and `validateCurrentAppliedResults`, both
of which require an appended Event of the family, and `unavailableJudgmentEvents`
makes that impossible until the unfencing task. Placing them at task 1 would
write nine grammars, nine branches, fourteen shapes and three fill helpers that
task 1's own gate could never execute, and would convert task 4's red test into
an inherited defect. The parallel the Tech Lead draws to DEC-229's Condition-type
rule is exact: registration lands where it first executes.

**Is there a proof that would land at task 1 and that the Tech Lead missed?** No.
I looked for the `current_event_columns.go` analogue - a test that iterates one
of these three registries and cross-checks it against something task 1 ships.
There is none: `currentValueShapes` is checked only by being *read* during an
append, the `fields` map only by `validateCurrentResultJSON` on an append, and
`currentExpectedResult` only by `validateCurrentAppliedResults` on an append. A
task-1-time parity test could only be written by inventing a new cross-check
between `currentEventRules` and the other three maps, which is new product
behavior, not registration. The Tech Lead's search finding no such test is
confirmed.

### 14.2 The Reachability Sweep

**Is there any path that reaches one of the three registries before its family's
fence key is deleted?** I checked all four the brief names, and found none.

| Path | Result | Evidence |
|---|---|---|
| A retained family projecting a new value | **Cannot reach them** | `validateCurrentEvent` routes a family to the retained branch when it is absent from `currentEventRules` and returns at `current_event_validation.go:64` - before `validateCurrentJSON` is ever called. `validateCurrentAppliedResults` independently gates on `currentEventRules` membership (`current_event_result_projection.go:12`). So `evidence.linked`, `envelope.activated`, `envelope.created` and `artifact.availability-observed` - the four families tasks 9 and 10 amend - never touch any of the three. |
| Export or import | **Different registry** | `internal/export` keys on `semanticEventRules` (44 entries, measured), which DEC-229 already assigned. Neither `currentValueShapes` nor either result registry is referenced from `internal/export`. |
| Rebuild | **Cannot reach them** | A rebuild replays the Events that exist. A fenced family cannot be appended, so no Event of it exists to replay. |
| Replay of a stored snapshot | **Cannot reach them** | Replay returns the stored `idempotency_results` row; it re-derives nothing, and the row can only exist if the original append already succeeded, which the fence prevents. |

The unfencing task is therefore not too late for any of the three. The
disposition holds on its own terms.

**But the sweep did not stop there.** The brief's own section 6 asks whether a
fourth closed registry on the strict path has no owning task. It does - see
R-17 - and a fifth shared registry is being edited in a way that drops an
accepted Phase 22 constraint - see R-18.

### 14.3 Mechanical Re-checks

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Nine fence keys map one-to-one onto six unfencing tasks | `execution-scope.md` 4 and 6 against `schema.md` 6.1 and the task specs | PASS | 4 `review.recorded`; 5 `finding.addressed`; 8 `concern.recorded`/`concern.disposed`/`deviation.recorded`; 9 `decision.referenced`; 10 `human-outcome.recorded`; 12 `requirement.referenced`/`requirement.linked`. Nine keys, six tasks, no key unowned and none owned twice. |
| The two result files are named by exactly those six tasks | file-name search across all eighteen specs | PASS | `current_event_json.go` and `current_event_result_projection.go` each appear in tasks 4, 5, 8, 9, 10, 12 and nowhere else. `current_event_value_catalog.go` appears in those six plus 3, 6 and 11, which is correct: 3 and 6 own the `Claim` and Budget sub-items and 11 owns the `Condition` grammar under DEC-229. |
| Value shapes registered once and reused | task-spec read | PASS | Task 4 `Review`, `Coverage`, `CoverageItem`, `CoverageExclusion`, `Finding`, `FindingDisposition`; task 8 `Concern`, `ConcernDisposition`, `Deviation`; task 9 `DecisionRef`, `ProtectedDecision`; task 10 `HumanOutcome`; task 12 `RequirementRef`, `RequirementLink`. Fourteen shapes, each with one first owner. Tasks 5 and 10 are explicitly told to reuse and to escalate rather than edit a shape they do not own. |
| Task 4 owning `FindingDisposition` | `schema.md` 6.2 check | PASS | `review.recorded`'s accepted rule lists `finding_dispositions/insert` as an optional effect for `resolves_finding_ids`, so task 4 genuinely writes that table before task 5 exists. Correct, and not merely convenient. |
| Composite fill helpers assigned | task-spec read against `currentProjectedValue` | PASS | `Review.findings[]` / `Review.coverage` and `Finding.dispositions[]` at task 4, `Concern.disposition` at task 8, each named in the same style as the existing `currentFillRun` / `currentFillStopAttention`. These are the only Phase 23 shapes with child-table or nested content. |
| Gate reach for every task that gained a file | per-task command-block check | PASS | Tasks 3, 4, 5, 6, 8, 9, 10 and 12 all carry the byte-identical `go test -race -timeout 30m -coverprofile=coverage.out ./...`, and `internal/journal` is in `go list ./...` with in-package tests. Confirmed directly, not taken from the brief. |
| Task 1 unchanged | `git diff 5a94a42..HEAD` | PASS | `task-1.md` is not in the commit's file list. Task 1 gains nothing from R-16, as DEC-231 states. |
| The DEC-229 tasks-4/5 rejection still stands | `execution-scope.md` 6 read | PASS | Explicitly restated: tasks 4 and 5 gain the three R-16 files, not `current_event_catalog.go`, and the rejection is called out as unaffected. Correct - they are different files. |
| Contract untouched at 79 rows | row count plus diff inspection | PASS | 79 rows. The diff against `5a94a42` touches **zero** `VAL-` lines; only the header provenance block changed. |
| I-27 applied | contract header read | PASS | The header now records the 2026-09-15 and 2026-09-19 refinements separately, names VAL-006 (I-22) and VAL-034 (I-23) as the only two rows changed, states that DEC-230 and DEC-231 changed no row, and replaces the satisfied prerequisite with the current one. |
| I-28 applied | tasks 6 and 11 read | PASS | Both vacuous `internal/journal/` bullets are replaced with explicit "no catalog edit; a missing row is a task 1 defect - stop and escalate" instructions, which is stronger than the deletion I proposed. |
| I-29 applied | `schema.md` 6.2 read | PASS | `schema.md:617` now reads "rows added by task 9; also raised from this family by task 10 (I-29)". |
| Sub-item: `Claim.mode` at task 3 | reachability check | PASS | `currentProjectedValue` is called with kind `Claim` from the `run.recovered` and `run.abandoned` branches of `currentExpectedResult`, both live families, so the Tech Lead's stated reachability is correct. Task 3 is where review-mode Claims first exist. |
| Sub-item: Budget domains at task 6 | task-spec read | PASS | `unit` gains `review-cycles` and `exhaustion_action` gains `review-rework-exhausted` in both `Budget` and `ImmutableBudget` at task 6, which is the task that adds both. |

### 14.4 The Two Items The Tech Lead Recorded Unprompted

**5.1, the `idempotency_results` optional-stage note: CONFIRMED pre-existing,
with one figure corrected.** Measured in `currentEventRules`: **15** of 22
families encode `currentStage(0, ...)` for `idempotency_results` and **7**
encode `currentStage(1, ...)`, while `validateCurrentAppliedResults` requires
the result row unconditionally for all 22. So the distinction is real,
pre-existing and uniform in the sense that matters - the stage encoding never
governs whether the snapshot is required. "Nearly every family" overstates a
15/22 split; the note reads better as "most". All nine Phase 23 rows in
`schema.md` 6.2 list `idempotency_results/insert` as optional, mirroring the
majority encoding, and none is designed to omit the row. Not a gap; the note is
worth keeping.

**5.2, the `Subject` widening: NOT CONFIRMED.** This is R-18 below. The
conditional the Tech Lead removed was the wrong one to remove: four new shapes
do reference `Subject`, but the shape being widened is not free to widen.

### 14.5 Revisions Requested From The Tech Lead (fourth pass)

| ID | Sev | Where | Finding and proposal |
|---|---|---|---|
| R-17 | M | `internal/journal/current_event_cardinality.go`; task 4; `execution-scope.md` 6 | **A fourth closed registry on the same strict path has no owning task, and it rejects an ordinary two-Finding Review.** `validateCurrentCardinality` (`current_event_cardinality.go:8`) is called from the strict branch at `current_event_validation.go:96`, so every Phase 23 family reaches it the moment task 1 registers it in `currentEventRules`. It defaults `maximum := 1` per `table/op` key and widens that only through a closed `switch effect.Table` (`:14-33`) listing Phase 21/22 tables - `budget_samples`, `run_budget_bindings`, `referenced_artifacts`, `receipt_artifacts`, `activation_*`, and the Stop bundle tables. **None of the ten new Phase 23 tables is in it.** `public-interface.md` 4 accepts `findings:FindingInput[0..64]` and `resolves_finding_ids:ID[0..64]` on `review.record`, and `schema.md` 6.2 makes `review_findings/insert` and `finding_dispositions/insert` repeated optional effects of `review.recorded`, so the ordinary case is many effects on one key. **Executable proof** (`go test -overlay`, `internal/journal`, working tree untouched): a `review.recorded` Event with two `review_findings/insert` effects returns `current Event review.recorded: too many review_findings/insert Effects`, and two `finding_dispositions/insert` returns the same for that key. Controls behave correctly - one Finding passes; two `stops/insert` pass, because the existing `case "stops", "attention_items", ...` branch already widens the maximum for any family other than `stop.raised`/`stop.resolved`/`attention.resolved`; and `requirement_refs` insert-plus-update passes because they are distinct keys. A search of all eighteen task specs and every phase document returns **zero** mentions of `current_event_cardinality.go`. This lands at **task 4**, on the first Review that carries more than one Finding, with its cause outside task 4's declared scope - the identical failure mode to R-09, R-14 and R-16, one validator further along the same path. It blocks at minimum VAL-017 (verdict/Finding consistency across all three verdicts), VAL-020 (Findings created only with their Review in one atomic Event), VAL-022, VAL-023 and VAL-026, none of which is meaningful with a one-Finding ceiling. **Proposal:** apply the DEC-231 standing rule to this file as well - the task that unfences a family also registers that family's repeated-effect cardinality - and add `current_event_cardinality.go` to task 4's scope for `review_findings` and `finding_dispositions`. Check each remaining family as you disposition: `concern.*`, `deviation.recorded`, `decision.referenced` and `human-outcome.recorded` each write one row per Event and need nothing, and `requirement.referenced`'s insert and update are separate keys, so on my reading task 4 is the only task that needs an entry - but that derivation should be yours, not mine. No assertion changes; no public surface changes. |
| R-18 | M | task 4; `execution-scope.md` 6; DEC-231 | **The `Subject` widening drops an accepted Phase 22 constraint, and its recorded rationale is factually wrong.** `task-4.md` instructs "widen the `Subject` shape's `type:=run,envelope,activation` to the complete accepted SubjectType set", justified as "monotonic and cannot loosen anything reachable earlier because the fence keeps every other Phase 23 family closed". The fence is irrelevant here. `currentValueShapes["Subject"]` is referenced by exactly three shapes and no others - `Condition` (`:15`), `Stop` (`:26`) and `Attention` (`:27`) - and those belong to `stop.raised`, `stop.resolved` and `attention.resolved`, which are **Phase 21/22 families in `currentEventRules` that are never fenced** and are live from task 1. So the widening takes effect immediately at task 4 on live families, not on fenced ones. What it removes is explicit accepted material: `phase-22/public-interface.md:67-68` states "**Subject types for Stop/Attention are only run, envelope, activation**", and no Phase 23 document amends it - Phase 23's `public-interface.md` 2 amends `StopCode` and says the blocked-action subset "adds nothing", `stop.list --subject <run\|envelope\|activation>` is unamended, and a search of the Phase 23 documents for any Stop/Attention subject amendment returns nothing. The widening is also **unnecessary**: the four new shapes that need a wide subject (`Finding`, `Concern`, `Deviation`, `HumanOutcome`) are new, so they can reference a new shape instead. It is not a reachable defect today - `internal/workflow/stop.go`'s `validSubject` (`:377`) still restricts raising to run/envelope/activation and no Phase 23 task widens it - so this is a silent loss of a journal-layer invariant plus a wrong recorded reason, not a live bug; but a wrong reason in a task spec propagates, and the invariant would be gone with nothing recording that it was dropped deliberately. **Proposal:** leave `currentValueShapes["Subject"]` at `type:=run,envelope,activation`, and have task 4 register a **separate** wide shape - say `JudgmentSubject`, `type:=` the full accepted SubjectType set - referenced by `Finding`, `Concern`, `Deviation` and `HumanOutcome`. Same one-edit ownership at task 4, same monotonicity, and the Phase 22 Stop/Attention constraint survives intact. Correct the rationale in `task-4.md` and `execution-scope.md` 6 either way: whatever is decided, "the fence keeps every other Phase 23 family closed" is not a true statement about `Subject`. |

Improvements (the Tech Lead may apply or decline; none blocks Ready).

| ID | Where | Proposal |
|---|---|---|
| I-30 | `execution-scope.md` 6; DEC-231 | The `idempotency_results` note says "nearly every family". Measured: 15 of 22 use the optional stage and 7 use the required one. Say "most families" so the note does not invite a Developer to check for a uniformity that is not there. The substantive point - that the stage encoding never governs whether the snapshot is required - is correct and worth keeping. |
| I-31 | `execution-scope.md` 6 | The DEC-231 rule is stated for the three registries it names. Given that R-16 and R-17 were both found the same way, state it once as a general obligation - *a task that unfences a family owns that family's entry in every closed journal registry the strict path consults* - and list the registries known today (`current_event_json.go`, `current_event_result_projection.go`, `current_event_value_catalog.go`, `current_event_cardinality.go`). That makes the next such gap a Developer's escalation rather than a fifth planning pass. |

### 14.6 Blocker Candidates

**None.** R-17 assigns an existing file to a task that already owns the
behavior; the effect counts it needs are fixed by the already-accepted
`public-interface.md` 4 bounds and `schema.md` 6.2 rules. R-18 makes a task
spec's edit match accepted Phase 22 material instead of contradicting it, which
narrows a planned change rather than requesting a new one. Neither changes a
protected definition, command, threshold, budget, tolerance, assertion
semantics, accepted status, public surface, or policy. Both sit within Tech Lead
planning authority, as R-08, R-09, R-14 and R-16 did. No `ESC` identity is
requested and none is minted here.

### 14.7 Readiness Verdict

**NOT READY.**

**DEC-231's central judgment is confirmed correct.** The owner of the three
typed-result registries is the unfencing task, not task 1; the executability
test is the right test; the reasoning is the consistent extension of DEC-229
rather than a departure from it; and the full reachability sweep the brief asked
for - retained families, export and import, rebuild, replay - finds no path that
reaches any of the three before the fence opens. The mechanical re-checks all
pass: nine keys onto six tasks with none unowned or doubly owned, fourteen value
shapes each with one first owner, composite fill helpers assigned, gate reach
confirmed directly, task 1 untouched, the DEC-229 tasks-4/5 rejection intact,
the contract still at 79 rows with zero `VAL-` lines changed, and I-27, I-28 and
I-29 all applied. The unprompted `idempotency_results` note is confirmed
pre-existing.

What blocks Ready is the Tech Lead's own stated test in section 6 of the brief:
a fourth registry on the strict path with no owning task. **R-17**:
`validateCurrentCardinality` defaults every new Phase 23 table to one effect per
key, so `review.recorded` carrying two Findings - the ordinary case, and the one
VAL-017, VAL-020, VAL-022, VAL-023 and VAL-026 all depend on - is rejected at
append time at task 4, with its cause in a file no document names. Proven by
execution. **R-18**: the `Subject` widening would drop the accepted Phase 22
constraint that Stop and Attention subjects are only run, envelope and
activation, on live unfenced families, for a reason that does not hold; a
separate shape for the four new judgment values achieves the same end with the
invariant intact.

Both are bounded, decision-free, and of the same class the last three passes
have dispositioned. No assertion was weakened to reach this verdict, and no
`VAL-NNN` was narrowed, re-scoped, or restated to reach Ready.

## 15. Readiness Confirmation - DEC-232 (2026-09-19)

Readiness condition 3, fourth cycle. Confirmation of the edits applied in
`bf38e50` under DEC-232, per the superseding `tl-qa-confirmation-brief.md`.
Base: branch `phase-23` at `bf38e50`, `git status --porcelain` empty before and
after this pass. Sections 12, 13 and 14 stand; none is reopened. No product file
changed - all 14 edited files are under `project-docs/docs/` per
`git show --stat` - so no budget, migration, raising sequence, failure-code set
or contract count was re-measured.

### 15.1 The Three Tech Lead Findings

| Finding | Verdict | Independent source evidence |
|---|---|---|
| **R-17b `validateCurrentUpdate`** | **CONFIRMED, and the owner set is complete** | `current_event_updates.go:44` is `default: return fmt.Errorf("current Event update table not admitted")`, and the `switch effect.Table` at `:14` lists only `run_runtime`, `runs`, `claims`, `kit_activations`, `activation_components`, `activation_resources`, `activation_recoveries`, `attention_items`, `idempotency_results` - no Phase 23 table. It is called at `current_event_validation.go:200`, inside the strict path. I derived the affected set from `schema.md` 6.2 independently: of the nine families, exactly two emit an update to a Phase 23 table - `concern.disposed`'s **required** `judgment_concerns/update` and `requirement.referenced`'s **optional** `requirement_refs/update`. The `attention_items/update` rows in `review.recorded`, `decision.referenced` and `human-outcome.recorded` are already admitted. **No third Phase 23 update Effect exists.** Tasks 8 and 12 are the right owners and task 15's rename carry-through is correct. The field and guard sets are derivable as directed: `currentKeys` (`:52`) demands an exact sorted-set match of both `effect.Values` and `effect.Preconditions`, and `schema.md` 2 fixes the disposition columns on `judgment_concerns` and `status`/`superseded_by_id` on `requirement_refs`, so each task can derive its own sets from accepted DDL without loosening anything. The `validateCurrentTerminalValues` note is accurate - that second switch (`:61`) leaves `valid` true for any table it does not name, so it fails open and blocks nothing. |
| **R-17c the `validateCurrentJSON` column switch** | **CONFIRMED, and I agree with not reassigning it; the enumeration is complete, but the brief's count is one low** | `current_event_json.go:75` returns `unregistered current JSON column %s.%s` for any `*_json` column not in the switch at `:14`, with only `start_repository_json` and `terminal_repository_json` falling through. The judgment call is right: the file is already owned by all six unfencing tasks under DEC-231, the failure is loud and lands inside the owning task's own first Event, so describing the columns is sufficient and reassignment would add nothing. **Enumeration checked against `schema.md` 2 by extracting every `_json` column from the DDL.** Fourteen exist on Phase 23 tables. Twelve need registration and are all assigned: task 4 five (`review_records.repository_json`, `coverage_json`, `evidence_refs_json`, `receipt_ids_json`, `review_findings.scope_json`), task 8 three (`judgment_concerns.scope_json`, `judgment_concerns.references_json`, `deviations.references_json`), task 9 two (`decision_refs.source_json`, `authority_capabilities_json`), task 10 one (`human_outcomes.references_json`), task 12 one (`requirement_refs.source_json`). The other two - `envelope_protected_decisions.required_capabilities_json` and `affected_actions_json` - correctly need nothing, because their only writer `envelope.created` is retained. `evidence_records.target_repository_json` likewise needs nothing, for the same reason, and `finding_dispositions` and `requirement_issue_links` carry no `_json` column at all, so task 5's recorded negative holds. The brief's prose says "eleven columns"; the assignment is **twelve** and is complete. Prose only - see I-32. |
| **R-18b `schema.md` 6.2's I-17 sentence** | **CONFIRMED, and the correction is right** | Verified directly in source rather than from the prose. `currentEventRules` spans `current_event_catalog.go:26-49`; `retainedCurrentRules` starts at `:77`. All four amended families are inside the retained map: `envelope.created` `:82`, `envelope.activated` `:83`, `evidence.linked` `:88`, `artifact.availability-observed` `:89`. **None is in `currentEventRules`.** The earlier I-17 text, which placed the first row in `currentEventRules`, was wrong, and the consequence the Tech Lead names is real: a Developer following it would have moved a live family onto the strict path, where its own long-standing Effects would then face the six registries in the inventory. The corrected sentence at `schema.md:607-613` states the classification and the line range accurately. None of the R-17 conclusions moves, because they were derived from `schema.md` 6.2's effect rows rather than from this sentence. |

### 15.2 The `ErrorSubjectType` Reading

**CONFIRMED. The reading is correct, and I would have reached it independently
on two decisive tests.**

1. *"`requirement` is retained."* `requirement` is in Phase 21's documented
   `ErrorSubjectType` (`phase-21/public-interface.md:45-48`) and in the
   implementation (`internal/cli/codes.go` `SubjectRequirement`). It is **not**
   in the Stop/Attention `Subject` set, and there was no prior judgment subject
   vocabulary at all - Concerns, Deviations and Findings are new this phase. So
   "retained" can only be true of the error vocabulary. No other reading
   survives this clause.
2. *Phase 23 amends `ReasonSubjectType` separately.* `public-interface.md` 5
   says `ReasonSubjectType` adds `review`, `finding`, `decision`, and
   `human-outcome` - four values, a different set from section 2's six. The two
   rows therefore cannot name the same vocabulary.

The six added values are also exactly the six new entity kinds that can appear
as an error subject, which is coherent with the failure-code table. And the
Tech Lead's reductio holds: read literally as a Stop subject amendment, a
resolution condition could name `type:"path"`.

**`JudgmentSubject` adds no accepted value.** `public-interface.md:66-67`
defines it as `type:"issue"|"envelope"|"run"|"review"|"finding"|"requirement"`.
That is byte-identical to the closed set already accepted under DEC-222 in
`contracts.md` 8 ("`issue`, `envelope`, `run`, `review`, `finding`,
`requirement`") and to the `subject_type` CHECK on all four judgment tables in
`schema.md` 2, which I extracted and compared: `review_findings`,
`judgment_concerns`, `deviations` and `human_outcomes` each carry
`('issue','envelope','run','review','finding','requirement')`. It is a name for
an accepted domain, not a new one.

**The stronger conclusion is also right, and better than my R-18 proposal.** The
two domains are genuinely non-nested: `stops` and `attention_items` CHECK
`('run','envelope','activation')` (`0015:56,82`), so Stop admits `activation`
and judgment does not, while judgment admits `issue`, `review`, `finding` and
`requirement` and Stop does not. No single shape can serve both, so "no widening
is possible" is correct where I had proposed only that widening was unnecessary.
Withdrawing the `Subject` edit entirely is the right call.

### 15.3 Inventory Completeness - The Seventh Registry

I walked `validateCurrentEvent` and `validateCurrentAppliedResults` end to end
and enumerated every closed-set rejection in `internal/journal` rather than
checking the inventory's six against itself. The six listed registries are
correct and correctly owned, and the two recorded as needing nothing
(`validateCurrentBindings`'s `entityTable` map and `currentResultEntity`'s map)
do both fail open to the aggregate ID, as recorded. I also checked three the
inventory does not list and confirmed each needs nothing:
`validateCurrentStopReferences`'s subject map
(`current_event_stop_references.go:14`) fails closed but already admits
`envelope`, which every Phase 23 Stop uses; `current_event_prior.go:75,81` keys
on `currentEventTables`, which task 1 owns; and `validateCurrentValue`'s
`unregistered current value type` default is satisfied by the same
`currentValueShapes` registration DEC-231 already assigns.

**But there is a seventh, and it is a hard blocker: R-19 below.** It was missed
by both of us for the same structural reason - it is not reached through a
Phase 23 family at all. It is reached when a Phase 23 **value** flows through a
live Phase 21/22 family, which is outside the "walk the strict path for the nine
families" frame both previous sweeps used.

### 15.4 Per-Family Completeness

For each family, the registries it needs against the registries its owning task
names. Derived from `schema.md` 6.2 and the DDL, not from the brief.

| Family | Task | `fields` | `expectedResult` | value shapes | cardinality | JSON cols | updates | Complete? |
|---|---|---|---|---|---|---|---|---|
| `review.recorded` | 4 | yes | yes | `Review`, `Coverage`, `CoverageItem`, `CoverageExclusion`, `Finding`, `FindingDisposition`, `JudgmentSubject` | `review_findings`, `finding_dispositions` | 5 | - | **yes** |
| `finding.addressed` | 5 | yes | yes | reuse only | none needed | none | - | **yes** |
| `concern.recorded` | 8 | yes | yes | `Concern`, `ConcernDisposition` | none needed | 2 | - | **yes** |
| `concern.disposed` | 8 | yes | yes | reuse | none needed | - | `judgment_concerns` | **yes** |
| `deviation.recorded` | 8 | yes | yes | `Deviation` | none needed | 1 | - | **yes** |
| `decision.referenced` | 9 | yes | yes | `DecisionRef`, `ProtectedDecision` | none needed | 2 | - | **yes** |
| `human-outcome.recorded` | 10 | yes | yes | `HumanOutcome` | none needed | 1 | - | **yes** |
| `requirement.referenced` | 12 | yes | yes | `RequirementRef` | none needed | 1 | `requirement_refs` | **yes** |
| `requirement.linked` | 12 | yes | yes | `RequirementLink` | none needed | none | - | **yes** |

**Task 5's "needs nothing" holds.** `finding.addressed` writes one
`finding_dispositions` row per Event, so cardinality's default maximum of 1 is
correct; `finding_dispositions` carries no `*_json` column in the accepted DDL;
it emits no update; and its two value shapes are registered at task 4 because
`review.recorded` already inserts that table. Its only entries are the `fields`
row and the `currentExpectedResult` branch, exactly as recorded.

**R-17's owner derivation re-derived independently: task 4 alone is right.**
Only `review.recorded` carries repeated effects on a Phase 23 table -
`findings:FindingInput[0..64]` and `resolves_finding_ids:ID[0..64]`. Every other
family writes one row per key, and `requirement.referenced`'s insert and update
are distinct keys.

### 15.5 Mechanical Re-checks

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Contract untouched | row count plus diff | PASS | **79** rows; `git diff a2cef73..HEAD` on the contract file is **empty**. No row added, removed, renumbered, reworded, or weakened. |
| No product file touched | `git show --stat` | PASS | All 14 files under `project-docs/docs/`. |
| I-30 applied | `execution-scope.md` 6 read | PASS | Now reads "15 of the 22 existing families and a required stage for the other 7", matching my measurement. |
| I-31 applied and extended | `execution-scope.md` 6 read | PASS | The general rule is stated, the inventory lists six registries with file, line, rejection message and owners, two are recorded as deliberately needing nothing, and a registry not on the list is declared a Developer escalation. Grew from my four to six. See I-33 for the one scoping word I would change. |
| Task 12's end-of-phase audit widened | `task-12.md` read | PASS | Now requires every family to have its entry in **all six** inventory registries before the fence map is removed, not three. |
| Task 15 rename carry-through | `task-15.md` read | PASS | The `judgment_concerns` rename now crosses five catalogs, not three, which matches the registries that key on the table name. |
| `Subject` untouched by every task | file-scope search | PASS | No task instructs an edit to `currentValueShapes["Subject"]`; tasks 8 and 10 are explicitly told to reuse `JudgmentSubject` and not to edit `Subject`. |

### 15.6 Revisions Requested From The Tech Lead (fifth pass)

| ID | Sev | Where | Finding and proposal |
|---|---|---|---|
| R-19 | M | `internal/journal/current_event_value_relations.go`; tasks 4, 5, 8, 9, 10; `execution-scope.md` 6 | **The journal holds a second, independent copy of the closed 23-name capability allowlist, and no document names the file.** `validateCurrentValueRelations`'s `case "ActorProfile"` (`current_event_value_relations.go:59-66`) contains a hard-coded `const actions = "artifact.verify,attention.resolve,claim.acquire,...,stop.raise,stop.resolve"` - exactly the same 23 names as `sortedUniqueCapabilities` in `internal/workflow/activation.go:832`, which R-08 already assigned to tasks 4, 5, 8, 9 and 10 to take to 30 (`public-interface.md` 4.1, VAL-067). An Actor profile's `capabilities` are validated against that constant, plus a sorted-unique check. **Reachability, and why both previous sweeps missed it:** this is not on the Phase 23 families' strict path. It is reached through `activation.component-staged` - a live, never-fenced Phase 21/22 family - whose `activation_components.evidence_json` is validated as `[]ActorProfile` when `component_type` is `authority-profile` (`current_event_json.go:52-63`). `public-interface.md` 4 requires every Phase 23 leaf's capability to be present in the Actor profile, and that profile reaches the journal only through Kit Activation. So the first Activation that stages a Phase 23-capable profile is rejected at append time, before task 4 can record a single Review, with its cause in a file no task names. Both earlier sweeps walked the nine families' own path; this is a Phase 23 **value** travelling through a Phase 21/22 **family**, which that frame does not reach. **Executable proof** (`go test -overlay`, `internal/journal`, working tree untouched): a control profile with `claim.acquire` and `evidence.link` is accepted, and each of the seven grant actions - `review.record`, `finding.address`, `concern.record`, `concern.dispose`, `deviation.record`, `decision.record-reference`, `human-outcome.record` - is rejected with `invalid closed enum` naming the 23-name list. `requirement.record-reference` and `requirement.link` are also rejected, which is **correct** and serves as a second control: `public-interface.md` 4.1 deliberately keeps them out of the allowlist as operation IDs only, so the list is confirmed to be the capability allowlist and not something broader. **Proposal:** add `internal/journal/current_event_value_relations.go` to the file scope of tasks 4, 5, 8, 9 and 10 alongside the `internal/workflow/activation.go` entry R-08 already gives each of them, so one task adds one name to both copies in the same change; the journal list goes 23 -> 30 with exactly the seven grant actions and not the two requirement operations. Record it in the `execution-scope.md` 6 inventory as a seventh row, noting that it is the only entry reached through a Phase 21/22 family. Consider also whether task 18's 30-name audit should compare the two copies against each other, since nothing does today - that is the same silent-divergence risk R-09 was accepted to prevent for `semanticEventRules`. No assertion changes; no public value changes; the 30-name set is already accepted. |

Improvements (the Tech Lead may apply or decline; none blocks Ready).

| ID | Where | Proposal |
|---|---|---|
| I-32 | brief section 3, R-17c | The JSON column enumeration is described as "eleven columns" but the task specs correctly assign **twelve**: task 4 five, task 8 three, task 9 two, task 10 one, task 12 one. The assignment is complete; only the count in the prose is low. Correct it wherever it is restated so a later pass does not go looking for the missing one. |
| I-33 | `execution-scope.md` 6 | The inventory rule is scoped to registries "the strict path consults" for a Phase 23 family. R-19 is proof that the frame is one step too narrow: it sits on the value-relations path, reached when a Phase 23 **value** flows through a Phase 21/22 family. Restate the rule as *any closed journal registry that a Phase 23 value or Effect reaches, whichever family carries it*, and note the two known entry points - a Phase 23 family on the strict path, and a Phase 23 value inside a retained or live Phase 21/22 family. |
| I-34 | `public-interface.md` 2 | The renamed `ErrorSubjectType` row says "Retain the Phase 21/22 set", and cites `internal/cli/codes.go:88` as the implementation. The two have drifted: the documented Phase 21 set has 22 values including `actor`, `idempotency`, `obligation`, `path` and `repository`; the implemented `SubjectType` has 19 and adds `stop` and `attention`. The rename is still right and the instruction to add six is unambiguous and monotonic, but say which set is authoritative so a Developer reconciling them knows which way to move. |

### 15.7 Blocker Candidates

**None.** R-19 adds an existing file to five task scopes that already own the
identical change in the sibling file under R-08, and the 30-name set it must
reach is already accepted in `public-interface.md` 4.1 and asserted by VAL-067.
It changes no protected definition, command, threshold, budget, tolerance,
assertion semantics, accepted status, public value, or policy. Tech Lead
planning authority, as R-08, R-09, R-14, R-16, R-17 and R-18 were. No `ESC`
identity is requested and none is minted here.

### 15.8 Readiness Verdict

**NOT READY.**

**Everything DEC-232 did is confirmed correct.** All three Tech Lead findings
survive independent source verification: `validateCurrentUpdate` is a real
unowned hard-failure registry and tasks 8 and 12 are the complete owner set with
no third Phase 23 update Effect; the `validateCurrentJSON` column switch is
correctly described rather than reassigned, and its twelve columns are all
assigned; and the I-17 misclassification was real, with all four amended
families verified at `current_event_catalog.go:82,83,88,89` inside
`retainedCurrentRules`. The `ErrorSubjectType` reading is right on two
independent decisive tests, `JudgmentSubject` adds no accepted value, and the
stronger conclusion - that no widening of `Subject` is possible because the two
domains are not nested - is better than the R-18 remedy I proposed. Per-family
completeness closes for all nine families, task 5's recorded negative holds, and
R-17's owner derivation re-derives to task 4 alone. The contract is untouched at
79 rows with an empty diff, and no product file changed.

What blocks Ready is **R-19**, a seventh registry, and the honest prior in the
brief's own section 7 was right. `current_event_value_relations.go` holds a
byte-identical second copy of the 23-name capability allowlist that R-08 already
assigns to tasks 4, 5, 8, 9 and 10 in `internal/workflow/activation.go`, and no
document names the journal copy. All seven Phase 23 grant actions are rejected
by it, proven by execution with two passing controls. It is a hard blocker
before task 4's first Review, because a Phase 23 leaf's capability reaches the
journal only inside an Actor profile staged through Kit Activation.

It was missed twice because it is structurally different from the previous six:
it is reached by a Phase 23 **value** travelling through a live Phase 21/22
**family**, not by a Phase 23 family on the strict path. That is why I-33 asks
for the inventory rule to be widened rather than only extended by one row - the
frame, not just the list, was incomplete.

R-19 is bounded and decision-free. No assertion was weakened to reach this
verdict, and no `VAL-NNN` was narrowed, re-scoped, or restated to reach Ready.

## 17. Readiness Confirmation - DEC-233 And ESC-23-001 (2026-09-19)

Sixth confirmation pass, fresh session, against the Tech Lead brief for DEC-233
(`tl-qa-confirmation-brief.md`, rewritten for this pass). Base: branch
`phase-23` at `c27c1f4`, `git status --porcelain` empty before and after.
Sections 12 through 16 are not reopened; the contract is untouched at **79**
rows. Every measurement below was taken with `go test -overlay`, so no product
file was written at any point - re-verified clean after the last probe.

**Verdict: READY EXCEPT ESC-23-001.** I agree that tasks 1-8 are unaffected and
may dispatch. ESC-23-001 is verified in every particular, is *stronger* than
recorded in two respects, and its recommendation is right in direction but wrong
in one stated detail that would destroy the remedy's whole benefit if a
Developer implemented it literally - R-20 below.

### 17.1 The Four Asks, In One Line Each

| Ask | Result |
|---|---|
| 1. ESC-23-001 hard, and argue with Option A | **Confirmed, and extended.** Both consequences reproduced end to end, plus a third the escalation does not record. Option A is materially cheaper than B - measured - but its own sentence has two readings and only one works (R-20) |
| 2. Scope of the block | **Confirmed.** All eighteen specs checked: no task before 9 touches the `Condition` shape or struct, and task 6 explicitly forbids it. Tasks 9, 10, 11 held; 12-18 by ordering only, and task 12 is substantively independent |
| 3. R-19's application and the two other findings | **Confirmed in full.** Both allowlist copies verified at 23, the seven owners' names and the 30 total re-derived, five task specs and row 7 correct, task 18's audit is now a real parity assertion. Profile-14 exclusion and row 8 both verified by execution |
| 4. Inventory completeness on the widened frame | **A ninth item exists**, of the fails-open kind, like row 8: the journal's budget-fact validation does not reach any Phase 23 family. Recorded as R-21. It blocks no task. Also one validation **void** (R-22) and one new positive |

### 17.2 ESC-23-001 - Evidence Verification

Method for the replay arm: a real Phase 22 `stop.raised` Event was committed
through `journal.WithLock` / `PrepareEvent` / `CommitEvent` into a pristine
workspace by the **baseline** binary, and the resulting journal bytes were
recorded. Those same bytes were then replayed through `journal.ReplayJournal`
against a fresh pristine workspace under three different binaries. The recorded
line carries
`{"type":"work-ended","subject":{"type":"run","id":"RUN-1"},"budget_name":null,"capability_id":null}`.

| Claim | Method | Result | Evidence |
|---|---|---|---|
| `validateCurrentShape` requires an exact field count | source | **PASS - citation exact** | `current_event_value_types.go:112` is `if !ok || len(object) != len(fields)` |
| A pre-widening Condition is rejected by the widened shape | overlay execution | **PASS** | all three Phase 22 forms (`budget-observed`, `executor-terminated`, `capability-available`) return `Condition: closed object has wrong shape`; all three are accepted on the baseline |
| `ReplayJournal` revalidates historical Events against current shapes | **end-to-end execution**, not source reading | **PASS - stronger than recorded** | baseline replay of the recorded journal: `REBUILD OK: events=1 stops=1 attention=1`. Widened replay of the **same bytes**: `malformed journal entry at .fid/journal/2026-09-19.jsonl:1: project event "EVT-20260919-000002": current Event stop.raised stops/insert: invalid JSON column "resolution_condition_json": Condition: closed object has wrong shape`. Note the failure class: a `CorruptJournalError`, so the operator is told the journal is **malformed**, not that the binary is too new |
| There is genuinely no upgrade hook | source | **PASS** | `mustNewDefaultUpgradeChain` (`upgrade.go:97-118`) registers exactly two upgraders, `V1->V1` and `V2->V2`, both `func(event Event) (Event, error) { return event, nil }`. `CurrentEventSchema` is still `EventSchemaV1` and the chain's current is `EventSchemaV2`. Nothing in the chain can rewrite a payload |
| `attentionSemanticKey` hashes the whole struct | source + execution | **PASS** | `stop.go:422` marshals `map[string]any{"subject":…,"code":…,"resolution_condition":condition}` through `MarshalCanonicalJSON`, which round-trips `json.Marshal` and so preserves explicit nulls |
| Every stored `semantic_key` stops matching | measurement | **PASS** | three representative Conditions, baseline -> widened: `da2f2a11ead9bd44…` -> `d9858e43e38ffd7f…`; `38cda009a5f63497…` -> `5c48f0fa952948db…`; `28607d77c9f106d8…` -> `d263b7114f7965e3…`. The canonical form gains `"decision_class":null,"finding_id":null,"outcome_class":null` unconditionally, so the shift is **input-independent**: no Condition keeps its key |
| An identical Stop raises a duplicate Attention | **end-to-end execution across two binaries sharing one workspace** | **PASS** | baseline raise: `stop=STP-c6d8fd attention=ATT-e01a1c key=4ba154fa… openAttention=1`. Second identical raise under the baseline (control): `stop=STP-8a72e2 attention=ATT-e01a1c openAttention=1` - joins the open Attention. Second identical raise under the widened struct: `stop=STP-9286c4 attention=ATT-49846c key=a97a66c8… openAttention=**2**`, with both keys stored side by side |

Two corrections to the record, neither weakening the conclusion:

1. **The escalation names one strict reach point; there are four.** It cites
   `stops.resolution_condition_json` (`current_event_json.go:27`, not `:26`).
   `Condition` is also reached as a nested field of the `Stop` value, which
   `validateCurrentResultJSON` (`current_event_json.go:134-141`) requires in the
   typed result of `stop.raised` (`stop:Stop attention:Attention`),
   `stop.resolved` (`stop:Stop`), and `run.budget-observed`
   (`stops:[]Stop attention:[]Attention`). Measured: a Phase 22 `Stop` value
   carrying a four-field Condition is rejected under the widened shape with
   `Stop: resolution_condition: Condition: closed object has wrong shape`. This
   matters only for Option B, whose upgrader must therefore rewrite Conditions in
   **four** places rather than one.
2. **`openAttentionBySemanticKey` does not exist.** The function is
   `loadOpenAttentionByKey` at `stop.go:557` (`:558` is its `SELECT`), called
   from `stop.go:172`, `budget.go:253`, and `activation_recovery.go:334`. Its
   predicate is a literal `WHERE semantic_key = ? AND status = 'open'`, so the
   reasoning is right; the name in an escalation the Human will read should be
   the real one.

**A third consequence, not in the escalation, measured.** `StopRaiseRequest`
(`stop.go:72-79`) embeds `ResolutionCondition Condition`, and `stopDigest`
(`stop.go:431`) canonicalises the whole request into
`idempotency_results.request_digest`. Widening the struct therefore changes the
request digest of every `stop.raise`. Measured on a workspace whose Stop was
raised by the baseline binary: replaying the **same idempotency key** under the
widened struct returns
`idempotency conflict for stop.raise:project-39a67eeb55ab37a49660127379b9876d key "retained-test-1"`
instead of the stored result. A retry that was idempotent before the upgrade
becomes a hard refusal after it. This strengthens the escalation and it enlarges
Option B, whose migration must recompute `request_digest` as well as
`semantic_key`.

### 17.3 Arguing The Recommendation

**I agree with Option A, and I disagree with three things about how it is
stated.** Option A was implemented under overlay - per-type Condition shapes
keyed on `type` in the catalog, one dispatch branch, `omitempty` on the three new
fields - and measured.

| Measurement | Result |
|---|---|
| The recorded pre-widening journal replays under Option A | **`REBUILD OK: events=1 stops=1 attention=1`** - the identical bytes that fail under the plain widening |
| The three Phase 22 semantic keys under Option A | `da2f2a11ead9bd44…`, `38cda009a5f63497…`, `28607d77c9f106d8…` - **byte-identical to the baseline**. Both recorded consequences and the third one above are eliminated |
| Exclusivity becomes structural | wrong field -> `Condition: missing required field "outcome_class"`; own field null -> `Condition: outcome_class: invalid ID`; two fields -> `Condition: closed object has wrong shape`; own field only -> accepted. **All four cases inventory row 8 leaves open are closed by the shape itself** |
| Option A against the shipped corpus | `internal/journal`, `internal/workflow`, `internal/export` all `ok` (7.7 s / 32.8 s / 13.2 s), **zero failures, zero fixture edits** |
| The accepted present-and-null encoding against the shipped corpus | **11 failing tests and subtests** across 7 top-level tests, all Condition-carrying, before any Option B machinery is written: `TestExecutionBlockedEventProjectsIdenticallyLiveAndReplay`, `TestActivationRecoveryEventsProjectIdenticallyLiveAndReplay/activation.recovery-selected`, `TestCurrentEventRejectsConditionalBindingCorruption` (3 cases x live/replay), `TestTerminalSuccessEvaluationIncludesEnvelopeStop`, `TestRecoverRunRejectsApplicableStopWithoutEffects`, `TestStopAndAttentionListsFilterLimitAndRemainReadOnly`, `TestEquivalentStopsShareAttentionButKeepDistinctFacts`. `internal/export` is unaffected, consistent with positive finding 5 |
| Option A's size | 16 lines of catalog map plus a 15-line dispatch branch |

So on cost, Option A wins decisively, and by more than the escalation claims: it
also **removes** work, because the exclusivity rule inventory row 8 assigns to
tasks 9, 10, and 11 becomes structural and those three relation edits disappear.
Option B keeps row 8's work *and* adds a V3 bump, a four-site payload upgrader,
two column recomputes in migration 0017, a mid-life Attention key rule, and the
11 fixture repairs above.

**Three things about Option A are wrong or missing as written.**

1. **"Each new type carries only its own field" is the wrong reading of its own
   remedy, and it is the reading that breaks it.** With one Go struct there are
   two ways to read that sentence. (a) `omitempty` on the three *new* fields
   only: a new type marshals `{type, subject, budget_name:null,
   capability_id:null, <own field>}` - five fields - and the six Phase 22 types
   are untouched. (b) A new type carries *literally* only its own field, which
   requires `omitempty` on `budget_name` and `capability_id` too. Reading (b) was
   measured: the three keys become `ef2f1a6a1a938e0a…`, `cf865c69f70fcdd5…`,
   `9fb954bc97298c2a…` - **every one different from the baseline**, because
   `work-ended` drops from `{"budget_name":null,"capability_id":null,…}` to
   `{"subject":…,"type":"work-ended"}`. Reading (b) reproduces exactly the
   breakage Option A exists to prevent. Only reading (a) delivers the promised
   byte-identity, and every measurement in the table above used (a). The
   escalation must say (a) explicitly before the Human chooses, because a
   Developer given the current sentence can reasonably implement (b).
2. **Option A touches a file no task owns, and the escalation does not name
   it.** The analogy to `currentStopDetailShapes` is sound in spirit but differs
   in one way that costs a file: details shapes are selected at the call site
   from the **row's** `code` column (`current_event_json.go:35-40`), whereas a
   Condition's discriminant lives **inside the value**, and the value is also
   reached generically as a field of `Stop` through
   `if shape, ok := currentValueShapes[kind]; ok`
   (`current_event_value_types.go:38`). Per-type selection therefore requires a
   `Condition` branch in `validateCurrentValue` - that is, in
   `current_event_value_types.go`, which the inventory names only as a *reader*
   of `currentStopDetailShapes` and which no task's file scope lists. Under the
   phase's own ownership discipline that file must join the scope of whichever
   task implements the remedy.
3. **A naive implementation of that branch panics on adversarial input.**
   Selecting the shape requires reading `object["type"]` **before**
   `validateCurrentShape` has established it is a string. Measured: the baseline
   dispatch degrades gracefully -
   `Condition: type: invalid closed enum =executor-terminated,…` - while the
   obvious `object["type"].(string)` panics with
   `interface conversion: interface {} is int, not string`. This code validates
   untrusted journal and archive bytes during replay and `import`, so the lookup
   must use a checked assertion and return the ordinary rejection. One line of
   implementation, worth one sentence in the decision, because a panic here is a
   crash on a crafted archive.

**Option C's rejection is sound, and stronger than stated.** The escalation
rejects C because two human-outcome Attentions of different classes would
collide on one semantic key. Verified: the key is derived from
`(subject, code, condition)` only, so with no class field in the Condition two
Stops for different outcome classes on the same Envelope produce the same key,
`loadOpenAttentionByKey` joins them to one Attention, and resolving one class
closes the Attention covering both. The same collision applies a second time to
`decision-satisfied` and `decision_class`, which C does not mention. C is
correctly rejected.

**Is B ever preferable?** One honest argument for B: it keeps one shape per value
kind, which is the pattern every other entry in `currentValueShapes` follows, and
it preserves the accepted section 2 row verbatim. And one point in B's favour on
feasibility that the escalation does not claim: `ReplayJournal` hashes the
**stored** Event before upgrading it (`replay.go:141` then `:147`), so a
payload-rewriting upgrader does not break hash chains. B is implementable. It is
simply four to five times the work, in the two tasks already carrying the most
migration risk, for an encoding difference that no consumer has been shown to
need.

### 17.4 What The Harm Actually Is - The One Place I Think The Record Overstates

The escalation says "no workspace holding a pre-widening Stop can rebuild". The
mechanism is exactly as described; the **population** deserves a sentence,
because it changes what the Human is deciding.

`schema.md` section 1 states that task 1's edits to the shipped `0013` and `0015`
files are "a deliberate pre-adoption baseline replacement requiring fresh init
and rebuild; it is not an upgrade of old fid data". Migrations run under
`golang-migrate` (`internal/db/migrate.go:8`), which records only version and
dirty state and never re-runs or checksums an applied file, so an existing
workspace keeps its original `stops.code` and `work_claims.mode` CHECK text while
reporting schema 17. Under that accepted policy there is no released population
of Phase 22 workspaces expected to run Phase 23 features.

That does **not** dissolve ESC-23-001, for two reasons:

1. **The incompatibility is in-phase, not only post-release.** The widening lands
   at tasks 9-11, and task 6 raises `review-rework-exhausted` Stops carrying
   ordinary `budget-observed` Conditions. Any workspace - developer, gate, or
   fixture - whose journal was written by a tasks 1-8 binary and is rebuilt by a
   task 9+ binary hits the measured `CorruptJournalError`. That is precisely the
   TL's "VAL-004 and VAL-064 at risk at tasks 15/16", and it is a real risk
   rather than a theoretical one. One mitigating fact, checked: **no** pinned
   journal fixture in `testdata/` contains a `resolution_condition` today, so the
   exposure is live workspaces and any fixture task 6 chooses to pin.
2. **The Attention and idempotency consequences do not need a rebuild at all.**
   They are stored-column comparisons. Any workspace carrying an open Attention
   or a live `stop.raise` idempotency record across the upgrade is affected
   whether or not it ever replays.

So the decision the Human is making is better framed as "what encoding do we want
for V1, given that either choice costs us in-phase churn" than as "how do we
protect existing user data". Framed that way Option B's migration machinery is
work for a population the accepted policy says does not exist, while Option A's
cost is 31 lines and buys the row 8 invariant. That strengthens the
recommendation; it does not change it.

### 17.5 Scope Of The Block - Confirmed

All eighteen task specs searched for `Condition`, `resolution_condition`,
`semantic_key`, `attentionSemanticKey`, `stop.go`, and `currentValueShapes`.

| Task | Touches the `Condition` shape or struct? | Evidence |
|---|---|---|
| 1 | No | no match; its journal work is `current_event_columns.go`, `current_event_catalog.go`, the 0013/0015 CHECK text, and the `semanticPhase22Events` addition |
| 2 | No | no match |
| 3 | No | widens `currentValueShapes["Claim"]`'s `mode` enum only (`task-3.md:40`) - a closed-enum widening, monotonic, field count unchanged |
| 4 | No | adds `Finding`/`FindingDisposition`/`JudgmentSubject` shapes and asserts `Subject` byte-identical; `Condition` appears only as read context (`task-4.md:106,111,166`) |
| 5 | No | no match |
| **6** | **No - explicitly forbidden** | `task-6.md:72-74`: "**Add no Condition type here**: the resolution condition is the retained `budget-observed`." Its `current_event_value_catalog.go` edits are a new `currentStopDetailShapes` key and the `Budget`/`ImmutableBudget` enum widenings, all monotonic |
| 7 | No | no match |
| 8 | No | adds `Concern`/`Deviation` shapes, reuses `JudgmentSubject`, asserts `Subject` unchanged |
| 9, 10, 11 | **Yes** | each adds one field to the struct, to `currentValueShapes["Condition"]`, and to `validateCondition`, and each already carries "**Read ESC-23-001 before widening the shape at all**" |
| 12-18 | No | no match in any of the seven |

**I agree: tasks 1-8 are unaffected and the block is exactly tasks 9, 10, 11.**
Two notes. Task 6 co-edits `current_event_value_catalog.go`, the file Option A
would extend, but on unrelated lines - serial co-editing, not a conflict. And
tasks 12-18 contain no textual dependency on tasks 9-11 at all; task 12
(Requirement References) is substantively independent of the `Condition` shape,
so if the Human's decision takes time, task 12 is available out of order at the
Runner's discretion. Tasks 13-18 depend substantively on the features, not on the
shape.

Independently re-derived: every other section 2 amendment is monotonic. `Claim`,
`Budget`, `ImmutableBudget`, and `StopCode` are closed-enum or map-key additions
on live families; `Obligation`, `ObligationInput`, `Envelope`, `Evidence`, and
`Evaluation` land on `retainedCurrentRules` families or on tables absent from
`currentEventTables`, and the retained branch returns at
`current_event_validation.go:64` - verified as the exact line - before any value
validation; `ErrorSubjectType`, `Resume`, and `TrustSource` reach no journal
registry. **`Condition` is the only field-adding widening of a live journal value
shape in the phase.** Confirmed.

**New positive.** `ActivityEvent` is the only section 2 row that *removes*
names, which is the other non-monotonic direction and was not checked anywhere.
It is safe: `internal/activity/activity.go` holds no closed set and only reads
rows, `applyActivity` (`project_effects.go:49-62`) normalizes and inserts with no
vocabulary check, and `activity_log.event` is plain `TEXT NOT NULL` with no CHECK
(`migrations/0002_v01_tables.up.sql:172-181`). Removing the four V1 names at task
15 cannot break replay, and historical rows stay readable and filterable. Record
this as a positive so a later pass does not re-derive it.

### 17.6 R-19's Application - Confirmed

| Check | Result | Evidence |
|---|---|---|
| Both copies exist and are the same 23 names | PASS | `sortedUniqueCapabilities` (`activation.go:831-832`, a 23-key `map[string]bool`) and `const actions` (`current_event_value_relations.go:63`, a 23-member comma string). Parsed and compared: identical sets, both sorted |
| There is no **third** copy | PASS | `sortedUniqueEnvelopeActions` (`envelope.go:540-552`) delegates to `sortedUniqueCapabilities([]string{action})`; a repository-wide search for `run.force-recover` as a list member returns exactly the two copies. Row 7's "second copy" is exact |
| 23 -> 30 is the right split | PASS | `public-interface.md` 4.1 lists exactly seven added names across five owners: task 4 `review.record`; task 5 `finding.address`; task 8 `concern.record`, `concern.dispose`, `deviation.record`; task 9 `decision.record-reference`; task 10 `human-outcome.record`. 23 + 7 = 30, and `requirement.record-reference` / `requirement.link` are correctly excluded as operation IDs |
| All five task specs name the journal file | PASS | `task-4.md:63`, `task-5.md:58`, `task-8.md:76`, `task-9.md:63`, `task-10.md:76`, each citing `:63`, each naming its own action(s), each requiring "both copies sorted and identical" in one change, each stating the `activation.component-staged` reach |
| Inventory row 7 | PASS | `execution-scope.md` section 6 row 7 with the right file, message, and owner set, plus the narrative paragraph. `activation.component-staged` verified in `currentEventRules` (`current_event_catalog.go:28`) and behind no fence; `[]ActorProfile` selection verified at `current_event_json.go:52-63` |
| Task 18's VAL-067 audit | PASS - my follow-on applied as asked | `task-18.md` now requires the audit to "Assert set equality between the two copies directly, not each against a literal", with a matching acceptance criterion. This is deterministic structural conformance, which is the right evidence class for a parity assertion |
| Contract row VAL-067 unchanged | PASS | the contract file's last commit is `8cd613e`, two dispositions ago; no VAL row was touched to accommodate this |

### 17.7 The Tech Lead's Two Other Findings - Confirmed By Execution

**Finding 1, the missing profile-14 exclusion.** The boolean reading is exactly
right. `validateEvent` (`archive_semantics.go:323`) rejects on
`!found || (validator.profile == 14 && semanticPhase22Events[event.EventType])`
at `:328` - both citations exact. Once task 1 mirrors the nine families into
`semanticEventRules`, `found` is true at every profile, and the profile-14 guard
only excludes the 14 names currently in `semanticPhase22Events`
(`archive_semantics_catalog.go:85`), so a legacy projection-14 archive carrying
`review.recorded` would pass this check. The claim that VAL-065 is unaffected is
also right, and for a sharper reason than stated: the legacy branch at
`archive_validation.go:127-129` *deliberately admits* `ProjectionSchema == 14`
archives, so the manifest check at `:130` never sees them - which is precisely
why `semanticPhase22Events` has to do the work. Task 1 is the right owner and the
addition needs no Event. **PASS.** One gap, I-35 below: the addition has no
acceptance criterion and no VAL row.

**Finding 2, the `Condition` relation fails open.** Reproduced exactly, with the
shape widened under overlay and the tree untouched:

| Condition | Outcome |
|---|---|
| `human-outcome-recorded` with all three new fields non-null | **ACCEPTED** |
| `human-outcome-recorded` with none of them non-null | **ACCEPTED** |
| `decision-satisfied` carrying `outcome_class` instead of `decision_class` | **ACCEPTED** |
| Phase 22 control: `budget-observed` with a null `budget_name` | REJECTED `condition nullable binding mismatch` |

Row 8 is real, it fails open, it blocks nothing, and it is the correct class of
finding. Its contract coverage exists and is correctly placed: **VAL-078**
already requires "closed-schema rejection of a Condition carrying a second
non-null field or none" with the fail criterion "a Condition with two or zero
non-null fields is stored", owned by task 11. Note for the disposition of
ESC-23-001: if the Human chooses Option A, row 8's assigned work disappears
because the shape enforces it, and VAL-078's evidence becomes shape rejection
rather than relation rejection. The row should then be restated, not deleted.

### 17.8 Inventory Completeness On The Widened Frame

I walked the frame's second entry point the way the brief asked: for every value
in `public-interface.md` section 2, which live Phase 21/22 family carries it, and
what does that family's path validate.

| Section 2 value | Carrier | What validates it | Owner / result |
|---|---|---|---|
| `ErrorSubjectType` | none - CLI failure envelope | `internal/cli/codes.go` only | no journal reach. 19 implemented constants verified at `codes.go:89-108`; +6 = 25 |
| `ObligationInput`, `Obligation`, `Envelope` | `envelope.created` / `envelope.activated` (**retained**) | retained branch returns at `current_event_validation.go:64`; `work_envelopes` and `envelope_obligations` are absent from `currentEventTables` | no journal registry |
| `Evidence`, `Evaluation` | `evidence.linked` (**retained**); evaluations are a read model | same; `evidence_records` absent from `currentEventTables` | no journal registry |
| `Claim` | `run.recovered`, `run.abandoned`, `claim.expired` (**live**) | `currentValueShapes["Claim"]` | task 3 - enum widening |
| `Budget`, `ImmutableBudget` | `run.started`, `run.budget-observed` (**live**) | `currentValueShapes` + `run_budget_bindings.budget_json` as `ImmutableBudget` | task 6 - enum widening |
| `Condition` | `stop.raised`, `stop.resolved`, `attention.resolved`, `run.budget-observed`, `activation.recovery-*` (**live**) | `currentValueShapes["Condition"]` + row 8 relation, via **four** reach points (17.2) | tasks 9/10/11 - **ESC-23-001** |
| `StopCode`, StopDetails | same live families | `currentStopDetailShapes`, which `StopCode` itself keys on (`current_event_value_types.go:89`) | tasks 6/9/10/11 - map-key additions |
| Stop blocked actions | same | `validateCurrentBlockedActions` (`current_event_json.go:117`) closes six names; both names the new codes use are already members; the `budget-` prefix relation is at `:77` | no entry needed - positive holds |
| `Resume` | none - CLI output | - | no journal reach |
| `TrustSource` | `Receipt.trust_source:=fid-executed` is the only journal enum, and no Phase 23 receipt carries a new source | - | no entry needed |
| `ActivityEvent` | `activity_log` writes on retained issue/link families | nothing - free text, no CHECK | no entry needed; **new positive, 17.5** |

**R-21, the ninth item.** It is the fails-open kind, like row 8.
`validateCurrentBudgetFacts` (`current_event_budget_values.go:34-50`) switches on
`event.EventType` over exactly `run.budget-observed`, `run.execution-blocked`,
`run.recovered`, `gate.execution-started`, `gate.receipt-recorded`, and returns
`nil` at its `default`. Task 6 makes a **Phase 23 value** - a `review-cycles`
budget sample - reachable inside `review.recorded`, a family that will be in
`currentEventRules` after task 1, so `validateCurrentAppliedResults` will call
this function for it. Measured, with a `review-cycles` binding seeded and an
Event carrying a `budget_samples` insert:

```
BUDGET-PROBE review.recorded            -> returned nil - budget facts NOT validated
BUDGET-PROBE finding.addressed          -> returned nil - budget facts NOT validated
... all nine Phase 23 families identical ...
BUDGET-PROBE run.budget-observed        -> reached validation: current Event: missing reference to runs map[id:RUN-1]
```

So the journal-side sample/state consistency, automatic-charge rule
(`validateCurrentAutomaticCharge`, `:91-93`, itself scoped to three families),
and exhaustion-Stop derivation that Phase 22 enforces for `attempts` and
`interventions` do **not** reach Phase 23's `review-cycles` charge. VAL-027 and
VAL-029 are therefore producer-side assertions only. This blocks nothing and is
not a dispatch blocker - but the inventory's own closing rule says a registry
that "silently fails to constrain a Phase 23 value" belongs on the list, and a
Developer at task 6 reading the inventory today finds no row for it. Requested in
R-21.

**R-22, a validation void rather than a missing registry.** All four DEC-224
raising sequences were checked. Three of them - `envelope.activated`,
`evidence.linked`, `artifact.availability-observed` - are in
`retainedCurrentRules`, and both journal checks exit early for retained
families: `validateCurrentEvent` returns at `current_event_validation.go:64`, and
`validateCurrentAppliedResults` returns at
`current_event_result_projection.go:11-13`. Measured, in two arms:

| Arm | Result |
|---|---|
| Today's catalog | all three reject the readiness Stop row: `current Event <family>: unrelated domain Effect` - which is exactly why DEC-224 assigns the effect rows |
| With DEC-224's rows added (tasks 9/10's change) | **all three ACCEPT** a `human-outcome-required` Stop whose `resolution_condition_json` is `{"type":"human-outcome-recorded","subject":{…},"nonsense":true}` and whose `details_json` is `{"totally":"wrong"}` |

Only the `stops.code` CHECK and `json_valid()` constrain such a row. Two
consequences worth recording. The journal-side grammar tasks 9 and 10 add is
**not exercised by the raising path they own** - it first binds when the Stop is
read back through a live family, and task 9's `protected-decision-required` Stop
is raised only at `envelope.activated` (VAL-041), so it has no strict raising
path at all. That undercuts the inventory's ownership rationale - "the task's own
first successful Event proves them together" - for exactly these rows, and it
means a Developer could omit the `currentStopDetailShapes` entry, watch the raise
tests pass, and fail later at task 11 or at tasks 15/16, with the cause outside
the failing task. The remedy is documentation, not code: name the proving Event.

Per-value and per-family completeness otherwise closes, and all six DEC-233
positives were re-derived and hold, with the two anchor corrections in I-36.

### 17.9 Mechanical Re-checks

| Check | Result |
|---|---|
| Contract row count | **PASS** - 79 rows, 79 distinct `VAL-NNN`, no duplicates |
| No VAL renumbered or reworded | **PASS** - `validation-contract.md` has not been committed since `8cd613e`, two dispositions before this one |
| No product file touched | **PASS** - all 15 files in `c27c1f4` are under `project-docs/docs/`; `git diff --name-only HEAD~3 HEAD` outside `project-docs/` is empty |
| Working tree | **PASS** - `git status --porcelain` empty before and after every overlay probe; no `zz_qa_*` file exists in the tree |
| The one `public-interface.md` edit adds and removes no name | **PASS** - one row replaced, `ErrorSubjectType`; the six added names are unchanged; the edit only re-seats authority on the implementation and records the 19 -> 25 count. 19 implemented `SubjectType` constants counted at `codes.go:89-108`, and `codes.go:88` is the `type SubjectType string` declaration the row cites |
| I-32 applied | **PASS** - twelve columns in DEC-232's R-17c paragraph, assignment unchanged |
| I-34 applied | **PASS** - implementation authoritative, prose set recorded as drifted, reconciliation a wrap-up spec-sync item |
| No protected definition, command, threshold, or budget changed | **PASS** |
| Escalation mirrored | **PASS** - `tl_escalations_current.md` and `phase_brief.md` both carry ESC-23-001 with the same scope and recommendation as `tl-escalations.md` |

### 17.10 Revisions Requested From The Tech Lead (sixth pass)

| ID | Severity | Subject | Request |
|---|---|---|---|
| R-20 | M | ESC-23-001, Option A and Option B | Three amendments to the escalation before the Human reads it, all measured in 17.3: (i) state Option A's encoding unambiguously as `omitempty` on the **three new fields only**, with the six Phase 22 types keeping `budget_name` and `capability_id` present-and-null, and say that the alternative reading changes every existing semantic key; (ii) name `internal/journal/current_event_value_types.go` as the file the per-type dispatch requires and say it must join the implementing task's scope, and require a **checked** type assertion because the naive one panics on a non-string discriminant; (iii) correct Option B's cost to four payload sites plus an `idempotency_results.request_digest` recompute plus the 11 named corpus failures. Also fix two citations: the function is `loadOpenAttentionByKey` at `stop.go:557`, and the Condition column case is `current_event_json.go:27` |
| R-21 | M | `execution-scope.md` section 6 inventory; task 6 | Record the ninth item. `validateCurrentBudgetFacts` (`current_event_budget_values.go:34-50`) returns `nil` at its `default` for all nine Phase 23 families, so the `review-cycles` charge task 6 introduces has no journal-side validation of sample/state consistency, automatic charge, or exhaustion-Stop derivation - measured. Either add it as row 9 with task 6 as owner, or record it as an explicit positive stating that Phase 23 charges are deliberately producer-side only, and say so in task 6 so VAL-027 and VAL-029 evidence is not mistaken for journal-layer proof. Not a dispatch blocker: it fails open |
| R-22 | M | `execution-scope.md` section 6; tasks 9 and 10 | Record the validation void measured in 17.8 and name the proving Event. Three of DEC-224's four raising families are retained, so the `currentStopDetailShapes` and `Condition` entries tasks 9 and 10 add are not exercised by the raise they own; task 9's `protected-decision-required` has no strict raising path at all. State in both task specs which Event proves those entries, so the entry is not dropped as unnecessary when the raise tests pass green without it. Not a dispatch blocker |
| R-23 | M | `task-11.md` | VAL-078 is listed in task 11's Validation Contract Coverage and assigned to task 11 by the contract, but **no acceptance criterion cites it** - the criteria cover VAL-047, 048, 049, 050 only. Add the missing criterion, including the closed-schema half (a Condition with two or zero non-null fields is rejected with no Event), phrased to match whichever encoding the Human selects. Bounded and mechanical; task 11 is held by ESC-23-001 in any case |

### 17.11 Improvements (not blocking)

| ID | Subject | Note |
|---|---|---|
| I-35 | `task-1.md` | The `semanticPhase22Events` addition has no acceptance criterion and no VAL row, and VAL-065 correctly does not cover it (the legacy branch admits profile 14 by design). Add a criterion with no new VAL ID - "a crafted legacy projection-14 archive carrying `review.recorded` is rejected" - so the defensive addition is asserted rather than assumed |
| I-36 | anchor drift | Five citations are off by one or wrong, all in task-facing text: inventory row 8 says `current_event_value_relations.go:48`, actual `:49`; positive 2 says `:76`, actual `:77`; positive 3 says `current_event_value_types.go:88`, actual `:89`; ESC-23-001 says `current_event_json.go:26`, actual `:27`; ESC-23-001 names `openAttentionBySemanticKey`, actual `loadOpenAttentionByKey` at `stop.go:557`. Exact anchors verified correct: `current_event_value_types.go:112`, `current_event_value_relations.go:63`, `stop.go:422`, `replay.go:153`, `project_effects.go:11,20`, `upgrade.go:50`, `current_event_validation.go:26,64`, `archive_semantics.go:328`, `archive_semantics_catalog.go:85`, `archive_validation.go:130`, `current_event_value_catalog.go:6`, `activation.go:832` |
| I-37 | `schema.md` section 1 boundary | Under `golang-migrate` an applied migration is never re-run or checksummed, so a workspace migrated in place to 17 keeps the original seven-member `stops.code` CHECK and `mode = 'mutation'` while reporting schema 17, and a fresh init at 17 does not. The accepted policy ("fresh init, not an upgrade of old fid data") makes that intentional, and DEC-223 is not reopened - but the divergence is silent, and the first Phase 23 Stop on such a workspace fails with a CHECK-constraint error rather than a re-init instruction. Worth one recorded sentence, and it bounds ESC-23-001's harm population (17.4) |
| I-38 | `task-10.md:119` | "so tasks 11 extend a mapping" - grammar slip in text a Developer will read |

### 17.12 Blocker Candidates

One, and it is not mine to allocate: **ESC-23-001**, already raised by the Tech
Lead and awaiting the Human. Category: accepted public-value encoding with
durable-replay, Attention-identity, and idempotency-identity consequences.
Evidence: 17.2 and 17.3 above, all reproduced. Authority boundary: the remedy
changes what an accepted Phase 23 value looks like on the wire, which is Human
authority, not Tech Lead planning authority. Requested action: choose A, B, or
another route, with R-20's three amendments applied to the options first so the
choice is made on correct costs.

I raise **no new blocker candidate**. R-20 through R-23 are bounded and
decision-free, and none of them holds a task that ESC-23-001 does not already
hold: R-20 is an amendment to a record the Human has not yet acted on, R-21
concerns a registry that fails open, R-22 is documentation for two already-held
tasks, and R-23 is one missing checkbox in a held task.

### 17.13 Readiness Verdict

**READY EXCEPT ESC-23-001.**

**Tasks 1-8 are unaffected and may dispatch.** I checked all eighteen specs, not
only the three the escalation names: no task before 9 touches the `Condition`
shape or struct, and task 6 - the one that comes closest, because it owns
`internal/workflow/stop.go`'s closed registries and edits the same catalog file -
explicitly forbids adding a Condition type and uses the retained
`budget-observed` condition. Tasks 12-18 are held by ordering alone, and task 12
is substantively independent if the Runner wants work available while the Human
decides.

**Everything DEC-233 did is confirmed correct.** R-19's application is right in
both copies, in all five task specs, in inventory row 7, and in task 18's audit,
which is now a genuine two-copy parity assertion rather than two comparisons
against a literal. The `semanticPhase22Events` boolean reading is right and task
1 is the right owner. Row 8 is real, fails open, and already has contract
coverage in VAL-078. I-32 and I-34 are applied as accepted, the contract is
untouched at 79 rows, and no product file changed.

**ESC-23-001 is correctly raised and I would not have dispositioned it either.**
Its two recorded consequences reproduce end to end - the same recorded journal
bytes rebuild clean on the baseline and fail as a `CorruptJournalError` under the
widening, and the same identical Stop joins one Attention on the baseline and
opens a second one under the widening. There is a third consequence the record
does not carry: the `stop.raise` idempotency request digest also shifts, turning
an idempotent retry into a hard conflict. And there is a fourth reach point the
record does not carry: `Condition` is validated again inside the `Stop` value of
three families' typed results, which is where Option B's real cost lives.

**On the recommendation: Option A, with its statement corrected.** I implemented
it and measured it. It replays the pre-widening journal clean, it leaves all
three measured semantic keys byte-identical, it closes all four of inventory row
8's open cases structurally - which *removes* work from tasks 9, 10, and 11 - and
it passes `internal/journal`, `internal/workflow`, and `internal/export` with
zero failures and zero fixture edits, against 11 failing tests for the accepted
present-and-null encoding. That is not a marginal difference. But the sentence
"each new type carries only its own field" has a second reading that requires
`omitempty` on `budget_name` and `capability_id`, and I measured that reading
changing every existing key - it reproduces the exact harm Option A exists to
prevent. Option A is also one file larger than stated and carries a panic hazard
in the obvious implementation. Those three corrections are R-20, and they matter
more than the choice itself, because a correctly-chosen option implemented from
an ambiguous sentence lands in the same place as the wrong choice.

The four findings this pass adds - a ninth fails-open registry, a retained-family
validation void, one missing acceptance criterion in task 11, and one unasserted
addition in task 1 - are all bounded, none requires a decision, and none of them
holds a task that is not already held. No assertion was weakened, no `VAL-NNN`
was narrowed, re-scoped, or restated, and no product code was read as a
substitute for running it.
