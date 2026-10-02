"""Real CLI replay smoke; requires Python 3, SQLite CLI 3.47+, and /tmp/fid.
Run from the repo root after building /tmp/fid. Uses fresh scratch storage.
"""
import hashlib
import json
from pathlib import Path
import subprocess
import tempfile

workspace = Path(tempfile.mkdtemp(prefix="fid-esc-23-002-smoke-"))
print(f"workspace={workspace}", flush=True)


def run(*args):
    command = ["/tmp/fid", "--project", str(workspace), "--json", *args]
    result = subprocess.run(command, capture_output=True, text=True)
    print(json.dumps({"command": command, "exit_code": result.returncode,
                      "stdout": result.stdout, "stderr": result.stderr}), flush=True)
    result.check_returncode()
    return json.loads(result.stdout)


def journal_hashes():
    return {p.name: hashlib.sha256(p.read_bytes()).hexdigest()
            for p in (workspace / ".fid/journal").glob("*.jsonl")}


run("init", "--prefix", "FID")
created = run("issue", "create", "--kind", "task", "--title", "Rebuild smoke",
    "--actor", "agent:developer", "--idempotency-key", "esc-23-002-smoke-create")
issue_id = created["result"]["issue"]["id"]
before = run("issue", "get", issue_id)
hashes = journal_hashes()
run("rebuild")
after = run("issue", "get", issue_id)
assert before == after, "read changed after rebuild"
assert hashes == journal_hashes(), "journal bytes changed"
for query, expected in [("PRAGMA foreign_key_check;", ""),
                        ("PRAGMA integrity_check;", "ok"),
                        ("SELECT version, dirty FROM schema_migrations;", "17|0")]:
    command = ["sqlite3", str(workspace / ".fid/fid.sqlite"), query]
    result = subprocess.run(command, capture_output=True, text=True)
    print(json.dumps({"command": command, "exit_code": result.returncode,
                      "stdout": result.stdout, "stderr": result.stderr}), flush=True)
    result.check_returncode()
    assert result.stdout.strip() == expected, (query, result.stdout)
print("PASS real CLI rebuild: read parity, unchanged journal, schema 17, FK/integrity clean")
