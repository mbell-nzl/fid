#!/usr/bin/env python3
"""Render the reviewed manifest and observe every published fixture invocation."""
from pathlib import Path
import hashlib, json, os, re, shlex, shutil, subprocess, sys, tempfile, time

OUT=Path(__file__).resolve().parent
QA=OUT.parent
ROOT=OUT.parents[5]
manifest_path=Path(sys.argv[1])
m=json.loads(manifest_path.read_text())
head='32a55a94efb2a89dbf50b29accb434b9c753418c'
assert m['source_commit']==head
binary=Path(m['binary'])
assert 'sha256:'+hashlib.sha256(binary.read_bytes()).hexdigest()==m['binary_sha256']
(OUT/'human-fixture-manifest.json').write_bytes(manifest_path.read_bytes())
generation="FID_P23_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-23-human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23HumanFixtures$' -count=1 -v ./cmd/fid/integration"
commands=[]
def invocation(cwd,args,expected=0,human=None,facts=None,label=''):
    row=dict(cwd=cwd,args=[str(binary),*args],expected_exit=expected,label=label)
    if human is not None: row['expected_human']=human
    if facts is not None: row['json_facts']=facts
    commands.append(row)
    return f"Working directory: \x60{cwd}\x60. Expected exit: \x60{expected}\x60.\n\n\x60\x60\x60bash\n{shlex.join(row['args'])}\n\x60\x60\x60\n"
def raw_block(text,kind='text'):
    return '\n\x60\x60\x60'+kind+'\n'+text.rstrip()+'\n\x60\x60\x60\n'
identity=f"""- Reviewed source HEAD: \x60{head}\x60.
- Generated at: \x60{m['generated_at']}\x60.
- Checkpoint root: \x60{manifest_path.parent}\x60.
- Manifest: \x60{manifest_path}\x60.
- Binary: \x60{binary}\x60.
- Binary SHA-256: \x60{m['binary_sha256']}\x60.
- Source dirty flag: \x60{str(m['source_dirty']).lower()}\x60; only workflow/QA documents are dirty. Product source is the reviewed HEAD.
- Durable manifest: [human-fixture-manifest.json](final-evidence-1/human-fixture-manifest.json).
"""
guide=f"""# Phase 23 Human Fixtures

Reviewed on 2026-10-03 by gpt-6-astra QA Lead. This guide replaces the earlier
fixture snapshots with one fresh manifest at the reviewed HEAD, including
DEC-267 and tasks 19–21 (DEC-269/DEC-270).

## Generate a fresh set

Working directory: \x60/home/mark/projects/fid\x60.

\x60\x60\x60bash
{generation}
\x60\x60\x60

Wait for the command to finish. Its output names a new manifest with current
paths and IDs. Seeded active Claims last one hour, so generate immediately before a later
Human session and use that manifest's commands. Exact expected text below is
bound to this recorded generation. The original [human-checklist.md](human-checklist.md)
is the immutable historical Human result. Record only the seven re-checks in
[human-checklist-2.md](human-checklist-2.md).

## Reviewed fixture identity

{identity}
Generation receipt, including the complete terminal result and wall-time
observation: [human-fixture-generation.json](final-evidence-1/human-fixture-generation.json).
Every invocation below was executed by QA; receipt:
[human-guide-verification.json](final-evidence-1/human-guide-verification.json).
These receipts establish runnable mechanics and exact text, not Human judgment.

The protected-Decision and rework steps change their fixtures. Run inspections
before the steps and use a fresh generation when repeating a complete sequence.
The leaf mutations repeat existing requests idempotently.
"""
rework_blocks=[]
for cp in m['checkpoints']:
    guide+='\n## '+cp['name']+'\n\n'
    guide+='Fixture IDs:\n'+raw_block(json.dumps(cp['ids'],indent=2,sort_keys=True),'json')
    if 'VAL-072' in cp['name']:
        guide+='\nRead the denial and Stop/Attention before consulting the resolution step. Record the three rework judgments from what fid displays.\n\n'
        step=cp['steps'][0]
        human_args=[a for a in step['args'] if a!='--json']
        block=invocation(step['cwd'],human_args,step['expected_exit'],label='rework Human denial')
        guide+='\n### Human rework denial\n\n'+block
        rework_blocks.append(('Human rework denial',block))
    for index,inspection in enumerate(cp['inspections'],1):
        args=inspection['args']
        human=cp['expected_human'].get('\0'.join(args))
        block=invocation(inspection['cwd'],args,human=human,label=cp['name']+' inspection '+str(index))
        guide+=f'\n### Inspection {index}\n\n'+block
        if human is not None:guide+='\nExpected Human output:\n'+raw_block(human)
        if 'VAL-072' in cp['name'] and '--json' not in args:
            rework_blocks.append(('Inspect '+('resume' if 'resume' in args else 'Stop' if 'stop' in args else 'Attention'),block))
    for index,step in enumerate(cp.get('steps',[]),1):
        block=invocation(step['cwd'],step['args'],step['expected_exit'],label=cp['name']+' step '+step['name'])
        guide+=f'\n### Step {index}: {step["name"]}\n\n'+block
        if 'VAL-072' in cp['name']:
            if index==2:
                guide+='\nTask 21 resolution: this activated Human records an approved \x60review-rework-extension\x60 Outcome. Its Event resolves the matching rework Stop and Attention atomically. No internal file editing is involved.\n'
            rework_blocks.append((f'Step {index}: {step["name"]}',block))
guide+='\n## Every new leaf — Human and JSON pairs\n'
leaf_blocks={}
for index,pair in enumerate(m['leaf_pairs'],1):
    guide+=f'\n### {index:02}. {pair["operation"]}\n\n'
    block=invocation(pair['cwd'],pair['human_args'],human=pair['expected_human'],label=pair['operation']+' Human')
    leaf_blocks[pair['operation']]=block
    guide+=block+'\nExpected Human output:\n'+raw_block(pair['expected_human'])
    guide+='\nJSON counterpart:\n\n'+invocation(pair['cwd'],pair['json_args'],facts=pair['json_facts'],label=pair['operation']+' JSON')
    guide+='\nExpected JSON result facts:\n'+raw_block(json.dumps(pair['json_facts'],indent=2,sort_keys=True),'json')
(QA/'human-fixtures.md').write_text(guide)
(OUT/'human-guide-commands.json').write_text(json.dumps(commands,indent=2)+'\n')

old=(QA/'human-checklist.md').read_text()
selected=['HUMAN-03a','HUMAN-03b-resolution','HUMAN-03b-internals','HUMAN-04a-17','HUMAN-04b-17','HUMAN-04c-17','HUMAN-04a-18']
rowmap={}
for line in old.splitlines():
    if line.startswith('| HUMAN-'):
        cells=line.split('|')
        if cells[1].strip() in selected:
            rowmap[cells[1].strip()]='|'+'|'.join(cells[1:5])+ '| | |'
assert set(rowmap)==set(selected)
table='| Test ID | Covers | Test | Expected Result | Outcome | Notes |\n|---|---|---|---|---|---|\n'
checklist=f"""# Phase 23 Human Verification Checklist — Re-check 2

> Phase: 23 — Review and Human judgment
> Date: 2026-10-03
> Tester: _______________

**Pending Human input: seven blank judgments.** This is the DEC-269 re-check
scope. The [original Human checklist](human-checklist.md) remains immutable,
including its 55 pass, 2 fail and 6 partial results. Deterministic evidence
belongs in [final-review-1.md](final-review-1.md); it does not fill these cells.
That review requires task-16 repair of the full-checkptr gate before phase
acceptance. These fixtures still document the actual reviewed HEAD.

## Fixture identity

{identity}
Use the [fresh-generation instructions](human-fixtures.md#generate-a-fresh-set)
if the one-hour Claims have expired. After regeneration, substitute the new
manifest's binary, working directories, IDs and input paths.

## Review rework exhaustion — VAL-072

First run the Human denial and inspect resume, Stop and Attention. Record what
the displayed output tells you before reading or executing the resolution step.

"""
for title,block in rework_blocks[:4]:
    checklist+='### '+title+'\n\n'+block+'\n'
checklist+=table+'\n'.join(rowmap[x] for x in selected[:3])+'\n'
checklist+='\n### Task 21 resolution and admitted retry\n\nFollow these steps in order after recording the initial observations. The approved Human Outcome resolves the matching Stop and Attention atomically; the retry then acquires a Claim.\n\n'
for title,block in rework_blocks[4:]:
    checklist+='#### '+title+'\n\n'+block+'\n'
checklist+='\n## Requirement get — VAL-073\n\n'+leaf_blocks['requirement.get']+'\n'+table+'\n'.join(rowmap[x] for x in selected[3:6])+'\n'
checklist+='\n## Requirement list — VAL-073\n\n'+leaf_blocks['requirement.list']+'\n'+table+rowmap[selected[6]]+'\n'
checklist+="""
## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Review rework exhaustion | 3 | 0 | 0 | 0 | 3 |
| Requirement get | 3 | 0 | 0 | 0 | 3 |
| Requirement list | 1 | 0 | 0 | 0 | 1 |
| **Total** | **7** | **0** | **0** | **0** | **7** |

| | |
|---|---|
| **Overall result** | Pending Human input |
| **Date completed** | |
| **Notes** | ESC-23-012 remains open. A4–A8 and remaining original PARTIAL items stay on the Tech Lead wrap-up agenda; this checklist does not accept a carry-forward. |

Outcome values: pass, fail, partial, blank. Only the Human records judgments.
"""
(QA/'human-checklist-2.md').write_text(checklist)

# Observe the exact published argv against the retained fixture paths, then
# restore our newly generated fixtures for the Human. Never touch prior sets.
assert manifest_path.parent.parent==ROOT/'.cache/phase-23-human'
scratch=Path(tempfile.mkdtemp(prefix='fid-p23-review1-fixture-observation-'))
projects=sorted({Path(row['cwd']) for row in commands})
def tree_hash(p):
    return {str(f.relative_to(p)):hashlib.sha256(f.read_bytes()).hexdigest()
            for f in sorted(p.rglob('*')) if f.is_file()}
original={}
for i,p in enumerate(projects):
    assert p.parent==manifest_path.parent
    original[str(p)]=tree_hash(p)
    shutil.copytree(p,scratch/str(i))
receipts=[]
try:
    for row in commands:
        assert Path(row['cwd']).is_dir()
        if '--input' in row['args']:
            assert Path(row['args'][row['args'].index('--input')+1]).is_file()
        started=time.monotonic()
        r=subprocess.run(row['args'],cwd=row['cwd'],text=True,capture_output=True)
        receipt=dict(command=row['args'],cwd=row['cwd'],label=row['label'],
                     expected_exit=row['expected_exit'],exit_code=r.returncode,
                     stdout=r.stdout,stderr=r.stderr,wall_time_seconds=time.monotonic()-started)
        receipts.append(receipt)
        (OUT/'human-guide-invocation-receipts.json').write_text(json.dumps(receipts,indent=2)+'\n')
        assert r.returncode==row['expected_exit'],receipt
        if 'expected_human' in row:assert r.stdout==row['expected_human'] and not r.stderr,receipt
        if 'json_facts' in row:assert json.loads(r.stdout)['result']==row['json_facts'],receipt
        if row['label']=='rework Human denial':
            assert not r.stdout and len(r.stderr.splitlines())==2,receipt
            assert 'review-rework-extension' in r.stderr and 'fid stop list --subject envelope:' in r.stderr,receipt
        if '--json' not in row['args']:
            assert 'PRIVATE_' not in r.stdout,receipt
        print(row['label'],f'expected={row["expected_exit"]} observed exit_code={r.returncode}',flush=True)
finally:
    for i,p in enumerate(projects):
        shutil.rmtree(p)
        shutil.copytree(scratch/str(i),p)
    restoration={str(p):tree_hash(p)==original[str(p)] for p in projects}
    (OUT/'human-fixture-restoration.json').write_text(json.dumps(dict(backup=str(scratch),restored=restoration),indent=2)+'\n')
    assert all(restoration.values())
blocks=re.findall(r'\x60\x60\x60bash\n(.*?)\n\x60\x60\x60',guide,re.S)
assert blocks[0]==generation and [shlex.split(b) for b in blocks[1:]]==[x['args'] for x in commands]
for block in re.findall(r'\x60\x60\x60bash\n(.*?)\n\x60\x60\x60',checklist,re.S):
    assert shlex.split(block) in [x['args'] for x in commands]
assert len(re.findall(r'^\| HUMAN-',checklist,re.M))==7
print(f'PASS {len(commands)} exact guide invocations; Human/JSON facts match; seven Human judgments remain blank; fresh fixtures restored for handoff')
