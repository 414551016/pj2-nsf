#!/usr/bin/env python3
"""Behavior checks using real HTTP and a restartable NSF process (Linux).

Usage: python3 scripts/test_api.py --binary build/nsf-integration
Build an isolated binary with -ldflags '-X nsf/internal/db.NSF_DB_PATH=/tmp/nsf-pj2-test.db'.
The script never removes a database file; generated service names are unique.
"""
import argparse
import json
import os
from pathlib import Path
import signal
import subprocess
import tempfile
import time
import urllib.error
import urllib.request
import uuid

parser = argparse.ArgumentParser()
parser.add_argument('--binary', default='build/nsf-integration')
parser.add_argument('--db', choices=['bbolt', 'mongodb'], default='bbolt')
parser.add_argument('--port', type=int, default=18081)
args = parser.parse_args()
binary = str(Path(args.binary).resolve())
base = f'http://127.0.0.1:{args.port}/sbi'
opener = urllib.request.build_opener(urllib.request.ProxyHandler({}))
process = None
checks = 0

def request(method, path, body=None):
    req = urllib.request.Request(base+path, method=method,
        data=None if body is None else json.dumps(body).encode(),
        headers={'Content-Type':'application/json'})
    try:
        response = opener.open(req, timeout=5)
    except urllib.error.HTTPError as error:
        response = error
    with response:
        raw = response.read()
        return response.status, json.loads(raw) if raw else None

def check(condition, label):
    global checks
    if not condition:
        raise AssertionError(label)
    checks += 1
    print(f'PASS {checks:02d}: {label}', flush=True)

with tempfile.TemporaryDirectory(prefix='nsf-api-') as temp:
    config = Path(temp)/'config.yaml'
    config.write_text(f'nsf:\n  sbi_address: "127.0.0.1"\n  sbi_port: {args.port}\n  db_type: "{args.db}"\nlogger:\n  level: "error"\n')
    log = open(Path(temp)/'server.log','w+')
    def start():
        global process
        process = subprocess.Popen([binary,'-c',str(config)],stdout=log,stderr=subprocess.STDOUT,env=os.environ.copy())
        for _ in range(50):
            if process.poll() is not None:
                log.flush();log.seek(0)
                raise RuntimeError('NSF exited: '+log.read())
            try:
                if request('GET','/heartbeat')[0] == 204:
                    return
            except (OSError, urllib.error.URLError):
                pass
            time.sleep(.1)
        raise RuntimeError('NSF failed to start')
    def stop(crash=False):
        global process
        if process is not None and process.poll() is None:
            process.send_signal(signal.SIGKILL if crash else signal.SIGINT)
            process.wait(timeout=10)
        process = None
    tag = uuid.uuid4().hex[:10]
    a,b = 'alpha-'+tag,'beta-'+tag
    nf1,nf2 = 'NF1-'+tag,'NF2-'+tag
    try:
        start()
        check(request('GET','/heartbeat')[0]==204,'heartbeat returns 204')
        check(request('GET','/nnsf-disc',{'service':'missing-'+tag})[0]==404,'unknown service returns 404')
        code,data=request('POST','/nnsf-nfm',{'nf_name':nf1,'sbi_address':'127.0.0.11','sbi_port':8101,'services':[a,b]})
        check(code==200 and bool(data.get('service_token')) and bool(data.get('service_uuid')),'register returns token and UUID')
        code,data=request('GET','/nnsf-disc',{'service':a})
        check(code==200 and data=={'service':a,'nf_name':nf1,'sbi_address':'127.0.0.11','sbi_port':8101},'discovery returns exact endpoint')
        code,data=request('GET','/nnsf-disc/group',{'services':[a,b]})
        check(code==200 and {x['service'] for x in data['services']}=={a,b} and all(x['nf_name']==nf1 for x in data['services']),'group discovery returns both services')
        stop(crash=True)
        start()
        code,data=request('GET','/nnsf-disc',{'service':a})
        check(code==200 and data['nf_name']==nf1,'registration survives abrupt process kill and restart')
        check(request('POST','/nnsf-nfm',{'nf_name':nf2,'sbi_address':'127.0.0.12','sbi_port':8102,'services':[a]})[0]==200,'second provider registers same service')
        check(request('DELETE','/nnsf-nfm',{'nf_name':nf1,'services':[a]})[0]==204,'partial deregistration returns 204')
        code,data=request('GET','/nnsf-disc',{'service':a})
        check(code==200 and data['nf_name']==nf2 and data['sbi_port']==8102,'remaining provider is discoverable')
        code,data=request('GET','/nnsf-disc',{'service':b})
        check(code==200 and data['nf_name']==nf1,'unrelated service is preserved')
        check(request('DELETE','/nnsf-nfm',{'nf_name':nf2,'services':[a]})[0]==204,'last provider deregisters')
        stop()
        start()
        check(request('GET','/nnsf-disc',{'service':a})[0]==404,'deletion survives graceful restart')
        check(request('GET','/nnsf-disc',{'service':b})[0]==200,'remaining registration survives graceful restart')
        check(request('DELETE','/nnsf-nfm',{'nf_name':nf1,'services':[b]})[0]==204,'test registration cleaned up')
        print(f'RESULT: {checks} PASS, 0 FAIL; backend={args.db}',flush=True)
    finally:
        stop()
        log.close()
