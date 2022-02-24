from flask import Flask, request, render_template, session, make_response
import time
import sys
from datetime import datetime, timedelta
from model import AttestationLogForVerifier
import secrets
import json

sys.path.append('/attestation/common/')
from wrapper import verify, parse_k


ORIGIN = "http://localhost:5000"

app = Flask(__name__)

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

    now = datetime.now()
    now = now.replace(hour=now.hour, minute=now.minute, second=now.second // 20, microsecond=0)
    now = now.timestamp()
    now = int(now)

    signature = attestation['signature']
    print("signature: ", signature)

    if signature is None:
        return False

    signature = signature.encode('utf-8')
    
    basename = f"{ORIGIN}@{now}".encode('utf-8')
    
    result = verify(signature, basename)
    k = parse_k(signature)
    
    print("basename:", basename)
    print("result:", result)
    return result

    print("k:", k)

    attestation_log = AttestationLogForVerifier()
    attestation_log.k = k 
    attestation_log.basename = basename

    has_exist = AttestationLogForVerifier.exist(attestation_log) 
    AttestationLogForVerifier.save(attestation_log)
    
    print('has_exist: ', has_exist) 

    return result and not has_exist

@app.route("/", methods=['GET'])
def index():
    nonce = make_nonce()
    return render_template('index.html', nonce=nonce)

@app.route("/fast", methods=['GET'])
def index_fast():
    return render_template('hello.html')


@app.route("/slow_without_attest", methods=['POST'])
def index_with_slow():
    # something heavy
    time.sleep(5)
    return render_template('hello.html')   

@app.route("/slow_with_attest", methods=['POST'])
def index_without_slow():
    if not verify_attest():
        return render_template('error.html')

    # something heavy
    time.sleep(5)

    return render_template('hello.html')

if __name__ == "__main__":
    app.run(
        threaded=False,
        debug=True,
        host="0.0.0.0",
        port=5000)