#!/usr/bin/env python3
"""Capture full journal rebuild and pragma evidence in a fresh QA copy."""
from pathlib import Path
import hashlib,json,shutil,subprocess,tempfile
OUT=Path(__file__).resolve().parent
ROOT=OUT.parents[5]
source=ROOT/'.cache/phase-23-human/checkpoints-3245767079/mixed-findings'
scratch=Path(tempfile.mkdtemp(prefix='fid-p23-qa-rebuild-'))
workspace=scratch/'workspace'
shutil.copytree(source,workspace)
receipts=[]
def call(args):
 r=subprocess.run(args,cwd=workspace,text=True,capture_output=True)
 row={'command':args,'cwd':str(workspace),'stdout':r.stdout,'stderr':r.stderr,'exit_code':r.returncode}
 receipts.append(row)
 (OUT/'rebuild-pragmas-receipts.json').write_text(json.dumps(receipts,indent=2)+'\n')
 print(json.dumps(row))
 assert r.returncode==0
 return r.stdout
def journal_hash():
 return {p.name:hashlib.sha256(p.read_bytes()).hexdigest() for p in (workspace/'.fid/journal').glob('*.jsonl')}
before=journal_hash()
result=json.loads(call(['/tmp/fid','--project',str(workspace),'--json','rebuild']))
assert result['ok']
assert journal_hash()==before
db=str(workspace/'.fid/fid.sqlite')
assert call(['sqlite3','-batch',db,'PRAGMA foreign_key_check;'])==''
assert call(['sqlite3','-batch',db,'PRAGMA integrity_check;'])=='ok\n'
print('Journal SHA-256 unchanged:',json.dumps(before,sort_keys=True))
print('PASS full journal rebuild: foreign_key_check=[]; integrity_check=ok; source journal unchanged')
