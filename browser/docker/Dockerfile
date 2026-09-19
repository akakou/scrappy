FROM --platform=linux/amd64 ubuntu

RUN export DEBIAN_FRONTEND=noninteractive && apt-get update && apt-get install -y wget tzdata gnupg python3 cmake pkg-config m4 libtool automake autoconf libtss2-dev tpm2-tools git lsb-release sqlite3

RUN wget https://go.dev/dl/go1.19.1.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.19.1.linux-amd64.tar.gz

ADD build_chrome.sh build_chrome.sh
RUN sh build_chrome.sh

# ADD /dev/tpm0
# RUN chmod 777 /dev/

USER scrappy
WORKDIR /scrappy/