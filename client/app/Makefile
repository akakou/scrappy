FLAGS = `pkg-config ecdaa --libs --cflags`
INPUT = ./sign/*.c
SHARED_TARGET = ./bin

build:
	mkdir -p bin && gcc $(INPUT) -o ./bin/libsign.so $(FLAGS) -shared -fPIC 
test_build: 
	mkdir -p bin && gcc $(INPUT) -o ./bin/sign $(FLAGS)
install: ./attest.py ./bin/sign
	cp ./attestation.json ~/.config/google-chrome/NativeMessagingHosts/com.akakou.attestation.json
