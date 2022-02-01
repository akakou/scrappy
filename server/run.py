from flask import Flask, request, render_template
import time
import sys

from c_lib import parse_k, verify

app = Flask(__name__)
name = 'taro'

def verify_attest():
    attestation = request.cookies.get('attestation', None)
    print("attestation", attestation)

    if attestation is None:
        return False
    
    encoded = attestation.encode('utf-8')

    result = verify(encoded, b"hogehoge", b"hogehoge")
    k = parse_k(encoded)
    print("result:", result)
    print("k:", k)

    return result

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
        threaded=False,
        debug=True,
        host="0.0.0.0",
        port=5000)