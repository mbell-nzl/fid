#!/usr/bin/env python3
"""Observe authored QA artifacts and protected source state; no product mutation."""
from pathlib import Path
import collections, hashlib, json, re, subprocess

OUT=Path(__file__).resolve().parent
QA=OUT.parent
ROOT=OUT.parents[5]
def data(name): return json.loads((OUT/(name+'.json')).read_text())
initial=data('protected-state-initial')
receipts=[]
def git(*args):
    r=subprocess.run(['git',*args],cwd=ROOT,text=True,capture_output=True)
    receipts.append(dict(argv=['git',*args],cwd=str(ROOT),stdout=r.stdout,stderr=r.stderr,exit_code=r.returncode))
    assert r.returncode==0,receipts[-1]
    return r.stdout
assert git('rev-parse','HEAD').strip()==initial['reviewed_head']
hashes={name:hashlib.sha256((QA.parent/name).read_bytes()).hexdigest() for name in initial['sha256']}
assert hashes==initial['sha256'],hashes
allowed={str((QA/name).relative_to(ROOT)) for name in ['final-review-1.md','human-checklist-2.md','human-fixtures.md']}
prefix=str(OUT.relative_to(ROOT))+'/'
changed=git('diff','HEAD','--name-only').splitlines()
assert set(changed)<=allowed|set(initial['initial_dirty']),changed
untracked=git('ls-files','--others','--exclude-standard').splitlines()
assert all(name in allowed or name.startswith(prefix) for name in untracked),untracked
assert not git('diff','HEAD','--','cmd','internal','migrations','go.mod','go.sum')
final=dict(reviewed_head=initial['reviewed_head'],sha256=hashes,all_protected_hashes_unchanged=True,
           tracked_changes=changed,untracked=untracked,product_diff_empty=True,receipts=receipts)
(OUT/'protected-state-final.json').write_text(json.dumps(final,indent=2)+'\n')
review=(QA/'final-review-1.md').read_text()
checklist=(QA/'human-checklist-2.md').read_text()
rows=re.findall(r'^\| (VAL-\d{3}) \| (PASS|FAIL|BLOCKED|CARRY-FORWARD) \| (.+)$',review,re.M)
assert [x[0] for x in rows]==[f'VAL-{n:03}' for n in range(1,90)]
assert collections.Counter(x[1] for x in rows)==dict(PASS=86,BLOCKED=3)
assert all('](' in x[2] for x in rows)
assert all('Human' in rows[n-1][2] for n in range(70,74))
assert len(re.findall(r'^\| QA-\d{2} ',review,re.M))==52
expected=['HUMAN-03a','HUMAN-03b-resolution','HUMAN-03b-internals','HUMAN-04a-17','HUMAN-04b-17','HUMAN-04c-17','HUMAN-04a-18']
human=[line.split('|') for line in checklist.splitlines() if line.startswith('| HUMAN-')]
assert [x[1].strip() for x in human]==expected
assert all('(human)' in '|'.join(x) and not x[5].strip() and not x[6].strip() for x in human)
assert 'Task 21 resolution' in checklist and 'human-outcome' in checklist
assert '**REPAIR.**' in review and 'OUTCOME REPAIR' in review
assert data('full-checkptr')['results'][-1]['exit_code']==1
assert 'release imported Claim: claim expired' in review
links=[]
for file in ['final-review-1.md','human-checklist-2.md','human-fixtures.md']:
    text=(QA/file).read_text()
    for target in re.findall(r'\]\(([^)]+)\)',text):
        if '://' in target: continue
        path=target.split('#',1)[0]
        if not path: continue
        resolved=(QA/path).resolve()
        assert resolved.exists(),(file,target)
        links.append(dict(source=file,target=target))
        if path.startswith('final-evidence-1/') and path.endswith('.json'):
            d=json.loads(resolved.read_text())
            if isinstance(d,dict) and 'results' in d:
                assert d['results'] and 'exit_code' in d['results'][-1],target
assert data('human-fixture-manifest')['source_commit']==initial['reviewed_head']
assert len(data('human-guide-invocation-receipts'))==63
assert all(data('human-fixture-restoration')['restored'].values())
(OUT/'artifact-link-checks.json').write_text(json.dumps(links,indent=2)+'\n')
print('PASS 89 VAL outcomes; 52 deterministic rows; exactly seven blank Human rows; all local evidence links resolve; observed terminations retained; protected hashes/source unchanged')
