#!/usr/bin/env python3
"""Compare task-11-close and final capability discovery on fresh workspaces."""
from pathlib import Path
import json,subprocess,tempfile
OUT=Path(__file__).resolve().parent
refs={'task11':'3dff4574a71438900100449e36d7044969ba9162','current':'32a55a94efb2a89dbf50b29accb434b9c753418c'}
print('Reference worktrees:',json.dumps(refs,sort_keys=True))
root=Path('/tmp/fid-p23-qa-reference-PgIz9Jst')
receipts=[]
def call(args,cwd):
 r=subprocess.run(args,cwd=cwd,text=True,capture_output=True)
 row={'command':args,'cwd':str(cwd),'stdout':r.stdout,'stderr':r.stderr,'exit_code':r.returncode}
 receipts.append(row)
 (OUT/'capability-staging-receipts.json').write_text(json.dumps(receipts,indent=2)+'\n')
 print(json.dumps(row))
 assert r.returncode==0,(args,r.stderr)
 return r.stdout
for label,binary,checkout,want in [('task11','/tmp/fid-p23-final-review-1-task11',root/'task11',False),('final','/tmp/fid',OUT.parents[5],True)]:
 head=call(['git','rev-parse','HEAD'],checkout).strip()
 workspace=Path(tempfile.mkdtemp(prefix='fid-p23-qa-capability-'+label+'-'))
 call([binary,'--project',str(workspace),'--json','init','--prefix','QA'],checkout)
 data=json.loads(call([binary,'--project',str(workspace),'--json','capabilities'],checkout))
 # Find the unique option rather than rely on row order.
 matching=[]
 def walk(value):
  if isinstance(value,dict):
   if 'human_outcome_resolution' in value: matching.append(value['human_outcome_resolution'])
   for child in value.values(): walk(child)
  elif isinstance(value,list):
   for child in value: walk(child)
 walk(data)
 assert matching==[want],matching
 print(label,head,'human_outcome_resolution='+str(want).lower())
print('PASS task-11-close false and final true capability states on fresh initialized workspaces')
