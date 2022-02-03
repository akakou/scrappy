import ctypes
import base64


GPK_PATH = "/attestation/ignored_workspace/group_public.bin".encode('utf-8')
SECRET_KEY_PATH = "/attestation/ignored_workspace/member_private.bin".encode('utf-8')
CREDENTIAL_PATH = "/attestation/ignored_workspace/member_credential.bin".encode('utf-8')

SIMPLE_SIG_PATH = "/attestation/ignored_workspace/encoded_sig.bin".encode('utf-8')
SIG_SIZE = 421
K_SIZE = 65

MAX_COUNTER = 5
TERM_LEN = 50

char_ptr = ctypes.POINTER(ctypes.c_char)

lib = ctypes.CDLL("/attestation/common/bin/libattest.so")

def libsign():
    lib.sign.argtypes = (
            char_ptr,
            char_ptr,
            ctypes.c_int,
            char_ptr,
            ctypes.c_int,
            char_ptr,
            char_ptr)
    lib.sign.restype = ctypes.c_int

    def sign(message, basename, secret_key_path=SECRET_KEY_PATH, credential_path=CREDENTIAL_PATH):
        result = ctypes.create_string_buffer(SIG_SIZE)
        lib.sign(result, message, len(message), basename, len(basename), SECRET_KEY_PATH, CREDENTIAL_PATH)
        encoded = base64.b64encode(result.raw).decode('utf-8')
        return encoded
    
    return sign

def libverify():
    lib.verify.argtypes = (
            char_ptr,
            char_ptr,
            ctypes.c_int,
            char_ptr,
            ctypes.c_int,
            char_ptr)
    lib.verify.restype = ctypes.c_int

    def verify(signature, message, basename, gpk_path=GPK_PATH):
        decoded = base64.b64decode(signature)

        status = lib.verify(decoded, message, len(message), basename, len(basename), gpk_path)
        status = int(status)
        
        if status == 0:
            return True
        elif status == 1:
            return False
        else:
            raise Exception(f"Error verify ({status})")

    
    return verify

def libparse_k():
    lib.parse_k.argtypes = (
            char_ptr,
            char_ptr)
    lib.parse_k.restype = ctypes.c_int

    def parse_k(signature):
        result = ctypes.create_string_buffer(K_SIZE)
        decoded = base64.b64decode(signature)

        lib.parse_k(result, decoded)
        
        encoded = base64.b64encode(result.raw).decode('utf-8')
        return encoded
    
    return parse_k

sign = libsign()
parse_k = libparse_k()
verify = libverify()

if __name__ == '__main__':
    message = b"hogehoge"
    basename = b"hogehoge"

    sig = sign(message, basename)
    
    print("sig: ", sig)

    k = parse_k(sig)
    print("k: ", k)

    result = verify(sig, message, basename)
    assert result

    result = verify(sig, b'piyopiyo', basename)
    assert not result

    result = verify(sig, message, b'piyopiyo')
    assert not result

