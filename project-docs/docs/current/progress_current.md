# Progress - Current

> In-flight status only. Historical progress remains in phase archives.

## Project State

**Phase 21 is active in Group 2 with the ledger STOPPED after TRG-21-020 /
ESC-21-014 at task 20; DEC-173 resolves the contradiction between task 20's
explicit file list and the reach of its own absence controls and gates.** Tasks
1-19, reviews 1/1b through 6, and the Group 1 checkpoint are mechanically
complete. Phase 20 completed/wrapped on 2026-08-22; its retrospective is in
`../reviews/2026-08-22-phase-20-retrospective.md`.

Phase 21 planning decisions: DEC-153 (Gate D3.2), DEC-154 (slice/spec
contract), DEC-155 (carry-forward), DEC-156 (terminal gate), DEC-157 (complete
interface, bootstrap, inventory, sequencing), DEC-158 (QA refinements, rebuild
target), DEC-159 (QA readiness and Human approval).

Execution decisions: DEC-160 bounds the temporary V1 policy bridge; DEC-161
selects the valid ledger recovery record; DEC-162 converges retained V1 tests
inside the invalidating task; DEC-163 gives QA `REPAIR` an ordered owner list;
DEC-164 keeps schema-v2 repair authority byte-exact; DEC-165/166 make
convergence ownership follow the invalidating change, never a file subject or
historical label; DEC-167 requires a task's gates to exercise every scope it
owns; DEC-168 removes lifetime budgets from QA child sessions; DEC-169 resizes
the checkpoint repairs; DEC-170/171 make externally interrupted attempts
non-consuming and executable in the Runner controls; DEC-172 keeps the
unavailable `attention.v1` row and replaces the zero-match search with a
boundary test.

DEC-173 (2026-09-02) closes ESC-21-014: task 20's absence controls walk
`cmd/fid internal testdata` and its coverage gate compiles `./...`, but its file
list omitted `internal/bootstrap/`, `internal/journal/`, `internal/labels/`,
`internal/links/`, `internal/requirements/`, `internal/db/`,
`testdata/journal/`, and `cmd/fid/integration/`. The inventory already assigned
task 20 every residue, so the scope is amended with bounded subject-preserving
dispositions (delete `LockOwner.ActorRole`/`actor_role`, the bootstrap
fixed-role registry, and the `ActorRole`/`legacyRole` fallback readers; converge
the `dev.*` golden, every `roles` fixture, and the V1 closure journey) and task
20 gains the tagged integration gate under DEC-167. No command, expected status,
VAL, 0014, frozen-table, or public contract changes. `08_QUALITY_PLAN.md` §3
records the converse executable-ownership rule. TRG-21-020, ESC-21-014, and all
failed evidence remain historical.

## Planned Boundary

One Heavy Phase 21 uses 22 bounded Developer tasks in two internal groups. An
independent checkpoint after task 12 proved Group 1 contract closure and fresh
Group 2 executor/verifier feasibility before task 13. The validation contract
contains VAL-001..087. Seven three-task review windows, Phase Runner
two-argument post-commit report lint before each next dispatch, final QA, and
one Human terminal-output parent assertion with four Human subchecks apply.

## Approval State

- Human approval of revised Heavy mode, exact contracts, VAL contract, and task
  pack: **APPROVED 2026-08-23 / DEC-159**.
- QA Lead Heavy QA plan with fresh feasibility: **READY** at `e8a636e`; no
  outstanding Must-fix.
- Product implementation: **AUTHORIZED UNDER THE ORDERED TASK PACK**.
- Dispatch state: **TASKS 1-19 COMPLETE; TASK 20 STOPPED AFTER TRG-21-020 /
  ESC-21-014 WITH RECOVERY AUTHORIZED BY DEC-173**. The Phase Runner appends
  `RUN-RESUME` retaining those IDs, the exact recorded trigger, attempt 2, and
  resolution authority `DEC-173`, then one `TASK-RECOVERY task-20` with
  `in-scope-implementation`, `task-scope`, and no-progress `0/2`.
- Planning verification: **PASS** - deterministic coverage/inventory audit,
  `project-docs/tools/workflow-lint.sh`, and `git diff --check` all terminated
  `exit_code: 0` after the DEC-159 approval-state update.

## Carry-Forward

DEC-155 schedules CF-020/027 in task 8, CF-032 in task 12, and CF-024 in task
21. CF-008/018/023/026a/026b/028/029/030 are explicitly re-deferred; closed or
cancelled rows remain closed/cancelled.

## Next Work

1. Phase Runner appends the DEC-173 recovery sequence: `RUN-RESUME` bound to
   TRG-21-020 / ESC-21-014 with the recorded trigger, attempt 2, resolution
   authority `DEC-173`, and the exact next-recovery tuple from DEC-173 item 7;
   it then appends the matching `TASK-RECOVERY task-20` with no-progress `0/2`.
2. The Developer continues task 20 from `aef03cf` under the amended spec:
   deletes the enumerated residues, converges the named fixtures and the tagged
   closure journey with per-row provenance, commits durable work before long
   commands, and runs the complete byte-identical amended command block
   including `go test -tags integration -race ./cmd/fid/integration/...`.
3. After the report commit and two-argument report-lint PASS, the Runner records
   `REPAIR task-20 RESOLVED` and `TASK task-20 COMPLETE`, then dispatches task
   21. Window 7 review remains gated on completion of tasks 19-21.
