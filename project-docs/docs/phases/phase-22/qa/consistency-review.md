# Phase 22 Independent Cross-Product Consistency Review

> QA Lead (Claude Opus 5), 2026-09-05. Input to `qa-plan.md`; covers VAL-080.
> This review is independent evidence. The Tech Lead's own planning comparison
> was read as an input, not relied on as proof.

## 1. Sources And Version Currency

| Source | Exact reference | How it was read |
|---|---|---|
| Kit repository | `/home/mark/projects/agentic-coding-project-starter-kit` | read-only, normal specification material only |
| Kit pinned commit | `c5229617220f97ff9228df261e24718b1d1d6d89` (`docs(qa): close v2-2 correction`, 2026-08-16) | `git log -1 --format='%H %ad %s'` |
| Kit port contract | `project-docs/specs/trd/02_INTERFACES.md` sections 3.1-3.7 and 4/4.1 | direct read |
| Kit compatibility matrix | `project-docs/specs/FID_COMPATIBILITY.md` (observed 2026-08-16) | direct read |
| Kit decision | `project-docs/specs/decisions/0005-linked-run-interruption-recovery.md` | direct read |
| fid proposal | `contracts.md`, `public-interface.md`, `schema.md` under `phase-22/` | direct read |
| fid running baseline | `/tmp/fid` built from planning base `89d964a` | `fid --json capabilities`, `--help` trees |

Version currency check (required before this review can be final evidence):

```text
git -C /home/mark/projects/agentic-coding-project-starter-kit rev-parse HEAD
c5229617220f97ff9228df261e24718b1d1d6d89
git -C ... status --porcelain      -> (no output; clean worktree)
git -C ... rev-list --count c5229617220f97ff9228df261e24718b1d1d6d89..HEAD
0
```

The kit worktree is **exactly at the pinned commit**, clean, with zero commits
since. The pin is therefore current and usable as final VAL-080 evidence.
Final QA must repeat this three-command check; if it returns a different HEAD
or a nonzero count, VAL-080 evidence is stale and must be re-derived.

**No held-out scenario, evaluator, or excluded-scenario material was opened.**
Only the four normal specification files above were read.

## 2. Preserved Qualification Distinction

`FID_COMPATIBILITY.md` states `Status: Target contract aligned; no production
candidate qualified`, and every one of the twelve rows records
`production BLOCKED`. Phase 22 flipping `attention.v1`, `kit.activation.v1`,
and `export.portable.v1` to `available:true` changes **fid's own binary
capability advertisement**, not kit qualification. Kit qualification requires
all twelve rows plus the seven-item `Qualification contract`, and both
`contracts.md` section 11 and `phase_brief.md` assign that to fid Phase 25.

This review requests **no change** to that distinction and records it so no
later artifact can read three `true` rows as adoption evidence.

Observed baseline for the record (`fid --json capabilities`, planning base):

| Capability | Now | Trust source now | Phase 22 target |
|---|---|---|---|
| attention.v1 | false | phase-22-unavailable | true, `{state:true,human_outcome_resolution:false}` |
| kit.activation.v1 | false | phase-22-unavailable | true, `recover:true` |
| export.portable.v1 | false | phase-22-unavailable | true, adds `import_mode:"empty-target"` |
| work.claims.v1 | false | phase-22-unavailable | false, `phase-24-unavailable`, `recover:true` |

The other eight rows (`evidence.freshness.v1`, `evidence.gate-definitions.v1`,
`evidence.receipts.v1`, `human.outcomes.v1`, `reviews.coverage.v1`,
`telemetry.trajectory.v1`, `workflow.state.v1`, `workflow.transitions.v1`) are
unchanged, matching `public-interface.md` section 5. Re-owning Claims to
`phase-24-unavailable` requires the `system_capabilities.trust_source` CHECK
edit in `schema.md` section 1, which is present and correct.

## 3. Findings On The Named Comparison Axes

Severity: **AGREE** (no action), **RECORD** (accepted difference, must be
written down before execution), **DECIDE** (semantic disagreement needing an
explicit Tech Lead disposition).

### C-01 Closed terminated/live/unknown liveness - AGREE

Kit `02_INTERFACES.md` 3.2: "Executor liveness is the closed `terminated`,
`live`, or `unknown` set: ordinary takeover requires confirmed `terminated`,
`unknown` always stops, and live takeover requires an authorized force-recovery
fact bound to the exact referenced current grant." Decision 0005 adds "Lease
expiry alone is not proof that filesystem mutation stopped."

fid `contracts.md` 2/4 and `public-interface.md` `Liveness`
(`"live"|"terminated"|"unknown"`) match exactly, including "Unknown liveness
always returns `executor-liveness-unknown`, including a force request" and
"Effective lease validity nevertheless uses `now < expires_at`". No difference.

### C-02 Exact Human force grant binding - AGREE, one RECORD

Kit requires the force fact "bound to the exact referenced current grant";
Decision 0005 restricts Human authority to force takeover and
ambiguous/protected repository state.

fid is strictly more specific: activated `human` Actor, explicit
`run.force-recover` capability, Envelope grant, reason, reference, exact
current `grant_ref{envelope_id,revision,ordinal}`, and expected Run/Claim
versions. It also records the live observation permanently and fences later
predecessor writes while explicitly refusing to claim the old process was
killed. This satisfies the kit contract.

**RECORD:** `grant_ref.revision` is pinned to the literal `1`. That is the
one-revision slice boundary; kit Decision 0007 multi-Issue revision is fid
Phase 24. Acceptable, but the pin should be named as a scope-limited
difference rather than read as the general contract.

### C-03 Budget measurement evidence - AGREE, one RECORD

Kit 3.2: "A caller cannot declare attempt-unit consumption ... Exact measured
time/token/spend consumption requires a bound evidence reference; unknown
measurements stay null."

fid `contracts.md` 6 matches: `consumed`/`remaining` rejected as caller input
and null on Envelope input; five fid-owned units cannot be supplied by a
caller; a numeric external sample requires a `Source` reference/digest;
fid-owned samples carry `measurement_ref:null` and are bound through source
sample identity; "Unknown metrics remain null, not zero."

**RECORD:** the kit's `derived` / `observed` charge model on `review.recorded`
is Review-scoped (kit Decision 0006) and has no fid counterpart in this slice;
fid has no `review_cycles` unit. Already deferred to Phase 23 by the validation
contract's Human/Review exclusion, but the *budget-unit* half of that deferral
is not written down anywhere.

### C-04 Exhaustion action and Attention semantics - RECORD

Kit 3.3: "V2-2 accepts `budget-exhausted`, `retry-exhausted`, and the
review-cycle-specific `review-rework-exhausted` as exhaustion actions; each
carries `circuit-breaker` attention semantics."

fid `contracts.md` 6 accepts only `exhaustion_action:"budget-exhausted"`, and
fid Attention has no `circuit-breaker` classification; deduplication is by
semantic key alone.

The atomic part agrees exactly: kit "persists the consumed value, boundary
event, blocked state, typed stop, and deduplicated attention atomically" versus
fid "Admission of the last allowed attempt is atomic with charging that Attempt
and raising its exhaustion Stop/Attention".

**RECORD:** `retry-exhausted` and `review-rework-exhausted` are Phase 23/24
work and `circuit-breaker` Attention classification is not adopted. Neither
exclusion is currently written in the validation contract's exclusion list.

### C-05 Immutable history - AGREE

Kit 3.6 "Import preserves identifiers and history ... Duplicate idempotency or
event IDs do not duplicate facts"; 3.7 "active and aborted terminal records are
immutable". fid: Stop is "an immutable fact plus append-only disposition"
(section 5); "Terminal history remains immutable in Events" and "The
predecessor's historical finalized Event is never edited" (section 8);
"Preserve project/entity/Event identity, order, hashes, history, and policy"
(section 10). No difference.

### C-06 Cross-resource agreement - AGREE on semantics, three RECORD/DECIDE

Kit 3.7 requires one activation record binding old/new versions, component
states, recovery export, point-of-no-return state, and authoritative version;
commit only after every required component is prepared and validated; abort
only before the point of no return; finalize only after every new component
agrees; recovery old-directed before the point of no return and new-directed
after it. fid `contracts.md` 8 matches all of that with five named resources
plus immutable internal components, commit as the durable point of no return,
`restore-prior` before and `complete-forward` after, and
`activation.recover-complete` verifying exact agreement before ending recovery.

**RECORD (a):** operation shape differs. Kit has one
`recoverActivation(activationId, observedComponents, idempotencyKey)`; fid
splits it into `activation.resource`, `activation.recover`, and
`activation.recover-complete`. Semantically equivalent; the V2-3 adapter owns
the mapping. This is the "spelling may differ" case, not a disagreement.

**DECIDE (b):** the kit's activation record binds a **recovery export**. fid's
`ActivationRecovery` is `{direction, actor_attribution, reason,
selected_event_id, completed_event_id}` and binds no export identity, and
`activation_recoveries` has no such column. Either the obligation is satisfied
elsewhere (fid's separate `export.portable.v1`) and should be recorded as such,
or it is a gap. This needs an explicit Tech Lead disposition before execution.

**RECORD (c):** kit's `activation-conflict` covers version, component state and
recovery direction together; fid splits direction into
`activation-direction-conflict`. Kit section 4 says "Unknown codes in a
supported contract version are treated as `unsupported-capability`", so an
unmapped fid-only code degrades incorrectly. The V2-3 adapter owns the mapping;
record it so the mapping is not silently assumed.

### C-07 Validate-only import semantics - RECORD

Kit 3.6 exposes `exportPortable(scope, asOf?)`,
`validatePortableExport(export)`, and `importPortable(export, mode,
idempotencyKey)`, and states "`validate-only` is mandatory."

fid exposes `export create`, `export validate` (read-only, non-mutating), and
`export.import` with **no `mode` argument**; `import_mode:"empty-target"` is a
capability limit, not a request field. `contracts.md` 10 makes validation
mandatory before publication and non-mutating, and VAL-057 proves validation has
no project or destination effects.

The *semantic* requirement is satisfied: a mandatory, non-mutating validate-only
path exists and always runs before publication. The *shape* differs.

**RECORD:** the adapter maps kit `mode="validate-only"` onto `fid export
validate`, and `scope`/`asOf` are not exposed in this one-project slice. Write
both down; otherwise a later reviewer can read fid's missing `mode` as a
missing mandatory capability.

### C-08 Caller-declared terminal Run state - DECIDE

This is the one material semantic disagreement.

Kit 3.2 still names `finishRun(runId, terminalState, idempotencyKey)` - an
operation whose signature accepts a **caller-declared terminal state**.

fid Phase 22 deliberately deletes exactly that: `contracts.md` 7 "Replace
caller-controlled `run.finish` outcomes with `run.close` and `run.abandon`",
VAL-033 (success derived only from qualifying trusted facts) and VAL-039
(`run.finish` absent after derived closure lands). `fid run finish` exists on
the planning base today and is removed in task 9.

fid's direction matches the kit's stated *intent* - kit 3.4 "A caller cannot
submit that terminal Receipt", and `FID_COMPATIBILITY.md` qualification item 5
forbids "reported-evidence downgrade" and "caller-supplied terminal Receipt" -
but contradicts the kit's *named signature*.

**DECIDE:** the Tech Lead must record the disposition explicitly, e.g. "fid
derives the terminal state; the V2-3 adapter maps `finishRun` onto
`run.close`/`run.abandon` and treats `terminalState` as an assertion checked
against derived truth, never as an input." Without a recorded disposition,
VAL-080 cannot pass, because an unresolved contract difference is its exact
fail criterion.

### C-09 Attention deduplication identity - RECORD (kit is internally inconsistent)

Kit 3.5 signature is `raiseAttention(attentionItem, deduplicationKey,
idempotencyKey)` - a **caller-supplied** dedup key. Kit 3.3 says "Attention
uses its own stable semantic deduplication identity, so repeated equivalent
interruption requests collapse without discarding stop facts" - a
**backend-derived** identity.

fid follows 3.3: `semantic_key` is SHA-256 of canonical
`{subject,code,resolution_condition}`, `stop.raise` accepts no dedup key, and
`attention_one_open_semantic_key` enforces one open item per key.

fid's choice is the safer of the two and matches the kit's own invariant
statement. **RECORD** that fid derives the key and that the adapter must derive
or discard `deduplicationKey`; flag the kit-internal inconsistency back to the
kit owner rather than changing fid.

### C-10 Attention resolution shape - RECORD

Kit 3.5 `resolveAttention(attentionId, outcomeRef, idempotencyKey)`; fid
`attention.resolve{attention_id, expected_version, fact_event_ids:ID[]}`. fid
resolves from durable qualifying facts and defers Human Outcomes to Phase 23,
which it declares honestly through the `attention.v1` limit
`human_outcome_resolution:false` (consistent with DEC-172 negative metadata).
Record the shape difference and the declared limit as the disclosure of it.

### C-11 Who observes the repository - RECORD, fid is stricter

Kit Decision 0005: "the kit/harness captures and reconciles the current
worktree"; kit 3.2 requires "a repository-port reconciliation result".

fid `contracts.md` 3 puts observation inside fid (`internal/project` adapter,
structured `git ... -z`) and, at recovery, "the adapter observes again under
the fid write lock and compares the supplied observation digest with the actual
state." The two are compatible: the harness may supply a digest, fid verifies
it and rejects any intervening change with `recovery-state-conflict`.

**RECORD:** fid never trusts a supplied reconciliation as fact. This is a
strengthening, not a divergence, and should be stated so a kit adapter author
does not expect fid to accept harness-computed reconciliation.

Both agree on non-mutation. Decision 0005: "Recovery never automatically
resets, cleans, stashes, checks out, commits, or discards Git state." fid
`contracts.md` 3 uses almost the same sentence and VAL-016 proves it on both
success and failure paths.

### C-12 Resumption telemetry - AGREE by derivation, RECORD

Kit Decision 0005 consequence: the recovery transition "records the interrupted
Run, superseded Claim, new Claim, linked Run, repository observation, and
resumption/intervention telemetry without resetting budgets."

fid charges `interventions` for force recoveries only; there is no
`resumptions` unit in the closed unit list, and `schema.md` 2 states "No budget
total, Next Action, or second recovery-status aggregate is stored." Resumption
count is instead **derivable** from the `run_recoveries` chain, surfaced as
`Resume.recoveries:Recovery[]`. Budget non-reset is explicit and proven by
VAL-029.

**RECORD:** resumption telemetry is derived, not stored, and
`telemetry.trajectory.v1` stays Phase 25. Without that note, a reader comparing
Decision 0005 line by line will read a missing counter as a missing obligation.

### C-13 Stop-code mapping - RECORD

| fid Phase 22 StopCode | Kit stable stop code | Note |
|---|---|---|
| executor-liveness-unknown | executor-liveness-unknown | exact |
| recovery-state-conflict | recovery-state-conflict | exact |
| budget-exhausted | budget-exhausted | exact |
| activation-incomplete | activation-incomplete | exact |
| capability-unavailable | unsupported-capability | spelling difference only |
| budget-data-unavailable | **none** | kit default maps unknown codes to `unsupported-capability`, which is semantically wrong for a data-absence stop |
| operator-stop | **none** | no kit equivalent |

**RECORD** this table. Two fid stop codes have no kit port code and the kit's
declared fallback mis-classifies at least one of them; the V2-3 adapter owns
the mapping, and it must be written down rather than inferred.

## 4. Verdict

| Axis | Result |
|---|---|
| Closed liveness | AGREE |
| Force grant binding | AGREE (1 RECORD) |
| Budget measurement evidence | AGREE (1 RECORD) |
| Exhaustion / Attention | RECORD |
| Immutable history | AGREE |
| Resource agreement | AGREE (2 RECORD, 1 DECIDE) |
| Validate-only import | RECORD |
| Caller-declared terminal state | **DECIDE** |
| Attention dedup identity | RECORD |
| Attention resolution shape | RECORD |
| Repository observation ownership | RECORD |
| Resumption telemetry | AGREE by derivation (RECORD) |
| Stop-code mapping | RECORD |
| Qualification distinction | PRESERVED, no change requested |
| Kit version currency | CURRENT at the pin, 0 commits since |

No unresolved *contradiction* was found between fid's proposed recovery
semantics and the pinned kit port. Two items require an explicit Tech Lead
disposition before execution (C-06b activation recovery-export binding, C-08
caller-declared terminal state) and nine accepted differences must be written
down. VAL-080's fail criterion is "unresolved contract difference", so the
Tech Lead's dispositions for C-06b and C-08, plus the recorded RECORD items,
are prerequisites for a VAL-080 `PASS` at final review.

## 5. Re-check Obligation At Final Review

VAL-080 also requires that final changes are rechecked. Final QA must:

1. Re-run the three version-currency commands in section 1.
2. Re-read kit `02_INTERFACES.md` 3.2-3.7 and `FID_COMPATIBILITY.md` if the
   pin moved, and re-derive every row of section 3.
3. Compare the **as-built** wire surface (task 18 corpus manifest) against this
   review, not only the planning documents.
4. Confirm the C-06b and C-08 dispositions exist in `phase-22/decisions.md`
   and that no RECORD item was silently dropped.
