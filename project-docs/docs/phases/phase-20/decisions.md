# Phase 20 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes, and
> wrap-up decisions. DEC numbering is global and never resets. Phase 19 ended at
> DEC-106.

## Decisions

### DEC-107: Phase 20 operating mode is **Heavy**

- **Date:** 2026-08-09
- **Context:** The playbook mode test promotes a phase to heavy on any schema,
  public-interface, dependency, or security change. Phase 20 does all three of
  the first: it changes the `issues` table, replaces the entire machine output
  contract, and removes a runtime dependency. It is also the first phase of a
  newly reset delivery direction.
- **Decision:** Heavy. Full artifact set: validation contract, QA plan,
  mid-phase reviews, final review, human checklist, run ledger, triage log.
  Phase Runner orchestration.
- **Consequences:** Eighteen bounded Developer tasks with mid-phase reviews
  after tasks 3, 6, 9, 12, 15, and 17. No task may close on attested or
  remembered gates.
- **Refs:** `playbook.md` §1; roadmap §5.

### DEC-108: The Phase 20 slice contract lives in the phase directory, referenced from the CLI spec

- **Date:** 2026-08-09
- **Context:** ADR-0011 requires accepted target JSON contracts to be recorded
  in specs before implementation. `06_DESIGN_PLAN_CLI.md` §1 already states that
  exact leaf schemas are accepted with their owning slice. Inlining every leaf
  schema for six slices would turn the target CLI spec into a changelog.
- **Decision:** Cross-cutting rules (envelope shape, category set, slice code
  list) are amended into `06_DESIGN_PLAN_CLI.md` §2, §5, and a new §8. Exact
  leaf schemas, the removal inventory, and the legacy disposition live in
  `phase-20/contracts.md`, which §8 names as proposed until Human approval. The
  physical DDL is proposed in `03_TRD_SCHEMA.md` §8 per its section 7 gate.
- **Consequences:** Each future slice adds one row to `06` §8 and one subsection
  to `03` §8. Phase 25 consolidates them into the first supported baseline.
- **Refs:** ADR-0011; `03_TRD_SCHEMA.md` §7; `06_DESIGN_PLAN_CLI.md` §1.

### DEC-109: The JSON wrapper is named **Response Envelope**

- **Date:** 2026-08-09
- **Context:** `Work Envelope` is a core durable entity. The roadmap and CLI
  spec both call the JSON wrapper an "envelope." An implementation task that
  says "the envelope" would be ambiguous, which is an `AGENTS.md` hard stop.
- **Decision:** The wrapper is the **Response Envelope**, defined in
  `GLOSSARY.md` under a new "Machine Interface" section alongside *operation
  name*, *reason code*, and *project root*. The bare word "envelope" continues
  to mean Work Envelope.
- **Consequences:** Package, type, and test names in Phase 20 use "response
  envelope" or `internal/cli`, never `Envelope` alone.
- **Refs:** `GLOSSARY.md`; `AGENTS.md` hard stops.

### DEC-110: Only the retained target command tree gets the Response Envelope in Phase 20

- **Date:** 2026-08-09
- **Context:** Removing `--format` touches every command, including families
  that Phases 21-25 delete outright (`policy`, `dev`, `claim`, `run`, `concern`,
  `requirement`, `stats`, `sql`, `activity`, `maintenance vacuum`). Converting
  them to a contract they will never honor is wasted work; leaving `--format` in
  place violates the Phase 20 removal column.
- **Decision:** Twenty-one operations, listed in `contracts.md` §3, emit the
  Response Envelope. Every remaining `--format` is removed. Legacy leaves keep
  their established unenveloped JSON or Human mode, listed exactly in
  `contracts.md` §6 with an expiry phase.
  A drift-guard test (VAL-035) holds both lists as committed literals.
- **Consequences:** Two machine output styles coexist during convergence,
  alongside existing legacy Human commands. This is explicitly residue with a
  named expiry, not a second machine contract. The
  guard makes silent drift a test failure. Explicit `--project` or `--json` on
  a legacy leaf rejects rather than being accepted and ignored. Legacy JSON has
  no byte-compatibility promise.
- **Refs:** roadmap §5, §6; `contracts.md` §3, §6.

### DEC-111: Issue-path fields removed in Phase 20

- **Date:** 2026-08-09
- **Context:** The roadmap assigns "Issue-path fields or commands with no
  retained meaning" to Phase 20 without enumerating them. Some current fields
  are genuinely orphaned; others are role-bound or execution-bound and are owned
  by Phases 21-23.
- **Decision:** Remove `in_inbox` (no target concept), `est_tokens`
  (contradicts PRD-TELEM-003's observed-telemetry rule), `due_date` (scheduling
  is kit-owned under ADR-0004), `triage_status` from the contract (`Triage
  Entry` is a retired term), `completions` from the contract
  (completion-report-as-state), and `fid issue get --section` with
  `issues.SectionValue` (a sub-selector over a Markdown-oriented composite).
  Keep `spec_md_path`. Internal owner/concern storage and loaders may remain for
  their replacement-bound paths, but `owner_role`, `owner_actor`,
  `concerns_raised`, and `open_concerns_count` are omitted from canonical
  `fid-cli.v1` Issue JSON.
- **Consequences:** One migration and a contract-level trim. The `completions`
  and `gate_runs` writer paths survive until their owning phases, so nothing
  becomes unreachable in this phase.
- **Refs:** `contracts.md` §5.5; PRD-TELEM-003; ADR-0004; `03_TRD_SCHEMA.md`
  §3.2, §6.

### DEC-112: `revision` lands now; idempotency keys wait for Phase 21

- **Date:** 2026-08-09
- **Context:** The roadmap column says "revision/idempotency conventions where
  used." Concurrent Issue updates are part of the standalone journey; retried
  port mutations are not.
- **Decision:** Add `issues.revision` and `--expected-revision` on
  `issue update` and `issue transition` only. Do not add `--idempotency-key` or
  the `idempotency_results` table. Both belong to the Phase 21 port surface.
  The caller supplies only an expected prior revision. The Event records the
  resulting absolute revision in `Values` and the observed prior in both
  `Prior` and `Preconditions`, using existing deterministic Effect semantics.
- **Consequences:** `state-conflict` is exercised in Phase 20;
  `idempotency-conflict` is not, and its absence is recorded in the validation
  contract's uncovered-requirements table.
- **Refs:** `contracts.md` §3.6, §1; `03_TRD_SCHEMA.md` §3.1; DEC-122.

### DEC-113: Fail-closed V1 authority is preserved and given a typed shape

- **Date:** 2026-08-09
- **Context:** The roadmap forbids removing fail-closed V1 authority before
  Phase 21 proves the generic replacement. But the Phase 20 error contract must
  still be machine-readable, and today's authority failures are untyped text.
- **Decision:** `internal/policy`, `--by`, the `roles` table, and every existing
  permission check stay. Their failures map to three semantic reason
  codes: `authority-missing`, `authority-unknown-actor`, `authority-denied`.
  VAL-006 proves the rejection path still rejects and still performs no work.
- **Consequences:** Phase 21 removes the V1 checks, not automatically the
  general codes. Its generic authority contract decides which semantic codes
  remain or are refined. Mapping an untyped failure is not a weakening.
- **Refs:** roadmap §5 Phase 20 removal column; `06_DESIGN_PLAN_CLI.md` §5.

### DEC-114: Phase 20 owns dependency-cycle rejection over existing links

- **Date:** 2026-08-09
- **Context:** PRD-ISSUE-004 requires rejecting parent cycles, dependency
  cycles, invalid hierarchy, and missing references. Planning found no cycle
  detection anywhere in the codebase. "Relate" and "explain failure" are both
  Phase 20 journey steps, so the gap is in this slice's acceptance boundary.
  Phase 24 owns the normalized `issue_dependencies` table.
- **Decision:** Implement cycle, self-link, missing-reference, and hierarchy
  rejection over the existing `links` mechanism as task 11. Do **not** introduce
  the `issue_dependencies` table or `fid issue dependency` command family; those
  stay with Phase 24. Human direction confirmed the integrity invariant remains
  in Phase 20; it is not a carry-forward candidate.
- **Consequences:** A real domain invariant exercises the typed error contract
  rather than only not-found and authority paths.
- **Refs:** PRD-ISSUE-002, PRD-ISSUE-003, PRD-ISSUE-004; roadmap §5 Phase 24.

### DEC-115: Human success output moves from stderr to stdout on enveloped commands

- **Date:** 2026-08-09
- **Context:** Today `fid issue create` and `fid init` write their confirmation
  to stderr. `05_DESIGN_PLAN.md` §3 defines human mode as concise terminal
  output with errors on the error channel, and the Phase 20 contract puts JSON
  success on stdout and JSON errors on stderr. Keeping success on stderr would
  make the two modes inconsistent.
- **Decision:** For the twenty-one enveloped operations, human-mode success text
  goes to stdout and only errors go to stderr. Legacy commands are unchanged.
- **Consequences:** Test helpers that parse created Issue IDs from stderr must
  move to stdout as their command families migrate in tasks 2-12 and are
  corrected centrally rather than preserving the historical channel mistake.
- **Refs:** `05_DESIGN_PLAN.md` §3; `contracts.md` §2.2.

### DEC-116: Deferred items carry a named expiry phase

- **Date:** 2026-08-09
- **Context:** Roadmap §5 says the removal column is an ownership map, not
  permission to keep a retired surface until Phase 25, and that a moved item
  needs a reason and a concrete new expiry point.
- **Decision:** Every Phase 20 deferral is recorded with an owning phase in
  `contracts.md` §1 and §6. The uncovered-requirements table in the validation
  contract records why each unasserted PRD requirement has no Phase 20
  assertion.
- **Consequences:** Phase 21 planning starts from a written deferral ledger
  rather than rediscovery.
- **Refs:** roadmap §5; `contracts.md` §1, §6.

### DEC-117: No Phase 21 work is pulled forward

- **Date:** 2026-08-09
- **Context:** The Phase 20 contract needs project identity and a machine
  envelope, which are adjacent to the Phase 21 backend port. `fid capabilities`
  in particular is described as "the first adapter operation."
- **Decision:** `fid capabilities`, Work Envelopes, Envelope Revisions, generic
  Actor authority, Claims/Runs reshape, Gate Registrations, Receipts, Evidence,
  and idempotency keys are all excluded. Gate D3.2 remains open and is not
  pre-empted. Project identity in Phase 20 means only `--project`, upward
  discovery, and a stable `[project].id` in `.fid/config.toml` per
  `04_TRD_CONFIG.md` §2.
- **Consequences:** No roadmap deviation to escalate. If implementation
  discovers that the standalone journey cannot be proven without a Phase 21
  surface, that is an escalation, not a scope decision.
- **Refs:** roadmap §5 Phases 20-21; ADR-0012; `decisions_active.md` Gate D3.2.

### DEC-118: Restored V1 workflow self-check re-slices Phase 20 before approval

- **Date:** 2026-08-09
- **Context:** Phase planning overlapped a mistaken rewrite of
  `09_AI_WORKFLOW.md` as the not-yet-adopted V2 target. Commit `f04007d`
  restored the binding V1 workflow. Audit found oversized tasks, a missing
  review window, premature "accepted" labels, an impossible projection-only
  revision increment, a non-exhaustive legacy command list, and stale workflow
  commentary.
- **Decision:** Retain the approved vertical Phase 20 outcome but re-slice ten
  tasks into eighteen one-session tasks after independent audit. Reviews follow
  tasks 3, 6, 9, 12, 15, and 17. Exact request/result/error schemas, revision
  Event semantics, atomic init, legacy flag rejection, exact leaf inventories,
  and proposal/approval markers are repaired before Human approval.
- **Consequences:** No product code or Phase 21 capability is pulled forward.
  The larger task count is execution granularity, not a broader phase outcome.
- **Refs:** `09_AI_WORKFLOW.md` §§2.2, 4; phase `contracts.md`; validation contract.

### DEC-119: Workflow project scope uses `--project-wide`

- **Date:** 2026-08-09
- **Context:** Existing `workflow assign --project` is a local boolean scope
  selector. Phase 20 introduces global `--project <path>`. Cobra persistent and
  local flags cannot safely give the same spelling incompatible types/meanings.
- **Decision:** Task 2 renames the assignment scope selector to `--project-wide` before adding the root flag, and
  retain `--kind <kind>` as its mutually exclusive alternative. Global
  `--project <path>` always means fid project-root selection.
- **Consequences:** The standalone contract has no ambiguous/shadowed flag and
  may combine explicit root selection with project-wide workflow assignment.
  There is no compatibility alias because fid has no users.
- **Refs:** contracts §§2, 3.3; ADR-0011 no-compatibility decision.

### DEC-120: Project config uses `issue_prefix`, not a compatibility alias

- **Date:** 2026-08-09
- **Context:** Accepted `04_TRD_CONFIG.md` defines `[project].id` and
  `[project].issue_prefix`; current bootstrap/parser still use `[project].prefix`.
  Phase 20 owns project identity and initialization, and no live config requires
  migration.
- **Decision:** Task 2 changes writer, reader, fixtures, and result vocabulary to
  `issue_prefix`. The old key rejects as unknown rather than being accepted as
  an alias.
- **Consequences:** Runtime configuration and Response Envelope use one term and
  the first supported baseline carries no compatibility parser.
- **Refs:** `04_TRD_CONFIG.md` §2; contracts §§2.1, 5.5.

### DEC-121: Canonical Issue mutations cannot write hidden V1 owner state

- **Date:** 2026-08-09
- **Context:** Removing owner fields from canonical Issue JSON while retaining
  `issue create/update --owner-role/--owner-actor` would expose a mutation whose
  resulting state cannot be read and whose semantics are removed in Phase 21.
- **Decision:** Remove those four public flag registrations/request paths in
  Task 9. Retain internal default owner values, columns, and loaders only where
  replacement-bound V1 safety behavior still needs them.
- **Consequences:** `fid-cli.v1` has no hidden role-owner mutation. Phase 21 can
  replace internal state with generic Actor attribution without a published
  compatibility promise.
- **Refs:** ADR-0011; contracts §§3.2, 3.4, 5.5; DEC-111/113.

### DEC-122: Revision covers the complete Issue aggregate

- **Date:** 2026-08-09
- **Context:** Row-only revision would not detect acceptance-criterion, Issue-
  label, or dependency changes even though those facts alter the Issue consumed
  by readiness and future Work Envelopes.
- **Decision:** `issues.revision` is aggregate-wide. Create starts at 1;
  update/transition, AC add, Issue-label add/remove, and link add/remove advance
  every affected Issue exactly once. Until Phase 21 deletes them, legacy
  `dev start`/`finish`/`block` status mutations do the same without gaining a
  canonical transport. Results and Event fixtures expose revisions where the
  canonical operation has a result.
- **Consequences:** Tasks 10-11 share Task 6's revision helper and cite the same
  VALs. A two-Issue link advances both Issues under the single write lock.
- **Refs:** contracts §§3.4, 3.6, 7; VAL-010/011.

### DEC-123: Journal inspection stays legacy in Phase 20

- **Date:** 2026-08-09
- **Context:** Publishing current Event inspection would make V1 role-shaped
  attribution and speculative revert verdicts part of `fid-cli.v1`, immediately
  before generic Actor work and the clean Event baseline.
- **Decision:** Remove `journal.inspect` from the retained operation set. Keep
  raw event/location lookup as unenveloped legacy JSON until Phase 25, but delete
  the speculative `revert` verdict projection/helpers now. Root `--json`
  remains rejected.
- **Consequences:** The exact inventory becomes 21 retained/21 legacy. Phase 20
  still contracts internal Event changes for deterministic replay without
  making the raw journal a public machine schema.
- **Refs:** ADR-0011 candidate ledger; contracts §§3, 6-7; VAL-028/035.

### DEC-124: Link endpoint type is explicit and IDs remain opaque

- **Date:** 2026-08-09
- **Context:** Existing link code infers requirement type from `PRD-`/`TRD-`
  prefixes, contradicting the accepted opaque-ID rule and misclassifying missing
  requirements as missing Issues.
- **Decision:** Link requests require explicit endpoint types and resolve those
  types against stored identity. Missing Issues and requirements use distinct
  codes; no behavior parses an ID prefix.
- **Consequences:** Link CLI requests change sharply without an alias. Requirement
  endpoints remain supported while requirement authoring itself stays legacy.
- **Refs:** `02_TRD.md` opaque-ID rule; contracts §§2.4, 3.2; VAL-046.

### DEC-125: Initialization is an atomic task boundary

- **Date:** 2026-08-09
- **Context:** Accepted config requirements demand staging, genesis append,
  rebuild verification, fsync, and atomic activation. Combining this with root
  discovery/config/flag work exceeds one Developer session and invites partial
  workspace states.
- **Decision:** Task 2 owns discovery/config namespace and changes the existing
  direct bootstrap path/fixtures to emit and consume the strict target config.
  Task 3 owns staged atomic activation and the `project.init` envelope.
  Injected failures must prove no target or staging residue.
- **Consequences:** Phase 20 has eighteen tasks and six review windows. Init is
  not considered successful until rebuild equivalence passes before activation.
- **Refs:** `04_TRD_CONFIG.md` §5; contracts §2.1; VAL-042.

### DEC-126: Phase 20 machine contracts are closed and bounded

- **Date:** 2026-08-09
- **Context:** Result-only tables left Developers to infer request flags, error
  details, ordering, redaction, and unbounded-list behavior from V1 code.
- **Decision:** The slice contract closes all 21 request schemas, input value
  domains and ranges, error subject/details, deterministic ordering, explicit
  collection limits, maintenance check names, redaction, and same-version
  unknown-field behavior. `started_without_claim` temporarily uses a fixed
  300-second rule rather than aliasing future Attempt configuration. The
  operation name is the Phase 20 CLI capability identifier;
  discovery/advertisement remains Phase 21.
- **Consequences:** Unknown inputs reject, collection overflow is typed rather
  than truncated, and changes to closed `fid-cli.v1` objects require a new
  accepted version.
- **Refs:** `08_QUALITY_PLAN.md` §6; contracts §§2.5, 3.2, 3.7; VAL-044/045.

### DEC-127: Independent audit drives the approval-ready reslice

- **Date:** 2026-08-09
- **Context:** The restored V1 self-check still missed per-criterion VAL
  citations, phase-brief control sections, atomic init, aggregate revision,
  opaque link typing, and exact machine-contract obligations.
- **Decision:** Repair those gaps before Human approval, split legacy JSON
  conversion from broad test cleanup, and require an independent consistency
  check in review 6 plus final QA.
- **Consequences:** The phase outcome is unchanged; task/validation granularity
  increases to make the accepted slice implementable without Developer inference.
- **Refs:** `09_AI_WORKFLOW.md` §§3.6, 6; `08_QUALITY_PLAN.md` §§6, 8.

### DEC-128: Canonical workflow authoring excludes V1 assurance choreography

- **Date:** 2026-08-09
- **Context:** Independent remediation review found that publishing
  `all_acs_verified` roles, completion-report, concern, and child-task process
  rules in `fid-cli.v1` would bless the role/assurance coupling forbidden by
  PRD-WF-005. The same review found aggregate revision could be bypassed through
  retained `fid dev` status mutations and that strict config activation was not
  independently implementable across Tasks 2-3.
- **Decision:** Canonical `workflow.define` authors states/transitions only and
  does not register `--require`; bootstrap-seeded V1 requirement rows remain internal to
  legacy `fid dev` until Phases 21/23. Task 6 revisions all retained dev status
  mutations until their Phase 21 deletion. Task 2 makes strict config live in
  the transitional direct bootstrap path; Task 3 changes only staging and
  activation mechanics.
- **Consequences:** Phase 20 neither publishes the retired V1 choreography nor
  introduces premature generic-precondition/Obligation semantics. Every live
  Issue status mutation preserves aggregate freshness, and every task can pass
  its full gate independently.
- **Refs:** PRD-WF-005; contracts §§3.2.1, 3.6, 7; tasks 2, 3, 6, 12;
  VAL-010/011/026/043/044.

### DEC-129: Human accepts the Phase 20 implementation pack

- **Date:** 2026-08-09
- **Context:** The revised heavy-mode slice contract, validation assertions,
  eighteen task specs, six review windows, canonical-spec amendments, and
  independent remediation review were presented for approval.
- **Decision:** Human approval accepts the complete Phase 20 implementation
  pack and its amendments to schema, config, CLI, roadmap, glossary, and current
  context. Validation rows remain `Draft` until the QA plan demonstrates
  feasible evidence and is accepted under the binding V1 workflow.
- **Consequences:** Phase Runner may initiate QA planning. No Developer task
  starts until QA-plan acceptance; after that, Task 1 proceeds through the
  accepted review cadence without reinterpreting the contract.
- **Refs:** `contracts.md`; `qa/validation-contract.md`; `current/phase_brief.md`;
  `09_AI_WORKFLOW.md` §§3.6, 4, 6.

## Spec Sync Verdict

Recorded at phase wrap-up. Planning-time spec edits are noted for continuity.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | pending wrap-up | - |
| `01_PRD.md` | pending wrap-up | - |
| `02_TRD.md` | pending wrap-up | - |
| `03_TRD_SCHEMA.md` | accepted at planning: see DEC-108/129 (new §8 slice DDL) | `03_TRD_SCHEMA.md` §8.1 |
| `04_TRD_CONFIG.md` | accepted at planning: versioned strict config and legacy pricing expiry | §§2, 5; DEC-120/125/126/129 |
| `05_DESIGN_PLAN.md` | pending wrap-up | - |
| `06_DESIGN_PLAN_CLI.md` | accepted at planning: see DEC-108/109/129 (§2, §5, new §8) | `06_DESIGN_PLAN_CLI.md` §8 |
| `07_DESIGN_PLAN_FLOWS.md` | pending wrap-up | - |
| `08_QUALITY_PLAN.md` | pending wrap-up | - |
| `09_AI_WORKFLOW.md` | pending wrap-up | No fid product feature may be added here |
| `10_IMPLEMENTATION_ROADMAP.md` | pending wrap-up | §5.1 ownership map to be replaced by the delivered inventory |
| `GLOSSARY.md` | accepted at planning: see DEC-109/129 (new Machine Interface section) | `GLOSSARY.md` |

## Triage-Log Audit

- [ ] Walked every TRG-20-NNN in `triage-log.md`.
- [ ] Confirmed each RESUME verdict was correct in hindsight.
- [ ] Flagged any TRG that should have escalated.
- [ ] Confirmed autonomous TL-Lite resume remains in effect (DEC-002; no warm-up tally).

## Workflow Notes

Completed at phase wrap-up.

- **What worked well?** <pending>
- **What was inefficient?** <pending>
- **Were tasks the right size?** <pending>
- **Did the QA review cadence feel right?** <pending>
- **Were there token-efficiency opportunities missed?** <pending>

### Deep Retrospective Recommendation

<pending>
