# Decisions - Active

> Current commitments and carry-forwards. Historical reasoning stays in ADRs
> and phase decisions; standing delivery rules also live in decisions_principles.md.

## Durable Architecture

- ADR-0001: journal canonical; SQLite is a rebuildable projection.
- ADR-0004/0005: project/kit own intent and policy; fid owns operational facts;
  Issues are the only durable work graph.
- ADR-0006/0007: parallel Builders need distinct Issues/worktrees/disjoint scope;
  interruption recovery uses a new linked Run and non-destructive reconciliation.
- ADR-0010/0011: generic activated Actor/Envelope authority replaces global role
  policy; no pre-dogfood compatibility surface or old fid-data migration;
  JSON is the sole machine format; optional products remain evidence-gated.
- ADR-0012: six vertical scenario phases, 20..25, with explicit removal owners.
  All ADRs live under `../../../docs/adr/`.

## Ratified Product Commitments

- Gates A/B1: fid is the durable operational backend; kit compiles approved
  policy into immutable Envelope Revisions and fid enforces generic constraints.
- Gates B2/B3: immutable activated Gate Registrations drive fid execution.
  Receipts come from observed termination; reported results cannot satisfy proof.
- Gates B4/B4a: Issues are durable work units; exact Issue/Envelope/Actor/
  repository/optional-Obligation context binds Claims and Runs; one mutating
  Claim per Issue and isolated disjoint parallel work.
- Gates B4b/B4c: interrupted Runs never reopen; recovery records liveness and
  reconciliation. Builder rework stays on the same Issue; independent Review
  alone resolves Findings.
- Gate B4d: progress is a projection; material contract/selection changes create
  a revision; removing unfinished work needs an authorized disposition.
- Gates C/D1/D2: all twelve port capabilities qualify before live adoption;
  no global role-policy model; no duplicate process aggregate or generic format.
- Gates D3.1/D3.2: vertical delivery; Phase 21 delivered the one-Issue skeleton;
  Phase 22 delivered failure/recovery and its qualified capabilities (DEC-220).
- DEC-156: qualifying Evidence guards terminal-success Issue transition;
  Envelope completion reuses the same evaluator after the Issue is terminal.
- DEC-157/177: Phase 21 delivered wire authority is
  `../phases/phase-21/public-interface.md`; its removal inventory is
  `v1-convergence-inventory.md`. Public enums are kebab-case; bootstrap has no role.

## Standing Delivery Constraints

- DEC-238/239 recoveries are complete. Their earlier task holds are historical;
  see [phase progress](../phases/phase-23/progress.md). Keep the 30-second rebuild,
  600-second ordinary-package and 900-second tagged budgets. DEC-239 preserves
  Run/Claim Actor binding, durable history, real CLI Claim rejection and
  same-Reviewer rework while using an isolated Run-holder history query.

- DEC-158: report lint is a two-argument Runner pre-dispatch control, not a
  Developer task command; 100k rebuild needs three fresh <30 s samples (CF-024).
- DEC-162/165/166: retained-test convergence follows the accepted invalidating
  change; establish introducing commit provenance, never filename/task labels.
- DEC-167/173: a task's required gates execute every owned scope, and its scope
  covers every path its gates walk. Amend incomplete assignments before dispatch.
- DEC-175: a latent implementation defect first observable under a new gate
  belongs to that task with minimum causal source/regression scope.
- DEC-189: a nondeterministic retained test that first fails under an existing
  required gate is neither convergence nor DEC-175. Assign it to the observing
  task with a scope amendment, repair the nondeterminism at its source, and
  never re-baseline a golden to a load- or instrumentation-dependent value.
- DEC-170/171: external interruptions do not consume no-progress budget; polling
  timeout means wait. Only declared command budgets apply; preserve real outcomes.
- DEC-172: an unavailable capability row is required negative metadata, not an
  implementation of that capability. Replace only when owning semantics qualify.
- DEC-178: shipped migration SQL and canonical slice DDL must agree exactly.
- DEC-002: autonomous TL-Lite resume remains enabled under the current workflow.
- DEC-190: separate Human-authorized verifier maintenance resolves ESC-22-001;
  ordinary pending-triage correction and normal restart retain task/attempt
  authority, and active owning-review triggers work across recovery boundaries.
- DEC-191: canonical prompts retain the committed Codex model choices and
  the TL Bash capability; generated provider copies are synchronized.
- DEC-192: a progressive availability fence is removed key-by-key by the task
  owning the fenced Event's table; the last owner removes the empty fence and
  guard. This implements the contract; displaced assertions are convergence.
- DEC-193: a calendar-dependent retained test is DEC-189 nondeterminism, not
  convergence and not DEC-175. Assign it to the observing task, repair the date
  binding at its source with recorded introducing-commit provenance, and never
  pin an absolute date on an Event appended through the journal writer into a
  live-clock workspace.
- DEC-194: an accepted allocation that exists but was never delivered is bounded
  repair, not an ambiguous owner. Read the task spec, schema owner column,
  public-interface task column, and contract subjects before classifying.
- DEC-195: new closed-domain members require matching schema rows. Amend missing
  rows, retain assertions, and assign the first task that needs them, recording
  the original owner's provenance. Never substitute an admitted value.
- Current V1 delivery scaffolding is binding engineering procedure, not forward
  fid product scope. No failed gate authorizes a protected definition change.

## Current Phase

- DEC-222 (ACCEPTED by Human, 2026-09-14): Phase 23 "Review and Human judgment"
  is Heavy: accepted contracts, 84-assertion validation contract, and 18-task pack
  are indexed in [phase brief](phase_brief.md). Owning tasks apply the six
  accepted canonical spec amendments.
- DEC-222 Human planning decisions: Concern dispositions are exactly
  `carried-forward`, `rejected`, `superseded`, with self-disposal permitted and
  both Actors recorded; `links` is restricted to Issue endpoints so
  `requirement_issue_links` is the single traceability model, retaining the new
  `verifies` relation.
- DEC-223..227 (ACCEPTED by Human 2026-09-19): the link-preserving
  `evidence_records` rebuild, since 0017 fails wherever Evidence has Artifact
  links; four named Events raise `human-outcome-required`, two gaining effect
  rows; independence conflicts with mutating Claims and their Runs only, so one
  Reviewer may review successive rework cycles; named failure-code producers,
  denial precedence and non-committing denials; and an explicit `-timeout 30m`
  plus measured wall-time budgets.
- DEC-226 vocabulary (standing): a **Failure code** names one refused request
  and carries category, retryability, and exit status; an **Evaluation reason**
  says why one Obligation is unsatisfied and carries none of those. File each
  identifier in its glossary-defined set.
- DEC-228 (ACCEPTED, Tech Lead) carries the remaining dispositions, including
  three DEC-167/173 scope corrections and five new assertions.
- DEC-235 (ACCEPTED by Human 2026-09-20) resolves ESC-23-001 with Option A: the
  three new `Condition` fields are `omitempty` and absent on the types they do
  not apply to, the six Phase 22 types stay byte-identical, and exclusivity is
  structural in per-type journal shapes keyed on `type`.
- DEC-236: task 9's `case "Condition"` must call
  `validateCurrentValueRelations("Condition", object)` to retain nullable binding.
  `review-recorded` has no discriminator; `finding_id` is for `finding-resolved`.
- DEC-237: resolution (VAL-078) and Condition encoding (VAL-080) are separate.
  VAL-081 preserves nullable binding; VAL-082 preserves retained-type bytes.
- DEC-240: VAL-082 splits into VAL-082/083/084. Current count is **84**:
  80 deterministic, 4 Human.
- DEC-242: Task 5 adds optional reasoned `inapplicable_findings` to
  `review.record`; addressed is nonterminal, terminal dispositions belong to
  the qualifying Review Event.
- DEC-243: VAL-025 keeps all nine compatibility cells; Task 6 delivers the
  receipt/review links and Task 10 the real Human link and final closure.
- DEC-249 (Tech Lead 2026-09-26): an Obligation blocks only the guarded
  actions it declares. "Applicable" means `required` and guarding that action.
  Task 2 owns the evaluator conformance repair.
- DEC-250 (Tech Lead 2026-09-26): `human-outcome-recorded` and
  `decision-satisfied` qualifying facts. Resolution is keyed on the Condition,
  not the raising code, with one predicate shared by the manual and atomic
  paths. Task 11 implements it; DEC-253 ratifies its `stop_conditions.go`
  scope.
- DEC-251 (Human, Option B, 2026-09-26): `human-outcome.record` needs no
  Envelope grant. Human authority comes from the Activation profile. The
  section 4.1 cell is corrected.
- DEC-252..257 closed ESC-23-004's repairs (tasks 6, 2, 10, 11 complete);
  DEC-254 Part 1 (`cmd/fid` breach attribution) still stands.
- DEC-258 (Human-delegated): parallel tests, `paralleltest`, tiered test loop,
  80% headroom warning. DEC-259 (Human): both race gates carry
  `-gcflags=modernc.org/...=-d=checkptr=0`, plus a phase-end full-`checkptr` run.
- DEC-260 (Tech Lead, Human-delegated) resolves ESC-23-006. Requirement writes
  need a trusted Actor but no Actor-profile capability or grant. The closed
  capability set stays at thirty. Report-lint grammar is unchanged; task 12's
  declaration is corrected to one line.
- DEC-261 routes review-9's repairs. Reopened tasks run the DEC-259 race gates.
- DEC-262 (Tech Lead, Human-delegated) resolves ESC-23-007. Task 13 owns the
  strict-path `activity_log` registry entry (execution-scope row 11) and the
  `issue`-subject ActivityEvent rows. Task 15 writes the
  `requirement`/`decision`-subject rows with 0018 and closes VAL-055. The
  subjects are in `schema.md` section 6.2.
- DEC-263 (Tech Lead, Human-delegated) resolves ESC-23-008. The
  `protected-decision-required` and `human-outcome-required` failure details
  are their StopDetails plus `stop_ids`, `attention_ids`, `event` (null), and
  `replayed` (`false`), as DEC-226 part 4 intended. The envelope, StopDetails,
  and VAL-056/057 text are unchanged.

## Prior Phase and Future Qualification

- DEC-220/221: Phase 22 and its CI repair are complete; see the
  [wrap audit](../phases/phase-22/wrap-up-audit.md) and
  [hosted evidence](../phases/phase-22/qa/task-20-hosted-ci.md).
- Claims qualification remains Phase 24; external adapter qualification Phase 25.

## Open Carry-Forward Register

| ID | Disposition / trigger / owner |
|---|---|
| CF-008 | Future TL: retained Issue-cancellation vocabulary use case; Run abandonment does not trigger it. |
| CF-018 | Future editor: split production internal/sqlquery/sqlquery.go on its next causal edit; unchanged this phase. |
| CF-023 | Workflow maintainer: replace/reconcile temporary repository delivery tooling in dedicated maintenance. |
| CF-026a | Future product design: optional inspection projection beyond accepted recovery resume; unscheduled candidate. |
| CF-026b | Human/future TL: separate Issue lifecycle product decision; deferred. |
| CF-029 | Workflow maintainer: QA evidence-token validator improvements in dedicated maintenance. |
| CF-030 | Workflow maintainer: report-lint rejected-token diagnostics in dedicated maintenance. |
| CF-033 | Workflow maintainer: AGENTS package inventory and nonexistent quality section locator; execution-scope/TRD part delivered. |
| CF-036 | Dedicated workflow retrospective/maintenance: cross-phase/task ledger identities, fail-closed ESC registration and observed control friction. DEC-190/191 bounded repairs do not close it. |
| CF-037 | Phase 25 TL/adapter owner: exact coordinator/port obligations in Phase 22 kit-contract-mapping.md; fid conformance alone is insufficient. |
| CF-038 | Future causal editor of cmd/fid/live_output_test.go: same-stream coalescing hardening. Post-truncation output_bytes documentation half closed under DEC-189/220. |
| CF-039 | Dedicated workflow maintainer: executable cross-check of contract clauses, schema, task ownership and acceptance criteria; existing scope-audit item now formally indexed. |
| CF-040 | Partially scheduled in Phase 23 task 14 (VAL-073) for the Human views that phase introduces; polish/verbosity review of pre-existing views stays with Phase 25 operating-experience qualification. |
| CF-041 | Workflow/CI maintainer: hosted CI runs only format, lint, build, shard verification, and the ordinary race/coverage shards. The tagged integration suite and the rebuild benchmark are local-only gates, so a green PR is not phase-level evidence. Dedicated maintenance. |
| CF-042 | Workflow/test-infra maintainer: `cmd/fid/integration/main_test.go` leaks built-binary directories. Cause corrected under R-15/DEC-229: five tests re-execute the test binary, each child rebuilds its own ~16 MB binary, and the crash-process children never reach `RemoveAll`, so one clean run (exit 0, no kill) leaks 20 directories and 306 MB into a RAM-backed `/tmp`. Phase 23 task 1 makes children reuse the parent's published binary and keeps a 24-hour prune only as the residual safety net for a parent killed mid-run; the durable cleanup-on-kill design stays in dedicated maintenance. |
| CF-043 | Future TL (after Phase 23): open one SQLite connection per command. Each command re-opens the projection several times (`journal.ensureProjectionFresh`, `workflow.openWorkflowDB` at 32 sites), and each open re-parses the schema, about 33% of race-test CPU (`qa/esc-23-005-repair.md`). Needs planning plus VAL coverage for rebuild/catch-up file replacement, transaction scope and locking. |
| CF-044 | Future causal editor of `cmd/fid` test helpers: `mustJournalPath` assumes one journal file, but journals split per UTC day. A gate crossing 00:00 UTC fails with `journal files = 2, want 1` (ESC-23-005 attempt 1). |

CF-024/034/035 closed (IDs retired). DEC-222: CF-040/041; DEC-227/229: CF-042;
DEC-258/259: CF-043/044. See [phase decisions](../phases/phase-23/decisions.md).

## Environment

Go/tool caches remain redirected into gitignored repository directories by
.codex/config.toml; configuration changes require a fresh session.
