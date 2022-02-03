FLAGS = `pkg-config ecdaa --libs --cflags`
INPUT = ./sign/*.c
SHARED_TARGET = ./bin

build:
	mkdir -p bin && gcc $(INPUT) -o ./bin/libsign.so $(FLAGS) -shared -fPIC 
test: 
	mkdir -p bin && gcc $(INPUT) -o ./bin/sign -DMAIN_IS_SIGN=1 $(FLAGS) && ./bin/sign 
install: ./attest.py ./bin/sign
	cp ./attestation.json ~/.config/google-chrome/NativeMessagingHosts/com.akakou.attestation.json
