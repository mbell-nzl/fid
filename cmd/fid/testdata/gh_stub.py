#!/usr/bin/python3
"""Exact, read-only gh responses. Each test supplies only its child's env."""
import json
import os
import subprocess
import sys

with open(os.environ['FID_TEST_GH_RESPONSES'], encoding='utf-8') as stream:
    fixture = json.load(stream)
key = ' '.join(sys.argv[1:])
if key not in fixture['requests']:
    sys.exit('unknown gh request: ' + key)
with open(os.environ['FID_TEST_GH_LOG'], 'a', encoding='utf-8') as stream:
    stream.write(key + '\n')
if key.startswith('pr view ') and fixture.get('action'):
    subprocess.run(fixture['action'], check=True, stdout=subprocess.DEVNULL)
response = fixture['requests'][key]
if 'error' in response:
    sys.exit(response['error'])
print(json.dumps(response['body']))
