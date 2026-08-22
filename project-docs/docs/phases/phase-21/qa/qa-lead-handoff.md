# Phase 21 QA Lead Handoff

> Requested artifact: **Ready Heavy QA plan**
> Current state: planning only; task 1 is blocked
> Tech Lead authority: DEC-152..157

The Tech Lead requests a QA Lead plan; this file is not the QA plan and does
not assign QA verdicts. Read only the following planning authorities plus the
QA-role/workflow sections required to author the plan:

- `../contracts.md`
- `../public-interface.md`
- `../v1-convergence-inventory.md`
- `validation-contract.md` (VAL-001..087)
- `group-1-checkpoint-contract.md`
- `../tasks/task-1-*.md` through `../tasks/task-22-*.md`
- `../../../../specs/08_QUALITY_PLAN.md` applicable Heavy/phase sections
- `../../../../specs/09_AI_WORKFLOW.md` QA Lead, evidence, review, and gate sections

## Required Ready Plan Content

1. Map every VAL-001..087 exactly once to a QA step, task review, checkpoint,
   final-QA step, or Human checklist. Do not merge IDs. Preserve the validator
   role in each contract row.
2. Provide concrete bounded commands, working directories, timeouts where
   needed, pass tokens, and expected termination for every deterministic step.
   Declare every expected non-zero command and its exact accepted exit code.
3. Insert any QA-plan command a Developer must run into the owning task gate
   proposal. Identify the task and exact command so the Tech Lead can amend the
   spec before authorization.
4. Schedule QA Review windows after tasks 1-3, 4-6, 7-9, 10-12, 13-15,
   16-18, and 19-21. Findings return to the owning task; task 22 follows the
   seventh window. Apply task-local completion-report lint before every next
   dispatch.
5. Define the independent Group 1 review after task 12 and before task 13.
   It must separately decide VAL-042 and VAL-043, classify the exact temporary
   V1 inventory, and prove no new generic path consults it. Do not require
   final V1 absence before task 20.
6. Run fresh Group 2 feasibility probes in the checkpoint session for:
   structured argv; separate interleaved streams; invalid UTF-8/control bytes;
   timeout/signal/spawn failure; high-output draining; closed environment;
   credential/literal canaries across every surface; cwd/Artifact containment;
   changed-during-copy; atomic content storage; observed termination;
   tree freshness; Artifact loss/restoration; deterministic ZIP; non-mutating
   corruption validation; expected-nonzero absence commands; and the bounded
   tagged real-process journey. Transcribed Phase 20 evidence is not fresh.
7. At task 20, use the classified inventory rather than a raw no-match rule.
   Frozen/diagnostic readers are permitted only where
   `v1-convergence-inventory.md` says so. Prove migration 0014 and final
   commands/writers/Event acceptance agree with the inventory.
8. At task 21, run final-binary fresh init after 0014, delete the projection,
   rebuild without `roles`, compare named canonical reads and export bytes,
   exercise the three maintenance finding classes, and capture the complete
   100k Event-v2 benchmark termination plus hardware/toolchain identity.
9. At task 22, generate or enumerate conformance cases for every operation row
   in `public-interface.md`, including args, closed input, result, Event,
   idempotent replay, Human rendering, every failure details/channel/exit, and
   exact live-frame variants. Include an inventory that detects snake_case in
   closed public values while excluding JSON field names and opaque IDs.
10. The real-process journey must mutate once in this order: nonterminal
    transition; Gate execution; Receipt link/evaluation; terminal-success
    transition; Run finish/Claim release; Envelope completion; resume;
    export/validate. After delete/rebuild, repeat only named reads, evaluation,
    completion state, resume, export, and negative checks. Assert no second
    Envelope or other post-rebuild mutation.
11. Define final QA with the full phase command gate from `AGENTS.md`, PRD
    compliance, security review of executor/redaction/path handling,
    performance spot-check, documentation review, every VAL outcome, task/report
    audit, and final no-unclassified-inventory-match proof.
12. Create exactly one Human checklist item for VAL-036 using the deterministic
    recorded terminal fixture. Ask the Human to confirm stream identity, typed
    controls, and terminal result are understandable and no raw control affects
    display. Do not delegate or pre-answer it.

## Readiness Response

Return `Ready` only when all commands are concrete/bounded, every VAL has one
mapping, all seven review windows and both checkpoint assertions are present,
fresh feasibility is executable, final QA is complete, and VAL-036 is the sole
Human item. Otherwise return `Needs Refinement` with exact contract/task
changes requested.

Task 1 remains unauthorized after the QA Lead response. The Tech Lead must
apply accepted refinements, then the Human must approve execution.
