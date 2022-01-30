### ecdaa
export DIST=$(lsb_release -cs)
export ECDAA_CURVES=FP256BN
export CMAKE_PREFIX_PATH=/home/attestation/ecdaa/build/deps

git clone https://github.com/xaptum/ecdaa/ /home/attestation/ecdaa/

mkdir -p /home/attestation/ecdaa/build/deps
cd /home/attestation/ecdaa/build

../.travis/install-amcl.sh ./amcl ./deps FP256BN

cmake .. -DCMAKE_BUILD_TYPE=Release -DECDAA_CURVES=FP256BN
make -j
make install