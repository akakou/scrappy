FLAGS = `pkg-config ecdaa --libs --cflags`
INPUT = ./verify/*.c
SHARED_TARGET = ./bin

build:
	mkdir -p bin && gcc $(INPUT) -o ./bin/libverify.so $(FLAGS) -shared -fPIC 
test:
	/attestation/client/app/bin/sign | ./bin/verify
test_build: 
	mkdir -p bin && gcc $(INPUT) -o ./bin/verify $(FLAGS)
