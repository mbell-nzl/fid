#!/usr/bin/env python3
"""Independent SQLite inventory/pragma probes using the accepted migration files."""
from pathlib import Path
import json,subprocess,tempfile
ROOT=Path(__file__).resolve().parents[6]
OUT=Path(__file__).resolve().parent
scratch=Path(tempfile.mkdtemp(prefix='fid-p23-qa-schema-'))
receipts=[]
def sql(db,text):
 args=['sqlite3','-batch',str(db)]
 r=subprocess.run(args,input='.bail on\n'+text,text=True,capture_output=True)
 receipts.append({'command':args,'stdin':text,'stdout':r.stdout,'stderr':r.stderr,'exit_code':r.returncode})
 (OUT/'schema-sql-receipts.json').write_text(json.dumps(receipts,indent=2)+'\n')
 assert r.returncode==0,(db,r.stderr)
 return r.stdout
for version in (17,18):
 db=scratch/f'schema-{version}.sqlite'
 for migration in sorted((ROOT/'migrations').glob('*.up.sql')):
  if int(migration.name[:4])>version: break
  sql(db,'PRAGMA foreign_keys=ON;\nBEGIN;\n'+migration.read_text()+'\nCOMMIT;\n')
 foreign=sql(db,'PRAGMA foreign_key_check;\n')
 integrity=sql(db,'PRAGMA integrity_check;\n')
 assert foreign=='' and integrity=='ok\n'
 cols=sql(db,'SELECT name FROM pragma_table_info("evidence_records");\n')
 objects=sql(db,'SELECT type||":"||name FROM sqlite_schema WHERE name NOT LIKE "sqlite_%" ORDER BY type,name;\n')
 print(f'schema {version}: foreign_key_check=[]; integrity_check={integrity.strip()}')
 print(f'schema {version} evidence columns:\n'+cols)
 print(f'schema {version} objects:\n'+objects)
print('PASS independently applied schema-17 and schema-18 pragma and inventory probes')
