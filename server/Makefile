FLAGS = `pkg-config ecdaa --libs --cflags`
# `pkg-config acml --libs --cflags`
INPUT = ./verify/*.c
SHARED_TARGET = ./bin

build:
	mkdir -p bin && gcc $(INPUT) -o ./bin/libverify.so $(FLAGS) -shared -fPIC  
test:
	mkdir -p bin && gcc $(INPUT) -o ./bin/verify $(FLAGS) && ./bin/verify
	
