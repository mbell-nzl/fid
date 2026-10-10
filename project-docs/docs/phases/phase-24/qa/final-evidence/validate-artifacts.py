"""Check completeness and local-link feasibility of the two new QA artifacts.
Run from repo root after assemble.py. This is an artifact check, not a
substitute for a product quality gate.
"""
import pathlib,re,json
ROOT=pathlib.Path(__file__).resolve().parents[6]
BASE=ROOT/'project-docs/docs/phases/phase-24/qa'
review=(BASE/'final-review.md').read_text()
human=(BASE/'human-checklist.md').read_text()
vals=re.findall(r'^\| (VAL-\d{3}) \| (PASS|FAIL|BLOCKED|CARRY-FORWARD) \|',review,re.M)
assert len(vals)==68 and len({r[0] for r in vals})==68
assert {r[0] for r in vals}=={'VAL-%03d'%i for i in range(1,69)}
qa=re.findall(r'^\| (QA-\d{2}) \| (PASS|FAIL|BLOCKED) \|',review,re.M)
assert len(qa)==49 and len({r[0] for r in qa})==49
assert {r[0] for r in qa}=={'QA-%02d'%i for i in range(1,50)}
assert review.count('Outcome: BLOCK')==1
hrows=[l for l in human.splitlines() if l.startswith('| HUMAN-')]
assert hrows and all('(human)' in l for l in hrows)
assert not re.search(r'\(qa-lead\)|^\| QA-|^\| Gate',human,re.M)
assert {r for l in hrows for r in re.findall(r'VAL-\d{3}',l)}=={'VAL-%03d'%i for i in range(53,58)}
ids=[l.split('|')[1].strip() for l in hrows]; assert len(set(ids))==len(ids)
links=[]
for name,data in [('final-review.md',review),('human-checklist.md',human)]:
    for target in re.findall(r'\]\(([^)]+)\)',data):
        if target.startswith(('https://','http://')):continue
        path=(BASE/target.split('#')[0]).resolve()
        assert path.exists(),(name,target)
        links.append(dict(artifact=name,target=target,exists=True))
summary=dict(VAL_outcomes={k:sum(v==k for _,v in vals) for k in ['PASS','FAIL','BLOCKED','CARRY-FORWARD']},QA_outcomes={k:sum(v==k for _,v in qa) for k in ['PASS','FAIL','BLOCKED']},human_rows=len(hrows),human_types_only=True,local_links=links)
(BASE/'final-evidence/artifact-validation.json').write_text(json.dumps(summary,indent=2)+'\n')
print('PASS artifact completeness: 68 VAL, 49 QA, '+str(len(hrows))+' Human-only rows, '+str(len(links))+' valid local links')
