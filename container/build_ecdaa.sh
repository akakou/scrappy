### ecdaa
export ECDAA_CURVES=FP256BN
export CMAKE_PREFIX_PATH=/usr/local/
export LD_LIBRARY_PATH=/usr/local/lib/

git clone https://github.com/xaptum/amcl /attestation/attestation/thirdparty/amcl
mkdir -p /attestation/attestation/thirdparty/amcl/target/build
cd /attestation/attestation/thirdparty/amcl/target/build
cmake -D CMAKE_INSTALL_PREFIX=/opt/amcl ../..
make -j
make install


git clone https://github.com/xaptum/ecdaa/ /attestation/attestation/thirdparty/ecdaa
mkdir -p /attestation/attestation/thirdparty/ecdaa/build
cd /attestation/attestation/thirdparty/ecdaa/build
../.travis/install-amcl.sh ./amcl ${CMAKE_PREFIX_PATH} ${ECDAA_CURVES}
cmake .. -DCMAKE_BUILD_TYPE=Release -DECDAA_CURVES=FP256BN
cmake --build . --target install