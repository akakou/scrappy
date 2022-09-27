FROM --platform=linux/amd64 ubuntu

RUN apt-get update
RUN export DEBIAN_FRONTEND=noninteractive && apt-get install -y wget tzdata gnupg python3 cmake pkg-config m4 libtool automake autoconf libtss2-dev tpm2-tools git lsb-release 
WORKDIR /attestation/
ADD build_chrome.sh build_chrome.sh

# ADD /dev/tpm0
# RUN chmod 777 /dev/tpm0
USER attestation
