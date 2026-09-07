# Phase 22 Final QA Preflight

Reviewer: independent QA Lead `checkpoint12_final`
Date: 2026-09-08
Status: PREPARATION ONLY — tasks 12–18 and final phase gates remain pending.

This is the early QA-84 reconciliation requested while task 12 completes. It
does not declare final QA, re-attest another session's execution, or replace a
product gate. Report tokens and termination fields belong to the originating
Developer/root sessions. This reviewer inspected report bytes and task specs;
no earlier product command was rerun for this preflight.

## Task 1–11 exact-command reconciliation

The exact-command sections contain 70 required occurrences. All are represented
by closing successful receipts with a nonempty token and observed-termination
field. Required repetitions are counted separately. Every task requires:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
```

| Task | Occurrences | Additional obligations inspected | Closing report receipts / result |
|---|---:|---|---|
| 1 | 12 | Environment six commands and 100k Event rebuild benchmark below | Lines 60–114; presence complete |
| 2 | 7 | Workflow package coverage before edits and after implementation | Baseline line 48; recovery baseline line 61; closing lines 120–131; PF-01 encoding finding |
| 3 | 5 | None | Closing QA-rework lines 90–94; PF-01 |
| 4 | 6 | Exact equivalent-Stop regression with count 5 below | Lines 67–73; PF-01 |
| 5 | 5 | None | Lines 59–63; presence complete |
| 6 | 5 | None | Lines 96–102; failed race receipt line 98 preserved; presence complete |
| 7 | 7 | Evidence package coverage before edits and after implementation | Original baseline line 57; recovery baselines 58/83/95; closing lines 105–112; PF-01 |
| 8 | 8 | Exact terminal test count 15 and two full benchmark invocations below | Closing lines 94–101; both distinct benchmark receipts present; PF-01 |
| 9 | 5 | None | Final-source lines 65–69; earlier quotations attributed; presence complete |
| 10 | 5 | None | Lines 53–57; presence complete |
| 11 | 5 | None | Lines 53–57; presence complete |

Line numbers refer to the inspected pre-correction report bytes. Exact extras:

```text
# task 1
go version
go env GOOS GOARCH
uname -srm
nproc
grep -m1 '^model name' /proc/cpuinfo
grep -m1 '^MemTotal:' /proc/meminfo
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal
# task 2: before edits, then after implementation
go test -race -cover ./internal/workflow
# task 4
go test -race ./internal/workflow -run '^TestEquivalentStopsShareAttentionButKeepDistinctFacts$' -count=5
# task 7: before edits, then after implementation
go test -race -cover ./internal/evidence
# task 8
go test -race -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' -count=15 ./cmd/fid
go test -bench . -benchmem ./...
go test -bench . -benchmem ./...
```

No task declares a nonzero successful command. Task 8's eight historical
BLOCKED rows explicitly say the commands had not run during the cancelled
Receipt contradiction; they have no invented termination. Its later eight
PASS receipts close the gate after DEC-195.

Task 2 retains the 27.4% original workflow baseline, recovery baselines, failed
full gate, DEC-189 explanation and final 33.7% scoped result. Task 7 retains the
25.2% original evidence baseline and recovery baselines; its closing scoped
result is 24.0%. These are attributed historical measurements, not current-tree
measurements. Tasks 1,2,6,7,8,9,10,11 retain required-gate failure/deviation
history. Task 9 distinguishes originating Developer quotations from final
continuation receipts and root's coverage formatting correction. Tasks 10/11
retain report-control failures and root-assistance attribution.

## PF-01 — encoded tokens in five reports (original finding)

Classification: bounded report repair, no VAL / QA-84. This preparation finding
is not a final-review verdict. Root accepted the display-only repair ownership.

At observed HEAD `0699b660de9890bd769ac2dbceac21934f46d1e4`, these report
lines contained literal backslash-plus-t bytes in captured Go/coverage output.
The displayed tokens were encoded rather than verbatim. No missing execution
was inferred: commands, outcomes, values and termination fields were present.

| Report | Exact affected original lines |
|---|---|
| task-2-done.md | 48,54,57,61,63,69,71,73,75,77,88,93,95,110,114,115,117,118,120,126,128,130 |
| task-3-done.md | 56,60,66,69,72,75,76,79,85,86,87,92,93,94 |
| task-4-done.md | 49,54,59,62,70,71,72,73 |
| task-7-done.md | 57,58,63,66,67,70,71,72,76,79,80,81,82,83,86,91,94,95,100,101,105,110,111,112,127 |
| task-8-done.md | 82,96,97,98,99,100,101 |

The authorized repair decodes only known receipt display escapes, preserving
all original commands, outcomes, times, percentages, termination fields,
ordering, attribution and failure history. It records that this is rendering
correction, not new execution or a status upgrade. No product/test/gate change
is assigned. Resolution must separately cite the correction commit and byte
comparison against these originals.

### PF-01 resolution

RESOLVED by independent byte comparison of root's report-only commit
`be53edadaec38d12c7df6c0eaf43f1bf003f476d`
(`docs(phase-22): decode historical gate token rendering`). The changed-path
set is exactly the five named reports. For each, removing only the appended
`Root report display correction — 2026-09-08` paragraph leaves the original
bytes with literal backslash-t decoded to tab, and no other original-byte
change. Each appended paragraph explicitly preserves originating-session
attribution and disclaims new execution. The original finding remains above.

The Python byte-comparison command used `git diff-tree --no-commit-id
--name-only -r` to assert the exact five-file set, `git show <commit>^:<path>`
and `git show <commit>:<path>` for each version, and asserted:

```python
body, notes = after.split('\n## Root report display correction — 2026-09-08\n', 1)
assert body.rstrip('\n') == before.replace('\\t', '\t').rstrip('\n')
assert 'no new execution or attestation' in notes
```

Observed final token:
`PF-01 RESOLVED: five report-only changes exactly decode tabs and retain attribution; no command, outcome, value or termination changed`
(terminated `exit_code: 0`). The root's applicable report-control receipts are
separate Runner evidence; they have not been executed or re-attested here.

## Command and feasibility validation

Working directory: /home/mark/projects/fid. All referenced task/report files
exist and their exact-command sections contain concrete commands. No historical
binary or deleted test directory is needed to repeat this static inspection.

A Python stdlib inventory extracted fenced commands only from each exact-command
section, retained duplicates with Counter, and emitted each matching receipt's
line, status, token presence and termination-marker presence, plus report SHA-256
and latest report commit/author from git log. It terminated exit_code: 0.
All 70 occurrences had closing successful receipts; prefixed pre-edit baseline
rows were additionally inspected directly. This is file-state evidence, not
product PASS evidence.

The encoding inventory was:

```bash
python3 - <<'PY'
from pathlib import Path
for n in (2,3,4,7,8):
    p=Path(f'project-docs/docs/phases/phase-22/tasks/task-{n}-done.md')
    hits=[str(i) for i,line in enumerate(p.read_text().splitlines(),1)
          if '\\t' in line]
    print(f'{p}: literal-backslash-t lines '+','.join(hits))
PY
```

Observed final token:
`project-docs/docs/phases/phase-22/tasks/task-8-done.md: literal-backslash-t lines 82,96,97,98,99,100,101`
(terminated `exit_code: 0`). Original line lists are retained above; a
later correction does not erase the finding.

## Remaining work

Verified task-12 completion and its independent checkpoint take priority.
Tasks 12–18 still require reconciliation. Final QA must execute its complete
phase commands and deterministic checklist, verify review/repair closures, and
record each validation outcome independently. The separate human checklist has
six blank judgments; task 18 must supply, and QA must verify, its retained
fixture manifest before Human execution. No Human outcome is recorded here.


## Post-task-12 update and final QA queue — 2026-09-08

This append supersedes the earlier pending task-12 status without rewriting its
history. Tasks 1–12 are COMPLETE. Review4 is PASS at
`3addf7d6886b9c2b94ed8e5baeeff8c2dcbe1834`; the additional independent
checkpoint is PASS at `a9bf3bf48442d587f70b824a54254f09ccd61d66`.
Task 13 is active. This queue is preparation only: no final-phase product gate,
timing qualification, final VAL outcome or Human judgment is claimed.

Task 12 adds five reconciled exact gate occurrences to the preceding 70.
The checkpoint independently read its final report
`f57b4aa307fe5d3bb700ac04e8229070422b2d0e` and Developer addendum
`833b4fc2cefcdade9fd36c74e9b14573caa682fa`, verified product commit
`c2d80980b9bfa9c1c20723cc7751526eb8baf2bf`, and ran report lint to
observed exit_code 0 / OUTCOME COMPLETE. Those attributed task-gate receipts
remain distinct from the checkpoint's own fresh focused product checks.

### DEC-202 import qualification allocation

Verified in decision DEC-202, task-14 Binding Constraints/Scope/Acceptance
Criteria and commit `af820bf0a95d290fab9ba9777e50f717dbe0d437`
(`docs(phase-22): bind imported destination recovery scope`).
Contracts section 10 supplies the unchanged portable-recovery requirement.

Task 14 owns the added destination reconciliation, inspection, recovery,
observation, Run and evaluation paths, with tasks 3/4/5/6/12 provenance.
Installed-import delivery metadata may qualify the preserved Project/history
at the actual destination; it is not an authority grant. Required positive
evidence is fresh destination liveness plus exact reconciliation, then a new
successor Claim/Run and qualifying execution. Required negatives include
unqualified imported mutation, replay of an imported observation, ordinary
foreign-worktree requests without the installed association, unknown liveness
and existing applicable authority/Stop/Budget/freshness guards. Existing Git
metadata and source files must survive import into a target with no .fid.
Original history, repository values, policy, Event order/bytes and proof IDs
remain immutable. Task 15 carries the route through the real CLI journey.

No missing allocation or new protected requirement was identified in this
check. The private publication marker remains delivery metadata excluded from
canonical export; no public/schema/Event member is added by DEC-202.

### Remaining task command obligations

The current task specs, including accepted amendments, have 126 ordered command
occurrences overall: 75 reconciled for tasks 1–12 and 51 future occurrences
below. Historical QA-plan readiness totals predate task 4's added repeated Stop
regression and DEC-196's task-18 generator command; use each current task spec.
No required command is a placeholder.

| Task | Ordered required commands | Occurrences | Additional final evidence obligation |
|---|---|---:|---|
| 13 | Five-command base printed above | 5 | Actual ZIP boundary matrix 056-a..k, measured expansion, complete valid controls, semantic rejection and no-effect snapshots |
| 14 | Five-command base | 5 | Real lock/stage/reconstruct/fsync/rename/publication faults and retries; DEC-202 qualified/unqualified destination route |
| 15 | Five-command base | 5 | Source/live/rebuilt/imported parity with declared exclusions, re-export and real CLI journey; then Review5 |
| 16 | Base, six reference-instance commands, one exact 100k benchmark | 12 | Schema-16 removal/classification, provenance and fresh early timing sample |
| 17 | Base, six reference commands, three exact 100k processes, full benchmark twice, build | 17 | Final schema/parity/maintenance and exact capability matrix; each timing sample below 30000 ms |
| 18 | Base, build, exact retained Human-fixture generator below | 7 | Final shared corpus, 23 independently pinned actions, emptiness guard, Human/JSON parity, retained fixture survival; then Review6 |

The exact reference command and instance capture are already printed above
under tasks 1–11. They stay byte-identical. Task 16's early sample does not
replace any of task 17's three samples or the separate final-QA three samples.
Task 17's reference identity remains Go 1.26.2; linux/amd64;
Linux 6.6.114.1-microsoft-standard-WSL2 x86_64; six logical CPUs;
AMD Ryzen 7 5700G with Radeon Graphics; MemTotal 8133236 kB.

Task 18's additional exact delivery command, after its build:

```text
FID_P22_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-22-human go test -tags integration ./cmd/fid/integration/... -run '^TestP22HumanFixtures$' -count=1 -v
```

The full tagged gate must also execute generator setup modes and a subprocess
smoke check. The retained manifest must name a durable binary, exact projects
and entity IDs, and neutral inspection commands usable after generation exits.
Fresh uniquely named children must preserve existing fixtures. QA must run the
delivered inspection commands before replacing the checklist's pending locations.
Only the six HUMAN-01a..03b prompts belong in the user-facing checklist; expected
judgments stay separate and every Human evidence/outcome cell remains blank.

### Final execution order after verified task 18 and Review6

1. Verify all 18 completion reports/commits, reviews 1b/2b/3/4/5/6 and the
   additional task-12 checkpoint, including every earlier REPAIR's independent
   closure. Reconcile all 126 command occurrences and retained-test provenance.
2. Execute the full phase gate in its accepted order, recording each command's
   own token and actual termination. This is the five-command base, then
   `go test -bench . -benchmem ./...`, then
   `go build -o /tmp/fid ./cmd/fid`. QA-78 separately requires a second
   consecutive clean full benchmark invocation. Preserve failed attempts;
   never select a green from a failing series.
3. Capture all six reference-instance commands in that final-QA session.
   Execute three consecutive fresh processes of the exact 100k benchmark,
   preserving all three parity/self-asserting transcripts and each measured
   ms/rebuild. Each must remain strictly below 30000 ms on the accepted
   reference. Run the timing block without overlapping other heavy work;
   earlier task/checkpoint measurements cannot substitute.
4. Execute and map every deterministic QA-01..84 item and all 77 QA-owned VAL
   rows against final source. The full ordinary/tagged gates supply behavioral
   execution, with focused demonstrations where a row additionally needs an
   emptiness guard, filesystem comparison, provenance or other explicit proof.
5. Repeat the three external source-currency commands below; complete the
   final as-built task-18 corpus comparison against every C-01..C-13 item and
   DEC-187 disposition. Source-currency preparation is not final VAL-080 proof.
6. Complete security, requirement and documentation reviews; verify the actual
   retained Human fixtures and finalize neutral checklist locations. Record
   all 80 VAL rows in the final artifact with evidence and actual ownership.
   Human rows must remain explicitly unjudged until the Human records answers;
   deterministic completion does not claim phase sign-off or release readiness.

Required command results must be observed in the recording session. Keep each
long-running nested result/session until actual termination; coverage.out must
come from that session's successful coverage gate. Report-control receipts do
not replace product commands or the complete repeated-command reconciliation.

### Security and documentation review prerequisites

| Final check | Required evidence source |
|---|---|
| Archive expansion and entry containment | Task13's real measured 056 matrix; unsafe entries, duplicate keys/IDs and unknown Event/Effect/schema rejection; complete valid controls and unchanged source/destination snapshots |
| Atomic publication and target preservation | Task14 fault-point/state snapshots, populated-.fid rejection, preserved preexisting source files, same-key response-loss retry and changed digest/key rejection |
| Imported authority and old executor fencing | Task14/15 DEC-202 positive/negative journeys plus retained recovery force/unknown/replay/Stop/Budget cases; no import-created domain Event or automatic spawn |
| Force authority, unknown liveness, path scope | Existing task4–9 matrices rerun on final source; exact Human grant, scoped adoption, old fid-write fencing and no side effects on rejection |
| Contract/docs conformance | contracts.md, public-interface.md, schema.md, canonical 03_TRD_SCHEMA.md §8.3 and AGENTS.md command list against final CLI/catalog/migrations/discovery and task18 manifest |
| Removal/provenance | Task16's complete classified removal inventory and introducing commits; historical fixtures retain their original subject while latest-schema assertions prove absence |
| Capability honesty and exclusions | Exactly activation/Attention/export become available; Claims remains false with Phase24 ownership; no Review/Human/telemetry or full kit-production qualification inferred |

These are QA-79/80/81 preparation inputs, not completed security or documentation
findings. The quality plan requires written, evidence-linked review; this is not
an external-release phase and no new SAST/dependency/release gate is introduced.
Held-out scenarios remain outside this task and were not accessed.

### Future artifact inventory

The read-only inventory observed all task13–18 reports, review-5.md,
review-6.md, human-fixtures.md and final-review.md absent at this preparation
point. Each has an explicit future owner: Developers 13–18, the two window
reviewers, task18 fixture author, and final QA Lead respectively. The current
human-checklist.md exists with six neutral blank judgments. The named
p22_human_fixtures_test.go is not yet present and is allocated to task18.
These are scheduled deliveries, not present-readiness evidence or unowned
missing-artifact blockers.

The remaining-command/artifact Python inventory terminated exit_code 0 with:
`FINAL QUEUE INVENTORY COMPLETE: future artifacts classified by producing tasks; no product execution`.
The benchmark definition exists at internal/journal/rebuild_bench_test.go:19;
its 100000-Event path retains the 30-second self-assertion. `go help testflag`
terminated exit_code 0 and documents -bench, -benchtime 1x, -count, -run,
-timeout and explicit cache disabling. This validates syntax/reach only;
timing and retained-fixture usability remain unexecuted.

## External source-currency drift preflight — 2026-09-08

Own read-only commands, each with actual exit_code 0:

| Exact command | Observed output |
|---|---|
| `git -C /home/mark/projects/agentic-coding-project-starter-kit rev-parse HEAD` | `8fdd7d22c80ec508646ce53d0fae112f5b9216a2` |
| `git -C /home/mark/projects/agentic-coding-project-starter-kit status --porcelain` | `exit 0 (no output)` |
| `git -C /home/mark/projects/agentic-coding-project-starter-kit rev-list --count c5229617220f97ff9228df261e24718b1d1d6d89..HEAD` | `10` |

The kickoff's current-at-pin claim is stale. The accepted task18 comparison
pin remains c5229617220f97ff9228df261e24718b1d1d6d89; this review does not
silently replace it. Per consistency-review.md §5, the final comparison must
be re-derived and source currency rerun.

Only named normal specification sources were inspected. The exact read-only
diff command was:

```text
git -C /home/mark/projects/agentic-coding-project-starter-kit diff c5229617220f97ff9228df261e24718b1d1d6d89 HEAD -- project-docs/specs/trd/02_INTERFACES.md project-docs/specs/FID_COMPATIBILITY.md project-docs/specs/decisions/0005-linked-run-interruption-recovery.md
```

Observed exit_code 0. The port changes only add §7.1 offline ordinary-frame
binding; that section explicitly changes no public CLI/schema/backend-port
version and grants no production qualification. The current compatibility
matrix describes receiving/reconciliation obligations, an unqualified
production state, the discovery negotiation gap and offline decoder progress.
It retains DEC-187/CF-037 obligations rather than resolving them by fallback.

A separate Python comparison loaded both pinned Git blobs with git show.
The interval from the §3 Backend Port Operations heading through the end of
§4 Error And Stop Contract is byte-identical; Decision 0005 is byte-identical.
Observed output (exit_code 0):

```text
kit normative port sections 3/4 byte-identical: sha256=affa51f5a4658953b645a34945cb857500b419479601719342b7b81c324f4687
kit Decision 0005 byte-identical: sha256=e5725b67efeec17f87bd71cf46573560abfb4e8e1ff9f574fb3182cd6ab55820
SOURCE DRIFT CLASSIFIED: named recovery operations/errors unchanged; compatibility/receiving-status update requires final as-built review, not silent pin replacement
```

Thus the C-01..C-13 normative comparison inputs remain unchanged for liveness,
exact force grants, scoped/deferred budget units, Stop mappings, identity,
activation directions/export binding, validate-only/import, derived terminal
truth, Attention identity/disposition, repository observation and telemetry
ownership. DEC-187 and kit-contract-mapping.md still disposition every row.
The new compatibility status does not qualify any fid row and does not assign
new Phase22 product behavior. No new unresolved Phase22 contract difference was
identified in this preflight. Final VAL-080 remains pending comparison against
the completed task18 manifest and a fresh currency check, with any subsequently
changed normative requirement returned to the owning authority.
