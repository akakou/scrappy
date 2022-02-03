#!/usr/bin/env python3
import nativemessaging
import ctypes
import base64
import sys

from model import AttestationLog

nonce = "hoge".encode('utf-8')
basename = "hoge".encode('utf-8')

SECRET_KEY_PATH = "/attestation/ignored_workspace/member_private.bin".encode('utf-8')
CREDENTIAL_PATH = "/attestation/ignored_workspace/member_credential.bin".encode('utf-8')
MAX_SIZE = 421

char_ptr = ctypes.POINTER(ctypes.c_char)
send_native_message = lambda x: nativemessaging.send_message(nativemessaging.encode_message(x))


def libsign():
    lib = ctypes.CDLL("./bin/libsign.so")

    lib.sign.argtypes = (
            char_ptr,
            char_ptr,
            ctypes.c_int,
            char_ptr,
            ctypes.c_int,
            char_ptr,
            char_ptr)
    lib.sign.restype = ctypes.c_int

    def sign(message, basename, secret_key_path, credential_path):
        result = ctypes.create_string_buffer(MAX_SIZE)
        lib.sign(result, message, len(message), basename, len(basename), SECRET_KEY_PATH, CREDENTIAL_PATH)
        encoded = base64.b64encode(result.raw).decode('utf-8')
        return encoded
    
    return sign

sign = libsign()


def main_loop():
    message = nativemessaging.get_message()
    
    if "nonce" in message:
        nonce = message['nonce']
    else:
        send_native_message("[error] nonce must needed")
        return 

    if "domain" in message:
        domain = message['domain']
    else:
        send_native_message("[error] domain must needed")
        return
    
    
    try:
        attestation_log = AttestationLog.gen_attestation_log(domain)
    except Exception as e:
        send_native_message(f"[error] siging error : {str(e)}")
        return
    

    nonce = nonce.encode('utf-8')
    domain = domain.encode('utf-8')


    if attestation_log is None:
        send_native_message(f"[error] have generate too many attestations")
        return

    try:
        signature = sign(nonce, domain, SECRET_KEY_PATH, CREDENTIAL_PATH)
        AttestationLog.save_attestation_log(attestation_log)

        send_native_message({'signature': signature, 'counter': attestation_log.counter})
    except Exception as e:
        send_native_message(f"[error] siging error : {str(e)}")
        return


if __name__ == '__main__':
    print(sign(nonce, basename, SECRET_KEY_PATH, CREDENTIAL_PATH))

    while True:
        try:
            main_loop()
        except Exception as e:
            send_native_message(f"[error] other error : {str(e)}")

