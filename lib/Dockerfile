FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
    vim \
    curl \
    wget \
    unzip \
    python3 \
    build-essential

ENV RUST_HOME /usr/local/lib/rust
ENV RUSTUP_HOME ${RUST_HOME}/rustup
ENV CARGO_HOME ${RUST_HOME}/cargo
RUN mkdir /usr/local/lib/rust && \
    chmod 0755 $RUST_HOME
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs > ${RUST_HOME}/rustup.sh \
    && chmod +x ${RUST_HOME}/rustup.sh \
    && ${RUST_HOME}/rustup.sh -y --default-toolchain nightly --no-modify-path
ENV PATH $PATH:$CARGO_HOME/bin

RUN rustup target add armv7-linux-androideabi
RUN rustup target add aarch64-linux-android
RUN rustup target add thumbv7neon-linux-androideabi
RUN rustup target add i686-linux-android
RUN rustup target add x86_64-linux-android

RUN wget https://dl.google.com/android/repository/android-ndk-r25c-linux.zip
RUN unzip android-ndk-r25c-linux.zip

RUN /android-ndk-r25c/build/tools/make_standalone_toolchain.py --api 33 --arch arm64 --install-dir /ndk-bin-build/arm64
RUN /android-ndk-r25c/build/tools/make_standalone_toolchain.py --api 33 --arch arm --install-dir /ndk-bin-build/arm
RUN /android-ndk-r25c/build/tools/make_standalone_toolchain.py --api 33 --arch x86 --install-dir /ndk-bin-build/x86
