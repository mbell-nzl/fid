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
12. Keep VAL-036 as the sole Human parent assertion and, per Workflow §6.1,
    fan it out into HUMAN-01a..01d over the deterministic terminal fixture:
    stream identity, typed controls, terminal result, and no raw-control display
    effect. All four cite VAL-036; do not delegate or pre-answer them.

## Readiness Response

Return `Ready` only when all commands are concrete/bounded, every VAL has one
mapping, all seven review windows and both checkpoint assertions are present,
fresh feasibility is executable, final QA is complete, and VAL-036 is the sole
Human item. Otherwise return `Needs Refinement` with exact contract/task
changes requested.

Task 1 remains unauthorized after the QA Lead response. The Tech Lead must
apply accepted refinements, then the Human must approve execution.

## Tech Lead Refinement Response - 2026-08-23

DEC-158 applies BC-01/BC-02 and R-01..R-18, except R-09's final-status premise:
final role-policy and `--by` searches must terminate 1, while broad retained
reader/Event searches terminate 0 and all matches classify to
`UNCLASSIFIED=0`. R-19 remains deferred until Human execution approval. QA Lead
must re-review the amended contract/tasks and update its own plan before a
`Ready` verdict. In particular it must:

1. replace BC-01/QA-80/final-performance language with DEC-158's exact reference
   identity and three consecutive fresh-process, all-samples-`<30 s` rule;
2. close BC-02 and use only the Phase Runner's two-argument post-commit report
   lint before the next dispatch, including in VAL-042/VAL-087 cadence text;
3. bind the task 14 Human fixtures, task 22 committed corpus/non-vacuity, exact
   20-Event/one-`envelope.created` positive control, containment lexicon, and
   Group 1 environment-only feasibility wording; and
4. correct R-09/QA-41/QA-93 status expectations without narrowing the two broad
   classified inventories.

Mark `qa-plan.md` `Ready` only if the independent re-review finds no remaining
Must-fix. Task 1 stays unauthorized.
