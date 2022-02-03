from flask import Flask, request, render_template, session, make_response
import time
import sys
from datetime import datetime, timedelta
from model import AttestationLogForVerifier
import secrets
from c_lib import parse_k, verify

DOMAIN = "localhost"
MAX_COUNTER = 5

app = Flask(__name__)
name = 'taro'

app.permanent_session_lifetime = timedelta(minutes=5) 

app.secret_key = 'hogehogehoge'


def setup_nonce(body):
    session.permanent = True
    nonce = secrets.token_urlsafe(32)
    session["nonce"] = nonce
    resp = make_response(body)
    resp.set_cookie('nonce', nonce)
    return resp


def verify_attest():
    if "nonce" not in session:
        return False

    nonce = session["nonce"].encode('utf-8')
    now = datetime.now()
    now = now.replace(hour=now.hour, minute=now.minute, second=0, microsecond=0)

    signature = request.cookies.get('signature', None)
    print("signature: ", signature)

    counter = request.cookies.get('counter', None)
    print("counter: ", counter)

    if signature is None or counter is None :
        return False
    
    if int(counter) > MAX_COUNTER:
        return False

    signature = signature.encode('utf-8')
    basename = f"{DOMAIN}@{now}@{counter}".encode('utf-8')

    result = verify(signature, b'hello', basename)
    k = parse_k(signature)
    
    print("basename:", basename)

    print("result:", result)
    print("k:", k)

    attestation_log = AttestationLogForVerifier()
    attestation_log.k = k 
    attestation_log.basename = basename

    has_exist = AttestationLogForVerifier.exist(attestation_log) 
    AttestationLogForVerifier.save(attestation_log)
    
    print('has_exist: ', has_exist)

    return result and not has_exist

@app.route("/", methods=['GET'])
def get_index():
    body = render_template('index.html', name=name)
    return setup_nonce(body)

@app.route("/slow", methods=['GET'])
def post_index():
    if not verify_attest():
        return "error"

    # something heavy
    time.sleep(5)

    body = render_template('index.html', name=name)
    return setup_nonce(body)

if __name__ == "__main__":
    app.run(
        threaded=True,
        debug=True,
        host="0.0.0.0",
        port=5000)