REPOSITORY_URL = git@github.com:akakou/ecdaa.git

all: ./ecdaa64 ./ecdaa32
	echo "Done"

./ecdaa64:
	git clone ${REPOSITORY_URL} ecdaa64 \
	&& cd ecdaa64/thirdparty && export BITS=64 && make

./ecdaa32:
	git clone ${REPOSITORY_URL} ecdaa32 \
	&& cd ecdaa32/thirdparty && export BITS=32 && make
 
clean:
	rm -rf ./ecdaa64 ./ecdaa32
