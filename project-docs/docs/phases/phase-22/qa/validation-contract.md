# Phase 22 Validation Contract - Failure And Recovery

> Author: Tech Lead (draft), 2026-09-05.
> Status: DRAFT. No outcomes have been executed or attested.
> Decision: DEC-180 (proposed slice); task ownership: DEC-181.

## Purpose

This finite contract states what the phase must prove. QA Lead authors the
evidence strategy and proposes refinements; Tech Lead owns changes here.
One role validates each row. Developers produce evidence but do not declare
Human or independent QA outcomes. A missing required proof is never PASS.

Source abbreviations name the canonical specs: PRD = 01_PRD.md,
TRD = 02_TRD.md, SCHEMA = 03_TRD_SCHEMA.md, CONFIG = 04_TRD_CONFIG.md,
DESIGN = 05_DESIGN_PLAN.md, CLI = 06_DESIGN_PLAN_CLI.md,
FLOWS = 07_DESIGN_PLAN_FLOWS.md, QUALITY = 08_QUALITY_PLAN.md.
Phase-local contracts/public-interface/schema references are normative only
after Human acceptance. The kit source is pinned in qa/qa-lead-handoff.md.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | SCHEMA §7; PRD-OPS-001 | Fresh init installs the recovery schema and enforces its identity constraints. | SQLite migration tests and built-binary init | QA Lead | task-1 / window-1 | All declared tables/constraints present at schema 15; invalid duplicate identities reject. | Missing constraint, unusable init, or unaccepted schema. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-002 | SCHEMA §7 | Disposable full down/up restores the same accepted empty baseline. | Migration reset tests | QA Lead | task-1 / window-1 | Down to zero and up reproduce schema and role-free bootstrap. | Partial/failed reset or retired authority reappears. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-003 | PRD-REC-004; TRD §§6-7 | The same closed Event semantics govern live apply and replay. | Per-family positive/corrupt Event fixtures | QA Lead | task-1 / window-1 | Accepted Effects reproduce state; unknown or invalid Effects reject atomically. | Replay-only bypass, partial apply, or divergent state. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-004 | PRD-RUN-012 | A Stop's original fact remains immutable after resolution. | Domain history and replay tests | QA Lead | task-2 / window-1 | Original ID/details unchanged; separate resolution references qualifying fact. | Edited/deleted fact or prose-only resolution. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-005 | PRD-RUN-013 | Equivalent Attention requests share one open identity without losing distinct Stops. | Same/different semantic-key concurrency cases | QA Lead | task-2 / window-1 | One open item per semantic key; all distinct Stops remain linked. | Duplicate open items or lost distinct Stop. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-006 | PRD-HUM-003 | Attention resolution requires qualifying facts for every linked Stop. | Valid/wrong-subject/missing-reference cases | QA Lead | task-2 / window-1 | Valid facts resolve once; invalid facts have no Effects. | Dismissal or unqualified resolution accepted. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-007 | PRD-HUM-002 | Routine status and warnings do not create Attention. | Status/warning versus policy Stop cases | QA Lead | task-2 / window-1 | Routine queries/warnings leave count unchanged; real Stop raises Attention. | Routine interruption or missing real intervention. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-008 | PRD-RUN-004 | Lease validity changes exactly at expiry and cannot be revived by renewal. | Clock-controlled before/equal/after tests | QA Lead | task-3 / window-1 | Unexpired current version works; equal/expired/released rejects renewal. | Boundary grants expired authority. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-009 | PRD-RUN-002; PRD-OPS-001 | Concurrent same-Issue mutation acquisition has at most one winner. | Two-process Claim race | QA Lead | task-3 / window-1 | One active mutating Claim; loser returns typed conflict with no extra Event. | Multiple owners or partial losing write. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-010 | PRD-RUN-008; PRD-AUTH-001 | Only a trusted bound observer can establish executor liveness. | Activated harness versus Builder/forged-ID cases | QA Lead | task-3 / window-1 | Accepted observation binds exact Run/executor; other actors reject. | Self-declared termination grants authority. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-011 | PRD-RUN-001/004 | Claim renewal preserves Run lineage while fencing obsolete lease versions. | Renew then mutate with current/stale version | QA Lead | task-3 / window-1 | Run identity stable; current unexpired lease authorizes; stale version rejects. | Renewal bypass or stale version works. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-012 | PRD-RUN-005 | A Run binds exact executor, optional Obligation, and originating authority context. | Start binding table tests | QA Lead | task-3 / window-1 | Matching binding accepted, mismatched Issue/revision/Claim/Obligation rejects. | Unbound or mismatched Run accepted. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-013 | PRD-RUN-008 | Clean unchanged repository state qualifies for ordinary reconciliation. | Git fixture with unchanged tree/checkpoint | QA Lead | task-4 / window-2 | Complete deterministic observation qualifies. | False conflict or missing observation dimensions. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-014 | PRD-RUN-011 | Dirty in-scope work requires exact recorded adoption. | Adopt/no-adopt/changed-digest cases | QA Lead | task-4 / window-2 | Exact set accepted with provenance; omitted or stale set rejects. | Implicit adoption or forged digest accepted. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-015 | PRD-RUN-009 | Unsafe or ambiguous repository state blocks recovery. | Unmerged/history/path/worktree/operation-marker cases | QA Lead | task-4 / window-2 | Every unsafe state returns recovery-state-conflict without takeover. | Unsafe state proceeds or silently discarded detail. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-016 | PRD-RUN-010 | Recovery never changes source files or Git state. | Before/after refs/index/stash/file digests on success and failure | QA Lead | task-4 / window-2 | Source bytes and Git refs/index/stash unchanged; only normal fid operational metadata may change. | Reset/clean/stash/checkout/commit/discard or other mutation. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-017 | PRD-RUN-007 | Recovery creates a distinct linked successor and keeps the predecessor terminal. | Domain plus real-process successor chain | QA Lead | task-5 / window-2 | One new ID, exact predecessor/successor links, predecessor interrupted. | Old Run reused/reopened or duplicate successor. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-018 | PRD-OPS-001/002; PRD-RUN-007 | Recovery is atomic and idempotent across response loss and projection interruption. | Fault injection after append and competing-key calls | QA Lead | task-5 / window-2 | One full recovery transaction; same-key replay original; different digest conflicts. | Partial facts or duplicated mutation. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-019 | PRD-RUN-008/009; kit INTERFACES §3.2 | Unknown liveness always blocks ordinary and force recovery. | Closed liveness matrix | QA Lead | task-5 / window-2 | Terminated ordinary can proceed; unknown always executor-liveness-unknown. | Lease expiry or force bypasses unknown state. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-020 | PRD-AUTH-005; PRD-RUN-008; kit INTERFACES §3.2 | Known-live takeover requires exact current Human force authority. | Grant/type/reference mismatch matrix | QA Lead | task-5 / window-2 | Only bound Human force qualifies; historical live fact remains and old fid writes fence. | Builder/forged/stale authority succeeds. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-021 | SCHEMA §4; PRD-RUN-005/007 | Recovery preserves the Issue and immutable Envelope execution contract. | Before/after semantic snapshot | QA Lead | task-5 / window-2 | Issue/revision/Obligations/scope/policy/evidence identities unchanged. | Recovery resets or revises durable work. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-022 | Roadmap §5 Phase 22; DEC-181 | A built-binary recovery journey works by task 6. | Tagged subprocess journey | QA Lead | task-6 / window-2 | Initialize, activate, start, observe termination, expire/release, reconcile, recover, inspect succeed. | In-process substitute, stub, or unexercised recovery. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-023 | PRD-OPS-008; CLI §§1-4 | Recovery leaves enforce their closed wire schemas. | Per-leaf accepted/unknown-field request/result fixtures | QA Lead | task-6 / window-2 | Every section-3 recovery leaf matches shape, channel, status, and Event. | Unknown field accepted or missing leaf/result field. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-024 | PRD-OPS-008; CLI §5 | Recovery failures have deterministic typed process exits. | Real-process failure corpus | QA Lead | task-6 / window-2 | Not-found/authority/liveness/revision/conflict map to exact codes and exits. | Human text parsed or incorrect/zero failure exit. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-025 | PRD-RUN-006; contracts §6 | fid-owned resource units are charged only from observed execution facts. | Automatic-charge versus caller-submission cases | QA Lead | task-7 / window-3 | Exactly matching Attempts/Receipts/force events charge; caller substitution rejects. | Missing/double/invented charge. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-026 | PRD-RUN-006; PRD-OPS-002 | Budget samples are idempotent and ordered across recovery. | Same sample/different key/stale baseline cases | QA Lead | task-7 / window-3 | Same source sample counts once; invalid predecessor/baseline conflicts. | Replayed sample adds consumption or total decreases. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-027 | PRD-TELEM-001; PRD-RUN-006 | Unknown budget measurements remain null and block admission. | Unknown, zero, numeric baseline cases | QA Lead | task-7 / window-3 | Unknown differs from zero; qualifying baseline resolves data absence only. | Unknown silently zero or work starts with missing required data. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-028 | PRD-RUN-006 | Budget exhaustion commits its charge, Stop, and Attention atomically. | Exact-boundary plus append/projection fault cases | QA Lead | task-7 / window-3 | Last admission charged once; Stop/Attention co-commit; later spawn denied. | Charge without Stop/Attention, over-budget new spawn, or partial commit. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-029 | PRD-RUN-006; SCHEMA §4 | Consumed budgets survive successor Runs and rebuild. | Linked recovery and replay budget snapshots | QA Lead | task-7 / window-3 | Exact totals/sample chain retained; maximum unchanged. | Reset counter or new recovery allowance. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-030 | PRD-RUN-006; PRD-HUM-002 | Budget warning and exhaustion states have distinct behavior. | Below/equal-warning/equal-maximum cases | QA Lead | task-7 / window-3 | Warning is queryable; exhaustion blocks configured further work. | Warning blocks or maximum is ignored. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-031 | PRD-GATE-004/005 | Observed unsuccessful termination cannot produce a passing Receipt. | Cancel/timeout/signal/spawn/harness fault tests | QA Lead | task-8 / window-3 | One observed non-success Receipt, or no Receipt if observation lost. | Synthetic pass, duplicate Receipt, or outcome mislabeled. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-032 | PRD-GATE-005/009 | An abandoned Attempt cannot be replayed into an automatic second spawn. | Kill executor parent and retry original key | QA Lead | task-8 / window-3 | Original Attempt visible, no invented Receipt, same-key retry never spawns. | Retry reruns process or supplies result. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-033 | PRD-GATE-009; PRD-EVID-006 | Run success is derived from qualifying trusted facts. | Satisfied versus missing/stale/failed proof cases | QA Lead | task-8 / window-3 | Only qualifying completion derives success; failed/missing states cannot. | Caller result accepted or success inferred from absence. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-034 | PRD-RUN-007/008; Glossary Interrupted Run | Intentional abandonment is distinct from interruption and proof. | Authorized abandon versus unknown/unauthorized cases | QA Lead | task-8 / window-3 | Abandoned Run releases Claim with provenance and no new Evidence. | Abandon produces success evidence or masquerades as interruption. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-035 | CF-034 | Reopened evidence recovery branches have direct package tests and measured coverage. | Before/after go test -race -cover ./internal/evidence outputs | QA Lead | task-7 / window-3 | Both coverage values and direct positive/rejection tests reported. | Missing baseline/final measurement or only command-layer tests. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-036 | PRD-RUN-012/013; CLI §3 | Stop, Attention, budget, and closure leaves implement their declared contracts. | Per-leaf JSON and Human text cases | QA Lead | task-9 / window-3 | Every new leaf works and rejects invalid inputs as specified. | Missing leaf, undeclared option, or divergent representation. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-037 | PRD-EVID-006; PRD-RUN-012 | All guarded actions apply the same scoped Stop predicates. | Evaluator/Issue/Envelope/execute cross-surface cases | QA Lead | task-9 / window-3 | Applicable Stop blocks; unrelated Stop does not; all reasons returned. | Overblocking or bypass via another adapter. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-038 | PRD-OPS-008; contracts §6 | A committed execution denial names its durable effects without implying execution success. | Live/JSON/exit/idempotent replay fixtures | QA Lead | task-9 / window-3 | Nonzero failure names Event/Stops; no start/output frame or spawn. | Hidden partial effects, zero exit, or replay duplicates Stop. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-039 | Roadmap §5 removal; PRD-GATE-009 | Caller-reported run.finish is absent after derived closure lands. | Command tree and current Event fixture inventory tests | QA Lead | task-9 / window-3 | finish rejects as absent; close/abandon remain; no current run.finished writer. | Legacy outcome path or compatibility alias remains. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-040 | PRD-OPS-004 | Activation prepare records the exact old and new resource sets. | First activation and upgrade cases | QA Lead | task-10 / window-4 | Five resources and immutable internal components bind one predecessor. | Missing old identity or mixed mutable component set. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-041 | PRD-OPS-006 | Activation recovery direction respects the point of no return. | Before/after commit direction matrix | QA Lead | task-10 / window-4 | Restore prior before commit; complete forward after; wrong direction rejects. | Post-commit rollback or new recovery identity. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-042 | PRD-OPS-004/005 | Activation cannot finalize until every required resource agrees. | All-agree versus each missing/mismatched resource | QA Lead | task-10 / window-4 | Only complete exact agreement reaches finalized trust. | Partial activation accepted. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-043 | PRD-OPS-005; PRD-AUTH-001 | Incomplete activation blocks normal work and preserves recovery inspection. | Normal mutations versus allowed inspection/recovery matrix | QA Lead | task-10 / window-4 | Normal operations reject; inspection and recovery remain usable. | Mixed authority usable or repair inspection unavailable. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-044 | PRD-OPS-002/004 | Activation recovery retries preserve one durable direction and identity. | Replay/conflicting-key/failure-after-append cases | QA Lead | task-10 / window-4 | Same recovery identity; stored direction immutable; no duplicate terminal Effects. | New activation or conflicting direction on retry. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-045 | PRD-OPS-008; CLI §3 | Activation recovery commands expose closed results and typed direction errors. | CLI schema and exit corpus | QA Lead | task-11 / window-4 | Every new/amended leaf matches section 3/4; old unavailability error removed. | Undefined shape or obsolete error used. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-046 | PRD-OPS-004/006; FLOWS §2 | Activation recovers at every durable boundary through the built binary. | Fault-point process journey | QA Lead | task-11 / window-4 | Prepare/stage/commit/finalize interruptions recover in required direction. | Only happy path exercised or trust disagreement persists. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-047 | CF-035; CONFIG §2; contracts §12 | Whole-config validation occurs before a consuming operation uses values. | Unknown-key consuming/non-consuming paired cases | QA Lead | task-11 / window-4 | Consuming command rejects unknown key; help/archive validation remains available. | Unknown key used or global parser blocks diagnostics. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-048 | PRD-EVID-003; DEC-185; contracts §3.1 | A changed repository target stales prior repository-bound proof. | Changed tree and real same-HEAD one-byte dirty change using canonical vectors | QA Lead | task-12 / window-4 | Same Evidence identity projects stale against the changed observed target. | Old proof qualifies on changed target or fixture bypasses actual observation. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-049 | PRD-EVID-003; DEC-185; contracts §3.1 | An unchanged target retains proof through recovery. | Real clean and dirty targets through recovery/branch rename | QA Lead | task-12 / window-4 | Proof remains valid when canonical observed identities match in both clean and dirty cases. | Recovery alone stales or refreshes proof, or matching synthetic digests replace observation. | Draft | No conversion of historical nonconforming digests. |
| VAL-050 | PRD-EVID-002/004/005 | Artifact loss preserves Evidence identity while changing only availability and policy outcome. | Required-material versus metadata-sufficient cases | QA Lead | task-12 / window-4 | Digest/provenance retained; outcome follows explicit Artifact policy. | Missing artifact erases proof or availability conflated with validity. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-051 | PRD-RUN-014; PRD-PROJ-002/003 | Fresh resume reconstructs bounded recovery context without prior chat. | New process after interruption and rebuild | QA Lead | task-12 / window-4 | All specified context dimensions present from durable facts. | Stored next-action state, omitted Stops/budgets, or prose dependency. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-052 | PRD-RUN-014; PRD-PROJ-001 | Resume action eligibility matches actual mutation predicates. | Eligible/ineligible action matrix | QA Lead | task-12 / window-4 | Advertised executable actions succeed at matching revision; blocked actions carry reasons. | Impossible action advertised as executable. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-053 | PRD-PROJ-005; CLI §7 | Human and JSON recovery views use the same semantic projection. | Paired render fixtures | QA Lead | task-18 / window-6 | IDs, status, reasons, budgets, and uncertainty agree across renderings. | Human rendering invents different state. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-054 | PRD-REC-001 | Manifest v2 contains all facts needed for recovery without project prose or SQLite. | Archive entry and source-reference corpus | QA Lead | task-13 / window-5 | All accepted recovery identities/history/config present; forbidden content absent. | Missing reconstruction data or raw project/DB copy. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-055 | PRD-REC-001; FLOWS §10 | Valid hashes cannot conceal invalid recovery history or references. | Integrity-valid but semantic-invalid archive cases | QA Lead | task-13 / window-5 | Unknown Effects/schema, duplicate IDs, broken references/order reject. | Malformed history validates or mutates state. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-056 | PRD-OPS-007; TRD §8; contracts §10; DEC-186 | Archive validation obeys every case in the explicit boundary acceptance matrix below. | Named per-boundary cases 056-a..k, measured decompression and filesystem snapshots | QA Lead | task-13 / window-5 | Every permitted boundary passes its limit check; every excess/unsafe case rejects before publication; cumulative cap uses measured expanded bytes. | Missing matrix row, unbounded expansion, unsafe entry accepted, or only declared sizes checked. | Draft | No forged ZIP metadata bypass required; all resource ceilings unchanged. |
| VAL-057 | PRD-REC-002 | Archive validation has no project or destination effects. | Before/after file/journal/projection snapshots | QA Lead | task-13 / window-5 | Valid and invalid validation leave all state unchanged. | Validation imports, repairs, writes, or silently catches up project. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-058 | PRD-REC-003 | Import preserves stable entity and Event identity/history. | Ordered source/import identity comparison | QA Lead | task-14 / window-5 | IDs, Event bytes/order/hashes, idempotency history identical. | Identity remap or rewritten/dropped history. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-059 | PRD-REC-004; PRD-OPS-001 | Import publishes all reconstructed state or no target project. | Fault injection at staging/fsync/rename boundaries | QA Lead | task-14 / window-5 | No partial .fid before publication; complete recoverable state after. | Visible partial installation or overwritten project. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-060 | PRD-OPS-002; PRD-REC-003 | Import response-loss retry does not duplicate facts. | Same and changed digest/key process retries | QA Lead | task-14 / window-5 | Same installed import replays; changed request conflicts. | Duplicate Events or second installation. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-061 | PRD-REC-003; contracts §10 | Import cannot overwrite or merge a populated target. | Empty/populated/symlinked target cases | QA Lead | task-14 / window-5 | Only safe empty target succeeds; others preserve all existing bytes. | Overwrite, implicit merge, or unsafe target accepted. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-062 | PRD-RUN-008; PRD-REC-004 | Imported active work does not automatically regain mutation authority. | Import active Run then attempt mutation/recovery | QA Lead | task-14 / window-5 | Requires current liveness and reconciliation; no spawn on import. | Historical lease/observation alone grants local takeover. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-063 | PRD-REC-004 | Live, rebuilt, and imported current state are semantically equivalent. | Named JSON projection three-way comparison | QA Lead | task-15 / window-5 | All promised dimensions agree under declared environmental normalization. | Unexplained difference or omitted dimension. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-064 | PRD-REC-001/004 | Re-export of reconstructed state preserves canonical recovery content. | Source/rebuilt/imported archive comparison | QA Lead | task-15 / window-5 | Manifest and entry bytes agree except declared path/publication exclusions. | Nondeterministic identity/order/content. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-065 | Roadmap §5 Phase 22; PRD-REC-004 | The built CLI completes export, validate, import, rebuild, and resume. | Tagged multi-repository journey | QA Lead | task-15 / window-5 | Same interrupted recovery state survives all steps. | Library-only substitute or last-step failure. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-066 | Roadmap §5 removal; SCHEMA §6 | Latest schema has no escalation entity or foreign-key references. | Latest-schema/catalog/command tests | QA Lead | task-16 / window-6 | Escalation table/seed/refs/commands absent; historical migration fixtures classified. | Retired live schema or command remains. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-067 | Roadmap §5 removal; PRD-PROJ-002 | Recovery/activity expose one Event-derived source of truth. | Writer/reader inventory and replay tests | QA Lead | task-16 / window-6 | Run/Stop/Attention facts drive recovery; activity is inspection only. | Caller-reported/duplicate recovery state or mutable activity authority. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-068 | QUALITY §3; DEC-167/173 | Retained assertions displaced by removal converge under executing gates. | Per-row provenance and ordinary/tagged test evidence | QA Lead | task-16 / window-6 | Introducing commit, invalidating task, changed row, and executing command recorded. | Filename ownership, skipped/weak test, or unexecutable assignment. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-069 | PRD-OPS-003 | Maintenance detects recovery corruption without repairing it. | Corrupt/healthy fixture matrix | QA Lead | task-17 / window-6 | Broken links, abandoned Attempts, stale Claims, budgets, activation, and missing artifacts classified; healthy controls clean. | Missed class, false healthy finding, or mutation. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-070 | TRD §3; Roadmap §5 conditional availability | Only fully qualified capabilities are advertised available. | Complete final discovery/export/fake matrix | QA Lead | task-17 / window-6 | Three completed capabilities true; Claims false with Phase 24 limit/owner; other rows unchanged. | Premature capability or missing negative metadata. | Draft | Parallel Claims qualify in Phase 24. |
| VAL-071 | CF-024; TRD §9; Phase 21 VAL-083 | 100k Event rebuild meets the unchanged reference target in three fresh samples. | Exact three-process self-asserting benchmark | QA Lead | task-17 / window-6 | All three preserve parity and self-assert <30 seconds on recorded instance. | Any slow/failed sample or unavailable reference instance. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-072 | SCHEMA §7; PRD-REC-004 | Final schema 16 initializes and rebuilds without retired authority. | Fresh init/reset/rebuild process tests | QA Lead | task-17 / window-6 | Final schema/current bootstrap/capability paths agree. | Schema-14 assumption, missing migration, or legacy bootstrap. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-073 | Roadmap §7; TRD §3; kit backend-port.v1 | One shared semantic case set passes against fake and fid recovery adapters. | Tagged dual-adapter conformance | QA Lead | task-18 / window-6 | Same inputs/expected semantic outputs and negative cases; no fake backend access. | Separate weakened expectations or production fake fallback. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-074 | Roadmap §5 Phase 22 | The full one-Issue failure/recovery journey completes through public operations. | Tagged final journey with failure injection | QA Lead | task-18 / window-6 | Interrupt, reconcile, linked retry, Stop, fresh proof, export/import, rebuild, resume are coherent. | Disconnected demos or skipped failed step. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-075 | PRD-OPS-008; CLI contract amendment; DEC-184 | Every amended wire variant has a non-vacuous deterministic conformance fixture. | Operation/error/frame inventory and explicit 23-action set compared with independently pinned corpus members and actual requests | QA Lead | task-18 / window-6 | No missing leaf, closed value, error, replay or live-frame variant; retired/unknown actions reject. | Empty corpus, self-derived expectation, uncovered variant, or mismatching fixture. | Draft | Activation operations are bootstrap operations, not activated grant actions. |
| VAL-076 | DESIGN §§4/6; phase brief | A Human can choose the next recovery action from a fresh resume view. | Human recovery checkpoint | Human | task-18 / window-6 | Human records a clear actionable choice without prior chat. | Human cannot determine action or mistakes adoption for proof. | Draft | Developer supplies fixture only. |
| VAL-077 | PRD-OPS-006; phase brief | A Human can distinguish restore-prior from complete-forward recovery. | Human paired activation checkpoints | Human | task-18 / window-6 | Direction and irreversible commit consequence are understood. | Human believes post-commit rollback occurred. | Draft | Developer supplies fixture only. |
| VAL-078 | PRD-REC-003; phase brief | A Human can recognize restored work and its remaining intervention need. | Human imported-project checkpoint | Human | task-18 / window-6 | Identity and remaining work are recognizable from restored views. | Restored project appears complete or unrelated incorrectly. | Draft | Developer supplies fixture only. |
| VAL-079 | CF-034 | Reopened workflow recovery branches have direct package tests and measured coverage. | Before/after go test -race -cover ./internal/workflow outputs | QA Lead | task-2 / window-1 | Both measured values and direct positive/rejection tests reported. | Missing baseline/final coverage or only CLI tests. | Draft | Phase 22 one-Issue scope; see exclusions below. |
| VAL-080 | Roadmap §7; QUALITY §8 | Independent review finds no unresolved cross-product recovery contract mismatch. | QA comparison with pinned kit interface/decisions | QA Lead | task-18 / window-6 | Exact semantic mapping and any accepted differences recorded before execution; final changes rechecked. | Unknown source version or unresolved contract difference. | Draft | Independent QA; TL planning is not its evidence. |

## VAL-056 Boundary Acceptance Matrix

VAL-056's single observable is conformance to this finite acceptance matrix,
following the same matrix treatment as VAL-024/069/075. Each case has separate
evidence; one success cannot stand in for the others. Sizes use binary units.
For a numeric ceiling L, exercise L-1, L, and L+1. Permitted means that limit
check accepts; full archive acceptance still requires every semantic check.
Use valid complete controls as well, and identify which boundary each negative
actually reaches so an earlier unrelated rejection cannot stand in for proof.

| Case | Boundary | Required result |
|---|---|---|
| 056-a | 1000000 entries | Count at/below ceiling accepted; next entry rejected. |
| 056-b | 64 MiB per non-Artifact entry | Actual streamed bytes at/below accepted; excess rejected. |
| 056-c | 1 GiB per Artifact | Actual streamed bytes at/below accepted; excess rejected. |
| 056-d | 8 GiB cumulative expanded bytes | Many individually permitted entries reach the total; measured decompressed byte L+1 trips the cap. |
| 056-e | 1 MiB per Event line | Exact line boundary accepted; overlong line rejected even across read/chunk boundaries. |
| 056-f | Relative normalized safe path versus traversal | Safe file accepted; parent traversal and escaping components rejected. |
| 056-g | Relative versus absolute path | Safe relative file accepted; absolute/drive-rooted paths rejected. |
| 056-h | Regular file versus symlink | Regular file accepted; symlink entry rejected without following it. |
| 056-i | Regular file versus device | Regular file accepted; device entry rejected. |
| 056-j | File versus explicit directory entry | Regular file accepted; directory entry rejected. |
| 056-k | Unique versus duplicate entry name | Unique accepted; duplicate rejected before extraction/publication. |

Generate large compressed/count fixtures per run, never commit them or
materialize expanded gigabytes. Stream while counting/discarding; for the
cumulative case keep individual entries within their own limits. Metadata
limits may reject early, but evidence must also exercise the measured-byte
limiter on actual decompression. A malformed raw ZIP rejected by archive/zip
before output is corruption evidence only, not measured-cap evidence.
Use one filesystem for import staging and its target; avoid RAM-backed /tmp
for fixtures that truly require large materialized files.

## Refinement Evidence Obligations

- VAL-001/003: non-null recovery observation FK and equal Stop/Attention code
  CHECKs in the amended exact DDL, including corrupt/replay rejection.
- VAL-006: exact disposition-containing Event set, underlying facts revalidated;
  missing/extra/duplicate/wrong-subject/underlying-only requests reject.
- VAL-012/023: reused executor identity has the declared duplicate-entity shape;
  same-key replay wins before conflict checking, with no extra effects.
- VAL-014: task 4's canonical byte/hash vectors bind staged/unstaged/mode/rename/
  deletion/untracked cases and exclusions; no observer-derived expectation.
- VAL-019/020/034: null/missing/nonexistent/wrong-subject observation references
  reject, terminated ordinary qualifies, known-live requires force, unknown
  always rejects. No rejection creates a Stop as a side effect.
- VAL-031: exact signal/timeout/cancelled outcomes with synchronized real
  subprocesses under bounded load; DEC-186 assigns the known fixture repair.
- VAL-036/037/039: distinguish explicit Stop facts from non-mutating recovery
  errors, enforce the action table, and remove run.finish leaf/action/Event.
- VAL-045: activation-version-conflict, competing-prepare activation-conflict,
  and activation-direction-conflict each have one closed details shape.
- VAL-071: task 1/16 single early samples supplement, never replace, task 17
  and final QA's three fresh samples at the unchanged reference/threshold.
- VAL-080: audit every C-01..C-13 disposition in kit-contract-mapping.md against
  as-built behavior and current pin; a documented Phase 22 exclusion does not
  prove the corresponding production kit obligation is implemented.

## Requirement Coverage And Explicit Exclusions

- PRD-RUN-001/002/004..014 map to VAL-004..034, VAL-051/052 and VAL-062.
  PRD-RUN-003 parallel acceptance stays Phase 24; same-Issue exclusion is
  VAL-009 and unknown/overlapping scope still fails closed.
- PRD-GATE-004/005/009 map to VAL-031..034/039; preserved immutable registration,
  redaction, streaming and Receipt bindings retain Phase 21 regression coverage.
  PRD-GATE-008 external attestation is not introduced.
- PRD-EVID-002..006 map to VAL-033/037/048..050. PRD-EVID-003's external-artifact
  exception remains deferred because this phase enables no such target type.
- PRD-HUM-002/003 map to VAL-005..007. Protected Human Outcome and Review
  implementation remain Phase 23; operational fact resolution is in scope.
- DEC-187: Review budget charging/review_cycles and review-rework-exhausted are
  Phase 23; generic retry-exhausted is Phase 24. circuit-breaker classification
  belongs to the kit coordinator, not fid Attention in this slice.
- DEC-187/CF-037: kit recovery-export binding, scope/asOf mapping, unresolved
  kit code representations and terminalState adapter expectation checking are
  Phase 25 qualification obligations. The fid semantics and disclosed limits
  are in scope for VAL-080; the production adapter is not.
- PRD-REC-001..004 map to VAL-054..065. PRD-PROJ-001..003/005 are covered only
  for recovery/resume/Attention by VAL-051..053; other named projections and
  trajectory summaries retain their roadmap owners.
- PRD-OPS-001..006/008 map to VAL-003/018/024/038/040..047/059/060/069.
  PRD-OPS-007 and TRD section 8 archive trust boundaries map to VAL-056.
  Existing sensitive-output protection remains a regression obligation.
- PRD-AUTH-001/005 force/observer authority map to VAL-010/020/043.
  PRD-TELEM-001 is covered only for unknown and observed budget metrics by
  VAL-025..027; telemetry capability remains Phase 25.
- Schema/Events, CLI amendment, flow sections 2/6/10, legacy removal, and
  CF-024/034/035 have explicit rows. No project-prose, role-policy, Review,
  multi-Issue revision, first-adoption or held-out-scenario promise is implied.

## Coverage Rules

All 80 IDs are referenced in an owning task's acceptance criteria.
Each has exactly one validator: 77 QA Lead assertions and three Human
assertions. QA Reviewer examines three-task windows without replacing the
final validator. QA Lead must map every ID into its Ready QA plan.
Conditional rows require both qualifying and rejecting cases.
PARITY requires deterministic structural comparison, never a prose assertion.

## Outcomes

No assertion has been run in this Tech Lead planning session. Keep Draft until
the QA plan and Human acceptance make the pack Ready. Final review records
PASS, FAIL, BLOCKED, or CARRY-FORWARD for each original assertion with evidence.
Do not weaken an assertion to accommodate an implementation or fixture failure.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-09-05 | VAL-001..080 | Drafted with task/evidence ownership | Phase 22 kickoff | DEC-180/181 |
| 2026-09-05 | VAL-001/003/006/012/014/019/020/023/031/034/036/037/039/045/048/049/056/071/075/080 | Pin exact evidence cases and accepted differences; enumerate VAL-056 without changing ceilings or adding/removing IDs | QA R-01..14 and C-01..13 | DEC-184..187 |
