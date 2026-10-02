#!/usr/bin/env python3
"""Run the exact reference Go argv with Git metadata bound to its worktree."""
from pathlib import Path
import json,os,subprocess,sys
OUT=Path(__file__).resolve().parent
ROOT=Path('/tmp/fid-p23-qa-reference-PgIz9Jst')
which=sys.argv[1]
workspace=ROOT/which
line=(workspace/'.git').read_text().strip()
assert line.startswith('gitdir: ')
gitdir=line.removeprefix('gitdir: ')
env=dict(os.environ,GIT_DIR=gitdir,GIT_WORK_TREE=str(workspace))
args=sys.argv[2:]
assert args[0]=='go'
for check in [['git','rev-parse','HEAD'],['git','status','--porcelain']]:
 r=subprocess.run(check,cwd=workspace,env=env,text=True,capture_output=True)
 print(json.dumps({'command':check,'cwd':str(workspace),'git_dir':gitdir,'stdout':r.stdout,'stderr':r.stderr,'exit_code':r.returncode}),flush=True)
 assert r.returncode==0
 if check[1]=='status': assert r.stdout==''
result=subprocess.run(args,cwd=workspace,env=env,text=True,capture_output=True)
record={'command':args,'cwd':str(workspace),'environment':{'GIT_DIR':gitdir,'GIT_WORK_TREE':str(workspace)},'stdout':result.stdout,'stderr':result.stderr,'exit_code':result.returncode}
(OUT/('reference-'+which+'-bound-command.json')).write_text(json.dumps(record,indent=2)+'\n')
print(json.dumps(record),flush=True)
sys.exit(result.returncode)
