# Phase 24 Executable Inventory

> **Status: ACCEPTED with the contract (DEC-277).**
> Measured at `4b86710` (branch `phase-24`). Phase 23's five planning BLOCKs
> came from registries missed during planning. This file therefore lists
> every closed registry and fixture family the phase reaches, each with the
> command that regenerates it. Before task close, a Developer reruns the
> command for each registry the task owns and records the result. QA reruns
> them all at every review window. Task numbers refer to the approved
> task pack (`tasks/`).

## 1. `--input` Leaves (39)

Regenerate the live list (scratch binary; no workspace needed):

```bash
go build -o "$S/fid" ./cmd/fid
for g in $("$S/fid" --help | awk '/Available Commands:/{f=1;next} /^$/{f=0} f{print $1}'); do
  for s in $("$S/fid" "$g" --help | awk '/Available Commands:/{f=1;next} /^$/{f=0} f{print $1}'); do
    "$S/fid" "$g" "$s" --help 2>/dev/null | grep -q -- '--input' && echo "$g $s"
  done
done
```

The expected result is exactly the 39 rows below. `Shape` is `changes` when
public-interface section 4 amends the input. `No Short` marks leaves whose
help has no one-line description today.

| Leaf | Shape | No Short | Descriptor and convergence task |
|---|---|---|---|
| activation abort, commit, component, finalize, prepare, recover, recover-complete, resource (8) | same | - | 9 |
| artifact verify | same | - | 8 |
| attention resolve | same | yes; help says "Run JSON input" | 3 |
| claim acquire | changes (`worktree_id`) | - | 6 |
| claim expire, release, renew | same | - | 6 |
| concern record | changes (references) | yes | 4 |
| concern dispose | same | yes | 4 |
| decision record-reference | changes (source) | yes | 4 |
| deviation record | changes (references) | yes | 4 |
| envelope create | changes (scope, Sources) | - | 6 |
| envelope complete | changes (`repository`) | - | 6 |
| envelope activate | same | - | 6 |
| evidence link | changes (`target_repository`) | - | 8 |
| evidence evaluate | changes (`target_repository`) | - | 8 |
| finding address | same | yes | 5 |
| gate register | changes (`worktree_id`, `repository`) | - | 6 |
| gate run | changes (`repository`) | - | 8 |
| human-outcome record | changes (references) | yes | 4 |
| requirement record-reference | changes (source) | yes | 4 |
| requirement link | same | yes | 4 |
| review record | changes (`repository`, Sources) | yes | 5 |
| run start, close, abandon | changes (`repository`) | - | 7 |
| run observe-budget | changes (`measurement_ref`) | - | 7 |
| run observe, checkpoint, recover | same | - | 7 |
| stop raise, resolve | same | yes | 3 |

Read leaves without a Short today: `attention list`, `concern list`,
`decision list`, `deviation list`, `finding list`, `human-outcome list`,
`requirement get`, `requirement list` and `stop list`. Task 3 owns Short
descriptions for the whole live tree. A pending-descriptor fence lists the 36
leaves task 3 does not adopt. Each owning task removes its leaves from the
fence, and task 9 removes the empty fence (the DEC-192 pattern).

## 2. Closed Registries

| ID | Registry | Regenerate | Count at `4b86710` | Owner |
|---|---|---|---|---|
| R1 | Request-digest call sites (`RequestDigest:`) | `rg -n "RequestDigest:" internal cmd -g '!*_test.go'` | 43 sites in 32 files | 2 |
| R2 | Request-digest helpers | `rg -n "func [a-zA-Z]*(RequestDigest\|requestDigest\|judgmentDigest)\(" internal -g '!*_test.go'` | 11 helpers | 2 |
| R3 | Idempotency storage: `idempotency_results` DDL, idempotency Effect validation, projection | `rg -n "idempotency_results" migrations internal -g '!*_test.go'` | 0013 DDL; `internal/journal/current_event_idempotency.go`, `current_event_validation.go`, catalog rows | 1 |
| R4 | Migrations and projection schema constants | `ls migrations; rg -n "ProjectionSchemaVersion\|projection_schema" internal -g '!*_test.go'` | latest 0018; schema 18 | 1 |
| R5 | Closed `--input` decoders | `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid` | 7 helpers in 14 command files | 3-9 |
| R6 | Repository comparison and validation | `rg -n "func (same[A-Za-z]*Repository\|validateRepositoryBinding\|validEnvelopeRepository\|validEnvelopeSource\|validEvidenceRepository\|validateRunTerminalRepository\|verifyDecisionSource\|importedGateDestination)" internal -g '!*_test.go'` | 4 `same*Repository`, 6 validators | 4-8 |
| R7 | Caller-repository writers into stored facts or Events | `rg -n "Request\.Repository\|request\.Repository\|Request\.TargetRepository\|request\.RepositoryScope\|Request\.WorktreeID\|request\.WorktreeID" internal -g '!*_test.go'` | registration, execution (Attempt, Receipt, Events), evidence, envelope, run, run_close, claims, review | 6-8 |
| R8 | Failure codes and details | `internal/cli/codes.go` (`definition(` rows); conformance `internal/cli/*codes_test.go` | `repository-unavailable` absent | 4 (unavailable, source), 5 (mismatch), 2 (conflict), 3 (invalid-argument) |
| R9 | Human failure renderers | `rg -n "func (write[A-Za-z]*Failure\|render[A-Za-z]*Error)\(" cmd/fid -g '!*_test.go'` | 41 functions | 10 |
| R10 | Human failure assertions in tests | `rg -n '"error: ' cmd internal -g '*_test.go'` | 32 in 13 files | 10 |
| R11 | Human views printing repository identity | `rg -n -i "worktree\|dirty digest\|HeadOID\|head_oid" cmd/fid -g '!*_test.go'` | `run_inspect_human.go`, `resume_human.go`, `p23_human.go`, `review.go` | 11 |
| R12 | Stop/Attention Human renderers | `cmd/fid/stop_human.go`, `resume_human.go`, `resume_human_guidance.go` | 3 files | 11 |
| R13 | Judgment Human renderers, disclaimers, leak guard | `rg -n "p23CheckJudgmentQualifiers\|PRIVATE_" cmd` | `p23_human.go`; `integration/p23_human_leaf_pairs_test.go:123`, `p23_human_fixtures_test.go:143` | 13 |
| R14 | Live leaf inventory test | `rg -n "func TestLiveLeafCommandsMatchContractInventories" cmd/fid` | `contract_surface_test.go:128` | 12 (DEC-274) |
| R15 | Cobra Short descriptions | the section 1 walk with `head -1` per leaf | 9 read + 12 `--input` leaves start with `Usage:` | 3 |

## 3. Fixture And Test Convergence

Displaced assertions follow DEC-162/165/166. Each Part 1 task converges the
retained tests its accepted change displaces and records per-row provenance
(introducing commit, phase and task).

| Family | Regenerate | Count | Notes |
|---|---|---|---|
| Tests naming caller repository fields | `rg -l "worktree_id\|capture_source\|repository_scope\|target_repository" cmd internal -g '*_test.go'` | 56 files (29 in `cmd/fid/integration`) | mostly via helpers below |
| Repository-input helpers | `rg -n "func [a-zA-Z0-9_]*(Repository\|Repo)[a-zA-Z0-9_]*\(" cmd -g '*_test.go'` | `runRepositoryForTest`, `testEnvelopeRepository`, `p22CorpusRepository`, `p23CurrentRepository`, `cloneP23Repository`, `task18EvidenceRepository` and others | converging a helper converges its callers |
| Shared fake/fid conformance adapters | `ls cmd/fid/integration/p2*_fake_test.go` | Phase 22 fake corpus and activation fake | the fake must accept the converged inputs (roadmap section 7) |
| Pinned literal digests | `rg -n '"sha256:[0-9a-f]{64}"' cmd internal -g '*_test.go'` | 8 in 7 files | recompute only if the input as written changed; never re-baseline to an observed value |
| Gate fixtures writing Artifacts | `rg -n "gitignore" cmd/fid/integration -g '*_test.go'` | fixtures ignore `*.json` Artifacts | under real observation a non-ignored Artifact dirties the worktree and stales the gate's own Evidence; fixtures must keep Artifacts ignored |
| Phase 23 Human fixtures | `qa/human-fixtures.md` in Phase 23 | `.cache/phase-23-human/` | regenerate after task 13; Claims expire |

## 4. Known Risks

- **Freshness becomes real (task 8).** Today Receipt and Evidence repositories
  are the caller's block, so Phase 21-23 journeys stay "fresh" by
  construction. After task 8 they compare real observations. A journey that
  edits files between a gate run and its guarded action will now (correctly)
  see `evidence-stale`. Such a row is convergence of the accepted change, not
  a defect to weaken.
- **Imported destinations.** `importedGateDestination`
  (`internal/evidence/execution.go:338`) lets an imported workspace run a
  gate whose Registration names the source worktree. Task 8 must preserve
  this under observation (Phase 22 import semantics).
- **`run.close` base OID.** Today the caller's `base_oid` is retained and
  checked for existence. After task 7 `base_oid` is the observed HEAD, as at
  every other observation site.
- **Digest helpers that mutate requests.** `envelope.complete` overwrites
  `request.Repository` after hashing. The field disappears in task 6. Task 2
  pins the rule "digest before observe; never mutate afterwards".
