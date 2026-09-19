# TL Escalations — Current

## Open escalations

None.

## Recently resolved

ESC-23-002 is **RESOLVED 2026-09-20 under DEC-238**. The Human delegated
bounded repair authority above the fixed agent roles. Rebuild statement reuse
and fresh/isolated test setup repair the causes without changing any gate or
assertion. Fresh exact gates pass: CLI **527.099s**, workflow **231.010s**,
tagged integration **669.586s**; rebuild samples **22.728s / 22.623s / 22.489s**.
Integration binary directories stay **0 → 0**. Task 1 is complete in its report;
the Runner must reconcile the retained STOP before advancing. Tasks 9-11 stay
held. See [repair evidence](../phases/phase-23/qa/esc-23-002-repair.md) and
[resume handoff](../phases/phase-23/qa/resume-after-esc-23-002.md).

ESC-23-001 is **resolved**: the Human chose Option A on 2026-09-20 - `omitempty`
on the three new `Condition` fields only, the six Phase 22 types keeping
`budget_name` and `capability_id` present-and-null, with per-type journal shapes
keyed on `type`. Recorded as DEC-235, which amends `public-interface.md` section
2, restates inventory row 8, creates row 10, and rewrites tasks 9, 10, and 11.
The three held revisions (R-22, R-23, I-37) shipped with it. **ESC-23-001
blocks no task**; tasks 9-11 are held only by the Human's dispatch instruction
and one bounded QA confirmation of the DEC-237 contract edits. See
[Phase 23 escalations](../phases/phase-23/tl-escalations.md).

QA pass 8 raised one **blocker candidate with no identifier allocated** (R-28,
VAL-078's granularity), which is the correct handling: the repair is a Tech Lead
contract edit, not a Human decision. Dispositioned under DEC-237 by splitting
the row. No escalation was or should have been opened.

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
