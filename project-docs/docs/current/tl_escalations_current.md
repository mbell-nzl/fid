# TL Escalations — Current

## Open escalations

**ESC-23-001 - widening the journal `Condition` value breaks journal replay and
Attention deduplication.** Raised by the Tech Lead 2026-09-19 while
dispositioning R-19 under the widened inventory frame. `public-interface.md`
section 2 adds three nullable fields to `Condition`; it is the only field-adding
widening of a live journal value shape in the phase. `validateCurrentShape`
matches field count exactly, and `ReplayJournal` revalidates every historical
Event against current shapes, so no workspace holding a pre-widening Stop can
rebuild; and `attentionSemanticKey` hashes the whole struct, so every open
Attention's stored `semantic_key` stops matching and an identical Stop raises a
duplicate. Both measured under overlay, working tree untouched. Every remedy
changes an accepted value's encoding or adds an `EventSchema` bump plus a
semantic-key migration, so the choice is the Human's. **Blocks tasks 9, 10, and
11 only**; tasks 1-8 may dispatch. Recommended Option A: non-applicable fields
absent rather than null, with per-type Condition shapes keyed on `type` exactly
as `currentStopDetailShapes` is keyed on `code`. Full record, evidence, and
options: [Phase 23 escalations](../phases/phase-23/tl-escalations.md).

ESC-22-005 is resolved by the Human's performance-repair approval under
DEC-221; its proposed timeout increase was not approved or applied. Task 20 is
complete; hosted run34805256766 passed on6b8fe50. See
[Phase 22 escalations](../phases/phase-22/tl-escalations.md) for all five records.
The wrap's per-TRG hindsight audit is in
[wrap-up-audit.md](../phases/phase-22/wrap-up-audit.md).

## Standing handling

Register an escalation when allocated, not after resolution. Ordinary causal
repair and incomplete in-scope work do not themselves require a Human stop.
Read accepted ownership before calling it ambiguous (DEC-194). A real contract/
schema contradiction needs the missing authority, not a weakened assertion
(DEC-195). External interruptions do not consume repair no-progress budget.

Prior archives: [Phase 21](../phases/phase-21/tl-escalations.md),
[Phase 20](../phases/phase-20/tl-escalations.md),
[Phase 19](../phases/phase-19/tl-escalations.md).
