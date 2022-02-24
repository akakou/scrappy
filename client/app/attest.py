#!/usr/bin/env python3
import ctypes
import base64
import sys

import json

import subprocess
import nativemessaging
from model import AttestationLog

message = b"hogehoge"
basename = b"hogehoge"

# todo: escape
def sign(message, basename):
    return subprocess.getoutput(f'python3 /attestation/client/app/interface.py {message} {basename}')


send_native_message = lambda x: nativemessaging.send_message(nativemessaging.encode_message(x))

def main_loop():
    message = nativemessaging.get_message()
    
    if "nonce" in message:
        nonce = message['nonce'].encode('utf-8')
    else:
        send_native_message("[error] nonce must needed")
        return 

    if "origin" in message:
        origin = message['origin']
    else:
        send_native_message("[error] origin must needed")
        return
    
    try:
        attestation_log = AttestationLog.gen_attestation_log(origin)
    except Exception as e:
        send_native_message(f"[error] siging error : {str(e)}")
        return


    if attestation_log is None:
        send_native_message(f"[error] have generate too many attestations")
        return

    try:
        signature = sign(nonce, str(attestation_log).encode('utf-8'))
        AttestationLog.save_attestation_log(attestation_log)

        send_native_message({'signature': signature, 'counter': attestation_log.counter, 'basename': str(attestation_log)})
    except Exception as e:
        send_native_message(f"[error] siging error : {str(e)}")
        return


if __name__ == '__main__':
    # message = b"hogehoge"
    # basename = b"hogehoge"

    # sig = sign(message, basename)
    
    # print("sig: ", sig)

    while True:
        try:
            main_loop()
        except Exception as e:
            send_native_message(f"[error] other error : {str(e)}")

#