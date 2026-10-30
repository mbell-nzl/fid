
"""Assemble the initial final review and Human-only checklist from captured
Phase 24 evidence. Run only after all gate/probe receipts have terminated.
"""
import pathlib,json,re,shlex
ROOT=pathlib.Path(__file__).resolve().parents[6]
BASE=ROOT/'project-docs/docs/phases/phase-24'
OUT=BASE/'qa/final-evidence'
def read(name): return json.loads((OUT/name).read_text())
def code(s): return chr(96)+s.replace('|','\\|')+chr(96)
def link(name): return '['+name+'](final-evidence/'+name+')'
gates=[read('G%02d.json'%i) for i in range(1,15)]
for g in gates: assert g['chunks'][-1].get('exit_code')==0,g['id']
for i in [1,2]: assert read('supplemental-%d.json'%i)['exit_code']==0
probes=read('binary-probes.json'); docs=read('document-audit.json'); inv=read('inventory-and-help.json')
for r in inv:
    # Help usage observations are classified separately, not passing product gates.
    if r.get('registry') or r['command'].startswith('rg -n "PRIVATE_') or r['command'].startswith('rg -n "This '):
        assert r['exit_code']==r['expected_exit_code'],r
assert probes['product_identity']['exit_code']==0 and not probes['product_identity']['stdout']
assert all(r['exit_code']==0 for r in read('report-lint.json'))
head=docs['head']['stdout'].strip()
s1=link('supplemental-1.log'); s2=link('supplemental-2.log')
fixture=read('human-fixture-manifest.json')
validators={}
for line in (BASE/'qa/validation-contract.md').read_text().splitlines():
    if re.match(r'\| VAL-\d{3} \|',line):
        cells=[c.strip() for c in re.split(r'(?<!\\)\|',line)[1:-1]]
        validators[cells[0]]=cells[4]
assert len(validators)==68
def output(row): return '[Human output](final-evidence/human-output/'+row['name']+'.txt)'
def gate_output(g): return ''.join(c.get('output','') for c in g['chunks'])
bench=[re.search(r'([0-9.]+) ms/rebuild',gate_output(g)).group(1) for g in gates[10:13]]
assert all(float(x)<30000 for x in bench)
qa_notes={
1:'G01-G07: every exact phase command terminates 0; see gate table.',
2:'G08 full checkptr; G09 full archive-scale alternation; both terminate 0.',
3:'G10 reference instance; G11-G13 fresh processes: '+', '.join(bench)+' ms/rebuild; parity asserted by benchmark.',
4:link('schema-proof.json')+'; schema-18/schema-19/down dumps; cmp exit 0; '+link('binary-probes.json')+' fresh CLI init has all 102 identical projection objects plus the explicitly asserted migrator metadata index; '+s1+' TestMigration0019AddsOnlyFieldDigestColumn also checks every migrator object on up/down.',
5:link('binary-probes.json')+' live/rebuilt map SQL byte equality; '+s2+' TestP24All49KeyedLeavesRebuildInOneJournal and portable round-trip.',
6:s1+' TestCurrentIdempotencyFieldDigestsReplayFailClosed and TestGateReceiptFieldDigestGuardAndMalformedInsertFailClosedThroughCLI: valid receipt, missing/array/non-digest inserts, bad guard, exact stop line/exit 6 and unchanged tree.',
7:s1+' TestSchema19ArchiveImportsAndSchema18ManifestFailsClosed: 19 success; 18 archive-invalid/schema-unsupported manifest.json; unchanged target.',
8:s2+' all per-family keyed corpora and TestP24ExecutedKeyedCorpusMatches49LeafInventory; '+s1+' live keyed-tree checks; original EventRef, conflict sorted fields/exit 4, canonical replay/no append.',
9:s2+' TestP24ReviewOptionalPresenceAndEmptyException and Source-presence conflict rows: all three states.',
10:s2+' TestP24ObservingJudgmentRetriesAfterSourceAndHEADChange, Run replay rows and gate/binding replay corpora; unchanged Event count.',
11:link('inventory-no-short.txt')+' empty; 95 live leaves; '+s1+' TestEveryLiveLeafHasOneLineShort.',
12:s1+' TestInputDescriptorRegistryMatchesLiveTreeAndTypes: 39 descriptors, no fence, Go reflection and seeded mismatch.',
13:link('binary-probes.json')+' captures all 39 help/templates; '+s1+' per-leaf descriptor process corpus, full help pins, types/order/indent and Git meanings.',
14:s1+' TestAdoptedExampleProcessExclusivity and TestExampleRejectsOtherInputsInFixedOrder; all adopted leaves, including binary JSON/project conflicts.',
15:s1+' structural process and variant corpora; nested unknown/enum/null/type/required/duplicate/format errors; sensitive value null/no Event.',
16:s1+' TestP24RemovedBindingMembers, TestP24RemovedRunRepositoryMembers and TestP24RemovedGateEvidenceMembers: all 17 leaves.',
17:s1+' source resolver/descriptor Source-position matrices; '+s2+' judgment/review/envelope/budget Source binary rows cover match, mismatch, omitted-existing/missing and all path/external states.',
18:s1+' task 4-8 observation tests; '+s2+' persisted records and Events versus independent Git/config/root oracle; eight repository members, including non-null dirty digest.',
19:s1+' TestObserveRepositoryRejectsUnsafeImages/RejectsChangeDuringCapture and P24 observing classifier tests; '+s2+' normal, unmerged and non-Git rows per applicable leaf.',
20:s1+' review/Run/gate expectation matrices; '+s2+' terminal expectation/unsafe process corpus; five leaves x three expectation states, exact two-field expected object.',
21:s1+' gate binding digest, original/copy/import tests; '+s2+' TestP24PlainCopyFailurePairs; original success/import equivalence, four copy refusals and exact Claim-first Run cause.',
22:s1+' TestP24RunStartAfterCommitAndEdit and TestP24RunCloseAfterCommitAndEdit; '+s2+' terminal observation process rows.',
23:s1+' TestP24GatePreSpawnObservationAfterEdit/TestP24GateExpectationBeforeAttempt; eight-member independent oracle, Attempt/Receipt/two Events identical; no-spawn marker and zero new rows on mismatch.',
24:s2+' TestP24ReceiptFreshnessRealGuardedActions: unchanged/edited/reverted x evaluate/terminal transition/envelope complete.',
25:s1+' TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation; '+s2+' TestP24ReviewAndHumanEvidenceSubjectRepositories: three subject types and Event repository.',
26:s2+' changed fake/fid corpora for Source, review, binding, Run and gate; dispatch-derived completeness, acceptance/rejection parity.',
27:s2+' RegisteredCodePairs/Completeness, protected-decision and post-frame gate failure; '+s1+' HumanFailureLines/ReasonCut and Grammar matrices; next-table positive/negative rows and 200/201 boundaries.',
28:s1+' Stop corpus and completeness; '+s2+' TestP24StopHumanRealCorpus: all ten code/Condition rows, open/resolved; '+link('binary-probes.json')+' clean/dirty x same/different view captures and forbidden-token assertions.',
29:s2+' TestP24JudgmentGetRecords/Rejections: all three reads versus list, sole member, event null; unknown ID and grammar inside/outside project.',
30:s1+' TestP24JudgmentSectionDisclaimers/RecordedReasoningEveryFieldAndView: 0/1/3, all closed fields, 80/81 non-ASCII/multiline/control rows; '+s2+' detail/resume and every retained fixture sentinel guard.',
31:link('glossary-comparison.md')+' quotes GLOSSARY.md:21/32 and public-interface section 10, and compares all five recorded-field groups and referenced-file meaning.',
32:s2+' HumanTakeControlJourney, TakeControlTerminalStateProvenance, JourneyProvenanceKeepsExactSourceLine, TakeControlRebuildExportImportParity; verbose exact source/class logs, completion and schema 19. VAL-068: '+s1+' TestRegistryNextStates/TestP24ResumeIssuePosition/TestP24ResumeIssueNextStates/TestP24IssueTransitionPositionHelp and '+s2+' TestP24TakeControlResumePair prove revision 4 with sorted cancelled(cancelled), qa_block(started), qa_pass(completed), revision 5 with JSON []/Human none, nil/empty Human rows, exact key/line placement and both help pointers.',
33:s1+' descriptor forbidden-member scan; '+s2+' real Stop no-JSON assertions; '+link('inventory-and-help.json')+' disclaimer absence expected/observed exit 1.',
34:link('document-audit.json')+' every report timing section independently inspected, fresh versus cached kept distinct; no accepted final task sequence exceeds 600/900; warnings retained. G03/G05 meet phase budgets.',
35:'Security boundaries reviewed at envelope.go:279-312, run.go:165-194, evidence/execution.go:185-222, project/source.go:67-123, workflow/import_binding.go:18-51; replay/containment/sensitive/import regression output in S1/S2. See security findings below.',
36:link('report-lint.json')+' all 15 exact spec lists PASS/OUTCOME COMPLETE; '+link('document-audit.json')+' finds missing owned R4 and R10 receipts (F1/F2). Other owned commands present. '+link('provenance-commits.json')+' resolves all 168 introducing-change/phase-history references with observed status 0; tokens, termination declarations, expected non-zero, convergence and clean identity inspected.',
37:link('document-audit.json')+' dependency diff 304aa45 empty, exit 0; added files stay within existing cmd/fid and internal owners; TRD section 2.',
38:link('inventory-and-help.json')+' complete R1-R20 regeneration; labelled section-1 lists; deltas discussed below. R5 expected/observed exit 1; R19 four fallback rows.',
39:link('document-audit.json')+' task amendments in schema/design/CLI/flows/quality/glossary; accepted PRD IDs retained; wording requiring TL wrap-up consolidation listed below.',
40:'Reviews 2/4/6/8/10/12 close A/B/C/D1/D2/E respectively; independent closure of every original finding; Mark Bell walkthrough and original/clarification commits verified. See disposition table.',
41:s2+' judgment/binding/Run/gate replay-never-observes matrices: unobservable/deleted states, unchanged journal count and valid new-key controls; '+s1+' normal controls and owning independent oracle checks.',
42:link('binary-probes.json')+' attention.resolve/stop.raise/stop.resolve help captures; '+s1+' help/Short tests; no Run wording.',
43:link('binary-probes.json')+' all 39 --example exit 0, empty stderr, empty before/after tree; descriptor type/format assertion output in S1.',
44:s2+' TestP24RunTerminalExpectationObservationAndUnsafe: close/abandon mismatch repository-mismatch versus unmerged recovery-state-conflict; no append.',
45:s1+' Run obligation terminal observation fixtures: unchanged succeeded versus tracked edit failed; '+s2+' built-process close controls.',
46:s1+' Grammar whole-tree/class/domain/completeness/process suites; four grammar classes, every typed/cited string domain, both modes, no-project repetitions, DEC-282 short/long/512/513/non-ASCII/escape-cut boundaries.',
47:s1+' Grammar groups/built-binary/no-word controls; root bogus/isue, all 29 groups, operation/allowed/exit 2, exact paired messages and DEC-282 boundaries.',
48:link('inventory-and-help.json')+' exact PRIVATE_ absence: exit 1 (no output), expected and observed exit_code 1.',
49:s2+' GenericFallbackCorpusCompleteness/DomainPairs/SensitiveProcess; '+s1+' GenericFallback and BoundedMessage units: all 18 family sites, no exclusions, unchanged metadata, 512/513, blank/whitespace/escape-cut; copy Run claim cause.'
}
# Final review retains the template sections. Evidence in this session is
# distinguished from historical report statements and old window executions.
r=['# Phase 24 Final Review - Input authoring and Human take-control','',
'Reviewer: Codex QA Lead','Date: 2026-10-11','Review type: INITIAL','Prior review: none','',
'## Verdict','','BLOCK','',
'Phase advancement is blocked by the five outstanding Human-only assertions (B1-H). Two additional Must-fix findings are bounded report-only repairs in Tasks 1 and 10 (F1/F2). Their severity does not itself establish a protected blocker. All deterministic product assertions pass; the report obligations must still be repaired and independently re-reviewed.',
'','## Scope And Evidence','',
'Reviewed phase-24 at '+code(head)+' after committed Window E review-12. Loaded AGENTS.md, all four current context documents, the quality plan, workflow sections 2.6 and 8, the QA plan and all 68 contract assertions, all 15 task specs/reports, inventory R1-R20, reviews 1-12, the recorded Task 11 walkthrough, and relevant accepted Phase 24 decisions including DEC-291. Applied validation-contract and qa-evidence skills. Required templates govern these two artifacts.',
'',
'Current docs still schedule review-12; the committed review and task-report receipts govern the actual completed checkpoint. No product, test, task, decision, contract, QA plan, prior review, current document or Runner ledger was edited. No held-out/evaluator/excluded material was read.',
'',
'Evidence paths below are durable artifacts of this final-review session. Each G-file preserves every complete nested execution/poll result, including session IDs and observed terminal exit_code. S1/S2 JSON receipts retain the exact full child stdout/stderr, actual subprocess termination and command/cwd; their outer execution result is also preserved. The .log files are display copies with line-ending whitespace trimmed; '+link('log-display-provenance.json')+' binds their hashes to the immutable raw JSON streams. Report-lint is an independently rerun Runner control, not a replacement product gate.',
'','## Phase-Level Gate Results','',
'Commands ran serially in /home/mark/projects/fid. QA-plan section 5 explicitly selects the approved -timeout 30m forms for ordinary and tagged gates (DEC-227/259). These are the exact task/phase forms; no scope was narrowed.',
'','| Gate | Command | Result | Evidence |','|---|---|---|---|']
labels=['Format','Lint','Race + coverage','Coverage report','Tagged integration','Benchmarks','Build','Full-checkptr control','Archive-scale']
for label,g in zip(labels,gates[:9]):
    lines=gate_output(g).strip().splitlines()
    token=lines[-1] if lines else 'exit 0 (no output)'
    r.append('| '+label+' | '+code(g['cmd'])+' | PASS | '+link(g['id']+'.json')+'; '+code(token)+'; observed exit_code: 0 |')
r+=['| Full phase gate | G01 through G07, exactly in the declared order | PASS | All seven observed exit_code: 0 |',
'| Held-out / end-to-end | N/A | N/A | No harness accepted; phase brief/quality plan maturity deferral |',
'| Reference instance | Six commands from QA-plan section 5 | PASS | '+link('G10.json')+'; recorded in this session before all three samples |']
for g,value in zip(gates[10:13],bench):
    r.append('| Fresh 100k sample '+g['id']+' | '+code(g['cmd'])+' | PASS | '+link(g['id']+'.json')+'; '+value+' ms/rebuild (<30000); count parity; exit_code: 0 |')
r+=['| Human fixture preparation | '+code(gates[13]['cmd'])+' | PASS | '+link('G14.json')+'; 59-row manifest independently validated below; exit_code: 0 |','',
'G04’s tool output was truncated by the execution surface. Its complete nested result and terminal total are retained. The identical cover-report command was repeated by the inventory helper to retain the entire report in '+link('coverage-functions.txt')+' and '+link('coverage-report-complete.json')+'. Total coverage: 62.8%. No coverage threshold is invented. Cached G03 package results are passing captured command results, not fresh numeric timings. G05 reports 440.177 s, below 900 s. G08 reports cmd/fid at 712.964 s; the full-checkptr control is explicitly exempt from VAL-052 ordinary budgets under DEC-259 and terminates 0. Historical ordinary 480 s warnings remain in Tasks 10/12/13/14 reports.',
'','## Validation Contract Outcomes','',
'63 deterministic assertions: PASS. Five Human assertions: BLOCKED pending actual Human judgments. No CARRY-FORWARD or silent deferral. Each QA item below cites its current transcript/file proof; S1/S2 enumerate named tests and every subcase with PASS output, including positive and negative states.',
'','| VAL ID | Outcome | Validator | Evidence |','|---|---|---|---|']
mapping={}
for a,nums in {4:[1,58],5:[2],6:[3],7:[4],3:[5],8:[6,7],9:[8],10:[9],11:[10],12:[11],13:[12,13,14],14:[15],15:[16],16:[17],17:[18,19,20,21],18:[22],19:[23],20:[24,25],21:[26,27,28,29],22:[30,31],23:[32],24:[33],25:[34],26:[35],27:[36,37,38,39],28:[40,41],29:[42,43],30:[44,45,46,47],31:[48],32:[49,50,68],33:[51],34:[52],42:[59],43:[60],44:[61],45:[62],48:[63],46:[64],47:[65],41:[66],49:[67]}.items():
    for n in nums:mapping[n]=a
for n in range(1,69):
    ident='VAL-%03d'%n
    if 53<=n<=57:
        r.append('| '+ident+' | BLOCKED | Human | B1-H; [human-checklist.md](human-checklist.md) mapped rows. G14 and binary-probes prove feasibility only. Interim 054/055 answers are preserved in review-10, explicitly not phase-end confirmation. |')
    else:
        q=mapping[n]
        r.append('| '+ident+' | PASS | '+validators[ident]+' (QA Lead final audit) | QA-%02d: '%q+qa_notes[q]+' |')
r+=['','## Requirement Compliance','',
'| Requirement | Status | Evidence |','|---|---|---|',
'| PRD-FOUND-003; PRD-REC-001/004; schema section 8.5 | MET | VAL-001..004/050/058; one-column and down dumps; replay/manifest/parity tests |',
'| PRD-OPS-011; Source/repository contract sections 2-6 | MET | VAL-017..032/061; every changed/observing leaf and unsafe/refusal controls |',
'| PRD-OPS-012 | MET | VAL-002/006..009/066; live 49-leaf corpus, canonical/presence/replay controls |',
'| PRD-OPS-013 | PARTIAL | Deterministic VAL-010..016/049/059/060 PASS; Human VAL-057 and take-control judgment 053 pending |',
'| PRD-OPS-014 | PARTIAL | VAL-036..039/049/064/065/067 PASS; phase-end Human 054 pending |',
'| PRD-PROJ-006 | PARTIAL | VAL-040/041/068 PASS; phase-end Human 055 pending |',
'| PRD-PROJ-007; PRD-OPS-007 | MET (deterministic scope) | VAL-045..048/063; referenced-file sentinel, escaped controls and sensitive generic-fallback guards |',
'| PRD-HUM-001; Phase 23 restated guidance | PARTIAL | VAL-044 PASS; VAL-056 readability/section judgment pending |',
'| PRD-EVID-001/003; PRD-GATE-003/006; PRD-RUN-004/005; PRD-ENV-003; PRD-ISSUE-011 | MET (phase scope) | VAL-018/026..034/062; recorded observation/freshness and dual-state controls |',
'| TRD section 2; dependency budget | MET | Empty go.mod/go.sum diff; no new package boundary; existing owning domain/adapters retained |',
'| Workflow 2.6/8; accepted inventory header; quality evidence rules | NOT MET | QA-36 F1/F2 and B1-H; deterministic gate successes do not close missing owner receipts or Human judgments |',
'','## Deterministic Checklist Execution','',
'Every QA-01 through QA-49 was executed or independently verified against captured command/file-state evidence. Test-based checks cite their fresh executable assertions, not code-reading alone.',
'','| QA ID | Result | Evidence |','|---|---|---|']
for i in range(1,50):r.append('| QA-%02d | %s | %s (qa-lead) |'%(i,'FAIL' if i==36 else 'PASS',qa_notes[i]))
r+=['','## Command And Feasibility Validation Result','',
'Planning validation completed before implementation: YES, qa-plan section 6 and readiness section 12.7. This review additionally executes every required gate and checks the actual final scopes, registries, fixture paths and steps.',
'','| Command/step | Validation method | Result/evidence |','|---|---|---|',
'| Exact full phase, checkptr, archive scale and three samples | Execute approved commands, wait for real termination | G01-G13; every exit_code 0; archive table escape decoded to regex alternation per DEC-291 |',
'| Go race/tags/gcflags/timeout/coverage/bench/build, lint, format, rg | Current applicable help plus actual executions | inventory-and-help.json; go/lint/rg help exit 0; cover help exits 2 with usage; gofmt help exits 0; usage observations are not product-gate FAIL/PASS |',
'| Report-lint syntax/coverage/commit identity | Current usage and all 15 exact report/spec invocations | report-lint.json and document-audit.json; all OUTCOME COMPLETE, terminating 0; F1/F2 are additional inventory-header duties outside the spec parser |',
'| R1-R20 and live walk | Execute rendered commands and verbatim section-1 Bash block | inventory-and-help.json, inventory-*.txt; expected absence statuses recorded without rewriting them to zero |',
'| Fresh schema fixtures | Apply all 19 actual migration files with SQLite CLI, dump objects and cmp down | schema-proof.json and three dumps; CLI supports accepted trigger syntax; migration unit test separately proves fresh-init/migrator equivalence |',
'| Named tests and package paths | Existing file/function inventory plus fresh -v execution | S1/S2; no no-tests-to-run stand-in, no omitted wider gate |',
'| Human manifest/binary/workspaces/input paths/argv | File checks, SHA-256, go version -m, command-to-argv comparison, independent Human replay | binary-probes.json: 59 unique rows, 56 exact output/status replays, 3 authoring help/templates; prepared authoring state retained by generator rehearsal |',
'| All 39 help/example steps | Real binary from non-project fresh directories, parse JSON template, before/after tree snapshots | binary-probes.json; exit 0, empty stderr, no filesystem effect; descriptor/reflection/pinned-help evidence in S1 |',
'| Journey and authoring action feasibility | Generator rehearses accepted actions then restores pre-state; real Human-only journey run with source ledger | G14 and S2; Actor IDs/project directory are supplied setup, all other values from allowed provenance classes |',
'',
'QA support setup failures remain recorded: A1 assumed every old review had a Verdict heading (AttributeError, exit 1); A2 used Python’s older SQLite, which rejected a retained trigger CASE expression (OperationalError, exit 1); A3 terminated 0 with an incomplete registry inventory; A4 constructed an authoring --example command at the root rather than the leaf (AssertionError, exit 1); A5 expected direct SQL application to create the migrator’s metadata index (AssertionError, exit 1). The helpers were corrected without changing any protected definition, and the byte-identical helper invocations subsequently terminate 0. The final schema proof retains and explicitly asserts the sole bookkeeping index difference; the fresh migration test independently checks all migrator objects. Original failures, diagnostic reproduction and partial receipts are preserved in '+link('support-attempts.md')+', attempt JSON files and complete outer execution receipts. A6 preserves an additional pre-commit whitespace-check failure (exit 2) and the formatting repair, with exact test streams retained in JSON; its byte-identical passing check receipt is '+link('precommit-check-final.json')+'. None is reclassified as a successful product gate. Earlier exploratory nonexistent-path reads are not gate substitutions.',
'',
'Inventory deltas: task 12 adds exactly concern.get/deviation.get/human-outcome.get to the 92 baseline leaves (95 now); groups 29, inputs 39, keyed 50 including export import (49 digest leaves), no-short 0. The 39 example flags add booleans to the 22 earlier typed flags (61 now). R2 now lists only ComputeRequestDigest; R5 no closed decoder remains. R6/R7 residual domain fact/validation references use observed repositories, not removed public identity fields. R8 has accepted repository-unavailable and operation-unavailable allocations. R9/R10 changes reflect renderer/converged assertions; R11/R12 include labelled repository/Stop views; R13 returns qualifier-test references, no PRIVATE_ values; R14 includes the three detail reads. R19 is exactly the four retained family request fallbacks. R20’s 18 generic sites are covered; removed grammar/specific sites are individually accounted by the corpus.',
'','## Mid-Phase Review Disposition','',
'| Review/window | Outcome | Must-fix disposition | Fresh re-review evidence | Blocker / existing ESC |','|---|---|---|---|']
reviewrows=[
(1,'A','REPAIR','Task 1 F1/F2 replay guard/archive entry; Task 2 F3 deleted-file control; Task 3 F4/F5 exclusivity/structure','review-2 resolves all five','none'),
(2,'A','PASS','All review-1 Must-fix closed','Fresh independent re-review','none'),
(3,'B','REPAIR','Tasks 4/5/6 eight evidence/registry/oracle/Source/help findings','review-4 resolves F1..F8','none'),
(4,'B','PASS','All review-3 Must-fix closed','Fresh independent re-review','none'),
(5,'C','REPAIR','Task 7 dirty observation oracle; Task 8 pre-spawn/persisted oracle','review-6 resolves both','none'),
(6,'C','PASS','All review-5 Must-fix closed','Fresh independent re-review','none'),
(7,'D1','REPAIR','Task 9 activation identity Unicode scalar boundary','review-8 resolves F1','none'),
(8,'D1','PASS','Review-7 Must-fix closed','Fresh independent re-review','none'),
(9,'D2','BLOCK','Missing interim Human answers only','review-10 independently verifies original answers and clarification','ESC-24-007 closed by DEC-288'),
(10,'D2','PASS','Interim checkpoint closed; phase-end Human scope retained','Mark Bell original a23b9f33 and fcc60def clarification; walkthrough unchanged','No remaining interim blocker'),
(11,'E','REPAIR','Task 13 owned R13 receipt (no VAL)','review-12 resolves F1 only; old failed literal-pipe receipt preserved','ESC-24-009 closed by DEC-291'),
(12,'E','PASS','Review-11 F1 closed','Committed d63c9bae6a343e49ae46146a597edf7f0620281c; report a9aa99f93ec753ff954711f39b5467507385e11e','none')]
for i,w,o,d,e,b in reviewrows:r.append('| [review-%d](review-%d.md) / %s | %s | %s | %s | %s |'%(i,i,w,o,d,e,b))
r+=['',
'Document-audit.json records each review’s actual Git receipt and clean state. Historical REPAIR/BLOCK artifacts remain unchanged. No finding is closed here merely by a Developer report. The new F1/F2 are separate final-review defects, not reopened or edited mid-phase verdicts.',
'','## Documentation Review','',
'All 15 completion reports exist and equal their committed HEAD blobs. Each claimed implementation hash resolves, its subject matches and its Refs footer names the correct phase/task. Report-lint independently checks exact task-spec command coverage and expected non-zero/termination receipt grammar. Inspecting the additional accepted inventory ownership rule found F1/F2; therefore completion-report compliance is incomplete.',
'',
'Task amendments are present in 03_TRD_SCHEMA section 8.5 (schema 19/guard/down), 05_DESIGN_PLAN section 6 (Human failure order), 06_DESIGN_PLAN_CLI (all descriptors/examples/grammar/cause/gets/recorded reasoning and DEC-290 resume), 07_DESIGN_PLAN_FLOWS (observed repository and expectation flows), 08_QUALITY_PLAN section 4 (retained Human journey/parity), and GLOSSARY (SourceInput and reasoning/prose). The phase PRD identifiers remain the durable requirements. The QA-plan’s old “quality section 5.2” locator has no corresponding section; binding current sections 3-8 are used and known CF-033 retains the locator cleanup.',
'',
'Glossary comparison: Concern summary/reasoning/disposition reason; Finding summary/every disposition reason; Deviation expected/observed/reasoning; Human Outcome reasoning all match contract section 10. Referenced project-file content is project-owned prose and is excluded from Human output; metadata remains visible. Sentinel and escaping matrices independently pass.',
'',
'TL wrap-up inputs: consolidate 06_DESIGN_PLAN_CLI’s “delivery in progress” amendment language after actual final acceptance, update GLOSSARY’s stale last-updated date, correct the quality locator/current progress scheduling under CF-033, and reconcile seven-spec/glossary phase-end state. These are scheduled wrap-up synchronization, not authority for QA to edit specs.',
'','## Human Checklist Readiness','',
'Generated [human-checklist.md](human-checklist.md). Contains only human perception/product/dogfooding rows: YES. One observable per row: YES. VAL-053..057 mapped: YES. Each judgment view and each repository-view state is fanned out; deterministic counts, grammar, template effects and sentinel scans remain in this review.',
'',
'The fresh fixture snapshot is '+code(probes['manifest_path'])+'. Manifest SHA-256 '+code(probes['manifest_sha256'])+'; binary SHA-256 '+code(probes['binary_sha256'])+'. Embedded binary identity and dirt flag are retained in binary-probes.json. Task 14 implementation to reviewed HEAD has no product diff. G14 rehearses three authoring requests in the bound physical root and restores their pre-action .fid state; the Human receives usable first submissions. Claims are generated with 86400-second TTL; regenerate with the exact G14 command before using expired workspaces. Captured Human outputs remain durable for perception checks.',
'',
'Interim Task 11 answers explicitly say they do not replace phase-end HUMAN-02/HUMAN-03. No final Human result exists in inspected artifacts. Agent simulation and deterministic journey success do not supply that judgment.',
'','## Coverage Caveat','',
'Deterministic PASS proves the listed mechanics, data states, captured command outputs and structural conformance. It does not prove appearance, tone, workflow feel, product acceptance or a Human’s ability to complete/author the flow. Those remain unconfirmed until the Human checklist records outcomes.',
'','## Release Maturity Status','',
'| Maturity area | Status | Evidence or deferral |','|---|---|---|',
'| Baseline artifacts | BLOCKED | QA-36 owner receipts plus five Human assertions unresolved |',
'| Automated phase gates | READY | All G01-G13; current coverage diagnostics; security regression evidence |',
'| Held-out scenarios | N/A | No accepted harness; no excluded material accessed |',
'| External release readiness | DEFERRED | This phase is pre-release; SAST/scenario qualification, threat model/release checklist and capability work remain roadmap maturity scope |',
'','## Findings','','### Must-fix','',
'**F1 — Task 1 missing owned R4 receipt (no VAL).** Reproduction: inspect input-inventory.md section 2 owner R4 = Task 1 and its accepted header requiring a command receipt before close; scan the entire committed task-1-done.md, including retained history, for the rendered command. Document-audit.json records receipt_present=false. The report includes R3 but neither ls migrations nor the ProjectionSchemaVersion/projection_schema search. Current report-lint PASS does not cover this separate inventory rule. Same bounded report-evidence category as review-11 F1.',
'',
'**F2 — Task 10 missing owned R10 receipt (no VAL).** The accepted header and section 2 assign R10 to Task 10. Its committed report has R9/R20 receipts but no rendered R10 Human-failure assertion search, in either current or quoted prior evidence. Document-audit.json records receipt_present=false. Deterministic renderer/corpus behavior passes; the owner command/evidence requirement remains unfulfilled.',
'',
'**B1-H — Phase-end Human-only validation unresolved.** VAL-053..057 require actual Human outcomes; current checklist cells are blank. Review-10 and the walkthrough establish interim 054/055 only. Workflow section 2.6 explicitly reserves BLOCK for unresolved Human-only validation. This is the sole protected reason for the final BLOCK; F1/F2 are bounded repairs.',
'','### Improvement','',
'Keep existing ordinary-gate headroom warnings visible: Task 10 484.734 s; Task 12 reports earlier 484.599 s; Task 13 preserves 530.445 s; Task 14 final 488.166 s and earlier 530.567 s. All are below 600; do not turn cached runs into fresh measurements or change budgets. The protected 100k benchmark contains no keyed Events (DEC-158, accepted I-07); its PASS does not establish keyed-event scale performance.',
'',
'Wrap-up doc synchronization above and the already recorded Runner task-3 hash/coverage transcription improvements remain nonblocking provenance cleanup; QA has not edited the ledger or current docs.',
'','## Prior Finding Resolution','',
'N/A — initial final review. The mid-phase closure table cites independent reviewers. Any follow-up must cite this artifact unchanged and write the next unused final-review-R.md.',
'','## Repair Plan','',
'| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |','|---|---|---|---|---|---|',
'| F1 | phase-24/task-1 | no VAL; QA-36 | Only tasks/task-1-done.md; execute rendered R4 and add actual token/status/observed termination, root cwd and missing-receipt note; preserve all history. No product/spec/test/command change. | '+code('ls migrations; rg -n "ProjectionSchemaVersion|projection_schema" internal -g \'!*_test.go\'')+'; nonempty output/exit 0; commit report under own authorship; exact report-lint invocation against task-1.md, OUTCOME COMPLETE/exit 0 | Fresh independent QA Lead reruns R4 and report-lint, verifies only authorized report delta and immutable history, QA-36. Report-only change does not invalidate full product gates. |',
'| F2 | phase-24/task-10 | no VAL; QA-36 | Only tasks/task-10-done.md; execute rendered R10 and add actual token/status/observed termination, cwd and missing-receipt note; preserve failures, scope/provenance and gates. | '+code('rg -n \'"error: \' cmd internal -g \'*_test.go\'')+'; nonempty output/exit 0; commit report; exact report-lint against task-10.md, OUTCOME COMPLETE/exit 0 | Fresh independent QA Lead reruns R10 and report-lint, report delta/history/identity and QA-36. No invalidated full phase gate for a pure report repair. |',
'',
'These findings remain open until that new reviewer verifies the evidence in a new final-review artifact. QA executions in inventory-and-help.json do not authorize the Developer to transcribe another session’s receipt as its own.',
'','## Blocker Candidates','',
'**B1-H: protected product-owner judgment / Human-only validation.** Evidence: VAL-053..057 and QA-plan section 9; blank phase-end checklist; review-10’s explicit interim limitation; workflow 2.6 final outcome rule. Authority boundary: an agent cannot approve Human usability, dogfooding or authoring judgments. Requested action: Human completes the prepared checklist from the fresh public fixtures, records their own outcomes/notes, and the Tech Lead/Runner arranges fresh independent final review after both report continuations. QA allocates no ESC identity and edits no escalation register. Existing ESC-24-007 is closed and is not reused as a new unresolved identity.',
'','## Unresolved Risks And Carry-Forward Candidates','',
'No new CARRY-FORWARD replaces required phase assertions. Known CF-040 layout/evidence-stale guidance, CF-041 hosted integration automation, CF-046 parallel-worktree refusal and CF-047 flag-combination scope retain their existing later-phase ownership. Missing final Human acceptance and F1/F2 must close before phase advancement.',
'','## Final Decision','',
'Must-fix bounded findings: 2. Protected blocker candidates: 1 (Human-only). Deterministic VAL outcomes: 63 PASS. Human VAL outcomes: 5 BLOCKED. QA checklist: 48 PASS, QA-36 FAIL. Historical window repairs are independently closed; this review does not alter those artifacts.',
'',
'Outcome: BLOCK','']
(BASE/'qa/final-review.md').write_text('\n'.join(r))
# Human checklist: setup text contains no gate rows, deterministic judgments,
# raw JSON internals, pre-authored solutions or reference input content.
h=['# Phase 24 Human Verification Checklist','',
'> Phase: 24 - Input authoring and Human take-control','> Date: 2026-10-11','> Tester: _______________','',
'This is the phase-end Human subset. Record pass/fail/partial, quoted output, attempts and your own explanation in each row. Prior Task 11 answers remain interim. Exact output/count/grammar/privacy checks are recorded by QA in final-review.md.',
'',
'Use the fresh binary '+code(fixture['binary'])+' and each workspace below. For live journey/authoring, stay in that original physical root. Source project files belong inside that root; write temporary --input JSON files outside it, for example in a fresh /tmp directory. Read only Human fid output, --help/--example, and Git commands the help names. Acting Actor IDs and project directory are supplied setup; new record IDs, keys, reasoning and files are your choices.',
'',
'If the prepared Claims have expired, QA reruns the exact generator command in final-review.md and refreshes workspace references. Captured Human-only text remains available below for view judgments. Do not use the internal QA manifest or completed test input files to obtain journey/authoring values.',
'','## Human takes control - VAL-053','']
start=next(f for f in fixture['fixtures'] if f['name']=='take-control-start')
h+=['Workspace: '+code(start['workspace'])+'. Acting Actors: '+code('human')+' and '+code('agent:builder')+'. Begin with '+output(start)+' showing the already-refused terminal-success step. Discover the rest from fid --help, the refusal, and subsequent Human reads. Do not copy the setup refusal’s transition arguments as your own solution. Record each source line used.',
'','| Test ID | Covers | Test | Expected Result | Outcome | Notes |','|---|---|---|---|---|---|']
journey=[
('a','From the initial refusal alone, explain why the step was refused.','You can state the protected Decision reason.'),
('b','From that refusal alone, choose the next legitimate command.','The output gives you an actionable next step.'),
('c1','Find the acting Issue using Human fid reads.','The Issue identifier is understandable without outside knowledge.'),
('c2','Find its Envelope using Human fid reads.','The Envelope identifier and scope are understandable without outside knowledge.'),
('d','Author and submit the Decision Reference using discovered help/template and your own contained project file.','You can supply all required values from permitted sources.'),
('e','Author and submit the approving Human Outcome using discovered help/template.','You can record the intended approval without outside assistance.'),
('f','Read the Stops after the approval pair and judge whether resolution is clear.','You can tell the blocker is resolved.'),
('g','Discover and link the qualifying Human Outcome Evidence.','The Human output and help let you choose the required linkage.'),
('h','Use the latest Human resume to choose the terminal-success state.','The available completed state is understandable from the output.'),
('i','Use the latest Human resume to choose the current expected Issue revision.','You can find the precondition without JSON or guessing.'),
('j','Retry the refused terminal transition using your discovered values.','You can move the Issue to terminal success.'),
('k','Discover and submit envelope completion.','You complete the Envelope from fid’s own output/help.'),
('l','Identify any step that required source, specs, JSON internals, .fid contents or other outside knowledge.','No step requires a forbidden source; name any gap.')]
for sub,test,expected in journey:h.append('| HUMAN-01'+sub+' | VAL-053 | '+test+' (human) | '+expected+' | | |')
h+=['','## Failure clarity - VAL-054','',
'For each captured denial, write the reason and the next action separately. Judge from the Human text alone.',
'','| Test ID | Covers | Test | Expected Result | Outcome | Notes |','|---|---|---|---|---|---|']
denial_codes=['idempotency-conflict','repository-mismatch','invalid-argument','protected-decision-required','source-digest-mismatch']
denials=[f for f in fixture['fixtures'] if f.get('code') in denial_codes and f['name']!='take-control-start']
for i,f in enumerate(denials,1):
    for sub,test,expected in [('r','Explain why '+f['code']+' refused this request.','Reason is clear from the output.'),('n','Choose the next legitimate action after '+f['code']+'.','Next action is clear from the output.')]:
        h.append('| HUMAN-02-%02d%s | VAL-054 | %s %s (human) | %s | | |'%(i,sub,test,output(f),expected))
h+=['','## Stop, Attention and repository views - VAL-055','',
'Accept or reject each individual view/state, explaining what made it clear or unclear.',
'','| Test ID | Covers | Test | Expected Result | Outcome | Notes |','|---|---|---|---|---|---|']
views=[f for f in fixture['fixtures'] if not f.get('code') and not f['name'].startswith(('judgment-','authoring-')) and f['operation'] in ['stop.list','attention.list','resume','review.get']]
for i,f in enumerate(views,1): h.append('| HUMAN-03-%02d | VAL-055 | Judge %s in fixture %s: %s (human) | You accept this view’s readability and repository context. | | |'%(i,f['operation'],f['name'],output(f)))
h+=['','## Judgment presentation - VAL-056','',
'Review each individual judgment view below. Structural disclaimer counts and escaping already have deterministic evidence; here judge whether the section labels, purpose, qualifier and reasoning communicate their intended meaning.',
'','| Test ID | Covers | Test | Expected Result | Outcome | Notes |','|---|---|---|---|---|---|']
judgments=[f for f in fixture['fixtures'] if f['name'].startswith('judgment-')]
criteria=[('a','Do the labels make the judgment section boundaries clear?','Sections are easy to identify.'),
('b','Is each displayed section’s purpose understandable?','Purpose is clear.'),
('d','Does the view keep routine progress narration separate from the judgment you are reading?','Judgment content is easy to find without routine narration.'),
('e','Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?','You can understand the recorded reasoning or its absent presentation.')]
for i,f in enumerate(judgments,1):
    for sub,test,expected in criteria: h.append('| HUMAN-04-%02d%s | VAL-056 | %s %s: %s (human) | %s | | |'%(i,sub,f['operation'],test,output(f),expected))
    kinds={'concern':'Concern','deviation':'Deviation','human-outcome':'Human Outcome'}
    kind=f['operation'].split('.')[0]
    qualified=list(kinds.values()) if f['operation']=='resume' else [kinds[kind]] if kind in kinds else []
    for j,section in enumerate(qualified,1):
        h.append('| HUMAN-04-%02dc%d | VAL-056 | Judge how the %s section communicates its judgment/proof limits (or its empty state): %s (human) | The qualifier is clear without distracting repetition; an empty section is understandable. | | |'%(i,j,section,output(f)))
h+=['','## Author inputs from help - VAL-057','',
'Use your own input files and new keys; do not read the generator’s completed input. Each workspace is restored to its pre-action state. The acting Actor is supplied per row below. Discover domain identifiers from Human reads. For the third leaf, a concern.record workspace is prepared; choose it or another valid --input leaf available in that workspace and record your choice.',
'','| Test ID | Covers | Test | Expected Result | Outcome | Notes |','|---|---|---|---|---|---|']
author=[f for f in fixture['fixtures'] if f.get('authoring')]
for i,f in enumerate(author,1):
    clause={'authoring-run-close':'Author run close with expected_repository.','authoring-decision-reference':'Author decision record-reference with project-file digest omitted.','authoring-third-leaf':'Choose and author a third --input leaf.'}[f['name']]
    actor=f['human_args'][f['human_args'].index('--actor')+1]
    h.append('| HUMAN-05'+chr(96+i)+' | VAL-057 | '+clause+' Workspace '+code(f['workspace'])+'. Acting Actor '+code(actor)+'. Start from '+code(fixture['binary']+' '+f['operation'].replace('.',' ')+' --help')+' and --example. (human) | Accepted on your first or second try using only permitted guidance. | | Record chosen leaf, attempts, sources and final input. |')
h+=['','## Summary','',
'| Section | Total | Pass | Fail | Partial | Blank |','|---|---:|---:|---:|---:|---:|']
judgment_rows=sum(1 for line in h if line.startswith('| HUMAN-04-'))
counts=[('Take control',len(journey)),('Failure clarity',len(denials)*2),('Repository views',len(views)),('Judgment presentation',judgment_rows),('Input authoring',3)]
for name,n in counts:h.append('| '+name+' | '+str(n)+' | 0 | 0 | 0 | '+str(n)+' |')
total=sum(n for _,n in counts)
h+=['| **Total** | '+str(total)+' | 0 | 0 | 0 | '+str(total)+' |','',
'| | |','|---|---|','| **Overall result** | Pending Human |','| **Date completed** | |','| **Notes** | All cells intentionally blank; QA supplies feasibility, the Human supplies judgment. |','']
(BASE/'qa/human-checklist.md').write_text('\n'.join(h))
print('assembled final review: 68 VAL rows; 49 QA rows; Human-only rows',total)
