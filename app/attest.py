#!/usr/bin/env python3
import nativemessaging

while True:
    message = nativemessaging.get_message()
    nativemessaging.send_message(nativemessaging.encode_message("world"))