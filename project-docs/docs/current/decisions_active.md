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
- DEC-179: Phase 21 wrapped with 87/87 VAL and 4/4 Human PASS. Its local
  DEC-152..178 decisions are archived in `../phases/phase-21/decisions.md`.

## Standing Delivery Constraints

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
  that owns the fenced Event; the owning table, not the fence file's author,
  fixes ownership, and the last owner removes the empty fence and its guard.
  Fence removal is implementation of an accepted contract, never a boundary
  change, and the displaced fence assertions are scheduled convergence.
- DEC-193: a calendar-dependent retained test is DEC-189 nondeterminism, not
  convergence and not DEC-175. Assign it to the observing task, repair the date
  binding at its source with recorded introducing-commit provenance, and never
  pin an absolute date on an Event appended through the journal writer into a
  live-clock workspace.
- DEC-194: an accepted allocation that exists but was never delivered is bounded
  repair, not an ambiguous owner. Read the task spec, schema owner column,
  public-interface task column, and contract subjects before classifying.
- DEC-195: an accepted contract clause that adds a closed-domain member requires
  a matching schema row. Where the row is missing, amend the schema and keep the
  assertion; never substitute an admitted value or weaken the assertion. Assign
  execution to the first task that cannot pass without it, with provenance to the
  original owner recorded rather than relabelled.
- Current V1 delivery scaffolding is binding engineering procedure, not forward
  fid product scope. No failed gate authorizes a protected definition change.

## Current Phase

- DEC-222 (ACCEPTED by Human, 2026-09-14): Phase 23 "Review and Human judgment"
  is approved and Heavy. The accepted authority is `contracts.md`,
  `public-interface.md`, `schema.md`, `execution-scope.md`,
  `v1-convergence-inventory.md`, `contract-coverage.md`, the now 79-assertion
  `qa/validation-contract.md`, and the ordered 18-task pack, all under
  `../phases/phase-23/`. Six canonical spec amendments are accepted with it and
  are applied by their owning implementation tasks, not in advance.
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
  describes why one Obligation is unsatisfied and carries none of those. One
  glossary term named both sets, which is how `finding-unresolved` was filed as
  a failure code. `GLOSSARY.md` now defines each; register every new identifier
  in the set it belongs to.
- DEC-228 (ACCEPTED, Tech Lead) carries the remaining dispositions, including
  three DEC-167/173 scope corrections and five new assertions.
- DEC-235 (ACCEPTED by Human 2026-09-20) resolves ESC-23-001 with Option A: the
  three new `Condition` fields are `omitempty` and absent on the types they do
  not apply to, the six Phase 22 types stay byte-identical, and exclusivity is
  structural in per-type journal shapes keyed on `type`.
- DEC-236 (ACCEPTED, Tech Lead 2026-09-20) dispositions QA pass 7. Two standing
  consequences: task 9's new `case "Condition"` dispatch must **end by calling**
  `validateCurrentValueRelations("Condition", object)` - removing the
  `currentValueShapes["Condition"]` key deletes the Phase 22 nullable binding's
  only call site, and no test asserts its message; and `review-recorded` carries
  **no** discriminator field, per accepted `contracts.md` section 11, so
  `finding_id` belongs to `finding-resolved` alone. VAL-078 is refined to the
  accepted pairing and strengthened; the contract stays at 79 rows.
- Execution has not started. No decision remains open and no escalation is
  registered; tasks 1-8 are dispatchable now and 9-11 once the QA Lead confirms
  this batch. See [phase brief](phase_brief.md) and
  [Phase 23 decisions](../phases/phase-23/decisions.md).

## Latest Phase Closure

- DEC-220: Phase 22 complete, Human-accepted; 18+1 tasks, 80/80 VAL,
  84 deterministic QA items and 6/6 Human checks. Post-wrap PR CI timeout
  follow-up is ESC-22-005; merge awaits its disposition and successful CI.
  See [acceptance](../phases/phase-22/qa/final-acceptance.md) and
  [wrap audit](../phases/phase-22/wrap-up-audit.md). No phase currently active.
- Delivered contracts: Phase 22 contracts.md/public-interface.md/schema.md amend
  the Phase 21 baseline. DEC-195 retains cancelled Receipts as their own outcome.
  DEC-216 requires typed idempotency snapshots; DEC-217 retains full product
  archive limits with small routine and real-scale phase/affected-change tests.
- Recovery retains exact authority, unknown-liveness rejection, immutable budget
  charges and original replay results. Adoption is not proof; imported history
  grants no destination authority. DEC-219 Human guidance is accepted.
- Claims capability remains unavailable until Phase 24 parallel qualification;
  external adapter/all-twelve-capability qualification remains Phase 25.

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

CF-024/034/035 closed with Phase 22 evidence; historical dispositions and closed
IDs remain in the wrap audit and must not be reused. All triggers checked
against the actual phase diff.

Under DEC-222, CF-008/018/023/026a/026b/029/030/033/036/037/038/039 are
re-deferred unchanged, CF-040 is partially scheduled, and CF-041 is registered.
CF-042 is registered under DEC-227/DEC-228 from the QA plan's I-16, with its
cause and mitigation corrected under DEC-229 from the readiness refresh's R-15.
The full disposition table is in
[Phase 23 decisions](../phases/phase-23/decisions.md).

## Environment

Go/tool caches remain redirected into gitignored repository directories by
.codex/config.toml; configuration changes require a fresh session.

- DEC-221: Human approved task 20 performance repair instead of extending the
  timeout. Immutable fixture copies, one full race/coverage execution partitioned
  exhaustively across four CI workers, combined coverage and retained quality
  failure propagation. ESC-22-005 resolved; hosted run34805256766 PASS on6b8fe50.
