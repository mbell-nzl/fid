# Phase 19 Final Review - Subsequent Disposition

Reviewer: QA Lead (Claude Opus 5)
Date: 2026-08-07
Operating mode: Heavy
Supersedes the verdict of: `final-review.md` (2026-08-02, BLOCK)

## Verdict

**PASS WITH IMPROVEMENTS**

Four of the five blocking findings from `final-review.md` are **closed on
verified evidence**. The fifth, **QA-BLK-FINAL-03, is sustained in part**: its
canonical-spec half is closed and verified against shipped code, but
`phase_brief.md`'s task table still records Task 4 as needing a closure session
and Task 6 as `Pending`, contradicting the same file's own Status section, and
`02_TRD.md` §3.1's lead sentence remains false for ten packages. Both are
single-edit documentary defects that misdescribe **no Phase 19 deliverable**;
they are named below as **preconditions on the wrap commit**, not as
carry-forward, so heavy mode's no-carry-forward rule is respected.

With the Human's VAL-026 result recorded, the validation contract is
**70 PASS / 0 FAIL / 0 BLOCKED / 0 CARRY-FORWARD**.

Two of the dispositions rest on evidence that does not prove what it is cited
for. In both cases I reached the **same outcome the Tech Lead reached** by a
different and stronger route, and I record the correction rather than the
original sentence, because the sentences have already been copied into
default-loaded files.

## Scope Of This Review

Per the Human's instruction this disposition is **scoped to the five findings**
in `final-review.md`. The 69 deterministic assertions, the eight phase-end gate
commands, the six `report-lint` runs, and the mid-phase review dispositions are
owned by `review-1.md`, `review-2.md`, and `final-review.md` and are **not
re-run here**.

That inheritance is only legitimate if the tree they were run against has not
moved. It has not:

| Evidence | Command | Result |
|---|---|---|
| E-01 | `git diff --name-only cfba808..HEAD -- cmd internal go.mod go.sum` | **no output** — zero product, test, or dependency change since the final review commit |
| E-02 | `git show --stat bb7d61a` | the only post-review commit; touches `project-docs/` specs and records exclusively |
| E-03 | `git status --porcelain` | clean except one untracked `project-docs/docs/reviews/*.html` (outside phase scope) |

Every PASS in `final-review.md` therefore still stands on its recorded evidence.

Documents loaded: `AGENTS.md`; `phase_brief.md`, `decisions_active.md`,
`progress_current.md`, `tl_escalations_current.md`; `08_QUALITY_PLAN.md`
§STRUCT-001/002; `qa/validation-contract.md` VAL-026;
`qa/final-review.md`; `qa/human-checklist.md`;
`qa/human-checklist-run-sheet.md`; `phase-19/decisions.md` DEC-105 (and
DEC-076/089/096/102/103/104 as cited); `02_TRD.md` §3/§3.1;
`03_TRD_SCHEMA.md` §14.11; `06_DESIGN_PLAN_CLI.md` `fid workflow list`;
`10_IMPLEMENTATION_ROADMAP.md` Phase 19.

No held-out scenario repository, evaluator prompt, or excluded scenario
material was read.

## Finding Dispositions

### QA-BLK-FINAL-01 (VAL-026) — **CLOSED**

`human-checklist.md` records HUMAN-26a..26i as **9 PASS / 0 FAIL / 0 PARTIAL /
0 blank**, overall `PASS — VAL-026 is PASS`, dated 2026-08-07, tester MBEL. The
file contains only `(human)` rows, nine of them, one judgment each, every row
mapped to VAL-026. VAL-026's pass criterion is "the Human records PASS for every
fanned-out sub-check"; its fail condition is "any sub-check recorded FAIL with a
described defect." Nine PASS, no FAIL. **VAL-026 = PASS.**

I did not substitute judgment for the Human's. I did corroborate that the
recorded evidence came from the shipped binary, because three rows quote exact
error text:

| Evidence | Check | Result |
|---|---|---|
| E-04 | Rebuilt CLI, fresh `mktemp -d` workspace, `fid workflow define --name "Invalid Workflow" …` | `error: invalid workflow slug "Invalid Workflow"`, exit 1 — **verbatim match** to HUMAN-26g |
| E-05 | Same workspace, `fid workflow assign --name task-default --by tl` | `error: workflow assignment requires a kind or project scope`, exit 1 — **verbatim match** to HUMAN-26h |
| E-06 | `grep` for HUMAN-26i's quoted string | `internal/policy/authorize.go:27` — `authorize action %q: no permissions matrix entry`, verbatim; both states already proven deterministically by VAL-024 (F-08) |

The run sheet (`human-checklist-run-sheet.md`) is a concrete, executable command
sheet against a disposable `mktemp -d` workspace and instructs the tester not to
read source before recording 26a/26b. The scratch workspace was removed after
this check.

**CF-026 is a finding, not unshipped scope — confirmed.** HUMAN-26c passed and
noted that full named-workflow graph inspection is absent.
`06_DESIGN_PLAN_CLI.md:2238-2258` documents the `fid workflow list` payload as
`id`, `name`, `applies_to`, `active`, `states` **count**, `transitions`
**count**, `assigned_to` — counts, not a graph. The shipped surface matches the
documented surface exactly (VAL-018/060..067). No spec, phase brief, roadmap exit
clause, or contract row promised `workflow show`, workflow revision,
deactivation, or slug reuse. CF-026a (read surface) and CF-026b (lifecycle
model, design decision first) are therefore **new product gaps discovered by
dogfooding**, which is what a dogfood row is for. The split into a schedulable
half and a decide-first half is the right shape.

### QA-BLK-FINAL-02 (architecture contradiction) — **CLOSED**

The disposition resolves this against the decision rather than the code, and
that is the correct direction. `02_TRD.md` §3.1 (lines 120-149) now states that
`internal/workflow` → `internal/journal` **includes the journal write API**,
names the six symbols, cites DEC-089 as the authority, records DEC-096's
`Effect`-only phrase as a Task-2 snapshot wrongly stated as permanent, and
preserves DEC-096's operative rule (no in-package `package db` or
`package journal` test file may import `internal/workflow`).

Verified symbol-for-symbol against shipped code:

| Evidence | Symbol named in §3.1 | Shipped call sites |
|---|---|---|
| E-07 | `WithLock` | `define.go:57`, `assign.go:49` |
| E-07 | `EnsureProjectionFreshWithLock` | `define.go:75`, `assign.go:67` |
| E-07 | `PrepareEvent` | `define.go:87`, `assign.go:82` |
| E-07 | `DryRunProjectEvent` | `define.go:101`, `assign.go:96` |
| E-07 | `AppendEvent` | `define.go:104`, `assign.go:99` |
| E-07 | `ProjectEvent` | `define.go:133` |

`02_TRD.md:96` (§3 layer table) independently matches: the `internal/workflow/`
row now reads "workflow-assignment validation and **write paths**, and
assignment resolution consumed by `internal/issues`", and its must-not column is
exactly the boundary DEC-105 says survives — issue creation, AC verification,
output formatting, schema/migration ownership. `internal/workflow` creates no
issues and owns no migrations. **The contradiction I raised is gone.**

One residual, raised as a **new** finding rather than a sustain, because it is
about packages outside Phase 19's deliverable: see MAJOR-01.

### QA-BLK-FINAL-03 (stale canonical records) — **SUSTAINED IN PART**

Four of the five named records are updated and verified against shipped
behaviour. One is not.

| Record | Verdict | Evidence |
|---|---|---|
| `02_TRD.md` §3/§3.1 | **CLOSED** | E-07 above; §3 row at line 96 |
| `03_TRD_SCHEMA.md` §14.11 | **CLOSED** | E-08 below |
| `10_IMPLEMENTATION_ROADMAP.md` Phase 19 closure | **CLOSED** | E-09 below |
| `progress_current.md` | **CLOSED** | reflects six tasks complete, VAL-026 9 PASS, CF register current |
| `phase_brief.md` | **OPEN** | E-10 below |

**E-08 — `03_TRD_SCHEMA.md` §14.11 (line 1228) matches the code.** The section
now carries the workflow-assignment compatibility bullet and a "Table scope
(corrected v0.2a, DEC-105)" paragraph replacing the v0.1-only exclusion with
"a projected table is audited once it has a structural invariant that a corrupt
or externally modified projection could violate." Checked against
`internal/maintenance/workflow_assignments.go`:

- kind-scope row flagged unless `applies_to` is NULL (`any`) or equals the kind
  (`workflow_assignments.go:42`) — matches "a kind-scope row whose workflow is
  bound to a different kind";
- project-scope row flagged unless `applies_to` is NULL (`workflow_assignments.go:45`)
  — matches "a project-scope row whose workflow is not `any`";
- the finding message names scope, entity, workflow, and the offending
  `applies_to` (`workflow_assignments.go:57-61`) — matches "naming the offending
  assignment";
- registered as a structural check (`check.go:67`) with deterministic sort
  (`check.go:93`) — matches "deterministic by check, entity ID, then message";
- read-only: the file contains one `db.Query` and no write — matches
  "never repairs projection rows or appends a journal event".

The claim that FK integrity is covered by the `PRAGMA foreign_key_check` pass
and deliberately not duplicated is accurate and is the right call.

**E-09 — roadmap closure disposition is present**
(`10_IMPLEMENTATION_ROADMAP.md:542`) and accurate: exit criteria met, six tasks,
70-row contract, VAL-026 9 PASS, zero product defects, the three phase-level
lessons, CF-026a/CF-026b/CF-027, and the deep-retrospective recommendation.

**E-10 — `phase_brief.md` is internally contradictory and its task table is
stale.** The file's Status section (lines 25-30) says "All six tasks complete"
and lists which records are synchronized — notably **not itself**. Its task
table 55 lines later says otherwise:

- `phase_brief.md:86` — Task 4: "**Needs a closure evidence session** (DEC-103)
  to produce a `report-lint`-clean report with its own tokens." That session
  happened: `ba65078` is committed and `final-review.md` records
  `PASS report-lint completed with no failures` against it.
- `phase_brief.md:88` — Task 6: "Pending — **CF-025 first and mandatory**." Task
  6 shipped as `09aa0a0` with report `988737b`; `progress_current.md:33` and
  `decisions_active.md` both record CF-025 as **closed in Task 6**.

This is the default-loaded execution brief. A Phase 20 planning session reads it
before anything else and would read that Phase 19's closure task never ran.
I am not softening this because the wrap will replace the file: the wrap has not
happened, and DEC-105 item 3 scheduled the edit rather than performing it while
`bb7d61a` performed the other four.

**Severity judgment.** I am **not** re-blocking the phase on it, and the
distinction is principled rather than convenient. The original QA-BLK-FINAL-03
was blocking because **permanent canonical specs actively contradicted shipped
code** — a reader auditing `internal/workflow` or `internal/maintenance` against
the TRD would have found a false violation. That half is closed and verified.
What remains misdescribes **no shipped behaviour and no Phase 19 deliverable**;
it is a status field in a transient current-context file. See PRECON-01.

### QA-BLK-FINAL-04 (missing `Refs:` on `32614bb`) — **CLOSED; evidence corrected**

**The outcome is right and I sustain nothing.** Immutable history plus a widely
cited hash makes rewriting a strictly negative trade, and the disposition's
remedy — record the deviation, move enforcement to CF-023 as a `commit-msg` hook
or `report-lint.sh` check — is the correct one. Exhortation has already failed
here; mechanism is the answer.

**The cited evidence does not establish the cited conclusion.** DEC-105 item 4
and the CF-023 register line in `decisions_active.md:42` both say
`git log --format='%h %s%n%b' 50533e2..HEAD | grep -c 'Refs:'` → 10, "so the
convention holds and `32614bb` is an outlier." I reproduce the 10. But the range
holds **18 commits** (`git log --oneline 50533e2..HEAD | wc -l`), so **eight**
lack the trailer, not one. 10-of-18 does not demonstrate that a convention
holds.

**The correct statement is narrower and stronger, and it supports the same
disposition.** Partitioning the eight (E-11):

| Commit | Subject | Class |
|---|---|---|
| `32614bb` | `feat(workflow): add assign command and assignment-aware creation` | **Task commit — Task 4 irregular session** |
| `ce4cf5b` | `docs(task): finalize completion report for phase-19 task-4` | **Task commit — Task 4 irregular session** |
| `f7d7c04`, `3a6ca15`, `2711f09`, `4945f34`, `bb7d61a` | `docs(tl): …` | Tech Lead governance commits |
| `ebbcd16` | `Phase 10 task 4 unblock` | Raw human commit; see MINOR-02 |

Among the range's Developer and QA-Lead **task** commits, **every one carries
`Refs:` except two — and both are products of the same irregular human-executed
Task 4 session.** That is a far better-founded version of the disposition's own
argument, and it raises that session's defect count from four to five.

Two corrections follow that should land with the wrap edits:

1. `progress_current.md` states "`32614bb`, the **single commit** produced by the
   irregular human-executed Task 4 session … **No other Phase 19 commit produced
   any**." Both clauses are false: `ce4cf5b` is from the same session and carries
   the same defect.
2. `decisions_active.md:42`'s CF-023 line repeats "`32614bb` missed it while 10
   sibling commits carried it, so the convention holds." Replace with the
   partition above, and scope the mechanical check to **all** phase commits —
   `AGENTS.md:64` requires `Refs: phase-N/task-M` for "phase work" without
   exempting `docs(tl)`, so either the hook covers them or the rule should say
   so explicitly. CF-023 is the right home for that call.

The "strictly negative trade" argument is also **understated in its own
favour**: DEC-105 says the hash is cited in eight places; it is actually cited
**30 times across 12 files** (E-12), including `task-4-done.md`, `task-4.md`,
`task-5-done.md`, `review-2.md`, `run-ledger.md`, `triage-log.md`, four
`decisions`/`current` files, and `final-review.md`. No rewrite. Confirmed.

### QA-BLK-FINAL-05 (structure budgets) — **CLOSED as inverted; scope corrected**

The provenance inversion is **correct and reproducible** (E-13):

| Selector | `internal/issues/create.go` |
|---|---|
| `git show 196a31d:internal/issues/create.go \| wc -l` (phase start) | **419** |
| `git show 50533e2:internal/issues/create.go \| wc -l` (Task 3) | **419** |
| `wc -l internal/issues/create.go` (today) | **412** |

The STRUCT-002 breach predates the phase and Phase 19 made the file **seven
lines smaller**. Applying DEC-102's rule — "pre-existing" is a provenance claim
requiring provenance evidence against the phase-start commit — is exactly what
my own finding failed to do. **My QA-29 classification was wrong on provenance;
the measurement was right and the classification was not.** Deferring a 228-line
decomposition of `internal/issues.Create` — the most exercised write path in the
product — to land **behind** CF-025's byte-level parity test is the correct
sequencing, and CF-027 records it with that sequencing attached.

I also confirmed the disposition's supporting claim that `create.go` is the only
task-touched file over budget. It is, for STRUCT-002 (E-14): scanning every
non-test `.go` file changed in `196a31d..HEAD` under `cmd/` and `internal/`,
`internal/issues/create.go` at 412 is the **only** file above 400. Next largest
is `cmd/fid/workflow.go` at 261.

**One correction to my own finding, in the other direction.** QA-29 checked
STRUCT-002 file length and only `Create` for STRUCT-001. Scanning function
lengths across all phase-touched files surfaces a second STRUCT-001 breach that
neither my final review nor DEC-105 named: `internal/format/markdown.go`
`writeMarkdown` is **68 lines** against the 60-line target. Provenance, checked
by the DEC-102 rule before I classify it: it was **66 lines at phase start**
(`196a31d`) and Task 5 (`8f79013`) added **two** lines — a single
`case []workflow.ListItem:` arm to a type-switch dispatch. That is the minimal
legal way to register the new renderer; refusing it would have meant not shipping
`workflow list --format markdown`. **Inherited breach, marginally grown by an
unavoidable dispatch arm.** Non-blocking. It belongs on the CF-027 register line
so the next `internal/format` reopen (already scheduled by CF-022's deferred
renderer backfill) picks it up in the same pass. See MINOR-01.

## Updated Outcome Ledger

### Validation contract

`final-review.md` recorded 69 PASS and VAL-026 BLOCKED. With the Human's result:

| VAL ID | Outcome | Validator | Evidence |
|---|---|---|---|
| VAL-001..025, VAL-027..070 | **PASS** (69) | QA Lead | Unchanged; inherited from `final-review.md` on E-01 (no product change since `cfba808`) |
| VAL-026 | **PASS** | Human | `human-checklist.md` HUMAN-26a..26i, 9 PASS / 0 FAIL / 0 PARTIAL, 2026-08-07; corroborated by E-04/E-05/E-06 |

**Assertion summary: 70 PASS, 0 FAIL, 0 BLOCKED, 0 CARRY-FORWARD.**

### Deterministic checklist

Only the three FAILs from `final-review.md` are re-dispositioned; QA-01..QA-28
and QA-32..QA-75 are unchanged at PASS.

| QA ID | Result | Disposition |
|---|---|---|
| QA-29 | **FAIL (accepted → CF-027)** | Measurement unchanged (412 lines; `Create` 19-246). Reclassified on provenance: inherited breach, phase-start 419 → 412. Not a Phase 19 regression. Sequenced behind CF-025. Register line should also absorb `writeMarkdown` (MINOR-01). |
| QA-30 | **FAIL (narrowed)** | Canonical-spec half closed and verified (E-07/E-08/E-09). Open: `phase_brief.md:86,88` and `02_TRD.md:121-123`. See PRECON-01 and MAJOR-01. |
| QA-31 | **FAIL (accepted deviation → CF-023)** | Observable unchanged: `32614bb` carries no `Refs:`. No rewrite; enforcement moved to a mechanical check. Evidence and defect population corrected (E-11/E-12). |

**Checklist summary: 72 PASS, 3 FAIL** — one narrowed to two documentary edits,
two accepted onto the carry-forward register with named targets and sequencing.

Register items with named targets are the normal mechanism and are not
carry-forward of unfinished phase scope; heavy mode's prohibition is not
engaged by CF-026a/b or CF-027.

## Findings

### Preconditions On The Wrap Commit

- **PRECON-01 — `phase_brief.md` task table is stale and self-contradictory.**
  Reproduction: read `phase_brief.md:25-30` ("All six tasks complete"), then
  `phase_brief.md:86` and `:88`. **Closes with:** Task 4's row reading
  `**Complete** (32614bb, closure report ba65078)`, and Task 6's row reading
  `**Complete** (09aa0a0, report 988737b) — CF-025 closed`. Two cells. This is
  the last open piece of QA-BLK-FINAL-03 and must land **in** the wrap commit,
  not after it.
- **PRECON-02 — two false statements about the Task 4 session, already copied
  into default-loaded files.** Reproduction: `progress_current.md` "the single
  commit produced by the irregular human-executed Task 4 session … No other
  Phase 19 commit produced any", and `decisions_active.md:42` "10 sibling
  commits carried it, so the convention holds". **Closes with:** the E-11
  partition — every Developer/QA task commit in `50533e2..HEAD` carries `Refs:`
  except `32614bb` and `ce4cf5b`, both from the Task 4 session; the other six
  gaps are five `docs(tl)` commits and one raw human commit.

### Major

- **MAJOR-01 — `02_TRD.md` §3.1's lead sentence is false for ten packages and
  contradicts its own next bullet.** Reproduction: `02_TRD.md:121-123` states
  "Domain packages that produce journal events depend on `internal/journal` for
  the `Effect` data contract **only**"; the bullet immediately below exempts
  `internal/workflow`. But
  `grep -rl 'journal\.\(WithLock\|AppendEvent\|ProjectEvent\)' internal --include=*.go`
  (non-test) returns **ten** packages — `bootstrap`, `claims`, `dev`, `issues`,
  `labels`, `links`, `maintenance`, `requirements`, `runs`, `workflow` — e.g.
  `internal/issues/ac.go:61,87,133,150`. DEC-089 clause 2 already said as much
  ("`internal/issues` and `internal/dev` already orchestrate
  `journal.WithLock` → `PrepareEvent` → `AppendEvent` → `ProjectEvent`"). The
  sentence is **Phase 19-introduced** (`65f6710`, the DEC-096 sync; §3.1 did not
  exist at `196a31d`) and `bb7d61a` corrected only the workflow bullet beneath
  it. An auditor checking `internal/issues` against §3.1 today finds a false
  violation — the same defect class as the original QA-BLK-FINAL-02, aimed at
  packages outside this phase's deliverable. **Closes with:** restating the lead
  sentence as the actual norm — domain packages depend on `internal/journal` for
  the event/effect data contract **and** the write API — and keeping the
  one-way `workflow → journal → db` chain, which is the constraint the section
  exists to protect. While there, §3.1 enumerates the write API for `define.go`
  and `assign.go` only; `internal/workflow/list.go:47` also calls
  `journal.EnsureProjectionFresh` on the read path, which the enumeration should
  mention or the prose should generalise.

### Minor

- **MINOR-01 — second STRUCT-001 breach not named by either document.**
  `internal/format/markdown.go` `writeMarkdown` is 68 lines (target 60). Was 66
  at `196a31d`; `8f79013` added two lines as one type-switch dispatch arm.
  Inherited, unavoidable growth, zero product risk. **Closes with:** adding it to
  the CF-027 register line so the next `internal/format` reopen — already
  scheduled by CF-022's deferred renderer backfill — takes both in one pass.
- **MINOR-02 — one commit in the phase range is outside the commit convention
  entirely.** `ebbcd16`, authored by the Human, subject `Phase 10 task 4
  unblock`: no Conventional Commit type, no `Refs:`, and the **wrong phase
  number** — it touches only `phase-19/` and `current/` files. Same
  immutable-history logic applies; no rewrite. **Closes with:** noting it under
  CF-023, where a `commit-msg` hook would catch subject format and trailer in
  the same check.

### Improvement

- Carried unchanged from `final-review.md`: the 100k rebuild control has ~20%
  headroom (`23.935 s` against `<30 s`) and is worth watching as v0.2b adds
  receipt volume. The target must remain unchanged.

## Human Checklist Readiness

- Contains only `(human)` items: **YES**.
- Contains exactly HUMAN-26a..HUMAN-26i: **YES, nine rows**.
- One product/perception judgment per row: **YES**.
- Human-only VAL-026 mapped on every row: **YES**.
- Final status: **9 PASS, 0 FAIL, 0 PARTIAL, 0 blank; VAL-026 PASS**.
- Product risks called out at generation — discoverability, four-format
  coherence, corrective error quality, DEC-079 understandability — were all
  exercised and all passed. The dogfood produced one genuine product gap
  (CF-026a/b) that no deterministic assertion could have found, which is the
  strongest available argument for keeping the Human row in the contract.

## Coverage Caveat

VAL-026's nine PASS outcomes close the operator-coherence gap the deterministic
suite could not reach. They do not extend to workflow inspection depth or
lifecycle management, which the dogfood explicitly found missing and which
CF-026a/CF-026b now carry. Nothing in this review re-derives the 69
deterministic outcomes; they rest on `final-review.md`'s evidence and on E-01's
proof that the tree has not moved.

## Release Maturity Status

| Maturity area | Status | Evidence or deferral |
|---|---|---|
| Baseline artifacts | **READY** | 70/70 assertions have outcomes; all five blocking findings dispositioned; two documentary preconditions named against the wrap commit |
| Automated gates | **READY** | Eight-command phase gate 8/8 PASS (`final-review.md`); tree unchanged since (E-01) |
| Held-out scenarios | DEFERRED | No accepted harness exists; none was read. The Human dogfood is not represented as held-out evidence. |
| Release readiness | DEFERRED | Phase 19 is not an external-release phase; security/dependency/release commands stay governed by the roadmap and quality plan |

## Note On The Two Corrected Evidence Claims

Both corrections point the same way, and it is worth naming the pattern for the
retrospective DEC-104 recommends. QA-BLK-FINAL-04's disposition asserted a
conclusion from a bare `grep -c` whose denominator was never taken;
QA-BLK-FINAL-05's disposition was right precisely **because** it took the
denominator (phase-start line count) that my finding had skipped. The same
phase produced both. DEC-102's provenance rule is the thing that separates them,
and it now has a third case: **a count is not a rate until the population is
stated.** That is one sentence and it is cheaply mechanical — the same class as
the four drill cases already accumulated on CF-023, and it belongs there.

Recorded plainly and for the record: the phase-19 dispositions were argued
honestly and reached the right outcomes on all five findings. Two of them were
reached on evidence weaker than the argument they were used to support, and one
of my own three deterministic FAILs was misclassified because I skipped the
provenance check my own quality plan requires. Neither of those changes any
outcome; both change what the record should say.

---

**VERDICT: PASS WITH IMPROVEMENTS** — all five blocking findings from
`final-review.md` are dispositioned; VAL-026 is PASS and the contract closes
70 PASS / 0 FAIL / 0 BLOCKED / 0 CARRY-FORWARD; QA-BLK-FINAL-03 is sustained in
part and closes with PRECON-01 and PRECON-02 landing **in** the wrap commit;
MAJOR-01 closes with one corrected sentence in `02_TRD.md` §3.1; MINOR-01 and
MINOR-02 are register additions to CF-027 and CF-023. No product code change is
required and none is recommended.
