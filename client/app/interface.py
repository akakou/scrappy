#!/usr/bin/env python3
import ctypes
import base64
import sys

sys.path.append('/attestation/common/')
from wrapper import libsign
sign = libsign()

basename = sys.argv[1].encode('utf-8')
res = sign(basename)

print(res)