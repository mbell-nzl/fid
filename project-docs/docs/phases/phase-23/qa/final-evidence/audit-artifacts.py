#!/usr/bin/env python3
"""Validate the authored final review, receipt links and preserved Human scope."""
from pathlib import Path
from collections import Counter
import hashlib, json, re

OUT=Path(__file__).resolve().parent
QA=OUT.parent
review=(QA/'final-review.md').read_text()
vals=re.findall(r'^\| (VAL-\d{3}) \| (PASS|FAIL|BLOCKED|CARRY-FORWARD) \|',review,re.M)
assert len(vals)==84 and len(dict(vals))==84
assert Counter(result for _,result in vals)=={'PASS':79,'BLOCKED':5}
assert {ident for ident,result in vals if result=='BLOCKED'}=={'VAL-005','VAL-070','VAL-071','VAL-072','VAL-073'}
checks=re.findall(r'^\| (QA-\d{2}) `\(qa-lead\)` \| (PASS|FAIL|BLOCKED) \|',review,re.M)
assert len(checks)==52 and Counter(result for _,result in checks)=={'PASS':51,'BLOCKED':1}
assert dict(checks)['QA-07']=='BLOCKED'
assert review.rstrip().endswith('\nOUTCOME BLOCKED')
assert len(re.findall(r'^OUTCOME ',review,re.M))==1
links=0
for file in (QA/'final-review.md',QA/'human-checklist.md',OUT/'command-index.md'):
    for target in re.findall(r'\]\(([^)]+)\)',file.read_text()):
        if '://' in target or target.startswith('#'): continue
        path=target.split('#',1)[0]
        assert (file.parent/path).exists(),(file.name,target)
        links+=1
print(f'PASS {links} local review/checklist/index links resolve')
original=[s for s in (OUT/'human-checklist-before-refresh.md').read_text().splitlines() if s.startswith('| HUMAN-')]
current=[s for s in (QA/'human-checklist.md').read_text().splitlines() if s.startswith('| HUMAN-')]
assert len(current)==63 and current==original
assert all(s.endswith('| | |') and '(human)' in s for s in current)
print('PASS all 63 original Human rows remain blank and byte-identical')
receipt_ids=re.findall(r'final-evidence/([a-z0-9-]+)\.json',review)
for ident in set(receipt_ids):
    data=json.loads((OUT/(ident+'.json')).read_text())
    if isinstance(data,dict) and {'id','cmd','results'}<=data.keys():
        assert 'exit_code' in data['results'][-1],ident
print('PASS all linked execution receipts have an observed terminal exit_code')
for ident in ('gate-01-format','gate-02-lint','gate-03-coverage','gate-04-coverage-summary','gate-05-integration','gate-06-benchmarks','gate-07-build','budget-ordinary','budget-integration','full-checkptr','benchmark-1','benchmark-2','benchmark-3','archive-scale'):
    data=json.loads((OUT/(ident+'.json')).read_text())
    assert data['cmd'].replace('|','\\|') in review,ident
    result=data['results'][-1]
    assert result['exit_code']==(1 if ident=='full-checkptr' else 0)
print('PASS exact phase/final-control command text and real statuses match the receipts')
assert 'No retry of the failed full gate was used' in review
assert 'VERDICT: BLOCK' in review
print('PASS final artifacts: 84 VAL outcomes (79 PASS, 5 BLOCKED), 52 QA checks, open full-checkptr failure, OUTCOME BLOCKED')
