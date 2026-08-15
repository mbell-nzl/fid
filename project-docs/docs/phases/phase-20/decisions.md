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

### DEC-134: ESC-20-004 sustains the review finding and applies DEC-133 repository-wide

- **Date:** 2026-08-10
- **Context:** Infrastructure Review 1 found
  `internal/issues/transition_rules_test.go` retained its integration build
  constraint after Task 3f. The classification guard inspected only the direct
  `cmd/fid` parent, the ordinary gate therefore excluded this test, and the
  dedicated command could not reach `internal/issues`. The Human-approved
  pack requires every Blocking finding to be corrected and returned to the
  same reviewer before Task 4.
- **Decision:**
  1. Sustain ESC-20-004 as a gate-definition regression, not a product defect
     and not grounds to weaken the retained persistence assertion.
  2. Apply DEC-133's existing classification across the whole repository:
     integration-tagged tests may exist only in the dedicated real-process
     package; in-process tests belong to the ordinary suite.
  3. Use a failing repository-wide source inventory to reproduce the omission,
     remove only the orphaned file's build constraint, and preserve its test
     body unchanged. Correction commit `75c856a` implements this path.
  4. Rerun the complete phase-level gate on the corrected tree. ESC-20-004
     remains open and Task 4 remains blocked until the same independent
     reviewer returns a non-BLOCK verdict.
- **Consequences:** The retained concern-write transition test becomes part of
  ordinary race/coverage, and a future tagged test anywhere outside
  `cmd/fid/integration` fails the dedicated classifier. No product code,
  public contract, VAL, dependency, or accepted gate command changes.
- **Refs:** `TRG-20-004`; `ESC-20-004`; `qa/infra-review-1.md`; DEC-133;
  correction `75c856a`; `infra-improvement-pack.md`.

### DEC-135: ESC-20-005 deletes the orphaned lock-holder formatter and re-runs Task 4 in a fresh session

- **Date:** 2026-08-11
- **Context:** Task 4's Developer implemented rebuild project selection,
  Response Envelope rendering, typed V1 authority mapping, and typed retryable
  lock handling, then the required `golangci-lint run ./...` gate failed with
  `internal/journal/lock.go:185:6: func describeLockHolder is unused (unused)`.
  The typed `LockUnavailableError` carries only `LockPath`, so the pre-existing
  holder formatter lost its only caller. `AGENTS.md` forbids repairing and
  rerunning a failed required gate in the session that failed it, so the
  session produced a blocker report with no commit, and the coverage,
  coverage-report, dedicated integration, and build gates never ran. TL-Lite
  escalated as `TRG-20-005` / `ESC-20-005`.
- **Decision:**
  1. Sustain the hard stop. A real required gate failed and no commit exists.
     This is dead code left by a contract-correct narrowing, not a product
     defect, and it is not grounds to weaken any assertion, gate, or VAL.
  2. Delete `describeLockHolder` from `internal/journal/lock.go`. `contracts.md`
     §2.4/§2.5 close `lock-unavailable` to `conflict` / `retryable: true` with
     the exact detail object `{lock_path: string}`, and §2.5 forbids
     environment values in detail objects; holder PID, hostname, command, and
     actor role therefore have no contracted destination in Phase 20.
     Suppressing the finding with `//nolint`, a blank assignment, or any other
     fake reference is prohibited. Adding holder identity to the envelope would
     be a closed-contract change requiring a spec amendment plus Human
     approval and is not authorized here.
  3. Keep `.fid/write.lock/owner.json` behavior. Writing lock metadata is
     retained operator diagnostics, `readLockMetadata` keeps its `lock_test.go`
     caller, and lint reported exactly one unused symbol, so the deletion is
     minimal and complete. If deletion cascades to further unexported symbols
     in `lock.go`, delete those too; if it reaches exported API or another
     package, stop and report a blocker instead of widening scope.
  4. Task 4 is re-run, not re-sliced: no `task-4a`, and `task-4.md` scope,
     acceptance criteria, VAL coverage, and command list are unchanged. The
     recovery session inherits the existing uncommitted working tree, makes the
     deletion, confirms the inherited tests still assert every acceptance
     criterion, runs the exact command sequence in order with captured
     verbatim tokens under the DEC-131/DEC-132 evidence rules, and produces one
     commit. Discarding the inherited tree is prohibited.
  5. `task-4-done.md` is rewritten in place as the completion report. It must
     carry a deviation note quoting the first session's failing token verbatim,
     naming the fix, and stating that the repair and the passing gates happened
     in a later session. The blocker version is preserved by Git history, the
     run ledger, and `TRG-20-005`.
  6. Human-mode rebuild failure prose is unchanged and out of scope. `contracts.md`
     §2.3 makes no Human message substring load-bearing, and §2.5 governs JSON
     detail objects only.
- **Consequences:** Task 4 costs one extra Developer session and no scope,
  contract, VAL, dependency, product behavior, or gate command changes. All
  seven `BLOCKED` rows in the blocker report (VAL-001, 002, 006, 007, 008, 017,
  and the VAL-044 rebuild row) must be re-recorded from the fresh session's own
  full-gate evidence; the earlier focused-run tokens do not carry over.
- **Refs:** `TRG-20-005`; `ESC-20-005`; `tasks/task-4.md`; `tasks/task-4-done.md`
  (blocker version); `contracts.md` §2.4-§2.5; `AGENTS.md` (Commands; Hard
  Stops); DEC-131; DEC-132.
- **Superseded in part:** decision 5's single-commit instruction is wrong and is
  replaced by DEC-136. Decisions 1-4 and 6 stand.

### DEC-136: ESC-20-006 corrects a Tech Lead commit-protocol error and re-captures Task 4's gate tokens

- **Date:** 2026-08-11
- **Context:** The Task 4 recovery session produced verified commit
  `e56b3a0ce0433f1044f70b7969c14e19234dfdf8`
  (`feat(rebuild): map typed authority and lock failures`), but
  `report-lint.sh` returned four failures against the committed report:
  `## Commit hash` and `## Commit subject` still read `Pending commit`, and
  three `go test` records carry the fragments `0.009s`,
  `coverage: 80.5% of statements`, and `4.625s`. TL-Lite escalated as
  `TRG-20-006` / `ESC-20-006`.
- **Decision:**
  1. **The pending-commit failure is a Tech Lead defect, not a Developer
     defect.** DEC-135 decision 5 ordered "exactly one commit containing the
     code, tests, and `task-4-done.md`". That contradicts the canonical
     Developer prompt's Report And Commit Sequence steps 3-6 and the protocol
     documented in `report-lint.sh` itself: a report cannot carry the hash of
     the commit that introduces its own final content, so the task commit is
     followed by a `docs(task): ...` commit that finalizes the hash-bearing
     report. The task spec binds the Developer, so it complied and stopped at
     step 2 with the placeholders intact. DEC-135 decision 5 is withdrawn and
     `tasks/task-4.md` is corrected.
  2. **Standing rule:** a task spec never restates, paraphrases, or overrides
     the commit, report, or gate-evidence sequence owned by `AGENTS.md` and the
     role prompts. Task specs carry scope, acceptance criteria, VAL coverage,
     and the exact command list only. A future conflict between a task spec and
     a role prompt on process is a Tech Lead defect to be fixed in the spec.
  3. **Sustain the three token failures.** `AGENTS.md` requires one short
     verbatim token from the command's own output — its last printed line.
     `0.009s`, `coverage: 80.5% of statements`, and `4.625s` are truncations
     that drop the `ok`/`FAIL`/`PASS` marker, so `report-lint.sh` cannot
     distinguish them from a token copied from elsewhere. The control is
     working as DEC-131 intended.
  4. **Do not weaken the control.** Relaxing the token heuristic to accept a
     bare duration or coverage fragment, or granting Task 4 an exception, is
     refused; it would reopen precisely the gap DEC-131 closed.
  5. **Re-run, do not reconstruct.** A later session may not expand a fragment
     into the full line it remembers; `AGENTS.md` classifies attestation from
     memory and transcription as `UNVERIFIED`. The repair session therefore
     re-runs the exact required command sequence itself and records its own
     captured lines. This is now cheap: DEC-133 cut the combined gate to a
     75-second median. Real re-run values will differ from the earlier
     fragments; the new values are authoritative and must never be reconciled
     to the old ones.
  6. **The evidence must apply to the reviewed commit.** Before running gates,
     the repair session verifies `git rev-parse HEAD` is `e56b3a0` and that
     `git status --porcelain` reports no source modification. If the tree is not
     identical to the reviewed commit, it stops with a blocker instead of
     recording gate results against a different state.
  7. **Task 4's commit stands.** No code change, revert, amend, or force-push.
     The repair is report-only: rewrite `task-4-done.md`, name `e56b3a0` and its
     exact subject, then land it in `docs(task): finalize completion report for
     phase-20 task-4` and rerun `report-lint.sh` against the committed file.
- **Consequences:** Task 4 costs a third session, all of it evidence capture.
  No product code, contract, VAL, dependency, gate command, or tooling
  threshold changes. The corrected `task-4.md` and the standing rule in
  decision 2 stop the same conflict from recurring in Tasks 5-18. **CF-030**
  records the optional tooling follow-up.
- **Refs:** `TRG-20-006`; `ESC-20-006`; DEC-135 (decision 5 withdrawn); DEC-131;
  DEC-133; `project-docs/agents/developer.md` (Report And Commit Sequence);
  `project-docs/tools/report-lint.sh` (header, checks 4-6); `AGENTS.md`
  (Commands: evidence tokens); commit `e56b3a0`.

### DEC-137: ESC-20-007 replaces the classification test's exact journey list with the invariant it actually protects

- **Date:** 2026-08-11
- **Context:** Task 5's implementation passed format, lint, race/coverage
  (77.7%), and every focused test, then the required dedicated integration gate
  failed: `TestIntegrationClassification` compares the dedicated package's test
  inventory against a hardcoded `wantJourneys` list, and Task 5 legitimately
  added `TestV01ParityFixtureRebuildPreservesEveryProjectionTable` to prove
  VAL-012's rebuild parity. Nothing was wrong with the product or the new
  journey; the list was stale by one name. This same test produced ESC-20-004.
- **Decision:**
  1. Sustain the gate failure as correct behavior — the list *was* stale — but
     record the assertion itself as defective. Exact equality against a
     hardcoded name list is a change-detector: it fails on legitimate growth,
     and Tasks 14-18 (especially Task 18, journey acceptance) would trip it
     again.
  2. No accepted decision requires "no new integration journey may ever exist".
     DEC-133's real invariants are that integration-tagged tests live only in
     the dedicated package, that in-process tests stay in the ordinary suite,
     and that the four journeys Task 3f migrated are not lost. Relaxing exact
     equality removes an unintended assertion, it does not weaken an accepted
     one.
  3. Change `wantJourneys` from exact equality to **required presence**: the
     four migrated journeys must still be there (a removal fails), additions
     are allowed. Keep `taggedTestFilesOutside` and its repository-wide walk
     exactly as DEC-134 left them, and keep `TestIntegrationBinaryOwnedByPackage`
     and `TestClosureJourneyHasNoWallClockSleep` unchanged. Make the failure
     message name the missing journey and the required action.
  4. Dedicated-gate runtime stays a reviewed property, not a name list: each
     report records the gate's own token, and the QA review windows read the
     diff when a journey is added.
  5. Task 5 resumes in a fresh session from the inherited tree, including the
     regenerated parity fixtures and `SHA256SUMS`, and runs the full required
     sequence. Under DEC-138 the classification edit would have been an
     in-session bookkeeping repair; the semantic change in decision 3 is a Tech
     Lead call and is authorized here.
- **Consequences:** One recurring tripwire is removed for the remaining
  thirteen tasks. No product code, contract, VAL, dependency, or gate command
  changes, and the dedicated gate keeps every assertion that protects the
  DEC-133/DEC-134 classification.
- **Refs:** `TRG-20-007`; `ESC-20-007`; `cmd/fid/integration/classification_test.go`;
  DEC-133; DEC-134; `tasks/task-5-done.md` (blocker).

### DEC-138: A stale inventory in the verification apparatus is a bookkeeping repair, not a failed gate

- **Date:** 2026-08-11
- **Context:** ESC-20-005, ESC-20-006, and ESC-20-007 are three consecutive
  escalations in which a Developer held finished, passing work and stopped
  because it was forbidden to make a one-line, self-describing correction to
  the verification apparatus: delete a function the linter named as unused, add
  a test name the classifier named as missing. Each cost a blocker report, a
  TL-Lite triage, a full Tech Lead invocation with a decision record, a Phase
  Runner stop, and a human relay, then a fresh Developer session — for edits
  smaller than the paperwork describing them. Across Phases 17-20, twenty-five
  of twenty-five escalations landed in the verification apparatus and zero in
  product behavior. DEC-100 already carved out a task-authored harness that
  fails while being authored; its exclusion of "a pre-existing committed gate or
  tool" is what these three hit.
- **Decision:**
  1. **Extend DEC-100 with a second narrow case.** A required command that fails
     only because a **bookkeeping inventory in the verification apparatus is
     stale with respect to a change the task's accepted scope legitimately made**
     may be repaired in the same session and rerun to a pass, recorded `PASS`
     with a mandatory deviation note giving the first failure's verbatim token,
     the root cause, and the fix.
  2. **All of these conditions must hold.** The failing assertion tracks an
     inventory of the repository's own artifacts — a list of test names, files,
     registered codes, or generated fixtures — and **not** product behavior. The
     failure output **names the exact discrepancy**, so the repair needs no
     judgment. The repair only reconciles the inventory with reality the task
     was authorized to create: **no threshold, budget, or tolerance is changed;
     no test or assertion is deleted, weakened, narrowed, or re-pointed; no
     product-behavior assertion is touched**. The rerun command is
     byte-identical. Both the failure and the final result are recorded with
     their own tokens.
  3. **Changing what a gate checks is still a hard stop.** Reconciling a stale
     list is bookkeeping; altering the assertion's semantics, scope, or
     strictness is a gate change and escalates. Where a task cannot tell which
     it is facing, it stops. **You may fix the bookkeeping; you may never fix
     the verdict.**
  4. **TL-Lite gains a fourth RESUME ground** with the identical test: a hard
     stop whose only unresolved element is a bookkeeping repair meeting
     decision 2 resumes with a TRG record and no escalation. Any doubt is BLOCK.
  5. **Anti-gaming is unchanged.** A failing product test, a coverage or
     performance threshold, a flaky rerun, a narrowed command, and an uncaptured
     token remain failed gates, and DEC-131's evidence-token rules are untouched.
- **Consequences:** The most common escalation class in this project's history
  resolves inside the session that hits it. Affected canonical docs: `AGENTS.md`
  (Commands), `09_AI_WORKFLOW.md` §2.3 and §8, and the canonical Developer and
  TL-Lite prompts; generated provider prompts need `sync-agents.sh` and
  `workflow-lint.sh`. The trade is a real, bounded loss of friction: a Developer
  can now edit apparatus a reviewer would previously have seen in an escalation,
  so every such repair carries a mandatory deviation note and is read at the
  next QA review window.
- **Refs:** `ESC-20-005`; `ESC-20-006`; `ESC-20-007`; DEC-100 (Phase 19, kept at
  retrospective); DEC-002; Human direction 2026-08-11 that the escalation
  overhead is disproportionate.

### DEC-139: ESC-20-008 is an unpublished evidence requirement, not a failed gate

- **Date:** 2026-08-11
- **Context:** Task 8's completion report records all six required commands
  `PASS`, each with one verbatim token from the command's own output
  (`0 issues.`; `ok  …/internal/workflow  (cached)  coverage: 80.5% of
  statements`; `total: … 77.5%`; `ok  …/cmd/fid/integration  4.971s`; and
  `exit 0 (no output)` for the two silent-success commands). `report-lint.sh`
  passed, implementation commit `4df4d66` and finalization commit `ce282c2`
  both verified. The Phase Runner nonetheless treated the gates as failed
  because `project-docs/agents/phase-runner.md` requires a completion report's
  evidence to include the nested command's `exit_code` or an explicit
  signal/termination result, and Task 8's rows carry no such annotation. Tasks
  4-7 carry `(terminated \`exit_code: 0\`)`; Task 8 followed the
  completion-report template, which has no field for it. TL-Lite returned BLOCK
  because no RESUME ground fits: DEC-138's fourth ground covers a stale
  inventory in the verification apparatus, not an evidence record.
- **Decision:**
  1. **The routing was correct and is sustained.** The Phase Runner and TL-Lite
     each applied a rule that was in force and unmet, and TL-Lite was right that
     DEC-138 does not reach this. `TRG-20-008` needed no different verdict from
     the information it had.
  2. **Task 8's report satisfies every canonical requirement in force when it
     was written**, and Task 8 is complete at `4df4d66`. All six required
     commands appear verbatim with one real token each, none of them in a
     prohibited token class, the report is committed and clean, and both commits
     verify. **No rerun, no re-recording, and no evidence-recovery session.**
     The tokens are also mutually distinct and consistent with the task's own
     diff — 77.5% total coverage against Task 7's 77.7%, integration 4.971s
     against 4.767s — so nothing independent casts doubt on their capture.
  3. **A downstream control may not demand an artifact field the producing role
     was never told to produce.** DEC-136 decision 2 forbids a task spec from
     restating process owned by `AGENTS.md` and the role prompts; the same rule
     binds controls. Where a Phase Runner, QA, or lint rule is stricter than
     `AGENTS.md`, `09_AI_WORKFLOW.md` §8, and the completion-report template,
     the canonical contract governs and the divergence is a Tech Lead defect to
     fix — not a Developer failure to correct. This is cause E of
     `docs/reviews/2026-08-11-developer-stoppage-analysis.md` recurring in a
     control instead of a task spec.
  4. **Keep the stricter rule, promote it to canon, forward-effective.** It
     carries real information. A real last-printed line proves a command
     produced output; it does not by itself prove the process exited, and for an
     `exit 0 (no output)` row the observed exit status is the *only* thing
     separating a genuine silent success from the discarded-session failure mode
     of ESC-20-003. From this decision onward, each recorded required-command
     result carries both the verbatim token and the observed termination result,
     in the shape Tasks 5-7 already use:
     `` - `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`) ``.
     Where the execution surface exposes no exit code, record the explicit
     termination signal it does expose and name it. Never infer termination from
     wrapper completion, an empty output field, or absence of new output.
  5. **Canonical docs updated now:** `AGENTS.md` (Commands),
     `09_AI_WORKFLOW.md` §8, `project-docs/templates/completion-report-template.md`,
     `project-docs/agents/developer.md`, and `project-docs/agents/phase-runner.md`
     (which now cites the canonical rule instead of holding a private one).
     Generated provider prompts require `project-docs/tools/sync-agents.sh` then
     `project-docs/tools/workflow-lint.sh`; until both run, the rule is inert.
  6. **No retroactive re-recording.** Every report written before this decision
     — Phase 20 Tasks 1-8 and all Phase 17-19 reports — is judged on the token
     contract in force when it was written and remains valid. This decision is
     not grounds to reopen any of them. DEC-131's `UNVERIFIED` lines stay
     `UNVERIFIED` for their own separate reason.
  7. **Anti-gaming is unchanged.** UNVERIFIED-not-`PASS`, DEC-131's token rules,
     `exit 0 (no output)` being invalid for any `go test`, coverage and
     performance thresholds, the ban on weakening or re-pointing tests, and the
     ban on mining a flaky green all stand. An `exit_code` asserted rather than
     observed is fabricated evidence and a hard stop, exactly like a fabricated
     token: copy it from the command's own result or record `BLOCKED`.
  8. **CF-031 opened.** `report-lint.sh` should mechanically require a
     termination result on every required-command row, with a self-test case,
     alongside CF-030's token-shape message. Until it lands the check is a Phase
     Runner read, and a missing annotation on a report written after the synced
     prompts land is a normal triage-gated hard stop.
- **Consequences:** Task 8 closes with no further work and Phase 20's eighth
  escalation costs no Developer session. The permanent cost is one clause per
  gate row, which three of the last four reports already paid voluntarily. The
  evidence contract gets stronger, not weaker: termination becomes visible in
  the artifact instead of trusted. The general lesson is recorded as decision 3
  — controls are held to the same "no private process" rule as task specs.
- **Refs:** `TRG-20-008`; `ESC-20-008`; `tasks/task-8-done.md`; `tasks/task-8.md`;
  DEC-131; DEC-132; DEC-136 (decision 2); DEC-138; `AGENTS.md` (Commands);
  `09_AI_WORKFLOW.md` §8; `project-docs/agents/phase-runner.md`;
  `docs/reviews/2026-08-11-developer-stoppage-analysis.md` (cause E).

### DEC-140: ESC-20-009 sustains all three Review 3 findings and recovers them in one bounded test-led task

- **Date:** 2026-08-12
- **Context:** Review 3 covered Tasks 7-9 and returned BLOCK with three Task 9
  Blocking findings and zero findings against Tasks 7 and 8. Independent
  verification of the tree confirms each one. `--expected-revision` is wired on
  transition (`cmd/fid/issue_transition.go:28,39,50`), but no test anywhere in
  the repository references it outside `cmd/fid/issue_update_test.go` and
  `internal/cli/codes_test.go`, so the transition side of the three-state matrix
  has no assertion at all. The only canonical `--json` transition test is
  `TestIssueTransitionJSONRejectsUnsatisfiedClaimPrecondition`; the remaining
  transition tests are legacy Human-mode substring checks, and no test compares
  a transition envelope's Event identity against the journal tail.
  `TestIssueCreateRejectsRemovedOwnerFlagsWithJSONFailure` loops both flags over
  `issue create` only. `internal/issues/update_test.go:244` asserts a projected
  activity actor and unchanged owner fields but makes no forbidden-table or
  Event-effect assertion. The Task 9 gates, report-lint, commit verification,
  and the independent integration run all pass; none of them can prove a
  conditional state that was never asserted.
- **Decision:**
  1. **All three findings are sustained and `TRG-20-009`'s BLOCK was correct.**
     None is misclassified, no verified corrective commit exists, and no accepted
     decision covers them. This is the escalation path working as designed.
  2. **Nothing is weakened.** VAL-018's three-state update *and* transition
     matrix with journal counts, VAL-019's four-way transition table with journal
     counts, VAL-040's flag-tree plus mutation-result plus retained-path
     evidence, and VAL-041's forbidden-table/journal-effect assertions stand
     exactly as Human-approved.
  3. **This is an evidence gap, not a demonstrated product defect** — and the
     distinction is the point. The implementation exists; whether it is correct
     is precisely what is unproven. Task 9a is test-led: an assertion that fails
     reveals a real defect, fixed in the same session inside the existing
     contract under red/green with a mandatory deviation note. A contract that
     looks wrong is a hard stop, not a local edit.
  4. **Task 9 is not reopened.** Its implementation commit `86d927d` stands and
     its report is not rewritten. Task 9a is additive evidence on top of it, in
     the same shape as Tasks 3a-3d. The affected VAL rows — VAL-017 transition
     portion, VAL-018, VAL-019, VAL-040, VAL-041, and the owner/transition rows
     of VAL-044 — re-record from Task 9a's own session evidence.
  5. **Responsibility: this is a Developer evidence shortfall, and DEC-139
     decision 3 does not apply.** The requirement was published in three places
     the Developer independently loads: the Task 9 acceptance criteria naming
     VAL-018/019/040/041, the validation contract's Evidence and Pass columns
     (which `project-docs/agents/developer.md` requires loading for every VAL the
     task spec names), and the phase brief's Known Risk restating DEC-130's
     no-sampling rule. Unlike ESC-20-008, no control demanded something the
     producing role was never told to produce.
  6. **Contributing Tech Lead factor, corrected forward.** Task 9 is the phase's
     widest product task: eleven VAL IDs and three commands, expressed as
     one-line acceptance criteria with a single update-only First Failing Test.
     A one-line criterion that summarizes a matrix reads like a replacement for
     it. The fix is **not** to copy each matrix into each task spec — that
     duplication is its own defect source, per DEC-136 decision 2 — but to make
     the precedence explicit. New standing constraint, added to the phase brief
     and effective for Tasks 10-18: **a task spec's acceptance criterion
     summarizes the validation contract row, it never replaces it.** Where a VAL's
     Evidence or Pass column names a matrix, table, or "every X", assert one
     observable case per cell for every named operation, and cite the
     cell-to-test mapping in the completion report. Representative or
     single-operation evidence is a failed acceptance criterion.
  7. **Recovery path.** `tasks/task-9a.md` is the single bounded recovery task.
     After it closes, the same reviewer performs a bounded re-review written to
     `qa/review-3b.md`, scoped to the three Review 3 Blocking findings and the
     Task 9a gate record only. `qa/review-3.md` stays as written. Phase Runner
     resumes at Task 9a, then review-3b, then Task 10.
  8. **Anti-gaming unchanged.** DEC-131's token rules, DEC-132's termination
     rules, DEC-139's forward-effective termination annotation, the ban on
     weakening or re-pointing tests, and UNVERIFIED-not-`PASS` all apply. New
     assertions must be observable at a real boundary — decoded canonical JSON,
     journal line counts, projection reads — never source constants or
     unexported internals that merely restate the implementation.
- **Consequences:** One extra Developer session plus one bounded re-review, and
  Task 10 starts from a mutation surface whose conditional states are actually
  proven rather than assumed. The standing constraint in decision 6 costs each
  remaining task a short cell-to-test mapping in its report and removes the
  reading under which a compressed criterion looks complete. This is the second
  occurrence of DEC-130's realized risk — a wide task producing representative
  rather than exhaustive evidence — so the correction targets that pattern, not
  this instance.
- **Refs:** `TRG-20-009`; `ESC-20-009`; `qa/review-3.md`; `tasks/task-9.md`;
  `tasks/task-9-done.md`; `tasks/task-9a.md`; DEC-130; DEC-136 (decision 2);
  DEC-139 (decision 3); `qa/validation-contract.md` VAL-017/018/019/040/041/044;
  `contracts.md` §2.4, §3.2, §3.6.

### DEC-141: Recover Task 10 under the typed workflow-repair model and close ESC-20-010

- **Date:** 2026-08-15
- **Context:** The Task 10 Developer session ended normally after 616 seconds;
  no timeout, cancellation, kill, or Phase Runner interrupt occurred. It left a
  durable partial migration and a focused package failure caused by remaining
  contract-obsolete tests, then classified that incomplete in-scope state as a
  blocker. `TRG-20-010` and `ESC-20-010` faithfully applied the old binary
  `RESUME | BLOCK` process, but no product, architecture, data-model,
  public-interface, dependency, security, specification, or external
  infrastructure decision was required. The Human authorized an isolated
  principal-engineer recovery and a systemic workflow-reliability overhaul.
- **Decision:**
  1. `ESC-20-010` is reclassified as a decision-free same-task recovery. The
     original red suite and incomplete report remain truthful evidence, but
     neither condition requires Human attention under the accepted
     `RESUME | REPAIR | BLOCK` model.
  2. Recovery uses the preserved stash `779beab1` on the isolated
     `workflow-reliability-overhaul` branch. The v2-marked run ledger records
     `RUN-RESUME`, `TASK-RECOVERY`, `REPAIR`, and Task 10 completion in order;
     failed attempts are retained rather than rewritten.
  3. Task 10 remains one task. Its accepted contract and VAL rows are unchanged.
     Recovery completes the six operation transports, exact response/error
     shapes, aggregate revisions, Event identities and golden replay, closed
     value domains, deterministic ordering, collection limits, selected-project
     routing, Human channels, and the guarded legacy `issue ac verify` leaf.
  4. The first restored focused failures are causal red/green evidence. Every
     obsolete assertion is migrated to the accepted contract; no test,
     threshold, command, or acceptance criterion is weakened or narrowed.
  5. The repository workflow overhaul is accepted as the v1 reliability patch:
     typed report outcomes, protected blocker categories, same-task continuation,
     strict recovery-ledger validation, bounded QA repair/re-review, exact PASS
     termination parsing, anti-churn progress accounting, and retained polling
     of silent live commands. Independent review at `06a413a` returned PASS
     after all Must-fix findings were repaired.
  6. Broader test-performance refactors and edits to the separate v2 starter kit
     remain out of this recovery. The audit's v2 adoption deltas are recorded for
     a later isolated workstream.
- **Consequences:** Task 10 may close after its exact six-command gate,
  candidate commit verification, report finalization, and report-lint pass.
  `ESC-20-010` is resolved without a product exception and Phase Runner may
  continue at Task 11 in a fresh session. Future ordinary red tests, in-scope
  gate repairs, report defects, bounded QA rework, and session continuation use
  `REPAIR`; only protected or budget-exhausted conditions produce a STOP.
- **Refs:** `TRG-20-010`; `ESC-20-010`; `tasks/task-10.md`;
  `tasks/task-10-done.md`; `run-ledger.md`; `contracts.md` §§2.4, 3, 7;
  `qa/validation-contract.md` VAL-008/010/011/017/021/022/041/044/045;
  `docs/reviews/2026-08-15-workflow-reliability-overhaul-audit.md`;
  `docs/reviews/2026-08-15-workflow-reliability-overhaul-plan.md`.

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
