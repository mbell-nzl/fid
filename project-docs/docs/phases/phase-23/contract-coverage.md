# Phase 23 Contract Coverage Matrix

> Tech Lead planning cross-check, 2026-09-14. ACCEPTED under DEC-222.
> Amended 2026-09-15 for the QA Lead's independent re-derivation: every clause
> the re-derivation found mapped-but-unasserted now names a row that evidences
> it, and the duplicated `links` row was removed. The manual cross-check worked
> as intended - it is what surfaced R-02 and R-12 before dispatch - and that
> outcome is recorded for CF-039 at wrap-up.
>
> **Why this file exists.** Phase 22 shipped a contract clause that had no
> schema row (DEC-195) and an accepted allocation that was never delivered
> (DEC-194). Both were planning gaps discovered during implementation. CF-039
> records the executable version of this cross-check as workflow maintenance;
> this file is the manual Phase 23 mitigation.
>
> **Rule.** Every binding clause in `contracts.md` has a row here naming its
> schema row (or an explicit `domain-only`), exactly one owning task, and at
> least one `VAL-NNN`. A clause with no owning task, no enforcement point, or
> no assertion is a planning defect that blocks execution. QA Lead re-derives
> this matrix independently from `contracts.md` before execution approval.

## Matrix

| Contract clause | Enforcement point | Owning task | VAL |
|---|---|---|---|
| 2: three Obligation types | `0013` `envelope_obligations.obligation_type` CHECK | 1 defines, 2 uses | VAL-002, VAL-006 |
| 2: registration null exactly for the new types | `0013` registration discriminator CHECK | 1 defines, 2 uses | VAL-002, VAL-006 |
| 2: `review_profile` / `outcome_class` present and unique | `0017` ALTER plus two partial unique indexes | 1 defines, 2 uses | VAL-006, VAL-007 |
| 2: at most one `gate-evidence` Obligation | `envelope_obligations_one_gate_evidence` | 1 defines, 2 uses | VAL-007 |
| 2: 1..16 Obligations, unique IDs, primary-Issue subject | domain-only (`internal/workflow/envelope.go`) | 2 | VAL-006, VAL-007 |
| 2: `guarded_actions` subset and per-type limits | domain-only | 2 | VAL-009 |
| 2: one Evaluation per Obligation in ID order | domain-only (`internal/evidence`) | 2 | VAL-008 |
| 2: three independent evidence dimensions and all reasons | domain-only | 2 | VAL-008, VAL-010 |
| 3: `review` Claim mode | `0013` `work_claims.mode` CHECK | 1 defines, 3 uses | VAL-011 |
| 3: one active Claim per Issue across modes | retained `claims_one_active_mutation` index | 3 | VAL-011 |
| 3: review Claim cannot adopt dirty paths | domain-only plus `claim-mode-invalid` | 3 | VAL-012 |
| 3: review Run lifecycle otherwise unchanged | domain-only | 3 | VAL-013 |
| 3: independence computed from Claim and Run history, mutating side only (DEC-225) | `review_records.independence` CHECK plus domain | 4 | VAL-015 |
| 13: link-time independence re-evaluation | domain-only plus `review-not-independent` | 6 | VAL-077 |
| 4: immutable Review binding | `review_records` columns and FKs | 4 | VAL-014 |
| 4: repository re-observation at record time | domain-only plus `repository-mismatch` | 4 | VAL-016 |
| 4: closed coverage structure | `review_records.coverage_json` plus domain | 4 | VAL-014 |
| 4: three verdicts and their Finding consistency | `review_records.verdict` CHECK plus domain | 4 | VAL-017 |
| 4: one Review per review Run | `review_records.run_id UNIQUE` | 4 | VAL-018 |
| 4: referenced Receipts belong to the recording Run | domain-only | 4 | VAL-014 |
| 4: a `pass` verdict is not an Issue mutation | domain-only | 4 | VAL-024 |
| 5: Findings created only with their Review, atomically | `review.recorded` effect rule stage 1 plus stage 0 | 5 | VAL-020 |
| 5: severity and blocking must agree | `review_findings` composite CHECK | 4 defines, 5 uses | VAL-017 |
| 5: `addressed` requires an active mutating Claim | `finding_dispositions` action CHECK plus domain | 5 | VAL-021 |
| 5: only a later independent Review resolves | `finding_dispositions.qualifying_review_id` CHECK plus domain | 5 | VAL-022 |
| 5: at most one terminal disposition | `finding_dispositions_one_terminal` | 5 | VAL-022 |
| 5: `Finding.state` derived, never supplied | domain-only | 5 | VAL-023 |
| 5: unresolved blocking Finding keeps the Obligation unsatisfied | domain-only evaluation | 6 | VAL-026 |
| 6: `review-cycles` unit | domain-only budget validation | 6 | VAL-030 |
| 6: `review-rework-exhausted` exhaustion action | domain-only plus `stops.code` CHECK | 1 defines code, 6 uses | VAL-029, VAL-030 |
| 6: exactly one charge per `changes-required` Review, keyed on Review ID | `budget_samples` UNIQUE (envelope, revision, name, source, source_sample_id) | 6 | VAL-027 |
| 6: `run.observe-budget` rejects `review-cycles` | domain-only | 6 | VAL-028 |
| 6: atomic admission, charge, Stop, and Attention | `review.recorded` effect rule optional rows | 6 | VAL-029 |
| 6: blocked actions exactly `claim.acquire`, `run.start` | `stops.blocked_actions_json` plus domain | 6 | VAL-029 |
| 6: unit and exhaustion action require each other; budget requires a `review` Obligation | domain-only | 6 | VAL-030 |
| 7: target change stales review Evidence, Record retained | domain-only freshness evaluation | 7 | VAL-031 |
| 7: `human-outcome` Obligations are never tree-staled | domain-only | 7 | VAL-032 |
| 7: later Outcome supersedes for evaluation, both retained | domain-only plus `human_outcomes` Event order | 7 | VAL-033 |
| 7: `not-applicable` is a distinct third availability value | domain-only | 2 defines, 7 confirms | VAL-010, VAL-032 |
| 8: generic Concern subject set | `judgment_concerns.subject_type` CHECK | 8 | VAL-034 |
| 8: three terminal dispositions, single and irreversible | `judgment_concerns` status and disposition discriminator CHECKs | 8 | VAL-035 |
| 8: `carried-forward` binds an existing Issue; `superseded` binds an existing Concern or Finding | `judgment_concerns.disposition_issue_id`, `supersedes_concern_id`, `supersedes_finding_id` FKs | 8 | VAL-035 |
| 8: self-disposal permitted with both Actors recorded | `judgment_concerns.actor_id` and `disposed_actor_id`; domain-only permission | 8 | VAL-074 |
| 8: an open Concern never blocks a guarded action | domain-only | 8 | VAL-034 |
| 8: Deviations append-only with explicit supersession | `deviations` self FK; no update effect rule | 8 | VAL-036 |
| 8: neither Concern nor Deviation is Evidence | `evidence_records.evidence_type` CHECK | 8 | VAL-037 |
| 9: Decision Reference stores no body | `decision_refs` column set | 9 | VAL-038 |
| 9: Source digest verification | domain-only plus `source-digest-mismatch` | 9 | VAL-038 |
| 9: 0..16 protected Decisions per revision | `envelope_protected_decisions` PK and UNIQUE | 9 | VAL-039 |
| 9: blocking exactly the declared affected actions | `affected_actions_json` plus domain | 9 | VAL-039 |
| 9: satisfaction requires the approved reference and Outcome pair | domain-only | 9 | VAL-040 |
| 9: one Stop per unsatisfied protected Decision, resolved with the pair | `stops.code` CHECK; `envelope.activated` and `decision.referenced` effect rules | 1 defines code, 9 uses | VAL-041 |
| 9: `linked_to_decision` becomes implemented | `workflow_transition_requires` rule name plus evaluator | 9 | VAL-042 |
| 10: human Actor type, finalized Activation, capability | `human_outcomes` provenance columns plus domain | 10 | VAL-043 |
| 10: satisfies only a matching class with `approved` | `evidence_records` discriminator CHECK plus domain | 10 | VAL-044 |
| 10: never serialized as deterministic proof | `evidence_records` producer and trust-source CHECK | 10 | VAL-045 |
| 10: `human-outcome-required` Stop conditions | `stops.code` CHECK plus domain | 1 defines code, 10 uses | VAL-046 |
| 10: the four named raising Events for that Stop (DEC-224) | `envelope.activated`, `evidence.linked`, `review.recorded`, and `artifact.availability-observed` effect rules | 10 (task 9 already adds the activation rows) | VAL-046 |
| 9: 0..16 protected Decisions, unique class and ordinal, affected-action subset | `envelope_protected_decisions` PK and UNIQUE plus domain authoring validation | 2 authors, 9 enforces | VAL-006 |
| 11: `human_outcome_resolution:true` | `system_capabilities` materialization | 11 implements, 17 advertises | VAL-047 |
| 11: four new Condition types with exclusive fields | `stops.resolution_condition_json` plus the closed `conditionTypes`/`Condition` registries in `internal/workflow/stop.go` and `internal/journal/current_event_value_catalog.go` | 9 and 10 register the type each raises; 11 registers the other two and owns resolution (R-14) | VAL-048, VAL-078 |
| 11: `attention.resolve` input shape unchanged | domain-only | 11 | VAL-049 |
| 11: only matching Stops resolve | domain-only | 11 | VAL-050 |
| 12: `requirement_refs` with digest and supersession | `requirement_refs` CHECKs and `requirement_refs_one_active` | 12 | VAL-051 |
| 12: `requirement_issue_links` traceability | `requirement_issue_links` plus its guard trigger | 12 | VAL-052 |
| 12: V1 requirement authoring removed | `0018` DROP plus fence and Event removal | 12 | VAL-053 |
| 13: three evidence types, one non-null subject each | `evidence_records` rebuild discriminator CHECK | 1 defines, 6 and 10 use | VAL-025 |
| 13: one record per Obligation and subject, never overwritten | three partial unique indexes on `evidence_records` | 6 | VAL-024 |
| 13: `evidence.evaluate` optional Obligation ID | domain-only | 6 defines, 13 exposes | VAL-008, VAL-055 |
| 14: Human text obligations per leaf | domain-only; no schema | 14 | VAL-058, VAL-060 |
| 14: resume Findings and judgment sections | domain-only; no schema | 14 | VAL-059 |
| 14: CF-040 guidance applied to new views | domain-only; no schema | 14 | VAL-073 |
| 15: final capability matrix | `system_capabilities` materialization and export profile | 17 | VAL-066, VAL-047 |
| 15: owned removals proven absent | `0018` plus `v1-convergence-inventory.md` | 15 | VAL-061, VAL-062, VAL-063 |
| 12: `links` restricted to Issue endpoints, single traceability model | `0018` rebuilt link triggers plus `internal/links` registry | 12 | VAL-054 |
| 15: accepted spec amendments applied | canonical spec files | 1 (schema 8.4), 8 (schema 5 and 3.7, CLI 3), 12 (schema 3.2 and CLI 3), 9 and 13 (CLI 3 and 5), 17 (TRD 4) | wrap-up spec sync |
| public-interface 4.1: thirty grant actions | `envelope_authority_grants` plus profile validation | each listed owner; 18 audits | VAL-067 |
| public-interface 5: added failure codes | `internal/cli/codes.go` plus adapters | 13 | VAL-056 |
| public-interface 5.1: named producers and denial precedence (DEC-226) | domain-only; one code per denied request | 13 audits; each producing task implements | VAL-039, VAL-056 |
| public-interface 5: Phase 23 denials commit nothing (DEC-226) | no denial Event; no `idempotency_results` row | 13 | VAL-057 |
| public-interface 5: added ReasonCodes | evaluation reason precedence in `internal/workflow/evaluation.go` | 6 adds the codes; 10 and 13 exercise them | VAL-008, VAL-026, VAL-033, VAL-046 |
| public-interface 2: added ActivityEvent names | `internal/activity` plus the owning effect rules | 8 and 13 add; 15 removes the four V1 names | VAL-055, VAL-061 |
| public-interface 6: manifest projection schema | `internal/export` plus manifest fixtures | 1, 15, 16 | VAL-064, VAL-065 |
| schema 6.1: progressive Event fence | `internal/journal/project_effects.go` | 1 creates; 4, 5, 8, 9, 10, 12 remove keys | VAL-003 |
| schema 6.2: effect families | `internal/journal/current_event_catalog.go` | each Event's owner; 15 renames the concern rows | VAL-003, VAL-055 |
| schema 6.3: typed result snapshots | idempotency snapshot writers | 13 | VAL-057 |
| schema 7: parity and corruption proof | fresh-workspace rebuild, pragma, and export probes | 1, 16 | VAL-004, VAL-005, VAL-064 |
| schema 7: corrupted or unknown Phase 23 payload rejected, never partially applied | current-Event validator | 16 (all nine types; task 1 can only reach the fence) | VAL-076 |
| schema 1: baseline hunks exactly the accepted seven | `0013` and `0015` file state | 1 | VAL-075 |
| schema 2: link-preserving `evidence_records` rebuild (DEC-223) | `0017` staging of `evidence_artifact_links` | 1 | VAL-005 |
| 8: Concern and Deviation subjects need not resolve (DEC-228) | `subject_type` CHECK only; no reference trigger | 8 | VAL-034 |

## Planning Checks Performed

- Every `contracts.md` section 2 through 15 clause above has an enforcement
  point, one owning task, and at least one assertion.
- Every new closed-domain member named in `public-interface.md` section 2 has a
  matching schema row or an explicit `domain-only` entry here: the three Stop
  codes, the `review-cycles` unit, the `review-rework-exhausted` exhaustion
  action, the `review` Claim mode, the three Obligation types, the three
  evidence types, the `not-applicable` availability value, and the four
  Condition types (DEC-195).
- Every table, index, and trigger in `schema.md` sections 2 and 4 is referenced
  by at least one clause row above; none is created without a consumer.
- Every task numbered 1 through 18 appears at least once as an owner.
- No clause is owned by two tasks without an explicit "defines / uses" split.
