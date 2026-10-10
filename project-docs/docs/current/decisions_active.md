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
- ADR-0012: vertical scenario phases 20..26 (DEC-272), explicit removal owners.
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
- Gates D3.1/D3.2: vertical delivery. Phases 21, 22 and 23 delivered the
  one-Issue skeleton, failure/recovery (DEC-220), and Review/Human judgment
  (DEC-271) with their qualified capabilities.
- DEC-156: qualifying Evidence guards terminal-success Issue transition;
  Envelope completion reuses the same evaluator after the Issue is terminal.
- DEC-157/177: Phase 21 delivered wire authority is
  `../phases/phase-21/public-interface.md`; its removal inventory is
  `v1-convergence-inventory.md`. Public enums are kebab-case; bootstrap has no role.

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
- DEC-002: autonomous TL-Lite resume remains permanently enabled.
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

## Phase 23 Delivered Commitments (DEC-271)

Phase 23 is complete; full reasoning is in
[phase decisions](../phases/phase-23/decisions.md). Standing rules:

- DEC-225: Review independence conflicts only with mutating Claims and their
  Runs; the same independent Reviewer may review successive rework cycles.
- DEC-226: a **Failure code** names one refused request; an **Evaluation
  reason** says why one Obligation is unsatisfied. File each in its set.
- DEC-235/236: the `Condition` value uses `omitempty` only on
  `outcome_class`, `decision_class`, `finding_id`; tagging any other field
  changes every stored `semantic_key` and `request_digest`.
- DEC-249: an Obligation blocks only the guarded actions it declares.
- DEC-250/269: Stop resolution is keyed on the recorded Condition with one
  predicate shared by manual and atomic paths. The `review-rework-exhausted`
  Stop resolves through an approved `review-rework-extension` Human Outcome,
  once per Envelope revision; generic admission ignores the rework budget.
- DEC-251/260: Human Outcomes need no Envelope grant (Activation supplies
  authority); requirement writes need a trusted Actor and no capability.
- DEC-258/259: tests default to `t.Parallel()`; race gates carry the modernc
  `checkptr` exclusion plus a phase-end full-`checkptr` run. DEC-227's
  `-timeout 30m` is added per phase when package time passes the 80% warning.
- DEC-263: Stop-blocked denials carry `stop_ids`, `attention_ids`, `event`
  (null), `replayed` in failure details; StopDetails unchanged.
- DEC-264: a down migration restores the prior live schema state exactly.
- DEC-270: journals from earlier development binaries fail closed under
  ADR-0011; no legacy replay path is added for pre-adoption history.

## Prior Phase and Future Qualification

- Wrap audits: [Phase 23](../phases/phase-23/wrap-up-audit.md), [Phase 22](../phases/phase-22/wrap-up-audit.md).
- DEC-272: new Phase 24 operating experience and input convergence (scope in
  [phase decisions](../phases/phase-23/decisions.md)); Claims qualification is
  Phase 25; external adapter qualification Phase 26.

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
| CF-033 | Workflow maintainer: AGENTS package inventory (lists nonexistent packages) and nonexistent quality section locator. `internal/review/` added at Phase 23 wrap. |
| CF-036 | Dedicated workflow retrospective/maintenance: cross-phase/task ledger identities, fail-closed ESC registration and observed control friction. DEC-190/191 bounded repairs do not close it. DEC-286: the sequence verifier accepts any existing DEC as RUN-RESUME resolution authority. DEC-288: a planned mid-phase Human checkpoint can only resume from its STOP through an evidence-only task recovery; there is no Human-evidence route to REVIEW. |
| CF-037 | Phase 26 TL/adapter owner: exact coordinator/port obligations in Phase 22 kit-contract-mapping.md; fid conformance alone is insufficient. |
| CF-038 | Future causal editor of cmd/fid/live_output_test.go: same-stream coalescing hardening. Post-truncation output_bytes documentation half closed under DEC-189/220. |
| CF-039 | Dedicated workflow maintainer: executable cross-check of contract clauses, schema, task ownership and acceptance criteria; existing scope-audit item now formally indexed. |
| CF-040 | Phase 26 operating-experience qualification: polish of pre-existing views and the A6 hint on/off switch; then, after Phase 26 and live use, a full Human-output design revisit (layout, spacing, tables, scannability) (DEC-272). Also (DEC-280): `evidence-stale` has no Phase 24 `next:` row or `reason:` override, so a stale-Evidence denial names no next action; the Phase 24 journey avoids it by design. |
| CF-041 | Workflow/CI maintainer: hosted CI runs only format, lint, build, shard verification, and the ordinary race/coverage shards. The tagged integration suite and the rebuild benchmark are local-only gates, so a green PR is not phase-level evidence. Dedicated maintenance. |
| CF-042 | Workflow/test-infra maintainer: `cmd/fid/integration/main_test.go` binary-directory leak. Phase 23 task 1 made children reuse the parent binary and added a 24-hour prune; durable cleanup-on-kill stays in dedicated maintenance. |
| CF-043 | Future TL (after Phase 23): open one SQLite connection per command. Each command re-opens the projection several times (`journal.ensureProjectionFresh`, `workflow.openWorkflowDB` at 32 sites), and each open re-parses the schema, about 33% of race-test CPU (`qa/esc-23-005-repair.md`). Needs planning plus VAL coverage for rebuild/catch-up file replacement, transaction scope and locking. |
| CF-044 | Future causal editor of `cmd/fid` test helpers: `mustJournalPath` assumes one journal file, but journals split per UTC day. A gate crossing 00:00 UTC fails with `journal files = 2, want 1` (ESC-23-005 attempt 1). |
| CF-045 | SCHEDULED in Phase 24 (DEC-272; plan approved DEC-277): A4 Stop/Attention rendering, A5 Human failure reason lines, A6 once-per-section disclaimers, A7 recorded reasoning in Human views, A8 option C (`--help` input docs, `--example`, guidance), and repository-input convergence with the input-as-written digest rule. Closes when Phase 24 is accepted. |
| CF-046 | Phase 25 TL (parallel worktrees), recorded under DEC-278. (a) A plain-copy worktree refusal returns each leaf's generic fallback `invalid-argument` (`argument:"input"`), as Phase 24 contract section 5 accepts ("refused as today"); a dedicated worktree-binding failure belongs with parallel worktrees. (b) Leaves that observe outside the Phase 24 observing-leaf list (`run.recover`, `run.checkpoint`, `artifact.verify`, terminal `issue.transition`, `resume`, `run.inspect`) keep their earlier unobservable-repository codes; decide whether they should return `repository-unavailable`. (c) `run.inspect` refuses a plain copy with `recovery-state-conflict` (`repository identity changed`, exit 4), so its Human comparison line only says `same worktree` (DEC-287). Decide, with CF-026a, whether inspection should show a different worktree. |
| CF-047 | Phase 26 TL (qualification and baseline convergence), recorded under DEC-280: Phase 21 section 2 grammar items Phase 24 task 9a does not deliver. (a) Flag-combination rules keep pre-Phase-24 failures, often `argument:"request"`: exclusive selectors such as `issue create --epic`/`--parent`, exactly-one or at-least-one groups, and kind-conditional flags. (b) Live flag sets that differ from their accepted rows: `activity list` has live `--on` and no `--limit`, while the Phase 21 row lists `--ref` and `--limit <1..10000>`; plus any others the task 9a report lists. Until then, malformed `--since`/`--on` values reach the retained-read fallback (`argument:"request"`, DEC-285); declare their domains with that row. (c) `stats.tokens`/`stats.cost` Human success prints the JSON success envelope, despite the row's Human obligation (DEC-285). (d) The `sql`/`journal inspect` root-flag refusal still says "until its fid-cli.v1 transport is migrated", though Phase 21 section 5 keeps them outside `fid-cli.v1` (DEC-286). |

CF-024/034/035 closed (IDs retired). DEC-222: CF-040/041; DEC-227/229: CF-042; DEC-258/259: CF-043/044; DEC-271: CF-045; DEC-278/287: CF-046; DEC-280/285/286: CF-047. See [Phase 23](../phases/phase-23/decisions.md) and [Phase 24](../phases/phase-24/decisions.md) decisions.

## Environment

Go/tool caches stay redirected into gitignored repo directories by .codex/config.toml; config changes need a fresh session.
