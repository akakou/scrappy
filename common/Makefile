FLAGS = `pkg-config ecdaa --libs --cflags` `pkg-config ecdaa-tpm --libs --cflags` `pkg-config tss2-tcti-device --libs` `pkg-config tss2-sys --libs --cflags`
INPUT = ./*.h ./*.c
SHARED_TARGET = ./bin

build:
	mkdir -p bin && gcc $(INPUT) -o ./bin/libattest.so $(FLAGS) -shared -fPIC 
test: 
	mkdir -p bin && gcc $(INPUT) -o ./bin/test_attest $(FLAGS) && ./bin/test_attest 
install: ./attest.py ./bin/sign
	cp ./attestation.json ~/.config/google-chrome/NativeMessagingHosts/com.akakou.attestation.json
clean:
	rm -rf ./bin