import ctypes
import base64

char_ptr = ctypes.POINTER(ctypes.c_char)
lib = ctypes.CDLL("./bin/libverify.so")

GPK_PATH = "/attestation/ignored_workspace/group_public.bin".encode('utf-8')
SIMPLE_SIG_PATH = "/attestation/ignored_workspace/encoded_sig.bin".encode('utf-8')

MAX_SIZE = 65


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
        result = ctypes.create_string_buffer(MAX_SIZE)
        decoded = base64.b64decode(signature)

        lib.parse_k(result, decoded)
        
        encoded = base64.b64encode(result.raw).decode('utf-8')
        return encoded
    
    return parse_k

parse_k = libparse_k()
verify = libverify()

if __name__ == '__main__':
    message = b"hogehoge"
    basename = b"hogehoge"

    with open(SIMPLE_SIG_PATH) as f:
        sig = f.read().encode('utf-8')
    
    print("sig", sig)

    k = parse_k(sig)
    print("k", k)

    result = verify(sig, message, basename)

    print(result)