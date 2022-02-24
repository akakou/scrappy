#!/usr/bin/env python3
import ctypes
import base64
import sys

sys.path.append('/attestation/common/')
from wrapper import libsign
sign = libsign()

message = sys.argv[1].encode('utf-8')
basename = sys.argv[2].encode('utf-8')

res = sign(message, basename)

print(res)