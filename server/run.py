from flask import Flask, request, render_template
import time
import sys
from datetime import datetime
from model import AttestationLogForVerifier

from c_lib import parse_k, verify

DOMAIN = "localhost"
MAX_COUNTER = 5

app = Flask(__name__)
name = 'taro'

def verify_attest():
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

    result = verify(signature, b"hoge", basename)
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
    return render_template('index.html', name=name)

@app.route("/slow", methods=['GET'])
def post_index():
    if not verify_attest():
        return "error"

    # something heavy
    time.sleep(5)

    return render_template('index.html', name=name)

if __name__ == "__main__":
    app.run(
        threaded=True,
        debug=True,
        host="0.0.0.0",
        port=5000)