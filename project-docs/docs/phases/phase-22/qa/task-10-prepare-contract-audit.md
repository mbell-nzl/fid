# Task 10 Activation Prepare Contract Audit

Date: 2026-09-07. Author: delegated Tech Lead reviewer.
Inspected HEAD: `7467e1058ebb51c6d7793905b74fe2cef5c69415`.
Scope: whether accepted prepare/staging records resolve the timing of immutable
authority-profile and Gate Definition content. Read-only contract/source audit;
no quality gate, implementation edit, test edit, or VAL outcome.

## Finding

**A genuine timing contradiction requires an explicit recorded resolution
before task 10 relies on either interpretation.** The accepted closed wire and
Event contracts retain component staging after prepare. Two prose clauses
instead require the new internal component contents to be immutable prepare
inputs. No inspected DEC, glossary definition or QA disposition resolves that
contradiction. It is not an unassigned function or an ordinary missing test.

The minimum correction can preserve the existing wire, DDL, Event families,
recovery direction and trust rules: explicitly distinguish prepare's five
resource bindings from later immutable internal-component staging and exact
selection at commit. That correction still needs written authority; a Developer
must not silently reinterpret the contradictory clauses or invent fields.

## Exact Accepted Sources

| Source at inspected HEAD | What it says | Consequence |
|---|---|---|
| Phase 22 `contracts.md:307-311`, section 8 | Prepare freezes the old/new complete sets, enumerating kit, overlay, provider, project, backend, authority profile and Gate Definitions. | Literal single-operation reading binds both new internal contents at prepare. |
| Phase 22 `tasks/task-10.md:28` | “Five resource identities plus authority/Gate content are immutable prepare inputs”. | Explicitly repeats that temporal requirement; cannot be treated as a casual example. |
| Phase 22 `public-interface.md:160-161` | Prepare is retained Phase 21 input plus predecessor ID/version and exactly five resource identities; component retains its closed Phase 21 input plus expected activation version. | Prepare has no new internal content or internal digest commitment field. Component still receives content separately. |
| Phase 21 `public-interface.md:167-169` | Prepare carries kit/backend identity, phase and direction. Component carries type, identity/version, expected prior digest and content. Commit carries the exact two sorted unique component digests. | New authority/Gate content first enters through component staging; commit selects the supplied pair. |
| Phase 22 `schema.md:306` and section 5 introduction | `activation.prepared` extends existing prepare with predecessor identity and five pending `activation_resources`; only listed domain Effects are allowed. | No prepare-time profile/Gate insertion or hidden content commitment is authorized. |
| Phase 21 `contracts.md:432-436` | Prepare inserts `kit_activations`; separate component-staged Events insert the component plus profiles/definitions; commit updates activation/components. | The retained Event grammar expressly separates those boundaries. |
| Phase 21 `contracts.md:133-145`, section 4 | Stage validates/content-addresses exact components; commit atomically selects their identities; finalize activates that committed identity. | Existing recorded behavior supplies a coherent staged protocol, not prepare-time content. |
| Phase 22 `schema.md:311-320` | Existing commit/finalize additionally validate resource agreement; every later activation Event increments its version. | Versioning fences later calls; it does not retroactively supply absent prepare data. |
| Canonical `07_DESIGN_PLAN_FLOWS.md:18-24` | Prepare names old/new identities, then each participant stages/validates under that activation ID. | Supports staged preparation; does not erase the more explicit contradictory Phase 22 prose. |

`03_TRD_SCHEMA.md` section 8.3 incorporates the Phase 22 wire/semantic/Event
documents together. Phase 22 `contracts.md` section 1 retains Phase 21 except
explicit replacements; it does not define a precedence rule that permits
ignoring one of these jointly accepted requirements.

## Why The Contradiction Cannot Be Implemented Away

Consider a first activation with a valid prepare request P, followed by a valid
authority-profile payload A or another valid payload B. Both histories have
identical P: neither new profile content nor its digest occurs in that request
or its allowed prepare Effects. Before a component call, fid has no accepted
input distinguishing A from B. Therefore P cannot already freeze the exact new
profile content. The same argument applies to Gate Definitions.

Opaque `new_identity` text on one of the five external resources is not a
declared internal-content digest or an authorized resolver. Reading project or
provider files to obtain that missing content would invent a new authority
path. Copying the prior components would fail first activation and prevent
an intended profile/Gate update. Treating future Event contents as if they were
present at prepare would falsify history. None is an accepted solution.

Current `internal/workflow/activation.go` agrees with the retained staging
shape: `ActivationPrepareRequest` has no component content; the component
request carries `Content`; `StageActivationComponent` constructs inserts;
`CommitActivation` checks the supplied two digests. This is corroborating
source inspection, not authority to prefer existing code over accepted specs.

## Decision And QA Record Audit

- DEC-180 proposes and DEC-188 jointly accepts the slice/wire/DDL/task pack.
  Neither disambiguates this prepare boundary. `git log -S` identifies both
  conflicting prose clauses in planning commit
  `6563fb4b4b92402c9df7826c577041bcac4e65b3`
  (`Refs: phase-22/task-planning`); they were not introduced by a Developer.
- DEC-184/R-10 resolves activation version/error shapes and bootstrap authority,
  not prepare-time contents. DEC-187/C-06b explicitly excludes coordinator-owned
  recovery-export binding and says that disposition needs no activation-input
  extension. It does not secretly supply new internal component contents.
- `qa/consistency-review.md` C-06 describes five resources plus immutable
  internal components and commit after preparation/validation. Its AGREE
  assessment does not name or resolve the timing discrepancy. C-06a concerns
  recovery-operation spelling, not prepare/component staging.
- VAL-040 names exact old/new resources and requires “Five resources and
  immutable internal components bind one predecessor.” QA plan EVID-35/QA-41
  asks first-activation/upgrade prepare cases for exactly five resources and
  one predecessor. Those criteria do not provide the missing prepare input.
- VAL-041..044 retain the before/after commit directions, exact agreement,
  normal-work block, and immutable recovery identity/direction. All can be
  evidenced with the retained staged protocol; none chooses between the
  contradictory prepare readings. Do not mark VAL-040 satisfied by silently
  dropping internal-component immutability.

## Proposed Bounded Written Clarification

For the root's recorded resolution, replace the conflicting prepare paragraph
in `contracts.md` section 8 with wording equivalent to:

> Prepare records one exact predecessor activation ID/version, or both null
> for first activation, and freezes the old/new identities of the five declared
> resources: kit, overlay, provider, project and backend. The predecessor's
> immutable internal authority-profile and Gate Definition components remain
> identified through that predecessor. New internal component content is
> supplied by the retained `activation.component` operations under the same
> new activation ID; staging validates and content-addresses that content,
> which remains immutable. Commit selects the exact two staged component
> digests after the required resource agreement checks and is the durable
> point of no return. Finalization/recovery cannot mix the predecessor's
> components with the new activation's selected set.

Replace task 10's conflicting binding constraint with:

> The five old/new resource identities and predecessor binding are immutable
> prepare inputs. Authority-profile/Gate content is supplied and made immutable
> through the retained component-staged operations under that same activation
> ID; commit binds the exact two staged digests. Every mutation checks the
> accepted activation/resource versions. No new prepare field, DDL column,
> Event family or provider-file read is authorized.

Add a matching task-10 acceptance criterion that fails if the distinction is
lost: prepare stores exactly the declared resource/predecessor bindings;
separate component calls record immutable content under the same activation;
commit/finalize reject incomplete, mismatched or mixed sets; pre-commit
restore and post-commit forward recovery preserve the corresponding selected
trust. Cite VAL-040..044 and retain the existing exact gates. A short note in
`schema.md` section 5 and `public-interface.md` section 3 can cross-reference
that resolution while leaving their normative rows byte-identical.

This recommendation preserves the explicit closed operation protocol and
strengthens its timing explanation. It is a proposed resolution of contradictory
accepted prose, not a claim that the original pack was already unambiguous.
If the authority instead intends exact new internal contents to be committed
in the single prepare operation, that requires a separately accepted wire and
Event amendment before implementation; it cannot be achieved by wording alone.

## Handback

The root owns the disposition and any spec/task edits. Task 10 should receive
the recorded clarification before dispatching dependent implementation. This
audit does not allocate an ESC, consume recovery budget, change a gate, or
claim any implementation/QA success. Only this new audit artifact was edited.

## Subsequent Authority Disposition

After the independent audit commit `8a4ff6d`, the root explicitly accepted its
closed-interface-preserving recommendation under the Human's super-delegate
authority and assigned this reviewer the bounded documentation correction.
DEC-198 now records that resolution; `contracts.md` section 8 and task 10's
binding/acceptance language distinguish prepare bindings from immutable staged
content and exact commit selection. The original finding above remains the
record of the inspected baseline. Its specification contradiction is resolved
for dispatch; runtime conformance remains unverified by this audit.
