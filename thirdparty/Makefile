REPOSITORY_URL = https://github.com/akakou/ecdaa

all: ./ecdaa64 ./ecdaa32
	echo "Done"

./ecdaa64:
	git clone ${REPOSITORY_URL} ecdaa64 \
	&& cd ecdaa64/thirdparty && make BITS=64

./ecdaa32:
	git clone ${REPOSITORY_URL} ecdaa32 \
	&& cd ecdaa32/thirdparty && make BITS=32
 
clean:
	rm -rf ./ecdaa64 ./ecdaa32
