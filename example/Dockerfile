FROM ubuntu:latest

RUN apt-get update && apt-get install -y sqlite3 gcc wget
RUN wget https://go.dev/dl/go1.19.1.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.19.1.linux-amd64.tar.gz

ENV GOPATH /root/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"