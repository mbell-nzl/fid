#!/usr/bin/env python3
"""Refresh only fixture identities/invocations; preserve all Human judgments."""
from pathlib import Path
import hashlib,json,re,sys
OUT=Path(__file__).resolve().parent
path=OUT.parent/'human-checklist.md'
old_text=path.read_text()
old_path=Path(re.search(r'- Manifest: `([^`]+)`',old_text).group(1))
new_path=Path(sys.argv[1])
old=json.loads(old_path.read_text());new=json.loads(new_path.read_text())
assert len(old['checkpoints'])==len(new['checkpoints'])==3
assert len(old['leaf_pairs'])==len(new['leaf_pairs'])==18
replacements={str(old_path.parent):str(new_path.parent),old['binary_sha256']:new['binary_sha256']}
def pair_args(a,b):
 assert len(a)==len(b),(a,b)
 for x,y in zip(a,b):
  if x!=y:
   if x in replacements: assert replacements[x]==y,(x,y)
   replacements[x]=y
for a,b in zip(old['checkpoints'],new['checkpoints']):
 assert a['name']==b['name']
 for key,value in a['ids'].items():
  if value!=b['ids'][key]: replacements[value]=b['ids'][key]
 assert len(a['inspections'])==len(b['inspections'])
 for x,y in zip(a['inspections'],b['inspections']): pair_args(x['args'],y['args'])
 assert len(a.get('steps',[]))==len(b.get('steps',[]))
 for x,y in zip(a.get('steps',[]),b.get('steps',[])): pair_args(x['args'],y['args'])
for a,b in zip(old['leaf_pairs'],new['leaf_pairs']):
 assert a['operation']==b['operation']
 pair_args(a['human_args'],b['human_args']);pair_args(a['json_args'],b['json_args'])
pattern=re.compile('|'.join(re.escape(k) for k in sorted(replacements,key=len,reverse=True)))
new_text=pattern.sub(lambda m:replacements[m[0]],old_text)
note='- These fixtures contain short-lived Claims. Generate a fresh set immediately before the Human run using the guide linked below, then use that manifest\'s current paths and IDs. Exact expected output is bound to generation time.\n'
if note not in new_text:
 new_text=new_text.replace('- Run each section from its stated directory.',note+'- Run each section from its stated directory.')
new_text=new_text.replace('The recorded dirty state consists of pre-existing workflow documents; product files were unchanged.','The recorded dirty state consists of workflow documents and these QA artifacts; product files were unchanged.')
before=[line for line in old_text.splitlines() if line.startswith('| HUMAN-')]
after=[line for line in new_text.splitlines() if line.startswith('| HUMAN-')]
assert before==after and len(after)==63
assert all(line.endswith('| | |') for line in after)
backup=OUT/'human-checklist-before-refresh.md'
version=2
while backup.exists():
 backup=OUT/f'human-checklist-before-refresh-{version}.md'
 version+=1
backup.write_text(old_text)
path.write_text(new_text)
(OUT/'human-fixture-manifest.json').write_text(json.dumps(new,indent=2)+'\n')
print('Old manifest:',old_path)
print('New manifest:',new_path)
print('Binary digest:',new['binary_sha256'])
print('Human row SHA-256:',hashlib.sha256(('\n'.join(after)+'\n').encode()).hexdigest())
print('PASS fixture identities refreshed; all 63 blank Human rows preserved byte-for-byte')
