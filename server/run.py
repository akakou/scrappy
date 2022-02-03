from flask import Flask, request, render_template, session, make_response
import time
import sys
from datetime import datetime, timedelta
from model import AttestationLogForVerifier
import secrets
import json

sys.path.append('/attestation/common/')
from wrapper import verify, parse_k, MAX_COUNTER


ORIGIN = "http://localhost:5000"

app = Flask(__name__)
name = 'taro'

app.permanent_session_lifetime = timedelta(minutes=5) 
app.secret_key = 'hogehogehoge'


def make_nonce():
    session.permanent = True
    nonce = secrets.token_urlsafe(32)
    session["nonce"] = nonce

    return nonce


def verify_attest():
    if "attestation" not in request.form:
        return False
    
    attestation = request.form["attestation"]
    if not len(attestation):
        return False
    
    attestation = json.loads(attestation)

    if "nonce" not in session:
        return False

    nonce = session["nonce"].encode('utf-8')
    # print("nonce:", nonce)

    now = datetime.now()
    now = now.replace(hour=now.hour, minute=now.minute, second=0, microsecond=0)

    signature = attestation['signature']
    print("signature: ", signature)

    counter = attestation['counter']
    print("counter: ", counter)

    if signature is None or counter is None :
        return False
    
    if int(counter) > MAX_COUNTER:
        return False

    signature = signature.encode('utf-8')
    basename = f"{ORIGIN}@{now}@{counter}".encode('utf-8')

    result = verify(signature, nonce, basename)
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

    return result and not False

@app.route("/", methods=['GET'])
def index():
    nonce = make_nonce()
    body = render_template('index.html', nonce=nonce)
    return body

@app.route("/fast", methods=['GET'])
def index_fast():
    body = render_template('hello.html')
    return body


@app.route("/slow", methods=['POST'])
def index_slow():
    if not verify_attest():
        return "error"

    # something heavy
    time.sleep(5)

    nonce = make_nonce()
    return render_template('hello.html')

if __name__ == "__main__":
    app.run(
        threaded=True,
        debug=True,
        host="0.0.0.0",
        port=5000)