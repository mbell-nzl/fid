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

### DEC-130: ESC-20-001 resolved by two bounded evidence-recovery tasks before Task 4

- **Date:** 2026-08-09
- **Context:** Review 1 returned BLOCK with two Blocking findings, and
  `TRG-20-001` escalated them as `ESC-20-001` because both span completed tasks
  and leave an accepted VAL unproven. Task 1's registry test compares only
  flattened `DetailFields` names, so no test can detect a changed detail value
  type, nullable rule, or nested schema for any of the 23 codes (VAL-002).
  Tasks 2-3 prove representative malformed config input and one 16-character
  prefix, but not the required just-inside/outside domain matrix for the
  one-character prefix, the 128-byte project ID, `math.MaxInt64` executor and
  maintenance values, or overflow rejection (VAL-043). Inspection of
  `internal/cli/codes.go` shows the production registry already models typed,
  nullable, and nested detail fields with per-code subject allow-lists, so this
  is an evidence gap rather than a suspected product defect.
- **Decision:**
  1. **VAL-002 and VAL-043 keep their accepted pass/fail criteria unchanged.**
     No assertion is weakened, deferred to a later window, or closed on source
     reading. Both findings are sustained as Blocking.
  2. Recovery is **two bounded Developer tasks**, not one, because a single
     session covering both matrices is the same oversized-evidence shape that
     produced the gap: **task-3a** (exhaustive reason-code schema, subject, and
     redaction matrix, `internal/cli` only) and **task-3b** (strict-config and
     init boundary matrix, plus the Task 3 collision regression).
  3. Both are **test-first and test-led**. Product code changes only where an
     exhaustive matrix exposes a registry or parser that accepts an
     out-of-contract input; any repair stays inside the owning package and any
     needed behavior beyond contracts §2.4/§2.5 or `04_TRD_CONFIG.md` §2
     escalates instead of being invented.
  4. **Contracts §2.5 is clarified, not changed:** `duplicate-entity`
     `identity.fields`/`identity.values` are equal-length parallel arrays, and a
     length mismatch is invalid. This makes an already-implied rule assertable;
     enforcement is authorized in task-3a behind a red test.
  5. Review 1's Medium improvement (dangling `.fid` symlink collision) is
     **do now** in task-3b. QA reproduced the correct behavior manually, so only
     the committed regression is missing and the guard is one `Lstat` away from
     silent loss of a path entry.
  6. Review 1's Low improvement (`internal/config/config.go` at 512 lines) is
     **do now** in task-3b as a pure file move with no API change, because
     task-3b is the config-owning session and Task 2 enlarged the breach.
  7. `internal/bootstrap/init.go` at 427 lines is **carry-forward as CF-028**,
     trigger-based: split staging/fsync helpers only when a later task changes
     that file for another reason. QA's own note is that splitting now would
     obscure the intentionally linear activation protocol. CF-028 is registered
     in `decisions_active.md` at wrap-up.
  8. Phase Runner resumes at **task-3a → task-3b → bounded QA re-review
     `qa/review-1b.md` → Task 4**. The re-review is scoped to the two recovery
     commits, the two sustained Blocking findings, and the Medium improvement;
     it does not re-open unrelated Tasks 1-3 surface. Task 4 starts only on a
     non-BLOCK re-review verdict. The six accepted review windows and Review 2's
     scope are unchanged; `review-1b.md` is an escalation-resolution artifact,
     not a new cadence review, so the QA plan needs no amendment.
- **Consequences:** VAL-002 and VAL-043 gain qualifying and non-qualifying
  evidence inside their accepted owning windows, so neither has to be carried to
  Task 18 or the final review. The phase costs two extra Developer sessions and
  one bounded review. The phase task index gains 3a/3b; the eighteen accepted
  tasks and their numbering are untouched.
- **Refs:** `qa/review-1.md`; `triage-log.md` `TRG-20-001`;
  `tl-escalations.md` `ESC-20-001`; `qa/validation-contract.md` VAL-002/042/043;
  `contracts.md` §§2.4, 2.5, 3.2.1; `04_TRD_CONFIG.md` §2; `AGENTS.md`
  (file/function size budgets); tasks 3a and 3b.

### DEC-131: ESC-20-002 resolved; the real defect is the gate-evidence record, not the product

- **Date:** 2026-08-09
- **Context:** Review 1b sustained the recovery but returned BLOCK on two
  findings: VAL-043 never exercises rejected lower-bound/type/overflow values
  for `executor_limits.max_excerpt_bytes`, and the required
  `go test -tags integration -race ./...` gate fails in three `cmd/fid`
  journeys with `error: invalid project initialization argument` despite
  task-3b recording it `PASS`. Investigating the second finding exposed a
  larger problem than the review scoped:
  1. **The product is correct.** `contracts.md` §3.2.1 requires that an omitted
     `--prefix` uppercase the directory basename and be accepted **only** if it
     matches `^[A-Z][A-Z0-9]{0,15}$`. `cmd/fid/init.go:34-41` does exactly that.
     The three failures are `t.TempDir()` roots whose basename is `001`, which
     must reject.
  2. **The fixtures never converged.** `cmd/fid/e2e_v01_closure_test.go:19` and
     the other binary-driven journeys still call bare `init`. The unit-suite
     helper `runInitCommand` (`cmd/fid/init_test.go:704-706`) silently appends
     `--prefix FID`, which is why `go test -race ./...` stayed green while the
     integration journeys did not. Strict rejection has existed since Task 2's
     `TestInitCommandRequiresExplicitPrefixForInvalidDirectoryBasename`.
  3. **Every full-suite gate record in this phase is unverifiable.** Tasks 1, 2,
     3, 3a, and 3b all record `go test ... ./...` with the token
     `exit 0 (no output)`. `go test ./...` prints an `ok`/`?` line per package,
     so that token cannot come from that command; `AGENTS.md` (Commands) says
     `exit 0 (no output)` for a command that prints **is fabricated evidence**.
     `task-3b-done.md:70` goes further and records a `go test` command with the
     token `  295 internal/config/project_config.go`, which is `wc -l` output.
     Review 1 also claimed an independent integration `PASS` with the same
     impossible token. Review 1b is the first artifact in Phase 20 that actually
     ran the integration gate, and it failed.
  4. **Confirmed by Human reproduction (2026-08-09).** `go test -tags
     integration -race ./...` run by the Human from the repository root
     reproduced exactly three failures — `TestV01ClosureJourneyAndRebuildParity`,
     `TestV02aWorkflowAuthoringExitScenarioAndRebuildParity`, and
     `TestMaintenanceCheckReportsIncompatibleWorkflowAssignment` — each
     `.../001/fid init: exit status 1` with `error: invalid project
     initialization argument`. Every other package reported `ok`. The
     diagnosis is verified, not inferred.
  5. **Named probable cause of the bad token.** That run took **184.9s for
     `cmd/fid` alone** and several minutes overall, and Go prints nothing until
     packages complete. The likely mechanism is an agent whose long-running
     command was timed out or killed, then written up as `exit 0 (no output)`
     because no output had appeared. `AGENTS.md` already governs this: a hung,
     interrupted, or killed command is `BLOCKED`, never `PASS`. This is a
     probable cause, not a proven one; the disposition below is unchanged
     either way, but task-3d's failure message must name it, because telling a
     future agent "your command was probably killed — record `BLOCKED`" is far
     more useful than "invalid token".
- **Decision:**
  1. **Both Review 1b Blocking findings are sustained.** VAL-043 keeps its
     accepted criteria; the exact integration command is not narrowed,
     substituted, or waived.
  2. **No product change.** The init behavior matches the accepted contract.
     Any proposal to relax derived-prefix validation so the old journeys pass is
     a contract change and escalates; it is not authorized here.
  3. **Fixture convergence, not weakening:** repair every stale integration
     journey that initializes with a non-conforming derived prefix by supplying
     an explicit valid `--prefix`, preserving each journey's behavioral
     assertions. `runInitCommand`'s `--prefix FID` injection is **accepted and
     recorded** as a deliberate fixture default, valid because dedicated tests
     own the derived-prefix accept and reject paths; it is not to be
     re-litigated, but it may never be the only reason a suite is green.
  4. **All Phase 20 full-suite gate lines recorded with `exit 0 (no output)` are
     `UNVERIFIED`, not `PASS`,** per the DEC-102 rule that an uncaptured gate is
     never a passing state. History is **not** re-run at each old commit; that
     buys little. The honest baseline is one clean full-gate run on the current
     tree, captured with real tokens, in task-3c.
  5. **Recovery is two bounded tasks:** **task-3c** (integration-journey
     convergence, the missing `max_excerpt_bytes` rejection rows, provenance
     for the failure, and a real-token full gate) and **task-3d** (the
     evidence-token control in `report-lint.sh`, plus correction notes on the
     five affected completion reports).
  6. **A control lands now, not as a promise.** `report-lint.sh` must reject
     `exit 0 (no output)` as the token for any recorded command containing
     `go test`, and must require such a command's token to look like `go test`
     output. This phase has now lost two review cycles to unverified gates;
     a named future control is not a substitute (DEC-102).
  7. **QA correction:** the next review appends a `## Corrections` record to
     `qa/review-1.md` retracting its independent integration-reproduction row.
     QA independent reproduction must cite the real output tail.
  8. **CF-029:** extend the evidence-token control to QA review artifacts, so a
     review cannot record an impossible token either. Trigger-based; registered
     in `decisions_active.md` at wrap-up.
  9. Phase Runner resumes at **task-3c → task-3d → bounded re-review
     `qa/review-1c.md` → Task 4**, on a non-BLOCK verdict. A third BLOCK returns
     to the Human, not to another recovery task.
- **Consequences:** Phase 20 gets a truthful gate baseline and a mechanical
  control against recurrence, and VAL-043 closes inside its owning window. The
  cost is two more Developer sessions and one bounded review. No accepted
  contract, VAL, or command changed. The three journeys keep their coverage.
- **Refs:** `qa/review-1b.md`; `triage-log.md` `TRG-20-002`; `tl-escalations.md`
  `ESC-20-002`; DEC-130; phase-19 DEC-102/DEC-103; `AGENTS.md` (Commands);
  `contracts.md` §3.2.1; `qa/validation-contract.md` VAL-043; tasks 3c and 3d.

### DEC-132: Human-authorized ESC-20-003 recovery preserves Task 3d and repairs long-running gate handling

- **Date:** 2026-08-10
- **Context:** Task 3d stopped after its execution wrapper returned only the
  nested command's empty `output` field and discarded the rest of the result,
  including the `session_id` needed to keep polling. The underlying full-suite
  command was still running; completion of the outer execution cell was
  mistaken for termination of the nested command. The Task 3d blocker report
  then exposed a second tooling defect: `report-lint.sh` treated wrapped
  explanatory prose containing `go test` as command-result records. The Human
  authorized a recovery super-agent outside the normal one-task/session and
  role boundaries solely to preserve and complete Task 3d, repair these
  workflow controls, capture the exact gates, commit the recovery, and close
  ESC-20-003.
- **Decision:**
  1. Preserve the existing Task 3d implementation and historical report
     corrections. Repair the false match by validating only recorded
     command-result bullets and excluding legitimate `BLOCKED`/`UNVERIFIED`
     records from PASS-token plausibility checks. Keep every existing rejection
     case and add a wrapped-blocker regression.
  2. Canonical Codex prompts for Developer, QA Lead, and QA Reviewer must
     preserve the complete nested execution result, retain and poll any
     `session_id` with `write_stdin`, and claim termination only from an
     `exit_code` or explicit termination signal. Phase Runner must reject a
     report that confuses outer execution-cell completion with nested-command
     termination. Generated provider prompts and `workflow-lint.sh` enforce the
     rule.
  3. Run Task 3d's exact command sequence without narrowing or substitution,
     capture real `go test` package-output tokens, finalize its completion
     report, and use the normal task-commit/final-report-commit sequence.
  4. The recovery may also commit its diagnosis, canonical workflow rule,
     generated prompts, decision, escalation closure, and append-only ledger
     reconciliation. No product code or accepted VAL/contract changes are
     authorized.
  5. **Do not begin the broader test-performance refactors.** They belong to a
     separate Human-requested infrastructure follow-up pack inserted before
     Task 4, after Task 3d and ESC-20-003 are closed.
- **Consequences:** The third BLOCK is resolved without discarding completed
  work or weakening any gate. Long-running silence remains expected behavior;
  agents retain the actual session until it terminates. Task 4 remains blocked
  behind review-1c and the separately authorized follow-up pack.
- **Refs:** Human authorization dated 2026-08-10; `TRG-20-003`;
  `ESC-20-003`; task-3d; `AGENTS.md` (Commands); `09_AI_WORKFLOW.md` §8;
  `report-lint.sh`; `workflow-lint.sh`.

### DEC-133: Human-approved infrastructure pack isolates integration work without weakening gates

- **Date:** 2026-08-10
- **Context:** The Human approved the remediation plan recorded in
  `reviews/2026-08-10-gate-execution-and-test-runtime-diagnosis.md` after
  Review 1c closed Tasks 3c/3d with a non-BLOCK verdict. The current race gates
  spend about 294 seconds combined on this host. The integration command reruns
  the ordinary suite, four journeys rebuild one binary, command tests repeatedly
  initialize durable workspaces, process-wide `os.Chdir` prevents parallelism,
  and four fixed sleeps consume 7.3 seconds.
- **Decision:**
  1. Insert Tasks 3e-3l before Task 4, tracked by
     `infra-improvement-pack.md`: initialization benchmark; integration
     classification/package and gate migration; shared integration binary;
     reusable workspace fixture; working-directory injection; bounded
     parallelism; claims clock; maintenance clock-fixture convergence.
  2. The dedicated integration package owns real-process journeys only. The
     twelve integration-tagged in-process tests become ordinary `cmd/fid`
     tests; the four binary journeys move with their helpers to
     `cmd/fid/integration`. No test or assertion is removed.
  3. After Task 3f, the accepted integration gate is
     `go test -tags integration -race ./cmd/fid/integration/...`. Task 3f runs
     both old and new commands as transition evidence and updates every
     canonical gate surface plus remaining Phase 20 task/QA commands.
     Completed historical specs and reports are not rewritten.
  4. A pristine workspace fixture is test-only and may replace initialization
     only where initialization is not the subject. Init, bootstrap, migration,
     rebuild, compatibility, and durability tests keep the production path;
     parity and per-test isolation are mandatory.
  5. Production CLI behavior is unchanged. Tests receive an explicit working
     directory through the root-command boundary before any `t.Parallel` is
     added. Parallelism is bounded to tests with isolated filesystem, database,
     environment, stream, and signal state and must survive repeated race runs.
  6. Claims receives a narrow optional time input with real UTC time as the
     zero-value production default. Maintenance already has `CheckOptions.Now`;
     Task 3l uses that accepted boundary instead of adding another abstraction.
  7. `InitializeWorkspace` durability semantics are benchmarked, not optimized.
     No rebuild verification, migration, fsync, activation, journal, or DB
     parity behavior may be weakened for test speed.
  8. No dependency is added. No product contract, VAL, schema, journal format,
     or CLI output changes. The same-host target is at least 40% lower median
     combined race/coverage plus dedicated-integration wall time than the
     roughly 294-second baseline; correctness remains the hard requirement.
  9. After final gates, a fresh QA Reviewer independently audits the frozen
     commit range and writes `qa/infra-review-1.md`. Task 4 remains blocked
     until every Blocking finding is fixed and the pack is complete.
- **Consequences:** The integration gate ceases to duplicate ordinary test
  execution, fixtures become reusable without mocking persistence,
  process-global directory mutation no longer blocks safe concurrency, and
  wall-clock waits become deterministic. The campaign costs eight separately
  committed TDD tasks plus an independent review. The anchor, not chat history,
  is authoritative across compaction or session boundaries.
- **Refs:** Human approval dated 2026-08-10;
  `reviews/2026-08-10-gate-execution-and-test-runtime-diagnosis.md`;
  `qa/review-1c.md`; tasks 3e-3l; `AGENTS.md`; `08_QUALITY_PLAN.md`;
  `09_AI_WORKFLOW.md` §8.

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
