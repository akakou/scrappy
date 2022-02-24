FROM --platform=linux/amd64 ubuntu

RUN apt-get update
RUN export DEBIAN_FRONTEND=noninteractive && apt-get install -y wget tzdata gnupg python3 cmake pkg-config m4 libtool automake autoconf libtss2-dev tpm2-tools git lsb-release 
ADD build_chrome.sh /attestation/build_chrome.sh
RUN sh /attestation/build_chrome.sh
ADD build_ecdaa.sh /attestation/build_ecdaa.sh
RUN sh /attestation/build_ecdaa.sh
RUN pip3 install nativemessaging sqlalchemy flask
RUN apt-get install -y sqlite3
# ADD /dev/tpm0
# RUN chmod 777 /dev/tpm0
USER attestation


