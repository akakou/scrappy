#!/usr/bin/env python3
import nativemessaging
import ctypes
import base64
import sys

nonce = "hoge".encode('utf-8')
basename = "hoge".encode('utf-8')

SECRET_KEY_PATH = "./ignored_workspace/member_private.bin".encode('utf-8')
CREDENTIAL_PATH = "./ignored_workspace/member_credential.bin".encode('utf-8')
MAX_SIZE = 512

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

    if "nonce" not in message:
        send_native_message("[error] nonce must needed")
        return

    if "domain" not in message:
        send_native_message("[error] domain must needed")
        return

    nonce = message["nonce"].encode('utf-8')
    domain = message["domain"].encode('utf-8')
    
    try:
        signature = sign(nonce, domain, SECRET_KEY_PATH, CREDENTIAL_PATH)
        send_native_message(signature)
    except Exception as e:
        send_native_message(f"[error] siging error : {str(e)}")


if __name__ == '__main__':
    print(sign(nonce, basename, SECRET_KEY_PATH, CREDENTIAL_PATH))

    while True:
        main_loop()
