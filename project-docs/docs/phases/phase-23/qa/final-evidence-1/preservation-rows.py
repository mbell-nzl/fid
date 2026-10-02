#!/usr/bin/env python3
"""Supplement the modernc named test with explicit linked-seed row dumps."""
from pathlib import Path
import json,re,subprocess,tempfile
OUT=Path(__file__).resolve().parent
ROOT=OUT.parents[5]
db=Path(tempfile.mkdtemp(prefix='fid-p23-qa-preservation-'))/'fid.sqlite'
receipts=[]
def sql(text,json_mode=False):
 args=['sqlite3','-batch']+(['-json'] if json_mode else [])+[str(db)]
 r=subprocess.run(args,input='.bail on\n'+text,text=True,capture_output=True)
 receipts.append({'command':args,'stdin':text,'stdout':r.stdout,'stderr':r.stderr,'exit_code':r.returncode})
 (OUT/'preservation-rows-receipts.json').write_text(json.dumps(receipts,indent=2)+'\n')
 assert r.returncode==0,r.stderr
 return json.loads(r.stdout) if json_mode else r.stdout
for f in sorted((ROOT/'migrations').glob('*.up.sql')):
 if int(f.name[:4])>16: break
 sql('PRAGMA foreign_keys=ON; BEGIN;\n'+f.read_text()+'\nCOMMIT;')
source=(ROOT/'internal/db/review_judgment_migration_test.go').read_text()
section=source.split('func seedSchema16ReceiptEvidence',1)[1]
seed=re.findall(r'`(INSERT[^\`]+)`',section)
assert len(seed)==14,len(seed)
sql('PRAGMA foreign_keys=ON; BEGIN;\n'+';\n'.join(seed)+';\nCOMMIT;')
before=sql('SELECT * FROM evidence_records ORDER BY id;',True)
links_before=sql('SELECT * FROM evidence_artifact_links ORDER BY evidence_id,artifact_id;',True)
sql('PRAGMA foreign_keys=ON; BEGIN;\n'+(ROOT/'migrations/0017_p23_review_judgment.up.sql').read_text()+'\nCOMMIT;')
after=sql('SELECT * FROM evidence_records ORDER BY id;',True)
links_after=sql('SELECT * FROM evidence_artifact_links ORDER BY evidence_id,artifact_id;',True)
print('schema-16 Evidence:',json.dumps(before,sort_keys=True))
print('schema-17 Evidence:',json.dumps(after,sort_keys=True))
print('schema-16 links:',json.dumps(links_before,sort_keys=True))
print('schema-17 links:',json.dumps(links_after,sort_keys=True))
assert len(before)==len(after)==1
assert before==[{k:v for k,v in row.items() if k in before[0]} for row in after]
assert links_before==links_after
added=set(after[0])-set(before[0])
assert added=={'review_id','human_outcome_id'}
assert all(after[0][k] is None for k in added)
assert sql('PRAGMA foreign_key_check;')==''
assert sql('PRAGMA integrity_check;')=='ok\n'
print('PASS linked-seed row preservation; exactly two added columns are null; foreign_key_check=[]; integrity_check=ok')
