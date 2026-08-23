# Phase 21 V1 Convergence Inventory

> Status: **ACCEPTED FOR EXECUTION - DEC-159**
> Owner: Tech Lead
> Live inventory base: Phase 21 planning HEAD before implementation
> Governing decisions: DEC-153, DEC-156, DEC-157, DEC-158, DEC-159

This is the exhaustive removal/disposition authority for Phase 21. It was
derived from the live `cmd/fid`, `internal`, `migrations`, and journal fixture
tree, not from the roadmap alone. Task 20 may remove only replacement-bound
items listed here. Task 21 proves the final reader/writer/bootstrap/rebuild
state. Phase 23/25-owned state may remain only under the explicit freeze rules.

## 1. Authority Callers And Mutation Flags

| Live surface | Current V1 dependency | Final Phase 21 disposition | Owner |
|---|---|---|---|
| Issue create/update | `--by`, `policy.Authorize`, role Actor/Event fields | retain; rename `--actor`; project-local attribution; keyed Event v2; never kit authority | task 8 |
| Issue transition | `--by`, role policy, role Claim, role-dependent workflow requirements | retain; Envelope-bound paths resolve active Actor/Grant/Claim/Run; terminal-success also requires shared Obligation evaluator; standalone structural transitions use attribution only and cannot enter terminal-success while an active Envelope binds the Issue | tasks 8, 12, 18 |
| Issue AC add | `--by`, role policy | retain as inventory mutation with attribution/key; no verification authority | task 9 |
| Issue AC verify | role-specific policy and `verifier_role`; writes `ac_verifications` | fail closed `capability-unavailable` until Phase 23; no Phase 21 writer | tasks 9, 20 |
| Issue label add/remove | `--by`, role policy | retain with attribution/key | task 9 |
| Label add | `--by`, role policy | retain with attribution/key | task 9 |
| Link add/remove | `--by`, role policy | retain with attribution/key | task 9 |
| Requirement sync | `--by`, role policy | retain with attribution/key; source inventory only | task 9 |
| Workflow define/assign | `--by`, role policy; authoring can carry requirement rows internally | retain state/edge authoring with attribution/key; `--require` stays absent; no role/judgment predicate can be authored | task 9 |
| Concern raise/address | role policy and role attribution; completion foreign key | fail closed until Phase 23; list remains read-only legacy judgment inventory | tasks 9, 20 |
| Issue claim/unclaim and Claim list | role authority; V1 `claims` table/readers/writers | replace with generic Claim leaves/table; old spelling absent | task 10, final table rename task 20 |
| Run record | role policy; V1 `agent_runs` writer | replace with generic Runs; old spelling absent; historical table frozen | tasks 11, 20 |
| Dev start/gate/finish/block | role policy, V1 Claim/Run/Gate/completion choreography | temporarily retained through Group 1 only; remove after trusted replacements pass | task 20 |
| Issue complete | role policy; writes `task_completions` | temporarily retained through Group 1 only; remove after shared evaluator/Envelope completion passes | task 20 |
| Rebuild | `--by`, `policy.Authorize`, ActorRole lock metadata | retain projection operation with no Actor flag or authority claim | tasks 2, 21 |
| Maintenance vacuum | role policy | fail closed until Phase 25 maintenance authority | tasks 9, 20 |
| Maintenance check | reads legacy tables | retain; readers become integrity/frozen-state observations only | tasks 20, 21 |
| Stats tokens/cost | `--by` means grouping and reads `agent_runs.role` | rename `--group-by`; `actor` exposes historical `actor_attribution`; no authority/capability inference | tasks 9, 20 |
| Policy show and action registry | public V1 role/action inspection plus policy used by dev/completion | remove public command in task 9; isolate implementation to enumerated legacy callers; delete implementation only after all caller dispositions pass in task 20 | tasks 9, 20 |
| Bootstrap/init | genesis role Effects and rebuild as `human` | Event-v2 empty domain Effects, service Actor, no role projection; bounded nonpersistent V1 policy vocabulary bridge through task 19; final fresh init proven after 0014 | tasks 2, 20, 21 |

### 1.1 Exact live caller/file classification

The planning-base `policy.Authorize|KnownRoles|ActorRole|Role:` and public
`--by` inventories classify every production match as follows. A later file
move remains with the named surface/task; an unlisted match is a planning
contract failure.

| Disposition / owner | Exact planning-base production files |
|---|---|
| role-free bootstrap/rebuild/Event lock metadata; tasks 2, 20, 21 | `cmd/fid/rebuild.go`; `internal/bootstrap/init.go`; `internal/journal/lock.go`; `internal/journal/projection.go` |
| retained Issue attribution; task 8 | `cmd/fid/issue_create.go`; `cmd/fid/issue_helpers.go`; `cmd/fid/issue_mutation_response.go`; `cmd/fid/issue_transition.go`; `cmd/fid/issue_update.go`; `internal/issues/create.go`; `internal/issues/issue.go`; `internal/issues/transition.go`; `internal/issues/update.go`; `internal/issues/update_request.go` |
| retained/fail-closed AC, judgment, label, link, requirement, workflow, stats, maintenance; tasks 9, 20 | `cmd/fid/concern.go`; `cmd/fid/issue_ac.go`; `cmd/fid/issue_label.go`; `cmd/fid/label.go`; `cmd/fid/link.go`; `cmd/fid/maintenance.go`; `cmd/fid/requirement.go`; `cmd/fid/stats.go`; `cmd/fid/workflow.go`; `internal/issues/ac.go`; `internal/issues/ac_verify.go`; `internal/issues/concern_address.go`; `internal/issues/concern_raise.go`; `internal/issues/label.go`; `internal/issues/transition_rules.go`; `internal/labels/add.go`; `internal/labels/types.go`; `internal/links/add.go`; `internal/links/remove.go`; `internal/links/types.go`; `internal/maintenance/vacuum.go`; `internal/requirements/sync.go`; `internal/requirements/types.go`; `internal/workflow/assign.go`; `internal/workflow/define.go` |
| generic authority replacement then final legacy deletion; tasks 7, 9, 20 | `cmd/fid/policy.go`; `internal/policy/authorize.go`; `internal/policy/matrix.go`; `internal/policy/policy.go` |
| V1 Claim replacement; tasks 10, 20 | `cmd/fid/claim.go`; `cmd/fid/issue_claim.go`; `internal/claims/acquire.go`; `internal/claims/claims.go`; `internal/claims/release.go`; `internal/claims/renew.go` |
| V1 Run replacement; tasks 11, 20 | `cmd/fid/run.go`; `internal/runs/record.go` |
| dev/completion removal after replacement; task 20 | `cmd/fid/dev.go`; `cmd/fid/issue_complete.go`; `internal/dev/block.go`; `internal/dev/dev.go`; `internal/dev/finish.go`; `internal/dev/gate.go`; `internal/dev/start.go`; `internal/issues/complete.go`; `internal/issues/complete_support.go` |

The exact command adapter inventory above covers all current mutation/grouping
flag declarations. The domain/file inventory covers every current role-policy
call or role-bearing writer/lock value. Tests and fixtures are not production
callers; their final retained/replaced/removed Event disposition is section 3.

No final retained mutation accepts `--by`. The unrelated grouping flag also
disappears. Actor-looking strings in activity, legacy concern, AC, completion,
telemetry, or journal facts are attribution only. They never satisfy an
Authority Grant, Evidence trust source, Human Outcome, protected transition,
or Envelope completion.

## 2. Workflow Preconditions And Judgment State

| Rule/state | Current readers/writers | Final disposition | Owner |
|---|---|---|---|
| `all_acs_verified` with role params | workflow registry and Issue transition read `ac_verifications` | migration 0014 deletes seeded/legacy rule rows; generic transition never reads it; Phase 23 redesigns judgment | tasks 12, 20 |
| `has_completion_report` | Issue transition reads `task_completions` | migration 0014 deletes rule rows; shared evaluator uses Obligation + terminal Issue, never completion table | tasks 12, 18, 20 |
| `concerns_addressed` | Issue transition reads `concerns` | migration 0014 deletes any such rows; Phase 23 owns judgment | tasks 12, 20 |
| `all_child_tasks_done_or_cancelled` | Issue workflow structural reader | retain for standalone Epic inventory; no role or kit authority | task 9 |
| reserved `required_gates_current` and `has_completion` | registry constants reject authoring and have no evaluator/seeded row | retain only as rejected legacy identifiers; Phase 21 Gate Evidence/Envelope completion never consult them | tasks 9, 12, 20 |
| frozen `linked_to_decision` escalation row | seeded V1 escalation workflow; no Issue evaluator | retain as inaccessible Phase 23 judgment history only; it cannot enter an Issue workflow, authorize kit action, satisfy Evidence, or guard Phase 21 transition/completion | tasks 9, 20 |
| terminal-success guard | absent in V1 | add one shared Obligation evaluator used by terminal Issue transition and Envelope completion | tasks 12, 17, 18 |
| `workflow_transition_requires` authoring | internal builder supports generic requirement Effects although public `--require` is absent | new Issue authoring is structural-only; reject every role/judgment rule; delete displaced Issue rows in 0014; freeze the inaccessible Phase 23 escalation row | tasks 9, 20 |

`ready` may report structural workflow actionability. For an Envelope-bound
terminal-success transition it returns the same evaluator reasons and cannot
derive readiness from role-shaped AC/completion/concern records.

## 3. Event Families

| Family in live tree | Final disposition | Owner |
|---|---|---|
| `project.initialized` | retain as exact Event-v2 bootstrap fixture with zero domain Effects and aggregate sequence 1 | task 2 |
| `issue.created`, `issue.updated`, `issue.transitioned` | retain/convert to v2; transition uses shared evaluator for terminal success | tasks 8, 12, 18 |
| `issue.ac_added` | rename public Event value to `issue.ac-added`; retain inventory semantics | task 9 |
| `issue.ac_verified` | no new Phase 21 writer; retained history replayable but not authority; any accepted activity/read rendering normalizes the public value to `issue.ac-verified`; deletion/redesign Phase 23 | task 20 |
| `issue.label.added`, `issue.label.removed`, `label.added`, `link.added`, `link.removed`, `requirement.synced`, `workflow.defined`, `workflow.assigned` | retain/convert v2 and project-local attribution | task 9 |
| `issue.concern_raised`, `issue.concern_addressed` | no new Phase 21 writer; history is non-authoritative; accepted activity/read rendering uses `issue.concern-raised` / `issue.concern-addressed`; Phase 23 | task 20 |
| `issue.claim.acquired`, `issue.claim.renewed`, `issue.claim.released` | remove V1 fixtures/writers; generic `claim.*` replaces | tasks 10, 20 |
| `run.recorded` | remove writer/fixture; generic `run.started|finished` replaces | tasks 11, 20 |
| `dev.started`, `dev.gate.recorded`, `dev.finished`, `dev.blocked` | may remain through task 19; remove source/fixture acceptance in task 20 | task 20 |
| `issue.completed` and V1 completion effects | may remain through task 19; remove writer/fixture acceptance in task 20 | task 20 |
| Activation, Registration, Envelope, generic Claim/Run, Attempt/Receipt, Evidence, Artifact availability | add exact v2 families from `contracts.md` §14 | tasks 2-18 |
| `activity_log.event` values such as `status_changed`, `ac_verified`, `gate_recorded`, `task_completed`, and concern/claim/link variants | retained rows remain historical attribution only; new writers and accepted reads use the exact kebab-case `ActivityEvent` registry; rebuild normalizes frozen values deterministically | tasks 8, 9, 20, 21 |

Event names are contract values and therefore kebab-case within dotted family
segments. A v1 historical name may replay during development only until task
20 removes its final accepted fixture; it is never emitted by the final binary.

## 4. Table Readers And Writers

| Table | Live readers/writers | Final Phase 21 rule | Owner |
|---|---|---|---|
| `roles` | bootstrap writer; policy journal scanner; SQL/tests | task 2 removes the writer/Effect and permits only the exact nonpersistent V1 compatibility lookup in section 5; task 20 drops the table and lookup; no final reader/writer/genesis Effect | tasks 2, 9, 20, 21 |
| V1 `claims` | Issue Claim, dev, transition, maintenance readers/writers | generic replacement uses `work_claims` through task 19; task 20 drops V1 and renames generic table | tasks 10, 12, 20 |
| `gate_runs` | dev gate writer; dev/Issue detail/transition/maintenance readers | may remain replacement-bound through task 19; drop in 0014 after Receipt/Evidence | tasks 15-18, 20 |
| `task_completions` | Issue complete/dev finish writers; Issue detail/transition/concern/maintenance readers | remove all writers and authority/completion readers; freeze table for Phase 23 FKs; maintenance may report nonempty/frozen writes | tasks 18, 20, 21 |
| `agent_runs` | run record/dev writers; dev/stats/maintenance readers | remove all writers and generic-Run readers; freeze for Phase 25 telemetry; stats may expose historical attribution only | tasks 11, 20, 21 |
| `ac_verifications` | AC verify writer; Issue detail/transition/maintenance readers | no Phase 21 writer or terminal/authority read; retain read-only Phase 23 state and maintenance integrity observation | tasks 9, 12, 20 |
| `concerns` | concern writers; Issue detail/transition/list/maintenance readers | list/history only; no Phase 21 writer or terminal/authority read; Phase 23 | tasks 9, 12, 20 |
| `workflow_transition_requires` | workflow authoring/registry and Issue transition | structural Issue rows plus explicitly frozen non-Issue Phase 23 history only; role/judgment Issue rows deleted in 0014; no generic terminal evaluator read | tasks 9, 12, 20 |

Maintenance is not an authority reader: it may inspect frozen/deleted-state
violations solely to emit deterministic integrity findings. `fid sql` is an
unstable Human diagnostic and cannot be cited as product authority/evidence.

## 5. Group 1 Temporary Replacement Boundary

After task 12 and before task 13, the following replacement-bound state may
still exist because final deletion depends on Group 2 trusted execution:

- tables: `roles`, V1 `claims`, `gate_runs`, `task_completions`, `agent_runs`,
  `ac_verifications`, `concerns`, and role/judgment rows in
  `workflow_transition_requires`;
- commands/writers: `fid dev start|gate|finish|block`, `fid issue complete`,
  and their `dev.*`/completion Events; the `roles` table may remain but task 2
  has already removed every bootstrap role writer/Effect; and
- source packages: `internal/dev`, legacy completion files, and frozen-table
  readers needed for removal conversion; and
- internal V1 policy implementation reachable only by the enumerated
  dev/completion compatibility callers, with `fid policy` already absent; and
- the task-2 V1 compatibility role vocabulary inside `internal/policy`, which
  is available only when the exact role-free Event-v2 bootstrap exists and no
  legacy role Effect exists. It recognizes the closed seven V1 names for the
  retained V1 matrix, persists nothing, and is unavailable to every new
  generic path.

Checkpoint PASS does not assert their absence. It proves by dependency/call
graph and adversarial tests that no new Activation, Registration, Envelope,
Grant, generic Claim/Run, nonterminal execution-bound transition, capability,
or idempotency path consults them. It also proves judgment mutations are
fail-closed. Final no-command/no-writer/no-role-bootstrap absence is task 20;
final-schema init/rebuild/export/maintenance parity is task 21.

## 6. Final Removal And Deferral Inventory

Task 20 removes: `internal/dev`; legacy `internal/runs`; `fid policy`; `fid
dev`; legacy Issue Claim/unclaim and Run record leaves; Issue complete files;
V1 policy implementation; every writer of `roles`, V1 `claims`, `gate_runs`,
`task_completions`, `agent_runs`, `ac_verifications`, and `concerns`; accepted
V1 dev/completion/claim/run fixtures; every mutation `--by`; stats `--by`; and
role/judgment workflow requirement rows; and the task-2 nonpersistent V1 role
vocabulary bridge. Migration 0014 drops `roles`, V1
`claims`, and `gate_runs`, renames `work_claims`, and deletes the displaced
requirement rows.

Task 20 does not drop `task_completions`, `agent_runs`, `ac_verifications`, or
`concerns`. Those tables are frozen for Phase 23/25 foreign-key/history
ownership; they have no final Phase 21 writer and cannot satisfy protected
semantics. Task 21 proves the final schema, fresh init, delete/rebuild,
maintenance, capabilities, evaluation/completion/resume, and export parity
without `roles`, V1 Claims/Gates, or a V1 bootstrap branch, then runs the
100k Event-v2 rebuild threshold and parity check on the same final schema.

Phase 22 owns recovery, Attention, Stops, budget use, respawn, and import.
Phase 23 owns AC verification, concerns, Review, Human Outcome, and deletion or
migration of `task_completions`/judgment state. Phase 25 owns telemetry and
`agent_runs` deletion/migration.

## 7. Deterministic Inventory Commands

The QA plan must bind exact bounded variants of these inventories at the
Group 1 checkpoint, task 20, task 21, and final QA:

```text
rg -n --glob '*.go' --glob '!**/*_test.go' 'policy\.Authorize|KnownRoles|ActorRole|Role:' cmd/fid internal
rg -n --glob '*.go' --glob '!**/*_test.go' 'String(Var)?\([^\n]*"by"|GetString\("by"|--by' cmd/fid internal
rg -n --glob '*.go' --glob '!**/*_test.go' '(FROM|INTO|UPDATE|JOIN|DELETE FROM) (roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)\b|Table:\s*"(roles|claims|agent_runs|task_completions|gate_runs|ac_verifications|concerns|workflow_transition_requires)"' cmd/fid internal
rg -n --glob '*.go' --glob '!**/*_test.go' 'EventType:\s*"' cmd/fid internal
```

The report must classify every match against this file. An unclassified match
fails; a raw no-match is not sufficient where a frozen or diagnostic reader is
explicitly retained. At final task-20 HEAD, the first two commands are absence
controls and must terminate with the predeclared `exit_code: 1`; the broad
table-reader and Event-family commands are classified inventories and must
terminate `exit_code: 0`. Every line produced by either classified inventory
must map to an explicit retained, frozen, diagnostic, or removal disposition,
and the completion report must record `UNCLASSIFIED=0`.
