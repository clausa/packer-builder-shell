FROM golang:1.3.3
MAINTAINER csa@csa-net.dk

RUN go get -u github.com/mitchellh/gox
RUN go get -u github.com/clausa/packer-builder-shell
WORKDIR /go/bin
RUN go get -u -v ./...

